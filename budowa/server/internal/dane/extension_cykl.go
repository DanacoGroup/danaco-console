// Odpowiedzialność pliku: rodzina `extension.*` — cykl życia pozycji katalogu.
// Kolekcje kuratorskie, dziennik cyklu życia, wersje pozycji wraz z przypięciem
// oraz paczki przesłane instalacją Personal
// (`store/migracja_207_rozszerzenia_cykl_zycia.sql`).
//
// Interfejs `RepozytoriumRozszerzen` deklaruje `extension.go`; ten plik
// i trzy sąsiednie (`extension_protokol.go`, `extension_integracje.go`,
// `extension_zaufanie.go`) dokładają mu metody — jedno repozytorium, cztery
// pliki wedle odpowiedzialności, tak jak w obszarze Apps.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// KolekcjaRozszerzen to wiersz tabeli `kolekcja_rozszerzen` wraz z kodami
// pozycji, które ją tworzą — kontrakt oddaje kolekcję zawsze razem z nimi.
type KolekcjaRozszerzen struct {
	ID               int64
	Kod              string
	Nazwa            string
	Opis             *string
	OznaczenieBarwne *string
	KodyPozycji      []string
	Zaktualizowano   int64
}

// WpisHistoriiRozszerzenia to wiersz tabeli `historia_rozszerzenia`.
type WpisHistoriiRozszerzenia struct {
	ID              int64
	Kod             string
	RozszerzenieKod string
	Czynnosc        string
	WersjaPrzed     *string
	WersjaPo        *string
	Szczegol        *string
	Zaszlo          int64
}

// WersjaRozszerzenia to wiersz tabeli `wersja_rozszerzenia`.
type WersjaRozszerzenia struct {
	ID              int64
	RozszerzenieKod string
	Wersja          string
	DziennikZmian   *string
	PaczkaOdwolanie *string
	Utworzono       int64
}

// PaczkaRozszerzenia to wiersz tabeli `paczka_rozszerzenia` — bajty przesłane
// instalacją Personal leżą pod `Sciezka` w magazynie treści rdzenia.
type PaczkaRozszerzenia struct {
	ID            int64
	Kod           string
	NazwaPliku    string
	Sciezka       string
	Rozmiar       int64
	SumaKontrolna string
	Utworzono     int64
}

