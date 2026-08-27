// Odpowiedzialność pliku: sześć czynności Execution Monitora sięgających
// trwałego zapisu przebiegu — log, kroki, punkty wznowienia, wznowienie,
// podgląd ładunku i odtworzenie. Wszystkie czytają bazę, nie pamięć procesu.
package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// nazwyPolWrazliwych wylicza pola ładunku maskowane regułą redakcji sekretów.
// Dopasowanie idzie po fragmencie nazwy pola bez względu na wielkość liter, bo
// ładunek kroku bywa nazwany przez interfejs zewnętrzny, a nie przez rdzeń.
var nazwyPolWrazliwych = []string{
	"secret", "sekret", "password", "haslo", "token", "apikey", "api_key",
	"authorization", "credential", "poswiadczenie", "privatekey",
}

// DziennikPrzebiegu oddaje pełny zapis zdarzeń jednego uruchomienia, czytany
// z bazy, nie z telemetrii WebSocket ograniczonej do otwartego okna.
func (a *adapterAutomatyk) DziennikPrzebiegu(ctx context.Context,
	z shared.AutomationExecutionLogRequest) (shared.AutomationExecutionLogResponse, error) {

	przebieg, err := a.wierszPrzebiegu(ctx, z.ExecutionId)
	if err != nil {
		return shared.AutomationExecutionLogResponse{}, err
	}
	poziom := poziomBazy(z.Level)
	granica := wartoscLiczby(z.Limit)
	// Granica jest podnoszona o jeden, żeby odróżnić „tyle jest” od „przycięto”.
	granicaOdczytu := granica
	if granicaOdczytu > 0 {
		granicaOdczytu++
	}
	wiersze, err := a.repozytorium.LogPrzebiegu(ctx, przebieg.ID, poziom,
		wartoscTekstu(z.StepId), granicaOdczytu)
	if err != nil {
		return shared.AutomationExecutionLogResponse{}, bladAutomatyki(err)
	}
	przyciety := granica > 0 && len(wiersze) > granica
	if przyciety {
		wiersze = wiersze[:granica]
	}
	wpisy := make([]shared.AutomationLogEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wpisy = append(wpisy, shared.AutomationLogEntry{
			ExecutionId: przebieg.Kod, StepId: wiersz.KrokKod,
			Level: poziomKontraktu(wiersz.Poziom), Message: wiersz.Tresc,
			At: chwilaBazy(wiersz.Chwila),
		})
	}
	return shared.AutomationExecutionLogResponse{Entries: wpisy, Truncated: przyciety}, nil
}

// KrokiPrzebiegu oddaje stan każdego kroku przebiegu osobno: rozpoczęcie,
// zakończenie, błąd i ślad stosu, gdy krok zawiódł.
func (a *adapterAutomatyk) KrokiPrzebiegu(ctx context.Context,
	z shared.AutomationExecutionStepsRequest) (shared.AutomationExecutionStepsResponse, error) {

	przebieg, err := a.wierszPrzebiegu(ctx, z.ExecutionId)
	if err != nil {
		return shared.AutomationExecutionStepsResponse{}, err
	}
	wiersze, err := a.repozytorium.KrokiPrzebiegu(ctx, przebieg.ID)
	if err != nil {
		return shared.AutomationExecutionStepsResponse{}, bladAutomatyki(err)
	}
	kroki := make([]shared.AutomationExecutionStep, 0, len(wiersze))
	for _, wiersz := range wiersze {
		proba := wiersz.Proba
		krok := shared.AutomationExecutionStep{
			ExecutionId: przebieg.Kod, StepId: wiersz.KrokKod,
			Status: shared.AutomationExecutionStatus(wiersz.Stan), Attempt: &proba,
			ErrorMessage: wiersz.KomunikatBledu, StackTrace: wiersz.SladStosu,
			StartedAt: chwilaBazy(wiersz.Rozpoczeto),
		}
		if wiersz.Zakonczono != nil {
			zakonczono := chwilaBazy(*wiersz.Zakonczono)
			krok.FinishedAt = &zakonczono
		}
		kroki = append(kroki, krok)
	}
	return shared.AutomationExecutionStepsResponse{Steps: kroki}, nil
}

