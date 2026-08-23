// Odpowiedzialność pliku: dziennik działań asystenta — zapis wpisu, odczyt
// chronologiczny w oknie i odczyt spięty z jednym zleceniem. Zasila komendę
// `assistant.activity.list` (plik osobny od `asystent.go`).
//
// Wpis bez zlecenia jest wpisem prawdziwym, nie półwpisem. Asystent notuje też
// zdarzenia własne (`rodzaj = note`) bez zlecenia w tle — kolumna `zlecenie_kod`
// jest opcjonalna z migracji 050, więc metody tego pliku nie wymuszają
// powiązania i nie odrzucają wpisu, któremu go brakuje.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// WpisDziennikaAsystenta to wiersz tabeli `wpis_dziennika_asystenta` — jedna linia
// chronologicznej rozmowy widoczna w Actions Monitor.
type WpisDziennikaAsystenta struct {
	Kod              string
	OknoKod          string
	ZlecenieKod      *string
	Rodzaj           string
	Tresc            string
	NagranieOdnosnik *string
	Utworzono        int64
	// Wazny — wyróżnienie wpisu w Activity Feed (`assistant.activity.flag`,
	// migracja 296). Bez tej kolumny okno pokazywałoby wyróżnienie, którego
	// rdzeń nie pamięta.
	Wazny bool
	// NotatkaWyroznienia — powód wyróżnienia zapisany razem ze znacznikiem.
	NotatkaWyroznienia *string
}

