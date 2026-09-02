// Adapter prowenancji wywołań modelu i rozliczenia zużycia. Koszt zapisuje ten,
// kto wywołanie wykonał i zna odpowiedź dostawcy, nie cennik tego adaptera.
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

// Provenance Explorer pokazuje stronę, nie całość tabeli.
const granicaWykazuWywolan = 500

type adapterProwenancji struct {
	repozytorium dane.RepozytoriumProwenancji
	kanaly       *models.Rejestr
	// repozytoriumKanalow rozstrzyga własność kanału z żądania (decyzja 34).
	repozytoriumKanalow dane.RepozytoriumKanalow
}

func nowyAdapterProwenancji(r dane.RepozytoriumProwenancji) *adapterProwenancji {
	return &adapterProwenancji{repozytorium: r}
}

// Nazwa inna niż `bladProwenancji` adaptera konfiguracji, bo tamten dotyczy innego pojęcia o tej samej nazwie.
func odmowaSladu(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "ślad wywołań: "+powod))
}

func odmowaMagazynuSladu(err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return odmowaSladu(shared.ErrorCodeNotFound,
			"wywołania o tym identyfikatorze nie ma w śladzie serwera")
	}
	return odmowaSladu(shared.ErrorCodeInternalError, err.Error())
}

// Pole `truncated` mówi prawdę o obcięciu; bez niego strona pierwsza wyglądałaby jak komplet.
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

// Treść idzie wyłącznie na żądanie i gdy została zapisana; `redacted` odróżnia treść zredagowaną od niepodanej.
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

// Postać nieznana jest odmową — wyniesienie w innej postaci wygląda identycznie jak zamówiona.
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
		// Rdzeń nie niesie odwzorowania na schemat OpenTelemetry: odmowa nazywa brak.
		return "", odmowaSladu(shared.ErrorCodeNotFound,
			"postać OTLP wymaga odwzorowania na schemat OpenTelemetry, którego serwer "+
				"nie niesie; wyniesienie w postaci json, jsonl albo csv jest dostępne")
	}
	return "", odmowaSladu(shared.ErrorCodeValidationFailed,
		"nieznana postać wyniesienia: "+string(postac))
}

type adapterZuzycia struct {
	repozytorium dane.RepozytoriumProwenancji
}

func nowyAdapterZuzycia(r dane.RepozytoriumProwenancji) *adapterZuzycia {
	return &adapterZuzycia{repozytorium: r}
}

// `priceCoverage` mówi, jaka część wierszy niosła cenę; suma kosztu bez tej liczby jest myląca.
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
		"postać "+string(postac)+" nie ma w serwerze odwzorowania dla raportu rozliczenia")
}

// Odwzorowanie stanu stoi w kontrakcie (pole `baza` przy wartości wyliczenia).
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

// Dziesięć sekund jest progiem rdzenia, nie prawdą o modelach; jedno miejsce do zmiany.
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

// Kolumna trzyma napis, bo baza nie zna typu kontraktu; przekład stoi w jednym miejscu.
func kodOdmowy(wartosc *string) *shared.ErrorCode {
	if wartosc == nil || *wartosc == "" {
		return nil
	}
	kod := shared.ErrorCode(*wartosc)
	return &kod
}

// Zero znaczy „bez granicy” i jest tą samą wartością, którą przyjęło zapytanie.
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
