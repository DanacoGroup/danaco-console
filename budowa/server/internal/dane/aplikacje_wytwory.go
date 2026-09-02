// Plik definiuje obszar Apps: wytwory pracy modułu — serwer podglądu, motyw
// produktu, dziennik usług i wdrożeń, artefakty budowania oraz pakiety
// rozszerzenia.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

type PodgladApp struct {
	Okno       string
	Warstwa    string
	Adres      string
	Stan       string
	Rozpoczeto int64
	Zatrzymano *int64
}

type MotywApp struct {
	Okno           string
	Tresc          string
	Zaktualizowano string
}

type WierszDziennikaApp struct {
	ID           int64
	Okno         string
	WdrozenieKod *string
	KomponentKod *string
	Chwila       int64
	Tresc        string
}

type ArtefaktApp struct {
	ID            int64
	Kod           string
	Okno          string
	WdrozenieKod  *string
	Rodzaj        string
	Sciezka       string
	Rozmiar       *int64
	SumaKontrolna *string
	Utworzono     string
}

type PakietApp struct {
	ID                int64
	Kod               string
	Okno              string
	Manifest          *string
	ArtefaktOdwolanie *string
	Format            string
	Sciezka           *string
	Rozmiar           *int64
	Podpis            *string
	RozszerzenieKod   *string
	Utworzono         string
	Zaktualizowano    string
}

