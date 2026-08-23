// Odpowiedzialność pliku: trwałość wskaźnika znaczenia — zapis, odczyt
// i czyszczenie pozycji tabeli `fragment_wiedzy`.
//
// Wskaźnik leży przy bazie rdzenia, a nie w osobnym pliku: druga baza obok
// pierwszej to drugi plik do przeniesienia, drugi do kopii zapasowej i drugi,
// który da się zgubić osobno. Wektor jest bytem wtórnym — odtwarzalnym z treści
// jednym przebiegiem `knowledge.index` — więc jego utrata kosztuje czas
// procesora, a nie wiedzę Operatora. Tak samo wtórny jest indeks FTS5 treści
// biblioteki i leży w tej samej bazie.
//
// Baza nadal nie przechowuje pliku: bajty treści leżą w magazynie biblioteki
// pod sumą sha256, a tu leży wektor i fragment tekstu, z którego go policzono.
// Fragment jest w tabeli z jednego powodu:
// kontrakt każe oddać w trafieniu `text`, a odtwarzanie go przy każdym zapytaniu
// znaczyłoby otwieranie plików źródłowych i ponowne dzielenie ich na fragmenty —
// czyli wykonywanie całej pracy indeksowania po to, żeby oddać dwa zdania.
//
// Składnica mówi własnymi typami, nie typami kontraktu. Pakiet `wiedza` nie zna
// kontraktu i nie ma go poznać: przekład na `shared.KnowledgeHit` należy do
// adaptera rdzenia, tak samo jak w silniku mowy.
package wiedza

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Pozycja to jeden wiersz wskaźnika: fragment treści, jego pochodzenie i wektor.
type Pozycja struct {
	// Zakres — `library`, `history` albo `workspace`. Napisy przychodzą
	// z adaptera, bo to on zna wyliczenie kontraktu.
	Zakres string
	// Zrodlo — nazwa czytelna dla człowieka: nazwa pliku biblioteki, tytuł okna
	// rozmowy, ścieżka względna pliku przestrzeni. Wchodzi wprost do trafienia,
	// bo bez wskazania źródła model cytowałby bez możliwości sprawdzenia.
	Zrodlo string
	// ZrodloKod — identyfikator, którym da się po źródło sięgnąć (kod pliku
	// biblioteki, identyfikator wiadomości). Pusty jest stanem poprawnym dla
	// źródeł, które kodu nie mają — wtedy trafienie oddaje samą nazwę.
	ZrodloKod string
	// Kolejnosc — numer fragmentu w dokumencie źródłowym.
	Kolejnosc int
	// Tresc — sam fragment, ten, który wróci Operatorowi jako cytat.
	Tresc string
	// Model — nazwa modelu, którym policzono wektor. Bez niej nie wiadomo,
	// z czym wolno ten wektor porównywać.
	Model string
	// Wektor — współrzędne znormalizowane (patrz `podobienstwo.go`).
	Wektor []float32
}

// Skladnica prowadzi tabelę wskaźnika.
type Skladnica struct {
	baza *sql.DB
}

// NowaSkladnica zakłada składnicę nad bazą rdzenia. Baza pusta oddaje nil —
// wołający znosi to sam, tak jak dziennik transkrypcji.
func NowaSkladnica(baza *sql.DB) *Skladnica {
	if baza == nil {
		return nil
	}
	return &Skladnica{baza: baza}
}

