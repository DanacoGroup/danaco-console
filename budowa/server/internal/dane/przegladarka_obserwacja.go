// Odpowiedzialność pliku: to, co moduł Browser obserwuje między jedną sesją
// a drugą — monitory zmian stron (migracja 172), kanały RSS/Atom wraz z ich
// wpisami (173), kolejka czytania (174) i zakładki (171), zasilane z Capture
// & Monitor Panel.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// MonitorPrzegladania to wiersz tabeli `monitor_przegladania`, niosący adres,
// próg zmiany, kanał powiadomienia i odwołanie do treści pilnowanej strony.
type MonitorPrzegladania struct {
	ID                   int64
	Kod                  string
	Okno                 string
	Url                  string
	Selektor             *string
	InterwalSekund       int64
	ProgZmiany           *int64
	KanalPowiadomienia   *string
	Wlaczony             bool
	Stan                 string
	OdniesienieOdwolanie *string
	OdniesienieDlugosc   *int64
	Sprawdzono           *string
	Zmieniono            *string
	Utworzono            string
}

// KanalPrzegladania to wiersz tabeli `kanal_przegladania` — jedna subskrypcja
// kanału RSS albo Atom, niosąca adres, tytuł, postać i chwilę ostatniego pobrania.
type KanalPrzegladania struct {
	ID             int64
	Kod            string
	Okno           string
	Url            string
	Tytul          *string
	Postac         *string
	InterwalSekund *int64
	Pobrano        *string
	Utworzono      string
}

// WpisKanalu to wiersz tabeli `wpis_kanalu_przegladania`, niosący adres,
// tytuł, streszczenie i znamię przeczytania jednej pozycji kanału.
type WpisKanalu struct {
	ID           int64
	Kod          string
	Kanal        string
	Url          string
	Tytul        *string
	Streszczenie *string
	Przeczytany  bool
	Opublikowano *string
	Utworzono    string
}

// PozycjaCzytania to wiersz tabeli `pozycja_czytania_przegladania`, niosący
// adres, notatkę, znamię przeczytania i chwilę przypomnienia.
type PozycjaCzytania struct {
	ID            int64
	Kod           string
	Okno          string
	Url           string
	Tytul         *string
	Notatka       *string
	Przeczytana   bool
	Przypomnienie *string
	Utworzono     string
}

// ZakladkaPrzegladania to wiersz tabeli `zakladka_przegladania`, niosący
// adres, tytuł, folder, etykiety tekstem JSON i notatkę.
type ZakladkaPrzegladania struct {
	ID           int64
	Kod          string
	Okno         string
	Url          string
	Tytul        *string
	Folder       *string
	EtykietyJson *string
	Notatka      *string
	Utworzono    string
}

// FiltrZakladek zawęża wykaz zakładek — obsługuje pola `browser.bookmark.list`.
// Szukanie idzie po tytule, adresie i notatce naraz, bo tak też szuka Operator:
// pamięta jedno z trzech, nie wie które.
type FiltrZakladek struct {
	Okno   string
	Folder string
	Szukaj string
	Limit  int
}

