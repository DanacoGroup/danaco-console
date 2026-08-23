// Odpowiedzialność pliku: MIEJSCE, W KTÓRYM ZESTAW NARZĘDZI TURY PRZESTAJE BYĆ
// STAŁY. Tu spotykają się dwa źródła — podstawa wnoszona przez definicję
// eksperta i doraźne dołożenia sesji — i tu wchodzą do wpisu `danaco`
// konfiguracji MCP okna, czyli do wiersza uruchomienia serwera narzędzi.
//
// ── DLACZEGO OPIS JEDZIE ARGUMENTEM, A NIE GOTOWĄ LISTĄ NAZW ────────────────
// Rdzeń nie zna nazw narzędzi eksperta i znać ich nie ma. Przełożenie kodu
// z definicji eksperta na pozycje wykazu (albo na całą ich grupę) należy do
// serwera narzędzi — `narzedzia/ekspert_wykaz.go` — bo tylko on trzyma wykaz
// kontraktu. Gdyby rdzeń wyliczał listę sam, powstałby drugi wykaz obok
// kontraktowego. Stąd wychodzi więc wskazanie: kod eksperta i nazwy
// dołożeń; rozwinięcie wskazania w wykaz robi strona przeciwna.
//
// ── DLACZEGO PRZEPISANIEM WPISU, A NIE DRUGIM SKŁADACZEM ──────────
// Konfigurację MCP składa `mostyOkna.tekstZNarzedziami` — jedyna droga, którą
// jedzie rozmowa. Ten kod jej nie powtarza: bierze jej wynik i dokłada do
// gotowego wpisu `danaco` argumenty zestawu. To ten sam wzorzec, którym okno
// asystenta dokłada sobie zasięg klawiatury (`zZasiegiemKlawiatury`
// w `adapter_modul_asystent_sterowanie.go`). Tekst nieczytelny albo wpis
// `danaco` nieobecny (brak binarium serwera narzędzi — `most_narzedzi.go`
// melduje ten powód osobno) zostaje nietknięty: tura idzie z zestawem, jaki
// jest, zamiast paść na składaniu konfiguracji.
//
// ── ARGUMENT NIEZNANY NIE JEST POMINIĘTY, TYLKO ZABÓJCZY ────────────────────
// Powód, dla którego dokładanie argumentów jest tu wybiórcze, a nie hurtowe,
// i dla którego każdy nowy argument wymaga umowy, a nie założenia.
// `cmd/danaco-narzedzia/main.go` woła `flag.Parse()` na domyślnym
// `flag.CommandLine`, czyli z `ExitOnError`: przełącznik, którego tamta strona
// nie zna, kończy proces serwera narzędzi. Skutek nie jest wtedy „tura bez kilku
// narzędzi", tylko „tura bez wykazu w całości" — i wygląda dla Operatora jak
// awaria modelu, nie jak rozjazd dwóch pakietów.
//
// Oba argumenty składane niżej mają odczyt po drugiej stronie:
//
//	`--zasieg ekspert --ekspert <kod>` — `narzedzia/zasieg_eksperta.go`,
//	    `ekspert_definicja.go`, `ekspert_wykaz.go`.
//	`--dolozenia <lista>`             — `narzedzia.RozbijDolozenia`
//	    i `WykazEksperta.ZDolozeniami`; wartość składa `injection.PrzelacznikDolozen`,
//	    czyli jedna definicja nazwy po obu stronach.
//
// ── PORT NIEWPIĘTY MELDUJE, ZAMIAST MILCZEĆ ───────────────────────
// Brak źródła dołożeń NIE odbiera modelowi narzędzi: tura jedzie podstawą, jaką
// ma, i wpisuje do dziennika rdzenia POWÓD. Meldunek idzie raz na powód, nie raz
// na turę — tak samo jak meldunek o braku binarium serwera narzędzi
// (`most_narzedzi.go`) i z tego samego powodu: tur bywa kilkaset dziennie,
// a powód się między nimi nie zmienia.
package core

import (
	"context"
	"encoding/json"

	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/narzedzia"
	"danacoconsole/server/internal/session"
)

