// Report Builder poza samą budową raportu: odczyt, szablony struktury,
// streszczenie zarządcze, bibliografia, przypisy, wstawki, wersje i różnice
// oraz tryb recenzji. Migawkę zakłada każda operacja zmieniająca treść raportu.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// szablonWbudowanyBadania opisuje jeden wzorzec struktury raportu wraz
// z kolejnością jego tytułów sekcji.
type szablonWbudowanyBadania struct {
	kod    string
	nazwa  string
	sekcje []string
}

// szablonyWbudowaneBadania wymieniają wzorce struktury raportu zaczerpnięte
// z opracowania modułu badań.
var szablonyWbudowaneBadania = []szablonWbudowanyBadania{
	{kod: "streszczenie-zarzadcze", nazwa: "Streszczenie zarządcze",
		sekcje: []string{"Streszczenie", "Kluczowe ustalenia", "Rekomendacje"}},
	{kod: "analiza-konkurencyjna", nazwa: "Analiza konkurencyjna",
		sekcje: []string{"Streszczenie", "Kontekst rynku", "Konkurenci", "Porównanie", "Wnioski"}},
	{kod: "analiza-swot", nazwa: "Analiza SWOT",
		sekcje: []string{"Wprowadzenie", "Mocne strony", "Słabe strony", "Szanse", "Zagrożenia", "Wnioski"}},
	{kod: "raport-benchmarkowy", nazwa: "Raport benchmarkowy",
		sekcje: []string{"Cel i metoda", "Kryteria", "Zestawienie", "Wnioski"}},
	{kod: "nota-badawcza", nazwa: "Nota badawcza",
		sekcje: []string{"Pytanie badawcze", "Ustalenia", "Ograniczenia"}},
	{kod: "przeglad-literatury", nazwa: "Przegląd literatury",
		sekcje: []string{"Protokół i kryteria", "Przesiew PRISMA", "Tabela dowodów",
			"Synteza narracyjna", "Luki badawcze"}},
}

// PobierzRaport obsługuje `research.report.get`. Żądanie bez wskazania raportu
// pyta o raport bieżący, a bieżącym jest ten zmieniony ostatnio.
func (a *adapterBadan) PobierzRaport(ctx context.Context,
	z shared.ResearchReportGetRequest) (shared.ResearchReportGetResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchReportGetResponse{}, bladWskazaniaBadan("report.get bez okna badania")
	}
	if z.ReportId != nil && strings.TrimSpace(*z.ReportId) != "" {
		raport, err := a.repozytorium.Raport(ctx, *z.ReportId)
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.ResearchReportGetResponse{}, bladNieznanegoRaportuBadania(*z.ReportId)
		}
		if err != nil {
			return shared.ResearchReportGetResponse{}, bladBadan(err)
		}
		przelozony, err := a.przelozRaport(ctx, raport)
		if err != nil {
			return shared.ResearchReportGetResponse{}, err
		}
		return shared.ResearchReportGetResponse{Report: &przelozony}, nil
	}

	raporty, err := a.repozytorium.Raporty(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchReportGetResponse{}, bladBadan(err)
	}
	if len(raporty) == 0 {
		// Brak raportu jest stanem badania, nie usterką: Report Builder pokazuje
		// pusty konspekt.
		return shared.ResearchReportGetResponse{}, nil
	}
	przelozony, err := a.przelozRaport(ctx, raporty[0])
	if err != nil {
		return shared.ResearchReportGetResponse{}, err
	}
	return shared.ResearchReportGetResponse{Report: &przelozony}, nil
}