const (
	kolumnyMonitora = `id, identyfikator_zewnetrzny, okno, url, selektor, interwal_sekund,
	                   prog_zmiany, kanal_powiadomienia, wlaczony, stan,
	                   odniesienie_odwolanie, odniesienie_dlugosc, sprawdzono, zmieniono, utworzono`

	zapiszMonitor = `INSERT INTO monitor_przegladania
	                 (identyfikator_zewnetrzny, okno, url, selektor, interwal_sekund, prog_zmiany,
	                  kanal_powiadomienia, wlaczony, stan, odniesienie_odwolanie,
	                  odniesienie_dlugosc, sprawdzono, zmieniono)
	                  konto_id)
	                 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                 ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                     okno = excluded.okno,
	                     url = excluded.url,
	                     selektor = excluded.selektor,
	                     interwal_sekund = excluded.interwal_sekund,
	                     prog_zmiany = excluded.prog_zmiany,
	                     kanal_powiadomienia = excluded.kanal_powiadomienia,
	                     wlaczony = excluded.wlaczony,
	                     stan = excluded.stan,
	                     odniesienie_odwolanie = excluded.odniesienie_odwolanie,
	                     odniesienie_dlugosc = excluded.odniesienie_dlugosc,
	                     sprawdzono = excluded.sprawdzono,
	                     zmieniono = excluded.zmieniono
	                 WHERE ` + WarunekKonta

	pobierzMonitor = `SELECT ` + kolumnyMonitora + `
	                  FROM monitor_przegladania
	                  WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaMonitorow = `SELECT ` + kolumnyMonitora + ` FROM monitor_przegladania
	                  WHERE (? = '' OR okno = ?) AND (? = 0 OR wlaczony = 1)
	                    AND ` + WarunekKonta + `
	                  ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunMonitor = `DELETE FROM monitor_przegladania
	               WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	kolumnyKanaluObserwacji = `id, identyfikator_zewnetrzny, okno, url, tytul, postac,
	                 interwal_sekund, pobrano, utworzono`

	zapiszKanalObserwacji = `INSERT INTO kanal_przegladania
	               (identyfikator_zewnetrzny, okno, url, tytul, postac, interwal_sekund, pobrano)
	                konto_id)
	               VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	               ON CONFLICT(okno, url) DO UPDATE SET
	                   tytul = excluded.tytul,
	                   postac = excluded.postac,
	                   interwal_sekund = excluded.interwal_sekund,
	                   pobrano = excluded.pobrano
	               WHERE ` + WarunekKonta

	pobierzKanalObserwacji = `SELECT ` + kolumnyKanaluObserwacji + ` FROM kanal_przegladania
	                WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	pobierzKanalObserwacjiPoAdresie = `SELECT ` + kolumnyKanaluObserwacji + ` FROM kanal_przegladania
	                         WHERE okno = ? AND url = ? AND ` + WarunekKonta

	listaKanalowObserwacji = `SELECT ` + kolumnyKanaluObserwacji + ` FROM kanal_przegladania
	                WHERE (? = '' OR okno = ?) AND (? = '' OR identyfikator_zewnetrzny = ?)
	                  AND ` + WarunekKonta + `
	                ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunKanalObserwacji = `DELETE FROM kanal_przegladania
	                       WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	kolumnyWpisuKanalu = `id, identyfikator_zewnetrzny, kanal_zewnetrzny_id, url, tytul,
	                      streszczenie, przeczytany, opublikowano, utworzono`

	zapiszWpisKanalu = `INSERT INTO wpis_kanalu_przegladania
	                    (identyfikator_zewnetrzny, kanal_zewnetrzny_id, url, tytul,
	                     streszczenie, przeczytany, opublikowano)
	                    VALUES (?, ?, ?, ?, ?, ?, ?)
	                    ON CONFLICT(kanal_zewnetrzny_id, url) DO UPDATE SET
	                        tytul = excluded.tytul,
	                        streszczenie = excluded.streszczenie,
	                        opublikowano = excluded.opublikowano`

	listaWpisowKanalu = `SELECT ` + kolumnyWpisuKanalu + ` FROM wpis_kanalu_przegladania
	                     WHERE kanal_zewnetrzny_id = ? AND (? = 0 OR przeczytany = 0)
	                     ORDER BY opublikowano DESC, id DESC LIMIT ?`

	liczbaNieprzeczytanych = `SELECT COUNT(*) FROM wpis_kanalu_przegladania
	                          WHERE kanal_zewnetrzny_id = ? AND przeczytany = 0`

	kolumnyPozycjiCzytania = `id, identyfikator_zewnetrzny, okno, url, tytul, notatka,
	                          przeczytana, przypomnienie, utworzono`

	zapiszPozycjeCzytania = `INSERT INTO pozycja_czytania_przegladania
	                         (identyfikator_zewnetrzny, okno, url, tytul, notatka,
	                          przeczytana, przypomnienie)
	                          konto_id)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                         ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                             url = excluded.url,
	                             tytul = excluded.tytul,
	                             notatka = excluded.notatka,
	                             przeczytana = excluded.przeczytana,
	                             przypomnienie = excluded.przypomnienie
	                         WHERE ` + WarunekKonta

	pobierzPozycjeCzytania = `SELECT ` + kolumnyPozycjiCzytania + `
	                          FROM pozycja_czytania_przegladania
	                          WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaPozycjiCzytania = `SELECT ` + kolumnyPozycjiCzytania + `
	                        FROM pozycja_czytania_przegladania
	                        WHERE (? = '' OR okno = ?) AND (? = 0 OR przeczytana = 0)
	                          AND ` + WarunekKonta + `
	                        ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunPozycjeCzytania = `DELETE FROM pozycja_czytania_przegladania
	                       WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	kolumnyZakladki = `id, identyfikator_zewnetrzny, okno, url, tytul, folder,
	                   etykiety_json, notatka, utworzono`

	// Więz UNIQUE na `identyfikator_zewnetrzny` obejmuje całą tabelę, więc
	// warunek konta w gałęzi DO UPDATE zostawia wiersz cudzy nietknięty.
	zapiszZakladke = `INSERT INTO zakladka_przegladania
	                  (identyfikator_zewnetrzny, okno, url, tytul, folder, etykiety_json, notatka,
	                   konto_id)
	                  VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                  ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                      url = excluded.url,
	                      tytul = excluded.tytul,
	                      folder = excluded.folder,
	                      etykiety_json = excluded.etykiety_json,
	                      notatka = excluded.notatka
	                  WHERE ` + WarunekKonta

	pobierzZakladke = `SELECT ` + kolumnyZakladki + ` FROM zakladka_przegladania
	                   WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaZakladek = `SELECT ` + kolumnyZakladki + ` FROM zakladka_przegladania
	                 WHERE (? = '' OR okno = ?) AND (? = '' OR folder = ?)
	                   AND (? = '' OR url LIKE ? OR IFNULL(tytul,'') LIKE ? OR IFNULL(notatka,'') LIKE ?)
	                   AND ` + WarunekKonta + `
	                 ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunZakladke = `DELETE FROM zakladka_przegladania
	                WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

