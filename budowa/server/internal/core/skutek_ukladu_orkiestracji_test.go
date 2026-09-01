// Skutek dopełnień układu: czy bramka, grupa, kompensacja i spięcie zostawiają po sobie wiersz w bazie danych.
package core

import (
	"context"
	"testing"

	"danacoconsole/shared"
)

// automatykaUkladuSprawdzianu zakłada automatykę z trzema krokami i dwoma
// torami schodzącymi się w kroku trzecim. Oddaje jej kod.
func automatykaUkladuSprawdzianu(t *testing.T, zmontowany *Zmontowany,
	zycie context.Context) string {
	t.Helper()

	var zapis shared.AutomationWorkflowSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAutomationWorkflowSave,
		shared.AutomationWorkflowSaveRequest{
			Name: "Układ sprawdzianu",
			Steps: []shared.AutomationStep{
				{Id: "krok-a", Kind: shared.AutomationStepKindCommand},
				{Id: "krok-b", Kind: shared.AutomationStepKindCommand},
				{Id: "krok-scalajacy", Kind: shared.AutomationStepKindCommand},
			},
		}, &zapis)

	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationDependencySet,
		shared.OrchestrationDependencySetRequest{
			WorkflowId: zapis.Workflow.Id,
			Dependency: shared.AutomationDependency{
				FromStepId: "krok-a", ToStepId: "krok-scalajacy",
				Kind: shared.AutomationDependencyKindParallel,
			},
		}, nil)
	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationDependencySet,
		shared.OrchestrationDependencySetRequest{
			WorkflowId: zapis.Workflow.Id,
			Dependency: shared.AutomationDependency{
				FromStepId: "krok-b", ToStepId: "krok-scalajacy",
				Kind: shared.AutomationDependencyKindParallel,
			},
		}, nil)
	return zapis.Workflow.Id
}

