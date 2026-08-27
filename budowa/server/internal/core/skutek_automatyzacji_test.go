// Plik mierzy skutek dobudowy modułu Automations: czy za odpowiedzią komendy
// w bazie rzeczywiście leży wiersz, niezależnie od treści samej odpowiedzi.
package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"danacoconsole/shared"
)

// automatykaSprawdzianuDobudowy zakłada automatykę o trzech krokach i oddaje
// jej kod wraz z kluczem wiersza — sprawdziany schodzą do bazy po kluczu.
func automatykaSprawdzianuDobudowy(t *testing.T, zmontowany *Zmontowany,
	zycie context.Context, baza *sql.DB, nazwa string) (string, int64) {
	t.Helper()

	var zapis shared.AutomationWorkflowSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowSave,
		shared.AutomationWorkflowSaveRequest{
			Name: nazwa,
			Steps: []shared.AutomationStep{
				{Id: "krok-pobierz", Kind: shared.AutomationStepKindCommand,
					Command: wskaznik("library.file.list")},
				{Id: "krok-waliduj", Kind: shared.AutomationStepKindCondition,
					Condition: wskaznik("wynik.liczba > 0")},
				{Id: "krok-wyslij", Kind: shared.AutomationStepKindCommand,
					Command: wskaznik("mail.message.send")},
			},
		}, &zapis)

	var id int64
	err := baza.QueryRow(`SELECT id FROM automatyka WHERE identyfikator_zewnetrzny = ?`,
		zapis.Workflow.Id).Scan(&id)
	if err != nil {
		t.Fatalf("automatyki %q nie ma w bazie po zapisie: %v", zapis.Workflow.Id, err)
	}
	return zapis.Workflow.Id, id
}

// TestSkutekWersjiDefinicjiWBazie mierzy, czy zapis definicji zostawia migawkę,
// czy panel „Wersje” pokazuje historię, której nie ma.
func TestSkutekWersjiDefinicjiWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := automatykaSprawdzianuDobudowy(t, zmontowany, zycie, baza, "Raport tygodniowy")

	var migawek int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM wersja_automatyki WHERE automatyka_id = ?`,
		id).Scan(&migawek); err != nil {
		t.Fatalf("nie można policzyć wersji w bazie: %v", err)
	}
	if migawek != 1 {
		t.Fatalf("zapis definicji zostawił %d migawek wersji, oczekiwano 1 — panel „Wersje” "+
			"pokazywałby historię, której w bazie nie ma", migawek)
	}

	// Drugi zapis podnosi wersję, więc migawek ma być dwie, nie jedna nadpisana.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowSave,
		shared.AutomationWorkflowSaveRequest{
			WorkflowId: wskaznik(kod), Name: "Raport tygodniowy",
			Steps: []shared.AutomationStep{
				{Id: "krok-pobierz", Kind: shared.AutomationStepKindCommand,
					Command: wskaznik("library.file.list")},
			},
		}, nil)

	if err := baza.QueryRow(`SELECT COUNT(*) FROM wersja_automatyki WHERE automatyka_id = ?`,
		id).Scan(&migawek); err != nil {
		t.Fatalf("nie można policzyć wersji po drugim zapisie: %v", err)
	}
	if migawek != 2 {
		t.Fatalf("po dwóch zapisach w bazie stoi %d migawek, oczekiwano 2", migawek)
	}

	// Porównanie ma zobaczyć dwa kroki usunięte — mierzymy skutek zapisu, a nie
	// samo istnienie wierszy.
	var roznica shared.AutomationWorkflowVersionDiffResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowVersionDiff,
		shared.AutomationWorkflowVersionDiffRequest{
			WorkflowId: kod, FromVersion: 1, ToVersion: 2,
		}, &roznica)
	usuniete := 0
	for _, zmiana := range roznica.Changes {
		if zmiana.Change == shared.ChangeKindDeleted {
			usuniete++
		}
	}
	if usuniete != 2 {
		t.Fatalf("porównanie wersji nazwało %d kroków usuniętych, oczekiwano 2: %+v",
			usuniete, roznica.Changes)
	}
}

// TestSkutekPrzywroceniaWersjiWBazie mierzy, czy przywrócenie wersji naprawdę
// odtwarza kroki w tabeli `krok_automatyki`, czy oddaje je wyłącznie w odpowiedzi.
func TestSkutekPrzywroceniaWersjiWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := automatykaSprawdzianuDobudowy(t, zmontowany, zycie, baza, "Odtworzenie")

	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowSave,
		shared.AutomationWorkflowSaveRequest{
			WorkflowId: wskaznik(kod), Name: "Odtworzenie",
			Steps: []shared.AutomationStep{
				{Id: "krok-jedyny", Kind: shared.AutomationStepKindCommand},
			},
		}, nil)

	if krokow := liczbaKrokowWBazie(t, baza, id); krokow != 1 {
		t.Fatalf("po zapisie skracającym definicję w bazie stoi %d kroków, oczekiwano 1", krokow)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowVersionRestore,
		shared.AutomationWorkflowVersionRestoreRequest{WorkflowId: kod, Version: 1}, nil)

	if krokow := liczbaKrokowWBazie(t, baza, id); krokow != 3 {
		t.Fatalf("przywrócenie wersji zostawiło w bazie %d kroków, oczekiwano 3 — "+
			"odpowiedź komendy oddawałaby definicję, której w bazie nie ma", krokow)
	}
}

// liczbaKrokowWBazie liczy kroki automatyki własnym zapytaniem SQL, niezależnie od odpowiedzi komendy.
func liczbaKrokowWBazie(t *testing.T, baza *sql.DB, automatykaID int64) int {
	t.Helper()

	var krokow int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM krok_automatyki WHERE automatyka_id = ?`,
		automatykaID).Scan(&krokow); err != nil {
		t.Fatalf("nie można policzyć kroków w bazie: %v", err)
	}
	return krokow
}

