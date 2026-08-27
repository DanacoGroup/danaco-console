// Pętla wykonawcza modułu Studio obsługuje pięć czynności rodziny
// studio.plan.*: rozkład zlecenia na zadania, odczyt rozkładu, przestawienie
// zadania, uruchomienie pętli i jej zatrzymanie; domyślnie jest wyłączona,
// włącza ją Operator nastawą.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Granice obiegów pętli, gdy nastawa ich nie podaje. Są granicami, nie
// zaleceniami: pętla bez górnego ograniczenia obiegów jest pętlą, która przy
// pierwszym zadaniu odmawiającym w kółko kręci się na koszt Operatora.
const (
	petlaObiegiDomyslnie     = 12
	petlaBezPostepuDomyslnie = 2
	petlaWykonawcowDomyslnie = 1
)

// adapterPetliStudia prowadzi pętlę wykonawczą modułu Studio: typ osobny od
// adapterStudia, bo pętla ma stan własny (magazyn rozkładów) i własną
// zależność (rejestr komend, którym czyta nastawy).
type adapterPetliStudia struct {
	studio  *adapterStudia
	magazyn *petlaMagazynStudia
	rejestr *Rejestr
	emiter  *emiter
}

// nowyAdapterPetliStudia wiąże pętlę z adapterem Studia i rejestrem komend,
// tworząc pusty magazyn rozkładów gotowy do przyjęcia pierwszego zlecenia.
func nowyAdapterPetliStudia(studio *adapterStudia, rejestr *Rejestr, e *emiter) *adapterPetliStudia {
	return &adapterPetliStudia{
		studio: studio, magazyn: petlaNowyMagazyn(), rejestr: rejestr, emiter: e,
	}
}

// ── studio.plan.create ──────────────────────────────────────────────────────

// RozlozZlecenie obsługuje studio.plan.create: rozkład idzie trzema drogami
// po kolei, a pierwsza, która da zadania, wygrywa — zadania wskazane wprost,
// rozkład ułożony modelem, rozkład własny rdzenia, który nie zawodzi nigdy.
func (p *adapterPetliStudia) RozlozZlecenie(ctx context.Context,
	z shared.StudioPlanCreateRequest) (shared.StudioPlanCreateResponse, error) {

	if strings.TrimSpace(z.Order) == "" {
		return shared.StudioPlanCreateResponse{}, bladWskazaniaStudio(
			"rozkład bez treści zlecenia — pętla nie ma czego rozkładać")
	}
	dokument, err := p.studio.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPlanCreateResponse{}, err
	}

	surowe := p.petlaZadaniaZlecenia(ctx, z, dokument.Okno)
	if len(surowe) == 0 {
		// Rozkład własny zawsze daje co najmniej jedno zadanie — odmowa stoi tu
		// jako zapora, nie ścieżka.
		return shared.StudioPlanCreateResponse{}, bladStudio(
			protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
				"moduł Studio: rozkład zlecenia nie dał ani jednego zadania")))
	}

	kod := nowyIdentyfikator(przedrostekRozkladuStudia)
	teraz := petlaTeraz()
	okno := dokument.Okno
	rozklad := &petlaRozkladStudia{
		Kod:            kod,
		KodDokumentu:   z.DocumentId,
		Zlecenie:       z.Order,
		Stan:           shared.StudioPlanStateDraft,
		Ulozyl:         petlaWykonawcaZlecenia(z.AgentId, z.AgentName, &okno),
		Utworzono:      teraz,
		Zaktualizowano: teraz,
	}
	// Zakres zlecenia schodzi na każde zadanie bez własnego zakresu, nie
	// tylko na pierwsze.
	for numer := range surowe {
		if surowe[numer].ZakresOd == nil && z.RangeStart != nil {
			surowe[numer].ZakresOd = z.RangeStart
		}
		if surowe[numer].ZakresDo == nil && z.RangeEnd != nil {
			surowe[numer].ZakresDo = z.RangeEnd
		}
	}
	rozklad.Zadania = petlaZlozZadania(kod, z.DocumentId, surowe)
	p.magazyn.petlaDolozRozklad(rozklad)

	return shared.StudioPlanCreateResponse{Plan: petlaZlozRozklad(rozklad, nil)}, nil
}

// petlaZadaniaZlecenia oddaje zadania rozkładu pierwszą drogą, która je da:
// wskazaniem Operatora, rozkładem modelu, rozkładem własnym rdzenia.
func (p *adapterPetliStudia) petlaZadaniaZlecenia(ctx context.Context,
	z shared.StudioPlanCreateRequest, okno string) []petlaZadanieStudia {

	if wprost := petlaZadaniaZLadunku(z.Tasks); len(wprost) > 0 {
		return wprost
	}
	if modelem := p.petlaRozkladModelem(ctx, z.Order, okno); len(modelem) > 0 {
		return modelem
	}
	return petlaRozlozZlecenie(z.Order)
}

// petlaZadaniaZLadunku czyta zadania podane wprost w kształcie
// StudioDocumentTask[]. Ładunek nieczytelny nie jest odmową: kontrakt nazywa
// to pole nieobowiązkowym, a rozkład idzie wtedy drogą następną.
func petlaZadaniaZLadunku(ladunek json.RawMessage) []petlaZadanieStudia {
	if len(ladunek) == 0 || strings.TrimSpace(string(ladunek)) == "null" {
		return nil
	}
	var podane []shared.StudioDocumentTask
	if err := json.Unmarshal(ladunek, &podane); err != nil {
		return nil
	}
	zadania := make([]petlaZadanieStudia, 0, len(podane))
	for numer, wpis := range podane {
		if strings.TrimSpace(wpis.Title) == "" && strings.TrimSpace(string(wpis.Kind)) == "" {
			continue
		}
		zadanie := petlaZadanieStudia{
			Kolejnosc: wpis.Order,
			Rodzaj:    wpis.Kind,
			Nazwa:     wpis.Title,
			ZakresOd:  wpis.RangeStart,
			ZakresDo:  wpis.RangeEnd,
			Stan:      shared.StudioTaskStatePending,
			StoiNa:    append([]string(nil), wpis.DependsOn...),
		}
		if zadanie.Kolejnosc == 0 {
			zadanie.Kolejnosc = numer + 1
		}
		if wpis.Instruction != nil {
			zadanie.Polecenie = *wpis.Instruction
		}
		zadania = append(zadania, zadanie)
	}
	return zadania
}

