// Odpowiedzialność pliku: tury debaty (tabela `debata_tura`) i wypowiedzi w nich
// (tabela `debata_wypowiedz`) — `store/migracja_044_roundtable.sql`.
//
// NUMER TURY NADAJE BAZA, NIE RDZEŃ. Numer powstaje jako „największy dotychczas
// w tym oknie plus jeden” WEWNĄTRZ transakcji zakładającej wiersz. Gdyby liczył
// go rdzeń, dwie tury uruchomione w tej samej chwili z dwóch urządzeń tego
// konta dostałyby ten sam numer, a więz UNIQUE(okno, numer) odrzuciłby drugą.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	kolumnyTury = `identyfikator_zewnetrzny, okno, numer, zagadnienie, pytanie, format,
	               stan, granica_tur, tura_nadrzedna, granica_czasu_ms, granica_znakow,
	               anonimowa, rozpoczeto, zamknieto`

	zalozTureDebaty = `INSERT INTO debata_tura
	                   (identyfikator_zewnetrzny, okno, numer, zagadnienie, pytanie,
	                    format, stan, granica_tur, tura_nadrzedna, granica_czasu_ms,
	                    granica_znakow, anonimowa)
	                   VALUES (?, ?,
	                           (SELECT COALESCE(MAX(numer), 0) + 1 FROM debata_tura WHERE okno = ?),
	                           ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	zmienTureDebaty = `UPDATE debata_tura
	                   SET zagadnienie = ?, stan = ?, zamknieto = ?
	                   WHERE identyfikator_zewnetrzny = ?`

	pobierzTure = `SELECT ` + kolumnyTury + `
	               FROM debata_tura WHERE identyfikator_zewnetrzny = ?`

	// Zero w granicy znaczy „bez granicy”, więc jedno przygotowane zapytanie
	// obsługuje wykaz pełny i wykaz przycięty.
	pobierzTury = `SELECT ` + kolumnyTury + `
	               FROM debata_tura WHERE okno = ?
	               ORDER BY numer DESC LIMIT (CASE WHEN ? > 0 THEN ? ELSE -1 END)`

	zapiszWypowiedzDebaty = `INSERT INTO debata_wypowiedz
	                         (identyfikator_zewnetrzny, tura_id, uczestnik, tresc,
	                          odpowiedz_na, akt_mowy, pewnosc, redakcja)
	                         SELECT ?, t.id, ?, ?, ?, ?, ?, ?
	                           FROM debata_tura t
	                          WHERE t.identyfikator_zewnetrzny = ?`

	uzupelnijWypowiedzDebaty = `UPDATE debata_wypowiedz SET tresc = ?
	                            WHERE identyfikator_zewnetrzny = ?`

	// ZastapWypowiedz podnosi numer redakcji razem z treścią: regeneracja jest
	// zastąpieniem, a nie dopisaniem, więc licznik redakcji jest jedynym śladem
	// tego, że wypowiedź już raz padła inaczej.
	zastapWypowiedzDebaty = `UPDATE debata_wypowiedz
	                         SET tresc = ?, redakcja = redakcja + 1
	                         WHERE identyfikator_zewnetrzny = ?`

	oznaczWypowiedzDebaty = `UPDATE debata_wypowiedz SET akt_mowy = ?, pewnosc = ?
	                         WHERE identyfikator_zewnetrzny = ?`

	kolumnyWypowiedziDebaty = `w.identyfikator_zewnetrzny, t.identyfikator_zewnetrzny,
	                     w.uczestnik, w.tresc, w.odpowiedz_na, w.akt_mowy, w.pewnosc,
	                     w.redakcja, w.utworzono`

	pobierzWypowiedzi = `SELECT ` + kolumnyWypowiedziDebaty + `
	                       FROM debata_wypowiedz w
	                       JOIN debata_tura t ON t.id = w.tura_id
	                      WHERE t.identyfikator_zewnetrzny = ?
	                      ORDER BY w.id ASC`

	pobierzWypowiedzDebaty = `SELECT ` + kolumnyWypowiedziDebaty + `
	                      FROM debata_wypowiedz w
	                      JOIN debata_tura t ON t.id = w.tura_id
	                     WHERE w.identyfikator_zewnetrzny = ?`

	pobierzWypowiedziOknaDebaty = `SELECT ` + kolumnyWypowiedziDebaty + `
	                           FROM debata_wypowiedz w
	                           JOIN debata_tura t ON t.id = w.tura_id
	                          WHERE t.okno = ?
	                          ORDER BY t.numer ASC, w.id ASC`
)

