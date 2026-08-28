// Repozytorium przechowuje notatki i strony wiki projektu w tabeli
// `notatka_projektu`, wraz z odnośnikami treści w tabeli
// `odnosnik_notatki_projektu`.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// NotatkaWorkspace to wiersz tabeli `notatka_projektu`, reprezentujący
// notatkę albo stronę wiki projektu.
type NotatkaWorkspace struct {
	ID               int64
	ProjektID        int64
	ProjektKod       string
	Identyfikator    string
	Tytul            string
	Tresc            string
	NotatkaNadrzedna string
	Etykiety         []string
	Naglowki         []string
	RodzajAutora     shared.WorkspaceAssigneeKind
	Autor            string
	Utworzono        string
	Zaktualizowano   string
}

// OdnosnikWorkspace to wiersz odnośnika treści prowadzącego od jednej strony
// do nazwy innej strony projektu.
type OdnosnikWorkspace struct {
	ProjektID       int64
	NotatkaZrodlowa string
	NazwaDocelowa   string
	NotatkaDocelowa string
	Rodzaj          shared.WorkspaceGraphEdgeKind
	Kontekst        string
}

const (
	kolumnyNotatkiWorkspace = `n.id, p.kod, n.projekt_id, n.identyfikator_zewnetrzny, n.tytul,
	                           n.tresc, n.notatka_nadrzedna, n.etykiety, n.naglowki,
	                           n.rodzaj_autora, n.autor, n.utworzono, n.zaktualizowano`

	zrodloNotatkiWorkspace = ` FROM notatka_projektu n JOIN projekt p ON p.id = n.projekt_id`

	zapiszNotatkeWorkspace = `INSERT INTO notatka_projektu
	    (identyfikator_zewnetrzny, projekt_id, tytul, tresc, notatka_nadrzedna, etykiety,
	     naglowki, rodzaj_autora, autor)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	        tytul = excluded.tytul, tresc = excluded.tresc,
	        notatka_nadrzedna = excluded.notatka_nadrzedna, etykiety = excluded.etykiety,
	        naglowki = excluded.naglowki, rodzaj_autora = excluded.rodzaj_autora,
	        autor = excluded.autor,
	        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzNotatkeWorkspace = `SELECT ` + kolumnyNotatkiWorkspace + zrodloNotatkiWorkspace +
		` WHERE n.identyfikator_zewnetrzny = ?`

	listaNotatekWorkspace = `SELECT ` + kolumnyNotatkiWorkspace + zrodloNotatkiWorkspace +
		` WHERE n.projekt_id = ? ORDER BY n.tytul, n.id`

	usunNotatkeWorkspace = `DELETE FROM notatka_projektu WHERE identyfikator_zewnetrzny = ?`

	przeniesPodrzedneWorkspace = `UPDATE notatka_projektu SET notatka_nadrzedna = ?
	                              WHERE notatka_nadrzedna = ?`

	kasujOdnosnikiWorkspace = `DELETE FROM odnosnik_notatki_projektu WHERE notatka_zrodlowa = ?`

	wstawOdnosnikWorkspace = `INSERT INTO odnosnik_notatki_projektu
	    (projekt_id, notatka_zrodlowa, nazwa_docelowa, notatka_docelowa, rodzaj, kontekst)
	    VALUES (?, ?, ?, ?, ?, ?)`

	listaOdnosnikowWorkspace = `SELECT projekt_id, notatka_zrodlowa, nazwa_docelowa,
	                            notatka_docelowa, rodzaj, kontekst
	                            FROM odnosnik_notatki_projektu WHERE projekt_id = ? ORDER BY id`

	// Odnośnik wskazujący stronę jeszcze niezałożoną wiąże się z nią nazwą,
	// dlatego wskazanie po założeniu strony domyka się jednym poleceniem.
	domknijOdnosnikiWorkspace = `UPDATE odnosnik_notatki_projektu SET notatka_docelowa = ?
	                             WHERE projekt_id = ? AND notatka_docelowa = ''
	                               AND lower(nazwa_docelowa) = lower(?)`

	// Usunięcie strony nie kasuje odnośników do niej: zostają jako wskazania
	// nazwy bez strony, czyli w stanie, w którym wiki je zna.
	osierocOdnosnikiWorkspace = `UPDATE odnosnik_notatki_projektu SET notatka_docelowa = ''
	                             WHERE notatka_docelowa = ?`
)

// ZapiszNotatkeWorkspace zapisuje notatkę wraz z kompletem jej odnośników
// w jednej transakcji i oddaje stan po zapisie.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszNotatkeWorkspace(ctx context.Context,
	notatka NotatkaWorkspace, odnosniki []OdnosnikWorkspace) (NotatkaWorkspace, error) {

	if notatka.Identyfikator == "" {
		return NotatkaWorkspace{}, fmt.Errorf("dane: notatka projektu bez identyfikatora")
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		_, err := transakcja.ExecContext(ctx, zapiszNotatkeWorkspace, notatka.Identyfikator,
			notatka.ProjektID, notatka.Tytul, notatka.Tresc, notatka.NotatkaNadrzedna,
			strings.Join(notatka.Etykiety, "\n"), strings.Join(notatka.Naglowki, "\n"),
			string(rodzajAutoraWorkspace(notatka.RodzajAutora)), notatka.Autor)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać notatki %q: %w", notatka.Identyfikator, err)
		}
		if _, err := transakcja.ExecContext(ctx, kasujOdnosnikiWorkspace, notatka.Identyfikator); err != nil {
			return fmt.Errorf("dane: nie można wymienić odnośników notatki %q: %w",
				notatka.Identyfikator, err)
		}
		for _, odnosnik := range odnosniki {
			_, err := transakcja.ExecContext(ctx, wstawOdnosnikWorkspace, notatka.ProjektID,
				notatka.Identyfikator, odnosnik.NazwaDocelowa, odnosnik.NotatkaDocelowa,
				string(rodzajKrawedziWorkspace(odnosnik.Rodzaj)), odnosnik.Kontekst)
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać odnośnika notatki %q: %w",
					notatka.Identyfikator, err)
			}
		}
		// Strona właśnie zapisana domyka odnośniki, które czekały na jej nazwę.
		_, err = transakcja.ExecContext(ctx, domknijOdnosnikiWorkspace, notatka.Identyfikator,
			notatka.ProjektID, notatka.Tytul)
		if err != nil {
			return fmt.Errorf("dane: nie można domknąć odnośników do strony %q: %w",
				notatka.Tytul, err)
		}
		return nil
	})
	if err != nil {
		return NotatkaWorkspace{}, err
	}
	return r.NotatkaWorkspace(ctx, notatka.Identyfikator)
}

// NotatkaWorkspace zwraca z tabeli `notatka_projektu` jedną notatkę projektu
// wraz z jej pełną treścią wpisu.
func (r *repozytoriumPrzestrzeniRoboczej) NotatkaWorkspace(ctx context.Context,
	identyfikator string) (NotatkaWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzNotatkeWorkspace)
	if err != nil {
		return NotatkaWorkspace{}, err
	}
	notatka, err := odczytajNotatkeWorkspace(polecenie.QueryRowContext(ctx, identyfikator))
	if errors.Is(err, sql.ErrNoRows) {
		return NotatkaWorkspace{}, ErrBrakWiersza
	}
	if err != nil {
		return NotatkaWorkspace{}, fmt.Errorf("dane: nieczytelna notatka %q: %w", identyfikator, err)
	}
	return notatka, nil
}

// NotatkiWorkspace zwraca z tabeli `notatka_projektu` komplet notatek i stron
// wiki całego projektu, bez zawężenia.
func (r *repozytoriumPrzestrzeniRoboczej) NotatkiWorkspace(ctx context.Context,
	projektID int64) ([]NotatkaWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaNotatekWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać notatek projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []NotatkaWorkspace{}
	for wiersze.Next() {
		notatka, err := odczytajNotatkeWorkspace(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz notatki projektu: %w", err)
		}
		lista = append(lista, notatka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt notatek projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// UsunNotatkiWorkspace usuwa wskazane notatki. Strony podrzędne wskazane
// osobno nie są: rdzeń podaje pełny wykaz albo przenosi je pod stronę nadrzędną
// usuwanej — o tym rozstrzyga żądanie, nie repozytorium.
func (r *repozytoriumPrzestrzeniRoboczej) UsunNotatkiWorkspace(ctx context.Context,
	identyfikatory []string, nadrzednaDlaSierot string) (int, error) {

	usuniete := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		for _, identyfikator := range identyfikatory {
			_, err := transakcja.ExecContext(ctx, przeniesPodrzedneWorkspace,
				nadrzednaDlaSierot, identyfikator)
			if err != nil {
				return fmt.Errorf("dane: nie można przenieść stron podrzędnych %q: %w",
					identyfikator, err)
			}
			wynik, err := transakcja.ExecContext(ctx, usunNotatkeWorkspace, identyfikator)
			if err != nil {
				return fmt.Errorf("dane: nie można usunąć notatki %q: %w", identyfikator, err)
			}
			zmienione, err := wynik.RowsAffected()
			if err != nil {
				return fmt.Errorf("dane: nieznany wynik usunięcia notatki %q: %w", identyfikator, err)
			}
			usuniete += int(zmienione)
			if _, err := transakcja.ExecContext(ctx, kasujOdnosnikiWorkspace, identyfikator); err != nil {
				return fmt.Errorf("dane: nie można usunąć odnośników notatki %q: %w",
					identyfikator, err)
			}
			if _, err := transakcja.ExecContext(ctx, osierocOdnosnikiWorkspace, identyfikator); err != nil {
				return fmt.Errorf("dane: nie można osierocić odnośników do notatki %q: %w",
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

// OdnosnikiWorkspace zwraca komplet odnośników treści projektu — materiał
// panelu „co linkuje tutaj" oraz grafu wiedzy.
func (r *repozytoriumPrzestrzeniRoboczej) OdnosnikiWorkspace(ctx context.Context,
	projektID int64) ([]OdnosnikWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaOdnosnikowWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać odnośników projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []OdnosnikWorkspace{}
	for wiersze.Next() {
		var odnosnik OdnosnikWorkspace
		var rodzaj string
		err := wiersze.Scan(&odnosnik.ProjektID, &odnosnik.NotatkaZrodlowa,
			&odnosnik.NazwaDocelowa, &odnosnik.NotatkaDocelowa, &rodzaj, &odnosnik.Kontekst)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz odnośnika notatki: %w", err)
		}
		odnosnik.Rodzaj = shared.WorkspaceGraphEdgeKind(rodzaj)
		lista = append(lista, odnosnik)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt odnośników projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// odczytajNotatkeWorkspace składa strukturę NotatkaWorkspace z jednego
// wiersza wyniku zapytania do bazy.
func odczytajNotatkeWorkspace(wiersz skaner) (NotatkaWorkspace, error) {
	var notatka NotatkaWorkspace
	var etykiety, naglowki, rodzajAutora string
	err := wiersz.Scan(&notatka.ID, &notatka.ProjektKod, &notatka.ProjektID,
		&notatka.Identyfikator, &notatka.Tytul, &notatka.Tresc, &notatka.NotatkaNadrzedna,
		&etykiety, &naglowki, &rodzajAutora, &notatka.Autor,
		&notatka.Utworzono, &notatka.Zaktualizowano)
	if err != nil {
		return NotatkaWorkspace{}, err
	}
	notatka.Etykiety = rozdzielWierszamiWorkspace(etykiety)
	notatka.Naglowki = rozdzielWierszamiWorkspace(naglowki)
	notatka.RodzajAutora = shared.WorkspaceAssigneeKind(rodzajAutora)
	return notatka, nil
}

// rodzajAutoraWorkspace sprowadza rodzaj autora do wartości znanej kolumnie.
// Pusty rodzaj znaczy zapis Operatora — to jest stan domyślny, nie brak wiedzy.
func rodzajAutoraWorkspace(rodzaj shared.WorkspaceAssigneeKind) shared.WorkspaceAssigneeKind {
	if rodzaj == shared.WorkspaceAssigneeKindAgent {
		return shared.WorkspaceAssigneeKindAgent
	}
	return shared.WorkspaceAssigneeKindOperator
}

// rodzajKrawedziWorkspace sprowadza rodzaj odnośnika do wartości znanej
// kolumnie; brak wskazania znaczy odnośnik wiki.
func rodzajKrawedziWorkspace(rodzaj shared.WorkspaceGraphEdgeKind) shared.WorkspaceGraphEdgeKind {
	switch rodzaj {
	case shared.WorkspaceGraphEdgeKindReference, shared.WorkspaceGraphEdgeKindEmbed,
		shared.WorkspaceGraphEdgeKindAttachment:
		return rodzaj
	}
	return shared.WorkspaceGraphEdgeKindWikilink
}
