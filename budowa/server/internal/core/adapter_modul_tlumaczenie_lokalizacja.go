// Odpowiedzialność pliku: lokalizacja oprogramowania w module Translate —
// `translate.resource.import`, `.export`, `.plural.apply`, `.key.context.set`
// oraz `translate.xliff.import`.
//
// Formaty zasobów (JSON, YAML, properties, Android XML, iOS strings
// i stringsdict, RESX, gettext PO) czyta i pisze ten rdzeń sam, bibliotekami
// wkompilowanymi. Żaden z nich nie wymaga programu zewnętrznego i żaden go tu
// nie dostaje.
//
// Formy mnogie idą regułami CLDR wpisanymi w `formyMnogieJezyka`. Reguła
// mnogości jest własnością języka, nie tłumaczenia: polski ma trzy formy,
// angielski dwie, czeski trzy, rosyjski trzy, arabski sześć. Zastosowanie form
// (`resource.plural.apply`) zakłada klucze wariantów, których język docelowy
// wymaga, i zdejmuje te, których nie zna — inaczej plik wyniku miałby formę
// `few` w języku, który jej nie ma, i program lokalizowany nigdy by jej nie użył.
package core

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekZasobuLokalizacji znakuje identyfikator zasobu lokalizacyjnego.
const przedrostekZasobuLokalizacji = "lok-"

// formyMnogieJezyka oddaje nazwy form mnogich CLDR dla języka. Język spoza
// wykazu dostaje parę `one`/`other` — najwęższy zestaw, który ma każdy język
// świata; wymyślanie mu form, których reguły nie znamy, byłoby zgadywaniem
// gramatyki cudzego języka.
func formyMnogieJezyka(jezyk string) []string {
	kod := strings.ToLower(strings.SplitN(strings.TrimSpace(jezyk), "-", 2)[0])
	switch kod {
	case "pl", "cs", "sk", "ru", "uk", "be", "hr", "sr", "lt":
		return []string{"one", "few", "many", "other"}
	case "ar":
		return []string{"zero", "one", "two", "few", "many", "other"}
	case "ja", "ko", "zh", "vi", "th":
		return []string{"other"}
	}
	return []string{"one", "other"}
}

