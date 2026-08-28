// Plik przyjmuje polecenie asystenta z okna monitorowania działań: zapisuje
// zlecenie i pierwszy wpis dziennika rozmowy w jednej transakcji.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PrzyjmijPolecenie zakłada zlecenie asystenta i zapisuje pierwszy wpis jego
// dziennika (rozmowy) w jednej transakcji. Zwraca oba wiersze po zapisie, ze
// wszystkimi polami uzupełnionymi przez bazę (identyfikator wewnętrzny,
// chwile czasu).
func (r *repozytoriumAsystenta) PrzyjmijPolecenie(ctx context.Context, zlecenie ZlecenieAsystenta,
	wpis WpisDziennikaAsystenta) (ZlecenieAsystenta, WpisDziennikaAsystenta, error) {

	if zlecenie.Kod == "" {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf("dane: polecenie asystenta bez kodu zlecenia")
	}
	if zlecenie.OknoKod == "" {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf("dane: polecenie asystenta %q bez okna", zlecenie.Kod)
	}
	if wpis.Kod == "" {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf("dane: polecenie asystenta %q bez kodu wpisu dziennika", zlecenie.Kod)
	}
	// Ani nagranie, ani tekst nie dają czym złożyć wpisu; to błąd żądania, nie zapis pustego wiersza.
	brakTresci := wpis.Tresc == ""
	brakNagrania := wpis.NagranieOdnosnik == nil || *wpis.NagranieOdnosnik == ""
	if brakTresci && brakNagrania {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf(
			"dane: polecenie asystenta %q bez treści i bez odnośnika do nagrania", zlecenie.Kod)
	}

	teraz := time.Now().UnixMilli()
	utworzonoZlecenia := zlecenie.Utworzono
	if utworzonoZlecenia == 0 {
		utworzonoZlecenia = teraz
	}
	utworzonoWpisu := wpis.Utworzono
	if utworzonoWpisu == 0 {
		utworzonoWpisu = teraz
	}
	// Wpis dziennika wisi na kodzie zlecenia przyjmowanego w tej transakcji, nie na kodzie parametru wpis.
	kodZlecenia := zlecenie.Kod

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapisZlecenia, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszZlecenieAsystenta)
		if err != nil {
			return err
		}
		if _, err := zapisZlecenia.ExecContext(ctx, zlecenie.Kod, zlecenie.OknoKod,
			tekstDoKolumny(zlecenie.Tytul), zlecenie.Stan, zlecenie.Droga,
			liczbaDoKolumny(zlecenie.EtapBiezacy), liczbaDoKolumny(zlecenie.LiczbaEtapow),
			liczbaDoKolumny(zlecenie.Priorytet), tekstDoKolumny(zlecenie.Wynik),
			tekstDoKolumny(zlecenie.ProfilKod), utworzonoZlecenia, teraz); err != nil {
			return fmt.Errorf("dane: nie można założyć zlecenia asystenta %q: %w", zlecenie.Kod, err)
		}

		zapisWpisu, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWpisDziennikaAsystenta)
		if err != nil {
			return err
		}
		if _, err := zapisWpisu.ExecContext(ctx, wpis.Kod, zlecenie.OknoKod, kodZlecenia,
			wpis.Rodzaj, wpis.Tresc, tekstDoKolumny(wpis.NagranieOdnosnik), utworzonoWpisu); err != nil {
			return fmt.Errorf("dane: nie można zapisać wpisu dziennika %q polecenia asystenta %q: %w",
				wpis.Kod, zlecenie.Kod, err)
		}
		return nil
	})
	if err != nil {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, err
	}

	zapisaneZlecenie, err := r.Zlecenie(ctx, zlecenie.Kod)
	if err != nil {
		return ZlecenieAsystenta{}, WpisDziennikaAsystenta{}, fmt.Errorf(
			"dane: nie można odczytać zapisanego zlecenia asystenta %q: %w", zlecenie.Kod, err)
	}
	zapisanyWpis := wpis
	zapisanyWpis.OknoKod = zlecenie.OknoKod
	zapisanyWpis.ZlecenieKod = &kodZlecenia
	zapisanyWpis.Utworzono = utworzonoWpisu
	return zapisaneZlecenie, zapisanyWpis, nil
}
