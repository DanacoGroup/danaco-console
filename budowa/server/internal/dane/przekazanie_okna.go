// Odpowiedzialność pliku: obszar window.* — deklaracja całego interfejsu
// `RepozytoriumPrzekazan` wraz z typem repozytorium i konstruktorem,
// oraz implementacja zlecenia przekazania (tabela `zlecenie_przekazania`),
// które zasila komendę `window.handoff`. Więź koordynator–wykonawca leży
// w `przekazanie_okna_wiez.go`, dziennik akcji w `przekazanie_okna_akcje.go`
// — jedno repozytorium, trzy pliki wedle odpowiedzialności, tak jak
// `dane/automations*.go` i `dane/library*.go`.
//
// Interfejs deklaruje wyłącznie ten plik, w całości — wraz z metodami, które
// implementują pozostałe pliki obszaru. Interfejs rozdzielony na trzy pliki
// byłby trzema prawdami o jednym kontrakcie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ZleceniePrzekazania to wiersz tabeli `zlecenie_przekazania` — treść jednego
// wywołania `window.handoff`. `PozycjaKolejkiID` nie jest zakładana tu własnym
// zapisem SQL: pozycję kolejki zakłada jedyny silnik kolejek, a adapter rdzenia
// wypełnia to pole gotowym identyfikatorem po założeniu pozycji, w tej samej
// transakcji co `ZapiszZlecenie`. `KompletKontekstu` niesie `ContextBundle`
// w całości jako surowy zapis JSON — warstwa danych go nie interpretuje, nie
// rozbiera na pola, tylko przechowuje i oddaje.
type ZleceniePrzekazania struct {
	ID               int64
	SesjaID          int64
	OknoZrodloweID   int64
	OknoDoceloweID   int64
	PozycjaKolejkiID int64
	Polecenie        string
	KompletKontekstu *string
	Utworzono        string
}

// RepozytoriumPrzekazan jest kontraktem obszaru window.*.
type RepozytoriumPrzekazan interface {
	// --- zlecenie ---
	ZapiszZlecenie(ctx context.Context, zlecenie ZleceniePrzekazania) (ZleceniePrzekazania, error)
	Zlecenie(ctx context.Context, id int64) (ZleceniePrzekazania, error)
	ZleceniaSesji(ctx context.Context, sesja string) ([]ZleceniePrzekazania, error)

	// --- więź koordynator–wykonawca ---
	UstawKoordynatora(ctx context.Context, oknoWykonawcy, oknoKoordynatora string) error
	Koordynator(ctx context.Context, oknoWykonawcy string) (string, error)
	Wykonawcy(ctx context.Context, oknoKoordynatora string) ([]string, error)

	// --- dziennik akcji ---
	ZapiszAkcje(ctx context.Context, akcja AkcjaOkna) (AkcjaOkna, error)
	AkcjeOkna(ctx context.Context, okno string, limit int) ([]AkcjaOkna, error)
}

const (
	kolumnyZleceniaPrzekazania = `id, sesja_id, okno_zrodlowe_id, okno_docelowe_id,
	                              pozycja_kolejki_id, polecenie, komplet_kontekstu, utworzono`

	zapiszZleceniePrzekazania = `INSERT INTO zlecenie_przekazania
	                             (sesja_id, okno_zrodlowe_id, okno_docelowe_id, pozycja_kolejki_id,
	                              polecenie, komplet_kontekstu)
	                             VALUES (?, ?, ?, ?, ?, ?)`

	pobierzZleceniePrzekazania = `SELECT ` + kolumnyZleceniaPrzekazania + `
	                              FROM zlecenie_przekazania WHERE id = ?`

	pobierzZleceniaSesji = `SELECT ` + kolumnyZleceniaPrzekazania + `
	                        FROM zlecenie_przekazania WHERE sesja_id = ?
	                        ORDER BY utworzono DESC, id DESC`
)

type repozytoriumPrzekazan struct {
	zapytania *zapytania
	db        *sql.DB
}

// noweRepozytoriumPrzekazan konstruuje repozytorium obszaru window.*.
// Bez odbiorcy w tym pliku — wpina go koordynator (`dane/zestaw.go`).
func noweRepozytoriumPrzekazan(z *zapytania, db *sql.DB) *repozytoriumPrzekazan {
	return &repozytoriumPrzekazan{zapytania: z, db: db}
}

