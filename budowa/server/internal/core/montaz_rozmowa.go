package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
)

// skladRozmowy jest kompletem, z którego powstaje domena rozmowy: pętla
// koordynator–wykonawca wraz z adapterem i wszystkimi jego portami.
//
// Struktura zamiast dziesięciu argumentów, bo składanie rozmowy wiąże ze sobą
// najwięcej bytów w całym montażu, a wywołanie pozycyjne o dziesięciu polach
// jest nieczytelne i łatwe do pomylenia przy zmianie kolejności.
type skladRozmowy struct {
	zycie        context.Context
	montaz       Montaz
	repozytoria  *dane.Zestaw
	nadzorca     *session.Nadzorca
	kanaly       *models.Rejestr
	nadajnik     Nadajnik
	dziennik     *dziennikRozmowy
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	tozsamosc    Tozsamosc
	biegi        *rejestrBiegow
	obecnosc     *rejestrObecnosci
	// ustawienia odczytuje konfigurację obowiązującą sesji (KonfiguracjaSesjiOkna).
	// Ten sam adapter osi wypełnia port Ustawienia rdzenia — jeden czytelnik
	// konfiguracji na całą platformę.
	ustawienia *adapterUstawienOsi
	// zdarzenia utrwala zamknięcia tur i zdarzenia zaczepów. Ten sam
	// odbiornik dostaje fabryka kanału głównego — jedna droga śladu.
	zdarzenia *zdarzeniaWykonawcze
	// dolozenia wnosi doraźne dołożenia sesji do zestawu narzędzi tury. Ten sam
	// adapter wypełnia port `NarzedziaSesji` rdzenia (montaz_porty.go): dołożenie
	// widziane przez model i dołożenie widziane przez Operatora to jeden byt.
	dolozenia *adapterNarzedziSesji
}

// zlozRozmowe składa adapter rozmowy wraz z pętlą koordynator–wykonawca i wpina
// w niego komplet portów.
//
// Funkcja ma jedną odpowiedzialność: wszystko, co adapter rozmowy musi dostać,
// zanim zacznie prowadzić tury.
//
// Wiązanie z pętlą jest dwuetapowe, bo adapter i pętla znają się nawzajem:
// pętla potrzebuje portu rozpoczynania obiegu, którym jest ten adapter.
func zlozRozmowe(s skladRozmowy) (*adapterRozmowy, *session.Petla) {
	rozmowa := nowyAdapterRozmowy(s.zycie, s.nadzorca, s.kanaly, s.nadajnik, s.dziennik).
		ZCiagloscia(ciagloscRozmowy(s.repozytoria)).
		ZZdarzeniami(s.zdarzenia).
		// Rejestrator bloków dostaje kontekst życia rdzenia — zapis
		// fragmentu ma dobiec końca także po rozłączeniu klienta.
		ZBlokami(nowyRejestratorBlokow(s.zycie, s.repozytoria.Bloki, s.montaz.Dziennik)).
		// Magazyn załączników rozmowy stoi na tym samym katalogu danych, co sejf
		// poświadczeń i magazyn biblioteki (`montaz_porty.go`) — katalog
		// obowiązujący zna wyłącznie montaż.
		ZZalacznikami(s.montaz.Konfiguracja.KatalogDanych)

	petla := session.NowaPetla(s.nadzorca,
		session.UruchomienieFunkcja(rozmowa.RozpocznijObieg), session.UstawieniaPetli{})
	rozmowa.UstawPetle(petla)
	petla.Obserwuj(session.ObserwatorFunkcja(sladObiegu(s.montaz.Dziennik)))
	// Bieg naprawczy idzie dwiema drogami: do dziennika rdzenia i do kontrolki
	// sesji trwającej w tle. Rejestr biegów jest tą drugą.
	petla.Obserwuj(s.biegi)
	s.obecnosc.ZeStrumieniem(rozmowa.CzyTuraWBiegu)

	// Tożsamość dojeżdża do procesu modelu tą samą drogą, którą pokazuje ją okno
	// konfiguracji: jeden składacz, jeden przekład trybu.
	rozmowa.ZTozsamoscia(s.tozsamosc).
		// Ekspert wskazany w oknie wnosi warstwy promptu, model bazowy i nastawy
		// procesu. Wchodzi do tego samego składacza co osie — drugiej drogi do
		// promptu nie ma.
		ZAgentami(nowyAdapterTozsamosciAgenta(s.repozytoria.Agenci, s.repozytoria.WarstwyAgenta)).
		ZKatalogiemRoboczym(s.katalog).
		// Dziennik idzie do składacza mostów po jedno: gdy w produkcie brakuje
		// binarium serwera narzędzi, rdzeń zapisuje to w dzienniku, zamiast
		// wpisywać do MCP ścieżkę, której nie ma (most_narzedzi.go).
		ZMostami(noweMostyOkna(s.repozytoria.Nadania, s.repozytoria.PunktyDostepu, s.repozytoria.Okna).
			ZDziennikiem(s.montaz.Dziennik)).
		// Spoina konfiguracja → wykonanie. Rozstrzygacz podaje nakład rozumowania
		// i kanał zapasowy, rejestr kanałów tłumaczy kod kanału na identyfikator
		// modelu, którego oczekuje przełącznik dostawcy.
		ZParametramiWykonania(noweWykonanieZKonfiguracji(s.rozstrzygacz, s.repozytoria.Kanaly))

	// Czytelnik obowiązującej konfiguracji sesji. Ten sam adapter osi, który
	// wypełnia port Ustawienia rdzenia, tłumaczy tu obszary konfiguracji na
	// wejście procesu modelu — config.session.set zmienia więc zbudowane
	// wywołanie. Port pusty nie blokuje tury.
	if s.ustawienia != nil {
		rozmowa.ZKonfiguracjaSesji(s.ustawienia)
	}

	// Doraźne dołożenia sesji — drugie źródło zestawu narzędzi tury
	// (`adapter_rozmowa_zestaw.go`). Pierwszym jest ekspert okna, który jedzie
	// polem `Ustawienia.Agent` i portu nie potrzebuje.
	//
	// WPIĘCIE JEST WARUNKOWE I MUSI TAKIE BYĆ. Port trzymany jest interfejsem,
	// a wskaźnik zerowy włożony do interfejsu NIESIE TYP, więc porównanie
	// `port == nil` u wołającego przechodzi i tura woła metodę na niczym.
	// Wołanie bezwarunkowe zamieniało tu szczery brak wpięcia — który adapter
	// rozmowy umie zameldować i przeżyć — w padnięcie całego rdzenia przy
	// pierwszym `message.send`. Ta jedna gwiazdka rozstrzyga o tym, czy brak
	// portu jest zdaniem w dzienniku, czy zgaszonym procesem.
	if s.dolozenia != nil {
		rozmowa.ZDolozeniamiSesji(s.dolozenia)
	}

	return rozmowa, petla
}
