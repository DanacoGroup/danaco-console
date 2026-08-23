// Odpowiedzialność pliku: kolumny dopuszczające NULL. Struktury repozytoriów
// używają wskaźników, baza — wartości pustych; tu leży jedyne miejsce przekładu.
package dane

import "database/sql"

// tekstZKolumny zwraca wskaźnik na napis albo nil dla kolumny pustej.
func tekstZKolumny(kolumna sql.NullString) *string {
	if !kolumna.Valid {
		return nil
	}
	wartosc := kolumna.String
	return &wartosc
}

// liczbaZKolumny zwraca wskaźnik na liczbę albo nil dla kolumny pustej.
func liczbaZKolumny(kolumna sql.NullInt64) *int64 {
	if !kolumna.Valid {
		return nil
	}
	wartosc := kolumna.Int64
	return &wartosc
}

// liczbaRzeczywistaZKolumny zwraca wskaźnik na liczbę rzeczywistą albo nil dla
// kolumny pustej — granice i skok pozycji katalogu ustawień są opcjonalne.
func liczbaRzeczywistaZKolumny(kolumna sql.NullFloat64) *float64 {
	if !kolumna.Valid {
		return nil
	}
	wartosc := kolumna.Float64
	return &wartosc
}

// wartoscLogiczna zwraca wskaźnik na wartość logiczną. Pole opcjonalne
// kontraktu niesie wtedy rozstrzygnięcie wprost, a nie brak wartości.
func wartoscLogiczna(wartosc bool) *bool {
	kopia := wartosc
	return &kopia
}

// tekstDoKolumny przekłada wskaźnik na argument zapytania; nil daje NULL.
func tekstDoKolumny(wartosc *string) any {
	if wartosc == nil {
		return nil
	}
	return *wartosc
}

// liczbaDoKolumny przekłada wskaźnik na argument zapytania; nil daje NULL.
func liczbaDoKolumny(wartosc *int64) any {
	if wartosc == nil {
		return nil
	}
	return *wartosc
}

// liczbaLogiczna przekłada wartość logiczną na kolumnę
// INTEGER NOT NULL CHECK(... IN (0,1)) — konwencja schematu Danaco Console.
func liczbaLogiczna(wartosc bool) int {
	if wartosc {
		return 1
	}
	return 0
}

// skaner obejmuje `sql.Row` i `sql.Rows` — odczyt wiersza wygląda tak samo
// niezależnie od tego, czy zapytanie zwraca jeden wiersz, czy wiele.
type skaner interface {
	Scan(cele ...any) error
}
