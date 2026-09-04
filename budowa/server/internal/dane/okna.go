// Odpowiedzialność pliku: dostęp do obszaru okien komunikacji. Okno
// jest bytem pośrednim między sesją a wiadomością: niesie moduł, kanał modelu, środowisko wykonania,
// tryb uprawnień i rolę w pętli.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// Okno to wiersz tabeli `okno_komunikacji` wraz z pełną listą katalogów roboczych, jakie niesie to okno.
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
	// AgentKod niesie eksperta nałożonego na kanał modelu tego okna; pusty wskaźnik znaczy model surowy.
	AgentKod *string
	// IdentyfikatorZewnetrzny wiąże wiersz z oknem rdzenia, które żyje pod identyfikatorem tekstowym.
	IdentyfikatorZewnetrzny *string
	Kolejnosc               int
	Utworzono               string
	Zaktualizowano          string
}

// RepozytoriumOkien jest kontraktem obszaru okien komunikacji dla warstw wyższych całej tej platformy.
type RepozytoriumOkien interface {
	Utworz(ctx context.Context, okno Okno) (int64, error)
	Pobierz(ctx context.Context, id int64) (Okno, error)
	PoIdentyfikatorze(ctx context.Context, identyfikator string) (Okno, error)
	ListaSesji(ctx context.Context, sesjaID int64) ([]Okno, error)
	Aktualizuj(ctx context.Context, okno Okno) error
	ZmienStan(ctx context.Context, id int64, stan shared.WindowStatus) error
	// ZapiszRozmoweCLI utrwala identyfikator rozmowy nadany przez program CLI modelu.
	ZapiszRozmoweCLI(ctx context.Context, id int64, idRozmowy string) error
	// RozmowaCLI zwraca identyfikator rozmowy okna; pusty napis znaczy okno bez rozmowy z modelem.
	RozmowaCLI(ctx context.Context, id int64) (string, error)
	// LiczbaOtwartych liczy okna komunikacji o stanie otwartym — miara stanu platformy dla dowodzenia.
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
)

/*
warunekKontaOkna zwraca sprawdzenie własności okna dla okna nazwanego w zapytaniu
tabelą albo aliasem. Tabela `okno_komunikacji` nie ma kolumny `konto_id`: granica
dochodzi do niej przez sesję i kartę sesji (migracja 407), więc każde zapytanie
sięgające okna dokłada ten warunek zamiast porównania kolumny.

Aliasy wewnętrzne są własne (`so`, `ko`), żeby warunek wszedł także do zapytania,
które samo używa aliasów `s` i `k`.
*/
func warunekKontaOkna(okno string) string {
	return `EXISTS (SELECT 1 FROM sesja so
	                  JOIN karta_sesji ko ON ko.id = so.karta_sesji_id
	                 WHERE so.id = ` + okno + `.sesja_id
	                   AND ` + strings.ReplaceAll(WarunekKonta, "konto_id", "ko.konto_id") + `)`
}

var (
	kontoOknaWlasnego = warunekKontaOkna("okno_komunikacji")

	pobierzOkno = `SELECT ` + kolumnyOkna + ` FROM okno_komunikacji
	               WHERE id = ? AND ` + kontoOknaWlasnego

	oknoPoIdentyfikatorze = `SELECT ` + kolumnyOkna + ` FROM okno_komunikacji
	                         WHERE identyfikator_zewnetrzny = ? AND ` + kontoOknaWlasnego

	listaOkienSesji = `SELECT ` + kolumnyOkna + ` FROM okno_komunikacji
	                   WHERE sesja_id = ? AND ` + kontoOknaWlasnego + `
	                   ORDER BY kolejnosc, id`

	aktualizujOkno = `UPDATE okno_komunikacji
	                  SET modul_id = ?, kanal_modelu_id = ?, tytul = ?, srodowisko_wykonania = ?,
	                      tryb_uprawnien = ?, rola_okna = ?, okno_koordynatora_id = ?,
	                      tryb_komunikacji = ?, agent_kod = ?,
	                      zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                  WHERE id = ? AND ` + kontoOknaWlasnego

	zmienStanOkna = `UPDATE okno_komunikacji
	                 SET stan = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                 WHERE id = ? AND ` + kontoOknaWlasnego

	// Ciągłość rozmowy. Kolumna jest czytana i zapisywana osobnymi poleceniami,
	// a nie razem z resztą okna: identyfikator nadaje program `claude` w trakcie
	// tury, więc zmienia się niezależnie od pozostałych pól okna. Wspólny UPDATE
	// nadpisywałby jedno drugim.
	zapiszRozmoweCLI = `UPDATE okno_komunikacji
	                    SET id_rozmowy_cli = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                    WHERE id = ? AND ` + kontoOknaWlasnego

	odczytajRozmoweCLI = `SELECT id_rozmowy_cli FROM okno_komunikacji
	                      WHERE id = ? AND ` + kontoOknaWlasnego
)

