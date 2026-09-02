// Odpowiedzialność pliku: zakres działania eksperta — moduły zastosowania,
// osiem zakresów izolacji technicznej, granica Subagent Network, zdjęcie
// wpisów uprawnień oraz odczyt konektorów i przypisań widziany od strony
// eksperta.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type PrzelacznikIzolacjiAgenta struct {
	Zakres  string
	Odciety bool
}

type PrzypisanieEksperta struct {
	AgentKod          string
	Rodzaj            string
	CelKod            string
	CelNazwa          string
	Rola              string
	DomyslnyWykonawca bool
	Przypisano        string
}

type RepozytoriumZakresuAgenta interface {
	ModulyAgenta(ctx context.Context, kodAgenta string) ([]string, error)
	UstawModulyAgenta(ctx context.Context, kodAgenta string, kody []string) error
	IzolacjaAgenta(ctx context.Context, kodAgenta string) ([]PrzelacznikIzolacjiAgenta, error)
	UstawIzolacjeAgenta(ctx context.Context, kodAgenta string, przelaczniki []PrzelacznikIzolacjiAgenta) error
	GranicaPodagentow(ctx context.Context, kodAgenta string) (int, error)
	UstawGranicePodagentow(ctx context.Context, kodAgenta string, granica int) error
	UsunUprawnieniaAgenta(ctx context.Context, kodAgenta, grupa, zakres string) (int, error)
	KonektoryAgenta(ctx context.Context, kodAgenta string) ([]KonektorAgenta, []string, error)
	UsunKonektorAgenta(ctx context.Context, kodAgenta, kodKonektora string) (bool, error)
	ZapiszKonektorAgenta(ctx context.Context, kodAgenta, kodKonektora string,
		punktID *int64, konfiguracja *string, aktywny *bool) (KonektorAgenta, string, error)
	UsunUmiejetnoscAgenta(ctx context.Context, kodAgenta, kodUmiejetnosci string) (bool, error)
	PrzypisaniaEkspertow(ctx context.Context, kodAgenta, rodzaj string) ([]PrzypisanieEksperta, error)
}

const (
	kolumnyKonektoraZakresu = `k.id, k.kod, k.agent_id, a.kod, k.nazwa, k.rodzaj,
	                           k.punkt_dostepu_id, k.konfiguracja, k.aktywny, k.utworzono,
	                           COALESCE(p.kod, '')`

	czyscModulyZakresuAgenta = `DELETE FROM agent_modul_zastosowania WHERE agent_id = ?`

	wstawModulZakresuAgenta = `INSERT INTO agent_modul_zastosowania (agent_id, kod_modulu)
	                           VALUES (?, ?) ON CONFLICT(agent_id, kod_modulu) DO NOTHING`

	numerEksperta = `SELECT id FROM agent WHERE kod = ? AND ` + WarunekKonta

	ekspertIstnieje = `SELECT 1 FROM agent WHERE kod = ?`

	granicaPodagentowAgenta = `SELECT limit_podagentow FROM agent WHERE kod = ? AND ` + WarunekKonta

	zapiszGranicePodagentow = `UPDATE agent SET limit_podagentow = ?,
	                               zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                            WHERE kod = ? AND ` + WarunekKonta

	usunUprawnieniaZakresuAgenta = `DELETE FROM agent_uprawnienie
	                                 WHERE agent_id = ?
	                                   AND (? = '' OR grupa = ?)
	                                   AND (? = '' OR zakres = ?)`

	// Zawężenie stoi w zapytaniu podrzędnym, bo `agent` i `punkt_dostepu` mają
	// obie kolumnę `konto_id`; punkt innego konta wchodzi do złączenia bez kodu.
	punktKonektoraZakresu = ` LEFT JOIN (SELECT id, kod FROM punkt_dostepu
	                                      WHERE ` + WarunekKonta + `) p
	                                 ON p.id = k.punkt_dostepu_id`

	usunKonektorZakresuAgenta = `DELETE FROM agent_konektor
	                              WHERE kod = ? AND agent_id = (SELECT id FROM agent
	                                                            WHERE kod = ? AND ` + WarunekKonta + `)`

	usunUmiejetnoscZakresuAgenta = `DELETE FROM agent_umiejetnosc
	                                 WHERE kod = ? AND agent_id = (SELECT id FROM agent
	                                                               WHERE kod = ? AND ` + WarunekKonta + `)`
)

