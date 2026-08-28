// Pakiet obsługuje panel ustaleń badania: przegląd i zmianę ustaleń, kodowanie
// jakościowe wraz z książką kodów i macierzą kod na źródło, wykrywanie
// i rozstrzyganie sprzeczności, grupowanie w wątki, scalanie oraz ślad prowenancji.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// WypiszUstalenia obsługuje `research.finding.list`: wykaz ustaleń zapisanych
// dotąd przy wskazanym oknie badania.
func (a *adapterBadan) WypiszUstalenia(ctx context.Context,
	z shared.ResearchFindingListRequest) (shared.ResearchFindingListResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchFindingListResponse{}, bladWskazaniaBadan("finding.list bez okna badania")
	}
	ustalenia, err := a.repozytorium.Ustalenia(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchFindingListResponse{}, bladBadan(err)
	}

	dopasowane := []shared.ResearchFinding{}
	for _, ustalenie := range ustalenia {
		pasuje, przelozone, err := a.ustaleniePasujeBadania(ctx, ustalenie, z)
		if err != nil {
			return shared.ResearchFindingListResponse{}, err
		}
		if pasuje {
			dopasowane = append(dopasowane, przelozone)
		}
	}
	razem := len(dopasowane)
	wycinek := wycinekListyBadania(razem, z.Offset, z.Limit)
	return shared.ResearchFindingListResponse{
		Findings: dopasowane[wycinek.odPozycji:wycinek.doPozycji], Total: razem,
	}, nil
}

// ustaleniePasujeBadania rozstrzyga zawężenia i przy okazji składa byt kontraktu
// — źródła ustalenia i tak trzeba odczytać, więc drugi przebieg byłby drugim
// odczytem tego samego.
func (a *adapterBadan) ustaleniePasujeBadania(ctx context.Context, ustalenie dane.UstalenieBadania,
	z shared.ResearchFindingListRequest) (bool, shared.ResearchFinding, error) {

	zrodla, err := a.repozytorium.ZrodlaUstalenia(ctx, ustalenie.ID)
	if err != nil {
		return false, shared.ResearchFinding{}, bladBadan(err)
	}
	przelozone := zlozUstalenieBadania(ustalenie, zrodla)

	if z.Query != nil && strings.TrimSpace(*z.Query) != "" {
		if !strings.Contains(strings.ToLower(przelozone.Content),
			strings.ToLower(strings.TrimSpace(*z.Query))) {
			return false, przelozone, nil
		}
	}
	if z.Status != nil && *z.Status != "" && ustalenie.Stan != *z.Status {
		return false, przelozone, nil
	}
	if z.SourceId != nil && *z.SourceId != "" && !zawieraNapisBadania(przelozone.SourceIds, *z.SourceId) {
		return false, przelozone, nil
	}
	if (z.Kind != nil && *z.Kind != "") || (z.Weight != nil && *z.Weight != "") {
		szczegoly, err := a.repozytorium.SzczegolyUstaleniaBadania(ctx, ustalenie.Kod)
		if err != nil {
			return false, przelozone, bladBadan(err)
		}
		if z.Kind != nil && *z.Kind != "" && szczegoly.Rodzaj != string(*z.Kind) {
			return false, przelozone, nil
		}
		if z.Weight != nil && *z.Weight != "" && szczegoly.Waga != string(*z.Weight) {
			return false, przelozone, nil
		}
	}
	if len(z.CodeIds) > 0 {
		kody, err := a.repozytorium.KodyUstalenia(ctx, ustalenie.Kod)
		if err != nil {
			return false, przelozone, bladBadan(err)
		}
		posiadane := make([]string, 0, len(kody))
		for _, kod := range kody {
			posiadane = append(posiadane, kod.Kod)
		}
		for _, zadany := range z.CodeIds {
			if !zawieraNapisBadania(posiadane, zadany) {
				return false, przelozone, nil
			}
		}
	}
	return true, przelozone, nil
}