// petlaRozkladModelem prosi model o rozkład zlecenia na zadania w postaci
// tablicy JSON o polach kind, title i instruction. Odpowiedź innego kształtu
// nie jest odmową — rozkład idzie wtedy drogą własną rdzenia.
func (p *adapterPetliStudia) petlaRozkladModelem(ctx context.Context,
	zlecenie, okno string) []petlaZadanieStudia {

	if p.studio == nil || p.studio.kanaly == nil || p.studio.okna == nil || okno == "" {
		return nil
	}
	opisOkna, err := p.studio.okna.Okno(okno)
	if err != nil || opisOkna.KanalModelu == "" {
		return nil
	}

	var odpowiedz strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			odpowiedz.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	rodzaje := make([]string, 0, len(shared.WartosciStudioTaskKind()))
	for _, rodzaj := range shared.WartosciStudioTaskKind() {
		rodzaje = append(rodzaje, string(rodzaj))
	}
	tresc := "Rozłóż zlecenie dokumentowe na zadania jednostkowe.\n\n" +
		"Zlecenie Operatora:\n" + zlecenie + "\n\n" +
		"Oddaj WYŁĄCZNIE tablicę JSON. Każdy wpis ma pola: " +
		"\"kind\" (jedna z wartości: " + strings.Join(rodzaje, ", ") + "), " +
		"\"title\" (zadanie nazwane jednym zdaniem), " +
		"\"instruction\" (polecenie dla wykonawcy). " +
		"Zadania ustaw w kolejności wykonania. Bez komentarza wokół tablicy."

	if err := p.studio.kanaly.Wyslij(ctx, models.Zapytanie{
		Zasiegi:   models.Zasiegi{Okno: okno},
		Wiadomosc: "studio.plan.create",
		Tresc:     tresc,
		Kanal:     opisOkna.KanalModelu,
	}, ujscie); err != nil {
		return nil
	}
	return petlaZadaniaZLadunku(petlaWyjmijTablice(odpowiedz.String()))
}

// petlaWyjmijTablice wyłuskuje tablicę JSON z odpowiedzi modelu. Model bywa
// rozmowny i opakowuje tablicę zdaniem albo płotkiem ```json — wyjmujemy ją po
// pierwszym `[` i ostatnim `]`, zamiast odrzucać odpowiedź za obudowę.
func petlaWyjmijTablice(tresc string) json.RawMessage {
	od := strings.Index(tresc, "[")
	do := strings.LastIndex(tresc, "]")
	if od < 0 || do <= od {
		return nil
	}
	return json.RawMessage(tresc[od : do+1])
}

// ── studio.plan.get ─────────────────────────────────────────────────────────

// Rozklad obsługuje studio.plan.get: oddaje rozkład wskazany identyfikatorem
// albo wykaz rozkładów dokumentu wraz z rozkładem najświeższym.
func (p *adapterPetliStudia) Rozklad(_ context.Context,
	z shared.StudioPlanGetRequest) (shared.StudioPlanGetResponse, error) {

	if z.PlanId != nil && strings.TrimSpace(*z.PlanId) != "" {
		rozklad, jest := p.magazyn.petlaRozkladPoKodzie(strings.TrimSpace(*z.PlanId))
		if !jest {
			return shared.StudioPlanGetResponse{},
				bladBrakuStudio("rozkład nie istnieje: " + *z.PlanId)
		}
		return shared.StudioPlanGetResponse{Plan: petlaZlozRozklad(rozklad, z.State)}, nil
	}

	if z.DocumentId == nil || strings.TrimSpace(*z.DocumentId) == "" {
		return shared.StudioPlanGetResponse{}, bladWskazaniaStudio(
			"odczyt rozkładu bez wskazania rozkładu ani dokumentu")
	}
	znalezione := p.magazyn.petlaRozkladyDokumentu(strings.TrimSpace(*z.DocumentId))
	if len(znalezione) == 0 {
		return shared.StudioPlanGetResponse{}, bladBrakuStudio(
			"dokument " + *z.DocumentId + " nie ma ani jednego rozkładu zlecenia")
	}
	rozklady := make([]shared.StudioTaskPlan, 0, len(znalezione))
	for _, rozklad := range znalezione {
		rozklady = append(rozklady, petlaZlozRozklad(rozklad, z.State))
	}
	// Rozkład bieżący jest najświeższym — magazyn oddaje wykaz w tej
	// kolejności.
	return shared.StudioPlanGetResponse{Plan: rozklady[0], Plans: rozklady}, nil
}

// ── studio.plan.task.update ─────────────────────────────────────────────────

// PrzestawZadanie obsługuje studio.plan.task.update: przejście stanu jest
// sprawdzane, a nie przyjmowane na słowo, bo zadanie domknięte wracające do
// biegu bez decyzji o ponowieniu zgubiłoby wynik poprzedniej pracy.
func (p *adapterPetliStudia) PrzestawZadanie(_ context.Context,
	z shared.StudioPlanTaskUpdateRequest) (shared.StudioPlanTaskUpdateResponse, error) {

	if strings.TrimSpace(z.TaskId) == "" {
		return shared.StudioPlanTaskUpdateResponse{}, bladWskazaniaStudio(
			"przestawienie bez wskazania zadania")
	}
	zadanie, rozklad, jest := p.magazyn.petlaZadaniePoKodzie(strings.TrimSpace(z.TaskId))
	if !jest {
		return shared.StudioPlanTaskUpdateResponse{},
			bladBrakuStudio("zadanie nie istnieje: " + z.TaskId)
	}
	if z.State == nil && z.Result == nil && z.FailureReason == nil &&
		z.AgentId == nil && z.AgentName == nil && z.SubagentId == nil && z.Instruction == nil {
		return shared.StudioPlanTaskUpdateResponse{}, bladWskazaniaStudio(
			"przestawienie zadania " + z.TaskId + " bez ani jednej zmiany")
	}

	p.magazyn.zamek.Lock()
	if z.State != nil {
		if !petlaWolnoPrzestawic(zadanie.Stan, *z.State) {
			p.magazyn.zamek.Unlock()
			return shared.StudioPlanTaskUpdateResponse{}, bladWskazaniaStudio(
				"zadanie " + zadanie.Kod + " stoi w stanie „" + string(zadanie.Stan) +
					"” i nie przechodzi wprost w „" + string(*z.State) + "”")
		}
		petlaPrzestawStan(zadanie, *z.State)
	}
	if z.Result != nil {
		wynik := *z.Result
		zadanie.Wynik = &wynik
	}
	if z.FailureReason != nil {
		powod := *z.FailureReason
		zadanie.PowodPorazki = &powod
	}
	if z.Instruction != nil {
		zadanie.Polecenie = *z.Instruction
	}
	if z.AgentId != nil || z.AgentName != nil || z.SubagentId != nil {
		zadanie.Wykonawca = petlaWykonawcaZadania(z.AgentId, z.AgentName, z.SubagentId, rozklad.KodDokumentu)
	}
	rozklad.Zaktualizowano = petlaTeraz()
	petlaPrzelicz(rozklad)
	p.magazyn.zamek.Unlock()

	return shared.StudioPlanTaskUpdateResponse{
		Task: petlaZlozZadanie(zadanie),
		Plan: petlaZlozRozklad(rozklad, nil),
	}, nil
}

