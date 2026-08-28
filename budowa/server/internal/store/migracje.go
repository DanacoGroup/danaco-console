package store

import "fmt"

// schematRejestruMigracji — rejestr zastosowanych kroków. Zakładany przed
// pierwszą migracją, dlatego jako jedyny nie pochodzi z pliku .sql.
const schematRejestruMigracji = `
CREATE TABLE IF NOT EXISTS migracja (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    wersja         INTEGER NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    suma_kontrolna TEXT    NOT NULL,
    zastosowano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
)`

// Migruj doprowadza schemat do najnowszej wersji. Każdy krok wykonywany jest
// w jednej transakcji razem z wpisem do rejestru — schemat nigdy nie zostaje
// zastosowany połowicznie.
func (b *Baza) Migruj() error {
	kroki, err := wczytajMigracje()
	if err != nil {
		return err
	}
	if _, err := b.DB.Exec(schematRejestruMigracji); err != nil {
		return fmt.Errorf("store: nie można założyć rejestru migracji: %w", err)
	}
	zastosowane, err := b.zastosowaneMigracje()
	if err != nil {
		return err
	}
	for _, krok := range kroki {
		suma, jest := zastosowane[krok.Wersja]
		if jest {
			if suma != krok.SumaKontrolna {
				return fmt.Errorf("store: migracja %03d (%s) zmieniła treść po zastosowaniu",
					krok.Wersja, krok.Nazwa)
			}
			continue
		}
		if err := b.zastosujMigracje(krok); err != nil {
			return err
		}
	}
	return nil
}

// Metoda zastosowaneMigracje zwraca mapę wersja-suma kontrolna dla kroków migracji już wykonanych w tej bazie.
func (b *Baza) zastosowaneMigracje() (map[int]string, error) {
	wiersze, err := b.DB.Query("SELECT wersja, suma_kontrolna FROM migracja")
	if err != nil {
		return nil, fmt.Errorf("store: nie można odczytać rejestru migracji: %w", err)
	}
	defer wiersze.Close()

	zastosowane := map[int]string{}
	for wiersze.Next() {
		var wersja int
		var suma string
		if err := wiersze.Scan(&wersja, &suma); err != nil {
			return nil, fmt.Errorf("store: uszkodzony wpis rejestru migracji: %w", err)
		}
		zastosowane[wersja] = suma
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("store: przerwany odczyt rejestru migracji: %w", err)
	}
	return zastosowane, nil
}

// Metoda zastosujMigracje wykonuje treść pojedynczego kroku migracji i odnotowuje go w rejestrze kroków.
func (b *Baza) zastosujMigracje(krok migracja) error {
	transakcja, err := b.DB.Begin()
	if err != nil {
		return fmt.Errorf("store: nie można otworzyć transakcji migracji %03d: %w", krok.Wersja, err)
	}
	defer transakcja.Rollback()

	if _, err := transakcja.Exec(krok.Tresc); err != nil {
		return fmt.Errorf("store: migracja %03d (%s) nie powiodła się: %w", krok.Wersja, krok.Nazwa, err)
	}
	_, err = transakcja.Exec(
		"INSERT INTO migracja (wersja, nazwa, suma_kontrolna) VALUES (?, ?, ?)",
		krok.Wersja, krok.Nazwa, krok.SumaKontrolna)
	if err != nil {
		return fmt.Errorf("store: nie można odnotować migracji %03d: %w", krok.Wersja, err)
	}
	if err := transakcja.Commit(); err != nil {
		return fmt.Errorf("store: nie można zatwierdzić migracji %03d: %w", krok.Wersja, err)
	}
	return nil
}
