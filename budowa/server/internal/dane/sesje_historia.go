package dane

import (
	"context"
	"fmt"
)

// Polecenia historii sesji. Osobno od sesje.go, bo dotyczą czynności Operatora
// na wykazie sesji, a nie cyklu życia sesji w pracy bieżącej.
const (
	zmienTytulSesji = `UPDATE sesja
	                   SET tytul = ?,
	                       zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                   WHERE id = ?`

	zmienProjektSesji = `UPDATE sesja
	                     SET projekt = ?,
	                         zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                     WHERE id = ?`
)

// ZmienTytul zmienia nazwę sesji w historii.
//
// Tytuł jest kolumną obowiązkową, więc nazwa pusta jest odrzucana: sesja bez
// nazwy zniknęłaby z wykazu jako wiersz bez etykiety.
func (r *repozytoriumSesji) ZmienTytul(ctx context.Context, id int64, tytul string) error {
	if tytul == "" {
		return fmt.Errorf("dane: nazwa sesji %d nie może być pusta", id)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienTytulSesji)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, tytul, id)
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić nazwy sesji %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "sesja", id)
}

// ZmienProjekt wiąże sesję z projektem albo wyjmuje ją z projektu.
//
// Wskazanie puste zapisuje NULL — sesja bez projektu jest stanem poprawnym,
// a nie brakiem do uzupełnienia. Kolumna dopuszcza NULL z tego samego powodu.
func (r *repozytoriumSesji) ZmienProjekt(ctx context.Context, id int64, projekt string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zmienProjektSesji)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, tekstDoKolumny(wskaznikTekstuLubNil(projekt)), id)
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić projektu sesji %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "sesja", id)
}

// wskaznikTekstuLubNil zamienia napis pusty na brak wartości, ponieważ
// kolumna dopuszcza NULL zamiast pustego łańcucha znaków.
func wskaznikTekstuLubNil(wartosc string) *string {
	if wartosc == "" {
		return nil
	}
	return &wartosc
}
