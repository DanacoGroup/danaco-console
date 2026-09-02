package store_test

import (
	"testing"

	"danacoconsole/server/internal/wiedza"
)

// TestNastawyWiedzyWskazujaWagiStojace sprawdza, że nastawy wskaźnika znaczenia i stałe pakietu wskazują te same wagi.
func TestNastawyWiedzyWskazujaWagiStojace(t *testing.T) {
	baza := swiezaBazaWiedzy(t)

	oczekiwane := map[string]string{
		wiedza.KluczModel:         wiedza.ModelDomyslny,
		wiedza.KluczKatalogModeli: wiedza.KatalogModeliDomyslny,
	}
	for klucz, stala := range oczekiwane {
		var domyslna string
		err := baza.DB.QueryRow(
			`SELECT wartosc_domyslna FROM definicja_ustawienia WHERE klucz = ?`,
			klucz).Scan(&domyslna)
		if err != nil {
			t.Fatalf("nastawa %s nie ma wiersza w katalogu ustawień: %v", klucz, err)
		}
		if domyslna != stala {
			t.Errorf("nastawa %s: baza mówi %q, stała pakietu wiedza mówi %q",
				klucz, domyslna, stala)
		}
	}

	// Sonda dodatnia dla samego odczytu: klucz, którego w katalogu nie ma, musi się nie odczytać.
	var nic string
	if err := baza.DB.QueryRow(
		`SELECT wartosc_domyslna FROM definicja_ustawienia WHERE klucz = ?`,
		"wiedza_klucz_ktorego_nie_ma").Scan(&nic); err == nil {
		t.Fatal("odczyt oddał wartość dla klucza spoza katalogu — mierzy co innego, niż sądzi")
	}

	// Katalog wag pusty znaczy pobranie wag od nowa do katalogu danych rdzenia przy starcie.
	if wiedza.KatalogModeliDomyslny == "" {
		t.Error("domyślny katalog wag jest pusty — świeże wdrożenie pobierze model z sieci")
	}
}
