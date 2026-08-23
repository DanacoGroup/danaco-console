// Odpowiedzialność pliku: `context.usage.get` — zajętość okna kontekstu
// rozmowy: ile żetonów zużyto, gdzie jest granica i jak rozkłada się to na
// warstwy.
//
// ── Tokenizator jest podsystemem rdzenia, nie oszacowaniem ──────────────────
// Kontrakt tej komendy stawia sprawę wprost: bez tokenizatora liczba żetonów
// byłaby wartością wziętą znikąd. Rdzeń ma więc tokenizator wkompilowany
// (`server/internal/tokenizator`) i liczy nim NAPRAWDĘ — słownikiem BPE, tym
// samym podziałem, którym liczy model. Odpowiedź niesie nazwę słownika
// (`ContextUsage.tokenizer`), więc czytelnik wie, czym zmierzono.
//
// Przybliżenie „znaki podzielone przez cztery" byłoby tu gorsze niż brak
// odpowiedzi: myli się na polszczyźnie o kilkadziesiąt procent, a Operator
// czytający „w normie" traci turę na przepełnieniu okna.
//
// ── Granica okna jest odczytana, nie zmyślona ───────────────────────────────
// Granicę bierzemy z parametru kanału (`contextWindow` w konfiguracji wiersza
// rejestru kanałów). Kanał, który jej nie podaje, daje odpowiedź
// `available: false` wraz z powodem i wskazaniem naprawy — bo pasek zajętości
// wobec granicy wziętej z głowy pokazywałby „w normie" albo „prawie pełne"
// zależnie od tego, co rdzeń akurat zgadł.
//
// ── Trzy warstwy, trzy osobne pomiary ───────────────────────────────────────
// Prompt systemowy, historia rozmowy i pamięć są liczone osobno, każda swoim
// tekstem. Suma jest sumą tych trzech, a nie osobnym pomiarem — inaczej części
// nie sumowałyby się do całości i okno pokazywałoby dwie prawdy naraz.
package core

import (
	"context"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/tokenizator"
	"danacoconsole/shared"
)

// parametrGranicyOkna to nazwa parametru konfiguracji kanału niosącego wielkość
// okna kontekstu modelu w żetonach. Kontrakt `channel.add` nie ma osobnego pola
// na tę liczbę, więc jedzie ona parametrem — tą samą drogą co `credentialRef`.
const parametrGranicyOkna = "contextWindow"

// adapterZajetosciKontekstu wypełnia port `ZajetoscKontekstu`.
type adapterZajetosciKontekstu struct {
	okna       dane.RepozytoriumOkien
	wiadomosci dane.RepozytoriumWiadomosci
	pamiec     dane.RepozytoriumPamieci
	kanaly     *models.Rejestr
	// tozsamosc oddaje prompt systemowy obowiązujący w oknie — ten sam, który
	// pojedzie do modelu. Drugi składacz promptu dałby drugą liczbę żetonów
	// warstwy systemowej i rozjechał się z pierwszym przy pierwszej zmianie
	// warstw.
	tozsamosc Tozsamosc
}

// nowyAdapterZajetosciKontekstu wiąże port z magazynami, z których składa się
// treść okna kontekstu.
func nowyAdapterZajetosciKontekstu(okna dane.RepozytoriumOkien,
	wiadomosci dane.RepozytoriumWiadomosci, pamiec dane.RepozytoriumPamieci,
	kanaly *models.Rejestr, tozsamosc Tozsamosc) *adapterZajetosciKontekstu {

	return &adapterZajetosciKontekstu{
		okna: okna, wiadomosci: wiadomosci, pamiec: pamiec, kanaly: kanaly,
		tozsamosc: tozsamosc,
	}
}

// ZajetoscKontekstu obsługuje `context.usage.get`.
func (a *adapterZajetosciKontekstu) ZajetoscKontekstu(ctx context.Context,
	z shared.ContextUsageGetRequest) (shared.ContextUsageGetResponse, error) {

	kodOkna := strings.TrimSpace(z.WindowId)
	if kodOkna == "" {
		return shared.ContextUsageGetResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeValidationFailed,
			"zajętość kontekstu: żądanie bez wskazania okna — kontekst należy do okna"))
	}
	if a.okna == nil || a.wiadomosci == nil {
		return niezmierzonaZajetosc("rdzeń nie ma wpiętych magazynów okna i historii rozmowy, " +
			"więc nie ma czego zmierzyć"), nil
	}

	okno, err := a.okna.PoIdentyfikatorze(ctx, kodOkna)
	if err != nil {
		return shared.ContextUsageGetResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "zajętość kontekstu: nie ma okna o identyfikatorze "+kodOkna))
	}

	model := strings.TrimSpace(wartoscTekstu(z.Model))
	granica := 0
	if model == "" || granica == 0 {
		modelOkna, granicaOkna := a.modelIGranicaOkna(okno)
		if model == "" {
			model = modelOkna
		}
		granica = granicaOkna
	}
	if granica <= 0 {
		return niezmierzonaZajetosc("kanał modelu tego okna nie podaje wielkości okna kontekstu; " +
			"naprawa: dopisać parametr " + parametrGranicyOkna + " do konfiguracji kanału " +
			"(channel.update), bo pasek zajętości wobec granicy zgadniętej przez rdzeń " +
			"pokazywałby liczbę, której nikt nie ustalił"), nil
	}

	licznik, err := tokenizator.DlaModelu(model)
	if err != nil {
		return niezmierzonaZajetosc("tokenizator rdzenia nie zbudował słownika dla modelu „" +
			model + "”: " + err.Error()), nil
	}

	zetonyPromptu := licznik.Policz(a.promptSystemowyOkna(ctx, kodOkna))
	zetonyHistorii := a.zetonyHistorii(ctx, licznik, okno)
	zetonyPamieci := a.zetonyPamieci(ctx, licznik, okno)
	razem := zetonyPromptu + zetonyHistorii + zetonyPamieci

	nazwaSlownika := licznik.Nazwa()
	return shared.ContextUsageGetResponse{
		Usage: shared.ContextUsage{
			WindowId:           kodOkna,
			UsedTokens:         razem,
			LimitTokens:        granica,
			SystemPromptTokens: &zetonyPromptu,
			HistoryTokens:      &zetonyHistorii,
			MemoryTokens:       &zetonyPamieci,
			Tokenizer:          nazwaSlownika,
			MeasuredAt:         time.Now().UnixMilli(),
		},
		Available: true,
	}, nil
}

