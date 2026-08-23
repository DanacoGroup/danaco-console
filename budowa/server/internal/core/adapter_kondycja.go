// Odpowiedzialność pliku: rodzina `health.*` — definicje sond kondycji, ich
// przebiegi, seria wyników i dostępność liczona z tej serii.
//
// ── Zasada tej rodziny: każda wartość pochodzi z pomiaru ────────────────────
// Komenda zdrowia, która oddaje „w porządku" nie zmierzywszy niczego, jest
// gorsza niż jej brak: to fasada, przez którą awaria przechodzi niezauważona.
// Dlatego w tym pliku nie ma ani jednej ścieżki oddającej stan `up` bez
// wykonanego pomiaru. Sonda, której rdzeń nie ma czym wykonać, kończy się
// stanem `unknown` wraz z powodem — a nie stanem `up` „bo nic nie zawiodło".
//
// ── Dostępność liczy się z wierszy, nie z licznika ──────────────────────────
// `health.uptime.get` nie czyta żadnej kolumny „dostępność". Bierze serię
// pomiarów z zakresu czasu i liczy udział wyników udanych. Licznik podnoszony
// przy zapisie rozjechałby się z serią przy pierwszym usunięciu wyników albo
// zmianie zakresu, a rozjazd byłby niewidoczny — obie liczby wyglądają tak samo.
//
// Sam przebieg sondy (co właściwie mierzy każdy rodzaj) leży
// w `adapter_kondycja_pomiar.go`; port i wpięcie w `handlers_kondycja.go`.
package core

