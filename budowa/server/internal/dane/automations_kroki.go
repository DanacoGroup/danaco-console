// Plik prowadzi kroki automatyki i układ zależności między nimi: trwałość okien Workflow Builder i Orchestrator;
// zapis jest wymianą, nie dokładaniem, więc kroki i zależności podmieniają zestaw w jednej transakcji.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// KrokAutomatyki to wiersz tabeli `krok_automatyki`. Rodzaj i parametry są
// wartościami kontraktu (AutomationStepKind, AutomationStep.params).
type KrokAutomatyki struct {
	ID        int64
	Kod       string
	Nazwa     *string
	Rodzaj    string
	Komenda   *string
	Parametry *string
	Warunek   *string
	Kolejnosc int
	// OdwolaniaSekretow to wykaz referencji skarbca; znika razem z krokiem, który przestał go wołać.
	OdwolaniaSekretow *string
}

// ZaleznoscKroku to jeden łuk układu zależności — jedyne miejsce zapisu
// kolejności wykonania (pole `dependsOn` kontraktu powstaje z tych wierszy).
type ZaleznoscKroku struct {
	KrokZ   string
	KrokDo  string
	Rodzaj  string
	Warunek *string
}

const (
	usunKrokiAutomatyki = `DELETE FROM krok_automatyki WHERE automatyka_id = ?`

	wstawKrokAutomatyki = `INSERT INTO krok_automatyki
	                       (automatyka_id, identyfikator_zewnetrzny, nazwa, rodzaj,
	                        komenda, parametry, warunek, kolejnosc, odwolania_sekretow)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	listaKrokowAutomatyki = `SELECT id, identyfikator_zewnetrzny, nazwa, rodzaj, komenda,
	                                parametry, warunek, kolejnosc, odwolania_sekretow
	                         FROM krok_automatyki WHERE automatyka_id = ?
	                         ORDER BY kolejnosc, id`

	usunZaleznosciAutomatyki = `DELETE FROM zaleznosc_kroku_automatyki WHERE automatyka_id = ?`

	wstawZaleznoscKroku = `INSERT INTO zaleznosc_kroku_automatyki
	                       (automatyka_id, krok_z, krok_do, rodzaj, warunek)
	                       VALUES (?, ?, ?, ?, ?)
	                       ON CONFLICT(automatyka_id, krok_z, krok_do) DO UPDATE SET
	                           rodzaj = excluded.rodzaj,
	                           warunek = excluded.warunek`

	listaZaleznosciAutomatyki = `SELECT krok_z, krok_do, rodzaj, warunek
	                             FROM zaleznosc_kroku_automatyki WHERE automatyka_id = ?
	                             ORDER BY krok_do, krok_z`
)

// ZapiszKroki podmienia komplet kroków automatyki. Wykaz pusty zostawia
// automatykę bez kroków — definicja szkicowa jest stanem poprawnym.
func (r *repozytoriumAutomatyk) ZapiszKroki(ctx context.Context,
	automatykaID int64, kroki []KrokAutomatyki) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunKrokiAutomatyki)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, automatykaID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić kroków automatyki %d: %w", automatykaID, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawKrokAutomatyki)
		if err != nil {
			return err
		}
		for numer, krok := range kroki {
			kolejnosc := krok.Kolejnosc
			if kolejnosc == 0 {
				kolejnosc = numer + 1
			}
			_, err := wstawienie.ExecContext(ctx, automatykaID, krok.Kod,
				tekstDoKolumny(krok.Nazwa), krok.Rodzaj, tekstDoKolumny(krok.Komenda),
				tekstDoKolumny(krok.Parametry), tekstDoKolumny(krok.Warunek), kolejnosc,
				tekstDoKolumny(krok.OdwolaniaSekretow))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać kroku %q automatyki %d: %w",
					krok.Kod, automatykaID, err)
			}
		}
		return nil
	})
}

// Kroki zwraca wszystkie kroki automatyki w zapisanej kolejności ich wykonania z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) Kroki(ctx context.Context, automatykaID int64) ([]KrokAutomatyki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKrokowAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kroków automatyki %d: %w", automatykaID, err)
	}
	defer wiersze.Close()

	lista := []KrokAutomatyki{}
	for wiersze.Next() {
		var krok KrokAutomatyki
		var nazwa, komenda, parametry, warunek, odwolania sql.NullString
		err := wiersze.Scan(&krok.ID, &krok.Kod, &nazwa, &krok.Rodzaj, &komenda,
			&parametry, &warunek, &krok.Kolejnosc, &odwolania)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kroku automatyki: %w", err)
		}
		krok.Nazwa, krok.Komenda = tekstZKolumny(nazwa), tekstZKolumny(komenda)
		krok.Parametry, krok.Warunek = tekstZKolumny(parametry), tekstZKolumny(warunek)
		krok.OdwolaniaSekretow = tekstZKolumny(odwolania)
		lista = append(lista, krok)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kroków automatyki %d: %w", automatykaID, err)
	}
	return lista, nil
}

// ZapiszZaleznosci podmienia układ zależności automatyki. Cyklu ani łuku do
// kroku nieistniejącego repozytorium nie ocenia — to zadanie walidacji układu
// w rdzeniu, która ma oddać Operatorowi zastrzeżenia, a nie odmowę zapisu.
func (r *repozytoriumAutomatyk) ZapiszZaleznosci(ctx context.Context,
	automatykaID int64, zaleznosci []ZaleznoscKroku) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunZaleznosciAutomatyki)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, automatykaID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić zależności automatyki %d: %w", automatykaID, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawZaleznoscKroku)
		if err != nil {
			return err
		}
		for _, zaleznosc := range zaleznosci {
			_, err := wstawienie.ExecContext(ctx, automatykaID, zaleznosc.KrokZ,
				zaleznosc.KrokDo, zaleznosc.Rodzaj, tekstDoKolumny(zaleznosc.Warunek))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać zależności %q→%q automatyki %d: %w",
					zaleznosc.KrokZ, zaleznosc.KrokDo, automatykaID, err)
			}
		}
		return nil
	})
}

// Zaleznosci zwraca cały układ zależności między krokami automatyki wprost z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) Zaleznosci(ctx context.Context, automatykaID int64) ([]ZaleznoscKroku, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaZaleznosciAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zależności automatyki %d: %w", automatykaID, err)
	}
	defer wiersze.Close()

	lista := []ZaleznoscKroku{}
	for wiersze.Next() {
		var zaleznosc ZaleznoscKroku
		var warunek sql.NullString
		if err := wiersze.Scan(&zaleznosc.KrokZ, &zaleznosc.KrokDo, &zaleznosc.Rodzaj, &warunek); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zależności automatyki: %w", err)
		}
		zaleznosc.Warunek = tekstZKolumny(warunek)
		lista = append(lista, zaleznosc)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zależności automatyki %d: %w", automatykaID, err)
	}
	return lista, nil
}
