package dane

import (
	"context"
	"fmt"
)

// Dziennik zdarzeń zaczepów i zamknięcia tur; zapis następuje po zdarzeniu.

type ZdarzenieZaczepuWiersz struct {
	ID           int64
	Chwila       int64
	OknoKod      string
	WiadomoscKod string
	Rodzaj       string // 'zaczep_start' | 'zaczep_odpowiedz'
	ZaczepID     string
	Zaczep       string
	Zdarzenie    string
	Wynik        string
	KodWyjscia   *int64
	Tresc        string
	Ladunek      string
	SesjaCLI     string
}

type ZamkniecieTuryWiersz struct {
	Chwila             int64
	OknoKod            string
	WiadomoscKod       string
	Podtyp             string
	Blad               bool
	StanNadany         string // słownik `wiadomosc.stan`: zakonczona/zatrzymana/bledna
	KosztUSD           float64
	Tury               int
	CzasMs             int64
	Konto              string
	SesjaCLI           string
	TypyZdarzen        string // CSV w kolejności pierwszego wystąpienia
	LinieNierozpoznane int
}

type PrzelaczenieKanaluWiersz struct {
	Chwila       int64
	OknoKod      string
	WiadomoscKod string
	ZKanalu      string
	NaKanal      string
	Powod        string
}

type RepozytoriumZdarzenWykonawczych interface {
	ZapiszZaczep(ctx context.Context, z ZdarzenieZaczepuWiersz) error
	ZapiszZamkniecie(ctx context.Context, z ZamkniecieTuryWiersz) error
	ZapiszPrzelaczenie(ctx context.Context, z PrzelaczenieKanaluWiersz) error
	ZaczepyOkna(ctx context.Context, oknoKod string, granica int) ([]ZdarzenieZaczepuWiersz, error)
}

const (
	zapiszZaczepSQL = `INSERT INTO dziennik_zdarzen
		(chwila, okno_kod, wiadomosc_kod, rodzaj, zaczep_id, zaczep, zdarzenie,
		 wynik, kod_wyjscia, tresc, ladunek, sesja_cli, konto_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	// INSERT OR REPLACE, bo `wiadomosc_kod` jest UNIQUE: gdyby tura z jakiegoś
	// powodu domknęła się dwa razy, prawdą zostaje zamknięcie ostatnie.
	zapiszZamkniecieSQL = `INSERT OR REPLACE INTO zamkniecie_tury
		(chwila, okno_kod, wiadomosc_kod, podtyp, blad, stan_nadany, koszt_usd,
		 tury, czas_ms, konto, sesja_cli, typy_zdarzen, linie_nierozpoznane)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	zapiszPrzelaczenieSQL = `INSERT INTO przelaczenie_kanalu
		(chwila, okno_kod, wiadomosc_kod, z_kanalu, na_kanal, powod)
		VALUES (?, ?, ?, ?, ?, ?)`

	zaczepyOknaSQL = `SELECT id, chwila, okno_kod, wiadomosc_kod, rodzaj,
		 zaczep_id, zaczep, zdarzenie, wynik, kod_wyjscia, tresc, ladunek, sesja_cli
		FROM dziennik_zdarzen WHERE okno_kod = ? AND ` + WarunekKonta + `
		ORDER BY chwila DESC, id DESC LIMIT ?`
)

const granicaZaczepowDomyslna = 100

type repozytoriumZdarzenWykonawczych struct {
	zapytania *zapytania
}

func noweRepozytoriumZdarzenWykonawczych(zapytania *zapytania) RepozytoriumZdarzenWykonawczych {
	return &repozytoriumZdarzenWykonawczych{zapytania: zapytania}
}

func (r *repozytoriumZdarzenWykonawczych) ZapiszZaczep(ctx context.Context, z ZdarzenieZaczepuWiersz) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZaczepSQL)
	if err != nil {
		return err
	}
	var kodWyjscia any
	if z.KodWyjscia != nil {
		kodWyjscia = *z.KodWyjscia
	}
	if _, err := polecenie.ExecContext(ctx, z.Chwila, z.OknoKod, z.WiadomoscKod,
		z.Rodzaj, z.ZaczepID, z.Zaczep, z.Zdarzenie, z.Wynik, kodWyjscia,
		z.Tresc, z.Ladunek, z.SesjaCLI, KontoOperatora(ctx)); err != nil {
		return fmt.Errorf("dane: nie można zapisać zdarzenia zaczepu %q: %w", z.Zaczep, err)
	}
	return nil
}

func (r *repozytoriumZdarzenWykonawczych) ZapiszZamkniecie(ctx context.Context, z ZamkniecieTuryWiersz) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZamkniecieSQL)
	if err != nil {
		return err
	}
	blad := 0
	if z.Blad {
		blad = 1
	}
	if _, err := polecenie.ExecContext(ctx, z.Chwila, z.OknoKod, z.WiadomoscKod,
		z.Podtyp, blad, z.StanNadany, z.KosztUSD, z.Tury, z.CzasMs, z.Konto,
		z.SesjaCLI, z.TypyZdarzen, z.LinieNierozpoznane); err != nil {
		return fmt.Errorf("dane: nie można zapisać zamknięcia tury %q: %w", z.WiadomoscKod, err)
	}
	return nil
}

func (r *repozytoriumZdarzenWykonawczych) ZapiszPrzelaczenie(ctx context.Context, z PrzelaczenieKanaluWiersz) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPrzelaczenieSQL)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, z.Chwila, z.OknoKod, z.WiadomoscKod,
		z.ZKanalu, z.NaKanal, z.Powod); err != nil {
		return fmt.Errorf("dane: nie można zapisać przełączenia kanału %q→%q: %w", z.ZKanalu, z.NaKanal, err)
	}
	return nil
}

func (r *repozytoriumZdarzenWykonawczych) ZaczepyOkna(ctx context.Context, oknoKod string, granica int) ([]ZdarzenieZaczepuWiersz, error) {
	if granica <= 0 {
		granica = granicaZaczepowDomyslna
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zaczepyOknaSQL)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKod, KontoOperatora(ctx), granica)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zdarzeń zaczepów okna %q: %w", oknoKod, err)
	}
	defer wiersze.Close()
	var wynik []ZdarzenieZaczepuWiersz
	for wiersze.Next() {
		var w ZdarzenieZaczepuWiersz
		if err := wiersze.Scan(&w.ID, &w.Chwila, &w.OknoKod, &w.WiadomoscKod,
			&w.Rodzaj, &w.ZaczepID, &w.Zaczep, &w.Zdarzenie, &w.Wynik,
			&w.KodWyjscia, &w.Tresc, &w.Ladunek, &w.SesjaCLI); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz dziennika zdarzeń: %w", err)
		}
		wynik = append(wynik, w)
	}
	return wynik, wiersze.Err()
}