const (
	kolumnyKolekcjiRozszerzen = `id, identyfikator_zewnetrzny, nazwa, opis,
	                             oznaczenie_barwne, zaktualizowano`

	zapiszKolekcjeRozszerzen = `INSERT INTO kolekcja_rozszerzen
	                            (identyfikator_zewnetrzny, nazwa, opis, oznaczenie_barwne, zaktualizowano)
	                            VALUES (?, ?, ?, ?, ?)
	                            ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                nazwa = excluded.nazwa,
	                                opis = excluded.opis,
	                                oznaczenie_barwne = excluded.oznaczenie_barwne,
	                                zaktualizowano = excluded.zaktualizowano`

	pobierzKolekcjeRozszerzen = `SELECT ` + kolumnyKolekcjiRozszerzen + `
	                             FROM kolekcja_rozszerzen WHERE identyfikator_zewnetrzny = ?`

	listaKolekcjiRozszerzen = `SELECT ` + kolumnyKolekcjiRozszerzen + `
	                           FROM kolekcja_rozszerzen ORDER BY nazwa, id`

	usunPozycjeKolekcjiRozszerzen = `DELETE FROM pozycja_kolekcji_rozszerzen WHERE kolekcja_id = ?`

	wstawPozycjeKolekcjiRozszerzen = `INSERT INTO pozycja_kolekcji_rozszerzen
	                                  (kolekcja_id, rozszerzenie_kod, kolejnosc)
	                                  VALUES (?, ?, ?)
	                                  ON CONFLICT(kolekcja_id, rozszerzenie_kod) DO NOTHING`

	listaPozycjiKolekcjiRozszerzen = `SELECT k.identyfikator_zewnetrzny, p.rozszerzenie_kod
	                                  FROM pozycja_kolekcji_rozszerzen p
	                                  JOIN kolekcja_rozszerzen k ON k.id = p.kolekcja_id
	                                  ORDER BY p.kolejnosc, p.rozszerzenie_kod`

	wstawHistorieRozszerzenia = `INSERT INTO historia_rozszerzenia
	                             (identyfikator_zewnetrzny, rozszerzenie_kod, czynnosc,
	                              wersja_przed, wersja_po, szczegol, zaszlo)
	                             VALUES (?, ?, ?, ?, ?, ?, ?)`

	listaHistoriiRozszerzenia = `SELECT id, identyfikator_zewnetrzny, rozszerzenie_kod, czynnosc,
	                                    wersja_przed, wersja_po, szczegol, zaszlo
	                             FROM historia_rozszerzenia
	                             WHERE (? = '' OR rozszerzenie_kod = ?) AND zaszlo >= ?
	                             ORDER BY zaszlo DESC, id DESC LIMIT ?`

	policzHistorieRozszerzenia = `SELECT COUNT(*) FROM historia_rozszerzenia
	                              WHERE (? = '' OR rozszerzenie_kod = ?) AND zaszlo >= ?`

	zapiszWersjeRozszerzenia = `INSERT INTO wersja_rozszerzenia
	                            (rozszerzenie_kod, wersja, dziennik_zmian, paczka_odwolanie, utworzono)
	                            VALUES (?, ?, ?, ?, ?)
	                            ON CONFLICT(rozszerzenie_kod, wersja) DO UPDATE SET
	                                dziennik_zmian = IFNULL(excluded.dziennik_zmian,
	                                                        wersja_rozszerzenia.dziennik_zmian),
	                                paczka_odwolanie = IFNULL(excluded.paczka_odwolanie,
	                                                          wersja_rozszerzenia.paczka_odwolanie)`

	listaWersjiRozszerzenia = `SELECT id, rozszerzenie_kod, wersja, dziennik_zmian,
	                                  paczka_odwolanie, utworzono
	                           FROM wersja_rozszerzenia WHERE rozszerzenie_kod = ?
	                           ORDER BY utworzono DESC, id DESC`

	przypnijWersjeRozszerzenia = `UPDATE rozszerzenie SET wersja_przypieta = ?, zaktualizowano = ?
	                              WHERE identyfikator_zewnetrzny = ?`

	pobierzPrzypiecieRozszerzenia = `SELECT IFNULL(wersja_przypieta,'') FROM rozszerzenie
	                                 WHERE identyfikator_zewnetrzny = ?`

	wstawPaczkeRozszerzenia = `INSERT INTO paczka_rozszerzenia
	                           (identyfikator_zewnetrzny, nazwa_pliku, sciezka, rozmiar,
	                            suma_kontrolna, utworzono)
	                           VALUES (?, ?, ?, ?, ?, ?)`

	pobierzPaczkeRozszerzenia = `SELECT id, identyfikator_zewnetrzny, nazwa_pliku, sciezka,
	                                    rozmiar, suma_kontrolna, utworzono
	                             FROM paczka_rozszerzenia WHERE identyfikator_zewnetrzny = ?`
)

