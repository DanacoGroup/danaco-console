// Plik prowadzi sugestie porządkujące biblioteki; sugestia jest bytem trwałym, nie wynikiem oddanym i zapomnianym
// — klasyfikacja wsadowa ją wytwarza, a decyzja Operatora zapada osobnym żądaniem, często znacznie później.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SugestiaBiblioteki to wiersz tabeli `sugestia_biblioteki` niosący propozycję porządkującą wraz z jej stanem.
type SugestiaBiblioteki struct {
	ID             int64
	Kod            string
	PlikKod        string
	Rodzaj         string
	Wartosc        *string
	Uzasadnienie   string
	Pewnosc        *int64
	Stan           string
	Rozstrzygnieto *string
	Utworzono      string
}

// Stany sugestii porządkującej biblioteki: oczekująca na decyzję, przyjęta albo odrzucona przez Operatora.
const (
	StanSugestiiOczekujaca = "oczekujaca"
	StanSugestiiPrzyjeta   = "przyjeta"
	StanSugestiiOdrzucona  = "odrzucona"
)

const (
	kolumnySugestiiBiblioteki = `id, identyfikator_zewnetrzny, plik_kod, rodzaj, wartosc,
	                             uzasadnienie, pewnosc, stan, rozstrzygnieto, utworzono`

	zapiszSugestieBiblioteki = `INSERT INTO sugestia_biblioteki
	                            (identyfikator_zewnetrzny, plik_kod, rodzaj, wartosc,
	                             uzasadnienie, pewnosc, stan)
	                            VALUES (?, ?, ?, ?, ?, ?, ?)
	                            ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                wartosc = excluded.wartosc,
	                                uzasadnienie = excluded.uzasadnienie,
	                                pewnosc = excluded.pewnosc,
	                                stan = excluded.stan`

	pobierzSugestieBiblioteki = `SELECT ` + kolumnySugestiiBiblioteki + ` FROM sugestia_biblioteki
	                             WHERE identyfikator_zewnetrzny = ?`

	rozstrzygnijSugestieBiblioteki = `UPDATE sugestia_biblioteki
	                                  SET stan = ?,
	                                      rozstrzygnieto = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                  WHERE identyfikator_zewnetrzny = ? AND stan = 'oczekujaca'`
)

// ZapiszSugestie zakłada sugestię porządkującą albo zmienia zastaną sugestię tego samego kodu w bazie.
func (r *repozytoriumBiblioteki) ZapiszSugestie(ctx context.Context,
	sugestia SugestiaBiblioteki) (SugestiaBiblioteki, error) {

	if strings.TrimSpace(sugestia.Kod) == "" {
		return SugestiaBiblioteki{}, fmt.Errorf("dane: sugestia biblioteki bez identyfikatora")
	}
	if strings.TrimSpace(sugestia.Uzasadnienie) == "" {
		return SugestiaBiblioteki{}, fmt.Errorf("dane: sugestia %q bez uzasadnienia", sugestia.Kod)
	}
	stan := sugestia.Stan
	if stan == "" {
		stan = StanSugestiiOczekujaca
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSugestieBiblioteki)
	if err != nil {
		return SugestiaBiblioteki{}, err
	}
	_, err = polecenie.ExecContext(ctx, sugestia.Kod, sugestia.PlikKod, sugestia.Rodzaj,
		tekstDoKolumny(sugestia.Wartosc), sugestia.Uzasadnienie,
		liczbaDoKolumny(sugestia.Pewnosc), stan)
	if err != nil {
		return SugestiaBiblioteki{}, fmt.Errorf("dane: nie można zapisać sugestii %q: %w",
			sugestia.Kod, err)
	}
	return r.Sugestia(ctx, sugestia.Kod)
}

// Sugestia zwraca sugestię porządkującą o wskazanym kodzie zewnętrznym wprost z bazy danych repozytorium.
func (r *repozytoriumBiblioteki) Sugestia(ctx context.Context, kod string) (SugestiaBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSugestieBiblioteki)
	if err != nil {
		return SugestiaBiblioteki{}, err
	}
	sugestia, err := odczytajSugestieBiblioteki(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SugestiaBiblioteki{}, ErrBrakWiersza
	}
	if err != nil {
		return SugestiaBiblioteki{}, fmt.Errorf("dane: nieczytelna sugestia %q: %w", kod, err)
	}
	return sugestia, nil
}

