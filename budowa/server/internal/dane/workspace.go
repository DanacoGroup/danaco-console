// Plik prowadzi obszar projektu przestrzeni roboczej w tabeli projekt oraz
// odczyt i odpięcie kart sesji pracujących w projekcie; pamięć projektu
// i przypisania ekspertów leżą w osobnych plikach tego repozytorium.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// Projekt to wiersz tabeli `projekt`. Kod jest identyfikatorem, którym projekt
// wychodzi kontraktem i którym posługuje się kolumna `sesja.projekt`.
type Projekt struct {
	ID             int64
	Kod            string
	Nazwa          string
	Opis           *string
	Stan           shared.WorkspaceProjectStatus
	Utworzono      string
	Zaktualizowano string
}

// RepozytoriumPrzestrzeniRoboczej jest kontraktem obszaru Workspace: projekty,
// pamięć, zadania, notatki, materiały i ślad zdarzeń.
type RepozytoriumPrzestrzeniRoboczej interface {
	ZapewnijProjekt(ctx context.Context, kod, nazwa string) (Projekt, bool, error)
	ZalozProjekt(ctx context.Context, kod, nazwa string, opis *string) (Projekt, error)
	PrzemianujProjekt(ctx context.Context, kod, nazwa string) (Projekt, error)
	UsunProjekt(ctx context.Context, kod string) ([]string, error)
	Projekt(ctx context.Context, kod string) (Projekt, error)
	Projekty(ctx context.Context, zArchiwalnymi bool) ([]Projekt, error)
	OdnotujCzynnosc(ctx context.Context, projektID int64) error
	SesjeProjektu(ctx context.Context, kod string) ([]string, error)

	ZapiszWpisPamieci(ctx context.Context, wpis WpisPamieciProjektu) (WpisPamieciProjektu, error)
	WpisyPamieci(ctx context.Context, projektID int64, limit int) ([]WpisPamieciProjektu, error)
	WpisyPamieciWspoldzielone(ctx context.Context, projektID int64, limit int) ([]WpisPamieciProjektu, error)

	PrzypiszAgenta(ctx context.Context, przypisanie PrzypisanieAgenta) (PrzypisanieAgenta, error)
	PrzypisaniaAgentow(ctx context.Context, projektID int64) ([]PrzypisanieAgenta, error)
	OdlaczAgentaWorkspace(ctx context.Context, projektID int64, agent string) (bool, error)
	UstawStanProjektuWorkspace(ctx context.Context, projektID int64,
		stan shared.WorkspaceProjectStatus) error

	// Hub planowania: zadania, ich zależności i kolumny tablicy kanban.
	ZapiszZadanieWorkspace(ctx context.Context, zadanie ZadanieWorkspace) (ZadanieWorkspace, error)
	ZadanieWorkspace(ctx context.Context, identyfikator string) (ZadanieWorkspace, error)
	ZadaniaWorkspace(ctx context.Context, projektID int64) ([]ZadanieWorkspace, error)
	UsunZadaniaWorkspace(ctx context.Context, identyfikatory []string) (int, error)
	ZapiszKolumneTablicyWorkspace(ctx context.Context, kolumna KolumnaTablicyWorkspace) error
	KolumnyTablicyWorkspace(ctx context.Context, projektID int64) ([]KolumnaTablicyWorkspace, error)
	ZapiszZaleznoscWorkspace(ctx context.Context, zaleznosc ZaleznoscWorkspace) (ZaleznoscWorkspace, error)
	ZaleznoscWorkspace(ctx context.Context, identyfikator string) (ZaleznoscWorkspace, error)
	ZaleznosciWorkspace(ctx context.Context, projektID int64) ([]ZaleznoscWorkspace, error)
	UsunZaleznoscWorkspace(ctx context.Context, identyfikator string) (bool, error)

	// Notatki, strony wiki i ich odnośniki treści.
	ZapiszNotatkeWorkspace(ctx context.Context, notatka NotatkaWorkspace,
		odnosniki []OdnosnikWorkspace) (NotatkaWorkspace, error)
	NotatkaWorkspace(ctx context.Context, identyfikator string) (NotatkaWorkspace, error)
	NotatkiWorkspace(ctx context.Context, projektID int64) ([]NotatkaWorkspace, error)
	UsunNotatkiWorkspace(ctx context.Context, identyfikatory []string,
		nadrzednaDlaSierot string) (int, error)
	OdnosnikiWorkspace(ctx context.Context, projektID int64) ([]OdnosnikWorkspace, error)

	// Materiał projektu: tablica wizualna, wskaźnik treści plików, kalendarz.
	ZapiszTabliceWorkspace(ctx context.Context, tablica TablicaWorkspace) (TablicaWorkspace, error)
	TablicaWorkspace(ctx context.Context, identyfikator string) (TablicaWorkspace, error)
	TabliceWorkspace(ctx context.Context, projektID int64) ([]TablicaWorkspace, error)
	ZapiszWyciagWorkspace(ctx context.Context, wyciag WyciagWorkspace) error
	WyciagWorkspace(ctx context.Context, projektID int64, plik string) (WyciagWorkspace, error)
	WyciagiWorkspace(ctx context.Context, projektID int64) ([]WyciagWorkspace, error)
	UsunWyciagWorkspace(ctx context.Context, projektID int64, plik string) error
	ZapiszPozycjeKalendarzaWorkspace(ctx context.Context, pozycja PozycjaKalendarzaWorkspace) error
	PozycjeKalendarzaWorkspace(ctx context.Context, projektID int64) ([]PozycjaKalendarzaWorkspace, error)

	// Ślad: oś czasu, komentarze i historia instrukcji.
	ZapiszZdarzenieWorkspace(ctx context.Context, zdarzenie ZdarzenieWorkspace) error
	ZdarzeniaWorkspace(ctx context.Context, projektID int64) ([]ZdarzenieWorkspace, error)
	ZapiszKomentarzWorkspace(ctx context.Context, komentarz KomentarzWorkspace) (KomentarzWorkspace, error)
	KomentarzWorkspace(ctx context.Context, identyfikator string) (KomentarzWorkspace, error)
	KomentarzeWorkspace(ctx context.Context, projektID int64) ([]KomentarzWorkspace, error)
	UsunKomentarzeWorkspace(ctx context.Context, identyfikatory []string) (int, error)
	ZapiszWersjeInstrukcjiWorkspace(ctx context.Context,
		wersja WersjaInstrukcjiWorkspace) (WersjaInstrukcjiWorkspace, error)
	WersjaInstrukcjiWorkspace(ctx context.Context, identyfikator string) (WersjaInstrukcjiWorkspace, error)
	WersjeInstrukcjiWorkspace(ctx context.Context, projektID int64) ([]WersjaInstrukcjiWorkspace, error)
}

