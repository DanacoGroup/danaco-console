// Odpowiedzialność pliku: warstwa danych odcinka kontroli pracy modułu Studio,
// część druga — znakowanie fragmentów wraz z rodzajami znaczników własnych
// (migracja 365), zajęcia fragmentów przez wykonawców i spięcia o ten sam
// fragment (370), oraz wersje w rozbiciu na szereg Operatora i szereg
// autozapisu (kolumny dołożone migracją 367).
//
// ── Dlaczego wersje w szeregach czyta ten plik, a nie `studio_wersje.go` ─────
// `studio_wersje.go` zna wersję sprzed dobudowy: treść, etykietę, kamień
// milowy. Kolumny `szereg`, `postac_json` i tożsamości wykonawcy dołożyła
// migracja 367 i pyta o nie WYŁĄCZNIE ten odcinek — wykaz historii z
// przełącznikiem „pokaż także zapisy samoczynne" oraz porównanie POSTACI dwóch
// wersji. Dopisanie ich do tamtego pliku byłoby wejściem w plik cudzego
// odcinka; osobny odczyt tych samych wierszy nie zakłada drugiego pojęcia
// wersji, bo tabela jest jedna i przywraca się ją tą samą drogą.
//
// ── Dlaczego zajęcie fragmentu NIE jest blokadą ─────────────────────────────
// Blokada Operatora jest trwała i skierowana przeciw wykonawcom: zdejmuje ją
// wyłącznie Operator. Zajęcie fragmentu jest chwilowe, WYGASA samo i chroni
// przed drugim wykonawcą. Dwa różne byty, dwie tabele — zlanie ich dałoby
// blokadę, która wygasa (czyli żadną), albo zajęcie, którego nikt nie zdejmie
// po agencie ubitym w pół pracy.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ZnakowanieStudia to wiersz tabeli `znakowanie_studio` — wyróżnienie barwą,
// znacznik własny Operatora albo propozycja brzmienia na marginesie.
type ZnakowanieStudia struct {
	ID                   int64
	Kod                  string
	DokumentKod          string
	Rodzaj               string
	AutorRodzaj          string
	AutorAgentKod        *string
	AutorAgentNazwa      *string
	AutorAgentWersja     *string
	AutorPodagentKod     *string
	ZakresOd             int64
	ZakresDo             int64
	Barwa                *string
	ZnacznikNazwa        *string
	Tresc                *string
	BrzmienieProponowane *string
	BrzmienieZastane     *string
	Stan                 string
	CzynnoscKod          *string
	Utworzono            string
}

// RodzajZnacznikaStudia to wiersz tabeli `rodzaj_znacznika_studio`.
type RodzajZnacznikaStudia struct {
	Nazwa         string
	NazwaWidoczna string
	Barwa         *string
	Fabryczny     bool
	IleUzyc       int64
}

// ZajecieFragmentuStudia to wiersz tabeli `zajecie_fragmentu_studio`.
type ZajecieFragmentuStudia struct {
	ID              int64
	Kod             string
	DokumentKod     string
	WykonawcaRodzaj string
	AgentKod        *string
	AgentNazwa      *string
	PodagentKod     *string
	Stan            string
	ZakresOd        int64
	ZakresDo        int64
	ZadanieKod      *string
	Zajeto          string
	Wygasa          *string
	// Wygasle mówi, czy zajęcie przeterminowało się względem zegara BAZY.
	// Liczone w SQL, nie w rdzeniu: zegar rdzenia i zegar bazy rozjadą się przy
	// pierwszej różnicy strefy, a wtedy zajęcie trzymałoby fragment dłużej albo
	// krócej, niż mówi jego własna kolumna `wygasa`.
	Wygasle bool
}

// SpiecieWykonawcowStudia to wiersz tabeli `spiecie_wykonawcow_studio` — zapis
// tego, czyja zmiana weszła, czyja została odłożona i dlaczego.
type SpiecieWykonawcowStudia struct {
	ID                int64
	Kod               string
	DokumentKod       string
	ZakresOd          int64
	ZakresDo          int64
	WeszlaRodzaj      string
	WeszlaAgentKod    *string
	WeszlaAgentNazwa  *string
	OdlozonaRodzaj    string
	OdlozonaAgentKod  *string
	OdlozonaAgentNaz  *string
	Nastawa           string
	Powod             string
	BrzmienieOdlozone *string
	ZmianaSledzonaK   *string
	ZnakowanieKod     *string
	Domkniete         bool
	Utworzono         string
}

// WersjaSzereguStudia to wiersz `wersja_dokumentu_studio` widziany razem
// z kolumnami dołożonymi migracją 367: szeregiem, postacią i tożsamością
// wykonawcy.
type WersjaSzereguStudia struct {
	ID              int64
	Kod             string
	DokumentKod     string
	Etykieta        *string
	Podsumowanie    *string
	SkrotTresci     *string
	Tresc           *string
	PostacJSON      *string
	Autor           *string
	AutorAgentKod   *string
	AutorAgentNazwa *string
	KamienMilowy    bool
	GalazKod        *string
	PropozycjaKod   *string
	Szereg          string
	Utworzono       string
}

