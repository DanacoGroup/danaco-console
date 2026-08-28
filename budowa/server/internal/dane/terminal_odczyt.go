// Odpowiedzialność pliku: odczyt obszaru Terminal — karta powłoki terminala po
// kodzie i w komplecie oraz dziennik procesów zawężany parametrycznym filtrem SQL.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	kolumnyKartyTerminala = `kod, okno_kod, powloka, tytul, katalog_roboczy, stan, utworzono,
	                         cel_zdalny, port_zdalny, host_kod`

	pobierzKarteTerminala = `SELECT ` + kolumnyKartyTerminala + `
	                         FROM terminal_karta WHERE kod = ?`

	pobierzKartyTerminala = `SELECT ` + kolumnyKartyTerminala + `
	                         FROM terminal_karta ORDER BY utworzono ASC, id ASC`

	kolumnyProcesuTerminala = `kod, karta_kod, okno_kod, pid, pid_nadrzedny, polecenie,
	                           inicjator, stan, kod_wyjscia, uruchomiono, zakonczono`

	pobierzProcesTerminala = `SELECT ` + kolumnyProcesuTerminala + `
	                          FROM terminal_proces WHERE kod = ?`

	// Zawężenia wchodzą parametrem, nie sklejaniem tekstu: pusty parametr znaczy
	// „nie zawężaj”, więc jedno przygotowane zapytanie obsługuje cały filtr.
	pobierzProcesyTerminala = `SELECT ` + kolumnyProcesuTerminala + `
	                           FROM terminal_proces
	                           WHERE (? = '' OR okno_kod = ?)
	                             AND (? = '' OR karta_kod = ?)
	                             AND (? = '' OR stan = ?)
	                             AND (? = '' OR inicjator = ?)
	                           ORDER BY uruchomiono DESC, id DESC`
)

// repozytoriumTerminala obsługuje cały obszar modułu Terminal — dziennik kart
// i procesów wraz z wyposażeniem — i trzyma uchwyt bazy obok przygotowanych
// zapytań na potrzeby poleceń zapisu obejmujących więcej niż jedno zapytanie
// w jednej transakcji.
type repozytoriumTerminala struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumTerminala(z *zapytania, db *sql.DB) *repozytoriumTerminala {
	return &repozytoriumTerminala{zapytania: z, db: db}
}

// Karta zwraca kartę o wskazanym kodzie. Brak wiersza wraca jako ErrBrakWiersza
// — rdzeń odróżnia „nie ma takiej karty” od „odczyt się nie powiódł”.
func (r *repozytoriumTerminala) Karta(ctx context.Context, kod string) (KartaTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKarteTerminala)
	if err != nil {
		return KartaTerminala{}, err
	}
	karta, err := odczytajKarteTerminala(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return KartaTerminala{}, ErrBrakWiersza
	}
	if err != nil {
		return KartaTerminala{}, fmt.Errorf("dane: nieczytelny wiersz karty terminala %q: %w", kod, err)
	}
	return karta, nil
}

// Karty zwraca wszystkie karty w kolejności powstania. Czyta to rdzeń przy
// starcie, żeby odtworzyć profile powłok kart otwartych przed restartem.
func (r *repozytoriumTerminala) Karty(ctx context.Context) ([]KartaTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKartyTerminala)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kart terminala: %w", err)
	}
	defer wiersze.Close()

	karty := make([]KartaTerminala, 0, 8)
	for wiersze.Next() {
		karta, err := odczytajKarteTerminala(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz karty terminala: %w", err)
		}
		karty = append(karty, karta)
	}
	return karty, wiersze.Err()
}

// Proces zwraca wiersz dziennika procesu terminala o wskazanym kodzie, odczytany
// z tabeli terminal_proces wraz z jego bieżącym stanem i kodem wyjścia.
func (r *repozytoriumTerminala) Proces(ctx context.Context, kod string) (ProcesTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzProcesTerminala)
	if err != nil {
		return ProcesTerminala{}, err
	}
	proces, err := odczytajProcesTerminala(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ProcesTerminala{}, ErrBrakWiersza
	}
	if err != nil {
		return ProcesTerminala{}, fmt.Errorf("dane: nieczytelny wiersz procesu %q: %w", kod, err)
	}
	return proces, nil
}

// Procesy zwraca dziennik procesów terminala zawężony filtrem po oknie, karcie,
// stanie i inicjatorze, w kolejności od najnowszego przebiegu.
func (r *repozytoriumTerminala) Procesy(ctx context.Context, filtr FiltrProcesow) ([]ProcesTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzProcesyTerminala)
	if err != nil {
		return nil, err
	}
	stan, inicjator := string(filtr.Stan), string(filtr.Inicjator)
	wiersze, err := polecenie.QueryContext(ctx,
		filtr.OknoKod, filtr.OknoKod,
		filtr.KartaKod, filtr.KartaKod,
		stan, stan,
		inicjator, inicjator)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dziennika procesów: %w", err)
	}
	defer wiersze.Close()

	procesy := make([]ProcesTerminala, 0, 16)
	for wiersze.Next() {
		proces, err := odczytajProcesTerminala(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz procesu: %w", err)
		}
		procesy = append(procesy, proces)
	}
	return procesy, wiersze.Err()
}

// odczytajKarteTerminala składa wartość karty powłoki terminala z jednego wiersza
// wyniku zapytania, przypisując kolumny tabeli terminal_karta do pól struktury.
func odczytajKarteTerminala(wiersz interface{ Scan(...any) error }) (KartaTerminala, error) {
	var karta KartaTerminala
	err := wiersz.Scan(&karta.Kod, &karta.OknoKod, &karta.Powloka, &karta.Tytul,
		&karta.KatalogRoboczy, &karta.Stan, &karta.Utworzono,
		&karta.CelZdalny, &karta.PortZdalny, &karta.HostKod)
	return karta, err
}

// odczytajProcesTerminala składa wartość procesu terminala z jednego wiersza
// wyniku zapytania, przypisując kolumny tabeli terminal_proces do pól struktury.
func odczytajProcesTerminala(wiersz interface{ Scan(...any) error }) (ProcesTerminala, error) {
	var proces ProcesTerminala
	err := wiersz.Scan(&proces.Kod, &proces.KartaKod, &proces.OknoKod, &proces.Pid,
		&proces.PidNadrzedny, &proces.Polecenie, &proces.Inicjator, &proces.Stan,
		&proces.KodWyjscia, &proces.Uruchomiono, &proces.Zakonczono)
	return proces, err
}
