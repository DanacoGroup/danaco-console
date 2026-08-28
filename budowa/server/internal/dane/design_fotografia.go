// Plik prowadzi warsztat fotografii modułu Design: łańcuch edycji zasobu obrazowego i nastawy powtarzalne; ten plik
// przechowuje to, czego kolumna wariantu zasobu nie niesie — czynność i jej nastawy, a nastawa jest zapisem JSON, nie wierszem na pole.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// CzynnoscFotografiiDesignu to wiersz tabeli `czynnosc_fotografii_design` —
// jedno ogniwo łańcucha edycji zasobu.
type CzynnoscFotografiiDesignu struct {
	ID             int64
	ZasobID        int64
	ZasobZrodlaID  *int64
	Komenda        string
	NastawyJSON    *string
	PoliczonePrzez *string
	Utworzono      string
}

// NastawaFotografiiDesignu to wiersz tabeli `nastawa_fotografii_design` — zestaw
// czynności zapisany pod nazwą.
type NastawaFotografiiDesignu struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	CzynnosciJSON  string
	Zaktualizowano string
}

const (
	kolumnyCzynnosciFotografiiDesignu = `id, zasob_id, zasob_zrodla_id, komenda, nastawy_json,
	                                     policzone_przez, utworzono`

	// Czynność ZAWSZE wstawia nowy wiersz: łańcuch edycji jest historią, a
	// historia się nie nadpisuje. Ta sama czynność puszczona dwa razy na tym
	// samym zasobie to dwa ogniwa i taka jest prawda o pracy Operatora.
	zapiszCzynnoscFotografiiDesignu = `INSERT INTO czynnosc_fotografii_design
	                                   (zasob_id, zasob_zrodla_id, komenda, nastawy_json,
	                                    policzone_przez)
	                                   VALUES (?, ?, ?, ?, ?)`

	listaCzynnosciFotografiiDesignu = `SELECT ` + kolumnyCzynnosciFotografiiDesignu +
		` FROM czynnosc_fotografii_design WHERE zasob_id = ? ORDER BY id`

	// Łańcuch wstecz: czynność, z której powstał wskazany zasób. Wiersz jest
	// najwyżej jeden, ale zapytanie oddaje wykaz — ten sam zasób mógłby
	// teoretycznie powstać dwiema drogami, a odczyt ma pokazać stan bazy, nie
	// założenie o niej.
	czynnoscFotografiiDesignuZasobu = `SELECT ` + kolumnyCzynnosciFotografiiDesignu +
		` FROM czynnosc_fotografii_design WHERE zasob_id = ? ORDER BY id DESC LIMIT 1`

	kolumnyNastawyFotografiiDesignu = `id, identyfikator_zewnetrzny, okno, nazwa,
	                                   czynnosci_json, zaktualizowano`

	zapiszNastaweFotografiiDesignu = `INSERT INTO nastawa_fotografii_design
	                                  (identyfikator_zewnetrzny, okno, nazwa, czynnosci_json,
	                                   zaktualizowano)
	                                  VALUES (?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                                  ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                      nazwa = excluded.nazwa,
	                                      czynnosci_json = excluded.czynnosci_json,
	                                      zaktualizowano = excluded.zaktualizowano`

	pobierzNastaweFotografiiDesignu = `SELECT ` + kolumnyNastawyFotografiiDesignu +
		` FROM nastawa_fotografii_design WHERE identyfikator_zewnetrzny = ?`

	listaNastawFotografiiDesignu = `SELECT ` + kolumnyNastawyFotografiiDesignu +
		` FROM nastawa_fotografii_design WHERE okno = ? ORDER BY nazwa, id`
)

// ZapiszCzynnoscFotografiiDesignu dokłada ogniwo do łańcucha edycji; niepowodzenie zapisu jest błędem oddanym wołającemu, nie milczeniem.
func (r *repozytoriumDesignu) ZapiszCzynnoscFotografiiDesignu(ctx context.Context,
	czynnosc CzynnoscFotografiiDesignu) error {

	if czynnosc.ZasobID == 0 {
		return fmt.Errorf("dane: czynność fotografii design bez zasobu wynikowego")
	}
	if czynnosc.Komenda == "" {
		return fmt.Errorf("dane: czynność fotografii design zasobu %d bez nazwy komendy",
			czynnosc.ZasobID)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszCzynnoscFotografiiDesignu)
	if err != nil {
		return err
	}
	var zrodlo any
	if czynnosc.ZasobZrodlaID != nil {
		zrodlo = *czynnosc.ZasobZrodlaID
	}
	if _, err := polecenie.ExecContext(ctx, czynnosc.ZasobID, zrodlo, czynnosc.Komenda,
		tekstDoKolumny(czynnosc.NastawyJSON), tekstDoKolumny(czynnosc.PoliczonePrzez)); err != nil {

		return fmt.Errorf("dane: nie można zapisać czynności fotografii design zasobu %d: %w",
			czynnosc.ZasobID, err)
	}
	return nil
}

// CzynnosciFotografiiDesignuZasobu oddaje ogniwa łańcucha wskazanego zasobu, od najstarszego, wprost z bazy.
func (r *repozytoriumDesignu) CzynnosciFotografiiDesignuZasobu(ctx context.Context,
	zasobID int64) ([]CzynnoscFotografiiDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaCzynnosciFotografiiDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, zasobID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać czynności fotografii design zasobu %d: %w",
			zasobID, err)
	}
	defer wiersze.Close()

	czynnosci := []CzynnoscFotografiiDesignu{}
	for wiersze.Next() {
		czynnosc, err := odczytajCzynnoscFotografiiDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz czynności fotografii design: %w", err)
		}
		czynnosci = append(czynnosci, czynnosc)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt czynności fotografii design: %w", err)
	}
	return czynnosci, nil
}

