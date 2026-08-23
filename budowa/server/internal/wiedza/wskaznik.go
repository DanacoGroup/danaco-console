// Odpowiedzialność pliku: CAŁA DROGA od dokumentu do trafienia w jednym
// miejscu — wniesienie treści do wskaźnika znaczenia i przeszukanie go.
//
// Ten plik zbiera drogę w jeden byt. Adapter rdzenia zostaje przy swojej
// robocie: przekład kontraktu, zasięg izolacji, kody odmów.
//
// ── DLACZEGO WNOSZENIE IDZIE DOKUMENT PO DOKUMENCIE ─────────────────────────
// Biblioteka Operatora bywa gigabajtem tekstu. Jedna transakcja na całość
// znaczyłaby komplet wektorów w pamięci rdzenia, a przerwanie w połowie —
// cofnięcie wszystkiego, co już policzono kwadransem procesora. Dlatego granica
// niepodzielności biegnie po DOKUMENCIE: każdy, który wszedł, wszedł w całości
// (patrz `Skladnica.Zapisz`), przerwany przebieg zostawia wskaźnik NIEPEŁNY,
// ale SPÓJNY, a powtórzenie dokańcza resztę.
//
// ── DLACZEGO KASOWANIE ŹRÓDŁA POPRZEDZA ZAPIS ───────────────────────────────
// Dokument SKRÓCONY od poprzedniego przebiegu zostawiłby inaczej fragmenty
// treści, której już w nim nie ma — a wracałyby jako cytat z dokumentu, w
// którym ich nie ma. Sam zapis jest wprawdzie nadpisujący (warunek
// jednoznaczności: zakres, kod źródła, kolejność, model), ale nadpisuje tylko
// te numery fragmentów, które przyszły; nadmiarowe zostają. Kasowanie jest
// w TEJ SAMEJ chwili co zapis i dla tego samego źródła, więc okno, w którym
// dokument jest niewidoczny, trwa tyle, ile jedna transakcja.
package wiedza

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Dokument to jedna jednostka treści przed podziałem na fragmenty.
//
// TREŚĆ PRZYCHODZI GOTOWA, a wskaźnik nie wie, skąd. Czytanie źródeł —
// repozytorium biblioteki, historii rozmów, plików przestrzeni roboczej —
// należy do warstwy, która te repozytoria zna (`core/adapter_modul_wiedza_zrodla.go`).
// Druga droga do wierszy biblioteki założona tutaj byłaby drugą prawdą o tym,
// co Operator w niej ma.
type Dokument struct {
	// Zakres — `library`, `history` albo `workspace`. Napis przychodzi
	// z adaptera, bo to on zna wyliczenie kontraktu.
	Zakres string
	// Zrodlo — nazwa CZYTELNA dla człowieka; wchodzi wprost do trafienia.
	Zrodlo string
	// ZrodloKod — identyfikator, którym da się po dokument sięgnąć. Pusty jest
	// stanem poprawnym dla źródeł, które kodu nie mają.
	ZrodloKod string
	// Tresc — cały tekst dokumentu.
	Tresc string
}

// Wskaznik prowadzi wskaźnik znaczenia: wnosi do niego treść i przeszukuje go.
type Wskaznik struct {
	// osadzarka — źródło wektorów wraz z nazwą modelu (patrz `osadzarka.go`).
	osadzarka Osadzarka
	// skladnica — trwałość przy bazie rdzenia (`store/migracja_115_wskaznik_znaczenia.sql`).
	skladnica *Skladnica
	// dlugoscFragmentu — docelowa długość fragmentu w znakach. Wartość spoza
	// przedziału sensownego sprowadza `Podziel` do domyślnej — jeden strażnik
	// tej liczby, nie dwóch.
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

// ZZegarem podstawia źródło czasu zapisu.
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

// Gotowy sprawdza, czy jest czym liczyć i gdzie zapisać, NIE licząc niczego.
//
// SPRAWDZENIE IDZIE PRZED CZYTANIEM TREŚCI i to jest rozstrzygnięcie:
// odwrotnie byłoby taniej w przypadku szczęśliwym i znacznie gorzej w
// przypadku brakującego silnika — rdzeń przeczytałby całą bibliotekę z dysku,
// podzielił ją na fragmenty i dopiero wtedy powiedział „nie ma czym liczyć".
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

// Policz oddaje liczbę pozycji wskaźnika policzonych modelem tego wskaźnika.
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
			// Liczba wniesionych do tej chwili wraca RAZEM z odmową: przebieg
			// przerwany w połowie zostawia wskaźnik niepełny i wołający ma
			// wiedzieć, ile weszło, zamiast zgadywać, czy weszło cokolwiek.
			return wniesione, err
		}
		wniesione += ile
	}
	return wniesione, nil
}

// wniesDokument przeprowadza jeden dokument przez całą drogę.
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
//
// PYTANIE OSADZA SIĘ TYM SAMYM MODELEM, CO DOKUMENTY, i nie jest to szczegół
// techniczny: wektor pytania z modelu innego niż wektory wskaźnika daje iloczyn
// skalarny, który jest liczbą i nie znaczy nic. Zawężenie odczytu po nazwie
// modelu jest jedyną obroną przed tym po zmianie ustawienia — a nazwa jest
// jedna, bo pyta się o nią tego, kto liczy.
//
// WYNIK PUSTY JEST ODPOWIEDZIĄ, NIE ODMOWĄ: wskaźnik pusty albo
// wiedza bez związku z pytaniem znaczą „nie mam na to nic". Inaczej niż BRAK
// SILNIKA, który jest odmową, bo wtedy rdzeń nie wie, czy ma coś, czy nie ma.
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
