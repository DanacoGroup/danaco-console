// Błędy modułu Diagnostics: dopisanie wystąpienia, odczyt zawężony filtrem, odczyt po kodach i rozkład stanów.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

const (
	kolumnyBleduDiagnostycznego = `kod, odcisk, tresc, zrodlo, kod_bledu, stan, priorytet,
	                               wystapienia, notatka, kontekst, pierwsze, ostatnie`

	// Wstawienie z rozstrzygnięciem kolizji odcisku w koncie: powtórzenie podnosi licznik i przesuwa chwilę.
	wstawBladDiagnostyczny = `INSERT INTO diagnostyka_blad
	                          (kod, odcisk, tresc, zrodlo, kod_bledu, stan, priorytet,
	                           wystapienia, notatka, kontekst, pierwsze, ostatnie, konto_id)
	                          VALUES (?, ?, ?, ?, ?, ?, ?, 1, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                          ON CONFLICT(odcisk, COALESCE(konto_id, 0)) DO UPDATE SET
	                              wystapienia = wystapienia + 1,
	                              ostatnie    = excluded.ostatnie,
	                              tresc       = excluded.tresc,
	                              kontekst    = excluded.kontekst`

	pobierzBladPoOdcisku = `SELECT ` + kolumnyBleduDiagnostycznego + `
	                        FROM diagnostyka_blad WHERE odcisk = ? AND ` + WarunekKonta

	pobierzBledy = `SELECT ` + kolumnyBleduDiagnostycznego + ` FROM diagnostyka_blad
	                WHERE (? = '' OR stan = ?)
	                  AND (? = '' OR priorytet = ?)
	                  AND (? = 0 OR ostatnie >= ?)
	                  AND (? = 0 OR pierwsze <= ?) AND ` + WarunekKonta + `
	                ORDER BY ostatnie DESC, id DESC
	                LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`

	// policzBledy liczy wiersze filtru bez granicy, by wykaz podał liczbę całkowitą niezależną od długości listy.
	policzBledy = `SELECT COUNT(*) FROM diagnostyka_blad
	                WHERE (? = '' OR stan = ?)
	                  AND (? = '' OR priorytet = ?)
	                  AND (? = 0 OR ostatnie >= ?)
	                  AND (? = 0 OR pierwsze <= ?) AND ` + WarunekKonta

	policzStanyBledow = `SELECT stan, COUNT(*) FROM diagnostyka_blad
	                     WHERE (? = 0 OR ostatnie >= ?) AND (? = 0 OR pierwsze <= ?) AND ` + WarunekKonta + `
	                     GROUP BY stan`
)

// ZapiszBlad dopisuje wystąpienie błędu i zwraca wiersz po zapisie, z podniesionym już licznikiem.
func (r *repozytoriumDiagnostyki) ZapiszBlad(ctx context.Context,
	blad BladDiagnostyczny) (BladDiagnostyczny, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wstawBladDiagnostyczny)
	if err != nil {
		return BladDiagnostyczny{}, err
	}
	_, err = polecenie.ExecContext(ctx, blad.Kod, blad.Odcisk, blad.Tresc,
		tekstDoKolumny(blad.Zrodlo), string(blad.KodBledu), string(blad.Stan),
		string(blad.Priorytet), tekstDoKolumny(blad.Notatka), tekstDoKolumny(blad.Kontekst),
		blad.Pierwsze, blad.Ostatnie, KontoOperatora(ctx))
	if err != nil {
		return BladDiagnostyczny{}, fmt.Errorf("dane: nie można zapisać błędu %q: %w", blad.Odcisk, err)
	}

	odczyt, err := r.zapytania.przygotuj(ctx, pobierzBladPoOdcisku)
	if err != nil {
		return BladDiagnostyczny{}, err
	}
	zapisany, err := odczytajBladDiagnostyczny(odczyt.QueryRowContext(ctx, blad.Odcisk, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return BladDiagnostyczny{}, ErrBrakWiersza
	}
	if err != nil {
		return BladDiagnostyczny{}, fmt.Errorf("dane: nieczytelny wiersz błędu %q: %w", blad.Odcisk, err)
	}
	return zapisany, nil
}

