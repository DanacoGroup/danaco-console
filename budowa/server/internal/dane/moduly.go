// Odpowiedzialność pliku: dostęp do słownika modułów (tabele `modul`
// i `srodowisko_modul`). Moduł jest jednostką funkcjonalną osadzaną w środowiskach i należy do okna
// komunikacji, nie do sesji.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Modul to wiersz tabeli `modul`, niosący nazwę, opis, ikonę i rodzaj modułu platformy w całym systemie.
type Modul struct {
	ID    int64
	Kod   string
	Nazwa string
	Opis  *string
	// Ikona to nazwa ikony interfejsu; pusta wartość znaczy ikonę zastępczą dobieraną przez klienta.
	Ikona *string
	// Rodzaj rozstrzyga, czym moduł jest — wartość ustalana treścią migracji, nie warunkiem schematu.
	Rodzaj *string
	// KonfigurowanyNaStronieGlownej mówi, czy moduł nastawia się w Strefie 2 Strony głównej.
	KonfigurowanyNaStronieGlownej bool
	Aktywny                       bool
}

// RepozytoriumModulow jest kontraktem słownika modułów, określającym operacje dostępne na całym wykazie.
type RepozytoriumModulow interface {
	Lista(ctx context.Context) ([]Modul, error)
	PoKodzie(ctx context.Context, kod string) (Modul, error)
	ListaSrodowiska(ctx context.Context, srodowiskoID int64) ([]Modul, error)
}

const (
	kolumnyModulu = `m.id, m.kod, m.nazwa, m.opis, m.ikona, m.rodzaj,
	                 m.konfigurowany_na_stronie_glownej, m.aktywny`

	listaModulow = `SELECT ` + kolumnyModulu + ` FROM modul m ORDER BY m.kod`

	modulPoKodzie = `SELECT ` + kolumnyModulu + ` FROM modul m WHERE m.kod = ?`

	// Widoczność modułu w środowisku niesie macierz `srodowisko_modul`; kolejność
	// jej wiersza jest kolejnością pozycji w nawigacji bocznej.
	listaModulowSrodowiska = `SELECT ` + kolumnyModulu + ` FROM modul m
	                          JOIN srodowisko_modul sm ON sm.modul_id = m.id
	                          WHERE sm.srodowisko_id = ? AND sm.widoczny = 1
	                          ORDER BY sm.kolejnosc, m.kod`
)

type repozytoriumModulow struct {
	zapytania *zapytania
}

func noweRepozytoriumModulow(z *zapytania) *repozytoriumModulow {
	return &repozytoriumModulow{zapytania: z}
}

// Lista zwraca wszystkie moduły platformy zapisane w słowniku modułów systemu rdzenia całej tej aplikacji.
func (r *repozytoriumModulow) Lista(ctx context.Context) ([]Modul, error) {
	return r.wykaz(ctx, listaModulow, "modułów")
}

// ListaSrodowiska zwraca moduły widoczne w środowisku, w kolejności ustalonej nawigacją boczną interfejsu.
func (r *repozytoriumModulow) ListaSrodowiska(ctx context.Context, srodowiskoID int64) ([]Modul, error) {
	return r.wykaz(ctx, listaModulowSrodowiska, "modułów środowiska", srodowiskoID)
}

// PoKodzie zwraca moduł ze słownika modułów wskazany jego unikalnym kodem tekstowym zapisanym w bazie.
func (r *repozytoriumModulow) PoKodzie(ctx context.Context, kod string) (Modul, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, modulPoKodzie)
	if err != nil {
		return Modul{}, err
	}
	modul, err := odczytajModul(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return Modul{}, fmt.Errorf("%w: moduł o kodzie %q", ErrBrakWiersza, kod)
	}
	return modul, err
}

// wykaz wykonuje zapytanie do bazy danych zwracające wiele wierszy słownika modułów całej tej platformy.
func (r *repozytoriumModulow) wykaz(ctx context.Context, zapytanie, opis string, argumenty ...any) ([]Modul, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać %s: %w", opis, err)
	}
	defer wiersze.Close()

	lista := []Modul{}
	for wiersze.Next() {
		modul, err := odczytajModul(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, modul)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt %s: %w", opis, err)
	}
	return lista, nil
}

// odczytajModul składa pełną strukturę modułu z jednego wiersza wyniku zapytania do bazy danych rdzenia.
func odczytajModul(wiersz skaner) (Modul, error) {
	var modul Modul
	var opis, ikona, rodzaj sql.NullString
	var naStronieGlownej, aktywny int
	if err := wiersz.Scan(&modul.ID, &modul.Kod, &modul.Nazwa, &opis, &ikona, &rodzaj,
		&naStronieGlownej, &aktywny); err != nil {
		return Modul{}, err
	}
	modul.Opis = tekstZKolumny(opis)
	modul.Ikona = tekstZKolumny(ikona)
	modul.Rodzaj = tekstZKolumny(rodzaj)
	modul.KonfigurowanyNaStronieGlownej = naStronieGlownej != 0
	modul.Aktywny = aktywny != 0
	return modul, nil
}
