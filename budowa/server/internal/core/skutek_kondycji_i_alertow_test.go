// Sprawdza skutek rodzin kondycji i alertów: każdy test odczytuje bezpośrednio
// wiersze tabel sonda_kondycji, wynik_sondy_kondycji, regula_alertu
// i wyzwolenie_alertu, nie tylko odpowiedź komendy.
package core

import (
	"context"
	"database/sql"
	"testing"

	"danacoconsole/shared"
)

// TestSondaWewnetrznaZostawiaZmierzonyWiersz wykazuje, że przebieg sondy
// naprawdę czegoś dotknął: w tabeli serii leży wiersz z niezerowym czasem
// odpowiedzi i ze zdaniem mówiącym, co zmierzono.
func TestSondaWewnetrznaZostawiaZmierzonyWiersz(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var zapis shared.HealthProbeSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeSave,
		shared.HealthProbeSaveRequest{
			Name:       "baza stanu rdzenia",
			Kind:       shared.HealthProbeKindInternal,
			Target:     celSondyWewnetrznejBaza,
			IntervalMs: 60_000,
		}, &zapis)
	if !zapis.Created {
		t.Fatal("pierwszy zapis sondy nie zameldował założenia nowej")
	}

	var przebieg shared.HealthProbeRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeRun,
		shared.HealthProbeRunRequest{ProbeId: zapis.Probe.Id}, &przebieg)

	if przebieg.Result.Status != shared.HealthProbeStatusUp {
		t.Fatalf("sonda bazy oddała stan %q, a baza sprawdzianu działa", przebieg.Result.Status)
	}

	// Świadek niezależny: wiersz serii wraz z tym, co naprawdę zmierzono.
	var stan, szczegol string
	var wykonano int64
	var czas sql.NullInt64
	err := baza.QueryRow(
		`SELECT stan, COALESCE(szczegol, ''), wykonano, czas_odpowiedzi_ms
		   FROM wynik_sondy_kondycji WHERE sonda_kod = ?`, zapis.Probe.Id).
		Scan(&stan, &szczegol, &wykonano, &czas)
	if err != nil {
		t.Fatalf("po przebiegu sondy nie ma wiersza serii — pomiar, którego dostępność nie zobaczy: %v", err)
	}
	if stan != shared.HealthProbeStatusUp {
		t.Fatalf("wiersz serii niesie stan %q, a odpowiedź meldowała sprawność", stan)
	}
	if szczegol == "" {
		t.Fatal("wiersz serii nie mówi, co zmierzono — liczba bez świadka")
	}
	if wykonano <= 0 {
		t.Fatal("wiersz serii nie ma chwili pomiaru")
	}
	if !czas.Valid {
		t.Fatal("wiersz serii nie ma czasu odpowiedzi — pomiar bez pomiaru")
	}

	var odbicieStanu string
	var odbicieCzasu int64
	if err := baza.QueryRow(
		`SELECT COALESCE(ostatni_stan, ''), COALESCE(ostatni_przebieg, 0)
		   FROM sonda_kondycji WHERE identyfikator_zewnetrzny = ?`, zapis.Probe.Id).
		Scan(&odbicieStanu, &odbicieCzasu); err != nil {
		t.Fatalf("nie można odczytać definicji sondy: %v", err)
	}
	// Odbicie w definicji ma się zgadzać z ostatnim wierszem serii.
	if odbicieStanu != stan || odbicieCzasu != wykonano {
		t.Fatalf("wykaz sond pokazałby stan %q z chwili %d, a seria niesie %q z chwili %d",
			odbicieStanu, odbicieCzasu, stan, wykonano)
	}
}

