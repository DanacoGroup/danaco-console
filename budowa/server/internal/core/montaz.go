package core

import (
	"context"
	"log"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/store"
	"danacoconsole/server/internal/transport"
)

// Montaz jest kompletem tego, co punkt wejścia wie o świecie: ustawieniami
// procesu, otwartą bazą i dziennikiem. Wszystko pozostałe powstaje tutaj.
type Montaz struct {
	Konfiguracja konfiguracja.Konfiguracja
	Baza         *store.Baza
	Dziennik     *log.Logger
	// KatalogKlienta wskazuje pakiet interfejsu serwowany obok gniazda.
	KatalogKlienta string
	// KatalogProfili wskazuje katalog profili kanału głównego. Poświadczenia
	// zostają w profilach na dysku i nigdy nie wchodzą do repozytorium ani do
	// bazy — rdzeń zna wyłącznie odwołanie.
	KatalogProfili string
}

// Zmontowany niesie rdzeń wraz z zasobami, które trzeba zwolnić po zatrzymaniu.
type Zmontowany struct {
	Rdzen *Rdzen

	// Pola katalogu roboczego tu nie ma z zamysłem. Katalog sesji dojeżdża do
	// procesu modelu drogą tury — adapter rozmowy ustala go i wkłada do
	// zapytania kanału — więc wystawianie go drugi raz na Zmontowanym byłoby
	// eksportem bez odbiorcy.

	dane     *dane.Zestaw
	kanaly   *models.Rejestr
	nadzorca *session.Nadzorca
	terminal *adapterTerminala
	// developer trzeba zwolnić osobno: przebieg budowania biegnie poza żądaniem
	// i bez tego przeżyłby zatrzymanie rdzenia jako sierota.
	developer *adapterDevelopera
	// ustawienia to ten sam adapter nastaw, który idzie do portów. Trzymany tu,
	// bo drogi wewnętrzne rdzenia (np. zgłoszenie do centrum powiadomień) czytają
	// nim nastawy Operatora poza obsługą komendy.
	ustawienia *adapterUstawienOsi
	// aplikacje trzeba zwolnić osobno: `apps.preview.start` podnosi nasłuch HTTP
	// serwera podglądu, który bez zamknięcia przeżyłby zatrzymanie rdzenia
	// i zostawił zajęty port.
	aplikacje *adapterAplikacji
}

