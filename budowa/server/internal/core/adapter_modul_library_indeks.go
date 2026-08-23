// Odpowiedzialność pliku: zasilanie indeksu treści modułu Library — co z bajtów
// pliku trafia do indeksu pełnotekstowego FTS5 i kiedy. Samo wyszukiwanie leży
// po stronie danych (`dane/biblioteka_indeks_tresci.go`, `Szukaj`); tu jest
// odczyt treści spod odwołania i rozstrzygnięcie, czy to w ogóle jest tekst.
//
// Indeks zasila się przy zapisie, nie przy odczycie: skanowanie blobów przy
// każdym żądaniu wyszukiwania otwierałoby wszystkie pliki repozytorium na każde
// naciśnięcie klawisza w Library Explorerze, a koszt rósłby z rozmiarem
// biblioteki zamiast z liczbą trafień. Wpięcia są trzy i wszystkie tam, gdzie
// zmienia się treść bieżąca pliku: `Wgraj`, `DolozWersje`, `PrzywrocWersje`.
//
// Nieudane zaindeksowanie nie jest odmową komendy. W chwili indeksowania bajty
// leżą już na nośniku, a wiersz pliku w bazie; odmowa dawałaby błąd przy pliku,
// który jest wgrany i widoczny w wykazie. Brak wiersza indeksu odbiera tylko
// trafność wyszukiwania po treści — nazwa dopasowuje się dalej, bo `Szukaj`
// trzyma oba człony w alternatywie. Niepowodzenie idzie więc do dziennika
// rdzenia, a komenda kończy się powodzeniem.
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

// granicaIndeksowaniaTresci obcina wyciąg indeksowany z jednego pliku.
// Megabajt tekstu to kilkaset stron; bez tej granicy jeden plik dziennika na
// kilka gigabajtów wciągnąłby się do bazy w całości. To granica indeksu, nie
// treści: bajty pozostają na nośniku nietknięte, a podgląd i wersje czytają je
// dalej w całości.
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
	// Wyciąg pusty (treść nietekstowa albo plik bez treści) zdejmuje wiersz
	// z indeksu — inaczej plik, którego nową wersją jest obraz, trafiałby
	// nadal we frazy z treści poprzedniej, tekstowej.
	if err := a.repozytorium.ZapiszIndeksTresci(ctx, plik.ID, wyciag); err != nil {
		log.Printf("moduł Library: nie można zaindeksować treści pliku %s: %v", plik.Kod, err)
	}
}

// wyciagTekstowy oddaje początek treści pliku jako tekst do zaindeksowania,
// a dla treści, która tekstem nie jest, oddaje pustkę.
//
// Rozstrzyga treść, nie typ MIME: `mimeType` jest polem żądania i bywa go brak
// albo bywa niezgodny z bajtami, więc indeks oparty na tej deklaracji wpuściłby
// bajty obrazu jako tekst. Treść z bajtem zerowym albo z niepoprawnym UTF-8 nie
// jest tekstem i nie ma czego wnieść do wyszukiwania po słowach.
//
// Nieczytelne odwołanie oddaje pustkę bez zgłaszania błędu: sprawa czytelności
// treści należy do podglądu (`trescPodgladuBiblioteki`, odmowa wprost), a nie
// do zasilania indeksu, które komendy nie wywraca.
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
	// Obcięcie granicą mogło rozciąć ostatni znak wielobajtowy — sam ogon
	// znaku nie jest „niepoprawnym UTF-8" pliku, tylko skutkiem naszego cięcia,
	// więc odpada przed rozstrzygnięciem, czy treść jest tekstem.
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