// Bledy zwraca wykaz błędów diagnostycznych zawężony filtrem wyszukiwania, uporządkowany od wystąpienia najświeższego.
func (r *repozytoriumDiagnostyki) Bledy(ctx context.Context, filtr FiltrBledow) ([]BladDiagnostyczny, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzBledy)
	if err != nil {
		return nil, err
	}
	stan, priorytet := string(filtr.Stan), string(filtr.Priorytet)
	granica := filtr.Granica
	if granica <= 0 {
		granica = granicaDziennika
	}
	wiersze, err := polecenie.QueryContext(ctx, stan, stan, priorytet, priorytet,
		filtr.Od, filtr.Od, filtr.Do, filtr.Do, KontoOperatora(ctx), granica, granica)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać błędów diagnostyki: %w", err)
	}
	defer wiersze.Close()
	return zbierzBledy(wiersze)
}

// LiczbaBledow oddaje liczbę błędów pasujących do filtru bez granicy.
func (r *repozytoriumDiagnostyki) LiczbaBledow(ctx context.Context, filtr FiltrBledow) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, policzBledy)
	if err != nil {
		return 0, err
	}
	stan, priorytet := string(filtr.Stan), string(filtr.Priorytet)
	var razem int
	err = polecenie.QueryRowContext(ctx, stan, stan, priorytet, priorytet,
		filtr.Od, filtr.Od, filtr.Do, filtr.Do, KontoOperatora(ctx)).Scan(&razem)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć błędów diagnostyki: %w", err)
	}
	return razem, nil
}

// BledyPoKodach zwraca błędy wskazane wprost; wykaz pusty daje wynik pusty, nie pełny.
func (r *repozytoriumDiagnostyki) BledyPoKodach(ctx context.Context, kody []string) ([]BladDiagnostyczny, error) {
	if len(kody) == 0 {
		return nil, nil
	}
	zapytanie := `SELECT ` + kolumnyBleduDiagnostycznego + ` FROM diagnostyka_blad
	              WHERE kod IN (` + strings.TrimSuffix(strings.Repeat("?,", len(kody)), ",") + `)
	                AND ` + WarunekKonta + `
	              ORDER BY ostatnie DESC, id DESC`
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	argumenty := make([]any, 0, len(kody)+1)
	for _, kod := range kody {
		argumenty = append(argumenty, kod)
	}
	argumenty = append(argumenty, KontoOperatora(ctx))
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wskazanych błędów: %w", err)
	}
	defer wiersze.Close()
	return zbierzBledy(wiersze)
}

// StanyBledow zwraca rozkład liczby błędów diagnostycznych pogrupowanych według stanu, w podanym zakresie czasu.
func (r *repozytoriumDiagnostyki) StanyBledow(ctx context.Context, od, do int64) (LicznikStanow, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, policzStanyBledow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, od, od, do, do, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można policzyć stanów błędów: %w", err)
	}
	defer wiersze.Close()

	licznik := LicznikStanow{}
	for wiersze.Next() {
		var stan shared.DiagnosticErrorStatus
		var liczba int
		if err := wiersze.Scan(&stan, &liczba); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny rozkład stanów: %w", err)
		}
		licznik[stan] = liczba
	}
	return licznik, wiersze.Err()
}

func zbierzBledy(wiersze *sql.Rows) ([]BladDiagnostyczny, error) {
	bledy := make([]BladDiagnostyczny, 0, 16)
	for wiersze.Next() {
		blad, err := odczytajBladDiagnostyczny(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz błędu diagnostyki: %w", err)
		}
		bledy = append(bledy, blad)
	}
	return bledy, wiersze.Err()
}

func odczytajBladDiagnostyczny(s skaner) (BladDiagnostyczny, error) {
	var blad BladDiagnostyczny
	err := s.Scan(&blad.Kod, &blad.Odcisk, &blad.Tresc, &blad.Zrodlo, &blad.KodBledu,
		&blad.Stan, &blad.Priorytet, &blad.Wystapienia, &blad.Notatka, &blad.Kontekst,
		&blad.Pierwsze, &blad.Ostatnie)
	return blad, err
}
