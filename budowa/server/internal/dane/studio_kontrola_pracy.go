// Odpowiedzialność pliku: warstwa danych odcinka kontroli pracy modułu Studio —
// blokady fragmentów (migracja 363), odwracalny dziennik czynności wraz
// z zależnościami (364) oraz kopie zapasowe i nastawy pracy okna (366).
//
// ── Dlaczego te trzy obszary leżą w jednym pliku ─────────────────────────────
// Wiąże je jedno pytanie, zadawane w jednym miejscu: „czy tę zmianę wolno
// wnieść, a jeśli tak, to czym ją potem cofnąć". Blokada odpowiada na pierwszą
// połowę, dziennik na drugą, a kopia zapasowa jest siatką pod obiema — zakłada
// się ją PRZED czynnością nieodwracalną i po nieudanym zapisie. Rozdzielenie
// ich na trzy pliki rozdzieliłoby zapytania, które i tak padają razem.
//
// Znakowanie fragmentów, zajęcia wykonawców i spięcia leżą osobno
// (`studio_znakowanie_wykonawcy.go`): tamte opisują, co ktoś o dokumencie
// POWIEDZIAŁ i kto nad nim PRACUJE, a nie czego nie wolno tknąć.
//
// ── Dlaczego czas mierzy baza, nie rdzeń ─────────────────────────────────────
// Wygasanie zajęcia i wygasanie kopii zapasowej liczy się w SQL
// (`strftime('now')`), bo obie wielkości muszą być tym samym zegarem co kolumny
// `utworzono` i `wygasa` w tych wierszach. Zegar rdzenia i zegar bazy rozjadą
// się przy pierwszej różnicy strefy, a wtedy kopia wygasłaby wcześniej albo
// później, niż mówi nastawa Operatora.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// BlokadaFragmentuStudia to wiersz tabeli `blokada_fragmentu_studio`.
//
// `Zasieg` rozstrzyga, kogo blokada dotyczy: `model` (postać domyślna) wiąże
// wyłącznie wykonawców, `everyone` wiąże także Operatora. Operator zmienia
// fragment zablokowany bez przeszkód, dopóki nie zażąda tego drugiego jawnie.
type BlokadaFragmentuStudia struct {
	ID              int64
	Kod             string
	DokumentKod     string
	Nazwa           string
	Powod           *string
	ZakresOd        int64
	ZakresDo        int64
	Zasieg          string
	ZalozylRodzaj   string
	ZalozylAgentKod *string
	ZalozylAgentNaz *string
	ZSzablonu       bool
	SzablonKod      *string
	Przesuniecia    int64
	Utworzono       string
}

// BlokadaSzablonuStudia to wiersz tabeli `blokada_szablonu_studio` — blokada
// WZORCOWA, kopiowana do każdego dokumentu zakładanego z szablonu.
type BlokadaSzablonuStudia struct {
	Kod        string
	SzablonKod string
	Nazwa      string
	Powod      *string
	ZakresOd   int64
	ZakresDo   int64
	Zasieg     string
}

// CzynnoscDokumentuStudia to wiersz tabeli `czynnosc_dokumentu_studio`.
//
// `StanPrzed` i `StanPo` niosą WYCINEK objęty czynnością, nie migawkę całego
// dokumentu — inaczej cofnięcie czynności ze środka dziennika zabrałoby ze sobą
// wszystko, co po niej weszło, czyli byłoby przywróceniem wersji.
type CzynnoscDokumentuStudia struct {
	ID               int64
	Kod              string
	DokumentKod      string
	Kolejnosc        int64
	Rodzaj           string
	AutorRodzaj      string
	AutorAgentKod    *string
	AutorAgentNazwa  *string
	AutorAgentWersja *string
	AutorPodagentKod *string
	Opis             string
	ZakresOd         *int64
	ZakresDo         *int64
	StanPrzed        *string
	StanPo           *string
	ZmianaSledzonaK  *string
	ZadanieKod       *string
	Stan             string
	Utworzono        string
	// PodstawyKody i StojaceNaNiej wypełnia odczyt z tabeli zależności. Puste
	// znaczy „czynność samodzielna", nie „nie wiadomo".
	PodstawyKody  []string
	StojaceNaNiej []string
}

// KopiaZapasowaStudia to wiersz tabeli `kopia_zapasowa_studio`.
//
// `UdaloSie` i `PowodNiepowodzenia` istnieją, bo wskaźnik „zapisano" pokazany
// przy zapisie nieudanym jest najgorszym możliwym błędem tego modułu: Operator
// zamknie okno i straci pracę. Kopia nieudana ZOSTAJE wierszem — musi być
// widoczna i nazwana, a nie zniknąć razem z niepowodzeniem.
type KopiaZapasowaStudia struct {
	ID                 int64
	Kod                string
	DokumentKod        string
	Powod              string
	Tresc              *string
	TrescOdwolanie     *string
	PostacJSON         *string
	RozmiarBajtow      int64
	UdaloSie           bool
	PowodNiepowodzenia *string
	ZmianyNiezapisane  bool
	Utworzono          string
}

// NastawaPracyStudia to wiersz tabeli `nastawa_pracy_studio` — autozapis,
// wygasanie kopii i nastawy widoku, pamiętane przy dokumencie albo przy oknie.
type NastawaPracyStudia struct {
	ID                       int64
	Okno                     string
	DokumentID               *int64
	AutozapisCzynny          bool
	AutozapisOdstepSekund    int64
	AutozapisPrzyOdejsciu    bool
	AutozapisPrzyZamknieciu  bool
	AutozapisPrzyPrzelacz    bool
	KopieIleZachowac         int64
	KopieWygasanieGodzin     int64
	OstatniZapis             *string
	OstatniZapisNieudany     bool
	OstatniPowodNiepowodz    *string
	PodswietlenieZmianModelu bool
}

