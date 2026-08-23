// Odpowiedzialność pliku: sterowanie pojedynczym krokiem zlecenia — tabela
// `wstrzymanie_kroku`. Krokiem jest wiersz `pozycja_kolejki`; tutaj nie ma
// drugiego bytu kroku ani drugiej kopii jego stanu pracy. Jest wyłącznie to,
// czego `pozycja_kolejki` nie umie powiedzieć: że krok stoi, na jaką decyzję
// czeka, jaka decyzja zapadła i czy dojechała do wykonawcy.
//
// Metody siedzą na `repozytoriumKolejek`, a nie na własnym typie, bo wstrzymanie
// kroku i stan kroku to dwa pytania o tę samą tabelę `pozycja_kolejki` i ten sam
// dziennik `log_akcji_kolejki`. Osobne repozytorium musiałoby powtórzyć odczyt
// położenia pozycji i własnym zapisem dziennika rozjechać się z zapisem silnika.
// Wzorem par Moduly/Macierz i Sesje/KoszSesji: jedna implementacja, dwa widoki.
// Rozszerzenie widać przez `RepozytoriumWstrzymanKroku` — port sięga po nie
// asercją typu, tak jak rdzeń sięga po `wiazaneKolejki` przy
// `queue.list`/`queue.link`.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Stany sterowania krokiem — słownik kolumny `wstrzymanie_kroku.stan`.
// Wartości `czeka` i `biegnie` tu nie należą: czyta się je z
// `pozycja_kolejki.stan`, a zapisane drugi raz byłyby drugą prawdą.
const (
	StanKrokuWstrzymany   = "wstrzymany"
	StanKrokuZatwierdzony = "zatwierdzony"
	StanKrokuOdrzucony    = "odrzucony"
)

// Akcje dziennika kolejki zapisywane przez sterowanie krokiem. Idą do tego
// samego `log_akcji_kolejki`, co działania silnika — jeden ślad, nie dwa.
const (
	AkcjaWstrzymanieKroku = "wstrzymanie_kroku"
	AkcjaDecyzjaOKroku    = "decyzja_o_kroku"
	AkcjaZastosowanie     = "zastosowanie_decyzji_kroku"
)

// WstrzymanieKroku to wiersz tabeli `wstrzymanie_kroku` — jeden epizod
// wstrzymania jednego kroku wraz z losem decyzji Operatora.
type WstrzymanieKroku struct {
	ID               int64
	PozycjaID        int64
	Stan             string
	StanPozycjiPrzed string
	Powod            *string
	Uzasadnienie     *string
	WstrzymanoO      string
	ZdecydowanoO     *string
	ZastosowanoO     *string
	DoreczonoO       *string
}

// CzyZdecydowane mówi, czy Operator już rozstrzygnął ten epizod.
func (w WstrzymanieKroku) CzyZdecydowane() bool {
	return w.Stan != StanKrokuWstrzymany
}

// CzyZastosowane mówi, czy decyzja zmieniła już los kroku. Epizod zdecydowany,
// lecz niezastosowany, to decyzja w drodze — nie wolno jej zgubić.
func (w WstrzymanieKroku) CzyZastosowane() bool {
	return w.ZastosowanoO != nil
}

// RepozytoriumWstrzymanKroku jest rozszerzeniem obszaru kolejek o sterowanie
// pojedynczym krokiem. Wypełnia je ta sama implementacja, co
// `RepozytoriumKolejek` — drugiej nie ma.
type RepozytoriumWstrzymanKroku interface {
	// WstrzymajKrok zakłada epizod wstrzymania. `stanPozycjiPrzed` musi być
	// stanem roboczym — stanu końcowego schemat odmawia.
	WstrzymajKrok(ctx context.Context, pozycjaID int64, stanPozycjiPrzed string,
		powod *string) (WstrzymanieKroku, error)

	// CzynneWstrzymanie zwraca epizod niezastosowany, jeśli taki jest.
	CzynneWstrzymanie(ctx context.Context, pozycjaID int64) (WstrzymanieKroku, bool, error)

	// ZapiszDecyzje odnotowuje rozstrzygnięcie Operatora. Nie stosuje go —
	// zastosowanie jest osobnym faktem i osobnym zapisem.
	ZapiszDecyzje(ctx context.Context, wstrzymanieID int64, stan string,
		uzasadnienie *string) (WstrzymanieKroku, error)

	// OznaczZastosowanie zamyka epizod: decyzja zmieniła los kroku.
	// `doreczono` mówi, czy trafiła też do wykonawcy — fałsz zostawia
	// `doreczono_o` pusty i to jest prawda o braku wykonawcy, nie usterka.
	OznaczZastosowanie(ctx context.Context, wstrzymanieID int64,
		doreczono bool) (WstrzymanieKroku, error)

	// CzynneWstrzymaniaKolejki zwraca epizody niezastosowane wszystkich kroków
	// kolejki, po identyfikatorze pozycji. Kolejka bez wstrzymań daje mapę pustą.
	CzynneWstrzymaniaKolejki(ctx context.Context, kolejkaID int64) (map[int64]WstrzymanieKroku, error)

	// HistoriaWstrzymanKroku zwraca wszystkie epizody kroku, od najstarszego.
	HistoriaWstrzymanKroku(ctx context.Context, pozycjaID int64) ([]WstrzymanieKroku, error)
}

