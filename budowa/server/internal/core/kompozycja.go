package core

import (
	"log"

	"danacoconsole/shared"
)

// Porty jest kompletem zależności rdzenia: każde pole jest interfejsem z
// porty.go, a rdzeń nie wie, kto je wypełnia. Pole niewypełnione nie
// zatrzymuje startu — domena bez portu nie rejestruje obsługiwaczy, a jej
// komendy odpowiadają `*.unknown`.
type Porty struct {
	Sesje         Sesje
	Okna          Okna
	Rozmowa       Rozmowa
	Ustawienia    Ustawienia
	Kanaly        Kanaly
	Kolejki       Kolejki
	Przenoszenie  Przenoszenie
	Nawigacja     Nawigacja
	WiazanieSesji WiazanieSesji
	Akcje         Akcje
	// CentrumPowiadomien jest rejestrem zdarzeń centrum powiadomień, zgłaszanym
	// przez sam rdzeń.
	CentrumPowiadomien CentrumPowiadomien
	// Pięć portów okna konfiguracji: katalog pozycji, dostęp i nadania rozmowy,
	// konta, tożsamość modelu.
	KatalogUstawien KatalogUstawien
	PunktyDostepu   PunktyDostepu
	NadaniaDostepu  NadaniaDostepu
	Konta           Konta
	Tozsamosc       Tozsamosc
	// Agenci obsługuje moduł Agents: bibliotekę ekspertów, model bazowy,
	// umiejętności, uprawnienia.
	Agenci Agenci
	// PrzestrzenRobocza obsługuje moduł Workspace: projekt, instrukcje, pamięć,
	// bibliotekę, ekspertów.
	PrzestrzenRobocza PrzestrzenRobocza
	// Terminal obsługuje moduł Terminal: karty powłok, rejestr procesów rdzenia
	// i strumień wyjścia.
	Terminal Terminal
	// Developer obsługuje moduł Developer: plik repozytorium, drzewo projektu
	// i przebiegi budowania.
	Developer Developer
	// Automatyki obsługuje moduł Automations: definicje, harmonogramy,
	// zależności, przebiegi.
	Automatyki Automatyki
	// Biblioteka obsługuje moduł Library: pliki repozytorium wiedzy, ich wersje,
	// etykiety i kolekcje.
	Biblioteka Biblioteka
	// Studio obsługuje moduł Studio: dokumenty edytora, ich wersje i propozycje
	// zmian.
	Studio Studio
	// Przegladarka obsługuje moduł Browser: migawki stron, zebrane źródła
	// i notatki.
	Przegladarka Przegladarka
	// Design obsługuje moduł Design: prompty strukturalne, zasoby wizualne
	// i kompozycje.
	Design Design
	// Badania obsługuje moduł Research: źródła, ustalenia, raport i przestrzeń.
	Badania Badania
	// Asystent obsługuje moduł Assistant: polecenia, stan zleceń i dziennik
	// czynności.
	Asystent Asystent
	// Aplikacje obsługuje moduł Apps: architektura, warsztat i wdrożenia.
	Aplikacje Aplikacje
	// Komponenty obsługuje rodzinę `component.*` — kafle Strefy 2 Strony
	// głównej.
	Komponenty Komponenty
	// PrzekazanieOkna obsługuje obszar window.*: przekazanie zlecenia między
	// oknami i akcje panelu okna.
	PrzekazanieOkna PrzekazanieOkna
	// Tlumaczenie obsługuje moduł Translate: szesnaście komend przekładu,
	// słownika, jakości i mowy.
	Tlumaczenie Tlumaczenie
	// Diagnostyka obsługuje moduł Diagnostics: dziennik rdzenia, wykaz błędów,
	// analizy, rekomendacje.
	Diagnostyka Diagnostyka
	// Prowenancja i Zuzycie są portami przekrojowymi, nie modułowymi — jeden
	// magazyn, dwa porty.
	Prowenancja Prowenancja
	Zuzycie     Zuzycie
	// Kondycja i Alerty opisują stan produktu: sondy z pomiarami oraz reguły
	// wyzwalania z rejestrem.
	Kondycja Kondycja
	Alerty   Alerty
	// Schowek, SkrotyTekstowe i KontekstyPamieci obsługują tekst poza jednym
	// oknem platformy.
	Schowek          Schowek
	SkrotyTekstowe   SkrotyTekstowe
	KontekstyPamieci KontekstyPamieci
	// Wywolywacz obsługuje rodzinę `launcher.*`: skrót globalny otwierający
	// wywoływacz poleceń.
	Wywolywacz Wywolywacz
	// ZajetoscKontekstu obsługuje `context.usage.get`: zajętość okna kontekstu
	// tokenizatorem rdzenia.
	ZajetoscKontekstu ZajetoscKontekstu
	// zdolnosciKlienta jest nieeksportowaną mapą deklaracji klientów spisanych
	// w powitaniu.
	zdolnosciKlienta *zdolnosciKlientow
	// WarsztatPdf jest częścią modułu Studio o innym komplecie zależności:
	// biblioteka w rdzeniu.
	WarsztatPdf WarsztatPdf
	// BezpieczenstwoDokumentu stoi obok warsztatu na tych samych zależnościach,
	// ale w osobnym porcie.
	BezpieczenstwoDokumentu BezpieczenstwoDokumentu
	// Debata obsługuje moduł Roundtable: skład debaty, tury, wypowiedzi
	// i stanowisko końcowe.
	Debata Debata
	// Izolacja obsługuje rodzinę `isolation.*`: profile, warstwy, zakresy
	// techniczne, podgląd polityki.
	Izolacja Izolacja
	// Uwierzytelnianie obsługuje rodzinę `auth.*`: bramka, metody i sesja
	// Operatora (`dane/auth.go`).
	Uwierzytelnianie Uwierzytelnianie
	// Urzadzenia obsługuje rodzinę `device.*`: wykaz maszyn konta i unieważnienie
	// ich tokenów.
	Urzadzenia Urzadzenia
	// NakladkaAod obsługuje rodzinę `aod.*`: nakładkę zawsze-na-wierzchu nad
	// stanem biegnącej pracy.
	NakladkaAod NakladkaAod
	// Rozszerzenia obsługuje rodzinę `extension.*`: katalog rozszerzeń.
	Rozszerzenia Rozszerzenia
	// Role obsługuje rodzinę `role.*`: rolę okna, więź koordynatora, wykaz nadań
	// i zdjęcie roli.
	Role RoleWykaz
	// Panele obsługuje `panel.sections.*`: kolejność, zwinięcie i zdjęcie
	// sekcji panelu okna.
	Panele SekcjePaneli
	// Historia obsługuje `history.*` i `retention.set`: wykaz pozycji rozmowy
	// i zasadę przechowywania.
	Historia Historia

	// Podagenci obsługuje rodzinę `subagent.*`: powołanie podagentów okna
	// wykonawcy i zbieranie wyników.
	Podagenci Podagenci

	// WersjeEksperta obsługuje historię tożsamości i archiwum eksperta,
	// w osobnym repozytorium.
	WersjeEksperta WersjeEksperta

	// ZakresEksperta obsługuje zakres działania eksperta: Permissions Center
	// i dopełnienia okien.
	ZakresEksperta ZakresEksperta

	// ZakresyNarzedzi obsługuje rodzinę `tools.scope.*`: zakresy uprawnień
	// i limity wywołań katalogu.
	ZakresyNarzedzi ZakresyNarzedzi

	// Zespoly obsługuje rodzinę `team.*`: zapisane składy ekspertów, byt
	// Operatora nad biblioteką.
	Zespoly Zespoly
	// Mowa obsługuje rodzinę `speech.*`: gotowość silnika rozpoznawania
	// i zamianę nagrania na tekst.
	Mowa Mowa
	// Wiedza obsługuje rodzinę `knowledge.*`: wskaźnik znaczenia i wyszukiwanie
	// po sensie, nie literach.
	Wiedza Wiedza
	// NarzedziaObrazu obsługuje rodzinę `image.*`: cztery narzędzia modelu do
	// pracy na obrazie.
	NarzedziaObrazu NarzedziaObrazu
	// NarzedziaObrazuModelu obsługuje powiększanie i usuwanie tła siecią
	// neuronową, osobno od retuszu.
	NarzedziaObrazuModelu NarzedziaObrazuModelu
	// Doradcy obsługuje `advisor.consult`: konsultację modelu u doradcy
	// w trakcie tury.
	Doradcy Doradcy
	// NarzedziaMedia obsługuje rodzinę `media.*`: pomiar i przetworzenie dźwięku
	// oraz filmu.
	NarzedziaMedia NarzedziaMedia
	// Dokumenty obsługuje rodzinę `document.*`: zamianę formatu i odczyt treści
	// dokumentu.
	Dokumenty Dokumenty
	// NarzedziaArchiwum obsługuje rodzinę `archive.*`: złożenie i otwarcie
	// archiwum plików.
	NarzedziaArchiwum NarzedziaArchiwum
	// Poczta obsługuje rodzinę `mail.*`: dziesięć komend klienta skrzynki
	// Operatora.
	Poczta Poczta
	// NarzedziaSesji obsługuje `session.tool.*` i `tools.catalog.list`:
	// narzędzia doraźne sesji.
	NarzedziaSesji NarzedziaSesji
	Nadajnik       Nadajnik
	Nasluch        Nasluch
	Dziennik       *log.Logger
}

