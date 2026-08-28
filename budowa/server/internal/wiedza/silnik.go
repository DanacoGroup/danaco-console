// Odpowiedzialność pliku: złożenie silnika osadzeń w jeden byt i przeprowadzenie
// jednego zlecenia — od tekstów do wektorów. Model jest jednym polem jednego
// bytu, bo osadzenia wskaźnika i zapytania musi liczyć ten sam model.
package wiedza

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
)

const (
	// LimitBudowania — granica czasu jednej partii osadzeń przy budowaniu
	// wskaźnika, obejmująca pobranie wag przy pierwszym uruchomieniu procesu.
	LimitBudowania = 20 * time.Minute
	// LimitZapytania — granica czasu osadzenia jednego pytania, krótsza niż
	// budowania, bo wagi leżą już na dysku po pierwszym uruchomieniu.
	LimitZapytania = 2 * time.Minute
	// wielkoscPartii — ile fragmentów idzie do pomocnika za jednym razem,
	// dobrane duże, bo start procesu i wczytanie wag kosztuje kilka sekund.
	wielkoscPartii = 256
)

// Silnik liczy osadzenia pomocnikiem lokalnym. Ustawienia trzymane są
// w silniku, a nie odczytywane przy każdym zleceniu: wołający składa je
// z konfiguracji zasięgu i podaje gotowe.
type Silnik struct {
	// uruchamiacz — jedyna droga startu procesu w drzewie.
	uruchamiacz session.Uruchamiacz
	// katalogDanych — ten sam katalog, w którym leżą baza i magazyn biblioteki.
	katalogDanych string
	// ustawienia — komplet nastaw obowiązujący dla zleceń tego silnika.
	ustawienia Ustawienia
}

// NowySilnik zakłada silnik na uruchamiaczu procesów i katalogu danych,
// z domyślnym modelem i domyślnymi limitami czasu, gotowy do użycia.
func NowySilnik(uruchamiacz session.Uruchamiacz, katalogDanych string) *Silnik {
	return &Silnik{
		uruchamiacz:   uruchamiacz,
		katalogDanych: katalogDanych,
		ustawienia:    UstawieniaDomyslne(),
	}
}

// ZUstawieniami oddaje silnikowi komplet nastaw złożony z konfiguracji,
// zastępując wartości domyślne modelu i limitów czasu.
func (s *Silnik) ZUstawieniami(u Ustawienia) *Silnik {
	s.ustawienia = u
	return s
}

// Ustawienia oddaje nastawy, którymi silnik dziś pracuje: model i limity
// czasu, tak jak zostały mu ustawione przy montażu albo domyślnie założone.
func (s *Silnik) Ustawienia() Ustawienia {
	return s.ustawienia
}

// Model oddaje nazwę modelu, którym liczone są wektory. Wchodzi do odpowiedzi
// `knowledge.index` i do wiersza wskaźnika — wektor bez nazwy modelu jest
// liczbami, o których nie wiadomo, z czym wolno je porównywać.
func (s *Silnik) Model() string {
	return s.ustawienia.Model
}

// zlecenie i odpowiedz to kształt rozmowy z pomocnikiem. Pola odpowiadają co do
// znaku kluczom w `pomocnik_osadzen.py` — rozjazd zamieniłby nazwany brak
// w brak nierozpoznany.
type zlecenie struct {
	Model         string   `json:"model"`
	KatalogModeli string   `json:"katalogModeli"`
	Teksty        []string `json:"teksty"`
	WagaMb        int      `json:"wagaMb"`
}

type odpowiedz struct {
	Ok      bool        `json:"ok"`
	Model   string      `json:"model"`
	Wymiar  int         `json:"wymiar"`
	Wektory [][]float32 `json:"wektory"`
	Brak    string      `json:"brak"`
	Powod   string      `json:"powod"`
	WagaMb  int         `json:"wagaMb"`
}

// Gotowy sprawdza, czy silnik ma czym liczyć, nie licząc niczego naprawdę.
// Kosztuje jedno uruchomienie pomocnika z pustym wykazem tekstów.
func (s *Silnik) Gotowy(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar, limit time.Duration) error {

	_, err := s.wolaj(ctx, okno, zasady, obszar, nil, limit)
	return err
}

