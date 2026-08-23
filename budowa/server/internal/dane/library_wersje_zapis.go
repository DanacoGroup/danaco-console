// Odpowiedzialność pliku: dołożenie kolejnej wersji pliku repozytorium wiedzy
// (`library.version.add`) — jedyna droga, którą historia pliku rośnie
// ponad wiersz założony przy wgraniu.
//
// Dlaczego osobny plik od `library_wersje.go`. Tamten plik odpowiada za odczyt
// historii i za przywrócenie wersji zastanej; ten za jej dołożenie. Rozdział
// idzie wzdłuż odpowiedzialności, a nie wzdłuż tabeli — SQL obu
// stron jest ten sam i mieszka nadal w `library_wersje.go`, żeby nie było
// dwóch prawd o jednym poleceniu.
//
// Dołożenie jest nierozdzielne. Wstawienie wiersza historii i przestawienie
// pliku macierzystego na tę wersję to jedna zmiana stanu: plik, którego
// `wersja_biezaca_id` wskazuje wiersz nieistniejący, albo historia z wersją,
// której plik nigdy nie przyjął, to schemat rozjechany w połowie.
// Stąd transakcja, tak samo jak przy `PrzywrocWersje`.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// DolozWersje wstawia kolejny wiersz historii pliku i przestawia na niego plik
// macierzysty: wskaźnik wersji bieżącej, sumę kontrolną, odwołanie do treści
// i rozmiar. Zwraca wersję po zapisie oraz plik po zmianie, bo `library.version.add`
// obiecuje kontraktem oba byty naraz.
//
// Liczba wersji poprzednich jest oddana wołającemu. Tabela nie ma kolumny
// numeru wersji — porządek historii daje `utworzono DESC, id DESC`.
// Rdzeń, który chce nazwać wersję jej kolejnością, bierze ją z `Wersje`; ten
// zapis niczego nie numeruje, bo numer nie jest tu bytem trwałym.
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
		if err := kod.QueryRowContext(ctx, plikID).Scan(&kodPliku); err != nil {
			return fmt.Errorf("dane: nie można odczytać kodu pliku %d po dołożeniu wersji: %w", plikID, err)
		}
		return nil
	})
	if err != nil {
		return WersjaPlikuBiblioteki{}, PlikBiblioteki{}, err
	}

	zapisana, err := r.jednaWersja(ctx, pobierzWersjePlikuBiblioteki, wersja.Kod, "wersja "+wersja.Kod)
	if err != nil {
		return WersjaPlikuBiblioteki{}, PlikBiblioteki{}, err
	}
	plik, err := r.Plik(ctx, kodPliku)
	if err != nil {
		return WersjaPlikuBiblioteki{}, PlikBiblioteki{}, err
	}
	return zapisana, plik, nil
}

// wstawWersje dokłada wiersz historii wewnątrz transakcji i oddaje jego klucz —
// klucz, a nie kod, bo to nim plik macierzysty wskazuje wersję bieżącą.
func (r *repozytoriumBiblioteki) wstawWersje(ctx context.Context, transakcja *sql.Tx,
	plikID int64, wersja WersjaPlikuBiblioteki) (int64, error) {

	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWersjePlikuBiblioteki)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, wersja.Kod, plikID, tekstDoKolumny(wersja.Etykieta),
		tekstDoKolumny(wersja.Autor), liczbaDoKolumny(wersja.RozmiarBajtow),
		tekstDoKolumny(wersja.SumaKontrolna), tekstDoKolumny(wersja.TrescOdwolanie))
	if err != nil {
		return 0, fmt.Errorf("dane: nie można dołożyć wersji %q pliku %d: %w", wersja.Kod, plikID, err)
	}
	wersjaID, err := wynik.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("dane: nie można ustalić klucza wersji %q pliku %d: %w", wersja.Kod, plikID, err)
	}
	return wersjaID, nil
}

// przestawPlikNaWersje aktualizuje plik macierzysty tym samym poleceniem, co
// przywrócenie wersji zastanej (`przywrocBiezacaWersje`) — bo to jest ta sama
// zmiana: „plik od teraz niesie treść tej wersji". Osobne polecenie o tym
// samym skutku byłoby drugą prawdą o jednym zapisie.
func (r *repozytoriumBiblioteki) przestawPlikNaWersje(ctx context.Context, transakcja *sql.Tx,
	plikID, wersjaID int64, wersja WersjaPlikuBiblioteki) error {

	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, przywrocBiezacaWersje)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wersjaID, tekstDoKolumny(wersja.SumaKontrolna),
		tekstDoKolumny(wersja.TrescOdwolanie), liczbaDoKolumny(wersja.RozmiarBajtow), plikID)
	if err != nil {
		return fmt.Errorf("dane: nie można przestawić pliku %d na wersję %q: %w", plikID, wersja.Kod, err)
	}
	return nil
}
