// Plik obsługuje cztery komendy dopełniające układ zależności: ustawienie bramki, grupy, kompensacji oraz spięcie z kolejką MultitaskingAI, prowadzone na adapterze modułu Automations.
package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekGrupyUkladu znakuje kod grupy kroków nadany przez rdzeń, aby odróżnić go od kodów bramek i kompensacji zapisanych w tej samej tabeli układu.
const przedrostekGrupyUkladu = "og-"

// ZUkladem wpina repozytorium dopełnień układu. Bez niego cztery komendy
// odmawiają wprost — odmowa jest odpowiedzią, a milczący zapis nie.
func (a *adapterAutomatyk) ZUkladem(uklad dane.RepozytoriumUkladuOrkiestracji) *adapterAutomatyk {
	a.uklad = uklad
	return a
}

// ── orchestration.gate.set ───────────────────────────────────────────────────

// UstawBramke ustala regułę scalenia torów równoległych na wskazanym kroku i oddaje bramki układu po zapisie wraz z jego oceną. Ocena nie blokuje zapisu.
func (a *adapterAutomatyk) UstawBramke(ctx context.Context,
	z shared.OrchestrationGateSetRequest) (shared.OrchestrationGateSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.OrchestrationGateSetResponse{}, err
	}
	if a.uklad == nil {
		return shared.OrchestrationGateSetResponse{}, bladBrakuKatalogu("układu orkiestracji")
	}
	krok := strings.TrimSpace(z.StepId)
	if krok == "" {
		return shared.OrchestrationGateSetResponse{}, bladWskazaniaAutomatyki("bramka bez wskazania kroku")
	}
	if !regulaBramkiZnana(z.Rule) {
		return shared.OrchestrationGateSetResponse{},
			bladWskazaniaAutomatyki("reguła scalenia " + string(z.Rule) + " nie należy do kontraktu")
	}
	licznik := 0
	if z.Count != nil {
		licznik = *z.Count
	}
	// Reguła licznikowa bez liczby torów nie wskazuje, kiedy się otwiera — to błąd żądania.
	if z.Rule == shared.OrchestrationGateRuleCount && licznik < 1 {
		return shared.OrchestrationGateSetResponse{},
			bladWskazaniaAutomatyki("reguła licznikowa wymaga liczby torów większej od zera")
	}
	if z.Rule != shared.OrchestrationGateRuleCount {
		licznik = 0
	}
	if err := a.uklad.ZapiszBramkeUkladu(ctx, wiersz.ID,
		dane.BramkaUkladu{Krok: krok, Regula: string(z.Rule), Licznik: licznik}); err != nil {
		return shared.OrchestrationGateSetResponse{}, bladAutomatyki(err)
	}
	bramki, err := a.uklad.BramkiUkladu(ctx, wiersz.ID)
	if err != nil {
		return shared.OrchestrationGateSetResponse{}, bladAutomatyki(err)
	}
	zastrzezenia, err := a.zastrzezeniaBramek(ctx, wiersz.ID, bramki)
	if err != nil {
		return shared.OrchestrationGateSetResponse{}, err
	}
	return shared.OrchestrationGateSetResponse{
		Gates:  bramkiKontraktu(z.WorkflowId, bramki),
		Valid:  len(zastrzezenia) == 0,
		Issues: zastrzezenia,
	}, nil
}

// zastrzezeniaBramek nazywa bramki, których krok nie istnieje w automatyce albo
// które nie mają czego scalać.
func (a *adapterAutomatyk) zastrzezeniaBramek(ctx context.Context, automatykaID int64,
	bramki []dane.BramkaUkladu) ([]string, error) {

	kroki, zaleznosci, err := a.ukladPoZapisie(ctx, automatykaID)
	if err != nil {
		return nil, err
	}
	znane := map[string]struct{}{}
	for _, krok := range kroki {
		znane[krok.Kod] = struct{}{}
	}
	dochodzace := map[string]int{}
	for _, luk := range zaleznosci {
		dochodzace[luk.KrokDo]++
	}
	zastrzezenia := []string{}
	for _, bramka := range bramki {
		if _, jest := znane[bramka.Krok]; !jest {
			zastrzezenia = append(zastrzezenia,
				"bramka wskazuje krok "+bramka.Krok+", którego w układzie nie ma")
			continue
		}
		if bramka.Regula == string(shared.OrchestrationGateRuleCount) &&
			bramka.Licznik > dochodzace[bramka.Krok] {
			zastrzezenia = append(zastrzezenia,
				"bramka kroku "+bramka.Krok+" wymaga "+strconv.Itoa(bramka.Licznik)+
					" torów, a dochodzi ich "+strconv.Itoa(dochodzace[bramka.Krok]))
		}
	}
	return zastrzezenia, nil
}