// CzynnoscFotografiiDesignuWyniku oddaje czynność, z której powstał wskazany
// zasób. Brak wiersza wraca jako ErrBrakWiersza — warstwa wyżej odróżnia „zasób
// nie powstał z obróbki" od usterki odczytu.
func (r *repozytoriumDesignu) CzynnoscFotografiiDesignuWyniku(ctx context.Context,
	zasobID int64) (CzynnoscFotografiiDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, czynnoscFotografiiDesignuZasobu)
	if err != nil {
		return CzynnoscFotografiiDesignu{}, err
	}
	czynnosc, err := odczytajCzynnoscFotografiiDesignu(polecenie.QueryRowContext(ctx, zasobID))
	if errors.Is(err, sql.ErrNoRows) {
		return CzynnoscFotografiiDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return CzynnoscFotografiiDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz czynności fotografii design zasobu %d: %w", zasobID, err)
	}
	return czynnosc, nil
}

// ZapiszNastaweFotografiiDesignu zakłada nastawę albo nadpisuje zastaną po
// identyfikatorze zewnętrznym.
func (r *repozytoriumDesignu) ZapiszNastaweFotografiiDesignu(ctx context.Context,
	nastawa NastawaFotografiiDesignu) (NastawaFotografiiDesignu, error) {

	if nastawa.Kod == "" {
		return NastawaFotografiiDesignu{}, fmt.Errorf(
			"dane: nastawa fotografii design bez identyfikatora")
	}
	if nastawa.Okno == "" {
		return NastawaFotografiiDesignu{}, fmt.Errorf(
			"dane: nastawa fotografii design %q bez okna", nastawa.Kod)
	}
	if nastawa.Nazwa == "" {
		return NastawaFotografiiDesignu{}, fmt.Errorf(
			"dane: nastawa fotografii design %q bez nazwy", nastawa.Kod)
	}
	if nastawa.CzynnosciJSON == "" {
		return NastawaFotografiiDesignu{}, fmt.Errorf(
			"dane: nastawa fotografii design %q bez czynności", nastawa.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszNastaweFotografiiDesignu)
	if err != nil {
		return NastawaFotografiiDesignu{}, err
	}
	if _, err := polecenie.ExecContext(ctx, nastawa.Kod, nastawa.Okno, nastawa.Nazwa,
		nastawa.CzynnosciJSON); err != nil {

		return NastawaFotografiiDesignu{}, fmt.Errorf(
			"dane: nie można zapisać nastawy fotografii design %q: %w", nastawa.Kod, err)
	}
	return r.NastawaFotografiiDesignuPoKodzie(ctx, nastawa.Kod)
}

// NastawaFotografiiDesignuPoKodzie czyta nastawę po identyfikatorze zewnętrznym wprost z bazy danych repozytorium.
func (r *repozytoriumDesignu) NastawaFotografiiDesignuPoKodzie(ctx context.Context,
	kod string) (NastawaFotografiiDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzNastaweFotografiiDesignu)
	if err != nil {
		return NastawaFotografiiDesignu{}, err
	}
	nastawa, err := odczytajNastaweFotografiiDesignu(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return NastawaFotografiiDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return NastawaFotografiiDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz nastawy fotografii design %q: %w", kod, err)
	}
	return nastawa, nil
}

// NastawyFotografiiDesignu oddaje nastawy okna w kolejności ich nazw wprost z bazy danych repozytorium.
func (r *repozytoriumDesignu) NastawyFotografiiDesignu(ctx context.Context,
	okno string) ([]NastawaFotografiiDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaNastawFotografiiDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać nastaw fotografii design okna %q: %w",
			okno, err)
	}
	defer wiersze.Close()

	nastawy := []NastawaFotografiiDesignu{}
	for wiersze.Next() {
		nastawa, err := odczytajNastaweFotografiiDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz nastawy fotografii design: %w", err)
		}
		nastawy = append(nastawy, nastawa)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt nastaw fotografii design: %w", err)
	}
	return nastawy, nil
}

// odczytajCzynnoscFotografiiDesignu składa strukturę czynności wprost z jednego wiersza wyniku zapytania SQL.
func odczytajCzynnoscFotografiiDesignu(wiersz skaner) (CzynnoscFotografiiDesignu, error) {
	var czynnosc CzynnoscFotografiiDesignu
	var zrodlo sql.NullInt64
	var nastawy, policzonePrzez sql.NullString
	err := wiersz.Scan(&czynnosc.ID, &czynnosc.ZasobID, &zrodlo, &czynnosc.Komenda,
		&nastawy, &policzonePrzez, &czynnosc.Utworzono)
	if err != nil {
		return CzynnoscFotografiiDesignu{}, err
	}
	if zrodlo.Valid {
		wartosc := zrodlo.Int64
		czynnosc.ZasobZrodlaID = &wartosc
	}
	czynnosc.NastawyJSON = tekstZKolumny(nastawy)
	czynnosc.PoliczonePrzez = tekstZKolumny(policzonePrzez)
	return czynnosc, nil
}

// odczytajNastaweFotografiiDesignu składa strukturę nastawy wprost z jednego wiersza wyniku zapytania SQL.
func odczytajNastaweFotografiiDesignu(wiersz skaner) (NastawaFotografiiDesignu, error) {
	var nastawa NastawaFotografiiDesignu
	err := wiersz.Scan(&nastawa.ID, &nastawa.Kod, &nastawa.Okno, &nastawa.Nazwa,
		&nastawa.CzynnosciJSON, &nastawa.Zaktualizowano)
	if err != nil {
		return NastawaFotografiiDesignu{}, err
	}
	return nastawa, nil
}