var (
	warunekKontaEksperta = strings.ReplaceAll(WarunekKonta, "konto_id", "a.konto_id")
	warunekKontaProjektu = strings.ReplaceAll(WarunekKonta, "konto_id", "p.konto_id")

	modulyZakresuAgenta = `SELECT m.kod_modulu FROM agent_modul_zastosowania m
	                         JOIN agent a ON a.id = m.agent_id
	                        WHERE a.kod = ? AND ` + warunekKontaEksperta + `
	                        ORDER BY m.kod_modulu`

	izolacjaZakresuAgenta = `SELECT i.zakres, i.odciety FROM agent_izolacja_techniczna i
	                           JOIN agent a ON a.id = i.agent_id
	                          WHERE a.kod = ? AND ` + warunekKontaEksperta + `
	                          ORDER BY i.zakres`

	// agent_izolacja_techniczna nie ma konto_id; granica idzie przez korzeń agent.
	zapiszIzolacjeZakresuAgenta = `INSERT INTO agent_izolacja_techniczna (agent_id, zakres, odciety)
	                               VALUES (?, ?, ?)
	                               ON CONFLICT(agent_id, zakres) DO UPDATE SET
	                                   odciety = excluded.odciety,
	                                   zapisano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                               WHERE EXISTS (SELECT 1 FROM agent a
	                                             WHERE a.id = agent_izolacja_techniczna.agent_id
	                                               AND ` + warunekKontaEksperta + `)`

	konektoryZakresuAgenta = `SELECT ` + kolumnyKonektoraZakresu + `
	                            FROM agent_konektor k
	                            JOIN agent a ON a.id = k.agent_id` + punktKonektoraZakresu + `
	                           WHERE a.kod = ? AND ` + warunekKontaEksperta + `
	                           ORDER BY k.utworzono, k.id`

	konektorZakresuPoKodzie = `SELECT ` + kolumnyKonektoraZakresu + `
	                             FROM agent_konektor k
	                             JOIN agent a ON a.id = k.agent_id` + punktKonektoraZakresu + `
	                            WHERE k.kod = ? AND a.kod = ? AND ` + warunekKontaEksperta

	// Przypisanie do projektu: Agent Manager modułu Workspace (migracja 035); projekt niesie konto_id od migracji 484.
	przypisaniaProjektoweEkspertow = `SELECT pa.agent_kod, p.kod, p.nazwa,
	                                         COALESCE(pa.rola, ''), pa.domyslny_wykonawca, pa.przypisano
	                                    FROM przypisanie_agenta_projektu pa
	                                    JOIN projekt p ON p.id = pa.projekt_id
	                                   WHERE (? = '' OR pa.agent_kod = ?) AND ` + warunekKontaProjektu + `
	                                   ORDER BY pa.przypisano, p.kod`

	// Przypisanie do roli: stanowisko obsady biegu orkiestracji (migracja 077).
	przypisaniaRolowaEkspertow = `SELECT a.kod, o.id, COALESCE(b.nazwa, ''), o.rola, o.utworzono
	                                FROM obsada_biegu o
	                                JOIN agent a ON a.id = o.agent_id
	                                LEFT JOIN bieg_orkiestracji b ON b.id = o.bieg_id
	                               WHERE o.agent_id IS NOT NULL AND (? = '' OR a.kod = ?)
	                                 AND ` + warunekKontaEksperta + `
	                               ORDER BY o.utworzono, o.id`
)

type repozytoriumZakresuAgenta struct {
	zapytania *zapytania
	db        *sql.DB
}

var _ RepozytoriumZakresuAgenta = (*repozytoriumZakresuAgenta)(nil)

func noweRepozytoriumZakresuAgenta(z *zapytania, db *sql.DB) *repozytoriumZakresuAgenta {
	return &repozytoriumZakresuAgenta{zapytania: z, db: db}
}

// numerEkspertaZakresu przekłada kod eksperta na klucz wiersza konta Operatora.
func (r *repozytoriumZakresuAgenta) numerEkspertaZakresu(ctx context.Context, kod string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, numerEksperta)
	if err != nil {
		return 0, err
	}
	var id int64
	switch err := polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&id); {
	case err == sql.ErrNoRows:
		return 0, odmowaEksperta(ctx, r.zapytania, kod)
	case err != nil:
		return 0, fmt.Errorf("dane: nie można odczytać eksperta %q: %w", kod, err)
	}
	return id, nil
}

// odmowaEksperta odróżnia eksperta innego konta (ErrKolizjaWiersza) od kodu nieznanego (ErrBrakWiersza).
func odmowaEksperta(ctx context.Context, z *zapytania, kod string) error {
	polecenie, err := z.przygotuj(ctx, ekspertIstnieje)
	if err != nil {
		return err
	}
	var jest int
	switch err := polecenie.QueryRowContext(ctx, kod).Scan(&jest); {
	case err == sql.ErrNoRows:
		return ErrBrakWiersza
	case err != nil:
		return fmt.Errorf("dane: nie można sprawdzić eksperta %q: %w", kod, err)
	}
	return fmt.Errorf("dane: ekspert %q należy do innego konta: %w", kod, ErrKolizjaWiersza)
}