import (
	"context"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// przedrostekSondyKondycji i przedrostekWynikuKondycji znakują identyfikatory
// nadawane przez rdzeń.
const (
	przedrostekSondyKondycji  = "sonda-"
	przedrostekWynikuKondycji = "pomiar-"

	// domyslnyLimitCzasuSondy jest brany, gdy definicja nie stawia własnego.
	// Pięć sekund to granica, po której odpowiedź i tak jest bezużyteczna dla
	// tego, kto na nią czeka.
	domyslnyLimitCzasuSondy = 5 * time.Second

	// gornyLimitCzasuSondy zamyka limit podany w definicji. Sonda z limitem
	// godzinnym zablokowałaby wywołanie `health.probe.run` na godzinę.
	gornyLimitCzasuSondy = 2 * time.Minute

	// domyslnyCelDostepnosci jest brany, gdy ani sonda, ani żądanie nie podają
	// celu. Nie jest to pomiar, tylko punkt odniesienia dla budżetu błędów,
	// więc wolno mu mieć wartość domyślną.
	domyslnyCelDostepnosci = 99.0

	// domyslnyZakresDostepnosci obejmuje dobę wstecz, gdy żądanie nie stawia
	// własnych granic czasu.
	domyslnyZakresDostepnosci = 24 * time.Hour
)

// adapterKondycji wypełnia port `Kondycja`.
//
// Trzy zależności ponad repozytorium to trzy rodzaje pomiaru, których rdzeń nie
// wykona sam: uruchamiacz procesów (sonda `command`), rejestr kanałów (sonda
// `modelCall`) oraz izolacja wraz z katalogiem roboczym — ta sama droga, którą
// idą wszystkie wołania arsenału. Brak którejkolwiek nie psuje montażu: psuje
// jeden rodzaj sondy, która wtedy oddaje `unknown` wraz z powodem.
type adapterKondycji struct {
	repozytorium dane.RepozytoriumKondycji
	uruchamiacz  session.Uruchamiacz
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	kanaly       *models.Rejestr
}

// nowyAdapterKondycji wiąże port z magazynem sond.
func nowyAdapterKondycji(repozytorium dane.RepozytoriumKondycji) *adapterKondycji {
	return &adapterKondycji{repozytorium: repozytorium}
}

// ZArsenalem wpina drogę startu procesu dla sond rodzaju `command`.
func (a *adapterKondycji) ZArsenalem(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterKondycji {

	a.uruchamiacz = uruchamiacz
	a.rozstrzygacz = rozstrzygacz
	a.katalog = katalog
	return a
}

// ZKanalami wpina rejestr kanałów dla sond rodzaju `modelCall`.
func (a *adapterKondycji) ZKanalami(kanaly *models.Rejestr) *adapterKondycji {
	a.kanaly = kanaly
	return a
}

// ZapiszSonde obsługuje `health.probe.save`.
func (a *adapterKondycji) ZapiszSonde(ctx context.Context,
	z shared.HealthProbeSaveRequest) (shared.HealthProbeSaveResponse, error) {

	if a.repozytorium == nil {
		return shared.HealthProbeSaveResponse{}, bladZapleczaKondycji()
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.HealthProbeSaveResponse{}, bladWskazaniaKondycji("sonda bez nazwy")
	}
	if strings.TrimSpace(z.Target) == "" {
		return shared.HealthProbeSaveResponse{}, bladWskazaniaKondycji(
			"sonda bez celu — nie wiadomo, co miałaby mierzyć")
	}
	if err := sprawdzRodzajSondy(z.Kind); err != nil {
		return shared.HealthProbeSaveResponse{}, err
	}
	if z.IntervalMs <= 0 {
		return shared.HealthProbeSaveResponse{}, bladWskazaniaKondycji(
			"odstęp między przebiegami musi być dodatni — sonda bez odstępu nigdy by nie ruszyła")
	}

	teraz := terazWMilisekundachKondycji()
	kod := strings.TrimSpace(wartoscTekstu(z.ProbeId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekSondyKondycji)
	}
	sonda := dane.SondaKondycji{
		Kod:            kod,
		Nazwa:          strings.TrimSpace(z.Name),
		Rodzaj:         string(z.Kind),
		Cel:            strings.TrimSpace(z.Target),
		KomponentKod:   z.ComponentId,
		OdstepMs:       z.IntervalMs,
		LimitCzasuMs:   z.TimeoutMs,
		OczekiwanyKod:  liczbaCalkowitaZeWskaznikaKondycji(z.ExpectedStatus),
		TrescWysylana:  z.Payload,
		CelDostepnosci: z.Objective,
		Czynna:         z.Enabled == nil || *z.Enabled,
		Utworzono:      teraz,
		Zaktualizowano: &teraz,
	}
	zapisana, powstala, err := a.repozytorium.ZapiszSonde(ctx, sonda)
	if err != nil {
		return shared.HealthProbeSaveResponse{}, bladMagazynuKondycji(err)
	}
	return shared.HealthProbeSaveResponse{
		Probe: sondaKontraktu(zapisana), Created: powstala,
	}, nil
}

// WykazSond obsługuje `health.probe.list`.
func (a *adapterKondycji) WykazSond(ctx context.Context,
	z shared.HealthProbeListRequest) (shared.HealthProbeListResponse, error) {

	if a.repozytorium == nil {
		return shared.HealthProbeListResponse{}, bladZapleczaKondycji()
	}
	rodzaj := ""
	if z.Kind != nil {
		rodzaj = string(*z.Kind)
	}
	tylkoCzynne := z.Enabled != nil && *z.Enabled
	sondy, err := a.repozytorium.Sondy(ctx, rodzaj, wartoscTekstu(z.ComponentId),
		tylkoCzynne, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.HealthProbeListResponse{}, bladMagazynuKondycji(err)
	}
	wykaz := make([]shared.HealthProbe, 0, len(sondy))
	for _, sonda := range sondy {
		wykaz = append(wykaz, sondaKontraktu(sonda))
	}
	liczba := len(wykaz)
	return shared.HealthProbeListResponse{Probes: wykaz, Total: &liczba}, nil
}

// UsunSonde obsługuje `health.probe.remove`.
func (a *adapterKondycji) UsunSonde(ctx context.Context,
	z shared.HealthProbeRemoveRequest) (shared.HealthProbeRemoveResponse, error) {

	if a.repozytorium == nil {
		return shared.HealthProbeRemoveResponse{}, bladZapleczaKondycji()
	}
	usunietych, err := a.repozytorium.UsunSonde(ctx, strings.TrimSpace(z.ProbeId))
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.HealthProbeRemoveResponse{}, bladNieznanejSondy(z.ProbeId)
		}
		return shared.HealthProbeRemoveResponse{}, bladMagazynuKondycji(err)
	}
	return shared.HealthProbeRemoveResponse{
		ProbeId: z.ProbeId, RemovedResults: usunietych,
	}, nil
}