// TestSondaGniazdaBezNasluchuJestNiesprawna wykazuje, że pomiar naprawdę
// sięgnął sieci: port, na którym nikt nie słucha, daje stan down, a nie up.
// Adapter, który zawsze oddaje up, przeszedłby poprzedni sprawdzian bez trudu.
func TestSondaGniazdaBezNasluchuJestNiesprawna(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var zapis shared.HealthProbeSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeSave,
		shared.HealthProbeSaveRequest{
			Name: "gniazdo, na którym nikt nie słucha",
			Kind: shared.HealthProbeKindTcp,
			// Port 1 na pętli zwrotnej: numer zarezerwowany, na którym nie nasłuchuje nic.
			Target:     "127.0.0.1:1",
			IntervalMs: 60_000,
			TimeoutMs:  wskaznik(int64(1000)),
		}, &zapis)

	var przebieg shared.HealthProbeRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeRun,
		shared.HealthProbeRunRequest{ProbeId: zapis.Probe.Id}, &przebieg)

	if przebieg.Result.Status == shared.HealthProbeStatusUp {
		t.Fatal("sonda gniazda bez nasłuchu oddała sprawność — to fasada, " +
			"przez którą awaria przechodzi niezauważona")
	}

	var stan string
	if err := baza.QueryRow(`SELECT stan FROM wynik_sondy_kondycji WHERE sonda_kod = ?`,
		zapis.Probe.Id).Scan(&stan); err != nil {
		t.Fatalf("nie ma wiersza serii nieudanego pomiaru: %v", err)
	}
	if stan == shared.HealthProbeStatusUp {
		t.Fatalf("wiersz serii niesie stan %q dla gniazda, które nie przyjęło połączenia", stan)
	}
}

