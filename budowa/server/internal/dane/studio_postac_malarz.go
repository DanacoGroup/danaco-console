// Repozytorium przechowuje w tabeli `postac_malarza_studio` z migracji 368 postać
// formatów skopiowaną malarzem formatów modułu Studio do naniesienia w innym miejscu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// PostacMalarzaStudia to wiersz tabeli `postac_malarza_studio`, niosący postać
// zabraną malarzem formatów wraz z chwilą wygaśnięcia.
type PostacMalarzaStudia struct {
	ID          int64
	Kod         string
	Okno        string
	DokumentKod *string
	// PostacZnakuJSON i PostacAkapituJSON niosą postać kontraktu w zapisie JSON, nie
	// w osobnych kolumnach.
	PostacZnakuJSON   *string
	PostacAkapituJSON *string
	StylNazwany       *string
	Utworzono         string
	// Wygasa jest chwilą, po której wpisu nie wolno nanieść; wartość pusta oznacza
	// brak wygaśnięcia.
	Wygasa *string
}

const (
	malarzKolumny = `id, identyfikator_zewnetrzny, okno, dokument_kod, postac_znaku_json,
	                 postac_akapitu_json, styl_nazwany, utworzono, wygasa`

	malarzZapisz = `INSERT INTO postac_malarza_studio
	                (identyfikator_zewnetrzny, okno, dokument_kod, postac_znaku_json,
	                 postac_akapitu_json, styl_nazwany, wygasa)
	                VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Odczyt po uchwycie NIE pyta o okno: uchwyt jest identyfikatorem nadanym
	// przez rdzeń i dowodzi sam z siebie. Warunek okna kazałby naniesieniu
	// podawać okno drugi raz, a model, który pobrał postać jednym narzędziem,
	// niesie do drugiego uchwyt, nie okno.
	malarzPobierz = `SELECT ` + malarzKolumny + ` FROM postac_malarza_studio
	                 WHERE identyfikator_zewnetrzny = ?`

	// Najświeższa postać okna — droga naniesienia bez podanego uchwytu
	// („nanieś to, co ostatnio zabrałem"). Wygasłych nie oddaje.
	malarzNajswiezsza = `SELECT ` + malarzKolumny + ` FROM postac_malarza_studio
	                     WHERE okno = ?
	                       AND (wygasa IS NULL OR wygasa > strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                     ORDER BY utworzono DESC, id DESC LIMIT 1`

	// Sprzątanie zdejmuje wpisy wygasłe. Idzie osobnym wywołaniem, nie w odczycie:
	// odczyt kasujący wiersze byłby odczytem zmieniającym stan, a przy dwóch
	// oknach pytających naraz jedno zabierałoby postać drugiemu.
	malarzSprzataj = `DELETE FROM postac_malarza_studio
	                  WHERE wygasa IS NOT NULL
	                    AND wygasa <= strftime('%Y-%m-%dT%H:%M:%fZ','now')`
)

// MalarzFormatowStudia jest kontraktem warstwy danych, deklarującym zapis, odczyt
// i sprzątanie postaci zabranej malarzem formatów.
type MalarzFormatowStudia interface {
	ZapiszPostacMalarza(ctx context.Context,
		postac PostacMalarzaStudia) (PostacMalarzaStudia, error)
	PostacMalarza(ctx context.Context, kod string) (PostacMalarzaStudia, error)
	NajswiezszaPostacMalarza(ctx context.Context, okno string) (PostacMalarzaStudia, error)
	SprzatnijPostacieMalarza(ctx context.Context) (int, error)
}

// ZapiszPostacMalarza odkłada w tabeli `postac_malarza_studio` postać zabraną
// malarzem formatów i oddaje ją zapisaną.
func (r *repozytoriumStudia) ZapiszPostacMalarza(ctx context.Context,
	postac PostacMalarzaStudia) (PostacMalarzaStudia, error) {

	kod := strings.TrimSpace(postac.Kod)
	if kod == "" {
		return PostacMalarzaStudia{}, fmt.Errorf("dane: postać malarza bez identyfikatora")
	}
	// Postać bez żadnej z trzech treści nie jest postacią i nie zmienia stanu.
	if postac.PostacZnakuJSON == nil && postac.PostacAkapituJSON == nil &&
		postac.StylNazwany == nil {

		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: postać malarza %q nie niesie ani postaci znaku, ani postaci akapitu, "+
				"ani stylu nazwanego", kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, malarzZapisz)
	if err != nil {
		return PostacMalarzaStudia{}, err
	}
	if _, err := polecenie.ExecContext(ctx, kod, strings.TrimSpace(postac.Okno),
		tekstDoKolumny(postac.DokumentKod), tekstDoKolumny(postac.PostacZnakuJSON),
		tekstDoKolumny(postac.PostacAkapituJSON), tekstDoKolumny(postac.StylNazwany),
		tekstDoKolumny(postac.Wygasa)); err != nil {

		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: nie można zapisać postaci malarza %q: %w", kod, err)
	}
	return r.PostacMalarza(ctx, kod)
}

// PostacMalarza oddaje postać o wskazanym uchwycie.
//
// Wpis wygasły jest tu BRAKIEM WIERSZA, nie wierszem z datą w przeszłości:
// rdzeń ma powiedzieć „ta postać już nie obowiązuje", a nie nanieść postać,
// której Operator nie pamięta.
func (r *repozytoriumStudia) PostacMalarza(ctx context.Context,
	kod string) (PostacMalarzaStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, malarzPobierz)
	if err != nil {
		return PostacMalarzaStudia{}, err
	}
	postac, err := malarzOdczytaj(polecenie.QueryRowContext(ctx, strings.TrimSpace(kod)))
	if errors.Is(err, sql.ErrNoRows) {
		return PostacMalarzaStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz postaci malarza %q: %w", kod, err)
	}
	return postac, nil
}

// NajswiezszaPostacMalarza oddaje ostatnią niewygasłą postać zabraną malarzem
// formatów we wskazanym oknie.
func (r *repozytoriumStudia) NajswiezszaPostacMalarza(ctx context.Context,
	okno string) (PostacMalarzaStudia, error) {

	if strings.TrimSpace(okno) == "" {
		return PostacMalarzaStudia{}, ErrBrakWiersza
	}
	polecenie, err := r.zapytania.przygotuj(ctx, malarzNajswiezsza)
	if err != nil {
		return PostacMalarzaStudia{}, err
	}
	postac, err := malarzOdczytaj(polecenie.QueryRowContext(ctx, strings.TrimSpace(okno)))
	if errors.Is(err, sql.ErrNoRows) {
		return PostacMalarzaStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz najświeższej postaci malarza okna %q: %w", okno, err)
	}
	return postac, nil
}

// SprzatnijPostacieMalarza zdejmuje z tabeli `postac_malarza_studio` wpisy
// wygasłe i oddaje liczbę zdjętych wierszy.
func (r *repozytoriumStudia) SprzatnijPostacieMalarza(ctx context.Context) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, malarzSprzataj)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zdjąć wygasłych postaci malarza: %w", err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany skutek zdjęcia wygasłych postaci malarza: %w", err)
	}
	return int(zdjete), nil
}

// malarzOdczytaj składa strukturę PostacMalarzaStudia z jednego wiersza wyniku
// zapytania SQL do tabeli `postac_malarza_studio`.
func malarzOdczytaj(wiersz interface{ Scan(...any) error }) (PostacMalarzaStudia, error) {
	var postac PostacMalarzaStudia
	var dokument, znak, akapit, styl, wygasa sql.NullString
	if err := wiersz.Scan(&postac.ID, &postac.Kod, &postac.Okno, &dokument, &znak,
		&akapit, &styl, &postac.Utworzono, &wygasa); err != nil {

		return PostacMalarzaStudia{}, err
	}
	postac.DokumentKod = tekstZKolumny(dokument)
	postac.PostacZnakuJSON = tekstZKolumny(znak)
	postac.PostacAkapituJSON = tekstZKolumny(akapit)
	postac.StylNazwany = tekstZKolumny(styl)
	postac.Wygasa = tekstZKolumny(wygasa)
	return postac, nil
}
