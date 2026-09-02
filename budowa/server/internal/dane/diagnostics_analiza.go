// Analiza modułu Diagnostics wraz z rekomendacjami z niej wyprowadzonymi: zapis, odczyt i wykaz.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	kolumnyAnalizy = `kod, okno_kod, zakres_od, zakres_do, podsumowanie, bledy,
	                  porownana_kod, utworzono`

	wstawAnalize = `INSERT INTO diagnostyka_analiza
	                (kod, okno_kod, zakres_od, zakres_do, podsumowanie, bledy,
	                 porownana_kod, utworzono, konto_id)
	                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzAnalize = `SELECT ` + kolumnyAnalizy + ` FROM diagnostyka_analiza WHERE kod = ? AND ` + WarunekKonta

	kolumnyRekomendacji = `kod, analiza_kod, blad_kod, tytul, szczegol, priorytet, stan,
	                       sciezka, poprawka, utworzono`

	wstawRekomendacje = `INSERT INTO diagnostyka_rekomendacja
	                     (kod, analiza_kod, blad_kod, tytul, szczegol, priorytet, stan,
	                      sciezka, poprawka, utworzono)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzRekomendacje = `SELECT ` + kolumnyRekomendacji + ` FROM diagnostyka_rekomendacja
	                       WHERE (? = '' OR analiza_kod = ?)
	                         AND (? = '' OR stan = ?)
	                         AND (? = '' OR priorytet = ?)
	                         AND EXISTS (SELECT 1 FROM diagnostyka_analiza a
	                                     WHERE a.kod = diagnostyka_rekomendacja.analiza_kod AND ` + WarunekKonta + `)
	                       ORDER BY utworzono DESC, id DESC
	                       LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`
)

// ZapiszAnalize utrwala migawkę diagnostyczną wraz z jej rekomendacjami w jednej transakcji zapisu do bazy danych.
func (r *repozytoriumDiagnostyki) ZapiszAnalize(ctx context.Context,
	analiza AnalizaDiagnostyczna, rekomendacje []RekomendacjaDiagnostyczna) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawAnalize)
		if err != nil {
			return err
		}
		_, err = polecenie.ExecContext(ctx, analiza.Kod, tekstDoKolumny(analiza.OknoKod),
			liczbaDoKolumny(analiza.ZakresOd), liczbaDoKolumny(analiza.ZakresDo),
			tekstDoKolumny(analiza.Podsumowanie), tekstDoKolumny(analiza.Bledy),
			tekstDoKolumny(analiza.PorownanaKod), analiza.Utworzono, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać analizy %q: %w", analiza.Kod, err)
		}

		wpis, err := r.zapytania.wTransakcji(ctx, transakcja, wstawRekomendacje)
		if err != nil {
			return err
		}
		for _, rekomendacja := range rekomendacje {
			_, err := wpis.ExecContext(ctx, rekomendacja.Kod, analiza.Kod,
				tekstDoKolumny(rekomendacja.BladKod), rekomendacja.Tytul,
				tekstDoKolumny(rekomendacja.Szczegol), string(rekomendacja.Priorytet),
				string(rekomendacja.Stan), tekstDoKolumny(rekomendacja.Sciezka),
				tekstDoKolumny(rekomendacja.Poprawka), rekomendacja.Utworzono)
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać rekomendacji %q: %w", rekomendacja.Kod, err)
			}
		}
		return nil
	})
}

// Analiza zwraca migawkę po kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumDiagnostyki) Analiza(ctx context.Context, kod string) (AnalizaDiagnostyczna, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzAnalize)
	if err != nil {
		return AnalizaDiagnostyczna{}, err
	}
	analiza, err := odczytajAnalize(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return AnalizaDiagnostyczna{}, ErrBrakWiersza
	}
	if err != nil {
		return AnalizaDiagnostyczna{}, fmt.Errorf("dane: nieczytelny wiersz analizy %q: %w", kod, err)
	}
	return analiza, nil
}

// Rekomendacje zwraca wykaz rekomendacji diagnostycznych zawężony filtrem wyszukiwania, uporządkowany od najświeższej.
func (r *repozytoriumDiagnostyki) Rekomendacje(ctx context.Context,
	filtr FiltrRekomendacji) ([]RekomendacjaDiagnostyczna, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzRekomendacje)
	if err != nil {
		return nil, err
	}
	stan, priorytet := string(filtr.Stan), string(filtr.Priorytet)
	granica := filtr.Granica
	if granica <= 0 {
		granica = granicaDziennika
	}
	wiersze, err := polecenie.QueryContext(ctx, filtr.AnalizaKod, filtr.AnalizaKod,
		stan, stan, priorytet, priorytet, KontoOperatora(ctx), granica, granica)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać rekomendacji: %w", err)
	}
	defer wiersze.Close()

	rekomendacje := make([]RekomendacjaDiagnostyczna, 0, 16)
	for wiersze.Next() {
		rekomendacja, err := odczytajRekomendacje(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz rekomendacji: %w", err)
		}
		rekomendacje = append(rekomendacje, rekomendacja)
	}
	return rekomendacje, wiersze.Err()
}

// odczytajAnalize składa migawkę diagnostyczną ze struktury opartej na jednym wierszu wyniku zapytania do bazy danych.
func odczytajAnalize(s skaner) (AnalizaDiagnostyczna, error) {
	var analiza AnalizaDiagnostyczna
	err := s.Scan(&analiza.Kod, &analiza.OknoKod, &analiza.ZakresOd, &analiza.ZakresDo,
		&analiza.Podsumowanie, &analiza.Bledy, &analiza.PorownanaKod, &analiza.Utworzono)
	return analiza, err
}

// odczytajRekomendacje składa rekomendację diagnostyczną ze struktury opartej na jednym wierszu wyniku zapytania.
func odczytajRekomendacje(s skaner) (RekomendacjaDiagnostyczna, error) {
	var rekomendacja RekomendacjaDiagnostyczna
	err := s.Scan(&rekomendacja.Kod, &rekomendacja.AnalizaKod, &rekomendacja.BladKod,
		&rekomendacja.Tytul, &rekomendacja.Szczegol, &rekomendacja.Priorytet,
		&rekomendacja.Stan, &rekomendacja.Sciezka, &rekomendacja.Poprawka,
		&rekomendacja.Utworzono)
	return rekomendacja, err
}
