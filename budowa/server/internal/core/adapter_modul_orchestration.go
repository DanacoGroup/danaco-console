// Rodzina orchestration.* — cztery komendy ukladu zaleznosci, wniesione do
// kontraktu osobno od okna Orchestrator; jedzie na tej samej maszynerii, ale
// pracuje pojedynczym lukiem zamiast kompletem.
package core

import (
	"context"
	"sort"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Asercja kompilatora: adapter modułu Automations wypełnia rozszerzony port.
// Bez niej rozjazd nazwy albo podpisu metody wyszedłby dopiero przy montażu,
// jako panika asercji typu na uruchomionym rdzeniu, zamiast przy budowaniu.
var _ Orkiestracja = (*adapterAutomatyk)(nil)

// repozytoriumZaleznosciKrokow to rozszerzenie repozytorium automatyk o zapis
// i usuniecie pojedynczego luku ukladu; deklaracja stoi po stronie czytelnika,
// bo wymaga jej wylacznie ta rodzina komend.
type repozytoriumZaleznosciKrokow interface {
	ZapiszZaleznosc(ctx context.Context, automatykaID int64, zaleznosc dane.ZaleznoscKroku) error
	UsunZaleznosc(ctx context.Context, automatykaID int64, krokZ, krokDo string) (bool, error)
}

// ── orchestration.dependency.set ─────────────────────────────────────────────

// UstawZaleznosc zapisuje jeden łuk układu zależności i oddaje układ po zapisie
// wraz z jego oceną. Łuk zastany zmienia, nowego dokłada.
func (a *adapterAutomatyk) UstawZaleznosc(ctx context.Context,
	z shared.OrchestrationDependencySetRequest) (shared.OrchestrationDependencySetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.OrchestrationDependencySetResponse{}, err
	}
	luk, err := lukUkladu(z.Dependency)
	if err != nil {
		return shared.OrchestrationDependencySetResponse{}, err
	}
	luki, err := a.lukiUkladu()
	if err != nil {
		return shared.OrchestrationDependencySetResponse{}, err
	}
	if err := luki.ZapiszZaleznosc(ctx, wiersz.ID, luk); err != nil {
		return shared.OrchestrationDependencySetResponse{}, bladAutomatyki(err)
	}
	kroki, zaleznosci, err := a.ukladPoZapisie(ctx, wiersz.ID)
	if err != nil {
		return shared.OrchestrationDependencySetResponse{}, err
	}
	zastrzezenia := zastrzezeniaUkladu(kroki, zaleznosci)
	return shared.OrchestrationDependencySetResponse{
		Dependencies: zaleznosciKontraktu(zaleznosci),
		Valid:        len(zastrzezenia) == 0,
		Issues:       zastrzezenia,
	}, nil
}

// ── orchestration.dependency.list ────────────────────────────────────────────

// WykazZaleznosci zwraca uklad zaleznosci automatyki, w razie wskazania kroku
// zawezony do lukow, ktore go dotykaja; krok wskazany musi istniec, inaczej
// odmowa niesie kod not_found z jego nazwa.
func (a *adapterAutomatyk) WykazZaleznosci(ctx context.Context,
	z shared.OrchestrationDependencyListRequest) (shared.OrchestrationDependencyListResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.OrchestrationDependencyListResponse{}, err
	}
	kroki, zaleznosci, err := a.ukladPoZapisie(ctx, wiersz.ID)
	if err != nil {
		return shared.OrchestrationDependencyListResponse{}, err
	}
	kodKroku := wartoscTekstu(z.StepId)
	if kodKroku != "" {
		if !czyKrokUkladu(kroki, kodKroku) {
			return shared.OrchestrationDependencyListResponse{},
				bladBrakuKrokuUkladu(z.WorkflowId, kodKroku)
		}
		zaleznosci = lukiKroku(zaleznosci, kodKroku)
	}
	return shared.OrchestrationDependencyListResponse{
		Dependencies: zaleznosciKontraktu(zaleznosci),
	}, nil
}

// ── orchestration.dependency.remove ──────────────────────────────────────────

