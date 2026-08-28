// Odpowiedzialność pliku: silnik przesiewu — drugi przebieg wyszukiwania po
// znaczeniu, w którym krzyżowy koder układa kandydatów pierwszego przebiegu
// na nowo.
package wiedza

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
)

const (
	// LimitPrzesiewu — granica czasu jednego przesiewu. Obejmuje pobranie wag
	// przy pierwszym uruchomieniu, dlatego liczona jest w minutach; samo
	// ocenienie kilkudziesięciu kandydatów na wagach stojących zajmuje sekundy.
	LimitPrzesiewu = 20 * time.Minute
	// KandydaciDomyslni — ilu kandydatów pierwszego przebiegu wchodzi do
	// przesiewu, gdy żądanie nie mówi ilu.
	KandydaciDomyslni = 50
	// GranicaKandydatow — sufit liczby ocenianych par. Koszt rośnie wprost
	// proporcjonalnie do liczby kandydatów, więc żądanie o cały wskaźnik
	// zamieniłoby zapytanie w budowanie wskaźnika.
	GranicaKandydatow = 200
)

// SilnikPrzesiewu ocenia pary pytanie–fragment pomocnikiem lokalnym.
//
// Nastawy trzymane są w silniku, a nie odczytywane przy każdym zleceniu — z tego
// samego powodu co w silniku osadzeń: dwie drogi do tej samej wartości byłyby
// dwiema prawdami.
type SilnikPrzesiewu struct {
	uruchamiacz   session.Uruchamiacz
	katalogDanych string
	ustawienia    Ustawienia
}

// NowySilnikPrzesiewu zakłada silnik na uruchamiaczu procesów i katalogu
// danych, z nastawami domyślnymi.
func NowySilnikPrzesiewu(uruchamiacz session.Uruchamiacz, katalogDanych string) *SilnikPrzesiewu {
	return &SilnikPrzesiewu{
		uruchamiacz:   uruchamiacz,
		katalogDanych: katalogDanych,
		ustawienia:    UstawieniaDomyslne(),
	}
}

// ZUstawieniami oddaje silnikowi komplet nastaw złożony z konfiguracji
// poziomu okna, zamiast domyślnych.
func (s *SilnikPrzesiewu) ZUstawieniami(u Ustawienia) *SilnikPrzesiewu {
	s.ustawienia = u
	return s
}

// Model oddaje nazwę krzyżowego kodera, którym liczone są oceny. Wchodzi do
// treści odmowy — „nie ma czym przesiać" bez nazwy modelu nie mówi, czego
// brakuje.
func (s *SilnikPrzesiewu) Model() string {
	return s.ustawienia.ModelPrzesiewu
}

// zleceniePrzesiewu i odpowiedzPrzesiewu to kształt rozmowy z pomocnikiem. Pola
// odpowiadają co do znaku kluczom w `pomocnik_przesiewu.py` — rozjazd zamieniłby
// nazwany brak w brak nierozpoznany.
type zleceniePrzesiewu struct {
	Model         string   `json:"model"`
	KatalogModeli string   `json:"katalogModeli"`
	Pytanie       string   `json:"pytanie"`
	Teksty        []string `json:"teksty"`
	WagaMb        int      `json:"wagaMb"`
	OknoTokenow   int      `json:"oknoTokenow"`
}

type odpowiedzPrzesiewu struct {
	Ok     bool      `json:"ok"`
	Model  string    `json:"model"`
	Oceny  []float32 `json:"oceny"`
	Brak   string    `json:"brak"`
	Powod  string    `json:"powod"`
	WagaMb int       `json:"wagaMb"`
}

// Gotowy sprawdza, czy jest czym przesiewać, nie oceniając niczego. Osobne
// pytanie, tak samo jak w silniku osadzeń: odmowa „nie ma czym" jest dla
// Operatora czymś innym niż „liczyło i się wywróciło".
func (s *SilnikPrzesiewu) Gotowy(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar, limit time.Duration) error {

	_, err := s.wolaj(ctx, okno, zasady, obszar, "", nil, limit)
	return err
}

// Przesiej oddaje po jednej ocenie na tekst, w kolejności tekstów. Kolejność
// jest warunkiem: wołający wiąże ocenę z kandydatem pozycją, a nie treścią.
func (s *SilnikPrzesiewu) Przesiej(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	pytanie string, teksty []string, limit time.Duration) ([]float32, error) {

	if len(teksty) == 0 {
		return nil, nil
	}
	oceny, err := s.wolaj(ctx, okno, zasady, obszar, pytanie, teksty, limit)
	if err != nil {
		return nil, err
	}
	if len(oceny) != len(teksty) {
		return nil, errors.New("wskaźnik znaczenia: pomocnik przesiewu oddał " +
			liczba(len(oceny)) + " ocen na " + liczba(len(teksty)) +
			" kandydatów — wiązanie po pozycji przestałoby cokolwiek znaczyć; " +
			"naprawa: zgłosić usterkę pomocnika przesiewu")
	}
	return oceny, nil
}

