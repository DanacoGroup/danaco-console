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

// Struktura zamiast dwudziestu pozycyjnych argumentów podatnych na cichą pomyłkę kolejności.
type skladPortow struct {
	zycie           context.Context
	montaz          Montaz
	repozytoria     *dane.Zestaw
	nadzorca        *session.Nadzorca
	kanaly          *models.Rejestr
	rozstrzygacz    *konfig.Rozstrzygacz
	katalogRoboczy  *KatalogRoboczy
	trwalosc        *utrwalaczStanow
	obecnosc        *rejestrObecnosci
	biegi           *rejestrBiegow
	telemetria      *telemetriaPostepu
	rozmowa         *adapterRozmowy
	petla           *session.Petla
	ustawienia      *adapterUstawienOsi
	tozsamosc       Tozsamosc
	klienci         *wieziKlientow
	utrwalacz       *dane.UtrwalaczRozmowy
	dziennikRozmow  *dziennikRozmowy
	terminal        *adapterTerminala
	developer       *adapterDevelopera
	aplikacje       *adapterAplikacji
	diagnostyka     *adapterDiagnostyki
	kolejki         *adapterKolejek
	automatyki      *adapterAutomatyk
	mowa            *adapterMowy
	doradcy         *adapterDoradcow
	szynaZdarzen    Nadajnik
	nasluch         nasluchTransportu
	dolozenia       *adapterNarzedziSesji
	zakresyNarzedzi *adapterZakresowNarzedzi
	strazEkspertow  StrazEksperta
	nadajnik        nadajnik.Nastawy
}

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