// ZmienUstalenie obsługuje `research.finding.update`. Każda zmiana dopisuje wpis
// do śladu prowenancji — historia ustalenia jest wymogiem opracowania modułu,
// a nie dodatkiem.
func (a *adapterBadan) ZmienUstalenie(ctx context.Context,
	z shared.ResearchFindingUpdateRequest) (shared.ResearchFindingUpdateResponse, error) {

	if z.FindingId == "" {
		return shared.ResearchFindingUpdateResponse{}, bladWskazaniaBadan("finding.update bez ustalenia")
	}
	zastane, err := a.repozytorium.Ustalenie(ctx, z.FindingId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchFindingUpdateResponse{}, bladNieznanegoUstaleniaRaportu(z.FindingId)
	}
	if err != nil {
		return shared.ResearchFindingUpdateResponse{}, bladBadan(err)
	}

	if z.Content != nil {
		tresc := *z.Content
		zastane.Tresc = &tresc
	}
	if z.Status != nil && *z.Status != "" {
		zastane.Stan = *z.Status
	}
	zrodla, err := a.repozytorium.ZrodlaUstalenia(ctx, zastane.ID)
	if err != nil {
		return shared.ResearchFindingUpdateResponse{}, bladBadan(err)
	}
	kodyZrodel := make([]string, 0, len(zrodla))
	for _, zrodlo := range zrodla {
		kodyZrodel = append(kodyZrodel, zrodlo.Kod)
	}
	zapisane, err := a.repozytorium.ZapiszUstalenie(ctx, zastane, kodyZrodel)
	if err != nil {
		return shared.ResearchFindingUpdateResponse{}, bladBadan(err)
	}

	szczegoly := dane.SzczegolyUstaleniaBadania{Notatka: z.Note}
	if z.Kind != nil {
		szczegoly.Rodzaj = string(*z.Kind)
	}
	if z.Weight != nil {
		szczegoly.Waga = string(*z.Weight)
	}
	if z.NeedsConfirmation != nil {
		szczegoly.WymagaPotwierdzeni = *z.NeedsConfirmation
	} else {
		zastaneSzczegoly, err := a.repozytorium.SzczegolyUstaleniaBadania(ctx, z.FindingId)
		if err != nil {
			return shared.ResearchFindingUpdateResponse{}, bladBadan(err)
		}
		szczegoly.WymagaPotwierdzeni = zastaneSzczegoly.WymagaPotwierdzeni
	}
	if err := a.repozytorium.UstawSzczegolyUstalenia(ctx, z.FindingId, szczegoly); err != nil {
		return shared.ResearchFindingUpdateResponse{}, bladBadan(err)
	}
	_ = a.repozytorium.ZapiszProwenancje(ctx, dane.WpisProwenancjiBadania{
		UstalenieKod: z.FindingId, Aktor: string(shared.ActorKindOperator), Czynnosc: "zmiana ustalenia",
	})

	zrodlaPoZapisie, err := a.repozytorium.ZrodlaUstalenia(ctx, zapisane.ID)
	if err != nil {
		return shared.ResearchFindingUpdateResponse{}, bladBadan(err)
	}
	return shared.ResearchFindingUpdateResponse{
		Finding: zlozUstalenieBadania(zapisane, zrodlaPoZapisie),
	}, nil
}

// UsunUstalenie obsługuje `research.finding.remove`: usuwa ustalenie wskazane
// żądaniem z okna badania.
func (a *adapterBadan) UsunUstalenie(ctx context.Context,
	z shared.ResearchFindingRemoveRequest) (shared.ResearchFindingRemoveResponse, error) {

	if z.FindingId == "" {
		return shared.ResearchFindingRemoveResponse{}, bladWskazaniaBadan("finding.remove bez ustalenia")
	}
	odwiazane, err := a.repozytorium.UsunUstalenie(ctx, z.FindingId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchFindingRemoveResponse{}, bladNieznanegoUstaleniaRaportu(z.FindingId)
	}
	if err != nil {
		return shared.ResearchFindingRemoveResponse{}, bladBadan(err)
	}
	return shared.ResearchFindingRemoveResponse{
		FindingId: z.FindingId, DetachedSections: odwiazane,
	}, nil
}

// ScalUstalenia obsługuje `research.finding.merge`: przenosi źródła ustaleń
// scalanych na ustalenie docelowe.
func (a *adapterBadan) ScalUstalenia(ctx context.Context,
	z shared.ResearchFindingMergeRequest) (shared.ResearchFindingMergeResponse, error) {

	if z.TargetFindingId == "" {
		return shared.ResearchFindingMergeResponse{},
			bladWskazaniaBadan("finding.merge bez ustalenia docelowego")
	}
	if len(z.MergedFindingIds) == 0 {
		return shared.ResearchFindingMergeResponse{},
			bladWskazaniaBadan("finding.merge bez ustaleń do scalenia — nie ma czego przenieść")
	}
	przeniesione, err := a.repozytorium.PrzeniesZrodlaUstalen(ctx, z.MergedFindingIds, z.TargetFindingId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchFindingMergeResponse{}, bladNieznanegoUstaleniaRaportu(z.TargetFindingId)
	}
	if err != nil {
		return shared.ResearchFindingMergeResponse{}, bladBadan(err)
	}
	ustalenie, err := a.repozytorium.Ustalenie(ctx, z.TargetFindingId)
	if err != nil {
		return shared.ResearchFindingMergeResponse{}, bladBadan(err)
	}
	zrodla, err := a.repozytorium.ZrodlaUstalenia(ctx, ustalenie.ID)
	if err != nil {
		return shared.ResearchFindingMergeResponse{}, bladBadan(err)
	}
	_ = a.repozytorium.ZapiszProwenancje(ctx, dane.WpisProwenancjiBadania{
		UstalenieKod: z.TargetFindingId, Aktor: string(shared.ActorKindOperator),
		Czynnosc: "scalenie ustaleń: " + strings.Join(z.MergedFindingIds, ", "),
	})
	return shared.ResearchFindingMergeResponse{
		Finding: zlozUstalenieBadania(ustalenie, zrodla), MovedSources: przeniesione,
	}, nil
}