// petlaWykonawcaZadania składa tożsamość wykonawcy, któremu zadanie
// powierzono, z identyfikatora i nazwy eksperta oraz identyfikatora
// podagenta z żądania.
func petlaWykonawcaZadania(kodEksperta, nazwaEksperta, kodPodagenta *string, _ string) *shared.StudioActor {
	wykonawca := shared.StudioActor{Kind: shared.StudioAuthorModel}
	if kodEksperta != nil && strings.TrimSpace(*kodEksperta) != "" {
		wykonawca.AgentId = kodEksperta
	}
	if nazwaEksperta != nil && strings.TrimSpace(*nazwaEksperta) != "" {
		wykonawca.AgentName = nazwaEksperta
	}
	if kodPodagenta != nil && strings.TrimSpace(*kodPodagenta) != "" {
		wykonawca.SubagentId = kodPodagenta
	}
	return &wykonawca
}

// petlaWolnoPrzestawic rozstrzyga, czy zadanie przechodzi z jednego stanu
// w drugi. Przejścia dozwolone wypisane są wprost, bo zbiór jest mały, a mapa
// czytelniejsza niż łańcuch warunków.
func petlaWolnoPrzestawic(z, na shared.StudioTaskState) bool {
	if z == na {
		return true
	}
	dozwolone := map[shared.StudioTaskState][]shared.StudioTaskState{
		shared.StudioTaskStatePending: {
			shared.StudioTaskStateRunning, shared.StudioTaskStateBlocked,
			shared.StudioTaskStateSkipped, shared.StudioTaskStateFailed,
		},
		shared.StudioTaskStateRunning: {
			shared.StudioTaskStateDone, shared.StudioTaskStateFailed,
			shared.StudioTaskStatePending, shared.StudioTaskStateSkipped,
		},
		shared.StudioTaskStateBlocked: {
			shared.StudioTaskStatePending, shared.StudioTaskStateSkipped,
			shared.StudioTaskStateFailed,
		},
		// Zadanie nieudane wolno ponowić — na tym stoi kontrola jakości pętli.
		shared.StudioTaskStateFailed: {
			shared.StudioTaskStatePending, shared.StudioTaskStateRunning,
			shared.StudioTaskStateSkipped,
		},
		// Zadanie pominięte wolno przywrócić do kolejki; domkniętego nie
		// ruszamy wstecz.
		shared.StudioTaskStateSkipped: {shared.StudioTaskStatePending},
		shared.StudioTaskStateDone:    {},
	}
	for _, wolno := range dozwolone[z] {
		if wolno == na {
			return true
		}
	}
	return false
}

// petlaPrzestawStan przestawia stan zadania wraz z jego znacznikami czasu:
// podjęcia przy przejściu w bieg i domknięcia przy zakończeniu.
func petlaPrzestawStan(zadanie *petlaZadanieStudia, stan shared.StudioTaskState) {
	zadanie.Stan = stan
	teraz := petlaTeraz()
	switch stan {
	case shared.StudioTaskStateRunning:
		zadanie.Podjeto = &teraz
		zadanie.Domknieto = nil
	case shared.StudioTaskStateDone, shared.StudioTaskStateFailed, shared.StudioTaskStateSkipped:
		zadanie.Domknieto = &teraz
	}
}

// petlaPrzelicz przestawia stan rozkładu wedle stanów jego zadań — wołający
// trzyma zamek magazynu przez cały czas przeliczania.
func petlaPrzelicz(rozklad *petlaRozkladStudia) {
	if rozklad.Stan == shared.StudioPlanStateStopped {
		return
	}
	wBiegu, doZrobienia := 0, 0
	for _, zadanie := range rozklad.Zadania {
		switch zadanie.Stan {
		case shared.StudioTaskStateRunning:
			wBiegu++
		case shared.StudioTaskStatePending, shared.StudioTaskStateBlocked:
			doZrobienia++
		}
	}
	switch {
	case wBiegu > 0:
		rozklad.Stan = shared.StudioPlanStateRunning
	case doZrobienia == 0:
		rozklad.Stan = shared.StudioPlanStateDone
	case rozklad.Obiegi > 0:
		rozklad.Stan = shared.StudioPlanStatePaused
	default:
		rozklad.Stan = shared.StudioPlanStateDraft
	}
}

// ── studio.plan.run ─────────────────────────────────────────────────────────