// UsunZaleznoscUkladu zdejmuje jeden łuk układu i oddaje układ po usunięciu.
//
// Łuku, którego nie ma, nie da się usunąć, a odpowiedź `removed: false` byłaby
// tu ciszą udającą wynik. Odmowa niesie kod `not_found` i nazwy obu kroków.
func (a *adapterAutomatyk) UsunZaleznoscUkladu(ctx context.Context,
	z shared.OrchestrationDependencyRemoveRequest) (shared.OrchestrationDependencyRemoveResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.OrchestrationDependencyRemoveResponse{}, err
	}
	if z.FromStepId == "" || z.ToStepId == "" {
		return shared.OrchestrationDependencyRemoveResponse{},
			bladWskazaniaAutomatyki("usunięcie zależności bez wskazania obu kroków")
	}
	luki, err := a.lukiUkladu()
	if err != nil {
		return shared.OrchestrationDependencyRemoveResponse{}, err
	}
	zdjeta, err := luki.UsunZaleznosc(ctx, wiersz.ID, z.FromStepId, z.ToStepId)
	if err != nil {
		return shared.OrchestrationDependencyRemoveResponse{}, bladAutomatyki(err)
	}
	if !zdjeta {
		return shared.OrchestrationDependencyRemoveResponse{},
			bladBrakuZaleznosci(z.WorkflowId, z.FromStepId, z.ToStepId)
	}
	zaleznosci, err := a.repozytorium.Zaleznosci(ctx, wiersz.ID)
	if err != nil {
		return shared.OrchestrationDependencyRemoveResponse{}, bladAutomatyki(err)
	}
	return shared.OrchestrationDependencyRemoveResponse{
		Dependencies: zaleznosciKontraktu(zaleznosci),
		Removed:      true,
	}, nil
}

// ── orchestration.validate ───────────────────────────────────────────────────

// SprawdzUkladZaleznosci ocenia układ zależności automatyki: cykle, kroki
// osierocone i ścieżkę krytyczną. Niczego nie zapisuje.
func (a *adapterAutomatyk) SprawdzUkladZaleznosci(ctx context.Context,
	z shared.OrchestrationValidateRequest) (shared.OrchestrationValidateResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.OrchestrationValidateResponse{}, err
	}
	kroki, zaleznosci, err := a.ukladPoZapisie(ctx, wiersz.ID)
	if err != nil {
		return shared.OrchestrationValidateResponse{}, err
	}
	zastrzezenia := zastrzezeniaUkladu(kroki, zaleznosci)
	return shared.OrchestrationValidateResponse{
		Valid:               len(zastrzezenia) == 0,
		Issues:              zastrzezenia,
		CriticalPathStepIds: sciezkaKrytyczna(kroki, zaleznosci),
	}, nil
}

// ── wspólne ustalenia rodziny ────────────────────────────────────────────────

// zastrzezeniaUkladu skada komplet zastrzezen rodziny orchestration.*:
// zastrzezenia okna Orchestrator oraz kroki osierocone, ktorych wprost wymaga
// kontrakt orchestration.validate; dependency.set i validate licza pole valid
// jedna miara.
func zastrzezeniaUkladu(kroki []dane.KrokAutomatyki, zaleznosci []dane.ZaleznoscKroku) []string {
	zastrzezenia := sprawdzUklad(kroki, zaleznosci)
	return append(zastrzezenia, krokiOsierocone(kroki, zaleznosci)...)
}

// krokiOsierocone wskazuje kroki, ktorych nie dotyka zaden luk ukladu.
// Automatyka z samymi krokami, jeszcze niepowiazanymi, jest stanem poprawnym;
// sierota istnieje dopiero, gdy uklad ma luki i wiecej niz jeden krok.
func krokiOsierocone(kroki []dane.KrokAutomatyki, zaleznosci []dane.ZaleznoscKroku) []string {
	if len(zaleznosci) == 0 || len(kroki) < 2 {
		return nil
	}
	dotkniete := make(map[string]bool, len(zaleznosci)*2)
	for _, zaleznosc := range zaleznosci {
		dotkniete[zaleznosc.KrokZ] = true
		dotkniete[zaleznosc.KrokDo] = true
	}
	osierocone := []string{}
	for _, krok := range kroki {
		if !dotkniete[krok.Kod] {
			osierocone = append(osierocone, "krok nie wchodzi w żaden łuk układu: "+krok.Kod)
		}
	}
	sort.Strings(osierocone)
	return osierocone
}