const (
	kolumnyWstrzymania = `id, pozycja_kolejki_id, stan, stan_pozycji_przed, powod,
	                      uzasadnienie, wstrzymano_o, zdecydowano_o, zastosowano_o, doreczono_o`

	wstawWstrzymanieKroku = `INSERT INTO wstrzymanie_kroku
	                         (pozycja_kolejki_id, stan, stan_pozycji_przed, powod)
	                         VALUES (?, '` + StanKrokuWstrzymany + `', ?, ?)`

	pobierzWstrzymanie = `SELECT ` + kolumnyWstrzymania + `
	                      FROM wstrzymanie_kroku WHERE id = ?`

	czynneWstrzymaniePozycji = `SELECT ` + kolumnyWstrzymania + `
	                            FROM wstrzymanie_kroku
	                            WHERE pozycja_kolejki_id = ? AND zastosowano_o IS NULL`

	czynneWstrzymaniaKolejki = `SELECT w.id, w.pozycja_kolejki_id, w.stan, w.stan_pozycji_przed,
	                                   w.powod, w.uzasadnienie, w.wstrzymano_o,
	                                   w.zdecydowano_o, w.zastosowano_o, w.doreczono_o
	                            FROM wstrzymanie_kroku w
	                            JOIN pozycja_kolejki p ON p.id = w.pozycja_kolejki_id
	                            WHERE p.kolejka_id = ? AND w.zastosowano_o IS NULL
	                            ORDER BY w.pozycja_kolejki_id, w.id`

	historiaWstrzymanKroku = `SELECT ` + kolumnyWstrzymania + `
	                          FROM wstrzymanie_kroku
	                          WHERE pozycja_kolejki_id = ? ORDER BY id`

	// Zapis decyzji wchodzi wyłącznie na epizod jeszcze nierozstrzygnięty
	// i jeszcze niezastosowany. Warunek stoi w SQL, nie w warstwie wyżej —
	// inaczej dwa równoległe rozstrzygnięcia nadpisałyby się nawzajem.
	zapiszDecyzjeKroku = `UPDATE wstrzymanie_kroku
	                      SET stan = ?, uzasadnienie = ?,
	                          zdecydowano_o = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                      WHERE id = ? AND zdecydowano_o IS NULL AND zastosowano_o IS NULL`

	// Zastosowanie wchodzi wyłącznie na epizod rozstrzygnięty i jeszcze
	// niezastosowany — powtórne wywołanie nie zrobi drugiego zastosowania.
	oznaczZastosowanieKroku = `UPDATE wstrzymanie_kroku
	                           SET zastosowano_o = strftime('%Y-%m-%dT%H:%M:%fZ','now'),
	                               doreczono_o = CASE WHEN ? = 1
	                                                  THEN strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                                  ELSE NULL END
	                           WHERE id = ? AND zdecydowano_o IS NOT NULL AND zastosowano_o IS NULL`
)

// WstrzymajKrok zakłada epizod wstrzymania i odnotowuje go w dzienniku kolejki.
// Krok w stanie końcowym odbija się o CHECK schematu — odmowa jest wtedy
// prawdą o kroku, nie awarią zapisu.
func (r *repozytoriumKolejek) WstrzymajKrok(ctx context.Context, pozycjaID int64,
	stanPozycjiPrzed string, powod *string) (WstrzymanieKroku, error) {

	var id int64
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		kolejkaID, obieg, err := polozeniePozycji(ctx, r.zapytania, transakcja, pozycjaID)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWstrzymanieKroku)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, pozycjaID, stanPozycjiPrzed, tekstDoKolumny(powod))
		if err != nil {
			return fmt.Errorf("dane: nie można wstrzymać kroku %d: %w", pozycjaID, err)
		}
		if id, err = wynik.LastInsertId(); err != nil {
			return fmt.Errorf("dane: nieznany identyfikator wstrzymania kroku: %w", err)
		}
		stan := StanKrokuWstrzymany
		przed := stanPozycjiPrzed
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: kolejkaID, PozycjaID: &pozycjaID, Akcja: AkcjaWstrzymanieKroku,
			StanPrzed: &przed, StanPo: &stan, NumerObiegu: obieg, Szczegoly: powod,
		})
	})
	if err != nil {
		return WstrzymanieKroku{}, err
	}
	return r.wstrzymanie(ctx, id)
}

