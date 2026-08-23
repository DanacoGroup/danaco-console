// Odpowiedzialność pliku: moduł Terminal — wypełnienie portu Terminal czterema
// komendami obszaru `terminal.*`. Powłoki leżą w
// `adapter_modul_terminal_powloki.go`, tryb uprawnień w
// `adapter_modul_terminal_uprawnienia.go`, ewidencja procesów w
// `adapter_modul_terminal_rejestr.go`, strumień wyjścia w
// `adapter_modul_terminal_strumien.go`, a uruchomienie procesu
// w `adapter_modul_terminal_wykonanie.go`.
//
// Przed każdym procesem stoją trzy sprawdzenia, w tej kolejności:
//  1. tryb uprawnień okna — czy wolno w ogóle uruchomić proces i czyim
//     poleceniem (PermissionMode, `adapter_modul_terminal_uprawnienia.go`);
//  2. egzekutor izolacji — czy polecenie mieści się w obszarze okna
//     (session.SprawdzPolecenie nad zasadami z `izolacja.go`);
//  3. uruchamiacz warstwy kanału — jedyna droga startu procesu w drzewie
//     (port session.Uruchamiacz).
//
// Rdzeń nie buduje własnego `exec.Cmd` ani drugiego egzekutora — korzysta
// wyłącznie z tych dwóch warstw.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Zgodność adaptera z portem sprawdzana jest przy kompilacji.
var _ Terminal = (*adapterTerminala)(nil)

// przedrostekKarty znakuje identyfikator karty powłoki nadany przez rdzeń.
const przedrostekKarty = "term-"

// przedrostekProcesuTerminala znakuje identyfikator procesu terminala.
// Odrębny od `proc-` telemetrii: tamten opisuje bieg tury modelu, ten —
// proces urządzenia, i mylenie ich w Process Monitorze byłoby kosztowne.
const przedrostekProcesuTerminala = "tproc-"

// adapterTerminala wypełnia port Terminal.
type adapterTerminala struct {
	repozytorium dane.RepozytoriumTerminala
	rejestr      *rejestrTerminala
	// tunele i obserwacje trzymają wyposażenie DŁUGOŻYJĄCE tego biegu rdzenia.
	// Wiersz w bazie mówi, że tunel albo obserwacja są; uchwyt tutaj jest
	// jedynym, czym da się je zatrzymać — tak samo jak przy procesach.
	tunele     *rejestrTuneli
	obserwacje *rejestrObserwacji
	// okna daje tryb uprawnień okna i jego sesję. Bez niego moduł nie ruszy
	// ani jednego procesu: uruchomienie bez znajomości trybu uprawnień byłoby
	// obejściem bramy, a nie pracą w warunkach niepełnych danych.
	okna *session.Rejestr
	// uruchamiacz jest portem warstwy kanału — jedyną drogą startu procesu.
	uruchamiacz session.Uruchamiacz
	// rozstrzygacz i katalog składają zasady izolacji obowiązujące w oknie.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// katalogDanych jest katalogiem danych rdzenia. Moduł kładzie w nim klucze
	// SSH, które sam wytworzył (`adapter_modul_terminal_klucze.go`) — poza
	// konfiguracją OpenSSH Operatora, żeby produkt nie zmieniał jej sam.
	katalogDanych string
	// wyjscie rozsyła fragmenty strumienia do okna Output Console.
	wyjscie *nadawcaWyjscia
	// zmiana rozgłasza `terminal.process.changed`. Podpina ją obsługiwacz.
	zmiana func(shared.ChangeKind, shared.TerminalProcess)
}

// nowyAdapterTerminala wiąże port z rejestrem okien i uruchamiaczem procesów.
func nowyAdapterTerminala(okna *session.Rejestr, uruchamiacz session.Uruchamiacz) *adapterTerminala {
	return &adapterTerminala{
		rejestr:     nowyRejestrTerminala(),
		tunele:      nowyRejestrTuneli(),
		obserwacje:  nowyRejestrObserwacji(),
		okna:        okna,
		uruchamiacz: uruchamiacz,
	}
}

// ZTrwaloscia podpina dziennik kart i procesów. Bez niego moduł
// pracuje w pamięci jednego biegu rdzenia, a Process Monitor pokazuje
// wyłącznie procesy czynne.
func (a *adapterTerminala) ZTrwaloscia(repozytorium dane.RepozytoriumTerminala) *adapterTerminala {
	a.repozytorium = repozytorium
	return a
}