// WykonajSonde obsługuje `health.probe.run` — pomiar na żądanie.
//
// Wynik zapisujemy zawsze, także gdy sonda zawiodła i także gdy rdzeń nie miał
// czym jej wykonać. Pomiar niezapisany byłby pomiarem, którego dostępność nie
// zobaczy — a właśnie z tych wierszy liczy się dostępność.
func (a *adapterKondycji) WykonajSonde(ctx context.Context,
	z shared.HealthProbeRunRequest) (shared.HealthProbeRunResponse, error) {

	if a.repozytorium == nil {
		return shared.HealthProbeRunResponse{}, bladZapleczaKondycji()
	}
	sonda, err := a.repozytorium.Sonda(ctx, strings.TrimSpace(z.ProbeId))
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.HealthProbeRunResponse{}, bladNieznanejSondy(z.ProbeId)
		}
		return shared.HealthProbeRunResponse{}, bladMagazynuKondycji(err)
	}

	pomiar := a.zmierz(ctx, sonda)
	pomiar.Kod = nowyIdentyfikator(przedrostekWynikuKondycji)
	pomiar.SondaKod = sonda.Kod
	zapisany, err := a.repozytorium.ZapiszWynik(ctx, pomiar)
	if err != nil {
		return shared.HealthProbeRunResponse{}, bladMagazynuKondycji(err)
	}
	poPrzebiegu, err := a.repozytorium.Sonda(ctx, sonda.Kod)
	if err != nil {
		return shared.HealthProbeRunResponse{}, bladMagazynuKondycji(err)
	}
	return shared.HealthProbeRunResponse{
		Result: wynikKontraktu(zapisany), Probe: sondaKontraktu(poPrzebiegu),
	}, nil
}

// WykazWynikow obsługuje `health.result.list`.
func (a *adapterKondycji) WykazWynikow(ctx context.Context,
	z shared.HealthResultListRequest) (shared.HealthResultListResponse, error) {

	if a.repozytorium == nil {
		return shared.HealthResultListResponse{}, bladZapleczaKondycji()
	}
	sito := dane.SitoWynikowKondycji{
		SondaKod: wartoscTekstu(z.ProbeId),
		OdCzasu:  wartoscLiczbyDlugiej(z.FromTime),
		DoCzasu:  wartoscLiczbyDlugiej(z.ToTime),
		Granica:  wartoscLiczby(z.Limit),
	}
	if z.Status != nil {
		sito.Stan = string(*z.Status)
	}
	wyniki, wszystkich, err := a.repozytorium.Wyniki(ctx, sito)
	if err != nil {
		return shared.HealthResultListResponse{}, bladMagazynuKondycji(err)
	}
	wykaz := make([]shared.HealthProbeResult, 0, len(wyniki))
	for _, wynik := range wyniki {
		wykaz = append(wykaz, wynikKontraktu(wynik))
	}
	przyciety := wszystkich > len(wykaz)
	return shared.HealthResultListResponse{
		Results: wykaz, Total: &wszystkich, Truncated: &przyciety,
	}, nil
}

// Dostepnosc obsługuje `health.uptime.get`.
//
// Dostępność liczy się z serii pomiarów, po jednej sondzie naraz. Wynik
// `degraded` liczy się jako połowa udanego: sonda, która odpowiedziała za
// wolno, odpowiedziała — ale nie tak, żeby nazwać to pełną dostępnością.
func (a *adapterKondycji) Dostepnosc(ctx context.Context,
	z shared.HealthUptimeGetRequest) (shared.HealthUptimeGetResponse, error) {

	if a.repozytorium == nil {
		return shared.HealthUptimeGetResponse{}, bladZapleczaKondycji()
	}
	doCzasu := wartoscLiczbyDlugiej(z.ToTime)
	if doCzasu <= 0 {
		doCzasu = terazWMilisekundachKondycji()
	}
	odCzasu := wartoscLiczbyDlugiej(z.FromTime)
	if odCzasu <= 0 {
		odCzasu = doCzasu - domyslnyZakresDostepnosci.Milliseconds()
	}

	sondy, err := a.sondyDoLiczenia(ctx, wartoscTekstu(z.ProbeId))
	if err != nil {
		return shared.HealthUptimeGetResponse{}, err
	}

	dostepnosci := make([]shared.HealthUptime, 0, len(sondy))
	for _, sonda := range sondy {
		wyniki, _, err := a.repozytorium.Wyniki(ctx, dane.SitoWynikowKondycji{
			SondaKod: sonda.Kod, OdCzasu: odCzasu, DoCzasu: doCzasu,
		})
		if err != nil {
			return shared.HealthUptimeGetResponse{}, bladMagazynuKondycji(err)
		}
		dostepnosci = append(dostepnosci,
			policzDostepnosc(sonda, wyniki, odCzasu, doCzasu, z.Objective))
	}
	return shared.HealthUptimeGetResponse{
		Uptimes: dostepnosci, FromTime: odCzasu, ToTime: doCzasu,
	}, nil
}

