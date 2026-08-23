// Odpowiedzialność pliku: dostęp do obszaru okien komunikacji. Okno
// jest bytem pośrednim między sesją a wiadomością: niesie moduł, kanał modelu,
// środowisko wykonania, tryb uprawnień i rolę w pętli. Zapis dotyka
// dwóch tabel (`okno_komunikacji`, `katalog_okna`), więc idzie w transakcji.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// Okno to wiersz tabeli `okno_komunikacji` wraz z listą katalogów roboczych.
type Okno struct {
	ID                  int64
	SesjaID             int64
	ModulID             int64
	KanalModeluID       int64
	Tytul               *string
	KatalogiRobocze     []string
	SrodowiskoWykonania shared.ExecutionEnv
	TrybUprawnien       shared.PermissionMode
	RolaOkna            shared.WindowRole
	OknoKoordynatoraID  *int64
	TrybKomunikacji     string
	Stan                shared.WindowStatus
	// AgentKod niesie eksperta nałożonego na kanał modelu tego okna. Pusty
	// wskaźnik znaczy model surowy. Kod, nie identyfikator wiersza — ekspert
	// bywa kasowany niezależnie od okien, w których pracował.
	AgentKod *string
	// IdentyfikatorZewnetrzny wiąże wiersz z oknem rdzenia, które żyje pod
	// identyfikatorem tekstowym. Bez niego po restarcie rdzenia nie da się
	// połączyć okna wskazanego przez klienta z jego historią.
	IdentyfikatorZewnetrzny *string
	Kolejnosc               int
	Utworzono               string
	Zaktualizowano          string
}

// RepozytoriumOkien jest kontraktem obszaru okien dla warstw wyższych.
type RepozytoriumOkien interface {
	Utworz(ctx context.Context, okno Okno) (int64, error)
	Pobierz(ctx context.Context, id int64) (Okno, error)
	PoIdentyfikatorze(ctx context.Context, identyfikator string) (Okno, error)
	ListaSesji(ctx context.Context, sesjaID int64) ([]Okno, error)
	Aktualizuj(ctx context.Context, okno Okno) error
	ZmienStan(ctx context.Context, id int64, stan shared.WindowStatus) error
	// ZapiszRozmoweCLI utrwala identyfikator rozmowy nadany przez program
	// `claude`, dzięki któremu następna tura wznawia rozmowę zamiast zaczynać
	// od zera.
	ZapiszRozmoweCLI(ctx context.Context, id int64, idRozmowy string) error
	// RozmowaCLI zwraca identyfikator rozmowy okna. Pusty napis znaczy „okno
	// nie rozmawiało jeszcze z modelem" i jest stanem poprawnym.
	RozmowaCLI(ctx context.Context, id int64) (string, error)
	// LiczbaOtwartych liczy okna o stanie `otwarte` — miara stanu platformy dla
	// mobilnego centrum dowodzenia. Rachunek stoi tutaj, bo tabelę
	// `okno_komunikacji` prowadzi to repozytorium i drugiego czytelnika mieć nie
	// będzie; ciało metody leży w `mobile.go`.
	LiczbaOtwartych(ctx context.Context) (int, error)
}