// WykazPunktowWznowienia oddaje punkty wznowienia przebiegu wraz z kodami
// kroków, które każdy punkt uznaje za już ukończone.
func (a *adapterAutomatyk) WykazPunktowWznowienia(ctx context.Context,
	z shared.AutomationExecutionCheckpointListRequest) (shared.AutomationExecutionCheckpointListResponse, error) {

	przebieg, err := a.wierszPrzebiegu(ctx, z.ExecutionId)
	if err != nil {
		return shared.AutomationExecutionCheckpointListResponse{}, err
	}
	wiersze, err := a.repozytorium.PunktyWznowienia(ctx, przebieg.ID)
	if err != nil {
		return shared.AutomationExecutionCheckpointListResponse{}, bladAutomatyki(err)
	}
	punkty := make([]shared.AutomationCheckpoint, 0, len(wiersze))
	for _, wiersz := range wiersze {
		punkty = append(punkty, shared.AutomationCheckpoint{
			Id: wiersz.Kod, ExecutionId: przebieg.Kod, StepId: wiersz.KrokKod,
			CompletedStepIds: kodyKrokowZZapisu(wiersz.KrokiUkonczone),
			CreatedAt:        chwilaBazy(wiersz.Utworzono),
		})
	}
	return shared.AutomationExecutionCheckpointListResponse{Checkpoints: punkty}, nil
}

// WznowPrzebieg wznawia przebieg od punktu wznowienia, bez powtarzania kroków
// ukończonych. Brak wskazania punktu bierze najnowszy — tak mówi kontrakt.
func (a *adapterAutomatyk) WznowPrzebieg(ctx context.Context,
	z shared.AutomationExecutionResumeRequest) (shared.AutomationExecutionResumeResponse, error) {

	przebieg, err := a.wierszPrzebiegu(ctx, z.ExecutionId)
	if err != nil {
		return shared.AutomationExecutionResumeResponse{}, err
	}
	punkt, err := a.punktWznowienia(ctx, przebieg, wartoscTekstu(z.CheckpointId))
	if err != nil {
		return shared.AutomationExecutionResumeResponse{}, err
	}
	// Kroki ukończone przed punktem nie idą po raz drugi: ich zlecenia
	// zamykane są, zanim kolejka ruszy.
	if err := a.domknijKrokiSprzedPunktu(ctx, przebieg, punkt); err != nil {
		return shared.AutomationExecutionResumeResponse{}, bladAutomatyki(err)
	}
	kolejka, err := a.ruszKolejkePrzebiegu(ctx, przebieg, shared.QueueActionResume)
	if err != nil {
		return shared.AutomationExecutionResumeResponse{}, err
	}
	a.zapiszLog(ctx, przebieg.ID, &punkt.KrokKod, shared.AutomationLogLevelInfo,
		"wznowienie przebiegu od punktu "+punkt.Kod)
	a.zapisAudytu(ctx, &przebieg.AutomatykaID, "wznowienie przebiegu",
		map[string]any{"przebieg": przebieg.Kod, "punkt": punkt.Kod, "kolejka": kolejka.Id})

	odswiezony, err := a.repozytorium.Przebieg(ctx, przebieg.Kod)
	if err != nil {
		return shared.AutomationExecutionResumeResponse{}, bladAutomatyki(err)
	}
	return shared.AutomationExecutionResumeResponse{Execution: przebiegKontraktu(odswiezony)}, nil
}