// sondyDoLiczenia zwraca sondy objęte wyliczeniem dostępności: jedną wskazaną
// albo wszystkie.
func (a *adapterKondycji) sondyDoLiczenia(ctx context.Context, kod string) ([]dane.SondaKondycji, error) {
	if strings.TrimSpace(kod) != "" {
		sonda, err := a.repozytorium.Sonda(ctx, strings.TrimSpace(kod))
		if err != nil {
			if errors.Is(err, dane.ErrBrakWiersza) {
				return nil, bladNieznanejSondy(kod)
			}
			return nil, bladMagazynuKondycji(err)
		}
		return []dane.SondaKondycji{sonda}, nil
	}
	sondy, err := a.repozytorium.Sondy(ctx, "", "", false, 0)
	if err != nil {
		return nil, bladMagazynuKondycji(err)
	}
	return sondy, nil
}

// policzDostepnosc liczy dostępność i budżet błędów jednej sondy z serii jej
// pomiarów.
//
// Seria pusta daje dostępność pustą, a nie sto procent: brak pomiarów nie jest
// dowodem sprawności. To jest właśnie ta różnica, przez którą awaria przechodzi
// niezauważona — sonda, która nigdy nie ruszyła, pokazałaby wtedy „100%".
func policzDostepnosc(sonda dane.SondaKondycji, wyniki []dane.WynikSondyKondycji,
	odCzasu, doCzasu int64, celZadania *float64) shared.HealthUptime {

	nazwa := sonda.Nazwa
	dostepnosc := shared.HealthUptime{
		ProbeId: sonda.Kod, ProbeName: &nazwa,
		FromTime: odCzasu, ToTime: doCzasu, Samples: len(wyniki),
	}
	cel := domyslnyCelDostepnosci
	switch {
	case celZadania != nil && *celZadania > 0:
		cel = *celZadania
	case sonda.CelDostepnosci != nil && *sonda.CelDostepnosci > 0:
		cel = *sonda.CelDostepnosci
	}
	dostepnosc.Objective = &cel

	if len(wyniki) == 0 {
		return dostepnosc
	}

	udane, oslabione, nieudane, zdarzenia := 0, 0, 0, 0
	var ostatniUpadek *int64
	poprzedniUpadek := false
	// Seria przychodzi od najnowszego; incydentem jest przejście ze stanu
	// sprawnego w niesprawny, więc idziemy od najstarszego.
	for numer := len(wyniki) - 1; numer >= 0; numer-- {
		wynik := wyniki[numer]
		switch wynik.Stan {
		case shared.HealthProbeStatusUp:
			udane++
			poprzedniUpadek = false
		case shared.HealthProbeStatusDegraded:
			oslabione++
			poprzedniUpadek = false
		default:
			nieudane++
			chwila := wynik.Wykonano
			ostatniUpadek = &chwila
			if !poprzedniUpadek {
				zdarzenia++
			}
			poprzedniUpadek = true
		}
	}

	dostepnosc.UpSamples = udane
	dostepnosc.DegradedSamples = &oslabione
	dostepnosc.DownSamples = &nieudane
	dostepnosc.LastDownAt = ostatniUpadek
	dostepnosc.Incidents = &zdarzenia

	procent := (float64(udane) + float64(oslabione)/2) * 100 / float64(len(wyniki))
	dostepnosc.UptimePercent = &procent

	// Budżet błędów: ile z dopuszczalnej niedostępności jeszcze zostało.
	// Cel stu procent nie zostawia budżetu, więc dzielenie przez zero zastępuje
	// jawna zerowa reszta.
	budzet := 0.0
	if dopuszczalne := 100 - cel; dopuszczalne > 0 {
		budzet = (dopuszczalne - (100 - procent)) * 100 / dopuszczalne
		if budzet < 0 {
			budzet = 0
		}
	}
	dostepnosc.ErrorBudgetPercent = &budzet
	return dostepnosc
}