// Osadz zamienia teksty na wektory, zachowując ich kolejność. Idzie
// partiami, nie wszystko naraz. Wykaz pusty oddaje wykaz pusty bez
// uruchamiania procesu.
func (s *Silnik) Osadz(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	teksty []string, limit time.Duration) ([][]float32, error) {

	if len(teksty) == 0 {
		return nil, nil
	}
	wektory := make([][]float32, 0, len(teksty))
	for poczatek := 0; poczatek < len(teksty); poczatek += wielkoscPartii {
		koniec := poczatek + wielkoscPartii
		if koniec > len(teksty) {
			koniec = len(teksty)
		}
		partia, err := s.wolaj(ctx, okno, zasady, obszar, teksty[poczatek:koniec], limit)
		if err != nil {
			return nil, err
		}
		if len(partia) != koniec-poczatek {
			return nil, errors.New("wskaźnik znaczenia: pomocnik oddał " +
				liczba(len(partia)) + " wektorów na " + liczba(koniec-poczatek) +
				" tekstów — wiązanie po pozycji przestałoby cokolwiek znaczyć; " +
				"naprawa: zgłosić usterkę pomocnika osadzeń")
		}
		wektory = append(wektory, partia...)
	}
	return wektory, nil
}

// wolaj przeprowadza jedno uruchomienie pomocnika i czyta jego odpowiedź,
// z granicą czasu wskazaną wołaniem, nigdy bez niej.
func (s *Silnik) wolaj(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	teksty []string, limit time.Duration) ([][]float32, error) {

	skrypt, err := wylozSkrypt(s.katalogDanych)
	if err != nil {
		return nil, err
	}
	tresc, err := json.Marshal(zlecenie{
		Model:         s.ustawienia.Model,
		KatalogModeli: katalogWag(s.katalogDanych, s.ustawienia.KatalogModeli),
		Teksty:        teksty,
		WagaMb:        WagaModeluMb,
	})
	if err != nil {
		return nil, errors.New("wskaźnik znaczenia: nie da się złożyć zlecenia pomocnika: " + err.Error())
	}
	sciezkaZlecenia, sprzatnij, err := zapiszZlecenie(s.katalogDanych, tresc)
	if err != nil {
		return nil, err
	}
	defer sprzatnij()

	narzedzie := zewnetrzne.Narzedzie{
		Nazwa:   "pomocnik osadzeń (Python)",
		Program: odnajdzInterpreter(s.ustawienia.Program),
		Pakiet:  "python3 wraz z biblioteką fastembed",
	}
	// `-X utf8` idzie zawsze, bo Python bez niego koduje wyjście regionalnie.
	wynik, err := zewnetrzne.Wolaj(ctx, s.uruchamiacz, okno, zasady, obszar,
		narzedzie, []string{"-X", "utf8", skrypt, sciezkaZlecenia}, "", limit)
	if err != nil {
		var brakNarzedzia *zewnetrzne.BrakNarzedzia
		if errors.As(err, &brakNarzedzia) {
			return nil, &BrakSilnika{Rodzaj: BrakInterpretera, Model: s.ustawienia.Model,
				WagaMb: WagaModeluMb, Powod: err.Error()}
		}
		return nil, errors.New("wskaźnik znaczenia: " + err.Error())
	}
	return s.odczytaj(wynik)
}

// odczytaj rozbiera odpowiedź pomocnika i zamienia nazwany brak na typowaną
// odmowę. Odpowiedź nieczytelna jest usterką, a nie brakiem: brak pomocnik umie
// nazwać sam, więc nieczytelność znaczy, że coś innego niż on pisało na wyjście.
func (s *Silnik) odczytaj(wynik zewnetrzne.Wynik) ([][]float32, error) {
	var wczytana odpowiedz
	if err := json.Unmarshal(wynik.Wyjscie, &wczytana); err != nil {
		return nil, errors.New("wskaźnik znaczenia: odpowiedź pomocnika osadzeń " +
			"jest nieczytelna (" + err.Error() + ")" +
			opisDiagnostyki(wynik.Diagnostyka) +
			"; naprawa: zgłosić usterkę pomocnika osadzeń")
	}
	if !wczytana.Ok {
		return nil, &BrakSilnika{
			Rodzaj: wczytana.Brak,
			Model:  s.ustawienia.Model,
			WagaMb: wczytana.WagaMb,
			Powod:  wczytana.Powod,
		}
	}
	return wczytana.Wektory, nil
}

// opisDiagnostyki dokłada do odmowy to, co pomocnik powiedział o sobie sam —
// bez tego członu zostaje „odpowiedź nieczytelna" i żadnej wskazówki.
func opisDiagnostyki(diagnostyka string) string {
	tresc := strings.TrimSpace(diagnostyka)
	if tresc == "" {
		return ""
	}
	return "; pomocnik powiedział: " + skroc(tresc, 400)
}
