// Moduł Library — słownik etykiet, tezaurus i wykaz kolekcji:
// `library.tag.list`, `library.tag.update`, `library.tag.merge`,
// `library.tag.remove`, `library.collection.list`,
// `library.thesaurus.relate`, `library.thesaurus.export`.
//
// Etykieta ma tu tożsamość, której `library.tag.set` jej nie daje: barwę, czas
// założenia i istnienie niezależne od tego, czy nosi ją jakikolwiek zasób.
// Dzięki temu okno Tags & Collections pokazuje słownik, a nie próbkę zebraną
// z odczytanej strony wykazu.
//
// Wywóz tezaurusa idzie w SKOS/RDF w trzech serializacjach i powstaje w rdzeniu,
// bo relacje leżą w bazie i klient nie ma ich skąd wziąć. Zapis składa się
// z tekstu — żadnego programu zewnętrznego: RDF w postaci Turtle, RDF/XML
// i JSON-LD to formaty tekstowe o znanym kształcie.
package core

import (
	"context"
	"danacoconsole/server/internal/dane"
	"encoding/json"
	"encoding/xml"
	"sort"
	"strings"

	"danacoconsole/shared"
)

// SlownikEtykiet obsługuje `library.tag.list`.
func (a *adapterBiblioteki) SlownikEtykiet(ctx context.Context,
	z shared.LibraryTagListRequest) (shared.LibraryTagListResponse, error) {

	granica := 0
	if z.Limit != nil {
		granica = *z.Limit
	}
	tylkoNieuzywane := z.UnusedOnly != nil && *z.UnusedOnly
	wiersze, lacznie, err := a.repozytorium.EtykietySlownika(ctx, z.Query, tylkoNieuzywane, granica)
	if err != nil {
		return shared.LibraryTagListResponse{}, bladBiblioteki(err)
	}
	etykiety := make([]shared.LibraryTag, 0, len(wiersze))
	for _, wiersz := range wiersze {
		etykiety = append(etykiety, shared.LibraryTag{
			Name: wiersz.Nazwa, Color: wiersz.Barwa, FileCount: wiersz.LiczbaPlikow,
			CreatedAt: chwilaBazy(wiersz.Utworzono),
		})
	}
	return shared.LibraryTagListResponse{Tags: etykiety, Total: lacznie}, nil
}

// ZmienEtykiete obsługuje `library.tag.update`.
//
// Zmiana nazwy przechodzi na wszystkie zasoby noszące etykietę — inaczej
// powstałaby druga etykieta o tym samym znaczeniu, a pierwsza zostałaby przy
// zasobach jako sierota.
func (a *adapterBiblioteki) ZmienEtykiete(ctx context.Context,
	z shared.LibraryTagUpdateRequest) (shared.LibraryTagUpdateResponse, error) {

	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.LibraryTagUpdateResponse{}, bladWskazaniaBiblioteki("komenda bez wskazania etykiety")
	}
	dotkniete := 0
	docelowa := nazwa
	if z.NewName != nil && strings.TrimSpace(*z.NewName) != "" &&
		strings.TrimSpace(*z.NewName) != nazwa {

		docelowa = strings.TrimSpace(*z.NewName)
		// Wpis słownika zakłada się przed zmianą nazwy: etykieta nadana przy
		// zasobie mogła nigdy nie mieć wiersza słownikowego, a bez niego zmiana
		// nazwy przeszłaby po zasobach i zgubiła barwę.
		if _, err := a.repozytorium.ZapiszEtykieteSlownika(ctx, nazwa, nil); err != nil {
			return shared.LibraryTagUpdateResponse{}, bladBiblioteki(err)
		}
		liczba, err := a.repozytorium.PrzemianujEtykiete(ctx, nazwa, docelowa)
		if err != nil {
			return shared.LibraryTagUpdateResponse{}, bladBiblioteki(err)
		}
		dotkniete = liczba
	}
	if z.Color != nil || docelowa == nazwa {
		if _, err := a.repozytorium.ZapiszEtykieteSlownika(ctx, docelowa, z.Color); err != nil {
			return shared.LibraryTagUpdateResponse{}, bladBiblioteki(err)
		}
	}
	etykieta, err := a.repozytorium.EtykietaSlownika(ctx, docelowa)
	if err != nil {
		return shared.LibraryTagUpdateResponse{}, bladBiblioteki(err)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "zmiana etykiety "+nazwa)
	return shared.LibraryTagUpdateResponse{
		Tag: shared.LibraryTag{
			Name: etykieta.Nazwa, Color: etykieta.Barwa, FileCount: etykieta.LiczbaPlikow,
			CreatedAt: chwilaBazy(etykieta.Utworzono),
		},
		AffectedFiles: dotkniete,
	}, nil
}

