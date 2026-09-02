// Odpowiedzialność pliku: stanowisko końcowe debaty, trwałość okna Consensus Panel; wersja rośnie przy zmianie treści, nie przy każdym odczycie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	kolumnyStanowiska = `identyfikator_zewnetrzny, okno, tura, tresc, wersja, zaktualizowano,
	                     redagowane, zaakceptowane, kontekst, warianty, konsekwencje, tury`

	// Warunek redagowane = 0 broni redakcji Operatora: od chwili nadania własnej treści, złożenie z zapisu tur nie ma prawa jej nadpisać.
	zapiszStanowiskoDebaty = `INSERT INTO debata_stanowisko
	                          (identyfikator_zewnetrzny, okno, tura, tresc, konto_id)
	                          VALUES (?, ?, ?, ?, ` + WskazanieKonta + `)
	                          ON CONFLICT(okno, tura) DO UPDATE SET
	                              tresc = excluded.tresc,
	                              wersja = debata_stanowisko.wersja + 1,
	                              zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE debata_stanowisko.redagowane = 0
	                            AND COALESCE(debata_stanowisko.tresc, '') <> COALESCE(excluded.tresc, '')
	                            AND ` + WarunekKonta

	// Redakcja Operatora podnosi wersję zawsze — także wtedy, gdy treść wyszła
	// ta sama. Zapisanie tej samej treści jest czynnością zamierzoną (na przykład
	// samą akceptacją), a wersja liczy redakcje, nie różnice.
	redagujStanowiskoDebaty = `INSERT INTO debata_stanowisko
	                           (identyfikator_zewnetrzny, okno, tura, tresc, redagowane,
	                            zaakceptowane, kontekst, warianty, konsekwencje, tury,
	                            konto_id)
	                           VALUES (?, ?, ?, ?, 1, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                           ON CONFLICT(okno, tura) DO UPDATE SET
	                               tresc = excluded.tresc,
	                               redagowane = 1,
	                               zaakceptowane = excluded.zaakceptowane,
	                               kontekst = excluded.kontekst,
	                               warianty = excluded.warianty,
	                               konsekwencje = excluded.konsekwencje,
	                               tury = excluded.tury,
	                               wersja = debata_stanowisko.wersja + 1,
	                               zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                           WHERE ` + WarunekKonta

	pobierzStanowisko = `SELECT ` + kolumnyStanowiska + `
	                     FROM debata_stanowisko
	                     WHERE okno = ? AND tura = ? AND ` + WarunekKonta

	pobierzStanowiskoDebatyPoKodzie = `SELECT ` + kolumnyStanowiska + `
	                             FROM debata_stanowisko
	                             WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

// ZapiszStanowisko utrwala stanowisko okna albo tury. Treść niezmieniona nie
// podbija wersji ani znacznika czasu; wiersz zostaje jak stał.
func (r *repozytoriumRoundtable) ZapiszStanowisko(ctx context.Context,
	stanowisko StanowiskoDebaty) (StanowiskoDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszStanowiskoDebaty)
	if err != nil {
		return StanowiskoDebaty{}, err
	}
	if _, err := polecenie.ExecContext(ctx, stanowisko.Kod, stanowisko.Okno,
		stanowisko.Tura, stanowisko.Tresc, KontoOperatora(ctx),
		KontoOperatora(ctx)); err != nil {
		return StanowiskoDebaty{}, fmt.Errorf("dane: nie można zapisać stanowiska debaty okna %q: %w",
			stanowisko.Okno, err)
	}
	zapisane, err := r.Stanowisko(ctx, stanowisko.Okno, stanowisko.Tura)
	if errors.Is(err, ErrBrakWiersza) {
		// Zero zmienionych wierszy jest tu stanem normalnym: treść niezmieniona
		// albo wiersz redagowany. Brak wiersza po zapisie normalny nie jest —
		// parę (okno, tura) zajmuje wtedy wiersz konta innego.
		return StanowiskoDebaty{}, fmt.Errorf("dane: stanowisko debaty %q należy do innego konta: %w",
			stanowisko.Kod, ErrKolizjaWiersza)
	}
	return zapisane, err
}

// Stanowisko zwraca stanowisko wskazanego okna albo jednej jego tury; pusta tura znaczy całą tę debatę.
func (r *repozytoriumRoundtable) Stanowisko(ctx context.Context,
	okno, tura string) (StanowiskoDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzStanowisko)
	if err != nil {
		return StanowiskoDebaty{}, err
	}
	stanowisko, err := odczytajStanowiskoDebaty(polecenie.QueryRowContext(ctx, okno, tura,
		KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return StanowiskoDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return StanowiskoDebaty{}, fmt.Errorf("dane: nieczytelny wiersz stanowiska debaty okna %q: %w",
			okno, err)
	}
	return stanowisko, nil
}

// StanowiskoPoKodzie zwraca stanowisko wskazane identyfikatorem. Zdanie odrębne
// i przekazanie wskazują stanowisko kodem, nie parą (okno, tura), więc odczyt
// po kodzie jest tu drogą pierwszą, a nie skrótem.
func (r *repozytoriumRoundtable) StanowiskoPoKodzie(ctx context.Context,
	kod string) (StanowiskoDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzStanowiskoDebatyPoKodzie)
	if err != nil {
		return StanowiskoDebaty{}, err
	}
	stanowisko, err := odczytajStanowiskoDebaty(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return StanowiskoDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return StanowiskoDebaty{}, fmt.Errorf("dane: nieczytelny wiersz stanowiska %q: %w", kod, err)
	}
	return stanowisko, nil
}

// RedagujStanowisko zapisuje treść nadaną przez Operatora wraz z zapisem decyzji
// i znakuje wiersz jako redagowany.
func (r *repozytoriumRoundtable) RedagujStanowisko(ctx context.Context,
	stanowisko StanowiskoDebaty) (StanowiskoDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, redagujStanowiskoDebaty)
	if err != nil {
		return StanowiskoDebaty{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, stanowisko.Kod, stanowisko.Okno, stanowisko.Tura,
		stanowisko.Tresc, stanowisko.Zaakceptowane, stanowisko.Kontekst, stanowisko.Warianty,
		stanowisko.Konsekwencje, stanowisko.Tury, KontoOperatora(ctx),
		KontoOperatora(ctx))
	if err != nil {
		return StanowiskoDebaty{}, fmt.Errorf(
			"dane: nie można zredagować stanowiska debaty okna %q: %w", stanowisko.Okno, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "stanowisko debaty", stanowisko.Kod); err != nil {
		return StanowiskoDebaty{}, err
	}
	return r.Stanowisko(ctx, stanowisko.Okno, stanowisko.Tura)
}

// odczytajStanowiskoDebaty składa stanowisko z jednego wiersza wyniku zapytania, kolumna po kolumnie SQL.
func odczytajStanowiskoDebaty(wiersz interface{ Scan(...any) error }) (StanowiskoDebaty, error) {
	var stanowisko StanowiskoDebaty
	err := wiersz.Scan(&stanowisko.Kod, &stanowisko.Okno, &stanowisko.Tura, &stanowisko.Tresc,
		&stanowisko.Wersja, &stanowisko.Zaktualizowano, &stanowisko.Redagowane,
		&stanowisko.Zaakceptowane, &stanowisko.Kontekst, &stanowisko.Warianty,
		&stanowisko.Konsekwencje, &stanowisko.Tury)
	return stanowisko, err
}
