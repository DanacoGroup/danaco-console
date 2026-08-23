// Odpowiedzialność pliku: treść kategorii zasad zapisana per oś (tabela
// `dokument_tozsamosci`). Oś mówi, dla czego treść obowiązuje — dla platformy,
// dla wskazanego modelu albo dla wskazanego konta. Która oś wygrywa, rozstrzyga
// warstwa wyższa (`core/tozsamosc_wybor.go`); repozytorium wyłącznie oddaje
// wiersze i zapisuje zmianę Operatora.
//
// Klucz zapisu jest trójką kategoria + oś + byt osi. Ta sama trójka jest
// warunkiem UNIQUE w schemacie, więc zapis jest nadpisaniem, nie mnożeniem
// wierszy o tym samym znaczeniu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

const (
	kolumnyDokumentuTozsamosci = `d.id, d.kategoria_id, k.kod, d.os, d.os_byt, d.tryb,
	                              d.tresc, d.odcisk_tresci, d.aktywny, d.zaktualizowano`

	zrodloDokumentuTozsamosci = ` FROM dokument_tozsamosci d
	                              JOIN kategoria_tozsamosci k ON k.id = d.kategoria_id`

	listaDokumentowTozsamosci = `SELECT ` + kolumnyDokumentuTozsamosci + zrodloDokumentuTozsamosci + `
	                             WHERE (? = '' OR k.kod = ?)
	                               AND (? = '' OR d.os = ?)
	                               AND (? = '' OR d.os_byt = ?)
	                               AND (? = 0 OR d.aktywny = 1)
	                             ORDER BY k.kolejnosc, k.kod, d.os, d.os_byt, d.id`

	pobierzDokumentTozsamosci = `SELECT ` + kolumnyDokumentuTozsamosci + zrodloDokumentuTozsamosci + `
	                             WHERE k.kod = ? AND d.os = ? AND d.os_byt = ?`

	zapiszDokumentTozsamosci = `INSERT INTO dokument_tozsamosci
	                            (kategoria_id, os, os_byt, tryb, tresc, odcisk_tresci, aktywny)
	                            VALUES ((SELECT id FROM kategoria_tozsamosci WHERE kod = ?),
	                                    ?, ?, ?, ?, ?, ?)
	                            ON CONFLICT(kategoria_id, os, os_byt) DO UPDATE SET
	                                tryb = excluded.tryb,
	                                tresc = excluded.tresc,
	                                odcisk_tresci = excluded.odcisk_tresci,
	                                aktywny = excluded.aktywny,
	                                zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	usunDokumentTozsamosci = `DELETE FROM dokument_tozsamosci WHERE id = ?`
)

// Dokumenty zwraca zapisy treści zawężone filtrem. Brak zapisów nie jest błędem
// — znaczy „treść bierze się z osi szerszej”, a w jej braku kategoria po prostu
// nie wchodzi do nakładki.
func (r *repozytoriumTozsamosci) Dokumenty(ctx context.Context, filtr FiltrTozsamosci) ([]DokumentTozsamosci, error) {
	os := ""
	if filtr.Os != "" {
		kolumna, err := osTozsamosciNaBaze(filtr.Os)
		if err != nil {
			return nil, err
		}
		os = kolumna
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaDokumentowTozsamosci)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx,
		filtr.KodKategorii, filtr.KodKategorii, os, os, filtr.OsByt, filtr.OsByt,
		liczbaLogiczna(filtr.TylkoAktywne))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać treści tożsamości modelu: %w", err)
	}
	defer wiersze.Close()

	lista := []DokumentTozsamosci{}
	for wiersze.Next() {
		dokument, err := odczytajDokumentTozsamosci(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, dokument)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt treści tożsamości modelu: %w", err)
	}
	return lista, nil
}

// ZapiszDokument zapisuje treść kategorii dla jednej osi i oddaje wiersz po
// zmianie. Wiersz istniejący nadpisuje — Operator zmienia treść, nie zakłada
// drugiej obok. Kategoria spoza katalogu jest błędem wskazania, nie awarią:
// katalog jest zamkniętym zbiorem wierszy migracji.
func (r *repozytoriumTozsamosci) ZapiszDokument(ctx context.Context, dokument DokumentTozsamosci) (DokumentTozsamosci, error) {
	os, err := osTozsamosciNaBaze(dokument.Os)
	if err != nil {
		return DokumentTozsamosci{}, err
	}
	tryb, err := trybTozsamosciNaBaze(dokument.Tryb)
	if err != nil {
		return DokumentTozsamosci{}, err
	}
	if err := sprawdzBytOsiTozsamosci(shared.ConfigAxis(os), dokument.OsByt); err != nil {
		return DokumentTozsamosci{}, err
	}
	if _, err := r.KategoriaPoKodzie(ctx, dokument.KodKategorii); err != nil {
		return DokumentTozsamosci{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszDokumentTozsamosci)
	if err != nil {
		return DokumentTozsamosci{}, err
	}
	_, err = polecenie.ExecContext(ctx, dokument.KodKategorii, os, dokument.OsByt, tryb,
		dokument.Tresc, dokument.OdciskTresci, liczbaLogiczna(dokument.Aktywny))
	if err != nil {
		return DokumentTozsamosci{}, fmt.Errorf(
			"dane: nie można zapisać treści kategorii %q dla osi %q: %w",
			dokument.KodKategorii, os, err)
	}
	return r.dokument(ctx, dokument.KodKategorii, os, dokument.OsByt)
}

// UsunDokument kasuje zapis treści. Drugi wynik mówi, czy wiersz w ogóle
// istniał — usunięcie zapisu, którego nie ma, nie jest awarią.
func (r *repozytoriumTozsamosci) UsunDokument(ctx context.Context, id int64) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunDokumentTozsamosci)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, id)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć treści tożsamości %d: %w", id, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieczytelny wynik usunięcia treści tożsamości %d: %w", id, err)
	}
	return usuniete > 0, nil
}

// dokument zwraca jeden zapis wskazany trójką klucza — po zapisie oddaje wiersz
// z czasem nadanym przez bazę, żeby warstwa wyższa nie zgadywała znacznika.
func (r *repozytoriumTozsamosci) dokument(ctx context.Context, kod, os, osByt string) (DokumentTozsamosci, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzDokumentTozsamosci)
	if err != nil {
		return DokumentTozsamosci{}, err
	}
	dokument, err := odczytajDokumentTozsamosci(polecenie.QueryRowContext(ctx, kod, os, osByt))
	if errors.Is(err, sql.ErrNoRows) {
		return DokumentTozsamosci{}, fmt.Errorf(
			"dane: treść kategorii %q dla osi %q nie istnieje: %w", kod, os, ErrBrakWiersza)
	}
	return dokument, err
}

// odczytajDokumentTozsamosci składa strukturę z jednego wiersza wyniku.
func odczytajDokumentTozsamosci(wiersz skaner) (DokumentTozsamosci, error) {
	var dokument DokumentTozsamosci
	var os, tryb string
	var aktywny int
	err := wiersz.Scan(&dokument.ID, &dokument.KategoriaID, &dokument.KodKategorii, &os,
		&dokument.OsByt, &tryb, &dokument.Tresc, &dokument.OdciskTresci, &aktywny,
		&dokument.Zaktualizowano)
	if err != nil {
		return DokumentTozsamosci{}, err
	}
	if dokument.Os, err = osTozsamosciZBazy(os); err != nil {
		return DokumentTozsamosci{}, err
	}
	if dokument.Tryb, err = trybTozsamosciZBazy(tryb); err != nil {
		return DokumentTozsamosci{}, err
	}
	dokument.Aktywny = aktywny != 0
	return dokument, nil
}
