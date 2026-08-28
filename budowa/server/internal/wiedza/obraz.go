// Odpowiedzialność pliku: silnik osi obrazu — porównanie zdania Operatora
// z obrazami jego biblioteki, przez model z dwiema wieżami, jedną dla pikseli
// i drugą dla słów, uczonymi tak, żeby kończyły w jednej przestrzeni.
package wiedza

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
)

const (
	// LimitOsiObrazu — granica czasu jednego zapytania osi obrazu. Obejmuje
	// pobranie wag przy pierwszym uruchomieniu oraz przebieg wieży obrazu przez
	// wszystkie porównywane pliki.
	LimitOsiObrazu = 20 * time.Minute
	// GranicaObrazow — sufit liczby obrazów jednego zapytania. Koszt rośnie
	// wprost proporcjonalnie do liczby plików, a biblioteka Operatora bywa
	// tysiącami zdjęć — bez sufitu jedno zapytanie zajęłoby maszynę na godziny.
	GranicaObrazow = 200
)

// SilnikObrazu porównuje zdanie z obrazami pomocnikiem lokalnym; niesie
// uruchamiacz procesów, katalog danych na maszynie rdzenia i nastawy modelu.
type SilnikObrazu struct {
	uruchamiacz   session.Uruchamiacz
	katalogDanych string
	ustawienia    Ustawienia
}

// NowySilnikObrazu zakłada silnik na uruchamiaczu procesów i katalogu danych,
// z nastawami domyślnymi; do innych nastaw służy metoda ZUstawieniami.
func NowySilnikObrazu(uruchamiacz session.Uruchamiacz, katalogDanych string) *SilnikObrazu {
	return &SilnikObrazu{
		uruchamiacz:   uruchamiacz,
		katalogDanych: katalogDanych,
		ustawienia:    UstawieniaDomyslne(),
	}
}

// ZUstawieniami oddaje silnikowi komplet nastaw złożony z konfiguracji,
// zamiast nastaw domyślnych założonych przy tworzeniu silnika.
func (s *SilnikObrazu) ZUstawieniami(u Ustawienia) *SilnikObrazu {
	s.ustawienia = u
	return s
}

// Model oddaje nazwę modelu osi obrazu. Wchodzi do odpowiedzi komendy i do
// treści odmowy, gdy silnik nie potrafi porównać zdania z obrazami.
func (s *SilnikObrazu) Model() string {
	return s.ustawienia.ModelObrazu
}

// zlecenieObrazu i odpowiedzObrazu to kształt rozmowy z pomocnikiem. Pola
// odpowiadają co do znaku kluczom w `pomocnik_obrazu.py`.
type zlecenieObrazu struct {
	Model         string   `json:"model"`
	KatalogModeli string   `json:"katalogModeli"`
	Pytanie       string   `json:"pytanie"`
	Obrazy        []string `json:"obrazy"`
	WagaMb        int      `json:"wagaMb"`
}

type odpowiedzObrazu struct {
	Ok        bool             `json:"ok"`
	Model     string           `json:"model"`
	Oceny     []float32        `json:"oceny"`
	Pominiete []PominietyObraz `json:"pominiete"`
	Brak      string           `json:"brak"`
	Powod     string           `json:"powod"`
	WagaMb    int              `json:"wagaMb"`
}

// PominietyObraz nazywa plik, którego pomocnik nie otworzył, wraz z powodem.
// Wykaz wraca do wołającego, bo obraz pominięty i obraz niepodobny do zdania
// wyglądają w odpowiedzi tak samo — oba po prostu w niej nie stoją.
type PominietyObraz struct {
	Obraz string `json:"obraz"`
	Powod string `json:"powod"`
}

// Gotowy sprawdza, czy jest czym porównywać, nie porównując niczego; woła
// pomocnika bez pytania i bez plików, tylko po to, żeby sprawdzić dostępność.
func (s *SilnikObrazu) Gotowy(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar, limit time.Duration) error {

	_, _, err := s.wolaj(ctx, okno, zasady, obszar, "", nil, limit)
	return err
}

// Dopasuj oddaje po jednej ocenie na obraz, w kolejności ścieżek, wraz
// z wykazem plików, których pomocnik nie otworzył. Kolejność jest warunkiem:
// wołający wiąże ocenę z plikiem pozycją.
func (s *SilnikObrazu) Dopasuj(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	pytanie string, sciezki []string, limit time.Duration) ([]float32, []PominietyObraz, error) {

	if len(sciezki) == 0 {
		return nil, nil, nil
	}
	oceny, pominiete, err := s.wolaj(ctx, okno, zasady, obszar, pytanie, sciezki, limit)
	if err != nil {
		return nil, nil, err
	}
	if len(oceny) != len(sciezki) {
		return nil, nil, errors.New("wskaźnik znaczenia: pomocnik osi obrazu oddał " +
			liczba(len(oceny)) + " ocen na " + liczba(len(sciezki)) +
			" obrazów — wiązanie po pozycji przestałoby cokolwiek znaczyć; " +
			"naprawa: zgłosić usterkę pomocnika osi obrazu")
	}
	return oceny, pominiete, nil
}