// sondaKontraktu przekłada wiersz na sondę kontraktu.
func sondaKontraktu(s dane.SondaKondycji) shared.HealthProbe {
	sonda := shared.HealthProbe{
		Id: s.Kod, Name: s.Nazwa, Kind: shared.HealthProbeKind(s.Rodzaj), Target: s.Cel,
		ComponentId: s.KomponentKod, IntervalMs: s.OdstepMs, TimeoutMs: s.LimitCzasuMs,
		Payload: s.TrescWysylana, Enabled: s.Czynna, Objective: s.CelDostepnosci,
		CreatedAt: s.Utworzono, UpdatedAt: s.Zaktualizowano, LastRunAt: s.OstatniPrzebieg,
	}
	if s.OczekiwanyKod != nil {
		kod := int(*s.OczekiwanyKod)
		sonda.ExpectedStatus = &kod
	}
	if s.OstatniStan != nil {
		stan := shared.HealthProbeStatus(*s.OstatniStan)
		sonda.LastStatus = &stan
	}
	return sonda
}

// wynikKontraktu przekłada wiersz serii na wynik kontraktu.
func wynikKontraktu(w dane.WynikSondyKondycji) shared.HealthProbeResult {
	wynik := shared.HealthProbeResult{
		Id: w.Kod, ProbeId: w.SondaKod, Status: shared.HealthProbeStatus(w.Stan),
		RanAt: w.Wykonano, Detail: w.Szczegol, ErrorId: w.BladKod,
	}
	if w.CzasOdpowiedziMs != nil {
		czas := int(*w.CzasOdpowiedziMs)
		wynik.LatencyMs = &czas
	}
	if w.StatusHttp != nil {
		status := int(*w.StatusHttp)
		wynik.HttpStatus = &status
	}
	return wynik
}

// sprawdzRodzajSondy odbija rodzaj spoza wyliczenia kontraktu. Rodzaj nieznany
// zapisany do bazy dałby sondę, której przebieg zawsze kończy się `unknown` —
// czyli sondę udającą, że mierzy.
func sprawdzRodzajSondy(rodzaj shared.HealthProbeKind) error {
	switch rodzaj {
	case shared.HealthProbeKindHttp, shared.HealthProbeKindTcp, shared.HealthProbeKindCommand,
		shared.HealthProbeKindModelCall, shared.HealthProbeKindInternal:
		return nil
	default:
		return bladWskazaniaKondycji("nie znam rodzaju sondy „" + string(rodzaj) +
			"” — rdzeń zna: http, tcp, command, modelCall, internal")
	}
}

// terazWMilisekundachKondycji oddaje bieżącą chwilę w milisekundach epoki.
func terazWMilisekundachKondycji() int64 {
	return time.Now().UnixMilli()
}

// liczbaCalkowitaZeWskaznikaKondycji rozszerza wskaźnik na `int` do `int64`
// kolumny.
func liczbaCalkowitaZeWskaznikaKondycji(wartosc *int) *int64 {
	if wartosc == nil {
		return nil
	}
	rozszerzona := int64(*wartosc)
	return &rozszerzona
}

// wartoscLiczbyDlugiej oddaje wartość wskaźnika na `int64` albo zero.
func wartoscLiczbyDlugiej(wartosc *int64) int64 {
	if wartosc == nil {
		return 0
	}
	return *wartosc
}

// bladZapleczaKondycji nazywa brak magazynu sond po stronie rdzenia.
func bladZapleczaKondycji() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"kondycja: rdzeń nie ma wpiętego magazynu sond — naprawa: podpiąć "+
			"repozytorium kondycji przy składaniu rdzenia"))
}

// bladWskazaniaKondycji nazywa niepoprawne żądanie.
func bladWskazaniaKondycji(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "kondycja: "+powod))
}

// bladNieznanejSondy nazywa wskazanie sondy, której nie ma.
func bladNieznanejSondy(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"kondycja: nie ma sondy o identyfikatorze "+strings.TrimSpace(kod)))
}

// bladMagazynuKondycji nazywa niepowodzenie zapisu albo odczytu.
func bladMagazynuKondycji(err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"kondycja: "+err.Error()))
}
