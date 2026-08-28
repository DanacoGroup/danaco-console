// Plik obsługuje przestrzeń badania: pytania badawcze i ich pokrycie, notatkę
// roboczą, świeżość źródeł, luki badawcze, liczniki przesiewu PRISMA oraz graf
// dowodów łączący źródła, ustalenia i sprzeczności.
package core

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// progSwiezosciBadania jest liczbą dni, po których badanie uchodzi za wymagające
// odświeżenia źródeł — trzydzieści dni starcza, żeby literatura zdążyła się
// zmienić.
const progSwiezosciBadania = 30

// PobierzPrzestrzen obsługuje komendę research.workspace.get: zwraca zakres,
// etapy, pytania badawcze, odbiorcę, protokół i notatkę roboczą jednej
// przestrzeni badania.
func (a *adapterBadan) PobierzPrzestrzen(ctx context.Context,
	_ shared.ResearchWorkspaceGetRequest) (shared.ResearchWorkspaceGetResponse, error) {

	zakres, etapy, err := a.repozytorium.Przestrzen(ctx)
	if err != nil {
		return shared.ResearchWorkspaceGetResponse{}, bladBadan(err)
	}
	szczegoly, err := a.repozytorium.SzczegolyPrzestrzeniBadania(ctx)
	if err != nil {
		return shared.ResearchWorkspaceGetResponse{}, bladBadan(err)
	}
	pytania, err := a.repozytorium.PytaniaBadania(ctx)
	if err != nil {
		return shared.ResearchWorkspaceGetResponse{}, bladBadan(err)
	}
	if etapy == nil {
		etapy = []string{}
	}
	return shared.ResearchWorkspaceGetResponse{
		Scope: zakres, Stages: etapy, Questions: zlozPytaniaBadania(pytania, nil),
		Audience: szczegoly.Odbiorca, Protocol: szczegoly.Protokol, Note: szczegoly.Notatka,
		UpdatedAt: time.Now().UnixMilli(),
	}, nil
}

// UstawPytaniaBadania obsługuje komendę research.workspace.question.set:
// zapisuje wykaz pytań badawczych, nadając nowy identyfikator każdemu pytaniu
// bez własnego kodu.
func (a *adapterBadan) UstawPytaniaBadania(ctx context.Context,
	z shared.ResearchWorkspaceQuestionSetRequest) (shared.ResearchWorkspaceQuestionSetResponse, error) {

	doZapisu := make([]dane.PytanieBadania, 0, len(z.Questions))
	for numer, pytanie := range z.Questions {
		if strings.TrimSpace(pytanie.Text) == "" {
			return shared.ResearchWorkspaceQuestionSetResponse{},
				bladWskazaniaBadan("pytanie badawcze bez treści — pusty wiersz nie jest pytaniem")
		}
		kod := pytanie.Id
		if strings.TrimSpace(kod) == "" {
			kod = nowyIdentyfikator(przedrostekPytaniaBadania)
		}
		kolejnosc := numer + 1
		if pytanie.Order != nil {
			kolejnosc = *pytanie.Order
		}
		doZapisu = append(doZapisu, dane.PytanieBadania{
			Kod: kod, Tekst: pytanie.Text, Kolejnosc: kolejnosc,
		})
	}
	zapisane, err := a.repozytorium.UstawPytaniaBadania(ctx, doZapisu)
	if err != nil {
		return shared.ResearchWorkspaceQuestionSetResponse{}, bladBadan(err)
	}
	return shared.ResearchWorkspaceQuestionSetResponse{
		Questions: zlozPytaniaBadania(zapisane, nil),
	}, nil
}

// PokryciePytan obsługuje komendę research.workspace.coverage: zwraca pytania
// badawcze wraz z kodami źródeł, które je pokrywają, oraz liczbę pytań bez ani
// jednego źródła.
func (a *adapterBadan) PokryciePytan(ctx context.Context,
	z shared.ResearchWorkspaceCoverageRequest) (shared.ResearchWorkspaceCoverageResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchWorkspaceCoverageResponse{},
			bladWskazaniaBadan("workspace.coverage bez okna badania")
	}
	pytania, err := a.repozytorium.PytaniaBadania(ctx)
	if err != nil {
		return shared.ResearchWorkspaceCoverageResponse{}, bladBadan(err)
	}
	pokrycie, err := a.pokrycieOknaBadania(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchWorkspaceCoverageResponse{}, err
	}
	przelozone := zlozPytaniaBadania(pytania, pokrycie)
	niepokryte := 0
	for _, pytanie := range przelozone {
		if len(pytanie.SourceIds) == 0 {
			niepokryte++
		}
	}
	return shared.ResearchWorkspaceCoverageResponse{
		Questions: przelozone, UncoveredCount: niepokryte,
	}, nil
}

