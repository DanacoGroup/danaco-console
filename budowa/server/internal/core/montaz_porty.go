package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/nadajnik"
	"danacoconsole/server/internal/session"
)

// skladPortow jest kompletem bytow, z ktorych powstaja porty rdzenia.
//
// Struktura zamiast dwudziestu argumentow. Przy wywolaniu pozycyjnym o tylu
// polach kazda zmiana kolejnosci bylaby cicha pomylka nie do wylapania przez
// kontrole typow — polowa tych pol ma ten sam typ.
type skladPortow struct {
	zycie          context.Context
	montaz         Montaz
	repozytoria    *dane.Zestaw
	nadzorca       *session.Nadzorca
	kanaly         *models.Rejestr
	rozstrzygacz   *konfig.Rozstrzygacz
	katalogRoboczy *KatalogRoboczy
	trwalosc       *utrwalaczStanow
	obecnosc       *rejestrObecnosci
	biegi          *rejestrBiegow
	telemetria     *telemetriaPostepu
	rozmowa        *adapterRozmowy
	petla          *session.Petla
	ustawienia     *adapterUstawienOsi
	tozsamosc      Tozsamosc
	klienci        *wieziKlientow
	utrwalacz      *dane.UtrwalaczRozmowy
	dziennikRozmow *dziennikRozmowy
	terminal       *adapterTerminala
	developer      *adapterDevelopera
	// aplikacje trzeba zwolnić osobno: `apps.preview.start` podnosi nasłuch HTTP,
	// który bez zamknięcia przeżyłby zatrzymanie rdzenia i zostawił zajęty port.
	aplikacje    *adapterAplikacji
	diagnostyka  *adapterDiagnostyki
	kolejki      *adapterKolejek
	automatyki   *adapterAutomatyk
	mowa         *adapterMowy
	doradcy      *adapterDoradcow
	szynaZdarzen Nadajnik
	nasluch      nasluchTransportu
	// dolozenia jest tym samym adapterem, ktory wnosi dolozenia sesji do zestawu
	// narzedzi tury (skladRozmowy) — jeden byt na dwoch czytelnikow.
	dolozenia *adapterNarzedziSesji
	// zakresyNarzedzi jest tym samym adapterem, który wchodzi do rdzenia jako
	// straż zakresów — jeden byt na dwóch czytelników.
	zakresyNarzedzi *adapterZakresowNarzedzi
	// strazEkspertow czyta zakres eksperta w dwóch miejscach odmowy.
	strazEkspertow StrazEksperta
	// nadajnik to konto nadawcze platformy — nim ida dwa listy systemowe.
	nadajnik nadajnik.Nastawy
}

// nastawyNadajnika przeklada konfiguracje startu na konto nadawcze platformy.
//
// Szyfrowanie jest wlaczone, dopoki Operator jawnie go nie zdejmie: wartosc
// domyslna ma chronic, a nie ulatwiac. Zejscie do rozmowy otwartym tekstem ma
// sens wylacznie dla przekaznika na tej samej maszynie i wymaga jawnego zapisu.
func nastawyNadajnika(k konfiguracja.Konfiguracja) nadajnik.Nastawy {
	startTLS := true
	if k.NadawcaStartTLS != nil {
		startTLS = *k.NadawcaStartTLS
	}
	return nadajnik.Nastawy{
		Host:             k.NadawcaHost,
		Port:             k.NadawcaPort,
		Uzytkownik:       k.NadawcaUzytkownik,
		Sekret:           k.NadawcaSekret,
		Adres:            k.NadawcaAdres,
		NazwaWyswietlana: k.NadawcaNazwa,
		SzyfrujStartTLS:  startTLS,
		WeryfikujTLS:     true,
	}
}

