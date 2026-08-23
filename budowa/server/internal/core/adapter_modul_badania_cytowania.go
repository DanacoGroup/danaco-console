// Odpowiedzialność pliku: cytowania i bibliografia — złożenie cytatu w stylu,
// wykaz stylów, kontrola kompletności metadanych (Citation Health Check) oraz
// flaga wycofania pracy.
//
// ── Style są wkompilowane, nie doczytywane ─────────────────────────────────
// Repozytorium CSL liczy około dwóch i pół tysiąca plików XML i procesor, który
// je czyta, jest biblioteką JavaScriptu. Oparcie cytowań o taki procesor
// znaczyłoby, że u Operatora, który nie doinstalował środowiska JS, cytowanie
// nie działa wcale — a cytowanie jest w module badawczym czynnością codzienną,
// nie ozdobą. Dlatego pięć stylów najczęściej wymaganych (APA, MLA, Chicago,
// IEEE, Vancouver) jest złożonych tutaj, w Go, i działa od pierwszego
// uruchomienia. Styl własny Operatora dokłada się obok jako nazwany wariant.
//
// ── Kontrola kompletności mówi, czego brak ─────────────────────────────────
// `research.citation.check` nie mówi „metadane niekompletne". Mówi, które pole
// brakuje przy którym źródle i skąd da się je uzupełnić — bo to jest jedyna
// postać tej informacji, z którą Operator może cokolwiek zrobić.
package core

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// stylWbudowanyBadania opisuje jeden styl cytowania złożony w rdzeniu.
type stylWbudowanyBadania struct {
	kod   string
	nazwa string
}

// styleWbudowaneBadania wymienia style dostępne bez żadnej konfiguracji.
var styleWbudowaneBadania = []stylWbudowanyBadania{
	{kod: "apa", nazwa: "APA (7. wydanie)"},
	{kod: "mla", nazwa: "MLA (9. wydanie)"},
	{kod: "chicago", nazwa: "Chicago (notes-bibliography)"},
	{kod: "ieee", nazwa: "IEEE"},
	{kod: "vancouver", nazwa: "Vancouver"},
}

// metadaneCytowaniaBadania to pola pozycji potrzebne do złożenia cytatu.
type metadaneCytowaniaBadania struct {
	kodZrodla string
	tytul     string
	autorzy   []string
	rok       string
	wydawca   string
	adres     string
	doi       string
}

// ZlozCytowania obsługuje `research.citation.render`.
func (a *adapterBadan) ZlozCytowania(ctx context.Context,
	z shared.ResearchCitationRenderRequest) (shared.ResearchCitationRenderResponse, error) {

	if len(z.SourceIds) == 0 {
		return shared.ResearchCitationRenderResponse{},
			bladWskazaniaBadan("citation.render bez źródeł — nie ma czego cytować")
	}
	styl := z.StyleId
	if strings.TrimSpace(styl) == "" {
		styl = "apa"
	}
	cytaty := make([]shared.ResearchCitation, 0, len(z.SourceIds))
	niekompletne := []string{}
	for _, kod := range z.SourceIds {
		metadane, err := a.metadaneCytowaniaBadania(ctx, kod)
		if err != nil {
			return shared.ResearchCitationRenderResponse{}, err
		}
		if len(brakiCytowaniaBadania(metadane)) > 0 {
			niekompletne = append(niekompletne, kod)
		}
		cytat := shared.ResearchCitation{SourceId: kod, StyleId: styl}
		if z.Mode == shared.ResearchCitationModeInText || z.Mode == shared.ResearchCitationModeBoth {
			wTekscie := cytatWTekscieBadania(styl, metadane)
			cytat.InText = &wTekscie
		}
		if z.Mode == shared.ResearchCitationModeBibliography || z.Mode == shared.ResearchCitationModeBoth {
			pozycja := pozycjaBibliograficznaBadania(styl, metadane)
			cytat.Bibliography = &pozycja
		}
		cytaty = append(cytaty, cytat)
	}
	return shared.ResearchCitationRenderResponse{
		Citations: cytaty, IncompleteSourceIds: niekompletne,
	}, nil
}

