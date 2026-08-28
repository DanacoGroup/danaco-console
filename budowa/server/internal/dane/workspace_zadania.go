// Odpowiedzialność pliku: zadania projektu (tabela `zadanie_projektu`) oraz kolumny tablicy
// kanban (tabela `kolumna_tablicy_projektu`) — hub planowania modułu Workspace.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// ZadanieWorkspace to wiersz zadania projektu w tabeli `zadanie_projektu` modułu Workspace, złączony z projektem.
type ZadanieWorkspace struct {
	ID                   int64
	ProjektID            int64
	ProjektKod           string
	Identyfikator        string
	Tytul                string
	Opis                 string
	Stan                 shared.WorkspaceTaskStatus
	Waga                 shared.WorkspaceTaskPriority
	RodzajWykonawcy      string
	Wykonawca            string
	ZadanieNadrzedne     string
	KolumnaTablicy       string
	KluczPorzadkowy      string
	PoczatekMs           int64
	TerminMs             int64
	UkonczonoMs          int64
	SzacunekMinut        int
	SpedzonoMinut        int
	PostepProcent        int
	KamienMilowy         bool
	RegulaPowtarzalnosci string
	Etykiety             []string
	ListaKontrolnaJson   string
	Utworzono            string
	Zaktualizowano       string
}

// KolumnaTablicyWorkspace to wiersz kolumny tablicy kanban projektu w tabeli `kolumna_tablicy_projektu`.
type KolumnaTablicyWorkspace struct {
	ProjektID     int64
	Identyfikator string
	Nazwa         string
	Stan          shared.WorkspaceTaskStatus
	Kolejnosc     int
	GranicaWip    int
}

