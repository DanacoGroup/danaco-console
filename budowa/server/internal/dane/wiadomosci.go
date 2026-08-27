// Odpowiedzialność pliku: dostęp do obszaru wiadomości (tabela `wiadomosc`).
// Wiadomość należy do okna komunikacji, nie wprost do sesji — dwa okna
// tej samej sesji mają rozdzielone historie. Treść obszerna trafia do pliku,
// baza trzyma odwołanie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// Wiadomosc to wiersz tabeli wiadomosc, niosący treść i metadane jednej
// pozycji historii okna komunikacji.
type Wiadomosc struct {
	ID             int64
	OknoID         int64
	Rola           shared.MessageRole
	Persona        *string
	OknoZrodloweID *int64
	RodzajTresci   shared.ChunkKind
	Stan           shared.MessageStatus
	Tresc          *string
	TrescOdwolanie *string
	TokenyWejscia  int
	TokenyWyjscia  int
	// IdentyfikatorZewnetrzny wiąże wiersz z wiadomością rdzenia, znaną
	// klientowi identyfikatorem.
	IdentyfikatorZewnetrzny *string
	Kolejnosc               int
	Utworzono               string
	// Zalaczniki niesie tablicę JSON odwołań; nil znaczy brak rozpoznania,
	// nie pustą tablicę.
	Zalaczniki *string
}

// RepozytoriumWiadomosci jest kontraktem obszaru wiadomości: zapis, odczyt
// historii okna i zmianę stanu.
type RepozytoriumWiadomosci interface {
	Dopisz(ctx context.Context, wiadomosc Wiadomosc) (int64, error)
	ListaOkna(ctx context.Context, oknoID int64, limit int) ([]Wiadomosc, error)
	PoIdentyfikatorze(ctx context.Context, identyfikator string) (Wiadomosc, error)
	ZapiszWynik(ctx context.Context, wiadomosc Wiadomosc) error
	ZmienStan(ctx context.Context, id int64, stan shared.MessageStatus) error
}

const (
	kolumnyWiadomosci = `id, okno_komunikacji_id, rola, persona, okno_zrodlowe_id, rodzaj_tresci,
	                     stan, tresc, tresc_odwolanie, tokeny_wejscia, tokeny_wyjscia,
	                     identyfikator_zewnetrzny, kolejnosc, utworzono, zalaczniki`

	// Kolejność wyliczana w tym samym poleceniu — numeracja historii okna nie
	// wymaga osobnego odczytu ani transakcji.
	wstawWiadomosc = `INSERT INTO wiadomosc
	                  (okno_komunikacji_id, rola, persona, okno_zrodlowe_id, rodzaj_tresci, stan,
	                   tresc, tresc_odwolanie, tokeny_wejscia, tokeny_wyjscia,
	                   identyfikator_zewnetrzny, zalaczniki, kolejnosc)
	                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	                          (SELECT COALESCE(MAX(kolejnosc) + 1, 0) FROM wiadomosc
	                           WHERE okno_komunikacji_id = ?))`

	wiadomoscPoIdentyfikatorze = `SELECT ` + kolumnyWiadomosci + ` FROM wiadomosc
	                              WHERE identyfikator_zewnetrzny = ?`

	listaWiadomosciOkna = `SELECT ` + kolumnyWiadomosci + ` FROM wiadomosc
	                       WHERE okno_komunikacji_id = ?
	                       ORDER BY kolejnosc, id
	                       LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`

	zapiszWynikWiadomosci = `UPDATE wiadomosc
	                         SET tresc = ?, tresc_odwolanie = ?, rodzaj_tresci = ?, stan = ?,
	                             tokeny_wejscia = ?, tokeny_wyjscia = ?
	                         WHERE id = ?`

	zmienStanWiadomosci = `UPDATE wiadomosc SET stan = ? WHERE id = ?`
)

type repozytoriumWiadomosci struct {
	zapytania *zapytania
}

func noweRepozytoriumWiadomosci(z *zapytania) *repozytoriumWiadomosci {
	return &repozytoriumWiadomosci{zapytania: z}
}

// Dopisz dokłada wiadomość na koniec historii okna i zwraca jej identyfikator
// wiersza, nadając kolejny numer kolejności.
func (r *repozytoriumWiadomosci) Dopisz(ctx context.Context, wiadomosc Wiadomosc) (int64, error) {
	rola, err := rolaWiadomosciNaBaze(wiadomosc.Rola)
	if err != nil {
		return 0, err
	}
	rodzaj, err := rodzajTresciNaBaze(wiadomosc.RodzajTresci)
	if err != nil {
		return 0, err
	}
	stan, err := stanWiadomosciNaBaze(wiadomosc.Stan)
	if err != nil {
		return 0, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawWiadomosc)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, wiadomosc.OknoID, rola,
		tekstDoKolumny(wiadomosc.Persona), liczbaDoKolumny(wiadomosc.OknoZrodloweID), rodzaj, stan,
		tekstDoKolumny(wiadomosc.Tresc), tekstDoKolumny(wiadomosc.TrescOdwolanie),
		wiadomosc.TokenyWejscia, wiadomosc.TokenyWyjscia,
		tekstDoKolumny(wiadomosc.IdentyfikatorZewnetrzny), tekstDoKolumny(wiadomosc.Zalaczniki),
		wiadomosc.OknoID)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zapisać wiadomości okna %d: %w", wiadomosc.OknoID, err)
	}
	return wynik.LastInsertId()
}

