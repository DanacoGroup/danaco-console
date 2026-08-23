// Odpowiedzialność pliku: wydanie grafu argumentów — `roundtable.argument.export`
// w sześciu formatach kontraktu: DOT, GraphML, Argdown, AIF, SVG i PNG.
//
// Wszystkie sześć składa rdzeń sam, bez ani jednego programu z zewnątrz.
// Cztery pierwsze są formatami tekstowymi i pisze się je wprost. SVG jest
// dokumentem XML, więc też. PNG powstaje rysowaniem po mapie bitowej
// biblioteką standardową — rasteryzator zewnętrzny byłby zależnością, której
// instalka nie niesie, po to, żeby narysować prostokąty i podpisy.
//
// Układ jest kolumnowy i wynika z treści: węzły stoją w kolumnach według aktu
// mowy (teza, argument, kontrargument, …), więc czytelnik widzi strukturę
// sporu, zanim przeczyta choć jedno zdanie.
package core

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"danacoconsole/shared"
)

const (
	// Wymiary rysunku. Węzeł ma stałą wysokość, a szerokość kolumny wynika
	// z liczby kolumn — graf o dwóch aktach mowy nie ma powodu być wąski.
	szerokoscRysunkuGrafu = 1200
	wysokoscWezlaGrafu    = 64
	odstepWezlowGrafu     = 16
	marginesRysunkuGrafu  = 24
)

// WydajGraf wydaje graf argumentów w formacie wymiany albo jako obraz.
func (a *adapterDebaty) WydajGraf(ctx context.Context,
	z shared.RoundtableArgumentExportRequest) (shared.RoundtableArgumentExportResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableArgumentExportResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	format := strings.TrimSpace(string(z.Format))

	graf, err := a.zlozGraf(ctx, okno, strings.TrimSpace(wartoscTekstu(z.TurnId)), false)
	if err != nil {
		return shared.RoundtableArgumentExportResponse{}, err
	}
	if len(graf.Nodes) == 0 {
		return shared.RoundtableArgumentExportResponse{}, odmowaWydaniaBezGrafu(okno)
	}

	var bajty []byte
	switch format {
	case shared.RoundtableArgumentFormatDot:
		bajty = []byte(grafWDot(graf))
	case shared.RoundtableArgumentFormatGraphml:
		bajty, err = grafWGraphml(graf)
	case shared.RoundtableArgumentFormatArgdown:
		bajty = []byte(grafWArgdown(graf))
	case shared.RoundtableArgumentFormatAif:
		bajty, err = grafWAif(graf)
	case shared.RoundtableArgumentFormatSvg:
		bajty = []byte(grafWSvg(graf))
	case shared.RoundtableArgumentFormatPng:
		bajty, err = grafWPng(graf)
	default:
		return shared.RoundtableArgumentExportResponse{},
			bladWskazaniaDebaty("format grafu " + format + " nie jest formatem znanym kontraktowi")
	}
	if err != nil {
		return shared.RoundtableArgumentExportResponse{}, err
	}

	artefakt, err := a.wydajArtefaktDebaty(ctx, okno, rodzajArtefaktuGrafu, format, bajty, 0)
	if err != nil {
		return shared.RoundtableArgumentExportResponse{}, err
	}
	return shared.RoundtableArgumentExportResponse{
		ArtifactId: artefakt.Kod, Uri: odwolanieArtefaktu(artefakt),
	}, nil
}

// grafWDot zapisuje graf w języku DOT.
func grafWDot(graf shared.RoundtableArgumentGraph) string {
	var zapis strings.Builder
	zapis.WriteString("digraph debata {\n  rankdir=LR;\n  node [shape=box];\n")
	for _, wezel := range graf.Nodes {
		zapis.WriteString("  \"" + wezel.Id + "\" [label=\"" + wCudzyslowieDot(wezel.Text) +
			"\", xlabel=\"" + string(wezel.SpeechAct) + "\"];\n")
	}
	for _, krawedz := range graf.Edges {
		zapis.WriteString("  \"" + krawedz.FromNodeId + "\" -> \"" + krawedz.ToNodeId +
			"\" [label=\"" + string(krawedz.Relation) + "\"];\n")
	}
	zapis.WriteString("}\n")
	return zapis.String()
}

