// Odpowiedzialność pliku: archiwum eksperta — odłożenie go poza wykaz czynnych
// bez utraty definicji i historii, powrót z archiwum oraz wykaz zarchiwizowanych
// (kolumny `agent.zarchiwizowano_o` i `agent.aktywny_przed_archiwum`). Historia
// wersji leży w `agent_wersje.go` — to ta sama implementacja repozytorium,
// rozdzielona na dwa pliki.
//
// Archiwizacja nie jest usunięciem i go nie zastępuje. `agent.delete` zostaje
// tym, czym był: jedyną drogą utraty definicji eksperta (wzorzec kosza sesji).
// Archiwum jest drogą drugą — odwracalną i niczego nietracącą.
// Wybór między nimi należy do Operatora, nie do rdzenia, więc rdzeń nie
// przekierowuje jednej komendy na drugą.
//
// Stan czynności wraca taki, jaki był. Archiwizacja zapamiętuje `aktywny`
// w `aktywny_przed_archiwum` i gasi ekspertowi czynność, żeby zniknął z wykazu
// `enabledOnly`. Powrót oddaje zapamiętaną wartość zamiast zakładać, że każdy
// zarchiwizowany był czynny — ekspert wyłączony przed archiwizacją wróciłby
// inaczej, niż go odkładano.
//
// Migawki archiwizacja nie mnoży. Zapis nie rusza licznika `agent.wersja`, więc
// wyzwalacz bazy uzupełnia migawkę wersji bieżącej zamiast zakładać nową.
// Odłożenie eksperta na półkę nie jest zmianą jego tożsamości.
package dane

import (
	"context"
	"fmt"
)

// RepozytoriumArchiwumAgentow jest kontraktem archiwum biblioteki ekspertów.
// Wypełnia go ta sama implementacja co historię wersji.
type RepozytoriumArchiwumAgentow interface {
	// Zarchiwizuj odkłada eksperta poza wykaz czynnych. Wynik `false` znaczy
	// „nie było czego archiwizować albo już leżał w archiwum”.
	Zarchiwizuj(ctx context.Context, kodAgenta string) (bool, error)
	// PrzywrocZArchiwum oddaje eksperta wykazowi czynnych wraz z zapamiętanym
	// stanem czynności.
	PrzywrocZArchiwum(ctx context.Context, kodAgenta string) (bool, error)
	// Archiwum oddaje ekspertów odłożonych, od odłożonego najpóźniej.
	Archiwum(ctx context.Context) ([]Agent, error)
}

const (
	zarchiwizujAgenta = `UPDATE agent
	                     SET zarchiwizowano_o = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
	                         aktywny_przed_archiwum = aktywny,
	                         aktywny = 0
	                     WHERE kod = ? AND zarchiwizowano_o IS NULL`

	przywrocAgentaZArchiwum = `UPDATE agent
	                           SET zarchiwizowano_o = NULL,
	                               aktywny = COALESCE(aktywny_przed_archiwum, 1),
	                               aktywny_przed_archiwum = NULL
	                           WHERE kod = ? AND zarchiwizowano_o IS NOT NULL`

	listaArchiwumAgentow = `SELECT ` + kolumnyAgenta + ` FROM agent
	                        WHERE zarchiwizowano_o IS NOT NULL
	                        ORDER BY zarchiwizowano_o DESC, nazwa`
)

// Zarchiwizuj odkłada eksperta poza wykaz czynnych. Definicja, powiązania
// i cała historia wersji zostają nietknięte — nie kasujemy ani jednego wiersza.
func (r *repozytoriumWersjiAgenta) Zarchiwizuj(ctx context.Context, kodAgenta string) (bool, error) {
	return r.przestawArchiwum(ctx, zarchiwizujAgenta, kodAgenta, "zarchiwizować")
}

// PrzywrocZArchiwum oddaje eksperta wykazowi czynnych. Wynik `false` znaczy
// „nie leżał w archiwum”, co nie jest odmową.
func (r *repozytoriumWersjiAgenta) PrzywrocZArchiwum(ctx context.Context, kodAgenta string) (bool, error) {
	return r.przestawArchiwum(ctx, przywrocAgentaZArchiwum, kodAgenta, "przywrócić z archiwum")
}

// przestawArchiwum wykonuje jeden z dwóch zapisów znacznika i oddaje informację,
// czy trafił w wiersz. Oba zapisy różnią się wyłącznie treścią polecenia, więc
// drugiego obiegu błędów dla nich nie zakładamy.
func (r *repozytoriumWersjiAgenta) przestawArchiwum(ctx context.Context, zapytanie, kodAgenta,
	czynnosc string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kodAgenta)
	if err != nil {
		return false, fmt.Errorf("dane: nie można %s eksperta %q: %w", czynnosc, kodAgenta, err)
	}
	liczba, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany wynik archiwizacji eksperta %q: %w", kodAgenta, err)
	}
	return liczba > 0, nil
}

// Archiwum oddaje ekspertów odłożonych. Wiersz niesie samą definicję bez
// powiązań: okno archiwum wybiera eksperta do przywrócenia po nazwie i dacie
// odłożenia, a komplet umiejętności, konektorów i uprawnień dostaje z
// `agent.list` w chwili, gdy ekspert wróci między czynnych. Dobieranie ich tutaj
// byłoby trzema zapytaniami na widok, którego nikt o nie nie pyta.
func (r *repozytoriumWersjiAgenta) Archiwum(ctx context.Context) ([]Agent, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaArchiwumAgentow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać archiwum ekspertów: %w", err)
	}
	defer wiersze.Close()

	odlozeni := []Agent{}
	for wiersze.Next() {
		agent, err := odczytajAgenta(wiersze)
		if err != nil {
			return nil, err
		}
		odlozeni = append(odlozeni, agent)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt archiwum ekspertów: %w", err)
	}
	return odlozeni, nil
}
