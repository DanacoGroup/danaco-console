// Odpowiedzialność pliku: jedno pytanie, na które musi umieć odpowiedzieć cały
// obszar Browser — czy moduł w ogóle zna wskazane okno operacyjne.
//
// Wykazy `browser.source.list` i `browser.note.list` rozróżniają dzięki temu
// pytaniu dwie sytuacje, które bez niego wyglądałyby identycznie: „okno
// przeglądania jest, tylko nic w nim jeszcze nie zebrano" (wykaz pusty,
// odpowiedź udana) oraz „takiego okna moduł nigdy nie widział" (odmowa
// `not_found`). Bez osobnego sprawdzenia obie drogi kończyłyby się tą samą
// pustą tablicą.
//
// Śladem jest każda z trzech tabel modułu, nie sama migawka. Okno staje się
// znane pierwszą nawigacją, ale równie dobrze pierwszym `browser.source.add`
// albo `browser.note.add` — żadna z tych komend nie wymaga poprzedniczki.
// Pytanie patrzące wyłącznie na `migawka_strony` odrzuciłoby jako nieznane
// okno zasilone samym źródłem.
//
// Pytanie nie sięga katalogu okien operacyjnych. Adapter modułu Browser ma
// jedną zależność — własne repozytorium — więc „znane" znaczy „znane modułowi
// Browser", nie „istniejące w tabeli okien".
package dane

import (
	"context"
	"fmt"
)

// Ślad okna szukany jest w trzech tabelach naraz jednym zapytaniem — trzy
// osobne przebiegi dałyby ten sam wynik trzema odczytami.
const sladOkna = `SELECT
	EXISTS(SELECT 1 FROM migawka_strony WHERE okno = ?)
	OR EXISTS(SELECT 1 FROM zrodlo_przegladania WHERE okno = ?)
	OR EXISTS(SELECT 1 FROM notatka_przegladania WHERE okno = ?)`

// OknoZnane mówi, czy moduł Browser zetknął się kiedykolwiek ze wskazanym
// oknem: czy zostawiono w nim migawkę, źródło albo notatkę. Okno bez nazwy
// nie jest znane nikomu — pytanie o nie nie idzie nawet do bazy.
func (r *repozytoriumPrzegladania) OknoZnane(ctx context.Context, okno string) (bool, error) {
	if okno == "" {
		return false, nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, sladOkna)
	if err != nil {
		return false, err
	}
	var znane bool
	if err := polecenie.QueryRowContext(ctx, okno, okno, okno).Scan(&znane); err != nil {
		return false, fmt.Errorf("dane: nie można sprawdzić okna przeglądania %q: %w", okno, err)
	}
	return znane, nil
}