// regulaBramkiZnana sprawdza wskazanie reguły scalenia wobec trzech reguł dopuszczonych kontraktem, odrzucając wartość spoza tego zbioru.
func regulaBramkiZnana(regula shared.OrchestrationGateRule) bool {
	for _, znana := range shared.WartosciOrchestrationGateRule() {
		if znana == regula {
			return true
		}
	}
	return false
}

// bramkiKontraktu przekłada wiersze bramek odczytane z repozytorium układu na kształt bramek zgodny z kontraktem odpowiedzi.
func bramkiKontraktu(idUkladu string, wiersze []dane.BramkaUkladu) []shared.OrchestrationGate {
	bramki := make([]shared.OrchestrationGate, 0, len(wiersze))
	for _, wiersz := range wiersze {
		bramka := shared.OrchestrationGate{
			WorkflowId: idUkladu,
			StepId:     wiersz.Krok,
			Rule:       shared.OrchestrationGateRule(wiersz.Regula),
		}
		if wiersz.Licznik > 0 {
			licznik := wiersz.Licznik
			bramka.Count = &licznik
		}
		bramki = append(bramki, bramka)
	}
	return bramki
}

// ── orchestration.group.set ──────────────────────────────────────────────────

// UstawGrupe oznacza zbiór kroków jako wykonywanych równolegle albo w ścisłej
// kolejności. Wykaz kroków pusty USUWA grupę — tak stanowi kontrakt.
func (a *adapterAutomatyk) UstawGrupe(ctx context.Context,
	z shared.OrchestrationGroupSetRequest) (shared.OrchestrationGroupSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.OrchestrationGroupSetResponse{}, err
	}
	if a.uklad == nil {
		return shared.OrchestrationGroupSetResponse{}, bladBrakuKatalogu("układu orkiestracji")
	}
	if !rodzajGrupowaniaZnany(z.Kind) {
		return shared.OrchestrationGroupSetResponse{},
			bladWskazaniaAutomatyki("rodzaj grupowania " + string(z.Kind) + " nie należy do kontraktu")
	}
	kod := strings.TrimSpace(wartoscTekstu(z.GroupId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekGrupyUkladu)
	}
	kroki := make([]string, 0, len(z.StepIds))
	for _, krok := range z.StepIds {
		if przyciety := strings.TrimSpace(krok); przyciety != "" {
			kroki = append(kroki, przyciety)
		}
	}
	if err := a.uklad.ZapiszGrupeUkladu(ctx, wiersz.ID, dane.GrupaUkladu{
		Kod: kod, Nazwa: strings.TrimSpace(z.Name), Rodzaj: string(z.Kind), Kroki: kroki,
	}); err != nil {
		return shared.OrchestrationGroupSetResponse{}, bladAutomatyki(err)
	}
	grupy, err := a.uklad.GrupyUkladu(ctx, wiersz.ID)
	if err != nil {
		return shared.OrchestrationGroupSetResponse{}, bladAutomatyki(err)
	}
	return shared.OrchestrationGroupSetResponse{Groups: grupyKontraktu(z.WorkflowId, grupy)}, nil
}

// rodzajGrupowaniaZnany sprawdza wskazanie rodzaju grupowania wobec trzech rodzajów dopuszczonych kontraktem, odrzucając wartość spoza tego zbioru.
func rodzajGrupowaniaZnany(rodzaj shared.AutomationDependencyKind) bool {
	for _, znany := range shared.WartosciAutomationDependencyKind() {
		if znany == rodzaj {
			return true
		}
	}
	return false
}

// grupyKontraktu przekłada wiersze grup odczytane z repozytorium układu na kształt grup zgodny z kontraktem odpowiedzi.
func grupyKontraktu(idUkladu string, wiersze []dane.GrupaUkladu) []shared.OrchestrationGroup {
	grupy := make([]shared.OrchestrationGroup, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kroki := wiersz.Kroki
		if kroki == nil {
			kroki = []string{}
		}
		grupy = append(grupy, shared.OrchestrationGroup{
			Id:         wiersz.Kod,
			WorkflowId: idUkladu,
			Name:       wiersz.Nazwa,
			StepIds:    kroki,
			Kind:       shared.AutomationDependencyKind(wiersz.Rodzaj),
		})
	}
	return grupy
}

// ── orchestration.compensation.set ───────────────────────────────────────────