// ZapiszMonitor zakłada monitor albo nadpisuje zastany, przyjmując stan
// „pending", gdy wołający go nie poda.
func (r *repozytoriumPrzegladania) ZapiszMonitor(ctx context.Context,
	monitor MonitorPrzegladania) (MonitorPrzegladania, error) {

	if monitor.Kod == "" || monitor.Okno == "" || monitor.Url == "" {
		return MonitorPrzegladania{}, fmt.Errorf("dane: monitor przeglądania bez identyfikatora, okna albo adresu")
	}
	if monitor.Stan == "" {
		monitor.Stan = "pending"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszMonitor)
	if err != nil {
		return MonitorPrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, monitor.Kod, monitor.Okno, monitor.Url,
		tekstDoKolumny(monitor.Selektor), monitor.InterwalSekund, liczbaDoKolumny(monitor.ProgZmiany),
		tekstDoKolumny(monitor.KanalPowiadomienia), liczbaLogiczna(monitor.Wlaczony), monitor.Stan,
		tekstDoKolumny(monitor.OdniesienieOdwolanie), liczbaDoKolumny(monitor.OdniesienieDlugosc),
		tekstDoKolumny(monitor.Sprawdzono), tekstDoKolumny(monitor.Zmieniono),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return MonitorPrzegladania{}, fmt.Errorf("dane: nie można zapisać monitora %q: %w", monitor.Kod, err)
	}
	return r.Monitor(ctx, monitor.Kod)
}

