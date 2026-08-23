// Katalog okien operacyjnych modułu na potrzeby nawigacji.
//
// Kody okien wracają z rejestru `okno_operacyjne` + `okno_operacyjne_modul`,
// a nie z pola modułu — moduł nie nosi listy swoich okien, bo ta sama definicja
// okna należy do wielu modułów (Chat Window do wszystkich piętnastu,
// Preview Window do Studio i Design). Rejestr jest wykazem informacyjnym:
// moduł spoza niego daje wykaz pusty, nie błąd.
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
