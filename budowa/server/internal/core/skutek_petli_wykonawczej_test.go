package core

import (
	"context"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Sprawdziany SKUTKU pętli wykonawczej modułu Studio.
//
// Mierzą to, co po komendzie ZOSTAJE — stan zadań odczytany osobnym wywołaniem,
// bilans przebiegu, powód odmowy — a nie kopertę odpowiedzi. Każdy z nich
// wyklucza jedną konkretną szkodę, którą ten produkt już raz poniósł:
//
//  1. plan „uruchomiony", który nie wykonał ani jednego zadania i nie powiedział
//     o tym ani słowa — `status: ok` z pustym skutkiem;
//  2. pętla puszczona przy WYŁĄCZONEJ nastawie: albo cicho nic nie robi, albo
//     melduje powodzenie, którego nie ma;
//  3. wsad na wielu dokumentach meldujący „przyjęto 10" i milczący o tym, co się
//     stało z każdym z nich osobno;
//  4. zadanie pominięte przez blokadę fragmentu, o którym Operator nie dowiaduje
//     się z kolejki, bo powód został przemilczany;
//  5. zatrzymanie, po którym zadanie w biegu ZNIKA, zamiast powiedzieć, że
//     zostało przerwane.

// petlaOknoSprawdzianu jest oknem, w którym stoją dokumenty tych sprawdzianów.
const petlaOknoSprawdzianu = "okno-petli-wykonawczej"

// petlaWlaczNastawe włącza pętlę wykonawczą na zasięgu okna dokumentu.
//
// Nastawa idzie drogą Operatora — komendą `studio.agents.settings.set` — a nie
// zapisem do magazynu na skróty. Sprawdzian, który włączałby pętlę inaczej niż
// Operator, mierzyłby drogę, której w produkcie nie ma.
func petlaWlaczNastawe(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno string, wieluAgentow bool) {
	t.Helper()

	prawda := true
	zasieg := shared.ConfigScope(shared.ConfigScopeWindow)
	var odpowiedz shared.StudioAgentsSettingsSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioAgentsSettingsSet,
		shared.StudioAgentsSettingsSetRequest{
			Scope: &zasieg, ScopeId: &okno,
			ExecutionLoopEnabled: &prawda,
			MultiAgentEnabled:    &wieluAgentow,
		}, &odpowiedz)
	if !odpowiedz.Settings.ExecutionLoopEnabled {
		t.Fatalf("nastawa pętli po zapisie nadal stoi na „nie”: %+v", odpowiedz.Settings)
	}
}

// petlaRozkladDokumentu odczytuje rozkład OSOBNYM wywołaniem — to jest miara
// właściwa, bo Operator zobaczy stan zadań przez `studio.plan.get`, a nie przez
// odpowiedź uruchomienia.
func petlaRozkladDokumentu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	kodRozkladu string) shared.StudioTaskPlan {
	t.Helper()

	var odpowiedz shared.StudioPlanGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanGet,
		shared.StudioPlanGetRequest{PlanId: &kodRozkladu}, &odpowiedz)
	return odpowiedz.Plan
}

// petlaZadaniePoNazwie wyszukuje zadanie rozkładu po jego nazwie.
func petlaZadaniePoNazwie(rozklad shared.StudioTaskPlan,
	nazwa string) (shared.StudioDocumentTask, bool) {

	for _, zadanie := range rozklad.Tasks {
		if zadanie.Title == nazwa {
			return zadanie, true
		}
	}
	return shared.StudioDocumentTask{}, false
}

// ── Szkoda pierwsza i druga: pętla wyłączona nastawą ────────────────────────

