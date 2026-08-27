// Odpowiedzialność pliku: sprzątanie magazynu treści modułu Library —
// usunięcie z nośnika blobów bez żywego odwołania oraz plików tymczasowych po
// zapisach przerwanych w połowie.
package core

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
)

// karencjaPlikuCzesciowego chroni plik tymczasowy zapisu, który może trwać
// w tej właśnie chwili — drugi rdzeń Operatora na tym samym katalogu danych
// może akurat wciągać wielki plik. Doba jest granicą z ogromnym zapasem.
const karencjaPlikuCzesciowego = 24 * time.Hour

// przyrostekPlikuCzesciowego znakuje pliki tymczasowe magazynu
// (`tresc-*.czesciowa`, `wciagana-*.czesciowa` — `Zapisz` i `ZapiszZePliku`).
const przyrostekPlikuCzesciowego = ".czesciowa"

// posprzatajMagazynTresci kasuje bloby bez żywego odwołania oraz przeterminowane
// pliki tymczasowe. Woła się raz, przy montażu portu `ZKatalogiemDanych`.
func posprzatajMagazynTresci(ctx context.Context, repozytorium dane.RepozytoriumBiblioteki,
	magazyn *magazynTresciBiblioteki) {

	if repozytorium == nil || magazyn == nil || magazyn.katalog == "" {
		return
	}
	// Katalog jeszcze nieistniejący znaczy magazyn pusty — nie ma czego
	// sprzątać i nie ma o czym meldować.
	if _, err := os.Stat(magazyn.katalog); err != nil {
		return
	}
	odwolania, err := repozytorium.OdwolaniaTresci(ctx)
	if err != nil {
		log.Printf("moduł Library: sprzątanie magazynu treści wstrzymane — nie można odczytać "+
			"żywych odwołań: %v", err)
		return
	}

	zywe := make(map[string]struct{}, len(odwolania))
	for _, odwolanie := range odwolania {
		// Klucz jest ścieżką oczyszczoną, nie surowym napisem, inaczej porównanie
		// zgubi żywy blob.
		zywe[filepath.Clean(odwolanie)] = struct{}{}
	}

	usuniete, odzyskane := przemiecKatalogMagazynu(magazyn.katalog, zywe)
	if usuniete > 0 {
		log.Printf("moduł Library: magazyn treści — usunięto %d porzuconych blobów (%d bajtów)",
			usuniete, odzyskane)
	}
}

// przemiecKatalogMagazynu obchodzi magazyn i kasuje to, czego nikt nie trzyma.
// Oddaje liczbę usuniętych plików i odzyskane bajty — meldunek ma mówić, co
// naprawdę zaszło, a nie że „sprzątnięto".
func przemiecKatalogMagazynu(katalog string, zywe map[string]struct{}) (int, int64) {
	usuniete, odzyskane := 0, int64(0)
	granicaCzesciowych := time.Now().Add(-karencjaPlikuCzesciowego)

	// Błąd obejścia pojedynczego wpisu jest pominięciem, obchód idzie dalej.
	_ = filepath.WalkDir(katalog, func(sciezka string, wpis os.DirEntry, err error) error {
		if err != nil || wpis == nil || wpis.IsDir() {
			return nil
		}
		stan, err := wpis.Info()
		if err != nil {
			return nil
		}
		if strings.HasSuffix(wpis.Name(), przyrostekPlikuCzesciowego) {
			// Plik tymczasowy nigdy nie jest odwołaniem w bazie, rozstrzyga sam wiek.
			if stan.ModTime().Before(granicaCzesciowych) && os.Remove(sciezka) == nil {
				usuniete++
				odzyskane += stan.Size()
			}
			return nil
		}
		if _, jest := zywe[filepath.Clean(sciezka)]; jest {
			return nil
		}
		if err := os.Remove(sciezka); err != nil {
			log.Printf("moduł Library: nie można usunąć porzuconego bloba %s: %v", sciezka, err)
			return nil
		}
		usuniete++
		odzyskane += stan.Size()
		return nil
	})
	return usuniete, odzyskane
}
