// Odpowiedzialność pliku: dostęp do obszaru konfiguracji (tabele `ustawienie`
// i `poziom_zasiegu`). Repozytorium czyta i zapisuje wartości na wskazanym poziomie zasięgu.
package dane

import (
	"context"
	"database/sql"
	"fmt"

	"danacoconsole/shared"
)

// Ustawienie to wiersz tabeli `ustawienie` opisany poziomem zasięgu kontraktu oraz osią rozstrzygania,
// dla czego wartość obowiązuje.
type Ustawienie struct {
	Poziom         shared.ConfigScope
	KluczZasiegu   string
	Os             shared.ConfigAxis
	KluczOsi       string
	Klucz          string
	Wartosc        *string
	RodzajWartosci string
	Zaktualizowano string
}

// RepozytoriumKonfiguracji jest kontraktem obszaru konfiguracji, opisującym metody dotyczące osi platformy.
type RepozytoriumKonfiguracji interface {
	Ustaw(ctx context.Context, ustawienie Ustawienie) error
	Odczytaj(ctx context.Context, poziom shared.ConfigScope, kluczZasiegu, klucz string) (Ustawienie, bool, error)
	ListaPoziomu(ctx context.Context, poziom shared.ConfigScope, kluczZasiegu string) ([]Ustawienie, error)
	Usun(ctx context.Context, poziom shared.ConfigScope, kluczZasiegu, klucz string) error
}

const domyslnyRodzajWartosci = "tekst"

const (
	kolumnyUstawienia = `p.kod, u.klucz_zasiegu, u.os, u.klucz_osi, u.klucz, u.wartosc,
	                     u.rodzaj_wartosci, u.zaktualizowano`

	zrodloUstawienia = ` FROM ustawienie u JOIN poziom_zasiegu p ON p.id = u.poziom_zasiegu_id`

	// Więz UNIQUE adresu obejmuje całą tabelę, więc konflikt trafia i w wiersz
	// cudzego konta; WarunekKonta przy DO UPDATE odcina tam zapis.
	zapiszUstawienie = `INSERT INTO ustawienie
	                    (poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz, wartosc,
	                     rodzaj_wartosci, konto_id)
	                    VALUES ((SELECT id FROM poziom_zasiegu WHERE kod = ?), ?, ?, ?, ?, ?, ?,
	                            ` + WskazanieKonta + `)
	                    ON CONFLICT(poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz,
	                                COALESCE(konto_id, 0)) DO UPDATE SET
	                        wartosc = excluded.wartosc,
	                        rodzaj_wartosci = excluded.rodzaj_wartosci,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                    WHERE ` + WarunekKonta

	pobierzUstawienie = `SELECT ` + kolumnyUstawienia + zrodloUstawienia +
		` WHERE p.kod = ? AND u.klucz_zasiegu = ? AND u.os = ? AND u.klucz_osi = ? AND u.klucz = ?
		  AND ` + WarunekKonta

	listaUstawienPoziomu = `SELECT ` + kolumnyUstawienia + zrodloUstawienia +
		` WHERE p.kod = ? AND u.klucz_zasiegu = ? AND u.os = ? AND u.klucz_osi = ?
		  AND ` + WarunekKonta + `
		  ORDER BY u.klucz`

	usunUstawienie = `DELETE FROM ustawienie
	                  WHERE poziom_zasiegu_id = (SELECT id FROM poziom_zasiegu WHERE kod = ?)
	                    AND klucz_zasiegu = ? AND os = ? AND klucz_osi = ? AND klucz = ?
	                    AND ` + WarunekKonta
)

type repozytoriumKonfiguracji struct {
	zapytania *zapytania
}

func noweRepozytoriumKonfiguracji(z *zapytania) *repozytoriumKonfiguracji {
	return &repozytoriumKonfiguracji{zapytania: z}
}

// Ustaw zapisuje wartość pod adresem złożonym z poziomu zasięgu i osi; wiersz
// istniejący nadpisuje. Oś pusta znaczy `platform`.
func (r *repozytoriumKonfiguracji) Ustaw(ctx context.Context, ustawienie Ustawienie) error {
	poziom, err := poziomZasieguNaBaze(ustawienie.Poziom)
	if err != nil {
		return err
	}
	os, err := osZasieguNaBaze(ustawienie.Os)
	if err != nil {
		return err
	}
	rodzaj := ustawienie.RodzajWartosci
	if rodzaj == "" {
		rodzaj = domyslnyRodzajWartosci
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszUstawienie)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, poziom, ustawienie.KluczZasiegu, os, ustawienie.KluczOsi,
		ustawienie.Klucz, tekstDoKolumny(ustawienie.Wartosc), rodzaj, KontoOperatora(ctx),
		KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać ustawienia %q na poziomie %q osi %q: %w",
			ustawienie.Klucz, poziom, os, err)
	}
	return sprawdzTrafienieZapisu(wynik, "ustawienie", ustawienie.Klucz)
}

// Odczytaj zwraca ustawienie osi platformy z jednego poziomu. Drugi wynik mówi,
// czy wartość w ogóle ustawiono — brak wiersza nie jest błędem.
func (r *repozytoriumKonfiguracji) Odczytaj(ctx context.Context, poziom shared.ConfigScope,
	kluczZasiegu, klucz string) (Ustawienie, bool, error) {

	return r.OdczytajOsi(ctx, poziom, kluczZasiegu, shared.ConfigAxisPlatform, "", klucz)
}

// ListaPoziomu zwraca ustawienia osi platformy dla jednego bytu danego poziomu zasięgu konfiguracji systemu.
func (r *repozytoriumKonfiguracji) ListaPoziomu(ctx context.Context, poziom shared.ConfigScope,
	kluczZasiegu string) ([]Ustawienie, error) {

	return r.ListaOsi(ctx, poziom, kluczZasiegu, shared.ConfigAxisPlatform, "")
}

// Usun kasuje ustawienie osi platformy; wartość wraca wtedy do wartości domyślnej tej konfiguracji systemu.
func (r *repozytoriumKonfiguracji) Usun(ctx context.Context, poziom shared.ConfigScope,
	kluczZasiegu, klucz string) error {

	return r.UsunOsi(ctx, poziom, kluczZasiegu, shared.ConfigAxisPlatform, "", klucz)
}

// odczytajUstawienie składa pełną strukturę ustawienia z jednego wiersza wyniku zapytania do bazy danych.
func odczytajUstawienie(wiersz skaner) (Ustawienie, error) {
	var ustawienie Ustawienie
	var kod, os string
	var wartosc sql.NullString
	err := wiersz.Scan(&kod, &ustawienie.KluczZasiegu, &os, &ustawienie.KluczOsi,
		&ustawienie.Klucz, &wartosc, &ustawienie.RodzajWartosci, &ustawienie.Zaktualizowano)
	if err != nil {
		return Ustawienie{}, err
	}
	ustawienie.Wartosc = tekstZKolumny(wartosc)
	if ustawienie.Poziom, err = poziomZasieguZBazy(kod); err != nil {
		return Ustawienie{}, err
	}
	if ustawienie.Os, err = osZasieguZBazy(os); err != nil {
		return Ustawienie{}, err
	}
	return ustawienie, nil
}
