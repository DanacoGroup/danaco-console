// Odpowiedzialność pliku: dwa obszary modułu Studio, komentarze redakcyjne z adnotacjami przy różnicach oraz zmiany śledzone dokumentu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// KomentarzStudia to wiersz tabeli komentarz_studio, komentarz redakcyjny albo adnotacja przy fragmencie różnicy dokumentu.
type KomentarzStudia struct {
	ID                  int64
	Kod                 string
	DokumentKod         string
	Rodzaj              string
	WersjaKod           *string
	WatekNadrzednyKod   *string
	Autor               string
	ZakresOd            *int64
	ZakresDo            *int64
	FragmentNumer       *int64
	WersjaOdniesieniaID *string
	WersjaPorownywanaID *string
	PropozycjaID        *string
	Tresc               string
	Rozwiazany          bool
	Utworzono           string
}

// ZmianaSledzona to wiersz tabeli zmiana_sledzona_studio, niosący rodzaj, zakres i decyzję o tej zmianie.
type ZmianaSledzona struct {
	ID          int64
	Kod         string
	DokumentKod string
	Rodzaj      string
	Autor       string
	ZakresOd    int64
	ZakresDo    int64
	TrescPrzed  *string
	TrescPo     *string
	Decyzja     string
	Utworzono   string
}