// PolaczEtykiety obsługuje `library.tag.merge`.
//
// Łączenie jest zmianą nazwy wykonaną wielokrotnie: zasób noszący etykietę
// źródłową dostaje docelową, a źródłowa znika ze słownika. Etykieta wskazana
// jako źródłowa i docelowa naraz jest pomijana — wchłonięcie siebie samej
// zdjęłoby etykietę z zasobów bez powodu.
func (a *adapterBiblioteki) PolaczEtykiety(ctx context.Context,
	z shared.LibraryTagMergeRequest) (shared.LibraryTagMergeResponse, error) {

	docelowa := strings.TrimSpace(z.TargetName)
	if docelowa == "" {
		return shared.LibraryTagMergeResponse{}, bladWskazaniaBiblioteki(
			"łączenie bez wskazania etykiety docelowej")
	}
	if len(z.SourceNames) == 0 {
		return shared.LibraryTagMergeResponse{}, bladWskazaniaBiblioteki(
			"łączenie bez wskazania etykiet wchłanianych")
	}
	if _, err := a.repozytorium.ZapiszEtykieteSlownika(ctx, docelowa, nil); err != nil {
		return shared.LibraryTagMergeResponse{}, bladBiblioteki(err)
	}

	dotkniete := 0
	wchloniete := []string{}
	for _, zrodlowa := range z.SourceNames {
		nazwa := strings.TrimSpace(zrodlowa)
		if nazwa == "" || nazwa == docelowa {
			continue
		}
		liczba, err := a.repozytorium.PrzemianujEtykiete(ctx, nazwa, docelowa)
		if err != nil {
			return shared.LibraryTagMergeResponse{}, bladBiblioteki(err)
		}
		// Wiersz słownika etykiety źródłowej znika razem z nią: zmiana nazwy
		// mogła go nie objąć, gdy etykieta docelowa miała już swój wpis.
		if _, err := a.repozytorium.UsunEtykieteZeSlownika(ctx, nazwa); err != nil {
			return shared.LibraryTagMergeResponse{}, bladBiblioteki(err)
		}
		dotkniete += liczba
		wchloniete = append(wchloniete, nazwa)
	}
	etykieta, err := a.repozytorium.EtykietaSlownika(ctx, docelowa)
	if err != nil {
		return shared.LibraryTagMergeResponse{}, bladBiblioteki(err)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil,
		"połączenie etykiet w "+docelowa+": "+strings.Join(wchloniete, ", "))
	return shared.LibraryTagMergeResponse{
		Tag: shared.LibraryTag{
			Name: etykieta.Nazwa, Color: etykieta.Barwa, FileCount: etykieta.LiczbaPlikow,
			CreatedAt: chwilaBazy(etykieta.Utworzono),
		},
		AffectedFiles: dotkniete, MergedNames: wchloniete,
	}, nil
}

// UsunEtykiete obsługuje `library.tag.remove`.
//
// Usunięcie etykiety noszonej przez zasoby żąda potwierdzenia, bo zdejmuje ją
// z zasobów, których żądanie nie wymienia. Etykieta nieużywana schodzi bez
// pytania — nie ma czego stracić.
func (a *adapterBiblioteki) UsunEtykiete(ctx context.Context,
	z shared.LibraryTagRemoveRequest) (shared.LibraryTagRemoveResponse, error) {

	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.LibraryTagRemoveResponse{}, bladWskazaniaBiblioteki("komenda bez wskazania etykiety")
	}
	etykieta, err := a.repozytorium.EtykietaSlownika(ctx, nazwa)
	if err != nil {
		return shared.LibraryTagRemoveResponse{}, bladNieznanejEtykietyBiblioteki(nazwa, err)
	}
	potwierdzone := z.Confirm != nil && *z.Confirm
	if etykieta.LiczbaPlikow > 0 && !potwierdzone {
		return shared.LibraryTagRemoveResponse{}, bladWskazaniaBiblioteki(
			"etykieta " + nazwa + " stoi przy zasobach; usunięcie zdejmie ją z nich — " +
				"potwierdź żądanie polem confirm")
	}
	zdjete, err := a.repozytorium.UsunEtykieteZeSlownika(ctx, nazwa)
	if err != nil {
		return shared.LibraryTagRemoveResponse{}, bladBiblioteki(err)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "usunięcie etykiety "+nazwa)
	return shared.LibraryTagRemoveResponse{Removed: true, AffectedFiles: zdjete}, nil
}