// PuscPetle obsługuje studio.plan.run: sprawdza nastawy, dobiera zadania
// gotowe w kolejnych obiegach i wykonuje je aż do wyczerpania albo
// zatrzymania.
func (p *adapterPetliStudia) PuscPetle(ctx context.Context,
	z shared.StudioPlanRunRequest) (shared.StudioPlanRunResponse, error) {

	if strings.TrimSpace(z.PlanId) == "" {
		return shared.StudioPlanRunResponse{}, bladWskazaniaStudio(
			"uruchomienie pętli bez wskazania rozkładu")
	}
	rozklad, jest := p.magazyn.petlaRozkladPoKodzie(strings.TrimSpace(z.PlanId))
	if !jest {
		return shared.StudioPlanRunResponse{}, bladBrakuStudio("rozkład nie istnieje: " + z.PlanId)
	}

	nastawy, powodBraku := p.petlaNastawy(ctx, rozklad.KodDokumentu)
	if powodBraku != "" {
		// Odmowa nazywa brak nastawy i drogę jej włączenia, polem
		// refusalReason kontraktu.
		powod := powodBraku
		return shared.StudioPlanRunResponse{
			Plan: petlaZlozRozklad(rozklad, nil), Started: false, RefusalReason: &powod,
		}, nil
	}

	dokument, err := p.studio.dokumentDoCzynnosci(ctx, rozklad.KodDokumentu)
	if err != nil {
		return shared.StudioPlanRunResponse{}, err
	}

	granicaObiegow := petlaGranicaObiegow(nastawy, z.MaxIterations)
	granicaBezPostepu := petlaGranicaBezPostepu(nastawy)
	naraz := petlaWykonawcowNaraz(nastawy)
	dopuszczeni := petlaDopuszczeni(z.AgentIds)
	probny := z.DryRun != nil && *z.DryRun

	p.magazyn.zamek.Lock()
	rozklad.zatrzymanie = false
	rozklad.PowodZatrzymania = nil
	rozklad.Stan = shared.StudioPlanStateRunning
	p.magazyn.zamek.Unlock()

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	for obieg := 0; obieg < granicaObiegow; obieg++ {
		if p.petlaCzyZatrzymana(rozklad) {
			break
		}
		gotowe := p.petlaZadaniaGotowe(rozklad, naraz, dopuszczeni)
		if len(gotowe) == 0 {
			p.magazyn.zamek.Lock()
			rozklad.Obiegi++
			rozklad.ObiegiBezPostepu++
			bezPostepu := rozklad.ObiegiBezPostepu
			p.magazyn.zamek.Unlock()
			if bezPostepu >= granicaBezPostepu {
				p.petlaZatrzymaj(rozklad, "pętla stanęła: "+
					"żadne zadanie nie było gotowe do podjęcia w dwóch obiegach z rzędu")
				break
			}
			continue
		}

		postep := 0
		for _, zadanie := range gotowe {
			if p.petlaCzyZatrzymana(rozklad) {
				break
			}
			if probny {
				bilans.SkippedCount++
				bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
					Reason: "przebieg próbny — zadanie „" + zadanie.Nazwa +
						"” było gotowe do podjęcia, dokumentu nie tknięto",
				})
				postep++
				continue
			}
			if p.petlaWykonajZadanie(ctx, rozklad, zadanie, dokument.Okno, &bilans) {
				postep++
			}
		}

		p.magazyn.zamek.Lock()
		rozklad.Obiegi++
		if postep == 0 {
			rozklad.ObiegiBezPostepu++
		} else {
			rozklad.ObiegiBezPostepu = 0
		}
		rozklad.Zaktualizowano = petlaTeraz()
		bezPostepu := rozklad.ObiegiBezPostepu
		p.magazyn.zamek.Unlock()

		if bezPostepu >= granicaBezPostepu {
			p.petlaZatrzymaj(rozklad, "pętla stanęła: dwa obiegi bez postępu")
			break
		}
		if probny {
			// Przebieg próbny sprawdza rozkład raz — drugi obieg oglądałby te
			// same zadania.
			break
		}
	}

	p.magazyn.zamek.Lock()
	petlaPrzelicz(rozklad)
	p.magazyn.zamek.Unlock()

	if probny {
		nota := "przebieg próbny: rozkład i zależności sprawdzone, dokument nietknięty"
		bilans.Note = &nota
	} else if bilans.Note == nil {
		nota := "pętla wykonała " + petlaLiczba(bilans.Applied) + " zadań; pominięto " +
			petlaLiczba(bilans.SkippedCount)
		bilans.Note = &nota
	}
	return shared.StudioPlanRunResponse{
		Plan: petlaZlozRozklad(rozklad, nil), Started: true, Balance: &bilans,
	}, nil
}

// petlaLiczba zamienia liczbę na napis bez sięgania po strconv w treści
// zdania bilansu, budując cyfry od końca liczby.
func petlaLiczba(ile int) string {
	if ile == 0 {
		return "0"
	}
	cyfry := []byte{}
	for ile > 0 {
		cyfry = append([]byte{byte('0' + ile%10)}, cyfry...)
		ile /= 10
	}
	return string(cyfry)
}

// petlaCzyZatrzymana mówi, czy Operator zatrzymał pętlę tego rozkładu,
// czytając znacznik pod zamkiem magazynu.
func (p *adapterPetliStudia) petlaCzyZatrzymana(rozklad *petlaRozkladStudia) bool {
	p.magazyn.zamek.Lock()
	defer p.magazyn.zamek.Unlock()
	return rozklad.zatrzymanie
}

// petlaZatrzymaj przestawia rozkład w stan zatrzymany wraz z powodem i
// znacznikiem czasu ostatniej zmiany.
func (p *adapterPetliStudia) petlaZatrzymaj(rozklad *petlaRozkladStudia, powod string) {
	p.magazyn.zamek.Lock()
	defer p.magazyn.zamek.Unlock()
	rozklad.Stan = shared.StudioPlanStateStopped
	rozklad.PowodZatrzymania = &powod
	rozklad.Zaktualizowano = petlaTeraz()
}

// petlaZadaniaGotowe wybiera zadania, które wolno podjąć w tym obiegu:
// czekające i takie, których wszystkie zależności są domknięte. Zadanie
// czekające na zależność nieudaną przechodzi w stan blocked.
func (p *adapterPetliStudia) petlaZadaniaGotowe(rozklad *petlaRozkladStudia,
	naraz int, dopuszczeni map[string]bool) []*petlaZadanieStudia {

	p.magazyn.zamek.Lock()
	defer p.magazyn.zamek.Unlock()

	stany := make(map[string]shared.StudioTaskState, len(rozklad.Zadania))
	for _, zadanie := range rozklad.Zadania {
		stany[zadanie.Kod] = zadanie.Stan
	}

	gotowe := []*petlaZadanieStudia{}
	for _, zadanie := range rozklad.Zadania {
		if zadanie.Stan != shared.StudioTaskStatePending &&
			zadanie.Stan != shared.StudioTaskStateBlocked {
			continue
		}
		if len(dopuszczeni) > 0 && zadanie.Wykonawca != nil &&
			zadanie.Wykonawca.AgentId != nil && !dopuszczeni[*zadanie.Wykonawca.AgentId] {
			continue
		}
		czeka, wstrzymane := false, false
		for _, kodZaleznosci := range zadanie.StoiNa {
			switch stany[kodZaleznosci] {
			case shared.StudioTaskStateDone, shared.StudioTaskStateSkipped:
			case shared.StudioTaskStateFailed:
				wstrzymane = true
			default:
				czeka = true
			}
		}
		if wstrzymane || czeka {
			if zadanie.Stan == shared.StudioTaskStatePending {
				zadanie.Stan = shared.StudioTaskStateBlocked
			}
			continue
		}
		zadanie.Stan = shared.StudioTaskStatePending
		gotowe = append(gotowe, zadanie)
		if len(gotowe) >= naraz {
			break
		}
	}
	return gotowe
}

