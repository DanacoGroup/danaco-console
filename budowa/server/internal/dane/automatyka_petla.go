// Odpowiedzialność pliku: pętla wykonawcza automatyki — obsada uczestników,
// wiązanie kroku z agentem portfolio i oczekiwanie biegu na sygnał z zewnątrz.
// Definicja, harmonogram i przebieg opisują automatykę; te trzy byty opisują
// jej obieg: kto bierze udział w pętli, który agent prowadzi krok i na co bieg
// czeka.
//
// Rozszerzenie idzie osobnym interfejsem `RepozytoriumPetli`, po który rdzeń
// sięga asercją typu na porcie automatyk — tak jak po `Harmonogramy`
// i `Orkiestracja` w `kompozycja.go`.
//
// Bieg oczekujący czeka na sygnał z zewnątrz, nie na zegar ani na Operatora,
// więc jego stan musi przeżyć restart rdzenia — stąd wiersz w bazie, a nie wpis
// w mapie adaptera.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// UczestnikPetli to wiersz tabeli `obsada_biegu` — jedno miejsce obsady.
//
// Tabela `obsada_biegu` niesie ten sam byt dla biegu automatyki i biegu
// orkiestracji, który wiersza w `automatyka` nie ma. Zapytania tego pliku niosą
// `automatyka_id`, więc widzą wyłącznie obsadę automatyk; obsadę orkiestracji
// czyta osobny plik.
type UczestnikPetli struct {
	ID           int64
	AutomatykaID int64
	AgentID      *int64
	AgentKod     *string
	Model        *string
	Rola         string
	Miejsce      int
}

// OczekiwanieBiegu to wiersz tabeli `oczekiwanie_biegu` — bieg zawieszony
// w konkretnym kroku, czekający na konkretny sygnał.
type OczekiwanieBiegu struct {
	ID              int64
	PrzebiegID      int64
	PrzebiegKod     string
	KrokZewnetrzny  string
	Sygnal          string
	Termin          *string
	PoTerminie      string
	Zalozono        string
	Wybudzono       *string
	PowodWybudzenia *string
	TrescSygnalu    *string
}

// RepozytoriumPetli — odczyt i zapis bytów pętli wykonawczej. Rdzeń bierze je
// asercją typu na porcie automatyk.
type RepozytoriumPetli interface {
	ZapiszObsade(ctx context.Context, automatykaID int64, obsada []UczestnikPetli) error
	Obsada(ctx context.Context, automatykaID int64) ([]UczestnikPetli, error)

	PrzypiszAgentaDoKroku(ctx context.Context, automatykaID int64, krokZewnetrzny string, agentID int64) error
	OdepnijAgentaOdKroku(ctx context.Context, automatykaID int64, krokZewnetrzny string) error
	AgenciKrokow(ctx context.Context, automatykaID int64) (map[string]int64, error)

	ZalozOczekiwanie(ctx context.Context, oczekiwanie OczekiwanieBiegu) (OczekiwanieBiegu, error)
	OczekiwanieCzynne(ctx context.Context, przebiegID int64) (OczekiwanieBiegu, error)
	OczekiwaniaNaSygnal(ctx context.Context, sygnal string) ([]OczekiwanieBiegu, error)
	OczekiwaniaPoTerminie(ctx context.Context, chwila string) ([]OczekiwanieBiegu, error)
	ZamknijOczekiwanie(ctx context.Context, id int64, powod string, tresc *string) error
}