// modelIGranicaOkna odczytuje model kanału okna i wielkość jego okna kontekstu.
func (a *adapterZajetosciKontekstu) modelIGranicaOkna(okno dane.Okno) (string, int) {
	if a.kanaly == nil {
		return "", 0
	}
	// Rejestr kanałów jest kluczowany kodem wiersza, a okno trzyma jego numer,
	// więc wiersz odnajdujemy po numerze. Przejście po wykazie jest tu tańsze
	// niż drugie odpytanie bazy: rejestr stoi w pamięci i liczy pozycje
	// w dziesiątkach.
	var definicja models.Definicja
	znaleziona := false
	for _, wpis := range a.kanaly.Wykaz() {
		if wpis.Id == okno.KanalModeluID {
			definicja, znaleziona = wpis, true
			break
		}
	}
	if !znaleziona {
		return "", 0
	}
	granica, err := strconv.Atoi(strings.TrimSpace(definicja.Parametr(parametrGranicyOkna)))
	if err != nil || granica <= 0 {
		return definicja.Model, 0
	}
	return definicja.Model, granica
}

// promptSystemowyOkna składa treść, którą okno wkłada przed rozmową.
//
// Warstwy bierzemy z tego samego portu tożsamości, którym składa je tura
// (`adapter_rozmowa_tozsamosc.go`). Własne sklejanie warstw byłoby drugim
// składaczem promptu — i pierwszą różnicą między tym, co zmierzono, a tym, co
// pojechało do modelu.
func (a *adapterZajetosciKontekstu) promptSystemowyOkna(ctx context.Context, kodOkna string) string {
	if a.tozsamosc == nil {
		return ""
	}
	okno := kodOkna
	obowiazujaca, err := a.tozsamosc.Obowiazujaca(ctx,
		shared.IdentityEffectiveGetRequest{WindowId: &okno})
	if err != nil {
		return ""
	}
	warstwy := make([]string, 0, len(obowiazujaca.Layers))
	for _, warstwa := range obowiazujaca.Layers {
		if strings.TrimSpace(warstwa.Content) != "" {
			warstwy = append(warstwy, warstwa.Content)
		}
	}
	return strings.Join(warstwy, "\n\n")
}

// zetonyHistorii mierzy całą historię rozmowy okna.
//
// Historia idzie bez granicy liczby wiadomości: pytanie brzmi „ile zajmuje
// kontekst", a kontekst obejmuje wszystko, co wejdzie do wywołania. Ucięcie
// wykazu dałoby liczbę mniejszą od prawdy przy każdej dłuższej rozmowie —
// czyli dokładnie tam, gdzie odpowiedź jest potrzebna.
func (a *adapterZajetosciKontekstu) zetonyHistorii(ctx context.Context,
	licznik tokenizator.Licznik, okno dane.Okno) int {

	wiadomosci, err := a.wiadomosci.ListaOkna(ctx, okno.ID, 0)
	if err != nil {
		return 0
	}
	razem := 0
	for _, wiadomosc := range wiadomosci {
		if wiadomosc.Tresc == nil {
			continue
		}
		razem += licznik.Policz(*wiadomosc.Tresc)
	}
	return razem
}

// zetonyPamieci mierzy wpisy pamięci wchodzące do wywołania.
func (a *adapterZajetosciKontekstu) zetonyPamieci(ctx context.Context,
	licznik tokenizator.Licznik, okno dane.Okno) int {

	if a.pamiec == nil {
		return 0
	}
	kodOkna := ""
	if okno.IdentyfikatorZewnetrzny != nil {
		kodOkna = *okno.IdentyfikatorZewnetrzny
	}
	if kodOkna == "" {
		return 0
	}
	wpisy, err := a.pamiec.ListaPoziomu(ctx, dane.PoziomPamieci(shared.ConfigScopeWindow), kodOkna)
	if err != nil {
		return 0
	}
	razem := 0
	for _, wpis := range wpisy {
		if wpis.Tresc == nil {
			continue
		}
		razem += licznik.Policz(*wpis.Tresc)
	}
	return razem
}

// niezmierzonaZajetosc oddaje uczciwe „nie zmierzyłem" wraz z powodem.
//
// To NIE jest odmowa: kontrakt tej komendy przewiduje `available: false` wprost
// i chce, żeby okno pokazało wtedy brak miernika, a nie błąd. Pomiar zerowy
// podany jako zmierzony byłby paskiem „pusto" przy oknie pełnym po brzegi.
func niezmierzonaZajetosc(powod string) shared.ContextUsageGetResponse {
	return shared.ContextUsageGetResponse{Available: false, Reason: &powod}
}
