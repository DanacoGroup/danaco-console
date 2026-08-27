package core

import "danacoconsole/shared"

// Stan pojedynczego procesu telemetrii postępu i jego przekład na ładunek kontraktu.

// procesPostepu jest stanem jednego procesu prowadzonym między zgłoszeniami punktów pracy, aż do jego zamknięcia.
type procesPostepu struct {
	Id        string
	IdOkna    string
	IdSesji   string
	Etap      int
	Etapow    int
	Stan      shared.ProgressStatus
	Nazwa     string
	Zamkniety bool
}

// nanies przenosi na proces to, co opis rzeczywiście niesie. Pole puste zostawia
// wartość zastaną, bo kolejne zgłoszenia tego samego procesu bywają uboższe
// od pierwszego.
func nanies(proces *procesPostepu, o opisProcesu) {
	if o.IdOkna != "" {
		proces.IdOkna = o.IdOkna
	}
	if o.IdSesji != "" {
		proces.IdSesji = o.IdSesji
	}
	if o.Etap > 0 {
		proces.Etap = o.Etap
	}
	if o.Etapow > 0 {
		proces.Etapow = o.Etapow
	}
}

// czyStanKoncowy mówi, czy stan zamyka proces, czyli czy dalsze zgłoszenia tego procesu nie są już oczekiwane.
func czyStanKoncowy(stan shared.ProgressStatus) bool {
	switch stan {
	case shared.ProgressStatusDone, shared.ProgressStatusFailed, shared.ProgressStatusStopped:
		return true
	default:
		return false
	}
}

// stopienUkonczenia liczy stopień ukończenia od 0 do 100.
//
// Przy nieznanej liczbie etapów stopnia nie zmyślamy: zostaje zerem do chwili
// domknięcia procesu. Etap bieżący liczy się od jedynki, więc
// etapów zakończonych jest o jeden mniej niż numer etapu.
func stopienUkonczenia(p procesPostepu) float64 {
	if p.Stan == shared.ProgressStatusDone {
		return 100
	}
	if p.Etapow <= 0 {
		return 0
	}
	zakonczone := p.Etap - 1
	if zakonczone < 0 {
		zakonczone = 0
	}
	if zakonczone > p.Etapow {
		zakonczone = p.Etapow
	}
	return float64(zakonczone) * 100 / float64(p.Etapow)
}

// ladunek składa treść zdarzenia progress.changed ze stanu procesu, obliczając przy tym stopień ukończenia.
func (p procesPostepu) ladunek() shared.ProgressChangedEvent {
	zdarzenie := shared.ProgressChangedEvent{
		ProcessId:   p.Id,
		CurrentStep: p.Etap,
		TotalSteps:  p.Etapow,
		Percent:     stopienUkonczenia(p),
		Status:      p.Stan,
	}
	if p.IdOkna != "" {
		okno := p.IdOkna
		zdarzenie.WindowId = &okno
	}
	if p.Nazwa != "" {
		nazwa := p.Nazwa
		zdarzenie.StepLabel = &nazwa
	}
	return zdarzenie
}