// TestPetlaPrzyWylaczonejNastawieOdmawiaNazywajacBrak wykazuje, że uruchomienie
// pętli przy wyłączonej nastawie wraca ODMOWĄ NAZYWAJĄCĄ BRAK NASTAWY — i że
// przy tym NIC nie robi.
//
// Dwie miary naraz, bo pojedyncza dałaby się obejść. Powód odmowy sprawdza, że
// Operator dowiaduje się, czego brakuje; stan zadań odczytany osobnym wywołaniem
// sprawdza, że pętla naprawdę nie tknęła kolejki. Odpowiedź „started: false"
// przy zadaniach przestawionych na „done" byłaby kłamstwem w drugą stronę.
func TestPetlaPrzyWylaczonejNastawieOdmawiaNazywajacBrak(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, petlaOknoSprawdzianu,
		"Pismo w sprawie rozgraniczenia nieruchomości.")

	var rozlozony shared.StudioPlanCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanCreate,
		shared.StudioPlanCreateRequest{
			DocumentId: dokument.Id,
			Order: "napisz pismo w tej sprawie na podstawie tych materiałów, " +
				"sprawdź terminologię i przygotuj wersję do druku",
		}, &rozlozony)

	// Rozkład sam jest już skutkiem: zlecenie jednym zdaniem daje kilka zadań.
	if len(rozlozony.Plan.Tasks) < 3 {
		t.Fatalf("zlecenie o czterech czynnościach rozłożyło się na %d zadań — "+
			"rozkład nie rozpoznaje czynności nazwanych w zleceniu: %+v",
			len(rozlozony.Plan.Tasks), rozlozony.Plan.Tasks)
	}

	var puszczony shared.StudioPlanRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanRun,
		shared.StudioPlanRunRequest{PlanId: rozlozony.Plan.Id}, &puszczony)

	if puszczony.Started {
		t.Fatal("pętla ruszyła przy WYŁĄCZONEJ nastawie — narzędzie nieuruchamiane " +
			"na starcie uruchomiło się samo")
	}
	if puszczony.RefusalReason == nil || strings.TrimSpace(*puszczony.RefusalReason) == "" {
		t.Fatal("pętla nie ruszyła i NIE POWIEDZIAŁA dlaczego — cisza w tym miejscu " +
			"jest dokładnie tym wzorcem szkody, który ten sprawdzian ma wykluczyć")
	}
	powod := *puszczony.RefusalReason
	if !strings.Contains(powod, "executionLoopEnabled") {
		t.Errorf("powód odmowy nie nazywa nastawy, której brakuje: %q", powod)
	}
	if !strings.Contains(powod, "studio.agents.settings.set") {
		t.Errorf("powód odmowy nie mówi, CZYM nastawę włączyć: %q", powod)
	}

	// Miara druga: kolejka nietknięta. Odczyt osobnym wywołaniem, nie z odpowiedzi.
	po := petlaRozkladDokumentu(t, zmontowany, zycie, rozlozony.Plan.Id)
	for _, zadanie := range po.Tasks {
		if zadanie.State != shared.StudioTaskStatePending {
			t.Errorf("zadanie %q stoi w stanie %q po odmowie uruchomienia — pętla "+
				"wyłączona nastawą tknęła kolejkę", zadanie.Title, zadanie.State)
		}
	}
}

// ── Szkoda pierwsza: plan uruchomiony bez pracy ─────────────────────────────