// Monitor oddaje monitor o wskazanym kodzie zewnętrznym albo błąd
// ErrBrakWiersza, gdy monitor o tym kodzie nie istnieje.
func (r *repozytoriumPrzegladania) Monitor(ctx context.Context, kod string) (MonitorPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzMonitor)
	if err != nil {
		return MonitorPrzegladania{}, err
	}
	monitor, err := odczytajMonitorPrzegladania(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return MonitorPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return MonitorPrzegladania{}, fmt.Errorf("dane: nieczytelny monitor %q: %w", kod, err)
	}
	return monitor, nil
}

// Monitory oddaje monitory, opcjonalnie zawężone do okna i do włączonych,
// od najświeższego, do granicy podanego limitu.
func (r *repozytoriumPrzegladania) Monitory(ctx context.Context, okno string,
	tylkoWlaczone bool, limit int) ([]MonitorPrzegladania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaMonitorow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, okno,
		liczbaLogiczna(tylkoWlaczone), KontoOperatora(ctx), granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać monitorów: %w", err)
	}
	defer wiersze.Close()

	lista := []MonitorPrzegladania{}
	for wiersze.Next() {
		monitor, err := odczytajMonitorPrzegladania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz monitorów: %w", err)
		}
		lista = append(lista, monitor)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt monitorów: %w", err)
	}
	return lista, nil
}

// UsunMonitor zdejmuje monitor wraz z jego odniesieniem do treści pilnowanej
// strony i mówi, czy monitor o wskazanym kodzie istniał.
func (r *repozytoriumPrzegladania) UsunMonitor(ctx context.Context, kod string) (bool, error) {
	return r.usunWiersz(ctx, usunMonitor, kod, "monitor", KontoOperatora(ctx))
}

// ZapiszKanal zakłada subskrypcję albo odświeża zastaną, wskazaną parą
// okno i adres, i oddaje zapisany wiersz.
func (r *repozytoriumPrzegladania) ZapiszKanal(ctx context.Context,
	kanal KanalPrzegladania) (KanalPrzegladania, error) {

	if kanal.Kod == "" || kanal.Okno == "" || kanal.Url == "" {
		return KanalPrzegladania{}, fmt.Errorf("dane: kanał przeglądania bez identyfikatora, okna albo adresu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKanalObserwacji)
	if err != nil {
		return KanalPrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, kanal.Kod, kanal.Okno, kanal.Url,
		tekstDoKolumny(kanal.Tytul), tekstDoKolumny(kanal.Postac),
		liczbaDoKolumny(kanal.InterwalSekund), tekstDoKolumny(kanal.Pobrano),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return KanalPrzegladania{}, fmt.Errorf("dane: nie można zapisać kanału %q: %w", kanal.Kod, err)
	}
	// Odczyt idzie po parze okno + adres, bo subskrypcja odświeża wiersz
	// zastany pod jego dawnym kodem.
	return r.KanalPoAdresie(ctx, kanal.Okno, kanal.Url)
}

// Kanal oddaje kanał o wskazanym kodzie zewnętrznym albo błąd ErrBrakWiersza,
// gdy kanał o tym kodzie nie istnieje.
func (r *repozytoriumPrzegladania) Kanal(ctx context.Context, kod string) (KanalPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKanalObserwacji)
	if err != nil {
		return KanalPrzegladania{}, err
	}
	kanal, err := odczytajKanalObserwacji(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return KanalPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return KanalPrzegladania{}, fmt.Errorf("dane: nieczytelny kanał %q: %w", kod, err)
	}
	return kanal, nil
}

// KanalPoAdresie oddaje kanał subskrybowany w oknie pod wskazanym adresem
// albo błąd ErrBrakWiersza, gdy taka subskrypcja nie istnieje.
func (r *repozytoriumPrzegladania) KanalPoAdresie(ctx context.Context, okno, url string) (KanalPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKanalObserwacjiPoAdresie)
	if err != nil {
		return KanalPrzegladania{}, err
	}
	kanal, err := odczytajKanalObserwacji(polecenie.QueryRowContext(ctx, okno, url, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return KanalPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return KanalPrzegladania{}, fmt.Errorf("dane: nieczytelny kanał %q okna %q: %w", url, okno, err)
	}
	return kanal, nil
}

// Kanaly oddaje kanały, opcjonalnie zawężone do okna albo do jednego kanału,
// od najświeższego, do granicy podanego limitu.
func (r *repozytoriumPrzegladania) Kanaly(ctx context.Context, okno, kod string, limit int) ([]KanalPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKanalowObserwacji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, okno, kod, kod,
		KontoOperatora(ctx), granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kanałów: %w", err)
	}
	defer wiersze.Close()

	lista := []KanalPrzegladania{}
	for wiersze.Next() {
		kanal, err := odczytajKanalObserwacji(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kanałów: %w", err)
		}
		lista = append(lista, kanal)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kanałów: %w", err)
	}
	return lista, nil
}