func (r *repozytoriumZakresuAgenta) ModulyAgenta(ctx context.Context, kodAgenta string) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, modulyZakresuAgenta)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta, KontoOperatora(ctx))
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

func (r *repozytoriumZakresuAgenta) IzolacjaAgenta(ctx context.Context,
	kodAgenta string) ([]PrzelacznikIzolacjiAgenta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, izolacjaZakresuAgenta)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta, KontoOperatora(ctx))
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
		wynik, err := transakcja.ExecContext(ctx, zapiszIzolacjeZakresuAgenta,
			id, zakres, liczbaLogiczna(przelacznik.Odciety), KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać zakresu %q eksperta %q: %w", zakres, kodAgenta, err)
		}
		if err := sprawdzTrafienieZapisu(wynik, "zakres izolacji eksperta", zakres); err != nil {
			return err
		}
	}
	if err := transakcja.Commit(); err != nil {
		return fmt.Errorf("dane: nie można domknąć zapisu izolacji eksperta %q: %w", kodAgenta, err)
	}
	return nil
}

func (r *repozytoriumZakresuAgenta) GranicaPodagentow(ctx context.Context, kodAgenta string) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, granicaPodagentowAgenta)
	if err != nil {
		return 0, err
	}
	var granica int
	switch err := polecenie.QueryRowContext(ctx, kodAgenta, KontoOperatora(ctx)).Scan(&granica); {
	case err == sql.ErrNoRows:
		return 0, ErrBrakWiersza
	case err != nil:
		return 0, fmt.Errorf("dane: nie można odczytać granicy podagentów eksperta %q: %w", kodAgenta, err)
	}
	return granica, nil
}

func (r *repozytoriumZakresuAgenta) UstawGranicePodagentow(ctx context.Context,
	kodAgenta string, granica int) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszGranicePodagentow)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, granica, kodAgenta, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać granicy podagentów eksperta %q: %w", kodAgenta, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err == nil && zmienione == 0 {
		return odmowaEksperta(ctx, r.zapytania, kodAgenta)
	}
	return nil
}

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

func (r *repozytoriumZakresuAgenta) KonektoryAgenta(ctx context.Context,
	kodAgenta string) ([]KonektorAgenta, []string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, konektoryZakresuAgenta)
	if err != nil {
		return nil, nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, KontoOperatora(ctx), kodAgenta, KontoOperatora(ctx))
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

type czytelnikWierszaZakresu interface {
	Scan(cele ...any) error
}

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
	wynik, err := polecenie.ExecContext(ctx, kodKonektora, kodAgenta, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można odłączyć konektora %q: %w", kodKonektora, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć odłączeń konektora %q: %w", kodKonektora, err)
	}
	return zdjete > 0, r.kolizjaEksperta(ctx, zdjete, kodAgenta)
}

// kolizjaEksperta oddaje ErrKolizjaWiersza, gdy zapis niczego nie zdjął, a ekspert o tym kodzie stoi na innym koncie; brak eksperta nie jest tu odmową.
func (r *repozytoriumZakresuAgenta) kolizjaEksperta(ctx context.Context, zdjete int64, kodAgenta string) error {
	if zdjete > 0 {
		return nil
	}
	err := odmowaEksperta(ctx, r.zapytania, kodAgenta)
	if errors.Is(err, ErrBrakWiersza) {
		return nil
	}
	return err
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

func (r *repozytoriumZakresuAgenta) konektorZakresu(ctx context.Context,
	kodAgenta, kodKonektora string) (KonektorAgenta, string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, konektorZakresuPoKodzie)
	if err != nil {
		return KonektorAgenta{}, "", err
	}
	return odczytajKonektorZakresu(polecenie.QueryRowContext(ctx, KontoOperatora(ctx),
		kodKonektora, kodAgenta, KontoOperatora(ctx)))
}

func (r *repozytoriumZakresuAgenta) UsunUmiejetnoscAgenta(ctx context.Context,
	kodAgenta, kodUmiejetnosci string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunUmiejetnoscZakresuAgenta)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kodUmiejetnosci, kodAgenta, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można zdjąć umiejętności %q: %w", kodUmiejetnosci, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć zdjęć umiejętności %q: %w", kodUmiejetnosci, err)
	}
	return zdjete > 0, r.kolizjaEksperta(ctx, zdjete, kodAgenta)
}

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

func (r *repozytoriumZakresuAgenta) przypisaniaProjektowe(ctx context.Context,
	kodAgenta string) ([]PrzypisanieEksperta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, przypisaniaProjektoweEkspertow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta, kodAgenta, KontoOperatora(ctx))
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

func (r *repozytoriumZakresuAgenta) przypisaniaRolowe(ctx context.Context,
	kodAgenta string) ([]PrzypisanieEksperta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, przypisaniaRolowaEkspertow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta, kodAgenta, KontoOperatora(ctx))
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