// UstawKompensacje ustala krok wycofujący skutki kroku głównego przy błędzie
// w połowie przebiegu. Krok wycofujący pusty zdejmuje kompensację.
func (a *adapterAutomatyk) UstawKompensacje(ctx context.Context,
	z shared.OrchestrationCompensationSetRequest) (shared.OrchestrationCompensationSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.OrchestrationCompensationSetResponse{}, err
	}
	if a.uklad == nil {
		return shared.OrchestrationCompensationSetResponse{}, bladBrakuKatalogu("układu orkiestracji")
	}
	krok := strings.TrimSpace(z.StepId)
	if krok == "" {
		return shared.OrchestrationCompensationSetResponse{},
			bladWskazaniaAutomatyki("kompensacja bez wskazania kroku głównego")
	}
	wycofujacy := strings.TrimSpace(z.CompensationStepId)
	if wycofujacy == krok {
		return shared.OrchestrationCompensationSetResponse{},
			bladWskazaniaAutomatyki("krok " + krok + " nie wycofuje sam siebie")
	}
	if err := a.uklad.ZapiszKompensacjeUkladu(ctx, wiersz.ID,
		dane.KompensacjaUkladu{Krok: krok, KrokWycofu: wycofujacy}); err != nil {
		return shared.OrchestrationCompensationSetResponse{}, bladAutomatyki(err)
	}
	kompensacje, err := a.uklad.KompensacjeUkladu(ctx, wiersz.ID)
	if err != nil {
		return shared.OrchestrationCompensationSetResponse{}, bladAutomatyki(err)
	}
	wynik := make([]shared.OrchestrationCompensation, 0, len(kompensacje))
	for _, pozycja := range kompensacje {
		wynik = append(wynik, shared.OrchestrationCompensation{
			WorkflowId: z.WorkflowId, StepId: pozycja.Krok, CompensationStepId: pozycja.KrokWycofu,
		})
	}
	return shared.OrchestrationCompensationSetResponse{Compensations: wynik}, nil
}

// ── orchestration.multitasking.link ──────────────────────────────────────────

// SpnijZMultitaskingiem spina silnik kolejek Automations z silnikiem kolejek
// środowiska MultitaskingAI. Stan wyjściowy: rozłączone.
func (a *adapterAutomatyk) SpnijZMultitaskingiem(ctx context.Context,
	z shared.OrchestrationMultitaskingLinkRequest) (shared.OrchestrationMultitaskingLinkResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.OrchestrationMultitaskingLinkResponse{}, err
	}
	if a.uklad == nil {
		return shared.OrchestrationMultitaskingLinkResponse{}, bladBrakuKatalogu("układu orkiestracji")
	}
	oknoRoli, err := a.oknoRoliSpiecia(ctx, z.RoleId)
	if err != nil {
		return shared.OrchestrationMultitaskingLinkResponse{}, err
	}
	kolejki, err := a.uklad.ZapiszSpiecieAutomatyki(ctx, wiersz.ID, wiersz.Kod, z.Linked, oknoRoli)
	if err != nil {
		return shared.OrchestrationMultitaskingLinkResponse{}, bladAutomatyki(err)
	}
	identyfikatory := make([]string, 0, len(kolejki))
	for _, id := range kolejki {
		identyfikatory = append(identyfikatory, strconv.FormatInt(id, 10))
	}
	return shared.OrchestrationMultitaskingLinkResponse{
		Linked: z.Linked, QueueIds: identyfikatory,
	}, nil
}

// oknoRoliSpiecia przekłada wskazanie roli na klucz wiersza okna. Rola nierozpoznana wraca odmową, ponieważ nie ma czego spiąć.
func (a *adapterAutomatyk) oknoRoliSpiecia(ctx context.Context, wskazanie *string) (*int64, error) {
	kod := strings.TrimSpace(wartoscTekstu(wskazanie))
	if kod == "" {
		return nil, nil
	}
	if a.okna == nil {
		return nil, bladBrakuKatalogu("okien komunikacji")
	}
	okno, err := a.okna.PoIdentyfikatorze(ctx, kod)
	if err != nil {
		return nil, bladWskazaniaAutomatyki("rola " + kod + " nie wskazuje żadnego okna środowiska")
	}
	numer := okno.ID
	return &numer, nil
}

// ZOknami wpina rejestr okien komunikacji. Bez niego spięcie ze wskazaniem roli
// odmawia — rola bez okna nie ma czego spiąć.
func (a *adapterAutomatyk) ZOknami(okna dane.RepozytoriumOkien) *adapterAutomatyk {
	a.okna = okna
	return a
}
