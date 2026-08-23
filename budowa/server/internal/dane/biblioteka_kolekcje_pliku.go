// Przynależność pliku do kolekcji widziana od strony pliku (tabela
// `przypisanie_kolekcji_biblioteki`): odczyt kolekcji danego pliku oraz
// ustawienie kompletu kolekcji pliku, czyli zapis, który potrafi z kolekcji
// zdjąć. Kolekcja jako byt (`kolekcja_biblioteki`) i przypisanie widziane od
// strony kolekcji (`PrzypiszDoKolekcji`) leżą w `library_kolekcje.go`.
//
// `UstawKolekcjePliku` ma semantykę „ustaw": stan po zapisie jest wykazem
// z żądania, więc `library.tag.set` zdejmuje plik z kolekcji pominiętych
// w wykazie, zamiast tylko dokładać nowe.
//
// Ustawienie jest wymianą w jednej transakcji, tak samo jak `UstawEtykiety`
// (`library_kolekcje.go`). Zapis przerwany w połowie nie zostawia pliku
// w stanie przejściowym, widocznym dla wykazu filtrowanego po kolekcji.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	// Kolekcje pliku wychodzą kodami, bo tym plik i kolekcja wychodzą
	// kontraktem (`LibraryFile.collectionIds`) — klucz wiersza zostaje
	// wewnątrz warstwy danych.
	listaKolekcjiPliku = `SELECT k.identyfikator_zewnetrzny
	                      FROM przypisanie_kolekcji_biblioteki pk
	                      JOIN kolekcja_biblioteki k ON k.id = pk.kolekcja_id
	                      WHERE pk.plik_id = ?
	                      ORDER BY k.identyfikator_zewnetrzny`

	usunPrzypisaniaPliku = `DELETE FROM przypisanie_kolekcji_biblioteki WHERE plik_id = ?`

	idKolekcjiPoKodzie = `SELECT id FROM kolekcja_biblioteki WHERE identyfikator_zewnetrzny = ?`
)

// KolekcjePliku zwraca kody kolekcji, do których plik należy. Porządek jest
// ustalony (po kodzie), żeby odpowiedź kontraktu nie zmieniała kolejności
// między dwoma odczytami tego samego stanu.
func (r *repozytoriumBiblioteki) KolekcjePliku(ctx context.Context, plikID int64) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKolekcjiPliku)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, plikID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kolekcji pliku %d: %w", plikID, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var kod string
		if err := wiersze.Scan(&kod); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kolekcji pliku %d: %w", plikID, err)
		}
		lista = append(lista, kod)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kolekcji pliku %d: %w", plikID, err)
	}
	return lista, nil
}

// UstawKolekcjePliku czyni wykaz kolekcji pliku dokładnie takim, jaki podano:
// kolekcje spoza wykazu zostają zdjęte, kolekcje z wykazu dopięte. Zwraca stan
// po zapisie odczytany z bazy, a nie powtórzone żądanie.
//
// Kolekcja nieznana jest odmową całości, tak samo jak przy `PrzypiszDoKolekcji`:
// wykaz z kodem, którego nie ma, wskazuje przynależność nieosiągalną,
// a wykonanie reszty zdjęłoby plik z kolekcji zastanych na podstawie żądania
// zrozumianego tylko częściowo. Sprawdzenie idzie przed usunięciem, wewnątrz
// tej samej transakcji, więc odmowa nie zostawia pliku bez przypisań.
func (r *repozytoriumBiblioteki) UstawKolekcjePliku(ctx context.Context,
	kodPliku string, kodyKolekcji []string) ([]string, error) {

	plik, err := r.Plik(ctx, kodPliku)
	if err != nil {
		return nil, err
	}

	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		identyfikatory, err := r.kluczeKolekcji(ctx, transakcja, kodyKolekcji)
		if err != nil {
			return err
		}
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunPrzypisaniaPliku)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, plik.ID); err != nil {
			return fmt.Errorf("dane: nie można zdjąć pliku %q z kolekcji: %w", kodPliku, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawPrzypisanieKolekcji)
		if err != nil {
			return err
		}
		for _, kolekcjaID := range identyfikatory {
			if _, err := wstawienie.ExecContext(ctx, kolekcjaID, plik.ID); err != nil {
				return fmt.Errorf("dane: nie można przypisać pliku %q do kolekcji %d: %w",
					kodPliku, kolekcjaID, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.KolekcjePliku(ctx, plik.ID)
}

// kluczeKolekcji przekłada kody kolekcji na klucze wierszy i odmawia przy
// pierwszym kodzie bez wiersza (ErrBrakWiersza — adapter odróżni „nie ma
// takiej kolekcji" od awarii odczytu). Kod pusty jest pominięciem, nie
// odmową: puste pole w wykazie żądania nie wskazuje żadnej kolekcji, więc nie
// ma czego nie znaleźć.
func (r *repozytoriumBiblioteki) kluczeKolekcji(ctx context.Context, transakcja *sql.Tx,
	kody []string) ([]int64, error) {

	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, idKolekcjiPoKodzie)
	if err != nil {
		return nil, err
	}
	klucze := make([]int64, 0, len(kody))
	for _, kod := range kody {
		if kod == "" {
			continue
		}
		var kolekcjaID int64
		err := polecenie.QueryRowContext(ctx, kod).Scan(&kolekcjaID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBrakWiersza
		}
		if err != nil {
			return nil, fmt.Errorf("dane: nie można odnaleźć kolekcji %q: %w", kod, err)
		}
		klucze = append(klucze, kolekcjaID)
	}
	return klucze, nil
}
