// Odpowiedzialność pliku: szablony moderacji i artefakty wydane z debaty; artefakt niesie odwołanie do bajtów w magazynie treści, nie bajty.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SzablonModeracjiDebaty to zapisany format prowadzenia debaty, gotowy do ponownego użycia przy kolejnej debacie.
type SzablonModeracjiDebaty struct {
	Kod            string
	Nazwa          string
	Format         string
	GranicaTur     int
	GranicaCzasuMs int
	KolejnoscGlosu []string
	Utworzono      string
}

// ArtefaktDebaty to wydany zapis debaty: transkrypt, graf albo nagranie, z odwołaniem do jego bajtów treści.
type ArtefaktDebaty struct {
	Kod           string
	Okno          string
	Rodzaj        string
	Format        string
	Odwolanie     string
	Rozmiar       int64
	SumaKontrolna string
	DlugoscMs     int
	Utworzono     string
}

// RepozytoriumDebatyWydania jest częścią kontraktu obszaru odpowiadającą za
// szablony moderacji i wydane artefakty.
type RepozytoriumDebatyWydania interface {
	ZapiszSzablonDebaty(ctx context.Context,
		szablon SzablonModeracjiDebaty) (SzablonModeracjiDebaty, error)
	SzablonyDebaty(ctx context.Context, fraza string) ([]SzablonModeracjiDebaty, error)

	ZapiszArtefaktDebaty(ctx context.Context, artefakt ArtefaktDebaty) error
	ArtefaktDebatyPoKodzie(ctx context.Context, kod string) (ArtefaktDebaty, error)
	ArtefaktyDebaty(ctx context.Context, okno string) ([]ArtefaktDebaty, error)
}

const (
	kolumnySzablonuDebaty = `identyfikator_zewnetrzny, nazwa, format, granica_tur,
	                         granica_czasu_ms, kolejnosc_glosu, utworzono`

	zapiszSzablonDebaty = `INSERT INTO debata_szablon
	                       (identyfikator_zewnetrzny, nazwa, format, granica_tur,
	                        granica_czasu_ms, kolejnosc_glosu, konto_id)
	                       VALUES (?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzSzablonDebaty = `SELECT ` + kolumnySzablonuDebaty + `
	                        FROM debata_szablon
	                        WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	pobierzSzablonyDebaty = `SELECT ` + kolumnySzablonuDebaty + `
	                         FROM debata_szablon
	                         WHERE nazwa LIKE '%' || ? || '%' AND ` + WarunekKonta + `
	                         ORDER BY id DESC`

	kolumnyArtefaktuDebaty = `identyfikator_zewnetrzny, okno, rodzaj, format, odwolanie, rozmiar,
	                          suma_kontrolna, dlugosc_ms, utworzono`

	zapiszArtefaktDebaty = `INSERT INTO debata_artefakt
	                        (identyfikator_zewnetrzny, okno, rodzaj, format, odwolanie, rozmiar,
	                         suma_kontrolna, dlugosc_ms, konto_id)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzArtefaktDebaty = `SELECT ` + kolumnyArtefaktuDebaty + `
	                         FROM debata_artefakt
	                         WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	pobierzArtefaktyDebaty = `SELECT ` + kolumnyArtefaktuDebaty + `
	                          FROM debata_artefakt WHERE okno = ? AND ` + WarunekKonta + `
	                          ORDER BY id DESC`
)

// ZapiszSzablonDebaty utrwala szablon moderacji jako osobny wpis gotowy do ponownego użycia w debacie.
func (r *repozytoriumRoundtable) ZapiszSzablonDebaty(ctx context.Context,
	szablon SzablonModeracjiDebaty) (SzablonModeracjiDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSzablonDebaty)
	if err != nil {
		return SzablonModeracjiDebaty{}, err
	}
	if _, err := polecenie.ExecContext(ctx, szablon.Kod, szablon.Nazwa, szablon.Format,
		szablon.GranicaTur, szablon.GranicaCzasuMs,
		strings.Join(szablon.KolejnoscGlosu, "\n"), KontoOperatora(ctx)); err != nil {
		// Więz UNIQUE na identyfikatorze obejmuje całą tabelę: kod zajęty przez konto inne rozbija zapis.
		if czyKolizja(err) {
			return SzablonModeracjiDebaty{}, fmt.Errorf(
				"dane: kod %q nosi szablon moderacji innego konta: %w", szablon.Kod, ErrKolizjaWiersza)
		}
		return SzablonModeracjiDebaty{}, fmt.Errorf("dane: nie można zapisać szablonu moderacji %q: %w",
			szablon.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzSzablonDebaty)
	if err != nil {
		return SzablonModeracjiDebaty{}, err
	}
	return odczytajSzablonDebaty(odczyt.QueryRowContext(ctx, szablon.Kod, KontoOperatora(ctx)))
}

