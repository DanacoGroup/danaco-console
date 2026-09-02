// Odpowiedzialność pliku: dołożenie kolejnej wersji pliku repozytorium wiedzy — jedyna droga, którą historia
// pliku rośnie ponad wiersz założony przy wgraniu.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *repozytoriumBiblioteki) DolozWersje(ctx context.Context, plikID int64,
	wersja WersjaPlikuBiblioteki) (WersjaPlikuBiblioteki, PlikBiblioteki, error) {

	if wersja.Kod == "" || plikID == 0 {
		return WersjaPlikuBiblioteki{}, PlikBiblioteki{},
			fmt.Errorf("dane: dołożenie wersji bez identyfikatora albo bez pliku")
	}

	var kodPliku string
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wersjaID, err := r.wstawWersje(ctx, transakcja, plikID, wersja)
		if err != nil {
			return err
		}
		if err := r.przestawPlikNaWersje(ctx, transakcja, plikID, wersjaID, wersja); err != nil {
			return err
		}
		kod, err := r.zapytania.wTransakcji(ctx, transakcja, kodPlikuPoId)
		if err != nil {
			return err
		}
		if err := kod.QueryRowContext(ctx, plikID, KontoOperatora(ctx)).Scan(&kodPliku); err != nil {
			return fmt.Errorf("dane: nie można odczytać kodu pliku %d po dołożeniu wersji: %w", plikID, err)
		}
		return nil
	})
	if err != nil {
		return WersjaPlikuBiblioteki{}, PlikBiblioteki{}, err
	}

	zapisana, err := r.jednaWersja(ctx, pobierzWersjePlikuBiblioteki, []any{wersja.Kod, KontoOperatora(ctx)},
		"wersja "+wersja.Kod)
	if err != nil {
		return WersjaPlikuBiblioteki{}, PlikBiblioteki{}, err
	}
	plik, err := r.Plik(ctx, kodPliku)
	if err != nil {
		return WersjaPlikuBiblioteki{}, PlikBiblioteki{}, err
	}
	return zapisana, plik, nil
}

// wstawWersje oddaje klucz wiersza, nie kod — kluczem plik macierzysty wskazuje wersję bieżącą.
func (r *repozytoriumBiblioteki) wstawWersje(ctx context.Context, transakcja *sql.Tx,
	plikID int64, wersja WersjaPlikuBiblioteki) (int64, error) {

	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWersjePlikuBiblioteki)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, wersja.Kod, tekstDoKolumny(wersja.Etykieta),
		tekstDoKolumny(wersja.Autor), liczbaDoKolumny(wersja.RozmiarBajtow),
		tekstDoKolumny(wersja.SumaKontrolna), tekstDoKolumny(wersja.TrescOdwolanie),
		plikID, KontoOperatora(ctx))
	if err != nil {
		return 0, fmt.Errorf("dane: nie można dołożyć wersji %q pliku %d: %w", wersja.Kod, plikID, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "plik biblioteki", fmt.Sprint(plikID)); err != nil {
		return 0, err
	}
	wersjaID, err := wynik.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("dane: nie można ustalić klucza wersji %q pliku %d: %w", wersja.Kod, plikID, err)
	}
	return wersjaID, nil
}

func (r *repozytoriumBiblioteki) przestawPlikNaWersje(ctx context.Context, transakcja *sql.Tx,
	plikID, wersjaID int64, wersja WersjaPlikuBiblioteki) error {

	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, przywrocBiezacaWersje)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, wersjaID, tekstDoKolumny(wersja.SumaKontrolna),
		tekstDoKolumny(wersja.TrescOdwolanie), liczbaDoKolumny(wersja.RozmiarBajtow),
		plikID, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można przestawić pliku %d na wersję %q: %w", plikID, wersja.Kod, err)
	}
	return sprawdzTrafienieZapisu(wynik, "plik biblioteki", fmt.Sprint(plikID))
}
