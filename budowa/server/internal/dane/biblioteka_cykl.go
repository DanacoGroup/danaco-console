// Odpowiedzialność pliku: cykl życia zasobu (stan, miejsce w strukturze, nazwa,
// usunięcie trwałe) oraz pulpit stanu repozytorium — kolumny dołożone migracją
// 180 i agregaty liczone po całym zbiorze.
//
// Archiwizacja i przywrócenie są jedną czynnością o dwóch kierunkach, więc mają
// jedną metodę z parametrem stanu. Rozdzielenie ich dałoby dwa zapytania
// różniące się jednym łańcuchem.
//
// Usunięcie trwałe zdejmuje wiersz zasobu; wersje, etykiety i przypisania
// znikają kaskadą schematu (migracja 045). Bajty w magazynie treści zostają —
// sprząta je obchód magazynu przy starcie rdzenia
// (`core/adapter_modul_library_sprzatanie.go`), bo ta sama treść bywa
// współdzielona przez inny zasób pod tą samą sumą kontrolną.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// LiczbaWedlugKlucza to jedna pozycja rozkładu pulpitu stanu.
type LiczbaWedlugKlucza struct {
	Klucz  string
	Liczba int
}

// StatystykiBiblioteki to pulpit stanu repozytorium liczony po całym zbiorze.
type StatystykiBiblioteki struct {
	LiczbaZasobow       int
	LiczbaArchiwalnych  int
	LacznyRozmiar       int64
	LiczbaOsieroconych  int
	LiczbaDuplikatow    int
	BezSumyKontrolnej   int
	WedlugRodzajuTresci []LiczbaWedlugKlucza
	WedlugModuluZrodla  []LiczbaWedlugKlucza
	UzycieEtykiet       []LiczbaWedlugKlucza
	UzycieKolekcji      []LiczbaWedlugKlucza
}

const (
	ustawStanPlikuBiblioteki = `UPDATE plik_biblioteki
	                            SET stan = ?,
	                                zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                            WHERE identyfikator_zewnetrzny = ?`

	ustawSciezkePlikuBiblioteki = `UPDATE plik_biblioteki
	                               SET sciezka_repozytorium = ?,
	                                   zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                               WHERE identyfikator_zewnetrzny = ?`

	ustawNazwePlikuBiblioteki = `UPDATE plik_biblioteki
	                             SET nazwa = ?,
	                                 zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             WHERE identyfikator_zewnetrzny = ?`

	policzWersjePlikuBiblioteki = `SELECT COUNT(*) FROM wersja_pliku_biblioteki WHERE plik_id = ?`

	usunPlikBibliotekiTrwale = `DELETE FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`

	usunIndeksTresciPlikuBiblioteki = `DELETE FROM indeks_tresci_biblioteki WHERE rowid = ?`
)

// UstawStanPlikow przenosi zasoby między wykazem czynnym a archiwum i oddaje
// zasoby po zmianie. Zasób nieznany jest pomijany, a nie wywraca całości —
// odpowiedź niesie liczbę faktycznie przeniesionych.
func (r *repozytoriumBiblioteki) UstawStanPlikow(ctx context.Context, kody []string,
	stan string) ([]PlikBiblioteki, error) {

	return r.zmienPlikiZapytaniem(ctx, kody, ustawStanPlikuBiblioteki, stan,
		"nie można zmienić stanu pliku")
}

// UstawSciezkeRepozytorium przenosi zasoby pod wskazaną ścieżkę wewnątrz
// repozytorium. Treść i historia wersji zostają nietknięte — przeniesienie
// zmienia porządek, nie zasób.
func (r *repozytoriumBiblioteki) UstawSciezkeRepozytorium(ctx context.Context, kody []string,
	sciezka string) ([]PlikBiblioteki, error) {

	return r.zmienPlikiZapytaniem(ctx, kody, ustawSciezkePlikuBiblioteki, sciezka,
		"nie można przenieść pliku")
}

