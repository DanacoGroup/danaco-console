// Plik prowadzi obszar Apps — plik warsztatu, trwałość okna Workspace modułu Apps; zapis jest operacją UPSERT po
// kluczu okna, warstwy i ścieżki, więc tabela ma jeden wiersz na tę trójkę, a historii wersji nie prowadzi.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// PlikWarsztatu to wiersz tabeli `plik_warsztatu_apps` — bieżąca treść pliku
// warstwy frontendu albo backendu produktu w oknie Apps.
type PlikWarsztatu struct {
	ID             int64
	Okno           string
	Warstwa        shared.AppWorkspaceLayer
	Sciezka        string
	Tresc          string
	Rozmiar        int64
	KomponentID    *string
	Utworzono      string
	Zaktualizowano string
}

const (
	kolumnyPlikuWarsztatu = `id, okno, warstwa, sciezka, tresc, rozmiar, komponent_id,
	                         utworzono, zaktualizowano`

	// UPSERT po kluczu okna, warstwy i ścieżki: zapis nadpisuje stan bieżący pliku warstwy, nie zakłada nowego wiersza historii.
	zapiszPlikWarsztatuApps = `INSERT INTO plik_warsztatu_apps
	                    (okno, warstwa, sciezka, tresc, rozmiar, komponent_id, konto_id)
	                    VALUES (?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                    ON CONFLICT(okno, warstwa, sciezka) DO UPDATE SET
	                        tresc = excluded.tresc,
	                        rozmiar = excluded.rozmiar,
	                        komponent_id = excluded.komponent_id,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                    WHERE ` + WarunekKonta

	pobierzPlikWarsztatuApps = `SELECT ` + kolumnyPlikuWarsztatu + `
	                    FROM plik_warsztatu_apps
	                    WHERE okno = ? AND warstwa = ? AND sciezka = ? AND ` + WarunekKonta

	listaPlikowWarsztatuApps = `SELECT ` + kolumnyPlikuWarsztatu + `
	                    FROM plik_warsztatu_apps
	                    WHERE okno = ? AND ` + WarunekKonta + `
	                    ORDER BY warstwa, sciezka`
)

// ZapiszPlikWarsztatu zapisuje plik warstwy frontendu albo backendu produktu
// jako stan bieżący (UPSERT po okno+warstwa+ścieżka) i zwraca wiersz po
// zapisie, z rozmiarem i znacznikiem czasu policzonymi przez bazę.
func (r *repozytoriumAplikacji) ZapiszPlikWarsztatu(ctx context.Context, plik PlikWarsztatu) (PlikWarsztatu, error) {
	if plik.Okno == "" || plik.Sciezka == "" {
		return PlikWarsztatu{}, fmt.Errorf("dane: plik warsztatu bez okna albo bez ścieżki")
	}
	if plik.Warstwa != shared.AppWorkspaceLayerFrontend && plik.Warstwa != shared.AppWorkspaceLayerBackend {
		return PlikWarsztatu{}, fmt.Errorf("dane: plik warsztatu %q ma nieznaną warstwę %q", plik.Sciezka, plik.Warstwa)
	}
	rozmiar := plik.Rozmiar
	if rozmiar == 0 {
		rozmiar = int64(len(plik.Tresc))
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPlikWarsztatuApps)
	if err != nil {
		return PlikWarsztatu{}, err
	}
	konto := KontoOperatora(ctx)
	wynik, err := polecenie.ExecContext(ctx, plik.Okno, string(plik.Warstwa), plik.Sciezka,
		plik.Tresc, rozmiar, tekstDoKolumny(plik.KomponentID), konto, konto)
	if err != nil {
		return PlikWarsztatu{}, fmt.Errorf("dane: nie można zapisać pliku warsztatu %q okna %q: %w", plik.Sciezka, plik.Okno, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "plik warsztatu", plik.Sciezka); err != nil {
		return PlikWarsztatu{}, err
	}
	return r.jedenPlikWarsztatu(ctx, plik.Okno, plik.Warstwa, plik.Sciezka)
}

// PlikiWarsztatu zwraca wszystkie pliki warsztatu okna, obu warstw razem —
// Workspace rozdziela je po stronie widoku po kolumnie `Warstwa`.
func (r *repozytoriumAplikacji) PlikiWarsztatu(ctx context.Context, okno string) ([]PlikWarsztatu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPlikowWarsztatuApps)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać plików warsztatu okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []PlikWarsztatu{}
	for wiersze.Next() {
		plik, err := odczytajPlikWarsztatu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny plik warsztatu okna %q: %w", okno, err)
		}
		lista = append(lista, plik)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt plików warsztatu okna %q: %w", okno, err)
	}
	return lista, nil
}

// PlikWarsztatu oddaje na zewnątrz odczyt po kluczu naturalnym, wołając ten sam odczyt, którym repozytorium zwraca stan po zapisie.
func (r *repozytoriumAplikacji) PlikWarsztatu(ctx context.Context, okno string,
	warstwa shared.AppWorkspaceLayer, sciezka string) (PlikWarsztatu, error) {

	return r.jedenPlikWarsztatu(ctx, okno, warstwa, sciezka)
}

// jedenPlikWarsztatu odczytuje wiersz po kluczu naturalnym — używane po
// zapisie, żeby zwrócić stan policzony przez bazę (rozmiar, zaktualizowano).
func (r *repozytoriumAplikacji) jedenPlikWarsztatu(ctx context.Context, okno string,
	warstwa shared.AppWorkspaceLayer, sciezka string) (PlikWarsztatu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPlikWarsztatuApps)
	if err != nil {
		return PlikWarsztatu{}, err
	}
	plik, err := odczytajPlikWarsztatu(polecenie.QueryRowContext(ctx, okno, string(warstwa), sciezka,
		KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return PlikWarsztatu{}, ErrBrakWiersza
	}
	if err != nil {
		return PlikWarsztatu{}, fmt.Errorf("dane: nieczytelny plik warsztatu %q okna %q: %w", sciezka, okno, err)
	}
	return plik, nil
}

// odczytajPlikWarsztatu składa strukturę pliku warsztatu wprost z jednego wiersza wyniku zapytania SQL.
func odczytajPlikWarsztatu(wiersz skaner) (PlikWarsztatu, error) {
	var plik PlikWarsztatu
	var warstwa string
	var komponentID sql.NullString
	err := wiersz.Scan(&plik.ID, &plik.Okno, &warstwa, &plik.Sciezka, &plik.Tresc,
		&plik.Rozmiar, &komponentID, &plik.Utworzono, &plik.Zaktualizowano)
	if err != nil {
		return PlikWarsztatu{}, err
	}
	plik.Warstwa = shared.AppWorkspaceLayer(warstwa)
	plik.KomponentID = tekstZKolumny(komponentID)
	return plik, nil
}
