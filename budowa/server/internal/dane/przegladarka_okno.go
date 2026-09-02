// Plik sprawdza, czy moduł przeglądania stron zna wskazane okno operacyjne,
// przez jedno zapytanie o ślad w trzech tabelach: migawek, źródeł i notatek.
package dane

import (
	"context"
	"fmt"
)

// Ślad okna szukany jest w trzech tabelach naraz jednym zapytaniem — trzy
// osobne przebiegi dałyby ten sam wynik trzema odczytami.
const sladOkna = `SELECT
	EXISTS(SELECT 1 FROM migawka_strony WHERE okno = ? AND ` + WarunekKonta + `)
	OR EXISTS(SELECT 1 FROM zrodlo_przegladania WHERE okno = ? AND ` + WarunekKonta + `)
	OR EXISTS(SELECT 1 FROM notatka_przegladania WHERE okno = ? AND ` + WarunekKonta + `)`

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
	konto := KontoOperatora(ctx)
	if err := polecenie.QueryRowContext(ctx, okno, konto, okno, konto, okno, konto).Scan(&znane); err != nil {
		return false, fmt.Errorf("dane: nie można sprawdzić okna przeglądania %q: %w", okno, err)
	}
	return znane, nil
}