const (
	kolumnyZadaniaWorkspace = `z.id, p.kod, z.projekt_id, z.identyfikator_zewnetrzny, z.tytul, z.opis,
	                           z.stan, z.waga, z.rodzaj_wykonawcy, z.wykonawca, z.zadanie_nadrzedne,
	                           z.kolumna_tablicy, z.klucz_porzadkowy, z.poczatek_ms, z.termin_ms,
	                           z.ukonczono_ms, z.szacunek_minut, z.spedzono_minut, z.postep_procent,
	                           z.kamien_milowy, z.regula_powtarzalnosci, z.etykiety,
	                           z.lista_kontrolna, z.utworzono, z.zaktualizowano`

	zrodloZadaniaWorkspace = ` FROM zadanie_projektu z JOIN projekt p ON p.id = z.projekt_id`

	zapiszZadanieWorkspace = `INSERT INTO zadanie_projektu
	    (identyfikator_zewnetrzny, projekt_id, tytul, opis, stan, waga, rodzaj_wykonawcy,
	     wykonawca, zadanie_nadrzedne, kolumna_tablicy, klucz_porzadkowy, poczatek_ms,
	     termin_ms, ukonczono_ms, szacunek_minut, spedzono_minut, postep_procent,
	     kamien_milowy, regula_powtarzalnosci, etykiety, lista_kontrolna)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        tytul = excluded.tytul, opis = excluded.opis, stan = excluded.stan,
	        waga = excluded.waga, rodzaj_wykonawcy = excluded.rodzaj_wykonawcy,
	        wykonawca = excluded.wykonawca, zadanie_nadrzedne = excluded.zadanie_nadrzedne,
	        kolumna_tablicy = excluded.kolumna_tablicy,
	        klucz_porzadkowy = excluded.klucz_porzadkowy,
	        poczatek_ms = excluded.poczatek_ms, termin_ms = excluded.termin_ms,
	        ukonczono_ms = excluded.ukonczono_ms, szacunek_minut = excluded.szacunek_minut,
	        spedzono_minut = excluded.spedzono_minut, postep_procent = excluded.postep_procent,
	        kamien_milowy = excluded.kamien_milowy,
	        regula_powtarzalnosci = excluded.regula_powtarzalnosci,
	        etykiety = excluded.etykiety, lista_kontrolna = excluded.lista_kontrolna,
	        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzZadanieWorkspace = `SELECT ` + kolumnyZadaniaWorkspace + zrodloZadaniaWorkspace +
		` WHERE z.identyfikator_zewnetrzny = ?`

	listaZadanWorkspace = `SELECT ` + kolumnyZadaniaWorkspace + zrodloZadaniaWorkspace +
		` WHERE z.projekt_id = ? ORDER BY z.kamien_milowy DESC, z.termin_ms > 0 DESC,
		   z.termin_ms, z.id`

	usunZadanieWorkspace = `DELETE FROM zadanie_projektu WHERE identyfikator_zewnetrzny = ?`

	zapiszKolumneTablicyWorkspace = `INSERT INTO kolumna_tablicy_projektu
	    (identyfikator_zewnetrzny, projekt_id, nazwa, stan, kolejnosc, granica_wip)
	    VALUES (?, ?, ?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        nazwa = excluded.nazwa, stan = excluded.stan,
	        kolejnosc = excluded.kolejnosc, granica_wip = excluded.granica_wip`

	listaKolumnTablicyWorkspace = `SELECT projekt_id, identyfikator_zewnetrzny, nazwa, stan,
	                               kolejnosc, granica_wip
	                               FROM kolumna_tablicy_projektu
	                               WHERE projekt_id = ? ORDER BY kolejnosc, id`
)

// ZapiszZadanieWorkspace zakłada zadanie albo zmienia zadanie o tym samym
// identyfikatorze i oddaje jego stan po zapisie.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszZadanieWorkspace(ctx context.Context,
	zadanie ZadanieWorkspace) (ZadanieWorkspace, error) {

	if zadanie.Identyfikator == "" {
		return ZadanieWorkspace{}, fmt.Errorf("dane: zadanie projektu bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZadanieWorkspace)
	if err != nil {
		return ZadanieWorkspace{}, err
	}
	_, err = polecenie.ExecContext(ctx, zadanie.Identyfikator, zadanie.ProjektID, zadanie.Tytul,
		zadanie.Opis, string(zadanie.Stan), string(zadanie.Waga), zadanie.RodzajWykonawcy,
		zadanie.Wykonawca, zadanie.ZadanieNadrzedne, zadanie.KolumnaTablicy,
		zadanie.KluczPorzadkowy, zadanie.PoczatekMs, zadanie.TerminMs, zadanie.UkonczonoMs,
		zadanie.SzacunekMinut, zadanie.SpedzonoMinut, zadanie.PostepProcent,
		liczbaLogiczna(zadanie.KamienMilowy), zadanie.RegulaPowtarzalnosci,
		strings.Join(zadanie.Etykiety, "\n"), zadanie.ListaKontrolnaJson)
	if err != nil {
		return ZadanieWorkspace{}, fmt.Errorf("dane: nie można zapisać zadania %q: %w",
			zadanie.Identyfikator, err)
	}
	return r.ZadanieWorkspace(ctx, zadanie.Identyfikator)
}

// ZadanieWorkspace zwraca jedno zadanie po jego identyfikatorze zewnętrznym; brak wraca jako ErrBrakWiersza.
func (r *repozytoriumPrzestrzeniRoboczej) ZadanieWorkspace(ctx context.Context,
	identyfikator string) (ZadanieWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZadanieWorkspace)
	if err != nil {
		return ZadanieWorkspace{}, err
	}
	zadanie, err := odczytajZadanieWorkspace(polecenie.QueryRowContext(ctx, identyfikator))
	if errors.Is(err, sql.ErrNoRows) {
		return ZadanieWorkspace{}, ErrBrakWiersza
	}
	if err != nil {
		return ZadanieWorkspace{}, fmt.Errorf("dane: nieczytelny wiersz zadania %q: %w",
			identyfikator, err)
	}
	return zadanie, nil
}

// ZadaniaWorkspace zwraca wszystkie zadania projektu w całości, uporządkowane według kamienia milowego i terminu.
func (r *repozytoriumPrzestrzeniRoboczej) ZadaniaWorkspace(ctx context.Context,
	projektID int64) ([]ZadanieWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZadanWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zadań projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []ZadanieWorkspace{}
	for wiersze.Next() {
		zadanie, err := odczytajZadanieWorkspace(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zadania projektu: %w", err)
		}
		lista = append(lista, zadanie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zadań projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// UsunZadaniaWorkspace usuwa wskazane zadania wraz z ich zależnościami
// w jednej transakcji i oddaje liczbę zadań naprawdę usuniętych.
func (r *repozytoriumPrzestrzeniRoboczej) UsunZadaniaWorkspace(ctx context.Context,
	identyfikatory []string) (int, error) {

	usuniete := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		for _, identyfikator := range identyfikatory {
			wynik, err := transakcja.ExecContext(ctx, usunZadanieWorkspace, identyfikator)
			if err != nil {
				return fmt.Errorf("dane: nie można usunąć zadania %q: %w", identyfikator, err)
			}
			zmienione, err := wynik.RowsAffected()
			if err != nil {
				return fmt.Errorf("dane: nieznany wynik usunięcia zadania %q: %w", identyfikator, err)
			}
			usuniete += int(zmienione)
			_, err = transakcja.ExecContext(ctx,
				`DELETE FROM zaleznosc_zadan_projektu WHERE poprzednik = ? OR nastepnik = ?`,
				identyfikator, identyfikator)
			if err != nil {
				return fmt.Errorf("dane: nie można usunąć zależności zadania %q: %w",
					identyfikator, err)
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return usuniete, nil
}

// ZapiszKolumneTablicyWorkspace zakłada albo zmienia kolumnę tablicy kanban projektu w bazie danych aplikacji.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszKolumneTablicyWorkspace(ctx context.Context,
	kolumna KolumnaTablicyWorkspace) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKolumneTablicyWorkspace)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, kolumna.Identyfikator, kolumna.ProjektID, kolumna.Nazwa,
		string(kolumna.Stan), kolumna.Kolejnosc, kolumna.GranicaWip)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać kolumny tablicy %q: %w",
			kolumna.Identyfikator, err)
	}
	return nil
}

// KolumnyTablicyWorkspace zwraca kolumny tablicy kanban projektu w kolejności nastawionej pola `Kolejnosc`.
func (r *repozytoriumPrzestrzeniRoboczej) KolumnyTablicyWorkspace(ctx context.Context,
	projektID int64) ([]KolumnaTablicyWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaKolumnTablicyWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kolumn tablicy projektu %d: %w",
			projektID, err)
	}
	defer wiersze.Close()

	lista := []KolumnaTablicyWorkspace{}
	for wiersze.Next() {
		var kolumna KolumnaTablicyWorkspace
		var stan string
		err := wiersze.Scan(&kolumna.ProjektID, &kolumna.Identyfikator, &kolumna.Nazwa, &stan,
			&kolumna.Kolejnosc, &kolumna.GranicaWip)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kolumny tablicy: %w", err)
		}
		kolumna.Stan = shared.WorkspaceTaskStatus(stan)
		lista = append(lista, kolumna)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kolumn tablicy projektu %d: %w",
			projektID, err)
	}
	return lista, nil
}

// odczytajZadanieWorkspace składa strukturę zadania z jednego wiersza wyniku zapytania SQL bazy danych.
func odczytajZadanieWorkspace(wiersz skaner) (ZadanieWorkspace, error) {
	var zadanie ZadanieWorkspace
	var stan, waga, etykiety string
	var kamien int
	err := wiersz.Scan(&zadanie.ID, &zadanie.ProjektKod, &zadanie.ProjektID,
		&zadanie.Identyfikator, &zadanie.Tytul, &zadanie.Opis, &stan, &waga,
		&zadanie.RodzajWykonawcy, &zadanie.Wykonawca, &zadanie.ZadanieNadrzedne,
		&zadanie.KolumnaTablicy, &zadanie.KluczPorzadkowy, &zadanie.PoczatekMs,
		&zadanie.TerminMs, &zadanie.UkonczonoMs, &zadanie.SzacunekMinut,
		&zadanie.SpedzonoMinut, &zadanie.PostepProcent, &kamien,
		&zadanie.RegulaPowtarzalnosci, &etykiety, &zadanie.ListaKontrolnaJson,
		&zadanie.Utworzono, &zadanie.Zaktualizowano)
	if err != nil {
		return ZadanieWorkspace{}, err
	}
	zadanie.Stan = shared.WorkspaceTaskStatus(stan)
	zadanie.Waga = shared.WorkspaceTaskPriority(waga)
	zadanie.KamienMilowy = kamien == 1
	zadanie.Etykiety = rozdzielWierszamiWorkspace(etykiety)
	return zadanie, nil
}

// rozdzielWierszamiWorkspace rozbija kolumnę wielowartościową na wykaz. Pusta
// kolumna daje wykaz pusty, a nie wykaz z jedną pustą pozycją.
func rozdzielWierszamiWorkspace(kolumna string) []string {
	kolumna = strings.TrimSpace(kolumna)
	if kolumna == "" {
		return []string{}
	}
	czesci := strings.Split(kolumna, "\n")
	wykaz := make([]string, 0, len(czesci))
	for _, czesc := range czesci {
		if czesc = strings.TrimSpace(czesc); czesc != "" {
			wykaz = append(wykaz, czesc)
		}
	}
	return wykaz
}
