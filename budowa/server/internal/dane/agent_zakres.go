// Odpowiedzialność pliku: zakres działania eksperta — moduły zastosowania,
// osiem zakresów izolacji technicznej, granica Subagent Network, zdjęcie
// wpisów uprawnień oraz odczyt konektorów i przypisań widziany od strony
// eksperta.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
)

// PrzelacznikIzolacjiAgenta to jeden z ośmiu zakresów izolacji technicznej
// zapisany przy ekspercie jako jego wartość wyjściowa.
type PrzelacznikIzolacjiAgenta struct {
	// Zakres niesie wartość wyliczenia `IsolationTechnicalScope` kontraktu.
	Zakres string
	// Odciety mówi, czy zakres jest odcinany. Fałsz jest stanem wyjściowym.
	Odciety bool
}

// PrzypisanieEksperta to jedno przypisanie widziane od strony eksperta: projekt
// modułu Workspace albo rola środowiska MultitaskingAI.
type PrzypisanieEksperta struct {
	AgentKod string
	// Rodzaj niesie wartość wyliczenia `AgentAssignmentKind`: `project` albo
	// `role`.
	Rodzaj string
	// CelKod jest kodem projektu albo identyfikatorem stanowiska obsady biegu.
	CelKod string
	// CelNazwa jest nazwą projektu albo biegu; pusta oznacza byt bez nazwy
	// własnej.
	CelNazwa string
	// Rola niesie rolę eksperta w projekcie albo wcielenie roli w biegu.
	Rola string
	// DomyslnyWykonawca dotyczy wyłącznie przypisania do projektu.
	DomyslnyWykonawca bool
	Przypisano        string
}

// RepozytoriumZakresuAgenta jest kontraktem zakresu działania eksperta:
// modułów, izolacji, granicy podagentów, uprawnień, konektorów i przypisań.
type RepozytoriumZakresuAgenta interface {
	// ModulyAgenta oddaje kody modułów zastosowania; pusty wycinek znaczy brak
	// ograniczenia.
	ModulyAgenta(ctx context.Context, kodAgenta string) ([]string, error)
	// UstawModulyAgenta zastępuje komplet modułów zastosowania; wycinek pusty
	// zdejmuje ograniczenie.
	UstawModulyAgenta(ctx context.Context, kodAgenta string, kody []string) error
	// IzolacjaAgenta oddaje zapisane przełączniki izolacji; komplet ośmiu
	// zakresów składa warstwa wyższa.
	IzolacjaAgenta(ctx context.Context, kodAgenta string) ([]PrzelacznikIzolacjiAgenta, error)
	// UstawIzolacjeAgenta zapisuje wskazane przełączniki; zakres pominięty
	// w wycinku zostaje bez zmiany.
	UstawIzolacjeAgenta(ctx context.Context, kodAgenta string, przelaczniki []PrzelacznikIzolacjiAgenta) error
	// GranicaPodagentow oddaje górną liczbę podagentów eksperta; zero znaczy
	// Subagent Network wyłączony.
	GranicaPodagentow(ctx context.Context, kodAgenta string) (int, error)
	// UstawGranicePodagentow zapisuje tę granicę.
	UstawGranicePodagentow(ctx context.Context, kodAgenta string, granica int) error
	// UsunUprawnieniaAgenta zdejmuje wpisy uprawnień i oddaje liczbę zdjętych,
	// zawężone grupą i zakresem.
	UsunUprawnieniaAgenta(ctx context.Context, kodAgenta, grupa, zakres string) (int, error)
	// KonektoryAgenta oddaje konektory eksperta wraz z kodem punktu dostępu.
	KonektoryAgenta(ctx context.Context, kodAgenta string) ([]KonektorAgenta, []string, error)
	// UsunKonektorAgenta odłącza konektor od eksperta.
	UsunKonektorAgenta(ctx context.Context, kodAgenta, kodKonektora string) (bool, error)
	// ZapiszKonektorAgenta zmienia punkt dostępu, konfigurację i stan czynności
	// konektora eksperta.
	ZapiszKonektorAgenta(ctx context.Context, kodAgenta, kodKonektora string,
		punktID *int64, konfiguracja *string, aktywny *bool) (KonektorAgenta, string, error)
	// UsunUmiejetnoscAgenta zdejmuje umiejętność z definicji eksperta.
	UsunUmiejetnoscAgenta(ctx context.Context, kodAgenta, kodUmiejetnosci string) (bool, error)
	// PrzypisaniaEkspertow oddaje przypisania, zawężone kodem eksperta
	// i rodzajem, gdy podane.
	PrzypisaniaEkspertow(ctx context.Context, kodAgenta, rodzaj string) ([]PrzypisanieEksperta, error)
}

