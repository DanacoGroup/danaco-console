// Odpowiedzialność pliku: kolejka wczytywania i cyfryzacji modułu Studio
// (tabela `pozycja_wczytywania_studio`) — Ingest/OCR Panel.
//
// Kolejka należy do OKNA, nie do dokumentu: pozycja istnieje, zanim jakikolwiek
// dokument z niej powstanie, i bywa odrzucona, zanim powstanie jakikolwiek.
// Wiązanie jej z dokumentem wymagałoby zakładania dokumentu pustego przy każdym
// wskazaniu pliku — także tym, które skończy się odmową rozpoznania.
//
// Warstwa słów rozpoznanych i bloki układu stoją tekstem w formacie JSON. Nie są
// bytem samodzielnym: nie mają cyklu życia, nikt się do nich nie odwołuje
// z zewnątrz i giną razem z pozycją. Tabela podrzędna dałaby złączenie przy
// każdym odczycie i nie dałaby w zamian niczego.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// PozycjaWczytywania to wiersz tabeli `pozycja_wczytywania_studio`.
type PozycjaWczytywania struct {
	ID               int64
	Kod              string
	Okno             string
	SciezkaZrodlowa  *string
	ZasobID          *string
	Stan             string
	Tekst            *string
	Stron            *int64
	UzytoRozpoznania bool
	Pewnosc          *float64
	PowodOdmowy      *string
	SlowaJSON        *string
	UkladJSON        *string
	NastawyJSON      *string
	Utworzono        string
}

const (
	kolumnyPozycjiWczytywania = `id, identyfikator_zewnetrzny, okno, sciezka_zrodlowa, zasob_id,
	                             stan, tekst, stron, uzyto_rozpoznania, pewnosc, powod_odmowy,
	                             slowa_json, uklad_json, nastawy_json, utworzono`

	zapiszPozycjeWczytywania = `INSERT INTO pozycja_wczytywania_studio
	                            (identyfikator_zewnetrzny, okno, sciezka_zrodlowa, zasob_id, stan,
	                             tekst, stron, uzyto_rozpoznania, pewnosc, powod_odmowy,
	                             slowa_json, uklad_json, nastawy_json)
	                            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                            ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                stan = excluded.stan,
	                                tekst = excluded.tekst,
	                                stron = excluded.stron,
	                                uzyto_rozpoznania = excluded.uzyto_rozpoznania,
	                                pewnosc = excluded.pewnosc,
	                                powod_odmowy = excluded.powod_odmowy,
	                                slowa_json = excluded.slowa_json,
	                                uklad_json = excluded.uklad_json,
	                                nastawy_json = excluded.nastawy_json`

	pobierzPozycjeWczytywania = `SELECT ` + kolumnyPozycjiWczytywania + `
	                             FROM pozycja_wczytywania_studio
	                             WHERE identyfikator_zewnetrzny = ?`

	listaPozycjiWczytywania = `SELECT ` + kolumnyPozycjiWczytywania + `
	                           FROM pozycja_wczytywania_studio
	                           WHERE okno = ? ORDER BY id`

	listaPozycjiOczekujacych = `SELECT ` + kolumnyPozycjiWczytywania + `
	                            FROM pozycja_wczytywania_studio
	                            WHERE okno = ? AND stan IN ('oczekuje','ponowienie','przetwarzanie')
	                            ORDER BY id`
)