const (
	usunObsade = `DELETE FROM obsada_biegu WHERE automatyka_id = ?`

	wstawUczestnika = `INSERT INTO obsada_biegu
	                   (automatyka_id, agent_id, model, rola, miejsce)
	                   VALUES (?, ?, ?, ?, ?)`

	listaObsady = `SELECT o.id, o.automatyka_id, o.agent_id, a.kod, o.model, o.rola, o.miejsce
	                 FROM obsada_biegu o
	                 LEFT JOIN agent a ON a.id = o.agent_id
	                WHERE o.automatyka_id = ?
	                ORDER BY o.miejsce`

	przypiszAgentaKroku = `INSERT INTO agent_kroku_automatyki
	                       (automatyka_id, krok_zewnetrzny, agent_id)
	                       VALUES (?, ?, ?)
	                       ON CONFLICT(automatyka_id, krok_zewnetrzny) DO UPDATE SET
	                           agent_id = excluded.agent_id`

	odepnijAgentaKroku = `DELETE FROM agent_kroku_automatyki
	                       WHERE automatyka_id = ? AND krok_zewnetrzny = ?`

	listaAgentowKrokow = `SELECT krok_zewnetrzny, agent_id
	                        FROM agent_kroku_automatyki WHERE automatyka_id = ?`

	kolumnyOczekiwania = `o.id, o.przebieg_id, p.identyfikator_zewnetrzny, o.krok_zewnetrzny,
	                      o.sygnal, o.termin, o.po_terminie, o.zalozono, o.wybudzono,
	                      o.powod_wybudzenia, o.tresc_sygnalu`

	zrodloOczekiwania = ` FROM oczekiwanie_biegu o
	                      JOIN przebieg_automatyki p ON p.id = o.przebieg_id`

	zalozOczekiwanie = `INSERT INTO oczekiwanie_biegu
	                    (przebieg_id, krok_zewnetrzny, sygnal, termin, po_terminie)
	                    VALUES (?, ?, ?, ?, ?)`

	pobierzOczekiwanieCzynne = `SELECT ` + kolumnyOczekiwania + zrodloOczekiwania +
		` WHERE o.przebieg_id = ? AND o.wybudzono IS NULL`

	pobierzOczekiwaniePoID = `SELECT ` + kolumnyOczekiwania + zrodloOczekiwania + ` WHERE o.id = ?`

	// Doręczenie sygnału. Biegi zamknięte są poza wykazem — sygnał trafia
	// wyłącznie do tych, które nadal czekają.
	listaNaSygnal = `SELECT ` + kolumnyOczekiwania + zrodloOczekiwania +
		` WHERE o.sygnal = ? AND o.wybudzono IS NULL ORDER BY o.id`

	// Budzik terminów. Oczekiwanie BEZ terminu nie wchodzi do wykazu — czekanie
	// bez zegara jest wyborem Operatora, nie zaległością do posprzątania.
	listaPoTerminie = `SELECT ` + kolumnyOczekiwania + zrodloOczekiwania +
		` WHERE o.wybudzono IS NULL AND o.termin IS NOT NULL AND o.termin <= ?
		  ORDER BY o.termin`

	zamknijOczekiwanie = `UPDATE oczekiwanie_biegu
	                         SET wybudzono = strftime('%Y-%m-%dT%H:%M:%fZ','now'),
	                             powod_wybudzenia = ?,
	                             tresc_sygnalu = ?
	                       WHERE id = ? AND wybudzono IS NULL`
)

// ZapiszObsade podmienia komplet obsady automatyki. Podmiana, a nie dopisywanie:
// obsada jest wykazem zamkniętym, a scalanie wierszy zostawiałoby uczestników,
// których Operator z niej usunął.
//
// Obsada pusta jest poprawna — automatyka bez obsady biegnie tak jak
// dotąd, na modelu wskazanym w kroku.
func (r *repozytoriumAutomatyk) ZapiszObsade(ctx context.Context,
	automatykaID int64, obsada []UczestnikPetli) error {

	if automatykaID == 0 {
		return fmt.Errorf("dane: obsada bez wskazania automatyki")
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunObsade)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, automatykaID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić obsady automatyki %d: %w", automatykaID, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawUczestnika)
		if err != nil {
			return err
		}
		for numer, uczestnik := range obsada {
			miejsce := uczestnik.Miejsce
			if miejsce == 0 {
				miejsce = numer + 1
			}
			_, err := wstawienie.ExecContext(ctx, automatykaID,
				liczbaDoKolumny(uczestnik.AgentID), tekstDoKolumny(uczestnik.Model),
				uczestnik.Rola, miejsce)
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać miejsca %d obsady automatyki %d: %w",
					miejsce, automatykaID, err)
			}
		}
		return nil
	})
}

// Obsada zwraca uczestników pętli w kolejności miejsc.
func (r *repozytoriumAutomatyk) Obsada(ctx context.Context, automatykaID int64) ([]UczestnikPetli, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaObsady)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać obsady automatyki %d: %w", automatykaID, err)
	}
	defer wiersze.Close()

	lista := []UczestnikPetli{}
	for wiersze.Next() {
		var uczestnik UczestnikPetli
		var agent sql.NullInt64
		var agentKod, model sql.NullString
		if err := wiersze.Scan(&uczestnik.ID, &uczestnik.AutomatykaID, &agent, &agentKod,
			&model, &uczestnik.Rola, &uczestnik.Miejsce); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz obsady: %w", err)
		}
		uczestnik.AgentID = liczbaZKolumny(agent)
		uczestnik.AgentKod, uczestnik.Model = tekstZKolumny(agentKod), tekstZKolumny(model)
		lista = append(lista, uczestnik)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt obsady: %w", err)
	}
	return lista, nil
}