type repozytoriumOkien struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumOkien(z *zapytania, db *sql.DB) *repozytoriumOkien {
	return &repozytoriumOkien{zapytania: z, db: db}
}

// Utworz zakłada nowe okno komunikacji wraz z jego listą katalogów roboczych w jednej transakcji zapisu.
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

// Pobierz zwraca okno komunikacji wraz z listą jego katalogów roboczych zapisanych w bazie danych rdzenia.
func (r *repozytoriumOkien) Pobierz(ctx context.Context, id int64) (Okno, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOkno)
	if err != nil {
		return Okno{}, err
	}
	okno, err := odczytajOkno(polecenie.QueryRowContext(ctx, id, KontoOperatora(ctx)))
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
	okno, err := odczytajOkno(polecenie.QueryRowContext(ctx, identyfikator, KontoOperatora(ctx)))
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

// ListaSesji zwraca okna danej sesji rozmowy; każde ma własny kanał modelu i własne katalogi robocze okna.
func (r *repozytoriumOkien) ListaSesji(ctx context.Context, sesjaID int64) ([]Okno, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaOkienSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, sesjaID, KontoOperatora(ctx))
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

// Aktualizuj zapisuje zmienione parametry okna razem z listą jego katalogów w jednej transakcji zapisu.
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
			tekstDoKolumny(okno.AgentKod), okno.ID, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać zmian okna %d: %w", okno.ID, err)
		}
		if err := sprawdzTrafienie(wynik, "okno_komunikacji", okno.ID); err != nil {
			return err
		}
		return zapiszKatalogiOkna(ctx, r.zapytania, transakcja, okno.ID, okno.KatalogiRobocze)
	})
}

// ZmienStan otwiera albo zamyka okno komunikacji, zapisując jego nowy stan wprost w bazie danych rdzenia.
func (r *repozytoriumOkien) ZmienStan(ctx context.Context, id int64, stan shared.WindowStatus) error {
	kolumna, err := stanOknaNaBaze(stan)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienStanOkna)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kolumna, id, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić stanu okna %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "okno_komunikacji", id)
}

// ZapiszRozmoweCLI utrwala identyfikator rozmowy programu wiersza poleceń przy danym oknie komunikacji.
func (r *repozytoriumOkien) ZapiszRozmoweCLI(ctx context.Context, id int64, idRozmowy string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszRozmoweCLI)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, idRozmowy, id, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać rozmowy CLI okna %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "okno_komunikacji", id)
}

// RozmowaCLI zwraca identyfikator rozmowy okna komunikacji albo pusty napis, jeśli jej jeszcze nie było.
func (r *repozytoriumOkien) RozmowaCLI(ctx context.Context, id int64) (string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, odczytajRozmoweCLI)
	if err != nil {
		return "", err
	}
	var idRozmowy string
	if err := polecenie.QueryRowContext(ctx, id, KontoOperatora(ctx)).Scan(&idRozmowy); err != nil {
		return "", fmt.Errorf("dane: nie można odczytać rozmowy CLI okna %d: %w", id, err)
	}
	return idRozmowy, nil
}