// SzablonyDebaty zwraca szablony moderacji tego okna od najnowszego do najstarszego zapisanego szablonu.
func (r *repozytoriumRoundtable) SzablonyDebaty(ctx context.Context,
	fraza string) ([]SzablonModeracjiDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSzablonyDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, strings.TrimSpace(fraza), KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać szablonów moderacji: %w", err)
	}
	defer wiersze.Close()

	szablony := make([]SzablonModeracjiDebaty, 0, 8)
	for wiersze.Next() {
		szablon, err := odczytajSzablonDebaty(wiersze)
		if err != nil {
			return nil, err
		}
		szablony = append(szablony, szablon)
	}
	return szablony, wiersze.Err()
}

// ZapiszArtefaktDebaty odnotowuje wydany artefakt wraz z odwołaniem do jego bajtów w magazynie treści.
func (r *repozytoriumRoundtable) ZapiszArtefaktDebaty(ctx context.Context, artefakt ArtefaktDebaty) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszArtefaktDebaty)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, artefakt.Kod, artefakt.Okno, artefakt.Rodzaj,
		artefakt.Format, artefakt.Odwolanie, artefakt.Rozmiar, artefakt.SumaKontrolna,
		artefakt.DlugoscMs, KontoOperatora(ctx)); err != nil {
		// Więz UNIQUE na identyfikatorze obejmuje całą tabelę: kod zajęty przez konto inne rozbija zapis.
		if czyKolizja(err) {
			return fmt.Errorf("dane: kod %q nosi artefakt debaty innego konta: %w",
				artefakt.Kod, ErrKolizjaWiersza)
		}
		return fmt.Errorf("dane: nie można zapisać artefaktu debaty %q: %w", artefakt.Kod, err)
	}
	return nil
}

// ArtefaktDebatyPoKodzie zwraca artefakt po jego identyfikatorze zewnętrznym, wraz z całą jego treścią.
func (r *repozytoriumRoundtable) ArtefaktDebatyPoKodzie(ctx context.Context,
	kod string) (ArtefaktDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzArtefaktDebaty)
	if err != nil {
		return ArtefaktDebaty{}, err
	}
	return odczytajArtefaktDebaty(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
}

// ArtefaktyDebaty zwraca artefakty danego okna od najnowszego do najstarszego wydanego artefaktu debaty.
func (r *repozytoriumRoundtable) ArtefaktyDebaty(ctx context.Context,
	okno string) ([]ArtefaktDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzArtefaktyDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać artefaktów debaty okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	artefakty := make([]ArtefaktDebaty, 0, 8)
	for wiersze.Next() {
		artefakt, err := odczytajArtefaktDebaty(wiersze)
		if err != nil {
			return nil, err
		}
		artefakty = append(artefakty, artefakt)
	}
	return artefakty, wiersze.Err()
}

// odczytajSzablonDebaty składa szablon moderacji z jednego wiersza wyniku zapytania, kolumna po kolumnie.
func odczytajSzablonDebaty(wiersz interface{ Scan(...any) error }) (SzablonModeracjiDebaty, error) {
	var szablon SzablonModeracjiDebaty
	var kolejnosc string
	err := wiersz.Scan(&szablon.Kod, &szablon.Nazwa, &szablon.Format, &szablon.GranicaTur,
		&szablon.GranicaCzasuMs, &kolejnosc, &szablon.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return SzablonModeracjiDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return SzablonModeracjiDebaty{}, fmt.Errorf("dane: nieczytelny wiersz szablonu moderacji: %w", err)
	}
	szablon.KolejnoscGlosu = rozdzielWierszeDebaty(kolejnosc)
	return szablon, nil
}

// odczytajArtefaktDebaty składa artefakt debaty z jednego wiersza wyniku zapytania, kolumna po kolumnie.
func odczytajArtefaktDebaty(wiersz interface{ Scan(...any) error }) (ArtefaktDebaty, error) {
	var artefakt ArtefaktDebaty
	err := wiersz.Scan(&artefakt.Kod, &artefakt.Okno, &artefakt.Rodzaj, &artefakt.Format,
		&artefakt.Odwolanie, &artefakt.Rozmiar, &artefakt.SumaKontrolna, &artefakt.DlugoscMs,
		&artefakt.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return ArtefaktDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return ArtefaktDebaty{}, fmt.Errorf("dane: nieczytelny wiersz artefaktu debaty: %w", err)
	}
	return artefakt, nil
}