// TestPetlaUruchomionaWykonujeZadaniaAlboNazywaBrak wykazuje, że pętla puszczona
// przy WŁĄCZONEJ nastawie kończy każde zadanie stanem rozstrzygniętym — nigdy nie
// zostawia zadania w czekaniu i nie melduje powodzenia bez pracy.
//
// Zadanie rodzaju `export` jedzie rodziną wydania i wykonuje się do skutku.
// Zadania treści jadą kanałem modelu, którego stanowisko sprawdzianu nie ma —
// i właśnie dlatego są tu miarą najostrzejszą: pętla ma je zamknąć stanem
// `failed` z powodem NAZYWAJĄCYM brak, a nie zameldować „gotowe".
func TestPetlaUruchomionaWykonujeZadaniaAlboNazywaBrak(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, petlaOknoSprawdzianu,
		"Uzasadnienie wniosku o wpis do rejestru.\n\nCzęść druga uzasadnienia.")
	petlaWlaczNastawe(t, zmontowany, zycie, petlaOknoSprawdzianu, false)

	var rozlozony shared.StudioPlanCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanCreate,
		shared.StudioPlanCreateRequest{
			DocumentId: dokument.Id,
			Order:      "popraw język i wydaj do pdf",
			Tasks: []byte(`[
				{"kind":"proofread","title":"Poprawa językowa","instruction":"popraw język"},
				{"kind":"export","title":"Wydanie do pdf","instruction":"wydaj do pdf"}
			]`),
		}, &rozlozony)
	if len(rozlozony.Plan.Tasks) != 2 {
		t.Fatalf("zadania podane wprost nie weszły do rozkładu: %+v", rozlozony.Plan.Tasks)
	}

	var puszczony shared.StudioPlanRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanRun,
		shared.StudioPlanRunRequest{PlanId: rozlozony.Plan.Id}, &puszczony)
	if !puszczony.Started {
		t.Fatalf("pętla nie ruszyła przy włączonej nastawie: %v", puszczony.RefusalReason)
	}
	if puszczony.Balance == nil {
		t.Fatal("przebieg nie oddał bilansu — Operator nie dowie się, co pętla zrobiła")
	}

	po := petlaRozkladDokumentu(t, zmontowany, zycie, rozlozony.Plan.Id)

	// Miara pierwsza: ŻADNE zadanie nie zostało w czekaniu. Plan „uruchomiony"
	// z zadaniami nadal czekającymi to plan, który nie wykonał niczego.
	for _, zadanie := range po.Tasks {
		if zadanie.State == shared.StudioTaskStatePending ||
			zadanie.State == shared.StudioTaskStateRunning {

			t.Errorf("zadanie %q zostało w stanie %q po przebiegu — pętla zameldowała "+
				"uruchomienie i nie rozstrzygnęła zadania", zadanie.Title, zadanie.State)
		}
	}

	// Miara druga: każde zadanie NIEUDANE ma powód. Stan `failed` bez powodu jest
	// tym samym co cisza — Operator nie wie, czego brakuje.
	for _, zadanie := range po.Tasks {
		if zadanie.State != shared.StudioTaskStateFailed {
			continue
		}
		if zadanie.FailureReason == nil || strings.TrimSpace(*zadanie.FailureReason) == "" {
			t.Errorf("zadanie %q jest nieudane BEZ POWODU", zadanie.Title)
		}
	}

	// Miara trzecia: bilans przebiegu liczy pominięcia, a nie milczy o nich.
	if puszczony.Balance.Applied == 0 && puszczony.Balance.SkippedCount == 0 {
		t.Error("bilans przebiegu mówi zero wykonanych i zero pominiętych — " +
			"przebieg bez ani jednej liczby jest przebiegiem przemilczanym")
	}

	// Miara czwarta: stan rozkładu jest rozstrzygnięty, nie zostawiony w biegu.
	if po.State == shared.StudioPlanStateRunning {
		t.Errorf("rozkład został w stanie „running” po zakończonym przebiegu")
	}
}

// TestPetlaWydajeDokumentDoFormatuJakoZadanie wykazuje, że zadanie rodzaju
// `export` naprawdę wydaje dokument — a nie jest odkładane jako brak.
//
// Miara jest skutkiem: po przebiegu wydanie musi istnieć jako wynik zadania albo
// jako zasób w magazynie rdzenia. Zadanie zamknięte stanem `done` bez ani jednego
// śladu wydania byłoby meldunkiem bez pracy.
func TestPetlaWydajeDokumentDoFormatuJakoZadanie(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, petlaOknoSprawdzianu,
		"# Notatka służbowa\n\nTreść notatki w dwóch akapitach.\n\nAkapit drugi.")
	petlaWlaczNastawe(t, zmontowany, zycie, petlaOknoSprawdzianu, false)

	var rozlozony shared.StudioPlanCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanCreate,
		shared.StudioPlanCreateRequest{
			DocumentId: dokument.Id,
			Order:      "wydaj notatkę do txt",
			Tasks: []byte(`[{"kind":"export","title":"Wydanie do txt",` +
				`"instruction":"wydaj do txt"}]`),
		}, &rozlozony)

	var puszczony shared.StudioPlanRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanRun,
		shared.StudioPlanRunRequest{PlanId: rozlozony.Plan.Id}, &puszczony)
	if !puszczony.Started {
		t.Fatalf("pętla nie ruszyła: %v", puszczony.RefusalReason)
	}

	po := petlaRozkladDokumentu(t, zmontowany, zycie, rozlozony.Plan.Id)
	zadanie, jest := petlaZadaniePoNazwie(po, "Wydanie do txt")
	if !jest {
		t.Fatalf("zadania wydania nie ma w rozkładzie: %+v", po.Tasks)
	}
	if zadanie.State != shared.StudioTaskStateDone {
		powod := ""
		if zadanie.FailureReason != nil {
			powod = *zadanie.FailureReason
		}
		t.Fatalf("zadanie wydania skończyło w stanie %q (%s) — pętla nie wykonała "+
			"wydania, choć rodzina wydania stoi w rdzeniu", zadanie.State, powod)
	}
	if puszczony.Balance == nil || puszczony.Balance.Applied == 0 {
		t.Error("bilans przebiegu nie liczy wykonanego wydania — zadanie „done” " +
			"bez wpisu w bilansie jest meldunkiem bez pracy")
	}
}