const (
	kolumnyProjektu = `id, kod, nazwa, opis, stan, utworzono, zaktualizowano`

	wstawProjekt = `INSERT INTO projekt (kod, nazwa) VALUES (?, ?)
	                ON CONFLICT(kod) DO NOTHING`

	// Założenie projektu nie znosi konfliktu kodu: kod nadaje rdzeń, więc wiersz
	// zastany znaczy zderzenie identyfikatorów, a nie powtórzone żądanie.
	zalozProjekt = `INSERT INTO projekt (kod, nazwa, opis) VALUES (?, ?, ?)`

	przemianujProjekt = `UPDATE projekt
	                     SET nazwa = ?,
	                         zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                     WHERE kod = ?`

	usunProjekt = `DELETE FROM projekt WHERE kod = ?`

	// Sesje projektu nie są jego wierszami potomnymi — kolumna `sesja.projekt`
	// niesie sam kod, bez więzi obcej. Usunięcie projektu musi je odpiąć jawnie,
	// inaczej wskazywałyby w pustkę.
	odepnijSesjeProjektu = `UPDATE sesja
	                        SET projekt = NULL,
	                            zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                        WHERE projekt = ?`

	pobierzProjekt = `SELECT ` + kolumnyProjektu + ` FROM projekt WHERE kod = ?`

	listaProjektow = `SELECT ` + kolumnyProjektu + ` FROM projekt
	                  WHERE (? = 1 OR stan <> 'archiwalny')
	                  ORDER BY zaktualizowano DESC, id DESC`

	odnotujCzynnoscProjektu = `UPDATE projekt
	                           SET zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                           WHERE id = ?`

	// Karty sesji projektu pochodzą z kolumny `sesja.projekt` — jedynego miejsca,
	// w którym sesja mówi, nad czym pracuje. Sesja bez identyfikatora zewnętrznego
	// nie wyszła nigdy kontraktem, więc do wykazu nie wchodzi.
	sesjeProjektu = `SELECT identyfikator_zewnetrzny FROM sesja
	                 WHERE projekt = ? AND identyfikator_zewnetrzny IS NOT NULL
	                 ORDER BY zaktualizowano DESC`
)

