package store

import (
	"fmt"
	"strings"
)

// naruszenieKluczaObcego opisuje jeden wiersz wyniku zapytania sprawdzającego więzy kluczy obcych bazy.
type naruszenieKluczaObcego struct {
	Tabela        string
	Wiersz        int64
	TabelaRodzica string
	NumerKlucza   int64
}

func (n naruszenieKluczaObcego) String() string {
	return fmt.Sprintf("%s(rowid=%d) → %s (klucz %d)", n.Tabela, n.Wiersz, n.TabelaRodzica, n.NumerKlucza)
}

// SprawdzSpojnosc wykonuje pełną kontrolę: integralność pliku oraz zgodność
// kluczy obcych. Zwraca błąd opisujący pierwsze wykryte naruszenie.
func (b *Baza) SprawdzSpojnosc() error {
	if err := b.sprawdzIntegralnosc(); err != nil {
		return err
	}
	naruszenia, err := b.sprawdzKluczeObce()
	if err != nil {
		return err
	}
	if len(naruszenia) > 0 {
		opisy := make([]string, 0, len(naruszenia))
		for _, naruszenie := range naruszenia {
			opisy = append(opisy, naruszenie.String())
		}
		return fmt.Errorf("store: naruszone klucze obce: %s", strings.Join(opisy, "; "))
	}
	return nil
}

// Metoda sprawdzIntegralnosc wykonuje kontrolę integralności pliku bazy i oczekuje wyniku pozytywnego.
func (b *Baza) sprawdzIntegralnosc() error {
	wynik, err := b.PragmaTekstowa("integrity_check")
	if err != nil {
		return err
	}
	if wynik != "ok" {
		return fmt.Errorf("store: integrity_check zgłosił %q", wynik)
	}
	return nil
}

// sprawdzKluczeObce wykonuje PRAGMA foreign_key_check i zwraca listę naruszeń.
// Pusta lista oznacza bazę zgodną z więzami.
func (b *Baza) sprawdzKluczeObce() ([]naruszenieKluczaObcego, error) {
	wiersze, err := b.DB.Query("PRAGMA foreign_key_check")
	if err != nil {
		return nil, fmt.Errorf("store: foreign_key_check nie powiódł się: %w", err)
	}
	defer wiersze.Close()

	naruszenia := []naruszenieKluczaObcego{}
	for wiersze.Next() {
		var naruszenie naruszenieKluczaObcego
		if err := wiersze.Scan(&naruszenie.Tabela, &naruszenie.Wiersz,
			&naruszenie.TabelaRodzica, &naruszenie.NumerKlucza); err != nil {
			return nil, fmt.Errorf("store: nieczytelny wynik foreign_key_check: %w", err)
		}
		naruszenia = append(naruszenia, naruszenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("store: przerwany odczyt foreign_key_check: %w", err)
	}
	return naruszenia, nil
}