// wCudzyslowieDot chroni znaki, które w DOT kończyłyby etykietę.
func wCudzyslowieDot(tekst string) string {
	zamiennik := strings.NewReplacer(`"`, `\"`, "\n", `\n`, `\`, `\\`)
	return zamiennik.Replace(tekst)
}

// wezelGraphml i krawedzGraphml opisują graf w GraphML. Znaczniki XML składa
// biblioteka standardowa — ręczne sklejanie napisów rozsypałoby się przy
// pierwszym cudzysłowie w treści argumentu.
type wezelGraphml struct {
	XMLName xml.Name      `xml:"node"`
	Id      string        `xml:"id,attr"`
	Dane    []daneGraphml `xml:"data"`
}

type krawedzGraphml struct {
	XMLName xml.Name      `xml:"edge"`
	Zrodlo  string        `xml:"source,attr"`
	Cel     string        `xml:"target,attr"`
	Dane    []daneGraphml `xml:"data"`
}

type daneGraphml struct {
	Klucz   string `xml:"key,attr"`
	Wartosc string `xml:",chardata"`
}

type grafGraphml struct {
	XMLName    xml.Name         `xml:"graphml"`
	Przestrzen string           `xml:"xmlns,attr"`
	Klucze     []kluczGraphml   `xml:"key"`
	Graf       zawartoscGraphml `xml:"graph"`
}

type kluczGraphml struct {
	Id     string `xml:"id,attr"`
	Dla    string `xml:"for,attr"`
	Nazwa  string `xml:"attr.name,attr"`
	Rodzaj string `xml:"attr.type,attr"`
}

type zawartoscGraphml struct {
	Id        string           `xml:"id,attr"`
	Kierunek  string           `xml:"edgedefault,attr"`
	Wezly     []wezelGraphml   `xml:"node"`
	Krawedzie []krawedzGraphml `xml:"edge"`
}

// grafWGraphml zapisuje graf w GraphML.
func grafWGraphml(graf shared.RoundtableArgumentGraph) ([]byte, error) {
	dokument := grafGraphml{
		Przestrzen: "http://graphml.graphdrawing.org/xmlns",
		Klucze: []kluczGraphml{
			{Id: "tresc", Dla: "node", Nazwa: "text", Rodzaj: "string"},
			{Id: "akt", Dla: "node", Nazwa: "speechAct", Rodzaj: "string"},
			{Id: "relacja", Dla: "edge", Nazwa: "relation", Rodzaj: "string"},
		},
		Graf: zawartoscGraphml{Id: graf.WindowId, Kierunek: "directed"},
	}
	for _, wezel := range graf.Nodes {
		dokument.Graf.Wezly = append(dokument.Graf.Wezly, wezelGraphml{
			Id: wezel.Id,
			Dane: []daneGraphml{
				{Klucz: "tresc", Wartosc: wezel.Text},
				{Klucz: "akt", Wartosc: string(wezel.SpeechAct)},
			},
		})
	}
	for _, krawedz := range graf.Edges {
		dokument.Graf.Krawedzie = append(dokument.Graf.Krawedzie, krawedzGraphml{
			Zrodlo: krawedz.FromNodeId, Cel: krawedz.ToNodeId,
			Dane: []daneGraphml{{Klucz: "relacja", Wartosc: string(krawedz.Relation)}},
		})
	}
	bajty, err := xml.MarshalIndent(dokument, "", "  ")
	if err != nil {
		return nil, bladDebaty(err)
	}
	return append([]byte(xml.Header), bajty...), nil
}

// grafWArgdown zapisuje graf w składni Argdown: teza, a pod nią wcięte
// argumenty z plusem (wsparcie) albo minusem (podważenie).
func grafWArgdown(graf shared.RoundtableArgumentGraph) string {
	wychodzace := make(map[string][]shared.RoundtableArgumentEdge, len(graf.Edges))
	wskazywane := make(map[string]struct{}, len(graf.Edges))
	for _, krawedz := range graf.Edges {
		wychodzace[krawedz.ToNodeId] = append(wychodzace[krawedz.ToNodeId], krawedz)
		wskazywane[krawedz.FromNodeId] = struct{}{}
	}

	var zapis strings.Builder
	for _, wezel := range graf.Nodes {
		if _, podrzedny := wskazywane[wezel.Id]; podrzedny {
			continue // węzeł wskazujący inny wypisze się pod nim jako wcięcie
		}
		zapis.WriteString("[" + wezel.Id + "]: " + jednymWierszem(wezel.Text) + "\n")
		for _, krawedz := range wychodzace[wezel.Id] {
			znak := "  +"
			if krawedz.Relation == shared.RoundtableArgumentRelationAttacks {
				znak = "  -"
			}
			zapis.WriteString(znak + " <" + krawedz.FromNodeId + ">: " +
				jednymWierszem(trescWezla(graf, krawedz.FromNodeId)) + "\n")
		}
		zapis.WriteString("\n")
	}
	return zapis.String()
}

// pozycjaAif opisuje węzeł w formacie wymiany argumentów (AIF).
type pozycjaAif struct {
	NodeID   string `json:"nodeID"`
	Text     string `json:"text"`
	Type     string `json:"type"`
	Category string `json:"category,omitempty"`
}

// polaczenieAif opisuje krawędź AIF.
type polaczenieAif struct {
	EdgeID   string `json:"edgeID"`
	FromID   string `json:"fromID"`
	ToID     string `json:"toID"`
	FormEdge string `json:"formEdgeID,omitempty"`
}

// grafWAif zapisuje graf w formacie wymiany argumentów.
//
// Węzeł treści ma typ „I" (information), a relacja typ zależny od jej rodzaju:
// wsparcie „RA" (rule application), podważenie „CA" (conflict application),
// przeformułowanie „MA" (preference/restatement). To jest podział z samego AIF,
// nie nazwa wymyślona tutaj.
func grafWAif(graf shared.RoundtableArgumentGraph) ([]byte, error) {
	dokument := struct {
		Nodes []pozycjaAif    `json:"nodes"`
		Edges []polaczenieAif `json:"edges"`
	}{}
	for _, wezel := range graf.Nodes {
		dokument.Nodes = append(dokument.Nodes, pozycjaAif{
			NodeID: wezel.Id, Text: wezel.Text, Type: "I", Category: string(wezel.SpeechAct),
		})
	}
	for _, krawedz := range graf.Edges {
		typ := "RA"
		switch krawedz.Relation {
		case shared.RoundtableArgumentRelationAttacks:
			typ = "CA"
		case shared.RoundtableArgumentRelationRestates:
			typ = "MA"
		}
		dokument.Nodes = append(dokument.Nodes, pozycjaAif{
			NodeID: krawedz.Id, Text: string(krawedz.Relation), Type: typ,
		})
		dokument.Edges = append(dokument.Edges,
			polaczenieAif{EdgeID: krawedz.Id + "-we", FromID: krawedz.FromNodeId, ToID: krawedz.Id},
			polaczenieAif{EdgeID: krawedz.Id + "-wy", FromID: krawedz.Id, ToID: krawedz.ToNodeId})
	}
	bajty, err := json.MarshalIndent(dokument, "", "  ")
	if err != nil {
		return nil, bladDebaty(err)
	}
	return bajty, nil
}

// ukladGrafu rozstawia węzły w kolumnach według aktu mowy.
type ukladGrafu struct {
	Kolumny    []string
	Pozycje    map[string]struct{ X, Y int }
	Szerokosc  int
	Wysokosc   int
	SzerKolumn int
}

// rozstawGraf liczy położenie każdego węzła.
func rozstawGraf(graf shared.RoundtableArgumentGraph) ukladGrafu {
	kolumny := make([]string, 0, 6)
	wKolumnie := make(map[string][]string, 6)
	for _, wezel := range graf.Nodes {
		akt := string(wezel.SpeechAct)
		if _, jest := wKolumnie[akt]; !jest {
			kolumny = append(kolumny, akt)
		}
		wKolumnie[akt] = append(wKolumnie[akt], wezel.Id)
	}

	szerKolumn := (szerokoscRysunkuGrafu - 2*marginesRysunkuGrafu) / len(kolumny)
	pozycje := make(map[string]struct{ X, Y int }, len(graf.Nodes))
	najwyzsza := 0
	for numer, akt := range kolumny {
		for wiersz, kod := range wKolumnie[akt] {
			y := marginesRysunkuGrafu + 28 +
				wiersz*(wysokoscWezlaGrafu+odstepWezlowGrafu)
			pozycje[kod] = struct{ X, Y int }{
				X: marginesRysunkuGrafu + numer*szerKolumn, Y: y,
			}
			if dol := y + wysokoscWezlaGrafu; dol > najwyzsza {
				najwyzsza = dol
			}
		}
	}
	return ukladGrafu{
		Kolumny: kolumny, Pozycje: pozycje, Szerokosc: szerokoscRysunkuGrafu,
		Wysokosc: najwyzsza + marginesRysunkuGrafu, SzerKolumn: szerKolumn,
	}
}

// grafWSvg rysuje graf jako dokument SVG.
func grafWSvg(graf shared.RoundtableArgumentGraph) string {
	uklad := rozstawGraf(graf)
	var zapis strings.Builder
	zapis.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	zapis.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="` + itoa(uklad.Szerokosc) +
		`" height="` + itoa(uklad.Wysokosc) + `" viewBox="0 0 ` + itoa(uklad.Szerokosc) + ` ` +
		itoa(uklad.Wysokosc) + `">` + "\n")
	zapis.WriteString(`<rect width="100%" height="100%" fill="#ffffff"/>` + "\n")

	for numer, akt := range uklad.Kolumny {
		x := marginesRysunkuGrafu + numer*uklad.SzerKolumn
		zapis.WriteString(`<text x="` + itoa(x) + `" y="` + itoa(marginesRysunkuGrafu) +
			`" font-family="sans-serif" font-size="14" fill="#333333">` +
			wXml(akt) + `</text>` + "\n")
	}
	for _, krawedz := range graf.Edges {
		od, jestOd := uklad.Pozycje[krawedz.FromNodeId]
		do, jestDo := uklad.Pozycje[krawedz.ToNodeId]
		if !jestOd || !jestDo {
			continue
		}
		kolor := "#3a7bd5"
		if krawedz.Relation == shared.RoundtableArgumentRelationAttacks {
			kolor = "#d64545"
		}
		zapis.WriteString(`<line x1="` + itoa(od.X+uklad.SzerKolumn-32) +
			`" y1="` + itoa(od.Y+wysokoscWezlaGrafu/2) +
			`" x2="` + itoa(do.X) + `" y2="` + itoa(do.Y+wysokoscWezlaGrafu/2) +
			`" stroke="` + kolor + `" stroke-width="2"/>` + "\n")
	}
	for _, wezel := range graf.Nodes {
		pozycja := uklad.Pozycje[wezel.Id]
		zapis.WriteString(`<rect x="` + itoa(pozycja.X) + `" y="` + itoa(pozycja.Y) +
			`" width="` + itoa(uklad.SzerKolumn-32) + `" height="` + itoa(wysokoscWezlaGrafu) +
			`" rx="6" fill="#f4f6fb" stroke="#8a93a6"/>` + "\n")
		for wiersz, tekst := range podzielNaWierszeGrafu(wezel.Text, 3, (uklad.SzerKolumn-48)/7) {
			zapis.WriteString(`<text x="` + itoa(pozycja.X+8) + `" y="` +
				itoa(pozycja.Y+18+wiersz*16) +
				`" font-family="sans-serif" font-size="12" fill="#1b1f2a">` +
				wXml(tekst) + `</text>` + "\n")
		}
	}
	zapis.WriteString("</svg>\n")
	return zapis.String()
}