// WypiszSzablonyRaportu obsługuje `research.report.template.list`, łącząc
// wzorce wbudowane z własnymi.
func (a *adapterBadan) WypiszSzablonyRaportu(ctx context.Context,
	_ shared.ResearchReportTemplateListRequest) (shared.ResearchReportTemplateListResponse, error) {

	szablony := []shared.ResearchReportTemplate{}
	for _, szablon := range szablonyWbudowaneBadania {
		szablony = append(szablony, shared.ResearchReportTemplate{
			Id: szablon.kod, Name: szablon.nazwa, SectionTitles: szablon.sekcje, Custom: false,
		})
	}
	wlasne, err := a.repozytorium.SzablonyRaportu(ctx)
	if err != nil {
		return shared.ResearchReportTemplateListResponse{}, bladBadan(err)
	}
	for _, szablon := range wlasne {
		szablony = append(szablony, shared.ResearchReportTemplate{
			Id: szablon.Kod, Name: szablon.Nazwa, SectionTitles: szablon.TytulySekcji,
			Custom: szablon.Wlasny,
		})
	}
	return shared.ResearchReportTemplateListResponse{Templates: szablony}, nil
}

// StreszczRaport obsługuje `research.report.summarize` — streszczenie zarządcze
// z ustaleń wskazanej wagi. Sekcja wchodzi na początek raportu i zostaje
// zapisana, a nie tylko zwrócona.
func (a *adapterBadan) StreszczRaport(ctx context.Context,
	z shared.ResearchReportSummarizeRequest) (shared.ResearchReportSummarizeResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchReportSummarizeResponse{}, bladWskazaniaBadan("report.summarize bez raportu")
	}
	raport, err := a.repozytorium.Raport(ctx, z.ReportId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchReportSummarizeResponse{}, bladNieznanegoRaportuBadania(z.ReportId)
	}
	if err != nil {
		return shared.ResearchReportSummarizeResponse{}, bladBadan(err)
	}

	ustalenia, err := a.repozytorium.Ustalenia(ctx, raport.Okno)
	if err != nil {
		return shared.ResearchReportSummarizeResponse{}, bladBadan(err)
	}
	wybrane := []dane.UstalenieBadania{}
	for _, ustalenie := range ustalenia {
		if z.Scope != nil && *z.Scope != "" {
			szczegoly, err := a.repozytorium.SzczegolyUstaleniaBadania(ctx, ustalenie.Kod)
			if err != nil {
				return shared.ResearchReportSummarizeResponse{}, bladBadan(err)
			}
			if szczegoly.Waga != string(*z.Scope) {
				continue
			}
		}
		wybrane = append(wybrane, ustalenie)
	}
	if len(wybrane) == 0 {
		return shared.ResearchReportSummarizeResponse{}, protokolBladBadania(shared.ErrorCodeNotFound,
			"badanie nie ma ustaleń o wskazanej wadze — streszczenie zarządcze nie ma z czego powstać")
	}

	sekcja, err := a.sekcjaStreszczenia(ctx, raport.Okno, wybrane)
	if err != nil {
		return shared.ResearchReportSummarizeResponse{}, err
	}
	sekcje, err := a.repozytorium.Sekcje(ctx, raport.ID)
	if err != nil {
		return shared.ResearchReportSummarizeResponse{}, bladBadan(err)
	}
	zlozone := []dane.SekcjaRaportu{sekcja}
	for numer, zastana := range sekcje {
		if zastana.Tytul == tytulStreszczeniaRaportu {
			continue
		}
		zastana.Kolejnosc = numer + 2
		zlozone = append(zlozone, zastana)
	}
	if _, err := a.repozytorium.ZapiszRaport(ctx, raport, zlozone); err != nil {
		return shared.ResearchReportSummarizeResponse{}, bladBadan(err)
	}
	if err := a.zapiszWersjeRaportuBadania(ctx, z.ReportId, "streszczenie zarządcze"); err != nil {
		return shared.ResearchReportSummarizeResponse{}, err
	}
	return shared.ResearchReportSummarizeResponse{
		Section: shared.ResearchReportSection{
			Id: sekcja.Kod, Title: sekcja.Tytul, Content: sekcja.Tresc,
			FindingIds: sekcja.UstalenieKody, Order: wskaznikLiczby(sekcja.Kolejnosc),
		},
		FromModel: true,
	}, nil
}