// SladUstalenia obsługuje `research.finding.provenance`: oddaje ślad
// pochodzenia ustalenia od źródła do wpisu w bazie.
func (a *adapterBadan) SladUstalenia(ctx context.Context,
	z shared.ResearchFindingProvenanceRequest) (shared.ResearchFindingProvenanceResponse, error) {

	if z.FindingId == "" {
		return shared.ResearchFindingProvenanceResponse{},
			bladWskazaniaBadan("finding.provenance bez ustalenia")
	}
	wpisy, err := a.repozytorium.Prowenancja(ctx, z.FindingId)
	if err != nil {
		return shared.ResearchFindingProvenanceResponse{}, bladBadan(err)
	}
	przelozone := make([]shared.ResearchProvenanceEntry, 0, len(wpisy))
	for _, wpis := range wpisy {
		pozycja := shared.ResearchProvenanceEntry{
			At: chwilaBazy(wpis.OCzasie), Actor: shared.ActorKind(wpis.Aktor),
			Action: wpis.Czynnosc, SourceId: wpis.ZrodloKod,
		}
		if wpis.Strona != nil {
			strona := int(*wpis.Strona)
			pozycja.Anchor = &shared.ResearchAnchor{
				Kind: shared.ResearchAnchorKindPage, Page: &strona,
			}
		}
		przelozone = append(przelozone, pozycja)
	}
	return shared.ResearchFindingProvenanceResponse{Entries: przelozone}, nil
}

// ── Kodowanie jakościowe ───────────────────────────────────────────────────

// PobierzKsiazkeKodow obsługuje `research.codebook.get`: oddaje książkę
// kodów jakościowych zapisaną przy oknie badania.
func (a *adapterBadan) PobierzKsiazkeKodow(ctx context.Context,
	z shared.ResearchCodebookGetRequest) (shared.ResearchCodebookGetResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchCodebookGetResponse{}, bladWskazaniaBadan("codebook.get bez okna badania")
	}
	kody, err := a.repozytorium.KsiazkaKodow(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchCodebookGetResponse{}, bladBadan(err)
	}
	return shared.ResearchCodebookGetResponse{Codes: zlozKodyBadania(kody)}, nil
}

// UstawKsiazkeKodow obsługuje `research.codebook.set`: zapisuje książkę
// kodów jakościowych okna badania.
func (a *adapterBadan) UstawKsiazkeKodow(ctx context.Context,
	z shared.ResearchCodebookSetRequest) (shared.ResearchCodebookSetResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchCodebookSetResponse{}, bladWskazaniaBadan("codebook.set bez okna badania")
	}
	doZapisu := make([]dane.KodBadania, 0, len(z.Codes))
	for _, kod := range z.Codes {
		identyfikator := kod.Id
		if strings.TrimSpace(identyfikator) == "" {
			identyfikator = nowyIdentyfikator(przedrostekKoduBadania)
		}
		doZapisu = append(doZapisu, dane.KodBadania{
			Kod: identyfikator, Okno: z.WindowId, Nazwa: kod.Name,
			Opis: kod.Description, Nadrzedny: kod.ParentId,
		})
	}
	zapisane, err := a.repozytorium.ZapiszKsiazkeKodow(ctx, z.WindowId, doZapisu)
	if err != nil {
		return shared.ResearchCodebookSetResponse{}, bladBadan(err)
	}
	return shared.ResearchCodebookSetResponse{Codes: zlozKodyBadania(zapisane)}, nil
}

