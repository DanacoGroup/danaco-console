// Odpowiedzialność pliku: trwałość wskaźnika znaczenia — zapis, odczyt
// i czyszczenie pozycji tabeli fragment_wiedzy; wskaźnik leży przy bazie
// rdzenia.
package wiedza

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/server/internal/dane"
)

// Pozycja to jeden wiersz wskaźnika: fragment treści, jego pochodzenie
// i wektor, w kształcie zapisywanym do tabeli fragment_wiedzy.
type Pozycja struct {
	// Zakres — library, history albo workspace; napisy przychodzą z adaptera.
	Zakres string
	// Zrodlo — nazwa czytelna dla człowieka: pliku biblioteki, okna rozmowy
	// albo ścieżki przestrzeni.
	Zrodlo string
	// ZrodloKod — identyfikator, którym da się po źródło sięgnąć; pusty dla
	// źródeł bez kodu.
	ZrodloKod string
	// Kolejnosc — numer fragmentu w dokumencie źródłowym.
	Kolejnosc int
	// Tresc — sam fragment, ten, który wróci Operatorowi jako cytat.
	Tresc string
	// Model — nazwa modelu, którym policzono wektor.
	Model string
	// Wektor — współrzędne znormalizowane funkcją podobieństwa kosinusowego.
	Wektor []float32
}

// Skladnica prowadzi tabelę wskaźnika, oddając zapis, odczyt i czyszczenie
// pozycji rdzeniowi bez znajomości kontraktu.
type Skladnica struct {
	baza *sql.DB
}

// NowaSkladnica zakłada składnicę nad bazą rdzenia; baza pusta oddaje nil,
// wołający znosi to sam, tak jak dziennik transkrypcji.
func NowaSkladnica(baza *sql.DB) *Skladnica {
	if baza == nil {
		return nil
	}
	return &Skladnica{baza: baza}
}

// Zapisz wnosi komplet pozycji jednego źródła w jednej transakcji i zwraca
// liczbę wniesionych; źródło wchodzi w całości albo wcale.
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

	// Więz UNIQUE (zakres, zrodlo_kod, kolejnosc, model) obejmuje całą tabelę: człon DO UPDATE bez zawężenia nadpisałby fragment konta cudzego.
	polecenie, err := transakcja.PrepareContext(ctx, `
		INSERT INTO fragment_wiedzy (zakres, zrodlo, zrodlo_kod, kolejnosc, tresc,
		                             model, wymiar, wektor, utworzono, konto_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, `+dane.WskazanieKonta+`)
		ON CONFLICT(zakres, zrodlo_kod, kolejnosc, model,
		            COALESCE(konto_id, 0)) DO UPDATE SET
		    zrodlo    = excluded.zrodlo,
		    tresc     = excluded.tresc,
		    wymiar    = excluded.wymiar,
		    wektor    = excluded.wektor,
		    utworzono = excluded.utworzono
		WHERE `+dane.WarunekKonta)
	if err != nil {
		return 0, fmt.Errorf("wskaźnik znaczenia: przygotowanie zapisu: %w", err)
	}
	defer polecenie.Close()

	for _, pozycja := range pozycje {
		wynik, err := polecenie.ExecContext(ctx, pozycja.Zakres, pozycja.Zrodlo, pozycja.ZrodloKod,
			pozycja.Kolejnosc, pozycja.Tresc, pozycja.Model, len(pozycja.Wektor),
			NaBajty(pozycja.Wektor), chwila, dane.KontoOperatora(ctx), dane.KontoOperatora(ctx))
		if err != nil {
			return 0, fmt.Errorf("wskaźnik znaczenia: zapis fragmentu %d źródła %s: %w",
				pozycja.Kolejnosc, pozycja.Zrodlo, err)
		}
		zmienione, err := wynik.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("wskaźnik znaczenia: nieznany wynik zapisu fragmentu %d źródła %s: %w",
				pozycja.Kolejnosc, pozycja.Zrodlo, err)
		}
		if zmienione == 0 {
			return 0, fmt.Errorf("wskaźnik znaczenia: fragment %d źródła %s stoi na koncie innym: %w",
				pozycja.Kolejnosc, pozycja.Zrodlo, dane.ErrKolizjaWiersza)
		}
	}
	if err := transakcja.Commit(); err != nil {
		return 0, fmt.Errorf("wskaźnik znaczenia: domknięcie zapisu: %w", err)
	}
	return len(pozycje), nil
}

// UsunZrodlo kasuje wszystkie fragmenty jednego źródła w danym zakresie;
// wołane przed zapisem, żeby nie zostały fragmenty nieaktualne.
func (s *Skladnica) UsunZrodlo(ctx context.Context, zakres, zrodloKod string) error {
	if s == nil {
		return brakSkladnicy()
	}
	_, err := s.baza.ExecContext(ctx,
		`DELETE FROM fragment_wiedzy WHERE zakres = ? AND zrodlo_kod = ? AND `+dane.WarunekKonta,
		zakres, zrodloKod, dane.KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("wskaźnik znaczenia: usunięcie fragmentów źródła %s: %w", zrodloKod, err)
	}
	return nil
}

// UsunZakres czyści cały zakres — droga przebudowy od zera (rebuild);
// zakres pusty czyści wszystko, wołający rozstrzyga, czy tego chce.
func (s *Skladnica) UsunZakres(ctx context.Context, zakresy []string) error {
	if s == nil {
		return brakSkladnicy()
	}
	if len(zakresy) == 0 {
		_, err := s.baza.ExecContext(ctx, `DELETE FROM fragment_wiedzy WHERE `+dane.WarunekKonta,
			dane.KontoOperatora(ctx))
		return err
	}
	for _, zakres := range zakresy {
		if _, err := s.baza.ExecContext(ctx,
			`DELETE FROM fragment_wiedzy WHERE zakres = ? AND `+dane.WarunekKonta,
			zakres, dane.KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("wskaźnik znaczenia: czyszczenie zakresu %s: %w", zakres, err)
		}
	}
	return nil
}

// Pozycje odczytuje wiersze zakresów policzone wskazanym modelem;
// zawężenie po modelu jest warunkiem poprawności, nie optymalizacją.
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
	zapytanie += " AND " + dane.WarunekKonta
	argumenty = append(argumenty, dane.KontoOperatora(ctx))
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

// Policz zwraca liczbę pozycji wskaźnika dla modelu — pole total
// odpowiedzi knowledge.index, liczone zapytaniem, a nie długością wykazu.
func (s *Skladnica) Policz(ctx context.Context, model string) (int, error) {
	if s == nil {
		return 0, brakSkladnicy()
	}
	var ile int
	err := s.baza.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM fragment_wiedzy WHERE model = ? AND `+dane.WarunekKonta,
		model, dane.KontoOperatora(ctx)).Scan(&ile)
	if err != nil {
		return 0, fmt.Errorf("wskaźnik znaczenia: rachunek pozycji: %w", err)
	}
	return ile, nil
}

// znakiZapytania składa listę znaków zapytania dla klauzuli IN; argumenty
// wiązane, nie sklejane, bo zakres przychodzi z żądania klienta.
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

// brakSkladnicy nazywa jedyny stan, w którym składnica nie umie nic,
// i podaje wołającemu jednoznaczne rozpoznanie tego stanu.
func brakSkladnicy() error {
	return errors.New("wskaźnik znaczenia: serwer nie ma bazy, w której miałby leżeć " +
		"wskaźnik — wektorów nie ma gdzie zapisać ani skąd odczytać; " +
		"naprawa: podpiąć bazę przy składaniu serwera")
}
