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

// skladPortow jest kompletem bytów, z których powstają porty rdzenia — strukturą zamiast dwudziestu pozycyjnych argumentów podatnych na cichą pomyłkę kolejności.
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
	// Aplikacje trzeba zwolnić osobno: `apps.preview.start` podnosi nasłuch HTTP poza cyklem rdzenia.
	aplikacje    *adapterAplikacji
	diagnostyka  *adapterDiagnostyki
	kolejki      *adapterKolejek
	automatyki   *adapterAutomatyk
	mowa         *adapterMowy
	doradcy      *adapterDoradcow
	szynaZdarzen Nadajnik
	nasluch      nasluchTransportu
	// dolozenia to adapter dołożeń sesji współdzielony z zestawem narzędzi tury.
	dolozenia *adapterNarzedziSesji
	// zakresyNarzedzi to adapter straży zakresów współdzielony z zestawem narzędzi tury.
	zakresyNarzedzi *adapterZakresowNarzedzi
	// strazEkspertow czyta zakres eksperta w dwóch miejscach odmowy.
	strazEkspertow StrazEksperta
	// nadajnik to konto nadawcze platformy — nim ida dwa listy systemowe.
	nadajnik nadajnik.Nastawy
}

// nastawyNadajnika przekłada konfigurację startu na konto nadawcze platformy; szyfrowanie jest domyślnie włączone i wymaga jawnego zdjęcia.
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
	// Katalog danych czytany raz z konfiguracji — pominięcie rozdzieliłoby stan między katalogi.
	katalogDanych := s.montaz.Konfiguracja.KatalogDanych
	// Jeden sejf poświadczeń nad katalogiem danych skonfigurowanym: druga instancja ścigałaby się o zapis.
	sejf := dane.NowySejfPlikowy(katalogDanych)
	// Ten sam sejf idzie do warstwy modeli uchwytem pakietowym, zapisanym przed pierwszym żądaniem.
	models.UstawSejfPoswiadczen(sejf)
	// Katalog akcji i przenoszenie kontekstu powstają wcześniej — mają po dwóch czytelników.
	akcje := rejestrAkcji(s.zycie, s.repozytoria, s.montaz.Dziennik)

	// Centrum powiadomień powstaje wcześniej, bo ma dwóch czytelników: port i drogę wewnętrzną `Zglos`.
	centrumPowiadomien := nowyAdapterCentrumPowiadomien(s.repozytoria.CentrumPowiadomien).
		ZNadawca(nowyEmiter(s.szynaZdarzen)).
		ZNastawami(s.ustawienia)
	przenoszenie := nowyAdapterPrzenoszenia(s.nadzorca).
		ZTrwaloscia(s.zycie, s.repozytoria.Konfiguracja, s.montaz.Dziennik).
		ZHistoria(s.dziennikRozmow)

	// Adapter Asystenta powstaje wcześniej, bo ma dwóch czytelników: port `Asystent` i nakładkę AOD.
	asystent := nowyAdapterAsystenta(s.repozytoria.Asystent).
		ZKanalami(s.kanaly).ZSesjami(s.nadzorca).ZWyjsciem(s.szynaZdarzen, s.zycie).
		ZMowa(s.mowa).ZMostami(s.rozmowa.mosty).ZDomknieciemZastanych()

	// Nakładka AOD dostaje te same byty co reszta rdzenia; bez kompletu wpięć rodzina `aod.*` odmawia.
	nakladkaAod := nowyAdapterNakladkiAod(s.nadzorca, s.telemetria, s.obecnosc).
		ZRozmowa(s.telemetria.OwinRozmowe(s.rozmowa)).
		ZAsystentem(asystent).
		ZPodpowiedziami(akcje).
		ZKontekstem(przenoszenie).
		ZUrzadzeniami(s.repozytoria.Urzadzenia).
		ZWyciszeniami(s.repozytoria.WyciszeniaNakladki)

	// Jeden adapter okien zasila porty Okna i Role, by obie rodziny czytały okno tym kompletem wpięć.
	okna := nowyAdapterOkien(s.nadzorca).ZTrwaloscia(s.trwalosc).ZPetla(s.petla).
		ZPrzerwaniemTury(s.rozmowa.PrzerwijTure, s.rozmowa.CzyTuraWBiegu).
		ZWieziami(s.repozytoria.Przekazania).ZKanalami(s.repozytoria.Kanaly).
		ZModulami(s.repozytoria.Moduly).ZeStrazaEksperta(s.strazEkspertow)

	// Warsztat PDF stoi w zmiennej, bo bierze go też port bezpieczeństwa — jeden magazyn zasobów.
	warsztatPdf := nowyAdapterPdfStudia().
		ZKatalogiemDanych(katalogDanych).ZZasobami(s.repozytoria.Design)

	// Deklaracje zdolności klientów spisane w powitaniu; jedna mapa czytana przez `launcher.hotkey.*`.
	zdolnosciKlientow := noweZdolnosciKlientow()

	return Porty{
		zdolnosciKlienta: zdolnosciKlientow,
		// Kondycja mierzy adres, gniazdo, program i kanał modelu, stąd trzy zależności ponad magazyn sond.
		Kondycja: nowyAdapterKondycji(s.repozytoria.Kondycja).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy).
			ZKanalami(s.kanaly),
		// Alerty liczą miary z magazynów spoza warstwy alertu; brak jednego — adapter odmawia zapisu.
		Alerty: nowyAdapterAlertow(s.repozytoria.Alerty).
			ZeZrodlamiMiar(s.repozytoria.Prowenancja, s.repozytoria.Diagnostyka, s.rozstrzygacz).
			ZWyjsciem(nowyEmiter(s.szynaZdarzen)).
			ZCentrumPowiadomien(centrumPowiadomien),
		// Trzy rodziny obsługi tekstu dostają wyłącznie swój magazyn i nie dzielą stanu.
		Schowek:          nowyAdapterSchowka(s.repozytoria.Schowek),
		SkrotyTekstowe:   nowyAdapterSkrotowTekstowych(s.repozytoria.SkrotyTekstowe),
		KontekstyPamieci: nowyAdapterKontekstowPamieci(s.repozytoria.KontekstyPamieci),
		// Wywoływacz trzyma skrót w konfiguracji i czyta zdolności klientów, by nie obiecać nic bez odbiorcy.
		Wywolywacz: nowyAdapterWywolywacza(s.repozytoria.Konfiguracja, s.rozstrzygacz,
			zdolnosciKlientow),
		// Zajętość kontekstu liczy się z treści jadącej do modelu: prompt systemowy, historię i pamięć okna.
		ZajetoscKontekstu: nowyAdapterZajetosciKontekstu(s.repozytoria.Okna,
			s.repozytoria.Wiadomosci, s.repozytoria.Pamiec, s.kanaly, s.tozsamosc),
		Sesje: nowyAdapterSesji(s.nadzorca).ZTrwaloscia(s.trwalosc).
			ZObecnoscia(s.obecnosc).ZZapewnieniem(s.utrwalacz).
			ZProjektami(s.repozytoria.PrzestrzenRobocza).ZZestawem(s.repozytoria).
			ZZatrzymaniemTur(s.rozmowa.PrzerwijTure),
		Okna:       okna,
		Rozmowa:    s.telemetria.OwinRozmowe(s.rozmowa),
		Ustawienia: s.ustawienia.ZProwenancja(s.rozmowa),
		// Kanały biorą sejf wyłącznie do odczytu stanu poświadczenia dla `channel.credential.status`.
		Kanaly:       nowyAdapterKanalow(s.repozytoria.Kanaly, s.kanaly).ZSejfem(sejf),
		Kolejki:      s.kolejki,
		Przenoszenie: przenoszenie,
		// Warstwa mobilna rozszerza monitor, czytając ten sam rejestr telemetrii co `monitor.status`.
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
		// Moduł Agents stoi na rejestrze kanałów okien rozmowy i katalogu mostów składacza `mcpServers`.
		Agenci: nowyAdapterAgentow(s.repozytoria.Agenci).
			ZKanalami(s.kanaly).
			ZPunktamiDostepu(s.repozytoria.PunktyDostepu).
			// Warstwy promptu eksperta. Bez tego ogniwa komendy
			// agent.layer.* i agent.plugin.* odmawiaja.
			ZWarstwami(s.repozytoria.WarstwyAgenta),
		// Moduł Workspace bierze instrukcje z rozstrzygacza platformy i bibliotekę z katalogu roboczego sesji.
		PrzestrzenRobocza: nowyAdapterPrzestrzeniRoboczej(s.repozytoria.PrzestrzenRobocza).
			ZInstrukcjami(s.repozytoria.Konfiguracja, s.rozstrzygacz).
			ZKatalogiem(s.katalogRoboczy).
			// Wydobycie treści idzie warsztatem `document.text.extract`; `ZPamiecia` musi stać ostatnim ogniwem.
			ZDokumentamiWorkspace(nowyAdapterNarzedziDokumentu(injection.UruchamiaczOkien()).
				ZKatalogiemDanych(katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
				ZZasobamiDesignu(s.repozytoria.Design).
				ZBiblioteka(s.repozytoria.Biblioteka)).
			// `ZWylaczeniami` idzie po `ZPamiecia` — inna kolejność zostawia `memory.disable.*` bez magazynu.
			ZPamiecia(s.repozytoria.Pamiec, s.repozytoria.Sesje).
			ZWylaczeniami(s.repozytoria.WylaczeniaPamieci),
		Terminal:    s.terminal.ZDziennikiemWyjscia(),
		Diagnostyka: s.diagnostyka,
		// Ślad wywołań i zużycie stoją na jednym magazynie; Prowenancja bierze też kanały dla replay.
		Prowenancja: nowyAdapterProwenancji(s.repozytoria.Prowenancja).ZKanalami(s.kanaly),
		Zuzycie:     nowyAdapterZuzycia(s.repozytoria.Prowenancja),
		// Warsztat PDF pracuje biblioteką wkompilowaną i nie startuje procesu potomnego.
		WarsztatPdf: warsztatPdf,
		// Bezpieczeństwo sięga po ten sam warsztat PDF i po repozytorium Studia dla rozpoznania treści.
		BezpieczenstwoDokumentu: nowyAdapterBezpieczenstwaStudia(warsztatPdf).
			ZeStudiem(s.repozytoria.Studio),
		// Oba porty muszą być wypełnione — strażniki nil dają odpowiedź `*.unknown` zamiast realnej.
		Developer: s.developer,
		// Debata bierze katalog danych pod transkrypty i arsenał pod zamianę na dokument oraz odsłuch.
		Debata: nowyAdapterDebaty(s.zycie, s.repozytoria.Roundtable).
			ZKanalami(s.kanaly).
			ZKatalogiemDanych(katalogDanych).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz),
		// Automations wykonuje adapterem kolejek domeny kolejek; instancja dzieli budzik harmonogramu.
		Automatyki: s.automatyki,
		// Biblioteka, Studio, Przeglądarka i Design mają własne repozytoria; Biblioteka bierze katalog danych.
		Biblioteka: nowyAdapterBiblioteki(s.repozytoria.Biblioteka).
			ZKatalogiemDanych(katalogDanych).
			ZNarzedziami(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy).
			ZKanalamiModelu(s.kanaly),
		// Studio dzieli kanały z oknem rozmowy i Roundtable oraz magazyn zasobów z PDF i Design.
		Studio: nowyAdapterStudia(s.repozytoria.Studio).
			ZKanalami(s.kanaly, s.nadzorca.Rejestr()).
			ZNarzedziami(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy).
			ZZasobami(s.repozytoria.Design, katalogDanych).
			ZBiblioteka(s.repozytoria.Biblioteka),
		// Przeglądarka dostaje katalog danych i uruchamiacz z izolacją dla silnika Chromium, nic więcej.
		Przegladarka: nowyAdapterPrzegladarki(s.repozytoria.Przegladanie).
			ZMagazynem(katalogDanych).
			ZSilnikiem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		// Design bierze rejestr kanałów, katalog danych, sejf kluczy dostawców i drogę odczytu pisma.
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
		// Assistant dostaje kanały, nadzorcę sesji i nadajnik zdarzeń; `ZMowa` dzieli silnik z portem `Mowa`.
		Asystent: asystent,
		// Apps bierze rejestr okien, katalog danych z sejfem wydawcy i rejestr pozycji katalogu do publikacji.
		Aplikacje: s.aplikacje,
		// Komponenty własne dostają rejestr i trzy magazyny modułowe; profili asystenta platforma nie ma.
		Komponenty: nowyAdapterKomponentow(s.repozytoria.Komponenty).
			ZMagazynamiModulow(s.repozytoria.PrzestrzenRobocza, s.repozytoria.Agenci,
				s.repozytoria.Automatyki),
		// Translate bierze kanały, arsenał syntezy mowy z katalogiem danych i wiązania publikacji wytworów.
		Tlumaczenie: nowyAdapterTlumaczenia(s.repozytoria.Tlumaczenia).ZKanalami(s.kanaly).
			ZSynteza(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy, katalogDanych).
			ZWytworami(s.repozytoria.Biblioteka, katalogDanych),
		// Przekazanie okna bierze adapter kolejek oraz cztery repozytoria domykające lukę kontrakt–schemat.
		PrzekazanieOkna: nowyAdapterPrzekazaniaOkna(s.repozytoria.Przekazania).
			ZKolejkami(s.kolejki).
			ZBytamiOkien(s.repozytoria.Okna, s.repozytoria.Sesje,
				s.repozytoria.Moduly, s.repozytoria.Kanaly).
			ZKatalogiemAkcji(s.repozytoria.Akcje),
		// Izolacja bierze rozstrzygacz polityki, bo widoczność historii i pamięci rozstrzyga się warstwowo.
		Izolacja: nowyAdapterIzolacji(s.repozytoria, s.rozstrzygacz).ZRozgloszeniem(s.szynaZdarzen),
		// Uwierzytelnianie stoi na repozytorium budowanym przez zestaw danych
		// (dane/zestaw.go).
		Uwierzytelnianie: nowyAdapterUwierzytelnienia(s.repozytoria.Uwierzytelnienie).
			// Ten sam sejf co Konta i PunktyDostepu; bez niego sekret bramki nie ma gdzie leżeć.
			ZSejfem(sejf).
			// Tożsamość właściciela: login, adres i drogi potwierdzenia; bez niej rejestracja nie zapisze konta.
			ZKontem(s.repozytoria.KontoWlasciciela).
			// Konto nadawcze platformy, nie skrzynka Operatora — nim idą listy potwierdzenia i odzyskania.
			ZNadajnikiem(s.nadajnik).
			// Ten sam adapter ustawień co reszta platformy; zmiana konta nadawczego działa bez restartu rdzenia.
			ZNastawamiPlatformy(s.ustawienia),
		// Urządzenia stoją na repozytorium bramki: urządzeniem konta jest to, które przez nią weszło.
		Urzadzenia: nowyAdapterUrzadzen(s.repozytoria.Uwierzytelnienie),
		// Nakładka AOD dostaje nadzorcę sesji, telemetrię i rejestr obecności — te same, którymi jedzie rdzeń.
		NakladkaAod: nakladkaAod,
		// Katalog rozszerzeń wystawiony metodą; katalog dostępu i biblioteka ekspertów idą z rozgłoszeniem.
		Rozszerzenia: nowyAdapterRozszerzen(s.repozytoria.Rozszerzenia()).
			ZKatalogiemDostepu(s.repozytoria.PunktyDostepu, s.repozytoria.Agenci).
			ZMagazynemRozszerzen(katalogDanych).
			ZRozgloszeniem(s.szynaZdarzen),
		// Role okien biorą cały zestaw dla przekazań i adresów, na tym samym adapterze co port Okna.
		Role: okna.ZRolami(s.repozytoria),
		// Układ paneli wystawiony metodą, jak role i rozszerzenia — rejestr bez stanu poza pamięcią zapytań.
		Panele: nowyAdapterSekcjiPaneli(s.repozytoria.SekcjePaneli()),
		// Historia rozmowy okna bierze własne repozytorium, wiązane po identyfikatorze kontraktowym.
		Historia: NowyPortHistorii(s.repozytoria.Historia),

		// Podagenci biorą wiązania ze złożenia modułów: repozytoria, biegi, kolejki i żywotność procesów.
		Podagenci: zlozPodagentow(s),

		// Zespoły biorą własne repozytorium i dziennik rdzenia; skład bez usuniętego eksperta wraca krótszy.
		Zespoly: nowyAdapterZespolow(s.repozytoria.Zespoly).
			ZDziennikiem(s.montaz.Dziennik),
		// Historia i archiwum eksperta biorą dwa widoki wersji oraz bibliotekę dla kształtu kontraktu.
		WersjeEksperta: NowyPortWersjiEksperta(s.repozytoria.WersjeAgenta,
			s.repozytoria.ArchiwumAgentow, s.repozytoria.Agenci),
		// Zakres eksperta bierze pięć źródeł: repozytorium, bibliotekę, historię wersji, moduły i dostęp.
		ZakresEksperta: NowyPortZakresuEksperta(s.repozytoria.ZakresAgenta,
			s.repozytoria.Agenci, s.repozytoria.WersjeAgenta, s.repozytoria.Moduly,
			s.repozytoria.PunktyDostepu, s.rozstrzygacz),
		// Zakresy narzędzi stoją na własnym repozytorium i katalogu profili; profil pusty bierze domyślny.
		ZakresyNarzedzi: s.zakresyNarzedzi,
		// Silnik mowy dzieli uruchamiacz i izolację z Terminalem i Developerem; instancję bierze Assistant.
		Mowa: s.mowa,
		// Wskaźnik znaczenia dzieli uruchamiacz i izolację z Terminalem, Developerem i silnikiem mowy.
		Wiedza: nowyAdapterWiedzy(injection.UruchamiaczOkien(), katalogDanych).
			ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
			ZeSkladnica(skladnicaWiedzy(s.montaz)).
			ZeZrodlami(s.repozytoria.Biblioteka, s.repozytoria.Historia),
		// Narzędzia obrazu biorą repozytorium Designu, katalog danych i arsenał wspólny z Terminalem.
		NarzedziaObrazu: nowyAdapterNarzedziObrazu(s.repozytoria.Design, katalogDanych).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		// Silniki neuronowe obrazu dzielą zaplecze z rodziną `image.*`; druga instancja niczego nie rozdwaja.
		NarzedziaObrazuModelu: nowyAdapterNarzedziObrazuModelu(
			nowyAdapterNarzedziObrazu(s.repozytoria.Design, katalogDanych).
				ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy)),
		// Narzędzia mediów stoją na trójce co obraz: repozytorium, katalog danych, arsenał.
		NarzedziaMedia: nowyAdapterNarzedziMediow(injection.UruchamiaczOkien(),
			s.repozytoria.Design, katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy),
		// Narzędzia dokumentowe dzielą arsenał i katalog danych z rodzinami wyżej; Design do odczytu.
		Dokumenty: nowyAdapterNarzedziDokumentu(injection.UruchamiaczOkien()).
			ZKatalogiemDanych(katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
			ZZasobamiDesignu(s.repozytoria.Design).
			ZBiblioteka(s.repozytoria.Biblioteka),
		// Narzędzia archiwum dzielą trójkę z rodzinami wyżej; katalog roboczy służy `7z` i ścieżkom żądania.
		NarzedziaArchiwum: nowyAdapterNarzedziArchiwum(s.repozytoria.Design, katalogDanych).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		// Poczta bierze repozytorium skrzynek wystawione metodą, sejf, katalog danych i repozytorium Designu.
		Poczta: nowyAdapterPoczty(s.repozytoria.SkrzynkiOperatora()).
			ZSejfem(sejf).ZKatalogiemDanych(katalogDanych).
			ZZasobami(s.repozytoria.Design),
		// Dołożenia sesji: adapter z montażu dzielony z zestawem narzędzi tury; konstruktor bierze trzy byty.
		NarzedziaSesji: s.dolozenia,
		// Doradca powstaje w złożeniu modułów z mową: stoi na rejestrze kanałów okna rozmowy i dzienniku bazy.
		Doradcy:  s.doradcy,
		Nadajnik: s.szynaZdarzen,
		Nasluch:  s.nasluch,
		Dziennik: s.montaz.Dziennik,
	}
}
