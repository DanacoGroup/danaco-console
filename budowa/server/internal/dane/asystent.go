// Odpowiedzialność pliku: obszar Assistant — definicja zlecenia asystenta
// (tabela `zlecenie_asystenta`, `store/migracja_050_asystent.sql`) wraz
// z kontraktem całego obszaru.
// Dziennik działań leży w `asystent_dziennik.go`, przyjęcie polecenia w
// `asystent_polecenia.go` — jedno repozytorium, trzy pliki wedle
// odpowiedzialności, tak jak `dane/library.go` i `dane/automations*.go`.
//
// INTERFEJS DEKLARUJE WYŁĄCZNIE TEN PLIK, w całości — wraz z metodami, które
// implementują pozostałe pliki obszaru. Interfejs rozdzielony na trzy pliki
// byłby trzema prawdami o jednym kontrakcie.
//
// CZAS JEST LICZBĄ, NIE NAPISEM. `utworzono`/`zaktualizowano` niosą milisekundy
// epoki wprost (kolumna INTEGER) — bez przekładu przez `strftime`, jak
// w modułach trzymających czas tekstem.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ZlecenieAsystenta to wiersz tabeli `zlecenie_asystenta`. Kod jest
// identyfikatorem, którym zlecenie wychodzi kontraktem (`AssistantAction.id`).
// Wartości `Stan` i `Droga` są wartościami kontraktu wprost (queued/running/
// paused/done/failed/cancelled — voice/text), bez tłumaczenia.
type ZlecenieAsystenta struct {
	ID           int64
	Kod          string
	OknoKod      string
	Tytul        *string
	Stan         string
	Droga        string
	EtapBiezacy  *int64
	LiczbaEtapow *int64
	Priorytet    *int64
	Wynik        *string
	// ProfilKod wskazuje profil asystenta, którym polecenie zostało wydane
	// (`store/migracja_117_profil_asystenta.sql`). NIE jest ozdobą wiersza:
	// wykonawca zlecenia biegnie
	// gorutyną długo po tym, jak komenda `assistant.voice.command` już
	// odpowiedziała, więc czyta warunki tury z BAZY, nie z żądania.
	// Profil trzymany tylko w pamięci procesu nie przetrwałby ani odpowiedzi na
	// komendę, ani restartu rdzenia w trakcie zlecenia. NULL znaczy zlecenie
	// bez profilu — tura idzie wtedy bez warstwy promptu.
	ProfilKod      *string
	Utworzono      int64
	Zaktualizowano int64
}

