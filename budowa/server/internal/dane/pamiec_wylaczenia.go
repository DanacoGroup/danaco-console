// Odpowiedzialność pliku: wyłączenia pamięci w danym zasięgu (tabela `wylaczenie_pamieci`, migracja 373).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// WylaczeniePamieci to wiersz tabeli `wylaczenie_pamieci`, niosący jedno wyłączenie pamięci w zasięgu.
type WylaczeniePamieci struct {
	ID int64
	// Identyfikator jest identyfikatorem kontraktu, nie numerem wiersza; zniesienie idzie nim z okna.
	Identyfikator string
	// WpisIdentyfikator jest identyfikatorem kontraktu wyłączonego wpisu; pusty znaczy cały poziom.
	WpisIdentyfikator string
	// PoziomPamieci jest wartością kontraktu (`MemoryLevel`); pusty znaczy
	// wyłączenie pojedynczego wpisu.
	PoziomPamieci  shared.MemoryLevel
	Zasieg         shared.ConfigScope
	KluczZasiegu   string
	Utworzono      string
	Zaktualizowano string
}

// RepozytoriumWylaczenPamieci jest kontraktem wyłączeń pamięci dla warstw wyższych całej tej platformy.
type RepozytoriumWylaczenPamieci interface {
	// ZapiszWylaczeniePamieci zakłada wyłączenie albo oddaje zastane; drugi wynik mówi, czy powstało.
	ZapiszWylaczeniePamieci(ctx context.Context,
		wylaczenie WylaczeniePamieci) (WylaczeniePamieci, bool, error)
	// ZniesWylaczeniePamieci usuwa wyłączenie po identyfikatorze; drugi wynik mówi, czy było co znosić.
	ZniesWylaczeniePamieci(ctx context.Context, identyfikator string) (bool, error)
	// WylaczeniePamieciPoBycie odnajduje wyłączenie złożone z bytu i zasięgu, bez identyfikatora.
	WylaczeniePamieciPoBycie(ctx context.Context, wzor WylaczeniePamieci) (WylaczeniePamieci, bool, error)
	// WylaczeniaPamieci zwraca wyłączenia zawężone niepustymi polami wzoru; wzór pusty zwraca wszystkie.
	WylaczeniaPamieci(ctx context.Context, wzor WylaczeniePamieci) ([]WylaczeniePamieci, error)
}

const (
	kolumnyWylaczeniaPamieci = `w.id, w.identyfikator_zewnetrzny, IFNULL(p.identyfikator_zewnetrzny, ''),
	                            w.poziom_pamieci, z.kod, w.klucz_zasiegu,
	                            w.utworzono, w.zaktualizowano`

	zrodloWylaczeniaPamieci = ` FROM wylaczenie_pamieci w
	                            JOIN poziom_zasiegu z ON z.id = w.poziom_zasiegu_id
	                            LEFT JOIN wpis_pamieci_projektu p ON p.id = w.wpis_id`

	zapiszWylaczeniePamieci = `INSERT INTO wylaczenie_pamieci
	                           (identyfikator_zewnetrzny, wpis_id, poziom_pamieci,
	                            poziom_zasiegu_id, klucz_zasiegu)
	                           VALUES (?,
	                                   (SELECT id FROM wpis_pamieci_projektu
	                                     WHERE identyfikator_zewnetrzny = ?),
	                                   ?,
	                                   (SELECT id FROM poziom_zasiegu WHERE kod = ?),
	                                   ?)`

	pobierzWylaczeniePamieci = `SELECT ` + kolumnyWylaczeniaPamieci + zrodloWylaczeniaPamieci +
		` WHERE w.identyfikator_zewnetrzny = ?`

	// Wyłączenie po bycie: pusty identyfikator wpisu i pusty poziom pamięci
	// wchodzą do warunku wprost, bo schemat trzyma je jako '' i NULL — jedno
	// z dwóch pól jest zawsze puste.
	pobierzWylaczeniePamieciPoBycie = `SELECT ` + kolumnyWylaczeniaPamieci + zrodloWylaczeniaPamieci +
		` WHERE IFNULL(p.identyfikator_zewnetrzny, '') = ?
		    AND w.poziom_pamieci = ?
		    AND z.kod = ?
		    AND w.klucz_zasiegu = ?`

	zniesWylaczeniePamieci = `DELETE FROM wylaczenie_pamieci WHERE identyfikator_zewnetrzny = ?`

	listaWylaczenPamieci = `SELECT ` + kolumnyWylaczeniaPamieci + zrodloWylaczeniaPamieci +
		` WHERE (? = '' OR IFNULL(p.identyfikator_zewnetrzny, '') = ?)
		    AND (? = '' OR w.poziom_pamieci = ?)
		    AND (? = '' OR z.kod = ?)
		    AND (? = '' OR w.klucz_zasiegu = ?)
		  ORDER BY w.utworzono DESC, w.id DESC`
)

