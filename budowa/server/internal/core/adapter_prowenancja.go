// Odpowiedzialność pliku: adapter prowenancji wywołań modelu i rozliczenia
// zużycia — przekład między kontraktem a magazynem śladu.
//
// ── Czego ten adapter nie robi ───────────────────────────────────────────────
// Nie liczy kosztu z cennika. Koszt jest kolumną wiersza — zapisuje go ten, kto
// wywołanie wykonał i zna odpowiedź dostawcy. Wywołanie bez ceny nie jest
// wywołaniem darmowym, więc suma niesie osobno liczbę wierszy bez ceny: bez
// niej rozliczenie pokazywałoby kwotę mniejszą niż rzeczywista i wyglądałoby
// to identycznie jak kwota prawdziwa.
package core

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// granicaWykazuWywolan chroni odczyt przed wykazem wielkości całej tabeli.
// Wywołań przybywa w tempie pracy Operatora, a Provenance Explorer pokazuje
// stronę, nie całość.
const granicaWykazuWywolan = 500

type adapterProwenancji struct {
	repozytorium dane.RepozytoriumProwenancji
	// kanaly obsluguje wylacznie `provenance.call.replay`
	// (`adapter_prowenancja_powtorzenie.go`): powtorzenie jest NOWYM wywolaniem
	// kanalu modelu, a nie odczytem sladu.
	kanaly *models.Rejestr
}

func nowyAdapterProwenancji(r dane.RepozytoriumProwenancji) *adapterProwenancji {
	return &adapterProwenancji{repozytorium: r}
}

// odmowaSladu buduje odmowę kontraktu dla rodzin śladu i rozliczenia.
//
// Nazwa jest inna niż `bladProwenancji` z adaptera konfiguracji, bo tamten
// dotyczy prowenancji WYWOŁANIA MODELU w oknie konfiguracji — innego pojęcia
// o tej samej nazwie własnej. Dwie funkcje o jednej nazwie w jednym pakiecie
// byłyby dwiema prawdami o odmowie.
func odmowaSladu(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "ślad wywołań: "+powod))
}

// odmowaMagazynuSladu przekłada odmowę warstwy danych na odmowę kontraktu.
func odmowaMagazynuSladu(err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return odmowaSladu(shared.ErrorCodeNotFound,
			"wywołania o tym identyfikatorze nie ma w śladzie rdzenia")
	}
	return odmowaSladu(shared.ErrorCodeInternalError, err.Error())
}

// WykazWywolan oddaje ślad zawężony pytaniem Operatora.
//
// Pole `truncated` mówi prawdę o obcięciu: wykaz krótszy od liczby wszystkich
// pasujących znaczy, że Operator patrzy na wycinek. Bez tego pola strona
// pierwsza wyglądałaby jak komplet.
func (a *adapterProwenancji) WykazWywolan(ctx context.Context,
	z shared.ProvenanceCallListRequest) (shared.ProvenanceCallListResponse, error) {

	granica := granicaWykazuWywolan
	if z.Limit != nil && *z.Limit > 0 && *z.Limit < granicaWykazuWywolan {
		granica = *z.Limit
	}
	sito := dane.SitoWywolan{
		Od: z.FromTime, Do: z.ToTime,
		SesjaKod:  wartoscTekstu(z.SessionId),
		OknoKod:   wartoscTekstu(z.WindowId),
		ProcesKod: wartoscTekstu(z.ProcessId),
		SladKod:   wartoscTekstu(z.TraceId),
		KanalKod:  wartoscTekstu(z.ChannelId),
		Model:     wartoscTekstu(z.Model),
		KontoKod:  wartoscTekstu(z.AccountId),
		Fraza:     wartoscTekstu(z.Pattern),
		MinKoszt:  z.MinCost,
		Granica:   granica,
	}
	if z.Status != nil {
		sito.Stan = stanBazyWywolania(*z.Status)
	}

	wiersze, wszystkie, err := a.repozytorium.Wywolania(ctx, sito)
	if err != nil {
		return shared.ProvenanceCallListResponse{}, odmowaMagazynuSladu(err)
	}

	slady := make([]shared.ModelCallTrace, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if z.SlowOnly != nil && *z.SlowOnly && !wolneWywolanie(wiersz) {
			continue
		}
		slady = append(slady, zlozSladWywolania(wiersz))
	}
	obciete := len(slady) < wszystkie
	return shared.ProvenanceCallListResponse{
		Calls: slady, Total: &wszystkie, Truncated: &obciete,
	}, nil
}