// WczytajZasobLokalizacji obsługuje `translate.resource.import`.
func (a *adapterTlumaczenia) WczytajZasobLokalizacji(ctx context.Context,
	z shared.TranslateResourceImportRequest) (shared.TranslateResourceImportResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateResourceImportResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	sciezka := strings.TrimSpace(z.Path)
	if sciezka == "" {
		return shared.TranslateResourceImportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.resource.import bez ścieżki pliku zasobu")
	}
	format := shared.LocalizationResourceFormat("")
	if z.Format != nil {
		format = *z.Format
	}
	if strings.TrimSpace(string(format)) == "" {
		rozpoznany, jest := formatZasobuZeSciezki(sciezka)
		if !jest {
			return shared.TranslateResourceImportResponse{}, bladWskazaniaTlumaczenia(
				"nie da się rozpoznać formatu zasobu po końcówce nazwy — wskaż format wprost")
		}
		format = rozpoznany
	}

	klucze, err := wczytajKluczeZasobu(sciezka, format)
	if err != nil {
		return shared.TranslateResourceImportResponse{}, err
	}
	if len(klucze) == 0 {
		return shared.TranslateResourceImportResponse{}, bladWskazaniaTlumaczenia(
			"plik " + filepath.Base(sciezka) + " nie niesie ani jednego klucza")
	}

	zasob, err := a.repozytorium.ZapiszZasobLokalizacji(ctx, dane.ZasobLokalizacji{
		Kod:           nowyIdentyfikator(przedrostekZasobuLokalizacji),
		OknoID:        okno.ID,
		Sciezka:       sciezka,
		Format:        string(format),
		JezykZrodlowy: z.SourceLanguage,
	}, klucze)
	if err != nil {
		return shared.TranslateResourceImportResponse{}, bladTlumaczenia(err)
	}

	// Treść kluczy staje się tekstem źródłowym okna — bez tego zasób byłby
	// wczytany, a tłumaczyć nie byłoby czego.
	tresci := make([]string, 0, len(klucze))
	for _, klucz := range klucze {
		tresci = append(tresci, klucz.Tresc)
	}
	tekst := strings.Join(tresci, "\n\n")
	liczba := int64(len(klucze))
	if _, err := a.repozytorium.ZapiszOkno(ctx, dane.OknoTlumaczenia{
		Kod:             okno.Kod,
		TekstZrodlowy:   &tekst,
		JezykZrodlowy:   pierwszyJezyk(z.SourceLanguage, okno.JezykZrodlowy),
		LiczbaSegmentow: &liczba,
	}); err != nil {
		return shared.TranslateResourceImportResponse{}, bladTlumaczenia(err)
	}
	if err := a.repozytorium.UstawSegmentyOkna(ctx, okno.ID, tresci); err != nil {
		return shared.TranslateResourceImportResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslateResourceImportResponse{
		Resource: zlozZasobLokalizacji(zasob, len(klucze)),
		Keys:     zlozKluczeLokalizacji(klucze),
	}, nil
}

// pierwszyJezyk bierze język wskazany żądaniem, a przy jego braku zostawia
// język zastany okna — wczytanie zasobu nie ma prawa skasować rozpoznania.
func pierwszyJezyk(wskazany, zastany *string) *string {
	if wskazany != nil && strings.TrimSpace(*wskazany) != "" {
		return wskazany
	}
	return zastany
}

// WydajZasobLokalizacji obsługuje `translate.resource.export`. Zapisuje plik
// zasobu w języku panelu: klucze zostają te same, treść bierze się z przekładu.
func (a *adapterTlumaczenia) WydajZasobLokalizacji(ctx context.Context,
	z shared.TranslateResourceExportRequest) (shared.TranslateResourceExportResponse, error) {

	zasob, klucze, err := a.zasobZKluczami(ctx, z.ResourceId)
	if err != nil {
		return shared.TranslateResourceExportResponse{}, err
	}
	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateResourceExportResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	tresc, err := trescPanelu(panel)
	if err != nil {
		return shared.TranslateResourceExportResponse{}, err
	}
	sciezka := strings.TrimSpace(z.Path)
	if sciezka == "" {
		return shared.TranslateResourceExportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.resource.export bez ścieżki pliku wyniku")
	}
	format := shared.LocalizationResourceFormat(zasob.Format)
	if z.Format != nil && strings.TrimSpace(string(*z.Format)) != "" {
		format = *z.Format
	}

	// Przekład panelu jest ciągiem akapitów w kolejności kluczy — tak wszedł
	// do okna przy wczytaniu zasobu. Klucz bez odpowiadającego akapitu zostaje
	// z treścią źródłową, zamiast dostać pustkę.
	akapity := rozdzielAkapity(tresc)
	przelozone := make([]dane.KluczLokalizacji, 0, len(klucze))
	for numer, klucz := range klucze {
		wyjscie := klucz
		if numer < len(akapity) && strings.TrimSpace(akapity[numer]) != "" {
			wyjscie.Tresc = strings.TrimSpace(akapity[numer])
		}
		przelozone = append(przelozone, wyjscie)
	}

	if err := zapiszKluczeZasobu(sciezka, format, przelozone); err != nil {
		return shared.TranslateResourceExportResponse{}, err
	}
	return shared.TranslateResourceExportResponse{
		Path:          sciezka,
		ExportedCount: len(przelozone),
	}, nil
}

// ZastosujFormyMnogie obsługuje `translate.resource.plural.apply`.
func (a *adapterTlumaczenia) ZastosujFormyMnogie(ctx context.Context,
	z shared.TranslateResourcePluralApplyRequest) (shared.TranslateResourcePluralApplyResponse, error) {

	zasob, klucze, err := a.zasobZKluczami(ctx, z.ResourceId)
	if err != nil {
		return shared.TranslateResourcePluralApplyResponse{}, err
	}
	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateResourcePluralApplyResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	wskazane := map[string]struct{}{}
	for _, klucz := range z.Keys {
		wskazane[klucz] = struct{}{}
	}
	formy := formyMnogieJezyka(panel.Jezyk)

	zmienione := 0
	for _, klucz := range klucze {
		if len(wskazane) > 0 {
			if _, jest := wskazane[klucz.Klucz]; !jest {
				continue
			}
		}
		nowe := map[string]string{}
		zastane := map[string]string{}
		if klucz.FormyMnogie != nil {
			_ = json.Unmarshal([]byte(*klucz.FormyMnogie), &zastane)
		}
		for _, forma := range formy {
			if tresc, jest := zastane[forma]; jest {
				nowe[forma] = tresc
				continue
			}
			// Forma, której język wymaga, a której nie było: zakładana jest
			// z treści klucza. To jest miejsce do wypełnienia przez tłumacza,
			// a nie gotowa odmiana — ale bez niego plik wyniku nie ma gdzie jej
			// nawet zapisać.
			nowe[forma] = klucz.Tresc
		}
		if len(nowe) == len(zastane) && formyRowne(nowe, zastane) {
			continue
		}
		bajty, err := json.Marshal(nowe)
		if err != nil {
			return shared.TranslateResourcePluralApplyResponse{}, bladTlumaczenia(err)
		}
		zmieniony := klucz
		tresc := string(bajty)
		zmieniony.FormyMnogie = &tresc
		if err := a.repozytorium.ZapiszKluczLokalizacji(ctx, zasob.ID, zmieniony); err != nil {
			return shared.TranslateResourcePluralApplyResponse{}, bladTlumaczenia(err)
		}
		zmienione++
	}

	poZmianie, err := a.repozytorium.KluczeLokalizacji(ctx, zasob.ID)
	if err != nil {
		return shared.TranslateResourcePluralApplyResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateResourcePluralApplyResponse{
		Keys:         zlozKluczeLokalizacji(poZmianie),
		ChangedCount: zmienione,
	}, nil
}

// formyRowne mówi, czy dwa zestawy form są tożsame — inaczej zapis liczyłby
// zmianę tam, gdzie nic się nie zmieniło.
func formyRowne(pierwszy, drugi map[string]string) bool {
	for forma, tresc := range pierwszy {
		if drugi[forma] != tresc {
			return false
		}
	}
	return true
}

// UstawKontekstKlucza obsługuje `translate.resource.key.context.set`.
func (a *adapterTlumaczenia) UstawKontekstKlucza(ctx context.Context,
	z shared.TranslateResourceKeyContextSetRequest) (shared.TranslateResourceKeyContextSetResponse, error) {

	zasob, klucze, err := a.zasobZKluczami(ctx, z.ResourceId)
	if err != nil {
		return shared.TranslateResourceKeyContextSetResponse{}, err
	}
	for _, klucz := range klucze {
		if klucz.Klucz != z.Key {
			continue
		}
		zmieniony := klucz
		zmieniony.Kontekst = z.Context
		zmieniony.ZrzutZasobID = z.ScreenshotAssetId
		if err := a.repozytorium.ZapiszKluczLokalizacji(ctx, zasob.ID, zmieniony); err != nil {
			return shared.TranslateResourceKeyContextSetResponse{}, bladTlumaczenia(err)
		}
		return shared.TranslateResourceKeyContextSetResponse{
			Key: zlozKluczLokalizacji(zmieniony)}, nil
	}
	return shared.TranslateResourceKeyContextSetResponse{}, bladWskazaniaTlumaczenia(
		"zasób " + z.ResourceId + " nie ma klucza " + z.Key)
}

// zasobZKluczami odczytuje zasób wraz z kluczami — wspólne wejście trzech
// komend pracujących na kluczach.
func (a *adapterTlumaczenia) zasobZKluczami(ctx context.Context,
	kod string) (dane.ZasobLokalizacji, []dane.KluczLokalizacji, error) {

	zasob, err := a.repozytorium.ZasobLokalizacji(ctx, kod)
	if err != nil {
		return dane.ZasobLokalizacji{}, nil, bladWskazaniaTlumaczenia(
			"nie ma zasobu lokalizacyjnego o identyfikatorze " + kod)
	}
	klucze, err := a.repozytorium.KluczeLokalizacji(ctx, zasob.ID)
	if err != nil {
		return dane.ZasobLokalizacji{}, nil, bladTlumaczenia(err)
	}
	return zasob, klucze, nil
}

// zlozZasobLokalizacji przekłada wiersz zasobu na byt kontraktu.
func zlozZasobLokalizacji(zasob dane.ZasobLokalizacji, kluczy int) shared.LocalizationResource {
	return shared.LocalizationResource{
		Id:             zasob.Kod,
		WindowId:       zasob.OknoKod,
		Path:           zasob.Sciezka,
		Format:         shared.LocalizationResourceFormat(zasob.Format),
		SourceLanguage: zasob.JezykZrodlowy,
		KeyCount:       kluczy,
		UpdatedAt:      zasob.Zaktualizowano,
	}
}

// zlozKluczeLokalizacji przekłada wiersze kluczy na byty kontraktu.
func zlozKluczeLokalizacji(klucze []dane.KluczLokalizacji) []shared.LocalizationKey {
	wykaz := make([]shared.LocalizationKey, 0, len(klucze))
	for _, klucz := range klucze {
		wykaz = append(wykaz, zlozKluczLokalizacji(klucz))
	}
	return wykaz
}

// zlozKluczLokalizacji przekłada jeden wiersz klucza na byt kontraktu.
func zlozKluczLokalizacji(klucz dane.KluczLokalizacji) shared.LocalizationKey {
	byt := shared.LocalizationKey{
		Key:               klucz.Klucz,
		Text:              klucz.Tresc,
		Placeholders:      klucz.Znaczniki,
		Context:           klucz.Kontekst,
		ScreenshotAssetId: klucz.ZrzutZasobID,
	}
	if byt.Placeholders == nil {
		byt.Placeholders = []string{}
	}
	if klucz.FormyMnogie != nil {
		byt.PluralForms = json.RawMessage(*klucz.FormyMnogie)
	}
	return byt
}

// formatZasobuZeSciezki rozpoznaje format zasobu po końcówce nazwy.
func formatZasobuZeSciezki(sciezka string) (shared.LocalizationResourceFormat, bool) {
	nazwa := strings.ToLower(filepath.Base(sciezka))
	switch {
	case strings.HasSuffix(nazwa, ".json"):
		return shared.LocalizationResourceFormatJson, true
	case strings.HasSuffix(nazwa, ".yaml"), strings.HasSuffix(nazwa, ".yml"):
		return shared.LocalizationResourceFormatYaml, true
	case strings.HasSuffix(nazwa, ".properties"):
		return shared.LocalizationResourceFormatProperties, true
	case strings.HasSuffix(nazwa, ".stringsdict"):
		return shared.LocalizationResourceFormatIosStringsdict, true
	case strings.HasSuffix(nazwa, ".strings"):
		return shared.LocalizationResourceFormatIosStrings, true
	case strings.HasSuffix(nazwa, ".resx"):
		return shared.LocalizationResourceFormatResx, true
	case strings.HasSuffix(nazwa, ".po"), strings.HasSuffix(nazwa, ".pot"):
		return shared.LocalizationResourceFormatGettextPo, true
	case strings.HasSuffix(nazwa, ".xml"):
		return shared.LocalizationResourceFormatAndroidXml, true
	}
	return "", false
}

// wczytajKluczeZasobu czyta klucze z pliku wskazanego formatu.
func wczytajKluczeZasobu(sciezka string,
	format shared.LocalizationResourceFormat) ([]dane.KluczLokalizacji, error) {

	bajty, err := os.ReadFile(filepath.Clean(sciezka))
	if err != nil {
		return nil, bladPlikuTlumaczenia(sciezka, err)
	}
	pary := map[string]string{}

	switch format {
	case shared.LocalizationResourceFormatJson:
		surowe := map[string]any{}
		if err := json.Unmarshal(bajty, &surowe); err != nil {
			return nil, bladWskazaniaTlumaczenia("plik nie jest czytelnym JSON: " + err.Error())
		}
		splaszczKlucze("", surowe, pary)
	case shared.LocalizationResourceFormatYaml:
		surowe := map[string]any{}
		if err := yaml.Unmarshal(bajty, &surowe); err != nil {
			return nil, bladWskazaniaTlumaczenia("plik nie jest czytelnym YAML: " + err.Error())
		}
		splaszczKlucze("", surowe, pary)
	case shared.LocalizationResourceFormatProperties:
		for _, linia := range strings.Split(string(bajty), "\n") {
			linia = strings.TrimSpace(linia)
			if linia == "" || strings.HasPrefix(linia, "#") || strings.HasPrefix(linia, "!") {
				continue
			}
			klucz, tresc, jest := strings.Cut(linia, "=")
			if !jest {
				continue
			}
			pary[strings.TrimSpace(klucz)] = strings.TrimSpace(tresc)
		}
	case shared.LocalizationResourceFormatIosStrings:
		for _, linia := range strings.Split(string(bajty), "\n") {
			linia = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(linia), ";"))
			klucz, tresc, jest := strings.Cut(linia, "=")
			if !jest {
				continue
			}
			pary[zdejmijCudzyslowy(klucz)] = zdejmijCudzyslowy(tresc)
		}
	case shared.LocalizationResourceFormatAndroidXml, shared.LocalizationResourceFormatResx,
		shared.LocalizationResourceFormatIosStringsdict:
		wczytane, err := paryZXmlZasobu(bajty, format)
		if err != nil {
			return nil, err
		}
		pary = wczytane
	case shared.LocalizationResourceFormatGettextPo:
		pary = paryZGettext(string(bajty))
	default:
		return nil, bladWskazaniaTlumaczenia("nieznany format zasobu: " + string(format))
	}

	nazwy := make([]string, 0, len(pary))
	for nazwa := range pary {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)

	klucze := make([]dane.KluczLokalizacji, 0, len(nazwy))
	for numer, nazwa := range nazwy {
		klucze = append(klucze, dane.KluczLokalizacji{
			Klucz:     nazwa,
			Tresc:     pary[nazwa],
			Znaczniki: znacznikiKluczaLokalizacji(pary[nazwa]),
			Kolejnosc: int64(numer),
		})
	}
	return klucze, nil
}