// ZlozBibliografie obsługuje `research.report.bibliography`, składając
// pozycje w stylu cytowania raportu.
func (a *adapterBadan) ZlozBibliografie(ctx context.Context,
	z shared.ResearchReportBibliographyRequest) (shared.ResearchReportBibliographyResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchReportBibliographyResponse{},
			bladWskazaniaBadan("report.bibliography bez raportu")
	}
	raport, err := a.repozytorium.Raport(ctx, z.ReportId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchReportBibliographyResponse{}, bladNieznanegoRaportuBadania(z.ReportId)
	}
	if err != nil {
		return shared.ResearchReportBibliographyResponse{}, bladBadan(err)
	}

	kody := []string{}
	if z.Scope == shared.ResearchBibliographyScopeCited {
		uzyte, err := a.zrodlaRaportuBadania(ctx, z.ReportId)
		if err != nil {
			return shared.ResearchReportBibliographyResponse{}, err
		}
		kody = posortowaneKluczeBadania(uzyte)
	} else {
		zrodla, err := a.repozytorium.Zrodla(ctx, raport.Okno)
		if err != nil {
			return shared.ResearchReportBibliographyResponse{}, bladBadan(err)
		}
		for _, zrodlo := range zrodla {
			kody = append(kody, zrodlo.Kod)
		}
	}

	styl := z.StyleId
	if strings.TrimSpace(styl) == "" {
		styl = "apa"
	}
	pozycje := make([]string, 0, len(kody))
	for _, kod := range kody {
		metadane, err := a.metadaneCytowaniaBadania(ctx, kod)
		if err != nil {
			return shared.ResearchReportBibliographyResponse{}, err
		}
		pozycje = append(pozycje, pozycjaBibliograficznaBadania(styl, metadane))
	}
	sort.Strings(pozycje)
	return shared.ResearchReportBibliographyResponse{Entries: pozycje, SourceIds: kody}, nil
}

// UstawPrzypisyRaportu obsługuje `research.report.footnote.set`. Liczba przypisów
// jest policzona z sekcji, a nie zadeklarowana: menedżer przypisów ma mówić, ile
// odnośników naprawdę stoi w treści.
func (a *adapterBadan) UstawPrzypisyRaportu(ctx context.Context,
	z shared.ResearchReportFootnoteSetRequest) (shared.ResearchReportFootnoteSetResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchReportFootnoteSetResponse{},
			bladWskazaniaBadan("report.footnote.set bez raportu")
	}
	if strings.TrimSpace(z.Placement) == "" {
		return shared.ResearchReportFootnoteSetResponse{},
			bladWskazaniaBadan("report.footnote.set bez umiejscowienia przypisów")
	}
	skrocone := z.ShortForms != nil && *z.ShortForms
	err := a.repozytorium.UstawPrzypisyRaportu(ctx, z.ReportId, z.Placement, skrocone)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchReportFootnoteSetResponse{}, bladNieznanegoRaportuBadania(z.ReportId)
	}
	if err != nil {
		return shared.ResearchReportFootnoteSetResponse{}, bladBadan(err)
	}

	raport, err := a.repozytorium.Raport(ctx, z.ReportId)
	if err != nil {
		return shared.ResearchReportFootnoteSetResponse{}, bladBadan(err)
	}
	sekcje, err := a.repozytorium.Sekcje(ctx, raport.ID)
	if err != nil {
		return shared.ResearchReportFootnoteSetResponse{}, bladBadan(err)
	}
	przypisy := 0
	for _, sekcja := range sekcje {
		przypisy += len(sekcja.UstalenieKody)
	}
	return shared.ResearchReportFootnoteSetResponse{FootnoteCount: przypisy}, nil
}

