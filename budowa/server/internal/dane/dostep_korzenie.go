// Odpowiedzialność pliku: listy korzeni obszaru dostępów — korzenie punktu
// (`korzen_punktu_dostepu`) i korzenie nadania (`korzen_nadania`). Korzeń jest
// listą, nie pojedynczym polem, tak samo jak katalog roboczy okna;
// odpowiada zmiennej DANACO_MOST_KORZENIE mostu MCP, poza którą most nie wyjdzie.
//
// Nadanie może korzenie punktu wyłącznie ZAWĘZIĆ. Korzeń nadania spoza obszaru
// punktu byłby obietnicą dostępu, którego most i tak nie da — sprawdzenie leży
// tu, żeby nie powstał wiersz wprowadzający w błąd.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ErrPozaKorzeniami oznacza próbę nadania dostępu do ścieżki spoza obszaru
// punktu. Warstwa wyższa odróżnia ten przypadek przez errors.Is, nie przez treść
// komunikatu — to odmowa merytoryczna, nie awaria zapisu.
var ErrPozaKorzeniami = errors.New("dane: ścieżka poza korzeniami punktu dostępu")

const (
	usunKorzeniePunktu = `DELETE FROM korzen_punktu_dostepu WHERE punkt_dostepu_id = ?`

	wstawKorzenPunktu = `INSERT INTO korzen_punktu_dostepu (punkt_dostepu_id, sciezka, kolejnosc)
	                     VALUES (?, ?, ?)`

	listaKorzeniPunktu = `SELECT sciezka FROM korzen_punktu_dostepu
	                      WHERE punkt_dostepu_id = ? ORDER BY kolejnosc, id`

	usunKorzenieNadania = `DELETE FROM korzen_nadania WHERE nadanie_dostepu_id = ?`

	wstawKorzenNadania = `INSERT INTO korzen_nadania (nadanie_dostepu_id, sciezka, kolejnosc)
	                      VALUES (?, ?, ?)`

	listaKorzeniNadania = `SELECT sciezka FROM korzen_nadania
	                       WHERE nadanie_dostepu_id = ? ORDER BY kolejnosc, id`
)

// zapiszKorzenie wymienia listę korzeni w tabeli podrzędnej. Wywoływane wyłącznie
// wewnątrz transakcji zapisu właściciela — inaczej właściciel i jego korzenie
// mogłyby się rozjechać.
func zapiszKorzenie(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	kasowanieSQL, wstawianieSQL string, wlascicielID int64, korzenie []string, opis string) error {

	kasowanie, err := z.wTransakcji(ctx, transakcja, kasowanieSQL)
	if err != nil {
		return err
	}
	if _, err := kasowanie.ExecContext(ctx, wlascicielID); err != nil {
		return fmt.Errorf("dane: nie można wyczyścić korzeni %s %d: %w", opis, wlascicielID, err)
	}
	wstawianie, err := z.wTransakcji(ctx, transakcja, wstawianieSQL)
	if err != nil {
		return err
	}
	kolejnosc := 0
	for _, sciezka := range uporzadkujKorzenie(korzenie) {
		if _, err := wstawianie.ExecContext(ctx, wlascicielID, sciezka, kolejnosc); err != nil {
			return fmt.Errorf("dane: nie można zapisać korzenia %q %s %d: %w",
				sciezka, opis, wlascicielID, err)
		}
		kolejnosc++
	}
	return nil
}

// wczytajKorzenie zwraca listę korzeni jednego właściciela w zapisanej kolejności.
func wczytajKorzenie(ctx context.Context, z *zapytania, zapytanie string,
	wlascicielID int64, opis string) ([]string, error) {

	polecenie, err := z.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, wlascicielID)
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

// uporzadkujKorzenie odrzuca wpisy puste i powtórzone, zachowując kolejność
// podaną przez Operatora — kolejność korzeni ma znaczenie przy przekazaniu ich
// mostowi.
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

// sprawdzZawezenieKorzeni pilnuje, żeby korzenie nadania mieściły się w obszarze
// punktu. Lista pusta znaczy „komplet korzeni punktu" i jest poprawna zawsze.
// Punkt bez własnych korzeni nie ogranicza niczego — tak samo jak
// most z pustą zmienną DANACO_MOST_KORZENIE.
func sprawdzZawezenieKorzeni(korzeniePunktu, korzenieNadania []string) error {
	if len(korzeniePunktu) == 0 {
		return nil
	}
	for _, sciezka := range uporzadkujKorzenie(korzenieNadania) {
		if !wKtorymkolwiekKorzeniu(korzeniePunktu, sciezka) {
			return fmt.Errorf("dane: korzeń nadania %q leży poza obszarem punktu dostępu %v: %w",
				sciezka, korzeniePunktu, ErrPozaKorzeniami)
		}
	}
	return nil
}

// wKtorymkolwiekKorzeniu rozstrzyga, czy ścieżka mieści się w którymkolwiek
// z korzeni punktu.
func wKtorymkolwiekKorzeniu(korzenie []string, sciezka string) bool {
	for _, korzen := range korzenie {
		if wKorzeniu(korzen, sciezka) {
			return true
		}
	}
	return false
}

// wKorzeniu porównuje ścieżkę z korzeniem po ujednoliceniu separatora. Wielkość
// liter zostaje znacząca: maszyny mostu pracują na systemie plików rozróżniającym
// wielkość liter, a zrównanie liter przepuściłoby ścieżkę spoza obszaru.
func wKorzeniu(korzen, sciezka string) bool {
	k := normalizujKorzen(korzen)
	s := normalizujKorzen(sciezka)
	if k == "" || s == "" {
		return false
	}
	return s == k || strings.HasPrefix(s, k+"/")
}

// normalizujKorzen ujednolica separator i obcina separator końcowy, żeby
// `/opt/danaco` i `/opt/danaco/` znaczyły to samo.
func normalizujKorzen(sciezka string) string {
	ujednolicona := strings.ReplaceAll(strings.TrimSpace(sciezka), `\`, "/")
	obcieta := strings.TrimRight(ujednolicona, "/")
	if obcieta == "" {
		return ujednolicona
	}
	return obcieta
}