// Zapisz wnosi komplet pozycji jednego źródła w jednej transakcji i zwraca
// liczbę wniesionych.
//
// Źródło wchodzi w całości albo wcale. Fragmenty jednego dokumentu wniesione
// połowicznie dałyby wskaźnik, który o tym dokumencie wie, ale zna go do
// połowy — a Operator nie ma jak tego zobaczyć: zapytanie o drugą połowę wróci
// puste tak samo, jak wraca dla dokumentu nieindeksowanego. Dlatego przerwanie
// w środku cofa całość.
//
// Powtórne indeksowanie nadpisuje, a nie dokłada. Warunek jednoznaczności
// (zakres, kod źródła, kolejność, model) czyni z zapisu upsert, więc dokument
// zmieniony i zaindeksowany ponownie ma tyle fragmentów, ile ma treści — a nie
// sumę wszystkich swoich wersji. Fragmenty nadmiarowe po skróceniu dokumentu
// kasuje `UsunZrodlo` wołane przez adapter przed zapisem.
func (s *Skladnica) Zapisz(ctx context.Context, pozycje []Pozycja, chwila int64) (int, error) {
	if s == nil {
		return 0, brakSkladnicy()
	}
	if len(pozycje) == 0 {
		return 0, nil
	}
	transakcja, err := s.baza.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("wskaźnik znaczenia: otwarcie transakcji: %w", err)
	}
	defer transakcja.Rollback()

	polecenie, err := transakcja.PrepareContext(ctx, `
		INSERT INTO fragment_wiedzy (zakres, zrodlo, zrodlo_kod, kolejnosc, tresc,
		                             model, wymiar, wektor, utworzono)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(zakres, zrodlo_kod, kolejnosc, model) DO UPDATE SET
		    zrodlo    = excluded.zrodlo,
		    tresc     = excluded.tresc,
		    wymiar    = excluded.wymiar,
		    wektor    = excluded.wektor,
		    utworzono = excluded.utworzono`)
	if err != nil {
		return 0, fmt.Errorf("wskaźnik znaczenia: przygotowanie zapisu: %w", err)
	}
	defer polecenie.Close()

	for _, pozycja := range pozycje {
		_, err := polecenie.ExecContext(ctx, pozycja.Zakres, pozycja.Zrodlo, pozycja.ZrodloKod,
			pozycja.Kolejnosc, pozycja.Tresc, pozycja.Model, len(pozycja.Wektor),
			NaBajty(pozycja.Wektor), chwila)
		if err != nil {
			return 0, fmt.Errorf("wskaźnik znaczenia: zapis fragmentu %d źródła %s: %w",
				pozycja.Kolejnosc, pozycja.Zrodlo, err)
		}
	}
	if err := transakcja.Commit(); err != nil {
		return 0, fmt.Errorf("wskaźnik znaczenia: domknięcie zapisu: %w", err)
	}
	return len(pozycje), nil
}

// UsunZrodlo kasuje wszystkie fragmenty jednego źródła w danym zakresie.
// Wołane przed zapisem, żeby dokument skrócony nie zostawił po sobie fragmentów
// z treści, której już nie ma — a które wracałyby jako cytat z dokumentu.
func (s *Skladnica) UsunZrodlo(ctx context.Context, zakres, zrodloKod string) error {
	if s == nil {
		return brakSkladnicy()
	}
	_, err := s.baza.ExecContext(ctx,
		`DELETE FROM fragment_wiedzy WHERE zakres = ? AND zrodlo_kod = ?`, zakres, zrodloKod)
	if err != nil {
		return fmt.Errorf("wskaźnik znaczenia: usunięcie fragmentów źródła %s: %w", zrodloKod, err)
	}
	return nil
}

// UsunZakres czyści cały zakres — droga przebudowy od zera (`rebuild`).
// Zakres pusty czyści wszystko; wołający rozstrzyga, czy tego chce.
func (s *Skladnica) UsunZakres(ctx context.Context, zakresy []string) error {
	if s == nil {
		return brakSkladnicy()
	}
	if len(zakresy) == 0 {
		_, err := s.baza.ExecContext(ctx, `DELETE FROM fragment_wiedzy`)
		return err
	}
	for _, zakres := range zakresy {
		if _, err := s.baza.ExecContext(ctx,
			`DELETE FROM fragment_wiedzy WHERE zakres = ?`, zakres); err != nil {
			return fmt.Errorf("wskaźnik znaczenia: czyszczenie zakresu %s: %w", zakres, err)
		}
	}
	return nil
}