const (
	kolumnyWpisuDziennika = `kod, okno_kod, zlecenie_kod, rodzaj, tresc, nagranie_odnosnik,
	                         utworzono, wazny, notatka_wyroznienia`

	wstawWpisDziennikaAsystenta = `INSERT INTO wpis_dziennika_asystenta
	                               (` + kolumnyWpisuDziennika + `)
	                               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzWpisDziennikaAsystenta = `SELECT ` + kolumnyWpisuDziennika + `
	                                 FROM wpis_dziennika_asystenta WHERE kod = ?`

	oznaczWpisDziennikaAsystenta = `UPDATE wpis_dziennika_asystenta
	                                SET wazny = ?, notatka_wyroznienia = ?
	                                WHERE kod = ?`

	pobierzWpisyOkna = `SELECT ` + kolumnyWpisuDziennika + `
	                    FROM wpis_dziennika_asystenta
	                    WHERE okno_kod = ?
	                    ORDER BY utworzono DESC, id DESC
	                    LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`

	pobierzWpisyZlecenia = `SELECT ` + kolumnyWpisuDziennika + `
	                        FROM wpis_dziennika_asystenta
	                        WHERE zlecenie_kod = ?
	                        ORDER BY utworzono DESC, id DESC`
)

// ZapiszWpis dopisuje jeden wpis dziennika. Metoda nie otwiera własnej
// transakcji — `PrzyjmijPolecenie` wywołuje ją wewnątrz transakcji
// wspólnej ze zleceniem, tu wpis stoi też samodzielnie dla wpisów własnych.
func (r *repozytoriumAsystenta) ZapiszWpis(ctx context.Context, wpis WpisDziennikaAsystenta) (WpisDziennikaAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wstawWpisDziennikaAsystenta)
	if err != nil {
		return WpisDziennikaAsystenta{}, err
	}
	_, err = polecenie.ExecContext(ctx, wpis.Kod, wpis.OknoKod, tekstDoKolumny(wpis.ZlecenieKod),
		wpis.Rodzaj, wpis.Tresc, tekstDoKolumny(wpis.NagranieOdnosnik), wpis.Utworzono,
		liczbaLogiczna(wpis.Wazny), tekstDoKolumny(wpis.NotatkaWyroznienia))
	if err != nil {
		return WpisDziennikaAsystenta{}, fmt.Errorf("dane: nie można zapisać wpisu dziennika asystenta %q: %w", wpis.Kod, err)
	}
	return wpis, nil
}

// Wpisy zwraca dziennik okna od najnowszego, ucięty granicą. Zero albo liczba
// ujemna znaczy „bez granicy" — Actions Monitor czasem chce całą rozmowę.
func (r *repozytoriumAsystenta) Wpisy(ctx context.Context, okno string, limit int) ([]WpisDziennikaAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWpisyOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, limit, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dziennika asystenta okna %q: %w", okno, err)
	}
	defer wiersze.Close()
	return zebrzWpisyDziennika(wiersze)
}

// WpisyZlecenia zwraca dziennik jednego zlecenia od najnowszego — bez granicy,
// bo rozmowa spięta z konkretnym zleceniem jest z natury skończona.
func (r *repozytoriumAsystenta) WpisyZlecenia(ctx context.Context, kodZlecenia string) ([]WpisDziennikaAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWpisyZlecenia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodZlecenia)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dziennika zlecenia %q: %w", kodZlecenia, err)
	}
	defer wiersze.Close()
	return zebrzWpisyDziennika(wiersze)
}

// zebrzWpisyDziennika przenosi wiersze zapytania do bytu obszaru; obie metody
// odczytu dzielą ten sam kształt kolumn, więc dzielą też przejście po wynikach.
func zebrzWpisyDziennika(wiersze *sql.Rows) ([]WpisDziennikaAsystenta, error) {
	wpisy := make([]WpisDziennikaAsystenta, 0, 32)
	for wiersze.Next() {
		wpis, err := odczytajWpisDziennika(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wpis dziennika asystenta: %w", err)
		}
		wpisy = append(wpisy, wpis)
	}
	return wpisy, wiersze.Err()
}

// odczytajWpisDziennika przenosi jeden wiersz do bytu obszaru.
func odczytajWpisDziennika(s skaner) (WpisDziennikaAsystenta, error) {
	var wpis WpisDziennikaAsystenta
	var zlecenieKod, nagranieOdnosnik, notatka sql.NullString
	var wazny int
	err := s.Scan(&wpis.Kod, &wpis.OknoKod, &zlecenieKod, &wpis.Rodzaj, &wpis.Tresc,
		&nagranieOdnosnik, &wpis.Utworzono, &wazny, &notatka)
	wpis.ZlecenieKod = tekstZKolumny(zlecenieKod)
	wpis.NagranieOdnosnik = tekstZKolumny(nagranieOdnosnik)
	wpis.Wazny = wazny == 1
	wpis.NotatkaWyroznienia = tekstZKolumny(notatka)
	return wpis, err
}

// OznaczWpisDziennika wyróżnia wpis dziennika albo zdejmuje to oznaczenie.
//
// Notatka pusta przy zdejmowaniu wyróżnienia kasuje też powód: powód bez
// znacznika byłby notatką do wpisu, którego nikt nie wyróżnił. Notatka pusta
// przy nadawaniu wyróżnienia zostawia powód zastany — wyróżnienie ponowione bez
// słowa nie ma prawa skasować zdania zapisanego wcześniej.
func (r *repozytoriumAsystenta) OznaczWpisDziennika(ctx context.Context, kod string,
	wazny bool, notatka *string) (WpisDziennikaAsystenta, error) {

	zapis, err := r.zapytania.przygotuj(ctx, oznaczWpisDziennikaAsystenta)
	if err != nil {
		return WpisDziennikaAsystenta{}, err
	}
	zastany, err := r.WpisDziennika(ctx, kod)
	if err != nil {
		return WpisDziennikaAsystenta{}, err
	}
	docelowa := notatka
	switch {
	case !wazny:
		docelowa = nil
	case notatka == nil || strings.TrimSpace(*notatka) == "":
		docelowa = zastany.NotatkaWyroznienia
	}
	if _, err := zapis.ExecContext(ctx, liczbaLogiczna(wazny), tekstDoKolumny(docelowa), kod); err != nil {
		return WpisDziennikaAsystenta{}, fmt.Errorf(
			"dane: nie można oznaczyć wpisu dziennika asystenta %q: %w", kod, err)
	}
	return r.WpisDziennika(ctx, kod)
}

// WpisDziennika zwraca jeden wpis dziennika po kodzie.
func (r *repozytoriumAsystenta) WpisDziennika(ctx context.Context,
	kod string) (WpisDziennikaAsystenta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWpisDziennikaAsystenta)
	if err != nil {
		return WpisDziennikaAsystenta{}, err
	}
	wpis, err := odczytajWpisDziennika(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return WpisDziennikaAsystenta{}, fmt.Errorf(
			"dane: wpis dziennika asystenta %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	return wpis, err
}