// zlozPorty wypelnia porty rdzenia gotowymi adapterami.
//
// Rozdzial wzgledem Zmontuj idzie po odpowiedzialnosci: Zmontuj sklada byty
// i wiaze je ze soba, a ta funkcja wylacznie przeklada je na porty kontraktu.
func zlozPorty(s skladPortow) Porty {
	// Katalog danych rdzenia odczytany raz i przekazany do obu składów
	// trzymających stan poza bazą: sejfu poświadczeń i magazynu treści
	// biblioteki. Wartość niesie przełącznik `-dane` albo zmienna
	// `DANACO_KATALOG_DANYCH`; przy braku wskazania `konfiguracja.Domyslna()`
	// wstawia katalog domyślny. Konstruktory adapterów wołane bez tej wartości
	// stoją na `konfiguracja.KatalogDanychDomyslny()`, więc pominięcie jej
	// rozdziela stan rdzenia między katalog wskazany a domyślny.
	katalogDanych := s.montaz.Konfiguracja.KatalogDanych
	// Jeden sejf poświadczeń dla kont i punktów dostępu, zbudowany nad
	// skonfigurowanym katalogiem danych, a nie domyślnym. Jedna instancja pod
	// jednym zamkiem — dwa sejfy nad tym samym plikiem ścigałyby się o zapis.
	sejf := dane.NowySejfPlikowy(katalogDanych)
	// Ten sam sejf idzie do warstwy modeli. Kanał API rozwiązuje odwołanie
	// „sejf:<byt>" przez uchwyt pakietowy (models.UstawSejfPoswiadczen), bo kanały
	// powstają fabryką z samego wiersza rejestru i nie mają jak dostać sejfu
	// argumentem. Montaż jest jedyny w procesie i biegnie przed obsługą
	// pierwszego żądania, więc zapis uchwytu wyprzedza wszystkie odczyty.
	models.UstawSejfPoswiadczen(sejf)
	// Katalog akcji i adapter przenoszenia kontekstu powstają przed literałem
	// portów, bo mają po dwóch czytelników: własny port niżej i nakładkę AOD
	// (podpowiedzi biorą się z katalogu akcji, a `aod.context.get` z magazynu
	// kompletu okna). Druga instancja każdego z nich byłaby drugą prawdą o tym
	// samym bycie.
	akcje := rejestrAkcji(s.zycie, s.repozytoria, s.montaz.Dziennik)

	// Centrum powiadomień powstaje przed literałem portów, bo ma dwóch
	// czytelników: port `CentrumPowiadomien` (rodzina notification.*) i drogę
	// wewnętrzną `Zglos`, którą rdzeń wnosi do rejestru zdarzenia, gdy coś
	// zaszło. Druga instancja byłaby drugim rejestrem tego samego bytu.
	// Nastawy idą tym samym adapterem ustawień, co konto nadawcze: o tym, czy
	// klasa zdarzenia w ogóle powiadamia, rozstrzyga sekcja „Powiadomienia".
	centrumPowiadomien := nowyAdapterCentrumPowiadomien(s.repozytoria.CentrumPowiadomien).
		ZNadawca(nowyEmiter(s.szynaZdarzen)).
		ZNastawami(s.ustawienia)
	przenoszenie := nowyAdapterPrzenoszenia(s.nadzorca).
		ZTrwaloscia(s.zycie, s.repozytoria.Konfiguracja, s.montaz.Dziennik).
		ZHistoria(s.dziennikRozmow)

	// Adapter Asystenta powstaje przed literałem portów, bo ma dwóch czytelników:
	// port `Asystent` niżej i nakładkę AOD, która zleceń nie dubluje, tylko woła
	// ten sam moduł. Składacz mostów idzie ten sam, co do rozmowy
	// (`s.rozmowa.mosty`), bo wykaz narzędzi jest jeden; bez `ZMostami` tura
	// zlecenia idzie bez wpisu serwera narzędzi w konfiguracji MCP, a model
	// prowadzący zlecenie nie ma ani jednego narzędzia kontraktu.
	asystent := nowyAdapterAsystenta(s.repozytoria.Asystent).
		ZKanalami(s.kanaly).ZSesjami(s.nadzorca).ZWyjsciem(s.szynaZdarzen, s.zycie).
		ZMowa(s.mowa).ZMostami(s.rozmowa.mosty).ZDomknieciemZastanych()

	// Nakładka dostaje te same byty, którymi jedzie reszta rdzenia: port rozmowy
	// (droga `message.send`), moduł Assistant (droga zlecenia), katalog akcji
	// (podpowiedzi), magazyn kompletu kontekstu okna i katalog urządzeń. Bez
	// kompletu wpięć komendy rodziny `aod.*` odmawiają mimo rejestracji.
	nakladkaAod := nowyAdapterNakladkiAod(s.nadzorca, s.telemetria, s.obecnosc).
		ZRozmowa(s.telemetria.OwinRozmowe(s.rozmowa)).
		ZAsystentem(asystent).
		ZPodpowiedziami(akcje).
		ZKontekstem(przenoszenie).
		ZUrzadzeniami(s.repozytoria.Urzadzenia).
		ZWyciszeniami(s.repozytoria.WyciszeniaNakladki)

	// Jeden adapter okien na dwa porty — Okna i Role. Stoi w zmiennej, a nie
	// dwukrotnie w literale, bo obie rodziny komend muszą czytać okno tym samym
	// kompletem wpięć. Adapter zbudowany bez `ZWieziami` ma `wiezie == nil`,
	// więc `dolozWiezi` (adapter_okna.go) kończy się na pierwszym warunku
	// i `role.list` nie widzi więzi koordynator–wykonawca zapisanej w bazie.
	okna := nowyAdapterOkien(s.nadzorca).ZTrwaloscia(s.trwalosc).ZPetla(s.petla).
		ZPrzerwaniemTury(s.rozmowa.PrzerwijTure, s.rozmowa.CzyTuraWBiegu).
		ZWieziami(s.repozytoria.Przekazania).ZKanalami(s.repozytoria.Kanaly).
		ZModulami(s.repozytoria.Moduly).ZeStrazaEksperta(s.strazEkspertow)

	// Warsztat PDF stoi w zmiennej, bo bierze go także port bezpieczeństwa:
	// dwa adaptery na jednym magazynie zasobów, a nie dwa magazyny.
	warsztatPdf := nowyAdapterPdfStudia().
		ZKatalogiemDanych(katalogDanych).ZZasobami(s.repozytoria.Design)

	// Deklaracje zdolności klientów spisywane w powitaniu. Jedna mapa na rdzeń;
	// czyta ją wyłącznie `launcher.hotkey.*`, żeby odpowiedzieć uczciwie, czy
	// skrót globalny ma kto przechwycić.
	zdolnosciKlientow := noweZdolnosciKlientow()

	return Porty{
		zdolnosciKlienta: zdolnosciKlientow,
		// Kondycja mierzy naprawdę: adres, gniazdo, program i kanał modelu.
		// Stąd trzy zależności ponad magazyn sond — uruchamiacz procesów
		// (sonda programowa), izolacja wraz z katalogiem roboczym (ta sama
		// droga, co każde inne wołanie arsenału) oraz rejestr kanałów (sonda
		// wywołania modelu).
		Kondycja: nowyAdapterKondycji(s.repozytoria.Kondycja).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy).
			ZKanalami(s.kanaly),
		// Alerty liczą miary z magazynów, o których warstwa alertu nie ma prawa
		// wiedzieć sama: ślad wywołań modelu, dziennik błędów i nastawa pułapu
		// kosztu. Bez któregokolwiek reguła na tej mierze nie powstanie —
		// adapter odmawia jej zapisu, zamiast milczeć przy ewaluacji.
		Alerty: nowyAdapterAlertow(s.repozytoria.Alerty).
			ZeZrodlamiMiar(s.repozytoria.Prowenancja, s.repozytoria.Diagnostyka, s.rozstrzygacz).
			ZWyjsciem(nowyEmiter(s.szynaZdarzen)).
			ZCentrumPowiadomien(centrumPowiadomien),
		// Trzy rodziny obsługi tekstu poza jednym oknem. Każda dostaje wyłącznie
		// swój magazyn: nie znają się nawzajem i nie dzielą stanu.
		Schowek:          nowyAdapterSchowka(s.repozytoria.Schowek),
		SkrotyTekstowe:   nowyAdapterSkrotowTekstowych(s.repozytoria.SkrotyTekstowe),
		KontekstyPamieci: nowyAdapterKontekstowPamieci(s.repozytoria.KontekstyPamieci),
		// Wywoływacz trzyma nastawę skrótu w konfiguracji — tam, gdzie mieszka
		// każde inne ustawienie — i czyta deklaracje zdolności klientów, żeby
		// nie obiecywać skrótu, którego nie ma kto przechwycić.
		Wywolywacz: nowyAdapterWywolywacza(s.repozytoria.Konfiguracja, s.rozstrzygacz,
			zdolnosciKlientow),
		// Zajętość okna kontekstu liczy się z treści, która naprawdę pojedzie do
		// modelu: prompt systemowy z portu tożsamości, historia rozmowy okna
		// i wpisy pamięci okna. Granicę okna podaje parametr kanału.
		ZajetoscKontekstu: nowyAdapterZajetosciKontekstu(s.repozytoria.Okna,
			s.repozytoria.Wiadomosci, s.repozytoria.Pamiec, s.kanaly, s.tozsamosc),
		Sesje: nowyAdapterSesji(s.nadzorca).ZTrwaloscia(s.trwalosc).
			ZObecnoscia(s.obecnosc).ZZapewnieniem(s.utrwalacz).
			ZProjektami(s.repozytoria.PrzestrzenRobocza).ZZestawem(s.repozytoria).
			ZZatrzymaniemTur(s.rozmowa.PrzerwijTure),
		Okna:       okna,
		Rozmowa:    s.telemetria.OwinRozmowe(s.rozmowa),
		Ustawienia: s.ustawienia.ZProwenancja(s.rozmowa),
		// Kanały biorą sejf poświadczeń wyłącznie dla `channel.credential.status`
		// — i widzą z niego sam odczyt, bo stan poświadczenia jest pytaniem, czy
		// coś pod odwołaniem leży, a nie prośbą o treść.
		Kanaly:       nowyAdapterKanalow(s.repozytoria.Kanaly, s.kanaly).ZSejfem(sejf),
		Kolejki:      s.kolejki,
		Przenoszenie: przenoszenie,
		// Warstwa mobilna jest rozszerzeniem monitora, nie osobnym portem: rodzina
		// `mobile.*` czyta procesy z tego samego rejestru telemetrii, którym jedzie
		// `monitor.status`. Bez tego ogniwa asercja portu w `handlers_mobile.go`
		// nie przechodzi i wszystkie trzy komendy odmawiają, choć adapter jest
		// napisany — dlatego wpięcie idzie tutaj, zaraz za telemetrią.
		Nawigacja: nowyAdapterNawigacji(s.repozytoria, s.nadzorca, s.ustawienia, s.klienci,
			s.rozmowa.CzyTuraWBiegu).ZObecnoscia(s.obecnosc, s.biegi).
			ZTelemetriaProcesow(s.telemetria).
			ZWarstwaMobilna(s.telemetria.OwinRozmowe(s.rozmowa), s.kolejki),
		WiazanieSesji:      nowyAdapterWiazaniaSesji(s.repozytoria, s.nadzorca, s.klienci),
		Akcje:              akcje,
		CentrumPowiadomien: centrumPowiadomien,
		KatalogUstawien:    nowyAdapterKatalogUstawien(s.repozytoria.KatalogUstawien),
		PunktyDostepu:      nowyAdapterPunktowDostepu(s.repozytoria.PunktyDostepu).ZSejfem(sejf),
		NadaniaDostepu: nowyAdapterNadanDostepu(s.repozytoria.Nadania, s.repozytoria.Okna,
			s.repozytoria.PunktyDostepu),
		Konta:     nowyAdapterKont(s.repozytoria.Konta).ZSejfem(sejf),
		Tozsamosc: s.tozsamosc,
		// Moduł Agents stoi na tym samym rejestrze kanałów co okna rozmowy
		// i na tym samym katalogu mostów co składacz `mcpServers`.
		Agenci: nowyAdapterAgentow(s.repozytoria.Agenci).
			ZKanalami(s.kanaly).
			ZPunktamiDostepu(s.repozytoria.PunktyDostepu).
			// Warstwy promptu eksperta. Bez tego ogniwa komendy
			// agent.layer.* i agent.plugin.* odmawiaja.
			ZWarstwami(s.repozytoria.WarstwyAgenta),
		// Moduł Workspace bierze instrukcje warstwowe z tego rozstrzygacza —
		// jednego na całą platformę — a bibliotekę projektu z tego samego
		// ustalacza katalogu roboczego, którym jedzie sesja.
		PrzestrzenRobocza: nowyAdapterPrzestrzeniRoboczej(s.repozytoria.PrzestrzenRobocza).
			ZInstrukcjami(s.repozytoria.Konfiguracja, s.rozstrzygacz).
			ZKatalogiem(s.katalogRoboczy).
			// Wydobycie treści pliku projektu idzie tym samym warsztatem, co
			// `document.text.extract`: warstwa tekstowa dokumentu, a po jej
			// braku rozpoznanie pisma. Drugiego czytnika dokumentów w rdzeniu
			// nie ma i nie ma po co być.
			//
			// Kolejność ogniw nie jest dowolna: `ZPamiecia` oddaje adapter
			// OPAKOWANY o rodzinę `memory.*` i musi stać ostatnie. Ogniwo
			// dopięte po nim oddawałoby adapter wewnętrzny, a port straciłby
			// pamięć — rdzeń nie miałby wtedy pięciu komend `memory.*`.
			ZDokumentamiWorkspace(nowyAdapterNarzedziDokumentu(injection.UruchamiaczOkien()).
				ZKatalogiemDanych(katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
				ZZasobamiDesignu(s.repozytoria.Design).
				ZBiblioteka(s.repozytoria.Biblioteka)).
			// `ZWylaczeniami` idzie PO `ZPamiecia`, bo wyłączenia są ogniwem
			// adaptera opakowanego — bez tej kolejności rodzina
			// `memory.disable.*` odmawiałaby z nazwą niewpiętego magazynu.
			ZPamiecia(s.repozytoria.Pamiec, s.repozytoria.Sesje).
			ZWylaczeniami(s.repozytoria.WylaczeniaPamieci),
		Terminal:    s.terminal.ZDziennikiemWyjscia(),
		Diagnostyka: s.diagnostyka,
		// Ślad wywołań i rozliczenie zużycia stoją na jednym magazynie:
		// dwa porty, jedno repozytorium.
		// Prowenancja bierze ponadto rejestr kanałów — wyłącznie dla powtórzenia
		// wywołania (`provenance.call.replay`), które jest nowym wywołaniem
		// kanału, a nie odczytem śladu. Cztery komendy odczytu go nie dotykają.
		Prowenancja: nowyAdapterProwenancji(s.repozytoria.Prowenancja).ZKanalami(s.kanaly),
		Zuzycie:     nowyAdapterZuzycia(s.repozytoria.Prowenancja),
		// Warsztat PDF nie bierze uruchamiacza procesów ani zasad izolacji:
		// pracuje biblioteką wkompilowaną w rdzeń i nie startuje ani jednego
		// procesu potomnego.
		WarsztatPdf: warsztatPdf,
		// Bezpieczeństwo sięga po ten sam warsztat, bo materiał wchodzi tą samą
		// drogą, oraz po repozytorium Studia, bo rozpoznanie danych wrażliwych
		// czyta treść dokumentu, a nie zasób magazynu.
		BezpieczenstwoDokumentu: nowyAdapterBezpieczenstwaStudia(warsztatPdf).
			ZeStudiem(s.repozytoria.Studio),
		// Oba porty muszą być wypełnione: strażniki `if d == nil { return }`
		// w `zarejestrujDevelopera` i `zarejestrujDebate` wychodzą przed
		// rejestracją, więc port pusty daje odpowiedź `*.unknown` na komendy,
		// których klient ma komplet.
		Developer: s.developer,
		// Debata dostaje ponadto katalog danych i arsenał: pierwszy pod magazyn
		// wydanych transkryptów, grafów i nagrań, drugi pod dwie czynności
		// wymagające programu serwerowego — zamianę transkryptu na dokument
		// biurowy (Pandoc) i odsłuch debaty (silnik mowy). Reszta modułu ich nie
		// dotyka: debata, analiza i głosowanie nie startują ani jednego procesu.
		Debata: nowyAdapterDebaty(s.zycie, s.repozytoria.Roundtable).
			ZKanalami(s.kanaly).
			ZKatalogiemDanych(katalogDanych).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz),
		// Moduł Automations buduje definicje, a wykonuje je tym samym adapterem
		// kolejek, którym jedzie domena kolejek — silnik jest jeden.
		// Instancja powstaje w złożeniu modułów (montaz_moduly.go), bo tą samą
		// dzieli budzik harmonogramu: gdyby port i budzik miały osobne adaptery,
		// przypięcia obserwatorów przebiegów rozjechałyby się na dwie mapy.
		Automatyki: s.automatyki,
		// Biblioteka, Studio, Przegladarka i Design dostają wyłącznie swoje
		// repozytoria — nie znają się nawzajem i nie dzielą stanu.
		//
		// Biblioteka dostaje ponadto katalog danych — ten sam, nad którym stoi
		// sejf poświadczeń wyżej. Baza trzyma wyłącznie odwołanie do treści,
		// więc bajty wgranych plików muszą leżeć tam, gdzie reszta stanu rdzenia.
		// Arsenał i rejestr kanałów są modułowi potrzebne w dwóch miejscach:
		// pomiar czasu trwania nagrania przy odczycie metadanych osadzonych
		// oraz klasyfikacja wsadowa. Brak któregokolwiek zawęża zakres tych
		// dwóch czynności, a nie wyłącza modułu.
		Biblioteka: nowyAdapterBiblioteki(s.repozytoria.Biblioteka).
			ZKatalogiemDanych(katalogDanych).
			ZNarzedziami(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy).
			ZKanalamiModelu(s.kanaly),
		// Studio dostaje ten sam rejestr kanalow, ktorym jedzie okno rozmowy
		// i modul Roundtable — drugiego silnika modelu nie ma nigdzie. Magazyn
		// zasobów jest ten sam, którym jedzie warsztat PDF i moduł Design:
		// archiwum historii, paczka redakcyjna i strony podglądu są zasobami tej
		// samej platformy, więc leżą w jednym miejscu. Repozytorium Library
		// wchodzi wyłącznie do odczytu — `studio.diff.source` zestawia dokument
		// roboczy z materiałem wejściowym i bez niego nie miałoby czego czytać.
		Studio: nowyAdapterStudia(s.repozytoria.Studio).
			ZKanalami(s.kanaly, s.nadzorca.Rejestr()).
			ZNarzedziami(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy).
			ZZasobami(s.repozytoria.Design, katalogDanych).
			ZBiblioteka(s.repozytoria.Biblioteka),
		// Przeglądarka dostaje trzy rzeczy ponad własne repozytorium: katalog
		// danych (magazyn bajtów zrzutów, archiwów i rejestrów sieciowych — ten
		// sam korzeń, co magazyn biblioteki i zasobów Designu), uruchamiacz
		// procesów wraz z izolacją (silnik Chromium prowadzony protokołem CDP)
		// oraz nic więcej: moduł nie zna innych modułów i nie dzieli z nimi stanu.
		Przegladarka: nowyAdapterPrzegladarki(s.repozytoria.Przegladanie).
			ZMagazynem(katalogDanych).
			ZSilnikiem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		// Design bierze rejestr kanałów tym samym sposobem co Roundtable i Studio
		// (ZKanalami): most do modelu dla pracy tekstowej (budowa promptu,
		// opis zasobu). Rdzeń nie generuje obrazów — rejestr służy tylko
		// operacjom słownym modułu (patrz `adapter_modul_design.go`).
		// Design dostaje ponadto katalog danych — ten sam, którym jadą sejf
		// poświadczeń i magazyn biblioteki, bo `design.asset.upload` trzyma bajty
		// zasobów poza bazą.
		// Design dostaje też sejf — ten sam, co Konta i PunktyDostepu — bo klucze
		// darmowych baz zdjęciowych (`design.stock.*`) leżą w nim pod bytem
		// `design.stock.<dostawca>`. Moduł wyłącznie je CZYTA.
		// Design dostaje na koniec drogę ODCZYTU PISMA (ZOdczytemPisma) —
		// uruchamiacz i bramę izolacji dla jednej czynności:
		// `design.mockup.import` czyta treść napisów ze zrzutu programem pakietu
		// serwera, bo czytnika liter w czystym Go nie ma. Reszta modułu procesów
		// nie startuje i pilnuje tego zapora.
		Design: nowyAdapterDesignu(s.repozytoria.Design).ZKanalami(s.kanaly).
			ZKatalogiemDanych(katalogDanych).ZSejfem(sejf).
			ZOdczytemPisma(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		// Research dostaje dodatkowo rejestr kanałów, bo streszczanie
		// i porównanie źródeł to operacja modelu.
		Badania: nowyAdapterBadan(s.repozytoria.Badania).ZKanalami(s.kanaly).
			ZKatalogiemDanych(katalogDanych).
			ZDokumentami(nowyAdapterNarzedziDokumentu(injection.UruchamiaczOkien()).
				ZKatalogiemDanych(katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
				ZZasobamiDesignu(s.repozytoria.Design).
				ZBiblioteka(s.repozytoria.Biblioteka)).
			ZMowa(s.mowa),
		// Assistant dostaje wykonawcę zleceń: rejestr kanałów (droga modelu),
		// nadzorcę sesji (okno zlecenia) i nadajnik z kontekstem życia (strumień
		// tury i assistant.action.changed). Bez nich zlecenie zostaje `queued`
		// zamiast się wykonać. Rozpoznanie mowy (`ZMowa`) idzie tym samym
		// adapterem, który wypełnia port `Mowa` niżej — silnik jest jeden.
		Asystent: asystent,
		// Apps dostaje rejestr okien, bo wdrożenie bez istniejącego okna nie ma
		// przestrzeni roboczej, z której miałoby cokolwiek wziąć — bez rejestru
		// komenda zakłada przebieg dla okna, którego nie ma, i kończy go powodem
		// o pustym warsztacie zamiast o braku okna.
		// Apps dostaje ponadto katalog danych (magazyn bajtów eksportu, artefaktów
		// i pakietów oraz sejf, z którego bierze się klucz wydawcy) i rejestr
		// pozycji katalogu — `apps.package.publish` publikuje do TEGO rejestru,
		// nie do drugiego obok niego.
		Aplikacje: s.aplikacje,
		// Komponenty własne dostają swój rejestr oraz trzy
		// magazyny modułowe, w których `component.create` zakłada byt docelowy:
		// projekty, ekspertów i automatyki. Czwartego — profili asystenta —
		// platforma nie ma, więc rodzaj `assistant` odmawia z powodem, zamiast
		// zakładać kafel wskazujący na nic.
		Komponenty: nowyAdapterKomponentow(s.repozytoria.Komponenty).
			ZMagazynamiModulow(s.repozytoria.PrzestrzenRobocza, s.repozytoria.Agenci,
				s.repozytoria.Automatyki),
		// Translate dostaje rejestr kanałów (ZKanalami) — most do modelu, którym
		// tłumaczenie faktycznie woła kanał zamiast zakładać puste panele. Wybór
		// kanału niesie już żądanie kontraktu (pole `channelId`); kanał wskazany,
		// a nieznany albo nieczynny kończy się odmową nazwaną, nie cichym zejściem
		// na kanał domyślny.
		//
		// ZSynteza wpina silnik syntezy mowy dla `translate.speech.synthesize`:
		// ten sam uruchamiacz i te same dwa źródła izolacji, którymi jadą Terminal,
		// Developer i silnik rozpoznawania mowy, plus katalog danych rdzenia jako
		// miejsce na nagrania — ten sam, który dostają sejf poświadczeń i magazyn
		// treści biblioteki wyżej.
		// ZWytworami wpina drogę `translate.artifact.publish`: repozytorium
		// biblioteki (wiersz pliku) i katalog danych rdzenia (magazyn bajtów pod
		// sumą kontrolną). Bez niej wytwór nie miałby gdzie leżeć ani czym się
		// zgłosić reszcie platformy, więc sama ta komenda odmawia nazywając brak.
		Tlumaczenie: nowyAdapterTlumaczenia(s.repozytoria.Tlumaczenia).ZKanalami(s.kanaly).
			ZSynteza(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy, katalogDanych).
			ZWytworami(s.repozytoria.Biblioteka, katalogDanych),
		// Przekazanie okna zakłada pozycję kolejki, więc bierze ten sam adapter
		// kolejek, którym jedzie domena kolejek i moduł Automations — silnik
		// jest jeden. Bez niego `window.handoff` odmawia wprost.
		// Dwa wiązania, bo przekazanie okna domyka lukę między kontraktem
		// a schematem. Kontrakt niesie identyfikatory zewnętrzne (sesji, okien),
		// a schemat wiąże klucze wewnętrzne i słowniki modułów oraz kanałów —
		// bez tych czterech repozytoriów adapter nie umiałby ani odnaleźć okna
		// wskazanego przez klienta, ani oddać okna w kształcie kontraktu.
		PrzekazanieOkna: nowyAdapterPrzekazaniaOkna(s.repozytoria.Przekazania).
			ZKolejkami(s.kolejki).
			ZBytamiOkien(s.repozytoria.Okna, s.repozytoria.Sesje,
				s.repozytoria.Moduly, s.repozytoria.Kanaly).
			ZKatalogiemAkcji(s.repozytoria.Akcje),
		// Izolacja bierze rozstrzygacz polityki, nie samo repozytorium: profil
		// mówi, co okno widzi z historii, pamięci i kontekstu sąsiada, a to
		// rozstrzyga się warstwowo, nie jednym wierszem tabeli.
		Izolacja: nowyAdapterIzolacji(s.repozytoria, s.rozstrzygacz).ZRozgloszeniem(s.szynaZdarzen),
		// Uwierzytelnianie stoi na repozytorium budowanym przez zestaw danych
		// (dane/zestaw.go).
		Uwierzytelnianie: nowyAdapterUwierzytelnienia(s.repozytoria.Uwierzytelnienie).
			// Ten sam sejf, co Konta i PunktyDostepu — jedna instancja pod jednym
			// zamkiem. Bez tego ogniwa sekret bramki nie ma gdzie lezec, wiec
			// logowanie odmawia zawsze.
			ZSejfem(sejf).
			// Tozsamosc wlasciciela: login, adres uwierzytelniajacy i drogi
			// potwierdzenia. Bez tego ogniwa rejestracja odmawia, bo konta nie ma
			// gdzie zapisac.
			ZKontem(s.repozytoria.KontoWlasciciela).
			// Konto nadawcze PLATFORMY — nie skrzynka Operatora. Nim ida dwa listy
			// systemowe: potwierdzenie adresu i droga odzyskania konta.
			ZNadajnikiem(s.nadajnik).
			// Ten sam adapter ustawien, ktorym idzie kazda inna nastawa platformy.
			// Konto nadawcze zapisane w oknie Konfiguracji przeslania to ze startu,
			// wiec pomylke w adresie serwera poczty naprawia sie bez zatrzymywania
			// rdzenia (nastawy_nadajnika.go).
			ZNastawamiPlatformy(s.ustawienia),
		// Urzadzenia stoja na tym samym repozytorium: urzadzeniem konta jest to,
		// ktore weszlo przez bramke, wiec wykaz bierze sie z sesji bramki.
		Urzadzenia: nowyAdapterUrzadzen(s.repozytoria.Uwierzytelnienie),
		// Nakladka AOD nie dostaje repozytorium, tylko te same trzy byty, ktorymi
		// jedzie reszta rdzenia: nadzorce sesji, telemetrie postepu i rejestr
		// obecnosci. Nakladka pokazuje stan biegnacej pracy — wlasna tabela bylaby
		// druga prawda o tym samym.
		NakladkaAod: nakladkaAod,
		// Katalog rozszerzen — warstwa danych wystawia go metoda, nie polem
		// (dane/extension.go), bo rejestr nie trzyma stanu poza wskaznikiem
		// na wspolna pamiec zapytan.
		// Katalog dostępu i biblioteka ekspertów idą razem z rozgłoszeniem, bo bez
		// nich adapter odmawia `internal_error` każdemu żądaniu niosącemu
		// `accessPointId` albo `agentId`; wskazanie puste przechodzi.
		// Katalog danych wchodzi tą samą drogą, co do modułu Apps: paczka
		// przesłana instalacją Personal i dziennik piaskownicy są bajtami na
		// dysku, nie wpisem o bajtach.
		Rozszerzenia: nowyAdapterRozszerzen(s.repozytoria.Rozszerzenia()).
			ZKatalogiemDostepu(s.repozytoria.PunktyDostepu, s.repozytoria.Agenci).
			ZMagazynemRozszerzen(katalogDanych).
			ZRozgloszeniem(s.szynaZdarzen),
		// Role okien bierze caly zestaw, nie samo repozytorium: adapter musi
		// dosiegnac takze przekazan (wiez koordynatora) i adresow okien. Stoi na
		// tym samym adapterze okien, ktory wypelnia port Okna — drugi bylby
		// druga prawda o oknie.
		Role: okna.ZRolami(s.repozytoria),
		// Układ sekcji paneli bierze repozytorium wystawione metodą, nie polem
		// (dane/panele.go) — tak samo jak role i rozszerzenia, bo rejestr nie
		// trzyma stanu poza wskaźnikiem na wspólną pamięć zapytań.
		Panele: nowyAdapterSekcjiPaneli(s.repozytoria.SekcjePaneli()),
		// Historia rozmowy okna bierze jedno repozytorium — swoje. Okien ani sesji
		// nie dobiera: wiąże je zapytanie po identyfikatorze kontraktowym, a drugi
		// czytelnik tych samych wierszy byłby drugą prawdą o wypowiedzi.
		Historia: NowyPortHistorii(s.repozytoria.Historia),

		// Podagenci biorą komplet wiązań z jednego miejsca
		// (adapter_modul_orkiestracja_zlozenie.go: zlozPodagentow): repozytorium
		// podagentów i okien, biegi orkiestracji, ten sam adapter kolejek co domena
		// kolejek i Automations oraz ocenę żywotności z rejestru procesów
		// sesji. Montaż nie zna pakietu `podagenci` — wiedzę o nim trzyma adapter,
		// który jako jedyny go używa.
		Podagenci: zlozPodagentow(s),

		// Zespoły biorą własne repozytorium oraz dziennik rdzenia — ten sam,
		// którym mówi reszta montażu. Skład, z którego wypadł ekspert usunięty
		// albo zarchiwizowany, wraca do klienta krótszy, a kontrakt nie ma pola,
		// którym dałoby się o tym powiedzieć.
		Zespoly: nowyAdapterZespolow(s.repozytoria.Zespoly).
			ZDziennikiem(s.montaz.Dziennik),
		// Historia i archiwum eksperta biorą trzy repozytoria: dwa widoki wersji
		// oraz bibliotekę ekspertów, która służy wyłącznie oddaniu eksperta
		// w kształcie kontraktu po zmianie wersji.
		WersjeEksperta: NowyPortWersjiEksperta(s.repozytoria.WersjeAgenta,
			s.repozytoria.ArchiwumAgentow, s.repozytoria.Agenci),
		// Zakres działania eksperta bierze pięć źródeł: własne repozytorium
		// zakresu, bibliotekę (kształt kontraktu po zmianie), historię wersji
		// (podgląd migawki), katalog modułów (sprawdzenie wskazania) i katalog
		// punktów dostępu (konfiguracja instancji konektora). Rozstrzygacz
		// dokłada dziedziczenie izolacji do polityki efektywnej — ten sam,
		// którym jedzie okno konfiguracji punktów izolacji.
		ZakresEksperta: NowyPortZakresuEksperta(s.repozytoria.ZakresAgenta,
			s.repozytoria.Agenci, s.repozytoria.WersjeAgenta, s.repozytoria.Moduly,
			s.repozytoria.PunktyDostepu, s.rozstrzygacz),
		// Zakresy narzędzi stoją na własnym repozytorium i na katalogu profili
		// asystenta — profil pusty w żądaniu bierze profil domyślny.
		ZakresyNarzedzi: s.zakresyNarzedzi,
		// Silnik mowy powstaje w złożeniu modułów (montaz_moduly.go), bo dzieli
		// uruchamiacz procesów i oba źródła izolacji z Terminalem i Developerem,
		// a tę samą instancję bierze moduł Assistant wyżej — drugi silnik byłby
		// drugą prawdą o tym, czy platforma rozpoznaje mowę.
		Mowa: s.mowa,
		// Wskaźnik znaczenia stoi na tym samym uruchamiaczu i tych samych dwóch
		// źródłach izolacji, co Terminal, Developer i silnik mowy: pomocnik
		// osadzeń jest procesem drzewa jak każdy inny. Katalog danych jedzie ten
		// sam, którym jadą sejf, magazyn biblioteki i zasoby Designu — tam
		// wykłada się pomocnik i tam lądują wagi modelu. Składnica wektorów
		// bierze `*sql.DB` z montażu, tak samo jak dziennik transkrypcji
		// i dziennik doradcy, bo `dane.Zestaw` uchwytu bazy nie wystawia; stan
		// trzyma tabela `fragment_wiedzy`. Źródła treści idą przez repozytoria
		// modułu Library i historii — drugiej drogi do tych wierszy nie ma.
		Wiedza: nowyAdapterWiedzy(injection.UruchamiaczOkien(), katalogDanych).
			ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
			ZeSkladnica(skladnicaWiedzy(s.montaz)).
			ZeZrodlami(s.repozytoria.Biblioteka, s.repozytoria.Historia),
		// Narzędzia obrazu biorą repozytorium Designu i katalog danych — ten sam,
		// którym jadą sejf, magazyn biblioteki i zasoby Designu — a arsenałem
		// jedzie ten sam uruchamiacz i te same dwa źródła izolacji, co Terminal,
		// Developer i silniki mowy.
		NarzedziaObrazu: nowyAdapterNarzedziObrazu(s.repozytoria.Design, katalogDanych).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		// Silniki neuronowe obrazu (`image.upscale`, `image.background.remove`)
		// stoją na tym samym zapleczu, co rodzina `image.*` wyżej — to samo
		// repozytorium Designu, ten sam katalog danych, ten sam uruchamiacz
		// i te same dwa źródła izolacji. Zaplecze jest bezstanowe, więc druga
		// instancja niczego nie rozdwaja; dokładają się do niego wyłącznie wagi
		// sieci i granice czasu liczone w minutach.
		NarzedziaObrazuModelu: nowyAdapterNarzedziObrazuModelu(
			nowyAdapterNarzedziObrazu(s.repozytoria.Design, katalogDanych).
				ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy)),
		// Narzędzia mediów stoją na tej samej trójce co narzędzia obrazu wyżej:
		// repozytorium Designu (wynik jest zasobem), katalog danych (magazyn
		// bajtów) i arsenał wpięty tym samym uruchamiaczem oraz tymi samymi
		// dwoma źródłami izolacji.
		NarzedziaMedia: nowyAdapterNarzedziMediow(injection.UruchamiaczOkien(),
			s.repozytoria.Design, katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy),
		// Narzędzia dokumentowe stoją na tej samej trójce co dwie rodziny wyżej:
		// arsenał wpięty tym samym uruchamiaczem i tymi samymi dwoma źródłami
		// izolacji, katalog danych jako magazyn bajtów wyniku oraz repozytorium
		// Designu — tu wyłącznie do odczytu, żeby rozwiązać `assetId` żądania na
		// bajty.
		Dokumenty: nowyAdapterNarzedziDokumentu(injection.UruchamiaczOkien()).
			ZKatalogiemDanych(katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
			ZZasobamiDesignu(s.repozytoria.Design).
			ZBiblioteka(s.repozytoria.Biblioteka),
		// Narzędzia archiwum stoją na tej samej trójce co trzy rodziny wyżej:
		// repozytorium Designu (wytworzone archiwum jest zasobem, a wskazane
		// zasoby są materiałem pakowania), katalog danych (magazyn bajtów oraz
		// podstawa katalogów pracy pośredniej) i arsenał wpięty tym samym
		// uruchamiaczem oraz tymi samymi dwoma źródłami izolacji. Katalog
		// roboczy wchodzi tu podwójnie: raz jako obszar wołania `7z`, raz jako
		// jedyny układ odniesienia dla ścieżek żądania — bez niego `archive.*`
		// nie miałyby względem czego rozstrzygać `targetPath`.
		NarzedziaArchiwum: nowyAdapterNarzedziArchiwum(s.repozytoria.Design, katalogDanych).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		// Poczta bierze repozytorium skrzynek wystawione metodą, nie polem
		// (dane/poczta_skrzynki.go) — tak samo jak role, rozszerzenia i sekcje
		// paneli. Trzy dalsze wiązania są tymi samymi bytami, którymi jedzie
		// reszta rdzenia: ten sam sejf poświadczeń, co konta, punkty dostępu
		// i bramka, bo poświadczenie skrzynki nie ma innej drogi niż sejf; ten
		// sam katalog danych, co magazyn biblioteki i zasoby Designu; to samo
		// repozytorium Designu, do którego pisze `design.asset.upload` — bo
		// załącznik listu wciągnięty osobną drogą byłby zasobem, którego arsenał
		// obrazu i dokumentów nie widzi.
		Poczta: nowyAdapterPoczty(s.repozytoria.SkrzynkiOperatora()).
			ZSejfem(sejf).ZKatalogiemDanych(katalogDanych).
			ZZasobami(s.repozytoria.Design),
		// Doraźne dołożenia sesji. Adapter powstaje w montażu (montaz.go), bo ten
		// sam byt wnosi dołożenia do zestawu narzędzi tury: port oddaje je
		// Operatorowi, składacz tury oddaje je modelowi, a prawda ma być jedna.
		// Trzy wiązania niesie już konstruktor — dołożenia sesji (migracja 122),
		// wykaz sesji, bo dołożenie żyje w stanie sesji i bez przekładu
		// identyfikatora kontraktowego nie ma gdzie usiąść, oraz katalog
		// rozszerzeń, z którego bierze się druga połowa wykazu po ukośniku.
		NarzedziaSesji: s.dolozenia,
		// Doradca powstaje w złożeniu modułów (montaz_moduly.go) razem z mową:
		// stoi na tym samym rejestrze kanałów, którym jedzie okno rozmowy,
		// i na dzienniku bazy.
		Doradcy:  s.doradcy,
		Nadajnik: s.szynaZdarzen,
		Nasluch:  s.nasluch,
		Dziennik: s.montaz.Dziennik,
	}
}