// metadaneCytowaniaBadania składa pola pozycji z wiersza źródła i z metadanych
// CSL, gdy zostały zapisane przy rozstrzygnięciu identyfikatora.
func (a *adapterBadan) metadaneCytowaniaBadania(ctx context.Context,
	kodZrodla string) (metadaneCytowaniaBadania, error) {

	zrodlo, err := a.repozytorium.Zrodlo(ctx, kodZrodla)
	if err != nil {
		return metadaneCytowaniaBadania{}, bladNieznanegoZrodlaBadania(kodZrodla)
	}
	metadane := metadaneCytowaniaBadania{kodZrodla: kodZrodla, tytul: zrodlo.Tytul}
	if zrodlo.Adres != nil {
		metadane.adres = *zrodlo.Adres
	}
	if zrodlo.Pochodzenie != nil {
		metadane.wydawca = *zrodlo.Pochodzenie
	}
	lektura, err := a.repozytorium.LekturaZrodlaBadania(ctx, kodZrodla)
	if err != nil {
		return metadaneCytowaniaBadania{}, bladBadan(err)
	}
	if lektura.Identyfikat != nil {
		metadane.doi = *lektura.Identyfikat
	}
	if lektura.CslJson != nil && strings.TrimSpace(*lektura.CslJson) != "" {
		var csl struct {
			Title  string `json:"title"`
			Author []struct {
				Given  string `json:"given"`
				Family string `json:"family"`
			} `json:"author"`
			Issued struct {
				DateParts [][]int `json:"date-parts"`
			} `json:"issued"`
			Publisher string `json:"publisher"`
		}
		if err := json.Unmarshal([]byte(*lektura.CslJson), &csl); err == nil {
			if strings.TrimSpace(csl.Title) != "" {
				metadane.tytul = csl.Title
			}
			for _, autor := range csl.Author {
				metadane.autorzy = append(metadane.autorzy,
					strings.TrimSpace(autor.Family+", "+autor.Given))
			}
			if len(csl.Issued.DateParts) > 0 && len(csl.Issued.DateParts[0]) > 0 {
				metadane.rok = strconv.Itoa(csl.Issued.DateParts[0][0])
			}
			if strings.TrimSpace(csl.Publisher) != "" {
				metadane.wydawca = csl.Publisher
			}
		}
	}
	if metadane.rok == "" {
		if rok := rokZTekstuBadania(zrodlo.PozyskanoO); rok != nil {
			// Rok pozyskania nie jest rokiem wydania i tak jest traktowany:
			// wchodzi wyłącznie jako data dostępu do zasobu sieciowego, a przy
			// pozycji recenzowanej zostaje brakiem, który zgłosi kontrola.
			if zrodlo.Rodzaj == shared.ResearchSourceKindWeb {
				metadane.rok = strconv.Itoa(*rok)
			}
		}
	}
	return metadane, nil
}

// autorzyCytatuBadania oddaje autorów albo wskazanie ich braku.
func autorzyCytatuBadania(m metadaneCytowaniaBadania) string {
	if len(m.autorzy) == 0 {
		return "[brak autora]"
	}
	if len(m.autorzy) == 1 {
		return m.autorzy[0]
	}
	if len(m.autorzy) == 2 {
		return m.autorzy[0] + " i " + m.autorzy[1]
	}
	return m.autorzy[0] + " i in."
}

// nazwiskoPierwszegoBadania oddaje samo nazwisko pierwszego autora.
func nazwiskoPierwszegoBadania(m metadaneCytowaniaBadania) string {
	if len(m.autorzy) == 0 {
		return "[brak autora]"
	}
	return strings.TrimSpace(strings.Split(m.autorzy[0], ",")[0])
}

// rokCytatuBadania oddaje rok albo znacznik jego braku, przyjęty w stylach.
func rokCytatuBadania(m metadaneCytowaniaBadania) string {
	if strings.TrimSpace(m.rok) == "" {
		return "b.d."
	}
	return m.rok
}

// cytatWTekscieBadania składa odnośnik wstawiany w treść raportu.
func cytatWTekscieBadania(styl string, m metadaneCytowaniaBadania) string {
	switch strings.ToLower(styl) {
	case "mla":
		return "(" + nazwiskoPierwszegoBadania(m) + ")"
	case "ieee", "vancouver":
		return "[" + m.kodZrodla + "]"
	case "chicago":
		return "(" + nazwiskoPierwszegoBadania(m) + " " + rokCytatuBadania(m) + ")"
	default:
		return "(" + nazwiskoPierwszegoBadania(m) + ", " + rokCytatuBadania(m) + ")"
	}
}