// Wywolanie oddaje jeden ślad wraz z drzewem odcinków.
//
// Treść promptu i odpowiedzi idzie wyłącznie na wyraźne żądanie i wyłącznie
// wtedy, gdy została zapisana. Pole `redacted` odróżnia treść zredagowaną od
// niepodanej — bez niego Operator nie wie, czy patrzy na wywołanie bez treści,
// czy na treść, której mu nie pokazano.
func (a *adapterProwenancji) Wywolanie(ctx context.Context,
	z shared.ProvenanceCallGetRequest) (shared.ProvenanceCallGetResponse, error) {

	if strings.TrimSpace(z.CallId) == "" {
		return shared.ProvenanceCallGetResponse{}, odmowaSladu(
			shared.ErrorCodeValidationFailed, "odczyt śladu bez wskazania wywołania")
	}
	wiersz, err := a.repozytorium.Wywolanie(ctx, z.CallId)
	if err != nil {
		return shared.ProvenanceCallGetResponse{}, odmowaMagazynuSladu(err)
	}
	odcinki, err := a.repozytorium.Odcinki(ctx, z.CallId)
	if err != nil {
		return shared.ProvenanceCallGetResponse{}, odmowaMagazynuSladu(err)
	}

	odpowiedz := shared.ProvenanceCallGetResponse{
		Call:     zlozSladWywolania(wiersz),
		Spans:    zlozOdcinki(odcinki),
		Redacted: wiersz.Zredagowane,
	}
	if z.IncludeContent != nil && *z.IncludeContent && wiersz.TrescZapisana {
		odpowiedz.Prompt = wiersz.Prompt
		odpowiedz.Response = wiersz.Odpowiedz
	}
	return odpowiedz, nil
}

// OcenWywolanie zapisuje zdanie Operatora o wywołaniu.
func (a *adapterProwenancji) OcenWywolanie(ctx context.Context,
	z shared.ProvenanceCallRateRequest) (shared.ProvenanceCallRateResponse, error) {

	if strings.TrimSpace(z.CallId) == "" {
		return shared.ProvenanceCallRateResponse{}, odmowaSladu(
			shared.ErrorCodeValidationFailed, "ocena bez wskazania wywołania")
	}
	wiersz, err := a.repozytorium.OcenWywolanie(ctx, z.CallId,
		ocenaBazy(z.Quality), z.Note)
	if err != nil {
		return shared.ProvenanceCallRateResponse{}, odmowaMagazynuSladu(err)
	}
	return shared.ProvenanceCallRateResponse{Call: zlozSladWywolania(wiersz)}, nil
}

// WydajSlad składa wybrane wywołania w postać do wyniesienia poza produkt.
func (a *adapterProwenancji) WydajSlad(ctx context.Context,
	z shared.ProvenanceTraceExportRequest) (shared.ProvenanceTraceExportResponse, error) {

	sito := dane.SitoWywolan{Od: z.FromTime, Do: z.ToTime, Granica: granicaWykazuWywolan}
	wiersze, _, err := a.repozytorium.Wywolania(ctx, sito)
	if err != nil {
		return shared.ProvenanceTraceExportResponse{}, odmowaMagazynuSladu(err)
	}
	if len(z.CallIds) > 0 {
		wybrane := make(map[string]bool, len(z.CallIds))
		for _, kod := range z.CallIds {
			wybrane[kod] = true
		}
		zawezone := wiersze[:0]
		for _, wiersz := range wiersze {
			if wybrane[wiersz.Kod] {
				zawezone = append(zawezone, wiersz)
			}
		}
		wiersze = zawezone
	}

	zTrescia := z.IncludeContent != nil && *z.IncludeContent
	tresc, err := wydajSladWPostaci(z.Format, wiersze, zTrescia)
	if err != nil {
		return shared.ProvenanceTraceExportResponse{}, err
	}
	zredagowane := !zTrescia
	return shared.ProvenanceTraceExportResponse{
		Content: tresc, Format: z.Format,
		CallCount: len(wiersze), Redacted: &zredagowane,
	}, nil
}