// WstawBlokRaportu obsługuje `research.report.insert`, dokładając wstawkę
// z danymi badania do wskazanej sekcji.
func (a *adapterBadan) WstawBlokRaportu(ctx context.Context,
	z shared.ResearchReportInsertRequest) (shared.ResearchReportInsertResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchReportInsertResponse{}, bladWskazaniaBadan("report.insert bez raportu")
	}
	if z.SectionId == "" {
		return shared.ResearchReportInsertResponse{}, bladWskazaniaBadan("report.insert bez sekcji")
	}
	if z.Kind == "" {
		return shared.ResearchReportInsertResponse{}, bladWskazaniaBadan("report.insert bez rodzaju wstawki")
	}
	raport, err := a.repozytorium.Raport(ctx, z.ReportId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchReportInsertResponse{}, bladNieznanegoRaportuBadania(z.ReportId)
	}
	if err != nil {
		return shared.ResearchReportInsertResponse{}, bladBadan(err)
	}

	naglowki, wiersze, err := a.zawartoscBlokuBadania(ctx, raport.Okno, z)
	if err != nil {
		return shared.ResearchReportInsertResponse{}, err
	}
	surowe, err := json.Marshal(wiersze)
	if err != nil {
		return shared.ResearchReportInsertResponse{}, bladBadan(err)
	}
	blok := dane.BlokRaportuBadania{
		Kod: nowyIdentyfikator(przedrostekBlokuBadania), RaportKod: z.ReportId,
		SekcjaKod: z.SectionId, Rodzaj: string(z.Kind), Naglowki: naglowki,
		Wiersze: string(surowe), UstalenieKody: z.FindingIds,
	}
	if _, err := a.repozytorium.ZapiszBlokRaportu(ctx, blok); err != nil {
		return shared.ResearchReportInsertResponse{}, bladBadan(err)
	}
	if err := a.zapiszWersjeRaportuBadania(ctx, z.ReportId, "wstawka "+string(z.Kind)); err != nil {
		return shared.ResearchReportInsertResponse{}, err
	}
	return shared.ResearchReportInsertResponse{Block: shared.ResearchReportBlock{
		Id: blok.Kod, SectionId: z.SectionId, Kind: z.Kind, Headers: naglowki,
		Rows: surowe, FindingIds: z.FindingIds,
	}}, nil
}

// zawartoscBlokuBadania składa dane wstawki z bytów badania. Każdy rodzaj bierze
// dane z innego miejsca, ale żaden nie bierze ich znikąd: wstawka bez danych
// byłaby ramką, którą Operator wypełniałby ręcznie.
func (a *adapterBadan) zawartoscBlokuBadania(ctx context.Context, okno string,
	z shared.ResearchReportInsertRequest) ([]string, [][]string, error) {

	switch z.Kind {
	case shared.ResearchBlockKindMatrix:
		odpowiedz, err := a.MacierzKodowania(ctx,
			shared.ResearchFindingMatrixRequest{WindowId: okno})
		if err != nil {
			return nil, nil, err
		}
		naglowki := append([]string{"kod"}, odpowiedz.SourceIds...)
		wiersze := [][]string{}
		for _, kod := range odpowiedz.Codes {
			wiersz := []string{kod.Name}
			for _, kodZrodla := range odpowiedz.SourceIds {
				liczba := 0
				for _, komorka := range odpowiedz.Cells {
					if komorka.CodeId == kod.Id && komorka.SourceId == kodZrodla {
						liczba = komorka.Count
					}
				}
				wiersz = append(wiersz, strconv.Itoa(liczba))
			}
			wiersze = append(wiersze, wiersz)
		}
		return naglowki, wiersze, nil

	case shared.ResearchBlockKindTimeline:
		zrodla, err := a.repozytorium.Zrodla(ctx, okno)
		if err != nil {
			return nil, nil, bladBadan(err)
		}
		wiersze := [][]string{}
		for _, zrodlo := range zrodla {
			wiersze = append(wiersze, []string{zrodlo.PozyskanoO, zrodlo.Tytul})
		}
		sort.SliceStable(wiersze, func(i, j int) bool { return wiersze[i][0] < wiersze[j][0] })
		return []string{"data", "pozycja"}, wiersze, nil

	case shared.ResearchBlockKindChart:
		ustalenia, err := a.repozytorium.Ustalenia(ctx, okno)
		if err != nil {
			return nil, nil, bladBadan(err)
		}
		wiersze := [][]string{}
		for _, ustalenie := range ustalenia {
			if len(z.FindingIds) > 0 && !zawieraNapisBadania(z.FindingIds, ustalenie.Kod) {
				continue
			}
			liczby := liczbyZTekstuBadania(trescUstaleniaDoPromptu(ustalenie))
			if len(liczby) == 0 {
				continue
			}
			wiersze = append(wiersze, []string{
				skrocDoBadania(trescUstaleniaDoPromptu(ustalenie), 80),
				strconv.FormatFloat(liczby[0], 'f', -1, 64),
			})
		}
		if len(wiersze) == 0 {
			return nil, nil, protokolBladBadania(shared.ErrorCodeNotFound,
				"żadne ze wskazanych ustaleń nie niesie danej liczbowej — wykres nie ma z czego powstać")
		}
		return []string{"ustalenie", "wartosc"}, wiersze, nil

	default:
		ustalenia, err := a.repozytorium.Ustalenia(ctx, okno)
		if err != nil {
			return nil, nil, bladBadan(err)
		}
		wiersze := [][]string{}
		for _, ustalenie := range ustalenia {
			if len(z.FindingIds) > 0 && !zawieraNapisBadania(z.FindingIds, ustalenie.Kod) {
				continue
			}
			zrodla, err := a.repozytorium.ZrodlaUstalenia(ctx, ustalenie.ID)
			if err != nil {
				return nil, nil, bladBadan(err)
			}
			tytuly := make([]string, 0, len(zrodla))
			for _, zrodlo := range zrodla {
				tytuly = append(tytuly, zrodlo.Tytul)
			}
			szczegoly, err := a.repozytorium.SzczegolyUstaleniaBadania(ctx, ustalenie.Kod)
			if err != nil {
				return nil, nil, bladBadan(err)
			}
			wiersze = append(wiersze, []string{
				trescUstaleniaDoPromptu(ustalenie), strings.Join(tytuly, "; "),
				szczegoly.Rodzaj, szczegoly.Waga,
			})
		}
		return []string{"ustalenie", "zrodla", "rodzaj", "waga"}, wiersze, nil
	}
}

