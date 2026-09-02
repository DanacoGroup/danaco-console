// Odpowiedzialność pliku: listy korzeni obszaru dostępów — korzenie punktu
// (`korzen_punktu_dostepu`) i korzenie nadania (`korzen_nadania`).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ErrPozaKorzeniami odróżnia warstwa wyższa przez errors.Is, nie przez treść komunikatu.
var ErrPozaKorzeniami = errors.New("dane: ścieżka poza korzeniami punktu dostępu")

// Tabele korzeni konta nie niosą (krok 484 dodał `konto_id` punktowi, nie jego korzeniom):
// granica dochodzi przez `punkt_dostepu.konto_id`, a dla korzenia nadania przez nadanie i jego okno.
var (
	kontoKorzeniaPunktu = `EXISTS (SELECT 1 FROM punkt_dostepu
	                                WHERE punkt_dostepu.id = korzen_punktu_dostepu.punkt_dostepu_id
	                                  AND ` + WarunekKonta + `)`

	kontoKorzeniaNadania = `EXISTS (SELECT 1 FROM nadanie_dostepu
	                                 WHERE nadanie_dostepu.id = korzen_nadania.nadanie_dostepu_id
	                                   AND ` + kontoNadania + `)`

	usunKorzeniePunktu = `DELETE FROM korzen_punktu_dostepu
	                      WHERE punkt_dostepu_id = ? AND ` + kontoKorzeniaPunktu

	wstawKorzenPunktu = `INSERT INTO korzen_punktu_dostepu (punkt_dostepu_id, sciezka, kolejnosc)
	                     SELECT ?, ?, ?
	                      WHERE EXISTS (SELECT 1 FROM punkt_dostepu WHERE id = ? AND ` + WarunekKonta + `)`

	listaKorzeniPunktu = `SELECT sciezka FROM korzen_punktu_dostepu
	                      WHERE punkt_dostepu_id = ? AND ` + kontoKorzeniaPunktu + `
	                      ORDER BY kolejnosc, id`

	usunKorzenieNadania = `DELETE FROM korzen_nadania
	                       WHERE nadanie_dostepu_id = ? AND ` + kontoKorzeniaNadania

	wstawKorzenNadania = `INSERT INTO korzen_nadania (nadanie_dostepu_id, sciezka, kolejnosc)
	                      SELECT ?, ?, ?
	                       WHERE EXISTS (SELECT 1 FROM nadanie_dostepu
	                                      WHERE id = ? AND ` + kontoNadania + `)`

	listaKorzeniNadania = `SELECT sciezka FROM korzen_nadania
	                       WHERE nadanie_dostepu_id = ? AND ` + kontoKorzeniaNadania + `
	                       ORDER BY kolejnosc, id`
)

func zapiszKorzenie(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	kasowanieSQL, wstawianieSQL string, wlascicielID int64, korzenie []string, opis string) error {

	konto := KontoOperatora(ctx)
	kasowanie, err := z.wTransakcji(ctx, transakcja, kasowanieSQL)
	if err != nil {
		return err
	}
	if _, err := kasowanie.ExecContext(ctx, wlascicielID, konto); err != nil {
		return fmt.Errorf("dane: nie można wyczyścić korzeni %s %d: %w", opis, wlascicielID, err)
	}
	wstawianie, err := z.wTransakcji(ctx, transakcja, wstawianieSQL)
	if err != nil {
		return err
	}
	kolejnosc := 0
	for _, sciezka := range uporzadkujKorzenie(korzenie) {
		wynik, err := wstawianie.ExecContext(ctx, wlascicielID, sciezka, kolejnosc, wlascicielID, konto)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać korzenia %q %s %d: %w",
				sciezka, opis, wlascicielID, err)
		}
		if err := sprawdzTrafienieZapisu(wynik, "korzeń "+opis, fmt.Sprintf("%d", wlascicielID)); err != nil {
			return err
		}
		kolejnosc++
	}
	return nil
}

func wczytajKorzenie(ctx context.Context, z *zapytania, zapytanie string,
	wlascicielID int64, opis string, dalszeArgumenty ...any) ([]string, error) {

	polecenie, err := z.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, append([]any{wlascicielID}, dalszeArgumenty...)...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać korzeni %s %d: %w", opis, wlascicielID, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var sciezka string
		if err := wiersze.Scan(&sciezka); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny korzeń %s %d: %w", opis, wlascicielID, err)
		}
		lista = append(lista, sciezka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt korzeni %s %d: %w", opis, wlascicielID, err)
	}
	return lista, nil
}

func uporzadkujKorzenie(korzenie []string) []string {
	wynik := make([]string, 0, len(korzenie))
	widziane := map[string]struct{}{}
	for _, sciezka := range korzenie {
		sciezka = strings.TrimSpace(sciezka)
		if sciezka == "" {
			continue
		}
		if _, jest := widziane[sciezka]; jest {
			continue
		}
		widziane[sciezka] = struct{}{}
		wynik = append(wynik, sciezka)
	}
	return wynik
}

// Komunikat niesie samą ścieżkę odrzuconą: treść wychodzi kontraktem jako `validation_failed`.
func sprawdzZawezenieKorzeni(korzeniePunktu, korzenieNadania []string) error {
	if len(korzeniePunktu) == 0 {
		return nil
	}
	for _, sciezka := range uporzadkujKorzenie(korzenieNadania) {
		if !wKtorymkolwiekKorzeniu(korzeniePunktu, sciezka) {
			return fmt.Errorf("dane: korzeń nadania %q leży poza obszarem punktu dostępu: %w",
				sciezka, ErrPozaKorzeniami)
		}
	}
	return nil
}

func wKtorymkolwiekKorzeniu(korzenie []string, sciezka string) bool {
	for _, korzen := range korzenie {
		if wKorzeniu(korzen, sciezka) {
			return true
		}
	}
	return false
}

// Wielkość liter zostaje znacząca: maszyny mostu pracują na systemie plików, który ją rozróżnia.
func wKorzeniu(korzen, sciezka string) bool {
	k := normalizujKorzen(korzen)
	s := normalizujKorzen(sciezka)
	if k == "" || s == "" {
		return false
	}
	return s == k || strings.HasPrefix(s, k+"/")
}

func normalizujKorzen(sciezka string) string {
	ujednolicona := strings.ReplaceAll(strings.TrimSpace(sciezka), `\`, "/")
	obcieta := strings.TrimRight(ujednolicona, "/")
	if obcieta == "" {
		return ujednolicona
	}
	return obcieta
}
