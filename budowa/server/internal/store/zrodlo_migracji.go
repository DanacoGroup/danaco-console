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

// plikiMigracji zawiera schemat wkompilowany w binarium, żeby wdrożenie nie zależało od obecności plików obok programu; numeracja migracji może mieć luki z zamysłu.
//
//go:embed migracja_*.sql
var plikiMigracji embed.FS

const (
	przedrostekMigracji  = "migracja_"
	rozszerzenieMigracji = ".sql"
)

// migracja to pojedynczy krok schematu bazy danych odczytany z osadzonych zasobów tego pakietu aplikacji.
type migracja struct {
	Wersja        int
	Nazwa         string
	Tresc         string
	SumaKontrolna string
}

// Funkcja wczytajMigracje odczytuje wszystkie kroki schematu i porządkuje je rosnąco po numerze ich wersji.
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
// oraz liczy sumę kontrolną znormalizowanej treści kroku.
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
	suma := sha256.Sum256(normalizujTrescMigracji(tresc))
	return migracja{
		Wersja:        wersja,
		Nazwa:         czesci[1],
		Tresc:         string(tresc),
		SumaKontrolna: hex.EncodeToString(suma[:]),
	}, nil
}

// normalizujTrescMigracji sprowadza treść kroku do postaci mierzonej sumą kontrolną: poza literałami znakowymi usuwa komentarze `--`, końcowe białe znaki wiersza i wiersze puste, przez co redakcja komentarzy nie unieważnia kroków już zastosowanych.
func normalizujTrescMigracji(tresc []byte) []byte {
	var oczyszczona strings.Builder
	wLiterale := false
	for i := 0; i < len(tresc); i++ {
		znak := tresc[i]
		if wLiterale {
			if znak == '\'' {
				if i+1 < len(tresc) && tresc[i+1] == '\'' {
					oczyszczona.WriteByte(znak)
					oczyszczona.WriteByte(tresc[i+1])
					i++
					continue
				}
				wLiterale = false
			}
			oczyszczona.WriteByte(znak)
			continue
		}
		if znak == '\'' {
			wLiterale = true
			oczyszczona.WriteByte(znak)
			continue
		}
		if znak == '-' && i+1 < len(tresc) && tresc[i+1] == '-' {
			for i < len(tresc) && tresc[i] != '\n' {
				i++
			}
			if i < len(tresc) {
				oczyszczona.WriteByte('\n')
			}
			continue
		}
		oczyszczona.WriteByte(znak)
	}
	wiersze := strings.Split(oczyszczona.String(), "\n")
	zebrane := make([]string, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wiersz = strings.TrimRight(wiersz, " \t\r")
		if wiersz == "" {
			continue
		}
		zebrane = append(zebrane, wiersz)
	}
	return []byte(strings.Join(zebrane, "\n"))
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
