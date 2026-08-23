package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"testing"

	"danacoconsole/shared"
)

// Skutek zleceń kolejki: czy za odpowiedzią rodziny `queue.item.*` leży wiersz
// tabeli `zlecenie_kolejki`, a za `queue.policy.set` — kolumny tabeli `kolejka`.
//
// Każdy sprawdzian schodzi do bazy własnym zapytaniem. Odpowiedzi tej rodziny
// są szczególnie łatwe do sfałszowania: `queue.item.enqueue` oddaje zlecenie
// złożone z pól żądania, więc rdzeń, który niczego nie zapisał, oddałby
// dokładnie ten sam kształt.

// kolejkaSprawdzianuZlecen zakłada kolejkę i oddaje jej identyfikator kontraktu
// wraz z kluczem wiersza.
func kolejkaSprawdzianuZlecen(t *testing.T, zmontowany *Zmontowany,
	zycie context.Context) (string, int64) {
	t.Helper()

	var wynik shared.QueueCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueCreate,
		shared.QueueCreateRequest{SessionId: "sesja-sprawdzianu"}, &wynik)

	// Rdzeń oddaje kolejkę pod kluczem jej wiersza, więc sprawdzian schodzi do
	// bazy wprost po nim — bez tego nie dałoby się zmierzyć niczego niezależnie.
	id, err := strconv.ParseInt(wynik.Queue.Id, 10, 64)
	if err != nil {
		t.Fatalf("identyfikator kolejki %q nie jest kluczem wiersza: %v", wynik.Queue.Id, err)
	}
	return wynik.Queue.Id, id
}

// TestSkutekDolozeniaZleceniaWBazie mierzy `queue.item.enqueue` wraz z kluczem
// idempotencji: powtórzone wywołanie nie ma zakładać drugiego wiersza.
func TestSkutekDolozeniaZleceniaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := kolejkaSprawdzianuZlecen(t, zmontowany, zycie)

	var pierwsze shared.QueueItemEnqueueResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemEnqueue,
		shared.QueueItemEnqueueRequest{
			QueueId: kod, Payload: json.RawMessage(`{"zadanie":"raport"}`),
			Priority: wskaznik(5), IdempotencyKey: wskaznik("klucz-jeden"),
		}, &pierwsze)
	if pierwsze.Duplicate {
		t.Fatalf("pierwsze dołożenie zameldowało duplikat")
	}

	var ladunek string
	var priorytet int
	var stan string
	err := baza.QueryRow(`SELECT ladunek, priorytet, stan FROM zlecenie_kolejki
	                      WHERE identyfikator_zewnetrzny = ?`, pierwsze.Item.Id).
		Scan(&ladunek, &priorytet, &stan)
	if err != nil {
		t.Fatalf("zlecenia nie ma w bazie po dołożeniu: %v", err)
	}
	if ladunek != `{"zadanie":"raport"}` || priorytet != 5 || stan != "oczekuje" {
		t.Fatalf("w bazie stoi ładunek %q, priorytet %d, stan %q", ladunek, priorytet, stan)
	}

	// Powtórzenie tego samego klucza ma oddać zlecenie zastane i NIE założyć
	// drugiego wiersza — wywołanie przychodzące powtórzone przez nadawcę nie
	// może wykonać pracy dwa razy.
	var drugie shared.QueueItemEnqueueResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemEnqueue,
		shared.QueueItemEnqueueRequest{
			QueueId: kod, Payload: json.RawMessage(`{"zadanie":"raport"}`),
			IdempotencyKey: wskaznik("klucz-jeden"),
		}, &drugie)
	if !drugie.Duplicate {
		t.Fatalf("powtórzenie klucza idempotencji nie zameldowało duplikatu")
	}
	if drugie.Item.Id != pierwsze.Item.Id {
		t.Fatalf("powtórzenie oddało inne zlecenie (%q zamiast %q)",
			drugie.Item.Id, pierwsze.Item.Id)
	}
	if zlecen := liczbaZlecenWBazie(t, baza, id); zlecen != 1 {
		t.Fatalf("po powtórzeniu w bazie stoi %d zleceń, oczekiwano 1", zlecen)
	}
}

