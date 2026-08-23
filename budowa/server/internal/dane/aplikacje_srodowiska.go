// Odpowiedzialność pliku: obszar Apps — Deployment Panel. Środowiska
// wdrożeniowe (`srodowisko_apps`), ich zmienne (`zmienna_srodowiska_apps`),
// nastawy skalowania (`skalowanie_apps`) i wyniki sprawdzeń kondycji
// (`kondycja_wdrozenia_apps`) — `store/migracja_202_apps_srodowiska.sql`.
//
// Zmienna niesie wartość jawną ALBO odwołanie do sekretu; warunek CHECK
// schematu pilnuje tego po raz drugi, a warstwa `dane` nie przepuszcza obu
// naraz, żeby literówka wołającego wracała powodem, a nie treścią SQL-a.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// SrodowiskoApp to wiersz tabeli `srodowisko_apps`.
type SrodowiskoApp struct {
	ID             int64
	Kod            string
	Okno           string
	KodSrodowiska  string
	Nazwa          string
	Kolejnosc      int
	Domena         *string
	WpisyDNS       *string
	Utworzono      string
	Zaktualizowano string
}

// ZmiennaSrodowiskaApp to wiersz tabeli `zmienna_srodowiska_apps`.
type ZmiennaSrodowiskaApp struct {
	ID               int64
	Okno             string
	Srodowisko       string
	Nazwa            string
	Wartosc          *string
	OdwolanieSekretu *string
	Zaktualizowano   string
}

// SkalowanieApp to wiersz tabeli `skalowanie_apps`.
type SkalowanieApp struct {
	Okno           string
	Srodowisko     string
	Instancje      *int64
	MinInstancji   *int64
	MaksInstancji  *int64
	Reguly         *string
	Zaktualizowano string
}

// KondycjaWdrozeniaApp to wiersz tabeli `kondycja_wdrozenia_apps` — jeden
// wynik sprawdzenia, nie stan bieżący (czoło migracji 202).
type KondycjaWdrozeniaApp struct {
	ID         int64
	Okno       string
	Srodowisko string
	Dostepna   bool
	Szczegol   *string
	Sprawdzono int64
}

