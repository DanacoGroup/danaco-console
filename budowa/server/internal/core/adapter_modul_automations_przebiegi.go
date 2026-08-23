// Odpowiedzialność pliku: okno Execution Monitor — stan przebiegów automatyki,
// zapis obserwacji i wyprowadzenie stanu przebiegu z kolejki, która go wykonuje.
//
// Stan przebiegu jest wyprowadzony, nie deklarowany. Etapem przebiegu jest
// pozycja kolejki, więc etap bieżący, liczba etapów i numer próby biorą się
// z pozycji — dokładnie tak, jak telemetria postępu liczy etapy kolejki. Rdzeń
// nie prowadzi drugiego licznika i nie zmyśla stopnia ukończenia.
//
// Licznik obiegów nie ma granicy: numer próby to najwyższy licznik obiegów
// wśród pozycji, a progu, po którym rdzeń odmawia powtórzenia, nie ma —
// przerwanie należy do Operatora, a przycisk zatrzymania jest zawsze czynny.
package core

import (
	"context"
	"errors"
	"strconv"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Przebiegi zapisuje okno na telemetrię przebiegów i oddaje ich stan bieżący.
func (a *adapterAutomatyk) Przebiegi(ctx context.Context,
	z shared.AutomationExecutionSubscribeRequest) (shared.AutomationExecutionSubscribeResponse, error) {

	kodAutomatyki := wartoscTekstu(z.WorkflowId)
	var automatykaID int64
	if kodAutomatyki != "" {
		wiersz, err := a.wiersz(ctx, kodAutomatyki)
		if err != nil {
			return shared.AutomationExecutionSubscribeResponse{}, err
		}
		automatykaID = wiersz.ID
	}
	wiersze, err := a.wierszePrzebiegow(ctx, z, automatykaID)
	if err != nil {
		return shared.AutomationExecutionSubscribeResponse{}, err
	}
	zapisane := a.obserwatorzy.Zapamietaj(wartoscTekstu(z.WindowId), kodAutomatyki)
	przebiegi := make([]shared.AutomationExecution, 0, len(wiersze))
	for _, wiersz := range wiersze {
		przebiegi = append(przebiegi, przebiegKontraktu(wiersz))
	}
	return shared.AutomationExecutionSubscribeResponse{
		Executions: przebiegi, Subscribed: zapisane,
	}, nil
}

// wierszePrzebiegow dobiera przebiegi wedle wskazania żądania: pojedynczy
// przebieg, przebiegi automatyki albo wszystkie.
func (a *adapterAutomatyk) wierszePrzebiegow(ctx context.Context,
	z shared.AutomationExecutionSubscribeRequest, automatykaID int64) ([]dane.Przebieg, error) {

	if kod := wartoscTekstu(z.ExecutionId); kod != "" {
		wiersz, err := a.repozytorium.Przebieg(ctx, kod)
		if err != nil {
			return nil, bladNieznanegoPrzebiegu(kod, err)
		}
		return []dane.Przebieg{wiersz}, nil
	}
	wiersze, err := a.repozytorium.Przebiegi(ctx, automatykaID, wartoscLiczby(z.Limit))
	if err != nil {
		return nil, bladAutomatyki(err)
	}
	return wiersze, nil
}

// PrzebiegKolejki oddaje przebieg realizowany przez wskazaną kolejkę. Służy
// obsługiwaczowi do rozgłoszenia `automation.execution.status` po działaniu na
// kolejce — kolejka bez przebiegu nie rozgłasza niczego.
func (a *adapterAutomatyk) PrzebiegKolejki(ctx context.Context,
	idKolejki string) (shared.AutomationExecution, bool) {

	id, err := strconv.ParseInt(idKolejki, 10, 64)
	if err != nil {
		return shared.AutomationExecution{}, false
	}
	wiersz, err := a.repozytorium.PrzebiegKolejki(ctx, id)
	if err != nil {
		return shared.AutomationExecution{}, false
	}
	return przebiegKontraktu(wiersz), true
}

// odnotujPrzebieg zapisuje stan przebiegu po działaniu na kolejce. Kolejka bez
// wskazanej automatyki, która nie ma jeszcze przebiegu, przebiegu nie zakłada:
// nie każda kolejka wykonuje automatykę, bo silnik jest jeden.
func (a *adapterAutomatyk) odnotujPrzebieg(ctx context.Context, kolejkaID int64,
	automatyka *dane.Automatyka, kolejka shared.Queue) error {

	zastany, err := a.repozytorium.PrzebiegKolejki(ctx, kolejkaID)
	istnieje := err == nil
	if !istnieje && automatyka == nil {
		return nil
	}
	pozycje, err := a.kolejki.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return err
	}
	przebieg := dane.Przebieg{Kod: zastany.Kod, AutomatykaID: zastany.AutomatykaID, KolejkaID: &kolejkaID}
	if !istnieje {
		przebieg = dane.Przebieg{
			Kod: nowyIdentyfikator(przedrostekPrzebiegu), AutomatykaID: automatyka.ID,
			KolejkaID: &kolejkaID,
		}
	}
	naniesPostep(&przebieg, kolejka.Status, pozycje)
	_, err = a.repozytorium.ZapiszPrzebieg(ctx, przebieg)
	return err
}