// TestDostepnoscLiczySieZSeriiPomiarow wykazuje, że dostępność jest wyliczeniem
// z wierszy, a nie licznikiem: dwa pomiary o różnych wynikach dają wartość
// pośrednią, a nie sto ani zero procent.
func TestDostepnoscLiczySieZSeriiPomiarow(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var sprawna shared.HealthProbeSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeSave,
		shared.HealthProbeSaveRequest{
			Name: "baza", Kind: shared.HealthProbeKindInternal,
			Target: celSondyWewnetrznejBaza, IntervalMs: 60_000,
			Objective: wskaznik(99.0),
		}, &sprawna)

	// Dwa przebiegi udane i jeden nieudany, dołożony wprost do serii.
	for numer := 0; numer < 2; numer++ {
		wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeRun,
			shared.HealthProbeRunRequest{ProbeId: sprawna.Probe.Id}, nil)
	}
	if _, err := baza.Exec(
		`INSERT INTO wynik_sondy_kondycji
		     (identyfikator_zewnetrzny, sonda_kod, stan, wykonano, szczegol)
		 VALUES ('pomiar-sprawdzianu', ?, 'down', ?, 'wstawione przez sprawdzian')`,
		sprawna.Probe.Id, terazWMilisekundachKondycji()); err != nil {
		t.Fatalf("nie można dołożyć nieudanego pomiaru do serii: %v", err)
	}

	var dostepnosc shared.HealthUptimeGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthUptimeGet,
		shared.HealthUptimeGetRequest{ProbeId: wskaznik(sprawna.Probe.Id)}, &dostepnosc)

	if len(dostepnosc.Uptimes) != 1 {
		t.Fatalf("dostępność oddała %d pozycji dla jednej sondy", len(dostepnosc.Uptimes))
	}
	wynik := dostepnosc.Uptimes[0]
	if wynik.Samples != 3 {
		t.Fatalf("dostępność policzyła %d pomiarów, a w serii leżą 3", wynik.Samples)
	}
	if wynik.UptimePercent == nil {
		t.Fatal("dostępność nie oddała wartości procentowej")
	}
	// Dwa udane na trzy pomiary: dwie trzecie. Trzeci bywa degraded przy
	// obciążonej maszynie sprawdzianu.
	if *wynik.UptimePercent <= 50 || *wynik.UptimePercent >= 100 {
		t.Fatalf("dostępność wyniosła %.2f%% — przy dwóch udanych i jednym nieudanym "+
			"pomiarze wartość skrajna znaczy, że nikt serii nie policzył", *wynik.UptimePercent)
	}

	// Dostępność liczona jest z serii, a nie z kolumny odpowiedzi.
	var wierszy int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM wynik_sondy_kondycji WHERE sonda_kod = ?`,
		sprawna.Probe.Id).Scan(&wierszy); err != nil {
		t.Fatalf("nie można policzyć wierszy serii: %v", err)
	}
	if wierszy != wynik.Samples {
		t.Fatalf("odpowiedź mówi o %d pomiarach, a w bazie leży %d wierszy",
			wynik.Samples, wierszy)
	}
}

// TestSondaBezPomiarowNieUdajeSprawnosci wykazuje najważniejszą rzecz w tej
// rodzinie: sonda, która nigdy nie ruszyła, NIE pokazuje stu procent
// dostępności. Brak pomiarów nie jest dowodem sprawności.
func TestSondaBezPomiarowNieUdajeSprawnosci(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var zapis shared.HealthProbeSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeSave,
		shared.HealthProbeSaveRequest{
			Name: "sonda, która nigdy nie ruszyła", Kind: shared.HealthProbeKindInternal,
			Target: celSondyWewnetrznejBaza, IntervalMs: 60_000,
		}, &zapis)

	var dostepnosc shared.HealthUptimeGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthUptimeGet,
		shared.HealthUptimeGetRequest{ProbeId: wskaznik(zapis.Probe.Id)}, &dostepnosc)

	if len(dostepnosc.Uptimes) != 1 {
		t.Fatalf("dostępność oddała %d pozycji dla jednej sondy", len(dostepnosc.Uptimes))
	}
	wynik := dostepnosc.Uptimes[0]
	if wynik.Samples != 0 {
		t.Fatalf("sonda bez przebiegów melduje %d pomiarów", wynik.Samples)
	}
	if wynik.UptimePercent != nil {
		t.Fatalf("sonda bez ani jednego pomiaru pokazuje %.2f%% dostępności — "+
			"to jest właśnie fasada, przez którą awaria przechodzi niezauważona",
			*wynik.UptimePercent)
	}
}

// TestUsuniecieSondyZabieraSerie wykazuje, że po usunięciu definicji nie
// zostaje seria liczb bez pytania, na które odpowiadają.
func TestUsuniecieSondyZabieraSerie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var zapis shared.HealthProbeSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeSave,
		shared.HealthProbeSaveRequest{
			Name: "do usunięcia", Kind: shared.HealthProbeKindInternal,
			Target: celSondyWewnetrznejBaza, IntervalMs: 60_000,
		}, &zapis)
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeRun,
		shared.HealthProbeRunRequest{ProbeId: zapis.Probe.Id}, nil)

	var usuniecie shared.HealthProbeRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandHealthProbeRemove,
		shared.HealthProbeRemoveRequest{ProbeId: zapis.Probe.Id}, &usuniecie)
	if usuniecie.RemovedResults != 1 {
		t.Fatalf("usunięcie meldowało %d skasowanych pomiarów, a był jeden",
			usuniecie.RemovedResults)
	}

	var zostalo int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM wynik_sondy_kondycji WHERE sonda_kod = ?`,
		zapis.Probe.Id).Scan(&zostalo); err != nil {
		t.Fatalf("nie można policzyć wierszy serii: %v", err)
	}
	if zostalo != 0 {
		t.Fatalf("po usunięciu sondy w serii zostało %d wierszy", zostalo)
	}
}