const (
	zapiszPodgladApp = `INSERT INTO podglad_apps (okno, warstwa, adres, stan, rozpoczeto, zatrzymano, konto_id)
	                    VALUES (?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                    ON CONFLICT(okno) DO UPDATE SET
	                        warstwa = excluded.warstwa,
	                        adres = excluded.adres,
	                        stan = excluded.stan,
	                        rozpoczeto = excluded.rozpoczeto,
	                        zatrzymano = excluded.zatrzymano
	                    WHERE ` + WarunekKonta

	pobierzPodgladApp = `SELECT okno, warstwa, adres, stan, rozpoczeto, zatrzymano
	                     FROM podglad_apps WHERE okno = ? AND ` + WarunekKonta

	zapiszMotywApp = `INSERT INTO motyw_apps (okno, tresc, konto_id) VALUES (?, ?, ` + WskazanieKonta + `)
	                  ON CONFLICT(okno) DO UPDATE SET
	                      tresc = excluded.tresc,
	                      zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                  WHERE ` + WarunekKonta

	pobierzMotywApp = `SELECT okno, tresc, zaktualizowano FROM motyw_apps
	                   WHERE okno = ? AND ` + WarunekKonta

	wstawWierszDziennikaApp = `INSERT INTO wiersz_dziennika_apps
	                           (okno, wdrozenie_kod, komponent_kod, chwila, tresc, konto_id)
	                           VALUES (?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	// Warunek komponentu przepuszcza wszystko przy braku wskazania: jedno zapytanie zamiast dwóch bliźniaczych.
	listaDziennikaUslugiApp = `SELECT id, okno, wdrozenie_kod, komponent_kod, chwila, tresc
	                           FROM wiersz_dziennika_apps
	                           WHERE okno = ? AND chwila >= ?
	                             AND (? = '' OR komponent_kod = ?)
	                             AND ` + WarunekKonta + `
	                           ORDER BY chwila, id LIMIT ?`

	listaDziennikaWdrozeniaApp = `SELECT id, okno, wdrozenie_kod, komponent_kod, chwila, tresc
	                              FROM wiersz_dziennika_apps
	                              WHERE wdrozenie_kod = ? AND chwila >= ?
	                                AND ` + WarunekKonta + `
	                              ORDER BY chwila, id LIMIT ?`

	policzDziennikUslugiApp = `SELECT COUNT(*) FROM wiersz_dziennika_apps
	                           WHERE okno = ? AND chwila >= ? AND (? = '' OR komponent_kod = ?)
	                             AND ` + WarunekKonta

	policzDziennikWdrozeniaApp = `SELECT COUNT(*) FROM wiersz_dziennika_apps
	                              WHERE wdrozenie_kod = ? AND chwila >= ? AND ` + WarunekKonta

	kolumnyArtefaktuApp = `id, identyfikator_zewnetrzny, okno, wdrozenie_kod, rodzaj,
	                       sciezka, rozmiar, suma_kontrolna, utworzono`

	wstawArtefaktApp = `INSERT INTO artefakt_apps
	                    (identyfikator_zewnetrzny, okno, wdrozenie_kod, rodzaj, sciezka,
	                     rozmiar, suma_kontrolna, konto_id)
	                    VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	listaArtefaktowApp = `SELECT ` + kolumnyArtefaktuApp + ` FROM artefakt_apps
	                      WHERE okno = ? AND (? = '' OR wdrozenie_kod = ?)
	                        AND ` + WarunekKonta + `
	                      ORDER BY id DESC`

	// Artefakt ostatniego wdrożenia udanego — domyślne wejście
	// `apps.package.build`, gdy żądanie nie wskazuje artefaktu.
	ostatniArtefaktUdanegoApp = `SELECT ` + kolumnyArtefaktuApp + ` FROM artefakt_apps
	                             WHERE okno = ? AND wdrozenie_kod IN (
	                                 SELECT kod FROM wdrozenie_apps
	                                 WHERE okno = ? AND stan = 'succeeded')
	                               AND ` + WarunekKonta + `
	                             ORDER BY id DESC LIMIT 1`

	pobierzArtefaktApp = `SELECT ` + kolumnyArtefaktuApp + ` FROM artefakt_apps
	                      WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	kolumnyPakietuApp = `id, identyfikator_zewnetrzny, okno, manifest, artefakt_odwolanie,
	                     format, sciezka, rozmiar, podpis, rozszerzenie_kod,
	                     utworzono, zaktualizowano`

	zapiszPakietApp = `INSERT INTO pakiet_apps
	                   (identyfikator_zewnetrzny, okno, manifest, artefakt_odwolanie, format,
	                    sciezka, rozmiar, podpis, rozszerzenie_kod, konto_id)
	                   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                   ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                       manifest = IFNULL(excluded.manifest, pakiet_apps.manifest),
	                       artefakt_odwolanie = IFNULL(excluded.artefakt_odwolanie,
	                                                   pakiet_apps.artefakt_odwolanie),
	                       format = excluded.format,
	                       sciezka = IFNULL(excluded.sciezka, pakiet_apps.sciezka),
	                       rozmiar = IFNULL(excluded.rozmiar, pakiet_apps.rozmiar),
	                       podpis = IFNULL(excluded.podpis, pakiet_apps.podpis),
	                       rozszerzenie_kod = IFNULL(excluded.rozszerzenie_kod,
	                                                 pakiet_apps.rozszerzenie_kod),
	                       zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                   WHERE ` + WarunekKonta

	pobierzPakietApp = `SELECT ` + kolumnyPakietuApp + ` FROM pakiet_apps
	                    WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaPakietowApp = `SELECT ` + kolumnyPakietuApp + ` FROM pakiet_apps
	                    WHERE okno = ? AND ` + WarunekKonta + ` ORDER BY id`
)

func (r *repozytoriumAplikacji) ZapiszPodgladApp(ctx context.Context, podglad PodgladApp) error {
	if podglad.Okno == "" {
		return fmt.Errorf("dane: podgląd aplikacji bez okna")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPodgladApp)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, podglad.Okno, podglad.Warstwa, podglad.Adres,
		podglad.Stan, podglad.Rozpoczeto, liczbaDoKolumny(podglad.Zatrzymano),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać podglądu okna %q: %w", podglad.Okno, err)
	}
	return sprawdzTrafienieZapisu(wynik, "podgląd okna", podglad.Okno)
}

func (r *repozytoriumAplikacji) PodgladApp(ctx context.Context, okno string) (PodgladApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPodgladApp)
	if err != nil {
		return PodgladApp{}, err
	}
	var podglad PodgladApp
	var zatrzymano sql.NullInt64
	err = polecenie.QueryRowContext(ctx, okno, KontoOperatora(ctx)).Scan(&podglad.Okno, &podglad.Warstwa,
		&podglad.Adres, &podglad.Stan, &podglad.Rozpoczeto, &zatrzymano)
	if err == sql.ErrNoRows {
		return PodgladApp{}, ErrBrakWiersza
	}
	if err != nil {
		return PodgladApp{}, fmt.Errorf("dane: nieczytelny podgląd okna %q: %w", okno, err)
	}
	podglad.Zatrzymano = liczbaZKolumny(zatrzymano)
	return podglad, nil
}

