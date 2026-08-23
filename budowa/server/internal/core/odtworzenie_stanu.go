package core

import (
	"context"
	"log"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
)

// Odtworzenie rejestru sesji i okien po restarcie rdzenia.
//
// Rejestr nadzorcy jest pamięcią jednego uruchomienia. Bez tego kroku rozmowa
// sprzed restartu wraca wyłącznie historią wiadomości, a sam byt sesji i okna
// znika: klient trzyma identyfikator, pod którym nie ma już czego wskazać.
// Odtworzenie wnosi sesje i okna z bazy POD TYMI SAMYMI identyfikatorami
// (`identyfikator_zewnetrzny`, `store/migracja_006_identyfikatory_zewnetrzne.sql`),
// więc `session.bind`,
// `window.state.get` i `message.send` trafiają w te same byty.
//
// Czego odtworzenie NIE robi: nie startuje procesów okien. Proces ginie wraz
// z rdzeniem, a okno wraca jako byt bez procesu — rejestr procesów zgłosi wtedy
// stan `pending`, a pierwsza tura uruchomi proces zwykłą drogą.

// odtworzStanZBazy wnosi do rejestru rdzenia sesje i okna zapisane w bazie.
// Błąd odczytu nie przerywa montażu: rdzeń rusza z rejestrem pustym.
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
		// Brak słowników nie przerywa odtwarzania — okno wraca wtedy bez kodu
		// modułu i bez kodu kanału, a Operator wskaże je ponownie.
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
	}
	nadzorca.Rejestr().Odtworz(sesja, oknaOdtworzone(wiersze, sesja.Id, moduly, kanaly))
	return true
}

// odtworzSesjePoIdentyfikatorze wnosi do rejestru jedną sesję odczytaną z bazy
// po identyfikatorze rdzenia — droga POWROTU Z KOSZA (session.restore,
// adapter_sesje_kosz.go): usunięcie zdjęło sesję z rejestru żywego, a wykaz
// startowy sesji z kosza nie widział, więc przywrócenie musi ją wnieść samo,
// tą samą drogą, którą wnosi start.
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

// odnotujOdtworzenie zapisuje niepowodzenie odczytu w dzienniku rdzenia.
func odnotujOdtworzenie(dziennik *log.Logger, wzor string, argumenty ...any) {
	if dziennik == nil {
		return
	}
	dziennik.Printf(wzor, argumenty...)
}