// ── Szkoda piąta: zatrzymanie, po którym zadanie znika ──────────────────────

// TestZatrzymaniePetliZostawiaSladPrzerwaniaINiegubiPracy wykazuje trzy rzeczy
// naraz: rozkład wchodzi w stan zatrzymany, powód zatrzymania jest zapisany,
// a zadania już domknięte ZOSTAJĄ domknięte.
//
// Zatrzymanie, które kasuje pracę wykonaną przed nim, byłoby gorsze niż brak
// zatrzymania — Operator naciskający „przerwij" traciłby to, co pętla zrobiła
// dobrze.
func TestZatrzymaniePetliZostawiaSladPrzerwaniaINiegubiPracy(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, petlaOknoSprawdzianu,
		"Treść dokumentu do przebiegu zatrzymywanego.")
	petlaWlaczNastawe(t, zmontowany, zycie, petlaOknoSprawdzianu, false)

	var rozlozony shared.StudioPlanCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanCreate,
		shared.StudioPlanCreateRequest{
			DocumentId: dokument.Id,
			Order:      "wydaj do txt, potem przejrzyj",
			Tasks: []byte(`[
				{"kind":"export","title":"Wydanie przed zatrzymaniem","instruction":"wydaj do txt"},
				{"kind":"review","title":"Przejrzenie po wydaniu","instruction":"przejrzyj"}
			]`),
		}, &rozlozony)

	var puszczony shared.StudioPlanRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanRun,
		shared.StudioPlanRunRequest{PlanId: rozlozony.Plan.Id}, &puszczony)
	if !puszczony.Started {
		t.Fatalf("pętla nie ruszyła: %v", puszczony.RefusalReason)
	}
	przed := petlaRozkladDokumentu(t, zmontowany, zycie, rozlozony.Plan.Id)
	wydanie, jest := petlaZadaniePoNazwie(przed, "Wydanie przed zatrzymaniem")
	if !jest {
		t.Fatal("zadania wydania nie ma w rozkładzie")
	}
	stanPrzed := wydanie.State

	powodZatrzymania := "Operator przerwał przebieg w oknie pętli"
	var zatrzymany shared.StudioPlanStopResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanStop,
		shared.StudioPlanStopRequest{PlanId: rozlozony.Plan.Id, Reason: &powodZatrzymania},
		&zatrzymany)

	po := petlaRozkladDokumentu(t, zmontowany, zycie, rozlozony.Plan.Id)
	if po.State != shared.StudioPlanStateStopped {
		t.Errorf("rozkład po zatrzymaniu stoi w stanie %q, a nie „stopped”", po.State)
	}
	if po.StopReason == nil || !strings.Contains(*po.StopReason, "przerwał") {
		t.Errorf("rozkład po zatrzymaniu nie niesie powodu Operatora: %v", po.StopReason)
	}
	// Praca sprzed zatrzymania ZOSTAJE. Miara jest wprost: stan zadania wydania
	// nie może się po zatrzymaniu pogorszyć.
	poZatrzymaniu, _ := petlaZadaniePoNazwie(po, "Wydanie przed zatrzymaniem")
	if stanPrzed == shared.StudioTaskStateDone &&
		poZatrzymaniu.State != shared.StudioTaskStateDone {

		t.Errorf("zadanie domknięte przed zatrzymaniem zmieniło stan na %q — "+
			"zatrzymanie skasowało pracę, która się udała", poZatrzymaniu.State)
	}
	// Żadne zadanie nie zniknęło: liczba zadań przed i po jest ta sama.
	if len(po.Tasks) != len(przed.Tasks) {
		t.Errorf("po zatrzymaniu rozkład ma %d zadań, przed miał %d — zadanie zniknęło",
			len(po.Tasks), len(przed.Tasks))
	}
}

// ── Szkoda czwarta: blokada fragmentu przemilczana w kolejce ────────────────