// splaszczKlucze sprowadza drzewo JSON albo YAML do par klucz-treść z kropką
// jako separatorem poziomów. Zagnieżdżenie jest wygodą pliku, nie tożsamością
// klucza — program lokalizowany i tak adresuje je jedną nazwą.
func splaszczKlucze(przedrostek string, wezel map[string]any, pary map[string]string) {
	for nazwa, wartosc := range wezel {
		pelna := nazwa
		if przedrostek != "" {
			pelna = przedrostek + "." + nazwa
		}
		switch treść := wartosc.(type) {
		case string:
			pary[pelna] = treść
		case map[string]any:
			splaszczKlucze(pelna, treść, pary)
		}
	}
}

// zdejmijCudzyslowy zdejmuje cudzysłowy i odstępy wokół napisu formatu iOS.
func zdejmijCudzyslowy(tekst string) string {
	return strings.Trim(strings.TrimSpace(tekst), "\"")
}

// paryZXmlZasobu czyta pary z formatów opartych na XML.
func paryZXmlZasobu(bajty []byte,
	format shared.LocalizationResourceFormat) (map[string]string, error) {

	czytnik := xml.NewDecoder(strings.NewReader(string(bajty)))
	pary := map[string]string{}
	nazwaBiezaca := ""
	zbieraj := false
	var tresc strings.Builder
	for {
		znacznik, err := czytnik.Token()
		if err != nil {
			break
		}
		switch element := znacznik.(type) {
		case xml.StartElement:
			nazwaWezla := element.Name.Local
			interesujacy := (format == shared.LocalizationResourceFormatAndroidXml && nazwaWezla == "string") ||
				(format == shared.LocalizationResourceFormatResx && nazwaWezla == "data") ||
				(format == shared.LocalizationResourceFormatIosStringsdict && nazwaWezla == "key")
			if !interesujacy {
				continue
			}
			zbieraj = true
			tresc.Reset()
			for _, cecha := range element.Attr {
				if cecha.Name.Local == "name" {
					nazwaBiezaca = cecha.Value
				}
			}
		case xml.CharData:
			if zbieraj {
				tresc.Write(element)
			}
		case xml.EndElement:
			if !zbieraj {
				continue
			}
			zbieraj = false
			wartosc := strings.TrimSpace(tresc.String())
			if nazwaBiezaca == "" {
				// Stringsdict adresuje wartość kluczem podanym treścią węzła
				// `key` — wtedy nazwą jest sama treść.
				nazwaBiezaca = wartosc
			}
			if nazwaBiezaca != "" && wartosc != "" {
				pary[nazwaBiezaca] = wartosc
			}
			nazwaBiezaca = ""
		}
	}
	if len(pary) == 0 {
		return nil, bladWskazaniaTlumaczenia("plik XML nie niesie ani jednego klucza tego formatu")
	}
	return pary, nil
}

