package store

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"slices"
	"strings"
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

// schematZnacznikaUzgodnienia — znacznik jednorazowego uzgodnienia sum. Stoi
// w tabeli, bo PRAGMA user_version ginie przy zrzucie i wczytaniu bazy
// narzędziem sqlite3, a jedno polecenie na pliku bazy cofa ją do zera.
const schematZnacznikaUzgodnienia = `
CREATE TABLE IF NOT EXISTS uzgodnienie_sum (
    id          INTEGER PRIMARY KEY CHECK (id = 1),
    zastosowano TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
)`

// kluczNaruszenia grupuje wynik foreign_key_check po więzie, bez rowid: przebudowa
// tabeli w kroku migracji nadaje wierszom zastanym nowe rowid.
type kluczNaruszenia struct {
	Tabela        string
	TabelaRodzica string
	NumerKlucza   int64
}

type zrodloWierszy interface {
	Query(zapytanie string, argumenty ...any) (*sql.Rows, error)
}

// Migruj doprowadza schemat do najnowszej wersji. Krok idzie w jednej transakcji
// razem z wpisem do rejestru, więc nie zostaje zastosowany połowicznie.
func (b *Baza) Migruj() error {
	kroki, err := wczytajMigracje()
	if err != nil {
		return err
	}
	if _, err := b.DB.Exec(schematRejestruMigracji); err != nil {
		return fmt.Errorf("store: nie można założyć rejestru migracji: %w", err)
	}
	if _, err := b.DB.Exec(schematZnacznikaUzgodnienia); err != nil {
		return fmt.Errorf("store: nie można założyć znacznika uzgodnienia sum: %w", err)
	}
	zastosowane, err := b.zastosowaneMigracje()
	if err != nil {
		return err
	}
	if err := sprawdzZgodnoscNazw(kroki, zastosowane); err != nil {
		return err
	}
	uzgodnione, err := b.znacznikUzgodnieniaStoi()
	if err != nil {
		return err
	}
	if !uzgodnione {
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
	var zastaneNaruszenia map[kluczNaruszenia]int
	for _, krok := range kroki {
		wpis, jest := zastosowane[krok.Wersja]
		if jest {
			if wpis.SumaKontrolna != krok.SumaKontrolna {
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
		if zastaneNaruszenia == nil {
			if zastaneNaruszenia, err = b.zastaneNaruszeniaWiezow(); err != nil {
				return err
			}
		}
		if zastaneNaruszenia, err = b.zastosujMigracje(krok, zastaneNaruszenia); err != nil {
			return err
		}
	}
	return nil
}

// sprawdzZgodnoscNazw porównuje nazwę kroku zastosowanego z nazwą pliku o tym
// numerze i idzie przy każdym starcie: numer sam tożsamości kroku nie dowodzi,
// a suma w rejestrze może pochodzić z uzgodnienia. Obie nazwy biorą się z członu
// nazwy pliku migracja_NNN_nazwa.sql, więc porównanie idzie bez normalizacji.
func sprawdzZgodnoscNazw(kroki []migracja, zastosowane map[int]migracja) error {
	znane := make(map[int]struct{}, len(kroki))
	for _, krok := range kroki {
		znane[krok.Wersja] = struct{}{}
	}
	for wersja, wpis := range zastosowane {
		if _, jest := znane[wersja]; jest {
			continue
		}
		return fmt.Errorf("store: rejestr niesie migrację %03d (%s), której repozytorium nie zna — "+
			"baza idzie inną linią numeracji albo pochodzi z wersji nowszej niż ten rdzeń",
			wersja, wpis.Nazwa)
	}
	for _, krok := range kroki {
		wpis, jest := zastosowane[krok.Wersja]
		if !jest || wpis.Nazwa == krok.Nazwa {
			continue
		}
		return fmt.Errorf("store: migracja %03d stoi w rejestrze pod nazwą %q, a repozytorium ma pod tym numerem %q — "+
			"baza idzie inną linią numeracji i kroku %03d (%s) na niej nie wykonano",
			krok.Wersja, wpis.Nazwa, krok.Nazwa, krok.Wersja, krok.Nazwa)
	}
	return nil
}

// znacznikUzgodnieniaStoi mówi, czy sumy tej bazy zostały już uzgodnione. Bazy
// uzgodnione przed przeniesieniem znacznika do tabeli niosą go w PRAGMA
// user_version; powtórne uzgodnienie przykryłoby zmianę treści kroku.
func (b *Baza) znacznikUzgodnieniaStoi() (bool, error) {
	var wierszy int
	if err := b.DB.QueryRow("SELECT count(*) FROM uzgodnienie_sum").Scan(&wierszy); err != nil {
		return false, fmt.Errorf("store: nie można odczytać znacznika uzgodnienia sum: %w", err)
	}
	if wierszy > 0 {
		return true, nil
	}
	var znacznikWNaglowku int
	if err := b.DB.QueryRow("PRAGMA user_version").Scan(&znacznikWNaglowku); err != nil {
		return false, fmt.Errorf("store: nie można odczytać znacznika uzgodnienia sum z nagłówka pliku: %w", err)
	}
	if znacznikWNaglowku == 0 {
		return false, nil
	}
	if _, err := b.DB.Exec("INSERT INTO uzgodnienie_sum (id) VALUES (1)"); err != nil {
		return false, fmt.Errorf("store: nie można przenieść znacznika uzgodnienia sum do tabeli: %w", err)
	}
	return true, nil
}

// uzgodnijSumyKontrolne przepisuje sumy kroków już zastosowanych na wartości
// liczone z treści znormalizowanej. Idzie raz na bazę i wyłącznie po kontroli
// nazw, która orzeka, że rejestr opisuje te same kroki co repozytorium.
func (b *Baza) uzgodnijSumyKontrolne(kroki []migracja, zastosowane map[int]migracja) error {
	transakcja, err := b.DB.Begin()
	if err != nil {
		return fmt.Errorf("store: nie można otworzyć transakcji uzgodnienia sum: %w", err)
	}
	defer transakcja.Rollback()

	for _, krok := range kroki {
		wpis, jest := zastosowane[krok.Wersja]
		if !jest || wpis.SumaKontrolna == krok.SumaKontrolna {
			continue
		}
		if _, err := transakcja.Exec("UPDATE migracja SET suma_kontrolna = ? WHERE wersja = ?",
			krok.SumaKontrolna, krok.Wersja); err != nil {
			return fmt.Errorf("store: nie można uzgodnić sumy migracji %03d: %w", krok.Wersja, err)
		}
		wpis.SumaKontrolna = krok.SumaKontrolna
		zastosowane[krok.Wersja] = wpis
	}
	if _, err := transakcja.Exec("INSERT INTO uzgodnienie_sum (id) VALUES (1)"); err != nil {
		return fmt.Errorf("store: nie można postawić znacznika uzgodnienia sum: %w", err)
	}
	if err := transakcja.Commit(); err != nil {
		return fmt.Errorf("store: nie można zatwierdzić uzgodnienia sum: %w", err)
	}
	return nil
}

// zastosowaneMigracje zwraca wpisy rejestru po numerze wersji. Pole Tresc zostaje
// puste: rejestr niesie numer, nazwę i sumę, treści kroku nie przechowuje.
func (b *Baza) zastosowaneMigracje() (map[int]migracja, error) {
	wiersze, err := b.DB.Query("SELECT wersja, nazwa, suma_kontrolna FROM migracja")
	if err != nil {
		return nil, fmt.Errorf("store: nie można odczytać rejestru migracji: %w", err)
	}
	defer wiersze.Close()

	zastosowane := map[int]migracja{}
	for wiersze.Next() {
		var wpis migracja
		if err := wiersze.Scan(&wpis.Wersja, &wpis.Nazwa, &wpis.SumaKontrolna); err != nil {
			return nil, fmt.Errorf("store: uszkodzony wpis rejestru migracji: %w", err)
		}
		zastosowane[wpis.Wersja] = wpis
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("store: przerwany odczyt rejestru migracji: %w", err)
	}
	return zastosowane, nil
}

// zastosujMigracje wykonuje krok na wydzielonym połączeniu z wygaszonymi więzami:
// przy więzach czynnych DROP TABLE i przebudowa tabeli (jedyna droga zmiany więzu
// CHECK w SQLite) wykonują niejawne DELETE i odpalają kaskady ON DELETE. Wewnątrz
// transakcji ta pragma nie ma skutku, więc pada przed jej otwarciem.
func (b *Baza) zastosujMigracje(krok migracja, zastaneNaruszenia map[kluczNaruszenia]int) (
	poKroku map[kluczNaruszenia]int, blad error) {

	zycie := context.Background()
	polaczenie, err := b.DB.Conn(zycie)
	if err != nil {
		return nil, fmt.Errorf("store: nie można zająć połączenia dla migracji %03d: %w", krok.Wersja, err)
	}
	defer polaczenie.Close()

	if _, err := polaczenie.ExecContext(zycie, "PRAGMA foreign_keys = off"); err != nil {
		return nil, fmt.Errorf("store: nie można wygasić więzów na czas migracji %03d: %w", krok.Wersja, err)
	}
	defer func() {
		if _, err := polaczenie.ExecContext(zycie, "PRAGMA foreign_keys = on"); err != nil {
			odrzucPolaczenie(polaczenie)
			if blad == nil {
				blad = fmt.Errorf("store: nie można przywrócić więzów po migracji %03d: %w", krok.Wersja, err)
			}
		}
	}()
	// Krok kasujący wiersz zaczynu z adresem maszyny zostawia jego treść na
	// zwolnionych stronach pliku bazy, czytelną zwykłym grepem; secure_delete
	// zeruje ją przy kasowaniu. Wartość zastana wraca razem z połączeniem: pragma
	// należy do połączenia, a to wraca do puli i obsługuje pracę Operatora.
	var zerowanieZastane int
	if err := polaczenie.QueryRowContext(zycie, "PRAGMA secure_delete").Scan(&zerowanieZastane); err != nil {
		return nil, fmt.Errorf("store: nie można odczytać zerowania kasowanych stron przed migracją %03d: %w",
			krok.Wersja, err)
	}
	if _, err := polaczenie.ExecContext(zycie, "PRAGMA secure_delete = on"); err != nil {
		return nil, fmt.Errorf("store: nie można włączyć zerowania kasowanych stron na czas migracji %03d: %w",
			krok.Wersja, err)
	}
	defer func() {
		polecenie := fmt.Sprintf("PRAGMA secure_delete = %d", zerowanieZastane)
		if _, err := polaczenie.ExecContext(zycie, polecenie); err != nil {
			odrzucPolaczenie(polaczenie)
			if blad == nil {
				blad = fmt.Errorf("store: nie można przywrócić zerowania kasowanych stron po migracji %03d: %w",
					krok.Wersja, err)
			}
		}
	}()

	// Przemianowanie tabeli zastępczej na nazwę pierwotną każe SQLite sparsować
	// cały schemat, niespójny póki tabela pierwotna nie stoi; pragma to znosi.
	var przemianowanieZastane int
	if err := polaczenie.QueryRowContext(zycie, "PRAGMA legacy_alter_table").Scan(&przemianowanieZastane); err != nil {
		return nil, fmt.Errorf("store: nie można odczytać trybu przemianowania tabel przed migracją %03d: %w",
			krok.Wersja, err)
	}
	if _, err := polaczenie.ExecContext(zycie, "PRAGMA legacy_alter_table = on"); err != nil {
		return nil, fmt.Errorf("store: nie można włączyć trybu przemianowania tabel na czas migracji %03d: %w",
			krok.Wersja, err)
	}
	defer func() {
		polecenie := fmt.Sprintf("PRAGMA legacy_alter_table = %d", przemianowanieZastane)
		if _, err := polaczenie.ExecContext(zycie, polecenie); err != nil {
			odrzucPolaczenie(polaczenie)
			if blad == nil {
				blad = fmt.Errorf("store: nie można przywrócić trybu przemianowania tabel po migracji %03d: %w",
					krok.Wersja, err)
			}
		}
	}()

	transakcja, err := polaczenie.BeginTx(zycie, nil)
	if err != nil {
		return nil, fmt.Errorf("store: nie można otworzyć transakcji migracji %03d: %w", krok.Wersja, err)
	}
	defer transakcja.Rollback()

	if _, err := transakcja.Exec(krok.Tresc); err != nil {
		return nil, fmt.Errorf("store: migracja %03d (%s) nie powiodła się: %w", krok.Wersja, krok.Nazwa, err)
	}
	_, err = transakcja.Exec(
		"INSERT INTO migracja (wersja, nazwa, suma_kontrolna) VALUES (?, ?, ?)",
		krok.Wersja, krok.Nazwa, krok.SumaKontrolna)
	if err != nil {
		return nil, fmt.Errorf("store: nie można odnotować migracji %03d: %w", krok.Wersja, err)
	}
	poKroku, err = sprawdzWiezyPoKroku(transakcja, krok, zastaneNaruszenia)
	if err != nil {
		return nil, err
	}
	if err := transakcja.Commit(); err != nil {
		return nil, fmt.Errorf("store: nie można zatwierdzić migracji %03d: %w", krok.Wersja, err)
	}
	return poKroku, nil
}

// odrzucPolaczenie wyprowadza połączenie z puli. Pragma foreign_keys żyje tyle,
// co połączenie, więc połączenie z nieprzywróconymi więzami nie może wrócić do
// puli i obsługiwać zapisów aplikacji.
func odrzucPolaczenie(polaczenie *sql.Conn) {
	_ = polaczenie.Raw(func(any) error { return driver.ErrBadConn })
}

// zastaneNaruszeniaWiezow zdejmuje stan więzów sprzed pierwszego kroku przejazdu.
// Wiersz bez rodzica wstawiony wcześniej — ręczną naprawą albo przez sqlite3,
// który więzy ma domyślnie wyłączone — nie obciąża kroku migracji.
func (b *Baza) zastaneNaruszeniaWiezow() (map[kluczNaruszenia]int, error) {
	zastane, err := policzNaruszeniaWiezow(b.DB)
	if err != nil {
		return nil, fmt.Errorf("store: nie można zdjąć stanu więzów przed migracjami: %w", err)
	}
	if wierszy, tabele := opiszNaruszenia(zastane); wierszy > 0 {
		log.Printf("store: przed migracjami baza ma %d wierszy bez wskazywanego rodzica w tabelach %s; "+
			"odmową kończą się wyłącznie wiersze przybyłe po kroku", wierszy, tabele)
	}
	return zastane, nil
}

// sprawdzWiezyPoKroku wykazuje, że krok wykonany bez więzów nie zostawił wiersza
// wskazującego na rodzica, którego nie ma. Odmowa pada przed zatwierdzeniem, bo
// po nim wycofanie nie jest możliwe. Stan po kroku wraca do wołającego: krok,
// który zastaną sierotę skasował, obniża odniesienie krokowi następnemu.
func sprawdzWiezyPoKroku(transakcja *sql.Tx, krok migracja,
	zastane map[kluczNaruszenia]int) (map[kluczNaruszenia]int, error) {

	po, err := policzNaruszeniaWiezow(transakcja)
	if err != nil {
		return nil, fmt.Errorf("store: nie można sprawdzić więzów po migracji %03d: %w", krok.Wersja, err)
	}
	przybyle := map[kluczNaruszenia]int{}
	for klucz, ile := range po {
		if roznica := ile - zastane[klucz]; roznica > 0 {
			przybyle[klucz] = roznica
		}
	}
	if wierszy, tabele := opiszNaruszenia(przybyle); wierszy > 0 {
		return nil, fmt.Errorf("store: migracja %03d (%s) zostawiła %d wierszy bez wskazywanego rodzica w tabelach %s",
			krok.Wersja, krok.Nazwa, wierszy, tabele)
	}
	return po, nil
}

func policzNaruszeniaWiezow(zrodlo zrodloWierszy) (map[kluczNaruszenia]int, error) {
	wiersze, err := zrodlo.Query("PRAGMA foreign_key_check")
	if err != nil {
		return nil, fmt.Errorf("foreign_key_check nie powiódł się: %w", err)
	}
	defer wiersze.Close()

	policzone := map[kluczNaruszenia]int{}
	for wiersze.Next() {
		var tabela, rodzic sql.NullString
		var wiersz, numerWiezu sql.NullInt64
		if err := wiersze.Scan(&tabela, &wiersz, &rodzic, &numerWiezu); err != nil {
			return nil, fmt.Errorf("nieczytelny wynik foreign_key_check: %w", err)
		}
		policzone[kluczNaruszenia{Tabela: tabela.String, TabelaRodzica: rodzic.String, NumerKlucza: numerWiezu.Int64}]++
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("przerwany odczyt foreign_key_check: %w", err)
	}
	return policzone, nil
}

func opiszNaruszenia(policzone map[kluczNaruszenia]int) (int, string) {
	wierszy := 0
	tabele := make([]string, 0, len(policzone))
	for klucz, ile := range policzone {
		wierszy += ile
		tabele = append(tabele, klucz.Tabela)
	}
	slices.Sort(tabele)
	return wierszy, strings.Join(slices.Compact(tabele), ", ")
}
