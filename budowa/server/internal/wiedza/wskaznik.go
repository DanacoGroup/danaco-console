// Odpowiedzialność pliku: cała droga od dokumentu do trafienia w jednym
// miejscu — wniesienie treści do wskaźnika znaczenia i przeszukanie go.
// Adapter rdzenia zostaje przy swojej robocie: przekład kontraktu, zasięg
// izolacji, kody odmów.
package wiedza

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Dokument to jedna jednostka treści przed podziałem na fragmenty. Treść
// przychodzi gotowa, a wskaźnik nie wie, skąd — czytanie źródeł należy do
// warstwy, która te repozytoria zna.
type Dokument struct {
	// Zakres — `library`, `history` albo `workspace`.
	Zakres string
	// Zrodlo — nazwa CZYTELNA dla człowieka; wchodzi wprost do trafienia.
	Zrodlo string
	// ZrodloKod — identyfikator, którym da się po dokument sięgnąć.
	ZrodloKod string
	// Tresc — cały tekst dokumentu.
	Tresc string
}

// Wskaznik prowadzi wskaźnik znaczenia: wnosi do niego treść i przeszukuje
// go po pytaniu zadanym przez Operatora.
type Wskaznik struct {
	// osadzarka — źródło wektorów wraz z nazwą modelu.
	osadzarka Osadzarka
	// skladnica — trwałość przy bazie rdzenia (`store/migracja_115_wskaznik_znaczenia.sql`).
	skladnica *Skladnica
	// dlugoscFragmentu — docelowa długość fragmentu w znakach.
	dlugoscFragmentu int
	// teraz oddaje czas w milisekundach epoki — jeden zegar na byt.
	teraz func() int64
}

// NowyWskaznik wiąże źródło wektorów z trwałością.
//
// Zegar domyślny to zegar systemowy; podmienia go `ZZegarem`.
func NowyWskaznik(osadzarka Osadzarka, skladnica *Skladnica, dlugoscFragmentu int) *Wskaznik {
	return &Wskaznik{
		osadzarka:        osadzarka,
		skladnica:        skladnica,
		dlugoscFragmentu: dlugoscFragmentu,
		teraz:            func() int64 { return time.Now().UTC().UnixMilli() },
	}
}

// ZZegarem podstawia źródło czasu zapisu, zamiast zegara systemowego
// użytego domyślnie przy starcie modułu.
func (w *Wskaznik) ZZegarem(teraz func() int64) *Wskaznik {
	if teraz != nil {
		w.teraz = teraz
	}
	return w
}

// Model oddaje nazwę modelu, którym pracuje ten wskaźnik. Jedna prawda: bierze
// się z tego, kto liczy wektory.
func (w *Wskaznik) Model() string {
	if w == nil || w.osadzarka == nil {
		return ""
	}
	return w.osadzarka.Model()
}

// Gotowy sprawdza, czy jest czym liczyć i gdzie zapisać, nie licząc niczego.
// Sprawdzenie idzie przed czytaniem treści, bo inaczej rdzeń przeczytałby
// całą bibliotekę z dysku i dopiero wtedy powiedział „nie ma czym liczyć".
func (w *Wskaznik) Gotowy(ctx context.Context, limit time.Duration) error {
	if err := w.sprawny(); err != nil {
		return err
	}
	return w.osadzarka.Gotowy(ctx, limit)
}

// Wyczysc kasuje wskazane zakresy — droga przebudowy od zera (`rebuild`).
// Wykaz pusty czyści WSZYSTKO; wołający rozstrzyga, czy tego chce.
func (w *Wskaznik) Wyczysc(ctx context.Context, zakresy []string) error {
	if err := w.sprawny(); err != nil {
		return err
	}
	return w.skladnica.UsunZakres(ctx, zakresy)
}

// Policz oddaje liczbę pozycji wskaźnika policzonych modelem tego
// wskaźnika, obecnie ustawionym w konfiguracji.
func (w *Wskaznik) Policz(ctx context.Context) (int, error) {
	if err := w.sprawny(); err != nil {
		return 0, err
	}
	return w.skladnica.Policz(ctx, w.osadzarka.Model())
}

