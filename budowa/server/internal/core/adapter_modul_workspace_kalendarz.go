// Odpowiedzialność pliku: kalendarz projektu ze zadań i wciągniętych plików iCal —
// `workspace.calendar.get`, `workspace.calendar.import` i `workspace.calendar.export`.
package core

import (
	"context"
	"sort"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Kalendarz obsługuje `workspace.calendar.get` i zwraca pozycje z obu źródeł projektu naraz, posortowane.
func (a *adapterPrzestrzeniRoboczej) Kalendarz(ctx context.Context,
	z shared.WorkspaceCalendarGetRequest) (shared.WorkspaceCalendarGetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceCalendarGetResponse{}, err
	}
	od, do := okresSiatkiWorkspace(z.Span, z.AnchorAt)
	pozycje, err := a.pozycjeKalendarzaWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceCalendarGetResponse{}, err
	}
	wOkresie := []shared.WorkspaceCalendarEntry{}
	for _, pozycja := range pozycje {
		koniec := pozycja.StartAt
		if pozycja.EndAt != nil {
			koniec = *pozycja.EndAt
		}
		if koniec < od || pozycja.StartAt > do {
			continue
		}
		wOkresie = append(wOkresie, pozycja)
	}
	return shared.WorkspaceCalendarGetResponse{Entries: wOkresie, From: od, To: do}, nil
}

// WciagnijKalendarz obsługuje `workspace.calendar.import` i wciąga pozycje z przesłanego pliku iCal do projektu.
func (a *adapterPrzestrzeniRoboczej) WciagnijKalendarz(ctx context.Context,
	z shared.WorkspaceCalendarImportRequest) (shared.WorkspaceCalendarImportResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceCalendarImportResponse{}, err
	}
	if strings.TrimSpace(z.Content) == "" {
		return shared.WorkspaceCalendarImportResponse{},
			bladProjektu("wciągnięcie kalendarza bez treści pliku")
	}
	wydarzenia, powody := rozbierzIcalWorkspace(z.Content)
	jakoZadania := z.AsTasks == nil || *z.AsTasks

	wciagniete := 0
	kodyZadan := []string{}
	for _, wydarzenie := range wydarzenia {
		uid := pierwszaNiepustaWorkspace(wydarzenie.Uid,
			uidIcalWorkspace(wydarzenie.Tytul, wydarzenie.PoczatekMs))
		if jakoZadania {
			odpowiedz, err := a.ZalozZadanie(ctx, shared.WorkspaceTaskCreateRequest{
				ProjectId: projekt.Kod,
				Title:     wydarzenie.Tytul,
				StartAt:   chwilaKontraktuWorkspace(wydarzenie.PoczatekMs),
				DueAt: chwilaKontraktuWorkspace(pierwszaChwilaWorkspace(wydarzenie.KoniecMs,
					wydarzenie.PoczatekMs)),
				RecurrenceRule: tekstOpcjonalnyWorkspace(wydarzenie.Regula),
				Description:    tekstOpcjonalnyWorkspace("wciągnięte z kalendarza, UID " + uid),
			})
			if err != nil {
				powody = append(powody, "wydarzenia „"+wydarzenie.Tytul+
					"” nie dało się zapisać jako zadania: "+err.Error())
				continue
			}
			kodyZadan = append(kodyZadan, odpowiedz.Task.Id)
			wciagniete++
			continue
		}
		err := a.repozytorium.ZapiszPozycjeKalendarzaWorkspace(ctx, dane.PozycjaKalendarzaWorkspace{
			ProjektID: projekt.ID, Identyfikator: nowyIdentyfikator("wscal-"),
			Tytul: wydarzenie.Tytul, PoczatekMs: wydarzenie.PoczatekMs,
			KoniecMs: wydarzenie.KoniecMs, CalyDzien: wydarzenie.CalyDzien,
			RegulaPowtarzalnosci: wydarzenie.Regula, UidZewnetrzny: uid,
		})
		if err != nil {
			powody = append(powody, "pozycji „"+wydarzenie.Tytul+"” nie dało się zapisać: "+err.Error())
			continue
		}
		wciagniete++
	}
	if wciagniete > 0 {
		if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
			return shared.WorkspaceCalendarImportResponse{}, err
		}
		a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindTaskChanged,
			shared.ChangeKindCreated, shared.WorkspaceEntityKindProject, projekt.Kod,
			"wciągnięto kalendarz "+pierwszaNiepustaWorkspace(wartoscTekstuWorkspace(z.FileName),
				"iCal"))
	}
	return shared.WorkspaceCalendarImportResponse{
		Imported: wciagniete, Skipped: len(powody), SkippedReasons: powody, TaskIds: kodyZadan,
	}, nil
}

