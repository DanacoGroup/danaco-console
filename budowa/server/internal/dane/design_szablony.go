// Repozytorium obsługuje szablony promptu w tabeli `szablon_promptu_design` z migracji
// 231 oraz historię wydanych promptów w tabeli `prompt_design` wzbogaconej migracją 232.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// SzablonPromptuDesignu to wiersz tabeli `szablon_promptu_design`. Pola promptu
// powtarzają `PromptDesignu`, bo opisują ten sam kształt kontraktu; wskaźnik
// znaczy „nieustawione", nie „puste".
type SzablonPromptuDesignu struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Temat          string
	Styl           *string
	Kompozycja     *string
	Oswietlenie    *string
	Paleta         *string
	Proporcje      *string
	Wykluczenia    *string
	Ziarno         *int
	Warianty       *int
	Silnik         *string
	Kreatywnosc    *float64
	Zaktualizowano string
}

const (
	kolumnySzablonuPromptuDesignu = `id, identyfikator_zewnetrzny, okno, nazwa, temat, styl,
	                                 kompozycja, oswietlenie, paleta, proporcje_kadru,
	                                 wykluczenia, ziarno, warianty, silnik, kreatywnosc,
	                                 zaktualizowano`

	zapiszSzablonPromptuDesignu = `INSERT INTO szablon_promptu_design
	                               (identyfikator_zewnetrzny, okno, nazwa, temat, styl, kompozycja,
	                                oswietlenie, paleta, proporcje_kadru, wykluczenia, ziarno,
	                                warianty, silnik, kreatywnosc, zaktualizowano, konto_id)
	                               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	                                       strftime('%Y-%m-%dT%H:%M:%fZ','now'),
	                                       ` + WskazanieKonta + `)
	                               ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                   nazwa = excluded.nazwa,
	                                   temat = excluded.temat,
	                                   styl = excluded.styl,
	                                   kompozycja = excluded.kompozycja,
	                                   oswietlenie = excluded.oswietlenie,
	                                   paleta = excluded.paleta,
	                                   proporcje_kadru = excluded.proporcje_kadru,
	                                   wykluczenia = excluded.wykluczenia,
	                                   ziarno = excluded.ziarno,
	                                   warianty = excluded.warianty,
	                                   silnik = excluded.silnik,
	                                   kreatywnosc = excluded.kreatywnosc,
	                                   zaktualizowano = excluded.zaktualizowano
	                               WHERE ` + WarunekKonta

	pobierzSzablonPromptuDesignu = `SELECT ` + kolumnySzablonuPromptuDesignu +
		` FROM szablon_promptu_design
		  WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaSzablonowPromptuDesignu = `SELECT ` + kolumnySzablonuPromptuDesignu +
		` FROM szablon_promptu_design WHERE okno = ? AND ` + WarunekKonta + `
		  ORDER BY zaktualizowano DESC, id DESC`

	listaPromptowOknaDesignu = `SELECT ` + kolumnyPromptuDesign + ` FROM prompt_design
	                            WHERE okno = ? ORDER BY utworzono DESC, id DESC LIMIT ?`

	liczbaPromptowOknaDesignu = `SELECT COUNT(*) FROM prompt_design WHERE okno = ?`

	listaZasobowPromptuDesignu = `SELECT identyfikator_zewnetrzny FROM zasob_design
	                              WHERE prompt_id = ? AND ` + WarunekKonta + `
	                              ORDER BY utworzono, id`
)

// ZapiszSzablonPromptuDesignu zakłada szablon albo nadpisuje zastany po
// identyfikatorze zewnętrznym i oddaje stan po zapisie.
func (r *repozytoriumDesignu) ZapiszSzablonPromptuDesignu(ctx context.Context,
	szablon SzablonPromptuDesignu) (SzablonPromptuDesignu, error) {

	if szablon.Kod == "" {
		return SzablonPromptuDesignu{}, fmt.Errorf("dane: szablon promptu design bez identyfikatora")
	}
	if szablon.Okno == "" {
		return SzablonPromptuDesignu{}, fmt.Errorf("dane: szablon promptu design %q bez okna", szablon.Kod)
	}
	if szablon.Nazwa == "" {
		return SzablonPromptuDesignu{}, fmt.Errorf("dane: szablon promptu design %q bez nazwy", szablon.Kod)
	}
	if szablon.Temat == "" {
		return SzablonPromptuDesignu{}, fmt.Errorf("dane: szablon promptu design %q bez tematu", szablon.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSzablonPromptuDesignu)
	if err != nil {
		return SzablonPromptuDesignu{}, err
	}
	// Kontrakt niesie Ziarno i Warianty jako wskaźnik na int, kolumna SQL wymaga typu int64.
	var ziarno, warianty, kreatywnosc any
	if szablon.Ziarno != nil {
		ziarno = int64(*szablon.Ziarno)
	}
	if szablon.Warianty != nil {
		warianty = int64(*szablon.Warianty)
	}
	if szablon.Kreatywnosc != nil {
		kreatywnosc = *szablon.Kreatywnosc
	}
	wynik, err := polecenie.ExecContext(ctx, szablon.Kod, szablon.Okno, szablon.Nazwa, szablon.Temat,
		tekstDoKolumny(szablon.Styl), tekstDoKolumny(szablon.Kompozycja),
		tekstDoKolumny(szablon.Oswietlenie), tekstDoKolumny(szablon.Paleta),
		tekstDoKolumny(szablon.Proporcje), tekstDoKolumny(szablon.Wykluczenia),
		ziarno, warianty, tekstDoKolumny(szablon.Silnik), kreatywnosc,
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return SzablonPromptuDesignu{}, fmt.Errorf(
			"dane: nie można zapisać szablonu promptu design %q: %w", szablon.Kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return SzablonPromptuDesignu{}, fmt.Errorf(
			"dane: nieznana liczba zapisanych szablonów promptu design: %w", err)
	}
	if zmienione == 0 {
		return SzablonPromptuDesignu{}, fmt.Errorf(
			"dane: szablon promptu design %q należy do innego konta: %w",
			szablon.Kod, ErrKolizjaWiersza)
	}

	polecenieOdczytu, err := r.zapytania.przygotuj(ctx, pobierzSzablonPromptuDesignu)
	if err != nil {
		return SzablonPromptuDesignu{}, err
	}
	zapisany, err := odczytajSzablonPromptuDesignu(
		polecenieOdczytu.QueryRowContext(ctx, szablon.Kod, KontoOperatora(ctx)))
	if err != nil {
		return SzablonPromptuDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz szablonu promptu design %q: %w", szablon.Kod, err)
	}
	return zapisany, nil
}

// SzablonyPromptuDesignu zwraca szablony promptu przypisane do okna, uporządkowane od
// ostatnio zmienianego.
func (r *repozytoriumDesignu) SzablonyPromptuDesignu(ctx context.Context,
	okno string) ([]SzablonPromptuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaSzablonowPromptuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać szablonów promptu design okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []SzablonPromptuDesignu{}
	for wiersze.Next() {
		szablon, err := odczytajSzablonPromptuDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz szablonu promptu design okna %q: %w", okno, err)
		}
		lista = append(lista, szablon)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt szablonów promptu design okna %q: %w", okno, err)
	}
	return lista, nil
}

// PromptyOknaDesignu zwraca prompty wydane w oknie (od najświeższego) wraz
// z liczbą wszystkich promptów okna — wykaz bywa przycięty limitem, a licznik
// nie, więc panel pokazuje „X z Y" bez drugiego zapytania.
func (r *repozytoriumDesignu) PromptyOknaDesignu(ctx context.Context,
	okno string, limit int) ([]PromptDesignu, int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPromptowOknaDesignu)
	if err != nil {
		return nil, 0, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, granicaWykazu(limit))
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać promptów design okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []PromptDesignu{}
	for wiersze.Next() {
		prompt, err := odczytajPromptDesign(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz promptu design okna %q: %w", okno, err)
		}
		lista = append(lista, prompt)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt promptów design okna %q: %w", okno, err)
	}

	polecenieLiczby, err := r.zapytania.przygotuj(ctx, liczbaPromptowOknaDesignu)
	if err != nil {
		return nil, 0, err
	}
	var razem int
	if err := polecenieLiczby.QueryRowContext(ctx, okno).Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć promptów design okna %q: %w", okno, err)
	}
	return lista, razem, nil
}

// ZasobyPromptuDesignu zwraca identyfikatory zewnętrzne zasobów powstałych
// z promptu, w kolejności powstawania.
func (r *repozytoriumDesignu) ZasobyPromptuDesignu(ctx context.Context,
	promptID int64) ([]string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZasobowPromptuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, promptID, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zasobów promptu design %d: %w", promptID, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var kod string
		if err := wiersze.Scan(&kod); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zasobu promptu design %d: %w", promptID, err)
		}
		lista = append(lista, kod)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zasobów promptu design %d: %w", promptID, err)
	}
	return lista, nil
}

// odczytajSzablonPromptuDesignu składa strukturę SzablonPromptuDesignu z jednego wiersza
// wyniku zapytania SQL.
func odczytajSzablonPromptuDesignu(wiersz skaner) (SzablonPromptuDesignu, error) {
	var szablon SzablonPromptuDesignu
	var styl, kompozycja, oswietlenie, paleta, proporcje, wykluczenia, silnik sql.NullString
	var ziarno, warianty sql.NullInt64
	var kreatywnosc sql.NullFloat64
	err := wiersz.Scan(&szablon.ID, &szablon.Kod, &szablon.Okno, &szablon.Nazwa, &szablon.Temat,
		&styl, &kompozycja, &oswietlenie, &paleta, &proporcje, &wykluczenia, &ziarno,
		&warianty, &silnik, &kreatywnosc, &szablon.Zaktualizowano)
	if err != nil {
		return SzablonPromptuDesignu{}, err
	}
	szablon.Styl = tekstZKolumny(styl)
	szablon.Kompozycja = tekstZKolumny(kompozycja)
	szablon.Oswietlenie = tekstZKolumny(oswietlenie)
	szablon.Paleta = tekstZKolumny(paleta)
	szablon.Proporcje = tekstZKolumny(proporcje)
	szablon.Wykluczenia = tekstZKolumny(wykluczenia)
	szablon.Silnik = tekstZKolumny(silnik)
	if ziarno.Valid {
		wartosc := int(ziarno.Int64)
		szablon.Ziarno = &wartosc
	}
	if warianty.Valid {
		wartosc := int(warianty.Int64)
		szablon.Warianty = &wartosc
	}
	szablon.Kreatywnosc = liczbaRzeczywistaZKolumny(kreatywnosc)
	return szablon, nil
}
