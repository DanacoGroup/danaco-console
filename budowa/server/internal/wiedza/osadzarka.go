// Osadzarka jest portem źródła wektorów: jednym wąskim kształtem, którym
// wskaźnik znaczenia prosi o zamianę tekstu na liczby.
//
// Silnik osadzeń (`silnik.go`) startuje proces Pythona przez `zewnetrzne.Wolaj`
// i potrzebuje do tego trójki `session.Okno` + `session.Zasady` +
// `session.Obszar`. Trójka jest sprawą uruchamiania procesu, a nie sprawą
// wskaźnika, który pyta wyłącznie o zamianę tekstów na wektory. Port zdejmuje ze
// wskaźnika tę jedną zależność; podział na fragmenty, normalizacja, zapis do
// bazy, odczyt, iloczyn skalarny, ranking i przycięcie do `limit` zostają tym
// samym kodem, który pracuje na maszynie docelowej.
//
// Port ma w drzewie jedno wypełnienie produkcyjne — `Silnik`, co potwierdza
// `var _ Osadzarka` niżej; adapter rdzenia wiąże silnik i podaje go wskaźnikowi.
// Drugie wypełnienie żyje wyłącznie w pliku `_test.go`, więc nie da się go
// wkompilować w binarium.
//
// `Model()` należy do portu, bo nazwa modelu wchodzi do wiersza wskaźnika przy
// zapisie i zawęża odczyt przy szukaniu. Gdyby wskaźnik brał ją z ustawień,
// a wektory liczył kto inny, rozjazd zamieniłby wskaźnik w zbiór wektorów,
// o których nie wiadomo, z czym wolno je porównywać.
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
// procesu: okno, zasady izolacji i obszar dozwolony.
//
// Trójka dowiązywana jest raz, a nie podawana przy każdym zleceniu. Wskaźnik
// jest jeden na maszynę i wszystkie jego zlecenia jadą tym samym zasięgiem
// platformy; podawanie trójki przy każdym wołaniu sugerowałoby, że wolno ją
// zmieniać między fragmentami jednego dokumentu, a wtedy część wskaźnika
// powstawałaby pod inną izolacją niż reszta.
func (s *Silnik) Zwiazany(okno session.Okno, zasady session.Zasady,
	obszar session.Obszar) Osadzarka {

	return zwiazanySilnik{silnik: s, okno: okno, zasady: zasady, obszar: obszar}
}

// zwiazanySilnik jest silnikiem wraz z dowiązaną trójką uruchomienia.
type zwiazanySilnik struct {
	silnik *Silnik
	okno   session.Okno
	zasady session.Zasady
	obszar session.Obszar
}

// Model oddaje nazwę modelu silnika.
func (z zwiazanySilnik) Model() string { return z.silnik.Model() }

// Gotowy pyta silnik o gotowość dowiązaną trójką.
func (z zwiazanySilnik) Gotowy(ctx context.Context, limit time.Duration) error {
	return z.silnik.Gotowy(ctx, z.okno, z.zasady, z.obszar, limit)
}

// Osadz zleca silnikowi policzenie wektorów dowiązaną trójką.
func (z zwiazanySilnik) Osadz(ctx context.Context, teksty []string,
	limit time.Duration) ([][]float32, error) {

	return z.silnik.Osadz(ctx, z.okno, z.zasady, z.obszar, teksty, limit)
}

// Jedyne wypełnienie produkcyjne portu. Wiersz sprawia, że zmiana kształtu
// silnika przestaje się kompilować tutaj, a nie u wołającego.
var _ Osadzarka = zwiazanySilnik{}