// PoIdentyfikatorze zwraca wiadomość po identyfikatorze zewnętrznym nadanym
// przez rdzeń, którym klient rozpoznaje wiersz.
func (r *repozytoriumWiadomosci) PoIdentyfikatorze(ctx context.Context, identyfikator string) (Wiadomosc, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wiadomoscPoIdentyfikatorze)
	if err != nil {
		return Wiadomosc{}, err
	}
	wiadomosc, err := odczytajWiadomosc(polecenie.QueryRowContext(ctx, identyfikator))
	if errors.Is(err, sql.ErrNoRows) {
		return Wiadomosc{}, fmt.Errorf("%w: wiadomość %q", ErrBrakWiersza, identyfikator)
	}
	return wiadomosc, err
}

// ListaOkna zwraca historię jednego okna komunikacji w kolejności zapisu.
// Limit 0 oznacza całą historię.
func (r *repozytoriumWiadomosci) ListaOkna(ctx context.Context, oknoID int64, limit int) ([]Wiadomosc, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaWiadomosciOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoID, limit, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wiadomości okna %d: %w", oknoID, err)
	}
	defer wiersze.Close()

	lista := []Wiadomosc{}
	for wiersze.Next() {
		wiadomosc, err := odczytajWiadomosc(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, wiadomosc)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wiadomości okna %d: %w", oknoID, err)
	}
	return lista, nil
}

// ZapiszWynik utrwala treść wiadomości po zakończeniu strumienia modelu wraz
// ze zużyciem tokenów wejścia i wyjścia.
func (r *repozytoriumWiadomosci) ZapiszWynik(ctx context.Context, wiadomosc Wiadomosc) error {
	rodzaj, err := rodzajTresciNaBaze(wiadomosc.RodzajTresci)
	if err != nil {
		return err
	}
	stan, err := stanWiadomosciNaBaze(wiadomosc.Stan)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWynikWiadomosci)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, tekstDoKolumny(wiadomosc.Tresc),
		tekstDoKolumny(wiadomosc.TrescOdwolanie), rodzaj, stan,
		wiadomosc.TokenyWejscia, wiadomosc.TokenyWyjscia, wiadomosc.ID)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać treści wiadomości %d: %w", wiadomosc.ID, err)
	}
	return sprawdzTrafienie(wynik, "wiadomosc", wiadomosc.ID)
}

// ZmienStan zapisuje nowy stan wiadomości — także stan `zatrzymana`, bo przycisk
// zatrzymania jest zawsze czynny.
func (r *repozytoriumWiadomosci) ZmienStan(ctx context.Context, id int64, stan shared.MessageStatus) error {
	kolumna, err := stanWiadomosciNaBaze(stan)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienStanWiadomosci)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kolumna, id)
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić stanu wiadomości %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "wiadomosc", id)
}

// odczytajWiadomosc składa strukturę wiadomości z jednego wiersza wyniku
// zapytania, tłumacząc rolę, rodzaj treści i stan.
func odczytajWiadomosc(wiersz skaner) (Wiadomosc, error) {
	var wiadomosc Wiadomosc
	var persona, tresc, odwolanie, identyfikator, zalaczniki sql.NullString
	var oknoZrodlowe sql.NullInt64
	var rola, rodzaj, stan string
	err := wiersz.Scan(&wiadomosc.ID, &wiadomosc.OknoID, &rola, &persona, &oknoZrodlowe, &rodzaj,
		&stan, &tresc, &odwolanie, &wiadomosc.TokenyWejscia, &wiadomosc.TokenyWyjscia,
		&identyfikator, &wiadomosc.Kolejnosc, &wiadomosc.Utworzono, &zalaczniki)
	if errors.Is(err, sql.ErrNoRows) {
		return Wiadomosc{}, err
	}
	if err != nil {
		return Wiadomosc{}, fmt.Errorf("dane: nieczytelny wiersz wiadomości: %w", err)
	}
	wiadomosc.IdentyfikatorZewnetrzny = tekstZKolumny(identyfikator)
	wiadomosc.Persona = tekstZKolumny(persona)
	wiadomosc.Tresc = tekstZKolumny(tresc)
	wiadomosc.TrescOdwolanie = tekstZKolumny(odwolanie)
	wiadomosc.OknoZrodloweID = liczbaZKolumny(oknoZrodlowe)
	wiadomosc.Zalaczniki = tekstZKolumny(zalaczniki)
	if wiadomosc.Rola, err = rolaWiadomosciZBazy(rola); err != nil {
		return Wiadomosc{}, err
	}
	if wiadomosc.RodzajTresci, err = rodzajTresciZBazy(rodzaj); err != nil {
		return Wiadomosc{}, err
	}
	if wiadomosc.Stan, err = stanWiadomosciZBazy(stan); err != nil {
		return Wiadomosc{}, err
	}
	return wiadomosc, nil
}