// ── Wersje i różnice ───────────────────────────────────────────────────────

// zapiszWersjeRaportuBadania zakłada migawkę sekcji raportu pod etykietą
// operacji, która ją wywołała teraz.
func (a *adapterBadan) zapiszWersjeRaportuBadania(ctx context.Context, kodRaportu, etykieta string) error {
	raport, err := a.repozytorium.Raport(ctx, kodRaportu)
	if err != nil {
		return bladBadan(err)
	}
	sekcje, err := a.repozytorium.Sekcje(ctx, raport.ID)
	if err != nil {
		return bladBadan(err)
	}
	migawka := make([]map[string]string, 0, len(sekcje))
	for _, sekcja := range sekcje {
		tresc := ""
		if sekcja.Tresc != nil {
			tresc = *sekcja.Tresc
		}
		migawka = append(migawka, map[string]string{
			"id": sekcja.Kod, "tytul": sekcja.Tytul, "tresc": tresc,
		})
	}
	surowa, err := json.Marshal(migawka)
	if err != nil {
		return bladBadan(err)
	}
	nazwa := etykieta
	_, err = a.repozytorium.ZapiszWersjeRaportu(ctx, dane.WersjaRaportuBadania{
		Kod: nowyIdentyfikator(przedrostekWersjiRaportuBadania), RaportKod: kodRaportu,
		Etykieta: &nazwa, Migawka: string(surowa), LiczbaSekcji: len(sekcje),
	})
	if err != nil {
		return bladBadan(err)
	}
	return nil
}

// WypiszWersjeRaportu obsługuje `research.report.version.list`, oddając
// wersje raportu od najnowszej do najstarszej.
func (a *adapterBadan) WypiszWersjeRaportu(ctx context.Context,
	z shared.ResearchReportVersionListRequest) (shared.ResearchReportVersionListResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchReportVersionListResponse{},
			bladWskazaniaBadan("report.version.list bez raportu")
	}
	limit := 0
	if z.Limit != nil {
		limit = *z.Limit
	}
	wersje, err := a.repozytorium.WersjeRaportu(ctx, z.ReportId, limit)
	if err != nil {
		return shared.ResearchReportVersionListResponse{}, bladBadan(err)
	}
	przelozone := make([]shared.ResearchReportVersion, 0, len(wersje))
	for _, wersja := range wersje {
		przelozone = append(przelozone, shared.ResearchReportVersion{
			Id: wersja.Kod, ReportId: wersja.RaportKod, Label: wersja.Etykieta,
			SectionCount: wersja.LiczbaSekcji, CreatedAt: chwilaBazy(wersja.Utworzono),
		})
	}
	return shared.ResearchReportVersionListResponse{Versions: przelozone}, nil
}

