// Odpowiedzialność pliku: polityki przechowywania (`polityka_retencji_biblioteki`)
// i zapis utrwalenia archiwalnego (`zadanie_utrwalenia_biblioteki`) — migracja 184.
//
// Polityka opisuje regułę, nie zdarzenie: mówi, ile dni zasób ma być
// przechowywany i co ma się z nim stać po upływie okresu. Sam upływ liczy się
// przy odczycie raportu retencji, bo data graniczna wynika z chwili pytania —
// wartość zapisana starzałaby się w bazie.
//
// Zapis utrwalenia jest śladem po pracy, nie zleceniem do wykonania: powstaje po
// utrwaleniu i niesie wynik walidacji wraz z jej zapisem.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// PolitykaRetencjiBiblioteki to wiersz tabeli `polityka_retencji_biblioteki`.
type PolitykaRetencjiBiblioteki struct {
	ID                int64
	Kod               string
	Zasieg            string
	ZasiegID          *string
	DniPrzechowywania int
	Czynnosc          string
	Utworzono         string
}

// ZadanieUtrwaleniaBiblioteki to wiersz tabeli `zadanie_utrwalenia_biblioteki`.
type ZadanieUtrwaleniaBiblioteki struct {
	ID              int64
	Kod             string
	PlikKod         string
	Rodzaj          string
	PlikWynikowyKod *string
	Profil          *string
	Poprawne        bool
	Raport          string
	Utworzono       string
}

const (
	kolumnyPolitykiRetencjiBiblioteki = `id, identyfikator_zewnetrzny, zasieg, zasieg_id,
	                                     dni_przechowywania, czynnosc, utworzono`

	zapiszPolitykeRetencjiBiblioteki = `INSERT INTO polityka_retencji_biblioteki
	                                    (identyfikator_zewnetrzny, zasieg, zasieg_id,
	                                     dni_przechowywania, czynnosc)
	                                    VALUES (?, ?, ?, ?, ?)
	                                    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                        zasieg = excluded.zasieg,
	                                        zasieg_id = excluded.zasieg_id,
	                                        dni_przechowywania = excluded.dni_przechowywania,
	                                        czynnosc = excluded.czynnosc`

	pobierzPolitykeRetencjiBiblioteki = `SELECT ` + kolumnyPolitykiRetencjiBiblioteki + `
	                                     FROM polityka_retencji_biblioteki
	                                     WHERE identyfikator_zewnetrzny = ?`

	usunPolitykeRetencjiBiblioteki = `DELETE FROM polityka_retencji_biblioteki
	                                  WHERE identyfikator_zewnetrzny = ?`

	kolumnyUtrwaleniaBiblioteki = `id, identyfikator_zewnetrzny, plik_kod, rodzaj,
	                               plik_wynikowy_kod, profil, poprawne, raport, utworzono`

	zapiszUtrwalenieBiblioteki = `INSERT INTO zadanie_utrwalenia_biblioteki
	                              (identyfikator_zewnetrzny, plik_kod, rodzaj, plik_wynikowy_kod,
	                               profil, poprawne, raport)
	                              VALUES (?, ?, ?, ?, ?, ?, ?)`

	pobierzUtrwalenieBiblioteki = `SELECT ` + kolumnyUtrwaleniaBiblioteki + `
	                               FROM zadanie_utrwalenia_biblioteki
	                               WHERE identyfikator_zewnetrzny = ?`
)