const (
	kolumnyOkna = `id, sesja_id, modul_id, kanal_modelu_id, tytul, srodowisko_wykonania,
	               tryb_uprawnien, rola_okna, okno_koordynatora_id, tryb_komunikacji, stan,
	               identyfikator_zewnetrzny, agent_kod, kolejnosc, utworzono, zaktualizowano`

	wstawOkno = `INSERT INTO okno_komunikacji
	             (sesja_id, modul_id, kanal_modelu_id, tytul, srodowisko_wykonania, tryb_uprawnien,
	              rola_okna, okno_koordynatora_id, tryb_komunikacji, stan, identyfikator_zewnetrzny,
	              agent_kod, kolejnosc)
	             VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	                     (SELECT COALESCE(MAX(kolejnosc) + 1, 0) FROM okno_komunikacji WHERE sesja_id = ?))`

	pobierzOkno = `SELECT ` + kolumnyOkna + ` FROM okno_komunikacji WHERE id = ?`

	oknoPoIdentyfikatorze = `SELECT ` + kolumnyOkna + ` FROM okno_komunikacji
	                         WHERE identyfikator_zewnetrzny = ?`

	listaOkienSesji = `SELECT ` + kolumnyOkna + ` FROM okno_komunikacji
	                   WHERE sesja_id = ? ORDER BY kolejnosc, id`

	aktualizujOkno = `UPDATE okno_komunikacji
	                  SET modul_id = ?, kanal_modelu_id = ?, tytul = ?, srodowisko_wykonania = ?,
	                      tryb_uprawnien = ?, rola_okna = ?, okno_koordynatora_id = ?,
	                      tryb_komunikacji = ?, agent_kod = ?,
	                      zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                  WHERE id = ?`

	zmienStanOkna = `UPDATE okno_komunikacji
	                 SET stan = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                 WHERE id = ?`

	// Ciągłość rozmowy. Kolumna jest czytana i zapisywana osobnymi poleceniami,
	// a nie razem z resztą okna: identyfikator nadaje program `claude` w trakcie
	// tury, więc zmienia się niezależnie od pozostałych pól okna. Wspólny UPDATE
	// nadpisywałby jedno drugim.
	zapiszRozmoweCLI = `UPDATE okno_komunikacji
	                    SET id_rozmowy_cli = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                    WHERE id = ?`

	odczytajRozmoweCLI = `SELECT id_rozmowy_cli FROM okno_komunikacji WHERE id = ?`
)

type repozytoriumOkien struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumOkien(z *zapytania, db *sql.DB) *repozytoriumOkien {
	return &repozytoriumOkien{zapytania: z, db: db}
}

// Utworz zakłada okno wraz z jego listą katalogów roboczych — jedna transakcja.
func (r *repozytoriumOkien) Utworz(ctx context.Context, okno Okno) (int64, error) {
	wartosci, err := wartosciOkna(okno)
	if err != nil {
		return 0, err
	}
	var id int64
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawOkno)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, okno.SesjaID, okno.ModulID, okno.KanalModeluID,
			tekstDoKolumny(okno.Tytul), wartosci.srodowisko, wartosci.tryb, wartosci.rola,
			liczbaDoKolumny(okno.OknoKoordynatoraID), wartosci.trybKomunikacji, wartosci.stan,
			tekstDoKolumny(okno.IdentyfikatorZewnetrzny), tekstDoKolumny(okno.AgentKod),
			okno.SesjaID)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać okna sesji %d: %w", okno.SesjaID, err)
		}
		if id, err = wynik.LastInsertId(); err != nil {
			return fmt.Errorf("dane: nieznany identyfikator zapisanego okna: %w", err)
		}
		return zapiszKatalogiOkna(ctx, r.zapytania, transakcja, id, okno.KatalogiRobocze)
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Pobierz zwraca okno wraz z listą katalogów roboczych.
func (r *repozytoriumOkien) Pobierz(ctx context.Context, id int64) (Okno, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOkno)
	if err != nil {
		return Okno{}, err
	}
	okno, err := odczytajOkno(polecenie.QueryRowContext(ctx, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Okno{}, fmt.Errorf("dane: okno %d nie istnieje", id)
	}
	if err != nil {
		return Okno{}, err
	}
	okno.KatalogiRobocze, err = katalogiOkna(ctx, r.zapytania, id)
	if err != nil {
		return Okno{}, err
	}
	return okno, nil
}