// RepozytoriumAsystenta jest kontraktem obszaru Assistant.
type RepozytoriumAsystenta interface {
	// --- agent A: zlecenie ---
	ZapiszZlecenie(ctx context.Context, zlecenie ZlecenieAsystenta) (ZlecenieAsystenta, error)
	Zlecenie(ctx context.Context, kod string) (ZlecenieAsystenta, error)
	Zlecenia(ctx context.Context, okno string) ([]ZlecenieAsystenta, error)
	UstawStanZlecenia(ctx context.Context, kod, stan string) (ZlecenieAsystenta, error)
	// UstawPriorytetZlecenia zmienia kolejność obsługi zlecenia w Actions
	// Monitor. Osobna metoda, bo `assistant.action.status` pozwala zmienić sam
	// priorytet BEZ zmiany stanu — sterowanie stanem i sterowanie kolejnością
	// to dwie różne czynności Operatora.
	UstawPriorytetZlecenia(ctx context.Context, kod string, priorytet int64) (ZlecenieAsystenta, error)
	// ZakonczZlecenie domyka zlecenie po wykonaniu: ustawia stan końcowy
	// (done/failed) i wynik w JEDNYM zapisie, żeby Actions Monitor nie zobaczył
	// stanu bez wyniku ani wyniku bez stanu. Zasila wykonawcę zleceń
	// (`core/adapter_modul_asystent_wykonawca.go`); kod nieznany wraca jako
	// ErrBrakWiersza — domknięcie bytu, którego nie ma, byłoby cichą zgodą.
	ZakonczZlecenie(ctx context.Context, kod, stan, wynik string) (ZlecenieAsystenta, error)

	// --- agent B: dziennik ---
	ZapiszWpis(ctx context.Context, wpis WpisDziennikaAsystenta) (WpisDziennikaAsystenta, error)
	Wpisy(ctx context.Context, okno string, limit int) ([]WpisDziennikaAsystenta, error)
	WpisyZlecenia(ctx context.Context, kodZlecenia string) ([]WpisDziennikaAsystenta, error)
	WpisDziennika(ctx context.Context, kod string) (WpisDziennikaAsystenta, error)
	OznaczWpisDziennika(ctx context.Context, kod string, wazny bool, notatka *string) (WpisDziennikaAsystenta, error)

	// --- agent C: polecenia ---
	PrzyjmijPolecenie(ctx context.Context, zlecenie ZlecenieAsystenta,
		wpis WpisDziennikaAsystenta) (ZlecenieAsystenta, WpisDziennikaAsystenta, error)

	// --- profil asystenta (`store/migracja_117_profil_asystenta.sql`, `asystent_profil.go`) ---
	// Sam odczyt: zakładania, wykazu ani wskazania domyślnego nie ma, bo nie ma
	// komendy kontraktu, która by je wołała (`assistant.profile.*` nie istnieje).
	// Metoda bez wołającego byłaby drogą, której nikt nie przechodzi.
	Profil(ctx context.Context, kod string) (ProfilAsystenta, error)
	ProfilDomyslny(ctx context.Context) (ProfilAsystenta, error)
}