// Sugestie zwraca sugestie oczekujące na decyzję, od najnowszej; sugestia rozstrzygnięta jest wpisem historii.
func (r *repozytoriumBiblioteki) Sugestie(ctx context.Context, plikKod *string,
	rodzaje []string, limit int) ([]SugestiaBiblioteki, int, error) {

	warunki := []string{"stan = 'oczekujaca'"}
	argumenty := []any{}
	if plikKod != nil && *plikKod != "" {
		warunki = append(warunki, "plik_kod = ?")
		argumenty = append(argumenty, *plikKod)
	}
	if len(rodzaje) > 0 {
		miejsca := make([]string, 0, len(rodzaje))
		for _, rodzaj := range rodzaje {
			miejsca = append(miejsca, "?")
			argumenty = append(argumenty, rodzaj)
		}
		warunki = append(warunki, "rodzaj IN ("+strings.Join(miejsca, ", ")+")")
	}
	warunek := strings.Join(warunki, " AND ")

	zapytanie := `SELECT ` + kolumnySugestiiBiblioteki + ` FROM sugestia_biblioteki
	              WHERE ` + warunek + ` ORDER BY utworzono DESC, id DESC LIMIT ?`
	wiersze, err := r.db.QueryContext(ctx, zapytanie,
		append(append([]any{}, argumenty...), granicaWykazu(limit))...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać sugestii: %w", err)
	}
	defer wiersze.Close()

	lista := []SugestiaBiblioteki{}
	for wiersze.Next() {
		sugestia, err := odczytajSugestieBiblioteki(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz sugestii: %w", err)
		}
		lista = append(lista, sugestia)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt sugestii: %w", err)
	}

	var lacznie int
	zapytanieLiczby := `SELECT COUNT(*) FROM sugestia_biblioteki WHERE ` + warunek
	if err := r.db.QueryRowContext(ctx, zapytanieLiczby, argumenty...).Scan(&lacznie); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć sugestii: %w", err)
	}
	return lista, lacznie, nil
}

// RozstrzygnijSugestie znakuje sugestie jako przyjęte albo odrzucone i oddaje
// liczbę wierszy, które zmiana objęła. Sugestia już rozstrzygnięta nie liczy się
// po raz drugi — warunek zapytania pilnuje tego zamiast wołającego.
func (r *repozytoriumBiblioteki) RozstrzygnijSugestie(ctx context.Context, kody []string,
	przyjeto bool) (int, error) {

	stan := StanSugestiiOdrzucona
	if przyjeto {
		stan = StanSugestiiPrzyjeta
	}
	polecenie, err := r.zapytania.przygotuj(ctx, rozstrzygnijSugestieBiblioteki)
	if err != nil {
		return 0, err
	}
	rozstrzygniete := 0
	for _, kod := range kody {
		wynik, err := polecenie.ExecContext(ctx, stan, kod)
		if err != nil {
			return 0, fmt.Errorf("dane: nie można rozstrzygnąć sugestii %q: %w", kod, err)
		}
		zmienione, err := wynik.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("dane: nie można policzyć rozstrzygniętych sugestii: %w", err)
		}
		rozstrzygniete += int(zmienione)
	}
	return rozstrzygniete, nil
}

// odczytajSugestieBiblioteki składa sugestię wprost z jednego wiersza wyniku zapytania do bazy danych.
func odczytajSugestieBiblioteki(wiersz skaner) (SugestiaBiblioteki, error) {
	var sugestia SugestiaBiblioteki
	var wartosc, rozstrzygnieto sql.NullString
	var pewnosc sql.NullInt64
	err := wiersz.Scan(&sugestia.ID, &sugestia.Kod, &sugestia.PlikKod, &sugestia.Rodzaj,
		&wartosc, &sugestia.Uzasadnienie, &pewnosc, &sugestia.Stan, &rozstrzygnieto,
		&sugestia.Utworzono)
	if err != nil {
		return SugestiaBiblioteki{}, err
	}
	sugestia.Wartosc, sugestia.Rozstrzygnieto = tekstZKolumny(wartosc), tekstZKolumny(rozstrzygnieto)
	sugestia.Pewnosc = liczbaZKolumny(pewnosc)
	return sugestia, nil
}