// petlaWykonajZadanie wykonuje jedno zadanie i oddaje, czy pętla zrobiła
// postęp: komenda jedzie rejestrem z podpisem wykonawcy w kontekście;
// zadanie treści idzie studio.contextual.op, zadanie wydania —
// studio.document.export.format.
func (p *adapterPetliStudia) petlaWykonajZadanie(ctx context.Context,
	rozklad *petlaRozkladStudia, zadanie *petlaZadanieStudia, okno string,
	bilans *shared.StudioActionBalance) bool {

	p.magazyn.zamek.Lock()
	petlaPrzestawStan(zadanie, shared.StudioTaskStateRunning)
	numer, ile := petlaNumerZadania(rozklad, zadanie)
	p.magazyn.zamek.Unlock()
	p.petlaOglosPostep(rozklad, numer, ile, shared.StudioTaskStateRunning, nil)

	komenda, ladunek, blad := p.petlaZadanieNaKomende(rozklad, zadanie, okno)
	if blad != "" {
		p.petlaDomknijPorazke(rozklad, zadanie, numer, ile, blad, nil, bilans)
		return false
	}

	odpowiedz, powodBraku := p.petlaWolajRejestrem(ctx, rozklad, komenda, ladunek)
	if powodBraku != "" {
		p.petlaDomknijPorazke(rozklad, zadanie, numer, ile, powodBraku, nil, bilans)
		return false
	}
	if odpowiedz.Status != shared.EnvelopeStatusOk {
		powod := "czynność " + string(komenda) + " odmówiła"
		if odpowiedz.Blad != nil {
			powod = protocol.Opis(*odpowiedz.Blad)
		}
		// Odmowa zapory blokad trafia do kolejki jako powód przy zadaniu, nie
		// do dziennika.
		p.petlaDomknijPorazke(rozklad, zadanie, numer, ile, powod, nil, bilans)
		return false
	}

	wynik, kodPropozycji, pominiete := petlaCzytajWynik(odpowiedz.Wynik)

	p.magazyn.zamek.Lock()
	petlaPrzestawStan(zadanie, shared.StudioTaskStateDone)
	zadanie.PowodPorazki = nil
	if wynik != "" {
		tresc := wynik
		zadanie.Wynik = &tresc
	}
	if kodPropozycji != "" {
		zadanie.KodyCzynnosci = append(zadanie.KodyCzynnosci, kodPropozycji)
	}
	// Zadanie wykonane mimo napotkanej blokady niesie o tym zdanie przy
	// sobie.
	if len(pominiete) > 0 {
		nota := petlaZdanieOPominieciach(pominiete)
		zadanie.PowodPorazki = &nota
	}
	rozklad.Zaktualizowano = petlaTeraz()
	p.magazyn.zamek.Unlock()

	bilans.Applied++
	for _, pozycja := range pominiete {
		bilans.SkippedCount++
		bilans.Skipped = append(bilans.Skipped, petlaPominiecieZadania(zadanie, pozycja))
	}
	if len(pominiete) > 0 {
		licznik := bilans.DeferredCount
		if licznik == nil {
			zero := 0
			licznik = &zero
		}
		*licznik += len(pominiete)
		bilans.DeferredCount = licznik
	}

	var wskazanieKodu *string
	if kodPropozycji != "" {
		wskazanieKodu = &kodPropozycji
	}
	p.petlaOglosPostep(rozklad, numer, ile, shared.StudioTaskStateDone, wskazanieKodu)
	return true
}

// petlaZadanieNaKomende przekłada zadanie na komendę kontraktu wraz
// z ładunkiem. Trzeci wynik jest powodem, dla którego przekład się nie udał.
func (p *adapterPetliStudia) petlaZadanieNaKomende(rozklad *petlaRozkladStudia,
	zadanie *petlaZadanieStudia, okno string) (shared.MessageType, json.RawMessage, string) {

	if zadanie.Rodzaj == shared.StudioTaskKindExport {
		// Format wydania bierze się ze słów zlecenia; zlecenie bez nazwanego
		// formatu dostaje pdf.
		ladunek, err := json.Marshal(shared.StudioDocumentExportFormatRequest{
			DocumentId: rozklad.KodDokumentu,
			Format:     petlaFormatWydania(rozklad.Zlecenie + " " + zadanie.Polecenie),
		})
		if err != nil {
			return "", nil, "nie udało się złożyć żądania wydania: " + err.Error()
		}
		return shared.CommandStudioDocumentExportFormat, ladunek, ""
	}

	if okno == "" {
		return "", nil, "dokument nie ma wskazanego okna, więc pętla nie ma kanału " +
			"modelu, którym mogłaby wykonać zadanie"
	}
	// Wyliczenia kontraktu są stałymi nietypowanymi, więc typ podajemy
	// wprost.
	var zakres shared.StudioOperationScope = shared.StudioOperationScopeDocument
	if zadanie.ZakresOd != nil && zadanie.ZakresDo != nil {
		zakres = shared.StudioOperationScopeSelection
	}
	parametry, err := json.Marshal(map[string]string{
		"zadanie":   zadanie.Nazwa,
		"rodzaj":    string(zadanie.Rodzaj),
		"zlecenie":  rozklad.Zlecenie,
		"polecenie": zadanie.Polecenie,
	})
	if err != nil {
		parametry = nil
	}
	ladunek, err := json.Marshal(shared.StudioContextualOpRequest{
		WindowId:       okno,
		DocumentId:     rozklad.KodDokumentu,
		ActionId:       zadanie.Nazwa,
		Scope:          zakres,
		SelectionStart: zadanie.ZakresOd,
		SelectionEnd:   zadanie.ZakresDo,
		Params:         parametry,
	})
	if err != nil {
		return "", nil, "nie udało się złożyć żądania operacji: " + err.Error()
	}
	return shared.CommandStudioContextualOp, petlaPodpiszLadunek(ladunek, rozklad), ""
}

