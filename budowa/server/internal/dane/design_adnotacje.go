// Plik prowadzi obszar adnotacji kompozycji Design Board, część RepozytoriumDesignu; adnotacja przeżywa zapis układu
// warstw, bo pole adnotacji jedzie osobno od układu; wątek powstaje przez pole nadrzędnej, a porządek odczytu jest chronologiczny.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// AdnotacjaDesignu to wiersz tabeli `adnotacja_kompozycji_design`; WarstwaID i NadrzednaID są identyfikatorami zewnętrznymi, nie kluczami wierszy.
type AdnotacjaDesignu struct {
	ID           int64
	Kod          string
	KompozycjaID int64
	WarstwaID    *string
	NadrzednaID  *string
	Autor        *string
	Tresc        string
	Zamknieta    bool
	Utworzono    string
}

const (
	kolumnyAdnotacjiDesignu = `id, identyfikator_zewnetrzny, kompozycja_id, warstwa_id,
	                           nadrzedna_id, autor, tresc, zamknieta, utworzono`

	zapiszAdnotacjeDesignu = `INSERT INTO adnotacja_kompozycji_design
	                          (identyfikator_zewnetrzny, kompozycja_id, warstwa_id, nadrzedna_id,
	                           autor, tresc, zamknieta)
	                          VALUES (?, ?, ?, ?, ?, ?, ?)
	                          ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                              warstwa_id = excluded.warstwa_id,
	                              nadrzedna_id = excluded.nadrzedna_id,
	                              tresc = excluded.tresc,
	                              zamknieta = excluded.zamknieta`

	pobierzAdnotacjeDesignu = `SELECT ` + kolumnyAdnotacjiDesignu +
		` FROM adnotacja_kompozycji_design WHERE identyfikator_zewnetrzny = ?`

	listaAdnotacjiDesignu = `SELECT ` + kolumnyAdnotacjiDesignu +
		` FROM adnotacja_kompozycji_design WHERE kompozycja_id = ? ORDER BY id`

	listaAdnotacjiDesignuOtwartych = `SELECT ` + kolumnyAdnotacjiDesignu +
		` FROM adnotacja_kompozycji_design WHERE kompozycja_id = ? AND zamknieta = 0
		  ORDER BY id`
)

// ZapiszAdnotacjeDesignu zakłada adnotację albo nadpisuje zastaną po identyfikatorze zewnętrznym; autor zapisuje się wyłącznie przy założeniu.
func (r *repozytoriumDesignu) ZapiszAdnotacjeDesignu(ctx context.Context,
	adnotacja AdnotacjaDesignu) (AdnotacjaDesignu, error) {

	if adnotacja.Kod == "" {
		return AdnotacjaDesignu{}, fmt.Errorf("dane: adnotacja design bez identyfikatora")
	}
	if adnotacja.KompozycjaID == 0 {
		return AdnotacjaDesignu{}, fmt.Errorf("dane: adnotacja design %q bez kompozycji", adnotacja.Kod)
	}
	if adnotacja.Tresc == "" {
		return AdnotacjaDesignu{}, fmt.Errorf("dane: adnotacja design %q bez treści", adnotacja.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszAdnotacjeDesignu)
	if err != nil {
		return AdnotacjaDesignu{}, err
	}
	_, err = polecenie.ExecContext(ctx, adnotacja.Kod, adnotacja.KompozycjaID,
		tekstDoKolumny(adnotacja.WarstwaID), tekstDoKolumny(adnotacja.NadrzednaID),
		tekstDoKolumny(adnotacja.Autor), adnotacja.Tresc, liczbaLogiczna(adnotacja.Zamknieta))
	if err != nil {
		return AdnotacjaDesignu{}, fmt.Errorf("dane: nie można zapisać adnotacji design %q: %w",
			adnotacja.Kod, err)
	}
	return r.AdnotacjaDesignuPoKodzie(ctx, adnotacja.Kod)
}

// AdnotacjaDesignuPoKodzie zwraca adnotację o wskazanym identyfikatorze
// zewnętrznym. Brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumDesignu) AdnotacjaDesignuPoKodzie(ctx context.Context,
	kod string) (AdnotacjaDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzAdnotacjeDesignu)
	if err != nil {
		return AdnotacjaDesignu{}, err
	}
	adnotacja, err := odczytajAdnotacjeDesignu(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return AdnotacjaDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return AdnotacjaDesignu{}, fmt.Errorf("dane: nieczytelny wiersz adnotacji design %q: %w", kod, err)
	}
	return adnotacja, nil
}

// AdnotacjeKompozycjiDesignu zwraca adnotacje kompozycji w kolejności
// powstawania; zawężenie oddaje wyłącznie wątki niezamknięte.
func (r *repozytoriumDesignu) AdnotacjeKompozycjiDesignu(ctx context.Context,
	kompozycjaID int64, tylkoOtwarte bool) ([]AdnotacjaDesignu, error) {

	zapytanie := listaAdnotacjiDesignu
	if tylkoOtwarte {
		zapytanie = listaAdnotacjiDesignuOtwartych
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kompozycjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać adnotacji kompozycji design %d: %w",
			kompozycjaID, err)
	}
	defer wiersze.Close()

	lista := []AdnotacjaDesignu{}
	for wiersze.Next() {
		adnotacja, err := odczytajAdnotacjeDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz adnotacji kompozycji design %d: %w",
				kompozycjaID, err)
		}
		lista = append(lista, adnotacja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt adnotacji kompozycji design %d: %w",
			kompozycjaID, err)
	}
	return lista, nil
}

// odczytajAdnotacjeDesignu składa strukturę adnotacji wprost z jednego wiersza wyniku zapytania do bazy SQL.
func odczytajAdnotacjeDesignu(wiersz skaner) (AdnotacjaDesignu, error) {
	var adnotacja AdnotacjaDesignu
	var warstwaID, nadrzednaID, autor sql.NullString
	var zamknieta int
	err := wiersz.Scan(&adnotacja.ID, &adnotacja.Kod, &adnotacja.KompozycjaID, &warstwaID,
		&nadrzednaID, &autor, &adnotacja.Tresc, &zamknieta, &adnotacja.Utworzono)
	if err != nil {
		return AdnotacjaDesignu{}, err
	}
	adnotacja.WarstwaID = tekstZKolumny(warstwaID)
	adnotacja.NadrzednaID = tekstZKolumny(nadrzednaID)
	adnotacja.Autor = tekstZKolumny(autor)
	adnotacja.Zamknieta = zamknieta == 1
	return adnotacja, nil
}
