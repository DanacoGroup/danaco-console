// Postać formatów skopiowana malarzem formatów modułu Studio (tabela `postac_malarza_studio`) do naniesienia gdzie indziej.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// PostacMalarzaStudia to wiersz `postac_malarza_studio` — postać zabrana malarzem formatów, wraz z chwilą wygaśnięcia.
type PostacMalarzaStudia struct {
	ID                int64
	Kod               string
	Okno              string
	DokumentKod       *string
	PostacZnakuJSON   *string
	PostacAkapituJSON *string
	StylNazwany       *string
	Utworzono         string
	Wygasa            *string
}

const (
	malarzKolumny = `id, identyfikator_zewnetrzny, okno, dokument_kod, postac_znaku_json,
	                 postac_akapitu_json, styl_nazwany, utworzono, wygasa`

	malarzZapisz = `INSERT INTO postac_malarza_studio
	                (identyfikator_zewnetrzny, okno, dokument_kod, postac_znaku_json,
	                 postac_akapitu_json, styl_nazwany, wygasa, konto_id)
	                VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	malarzPobierz = `SELECT ` + malarzKolumny + ` FROM postac_malarza_studio
	                 WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	malarzNajswiezsza = `SELECT ` + malarzKolumny + ` FROM postac_malarza_studio
	                     WHERE okno = ? AND ` + WarunekKonta + `
	                       AND (wygasa IS NULL OR wygasa > strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                     ORDER BY utworzono DESC, id DESC LIMIT 1`

	malarzSprzataj = `DELETE FROM postac_malarza_studio
	                  WHERE wygasa IS NOT NULL
	                    AND wygasa <= strftime('%Y-%m-%dT%H:%M:%fZ','now') AND ` + WarunekKonta
)

type MalarzFormatowStudia interface {
	ZapiszPostacMalarza(ctx context.Context,
		postac PostacMalarzaStudia) (PostacMalarzaStudia, error)
	PostacMalarza(ctx context.Context, kod string) (PostacMalarzaStudia, error)
	NajswiezszaPostacMalarza(ctx context.Context, okno string) (PostacMalarzaStudia, error)
	SprzatnijPostacieMalarza(ctx context.Context) (int, error)
}

// ZapiszPostacMalarza odkłada postać zabraną malarzem formatów i oddaje ją zapisaną.
func (r *repozytoriumStudia) ZapiszPostacMalarza(ctx context.Context,
	postac PostacMalarzaStudia) (PostacMalarzaStudia, error) {

	kod := strings.TrimSpace(postac.Kod)
	if kod == "" {
		return PostacMalarzaStudia{}, fmt.Errorf("dane: postać malarza bez identyfikatora")
	}
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
		tekstDoKolumny(postac.Wygasa), KontoOperatora(ctx)); err != nil {

		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: nie można zapisać postaci malarza %q: %w", kod, err)
	}
	return r.PostacMalarza(ctx, kod)
}

// PostacMalarza oddaje postać o wskazanym uchwycie; wpis wygasły jest tu brakiem wiersza, nie datą w przeszłości.
func (r *repozytoriumStudia) PostacMalarza(ctx context.Context,
	kod string) (PostacMalarzaStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, malarzPobierz)
	if err != nil {
		return PostacMalarzaStudia{}, err
	}
	postac, err := malarzOdczytaj(polecenie.QueryRowContext(ctx, strings.TrimSpace(kod), KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return PostacMalarzaStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz postaci malarza %q: %w", kod, err)
	}
	return postac, nil
}

// NajswiezszaPostacMalarza oddaje ostatnią niewygasłą postać malarza we wskazanym oknie.
func (r *repozytoriumStudia) NajswiezszaPostacMalarza(ctx context.Context,
	okno string) (PostacMalarzaStudia, error) {

	if strings.TrimSpace(okno) == "" {
		return PostacMalarzaStudia{}, ErrBrakWiersza
	}
	polecenie, err := r.zapytania.przygotuj(ctx, malarzNajswiezsza)
	if err != nil {
		return PostacMalarzaStudia{}, err
	}
	postac, err := malarzOdczytaj(polecenie.QueryRowContext(ctx, strings.TrimSpace(okno), KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return PostacMalarzaStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return PostacMalarzaStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz najświeższej postaci malarza okna %q: %w", okno, err)
	}
	return postac, nil
}

// SprzatnijPostacieMalarza zdejmuje wpisy wygasłe i oddaje liczbę zdjętych wierszy.
func (r *repozytoriumStudia) SprzatnijPostacieMalarza(ctx context.Context) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, malarzSprzataj)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, KontoOperatora(ctx))
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zdjąć wygasłych postaci malarza: %w", err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany skutek zdjęcia wygasłych postaci malarza: %w", err)
	}
	return int(zdjete), nil
}

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
