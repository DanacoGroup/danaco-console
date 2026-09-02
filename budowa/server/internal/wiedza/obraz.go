// Odpowiedzialność pliku: silnik osi obrazu — porównanie zdania Operatora
// z obrazami jego biblioteki, przez model z dwiema wieżami, jedną dla pikseli
// i drugą dla słów, uczonymi tak, żeby kończyły w jednej przestrzeni.
package wiedza

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"danacoconsole/server/internal/dane"
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
	// skladnica trzyma wyniki policzone wcześniej, przy bazie rdzenia; pusta
	// oznacza pracę bez trwałości — silnik liczy każdy wynik od nowa, tak jak
	// przed dołożeniem trwałości.
	skladnica *Skladnica
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

// ZeSkladnica podpina trwałość wyników osi obrazu nad bazą rdzenia
// (`store/migracja_405_trwalosc_wektora_obrazu.sql`). Bez niej `Dopasuj`
// woła pomocnika przy każdym wywołaniu, tak jak przed dołożeniem trwałości.
func (s *SilnikObrazu) ZeSkladnica(skladnica *Skladnica) *SilnikObrazu {
	s.skladnica = skladnica
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
// wołający wiąże ocenę z plikiem pozycją. Obraz, którego wynik dla tego
// zdania i modelu leży w składnicy pod niezmienionym odciskiem pliku, nie
// wraca do pomocnika.
func (s *SilnikObrazu) Dopasuj(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	pytanie string, sciezki []string, limit time.Duration) ([]float32, []PominietyObraz, error) {

	if len(sciezki) == 0 {
		return nil, nil, nil
	}

	odciski := make([]string, len(sciezki))
	for i, sciezka := range sciezki {
		odciski[i] = odciskPliku(sciezka)
	}

	zSkladnicy, err := s.podobienstwaZeSkladnicy(ctx, sciezki, pytanie)
	if err != nil {
		return nil, nil, err
	}

	oceny := make([]float32, len(sciezki))
	var brakujaceSciezki []string
	var brakujaceIndeksy []int
	for i, sciezka := range sciezki {
		wpis, jest := zSkladnicy[sciezka]
		if jest && odciski[i] != "" && wpis.odcisk == odciski[i] {
			oceny[i] = wpis.podobienstwo
			continue
		}
		brakujaceSciezki = append(brakujaceSciezki, sciezka)
		brakujaceIndeksy = append(brakujaceIndeksy, i)
	}

	if len(brakujaceSciezki) == 0 {
		return oceny, nil, nil
	}

	swiezeOceny, pominiete, err := s.wolaj(ctx, okno, zasady, obszar, pytanie, brakujaceSciezki, limit)
	if err != nil {
		return nil, nil, err
	}
	if len(swiezeOceny) != len(brakujaceSciezki) {
		return nil, nil, errors.New("wskaźnik znaczenia: pomocnik osi obrazu oddał " +
			liczba(len(swiezeOceny)) + " ocen na " + liczba(len(brakujaceSciezki)) +
			" obrazów — wiązanie po pozycji przestałoby cokolwiek znaczyć; " +
			"naprawa: zgłosić usterkę pomocnika osi obrazu")
	}

	pominieteZbior := make(map[string]struct{}, len(pominiete))
	for _, obraz := range pominiete {
		pominieteZbior[obraz.Obraz] = struct{}{}
	}
	var doZapisu []wpisPodobienstwaObrazu
	for j, sciezka := range brakujaceSciezki {
		i := brakujaceIndeksy[j]
		oceny[i] = swiezeOceny[j]
		if _, pominietyObraz := pominieteZbior[sciezka]; pominietyObraz || odciski[i] == "" {
			continue
		}
		doZapisu = append(doZapisu, wpisPodobienstwaObrazu{
			sciezka: sciezka, odcisk: odciski[i], podobienstwo: swiezeOceny[j],
		})
	}
	if err := s.zapiszPodobienstwaDoSkladnicy(ctx, doZapisu, pytanie); err != nil {
		return nil, nil, err
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

// wpisPodobienstwaObrazu to jeden wiersz tabeli `podobienstwo_obrazu`: ścieżka
// obrazu, odcisk pliku w chwili liczenia i sam wynik.
type wpisPodobienstwaObrazu struct {
	sciezka      string
	odcisk       string
	podobienstwo float32
}

// odciskPliku znaczy plik jego rozmiarem i czasem modyfikacji — zmiana
// któregokolwiek daje inny odcisk. Plik nieodczytywalny oddaje napis pusty;
// wołający traktuje to jako brak potwierdzenia składnicy, nie jako błąd.
func odciskPliku(sciezka string) string {
	info, err := os.Stat(sciezka)
	if err != nil {
		return ""
	}
	return strconv.FormatInt(info.Size(), 10) + ":" + strconv.FormatInt(info.ModTime().UnixNano(), 10)
}

// podobienstwaZeSkladnicy czyta wyniki już policzone dla tego zdania i tego
// modelu, indeksowane po ścieżce. Składnica pusta oddaje wykaz pusty, nie
// odmowę — silnik bez trwałości ma prawo liczyć wszystko od nowa.
func (s *SilnikObrazu) podobienstwaZeSkladnicy(ctx context.Context, sciezki []string,
	pytanie string) (map[string]wpisPodobienstwaObrazu, error) {

	if s.skladnica == nil || len(sciezki) == 0 {
		return nil, nil
	}
	zapytanie := `SELECT sciezka, odcisk, podobienstwo FROM podobienstwo_obrazu
	                WHERE model = ? AND pytanie = ? AND sciezka IN (` +
		znakiZapytania(len(sciezki)) + `) AND ` + dane.WarunekKonta
	argumenty := make([]any, 0, len(sciezki)+3)
	argumenty = append(argumenty, s.ustawienia.ModelObrazu, pytanie)
	for _, sciezka := range sciezki {
		argumenty = append(argumenty, sciezka)
	}
	argumenty = append(argumenty, dane.KontoOperatora(ctx))
	wiersze, err := s.skladnica.baza.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, errors.New("wskaźnik znaczenia: odczyt wyników osi obrazu: " + err.Error())
	}
	defer wiersze.Close()

	wyniki := make(map[string]wpisPodobienstwaObrazu, len(sciezki))
	for wiersze.Next() {
		var wpis wpisPodobienstwaObrazu
		if err := wiersze.Scan(&wpis.sciezka, &wpis.odcisk, &wpis.podobienstwo); err != nil {
			return nil, errors.New("wskaźnik znaczenia: odczyt wiersza wyniku osi obrazu: " +
				err.Error())
		}
		wyniki[wpis.sciezka] = wpis
	}
	if err := wiersze.Err(); err != nil {
		return nil, errors.New("wskaźnik znaczenia: przerwany odczyt wyników osi obrazu: " +
			err.Error())
	}
	return wyniki, nil
}

// zapiszPodobienstwaDoSkladnicy wnosi świeżo policzone wyniki, nadpisując
// wynik poprzedni tej samej trójki ścieżka-model-pytanie. Składnica pusta jest
// przemilczana — bez trwałości nie ma czego zapisać.
func (s *SilnikObrazu) zapiszPodobienstwaDoSkladnicy(ctx context.Context,
	wpisy []wpisPodobienstwaObrazu, pytanie string) error {

	if s.skladnica == nil || len(wpisy) == 0 {
		return nil
	}
	transakcja, err := s.skladnica.baza.BeginTx(ctx, nil)
	if err != nil {
		return errors.New("wskaźnik znaczenia: otwarcie transakcji osi obrazu: " + err.Error())
	}
	defer transakcja.Rollback()

	// Więz UNIQUE (sciezka, model, pytanie) obejmuje całą tabelę: człon DO UPDATE bez zawężenia nadpisałby wynik konta cudzego.
	polecenie, err := transakcja.PrepareContext(ctx, `
		INSERT INTO podobienstwo_obrazu (sciezka, odcisk, model, pytanie, podobienstwo,
		                                 utworzono, konto_id)
		VALUES (?, ?, ?, ?, ?, ?, `+dane.WskazanieKonta+`)
		ON CONFLICT(sciezka, model, pytanie) DO UPDATE SET
		    odcisk       = excluded.odcisk,
		    podobienstwo = excluded.podobienstwo,
		    utworzono    = excluded.utworzono
		WHERE `+dane.WarunekKonta)
	if err != nil {
		return errors.New("wskaźnik znaczenia: przygotowanie zapisu osi obrazu: " + err.Error())
	}
	defer polecenie.Close()

	chwila := time.Now().UTC().UnixMilli()
	for _, wpis := range wpisy {
		wynik, err := polecenie.ExecContext(ctx, wpis.sciezka, wpis.odcisk,
			s.ustawienia.ModelObrazu, pytanie, wpis.podobienstwo, chwila,
			dane.KontoOperatora(ctx), dane.KontoOperatora(ctx))
		if err != nil {
			return errors.New("wskaźnik znaczenia: zapis wyniku osi obrazu dla " +
				wpis.sciezka + ": " + err.Error())
		}
		zmienione, err := wynik.RowsAffected()
		if err != nil {
			return errors.New("wskaźnik znaczenia: nieznany wynik zapisu osi obrazu dla " +
				wpis.sciezka + ": " + err.Error())
		}
		if zmienione == 0 {
			return fmt.Errorf("wskaźnik znaczenia: wynik osi obrazu dla %s stoi na koncie innym: %w",
				wpis.sciezka, dane.ErrKolizjaWiersza)
		}
	}
	if err := transakcja.Commit(); err != nil {
		return errors.New("wskaźnik znaczenia: domknięcie zapisu osi obrazu: " + err.Error())
	}
	return nil
}
