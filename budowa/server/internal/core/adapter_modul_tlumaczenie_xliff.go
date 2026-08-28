// Plik obsługuje standard wymiany XLIFF w module Translate: `translate.xliff.import` oraz skład pliku XLIFF na potrzeby pakietu przekazania (`handoff.build`).
package core

import (
	"context"
	"encoding/xml"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// jednostkaXliff to jedna para źródło-przekład niezależna od wersji standardu XLIFF, wspólna dla obu formatów odczytu.
type jednostkaXliff struct {
	Klucz  string
	Zrodlo string
	Cel    string
}

// plikXliff niesie wynik odczytu pliku XLIFF: język docelowy oraz wykaz jego jednostek źródło-przekład.
type plikXliff struct {
	JezykZrodlowy string
	JezykDocelowy string
	Jednostki     []jednostkaXliff
}

// WczytajXliff obsługuje `translate.xliff.import`, zakładając panel języka docelowego dla każdej jednostki pliku.
func (a *adapterTlumaczenia) WczytajXliff(ctx context.Context,
	z shared.TranslateXliffImportRequest) (shared.TranslateXliffImportResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateXliffImportResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	sciezka := strings.TrimSpace(z.Path)
	if sciezka == "" {
		return shared.TranslateXliffImportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.xliff.import bez ścieżki pliku")
	}
	bajty, err := os.ReadFile(filepath.Clean(sciezka))
	if err != nil {
		return shared.TranslateXliffImportResponse{}, bladPlikuTlumaczenia(sciezka, err)
	}
	odczytany := rozbierzXliff(bajty)
	if len(odczytany.Jednostki) == 0 {
		return shared.TranslateXliffImportResponse{}, bladWskazaniaTlumaczenia(
			"plik " + filepath.Base(sciezka) + " nie niesie ani jednej jednostki tłumaczeniowej")
	}

	uwagi := []string{}
	if odczytany.JezykDocelowy == "" {
		return shared.TranslateXliffImportResponse{}, bladWskazaniaTlumaczenia(
			"plik nie wskazuje języka docelowego — rdzeń nie zgaduje, do którego panelu wnieść przekład")
	}

	// Tekst źródłowy okna bierze się z jednostek: bez niego panel nie ma z czym zestawić przekładu.
	zrodlowe := make([]string, 0, len(odczytany.Jednostki))
	docelowe := make([]string, 0, len(odczytany.Jednostki))
	pustych := 0
	for _, jednostka := range odczytany.Jednostki {
		zrodlowe = append(zrodlowe, jednostka.Zrodlo)
		if strings.TrimSpace(jednostka.Cel) == "" {
			pustych++
			docelowe = append(docelowe, jednostka.Zrodlo)
			continue
		}
		docelowe = append(docelowe, jednostka.Cel)
	}
	if pustych > 0 {
		uwagi = append(uwagi, "jednostek bez przekładu: "+strconv.Itoa(pustych)+
			" — w panelu zostały w brzmieniu źródłowym")
	}

	tekstZrodlowy := strings.Join(zrodlowe, "\n\n")
	liczba := int64(len(zrodlowe))
	jezykZrodlowy := okno.JezykZrodlowy
	if odczytany.JezykZrodlowy != "" {
		jezykZrodlowy = &odczytany.JezykZrodlowy
	}
	if _, err := a.repozytorium.ZapiszOkno(ctx, dane.OknoTlumaczenia{
		Kod:             okno.Kod,
		TekstZrodlowy:   &tekstZrodlowy,
		JezykZrodlowy:   jezykZrodlowy,
		LiczbaSegmentow: &liczba,
	}); err != nil {
		return shared.TranslateXliffImportResponse{}, bladTlumaczenia(err)
	}
	if err := a.repozytorium.UstawSegmentyOkna(ctx, okno.ID, zrodlowe); err != nil {
		return shared.TranslateXliffImportResponse{}, bladTlumaczenia(err)
	}

	tresc := strings.Join(docelowe, "\n\n")
	panel, err := a.panelJezyka(ctx, okno, odczytany.JezykDocelowy,
		z.Overwrite != nil && *z.Overwrite, tresc)
	if err != nil {
		return shared.TranslateXliffImportResponse{}, err
	}
	if panel == nil {
		uwagi = append(uwagi, "panel języka "+odczytany.JezykDocelowy+
			" miał już treść, a żądanie nie zezwoliło na nadpisanie")
	}

	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateXliffImportResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateXliffImportResponse{
		Panels:        zlozPaneleTlumaczenia(panele),
		ImportedCount: len(odczytany.Jednostki),
		Notes:         uwagi,
	}, nil
}

// panelJezyka wnosi treść do panelu wskazanego języka: zakłada panel, gdy okno
// go jeszcze nie ma, a zastany nadpisuje wyłącznie za zgodą żądania. Wynik
// pusty (bez błędu) znaczy „panel zastany zostawiono nietknięty”.
func (a *adapterTlumaczenia) panelJezyka(ctx context.Context, okno dane.OknoTlumaczenia,
	jezyk string, nadpisuj bool, tresc string) (*dane.PanelTlumaczenia, error) {

	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return nil, bladTlumaczenia(err)
	}
	for _, panel := range panele {
		if panel.Jezyk != jezyk {
			continue
		}
		if panel.Tresc != nil && strings.TrimSpace(*panel.Tresc) != "" && !nadpisuj {
			return nil, nil
		}
		zmieniony, err := a.repozytorium.UstawTlumaczenie(ctx, panel.Kod, &tresc, nil)
		if err != nil {
			return nil, bladTlumaczenia(err)
		}
		a.rozglosZmianePanelu(shared.ChangeKindUpdated, zmieniony)
		return &zmieniony, nil
	}
	zalozony, err := a.repozytorium.ZapiszPanel(ctx, okno.ID, dane.PanelTlumaczenia{
		Kod:   nowyIdentyfikator(przedrostekPaneluTlumaczenia),
		Jezyk: jezyk,
		Tresc: &tresc,
	})
	if err != nil {
		return nil, bladTlumaczenia(err)
	}
	a.rozglosZmianePanelu(shared.ChangeKindCreated, zalozony)
	return &zalozony, nil
}

