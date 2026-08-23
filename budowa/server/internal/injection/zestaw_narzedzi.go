// Odpowiedzialność pliku: opis zestawu narzędzi jednej tury złożony z dwóch
// źródeł — podstawy wnoszonej przez definicję eksperta i doraźnych dołożeń
// sesji — oraz przełożenie tego opisu na argumenty uruchomienia.
//
// Dlaczego opis, a nie gotowa lista nazw.
// Rdzeń nie zna nazw narzędzi eksperta i znać ich nie ma. Podzbiór wskazany
// definicją eksperta rozstrzyga serwer narzędzi (`narzedzia/ekspert_wykaz.go`),
// bo tylko on trzyma wykaz kontraktu i tylko on umie przełożyć kod eksperta na
// pozycje albo na całą grupę. Gdyby rdzeń wyliczał tę listę sam, powstałby
// drugi wykaz obok kontraktowego.
// Dlatego stąd wychodzi opis: „podstawą jest ekspert o kodzie X, dołożeniami
// są te nazwy" — a rozwinięcie opisu w listę należy do strony przeciwnej.
//
// Dlaczego składacz siedzi w tym pakiecie.
// Ten pakiet składa wiersze uruchomienia procesów tury i nie zna ani rdzenia,
// ani bazy, ani kontraktu — tak jak `nakladka.go` składa warstwy
// promptu, nie znając ani jednego zdania promptu. Zestaw narzędzi jest bytem
// tej samej rodziny: opis wchodzi z zewnątrz, wychodzi gotowa para
// przełącznik–wartość.
package injection

import "strings"

const (
	// PrzelacznikDolozen jest nazwą przełącznika niosącego doraźne dołożenia
	// sesji do serwera narzędzi modelu.
	//
	// Umowa dwóch pakietów, i to umowa ostra. Ta stała jest jedyną definicją
	// nazwy: odczyt po drugiej stronie (`cmd/danaco-narzedzia/main.go`,
	// `narzedzia.RozbijDolozenia`) bierze ją stąd, zamiast zapisywać po swojemu.
	// Powód jest twardy — wiersz uruchomienia serwera narzędzi czyta
	// `flag.Parse` na domyślnym `flag.CommandLine`, czyli z `ExitOnError`, więc
	// przełącznik nieznany nie zostaje po cichu pominięty: kończy proces serwera
	// narzędzi, a wtedy tura traci nie kilka pozycji, tylko cały wykaz. Rozjazd
	// dwóch zapisów tej samej nazwy byłby więc nie usterką kosmetyczną, tylko
	// awarią wyglądającą jak awaria modelu.
	PrzelacznikDolozen = "dolozenia"
	// RozdzielnikDolozen rozdziela nazwy w wartości przełącznika. Przecinek, bo
	// nazwa narzędzia przecinka nie zawiera nigdy — jest identyfikatorem, nie
	// zdaniem.
	RozdzielnikDolozen = ","
)

// ZestawTury mówi, z czego składa się zestaw narzędzi tej jednej tury.
//
// Pola są dwa, bo źródła są dwa i są różnej natury: podstawa jest zawężeniem
// (definicja eksperta wybiera podzbiór wykazu), dołożenia są dokładaniem
// (Operator dorzuca pozycje na czas sesji). Jedno pole nie wyraziłoby obu:
// lista pusta znaczyłaby jednocześnie „bez zawężenia" i „zero narzędzi", a te
// dwie rzeczy różnią się wszystkim — pierwsza jest turą dzisiejszą, druga
// odebraniem modelowi całego wykazu i wygląda dla Operatora jak awaria modelu,
// nie jak brak wpięcia.
type ZestawTury struct {
	// Ekspert jest kodem eksperta wnoszącego podstawę zestawu. Pusty znaczy
	// „okno bez eksperta": podstawą jest wtedy pełny wykaz w zasięgu roli okna.
	Ekspert string
	// Dolozenia są nazwami doraźnie dołożonymi w sesji, w kolejności dokładania,
	// bez duplikatów i bez pozycji pustych.
	Dolozenia []string
}

// ZlozZestawTury składa opis zestawu tury z kodu eksperta i dołożeń sesji.
//
// Kolejność dołożeń zostaje kolejnością dokładania — tą samą, którą oddaje
// warstwa danych. To ta sama reguła, którą `narzedzia/wykaz.go` stosuje do
// rozszerzenia roli okna: dołożenie dokłada, nie przestawia.
//
// Odsiew duplikatów: nazwa powtórzona zostaje na pierwszej pozycji, na której
// się pojawiła. Dołożenie po raz drugi jest czynnością pustą, nie błędem.
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

// ArgumentyDolozen zwraca parę przełącznik–wartość niosącą dołożenia sesji.
// Brak dołożeń daje nic — nie przełącznik z wartością pustą, bo ten kazałby
// stronie przeciwnej rozstrzygać, czy „nic" znaczy brak dołożeń, czy dołożenie
// o nazwie pustej.
//
// Kształt wyniku odpowiada `narzedzia.ArgumentyZasiegu` i `ArgumentyEksperta`,
// bo wszystkie trzy jadą do tego samego wpisu `danaco` i tą samą drogą.
func (z ZestawTury) ArgumentyDolozen() []string {
	wartosc := z.WartoscDolozen()
	if wartosc == "" {
		return nil
	}
	return []string{"--" + PrzelacznikDolozen, wartosc}
}