// wolaj przeprowadza jedno uruchomienie pomocnika i czyta jego odpowiedź;
// wywołują ją metody Gotowy i Dopasuj, każda z inną treścią zlecenia.
func (s *SilnikObrazu) wolaj(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	pytanie string, sciezki []string,
	limit time.Duration) ([]float32, []PominietyObraz, error) {

	skrypt, err := wylozPomocnika(s.katalogDanych, nazwaSkryptuObrazu, skryptObrazu)
	if err != nil {
		return nil, nil, err
	}
	tresc, err := json.Marshal(zlecenieObrazu{
		Model: s.ustawienia.ModelObrazu,
		KatalogModeli: katalogWagOsobny(s.katalogDanych,
			s.ustawienia.KatalogObrazu, podkatalogObrazu),
		Pytanie: pytanie,
		Obrazy:  sciezki,
		WagaMb:  WagaObrazuMb,
	})
	if err != nil {
		return nil, nil, errors.New("wskaźnik znaczenia: nie da się złożyć zlecenia " +
			"osi obrazu: " + err.Error())
	}
	sciezkaZlecenia, sprzatnij, err := zapiszZlecenie(s.katalogDanych, tresc)
	if err != nil {
		return nil, nil, err
	}
	defer sprzatnij()

	narzedzie := zewnetrzne.Narzedzie{
		Nazwa:   "pomocnik osi obrazu (Python)",
		Program: odnajdzInterpreter(s.ustawienia.Program),
		Pakiet:  "python3 wraz z bibliotekami torch, transformers i pillow",
	}
	wynik, err := zewnetrzne.Wolaj(ctx, s.uruchamiacz, okno, zasady, obszar,
		narzedzie, []string{"-X", "utf8", skrypt, sciezkaZlecenia}, "", limit)
	if err != nil {
		var brakNarzedzia *zewnetrzne.BrakNarzedzia
		if errors.As(err, &brakNarzedzia) {
			return nil, nil, &BrakSilnika{Silnik: SilnikDlaObrazu, Rodzaj: BrakInterpretera,
				Model: s.ustawienia.ModelObrazu, WagaMb: WagaObrazuMb, Powod: err.Error()}
		}
		return nil, nil, errors.New("wskaźnik znaczenia: " + err.Error())
	}
	return s.odczytaj(wynik)
}

// odczytaj rozbiera odpowiedź pomocnika i zamienia nazwany brak na typowaną
// odmowę, którą wołający rozpozna po typie BrakSilnika.
func (s *SilnikObrazu) odczytaj(wynik zewnetrzne.Wynik) ([]float32, []PominietyObraz, error) {
	var wczytana odpowiedzObrazu
	if err := json.Unmarshal(wynik.Wyjscie, &wczytana); err != nil {
		return nil, nil, errors.New("wskaźnik znaczenia: odpowiedź pomocnika osi obrazu " +
			"jest nieczytelna (" + err.Error() + ")" +
			opisDiagnostyki(wynik.Diagnostyka) +
			"; naprawa: zgłosić usterkę pomocnika osi obrazu")
	}
	if !wczytana.Ok {
		return nil, nil, &BrakSilnika{
			Silnik: SilnikDlaObrazu,
			Rodzaj: wczytana.Brak,
			Model:  s.ustawienia.ModelObrazu,
			WagaMb: wczytana.WagaMb,
			Powod:  wczytana.Powod,
		}
	}
	return wczytana.Oceny, wczytana.Pominiete, nil
}

// Obraz wiąże obraz biblioteki z jego oceną wobec zdania, wraz ze ścieżką
// na maszynie rdzenia i identyfikatorem, którym da się po obraz sięgnąć.
type Obraz struct {
	// Zrodlo — nazwa pliku czytelna dla człowieka.
	Zrodlo string
	// ZrodloKod — identyfikator, którym da się po obraz sięgnąć.
	ZrodloKod string
	// Rodzaj — zapisany przy pliku rodzaj treści; pusty jest stanem poprawnym.
	Rodzaj string
	// Sciezka — droga do bajtów obrazu na maszynie rdzenia; nie wychodzi poza rdzeń.
	Sciezka string
	// Podobienst — kosinus zdania z obrazem w przestrzeni wspólnej.
	Podobienst float32
}

// NajblizszeObrazy układa obrazy według ocen i oddaje `ile` najbliższych.
// Ocena niedodatnia odpada: obraz, którego pomocnik nie otworzył, ma ocenę
// zerową i nie ma prawa wrócić jako trafienie.
func NajblizszeObrazy(obrazy []Obraz, oceny []float32, ile int) []Obraz {
	if len(oceny) != len(obrazy) {
		return nil
	}
	wybrane := make([]Obraz, 0, len(obrazy))
	for i, obraz := range obrazy {
		if oceny[i] <= 0 {
			continue
		}
		obraz.Podobienst = oceny[i]
		wybrane = append(wybrane, obraz)
	}
	sortujObrazy(wybrane)
	if ile > 0 && len(wybrane) > ile {
		wybrane = wybrane[:ile]
	}
	return wybrane
}

// sortujObrazy układa obrazy od najbliższego zdaniu, zachowując kolejność
// wejściową przy równej ocenie.
func sortujObrazy(obrazy []Obraz) {
	sort.SliceStable(obrazy, func(i, j int) bool {
		return obrazy[i].Podobienst > obrazy[j].Podobienst
	})
}
