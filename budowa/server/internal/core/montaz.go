package core

import (
	"context"
	"log"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/injection"
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
	// KatalogProfili wskazuje katalog profili kanału głównego przechowywany poza bazą i repozytorium.
	KatalogProfili string
}

// Zmontowany niesie rdzeń złożonego procesu wraz z zasobami, które wymagają osobnego zwolnienia po zatrzymaniu, bo przeżyłyby samo zamknięcie rdzenia jako sieroty.
type Zmontowany struct {
	Rdzen *Rdzen

	// Katalog roboczy sesji nie jest polem: ustala go adapter rozmowy i wkłada do zapytania kanału.

	dane     *dane.Zestaw
	kanaly   *models.Rejestr
	nadzorca *session.Nadzorca
	terminal *adapterTerminala
	// developer wymaga osobnego zwolnienia, bo przebieg budowania przeżyłby zatrzymanie rdzenia.
	developer *adapterDevelopera
	// ustawienia to adapter nastaw współdzielony z portami; czytają go też drogi wewnętrzne rdzenia.
	ustawienia *adapterUstawienOsi
	// aplikacje wymaga osobnego zwolnienia: nasłuch podglądu żyje poza żądaniem i poza rdzeniem.
	aplikacje *adapterAplikacji
}

// Zmontuj składa rdzeń z repozytoriów nad bazą, rozstrzygacza zasięgu, rejestru kanałów modelu, nadzorcy sesji, transportu i rejestru obsługiwaczy komend, w kolejności wymuszonej zależnościami między nimi.
func Zmontuj(kontekst context.Context, m Montaz) (*Zmontowany, error) {
	repozytoria, err := dane.Otworz(kontekst, m.Baza)
	if err != nil {
		return nil, err
	}

	// Rozpoznanie maszyny bieżącej zakłada wiersz urządzenia dla punktu dostępu localDirectory.
	odnotujUrzadzenieBiezace(kontekst, repozytoria.Urzadzenia, m.Dziennik)

	rozstrzygacz := konfig.Nowy(zrodloUstawienOsiZBazy(kontekst, repozytoria.KonfiguracjaOsi),
		rejestrUstawien(kontekst, repozytoria, m.Dziennik))

	// Katalog roboczy sesji; degradacja do lokalizacji zastępczej trafia do dziennika.
	katalogRoboczy := NowyKatalogRoboczy(rozstrzygacz,
		ZObserwatoremKatalogu(ObserwatorKataloguFunkcja(func(d DegradacjaKatalogu) {
			if m.Dziennik == nil {
				return
			}
			m.Dziennik.Printf("katalog roboczy %q niezdatny do zapisu (%s) — praca w %q (zdatny=%t)",
				d.Zadana, d.Powod, d.Zastepcza, d.Skuteczna)
		})))

	// Przejmowanie procesów powstaje przed kanałami i wiąże się z nadzorcą po jego złożeniu.
	przejmowanie := &przejmowanieProcesow{}
	// Odbiornik zdarzeń wykonawczych powstaje przed rejestrem kanałów, bo fabryka go potrzebuje.
	zdarzeniaWykonawcze := nowyOdbiorZdarzenWykonawczych(kontekst, repozytoria, m.Dziennik)
	kanaly := rejestrKanalow(kontekst, m, repozytoria, przejmowanie, zdarzeniaWykonawcze)
	nadzorca := session.NowyNadzorca()
	przejmowanie.Zwiaz(nadzorca.Procesy().Przejmij)
	// Sprzątanie stanu trwałego biegnie przed odtworzeniem rejestru, w ustalonej kolejności.
	usunSesjePoTerminie(kontekst, repozytoria, m.Dziennik)
	przemiecRetencjeHistorii(kontekst, repozytoria, m.Dziennik)
	odtworzStanZBazy(kontekst, repozytoria, nadzorca, m.Dziennik)
	serwer := transport.Nowy(ustawieniaTransportu(m, rozstrzygacz))
	nasluch := nowyNasluchTransportu(serwer)
	utrwalacz := utrwalaczRozmow(repozytoria, nadzorca)
	dziennikRozmow := nowyDziennikRozmowy(kontekst, utrwalacz, m.Dziennik)

	// Telemetria postępu czyta szynę zdarzeń i port rozmowy, więc powstaje przed rdzeniem.
	biegi := nowyRejestrBiegow()
	obecnosc := nowyRejestrObecnosci(kontekst, nadzorca, biegi).
		ZeZrodlami(zrodlaObecnosciZBazy(repozytoria))
	telemetria := nowaTelemetriePostepu(obecnosc.OwinNadajnik(nasluch))
	szynaZdarzen := telemetria.OwinNadajnik(nasluch)

	ustawienia := nowyAdapterUstawienOsi(repozytoria.KonfiguracjaOsi, rozstrzygacz)
	tozsamosc := nowyAdapterTozsamosci(repozytoria.Tozsamosc, rozstrzygacz).
		ZOknami(repozytoria.Okna, repozytoria.Kanaly)
	// Doraźne dołożenia narzędzi czyta port NarzedziaSesji rdzenia i składacz zestawu tury.
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
	// Apps składa się tu, nie w zlozPorty, bo montaż musi zamknąć jego nasłuch po zatrzymaniu.
	aplikacje := nowyAdapterAplikacji(repozytoria.Aplikacje).ZOknami(nadzorca.Rejestr()).
		ZMagazynem(m.Konfiguracja.KatalogDanych).ZKatalogiemRozszerzen(repozytoria.Rozszerzenia()).
		ZUruchamiaczem(injection.UruchamiaczOkien(), rozstrzygacz, katalogRoboczy)

	// Zakresy narzędzi ma dwóch czytelników: port ZakresyNarzedzi rdzenia i straż komend modelu.
	zakresyNarzedzi := NowyPortZakresowNarzedzi(repozytoria.ZakresyNarzedzi, repozytoria.Asystent)
	// Straż zakresu eksperta czyta bibliotekę i katalog modułów przy nakładaniu eksperta na okno.
	strazEkspertow := NowaStrazEksperta(repozytoria.Agenci, repozytoria.Moduly)

	porty := zlozPorty(skladPortow{
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
	})
	rdzen := Zloz(porty)
	/* Rozpoznanie konta wołającego bierze się z tego samego adaptera bramki,
	   który wydaje sesje — inaczej rdzeń nie miałby jak powiedzieć, czyje jest
	   żądanie, a karty sesji wróciłyby wspólne dla wszystkich kont. */
	if rozpoznanie, umie := porty.Uwierzytelnianie.(RozpoznanieKontaSesji); umie {
		rdzen.ZRozpoznaniemKontaSesji(rozpoznanie)
	}
	// Zrywanie gniazd po unieważnieniu sesji czyta więź rdzenia i wiersz sesji
	// z tego samego adaptera bramki; oba powstają dopiero tutaj.
	waznosc, _ := porty.Uwierzytelnianie.(RozpoznanieWaznosciSesji)
	nasluch.sesje.uzupelnij(rdzen.wiez, waznosc)
	if urzadzenia, umie := porty.Urzadzenia.(zRozlaczaniem); umie {
		urzadzenia.przyjmijRozlaczanie(nasluch)
	}
	// Tor strumieni wiąże turę z gniazdem zamawiającym; wpis robi adapter rozmowy przy otwarciu tury.
	rozmowa.ZTorem(nasluch.tor)
	rdzen.ZeStrazaZakresow(zakresyNarzedzi)
	rdzen.ZObserwatoremNiepowodzen(diagnostyka)
	// Diagnostyka jest odbiorcą odmów — zdarzenia zaczepów mają być w Errors
	// Panel faktami, nie ciszą.
	zdarzeniaWykonawcze.PodepnijDiagnostyke(diagnostyka)
	serwer.PodlaczRdzen(wejscieTransportu{rdzen: rdzen, tor: nasluch.tor})

	// Budzik harmonogramu odpala automatyki po terminie i żyje aż do zamknięcia kontekstu życia.
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

// ustawieniaTransportu przenosi nastawy rdzenia do warstwy nasłuchu transportu: port, TLS, wykaz pochodzeń i adres, a wymóg logowania rozstrzyga według pierwszeństwa startu, nastawy osi i adresu nasłuchu.
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