// liczbaZlecenWBazie liczy zlecenia kolejki własnym zapytaniem.
func liczbaZlecenWBazie(t *testing.T, baza *sql.DB, kolejkaID int64) int {
	t.Helper()

	var zlecen int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM zlecenie_kolejki WHERE kolejka_id = ?`,
		kolejkaID).Scan(&zlecen); err != nil {
		t.Fatalf("nie można policzyć zleceń w bazie: %v", err)
	}
	return zlecen
}

// TestSkutekCykluZleceniaWBazie mierzy sześć czynności posuwających jedno
// zlecenie: zdjęcie, odłożenie, warunek, podział, scalenie i skierowanie.
func TestSkutekCykluZleceniaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := kolejkaSprawdzianuZlecen(t, zmontowany, zycie)

	odlozone := dolozZlecenieSprawdzianu(t, zmontowany, zycie, kod, `{"nr":1}`)
	zdejmowane := dolozZlecenieSprawdzianu(t, zmontowany, zycie, kod, `{"nr":2}`)
	dzielone := dolozZlecenieSprawdzianu(t, zmontowany, zycie, kod, `{"nr":3}`)

	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemDelay,
		shared.QueueItemDelayRequest{QueueId: kod, ItemId: odlozone, DelaySeconds: 600}, nil)
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemDequeue,
		shared.QueueItemDequeueRequest{QueueId: kod, ItemId: zdejmowane}, nil)
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemCondition,
		shared.QueueItemConditionRequest{
			QueueId: kod, ItemId: odlozone, Condition: "wynik.liczba > 0",
		}, nil)

	if stan, termin := stanZleceniaWBazie(t, baza, odlozone); stan != "odlozone" || termin == "" {
		t.Fatalf("odłożone zlecenie ma w bazie stan %q i termin %q", stan, termin)
	}
	if stan, _ := stanZleceniaWBazie(t, baza, zdejmowane); stan != "zdjete" {
		t.Fatalf("zdjęte zlecenie ma w bazie stan %q, oczekiwano „zdjete”", stan)
	}
	var warunek sql.NullString
	if err := baza.QueryRow(`SELECT warunek FROM zlecenie_kolejki
	                         WHERE identyfikator_zewnetrzny = ?`, odlozone).Scan(&warunek); err != nil {
		t.Fatalf("nie można odczytać warunku zlecenia: %v", err)
	}
	if !warunek.Valid || warunek.String != "wynik.liczba > 0" {
		t.Fatalf("w bazie stoi warunek %+v, oczekiwano zapisanego", warunek)
	}

	// Podział zakłada podzadania i ZAMYKA źródło. Sprawdzian liczy jedno i drugie.
	var podzial shared.QueueItemSplitResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemSplit,
		shared.QueueItemSplitRequest{
			QueueId: kod, ItemId: dzielone,
			Payloads: json.RawMessage(`[{"czesc":1},{"czesc":2}]`),
		}, &podzial)
	if len(podzial.Items) != 2 {
		t.Fatalf("podział oddał %d podzadań, oczekiwano 2", len(podzial.Items))
	}
	if stan, _ := stanZleceniaWBazie(t, baza, dzielone); stan != "zakonczone" {
		t.Fatalf("zlecenie źródłowe podziału ma w bazie stan %q, oczekiwano zamkniętego", stan)
	}
	for _, podzadanie := range podzial.Items {
		var zrodlowe sql.NullInt64
		if err := baza.QueryRow(`SELECT zlecenie_zrodlowe_id FROM zlecenie_kolejki
		                         WHERE identyfikator_zewnetrzny = ?`, podzadanie.Id).
			Scan(&zrodlowe); err != nil {
			t.Fatalf("podzadania nie ma w bazie: %v", err)
		}
		if !zrodlowe.Valid {
			t.Fatalf("podzadanie %q nie wskazuje zlecenia źródłowego", podzadanie.Id)
		}
	}

	// Scalenie zamyka źródła i zakłada jedno zlecenie.
	var scalenie shared.QueueItemMergeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemMerge,
		shared.QueueItemMergeRequest{
			QueueId: kod, ItemIds: []string{podzial.Items[0].Id, podzial.Items[1].Id},
		}, &scalenie)
	for _, zrodlowe := range podzial.Items {
		if stan, _ := stanZleceniaWBazie(t, baza, zrodlowe.Id); stan != "zakonczone" {
			t.Fatalf("zlecenie źródłowe scalenia ma stan %q, oczekiwano zamkniętego", stan)
		}
	}
	if stan, _ := stanZleceniaWBazie(t, baza, scalenie.Item.Id); stan != "oczekuje" {
		t.Fatalf("zlecenie scalone ma w bazie stan %q, oczekiwano „oczekuje”", stan)
	}

	// Skierowanie przenosi zlecenie do INNEJ kolejki — mierzone kolumną
	// `kolejka_id`, a nie odpowiedzią komendy.
	kodDocelowej, idDocelowej := kolejkaSprawdzianuZlecen(t, zmontowany, zycie)
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemRoute,
		shared.QueueItemRouteRequest{
			QueueId: kod, ItemId: scalenie.Item.Id, TargetQueueId: wskaznik(kodDocelowej),
		}, nil)
	var kolejkaPoSkierowaniu int64
	if err := baza.QueryRow(`SELECT kolejka_id FROM zlecenie_kolejki
	                         WHERE identyfikator_zewnetrzny = ?`, scalenie.Item.Id).
		Scan(&kolejkaPoSkierowaniu); err != nil {
		t.Fatalf("nie można odczytać kolejki zlecenia po skierowaniu: %v", err)
	}
	if kolejkaPoSkierowaniu != idDocelowej {
		t.Fatalf("zlecenie stoi w kolejce %d, oczekiwano %d (źródłowa: %d)",
			kolejkaPoSkierowaniu, idDocelowej, id)
	}
}

// dolozZlecenieSprawdzianu dokłada zlecenie i oddaje jego identyfikator.
func dolozZlecenieSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	kolejka, ladunek string) string {
	t.Helper()

	var wynik shared.QueueItemEnqueueResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemEnqueue,
		shared.QueueItemEnqueueRequest{QueueId: kolejka, Payload: json.RawMessage(ladunek)}, &wynik)
	return wynik.Item.Id
}

// stanZleceniaWBazie oddaje stan i termin zlecenia odczytane wprost z bazy.
func stanZleceniaWBazie(t *testing.T, baza *sql.DB, kod string) (string, string) {
	t.Helper()

	var stan string
	var termin sql.NullString
	err := baza.QueryRow(`SELECT stan, termin FROM zlecenie_kolejki
	                      WHERE identyfikator_zewnetrzny = ?`, kod).Scan(&stan, &termin)
	if err != nil {
		t.Fatalf("zlecenia %q nie ma w bazie: %v", kod, err)
	}
	return stan, termin.String
}

// TestSkutekRozgalezieniaIZadanMartwychWBazie mierzy rozgałęzienie na tory
// równoległe oraz odczyt zadań martwych — ten drugi po ręcznym przestawieniu
// stanu w bazie, bo rdzeń przenosi zlecenie w stan martwy dopiero po wyczerpaniu
// prób, a sprawdzian ma mierzyć ODCZYT, nie politykę ponawiania.
func TestSkutekRozgalezieniaIZadanMartwychWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := kolejkaSprawdzianuZlecen(t, zmontowany, zycie)
	kodDocelowej, idDocelowej := kolejkaSprawdzianuZlecen(t, zmontowany, zycie)

	zrodlowe := dolozZlecenieSprawdzianu(t, zmontowany, zycie, kod, `{"nr":1}`)

	var rozgalezienie shared.QueueItemBranchResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemBranch,
		shared.QueueItemBranchRequest{
			QueueId: kod, ItemId: zrodlowe,
			Branches: []shared.QueueBranch{
				{Condition: wskaznik("tor-a")},
				{TargetQueueId: wskaznik(kodDocelowej), Condition: wskaznik("tor-b")},
			},
		}, &rozgalezienie)
	if len(rozgalezienie.Items) != 2 {
		t.Fatalf("rozgałęzienie oddało %d torów, oczekiwano 2", len(rozgalezienie.Items))
	}
	// Tor bez kolejki docelowej zostaje w źródłowej, tor z nią — przechodzi.
	if zlecen := liczbaZlecenWBazie(t, baza, id); zlecen != 2 {
		t.Fatalf("kolejka źródłowa ma w bazie %d zleceń, oczekiwano 2", zlecen)
	}
	if zlecen := liczbaZlecenWBazie(t, baza, idDocelowej); zlecen != 1 {
		t.Fatalf("kolejka docelowa toru ma w bazie %d zleceń, oczekiwano 1", zlecen)
	}

	if _, err := baza.Exec(`UPDATE zlecenie_kolejki SET stan = 'martwe'
	                        WHERE identyfikator_zewnetrzny = ?`, zrodlowe); err != nil {
		t.Fatalf("nie można przestawić stanu zlecenia w bazie: %v", err)
	}
	var martwe shared.QueueDeadListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueDeadList,
		shared.QueueDeadListRequest{}, &martwe)
	if len(martwe.Items) != 1 || martwe.Items[0].Id != zrodlowe {
		t.Fatalf("wykaz zadań martwych oddał %d pozycji, oczekiwano zlecenia %q",
			len(martwe.Items), zrodlowe)
	}
}

// TestSkutekPolitykiIWykazuZlecenWBazie mierzy `queue.policy.set` w kolumnach
// tabeli `kolejka` oraz `queue.item.list` wraz z licznikiem wszystkich zleceń.
func TestSkutekPolitykiIWykazuZlecenWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := kolejkaSprawdzianuZlecen(t, zmontowany, zycie)

	zasiegProjektu := shared.QueueScope(shared.QueueScopeProject)
	wycofanieStale := shared.QueueBackoffKind(shared.QueueBackoffKindFixed)
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueuePolicySet,
		shared.QueuePolicySetRequest{
			QueueId: kod,
			Policy: shared.QueuePolicy{
				Scope: &zasiegProjektu, ScopeId: wskaznik("proj-1"),
				MaxConcurrent: wskaznik(4), MaxAttempts: wskaznik(7),
				Backoff: &wycofanieStale, Jitter: wskaznik(false),
			},
		}, nil)

	var zasieg, zasiegID, wycofanie string
	var rownoleglych, prob, rozproszenie int
	err := baza.QueryRow(`SELECT zasieg, COALESCE(zasieg_id,''), limit_rownoleglych,
	                             limit_prob, wycofanie, rozproszenie
	                      FROM kolejka WHERE id = ?`, id).
		Scan(&zasieg, &zasiegID, &rownoleglych, &prob, &wycofanie, &rozproszenie)
	if err != nil {
		t.Fatalf("nie można odczytać polityki kolejki z bazy: %v", err)
	}
	if zasieg != "projektu" || zasiegID != "proj-1" || rownoleglych != 4 || prob != 7 {
		t.Fatalf("w bazie stoi polityka %q/%q, %d równoległych, %d prób",
			zasieg, zasiegID, rownoleglych, prob)
	}
	if wycofanie != string(shared.QueueBackoffKindFixed) || rozproszenie != 0 {
		t.Fatalf("w bazie stoi wycofanie %q i rozproszenie %d", wycofanie, rozproszenie)
	}

	// Pole nieobecne w żądaniu ma zostawić wartość zastaną, a nie wyzerować
	// dziewięciu pozostałych: kontrakt ma tu same pola opcjonalne.
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueuePolicySet,
		shared.QueuePolicySetRequest{
			QueueId: kod, Policy: shared.QueuePolicy{RatePerMinute: wskaznik(30)},
		}, nil)
	var tempo int
	err = baza.QueryRow(`SELECT limit_prob, tempo_na_minute FROM kolejka WHERE id = ?`, id).
		Scan(&prob, &tempo)
	if err != nil {
		t.Fatalf("nie można odczytać polityki po zmianie cząstkowej: %v", err)
	}
	if prob != 7 || tempo != 30 {
		t.Fatalf("zmiana cząstkowa dała %d prób i tempo %d, oczekiwano 7 i 30", prob, tempo)
	}

	dolozZlecenieSprawdzianu(t, zmontowany, zycie, kod, `{"nr":1}`)
	dolozZlecenieSprawdzianu(t, zmontowany, zycie, kod, `{"nr":2}`)

	var wykaz shared.QueueItemListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueItemList,
		shared.QueueItemListRequest{QueueId: kod}, &wykaz)
	if wykaz.Total != 2 || len(wykaz.Items) != 2 {
		t.Fatalf("wykaz zleceń oddał %d pozycji przy liczniku %d, oczekiwano 2 i 2",
			len(wykaz.Items), wykaz.Total)
	}
	if zlecen := liczbaZlecenWBazie(t, baza, id); zlecen != 2 {
		t.Fatalf("w bazie stoi %d zleceń, a wykaz oddał 2 — obraz rozjechał się ze stanem", zlecen)
	}

	// Głębokość liczy się z bazy i ma zobaczyć oba zlecenia oczekujące.
	var glebokosc shared.QueueDepthGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandQueueDepthGet,
		shared.QueueDepthGetRequest{QueueId: wskaznik(kod), FromAt: 0, BucketSeconds: 86400},
		&glebokosc)
	oczekujacych := 0
	for _, punkt := range glebokosc.Points {
		oczekujacych += punkt.Pending
	}
	if oczekujacych != 2 {
		t.Fatalf("głębokość kolejki naliczyła %d zleceń oczekujących, oczekiwano 2", oczekujacych)
	}
}
