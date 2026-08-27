// Plik prowadzi etykiety automatyki, zmienne przepływu i mapowanie danych między krokami oraz adnotacje kroków —
// notatkę i położenie węzła na kanwie; wszystko tutaj kluczuje się kodem kroku, nie kluczem wiersza kroku automatyki.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// ZmiennaAutomatyki to wiersz tabeli `zmienna_automatyki`. Odwołanie sekretu
// nigdy nie niesie wartości poświadczenia — wyłącznie referencję do skarbca.
type ZmiennaAutomatyki struct {
	Nazwa            string
	Rodzaj           string
	WartoscDomyslna  *string
	OdwolanieSekretu *string
	Kolejnosc        int
}

// MapowanieDanych to jeden łuk przepływu danych: wyjście kroku na wejście kroku następnego w automatyce.
type MapowanieDanych struct {
	KrokZ     string
	SciezkaZ  string
	KrokDo    string
	PoleDo    string
	Szablon   *string
	Kolejnosc int
}

// AdnotacjaKroku to notatka opisowa oraz położenie węzła kroku automatyki na kanwie edytora wizualnego.
type AdnotacjaKroku struct {
	KrokKod string
	Notatka *string
	X       int
	Y       int
}

const (
	usunEtykietyAutomatyki = `DELETE FROM etykieta_automatyki WHERE automatyka_id = ?`

	wstawEtykieteAutomatyki = `INSERT OR IGNORE INTO etykieta_automatyki
	                           (automatyka_id, etykieta) VALUES (?, ?)`

	listaEtykietAutomatyki = `SELECT etykieta FROM etykieta_automatyki
	                          WHERE automatyka_id = ? ORDER BY etykieta`

	usunZmienneAutomatyki = `DELETE FROM zmienna_automatyki WHERE automatyka_id = ?`

	wstawZmiennaAutomatyki = `INSERT INTO zmienna_automatyki
	                          (automatyka_id, nazwa, rodzaj, wartosc_domyslna,
	                           odwolanie_sekretu, kolejnosc)
	                          VALUES (?, ?, ?, ?, ?, ?)`

	listaZmiennychAutomatyki = `SELECT nazwa, rodzaj, wartosc_domyslna, odwolanie_sekretu, kolejnosc
	                            FROM zmienna_automatyki WHERE automatyka_id = ?
	                            ORDER BY kolejnosc, nazwa`

	usunMapowaniaAutomatyki = `DELETE FROM mapowanie_danych_automatyki WHERE automatyka_id = ?`

	wstawMapowanieAutomatyki = `INSERT INTO mapowanie_danych_automatyki
	                            (automatyka_id, krok_z, sciezka_z, krok_do, pole_do,
	                             szablon, kolejnosc)
	                            VALUES (?, ?, ?, ?, ?, ?, ?)`

	listaMapowanAutomatyki = `SELECT krok_z, sciezka_z, krok_do, pole_do, szablon, kolejnosc
	                          FROM mapowanie_danych_automatyki WHERE automatyka_id = ?
	                          ORDER BY kolejnosc, id`

	// Notatka i położenie zapisują się osobno, więc każdy zapis dotyka wyłącznie
	// swoich kolumn: ustawienie notatki nie przesuwa węzła, a przesunięcie węzła
	// nie kasuje notatki.
	zapiszNotatkeKroku = `INSERT INTO adnotacja_kroku_automatyki
	                      (automatyka_id, krok_kod, notatka) VALUES (?, ?, ?)
	                      ON CONFLICT(automatyka_id, krok_kod) DO UPDATE SET
	                          notatka = excluded.notatka`

	zapiszPolozenieKroku = `INSERT INTO adnotacja_kroku_automatyki
	                        (automatyka_id, krok_kod, wspolrzedna_x, wspolrzedna_y)
	                        VALUES (?, ?, ?, ?)
	                        ON CONFLICT(automatyka_id, krok_kod) DO UPDATE SET
	                            wspolrzedna_x = excluded.wspolrzedna_x,
	                            wspolrzedna_y = excluded.wspolrzedna_y`

	listaAdnotacjiKrokow = `SELECT krok_kod, notatka, wspolrzedna_x, wspolrzedna_y
	                        FROM adnotacja_kroku_automatyki WHERE automatyka_id = ?
	                        ORDER BY krok_kod`
)

