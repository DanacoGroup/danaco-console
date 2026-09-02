// Komponenty makiety modułu Design: komponent wisi na oknie, jego instancja na kompozycji konkretnej planszy.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// KomponentDesignu to wiersz `komponent_design`; WariantyJSON niesie zapis kontraktu bez rozkładania w warstwie danych.
type KomponentDesignu struct {
	ID               int64
	Kod              string
	Okno             string
	Nazwa            string
	WariantyJSON     *string
	ZestawZetonowKod *string
	Liczba           int
	Zaktualizowano   string
}

// InstancjaKomponentuDesignu to jedno wystąpienie komponentu na planszy.
type InstancjaKomponentuDesignu struct {
	ID           int64
	KomponentID  int64
	KompozycjaID int64
	WarstwaKod   string
	Wariant      *string
}

const (
	kolumnyKomponentuDesignu = `k.id, k.identyfikator_zewnetrzny, k.okno, k.nazwa, k.warianty_json,
	                            k.zestaw_zetonow_kod,
	                            (SELECT COUNT(*) FROM instancja_komponentu_design i
	                             WHERE i.komponent_id = k.id),
	                            k.zaktualizowano`

	zapiszKomponentDesignu = `INSERT INTO komponent_design
	                          (identyfikator_zewnetrzny, okno, nazwa, warianty_json,
	                           zestaw_zetonow_kod, zaktualizowano, konto_id)
	                          VALUES (?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'), ` + WskazanieKonta + `)
	                          ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                              nazwa = excluded.nazwa,
	                              warianty_json = excluded.warianty_json,
	                              zestaw_zetonow_kod = excluded.zestaw_zetonow_kod,
	                              zaktualizowano = excluded.zaktualizowano
	                          WHERE ` + WarunekKonta

	pobierzKomponentDesignu = `SELECT ` + kolumnyKomponentuDesignu +
		` FROM komponent_design k WHERE k.identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaKomponentowDesignu = `SELECT ` + kolumnyKomponentuDesignu +
		` FROM komponent_design k WHERE k.okno = ? AND ` + WarunekKonta +
		` ORDER BY k.zaktualizowano DESC, k.id DESC`

	listaKomponentowDesignuJeden = `SELECT ` + kolumnyKomponentuDesignu +
		` FROM komponent_design k WHERE k.okno = ? AND ` + WarunekKonta +
		` AND k.identyfikator_zewnetrzny = ?`

	wstawInstancjeKomponentuDesignu = `INSERT INTO instancja_komponentu_design
	                                   (komponent_id, kompozycja_id, warstwa_kod, wariant)
	                                   VALUES (?, ?, ?, ?)
	                                   ON CONFLICT(warstwa_kod) DO UPDATE SET
	                                       komponent_id = excluded.komponent_id,
	                                       kompozycja_id = excluded.kompozycja_id,
	                                       wariant = excluded.wariant`

	listaInstancjiKomponentuDesignu = `SELECT id, komponent_id, kompozycja_id, warstwa_kod, wariant
	                                   FROM instancja_komponentu_design
	                                   WHERE komponent_id = ? ORDER BY id`
)

// ZapiszKomponentDesignu zakłada komponent albo nadpisuje zastany po identyfikatorze zewnętrznym w bazie.
func (r *repozytoriumDesignu) ZapiszKomponentDesignu(ctx context.Context,
	komponent KomponentDesignu) (KomponentDesignu, error) {

	if komponent.Kod == "" {
		return KomponentDesignu{}, fmt.Errorf("dane: komponent design bez identyfikatora")
	}
	if komponent.Okno == "" {
		return KomponentDesignu{}, fmt.Errorf("dane: komponent design %q bez okna", komponent.Kod)
	}
	if komponent.Nazwa == "" {
		return KomponentDesignu{}, fmt.Errorf("dane: komponent design %q bez nazwy", komponent.Kod)
	}

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKomponentDesignu)
	if err != nil {
		return KomponentDesignu{}, err
	}
	_, err = polecenie.ExecContext(ctx, komponent.Kod, komponent.Okno, komponent.Nazwa,
		tekstDoKolumny(komponent.WariantyJSON), tekstDoKolumny(komponent.ZestawZetonowKod),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return KomponentDesignu{}, fmt.Errorf("dane: nie można zapisać komponentu design %q: %w",
			komponent.Kod, err)
	}
	return r.KomponentDesignuPoKodzie(ctx, komponent.Kod)
}

// KomponentDesignuPoKodzie zwraca komponent po kodzie wraz z liczbą instancji.
func (r *repozytoriumDesignu) KomponentDesignuPoKodzie(ctx context.Context,
	kod string) (KomponentDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKomponentDesignu)
	if err != nil {
		return KomponentDesignu{}, err
	}
	komponent, err := odczytajKomponentDesignu(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return KomponentDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return KomponentDesignu{}, fmt.Errorf("dane: nieczytelny wiersz komponentu design %q: %w",
			kod, err)
	}
	return komponent, nil
}

// KomponentyDesignu zwraca komponenty okna; wskazanie kodu zawęża wykaz do jednego komponentu.
func (r *repozytoriumDesignu) KomponentyDesignu(ctx context.Context, okno string,
	kod *string) ([]KomponentDesignu, error) {

	zapytanie := listaKomponentowDesignu
	argumenty := []any{okno, KontoOperatora(ctx)}
	if kod != nil && *kod != "" {
		zapytanie = listaKomponentowDesignuJeden
		argumenty = append(argumenty, *kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać komponentów design okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []KomponentDesignu{}
	for wiersze.Next() {
		komponent, err := odczytajKomponentDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz komponentu design okna %q: %w", okno, err)
		}
		lista = append(lista, komponent)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt komponentów design okna %q: %w", okno, err)
	}
	return lista, nil
}

// ZapiszInstancjeKomponentuDesignu utrwala wystąpienie komponentu na planszy wraz z jego pełnym położeniem.
func (r *repozytoriumDesignu) ZapiszInstancjeKomponentuDesignu(ctx context.Context,
	instancja InstancjaKomponentuDesignu) error {

	if instancja.WarstwaKod == "" {
		return fmt.Errorf("dane: instancja komponentu design bez warstwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawInstancjeKomponentuDesignu)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, instancja.KomponentID, instancja.KompozycjaID,
		instancja.WarstwaKod, tekstDoKolumny(instancja.Wariant))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać instancji komponentu design %q: %w",
			instancja.WarstwaKod, err)
	}
	return nil
}

// InstancjeKomponentuDesignu zwraca wystąpienia komponentu na planszach wprost z bazy danych repozytorium.
func (r *repozytoriumDesignu) InstancjeKomponentuDesignu(ctx context.Context,
	komponentID int64) ([]InstancjaKomponentuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaInstancjiKomponentuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, komponentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać instancji komponentu design %d: %w",
			komponentID, err)
	}
	defer wiersze.Close()

	lista := []InstancjaKomponentuDesignu{}
	for wiersze.Next() {
		var instancja InstancjaKomponentuDesignu
		var wariant sql.NullString
		if err := wiersze.Scan(&instancja.ID, &instancja.KomponentID, &instancja.KompozycjaID,
			&instancja.WarstwaKod, &wariant); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz instancji komponentu design %d: %w",
				komponentID, err)
		}
		instancja.Wariant = tekstZKolumny(wariant)
		lista = append(lista, instancja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt instancji komponentu design %d: %w",
			komponentID, err)
	}
	return lista, nil
}

// odczytajKomponentDesignu składa strukturę komponentu wprost z jednego wiersza wyniku zapytania do SQL.
func odczytajKomponentDesignu(wiersz skaner) (KomponentDesignu, error) {
	var komponent KomponentDesignu
	var warianty, zestaw sql.NullString
	err := wiersz.Scan(&komponent.ID, &komponent.Kod, &komponent.Okno, &komponent.Nazwa,
		&warianty, &zestaw, &komponent.Liczba, &komponent.Zaktualizowano)
	if err != nil {
		return KomponentDesignu{}, err
	}
	komponent.WariantyJSON = tekstZKolumny(warianty)
	komponent.ZestawZetonowKod = tekstZKolumny(zestaw)
	return komponent, nil
}
