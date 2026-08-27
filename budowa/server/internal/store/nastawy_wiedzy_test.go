package store

import (
	"testing"

	"danacoconsole/server/internal/wiedza"
)

// Nastawy wskaźnika znaczenia stoją w dwóch miejscach i muszą znaczyć to samo.
//
// Wartość domyślną Operator dostaje z bazy: rozstrzygacz zasięgu oddaje
// `definicja_ustawienia.wartosc_domyslna` jako rozstrzygnięcie o pochodzeniu
// „domyślna", a adapter wiedzy nanosi je na komplet nastaw. Stałe pakietu
// `wiedza` wchodzą tam, gdzie rozstrzygacza nie ma. Rozjazd między jednym
// a drugim nie wywraca niczego przy starcie — daje dwie różne odpowiedzi na
// pytanie „czym rdzeń liczy", zależne od drogi wywołania, a rozpoznać to można
// dopiero po pustym wyniku wyszukiwania. Sprawdzian wiąże więc obie prawdy
// w jedną.
func TestNastawyWiedzyWskazujaWagiStojace(t *testing.T) {
	baza := swiezaBaza(t)

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

	// Sonda dodatnia dla samego odczytu: klucz, którego w katalogu nie ma,
	// musi się nie odczytać. Bez niej sprawdzian przechodziłby także wtedy,
	// gdyby zapytanie milczało o każdym kluczu.
	var nic string
	if err := baza.DB.QueryRow(
		`SELECT wartosc_domyslna FROM definicja_ustawienia WHERE klucz = ?`,
		"wiedza_klucz_ktorego_nie_ma").Scan(&nic); err == nil {
		t.Fatal("odczyt oddał wartość dla klucza spoza katalogu — mierzy co innego, niż sądzi")
	}

	// Katalog wag pusty znaczy „pobierz wagi od nowa do katalogu danych rdzenia"
	// (`wiedza/pomocnik.go`, `katalogWag`). Świeże wdrożenie ma wystartować na
	// wagach rozłożonych na maszynie, więc wartość pusta jest tu regresją, a nie
	// wyborem.
	if wiedza.KatalogModeliDomyslny == "" {
		t.Error("domyślny katalog wag jest pusty — świeże wdrożenie pobierze model z sieci")
	}
}
