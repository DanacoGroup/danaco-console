// Odpowiedzialność pliku: sprzątanie magazynu treści modułu Library — usunięcie
// z nośnika blobów, na które nie wskazuje już ani plik, ani żadna jego wersja,
// oraz plików tymczasowych po zapisach przerwanych w połowie.
//
// Bez sprzątania magazyn treści (`adapter_modul_library_magazyn.go`) tylko
// przyrasta: wiersz pliku znika kaskadą `ON DELETE CASCADE` albo czyszczeniem
// stanu trwałego, wersja zostaje zastąpiona, zapis przerywa się na `os.Rename`
// i zostawia `tresc-*.czesciowa` — a bajty leżą dalej. To magazyn, a nie baza,
// wypełniałby wtedy dysk Operatora.
//
// Sprzątanie idzie przemiataniem przy starcie, a nie zliczaniem odwołań.
// Zliczanie wymaga miejsca, w którym się kasuje, a moduł Library takiego miejsca
// nie ma: wśród dziesięciu jego komend (`adapter_modul_library_uchwyty.go`) nie
// ma ani jednej usuwającej plik czy wersję. Bloby osierocają się więc wyłącznie
// drogami, których ten moduł nie kontroluje: kaskadą z wiersza usuniętego gdzie
// indziej, awarią zapisu, zastąpieniem treści bieżącej. Licznik wpięty
// w komendę, która nie istnieje, nie zliczyłby niczego.
//
// Do tego blob jest dzielony: jego nazwą jest suma sha256 treści, więc dwa
// wgrania tej samej zawartości wskazują jeden plik na nośniku, a plik i jego
// wersje wskazują go po wielekroć. Licznik musiałby być drugą, osobno
// utrzymywaną prawdą o tym, co baza i tak wie z kolumn `tresc_odwolanie`,
// i rozjechałby się przy pierwszym zapisie przerwanym w połowie.
//
// Przemiatanie startowe pyta o to samo wprost: żywe odwołania czyta z bazy
// (`dane.OdwolaniaTresci`, jedno zapytanie po obu tabelach), po czym obchodzi
// katalog magazynu i kasuje to, czego w tym wykazie nie ma. Nie potrzebuje
// budzika ani wątku — start jest jedynym momentem, w którym rdzeń i tak czyta
// stan trwały i w którym nikt równolegle nie wgrywa pliku; tak samo sprzątają
// kosz sesji i retencja historii (`trwalosc_kosza.go`). Za bezczynności rdzenia
// sprzątanie ma zero kosztu, a po awarii zapisu naprawia stan samo, bez komendy
// naprawczej.
//
// Cena jest jedna: blob osierocony między dwoma startami przeżyje do
// następnego. Miejsce na nośniku odzyskuje się z opóźnieniem, treść natomiast
// nie ginie przedwcześnie, bo o życiu bloba rozstrzyga baza, a nie licznik.
//
// Nieudany odczyt wykazu wstrzymuje sprzątanie w całości. Wykaz niepełny
// znaczyłby, że żywe bloby wyglądają na porzucone — sprzątanie skasowałoby
// wtedy treść biblioteki nie do odzyskania. Brak wiedzy nie jest wiedzą
// o braku.
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
// w tej właśnie chwili. Sprzątanie idzie przy starcie, więc własnych zapisów
// rdzeń jeszcze nie prowadzi — ale drugi rdzeń Operatora, wystartowany na tym
// samym katalogu danych, może akurat wciągać wielki plik. Doba jest granicą
// z ogromnym zapasem: zapis trwający dłużej niż dobę już się nie skończy.
const karencjaPlikuCzesciowego = 24 * time.Hour

// przyrostekPlikuCzesciowego znakuje pliki tymczasowe magazynu
// (`tresc-*.czesciowa`, `wciagana-*.czesciowa` — `Zapisz` i `ZapiszZePliku`).
const przyrostekPlikuCzesciowego = ".czesciowa"

// posprzatajMagazynTresci kasuje bloby bez żywego odwołania oraz przeterminowane
// pliki tymczasowe. Woła się raz, przy montażu portu — patrz `ZKatalogiemDanych`.
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
		// Klucz jest ścieżką oczyszczoną, nie surowym napisem z bazy: to samo
		// miejsce na dysku bywa zapisane inaczej (ukośniki, `.`), a porównanie
		// napisów uznałoby wtedy żywy blob za porzucony.
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

	// Błąd obejścia pojedynczego wpisu jest pominięciem, nie przerwaniem:
	// jeden nieczytelny katalog nie ma prawa zostawić całego magazynu
	// nieposprzątanego. Zwracamy `nil`, żeby obchód szedł dalej.
	_ = filepath.WalkDir(katalog, func(sciezka string, wpis os.DirEntry, err error) error {
		if err != nil || wpis == nil || wpis.IsDir() {
			return nil
		}
		stan, err := wpis.Info()
		if err != nil {
			return nil
		}
		if strings.HasSuffix(wpis.Name(), przyrostekPlikuCzesciowego) {
			// Plik tymczasowy nigdy nie jest odwołaniem w bazie (odwołaniem
			// staje się dopiero nazwa po przemianowaniu), więc rozstrzyga
			// o nim sam wiek.
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