// PorownajWersjeRaportu obsługuje `research.report.diff`, składając
// fragmenty różnicy z dwóch migawek.
func (a *adapterBadan) PorownajWersjeRaportu(ctx context.Context,
	z shared.ResearchReportDiffRequest) (shared.ResearchReportDiffResponse, error) {

	if z.BaseVersionId == "" || z.TargetVersionId == "" {
		return shared.ResearchReportDiffResponse{},
			bladWskazaniaBadan("report.diff bez obu wersji do porównania")
	}
	podstawa, err := a.repozytorium.WersjaRaportuBadania(ctx, z.BaseVersionId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchReportDiffResponse{}, protokolBladBadania(shared.ErrorCodeNotFound,
			"wersja raportu nie istnieje: "+z.BaseVersionId)
	}
	if err != nil {
		return shared.ResearchReportDiffResponse{}, bladBadan(err)
	}
	docelowa, err := a.repozytorium.WersjaRaportuBadania(ctx, z.TargetVersionId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchReportDiffResponse{}, protokolBladBadania(shared.ErrorCodeNotFound,
			"wersja raportu nie istnieje: "+z.TargetVersionId)
	}
	if err != nil {
		return shared.ResearchReportDiffResponse{}, bladBadan(err)
	}

	przed := sekcjeMigawkiBadania(podstawa.Migawka)
	po := sekcjeMigawkiBadania(docelowa.Migawka)
	fragmenty := []shared.StudioDiffHunk{}
	numer := 0
	widziane := map[string]bool{}

	for _, sekcja := range po {
		numer++
		widziane[sekcja["id"]] = true
		stara, jest := przed[sekcja["id"]]
		nowa := sekcja["tytul"] + "\n" + sekcja["tresc"]
		if !jest {
			fragmenty = append(fragmenty, shared.StudioDiffHunk{
				Index: numer, Kind: shared.DiffHunkKindAdded, After: &nowa,
			})
			continue
		}
		poprzednia := stara["tytul"] + "\n" + stara["tresc"]
		if poprzednia == nowa {
			fragmenty = append(fragmenty, shared.StudioDiffHunk{
				Index: numer, Kind: shared.DiffHunkKindContext, Before: &poprzednia, After: &nowa,
			})
			continue
		}
		fragmenty = append(fragmenty, shared.StudioDiffHunk{
			Index: numer, Kind: shared.DiffHunkKindChanged, Before: &poprzednia, After: &nowa,
		})
	}
	for kod, sekcja := range przed {
		if widziane[kod] {
			continue
		}
		numer++
		poprzednia := sekcja["tytul"] + "\n" + sekcja["tresc"]
		fragmenty = append(fragmenty, shared.StudioDiffHunk{
			Index: numer, Kind: shared.DiffHunkKindRemoved, Before: &poprzednia,
		})
	}
	sort.SliceStable(fragmenty, func(i, j int) bool { return fragmenty[i].Index < fragmenty[j].Index })
	return shared.ResearchReportDiffResponse{Hunks: fragmenty}, nil
}

// sekcjeMigawkiBadania rozkłada migawkę wersji na mapę sekcji po
// identyfikatorze każdej sekcji raportu.
func sekcjeMigawkiBadania(migawka string) map[string]map[string]string {
	var lista []map[string]string
	rozlozone := map[string]map[string]string{}
	if err := json.Unmarshal([]byte(migawka), &lista); err != nil {
		return rozlozone
	}
	for _, sekcja := range lista {
		rozlozone[sekcja["id"]] = sekcja
	}
	return rozlozone
}

// ── Tryb recenzji ──────────────────────────────────────────────────────────

