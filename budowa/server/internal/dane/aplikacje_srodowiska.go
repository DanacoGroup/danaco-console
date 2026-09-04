// Plik obsługuje obszar Apps: środowiska wdrożeniowe, ich zmienne, nastawy skalowania oraz
// wyniki sprawdzeń kondycji. Uzasadnienie wyłączności wartości i sekretu zmiennej niesie
// rozdział aplikacje_srodowiska.go dokumentacji architektury.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// SrodowiskoApp odwzorowuje wiersz tabeli srodowisko_apps: jedno środowisko wdrożeniowe
// okna wraz z domeną i wpisami DNS.
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

// ZmiennaSrodowiskaApp odwzorowuje wiersz tabeli zmienna_srodowiska_apps: jedną zmienną
// środowiskową z wartością jawną albo odwołaniem do sekretu.
type ZmiennaSrodowiskaApp struct {
	ID               int64
	Okno             string
	Srodowisko       string
	Nazwa            string
	Wartosc          *string
	OdwolanieSekretu *string
	Zaktualizowano   string
}

// SkalowanieApp odwzorowuje wiersz tabeli skalowanie_apps: nastawę liczby instancji
// i reguł skalowania jednego środowiska.
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
	                       (identyfikator_zewnetrzny, okno, kod, nazwa, kolejnosc, domena,
	                        wpisy_dns, konto_id)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                       ON CONFLICT(okno, kod) DO UPDATE SET
	                           nazwa = excluded.nazwa,
	                           kolejnosc = excluded.kolejnosc,
	                           domena = IFNULL(excluded.domena, srodowisko_apps.domena),
	                           wpisy_dns = IFNULL(excluded.wpisy_dns, srodowisko_apps.wpisy_dns),
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                       WHERE ` + WarunekKonta

	listaSrodowiskApp = `SELECT ` + kolumnySrodowiskaApp + ` FROM srodowisko_apps
	                     WHERE okno = ? AND ` + WarunekKonta + `
	                     ORDER BY kolejnosc, id`

	pobierzSrodowiskoApp = `SELECT ` + kolumnySrodowiskaApp + ` FROM srodowisko_apps
	                        WHERE okno = ? AND kod = ? AND ` + WarunekKonta

	kolumnyZmiennejSrodowiskaApp = `id, okno, srodowisko, nazwa, wartosc, odwolanie_sekretu,
	                                zaktualizowano`

	// Trójka okno-środowisko-nazwa jest UNIQUE w całej tabeli, więc warunek przy
	// DO UPDATE zatrzymuje nadpisanie wiersza należącego do innego konta.
	zapiszZmiennaSrodowiskaApp = `INSERT INTO zmienna_srodowiska_apps
	                              (okno, srodowisko, nazwa, wartosc, odwolanie_sekretu, konto_id)
	                              VALUES (?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                              ON CONFLICT(okno, srodowisko, nazwa) DO UPDATE SET
	                                  wartosc = excluded.wartosc,
	                                  odwolanie_sekretu = excluded.odwolanie_sekretu,
	                                  zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                              WHERE ` + WarunekKonta

	pobierzZmiennaSrodowiskaApp = `SELECT ` + kolumnyZmiennejSrodowiskaApp + `
	                               FROM zmienna_srodowiska_apps
	                               WHERE okno = ? AND srodowisko = ? AND nazwa = ?
	                                 AND ` + WarunekKonta

	listaZmiennychSrodowiskaApp = `SELECT ` + kolumnyZmiennejSrodowiskaApp + `
	                               FROM zmienna_srodowiska_apps
	                               WHERE okno = ? AND srodowisko = ? AND ` + WarunekKonta + `
	                               ORDER BY nazwa`

	zapiszSkalowanieApp = `INSERT INTO skalowanie_apps
	                       (okno, srodowisko, instancje, min_instancji, maks_instancji, reguly,
	                        konto_id)
	                       VALUES (?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                       ON CONFLICT(okno, srodowisko) DO UPDATE SET
	                           instancje = excluded.instancje,
	                           min_instancji = excluded.min_instancji,
	                           maks_instancji = excluded.maks_instancji,
	                           reguly = excluded.reguly,
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                       WHERE ` + WarunekKonta

	pobierzSkalowanieApp = `SELECT okno, srodowisko, instancje, min_instancji, maks_instancji,
	                               reguly, zaktualizowano
	                        FROM skalowanie_apps
	                        WHERE okno = ? AND srodowisko = ? AND ` + WarunekKonta

	wstawKondycjeApp = `INSERT INTO kondycja_wdrozenia_apps
	                    (okno, srodowisko, dostepna, szczegol, sprawdzono, konto_id)
	                    VALUES (?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	listaKondycjiApp = `SELECT id, okno, srodowisko, dostepna, szczegol, sprawdzono
	                    FROM kondycja_wdrozenia_apps
	                    WHERE okno = ? AND srodowisko = ? AND ` + WarunekKonta + `
	                    ORDER BY sprawdzono DESC, id DESC LIMIT ?`
)

// ZapiszSrodowiskoApp zakłada albo zmienia środowisko okna. Domena i wpisy DNS podane
// jako brak nie kasują wartości zastanej, ponieważ domenę nadaje osobna komenda i nie
// traci jej kolejny wykaz środowisk.
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
		tekstDoKolumny(srodowisko.WpisyDNS), KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return SrodowiskoApp{}, fmt.Errorf("dane: nie można zapisać środowiska %q okna %q: %w",
			srodowisko.KodSrodowiska, srodowisko.Okno, err)
	}
	return r.SrodowiskoApp(ctx, srodowisko.Okno, srodowisko.KodSrodowiska)
}

