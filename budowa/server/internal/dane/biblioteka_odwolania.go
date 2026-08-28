// Plik zestawia żywe odwołania do treści biblioteki z tabel pliku i jego
// wersji jednym zapytaniem UNION, ponieważ blob porzucony przez plik bieżący
// bywa nadal treścią wersji historycznej dostępnej do przywrócenia.
package dane

import (
	"context"
	"fmt"
)

// odwolaniaTresciBiblioteki zbiera odwołania z obu tabel. UNION (nie UNION ALL)
// zdejmuje powtórzenia — plik i jego wersja wskazują ten sam blob, a wgrania
// o identycznej treści dzielą blob po sumie kontrolnej.
const odwolaniaTresciBiblioteki = `SELECT tresc_odwolanie FROM plik_biblioteki
                                   WHERE tresc_odwolanie IS NOT NULL AND tresc_odwolanie <> ''
                                   UNION
                                   SELECT tresc_odwolanie FROM wersja_pliku_biblioteki
                                   WHERE tresc_odwolanie IS NOT NULL AND tresc_odwolanie <> ''`

// OdwolaniaTresci zwraca wszystkie żywe odwołania do treści. Błąd odczytu jest
// propagowany, nie pomijany: sprzątanie z wykazem niepełnym uznałoby żywe bloby
// za porzucone i skasowało treść biblioteki.
func (r *repozytoriumBiblioteki) OdwolaniaTresci(ctx context.Context) ([]string, error) {
	wiersze, err := r.db.QueryContext(ctx, odwolaniaTresciBiblioteki)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać odwołań do treści biblioteki: %w", err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var odwolanie string
		if err := wiersze.Scan(&odwolanie); err != nil {
			return nil, fmt.Errorf("dane: nieczytelne odwołanie do treści biblioteki: %w", err)
		}
		lista = append(lista, odwolanie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt odwołań do treści biblioteki: %w", err)
	}
	return lista, nil
}
