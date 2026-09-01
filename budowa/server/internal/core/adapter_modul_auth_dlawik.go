// Dławik prób wejścia przez bramkę: rosnąca zwłoka po próbach nieudanych, zerowana
// pierwszym wejściem udanym, ograniczająca wyłącznie prędkość zgadywania sekretu.
package core

import (
	"context"
	"sync"
	"time"
)

// zwlokaPierwszaDlawika jest zwłoką nałożoną na próbę następującą po pierwszej nieudanej,
// poniżej progu zauważalnego przez człowieka.
const zwlokaPierwszaDlawika = 250 * time.Millisecond

// zwlokaGranicznaDlawika jest sufitem, powyżej którego zwłoka nie rośnie, bo wzrost
// wykładniczy bez sufitu dałby czekanie liczone w godzinach.
const zwlokaGranicznaDlawika = 5 * time.Second

// dlawikWejscia trzyma licznik prób nieudanych osobno dla każdej drogi. Klucza drogi
// nie składa sam: podaje go dlawikDrog (kluczDlawika — czynność i konto albo połączenie),
// który też wygasza drogi bez prób; wartości z żądania kluczem nie są.
type dlawikWejscia struct {
	mu    sync.Mutex
	proby map[string]int

	// czekaj podstawia własne czekanie w miejsce zegara, żeby test mógł zmierzyć zwłokę.
	czekaj func(ctx context.Context, ile time.Duration)
}

// nowyDlawikWejscia zakłada dławik z pustym licznikiem prób nieudanych dla każdej drogi wejścia do bramki.
func nowyDlawikWejscia() *dlawikWejscia {
	return &dlawikWejscia{proby: map[string]int{}}
}

// Zaczekaj nakłada zwłokę należną drodze wejścia i zwraca jej długość; pierwsza próba nie
// czeka nigdy, a zerwanie kontekstu kończy czekanie natychmiast.
func (d *dlawikWejscia) Zaczekaj(ctx context.Context, droga string) time.Duration {
	if d == nil {
		return 0
	}
	d.mu.Lock()
	nieudane := d.proby[droga]
	d.mu.Unlock()

	zwloka := zwlokaPoProbach(nieudane)
	if zwloka <= 0 {
		return 0
	}
	if d.czekaj != nil {
		d.czekaj(ctx, zwloka)
		return zwloka
	}
	zegar := time.NewTimer(zwloka)
	defer zegar.Stop()
	select {
	case <-zegar.C:
	case <-ctx.Done():
	}
	return zwloka
}

// Niepowodzenie dolicza próbę nieudaną i zwraca zwłokę, która obejmie próbę następną na tej samej drodze wejścia.
func (d *dlawikWejscia) Niepowodzenie(droga string) time.Duration {
	if d == nil {
		return 0
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	// Licznik przestaje rosnąć tam, gdzie zwłoka i tak stoi na suficie.
	if zwlokaPoProbach(d.proby[droga]) < zwlokaGranicznaDlawika {
		d.proby[droga]++
	}
	return zwlokaPoProbach(d.proby[droga])
}

// Wyzeruj kasuje licznik drogi wejścia; woła się po wejściu udanym, gdy sekret był dobry, nie fałszywy.
func (d *dlawikWejscia) Wyzeruj(droga string) {
	if d == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.proby, droga)
}

// zwlokaPoProbach przekłada liczbę prób nieudanych na zwłokę rosnącą podwajaniem aż do sufitu wartości.
func zwlokaPoProbach(nieudane int) time.Duration {
	if nieudane <= 0 {
		return 0
	}
	zwloka := zwlokaPierwszaDlawika
	for i := 1; i < nieudane; i++ {
		zwloka *= 2
		if zwloka >= zwlokaGranicznaDlawika {
			return zwlokaGranicznaDlawika
		}
	}
	return zwloka
}
