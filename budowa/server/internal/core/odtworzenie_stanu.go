package core

import (
	"context"
	"log"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
)

func odtworzStanZBazy(kontekst context.Context, repozytoria *dane.Zestaw,
	nadzorca *session.Nadzorca, dziennik *log.Logger) {

	if repozytoria == nil || nadzorca == nil {
		return
	}
	konteksty, err := kontekstyKont(kontekst, repozytoria.KontoWlasciciela)
	if err != nil {
		odnotujOdtworzenie(dziennik, "konta właściciela przy odtwarzaniu: %v", err)
		return
	}
	for _, kontekstKonta := range konteksty {
		sesje, err := repozytoria.Sesje.Lista(kontekstKonta, 0)
		if err != nil {
			odnotujOdtworzenie(dziennik, "odtworzenie sesji z bazy: %v", err)
			continue
		}
		moduly, kanaly, err := slownikiOkna(kontekstKonta, repozytoria)
		if err != nil {
			// Brak słowników nie przerywa odtwarzania: okno wraca bez kodu modułu i kanału.
			odnotujOdtworzenie(dziennik, "słowniki okien przy odtwarzaniu: %v", err)
			moduly, kanaly = map[int64]string{}, map[int64]string{}
		}
		for _, wiersz := range sesje {
			wniesSesjeDoRejestru(kontekstKonta, repozytoria, nadzorca, dziennik, wiersz, moduly, kanaly)
		}
	}
}

// Praca procesu bez zamawiającego nie podszywa się pod konto — idzie po kontach
// właściciela po kolei (decyzja 34); pusta tabela daje jeden kontekst bez wskazania.
func kontekstyKont(kontekst context.Context,
	konta dane.RepozytoriumKontaWlasciciela) ([]context.Context, error) {

	if konta == nil {
		return []context.Context{kontekst}, nil
	}
	wiersze, err := konta.Konta(kontekst)
	if err != nil {
		return nil, err
	}
	if len(wiersze) == 0 {
		return []context.Context{kontekst}, nil
	}
	konteksty := make([]context.Context, 0, len(wiersze))
	for _, konto := range wiersze {
		konteksty = append(konteksty, dane.ZKontemOperatora(kontekst, konto.Id))
	}
	return konteksty, nil
}

// Wiersz bez identyfikatora rdzenia powstał wprost w bazie i klient nie ma w co trafić.
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
		// Bez znaczników z bazy sesja odtworzona pokazałaby Operatorowi rok pierwszy zamiast dnia startu pracy.
		Utworzono:      chwilaZBazy(wiersz.Utworzono),
		Zaktualizowano: chwilaZBazy(wiersz.Zaktualizowano),
		// Odtworzenie idzie kontekstem jednego konta; bez tego wpisu sesja
		// wracałaby do rejestru bez granicy i weszła do cudzego wykazu.
		KontoId: dane.KontoOperatora(kontekst),
	}
	nadzorca.Rejestr().Odtworz(sesja, oknaOdtworzone(wiersze, sesja.Id, moduly, kanaly))
	return true
}

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
				// Bez tego wybór eksperta przeżywa zapis, ale nie restart rdzenia.
				Agent: wartoscTekstu(o.AgentKod),
			},
			Utworzono: teraz, Zaktualizowano: teraz,
		})
	}
	return okna
}

func odnotujOdtworzenie(dziennik *log.Logger, wzor string, argumenty ...any) {
	if dziennik == nil {
		return
	}
	dziennik.Printf(wzor, argumenty...)
}

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