// PrzypiszKodyUstalenia obsługuje `research.finding.code`. Nazwy nowych kodów
// zakładają pozycje książki kodów — Operator nie musi zakładać ich osobno.
func (a *adapterBadan) PrzypiszKodyUstalenia(ctx context.Context,
	z shared.ResearchFindingCodeRequest) (shared.ResearchFindingCodeResponse, error) {

	if z.FindingId == "" {
		return shared.ResearchFindingCodeResponse{}, bladWskazaniaBadan("finding.code bez ustalenia")
	}
	ustalenie, err := a.repozytorium.Ustalenie(ctx, z.FindingId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchFindingCodeResponse{}, bladNieznanegoUstaleniaRaportu(z.FindingId)
	}
	if err != nil {
		return shared.ResearchFindingCodeResponse{}, bladBadan(err)
	}

	kody := append([]string{}, z.CodeIds...)
	if len(z.NewCodeNames) > 0 {
		nowe := make([]dane.KodBadania, 0, len(z.NewCodeNames))
		for _, nazwa := range z.NewCodeNames {
			if strings.TrimSpace(nazwa) == "" {
				continue
			}
			identyfikator := nowyIdentyfikator(przedrostekKoduBadania)
			nowe = append(nowe, dane.KodBadania{
				Kod: identyfikator, Okno: ustalenie.Okno, Nazwa: nazwa,
			})
			kody = append(kody, identyfikator)
		}
		if _, err := a.repozytorium.ZapiszKsiazkeKodow(ctx, ustalenie.Okno, nowe); err != nil {
			return shared.ResearchFindingCodeResponse{}, bladBadan(err)
		}
	}
	if err := a.repozytorium.UstawKodyUstalenia(ctx, z.FindingId, kody); err != nil {
		return shared.ResearchFindingCodeResponse{}, bladBadan(err)
	}

	przypisane, err := a.repozytorium.KodyUstalenia(ctx, z.FindingId)
	if err != nil {
		return shared.ResearchFindingCodeResponse{}, bladBadan(err)
	}
	zrodla, err := a.repozytorium.ZrodlaUstalenia(ctx, ustalenie.ID)
	if err != nil {
		return shared.ResearchFindingCodeResponse{}, bladBadan(err)
	}
	return shared.ResearchFindingCodeResponse{
		Finding: zlozUstalenieBadania(ustalenie, zrodla), Codes: zlozKodyBadania(przypisane),
	}, nil
}

// MacierzKodowania obsługuje `research.finding.matrix`: krzyżuje kody
// jakościowe ze źródłami ustaleń, którym te kody nadano.
func (a *adapterBadan) MacierzKodowania(ctx context.Context,
	z shared.ResearchFindingMatrixRequest) (shared.ResearchFindingMatrixResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchFindingMatrixResponse{}, bladWskazaniaBadan("finding.matrix bez okna badania")
	}
	komorki, err := a.repozytorium.MacierzKodow(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchFindingMatrixResponse{}, bladBadan(err)
	}
	kody, err := a.repozytorium.KsiazkaKodow(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchFindingMatrixResponse{}, bladBadan(err)
	}

	przelozone := []shared.ResearchMatrixCell{}
	zrodla := map[string]bool{}
	for _, komorka := range komorki {
		if len(z.CodeIds) > 0 && !zawieraNapisBadania(z.CodeIds, komorka.KodKodu) {
			continue
		}
		if len(z.SourceIds) > 0 && !zawieraNapisBadania(z.SourceIds, komorka.ZrodloKod) {
			continue
		}
		zrodla[komorka.ZrodloKod] = true
		przelozone = append(przelozone, shared.ResearchMatrixCell{
			CodeId: komorka.KodKodu, SourceId: komorka.ZrodloKod,
			Count: komorka.Liczba, FindingIds: komorka.UstalenieKody,
		})
	}
	kodyZrodel := make([]string, 0, len(zrodla))
	for kod := range zrodla {
		kodyZrodel = append(kodyZrodel, kod)
	}
	return shared.ResearchFindingMatrixResponse{
		Cells: przelozone, Codes: zlozKodyBadania(kody),
		SourceIds: posortowaneNapisyBadania(kodyZrodel),
	}, nil
}

// zlozKodyBadania przekłada wiersze książki kodów z bazy danych na byty
// kodów jakościowych zwracane kontraktem komunikacji.
func zlozKodyBadania(kody []dane.KodBadania) []shared.ResearchCode {
	przelozone := make([]shared.ResearchCode, 0, len(kody))
	for _, kod := range kody {
		wystapienia := kod.Wystapienia
		pozycja := shared.ResearchCode{
			Id: kod.Kod, Name: kod.Nazwa, Description: kod.Opis, ParentId: kod.Nadrzedny,
		}
		if wystapienia > 0 {
			pozycja.Occurrences = &wystapienia
		}
		przelozone = append(przelozone, pozycja)
	}
	return przelozone
}