// ZapiszPozycjeWczytywania zakłada pozycję kolejki albo nadpisuje jej stan.
func (r *repozytoriumStudia) ZapiszPozycjeWczytywania(ctx context.Context,
	pozycja PozycjaWczytywania) (PozycjaWczytywania, error) {

	if pozycja.Kod == "" {
		return PozycjaWczytywania{}, fmt.Errorf("dane: pozycja wczytywania bez identyfikatora")
	}
	if pozycja.Stan == "" {
		pozycja.Stan = "oczekuje"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPozycjeWczytywania)
	if err != nil {
		return PozycjaWczytywania{}, err
	}
	_, err = polecenie.ExecContext(ctx, pozycja.Kod, pozycja.Okno,
		tekstDoKolumny(pozycja.SciezkaZrodlowa), tekstDoKolumny(pozycja.ZasobID), pozycja.Stan,
		tekstDoKolumny(pozycja.Tekst), liczbaDoKolumny(pozycja.Stron),
		liczbaLogiczna(pozycja.UzytoRozpoznania), liczbaRzeczywistaDoKolumny(pozycja.Pewnosc),
		tekstDoKolumny(pozycja.PowodOdmowy), tekstDoKolumny(pozycja.SlowaJSON),
		tekstDoKolumny(pozycja.UkladJSON), tekstDoKolumny(pozycja.NastawyJSON))
	if err != nil {
		return PozycjaWczytywania{}, fmt.Errorf("dane: nie można zapisać pozycji wczytywania %q: %w",
			pozycja.Kod, err)
	}
	return r.PozycjaWczytywania(ctx, pozycja.Kod)
}

// PozycjaWczytywania zwraca pozycję kolejki o wskazanym kodzie.
func (r *repozytoriumStudia) PozycjaWczytywania(ctx context.Context,
	kod string) (PozycjaWczytywania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPozycjeWczytywania)
	if err != nil {
		return PozycjaWczytywania{}, err
	}
	pozycja, err := odczytajPozycjeWczytywania(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return PozycjaWczytywania{}, ErrBrakWiersza
	}
	if err != nil {
		return PozycjaWczytywania{}, fmt.Errorf("dane: nieczytelna pozycja wczytywania %q: %w", kod, err)
	}
	return pozycja, nil
}

// PozycjeWczytywania zwraca kolejkę okna; `tylkoNieprzetworzone` zawęża wykaz.
func (r *repozytoriumStudia) PozycjeWczytywania(ctx context.Context,
	okno string, tylkoNieprzetworzone bool) ([]PozycjaWczytywania, error) {

	zapytanie := listaPozycjiWczytywania
	if tylkoNieprzetworzone {
		zapytanie = listaPozycjiOczekujacych
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kolejki wczytywania okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []PozycjaWczytywania{}
	for wiersze.Next() {
		pozycja, err := odczytajPozycjeWczytywania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelna pozycja wczytywania: %w", err)
		}
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kolejki wczytywania okna %q: %w", okno, err)
	}
	return lista, nil
}

// odczytajPozycjeWczytywania składa strukturę z jednego wiersza wyniku.
func odczytajPozycjeWczytywania(wiersz skaner) (PozycjaWczytywania, error) {
	var pozycja PozycjaWczytywania
	var sciezka, zasob, tekst, powod, slowa, uklad, nastawy sql.NullString
	var stron sql.NullInt64
	var pewnosc sql.NullFloat64
	var uzyto int
	err := wiersz.Scan(&pozycja.ID, &pozycja.Kod, &pozycja.Okno, &sciezka, &zasob, &pozycja.Stan,
		&tekst, &stron, &uzyto, &pewnosc, &powod, &slowa, &uklad, &nastawy, &pozycja.Utworzono)
	if err != nil {
		return PozycjaWczytywania{}, err
	}
	pozycja.SciezkaZrodlowa = tekstZKolumny(sciezka)
	pozycja.ZasobID = tekstZKolumny(zasob)
	pozycja.Tekst = tekstZKolumny(tekst)
	pozycja.Stron = liczbaZKolumny(stron)
	pozycja.UzytoRozpoznania = uzyto == 1
	pozycja.Pewnosc = liczbaRzeczywistaZKolumny(pewnosc)
	pozycja.PowodOdmowy = tekstZKolumny(powod)
	pozycja.SlowaJSON = tekstZKolumny(slowa)
	pozycja.UkladJSON = tekstZKolumny(uklad)
	pozycja.NastawyJSON = tekstZKolumny(nastawy)
	return pozycja, nil
}

// liczbaRzeczywistaDoKolumny przekłada wskaźnik na argument zapytania; nil daje
// NULL. Odpowiednik `liczbaDoKolumny` dla pewności rozpoznania, która jest
// ułamkiem, a nie liczbą całkowitą.
func liczbaRzeczywistaDoKolumny(wartosc *float64) any {
	if wartosc == nil {
		return nil
	}
	return *wartosc
}