// SrodowiskoApp zwraca jedno środowisko okna wskazane kodem albo błąd ErrBrakWiersza,
// gdy nie istnieje.
func (r *repozytoriumAplikacji) SrodowiskoApp(ctx context.Context, okno, kod string) (SrodowiskoApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSrodowiskoApp)
	if err != nil {
		return SrodowiskoApp{}, err
	}
	srodowisko, err := odczytajSrodowiskoApp(polecenie.QueryRowContext(ctx, okno, kod, KontoOperatora(ctx)))
	if err == sql.ErrNoRows {
		return SrodowiskoApp{}, ErrBrakWiersza
	}
	if err != nil {
		return SrodowiskoApp{}, fmt.Errorf("dane: nieczytelne środowisko %q okna %q: %w", kod, okno, err)
	}
	return srodowisko, nil
}

// SrodowiskaApp zwraca wszystkie środowiska okna uporządkowane w kolejności wdrażania,
// zgodnie z polem kolejnosc.
func (r *repozytoriumAplikacji) SrodowiskaApp(ctx context.Context, okno string) ([]SrodowiskoApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSrodowiskApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
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

// ZapiszZmiennaSrodowiskaApp zapisuje zmienną środowiskową produktu, odrzucając zapis
// niosący naraz wartość i odwołanie do sekretu.
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
	wynik, err := polecenie.ExecContext(ctx, zmienna.Okno, zmienna.Srodowisko, zmienna.Nazwa,
		tekstDoKolumny(zmienna.Wartosc), tekstDoKolumny(zmienna.OdwolanieSekretu),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return ZmiennaSrodowiskaApp{}, fmt.Errorf("dane: nie można zapisać zmiennej %q: %w",
			zmienna.Nazwa, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "zmienna środowiskowa", zmienna.Nazwa); err != nil {
		return ZmiennaSrodowiskaApp{}, err
	}

	polecenieOdczytu, err := r.zapytania.przygotuj(ctx, pobierzZmiennaSrodowiskaApp)
	if err != nil {
		return ZmiennaSrodowiskaApp{}, err
	}
	zapisana, err := odczytajZmiennaSrodowiskaApp(
		polecenieOdczytu.QueryRowContext(ctx, zmienna.Okno, zmienna.Srodowisko, zmienna.Nazwa,
			KontoOperatora(ctx)))
	if err != nil {
		return ZmiennaSrodowiskaApp{}, fmt.Errorf("dane: nieczytelna zmienna %q po zapisie: %w",
			zmienna.Nazwa, err)
	}
	return zapisana, nil
}

// ZmienneSrodowiskaApp zwraca wszystkie zmienne jednego środowiska okna, uporządkowane
// według nazwy zmiennej.
func (r *repozytoriumAplikacji) ZmienneSrodowiskaApp(ctx context.Context,
	okno, srodowisko string) ([]ZmiennaSrodowiskaApp, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZmiennychSrodowiskaApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, srodowisko, KontoOperatora(ctx))
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

// ZapiszSkalowanieApp zapisuje nastawę skalowania jednego środowiska: liczbę instancji,
// granice i reguły.
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
		liczbaDoKolumny(skalowanie.MaksInstancji), tekstDoKolumny(skalowanie.Reguly),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return SkalowanieApp{}, fmt.Errorf("dane: nie można zapisać skalowania %q okna %q: %w",
			skalowanie.Srodowisko, skalowanie.Okno, err)
	}
	return r.SkalowanieApp(ctx, skalowanie.Okno, skalowanie.Srodowisko)
}

// SkalowanieApp zwraca nastawę skalowania jednego środowiska albo błąd ErrBrakWiersza,
// gdy nastawa nie istnieje.
func (r *repozytoriumAplikacji) SkalowanieApp(ctx context.Context,
	okno, srodowisko string) (SkalowanieApp, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSkalowanieApp)
	if err != nil {
		return SkalowanieApp{}, err
	}
	var nastawa SkalowanieApp
	var instancje, minimum, maksimum sql.NullInt64
	var reguly sql.NullString
	err = polecenie.QueryRowContext(ctx, okno, srodowisko, KontoOperatora(ctx)).Scan(&nastawa.Okno, &nastawa.Srodowisko,
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

// ZapiszKondycjeApp dopisuje wynik jednego sprawdzenia kondycji wdrożenia, nie
// nadpisując wyników wcześniejszych.
func (r *repozytoriumAplikacji) ZapiszKondycjeApp(ctx context.Context, kondycja KondycjaWdrozeniaApp) error {
	if kondycja.Okno == "" || kondycja.Srodowisko == "" {
		return fmt.Errorf("dane: wynik kondycji bez okna albo bez środowiska")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKondycjeApp)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, kondycja.Okno, kondycja.Srodowisko,
		liczbaLogiczna(kondycja.Dostepna), tekstDoKolumny(kondycja.Szczegol), kondycja.Sprawdzono,
		KontoOperatora(ctx))
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
	wiersze, err := polecenie.QueryContext(ctx, okno, srodowisko, KontoOperatora(ctx), granica)
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

// odczytajSrodowiskoApp składa strukturę SrodowiskoApp z jednego wiersza wyniku
// zapytania, niezależnie od jego źródła.
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

// odczytajZmiennaSrodowiskaApp składa strukturę ZmiennaSrodowiskaApp z jednego wiersza
// wyniku zapytania.
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
