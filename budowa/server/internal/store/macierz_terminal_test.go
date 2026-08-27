package store

import (
	"database/sql"
	"path/filepath"
	"testing"
)

// TestTerminalUkrytyWCodeStudioWbrewMacierzy jest zaporą pilnującą rozjazdu między macierzą dostępności modułów a ukryciem modułu Terminal w CodeStudio.
func TestTerminalUkrytyWCodeStudioWbrewMacierzy(t *testing.T) {
	baza, err := Otworz(filepath.Join(t.TempDir(), "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy sprawdzianu: %v", err)
	}
	defer func() { _ = baza.Zamknij() }()

	var widoczny int
	err = baza.DB.QueryRow(`
		SELECT sm.widoczny
		  FROM srodowisko_modul sm
		  JOIN srodowisko s ON s.id = sm.srodowisko_id
		  JOIN modul m      ON m.id = sm.modul_id
		 WHERE s.kod = 'codestudio' AND m.kod = 'terminal'`).Scan(&widoczny)
	if err == sql.ErrNoRows {
		t.Fatal("pary (codestudio, terminal) nie ma w macierzy — macierz przestała " +
			"być kompletem par; przeczytaj rozjazd od nowa")
	}
	if err != nil {
		t.Fatalf("nie można odczytać widoczności Terminala w CodeStudio: %v", err)
	}

	if widoczny == 1 {
		t.Fatal("USTERKA NAPRAWIONA — Terminal jest widoczny w CodeStudio zgodnie " +
			"z macierzą dostępności (rozdz. 10). Odwróć ten sprawdzian: zamień go " +
			"na straż pilnującą, że para (codestudio, terminal) ma widoczny = 1.")
	}

	t.Log("USTERKA CZYNNA: Terminal ukryty w CodeStudio (widoczny = 0) wbrew " +
		"macierzy dostępności, która daje mu TAK w tej kolumnie. Migracja 080 " +
		"gasi go rozstrzygnięciem wykonawcy przeciw dostawie.")
}
