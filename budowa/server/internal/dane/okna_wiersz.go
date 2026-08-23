// Odpowiedzialność pliku: przekład wiersza tabeli `okno_komunikacji` na
// strukturę Okno i z powrotem. Wartości słownikowe idą wyłącznie przez pakiet
// `shared`; brak wartości daje wartość domyślną, nie błąd.
package dane

import (
	"database/sql"

	"danacoconsole/shared"
)

// domyslnyTrybKomunikacji odpowiada wartości domyślnej kolumny w schemacie.
const domyslnyTrybKomunikacji = "tekst"

// wartosciOknaBazy to komplet wartości kolumn słownikowych okna.
type wartosciOknaBazy struct {
	srodowisko      string
	tryb            string
	rola            string
	trybKomunikacji string
	stan            string
}

// wartosciOkna przekłada pola wyliczeniowe struktury na wartości kolumn.
func wartosciOkna(okno Okno) (wartosciOknaBazy, error) {
	var wartosci wartosciOknaBazy
	var err error
	if wartosci.srodowisko, err = srodowiskoWykonaniaNaBaze(okno.SrodowiskoWykonania); err != nil {
		return wartosci, err
	}
	if wartosci.tryb, err = trybUprawnienNaBaze(okno.TrybUprawnien); err != nil {
		return wartosci, err
	}
	if wartosci.rola, err = rolaOknaNaBaze(okno.RolaOkna); err != nil {
		return wartosci, err
	}
	if wartosci.stan, err = stanOknaNaBaze(okno.Stan); err != nil {
		return wartosci, err
	}
	wartosci.trybKomunikacji = okno.TrybKomunikacji
	if wartosci.trybKomunikacji == "" {
		wartosci.trybKomunikacji = domyslnyTrybKomunikacji
	}
	return wartosci, nil
}

// odczytajOkno składa strukturę z jednego wiersza wyniku. Katalogi robocze
// dokłada repozytorium — leżą w tabeli podrzędnej.
func odczytajOkno(wiersz skaner) (Okno, error) {
	var okno Okno
	var tytul, identyfikator, agent sql.NullString
	var koordynator sql.NullInt64
	var srodowisko, tryb, rola, stan string
	err := wiersz.Scan(&okno.ID, &okno.SesjaID, &okno.ModulID, &okno.KanalModeluID, &tytul,
		&srodowisko, &tryb, &rola, &koordynator, &okno.TrybKomunikacji, &stan,
		&identyfikator, &agent, &okno.Kolejnosc, &okno.Utworzono, &okno.Zaktualizowano)
	if err != nil {
		return Okno{}, err
	}
	okno.Tytul = tekstZKolumny(tytul)
	okno.IdentyfikatorZewnetrzny = tekstZKolumny(identyfikator)
	okno.AgentKod = tekstZKolumny(agent)
	okno.OknoKoordynatoraID = liczbaZKolumny(koordynator)
	okno.KatalogiRobocze = []string{}
	return uzupelnijSlownikiOkna(okno, srodowisko, tryb, rola, stan)
}

// uzupelnijSlownikiOkna przekłada wartości kolumn słownikowych na kontrakt.
func uzupelnijSlownikiOkna(okno Okno, srodowisko, tryb, rola, stan string) (Okno, error) {
	var err error
	if okno.SrodowiskoWykonania, err = srodowiskoWykonaniaZBazy(srodowisko); err != nil {
		return Okno{}, err
	}
	if okno.TrybUprawnien, err = trybUprawnienZBazy(tryb); err != nil {
		return Okno{}, err
	}
	if okno.RolaOkna, err = rolaOknaZBazy(rola); err != nil {
		return Okno{}, err
	}
	if okno.Stan, err = stanOknaZBazy(stan); err != nil {
		return Okno{}, err
	}
	return okno, nil
}

// NoweOkno zwraca okno o wartościach domyślnych kontraktu. Warstwa wyższa
// nadpisuje wyłącznie pola ustawione jawnie; pole pominięte zostaje przy
// wartości domyślnej.
func NoweOkno(sesjaID, modulID, kanalID int64) Okno {
	return Okno{
		SesjaID:             sesjaID,
		ModulID:             modulID,
		KanalModeluID:       kanalID,
		KatalogiRobocze:     []string{},
		SrodowiskoWykonania: shared.ExecutionEnvLocal,
		TrybUprawnien:       shared.PermissionModeManual,
		RolaOkna:            shared.WindowRoleStandalone,
		TrybKomunikacji:     domyslnyTrybKomunikacji,
		Stan:                shared.WindowStatusOpen,
	}
}