// CzynneWstrzymanie zwraca epizod niezastosowany kroku. Brak epizodu nie jest
// błędem — większość kroków nigdy nie była wstrzymana.
func (r *repozytoriumKolejek) CzynneWstrzymanie(ctx context.Context,
	pozycjaID int64) (WstrzymanieKroku, bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, czynneWstrzymaniePozycji)
	if err != nil {
		return WstrzymanieKroku{}, false, err
	}
	wstrzymanie, err := odczytajWstrzymanie(polecenie.QueryRowContext(ctx, pozycjaID))
	if errors.Is(err, sql.ErrNoRows) {
		return WstrzymanieKroku{}, false, nil
	}
	if err != nil {
		return WstrzymanieKroku{}, false, err
	}
	return wstrzymanie, true, nil
}

// ZapiszDecyzje odnotowuje rozstrzygnięcie Operatora. Epizod już rozstrzygnięty
// albo już zastosowany nie zostaje ruszony — zapytanie ma warunek i brak
// trafienia jest tu odmową, nie ciszą.
func (r *repozytoriumKolejek) ZapiszDecyzje(ctx context.Context, wstrzymanieID int64,
	stan string, uzasadnienie *string) (WstrzymanieKroku, error) {

	if stan != StanKrokuZatwierdzony && stan != StanKrokuOdrzucony {
		return WstrzymanieKroku{}, fmt.Errorf(
			"dane: %q nie jest decyzją o kroku (zatwierdzony albo odrzucony)", stan)
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		pozycjaID, err := pozycjaWstrzymania(ctx, r.zapytania, transakcja, wstrzymanieID)
		if err != nil {
			return err
		}
		kolejkaID, obieg, err := polozeniePozycji(ctx, r.zapytania, transakcja, pozycjaID)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszDecyzjeKroku)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, stan, tekstDoKolumny(uzasadnienie), wstrzymanieID)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać decyzji o kroku %d: %w", pozycjaID, err)
		}
		zmienione, err := wynik.RowsAffected()
		if err != nil {
			return fmt.Errorf("dane: nieznany wynik zapisu decyzji o kroku %d: %w", pozycjaID, err)
		}
		if zmienione == 0 {
			return fmt.Errorf(
				"dane: wstrzymanie %d nie czeka na decyzję (już rozstrzygnięte albo zastosowane)",
				wstrzymanieID)
		}
		przed := StanKrokuWstrzymany
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: kolejkaID, PozycjaID: &pozycjaID, Akcja: AkcjaDecyzjaOKroku,
			StanPrzed: &przed, StanPo: &stan, NumerObiegu: obieg, Szczegoly: uzasadnienie,
		})
	})
	if err != nil {
		return WstrzymanieKroku{}, err
	}
	return r.wstrzymanie(ctx, wstrzymanieID)
}

// OznaczZastosowanie zamyka epizod. Doręczenie zapisuje się osobno od
// zastosowania, bo to dwa różne fakty: „krok ruszył zgodnie z decyzją" i
// „decyzja dojechała do tego, kto pracę wykonuje".
func (r *repozytoriumKolejek) OznaczZastosowanie(ctx context.Context, wstrzymanieID int64,
	doreczono bool) (WstrzymanieKroku, error) {

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		pozycjaID, err := pozycjaWstrzymania(ctx, r.zapytania, transakcja, wstrzymanieID)
		if err != nil {
			return err
		}
		kolejkaID, obieg, err := polozeniePozycji(ctx, r.zapytania, transakcja, pozycjaID)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, oznaczZastosowanieKroku)
		if err != nil {
			return err
		}
		znacznik := 0
		if doreczono {
			znacznik = 1
		}
		wynik, err := polecenie.ExecContext(ctx, znacznik, wstrzymanieID)
		if err != nil {
			return fmt.Errorf("dane: nie można zastosować decyzji o kroku %d: %w", pozycjaID, err)
		}
		zmienione, err := wynik.RowsAffected()
		if err != nil {
			return fmt.Errorf("dane: nieznany wynik zastosowania decyzji o kroku %d: %w", pozycjaID, err)
		}
		if zmienione == 0 {
			return fmt.Errorf(
				"dane: wstrzymanie %d nie ma decyzji do zastosowania albo jest już zastosowane",
				wstrzymanieID)
		}
		szczegoly := "decyzja doręczona wykonawcy"
		if !doreczono {
			szczegoly = "decyzja zastosowana bez wykonawcy — nie doręczona"
		}
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: kolejkaID, PozycjaID: &pozycjaID, Akcja: AkcjaZastosowanie,
			NumerObiegu: obieg, Szczegoly: &szczegoly,
		})
	})
	if err != nil {
		return WstrzymanieKroku{}, err
	}
	return r.wstrzymanie(ctx, wstrzymanieID)
}