// ZapiszKolekcjeRozszerzen zapisuje kolekcję wraz z kompletem jej pozycji
// w jednej transakcji — kontrakt nadsyła `extensionIds` bez trybu częściowej
// zmiany, więc związek wymienia się „usuń, wstaw od nowa".
func (r *repozytoriumRozszerzen) ZapiszKolekcjeRozszerzen(ctx context.Context,
	kolekcja KolekcjaRozszerzen) (KolekcjaRozszerzen, error) {

	if kolekcja.Kod == "" || kolekcja.Nazwa == "" {
		return KolekcjaRozszerzen{}, fmt.Errorf("dane: kolekcja rozszerzeń bez identyfikatora albo nazwy")
	}
	err := wTransakcji(ctx, r.baza, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszKolekcjeRozszerzen)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, kolekcja.Kod, kolekcja.Nazwa,
			tekstDoKolumny(kolekcja.Opis), tekstDoKolumny(kolekcja.OznaczenieBarwne),
			kolekcja.Zaktualizowano); err != nil {
			return fmt.Errorf("dane: nie można zapisać kolekcji %q: %w", kolekcja.Kod, err)
		}

		var kolekcjaID int64
		wiersz := transakcja.QueryRowContext(ctx,
			`SELECT id FROM kolekcja_rozszerzen WHERE identyfikator_zewnetrzny = ?`, kolekcja.Kod)
		if err := wiersz.Scan(&kolekcjaID); err != nil {
			return fmt.Errorf("dane: nie można odczytać id kolekcji %q: %w", kolekcja.Kod, err)
		}

		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunPozycjeKolekcjiRozszerzen)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, kolekcjaID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić pozycji kolekcji %q: %w", kolekcja.Kod, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawPozycjeKolekcjiRozszerzen)
		if err != nil {
			return err
		}
		for kolejnosc, pozycja := range kolekcja.KodyPozycji {
			if pozycja == "" {
				continue
			}
			if _, err := wstawienie.ExecContext(ctx, kolekcjaID, pozycja, kolejnosc); err != nil {
				return fmt.Errorf("dane: nie można zapisać pozycji %q kolekcji %q: %w",
					pozycja, kolekcja.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return KolekcjaRozszerzen{}, err
	}
	return r.KolekcjaRozszerzen(ctx, kolekcja.Kod)
}

// KolekcjaRozszerzen zwraca jedną kolekcję wraz z kodami jej pozycji.
func (r *repozytoriumRozszerzen) KolekcjaRozszerzen(ctx context.Context, kod string) (KolekcjaRozszerzen, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKolekcjeRozszerzen)
	if err != nil {
		return KolekcjaRozszerzen{}, err
	}
	kolekcja, err := odczytajKolekcjeRozszerzen(polecenie.QueryRowContext(ctx, kod))
	if err == sql.ErrNoRows {
		return KolekcjaRozszerzen{}, ErrBrakWiersza
	}
	if err != nil {
		return KolekcjaRozszerzen{}, fmt.Errorf("dane: nieczytelna kolekcja %q: %w", kod, err)
	}
	pozycje, err := r.pozycjeKolekcjiRozszerzen(ctx)
	if err != nil {
		return KolekcjaRozszerzen{}, err
	}
	kolekcja.KodyPozycji = pozycje[kolekcja.Kod]
	return kolekcja, nil
}

// KolekcjeRozszerzen zwraca wszystkie kolekcje wraz z ich pozycjami.
func (r *repozytoriumRozszerzen) KolekcjeRozszerzen(ctx context.Context) ([]KolekcjaRozszerzen, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKolekcjiRozszerzen)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kolekcji rozszerzeń: %w", err)
	}
	defer wiersze.Close()

	lista := []KolekcjaRozszerzen{}
	for wiersze.Next() {
		kolekcja, err := odczytajKolekcjeRozszerzen(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kolekcji rozszerzeń: %w", err)
		}
		lista = append(lista, kolekcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kolekcji rozszerzeń: %w", err)
	}
	pozycje, err := r.pozycjeKolekcjiRozszerzen(ctx)
	if err != nil {
		return nil, err
	}
	for indeks := range lista {
		lista[indeks].KodyPozycji = pozycje[lista[indeks].Kod]
	}
	return lista, nil
}

// pozycjeKolekcjiRozszerzen zwraca mapę kod kolekcji → kody pozycji. Jedno
// zapytanie na wszystkie kolekcje: wykaz ciągnąłby inaczej tyle zapytań, ile ma
// pozycji.
func (r *repozytoriumRozszerzen) pozycjeKolekcjiRozszerzen(ctx context.Context) (map[string][]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPozycjiKolekcjiRozszerzen)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pozycji kolekcji: %w", err)
	}
	defer wiersze.Close()

	mapa := map[string][]string{}
	for wiersze.Next() {
		var kolekcja, pozycja string
		if err := wiersze.Scan(&kolekcja, &pozycja); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pozycji kolekcji: %w", err)
		}
		mapa[kolekcja] = append(mapa[kolekcja], pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pozycji kolekcji: %w", err)
	}
	return mapa, nil
}

