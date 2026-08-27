package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
)

// skladRozmowy jest kompletem, z którego powstaje domena rozmowy: pętla koordynator–wykonawca
// wraz z adapterem i wszystkimi jego portami, zamiast dziesięciu argumentów pozycyjnych.
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
	// ustawienia odczytuje konfigurację obowiązującą sesji — jeden czytelnik konfiguracji platformy.
	ustawienia *adapterUstawienOsi
	// zdarzenia utrwala zamknięcia tur i zdarzenia zaczepów jedną drogą śladu z fabryki kanału.
	zdarzenia *zdarzeniaWykonawcze
	// dolozenia wnosi doraźne dołożenia sesji do zestawu narzędzi tym samym adapterem co port sesji.
	dolozenia *adapterNarzedziSesji
}

// zlozRozmowe składa adapter rozmowy wraz z pętlą koordynator–wykonawca i wpina komplet portów,
// bo funkcja ma jedną odpowiedzialność: wszystko, co adapter musi dostać przed prowadzeniem tur.
func zlozRozmowe(s skladRozmowy) (*adapterRozmowy, *session.Petla) {
	rozmowa := nowyAdapterRozmowy(s.zycie, s.nadzorca, s.kanaly, s.nadajnik, s.dziennik).
		ZCiagloscia(ciagloscRozmowy(s.repozytoria)).
		ZZdarzeniami(s.zdarzenia).
		// Rejestrator bloków dostaje kontekst życia rdzenia; zapis dobiega końca po rozłączeniu klienta.
		ZBlokami(nowyRejestratorBlokow(s.zycie, s.repozytoria.Bloki, s.montaz.Dziennik)).
		// Magazyn załączników rozmowy stoi na tym samym katalogu danych co sejf poświadczeń i biblioteki.
		ZZalacznikami(s.montaz.Konfiguracja.KatalogDanych)

	petla := session.NowaPetla(s.nadzorca,
		session.UruchomienieFunkcja(rozmowa.RozpocznijObieg), session.UstawieniaPetli{})
	rozmowa.UstawPetle(petla)
	petla.Obserwuj(session.ObserwatorFunkcja(sladObiegu(s.montaz.Dziennik)))
	// Bieg naprawczy idzie do dziennika rdzenia i do kontrolki sesji trwającej w tle.
	petla.Obserwuj(s.biegi)
	s.obecnosc.ZeStrumieniem(rozmowa.CzyTuraWBiegu)

	// Tożsamość dojeżdża do procesu modelu tą samą drogą, którą pokazuje ją okno konfiguracji.
	rozmowa.ZTozsamoscia(s.tozsamosc).
		// Ekspert wskazany w oknie wnosi warstwy promptu, model bazowy i nastawy tym samym składaczem.
		ZAgentami(nowyAdapterTozsamosciAgenta(s.repozytoria.Agenci, s.repozytoria.WarstwyAgenta)).
		ZKatalogiemRoboczym(s.katalog).
		// Dziennik idzie do składacza mostów, gdy brakuje binarium serwera narzędzi.
		ZMostami(noweMostyOkna(s.repozytoria.Nadania, s.repozytoria.PunktyDostepu, s.repozytoria.Okna).
			ZDziennikiem(s.montaz.Dziennik)).
		// Rozstrzygacz podaje nakład rozumowania i kanał zapasowy; rejestr tłumaczy kod na model.
		ZParametramiWykonania(noweWykonanieZKonfiguracji(s.rozstrzygacz, s.repozytoria.Kanaly))

	// Czytelnik konfiguracji sesji tłumaczy obszary na wejście modelu; port pusty nie blokuje tury.
	if s.ustawienia != nil {
		rozmowa.ZKonfiguracjaSesji(s.ustawienia)
	}

	// Wpięcie jest warunkowe: wskaźnik zerowy w interfejsie niesie typ, więc bezwarunek by padł.
	if s.dolozenia != nil {
		rozmowa.ZDolozeniamiSesji(s.dolozenia)
	}

	return rozmowa, petla
}