// naniesPostep wypełnia przebieg stanem wyprowadzonym z kolejki i jej pozycji.
func naniesPostep(przebieg *dane.Przebieg, stanKolejki shared.QueueStatus, pozycje []dane.Pozycja) {
	zakonczone, bledne, proba := 0, 0, 0
	for _, pozycja := range pozycje {
		if czyStanKoncowyPozycji(pozycja.Stan) {
			zakonczone++
		}
		if pozycja.Stan == stanPozycjiBledna {
			bledne++
		}
		if pozycja.LicznikObiegow > proba {
			proba = pozycja.LicznikObiegow
		}
	}
	przebieg.Etapow = len(pozycje)
	przebieg.EtapBiezacy = zakonczone
	if zakonczone < len(pozycje) {
		przebieg.EtapBiezacy = zakonczone + 1
	}
	przebieg.Proba = proba
	przebieg.Stan = stanPrzebiegu(stanKolejki, bledne)
	if czyStanKoncowyPrzebiegu(przebieg.Stan) {
		zakonczono := time.Now().UTC().Format(formatZnacznikaBazy)
		przebieg.Zakonczono = &zakonczono
	}
	if bledne > 0 {
		powod := "zleceń zakończonych błędem: " + strconv.Itoa(bledne)
		przebieg.KomunikatBledu = &powod
	}
}

// stanPrzebiegu przekłada stan kolejki na stan przebiegu automatyki. Kolejka
// wyczerpana, w której choć jedno zlecenie skończyło się błędem, daje przebieg
// nieudany — inaczej Execution Monitor meldowałby powodzenie wykonania, które
// się nie powiodło.
func stanPrzebiegu(stanKolejki shared.QueueStatus, bledne int) string {
	switch stanKolejki {
	case shared.QueueStatusRunning:
		return shared.AutomationExecutionStatusRunning
	case shared.QueueStatusPaused:
		return shared.AutomationExecutionStatusPaused
	case shared.QueueStatusStopped:
		return shared.AutomationExecutionStatusStopped
	case shared.QueueStatusDone:
		if bledne > 0 {
			return shared.AutomationExecutionStatusFailed
		}
		return shared.AutomationExecutionStatusSucceeded
	default:
		return shared.AutomationExecutionStatusPending
	}
}

// czyStanKoncowyPrzebiegu mówi, czy przebieg zamknął się na dobre.
func czyStanKoncowyPrzebiegu(stan string) bool {
	return stan == shared.AutomationExecutionStatusSucceeded ||
		stan == shared.AutomationExecutionStatusFailed ||
		stan == shared.AutomationExecutionStatusStopped
}

// przebiegKontraktu przekłada wiersz przebiegu na byt kontraktu.
func przebiegKontraktu(wiersz dane.Przebieg) shared.AutomationExecution {
	etap, etapow, proba := wiersz.EtapBiezacy, wiersz.Etapow, wiersz.Proba
	przebieg := shared.AutomationExecution{
		Id: wiersz.Kod, WorkflowId: wiersz.AutomatykaKod,
		Status:  shared.AutomationExecutionStatus(wiersz.Stan),
		Attempt: &proba, ErrorMessage: wiersz.KomunikatBledu,
		StartedAt: chwilaBazy(wiersz.Rozpoczeto),
	}
	if etapow > 0 {
		przebieg.CurrentStep, przebieg.TotalSteps = &etap, &etapow
	}
	if wiersz.Zakonczono != nil {
		zakonczono := chwilaBazy(*wiersz.Zakonczono)
		przebieg.FinishedAt = &zakonczono
	}
	return przebieg
}

// bladNieznanegoPrzebiegu odróżnia „przebiegu nie ma” od usterki odczytu.
func bladNieznanegoPrzebiegu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Automations: przebieg nie istnieje: "+kod))
	}
	return bladAutomatyki(err)
}