const (
	kolumnyBlokadyStudia = `b.id, b.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                        b.nazwa, b.powod, b.zakres_od, b.zakres_do, b.zasieg,
	                        b.zalozyl_rodzaj, b.zalozyl_agent_kod, b.zalozyl_agent_nazwa,
	                        b.z_szablonu, b.szablon_kod, b.przesuniecia, b.utworzono`

	zrodloBlokadyStudia = ` FROM blokada_fragmentu_studio b
	                        JOIN dokument_studio d ON d.id = b.dokument_id`

	blokadaStudiaZapisz = `INSERT INTO blokada_fragmentu_studio
	                       (identyfikator_zewnetrzny, dokument_id, nazwa, powod,
	                        zakres_od, zakres_do, zasieg, zalozyl_rodzaj,
	                        zalozyl_agent_kod, zalozyl_agent_nazwa, z_szablonu, szablon_kod)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	blokadaStudiaPobierz = `SELECT ` + kolumnyBlokadyStudia + zrodloBlokadyStudia +
		` WHERE b.identyfikator_zewnetrzny = ?`

	// Kolejność po początku zakresu, bo tak czyta się wykaz blokad w dokumencie
	// i tak liczy się bilans pominięć: od pierwszego fragmentu do ostatniego.
	blokadyStudiaLista = `SELECT ` + kolumnyBlokadyStudia + zrodloBlokadyStudia +
		` WHERE b.dokument_id = ? ORDER BY b.zakres_od, b.id`

	blokadaStudiaUsun = `DELETE FROM blokada_fragmentu_studio
	                     WHERE identyfikator_zewnetrzny = ?`

	// Przesunięcie zakresu po wpisie PRZED blokadą. Bez tego blokada zaczęłaby
	// po pierwszej edycji chronić nie ten fragment, co miała — i to bez słowa.
	// `przesuniecia` rośnie razem z zakresem: jest miarą zaufania do zakresu,
	// widoczną dla Operatora w wykazie.
	blokadyStudiaPrzesun = `UPDATE blokada_fragmentu_studio
	                        SET zakres_od = zakres_od + ?, zakres_do = zakres_do + ?,
	                            przesuniecia = przesuniecia + 1,
	                            zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                        WHERE dokument_id = ? AND zakres_od >= ?`

	blokadySzablonuLista = `SELECT identyfikator_zewnetrzny, szablon_kod, nazwa, powod,
	                               zakres_od, zakres_do, zasieg
	                        FROM blokada_szablonu_studio
	                        WHERE szablon_kod = ? ORDER BY zakres_od, id`

	blokadaSzablonuZapisz = `INSERT INTO blokada_szablonu_studio
	                         (identyfikator_zewnetrzny, szablon_kod, nazwa, powod,
	                          zakres_od, zakres_do, zasieg)
	                         VALUES (?, ?, ?, ?, ?, ?, ?)`

	kolumnyCzynnosciStudia = `c.id, c.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                          c.kolejnosc, c.rodzaj, c.autor_rodzaj, c.autor_agent_kod,
	                          c.autor_agent_nazwa, c.autor_agent_wersja, c.autor_podagent_kod,
	                          c.opis, c.zakres_od, c.zakres_do, c.stan_przed_json,
	                          c.stan_po_json, c.zmiana_sledzona_kod, c.zadanie_kod,
	                          c.stan, c.utworzono`

	zrodloCzynnosciStudia = ` FROM czynnosc_dokumentu_studio c
	                          JOIN dokument_studio d ON d.id = c.dokument_id`

	czynnoscStudiaZapisz = `INSERT INTO czynnosc_dokumentu_studio
	                        (identyfikator_zewnetrzny, dokument_id, kolejnosc, rodzaj,
	                         autor_rodzaj, autor_agent_kod, autor_agent_nazwa,
	                         autor_agent_wersja, autor_podagent_kod, opis,
	                         zakres_od, zakres_do, stan_przed_json, stan_po_json,
	                         zmiana_sledzona_kod, zadanie_kod, stan)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	czynnoscStudiaPobierz = `SELECT ` + kolumnyCzynnosciStudia + zrodloCzynnosciStudia +
		` WHERE c.identyfikator_zewnetrzny = ?`

	// Od najświeższej: dziennik czyta się od ostatniego ruchu wstecz, tak jak
	// go czyta Operator szukający tego, co właśnie zepsuł.
	czynnosciStudiaLista = `SELECT ` + kolumnyCzynnosciStudia + zrodloCzynnosciStudia +
		` WHERE c.dokument_id = ? ORDER BY c.kolejnosc DESC, c.id DESC`

	czynnoscStudiaNastepnaKolejnosc = `SELECT COALESCE(MAX(kolejnosc), 0) + 1
	                                   FROM czynnosc_dokumentu_studio WHERE dokument_id = ?`

	czynnoscStudiaPrzestawStan = `UPDATE czynnosc_dokumentu_studio
	                              SET stan = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                              WHERE identyfikator_zewnetrzny = ? AND stan = ?`

	zaleznoscStudiaZapisz = `INSERT OR IGNORE INTO zaleznosc_czynnosci_studio
	                         (czynnosc_id, podstawa_id, powod)
	                         VALUES ((SELECT id FROM czynnosc_dokumentu_studio
	                                  WHERE identyfikator_zewnetrzny = ?),
	                                 (SELECT id FROM czynnosc_dokumentu_studio
	                                  WHERE identyfikator_zewnetrzny = ?), ?)`

	// Dwa kierunki jednej relacji. „Na czym ta czynność stoi" mówi, czego nie
	// wolno cofnąć przed nią; „co stoi na niej" mówi, dlaczego jej samej cofnąć
	// nie można — i to jest treść odmowy nazywającej zależność.
	zaleznosciStudiaPodstawy = `SELECT p.identyfikator_zewnetrzny
	                            FROM zaleznosc_czynnosci_studio z
	                            JOIN czynnosc_dokumentu_studio c ON c.id = z.czynnosc_id
	                            JOIN czynnosc_dokumentu_studio p ON p.id = z.podstawa_id
	                            WHERE c.dokument_id = ?
	                            ORDER BY c.kolejnosc, p.kolejnosc`

	zaleznosciStudiaPary = `SELECT c.identyfikator_zewnetrzny, p.identyfikator_zewnetrzny
	                        FROM zaleznosc_czynnosci_studio z
	                        JOIN czynnosc_dokumentu_studio c ON c.id = z.czynnosc_id
	                        JOIN czynnosc_dokumentu_studio p ON p.id = z.podstawa_id
	                        WHERE c.dokument_id = ?`

	kolumnyKopiiStudia = `k.id, k.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny,
	                      k.powod, k.tresc, k.tresc_odwolanie, k.postac_json,
	                      k.rozmiar_bajtow, k.udalo_sie, k.powod_niepowodzenia,
	                      k.zmiany_niezapisane, k.utworzono`

	zrodloKopiiStudia = ` FROM kopia_zapasowa_studio k
	                      JOIN dokument_studio d ON d.id = k.dokument_id`

	kopiaStudiaZapisz = `INSERT INTO kopia_zapasowa_studio
	                     (identyfikator_zewnetrzny, dokument_id, powod, tresc,
	                      tresc_odwolanie, postac_json, rozmiar_bajtow, udalo_sie,
	                      powod_niepowodzenia, zmiany_niezapisane)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	kopiaStudiaPobierz = `SELECT ` + kolumnyKopiiStudia + zrodloKopiiStudia +
		` WHERE k.identyfikator_zewnetrzny = ?`

	kopieStudiaLista = `SELECT ` + kolumnyKopiiStudia + zrodloKopiiStudia +
		` WHERE k.dokument_id = ? ORDER BY k.utworzono DESC, k.id DESC`

	// Kopie niosące zmiany niezapisane, po WSZYSTKICH dokumentach okna. Tym
	// zapytaniem Studio samo zgłasza „mam niezapisany dokument z godziny X",
	// zamiast czekać, aż Operator się domyśli.
	kopieStudiaNiezapisane = `SELECT ` + kolumnyKopiiStudia + zrodloKopiiStudia +
		` WHERE k.zmiany_niezapisane = 1 AND (? = '' OR d.okno = ?)
		  ORDER BY k.utworzono DESC, k.id DESC`

	// Przemiatanie kopii wygasłych. Zasada wygasania jest jawnym, odwracalnym
	// ustawieniem Operatora, więc granice podaje wołający, a nie stała rdzenia.
	// Kopia NIEUDANA nie wygasa razem z udanymi: jest jedynym śladem, że praca
	// nie doszła na dysk, i ma zostać, dopóki Operator jej nie zobaczy.
	kopieStudiaPrzemiec = `DELETE FROM kopia_zapasowa_studio
	                       WHERE dokument_id = ? AND udalo_sie = 1 AND zmiany_niezapisane = 0
	                         AND (utworzono < strftime('%Y-%m-%dT%H:%M:%fZ','now', ?)
	                              OR id NOT IN (SELECT id FROM kopia_zapasowa_studio
	                                            WHERE dokument_id = ?
	                                            ORDER BY utworzono DESC, id DESC LIMIT ?))`

	kolumnyNastawyPracy = `id, okno, dokument_id, autozapis_czynny, autozapis_odstep_sekund,
	                       autozapis_przy_odejsciu, autozapis_przy_zamknieciu,
	                       autozapis_przy_przelaczeniu, kopie_ile_zachowac,
	                       kopie_wygasanie_godzin, ostatni_zapis, ostatni_zapis_nieudany,
	                       ostatni_powod_niepowodzenia, podswietlenie_zmian_modelu`

	nastawaPracyPobierzDokument = `SELECT ` + kolumnyNastawyPracy + `
	                               FROM nastawa_pracy_studio
	                               WHERE okno = ? AND dokument_id = ?`

	nastawaPracyPobierzOkno = `SELECT ` + kolumnyNastawyPracy + `
	                           FROM nastawa_pracy_studio
	                           WHERE okno = ? AND dokument_id IS NULL`

	// Dwa polecenia zapisu, nie jedno z warunkiem: więz UNIQUE stoi na dwóch
	// indeksach częściowych (osobno dla wiersza dokumentu i wiersza okna),
	// a `ON CONFLICT` w SQLite wskazuje jeden zbiór kolumn.
	nastawaPracyZalozDokument = `INSERT INTO nastawa_pracy_studio (okno, dokument_id)
	                             SELECT ?, ? WHERE NOT EXISTS
	                                 (SELECT 1 FROM nastawa_pracy_studio
	                                  WHERE okno = ? AND dokument_id = ?)`

	nastawaPracyZalozOkno = `INSERT INTO nastawa_pracy_studio (okno, dokument_id)
	                         SELECT ?, NULL WHERE NOT EXISTS
	                             (SELECT 1 FROM nastawa_pracy_studio
	                              WHERE okno = ? AND dokument_id IS NULL)`

	nastawaPracyZapiszAutozapis = `UPDATE nastawa_pracy_studio SET
	                                   autozapis_czynny = ?, autozapis_odstep_sekund = ?,
	                                   autozapis_przy_odejsciu = ?, autozapis_przy_zamknieciu = ?,
	                                   autozapis_przy_przelaczeniu = ?, kopie_ile_zachowac = ?,
	                                   kopie_wygasanie_godzin = ?,
	                                   zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                               WHERE id = ?`

	// Skutek zapisu odkłada się osobnym poleceniem, bo pisze go inna czynność
	// niż nastawy: nastawy stawia Operator, skutek — sam zapis. Jedno polecenie
	// kazałoby autozapisowi przepisywać nastawy, których nie zmieniał.
	nastawaPracyZapiszSkutek = `UPDATE nastawa_pracy_studio SET
	                                ostatni_zapis = ?, ostatni_zapis_nieudany = ?,
	                                ostatni_powod_niepowodzenia = ?,
	                                zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                            WHERE id = ?`
)

// ── Blokady fragmentów ──────────────────────────────────────────────────────

// ZapiszBlokadeFragmentu zakłada blokadę fragmentu dokumentu.
func (r *repozytoriumStudia) ZapiszBlokadeFragmentu(ctx context.Context, dokumentID int64,
	blokada BlokadaFragmentuStudia) (BlokadaFragmentuStudia, error) {

	if blokada.Kod == "" {
		return BlokadaFragmentuStudia{}, fmt.Errorf("dane: blokada fragmentu bez identyfikatora")
	}
	if blokada.Zasieg == "" {
		blokada.Zasieg = "model"
	}
	if blokada.ZalozylRodzaj == "" {
		blokada.ZalozylRodzaj = "uzytkownik"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, blokadaStudiaZapisz)
	if err != nil {
		return BlokadaFragmentuStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, blokada.Kod, dokumentID, blokada.Nazwa,
		tekstDoKolumny(blokada.Powod), blokada.ZakresOd, blokada.ZakresDo, blokada.Zasieg,
		blokada.ZalozylRodzaj, tekstDoKolumny(blokada.ZalozylAgentKod),
		tekstDoKolumny(blokada.ZalozylAgentNaz), liczbaLogiczna(blokada.ZSzablonu),
		tekstDoKolumny(blokada.SzablonKod))
	if err != nil {
		return BlokadaFragmentuStudia{}, fmt.Errorf(
			"dane: nie można zapisać blokady fragmentu %q: %w", blokada.Kod, err)
	}
	return r.BlokadaFragmentu(ctx, blokada.Kod)
}

// BlokadaFragmentu oddaje blokadę o wskazanym kodzie.
func (r *repozytoriumStudia) BlokadaFragmentu(ctx context.Context,
	kod string) (BlokadaFragmentuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, blokadaStudiaPobierz)
	if err != nil {
		return BlokadaFragmentuStudia{}, err
	}
	blokada, err := odczytajBlokadeFragmentuStudia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return BlokadaFragmentuStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return BlokadaFragmentuStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz blokady fragmentu %q: %w", kod, err)
	}
	return blokada, nil
}

// BlokadyFragmentow oddaje blokady dokumentu w kolejności położenia w treści.
func (r *repozytoriumStudia) BlokadyFragmentow(ctx context.Context,
	dokumentID int64) ([]BlokadaFragmentuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, blokadyStudiaLista)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać blokad fragmentów: %w", err)
	}
	defer wiersze.Close()

	lista := []BlokadaFragmentuStudia{}
	for wiersze.Next() {
		blokada, err := odczytajBlokadeFragmentuStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz blokady fragmentu: %w", err)
		}
		lista = append(lista, blokada)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt blokad fragmentów: %w", err)
	}
	return lista, nil
}

// UsunBlokadeFragmentu zdejmuje blokadę i mówi, czy była.
func (r *repozytoriumStudia) UsunBlokadeFragmentu(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, blokadaStudiaUsun)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zdjąć blokady fragmentu %q: %w", kod, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek zdjęcia blokady %q: %w", kod, err)
	}
	return zdjete > 0, nil
}

// PrzesunBlokadyFragmentow przesuwa zakresy blokad leżących ZA punktem edycji.
//
// Blokada leżąca przed punktem edycji zostaje nietknięta; leżąca za nim jedzie
// o różnicę długości. Blokada, w której środek trafiła edycja, też zostaje
// nietknięta — bo skoro edycja weszła w blokadę, to znaczy, że wolno jej było
// tam wejść (Operator albo zasięg `model` przy czynności Operatora), a wtedy
// zakres blokady ma zostać taki, jaki Operator ustawił.
func (r *repozytoriumStudia) PrzesunBlokadyFragmentow(ctx context.Context,
	dokumentID int64, odPozycji int64, przesuniecie int64) error {

	if przesuniecie == 0 {
		return nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, blokadyStudiaPrzesun)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, przesuniecie, przesuniecie,
		dokumentID, odPozycji); err != nil {

		return fmt.Errorf("dane: nie można przeliczyć zakresów blokad dokumentu %d: %w",
			dokumentID, err)
	}
	return nil
}

// ZapiszBlokadeSzablonu zakłada blokadę wzorcową szablonu pisma.
func (r *repozytoriumStudia) ZapiszBlokadeSzablonu(ctx context.Context,
	blokada BlokadaSzablonuStudia) error {

	if blokada.Kod == "" || blokada.SzablonKod == "" {
		return fmt.Errorf("dane: blokada szablonu bez identyfikatora albo bez szablonu")
	}
	if blokada.Zasieg == "" {
		blokada.Zasieg = "model"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, blokadaSzablonuZapisz)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, blokada.Kod, blokada.SzablonKod, blokada.Nazwa,
		tekstDoKolumny(blokada.Powod), blokada.ZakresOd, blokada.ZakresDo, blokada.Zasieg)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać blokady szablonu %q: %w", blokada.Kod, err)
	}
	return nil
}

// BlokadySzablonu oddaje blokady wzorcowe szablonu.
func (r *repozytoriumStudia) BlokadySzablonu(ctx context.Context,
	szablonKod string) ([]BlokadaSzablonuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, blokadySzablonuLista)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, szablonKod)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać blokad szablonu %q: %w", szablonKod, err)
	}
	defer wiersze.Close()

	lista := []BlokadaSzablonuStudia{}
	for wiersze.Next() {
		var blokada BlokadaSzablonuStudia
		var powod sql.NullString
		if err := wiersze.Scan(&blokada.Kod, &blokada.SzablonKod, &blokada.Nazwa, &powod,
			&blokada.ZakresOd, &blokada.ZakresDo, &blokada.Zasieg); err != nil {

			return nil, fmt.Errorf("dane: nieczytelny wiersz blokady szablonu: %w", err)
		}
		blokada.Powod = tekstZKolumny(powod)
		lista = append(lista, blokada)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt blokad szablonu: %w", err)
	}
	return lista, nil
}

// ── Dziennik czynności ──────────────────────────────────────────────────────

// ZapiszCzynnoscDokumentu odkłada wpis dziennika i oddaje go wraz z nadaną
// kolejnością.
//
// Kolejność nadaje baza, nie wołający: dwie czynności zapisane w tej samej
// chwili muszą dostać różne numery, a numer nadany przez rdzeń z odczytu „ile
// jest teraz" byłby wyścigiem. Więz UNIQUE(dokument, kolejnosc) zamienia ten
// wyścig w błąd zapisu, zamiast w dwa wpisy o tym samym miejscu w porządku.
func (r *repozytoriumStudia) ZapiszCzynnoscDokumentu(ctx context.Context, dokumentID int64,
	czynnosc CzynnoscDokumentuStudia) (CzynnoscDokumentuStudia, error) {

	if czynnosc.Kod == "" {
		return CzynnoscDokumentuStudia{}, fmt.Errorf("dane: czynność dokumentu bez identyfikatora")
	}
	if czynnosc.AutorRodzaj == "" {
		czynnosc.AutorRodzaj = "uzytkownik"
	}
	if czynnosc.Stan == "" {
		czynnosc.Stan = "active"
	}
	if czynnosc.Kolejnosc == 0 {
		kolejnosc, err := r.nastepnaKolejnoscCzynnosciStudia(ctx, dokumentID)
		if err != nil {
			return CzynnoscDokumentuStudia{}, err
		}
		czynnosc.Kolejnosc = kolejnosc
	}
	polecenie, err := r.zapytania.przygotuj(ctx, czynnoscStudiaZapisz)
	if err != nil {
		return CzynnoscDokumentuStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, czynnosc.Kod, dokumentID, czynnosc.Kolejnosc,
		czynnosc.Rodzaj, czynnosc.AutorRodzaj, tekstDoKolumny(czynnosc.AutorAgentKod),
		tekstDoKolumny(czynnosc.AutorAgentNazwa), tekstDoKolumny(czynnosc.AutorAgentWersja),
		tekstDoKolumny(czynnosc.AutorPodagentKod), czynnosc.Opis,
		liczbaDoKolumny(czynnosc.ZakresOd), liczbaDoKolumny(czynnosc.ZakresDo),
		tekstDoKolumny(czynnosc.StanPrzed), tekstDoKolumny(czynnosc.StanPo),
		tekstDoKolumny(czynnosc.ZmianaSledzonaK), tekstDoKolumny(czynnosc.ZadanieKod),
		czynnosc.Stan)
	if err != nil {
		return CzynnoscDokumentuStudia{}, fmt.Errorf(
			"dane: nie można zapisać czynności dokumentu %q: %w", czynnosc.Kod, err)
	}
	// Zależności zapisujemy po wierszu czynności, bo polecenie wiąże je po
	// identyfikatorze zewnętrznym — a ten musi już w tabeli stać.
	for _, podstawa := range czynnosc.PodstawyKody {
		if err := r.ZapiszZaleznoscCzynnosci(ctx, czynnosc.Kod, podstawa, nil); err != nil {
			return CzynnoscDokumentuStudia{}, err
		}
	}
	return r.CzynnoscDokumentu(ctx, czynnosc.Kod)
}

// nastepnaKolejnoscCzynnosciStudia liczy numer porządkowy następnej czynności.
func (r *repozytoriumStudia) nastepnaKolejnoscCzynnosciStudia(ctx context.Context,
	dokumentID int64) (int64, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, czynnoscStudiaNastepnaKolejnosc)
	if err != nil {
		return 0, err
	}
	var kolejnosc int64
	if err := polecenie.QueryRowContext(ctx, dokumentID).Scan(&kolejnosc); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć kolejności czynności dokumentu %d: %w",
			dokumentID, err)
	}
	return kolejnosc, nil
}

// CzynnoscDokumentu oddaje wpis dziennika o wskazanym kodzie, bez zależności.
func (r *repozytoriumStudia) CzynnoscDokumentu(ctx context.Context,
	kod string) (CzynnoscDokumentuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, czynnoscStudiaPobierz)
	if err != nil {
		return CzynnoscDokumentuStudia{}, err
	}
	czynnosc, err := odczytajCzynnoscDokumentuStudia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return CzynnoscDokumentuStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return CzynnoscDokumentuStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz czynności dokumentu %q: %w", kod, err)
	}
	return czynnosc, nil
}

// CzynnosciDokumentu oddaje dziennik dokumentu od najświeższego wpisu, wraz
// z obu kierunkami zależności wypełnionymi.
//
// Zależności doczytujemy JEDNYM zapytaniem dla całego dokumentu, nie zapytaniem
// na wpis: dziennik długiego dokumentu ma setki wpisów, a pytanie na wpis
// zamieniłoby odczyt wykazu w setki zapytań.
func (r *repozytoriumStudia) CzynnosciDokumentu(ctx context.Context,
	dokumentID int64) ([]CzynnoscDokumentuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, czynnosciStudiaLista)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dziennika czynności: %w", err)
	}
	defer wiersze.Close()

	lista := []CzynnoscDokumentuStudia{}
	for wiersze.Next() {
		czynnosc, err := odczytajCzynnoscDokumentuStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz dziennika czynności: %w", err)
		}
		lista = append(lista, czynnosc)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt dziennika czynności: %w", err)
	}

	naNiej, podstawy, err := r.zaleznosciCzynnosciStudia(ctx, dokumentID)
	if err != nil {
		return nil, err
	}
	for i := range lista {
		lista[i].PodstawyKody = podstawy[lista[i].Kod]
		lista[i].StojaceNaNiej = naNiej[lista[i].Kod]
	}
	return lista, nil
}

// zaleznosciCzynnosciStudia oddaje dwie mapy: „co stoi na tej czynności" oraz
// „na czym ta czynność stoi".
func (r *repozytoriumStudia) zaleznosciCzynnosciStudia(ctx context.Context,
	dokumentID int64) (map[string][]string, map[string][]string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zaleznosciStudiaPary)
	if err != nil {
		return nil, nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, nil, fmt.Errorf("dane: nie można odczytać zależności czynności: %w", err)
	}
	defer wiersze.Close()

	naNiej, podstawy := map[string][]string{}, map[string][]string{}
	for wiersze.Next() {
		var czynnosc, podstawa string
		if err := wiersze.Scan(&czynnosc, &podstawa); err != nil {
			return nil, nil, fmt.Errorf("dane: nieczytelny wiersz zależności czynności: %w", err)
		}
		naNiej[podstawa] = append(naNiej[podstawa], czynnosc)
		podstawy[czynnosc] = append(podstawy[czynnosc], podstawa)
	}
	if err := wiersze.Err(); err != nil {
		return nil, nil, fmt.Errorf("dane: przerwany odczyt zależności czynności: %w", err)
	}
	return naNiej, podstawy, nil
}

// ZapiszZaleznoscCzynnosci zapisuje, że czynność stoi na podstawie.
func (r *repozytoriumStudia) ZapiszZaleznoscCzynnosci(ctx context.Context,
	czynnoscKod, podstawaKod string, powod *string) error {

	if czynnoscKod == "" || podstawaKod == "" || czynnoscKod == podstawaKod {
		return nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zaleznoscStudiaZapisz)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, czynnoscKod, podstawaKod,
		tekstDoKolumny(powod)); err != nil {

		return fmt.Errorf("dane: nie można zapisać zależności czynności %q od %q: %w",
			czynnoscKod, podstawaKod, err)
	}
	return nil
}

// PrzestawStanCzynnosci przestawia stan wpisu dziennika i mówi, czy przestawił.
//
// Stan oczekiwany jest częścią warunku, żeby cofnięcie czynności już cofniętej
// oddało `false` zamiast cicho przejść — Operator musi wiedzieć, że drugi ruch
// nie zrobił nic.
func (r *repozytoriumStudia) PrzestawStanCzynnosci(ctx context.Context,
	kod, stanOczekiwany, stanNowy string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, czynnoscStudiaPrzestawStan)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, stanNowy, kod, stanOczekiwany)
	if err != nil {
		return false, fmt.Errorf("dane: nie można przestawić stanu czynności %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek przestawienia czynności %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// ── Kopie zapasowe ──────────────────────────────────────────────────────────

// ZapiszKopieZapasowa odkłada kopię zapasową dokumentu — także nieudaną.
func (r *repozytoriumStudia) ZapiszKopieZapasowa(ctx context.Context, dokumentID int64,
	kopia KopiaZapasowaStudia) (KopiaZapasowaStudia, error) {

	if kopia.Kod == "" {
		return KopiaZapasowaStudia{}, fmt.Errorf("dane: kopia zapasowa bez identyfikatora")
	}
	if kopia.Powod == "" {
		kopia.Powod = "manual"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, kopiaStudiaZapisz)
	if err != nil {
		return KopiaZapasowaStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, kopia.Kod, dokumentID, kopia.Powod,
		tekstDoKolumny(kopia.Tresc), tekstDoKolumny(kopia.TrescOdwolanie),
		tekstDoKolumny(kopia.PostacJSON), kopia.RozmiarBajtow, liczbaLogiczna(kopia.UdaloSie),
		tekstDoKolumny(kopia.PowodNiepowodzenia), liczbaLogiczna(kopia.ZmianyNiezapisane))
	if err != nil {
		return KopiaZapasowaStudia{}, fmt.Errorf(
			"dane: nie można zapisać kopii zapasowej %q: %w", kopia.Kod, err)
	}
	return r.KopiaZapasowa(ctx, kopia.Kod)
}

// KopiaZapasowa oddaje kopię o wskazanym kodzie.
func (r *repozytoriumStudia) KopiaZapasowa(ctx context.Context,
	kod string) (KopiaZapasowaStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, kopiaStudiaPobierz)
	if err != nil {
		return KopiaZapasowaStudia{}, err
	}
	kopia, err := odczytajKopieZapasowaStudia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return KopiaZapasowaStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return KopiaZapasowaStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz kopii zapasowej %q: %w", kod, err)
	}
	return kopia, nil
}

// KopieZapasowe oddaje kopie dokumentu, od najświeższej.
func (r *repozytoriumStudia) KopieZapasowe(ctx context.Context,
	dokumentID int64) ([]KopiaZapasowaStudia, error) {

	return r.kopieZapasoweStudia(ctx, kopieStudiaLista, dokumentID)
}

// KopieNiezapisane oddaje kopie niosące zmiany niezapisane. Okno puste znaczy
// wszystkie okna — po tym Studio zgłasza przywrócenie po nagłym zamknięciu.
func (r *repozytoriumStudia) KopieNiezapisane(ctx context.Context,
	okno string) ([]KopiaZapasowaStudia, error) {

	return r.kopieZapasoweStudia(ctx, kopieStudiaNiezapisane, okno, okno)
}

// kopieZapasoweStudia jest wspólnym odczytem wykazu kopii.
func (r *repozytoriumStudia) kopieZapasoweStudia(ctx context.Context,
	polecenieSQL string, argumenty ...any) ([]KopiaZapasowaStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, polecenieSQL)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kopii zapasowych: %w", err)
	}
	defer wiersze.Close()

	lista := []KopiaZapasowaStudia{}
	for wiersze.Next() {
		kopia, err := odczytajKopieZapasowaStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kopii zapasowej: %w", err)
		}
		lista = append(lista, kopia)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kopii zapasowych: %w", err)
	}
	return lista, nil
}

// PrzemiecKopieZapasowe usuwa kopie wygasłe wedle nastawy Operatora i oddaje,
// ile ich zdjęto. Kopia nieudana i kopia niosąca zmiany niezapisane zostają.
func (r *repozytoriumStudia) PrzemiecKopieZapasowe(ctx context.Context, dokumentID int64,
	ileZachowac, wygasanieGodzin int64) (int64, error) {

	if ileZachowac <= 0 {
		ileZachowac = 1
	}
	if wygasanieGodzin <= 0 {
		wygasanieGodzin = 1
	}
	polecenie, err := r.zapytania.przygotuj(ctx, kopieStudiaPrzemiec)
	if err != nil {
		return 0, err
	}
	// Modyfikator SQLite w postaci „-168 hours"; liczba i jednostka jedzą jednym
	// napisem, bo `strftime` nie przyjmuje ich osobno.
	granica := fmt.Sprintf("-%d hours", wygasanieGodzin)
	wynik, err := polecenie.ExecContext(ctx, dokumentID, granica, dokumentID, ileZachowac)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można przemieść kopii zapasowych dokumentu %d: %w",
			dokumentID, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany skutek przemiatania kopii dokumentu %d: %w",
			dokumentID, err)
	}
	return zdjete, nil
}

// ── Nastawy pracy: autozapis i wygasanie kopii ──────────────────────────────

// NastawaPracy oddaje nastawy pracy dla pary (okno, dokument) albo dla samego
// okna, zakładając wiersz domyślny, gdy go jeszcze nie ma.
//
// Wiersz zakłada się przy odczycie, nie przy zapisie, bo `autosave.get` ma
// oddać nastawy obowiązujące także wtedy, gdy Operator nigdy ich nie ruszał —
// a wtedy obowiązują wartości domyślne kolumn, i to jest odpowiedź, nie brak.
func (r *repozytoriumStudia) NastawaPracy(ctx context.Context, okno string,
	dokumentID *int64) (NastawaPracyStudia, error) {

	if okno == "" {
		return NastawaPracyStudia{}, fmt.Errorf("dane: nastawy pracy studio bez okna")
	}
	if err := r.zalozNastaweStudia(ctx, okno, dokumentID); err != nil {
		return NastawaPracyStudia{}, err
	}
	polecenieSQL, argumenty := nastawaPracyPobierzOkno, []any{okno}
	if dokumentID != nil {
		polecenieSQL, argumenty = nastawaPracyPobierzDokument, []any{okno, *dokumentID}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, polecenieSQL)
	if err != nil {
		return NastawaPracyStudia{}, err
	}
	nastawa, err := odczytajNastaweStudia(polecenie.QueryRowContext(ctx, argumenty...))
	if errors.Is(err, sql.ErrNoRows) {
		return NastawaPracyStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return NastawaPracyStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz nastaw pracy studio okna %q: %w", okno, err)
	}
	return nastawa, nil
}

// zalozNastaweStudia zakłada wiersz nastaw, jeśli go nie ma.
func (r *repozytoriumStudia) zalozNastaweStudia(ctx context.Context, okno string,
	dokumentID *int64) error {

	polecenieSQL, argumenty := nastawaPracyZalozOkno, []any{okno, okno}
	if dokumentID != nil {
		polecenieSQL = nastawaPracyZalozDokument
		argumenty = []any{okno, *dokumentID, okno, *dokumentID}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, polecenieSQL)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, argumenty...); err != nil {
		return fmt.Errorf("dane: nie można założyć nastaw pracy studio okna %q: %w", okno, err)
	}
	return nil
}

// ZapiszNastaweAutozapisu zapisuje nastawy autozapisu i wygasania kopii.
func (r *repozytoriumStudia) ZapiszNastaweAutozapisu(ctx context.Context,
	nastawa NastawaPracyStudia) error {

	if nastawa.ID == 0 {
		return fmt.Errorf("dane: zapis nastaw autozapisu bez wiersza nastaw")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, nastawaPracyZapiszAutozapis)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, liczbaLogiczna(nastawa.AutozapisCzynny),
		nastawa.AutozapisOdstepSekund, liczbaLogiczna(nastawa.AutozapisPrzyOdejsciu),
		liczbaLogiczna(nastawa.AutozapisPrzyZamknieciu),
		liczbaLogiczna(nastawa.AutozapisPrzyPrzelacz), nastawa.KopieIleZachowac,
		nastawa.KopieWygasanieGodzin, nastawa.ID)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać nastaw autozapisu %d: %w", nastawa.ID, err)
	}
	return nil
}

// ZapiszSkutekAutozapisu odkłada wynik zapisu samoczynnego — także nieudany.
//
// Osobne polecenie, bo skutek pisze zapis, a nastawy pisze Operator. Nieudany
// zapis musi zostawić ślad NAZWANY: wskaźnik „zapisano" pokazany po
// niepowodzeniu jest najgorszym możliwym błędem tego modułu.
func (r *repozytoriumStudia) ZapiszSkutekAutozapisu(ctx context.Context, nastawaID int64,
	chwila *string, nieudany bool, powod *string) error {

	if nastawaID == 0 {
		return fmt.Errorf("dane: zapis skutku autozapisu bez wiersza nastaw")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, nastawaPracyZapiszSkutek)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, tekstDoKolumny(chwila),
		liczbaLogiczna(nieudany), tekstDoKolumny(powod), nastawaID); err != nil {

		return fmt.Errorf("dane: nie można zapisać skutku autozapisu %d: %w", nastawaID, err)
	}
	return nil
}

// ── Składanie wierszy ───────────────────────────────────────────────────────

func odczytajBlokadeFragmentuStudia(wiersz skaner) (BlokadaFragmentuStudia, error) {
	var blokada BlokadaFragmentuStudia
	var powod, agentKod, agentNazwa, szablon sql.NullString
	var zSzablonu int
	err := wiersz.Scan(&blokada.ID, &blokada.Kod, &blokada.DokumentKod, &blokada.Nazwa,
		&powod, &blokada.ZakresOd, &blokada.ZakresDo, &blokada.Zasieg,
		&blokada.ZalozylRodzaj, &agentKod, &agentNazwa, &zSzablonu, &szablon,
		&blokada.Przesuniecia, &blokada.Utworzono)
	if err != nil {
		return BlokadaFragmentuStudia{}, err
	}
	blokada.Powod = tekstZKolumny(powod)
	blokada.ZalozylAgentKod = tekstZKolumny(agentKod)
	blokada.ZalozylAgentNaz = tekstZKolumny(agentNazwa)
	blokada.SzablonKod = tekstZKolumny(szablon)
	blokada.ZSzablonu = zSzablonu == 1
	return blokada, nil
}

func odczytajCzynnoscDokumentuStudia(wiersz skaner) (CzynnoscDokumentuStudia, error) {
	var czynnosc CzynnoscDokumentuStudia
	var agentKod, agentNazwa, agentWersja, podagent sql.NullString
	var przed, po, zmiana, zadanie sql.NullString
	var od, do sql.NullInt64
	err := wiersz.Scan(&czynnosc.ID, &czynnosc.Kod, &czynnosc.DokumentKod, &czynnosc.Kolejnosc,
		&czynnosc.Rodzaj, &czynnosc.AutorRodzaj, &agentKod, &agentNazwa, &agentWersja,
		&podagent, &czynnosc.Opis, &od, &do, &przed, &po, &zmiana, &zadanie,
		&czynnosc.Stan, &czynnosc.Utworzono)
	if err != nil {
		return CzynnoscDokumentuStudia{}, err
	}
	czynnosc.AutorAgentKod = tekstZKolumny(agentKod)
	czynnosc.AutorAgentNazwa = tekstZKolumny(agentNazwa)
	czynnosc.AutorAgentWersja = tekstZKolumny(agentWersja)
	czynnosc.AutorPodagentKod = tekstZKolumny(podagent)
	czynnosc.ZakresOd = liczbaZKolumny(od)
	czynnosc.ZakresDo = liczbaZKolumny(do)
	czynnosc.StanPrzed = tekstZKolumny(przed)
	czynnosc.StanPo = tekstZKolumny(po)
	czynnosc.ZmianaSledzonaK = tekstZKolumny(zmiana)
	czynnosc.ZadanieKod = tekstZKolumny(zadanie)
	return czynnosc, nil
}

func odczytajKopieZapasowaStudia(wiersz skaner) (KopiaZapasowaStudia, error) {
	var kopia KopiaZapasowaStudia
	var tresc, odwolanie, postac, powodNiepow sql.NullString
	var udaloSie, niezapisane int
	err := wiersz.Scan(&kopia.ID, &kopia.Kod, &kopia.DokumentKod, &kopia.Powod,
		&tresc, &odwolanie, &postac, &kopia.RozmiarBajtow, &udaloSie, &powodNiepow,
		&niezapisane, &kopia.Utworzono)
	if err != nil {
		return KopiaZapasowaStudia{}, err
	}
	kopia.Tresc = tekstZKolumny(tresc)
	kopia.TrescOdwolanie = tekstZKolumny(odwolanie)
	kopia.PostacJSON = tekstZKolumny(postac)
	kopia.PowodNiepowodzenia = tekstZKolumny(powodNiepow)
	kopia.UdaloSie = udaloSie == 1
	kopia.ZmianyNiezapisane = niezapisane == 1
	return kopia, nil
}

func odczytajNastaweStudia(wiersz skaner) (NastawaPracyStudia, error) {
	var nastawa NastawaPracyStudia
	var dokument sql.NullInt64
	var ostatniZapis, powodNiepow sql.NullString
	var czynny, przyOdejsciu, przyZamknieciu, przyPrzelacz, nieudany, podswietlenie int
	err := wiersz.Scan(&nastawa.ID, &nastawa.Okno, &dokument, &czynny,
		&nastawa.AutozapisOdstepSekund, &przyOdejsciu, &przyZamknieciu, &przyPrzelacz,
		&nastawa.KopieIleZachowac, &nastawa.KopieWygasanieGodzin, &ostatniZapis,
		&nieudany, &powodNiepow, &podswietlenie)
	if err != nil {
		return NastawaPracyStudia{}, err
	}
	nastawa.DokumentID = liczbaZKolumny(dokument)
	nastawa.OstatniZapis = tekstZKolumny(ostatniZapis)
	nastawa.OstatniPowodNiepowodz = tekstZKolumny(powodNiepow)
	nastawa.AutozapisCzynny = czynny == 1
	nastawa.AutozapisPrzyOdejsciu = przyOdejsciu == 1
	nastawa.AutozapisPrzyZamknieciu = przyZamknieciu == 1
	nastawa.AutozapisPrzyPrzelacz = przyPrzelacz == 1
	nastawa.OstatniZapisNieudany = nieudany == 1
	nastawa.PodswietlenieZmianModelu = podswietlenie == 1
	return nastawa, nil
}
