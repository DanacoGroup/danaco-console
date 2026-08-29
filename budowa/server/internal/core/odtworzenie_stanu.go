package core

import (
	"context"
	"log"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
)

// odtworzStanZBazy odtwarza w rejestrze rdzenia sesje i okna zapisane w bazie po restarcie, pod tymi samymi identyfikatorami zewnętrznymi, bez uruchamiania procesów okien.
func odtworzStanZBazy(kontekst context.Context, repozytoria *dane.Zestaw,
	nadzorca *session.Nadzorca, dziennik *log.Logger) {

	if repozytoria == nil || nadzorca == nil {
		return
	}
	sesje, err := repozytoria.Sesje.Lista(kontekst, 0)
	if err != nil {
		odnotujOdtworzenie(dziennik, "odtworzenie sesji z bazy: %v", err)
		return
	}
	moduly, kanaly, err := slownikiOkna(kontekst, repozytoria)
	if err != nil {
		// Brak słowników nie przerywa odtwarzania: okno wraca bez kodu modułu i kanału.
		odnotujOdtworzenie(dziennik, "słowniki okien przy odtwarzaniu: %v", err)
		moduly, kanaly = map[int64]string{}, map[int64]string{}
	}
	for _, wiersz := range sesje {
		wniesSesjeDoRejestru(kontekst, repozytoria, nadzorca, dziennik, wiersz, moduly, kanaly)
	}
}

// wniesSesjeDoRejestru odtwarza w rejestrze jedną sesję wraz z jej oknami.
// Wiersz bez identyfikatora rdzenia zostaje pominięty — powstał wprost w bazie
// i klient nie ma w co trafić. Zwraca, czy sesja weszła do rejestru.
func wniesSesjeDoRejestru(kontekst context.Context, repozytoria *dane.Zestaw,
	nadzorca *session.Nadzorca, dziennik *log.Logger, wiersz dane.Sesja,
	moduly, kanaly map[int64]string) bool {

	if wiersz.IdentyfikatorZewnetrzny == nil {
		return false
	}
	wiersze, err := repozytoria.Okna.ListaSesji(kontekst, wiersz.ID)
	if err != nil {
		odnotujOdtworzenie(dziennik, "odtworzenie okien sesji %d: %v", wiersz.ID, err)
		return false
	}
	sesja := session.Sesja{
		Id: *wiersz.IdentyfikatorZewnetrzny, Tytul: wiersz.Tytul,
		IdProjektu: wartoscTekstu(wiersz.Projekt), Stan: wiersz.Stan,
		/* Znaczniki czasu przenoszone z bazy: bez nich sesja odtworzona wchodzi
		   do rejestru z chwilą zerową i okno pokazuje Operatorowi rok pierwszy
		   zamiast dnia, w którym pracę zaczął. */
		Utworzono:      chwilaZBazy(wiersz.Utworzono),
		Zaktualizowano: chwilaZBazy(wiersz.Zaktualizowano),
	}
	nadzorca.Rejestr().Odtworz(sesja, oknaOdtworzone(wiersze, sesja.Id, moduly, kanaly))
	return true
}

// odtworzSesjePoIdentyfikatorze wnosi do rejestru jedną sesję odczytaną z bazy po identyfikatorze rdzenia, tą samą drogą, którą wnosi start rdzenia.
func odtworzSesjePoIdentyfikatorze(kontekst context.Context, repozytoria *dane.Zestaw,
	nadzorca *session.Nadzorca, dziennik *log.Logger, identyfikator string) bool {

	if repozytoria == nil || nadzorca == nil {
		return false
	}
	wiersz, err := repozytoria.Sesje.PoIdentyfikatorze(kontekst, identyfikator)
	if err != nil {
		odnotujOdtworzenie(dziennik, "odtworzenie sesji %s z bazy: %v", identyfikator, err)
		return false
	}
	moduly, kanaly, err := slownikiOkna(kontekst, repozytoria)
	if err != nil {
		odnotujOdtworzenie(dziennik, "słowniki okien przy odtwarzaniu %s: %v", identyfikator, err)
		moduly, kanaly = map[int64]string{}, map[int64]string{}
	}
	return wniesSesjeDoRejestru(kontekst, repozytoria, nadzorca, dziennik, wiersz, moduly, kanaly)
}

// oknaOdtworzone przekłada wiersze okien na byty rejestru. Wiersz bez
// identyfikatora rdzenia zostaje pominięty — powstał wprost w bazie i nie ma
// odpowiednika, w który klient mógłby trafić.
func oknaOdtworzone(wiersze []dane.Okno, idSesji string,
	moduly, kanaly map[int64]string) []session.Okno {

	teraz := time.Now().UTC()
	okna := make([]session.Okno, 0, len(wiersze))
	for _, o := range wiersze {
		if o.IdentyfikatorZewnetrzny == nil {
			continue
		}
		okna = append(okna, session.Okno{
			Id: *o.IdentyfikatorZewnetrzny, IdSesji: idSesji, Stan: o.Stan,
			Ustawienia: session.Ustawienia{
				Modul:               moduly[o.ModulID],
				KanalModelu:         kanaly[o.KanalModeluID],
				KatalogiRobocze:     o.KatalogiRobocze,
				SrodowiskoWykonania: o.SrodowiskoWykonania,
				TrybUprawnien:       o.TrybUprawnien,
				RolaOkna:            o.RolaOkna,
				Tytul:               wartoscTekstu(o.Tytul),
				// Bez tego wybór eksperta przeżywa zapis, ale nie restart
				// rdzenia — okno wraca wtedy na model surowy.
				Agent: wartoscTekstu(o.AgentKod),
			},
			Utworzono: teraz, Zaktualizowano: teraz,
		})
	}
	return okna
}

// odnotujOdtworzenie zapisuje w dzienniku rdzenia niepowodzenie odczytu napotkane podczas odtwarzania stanu z bazy po restarcie.
func odnotujOdtworzenie(dziennik *log.Logger, wzor string, argumenty ...any) {
	if dziennik == nil {
		return
	}
	dziennik.Printf(wzor, argumenty...)
}

/*
chwilaZBazy odczytuje znacznik czasu zapisany w bazie napisem.

Postać zapisu to RFC 3339 (`2026-08-29T20:02:04.260Z`). Napis nieczytelny daje
chwilę zerową — tak samo, jak działo się przed przeniesieniem znaczników; lepsza
jest jedna wartość pusta niż zatrzymanie odtworzenia całej sesji.
*/
func chwilaZBazy(zapis string) time.Time {
	if zapis == "" {
		return time.Time{}
	}
	chwila, err := time.Parse(time.RFC3339Nano, zapis)
	if err != nil {
		return time.Time{}
	}
	return chwila
}
