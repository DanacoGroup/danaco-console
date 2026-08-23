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

// nazwaSterownika — sterownik zarejestrowany przez modernc.org/sqlite.
const nazwaSterownika = "sqlite"

// maksPolaczen ogranicza pulę połączeń. SQLite w trybie WAL dopuszcza wielu
// czytelników, ale wyłącznie jednego pisarza naraz; nadmiar równoległych
// połączeń zamienia rywalizację o zapis w błędy „database is locked” zamiast
// czekać na busy_timeout. Skromny limit trzyma pulę w ryzach, a bezczynne
// połączenia utrzymuje ciepłe, żeby WAL nie był otwierany i zamykany bez końca.
const maksPolaczen = 4

// pragmyPolaczenia obowiązują każde połączenie z puli, dlatego trafiają do DSN,
// a nie do pojedynczego zapytania wykonanego po otwarciu bazy.
var pragmyPolaczenia = []string{
	"foreign_keys(1)",
	"busy_timeout(5000)",
	"journal_mode(WAL)",
	"synchronous(NORMAL)",
}

// Baza to otwarty plik bazy wraz z pulą połączeń.
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
		if err := os.MkdirAll(katalog, 0o755); err != nil {
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
	return baza, nil
}

// Zamknij zamyka pulę połączeń. Wywołanie na pustej bazie jest bezpieczne.
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