// grafWPng rysuje ten sam układ po mapie bitowej.
func grafWPng(graf shared.RoundtableArgumentGraph) ([]byte, error) {
	uklad := rozstawGraf(graf)
	plotno := image.NewRGBA(image.Rect(0, 0, uklad.Szerokosc, uklad.Wysokosc))
	draw.Draw(plotno, plotno.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	obramowanie := color.RGBA{R: 0x8a, G: 0x93, B: 0xa6, A: 0xff}
	wypelnienie := color.RGBA{R: 0xf4, G: 0xf6, B: 0xfb, A: 0xff}
	wsparcie := color.RGBA{R: 0x3a, G: 0x7b, B: 0xd5, A: 0xff}
	podwazenie := color.RGBA{R: 0xd6, G: 0x45, B: 0x45, A: 0xff}
	napis := color.RGBA{R: 0x1b, G: 0x1f, B: 0x2a, A: 0xff}

	for _, krawedz := range graf.Edges {
		od, jestOd := uklad.Pozycje[krawedz.FromNodeId]
		do, jestDo := uklad.Pozycje[krawedz.ToNodeId]
		if !jestOd || !jestDo {
			continue
		}
		kolor := wsparcie
		if krawedz.Relation == shared.RoundtableArgumentRelationAttacks {
			kolor = podwazenie
		}
		narysujOdcinek(plotno, od.X+uklad.SzerKolumn-32, od.Y+wysokoscWezlaGrafu/2,
			do.X, do.Y+wysokoscWezlaGrafu/2, kolor)
	}
	for _, wezel := range graf.Nodes {
		pozycja := uklad.Pozycje[wezel.Id]
		prostokat := image.Rect(pozycja.X, pozycja.Y,
			pozycja.X+uklad.SzerKolumn-32, pozycja.Y+wysokoscWezlaGrafu)
		draw.Draw(plotno, prostokat, &image.Uniform{C: wypelnienie}, image.Point{}, draw.Src)
		narysujRamke(plotno, prostokat, obramowanie)
		for wiersz, tekst := range podzielNaWierszeGrafu(wezel.Text, 3, (uklad.SzerKolumn-48)/7) {
			napiszTekst(plotno, pozycja.X+8, pozycja.Y+18+wiersz*16, tekst, napis)
		}
	}
	for numer, akt := range uklad.Kolumny {
		napiszTekst(plotno, marginesRysunkuGrafu+numer*uklad.SzerKolumn,
			marginesRysunkuGrafu, akt, napis)
	}

	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		return nil, bladDebaty(err)
	}
	return bufor.Bytes(), nil
}