// Wnies wnosi komplet dokumentów do wskaźnika i oddaje liczbę wniesionych
// fragmentów.
//
// Dokument bez treści dającej się podzielić jest POMIJANY, a nie odmawiany:
// plik pusty w bibliotece nie ma prawa odebrać Operatorowi wskaźnika
// pozostałych.
func (w *Wskaznik) Wnies(ctx context.Context, dokumenty []Dokument,
	limit time.Duration) (int, error) {

	if err := w.sprawny(); err != nil {
		return 0, err
	}
	wniesione := 0
	for _, dokument := range dokumenty {
		ile, err := w.wniesDokument(ctx, dokument, limit)
		if err != nil {
			// Liczba wniesionych do tej chwili wraca razem z odmową.
			return wniesione, err
		}
		wniesione += ile
	}
	return wniesione, nil
}

// wniesDokument przeprowadza jeden dokument przez całą drogę: podział,
// osadzenie, zapis do bazy danych.
func (w *Wskaznik) wniesDokument(ctx context.Context, dokument Dokument,
	limit time.Duration) (int, error) {

	fragmenty := Podziel(dokument.Tresc, w.dlugoscFragmentu)
	if len(fragmenty) == 0 {
		return 0, nil
	}
	teksty := make([]string, len(fragmenty))
	for i, fragment := range fragmenty {
		teksty[i] = fragment.Tresc
	}
	wektory, err := w.osadzarka.Osadz(ctx, teksty, limit)
	if err != nil {
		return 0, err
	}
	if len(wektory) != len(teksty) {
		return 0, errors.New("wskaźnik znaczenia: na " + liczba(len(teksty)) +
			" fragmentów dokumentu " + dokument.Zrodlo + " wróciło " + liczba(len(wektory)) +
			" wektorów — wiązanie po pozycji przestałoby cokolwiek znaczyć; " +
			"naprawa: zgłosić usterkę pomocnika osadzeń")
	}

	model := w.osadzarka.Model()
	pozycje := make([]Pozycja, len(fragmenty))
	for i, fragment := range fragmenty {
		pozycje[i] = Pozycja{
			Zakres:    dokument.Zakres,
			Zrodlo:    dokument.Zrodlo,
			ZrodloKod: dokument.ZrodloKod,
			Kolejnosc: fragment.Kolejnosc,
			Tresc:     fragment.Tresc,
			Model:     model,
			Wektor:    Znormalizuj(wektory[i]),
		}
	}
	if err := w.skladnica.UsunZrodlo(ctx, dokument.Zakres, dokument.ZrodloKod); err != nil {
		return 0, err
	}
	return w.skladnica.Zapisz(ctx, pozycje, w.teraz())
}

// Szukaj oddaje `ile` fragmentów najbliższych pytaniu, od najbliższego.
// Pytanie osadza się tym samym modelem, co dokumenty. Wynik pusty jest
// odpowiedzią, nie odmową.
func (w *Wskaznik) Szukaj(ctx context.Context, pytanie string, zakresy []string,
	ile int, limit time.Duration) ([]Trafienie, error) {

	if err := w.sprawny(); err != nil {
		return nil, err
	}
	tresc := strings.TrimSpace(pytanie)
	if tresc == "" {
		return nil, errors.New("wskaźnik znaczenia: pytanie puste — wyszukiwanie po " +
			"znaczeniu nie ma czego porównać z wiedzą Operatora; " +
			"naprawa: podać treść pytania")
	}
	wektory, err := w.osadzarka.Osadz(ctx, []string{tresc}, limit)
	if err != nil {
		return nil, err
	}
	if len(wektory) != 1 {
		return nil, errors.New("wskaźnik znaczenia: pomocnik osadzeń nie oddał wektora " +
			"pytania; naprawa: zgłosić usterkę pomocnika osadzeń")
	}
	pozycje, err := w.skladnica.Pozycje(ctx, zakresy, w.osadzarka.Model())
	if err != nil {
		return nil, err
	}
	return Najblizsze(pozycje, Znormalizuj(wektory[0]), ile), nil
}

// sprawny nazywa brak ogniwa, bez którego wskaźnik nie umie nic.
//
// Odmowa, a nie milcząca pustka: wskaźnik bez źródła wektorów oddawałby zero
// trafień nieodróżnialne od „w wiedzy Operatora tego nie ma".
func (w *Wskaznik) sprawny() error {
	if w == nil || w.osadzarka == nil {
		return errors.New("wskaźnik znaczenia: rdzeń nie ma czym liczyć wektorów — " +
			"wskaźnika nie da się ani zbudować, ani przeszukać; " +
			"naprawa: podpiąć silnik osadzeń przy składaniu rdzenia")
	}
	if w.skladnica == nil {
		return brakSkladnicy()
	}
	return nil
}