// ZapiszPolitykeRetencji zakłada politykę albo zmienia zastaną.
func (r *repozytoriumBiblioteki) ZapiszPolitykeRetencji(ctx context.Context,
	polityka PolitykaRetencjiBiblioteki) (PolitykaRetencjiBiblioteki, error) {

	if strings.TrimSpace(polityka.Kod) == "" {
		return PolitykaRetencjiBiblioteki{}, fmt.Errorf("dane: polityka retencji bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPolitykeRetencjiBiblioteki)
	if err != nil {
		return PolitykaRetencjiBiblioteki{}, err
	}
	_, err = polecenie.ExecContext(ctx, polityka.Kod, polityka.Zasieg,
		tekstDoKolumny(polityka.ZasiegID), polityka.DniPrzechowywania, polityka.Czynnosc)
	if err != nil {
		return PolitykaRetencjiBiblioteki{}, fmt.Errorf("dane: nie można zapisać polityki retencji %q: %w",
			polityka.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzPolitykeRetencjiBiblioteki)
	if err != nil {
		return PolitykaRetencjiBiblioteki{}, err
	}
	zapisana, err := odczytajPolitykeRetencjiBiblioteki(odczyt.QueryRowContext(ctx, polityka.Kod))
	if err != nil {
		return PolitykaRetencjiBiblioteki{}, fmt.Errorf("dane: nieczytelna polityka retencji %q: %w",
			polityka.Kod, err)
	}
	return zapisana, nil
}

// PolitykiRetencji zwraca polityki, zawężone poziomem zasięgu.
func (r *repozytoriumBiblioteki) PolitykiRetencji(ctx context.Context,
	zasieg *string) ([]PolitykaRetencjiBiblioteki, error) {

	warunek := "1 = 1"
	argumenty := []any{}
	if zasieg != nil && *zasieg != "" {
		warunek = "zasieg = ?"
		argumenty = append(argumenty, *zasieg)
	}
	zapytanie := `SELECT ` + kolumnyPolitykiRetencjiBiblioteki + `
	              FROM polityka_retencji_biblioteki WHERE ` + warunek + `
	              ORDER BY dni_przechowywania, id`

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać polityk retencji: %w", err)
	}
	defer wiersze.Close()

	lista := []PolitykaRetencjiBiblioteki{}
	for wiersze.Next() {
		polityka, err := odczytajPolitykeRetencjiBiblioteki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz polityki retencji: %w", err)
		}
		lista = append(lista, polityka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt polityk retencji: %w", err)
	}
	return lista, nil
}

// UsunPolitykeRetencji zdejmuje politykę.
func (r *repozytoriumBiblioteki) UsunPolitykeRetencji(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunPolitykeRetencjiBiblioteki)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć polityki retencji %q: %w", kod, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć usuniętych polityk: %w", err)
	}
	return zdjete > 0, nil
}

// ZapiszUtrwalenie odkłada ślad po utrwaleniu archiwalnym wraz z wynikiem
// walidacji.
func (r *repozytoriumBiblioteki) ZapiszUtrwalenie(ctx context.Context,
	zadanie ZadanieUtrwaleniaBiblioteki) (ZadanieUtrwaleniaBiblioteki, error) {

	if strings.TrimSpace(zadanie.Kod) == "" {
		return ZadanieUtrwaleniaBiblioteki{}, fmt.Errorf("dane: zapis utrwalenia bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszUtrwalenieBiblioteki)
	if err != nil {
		return ZadanieUtrwaleniaBiblioteki{}, err
	}
	_, err = polecenie.ExecContext(ctx, zadanie.Kod, zadanie.PlikKod, zadanie.Rodzaj,
		tekstDoKolumny(zadanie.PlikWynikowyKod), tekstDoKolumny(zadanie.Profil),
		liczbaLogiczna(zadanie.Poprawne), zadanie.Raport)
	if err != nil {
		return ZadanieUtrwaleniaBiblioteki{}, fmt.Errorf("dane: nie można zapisać utrwalenia %q: %w",
			zadanie.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzUtrwalenieBiblioteki)
	if err != nil {
		return ZadanieUtrwaleniaBiblioteki{}, err
	}
	zapisane, err := odczytajUtrwalenieBiblioteki(odczyt.QueryRowContext(ctx, zadanie.Kod))
	if err != nil {
		return ZadanieUtrwaleniaBiblioteki{}, fmt.Errorf("dane: nieczytelny zapis utrwalenia %q: %w",
			zadanie.Kod, err)
	}
	return zapisane, nil
}

// odczytajPolitykeRetencjiBiblioteki składa politykę z jednego wiersza wyniku.
func odczytajPolitykeRetencjiBiblioteki(wiersz skaner) (PolitykaRetencjiBiblioteki, error) {
	var polityka PolitykaRetencjiBiblioteki
	var zasiegID sql.NullString
	err := wiersz.Scan(&polityka.ID, &polityka.Kod, &polityka.Zasieg, &zasiegID,
		&polityka.DniPrzechowywania, &polityka.Czynnosc, &polityka.Utworzono)
	if err != nil {
		return PolitykaRetencjiBiblioteki{}, err
	}
	polityka.ZasiegID = tekstZKolumny(zasiegID)
	return polityka, nil
}

// odczytajUtrwalenieBiblioteki składa zapis utrwalenia z jednego wiersza wyniku.
func odczytajUtrwalenieBiblioteki(wiersz skaner) (ZadanieUtrwaleniaBiblioteki, error) {
	var zadanie ZadanieUtrwaleniaBiblioteki
	var wynikowy, profil sql.NullString
	var poprawne int
	err := wiersz.Scan(&zadanie.ID, &zadanie.Kod, &zadanie.PlikKod, &zadanie.Rodzaj,
		&wynikowy, &profil, &poprawne, &zadanie.Raport, &zadanie.Utworzono)
	if err != nil {
		return ZadanieUtrwaleniaBiblioteki{}, err
	}
	zadanie.PlikWynikowyKod, zadanie.Profil = tekstZKolumny(wynikowy), tekstZKolumny(profil)
	zadanie.Poprawne = poprawne == 1
	return zadanie, nil
}