// pozycjaBibliograficznaBadania składa pozycję listy literatury.
func pozycjaBibliograficznaBadania(styl string, m metadaneCytowaniaBadania) string {
	ogon := ""
	switch {
	case m.doi != "":
		ogon = " https://doi.org/" + m.doi
	case m.adres != "":
		ogon = " " + m.adres
	}
	wydawca := ""
	if strings.TrimSpace(m.wydawca) != "" {
		wydawca = " " + m.wydawca + "."
	}

	switch strings.ToLower(styl) {
	case "mla":
		return autorzyCytatuBadania(m) + ". \u201e" + m.tytul + ".\u201d" + wydawca + " " +
			rokCytatuBadania(m) + "." + ogon
	case "ieee":
		return autorzyCytatuBadania(m) + ", \u201e" + m.tytul + ",\u201d " +
			strings.TrimSpace(m.wydawca) + ", " + rokCytatuBadania(m) + "." + ogon
	case "vancouver":
		return autorzyCytatuBadania(m) + ". " + m.tytul + ". " +
			strings.TrimSpace(m.wydawca) + "; " + rokCytatuBadania(m) + "." + ogon
	case "chicago":
		return autorzyCytatuBadania(m) + ". " + rokCytatuBadania(m) + ". \u201e" + m.tytul +
			".\u201d" + wydawca + ogon
	default:
		return autorzyCytatuBadania(m) + " (" + rokCytatuBadania(m) + "). " + m.tytul + "." +
			wydawca + ogon
	}
}

// WypiszStyleCytowania obsługuje `research.citation.styles`.
func (a *adapterBadan) WypiszStyleCytowania(ctx context.Context,
	z shared.ResearchCitationStylesRequest) (shared.ResearchCitationStylesResponse, error) {

	style := []shared.ResearchCitationStyle{}
	for _, styl := range styleWbudowaneBadania {
		style = append(style, shared.ResearchCitationStyle{
			Id: styl.kod, Name: styl.nazwa, Custom: false,
		})
	}
	wlasne, err := a.repozytorium.StyleCytowania(ctx)
	if err != nil {
		return shared.ResearchCitationStylesResponse{}, bladBadan(err)
	}
	for _, styl := range wlasne {
		style = append(style, shared.ResearchCitationStyle{
			Id: styl.Kod, Name: styl.Nazwa, Custom: styl.Wlasny,
		})
	}
	if z.Query != nil && strings.TrimSpace(*z.Query) != "" {
		szukane := strings.ToLower(strings.TrimSpace(*z.Query))
		zawezone := style[:0]
		for _, styl := range style {
			if strings.Contains(strings.ToLower(styl.Name), szukane) ||
				strings.Contains(strings.ToLower(styl.Id), szukane) {
				zawezone = append(zawezone, styl)
			}
		}
		style = zawezone
	}
	razem := len(style)
	if z.Limit != nil && *z.Limit > 0 && *z.Limit < razem {
		style = style[:*z.Limit]
	}
	return shared.ResearchCitationStylesResponse{Styles: style, Total: razem}, nil
}

// SprawdzCytowania obsługuje `research.citation.check`.
func (a *adapterBadan) SprawdzCytowania(ctx context.Context,
	z shared.ResearchCitationCheckRequest) (shared.ResearchCitationCheckResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchCitationCheckResponse{}, bladWskazaniaBadan("citation.check bez okna badania")
	}
	zrodla, err := a.repozytorium.Zrodla(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchCitationCheckResponse{}, bladBadan(err)
	}

	uchybienia := []shared.ResearchCitationIssue{}
	for _, zrodlo := range zrodla {
		metadane, err := a.metadaneCytowaniaBadania(ctx, zrodlo.Kod)
		if err != nil {
			return shared.ResearchCitationCheckResponse{}, err
		}
		kod := zrodlo.Kod
		for _, pole := range brakiCytowaniaBadania(metadane) {
			nazwa := pole
			uchybienie := shared.ResearchCitationIssue{
				SourceId: &kod, Kind: "brak_pola", Field: &nazwa,
			}
			if metadane.doi != "" {
				skad := "crossref"
				uchybienie.FixableFrom = &skad
			}
			uchybienia = append(uchybienia, uchybienie)
		}
	}

	// Źródła niecytowane w raporcie: wykaz literatury, której nikt nie użył, jest
	// drugą połową kontroli kompletności — pierwszą są braki metadanych.
	if z.ReportId != nil && strings.TrimSpace(*z.ReportId) != "" {
		uzyte, err := a.zrodlaRaportuBadania(ctx, *z.ReportId)
		if err != nil {
			return shared.ResearchCitationCheckResponse{}, err
		}
		for _, zrodlo := range zrodla {
			if uzyte[zrodlo.Kod] {
				continue
			}
			kod := zrodlo.Kod
			uchybienia = append(uchybienia, shared.ResearchCitationIssue{
				SourceId: &kod, Kind: "zrodlo_niecytowane",
			})
		}
	}
	return shared.ResearchCitationCheckResponse{Issues: uchybienia}, nil
}