// ── Sprzeczności ───────────────────────────────────────────────────────────

// wzorzecLiczbyBadania wyławia liczbę z treści ustalenia — także z przecinkiem
// dziesiętnym, bo tak zapisuje się liczby w tekście polskim.
var wzorzecLiczbyBadania = regexp.MustCompile(`-?\d+(?:[.,]\d+)?`)

// SzukajSprzecznosci obsługuje `research.finding.contradictions`: wykrywa
// sprzeczności między ustaleniami heurystyką liczbową i modelem.
func (a *adapterBadan) SzukajSprzecznosci(ctx context.Context,
	z shared.ResearchFindingContradictionsRequest) (shared.ResearchFindingContradictionsResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchFindingContradictionsResponse{},
			bladWskazaniaBadan("finding.contradictions bez okna badania")
	}
	ustalenia, err := a.repozytorium.Ustalenia(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchFindingContradictionsResponse{}, bladBadan(err)
	}
	wybrane := []dane.UstalenieBadania{}
	for _, ustalenie := range ustalenia {
		if len(z.FindingIds) > 0 && !zawieraNapisBadania(z.FindingIds, ustalenie.Kod) {
			continue
		}
		wybrane = append(wybrane, ustalenie)
	}

	tolerancja := 10
	if z.NumericTolerance != nil && *z.NumericTolerance >= 0 {
		tolerancja = *z.NumericTolerance
	}

	for i := 0; i < len(wybrane); i++ {
		for j := i + 1; j < len(wybrane); j++ {
			streszczenie, roznica, jest := rozbieznoscLiczbowaBadania(wybrane[i], wybrane[j], tolerancja)
			if !jest {
				continue
			}
			kod := kodSprzecznosciBadania(wybrane[i].Kod, wybrane[j].Kod)
			zastana, err := a.repozytorium.Sprzecznosc(ctx, kod)
			if err == nil && zastana.Rozstrzygnieta {
				// Sprzeczność raz rozstrzygnięta nie wraca jako nowa przy kolejnym wykryciu.
				continue
			}
			if _, err := a.repozytorium.ZapiszSprzecznosc(ctx, dane.SprzecznoscBadania{
				Kod: kod, Okno: z.WindowId, Streszczenie: streszczenie,
				RoznicaLiczbowa: &roznica,
				UstalenieKody:   []string{wybrane[i].Kod, wybrane[j].Kod},
			}); err != nil {
				return shared.ResearchFindingContradictionsResponse{}, bladBadan(err)
			}
		}
	}

	sprzecznosci, err := a.repozytorium.Sprzecznosci(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchFindingContradictionsResponse{}, bladBadan(err)
	}
	przelozone := make([]shared.ResearchContradiction, 0, len(sprzecznosci))
	for _, sprzecznosc := range sprzecznosci {
		przelozone = append(przelozone, a.zlozSprzecznoscBadania(ctx, sprzecznosc))
	}
	return shared.ResearchFindingContradictionsResponse{Contradictions: przelozone}, nil
}

// kodSprzecznosciBadania nadaje sprzeczności kod wyprowadzony z pary ustaleń.
// Kod wyprowadzony, a nie losowy: powtórne wykrycie tej samej rozbieżności ma
// trafić w ten sam wiersz, zamiast mnożyć wpisy o tym samym.
func kodSprzecznosciBadania(pierwsze, drugie string) string {
	if pierwsze > drugie {
		pierwsze, drugie = drugie, pierwsze
	}
	return przedrostekSprzecznosciBadania + pierwsze + "-" + drugie
}

// rozbieznoscLiczbowaBadania wskazuje parę ustaleń mówiących o tym samym
// przedmiocie różnymi liczbami.
func rozbieznoscLiczbowaBadania(pierwsze, drugie dane.UstalenieBadania,
	tolerancjaProcent int) (string, string, bool) {

	trescA := trescUstaleniaDoPromptu(pierwsze)
	trescB := trescUstaleniaDoPromptu(drugie)
	if podobienstwoTekstuBadania(trescA, trescB) < 30 {
		return "", "", false
	}
	liczbyA := liczbyZTekstuBadania(trescA)
	liczbyB := liczbyZTekstuBadania(trescB)
	if len(liczbyA) == 0 || len(liczbyB) == 0 {
		return "", "", false
	}
	odniesienie := math.Max(math.Abs(liczbyA[0]), math.Abs(liczbyB[0]))
	if odniesienie == 0 {
		return "", "", false
	}
	roznica := math.Abs(liczbyA[0]-liczbyB[0]) / odniesienie * 100
	if roznica <= float64(tolerancjaProcent) {
		return "", "", false
	}
	streszczenie := "rozbieżne dane liczbowe w ustaleniach o tym samym przedmiocie: " +
		strconv.FormatFloat(liczbyA[0], 'f', -1, 64) + " wobec " +
		strconv.FormatFloat(liczbyB[0], 'f', -1, 64)
	return streszczenie, strconv.FormatFloat(roznica, 'f', 1, 64) + "%", true
}