func (r *repozytoriumAplikacji) ZapiszMotywApp(ctx context.Context, motyw MotywApp) (MotywApp, error) {
	if motyw.Okno == "" {
		return MotywApp{}, fmt.Errorf("dane: motyw aplikacji bez okna")
	}
	if motyw.Tresc == "" {
		motyw.Tresc = "{}"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszMotywApp)
	if err != nil {
		return MotywApp{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, motyw.Okno, motyw.Tresc, KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return MotywApp{}, fmt.Errorf("dane: nie można zapisać motywu okna %q: %w", motyw.Okno, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "motyw okna", motyw.Okno); err != nil {
		return MotywApp{}, err
	}
	return r.MotywApp(ctx, motyw.Okno)
}

func (r *repozytoriumAplikacji) MotywApp(ctx context.Context, okno string) (MotywApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzMotywApp)
	if err != nil {
		return MotywApp{}, err
	}
	var motyw MotywApp
	err = polecenie.QueryRowContext(ctx, okno, KontoOperatora(ctx)).Scan(&motyw.Okno, &motyw.Tresc, &motyw.Zaktualizowano)
	if err == sql.ErrNoRows {
		return MotywApp{}, ErrBrakWiersza
	}
	if err != nil {
		return MotywApp{}, fmt.Errorf("dane: nieczytelny motyw okna %q: %w", okno, err)
	}
	return motyw, nil
}

func (r *repozytoriumAplikacji) DopiszWierszDziennikaApp(ctx context.Context, wiersz WierszDziennikaApp) error {
	if wiersz.Okno == "" {
		return fmt.Errorf("dane: wiersz dziennika aplikacji bez okna")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawWierszDziennikaApp)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wiersz.Okno, tekstDoKolumny(wiersz.WdrozenieKod),
		tekstDoKolumny(wiersz.KomponentKod), wiersz.Chwila, wiersz.Tresc, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można dopisać wiersza dziennika okna %q: %w", wiersz.Okno, err)
	}
	return nil
}

// DziennikUslugiApp zwraca stronę dziennika usług okna oraz liczbę wszystkich
// wierszy spełniających te same warunki — `total` ma mówić o dzienniku, nie
// o długości oddanej strony.
func (r *repozytoriumAplikacji) DziennikUslugiApp(ctx context.Context, okno, komponent string,
	od int64, granica int) ([]WierszDziennikaApp, int, error) {

	lista, err := r.wierszeDziennikaApp(ctx, listaDziennikaUslugiApp,
		okno, od, komponent, komponent, KontoOperatora(ctx), granica)
	if err != nil {
		return nil, 0, err
	}
	razem, err := r.policzDziennikApp(ctx, policzDziennikUslugiApp, okno, od, komponent, komponent,
		KontoOperatora(ctx))
	if err != nil {
		return nil, 0, err
	}
	return lista, razem, nil
}

func (r *repozytoriumAplikacji) DziennikWdrozeniaApp(ctx context.Context, wdrozenie string,
	od int64, granica int) ([]WierszDziennikaApp, int, error) {

	lista, err := r.wierszeDziennikaApp(ctx, listaDziennikaWdrozeniaApp, wdrozenie, od,
		KontoOperatora(ctx), granica)
	if err != nil {
		return nil, 0, err
	}
	razem, err := r.policzDziennikApp(ctx, policzDziennikWdrozeniaApp, wdrozenie, od, KontoOperatora(ctx))
	if err != nil {
		return nil, 0, err
	}
	return lista, razem, nil
}

func (r *repozytoriumAplikacji) wierszeDziennikaApp(ctx context.Context, zapytanie string,
	argumenty ...any) ([]WierszDziennikaApp, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dziennika aplikacji: %w", err)
	}
	defer wiersze.Close()

	lista := []WierszDziennikaApp{}
	for wiersze.Next() {
		var wpis WierszDziennikaApp
		var wdrozenie, komponent sql.NullString
		if err := wiersze.Scan(&wpis.ID, &wpis.Okno, &wdrozenie, &komponent,
			&wpis.Chwila, &wpis.Tresc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz dziennika aplikacji: %w", err)
		}
		wpis.WdrozenieKod = tekstZKolumny(wdrozenie)
		wpis.KomponentKod = tekstZKolumny(komponent)
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt dziennika aplikacji: %w", err)
	}
	return lista, nil
}