// TestSkutekEtykietPublikacjiIUdostepnieniaWBazie mierzy trzy czynności naraz,
// bo wszystkie trzy oddają ten sam byt kontraktu i wszystkie trzy dałyby się
// udać bez zapisu.
func TestSkutekEtykietPublikacjiIUdostepnieniaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := automatykaSprawdzianuDobudowy(t, zmontowany, zycie, baza, "Etykiety")

	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowTagSet,
		shared.AutomationWorkflowTagSetRequest{
			WorkflowId: kod, Tags: []string{"raport", "tygodniowy"},
		}, nil)

	wiersze, err := baza.Query(
		`SELECT etykieta FROM etykieta_automatyki WHERE automatyka_id = ? ORDER BY etykieta`, id)
	if err != nil {
		t.Fatalf("nie można odczytać etykiet z bazy: %v", err)
	}
	defer wiersze.Close()
	etykiety := []string{}
	for wiersze.Next() {
		var etykieta string
		if err := wiersze.Scan(&etykieta); err != nil {
			t.Fatalf("nieczytelna etykieta w bazie: %v", err)
		}
		etykiety = append(etykiety, etykieta)
	}
	if len(etykiety) != 2 || etykiety[0] != "raport" || etykiety[1] != "tygodniowy" {
		t.Fatalf("w bazie stoją etykiety %v, oczekiwano [raport tygodniowy]", etykiety)
	}

	// Wykaz pusty ma zdjąć wszystkie etykiety, nie zostawić poprzednich.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowTagSet,
		shared.AutomationWorkflowTagSetRequest{WorkflowId: kod, Tags: []string{}}, nil)
	var pozostalo int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM etykieta_automatyki WHERE automatyka_id = ?`,
		id).Scan(&pozostalo); err != nil {
		t.Fatalf("nie można policzyć etykiet po zdjęciu: %v", err)
	}
	if pozostalo != 0 {
		t.Fatalf("wykaz pusty zostawił %d etykiet w bazie, oczekiwano zera", pozostalo)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowPublish,
		shared.AutomationWorkflowPublishRequest{WorkflowId: kod, Version: wskaznik(1)}, nil)
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowShare,
		shared.AutomationWorkflowShareRequest{WorkflowId: kod, Shared: true}, nil)
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationExecutionBudgetSet,
		shared.AutomationExecutionBudgetSetRequest{
			WorkflowId: kod, RunBudgetSeconds: wskaznik(900), StepBudgetSeconds: wskaznik(120),
		}, nil)

	var opublikowana sql.NullInt64
	var udostepniona, budzetPrzebiegu, budzetKroku int
	err = baza.QueryRow(`SELECT wersja_opublikowana, udostepniona,
	                            budzet_przebiegu_sekundy, budzet_kroku_sekundy
	                     FROM automatyka WHERE id = ?`, id).
		Scan(&opublikowana, &udostepniona, &budzetPrzebiegu, &budzetKroku)
	if err != nil {
		t.Fatalf("nie można odczytać kolumn automatyki: %v", err)
	}
	if !opublikowana.Valid || opublikowana.Int64 != 1 {
		t.Fatalf("publikacja nie zapisała numeru wersji w bazie: %+v", opublikowana)
	}
	if udostepniona != 1 {
		t.Fatalf("udostępnienie nie zostawiło śladu w bazie")
	}
	if budzetPrzebiegu != 900 || budzetKroku != 120 {
		t.Fatalf("budżety w bazie to %d/%d, oczekiwano 900/120", budzetPrzebiegu, budzetKroku)
	}

	// Wersja opublikowana ma wyjść kontraktem jako wersja wykonywana, nie robocza.
	var wykaz shared.AutomationWorkflowListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowList,
		shared.AutomationWorkflowListRequest{}, &wykaz)
	znaleziona := false
	for _, automatyka := range wykaz.Workflows {
		if automatyka.Id != kod {
			continue
		}
		znaleziona = true
		if automatyka.Version == nil || *automatyka.Version != 1 {
			t.Fatalf("wykaz oddał wersję %+v, oczekiwano opublikowanej 1", automatyka.Version)
		}
	}
	if !znaleziona {
		t.Fatalf("automatyki %q nie ma w wykazie", kod)
	}
}

// TestSkutekZmiennychNotatkiIUkladuWBazie mierzy trzy zapisy panelu „Zmienne ▼”
// i kanwy. Wszystkie trzy oddają byt, który powstaje z odczytu — więc odpowiedź
// udana bez zapisu byłaby tu nie do odróżnienia bez zejścia do bazy.
func TestSkutekZmiennychNotatkiIUkladuWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := automatykaSprawdzianuDobudowy(t, zmontowany, zycie, baza, "Zmienne")

	var zapis shared.AutomationWorkflowVariablesSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowVariablesSet,
		shared.AutomationWorkflowVariablesSetRequest{
			WorkflowId: kod,
			Variables: []shared.AutomationVariable{
				{Name: "odbiorca", Kind: "string"},
				{Name: "nieuzywana", Kind: "string"},
			},
			Mappings: []shared.AutomationDataMapping{
				{FromStepId: "krok-pobierz", FromPath: "$.wynik", ToStepId: "krok-waliduj",
					ToField: "odbiorca"},
				{FromStepId: "krok-pobierz", FromPath: "$.wynik", ToStepId: "krok-nieistniejacy",
					ToField: "cel"},
			},
		}, &zapis)

	var zmiennych, mapowan int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM zmienna_automatyki WHERE automatyka_id = ?`,
		id).Scan(&zmiennych); err != nil {
		t.Fatalf("nie można policzyć zmiennych: %v", err)
	}
	if err := baza.QueryRow(
		`SELECT COUNT(*) FROM mapowanie_danych_automatyki WHERE automatyka_id = ?`,
		id).Scan(&mapowan); err != nil {
		t.Fatalf("nie można policzyć mapowań: %v", err)
	}
	if zmiennych != 2 || mapowan != 2 {
		t.Fatalf("w bazie stoją %d zmienne i %d mapowania, oczekiwano 2 i 2", zmiennych, mapowan)
	}
	// Zastrzeżenia nazywają mapowanie do kroku nieistniejącego i zmienną nieużywaną.
	if len(zapis.Issues) < 2 {
		t.Fatalf("zapis oddał %d zastrzeżeń, oczekiwano co najmniej dwóch: %v",
			len(zapis.Issues), zapis.Issues)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationStepNoteSet,
		shared.AutomationStepNoteSetRequest{
			WorkflowId: kod, StepId: "krok-waliduj", Note: "Warunek ustalony z Operatorem",
		}, nil)
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationStepLayoutSet,
		shared.AutomationStepLayoutSetRequest{
			WorkflowId: kod,
			Positions: []shared.AutomationStepPosition{
				{StepId: "krok-pobierz", X: 40, Y: 120},
			},
		}, nil)

	var notatka sql.NullString
	var x, y int
	err := baza.QueryRow(`SELECT notatka, wspolrzedna_x, wspolrzedna_y
	                      FROM adnotacja_kroku_automatyki
	                      WHERE automatyka_id = ? AND krok_kod = ?`, id, "krok-waliduj").
		Scan(&notatka, &x, &y)
	if err != nil {
		t.Fatalf("notatki kroku nie ma w bazie: %v", err)
	}
	if !notatka.Valid || notatka.String != "Warunek ustalony z Operatorem" {
		t.Fatalf("w bazie stoi notatka %+v, oczekiwano treści zapisanej", notatka)
	}
	err = baza.QueryRow(`SELECT wspolrzedna_x, wspolrzedna_y FROM adnotacja_kroku_automatyki
	                     WHERE automatyka_id = ? AND krok_kod = ?`, id, "krok-pobierz").Scan(&x, &y)
	if err != nil {
		t.Fatalf("położenia węzła nie ma w bazie: %v", err)
	}
	if x != 40 || y != 120 {
		t.Fatalf("w bazie stoi położenie %d/%d, oczekiwano 40/120", x, y)
	}

	// Zapis definicji podmienia kroki w całości, a adnotacja ma to przeżyć nietknięta.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowSave,
		shared.AutomationWorkflowSaveRequest{
			WorkflowId: wskaznik(kod), Name: "Zmienne",
			Steps: []shared.AutomationStep{
				{Id: "krok-pobierz", Kind: shared.AutomationStepKindCommand,
					Command: wskaznik("library.file.list")},
			},
		}, nil)
	err = baza.QueryRow(`SELECT wspolrzedna_x FROM adnotacja_kroku_automatyki
	                     WHERE automatyka_id = ? AND krok_kod = ?`, id, "krok-pobierz").Scan(&x)
	if err != nil || x != 40 {
		t.Fatalf("zapis definicji skasował układ kanwy (x=%d, err=%v)", x, err)
	}
}

