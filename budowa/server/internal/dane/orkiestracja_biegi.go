// Odpowiedzialność pliku: bieg orkiestracji MultitaskingAI (tabela
// `bieg_orkiestracji`) wraz z jego obsadą (tabela `obsada_biegu`).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// BiegOrkiestracji to wiersz tabeli `bieg_orkiestracji` — nazwany przepływ
// prowadzony w sesji, grupujący modele pracujące pod orkiestratorem.
type BiegOrkiestracji struct {
	ID                 int64
	Kod                string
	SesjaID            int64
	SesjaKod           *string
	OknoKoordynatoraID *int64
	KolejkaID          *int64
	Nazwa              string
	Stan               string
	Etap               *string
	EtapBiezacy        int
	Etapow             int
	Rozpoczeto         string
	Zakonczono         *string
}

// MiejsceObsady to wiersz tabeli `obsada_biegu` widziany od strony biegu:
// jedno z czterech stanowisk sceny wraz z jego wyposażeniem.
type MiejsceObsady struct {
	ID               int64
	BiegID           int64
	AgentID          *int64
	AgentKod         *string
	Model            *string
	Rola             string
	Miejsce          int
	Prompt           *string
	Narzedzia        *string
	ProfilIzolacjiID *int64
}

// RepozytoriumBiegow jest kontraktem obszaru biegów orkiestracji.
// Rdzeń bierze je metodą `Zestaw.BiegiOrkiestracji()`.
type RepozytoriumBiegow interface {
	ZapiszBieg(ctx context.Context, bieg BiegOrkiestracji) (BiegOrkiestracji, error)
	Bieg(ctx context.Context, kod string) (BiegOrkiestracji, error)
	BiegOknaKoordynatora(ctx context.Context, oknoID int64) (BiegOrkiestracji, error)
	ZapiszObsadeBiegu(ctx context.Context, biegID int64, obsada []MiejsceObsady) error
	ObsadaBiegu(ctx context.Context, biegID int64) ([]MiejsceObsady, error)
}

// Stany biegu — słownik zamknięty więzem CHECK schematu. `oczekuje` niesie
// bieg zawieszony na sygnał ze świata.
const (
	stanBieguOczekuje   = "pending"
	stanBieguWBiegu     = "running"
	stanBieguWstrzymany = "paused"
	stanBieguNaSygnal   = "oczekuje"
	stanBieguUdany      = "succeeded"
	stanBieguBledny     = "failed"
	stanBieguZatrzymany = "stopped"
)