// wydajSladWPostaci składa treść wyniesienia. Postać nieznana jest odmową —
// wyniesienie w innej postaci niż zamówiona wygląda identycznie jak zamówiona.
func wydajSladWPostaci(postac shared.TelemetryFormat,
	wiersze []dane.WywolanieModelu, zTrescia bool) (string, error) {

	switch postac {
	case shared.TelemetryFormatJson:
		bajty, err := json.MarshalIndent(zlozSlady(wiersze, zTrescia), "", "  ")
		if err != nil {
			return "", odmowaSladu(shared.ErrorCodeInternalError, err.Error())
		}
		return string(bajty), nil

	case shared.TelemetryFormatJsonl:
		var budowniczy strings.Builder
		for _, slad := range zlozSlady(wiersze, zTrescia) {
			bajty, err := json.Marshal(slad)
			if err != nil {
				return "", odmowaSladu(shared.ErrorCodeInternalError, err.Error())
			}
			budowniczy.Write(bajty)
			budowniczy.WriteByte('\n')
		}
		return budowniczy.String(), nil

	case shared.TelemetryFormatCsv:
		var budowniczy strings.Builder
		zapis := csv.NewWriter(&budowniczy)
		_ = zapis.Write([]string{"id", "traceId", "channelId", "model", "status",
			"startedAt", "latencyMs", "totalTokens", "cost"})
		for _, wiersz := range wiersze {
			_ = zapis.Write([]string{
				wiersz.Kod, wartoscTekstu(wiersz.SladKod), wartoscTekstu(wiersz.KanalKod),
				wartoscTekstu(wiersz.Model), wiersz.Stan,
				strconv.FormatInt(wiersz.Poczatek, 10),
				liczbaAlboPusto(wiersz.OpoznienieMs), liczbaAlboPusto(wiersz.TokenyRazem),
				kwotaAlboPusto(wiersz.Koszt),
			})
		}
		zapis.Flush()
		return budowniczy.String(), zapis.Error()

	case shared.TelemetryFormatOtlp:
		// Postać OTLP jest w kontrakcie, a rdzeń nie ma czym jej złożyć: wymaga
		// odwzorowania śladu na schemat OpenTelemetry, którego w drzewie nie ma.
		// Odmowa nazywa brak zamiast oddać JSON pod cudzą nazwą.
		return "", odmowaSladu(shared.ErrorCodeNotFound,
			"postać OTLP wymaga odwzorowania na schemat OpenTelemetry, którego rdzeń "+
				"nie niesie; wyniesienie w postaci json, jsonl albo csv jest dostępne")
	}
	return "", odmowaSladu(shared.ErrorCodeValidationFailed,
		"nieznana postać wyniesienia: "+string(postac))
}

// ── Rozliczenie zużycia ─────────────────────────────────────────────────────

type adapterZuzycia struct {
	repozytorium dane.RepozytoriumProwenancji
}

func nowyAdapterZuzycia(r dane.RepozytoriumProwenancji) *adapterZuzycia {
	return &adapterZuzycia{repozytorium: r}
}