// Zmontuj składa cały rdzeń: repozytoria nad bazą, rozstrzygacz ośmiu poziomów
// zasięgu, rejestr kanałów modelu z wierszy bazy, nadzorcę sesji i okien,
// serwer transportu oraz rejestr obsługiwaczy komend.
//
// Kolejność jest wymuszona zależnościami, nie upodobaniem: repozytoria dają
// źródła konfiguracji i kanałów, transport daje nadajnik zdarzeń, a rdzeń
// powstaje na końcu, bo dopiero wtedy ma czym wypełnić porty.
//
// Brak elementu opcjonalnego nie przerywa montażu: pusty rejestr
// kanałów, brak profili kanału głównego i brak pakietu klienta zostawiają rdzeń
// zdolny do pracy w pozostałym zakresie.
func Zmontuj(kontekst context.Context, m Montaz) (*Zmontowany, error) {
	repozytoria, err := dane.Otworz(kontekst, m.Baza)
	if err != nil {
		return nil, err
	}

	// Rozpoznanie maszyny bieżącej: zakłada wiersz urządzenia, na którym stoi
	// rdzeń, tak by punkt dostępu localDirectory miał na co wskazać. Nieudane
	// rozpoznanie idzie do dziennika i nie przerywa montażu.
	odnotujUrzadzenieBiezace(kontekst, repozytoria.Urzadzenia, m.Dziennik)

	rozstrzygacz := konfig.Nowy(zrodloUstawienOsiZBazy(kontekst, repozytoria.KonfiguracjaOsi),
		rejestrUstawien(kontekst, repozytoria, m.Dziennik))

	// Katalog roboczy sesji. Degradacja do lokalizacji
	// zastępczej idzie do dziennika, bo Operator ma wiedzieć, że pracuje gdzie
	// indziej, niż ustawił.
	katalogRoboczy := NowyKatalogRoboczy(rozstrzygacz,
		ZObserwatoremKatalogu(ObserwatorKataloguFunkcja(func(d DegradacjaKatalogu) {
			if m.Dziennik == nil {
				return
			}
			m.Dziennik.Printf("katalog roboczy %q niezdatny do zapisu (%s) — praca w %q (zdatny=%t)",
				d.Zadana, d.Powod, d.Zastepcza, d.Skuteczna)
		})))

	// Przejmowanie procesów powstaje przed kanałami, bo kanał wkłada
	// haczyk do każdej tury, a wiązane jest po nadzorcy, bo to on ma rejestr
	// procesów. Ta jedna pośredniczka domyka różnicę kolejności.
	//
	// Nadzorca nie dostaje tu żadnego wypełnienia: proces tury startuje kanał
	// modelu własną drogą, a sesja obejmuje go uchwytem przez Przejmij. Drugiej
	// drogi startu procesu nie ma.
	przejmowanie := &przejmowanieProcesow{}
	// Odbiornik zdarzeń wykonawczych powstaje przed rejestrem kanałów,
	// bo fabryka kanału głównego dostaje go do ręki; diagnostyka dopina się
	// niżej, po złożeniu modułów — taka jest kolejność montażu.
	zdarzeniaWykonawcze := nowyOdbiorZdarzenWykonawczych(kontekst, repozytoria, m.Dziennik)
	kanaly := rejestrKanalow(kontekst, m, repozytoria, przejmowanie, zdarzeniaWykonawcze)
	nadzorca := session.NowyNadzorca()
	przejmowanie.Zwiaz(nadzorca.Procesy().Przejmij)
	// Sprzątanie startowe stanu trwałego idzie przed odtworzeniem rejestru: start
	// jest jedynym momentem, w którym rdzeń i tak czyta stan trwały. Opróżnienie
	// kosza sesji po terminie i przemiecenie retencji historii to dwie czynności
	// jednej drogi sprzątania; kolejność jest istotna, bo sesja skasowana z kosza
	// zabiera swoje okna wraz z wypowiedziami.
	usunSesjePoTerminie(kontekst, repozytoria, m.Dziennik)
	przemiecRetencjeHistorii(kontekst, repozytoria, m.Dziennik)
	odtworzStanZBazy(kontekst, repozytoria, nadzorca, m.Dziennik)
	serwer := transport.Nowy(ustawieniaTransportu(m, rozstrzygacz))
	nasluch := nasluchTransportu{serwer: serwer}
	utrwalacz := utrwalaczRozmow(repozytoria, nadzorca)
	dziennikRozmow := nowyDziennikRozmowy(kontekst, utrwalacz, m.Dziennik)

	// Telemetria postępu czyta szynę zdarzeń i port rozmowy, więc
	// powstaje przed rdzeniem i owija nadajnik transportu.
	//
	// Żywy stan sesji stoi pod telemetrią: producent telemetrii rozgłasza
	// `progress.changed` wprost tym nadajnikiem, który dostał, więc nasłuch
	// obecności musi być tym nadajnikiem. Tak domyka się łańcuch producent →
	// szyna → odbiorca kontrolki powrotu do sesji.
	biegi := nowyRejestrBiegow()
	obecnosc := nowyRejestrObecnosci(kontekst, nadzorca, biegi).
		ZeZrodlami(zrodlaObecnosciZBazy(repozytoria))
	telemetria := nowaTelemetriePostepu(obecnosc.OwinNadajnik(nasluch))
	szynaZdarzen := telemetria.OwinNadajnik(nasluch)

	ustawienia := nowyAdapterUstawienOsi(repozytoria.KonfiguracjaOsi, rozstrzygacz)
	tozsamosc := nowyAdapterTozsamosci(repozytoria.Tozsamosc, rozstrzygacz).
		ZOknami(repozytoria.Okna, repozytoria.Kanaly)
	// Doraźne dołożenia narzędzi mają dwóch czytelników: port `NarzedziaSesji`
	// rdzenia (rodzina `session.tool.*` i `tools.catalog.list`) oraz składacz
	// zestawu narzędzi tury, który dokłada je do wykazu eksperta
	// (`adapter_rozmowa_zestaw.go`). Instancja powstaje tu, bo składanie rozmowy
	// biegnie przed składaniem portów, a oba mają dostać ten sam adapter —
	// drugi byłby drugą prawdą o tym, czym model w sesji dysponuje.
	dolozeniaNarzedzi := nowyAdapterNarzedziSesji(repozytoria.NarzedziaSesji(),
		repozytoria.Sesje, repozytoria.Rozszerzenia())
	rozmowa, petla := zlozRozmowe(skladRozmowy{
		zycie: kontekst, montaz: m, repozytoria: repozytoria, nadzorca: nadzorca,
		kanaly: kanaly, nadajnik: szynaZdarzen, dziennik: dziennikRozmow,
		rozstrzygacz: rozstrzygacz, katalog: katalogRoboczy, tozsamosc: tozsamosc,
		biegi: biegi, obecnosc: obecnosc, ustawienia: ustawienia,
		zdarzenia: zdarzeniaWykonawcze, dolozenia: dolozeniaNarzedzi,
	})
	klienci := noweWieziKlientow()
	trwalosc := nowyUtrwalaczStanow(kontekst, repozytoria, m.Dziennik)

	moduly := zlozAdapteryModulow(kontekst, m, repozytoria, nadzorca,
		rozstrzygacz, katalogRoboczy, telemetria, szynaZdarzen, kanaly)
	terminal, developer := moduly.terminal, moduly.developer
	kolejki, diagnostyka := moduly.kolejki, moduly.diagnostyka
	// Apps składa się tu, a nie w `zlozPorty`, bo montaż musi go potem zamknąć:
	// serwer podglądu żyje poza żądaniem, tak samo jak przebieg budowania
	// Developera.
	aplikacje := nowyAdapterAplikacji(repozytoria.Aplikacje).ZOknami(nadzorca.Rejestr()).
		ZMagazynem(m.Konfiguracja.KatalogDanych).ZKatalogiemRozszerzen(repozytoria.Rozszerzenia())

	// Zakresy narzędzi mają dwóch czytelników: port `ZakresyNarzedzi` rdzenia
	// (rodzina `tools.scope.*`) oraz straż, którą rdzeń pyta przed skierowaniem
	// komendy ręki modelu. Jedna instancja na obie drogi — druga byłaby drugą
	// prawdą o tym, co profilowi wolno.
	zakresyNarzedzi := NowyPortZakresowNarzedzi(repozytoria.ZakresyNarzedzi, repozytoria.Asystent)
	// Straż zakresu eksperta czyta bibliotekę i katalog modułów. Wpina się
	// w dwa miejsca odmowy: nałożenie eksperta na okno i powołanie podagentów
	// (`straz_eksperta.go`).
	strazEkspertow := NowaStrazEksperta(repozytoria.Agenci, repozytoria.Moduly)

	rdzen := Zloz(zlozPorty(skladPortow{
		zycie: kontekst, kolejki: kolejki, diagnostyka: diagnostyka,
		montaz: m, repozytoria: repozytoria, nadzorca: nadzorca, kanaly: kanaly,
		rozstrzygacz: rozstrzygacz, katalogRoboczy: katalogRoboczy,
		trwalosc: trwalosc, obecnosc: obecnosc, biegi: biegi, telemetria: telemetria,
		rozmowa: rozmowa, petla: petla, ustawienia: ustawienia, tozsamosc: tozsamosc,
		klienci: klienci, utrwalacz: utrwalacz, dziennikRozmow: dziennikRozmow,
		terminal: terminal, developer: developer, aplikacje: aplikacje,
		automatyki:   moduly.automatyki,
		mowa:         moduly.mowa,
		doradcy:      moduly.doradcy,
		szynaZdarzen: szynaZdarzen, nasluch: nasluch,
		dolozenia:       dolozeniaNarzedzi,
		zakresyNarzedzi: zakresyNarzedzi,
		strazEkspertow:  strazEkspertow,
		nadajnik:        nastawyNadajnika(m.Konfiguracja),
	}))
	rdzen.ZeStrazaZakresow(zakresyNarzedzi)
	rdzen.ZObserwatoremNiepowodzen(diagnostyka)
	// Diagnostyka jest odbiorcą odmów — zdarzenia zaczepów mają być w Errors
	// Panel faktami, nie ciszą.
	zdarzeniaWykonawcze.PodepnijDiagnostyke(diagnostyka)
	serwer.PodlaczRdzen(wejscieTransportu{rdzen: rdzen})

	// Budzik harmonogramu startuje razem z rdzeniem: cyklicznie odpala
	// automatyki, których termin minął, tą samą drogą co Operator z Queue
	// Managera. Wątek żyje aż do zamknięcia kontekstu życia.
	nowyBudzikHarmonogramu(moduly.automatyki, m.Dziennik).Uruchom(kontekst)

	return &Zmontowany{Rdzen: rdzen, dane: repozytoria, kanaly: kanaly,
		nadzorca: nadzorca, terminal: terminal, developer: developer,
		aplikacje: aplikacje, ustawienia: ustawienia}, nil
}

