// Odpowiedzialność pliku: budowniczowie rejestrów i źródeł, z których korzysta
// montaż rdzenia — rejestr kanałów modelu wraz z pulą kont rotacji, źródło
// wierszy konfiguracji spod adresu złożonego, rejestr definicji z katalogu
// ustawień oraz katalog akcji.
//
// Wszystkie są sterowane danymi: nowy kanał, nowa pozycja
// okna konfiguracji i nowa akcja to nowy wiersz, nie nowa gałąź w kodzie.
// Niepowodzenie pierwszego odczytu nigdy nie przerywa startu — rdzeń rusza
// z zawartością uboższą i odbudowuje ją przy kolejnym odczycie.
//
// Osobno od montaz.go, bo montaż mówi, co z czym się wiąże, a ten plik — jak
// powstaje każde źródło.
package core

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// rejestrKanalow buduje rejestr kanałów z wierszy tabeli rejestru
// i wnosi do niego kanał główny. Niepowodzenie pierwszego odczytu nie
// przerywa startu — rejestr odbuduje się przy pierwszej zmianie wiersza.
func rejestrKanalow(kontekst context.Context, m Montaz, repozytoria *dane.Zestaw,
	przejmowanie *przejmowanieProcesow, zdarzenia *zdarzeniaWykonawcze) *models.Rejestr {
	rejestr := models.NowyRejestr(models.NoweZrodloBazy(m.Baza.DB), models.FabrykiWbudowane())
	// Kanał główny obsługuje oba rodzaje procesu lokalnego: „cli" (program code
	// CLI) i „lokalny" (proces lokalny rozmawiający strumieniem). Obie fabryki
	// dzielą jedną pulę kont — rotacja po wyczerpaniu limitu ma sens tylko przy
	// wspólnej pamięci wyczerpania. Rodzaj „sdk" nie ma tu fabryki:
	// bez dostawcy SDK kanał tego rodzaju nie ma jak działać, więc rejestr nie
	// pokaże go jako czynnego, zamiast udawać gotowość.
	fabrykaProcesu := fabrykaKanaluGlownego(pulaKont(kontekst, m, repozytoria), przejmowanie,
		zdarzenia, nazwaKontaZKatalogu(kontekst, repozytoria))
	rejestr.UstawFabryke(models.AdapterCLI, fabrykaProcesu)
	rejestr.UstawFabryke(rodzajKanaluLokalny, fabrykaProcesu)
	if err := rejestr.Odswiez(kontekst); err != nil && m.Dziennik != nil {
		m.Dziennik.Printf("rejestr kanałów: %v", err)
	}
	return rejestr
}

// rodzajKanaluLokalny jest wartością kolumny kanal_modelu.rodzaj_kanalu dla
// kanału procesu lokalnego (KnownChannelKinds). Jest wartością
// danych, nie nazwą typu — kanał tego rodzaju to wiersz, nie gałąź w kodzie.
const rodzajKanaluLokalny = "lokalny"

// nazwaKontaZKatalogu tłumaczy wskazanie konta na kod (nazwę) konta puli.
//
// Kontrakt identyfikuje konto liczbą (`Account.id` = klucz główny wiersza),
// a pula rotacji kodem (`konto.nazwa`) — wskazanie z obszaru account
// konfiguracji sesji i z kolumny `kanal_modelu.konto_id` przychodzi więc
// w innym słowniku niż ten, którym mówi pula. Wskazanie nieliczbowe
// jedzie dosłownie: Operator mógł wpisać nazwę wprost. Katalog pusty zostawia
// wskazanie bez tłumaczenia — tura powie wtedy, że konta nie zna.
func nazwaKontaZKatalogu(kontekst context.Context, repozytoria *dane.Zestaw) func(string) string {
	if repozytoria == nil || repozytoria.Konta == nil {
		return nil
	}
	repo := repozytoria.Konta
	return func(wskazanie string) string {
		id, err := strconv.ParseInt(wskazanie, 10, 64)
		if err != nil {
			return wskazanie
		}
		konto, err := repo.Pobierz(kontekst, id)
		if err != nil {
			return wskazanie
		}
		return konto.Nazwa
	}
}