// petlaPodpiszLadunek dokłada do ładunku podpis wykonawcy, w którego imieniu
// pętla woła komendę: zapora blokad fragmentu rozpoznaje wykonawcę z podpisu
// w ładunku, bo gniazdo mówi „Operator”, a pracuje model.
func petlaPodpiszLadunek(ladunek json.RawMessage, rozklad *petlaRozkladStudia) json.RawMessage {
	pola := map[string]json.RawMessage{}
	if err := json.Unmarshal(ladunek, &pola); err != nil {
		return ladunek
	}
	dopisz := func(nazwa string, wartosc string) {
		if strings.TrimSpace(wartosc) == "" {
			return
		}
		if zapis, err := json.Marshal(wartosc); err == nil {
			pola[nazwa] = zapis
		}
	}
	dopisz("author", string(shared.StudioAuthorModel))
	if rozklad.Ulozyl != nil {
		dopisz("agentId", wartoscTekstu(rozklad.Ulozyl.AgentId))
		dopisz("agentName", wartoscTekstu(rozklad.Ulozyl.AgentName))
		dopisz("subagentId", wartoscTekstu(rozklad.Ulozyl.SubagentId))
	}
	podpisany, err := json.Marshal(pola)
	if err != nil {
		return ladunek
	}
	return podpisany
}

// petlaFormatWydania rozpoznaje format wydania ze słów zlecenia, dobierając
// format najbliższy początkowi treści, a przy braku wskazania — pdf.
func petlaFormatWydania(tresc string) shared.StudioExportFormat {
	nazwy := map[string]shared.StudioExportFormat{
		"docx": shared.StudioExportFormatDocx,
		"word": shared.StudioExportFormatDocx,
		"odt":  shared.StudioExportFormatOdt,
		"html": shared.StudioExportFormatHtml,
		"rtf":  shared.StudioExportFormatRtf,
		"txt":  shared.StudioExportFormatTxt,
		"pdf":  shared.StudioExportFormatPdf,
	}
	male := strings.ToLower(tresc)
	najblizsze, wybrany := -1, shared.StudioExportFormat(shared.StudioExportFormatPdf)
	for slowo, format := range nazwy {
		miejsce := strings.Index(male, slowo)
		if miejsce >= 0 && (najblizsze < 0 || miejsce < najblizsze) {
			najblizsze, wybrany = miejsce, format
		}
	}
	// Markdown osobno, bo skrót „md" trafiłby w środek dowolnego wyrazu.
	if strings.Contains(male, "markdown") {
		return shared.StudioExportFormatMd
	}
	return wybrany
}

// petlaWolajRejestrem woła komendę kontraktu przez rejestr — tą samą drogą,
// którą jedzie klient. Drugi wynik jest powodem, gdy komendy nie ma czym zawołać.
func (p *adapterPetliStudia) petlaWolajRejestrem(ctx context.Context,
	rozklad *petlaRozkladStudia, komenda shared.MessageType,
	ladunek json.RawMessage) (protocol.Odpowiedz, string) {

	if p.rejestr == nil {
		return protocol.Odpowiedz{}, "rdzeń złożony bez rejestru komend — pętla nie ma " +
			"czym wykonać ani jednego zadania"
	}
	obsluga, jest := p.rejestr.Obsluga(komenda)
	if !jest {
		return protocol.Odpowiedz{}, "rdzeń nie ma obsługiwacza komendy " + string(komenda) +
			", więc tego zadania nie ma czym wykonać; pętla nie udaje wykonania"
	}
	return obsluga(p.petlaKontekstWykonawcy(ctx, rozklad), protocol.Request{
		Komenda: komenda, TypZadany: komenda, Znana: true, Ladunek: ladunek,
	}), ""
}

// petlaKontekstWykonawcy wkłada w kontekst tożsamość wykonawcy, w imieniu
// którego pętla pracuje. Bez tego zapora blokad wzięłaby pracę pętli za pracę
// Operatora i przepuściła ją przez fragment zablokowany przed wykonawcą.
func (p *adapterPetliStudia) petlaKontekstWykonawcy(ctx context.Context,
	rozklad *petlaRozkladStudia) context.Context {

	wykonawca := kontrolaWykonawca{Rodzaj: shared.StudioAuthorModel}
	if rozklad.Ulozyl != nil {
		wykonawca.AgentKod = rozklad.Ulozyl.AgentId
		wykonawca.AgentNazwa = rozklad.Ulozyl.AgentName
		wykonawca.AgentWersja = rozklad.Ulozyl.AgentVersion
		wykonawca.PodagentKod = rozklad.Ulozyl.SubagentId
		wykonawca.OknoID = rozklad.Ulozyl.WindowId
		wykonawca.KartaSesjiID = rozklad.Ulozyl.SessionId
		wykonawca.KanalID = rozklad.Ulozyl.ChannelId
	}
	return kontrolaZapiszWykonawce(ctx, wykonawca)
}

// petlaCzytajWynik wyjmuje z odpowiedzi to, czego pętla potrzebuje: treść
// wyniku, kod propozycji i wykaz pominięć, czytane po polach, bo pętla
// wykonuje zadania kilkoma różnymi komendami.
func petlaCzytajWynik(ladunek json.RawMessage) (string, string, []shared.StudioSkippedItem) {
	if len(ladunek) == 0 {
		return "", "", nil
	}
	var pola struct {
		ResultText *string                     `json:"resultText,omitempty"`
		ProposalId *string                     `json:"proposalId,omitempty"`
		Balance    *shared.StudioActionBalance `json:"balance,omitempty"`
		Result     *shared.StudioExportResult  `json:"result,omitempty"`
	}
	if err := json.Unmarshal(ladunek, &pola); err != nil {
		return "", "", nil
	}
	wynik, kod := "", ""
	if pola.ResultText != nil {
		wynik = *pola.ResultText
	}
	if pola.ProposalId != nil {
		kod = *pola.ProposalId
	}
	pominiete := []shared.StudioSkippedItem{}
	if pola.Balance != nil {
		pominiete = append(pominiete, pola.Balance.Skipped...)
	}
	// Wydanie do formatu uboższego niż dokument oddaje własny wykaz cech
	// pominiętych.
	if pola.Result != nil {
		if wynik == "" && pola.Result.Path != nil {
			wynik = *pola.Result.Path
		}
		pominiete = append(pominiete, pola.Result.DroppedFeatures...)
	}
	return wynik, kod, pominiete
}