// ZIzolacja podpina rozstrzygacz zasięgu i ustalacz katalogu roboczego —
// dwa źródła, z których powstaje obszar okna egzekwowany przy uruchomieniu.
func (a *adapterTerminala) ZIzolacja(rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterTerminala {
	a.rozstrzygacz, a.katalog = rozstrzygacz, katalog
	return a
}

// ZKatalogiemDanych wskazuje katalog danych rdzenia. Bez niego moduł pracuje
// dalej, lecz nie wytworzy klucza SSH: nie miałby gdzie go położyć, a klucz
// odłożony w katalogu tymczasowym byłby kluczem znikającym przy restarcie.
func (a *adapterTerminala) ZKatalogiemDanych(katalog string) *adapterTerminala {
	a.katalogDanych = katalog
	return a
}

// ZWyjsciem podpina nadajnik strumienia wyjścia procesów.
func (a *adapterTerminala) ZWyjsciem(nadajnik Nadajnik) *adapterTerminala {
	a.wyjscie = nowyNadawcaWyjscia(nadajnik)
	return a
}

// Przygotuj odtwarza karty poprzedniego biegu rdzenia i osierocą procesy, do
// których rdzeń stracił uchwyt. Wywołuje się to raz, przy montażu.
//
// Bez tego kroku po restarcie serwera Operator widziałby w Process Monitorze
// procesy oznaczone jako czynne, których nikt już nie prowadzi, a karta
// widoczna na ekranie odpowiadałaby `not_found` na pierwsze polecenie.
func (a *adapterTerminala) Przygotuj(ctx context.Context) error {
	if a.repozytorium == nil {
		return nil
	}
	if _, err := a.repozytorium.OsierociProcesy(ctx); err != nil {
		return err
	}
	karty, err := a.repozytorium.Karty(ctx)
	if err != nil {
		return err
	}
	for _, wiersz := range karty {
		if wiersz.Stan == shared.TerminalSessionStatusExited {
			continue
		}
		karta := &kartaTerminala{
			kod:        wiersz.Kod,
			oknoKod:    wiersz.OknoKod,
			idSesji:    a.sesjaOkna(wiersz.OknoKod),
			powloka:    wiersz.Powloka,
			tytul:      wartoscTekstu(wiersz.Tytul),
			katalog:    wartoscTekstu(wiersz.KatalogRoboczy),
			srodowisko: map[string]string{},
			stan:       shared.TerminalSessionStatusRunning,
			utworzono:  time.UnixMilli(chwilaBazy(wiersz.Utworzono)).UTC(),
			celZdalny:  wiersz.CelZdalny,
			hostKod:    wartoscTekstu(wiersz.HostKod),
		}
		if wiersz.PortZdalny != nil {
			karta.portZdalny = int(*wiersz.PortZdalny)
		}
		a.rejestr.ZapiszKarte(karta)
	}
	// Tunele i obserwacje poprzedniego biegu nie biegną: rdzeń stracił do nich
	// uchwyt tak samo jak do procesów. Wykazywanie ich jako czynnych byłoby
	// nieprawdą, więc schodzą na stany końcowe wraz z powodem.
	if _, err := a.repozytorium.OsierocTunele(ctx); err != nil {
		return err
	}
	if _, err := a.repozytorium.OsierocObserwacje(ctx); err != nil {
		return err
	}
	return nil
}

// Zamknij kończy procesy terminala czynne w chwili zatrzymania rdzenia.
func (a *adapterTerminala) Zamknij() {
	if a == nil {
		return
	}
	a.rejestr.Zamknij()
	// Tunel i obserwacja przeżywają każde żądanie, więc muszą zginąć wraz
	// z rdzeniem: proces `ssh` zostawiony bez uchwytu przenosiłby bajty, o których
	// nikt już nie wie, a pętla przeglądu uruchamiałaby polecenia po zatrzymaniu.
	a.tunele.zamknijWszystkie()
	a.obserwacje.zatrzymajWszystkie()
}

// OtworzKarte obsługuje `terminal.session.open` — zakłada kartę powłoki.
//
// Karta nie startuje procesu, więc tryb uprawnień sprawdza się dopiero przy
// poleceniu. Sprawdzamy natomiast okno: karta bez okna nie ma ani trybu
// uprawnień, ani obszaru izolacji, a więc nie da się jej później wykonać.
func (a *adapterTerminala) OtworzKarte(ctx context.Context,
	z shared.TerminalSessionOpenRequest) (shared.TerminalSessionOpenResponse, error) {

	oknoKod := strings.TrimSpace(z.WindowId)
	if oknoKod == "" {
		return shared.TerminalSessionOpenResponse{}, bladZadaniaTerminala("karta powłoki wymaga wskazania okna")
	}
	if !CzyPowlokaZnana(z.Shell) {
		return shared.TerminalSessionOpenResponse{}, bladZadaniaTerminala(
			"powłoka " + string(z.Shell) + " nie należy do słownika kontraktu")
	}
	okno, err := a.oknoWykonania(oknoKod)
	if err != nil {
		return shared.TerminalSessionOpenResponse{}, err
	}
	srodowisko, err := zmienneKarty(z.Environment)
	if err != nil {
		return shared.TerminalSessionOpenResponse{}, bladZadaniaTerminala(err.Error())
	}

	karta := &kartaTerminala{
		kod:        nowyIdentyfikator(przedrostekKarty),
		oknoKod:    oknoKod,
		idSesji:    okno.IdSesji,
		powloka:    z.Shell,
		tytul:      wartoscTekstu(z.Title),
		katalog:    strings.TrimSpace(wartoscTekstu(z.WorkingDir)),
		srodowisko: srodowisko,
		stan:       shared.TerminalSessionStatusRunning,
		utworzono:  time.Now().UTC(),
	}
	if err := a.wskazCelZdalny(ctx, karta, z); err != nil {
		return shared.TerminalSessionOpenResponse{}, err
	}
	wskazanieUrzadzenia(karta, z)
	if karta.katalog == "" {
		karta.katalog = a.obszarOkna(okno).KatalogRoboczy
	}
	a.rejestr.ZapiszKarte(karta)
	a.zapiszKarte(ctx, karta)

	return shared.TerminalSessionOpenResponse{Session: kartaKontraktu(karta)}, nil
}
