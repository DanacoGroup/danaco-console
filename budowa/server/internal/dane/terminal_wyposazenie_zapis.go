// Odpowiedzialność pliku: zapis wyposażenia modułu Terminal — wpisów książki
// hostów, pozycji biblioteki skryptów wraz z ich wersjami, kluczy SSH, tuneli
// i obserwacji plików.
//
// Usunięcie oddaje PRAWDĘ O SKUTKU (`bool`), a nie samo „nie było błędu”.
// Kontrakt komend `terminal.host.remove`, `terminal.script.remove`
// i `terminal.key.remove` mówi wprost: fałsz znaczy, że wpisu nie było, i nie
// jest błędem. Bez policzenia zmienionych wierszy rdzeń nie miałby czym tego
// rozróżnić i odpowiadałby prawdą zawsze.
package dane

import (
	"context"
	"database/sql"
	"fmt"

	"danacoconsole/shared"
)

const (
	wstawHostaTerminala = `INSERT INTO terminal_host
	                       (kod, nazwa, cel, port, grupa, katalog_roboczy, klucz_kod,
	                        host_posredni_kod, notatka)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(kod) DO UPDATE SET
	                         nazwa             = excluded.nazwa,
	                         cel               = excluded.cel,
	                         port              = excluded.port,
	                         grupa             = excluded.grupa,
	                         katalog_roboczy   = excluded.katalog_roboczy,
	                         klucz_kod         = excluded.klucz_kod,
	                         host_posredni_kod = excluded.host_posredni_kod,
	                         notatka           = excluded.notatka,
	                         zaktualizowano    = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	usunHostaTerminala = `DELETE FROM terminal_host WHERE kod = ?`

	// Odpięcie klucza czyta wpisy przed zmianą, żeby oddać ich kody; sama zmiana
	// idzie osobnym poleceniem w tej samej transakcji.
	hostyPoKluczuTerminala = `SELECT kod FROM terminal_host WHERE klucz_kod = ?`

	odepnijKluczTerminala = `UPDATE terminal_host
	                         SET klucz_kod = NULL,
	                             zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                         WHERE klucz_kod = ?`

	wstawSkryptTerminala = `INSERT INTO terminal_skrypt
	                        (kod, nazwa, rodzaj, powloka, tresc, znaczniki, alias, wersja)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, 1)`

	podmienSkryptTerminala = `UPDATE terminal_skrypt
	                          SET nazwa          = ?,
	                              rodzaj         = ?,
	                              powloka        = ?,
	                              tresc          = ?,
	                              znaczniki      = ?,
	                              alias          = ?,
	                              wersja         = wersja + 1,
	                              zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE kod = ?`

	wersjaSkryptuTerminala = `SELECT wersja FROM terminal_skrypt WHERE kod = ?`

	wstawWersjeSkryptuTerminala = `INSERT INTO terminal_skrypt_wersja (skrypt_kod, wersja, tresc)
	                               VALUES (?, ?, ?)
	                               ON CONFLICT(skrypt_kod, wersja) DO UPDATE SET tresc = excluded.tresc`

	usunSkryptTerminala = `DELETE FROM terminal_skrypt WHERE kod = ?`

	wstawKluczTerminala = `INSERT INTO terminal_klucz
	                       (kod, nazwa, rodzaj, odcisk, klucz_jawny, sciezka, haslo)
	                       VALUES (?, ?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(kod) DO UPDATE SET
	                         nazwa       = excluded.nazwa,
	                         rodzaj      = excluded.rodzaj,
	                         odcisk      = excluded.odcisk,
	                         klucz_jawny = excluded.klucz_jawny,
	                         sciezka     = excluded.sciezka,
	                         haslo       = excluded.haslo`

	usunKluczTerminala = `DELETE FROM terminal_klucz WHERE kod = ?`

	wstawTunelTerminala = `INSERT INTO terminal_tunel
	                       (kod, okno_kod, rodzaj, host_kod, cel, port_lokalny,
	                        host_docelowy, port_docelowy, stan, powod)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(kod) DO UPDATE SET
	                         port_lokalny = excluded.port_lokalny,
	                         stan         = excluded.stan,
	                         powod        = excluded.powod`

	// Zamknięcie znaczy chwilę tylko wtedy, gdy tunel faktycznie się domyka;
	// przestawienie stanu bez zamknięcia zostawia kolumnę `zamknieto` nietkniętą.
	zmienStanTuneluTerminala = `UPDATE terminal_tunel
	                            SET stan = ?, powod = ?,
	                                zamknieto = CASE WHEN ? = 1
	                                            THEN strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                            ELSE zamknieto END
	                            WHERE kod = ?`

	osierocTuneleTerminala = `UPDATE terminal_tunel
	                          SET stan = 'inactive',
	                              powod = 'rdzeń został uruchomiony ponownie, a tunel biegł w poprzednim biegu',
	                              zamknieto = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE stan = 'active'`

	wstawObserwacjeTerminala = `INSERT INTO terminal_obserwacja
	                            (kod, okno_kod, karta_kod, wzorzec, polecenie, tlumienie,
	                             rekurencyjnie, stan, powod)
	                            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                            ON CONFLICT(kod) DO UPDATE SET
	                              wzorzec   = excluded.wzorzec,
	                              polecenie = excluded.polecenie,
	                              tlumienie = excluded.tlumienie,
	                              stan      = excluded.stan,
	                              powod     = excluded.powod`

	zmienStanObserwacjiTerminala = `UPDATE terminal_obserwacja SET stan = ?, powod = ? WHERE kod = ?`

	odnotujWyzwolenieTerminala = `UPDATE terminal_obserwacja
	                              SET licznik = licznik + 1,
	                                  wyzwolono = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                              WHERE kod = ?`

	osierocObserwacjeTerminala = `UPDATE terminal_obserwacja
	                              SET stan = 'stopped',
	                                  powod = 'rdzeń został uruchomiony ponownie, a obserwacja biegła w poprzednim biegu'
	                              WHERE stan = 'active'`
)

// ZapiszHosta zakłada wpis książki hostów albo podmienia istniejący.
func (r *repozytoriumTerminala) ZapiszHosta(ctx context.Context, host HostTerminala) error {
	if host.Kod == "" || host.Nazwa == "" || host.Cel == "" {
		return fmt.Errorf("dane: wpis hosta bez identyfikatora, nazwy albo adresu celu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawHostaTerminala)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, host.Kod, host.Nazwa, host.Cel, host.Port,
		host.Grupa, host.KatalogRoboczy, host.KluczKod, host.HostPosredniKod,
		host.Notatka); err != nil {
		return fmt.Errorf("dane: nie można zapisać wpisu hosta %q: %w", host.Kod, err)
	}
	return nil
}

// UsunHosta zdejmuje wpis książki hostów. Fałsz znaczy, że wpisu nie było.
func (r *repozytoriumTerminala) UsunHosta(ctx context.Context, kod string) (bool, error) {
	return r.usunWpisTerminala(ctx, usunHostaTerminala, kod, "wpisu hosta")
}

// OdepnijKlucz zdejmuje wskazanie klucza z wpisów, które go używały, i oddaje
// ich kody. Odczyt i zmiana idą jedną transakcją: między nimi nie ma prawa
// wejść zapis, po którym oddany wykaz przestałby odpowiadać stanowi bazy.
func (r *repozytoriumTerminala) OdepnijKlucz(ctx context.Context, kluczKod string) ([]string, error) {
	if kluczKod == "" {
		return nil, nil
	}
	kody := make([]string, 0, 4)
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wiersze, err := transakcja.QueryContext(ctx, hostyPoKluczuTerminala, kluczKod)
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać wpisów wskazujących klucz %q: %w", kluczKod, err)
		}
		defer wiersze.Close()
		for wiersze.Next() {
			var kod string
			if err := wiersze.Scan(&kod); err != nil {
				return fmt.Errorf("dane: nieczytelny wiersz wpisu hosta: %w", err)
			}
			kody = append(kody, kod)
		}
		if err := wiersze.Err(); err != nil {
			return err
		}
		if _, err := transakcja.ExecContext(ctx, odepnijKluczTerminala, kluczKod); err != nil {
			return fmt.Errorf("dane: nie można odpiąć klucza %q od wpisów hostów: %w", kluczKod, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return kody, nil
}

// ZapiszSkrypt zakłada pozycję biblioteki albo dokłada jej kolejną wersję.
//
// Numer wersji nadaje baza wyrażeniem `wersja + 1` wykonanym w transakcji wraz
// z wpisem wersji. Odczytanie numeru osobnym zapytaniem przed zapisem dałoby
// dwóm równoległym zapisom ten sam numer, a warunek UNIQUE na parze
// (pozycja, wersja) odrzuciłby drugi z nich.
func (r *repozytoriumTerminala) ZapiszSkrypt(ctx context.Context,
	skrypt SkryptTerminala) (int64, bool, error) {

	if skrypt.Kod == "" || skrypt.Nazwa == "" {
		return 0, false, fmt.Errorf("dane: pozycja biblioteki bez identyfikatora albo nazwy")
	}
	rodzaj := skrypt.Rodzaj
	if rodzaj == "" {
		rodzaj = shared.TerminalScriptKindScript
	}
	var wersja int64
	powstala := false
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wynik, err := transakcja.ExecContext(ctx, podmienSkryptTerminala,
			skrypt.Nazwa, rodzaj, skrypt.Powloka, skrypt.Tresc, skrypt.Znaczniki,
			skrypt.Alias, skrypt.Kod)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać pozycji biblioteki %q: %w", skrypt.Kod, err)
		}
		zmienione, err := wynik.RowsAffected()
		if err != nil {
			return fmt.Errorf("dane: nie można policzyć zmienionych wierszy biblioteki: %w", err)
		}
		if zmienione == 0 {
			if _, err := transakcja.ExecContext(ctx, wstawSkryptTerminala,
				skrypt.Kod, skrypt.Nazwa, rodzaj, skrypt.Powloka, skrypt.Tresc,
				skrypt.Znaczniki, skrypt.Alias); err != nil {
				return fmt.Errorf("dane: nie można założyć pozycji biblioteki %q: %w", skrypt.Kod, err)
			}
			powstala = true
		}
		if err := transakcja.QueryRowContext(ctx, wersjaSkryptuTerminala, skrypt.Kod).
			Scan(&wersja); err != nil {
			return fmt.Errorf("dane: nie można odczytać wersji pozycji %q: %w", skrypt.Kod, err)
		}
		if _, err := transakcja.ExecContext(ctx, wstawWersjeSkryptuTerminala,
			skrypt.Kod, wersja, skrypt.Tresc); err != nil {
			return fmt.Errorf("dane: nie można zapisać wersji %d pozycji %q: %w", wersja, skrypt.Kod, err)
		}
		return nil
	})
	if err != nil {
		return 0, false, err
	}
	return wersja, powstala, nil
}

// UsunSkrypt usuwa pozycję biblioteki wraz ze wszystkimi jej wersjami
// (kasowanie kaskadowe, migracja 247).
func (r *repozytoriumTerminala) UsunSkrypt(ctx context.Context, kod string) (bool, error) {
	return r.usunWpisTerminala(ctx, usunSkryptTerminala, kod, "pozycji biblioteki")
}

// ZapiszKlucz wciąga klucz do wykazu albo odświeża jego wpis.
func (r *repozytoriumTerminala) ZapiszKlucz(ctx context.Context, klucz KluczTerminala) error {
	if klucz.Kod == "" || klucz.Nazwa == "" || klucz.Sciezka == "" {
		return fmt.Errorf("dane: wpis klucza bez identyfikatora, nazwy albo ścieżki")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKluczTerminala)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, klucz.Kod, klucz.Nazwa, klucz.Rodzaj,
		klucz.Odcisk, klucz.KluczJawny, klucz.Sciezka, klucz.Haslo); err != nil {
		return fmt.Errorf("dane: nie można zapisać klucza %q: %w", klucz.Kod, err)
	}
	return nil
}

// UsunKlucz zdejmuje klucz z wykazu. Pliki klucza na dysku to osobna czynność
// rdzenia — baza ich nie tyka.
func (r *repozytoriumTerminala) UsunKlucz(ctx context.Context, kod string) (bool, error) {
	return r.usunWpisTerminala(ctx, usunKluczTerminala, kod, "klucza")
}

// ZapiszTunel zakłada wiersz tunelu albo odświeża jego stan.
func (r *repozytoriumTerminala) ZapiszTunel(ctx context.Context, tunel TunelTerminala) error {
	if tunel.Kod == "" || tunel.OknoKod == "" {
		return fmt.Errorf("dane: tunel bez identyfikatora albo okna")
	}
	stan := tunel.Stan
	if stan == "" {
		stan = shared.TerminalTunnelStatusInactive
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawTunelTerminala)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, tunel.Kod, tunel.OknoKod, tunel.Rodzaj,
		tunel.HostKod, tunel.Cel, tunel.PortLokalny, tunel.HostDocelowy,
		tunel.PortDocelowy, stan, tunel.Powod); err != nil {
		return fmt.Errorf("dane: nie można zapisać tunelu %q: %w", tunel.Kod, err)
	}
	return nil
}

// ZmienStanTunelu przestawia stan tunelu wraz z powodem.
func (r *repozytoriumTerminala) ZmienStanTunelu(ctx context.Context, kod string,
	stan shared.TerminalTunnelStatus, powod string, zamkniety bool) error {

	if kod == "" {
		return fmt.Errorf("dane: zmiana stanu tunelu bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienStanTuneluTerminala)
	if err != nil {
		return err
	}
	znacznik := 0
	if zamkniety {
		znacznik = 1
	}
	if _, err := polecenie.ExecContext(ctx, stan, powod, znacznik, kod); err != nil {
		return fmt.Errorf("dane: nie można przestawić stanu tunelu %q: %w", kod, err)
	}
	return nil
}

// OsierocTunele przestawia tunele poprzedniego biegu rdzenia na `inactive`.
func (r *repozytoriumTerminala) OsierocTunele(ctx context.Context) (int64, error) {
	return r.osierocWpisyTerminala(ctx, osierocTuneleTerminala, "tuneli")
}

// ZapiszObserwacje zakłada obserwację plików albo odświeża jej wpis.
func (r *repozytoriumTerminala) ZapiszObserwacje(ctx context.Context, obserwacja ObserwacjaTerminala) error {
	if obserwacja.Kod == "" || obserwacja.KartaKod == "" {
		return fmt.Errorf("dane: obserwacja bez identyfikatora albo karty")
	}
	stan := obserwacja.Stan
	if stan == "" {
		stan = shared.TerminalWatchStatusActive
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawObserwacjeTerminala)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, obserwacja.Kod, obserwacja.OknoKod,
		obserwacja.KartaKod, obserwacja.Wzorzec, obserwacja.Polecenie,
		obserwacja.Tlumienie, obserwacja.Rekurencyjnie, stan, obserwacja.Powod); err != nil {
		return fmt.Errorf("dane: nie można zapisać obserwacji %q: %w", obserwacja.Kod, err)
	}
	return nil
}

// ZmienStanObserwacji przestawia stan obserwacji wraz z powodem.
func (r *repozytoriumTerminala) ZmienStanObserwacji(ctx context.Context, kod string,
	stan shared.TerminalWatchStatus, powod string) error {

	if kod == "" {
		return fmt.Errorf("dane: zmiana stanu obserwacji bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienStanObserwacjiTerminala)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, stan, powod, kod); err != nil {
		return fmt.Errorf("dane: nie można przestawić stanu obserwacji %q: %w", kod, err)
	}
	return nil
}

// OdnotujWyzwolenie podnosi licznik wyzwoleń obserwacji.
func (r *repozytoriumTerminala) OdnotujWyzwolenie(ctx context.Context, kod string) error {
	if kod == "" {
		return fmt.Errorf("dane: odnotowanie wyzwolenia bez identyfikatora obserwacji")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, odnotujWyzwolenieTerminala)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, kod); err != nil {
		return fmt.Errorf("dane: nie można odnotować wyzwolenia obserwacji %q: %w", kod, err)
	}
	return nil
}

// OsierocObserwacje przestawia obserwacje poprzedniego biegu rdzenia na `stopped`.
func (r *repozytoriumTerminala) OsierocObserwacje(ctx context.Context) (int64, error) {
	return r.osierocWpisyTerminala(ctx, osierocObserwacjeTerminala, "obserwacji")
}

// usunWpisTerminala wykonuje kasowanie po kodzie i oddaje prawdę o skutku.
func (r *repozytoriumTerminala) usunWpisTerminala(ctx context.Context,
	zapytanie, kod, czego string) (bool, error) {

	if kod == "" {
		return false, fmt.Errorf("dane: usunięcie %s bez identyfikatora", czego)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć %s %q: %w", czego, kod, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć usuniętych wierszy: %w", err)
	}
	return usuniete > 0, nil
}

// osierocWpisyTerminala wykonuje jedno osierocenie i oddaje liczbę wierszy.
func (r *repozytoriumTerminala) osierocWpisyTerminala(ctx context.Context,
	zapytanie, czego string) (int64, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można osierocić %s terminala: %w", czego, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return zmienione, nil
}