// DolozeniaNarzedziSesji podaje DORAŹNE dołożenia — narzędzia dorzucone do tury
// przez Operatora w obrębie jednej sesji, poza definicją eksperta.
//
// Port jest wąski z rozmysłu: pyta o jedno i oddaje jedno. Sesja
// wskazana jest identyfikatorem KONTRAKTOWYM, tym samym, który niesie
// `session.Okno.IdSesji` — przekład na klucz wiersza należy do strony
// odpowiadającej, tak jak przy każdym innym repozytorium.
//
// Dołożenia są listą jednorodną i nie niosą znacznika zawężenia, bo dołożenie
// z natury DOKŁADA: nie potrafi zawęzić czegoś, czego samo nie ustanowiło.
// Sesja bez dołożeń oddaje listę pustą i nie jest to ani błąd, ani brak.
type DolozeniaNarzedziSesji interface {
	DolozeniaNarzedzi(ctx context.Context, idSesji string) (nazwy []string, err error)
}

// ZDolozeniamiSesji wpina źródło doraźnych dołożeń.
//
// Adapter bez tego portu prowadzi turę samą podstawą eksperta i melduje powód,
// gdy okno ma podstawę, do której dołożenia mogłyby dojść — cisza w tym miejscu
// byłaby nie do odróżnienia od sesji, w której Operator nic nie dołożył
// .
func (a *adapterRozmowy) ZDolozeniamiSesji(d DolozeniaNarzedziSesji) *adapterRozmowy {
	a.dolozeniaSesji = d
	return a
}

// zestawTury składa opis zestawu narzędzi tej jednej tury z obu źródeł.
//
// Podstawa bierze się z pola `Agent` okna — tego samego, które nakłada eksperta
// na kanał modelu i którym `adapter_rozmowa_tozsamosc.go` wnosi
// jego warstwy promptu. Drugiego wskazania eksperta nie ma i mieć nie może:
// dwa wskazania to dwie prawdy o tym, kim tura jest.
//
// Dołożenia bierze się z sesji okna, bo dołożenie żyje sesją, nie oknem —
// Operator dokłada narzędzie do pracy, którą prowadzi, a nie do jednego okna.
func (a *adapterRozmowy) zestawTury(ctx context.Context, okno session.Okno) injection.ZestawTury {
	if a == nil {
		return injection.ZestawTury{}
	}
	return injection.ZlozZestawTury(okno.Agent, a.dolozeniaTury(ctx, okno))
}

// dolozeniaTury pyta port o doraźne dołożenia sesji.
//
// Meldunek o niewpiętym porcie wychodzi WYŁĄCZNIE przy oknie z ekspertem, bo
// tylko wtedy brak dołożeń ma skutek dający się nazwać: okno bez eksperta ma
// pełny wykaz kontraktu, więc nieodczytane dołożenie niczego mu nie odbiera.
// Meldunek bez skutku byłby hałasem, w którym giną meldunki ze skutkiem.
func (a *adapterRozmowy) dolozeniaTury(ctx context.Context, okno session.Okno) []string {
	if okno.IdSesji == "" {
		return nil
	}
	if a.dolozeniaSesji == nil {
		if okno.Agent != "" {
			a.zglosZestaw("brak-portu-dolozen-sesji",
				"doraźne dołożenia narzędzi poza turą: okno %q pracuje na zestawie eksperta %q,"+
					" ale port DolozeniaNarzedziSesji nie jest wpięty — narzędzia dołożone"+
					" w sesji %q do tury nie wejdą; tura jedzie samą podstawą eksperta",
				okno.Id, okno.Agent, okno.IdSesji)
		}
		return nil
	}
	nazwy, err := a.dolozeniaSesji.DolozeniaNarzedzi(ctx, okno.IdSesji)
	if err != nil {
		a.zglosZestaw("blad-dolozen-sesji",
			"doraźne dołożenia narzędzi nieodczytane (okno %q, sesja %q): %v —"+
				" tura idzie samą podstawą eksperta",
			okno.Id, okno.IdSesji, err)
		return nil
	}
	return nazwy
}

