// Dziennik audytu biblioteki: przyrostowy (dopisanie i odczyt); zasób wskazany kodem, nie kluczem obcym.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// WpisAudytuBiblioteki to wiersz dziennika audytu biblioteki niosący jedno zdarzenie już zaszłe wcześniej.
type WpisAudytuBiblioteki struct {
	ID       int64
	Kod      string
	PlikKod  *string
	Czynnosc string
	Sprawca  string
	Opis     *string
	Chwila   string
}

// FiltrAudytuBiblioteki niesie zawężenia wykazu wpisów dziennika audytu biblioteki po polach zdarzenia.
type FiltrAudytuBiblioteki struct {
	PlikKod   *string
	Czynnosci []string
	Sprawca   *string
	Od        *string
	Do        *string
	Limit     int
	Offset    int
}

const (
	kolumnyWpisuAudytuBiblioteki = `id, identyfikator_zewnetrzny, plik_kod, czynnosc, sprawca, opis, chwila`

	zapiszWpisAudytuBiblioteki = `INSERT INTO wpis_audytu_biblioteki
	                              (identyfikator_zewnetrzny, plik_kod, czynnosc, sprawca, opis, konto_id)
	                              VALUES (?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzWpisAudytuBiblioteki = `SELECT ` + kolumnyWpisuAudytuBiblioteki + `
	                               FROM wpis_audytu_biblioteki WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

// ZapiszWpisAudytu dopisuje zdarzenie do dziennika audytu biblioteki wraz z jego pełną treścią zapisu.
func (r *repozytoriumBiblioteki) ZapiszWpisAudytu(ctx context.Context,
	wpis WpisAudytuBiblioteki) (WpisAudytuBiblioteki, error) {

	if strings.TrimSpace(wpis.Kod) == "" {
		return WpisAudytuBiblioteki{}, fmt.Errorf("dane: wpis audytu biblioteki bez identyfikatora")
	}
	if strings.TrimSpace(wpis.Sprawca) == "" {
		return WpisAudytuBiblioteki{}, fmt.Errorf("dane: wpis audytu %q bez sprawcy", wpis.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWpisAudytuBiblioteki)
	if err != nil {
		return WpisAudytuBiblioteki{}, err
	}
	_, err = polecenie.ExecContext(ctx, wpis.Kod, tekstDoKolumny(wpis.PlikKod), wpis.Czynnosc,
		wpis.Sprawca, tekstDoKolumny(wpis.Opis), KontoOperatora(ctx))
	if err != nil {
		return WpisAudytuBiblioteki{}, fmt.Errorf("dane: nie można zapisać wpisu audytu %q: %w",
			wpis.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzWpisAudytuBiblioteki)
	if err != nil {
		return WpisAudytuBiblioteki{}, err
	}
	zapisany, err := odczytajWpisAudytuBiblioteki(odczyt.QueryRowContext(ctx, wpis.Kod, KontoOperatora(ctx)))
	if err != nil {
		return WpisAudytuBiblioteki{}, fmt.Errorf("dane: nieczytelny wpis audytu %q: %w", wpis.Kod, err)
	}
	return zapisany, nil
}

// WpisyAudytu zwraca wpisy dziennika audytu od najnowszego wraz z liczbą wpisów spełniających warunki.
func (r *repozytoriumBiblioteki) WpisyAudytu(ctx context.Context,
	filtr FiltrAudytuBiblioteki) ([]WpisAudytuBiblioteki, int, error) {

	warunki := []string{"1 = 1", WarunekKonta}
	argumenty := []any{KontoOperatora(ctx)}
	if filtr.PlikKod != nil && *filtr.PlikKod != "" {
		warunki = append(warunki, "plik_kod = ?")
		argumenty = append(argumenty, *filtr.PlikKod)
	}
	if filtr.Sprawca != nil && *filtr.Sprawca != "" {
		warunki = append(warunki, "sprawca = ?")
		argumenty = append(argumenty, *filtr.Sprawca)
	}
	if len(filtr.Czynnosci) > 0 {
		miejsca := make([]string, 0, len(filtr.Czynnosci))
		for _, czynnosc := range filtr.Czynnosci {
			miejsca = append(miejsca, "?")
			argumenty = append(argumenty, czynnosc)
		}
		warunki = append(warunki, "czynnosc IN ("+strings.Join(miejsca, ", ")+")")
	}
	if filtr.Od != nil && *filtr.Od != "" {
		warunki = append(warunki, "chwila >= ?")
		argumenty = append(argumenty, *filtr.Od)
	}
	if filtr.Do != nil && *filtr.Do != "" {
		warunki = append(warunki, "chwila <= ?")
		argumenty = append(argumenty, *filtr.Do)
	}
	warunek := strings.Join(warunki, " AND ")

	zapytanie := `SELECT ` + kolumnyWpisuAudytuBiblioteki + ` FROM wpis_audytu_biblioteki
	              WHERE ` + warunek + ` ORDER BY chwila DESC, id DESC LIMIT ? OFFSET ?`
	wiersze, err := r.db.QueryContext(ctx, zapytanie,
		append(append([]any{}, argumenty...), granicaWykazu(filtr.Limit),
			przesuniecieWykazu(filtr.Offset))...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać dziennika audytu: %w", err)
	}
	defer wiersze.Close()

	lista := []WpisAudytuBiblioteki{}
	for wiersze.Next() {
		wpis, err := odczytajWpisAudytuBiblioteki(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz dziennika audytu: %w", err)
		}
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt dziennika audytu: %w", err)
	}

	var lacznie int
	zapytanieLiczby := `SELECT COUNT(*) FROM wpis_audytu_biblioteki WHERE ` + warunek
	if err := r.db.QueryRowContext(ctx, zapytanieLiczby, argumenty...).Scan(&lacznie); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć wpisów audytu: %w", err)
	}
	return lista, lacznie, nil
}

// odczytajWpisAudytuBiblioteki składa wpis wprost z jednego wiersza wyniku zapytania SQL do bazy danych.
func odczytajWpisAudytuBiblioteki(wiersz skaner) (WpisAudytuBiblioteki, error) {
	var wpis WpisAudytuBiblioteki
	var plikKod, opis sql.NullString
	err := wiersz.Scan(&wpis.ID, &wpis.Kod, &plikKod, &wpis.Czynnosc, &wpis.Sprawca,
		&opis, &wpis.Chwila)
	if err != nil {
		return WpisAudytuBiblioteki{}, err
	}
	wpis.PlikKod, wpis.Opis = tekstZKolumny(plikKod), tekstZKolumny(opis)
	return wpis, nil
}
