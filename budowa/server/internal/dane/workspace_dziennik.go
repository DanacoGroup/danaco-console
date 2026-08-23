// Odpowiedzialność pliku: zapis zdarzeń osi czasu projektu (tabela
// `zdarzenie_projektu`), komentarze przy bytach projektu (tabela
// `komentarz_projektu`) oraz historia instrukcji systemowych (tabela
// `wersja_instrukcji_projektu`).
//
// Trzy tabele, jeden plik, bo wszystkie trzy są ZAPISEM TEGO, CO SIĘ WYDARZYŁO,
// a nie stanem bieżącym. Stan bieżący instrukcji leży w tabeli `ustawienie`,
// stan bieżący zadania w `zadanie_projektu` — tutaj leży ślad.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// ZdarzenieWorkspace to wiersz osi czasu projektu.
type ZdarzenieWorkspace struct {
	ProjektID     int64
	Identyfikator string
	Rodzaj        shared.WorkspaceActivityKind
	Zmiana        shared.ChangeKind
	RodzajBytu    shared.WorkspaceEntityKind
	Byt           string
	Opis          string
	RodzajAutora  shared.WorkspaceAssigneeKind
	Autor         string
	Zaszlo        string
}

// KomentarzWorkspace to wiersz komentarza przy bycie projektu.
type KomentarzWorkspace struct {
	ProjektID          int64
	ProjektKod         string
	Identyfikator      string
	RodzajBytu         shared.WorkspaceEntityKind
	Byt                string
	Tresc              string
	KomentarzNadrzedny string
	RodzajAutora       shared.WorkspaceAssigneeKind
	Autor              string
	Przywolania        []string
	Utworzono          string
	Zaktualizowano     string
}

// WersjaInstrukcjiWorkspace to wiersz historii instrukcji systemowych projektu.
type WersjaInstrukcjiWorkspace struct {
	ProjektID     int64
	ProjektKod    string
	Identyfikator string
	Tresc         string
	Odcisk        string
	Poziom        shared.ConfigScope
	KluczZasiegu  string
	RodzajAutora  shared.WorkspaceAssigneeKind
	Autor         string
	PrzywroconoZ  string
	Utworzono     string
}

