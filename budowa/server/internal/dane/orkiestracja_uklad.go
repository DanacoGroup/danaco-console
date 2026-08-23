// Odpowiedzialność pliku: trzy dopełnienia układu zależności — bramka
// dołączenia, grupa kroków i krok wycofujący — oraz spięcie kolejek automatyki
// z silnikiem kolejek środowiska MultitaskingAI (tabele `orkiestracja_bramka`,
// `orkiestracja_grupa`, `orkiestracja_grupa_krok`, `orkiestracja_kompensacja`
// i `orkiestracja_spiecie_multitasking` z migracji 279–280).
//
// Plik osobny od `orchestration.go`: tamten prowadzi łuki układu
// (`zaleznosc_kroku_automatyki`), tutaj mieszkają byty, których łuk nie wyraża.
// Bramka mówi, KIEDY tory scalają się w jednym kroku; grupa mówi, że zbiór
// kroków biegnie razem; kompensacja mówi, co zrobić, gdy przebieg pękł w pół.
//
// Wszystkie posługują się identyfikatorem ZEWNĘTRZNYM kroku, tak jak łuki —
// układ przepisuje się w całości, więc klucz wiersza kroku nie przeżywa
// przepisania, a napis Operatora przeżywa.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// BramkaUkladu to reguła scalenia torów równoległych na jednym kroku.
type BramkaUkladu struct {
	Krok    string
	Regula  string
	Licznik int
}

// GrupaUkladu to zbiór kroków wykonywanych równolegle albo w ścisłej kolejności.
type GrupaUkladu struct {
	Kod    string
	Nazwa  string
	Rodzaj string
	Kroki  []string
}

// KompensacjaUkladu to krok wycofujący skutki kroku głównego.
type KompensacjaUkladu struct {
	Krok       string
	KrokWycofu string
}

// SpiecieMultitaskingu to stan spięcia automatyki z kolejkami środowiska
// MultitaskingAI wraz z kolejkami objętymi spięciem.
type SpiecieMultitaskingu struct {
	Spiete     bool
	OknoRoliID *int64
	KolejkiID  []int64
}

// RepozytoriumUkladuOrkiestracji jest kontraktem trzech dopełnień układu
// i spięcia kolejek.
type RepozytoriumUkladuOrkiestracji interface {
	BramkiUkladu(ctx context.Context, automatykaID int64) ([]BramkaUkladu, error)
	ZapiszBramkeUkladu(ctx context.Context, automatykaID int64, bramka BramkaUkladu) error
	GrupyUkladu(ctx context.Context, automatykaID int64) ([]GrupaUkladu, error)
	// ZapiszGrupeUkladu zakłada grupę o pustym kodzie albo zmienia wskazaną.
	// Wykaz kroków pusty USUWA grupę — tak stanowi kontrakt komendy.
	ZapiszGrupeUkladu(ctx context.Context, automatykaID int64, grupa GrupaUkladu) error
	KompensacjeUkladu(ctx context.Context, automatykaID int64) ([]KompensacjaUkladu, error)
	// ZapiszKompensacjeUkladu zapisuje kompensację; krok wycofujący pusty ją
	// zdejmuje.
	ZapiszKompensacjeUkladu(ctx context.Context, automatykaID int64, kompensacja KompensacjaUkladu) error
	// SpiecieAutomatyki oddaje stan spięcia wraz z kolejkami objętymi.
	SpiecieAutomatyki(ctx context.Context, automatykaID int64, kodAutomatyki string) (SpiecieMultitaskingu, error)
	// ZapiszSpiecieAutomatyki spina albo rozłącza kolejki automatyki z silnikiem
	// kolejek środowiska MultitaskingAI i oddaje kolejki objęte zmianą.
	ZapiszSpiecieAutomatyki(ctx context.Context, automatykaID int64, kodAutomatyki string,
		spiete bool, oknoRoliID *int64) ([]int64, error)
}

