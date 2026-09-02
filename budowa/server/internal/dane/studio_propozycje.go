// Odpowiedzialność pliku: propozycje zmiany wypracowane przez operację kontekstową modułu Studio; propozycja nie jest wersją dokumentu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// PropozycjaZmiany to wiersz tabeli `propozycja_zmiany_studio`. WiadomoscID
// wiąże propozycję z wiadomością okna niosącą wynik operacji — pole opcjonalne,
// bo operacja może nie tworzyć wpisu w historii rozmowy.
type PropozycjaZmiany struct {
	ID                      int64
	IdentyfikatorZewnetrzny string
	DokumentID              int64
	AkcjaID                 string
	WiadomoscID             *string
	TrescWyniku             *string
	Tresc                   *string
	TrescOdwolanie          *string
	Utworzono               string
}

const (
	kolumnyPropozycjiZmiany = `id, identyfikator_zewnetrzny, dokument_id, akcja_id, wiadomosc_id,
	                           tresc_wyniku, tresc, tresc_odwolanie, utworzono`

	wstawPropozycjeZmiany = `INSERT INTO propozycja_zmiany_studio
	                         (identyfikator_zewnetrzny, dokument_id, akcja_id, wiadomosc_id,
	                          tresc_wyniku, tresc, tresc_odwolanie)
	                         VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Kod propozycji przychodzi z żądania niezależnie od kodu dokumentu
	// (`studio.diff.compare`, `studio.proposal.decide`), a własnej kolumny konta
	// tabela nie ma: granica idzie drogą po `dokument_id`.
	propozycjaZmianyPoKodzie = `SELECT ` + kolumnyPropozycjiZmiany + `
	                            FROM propozycja_zmiany_studio
	                            WHERE identyfikator_zewnetrzny = ?
	                              AND EXISTS (SELECT 1 FROM dokument_studio
	                                          WHERE dokument_studio.id = propozycja_zmiany_studio.dokument_id
	                                            AND ` + WarunekKonta + `)`

	listaPropozycjiZmianyDokumentu = `SELECT ` + kolumnyPropozycjiZmiany + `
	                                  FROM propozycja_zmiany_studio WHERE dokument_id = ?
	                                  ORDER BY utworzono DESC, id DESC`
)

// ZapiszPropozycje utrwala wynik operacji kontekstowej i oddaje wiersz z czasem nadanym przez bazę; zapis jest zawsze nowym wierszem.
func (r *repozytoriumStudia) ZapiszPropozycje(ctx context.Context, dokumentID int64,
	propozycja PropozycjaZmiany) (PropozycjaZmiany, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wstawPropozycjeZmiany)
	if err != nil {
		return PropozycjaZmiany{}, err
	}
	_, err = polecenie.ExecContext(ctx, propozycja.IdentyfikatorZewnetrzny, dokumentID,
		propozycja.AkcjaID, tekstDoKolumny(propozycja.WiadomoscID),
		tekstDoKolumny(propozycja.TrescWyniku), tekstDoKolumny(propozycja.Tresc),
		tekstDoKolumny(propozycja.TrescOdwolanie))
	if err != nil {
		return PropozycjaZmiany{}, fmt.Errorf(
			"dane: nie można zapisać propozycji zmiany dokumentu %d: %w", dokumentID, err)
	}
	return r.Propozycja(ctx, propozycja.IdentyfikatorZewnetrzny)
}

// Propozycja zwraca propozycję po identyfikatorze zewnętrznym — tak samo jak
// diff.compare adresuje ją w polu `proposalId`.
func (r *repozytoriumStudia) Propozycja(ctx context.Context, kod string) (PropozycjaZmiany, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, propozycjaZmianyPoKodzie)
	if err != nil {
		return PropozycjaZmiany{}, err
	}
	propozycja, err := odczytajPropozycjeZmiany(
		polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return PropozycjaZmiany{}, fmt.Errorf("%w: propozycja zmiany %q", ErrBrakWiersza, kod)
	}
	return propozycja, err
}

// Propozycje zwraca propozycje dokumentu od najnowszej — kolejność zgodna
// z przeglądem repozytorium sesji, żeby Operator widział najświeższy wynik
// operacji kontekstowej jako pierwszy.
func (r *repozytoriumStudia) Propozycje(ctx context.Context, dokumentID int64) ([]PropozycjaZmiany, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPropozycjiZmianyDokumentu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf(
			"dane: nie można odczytać propozycji zmiany dokumentu %d: %w", dokumentID, err)
	}
	defer wiersze.Close()

	lista := []PropozycjaZmiany{}
	for wiersze.Next() {
		propozycja, err := odczytajPropozycjeZmiany(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, propozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf(
			"dane: przerwany odczyt propozycji zmiany dokumentu %d: %w", dokumentID, err)
	}
	return lista, nil
}

// odczytajPropozycjeZmiany składa strukturę propozycji z jednego wiersza wyniku, kolumna po kolumnie SQL.
func odczytajPropozycjeZmiany(wiersz skaner) (PropozycjaZmiany, error) {
	var propozycja PropozycjaZmiany
	var wiadomoscID, trescWyniku, tresc, odwolanie sql.NullString
	err := wiersz.Scan(&propozycja.ID, &propozycja.IdentyfikatorZewnetrzny, &propozycja.DokumentID,
		&propozycja.AkcjaID, &wiadomoscID, &trescWyniku, &tresc, &odwolanie, &propozycja.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return PropozycjaZmiany{}, err
	}
	if err != nil {
		return PropozycjaZmiany{}, fmt.Errorf("dane: nieczytelny wiersz propozycji zmiany: %w", err)
	}
	propozycja.WiadomoscID = tekstZKolumny(wiadomoscID)
	propozycja.TrescWyniku = tekstZKolumny(trescWyniku)
	propozycja.Tresc = tekstZKolumny(tresc)
	propozycja.TrescOdwolanie = tekstZKolumny(odwolanie)
	return propozycja, nil
}