const (
	kolumnyKonektoraZakresu = `k.id, k.kod, k.agent_id, a.kod, k.nazwa, k.rodzaj,
	                           k.punkt_dostepu_id, k.konfiguracja, k.aktywny, k.utworzono,
	                           COALESCE(p.kod, '')`

	modulyZakresuAgenta = `SELECT m.kod_modulu FROM agent_modul_zastosowania m
	                         JOIN agent a ON a.id = m.agent_id
	                        WHERE a.kod = ? ORDER BY m.kod_modulu`

	czyscModulyZakresuAgenta = `DELETE FROM agent_modul_zastosowania WHERE agent_id = ?`

	wstawModulZakresuAgenta = `INSERT INTO agent_modul_zastosowania (agent_id, kod_modulu)
	                           VALUES (?, ?) ON CONFLICT(agent_id, kod_modulu) DO NOTHING`

	izolacjaZakresuAgenta = `SELECT i.zakres, i.odciety FROM agent_izolacja_techniczna i
	                           JOIN agent a ON a.id = i.agent_id
	                          WHERE a.kod = ? ORDER BY i.zakres`

	zapiszIzolacjeZakresuAgenta = `INSERT INTO agent_izolacja_techniczna (agent_id, zakres, odciety)
	                               VALUES (?, ?, ?)
	                               ON CONFLICT(agent_id, zakres) DO UPDATE SET
	                                   odciety = excluded.odciety,
	                                   zapisano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	granicaPodagentowAgenta = `SELECT limit_podagentow FROM agent WHERE kod = ?`

	zapiszGranicePodagentow = `UPDATE agent SET limit_podagentow = ?,
	                               zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                            WHERE kod = ?`

	usunUprawnieniaZakresuAgenta = `DELETE FROM agent_uprawnienie
	                                 WHERE agent_id = ?
	                                   AND (? = '' OR grupa = ?)
	                                   AND (? = '' OR zakres = ?)`

	konektoryZakresuAgenta = `SELECT ` + kolumnyKonektoraZakresu + `
	                            FROM agent_konektor k
	                            JOIN agent a ON a.id = k.agent_id
	                            LEFT JOIN punkt_dostepu p ON p.id = k.punkt_dostepu_id
	                           WHERE a.kod = ? ORDER BY k.utworzono, k.id`

	konektorZakresuPoKodzie = `SELECT ` + kolumnyKonektoraZakresu + `
	                             FROM agent_konektor k
	                             JOIN agent a ON a.id = k.agent_id
	                             LEFT JOIN punkt_dostepu p ON p.id = k.punkt_dostepu_id
	                            WHERE k.kod = ? AND a.kod = ?`

	usunKonektorZakresuAgenta = `DELETE FROM agent_konektor
	                              WHERE kod = ? AND agent_id = (SELECT id FROM agent WHERE kod = ?)`

	usunUmiejetnoscZakresuAgenta = `DELETE FROM agent_umiejetnosc
	                                 WHERE kod = ? AND agent_id = (SELECT id FROM agent WHERE kod = ?)`

	// Przypisanie do projektu: zapytanie czyta Agent Managera modułu Workspace,
	// zapisany migracją 035, wraz z rolą i domyślnym wykonawcą.
	przypisaniaProjektoweEkspertow = `SELECT pa.agent_kod, p.kod, p.nazwa,
	                                         COALESCE(pa.rola, ''), pa.domyslny_wykonawca, pa.przypisano
	                                    FROM przypisanie_agenta_projektu pa
	                                    JOIN projekt p ON p.id = pa.projekt_id
	                                   WHERE (? = '' OR pa.agent_kod = ?)
	                                   ORDER BY pa.przypisano, p.kod`

	// Przypisanie do roli: stanowisko obsady biegu orkiestracji (migracja 077).
	// Celem jest stanowisko, bo to ono niesie rolę i miejsce w zespole.
	przypisaniaRolowaEkspertow = `SELECT a.kod, o.id, COALESCE(b.nazwa, ''), o.rola, o.utworzono
	                                FROM obsada_biegu o
	                                JOIN agent a ON a.id = o.agent_id
	                                LEFT JOIN bieg_orkiestracji b ON b.id = o.bieg_id
	                               WHERE o.agent_id IS NOT NULL AND (? = '' OR a.kod = ?)
	                               ORDER BY o.utworzono, o.id`
)

// repozytoriumZakresuAgenta obsługuje zakres działania eksperta, wiążąc
// pamięć przygotowanych zapytań z bazą danych repozytorium.
type repozytoriumZakresuAgenta struct {
	zapytania *zapytania
	db        *sql.DB
}

var _ RepozytoriumZakresuAgenta = (*repozytoriumZakresuAgenta)(nil)

// noweRepozytoriumZakresuAgenta wiąże zakres działania eksperta z bazą danych
// i pamięcią przygotowanych zapytań.
func noweRepozytoriumZakresuAgenta(z *zapytania, db *sql.DB) *repozytoriumZakresuAgenta {
	return &repozytoriumZakresuAgenta{zapytania: z, db: db}
}

// numerEkspertaZakresu przekłada kod eksperta na klucz wiersza. Kod nieznany
// wraca jako ErrBrakWiersza — warstwa wyższa odróżnia „takiego eksperta nie ma"
// od „odczyt się nie powiódł".
func (r *repozytoriumZakresuAgenta) numerEkspertaZakresu(ctx context.Context, kod string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, `SELECT id FROM agent WHERE kod = ?`)
	if err != nil {
		return 0, err
	}
	var id int64
	switch err := polecenie.QueryRowContext(ctx, kod).Scan(&id); {
	case err == sql.ErrNoRows:
		return 0, ErrBrakWiersza
	case err != nil:
		return 0, fmt.Errorf("dane: nie można odczytać eksperta %q: %w", kod, err)
	}
	return id, nil
}

// ModulyAgenta oddaje kody modułów zastosowania wskazanego eksperta,
// uporządkowane alfabetycznie według kodu modułu.
func (r *repozytoriumZakresuAgenta) ModulyAgenta(ctx context.Context, kodAgenta string) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, modulyZakresuAgenta)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać modułów eksperta %q: %w", kodAgenta, err)
	}
	defer wiersze.Close()
	kody := make([]string, 0, 8)
	for wiersze.Next() {
		var kod string
		if err := wiersze.Scan(&kod); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny moduł eksperta %q: %w", kodAgenta, err)
		}
		kody = append(kody, kod)
	}
	return kody, wiersze.Err()
}

// UstawModulyAgenta zastępuje komplet modułów zastosowania jednym przebiegiem
// w transakcji: stan pośredni, w którym stary komplet już zszedł, a nowy
// jeszcze nie wszedł, byłby chwilą pełnego dostępu wziętą z pomyłki.
func (r *repozytoriumZakresuAgenta) UstawModulyAgenta(ctx context.Context,
	kodAgenta string, kody []string) error {

	id, err := r.numerEkspertaZakresu(ctx, kodAgenta)
	if err != nil {
		return err
	}
	transakcja, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("dane: nie można otworzyć zapisu modułów eksperta %q: %w", kodAgenta, err)
	}
	defer func() { _ = transakcja.Rollback() }()

	if _, err := transakcja.ExecContext(ctx, czyscModulyZakresuAgenta, id); err != nil {
		return fmt.Errorf("dane: nie można wyczyścić modułów eksperta %q: %w", kodAgenta, err)
	}
	for _, kod := range kody {
		przyciety := strings.TrimSpace(kod)
		if przyciety == "" {
			continue
		}
		if _, err := transakcja.ExecContext(ctx, wstawModulZakresuAgenta, id, przyciety); err != nil {
			return fmt.Errorf("dane: nie można zapisać modułu %q eksperta %q: %w", przyciety, kodAgenta, err)
		}
	}
	if err := transakcja.Commit(); err != nil {
		return fmt.Errorf("dane: nie można domknąć zapisu modułów eksperta %q: %w", kodAgenta, err)
	}
	return nil
}

// IzolacjaAgenta oddaje zapisane przełączniki izolacji technicznej wskazanego
// eksperta, uporządkowane po zakresie.
func (r *repozytoriumZakresuAgenta) IzolacjaAgenta(ctx context.Context,
	kodAgenta string) ([]PrzelacznikIzolacjiAgenta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, izolacjaZakresuAgenta)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać izolacji eksperta %q: %w", kodAgenta, err)
	}
	defer wiersze.Close()
	przelaczniki := make([]PrzelacznikIzolacjiAgenta, 0, 8)
	for wiersze.Next() {
		var przelacznik PrzelacznikIzolacjiAgenta
		var odciety int
		if err := wiersze.Scan(&przelacznik.Zakres, &odciety); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny przełącznik izolacji eksperta %q: %w", kodAgenta, err)
		}
		przelacznik.Odciety = odciety == 1
		przelaczniki = append(przelaczniki, przelacznik)
	}
	return przelaczniki, wiersze.Err()
}

// UstawIzolacjeAgenta zapisuje wskazane przełączniki izolacji technicznej
// eksperta w jednej transakcji.
func (r *repozytoriumZakresuAgenta) UstawIzolacjeAgenta(ctx context.Context,
	kodAgenta string, przelaczniki []PrzelacznikIzolacjiAgenta) error {

	id, err := r.numerEkspertaZakresu(ctx, kodAgenta)
	if err != nil {
		return err
	}
	transakcja, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("dane: nie można otworzyć zapisu izolacji eksperta %q: %w", kodAgenta, err)
	}
	defer func() { _ = transakcja.Rollback() }()

	for _, przelacznik := range przelaczniki {
		zakres := strings.TrimSpace(przelacznik.Zakres)
		if zakres == "" {
			continue
		}
		if _, err := transakcja.ExecContext(ctx, zapiszIzolacjeZakresuAgenta,
			id, zakres, liczbaLogiczna(przelacznik.Odciety)); err != nil {
			return fmt.Errorf("dane: nie można zapisać zakresu %q eksperta %q: %w", zakres, kodAgenta, err)
		}
	}
	if err := transakcja.Commit(); err != nil {
		return fmt.Errorf("dane: nie można domknąć zapisu izolacji eksperta %q: %w", kodAgenta, err)
	}
	return nil
}

// GranicaPodagentow oddaje górną liczbę jednoczesnych podagentów wskazanego
// eksperta zapisaną przy nim.
func (r *repozytoriumZakresuAgenta) GranicaPodagentow(ctx context.Context, kodAgenta string) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, granicaPodagentowAgenta)
	if err != nil {
		return 0, err
	}
	var granica int
	switch err := polecenie.QueryRowContext(ctx, kodAgenta).Scan(&granica); {
	case err == sql.ErrNoRows:
		return 0, ErrBrakWiersza
	case err != nil:
		return 0, fmt.Errorf("dane: nie można odczytać granicy podagentów eksperta %q: %w", kodAgenta, err)
	}
	return granica, nil
}

// UstawGranicePodagentow zapisuje granicę Subagent Network eksperta i mówi,
// czy wiersz eksperta istniał.
func (r *repozytoriumZakresuAgenta) UstawGranicePodagentow(ctx context.Context,
	kodAgenta string, granica int) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszGranicePodagentow)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, granica, kodAgenta)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać granicy podagentów eksperta %q: %w", kodAgenta, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err == nil && zmienione == 0 {
		return ErrBrakWiersza
	}
	return nil
}

// UsunUprawnieniaAgenta zdejmuje wpisy uprawnień i oddaje liczbę zdjętych.
// Zero nie jest odmową: znaczy, że zawężeń nie było.
func (r *repozytoriumZakresuAgenta) UsunUprawnieniaAgenta(ctx context.Context,
	kodAgenta, grupa, zakres string) (int, error) {

	id, err := r.numerEkspertaZakresu(ctx, kodAgenta)
	if err != nil {
		return 0, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, usunUprawnieniaZakresuAgenta)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, id, grupa, grupa, zakres, zakres)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zdjąć uprawnień eksperta %q: %w", kodAgenta, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć zdjętych uprawnień eksperta %q: %w", kodAgenta, err)
	}
	return int(zdjete), nil
}

// KonektoryAgenta oddaje konektory eksperta wraz z kodami punktów dostępu.
// Dwa wycinki tej samej długości: kontrakt niesie kod punktu, a nie jego numer
// wiersza, więc odczyt bierze go złączeniem zamiast zostawiać wołającemu.
func (r *repozytoriumZakresuAgenta) KonektoryAgenta(ctx context.Context,
	kodAgenta string) ([]KonektorAgenta, []string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, konektoryZakresuAgenta)
	if err != nil {
		return nil, nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta)
	if err != nil {
		return nil, nil, fmt.Errorf("dane: nie można odczytać konektorów eksperta %q: %w", kodAgenta, err)
	}
	defer wiersze.Close()
	konektory := make([]KonektorAgenta, 0, 8)
	punkty := make([]string, 0, 8)
	for wiersze.Next() {
		konektor, kodPunktu, err := odczytajKonektorZakresu(wiersze)
		if err != nil {
			return nil, nil, err
		}
		konektory = append(konektory, konektor)
		punkty = append(punkty, kodPunktu)
	}
	return konektory, punkty, wiersze.Err()
}

// czytelnikWierszaZakresu obejmuje `*sql.Row` i `*sql.Rows` jedną nazwą, żeby
// odczyt konektora stał w jednym miejscu dla obu dróg.
type czytelnikWierszaZakresu interface {
	Scan(cele ...any) error
}

// odczytajKonektorZakresu składa strukturę konektora wraz z kodem punktu
// dostępu z jednego wiersza wyniku.
func odczytajKonektorZakresu(zrodlo czytelnikWierszaZakresu) (KonektorAgenta, string, error) {
	var wpis KonektorAgenta
	var punkt sql.NullInt64
	var aktywny int
	var kodPunktu string
	err := zrodlo.Scan(&wpis.ID, &wpis.Kod, &wpis.AgentID, &wpis.AgentKod, &wpis.Nazwa,
		&wpis.Rodzaj, &punkt, &wpis.Konfiguracja, &aktywny, &wpis.Utworzono, &kodPunktu)
	if err == sql.ErrNoRows {
		return KonektorAgenta{}, "", ErrBrakWiersza
	}
	if err != nil {
		return KonektorAgenta{}, "", fmt.Errorf("dane: nieczytelny wiersz konektora eksperta: %w", err)
	}
	wpis.PunktDostepuID = liczbaZKolumny(punkt)
	wpis.Aktywny = aktywny == 1
	return wpis, kodPunktu, nil
}

// UsunKonektorAgenta odłącza konektor od eksperta. Brak wiersza nie jest
// odmową — wynik `false` mówi, że nie było czego odłączać.
func (r *repozytoriumZakresuAgenta) UsunKonektorAgenta(ctx context.Context,
	kodAgenta, kodKonektora string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunKonektorZakresuAgenta)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kodKonektora, kodAgenta)
	if err != nil {
		return false, fmt.Errorf("dane: nie można odłączyć konektora %q: %w", kodKonektora, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć odłączeń konektora %q: %w", kodKonektora, err)
	}
	return zdjete > 0, nil
}

// ZapiszKonektorAgenta zmienia konfigurację instancji konektora. Wskaźnik pusty
// zostawia pole bez zmiany — okno konfiguracji instancji zapisuje pojedyncze
// pole, a nie komplet.
func (r *repozytoriumZakresuAgenta) ZapiszKonektorAgenta(ctx context.Context,
	kodAgenta, kodKonektora string, punktID *int64, konfiguracja *string,
	aktywny *bool) (KonektorAgenta, string, error) {

	biezacy, _, err := r.konektorZakresu(ctx, kodAgenta, kodKonektora)
	if err != nil {
		return KonektorAgenta{}, "", err
	}
	nowyPunkt := biezacy.PunktDostepuID
	if punktID != nil {
		nowyPunkt = punktID
	}
	nowaKonfiguracja := biezacy.Konfiguracja
	if konfiguracja != nil {
		nowaKonfiguracja = *konfiguracja
	}
	nowyStan := biezacy.Aktywny
	if aktywny != nil {
		nowyStan = *aktywny
	}
	polecenie, err := r.zapytania.przygotuj(ctx,
		`UPDATE agent_konektor SET punkt_dostepu_id = ?, konfiguracja = ?, aktywny = ? WHERE id = ?`)
	if err != nil {
		return KonektorAgenta{}, "", err
	}
	if _, err := polecenie.ExecContext(ctx, liczbaDoKolumny(nowyPunkt), nowaKonfiguracja,
		liczbaLogiczna(nowyStan), biezacy.ID); err != nil {
		return KonektorAgenta{}, "", fmt.Errorf("dane: nie można zapisać konektora %q: %w", kodKonektora, err)
	}
	return r.konektorZakresu(ctx, kodAgenta, kodKonektora)
}

// konektorZakresu odczytuje jeden konektor wskazanego eksperta wraz z kodem
// punktu dostępu, po kodzie konektora.
func (r *repozytoriumZakresuAgenta) konektorZakresu(ctx context.Context,
	kodAgenta, kodKonektora string) (KonektorAgenta, string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, konektorZakresuPoKodzie)
	if err != nil {
		return KonektorAgenta{}, "", err
	}
	return odczytajKonektorZakresu(polecenie.QueryRowContext(ctx, kodKonektora, kodAgenta))
}

// UsunUmiejetnoscAgenta zdejmuje umiejętność z definicji wskazanego eksperta
// i mówi, czy umiejętność istniała.
func (r *repozytoriumZakresuAgenta) UsunUmiejetnoscAgenta(ctx context.Context,
	kodAgenta, kodUmiejetnosci string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunUmiejetnoscZakresuAgenta)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kodUmiejetnosci, kodAgenta)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zdjąć umiejętności %q: %w", kodUmiejetnosci, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć zdjęć umiejętności %q: %w", kodUmiejetnosci, err)
	}
	return zdjete > 0, nil
}

// PrzypisaniaEkspertow oddaje przypisania widziane od strony eksperta,
// złożone z dwóch osobnych zapytań o inne źródła zamiast zapytania z UNION.
func (r *repozytoriumZakresuAgenta) PrzypisaniaEkspertow(ctx context.Context,
	kodAgenta, rodzaj string) ([]PrzypisanieEksperta, error) {

	przypisania := make([]PrzypisanieEksperta, 0, 8)
	if rodzaj == "" || rodzaj == RodzajPrzypisaniaProjekt {
		wiersze, err := r.przypisaniaProjektowe(ctx, kodAgenta)
		if err != nil {
			return nil, err
		}
		przypisania = append(przypisania, wiersze...)
	}
	if rodzaj == "" || rodzaj == RodzajPrzypisaniaRola {
		wiersze, err := r.przypisaniaRolowe(ctx, kodAgenta)
		if err != nil {
			return nil, err
		}
		przypisania = append(przypisania, wiersze...)
	}
	sort.SliceStable(przypisania, func(i, j int) bool {
		return przypisania[i].Przypisano < przypisania[j].Przypisano
	})
	return przypisania, nil
}

// RodzajPrzypisaniaProjekt i RodzajPrzypisaniaRola powtarzają wartości
// wyliczenia `AgentAssignmentKind` kontraktu, wpisane napisem, bo pakiet
// `dane` nie sięga po pakiet `shared`.
const (
	RodzajPrzypisaniaProjekt = "project"
	RodzajPrzypisaniaRola    = "role"
)

// przypisaniaProjektowe czyta przypisania eksperta do projektów zapisane
// przez Agent Managera modułu Workspace.
func (r *repozytoriumZakresuAgenta) przypisaniaProjektowe(ctx context.Context,
	kodAgenta string) ([]PrzypisanieEksperta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, przypisaniaProjektoweEkspertow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta, kodAgenta)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać przypisań projektowych: %w", err)
	}
	defer wiersze.Close()
	wynik := make([]PrzypisanieEksperta, 0, 8)
	for wiersze.Next() {
		wpis := PrzypisanieEksperta{Rodzaj: RodzajPrzypisaniaProjekt}
		var domyslny int
		if err := wiersze.Scan(&wpis.AgentKod, &wpis.CelKod, &wpis.CelNazwa,
			&wpis.Rola, &domyslny, &wpis.Przypisano); err != nil {
			return nil, fmt.Errorf("dane: nieczytelne przypisanie projektowe: %w", err)
		}
		wpis.DomyslnyWykonawca = domyslny == 1
		wynik = append(wynik, wpis)
	}
	return wynik, wiersze.Err()
}

// przypisaniaRolowe czyta przypisania eksperta do ról zapisane w obsadzie
// biegu orkiestracji wraz z nazwą biegu.
func (r *repozytoriumZakresuAgenta) przypisaniaRolowe(ctx context.Context,
	kodAgenta string) ([]PrzypisanieEksperta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, przypisaniaRolowaEkspertow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta, kodAgenta)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać przypisań rolowych: %w", err)
	}
	defer wiersze.Close()
	wynik := make([]PrzypisanieEksperta, 0, 8)
	for wiersze.Next() {
		wpis := PrzypisanieEksperta{Rodzaj: RodzajPrzypisaniaRola}
		var stanowisko int64
		if err := wiersze.Scan(&wpis.AgentKod, &stanowisko, &wpis.CelNazwa,
			&wpis.Rola, &wpis.Przypisano); err != nil {
			return nil, fmt.Errorf("dane: nieczytelne przypisanie rolowe: %w", err)
		}
		wpis.CelKod = fmt.Sprintf("obsada-%d", stanowisko)
		wynik = append(wynik, wpis)
	}
	return wynik, wiersze.Err()
}
