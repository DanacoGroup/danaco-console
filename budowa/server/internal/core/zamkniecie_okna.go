package core

import "time"

// karencjaZamkniecia jest czasem, jaki tura dostaje na samodzielne domknięcie, zanim drzewo procesu zostanie ubite. Dwie sekundy pokrywają dopisanie ostatniego fragmentu i zamknięcie otwartych plików przez model.
const karencjaZamkniecia = 2 * time.Second

// krokKarencji wyznacza częstość sprawdzania, czy tura już się domknęła.
// Sprawdzanie jest tanie, więc krok jest krótki — okno ma zamknąć się od razu
// po ustaniu tury, a nie po upływie całej karencji.
const krokKarencji = 50 * time.Millisecond

// PrzerwanieTury przerywa turę okna i mówi, czy jakaś w ogóle biegła, zanim przerwanie zostało wykonane.
type PrzerwanieTury func(idOkna string) bool

// TuraWBiegu mówi, czy w oknie trwa jeszcze tura, sprawdzane po jej przerwaniu, w trakcie okresu karencji.
type TuraWBiegu func(idOkna string) bool

// ZPrzerwaniemTury wpina łagodny krok zamykania okna.
//
// Bez niego zamknięcie okna ubija drzewo procesu natychmiast, choć model może
// być w połowie zapisu pliku. Wiadomości są bezpieczne, bo idą do bazy, ale
// praca zostawiona przez model na dysku już nie.
func (a *adapterOkien) ZPrzerwaniemTury(przerwij PrzerwanieTury, wBiegu TuraWBiegu) *adapterOkien {
	a.przerwijTure = przerwij
	a.turaWBiegu = wBiegu
	return a
}

// domknijTureLagodnie przerywa turę okna i czeka, aż ustanie, nie dłużej niż karencja. Kolejność jest istotna: najpierw przerwanie tury, które daje procesowi szansę zamknąć się samemu, dopiero potem ubicie drzewa przez nadzorcę.
func (a *adapterOkien) domknijTureLagodnie(idOkna string) {
	if a.przerwijTure == nil {
		return
	}
	if !a.przerwijTure(idOkna) {
		return // tura nie biegła — nie ma czego domykać
	}
	if a.turaWBiegu == nil {
		return
	}
	granica := time.Now().Add(karencjaZamkniecia)
	for time.Now().Before(granica) {
		if !a.turaWBiegu(idOkna) {
			return
		}
		time.Sleep(krokKarencji)
	}
}