// TestRegulaAlertuWyzwalaSieNaZmierzonejWartosci wykazuje, że wyzwolenie
// powstaje z POMIARU, a nie z zapisu reguły: w rejestrze leży wiersz z wartością
// obserwowaną, którą da się niezależnie sprawdzić.
func TestRegulaAlertuWyzwalaSieNaZmierzonejWartosci(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	// Materiał pomiaru: dwa nieudane pomiary sondy w tabeli serii. Miara
	// `probeFailure` liczy właśnie je.
	teraz := terazWMilisekundachKondycji()
	for numer, kod := range []string{"pomiar-a", "pomiar-b"} {
		if _, err := baza.Exec(
			`INSERT INTO sonda_kondycji
			     (identyfikator_zewnetrzny, nazwa, rodzaj, cel, odstep_ms, utworzono)
			 VALUES (?, ?, 'internal', 'database', 60000, ?)
			 ON CONFLICT(identyfikator_zewnetrzny) DO NOTHING`,
			"sonda-alertu", "sonda alertu", teraz); err != nil {
			t.Fatalf("nie można założyć sondy materiału: %v", err)
		}
		if _, err := baza.Exec(
			`INSERT INTO wynik_sondy_kondycji
			     (identyfikator_zewnetrzny, sonda_kod, stan, wykonano)
			 VALUES (?, 'sonda-alertu', 'down', ?)`, kod, teraz-int64(numer)); err != nil {
			t.Fatalf("nie można dołożyć nieudanego pomiaru: %v", err)
		}
	}

	var regula shared.AlertRuleSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAlertRuleSave,
		shared.AlertRuleSaveRequest{
			Name:       "sondy padają",
			Kind:       shared.AlertRuleKindThreshold,
			Metric:     shared.AlertMetricProbeFailure,
			Comparison: wskaznik(shared.AlertComparison(shared.AlertComparisonGreaterOrEqual)),
			Threshold:  wskaznik(2.0),
			WindowMs:   3_600_000,
			Severity:   shared.DiagnosticPriorityHigh,
			Channels:   []shared.AlertChannel{shared.AlertChannelApp},
		}, &regula)

	var wykaz shared.AlertTriggerListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAlertTriggerList,
		shared.AlertTriggerListRequest{RuleId: wskaznik(regula.Rule.Id)}, &wykaz)

	if len(wykaz.Triggers) != 1 {
		t.Fatalf("reguła na dwóch nieudanych pomiarach dała %d wyzwoleń, a miała dać jedno",
			len(wykaz.Triggers))
	}
	if wykaz.Triggers[0].ObservedValue != 2 {
		t.Fatalf("wyzwolenie niesie wartość obserwowaną %.2f, a nieudanych pomiarów są dwa",
			wykaz.Triggers[0].ObservedValue)
	}

	// Świadek niezależny: wiersz rejestru wyzwoleń wraz z wartością.
	var wartosc float64
	var stan string
	if err := baza.QueryRow(
		`SELECT wartosc_obserwowana, stan FROM wyzwolenie_alertu WHERE regula_kod = ?`,
		regula.Rule.Id).Scan(&wartosc, &stan); err != nil {
		t.Fatalf("wyzwolenie nie zostawiło wiersza w rejestrze: %v", err)
	}
	if wartosc != 2 || stan != shared.AlertTriggerStatusFiring {
		t.Fatalf("wiersz rejestru niesie wartość %.2f i stan %q", wartosc, stan)
	}
}