// DodajKomentarzRaportu obsługuje `research.report.comment.add`, zapisując
// komentarz wątku recenzji raportu.
func (a *adapterBadan) DodajKomentarzRaportu(ctx context.Context,
	z shared.ResearchReportCommentAddRequest) (shared.ResearchReportCommentAddResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchReportCommentAddResponse{},
			bladWskazaniaBadan("report.comment.add bez raportu")
	}
	if strings.TrimSpace(z.Content) == "" {
		return shared.ResearchReportCommentAddResponse{},
			bladWskazaniaBadan("report.comment.add bez treści komentarza")
	}
	rozstrzygniety := z.Resolved != nil && *z.Resolved
	zapisany, err := a.repozytorium.ZapiszKomentarzRaportu(ctx, dane.KomentarzRaportuBadania{
		Kod: nowyIdentyfikator(przedrostekKomentarzaBadania), RaportKod: z.ReportId,
		SekcjaKod: z.SectionId, WatekKod: z.ThreadId, Tresc: z.Content,
		Cytat: z.Quote, Rozstrzygniety: rozstrzygniety,
	})
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchReportCommentAddResponse{}, bladNieznanegoRaportuBadania(z.ReportId)
	}
	if err != nil {
		return shared.ResearchReportCommentAddResponse{}, bladBadan(err)
	}
	return shared.ResearchReportCommentAddResponse{Comment: zlozKomentarzBadania(zapisany)}, nil
}

// WypiszKomentarzeRaportu obsługuje `research.report.comment.list`, oddając
// komentarze wybranego raportu.
func (a *adapterBadan) WypiszKomentarzeRaportu(ctx context.Context,
	z shared.ResearchReportCommentListRequest) (shared.ResearchReportCommentListResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchReportCommentListResponse{},
			bladWskazaniaBadan("report.comment.list bez raportu")
	}
	tylkoOtwarte := z.OpenOnly != nil && *z.OpenOnly
	komentarze, err := a.repozytorium.KomentarzeRaportu(ctx, z.ReportId, tylkoOtwarte)
	if err != nil {
		return shared.ResearchReportCommentListResponse{}, bladBadan(err)
	}
	przelozone := make([]shared.ResearchReportComment, 0, len(komentarze))
	for _, komentarz := range komentarze {
		przelozone = append(przelozone, zlozKomentarzBadania(komentarz))
	}
	return shared.ResearchReportCommentListResponse{Comments: przelozone}, nil
}

// zlozKomentarzBadania przekłada wiersz komentarza bazy danych na byt
// kontraktu okna recenzji raportu.
func zlozKomentarzBadania(k dane.KomentarzRaportuBadania) shared.ResearchReportComment {
	return shared.ResearchReportComment{
		Id: k.Kod, ReportId: k.RaportKod, SectionId: k.SekcjaKod, ThreadId: k.WatekKod,
		Content: k.Tresc, Quote: k.Cytat, Resolved: k.Rozstrzygniety,
		CreatedAt: chwilaBazy(k.Utworzono),
	}
}

