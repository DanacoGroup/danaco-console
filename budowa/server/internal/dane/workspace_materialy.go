// Plik niesie byty projektu będące materiałem: tablicę wizualną, wskaźnik
// treści plików i pozycje kalendarza wciągnięte z iCal. Niesie też zmianę
// stanu projektu i odłączenie eksperta.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// TablicaWorkspace to wiersz tablicy wizualnej projektu: karta ze sceną, nazwą
// i znacznikami czasu, powiązana z projektem przez jego identyfikator.
type TablicaWorkspace struct {
	ProjektID      int64
	ProjektKod     string
	Identyfikator  string
	Nazwa          string
	Scena          string
	Utworzono      string
	Zaktualizowano string
}

// WyciagWorkspace to wiersz wskaźnika treści jednego pliku projektu: sposób
// wydobycia, sama treść, liczba znaków oraz rozpoznane języki treści.
type WyciagWorkspace struct {
	ProjektID    int64
	Plik         string
	Sposob       shared.WorkspaceExtractionMethod
	Tresc        string
	LiczbaZnakow int
	Jezyki       []string
	Utworzono    string
}

// PozycjaKalendarzaWorkspace to wiersz kalendarza projektu wciągnięty z iCal:
// tytuł, granice czasowe, znacznik kamienia milowego i reguła powtarzalności.
type PozycjaKalendarzaWorkspace struct {
	ProjektID            int64
	Identyfikator        string
	Tytul                string
	PoczatekMs           int64
	KoniecMs             int64
	CalyDzien            bool
	KamienMilowy         bool
	RegulaPowtarzalnosci string
	UidZewnetrzny        string
}

var (
	kolumnyTablicyWorkspace = `t.projekt_id, p.kod, t.identyfikator_zewnetrzny, t.nazwa, t.scena,
	                           t.utworzono, t.zaktualizowano`

	zrodloTablicyWorkspace = ` FROM tablica_wizualna_projektu t JOIN projekt p ON p.id = t.projekt_id`

	zapiszTabliceWorkspace = `INSERT INTO tablica_wizualna_projektu
	    (identyfikator_zewnetrzny, projekt_id, nazwa, scena)
	    VALUES (?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        nazwa = excluded.nazwa, scena = excluded.scena,
	        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzTabliceWorkspace = `SELECT ` + kolumnyTablicyWorkspace + zrodloTablicyWorkspace +
		` WHERE t.identyfikator_zewnetrzny = ? AND ` + warunekKontaProjektu

	listaTablicWorkspace = `SELECT ` + kolumnyTablicyWorkspace + zrodloTablicyWorkspace +
		` WHERE t.projekt_id = ? AND ` + warunekKontaProjektu + ` ORDER BY t.id`

	zapiszWyciagWorkspace = `INSERT INTO wyciag_tekstu_projektu
	    (projekt_id, plik, sposob, tresc, liczba_znakow, jezyki)
	    VALUES (?, ?, ?, ?, ?, ?)
	    ON CONFLICT(projekt_id, plik) DO UPDATE SET
	        sposob = excluded.sposob, tresc = excluded.tresc,
	        liczba_znakow = excluded.liczba_znakow, jezyki = excluded.jezyki,
	        utworzono = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzWyciagWorkspace = `SELECT projekt_id, plik, sposob, tresc, liczba_znakow, jezyki, utworzono
	                          FROM wyciag_tekstu_projektu WHERE projekt_id = ? AND plik = ?`

	listaWyciagowWorkspace = `SELECT projekt_id, plik, sposob, tresc, liczba_znakow, jezyki, utworzono
	                          FROM wyciag_tekstu_projektu WHERE projekt_id = ? ORDER BY plik`

	usunWyciagWorkspace = `DELETE FROM wyciag_tekstu_projektu WHERE projekt_id = ? AND plik = ?`

	zapiszPozycjeKalendarzaWorkspace = `INSERT INTO pozycja_kalendarza_projektu
	    (identyfikator_zewnetrzny, projekt_id, tytul, poczatek_ms, koniec_ms, caly_dzien,
	     kamien_milowy, regula_powtarzalnosci, uid_zewnetrzny)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	    ON CONFLICT(projekt_id, uid_zewnetrzny) DO UPDATE SET
	        tytul = excluded.tytul, poczatek_ms = excluded.poczatek_ms,
	        koniec_ms = excluded.koniec_ms, caly_dzien = excluded.caly_dzien,
	        kamien_milowy = excluded.kamien_milowy,
	        regula_powtarzalnosci = excluded.regula_powtarzalnosci`

	listaPozycjiKalendarzaWorkspace = `SELECT projekt_id, identyfikator_zewnetrzny, tytul,
	                                   poczatek_ms, koniec_ms, caly_dzien, kamien_milowy,
	                                   regula_powtarzalnosci, uid_zewnetrzny
	                                   FROM pozycja_kalendarza_projektu
	                                   WHERE projekt_id = ? ORDER BY poczatek_ms, id`

	ustawStanProjektuWorkspace = `UPDATE projekt SET stan = ?,
	                              zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                              WHERE id = ? AND ` + WarunekKonta

	odlaczAgentaWorkspace = `DELETE FROM przypisanie_agenta_projektu
	                         WHERE projekt_id = ? AND agent_kod = ?`
)