// liczbyZTekstuBadania wyławia liczby z treści ustalenia, do porównania
// heurystyką sprzeczności liczbowej.
func liczbyZTekstuBadania(tekst string) []float64 {
	liczby := []float64{}
	for _, trafienie := range wzorzecLiczbyBadania.FindAllString(tekst, -1) {
		wartosc, err := strconv.ParseFloat(strings.ReplaceAll(trafienie, ",", "."), 64)
		if err == nil {
			liczby = append(liczby, wartosc)
		}
	}
	return liczby
}

// RozstrzygnijSprzecznosc obsługuje `research.contradiction.resolve`: zapisuje
// werdykt sprzeczności, żeby nie wracała jako otwarta.
func (a *adapterBadan) RozstrzygnijSprzecznosc(ctx context.Context,
	z shared.ResearchContradictionResolveRequest) (shared.ResearchContradictionResolveResponse, error) {

	if z.ContradictionId == "" {
		return shared.ResearchContradictionResolveResponse{},
			bladWskazaniaBadan("contradiction.resolve bez sprzeczności")
	}
	if strings.TrimSpace(z.Rationale) == "" {
		return shared.ResearchContradictionResolveResponse{},
			bladWskazaniaBadan("contradiction.resolve bez uzasadnienia — rozstrzygnięcie bez powodu " +
				"nie jest rozstrzygnięciem, tylko zamknięciem sprawy")
	}
	zastana, err := a.repozytorium.Sprzecznosc(ctx, z.ContradictionId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchContradictionResolveResponse{}, protokolBladBadania(
			shared.ErrorCodeNotFound, "sprzeczność nie istnieje: "+z.ContradictionId)
	}
	if err != nil {
		return shared.ResearchContradictionResolveResponse{}, bladBadan(err)
	}

	uzasadnienie := z.Rationale
	zastana.Rozstrzygnieta = true
	zastana.Uzasadnienie = &uzasadnienie
	zastana.UstalenieRozstrzyga = z.ResolutionFindingId
	zapisana, err := a.repozytorium.ZapiszSprzecznosc(ctx, zastana)
	if err != nil {
		return shared.ResearchContradictionResolveResponse{}, bladBadan(err)
	}
	return shared.ResearchContradictionResolveResponse{
		Contradiction: a.zlozSprzecznoscBadania(ctx, zapisana),
	}, nil
}

// zlozSprzecznoscBadania przekłada wiersz sprzeczności na byt kontraktu wraz ze
// źródłami ustaleń, których dotyczy.
func (a *adapterBadan) zlozSprzecznoscBadania(ctx context.Context,
	s dane.SprzecznoscBadania) shared.ResearchContradiction {

	kodyZrodel := map[string]bool{}
	for _, kodUstalenia := range s.UstalenieKody {
		ustalenie, err := a.repozytorium.Ustalenie(ctx, kodUstalenia)
		if err != nil {
			continue
		}
		zrodla, err := a.repozytorium.ZrodlaUstalenia(ctx, ustalenie.ID)
		if err != nil {
			continue
		}
		for _, zrodlo := range zrodla {
			kodyZrodel[zrodlo.Kod] = true
		}
	}
	lista := make([]string, 0, len(kodyZrodel))
	for kod := range kodyZrodel {
		lista = append(lista, kod)
	}
	return shared.ResearchContradiction{
		Id: s.Kod, FindingIds: s.UstalenieKody, SourceIds: posortowaneNapisyBadania(lista),
		Summary: s.Streszczenie, NumericDelta: s.RoznicaLiczbowa, Resolved: s.Rozstrzygnieta,
		ResolutionFindingId: s.UstalenieRozstrzyga, Rationale: s.Uzasadnienie,
	}
}

// ── Weryfikacja twierdzenia ────────────────────────────────────────────────