type repozytoriumWylaczenPamieci struct {
	zapytania *zapytania
}

// noweRepozytoriumWylaczenPamieci zakłada magazyn wyłączeń pamięci nad zapytaniami tego całego zestawu.
func noweRepozytoriumWylaczenPamieci(z *zapytania) *repozytoriumWylaczenPamieci {
	return &repozytoriumWylaczenPamieci{zapytania: z}
}

// ZapiszWylaczeniePamieci zakłada wyłączenie; wyłączenie tego samego bytu w tym samym zasięgu jest już
// zapisane, więc drugie żądanie oddaje wiersz zastany.
func (r *repozytoriumWylaczenPamieci) ZapiszWylaczeniePamieci(ctx context.Context,
	wylaczenie WylaczeniePamieci) (WylaczeniePamieci, bool, error) {

	if strings.TrimSpace(wylaczenie.Identyfikator) == "" {
		return WylaczeniePamieci{}, false,
			fmt.Errorf("dane: wyłączenie pamięci bez identyfikatora")
	}
	if err := sprawdzBytWylaczenia(wylaczenie); err != nil {
		return WylaczeniePamieci{}, false, err
	}
	zastane, jest, err := r.WylaczeniePamieciPoBycie(ctx, wylaczenie)
	if err != nil {
		return WylaczeniePamieci{}, false, err
	}
	if jest {
		return zastane, false, nil
	}
	poziom, err := poziomZasieguNaBaze(wylaczenie.Zasieg)
	if err != nil {
		return WylaczeniePamieci{}, false, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWylaczeniePamieci)
	if err != nil {
		return WylaczeniePamieci{}, false, err
	}
	_, err = polecenie.ExecContext(ctx, wylaczenie.Identyfikator, wylaczenie.WpisIdentyfikator,
		string(wylaczenie.PoziomPamieci), poziom, wylaczenie.KluczZasiegu)
	if err != nil {
		return WylaczeniePamieci{}, false,
			fmt.Errorf("dane: nie można zapisać wyłączenia pamięci %q: %w",
				wylaczenie.Identyfikator, err)
	}
	zapisane, jest, err := r.wylaczeniePoIdentyfikatorze(ctx, wylaczenie.Identyfikator)
	if err != nil {
		return WylaczeniePamieci{}, false, err
	}
	if !jest {
		// Wiersz wstawiony i nieodnaleziony znaczy, że wskazanego wpisu pamięci nie ma w bazie danych.
		return WylaczeniePamieci{}, false,
			fmt.Errorf("dane: wyłączenie pamięci %q zapisane i nieodczytane",
				wylaczenie.Identyfikator)
	}
	return zapisane, true, nil
}

