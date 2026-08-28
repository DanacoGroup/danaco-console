// Odpowiedzialność pliku: odczyt rejestru okien operacyjnych — definicji okna
// (tabela `okno_operacyjne`) wraz z jej przypięciem do modułu (macierz `okno_operacyjne_modul`).
package dane

import (
	"context"
	"fmt"
)

// oknoOperacyjne to definicja okna wraz z kolejnością właściwą kontekstowi
// odczytu: dla wykazu modułu jest to kolejność z macierzy, dla wykazu okien
// pozamodułowych — kolejność w obrębie kategorii.
type oknoOperacyjne struct {
	ID        int64
	Kod       string
	Nazwa     string
	Rola      string
	Kategoria string
	Kolejnosc int
	Aktywne   bool
}

// RepozytoriumOkienOperacyjnych jest kontraktem rejestru okien operacyjnych modułów całej tej platformy.
type RepozytoriumOkienOperacyjnych interface {
	ListaModulu(ctx context.Context, modulID int64) ([]oknoOperacyjne, error)
	Globalne(ctx context.Context) ([]oknoOperacyjne, error)
}

const (
	poleOknaOperacyjnego = `o.id, o.kod, o.nazwa, o.rola, o.kategoria`

	// Kolejność pochodzi z macierzy, bo to samo okno bywa przypięte do dwóch
	// modułów na różnych pozycjach.
	listaOkienModulu = `SELECT ` + poleOknaOperacyjnego + `, om.kolejnosc, o.aktywne
	                    FROM okno_operacyjne o
	                    JOIN okno_operacyjne_modul om ON om.okno_operacyjne_id = o.id
	                    WHERE om.modul_id = ? AND o.aktywne = 1
	                    ORDER BY om.kolejnosc, o.kod`

	// Okno pozamodułowe rozpoznajemy brakiem przypięcia, a nie osobną kolumną —
	// przynależność do modułu ma jedno miejsce zapisu.
	listaOkienGlobalnych = `SELECT ` + poleOknaOperacyjnego + `, o.kolejnosc, o.aktywne
	                        FROM okno_operacyjne o
	                        WHERE o.aktywne = 1
	                          AND NOT EXISTS (SELECT 1 FROM okno_operacyjne_modul om
	                                           WHERE om.okno_operacyjne_id = o.id)
	                        ORDER BY o.kategoria, o.kolejnosc, o.kod`
)

type repozytoriumOkienOperacyjnych struct {
	zapytania *zapytania
}

func noweRepozytoriumOkienOperacyjnych(z *zapytania) *repozytoriumOkienOperacyjnych {
	return &repozytoriumOkienOperacyjnych{zapytania: z}
}

// ListaModulu zwraca okna robocze jednego modułu w kolejności ich pozycji
// w tym module. Moduł spoza rejestru daje wykaz pusty.
func (r *repozytoriumOkienOperacyjnych) ListaModulu(ctx context.Context, modulID int64) ([]oknoOperacyjne, error) {
	return r.wykaz(ctx, listaOkienModulu, "okien modułu", modulID)
}

// Globalne zwraca okna spoza katalogu modułowego — te, których macierz
// `okno_operacyjne_modul` nie przypina do żadnego modułu.
func (r *repozytoriumOkienOperacyjnych) Globalne(ctx context.Context) ([]oknoOperacyjne, error) {
	return r.wykaz(ctx, listaOkienGlobalnych, "okien pozamodułowych")
}

// wykaz wykonuje zapytanie do bazy danych zwracające wiele wierszy rejestru okien operacyjnych modułu.
func (r *repozytoriumOkienOperacyjnych) wykaz(ctx context.Context, zapytanie, opis string,
	argumenty ...any) ([]oknoOperacyjne, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać %s: %w", opis, err)
	}
	defer wiersze.Close()

	lista := []oknoOperacyjne{}
	for wiersze.Next() {
		okno, err := odczytajOknoOperacyjne(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, okno)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt %s: %w", opis, err)
	}
	return lista, nil
}

// odczytajOknoOperacyjne składa pełną strukturę okna operacyjnego z jednego wiersza wyniku zapytania do bazy.
func odczytajOknoOperacyjne(wiersz skaner) (oknoOperacyjne, error) {
	var okno oknoOperacyjne
	var aktywne int
	err := wiersz.Scan(&okno.ID, &okno.Kod, &okno.Nazwa, &okno.Rola,
		&okno.Kategoria, &okno.Kolejnosc, &aktywne)
	if err != nil {
		return oknoOperacyjne{}, fmt.Errorf("dane: nieczytelny wiersz okna operacyjnego: %w", err)
	}
	okno.Aktywne = aktywne != 0
	return okno, nil
}
