// Rejestr komponentów własnych Strefy 2 Strony głównej (tabela `komponent`) — trwałość rodziny `component.*`.
package dane

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

// Komponent to wiersz `komponent`; `Kod` odpowiada `Component.id`, a `BytDocelowy` — `Component.targetId`.
type Komponent struct {
	ID             int64
	Kod            string
	Rodzaj         string
	Nazwa          string
	Opis           *string
	BytDocelowy    *string
	Czynny         bool
	Konfiguracja   string
	PoziomZasiegu  *string
	KluczZasiegu   string
	Utworzono      int64
	Zaktualizowano int64
}

// FiltrKomponentow zawęża wykaz komponentów; Rodzaj pusty znaczy „wszystkie”, DolaczWylaczone otwiera na niczynne.
type FiltrKomponentow struct {
	Rodzaj          string
	DolaczWylaczone bool
}

// ZmianaKomponentu niesie pola `component.update`; wskaźnik pusty znaczy „bez zmiany”.
type ZmianaKomponentu struct {
	Nazwa        *string
	Opis         *string
	Czynny       *bool
	Konfiguracja *string
}

// RepozytoriumKomponentow jest kontraktem rejestru komponentów, określającym operacje dostępne na wykazie.
type RepozytoriumKomponentow interface {
	ZalozKomponent(ctx context.Context, komponent Komponent) (Komponent, error)
	Komponent(ctx context.Context, kod string) (Komponent, error)
	Komponenty(ctx context.Context, filtr FiltrKomponentow) ([]Komponent, error)
	ZmienKomponent(ctx context.Context, kod string, zmiana ZmianaKomponentu, teraz int64) (Komponent, error)
	UsunKomponent(ctx context.Context, kod string) (bool, error)
	PrzypiszKomponent(ctx context.Context, kod, poziom, kluczZasiegu string, teraz int64) (Komponent, bool, error)
}