// lukUkladu przekłada zależność kontraktu na wiersz układu i odmawia żądaniom,
// których schemat nie zna. Łuk bez rodzaju jest sekwencyjny — to znaczenie
// domyślne „krok po kroku”, tak samo jak przy `automation.orchestrator.define`.
func lukUkladu(zaleznosc shared.AutomationDependency) (dane.ZaleznoscKroku, error) {
	if zaleznosc.FromStepId == "" || zaleznosc.ToStepId == "" {
		return dane.ZaleznoscKroku{},
			bladWskazaniaAutomatyki("zależność bez wskazania obu kroków")
	}
	if zaleznosc.FromStepId == zaleznosc.ToStepId {
		return dane.ZaleznoscKroku{}, bladWskazaniaAutomatyki(
			"krok nie może zależeć od siebie: " + zaleznosc.FromStepId)
	}
	rodzaj := string(zaleznosc.Kind)
	if rodzaj == "" {
		rodzaj = shared.AutomationDependencyKindSequential
	}
	if !czyRodzajZaleznosci(rodzaj) {
		return dane.ZaleznoscKroku{}, bladWskazaniaAutomatyki(
			"nieznany rodzaj zależności: " + rodzaj)
	}
	return dane.ZaleznoscKroku{
		KrokZ: zaleznosc.FromStepId, KrokDo: zaleznosc.ToStepId,
		Rodzaj: rodzaj, Warunek: zaleznosc.Condition,
	}, nil
}

// czyRodzajZaleznosci pilnuje trzech wartości kontraktu. Ten sam zbiór jest
// więzem schematu (`CHECK(rodzaj IN (...))` w
// `store/migracja_039_automatyki.sql`); sprawdzenie tutaj zamienia usterkę
// zapisu na czytelną odmowę żądania.
func czyRodzajZaleznosci(rodzaj string) bool {
	switch rodzaj {
	case shared.AutomationDependencyKindSequential,
		shared.AutomationDependencyKindParallel,
		shared.AutomationDependencyKindConditional:
		return true
	default:
		return false
	}
}

// lukiKroku zawęża układ do łuków dotykających wskazanego kroku — wchodzących
// i wychodzących. Okno Orchestrator zaznacza krok i pyta „co go wiąże”, a nie
// „co z niego wychodzi”.
func lukiKroku(zaleznosci []dane.ZaleznoscKroku, kod string) []dane.ZaleznoscKroku {
	zawezone := make([]dane.ZaleznoscKroku, 0, len(zaleznosci))
	for _, zaleznosc := range zaleznosci {
		if zaleznosc.KrokZ == kod || zaleznosc.KrokDo == kod {
			zawezone = append(zawezone, zaleznosc)
		}
	}
	return zawezone
}

// czyKrokUkladu mowi, czy dana automatyka niesie krok o wskazanym kodzie
// w wykazie swoich krokow ukladu.
func czyKrokUkladu(kroki []dane.KrokAutomatyki, kod string) bool {
	for _, krok := range kroki {
		if krok.Kod == kod {
			return true
		}
	}
	return false
}

// ukladPoZapisie odczytuje kroki i łuki automatyki — wszystkie cztery komendy
// rodziny potrzebują obu wykazów naraz.
func (a *adapterAutomatyk) ukladPoZapisie(ctx context.Context,
	automatykaID int64) ([]dane.KrokAutomatyki, []dane.ZaleznoscKroku, error) {

	kroki, err := a.repozytorium.Kroki(ctx, automatykaID)
	if err != nil {
		return nil, nil, bladAutomatyki(err)
	}
	zaleznosci, err := a.repozytorium.Zaleznosci(ctx, automatykaID)
	if err != nil {
		return nil, nil, bladAutomatyki(err)
	}
	return kroki, zaleznosci, nil
}

// lukiUkladu wydobywa z repozytorium modułu czynności pojedynczego łuku.
// Repozytorium bez tego rozszerzenia jest usterką montażu, nie stanem układu,
// więc odmowa niesie kod `internal_error`.
func (a *adapterAutomatyk) lukiUkladu() (repozytoriumZaleznosciKrokow, error) {
	luki, ok := a.repozytorium.(repozytoriumZaleznosciKrokow)
	if !ok {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Automations: repozytorium nie niesie czynności pojedynczej zależności"))
	}
	return luki, nil
}

// bladBrakuZaleznosci nazywa luk miedzy dwoma wskazanymi krokami, ktorego
// w ukladzie automatyki nie ma.
func bladBrakuZaleznosci(automatyka, krokZ, krokDo string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Automations: zależność "+krokZ+" → "+krokDo+
			" nie istnieje w układzie automatyki "+automatyka))
}

// bladBrakuKrokuUkladu nazywa krok, ktorego wskazana automatyka nie niesie
// w swoim ukladzie zaleznosci.
func bladBrakuKrokuUkladu(automatyka, krok string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Automations: krok "+krok+" nie istnieje w automatyce "+automatyka))
}