const (
	kolumnySrodowiskaApp = `id, identyfikator_zewnetrzny, okno, kod, nazwa, kolejnosc,
	                        domena, wpisy_dns, utworzono, zaktualizowano`

	zapiszSrodowiskoApp = `INSERT INTO srodowisko_apps
	                       (identyfikator_zewnetrzny, okno, kod, nazwa, kolejnosc, domena, wpisy_dns)
	                       VALUES (?, ?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(okno, kod) DO UPDATE SET
	                           nazwa = excluded.nazwa,
	                           kolejnosc = excluded.kolejnosc,
	                           domena = IFNULL(excluded.domena, srodowisko_apps.domena),
	                           wpisy_dns = IFNULL(excluded.wpisy_dns, srodowisko_apps.wpisy_dns),
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	listaSrodowiskApp = `SELECT ` + kolumnySrodowiskaApp + ` FROM srodowisko_apps
	                     WHERE okno = ? ORDER BY kolejnosc, id`

	pobierzSrodowiskoApp = `SELECT ` + kolumnySrodowiskaApp + ` FROM srodowisko_apps
	                        WHERE okno = ? AND kod = ?`

	kolumnyZmiennejSrodowiskaApp = `id, okno, srodowisko, nazwa, wartosc, odwolanie_sekretu,
	                                zaktualizowano`

	zapiszZmiennaSrodowiskaApp = `INSERT INTO zmienna_srodowiska_apps
	                              (okno, srodowisko, nazwa, wartosc, odwolanie_sekretu)
	                              VALUES (?, ?, ?, ?, ?)
	                              ON CONFLICT(okno, srodowisko, nazwa) DO UPDATE SET
	                                  wartosc = excluded.wartosc,
	                                  odwolanie_sekretu = excluded.odwolanie_sekretu,
	                                  zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzZmiennaSrodowiskaApp = `SELECT ` + kolumnyZmiennejSrodowiskaApp + `
	                               FROM zmienna_srodowiska_apps
	                               WHERE okno = ? AND srodowisko = ? AND nazwa = ?`

	listaZmiennychSrodowiskaApp = `SELECT ` + kolumnyZmiennejSrodowiskaApp + `
	                               FROM zmienna_srodowiska_apps
	                               WHERE okno = ? AND srodowisko = ? ORDER BY nazwa`

	zapiszSkalowanieApp = `INSERT INTO skalowanie_apps
	                       (okno, srodowisko, instancje, min_instancji, maks_instancji, reguly)
	                       VALUES (?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(okno, srodowisko) DO UPDATE SET
	                           instancje = excluded.instancje,
	                           min_instancji = excluded.min_instancji,
	                           maks_instancji = excluded.maks_instancji,
	                           reguly = excluded.reguly,
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzSkalowanieApp = `SELECT okno, srodowisko, instancje, min_instancji, maks_instancji,
	                               reguly, zaktualizowano
	                        FROM skalowanie_apps WHERE okno = ? AND srodowisko = ?`

	wstawKondycjeApp = `INSERT INTO kondycja_wdrozenia_apps
	                    (okno, srodowisko, dostepna, szczegol, sprawdzono)
	                    VALUES (?, ?, ?, ?, ?)`

	listaKondycjiApp = `SELECT id, okno, srodowisko, dostepna, szczegol, sprawdzono
	                    FROM kondycja_wdrozenia_apps
	                    WHERE okno = ? AND srodowisko = ?
	                    ORDER BY sprawdzono DESC, id DESC LIMIT ?`
)

// ZapiszSrodowiskoApp zakłada albo zmienia środowisko okna. Domena i wpisy DNS
// podane jako brak NIE kasują wartości zastanej — środowisko zakłada się przy
// pierwszym wykazie, a domenę nadaje osobna komenda i nie ma prawa jej stracić
// przy kolejnym wykazie.
func (r *repozytoriumAplikacji) ZapiszSrodowiskoApp(ctx context.Context,
	srodowisko SrodowiskoApp) (SrodowiskoApp, error) {

	if srodowisko.Okno == "" || srodowisko.KodSrodowiska == "" {
		return SrodowiskoApp{}, fmt.Errorf("dane: środowisko aplikacji bez okna albo bez kodu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSrodowiskoApp)
	if err != nil {
		return SrodowiskoApp{}, err
	}
	_, err = polecenie.ExecContext(ctx, srodowisko.Kod, srodowisko.Okno, srodowisko.KodSrodowiska,
		srodowisko.Nazwa, srodowisko.Kolejnosc, tekstDoKolumny(srodowisko.Domena),
		tekstDoKolumny(srodowisko.WpisyDNS))
	if err != nil {
		return SrodowiskoApp{}, fmt.Errorf("dane: nie można zapisać środowiska %q okna %q: %w",
			srodowisko.KodSrodowiska, srodowisko.Okno, err)
	}
	return r.SrodowiskoApp(ctx, srodowisko.Okno, srodowisko.KodSrodowiska)
}

// SrodowiskoApp zwraca jedno środowisko okna.
func (r *repozytoriumAplikacji) SrodowiskoApp(ctx context.Context, okno, kod string) (SrodowiskoApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSrodowiskoApp)
	if err != nil {
		return SrodowiskoApp{}, err
	}
	srodowisko, err := odczytajSrodowiskoApp(polecenie.QueryRowContext(ctx, okno, kod))
	if err == sql.ErrNoRows {
		return SrodowiskoApp{}, ErrBrakWiersza
	}
	if err != nil {
		return SrodowiskoApp{}, fmt.Errorf("dane: nieczytelne środowisko %q okna %q: %w", kod, okno, err)
	}
	return srodowisko, nil
}

// SrodowiskaApp zwraca środowiska okna w kolejności wdrażania.
func (r *repozytoriumAplikacji) SrodowiskaApp(ctx context.Context, okno string) ([]SrodowiskoApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSrodowiskApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać środowisk okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []SrodowiskoApp{}
	for wiersze.Next() {
		srodowisko, err := odczytajSrodowiskoApp(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz środowiska okna %q: %w", okno, err)
		}
		lista = append(lista, srodowisko)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt środowisk okna %q: %w", okno, err)
	}
	return lista, nil
}

// ZapiszZmiennaSrodowiskaApp zapisuje zmienną środowiskową produktu.
func (r *repozytoriumAplikacji) ZapiszZmiennaSrodowiskaApp(ctx context.Context,
	zmienna ZmiennaSrodowiskaApp) (ZmiennaSrodowiskaApp, error) {

	if zmienna.Okno == "" || zmienna.Srodowisko == "" || zmienna.Nazwa == "" {
		return ZmiennaSrodowiskaApp{}, fmt.Errorf("dane: zmienna środowiskowa bez okna, środowiska albo nazwy")
	}
	if zmienna.Wartosc != nil && zmienna.OdwolanieSekretu != nil {
		return ZmiennaSrodowiskaApp{}, fmt.Errorf(
			"dane: zmienna %q niesie naraz wartość jawną i odwołanie do sekretu", zmienna.Nazwa)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZmiennaSrodowiskaApp)
	if err != nil {
		return ZmiennaSrodowiskaApp{}, err
	}
	_, err = polecenie.ExecContext(ctx, zmienna.Okno, zmienna.Srodowisko, zmienna.Nazwa,
		tekstDoKolumny(zmienna.Wartosc), tekstDoKolumny(zmienna.OdwolanieSekretu))
	if err != nil {
		return ZmiennaSrodowiskaApp{}, fmt.Errorf("dane: nie można zapisać zmiennej %q: %w",
			zmienna.Nazwa, err)
	}

	polecenieOdczytu, err := r.zapytania.przygotuj(ctx, pobierzZmiennaSrodowiskaApp)
	if err != nil {
		return ZmiennaSrodowiskaApp{}, err
	}
	zapisana, err := odczytajZmiennaSrodowiskaApp(
		polecenieOdczytu.QueryRowContext(ctx, zmienna.Okno, zmienna.Srodowisko, zmienna.Nazwa))
	if err != nil {
		return ZmiennaSrodowiskaApp{}, fmt.Errorf("dane: nieczytelna zmienna %q po zapisie: %w",
			zmienna.Nazwa, err)
	}
	return zapisana, nil
}

// ZmienneSrodowiskaApp zwraca zmienne jednego środowiska, po nazwie.
func (r *repozytoriumAplikacji) ZmienneSrodowiskaApp(ctx context.Context,
	okno, srodowisko string) ([]ZmiennaSrodowiskaApp, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZmiennychSrodowiskaApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, srodowisko)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zmiennych środowiska %q okna %q: %w",
			srodowisko, okno, err)
	}
	defer wiersze.Close()

	lista := []ZmiennaSrodowiskaApp{}
	for wiersze.Next() {
		zmienna, err := odczytajZmiennaSrodowiskaApp(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zmiennej środowiskowej: %w", err)
		}
		lista = append(lista, zmienna)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zmiennych środowiska %q okna %q: %w",
			srodowisko, okno, err)
	}
	return lista, nil
}

// ZapiszSkalowanieApp zapisuje nastawę skalowania jednego środowiska.
func (r *repozytoriumAplikacji) ZapiszSkalowanieApp(ctx context.Context,
	skalowanie SkalowanieApp) (SkalowanieApp, error) {

	if skalowanie.Okno == "" || skalowanie.Srodowisko == "" {
		return SkalowanieApp{}, fmt.Errorf("dane: nastawa skalowania bez okna albo bez środowiska")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSkalowanieApp)
	if err != nil {
		return SkalowanieApp{}, err
	}
	_, err = polecenie.ExecContext(ctx, skalowanie.Okno, skalowanie.Srodowisko,
		liczbaDoKolumny(skalowanie.Instancje), liczbaDoKolumny(skalowanie.MinInstancji),
		liczbaDoKolumny(skalowanie.MaksInstancji), tekstDoKolumny(skalowanie.Reguly))
	if err != nil {
		return SkalowanieApp{}, fmt.Errorf("dane: nie można zapisać skalowania %q okna %q: %w",
			skalowanie.Srodowisko, skalowanie.Okno, err)
	}
	return r.SkalowanieApp(ctx, skalowanie.Okno, skalowanie.Srodowisko)
}

// SkalowanieApp zwraca nastawę skalowania jednego środowiska.
func (r *repozytoriumAplikacji) SkalowanieApp(ctx context.Context,
	okno, srodowisko string) (SkalowanieApp, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSkalowanieApp)
	if err != nil {
		return SkalowanieApp{}, err
	}
	var nastawa SkalowanieApp
	var instancje, minimum, maksimum sql.NullInt64
	var reguly sql.NullString
	err = polecenie.QueryRowContext(ctx, okno, srodowisko).Scan(&nastawa.Okno, &nastawa.Srodowisko,
		&instancje, &minimum, &maksimum, &reguly, &nastawa.Zaktualizowano)
	if err == sql.ErrNoRows {
		return SkalowanieApp{}, ErrBrakWiersza
	}
	if err != nil {
		return SkalowanieApp{}, fmt.Errorf("dane: nieczytelna nastawa skalowania %q okna %q: %w",
			srodowisko, okno, err)
	}
	nastawa.Instancje = liczbaZKolumny(instancje)
	nastawa.MinInstancji = liczbaZKolumny(minimum)
	nastawa.MaksInstancji = liczbaZKolumny(maksimum)
	nastawa.Reguly = tekstZKolumny(reguly)
	return nastawa, nil
}

// ZapiszKondycjeApp dopisuje wynik jednego sprawdzenia kondycji.
func (r *repozytoriumAplikacji) ZapiszKondycjeApp(ctx context.Context, kondycja KondycjaWdrozeniaApp) error {
	if kondycja.Okno == "" || kondycja.Srodowisko == "" {
		return fmt.Errorf("dane: wynik kondycji bez okna albo bez środowiska")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKondycjeApp)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, kondycja.Okno, kondycja.Srodowisko,
		liczbaLogiczna(kondycja.Dostepna), tekstDoKolumny(kondycja.Szczegol), kondycja.Sprawdzono)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać wyniku kondycji okna %q: %w", kondycja.Okno, err)
	}
	return nil
}

// KondycjeApp zwraca ostatnie wyniki sprawdzeń, od najnowszego. Udział
// dostępności liczy się z nich, a nie z jednego wiersza „stan bieżący".
func (r *repozytoriumAplikacji) KondycjeApp(ctx context.Context, okno, srodowisko string,
	granica int) ([]KondycjaWdrozeniaApp, error) {

	if granica <= 0 {
		granica = 50
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaKondycjiApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, srodowisko, granica)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kondycji %q okna %q: %w", srodowisko, okno, err)
	}
	defer wiersze.Close()

	lista := []KondycjaWdrozeniaApp{}
	for wiersze.Next() {
		var kondycja KondycjaWdrozeniaApp
		var dostepna int
		var szczegol sql.NullString
		err := wiersze.Scan(&kondycja.ID, &kondycja.Okno, &kondycja.Srodowisko, &dostepna,
			&szczegol, &kondycja.Sprawdzono)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kondycji: %w", err)
		}
		kondycja.Dostepna = dostepna == 1
		kondycja.Szczegol = tekstZKolumny(szczegol)
		lista = append(lista, kondycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kondycji %q okna %q: %w", srodowisko, okno, err)
	}
	return lista, nil
}

// odczytajSrodowiskoApp składa środowisko z jednego wiersza wyniku.
func odczytajSrodowiskoApp(wiersz skaner) (SrodowiskoApp, error) {
	var srodowisko SrodowiskoApp
	var domena, wpisy sql.NullString
	err := wiersz.Scan(&srodowisko.ID, &srodowisko.Kod, &srodowisko.Okno, &srodowisko.KodSrodowiska,
		&srodowisko.Nazwa, &srodowisko.Kolejnosc, &domena, &wpisy,
		&srodowisko.Utworzono, &srodowisko.Zaktualizowano)
	if err != nil {
		return SrodowiskoApp{}, err
	}
	srodowisko.Domena = tekstZKolumny(domena)
	srodowisko.WpisyDNS = tekstZKolumny(wpisy)
	return srodowisko, nil
}

// odczytajZmiennaSrodowiskaApp składa zmienną z jednego wiersza wyniku.
func odczytajZmiennaSrodowiskaApp(wiersz skaner) (ZmiennaSrodowiskaApp, error) {
	var zmienna ZmiennaSrodowiskaApp
	var wartosc, sekret sql.NullString
	err := wiersz.Scan(&zmienna.ID, &zmienna.Okno, &zmienna.Srodowisko, &zmienna.Nazwa,
		&wartosc, &sekret, &zmienna.Zaktualizowano)
	if err != nil {
		return ZmiennaSrodowiskaApp{}, err
	}
	zmienna.Wartosc = tekstZKolumny(wartosc)
	zmienna.OdwolanieSekretu = tekstZKolumny(sekret)
	return zmienna, nil
}