// TestRegulaPonizejProguNieWyzwalaSie wykazuje sprawdzian przeciwny: reguła,
// której próg nie został przekroczony, NIE zostawia wiersza. Bez tego
// sprawdzianu adapter zapisujący wyzwolenie zawsze przeszedłby poprzedni.
func TestRegulaPonizejProguNieWyzwalaSie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var regula shared.AlertRuleSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAlertRuleSave,
		shared.AlertRuleSaveRequest{
			Name:       "sto nieudanych pomiarów",
			Kind:       shared.AlertRuleKindThreshold,
			Metric:     shared.AlertMetricProbeFailure,
			Comparison: wskaznik(shared.AlertComparison(shared.AlertComparisonGreaterThan)),
			Threshold:  wskaznik(100.0),
			WindowMs:   3_600_000,
			Severity:   shared.DiagnosticPriorityLow,
			Channels:   []shared.AlertChannel{shared.AlertChannelApp},
		}, &regula)

	wykonajUdana(t, zmontowany, zycie, shared.CommandAlertTriggerList,
		shared.AlertTriggerListRequest{}, nil)

	var wierszy int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM wyzwolenie_alertu WHERE regula_kod = ?`,
		regula.Rule.Id).Scan(&wierszy); err != nil {
		t.Fatalf("nie można policzyć wyzwoleń: %v", err)
	}
	if wierszy != 0 {
		t.Fatalf("reguła o progu stu wyzwoliła się %d razy przy zerze nieudanych pomiarów", wierszy)
	}
}

// TestPotwierdzenieWyzwoleniaZapisujeChwileINotatke wykazuje, że potwierdzenie
// zostaje w wierszu, a nie tylko w odpowiedzi.
func TestPotwierdzenieWyzwoleniaZapisujeChwileINotatke(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var regula shared.AlertRuleSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAlertRuleSave,
		shared.AlertRuleSaveRequest{
			Name: "do potwierdzenia", Kind: shared.AlertRuleKindThreshold,
			Metric: shared.AlertMetricProbeFailure, WindowMs: 60_000,
			Severity: shared.DiagnosticPriorityMedium,
			Channels: []shared.AlertChannel{shared.AlertChannelApp},
		}, &regula)
	if _, err := baza.Exec(
		`INSERT INTO wyzwolenie_alertu
		     (identyfikator_zewnetrzny, regula_kod, stan, waga, miara,
		      wartosc_obserwowana, komunikat, wyzwolono)
		 VALUES ('wyzwolenie-sprawdzianu', ?, 'firing', 'medium', 'probeFailure',
		         7, 'sprawdzian', ?)`, regula.Rule.Id, terazWMilisekundachAlertow()); err != nil {
		t.Fatalf("nie można założyć wyzwolenia do potwierdzenia: %v", err)
	}

	var potwierdzenie shared.AlertTriggerAcknowledgeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAlertTriggerAcknowledge,
		shared.AlertTriggerAcknowledgeRequest{
			TriggerId: "wyzwolenie-sprawdzianu",
			Note:      wskaznik("obsłużone ręcznie"),
		}, &potwierdzenie)

	var stan, notatka string
	var potwierdzono sql.NullInt64
	if err := baza.QueryRow(
		`SELECT stan, COALESCE(notatka, ''), potwierdzono
		   FROM wyzwolenie_alertu WHERE identyfikator_zewnetrzny = 'wyzwolenie-sprawdzianu'`).
		Scan(&stan, &notatka, &potwierdzono); err != nil {
		t.Fatalf("nie można odczytać wyzwolenia po potwierdzeniu: %v", err)
	}
	if stan != shared.AlertTriggerStatusAcknowledged {
		t.Fatalf("wiersz niesie stan %q po potwierdzeniu", stan)
	}
	if notatka != "obsłużone ręcznie" {
		t.Fatalf("notatka Operatora nie doszła do wiersza: %q", notatka)
	}
	if !potwierdzono.Valid || potwierdzono.Int64 <= 0 {
		t.Fatal("wiersz nie ma chwili potwierdzenia — czas reakcji nie ma z czego się policzyć")
	}
}

// TestRegulaNaMiarzeBezZrodlaNiePowstaje wykazuje, że rdzeń nie pozwala założyć
// reguły, która nigdy by nie zawołała.
func TestRegulaNaMiarzeBezZrodlaNiePowstaje(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandAlertRuleSave,
		shared.AlertRuleSaveRequest{
			Name: "miara, której nie ma", Kind: shared.AlertRuleKindThreshold,
			Metric: shared.AlertMetric("wymyslonaMiara"), WindowMs: 60_000,
			Severity: shared.DiagnosticPriorityLow,
			Channels: []shared.AlertChannel{shared.AlertChannelApp},
		})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("odmowa nieznanej miary ma kod %q", odmowa.Code)
	}
}

// wykorzystanieKontekstu trzyma import kontekstu przy sprawdzianach, które
// posługują się nim wyłącznie przez uprząż.
var _ = context.Background