// punktWznowienia dobiera wskazany punkt albo, gdy wskazania brak, najnowszy
// zapisany punkt przebiegu.
func (a *adapterAutomatyk) punktWznowienia(ctx context.Context, przebieg dane.Przebieg,
	kod string) (dane.PunktWznowienia, error) {

	if kod != "" {
		punkt, err := a.repozytorium.PunktWznowieniaPoKodzie(ctx, kod)
		if err != nil {
			return dane.PunktWznowienia{},
				bladNieznanegoBytuAutomatyki(err, "punkt wznowienia nie istnieje: ", kod)
		}
		return punkt, nil
	}
	punkty, err := a.repozytorium.PunktyWznowienia(ctx, przebieg.ID)
	if err != nil {
		return dane.PunktWznowienia{}, bladAutomatyki(err)
	}
	if len(punkty) == 0 {
		return dane.PunktWznowienia{}, bladNieznanegoBytuAutomatyki(dane.ErrBrakWiersza,
			"przebieg nie ma ani jednego punktu wznowienia: ", przebieg.Kod)
	}
	return punkty[len(punkty)-1], nil
}

// domknijKrokiSprzedPunktu zamyka zlecenia kroków, które punkt wznowienia
// wymienia jako ukończone, zanim kolejka ruszy dalej.
func (a *adapterAutomatyk) domknijKrokiSprzedPunktu(ctx context.Context, przebieg dane.Przebieg,
	punkt dane.PunktWznowienia) error {

	if a.kolejki == nil || przebieg.KolejkaID == nil {
		return nil
	}
	ukonczone := map[string]bool{}
	for _, kod := range kodyKrokowZZapisu(punkt.KrokiUkonczone) {
		ukonczone[kod] = true
	}
	if len(ukonczone) == 0 {
		return nil
	}
	pozycje, err := a.kolejki.repozytorium.ListaPozycji(ctx, *przebieg.KolejkaID)
	if err != nil {
		return err
	}
	for _, pozycja := range pozycje {
		if !ukonczone[pozycja.Tytul] || czyStanKoncowyPozycji(pozycja.Stan) {
			continue
		}
		if err := a.kolejki.repozytorium.ZmienStanPozycji(ctx, pozycja.ID,
			stanPozycjiUkonczona, nil); err != nil {
			return err
		}
	}
	return nil
}

// PodgladLadunku oddaje dane wejściowe i wyjściowe kroku przebiegu, zredagowane
// regułą redakcji sekretów.
func (a *adapterAutomatyk) PodgladLadunku(ctx context.Context,
	z shared.AutomationExecutionPayloadGetRequest) (shared.AutomationExecutionPayloadGetResponse, error) {

	przebieg, err := a.wierszPrzebiegu(ctx, z.ExecutionId)
	if err != nil {
		return shared.AutomationExecutionPayloadGetResponse{}, err
	}
	if z.StepId == "" {
		return shared.AutomationExecutionPayloadGetResponse{},
			bladWskazaniaAutomatyki("podgląd ładunku wymaga wskazania kroku (stepId)")
	}
	krok, err := a.repozytorium.KrokPrzebieguPoKodzie(ctx, przebieg.ID, z.StepId)
	if err != nil {
		return shared.AutomationExecutionPayloadGetResponse{},
			bladNieznanegoBytuAutomatyki(err, "krok przebiegu nie istnieje: ", z.StepId)
	}
	wejscie, poleWejscia := zredagowanyLadunek(krok.LadunekWejscia)
	wyjscie, poleWyjscia := zredagowanyLadunek(krok.LadunekWyjscia)
	odpowiedz := shared.AutomationExecutionPayloadGetResponse{Input: wejscie, Output: wyjscie}
	if zamaskowane := append(poleWejscia, poleWyjscia...); len(zamaskowane) > 0 {
		odpowiedz.RedactedFields = zamaskowane
	}
	return odpowiedz, nil
}

