package store

import (
	"testing"

	"danacoconsole/server/internal/mowa"
	"danacoconsole/server/internal/wiedza"
)

// Cztery nastawy przesiewu i osi obrazu muszą mieć wiersz w katalogu ustawień.
//
// Bez wiersza zdolność dalej działa — rozstrzyganie nastawy czyta zapis
// niezależnie od katalogu definicji — ale `config.set` odmawia klucza spoza
// katalogu (`core/walidacja_klucza_ustawienia.go`), a okno konfiguracji
// wystawia wyłącznie pozycje katalogu. Brak wiersza znaczy więc dokładnie tyle:
// nastawa jest ustawialna wyłącznie ręcznym zapisem do bazy. Sprawdzian pilnuje
// jednego i drugiego naraz: że wiersz jest i że jego wartość domyślna mówi to
// samo, co stała pakietu, z którego liczy silnik.
func TestNastawyPrzesiewuIObrazuStojaWKatalogu(t *testing.T) {
	baza := swiezaBaza(t)

	oczekiwane := map[string]string{
		wiedza.KluczModelPrzesiewu:   wiedza.ModelPrzesiewuDomyslny,
		wiedza.KluczKatalogPrzesiewu: "",
		wiedza.KluczModelObrazu:      wiedza.ModelObrazuDomyslny,
		wiedza.KluczKatalogObrazu:    "",
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

	// Wiersz katalogu bez zasięgu i bez osi jest pozycją, której nie da się
	// zapisać: rozstrzygacz pyta o wartość dla poziomu i osi, a definicja
	// niedopuszczająca żadnego poziomu nie ma gdzie stanąć.
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

	// Sonda dodatnia dla samego odczytu: klucz, którego w katalogu nie ma, musi
	// się nie odczytać. Bez niej sprawdzian przechodziłby także wtedy, gdyby
	// zapytanie milczało o każdym kluczu.
	var nic string
	if err := baza.DB.QueryRow(
		`SELECT wartosc_domyslna FROM definicja_ustawienia WHERE klucz = ?`,
		"wiedza_klucz_ktorego_nie_ma").Scan(&nic); err == nil {
		t.Fatal("odczyt oddał wartość dla klucza spoza katalogu — mierzy co innego, niż sądzi")
	}
}

// Nastawa katalogu wag mowy stoi w dwóch miejscach i musi znaczyć to samo.
//
// Wartość domyślną Operator dostaje z bazy (rozstrzygacz zasięgu oddaje
// `definicja_ustawienia.wartosc_domyslna`), a stała pakietu `mowa` wchodzi tam,
// gdzie rozstrzygacza nie ma. Rozjazd daje dwie odpowiedzi na pytanie, gdzie
// leżą wagi, zależne od drogi wywołania — a rozpoznać go można dopiero po tym,
// że transkrypcja pobiera drugą kopię wag zamiast wystartować.
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

	// Katalog pusty znaczy „pamięć podręczna biblioteki w katalogu domowym
	// konta, które uruchomiło rdzeń" — wtedy widoczność wag zależy od tego, na
	// czyim koncie stoi proces. Wartość pusta jest tu regresją, a nie wyborem.
	if mowa.KatalogModeliDomyslny == "" {
		t.Error("domyślny katalog wag mowy jest pusty — wagi widzi tylko konto, " +
			"w którego katalogu domowym stoi pamięć podręczna")
	}
}
