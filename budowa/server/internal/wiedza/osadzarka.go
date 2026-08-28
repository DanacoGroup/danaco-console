// Osadzarka jest portem źródła wektorów: jednym wąskim kształtem, którym
// wskaźnik znaczenia prosi o zamianę tekstu na liczby, bez zależności od
// uruchamiania procesu.
package wiedza

import (
	"context"
	"time"

	"danacoconsole/server/internal/session"
)

// Osadzarka zamienia teksty na wektory tym samym modelem, którego nazwę podaje.
//
// Kolejność wektorów odpowiada kolejności tekstów: wołający wiąże wektor
// z fragmentem pozycją, a nie treścią.
type Osadzarka interface {
	// Model oddaje nazwę modelu, którym liczone są wektory.
	Model() string
	// Gotowy sprawdza, czy jest czym liczyć, bez liczenia czegokolwiek.
	Gotowy(ctx context.Context, limit time.Duration) error
	// Osadz oddaje po jednym wektorze na tekst, w tej samej kolejności.
	Osadz(ctx context.Context, teksty []string, limit time.Duration) ([][]float32, error)
}

// Zwiazany oddaje silnik jako Osadzarkę, dowiązując mu trójkę uruchomienia
// procesu: okno, zasady izolacji i obszar dozwolony. Trójka dowiązywana jest
// raz, a nie podawana przy każdym zleceniu.
func (s *Silnik) Zwiazany(okno session.Okno, zasady session.Zasady,
	obszar session.Obszar) Osadzarka {

	return zwiazanySilnik{silnik: s, okno: okno, zasady: zasady, obszar: obszar}
}

// zwiazanySilnik jest silnikiem wraz z dowiązaną trójką uruchomienia okna,
// zasad i obszaru, wypełniając port Osadzarka.
type zwiazanySilnik struct {
	silnik *Silnik
	okno   session.Okno
	zasady session.Zasady
	obszar session.Obszar
}

// Model oddaje nazwę modelu silnika, tego samego, którym liczone są wektory
// przy budowaniu wskaźnika i przy zapytaniu.
func (z zwiazanySilnik) Model() string { return z.silnik.Model() }

// Gotowy pyta silnik o gotowość, dowiązaną trójką uruchomienia procesu,
// przekazaną portowi przy założeniu.
func (z zwiazanySilnik) Gotowy(ctx context.Context, limit time.Duration) error {
	return z.silnik.Gotowy(ctx, z.okno, z.zasady, z.obszar, limit)
}

// Osadz zleca silnikowi policzenie wektorów, dowiązaną trójką uruchomienia
// procesu, przekazaną portowi przy założeniu.
func (z zwiazanySilnik) Osadz(ctx context.Context, teksty []string,
	limit time.Duration) ([][]float32, error) {

	return z.silnik.Osadz(ctx, z.okno, z.zasady, z.obszar, teksty, limit)
}

// Jedyne wypełnienie produkcyjne portu. Wiersz sprawia, że zmiana kształtu
// silnika przestaje się kompilować tutaj, a nie u wołającego.
var _ Osadzarka = zwiazanySilnik{}