// rozbierzXliff czyta plik obu wersji standardu XLIFF jednym przejściem strumienia pakietu encoding/xml.
func rozbierzXliff(bajty []byte) plikXliff {
	czytnik := xml.NewDecoder(strings.NewReader(string(bajty)))
	wynik := plikXliff{}
	var biezaca *jednostkaXliff
	zbieraj := ""
	var tresc strings.Builder

	for {
		znacznik, err := czytnik.Token()
		if err != nil {
			break
		}
		switch element := znacznik.(type) {
		case xml.StartElement:
			switch element.Name.Local {
			case "file", "xliff":
				for _, cecha := range element.Attr {
					switch cecha.Name.Local {
					case "source-language", "srcLang":
						wynik.JezykZrodlowy = cecha.Value
					case "target-language", "trgLang":
						wynik.JezykDocelowy = cecha.Value
					}
				}
			case "trans-unit", "unit":
				jednostka := jednostkaXliff{}
				for _, cecha := range element.Attr {
					if cecha.Name.Local == "id" {
						jednostka.Klucz = cecha.Value
					}
				}
				biezaca = &jednostka
			case "source", "target":
				zbieraj = element.Name.Local
				tresc.Reset()
			}
		case xml.CharData:
			if zbieraj != "" {
				tresc.Write(element)
			}
		case xml.EndElement:
			switch element.Name.Local {
			case "source", "target":
				if biezaca != nil && zbieraj == element.Name.Local {
					if element.Name.Local == "source" {
						biezaca.Zrodlo = strings.TrimSpace(tresc.String())
					} else {
						biezaca.Cel = strings.TrimSpace(tresc.String())
					}
				}
				zbieraj = ""
			case "trans-unit", "unit":
				if biezaca != nil && strings.TrimSpace(biezaca.Zrodlo) != "" {
					wynik.Jednostki = append(wynik.Jednostki, *biezaca)
				}
				biezaca = nil
			}
		}
	}
	return wynik
}

// zlozXliff składa plik XLIFF 1.2 z par źródło–przekład. Wersja 1.2, bo tę
// przyjmują wszystkie narzędzia wykonawców, do których ten pakiet jedzie.
func zlozXliff(jezykZrodlowy, jezykDocelowy string, jednostki []jednostkaXliff) []byte {
	var b strings.Builder
	b.WriteString(xml.Header)
	b.WriteString(`<xliff version="1.2" xmlns="urn:oasis:names:tc:xliff:document:1.2">` + "\n")
	b.WriteString(`  <file original="danaco-console" datatype="plaintext" source-language="` +
		jezykZrodlowy + `" target-language="` + jezykDocelowy + `">` + "\n    <body>\n")
	for numer, jednostka := range jednostki {
		klucz := jednostka.Klucz
		if klucz == "" {
			klucz = "u" + strconv.Itoa(numer+1)
		}
		b.WriteString(`      <trans-unit id="` + klucz + `">` + "\n")
		b.WriteString("        <source>" + zabezpieczHtml(jednostka.Zrodlo) + "</source>\n")
		b.WriteString("        <target>" + zabezpieczHtml(jednostka.Cel) + "</target>\n")
		b.WriteString("      </trans-unit>\n")
	}
	b.WriteString("    </body>\n  </file>\n</xliff>\n")
	return []byte(b.String())
}