// DopiszHistorieRozszerzenia odnotowuje jedno zdarzenie cyklu życia pozycji.
func (r *repozytoriumRozszerzen) DopiszHistorieRozszerzenia(ctx context.Context,
	wpis WpisHistoriiRozszerzenia) error {

	if wpis.Kod == "" || wpis.RozszerzenieKod == "" || wpis.Czynnosc == "" {
		return fmt.Errorf("dane: wpis historii rozszerzenia bez identyfikatora, pozycji albo czynności")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawHistorieRozszerzenia)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wpis.Kod, wpis.RozszerzenieKod, wpis.Czynnosc,
		tekstDoKolumny(wpis.WersjaPrzed), tekstDoKolumny(wpis.WersjaPo),
		tekstDoKolumny(wpis.Szczegol), wpis.Zaszlo)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać wpisu historii %q: %w", wpis.Kod, err)
	}
	return nil
}

// HistoriaRozszerzenia zwraca stronę dziennika cyklu życia wraz z liczbą
// wszystkich wpisów spełniających te same warunki.
func (r *repozytoriumRozszerzen) HistoriaRozszerzenia(ctx context.Context, rozszerzenie string,
	od int64, granica int) ([]WpisHistoriiRozszerzenia, int, error) {

	if granica <= 0 {
		granica = 200
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaHistoriiRozszerzenia)
	if err != nil {
		return nil, 0, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rozszerzenie, rozszerzenie, od, granica)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać historii rozszerzeń: %w", err)
	}
	defer wiersze.Close()

	lista := []WpisHistoriiRozszerzenia{}
	for wiersze.Next() {
		var wpis WpisHistoriiRozszerzenia
		var przed, po, szczegol sql.NullString
		err := wiersze.Scan(&wpis.ID, &wpis.Kod, &wpis.RozszerzenieKod, &wpis.Czynnosc,
			&przed, &po, &szczegol, &wpis.Zaszlo)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz historii rozszerzenia: %w", err)
		}
		wpis.WersjaPrzed = tekstZKolumny(przed)
		wpis.WersjaPo = tekstZKolumny(po)
		wpis.Szczegol = tekstZKolumny(szczegol)
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt historii rozszerzeń: %w", err)
	}

	liczenie, err := r.zapytania.przygotuj(ctx, policzHistorieRozszerzenia)
	if err != nil {
		return nil, 0, err
	}
	var razem int
	if err := liczenie.QueryRowContext(ctx, rozszerzenie, rozszerzenie, od).Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć wpisów historii: %w", err)
	}
	return lista, razem, nil
}

// ZapiszWersjeRozszerzenia odnotowuje wersję pozycji katalogu.
func (r *repozytoriumRozszerzen) ZapiszWersjeRozszerzenia(ctx context.Context,
	wersja WersjaRozszerzenia) error {

	if wersja.RozszerzenieKod == "" || wersja.Wersja == "" {
		return fmt.Errorf("dane: wersja rozszerzenia bez pozycji albo numeru")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWersjeRozszerzenia)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wersja.RozszerzenieKod, wersja.Wersja,
		tekstDoKolumny(wersja.DziennikZmian), tekstDoKolumny(wersja.PaczkaOdwolanie),
		wersja.Utworzono)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać wersji %q pozycji %q: %w",
			wersja.Wersja, wersja.RozszerzenieKod, err)
	}
	return nil
}

// WersjeRozszerzenia zwraca wersje pozycji, od najnowszej.
func (r *repozytoriumRozszerzen) WersjeRozszerzenia(ctx context.Context,
	rozszerzenie string) ([]WersjaRozszerzenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWersjiRozszerzenia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rozszerzenie)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wersji pozycji %q: %w", rozszerzenie, err)
	}
	defer wiersze.Close()

	lista := []WersjaRozszerzenia{}
	for wiersze.Next() {
		var wersja WersjaRozszerzenia
		var dziennik, paczka sql.NullString
		err := wiersze.Scan(&wersja.ID, &wersja.RozszerzenieKod, &wersja.Wersja,
			&dziennik, &paczka, &wersja.Utworzono)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji rozszerzenia: %w", err)
		}
		wersja.DziennikZmian = tekstZKolumny(dziennik)
		wersja.PaczkaOdwolanie = tekstZKolumny(paczka)
		lista = append(lista, wersja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wersji pozycji %q: %w", rozszerzenie, err)
	}
	return lista, nil
}

