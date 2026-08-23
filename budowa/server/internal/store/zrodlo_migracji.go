package store

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
)

// Numeracja migracji — luki są z zamysłu, nie z ubytku.
//
// Wykaz `migracja_*.sql` w tym katalogu niesie mniej plików, niż wskazuje
// najwyższy numer. Wolne numery nie oznaczają kroku usuniętego ani zgubionego:
// powstają przy pracy równoległej, gdy numer rezerwowany z góry dla kroku, który
// ostatecznie nie wszedł, zostaje pusty. Numeru zwolnionego nie wolno użyć
// powtórnie: bazy założone wcześniej mają już wyższą wersję schematu i krok
// wstawiony w lukę nigdy by się na nich nie wykonał.
//
// Co jest naprawdę wymagane od numeracji — i czego pilnuje kod poniżej:
//   - Jednoznaczność. Dwa pliki o tym samym numerze to awaria startu, bo
//     rejestr `migracja` ma na kolumnie `wersja` warunek UNIQUE: drugi krok
//     wykonałby swój schemat, ale nie zostałby odnotowany.
//   - Porządek rosnący. Kroki stosuje się po numerze rosnąco (sortowanie
//     w wczytajMigracje), więc kolejność w katalogu nie ma znaczenia.
//   - Niezmienność treści. Suma kontrolna kroku już zastosowanego musi się
//     zgadzać (Migruj w migracje.go).
//
// Ciągłość numeracji nie jest wymagana.
//
// plikiMigracji — schemat wkompilowany w binarium, żeby wdrożenie nie zależało
// od obecności plików obok programu.
//
//go:embed migracja_*.sql
var plikiMigracji embed.FS

const (
	przedrostekMigracji  = "migracja_"
	rozszerzenieMigracji = ".sql"
)

// migracja to pojedynczy krok schematu odczytany z zasobów pakietu.
type migracja struct {
	Wersja        int
	Nazwa         string
	Tresc         string
	SumaKontrolna string
}

// wczytajMigracje odczytuje wszystkie kroki schematu i porządkuje je rosnąco
// po numerze wersji.
func wczytajMigracje() ([]migracja, error) {
	wpisy, err := fs.ReadDir(plikiMigracji, ".")
	if err != nil {
		return nil, fmt.Errorf("store: nie można odczytać zasobów migracji: %w", err)
	}
	kroki := make([]migracja, 0, len(wpisy))
	for _, wpis := range wpisy {
		if wpis.IsDir() || !strings.HasSuffix(wpis.Name(), rozszerzenieMigracji) {
			continue
		}
		krok, err := zbudujKrok(wpis.Name())
		if err != nil {
			return nil, err
		}
		kroki = append(kroki, krok)
	}
	sort.Slice(kroki, func(i, j int) bool { return kroki[i].Wersja < kroki[j].Wersja })
	if err := sprawdzUnikalnoscWersji(kroki); err != nil {
		return nil, err
	}
	if len(kroki) == 0 {
		return nil, fmt.Errorf("store: brak plików migracji w zasobach pakietu")
	}
	return kroki, nil
}

// zbudujKrok wyprowadza wersję i nazwę z nazwy pliku `migracja_NNN_nazwa.sql`
// oraz liczy sumę kontrolną treści.
func zbudujKrok(nazwaPliku string) (migracja, error) {
	rdzen := strings.TrimSuffix(strings.TrimPrefix(nazwaPliku, przedrostekMigracji), rozszerzenieMigracji)
	czesci := strings.SplitN(rdzen, "_", 2)
	if !strings.HasPrefix(nazwaPliku, przedrostekMigracji) || len(czesci) != 2 {
		return migracja{}, fmt.Errorf("store: nazwa migracji %q nie ma postaci migracja_NNN_nazwa.sql", nazwaPliku)
	}
	wersja, err := strconv.Atoi(czesci[0])
	if err != nil {
		return migracja{}, fmt.Errorf("store: nieczytelny numer wersji w %q: %w", nazwaPliku, err)
	}
	tresc, err := plikiMigracji.ReadFile(nazwaPliku)
	if err != nil {
		return migracja{}, fmt.Errorf("store: nie można odczytać %q: %w", nazwaPliku, err)
	}
	suma := sha256.Sum256(tresc)
	return migracja{
		Wersja:        wersja,
		Nazwa:         czesci[1],
		Tresc:         string(tresc),
		SumaKontrolna: hex.EncodeToString(suma[:]),
	}, nil
}

// sprawdzUnikalnoscWersji nie dopuszcza dwóch kroków o tym samym numerze —
// kolejność zastosowania byłaby wtedy nieokreślona.
func sprawdzUnikalnoscWersji(kroki []migracja) error {
	for i := 1; i < len(kroki); i++ {
		if kroki[i].Wersja == kroki[i-1].Wersja {
			return fmt.Errorf("store: zduplikowana wersja migracji %03d (%q i %q)",
				kroki[i].Wersja, kroki[i-1].Nazwa, kroki[i].Nazwa)
		}
	}
	return nil
}