const (
	kolumnyZleceniaAsystenta = `id, identyfikator_zewnetrzny, okno_kod, tytul, stan, droga,
	                            etap_biezacy, liczba_etapow, priorytet, wynik, profil_kod,
	                            utworzono, zaktualizowano`

	zapiszZlecenieAsystenta = `INSERT INTO zlecenie_asystenta
	                           (identyfikator_zewnetrzny, okno_kod, tytul, stan, droga,
	                            etap_biezacy, liczba_etapow, priorytet, wynik, profil_kod,
	                            utworzono, zaktualizowano)
	                           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                           ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                               okno_kod = excluded.okno_kod,
	                               tytul = excluded.tytul,
	                               stan = excluded.stan,
	                               droga = excluded.droga,
	                               etap_biezacy = excluded.etap_biezacy,
	                               liczba_etapow = excluded.liczba_etapow,
	                               priorytet = excluded.priorytet,
	                               wynik = excluded.wynik,
	                               profil_kod = excluded.profil_kod,
	                               zaktualizowano = excluded.zaktualizowano`

	pobierzZlecenieAsystenta = `SELECT ` + kolumnyZleceniaAsystenta + ` FROM zlecenie_asystenta
	                            WHERE identyfikator_zewnetrzny = ?`

	pobierzZleceniaOkna = `SELECT ` + kolumnyZleceniaAsystenta + ` FROM zlecenie_asystenta
	                       WHERE okno_kod = ?
	                       ORDER BY zaktualizowano DESC, id DESC`

	ustawStanZleceniaAsystenta = `UPDATE zlecenie_asystenta SET stan = ?, zaktualizowano = ?
	                              WHERE identyfikator_zewnetrzny = ?`
	ustawPriorytetZleceniaAsystenta = `UPDATE zlecenie_asystenta SET priorytet = ?, zaktualizowano = ?
	                                   WHERE identyfikator_zewnetrzny = ?`
	zakonczZlecenieAsystenta = `UPDATE zlecenie_asystenta SET stan = ?, wynik = ?, zaktualizowano = ?
	                            WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumAsystenta struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumAsystenta(z *zapytania, db *sql.DB) *repozytoriumAsystenta {
	return &repozytoriumAsystenta{zapytania: z, db: db}
}

// ZapiszZlecenie zakłada wiersz zlecenia albo nadpisuje zastany po kodzie
// zewnętrznym. `Utworzono` nie wchodzi do klauzuli UPDATE — zapis powtórny nie
// ma prawa przesunąć chwili założenia zlecenia, tylko chwilę ostatniej zmiany.
func (r *repozytoriumAsystenta) ZapiszZlecenie(ctx context.Context, zlecenie ZlecenieAsystenta) (ZlecenieAsystenta, error) {
	if zlecenie.Kod == "" {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: zlecenie asystenta bez identyfikatora")
	}
	if zlecenie.OknoKod == "" {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: zlecenie asystenta %q bez okna", zlecenie.Kod)
	}
	teraz := time.Now().UnixMilli()
	utworzono := zlecenie.Utworzono
	if utworzono == 0 {
		utworzono = teraz
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZlecenieAsystenta)
	if err != nil {
		return ZlecenieAsystenta{}, err
	}
	_, err = polecenie.ExecContext(ctx, zlecenie.Kod, zlecenie.OknoKod, tekstDoKolumny(zlecenie.Tytul),
		zlecenie.Stan, zlecenie.Droga, liczbaDoKolumny(zlecenie.EtapBiezacy),
		liczbaDoKolumny(zlecenie.LiczbaEtapow), liczbaDoKolumny(zlecenie.Priorytet),
		tekstDoKolumny(zlecenie.Wynik), tekstDoKolumny(zlecenie.ProfilKod), utworzono, teraz)
	if err != nil {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: nie można zapisać zlecenia asystenta %q: %w", zlecenie.Kod, err)
	}
	return r.Zlecenie(ctx, zlecenie.Kod)
}

// Zlecenie zwraca zlecenie o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się nie powiódł”.
func (r *repozytoriumAsystenta) Zlecenie(ctx context.Context, kod string) (ZlecenieAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZlecenieAsystenta)
	if err != nil {
		return ZlecenieAsystenta{}, err
	}
	zlecenie, err := odczytajZlecenieAsystenta(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ZlecenieAsystenta{}, ErrBrakWiersza
	}
	if err != nil {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: nieczytelny wiersz zlecenia asystenta %q: %w", kod, err)
	}
	return zlecenie, nil
}

// Zlecenia zwraca wszystkie zlecenia okna, od najświeższej zmiany.
func (r *repozytoriumAsystenta) Zlecenia(ctx context.Context, okno string) ([]ZlecenieAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZleceniaOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zleceń okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []ZlecenieAsystenta{}
	for wiersze.Next() {
		zlecenie, err := odczytajZlecenieAsystenta(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zlecenia asystenta: %w", err)
		}
		lista = append(lista, zlecenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zleceń okna %q: %w", okno, err)
	}
	return lista, nil
}

// UstawStanZlecenia zmienia stan zlecenia po kodzie zewnętrznym. Kod nieznany
// wraca jako ErrBrakWiersza — cicha zgoda na zmianę stanu bytu, którego nie ma,
// byłaby potwierdzeniem czynności, która się nie odbyła.
func (r *repozytoriumAsystenta) UstawStanZlecenia(ctx context.Context, kod, stan string) (ZlecenieAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, ustawStanZleceniaAsystenta)
	if err != nil {
		return ZlecenieAsystenta{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, stan, time.Now().UnixMilli(), kod)
	if err != nil {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: nie można ustawić stanu zlecenia asystenta %q: %w", kod, err)
	}
	dotknietych, err := wynik.RowsAffected()
	if err != nil {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: nieczytelny wynik zmiany stanu zlecenia %q: %w", kod, err)
	}
	if dotknietych == 0 {
		return ZlecenieAsystenta{}, ErrBrakWiersza
	}
	return r.Zlecenie(ctx, kod)
}

// UstawPriorytetZlecenia zmienia priorytet zlecenia po kodzie zewnętrznym.
// Osobno od stanu, bo `assistant.action.status` pozwala zmienić sam priorytet
// bez zmiany stanu.
//
// Kod nieznany wraca jako ErrBrakWiersza — cicha zgoda na zmianę priorytetu
// bytu, którego nie ma, byłaby potwierdzeniem czynności, która się nie odbyła.
func (r *repozytoriumAsystenta) UstawPriorytetZlecenia(ctx context.Context,
	kod string, priorytet int64) (ZlecenieAsystenta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawPriorytetZleceniaAsystenta)
	if err != nil {
		return ZlecenieAsystenta{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, priorytet, time.Now().UnixMilli(), kod)
	if err != nil {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: nie można ustawić priorytetu zlecenia asystenta %q: %w", kod, err)
	}
	dotknietych, err := wynik.RowsAffected()
	if err != nil {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: nieczytelny wynik zmiany priorytetu zlecenia %q: %w", kod, err)
	}
	if dotknietych == 0 {
		return ZlecenieAsystenta{}, ErrBrakWiersza
	}
	return r.Zlecenie(ctx, kod)
}

// ZakonczZlecenie ustawia stan końcowy i wynik zlecenia jednym zapisem. Wywołuje
// je wykonawca po turze modelu: `done` z odpowiedzią albo `failed` z powodem.
// Oba pola idą razem, bo pochodzą z jednego zdarzenia (zakończenia tury) — dwa
// osobne UPDATE-y pokazałyby przez chwilę stan bez pasującego wyniku.
//
// Kod nieznany wraca jako ErrBrakWiersza — domknięcie zlecenia, którego nie ma,
// nie może wyglądać jak sukces (spójnie z UstawStanZlecenia).
func (r *repozytoriumAsystenta) ZakonczZlecenie(ctx context.Context, kod, stan, wynik string) (ZlecenieAsystenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zakonczZlecenieAsystenta)
	if err != nil {
		return ZlecenieAsystenta{}, err
	}
	wynikZapisu, err := polecenie.ExecContext(ctx, stan, wynik, time.Now().UnixMilli(), kod)
	if err != nil {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: nie można domknąć zlecenia asystenta %q: %w", kod, err)
	}
	dotknietych, err := wynikZapisu.RowsAffected()
	if err != nil {
		return ZlecenieAsystenta{}, fmt.Errorf("dane: nieczytelny wynik domknięcia zlecenia asystenta %q: %w", kod, err)
	}
	if dotknietych == 0 {
		return ZlecenieAsystenta{}, ErrBrakWiersza
	}
	return r.Zlecenie(ctx, kod)
}

// odczytajZlecenieAsystenta składa strukturę z jednego wiersza wyniku.
func odczytajZlecenieAsystenta(wiersz skaner) (ZlecenieAsystenta, error) {
	var zlecenie ZlecenieAsystenta
	var tytul, wynik, profil sql.NullString
	var etapBiezacy, liczbaEtapow, priorytet sql.NullInt64
	err := wiersz.Scan(&zlecenie.ID, &zlecenie.Kod, &zlecenie.OknoKod, &tytul, &zlecenie.Stan,
		&zlecenie.Droga, &etapBiezacy, &liczbaEtapow, &priorytet, &wynik, &profil,
		&zlecenie.Utworzono, &zlecenie.Zaktualizowano)
	if err != nil {
		return ZlecenieAsystenta{}, err
	}
	zlecenie.Tytul = tekstZKolumny(tytul)
	zlecenie.Wynik = tekstZKolumny(wynik)
	zlecenie.ProfilKod = tekstZKolumny(profil)
	zlecenie.EtapBiezacy = liczbaZKolumny(etapBiezacy)
	zlecenie.LiczbaEtapow = liczbaZKolumny(liczbaEtapow)
	zlecenie.Priorytet = liczbaZKolumny(priorytet)
	return zlecenie, nil
}