func (r *repozytoriumAplikacji) policzDziennikApp(ctx context.Context, zapytanie string,
	argumenty ...any) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return 0, err
	}
	var razem int
	if err := polecenie.QueryRowContext(ctx, argumenty...).Scan(&razem); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć wierszy dziennika aplikacji: %w", err)
	}
	return razem, nil
}

// ZalozArtefaktApp zapisuje artefakt budowania. Artefakt powstaje raz i się nie
// zmienia — bajty leżą już w magazynie, a ich suma kontrolna jest w wierszu.
func (r *repozytoriumAplikacji) ZalozArtefaktApp(ctx context.Context, artefakt ArtefaktApp) (ArtefaktApp, error) {
	if artefakt.Kod == "" || artefakt.Okno == "" || artefakt.Sciezka == "" {
		return ArtefaktApp{}, fmt.Errorf("dane: artefakt aplikacji bez identyfikatora, okna albo ścieżki")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawArtefaktApp)
	if err != nil {
		return ArtefaktApp{}, err
	}
	_, err = polecenie.ExecContext(ctx, artefakt.Kod, artefakt.Okno,
		tekstDoKolumny(artefakt.WdrozenieKod), artefakt.Rodzaj, artefakt.Sciezka,
		liczbaDoKolumny(artefakt.Rozmiar), tekstDoKolumny(artefakt.SumaKontrolna),
		KontoOperatora(ctx))
	if err != nil {
		return ArtefaktApp{}, fmt.Errorf("dane: nie można zapisać artefaktu %q: %w", artefakt.Kod, err)
	}
	return r.ArtefaktApp(ctx, artefakt.Kod)
}

func (r *repozytoriumAplikacji) ArtefaktApp(ctx context.Context, kod string) (ArtefaktApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzArtefaktApp)
	if err != nil {
		return ArtefaktApp{}, err
	}
	artefakt, err := odczytajArtefaktApp(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if err == sql.ErrNoRows {
		return ArtefaktApp{}, ErrBrakWiersza
	}
	if err != nil {
		return ArtefaktApp{}, fmt.Errorf("dane: nieczytelny artefakt %q: %w", kod, err)
	}
	return artefakt, nil
}

func (r *repozytoriumAplikacji) ArtefaktyApp(ctx context.Context, okno, wdrozenie string) ([]ArtefaktApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaArtefaktowApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, wdrozenie, wdrozenie, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać artefaktów okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []ArtefaktApp{}
	for wiersze.Next() {
		artefakt, err := odczytajArtefaktApp(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz artefaktu okna %q: %w", okno, err)
		}
		lista = append(lista, artefakt)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt artefaktów okna %q: %w", okno, err)
	}
	return lista, nil
}

func (r *repozytoriumAplikacji) OstatniArtefaktUdanegoWdrozeniaApp(ctx context.Context,
	okno string) (ArtefaktApp, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, ostatniArtefaktUdanegoApp)
	if err != nil {
		return ArtefaktApp{}, err
	}
	artefakt, err := odczytajArtefaktApp(polecenie.QueryRowContext(ctx, okno, okno,
		KontoOperatora(ctx)))
	if err == sql.ErrNoRows {
		return ArtefaktApp{}, ErrBrakWiersza
	}
	if err != nil {
		return ArtefaktApp{}, fmt.Errorf("dane: nieczytelny artefakt ostatniego wdrożenia okna %q: %w",
			okno, err)
	}
	return artefakt, nil
}

