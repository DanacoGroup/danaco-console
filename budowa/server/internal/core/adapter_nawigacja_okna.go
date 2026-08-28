// Katalog okien operacyjnych modułu dostarcza kody okien z rejestru
// okno_operacyjne i okno_operacyjne_modul, ponieważ pojedynczy moduł nie
// przechowuje własnej listy okien, a moduł spoza rejestru zwraca listę pustą
// zamiast błędu.
package core

import (
	"context"
)

// kodyOkienModulu zwraca kody okien operacyjnych otwieranych wraz z modułem,
// w kolejności ich pozycji w tym module. Okno rozmowy jest przypięte do każdego
// modułu, więc wykaz pusty oznacza moduł spoza rejestru, nie moduł bez okien.
func (a *adapterNawigacji) kodyOkienModulu(ctx context.Context, modulID int64) ([]string, error) {
	okna, err := a.zestaw.OknaOperacyjne.ListaModulu(ctx, modulID)
	if err != nil {
		return nil, err
	}
	kody := make([]string, 0, len(okna))
	for _, okno := range okna {
		kody = append(kody, okno.Kod)
	}
	return kody, nil
}
