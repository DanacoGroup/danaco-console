// Plik utrwala trwały ślad przejęć i oddań bezpośredniego sterowania
// zleceniem, zapisywany jako akcja w istniejącym dzienniku akcji okna, bez
// osobnej tabeli.
package dane

import (
	"context"
	"fmt"
)

// Stałe nazywają rodzaje wpisu dziennika akcji, którymi znaczy się przejęcie
// sterowania zleceniem, i brzmią jak komendy, które ten ślad wytwarzają.
const (
	// AkcjaPrzejeciaSterowania nazywa wpis dziennika zakładany, gdy sterowanie
	// zleceniem przejmuje osoba obsługująca sesję.
	AkcjaPrzejeciaSterowania = "control.takeover"
	// AkcjaOddaniaSterowania nazywa wpis dziennika zakładany, gdy sterowanie
	// zleceniem wraca do automatycznego prowadzenia.
	AkcjaOddaniaSterowania = "control.release"
)

// RepozytoriumSteru jest wąskim kontraktem odczytu śladu sterowania, osobnym
// od szerokiego kontraktu obszaru okien, choć stoi na tym samym typie
// repozytorium.
type RepozytoriumSteru interface {
	// SladySterowania zwraca przejęcia i oddania sterowania oknem, od
	// najnowszego, wykaz pusty bez błędu.
	SladySterowania(ctx context.Context, okno string, limit int) ([]AkcjaOkna, error)
}

// Typ repozytorium obszaru window.* wypełnia ten kontrakt. Gdyby metoda znikła
// albo zmieniła kształt, kompilacja stanie tutaj, a nie na martwej komendzie.
var _ RepozytoriumSteru = (*repozytoriumPrzekazan)(nil)

const pobierzSladySterowania = `SELECT log.id, ok.identyfikator_zewnetrzny, log.akcja_id, log.parametry,
                                       log.wynik, log.utworzono
                                FROM log_akcji_okna log
                                JOIN okno_komunikacji ok ON ok.id = log.okno_komunikacji_id
                                WHERE ok.identyfikator_zewnetrzny = ?
                                  AND log.akcja_id IN (?, ?)
                                ORDER BY log.utworzono DESC, log.id DESC
                                LIMIT ?`

// SladySterowania oddaje historię przejęć i oddań sterowania nad oknem.
//
// Limit niedodatni schodzi na domyślny — ten sam, co w dzienniku akcji. Odczyt
// bez granicy byłby pułapką wydajności przy oknie prowadzonym tygodniami.
func (r *repozytoriumPrzekazan) SladySterowania(ctx context.Context, okno string, limit int) ([]AkcjaOkna, error) {
	if okno == "" {
		return nil, fmt.Errorf("dane: ślad sterowania bez identyfikatora okna")
	}
	if limit <= 0 {
		limit = limitDomyslnyAkcjiOkna
	}

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSladySterowania)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno,
		AkcjaPrzejeciaSterowania, AkcjaOddaniaSterowania, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać śladu sterowania okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []AkcjaOkna{}
	for wiersze.Next() {
		akcja, err := odczytajAkcjeOkna(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny ślad sterowania okna %q: %w", okno, err)
		}
		lista = append(lista, akcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt śladu sterowania okna %q: %w", okno, err)
	}
	return lista, nil
}