type repozytoriumPrzestrzeniRoboczej struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumPrzestrzeniRoboczej(z *zapytania, db *sql.DB) *repozytoriumPrzestrzeniRoboczej {
	return &repozytoriumPrzestrzeniRoboczej{zapytania: z, db: db}
}

// ZapewnijProjekt zwraca projekt o wskazanym kodzie, zakładając go, gdy nie ma
// jeszcze wiersza. Drugi wynik mówi, czy projekt powstał teraz — rdzeń rozgłasza
// wtedy `workspace.project.changed` z rodzajem `created`.
func (r *repozytoriumPrzestrzeniRoboczej) ZapewnijProjekt(ctx context.Context,
	kod, nazwa string) (Projekt, bool, error) {

	if kod == "" {
		return Projekt{}, false, fmt.Errorf("dane: projekt bez identyfikatora")
	}
	if nazwa == "" {
		nazwa = kod
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawProjekt)
	if err != nil {
		return Projekt{}, false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, nazwa)
	if err != nil {
		return Projekt{}, false, fmt.Errorf("dane: nie można założyć projektu %q: %w", kod, err)
	}
	wstawione, err := wynik.RowsAffected()
	if err != nil {
		return Projekt{}, false, fmt.Errorf("dane: nieznany wynik założenia projektu %q: %w", kod, err)
	}
	projekt, err := r.Projekt(ctx, kod)
	if err != nil {
		return Projekt{}, false, err
	}
	return projekt, wstawione > 0, nil
}

// ZalozProjekt zakłada projekt o nadanym kodzie i zwraca go w kształcie, jaki
// wychodzi kontraktem. Kod zajęty kończy się błędem: zwrócenie wiersza zastanego
// pokazałoby Operatorowi cudzy projekt jako właśnie założony.
func (r *repozytoriumPrzestrzeniRoboczej) ZalozProjekt(ctx context.Context,
	kod, nazwa string, opis *string) (Projekt, error) {

	if kod == "" {
		return Projekt{}, fmt.Errorf("dane: projekt bez identyfikatora")
	}
	if nazwa == "" {
		return Projekt{}, fmt.Errorf("dane: projekt %q bez nazwy", kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zalozProjekt)
	if err != nil {
		return Projekt{}, err
	}
	if _, err := polecenie.ExecContext(ctx, kod, nazwa, tekstDoKolumny(opis)); err != nil {
		return Projekt{}, fmt.Errorf("dane: nie można założyć projektu %q: %w", kod, err)
	}
	return r.Projekt(ctx, kod)
}

// PrzemianujProjekt zmienia nazwę projektu, zostawiając kod bez zmiany: kodem
// wiążą się sesje i wiersze modułu, więc przemianowanie nie ma prawa ich zerwać.
// Brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumPrzestrzeniRoboczej) PrzemianujProjekt(ctx context.Context,
	kod, nazwa string) (Projekt, error) {

	if nazwa == "" {
		return Projekt{}, fmt.Errorf("dane: nazwa projektu %q nie może być pusta", kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, przemianujProjekt)
	if err != nil {
		return Projekt{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, nazwa, kod)
	if err != nil {
		return Projekt{}, fmt.Errorf("dane: nie można zmienić nazwy projektu %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return Projekt{}, fmt.Errorf("dane: nieznany wynik zmiany nazwy projektu %q: %w", kod, err)
	}
	if zmienione == 0 {
		return Projekt{}, ErrBrakWiersza
	}
	return r.Projekt(ctx, kod)
}

// UsunProjekt kasuje projekt wraz z jego wierszami potomnymi i zwraca karty
// sesji, które straciły przypisanie.
//
// Sesje zostają w historii — usunięcie projektu nie jest usunięciem rozmów,
// które w nim powstały. Odpięcie i skasowanie idą jedną transakcją, bo projekt
// zdjęty przy sesjach wskazujących na niego zostawiłby wykaz sesji z kodem bez
// pokrycia. Wykaz zwracany obejmuje sesje mające identyfikator zewnętrzny —
// tylko one wyszły kiedykolwiek kontraktem. Brak wiersza wraca jako
// ErrBrakWiersza.
func (r *repozytoriumPrzestrzeniRoboczej) UsunProjekt(ctx context.Context, kod string) ([]string, error) {
	if kod == "" {
		return nil, fmt.Errorf("dane: projekt bez identyfikatora")
	}
	odpiete := []string{}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, sesjeProjektu)
		if err != nil {
			return err
		}
		wiersze, err := odczyt.QueryContext(ctx, kod)
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać sesji projektu %q: %w", kod, err)
		}
		for wiersze.Next() {
			var identyfikator string
			if err := wiersze.Scan(&identyfikator); err != nil {
				_ = wiersze.Close()
				return fmt.Errorf("dane: nieczytelny wiersz sesji projektu: %w", err)
			}
			odpiete = append(odpiete, identyfikator)
		}
		if err := wiersze.Err(); err != nil {
			_ = wiersze.Close()
			return fmt.Errorf("dane: przerwany odczyt sesji projektu %q: %w", kod, err)
		}
		if err := wiersze.Close(); err != nil {
			return fmt.Errorf("dane: przerwany odczyt sesji projektu %q: %w", kod, err)
		}

		odpiecie, err := r.zapytania.wTransakcji(ctx, transakcja, odepnijSesjeProjektu)
		if err != nil {
			return err
		}
		if _, err := odpiecie.ExecContext(ctx, kod); err != nil {
			return fmt.Errorf("dane: nie można odpiąć sesji projektu %q: %w", kod, err)
		}

		kasowanie, err := r.zapytania.wTransakcji(ctx, transakcja, usunProjekt)
		if err != nil {
			return err
		}
		wynik, err := kasowanie.ExecContext(ctx, kod)
		if err != nil {
			return fmt.Errorf("dane: nie można usunąć projektu %q: %w", kod, err)
		}
		skasowane, err := wynik.RowsAffected()
		if err != nil {
			return fmt.Errorf("dane: nieznany wynik usunięcia projektu %q: %w", kod, err)
		}
		if skasowane == 0 {
			return ErrBrakWiersza
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return odpiete, nil
}

// Projekt zwraca projekt o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się nie powiódł”.
func (r *repozytoriumPrzestrzeniRoboczej) Projekt(ctx context.Context, kod string) (Projekt, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzProjekt)
	if err != nil {
		return Projekt{}, err
	}
	projekt, err := odczytajProjekt(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return Projekt{}, ErrBrakWiersza
	}
	if err != nil {
		return Projekt{}, fmt.Errorf("dane: nieczytelny wiersz projektu %q: %w", kod, err)
	}
	return projekt, nil
}

// Projekty zwraca projekty konta od najświeższego. Lewy panel ramy pokazuje
// całość dorobku Operatora, więc wykaz obejmuje także projekty bez sesji.
func (r *repozytoriumPrzestrzeniRoboczej) Projekty(ctx context.Context,
	zArchiwalnymi bool) ([]Projekt, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaProjektow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, liczbaLogiczna(zArchiwalnymi))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu projektów: %w", err)
	}
	defer func() { _ = wiersze.Close() }()
	wykaz := make([]Projekt, 0, 16)
	for wiersze.Next() {
		projekt, err := odczytajProjekt(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wykazu projektów: %w", err)
		}
		wykaz = append(wykaz, projekt)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wykazu projektów: %w", err)
	}
	return wykaz, nil
}