// narysujOdcinek kreśli linię prostą między dwoma punktami.
func narysujOdcinek(plotno *image.RGBA, x1, y1, x2, y2 int, kolor color.Color) {
	kroki := abs(x2-x1) + abs(y2-y1)
	if kroki == 0 {
		return
	}
	for krok := 0; krok <= kroki; krok++ {
		x := x1 + (x2-x1)*krok/kroki
		y := y1 + (y2-y1)*krok/kroki
		plotno.Set(x, y, kolor)
	}
}

// narysujRamke kreśli obwód prostokąta.
func narysujRamke(plotno *image.RGBA, prostokat image.Rectangle, kolor color.Color) {
	for x := prostokat.Min.X; x < prostokat.Max.X; x++ {
		plotno.Set(x, prostokat.Min.Y, kolor)
		plotno.Set(x, prostokat.Max.Y-1, kolor)
	}
	for y := prostokat.Min.Y; y < prostokat.Max.Y; y++ {
		plotno.Set(prostokat.Min.X, y, kolor)
		plotno.Set(prostokat.Max.X-1, y, kolor)
	}
}

// napiszTekst kładzie napis na mapie bitowej fontem wkompilowanym w bibliotekę.
func napiszTekst(plotno *image.RGBA, x, y int, tekst string, kolor color.Color) {
	rysownik := &font.Drawer{
		Dst: plotno, Src: &image.Uniform{C: kolor}, Face: basicfont.Face7x13,
		Dot: fixed.P(x, y),
	}
	rysownik.DrawString(tekst)
}

