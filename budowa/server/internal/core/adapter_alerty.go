// Odpowiedzialność pliku: rodzina `alert.*` — reguły wyzwalania, rejestr
// wyzwoleń i ich potwierdzanie. Sam pomiar miar leży
// w `adapter_alerty_miara.go`; port i wpięcie w `handlers_alerty.go`.
//
// ── Kiedy reguły są ewaluowane ──────────────────────────────────────────────
// Kontrakt nie ma komendy „przelicz reguły". Ewaluacja idzie więc tam, gdzie
// Operator pyta o wynik: przy odczycie rejestru wyzwoleń (`alert.trigger.list`).
// Wykaz powstaje z reguł przeliczonych W CHWILI ODPOWIEDZI, a nie z wierszy
// odłożonych kiedyś przez zegar, którego kontrakt nie zna. Alert, który
// pokazuje stan sprzed godziny jako stan bieżący, jest tą samą fasadą, co
// sonda oddająca „w porządku" bez pomiaru.
//
// ── Reguła, która nigdy by się nie wyzwoliła, nie powstaje ──────────────────
// `alert.rule.save` odmawia miary, której rdzeń nie umie zmierzyć, nazywając
// ją wprost. Zapisanie takiej reguły dałoby Operatorowi wiersz w wykazie
// i ciszę zamiast alertu — czyli dokładnie to złudzenie bezpieczeństwa, przed
// którym rodzina ma chronić.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	przedrostekReguluAlertu     = "regula-alertu-"
	przedrostekWyzwoleniaAlertu = "wyzwolenie-"

	// domyslnyZakresWyzwolen obejmuje dobę wstecz, gdy żądanie nie stawia
	// własnych granic czasu odczytu rejestru.
	domyslnyZakresWyzwolen = 24 * time.Hour
)

// adapterAlertow wypełnia port `Alerty`.
//
// Cztery źródła pomiaru ponad własny magazyn: ślad wywołań modelu (koszt,
// żetony, opóźnienie, udział niepowodzeń), dziennik błędów (liczba błędów),
// rozstrzygacz konfiguracji (pułap kosztu dla miary budżetu) oraz nadajnik
// zdarzeń, którym wyzwolenie dociera do okien. Brak któregokolwiek nie psuje
// montażu: psuje te miary, które z niego liczą — i te miary adapter wtedy
// odrzuca przy zapisie reguły, zamiast milczeć w trakcie ewaluacji.
type adapterAlertow struct {
	repozytorium dane.RepozytoriumAlertow
	prowenancja  dane.RepozytoriumProwenancji
	diagnostyka  dane.RepozytoriumDiagnostyki
	rozstrzygacz *konfig.Rozstrzygacz
	nadajnik     *emiter
	// centrum jest rejestrem centrum powiadomień. Wyzwolony alert jest
	// zdarzeniem klasy `blad` — zdarzeniem, po którym Operator ma sięgnąć do
	// platformy. Bez tego wpięcia alarm żyłby wyłącznie w oknie otwartym
	// w chwili wyzwolenia.
	centrum *adapterCentrumPowiadomien
}

// nowyAdapterAlertow wiąże port z magazynem reguł.
func nowyAdapterAlertow(repozytorium dane.RepozytoriumAlertow) *adapterAlertow {
	return &adapterAlertow{repozytorium: repozytorium}
}

// ZeZrodlamiMiar wpina magazyny, z których liczą się miary reguł.
func (a *adapterAlertow) ZeZrodlamiMiar(prowenancja dane.RepozytoriumProwenancji,
	diagnostyka dane.RepozytoriumDiagnostyki, rozstrzygacz *konfig.Rozstrzygacz) *adapterAlertow {

	a.prowenancja = prowenancja
	a.diagnostyka = diagnostyka
	a.rozstrzygacz = rozstrzygacz
	return a
}

// ZWyjsciem wpina nadajnik, którym idzie zdarzenie `alert.triggered`.
func (a *adapterAlertow) ZWyjsciem(e *emiter) *adapterAlertow {
	a.nadajnik = e
	return a
}

// ZCentrumPowiadomien wpina rejestr centrum. Bez niego alerty rozgłaszają się
// jak dotąd, ale nie zostawiają śladu w centrum.
func (a *adapterAlertow) ZCentrumPowiadomien(c *adapterCentrumPowiadomien) *adapterAlertow {
	a.centrum = c
	return a
}