// ZapiszTabliceWorkspace zakłada albo zmienia tablicę wizualną projektu:
// przy zgodnym identyfikatorze zewnętrznym nadpisuje nazwę i scenę tablicy.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszTabliceWorkspace(ctx context.Context,
	tablica TablicaWorkspace) (TablicaWorkspace, error) {

	if tablica.Identyfikator == "" {
		return TablicaWorkspace{}, fmt.Errorf("dane: tablica wizualna bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszTabliceWorkspace)
	if err != nil {
		return TablicaWorkspace{}, err
	}
	_, err = polecenie.ExecContext(ctx, tablica.Identyfikator, tablica.ProjektID, tablica.Nazwa,
		tablica.Scena)
	if err != nil {
		return TablicaWorkspace{}, fmt.Errorf("dane: nie można zapisać tablicy wizualnej %q: %w",
			tablica.Identyfikator, err)
	}
	return r.TablicaWorkspace(ctx, tablica.Identyfikator)
}

// TablicaWorkspace zwraca jedną tablicę wizualną po identyfikatorze zewnętrznym.
// Brak wiersza w bazie skutkuje błędem ErrBrakWiersza.
func (r *repozytoriumPrzestrzeniRoboczej) TablicaWorkspace(ctx context.Context,
	identyfikator string) (TablicaWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzTabliceWorkspace)
	if err != nil {
		return TablicaWorkspace{}, err
	}
	tablica, err := odczytajTabliceWorkspace(polecenie.QueryRowContext(ctx, identyfikator, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return TablicaWorkspace{}, ErrBrakWiersza
	}
	if err != nil {
		return TablicaWorkspace{}, fmt.Errorf("dane: nieczytelna tablica wizualna %q: %w",
			identyfikator, err)
	}
	return tablica, nil
}

// TabliceWorkspace zwraca komplet tablic wizualnych projektu, uporządkowany
// według kolejności założenia wierszy w bazie.
func (r *repozytoriumPrzestrzeniRoboczej) TabliceWorkspace(ctx context.Context,
	projektID int64) ([]TablicaWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaTablicWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać tablic projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []TablicaWorkspace{}
	for wiersze.Next() {
		tablica, err := odczytajTabliceWorkspace(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz tablicy wizualnej: %w", err)
		}
		lista = append(lista, tablica)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt tablic projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// ZapiszWyciagWorkspace odkłada treść wydobytą z pliku projektu, sposób jej
// wydobycia oraz rozpoznane języki; przy zgodnej ścieżce nadpisuje wcześniejszy wiersz.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszWyciagWorkspace(ctx context.Context,
	wyciag WyciagWorkspace) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWyciagWorkspace)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wyciag.ProjektID, wyciag.Plik, string(wyciag.Sposob),
		wyciag.Tresc, wyciag.LiczbaZnakow, strings.Join(wyciag.Jezyki, "\n"))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać wyciągu treści pliku %q: %w", wyciag.Plik, err)
	}
	return nil
}

// WyciagWorkspace zwraca wyciąg treści jednego pliku projektu, wskazanego
// identyfikatorem projektu i ścieżką pliku.
func (r *repozytoriumPrzestrzeniRoboczej) WyciagWorkspace(ctx context.Context,
	projektID int64, plik string) (WyciagWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWyciagWorkspace)
	if err != nil {
		return WyciagWorkspace{}, err
	}
	wyciag, err := odczytajWyciagWorkspace(polecenie.QueryRowContext(ctx, projektID, plik))
	if errors.Is(err, sql.ErrNoRows) {
		return WyciagWorkspace{}, ErrBrakWiersza
	}
	if err != nil {
		return WyciagWorkspace{}, fmt.Errorf("dane: nieczytelny wyciąg treści pliku %q: %w", plik, err)
	}
	return wyciag, nil
}

// WyciagiWorkspace zwraca komplet wyciągów treści projektu — materiał
// wyszukiwania po słowach nad plikami.
func (r *repozytoriumPrzestrzeniRoboczej) WyciagiWorkspace(ctx context.Context,
	projektID int64) ([]WyciagWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWyciagowWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wyciągów projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []WyciagWorkspace{}
	for wiersze.Next() {
		wyciag, err := odczytajWyciagWorkspace(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wyciągu treści: %w", err)
		}
		lista = append(lista, wyciag)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wyciągów projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// UsunWyciagWorkspace zdejmuje wyciąg pliku, który przestał istnieć — na
// przykład po scaleniu duplikatów.
func (r *repozytoriumPrzestrzeniRoboczej) UsunWyciagWorkspace(ctx context.Context,
	projektID int64, plik string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, usunWyciagWorkspace)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, projektID, plik); err != nil {
		return fmt.Errorf("dane: nie można usunąć wyciągu treści pliku %q: %w", plik, err)
	}
	return nil
}