// TestSkutekSzablonuPrzeplywuWBazie mierzy, czy szablon zostaje w bazie wraz ze
// swoimi parametrami i czy zastosowanie zakłada automatykę Z JEGO krokami.
func TestSkutekSzablonuPrzeplywuWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, _ := automatykaSprawdzianuDobudowy(t, zmontowany, zycie, baza, "Wzorzec")

	var zapis shared.AutomationTemplateSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationTemplateSave,
		shared.AutomationTemplateSaveRequest{
			WorkflowId: kod, Name: "Raport cykliczny",
			Parameters: []shared.AutomationTemplateParameter{
				{Name: "odbiorca", Required: wskaznik(true)},
			},
		}, &zapis)

	var szablonID int64
	var krokiSzablonu string
	err := baza.QueryRow(`SELECT id, kroki FROM szablon_automatyki
	                      WHERE identyfikator_zewnetrzny = ?`, zapis.Template.Id).
		Scan(&szablonID, &krokiSzablonu)
	if err != nil {
		t.Fatalf("szablonu nie ma w bazie: %v", err)
	}
	var parametrow int
	if err := baza.QueryRow(
		`SELECT COUNT(*) FROM parametr_szablonu_automatyki WHERE szablon_id = ?`,
		szablonID).Scan(&parametrow); err != nil {
		t.Fatalf("nie można policzyć parametrów szablonu: %v", err)
	}
	if parametrow != 1 {
		t.Fatalf("w bazie stoi %d parametrów szablonu, oczekiwano 1", parametrow)
	}
	kroki := []shared.AutomationStep{}
	if err := json.Unmarshal([]byte(krokiSzablonu), &kroki); err != nil {
		t.Fatalf("migawka kroków szablonu jest nieczytelna: %v", err)
	}
	if len(kroki) != 3 {
		t.Fatalf("szablon zapisał %d kroków, oczekiwano 3", len(kroki))
	}

	// Zastosowanie bez wartości parametru zakłada automatykę i nazywa brak, nie odmawia.
	var zastosowanie shared.AutomationTemplateApplyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationTemplateApply,
		shared.AutomationTemplateApplyRequest{
			TemplateId: zapis.Template.Id, Name: "Raport z szablonu",
		}, &zastosowanie)
	if len(zastosowanie.MissingParameters) != 1 {
		t.Fatalf("zastosowanie nazwało %d parametrów bez wartości, oczekiwano 1: %v",
			len(zastosowanie.MissingParameters), zastosowanie.MissingParameters)
	}

	var zalozonaID int64
	if err := baza.QueryRow(`SELECT id FROM automatyka WHERE identyfikator_zewnetrzny = ?`,
		zastosowanie.Workflow.Id).Scan(&zalozonaID); err != nil {
		t.Fatalf("automatyki z szablonu nie ma w bazie: %v", err)
	}
	if krokow := liczbaKrokowWBazie(t, baza, zalozonaID); krokow != 3 {
		t.Fatalf("automatyka z szablonu ma w bazie %d kroków, oczekiwano 3", krokow)
	}
}