// ZapiszRegule obsługuje `alert.rule.save`.
func (a *adapterAlertow) ZapiszRegule(ctx context.Context,
	z shared.AlertRuleSaveRequest) (shared.AlertRuleSaveResponse, error) {

	if a.repozytorium == nil {
		return shared.AlertRuleSaveResponse{}, bladZapleczaAlertow()
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.AlertRuleSaveResponse{}, bladWskazaniaAlertu("reguła bez nazwy")
	}
	if err := sprawdzRodzajReguly(z.Kind); err != nil {
		return shared.AlertRuleSaveResponse{}, err
	}
	if err := a.sprawdzMiareReguly(z.Metric); err != nil {
		return shared.AlertRuleSaveResponse{}, err
	}
	if z.WindowMs <= 0 {
		return shared.AlertRuleSaveResponse{}, bladWskazaniaAlertu(
			"okno czasu liczenia miary musi być dodatnie — bez niego nie wiadomo, " +
				"z jakiego okresu brać pomiar")
	}
	if len(z.Channels) == 0 {
		return shared.AlertRuleSaveResponse{}, bladWskazaniaAlertu(
			"reguła bez drogi dostarczenia nie zawoła nikogo — wskaż co najmniej jeden kanał")
	}
	for _, kanal := range z.Channels {
		if err := sprawdzKanalAlertu(kanal); err != nil {
			return shared.AlertRuleSaveResponse{}, err
		}
	}

	teraz := terazWMilisekundachAlertow()
	kod := strings.TrimSpace(wartoscTekstu(z.RuleId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekReguluAlertu)
	}
	regula := dane.RegulaAlertu{
		Kod:            kod,
		Nazwa:          strings.TrimSpace(z.Name),
		Opis:           z.Description,
		Rodzaj:         string(z.Kind),
		Miara:          string(z.Metric),
		Prog:           z.Threshold,
		OknoMs:         z.WindowMs,
		Waga:           string(z.Severity),
		KanalyJSON:     zapisKanalowAlertu(z.Channels),
		ZasiegKod:      z.ScopeId,
		Czynna:         z.Enabled == nil || *z.Enabled,
		WyciszonaDo:    z.MuteUntil,
		EskalacjaPoMs:  z.EscalateAfterMs,
		AdresZwrotny:   z.WebhookUrl,
		Utworzono:      teraz,
		Zaktualizowano: &teraz,
	}
	if z.Comparison != nil {
		porownanie := string(*z.Comparison)
		regula.Porownanie = &porownanie
	}
	if z.Scope != nil {
		zasieg := string(*z.Scope)
		regula.Zasieg = &zasieg
	}

	zapisana, powstala, err := a.repozytorium.ZapiszRegule(ctx, regula)
	if err != nil {
		return shared.AlertRuleSaveResponse{}, bladMagazynuAlertow(err)
	}
	return shared.AlertRuleSaveResponse{Rule: regulaKontraktu(zapisana), Created: powstala}, nil
}

// WykazRegul obsługuje `alert.rule.list`.
func (a *adapterAlertow) WykazRegul(ctx context.Context,
	z shared.AlertRuleListRequest) (shared.AlertRuleListResponse, error) {

	if a.repozytorium == nil {
		return shared.AlertRuleListResponse{}, bladZapleczaAlertow()
	}
	rodzaj, miara := "", ""
	if z.Kind != nil {
		rodzaj = string(*z.Kind)
	}
	if z.Metric != nil {
		miara = string(*z.Metric)
	}
	tylkoCzynne := z.Enabled != nil && *z.Enabled
	reguly, err := a.repozytorium.Reguly(ctx, rodzaj, miara, tylkoCzynne, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.AlertRuleListResponse{}, bladMagazynuAlertow(err)
	}
	wykaz := make([]shared.AlertRule, 0, len(reguly))
	for _, regula := range reguly {
		wykaz = append(wykaz, regulaKontraktu(regula))
	}
	liczba := len(wykaz)
	return shared.AlertRuleListResponse{Rules: wykaz, Total: &liczba}, nil
}

// UsunRegule obsługuje `alert.rule.remove`.
func (a *adapterAlertow) UsunRegule(ctx context.Context,
	z shared.AlertRuleRemoveRequest) (shared.AlertRuleRemoveResponse, error) {

	if a.repozytorium == nil {
		return shared.AlertRuleRemoveResponse{}, bladZapleczaAlertow()
	}
	usunietych, err := a.repozytorium.UsunRegule(ctx, strings.TrimSpace(z.RuleId))
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.AlertRuleRemoveResponse{}, bladNieznanejRegulyAlertu(z.RuleId)
		}
		return shared.AlertRuleRemoveResponse{}, bladMagazynuAlertow(err)
	}
	return shared.AlertRuleRemoveResponse{RuleId: z.RuleId, RemovedTriggers: usunietych}, nil
}