// petlaPominiecieZadania przenosi pominięcie do bilansu przebiegu, nazywając
// zadanie, przy którym zaszło. Wskazanie blokady zostaje nietknięte — po nim
// okno dochodzi do blokady, która czynność zatrzymała.
func petlaPominiecieZadania(zadanie *petlaZadanieStudia,
	pozycja shared.StudioSkippedItem) shared.StudioSkippedItem {

	przeniesione := pozycja
	przeniesione.Reason = "zadanie „" + zadanie.Nazwa + "”: " + pozycja.Reason
	if przeniesione.Detail == nil {
		nazwa := zadanie.Nazwa
		przeniesione.Detail = &nazwa
	}
	return przeniesione
}

// petlaZdanieOPominieciach składa zdanie o pominięciach dla wiersza kolejki,
// dołączając nazwę blokady do każdego powodu, gdy ta jest znana.
func petlaZdanieOPominieciach(pominiete []shared.StudioSkippedItem) string {
	czesci := make([]string, 0, len(pominiete))
	for _, pozycja := range pominiete {
		zdanie := pozycja.Reason
		if pozycja.LockName != nil && *pozycja.LockName != "" {
			zdanie += " (blokada „" + *pozycja.LockName + "”)"
		}
		czesci = append(czesci, zdanie)
	}
	return "zadanie wykonane z pominięciami: " + strings.Join(czesci, "; ")
}

// petlaDomknijPorazke zamyka zadanie stanem nieudanym i wpisuje powód do
// bilansu wraz ze wskazaniem blokady, gdy ta zaszła.
func (p *adapterPetliStudia) petlaDomknijPorazke(rozklad *petlaRozkladStudia,
	zadanie *petlaZadanieStudia, numer, ile int, powod string,
	blokada *shared.StudioSkippedItem, bilans *shared.StudioActionBalance) {

	p.magazyn.zamek.Lock()
	petlaPrzestawStan(zadanie, shared.StudioTaskStateFailed)
	zapisany := powod
	zadanie.PowodPorazki = &zapisany
	rozklad.Zaktualizowano = petlaTeraz()
	p.magazyn.zamek.Unlock()

	pozycja := shared.StudioSkippedItem{Reason: "zadanie „" + zadanie.Nazwa + "”: " + powod}
	if blokada != nil {
		pozycja.LockId, pozycja.LockName = blokada.LockId, blokada.LockName
		pozycja.RangeStart, pozycja.RangeEnd = blokada.RangeStart, blokada.RangeEnd
	}
	bilans.SkippedCount++
	bilans.Skipped = append(bilans.Skipped, pozycja)
	p.petlaOglosPostep(rozklad, numer, ile, shared.StudioTaskStateFailed, nil)
}

// petlaNumerZadania oddaje numer zadania w rozkładzie (od 1) i liczbę zadań.
// Wołający trzyma zamek magazynu.
func petlaNumerZadania(rozklad *petlaRozkladStudia, zadanie *petlaZadanieStudia) (int, int) {
	for numer, wpis := range rozklad.Zadania {
		if wpis == zadanie {
			return numer + 1, len(rozklad.Zadania)
		}
	}
	return 0, len(rozklad.Zadania)
}

// petlaOglosPostep rozgłasza `studio.chain.progressed` — zdarzenie, którym
// kontrakt zasila pętlę wykonawczą okna. Okno odświeża się nim, zamiast pytać
// rdzeń w kółko o stan każdego zadania.
func (p *adapterPetliStudia) petlaOglosPostep(rozklad *petlaRozkladStudia,
	numer, ile int, stan shared.StudioTaskState, kodPropozycji *string) {

	if p.emiter == nil || numer == 0 {
		return
	}
	p.emiter.wyslij(shared.EventStudioChainProgressed, "", shared.StudioChainProgressedEvent{
		RunId: rozklad.Kod, StepIndex: numer, StepCount: ile,
		State: string(stan), ProposalId: kodPropozycji,
	})
}

// ── studio.plan.stop ────────────────────────────────────────────────────────

// ZatrzymajPetle obsługuje studio.plan.stop: podnosi znacznik, który pętla
// sprawdza przed każdym zadaniem i po każdym obiegu, a zadania w biegu
// odstawia z powrotem do czekania z powodem przerwania.
func (p *adapterPetliStudia) ZatrzymajPetle(_ context.Context,
	z shared.StudioPlanStopRequest) (shared.StudioPlanStopResponse, error) {

	if strings.TrimSpace(z.PlanId) == "" {
		return shared.StudioPlanStopResponse{}, bladWskazaniaStudio(
			"zatrzymanie bez wskazania rozkładu")
	}
	rozklad, jest := p.magazyn.petlaRozkladPoKodzie(strings.TrimSpace(z.PlanId))
	if !jest {
		return shared.StudioPlanStopResponse{}, bladBrakuStudio("rozkład nie istnieje: " + z.PlanId)
	}

	powod := "zatrzymanie przez Operatora"
	if z.Reason != nil && strings.TrimSpace(*z.Reason) != "" {
		powod = strings.TrimSpace(*z.Reason)
	}

	p.magazyn.zamek.Lock()
	rozklad.zatrzymanie = true
	rozklad.Stan = shared.StudioPlanStateStopped
	rozklad.PowodZatrzymania = &powod
	rozklad.Zaktualizowano = petlaTeraz()
	odstawione := 0
	for _, zadanie := range rozklad.Zadania {
		if zadanie.Stan != shared.StudioTaskStateRunning {
			continue
		}
		zadanie.Stan = shared.StudioTaskStatePending
		zadanie.Domknieto = nil
		przerwane := "zadanie PRZERWANE zatrzymaniem pętli (" + powod +
			"); wróciło do czekania, praca wykonana przed przerwaniem została w dokumencie"
		zadanie.PowodPorazki = &przerwane
		odstawione++
	}
	p.magazyn.zamek.Unlock()

	return shared.StudioPlanStopResponse{
		Plan: petlaZlozRozklad(rozklad, nil), StoppedTasks: odstawione,
	}, nil
}

// ── Nastawy pętli ───────────────────────────────────────────────────────────