// WykazKolekcji obsługuje `library.collection.list`.
func (a *adapterBiblioteki) WykazKolekcji(ctx context.Context,
	z shared.LibraryCollectionListRequest) (shared.LibraryCollectionListResponse, error) {

	granica := 0
	if z.Limit != nil {
		granica = *z.Limit
	}
	wiersze, lacznie, err := a.repozytorium.KolekcjeWykaz(ctx, z.ParentId, z.Query, granica)
	if err != nil {
		return shared.LibraryCollectionListResponse{}, bladBiblioteki(err)
	}
	kolekcje := make([]shared.LibraryCollection, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kolekcje = append(kolekcje, shared.LibraryCollection{
			Id: wiersz.Kod, Name: wiersz.Nazwa, Description: wiersz.Opis,
			ParentId: wiersz.RodzicKod, RuleId: wiersz.RegulaKod,
			FileCount: wiersz.LiczbaPlikow,
			CreatedAt: chwilaBazy(wiersz.Utworzono), UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
		})
	}
	return shared.LibraryCollectionListResponse{Collections: kolekcje, Total: lacznie}, nil
}

// UstawRelacjeTezaurusa obsługuje `library.thesaurus.relate`.
//
// Relacja zapisuje się w jednym kierunku, tym wskazanym przez Operatora.
// Odwrotność wyprowadza odczyt (`nadrzędna` czytana od drugiej strony jest
// `podrzędną`) — zapis obu kierunków dałby dwa wiersze mówiące to samo i rozjazd
// przy zdjęciu jednego z nich.
func (a *adapterBiblioteki) UstawRelacjeTezaurusa(ctx context.Context,
	z shared.LibraryThesaurusRelateRequest) (shared.LibraryThesaurusRelateResponse, error) {

	zrodlo := strings.TrimSpace(z.SourceName)
	cel := strings.TrimSpace(z.TargetName)
	if zrodlo == "" || cel == "" {
		return shared.LibraryThesaurusRelateResponse{}, bladWskazaniaBiblioteki(
			"relacja tezaurusa bez wskazania obu etykiet")
	}
	if zrodlo == cel {
		return shared.LibraryThesaurusRelateResponse{}, bladWskazaniaBiblioteki(
			"etykieta " + zrodlo + " nie może stać w relacji sama ze sobą")
	}
	zdejmij := z.Remove != nil && *z.Remove
	if !zdejmij {
		// Obie etykiety wchodzą do słownika: relacja między pojęciami, z których
		// jedno nie istnieje, byłaby krawędzią donikąd.
		if _, err := a.repozytorium.ZapiszEtykieteSlownika(ctx, zrodlo, nil); err != nil {
			return shared.LibraryThesaurusRelateResponse{}, bladBiblioteki(err)
		}
		if _, err := a.repozytorium.ZapiszEtykieteSlownika(ctx, cel, nil); err != nil {
			return shared.LibraryThesaurusRelateResponse{}, bladBiblioteki(err)
		}
	}
	stoi, err := a.repozytorium.UstawRelacjeTezaurusa(ctx, zrodlo, cel,
		relacjaTezaurusaBazy(z.Relation), zdejmij)
	if err != nil {
		return shared.LibraryThesaurusRelateResponse{}, bladBiblioteki(err)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil,
		"relacja tezaurusa "+zrodlo+" → "+cel)
	return shared.LibraryThesaurusRelateResponse{RelationSet: stoi}, nil
}