// UstawEtykietyAutomatyki podmienia komplet etykiet. Wykaz pusty zdejmuje
// wszystkie — tak mówi kontrakt.
func (r *repozytoriumAutomatyk) UstawEtykietyAutomatyki(ctx context.Context,
	automatykaID int64, etykiety []string) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja,
			usunEtykietyAutomatyki, automatykaID); err != nil {
			return err
		}
		for _, etykieta := range etykiety {
			if etykieta == "" {
				continue
			}
			if err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja,
				wstawEtykieteAutomatyki, automatykaID, etykieta); err != nil {
				return err
			}
		}
		return nil
	})
}

// EtykietyAutomatyki zwraca wszystkie etykiety automatyki w porządku alfabetycznym wprost z bazy danych.
func (r *repozytoriumAutomatyk) EtykietyAutomatyki(ctx context.Context,
	automatykaID int64) ([]string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaEtykietAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać etykiet automatyki %d: %w", automatykaID, err)
	}
	defer wiersze.Close()

	etykiety := []string{}
	for wiersze.Next() {
		var etykieta string
		if err := wiersze.Scan(&etykieta); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna etykieta automatyki: %w", err)
		}
		etykiety = append(etykiety, etykieta)
	}
	return etykiety, wiersze.Err()
}

// ZapiszZmienneAutomatyki podmienia cały komplet zmiennych przepływu automatyki w jednej transakcji bazy.
func (r *repozytoriumAutomatyk) ZapiszZmienneAutomatyki(ctx context.Context,
	automatykaID int64, zmienne []ZmiennaAutomatyki) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja,
			usunZmienneAutomatyki, automatykaID); err != nil {
			return err
		}
		for numer, zmienna := range zmienne {
			if zmienna.Nazwa == "" {
				continue
			}
			err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja, wstawZmiennaAutomatyki,
				automatykaID, zmienna.Nazwa, zmienna.Rodzaj,
				tekstDoKolumny(zmienna.WartoscDomyslna),
				tekstDoKolumny(zmienna.OdwolanieSekretu), numer+1)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ZmienneAutomatyki zwraca wszystkie zmienne przepływu automatyki w kolejności ich zapisu do bazy danych.
func (r *repozytoriumAutomatyk) ZmienneAutomatyki(ctx context.Context,
	automatykaID int64) ([]ZmiennaAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZmiennychAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zmiennych automatyki %d: %w", automatykaID, err)
	}
	defer wiersze.Close()

	zmienne := []ZmiennaAutomatyki{}
	for wiersze.Next() {
		var zmienna ZmiennaAutomatyki
		var wartosc, odwolanie sql.NullString
		if err := wiersze.Scan(&zmienna.Nazwa, &zmienna.Rodzaj, &wartosc,
			&odwolanie, &zmienna.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zmiennej automatyki: %w", err)
		}
		zmienna.WartoscDomyslna = tekstZKolumny(wartosc)
		zmienna.OdwolanieSekretu = tekstZKolumny(odwolanie)
		zmienne = append(zmienne, zmienna)
	}
	return zmienne, wiersze.Err()
}

// ZapiszMapowaniaAutomatyki podmienia cały komplet mapowań danych automatyki w jednej transakcji bazy.
func (r *repozytoriumAutomatyk) ZapiszMapowaniaAutomatyki(ctx context.Context,
	automatykaID int64, mapowania []MapowanieDanych) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja,
			usunMapowaniaAutomatyki, automatykaID); err != nil {
			return err
		}
		for numer, mapowanie := range mapowania {
			if mapowanie.KrokZ == "" || mapowanie.KrokDo == "" {
				continue
			}
			err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja, wstawMapowanieAutomatyki,
				automatykaID, mapowanie.KrokZ, mapowanie.SciezkaZ, mapowanie.KrokDo,
				mapowanie.PoleDo, tekstDoKolumny(mapowanie.Szablon), numer+1)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// MapowaniaAutomatyki zwraca wszystkie mapowania danych automatyki w kolejności ich zapisu do bazy danych.