// ZapiszPozycjeKalendarzaWorkspace odkłada pozycję kalendarza wciągniętą z iCal;
// przy zgodnym identyfikatorze zewnętrznym nadpisuje jej pola.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszPozycjeKalendarzaWorkspace(ctx context.Context,
	pozycja PozycjaKalendarzaWorkspace) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPozycjeKalendarzaWorkspace)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, pozycja.Identyfikator, pozycja.ProjektID, pozycja.Tytul,
		pozycja.PoczatekMs, pozycja.KoniecMs, liczbaLogiczna(pozycja.CalyDzien),
		liczbaLogiczna(pozycja.KamienMilowy), pozycja.RegulaPowtarzalnosci, pozycja.UidZewnetrzny)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać pozycji kalendarza %q: %w", pozycja.Tytul, err)
	}
	return nil
}

// PozycjeKalendarzaWorkspace zwraca pozycje kalendarza projektu wciągnięte
// z iCal, w kolejności początków.
func (r *repozytoriumPrzestrzeniRoboczej) PozycjeKalendarzaWorkspace(ctx context.Context,
	projektID int64) ([]PozycjaKalendarzaWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPozycjiKalendarzaWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kalendarza projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []PozycjaKalendarzaWorkspace{}
	for wiersze.Next() {
		var pozycja PozycjaKalendarzaWorkspace
		var calyDzien, kamien int
		err := wiersze.Scan(&pozycja.ProjektID, &pozycja.Identyfikator, &pozycja.Tytul,
			&pozycja.PoczatekMs, &pozycja.KoniecMs, &calyDzien, &kamien,
			&pozycja.RegulaPowtarzalnosci, &pozycja.UidZewnetrzny)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kalendarza projektu: %w", err)
		}
		pozycja.CalyDzien, pozycja.KamienMilowy = calyDzien == 1, kamien == 1
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kalendarza projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// UstawStanProjektuWorkspace zapisuje stan projektu (aktywny, wstrzymany,
// zarchiwizowany) i przesuwa znacznik ostatniej zmiany.
func (r *repozytoriumPrzestrzeniRoboczej) UstawStanProjektuWorkspace(ctx context.Context,
	projektID int64, stan shared.WorkspaceProjectStatus) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawStanProjektuWorkspace)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, string(stan), projektID, KontoOperatora(ctx)); err != nil {
		return fmt.Errorf("dane: nie można zapisać stanu projektu %d: %w", projektID, err)
	}
	return nil
}

// OdlaczAgentaWorkspace znosi przypisanie eksperta do projektu. Fałsz znaczy,
// że takiego przypisania nie było.
func (r *repozytoriumPrzestrzeniRoboczej) OdlaczAgentaWorkspace(ctx context.Context,
	projektID int64, agent string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, odlaczAgentaWorkspace)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, projektID, agent)
	if err != nil {
		return false, fmt.Errorf("dane: nie można odłączyć eksperta %q: %w", agent, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany wynik odłączenia eksperta %q: %w", agent, err)
	}
	return zmienione > 0, nil
}

// odczytajTabliceWorkspace składa strukturę TablicaWorkspace z jednego wiersza
// wyniku zapytania, w kolejności kolumn kolumnyTablicyWorkspace.
func odczytajTabliceWorkspace(wiersz skaner) (TablicaWorkspace, error) {
	var tablica TablicaWorkspace
	err := wiersz.Scan(&tablica.ProjektID, &tablica.ProjektKod, &tablica.Identyfikator,
		&tablica.Nazwa, &tablica.Scena, &tablica.Utworzono, &tablica.Zaktualizowano)
	if err != nil {
		return TablicaWorkspace{}, err
	}
	return tablica, nil
}

// odczytajWyciagWorkspace składa strukturę WyciagWorkspace z jednego wiersza
// wyniku, rozdzielając zapisane w bazie języki treści na listę.
func odczytajWyciagWorkspace(wiersz skaner) (WyciagWorkspace, error) {
	var wyciag WyciagWorkspace
	var sposob, jezyki string
	err := wiersz.Scan(&wyciag.ProjektID, &wyciag.Plik, &sposob, &wyciag.Tresc,
		&wyciag.LiczbaZnakow, &jezyki, &wyciag.Utworzono)
	if err != nil {
		return WyciagWorkspace{}, err
	}
	wyciag.Sposob = shared.WorkspaceExtractionMethod(sposob)
	wyciag.Jezyki = rozdzielWierszamiWorkspace(jezyki)
	return wyciag, nil
}