// TestBlokadaFragmentuWidocznaWKolejceZadan wykazuje, że zadanie zatrzymane
// blokadą fragmentu jest WIDOCZNE w kolejce wraz z powodem, a nie przemilczane.
//
// Sprawdzian mierzy przy okazji rzecz nośną: pętla jedzie REJESTREM, nie
// wywołaniem adaptera wprost. Zapora blokad stoi w rejestrze — pętla omijająca
// rejestr przepisałaby zablokowany fragment i sprawdzian by tego nie zobaczył,
// bo zadanie skończyłoby się „gotowe".
func TestBlokadaFragmentuWidocznaWKolejceZadan(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, petlaOknoSprawdzianu,
		"Podstawa prawna: art. 5 ustawy. Dalsza treść pisma do poprawy.")
	petlaWlaczNastawe(t, zmontowany, zycie, petlaOknoSprawdzianu, false)

	nazwaBlokady := "podstawa prawna — nie zmieniać"
	powodBlokady := "cytat z ustawy musi zostać dosłownie"
	var zalozona shared.StudioLockAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioLockAdd,
		shared.StudioLockAddRequest{
			DocumentId: dokument.Id, RangeStart: 0, RangeEnd: 30,
			Name: nazwaBlokady, Reason: &powodBlokady,
		}, &zalozona)

	var rozlozony shared.StudioPlanCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanCreate,
		shared.StudioPlanCreateRequest{
			DocumentId: dokument.Id,
			Order:      "przepisz podstawę prawną własnymi słowami",
			Tasks: []byte(`[{"kind":"draft","title":"Przepisanie podstawy prawnej",` +
				`"instruction":"przepisz","rangeStart":0,"rangeEnd":30}]`),
		}, &rozlozony)

	var puszczony shared.StudioPlanRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanRun,
		shared.StudioPlanRunRequest{PlanId: rozlozony.Plan.Id}, &puszczony)
	if !puszczony.Started {
		t.Fatalf("pętla nie ruszyła: %v", puszczony.RefusalReason)
	}

	po := petlaRozkladDokumentu(t, zmontowany, zycie, rozlozony.Plan.Id)
	zadanie, jest := petlaZadaniePoNazwie(po, "Przepisanie podstawy prawnej")
	if !jest {
		t.Fatalf("zadania nie ma w rozkładzie: %+v", po.Tasks)
	}
	if zadanie.State == shared.StudioTaskStateDone {
		t.Fatal("zadanie godzące w zablokowany fragment skończyło się „gotowe” — " +
			"pętla ominęła zaporę blokad")
	}
	if zadanie.FailureReason == nil {
		t.Fatal("zadanie zatrzymane blokadą NIE MA powodu w kolejce — Operator " +
			"nie dowie się, że fragment został pominięty")
	}
	powod := *zadanie.FailureReason
	if !strings.Contains(powod, nazwaBlokady) {
		t.Errorf("powód przy zadaniu nie nazywa blokady, która je zatrzymała: %q", powod)
	}

	// Bilans przebiegu też o tym mówi — bo pominięcie dotyczy przebiegu, nie
	// tylko jednego wiersza kolejki.
	if puszczony.Balance == nil || puszczony.Balance.SkippedCount == 0 {
		t.Fatal("bilans przebiegu nie liczy pominięcia — zapora zatrzymała zadanie, " +
			"a przebieg tego nie zgłosił")
	}
	wBilansie := false
	for _, pozycja := range puszczony.Balance.Skipped {
		if strings.Contains(pozycja.Reason, nazwaBlokady) {
			wBilansie = true
		}
	}
	if !wBilansie {
		t.Errorf("bilans przebiegu nie nazywa blokady: %+v", puszczony.Balance.Skipped)
	}

	// Skutek na TREŚCI: zablokowany fragment został dosłownie taki, jaki był.
	// To jest miara ostateczna — reszta mówi o meldunkach, ta o dokumencie.
	tresc := trescDokumentu(t, zmontowany, zycie, petlaOknoSprawdzianu, dokument.Id)
	if !strings.HasPrefix(tresc, "Podstawa prawna: art. 5 ustawy.") {
		t.Errorf("zablokowany fragment został zmieniony przez pętlę; treść po przebiegu: %q",
			tresc)
	}
}

// ── Szkoda trzecia: wsad milczący o poszczególnych dokumentach ──────────────