func (r *repozytoriumAutomatyk) MapowaniaAutomatyki(ctx context.Context,
	automatykaID int64) ([]MapowanieDanych, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaMapowanAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać mapowań automatyki %d: %w", automatykaID, err)
	}
	defer wiersze.Close()

	mapowania := []MapowanieDanych{}
	for wiersze.Next() {
		var mapowanie MapowanieDanych
		var szablon sql.NullString
		if err := wiersze.Scan(&mapowanie.KrokZ, &mapowanie.SciezkaZ, &mapowanie.KrokDo,
			&mapowanie.PoleDo, &szablon, &mapowanie.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz mapowania automatyki: %w", err)
		}
		mapowanie.Szablon = tekstZKolumny(szablon)
		mapowania = append(mapowania, mapowanie)
	}
	return mapowania, wiersze.Err()
}

// UstawNotatkeKroku zapisuje notatkę opisową przy kroku automatyki; treść pusta zdejmuje ją z tego kroku.
func (r *repozytoriumAutomatyk) UstawNotatkeKroku(ctx context.Context, automatykaID int64,
	krokKod string, notatka *string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszNotatkeKroku)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, automatykaID, krokKod, tekstDoKolumny(notatka)); err != nil {
		return fmt.Errorf("dane: nie można zapisać notatki kroku %q: %w", krokKod, err)
	}
	return nil
}

// ZapiszPolozeniaKrokow zapisuje położenia węzłów kanwy. Zapis dokłada
// i nadpisuje wskazane, a nie kasuje pozostałych: Workflow Builder przesuwa
// zwykle jeden węzeł, nie całą kanwę.
func (r *repozytoriumAutomatyk) ZapiszPolozeniaKrokow(ctx context.Context, automatykaID int64,
	polozenia []AdnotacjaKroku) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		for _, polozenie := range polozenia {
			if polozenie.KrokKod == "" {
				continue
			}
			err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja, zapiszPolozenieKroku,
				automatykaID, polozenie.KrokKod, polozenie.X, polozenie.Y)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// AdnotacjeKrokow zwraca wszystkie notatki oraz położenia węzłów kroków automatyki na jej kanwie edytora.
func (r *repozytoriumAutomatyk) AdnotacjeKrokow(ctx context.Context,
	automatykaID int64) ([]AdnotacjaKroku, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaAdnotacjiKrokow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać adnotacji kroków automatyki %d: %w",
			automatykaID, err)
	}
	defer wiersze.Close()

	adnotacje := []AdnotacjaKroku{}
	for wiersze.Next() {
		var adnotacja AdnotacjaKroku
		var notatka sql.NullString
		if err := wiersze.Scan(&adnotacja.KrokKod, &notatka, &adnotacja.X, &adnotacja.Y); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz adnotacji kroku: %w", err)
		}
		adnotacja.Notatka = tekstZKolumny(notatka)
		adnotacje = append(adnotacje, adnotacja)
	}
	return adnotacje, wiersze.Err()
}

// wykonajWTransakcji przygotowuje i wykonuje jedno polecenie w transakcji.
// Nazwa niesie przedrostek obszaru, bo pakiet `dane` ma jedną przestrzeń nazw
// dzieloną przez wszystkie moduły produktu.
func wykonajWTransakcjiAutomatyzacji(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	tekst string, argumenty ...any) error {

	polecenie, err := z.wTransakcji(ctx, transakcja, tekst)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, argumenty...); err != nil {
		return fmt.Errorf("dane: nie można wykonać zapisu obszaru automatyk: %w", err)
	}
	return nil
}
