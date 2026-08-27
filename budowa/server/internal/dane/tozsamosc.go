// Plik prowadzi katalog kategorii zasad i tożsamości modelu w tabeli
// kategoria_tozsamosci, wyłącznie do odczytu, sterowany danymi migracji
// zamiast rozgałęzieniami kodu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// KategoriaTozsamosci to wiersz katalogu kategorii. Odpowiada strukturze
// shared.IdentityCategory; kolumny słownikowe niosą wartości kontraktu wprost,
// bo kontrakt nie daje dla nich słownika przekładu bazy.
type KategoriaTozsamosci struct {
	ID           int64
	Kod          string
	Nazwa        string
	Opis         string
	Warstwa      shared.IdentityLayer
	Kolejnosc    int
	Obowiazkowa  bool
	TrybDomyslny shared.IdentityMode
	Aktywna      bool
}

// DokumentTozsamosci to treść jednej kategorii zapisana dla jednej osi.
// Odpowiada strukturze shared.IdentityDocument.
type DokumentTozsamosci struct {
	ID             int64
	KategoriaID    int64
	KodKategorii   string
	Os             shared.ConfigAxis
	OsByt          string
	Tryb           shared.IdentityMode
	Tresc          string
	OdciskTresci   string
	Aktywny        bool
	Zaktualizowano string
}

// FiltrTozsamosci zawęża odczyt treści. Pole puste znaczy „bez zawężenia”, nie
// „brak wyniku” — odczyt bez filtra zwraca komplet zapisów.
type FiltrTozsamosci struct {
	KodKategorii string
	Os           shared.ConfigAxis
	OsByt        string
	TylkoAktywne bool
}

// RepozytoriumTozsamosci jest kontraktem obszaru tożsamości modelu: katalog
// kategorii jest wyłącznie do odczytu, a treść kategorii repozytorium
// zapisuje i usuwa.
type RepozytoriumTozsamosci interface {
	Kategorie(ctx context.Context, tylkoAktywne bool) ([]KategoriaTozsamosci, error)
	KategoriaPoKodzie(ctx context.Context, kod string) (KategoriaTozsamosci, error)
	Dokumenty(ctx context.Context, filtr FiltrTozsamosci) ([]DokumentTozsamosci, error)
	ZapiszDokument(ctx context.Context, dokument DokumentTozsamosci) (DokumentTozsamosci, error)
	UsunDokument(ctx context.Context, id int64) (bool, error)
}

const (
	kolumnyKategoriiTozsamosci = `k.id, k.kod, k.nazwa, k.opis, k.warstwa, k.kolejnosc,
	                              k.obowiazkowa, k.tryb_domyslny, k.aktywna`

	listaKategoriiTozsamosci = `SELECT ` + kolumnyKategoriiTozsamosci + `
	                            FROM kategoria_tozsamosci k
	                            WHERE (? = 0 OR k.aktywna = 1)
	                            ORDER BY k.warstwa, k.kolejnosc, k.kod`

	pobierzKategorieTozsamosci = `SELECT ` + kolumnyKategoriiTozsamosci + `
	                              FROM kategoria_tozsamosci k WHERE k.kod = ?`
)

type repozytoriumTozsamosci struct {
	zapytania *zapytania
}

// noweRepozytoriumTozsamosci zakłada repozytorium tożsamości modelu nad pamięcią
// przygotowanych zapytań zestawu.
func noweRepozytoriumTozsamosci(z *zapytania) *repozytoriumTozsamosci {
	return &repozytoriumTozsamosci{zapytania: z}
}

// Kategorie zwraca katalog kategorii w porządku jednoznacznym według warstwy,
// kolejności i kodu. Katalog pusty nie jest błędem.
func (r *repozytoriumTozsamosci) Kategorie(ctx context.Context, tylkoAktywne bool) ([]KategoriaTozsamosci, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKategoriiTozsamosci)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, liczbaLogiczna(tylkoAktywne))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać katalogu kategorii tożsamości: %w", err)
	}
	defer wiersze.Close()

	lista := []KategoriaTozsamosci{}
	for wiersze.Next() {
		kategoria, err := odczytajKategorieTozsamosci(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, kategoria)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt katalogu kategorii tożsamości: %w", err)
	}
	return lista, nil
}

// KategoriaPoKodzie zwraca jedną pozycję katalogu. Brak wiersza jest sygnałem
// ErrBrakWiersza, nie awarią odczytu.
func (r *repozytoriumTozsamosci) KategoriaPoKodzie(ctx context.Context, kod string) (KategoriaTozsamosci, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKategorieTozsamosci)
	if err != nil {
		return KategoriaTozsamosci{}, err
	}
	kategoria, err := odczytajKategorieTozsamosci(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return KategoriaTozsamosci{}, fmt.Errorf(
			"dane: kategoria tożsamości %q nie istnieje w katalogu: %w", kod, ErrBrakWiersza)
	}
	return kategoria, err
}

// odczytajKategorieTozsamosci składa strukturę kategorii tożsamości z jednego
// wiersza wyniku zapytania, tłumacząc kolumny słownikowe na wartości kontraktu.
func odczytajKategorieTozsamosci(wiersz skaner) (KategoriaTozsamosci, error) {
	var kategoria KategoriaTozsamosci
	var warstwa, tryb string
	var obowiazkowa, aktywna int
	err := wiersz.Scan(&kategoria.ID, &kategoria.Kod, &kategoria.Nazwa, &kategoria.Opis,
		&warstwa, &kategoria.Kolejnosc, &obowiazkowa, &tryb, &aktywna)
	if err != nil {
		return KategoriaTozsamosci{}, err
	}
	if kategoria.Warstwa, err = warstwaTozsamosciZBazy(warstwa); err != nil {
		return KategoriaTozsamosci{}, err
	}
	if kategoria.TrybDomyslny, err = trybTozsamosciZBazy(tryb); err != nil {
		return KategoriaTozsamosci{}, err
	}
	kategoria.Obowiazkowa = obowiazkowa != 0
	kategoria.Aktywna = aktywna != 0
	return kategoria, nil
}