// WywiezTezaurus obsługuje `library.thesaurus.export`.
func (a *adapterBiblioteki) WywiezTezaurus(ctx context.Context,
	z shared.LibraryThesaurusExportRequest) (shared.LibraryThesaurusExportResponse, error) {

	etykiety, _, err := a.repozytorium.EtykietySlownika(ctx, nil, false, 0)
	if err != nil {
		return shared.LibraryThesaurusExportResponse{}, bladBiblioteki(err)
	}
	relacje, err := a.repozytorium.RelacjeTezaurusa(ctx)
	if err != nil {
		return shared.LibraryThesaurusExportResponse{}, bladBiblioteki(err)
	}
	pojecia := make([]string, 0, len(etykiety))
	for _, etykieta := range etykiety {
		pojecia = append(pojecia, etykieta.Nazwa)
	}
	sort.Strings(pojecia)

	kolekcje := []shared.LibraryCollection{}
	if z.IncludeCollections != nil && *z.IncludeCollections {
		wykaz, err := a.WykazKolekcji(ctx, shared.LibraryCollectionListRequest{})
		if err != nil {
			return shared.LibraryThesaurusExportResponse{}, err
		}
		kolekcje = wykaz.Collections
	}

	postac := strings.ToLower(strings.TrimSpace(wartoscTekstu(z.Format)))
	if postac == "" {
		postac = "turtle"
	}
	var tresc string
	switch postac {
	case "turtle":
		tresc = tezaurusTurtleBiblioteki(pojecia, relacje, kolekcje)
	case "rdfxml":
		tresc, err = tezaurusRdfXmlBiblioteki(pojecia, relacje, kolekcje)
	case "jsonld":
		tresc, err = tezaurusJsonLdBiblioteki(pojecia, relacje, kolekcje)
	default:
		return shared.LibraryThesaurusExportResponse{}, bladWskazaniaBiblioteki(
			"serializacja " + postac + " nie jest znana; wywóz idzie w turtle, rdfxml albo jsonld")
	}
	if err != nil {
		return shared.LibraryThesaurusExportResponse{}, bladBiblioteki(err)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionExport, nil, "wywóz tezaurusa w postaci "+postac)
	return shared.LibraryThesaurusExportResponse{
		Content: tresc, Format: postac,
		ConceptCount: len(pojecia), RelationCount: len(relacje),
	}, nil
}

// przestrzenTezaurusaBiblioteki jest przestrzenią nazw pojęć repozytorium.
// Adres nie prowadzi do zasobu w sieci i prowadzić nie musi: w RDF przestrzeń
// nazw jest identyfikatorem, nie odnośnikiem do pobrania.
const przestrzenTezaurusaBiblioteki = "https://danaco-group.pl/biblioteka/tezaurus#"

// orzeczenieSkosBiblioteki przekłada rodzaj relacji na orzeczenie SKOS.
func orzeczenieSkosBiblioteki(rodzaj string) string {
	switch rodzaj {
	case "nadrzedna":
		return "skos:broader"
	case "podrzedna":
		return "skos:narrower"
	default:
		return "skos:related"
	}
}

// tezaurusTurtleBiblioteki składa wywóz w serializacji Turtle.
func tezaurusTurtleBiblioteki(pojecia []string, relacje []dane.RelacjaTezaurusaBiblioteki,
	kolekcje []shared.LibraryCollection) string {

	var zapis strings.Builder
	zapis.WriteString("@prefix skos: <http://www.w3.org/2004/02/skos/core#> .\n")
	zapis.WriteString("@prefix dnc: <" + przestrzenTezaurusaBiblioteki + "> .\n\n")
	for _, pojecie := range pojecia {
		zapis.WriteString("dnc:" + kodPojeciaTezaurusa(pojecie) + " a skos:Concept ;\n")
		zapis.WriteString("    skos:prefLabel \"" + ucieczkaTekstuTezaurusa(pojecie) + "\" .\n")
	}
	for _, relacja := range relacje {
		zapis.WriteString("dnc:" + kodPojeciaTezaurusa(relacja.Zrodlo) + " " +
			orzeczenieSkosBiblioteki(relacja.Rodzaj) + " dnc:" + kodPojeciaTezaurusa(relacja.Cel) + " .\n")
	}
	for _, kolekcja := range kolekcje {
		zapis.WriteString("dnc:kolekcja-" + kodPojeciaTezaurusa(kolekcja.Id) + " a skos:Collection ;\n")
		zapis.WriteString("    skos:prefLabel \"" + ucieczkaTekstuTezaurusa(kolekcja.Name) + "\" .\n")
	}
	return zapis.String()
}

