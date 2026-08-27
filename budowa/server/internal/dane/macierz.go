// Odpowiedzialność pliku: odczyt macierzy widoczności modułów w środowiskach
// (tabela `srodowisko_modul`). Macierz jest konfiguracją: wiersz mówi, czy moduł stoi w bocznej nawigacji
// danego środowiska i na którym miejscu.
package dane

import (
	"context"
	"fmt"
)

// WierszMacierzy to jeden wiersz tabeli `srodowisko_modul` wraz z kodami obu stron pary środowiska i modułu.
type WierszMacierzy struct {
	SrodowiskoID  int64
	SrodowiskoKod string
	ModulID       int64
	ModulKod      string
	// Kolejnosc to pozycja modułu w bocznej nawigacji środowiska; dla pary niewidocznej nie niesie treści.
	Kolejnosc int64
	Widoczny  bool
}

// RepozytoriumMacierzy jest kontraktem odczytu macierzy widoczności modułów; zapisu tu nie ma świadomie.
type RepozytoriumMacierzy interface {
	// Pelna zwraca wszystkie wiersze macierzy — widoczne i niewidoczne — w kolejności nawigacji środowisk.
	Pelna(ctx context.Context) ([]WierszMacierzy, error)
	// KodySrodowiskModulow zwraca odwzorowanie modułu na kody środowisk, w których jest widoczny.
	KodySrodowiskModulow(ctx context.Context) (map[int64][]string, error)
}

const (
	kolumnyMacierzy = `s.id, s.kod, m.id, m.kod, sm.kolejnosc, sm.widoczny`

	zlaczenieMacierzy = ` FROM srodowisko_modul sm
	                      JOIN srodowisko s ON s.id = sm.srodowisko_id
	                      JOIN modul m      ON m.id = sm.modul_id`

	// Porządek środowisk jest ten sam, co przy odczycie środowisk, a wewnątrz środowiska ten sam, co przy
	// odczycie modułów środowiska.
	porzadekMacierzy = ` ORDER BY s.kolejnosc, s.kod, sm.kolejnosc, m.kod`

	pelnaMacierz = `SELECT ` + kolumnyMacierzy + zlaczenieMacierzy + porzadekMacierzy

	widocznaMacierz = `SELECT ` + kolumnyMacierzy + zlaczenieMacierzy +
		` WHERE sm.widoczny = 1` + porzadekMacierzy
)

type repozytoriumMacierzy struct {
	zapytania *zapytania
}

func noweRepozytoriumMacierzy(z *zapytania) *repozytoriumMacierzy {
	return &repozytoriumMacierzy{zapytania: z}
}

// Pelna zwraca komplet wierszy macierzy widoczności modułów we wszystkich środowiskach całej platformy.
func (r *repozytoriumMacierzy) Pelna(ctx context.Context) ([]WierszMacierzy, error) {
	return r.wykaz(ctx, pelnaMacierz, "macierzy widoczności")
}

// KodySrodowiskModulow składa odwzorowanie modułu na kody środowisk, w których jest widoczny, jednym zapytaniem.
func (r *repozytoriumMacierzy) KodySrodowiskModulow(ctx context.Context) (map[int64][]string, error) {
	wiersze, err := r.wykaz(ctx, widocznaMacierz, "widocznej macierzy")
	if err != nil {
		return nil, err
	}
	kody := map[int64][]string{}
	for _, wiersz := range wiersze {
		kody[wiersz.ModulID] = append(kody[wiersz.ModulID], wiersz.SrodowiskoKod)
	}
	return kody, nil
}

// wykaz wykonuje zapytanie do bazy danych zwracające wiele wierszy macierzy widoczności modułów platformy.
func (r *repozytoriumMacierzy) wykaz(ctx context.Context, zapytanie, opis string) ([]WierszMacierzy, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać %s: %w", opis, err)
	}
	defer wiersze.Close()

	lista := []WierszMacierzy{}
	for wiersze.Next() {
		wiersz, err := odczytajWierszMacierzy(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, wiersz)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt %s: %w", opis, err)
	}
	return lista, nil
}

// odczytajWierszMacierzy składa pełną strukturę wiersza macierzy z jednego wyniku zapytania do bazy danych.
func odczytajWierszMacierzy(wiersz skaner) (WierszMacierzy, error) {
	var para WierszMacierzy
	var widoczny int
	if err := wiersz.Scan(
		&para.SrodowiskoID,
		&para.SrodowiskoKod,
		&para.ModulID,
		&para.ModulKod,
		&para.Kolejnosc,
		&widoczny,
	); err != nil {
		return WierszMacierzy{}, err
	}
	para.Widoczny = widoczny != 0
	return para, nil
}