// OdnotujCzynnosc przesuwa znacznik ostatniej zmiany projektu. Zestawienie
// Project Dashboard bierze z niego czas ostatniej czynności, więc znacznik
// przesuwa każdy zapis w obrębie projektu.
func (r *repozytoriumPrzestrzeniRoboczej) OdnotujCzynnosc(ctx context.Context, projektID int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, odnotujCzynnoscProjektu)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, projektID); err != nil {
		return fmt.Errorf("dane: nie można odnotować czynności projektu %d: %w", projektID, err)
	}
	return nil
}

// SesjeProjektu zwraca identyfikatory kart sesji pracujących w projekcie,
// od ostatnio zmienionej karty.
func (r *repozytoriumPrzestrzeniRoboczej) SesjeProjektu(ctx context.Context, kod string) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, sesjeProjektu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kod)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać sesji projektu %q: %w", kod, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var identyfikator string
		if err := wiersze.Scan(&identyfikator); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz sesji projektu: %w", err)
		}
		lista = append(lista, identyfikator)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt sesji projektu %q: %w", kod, err)
	}
	return lista, nil
}

// odczytajProjekt składa strukturę projektu z jednego wiersza wyniku
// zapytania, tłumacząc opis i stan.
func odczytajProjekt(wiersz skaner) (Projekt, error) {
	var projekt Projekt
	var opis sql.NullString
	var stan string
	err := wiersz.Scan(&projekt.ID, &projekt.Kod, &projekt.Nazwa, &opis, &stan,
		&projekt.Utworzono, &projekt.Zaktualizowano)
	if err != nil {
		return Projekt{}, err
	}
	projekt.Opis = tekstZKolumny(opis)
	projekt.Stan = shared.WorkspaceProjectStatus(stan)
	return projekt, nil
}