// Podsumowanie liczy sumy po wskazanym wymiarze.
//
// `priceCoverage` mówi, jaka część wierszy niosła cenę. Suma kosztu bez tej
// liczby jest myląca: wygląda tak samo przy pełnym cenniku i przy jednym
// wierszu wycenionym na dziesięć.
func (a *adapterZuzycia) Podsumowanie(ctx context.Context,
	z shared.UsageSummaryGetRequest) (shared.UsageSummaryGetResponse, error) {

	granica := 50
	if z.Limit != nil && *z.Limit > 0 {
		granica = *z.Limit
	}
	sumy, err := a.repozytorium.Zuzycie(ctx, string(z.Dimension), z.FromTime, z.ToTime, granica)
	if err != nil {
		return shared.UsageSummaryGetResponse{}, odmowaSladu(
			shared.ErrorCodeValidationFailed, err.Error())
	}
	if len(z.DimensionIds) > 0 {
		wybrane := make(map[string]bool, len(z.DimensionIds))
		for _, klucz := range z.DimensionIds {
			wybrane[klucz] = true
		}
		zawezone := sumy[:0]
		for _, suma := range sumy {
			if wybrane[suma.Klucz] {
				zawezone = append(zawezone, suma)
			}
		}
		sumy = zawezone
	}

	wykaz := make([]shared.UsageAggregate, 0, len(sumy))
	razem := shared.UsageAggregate{
		Dimension: z.Dimension, DimensionId: "*",
	}
	zCena, wszystkie := 0, 0
	for _, suma := range sumy {
		wykaz = append(wykaz, zlozSumeZuzycia(suma, z.Dimension))
		razem.Requests += suma.Zadania
		razem.PromptTokens += suma.TokenyPromptu
		razem.CompletionTokens += suma.TokenyOdpowiedzi
		razem.TotalTokens += suma.TokenyRazem
		wszystkie += suma.Zadania
		zCena += suma.Zadania - suma.BezCeny
	}
	pokrycie := 0.0
	if wszystkie > 0 {
		pokrycie = float64(zCena) / float64(wszystkie)
	}
	return shared.UsageSummaryGetResponse{
		Aggregates: wykaz, Overall: &razem,
		FromTime: chwilaAlboZero(z.FromTime), ToTime: chwilaAlboZero(z.ToTime),
		PriceCoverage: &pokrycie,
	}, nil
}

// Raport składa rozliczenie po wielu wymiarach naraz.
func (a *adapterZuzycia) Raport(ctx context.Context,
	z shared.UsageReportBuildRequest) (shared.UsageReportBuildResponse, error) {

	wymiary := z.Dimensions
	if len(wymiary) == 0 {
		wymiary = []shared.UsageDimension{shared.UsageDimensionChannel}
	}
	raport := map[string][]shared.UsageAggregate{}
	zadania := 0
	for _, wymiar := range wymiary {
		sumy, err := a.repozytorium.Zuzycie(ctx, string(wymiar), z.FromTime, z.ToTime, 100)
		if err != nil {
			return shared.UsageReportBuildResponse{}, odmowaSladu(
				shared.ErrorCodeValidationFailed, err.Error())
		}
		wykaz := make([]shared.UsageAggregate, 0, len(sumy))
		for _, suma := range sumy {
			wykaz = append(wykaz, zlozSumeZuzycia(suma, wymiar))
			zadania += suma.Zadania
		}
		raport[string(wymiar)] = wykaz
	}

	tresc, err := raportWPostaci(z.Format, raport)
	if err != nil {
		return shared.UsageReportBuildResponse{}, err
	}
	return shared.UsageReportBuildResponse{
		Content: tresc, Format: z.Format,
		GeneratedAt: protocol.Teraz(), Requests: zadania,
	}, nil
}

func raportWPostaci(postac shared.TelemetryFormat,
	raport map[string][]shared.UsageAggregate) (string, error) {

	switch postac {
	case shared.TelemetryFormatJson, shared.TelemetryFormatJsonl:
		bajty, err := json.MarshalIndent(raport, "", "  ")
		if err != nil {
			return "", odmowaSladu(shared.ErrorCodeInternalError, err.Error())
		}
		return string(bajty), nil
	case shared.TelemetryFormatCsv:
		var budowniczy strings.Builder
		zapis := csv.NewWriter(&budowniczy)
		_ = zapis.Write([]string{"wymiar", "klucz", "zadania", "tokeny", "koszt"})
		wymiary := make([]string, 0, len(raport))
		for wymiar := range raport {
			wymiary = append(wymiary, wymiar)
		}
		sort.Strings(wymiary)
		for _, wymiar := range wymiary {
			for _, suma := range raport[wymiar] {
				_ = zapis.Write([]string{wymiar, suma.DimensionId,
					strconv.Itoa(suma.Requests), strconv.Itoa(suma.TotalTokens),
					kwotaAlboPusto(suma.Cost)})
			}
		}
		zapis.Flush()
		return budowniczy.String(), zapis.Error()
	}
	return "", odmowaSladu(shared.ErrorCodeNotFound,
		"postać "+string(postac)+" nie ma w rdzeniu odwzorowania dla raportu rozliczenia")
}

