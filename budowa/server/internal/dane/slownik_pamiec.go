// Odpowiedzialność pliku: pamięć tłumaczeń (Translation Memory, tabela
// `pamiec_tlumaczen`) modułu Translate. Zasila `translate.memory.suggest`.
//
// Pamięć gromadzi się z zatwierdzonych par segmentów, nie z odczytu paneli na
// żywo. `ZapiszPamiec` dokłada wiersz — nie czyta bieżącej treści panelu
// w edycji; to wywołujący rozstrzyga, kiedy para segmentów jest zatwierdzona
// i warta zapamiętania.
//
// Dopasowanie jest przybliżone tylko na tyle, na ile pozwala baza. Kontrakt
// (`TranslateMemorySuggestRequest`) chce dopasowania po podobieństwie, nie po
// równości, a SQLite bez rozszerzenia (FTS5 albo funkcji trigramowej) nie ma
// wbudowanej miary podobieństwa napisów. `Podpowiedzi` wykonuje więc `LIKE`
// z frazą otoczoną znakami `%`, czyli dopasowanie podciągu segmentu źródłowego
// — nie dopasowanie znaczeniowe ani odległość edycyjną.
//
// Czas jest liczbą (ms epoki), wzorem `dane/asystent.go` i `dane/tlumaczenie.go`.
package dane

import (
	"context"
	"fmt"
	"time"
)

// WpisPamieciTlumaczen to wiersz tabeli `pamiec_tlumaczen` — zatwierdzona
// para (segment źródłowy, segment docelowy) w obrębie języka i panelu, z
// którego pochodzi zatwierdzenie.
type WpisPamieciTlumaczen struct {
	ID              int64
	Kod             string
	PanelID         int64
	Jezyk           string
	SegmentZrodlowy string
	SegmentDocelowy string
	Utworzono       int64
}

const (
	kolumnyPamieciTlumaczen = `id, identyfikator_zewnetrzny, panel_id, jezyk,
	                           segment_zrodlowy, segment_docelowy, utworzono`

	zapiszPamiecTlumaczen = `INSERT INTO pamiec_tlumaczen
	                         (identyfikator_zewnetrzny, panel_id, jezyk,
	                          segment_zrodlowy, segment_docelowy, utworzono)
	                         VALUES (?, ?, ?, ?, ?, ?)`

	// Zawężenie po (jezyk, segment_zrodlowy LIKE ?) korzysta z indeksu
	// idx_pamiec_tlumaczen_jezyk_segment na kolumnie jezyk —
	// SQLite wykorzystuje prefiks indeksu złożonego nawet gdy druga kolumna
	// idzie przez LIKE, bo warunek na jezyk jest równością. Najnowsze
	// zatwierdzenia na przodzie — świeższa podpowiedź jest zwykle trafniejsza.
	podpowiedziPamieciTlumaczen = `SELECT ` + kolumnyPamieciTlumaczen + ` FROM pamiec_tlumaczen
	                               WHERE jezyk = ? AND segment_zrodlowy LIKE ?
	                               ORDER BY utworzono DESC, id DESC LIMIT ?`
)

// ZapiszPamiec dokłada wiersz pamięci tłumaczeń dla zatwierdzonej pary
// segmentów. Nie nadpisuje istniejących wpisów (brak ON CONFLICT) — TM z
// definicji gromadzi historię zatwierdzeń, kolejne zatwierdzenie tego samego
// segmentu to nowy fakt, nie korekta poprzedniego.
func (r *repozytoriumTlumaczen) ZapiszPamiec(ctx context.Context, panelID int64, wpis WpisPamieciTlumaczen) (WpisPamieciTlumaczen, error) {
	if wpis.Kod == "" || panelID == 0 {
		return WpisPamieciTlumaczen{}, fmt.Errorf("dane: wpis pamięci tłumaczeń bez identyfikatora albo bez panelu")
	}
	if wpis.Jezyk == "" || wpis.SegmentZrodlowy == "" {
		return WpisPamieciTlumaczen{}, fmt.Errorf("dane: wpis pamięci tłumaczeń %q bez języka albo segmentu źródłowego", wpis.Kod)
	}
	utworzono := wpis.Utworzono
	if utworzono == 0 {
		utworzono = time.Now().UnixMilli()
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPamiecTlumaczen)
	if err != nil {
		return WpisPamieciTlumaczen{}, err
	}
	_, err = polecenie.ExecContext(ctx, wpis.Kod, panelID, wpis.Jezyk,
		wpis.SegmentZrodlowy, wpis.SegmentDocelowy, utworzono)
	if err != nil {
		return WpisPamieciTlumaczen{}, fmt.Errorf("dane: nie można zapisać wpisu pamięci tłumaczeń %q: %w", wpis.Kod, err)
	}
	return WpisPamieciTlumaczen{
		Kod: wpis.Kod, PanelID: panelID, Jezyk: wpis.Jezyk,
		SegmentZrodlowy: wpis.SegmentZrodlowy, SegmentDocelowy: wpis.SegmentDocelowy,
		Utworzono: utworzono,
	}, nil
}

// Podpowiedzi szuka wpisów pamięci tłumaczeń danego języka, których segment
// źródłowy zawiera frazę — patrz uwaga u góry pliku o ograniczeniu dopasowania
// przybliżonego do LIKE po podciągu. Pusta lista (nie błąd) oznacza brak
// dopasowań — `translate.memory.suggest` ma zwrócić pustą listę podpowiedzi,
// nie ErrBrakWiersza, bo brak podpowiedzi nie jest usterką zapytania.
func (r *repozytoriumTlumaczen) Podpowiedzi(ctx context.Context, jezyk, fraza string, limit int) ([]WpisPamieciTlumaczen, error) {
	if jezyk == "" {
		return nil, fmt.Errorf("dane: podpowiedzi pamięci tłumaczeń bez języka")
	}
	if limit <= 0 {
		limit = 5
	}
	polecenie, err := r.zapytania.przygotuj(ctx, podpowiedziPamieciTlumaczen)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, jezyk, "%"+fraza+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać podpowiedzi pamięci tłumaczeń dla języka %q: %w", jezyk, err)
	}
	defer wiersze.Close()

	lista := []WpisPamieciTlumaczen{}
	for wiersze.Next() {
		wpis, err := odczytajWpisPamieciTlumaczen(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pamięci tłumaczeń dla języka %q: %w", jezyk, err)
		}
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt podpowiedzi pamięci tłumaczeń dla języka %q: %w", jezyk, err)
	}
	return lista, nil
}

// odczytajWpisPamieciTlumaczen składa strukturę z jednego wiersza wyniku.
func odczytajWpisPamieciTlumaczen(wiersz skaner) (WpisPamieciTlumaczen, error) {
	var wpis WpisPamieciTlumaczen
	err := wiersz.Scan(&wpis.ID, &wpis.Kod, &wpis.PanelID, &wpis.Jezyk,
		&wpis.SegmentZrodlowy, &wpis.SegmentDocelowy, &wpis.Utworzono)
	if err != nil {
		return WpisPamieciTlumaczen{}, err
	}
	return wpis, nil
}