// Zmontuj składa byty i wiąże je ze sobą; ta funkcja wyłącznie przekłada je na porty kontraktu.
func zlozPorty(s skladPortow) Porty {
	katalogDanych := s.montaz.Konfiguracja.KatalogDanych
	// Jeden sejf poświadczeń nad katalogiem danych: druga instancja ścigałaby się o zapis.
	sejf := dane.NowySejfPlikowy(katalogDanych)
	if s.montaz.Dziennik != nil {
		s.montaz.Dziennik.Printf("sejf poświadczeń: %s", sejf.OpisKlucza())
	}
	models.UstawSejfPoswiadczen(sejf)
	sesje := nowyAdapterSesji(s.nadzorca).ZTrwaloscia(s.trwalosc).
		ZObecnoscia(s.obecnosc).ZZapewnieniem(s.utrwalacz).
		ZProjektami(s.repozytoria.PrzestrzenRobocza).ZZestawem(s.repozytoria).
		ZZatrzymaniemTur(s.rozmowa.PrzerwijTure).ZRozmowa(s.rozmowa.dziennik)
	akcje := rejestrAkcji(s.zycie, s.repozytoria, s.montaz.Dziennik)

	centrumPowiadomien := nowyAdapterCentrumPowiadomien(s.repozytoria.CentrumPowiadomien).
		ZNadawca(nowyEmiter(s.szynaZdarzen)).
		ZNastawami(s.ustawienia)
	przenoszenie := nowyAdapterPrzenoszenia(s.nadzorca).
		ZTrwaloscia(s.zycie, s.repozytoria.Konfiguracja, s.montaz.Dziennik).
		ZHistoria(s.dziennikRozmow)

	asystent := nowyAdapterAsystenta(s.repozytoria.Asystent).
		ZKanalami(s.kanaly).ZSesjami(s.nadzorca).ZWyjsciem(s.szynaZdarzen, s.zycie).
		ZMowa(s.mowa).ZMostami(s.rozmowa.mosty).ZDomknieciemZastanych()

	nakladkaAod := nowyAdapterNakladkiAod(s.nadzorca, s.telemetria, s.obecnosc).
		ZRozmowa(s.telemetria.OwinRozmowe(s.rozmowa)).
		ZAsystentem(asystent).
		ZPodpowiedziami(akcje).
		ZKontekstem(przenoszenie).
		ZUrzadzeniami(s.repozytoria.Urzadzenia).
		ZWyciszeniami(s.repozytoria.WyciszeniaNakladki)

	okna := nowyAdapterOkien(s.nadzorca).ZTrwaloscia(s.trwalosc).ZPetla(s.petla).
		ZPrzerwaniemTury(s.rozmowa.PrzerwijTure, s.rozmowa.CzyTuraWBiegu).
		ZWieziami(s.repozytoria.Przekazania).ZKanalami(s.repozytoria.Kanaly).
		ZModulami(s.repozytoria.Moduly).ZeStrazaEksperta(s.strazEkspertow)

	warsztatPdf := nowyAdapterPdfStudia().
		ZKatalogiemDanych(katalogDanych).ZZasobami(s.repozytoria.Design)

	zdolnosciKlientow := noweZdolnosciKlientow()

	return Porty{
		zdolnosciKlienta: zdolnosciKlientow,
		Kondycja: nowyAdapterKondycji(s.repozytoria.Kondycja).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy).
			ZKanalami(s.kanaly, s.repozytoria.Kanaly),
		Alerty: nowyAdapterAlertow(s.repozytoria.Alerty).
			ZeZrodlamiMiar(s.repozytoria.Prowenancja, s.repozytoria.Diagnostyka, s.rozstrzygacz).
			ZWyjsciem(nowyEmiter(s.szynaZdarzen)).
			ZCentrumPowiadomien(centrumPowiadomien),
		Schowek:          nowyAdapterSchowka(s.repozytoria.Schowek),
		SkrotyTekstowe:   nowyAdapterSkrotowTekstowych(s.repozytoria.SkrotyTekstowe),
		KontekstyPamieci: nowyAdapterKontekstowPamieci(s.repozytoria.KontekstyPamieci),
		Wywolywacz: nowyAdapterWywolywacza(s.repozytoria.Konfiguracja, s.rozstrzygacz,
			zdolnosciKlientow),
		ZajetoscKontekstu: nowyAdapterZajetosciKontekstu(s.repozytoria.Okna,
			s.repozytoria.Wiadomosci, s.repozytoria.Pamiec, s.kanaly, s.tozsamosc),
		Sesje:        sesje,
		Okna:         okna,
		Rozmowa:      s.telemetria.OwinRozmowe(s.rozmowa),
		Ustawienia:   s.ustawienia.ZProwenancja(s.rozmowa),
		Kanaly:       nowyAdapterKanalow(s.repozytoria.Kanaly, s.kanaly).ZSejfem(sejf),
		Kolejki:      s.kolejki,
		Przenoszenie: przenoszenie,
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
		Agenci: nowyAdapterAgentow(s.repozytoria.Agenci).
			ZKanalami(s.kanaly).
			ZPunktamiDostepu(s.repozytoria.PunktyDostepu).
			ZWarstwami(s.repozytoria.WarstwyAgenta),
		PrzestrzenRobocza: nowyAdapterPrzestrzeniRoboczej(s.repozytoria.PrzestrzenRobocza).
			ZSesjami(sesje).
			ZInstrukcjami(s.repozytoria.Konfiguracja, s.rozstrzygacz).
			ZKatalogiem(s.katalogRoboczy).
			// `ZPamiecia` musi stać ostatnim ogniwem.
			ZDokumentamiWorkspace(nowyAdapterNarzedziDokumentu(injection.UruchamiaczOkien()).
				ZKatalogiemDanych(katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
				ZZasobamiDesignu(s.repozytoria.Design).
				ZBiblioteka(s.repozytoria.Biblioteka)).
			// `ZWylaczeniami` idzie po `ZPamiecia` — inna kolejność zostawia `memory.disable.*` bez magazynu.
			ZPamiecia(s.repozytoria.Pamiec, s.repozytoria.Sesje).
			ZWylaczeniami(s.repozytoria.WylaczeniaPamieci),
		Terminal:    s.terminal.ZDziennikiemWyjscia(),
		Diagnostyka: s.diagnostyka,
		Prowenancja: nowyAdapterProwenancji(s.repozytoria.Prowenancja).
			ZKanalami(s.kanaly, s.repozytoria.Kanaly),
		Zuzycie:     nowyAdapterZuzycia(s.repozytoria.Prowenancja),
		WarsztatPdf: warsztatPdf,
		BezpieczenstwoDokumentu: nowyAdapterBezpieczenstwaStudia(warsztatPdf).
			ZeStudiem(s.repozytoria.Studio),
		// Port pozostawiony na nil daje odpowiedź `*.unknown` zamiast realnej.
		Developer: s.developer,
		Debata: nowyAdapterDebaty(s.zycie, s.repozytoria.Roundtable).
			ZKanalami(s.kanaly, s.repozytoria.Kanaly).
			ZKatalogiemDanych(katalogDanych).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz),
		Automatyki: s.automatyki,
		Biblioteka: nowyAdapterBiblioteki(s.repozytoria.Biblioteka).
			ZKatalogiemDanych(katalogDanych).
			ZNarzedziami(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy).
			ZKanalamiModelu(s.kanaly),
		Studio: nowyAdapterStudia(s.repozytoria.Studio).
			ZKanalami(s.kanaly, s.nadzorca.Rejestr()).
			ZNarzedziami(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy).
			ZZasobami(s.repozytoria.Design, katalogDanych).
			ZBiblioteka(s.repozytoria.Biblioteka),
		Przegladarka: nowyAdapterPrzegladarki(s.repozytoria.Przegladanie).
			ZMagazynem(katalogDanych).
			ZSilnikiem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		Design: nowyAdapterDesignu(s.repozytoria.Design).ZKanalami(s.kanaly, s.repozytoria.Kanaly).
			ZKatalogiemDanych(katalogDanych).ZSejfem(sejf).
			ZOdczytemPisma(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		Badania: nowyAdapterBadan(s.repozytoria.Badania).ZKanalami(s.kanaly).
			ZKatalogiemDanych(katalogDanych).
			ZDokumentami(nowyAdapterNarzedziDokumentu(injection.UruchamiaczOkien()).
				ZKatalogiemDanych(katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
				ZZasobamiDesignu(s.repozytoria.Design).
				ZBiblioteka(s.repozytoria.Biblioteka)).
			ZMowa(s.mowa),
		Asystent:  asystent,
		Aplikacje: s.aplikacje,
		Komponenty: nowyAdapterKomponentow(s.repozytoria.Komponenty).
			ZMagazynamiModulow(s.repozytoria.PrzestrzenRobocza, s.repozytoria.Agenci,
				s.repozytoria.Automatyki),
		Tlumaczenie: nowyAdapterTlumaczenia(s.repozytoria.Tlumaczenia).
			ZKanalami(s.kanaly, s.repozytoria.Kanaly).
			ZSynteza(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy, katalogDanych).
			ZWytworami(s.repozytoria.Biblioteka, katalogDanych),
		PrzekazanieOkna: nowyAdapterPrzekazaniaOkna(s.repozytoria.Przekazania).
			ZKolejkami(s.kolejki).
			ZBytamiOkien(s.repozytoria.Okna, s.repozytoria.Sesje,
				s.repozytoria.Moduly, s.repozytoria.Kanaly).
			ZKatalogiemAkcji(s.repozytoria.Akcje),
		Izolacja: nowyAdapterIzolacji(s.repozytoria, s.rozstrzygacz).ZRozgloszeniem(s.szynaZdarzen),
		Uwierzytelnianie: nowyAdapterUwierzytelnienia(s.repozytoria.Uwierzytelnienie).
			ZSejfem(sejf).
			ZKontem(s.repozytoria.KontoWlasciciela).
			ZNadajnikiem(s.nadajnik).
			ZAdresemKonsoli(s.montaz.Konfiguracja.AdresKonsoli).
			ZNastawamiPlatformy(s.ustawienia),
		Urzadzenia:  nowyAdapterUrzadzen(s.repozytoria.Uwierzytelnienie),
		NakladkaAod: nakladkaAod,
		Rozszerzenia: nowyAdapterRozszerzen(s.repozytoria.Rozszerzenia()).
			ZKatalogiemDostepu(s.repozytoria.PunktyDostepu, s.repozytoria.Agenci).
			ZMagazynemRozszerzen(katalogDanych).
			ZRozgloszeniem(s.szynaZdarzen),
		Role:     okna.ZRolami(s.repozytoria),
		Panele:   nowyAdapterSekcjiPaneli(s.repozytoria.SekcjePaneli()),
		Historia: NowyPortHistorii(s.repozytoria.Historia),

		Podagenci: zlozPodagentow(s),

		Zespoly: nowyAdapterZespolow(s.repozytoria.Zespoly).
			ZDziennikiem(s.montaz.Dziennik),
		WersjeEksperta: NowyPortWersjiEksperta(s.repozytoria.WersjeAgenta,
			s.repozytoria.ArchiwumAgentow, s.repozytoria.Agenci),
		ZakresEksperta: NowyPortZakresuEksperta(s.repozytoria.ZakresAgenta,
			s.repozytoria.Agenci, s.repozytoria.WersjeAgenta, s.repozytoria.Moduly,
			s.repozytoria.PunktyDostepu, s.rozstrzygacz),
		ZakresyNarzedzi: s.zakresyNarzedzi,
		Mowa:            s.mowa,
		Wiedza: nowyAdapterWiedzy(injection.UruchamiaczOkien(), katalogDanych).
			ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
			ZeSkladnica(skladnicaWiedzy(s.montaz)).
			ZeZrodlami(s.repozytoria.Biblioteka, s.repozytoria.Historia),
		NarzedziaObrazu: nowyAdapterNarzedziObrazu(s.repozytoria.Design, katalogDanych).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		NarzedziaObrazuModelu: nowyAdapterNarzedziObrazuModelu(
			nowyAdapterNarzedziObrazu(s.repozytoria.Design, katalogDanych).
				ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy)),
		NarzedziaMedia: nowyAdapterNarzedziMediow(injection.UruchamiaczOkien(),
			s.repozytoria.Design, katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy),
		Dokumenty: nowyAdapterNarzedziDokumentu(injection.UruchamiaczOkien()).
			ZKatalogiemDanych(katalogDanych).ZIzolacja(s.rozstrzygacz, s.katalogRoboczy).
			ZZasobamiDesignu(s.repozytoria.Design).
			ZBiblioteka(s.repozytoria.Biblioteka),
		NarzedziaArchiwum: nowyAdapterNarzedziArchiwum(s.repozytoria.Design, katalogDanych).
			ZArsenalem(injection.UruchamiaczOkien(), s.rozstrzygacz, s.katalogRoboczy),
		Poczta: nowyAdapterPoczty(s.repozytoria.SkrzynkiOperatora()).
			ZSejfem(sejf).ZKatalogiemDanych(katalogDanych).
			ZZasobami(s.repozytoria.Design),
		NarzedziaSesji: s.dolozenia,
		Doradcy:        s.doradcy,
		Nadajnik:       s.szynaZdarzen,
		Nasluch:        s.nasluch,
		Dziennik:       s.montaz.Dziennik,
	}
}