const (
	bramkiUkladuAutomatyki = `SELECT krok, regula, licznik FROM orkiestracja_bramka
	                           WHERE automatyka_id = ? ORDER BY krok`

	zapiszBramkeUkladuAutomatyki = `INSERT INTO orkiestracja_bramka (automatyka_id, krok, regula, licznik)
	                                VALUES (?, ?, ?, ?)
	                                ON CONFLICT(automatyka_id, krok) DO UPDATE SET
	                                    regula = excluded.regula,
	                                    licznik = excluded.licznik,
	                                    zapisano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	grupyUkladuAutomatyki = `SELECT g.id, g.kod, g.nazwa, g.rodzaj FROM orkiestracja_grupa g
	                          WHERE g.automatyka_id = ? ORDER BY g.nazwa, g.id`

	krokiGrupyUkladu = `SELECT krok FROM orkiestracja_grupa_krok
	                     WHERE grupa_id = ? ORDER BY kolejnosc, krok`

	kompensacjeUkladuAutomatyki = `SELECT krok, krok_wycofu FROM orkiestracja_kompensacja
	                                WHERE automatyka_id = ? ORDER BY krok`

	zapiszKompensacjeUkladuAutomatyki = `INSERT INTO orkiestracja_kompensacja
	                                     (automatyka_id, krok, krok_wycofu) VALUES (?, ?, ?)
	                                     ON CONFLICT(automatyka_id, krok) DO UPDATE SET
	                                         krok_wycofu = excluded.krok_wycofu,
	                                         zapisano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	// Kolejki automatyki rozpoznaje nazwa równa jej kodowi — tak zakłada je
	// `core/adapter_modul_automations_kolejka.go` i innej drogi nie ma.
	kolejkiAutomatykiUkladu = `SELECT id FROM kolejka WHERE nazwa = ? ORDER BY id`
)

// repozytoriumUkladuOrkiestracji obsługuje dopełnienia układu zależności.
type repozytoriumUkladuOrkiestracji struct {
	zapytania *zapytania
	db        *sql.DB
}

var _ RepozytoriumUkladuOrkiestracji = (*repozytoriumUkladuOrkiestracji)(nil)

// noweRepozytoriumUkladuOrkiestracji wiąże dopełnienia układu z bazą.
func noweRepozytoriumUkladuOrkiestracji(z *zapytania, db *sql.DB) *repozytoriumUkladuOrkiestracji {
	return &repozytoriumUkladuOrkiestracji{zapytania: z, db: db}
}

// BramkiUkladu oddaje bramki dołączenia automatyki.
func (r *repozytoriumUkladuOrkiestracji) BramkiUkladu(ctx context.Context,
	automatykaID int64) ([]BramkaUkladu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, bramkiUkladuAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać bramek układu %d: %w", automatykaID, err)
	}
	defer wiersze.Close()
	bramki := make([]BramkaUkladu, 0, 8)
	for wiersze.Next() {
		var bramka BramkaUkladu
		if err := wiersze.Scan(&bramka.Krok, &bramka.Regula, &bramka.Licznik); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna bramka układu %d: %w", automatykaID, err)
		}
		bramki = append(bramki, bramka)
	}
	return bramki, wiersze.Err()
}

// ZapiszBramkeUkladu zakłada albo zmienia bramkę dołączenia.
func (r *repozytoriumUkladuOrkiestracji) ZapiszBramkeUkladu(ctx context.Context,
	automatykaID int64, bramka BramkaUkladu) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszBramkeUkladuAutomatyki)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, automatykaID, bramka.Krok,
		bramka.Regula, bramka.Licznik); err != nil {
		return fmt.Errorf("dane: nie można zapisać bramki kroku %q: %w", bramka.Krok, err)
	}
	return nil
}

// GrupyUkladu oddaje grupy kroków wraz z ich składem.
func (r *repozytoriumUkladuOrkiestracji) GrupyUkladu(ctx context.Context,
	automatykaID int64) ([]GrupaUkladu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, grupyUkladuAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać grup układu %d: %w", automatykaID, err)
	}
	defer wiersze.Close()
	numery := make([]int64, 0, 8)
	grupy := make([]GrupaUkladu, 0, 8)
	for wiersze.Next() {
		var grupa GrupaUkladu
		var numer int64
		if err := wiersze.Scan(&numer, &grupa.Kod, &grupa.Nazwa, &grupa.Rodzaj); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna grupa układu %d: %w", automatykaID, err)
		}
		numery = append(numery, numer)
		grupy = append(grupy, grupa)
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	for pozycja := range grupy {
		kroki, err := r.krokiGrupy(ctx, numery[pozycja])
		if err != nil {
			return nil, err
		}
		grupy[pozycja].Kroki = kroki
	}
	return grupy, nil
}

// krokiGrupy oddaje skład jednej grupy w zapisanej kolejności.
func (r *repozytoriumUkladuOrkiestracji) krokiGrupy(ctx context.Context, grupaID int64) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, krokiGrupyUkladu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, grupaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kroków grupy %d: %w", grupaID, err)
	}
	defer wiersze.Close()
	kroki := make([]string, 0, 8)
	for wiersze.Next() {
		var krok string
		if err := wiersze.Scan(&krok); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny krok grupy %d: %w", grupaID, err)
		}
		kroki = append(kroki, krok)
	}
	return kroki, wiersze.Err()
}