const (
	kolumnyZnakowaniaStudia = `z.id, z.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                           z.rodzaj, z.autor_rodzaj, z.autor_agent_kod, z.autor_agent_nazwa,
	                           z.autor_agent_wersja, z.autor_podagent_kod, z.zakres_od,
	                           z.zakres_do, z.barwa, z.znacznik_nazwa, z.tresc,
	                           z.brzmienie_proponowane, z.brzmienie_zastane, z.stan,
	                           z.czynnosc_kod, z.utworzono`

	zrodloZnakowaniaStudia = ` FROM znakowanie_studio z
	                           JOIN dokument_studio d ON d.id = z.dokument_id`

	znakowanieStudiaZapisz = `INSERT INTO znakowanie_studio
	                          (identyfikator_zewnetrzny, dokument_id, rodzaj, autor_rodzaj,
	                           autor_agent_kod, autor_agent_nazwa, autor_agent_wersja,
	                           autor_podagent_kod, zakres_od, zakres_do, barwa,
	                           znacznik_nazwa, tresc, brzmienie_proponowane,
	                           brzmienie_zastane, stan, czynnosc_kod)
	                          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	znakowanieStudiaPobierz = `SELECT ` + kolumnyZnakowaniaStudia + zrodloZnakowaniaStudia +
		` WHERE z.identyfikator_zewnetrzny = ?`

	// Kolejność wystąpienia w treści, bo wykaz znakowań jest SPISEM DO PRZEJŚCIA:
	// Operator skacze po nim od góry dokumentu do dołu i odhacza pozycje.
	znakowaniaStudiaLista = `SELECT ` + kolumnyZnakowaniaStudia + zrodloZnakowaniaStudia +
		` WHERE z.dokument_id = ? ORDER BY z.zakres_od, z.id`

	znakowanieStudiaUsun = `DELETE FROM znakowanie_studio WHERE identyfikator_zewnetrzny = ?`

	znakowanieStudiaPrzestawStan = `UPDATE znakowanie_studio
	                                SET stan = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                WHERE identyfikator_zewnetrzny = ? AND stan = 'open'`

	// Liczba użyć liczy się w tym samym zapytaniu co wykaz rodzajów: pytanie
	// „ile znakowań tego rodzaju stoi w dokumencie" pada zawsze razem z wykazem,
	// a osobne zapytanie na rodzaj dałoby tyle zapytań, ile rodzajów.
	rodzajeZnacznikaLista = `SELECT r.nazwa, r.nazwa_widoczna, r.barwa, r.fabryczny,
	                                COUNT(z.id)
	                         FROM rodzaj_znacznika_studio r
	                         LEFT JOIN znakowanie_studio z
	                              ON z.znacznik_nazwa = r.nazwa
	                             AND (? = 0 OR z.dokument_id = ?)
	                         GROUP BY r.id ORDER BY r.fabryczny DESC, r.nazwa`

	rodzajZnacznikaPobierz = `SELECT nazwa, nazwa_widoczna, barwa, fabryczny, 0
	                          FROM rodzaj_znacznika_studio WHERE nazwa = ?`

	rodzajZnacznikaZapisz = `INSERT INTO rodzaj_znacznika_studio (nazwa, nazwa_widoczna, barwa)
	                         VALUES (?, ?, ?)
	                         ON CONFLICT(nazwa) DO UPDATE SET
	                             nazwa_widoczna = excluded.nazwa_widoczna,
	                             barwa = excluded.barwa,
	                             zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	// Warunek `fabryczny = 0` w poleceniu, nie w rdzeniu, byłby cichą odmową:
	// rdzeń ma powiedzieć, DLACZEGO nie usunął. Dlatego polecenie usuwa bez
	// warunku, a fabryczność sprawdza wołający i nazywa powód.
	rodzajZnacznikaUsun = `DELETE FROM rodzaj_znacznika_studio WHERE nazwa = ?`

	kolumnyZajeciaStudia = `j.id, j.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                        j.wykonawca_rodzaj, j.wykonawca_agent_kod, j.wykonawca_agent_nazwa,
	                        j.wykonawca_podagent_kod, j.stan, j.zakres_od, j.zakres_do,
	                        j.zadanie_kod, j.zajeto, j.wygasa,
	                        CASE WHEN j.wygasa IS NOT NULL
	                                  AND j.wygasa < strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             THEN 1 ELSE 0 END`

	zrodloZajeciaStudia = ` FROM zajecie_fragmentu_studio j
	                        JOIN dokument_studio d ON d.id = j.dokument_id`

	// Czas wygaśnięcia liczy baza: `?` niesie modyfikator w postaci „+90 seconds".
	zajecieStudiaZapisz = `INSERT INTO zajecie_fragmentu_studio
	                       (identyfikator_zewnetrzny, dokument_id, wykonawca_rodzaj,
	                        wykonawca_agent_kod, wykonawca_agent_nazwa,
	                        wykonawca_podagent_kod, stan, zakres_od, zakres_do,
	                        zadanie_kod, wygasa)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	                               strftime('%Y-%m-%dT%H:%M:%fZ','now', ?))`

	zajecieStudiaPobierz = `SELECT ` + kolumnyZajeciaStudia + zrodloZajeciaStudia +
		` WHERE j.identyfikator_zewnetrzny = ?`

	zajeciaStudiaLista = `SELECT ` + kolumnyZajeciaStudia + zrodloZajeciaStudia +
		` WHERE j.dokument_id = ? ORDER BY j.zakres_od, j.id`

	zajecieStudiaZwolnij = `DELETE FROM zajecie_fragmentu_studio
	                        WHERE identyfikator_zewnetrzny = ?`

	// Zwolnienie wszystkich zajęć wykonawcy oraz przemiecenie wygasłych jednym
	// poleceniem, bo obie czynności usuwają wiersze tej samej tabeli po tym
	// samym kluczu dokumentu. Pusty kod agenta znaczy „nie zawężaj po agencie".
	zajeciaStudiaZwolnijWykonawcy = `DELETE FROM zajecie_fragmentu_studio
	                                 WHERE dokument_id = ?
	                                   AND (? = '' OR wykonawca_agent_kod = ?)
	                                   AND (? = '' OR wykonawca_podagent_kod = ?)`

	zajeciaStudiaPrzemiec = `DELETE FROM zajecie_fragmentu_studio
	                         WHERE dokument_id = ? AND wygasa IS NOT NULL
	                           AND wygasa < strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	kolumnySpieciaStudia = `s.id, s.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                        s.zakres_od, s.zakres_do, s.weszla_rodzaj, s.weszla_agent_kod,
	                        s.weszla_agent_nazwa, s.odlozona_rodzaj, s.odlozona_agent_kod,
	                        s.odlozona_agent_nazwa, s.nastawa, s.powod, s.brzmienie_odlozone,
	                        s.zmiana_sledzona_kod, s.znakowanie_kod, s.domkniete, s.utworzono`

	zrodloSpieciaStudia = ` FROM spiecie_wykonawcow_studio s
	                        JOIN dokument_studio d ON d.id = s.dokument_id`

	spiecieStudiaZapisz = `INSERT INTO spiecie_wykonawcow_studio
	                       (identyfikator_zewnetrzny, dokument_id, zakres_od, zakres_do,
	                        weszla_rodzaj, weszla_agent_kod, weszla_agent_nazwa,
	                        odlozona_rodzaj, odlozona_agent_kod, odlozona_agent_nazwa,
	                        nastawa, powod, brzmienie_odlozone, zmiana_sledzona_kod,
	                        znakowanie_kod, domkniete)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	spieciaStudiaLista = `SELECT ` + kolumnySpieciaStudia + zrodloSpieciaStudia +
		` WHERE s.dokument_id = ? ORDER BY s.utworzono DESC, s.id DESC`

	kolumnyWersjiSzeregu = `w.id, w.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                        w.etykieta, w.podsumowanie, w.skrot_tresci, w.tresc,
	                        w.postac_json, w.autor, w.autor_agent_kod, w.autor_agent_nazwa,
	                        w.kamien_milowy, w.galaz_id, w.propozycja_id, w.szereg, w.utworzono`

	zrodloWersjiSzeregu = ` FROM wersja_dokumentu_studio w
	                        JOIN dokument_studio d ON d.id = w.dokument_id`

	wersjeSzereguLista = `SELECT ` + kolumnyWersjiSzeregu + zrodloWersjiSzeregu +
		` WHERE w.dokument_id = ? ORDER BY w.utworzono DESC, w.id DESC`

	wersjaSzereguPobierz = `SELECT ` + kolumnyWersjiSzeregu + zrodloWersjiSzeregu +
		` WHERE w.identyfikator_zewnetrzny = ?`

	// Wersja założycielska to NAJSTARSZY wiersz szeregu Operatora. Szereg
	// autozapisu odpada z rachunku z zamysłem: gdyby pierwszy zapis samoczynny
	// wypadł przed pierwszym zapisem Operatora, „powrót do stanu pierwotnego"
	// wracałby do stanu przypadkowego, a nie do tego, co Operator założył.
	wersjaZalozycielskaStudia = `SELECT ` + kolumnyWersjiSzeregu + zrodloWersjiSzeregu +
		` WHERE w.dokument_id = ? AND w.szereg = 'operator'
		  ORDER BY w.utworzono ASC, w.id ASC LIMIT 1`

	wersjaSzereguZapisz = `INSERT INTO wersja_dokumentu_studio
	                       (identyfikator_zewnetrzny, dokument_id, etykieta, podsumowanie,
	                        skrot_tresci, tresc, postac_json, autor, autor_agent_kod,
	                        autor_agent_nazwa, kamien_milowy, galaz_id, propozycja_id, szereg)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Przemiatanie szeregu autozapisu. Szeregu Operatora nie tyka: wersje nazwane
	// i kluczowe zakłada Operator i historia nic nie usuwa bez jego decyzji.
	wersjeAutozapisuPrzemiec = `DELETE FROM wersja_dokumentu_studio
	                            WHERE dokument_id = ? AND szereg = 'autosave'
	                              AND kamien_milowy = 0 AND etykieta IS NULL
	                              AND id NOT IN (SELECT id FROM wersja_dokumentu_studio
	                                             WHERE dokument_id = ? AND szereg = 'autosave'
	                                             ORDER BY utworzono DESC, id DESC LIMIT ?)`
)

// ── Znakowanie fragmentów ───────────────────────────────────────────────────

// ZapiszZnakowanie zakłada znakowanie fragmentu.
func (r *repozytoriumStudia) ZapiszZnakowanie(ctx context.Context, dokumentID int64,
	znakowanie ZnakowanieStudia) (ZnakowanieStudia, error) {

	if znakowanie.Kod == "" {
		return ZnakowanieStudia{}, fmt.Errorf("dane: znakowanie studio bez identyfikatora")
	}
	if znakowanie.AutorRodzaj == "" {
		znakowanie.AutorRodzaj = "uzytkownik"
	}
	if znakowanie.Stan == "" {
		znakowanie.Stan = "open"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, znakowanieStudiaZapisz)
	if err != nil {
		return ZnakowanieStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, znakowanie.Kod, dokumentID, znakowanie.Rodzaj,
		znakowanie.AutorRodzaj, tekstDoKolumny(znakowanie.AutorAgentKod),
		tekstDoKolumny(znakowanie.AutorAgentNazwa), tekstDoKolumny(znakowanie.AutorAgentWersja),
		tekstDoKolumny(znakowanie.AutorPodagentKod), znakowanie.ZakresOd, znakowanie.ZakresDo,
		tekstDoKolumny(znakowanie.Barwa), tekstDoKolumny(znakowanie.ZnacznikNazwa),
		tekstDoKolumny(znakowanie.Tresc), tekstDoKolumny(znakowanie.BrzmienieProponowane),
		tekstDoKolumny(znakowanie.BrzmienieZastane), znakowanie.Stan,
		tekstDoKolumny(znakowanie.CzynnoscKod))
	if err != nil {
		return ZnakowanieStudia{}, fmt.Errorf("dane: nie można zapisać znakowania %q: %w",
			znakowanie.Kod, err)
	}
	return r.Znakowanie(ctx, znakowanie.Kod)
}

// Znakowanie oddaje znakowanie o wskazanym kodzie.
func (r *repozytoriumStudia) Znakowanie(ctx context.Context, kod string) (ZnakowanieStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, znakowanieStudiaPobierz)
	if err != nil {
		return ZnakowanieStudia{}, err
	}
	znakowanie, err := odczytajZnakowanieStudia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ZnakowanieStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return ZnakowanieStudia{}, fmt.Errorf("dane: nieczytelny wiersz znakowania %q: %w", kod, err)
	}
	return znakowanie, nil
}

// Znakowania oddaje znakowania dokumentu w kolejności wystąpienia w treści.
func (r *repozytoriumStudia) Znakowania(ctx context.Context,
	dokumentID int64) ([]ZnakowanieStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, znakowaniaStudiaLista)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać znakowań: %w", err)
	}
	defer wiersze.Close()

	lista := []ZnakowanieStudia{}
	for wiersze.Next() {
		znakowanie, err := odczytajZnakowanieStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz znakowania: %w", err)
		}
		lista = append(lista, znakowanie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt znakowań: %w", err)
	}
	return lista, nil
}

// UsunZnakowanie zdejmuje znakowanie i mówi, czy było.
func (r *repozytoriumStudia) UsunZnakowanie(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, znakowanieStudiaUsun)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zdjąć znakowania %q: %w", kod, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek zdjęcia znakowania %q: %w", kod, err)
	}
	return zdjete > 0, nil
}

// PrzestawStanZnakowania rozstrzyga znakowanie i mówi, czy decyzja zapadła.
// Warunek „stan = open" jest w poleceniu: decyzja raz podjęta nie zmienia się
// drugim wywołaniem.
func (r *repozytoriumStudia) PrzestawStanZnakowania(ctx context.Context,
	kod, stan string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, znakowanieStudiaPrzestawStan)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, stan, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można rozstrzygnąć znakowania %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek rozstrzygnięcia znakowania %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// RodzajeZnacznika oddaje rodzaje znaczników wraz z liczbą użyć. Dokument zerowy
// znaczy „nie licz użyć w żadnym dokumencie" — wtedy licznik jest liczbą
// znakowań wszystkich dokumentów.
func (r *repozytoriumStudia) RodzajeZnacznika(ctx context.Context,
	dokumentID int64) ([]RodzajZnacznikaStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, rodzajeZnacznikaLista)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać rodzajów znacznika: %w", err)
	}
	defer wiersze.Close()

	lista := []RodzajZnacznikaStudia{}
	for wiersze.Next() {
		rodzaj, err := odczytajRodzajZnacznikaStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz rodzaju znacznika: %w", err)
		}
		lista = append(lista, rodzaj)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt rodzajów znacznika: %w", err)
	}
	return lista, nil
}

// RodzajZnacznika oddaje jeden rodzaj znacznika, bez licznika użyć.
func (r *repozytoriumStudia) RodzajZnacznika(ctx context.Context,
	nazwa string) (RodzajZnacznikaStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, rodzajZnacznikaPobierz)
	if err != nil {
		return RodzajZnacznikaStudia{}, err
	}
	rodzaj, err := odczytajRodzajZnacznikaStudia(polecenie.QueryRowContext(ctx, nazwa))
	if errors.Is(err, sql.ErrNoRows) {
		return RodzajZnacznikaStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return RodzajZnacznikaStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz rodzaju znacznika %q: %w", nazwa, err)
	}
	return rodzaj, nil
}

// ZapiszRodzajZnacznika zakłada rodzaj znacznika własnego albo zmienia zastany.
func (r *repozytoriumStudia) ZapiszRodzajZnacznika(ctx context.Context,
	rodzaj RodzajZnacznikaStudia) (RodzajZnacznikaStudia, error) {

	if rodzaj.Nazwa == "" {
		return RodzajZnacznikaStudia{}, fmt.Errorf("dane: rodzaj znacznika bez nazwy")
	}
	if rodzaj.NazwaWidoczna == "" {
		rodzaj.NazwaWidoczna = rodzaj.Nazwa
	}
	polecenie, err := r.zapytania.przygotuj(ctx, rodzajZnacznikaZapisz)
	if err != nil {
		return RodzajZnacznikaStudia{}, err
	}
	if _, err := polecenie.ExecContext(ctx, rodzaj.Nazwa, rodzaj.NazwaWidoczna,
		tekstDoKolumny(rodzaj.Barwa)); err != nil {

		return RodzajZnacznikaStudia{}, fmt.Errorf(
			"dane: nie można zapisać rodzaju znacznika %q: %w", rodzaj.Nazwa, err)
	}
	return r.RodzajZnacznika(ctx, rodzaj.Nazwa)
}

// UsunRodzajZnacznika usuwa rodzaj znacznika i mówi, czy był. Fabryczności NIE
// sprawdza: sprawdzenie należy do rdzenia, bo rdzeń ma nazwać powód odmowy.
func (r *repozytoriumStudia) UsunRodzajZnacznika(ctx context.Context, nazwa string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, rodzajZnacznikaUsun)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, nazwa)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć rodzaju znacznika %q: %w", nazwa, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia rodzaju %q: %w", nazwa, err)
	}
	return usuniete > 0, nil
}

// ── Zajęcia fragmentów i spięcia wykonawców ─────────────────────────────────

// ZapiszZajecieFragmentu zajmuje fragment dla wykonawcy na wskazaną liczbę
// sekund. Odstęp zerowy albo ujemny znaczy zajęcie bez wygasania — i takiego
// zajęcia rdzeń nie zakłada, bo agent ubity w pół pracy trzymałby fragment na
// zawsze; granicę podaje wołający.
func (r *repozytoriumStudia) ZapiszZajecieFragmentu(ctx context.Context, dokumentID int64,
	zajecie ZajecieFragmentuStudia, waznoscSekund int64) (ZajecieFragmentuStudia, error) {

	if zajecie.Kod == "" {
		return ZajecieFragmentuStudia{}, fmt.Errorf("dane: zajęcie fragmentu bez identyfikatora")
	}
	if zajecie.WykonawcaRodzaj == "" {
		zajecie.WykonawcaRodzaj = "model"
	}
	if zajecie.Stan == "" {
		zajecie.Stan = "working"
	}
	if waznoscSekund <= 0 {
		waznoscSekund = 1
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zajecieStudiaZapisz)
	if err != nil {
		return ZajecieFragmentuStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, zajecie.Kod, dokumentID, zajecie.WykonawcaRodzaj,
		tekstDoKolumny(zajecie.AgentKod), tekstDoKolumny(zajecie.AgentNazwa),
		tekstDoKolumny(zajecie.PodagentKod), zajecie.Stan, zajecie.ZakresOd, zajecie.ZakresDo,
		tekstDoKolumny(zajecie.ZadanieKod), fmt.Sprintf("+%d seconds", waznoscSekund))
	if err != nil {
		return ZajecieFragmentuStudia{}, fmt.Errorf(
			"dane: nie można zapisać zajęcia fragmentu %q: %w", zajecie.Kod, err)
	}
	polecenieOdczytu, err := r.zapytania.przygotuj(ctx, zajecieStudiaPobierz)
	if err != nil {
		return ZajecieFragmentuStudia{}, err
	}
	zapisane, err := odczytajZajecieFragmentuStudia(
		polecenieOdczytu.QueryRowContext(ctx, zajecie.Kod))
	if err != nil {
		return ZajecieFragmentuStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz zajęcia fragmentu %q: %w", zajecie.Kod, err)
	}
	return zapisane, nil
}

// ZajeciaFragmentow oddaje zajęcia dokumentu wraz z rozstrzygnięciem, które
// wygasły. Przemiatania NIE robi sama: wykaz ma pokazać także wygasłe, gdy
// Operator o to poprosi, a przemiatanie jest osobną czynnością.
func (r *repozytoriumStudia) ZajeciaFragmentow(ctx context.Context,
	dokumentID int64) ([]ZajecieFragmentuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zajeciaStudiaLista)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zajęć fragmentów: %w", err)
	}
	defer wiersze.Close()

	lista := []ZajecieFragmentuStudia{}
	for wiersze.Next() {
		zajecie, err := odczytajZajecieFragmentuStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zajęcia fragmentu: %w", err)
		}
		lista = append(lista, zajecie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zajęć fragmentów: %w", err)
	}
	return lista, nil
}

// ZwolnijZajecieFragmentu zwalnia jedno zajęcie po kodzie.
func (r *repozytoriumStudia) ZwolnijZajecieFragmentu(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zajecieStudiaZwolnij)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zwolnić zajęcia fragmentu %q: %w", kod, err)
	}
	zwolnione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek zwolnienia zajęcia %q: %w", kod, err)
	}
	return zwolnione > 0, nil
}

// ZwolnijZajeciaWykonawcy zwalnia wszystkie zajęcia wskazanego wykonawcy
// w dokumencie i oddaje, ile ich zwolniono.
func (r *repozytoriumStudia) ZwolnijZajeciaWykonawcy(ctx context.Context, dokumentID int64,
	agentKod, podagentKod string) (int64, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zajeciaStudiaZwolnijWykonawcy)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, dokumentID, agentKod, agentKod,
		podagentKod, podagentKod)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zwolnić zajęć wykonawcy w dokumencie %d: %w",
			dokumentID, err)
	}
	zwolnione, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany skutek zwolnienia zajęć w dokumencie %d: %w",
			dokumentID, err)
	}
	return zwolnione, nil
}

// PrzemiecZajeciaFragmentow zdejmuje zajęcia wygasłe. Bez tego agent ubity
// w pół pracy trzymałby fragment na zawsze.
func (r *repozytoriumStudia) PrzemiecZajeciaFragmentow(ctx context.Context,
	dokumentID int64) (int64, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zajeciaStudiaPrzemiec)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, dokumentID)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można przemieść zajęć dokumentu %d: %w", dokumentID, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany skutek przemiatania zajęć dokumentu %d: %w",
			dokumentID, err)
	}
	return zdjete, nil
}

// ZapiszSpiecieWykonawcow odkłada zapis spięcia o ten sam fragment.
func (r *repozytoriumStudia) ZapiszSpiecieWykonawcow(ctx context.Context, dokumentID int64,
	spiecie SpiecieWykonawcowStudia) error {

	if spiecie.Kod == "" {
		return fmt.Errorf("dane: spięcie wykonawców bez identyfikatora")
	}
	if spiecie.WeszlaRodzaj == "" {
		spiecie.WeszlaRodzaj = "model"
	}
	if spiecie.OdlozonaRodzaj == "" {
		spiecie.OdlozonaRodzaj = "model"
	}
	if spiecie.Nastawa == "" {
		spiecie.Nastawa = "queue"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, spiecieStudiaZapisz)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, spiecie.Kod, dokumentID, spiecie.ZakresOd,
		spiecie.ZakresDo, spiecie.WeszlaRodzaj, tekstDoKolumny(spiecie.WeszlaAgentKod),
		tekstDoKolumny(spiecie.WeszlaAgentNazwa), spiecie.OdlozonaRodzaj,
		tekstDoKolumny(spiecie.OdlozonaAgentKod), tekstDoKolumny(spiecie.OdlozonaAgentNaz),
		spiecie.Nastawa, spiecie.Powod, tekstDoKolumny(spiecie.BrzmienieOdlozone),
		tekstDoKolumny(spiecie.ZmianaSledzonaK), tekstDoKolumny(spiecie.ZnakowanieKod),
		liczbaLogiczna(spiecie.Domkniete))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać spięcia wykonawców %q: %w", spiecie.Kod, err)
	}
	return nil
}

// SpieciaWykonawcow oddaje spięcia dokumentu, od najświeższego.
func (r *repozytoriumStudia) SpieciaWykonawcow(ctx context.Context,
	dokumentID int64) ([]SpiecieWykonawcowStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, spieciaStudiaLista)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać spięć wykonawców: %w", err)
	}
	defer wiersze.Close()

	lista := []SpiecieWykonawcowStudia{}
	for wiersze.Next() {
		var spiecie SpiecieWykonawcowStudia
		var weszlaKod, weszlaNazwa, odlozonaKod, odlozonaNazwa sql.NullString
		var brzmienie, zmiana, znakowanie sql.NullString
		var domkniete int
		err := wiersze.Scan(&spiecie.ID, &spiecie.Kod, &spiecie.DokumentKod,
			&spiecie.ZakresOd, &spiecie.ZakresDo, &spiecie.WeszlaRodzaj, &weszlaKod,
			&weszlaNazwa, &spiecie.OdlozonaRodzaj, &odlozonaKod, &odlozonaNazwa,
			&spiecie.Nastawa, &spiecie.Powod, &brzmienie, &zmiana, &znakowanie,
			&domkniete, &spiecie.Utworzono)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz spięcia wykonawców: %w", err)
		}
		spiecie.WeszlaAgentKod = tekstZKolumny(weszlaKod)
		spiecie.WeszlaAgentNazwa = tekstZKolumny(weszlaNazwa)
		spiecie.OdlozonaAgentKod = tekstZKolumny(odlozonaKod)
		spiecie.OdlozonaAgentNaz = tekstZKolumny(odlozonaNazwa)
		spiecie.BrzmienieOdlozone = tekstZKolumny(brzmienie)
		spiecie.ZmianaSledzonaK = tekstZKolumny(zmiana)
		spiecie.ZnakowanieKod = tekstZKolumny(znakowanie)
		spiecie.Domkniete = domkniete == 1
		lista = append(lista, spiecie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt spięć wykonawców: %w", err)
	}
	return lista, nil
}

// ── Wersje w szeregach ──────────────────────────────────────────────────────

// ZapiszWersjeSzeregu zakłada wersję we wskazanym szeregu, wraz z postacią
// dokumentu i tożsamością wykonawcy. Wskaźnika wersji bieżącej NIE przestawia —
// tak samo jak `ZapiszWersje` z `studio_wersje.go`; decyzja, kiedy nowa wersja
// staje się bieżącą, należy do wołającego.
func (r *repozytoriumStudia) ZapiszWersjeSzeregu(ctx context.Context, dokumentID int64,
	wersja WersjaSzereguStudia) (WersjaSzereguStudia, error) {

	if wersja.Kod == "" {
		return WersjaSzereguStudia{}, fmt.Errorf("dane: wersja dokumentu bez identyfikatora")
	}
	if wersja.Szereg == "" {
		wersja.Szereg = "operator"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wersjaSzereguZapisz)
	if err != nil {
		return WersjaSzereguStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, wersja.Kod, dokumentID,
		tekstDoKolumny(wersja.Etykieta), tekstDoKolumny(wersja.Podsumowanie),
		tekstDoKolumny(wersja.SkrotTresci), tekstDoKolumny(wersja.Tresc),
		tekstDoKolumny(wersja.PostacJSON), tekstDoKolumny(wersja.Autor),
		tekstDoKolumny(wersja.AutorAgentKod), tekstDoKolumny(wersja.AutorAgentNazwa),
		liczbaLogiczna(wersja.KamienMilowy), tekstDoKolumny(wersja.GalazKod),
		tekstDoKolumny(wersja.PropozycjaKod), wersja.Szereg)
	if err != nil {
		return WersjaSzereguStudia{}, fmt.Errorf("dane: nie można zapisać wersji %q: %w",
			wersja.Kod, err)
	}
	return r.WersjaSzeregu(ctx, wersja.Kod)
}

// WersjaSzeregu oddaje wersję wraz z szeregiem i postacią.
func (r *repozytoriumStudia) WersjaSzeregu(ctx context.Context,
	kod string) (WersjaSzereguStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wersjaSzereguPobierz)
	if err != nil {
		return WersjaSzereguStudia{}, err
	}
	wersja, err := odczytajWersjeSzereguStudia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaSzereguStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaSzereguStudia{}, fmt.Errorf("dane: nieczytelny wiersz wersji %q: %w", kod, err)
	}
	return wersja, nil
}

// WersjeSzeregow oddaje wszystkie wersje dokumentu, od najświeższej, wraz
// z szeregiem — rozbicie na szeregi robi wołający, bo licznik obu szeregów
// oddaje jedna odpowiedź.
func (r *repozytoriumStudia) WersjeSzeregow(ctx context.Context,
	dokumentID int64) ([]WersjaSzereguStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wersjeSzereguLista)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wersji dokumentu: %w", err)
	}
	defer wiersze.Close()

	lista := []WersjaSzereguStudia{}
	for wiersze.Next() {
		wersja, err := odczytajWersjeSzereguStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji dokumentu: %w", err)
		}
		lista = append(lista, wersja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wersji dokumentu: %w", err)
	}
	return lista, nil
}

// WersjaZalozycielska oddaje najstarszą wersję szeregu Operatora — stan
// pierwotny, do którego wraca się jednym poleceniem, bez szukania w wykazie.
func (r *repozytoriumStudia) WersjaZalozycielska(ctx context.Context,
	dokumentID int64) (WersjaSzereguStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wersjaZalozycielskaStudia)
	if err != nil {
		return WersjaSzereguStudia{}, err
	}
	wersja, err := odczytajWersjeSzereguStudia(polecenie.QueryRowContext(ctx, dokumentID))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaSzereguStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaSzereguStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz wersji założycielskiej dokumentu %d: %w", dokumentID, err)
	}
	return wersja, nil
}

// PrzemiecWersjeAutozapisu zdejmuje nadmiar wersji szeregu autozapisu.
// Szeregu Operatora nie tyka, wersji nazwanych i kluczowych nie tyka nigdzie:
// historia nic nie usuwa bez decyzji Operatora.
func (r *repozytoriumStudia) PrzemiecWersjeAutozapisu(ctx context.Context, dokumentID int64,
	ileZachowac int64) (int64, error) {

	if ileZachowac <= 0 {
		ileZachowac = 1
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wersjeAutozapisuPrzemiec)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, dokumentID, dokumentID, ileZachowac)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można przemieść wersji autozapisu dokumentu %d: %w",
			dokumentID, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany skutek przemiatania wersji autozapisu %d: %w",
			dokumentID, err)
	}
	return zdjete, nil
}

// ── Składanie wierszy ───────────────────────────────────────────────────────

func odczytajZnakowanieStudia(wiersz skaner) (ZnakowanieStudia, error) {
	var znakowanie ZnakowanieStudia
	var agentKod, agentNazwa, agentWersja, podagent sql.NullString
	var barwa, znacznik, tresc, proponowane, zastane, czynnosc sql.NullString
	err := wiersz.Scan(&znakowanie.ID, &znakowanie.Kod, &znakowanie.DokumentKod,
		&znakowanie.Rodzaj, &znakowanie.AutorRodzaj, &agentKod, &agentNazwa, &agentWersja,
		&podagent, &znakowanie.ZakresOd, &znakowanie.ZakresDo, &barwa, &znacznik,
		&tresc, &proponowane, &zastane, &znakowanie.Stan, &czynnosc, &znakowanie.Utworzono)
	if err != nil {
		return ZnakowanieStudia{}, err
	}
	znakowanie.AutorAgentKod = tekstZKolumny(agentKod)
	znakowanie.AutorAgentNazwa = tekstZKolumny(agentNazwa)
	znakowanie.AutorAgentWersja = tekstZKolumny(agentWersja)
	znakowanie.AutorPodagentKod = tekstZKolumny(podagent)
	znakowanie.Barwa = tekstZKolumny(barwa)
	znakowanie.ZnacznikNazwa = tekstZKolumny(znacznik)
	znakowanie.Tresc = tekstZKolumny(tresc)
	znakowanie.BrzmienieProponowane = tekstZKolumny(proponowane)
	znakowanie.BrzmienieZastane = tekstZKolumny(zastane)
	znakowanie.CzynnoscKod = tekstZKolumny(czynnosc)
	return znakowanie, nil
}

func odczytajRodzajZnacznikaStudia(wiersz skaner) (RodzajZnacznikaStudia, error) {
	var rodzaj RodzajZnacznikaStudia
	var barwa sql.NullString
	var fabryczny int
	if err := wiersz.Scan(&rodzaj.Nazwa, &rodzaj.NazwaWidoczna, &barwa, &fabryczny,
		&rodzaj.IleUzyc); err != nil {

		return RodzajZnacznikaStudia{}, err
	}
	rodzaj.Barwa = tekstZKolumny(barwa)
	rodzaj.Fabryczny = fabryczny == 1
	return rodzaj, nil
}

func odczytajZajecieFragmentuStudia(wiersz skaner) (ZajecieFragmentuStudia, error) {
	var zajecie ZajecieFragmentuStudia
	var agentKod, agentNazwa, podagent, zadanie, wygasa sql.NullString
	var wygasle int
	err := wiersz.Scan(&zajecie.ID, &zajecie.Kod, &zajecie.DokumentKod,
		&zajecie.WykonawcaRodzaj, &agentKod, &agentNazwa, &podagent, &zajecie.Stan,
		&zajecie.ZakresOd, &zajecie.ZakresDo, &zadanie, &zajecie.Zajeto, &wygasa, &wygasle)
	if err != nil {
		return ZajecieFragmentuStudia{}, err
	}
	zajecie.AgentKod = tekstZKolumny(agentKod)
	zajecie.AgentNazwa = tekstZKolumny(agentNazwa)
	zajecie.PodagentKod = tekstZKolumny(podagent)
	zajecie.ZadanieKod = tekstZKolumny(zadanie)
	zajecie.Wygasa = tekstZKolumny(wygasa)
	zajecie.Wygasle = wygasle == 1
	return zajecie, nil
}

func odczytajWersjeSzereguStudia(wiersz skaner) (WersjaSzereguStudia, error) {
	var wersja WersjaSzereguStudia
	var etykieta, podsumowanie, skrot, tresc, postac sql.NullString
	var autor, agentKod, agentNazwa, galaz, propozycja sql.NullString
	var kamien int
	err := wiersz.Scan(&wersja.ID, &wersja.Kod, &wersja.DokumentKod, &etykieta,
		&podsumowanie, &skrot, &tresc, &postac, &autor, &agentKod, &agentNazwa,
		&kamien, &galaz, &propozycja, &wersja.Szereg, &wersja.Utworzono)
	if err != nil {
		return WersjaSzereguStudia{}, err
	}
	wersja.Etykieta = tekstZKolumny(etykieta)
	wersja.Podsumowanie = tekstZKolumny(podsumowanie)
	wersja.SkrotTresci = tekstZKolumny(skrot)
	wersja.Tresc = tekstZKolumny(tresc)
	wersja.PostacJSON = tekstZKolumny(postac)
	wersja.Autor = tekstZKolumny(autor)
	wersja.AutorAgentKod = tekstZKolumny(agentKod)
	wersja.AutorAgentNazwa = tekstZKolumny(agentNazwa)
	wersja.GalazKod = tekstZKolumny(galaz)
	wersja.PropozycjaKod = tekstZKolumny(propozycja)
	wersja.KamienMilowy = kamien == 1
	return wersja, nil
}
