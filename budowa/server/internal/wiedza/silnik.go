// Odpowiedzialność pliku: złożenie silnika osadzeń w jeden byt i przeprowadzenie
// jednego zlecenia — od tekstów do wektorów.
//
// Jeden silnik, nie dwa. Osadzenia liczy się w dwóch chwilach — przy budowaniu
// wskaźnika (`knowledge.index`) i przy zapytaniu (`knowledge.search`) — i musi
// je liczyć ten sam model tym samym sposobem. Wektor dokumentu policzony jednym
// modelem, a wektor pytania drugim, dają iloczyn skalarny, który jest liczbą
// i nawet wygląda sensownie, a nie znaczy nic. Dlatego model jest jednym polem
// jednego bytu, a nie dwoma stałymi w dwóch miejscach.
//
// Proces startuje wyłącznie przez `zewnetrzne.Wolaj`: w całym produkcie stoi
// dokładnie jedno `exec.Command`, a każde uruchomienie idzie tą samą bramą
// izolacji okna i tym samym obejmowaniem potomstwa. Proces Pythona liczący na
// wielu wątkach bez objęcia drzewem zostawiałby sieroty po każdym przekroczeniu
// czasu.
//
// Każde wołanie wczytuje model od nowa: prawie cały czas zlecenia to start
// procesu i wczytanie wag (rząd gigabajta) do pamięci, a nie samo porównanie
// wektorów. Proces rezydentny skróciłby zapytanie, ale wymaga dwukierunkowej
// rozmowy z procesem żyjącym między żądaniami, a `zewnetrzne.Wolaj` prowadzi
// rozmowę jednorazową i jest jedyną dozwoloną drogą startu procesu — zmiana tego
// jest zmianą w `zewnetrzne/**`, poza terytorium tego pliku.
//
// Granica czasu jest dwojaka. Pierwsze uruchomienie pobiera wagi modelu, więc
// granica budowania wskaźnika jest liczona w minutach; zapytanie ma wagi już na
// dysku i granica jest liczona w sekundach. Jedna wspólna granica byłaby albo za
// krótka na pobranie, albo tak długa, że zawieszone zapytanie wyglądałoby na
// pracujące.
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
	// wskaźnika. Obejmuje pobranie wag przy pierwszym uruchomieniu.
	LimitBudowania = 20 * time.Minute
	// LimitZapytania — granica czasu osadzenia jednego pytania.
	LimitZapytania = 2 * time.Minute
	// wielkoscPartii — ile fragmentów idzie do pomocnika za jednym razem.
	// Uruchomienie procesu kosztuje kilka sekund (wczytanie wag), więc partia
	// ma być duża; wykaz w pliku JSON o kilkuset fragmentach to kilkaset
	// kilobajtów, czyli nic.
	wielkoscPartii = 256
)

// Silnik liczy osadzenia pomocnikiem lokalnym.
//
// Ustawienia trzymane są w silniku, a nie odczytywane przy każdym zleceniu:
// wołający składa je z konfiguracji zasięgu i podaje gotowe — dwie drogi do tej
// samej wartości byłyby dwiema prawdami.
type Silnik struct {
	// uruchamiacz — jedyna droga startu procesu w drzewie.
	uruchamiacz session.Uruchamiacz
	// katalogDanych — ten sam katalog, w którym leżą baza i magazyn biblioteki.
	katalogDanych string
	// ustawienia — komplet nastaw obowiązujący dla zleceń tego silnika.
	ustawienia Ustawienia
}

// NowySilnik zakłada silnik na uruchamiaczu procesów i katalogu danych.
func NowySilnik(uruchamiacz session.Uruchamiacz, katalogDanych string) *Silnik {
	return &Silnik{
		uruchamiacz:   uruchamiacz,
		katalogDanych: katalogDanych,
		ustawienia:    UstawieniaDomyslne(),
	}
}

// ZUstawieniami oddaje silnikowi komplet nastaw złożony z konfiguracji.
func (s *Silnik) ZUstawieniami(u Ustawienia) *Silnik {
	s.ustawienia = u
	return s
}

// Ustawienia oddaje nastawy, którymi silnik dziś pracuje.
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

// Gotowy sprawdza, czy silnik ma czym liczyć, nie licząc niczego.
//
// Sprawdzenie idzie przed budowaniem wskaźnika i przed zapytaniem, bo odmowa
// „nie ma czym" jest dla Operatora czymś innym niż „liczyło i się wywróciło".
// Kosztuje jedno uruchomienie pomocnika z pustym wykazem tekstów — pomocnik
// przygotowuje model i wraca, więc przy pierwszym razie pobierze też wagi.
func (s *Silnik) Gotowy(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar, limit time.Duration) error {

	_, err := s.wolaj(ctx, okno, zasady, obszar, nil, limit)
	return err
}

// Osadz zamienia teksty na wektory, zachowując ich kolejność.
//
// Partiami, nie wszystko naraz: wykaz kilkudziesięciu tysięcy fragmentów
// w jednym pliku zlecenia zająłby pomocnikowi pamięć proporcjonalną do całej
// biblioteki. Kolejność wektorów odpowiada kolejności tekstów i to jest
// warunek — wołający wiąże je pozycją, nie treścią.
//
// Wykaz pusty oddaje wykaz pusty bez uruchamiania procesu: „osadź nic" nie jest
// pytaniem o gotowość silnika, tylko pracą, której nie ma (od pytania jest
// `Gotowy`).
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

// wolaj przeprowadza jedno uruchomienie pomocnika i czyta jego odpowiedź.
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
	// `-X utf8` idzie zawsze: pomocnik oddaje polski tekst w odpowiedzi, a Python
	// bez tego przełącznika koduje wyjście według ustawień regionalnych systemu —
	// na polskim Windowsie stroną 1250, w której JSON rozpada się na krzaki.
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