// ZalozTure zakłada turę o kolejnym numerze w oknie i oddaje ją po zapisie.
func (r *repozytoriumRoundtable) ZalozTure(ctx context.Context, tura TuraDebaty) (TuraDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zalozTureDebaty)
	if err != nil {
		return TuraDebaty{}, err
	}
	if _, err := polecenie.ExecContext(ctx, tura.Kod, tura.Okno, tura.Okno, tura.Zagadnienie,
		tura.Pytanie, tura.Format, tura.Stan, tura.GranicaTur, tura.TuraNadrzedna,
		tura.GranicaCzasuMs, tura.GranicaZnakow, tura.Anonimowa); err != nil {
		return TuraDebaty{}, fmt.Errorf("dane: nie można założyć tury debaty %q: %w", tura.Kod, err)
	}
	return r.Tura(ctx, tura.Kod)
}

// ZmienTure zapisuje zagadnienie, stan i chwilę zamknięcia tury.
func (r *repozytoriumRoundtable) ZmienTure(ctx context.Context, tura TuraDebaty) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zmienTureDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, tura.Zagadnienie, tura.Stan, tura.Zamknieto, tura.Kod)
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić tury debaty %q: %w", tura.Kod, err)
	}
	return trafienieDebaty(wynik)
}

// Tura zwraca turę po kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumRoundtable) Tura(ctx context.Context, kod string) (TuraDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzTure)
	if err != nil {
		return TuraDebaty{}, err
	}
	tura, err := odczytajTure(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return TuraDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return TuraDebaty{}, fmt.Errorf("dane: nieczytelny wiersz tury debaty %q: %w", kod, err)
	}
	return tura, nil
}

// Tury zwraca tury okna od najnowszej. Zero w granicy zwraca wszystkie.
func (r *repozytoriumRoundtable) Tury(ctx context.Context, okno string, limit int) ([]TuraDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzTury)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, limit, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać tur debaty okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	tury := make([]TuraDebaty, 0, 8)
	for wiersze.Next() {
		tura, err := odczytajTure(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz tury debaty: %w", err)
		}
		tury = append(tury, tura)
	}
	return tury, wiersze.Err()
}

// ZapiszWypowiedz dopisuje wypowiedź do tury. Wskazanie tury nieistniejącej nie
// dopisuje niczego i wraca jako ErrBrakWiersza — cichego zapisu w próżnię nie ma.
func (r *repozytoriumRoundtable) ZapiszWypowiedz(ctx context.Context,
	wypowiedz WypowiedzDebaty) (WypowiedzDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWypowiedzDebaty)
	if err != nil {
		return WypowiedzDebaty{}, err
	}
	if wypowiedz.Redakcja <= 0 {
		wypowiedz.Redakcja = 1
	}
	if wypowiedz.Pewnosc == 0 {
		wypowiedz.Pewnosc = -1 // zero jest deklaracją, brak deklaracji jest ujemny
	}
	wynik, err := polecenie.ExecContext(ctx, wypowiedz.Kod, wypowiedz.Uczestnik,
		wypowiedz.Tresc, wypowiedz.OdpowiedzNa, wypowiedz.AktMowy, wypowiedz.Pewnosc,
		wypowiedz.Redakcja, wypowiedz.TuraKod)
	if err != nil {
		return WypowiedzDebaty{}, fmt.Errorf("dane: nie można zapisać wypowiedzi %q: %w",
			wypowiedz.Kod, err)
	}
	if err := trafienieDebaty(wynik); err != nil {
		return WypowiedzDebaty{}, err
	}
	return wypowiedz, nil
}

// UzupelnijWypowiedz zapisuje treść wypowiedzi po zamknięciu strumienia modelu.
func (r *repozytoriumRoundtable) UzupelnijWypowiedz(ctx context.Context, kod, tresc string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, uzupelnijWypowiedzDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, tresc, kod)
	if err != nil {
		return fmt.Errorf("dane: nie można uzupełnić wypowiedzi %q: %w", kod, err)
	}
	return trafienieDebaty(wynik)
}