// pulaKont bierze konta rotacji z katalogu kont (tabela `konto`, migracja 014).
// Katalog jest źródłem pierwszym; katalog profili na dysku zostaje ścieżką
// zapasową, żeby instalacja bez wpisanych kont dalej pracowała.
func pulaKont(kontekst context.Context, m Montaz, repozytoria *dane.Zestaw) *injection.PulaKont {
	if repozytoria != nil && repozytoria.Konta != nil {
		konta, err := repozytoria.Konta.KontaRotacji(kontekst, shared.AccountKindCli)
		if err != nil && m.Dziennik != nil {
			m.Dziennik.Printf("katalog kont: %v", err)
		}
		// Katalog kont jest źródłem pierwszym, gdy repozytorium istnieje —
		// niezależnie od tego, czy na starcie ma już wpisy. Pulę budujemy
		// z bieżącego wykazu (może być pusty) i ZAWSZE wpinamy w nią źródło
		// oraz utrwalacz: wyczerpanie zapisuje się do bazy, a bieżący wykaz
		// odczytuje się na progu tury przez OdświeżZeŹródła, więc konto dodane
		// komendą account.* wchodzi do rotacji bez restartu — także gdy pula
		// wystartowała pusta. Ścieżka zapasowa z dysku zostaje
		// tylko dla instalacji bez katalogu kont (repozytorium kont brak).
		pula := injection.NowaPula(naPuleKont(konta)...)
		wyposazPuleKont(kontekst, pula, m, repozytoria.Konta)
		return pula
	}
	if m.KatalogProfili == "" {
		return injection.NowaPula()
	}
	konta, err := injection.KontaZKatalogu(m.KatalogProfili)
	if err != nil && m.Dziennik != nil {
		m.Dziennik.Printf("profile kanału głównego: %v", err)
	}
	return injection.NowaPula(konta...)
}

// wyposazPuleKont wpina w pulę trwałość wyczerpania i odświeżanie z katalogu.
// Trwałość zapisuje stan konta (konto.stan / konto.wyczerpane_do), żeby limit
// przeżył restart; odświeżanie odczytuje bieżący wykaz kont rotacji, żeby zmiana
// katalogu komendą account.* dotarła do rotacji bez restartu.
func wyposazPuleKont(kontekst context.Context, pula *injection.PulaKont, m Montaz,
	repo dane.RepozytoriumKont) {

	pula.UstawUtrwalanie(func(kod string, doChwili time.Time) {
		konto, err := repo.PobierzPoNazwie(kontekst, kod)
		if err != nil {
			if m.Dziennik != nil {
				m.Dziennik.Printf("utrwalenie wyczerpania konta %q: %v", kod, err)
			}
			return
		}
		iso := doChwili.UTC().Format(time.RFC3339)
		if err := repo.OznaczStan(kontekst, konto.ID, dane.StanKontaWyczerpane, &iso); err != nil && m.Dziennik != nil {
			m.Dziennik.Printf("utrwalenie wyczerpania konta %q: %v", kod, err)
		}
	})

	pula.UstawZrodlo(func() ([]injection.Konto, bool) {
		konta, err := repo.KontaRotacji(kontekst, shared.AccountKindCli)
		if err != nil {
			if m.Dziennik != nil {
				m.Dziennik.Printf("odświeżenie puli kont: %v", err)
			}
			return nil, false
		}
		return naPuleKont(konta), true
	})
}

// naPuleKont przekłada wiersze katalogu na konta puli rotacji. Konto bez
// katalogu konfiguracji do puli nie wchodzi — kanał główny bierze profil
// z katalogu na dysku, nie z bazy. Stan wyczerpania odczytany
// z katalogu (konto.stan / konto.wyczerpane_do) jedzie w polu WyczerpaneDo,
// żeby pula odtworzyła limit po restarcie.
func naPuleKont(konta []dane.Konto) []injection.Konto {
	pula := make([]injection.Konto, 0, len(konta))
	for _, konto := range konta {
		if konto.KatalogKonfiguracji == nil || *konto.KatalogKonfiguracji == "" {
			continue
		}
		wpis := injection.Konto{
			Kod:                 konto.Nazwa,
			KatalogKonfiguracji: *konto.KatalogKonfiguracji,
		}
		if konto.Stan == dane.StanKontaWyczerpane && konto.WyczerpaneDo != nil {
			if chwila, err := time.Parse(time.RFC3339, strings.TrimSpace(*konto.WyczerpaneDo)); err == nil {
				wpis.WyczerpaneDo = chwila
			}
		}
		pula = append(pula, wpis)
	}
	return pula
}

