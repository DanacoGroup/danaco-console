// Odpowiedzialność pliku: szukanie w treści rozmów — odczyt indeksu
// pełnotekstowego `wiadomosc_szukanie` (FTS5).
//
// Szukanie jedzie przez MATCH, snippet i rank, nie przez LIKE. Tokenizator
// indeksu sprowadza diakrytyki poza „ł" do liter podstawowych („zolta" trafi
// „żółta"), ale „ł" pozostaje osobną literą — „lodz" nie trafi w „łódź".
//
// Fraza jedzie w cudzysłowie: zapytanie Operatora jest frazą, nie wyrażeniem
// składni FTS5 — ujęcie w cudzysłów (z podwojeniem cudzysłowów wewnętrznych)
// zdejmuje z niego operatory AND/OR/NOT/*, żeby fraza z myślnikiem albo
// gwiazdką nie wywracała zapytania błędem składni.
package dane

import (
	"context"
	"fmt"
	"strings"
)

// TrafienieRozmowy to jeden wiersz wyniku szukania: wskazanie sesji i okna
// wraz z wycinkiem treści wokół trafienia.
type TrafienieRozmowy struct {
	// SesjaKod i OknoKod są identyfikatorami rdzenia (kontraktowymi); puste,
	// gdy wiersz powstał poza rdzeniem — trafienie zostaje, wskazanie nie
	// prowadzi wtedy do bytu rejestru.
	SesjaKod     string
	TytulSesji   string
	OknoKod      string
	WiadomoscKod string
	Rola         string
	// Wycinek treści wokół trafienia; słowa trafione ujęte w «…».
	Wycinek   string
	Utworzono string
}

// RepozytoriumSzukaniaRozmow jest kontraktem obszaru szukania.
type RepozytoriumSzukaniaRozmow interface {
	SzukajWiadomosci(ctx context.Context, fraza string, granica int) ([]TrafienieRozmowy, error)
}

// granicaSzukaniaDomyslna ogranicza wynik, gdy wywołujący nie poda własnej
// granicy — wykaz trafień jest podpowiedzią, nie zrzutem bazy.
const granicaSzukaniaDomyslna = 50

// Sesje w koszu nie wracają w trafieniach — ta sama zasada, która kryje je
// w wykazie sesji (dane/sesje.go).
const szukajWiadomosci = `SELECT COALESCE(s.identyfikator_zewnetrzny, ''), s.tytul,
                                 COALESCE(o.identyfikator_zewnetrzny, ''),
                                 COALESCE(w.identyfikator_zewnetrzny, ''),
                                 w.rola,
                                 snippet(wiadomosc_szukanie, 0, '«', '»', '…', 12),
                                 w.utworzono
                          FROM wiadomosc_szukanie
                          JOIN wiadomosc w         ON w.id = wiadomosc_szukanie.rowid
                          JOIN okno_komunikacji o  ON o.id = w.okno_komunikacji_id
                          JOIN sesja s             ON s.id = o.sesja_id
                          WHERE wiadomosc_szukanie MATCH ?
                            AND s.usunieto_o IS NULL
                          ORDER BY rank
                          LIMIT ?`

type repozytoriumSzukaniaRozmow struct {
	zapytania *zapytania
}

func noweRepozytoriumSzukaniaRozmow(z *zapytania) *repozytoriumSzukaniaRozmow {
	return &repozytoriumSzukaniaRozmow{zapytania: z}
}

// SzukajWiadomosci zwraca trafienia frazy w treści wiadomości, najlepsze
// najpierw (rank FTS5). Fraza pusta daje pusty wykaz — nie ma czego szukać
// i nie jest to błąd.
func (r *repozytoriumSzukaniaRozmow) SzukajWiadomosci(ctx context.Context, fraza string, granica int) ([]TrafienieRozmowy, error) {
	fraza = strings.TrimSpace(fraza)
	if fraza == "" {
		return []TrafienieRozmowy{}, nil
	}
	if granica <= 0 {
		granica = granicaSzukaniaDomyslna
	}
	polecenie, err := r.zapytania.przygotuj(ctx, szukajWiadomosci)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, frazaFTS(fraza), granica)
	if err != nil {
		return nil, fmt.Errorf("dane: szukanie %q nie powiodło się: %w", fraza, err)
	}
	defer wiersze.Close()

	trafienia := []TrafienieRozmowy{}
	for wiersze.Next() {
		var t TrafienieRozmowy
		if err := wiersze.Scan(&t.SesjaKod, &t.TytulSesji, &t.OknoKod,
			&t.WiadomoscKod, &t.Rola, &t.Wycinek, &t.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz trafienia: %w", err)
		}
		trafienia = append(trafienia, t)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt trafień %q: %w", fraza, err)
	}
	return trafienia, nil
}

// frazaFTS ujmuje zapytanie Operatora w cudzysłów składni FTS5 — fraza
// dosłowna, bez operatorów; cudzysłowy wewnętrzne podwojone.
func frazaFTS(fraza string) string {
	return `"` + strings.ReplaceAll(fraza, `"`, `""`) + `"`
}
