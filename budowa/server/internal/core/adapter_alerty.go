// Plik realizuje rodzinę adapterów `alert.*`: zapis i odczyt reguł wyzwalania
// oraz rejestru wyzwoleń kontraktu, wraz z potwierdzaniem pojedynczego wyzwolenia.
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

// adapterAlertow wypełnia port `Alerty`, licząc miary reguł z czterech źródeł
// ponad własny magazyn: śladu wywołań modelu, dziennika błędów, rozstrzygacza
// konfiguracji oraz nadajnika zdarzeń, którym wyzwolenie dociera do okien.
type adapterAlertow struct {
	repozytorium dane.RepozytoriumAlertow
	prowenancja  dane.RepozytoriumProwenancji
	diagnostyka  dane.RepozytoriumDiagnostyki
	rozstrzygacz *konfig.Rozstrzygacz
	nadajnik     *emiter
	// centrum jest rejestrem powiadomień, w którym wyzwolony alert zapisuje
	// się jako zdarzenie klasy błąd.
	centrum *adapterCentrumPowiadomien
}

// nowyAdapterAlertow wiąże port z magazynem reguł alertów, pozostawiając
// pozostałe źródła pomiaru i nadajnik zdarzeń do wpięcia osobnymi wywołaniami.
func nowyAdapterAlertow(repozytorium dane.RepozytoriumAlertow) *adapterAlertow {
	return &adapterAlertow{repozytorium: repozytorium}
}

// ZeZrodlamiMiar wpina magazyny prowenancji, diagnostyki i rozstrzygacza
// konfiguracji, z których liczą się miary progów zapisanych reguł.
func (a *adapterAlertow) ZeZrodlamiMiar(prowenancja dane.RepozytoriumProwenancji,
	diagnostyka dane.RepozytoriumDiagnostyki, rozstrzygacz *konfig.Rozstrzygacz) *adapterAlertow {

	a.prowenancja = prowenancja
	a.diagnostyka = diagnostyka
	a.rozstrzygacz = rozstrzygacz
	return a
}

// ZWyjsciem wpina nadajnik, którym idzie do okien zdarzenie `alert.triggered`
// po każdym wyzwoleniu reguły.
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

// ZapiszRegule obsługuje `alert.rule.save`, sprawdzając rodzaj, miarę, okno
// czasu i kanały dostarczenia reguły przed zapisem w magazynie.
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

// WykazRegul obsługuje `alert.rule.list`, filtrując reguły po rodzaju,
// mierze i stanie czynności zgodnie z żądaniem operatora.
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

// UsunRegule obsługuje `alert.rule.remove`, usuwając regułę razem z jej
// wyzwoleniami i zwracając liczbę usuniętych wierszy rejestru.
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

// WykazWyzwolen obsługuje `alert.trigger.list`. Przed odczytem rejestru
// reguły są przeliczane na bieżących pomiarach, więc wykaz jest stanem
// produktu w chwili odpowiedzi, a nie zapisem sprzed nieokreślonego czasu.
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

// PotwierdzWyzwolenie obsługuje `alert.trigger.acknowledge`, zapisując czas
// i treść potwierdzenia przy wskazanym wyzwoleniu rejestru.
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

// regulaKontraktu przekłada wiersz magazynu na regułę kontraktu, przenosząc
// porównanie i zasięg tylko wtedy, gdy wiersz je niesie.
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

// wyzwolenieKontraktu przekłada wiersz rejestru wyzwoleń na wyzwolenie
// kontraktu wraz z dekodowaniem kanałów, którymi zostało dostarczone.
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
// kolumny, oddając pustą tablicę, gdy reguła nie ma żadnego kanału.
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

// sprawdzRodzajReguly odbija rodzaj spoza wyliczenia kontraktu, nazywając
// w komunikacie błędu rodzaje, które rdzeń rozpoznaje.
func sprawdzRodzajReguly(rodzaj shared.AlertRuleKind) error {
	switch rodzaj {
	case shared.AlertRuleKindThreshold, shared.AlertRuleKindAnomaly,
		shared.AlertRuleKindNewFingerprint:
		return nil
	default:
		return bladWskazaniaAlertu("nie znam rodzaju reguły „" + string(rodzaj) +
			"” — serwer zna: threshold, anomaly, newFingerprint")
	}
}

// sprawdzKanalAlertu odbija drogę dostarczenia spoza wyliczenia kontraktu,
// nazywając w komunikacie błędu drogi, które rdzeń rozpoznaje.
func sprawdzKanalAlertu(kanal shared.AlertChannel) error {
	switch kanal {
	case shared.AlertChannelApp, shared.AlertChannelAlwaysOnDisplay,
		shared.AlertChannelMail, shared.AlertChannelWebhook:
		return nil
	default:
		return bladWskazaniaAlertu("nie znam drogi dostarczenia „" + string(kanal) +
			"” — serwer zna: app, alwaysOnDisplay, mail, webhook")
	}
}

// terazWMilisekundachAlertow oddaje bieżącą chwilę w milisekundach epoki,
// wartość zapisywaną przy tworzeniu i potwierdzaniu wierszy alertów.
func terazWMilisekundachAlertow() int64 {
	return time.Now().UnixMilli()
}

// zapisLiczbyMiary składa czytelny zapis wartości obserwowanej do treści
// komunikatu wyzwolenia, bez notacji wykładniczej właściwej zapisowi binarnemu.
func zapisLiczbyMiary(wartosc float64) string {
	return strconv.FormatFloat(wartosc, 'f', -1, 64)
}

// bladZapleczaAlertow nazywa brak magazynu reguł po stronie rdzenia, stan
// świadczący o niepełnym złożeniu adaptera przy starcie serwera.
func bladZapleczaAlertow() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"alerty: serwer nie ma wpiętego magazynu reguł — naprawa: podpiąć "+
			"repozytorium alertów przy składaniu serwera"))
}

// bladWskazaniaAlertu nazywa niepoprawne żądanie operatora, przenosząc powód
// odmowy wprost do treści komunikatu błędu.
func bladWskazaniaAlertu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "alerty: "+powod))
}

// bladNieznanejRegulyAlertu nazywa wskazanie reguły, której magazyn nie ma
// pod podanym identyfikatorem.
func bladNieznanejRegulyAlertu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"alerty: nie ma reguły o identyfikatorze "+strings.TrimSpace(kod)))
}

// bladMagazynuAlertow nazywa niepowodzenie zapisu albo odczytu magazynu reguł
// i wyzwoleń, przenosząc jego treść do komunikatu błędu rdzenia.
func bladMagazynuAlertow(err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"alerty: "+err.Error()))
}

// wyzwolenieAlertu rozgłasza `alert.triggered` bez wskazania karty sesji,
// ponieważ wyzwolenie jest bytem przekrojowym, a nie bytem jednej sesji, i
// niesie razem z sobą regułę, której próg został przekroczony.
func (e *emiter) wyzwolenieAlertu(ctx context.Context, wyzwolenie shared.AlertTrigger, regula shared.AlertRule) {
	e.wyslijDoKonta(ctx, shared.EventAlertTriggered, "", shared.AlertTriggeredEvent{
		Trigger: wyzwolenie, Rule: regula,
	})
}