// ZapiszZlecenie zakłada wiersz zlecenia przekazania i zwraca stan po zapisie.
// Wymaga gotowej `PozycjaKolejkiID` — zakładanie pozycji kolejki leży poza tym
// repozytorium (patrz komentarz typu).
func (r *repozytoriumPrzekazan) ZapiszZlecenie(ctx context.Context,
	zlecenie ZleceniePrzekazania) (ZleceniePrzekazania, error) {

	if zlecenie.SesjaID == 0 {
		return ZleceniePrzekazania{}, fmt.Errorf("dane: zlecenie przekazania bez sesji")
	}
	if zlecenie.OknoZrodloweID == 0 || zlecenie.OknoDoceloweID == 0 {
		return ZleceniePrzekazania{}, fmt.Errorf("dane: zlecenie przekazania bez okna źródłowego lub docelowego")
	}
	if zlecenie.PozycjaKolejkiID == 0 {
		return ZleceniePrzekazania{}, fmt.Errorf("dane: zlecenie przekazania bez pozycji kolejki")
	}
	if zlecenie.Polecenie == "" {
		return ZleceniePrzekazania{}, fmt.Errorf("dane: zlecenie przekazania bez polecenia")
	}

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZleceniePrzekazania)
	if err != nil {
		return ZleceniePrzekazania{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, zlecenie.SesjaID, zlecenie.OknoZrodloweID,
		zlecenie.OknoDoceloweID, zlecenie.PozycjaKolejkiID, zlecenie.Polecenie,
		tekstDoKolumny(zlecenie.KompletKontekstu))
	if err != nil {
		return ZleceniePrzekazania{}, fmt.Errorf("dane: nie można zapisać zlecenia przekazania: %w", err)
	}
	id, err := wynik.LastInsertId()
	if err != nil {
		return ZleceniePrzekazania{}, fmt.Errorf("dane: nie można odczytać id zlecenia przekazania: %w", err)
	}
	return r.Zlecenie(ctx, id)
}

// Zlecenie zwraca zlecenie przekazania o wskazanym id. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się nie powiódł”.
func (r *repozytoriumPrzekazan) Zlecenie(ctx context.Context, id int64) (ZleceniePrzekazania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZleceniePrzekazania)
	if err != nil {
		return ZleceniePrzekazania{}, err
	}
	zlecenie, err := odczytajZlecenieRekord(polecenie.QueryRowContext(ctx, id))
	if errors.Is(err, sql.ErrNoRows) {
		return ZleceniePrzekazania{}, ErrBrakWiersza
	}
	if err != nil {
		return ZleceniePrzekazania{}, fmt.Errorf("dane: nieczytelny wiersz zlecenia przekazania %d: %w", id, err)
	}
	return zlecenie, nil
}

// ZleceniaSesji zwraca zlecenia przekazania danej sesji uporządkowane od
// najnowszego — zasila Mission Control (przegląd zleceń w sesji).
func (r *repozytoriumPrzekazan) ZleceniaSesji(ctx context.Context, sesja string) ([]ZleceniePrzekazania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZleceniaSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, sesja)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zleceń przekazania sesji %q: %w", sesja, err)
	}
	defer wiersze.Close()

	lista := []ZleceniePrzekazania{}
	for wiersze.Next() {
		zlecenie, err := odczytajZlecenieRekord(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zlecenia przekazania: %w", err)
		}
		lista = append(lista, zlecenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zleceń przekazania sesji %q: %w", sesja, err)
	}
	return lista, nil
}

// odczytajZlecenieRekord składa strukturę z jednego wiersza wyniku.
func odczytajZlecenieRekord(wiersz skaner) (ZleceniePrzekazania, error) {
	var zlecenie ZleceniePrzekazania
	var kompletKontekstu sql.NullString
	err := wiersz.Scan(&zlecenie.ID, &zlecenie.SesjaID, &zlecenie.OknoZrodloweID,
		&zlecenie.OknoDoceloweID, &zlecenie.PozycjaKolejkiID, &zlecenie.Polecenie,
		&kompletKontekstu, &zlecenie.Utworzono)
	if err != nil {
		return ZleceniePrzekazania{}, err
	}
	zlecenie.KompletKontekstu = tekstZKolumny(kompletKontekstu)
	return zlecenie, nil
}