const (
	kolumnyKomentarzaStudia = `k.id, k.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                           k.rodzaj, k.wersja_id, k.watek_nadrzedny_id, k.autor,
	                           k.zakres_od, k.zakres_do, k.fragment_numer,
	                           k.wersja_odniesienia_id, k.wersja_porownywana_id, k.propozycja_id,
	                           k.tresc, k.rozwiazany, k.utworzono`

	zapiszKomentarzStudia = `INSERT INTO komentarz_studio
	                         (identyfikator_zewnetrzny, dokument_id, rodzaj, wersja_id,
	                          watek_nadrzedny_id, autor, zakres_od, zakres_do, fragment_numer,
	                          wersja_odniesienia_id, wersja_porownywana_id, propozycja_id,
	                          tresc, rozwiazany)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzKomentarzStudia = `SELECT ` + kolumnyKomentarzaStudia + `
	                          FROM komentarz_studio k
	                          JOIN dokument_studio d ON d.id = k.dokument_id
	                          WHERE k.identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaKomentarzyStudia = `SELECT ` + kolumnyKomentarzaStudia + `
	                         FROM komentarz_studio k
	                         JOIN dokument_studio d ON d.id = k.dokument_id
	                         WHERE k.dokument_id = ? AND k.rodzaj = ?
	                         ORDER BY k.zakres_od, k.fragment_numer, k.id`

	rozstrzygnijKomentarzStudia = `UPDATE komentarz_studio SET rozwiazany = ?
	                               WHERE identyfikator_zewnetrzny = ?
	                                 AND EXISTS (SELECT 1 FROM dokument_studio
	                                             WHERE dokument_studio.id = komentarz_studio.dokument_id
	                                               AND ` + WarunekKonta + `)`

	kolumnyZmianySledzonej = `z.id, z.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                          z.rodzaj, z.autor, z.zakres_od, z.zakres_do,
	                          z.tresc_przed, z.tresc_po, z.decyzja, z.utworzono`

	zapiszZmianeSledzona = `INSERT INTO zmiana_sledzona_studio
	                        (identyfikator_zewnetrzny, dokument_id, rodzaj, autor,
	                         zakres_od, zakres_do, tresc_przed, tresc_po, decyzja)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	listaZmianSledzonych = `SELECT ` + kolumnyZmianySledzonej + `
	                        FROM zmiana_sledzona_studio z
	                        JOIN dokument_studio d ON d.id = z.dokument_id
	                        WHERE z.dokument_id = ? ORDER BY z.zakres_od, z.id`

	rozstrzygnijZmianeSledzona = `UPDATE zmiana_sledzona_studio SET decyzja = ?
	                              WHERE identyfikator_zewnetrzny = ? AND decyzja = 'oczekuje'
	                                AND EXISTS (SELECT 1 FROM dokument_studio
	                                            WHERE dokument_studio.id = zmiana_sledzona_studio.dokument_id
	                                              AND ` + WarunekKonta + `)`

	ustawSledzenieDokumentu = `UPDATE dokument_studio SET sledzenie_zmian = ?
	                           WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	czytajSledzenieDokumentu = `SELECT sledzenie_zmian FROM dokument_studio
	                            WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

// ZapiszKomentarz zakłada komentarz redakcyjny albo adnotację różnicy i zwraca jego pełny stan po zapisie.
func (r *repozytoriumStudia) ZapiszKomentarz(ctx context.Context,
	dokumentID int64, komentarz KomentarzStudia) (KomentarzStudia, error) {

	if komentarz.Kod == "" {
		return KomentarzStudia{}, fmt.Errorf("dane: komentarz studio bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKomentarzStudia)
	if err != nil {
		return KomentarzStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, komentarz.Kod, dokumentID, komentarz.Rodzaj,
		tekstDoKolumny(komentarz.WersjaKod), tekstDoKolumny(komentarz.WatekNadrzednyKod),
		komentarz.Autor, liczbaDoKolumny(komentarz.ZakresOd), liczbaDoKolumny(komentarz.ZakresDo),
		liczbaDoKolumny(komentarz.FragmentNumer), tekstDoKolumny(komentarz.WersjaOdniesieniaID),
		tekstDoKolumny(komentarz.WersjaPorownywanaID), tekstDoKolumny(komentarz.PropozycjaID),
		komentarz.Tresc, liczbaLogiczna(komentarz.Rozwiazany))
	if err != nil {
		return KomentarzStudia{}, fmt.Errorf("dane: nie można zapisać komentarza studio %q: %w",
			komentarz.Kod, err)
	}
	return r.Komentarz(ctx, komentarz.Kod)
}

// Komentarz zwraca jeden komentarz o wskazanym kodzie zewnętrznym, wraz z jego pełną zapisaną treścią.
func (r *repozytoriumStudia) Komentarz(ctx context.Context, kod string) (KomentarzStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKomentarzStudia)
	if err != nil {
		return KomentarzStudia{}, err
	}
	komentarz, err := odczytajKomentarzStudia(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return KomentarzStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return KomentarzStudia{}, fmt.Errorf("dane: nieczytelny wiersz komentarza studio %q: %w", kod, err)
	}
	return komentarz, nil
}

// Komentarze zwraca komentarze albo adnotacje dokumentu, w kolejności ich położenia w treści dokumentu.
func (r *repozytoriumStudia) Komentarze(ctx context.Context,
	dokumentID int64, rodzaj string) ([]KomentarzStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaKomentarzyStudia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID, rodzaj)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać komentarzy studio: %w", err)
	}
	defer wiersze.Close()

	lista := []KomentarzStudia{}
	for wiersze.Next() {
		komentarz, err := odczytajKomentarzStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz komentarza studio: %w", err)
		}
		lista = append(lista, komentarz)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt komentarzy studio: %w", err)
	}
	return lista, nil
}

// RozstrzygnijKomentarz oznacza wątek komentarza jako rozwiązany albo cofa wcześniejsze jego oznaczenie.
func (r *repozytoriumStudia) RozstrzygnijKomentarz(ctx context.Context,
	kod string, rozwiazany bool) (KomentarzStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, rozstrzygnijKomentarzStudia)
	if err != nil {
		return KomentarzStudia{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, liczbaLogiczna(rozwiazany), kod, KontoOperatora(ctx))
	if err != nil {
		return KomentarzStudia{}, fmt.Errorf("dane: nie można rozstrzygnąć komentarza studio %q: %w", kod, err)
	}
	// Brak wiersza zmienionego znaczy komentarz nieistniejący, nie rozstrzygnięcie bez skutku.
	if zmienione, err := wynik.RowsAffected(); err == nil && zmienione == 0 {
		return KomentarzStudia{}, ErrBrakWiersza
	}
	return r.Komentarz(ctx, kod)
}

// ZapiszZmianeSledzona rejestruje jedną zmianę śledzoną tego samego dokumentu wraz z jej pełnym zakresem.
func (r *repozytoriumStudia) ZapiszZmianeSledzona(ctx context.Context,
	dokumentID int64, zmiana ZmianaSledzona) (ZmianaSledzona, error) {

	if zmiana.Kod == "" {
		return ZmianaSledzona{}, fmt.Errorf("dane: zmiana śledzona bez identyfikatora")
	}
	if zmiana.Decyzja == "" {
		zmiana.Decyzja = "oczekuje"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZmianeSledzona)
	if err != nil {
		return ZmianaSledzona{}, err
	}
	_, err = polecenie.ExecContext(ctx, zmiana.Kod, dokumentID, zmiana.Rodzaj, zmiana.Autor,
		zmiana.ZakresOd, zmiana.ZakresDo, tekstDoKolumny(zmiana.TrescPrzed),
		tekstDoKolumny(zmiana.TrescPo), zmiana.Decyzja)
	if err != nil {
		return ZmianaSledzona{}, fmt.Errorf("dane: nie można zapisać zmiany śledzonej %q: %w",
			zmiana.Kod, err)
	}
	return zmiana, nil
}

// ZmianySledzone zwraca wszystkie zmiany zarejestrowane dla wskazanego dokumentu tego całego modułu Studio.
func (r *repozytoriumStudia) ZmianySledzone(ctx context.Context,
	dokumentID int64) ([]ZmianaSledzona, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZmianSledzonych)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zmian śledzonych: %w", err)
	}
	defer wiersze.Close()

	lista := []ZmianaSledzona{}
	for wiersze.Next() {
		var zmiana ZmianaSledzona
		var przed, po sql.NullString
		err := wiersze.Scan(&zmiana.ID, &zmiana.Kod, &zmiana.DokumentKod, &zmiana.Rodzaj,
			&zmiana.Autor, &zmiana.ZakresOd, &zmiana.ZakresDo, &przed, &po,
			&zmiana.Decyzja, &zmiana.Utworzono)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zmiany śledzonej: %w", err)
		}
		zmiana.TrescPrzed = tekstZKolumny(przed)
		zmiana.TrescPo = tekstZKolumny(po)
		lista = append(lista, zmiana)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zmian śledzonych: %w", err)
	}
	return lista, nil
}

// RozstrzygnijZmianeSledzona zapisuje decyzję o zmianie i mówi, czy ta decyzja rzeczywiście zapadła przy tym wywołaniu.
func (r *repozytoriumStudia) RozstrzygnijZmianeSledzona(ctx context.Context,
	kod, decyzja string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, rozstrzygnijZmianeSledzona)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, decyzja, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można rozstrzygnąć zmiany śledzonej %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek rozstrzygnięcia zmiany %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// UstawSledzenie przestawia stan śledzenia zmian dla wskazanego dokumentu w tym module Studio operacyjnie.
func (r *repozytoriumStudia) UstawSledzenie(ctx context.Context, kodDokumentu string, czynne bool) error {
	polecenie, err := r.zapytania.przygotuj(ctx, ustawSledzenieDokumentu)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, liczbaLogiczna(czynne), kodDokumentu, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można przestawić śledzenia dokumentu %q: %w", kodDokumentu, err)
	}
	if zmienione, err := wynik.RowsAffected(); err == nil && zmienione == 0 {
		return ErrBrakWiersza
	}
	return nil
}

// Sledzenie zwraca bieżący stan śledzenia zmian dla wskazanego dokumentu w tym module Studio operacyjnie.
func (r *repozytoriumStudia) Sledzenie(ctx context.Context, kodDokumentu string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, czytajSledzenieDokumentu)
	if err != nil {
		return false, err
	}
	var czynne int
	err = polecenie.QueryRowContext(ctx, kodDokumentu, KontoOperatora(ctx)).Scan(&czynne)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrBrakWiersza
	}
	if err != nil {
		return false, fmt.Errorf("dane: nieczytelny stan śledzenia dokumentu %q: %w", kodDokumentu, err)
	}
	return czynne == 1, nil
}

// odczytajKomentarzStudia składa strukturę komentarza z jednego wiersza wyniku zapytania, kolumna po kolumnie.
func odczytajKomentarzStudia(wiersz skaner) (KomentarzStudia, error) {
	var komentarz KomentarzStudia
	var wersja, watek, odniesienie, porownywana, propozycja sql.NullString
	var od, do, fragment sql.NullInt64
	var rozwiazany int
	err := wiersz.Scan(&komentarz.ID, &komentarz.Kod, &komentarz.DokumentKod, &komentarz.Rodzaj,
		&wersja, &watek, &komentarz.Autor, &od, &do, &fragment,
		&odniesienie, &porownywana, &propozycja, &komentarz.Tresc, &rozwiazany, &komentarz.Utworzono)
	if err != nil {
		return KomentarzStudia{}, err
	}
	komentarz.WersjaKod = tekstZKolumny(wersja)
	komentarz.WatekNadrzednyKod = tekstZKolumny(watek)
	komentarz.ZakresOd = liczbaZKolumny(od)
	komentarz.ZakresDo = liczbaZKolumny(do)
	komentarz.FragmentNumer = liczbaZKolumny(fragment)
	komentarz.WersjaOdniesieniaID = tekstZKolumny(odniesienie)
	komentarz.WersjaPorownywanaID = tekstZKolumny(porownywana)
	komentarz.PropozycjaID = tekstZKolumny(propozycja)
	komentarz.Rozwiazany = rozwiazany == 1
	return komentarz, nil
}