// PoIdentyfikatorze zwraca okno po identyfikatorze nadanym przez rdzeń wraz
// z listą katalogów roboczych. Brak wiersza sygnalizuje ErrBrakWiersza — dla
// warstwy wyższej jest to stan normalny, nie awaria.
func (r *repozytoriumOkien) PoIdentyfikatorze(ctx context.Context, identyfikator string) (Okno, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, oknoPoIdentyfikatorze)
	if err != nil {
		return Okno{}, err
	}
	okno, err := odczytajOkno(polecenie.QueryRowContext(ctx, identyfikator))
	if errors.Is(err, sql.ErrNoRows) {
		return Okno{}, fmt.Errorf("%w: okno %q", ErrBrakWiersza, identyfikator)
	}
	if err != nil {
		return Okno{}, err
	}
	okno.KatalogiRobocze, err = katalogiOkna(ctx, r.zapytania, okno.ID)
	if err != nil {
		return Okno{}, err
	}
	return okno, nil
}

// ListaSesji zwraca okna sesji; każde ma własny kanał i własne katalogi.
func (r *repozytoriumOkien) ListaSesji(ctx context.Context, sesjaID int64) ([]Okno, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaOkienSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, sesjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać okien sesji %d: %w", sesjaID, err)
	}
	defer wiersze.Close()

	lista := []Okno{}
	for wiersze.Next() {
		okno, err := odczytajOkno(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, okno)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt okien sesji %d: %w", sesjaID, err)
	}
	katalogi, err := katalogiOkienSesji(ctx, r.zapytania, sesjaID)
	if err != nil {
		return nil, err
	}
	for i := range lista {
		lista[i].KatalogiRobocze = katalogi[lista[i].ID]
	}
	return lista, nil
}

// Aktualizuj zapisuje parametry okna razem z listą katalogów — jedna transakcja.
func (r *repozytoriumOkien) Aktualizuj(ctx context.Context, okno Okno) error {
	wartosci, err := wartosciOkna(okno)
	if err != nil {
		return err
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, aktualizujOkno)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, okno.ModulID, okno.KanalModeluID,
			tekstDoKolumny(okno.Tytul), wartosci.srodowisko, wartosci.tryb, wartosci.rola,
			liczbaDoKolumny(okno.OknoKoordynatoraID), wartosci.trybKomunikacji,
			tekstDoKolumny(okno.AgentKod), okno.ID)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać zmian okna %d: %w", okno.ID, err)
		}
		if err := sprawdzTrafienie(wynik, "okno_komunikacji", okno.ID); err != nil {
			return err
		}
		return zapiszKatalogiOkna(ctx, r.zapytania, transakcja, okno.ID, okno.KatalogiRobocze)
	})
}

// ZmienStan otwiera albo zamyka okno.
func (r *repozytoriumOkien) ZmienStan(ctx context.Context, id int64, stan shared.WindowStatus) error {
	kolumna, err := stanOknaNaBaze(stan)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienStanOkna)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kolumna, id)
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić stanu okna %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "okno_komunikacji", id)
}

// ZapiszRozmoweCLI utrwala identyfikator rozmowy programu `claude` przy oknie.
//
// Wołane po każdej turze, bo `claude` może nadać identyfikator dopiero w
// trakcie pierwszej wymiany. Zapis pustego napisu jest dozwolony i znaczy
// „zacznij następną turę od nowa" — na przykład po przeniesieniu kontekstu.
func (r *repozytoriumOkien) ZapiszRozmoweCLI(ctx context.Context, id int64, idRozmowy string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszRozmoweCLI)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, idRozmowy, id)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać rozmowy CLI okna %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "okno_komunikacji", id)
}

// RozmowaCLI zwraca identyfikator rozmowy okna albo pusty napis.
func (r *repozytoriumOkien) RozmowaCLI(ctx context.Context, id int64) (string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, odczytajRozmoweCLI)
	if err != nil {
		return "", err
	}
	var idRozmowy string
	if err := polecenie.QueryRowContext(ctx, id).Scan(&idRozmowy); err != nil {
		return "", fmt.Errorf("dane: nie można odczytać rozmowy CLI okna %d: %w", id, err)
	}
	return idRozmowy, nil
}