// Pozycje odczytuje wiersze zakresów policzone wskazanym modelem.
//
// Zawężenie po modelu jest warunkiem poprawności, nie optymalizacją. Operator,
// który zmienił ustawienie modelu, ma w tabeli wektory z dwóch przestrzeni;
// porównanie pytania z wektorem cudzego modelu daje liczbę, która wygląda jak
// trafność i nią nie jest. Stare wiersze zostają w tabeli świadomie — wracają
// do użytku, gdy Operator wróci do poprzedniego modelu, a `rebuild` je czyści.
func (s *Skladnica) Pozycje(ctx context.Context, zakresy []string, model string) ([]Pozycja, error) {
	if s == nil {
		return nil, brakSkladnicy()
	}
	zapytanie := `SELECT zakres, zrodlo, zrodlo_kod, kolejnosc, tresc, model, wektor
	                FROM fragment_wiedzy WHERE model = ?`
	argumenty := []any{model}
	if len(zakresy) > 0 {
		zapytanie += " AND zakres IN (" + znakiZapytania(len(zakresy)) + ")"
		for _, zakres := range zakresy {
			argumenty = append(argumenty, zakres)
		}
	}
	wiersze, err := s.baza.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("wskaźnik znaczenia: odczyt wskaźnika: %w", err)
	}
	defer wiersze.Close()

	pozycje := []Pozycja{}
	for wiersze.Next() {
		var pozycja Pozycja
		var bajty []byte
		if err := wiersze.Scan(&pozycja.Zakres, &pozycja.Zrodlo, &pozycja.ZrodloKod,
			&pozycja.Kolejnosc, &pozycja.Tresc, &pozycja.Model, &bajty); err != nil {
			return nil, fmt.Errorf("wskaźnik znaczenia: odczyt wiersza wskaźnika: %w", err)
		}
		wektor, err := ZBajtow(bajty)
		if err != nil {
			return nil, err
		}
		pozycja.Wektor = wektor
		pozycje = append(pozycje, pozycja)
	}
	return pozycje, wiersze.Err()
}

// Policz zwraca liczbę pozycji wskaźnika dla modelu — pole `total` odpowiedzi
// `knowledge.index`. Liczone zapytaniem, a nie długością odczytanego wykazu:
// wykaz bywa dziesiątkami tysięcy wierszy z wektorami, a pytanie brzmi „ile",
// nie „które".
func (s *Skladnica) Policz(ctx context.Context, model string) (int, error) {
	if s == nil {
		return 0, brakSkladnicy()
	}
	var ile int
	err := s.baza.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM fragment_wiedzy WHERE model = ?`, model).Scan(&ile)
	if err != nil {
		return 0, fmt.Errorf("wskaźnik znaczenia: rachunek pozycji: %w", err)
	}
	return ile, nil
}

// znakiZapytania składa listę znaków zapytania dla klauzuli IN. Argumenty
// wiązane, nie sklejane — zakres przychodzi z żądania i wklejony w SQL byłby
// dokładnie tym, przed czym chroni wiązanie.
func znakiZapytania(ile int) string {
	znaki := make([]byte, 0, ile*3)
	for i := 0; i < ile; i++ {
		if i > 0 {
			znaki = append(znaki, ',', ' ')
		}
		znaki = append(znaki, '?')
	}
	return string(znaki)
}

// brakSkladnicy nazywa jedyny stan, w którym składnica nie umie nic.
func brakSkladnicy() error {
	return errors.New("wskaźnik znaczenia: rdzeń nie ma bazy, w której miałby leżeć " +
		"wskaźnik — wektorów nie ma gdzie zapisać ani skąd odczytać; " +
		"naprawa: podpiąć bazę przy składaniu rdzenia")
}