// tezaurusRdfXmlBiblioteki składa wywóz w serializacji RDF/XML.
func tezaurusRdfXmlBiblioteki(pojecia []string, relacje []dane.RelacjaTezaurusaBiblioteki,
	kolekcje []shared.LibraryCollection) (string, error) {

	var zapis strings.Builder
	zapis.WriteString(xml.Header)
	zapis.WriteString(`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"` + "\n")
	zapis.WriteString(`         xmlns:skos="http://www.w3.org/2004/02/skos/core#">` + "\n")
	for _, pojecie := range pojecia {
		zapis.WriteString(`  <skos:Concept rdf:about="` + przestrzenTezaurusaBiblioteki +
			kodPojeciaTezaurusa(pojecie) + `">` + "\n")
		zapis.WriteString("    <skos:prefLabel>" + ucieczkaXmlTezaurusa(pojecie) + "</skos:prefLabel>\n")
		for _, relacja := range relacje {
			if relacja.Zrodlo != pojecie {
				continue
			}
			nazwa := strings.TrimPrefix(orzeczenieSkosBiblioteki(relacja.Rodzaj), "skos:")
			zapis.WriteString(`    <skos:` + nazwa + ` rdf:resource="` +
				przestrzenTezaurusaBiblioteki + kodPojeciaTezaurusa(relacja.Cel) + `"/>` + "\n")
		}
		zapis.WriteString("  </skos:Concept>\n")
	}
	for _, kolekcja := range kolekcje {
		zapis.WriteString(`  <skos:Collection rdf:about="` + przestrzenTezaurusaBiblioteki +
			"kolekcja-" + kodPojeciaTezaurusa(kolekcja.Id) + `">` + "\n")
		zapis.WriteString("    <skos:prefLabel>" + ucieczkaXmlTezaurusa(kolekcja.Name) +
			"</skos:prefLabel>\n")
		zapis.WriteString("  </skos:Collection>\n")
	}
	zapis.WriteString("</rdf:RDF>\n")
	return zapis.String(), nil
}

// tezaurusJsonLdBiblioteki składa wywóz w serializacji JSON-LD.
func tezaurusJsonLdBiblioteki(pojecia []string, relacje []dane.RelacjaTezaurusaBiblioteki,
	kolekcje []shared.LibraryCollection) (string, error) {

	graf := make([]map[string]any, 0, len(pojecia)+len(kolekcje))
	for _, pojecie := range pojecia {
		wpis := map[string]any{
			"@id":            "dnc:" + kodPojeciaTezaurusa(pojecie),
			"@type":          "skos:Concept",
			"skos:prefLabel": pojecie,
		}
		for _, relacja := range relacje {
			if relacja.Zrodlo != pojecie {
				continue
			}
			orzeczenie := orzeczenieSkosBiblioteki(relacja.Rodzaj)
			cele, _ := wpis[orzeczenie].([]string)
			wpis[orzeczenie] = append(cele, "dnc:"+kodPojeciaTezaurusa(relacja.Cel))
		}
		graf = append(graf, wpis)
	}
	for _, kolekcja := range kolekcje {
		graf = append(graf, map[string]any{
			"@id":            "dnc:kolekcja-" + kodPojeciaTezaurusa(kolekcja.Id),
			"@type":          "skos:Collection",
			"skos:prefLabel": kolekcja.Name,
		})
	}
	dokument := map[string]any{
		"@context": map[string]string{
			"skos": "http://www.w3.org/2004/02/skos/core#",
			"dnc":  przestrzenTezaurusaBiblioteki,
		},
		"@graph": graf,
	}
	tresc, err := json.MarshalIndent(dokument, "", "  ")
	if err != nil {
		return "", err
	}
	return string(tresc), nil
}

// kodPojeciaTezaurusa zamienia nazwę etykiety na człon identyfikatora: znaki
// spoza zakresu bezpiecznego ustępują myślnikowi, żeby wywóz był czytelny dla
// każdego czytnika RDF.
func kodPojeciaTezaurusa(nazwa string) string {
	var zapis strings.Builder
	for _, znak := range nazwa {
		switch {
		case znak >= 'a' && znak <= 'z', znak >= 'A' && znak <= 'Z',
			znak >= '0' && znak <= '9', znak == '-', znak == '_':
			zapis.WriteRune(znak)
		default:
			zapis.WriteRune('-')
		}
	}
	kod := zapis.String()
	if kod == "" {
		return "pojecie"
	}
	return kod
}

// ucieczkaTekstuTezaurusa zabezpiecza treść etykiety w zapisie Turtle.
func ucieczkaTekstuTezaurusa(tekst string) string {
	zamiana := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\r", `\r`)
	return zamiana.Replace(tekst)
}

// ucieczkaXmlTezaurusa zabezpiecza treść etykiety w zapisie XML.
func ucieczkaXmlTezaurusa(tekst string) string {
	var zapis strings.Builder
	_ = xml.EscapeText(&zapis, []byte(tekst))
	return zapis.String()
}

// relacjaTezaurusaBazy przekłada wyliczenie relacji na wartość kolumny
// (odwzorowanie: `relacja_tezaurusa_biblioteki.rodzaj`).
func relacjaTezaurusaBazy(relacja shared.LibraryThesaurusRelation) string {
	switch relacja {
	case shared.LibraryThesaurusRelationBroader:
		return "nadrzedna"
	case shared.LibraryThesaurusRelationNarrower:
		return "podrzedna"
	default:
		return "pokrewna"
	}
}