// petlaNastawy czyta nastawy pętli obsługiwaczem studio.agents.settings.get
// i oddaje je albo — gdy pętli włączyć nie sposób — powód odmowy nazywający
// brak, jako napis w polu refusalReason.
func (p *adapterPetliStudia) petlaNastawy(ctx context.Context,
	kodDokumentu string) (shared.StudioAgentSettings, string) {

	if p.rejestr == nil {
		return shared.StudioAgentSettings{}, "pętla wykonawcza jest wyłączona: " +
			"rdzeń złożony bez rejestru komend, więc nastawy nie ma czym odczytać"
	}
	obsluga, jest := p.rejestr.Obsluga(shared.CommandStudioAgentsSettingsGet)
	if !jest {
		return shared.StudioAgentSettings{}, "pętla wykonawcza jest wyłączona: " +
			"rdzeń nie ma jeszcze obsługiwacza komendy studio.agents.settings.get, " +
			"więc nastawy „executionLoopEnabled” nie da się ani odczytać, ani włączyć"
	}
	ladunek, err := json.Marshal(shared.StudioAgentsSettingsGetRequest{DocumentId: &kodDokumentu})
	if err != nil {
		return shared.StudioAgentSettings{}, "pętla wykonawcza jest wyłączona: " +
			"nie udało się złożyć zapytania o nastawy: " + err.Error()
	}
	odpowiedz := obsluga(ctx, protocol.Request{
		Komenda:   shared.CommandStudioAgentsSettingsGet,
		TypZadany: shared.CommandStudioAgentsSettingsGet,
		Znana:     true,
		Ladunek:   ladunek,
	})
	if odpowiedz.Status != shared.EnvelopeStatusOk {
		powod := "nastawy pętli nie dały się odczytać"
		if odpowiedz.Blad != nil {
			powod = protocol.Opis(*odpowiedz.Blad)
		}
		return shared.StudioAgentSettings{}, "pętla wykonawcza jest wyłączona: " + powod
	}
	var wynik shared.StudioAgentsSettingsGetResponse
	if err := json.Unmarshal(odpowiedz.Wynik, &wynik); err != nil {
		return shared.StudioAgentSettings{}, "pętla wykonawcza jest wyłączona: " +
			"odpowiedź o nastawach nie ma kształtu kontraktu: " + err.Error()
	}
	if !wynik.Settings.ExecutionLoopEnabled {
		return wynik.Settings, "pętla wykonawcza jest WYŁĄCZONA nastawą Operatora " +
			"(„executionLoopEnabled” = nie). Pętla jest narzędziem nieuruchamianym " +
			"na starcie; włącza ją komenda studio.agents.settings.set na wskazanym " +
			"zasięgu. Dopóki nastawa stoi na „nie”, pętla nie tknie dokumentu"
	}
	return wynik.Settings, ""
}

// petlaGranicaObiegow oddaje górną granicę obiegów tego uruchomienia:
// z żądania, z nastawy Operatora albo z wartości domyślnej.
func petlaGranicaObiegow(nastawy shared.StudioAgentSettings, zZadania *int) int {
	if zZadania != nil && *zZadania > 0 {
		return *zZadania
	}
	if nastawy.LoopMaxIterations != nil && *nastawy.LoopMaxIterations > 0 {
		return *nastawy.LoopMaxIterations
	}
	return petlaObiegiDomyslnie
}

// petlaGranicaBezPostepu oddaje liczbę obiegów bez postępu, po której pętla
// staje, z nastawy Operatora albo z wartości domyślnej.
func petlaGranicaBezPostepu(nastawy shared.StudioAgentSettings) int {
	if nastawy.LoopNoProgressThreshold != nil && *nastawy.LoopNoProgressThreshold > 0 {
		return *nastawy.LoopNoProgressThreshold
	}
	return petlaBezPostepuDomyslnie
}

// petlaWykonawcowNaraz oddaje, ilu wykonawców pętla podejmuje w jednym obiegu.
//
// Praca kilku wykonawców naraz jest osobnym narzędziem i osobną nastawą.
// Wyłączona znaczy JEDEN wykonawca na obieg — nie „tyle, ile się zmieści".
func petlaWykonawcowNaraz(nastawy shared.StudioAgentSettings) int {
	if !nastawy.MultiAgentEnabled {
		return petlaWykonawcowDomyslnie
	}
	if nastawy.MaxConcurrentAgents != nil && *nastawy.MaxConcurrentAgents > 0 {
		return *nastawy.MaxConcurrentAgents
	}
	return petlaWykonawcowDomyslnie
}

// petlaDopuszczeni składa zbiór wykonawców dopuszczonych do rozkładu
// z wykazu identyfikatorów żądania, pomijając wpisy puste.
func petlaDopuszczeni(kody []string) map[string]bool {
	if len(kody) == 0 {
		return nil
	}
	zbior := make(map[string]bool, len(kody))
	for _, kod := range kody {
		if przyciety := strings.TrimSpace(kod); przyciety != "" {
			zbior[przyciety] = true
		}
	}
	return zbior
}

// ── Rejestracja ─────────────────────────────────────────────────────────────

// zarejestrujPetleWykonawczaStudia wpina pięć komend rodziny studio.plan.*.
// Port Studia wchodzi rzutowaniem dwuwartościowym na typ adaptera modułu,
// bo pętla stoi na dokumencie i na silniku modelu tego adaptera.
func zarejestrujPetleWykonawczaStudia(r *Rejestr, m Studio, e *emiter) {
	if r == nil || m == nil {
		return
	}
	studio, jest := m.(*adapterStudia)
	if !jest {
		return
	}
	petla := nowyAdapterPetliStudia(studio, r, e)

	r.Zarejestruj(shared.CommandStudioPlanCreate, obsluz(petla.RozlozZlecenie))
	r.Zarejestruj(shared.CommandStudioPlanGet, obsluz(petla.Rozklad))
	r.Zarejestruj(shared.CommandStudioPlanTaskUpdate, obsluz(petla.PrzestawZadanie))
	r.Zarejestruj(shared.CommandStudioPlanStop, obsluz(petla.ZatrzymajPetle))

	// Uruchomienie pętli zmienia dokument; czytamy go ponownie, bo
	// odpowiedź komendy go nie niesie.
	r.Zarejestruj(shared.CommandStudioPlanRun,
		obsluz(func(ctx context.Context, z shared.StudioPlanRunRequest) (shared.StudioPlanRunResponse, error) {
			odpowiedz, err := petla.PuscPetle(ctx, z)
			if err != nil || !odpowiedz.Started {
				return odpowiedz, err
			}
			if odpowiedz.Balance == nil || odpowiedz.Balance.Applied == 0 {
				return odpowiedz, nil
			}
			kod := odpowiedz.Plan.DocumentId
			po, bladOdczytu := studio.OtworzDokument(ctx, shared.StudioDocumentOpenRequest{
				DocumentId: &kod,
			})
			if bladOdczytu == nil {
				e.dokumentStudio(shared.ChangeKindUpdated, po.Document)
			}
			return odpowiedz, nil
		}))
}