// TestWsadOddajeWynikKazdegoDokumentuOsobno wykazuje, że wsad na wielu
// dokumentach rozlicza się z KAŻDEGO osobno — a nie jedną liczbą „przyjęto".
//
// Miara jest zupełna: każdy dokument z żądania musi się znaleźć dokładnie raz
// po stronie przyjętych albo odrzuconych, a każde odrzucenie musi mieć powód.
// Dokument nieistniejący musi dostać powód INNY niż dokumenty istniejące —
// inaczej „odrzucono" byłoby jednym workiem na wszystko.
func TestWsadOddajeWynikKazdegoDokumentuOsobno(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	pierwszy := dokumentZTrescia(t, zmontowany, zycie, petlaOknoSprawdzianu,
		"Pierwszy dokument wsadu.")
	drugi := dokumentZTrescia(t, zmontowany, zycie, petlaOknoSprawdzianu,
		"Drugi dokument wsadu.")
	nieistniejacy := "studio-dok-nie-ma-takiego"

	var wsad shared.StudioBatchRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioBatchRun,
		shared.StudioBatchRunRequest{
			WindowId:    petlaOknoSprawdzianu,
			DocumentIds: []string{pierwszy.Id, drugi.Id, nieistniejacy},
			ActionId:    "studio.korekta.ortografia",
		}, &wsad)

	rozliczone := map[string]string{}
	for _, odrzucenie := range wsad.Rejected {
		if strings.TrimSpace(odrzucenie.Reason) == "" {
			t.Errorf("dokument %s odrzucony BEZ POWODU", odrzucenie.DocumentId)
		}
		if _, byl := rozliczone[odrzucenie.DocumentId]; byl {
			t.Errorf("dokument %s rozliczony dwa razy", odrzucenie.DocumentId)
		}
		rozliczone[odrzucenie.DocumentId] = odrzucenie.Reason
	}
	if wsad.Accepted+len(wsad.Rejected) != 3 {
		t.Errorf("wsad z trzech dokumentów rozliczył %d przyjętych i %d odrzuconych — "+
			"suma nie zgadza się z żądaniem, więc któryś dokument został przemilczany",
			wsad.Accepted, len(wsad.Rejected))
	}

	// Dokument nieistniejący musi mieć powód WŁASNY, nazywający brak dokumentu.
	powodNieistniejacego, byl := rozliczone[nieistniejacy]
	if !byl {
		if wsad.Accepted == 0 {
			t.Fatal("dokumentu nieistniejącego nie ma ani wśród przyjętych, ani wśród " +
				"odrzuconych — wsad go przemilczał")
		}
		t.Fatalf("dokument nieistniejący nie został odrzucony: przyjęto %d", wsad.Accepted)
	}
	if !strings.Contains(powodNieistniejacego, "nie istnieje") {
		t.Errorf("powód odrzucenia dokumentu nieistniejącego nie nazywa braku: %q",
			powodNieistniejacego)
	}
	// Dokumenty istniejące, gdy odmówiły, odmówiły z INNEGO powodu niż brak
	// dokumentu. Jeden powód dla wszystkiego znaczyłby bilans pozorny.
	for _, kod := range []string{pierwszy.Id, drugi.Id} {
		powod, odrzucony := rozliczone[kod]
		if !odrzucony {
			continue
		}
		if powod == powodNieistniejacego {
			t.Errorf("dokument istniejący %s dostał ten sam powód co nieistniejący (%q) — "+
				"bilans wsadu nie rozróżnia przypadków", kod, powod)
		}
	}
}

// ── Rozkład zlecenia: rozpoznanie czynności ze słów Operatora ───────────────

