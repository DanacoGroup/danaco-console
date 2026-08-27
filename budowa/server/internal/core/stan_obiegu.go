package core

import (
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przekład biegu naprawczego pętli koordynator-wykonawca na strukturę LoopState kontraktu.

// powodyBiegu wiążą nazwany powód zatrzymania pakietu sesji z wyliczeniem kontraktu, którego katalog wartości należy w całości do kontraktu. Wartości noProgress, manual i failure zastają zadanie przerwane, a completed skończone z wynikiem.
var powodyBiegu = map[session.PowodZatrzymania]shared.LoopStopReason{
	session.ZatrzymanieBrakPostepu: shared.LoopStopReasonNoProgress,
	session.ZatrzymanieRecznie:     shared.LoopStopReasonManual,
	session.ZatrzymanieUsterka:     shared.LoopStopReasonFailure,
	session.ZatrzymanieUkonczenie:  shared.LoopStopReasonCompleted,
}

// biegKontraktu przekłada odpis licznika obiegów na stan biegu kontraktu.
// Pola puste nie wychodzą wskaźnikiem — okno bez wykonawcy i bieg bez tury nie
// są bytami błędnymi, tylko biegiem przed pierwszym obiegiem.
func biegKontraktu(s session.StanObiegu) shared.LoopState {
	bieg := shared.LoopState{
		CoordinatorWindowId:  s.IdKoordynatora,
		Loops:                s.Obiegow,
		LoopsWithoutProgress: s.ObiegowBezPostepu,
		Threshold:            s.Prog,
		Stopped:              s.Zatrzymany,
		UpdatedAt:            s.Zaktualizowano.UnixMilli(),
	}
	if s.Zatrzymany {
		bieg.StopReason = powodBiegu(s.Powod)
	}
	if s.OstatniWykonawca != "" {
		wykonawca := s.OstatniWykonawca
		bieg.LastExecutorWindowId = &wykonawca
	}
	if s.OstatniPowodTury != "" {
		powod := s.OstatniPowodTury
		bieg.LastTurnReason = &powod
	}
	return bieg
}

// powodBiegu nazywa przyczynę zatrzymania biegu. Powód spoza katalogu schodzi
// na usterkę, aby zatrzymanie zostało widoczne w odpisie zamiast z niego
// zniknąć.
func powodBiegu(p session.PowodZatrzymania) *shared.LoopStopReason {
	powod, jest := powodyBiegu[p]
	if !jest {
		powod = shared.LoopStopReasonFailure
	}
	return &powod
}

// czyBiegZmieniony mówi, czy dwa odpisy różnią się czymś, co Operator widzi na
// kontrolce. Sam znacznik czasu zmianą nie jest — inaczej kontrolka migałaby
// przy każdym powtórzeniu tego samego stanu.
func czyBiegZmieniony(poprzedni, biezacy shared.LoopState) bool {
	return poprzedni.CoordinatorWindowId != biezacy.CoordinatorWindowId ||
		poprzedni.Loops != biezacy.Loops ||
		poprzedni.LoopsWithoutProgress != biezacy.LoopsWithoutProgress ||
		poprzedni.Threshold != biezacy.Threshold ||
		poprzedni.Stopped != biezacy.Stopped ||
		powodOdpisu(poprzedni.StopReason) != powodOdpisu(biezacy.StopReason)
}

// powodOdpisu odczytuje powód zatrzymania jako wartość porównywalną; brak powodu daje wartość pustą, gotową do porównania wskaźników.
func powodOdpisu(p *shared.LoopStopReason) shared.LoopStopReason {
	if p == nil {
		return ""
	}
	return *p
}
