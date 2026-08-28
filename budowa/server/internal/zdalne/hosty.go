// Pakiet zdalne prowadzi wybór operatora maszyny od nazwy hosta wykonania do decyzji toru zdalnego SSH.
package zdalne

import (
	"database/sql"
	"fmt"
	"strings"
)

// Host to wiersz tabeli hostów zdalnych, potrzebny torowi do zbudowania połączenia SSH z daną zdalną maszyną.
type Host struct {
	Id         int64
	Nazwa      string
	Adres      string
	Uzytkownik string
	Port       int
	Zgoda      bool
}

// AdresPolaczenia zwraca adres, z którym łączy się SSH: adres sieciowy wiersza,
// a gdy go nie wpisano — nazwę hosta (puste nie zatrzymuje pracy).
func (h Host) AdresPolaczenia() string {
	if adres := strings.TrimSpace(h.Adres); adres != "" {
		return adres
	}
	return h.Nazwa
}

// Funkcja hostOkna czyta nazwę hosta wykonania obowiązującą dla danego okna, po dwóch poziomach zasięgu.
func hostOkna(idOkna string) (string, error) {
	db := baza()
	if db == nil {
		return "", odmowaBrakuZasilenia()
	}
	const zapytanie = `
		SELECT COALESCE(u.wartosc, '')
		  FROM ustawienie u
		  JOIN poziom_zasiegu p ON p.id = u.poziom_zasiegu_id
		 WHERE u.klucz = 'host_wykonania'
		   AND u.os = 'platform' AND u.klucz_osi = ''
		   AND TRIM(COALESCE(u.wartosc, '')) <> ''
		   AND ((p.kod = 'okno' AND u.klucz_zasiegu = ?)
		     OR (p.kod = 'globalny' AND u.klucz_zasiegu = ''))
		 ORDER BY p.pierwszenstwo DESC
		 LIMIT 1`
	var nazwa string
	err := db.QueryRow(zapytanie, idOkna).Scan(&nazwa)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("zdalne: odczyt ustawienia host_wykonania okna %s: %w", idOkna, err)
	}
	return strings.TrimSpace(nazwa), nil
}

// hostZRejestru czyta wiersz hosta o podanej nazwie. Brak wiersza i brak zgody
// są odmowami trójczęściowymi — każda nazywa dokładnie ten ruch Operatora,
// który ją zdejmuje (instrukcja Danaco przy migracji 088).
func hostZRejestru(nazwa string) (Host, error) {
	db := baza()
	if db == nil {
		return Host{}, odmowaBrakuZasilenia()
	}
	const zapytanie = `SELECT id, nazwa, adres, uzytkownik, port, zgoda
	                     FROM host_zdalny WHERE nazwa = ?`
	var h Host
	var zgoda int
	err := db.QueryRow(zapytanie, nazwa).Scan(&h.Id, &h.Nazwa, &h.Adres, &h.Uzytkownik, &h.Port, &zgoda)
	if err == sql.ErrNoRows {
		return Host{}, fmt.Errorf("zdalne: połączenie z hostem %q nie zostało nawiązane, "+
			"bo hosta nie ma w wykazie hostów zdalnych (tabela host_zdalny, migracja 088); "+
			"wpisze go Operator instrukcją Danaco: INSERT INTO host_zdalny (nazwa, adres, "+
			"uzytkownik, zgoda, zgode_wydano) VALUES (…)", nazwa)
	}
	if err != nil {
		return Host{}, fmt.Errorf("zdalne: odczyt hosta %q z wykazu: %w", nazwa, err)
	}
	h.Zgoda = zgoda == 1
	if !h.Zgoda {
		return Host{}, fmt.Errorf("zdalne: połączenie z hostem %q nie zostało nawiązane, "+
			"bo Operator nie wydał zgody na inicjowanie połączeń z tą maszyną — to ochrona "+
			"jego maszyn, nie bramka; zgodę wydaje instrukcja Danaco: UPDATE host_zdalny "+
			"SET zgoda = 1, zgode_wydano = strftime('%%Y-%%m-%%dT%%H:%%M:%%fZ','now') "+
			"WHERE nazwa = '%s'", nazwa, nazwa)
	}
	return h, nil
}

// odmowaBrakuZasilenia opisuje stan przed wpięciem kompozycji — jednym zdaniem
// trójczęściowym, wspólnym dla wszystkich dróg pakietu.
func odmowaBrakuZasilenia() error {
	return fmt.Errorf("zdalne: tor do hosta zdalnego nie ma dostępu do bazy rdzenia, " +
		"bo kompozycja nie wywołała zdalne.Zasil(baza.DB) — wpięcie to jedna linia " +
		"w server/cmd/danaco-console/main.go po kontroli spójności bazy")
}