const (
	wstawZdarzenieWorkspace = `INSERT INTO zdarzenie_projektu
	    (identyfikator_zewnetrzny, projekt_id, rodzaj, zmiana, rodzaj_bytu, byt, opis,
	     rodzaj_autora, autor)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	listaZdarzenWorkspace = `SELECT projekt_id, identyfikator_zewnetrzny, rodzaj, zmiana,
	                         rodzaj_bytu, byt, opis, rodzaj_autora, autor, zaszlo
	                         FROM zdarzenie_projektu WHERE projekt_id = ? ORDER BY id DESC`

	kolumnyKomentarzaWorkspace = `k.projekt_id, p.kod, k.identyfikator_zewnetrzny, k.rodzaj_bytu,
	                              k.byt, k.tresc, k.komentarz_nadrzedny, k.rodzaj_autora, k.autor,
	                              k.przywolania, k.utworzono, k.zaktualizowano`

	zrodloKomentarzaWorkspace = ` FROM komentarz_projektu k JOIN projekt p ON p.id = k.projekt_id`

	wstawKomentarzWorkspace = `INSERT INTO komentarz_projektu
	    (identyfikator_zewnetrzny, projekt_id, rodzaj_bytu, byt, tresc, komentarz_nadrzedny,
	     rodzaj_autora, autor, przywolania)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzKomentarzWorkspace = `SELECT ` + kolumnyKomentarzaWorkspace + zrodloKomentarzaWorkspace +
		` WHERE k.identyfikator_zewnetrzny = ?`

	listaKomentarzyWorkspace = `SELECT ` + kolumnyKomentarzaWorkspace + zrodloKomentarzaWorkspace +
		` WHERE k.projekt_id = ? ORDER BY k.id`

	usunKomentarzWorkspace = `DELETE FROM komentarz_projektu WHERE identyfikator_zewnetrzny = ?`

	kolumnyWersjiInstrukcjiWorkspace = `w.projekt_id, p.kod, w.identyfikator_zewnetrzny, w.tresc,
	                                    w.odcisk, w.poziom, w.klucz_zasiegu, w.rodzaj_autora,
	                                    w.autor, w.przywrocono_z, w.utworzono`

	zrodloWersjiInstrukcjiWorkspace = ` FROM wersja_instrukcji_projektu w
	                                    JOIN projekt p ON p.id = w.projekt_id`

	wstawWersjeInstrukcjiWorkspace = `INSERT INTO wersja_instrukcji_projektu
	    (identyfikator_zewnetrzny, projekt_id, tresc, odcisk, poziom, klucz_zasiegu,
	     rodzaj_autora, autor, przywrocono_z)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzWersjeInstrukcjiWorkspace = `SELECT ` + kolumnyWersjiInstrukcjiWorkspace +
		zrodloWersjiInstrukcjiWorkspace + ` WHERE w.identyfikator_zewnetrzny = ?`

	listaWersjiInstrukcjiWorkspace = `SELECT ` + kolumnyWersjiInstrukcjiWorkspace +
		zrodloWersjiInstrukcjiWorkspace + ` WHERE w.projekt_id = ? ORDER BY w.id DESC`
)

// ZapiszZdarzenieWorkspace odkłada zdarzenie osi czasu projektu.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszZdarzenieWorkspace(ctx context.Context,
	zdarzenie ZdarzenieWorkspace) error {

	polecenie, err := r.zapytania.przygotuj(ctx, wstawZdarzenieWorkspace)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, zdarzenie.Identyfikator, zdarzenie.ProjektID,
		string(zdarzenie.Rodzaj), string(zdarzenie.Zmiana), string(zdarzenie.RodzajBytu),
		zdarzenie.Byt, zdarzenie.Opis, string(rodzajAutoraWorkspace(zdarzenie.RodzajAutora)),
		zdarzenie.Autor)
	if err != nil {
		return fmt.Errorf("dane: nie można odłożyć zdarzenia projektu %d: %w",
			zdarzenie.ProjektID, err)
	}
	return nil
}

// ZdarzeniaWorkspace zwraca oś czasu projektu od zdarzenia najnowszego.
func (r *repozytoriumPrzestrzeniRoboczej) ZdarzeniaWorkspace(ctx context.Context,
	projektID int64) ([]ZdarzenieWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZdarzenWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać osi czasu projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []ZdarzenieWorkspace{}
	for wiersze.Next() {
		var zdarzenie ZdarzenieWorkspace
		var rodzaj, zmiana, rodzajBytu, rodzajAutora string
		err := wiersze.Scan(&zdarzenie.ProjektID, &zdarzenie.Identyfikator, &rodzaj, &zmiana,
			&rodzajBytu, &zdarzenie.Byt, &zdarzenie.Opis, &rodzajAutora, &zdarzenie.Autor,
			&zdarzenie.Zaszlo)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz osi czasu projektu: %w", err)
		}
		zdarzenie.Rodzaj = shared.WorkspaceActivityKind(rodzaj)
		zdarzenie.Zmiana = shared.ChangeKind(zmiana)
		zdarzenie.RodzajBytu = shared.WorkspaceEntityKind(rodzajBytu)
		zdarzenie.RodzajAutora = shared.WorkspaceAssigneeKind(rodzajAutora)
		lista = append(lista, zdarzenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt osi czasu projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// ZapiszKomentarzWorkspace zakłada komentarz i oddaje jego stan po zapisie.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszKomentarzWorkspace(ctx context.Context,
	komentarz KomentarzWorkspace) (KomentarzWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wstawKomentarzWorkspace)
	if err != nil {
		return KomentarzWorkspace{}, err
	}
	_, err = polecenie.ExecContext(ctx, komentarz.Identyfikator, komentarz.ProjektID,
		string(komentarz.RodzajBytu), komentarz.Byt, komentarz.Tresc,
		komentarz.KomentarzNadrzedny, string(rodzajAutoraWorkspace(komentarz.RodzajAutora)),
		komentarz.Autor, strings.Join(komentarz.Przywolania, "\n"))
	if err != nil {
		return KomentarzWorkspace{}, fmt.Errorf("dane: nie można zapisać komentarza %q: %w",
			komentarz.Identyfikator, err)
	}
	return r.KomentarzWorkspace(ctx, komentarz.Identyfikator)
}

// KomentarzWorkspace zwraca jeden komentarz po identyfikatorze.
func (r *repozytoriumPrzestrzeniRoboczej) KomentarzWorkspace(ctx context.Context,
	identyfikator string) (KomentarzWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKomentarzWorkspace)
	if err != nil {
		return KomentarzWorkspace{}, err
	}
	komentarz, err := odczytajKomentarzWorkspace(polecenie.QueryRowContext(ctx, identyfikator))
	if errors.Is(err, sql.ErrNoRows) {
		return KomentarzWorkspace{}, ErrBrakWiersza
	}
	if err != nil {
		return KomentarzWorkspace{}, fmt.Errorf("dane: nieczytelny komentarz %q: %w",
			identyfikator, err)
	}
	return komentarz, nil
}

// KomentarzeWorkspace zwraca komplet komentarzy projektu w kolejności zapisu.
func (r *repozytoriumPrzestrzeniRoboczej) KomentarzeWorkspace(ctx context.Context,
	projektID int64) ([]KomentarzWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaKomentarzyWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać komentarzy projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []KomentarzWorkspace{}
	for wiersze.Next() {
		komentarz, err := odczytajKomentarzWorkspace(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz komentarza projektu: %w", err)
		}
		lista = append(lista, komentarz)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt komentarzy projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// UsunKomentarzeWorkspace usuwa wskazane komentarze i oddaje liczbę naprawdę
// usuniętych.
func (r *repozytoriumPrzestrzeniRoboczej) UsunKomentarzeWorkspace(ctx context.Context,
	identyfikatory []string) (int, error) {

	usuniete := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		for _, identyfikator := range identyfikatory {
			wynik, err := transakcja.ExecContext(ctx, usunKomentarzWorkspace, identyfikator)
			if err != nil {
				return fmt.Errorf("dane: nie można usunąć komentarza %q: %w", identyfikator, err)
			}
			zmienione, err := wynik.RowsAffected()
			if err != nil {
				return fmt.Errorf("dane: nieznany wynik usunięcia komentarza %q: %w",
					identyfikator, err)
			}
			usuniete += int(zmienione)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return usuniete, nil
}

// ZapiszWersjeInstrukcjiWorkspace odkłada wersję instrukcji systemowych.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszWersjeInstrukcjiWorkspace(ctx context.Context,
	wersja WersjaInstrukcjiWorkspace) (WersjaInstrukcjiWorkspace, error) {

	poziom, err := poziomZasieguNaBaze(wersja.Poziom)
	if err != nil {
		return WersjaInstrukcjiWorkspace{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawWersjeInstrukcjiWorkspace)
	if err != nil {
		return WersjaInstrukcjiWorkspace{}, err
	}
	_, err = polecenie.ExecContext(ctx, wersja.Identyfikator, wersja.ProjektID, wersja.Tresc,
		wersja.Odcisk, poziom, wersja.KluczZasiegu,
		string(rodzajAutoraWorkspace(wersja.RodzajAutora)), wersja.Autor, wersja.PrzywroconoZ)
	if err != nil {
		return WersjaInstrukcjiWorkspace{}, fmt.Errorf(
			"dane: nie można odłożyć wersji instrukcji projektu %d: %w", wersja.ProjektID, err)
	}
	return r.WersjaInstrukcjiWorkspace(ctx, wersja.Identyfikator)
}

// WersjaInstrukcjiWorkspace zwraca jedną wersję instrukcji po identyfikatorze.
func (r *repozytoriumPrzestrzeniRoboczej) WersjaInstrukcjiWorkspace(ctx context.Context,
	identyfikator string) (WersjaInstrukcjiWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWersjeInstrukcjiWorkspace)
	if err != nil {
		return WersjaInstrukcjiWorkspace{}, err
	}
	wersja, err := odczytajWersjeInstrukcjiWorkspace(polecenie.QueryRowContext(ctx, identyfikator))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaInstrukcjiWorkspace{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaInstrukcjiWorkspace{}, fmt.Errorf("dane: nieczytelna wersja instrukcji %q: %w",
			identyfikator, err)
	}
	return wersja, nil
}

// WersjeInstrukcjiWorkspace zwraca historię instrukcji projektu od najnowszej.
func (r *repozytoriumPrzestrzeniRoboczej) WersjeInstrukcjiWorkspace(ctx context.Context,
	projektID int64) ([]WersjaInstrukcjiWorkspace, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWersjiInstrukcjiWorkspace)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wersji instrukcji projektu %d: %w",
			projektID, err)
	}
	defer wiersze.Close()

	lista := []WersjaInstrukcjiWorkspace{}
	for wiersze.Next() {
		wersja, err := odczytajWersjeInstrukcjiWorkspace(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji instrukcji: %w", err)
		}
		lista = append(lista, wersja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wersji instrukcji projektu %d: %w",
			projektID, err)
	}
	return lista, nil
}

// odczytajKomentarzWorkspace składa strukturę z jednego wiersza wyniku.
func odczytajKomentarzWorkspace(wiersz skaner) (KomentarzWorkspace, error) {
	var komentarz KomentarzWorkspace
	var rodzajBytu, rodzajAutora, przywolania string
	err := wiersz.Scan(&komentarz.ProjektID, &komentarz.ProjektKod, &komentarz.Identyfikator,
		&rodzajBytu, &komentarz.Byt, &komentarz.Tresc, &komentarz.KomentarzNadrzedny,
		&rodzajAutora, &komentarz.Autor, &przywolania, &komentarz.Utworzono,
		&komentarz.Zaktualizowano)
	if err != nil {
		return KomentarzWorkspace{}, err
	}
	komentarz.RodzajBytu = shared.WorkspaceEntityKind(rodzajBytu)
	komentarz.RodzajAutora = shared.WorkspaceAssigneeKind(rodzajAutora)
	komentarz.Przywolania = rozdzielWierszamiWorkspace(przywolania)
	return komentarz, nil
}

// odczytajWersjeInstrukcjiWorkspace składa strukturę z jednego wiersza wyniku.
func odczytajWersjeInstrukcjiWorkspace(wiersz skaner) (WersjaInstrukcjiWorkspace, error) {
	var wersja WersjaInstrukcjiWorkspace
	var poziom, rodzajAutora string
	err := wiersz.Scan(&wersja.ProjektID, &wersja.ProjektKod, &wersja.Identyfikator, &wersja.Tresc,
		&wersja.Odcisk, &poziom, &wersja.KluczZasiegu, &rodzajAutora, &wersja.Autor,
		&wersja.PrzywroconoZ, &wersja.Utworzono)
	if err != nil {
		return WersjaInstrukcjiWorkspace{}, err
	}
	zasieg, err := poziomZasieguZBazy(poziom)
	if err != nil {
		return WersjaInstrukcjiWorkspace{}, err
	}
	wersja.Poziom = zasieg
	wersja.RodzajAutora = shared.WorkspaceAssigneeKind(rodzajAutora)
	return wersja, nil
}