// UsunKanal zdejmuje subskrypcję o wskazanym kodzie i mówi, czy istniała;
// wpisy kanału znikają kaskadą schematu bazy.
func (r *repozytoriumPrzegladania) UsunKanal(ctx context.Context, kod string) (bool, error) {
	return r.usunWiersz(ctx, usunKanalObserwacji, kod, "kanał", KontoOperatora(ctx))
}

// ZapiszWpisKanalu dopisuje wpis kanału albo odświeża zastany (para kanał +
// adres). Oznaczenie przeczytania zostaje przy odświeżeniu nietknięte — wpis
// przeczytany nie ma prawa wrócić jako nowy przy kolejnym odpytaniu kanału.
func (r *repozytoriumPrzegladania) ZapiszWpisKanalu(ctx context.Context, wpis WpisKanalu) error {
	if wpis.Kod == "" || wpis.Kanal == "" || wpis.Url == "" {
		return fmt.Errorf("dane: wpis kanału bez identyfikatora, kanału albo adresu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWpisKanalu)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wpis.Kod, wpis.Kanal, wpis.Url,
		tekstDoKolumny(wpis.Tytul), tekstDoKolumny(wpis.Streszczenie),
		liczbaLogiczna(wpis.Przeczytany), tekstDoKolumny(wpis.Opublikowano))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać wpisu kanału %q: %w", wpis.Kod, err)
	}
	return nil
}

// WpisyKanalu oddaje wpisy kanału od najświeższego, opcjonalnie zawężone do
// nieprzeczytanych, do granicy podanego limitu.
func (r *repozytoriumPrzegladania) WpisyKanalu(ctx context.Context, kanal string,
	tylkoNieprzeczytane bool, limit int) ([]WpisKanalu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWpisowKanalu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kanal, liczbaLogiczna(tylkoNieprzeczytane), granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wpisów kanału %q: %w", kanal, err)
	}
	defer wiersze.Close()

	lista := []WpisKanalu{}
	for wiersze.Next() {
		var wpis WpisKanalu
		var tytul, streszczenie, opublikowano sql.NullString
		var przeczytany int64
		if err := wiersze.Scan(&wpis.ID, &wpis.Kod, &wpis.Kanal, &wpis.Url, &tytul,
			&streszczenie, &przeczytany, &opublikowano, &wpis.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wpisów kanału %q: %w", kanal, err)
		}
		wpis.Tytul = tekstZKolumny(tytul)
		wpis.Streszczenie = tekstZKolumny(streszczenie)
		wpis.Opublikowano = tekstZKolumny(opublikowano)
		wpis.Przeczytany = przeczytany == 1
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wpisów kanału %q: %w", kanal, err)
	}
	return lista, nil
}

// NieprzeczytaneKanalu liczy wpisy wskazanego kanału bez oznaczenia
// przeczytania, do wyświetlenia jako odznaka liczby nowości kanału.
func (r *repozytoriumPrzegladania) NieprzeczytaneKanalu(ctx context.Context, kanal string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, liczbaNieprzeczytanych)
	if err != nil {
		return 0, err
	}
	var liczba int64
	if err := polecenie.QueryRowContext(ctx, kanal).Scan(&liczba); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć wpisów kanału %q: %w", kanal, err)
	}
	return liczba, nil
}