// ── Przekład wiersz ↔ kontrakt ──────────────────────────────────────────────

// stanBazyWywolania i stanKontraktuWywolania przekładają stan w obie strony.
// Odwzorowanie stoi w kontrakcie (pole `baza` przy wartości wyliczenia), a te
// dwie funkcje są jego jedynym miejscem użycia w rdzeniu.
var stanyWywolania = map[shared.ModelCallStatus]string{
	shared.ModelCallStatusRunning:   "biegnie",
	shared.ModelCallStatusOk:        "zakonczone",
	shared.ModelCallStatusFailed:    "bledne",
	shared.ModelCallStatusTimeout:   "limit_czasu",
	shared.ModelCallStatusCancelled: "przerwane",
}

func stanBazyWywolania(stan shared.ModelCallStatus) string {
	if wartosc, znany := stanyWywolania[stan]; znany {
		return wartosc
	}
	return ""
}

func stanKontraktuWywolania(stan string) shared.ModelCallStatus {
	for kontrakt, baza := range stanyWywolania {
		if baza == stan {
			return kontrakt
		}
	}
	return shared.ModelCallStatusRunning
}

var ocenyWywolania = map[shared.ModelCallQuality]string{
	shared.ModelCallQualityUnrated:    "bez_oceny",
	shared.ModelCallQualityAccurate:   "trafna",
	shared.ModelCallQualityPartial:    "czesciowa",
	shared.ModelCallQualityInaccurate: "nietrafna",
}

func ocenaBazy(ocena shared.ModelCallQuality) string {
	if wartosc, znana := ocenyWywolania[ocena]; znana {
		return wartosc
	}
	return "bez_oceny"
}

func ocenaKontraktu(ocena string) shared.ModelCallQuality {
	for kontrakt, baza := range ocenyWywolania {
		if baza == ocena {
			return kontrakt
		}
	}
	return shared.ModelCallQualityUnrated
}

var rodzajeOdcinka = map[string]shared.ModelCallSpanKind{
	"prompt":    shared.ModelCallSpanKindPrompt,
	"odpowiedz": shared.ModelCallSpanKindCompletion,
	"narzedzie": shared.ModelCallSpanKindTool,
	"podagent":  shared.ModelCallSpanKindSubagent,
	"kontekst":  shared.ModelCallSpanKindRetrieval,
	"cache":     shared.ModelCallSpanKindCache,
}

// progWolnegoWywolania oddziela wywołanie wolne od zwykłego.
//
// Dziesięć sekund jest progiem rdzenia, nie prawdą o modelach: wywołanie
// dłuższe znaczy, że Operator czekał, patrząc na okno. Wartość stoi tutaj,
// żeby dało się ją zmienić w jednym miejscu, gdy pojawi się ustawienie.
const progWolnegoWywolania = 10_000

func wolneWywolanie(w dane.WywolanieModelu) bool {
	return w.OpoznienieMs != nil && *w.OpoznienieMs >= progWolnegoWywolania
}

func zlozSladWywolania(w dane.WywolanieModelu) shared.ModelCallTrace {
	wolne := wolneWywolanie(w)
	ocena := ocenaKontraktu(w.Ocena)
	return shared.ModelCallTrace{
		Id: w.Kod, TraceId: w.SladKod, ParentCallId: w.RodzicKod,
		ProcessId: w.ProcesKod, SessionId: w.SesjaKod, WindowId: w.OknoKod,
		MessageId: w.WiadomoscKod, ChannelId: w.KanalKod, Provider: w.Dostawca,
		Model: w.Model, AccountId: w.KontoKod, ProjectId: w.ProjektKod,
		Environment: w.Srodowisko,
		Status:      stanKontraktuWywolania(w.Stan), ErrorCode: kodOdmowy(w.KodBledu),
		StartedAt: w.Poczatek, FinishedAt: w.Koniec, LatencyMs: w.OpoznienieMs,
		PromptTokens: w.TokenyPromptu, CompletionTokens: w.TokenyOdpowiedzi,
		CachedTokens: w.TokenyCache, TotalTokens: w.TokenyRazem,
		Cost: w.Koszt, Currency: w.Waluta,
		ToolCallCount: &w.LiczbaNarzedzi, SubagentCallCount: &w.LiczbaPodagentow,
		Slow: &wolne, Quality: &ocena, QualityNote: w.OcenaNotatka,
		ContentStored: w.TrescZapisana, Redacted: &w.Zredagowane,
	}
}