// argumentyZestawu przekłada opis zestawu na argumenty wpisu `danaco`.
//
// Osobno od przepisania tekstu, bo to jest cała DECYZJA tego pliku: co wolno
// dołożyć, a czego nie wolno, i dlaczego. Przepisanie tekstu niżej jest już
// tylko mechaniką.
//
// Dołożenia bez podstawy nie jadą i mówią, dlaczego. Serwer narzędzi ma dziś
// jedno miejsce, w którym dołożenie może dojść do wykazu — złożenie wykazu
// eksperta (`narzedzia/ekspert_wykaz.go`). Okno bez eksperta nie przechodzi tą
// drogą wcale, więc argument dołożeń nie miałby tam czego dołożyć. Cichy odrzut
// byłby tu gorszy niż brak: Operator widziałby narzędzie na wykazie sesji
// i nie widziałby go w turze.
func (a *adapterRozmowy) argumentyZestawu(zestaw injection.ZestawTury, okno session.Okno) []string {
	argumenty := narzedzia.ArgumentyEksperta(zestaw.Ekspert)
	if len(zestaw.Dolozenia) == 0 {
		return argumenty
	}
	if zestaw.Ekspert == "" {
		a.zglosZestaw("dolozenia-bez-podstawy",
			"doraźne dołożenia sesji %q nie weszły do tury okna %q: okno nie ma nałożonego"+
				" eksperta, a zawężony wykaz — jedyne miejsce, w którym dołożenie ma co"+
				" dołożyć — powstaje tylko dla okna z ekspertem; tura idzie pełnym wykazem",
			okno.IdSesji, okno.Id)
		return argumenty
	}
	return append(argumenty, zestaw.ArgumentyDolozen()...)
}

// zZestawemTury dopisuje do wpisu `danaco` argumenty niosące zestaw tury.
//
// Zestaw pusty nie dokłada NIC — wpis dotychczasowy zostaje wpisem
// dotychczasowym, co do znaku, a tura jedzie pełnym wykazem w zasięgu roli okna,
// jak przed wprowadzeniem składania.
//
// Zasięgu roli ta funkcja nie tyka. Rola okna (`--zasieg klawiatura`) jedzie
// osobną drogą okna asystenta i dokłada, a nie zawęża; zasięg eksperta zawęża
// i wchodzi tutaj. Gdyby oba spotkały się kiedyś na jednym wpisie, rozstrzyga
// `narzedzia.RozpoznajZasieg`, czytając ostatnie wystąpienie przełącznika —
// i to jest rozstrzygnięcie tamtej strony, nie tej.
func zZestawemTury(tekst string, argumenty []string) string {
	if len(argumenty) == 0 {
		return tekst
	}
	var konfiguracja KonfiguracjaMostu
	if err := json.Unmarshal([]byte(tekst), &konfiguracja); err != nil {
		return tekst
	}
	wpis, jest := konfiguracja.McpServers[narzedzia.KluczWpisu]
	if !jest {
		return tekst
	}
	wpis.Args = append(append([]string{}, wpis.Args...), argumenty...)
	konfiguracja.McpServers[narzedzia.KluczWpisu] = wpis
	tresc, err := json.MarshalIndent(konfiguracja, "", "  ")
	if err != nil {
		return tekst
	}
	return string(tresc)
}

// zKonfiguracjaZestawu jest całą drogą złożoną w jedno wywołanie, żeby miejsce
// wpięcia w `adapter_rozmowa_srodowisko.go` pozostało jedną linią.
func (a *adapterRozmowy) zKonfiguracjaZestawu(ctx context.Context, okno session.Okno, tekst string) string {
	zestaw := a.zestawTury(ctx, okno)
	if !zestaw.Skladany() {
		return tekst
	}
	return zZestawemTury(tekst, a.argumentyZestawu(zestaw, okno))
}

// zglosZestaw wpisuje do dziennika rdzenia powód, dla którego zestaw narzędzi
// tury jest taki, jaki jest.
//
// Raz na powód, nie raz na turę — powód nie zmienia się między turami (port
// jest wpięty albo nie), a tur bywa kilkaset dziennie; powtarzany meldunek
// zasłoniłby resztę dziennika i sam przestałby być czytany. Kluczem jest sam
// powód, bez identyfikatorów okna i sesji: te idą do TREŚCI meldunku, ale nie
// do klucza — inaczej każde nowe okno wywoływałoby ten sam meldunek od nowa.
//
// Dziennikiem jest dziennik składacza mostów — ten sam, do którego idzie
// meldunek o braku binarium serwera narzędzi. Jedna sprawa, jedno miejsce.
// Dziennik niewskazany nie zmienia przebiegu tury; znika wyłącznie meldunek.
func (a *adapterRozmowy) zglosZestaw(powod, wzorzec string, wartosci ...any) {
	if a == nil || a.mosty == nil || a.mosty.dziennik == nil || powod == "" {
		return
	}
	if _, juzBylo := a.zgloszoneZestawy.LoadOrStore(powod, true); juzBylo {
		return
	}
	a.mosty.dziennik.Printf(wzorzec, wartosci...)
}