// ZapiszPozycjeCzytania zakłada pozycję kolejki czytania albo nadpisuje
// zastaną o tym samym kodzie i oddaje zapisany wiersz.
func (r *repozytoriumPrzegladania) ZapiszPozycjeCzytania(ctx context.Context,
	pozycja PozycjaCzytania) (PozycjaCzytania, error) {

	if pozycja.Kod == "" || pozycja.Okno == "" || pozycja.Url == "" {
		return PozycjaCzytania{}, fmt.Errorf("dane: pozycja czytania bez identyfikatora, okna albo adresu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPozycjeCzytania)
	if err != nil {
		return PozycjaCzytania{}, err
	}
	_, err = polecenie.ExecContext(ctx, pozycja.Kod, pozycja.Okno, pozycja.Url,
		tekstDoKolumny(pozycja.Tytul), tekstDoKolumny(pozycja.Notatka),
		liczbaLogiczna(pozycja.Przeczytana), tekstDoKolumny(pozycja.Przypomnienie),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return PozycjaCzytania{}, fmt.Errorf("dane: nie można zapisać pozycji czytania %q: %w", pozycja.Kod, err)
	}
	return r.PozycjaCzytania(ctx, pozycja.Kod)
}

// PozycjaCzytania oddaje pozycję kolejki o wskazanym kodzie zewnętrznym albo
// błąd ErrBrakWiersza, gdy pozycja o tym kodzie nie istnieje.
func (r *repozytoriumPrzegladania) PozycjaCzytania(ctx context.Context, kod string) (PozycjaCzytania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPozycjeCzytania)
	if err != nil {
		return PozycjaCzytania{}, err
	}
	pozycja, err := odczytajPozycjeCzytania(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return PozycjaCzytania{}, ErrBrakWiersza
	}
	if err != nil {
		return PozycjaCzytania{}, fmt.Errorf("dane: nieczytelna pozycja czytania %q: %w", kod, err)
	}
	return pozycja, nil
}

// KolejkaCzytania oddaje pozycje kolejki, opcjonalnie zawężone do okna
// i do nieprzeczytanych, od najświeższej, do granicy podanego limitu.
func (r *repozytoriumPrzegladania) KolejkaCzytania(ctx context.Context, okno string,
	tylkoNieprzeczytane bool, limit int) ([]PozycjaCzytania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPozycjiCzytania)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, okno,
		liczbaLogiczna(tylkoNieprzeczytane), KontoOperatora(ctx), granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kolejki czytania: %w", err)
	}
	defer wiersze.Close()

	lista := []PozycjaCzytania{}
	for wiersze.Next() {
		pozycja, err := odczytajPozycjeCzytania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kolejki czytania: %w", err)
		}
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kolejki czytania: %w", err)
	}
	return lista, nil
}

// UsunPozycjeCzytania zdejmuje pozycję o wskazanym kodzie z kolejki czytania
// i mówi, czy pozycja istniała.
func (r *repozytoriumPrzegladania) UsunPozycjeCzytania(ctx context.Context, kod string) (bool, error) {
	return r.usunWiersz(ctx, usunPozycjeCzytania, kod, "pozycja czytania", KontoOperatora(ctx))
}