// Zamknij zwalnia zasoby złożonego rdzenia: procesy okien, kanały modelu
// i przygotowane zapytania. Niepowodzenie jednego zwolnienia nie wstrzymuje
// pozostałych; bazę zamyka ten, kto ją otworzył.
func (z *Zmontowany) Zamknij() {
	if z == nil {
		return
	}
	z.terminal.Zamknij()
	z.developer.Zamknij()
	z.aplikacje.Zamknij()
	_ = z.nadzorca.Zamknij()
	_ = z.kanaly.Zamknij()
	_ = z.dane.Zamknij()
}

// ustawieniaTransportu przenosi nastawy rdzenia do warstwy nasłuchu. Tędy idzie
// cały brzeg, nie sam port: montaż jest jedynym miejscem, przez które
// konfiguracja startu — port, TLS, wykaz pochodzeń, adres nasłuchu — dochodzi
// do transportu.
//
// Wymóg logowania ma dwa źródła i jedno pierwszeństwo. Wskazanie ze startu
// (przełącznik wiersza poleceń, zmienna środowiska) wygrywa zawsze. Dopiero
// jego brak oddaje głos nastawie poziomu `aplikacja` z tabeli `ustawienie`,
// a brak i jej — adresowi nasłuchu. To ten sam rozstrzygacz i ta sama tabela,
// którą widzi `config.get`.
func ustawieniaTransportu(m Montaz, rozstrzygacz *konfig.Rozstrzygacz) transport.Ustawienia {
	u := transport.Domyslne()
	u.Port = m.Konfiguracja.Port
	u.KatalogKlienta = m.KatalogKlienta
	u.Dziennik = m.Dziennik
	u.Adres = m.Konfiguracja.Adres
	u.WszystkieInterfejsy = m.Konfiguracja.WszystkieInterfejsy
	u.CertyfikatTLS = m.Konfiguracja.CertyfikatTLS
	u.KluczTLS = m.Konfiguracja.KluczTLS
	u.PochodzeniaDozwolone = m.Konfiguracja.PochodzeniaDozwolone
	u.WymogLogowania = m.Konfiguracja.WymogLogowania
	if u.WymogLogowania == nil {
		if wymog, wskazana := wartoscWymoguLogowania(
			rozstrzygacz.Rozstrzygnij(konfig.Kontekst{}, konfig.KluczWymogLogowania).Wartosc); wskazana {
			u.WymogLogowania = &wymog
		}
	}
	return u
}

// sladObiegu buduje obserwatora pętli koordynator–wykonawca zapisującego bieg
// w dzienniku rdzenia. Zatrzymanie biegu nigdy nie jest ciche — powód
// i licznik obiegów wychodzą wprost do dziennika. Brak dziennika nie wyłącza
// pętli, wyłącza wyłącznie ślad.
func sladObiegu(dziennik *log.Logger) func(session.ZdarzenieObiegu) {
	return func(z session.ZdarzenieObiegu) {
		if dziennik == nil {
			return
		}
		dziennik.Printf("pętla: koordynator=%s obieg=%d bez-postępu=%d rozpoczęty=%v zatrzymany=%v powód=%q błąd=%v",
			z.Stan.IdKoordynatora, z.Stan.Obiegow, z.Stan.ObiegowBezPostepu,
			z.Rozpoczety, z.Stan.Zatrzymany, z.Stan.Powod, z.Blad)
	}
}
