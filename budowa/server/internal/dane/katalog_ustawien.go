// Odpowiedzialność pliku: dostęp do katalogu ustawień (tabele
// `kategoria_ustawien`, `definicja_ustawienia`, `opcja_ustawienia`,
// `definicja_ustawienia_zasieg`, `definicja_ustawienia_os`) — część kategorii.
//
// Katalog jest sterowany danymi: nowa pozycja okna konfiguracji to nowy wiersz,
// nie nowa gałąź w kodzie. Wzorcem jest rejestr kanałów
// modelu i katalog akcji — repozytorium wyłącznie czyta wiersze, a rozstrzyganie
// wartości należy do pakietu `internal/konfig`.
//
// Brak wiersza w katalogu nie jest awarią: rezolwer schodzi wtedy na rejestr
// wbudowany rdzenia i pracuje dalej.
package dane

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// RepozytoriumKatalogUstawien jest kontraktem katalogu ustawień.
// Zapisu tu nie ma — wiersze wnosi zaczyn migracji, a rozszerzenie katalogu
// jest zmianą danych, nie czynnością kontraktu.
type RepozytoriumKatalogUstawien interface {
	Kategorie(ctx context.Context, tylkoAktywne bool) ([]shared.SettingCategory, error)
	Definicje(ctx context.Context, tylkoAktywne bool) ([]shared.SettingDefinition, error)
}

const listaKategoriiUstawien = `SELECT k.kod, k.nazwa, k.opis, k.ikona, r.kod,
                                       k.kolejnosc, k.aktywna
                                  FROM kategoria_ustawien k
                                  LEFT JOIN kategoria_ustawien r ON r.id = k.kategoria_nadrzedna_id
                                 WHERE (? = 0 OR k.aktywna = 1)
                                 ORDER BY k.kolejnosc, k.kod`

type repozytoriumKatalogUstawien struct {
	zapytania *zapytania
}

func noweRepozytoriumKatalogUstawien(z *zapytania) *repozytoriumKatalogUstawien {
	return &repozytoriumKatalogUstawien{zapytania: z}
}

// Kategorie zwraca kategorie okna konfiguracji uporządkowane kolejnością wiersza.
func (r *repozytoriumKatalogUstawien) Kategorie(ctx context.Context,
	tylkoAktywne bool) ([]shared.SettingCategory, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaKategoriiUstawien)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, liczbaLogiczna(tylkoAktywne))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kategorii ustawień: %w", err)
	}
	defer wiersze.Close()

	lista := []shared.SettingCategory{}
	for wiersze.Next() {
		kategoria, err := odczytajKategorieUstawien(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, kategoria)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kategorii ustawień: %w", err)
	}
	return lista, nil
}

// odczytajKategorieUstawien składa kategorię kontraktu z jednego wiersza.
func odczytajKategorieUstawien(wiersz skaner) (shared.SettingCategory, error) {
	var kategoria shared.SettingCategory
	var opis, ikona string
	var rodzic sql.NullString
	var aktywna int
	err := wiersz.Scan(&kategoria.Id, &kategoria.Name, &opis, &ikona, &rodzic,
		&kategoria.Order, &aktywna)
	if err != nil {
		return shared.SettingCategory{}, fmt.Errorf("dane: nieczytelny wiersz kategorii ustawień: %w", err)
	}
	kategoria.Description = tekstNiepusty(opis)
	kategoria.Icon = tekstNiepusty(ikona)
	kategoria.ParentId = tekstZKolumny(rodzic)
	kategoria.Enabled = aktywna == 1
	return kategoria, nil
}

// tekstNiepusty zwraca wskaźnik dla treści niepustej, a dla pustej nil —
// pole opcjonalne kontraktu nie niesie wtedy pustego napisu.
func tekstNiepusty(wartosc string) *string {
	if strings.TrimSpace(wartosc) == "" {
		return nil
	}
	kopia := wartosc
	return &kopia
}

// wartoscDomyslnaJSON koduje wartość domyślną katalogu na pole `defaultValue`
// koperty kontraktu. Rodzaje liczbowe, logiczne i złożone idą surowo, jeżeli są
// poprawnym JSON-em; wszystko pozostałe idzie napisem. Funkcja nigdy nie
// zawodzi — wartość nieczytelna trafia do kontraktu jako napis, nie jako błąd.
func wartoscDomyslnaJSON(wartosc string, rodzaj shared.SettingValueType) json.RawMessage {
	switch rodzaj {
	case shared.SettingValueTypeInt, shared.SettingValueTypeFloat, shared.SettingValueTypeBool,
		shared.SettingValueTypeJson, shared.SettingValueTypeEnumList, shared.SettingValueTypePathList:
		przyciete := strings.TrimSpace(wartosc)
		if przyciete != "" && json.Valid([]byte(przyciete)) {
			return json.RawMessage(przyciete)
		}
	}
	surowa, err := json.Marshal(wartosc)
	if err != nil {
		return json.RawMessage(`""`)
	}
	return surowa
}

// zasiegiKontraktu przekłada kody poziomów na wartości wyliczenia kontraktu.
// Kod nierozpoznany jest pomijany, nie przerywa odczytu katalogu.
func zasiegiKontraktu(kody []string) []shared.ConfigScope {
	poziomy := make([]shared.ConfigScope, 0, len(kody))
	for _, kod := range kody {
		if poziom, err := poziomZasieguZBazy(kod); err == nil {
			poziomy = append(poziomy, poziom)
		}
	}
	return poziomy
}

// osieKontraktu przekłada kody osi na wartości wyliczenia kontraktu.
func osieKontraktu(kody []string) []shared.ConfigAxis {
	osie := make([]shared.ConfigAxis, 0, len(kody))
	for _, kod := range kody {
		if os, err := osZasieguZBazy(kod); err == nil {
			osie = append(osie, os)
		}
	}
	return osie
}