// ZapiszZakladke zakłada zakładkę albo nadpisuje zastaną o tym samym kodzie
// zewnętrznym i oddaje zapisany wiersz z aktualną treścią.
func (r *repozytoriumPrzegladania) ZapiszZakladke(ctx context.Context,
	zakladka ZakladkaPrzegladania) (ZakladkaPrzegladania, error) {

	if zakladka.Kod == "" || zakladka.Okno == "" || zakladka.Url == "" {
		return ZakladkaPrzegladania{}, fmt.Errorf("dane: zakładka bez identyfikatora, okna albo adresu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZakladke)
	if err != nil {
		return ZakladkaPrzegladania{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, zakladka.Kod, zakladka.Okno, zakladka.Url,
		tekstDoKolumny(zakladka.Tytul), tekstDoKolumny(zakladka.Folder),
		tekstDoKolumny(zakladka.EtykietyJson), tekstDoKolumny(zakladka.Notatka),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return ZakladkaPrzegladania{}, fmt.Errorf("dane: nie można zapisać zakładki %q: %w", zakladka.Kod, err)
	}
	// Kod zewnętrzny zajęty przez wiersz konta obcego daje zero zmienionych
	// wierszy; odczyt zwrotny oddałby „brak wiersza” zamiast powodu odmowy.
	if err := sprawdzTrafienieZapisu(wynik, "zakładka", zakladka.Kod); err != nil {
		return ZakladkaPrzegladania{}, err
	}
	return r.Zakladka(ctx, zakladka.Kod)
}

// Zakladka oddaje zakładkę o wskazanym kodzie zewnętrznym albo błąd
// ErrBrakWiersza, gdy zakładka o tym kodzie nie istnieje.
func (r *repozytoriumPrzegladania) Zakladka(ctx context.Context, kod string) (ZakladkaPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZakladke)
	if err != nil {
		return ZakladkaPrzegladania{}, err
	}
	zakladka, err := odczytajZakladke(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return ZakladkaPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return ZakladkaPrzegladania{}, fmt.Errorf("dane: nieczytelna zakładka %q: %w", kod, err)
	}
	return zakladka, nil
}

// Zakladki oddaje zakładki zawężone filtrem żądania: oknem, folderem i szukaną
// frazą łączną dla tytułu, adresu i notatki.
func (r *repozytoriumPrzegladania) Zakladki(ctx context.Context, filtr FiltrZakladek) ([]ZakladkaPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaZakladek)
	if err != nil {
		return nil, err
	}
	wzorzec := ""
	if filtr.Szukaj != "" {
		wzorzec = "%" + filtr.Szukaj + "%"
	}
	wiersze, err := polecenie.QueryContext(ctx, filtr.Okno, filtr.Okno, filtr.Folder, filtr.Folder,
		filtr.Szukaj, wzorzec, wzorzec, wzorzec, KontoOperatora(ctx), granicaWykazu(filtr.Limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zakładek: %w", err)
	}
	defer wiersze.Close()

	lista := []ZakladkaPrzegladania{}
	for wiersze.Next() {
		zakladka, err := odczytajZakladke(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zakładek: %w", err)
		}
		lista = append(lista, zakladka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zakładek: %w", err)
	}
	return lista, nil
}

// UsunZakladke zdejmuje zakładkę o wskazanym kodzie zewnętrznym z wykazu okna
// i mówi, czy zakładka istniała.
func (r *repozytoriumPrzegladania) UsunZakladke(ctx context.Context, kod string) (bool, error) {
	return r.usunWiersz(ctx, usunZakladke, kod, "zakładka", KontoOperatora(ctx))
}

// usunWiersz jest wspólnym skasowaniem po kodzie zewnętrznym. Fałsz znaczy
// „nie było takiego wiersza" — wołający ma z czego zbudować odmowę `not_found`
// zamiast meldować usunięcie czegoś, czego nie było.
// Argumenty dalsze przyjmuje jako ogon, bo część zapytań niesie za kodem
// warunek konta, a część jeszcze nie.
func (r *repozytoriumPrzegladania) usunWiersz(ctx context.Context, zapytanie, kod, nazwa string,
	dalsze ...any) (bool, error) {

	if kod == "" {
		return false, fmt.Errorf("dane: %s bez identyfikatora do usunięcia", nazwa)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, append([]any{kod}, dalsze...)...)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć wiersza (%s) %q: %w", nazwa, kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku usunięcia (%s) %q: %w", nazwa, kod, err)
	}
	return zmienione > 0, nil
}

// odczytajMonitorPrzegladania składa strukturę monitora z jednego wiersza
// wyniku zapytania opartego na `kolumnyMonitora`.
func odczytajMonitorPrzegladania(wiersz skaner) (MonitorPrzegladania, error) {
	var monitor MonitorPrzegladania
	var selektor, kanal, odwolanie, sprawdzono, zmieniono sql.NullString
	var prog, dlugosc sql.NullInt64
	var wlaczony int64
	err := wiersz.Scan(&monitor.ID, &monitor.Kod, &monitor.Okno, &monitor.Url, &selektor,
		&monitor.InterwalSekund, &prog, &kanal, &wlaczony, &monitor.Stan,
		&odwolanie, &dlugosc, &sprawdzono, &zmieniono, &monitor.Utworzono)
	if err != nil {
		return MonitorPrzegladania{}, err
	}
	monitor.Selektor = tekstZKolumny(selektor)
	monitor.ProgZmiany = liczbaZKolumny(prog)
	monitor.KanalPowiadomienia = tekstZKolumny(kanal)
	monitor.Wlaczony = wlaczony == 1
	monitor.OdniesienieOdwolanie = tekstZKolumny(odwolanie)
	monitor.OdniesienieDlugosc = liczbaZKolumny(dlugosc)
	monitor.Sprawdzono = tekstZKolumny(sprawdzono)
	monitor.Zmieniono = tekstZKolumny(zmieniono)
	return monitor, nil
}

// odczytajKanalObserwacji składa strukturę kanału z jednego wiersza wyniku
// zapytania opartego na `kolumnyKanaluObserwacji`.
func odczytajKanalObserwacji(wiersz skaner) (KanalPrzegladania, error) {
	var kanal KanalPrzegladania
	var tytul, postac, pobrano sql.NullString
	var interwal sql.NullInt64
	if err := wiersz.Scan(&kanal.ID, &kanal.Kod, &kanal.Okno, &kanal.Url, &tytul,
		&postac, &interwal, &pobrano, &kanal.Utworzono); err != nil {
		return KanalPrzegladania{}, err
	}
	kanal.Tytul = tekstZKolumny(tytul)
	kanal.Postac = tekstZKolumny(postac)
	kanal.InterwalSekund = liczbaZKolumny(interwal)
	kanal.Pobrano = tekstZKolumny(pobrano)
	return kanal, nil
}

// odczytajPozycjeCzytania składa strukturę pozycji kolejki z jednego wiersza
// wyniku zapytania opartego na `kolumnyPozycjiCzytania`.
func odczytajPozycjeCzytania(wiersz skaner) (PozycjaCzytania, error) {
	var pozycja PozycjaCzytania
	var tytul, notatka, przypomnienie sql.NullString
	var przeczytana int64
	if err := wiersz.Scan(&pozycja.ID, &pozycja.Kod, &pozycja.Okno, &pozycja.Url, &tytul,
		&notatka, &przeczytana, &przypomnienie, &pozycja.Utworzono); err != nil {
		return PozycjaCzytania{}, err
	}
	pozycja.Tytul = tekstZKolumny(tytul)
	pozycja.Notatka = tekstZKolumny(notatka)
	pozycja.Przypomnienie = tekstZKolumny(przypomnienie)
	pozycja.Przeczytana = przeczytana == 1
	return pozycja, nil
}

// odczytajZakladke składa strukturę zakładki z jednego wiersza wyniku
// zapytania opartego na `kolumnyZakladki`.
func odczytajZakladke(wiersz skaner) (ZakladkaPrzegladania, error) {
	var zakladka ZakladkaPrzegladania
	var tytul, folder, etykiety, notatka sql.NullString
	if err := wiersz.Scan(&zakladka.ID, &zakladka.Kod, &zakladka.Okno, &zakladka.Url,
		&tytul, &folder, &etykiety, &notatka, &zakladka.Utworzono); err != nil {
		return ZakladkaPrzegladania{}, err
	}
	zakladka.Tytul = tekstZKolumny(tytul)
	zakladka.Folder = tekstZKolumny(folder)
	zakladka.EtykietyJson = tekstZKolumny(etykiety)
	zakladka.Notatka = tekstZKolumny(notatka)
	return zakladka, nil
}