// WykazWyzwolen obsługuje `alert.trigger.list`.
//
// Przed odczytem rejestru reguły są przeliczane na bieżących pomiarach. Dzięki
// temu wykaz jest stanem produktu TERAZ, a nie zapisem sprzed nieokreślonego
// czasu. Niepowodzenie ewaluacji nie przewraca odczytu: rejestr zastany jest
// wartościowszy niż odmowa, a wyzwolenie, którego nie dało się policzyć, i tak
// nie miałoby wartości obserwowanej.
func (a *adapterAlertow) WykazWyzwolen(ctx context.Context,
	z shared.AlertTriggerListRequest) (shared.AlertTriggerListResponse, error) {

	if a.repozytorium == nil {
		return shared.AlertTriggerListResponse{}, bladZapleczaAlertow()
	}
	a.przelicz(ctx)

	sito := dane.SitoWyzwolenAlertu{
		RegulaKod: wartoscTekstu(z.RuleId),
		OdCzasu:   wartoscLiczbyDlugiej(z.FromTime),
		DoCzasu:   wartoscLiczbyDlugiej(z.ToTime),
		Granica:   wartoscLiczby(z.Limit),
	}
	if z.Status != nil {
		sito.Stan = string(*z.Status)
	}
	if z.Severity != nil {
		sito.Waga = string(*z.Severity)
	}
	wyzwolenia, wszystkich, err := a.repozytorium.Wyzwolenia(ctx, sito)
	if err != nil {
		return shared.AlertTriggerListResponse{}, bladMagazynuAlertow(err)
	}
	wykaz := make([]shared.AlertTrigger, 0, len(wyzwolenia))
	for _, wyzwolenie := range wyzwolenia {
		wykaz = append(wykaz, wyzwolenieKontraktu(wyzwolenie))
	}
	przyciety := wszystkich > len(wykaz)
	return shared.AlertTriggerListResponse{
		Triggers: wykaz, Total: &wszystkich, Truncated: &przyciety,
	}, nil
}

// PotwierdzWyzwolenie obsługuje `alert.trigger.acknowledge`.
func (a *adapterAlertow) PotwierdzWyzwolenie(ctx context.Context,
	z shared.AlertTriggerAcknowledgeRequest) (shared.AlertTriggerAcknowledgeResponse, error) {

	if a.repozytorium == nil {
		return shared.AlertTriggerAcknowledgeResponse{}, bladZapleczaAlertow()
	}
	wyzwolenie, err := a.repozytorium.PotwierdzWyzwolenie(ctx, strings.TrimSpace(z.TriggerId),
		terazWMilisekundachAlertow(), z.Note)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.AlertTriggerAcknowledgeResponse{}, protocol.JakoError(
				protocol.NowyBlad(shared.ErrorCodeNotFound,
					"alerty: nie ma wyzwolenia o identyfikatorze "+strings.TrimSpace(z.TriggerId)))
		}
		return shared.AlertTriggerAcknowledgeResponse{}, bladMagazynuAlertow(err)
	}
	return shared.AlertTriggerAcknowledgeResponse{Trigger: wyzwolenieKontraktu(wyzwolenie)}, nil
}

// regulaKontraktu przekłada wiersz na regułę kontraktu.
func regulaKontraktu(r dane.RegulaAlertu) shared.AlertRule {
	licznik := r.LiczbaWyzwolen
	regula := shared.AlertRule{
		Id: r.Kod, Name: r.Nazwa, Description: r.Opis,
		Kind: shared.AlertRuleKind(r.Rodzaj), Metric: shared.AlertMetric(r.Miara),
		Threshold: r.Prog, WindowMs: r.OknoMs,
		Severity: shared.DiagnosticPriority(r.Waga),
		Channels: odczytKanalowAlertu(r.KanalyJSON),
		ScopeId:  r.ZasiegKod, Enabled: r.Czynna, MuteUntil: r.WyciszonaDo,
		EscalateAfterMs: r.EskalacjaPoMs, WebhookUrl: r.AdresZwrotny,
		CreatedAt: r.Utworzono, UpdatedAt: r.Zaktualizowano,
		LastTriggeredAt: r.OstatnieWyzwolenie, TriggerCount: &licznik,
	}
	if r.Porownanie != nil {
		porownanie := shared.AlertComparison(*r.Porownanie)
		regula.Comparison = &porownanie
	}
	if r.Zasieg != nil {
		zasieg := shared.ConfigScope(*r.Zasieg)
		regula.Scope = &zasieg
	}
	return regula
}

// wyzwolenieKontraktu przekłada wiersz rejestru na wyzwolenie kontraktu.
func wyzwolenieKontraktu(w dane.WyzwolenieAlertu) shared.AlertTrigger {
	return shared.AlertTrigger{
		Id: w.Kod, RuleId: w.RegulaKod, RuleName: w.RegulaNazwa,
		Status: shared.AlertTriggerStatus(w.Stan), Severity: shared.DiagnosticPriority(w.Waga),
		Metric: shared.AlertMetric(w.Miara), ObservedValue: w.WartoscObserwowana,
		Threshold: w.Prog, Message: w.Komunikat, FiredAt: w.Wyzwolono,
		AcknowledgedAt: w.Potwierdzono, ResolvedAt: w.Rozwiazano, Note: w.Notatka,
		ErrorId: w.BladKod, ProbeId: w.SondaKod, CallId: w.WywolanieKod,
		DeliveredChannels: odczytKanalowAlertu(w.KanalyDostarczoneJSON),
	}
}

