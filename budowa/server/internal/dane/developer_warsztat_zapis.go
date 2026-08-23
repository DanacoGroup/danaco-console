// Odpowiedzialność pliku: zapis warsztatu modułu Developer — punkty przerwania,
// kolekcje zapytań API, połączenia bazodanowe, przebiegi skanowania wraz ze
// znaleziskami oraz wyniki testów i pokrycie przebiegu budowania.
//
// Zapisy zbiorcze (znaleziska skanu, wyniki testów, pokrycie) idą jedną
// transakcją i zaczynają się od usunięcia poprzedniego pomiaru. Pomiar jest
// stanem z jednej chwili, nie przyrostem: dopisanie drugiego przebiegu do
// pierwszego dałoby wykaz, w którym ten sam test stoi dwa razy z dwoma różnymi
// wynikami i nie da się rozstrzygnąć, który jest dzisiejszy.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	wstawPunktPrzerwaniaDevelopera = `INSERT INTO developer_punkt_przerwania
	                                  (kod, okno_kod, sciezka, wiersz, rodzaj, warunek,
	                                   warunek_trafien, wpis, zweryfikowany)
	                                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                                  ON CONFLICT(okno_kod, sciezka, wiersz) DO UPDATE SET
	                                    rodzaj          = excluded.rodzaj,
	                                    warunek         = excluded.warunek,
	                                    warunek_trafien = excluded.warunek_trafien,
	                                    wpis            = excluded.wpis,
	                                    zweryfikowany   = excluded.zweryfikowany`

	usunPunktPrzerwaniaDevelopera = `DELETE FROM developer_punkt_przerwania
	                                 WHERE okno_kod = ? AND sciezka = ? AND wiersz = ?`

	wstawKolekcjeApiDevelopera = `INSERT INTO developer_kolekcja_api
	                              (kod, okno_kod, nazwa, zapytania, srodowiska, zmieniono)
	                              VALUES (?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                              ON CONFLICT(kod) DO UPDATE SET
	                                nazwa      = excluded.nazwa,
	                                zapytania  = excluded.zapytania,
	                                srodowiska = excluded.srodowiska,
	                                zmieniono  = excluded.zmieniono`

	wstawPolaczenieDanychDevelopera = `INSERT INTO developer_polaczenie_danych
	                                   (kod, okno_kod, nazwa, silnik, host, port, baza,
	                                    uzytkownik, poswiadczenie, tylko_odczyt, zmieniono)
	                                   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	                                           strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                                   ON CONFLICT(kod) DO UPDATE SET
	                                     nazwa         = excluded.nazwa,
	                                     silnik        = excluded.silnik,
	                                     host          = excluded.host,
	                                     port          = excluded.port,
	                                     baza          = excluded.baza,
	                                     uzytkownik    = excluded.uzytkownik,
	                                     poswiadczenie = excluded.poswiadczenie,
	                                     tylko_odczyt  = excluded.tylko_odczyt,
	                                     zmieniono     = excluded.zmieniono`

	wstawSkanDevelopera = `INSERT INTO developer_skan
	                       (kod, okno_kod, rodzaje, stan, znalezisk, zakonczono)
	                       VALUES (?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(kod) DO UPDATE SET
	                         stan       = excluded.stan,
	                         znalezisk  = excluded.znalezisk,
	                         zakonczono = excluded.zakonczono`

	wstawZnaleziskoDevelopera = `INSERT INTO developer_znalezisko
	                             (kod, skan_kod, rodzaj, waga, tytul, opis, sciezka, wiersz,
	                              regula, cve, pakiet, wersja_naprawy)
	                             VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	usunWynikiTestowDevelopera = `DELETE FROM developer_wynik_testu WHERE budowanie_kod = ?`

	wstawWynikTestuDevelopera = `INSERT INTO developer_wynik_testu
	                             (budowanie_kod, zestaw, nazwa, stan, czas_ms, tresc,
	                              sciezka, wiersz)
	                             VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	usunPokrycieDevelopera = `DELETE FROM developer_pokrycie WHERE budowanie_kod = ?`

	wstawPokrycieDevelopera = `INSERT INTO developer_pokrycie
	                           (budowanie_kod, sciezka, instrukcje, pokryte, procent,
	                            wiersze_bez_pokrycia)
	                           VALUES (?, ?, ?, ?, ?, ?)`
)

// ZapiszPunktPrzerwania zakłada albo odświeża punkt przerwania.
func (r *repozytoriumDevelopera) ZapiszPunktPrzerwania(ctx context.Context,
	punkt PunktPrzerwania) error {

	if punkt.Kod == "" || punkt.OknoKod == "" || punkt.Sciezka == "" {
		return fmt.Errorf("dane: punkt przerwania bez identyfikatora, okna albo ścieżki")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawPunktPrzerwaniaDevelopera)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, punkt.Kod, punkt.OknoKod, punkt.Sciezka,
		punkt.Wiersz, punkt.Rodzaj, punkt.Warunek, punkt.WarunekTrafien, punkt.Wpis,
		liczbaZPrawdy(punkt.Zweryfikowany)); err != nil {
		return fmt.Errorf("dane: nie można zapisać punktu przerwania %q: %w", punkt.Sciezka, err)
	}
	return nil
}

// UsunPunktPrzerwania zdejmuje punkt z wiersza pliku.
func (r *repozytoriumDevelopera) UsunPunktPrzerwania(ctx context.Context,
	oknoKod, sciezka string, wiersz int64) error {

	polecenie, err := r.zapytania.przygotuj(ctx, usunPunktPrzerwaniaDevelopera)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, oknoKod, sciezka, wiersz); err != nil {
		return fmt.Errorf("dane: nie można zdjąć punktu przerwania %q: %w", sciezka, err)
	}
	return nil
}

// ZapiszKolekcjeApi zakłada albo nadpisuje kolekcję zapytań.
func (r *repozytoriumDevelopera) ZapiszKolekcjeApi(ctx context.Context, kolekcja KolekcjaApi) error {
	if kolekcja.Kod == "" || kolekcja.OknoKod == "" || kolekcja.Nazwa == "" {
		return fmt.Errorf("dane: kolekcja zapytań bez identyfikatora, okna albo nazwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKolekcjeApiDevelopera)
	if err != nil {
		return err
	}
	zapytania := kolekcja.Zapytania
	if zapytania == "" {
		zapytania = "[]"
	}
	if _, err := polecenie.ExecContext(ctx, kolekcja.Kod, kolekcja.OknoKod, kolekcja.Nazwa,
		zapytania, kolekcja.Srodowiska); err != nil {
		return fmt.Errorf("dane: nie można zapisać kolekcji zapytań %q: %w", kolekcja.Nazwa, err)
	}
	return nil
}

// ZapiszPolaczenieDanych zakłada albo nadpisuje opis połączenia bazodanowego.
func (r *repozytoriumDevelopera) ZapiszPolaczenieDanych(ctx context.Context,
	polaczenie PolaczenieDanych) error {

	if polaczenie.Kod == "" || polaczenie.OknoKod == "" || polaczenie.Nazwa == "" {
		return fmt.Errorf("dane: połączenie bazodanowe bez identyfikatora, okna albo nazwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawPolaczenieDanychDevelopera)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, polaczenie.Kod, polaczenie.OknoKod, polaczenie.Nazwa,
		polaczenie.Silnik, polaczenie.Host, polaczenie.Port, polaczenie.Baza,
		polaczenie.Uzytkownik, polaczenie.Poswiadczenie,
		liczbaZPrawdy(polaczenie.TylkoOdczyt)); err != nil {
		return fmt.Errorf("dane: nie można zapisać połączenia %q: %w", polaczenie.Nazwa, err)
	}
	return nil
}

// ZapiszSkan zakłada albo domyka przebieg skanowania.
func (r *repozytoriumDevelopera) ZapiszSkan(ctx context.Context, skan PrzebiegSkanu) error {
	if skan.Kod == "" || skan.OknoKod == "" {
		return fmt.Errorf("dane: przebieg skanowania bez identyfikatora albo okna")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawSkanDevelopera)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, skan.Kod, skan.OknoKod, skan.Rodzaje, skan.Stan,
		skan.Znalezisk, skan.Zakonczono); err != nil {
		return fmt.Errorf("dane: nie można zapisać przebiegu skanowania %q: %w", skan.Kod, err)
	}
	return nil
}

// ZapiszZnaleziska dopisuje spostrzeżenia przebiegu skanowania jedną transakcją.
func (r *repozytoriumDevelopera) ZapiszZnaleziska(ctx context.Context,
	znaleziska []ZnaleziskoSkanu) error {

	if len(znaleziska) == 0 {
		return nil
	}
	return wTransakcji(ctx, r.zapytania.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawZnaleziskoDevelopera)
		if err != nil {
			return err
		}
		for _, z := range znaleziska {
			if _, err := polecenie.ExecContext(ctx, z.Kod, z.SkanKod, z.Rodzaj, z.Waga, z.Tytul,
				z.Opis, z.Sciezka, z.Wiersz, z.Regula, z.Cve, z.Pakiet,
				z.WersjaNaprawy); err != nil {
				return fmt.Errorf("dane: nie można zapisać znaleziska %q: %w", z.Tytul, err)
			}
		}
		return nil
	})
}

// ZapiszWynikiTestow zastępuje wyniki testów przebiegu budowania.
func (r *repozytoriumDevelopera) ZapiszWynikiTestow(ctx context.Context, budowanieKod string,
	wyniki []WynikTestu) error {

	if budowanieKod == "" {
		return fmt.Errorf("dane: wyniki testów bez identyfikatora przebiegu budowania")
	}
	return wTransakcji(ctx, r.zapytania.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunWynikiTestowDevelopera)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, budowanieKod); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić wyników testów %q: %w", budowanieKod, err)
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWynikTestuDevelopera)
		if err != nil {
			return err
		}
		for _, w := range wyniki {
			if _, err := polecenie.ExecContext(ctx, budowanieKod, w.Zestaw, w.Nazwa, w.Stan,
				w.CzasMs, w.Tresc, w.Sciezka, w.Wiersz); err != nil {
				return fmt.Errorf("dane: nie można zapisać wyniku testu %q: %w", w.Nazwa, err)
			}
		}
		return nil
	})
}

// ZapiszPokrycie zastępuje pomiar pokrycia przebiegu budowania.
func (r *repozytoriumDevelopera) ZapiszPokrycie(ctx context.Context, budowanieKod string,
	pokrycie []PokryciePliku) error {

	if budowanieKod == "" {
		return fmt.Errorf("dane: pokrycie bez identyfikatora przebiegu budowania")
	}
	return wTransakcji(ctx, r.zapytania.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunPokrycieDevelopera)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, budowanieKod); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić pokrycia %q: %w", budowanieKod, err)
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawPokrycieDevelopera)
		if err != nil {
			return err
		}
		for _, p := range pokrycie {
			if _, err := polecenie.ExecContext(ctx, budowanieKod, p.Sciezka, p.Instrukcje,
				p.Pokryte, p.Procent, p.WierszeBezPokrycia); err != nil {
				return fmt.Errorf("dane: nie można zapisać pokrycia %q: %w", p.Sciezka, err)
			}
		}
		return nil
	})
}

// liczbaZPrawdy przekłada prawdę na liczbę, którą SQLite trzyma zamiast typu
// logicznego.
func liczbaZPrawdy(prawda bool) int64 {
	if prawda {
		return 1
	}
	return 0
}