// TestRozkladZleceniaRozpoznajeCzynnosciZeSlowOperatora wykazuje, że rozkład
// jest rachunkiem, a nie jednym zadaniem z całym zleceniem w środku.
//
// Miara jest podwójna: liczba zadań i ICH RODZAJE. Rozkład oddający cztery
// zadania rodzaju `custom` byłby rozkładem pozornym — nazwałby czynności, których
// nie rozpoznał.
func TestRozkladZleceniaRozpoznajeCzynnosciZeSlowOperatora(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, petlaOknoSprawdzianu,
		"Materiał wejściowy pisma.")

	var rozlozony shared.StudioPlanCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanCreate,
		shared.StudioPlanCreateRequest{
			DocumentId: dokument.Id,
			Order: "napisz pismo w tej sprawie na podstawie tych materiałów, " +
				"sprawdź terminologię i przygotuj wersję do druku",
		}, &rozlozony)

	rodzaje := map[shared.StudioTaskKind]bool{}
	for _, zadanie := range rozlozony.Plan.Tasks {
		rodzaje[zadanie.Kind] = true
	}
	for _, oczekiwany := range []shared.StudioTaskKind{
		shared.StudioTaskKindResearch,
		shared.StudioTaskKindDraft,
		shared.StudioTaskKindProofread,
		shared.StudioTaskKindFormat,
	} {
		if !rodzaje[oczekiwany] {
			t.Errorf("rozkład nie rozpoznał czynności rodzaju %q w zleceniu; rozpoznane: %v",
				oczekiwany, rodzaje)
		}
	}

	// Zależności: rozkład zlecenia dokumentowego jest szeregowy — nie da się
	// poprawić języka pisma, którego jeszcze nie ma. Pierwsze zadanie na niczym
	// nie stoi, każde następne stoi na poprzednim.
	for numer, zadanie := range rozlozony.Plan.Tasks {
		if numer == 0 {
			if len(zadanie.DependsOn) != 0 {
				t.Errorf("pierwsze zadanie rozkładu stoi na %v — nie ma na czym stać",
					zadanie.DependsOn)
			}
			continue
		}
		if len(zadanie.DependsOn) == 0 {
			t.Errorf("zadanie %q nie stoi na niczym — rozkład bez zależności puściłby "+
				"korektę przed napisaniem treści", zadanie.Title)
			continue
		}
		poprzednie := rozlozony.Plan.Tasks[numer-1].Id
		if zadanie.DependsOn[0] != poprzednie {
			t.Errorf("zadanie %q stoi na %q, a nie na poprzednim zadaniu rozkładu",
				zadanie.Title, zadanie.DependsOn[0])
		}
	}
}

// TestZadanieRozkladuPrzestawiaSieTylkoDrogaDozwolona wykazuje, że przejścia
// stanów są sprawdzane, a nie przyjmowane na słowo.
//
// Zadanie domknięte, które wraca do biegu bez decyzji o ponowieniu, zgubiłoby
// wynik swojej pracy — dlatego to przejście musi wracać odmową NAZYWAJĄCĄ oba
// stany, a nie cichym przyjęciem.
func TestZadanieRozkladuPrzestawiaSieTylkoDrogaDozwolona(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, petlaOknoSprawdzianu,
		"Dokument do przestawiania stanów zadania.")

	var rozlozony shared.StudioPlanCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanCreate,
		shared.StudioPlanCreateRequest{
			DocumentId: dokument.Id,
			Order:      "zadanie własne Operatora",
			Tasks:      []byte(`[{"kind":"custom","title":"Zadanie własne"}]`),
		}, &rozlozony)
	pierwsze := rozlozony.Plan.Tasks[0]

	// Droga dozwolona: czekające → pominięte decyzją Operatora.
	pominiete := shared.StudioTaskState(shared.StudioTaskStateSkipped)
	var przestawione shared.StudioPlanTaskUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPlanTaskUpdate,
		shared.StudioPlanTaskUpdateRequest{TaskId: pierwsze.Id, State: &pominiete},
		&przestawione)
	if przestawione.Task.State != shared.StudioTaskStateSkipped {
		t.Fatalf("zadanie nie przeszło w stan pominięty: %q", przestawione.Task.State)
	}

	// Droga niedozwolona: pominięte → domknięte. Zadanie, którego nikt nie
	// wykonał, nie ma prawa stać się wykonanym.
	domkniete := shared.StudioTaskState(shared.StudioTaskStateDone)
	blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandStudioPlanTaskUpdate,
		shared.StudioPlanTaskUpdateRequest{TaskId: pierwsze.Id, State: &domkniete})
	if !strings.Contains(blad.Message, "skipped") || !strings.Contains(blad.Message, "done") {
		t.Errorf("odmowa przejścia nie nazywa obu stanów: %q", blad.Message)
	}

	// Stan zadania po odmowie jest nietknięty — odczyt osobnym wywołaniem.
	po := petlaRozkladDokumentu(t, zmontowany, zycie, rozlozony.Plan.Id)
	if po.Tasks[0].State != shared.StudioTaskStateSkipped {
		t.Errorf("zadanie po odmowionym przejściu stoi w stanie %q", po.Tasks[0].State)
	}
}