// zrodloUstawienOsiZBazy podaje rozstrzygaczowi wiersze spod adresu złożonego:
// poziom zasięgu razem z osią rozstrzygania (migracja 012). Błąd
// odczytu nie zatrzymuje rozstrzygania — ustawienia bez zapisu schodzą na
// wartości domyślne.
func zrodloUstawienOsiZBazy(kontekst context.Context,
	repozytorium dane.RepozytoriumKonfiguracjiOsi) konfig.Zrodlo {

	return konfig.NoweZrodloZOdczytuOsi(func(adres konfig.Adres) ([]konfig.Wpis, error) {
		ustawienia, err := repozytorium.ListaOsi(kontekst, adres.Poziom, adres.KluczZasiegu,
			adres.Os, adres.KluczOsi)
		if err != nil {
			return nil, err
		}
		wpisy := make([]konfig.Wpis, 0, len(ustawienia))
		for _, u := range ustawienia {
			wpisy = append(wpisy, konfig.Wpis{
				Poziom: u.Poziom, KluczZasiegu: u.KluczZasiegu,
				Os: u.Os, KluczOsi: u.KluczOsi, Klucz: u.Klucz,
				Wartosc: wartoscTekstu(u.Wartosc), Rodzaj: konfig.Rodzaj(u.RodzajWartosci),
			})
		}
		return wpisy, nil
	})
}

// katalogUstawienZBazy podaje rejestrowi definicji pozycje katalogu ustawień.
type katalogUstawienZBazy struct {
	kontekst     context.Context
	repozytorium dane.RepozytoriumKatalogUstawien
}

// Definicje zwraca pozycje aktywne katalogu — rejestr definicji nie zna wierszy
// wygaszonych, bo nie ma czego dla nich rozstrzygać.
func (k katalogUstawienZBazy) Definicje() ([]shared.SettingDefinition, error) {
	if k.repozytorium == nil {
		return nil, nil
	}
	return k.repozytorium.Definicje(k.kontekst, true)
}

// rejestrUstawien buduje rejestr definicji z katalogu ustawień:
// nowa pozycja okna konfiguracji to nowy wiersz migracji, nie nowa stała w kodzie.
//
// Katalog pusty albo niedostępny daje rejestr wbudowany rdzenia; wtedy — i tylko
// wtedy — dokładane są definicje katalogu roboczego, bo ich wiersze katalogu nie
// dojechały.
func rejestrUstawien(kontekst context.Context, repozytoria *dane.Zestaw,
	dziennik *log.Logger) *konfig.Rejestr {

	rejestr, zKatalogu := konfig.RejestrZKatalogu(katalogUstawienZBazy{kontekst, repozytoria.KatalogUstawien})
	if zKatalogu {
		return rejestr
	}
	if dziennik != nil {
		dziennik.Printf("katalog ustawień: pusty albo niedostępny, obowiązuje rejestr wbudowany")
	}
	rejestr.Dodaj(DefinicjeKataloguRoboczego()...)
	return rejestr
}

// rejestrAkcji buduje katalog akcji z wierszy tabeli `akcja`.
// Niepowodzenie pierwszego odczytu nie przerywa startu — rejestr odbuduje się
// przy pierwszym odczycie katalogu.
func rejestrAkcji(kontekst context.Context, repozytoria *dane.Zestaw, dziennik *log.Logger) *RejestrAkcji {
	rejestr := NowyRejestrAkcji(ZrodloAkcjiZRepozytorium(repozytoria.Akcje))
	if err := rejestr.Odswiez(kontekst); err != nil && dziennik != nil {
		dziennik.Printf("katalog akcji: %v", err)
	}
	return rejestr
}