// PrzemianujPlik zmienia nazwę jednego zasobu — droga normalizacji nazw.
func (r *repozytoriumBiblioteki) PrzemianujPlik(ctx context.Context, kod,
	nazwa string) (PlikBiblioteki, error) {

	if strings.TrimSpace(nazwa) == "" {
		return PlikBiblioteki{}, fmt.Errorf("dane: zmiana nazwy pliku %q na pustą", kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, ustawNazwePlikuBiblioteki)
	if err != nil {
		return PlikBiblioteki{}, err
	}
	if _, err := polecenie.ExecContext(ctx, nazwa, kod); err != nil {
		return PlikBiblioteki{}, fmt.Errorf("dane: nie można zmienić nazwy pliku %q: %w", kod, err)
	}
	return r.Plik(ctx, kod)
}

// UsunPliki usuwa zasoby trwale wraz z ich wersjami i wpisem indeksu treści.
// Oddaje liczbę usuniętych zasobów i liczbę usuniętych wersji.
func (r *repozytoriumBiblioteki) UsunPliki(ctx context.Context, kody []string) (int, int, error) {
	usuniete, wersjeRazem := 0, 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		liczenieWersji, err := r.zapytania.wTransakcji(ctx, transakcja, policzWersjePlikuBiblioteki)
		if err != nil {
			return err
		}
		czyszczenieIndeksu, err := r.zapytania.wTransakcji(ctx, transakcja, usunIndeksTresciPlikuBiblioteki)
		if err != nil {
			return err
		}
		usuwanie, err := r.zapytania.wTransakcji(ctx, transakcja, usunPlikBibliotekiTrwale)
		if err != nil {
			return err
		}
		poszukiwanie, err := r.zapytania.wTransakcji(ctx, transakcja, idPlikuBibliotekiPoKodzie)
		if err != nil {
			return err
		}
		for _, kod := range kody {
			var plikID int64
			if err := poszukiwanie.QueryRowContext(ctx, kod).Scan(&plikID); err != nil {
				// Zasób nieznany nie wywraca usunięcia pozostałych: odpowiedź
				// mówi, ile naprawdę ubyło.
				continue
			}
			var wersje int
			if err := liczenieWersji.QueryRowContext(ctx, plikID).Scan(&wersje); err != nil {
				return fmt.Errorf("dane: nie można policzyć wersji pliku %q: %w", kod, err)
			}
			// Indeks treści FTS5 nie ma kluczy obcych, więc kaskada go nie
			// obejmuje — wiersz zdejmuje się wprost, inaczej po usuniętym
			// zasobie zostawałoby trafienie wyszukiwania.
			if _, err := czyszczenieIndeksu.ExecContext(ctx, plikID); err != nil {
				return fmt.Errorf("dane: nie można zdjąć indeksu treści pliku %q: %w", kod, err)
			}
			wynik, err := usuwanie.ExecContext(ctx, kod)
			if err != nil {
				return fmt.Errorf("dane: nie można usunąć pliku %q: %w", kod, err)
			}
			zdjete, err := wynik.RowsAffected()
			if err != nil {
				return fmt.Errorf("dane: nie można policzyć usuniętych plików: %w", err)
			}
			if zdjete > 0 {
				usuniete++
				wersjeRazem += wersje
			}
		}
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	return usuniete, wersjeRazem, nil
}

// Statystyki liczy pulpit stanu repozytorium. Wszystkie liczby idą z bazy —
// żadna nie jest oszacowaniem ani liczeniem po stronie rdzenia z odczytanej
// strony wykazu.
func (r *repozytoriumBiblioteki) Statystyki(ctx context.Context, kolekcjaKod, projektID *string,
	topN int) (StatystykiBiblioteki, error) {

	zawezenie, argumenty := zawezenieStatystykBiblioteki(kolekcjaKod, projektID)
	if topN <= 0 {
		topN = 20
	}
	var stat StatystykiBiblioteki

	// Nazwa tabeli stoi przy kolumnie stanu jawnie, bo część agregatów łączy
	// tabele i sama „stan" byłaby wtedy dwuznaczna.
	czynne := zawezenie + ` AND plik_biblioteki.stan = 'aktywny'`
	// Każdy agregat sumujący idzie przez COALESCE: SUM po zbiorze pustym daje
	// w SQLite NULL, a docelowe pola pulpitu są liczbami całkowitymi bez stanu
	// pustego. Bez tej osłony repozytorium puste — czyli rdzeń świeżo założony —
	// wywracało odczyt pulpitu zamiast oddać zera.
	wiersz := r.db.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(SUM(rozmiar_bajtow), 0),
	                                            COALESCE(SUM(CASE WHEN suma_kontrolna IS NULL
	                                                       OR suma_kontrolna = '' THEN 1 ELSE 0 END), 0)
	                                     FROM plik_biblioteki WHERE `+czynne, argumenty...)
	if err := wiersz.Scan(&stat.LiczbaZasobow, &stat.LacznyRozmiar, &stat.BezSumyKontrolnej); err != nil {
		return StatystykiBiblioteki{}, fmt.Errorf("dane: nie można policzyć zasobów repozytorium: %w", err)
	}

	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM plik_biblioteki
	                                  WHERE `+zawezenie+` AND stan = 'zarchiwizowany'`,
		argumenty...).Scan(&stat.LiczbaArchiwalnych)
	if err != nil {
		return StatystykiBiblioteki{}, fmt.Errorf("dane: nie można policzyć archiwum: %w", err)
	}

	// Zasób osierocony: bez etykiety i bez kolekcji — reguła audytu z opracowania
	// modułu, przeniesiona wprost do zapytania.
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM plik_biblioteki
	                                 WHERE `+czynne+`
	                                   AND NOT EXISTS (SELECT 1 FROM etykieta_pliku_biblioteki e
	                                                    WHERE e.plik_id = plik_biblioteki.id)
	                                   AND NOT EXISTS (SELECT 1 FROM przypisanie_kolekcji_biblioteki k
	                                                    WHERE k.plik_id = plik_biblioteki.id)`,
		argumenty...).Scan(&stat.LiczbaOsieroconych)
	if err != nil {
		return StatystykiBiblioteki{}, fmt.Errorf("dane: nie można policzyć zasobów osieroconych: %w", err)
	}

	// Duplikat dokładny: zasób dzielący sumę kontrolną z innym. Liczy się
	// zasoby należące do grup, nie same grupy — kontrakt pyta o zasoby.
	err = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(ile), 0) FROM (
	                                     SELECT COUNT(*) AS ile FROM plik_biblioteki
	                                     WHERE `+czynne+` AND suma_kontrolna IS NOT NULL
	                                       AND suma_kontrolna <> ''
	                                     GROUP BY suma_kontrolna HAVING COUNT(*) > 1)`,
		argumenty...).Scan(&stat.LiczbaDuplikatow)
	if err != nil {
		return StatystykiBiblioteki{}, fmt.Errorf("dane: nie można policzyć duplikatów: %w", err)
	}

	stat.WedlugRodzajuTresci, err = r.rozkladBiblioteki(ctx,
		`SELECT COALESCE(NULLIF(mime_type, ''), 'nieznany') AS klucz, COUNT(*) AS ile
		 FROM plik_biblioteki WHERE `+czynne+` GROUP BY klucz ORDER BY ile DESC, klucz LIMIT ?`,
		append(append([]any{}, argumenty...), topN))
	if err != nil {
		return StatystykiBiblioteki{}, err
	}
	stat.WedlugModuluZrodla, err = r.rozkladBiblioteki(ctx,
		`SELECT COALESCE(NULLIF(modul_zrodlowy_id, ''), 'nieznany') AS klucz, COUNT(*) AS ile
		 FROM plik_biblioteki WHERE `+czynne+` GROUP BY klucz ORDER BY ile DESC, klucz LIMIT ?`,
		append(append([]any{}, argumenty...), topN))
	if err != nil {
		return StatystykiBiblioteki{}, err
	}
	stat.UzycieEtykiet, err = r.rozkladBiblioteki(ctx,
		`SELECT e.etykieta AS klucz, COUNT(*) AS ile
		 FROM etykieta_pliku_biblioteki e
		 JOIN plik_biblioteki ON plik_biblioteki.id = e.plik_id
		 WHERE `+czynne+` GROUP BY klucz ORDER BY ile DESC, klucz LIMIT ?`,
		append(append([]any{}, argumenty...), topN))
	if err != nil {
		return StatystykiBiblioteki{}, err
	}
	stat.UzycieKolekcji, err = r.rozkladBiblioteki(ctx,
		`SELECT k.nazwa AS klucz, COUNT(*) AS ile
		 FROM przypisanie_kolekcji_biblioteki p
		 JOIN kolekcja_biblioteki k ON k.id = p.kolekcja_id
		 JOIN plik_biblioteki ON plik_biblioteki.id = p.plik_id
		 WHERE `+czynne+` GROUP BY klucz ORDER BY ile DESC, klucz LIMIT ?`,
		append(append([]any{}, argumenty...), topN))
	if err != nil {
		return StatystykiBiblioteki{}, err
	}
	return stat, nil
}

