// Dostęp do obszaru bloków wiadomości (tabela `blok_wiadomosci`) —
// nietekstowych fragmentów strumienia odpowiedzi, zapisywanych w trakcie tury:
// toku rozumowania, wywołań narzędzi z wynikami, prowenancji i metadanych
// konta. Odczyt dokleja je do wiadomości kontraktu w `rozmowa_bloki.go`.
//
// Tekstu tutaj nie ma. Fragment rodzaju 'text' domyka dziennik rozmowy
// w kolumnie `wiadomosc.tresc`; zapisanie go drugi raz tutaj byłoby drugą
// prawdą o tej samej wypowiedzi, więc warunek CHECK schematu odbija taki zapis,
// a rejestrator bloków nawet go nie próbuje.
//
// Rodzaj jest wartością kontraktu (ChunkKind, angielską): kontrakt nie zapisuje
// polskiego odwzorowania dla 'provenance' i 'account', a przekład należy
// wyłącznie do kontraktu.
package dane

import (
	"context"
	"encoding/json"
	"fmt"

	"danacoconsole/shared"
)

// BlokWiadomosci to wiersz tabeli `blok_wiadomosci`.
type BlokWiadomosci struct {
	ID int64
	// Chwila zapisu w milisekundach epoki — czas nadejścia fragmentu, nie zapisu.
	Chwila int64
	// OknoKod i WiadomoscKod są identyfikatorami kontraktowymi (napisy), nie
	// kluczami obcymi — rejestrator strumienia innych nie zna.
	OknoKod      string
	WiadomoscKod string
	Kolejnosc    int
	Rodzaj       shared.ChunkKind
	// Tresc niesie pole `text` fragmentu (tak nadchodzi tok rozumowania),
	// Ladunek — surowe pole `data` (dowód pierwotny); nil znaczy brak ładunku.
	Tresc   string
	Ladunek json.RawMessage
}

// RepozytoriumBlokow jest kontraktem obszaru bloków wiadomości.
type RepozytoriumBlokow interface {
	Dopisz(ctx context.Context, blok BlokWiadomosci) error
	ListaOkna(ctx context.Context, oknoKod string) ([]BlokWiadomosci, error)
}

const (
	kolumnyBloku = `id, chwila, okno_kod, wiadomosc_kod, kolejnosc, rodzaj, tresc, ladunek`

	// Kolejność w obrębie wiadomości wyliczana w tym samym poleceniu — wzorem
	// numeracji historii okna (wiadomosci.go): bez osobnego odczytu i transakcji.
	wstawBlok = `INSERT INTO blok_wiadomosci
	             (chwila, okno_kod, wiadomosc_kod, kolejnosc, rodzaj, tresc, ladunek)
	             VALUES (?, ?, ?,
	                     (SELECT COALESCE(MAX(kolejnosc) + 1, 0) FROM blok_wiadomosci
	                      WHERE wiadomosc_kod = ?),
	                     ?, ?, ?)`

	listaBlokowOkna = `SELECT ` + kolumnyBloku + ` FROM blok_wiadomosci
	                   WHERE okno_kod = ?
	                   ORDER BY wiadomosc_kod, kolejnosc, id`
)

type repozytoriumBlokow struct {
	zapytania *zapytania
}

func noweRepozytoriumBlokow(z *zapytania) *repozytoriumBlokow {
	return &repozytoriumBlokow{zapytania: z}
}

// Dopisz dokłada blok na koniec zapisu wiadomości. Wywoływany w trakcie tury,
// fragment po fragmencie — niepowodzenie zgłasza wywołującemu, a ten degraduje
// zapis okna zamiast przerywać rozmowę.
func (r *repozytoriumBlokow) Dopisz(ctx context.Context, blok BlokWiadomosci) error {
	if blok.WiadomoscKod == "" {
		return fmt.Errorf("dane: nie można utrwalić bloku bez identyfikatora wiadomości")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawBlok)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, blok.Chwila, blok.OknoKod, blok.WiadomoscKod,
		blok.WiadomoscKod, string(blok.Rodzaj), blok.Tresc, ladunekDoKolumny(blok.Ladunek))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać bloku wiadomości %q: %w", blok.WiadomoscKod, err)
	}
	return nil
}

// ListaOkna zwraca wszystkie bloki okna, pogrupowane wiadomościami w kolejności
// zapisu. Puste okno daje pustą listę — brak nie jest błędem.
func (r *repozytoriumBlokow) ListaOkna(ctx context.Context, oknoKod string) ([]BlokWiadomosci, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaBlokowOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKod)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać bloków okna %q: %w", oknoKod, err)
	}
	defer wiersze.Close()

	lista := []BlokWiadomosci{}
	for wiersze.Next() {
		blok, err := odczytajBlok(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, blok)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt bloków okna %q: %w", oknoKod, err)
	}
	return lista, nil
}

// odczytajBlok składa strukturę z jednego wiersza wyniku.
func odczytajBlok(wiersz skaner) (BlokWiadomosci, error) {
	var blok BlokWiadomosci
	var rodzaj string
	var ladunek *string
	err := wiersz.Scan(&blok.ID, &blok.Chwila, &blok.OknoKod, &blok.WiadomoscKod,
		&blok.Kolejnosc, &rodzaj, &blok.Tresc, &ladunek)
	if err != nil {
		return BlokWiadomosci{}, fmt.Errorf("dane: nieczytelny wiersz bloku wiadomości: %w", err)
	}
	blok.Rodzaj = shared.ChunkKind(rodzaj)
	if ladunek != nil {
		blok.Ladunek = json.RawMessage(*ladunek)
	}
	return blok, nil
}

// ladunekDoKolumny przekłada surowy ładunek na wartość kolumny; pusty daje NULL.
func ladunekDoKolumny(ladunek json.RawMessage) any {
	if len(ladunek) == 0 {
		return nil
	}
	return string(ladunek)
}