// PrzypiszAgentaDoKroku wiąże krok z agentem portfolio. Wiązanie idzie po
// identyfikatorze zewnętrznym kroku, bo `ZapiszKroki` podmienia wiersze kroków
// w całości i klucz liczbowy by tego nie przeżył.
func (r *repozytoriumAutomatyk) PrzypiszAgentaDoKroku(ctx context.Context,
	automatykaID int64, krokZewnetrzny string, agentID int64) error {

	if automatykaID == 0 || krokZewnetrzny == "" || agentID == 0 {
		return fmt.Errorf("dane: wiązanie kroku z agentem bez automatyki, kroku albo agenta")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszAgentaKroku)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, automatykaID, krokZewnetrzny, agentID); err != nil {
		return fmt.Errorf("dane: nie można związać kroku %q z agentem %d: %w",
			krokZewnetrzny, agentID, err)
	}
	return nil
}

// OdepnijAgentaOdKroku zdejmuje wiązanie. Brak wiązania nie jest błędem —
// odpięcie kroku, który agenta nie miał, zostawia stan, o który chodziło.
func (r *repozytoriumAutomatyk) OdepnijAgentaOdKroku(ctx context.Context,
	automatykaID int64, krokZewnetrzny string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, odepnijAgentaKroku)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, automatykaID, krokZewnetrzny); err != nil {
		return fmt.Errorf("dane: nie można odpiąć agenta od kroku %q: %w", krokZewnetrzny, err)
	}
	return nil
}

// AgenciKrokow zwraca odwzorowanie krok → agent dla całej automatyki. Mapa
// zamiast wykazu, bo odbiorca pyta o krok, nie przegląda wierszy.
func (r *repozytoriumAutomatyk) AgenciKrokow(ctx context.Context,
	automatykaID int64) (map[string]int64, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaAgentowKrokow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać agentów kroków automatyki %d: %w",
			automatykaID, err)
	}
	defer wiersze.Close()

	wiazania := map[string]int64{}
	for wiersze.Next() {
		var krok string
		var agent int64
		if err := wiersze.Scan(&krok, &agent); err != nil {
			return nil, fmt.Errorf("dane: nieczytelne wiązanie kroku z agentem: %w", err)
		}
		wiazania[krok] = agent
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wiązań kroków: %w", err)
	}
	return wiazania, nil
}

