package core

import (
	"sync"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// rejestrBiegow trzyma ostatni znany stan biegu naprawczego każdego
// koordynatora i rozgłasza jego zmiany dalej.
//
// Pętla pakietu sesji zna stan biegu wyłącznie w chwili obiegu i wydaje go
// obserwatorom jednorazowo, a o bieżący stan koordynatora pyta się z zewnątrz:
// przy wejściu na stronę główną, przy stanie okna i przy telemetrii postępu.
// Rejestr jest pamiętającym pośrednikiem między jednym a drugim. Nie prowadzi
// biegu, nie zatrzymuje go i nie liczy obiegów — to robi wyłącznie pętla.
//
// Rejestr wypełnia interfejs session.ObserwatorObiegu, więc wpina się w pętlę
// tą samą drogą, którą wpina się ślad dziennika.
type rejestrBiegow struct {
	mu           sync.RWMutex
	stany        map[string]shared.LoopState
	obserwatorzy []obserwatorBiegu
}

// obserwatorBiegu przyjmuje stan biegu po zmianie. Odbiorcą jest rejestr
// obecności sesji, który zamienia zmianę na zdarzenie kontraktu.
type obserwatorBiegu func(bieg shared.LoopState)

// nowyRejestrBiegow zakłada pusty rejestr biegów.
func nowyRejestrBiegow() *rejestrBiegow {
	return &rejestrBiegow{stany: map[string]shared.LoopState{}}
}

// Obieg wypełnia interfejs session.ObserwatorObiegu: zapisuje odpis licznika
// i — gdy odpis różni się od poprzedniego — powiadamia obserwatorów.
//
// Ślad obiegu odmówionego i obiegu zakończonego usterką dochodzi tą samą drogą,
// bo zatrzymanie biegu jest właśnie tym, co kontrolka ma pokazać.
func (r *rejestrBiegow) Obieg(z session.ZdarzenieObiegu) {
	if r == nil || z.Stan.IdKoordynatora == "" {
		return
	}
	bieg := biegKontraktu(z.Stan)

	r.mu.Lock()
	poprzedni, znany := r.stany[bieg.CoordinatorWindowId]
	r.stany[bieg.CoordinatorWindowId] = bieg
	obserwatorzy := append([]obserwatorBiegu(nil), r.obserwatorzy...)
	r.mu.Unlock()

	if znany && !czyBiegZmieniony(poprzedni, bieg) {
		return
	}
	for _, obserwator := range obserwatorzy {
		obserwator(bieg)
	}
}

// Stan zwraca ostatni znany bieg koordynatora. Okno bez biegu daje fałsz —
// okno samodzielne i wykonawcze nie prowadzą pętli, więc kontrolka nie ma dla
// nich czego rysować.
func (r *rejestrBiegow) Stan(idKoordynatora string) (shared.LoopState, bool) {
	if r == nil || idKoordynatora == "" {
		return shared.LoopState{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	bieg, jest := r.stany[idKoordynatora]
	return bieg, jest
}

// Obserwuj dokłada odbiorcę zmian stanu biegu.
func (r *rejestrBiegow) Obserwuj(o obserwatorBiegu) {
	if r == nil || o == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.obserwatorzy = append(r.obserwatorzy, o)
}

// Zachowaj zostawia w rejestrze wyłącznie biegi okien wciąż otwartych i zwraca
// liczbę śladów wykreślonych.
//
// Sprzątanie idzie od strony okien żywych, a nie od zdarzenia zamknięcia okna,
// bo pętla o zamknięciu nie mówi, a rejestr biegów nie ma prawa trzymać okna
// przy życiu dłużej niż rejestr nadzorcy. Wykaz pusty niczego nie kasuje —
// brak wiedzy o oknach nie jest wiedzą o ich zamknięciu.
func (r *rejestrBiegow) Zachowaj(okna map[string]struct{}) int {
	if r == nil || len(okna) == 0 {
		return 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	wykreslone := 0
	for idOkna := range r.stany {
		if _, zyje := okna[idOkna]; zyje {
			continue
		}
		delete(r.stany, idOkna)
		wykreslone++
	}
	return wykreslone
}
