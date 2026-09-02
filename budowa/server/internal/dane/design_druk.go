// Plik prowadzi obszar profili wydania do druku oraz licencji zasobów wciągniętych z katalogów zewnętrznych, część
// RepozytoriumDesignu; profil jest bytem trwałym, bo czytają go naraz kontrola przeddrukowa i wydanie do druku.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ProfilDrukuDesignu to wiersz tabeli `profil_druku_design`. Brak spadu znaczy
// „bierz domyślny rdzenia", spad zerowy „drukuj bez spadu".
type ProfilDrukuDesignu struct {
	ID                 int64
	Kod                string
	Okno               string
	Nazwa              *string
	PrzestrzenBarw     string
	Norma              *string
	SpadMm             *float64
	ZnacznikiCiecia    bool
	ZnacznikiPasowania bool
	PasekBarw          bool
	Rozdzielczosc      *int64
	ProfilICC          *string
	NadrukCzerni       bool
	Nosnik             *string
	Zaktualizowano     string
}

// LicencjaZasobuDesignu to wiersz tabeli `licencja_zasobu_design`.
type LicencjaZasobuDesignu struct {
	ID                     int64
	ZasobID                int64
	Dostawca               string
	IdentyfikatorUDostawcy string
	Licencja               *string
	Autor                  *string
	Odsylacz               *string
	Utworzono              string
}