// ZweryfikujUstalenie obsługuje `research.finding.factCheck`. Weryfikacja idzie
// po korpusie badania: model dostaje treść ustalenia i fragmenty źródeł, a wynik
// zostaje zapisany przy ustaleniu.
func (a *adapterBadan) ZweryfikujUstalenie(ctx context.Context,
	z shared.ResearchFindingFactCheckRequest) (shared.ResearchFindingFactCheckResponse, error) {

	if z.FindingId == "" {
		return shared.ResearchFindingFactCheckResponse{},
			bladWskazaniaBadan("finding.factCheck bez ustalenia")
	}
	ustalenie, err := a.repozytorium.Ustalenie(ctx, z.FindingId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchFindingFactCheckResponse{}, bladNieznanegoUstaleniaRaportu(z.FindingId)
	}
	if err != nil {
		return shared.ResearchFindingFactCheckResponse{}, bladBadan(err)
	}
	twierdzenie := trescUstaleniaDoPromptu(ustalenie)
	if strings.TrimSpace(twierdzenie) == "" {
		return shared.ResearchFindingFactCheckResponse{}, bladWskazaniaBadan(
			"ustalenie " + z.FindingId + " nie ma treści — nie ma czego weryfikować")
	}

	zrodla, err := a.repozytorium.Zrodla(ctx, ustalenie.Okno)
	if err != nil {
		return shared.ResearchFindingFactCheckResponse{}, bladBadan(err)
	}
	trafienia := a.trafieniaKorpusuBadania(ctx, zrodla, nil, twierdzenie, 6)
	kanal, err := a.domyslnyKanalBadania()
	if err != nil {
		return shared.ResearchFindingFactCheckResponse{}, err
	}

	var polecenie strings.Builder
	polecenie.WriteString("Zweryfikuj twierdzenie względem podanych fragmentów źródeł. ")
	polecenie.WriteString("Odpowiedz WYŁĄCZNIE obiektem JSON o polach \"werdykt\" ")
	polecenie.WriteString("(supported, unsupported, contradicted albo inconclusive) ")
	polecenie.WriteString("oraz \"uzasadnienie\".\n\nTwierdzenie: ")
	polecenie.WriteString(twierdzenie)
	polecenie.WriteString("\n\nFragmenty:\n")
	for numer, trafienie := range trafienia {
		polecenie.WriteString("[" + strconv.Itoa(numer+1) + "] " + trafienie.tytul + ": " +
			trafienie.fragment + "\n")
	}
	if len(trafienia) == 0 {
		polecenie.WriteString("(korpus badania nie ma fragmentu pasującego do twierdzenia)\n")
	}

	odpowiedz, err := a.zapytajModel(ctx, ustalenie.Okno, kanal, polecenie.String())
	if err != nil {
		return shared.ResearchFindingFactCheckResponse{}, err
	}
	var rozlozone struct {
		Werdykt      string `json:"werdykt"`
		Uzasadnienie string `json:"uzasadnienie"`
	}
	if err := json.Unmarshal([]byte(wnetrzeJsonBadania(odpowiedz)), &rozlozone); err != nil {
		rozlozone.Werdykt = string(shared.ResearchFactCheckVerdictInconclusive)
		rozlozone.Uzasadnienie = strings.TrimSpace(odpowiedz)
	}
	if !werdyktZnanyBadania(rozlozone.Werdykt) {
		rozlozone.Werdykt = string(shared.ResearchFactCheckVerdictInconclusive)
	}

	sprawdzono, err := a.repozytorium.ZapiszWeryfikacje(ctx, z.FindingId,
		rozlozone.Werdykt, rozlozone.Uzasadnienie)
	if err != nil {
		return shared.ResearchFindingFactCheckResponse{}, bladBadan(err)
	}
	_ = a.repozytorium.ZapiszProwenancje(ctx, dane.WpisProwenancjiBadania{
		UstalenieKod: z.FindingId, Aktor: string(shared.ActorKindModel),
		Czynnosc: "weryfikacja twierdzenia: " + rozlozone.Werdykt,
	})

	dowody := make([]shared.ResearchExcerpt, 0, len(trafienia))
	for _, trafienie := range trafienia {
		dowody = append(dowody, shared.ResearchExcerpt{
			SourceId: trafienie.kodZrodla, SourceTitle: trafienie.tytul, Quote: trafienie.fragment,
			Anchor: shared.ResearchAnchor{Kind: shared.ResearchAnchorKindPage},
		})
	}
	return shared.ResearchFindingFactCheckResponse{Result: shared.ResearchFactCheck{
		FindingId: z.FindingId, Verdict: shared.ResearchFactCheckVerdict(rozlozone.Werdykt),
		Rationale: rozlozone.Uzasadnienie, Evidence: dowody, CheckedAt: chwilaBazy(sprawdzono),
	}}, nil
}