// OperacjaKontekstowaRaportu obsługuje `research.report.contextual.op`:
// korektę, streszczenie, rozwinięcie i zmianę stylu zaznaczonego fragmentu.
// Skutkiem jest treść sekcji zmieniona w bazie oraz nowa wersja raportu.
func (a *adapterBadan) OperacjaKontekstowaRaportu(ctx context.Context,
	z shared.ResearchReportContextualOpRequest) (shared.ResearchReportContextualOpResponse, error) {

	if z.ReportId == "" || z.SectionId == "" {
		return shared.ResearchReportContextualOpResponse{},
			bladWskazaniaBadan("report.contextual.op bez raportu albo bez sekcji")
	}
	polecenieOperacji, znane := polecenieOperacjiBadania(z.ActionId)
	if !znane {
		return shared.ResearchReportContextualOpResponse{}, bladWskazaniaBadan(
			"operacja „" + z.ActionId + "” nie jest jedną z: korekta, streszczenie, " +
				"rozwiniecie, zmiana-stylu")
	}
	raport, err := a.repozytorium.Raport(ctx, z.ReportId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchReportContextualOpResponse{}, bladNieznanegoRaportuBadania(z.ReportId)
	}
	if err != nil {
		return shared.ResearchReportContextualOpResponse{}, bladBadan(err)
	}
	sekcje, err := a.repozytorium.Sekcje(ctx, raport.ID)
	if err != nil {
		return shared.ResearchReportContextualOpResponse{}, bladBadan(err)
	}

	numerSekcji := -1
	for numer, sekcja := range sekcje {
		if sekcja.Kod == z.SectionId {
			numerSekcji = numer
		}
	}
	if numerSekcji < 0 {
		return shared.ResearchReportContextualOpResponse{}, protokolBladBadania(
			shared.ErrorCodeNotFound, "raport "+z.ReportId+" nie ma sekcji "+z.SectionId)
	}
	tresc := ""
	if sekcje[numerSekcji].Tresc != nil {
		tresc = *sekcje[numerSekcji].Tresc
	}
	if strings.TrimSpace(tresc) == "" {
		return shared.ResearchReportContextualOpResponse{}, bladWskazaniaBadan(
			"sekcja " + z.SectionId + " nie ma treści — nie ma na czym wykonać operacji")
	}

	od, do := zakresZaznaczeniaBadania(tresc, z.SelectionStart, z.SelectionEnd)
	zaznaczenie := tresc[od:do]
	kanal, err := a.domyslnyKanalBadania()
	if err != nil {
		return shared.ResearchReportContextualOpResponse{}, err
	}
	wynik, err := a.zapytajModel(ctx, raport.Okno, kanal, polecenieOperacji+"\n\n"+zaznaczenie)
	if err != nil {
		return shared.ResearchReportContextualOpResponse{}, err
	}

	nowa := tresc[:od] + strings.TrimSpace(wynik) + tresc[do:]
	sekcje[numerSekcji].Tresc = &nowa
	if _, err := a.repozytorium.ZapiszRaport(ctx, raport, sekcje); err != nil {
		return shared.ResearchReportContextualOpResponse{}, bladBadan(err)
	}
	if err := a.zapiszWersjeRaportuBadania(ctx, z.ReportId, "operacja "+z.ActionId); err != nil {
		return shared.ResearchReportContextualOpResponse{}, err
	}
	wynikTekst := strings.TrimSpace(wynik)
	return shared.ResearchReportContextualOpResponse{ResultText: &wynikTekst}, nil
}

// polecenieOperacjiBadania odwzorowuje nazwę operacji na polecenie dla modelu.
// Nazwy są po czynnościach, tak jak w module Studio — Operator wybiera „korekta",
// a nie numer operacji.
func polecenieOperacjiBadania(nazwa string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(nazwa)) {
	case "korekta":
		return "Popraw poniższy fragment raportu językowo i stylistycznie. Zachowaj sens " +
			"i wszystkie liczby. Odpowiedz samym poprawionym tekstem.", true
	case "streszczenie":
		return "Streść poniższy fragment raportu do jednego zwięzłego akapitu. " +
			"Odpowiedz samym streszczeniem.", true
	case "rozwiniecie":
		return "Rozwiń poniższy fragment raportu, nie dodając faktów spoza niego. " +
			"Odpowiedz samym rozwiniętym tekstem.", true
	case "zmiana-stylu":
		return "Przepisz poniższy fragment raportu w rejestrze formalnym, właściwym " +
			"raportowi badawczemu. Odpowiedz samym przepisanym tekstem.", true
	default:
		return "", false
	}
}

// zakresZaznaczeniaBadania sprowadza wskazanie zaznaczenia do granic treści.
// Wskazanie poza treścią obejmuje całość, zamiast wywracać operację — Operator
// zaznaczający „wszystko" nie ma podawać długości sekcji.
func zakresZaznaczeniaBadania(tresc string, poczatek, koniec *int) (int, int) {
	od := 0
	do := len(tresc)
	if poczatek != nil && *poczatek > 0 && *poczatek < len(tresc) {
		od = *poczatek
	}
	if koniec != nil && *koniec > od && *koniec <= len(tresc) {
		do = *koniec
	}
	return od, do
}