// TestSkutekBramkiDolaczeniaWBazie mierzy orchestration.gate.set wraz z jej trwałym zapisem w tej bazie.
func TestSkutekBramkiDolaczeniaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	uklad := automatykaUkladuSprawdzianu(t, zmontowany, zycie)

	dwa := 2
	var wynik shared.OrchestrationGateSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationGateSet,
		shared.OrchestrationGateSetRequest{
			WorkflowId: uklad, StepId: "krok-scalajacy",
			Rule: shared.OrchestrationGateRuleCount, Count: &dwa,
		}, &wynik)

	if !wynik.Valid {
		t.Errorf("bramka na dwóch dochodzących torach uznana za niepoprawną: %v", wynik.Issues)
	}
	liczba := liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM orkiestracja_bramka b JOIN automatyka a ON a.id = b.automatyka_id
		  WHERE a.identyfikator_zewnetrzny = ? AND b.krok = 'krok-scalajacy' AND b.regula = 'count' AND b.licznik = 2`,
		uklad)
	if liczba != 1 {
		t.Fatalf("bramka nie doszła do bazy: wierszy %d", liczba)
	}

	// Bramka licznikowa wymagająca więcej torów niż dochodzi zapisuje się, a zastrzeżenie wraca w wyniku.
	piec := 5
	var szeroka shared.OrchestrationGateSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationGateSet,
		shared.OrchestrationGateSetRequest{
			WorkflowId: uklad, StepId: "krok-scalajacy",
			Rule: shared.OrchestrationGateRuleCount, Count: &piec,
		}, &szeroka)
	if szeroka.Valid || len(szeroka.Issues) == 0 {
		t.Error("bramka wymagająca pięciu torów przy dwóch nie zgłosiła zastrzeżenia")
	}

	// Reguła licznikowa bez liczby torów jest błędem ŻĄDANIA, nie stanem układu.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandOrchestrationGateSet,
		shared.OrchestrationGateSetRequest{
			WorkflowId: uklad, StepId: "krok-scalajacy",
			Rule: shared.OrchestrationGateRuleCount,
		})
	if odmowa.Code == "" {
		t.Error("reguła licznikowa bez liczby torów przeszła bez odmowy")
	}
}

// TestSkutekGrupyKrokowWBazie mierzy `orchestration.group.set` wraz z jej
// usunięciem przez wykaz pusty.
func TestSkutekGrupyKrokowWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	uklad := automatykaUkladuSprawdzianu(t, zmontowany, zycie)

	var zapis shared.OrchestrationGroupSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationGroupSet,
		shared.OrchestrationGroupSetRequest{
			WorkflowId: uklad, Name: "Tory równoległe",
			StepIds: []string{"krok-a", "krok-b"},
			Kind:    shared.AutomationDependencyKindParallel,
		}, &zapis)
	if len(zapis.Groups) != 1 {
		t.Fatalf("odpowiedź niesie %d grup zamiast jednej", len(zapis.Groups))
	}
	kod := zapis.Groups[0].Id

	kroki := liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM orkiestracja_grupa_krok k
		   JOIN orkiestracja_grupa g ON g.id = k.grupa_id WHERE g.kod = ?`, kod)
	if kroki != 2 {
		t.Fatalf("skład grupy nie doszedł do bazy: kroków %d", kroki)
	}

	// Wykaz pusty usuwa grupę — tak stanowi kontrakt.
	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationGroupSet,
		shared.OrchestrationGroupSetRequest{
			WorkflowId: uklad, GroupId: &kod, Name: "Tory równoległe",
			StepIds: []string{}, Kind: shared.AutomationDependencyKindParallel,
		}, nil)
	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM orkiestracja_grupa WHERE kod = ?`, kod) != 0 {
		t.Fatal("grupa została w bazie mimo pustego wykazu kroków")
	}
}

// TestSkutekKompensacjiWBazie mierzy orchestration.compensation.set wraz z jej trwałym zapisem w bazie.
func TestSkutekKompensacjiWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	uklad := automatykaUkladuSprawdzianu(t, zmontowany, zycie)

	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationCompensationSet,
		shared.OrchestrationCompensationSetRequest{
			WorkflowId: uklad, StepId: "krok-a", CompensationStepId: tekstOpcjonalny("krok-b"),
		}, nil)
	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM orkiestracja_kompensacja k JOIN automatyka a ON a.id = k.automatyka_id
		  WHERE a.identyfikator_zewnetrzny = ? AND k.krok = 'krok-a' AND k.krok_wycofu = 'krok-b'`, uklad) != 1 {
		t.Fatal("kompensacja nie doszła do bazy")
	}

	// Krok wycofujący pusty zdejmuje kompensację.
	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationCompensationSet,
		shared.OrchestrationCompensationSetRequest{
			WorkflowId: uklad, StepId: "krok-a", CompensationStepId: nil,
		}, nil)
	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM orkiestracja_kompensacja k JOIN automatyka a ON a.id = k.automatyka_id
		  WHERE a.identyfikator_zewnetrzny = ?`, uklad) != 0 {
		t.Fatal("kompensacja została w bazie po zdjęciu")
	}

	// Krok wycofujący sam siebie nie wycofuje niczego.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandOrchestrationCompensationSet,
		shared.OrchestrationCompensationSetRequest{
			WorkflowId: uklad, StepId: "krok-a", CompensationStepId: tekstOpcjonalny("krok-a"),
		})
	if odmowa.Code == "" {
		t.Error("kompensacja kroku samym sobą przeszła bez odmowy")
	}
}

// TestSkutekSpieciaZMultitaskingiem mierzy rzecz, o którą naprawdę chodzi:
// czy spięcie PRZESTAWIA kolejki automatyki, czy tylko zapisuje znacznik.
func TestSkutekSpieciaZMultitaskingiem(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	uklad := automatykaUkladuSprawdzianu(t, zmontowany, zycie)

	// Kolejka automatyki jest zakładana wprost: sprawdzian mierzy skutek spięcia, nie jej powstanie.
	if _, err := baza.Exec(
		`INSERT INTO kolejka (nazwa, rodzaj, stan) VALUES (?, 'sesyjna', 'bezczynna')`,
		uklad); err != nil {
		t.Fatalf("nie można założyć kolejki sprawdzianu: %v", err)
	}

	var spiecie shared.OrchestrationMultitaskingLinkResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationMultitaskingLink,
		shared.OrchestrationMultitaskingLinkRequest{WorkflowId: uklad, Linked: true}, &spiecie)
	if !spiecie.Linked || len(spiecie.QueueIds) != 1 {
		t.Fatalf("odpowiedź spięcia: linked=%v kolejek=%d", spiecie.Linked, len(spiecie.QueueIds))
	}

	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM kolejka WHERE nazwa = ? AND rodzaj = 'multitasking'`, uklad) != 1 {
		t.Fatal("spięcie nie przestawiło rodzaju kolejki — znacznik bez skutku")
	}
	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM orkiestracja_spiecie_multitasking s
		   JOIN automatyka a ON a.id = s.automatyka_id WHERE a.identyfikator_zewnetrzny = ?`, uklad) != 1 {
		t.Fatal("zapis spięcia nie doszedł do bazy")
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandOrchestrationMultitaskingLink,
		shared.OrchestrationMultitaskingLinkRequest{WorkflowId: uklad, Linked: false}, nil)
	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM orkiestracja_spiecie_multitasking s
		   JOIN automatyka a ON a.id = s.automatyka_id WHERE a.identyfikator_zewnetrzny = ?`, uklad) != 0 {
		t.Fatal("zapis spięcia został w bazie po rozłączeniu")
	}
}