// ZapiszGrupeUkladu zakłada, zmienia albo usuwa grupę kroków.
func (r *repozytoriumUkladuOrkiestracji) ZapiszGrupeUkladu(ctx context.Context,
	automatykaID int64, grupa GrupaUkladu) error {

	transakcja, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("dane: nie można otworzyć zapisu grupy układu %d: %w", automatykaID, err)
	}
	defer func() { _ = transakcja.Rollback() }()

	// Wykaz kroków pusty usuwa grupę — tak stanowi kontrakt komendy. Grupa bez
	// ani jednego kroku nie mówi o niczym, a zostawiona w wykazie byłaby
	// wpisem, którego Operator nie umie ani wypełnić, ani skasować.
	if len(grupa.Kroki) == 0 {
		if grupa.Kod == "" {
			return nil
		}
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM orkiestracja_grupa WHERE kod = ? AND automatyka_id = ?`,
			grupa.Kod, automatykaID); err != nil {
			return fmt.Errorf("dane: nie można usunąć grupy %q: %w", grupa.Kod, err)
		}
		return transakcja.Commit()
	}

	var numer int64
	if grupa.Kod == "" {
		return fmt.Errorf("dane: grupa układu %d bez kodu", automatykaID)
	}
	wynik, err := transakcja.ExecContext(ctx,
		`INSERT INTO orkiestracja_grupa (kod, automatyka_id, nazwa, rodzaj) VALUES (?, ?, ?, ?)
		 ON CONFLICT(kod) DO UPDATE SET nazwa = excluded.nazwa, rodzaj = excluded.rodzaj`,
		grupa.Kod, automatykaID, grupa.Nazwa, grupa.Rodzaj)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać grupy %q: %w", grupa.Kod, err)
	}
	numer, err = wynik.LastInsertId()
	if err != nil || numer == 0 {
		if err := transakcja.QueryRowContext(ctx,
			`SELECT id FROM orkiestracja_grupa WHERE kod = ?`, grupa.Kod).Scan(&numer); err != nil {
			return fmt.Errorf("dane: nie można odnaleźć grupy %q po zapisie: %w", grupa.Kod, err)
		}
	}
	if _, err := transakcja.ExecContext(ctx,
		`DELETE FROM orkiestracja_grupa_krok WHERE grupa_id = ?`, numer); err != nil {
		return fmt.Errorf("dane: nie można wyczyścić składu grupy %q: %w", grupa.Kod, err)
	}
	for kolejnosc, krok := range grupa.Kroki {
		przyciety := strings.TrimSpace(krok)
		if przyciety == "" {
			continue
		}
		if _, err := transakcja.ExecContext(ctx,
			`INSERT INTO orkiestracja_grupa_krok (grupa_id, krok, kolejnosc) VALUES (?, ?, ?)
			 ON CONFLICT(grupa_id, krok) DO UPDATE SET kolejnosc = excluded.kolejnosc`,
			numer, przyciety, kolejnosc); err != nil {
			return fmt.Errorf("dane: nie można zapisać kroku %q grupy %q: %w", przyciety, grupa.Kod, err)
		}
	}
	return transakcja.Commit()
}

// KompensacjeUkladu oddaje kroki wycofujące automatyki.
func (r *repozytoriumUkladuOrkiestracji) KompensacjeUkladu(ctx context.Context,
	automatykaID int64) ([]KompensacjaUkladu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, kompensacjeUkladuAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kompensacji układu %d: %w", automatykaID, err)
	}
	defer wiersze.Close()
	kompensacje := make([]KompensacjaUkladu, 0, 8)
	for wiersze.Next() {
		var wpis KompensacjaUkladu
		if err := wiersze.Scan(&wpis.Krok, &wpis.KrokWycofu); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna kompensacja układu %d: %w", automatykaID, err)
		}
		kompensacje = append(kompensacje, wpis)
	}
	return kompensacje, wiersze.Err()
}

// ZapiszKompensacjeUkladu zapisuje albo zdejmuje krok wycofujący.
func (r *repozytoriumUkladuOrkiestracji) ZapiszKompensacjeUkladu(ctx context.Context,
	automatykaID int64, kompensacja KompensacjaUkladu) error {

	if strings.TrimSpace(kompensacja.KrokWycofu) == "" {
		polecenie, err := r.zapytania.przygotuj(ctx,
			`DELETE FROM orkiestracja_kompensacja WHERE automatyka_id = ? AND krok = ?`)
		if err != nil {
			return err
		}
		if _, err := polecenie.ExecContext(ctx, automatykaID, kompensacja.Krok); err != nil {
			return fmt.Errorf("dane: nie można zdjąć kompensacji kroku %q: %w", kompensacja.Krok, err)
		}
		return nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKompensacjeUkladuAutomatyki)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, automatykaID, kompensacja.Krok, kompensacja.KrokWycofu); err != nil {
		return fmt.Errorf("dane: nie można zapisać kompensacji kroku %q: %w", kompensacja.Krok, err)
	}
	return nil
}

// SpiecieAutomatyki oddaje stan spięcia wraz z kolejkami objętymi.
func (r *repozytoriumUkladuOrkiestracji) SpiecieAutomatyki(ctx context.Context,
	automatykaID int64, kodAutomatyki string) (SpiecieMultitaskingu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx,
		`SELECT okno_roli_id FROM orkiestracja_spiecie_multitasking WHERE automatyka_id = ?`)
	if err != nil {
		return SpiecieMultitaskingu{}, err
	}
	var okno sql.NullInt64
	stan := SpiecieMultitaskingu{}
	switch err := polecenie.QueryRowContext(ctx, automatykaID).Scan(&okno); {
	case err == sql.ErrNoRows:
		stan.Spiete = false
	case err != nil:
		return SpiecieMultitaskingu{}, fmt.Errorf("dane: nie można odczytać spięcia automatyki %d: %w",
			automatykaID, err)
	default:
		stan.Spiete = true
		stan.OknoRoliID = liczbaZKolumny(okno)
	}
	kolejki, err := r.kolejkiAutomatyki(ctx, kodAutomatyki)
	if err != nil {
		return SpiecieMultitaskingu{}, err
	}
	stan.KolejkiID = kolejki
	return stan, nil
}

// ZapiszSpiecieAutomatyki spina albo rozłącza kolejki automatyki.
//
// Skutek leży w tabeli `kolejka`, nie w samym znaczniku: spięcie przestawia
// kolejki automatyki na rodzaj `multitasking` i wiąże je z oknem roli, a
// rozłączenie zdejmuje wiązanie. Znacznik bez tego byłby zapisem, którego nikt
// nie czyta.
func (r *repozytoriumUkladuOrkiestracji) ZapiszSpiecieAutomatyki(ctx context.Context,
	automatykaID int64, kodAutomatyki string, spiete bool, oknoRoliID *int64) ([]int64, error) {

	transakcja, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można otworzyć zapisu spięcia automatyki %d: %w", automatykaID, err)
	}
	defer func() { _ = transakcja.Rollback() }()

	if spiete {
		if _, err := transakcja.ExecContext(ctx,
			`INSERT INTO orkiestracja_spiecie_multitasking (automatyka_id, okno_roli_id) VALUES (?, ?)
			 ON CONFLICT(automatyka_id) DO UPDATE SET okno_roli_id = excluded.okno_roli_id,
			     zapisano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`,
			automatykaID, liczbaDoKolumny(oknoRoliID)); err != nil {
			return nil, fmt.Errorf("dane: nie można zapisać spięcia automatyki %d: %w", automatykaID, err)
		}
		if _, err := transakcja.ExecContext(ctx,
			`UPDATE kolejka SET rodzaj = 'multitasking', okno_koordynatora_id = ?,
			     zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
			  WHERE nazwa = ?`, liczbaDoKolumny(oknoRoliID), kodAutomatyki); err != nil {
			return nil, fmt.Errorf("dane: nie można spiąć kolejek automatyki %q: %w", kodAutomatyki, err)
		}
	} else {
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM orkiestracja_spiecie_multitasking WHERE automatyka_id = ?`,
			automatykaID); err != nil {
			return nil, fmt.Errorf("dane: nie można rozłączyć automatyki %d: %w", automatykaID, err)
		}
		if _, err := transakcja.ExecContext(ctx,
			`UPDATE kolejka SET okno_koordynatora_id = NULL,
			     zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
			  WHERE nazwa = ?`, kodAutomatyki); err != nil {
			return nil, fmt.Errorf("dane: nie można rozłączyć kolejek automatyki %q: %w", kodAutomatyki, err)
		}
	}
	if err := transakcja.Commit(); err != nil {
		return nil, fmt.Errorf("dane: nie można domknąć zapisu spięcia automatyki %d: %w", automatykaID, err)
	}
	return r.kolejkiAutomatyki(ctx, kodAutomatyki)
}

// kolejkiAutomatyki oddaje klucze kolejek noszących pracę automatyki.
func (r *repozytoriumUkladuOrkiestracji) kolejkiAutomatyki(ctx context.Context,
	kodAutomatyki string) ([]int64, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, kolejkiAutomatykiUkladu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAutomatyki)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kolejek automatyki %q: %w", kodAutomatyki, err)
	}
	defer wiersze.Close()
	kolejki := make([]int64, 0, 4)
	for wiersze.Next() {
		var id int64
		if err := wiersze.Scan(&id); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna kolejka automatyki %q: %w", kodAutomatyki, err)
		}
		kolejki = append(kolejki, id)
	}
	return kolejki, wiersze.Err()
}