const (
	kolumnyKomponentu = `k.id, k.identyfikator_zewnetrzny, k.rodzaj, k.nazwa, k.opis,
	                     k.byt_docelowy, k.czynny, k.konfiguracja, p.kod, k.klucz_zasiegu,
	                     k.utworzono, k.zaktualizowano`

	zrodloKomponentu = ` FROM komponent k LEFT JOIN poziom_zasiegu p ON p.id = k.poziom_zasiegu_id`

	wstawKomponent = `INSERT INTO komponent
	                  (identyfikator_zewnetrzny, rodzaj, nazwa, opis, byt_docelowy,
	                   czynny, konfiguracja, utworzono, zaktualizowano, konto_id)
	                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzKomponent = `SELECT ` + kolumnyKomponentu + zrodloKomponentu +
		` WHERE k.identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	// Jedno zapytanie na cztery warianty żądania: puste zawężenie rodzaju wyłącza pierwszy warunek.
	listaKomponentow = `SELECT ` + kolumnyKomponentu + zrodloKomponentu +
		` WHERE (? = '' OR k.rodzaj = ?) AND (? = 1 OR k.czynny = 1) AND ` + WarunekKonta + `
		  ORDER BY k.rodzaj, k.nazwa, k.id`

	// NULL w argumencie zostawia kolumnę bez zmiany — „pola pominięte zostają bez zmian" jest własnością zapytania.
	zmienKomponent = `UPDATE komponent SET
	                     nazwa          = COALESCE(?, nazwa),
	                     opis           = CASE WHEN ? = 1 THEN ? ELSE opis END,
	                     czynny         = COALESCE(?, czynny),
	                     konfiguracja   = COALESCE(?, konfiguracja),
	                     zaktualizowano = ?
	                  WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	przypiszKomponent = `UPDATE komponent SET
	                        poziom_zasiegu_id = (SELECT id FROM poziom_zasiegu WHERE kod = ?),
	                        klucz_zasiegu     = ?,
	                        zaktualizowano    = ?
	                     WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	usunKomponent = `DELETE FROM komponent WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

type repozytoriumKomponentow struct {
	zapytania *zapytania
}

func noweRepozytoriumKomponentow(zapytania *zapytania) RepozytoriumKomponentow {
	if zapytania == nil {
		return nil
	}
	return &repozytoriumKomponentow{zapytania: zapytania}
}

// ZalozKomponent wstawia kafel Strefy 2; czas podaje warstwa wyższa w ms epoki, baza nie wstawia własnego „teraz".
func (r *repozytoriumKomponentow) ZalozKomponent(ctx context.Context,
	komponent Komponent) (Komponent, error) {

	if komponent.Kod == "" {
		return Komponent{}, fmt.Errorf("dane: komponent bez identyfikatora")
	}
	if komponent.Nazwa == "" {
		return Komponent{}, fmt.Errorf("dane: komponent %q bez nazwy", komponent.Kod)
	}
	konfiguracja, err := konfiguracjaKomponentu(komponent.Konfiguracja, komponent.Kod)
	if err != nil {
		return Komponent{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKomponent)
	if err != nil {
		return Komponent{}, err
	}
	if _, err := polecenie.ExecContext(ctx, komponent.Kod, komponent.Rodzaj, komponent.Nazwa,
		tekstDoKolumny(komponent.Opis), tekstDoKolumny(komponent.BytDocelowy),
		liczbaLogiczna(komponent.Czynny), konfiguracja,
		komponent.Utworzono, komponent.Zaktualizowano, KontoOperatora(ctx)); err != nil {

		return Komponent{}, fmt.Errorf("dane: nie można założyć komponentu %q: %w", komponent.Kod, err)
	}
	return r.Komponent(ctx, komponent.Kod)
}

// Komponent zwraca kafel po kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumKomponentow) Komponent(ctx context.Context, kod string) (Komponent, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKomponent)
	if err != nil {
		return Komponent{}, err
	}
	komponent, err := odczytajKomponent(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return Komponent{}, fmt.Errorf("dane: komponent %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	if err != nil {
		return Komponent{}, fmt.Errorf("dane: nie można odczytać komponentu %q: %w", kod, err)
	}
	return komponent, nil
}

func (r *repozytoriumKomponentow) Komponenty(ctx context.Context,
	filtr FiltrKomponentow) ([]Komponent, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaKomponentow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, filtr.Rodzaj, filtr.Rodzaj,
		liczbaLogiczna(filtr.DolaczWylaczone), KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu komponentów: %w", err)
	}
	defer wiersze.Close()

	komponenty := make([]Komponent, 0, 16)
	for wiersze.Next() {
		komponent, err := odczytajKomponent(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: uszkodzony wiersz komponentu: %w", err)
		}
		komponenty = append(komponenty, komponent)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu komponentów: %w", err)
	}
	return komponenty, nil
}

// ZmienKomponent zmienia tylko pola wskazane w żądaniu; opis ma osobny przełącznik (jego zmianą bywa wyczyszczenie).
func (r *repozytoriumKomponentow) ZmienKomponent(ctx context.Context, kod string,
	zmiana ZmianaKomponentu, teraz int64) (Komponent, error) {

	if _, err := r.Komponent(ctx, kod); err != nil {
		return Komponent{}, err
	}
	konfiguracja, err := zmienionaKonfiguracja(zmiana.Konfiguracja, kod)
	if err != nil {
		return Komponent{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienKomponent)
	if err != nil {
		return Komponent{}, err
	}
	var czynny any
	if zmiana.Czynny != nil {
		czynny = liczbaLogiczna(*zmiana.Czynny)
	}
	zmianaOpisu := liczbaLogiczna(zmiana.Opis != nil)
	if _, err := polecenie.ExecContext(ctx, tekstDoKolumny(zmiana.Nazwa),
		zmianaOpisu, tekstDoKolumny(zmiana.Opis), czynny, konfiguracja,
		teraz, kod, KontoOperatora(ctx)); err != nil {

		return Komponent{}, fmt.Errorf("dane: nie można zmienić komponentu %q: %w", kod, err)
	}
	return r.Komponent(ctx, kod)
}

// UsunKomponent zdejmuje kafel Strefy 2; bytu magazynu modułowego nie tyka.
func (r *repozytoriumKomponentow) UsunKomponent(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunKomponent)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć komponentu %q: %w", kod, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można potwierdzić usunięcia komponentu %q: %w", kod, err)
	}
	return usuniete > 0, nil
}

// PrzypiszKomponent zapisuje parę (poziom zasięgu, klucz); powtórzenie niczego nie zmienia i wraca jako `false`.
func (r *repozytoriumKomponentow) PrzypiszKomponent(ctx context.Context,
	kod, poziom, kluczZasiegu string, teraz int64) (Komponent, bool, error) {

	zastany, err := r.Komponent(ctx, kod)
	if err != nil {
		return Komponent{}, false, err
	}
	if zastany.PoziomZasiegu != nil && *zastany.PoziomZasiegu == poziom &&
		zastany.KluczZasiegu == kluczZasiegu {

		return zastany, false, nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszKomponent)
	if err != nil {
		return Komponent{}, false, err
	}
	if _, err := polecenie.ExecContext(ctx, poziom, kluczZasiegu,
		teraz, kod, KontoOperatora(ctx)); err != nil {

		return Komponent{}, false, fmt.Errorf("dane: nie można przypisać komponentu %q: %w", kod, err)
	}
	przypisany, err := r.Komponent(ctx, kod)
	if err != nil {
		return Komponent{}, false, err
	}
	if przypisany.PoziomZasiegu == nil {
		return Komponent{}, false,
			fmt.Errorf("dane: poziom zasięgu %q nie istnieje w słowniku", poziom)
	}
	return przypisany, true, nil
}

func konfiguracjaKomponentu(tresc, kod string) (string, error) {
	if tresc == "" {
		return "{}", nil
	}
	if !json.Valid([]byte(tresc)) {
		return "", fmt.Errorf("dane: konfiguracja komponentu %q nie jest poprawnym JSON-em", kod)
	}
	return tresc, nil
}

func zmienionaKonfiguracja(tresc *string, kod string) (any, error) {
	if tresc == nil {
		return nil, nil
	}
	sprawdzona, err := konfiguracjaKomponentu(*tresc, kod)
	if err != nil {
		return nil, err
	}
	return sprawdzona, nil
}

func odczytajKomponent(wiersz skaner) (Komponent, error) {
	var komponent Komponent
	var opis, bytDocelowy, poziom sql.NullString
	var czynny int
	if err := wiersz.Scan(&komponent.ID, &komponent.Kod, &komponent.Rodzaj, &komponent.Nazwa,
		&opis, &bytDocelowy, &czynny, &komponent.Konfiguracja, &poziom, &komponent.KluczZasiegu,
		&komponent.Utworzono, &komponent.Zaktualizowano); err != nil {

		return Komponent{}, err
	}
	komponent.Opis = tekstZKolumny(opis)
	komponent.BytDocelowy = tekstZKolumny(bytDocelowy)
	komponent.PoziomZasiegu = tekstZKolumny(poziom)
	komponent.Czynny = czynny == 1
	return komponent, nil
}
