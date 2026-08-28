// Plik opisuje zestaw narzędzi jednej tury, złożony z podstawy wnoszonej
// przez definicję eksperta i doraźnych dołożeń sesji, oraz przekłada ten
// opis na argumenty uruchomienia serwera narzędzi.
package injection

import "strings"

const (
	// PrzelacznikDolozen jest nazwą przełącznika niosącego doraźne dołożenia
	// sesji do serwera narzędzi modelu; definicja tej nazwy jest wspólna
	// z modułem cmd/danaco-narzedzia.
	PrzelacznikDolozen = "dolozenia"
	// RozdzielnikDolozen rozdziela nazwy w wartości przełącznika. Przecinek, bo
	// nazwa narzędzia przecinka nie zawiera nigdy — jest identyfikatorem, nie
	// zdaniem.
	RozdzielnikDolozen = ","
)

// ZestawTury opisuje zestaw narzędzi jednej tury, złożony z podstawy
// wskazanej definicją eksperta oraz dołożeń dorzuconych w sesji.
type ZestawTury struct {
	// Ekspert jest kodem eksperta wnoszącego podstawę zestawu; pusty oznacza
	// pełny wykaz roli okna.
	Ekspert string
	// Dolozenia są nazwami dołożonymi w sesji, w kolejności dokładania, bez
	// duplikatów i pozycji pustych.
	Dolozenia []string
}

// ZlozZestawTury składa opis zestawu tury z kodu eksperta i dołożeń sesji,
// zachowując kolejność dokładania i usuwając powtórzone nazwy.
func ZlozZestawTury(kodEksperta string, dolozenia []string) ZestawTury {
	zestaw := ZestawTury{Ekspert: strings.TrimSpace(kodEksperta)}
	widziane := make(map[string]bool, len(dolozenia))
	for _, nazwa := range dolozenia {
		nazwa = strings.TrimSpace(nazwa)
		if nazwa == "" || widziane[nazwa] {
			continue
		}
		widziane[nazwa] = true
		zestaw.Dolozenia = append(zestaw.Dolozenia, nazwa)
	}
	return zestaw
}

// Skladany mówi, czy zestaw w ogóle różni się od zestawu dotychczasowego.
// Fałsz znaczy turę bez składania: pełny wykaz w zasięgu roli okna.
func (z ZestawTury) Skladany() bool {
	return z.Ekspert != "" || len(z.Dolozenia) > 0
}

// WartoscDolozen zapisuje dołożenia jednym napisem. Brak dołożeń daje napis
// pusty, a napis pusty nie dokłada przełącznika.
func (z ZestawTury) WartoscDolozen() string {
	return strings.Join(z.Dolozenia, RozdzielnikDolozen)
}

// ArgumentyDolozen zwraca parę przełącznik-wartość niosącą dołożenia sesji;
// brak dołożeń zwraca pustą listę zamiast przełącznika z wartością pustą.
func (z ZestawTury) ArgumentyDolozen() []string {
	wartosc := z.WartoscDolozen()
	if wartosc == "" {
		return nil
	}
	return []string{"--" + PrzelacznikDolozen, wartosc}
}
