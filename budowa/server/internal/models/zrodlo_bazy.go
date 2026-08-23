package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// zapytanieOKanaly czyta rejestr kanałów wprost z tabeli kanal_modelu.
// Kolejność wiersza rozstrzyga kolumna kolejnosc, a przy jej równości kod —
// wykaz kanałów jest wtedy powtarzalny między uruchomieniami.
const zapytanieOKanaly = `
SELECT id, kod, nazwa, dostawca, identyfikator_modelu, rodzaj_kanalu,
       COALESCE(konto_id, 0), COALESCE(poswiadczenie_odwolanie, ''),
       COALESCE(parametry_json, '{}'), multimodalny, aktywny, kolejnosc, utworzono
  FROM kanal_modelu
 ORDER BY kolejnosc, kod`

// odpytujacy jest tą częścią *sql.DB, której źródło rzeczywiście używa.
// Dzięki temu rejestr da się zasilić także z transakcji.
type odpytujacy interface {
	QueryContext(ctx context.Context, zapytanie string, argumenty ...any) (*sql.Rows, error)
}

// zrodloBazy czyta wiersze rejestru kanałów z bazy produktu.
type zrodloBazy struct {
	db odpytujacy
}

// NoweZrodloBazy składa źródło nad otwartą bazą.
func NoweZrodloBazy(db odpytujacy) *zrodloBazy {
	return &zrodloBazy{db: db}
}

// Definicje odczytuje komplet wierszy tabeli kanal_modelu — także nieczynne,
// bo o pominięciu nieczynnego kanału rozstrzyga rejestr, nie zapytanie.
func (z *zrodloBazy) Definicje(ctx context.Context) ([]Definicja, error) {
	if z == nil || z.db == nil {
		return nil, fmt.Errorf("models: źródło kanałów bez bazy")
	}
	wiersze, err := z.db.QueryContext(ctx, zapytanieOKanaly)
	if err != nil {
		return nil, fmt.Errorf("models: odczyt rejestru kanałów: %w", err)
	}
	defer wiersze.Close()

	definicje := make([]Definicja, 0, 8)
	for wiersze.Next() {
		d, err := odczytajWiersz(wiersze)
		if err != nil {
			return nil, err
		}
		definicje = append(definicje, d)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("models: odczyt rejestru kanałów: %w", err)
	}
	return definicje, nil
}

// odczytajWiersz przenosi jeden rekord tabeli na definicję kanału.
func odczytajWiersz(wiersze *sql.Rows) (Definicja, error) {
	var (
		d             Definicja
		parametryJSON string
		multimodalny  int
		aktywny       int
	)
	err := wiersze.Scan(&d.Id, &d.Kod, &d.Nazwa, &d.Dostawca, &d.Model, &d.Rodzaj,
		&d.KontoId, &d.PoswiadczenieOdwolanie, &parametryJSON,
		&multimodalny, &aktywny, &d.Kolejnosc, &d.Utworzono)
	if err != nil {
		return Definicja{}, fmt.Errorf("models: wiersz rejestru kanałów: %w", err)
	}
	d.Multimodalny = multimodalny != 0
	d.Aktywny = aktywny != 0
	d.Parametry = odczytajParametry(parametryJSON)
	return d, nil
}

// odczytajParametry rozpakowuje kolumnę parametry_json. Treść nieczytelna daje
// pusty zestaw parametrów zamiast błędu całego odczytu: jeden zepsuty wiersz
// nie może wygasić rejestru kanałów.
func odczytajParametry(surowe string) map[string]any {
	parametry := map[string]any{}
	if surowe == "" {
		return parametry
	}
	if err := json.Unmarshal([]byte(surowe), &parametry); err != nil {
		return map[string]any{}
	}
	return parametry
}