// ZapiszKalendarz obsługuje `workspace.calendar.export` i składa kalendarz projektu w postaci pliku iCal.
func (a *adapterPrzestrzeniRoboczej) ZapiszKalendarz(ctx context.Context,
	z shared.WorkspaceCalendarExportRequest) (shared.WorkspaceCalendarExportResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceCalendarExportResponse{}, err
	}
	pozycje, err := a.pozycjeKalendarzaWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceCalendarExportResponse{}, err
	}
	wydarzenia := []wydarzenieIcalWorkspace{}
	for _, pozycja := range pozycje {
		if z.From != nil && pozycja.StartAt < *z.From {
			continue
		}
		if z.To != nil && pozycja.StartAt > *z.To {
			continue
		}
		wydarzenie := wydarzenieIcalWorkspace{
			Uid:        pierwszaNiepustaWorkspace(wartoscTekstuWorkspace(pozycja.ExternalUid), pozycja.Id),
			Tytul:      pozycja.Title,
			PoczatekMs: pozycja.StartAt,
			CalyDzien:  pozycja.AllDay != nil && *pozycja.AllDay,
			Regula:     wartoscTekstuWorkspace(pozycja.RecurrenceRule),
		}
		if pozycja.EndAt != nil {
			wydarzenie.KoniecMs = *pozycja.EndAt
		}
		wydarzenia = append(wydarzenia, wydarzenie)
	}
	tresc := zlozIcalWorkspace("Projekt "+projekt.Nazwa, wydarzenia)
	return shared.WorkspaceCalendarExportResponse{
		Content: tresc, FileName: "projekt-" + projekt.Kod + ".ics", Exported: len(wydarzenia),
	}, nil
}

// pozycjeKalendarzaWorkspace składa kalendarz projektu z obu jego źródeł:
// zadań z terminem oraz pozycji wciągniętych z iCal.
func (a *adapterPrzestrzeniRoboczej) pozycjeKalendarzaWorkspace(ctx context.Context,
	projektID int64) ([]shared.WorkspaceCalendarEntry, error) {

	zadania, err := a.repozytorium.ZadaniaWorkspace(ctx, projektID)
	if err != nil {
		return nil, err
	}
	pozycje := []shared.WorkspaceCalendarEntry{}
	for _, zadanie := range zadania {
		poczatek := pierwszaChwilaWorkspace(zadanie.PoczatekMs, zadanie.TerminMs)
		if poczatek == 0 {
			continue
		}
		kod := zadanie.Identyfikator
		kamien := zadanie.KamienMilowy
		pozycja := shared.WorkspaceCalendarEntry{
			Id: "zadanie-" + kod, TaskId: &kod, Title: zadanie.Tytul,
			StartAt: poczatek, Milestone: &kamien,
		}
		if zadanie.TerminMs != 0 {
			koniec := zadanie.TerminMs
			pozycja.EndAt = &koniec
		}
		if zadanie.RegulaPowtarzalnosci != "" {
			regula := zadanie.RegulaPowtarzalnosci
			pozycja.RecurrenceRule = &regula
		}
		pozycje = append(pozycje, pozycja)
	}
	wciagniete, err := a.repozytorium.PozycjeKalendarzaWorkspace(ctx, projektID)
	if err != nil {
		return nil, err
	}
	for _, wiersz := range wciagniete {
		calyDzien, kamien := wiersz.CalyDzien, wiersz.KamienMilowy
		pozycja := shared.WorkspaceCalendarEntry{
			Id: wiersz.Identyfikator, Title: wiersz.Tytul, StartAt: wiersz.PoczatekMs,
			AllDay: &calyDzien, Milestone: &kamien,
		}
		if wiersz.KoniecMs != 0 {
			koniec := wiersz.KoniecMs
			pozycja.EndAt = &koniec
		}
		if wiersz.RegulaPowtarzalnosci != "" {
			regula := wiersz.RegulaPowtarzalnosci
			pozycja.RecurrenceRule = &regula
		}
		if wiersz.UidZewnetrzny != "" {
			uid := wiersz.UidZewnetrzny
			pozycja.ExternalUid = &uid
		}
		pozycje = append(pozycje, pozycja)
	}
	sort.SliceStable(pozycje, func(i, j int) bool { return pozycje[i].StartAt < pozycje[j].StartAt })
	return pozycje, nil
}

// okresSiatkiWorkspace wyznacza granice okresu objętego siatką kalendarza.
// Chwila kotwicy należy do okresu, a granice liczone są w czasie UTC — tej
// samej mierze, w której kontrakt podaje wszystkie chwile.
func okresSiatkiWorkspace(siatka shared.WorkspaceCalendarSpan, kotwica int64) (int64, int64) {
	if kotwica == 0 {
		kotwica = terazWMilisekundachWorkspace()
	}
	chwila := time.UnixMilli(kotwica).UTC()
	dzien := time.Date(chwila.Year(), chwila.Month(), chwila.Day(), 0, 0, 0, 0, time.UTC)
	switch siatka {
	case shared.WorkspaceCalendarSpanDay:
		return dzien.UnixMilli(), dzien.AddDate(0, 0, 1).Add(-time.Millisecond).UnixMilli()
	case shared.WorkspaceCalendarSpanWeek:
		// Tydzień zaczyna się w poniedziałek, tak liczy kalendarz tego produktu.
		przesuniecie := (int(dzien.Weekday()) + 6) % 7
		poczatek := dzien.AddDate(0, 0, -przesuniecie)
		return poczatek.UnixMilli(), poczatek.AddDate(0, 0, 7).Add(-time.Millisecond).UnixMilli()
	default:
		poczatek := time.Date(chwila.Year(), chwila.Month(), 1, 0, 0, 0, 0, time.UTC)
		return poczatek.UnixMilli(), poczatek.AddDate(0, 1, 0).Add(-time.Millisecond).UnixMilli()
	}
}

// tekstOpcjonalnyWorkspace zwija pusty napis do braku pola kontraktu, zamiast oddawać pusty łańcuch znaków.
func tekstOpcjonalnyWorkspace(wartosc string) *string {
	if strings.TrimSpace(wartosc) == "" {
		return nil
	}
	tekst := wartosc
	return &tekst
}