// OdtworzPrzebieg zakłada przebieg nowy z ładunkiem przebiegu źródłowego.
// Przebieg źródłowy zostaje nietknięty — tak mówi kontrakt.
func (a *adapterAutomatyk) OdtworzPrzebieg(ctx context.Context,
	z shared.AutomationExecutionReplayRequest) (shared.AutomationExecutionReplayResponse, error) {

	zrodlowy, err := a.wierszPrzebiegu(ctx, z.ExecutionId)
	if err != nil {
		return shared.AutomationExecutionReplayResponse{}, err
	}
	automatyka, err := a.repozytorium.AutomatykaPoID(ctx, zrodlowy.AutomatykaID)
	if err != nil {
		return shared.AutomationExecutionReplayResponse{}, bladAutomatyki(err)
	}
	// Odtworzenie rusza tą samą drogą co Operator i budzik harmonogramu:
	// nowa kolejka, kroki, start.
	if _, err := a.UruchomAutomatyke(ctx, automatyka); err != nil {
		return shared.AutomationExecutionReplayResponse{}, err
	}
	przebiegi, err := a.repozytorium.Przebiegi(ctx, automatyka.ID, 1)
	if err != nil || len(przebiegi) == 0 {
		return shared.AutomationExecutionReplayResponse{},
			bladAutomatyki(err)
	}
	odtworzony := przebiegi[0]
	a.przeniesLadunki(ctx, zrodlowy, odtworzony, wartoscTekstu(z.FromStepId), z.PayloadOverride)
	a.zapiszLog(ctx, odtworzony.ID, nil, shared.AutomationLogLevelInfo,
		"przebieg odtworzony z "+zrodlowy.Kod)
	a.zapisAudytu(ctx, &automatyka.ID, "odtworzenie przebiegu",
		map[string]any{"zrodlowy": zrodlowy.Kod, "odtworzony": odtworzony.Kod})
	return shared.AutomationExecutionReplayResponse{Execution: przebiegKontraktu(odtworzony)}, nil
}

// przeniesLadunki kopiuje ładunki kroków przebiegu źródłowego do
// odtworzonego, podstawiając ładunek na kroku startowym. Nieudane
// przeniesienie nie wywraca odtworzenia: ładunek jest nośnikiem podglądu,
// nie warunkiem wykonania.
func (a *adapterAutomatyk) przeniesLadunki(ctx context.Context, zrodlowy, odtworzony dane.Przebieg,
	odKroku string, podstawienie json.RawMessage) {

	kroki, err := a.repozytorium.KrokiPrzebiegu(ctx, zrodlowy.ID)
	if err != nil {
		return
	}
	odNapotkany := odKroku == ""
	for _, krok := range kroki {
		if krok.KrokKod == odKroku {
			odNapotkany = true
		}
		if !odNapotkany {
			continue
		}
		kopia := krok
		kopia.PrzebiegID = odtworzony.ID
		kopia.Stan = shared.AutomationExecutionStatusPending
		kopia.Zakonczono, kopia.KomunikatBledu, kopia.SladStosu = nil, nil, nil
		if len(podstawienie) > 0 && (odKroku == "" || krok.KrokKod == odKroku) {
			zapis := string(podstawienie)
			kopia.LadunekWejscia = &zapis
		}
		_ = a.repozytorium.ZapiszKrokPrzebiegu(ctx, kopia)
	}
}

// ruszKolejkePrzebiegu wykonuje działanie silnika kolejek na kolejce
// przebiegu i odnotowuje wynikowy stan przebiegu.
func (a *adapterAutomatyk) ruszKolejkePrzebiegu(ctx context.Context, przebieg dane.Przebieg,
	dzialanie shared.QueueAction) (shared.Queue, error) {

	if a.kolejki == nil {
		return shared.Queue{}, errBrakSilnikaKolejek
	}
	if przebieg.KolejkaID == nil {
		return shared.Queue{}, bladWskazaniaAutomatyki(
			"przebieg nie ma kolejki wykonującej — nie ma czego wznowić")
	}
	wynik, err := a.kolejki.Wykonaj(ctx, shared.QueueActionRequest{
		QueueId: strconv.FormatInt(*przebieg.KolejkaID, 10), Action: dzialanie,
	})
	if err != nil {
		return shared.Queue{}, err
	}
	if err := a.odnotujPrzebieg(ctx, *przebieg.KolejkaID, nil, wynik.Queue); err != nil {
		return shared.Queue{}, bladAutomatyki(err)
	}
	return wynik.Queue, nil
}