// brakiCytowaniaBadania wymienia pola, bez których cytat będzie ułomny.
func brakiCytowaniaBadania(m metadaneCytowaniaBadania) []string {
	braki := []string{}
	if len(m.autorzy) == 0 {
		braki = append(braki, "autor")
	}
	if strings.TrimSpace(m.rok) == "" {
		braki = append(braki, "rok")
	}
	if strings.TrimSpace(m.tytul) == "" {
		braki = append(braki, "tytul")
	}
	if strings.TrimSpace(m.doi) == "" && strings.TrimSpace(m.adres) == "" {
		braki = append(braki, "identyfikator_albo_adres")
	}
	return braki
}

// SprawdzWycofania obsługuje `research.retraction.check`. Pytanie idzie do
// Crossref: praca wycofana ma tam wpis o rodzaju `retraction` powiązany z DOI.
func (a *adapterBadan) SprawdzWycofania(ctx context.Context,
	z shared.ResearchRetractionCheckRequest) (shared.ResearchRetractionCheckResponse, error) {

	if len(z.SourceIds) == 0 {
		return shared.ResearchRetractionCheckResponse{},
			bladWskazaniaBadan("retraction.check bez źródeł — nie ma czego sprawdzać")
	}
	flagi := make([]shared.ResearchRetractionFlag, 0, len(z.SourceIds))
	for _, kod := range z.SourceIds {
		lektura, err := a.repozytorium.LekturaZrodlaBadania(ctx, kod)
		if err != nil {
			return shared.ResearchRetractionCheckResponse{}, bladNieznanegoZrodlaBadania(kod)
		}
		stan := "nieznany"
		var adresNoty *string
		if lektura.Identyfikat != nil && strings.TrimSpace(*lektura.Identyfikat) != "" {
			praca, err := pracaCrossref(ctx, *lektura.Identyfikat)
			if err == nil {
				stan = "bez_zastrzezen"
				if strings.Contains(strings.ToLower(
					pierwszyNapisBadania(praca.Title, "")), "retracted") {
					stan = "wycofana"
					if praca.URL != "" {
						adres := praca.URL
						adresNoty = &adres
					}
				}
			}
		}
		flaga := dane.WycofanieZrodlaBadania{Stan: stan, Adres: adresNoty}
		if err := a.repozytorium.UstawWycofanieZrodla(ctx, kod, flaga); err != nil {
			return shared.ResearchRetractionCheckResponse{}, bladBadan(err)
		}
		zapisana, err := a.repozytorium.WycofanieZrodlaBadania(ctx, kod)
		if err != nil {
			return shared.ResearchRetractionCheckResponse{}, bladBadan(err)
		}
		flagi = append(flagi, shared.ResearchRetractionFlag{
			SourceId: kod, Status: zapisana.Stan, NoticeUrl: zapisana.Adres,
			CheckedAt: chwilaBazy(zapisana.Sprawdzono),
		})
	}
	return shared.ResearchRetractionCheckResponse{Flags: flagi}, nil
}

// zrodlaRaportuBadania oddaje zbiór źródeł, na których stoją ustalenia
// wykorzystane w sekcjach raportu.
func (a *adapterBadan) zrodlaRaportuBadania(ctx context.Context,
	kodRaportu string) (map[string]bool, error) {

	raport, err := a.repozytorium.Raport(ctx, kodRaportu)
	if err != nil {
		return nil, bladNieznanegoRaportuBadania(kodRaportu)
	}
	sekcje, err := a.repozytorium.Sekcje(ctx, raport.ID)
	if err != nil {
		return nil, bladBadan(err)
	}
	uzyte := map[string]bool{}
	for _, sekcja := range sekcje {
		for _, kodUstalenia := range sekcja.UstalenieKody {
			ustalenie, err := a.repozytorium.Ustalenie(ctx, kodUstalenia)
			if err != nil {
				continue
			}
			zrodla, err := a.repozytorium.ZrodlaUstalenia(ctx, ustalenie.ID)
			if err != nil {
				continue
			}
			for _, zrodlo := range zrodla {
				uzyte[zrodlo.Kod] = true
			}
		}
	}
	return uzyte, nil
}

// posortowaneKlucze oddaje klucze zbioru w kolejności ustalonej.
func posortowaneKluczeBadania(zbior map[string]bool) []string {
	klucze := make([]string, 0, len(zbior))
	for klucz := range zbior {
		klucze = append(klucze, klucz)
	}
	sort.Strings(klucze)
	return klucze
}