// Zloz składa rdzeń: zakłada rejestr, wpina obsługiwacze wszystkich domen
// kontraktu i buduje rozpoznanie nazw z tego, co rzeczywiście zostało wpięte —
// jedyne miejsce, w którym powstaje mapa komend rdzenia.
func Zloz(p Porty) *Rdzen {
	rejestr := NowyRejestr()
	nadawca := nowyEmiter(p.Nadajnik)
	// Więź połączenie-sesja bramki, jedna na rdzeń: zakłada ją powitanie,
	// nadpisuje bramka.
	wiez := nowaWiezBramki()
	// Deklaracje zdolności klientów spisane w powitaniu, jedna mapa na rdzeń,
	// czyta ją wywoływacz.
	zdolnosci := p.zdolnosciKlienta
	if zdolnosci == nil {
		zdolnosci = noweZdolnosciKlientow()
	}

	// Nastawy poziomu `aplikacja` czytane tym samym adapterem co cała
	// konfiguracja.
	nastawy, _ := p.Ustawienia.(NastawyAplikacji)
	zarejestrujPolaczenie(rejestr, p.Uwierzytelnianie, wiez, nastawy, zdolnosci)
	zarejestrujSesje(rejestr, p.Sesje, nadawca)
	zarejestrujOkna(rejestr, p.Okna, nadawca)
	zarejestrujKanalModelu(rejestr, p.Okna.(KanalModelu), nadawca)
	zarejestrujRozmowe(rejestr, p.Rozmowa, nadawca)
	zarejestrujUstawienia(rejestr, p.Ustawienia, nadawca)
	zarejestrujProwenancjeKonfiguracji(rejestr, p.Ustawienia.(ProwenancjaKonfiguracji))
	zarejestrujKanaly(rejestr, p.Kanaly)
	zarejestrujKolejki(rejestr, p.Kolejki, nadawca)
	zarejestrujWiazaniaKolejek(rejestr, p.Kolejki, nadawca)
	zarejestrujZleceniaKolejek(rejestr, p.Kolejki, nadawca)
	zarejestrujPrzenoszenie(rejestr, p.Przenoszenie, nadawca)
	zarejestrujNawigacje(rejestr, p.Nawigacja, nadawca)
	zarejestrujMonitor(rejestr, p.Nawigacja)
	zarejestrujWiazanieSesji(rejestr, p.WiazanieSesji, nadawca)
	zarejestrujAkcje(rejestr, p.Akcje, shared.CommandActionList)
	zarejestrujCentrumPowiadomien(rejestr, p.CentrumPowiadomien)
	zarejestrujKatalogUstawien(rejestr, p.KatalogUstawien)
	zarejestrujPunktyDostepu(rejestr, p.PunktyDostepu, nadawca)
	zarejestrujNadaniaDostepu(rejestr, p.NadaniaDostepu, nadawca)
	zarejestrujKonta(rejestr, p.Konta, nadawca)
	zarejestrujTozsamosc(rejestr, p.Tozsamosc, nadawca)
	zarejestrujAgentow(rejestr, p.Agenci, nadawca)
	zarejestrujWersjeEksperta(rejestr, p.WersjeEksperta, nadawca)
	zarejestrujZakresEksperta(rejestr, p.ZakresEksperta, nadawca)
	zarejestrujZakresyNarzedzi(rejestr, p.ZakresyNarzedzi)
	zarejestrujZespoly(rejestr, p.Zespoly, nadawca)
	zarejestrujPrzestrzenRobocza(rejestr, p.PrzestrzenRobocza, nadawca)
	// Moduł Workspace wchodzi trzema portami ponad `PrzestrzenRobocza`: pamięć,
	// planowanie, wiedza.
	if pamiec, jest := p.PrzestrzenRobocza.(PamiecPrzestrzeni); jest {
		zarejestrujPamiec(rejestr, pamiec, nadawca)
	}
	if planowanie, jest := p.PrzestrzenRobocza.(PlanowanieProjektu); jest {
		zarejestrujPlanowanieProjektu(rejestr, planowanie, nadawca)
	}
	if wiedza, jest := p.PrzestrzenRobocza.(WiedzaProjektu); jest {
		zarejestrujWiedzeProjektu(rejestr, wiedza, nadawca)
	}
	zarejestrujAutomatyki(rejestr, p.Automatyki, nadawca)
	zarejestrujDobudoweAutomatyk(rejestr, p.Automatyki.(AutomatykiDobudowa), nadawca)
	zarejestrujHarmonogramy(rejestr, p.Automatyki.(HarmonogramyNadzoru), nadawca)
	zarejestrujOrkiestracje(rejestr, p.Automatyki.(Orkiestracja), nadawca)
	zarejestrujTerminal(rejestr, p.Terminal, nadawca)
	zarejestrujWyjscieTerminala(rejestr, p.Terminal.(WyjscieTerminala))
	// Wyposażenie modułu wchodzi rzutowaniem dwuwartościowym, jak porty
	// Workspace'u wyżej.
	if wyposazenie, jest := p.Terminal.(WyposazenieTerminala); jest {
		zarejestrujWyposazenieTerminala(rejestr, wyposazenie)
	}
	zarejestrujDevelopera(rejestr, p.Developer, nadawca)
	zarejestrujDiagnostyke(rejestr, p.Diagnostyka, nadawca)
	zarejestrujProwenancje(rejestr, p.Prowenancja)
	zarejestrujKondycje(rejestr, p.Kondycja)
	zarejestrujAlerty(rejestr, p.Alerty)
	zarejestrujSchowek(rejestr, p.Schowek)
	zarejestrujSkrotyTekstowe(rejestr, p.SkrotyTekstowe)
	zarejestrujKontekstyPamieci(rejestr, p.KontekstyPamieci)
	zarejestrujWywolywacz(rejestr, p.Wywolywacz)
	zarejestrujZajetoscKontekstu(rejestr, p.ZajetoscKontekstu)
	zarejestrujZuzycie(rejestr, p.Zuzycie)
	zarejestrujWarsztatPdf(rejestr, p.WarsztatPdf)
	zarejestrujBezpieczenstwoDokumentu(rejestr, p.BezpieczenstwoDokumentu)
	zarejestrujBiblioteke(rejestr, p.Biblioteka, nadawca)
	zarejestrujStudio(rejestr, p.Studio, nadawca)
	// Pętla wykonawcza Studia stoi na tym samym porcie, ale ma stan własny,
	// więc wchodzi osobno.
	zarejestrujPetleWykonawczaStudia(rejestr, p.Studio, nadawca)
	// Zapora blokad fragmentu owija to, co w rejestrze już stoi — jedzie po
	// obu wpięciach Studia.
	if adapter, jest := p.Studio.(*adapterStudia); jest {
		zaporaBlokadStudia(rejestr, adapter)
		// Siatka śladu autora idzie po zaporze i owija ją z zewnątrz — kolejność
		// jest zamierzona.
		sladWykonawcyStudia(rejestr, adapter)
	}
	zarejestrujPrzegladarke(rejestr, p.Przegladarka, nadawca)
	zarejestrujDesign(rejestr, p.Design, nadawca)
	zarejestrujBadania(rejestr, p.Badania, nadawca)
	zarejestrujAsystenta(rejestr, p.Asystent, nadawca)
	zarejestrujAplikacje(rejestr, p.Aplikacje, nadawca)
	zarejestrujKomponenty(rejestr, p.Komponenty, nadawca)
	zarejestrujPrzekazanieOkna(rejestr, p.PrzekazanieOkna, nadawca)
	zarejestrujTlumaczenie(rejestr, p.Tlumaczenie, nadawca)
	zarejestrujDebate(rejestr, p.Debata, nadawca)
	zarejestrujIzolacje(rejestr, p.Izolacja)
	zarejestrujUwierzytelnianie(rejestr, p.Uwierzytelnianie, nadawca, wiez)
	zarejestrujUrzadzenia(rejestr, p.Urzadzenia, nadawca, wiez)
	zarejestrujNakladkeAod(rejestr, p.NakladkaAod, nadawca)
	// Dwa wpięcia bez nowego portu: obie funkcje stoją na portach, które rdzeń
	// już niesie.
	zarejestrujHistorieSesji(rejestr, p.Sesje, nadawca)
	zarejestrujWarstweMobilna(rejestr, p.Nawigacja)
	zarejestrujRozszerzenia(rejestr, p.Rozszerzenia)
	zarejestrujRole(rejestr, p.Role, nadawca)
	zarejestrujWykazRol(rejestr, p.Role, nadawca)
	zarejestrujSekcjePaneli(rejestr, p.Panele)
	zarejestrujHistorie(rejestr, p.Historia, nadawca)
	zarejestrujMowe(rejestr, p.Mowa)
	zarejestrujWiedze(rejestr, p.Wiedza)
	zarejestrujOrkiestracjePodagentow(rejestr, p.Podagenci)
	zarejestrujDoradcow(rejestr, p.Doradcy, nadawca)
	zarejestrujDokumenty(rejestr, p.Dokumenty)
	zarejestrujNarzedziaMediow(rejestr, p.NarzedziaMedia)
	zarejestrujNarzedziaObrazu(rejestr, p.NarzedziaObrazu)
	zarejestrujNarzedziaObrazuModelu(rejestr, p.NarzedziaObrazuModelu)
	zarejestrujNarzedziaArchiwum(rejestr, p.NarzedziaArchiwum)
	zarejestrujPoczte(rejestr, p.Poczta, nadawca)
	zarejestrujNarzedziaSesji(rejestr, p.NarzedziaSesji, nadawca)

	return &Rdzen{
		rejestr:  rejestr,
		komendy:  rejestr.RejestrKomend(),
		nasluch:  p.Nasluch,
		dziennik: p.Dziennik,
		wiez:     wiez,
	}
}