// ZniesWylaczeniePamieci usuwa wiersz wyłączenia z bazy danych rdzenia; treści wpisu nie dotyka wcale.
func (r *repozytoriumWylaczenPamieci) ZniesWylaczeniePamieci(ctx context.Context,
	identyfikator string) (bool, error) {

	if strings.TrimSpace(identyfikator) == "" {
		return false, fmt.Errorf("dane: zniesienie wyłączenia pamięci bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zniesWylaczeniePamieci)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, identyfikator)
	if err != nil {
		return false, fmt.Errorf("dane: nie można znieść wyłączenia pamięci %q: %w",
			identyfikator, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek zniesienia wyłączenia pamięci %q: %w",
			identyfikator, err)
	}
	return zmienione > 0, nil
}

// WylaczeniePamieciPoBycie odnajduje wyłączenie zapisane po bycie i jego zasięgu w bazie danych rdzenia.
func (r *repozytoriumWylaczenPamieci) WylaczeniePamieciPoBycie(ctx context.Context,
	wzor WylaczeniePamieci) (WylaczeniePamieci, bool, error) {

	poziom, err := poziomZasieguNaBaze(wzor.Zasieg)
	if err != nil {
		return WylaczeniePamieci{}, false, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWylaczeniePamieciPoBycie)
	if err != nil {
		return WylaczeniePamieci{}, false, err
	}
	wiersz := polecenie.QueryRowContext(ctx, wzor.WpisIdentyfikator,
		string(wzor.PoziomPamieci), poziom, wzor.KluczZasiegu)
	return odczytajJednoWylaczenie(wiersz)
}

// WylaczeniaPamieci zwraca wyłączenia zawężone niepustymi polami przekazanego wzoru tego wyszukiwania.
func (r *repozytoriumWylaczenPamieci) WylaczeniaPamieci(ctx context.Context,
	wzor WylaczeniePamieci) ([]WylaczeniePamieci, error) {

	zasieg := ""
	if wzor.Zasieg != "" {
		kod, err := poziomZasieguNaBaze(wzor.Zasieg)
		if err != nil {
			return nil, err
		}
		zasieg = kod
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaWylaczenPamieci)
	if err != nil {
		return nil, err
	}
	poziom := string(wzor.PoziomPamieci)
	wiersze, err := polecenie.QueryContext(ctx,
		wzor.WpisIdentyfikator, wzor.WpisIdentyfikator,
		poziom, poziom, zasieg, zasieg, wzor.KluczZasiegu, wzor.KluczZasiegu)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wyłączeń pamięci: %w", err)
	}
	defer wiersze.Close()

	wykaz := []WylaczeniePamieci{}
	for wiersze.Next() {
		wylaczenie, err := odczytajWylaczeniePamieci(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wyłączeń pamięci: %w", err)
		}
		wykaz = append(wykaz, wylaczenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wyłączeń pamięci: %w", err)
	}
	return wykaz, nil
}

// wylaczeniePoIdentyfikatorze zwraca jedno wyłączenie wskazane identyfikatorem kontraktu w bazie danych.
func (r *repozytoriumWylaczenPamieci) wylaczeniePoIdentyfikatorze(ctx context.Context,
	identyfikator string) (WylaczeniePamieci, bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWylaczeniePamieci)
	if err != nil {
		return WylaczeniePamieci{}, false, err
	}
	return odczytajJednoWylaczenie(polecenie.QueryRowContext(ctx, identyfikator))
}

// sprawdzBytWylaczenia pilnuje, żeby wyłączenie wskazywało dokładnie jeden byt:
// wpis albo poziom pamięci. Bez tego odmowa doszłaby do Operatora jako
// naruszenie warunku bazy, a nie jako nazwany brak w żądaniu.
func sprawdzBytWylaczenia(wylaczenie WylaczeniePamieci) error {
	maWpis := strings.TrimSpace(wylaczenie.WpisIdentyfikator) != ""
	maPoziom := strings.TrimSpace(string(wylaczenie.PoziomPamieci)) != ""
	if maWpis && maPoziom {
		return fmt.Errorf(
			"dane: wyłączenie pamięci wskazuje i wpis %q, i poziom %q — jedno albo drugie",
			wylaczenie.WpisIdentyfikator, wylaczenie.PoziomPamieci)
	}
	if !maWpis && !maPoziom {
		return fmt.Errorf("dane: wyłączenie pamięci bez wskazania wpisu ani poziomu pamięci")
	}
	return nil
}

// odczytajJednoWylaczenie zamienia brak wiersza na drugi wynik równy fałszowi:
// wyłączenia, którego nie ma, nie zgłasza się jako awarii odczytu.
func odczytajJednoWylaczenie(wiersz skaner) (WylaczeniePamieci, bool, error) {
	wylaczenie, err := odczytajWylaczeniePamieci(wiersz)
	if errors.Is(err, sql.ErrNoRows) {
		return WylaczeniePamieci{}, false, nil
	}
	if err != nil {
		return WylaczeniePamieci{}, false,
			fmt.Errorf("dane: nieczytelne wyłączenie pamięci: %w", err)
	}
	return wylaczenie, true, nil
}

// odczytajWylaczeniePamieci składa pełną strukturę wyłączenia z jednego wiersza wyniku tego zapytania.
func odczytajWylaczeniePamieci(wiersz skaner) (WylaczeniePamieci, error) {
	var wylaczenie WylaczeniePamieci
	var poziomPamieci, zasieg string
	err := wiersz.Scan(&wylaczenie.ID, &wylaczenie.Identyfikator, &wylaczenie.WpisIdentyfikator,
		&poziomPamieci, &zasieg, &wylaczenie.KluczZasiegu,
		&wylaczenie.Utworzono, &wylaczenie.Zaktualizowano)
	if err != nil {
		return WylaczeniePamieci{}, err
	}
	wylaczenie.PoziomPamieci = shared.MemoryLevel(poziomPamieci)
	if wylaczenie.Zasieg, err = poziomZasieguZBazy(zasieg); err != nil {
		return WylaczeniePamieci{}, err
	}
	return wylaczenie, nil
}