// pokrycieOknaBadania odwzorowuje pytanie badawcze na kody źródeł, które je
// pokrywają — jedno przejście po źródłach okna zamiast zapytania na pytanie.
func (a *adapterBadan) pokrycieOknaBadania(ctx context.Context,
	okno string) (map[string][]string, error) {

	zrodla, err := a.repozytorium.Zrodla(ctx, okno)
	if err != nil {
		return nil, bladBadan(err)
	}
	pokrycie := map[string][]string{}
	for _, zrodlo := range zrodla {
		katalog, err := a.repozytorium.KatalogZrodlaBadania(ctx, zrodlo.Kod)
		if err != nil {
			return nil, bladBadan(err)
		}
		for _, pytanie := range katalog.Pytania {
			pokrycie[pytanie] = append(pokrycie[pytanie], zrodlo.Kod)
		}
	}
	return pokrycie, nil
}

// zlozPytaniaBadania przekłada wiersze pytań na byty kontraktu, dokładając
// pokrycie, gdy zostało policzone.
func zlozPytaniaBadania(pytania []dane.PytanieBadania,
	pokrycie map[string][]string) []shared.ResearchQuestion {

	przelozone := make([]shared.ResearchQuestion, 0, len(pytania))
	for _, pytanie := range pytania {
		kolejnosc := pytanie.Kolejnosc
		pozycja := shared.ResearchQuestion{
			Id: pytanie.Kod, Text: pytanie.Tekst, Order: &kolejnosc,
		}
		if pokrycie != nil {
			pozycja.SourceIds = posortowaneNapisyBadania(pokrycie[pytanie.Kod])
		}
		przelozone = append(przelozone, pozycja)
	}
	return przelozone
}

// UstawNotatkeBadania obsługuje komendę research.workspace.note.set: zapisuje
// treść notatki roboczej przestrzeni badania wraz ze znacznikiem czasu zapisu.
func (a *adapterBadan) UstawNotatkeBadania(ctx context.Context,
	z shared.ResearchWorkspaceNoteSetRequest) (shared.ResearchWorkspaceNoteSetResponse, error) {

	chwila, err := a.repozytorium.UstawNotatkePrzestrzeni(ctx, z.Content)
	if err != nil {
		return shared.ResearchWorkspaceNoteSetResponse{}, bladBadan(err)
	}
	return shared.ResearchWorkspaceNoteSetResponse{
		Content: z.Content, UpdatedAt: chwilaBazy(chwila),
	}, nil
}

// SwiezoscBadania obsługuje komendę research.workspace.freshness: porównuje
// datę najnowszego źródła okna z progiem świeżości i zwraca obie wartości
// wprost.
func (a *adapterBadan) SwiezoscBadania(ctx context.Context,
	z shared.ResearchWorkspaceFreshnessRequest) (shared.ResearchWorkspaceFreshnessResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchWorkspaceFreshnessResponse{},
			bladWskazaniaBadan("workspace.freshness bez okna badania")
	}
	zrodla, err := a.repozytorium.Zrodla(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchWorkspaceFreshnessResponse{}, bladBadan(err)
	}
	odpowiedz := shared.ResearchWorkspaceFreshnessResponse{StaleDays: progSwiezosciBadania}
	if len(zrodla) == 0 {
		// Badanie bez źródeł jest nierozpoczęte, nie nieświeże — sugestia
		// odświeżenia tu nie ma sensu.
		return odpowiedz, nil
	}
	najnowsze := int64(0)
	for _, zrodlo := range zrodla {
		if chwila := chwilaBazy(zrodlo.PozyskanoO); chwila > najnowsze {
			najnowsze = chwila
		}
	}
	if najnowsze > 0 {
		odpowiedz.NewestSourceAt = &najnowsze
		wiek := time.Since(time.UnixMilli(najnowsze))
		odpowiedz.RefreshSuggested = wiek > progSwiezosciBadania*24*time.Hour
	}
	return odpowiedz, nil
}

