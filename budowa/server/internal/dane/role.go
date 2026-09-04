// Odpowiedzialność pliku: rola okna komunikacji widziana od strony identyfikatora zewnętrznego okna, droga potrzebna rodzinie komend role.*.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// RolaOkna to rola okna wraz z więzią, którą rola niesie; koordynator pusty znaczy, że okno nie podlega żadnemu koordynatorowi.
type RolaOkna struct {
	Rola shared.WindowRole
	// Koordynator jest identyfikatorem zewnętrznym okna koordynatora; tylko takim posługuje się kontrakt.
	Koordynator string
}

// RepozytoriumRolOkien jest kontraktem zapisu i odczytu roli okna po identyfikatorze zewnętrznym okna komunikacji.
type RepozytoriumRolOkien interface {
	// RolaOkna zwraca rolę okna wraz z więzią, z ErrBrakWiersza, gdy okno żyje w pamięci rdzenia.
	RolaOkna(ctx context.Context, okno string) (RolaOkna, error)
	// ZapiszRoleOkna zapisuje samą rolę; rola inna niż wykonawca zdejmuje więź koordynatora.
	ZapiszRoleOkna(ctx context.Context, okno string, rola shared.WindowRole) error
}

// Okno nie ma kolumny konto_id; własność sprawdza warunekKontaOkna.
var (
	rolaOknaZewnetrznego = `SELECT wykonawca.rola_okna, koordynator.identyfikator_zewnetrzny
	                        FROM okno_komunikacji wykonawca
	                        LEFT JOIN okno_komunikacji koordynator
	                            ON koordynator.id = wykonawca.okno_koordynatora_id
	                        WHERE wykonawca.identyfikator_zewnetrzny = ?
	                          AND ` + warunekKontaOkna("wykonawca")

	zapiszRoleOknaZewnetrznego = `UPDATE okno_komunikacji
	                              SET rola_okna = ?,
	                                  zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                              WHERE identyfikator_zewnetrzny = ?
	                                AND ` + kontoOknaWlasnego

	// Zapis roli innej niż wykonawca zdejmuje więź tym samym poleceniem, a nie drugim: rola i więź muszą się zgadzać w każdej chwili.
	zapiszRoleSamodzielnaOkna = `UPDATE okno_komunikacji
	                             SET rola_okna = ?,
	                                 okno_koordynatora_id = NULL,
	                                 zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             WHERE identyfikator_zewnetrzny = ?
	                               AND ` + kontoOknaWlasnego
)

type repozytoriumRolOkien struct {
	zapytania *zapytania
}

func noweRepozytoriumRolOkien(z *zapytania) *repozytoriumRolOkien {
	return &repozytoriumRolOkien{zapytania: z}
}

// RoleOkien oddaje repozytorium ról okien nad tą samą bazą, co reszta zestawu; jest metodą, nie polem struktury, bo rola okna nie jest osobnym obszarem danych.
func (z *Zestaw) RoleOkien() RepozytoriumRolOkien {
	if z == nil || z.zapytania == nil {
		return nil
	}
	return noweRepozytoriumRolOkien(z.zapytania)
}

// RolaOkna zwraca rolę okna wraz z identyfikatorem zewnętrznym koordynatora, odczytaną po identyfikatorze zewnętrznym okna.
func (r *repozytoriumRolOkien) RolaOkna(ctx context.Context, okno string) (RolaOkna, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, rolaOknaZewnetrznego)
	if err != nil {
		return RolaOkna{}, err
	}
	var rola string
	var koordynator sql.NullString
	err = polecenie.QueryRowContext(ctx, okno, KontoOperatora(ctx)).Scan(&rola, &koordynator)
	if errors.Is(err, sql.ErrNoRows) {
		return RolaOkna{}, fmt.Errorf("%w: okno %q", ErrBrakWiersza, okno)
	}
	if err != nil {
		return RolaOkna{}, fmt.Errorf("dane: nie można odczytać roli okna %q: %w", okno, err)
	}
	wartosc, err := rolaOknaZBazy(rola)
	if err != nil {
		return RolaOkna{}, err
	}
	return RolaOkna{Rola: wartosc, Koordynator: koordynator.String}, nil
}

// ZapiszRoleOkna zapisuje rolę okna wskazanego identyfikatorem zewnętrznym; okno bez wiersza wraca jako ErrBrakWiersza, nie cichy brak skutku.
func (r *repozytoriumRolOkien) ZapiszRoleOkna(ctx context.Context, okno string, rola shared.WindowRole) error {
	kolumna, err := rolaOknaNaBaze(rola)
	if err != nil {
		return err
	}
	zapytanie := zapiszRoleSamodzielnaOkna
	if rola == shared.WindowRoleExecutor {
		zapytanie = zapiszRoleOknaZewnetrznego
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kolumna, okno, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać roli okna %q: %w", okno, err)
	}
	trafione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznany skutek zapisu roli okna %q: %w", okno, err)
	}
	if trafione == 0 {
		return fmt.Errorf("%w: okno %q", ErrBrakWiersza, okno)
	}
	return nil
}
