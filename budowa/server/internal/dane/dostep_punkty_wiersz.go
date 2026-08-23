// Odpowiedzialność pliku: przekład wiersza tabeli `punkt_dostepu` na strukturę
// PunktDostepu i z powrotem oraz słownictwo trybu mostu (`argument_trybu_mostu`).
//
// Adres mostu nie ma własnej kolumny — składa się go z użytkownika, hosta i portu.
// Gdyby był osobnym polem, wiersz mógłby mieć adres sprzeczny z własnymi danymi
// połączenia, a taki rozjazd nie ujawniłby się aż do nieudanego uruchomienia.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"danacoconsole/shared"
)

const (
	listaArgumentowTrybu = `SELECT tryb, argument FROM argument_trybu_mostu
	                        WHERE punkt_dostepu_id = ? ORDER BY tryb`

	usunArgumentyTrybu = `DELETE FROM argument_trybu_mostu WHERE punkt_dostepu_id = ?`

	wstawArgumentTrybu = `INSERT INTO argument_trybu_mostu (punkt_dostepu_id, tryb, argument)
	                      VALUES (?, ?, ?)`
)

// wartosciPunktuBazy to komplet wartości kolumn słownikowych punktu dostępu.
type wartosciPunktuBazy struct {
	rodzaj string
	tryb   string
	stan   string
}

// wartosciPunktu przekłada pola wyliczeniowe struktury na wartości kolumn.
func wartosciPunktu(punkt PunktDostepu) (wartosciPunktuBazy, error) {
	var wartosci wartosciPunktuBazy
	var err error
	if wartosci.rodzaj, err = rodzajPunktuNaBaze(punkt.Rodzaj); err != nil {
		return wartosci, err
	}
	if wartosci.tryb, err = trybDostepuNaBaze(punkt.TrybDomyslny, "punkt_dostepu.tryb_domyslny"); err != nil {
		return wartosci, err
	}
	if wartosci.stan, err = stanPunktuNaBaze(punkt.Stan); err != nil {
		return wartosci, err
	}
	return wartosci, nil
}

// odczytajPunktDostepu składa strukturę z jednego wiersza wyniku. Listy podrzędne
// dokłada repozytorium — leżą w osobnych tabelach.
func odczytajPunktDostepu(wiersz skaner) (PunktDostepu, error) {
	var punkt PunktDostepu
	var urzadzenie sql.NullInt64
	var odwolanie, sprawdzono sql.NullString
	var rodzaj, tryb, stan string
	var aktywny int
	err := wiersz.Scan(&punkt.ID, &punkt.Kod, &punkt.Nazwa, &punkt.Opis, &rodzaj, &urzadzenie,
		&punkt.Host, &punkt.Port, &punkt.Uzytkownik, &punkt.SciezkaKlucza, &punkt.PolecenieStartu,
		&punkt.NazwaMostu, &odwolanie, &tryb, &stan, &sprawdzono, &aktywny, &punkt.Kolejnosc,
		&punkt.Utworzono, &punkt.Zaktualizowano)
	if err != nil {
		return PunktDostepu{}, err
	}
	punkt.UrzadzenieID = liczbaZKolumny(urzadzenie)
	punkt.PoswiadczenieOdwolanie = tekstZKolumny(odwolanie)
	punkt.Sprawdzono = tekstZKolumny(sprawdzono)
	punkt.Aktywny = aktywny != 0
	punkt.Korzenie = []string{}
	punkt.ArgumentyTrybu = map[shared.AccessMode]string{}
	return uzupelnijSlownikiPunktu(punkt, rodzaj, tryb, stan)
}

