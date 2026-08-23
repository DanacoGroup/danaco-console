package store

import (
	"path/filepath"
	"testing"
)

// Skutek migracji dobudowy modułu Studio: czy byty, na których stoją jego
// funkcje, naprawdę powstały.
//
// Migracja, która „przeszła", nie jest dowodem: krok wykonany bez błędu
// zostawia wersję w dzienniku migracji niezależnie od tego, czy polecenie
// czegokolwiek dokonało. Dlatego sprawdzian nie pyta o wersję schematu, tylko
// o tabele, kolumny i wiersz katalogu okien — czyli o to, po co ta migracja
// powstała.
func TestMigracjaDobudowyStudiaZakladaByty(t *testing.T) {
	baza, err := Otworz(filepath.Join(t.TempDir(), "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })

	for _, tabela := range []string{
		"komentarz_studio", "zmiana_sledzona_studio", "galaz_studio",
		"pozycja_wczytywania_studio", "operacja_studio", "lancuch_studio",
		"profil_wydania_studio", "szablon_studio", "odwolanie_wersji_studio",
	} {
		var nazwa string
		err := baza.DB.QueryRow(
			`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, tabela).Scan(&nazwa)
		if err != nil {
			t.Errorf("tabela %s nie powstała: %v", tabela, err)
		}
	}

	// Cztery kolumny wersji, bez których rozdział 3.6 opracowania jest
	// niewykonalny — a najważniejsza z nich jest kolumna autora.
	var kolumny int
	err = baza.DB.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('wersja_dokumentu_studio')
	                          WHERE name IN ('autor','kamien_milowy','galaz_id','propozycja_id')`).Scan(&kolumny)
	if err != nil || kolumny != 4 {
		t.Errorf("wersja dokumentu ma %d z 4 kolumn dobudowanych (błąd: %v)", kolumny, err)
	}

	// Wiersz katalogu okien: definicja ORAZ przypięcie do modułu. Sama
	// definicja bez przypięcia zostawiłaby okno poza zakresem modułu, czyli
	// dokładnie tam, gdzie było.
	var przypiete int
	err = baza.DB.QueryRow(`SELECT COUNT(*) FROM okno_operacyjne o
	                          JOIN okno_operacyjne_modul m ON m.okno_operacyjne_id = o.id
	                          JOIN modul d ON d.id = m.modul_id
	                          WHERE o.kod = 'ingest-ocr-panel' AND d.kod = 'studio'`).Scan(&przypiete)
	if err != nil || przypiete != 1 {
		t.Errorf("Ingest/OCR Panel nie jest przypięty do modułu Studio: %d (błąd: %v)", przypiete, err)
	}
}
