// Odpowiedzialność pliku: pytanie „które automatyki uruchamia list z tej
// skrzynki" — obsługa wyzwalacza rodzaju `mail`.
//
// Rozszerzenie idzie osobnym interfejsem, nie dopiskiem do
// `RepozytoriumAutomatyk` — z tego samego powodu, dla którego zrobiła to pętla
// wykonawcza (`automatyka_petla.go`): interfejs mieszka w `automations.go`.
// Rdzeń sięga po `RepozytoriumWyzwalaczyPoczty` asercją typu na repozytorium
// automatyk.
//
// Jedna tabela, drugie pytanie. Wiersze wyzwalaczy zapisuje wyłącznie
// `automations_harmonogram.go`; tu leży sam odczyt odwrotny — od skrzynki do
// automatyki — którego okno Scheduler nigdy nie zadaje.
package dane

import (
	"context"
	"fmt"
)

// RepozytoriumWyzwalaczyPoczty odpowiada obserwatorowi poczty na pytanie,
// które automatyki czekają na list. Rdzeń bierze je asercją typu.
type RepozytoriumWyzwalaczyPoczty interface {
	// AutomatykiNaList zwraca automatyki CZYNNE z czynnym wyzwalaczem rodzaju
	// `mail` pasującym do skrzynki: wyrażenie równe kodowi skrzynki albo `*`
	// albo puste (znaczy „każda"). Harmonogram wyłączony wyłącza i wyzwalacz.
	AutomatykiNaList(ctx context.Context, kodSkrzynki string) ([]Automatyka, error)
}

// Wyrażenie `*` i puste znaczą „każda skrzynka". Droga kontraktu zapisuje `*`,
// bo `automation.schedule.set` pomija wyzwalacze z pustym wyrażeniem.
const automatykiNaListSQL = `SELECT DISTINCT ` + kolumnyAutomatykiZPrzedrostkiem + `
	  FROM wyzwalacz_automatyki w
	  JOIN harmonogram_automatyki h ON h.id = w.harmonogram_id
	  JOIN automatyka a ON a.id = h.automatyka_id
	 WHERE w.rodzaj = 'mail' AND w.czynny = 1 AND h.czynny = 1 AND a.czynna = 1
	   AND (w.wyrazenie = ? OR w.wyrazenie = '*' OR w.wyrazenie = '')
	 ORDER BY a.id`

// kolumnyAutomatykiZPrzedrostkiem powtarza porządek `kolumnyAutomatyki`
// (automations.go) z przedrostkiem złączenia — wiersz czyta to samo
// `odczytajAutomatyke`, więc porządek kolumn ma jedną prawdę odczytu.
const kolumnyAutomatykiZPrzedrostkiem = `a.id, a.identyfikator_zewnetrzny, a.nazwa,
	a.opis, a.czynna, a.wersja, a.utworzono, a.zaktualizowano`

// AutomatykiNaList — patrz RepozytoriumWyzwalaczyPoczty.
func (r *repozytoriumAutomatyk) AutomatykiNaList(ctx context.Context,
	kodSkrzynki string) ([]Automatyka, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, automatykiNaListSQL)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodSkrzynki)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać automatyk wyzwalanych listem: %w", err)
	}
	defer wiersze.Close()
	lista := []Automatyka{}
	for wiersze.Next() {
		automatyka, err := odczytajAutomatyke(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz automatyki wyzwalanej listem: %w", err)
		}
		lista = append(lista, automatyka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt automatyk wyzwalanych listem: %w", err)
	}
	return lista, nil
}