// TestSkutekSkarbcaIAudytuWBazie mierzy, czy wiersz poświadczenia w bazie nie
// niesie jego wartości i czy dziennik audytu zapełnia się sam, przy okazji czynności.
func TestSkutekSkarbcaIAudytuWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := automatykaSprawdzianuDobudowy(t, zmontowany, zycie, baza, "Skarbiec")

	const wartosc = "wartosc-ktorej-nigdy-nie-wolno-zapisac-do-bazy"
	var zapis shared.AutomationSecretSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationSecretSet,
		shared.AutomationSecretSetRequest{Name: "klucz-poczty", Value: wartosc}, &zapis)

	if zapis.Secret.Ref == "" {
		t.Fatalf("zapis poświadczenia nie oddał referencji")
	}
	var nazwa string
	if err := baza.QueryRow(`SELECT nazwa FROM poswiadczenie_automatyki WHERE odwolanie = ?`,
		zapis.Secret.Ref).Scan(&nazwa); err != nil {
		t.Fatalf("referencji poświadczenia nie ma w bazie: %v", err)
	}
	if nazwa != "klucz-poczty" {
		t.Fatalf("w bazie stoi nazwa %q, oczekiwano „klucz-poczty”", nazwa)
	}

	var trafien int
	err := baza.QueryRow(`SELECT COUNT(*) FROM poswiadczenie_automatyki
	                      WHERE odwolanie LIKE '%' || ? || '%'
	                         OR nazwa LIKE '%' || ? || '%'
	                         OR COALESCE(zasieg,'') LIKE '%' || ? || '%'
	                         OR COALESCE(zasieg_id,'') LIKE '%' || ? || '%'`,
		wartosc, wartosc, wartosc, wartosc).Scan(&trafien)
	if err != nil {
		t.Fatalf("nie można przeszukać skarbca: %v", err)
	}
	if trafien != 0 {
		t.Fatalf("wartość poświadczenia znalazła się w bazie w %d wierszach — "+
			"skarbiec ma trzymać wyłącznie referencję", trafien)
	}

	// Usunięcie ma nazwać kroki, które straciły pokrycie, i mimo to usunąć.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowSave,
		shared.AutomationWorkflowSaveRequest{
			WorkflowId: wskaznik(kod), Name: "Skarbiec",
			Steps: []shared.AutomationStep{
				{Id: "krok-pobierz", Kind: shared.AutomationStepKindCommand,
					Command: wskaznik("library.file.list"), SecretRefs: []string{zapis.Secret.Ref}},
			},
		}, nil)

	var usuniecie shared.AutomationSecretRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationSecretRemove,
		shared.AutomationSecretRemoveRequest{SecretRef: zapis.Secret.Ref}, &usuniecie)
	if !usuniecie.Removed {
		t.Fatalf("usunięcie poświadczenia zameldowało brak skutku")
	}
	if len(usuniecie.ReferencingStepIds) != 1 || usuniecie.ReferencingStepIds[0] != "krok-pobierz" {
		t.Fatalf("usunięcie nazwało kroki %v, oczekiwano [krok-pobierz]",
			usuniecie.ReferencingStepIds)
	}
	var pozostalo int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM poswiadczenie_automatyki WHERE odwolanie = ?`,
		zapis.Secret.Ref).Scan(&pozostalo); err != nil {
		t.Fatalf("nie można sprawdzić skarbca po usunięciu: %v", err)
	}
	if pozostalo != 0 {
		t.Fatalf("referencja została w bazie mimo usunięcia — usunięcie nie może być " +
			"wstrzymywane przywołaniami")
	}

	// Dziennik audytu ma nieść ślad trzech czynności; sprawdzian liczy w bazie.
	var wpisow int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM wpis_audytu_automatyki`).Scan(&wpisow); err != nil {
		t.Fatalf("nie można policzyć wpisów audytu: %v", err)
	}
	if wpisow < 3 {
		t.Fatalf("dziennik audytu ma %d wpisów, oczekiwano co najmniej trzech — "+
			"audyt zapełnia się przy okazji czynności, nie na żądanie", wpisow)
	}

	var dziennik shared.AutomationAuditListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationAuditList,
		shared.AutomationAuditListRequest{WorkflowId: wskaznik(kod)}, &dziennik)
	if len(dziennik.Entries) == 0 {
		t.Fatalf("odczyt audytu automatyki %d oddał pusty dziennik przy wpisach w bazie", id)
	}
}