// zawezenieStatystykBiblioteki składa warunek wspólny wszystkim agregatom
// pulpitu. Nazwa tabeli stoi w warunku jawnie, bo część zapytań łączy tabele
// i „stan" bez wskazania tabeli byłby wtedy dwuznaczny.
func zawezenieStatystykBiblioteki(kolekcjaKod, projektID *string) (string, []any) {
	warunki := []string{"1 = 1"}
	argumenty := []any{}
	if projektID != nil && *projektID != "" {
		warunki = append(warunki, "plik_biblioteki.projekt_id = ?")
		argumenty = append(argumenty, *projektID)
	}
	if kolekcjaKod != nil && *kolekcjaKod != "" {
		warunki = append(warunki, `EXISTS (SELECT 1 FROM przypisanie_kolekcji_biblioteki pk
		                          JOIN kolekcja_biblioteki kb ON kb.id = pk.kolekcja_id
		                          WHERE pk.plik_id = plik_biblioteki.id
		                            AND kb.identyfikator_zewnetrzny = ?)`)
		argumenty = append(argumenty, *kolekcjaKod)
	}
	return strings.Join(warunki, " AND "), argumenty
}

// rozkladBiblioteki wykonuje zapytanie „klucz, liczba" i składa z niego rozkład.
func (r *repozytoriumBiblioteki) rozkladBiblioteki(ctx context.Context, zapytanie string,
	argumenty []any) ([]LiczbaWedlugKlucza, error) {

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można policzyć rozkładu repozytorium: %w", err)
	}
	defer wiersze.Close()

	lista := []LiczbaWedlugKlucza{}
	for wiersze.Next() {
		var pozycja LiczbaWedlugKlucza
		if err := wiersze.Scan(&pozycja.Klucz, &pozycja.Liczba); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz rozkładu repozytorium: %w", err)
		}
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt rozkładu repozytorium: %w", err)
	}
	return lista, nil
}

// zmienPlikiZapytaniem wykonuje jednakową zmianę na wielu zasobach i oddaje ich
// stan po zapisie — wspólne dla stanu i ścieżki repozytorium.
func (r *repozytoriumBiblioteki) zmienPlikiZapytaniem(ctx context.Context, kody []string,
	zapytanie, wartosc, powod string) ([]PlikBiblioteki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	zmienione := []PlikBiblioteki{}
	for _, kod := range kody {
		wynik, err := polecenie.ExecContext(ctx, wartosc, kod)
		if err != nil {
			return nil, fmt.Errorf("dane: %s %q: %w", powod, kod, err)
		}
		dotkniete, err := wynik.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("dane: nie można policzyć zmienionych plików: %w", err)
		}
		if dotkniete == 0 {
			continue
		}
		plik, err := r.Plik(ctx, kod)
		if err != nil {
			return nil, err
		}
		zmienione = append(zmienione, plik)
	}
	return zmienione, nil
}
