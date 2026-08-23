package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/store"
	"danacoconsole/shared"
)

// Skutek rodziny `schedule.*`: czy okno wykonania, nadzór obecności uruchomień
// i klucz podpisu webhooka PRZEŻYWAJĄ ponowne złożenie rdzenia.
//
// To jest właściwa miara tej rodziny. Wyzwalacz czasowy trzymany wyłącznie
// w pamięci procesu jest wyzwalaczem, który po restarcie nigdy nie zadziała,
// a okno Scheduler pokazywałoby go dalej jako obowiązujący. Sprawdzian, który
// zapisuje i odczytuje w jednym procesie, tego nie zauważy — więc tutaj rdzeń
// składa się DRUGI RAZ nad tą samą bazą i dopiero wtedy pada pytanie.

// zlozRdzenPonownieNadBaza składa rdzeń raz jeszcze nad katalogiem danych,
// w którym leży baza po pierwszym złożeniu. Odpowiada temu, co dzieje się przy
// ponownym uruchomieniu rdzenia na maszynie Operatora.
func zlozRdzenPonownieNadBaza(t *testing.T, katalog string) (*Zmontowany, context.Context) {
	t.Helper()

	baza, err := store.Otworz(filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy przy ponownym złożeniu: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })

	zycie, zakoncz := context.WithCancel(context.Background())
	t.Cleanup(zakoncz)

	ustawienia := konfiguracja.Domyslna()
	ustawienia.KatalogDanych = katalog
	ustawienia.KatalogKlienta = ""
	ustawienia.KatalogProfili = ""
	ustawienia.Port = 0

	zmontowany, err := Zmontuj(zycie, Montaz{
		Konfiguracja: ustawienia, Baza: baza, Dziennik: dziennikNiemy(),
	})
	if err != nil {
		t.Fatalf("ponowne złożenie rdzenia nie powiodło się: %v", err)
	}
	t.Cleanup(zmontowany.Zamknij)
	return zmontowany, zycie
}

// harmonogramSprawdzianu zakłada automatykę z harmonogramem i oddaje kod
// automatyki, kod harmonogramu oraz klucz wiersza harmonogramu.
func harmonogramSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	baza *sql.DB) (string, string, int64) {
	t.Helper()

	var zapis shared.AutomationWorkflowSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowSave,
		shared.AutomationWorkflowSaveRequest{
			Name: "Raport poniedziałkowy",
			Steps: []shared.AutomationStep{
				{Id: "krok-raport", Kind: shared.AutomationStepKindCommand,
					Command: wskaznik("research.report.build")},
			},
		}, &zapis)

	var harmonogram shared.AutomationScheduleSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationScheduleSet,
		shared.AutomationScheduleSetRequest{
			WorkflowId: zapis.Workflow.Id, Cron: wskaznik("0 7 * * 1"),
			TimeZone: wskaznik("CET"), Enabled: wskaznik(true),
		}, &harmonogram)

	var id int64
	err := baza.QueryRow(`SELECT id FROM harmonogram_automatyki
	                      WHERE identyfikator_zewnetrzny = ?`, harmonogram.Schedule.Id).Scan(&id)
	if err != nil {
		t.Fatalf("harmonogramu nie ma w bazie po zapisie: %v", err)
	}
	return zapis.Workflow.Id, harmonogram.Schedule.Id, id
}

