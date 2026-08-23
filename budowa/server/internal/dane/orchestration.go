// Odpowiedzialność pliku: pojedynczy łuk układu zależności — zapis jednej
// zależności i usunięcie jednej zależności w tabeli `zaleznosc_kroku_automatyki`.
//
// `automations_kroki.go` podmienia cały układ, bo Workflow Builder oddaje pełną
// definicję automatyki. Tu operacja obejmuje jeden łuk: podmiana kompletu przy
// dołożeniu pojedynczego łuku przepisywałaby pozostałe wiersze, gubiąc ich
// `utworzono`, a dwa okna pracujące równocześnie kasowałyby sobie zmiany
// nawzajem. Obie drogi sięgają tych samych wierszy tej samej tabeli.
package dane

import (
	"context"
	"fmt"
)

const (
	// Zapis jednego łuku. Klucz główny (automatyka_id, krok_z, krok_do) sprawia,
	// że powtórzone żądanie zmienia rodzaj i warunek zamiast wywracać się na
	// więzie — Operator przestawiający łuk z sekwencyjnego na warunkowy wykonuje
	// tę samą komendę drugi raz.
	wstawJedenLukUkladu = `INSERT INTO zaleznosc_kroku_automatyki
	                       (automatyka_id, krok_z, krok_do, rodzaj, warunek)
	                       VALUES (?, ?, ?, ?, ?)
	                       ON CONFLICT(automatyka_id, krok_z, krok_do) DO UPDATE SET
	                           rodzaj = excluded.rodzaj,
	                           warunek = excluded.warunek`

	usunJedenLukUkladu = `DELETE FROM zaleznosc_kroku_automatyki
	                      WHERE automatyka_id = ? AND krok_z = ? AND krok_do = ?`
)

// ZapiszZaleznosc zapisuje jeden łuk układu zależności; łuk zastany zmienia,
// nowego dokłada. Pozostałe łuki automatyki zostają nietknięte.
//
// Repozytorium nie ocenia układu: łuk do kroku nieistniejącego i łuk domykający
// cykl zapisują się tak samo jak każdy inny — ocena należy do walidacji w rdzeniu.
// Więzy pilnowane przez sam schemat, czyli pętla własna i rodzaj spoza wartości
// kontraktu, wracają stąd jako błąd zapisu.
func (r *repozytoriumAutomatyk) ZapiszZaleznosc(ctx context.Context,
	automatykaID int64, zaleznosc ZaleznoscKroku) error {

	if zaleznosc.KrokZ == "" || zaleznosc.KrokDo == "" {
		return fmt.Errorf("dane: zależność automatyki %d bez wskazania kroków", automatykaID)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawJedenLukUkladu)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, automatykaID, zaleznosc.KrokZ, zaleznosc.KrokDo,
		zaleznosc.Rodzaj, tekstDoKolumny(zaleznosc.Warunek))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać zależności %q→%q automatyki %d: %w",
			zaleznosc.KrokZ, zaleznosc.KrokDo, automatykaID, err)
	}
	return nil
}

// UsunZaleznosc zdejmuje jeden łuk układu. Wynik `false` znaczy, że łuku
// w układzie nie było.
func (r *repozytoriumAutomatyk) UsunZaleznosc(ctx context.Context,
	automatykaID int64, krokZ, krokDo string) (bool, error) {

	if krokZ == "" || krokDo == "" {
		return false, fmt.Errorf("dane: usunięcie zależności automatyki %d bez wskazania kroków", automatykaID)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, usunJedenLukUkladu)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, automatykaID, krokZ, krokDo)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć zależności %q→%q automatyki %d: %w",
			krokZ, krokDo, automatykaID, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia zależności %q→%q automatyki %d: %w",
			krokZ, krokDo, automatykaID, err)
	}
	return zdjete > 0, nil
}