// paryZGettext czyta pary `msgid`/`msgstr` z pliku PO.
func paryZGettext(tresc string) map[string]string {
	pary := map[string]string{}
	klucz := ""
	for _, linia := range strings.Split(tresc, "\n") {
		linia = strings.TrimSpace(linia)
		switch {
		case strings.HasPrefix(linia, "msgid "):
			klucz = zdejmijCudzyslowy(strings.TrimPrefix(linia, "msgid "))
		case strings.HasPrefix(linia, "msgstr "):
			wartosc := zdejmijCudzyslowy(strings.TrimPrefix(linia, "msgstr "))
			if klucz == "" {
				continue
			}
			if wartosc == "" {
				// Wpis nieprzetłumaczony niesie treść w kluczu — to ona idzie
				// do tłumaczenia.
				wartosc = klucz
			}
			pary[klucz] = wartosc
			klucz = ""
		}
	}
	return pary
}

// znacznikiKluczaLokalizacji wypisuje znaczniki obecne w treści klucza:
// `{nazwa}`, `%s`, `%1$s` i `%d`. Osobna od `znacznikiPodstawienia`
// (`*_jakosc_zrodlo.go`), która zna wyłącznie klamry: pliki zasobów niosą
// znaczniki w postaci języka programowania, a kontrola jakości panelu — nie.
func znacznikiKluczaLokalizacji(tresc string) []string {
	znaczniki := []string{}
	for i := 0; i < len(tresc); i++ {
		switch tresc[i] {
		case '{':
			koniec := strings.IndexByte(tresc[i:], '}')
			if koniec > 0 {
				znaczniki = append(znaczniki, tresc[i:i+koniec+1])
				i += koniec
			}
		case '%':
			koniec := i + 1
			for koniec < len(tresc) && strings.ContainsRune("0123456789$.", rune(tresc[koniec])) {
				koniec++
			}
			if koniec < len(tresc) {
				znaczniki = append(znaczniki, tresc[i:koniec+1])
				i = koniec
			}
		}
	}
	return znaczniki
}