const (
	kolumnyBiegu = `b.id, b.identyfikator_zewnetrzny, b.sesja_id, s.identyfikator_zewnetrzny,
	                b.okno_koordynatora_id, b.kolejka_id, b.nazwa, b.stan, b.etap,
	                b.etap_biezacy, b.etapow, b.rozpoczeto, b.zakonczono`

	zrodloBiegu = ` FROM bieg_orkiestracji b JOIN sesja s ON s.id = b.sesja_id`

	zapiszBiegOrkiestracji = `INSERT INTO bieg_orkiestracji
	                          (identyfikator_zewnetrzny, sesja_id, okno_koordynatora_id,
	                           kolejka_id, nazwa, stan, etap, etap_biezacy, etapow, zakonczono)
	                          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                          ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                              okno_koordynatora_id = excluded.okno_koordynatora_id,
	                              kolejka_id = excluded.kolejka_id,
	                              nazwa = excluded.nazwa,
	                              stan = excluded.stan,
	                              etap = excluded.etap,
	                              etap_biezacy = excluded.etap_biezacy,
	                              etapow = excluded.etapow,
	                              zakonczono = excluded.zakonczono`

	pobierzBieg = `SELECT ` + kolumnyBiegu + zrodloBiegu +
		` WHERE b.identyfikator_zewnetrzny = ?`

	// Bieg okna koordynatora: najświeższy, który jeszcze się nie domknął.
	// Biegi zamknięte są poza wykazem — nowo powołany podagent ma trafić do
	// pracy trwającej, a nie doczepić się do historii.
	pobierzBiegOkna = `SELECT ` + kolumnyBiegu + zrodloBiegu +
		` WHERE b.okno_koordynatora_id = ?
		    AND b.stan NOT IN ('succeeded','failed','stopped')
		  ORDER BY b.id DESC LIMIT 1`

	usunObsadeBiegu = `DELETE FROM obsada_biegu WHERE bieg_id = ?`

	wstawMiejsceBiegu = `INSERT INTO obsada_biegu
	                     (bieg_id, agent_id, model, rola, miejsce, prompt, narzedzia,
	                      profil_izolacji_id)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	listaObsadyBiegu = `SELECT o.id, o.bieg_id, o.agent_id, a.kod, o.model, o.rola,
	                           o.miejsce, o.prompt, o.narzedzia, o.profil_izolacji_id
	                      FROM obsada_biegu o
	                      LEFT JOIN agent a ON a.id = o.agent_id
	                     WHERE o.bieg_id = ?
	                     ORDER BY o.miejsce`
)

type repozytoriumBiegow struct {
	zapytania *zapytania
	db        *sql.DB
}

// Zgodność implementacji repozytorium biegów orkiestracji z kontraktem sprawdza kompilator, a nie montaż.
var _ RepozytoriumBiegow = (*repozytoriumBiegow)(nil)

// BiegiOrkiestracji oddaje repozytorium biegów nad pamięcią zapytań zestawu.
// Metoda, nie pole — tak samo jak `Zestaw.Rozszerzenia()`.
func (z *Zestaw) BiegiOrkiestracji() RepozytoriumBiegow {
	if z == nil || z.zapytania == nil {
		return nil
	}
	return &repozytoriumBiegow{zapytania: z.zapytania, db: z.zapytania.db}
}

// ZapiszBieg zakłada bieg orkiestracji albo nadpisuje zastany wiersz i oddaje jego pełny stan po zapisie.
func (r *repozytoriumBiegow) ZapiszBieg(ctx context.Context,
	bieg BiegOrkiestracji) (BiegOrkiestracji, error) {

	if bieg.Kod == "" || bieg.SesjaID == 0 || bieg.Nazwa == "" {
		return BiegOrkiestracji{}, fmt.Errorf("dane: bieg bez identyfikatora, sesji albo nazwy")
	}
	stan := bieg.Stan
	if stan == "" {
		stan = stanBieguOczekuje
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszBiegOrkiestracji)
	if err != nil {
		return BiegOrkiestracji{}, err
	}
	_, err = polecenie.ExecContext(ctx, bieg.Kod, bieg.SesjaID,
		liczbaDoKolumny(bieg.OknoKoordynatoraID), liczbaDoKolumny(bieg.KolejkaID),
		bieg.Nazwa, stan, tekstDoKolumny(bieg.Etap), bieg.EtapBiezacy, bieg.Etapow,
		tekstDoKolumny(bieg.Zakonczono))
	if err != nil {
		return BiegOrkiestracji{}, fmt.Errorf("dane: nie można zapisać biegu %q: %w", bieg.Kod, err)
	}
	return r.Bieg(ctx, bieg.Kod)
}

// Bieg zwraca bieg orkiestracji o wskazanym identyfikatorze zewnętrznym nadanym przez rdzeń całego systemu.
func (r *repozytoriumBiegow) Bieg(ctx context.Context, kod string) (BiegOrkiestracji, error) {
	return r.jedenBieg(ctx, pobierzBieg, kod, "bieg "+kod)
}

// BiegOknaKoordynatora zwraca bieg orkiestracji prowadzony przez wskazane okno komunikacji tego rdzenia.
func (r *repozytoriumBiegow) BiegOknaKoordynatora(ctx context.Context,
	oknoID int64) (BiegOrkiestracji, error) {

	if oknoID == 0 {
		return BiegOrkiestracji{}, ErrBrakWiersza
	}
	return r.jedenBieg(ctx, pobierzBiegOkna, oknoID,
		fmt.Sprintf("bieg okna %d", oknoID))
}

// ZapiszObsadeBiegu podmienia komplet obsady biegu; podmiana, a nie dopisywanie, bo obsada jest wykazem
// zamkniętym.
func (r *repozytoriumBiegow) ZapiszObsadeBiegu(ctx context.Context,
	biegID int64, obsada []MiejsceObsady) error {

	if biegID == 0 {
		return fmt.Errorf("dane: obsada bez wskazania biegu")
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunObsadeBiegu)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, biegID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić obsady biegu %d: %w", biegID, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawMiejsceBiegu)
		if err != nil {
			return err
		}
		for numer, stanowisko := range obsada {
			miejsce := stanowisko.Miejsce
			if miejsce == 0 {
				miejsce = numer + 1
			}
			rola := stanowisko.Rola
			if rola == "" {
				rola = "wykonawca"
			}
			_, err := wstawienie.ExecContext(ctx, biegID, liczbaDoKolumny(stanowisko.AgentID),
				tekstDoKolumny(stanowisko.Model), rola, miejsce,
				tekstDoKolumny(stanowisko.Prompt), tekstDoKolumny(stanowisko.Narzedzia),
				liczbaDoKolumny(stanowisko.ProfilIzolacjiID))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać miejsca %d obsady biegu %d: %w",
					miejsce, biegID, err)
			}
		}
		return nil
	})
}

// ObsadaBiegu zwraca stanowiska biegu orkiestracji w kolejności zajmowanych przez nie miejsc tej obsady.
func (r *repozytoriumBiegow) ObsadaBiegu(ctx context.Context,
	biegID int64) ([]MiejsceObsady, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaObsadyBiegu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, biegID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać obsady biegu %d: %w", biegID, err)
	}
	defer wiersze.Close()

	lista := []MiejsceObsady{}
	for wiersze.Next() {
		var m MiejsceObsady
		var agent, profil sql.NullInt64
		var agentKod, model, prompt, narzedzia sql.NullString
		err := wiersze.Scan(&m.ID, &m.BiegID, &agent, &agentKod, &model, &m.Rola,
			&m.Miejsce, &prompt, &narzedzia, &profil)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz obsady biegu: %w", err)
		}
		m.AgentID, m.ProfilIzolacjiID = liczbaZKolumny(agent), liczbaZKolumny(profil)
		m.AgentKod, m.Model = tekstZKolumny(agentKod), tekstZKolumny(model)
		m.Prompt, m.Narzedzia = tekstZKolumny(prompt), tekstZKolumny(narzedzia)
		lista = append(lista, m)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt obsady biegu: %w", err)
	}
	return lista, nil
}

// jedenBieg wykonuje odczyt pojedynczego wiersza biegu orkiestracji wspólny dla obu sposobów jego doboru.
func (r *repozytoriumBiegow) jedenBieg(ctx context.Context, zapytanie string,
	klucz any, opis string) (BiegOrkiestracji, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return BiegOrkiestracji{}, err
	}
	bieg, err := odczytajBieg(polecenie.QueryRowContext(ctx, klucz))
	if errors.Is(err, sql.ErrNoRows) {
		return BiegOrkiestracji{}, ErrBrakWiersza
	}
	if err != nil {
		return BiegOrkiestracji{}, fmt.Errorf("dane: nieczytelny %s: %w", opis, err)
	}
	return bieg, nil
}

// odczytajBieg składa pełną strukturę biegu orkiestracji z jednego wiersza wyniku zapytania do bazy danych.
func odczytajBieg(wiersz skaner) (BiegOrkiestracji, error) {
	var b BiegOrkiestracji
	var sesjaKod, etap, zakonczono sql.NullString
	var okno, kolejka sql.NullInt64
	err := wiersz.Scan(&b.ID, &b.Kod, &b.SesjaID, &sesjaKod, &okno, &kolejka, &b.Nazwa,
		&b.Stan, &etap, &b.EtapBiezacy, &b.Etapow, &b.Rozpoczeto, &zakonczono)
	if err != nil {
		return BiegOrkiestracji{}, err
	}
	b.SesjaKod = tekstZKolumny(sesjaKod)
	b.OknoKoordynatoraID, b.KolejkaID = liczbaZKolumny(okno), liczbaZKolumny(kolejka)
	b.Etap, b.Zakonczono = tekstZKolumny(etap), tekstZKolumny(zakonczono)
	return b, nil
}