// uzupelnijSlownikiPunktu przekłada wartości kolumn słownikowych na kontrakt.
func uzupelnijSlownikiPunktu(punkt PunktDostepu, rodzaj, tryb, stan string) (PunktDostepu, error) {
	var err error
	if punkt.Rodzaj, err = rodzajPunktuZBazy(rodzaj); err != nil {
		return PunktDostepu{}, err
	}
	if punkt.TrybDomyslny, err = trybDostepuZBazy(tryb, "punkt_dostepu.tryb_domyslny"); err != nil {
		return PunktDostepu{}, err
	}
	if punkt.Stan, err = stanPunktuZBazy(stan); err != nil {
		return PunktDostepu{}, err
	}
	return punkt, nil
}

// argumentyTrybuMostu zwraca słownictwo trybu jednego punktu. Brak wierszy nie
// jest błędem — most `mcp-danaco-pulpit-console` bez argumentu wchodzi w odczyt,
// więc uruchomienie ma dokąd sięgnąć po wartość domyślną.
func argumentyTrybuMostu(ctx context.Context, z *zapytania, punktID int64) (map[shared.AccessMode]string, error) {
	polecenie, err := z.przygotuj(ctx, listaArgumentowTrybu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, punktID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać słownictwa trybu punktu %d: %w", punktID, err)
	}
	defer wiersze.Close()

	slownik := map[shared.AccessMode]string{}
	for wiersze.Next() {
		var kolumna, argument string
		if err := wiersze.Scan(&kolumna, &argument); err != nil {
			return nil, fmt.Errorf("dane: nieczytelne słownictwo trybu punktu %d: %w", punktID, err)
		}
		tryb, err := trybDostepuZBazy(kolumna, "argument_trybu_mostu.tryb")
		if err != nil {
			return nil, err
		}
		slownik[tryb] = argument
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt słownictwa trybu punktu %d: %w", punktID, err)
	}
	return slownik, nil
}

// zapiszArgumentyTrybu wymienia słownictwo trybu punktu w transakcji jego zapisu.
func zapiszArgumentyTrybu(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	punktID int64, slownik map[shared.AccessMode]string) error {

	kasowanie, err := z.wTransakcji(ctx, transakcja, usunArgumentyTrybu)
	if err != nil {
		return err
	}
	if _, err := kasowanie.ExecContext(ctx, punktID); err != nil {
		return fmt.Errorf("dane: nie można wyczyścić słownictwa trybu punktu %d: %w", punktID, err)
	}
	wstawianie, err := z.wTransakcji(ctx, transakcja, wstawArgumentTrybu)
	if err != nil {
		return err
	}
	for _, tryb := range trybyDostepu {
		argument, jest := slownik[tryb]
		if !jest || argument == "" {
			continue
		}
		kolumna, err := trybDostepuNaBaze(tryb, "argument_trybu_mostu.tryb")
		if err != nil {
			return err
		}
		if _, err := wstawianie.ExecContext(ctx, punktID, kolumna, argument); err != nil {
			return fmt.Errorf("dane: nie można zapisać słownictwa trybu punktu %d: %w", punktID, err)
		}
	}
	return nil
}

// AdresMostu składa adres uruchomienia mostu z pól połączenia. Punkt, który nie
// jest mostem, adresu nie ma — katalog lokalny leży na urządzeniu, nie za SSH.
func (p PunktDostepu) AdresMostu() string {
	if p.Rodzaj != shared.AccessPointKindMcpBridge || p.Host == "" {
		return ""
	}
	adres := "ssh://"
	if p.Uzytkownik != "" {
		adres += p.Uzytkownik + "@"
	}
	adres += p.Host
	if p.Port > 0 {
		adres += ":" + strconv.Itoa(p.Port)
	}
	return adres
}

// ArgumentTrybu zwraca słowo, jakim ten most nazywa wskazany tryb. Brak wpisu
// daje napis pusty — uruchomienie poda wtedy sam skrypt, a most wejdzie w tryb
// odczytu, czyli w wariant bezpieczniejszy.
func (p PunktDostepu) ArgumentTrybu(tryb shared.AccessMode) string {
	return p.ArgumentyTrybu[tryb]
}