// CzynneWstrzymaniaKolejki zwraca epizody niezastosowane kroków jednej kolejki.
// Indeks częściowy schematu gwarantuje najwyżej jeden epizod na krok, więc mapa
// po identyfikatorze pozycji niczego nie gubi.
func (r *repozytoriumKolejek) CzynneWstrzymaniaKolejki(ctx context.Context,
	kolejkaID int64) (map[int64]WstrzymanieKroku, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, czynneWstrzymaniaKolejki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kolejkaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wstrzymań kroków kolejki %d: %w", kolejkaID, err)
	}
	defer wiersze.Close()

	wykaz := map[int64]WstrzymanieKroku{}
	for wiersze.Next() {
		wstrzymanie, err := odczytajWstrzymanie(wiersze)
		if err != nil {
			return nil, err
		}
		wykaz[wstrzymanie.PozycjaID] = wstrzymanie
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wstrzymań kroków kolejki %d: %w", kolejkaID, err)
	}
	return wykaz, nil
}

// HistoriaWstrzymanKroku zwraca wszystkie epizody kroku. Wiersze zastosowane
// zostają — ślad wstrzymania jest częścią przejrzystości pętli.
func (r *repozytoriumKolejek) HistoriaWstrzymanKroku(ctx context.Context,
	pozycjaID int64) ([]WstrzymanieKroku, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, historiaWstrzymanKroku)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, pozycjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać historii wstrzymań kroku %d: %w", pozycjaID, err)
	}
	defer wiersze.Close()

	lista := []WstrzymanieKroku{}
	for wiersze.Next() {
		wstrzymanie, err := odczytajWstrzymanie(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, wstrzymanie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt historii wstrzymań kroku %d: %w", pozycjaID, err)
	}
	return lista, nil
}

// wstrzymanie odczytuje jeden epizod po identyfikatorze.
func (r *repozytoriumKolejek) wstrzymanie(ctx context.Context, id int64) (WstrzymanieKroku, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWstrzymanie)
	if err != nil {
		return WstrzymanieKroku{}, err
	}
	wstrzymanie, err := odczytajWstrzymanie(polecenie.QueryRowContext(ctx, id))
	if errors.Is(err, sql.ErrNoRows) {
		return WstrzymanieKroku{}, fmt.Errorf("dane: wstrzymanie kroku %d nie istnieje", id)
	}
	if err != nil {
		return WstrzymanieKroku{}, err
	}
	return wstrzymanie, nil
}

// pozycjaWstrzymania odczytuje krok, którego dotyczy epizod — w tej samej
// transakcji, co zapis, żeby dziennik nie wskazał innego kroku niż zapis.
func pozycjaWstrzymania(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	wstrzymanieID int64) (int64, error) {

	polecenie, err := z.wTransakcji(ctx, transakcja,
		`SELECT pozycja_kolejki_id FROM wstrzymanie_kroku WHERE id = ?`)
	if err != nil {
		return 0, err
	}
	var pozycjaID int64
	err = polecenie.QueryRowContext(ctx, wstrzymanieID).Scan(&pozycjaID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("dane: wstrzymanie kroku %d nie istnieje", wstrzymanieID)
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odczytać wstrzymania kroku %d: %w", wstrzymanieID, err)
	}
	return pozycjaID, nil
}

// odczytajWstrzymanie składa strukturę z jednego wiersza wyniku.
func odczytajWstrzymanie(wiersz skaner) (WstrzymanieKroku, error) {
	var w WstrzymanieKroku
	var powod, uzasadnienie, zdecydowano, zastosowano, doreczono sql.NullString
	err := wiersz.Scan(&w.ID, &w.PozycjaID, &w.Stan, &w.StanPozycjiPrzed, &powod,
		&uzasadnienie, &w.WstrzymanoO, &zdecydowano, &zastosowano, &doreczono)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WstrzymanieKroku{}, err
		}
		return WstrzymanieKroku{}, fmt.Errorf("dane: nieczytelny wiersz wstrzymania kroku: %w", err)
	}
	w.Powod = tekstZKolumny(powod)
	w.Uzasadnienie = tekstZKolumny(uzasadnienie)
	w.ZdecydowanoO = tekstZKolumny(zdecydowano)
	w.ZastosowanoO = tekstZKolumny(zastosowano)
	w.DoreczonoO = tekstZKolumny(doreczono)
	return w, nil
}
