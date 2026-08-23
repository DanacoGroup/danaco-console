// Indeks treści modułu Library (tabela wirtualna `indeks_tresci_biblioteki`,
// FTS5): zapis wyciągu tekstowego pliku oraz zawężenie wyszukiwania do plików,
// których treść niesie szukaną frazę. Wiersz pliku, jego filtr i odczyt leżą
// w `library.go`; ten plik dokłada tam jedną klauzulę i jedno polecenie zapisu.
//
// Indeks jest bytem wtórnym wobec wiersza pliku: powstaje z treści leżącej
// w magazynie na dysku, a nie z kolumny tabeli, i rządzi się własnymi regułami
// (co się indeksuje, jak fraza Operatora staje się zapytaniem FTS5).
//
// Indeks nie jest drugą prawdą o treści. Bajty pliku leżą wyłącznie
// w magazynie treści rdzenia (`tresc_odwolanie`); tutaj leży wyciąg tekstowy
// treści bieżącej, obcięty granicą po stronie rdzenia, służący wyłącznie
// odnajdywaniu. Brak wiersza indeksu obniża trafność wyszukiwania i nic poza
// tym, dlatego zapis indeksu nie wywraca wgrania pliku.
package dane

import (
	"context"
	"fmt"
	"strings"
)

const (
	// Zapis indeksu jest podmianą wiersza pliku, nie dokładaniem: plik ma
	// jedną treść bieżącą, więc jeden wiersz indeksu. FTS5 nie zna
	// `ON CONFLICT`, stąd para „usuń po rowid, wstaw z rowid".
	usunIndeksTresciBiblioteki = `DELETE FROM indeks_tresci_biblioteki WHERE rowid = ?`

	wstawIndeksTresciBiblioteki = `INSERT INTO indeks_tresci_biblioteki (rowid, tresc)
	                               VALUES (?, ?)`

	// warunekTresciBiblioteki zawęża pliki do tych, których wyciąg treści
	// pasuje do zapytania FTS5. Podzapytanie po rowid, bo rowid indeksu jest
	// kluczem wiersza pliku.
	warunekTresciBiblioteki = `plik_biblioteki.id IN (
	                               SELECT rowid FROM indeks_tresci_biblioteki
	                               WHERE indeks_tresci_biblioteki MATCH ?)`
)

// ZapiszIndeksTresci podmienia wyciąg tekstowy pliku w indeksie treści.
// Wyciąg pusty zdejmuje plik z indeksu zamiast wpisywać pustkę: plik, którego
// treści nie da się przeczytać jako tekstu (obraz, archiwum), nie ma trafiać
// w żadną frazę.
func (r *repozytoriumBiblioteki) ZapiszIndeksTresci(ctx context.Context, plikID int64, wyciag string) error {
	if plikID == 0 {
		return fmt.Errorf("dane: indeks treści biblioteki bez pliku")
	}
	usuniecie, err := r.zapytania.przygotuj(ctx, usunIndeksTresciBiblioteki)
	if err != nil {
		return err
	}
	if _, err := usuniecie.ExecContext(ctx, plikID); err != nil {
		return fmt.Errorf("dane: nie można wyczyścić indeksu treści pliku %d: %w", plikID, err)
	}
	if wyciag == "" {
		return nil
	}
	wstawienie, err := r.zapytania.przygotuj(ctx, wstawIndeksTresciBiblioteki)
	if err != nil {
		return err
	}
	if _, err := wstawienie.ExecContext(ctx, plikID, wyciag); err != nil {
		return fmt.Errorf("dane: nie można zapisać indeksu treści pliku %d: %w", plikID, err)
	}
	return nil
}

// zapytanieTresci przekłada frazę Operatora na zapytanie FTS5.
//
// Fraza jest daną, nie składnią: w Library Explorer wpisuje się słowa, nie
// wyrażenie FTS5, a znaki `"`, `*`, `-`, `(`, `:` mają w tej składni znaczenie
// i surowa fraza z nawiasem wywracałaby wyszukiwanie błędem składni zamiast
// oddać zero trafień. Całość idzie więc jako cytowana fraza (cudzysłów
// wewnątrz podwojony), czyli wyszukiwanie sekwencji słów, nie wyrażenia
// logicznego.
//
// Gwiazdka na końcu zostawia dopasowanie przedrostkowe ostatniego słowa:
// „konfigur" trafia w „konfiguracja", więc trafienia widać przed dokończeniem
// frazy.
func zapytanieTresci(fraza string) string {
	oczyszczona := strings.TrimSpace(fraza)
	if oczyszczona == "" {
		return ""
	}
	return `"` + strings.ReplaceAll(oczyszczona, `"`, `""`) + `"*`
}
