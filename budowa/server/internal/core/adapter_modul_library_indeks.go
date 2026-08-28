// Plik obsługuje zasilanie indeksu treści modułu Library: co z bajtów pliku trafia do indeksu FTS5 i kiedy. Wyszukiwanie leży po stronie danych (`dane/biblioteka_indeks_tresci.go`); tu jest odczyt treści i rozstrzygnięcie, czy to tekst.
package core

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"unicode/utf8"

	"danacoconsole/server/internal/dane"
)

// granicaIndeksowaniaTresci obcina wyciąg indeksowany z jednego pliku. Megabajt tekstu to kilkaset stron; bez granicy plik na kilka gigabajtów wciągnąłby się do bazy w całości. To granica indeksu, nie treści: bajty pozostają na nośniku nietknięte.
const granicaIndeksowaniaTresci = 1 << 20

// zaindeksujTresc odczytuje treść spod odwołania pliku i podmienia jego wyciąg
// w indeksie. Woła się po udanym zapisie wiersza — indeks opisuje stan, który
// już jest, a nie zamiar.
func (a *adapterBiblioteki) zaindeksujTresc(ctx context.Context, plik dane.PlikBiblioteki) {
	if a == nil || a.repozytorium == nil || plik.ID == 0 {
		return
	}
	wyciag := ""
	if plik.TrescOdwolanie != nil && *plik.TrescOdwolanie != "" {
		wyciag = wyciagTekstowy(*plik.TrescOdwolanie)
	}
	// Wyciąg pusty zdejmuje wiersz z indeksu, inaczej obraz trafiałby we frazy treści poprzedniej.
	if err := a.repozytorium.ZapiszIndeksTresci(ctx, plik.ID, wyciag); err != nil {
		log.Printf("moduł Library: nie można zaindeksować treści pliku %s: %v", plik.Kod, err)
	}
}

// wyciagTekstowy oddaje początek treści pliku jako tekst do zaindeksowania, a dla treści, która tekstem nie jest, oddaje pustkę. Rozstrzyga treść, nie typ MIME, bo pole to bywa nieobecne albo niezgodne z bajtami.
func wyciagTekstowy(odwolanie string) string {
	plik, err := os.Open(odwolanie)
	if err != nil {
		return ""
	}
	defer plik.Close()

	bajty := make([]byte, granicaIndeksowaniaTresci)
	n, err := io.ReadFull(plik, bajty)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return ""
	}
	bajty = bajty[:n]
	for _, bajt := range bajty {
		if bajt == 0 {
			return ""
		}
	}
	// Obcięcie granicą mogło rozciąć ostatni znak wielobajtowy — to skutek cięcia, nie błędu UTF-8 pliku.
	bajty = bezObcietegoZnaku(bajty)
	if !utf8.Valid(bajty) {
		return ""
	}
	return string(bajty)
}

// bezObcietegoZnaku zdejmuje z końca bufora niedokończony znak UTF-8. Znak ma
// najwyżej cztery bajty, więc cofa się najwyżej o trzy — dalsze cofanie
// znaczyłoby, że treść jest niepoprawna sama z siebie, a to rozstrzyga
// `utf8.Valid` u wołającego.
func bezObcietegoZnaku(bajty []byte) []byte {
	for cofniecie := 0; cofniecie < 3 && len(bajty) > 0; cofniecie++ {
		if utf8.Valid(bajty) {
			return bajty
		}
		bajty = bajty[:len(bajty)-1]
	}
	return bajty
}