// SzukajLukBadania obsługuje `research.gap.find`. Luka jest pytaniem bez źródeł
// — to jest fakt policzony z bazy. Podpowiedź zapytania dokłada model, gdy jest
// dostępny; jego brak nie odbiera Operatorowi wykazu luk.
func (a *adapterBadan) SzukajLukBadania(ctx context.Context,
	z shared.ResearchGapFindRequest) (shared.ResearchGapFindResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchGapFindResponse{}, bladWskazaniaBadan("gap.find bez okna badania")
	}
	pytania, err := a.repozytorium.PytaniaBadania(ctx)
	if err != nil {
		return shared.ResearchGapFindResponse{}, bladBadan(err)
	}
	pokrycie, err := a.pokrycieOknaBadania(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchGapFindResponse{}, err
	}

	luki := []shared.ResearchGap{}
	for _, pytanie := range pytania {
		if len(pokrycie[pytanie.Kod]) > 0 {
			continue
		}
		kod := pytanie.Kod
		zapytanie := pytanie.Tekst
		luki = append(luki, shared.ResearchGap{
			Summary:        "pytanie badawcze bez ani jednego źródła: " + pytanie.Tekst,
			QuestionId:     &kod,
			SuggestedQuery: &zapytanie,
		})
	}
	if len(pytania) == 0 {
		luki = append(luki, shared.ResearchGap{
			Summary: "badanie nie ma zapisanych pytań badawczych — bez nich nie da się " +
				"zmierzyć pokrycia ani wskazać luk; naprawa: zapisać pytania przez " +
				"research.workspace.question.set",
		})
	}
	return shared.ResearchGapFindResponse{Gaps: luki}, nil
}

// PrzesiewPrisma obsługuje `research.prisma.get`. Liczniki idą z wierszy
// odkryć i źródeł, nie z pamięci sesji — diagram ma przeżyć zamknięcie okna.
func (a *adapterBadan) PrzesiewPrisma(ctx context.Context,
	z shared.ResearchPrismaGetRequest) (shared.ResearchPrismaGetResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchPrismaGetResponse{}, bladWskazaniaBadan("prisma.get bez okna badania")
	}
	liczniki, err := a.licznikiPrismaBadania(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchPrismaGetResponse{}, err
	}
	return shared.ResearchPrismaGetResponse{Counts: liczniki}, nil
}

// licznikiPrismaBadania liczy przesiew: zidentyfikowane, duplikaty, przesiane,
// wyłączone i włączone, wraz z wykazem powodów wyłączenia posortowanym
// alfabetycznie.
func (a *adapterBadan) licznikiPrismaBadania(ctx context.Context,
	okno string) (shared.ResearchPrismaCounts, error) {

	wyniki, err := a.repozytorium.WynikiOdkrycia(ctx, okno)
	if err != nil {
		return shared.ResearchPrismaCounts{}, bladBadan(err)
	}
	zrodla, err := a.repozytorium.Zrodla(ctx, okno)
	if err != nil {
		return shared.ResearchPrismaCounts{}, bladBadan(err)
	}

	liczniki := shared.ResearchPrismaCounts{Identified: len(wyniki), Included: len(zrodla)}
	powody := map[string]bool{}
	for _, wynik := range wyniki {
		if wynik.Duplikat {
			liczniki.DuplicatesRemoved++
		}
		if wynik.Odrzucony {
			liczniki.Excluded++
			if wynik.PowodOdrzucenia != nil && strings.TrimSpace(*wynik.PowodOdrzucenia) != "" {
				powody[*wynik.PowodOdrzucenia] = true
			}
		}
	}
	liczniki.Screened = liczniki.Identified - liczniki.DuplicatesRemoved
	for powod := range powody {
		liczniki.ExclusionReasons = append(liczniki.ExclusionReasons, powod)
	}
	sort.Strings(liczniki.ExclusionReasons)
	return liczniki, nil
}