// zapiszLog nanosi wiersz dziennika przebiegu. Nieudany zapis nie wywraca
// czynności — log opisuje fakt, który już zaszedł.
func (a *adapterAutomatyk) zapiszLog(ctx context.Context, przebiegID int64, krokKod *string,
	poziom shared.AutomationLogLevel, tresc string) {

	_ = a.repozytorium.DopiszLogPrzebiegu(ctx, dane.WpisLoguPrzebiegu{
		PrzebiegID: przebiegID, KrokKod: krokKod, Poziom: poziomBazy(&poziom), Tresc: tresc,
	})
}

// zredagowanyLadunek maskuje wartości pól wrażliwych i nazywa zamaskowane pola.
// Ładunek nieczytelny jako zapis strukturalny wychodzi nietknięty — nie ma jak
// wskazać w nim pól, więc redakcja po polach nie ma się czego chwycić.
func zredagowanyLadunek(zapis *string) (json.RawMessage, []string) {
	if zapis == nil || *zapis == "" {
		return nil, nil
	}
	pola := map[string]json.RawMessage{}
	if err := json.Unmarshal([]byte(*zapis), &pola); err != nil {
		return json.RawMessage(*zapis), nil
	}
	zamaskowane := []string{}
	for nazwa := range pola {
		if !czyPoleWrazliwe(nazwa) {
			continue
		}
		pola[nazwa] = json.RawMessage(`"***"`)
		zamaskowane = append(zamaskowane, nazwa)
	}
	tresc, err := json.Marshal(pola)
	if err != nil {
		return json.RawMessage(*zapis), nil
	}
	return json.RawMessage(tresc), zamaskowane
}

// czyPoleWrazliwe rozstrzyga, czy nazwa pola wskazuje wartość wrażliwą,
// dopasowaniem fragmentu bez względu na wielkość liter.
func czyPoleWrazliwe(nazwa string) bool {
	male := strings.ToLower(nazwa)
	for _, wzorzec := range nazwyPolWrazliwych {
		if strings.Contains(male, wzorzec) {
			return true
		}
	}
	return false
}

// kodyKrokowZZapisu odczytuje wykaz kodów kroków ukończonych z zapisu
// strukturalnego punktu wznowienia w bazie.
func kodyKrokowZZapisu(zapis string) []string {
	kody := []string{}
	if zapis == "" {
		return kody
	}
	if err := json.Unmarshal([]byte(zapis), &kody); err != nil {
		return []string{}
	}
	return kody
}

// poziomBazy przekłada poziom kontraktu na słownik bazy. Brak wskazania daje
// napis pusty, czyli brak zawężenia odczytu.
func poziomBazy(poziom *shared.AutomationLogLevel) string {
	if poziom == nil {
		return ""
	}
	switch *poziom {
	case shared.AutomationLogLevelWarning:
		return "ostrzezenie"
	case shared.AutomationLogLevelError:
		return "blad"
	case shared.AutomationLogLevelInfo:
		return "informacja"
	default:
		return ""
	}
}

// poziomKontraktu przekłada słownik bazy na poziom kontraktu; wartość
// nierozpoznana daje poziom informacyjny.
func poziomKontraktu(poziom string) shared.AutomationLogLevel {
	switch poziom {
	case "ostrzezenie":
		return shared.AutomationLogLevelWarning
	case "blad":
		return shared.AutomationLogLevelError
	default:
		return shared.AutomationLogLevelInfo
	}
}