// werdyktZnanyBadania sprawdza, czy model oddał werdykt sprzeczności z listy
// wartości, które kontrakt komunikacji dopuszcza.
func werdyktZnanyBadania(werdykt string) bool {
	switch shared.ResearchFactCheckVerdict(werdykt) {
	case shared.ResearchFactCheckVerdictSupported, shared.ResearchFactCheckVerdictUnsupported,
		shared.ResearchFactCheckVerdictContradicted, shared.ResearchFactCheckVerdictInconclusive:
		return true
	}
	return false
}

// ── Wątki tematyczne ───────────────────────────────────────────────────────

// PogrupujUstalenia obsługuje `research.finding.cluster`. Grupowanie idzie po
// wspólnych kodach i po podobieństwie treści — bez magazynu wektorów, ale też
// bez udawania klastrowania, którego nie ma.
func (a *adapterBadan) PogrupujUstalenia(ctx context.Context,
	z shared.ResearchFindingClusterRequest) (shared.ResearchFindingClusterResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchFindingClusterResponse{}, bladWskazaniaBadan("finding.cluster bez okna badania")
	}
	ustalenia, err := a.repozytorium.Ustalenia(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchFindingClusterResponse{}, bladBadan(err)
	}
	wybrane := []dane.UstalenieBadania{}
	for _, ustalenie := range ustalenia {
		if len(z.FindingIds) > 0 && !zawieraNapisBadania(z.FindingIds, ustalenie.Kod) {
			continue
		}
		wybrane = append(wybrane, ustalenie)
	}

	przypisane := map[string]bool{}
	watki := []dane.WatekBadania{}
	for i := 0; i < len(wybrane); i++ {
		if przypisane[wybrane[i].Kod] {
			continue
		}
		grupa := []string{wybrane[i].Kod}
		for j := i + 1; j < len(wybrane); j++ {
			if przypisane[wybrane[j].Kod] {
				continue
			}
			if podobienstwoTekstuBadania(trescUstaleniaDoPromptu(wybrane[i]),
				trescUstaleniaDoPromptu(wybrane[j])) >= 25 {
				grupa = append(grupa, wybrane[j].Kod)
				przypisane[wybrane[j].Kod] = true
			}
		}
		if len(grupa) < 2 {
			continue
		}
		przypisane[wybrane[i].Kod] = true
		watki = append(watki, dane.WatekBadania{
			Kod:          nowyIdentyfikator(przedrostekWatkuBadania),
			Nazwa:        nazwaWatkuBadania(trescUstaleniaDoPromptu(wybrane[i])),
			Automatyczny: true, UstalenieKody: grupa,
		})
	}
	if z.TargetThreads != nil && *z.TargetThreads > 0 && len(watki) > *z.TargetThreads {
		sort.SliceStable(watki, func(i, j int) bool {
			return len(watki[i].UstalenieKody) > len(watki[j].UstalenieKody)
		})
		for _, nadmiarowy := range watki[*z.TargetThreads:] {
			for _, kod := range nadmiarowy.UstalenieKody {
				przypisane[kod] = false
			}
		}
		watki = watki[:*z.TargetThreads]
	}
	if _, err := a.repozytorium.ZapiszWatki(ctx, z.WindowId, watki); err != nil {
		return shared.ResearchFindingClusterResponse{}, bladBadan(err)
	}

	przelozone := make([]shared.ResearchThread, 0, len(watki))
	for _, watek := range watki {
		przelozone = append(przelozone, shared.ResearchThread{
			Id: watek.Kod, Name: watek.Nazwa, FindingIds: watek.UstalenieKody, Automatic: true,
		})
	}
	poza := []string{}
	for _, ustalenie := range wybrane {
		if !przypisane[ustalenie.Kod] {
			poza = append(poza, ustalenie.Kod)
		}
	}
	return shared.ResearchFindingClusterResponse{
		Threads: przelozone, UngroupedFindingIds: poza,
	}, nil
}

// nazwaWatkuBadania nadaje wątkowi nazwę z pierwszych słów ustalenia wiodącego.
// Nazwa z treści, a nie „Wątek 1": numer nie mówi Operatorowi, czego wątek dotyczy.
func nazwaWatkuBadania(tresc string) string {
	slowa := strings.Fields(tresc)
	if len(slowa) > 6 {
		slowa = slowa[:6]
	}
	nazwa := strings.Join(slowa, " ")
	if strings.TrimSpace(nazwa) == "" {
		return "wątek bez treści wiodącej"
	}
	return nazwa
}