// PrzypnijWersjeRozszerzenia ustawia albo zdejmuje przypięcie wersji pozycji.
// Pusty numer zdejmuje przypięcie — kontrakt mówi wprost, że brak `version`
// „zdejmuje przypięcie".
func (r *repozytoriumRozszerzen) PrzypnijWersjeRozszerzenia(ctx context.Context,
	rozszerzenie, wersja string, teraz int64) error {

	polecenie, err := r.zapytania.przygotuj(ctx, przypnijWersjeRozszerzenia)
	if err != nil {
		return err
	}
	var wartosc any
	if wersja != "" {
		wartosc = wersja
	}
	if _, err := polecenie.ExecContext(ctx, wartosc, teraz, rozszerzenie); err != nil {
		return fmt.Errorf("dane: nie można przypiąć wersji pozycji %q: %w", rozszerzenie, err)
	}
	return nil
}

// PrzypieciaWersjiRozszerzenia zwraca numer wersji przypiętej; pusty znaczy
// „bez przypięcia".
func (r *repozytoriumRozszerzen) PrzypiecieWersjiRozszerzenia(ctx context.Context,
	rozszerzenie string) (string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPrzypiecieRozszerzenia)
	if err != nil {
		return "", err
	}
	var wersja string
	err = polecenie.QueryRowContext(ctx, rozszerzenie).Scan(&wersja)
	if err == sql.ErrNoRows {
		return "", ErrBrakWiersza
	}
	if err != nil {
		return "", fmt.Errorf("dane: nie można odczytać przypięcia pozycji %q: %w", rozszerzenie, err)
	}
	return wersja, nil
}

// ZalozPaczkeRozszerzenia zapisuje wiersz paczki przesłanej instalacją.
func (r *repozytoriumRozszerzen) ZalozPaczkeRozszerzenia(ctx context.Context,
	paczka PaczkaRozszerzenia) error {

	if paczka.Kod == "" || paczka.Sciezka == "" {
		return fmt.Errorf("dane: paczka rozszerzenia bez identyfikatora albo ścieżki")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawPaczkeRozszerzenia)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, paczka.Kod, paczka.NazwaPliku, paczka.Sciezka,
		paczka.Rozmiar, paczka.SumaKontrolna, paczka.Utworzono)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać paczki %q: %w", paczka.Kod, err)
	}
	return nil
}

// PaczkaRozszerzenia zwraca wiersz paczki po jej odwołaniu.
func (r *repozytoriumRozszerzen) PaczkaRozszerzenia(ctx context.Context, kod string) (PaczkaRozszerzenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPaczkeRozszerzenia)
	if err != nil {
		return PaczkaRozszerzenia{}, err
	}
	var paczka PaczkaRozszerzenia
	err = polecenie.QueryRowContext(ctx, kod).Scan(&paczka.ID, &paczka.Kod, &paczka.NazwaPliku,
		&paczka.Sciezka, &paczka.Rozmiar, &paczka.SumaKontrolna, &paczka.Utworzono)
	if err == sql.ErrNoRows {
		return PaczkaRozszerzenia{}, ErrBrakWiersza
	}
	if err != nil {
		return PaczkaRozszerzenia{}, fmt.Errorf("dane: nieczytelna paczka %q: %w", kod, err)
	}
	return paczka, nil
}

// odczytajKolekcjeRozszerzen składa kolekcję z jednego wiersza wyniku; kody
// pozycji dokłada wołający z osobnego zapytania.
func odczytajKolekcjeRozszerzen(wiersz skaner) (KolekcjaRozszerzen, error) {
	var kolekcja KolekcjaRozszerzen
	var opis, barwa sql.NullString
	err := wiersz.Scan(&kolekcja.ID, &kolekcja.Kod, &kolekcja.Nazwa, &opis, &barwa,
		&kolekcja.Zaktualizowano)
	if err != nil {
		return KolekcjaRozszerzen{}, err
	}
	kolekcja.Opis = tekstZKolumny(opis)
	kolekcja.OznaczenieBarwne = tekstZKolumny(barwa)
	return kolekcja, nil
}