// zapisKanalowAlertu składa wykaz dróg dostarczenia w zapis strukturalny
// kolumny.
func zapisKanalowAlertu(kanaly []shared.AlertChannel) string {
	if len(kanaly) == 0 {
		return "[]"
	}
	bajty, err := json.Marshal(kanaly)
	if err != nil {
		return "[]"
	}
	return string(bajty)
}

// odczytKanalowAlertu rozkłada zapis strukturalny kolumny na wykaz dróg.
// Zapis nieczytelny daje wykaz pusty, a nie awarię odczytu: reguła bez dróg
// jest widoczna i naprawialna, reguła nieodczytana — nie.
func odczytKanalowAlertu(zapis string) []shared.AlertChannel {
	if strings.TrimSpace(zapis) == "" {
		return nil
	}
	var kanaly []shared.AlertChannel
	if err := json.Unmarshal([]byte(zapis), &kanaly); err != nil {
		return nil
	}
	return kanaly
}

// sprawdzRodzajReguly odbija rodzaj spoza wyliczenia kontraktu.
func sprawdzRodzajReguly(rodzaj shared.AlertRuleKind) error {
	switch rodzaj {
	case shared.AlertRuleKindThreshold, shared.AlertRuleKindAnomaly,
		shared.AlertRuleKindNewFingerprint:
		return nil
	default:
		return bladWskazaniaAlertu("nie znam rodzaju reguły „" + string(rodzaj) +
			"” — rdzeń zna: threshold, anomaly, newFingerprint")
	}
}

// sprawdzKanalAlertu odbija drogę dostarczenia spoza wyliczenia kontraktu.
func sprawdzKanalAlertu(kanal shared.AlertChannel) error {
	switch kanal {
	case shared.AlertChannelApp, shared.AlertChannelAlwaysOnDisplay,
		shared.AlertChannelMail, shared.AlertChannelWebhook:
		return nil
	default:
		return bladWskazaniaAlertu("nie znam drogi dostarczenia „" + string(kanal) +
			"” — rdzeń zna: app, alwaysOnDisplay, mail, webhook")
	}
}

// terazWMilisekundachAlertow oddaje bieżącą chwilę w milisekundach epoki.
func terazWMilisekundachAlertow() int64 {
	return time.Now().UnixMilli()
}

// zapisLiczbyMiary składa czytelny zapis wartości obserwowanej do treści
// komunikatu wyzwolenia.
func zapisLiczbyMiary(wartosc float64) string {
	return strconv.FormatFloat(wartosc, 'f', -1, 64)
}

// bladZapleczaAlertow nazywa brak magazynu reguł po stronie rdzenia.
func bladZapleczaAlertow() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"alerty: rdzeń nie ma wpiętego magazynu reguł — naprawa: podpiąć "+
			"repozytorium alertów przy składaniu rdzenia"))
}

// bladWskazaniaAlertu nazywa niepoprawne żądanie.
func bladWskazaniaAlertu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "alerty: "+powod))
}

// bladNieznanejRegulyAlertu nazywa wskazanie reguły, której nie ma.
func bladNieznanejRegulyAlertu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"alerty: nie ma reguły o identyfikatorze "+strings.TrimSpace(kod)))
}

// bladMagazynuAlertow nazywa niepowodzenie zapisu albo odczytu.
func bladMagazynuAlertow(err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"alerty: "+err.Error()))
}

// wyzwolenieAlertu rozgłasza `alert.triggered`. Wyzwolenie jest bytem
// przekrojowym, nie bytem jednej sesji, więc zdarzenie idzie bez wskazania
// karty — wzorem zdarzeń modułu Design o zasobach.
//
// Sprawcy zdarzenie nie niesie i nieść nie może: wyzwolenie powstaje z pomiaru
// rdzenia, a nie z czyjegoś kliknięcia. Kontrakt też nie ma na nie pola.
// Reguła jedzie razem z wyzwoleniem, bo okno pokazujące alert musi wiedzieć,
// czyj to alarm i jakim progiem był postawiony.
func (e *emiter) wyzwolenieAlertu(wyzwolenie shared.AlertTrigger, regula shared.AlertRule) {
	e.wyslij(shared.EventAlertTriggered, "", shared.AlertTriggeredEvent{
		Trigger: wyzwolenie, Rule: regula,
	})
}
