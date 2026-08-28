// Odpowiedzialność pliku: rejestr nagrań mowy (tabela `nagranie_mowy`, migracja 295) — obsługa przesyłania
// i pobierania nagrania mowy.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// NagranieMowy to wiersz tabeli `nagranie_mowy`, niosący opis nagrania mowy zapisanego na dysku rdzenia.
type NagranieMowy struct {
	Kod           string
	Sciezka       string
	TypTresci     string
	RozmiarBajtow int64
	DlugoscMs     *int64
	SesjaKod      *string
	OknoKod       *string
	Trwale        bool
	Utworzono     int64
	Wygasa        *int64
}

// RepozytoriumNagranMowy jest kontraktem rejestru nagrań mowy, określającym dostępne operacje tego wykazu.
type RepozytoriumNagranMowy interface {
	ZapiszNagranieMowy(ctx context.Context, nagranie NagranieMowy) (NagranieMowy, error)
	NagranieMowyPoSciezce(ctx context.Context, sciezka string) (NagranieMowy, error)
	NagraniaMowyWygasle(ctx context.Context, chwila int64) ([]NagranieMowy, error)
	UsunNagranieMowy(ctx context.Context, sciezka string) error
}

const (
	kolumnyNagraniaMowy = `identyfikator_zewnetrzny, sciezka, typ_tresci, rozmiar_bajtow,
	                       dlugosc_ms, sesja_kod, okno_kod, trwale, utworzono, wygasa`

	zapiszNagranieMowy = `INSERT INTO nagranie_mowy (` + kolumnyNagraniaMowy + `)
	                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                      ON CONFLICT(sciezka) DO UPDATE SET
	                          typ_tresci = excluded.typ_tresci,
	                          rozmiar_bajtow = excluded.rozmiar_bajtow,
	                          dlugosc_ms = excluded.dlugosc_ms,
	                          trwale = excluded.trwale,
	                          wygasa = excluded.wygasa`

	pobierzNagranieMowy = `SELECT ` + kolumnyNagraniaMowy + ` FROM nagranie_mowy WHERE sciezka = ?`

	pobierzNagraniaMowyWygasle = `SELECT ` + kolumnyNagraniaMowy + ` FROM nagranie_mowy
	                              WHERE trwale = 0 AND wygasa IS NOT NULL AND wygasa <= ?`

	usunNagranieMowy = `DELETE FROM nagranie_mowy WHERE sciezka = ?`
)

type repozytoriumNagranMowy struct {
	zapytania *zapytania
	db        *sql.DB
}

// noweRepozytoriumNagranMowy zakłada rejestr nagrań mowy nad wspólną bazą danych całego tego zestawu repozytoriów.
func noweRepozytoriumNagranMowy(z *zapytania, db *sql.DB) *repozytoriumNagranMowy {
	return &repozytoriumNagranMowy{zapytania: z, db: db}
}

// ZapiszNagranieMowy dopisuje wiersz opisujący nagranie mowy, które już leży zapisane na dysku danych rdzenia.
func (r *repozytoriumNagranMowy) ZapiszNagranieMowy(ctx context.Context,
	nagranie NagranieMowy) (NagranieMowy, error) {

	if strings.TrimSpace(nagranie.Sciezka) == "" {
		return NagranieMowy{}, fmt.Errorf("dane: nagranie mowy bez ścieżki pliku")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszNagranieMowy)
	if err != nil {
		return NagranieMowy{}, err
	}
	_, err = polecenie.ExecContext(ctx, nagranie.Kod, nagranie.Sciezka, nagranie.TypTresci,
		nagranie.RozmiarBajtow, liczbaDoKolumny(nagranie.DlugoscMs),
		tekstDoKolumny(nagranie.SesjaKod), tekstDoKolumny(nagranie.OknoKod),
		liczbaLogiczna(nagranie.Trwale), nagranie.Utworzono, liczbaDoKolumny(nagranie.Wygasa))
	if err != nil {
		return NagranieMowy{}, fmt.Errorf("dane: nie można zapisać nagrania mowy %q: %w",
			nagranie.Kod, err)
	}
	return r.NagranieMowyPoSciezce(ctx, nagranie.Sciezka)
}

// NagranieMowyPoSciezce zwraca wiersz nagrania. Brak wiersza jest sygnałem
// ErrBrakWiersza — odnośnik spoza rejestru nie jest nagraniem rdzenia.
func (r *repozytoriumNagranMowy) NagranieMowyPoSciezce(ctx context.Context,
	sciezka string) (NagranieMowy, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzNagranieMowy)
	if err != nil {
		return NagranieMowy{}, err
	}
	nagranie, err := odczytajNagranieMowy(polecenie.QueryRowContext(ctx, sciezka))
	if errors.Is(err, sql.ErrNoRows) {
		return NagranieMowy{}, fmt.Errorf("dane: nagranie mowy %q nie jest w rejestrze: %w",
			sciezka, ErrBrakWiersza)
	}
	return nagranie, err
}

// NagraniaMowyWygasle zwraca nagrania, których czas minął. Kasowaniem zajmuje
// się adapter, bo to on wie, gdzie leżą bajty.
func (r *repozytoriumNagranMowy) NagraniaMowyWygasle(ctx context.Context,
	chwila int64) ([]NagranieMowy, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzNagraniaMowyWygasle)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, chwila)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wygasłych nagrań mowy: %w", err)
	}
	defer wiersze.Close()

	lista := []NagranieMowy{}
	for wiersze.Next() {
		nagranie, err := odczytajNagranieMowy(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, nagranie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wygasłych nagrań mowy: %w", err)
	}
	return lista, nil
}

// UsunNagranieMowy wykreśla wiersz rejestru nagrań mowy; pliku nie kasuje — kasuje go ten, kto go zapisał.
func (r *repozytoriumNagranMowy) UsunNagranieMowy(ctx context.Context, sciezka string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, usunNagranieMowy)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, sciezka); err != nil {
		return fmt.Errorf("dane: nie można usunąć nagrania mowy %q z rejestru: %w", sciezka, err)
	}
	return nil
}

// odczytajNagranieMowy przekłada wiersz wyniku zapytania na pełną strukturę opisu nagrania mowy w bazie.
func odczytajNagranieMowy(s skaner) (NagranieMowy, error) {
	var nagranie NagranieMowy
	var sesja, okno sql.NullString
	var dlugosc, wygasa sql.NullInt64
	var trwale int
	err := s.Scan(&nagranie.Kod, &nagranie.Sciezka, &nagranie.TypTresci, &nagranie.RozmiarBajtow,
		&dlugosc, &sesja, &okno, &trwale, &nagranie.Utworzono, &wygasa)
	if err != nil {
		return NagranieMowy{}, err
	}
	nagranie.DlugoscMs = liczbaZKolumny(dlugosc)
	nagranie.SesjaKod = tekstZKolumny(sesja)
	nagranie.OknoKod = tekstZKolumny(okno)
	nagranie.Trwale = trwale == 1
	nagranie.Wygasa = liczbaZKolumny(wygasa)
	return nagranie, nil
}