// TestSkutekRegulyAlarmowaniaWBazie mierzy, czy reguła zostaje w bazie wraz
// z kanałami — kanały leżą zapisem strukturalnym, więc zapis pusty byłby regułą
// bez adresata, o której nikt by się nie dowiedział.
func TestSkutekRegulyAlarmowaniaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	kod, id := automatykaSprawdzianuDobudowy(t, zmontowany, zycie, baza, "Alarmy")

	var zapis shared.AutomationAlertRuleSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationAlertRuleSet,
		shared.AutomationAlertRuleSetRequest{
			WorkflowId: kod, Trigger: shared.AutomationAlertTriggerFailure,
			Channels: []string{"mobile", "poczta"},
		}, &zapis)

	var wyzwalacz, kanaly string
	var automatykaID int64
	err := baza.QueryRow(`SELECT automatyka_id, wyzwalacz, kanaly
	                      FROM regula_alarmowania_automatyki
	                      WHERE identyfikator_zewnetrzny = ?`, zapis.Rule.Id).
		Scan(&automatykaID, &wyzwalacz, &kanaly)
	if err != nil {
		t.Fatalf("reguły alarmowania nie ma w bazie: %v", err)
	}
	if automatykaID != id || wyzwalacz != string(shared.AutomationAlertTriggerFailure) {
		t.Fatalf("reguła w bazie wskazuje automatykę %d i wyzwalacz %q", automatykaID, wyzwalacz)
	}
	adresaci := []string{}
	if err := json.Unmarshal([]byte(kanaly), &adresaci); err != nil {
		t.Fatalf("kanały reguły są w bazie nieczytelne: %v", err)
	}
	if len(adresaci) != 2 {
		t.Fatalf("w bazie stoi %d kanałów reguły, oczekiwano 2", len(adresaci))
	}
}
