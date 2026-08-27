package core

import (
	"testing"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przekład powodu zatrzymania biegu na wyliczenie kontraktu przez pole LoopState.stopReason.

// TestPowodyBieguRozrozniajaUkonczenieOdTrzechZatrzyman wykazuje, że cztery
// powody pakietu sesji dają cztery różne wartości kontraktu — żaden nie schodzi
// na cudzą i żaden nie ginie.
func TestPowodyBieguRozrozniajaUkonczenieOdTrzechZatrzyman(t *testing.T) {
	oczekiwane := map[session.PowodZatrzymania]shared.LoopStopReason{
		session.ZatrzymanieBrakPostepu: shared.LoopStopReasonNoProgress,
		session.ZatrzymanieRecznie:     shared.LoopStopReasonManual,
		session.ZatrzymanieUsterka:     shared.LoopStopReasonFailure,
		session.ZatrzymanieUkonczenie:  shared.LoopStopReasonCompleted,
	}

	wydane := map[shared.LoopStopReason]session.PowodZatrzymania{}
	for powod, oczekiwany := range oczekiwane {
		bieg := biegKontraktu(session.StanObiegu{
			IdKoordynatora: "okno-koordynatora",
			Obiegow:        1,
			Zatrzymany:     true,
			Powod:          powod,
		})
		if bieg.StopReason == nil {
			t.Fatalf("bieg zatrzymany powodem %q nie niesie powodu w odpisie kontraktu", powod)
		}
		if *bieg.StopReason != oczekiwany {
			t.Errorf("powód %q przełożył się na %q, oczekiwano %q",
				powod, *bieg.StopReason, oczekiwany)
		}
		if poprzedni, byl := wydane[*bieg.StopReason]; byl {
			t.Errorf("powody %q i %q dały tę samą wartość kontraktu %q — "+
				"rozróżnienia maszynowego nie ma", poprzedni, powod, *bieg.StopReason)
		}
		wydane[*bieg.StopReason] = powod
	}
}

// TestUkonczenieJestWartosciaKontraktu wykazuje, że `completed` nie jest nazwą
// wymyśloną przez rdzeń: wartość stoi w katalogu wyliczenia kontraktu, więc zna
// ją i klient, i narzędzia modelu.
func TestUkonczenieJestWartosciaKontraktu(t *testing.T) {
	for _, wartosc := range shared.WartosciLoopStopReason() {
		if wartosc == shared.LoopStopReasonCompleted {
			return
		}
	}
	t.Fatalf("wyliczenie LoopStopReason nie zna wartości %q; zna: %v",
		shared.LoopStopReasonCompleted, shared.WartosciLoopStopReason())
}

// TestZmianaPowoduNaUkonczenieDochodziDoKontrolki wykazuje, że przejście
// z biegu trwającego w ukończony jest zmianą WIDZIANĄ — inaczej Operator
// zobaczyłby kontrolkę stojącą na poprzednim stanie i nie dowiedziałby się
// o wyniku.
func TestZmianaPowoduNaUkonczenieDochodziDoKontrolki(t *testing.T) {
	wBiegu := biegKontraktu(session.StanObiegu{IdKoordynatora: "okno", Obiegow: 1})
	ukonczony := biegKontraktu(session.StanObiegu{
		IdKoordynatora: "okno", Obiegow: 1,
		Zatrzymany: true, Powod: session.ZatrzymanieUkonczenie,
	})
	przerwany := biegKontraktu(session.StanObiegu{
		IdKoordynatora: "okno", Obiegow: 1,
		Zatrzymany: true, Powod: session.ZatrzymanieRecznie,
	})

	if !czyBiegZmieniony(wBiegu, ukonczony) {
		t.Error("ukończenie biegu nie jest zmianą widzianą przez kontrolkę")
	}
	if !czyBiegZmieniony(przerwany, ukonczony) {
		t.Error("przejście z przerwania w ukończenie nie jest zmianą widzianą " +
			"przez kontrolkę — obie postawy biegu wyglądałyby tak samo")
	}
}