func zlozSlady(wiersze []dane.WywolanieModelu, zTrescia bool) []shared.ModelCallTrace {
	slady := make([]shared.ModelCallTrace, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if !zTrescia {
			wiersz.Prompt, wiersz.Odpowiedz = nil, nil
		}
		slady = append(slady, zlozSladWywolania(wiersz))
	}
	return slady
}

func zlozOdcinki(odcinki []dane.OdcinekWywolania) []shared.ModelCallSpan {
	wykaz := make([]shared.ModelCallSpan, 0, len(odcinki))
	for _, o := range odcinki {
		rodzaj, znany := rodzajeOdcinka[o.Rodzaj]
		if !znany {
			rodzaj = shared.ModelCallSpanKindTool
		}
		trafienie := o.TrafienieCache
		wykaz = append(wykaz, shared.ModelCallSpan{
			Id: o.Kod, CallId: o.WywolanieKod, ParentSpanId: o.RodzicKod,
			Name: o.Nazwa, Kind: rodzaj,
			StartedAt: o.Poczatek, FinishedAt: o.Koniec, DurationMs: o.CzasMs,
			Tokens: o.Tokeny, Cost: o.Koszt,
			Status: stanKontraktuWywolania(o.Stan), ErrorCode: kodOdmowy(o.KodBledu),
			CacheHit: &trafienie,
		})
	}
	return wykaz
}

func zlozSumeZuzycia(s dane.SumaZuzycia, wymiar shared.UsageDimension) shared.UsageAggregate {
	bledne, bezCeny, opoznienie := s.ZadaniaBledne, s.BezCeny, s.SrednieOpoznienie
	koszt := s.Koszt
	klucz := s.Klucz
	if klucz == "" {
		klucz = "(bez wskazania)"
	}
	return shared.UsageAggregate{
		Dimension: wymiar, DimensionId: klucz,
		Requests: s.Zadania, FailedRequests: &bledne,
		PromptTokens: s.TokenyPromptu, CompletionTokens: s.TokenyOdpowiedzi,
		CachedTokens: &s.TokenyCache, TotalTokens: s.TokenyRazem,
		Cost: &koszt, CostWithoutPrice: &bezCeny, AvgLatencyMs: &opoznienie,
	}
}

// chwilaAlboZero oddaje wskazaną granicę okna czasu albo zero.
//
// Zero znaczy „bez granicy" i jest tą samą wartością, którą przyjęło zapytanie —
// odpowiedź powtarza granice, na których naprawdę liczyła, a nie te, o które
// pytano. Przy braku granicy obie liczby są zerami i to jest prawda o wyniku.
// kodOdmowy przekłada kolumnę tekstową na kod odmowy kontraktu.
//
// Kolumna trzyma napis, bo baza nie zna typu kontraktu; przekład stoi tutaj,
// żeby nie powstał drugi w każdym miejscu odczytu.
func kodOdmowy(wartosc *string) *shared.ErrorCode {
	if wartosc == nil || *wartosc == "" {
		return nil
	}
	kod := shared.ErrorCode(*wartosc)
	return &kod
}

func chwilaAlboZero(wskazanie *int64) int64 {
	if wskazanie == nil {
		return 0
	}
	return *wskazanie
}

func liczbaAlboPusto(wartosc *int) string {
	if wartosc == nil {
		return ""
	}
	return strconv.Itoa(*wartosc)
}

func kwotaAlboPusto(wartosc *float64) string {
	if wartosc == nil {
		return ""
	}
	return fmt.Sprintf("%.6f", *wartosc)
}