// Wypowiedzi zwraca zapis tury w kolejności powstania wypowiedzi.
func (r *repozytoriumRoundtable) Wypowiedzi(ctx context.Context, turaKod string) ([]WypowiedzDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWypowiedzi)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, turaKod)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wypowiedzi tury %q: %w", turaKod, err)
	}
	defer wiersze.Close()

	return zbierzWypowiedziDebaty(wiersze)
}

// WypowiedziOkna zwraca zapis całej debaty okna w porządku tur i wypowiedzi.
// Bez niego każda analiza całej debaty czytałaby tury osobno i składała je
// zapytaniem na turę — czyli tym samym odczytem rozbitym na N wywołań.
func (r *repozytoriumRoundtable) WypowiedziOkna(ctx context.Context,
	okno string) ([]WypowiedzDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWypowiedziOknaDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wypowiedzi okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	return zbierzWypowiedziDebaty(wiersze)
}

// Wypowiedz zwraca jedną wypowiedź po kodzie.
func (r *repozytoriumRoundtable) Wypowiedz(ctx context.Context, kod string) (WypowiedzDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWypowiedzDebaty)
	if err != nil {
		return WypowiedzDebaty{}, err
	}
	wypowiedz, err := odczytajWypowiedzDebaty(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return WypowiedzDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return WypowiedzDebaty{}, fmt.Errorf("dane: nieczytelny wiersz wypowiedzi %q: %w", kod, err)
	}
	return wypowiedz, nil
}

// ZastapWypowiedz podmienia treść wypowiedzi i podnosi numer redakcji.
func (r *repozytoriumRoundtable) ZastapWypowiedz(ctx context.Context, kod, tresc string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zastapWypowiedzDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, tresc, kod)
	if err != nil {
		return fmt.Errorf("dane: nie można zastąpić wypowiedzi %q: %w", kod, err)
	}
	return trafienieDebaty(wynik)
}

// OznaczWypowiedz zapisuje akt mowy i pewność rozpoznane analizą.
func (r *repozytoriumRoundtable) OznaczWypowiedz(ctx context.Context,
	kod, aktMowy string, pewnosc float64) error {

	polecenie, err := r.zapytania.przygotuj(ctx, oznaczWypowiedzDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, aktMowy, pewnosc, kod)
	if err != nil {
		return fmt.Errorf("dane: nie można oznaczyć wypowiedzi %q: %w", kod, err)
	}
	return trafienieDebaty(wynik)
}

// zbierzWypowiedziDebaty czyta wiersze wypowiedzi jednym przebiegiem.
func zbierzWypowiedziDebaty(wiersze *sql.Rows) ([]WypowiedzDebaty, error) {
	wypowiedzi := make([]WypowiedzDebaty, 0, 8)
	for wiersze.Next() {
		wypowiedz, err := odczytajWypowiedzDebaty(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wypowiedzi debaty: %w", err)
		}
		wypowiedzi = append(wypowiedzi, wypowiedz)
	}
	return wypowiedzi, wiersze.Err()
}

// odczytajWypowiedzDebaty składa wypowiedź z jednego wiersza wyniku.
func odczytajWypowiedzDebaty(wiersz interface{ Scan(...any) error }) (WypowiedzDebaty, error) {
	var wypowiedz WypowiedzDebaty
	err := wiersz.Scan(&wypowiedz.Kod, &wypowiedz.TuraKod, &wypowiedz.Uczestnik,
		&wypowiedz.Tresc, &wypowiedz.OdpowiedzNa, &wypowiedz.AktMowy, &wypowiedz.Pewnosc,
		&wypowiedz.Redakcja, &wypowiedz.Utworzono)
	return wypowiedz, err
}

// odczytajTure składa turę z jednego wiersza wyniku.
func odczytajTure(wiersz interface{ Scan(...any) error }) (TuraDebaty, error) {
	var tura TuraDebaty
	err := wiersz.Scan(&tura.Kod, &tura.Okno, &tura.Numer, &tura.Zagadnienie, &tura.Pytanie,
		&tura.Format, &tura.Stan, &tura.GranicaTur, &tura.TuraNadrzedna, &tura.GranicaCzasuMs,
		&tura.GranicaZnakow, &tura.Anonimowa, &tura.Rozpoczeto, &tura.Zamknieto)
	return tura, err
}
