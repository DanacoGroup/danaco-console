// Plik prowadzi dostęp do obszaru bloków wiadomości: nietekstowych fragmentów strumienia odpowiedzi zapisywanych
// w trakcie tury; odczyt dokleja je do wiadomości kontraktu, a rodzaj jest wartością kontraktu wprost, bez polskiego przekładu.
package dane

import (
	"context"
	"encoding/json"
	"fmt"

	"danacoconsole/shared"
)

// BlokWiadomosci to wiersz tabeli `blok_wiadomosci` niosący jeden nietekstowy fragment strumienia odpowiedzi.
type BlokWiadomosci struct {
	ID int64
	// Chwila zapisu w milisekundach epoki — czas nadejścia fragmentu, nie zapisu.
	Chwila int64
	// OknoKod i WiadomoscKod są identyfikatorami kontraktowymi, nie kluczami obcymi.
	OknoKod      string
	WiadomoscKod string
	Kolejnosc    int
	Rodzaj       shared.ChunkKind
	// Tresc niesie pole text fragmentu, Ladunek niesie surowe pole data; pusty ładunek znaczy jego brak.
	Tresc   string
	Ladunek json.RawMessage
}

// RepozytoriumBlokow jest kontraktem obszaru bloków wiadomości: zapis i odczyt fragmentów strumienia odpowiedzi.
type RepozytoriumBlokow interface {
	Dopisz(ctx context.Context, blok BlokWiadomosci) error
	ListaOkna(ctx context.Context, oknoKod string) ([]BlokWiadomosci, error)
}

const (
	kolumnyBloku = `id, chwila, okno_kod, wiadomosc_kod, kolejnosc, rodzaj, tresc, ladunek`

	// Kolejność w obrębie wiadomości wyliczana w tym samym poleceniu — wzorem
	// numeracji historii okna (wiadomosci.go): bez osobnego odczytu i transakcji.
	// Numeracja liczy wyłącznie bloki konta żądania, bo tylko one wracają odczytem.
	wstawBlok = `INSERT INTO blok_wiadomosci
	             (chwila, okno_kod, wiadomosc_kod, kolejnosc, rodzaj, tresc, ladunek, konto_id)
	             VALUES (?, ?, ?,
	                     (SELECT COALESCE(MAX(kolejnosc) + 1, 0) FROM blok_wiadomosci
	                      WHERE wiadomosc_kod = ? AND ` + WarunekKonta + `),
	                     ?, ?, ?, ` + WskazanieKonta + `)`

	listaBlokowOkna = `SELECT ` + kolumnyBloku + ` FROM blok_wiadomosci
	                   WHERE okno_kod = ? AND ` + WarunekKonta + `
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
		blok.WiadomoscKod, KontoOperatora(ctx), string(blok.Rodzaj), blok.Tresc,
		ladunekDoKolumny(blok.Ladunek), KontoOperatora(ctx))
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
	wiersze, err := polecenie.QueryContext(ctx, oknoKod, KontoOperatora(ctx))
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

// odczytajBlok składa strukturę bloku wiadomości wprost z jednego wiersza wyniku zapytania do bazy SQL.
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

// ladunekDoKolumny przekłada surowy ładunek na wartość kolumny bazy danych; ładunek pusty daje wartość NULL.
func ladunekDoKolumny(ladunek json.RawMessage) any {
	if len(ladunek) == 0 {
		return nil
	}
	return string(ladunek)
}