// ZalozOczekiwanie zawiesza bieg na sygnał. Indeks częściowy dopuszcza jedno
// czynne oczekiwanie na bieg — bieg stoi w jednym miejscu, więc nie może czekać
// na dwie rzeczy naraz.
func (r *repozytoriumAutomatyk) ZalozOczekiwanie(ctx context.Context,
	oczekiwanie OczekiwanieBiegu) (OczekiwanieBiegu, error) {

	if oczekiwanie.PrzebiegID == 0 || oczekiwanie.KrokZewnetrzny == "" || oczekiwanie.Sygnal == "" {
		return OczekiwanieBiegu{}, fmt.Errorf("dane: oczekiwanie bez biegu, kroku albo sygnału")
	}
	poTerminie := oczekiwanie.PoTerminie
	if poTerminie == "" {
		poTerminie = "wznow"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zalozOczekiwanie)
	if err != nil {
		return OczekiwanieBiegu{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, oczekiwanie.PrzebiegID, oczekiwanie.KrokZewnetrzny,
		oczekiwanie.Sygnal, tekstDoKolumny(oczekiwanie.Termin), poTerminie)
	if err != nil {
		return OczekiwanieBiegu{}, fmt.Errorf("dane: nie można założyć oczekiwania biegu %d: %w",
			oczekiwanie.PrzebiegID, err)
	}
	id, err := wynik.LastInsertId()
	if err != nil {
		return OczekiwanieBiegu{}, fmt.Errorf("dane: nieznany numer założonego oczekiwania: %w", err)
	}
	return r.jednoOczekiwanie(ctx, pobierzOczekiwaniePoID, id, fmt.Sprintf("oczekiwanie %d", id))
}

// OczekiwanieCzynne zwraca oczekiwanie, na którym bieg stoi teraz.
// Bieg, który na nic nie czeka, wraca jako ErrBrakWiersza — to stan zwykły.
func (r *repozytoriumAutomatyk) OczekiwanieCzynne(ctx context.Context,
	przebiegID int64) (OczekiwanieBiegu, error) {

	return r.jednoOczekiwanie(ctx, pobierzOczekiwanieCzynne, przebiegID,
		fmt.Sprintf("czynne oczekiwanie biegu %d", przebiegID))
}

// OczekiwaniaNaSygnal zwraca biegi czekające na wskazany sygnał.
func (r *repozytoriumAutomatyk) OczekiwaniaNaSygnal(ctx context.Context,
	sygnal string) ([]OczekiwanieBiegu, error) {

	return r.wykazOczekiwan(ctx, listaNaSygnal, sygnal, "sygnał "+sygnal)
}

// OczekiwaniaPoTerminie zwraca oczekiwania, którym minął termin. To jest wykaz
// dla budzika: bieg czekający wiecznie jest wyciekiem, więc ktoś musi go
// policzyć.
func (r *repozytoriumAutomatyk) OczekiwaniaPoTerminie(ctx context.Context,
	chwila string) ([]OczekiwanieBiegu, error) {

	return r.wykazOczekiwan(ctx, listaPoTerminie, chwila, "termin "+chwila)
}

// ZamknijOczekiwanie odnotowuje wybudzenie. Warunek `wybudzono IS NULL`
// w poleceniu czyni je idempotentnym: sygnał doręczony dwa razy zamyka
// oczekiwanie raz, a drugie wywołanie nie nadpisze powodu ani chwili.
func (r *repozytoriumAutomatyk) ZamknijOczekiwanie(ctx context.Context,
	id int64, powod string, tresc *string) error {

	if powod == "" {
		return fmt.Errorf("dane: wybudzenie bez powodu — powód jest polem obowiązkowym")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zamknijOczekiwanie)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, powod, tekstDoKolumny(tresc), id); err != nil {
		return fmt.Errorf("dane: nie można zamknąć oczekiwania %d: %w", id, err)
	}
	return nil
}

// jednoOczekiwanie wykonuje odczyt pojedynczego wiersza wspólny wszystkim doborom.
func (r *repozytoriumAutomatyk) jednoOczekiwanie(ctx context.Context, zapytanie string,
	klucz any, opis string) (OczekiwanieBiegu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return OczekiwanieBiegu{}, err
	}
	oczekiwanie, err := odczytajOczekiwanie(polecenie.QueryRowContext(ctx, klucz))
	if errors.Is(err, sql.ErrNoRows) {
		return OczekiwanieBiegu{}, ErrBrakWiersza
	}
	if err != nil {
		return OczekiwanieBiegu{}, fmt.Errorf("dane: nieczytelne %s: %w", opis, err)
	}
	return oczekiwanie, nil
}

// wykazOczekiwan wykonuje odczyt wykazu wspólny obu doborom.
func (r *repozytoriumAutomatyk) wykazOczekiwan(ctx context.Context, zapytanie string,
	klucz any, opis string) ([]OczekiwanieBiegu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, klucz)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać oczekiwań na %s: %w", opis, err)
	}
	defer wiersze.Close()

	lista := []OczekiwanieBiegu{}
	for wiersze.Next() {
		oczekiwanie, err := odczytajOczekiwanie(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz oczekiwania: %w", err)
		}
		lista = append(lista, oczekiwanie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt oczekiwań: %w", err)
	}
	return lista, nil
}

// odczytajOczekiwanie składa strukturę z jednego wiersza wyniku.
func odczytajOczekiwanie(wiersz skaner) (OczekiwanieBiegu, error) {
	var o OczekiwanieBiegu
	var termin, wybudzono, powod, tresc sql.NullString
	err := wiersz.Scan(&o.ID, &o.PrzebiegID, &o.PrzebiegKod, &o.KrokZewnetrzny, &o.Sygnal,
		&termin, &o.PoTerminie, &o.Zalozono, &wybudzono, &powod, &tresc)
	if err != nil {
		return OczekiwanieBiegu{}, err
	}
	o.Termin, o.Wybudzono = tekstZKolumny(termin), tekstZKolumny(wybudzono)
	o.PowodWybudzenia, o.TrescSygnalu = tekstZKolumny(powod), tekstZKolumny(tresc)
	return o, nil
}