// TestSkutekOkienWykonaniaPrzezywaZlozenieRdzenia mierzy okna wykonania i nadzór
// obecności uruchomień PO ponownym złożeniu rdzenia.
func TestSkutekOkienWykonaniaPrzezywaZlozenieRdzenia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	_, kodHarmonogramu, idHarmonogramu := harmonogramSprawdzianu(t, zmontowany, zycie, baza)

	wykonajUdana(t, zmontowany, zycie, shared.CommandScheduleWindowSet,
		shared.ScheduleWindowSetRequest{
			ScheduleId: kodHarmonogramu,
			Windows: []shared.AutomationExecutionWindow{
				{DaysOfWeek: []int{1, 2, 3, 4, 5}, FromMinuteOfDay: 480, ToMinuteOfDay: 1020,
					TimeZone: wskaznik("CET")},
			},
		}, nil)
	wykonajUdana(t, zmontowany, zycie, shared.CommandScheduleHeartbeatSet,
		shared.ScheduleHeartbeatSetRequest{ScheduleId: kodHarmonogramu, ToleranceSeconds: 900}, nil)

	// Zamykamy rdzeń pierwszy i składamy drugi nad tą samą bazą. Dopiero to
	// odróżnia stan zapisany od stanu trzymanego w pamięci procesu.
	zmontowany.Zamknij()
	drugi, zycieDrugiego := zlozRdzenPonownieNadBaza(t, katalog)

	var dni string
	var od, do, tolerancja int
	err := baza.QueryRow(`SELECT dni_tygodnia, minuta_od, minuta_do
	                      FROM okno_wykonania_harmonogramu WHERE harmonogram_id = ?`,
		idHarmonogramu).Scan(&dni, &od, &do)
	if err != nil {
		t.Fatalf("okna wykonania nie ma w bazie po ponownym złożeniu rdzenia: %v", err)
	}
	if od != 480 || do != 1020 {
		t.Fatalf("okno wykonania w bazie to %d–%d, oczekiwano 480–1020", od, do)
	}
	dniTygodnia := []int{}
	if err := json.Unmarshal([]byte(dni), &dniTygodnia); err != nil {
		t.Fatalf("dni tygodnia okna są w bazie nieczytelne: %v", err)
	}
	if len(dniTygodnia) != 5 {
		t.Fatalf("w bazie stoi %d dni tygodnia okna, oczekiwano 5", len(dniTygodnia))
	}
	if err := baza.QueryRow(`SELECT tolerancja_sekundy FROM harmonogram_automatyki WHERE id = ?`,
		idHarmonogramu).Scan(&tolerancja); err != nil {
		t.Fatalf("nie można odczytać tolerancji nadzoru: %v", err)
	}
	if tolerancja != 900 {
		t.Fatalf("w bazie stoi tolerancja %d, oczekiwano 900", tolerancja)
	}

	// Rdzeń złożony na nowo ma oddać ten sam harmonogram — nie „harmonogram
	// nie istnieje”, bo to znaczyłoby, że stan zginął razem z procesem.
	var odczyt shared.ScheduleGetResponse
	wykonajUdana(t, drugi, zycieDrugiego, shared.CommandScheduleGet,
		shared.ScheduleGetRequest{ScheduleId: wskaznik(kodHarmonogramu)}, &odczyt)
	if len(odczyt.Schedules) != 1 || odczyt.Schedules[0].Id != kodHarmonogramu {
		t.Fatalf("rdzeń złożony ponownie oddał %d harmonogramów, oczekiwano zapisanego",
			len(odczyt.Schedules))
	}

	// Ta sama komenda na nowym rdzeniu ma przyjąć harmonogram, a nie odmówić
	// nieistnieniem: nadzór po restarcie musi dać się zmienić.
	wykonajUdana(t, drugi, zycieDrugiego, shared.CommandScheduleHeartbeatSet,
		shared.ScheduleHeartbeatSetRequest{ScheduleId: kodHarmonogramu, ToleranceSeconds: 0}, nil)
	if err := baza.QueryRow(`SELECT tolerancja_sekundy FROM harmonogram_automatyki WHERE id = ?`,
		idHarmonogramu).Scan(&tolerancja); err != nil {
		t.Fatalf("nie można odczytać tolerancji po zmianie na nowym rdzeniu: %v", err)
	}
	if tolerancja != 0 {
		t.Fatalf("zerowa tolerancja nie zeszła do bazy (stoi %d) — nadzór po restarcie "+
			"nie dawałby się wyłączyć", tolerancja)
	}
}

// TestSkutekKluczaPodpisuWebhookaWSejfie mierzy dwie rzeczy: że referencja
// klucza zostaje w bazie i przeżywa złożenie rdzenia, oraz że sama WARTOŚĆ
// klucza nie znalazła się w bazie ani w odpowiedzi.
func TestSkutekKluczaPodpisuWebhookaWSejfie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kodAutomatyki, _, idHarmonogramu := harmonogramSprawdzianu(t, zmontowany, zycie, baza)

	var pierwszy shared.ScheduleWebhookEndpointGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandScheduleWebhookEndpointGet,
		shared.ScheduleWebhookEndpointGetRequest{WorkflowId: kodAutomatyki}, &pierwszy)
	if pierwszy.EndpointUrl == "" || pierwszy.SignatureSecretRef == "" {
		t.Fatalf("odczyt webhooka oddał pusty adres albo pustą referencję: %+v", pierwszy)
	}
	if pierwszy.DeduplicationWindowSeconds <= 0 {
		t.Fatalf("okno deduplikacji wyszło jako %d — wywołanie powtórzone ruszyłoby "+
			"automatykę drugi raz", pierwszy.DeduplicationWindowSeconds)
	}

	var odwolanie sql.NullString
	if err := baza.QueryRow(`SELECT odwolanie_podpisu FROM harmonogram_automatyki WHERE id = ?`,
		idHarmonogramu).Scan(&odwolanie); err != nil {
		t.Fatalf("nie można odczytać odwołania podpisu z bazy: %v", err)
	}
	if !odwolanie.Valid || odwolanie.String != pierwszy.SignatureSecretRef {
		t.Fatalf("w bazie stoi odwołanie %+v, a komenda oddała %q",
			odwolanie, pierwszy.SignatureSecretRef)
	}

	// Odczyt powtórzony BEZ wymiany ma oddać tę samą referencję — inaczej każde
	// otwarcie okna unieważniałoby podpisy nadawcy.
	var powtorzony shared.ScheduleWebhookEndpointGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandScheduleWebhookEndpointGet,
		shared.ScheduleWebhookEndpointGetRequest{WorkflowId: kodAutomatyki}, &powtorzony)
	if powtorzony.SignatureSecretRef != pierwszy.SignatureSecretRef {
		t.Fatalf("odczyt bez wymiany zmienił referencję klucza podpisu")
	}

	// Wymiana ma oddać referencję INNĄ i zapisać ją w bazie.
	var wymieniony shared.ScheduleWebhookEndpointGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandScheduleWebhookEndpointGet,
		shared.ScheduleWebhookEndpointGetRequest{
			WorkflowId: kodAutomatyki, RotateSecret: wskaznik(true),
		}, &wymieniony)
	if wymieniony.SignatureSecretRef == pierwszy.SignatureSecretRef {
		t.Fatalf("wymiana klucza podpisu oddała tę samą referencję")
	}

	zmontowany.Zamknij()
	drugi, zycieDrugiego := zlozRdzenPonownieNadBaza(t, katalog)

	var poZlozeniu shared.ScheduleWebhookEndpointGetResponse
	wykonajUdana(t, drugi, zycieDrugiego, shared.CommandScheduleWebhookEndpointGet,
		shared.ScheduleWebhookEndpointGetRequest{WorkflowId: kodAutomatyki}, &poZlozeniu)
	if poZlozeniu.SignatureSecretRef != wymieniony.SignatureSecretRef {
		t.Fatalf("po ponownym złożeniu rdzenia referencja klucza podpisu to %q, "+
			"a przed nim była %q — klucz nie przeżył restartu",
			poZlozeniu.SignatureSecretRef, wymieniony.SignatureSecretRef)
	}
}

