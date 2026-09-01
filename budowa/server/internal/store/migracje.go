package store

import (
	"context"
	"database/sql"
	"fmt"
)

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
	// PRAGMA user_version równa zeru oznacza bazę sprzed uzgodnienia sum znormalizowanych.
	var znacznikUzgodnienia int
	if err := b.DB.QueryRow("PRAGMA user_version").Scan(&znacznikUzgodnienia); err != nil {
		return fmt.Errorf("store: nie można odczytać znacznika uzgodnienia sum: %w", err)
	}
	if znacznikUzgodnienia == 0 {
		if err := b.uzgodnijSumyKontrolne(kroki, zastosowane); err != nil {
			return err
		}
	}
	najwyzszaZastosowana := 0
	for wersja := range zastosowane {
		if wersja > najwyzszaZastosowana {
			najwyzszaZastosowana = wersja
		}
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
		if krok.Wersja < najwyzszaZastosowana {
			return fmt.Errorf("store: migracja %03d (%s) ma numer niższy od zastosowanego %03d — "+
				"krok dołożony w lukę numeracji rozjeżdża schemat instalacji świeżej i zaktualizowanej; "+
				"nadaj krokowi numer wyższy od %03d",
				krok.Wersja, krok.Nazwa, najwyzszaZastosowana, najwyzszaZastosowana)
		}
		if err := b.zastosujMigracje(krok); err != nil {
			return err
		}
	}
	return nil
}

// Metoda uzgodnijSumyKontrolne jednorazowo przepisuje sumy kroków już zastosowanych na wartości liczone z treści znormalizowanej i stawia PRAGMA user_version na 1, dzięki czemu bazy migrowane przed normalizacją wstają bez ręcznej naprawy.
func (b *Baza) uzgodnijSumyKontrolne(kroki []migracja, zastosowane map[int]string) error {
	transakcja, err := b.DB.Begin()
	if err != nil {
		return fmt.Errorf("store: nie można otworzyć transakcji uzgodnienia sum: %w", err)
	}
	defer transakcja.Rollback()

	for _, krok := range kroki {
		stara, jest := zastosowane[krok.Wersja]
		if !jest || stara == krok.SumaKontrolna {
			continue
		}
		if _, err := transakcja.Exec("UPDATE migracja SET suma_kontrolna = ? WHERE wersja = ?",
			krok.SumaKontrolna, krok.Wersja); err != nil {
			return fmt.Errorf("store: nie można uzgodnić sumy migracji %03d: %w", krok.Wersja, err)
		}
		zastosowane[krok.Wersja] = krok.SumaKontrolna
	}
	if _, err := transakcja.Exec("PRAGMA user_version = 1"); err != nil {
		return fmt.Errorf("store: nie można postawić znacznika uzgodnienia sum: %w", err)
	}
	if err := transakcja.Commit(); err != nil {
		return fmt.Errorf("store: nie można zatwierdzić uzgodnienia sum: %w", err)
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
//
// Krok idzie na wydzielonym połączeniu z wygaszonymi więzami kluczy obcych.
// Powód jest jeden i nie ma obejścia w treści kroku: przy włączonych więzach
// `DROP TABLE` wykonuje niejawne DELETE wszystkich wierszy, a to odpala kaskady
// ON DELETE i czyści tabele potomne bez jednego komunikatu — SQLite nie zna
// innego sposobu na zmianę więzu CHECK niż przebudowa tabeli, więc każdy taki
// krok kasowałby dane dzieci. Kroki 226, 269 i 378 stawiają wprawdzie w treści
// `PRAGMA foreign_keys = off`, ale wewnątrz transakcji ta pragma nie ma skutku,
// więc wygaszenie musi paść tutaj, przed otwarciem transakcji.
//
// Wygaszenie nie zdejmuje kontroli, tylko przesuwa ją na koniec kroku:
// `PRAGMA foreign_key_check` wykonany przed zatwierdzeniem wykazuje wiersz
// wskazujący na nieistniejący rodzica i kończy krok odmową. Różnica wobec więzów
// czynnych jest taka, że odmowa nazywa tabelę, a nie polecenie.
func (b *Baza) zastosujMigracje(krok migracja) error {
	zycie := context.Background()
	polaczenie, err := b.DB.Conn(zycie)
	if err != nil {
		return fmt.Errorf("store: nie można zająć połączenia dla migracji %03d: %w", krok.Wersja, err)
	}
	defer polaczenie.Close()

	if _, err := polaczenie.ExecContext(zycie, "PRAGMA foreign_keys = off"); err != nil {
		return fmt.Errorf("store: nie można wygasić więzów na czas migracji %03d: %w", krok.Wersja, err)
	}
	// Połączenie wraca do puli, więc więzy muszą wrócić także wtedy, gdy krok padł.
	defer polaczenie.ExecContext(zycie, "PRAGMA foreign_keys = on")

	transakcja, err := polaczenie.BeginTx(zycie, nil)
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
	if err := sprawdzWiezyPoKroku(transakcja, krok); err != nil {
		return err
	}
	if err := transakcja.Commit(); err != nil {
		return fmt.Errorf("store: nie można zatwierdzić migracji %03d: %w", krok.Wersja, err)
	}
	return nil
}

// sprawdzWiezyPoKroku wykazuje, że krok wykonany bez więzów nie zostawił wiersza
// wskazującego na rodzica, którego nie ma. Wynik niepusty kończy krok odmową
// jeszcze przed zatwierdzeniem, bo po zatwierdzeniu wycofanie nie jest możliwe.
func sprawdzWiezyPoKroku(transakcja *sql.Tx, krok migracja) error {
	wiersze, err := transakcja.Query("PRAGMA foreign_key_check")
	if err != nil {
		return fmt.Errorf("store: nie można sprawdzić więzów po migracji %03d: %w", krok.Wersja, err)
	}
	defer wiersze.Close()

	naruszenia := 0
	pierwszaTabela := ""
	for wiersze.Next() {
		var tabela, rodzic sql.NullString
		var wiersz, numerWiezu sql.NullInt64
		if err := wiersze.Scan(&tabela, &wiersz, &rodzic, &numerWiezu); err != nil {
			return fmt.Errorf("store: nieczytalny wynik sprawdzenia więzów po migracji %03d: %w", krok.Wersja, err)
		}
		if naruszenia == 0 {
			pierwszaTabela = tabela.String
		}
		naruszenia++
	}
	if err := wiersze.Err(); err != nil {
		return fmt.Errorf("store: przerwane sprawdzenie więzów po migracji %03d: %w", krok.Wersja, err)
	}
	if naruszenia > 0 {
		return fmt.Errorf("store: migracja %03d (%s) zostawiła %d wierszy bez wskazywanego rodzica (pierwszy w tabeli %q)",
			krok.Wersja, krok.Nazwa, naruszenia, pierwszaTabela)
	}
	return nil
}