// ZapiszPakietApp zapisuje pakiet rozszerzenia. Pola podane jako brak NIE
// kasują wartości zastanych: rodzina `apps.package.*` dokłada do pakietu po
// kolei — manifest, podpis, kod opublikowanej pozycji — a każde kolejne
// wywołanie zna tylko swoją część.
func (r *repozytoriumAplikacji) ZapiszPakietApp(ctx context.Context, pakiet PakietApp) (PakietApp, error) {
	if pakiet.Kod == "" || pakiet.Okno == "" {
		return PakietApp{}, fmt.Errorf("dane: pakiet aplikacji bez identyfikatora albo bez okna")
	}
	if pakiet.Format == "" {
		pakiet.Format = "zip"
	}
	if pakiet.ArtefaktOdwolanie != nil && *pakiet.ArtefaktOdwolanie != "" {
		if _, err := r.ArtefaktApp(ctx, *pakiet.ArtefaktOdwolanie); err != nil {
			return PakietApp{}, fmt.Errorf("dane: artefakt %q pakietu %q: %w",
				*pakiet.ArtefaktOdwolanie, pakiet.Kod, err)
		}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPakietApp)
	if err != nil {
		return PakietApp{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, pakiet.Kod, pakiet.Okno, tekstDoKolumny(pakiet.Manifest),
		tekstDoKolumny(pakiet.ArtefaktOdwolanie), pakiet.Format, tekstDoKolumny(pakiet.Sciezka),
		liczbaDoKolumny(pakiet.Rozmiar), tekstDoKolumny(pakiet.Podpis),
		tekstDoKolumny(pakiet.RozszerzenieKod), KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return PakietApp{}, fmt.Errorf("dane: nie można zapisać pakietu %q: %w", pakiet.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "pakiet", pakiet.Kod); err != nil {
		return PakietApp{}, err
	}
	return r.PakietApp(ctx, pakiet.Kod)
}

func (r *repozytoriumAplikacji) PakietApp(ctx context.Context, kod string) (PakietApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPakietApp)
	if err != nil {
		return PakietApp{}, err
	}
	pakiet, err := odczytajPakietApp(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if err == sql.ErrNoRows {
		return PakietApp{}, ErrBrakWiersza
	}
	if err != nil {
		return PakietApp{}, fmt.Errorf("dane: nieczytelny pakiet %q: %w", kod, err)
	}
	return pakiet, nil
}

// PakietyApp zwraca pakiety okna w kolejności powstawania — oś czasu projektu
// czyta z nich zdarzenia rodzaju `package`.
func (r *repozytoriumAplikacji) PakietyApp(ctx context.Context, okno string) ([]PakietApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPakietowApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pakietów okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []PakietApp{}
	for wiersze.Next() {
		pakiet, err := odczytajPakietApp(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pakietu okna %q: %w", okno, err)
		}
		lista = append(lista, pakiet)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pakietów okna %q: %w", okno, err)
	}
	return lista, nil
}

func odczytajPakietApp(wiersz skaner) (PakietApp, error) {
	var pakiet PakietApp
	var manifest, artefakt, sciezka, podpis, rozszerzenie sql.NullString
	var rozmiar sql.NullInt64
	err := wiersz.Scan(&pakiet.ID, &pakiet.Kod, &pakiet.Okno,
		&manifest, &artefakt, &pakiet.Format, &sciezka, &rozmiar, &podpis, &rozszerzenie,
		&pakiet.Utworzono, &pakiet.Zaktualizowano)
	if err != nil {
		return PakietApp{}, err
	}
	pakiet.Manifest = tekstZKolumny(manifest)
	pakiet.ArtefaktOdwolanie = tekstZKolumny(artefakt)
	pakiet.Sciezka = tekstZKolumny(sciezka)
	pakiet.Rozmiar = liczbaZKolumny(rozmiar)
	pakiet.Podpis = tekstZKolumny(podpis)
	pakiet.RozszerzenieKod = tekstZKolumny(rozszerzenie)
	return pakiet, nil
}

func odczytajArtefaktApp(wiersz skaner) (ArtefaktApp, error) {
	var artefakt ArtefaktApp
	var wdrozenie, suma sql.NullString
	var rozmiar sql.NullInt64
	err := wiersz.Scan(&artefakt.ID, &artefakt.Kod, &artefakt.Okno, &wdrozenie,
		&artefakt.Rodzaj, &artefakt.Sciezka, &rozmiar, &suma, &artefakt.Utworzono)
	if err != nil {
		return ArtefaktApp{}, err
	}
	artefakt.WdrozenieKod = tekstZKolumny(wdrozenie)
	artefakt.Rozmiar = liczbaZKolumny(rozmiar)
	artefakt.SumaKontrolna = tekstZKolumny(suma)
	return artefakt, nil
}