// zapiszKluczeZasobu wypisuje klucze do pliku wskazanego formatu.
func zapiszKluczeZasobu(sciezka string, format shared.LocalizationResourceFormat,
	klucze []dane.KluczLokalizacji) error {

	if err := os.MkdirAll(filepath.Dir(sciezka), 0o700); err != nil {
		return bladPlikuTlumaczenia(sciezka, err)
	}
	var tresc []byte
	switch format {
	case shared.LocalizationResourceFormatJson:
		pary := map[string]string{}
		for _, klucz := range klucze {
			pary[klucz.Klucz] = klucz.Tresc
		}
		bajty, err := json.MarshalIndent(pary, "", "  ")
		if err != nil {
			return bladTlumaczenia(err)
		}
		tresc = append(bajty, '\n')
	case shared.LocalizationResourceFormatYaml:
		pary := map[string]string{}
		for _, klucz := range klucze {
			pary[klucz.Klucz] = klucz.Tresc
		}
		bajty, err := yaml.Marshal(pary)
		if err != nil {
			return bladTlumaczenia(err)
		}
		tresc = bajty
	case shared.LocalizationResourceFormatProperties:
		var b strings.Builder
		for _, klucz := range klucze {
			b.WriteString(klucz.Klucz + "=" + klucz.Tresc + "\n")
		}
		tresc = []byte(b.String())
	case shared.LocalizationResourceFormatIosStrings:
		var b strings.Builder
		for _, klucz := range klucze {
			b.WriteString("\"" + klucz.Klucz + "\" = \"" + klucz.Tresc + "\";\n")
		}
		tresc = []byte(b.String())
	case shared.LocalizationResourceFormatAndroidXml:
		var b strings.Builder
		b.WriteString(xml.Header + "<resources>\n")
		for _, klucz := range klucze {
			b.WriteString("  <string name=\"" + klucz.Klucz + "\">" +
				zabezpieczHtml(klucz.Tresc) + "</string>\n")
		}
		b.WriteString("</resources>\n")
		tresc = []byte(b.String())
	case shared.LocalizationResourceFormatResx:
		var b strings.Builder
		b.WriteString(xml.Header + "<root>\n")
		for _, klucz := range klucze {
			b.WriteString("  <data name=\"" + klucz.Klucz + "\" xml:space=\"preserve\">" +
				zabezpieczHtml(klucz.Tresc) + "</data>\n")
		}
		b.WriteString("</root>\n")
		tresc = []byte(b.String())
	case shared.LocalizationResourceFormatGettextPo:
		var b strings.Builder
		for _, klucz := range klucze {
			b.WriteString("msgid \"" + klucz.Klucz + "\"\nmsgstr \"" + klucz.Tresc + "\"\n\n")
		}
		tresc = []byte(b.String())
	case shared.LocalizationResourceFormatIosStringsdict:
		// Stringsdict niesie formy mnogie, więc wypisujemy je, gdy klucz je ma.
		var b strings.Builder
		b.WriteString(xml.Header + "<plist version=\"1.0\"><dict>\n")
		for _, klucz := range klucze {
			b.WriteString("  <key>" + klucz.Klucz + "</key>\n  <dict>\n")
			formy := map[string]string{"other": klucz.Tresc}
			if klucz.FormyMnogie != nil {
				odczytane := map[string]string{}
				if err := json.Unmarshal([]byte(*klucz.FormyMnogie), &odczytane); err == nil &&
					len(odczytane) > 0 {
					formy = odczytane
				}
			}
			nazwy := make([]string, 0, len(formy))
			for nazwa := range formy {
				nazwy = append(nazwy, nazwa)
			}
			sort.Strings(nazwy)
			for _, nazwa := range nazwy {
				b.WriteString("    <key>" + nazwa + "</key><string>" +
					zabezpieczHtml(formy[nazwa]) + "</string>\n")
			}
			b.WriteString("  </dict>\n")
		}
		b.WriteString("</dict></plist>\n")
		tresc = []byte(b.String())
	default:
		return bladWskazaniaTlumaczenia("nieznany format zasobu: " + string(format))
	}
	return zapiszPlikWyniku(sciezka, tresc)
}
