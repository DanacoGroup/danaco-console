// Polityki przechowywania i ślad utrwalenia archiwalnego; upływ terminu liczy się przy odczycie raportu retencji.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// PolitykaRetencjiBiblioteki to wiersz tabeli `polityka_retencji_biblioteki` niosący regułę przechowywania zasobu.
type PolitykaRetencjiBiblioteki struct {
	ID                int64
	Kod               string
	Zasieg            string
	ZasiegID          *string
	DniPrzechowywania int
	Czynnosc          string
	Utworzono         string
}

// ZadanieUtrwaleniaBiblioteki to wiersz tabeli `zadanie_utrwalenia_biblioteki` niosący ślad wykonanego utrwalenia.
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
	                                     dni_przechowywania, czynnosc, konto_id)
	                                    VALUES (?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                                    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                        zasieg = excluded.zasieg,
	                                        zasieg_id = excluded.zasieg_id,
	                                        dni_przechowywania = excluded.dni_przechowywania,
	                                        czynnosc = excluded.czynnosc
	                                    WHERE ` + WarunekKonta

	pobierzPolitykeRetencjiBiblioteki = `SELECT ` + kolumnyPolitykiRetencjiBiblioteki + `
	                                     FROM polityka_retencji_biblioteki
	                                     WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	usunPolitykeRetencjiBiblioteki = `DELETE FROM polityka_retencji_biblioteki
	                                  WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	kolumnyUtrwaleniaBiblioteki = `id, identyfikator_zewnetrzny, plik_kod, rodzaj,
	                               plik_wynikowy_kod, profil, poprawne, raport, utworzono`

	zapiszUtrwalenieBiblioteki = `INSERT INTO zadanie_utrwalenia_biblioteki
	                              (identyfikator_zewnetrzny, plik_kod, rodzaj, plik_wynikowy_kod,
	                               profil, poprawne, raport, konto_id)
	                              VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzUtrwalenieBiblioteki = `SELECT ` + kolumnyUtrwaleniaBiblioteki + `
	                               FROM zadanie_utrwalenia_biblioteki
	                               WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

// ZapiszPolitykeRetencji zakłada politykę retencji albo zmienia zastaną politykę tego samego zasięgu w bazie.
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
		tekstDoKolumny(polityka.ZasiegID), polityka.DniPrzechowywania, polityka.Czynnosc,
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return PolitykaRetencjiBiblioteki{}, fmt.Errorf("dane: nie można zapisać polityki retencji %q: %w",
			polityka.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzPolitykeRetencjiBiblioteki)
	if err != nil {
		return PolitykaRetencjiBiblioteki{}, err
	}
	zapisana, err := odczytajPolitykeRetencjiBiblioteki(odczyt.QueryRowContext(ctx, polityka.Kod, KontoOperatora(ctx)))
	if err != nil {
		return PolitykaRetencjiBiblioteki{}, fmt.Errorf("dane: nieczytelna polityka retencji %q: %w",
			polityka.Kod, err)
	}
	return zapisana, nil
}

// PolitykiRetencji zwraca wszystkie polityki retencji, zawężone poziomem zasięgu, wprost z bazy danych.
func (r *repozytoriumBiblioteki) PolitykiRetencji(ctx context.Context,
	zasieg *string) ([]PolitykaRetencjiBiblioteki, error) {

	warunek := WarunekKonta
	argumenty := []any{KontoOperatora(ctx)}
	if zasieg != nil && *zasieg != "" {
		warunek += " AND zasieg = ?"
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

// UsunPolitykeRetencji zdejmuje politykę retencji po jej identyfikatorze trwałym z bazy danych repozytorium.
func (r *repozytoriumBiblioteki) UsunPolitykeRetencji(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunPolitykeRetencjiBiblioteki)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć polityki retencji %q: %w", kod, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć usuniętych polityk: %w", err)
	}
	return zdjete > 0, nil
}

// ZapiszUtrwalenie odkłada ślad po utrwaleniu archiwalnym wraz z wynikiem walidacji i jej pełnym zapisem.
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
		liczbaLogiczna(zadanie.Poprawne), zadanie.Raport, KontoOperatora(ctx))
	if err != nil {
		return ZadanieUtrwaleniaBiblioteki{}, fmt.Errorf("dane: nie można zapisać utrwalenia %q: %w",
			zadanie.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzUtrwalenieBiblioteki)
	if err != nil {
		return ZadanieUtrwaleniaBiblioteki{}, err
	}
	zapisane, err := odczytajUtrwalenieBiblioteki(odczyt.QueryRowContext(ctx, zadanie.Kod, KontoOperatora(ctx)))
	if err != nil {
		return ZadanieUtrwaleniaBiblioteki{}, fmt.Errorf("dane: nieczytelny zapis utrwalenia %q: %w",
			zadanie.Kod, err)
	}
	return zapisane, nil
}

// odczytajPolitykeRetencjiBiblioteki składa politykę wprost z jednego wiersza wyniku zapytania do bazy SQL.
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