// podzielNaWierszeGrafu przycina treść węzła do rysunku: najwyżej tyle wierszy,
// ile mieści prostokąt, ostatni zakończony wielokropkiem, gdy treść się urywa.
func podzielNaWierszeGrafu(tekst string, ileWierszy, znakow int) []string {
	if znakow < 8 {
		znakow = 8
	}
	wiersze := zawinWiersze(jednymWierszem(tekst), znakow)
	if len(wiersze) <= ileWierszy {
		return wiersze
	}
	przyciete := wiersze[:ileWierszy]
	przyciete[ileWierszy-1] += "…"
	return przyciete
}

// jednymWierszem zbija treść do jednego wiersza — łamanie wierszy w etykiecie
// węzła rozsypałoby zarówno DOT, jak i rysunek.
func jednymWierszem(tekst string) string {
	return strings.Join(strings.Fields(tekst), " ")
}

// trescWezla odnajduje treść węzła po kodzie.
func trescWezla(graf shared.RoundtableArgumentGraph, kod string) string {
	for _, wezel := range graf.Nodes {
		if wezel.Id == kod {
			return wezel.Text
		}
	}
	return kod
}

// wXml chroni znaki, które w dokumencie XML mają znaczenie składniowe.
func wXml(tekst string) string {
	var bufor bytes.Buffer
	_ = xml.EscapeText(&bufor, []byte(tekst))
	return bufor.String()
}

// abs oddaje wartość bezwzględną liczby całkowitej.
func abs(wartosc int) int {
	if wartosc < 0 {
		return -wartosc
	}
	return wartosc
}
