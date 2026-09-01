// Pakiet store odpowiada za trwałość Danaco Console: otwarcie pliku
// SQLite, doprowadzenie schematu do bieżącej wersji migracjami oraz kontrolę
// spójności bazy. Sterownik jest czystym Go (modernc.org/sqlite), bez CGO.
package store

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// nazwaSterownika wskazuje sterownik bazy danych zarejestrowany przez bibliotekę modernc.org/sqlite dla języka Go.
const nazwaSterownika = "sqlite"

// maksPolaczen ogranicza pulę połączeń, ponieważ SQLite w trybie WAL dopuszcza wielu czytelników, ale wyłącznie jednego pisarza naraz.
const maksPolaczen = 4

// trybTransakcji każe sterownikowi otwierać transakcje poleceniem BEGIN
// IMMEDIATE. Przy BEGIN DEFERRED zapis poprzedzony odczytem podnosi blokadę
// dopiero przy pierwszym zapisie, więc dwie transakcje czytają tę samą wartość
// i jedna z nich ginie — mierzone 14 przyrostów zamiast 20, z błędami 517 i 5
// mimo busy_timeout. Blokada wzięta od razu zamienia zgubiony zapis na czekanie.
const trybTransakcji = "immediate"

// prawaPlikuBazy i prawaKatalogu odcinają grupę i pozostałych od pliku, który
// niesie rozmowy, dokumenty, skróty sesji bramki i żetony udostępnień.
const (
	prawaPlikuBazy = 0o600
	prawaKatalogu  = 0o700
)

// pragmyPolaczenia obowiązują każde połączenie z puli, dlatego trafiają do DSN,
// a nie do pojedynczego zapytania wykonanego po otwarciu bazy.
var pragmyPolaczenia = []string{
	"foreign_keys(1)",
	"busy_timeout(5000)",
	"journal_mode(WAL)",
	"synchronous(NORMAL)",
}

// Baza to otwarty plik bazy danych SQLite wraz z pulą jego aktywnych połączeń gotowych do wykonywania zapytań.
type Baza struct {
	DB      *sql.DB
	Sciezka string
}

// Otworz zakłada brakujący katalog i plik bazy, ustawia pragmy oraz stosuje
// wszystkie niezastosowane migracje. Zwraca gotową do pracy bazę.
func Otworz(sciezka string) (*Baza, error) {
	if strings.TrimSpace(sciezka) == "" {
		return nil, fmt.Errorf("store: pusta ścieżka pliku bazy")
	}
	if katalog := filepath.Dir(sciezka); katalog != "" && katalog != "." {
		if err := os.MkdirAll(katalog, prawaKatalogu); err != nil {
			return nil, fmt.Errorf("store: nie można założyć katalogu %q: %w", katalog, err)
		}
	}
	db, err := sql.Open(nazwaSterownika, zbudujDSN(sciezka))
	if err != nil {
		return nil, fmt.Errorf("store: nie można otworzyć bazy %q: %w", sciezka, err)
	}
	db.SetMaxOpenConns(maksPolaczen)
	db.SetMaxIdleConns(maksPolaczen)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("store: baza %q nie odpowiada: %w", sciezka, err)
	}
	baza := &Baza{DB: db, Sciezka: sciezka}
	if err := baza.Migruj(); err != nil {
		db.Close()
		return nil, err
	}
	// Prawa nadaje się po migracjach, bo dziennik zapisu wyprzedzającego i plik
	// pamięci wspólnej powstają dopiero przy pierwszym zapisie do bazy.
	if err := zawezPrawaPlikow(sciezka); err != nil {
		db.Close()
		return nil, err
	}
	return baza, nil
}

// zawezPrawaPlikow odbiera grupie i pozostałym dostęp do pliku bazy oraz do jego
// dziennika zapisu wyprzedzającego i pliku pamięci wspólnej. Plik nieistniejący
// nie jest błędem: dziennik i pamięć wspólna znikają przy czystym zamknięciu bazy.
func zawezPrawaPlikow(sciezka string) error {
	for _, plik := range []string{sciezka, sciezka + "-wal", sciezka + "-shm"} {
		if err := os.Chmod(plik, prawaPlikuBazy); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("store: nie można zawęzić praw pliku %q: %w", plik, err)
		}
	}
	return nil
}

// Metoda Zamknij zamyka pulę połączeń bazy; wywołanie na pustej, niezainicjowanej bazie jest bezpieczne.
func (b *Baza) Zamknij() error {
	if b == nil || b.DB == nil {
		return nil
	}
	return b.DB.Close()
}

// zbudujDSN składa ścieżkę pliku z listą pragm. Bez przedrostka „file:”, aby
// sterownik odciął część zapytania i przekazał do SQLite czystą ścieżkę.
func zbudujDSN(sciezka string) string {
	parametry := url.Values{}
	for _, pragma := range pragmyPolaczenia {
		parametry.Add("_pragma", pragma)
	}
	parametry.Set("_txlock", trybTransakcji)
	return filepath.ToSlash(sciezka) + "?" + parametry.Encode()
}

// PragmaTekstowa zwraca wartość pragmy zwracającej pojedynczy wiersz —
// używana przez kontrolę spójności i diagnostykę.
func (b *Baza) PragmaTekstowa(nazwa string) (string, error) {
	var wartosc string
	if err := b.DB.QueryRow("PRAGMA " + nazwa).Scan(&wartosc); err != nil {
		return "", fmt.Errorf("store: PRAGMA %s nie powiodła się: %w", nazwa, err)
	}
	return wartosc, nil
}
