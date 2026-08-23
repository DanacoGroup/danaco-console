package store

import (
	"database/sql"
	"path/filepath"
	"testing"
)

// ZAPORA, NIE ZGODA — Terminal w CodeStudio.
//
// Macierz dostępności modułów w opracowaniu Właściciela
// (`docs/architektura/koncepcja-platformy.md` rozdz. 10, wiersz „Terminal
// (11.11) | NIE | NIE | TAK") stawia Terminal jako moduł widoczny w CodeStudio,
// z własnymi oknami (rozdz. 11.11: Terminal Tabs, Output Console, Process
// Monitor). Migracja 080 gasi go z `widoczny = 0` z uzasadnieniem „Terminal nie
// jest samodzielnym modułem" — rozstrzygnięcie wykonawcy PRZECIW dostawie.
//
// Sprawdzian zapisuje ten rozjazd tak, żeby nie zniknął po cichu ani nie
// pogłębił się po cichu: wypada niepomyślnie dopóki Terminal jest ukryty
// w CodeStudio, a gdy zostanie odsłonięty zgodnie z macierzą, każe się odwrócić
// w straż. Naprawa należy do warstwy modułów (macierz `srodowisko_modul`),
// nie tutaj — sprawdzian pomiar utrwala, nie rozstrzyga.
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