// GrafDowodow obsługuje `research.evidence.graph`. Węzły to źródła, ustalenia
// i sprzeczności; krawędzie to poparcie i rozbieżność.
func (a *adapterBadan) GrafDowodow(ctx context.Context,
	z shared.ResearchEvidenceGraphRequest) (shared.ResearchEvidenceGraphResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchEvidenceGraphResponse{}, bladWskazaniaBadan("evidence.graph bez okna badania")
	}
	zrodla, err := a.repozytorium.Zrodla(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchEvidenceGraphResponse{}, bladBadan(err)
	}
	ustalenia, err := a.repozytorium.Ustalenia(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchEvidenceGraphResponse{}, bladBadan(err)
	}
	sprzecznosci, err := a.repozytorium.Sprzecznosci(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchEvidenceGraphResponse{}, bladBadan(err)
	}

	wezly := []shared.ResearchEvidenceNode{}
	krawedzie := []shared.ResearchEvidenceEdge{}
	for _, zrodlo := range zrodla {
		wezly = append(wezly, shared.ResearchEvidenceNode{
			Id: zrodlo.Kod, Kind: "zrodlo", Label: zrodlo.Tytul,
		})
	}
	for _, ustalenie := range ustalenia {
		wezly = append(wezly, shared.ResearchEvidenceNode{
			Id: ustalenie.Kod, Kind: "ustalenie",
			Label: skrocDoBadania(trescUstaleniaDoPromptu(ustalenie), 120),
		})
		powiazane, err := a.repozytorium.ZrodlaUstalenia(ctx, ustalenie.ID)
		if err != nil {
			return shared.ResearchEvidenceGraphResponse{}, bladBadan(err)
		}
		for _, zrodlo := range powiazane {
			krawedzie = append(krawedzie, shared.ResearchEvidenceEdge{
				FromId: zrodlo.Kod, ToId: ustalenie.Kod, Relation: "poparcie",
			})
		}
	}
	for _, sprzecznosc := range sprzecznosci {
		wezly = append(wezly, shared.ResearchEvidenceNode{
			Id: sprzecznosc.Kod, Kind: "sprzecznosc", Label: sprzecznosc.Streszczenie,
		})
		for _, kod := range sprzecznosc.UstalenieKody {
			krawedzie = append(krawedzie, shared.ResearchEvidenceEdge{
				FromId: sprzecznosc.Kod, ToId: kod, Relation: "rozbieznosc",
			})
		}
	}
	return shared.ResearchEvidenceGraphResponse{Nodes: wezly, Edges: krawedzie}, nil
}

// UstawPrzestrzenPelna rozszerza research.workspace.set o pola, których zapis
// zastany nie obejmował: pytania badawcze, odbiorcę, protokół i granice
// tematu — to ta sama komenda, nie druga jej odmiana.
func (a *adapterBadan) UstawPrzestrzenPelna(ctx context.Context,
	z shared.ResearchWorkspaceSetRequest) (shared.ResearchWorkspaceSetResponse, error) {

	odpowiedz, err := a.UstawPrzestrzen(ctx, z)
	if err != nil {
		return shared.ResearchWorkspaceSetResponse{}, err
	}
	if err := a.repozytorium.UstawSzczegolyPrzestrzeni(ctx, dane.SzczegolyPrzestrzeniBadania{
		Odbiorca: z.Audience, Protokol: z.Protocol, Granice: z.Boundaries,
	}); err != nil {
		return shared.ResearchWorkspaceSetResponse{}, bladBadan(err)
	}
	if z.Questions != nil {
		pytania, err := a.UstawPytaniaBadania(ctx,
			shared.ResearchWorkspaceQuestionSetRequest{Questions: z.Questions})
		if err != nil {
			return shared.ResearchWorkspaceSetResponse{}, err
		}
		odpowiedz.Questions = pytania.Questions
	} else {
		zapisane, err := a.repozytorium.PytaniaBadania(ctx)
		if err != nil {
			return shared.ResearchWorkspaceSetResponse{}, bladBadan(err)
		}
		odpowiedz.Questions = zlozPytaniaBadania(zapisane, nil)
	}
	odpowiedz.Audience = z.Audience
	odpowiedz.Protocol = z.Protocol
	return odpowiedz, nil
}

// opisLicznikaBadania składa czytelny opis licznika — używany przy składaniu
// diagramu przesiewu do treści raportu i eksportu.
func opisLicznikaBadania(nazwa string, wartosc int) string {
	return nazwa + ": " + strconv.Itoa(wartosc)
}