const (
	kolumnyProfiluDrukuDesignu = `id, identyfikator_zewnetrzny, okno, nazwa, przestrzen_barw,
	                              norma, spad_mm, znaczniki_ciecia, znaczniki_pasowania,
	                              pasek_barw, rozdzielczosc, profil_icc, nadruk_czerni, nosnik,
	                              zaktualizowano`

	zapiszProfilDrukuDesignu = `INSERT INTO profil_druku_design
	                            (identyfikator_zewnetrzny, okno, nazwa, przestrzen_barw, norma,
	                             spad_mm, znaczniki_ciecia, znaczniki_pasowania, pasek_barw,
	                             rozdzielczosc, profil_icc, nadruk_czerni, nosnik, zaktualizowano,
	                             konto_id)
	                            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	                                    strftime('%Y-%m-%dT%H:%M:%fZ','now'), ` + WskazanieKonta + `)
	                            ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                nazwa = excluded.nazwa,
	                                przestrzen_barw = excluded.przestrzen_barw,
	                                norma = excluded.norma,
	                                spad_mm = excluded.spad_mm,
	                                znaczniki_ciecia = excluded.znaczniki_ciecia,
	                                znaczniki_pasowania = excluded.znaczniki_pasowania,
	                                pasek_barw = excluded.pasek_barw,
	                                rozdzielczosc = excluded.rozdzielczosc,
	                                profil_icc = excluded.profil_icc,
	                                nadruk_czerni = excluded.nadruk_czerni,
	                                nosnik = excluded.nosnik,
	                                zaktualizowano = excluded.zaktualizowano
	                            WHERE ` + WarunekKonta

	pobierzProfilDrukuDesignu = `SELECT ` + kolumnyProfiluDrukuDesignu +
		` FROM profil_druku_design WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaProfiliDrukuDesignu = `SELECT ` + kolumnyProfiluDrukuDesignu +
		` FROM profil_druku_design WHERE okno = ? AND ` + WarunekKonta +
		` ORDER BY zaktualizowano DESC, id DESC`

	// Licencja własnej kolumny konta nie ma: granica dochodzi przez `zasob_id`
	// do korzenia zasob_design (migracja 484).
	zasobDesignWGranicyKonta = `EXISTS (SELECT 1 FROM zasob_design
	                                     WHERE zasob_design.id = licencja_zasobu_design.zasob_id
	                                       AND ` + WarunekKonta + `)`

	zapiszLicencjeZasobuDesignu = `INSERT INTO licencja_zasobu_design
	                               (zasob_id, dostawca, identyfikator_u_dostawcy, licencja,
	                                autor, odsylacz)
	                               SELECT ?, ?, ?, ?, ?, ?
	                                WHERE EXISTS (SELECT 1 FROM zasob_design
	                                               WHERE zasob_design.id = ? AND ` + WarunekKonta + `)
	                               ON CONFLICT(zasob_id) DO UPDATE SET
	                                   dostawca = excluded.dostawca,
	                                   identyfikator_u_dostawcy = excluded.identyfikator_u_dostawcy,
	                                   licencja = excluded.licencja,
	                                   autor = excluded.autor,
	                                   odsylacz = excluded.odsylacz
	                               WHERE ` + zasobDesignWGranicyKonta

	pobierzLicencjeZasobuDesignu = `SELECT id, zasob_id, dostawca, identyfikator_u_dostawcy,
	                                       licencja, autor, odsylacz, utworzono
	                                FROM licencja_zasobu_design
	                                WHERE zasob_id = ? AND ` + zasobDesignWGranicyKonta
)

func (r *repozytoriumDesignu) ZapiszProfilDrukuDesignu(ctx context.Context,
	profil ProfilDrukuDesignu) (ProfilDrukuDesignu, error) {

	if profil.Kod == "" {
		return ProfilDrukuDesignu{}, fmt.Errorf("dane: profil druku design bez identyfikatora")
	}
	if profil.Okno == "" {
		return ProfilDrukuDesignu{}, fmt.Errorf("dane: profil druku design %q bez okna", profil.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszProfilDrukuDesignu)
	if err != nil {
		return ProfilDrukuDesignu{}, err
	}
	var spad any
	if profil.SpadMm != nil {
		spad = *profil.SpadMm
	}
	wynik, err := polecenie.ExecContext(ctx, profil.Kod, profil.Okno, tekstDoKolumny(profil.Nazwa),
		profil.PrzestrzenBarw, tekstDoKolumny(profil.Norma), spad,
		liczbaLogiczna(profil.ZnacznikiCiecia), liczbaLogiczna(profil.ZnacznikiPasowania),
		liczbaLogiczna(profil.PasekBarw), liczbaDoKolumny(profil.Rozdzielczosc),
		tekstDoKolumny(profil.ProfilICC), liczbaLogiczna(profil.NadrukCzerni),
		tekstDoKolumny(profil.Nosnik), KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return ProfilDrukuDesignu{}, fmt.Errorf("dane: nie można zapisać profilu druku design %q: %w",
			profil.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "profil druku design", profil.Kod); err != nil {
		return ProfilDrukuDesignu{}, err
	}
	return r.ProfilDrukuDesignuPoKodzie(ctx, profil.Kod)
}

func (r *repozytoriumDesignu) ProfilDrukuDesignuPoKodzie(ctx context.Context,
	kod string) (ProfilDrukuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzProfilDrukuDesignu)
	if err != nil {
		return ProfilDrukuDesignu{}, err
	}
	profil, err := odczytajProfilDrukuDesignu(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return ProfilDrukuDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return ProfilDrukuDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz profilu druku design %q: %w", kod, err)
	}
	return profil, nil
}

func (r *repozytoriumDesignu) ProfileDrukuDesignu(ctx context.Context,
	okno string) ([]ProfilDrukuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaProfiliDrukuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać profili druku design okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []ProfilDrukuDesignu{}
	for wiersze.Next() {
		profil, err := odczytajProfilDrukuDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz profilu druku design okna %q: %w",
				okno, err)
		}
		lista = append(lista, profil)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt profili druku design okna %q: %w", okno, err)
	}
	return lista, nil
}

func (r *repozytoriumDesignu) ZapiszLicencjeZasobuDesignu(ctx context.Context,
	licencja LicencjaZasobuDesignu) error {

	if licencja.ZasobID == 0 {
		return fmt.Errorf("dane: licencja zasobu design bez wskazania zasobu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszLicencjeZasobuDesignu)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, licencja.ZasobID, licencja.Dostawca,
		licencja.IdentyfikatorUDostawcy, tekstDoKolumny(licencja.Licencja),
		tekstDoKolumny(licencja.Autor), tekstDoKolumny(licencja.Odsylacz),
		licencja.ZasobID, KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać licencji zasobu design %d: %w",
			licencja.ZasobID, err)
	}
	return sprawdzTrafienieZapisu(wynik, "licencja zasobu design", fmt.Sprint(licencja.ZasobID))
}

// LicencjaZasobuDesignuPoZasobie oddaje ErrBrakWiersza także dla zasobu spoza
// katalogu zewnętrznego — brak licencji jest stanem poprawnym.
func (r *repozytoriumDesignu) LicencjaZasobuDesignuPoZasobie(ctx context.Context,
	zasobID int64) (LicencjaZasobuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzLicencjeZasobuDesignu)
	if err != nil {
		return LicencjaZasobuDesignu{}, err
	}
	var licencja LicencjaZasobuDesignu
	var tresc, autor, odsylacz sql.NullString
	err = polecenie.QueryRowContext(ctx, zasobID, KontoOperatora(ctx)).Scan(&licencja.ID, &licencja.ZasobID,
		&licencja.Dostawca, &licencja.IdentyfikatorUDostawcy, &tresc, &autor, &odsylacz,
		&licencja.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return LicencjaZasobuDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return LicencjaZasobuDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz licencji zasobu design %d: %w", zasobID, err)
	}
	licencja.Licencja = tekstZKolumny(tresc)
	licencja.Autor = tekstZKolumny(autor)
	licencja.Odsylacz = tekstZKolumny(odsylacz)
	return licencja, nil
}

func odczytajProfilDrukuDesignu(wiersz skaner) (ProfilDrukuDesignu, error) {
	var profil ProfilDrukuDesignu
	var nazwa, norma, profilICC, nosnik sql.NullString
	var spad sql.NullFloat64
	var rozdzielczosc sql.NullInt64
	var ciecia, pasowania, pasek, nadruk int
	err := wiersz.Scan(&profil.ID, &profil.Kod, &profil.Okno, &nazwa, &profil.PrzestrzenBarw,
		&norma, &spad, &ciecia, &pasowania, &pasek, &rozdzielczosc, &profilICC, &nadruk,
		&nosnik, &profil.Zaktualizowano)
	if err != nil {
		return ProfilDrukuDesignu{}, err
	}
	profil.Nazwa = tekstZKolumny(nazwa)
	profil.Norma = tekstZKolumny(norma)
	profil.SpadMm = liczbaRzeczywistaZKolumny(spad)
	profil.Rozdzielczosc = liczbaZKolumny(rozdzielczosc)
	profil.ProfilICC = tekstZKolumny(profilICC)
	profil.Nosnik = tekstZKolumny(nosnik)
	profil.ZnacznikiCiecia = ciecia == 1
	profil.ZnacznikiPasowania = pasowania == 1
	profil.PasekBarw = pasek == 1
	profil.NadrukCzerni = nadruk == 1
	return profil, nil
}
