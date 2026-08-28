package store

import (
	"testing"

	"danacoconsole/server/internal/mowa"
	"danacoconsole/server/internal/wiedza"
)

// TestNastawyPrzesiewuIObrazuStojaWKatalogu sprawdza, że cztery nastawy przesiewu i osi obrazu mają wiersz w katalogu ustawień.
func TestNastawyPrzesiewuIObrazuStojaWKatalogu(t *testing.T) {
	baza := swiezaBaza(t)

	oczekiwane := map[string]string{
		wiedza.KluczModelPrzesiewu:   wiedza.ModelPrzesiewuDomyslny,
		wiedza.KluczKatalogPrzesiewu: wiedza.KatalogPrzesiewuDomyslny,
		wiedza.KluczModelObrazu:      wiedza.ModelObrazuDomyslny,
		wiedza.KluczKatalogObrazu:    wiedza.KatalogObrazuDomyslny,
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

	// Wiersz katalogu bez zasięgu i bez osi jest pozycją, której nie da się zapisać.
	for klucz := range oczekiwane {
		var poziomy, osie int
		if err := baza.DB.QueryRow(`
			SELECT (SELECT COUNT(*) FROM definicja_ustawienia_zasieg z
			         JOIN definicja_ustawienia d ON d.id = z.definicja_id
			        WHERE d.klucz = ?),
			       (SELECT COUNT(*) FROM definicja_ustawienia_os w
			         JOIN definicja_ustawienia d ON d.id = w.definicja_id
			        WHERE d.klucz = ?)`, klucz, klucz).Scan(&poziomy, &osie); err != nil {
			t.Fatalf("nie można policzyć zasięgów i osi nastawy %s: %v", klucz, err)
		}
		if poziomy == 0 || osie == 0 {
			t.Errorf("nastawa %s: zasięgów %d, osi %d — pozycja bez zasięgu albo bez osi "+
				"jest w oknie konfiguracji widoczna i niezapisywalna", klucz, poziomy, osie)
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
	if wiedza.KatalogPrzesiewuDomyslny == "" {
		t.Error("domyślny katalog wag przesiewu jest pusty — świeże wdrożenie pobierze model z sieci")
	}
	if wiedza.KatalogObrazuDomyslny == "" {
		t.Error("domyślny katalog wag osi obrazu jest pusty — świeże wdrożenie pobierze model z sieci")
	}
}

// TestNastawaWagMowyWskazujeWagiStojace sprawdza, że nastawa katalogu wag mowy i stała pakietu wskazują te same wagi.
func TestNastawaWagMowyWskazujeWagiStojace(t *testing.T) {
	baza := swiezaBaza(t)

	var domyslna string
	if err := baza.DB.QueryRow(
		`SELECT wartosc_domyslna FROM definicja_ustawienia WHERE klucz = ?`,
		mowa.KluczKatalogModeli).Scan(&domyslna); err != nil {
		t.Fatalf("nastawa %s nie ma wiersza w katalogu ustawień: %v", mowa.KluczKatalogModeli, err)
	}
	if domyslna != mowa.KatalogModeliDomyslny {
		t.Errorf("nastawa %s: baza mówi %q, stała pakietu mowa mówi %q",
			mowa.KluczKatalogModeli, domyslna, mowa.KatalogModeliDomyslny)
	}

	// Katalog pusty znaczy pamięć podręczną biblioteki w katalogu domowym konta uruchamiającego rdzeń.
	if mowa.KatalogModeliDomyslny == "" {
		t.Error("domyślny katalog wag mowy jest pusty — wagi widzi tylko konto, " +
			"w którego katalogu domowym stoi pamięć podręczna")
	}
}
