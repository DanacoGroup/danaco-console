// Odpowiedzialność pliku: przebieg wsadu modułu Studio — tabele
// przebieg_wsadu_studio i pozycja_wsadu_studio — jako trwały ślad uruchomienia
// wsadu, zapisywany jedną transakcją obejmującą nagłówek przebiegu wraz
// z kompletem jego pozycji.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// PrzebiegWsaduStudia to wiersz tabeli przebieg_wsadu_studio, niosący nagłówek
// przebiegu wsadu wraz z sumaryczną liczbą pozycji przyjętych i odrzuconych.
type PrzebiegWsaduStudia struct {
	ID            int64
	Kod           string
	Okno          string
	AkcjaID       string
	ParametryJSON *string
	Przyjete      int
	Odrzucone     int
	Utworzono     string
}

// PozycjaWsaduStudia to wiersz tabeli `pozycja_wsadu_studio` — rozstrzygnięcie
// wsadu dla JEDNEGO dokumentu. Stan „odrzucona" niesie powód własnymi słowami,
// bo Operator ma przeczytać, czemu ten dokument wypadł, a nie domyślać się go
// z liczby.
type PozycjaWsaduStudia struct {
	ID            int64
	DokumentKod   string
	Stan          string
	Powod         *string
	PropozycjaKod *string
	Utworzono     string
}

const (
	kolumnyPrzebieguWsaduStudia = `id, identyfikator_zewnetrzny, okno, akcja_id, parametry_json,
	                               przyjete, odrzucone, utworzono`

	zapiszPrzebiegWsaduStudia = `INSERT INTO przebieg_wsadu_studio
	                             (identyfikator_zewnetrzny, okno, akcja_id, parametry_json,
	                              przyjete, odrzucone)
	                             VALUES (?, ?, ?, ?, ?, ?)`

	pobierzPrzebiegWsaduStudia = `SELECT ` + kolumnyPrzebieguWsaduStudia +
		` FROM przebieg_wsadu_studio WHERE identyfikator_zewnetrzny = ?`

	zapiszPozycjeWsaduStudia = `INSERT INTO pozycja_wsadu_studio
	                            (przebieg_id, dokument_kod, stan, powod, propozycja_kod)
	                            VALUES (?, ?, ?, ?, ?)`

	listaPozycjiWsaduStudia = `SELECT p.id, p.dokument_kod, p.stan, p.powod, p.propozycja_kod, p.utworzono
	                           FROM pozycja_wsadu_studio p
	                           JOIN przebieg_wsadu_studio w ON w.id = p.przebieg_id
	                           WHERE w.identyfikator_zewnetrzny = ? ORDER BY p.id`
)

// ZapiszPrzebiegWsadu zakłada przebieg wsadu wraz z kompletem jego pozycji jedną
// transakcją i zwraca zapisany nagłówek przebiegu odczytany po zapisie.
func (r *repozytoriumStudia) ZapiszPrzebiegWsadu(ctx context.Context,
	przebieg PrzebiegWsaduStudia, pozycje []PozycjaWsaduStudia) (PrzebiegWsaduStudia, error) {

	if przebieg.Kod == "" {
		return PrzebiegWsaduStudia{}, fmt.Errorf("dane: przebieg wsadu studio bez identyfikatora")
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszPrzebiegWsaduStudia)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, przebieg.Kod, przebieg.Okno, przebieg.AkcjaID,
			tekstDoKolumny(przebieg.ParametryJSON), przebieg.Przyjete, przebieg.Odrzucone)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać przebiegu wsadu studio %q: %w",
				przebieg.Kod, err)
		}
		przebiegID, err := wynik.LastInsertId()
		if err != nil {
			return fmt.Errorf("dane: nieczytelny klucz przebiegu wsadu studio %q: %w",
				przebieg.Kod, err)
		}

		poleceniePozycji, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszPozycjeWsaduStudia)
		if err != nil {
			return err
		}
		for _, pozycja := range pozycje {
			stan := pozycja.Stan
			if stan == "" {
				stan = "przyjeta"
			}
			_, err = poleceniePozycji.ExecContext(ctx, przebiegID, pozycja.DokumentKod, stan,
				tekstDoKolumny(pozycja.Powod), tekstDoKolumny(pozycja.PropozycjaKod))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać pozycji wsadu studio %q: %w",
					pozycja.DokumentKod, err)
			}
		}
		return nil
	})
	if err != nil {
		return PrzebiegWsaduStudia{}, err
	}
	return r.PrzebiegWsadu(ctx, przebieg.Kod)
}

// PrzebiegWsadu zwraca nagłówek przebiegu wsadu o wskazanym kodzie wraz z liczbą
// pozycji przyjętych i odrzuconych, bez samych pozycji.
func (r *repozytoriumStudia) PrzebiegWsadu(ctx context.Context,
	kod string) (PrzebiegWsaduStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPrzebiegWsaduStudia)
	if err != nil {
		return PrzebiegWsaduStudia{}, err
	}
	var przebieg PrzebiegWsaduStudia
	var parametry sql.NullString
	err = polecenie.QueryRowContext(ctx, kod).Scan(&przebieg.ID, &przebieg.Kod, &przebieg.Okno,
		&przebieg.AkcjaID, &parametry, &przebieg.Przyjete, &przebieg.Odrzucone, &przebieg.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return PrzebiegWsaduStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return PrzebiegWsaduStudia{}, fmt.Errorf("dane: nieczytelny przebieg wsadu studio %q: %w", kod, err)
	}
	przebieg.ParametryJSON = tekstZKolumny(parametry)
	return przebieg, nil
}

// PozycjeWsadu zwraca rozstrzygnięcia przebiegu wsadu dla poszczególnych
// dokumentów w kolejności ich zapisu, wskazanego kodem przebiegu.
func (r *repozytoriumStudia) PozycjeWsadu(ctx context.Context,
	kodPrzebiegu string) ([]PozycjaWsaduStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPozycjiWsaduStudia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodPrzebiegu)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pozycji wsadu studio %q: %w", kodPrzebiegu, err)
	}
	defer wiersze.Close()

	lista := []PozycjaWsaduStudia{}
	for wiersze.Next() {
		var pozycja PozycjaWsaduStudia
		var powod, propozycja sql.NullString
		if err := wiersze.Scan(&pozycja.ID, &pozycja.DokumentKod, &pozycja.Stan,
			&powod, &propozycja, &pozycja.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pozycji wsadu studio: %w", err)
		}
		pozycja.Powod = tekstZKolumny(powod)
		pozycja.PropozycjaKod = tekstZKolumny(propozycja)
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pozycji wsadu studio %q: %w", kodPrzebiegu, err)
	}
	return lista, nil
}