// wolaj przeprowadza jedno uruchomienie pomocnika i czyta jego odpowiedź,
// zwracając oceny albo błąd wywołania.
func (s *SilnikPrzesiewu) wolaj(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	pytanie string, teksty []string, limit time.Duration) ([]float32, error) {

	skrypt, err := wylozPomocnika(s.katalogDanych, nazwaSkryptuPrzesiewu, skryptPrzesiewu)
	if err != nil {
		return nil, err
	}
	tresc, err := json.Marshal(zleceniePrzesiewu{
		Model: s.ustawienia.ModelPrzesiewu,
		KatalogModeli: katalogWagOsobny(s.katalogDanych,
			s.ustawienia.KatalogPrzesiewu, podkatalogPrzesiewu),
		Pytanie:     pytanie,
		Teksty:      teksty,
		WagaMb:      WagaPrzesiewuMb,
		OknoTokenow: oknoPrzesiewuTokenow,
	})
	if err != nil {
		return nil, errors.New("wskaźnik znaczenia: nie da się złożyć zlecenia przesiewu: " +
			err.Error())
	}
	sciezkaZlecenia, sprzatnij, err := zapiszZlecenie(s.katalogDanych, tresc)
	if err != nil {
		return nil, err
	}
	defer sprzatnij()

	narzedzie := zewnetrzne.Narzedzie{
		Nazwa:   "pomocnik przesiewu (Python)",
		Program: odnajdzInterpreter(s.ustawienia.Program),
		Pakiet:  "python3 wraz z bibliotekami torch i transformers",
	}
	// `-X utf8` idzie zawsze, bo pomocnik oddaje polski tekst.
	wynik, err := zewnetrzne.Wolaj(ctx, s.uruchamiacz, okno, zasady, obszar,
		narzedzie, []string{"-X", "utf8", skrypt, sciezkaZlecenia}, "", limit)
	if err != nil {
		var brakNarzedzia *zewnetrzne.BrakNarzedzia
		if errors.As(err, &brakNarzedzia) {
			return nil, &BrakSilnika{Silnik: SilnikDlaPrzesiewu, Rodzaj: BrakInterpretera,
				Model: s.ustawienia.ModelPrzesiewu, WagaMb: WagaPrzesiewuMb, Powod: err.Error()}
		}
		return nil, errors.New("wskaźnik znaczenia: " + err.Error())
	}
	return s.odczytaj(wynik)
}

// odczytaj rozbiera odpowiedź pomocnika i zamienia nazwany brak na typowaną
// odmowę. Odpowiedź nieczytelna jest usterką, a nie brakiem — tak samo jak
// w silniku osadzeń.
func (s *SilnikPrzesiewu) odczytaj(wynik zewnetrzne.Wynik) ([]float32, error) {
	var wczytana odpowiedzPrzesiewu
	if err := json.Unmarshal(wynik.Wyjscie, &wczytana); err != nil {
		return nil, errors.New("wskaźnik znaczenia: odpowiedź pomocnika przesiewu " +
			"jest nieczytelna (" + err.Error() + ")" +
			opisDiagnostyki(wynik.Diagnostyka) +
			"; naprawa: zgłosić usterkę pomocnika przesiewu")
	}
	if !wczytana.Ok {
		return nil, &BrakSilnika{
			Silnik: SilnikDlaPrzesiewu,
			Rodzaj: wczytana.Brak,
			Model:  s.ustawienia.ModelPrzesiewu,
			WagaMb: wczytana.WagaMb,
			Powod:  wczytana.Powod,
		}
	}
	return wczytana.Oceny, nil
}

// GranicaKandydatowZadania sprowadza wskazanie żądania do liczby kandydatów.
// Wartość mniejsza od granicy trafień nie jest podnoszona do niej po cichu.
func GranicaKandydatowZadania(wskazanie *int) int {
	if wskazanie == nil || *wskazanie <= 0 {
		return KandydaciDomyslni
	}
	if *wskazanie > GranicaKandydatow {
		return GranicaKandydatow
	}
	return *wskazanie
}

// PoPrzesiewie układa trafienia według ocen kodera i przycina do `ile`. Ocena
// zastępuje podobieństwo, a nie dokłada się do niego, bo obie miary nie leżą
// w jednej skali.
func PoPrzesiewie(kandydaci []Trafienie, oceny []float32, ile int) []Trafienie {
	if len(oceny) != len(kandydaci) {
		return kandydaci
	}
	przesiane := make([]Trafienie, len(kandydaci))
	for i, kandydat := range kandydaci {
		kandydat.Podobienst = oceny[i]
		przesiane[i] = kandydat
	}
	posortujMalejaco(przesiane)
	if ile > 0 && len(przesiane) > ile {
		przesiane = przesiane[:ile]
	}
	return przesiane
}