// TestSkutekUruchomieniaWstecznegoWBazie mierzy `schedule.backfill.run`:
// czy przebiegi naprawdę powstały i czy historia wyzwoleń je odnotowała.
func TestSkutekUruchomieniaWstecznegoWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kodAutomatyki, kodHarmonogramu, idHarmonogramu :=
		harmonogramSprawdzianu(t, zmontowany, zycie, baza)

	// Zakres obejmujący trzy poniedziałki sierpnia 2026 (3, 10 i 17 sierpnia).
	const odMilisekund = int64(1753920000000) // 2025-07-31, na pewno przed zakresem
	const doMilisekund = int64(1755561600000) // 2025-08-19

	var wynik shared.ScheduleBackfillRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandScheduleBackfillRun,
		shared.ScheduleBackfillRunRequest{
			ScheduleId: kodHarmonogramu, FromAt: odMilisekund, ToAt: doMilisekund,
			MaxRuns: wskaznik(3),
		}, &wynik)
	if len(wynik.QueuedAt) == 0 {
		t.Fatalf("uruchomienie wsteczne nie zakolejkowało ani jednego terminu")
	}

	// Przebiegi mają być w bazie, a nie tylko w odpowiedzi.
	var automatykaID int64
	if err := baza.QueryRow(`SELECT id FROM automatyka WHERE identyfikator_zewnetrzny = ?`,
		kodAutomatyki).Scan(&automatykaID); err != nil {
		t.Fatalf("automatyki nie ma w bazie: %v", err)
	}
	var przebiegow int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM przebieg_automatyki WHERE automatyka_id = ?`,
		automatykaID).Scan(&przebiegow); err != nil {
		t.Fatalf("nie można policzyć przebiegów: %v", err)
	}
	if przebiegow != len(wynik.QueuedAt) {
		t.Fatalf("uruchomienie wsteczne zameldowało %d terminów, a w bazie stoi %d przebiegów",
			len(wynik.QueuedAt), przebiegow)
	}

	var wyzwolen int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM wyzwolenie_automatyki
	                         WHERE harmonogram_id = ? AND przyczyna = 'backfill'`,
		idHarmonogramu).Scan(&wyzwolen); err != nil {
		t.Fatalf("nie można policzyć wyzwoleń wstecznych: %v", err)
	}
	if wyzwolen != len(wynik.QueuedAt) {
		t.Fatalf("historia odnotowała %d wyzwoleń wstecznych, oczekiwano %d",
			wyzwolen, len(wynik.QueuedAt))
	}

	var historia shared.ScheduleTriggerHistoryResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandScheduleTriggerHistory,
		shared.ScheduleTriggerHistoryRequest{ScheduleId: wskaznik(kodHarmonogramu)}, &historia)
	if len(historia.Entries) != wyzwolen {
		t.Fatalf("odczyt historii oddał %d wpisów, a w bazie stoi %d",
			len(historia.Entries), wyzwolen)
	}
	for _, wpis := range historia.Entries {
		if wpis.Cause != shared.AutomationTriggerCauseBackfill {
			t.Fatalf("wpis historii niesie przyczynę %q, oczekiwano „backfill”", wpis.Cause)
		}
	}
}
