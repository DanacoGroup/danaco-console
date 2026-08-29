// Odpowiedzialność pliku: rachunek OOXML modułu Studio — odczyt dokumentów
// `.docx` i `.dotx` do wspólnej postaci dokumentu oraz złożenie tej postaci
// z powrotem do archiwum biurowego, bez biblioteki obcej ponad standardową.
package core

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// granicaArchiwumStudia chroni rdzeń przed archiwum, które rozpakowuje się do
// rozmiaru pamięci maszyny. Plik Operatora bywa duży; archiwum rozprężające się
// do gigabajtów bywa złośliwe.
const granicaArchiwumStudia = 256 << 20

// granicaSkladnikaArchiwum jest granicą jednego składnika archiwum: plik
// przekraczający ją odpada z rozbioru, mimo że całe archiwum mieści się
// w granicy ogólnej.
const granicaSkladnikaArchiwum = 64 << 20

// Nazwy składników archiwum OOXML, których dotyka ten rachunek: dokument
// główny, arkusz stylów, numeracja list, powiązania dokumentu i katalog
// nośników.
const (
	ooxmlSkladnikDokumentu    = "word/document.xml"
	ooxmlSkladnikStylow       = "word/styles.xml"
	ooxmlSkladnikNumeracji    = "word/numbering.xml"
	ooxmlSkladnikPowiazan     = "word/_rels/document.xml.rels"
	ooxmlKatalogNosnikow      = "word/media/"
	ooxmlPrzestrzenGlowna     = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
	ooxmlPrzestrzenPowiazania = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
)

// ── Drzewo XML ──────────────────────────────────────────────────────────────

// ooxmlWezel jest węzłem drzewa XML: nazwa wraz z przestrzenią, atrybuty,
// dzieci i treść tekstowa zebrana z tego węzła (bez dzieci).
type ooxmlWezel struct {
	Nazwa     xml.Name
	Atrybuty  []xml.Attr
	Dzieci    []*ooxmlWezel
	Tekst     string
	rodzicPtr *ooxmlWezel
}

// ooxmlRozbierz buduje drzewo z bajtów XML. Rozbiór idzie `xml.Decoder` ze
// wskazanym rozstrzygaczem znaków, bo plik biurowy zapisany kodowaniem
// starszym niż utf-8 inaczej odmówiłby odczytu na nagłówku.
func ooxmlRozbierz(bajty []byte) (*ooxmlWezel, error) {
	rozbierak := xml.NewDecoder(bytes.NewReader(bajty))
	rozbierak.CharsetReader = wejscieRozstrzygaczZnakowXml
	rozbierak.Strict = false

	korzen := &ooxmlWezel{}
	biezacy := korzen
	for {
		token, err := rozbierak.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("nieczytelny XML dokumentu: %w", err)
		}
		switch wartosc := token.(type) {
		case xml.StartElement:
			dziecko := &ooxmlWezel{
				Nazwa:     wartosc.Name,
				Atrybuty:  append([]xml.Attr(nil), wartosc.Attr...),
				rodzicPtr: biezacy,
			}
			biezacy.Dzieci = append(biezacy.Dzieci, dziecko)
			biezacy = dziecko
		case xml.EndElement:
			if biezacy.rodzicPtr != nil {
				biezacy = biezacy.rodzicPtr
			}
		case xml.CharData:
			biezacy.Tekst += string(wartosc)
		}
	}
	if len(korzen.Dzieci) == 0 {
		return nil, fmt.Errorf("dokument XML jest pusty")
	}
	return korzen.Dzieci[0], nil
}

// dziecko oddaje pierwsze dziecko węzła o wskazanej nazwie lokalnej albo nic,
// gdy węzeł go nie niesie lub gdy wskazany węzeł w ogóle nie istnieje.
func (w *ooxmlWezel) dziecko(nazwa string) *ooxmlWezel {
	if w == nil {
		return nil
	}
	for _, dziecko := range w.Dzieci {
		if dziecko.Nazwa.Local == nazwa {
			return dziecko
		}
	}
	return nil
}

// dzieci oddaje wszystkie dzieci węzła o wskazanej nazwie lokalnej; gdy
// takiego dziecka nie ma, oddaje wykaz pusty, a nie wartość nieokreśloną.
func (w *ooxmlWezel) dzieci(nazwa string) []*ooxmlWezel {
	if w == nil {
		return nil
	}
	wynik := make([]*ooxmlWezel, 0, len(w.Dzieci))
	for _, dziecko := range w.Dzieci {
		if dziecko.Nazwa.Local == nazwa {
			wynik = append(wynik, dziecko)
		}
	}
	return wynik
}

// sciezka schodzi po kolejnych nazwach lokalnych. Brak któregokolwiek ogniwa
// oddaje nic — wołający sprawdza jeden raz, nie po każdym kroku.
func (w *ooxmlWezel) sciezka(nazwy ...string) *ooxmlWezel {
	biezacy := w
	for _, nazwa := range nazwy {
		biezacy = biezacy.dziecko(nazwa)
		if biezacy == nil {
			return nil
		}
	}
	return biezacy
}

// atrybut oddaje wartość atrybutu węzła o wskazanej nazwie lokalnej; gdy
// atrybutu nie ma, oddaje napis pusty, nigdy wartość nieokreśloną.
func (w *ooxmlWezel) atrybut(nazwa string) string {
	if w == nil {
		return ""
	}
	for _, atrybut := range w.Atrybuty {
		if atrybut.Name.Local == nazwa {
			return atrybut.Value
		}
	}
	return ""
}

// atrybutCalkowity oddaje atrybut jako liczbę całkowitą wraz z informacją,
// czy atrybut w ogóle stał w pliku. Brak i zero są dwiema różnymi rzeczami:
// margines nieustawiony to nie margines zerowy.
func (w *ooxmlWezel) atrybutCalkowity(nazwa string) (int, bool) {
	tekst := strings.TrimSpace(w.atrybut(nazwa))
	if tekst == "" {
		return 0, false
	}
	// Wartości bywają zapisane z częścią dziesiętną, choć nie powinny.
	liczba, err := strconv.ParseFloat(tekst, 64)
	if err != nil {
		return 0, false
	}
	return int(liczba), true
}

// wlaczonyPrzelacznik rozstrzyga przełącznik OOXML: sam obecny węzeł znaczy
// „włączone", a `w:val` równe `0`, `false` albo `off` znaczy „wyłączone" —
// tak wyłącza się formatowanie odziedziczone po stylu.
func (w *ooxmlWezel) wlaczonyPrzelacznik() bool {
	if w == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(w.atrybut("val"))) {
	case "0", "false", "off", "none":
		return false
	}
	return true
}

// ── Archiwum ────────────────────────────────────────────────────────────────

// wejscieOtworzArchiwum rozkłada archiwum na składniki. Rozpakowanie idzie
// `archive/zip` z biblioteki wzorcowej, bo plik biurowy JEST archiwum ZIP,
// a osobne narzędzie archiwizujące serwera służy innym formatom.
func wejscieOtworzArchiwum(bajty []byte) (map[string][]byte, error) {
	if len(bajty) > granicaArchiwumStudia {
		return nil, bladWskazaniaStudio("plik jest większy niż " +
			strconv.Itoa(granicaArchiwumStudia>>20) + " MB i serwer go nie rozpakowuje")
	}
	archiwum, err := zip.NewReader(bytes.NewReader(bajty), int64(len(bajty)))
	if err != nil {
		return nil, bladWskazaniaStudio(
			"plik nie jest archiwum biurowym (docx, dotx, odt, ott): " + err.Error())
	}
	skladniki := make(map[string][]byte, len(archiwum.File))
	for _, wpis := range archiwum.File {
		if wpis.FileInfo().IsDir() {
			continue
		}
		// Nazwa składnika mogłaby wyprowadzać poza archiwum („../") — taki
		// wpis jest odrzucany jako złośliwy.
		nazwa := path.Clean(strings.ReplaceAll(wpis.Name, "\\", "/"))
		if strings.HasPrefix(nazwa, "../") || nazwa == ".." {
			continue
		}
		if wpis.UncompressedSize64 > granicaSkladnikaArchiwum {
			continue
		}
		strumien, err := wpis.Open()
		if err != nil {
			return nil, bladWskazaniaStudio(
				"nie można odczytać składnika " + nazwa + " archiwum: " + err.Error())
		}
		tresc, err := io.ReadAll(io.LimitReader(strumien, granicaSkladnikaArchiwum))
		_ = strumien.Close()
		if err != nil {
			return nil, bladWskazaniaStudio(
				"nie można odczytać składnika " + nazwa + " archiwum: " + err.Error())
		}
		skladniki[nazwa] = tresc
	}
	return skladniki, nil
}

// ── Miary ───────────────────────────────────────────────────────────────────

// ooxmlTwipyNaMilimetry przelicza twipy, czyli 1/1440 cala, na milimetry —
// jednostkę, w której kontrakt platformy wyraża odległości strony.
func ooxmlTwipyNaMilimetry(twipy int) float64 {
	return float64(twipy) / 1440 * 25.4
}

// ooxmlMilimetryNaTwipy przelicza milimetry kontraktu platformy z powrotem
// na twipy, czyli 1/1440 cala — jednostkę odległości używaną przez OOXML.
func ooxmlMilimetryNaTwipy(milimetry float64) int {
	return int(milimetry / 25.4 * 1440)
}

// ooxmlPolpunktyNaPunkty przelicza stopień pisma OOXML z półpunktów na
// punkty — jednostkę, w której kontrakt platformy wyraża wielkość czcionki.
func ooxmlPolpunktyNaPunkty(polpunkty int) float64 {
	return float64(polpunkty) / 2
}

// ooxmlTwipyNaPunkty przelicza odstęp akapitowy OOXML z twipów na punkty —
// jednostkę, w której kontrakt platformy wyraża odstępy akapitu.
func ooxmlTwipyNaPunkty(twipy int) float64 {
	return float64(twipy) / 20
}

// ── Odczyt ──────────────────────────────────────────────────────────────────

// wejscieCzytajOoxml rozbiera `.docx` albo `.dotx` na postać dokumentu, treść
// płaską i bilans wniesienia, który nazywa wprost to, czego rachunek nie
// odczytuje, zamiast przemilczeć okaleczenie dokumentu.
func wejscieCzytajOoxml(kodDokumentu string, bajty []byte,
	format shared.StudioImportFormat) (shared.StudioDocumentForm, string, shared.StudioImportBalance, error) {

	bilans := shared.StudioImportBalance{Format: format}
	skladniki, err := wejscieOtworzArchiwum(bajty)
	if err != nil {
		return shared.StudioDocumentForm{}, "", bilans, err
	}
	surowyDokument, jest := skladniki[ooxmlSkladnikDokumentu]
	if !jest {
		return shared.StudioDocumentForm{}, "", bilans, bladWskazaniaStudio(
			"archiwum nie niesie składnika " + ooxmlSkladnikDokumentu +
				" — to nie jest dokument Word (docx, dotx)")
	}
	dokument, err := ooxmlRozbierz(surowyDokument)
	if err != nil {
		return shared.StudioDocumentForm{}, "", bilans, bladWskazaniaStudio(err.Error())
	}

	postac := shared.StudioDocumentForm{DocumentId: kodDokumentu}

	// Arkusz stylów wchodzi PIERWSZY: akapity odwołują się do stylów nazwą,
	// którą styl musi już mieć.
	postac.Styles = ooxmlCzytajStyle(skladniki[ooxmlSkladnikStylow])
	if len(postac.Styles) == 0 {
		postac.Styles = wejscieDomyslnyArkuszStylow()
	}
	bilans.StylesRecovered = wejscieWskaznikCalkowity(len(postac.Styles))

	postac.Lists = ooxmlCzytajNumeracje(skladniki[ooxmlSkladnikNumeracji])
	powiazania := ooxmlCzytajPowiazania(skladniki[ooxmlSkladnikPowiazan])

	cialo := dokument.dziecko("body")
	if cialo == nil {
		return shared.StudioDocumentForm{}, "", bilans, bladWskazaniaStudio(
			"dokument Word bez ciała (w:body) — plik jest uszkodzony")
	}

	stan := &ooxmlStanOdczytu{
		postac:     &postac,
		bilans:     &bilans,
		skladniki:  skladniki,
		powiazania: powiazania,
	}
	stan.zalozSekcje(cialo.dziecko("sectPr"))
	stan.czytajCialo(cialo)

	// Sekcja ostatnia bierze nastawy z `w:sectPr` ciała, gdzie stoją nastawy
	// ostatniej sekcji dokumentu.
	if len(postac.Sections) > 0 && postac.PageSetup == nil {
		postac.PageSetup = postac.Sections[len(postac.Sections)-1].PageSetup
	}
	if postac.PageSetup == nil {
		nastawy := wejscieDomyslneNastawyStrony("", nil)
		postac.PageSetup = &nastawy
	}
	bilans.SectionsRecovered = wejscieWskaznikCalkowity(len(postac.Sections))
	bilans.ParagraphsRecovered = wejscieWskaznikCalkowity(len(postac.Blocks))
	bilans.TablesRecognized = wejscieWskaznikCalkowity(len(postac.Tables))

	tresc := wejscieTrescZPostaci(&postac)
	wejscieUzycieStylow(&postac)
	bilans.Note = wejscieWskaznikTekstu(ooxmlZdanieBilansu(&bilans))
	return postac, tresc, bilans, nil
}

// ooxmlStanOdczytu zbiera to, co rozbiór ciała dokumentu musi mieć pod ręką.
// Bez niego każda funkcja rozbioru brałaby te same pięć argumentów.
type ooxmlStanOdczytu struct {
	postac        *shared.StudioDocumentForm
	bilans        *shared.StudioImportBalance
	skladniki     map[string][]byte
	powiazania    map[string]string
	biezacaSekcja string
}

// zalozSekcje zakłada kolejną sekcję dokumentu wraz z jej nastawami strony,
// nagłówkami, stopkami i numeracją stron, i ustawia ją jako sekcję bieżącą.
func (s *ooxmlStanOdczytu) zalozSekcje(wezelSekcji *ooxmlWezel) {
	nastawy := ooxmlCzytajNastawyStrony(wezelSekcji)
	sekcja := shared.StudioSection{
		Id:        nowyIdentyfikator(przedrostekSekcjiStudia),
		Index:     len(s.postac.Sections),
		Start:     wejscieWskaznikPoczatkuSekcji(ooxmlPoczatekSekcji(wezelSekcji)),
		PageSetup: &nastawy,
	}
	if naglowki := ooxmlCzytajNaglowkiStopki(wezelSekcji, s.skladniki, s.powiazania); len(naglowki) > 0 {
		sekcja.HeadersFooters = naglowki
	}
	if numeracja := ooxmlCzytajNumeracjeStron(wezelSekcji); numeracja != nil {
		sekcja.Numbering = numeracja
	}
	s.postac.Sections = append(s.postac.Sections, sekcja)
	s.biezacaSekcja = sekcja.Id
}

// czytajCialo przechodzi ciało dokumentu blok po bloku w kolejności czytania,
// rozpoznając akapity, tabele i kontrolki treści opakowujące zwykłe akapity.
func (s *ooxmlStanOdczytu) czytajCialo(cialo *ooxmlWezel) {
	for _, wezel := range cialo.Dzieci {
		switch wezel.Nazwa.Local {
		case "p":
			s.czytajAkapit(wezel)
		case "tbl":
			s.czytajTabele(wezel)
		case "sdt":
			// Kontrolka treści (`w:sdt`) opakowuje zwykłą treść — schodzi się
			// do jej wnętrza, zamiast ją pomijać.
			if wnetrze := wezel.dziecko("sdtContent"); wnetrze != nil {
				s.czytajCialo(wnetrze)
			}
		case "sectPr":
			// Nastawy ostatniej sekcji — już wzięte przy zakładaniu sekcji.
		}
	}
}

// czytajAkapit składa blok akapitu wraz z postacią akapitu i fragmentami.
// Podział sekcji rozpoznaje się po `w:sectPr` wewnątrz akapitu: taki akapit
// jest ostatnim akapitem sekcji, która się właśnie skończyła.
func (s *ooxmlStanOdczytu) czytajAkapit(wezel *ooxmlWezel) {
	postacAkapitu := wezel.dziecko("pPr")
	blok := shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuAkapit,
		SectionId: wejscieWskaznikTekstu(s.biezacaSekcja),
		Paragraph: ooxmlCzytajPostacAkapitu(postacAkapitu),
	}
	blok.Runs = s.czytajFragmenty(wezel)
	s.postac.Blocks = append(s.postac.Blocks, blok)

	if postacAkapitu != nil {
		if wezelSekcji := postacAkapitu.dziecko("sectPr"); wezelSekcji != nil {
			// Nastawy zamykają sekcję bieżącą, a nowa sekcja zaczyna się od
			// akapitu następnego.
			if ostatnia := len(s.postac.Sections) - 1; ostatnia >= 0 {
				nastawy := ooxmlCzytajNastawyStrony(wezelSekcji)
				s.postac.Sections[ostatnia].PageSetup = &nastawy
			}
			s.zalozSekcje(nil)
		}
	}
}

// czytajFragmenty składa fragmenty tekstu akapitu wraz z postacią znaku.
// Obraz osadzony w akapicie (`w:drawing`) wchodzi do wykazu obiektów, a nie
// do treści — treść dokumentu to litery, obraz jest obiektem zakotwiczonym.
func (s *ooxmlStanOdczytu) czytajFragmenty(wezel *ooxmlWezel) []shared.StudioDocumentRun {
	fragmenty := make([]shared.StudioDocumentRun, 0, len(wezel.Dzieci))
	for _, dziecko := range wezel.Dzieci {
		switch dziecko.Nazwa.Local {
		case "r":
			postacZnaku := ooxmlCzytajPostacZnaku(dziecko.dziecko("rPr"))
			tekst := ooxmlTekstFragmentu(dziecko)
			if obrazy := dziecko.dziecko("drawing"); obrazy != nil {
				s.czytajObraz(obrazy)
			}
			if tekst == "" {
				continue
			}
			fragmenty = append(fragmenty, shared.StudioDocumentRun{
				Text: tekst, Format: postacZnaku,
			})
		case "hyperlink":
			// Odsyłacz niesie fragmenty w środku; adres wchodzi osobno do
			// aparatu dokumentu.
			wewnetrzne := s.czytajFragmenty(dziecko)
			fragmenty = append(fragmenty, wewnetrzne...)
			s.czytajOdsylacz(dziecko, wewnetrzne)
		case "ins", "del", "smartTag":
			// Zmiana śledzona i znacznik inteligentny opakowują zwykłe
			// fragmenty; treść usunięta odpada sama.
			fragmenty = append(fragmenty, s.czytajFragmenty(dziecko)...)
		}
	}
	return fragmenty
}

// czytajObraz dokłada obraz osadzony do wykazu obiektów i liczy go w bilansie.
// Bajty obrazu zostają w archiwum — tu powstaje obiekt wraz z nazwą składnika,
// a wołający odkłada bajty do magazynu zasobów po odczycie.
func (s *ooxmlStanOdczytu) czytajObraz(wezelRysunku *ooxmlWezel) {
	odwolanie := ooxmlPierwszeOdwolanieObrazu(wezelRysunku)
	if odwolanie == "" {
		s.pomin("obraz osadzony bez odwołania do składnika archiwum", "")
		return
	}
	sciezka := s.powiazania[odwolanie]
	if sciezka == "" {
		s.pomin("obraz osadzony wskazuje powiązanie, którego archiwum nie niesie", odwolanie)
		return
	}
	obiekt := shared.StudioDocumentObject{
		Id:     nowyIdentyfikator(przedrostekObiektuStudia),
		Kind:   shared.StudioObjectKindImage,
		Source: wejscieWskaznikZrodlaObiektu(shared.StudioObjectSourceFile),
		// Nazwa składnika archiwum jedzie w tekście zastępczym, dopóki
		// wołający nie wpisze tu kodu zasobu.
		AltText: wejscieWskaznikTekstu(ooxmlOpisObrazu(wezelRysunku, sciezka)),
	}
	if szerokosc, wysokosc, jest := ooxmlWymiaryRysunku(wezelRysunku); jest {
		obiekt.WidthMm = wejscieWskaznikRzeczywisty(szerokosc)
		obiekt.HeightMm = wejscieWskaznikRzeczywisty(wysokosc)
	}
	s.postac.Objects = append(s.postac.Objects, obiekt)
	s.postac.Blocks = append(s.postac.Blocks, shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuObiekt,
		SectionId: wejscieWskaznikTekstu(s.biezacaSekcja),
		ObjectId:  wejscieWskaznikTekstu(obiekt.Id),
	})
	s.bilans.ImagesEmbedded = wejscieWskaznikCalkowity(wartoscCalkowita(s.bilans.ImagesEmbedded) + 1)
}

// czytajOdsylacz dokłada odsyłacz do aparatu dokumentu wraz z jego adresem
// docelowym i etykietą złożoną z tekstu fragmentów, które opakowuje.
func (s *ooxmlStanOdczytu) czytajOdsylacz(wezel *ooxmlWezel, fragmenty []shared.StudioDocumentRun) {
	odwolanie := ""
	for _, atrybut := range wezel.Atrybuty {
		if atrybut.Name.Local == "id" && strings.Contains(atrybut.Name.Space, "relationships") {
			odwolanie = atrybut.Value
		}
	}
	adres := s.powiazania[odwolanie]
	if adres == "" {
		adres = wezel.atrybut("anchor")
	}
	if adres == "" {
		return
	}
	etykieta := make([]string, 0, len(fragmenty))
	for _, fragment := range fragmenty {
		etykieta = append(etykieta, fragment.Text)
	}
	s.postac.Apparatus = append(s.postac.Apparatus, shared.StudioApparatusItem{
		Id:        nowyIdentyfikator(przedrostekAparatuWejscia),
		Kind:      shared.StudioApparatusKindHyperlink,
		Label:     wejscieWskaznikTekstu(strings.Join(etykieta, "")),
		TargetUrl: wejscieWskaznikTekstu(adres),
	})
}

// przedrostekAparatuWejscia znakuje elementy aparatu odczytane z pliku,
// odróżniając je od elementów, które powstają dopiero przy złożeniu dokumentu.
const przedrostekAparatuWejscia = "studio-apar-"

// pomin dokłada pozycję do wykazu tego, czego rachunek nie odzyskał, wraz
// z powodem i szczegółem, które trafiają do bilansu wniesienia dokumentu.
func (s *ooxmlStanOdczytu) pomin(powod, szczegol string) {
	s.bilans.Skipped = append(s.bilans.Skipped, shared.StudioSkippedItem{
		Reason: powod, Detail: wejscieWskaznikTekstu(szczegol),
	})
}

// czytajTabele składa tabelę dokumentu wraz z komórkami, scaleniami,
// szerokościami kolumn, obramowaniem i wierszem nagłówkowym. Komórka
// wchłonięta scaleniem pionowym wchodzi do wykazu jako scalona, inaczej
// tabela zgubiłaby kolumnę przy zapisie.
func (s *ooxmlStanOdczytu) czytajTabele(wezel *ooxmlWezel) {
	wiersze := wezel.dzieci("tr")
	tabela := shared.StudioDocumentTable{
		Id:   nowyIdentyfikator(przedrostekTabeliStudia),
		Rows: len(wiersze),
	}
	if styl := wezel.sciezka("tblPr", "tblStyle"); styl != nil {
		tabela.StyleName = wejscieWskaznikTekstu(styl.atrybut("val"))
	}
	if obramowanie := ooxmlCzytajObramowanie(wezel.sciezka("tblPr", "tblBorders")); obramowanie != nil {
		tabela.Border = obramowanie
	}
	if szerokosc, jest := wezel.sciezka("tblPr", "tblW").atrybutCalkowity("w"); jest && szerokosc > 0 {
		tabela.WidthMm = wejscieWskaznikRzeczywisty(ooxmlTwipyNaMilimetry(szerokosc))
	}

	// Szerokości kolumn biorą się ze siatki tabeli (`w:tblGrid`), nie
	// z wiersza scalonego poziomo.
	for _, kolumna := range wezel.sciezka("tblGrid").dzieci("gridCol") {
		if twipy, jest := kolumna.atrybutCalkowity("w"); jest {
			tabela.ColumnWidthsMm = append(tabela.ColumnWidthsMm,
				ooxmlTwipyNaMilimetry(twipy))
		}
	}

	liczbaNaglowkow := 0
	for numerWiersza, wiersz := range wiersze {
		naglowkowy := wiersz.sciezka("trPr", "tblHeader") != nil
		if naglowkowy && numerWiersza == liczbaNaglowkow {
			liczbaNaglowkow++
		}
		kolumna := 0
		for _, komorka := range wiersz.dzieci("tc") {
			zlozona := ooxmlCzytajKomorke(komorka, numerWiersza, kolumna)
			tabela.Cells = append(tabela.Cells, zlozona)
			rozpietosc := 1
			if zlozona.ColumnSpan != nil && *zlozona.ColumnSpan > 1 {
				rozpietosc = *zlozona.ColumnSpan
			}
			// Kolumny wchłonięte scaleniem poziomym dostają własne komórki
			// oznaczone jako scalone.
			for przesuniecie := 1; przesuniecie < rozpietosc; przesuniecie++ {
				tabela.Cells = append(tabela.Cells, shared.StudioTableCell{
					Row: numerWiersza, Column: kolumna + przesuniecie,
					Merged: wejscieWskaznikLogiczny(true),
				})
			}
			kolumna += rozpietosc
		}
		if kolumna > tabela.Columns {
			tabela.Columns = kolumna
		}
	}
	if liczbaNaglowkow > 0 {
		tabela.HeaderRows = wejscieWskaznikCalkowity(liczbaNaglowkow)
		tabela.RepeatHeader = wejscieWskaznikLogiczny(true)
	}
	if len(tabela.ColumnWidthsMm) == 0 && tabela.Columns > 0 {
		// Szerokości nieobecne w pliku liczy się z podziału obszaru pisania,
		// nie z zera ani z braku.
		szerokoscObszaru := ooxmlSzerokoscObszaruPisania(s.postac)
		rowna := szerokoscObszaru / float64(tabela.Columns)
		for i := 0; i < tabela.Columns; i++ {
			tabela.ColumnWidthsMm = append(tabela.ColumnWidthsMm, rowna)
		}
	}

	s.postac.Tables = append(s.postac.Tables, tabela)
	s.postac.Blocks = append(s.postac.Blocks, shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuTabela,
		SectionId: wejscieWskaznikTekstu(s.biezacaSekcja),
		TableId:   wejscieWskaznikTekstu(tabela.Id),
	})
}

// ooxmlSzerokoscObszaruPisania liczy szerokość obszaru pisania z nastaw strony.
// Nastawy nieznane dają szerokość A4 z marginesami 25 mm.
func ooxmlSzerokoscObszaruPisania(postac *shared.StudioDocumentForm) float64 {
	szerokosc, lewy, prawy := 210.0, 25.0, 25.0
	nastawy := postac.PageSetup
	if nastawy == nil && len(postac.Sections) > 0 {
		nastawy = postac.Sections[len(postac.Sections)-1].PageSetup
	}
	if nastawy != nil {
		if nastawy.WidthMm != nil && *nastawy.WidthMm > 0 {
			szerokosc = *nastawy.WidthMm
		} else if nastawy.PageSize != nil {
			if wymiar, jest := wejscieWymiaryNosnika(*nastawy.PageSize); jest {
				szerokosc = wymiar.szerokosc
			}
		}
		if nastawy.Orientation != nil && *nastawy.Orientation == shared.StudioPageOrientationPozioma {
			if wymiar, jest := wejscieWymiaryNosnika(wartoscTekstu(nastawy.PageSize)); jest {
				szerokosc = wymiar.wysokosc
			}
		}
		if nastawy.MarginLeft != nil {
			lewy = float64(*nastawy.MarginLeft)
		}
		if nastawy.MarginRight != nil {
			prawy = float64(*nastawy.MarginRight)
		}
	}
	obszar := szerokosc - lewy - prawy
	if obszar < 10 {
		return 10
	}
	return obszar
}

// ooxmlCzytajKomorke składa komórkę tabeli wraz z jej postacią, obramowaniem,
// scaleniem, cieniowaniem i treścią akapitów, które niesie.
func ooxmlCzytajKomorke(wezel *ooxmlWezel, wiersz, kolumna int) shared.StudioTableCell {
	komorka := shared.StudioTableCell{Row: wiersz, Column: kolumna}
	postacKomorki := wezel.dziecko("tcPr")
	if rozpietosc, jest := postacKomorki.sciezka("gridSpan").atrybutCalkowity("val"); jest && rozpietosc > 1 {
		komorka.ColumnSpan = wejscieWskaznikCalkowity(rozpietosc)
	}
	if scalenie := postacKomorki.sciezka("vMerge"); scalenie != nil {
		if strings.EqualFold(scalenie.atrybut("val"), "restart") {
			komorka.RowSpan = wejscieWskaznikCalkowity(2)
		} else {
			komorka.Merged = wejscieWskaznikLogiczny(true)
		}
	}
	if wyrownanie := postacKomorki.sciezka("vAlign"); wyrownanie != nil {
		komorka.VerticalAlign = ooxmlWyrownaniePionowe(wyrownanie.atrybut("val"))
	}
	if cieniowanie := postacKomorki.sciezka("shd"); cieniowanie != nil {
		if barwa := ooxmlBarwa(cieniowanie.atrybut("fill")); barwa != "" {
			komorka.ShadingColor = wejscieWskaznikTekstu(barwa)
		}
	}
	if obramowanie := ooxmlCzytajObramowanie(postacKomorki.sciezka("tcBorders")); obramowanie != nil {
		komorka.Border = obramowanie
	}

	// Treść komórki to akapity złożone końcem wiersza; postać pierwszego
	// staje się postacią komórki.
	akapity := wezel.dzieci("p")
	czesci := make([]string, 0, len(akapity))
	for i, akapit := range akapity {
		czesci = append(czesci, ooxmlTekstAkapitu(akapit))
		if i == 0 {
			komorka.Paragraph = ooxmlCzytajPostacAkapitu(akapit.dziecko("pPr"))
			if fragment := akapit.dziecko("r"); fragment != nil {
				komorka.Character = ooxmlCzytajPostacZnaku(fragment.dziecko("rPr"))
			}
		}
	}
	komorka.Text = wejscieWskaznikTekstu(strings.Join(czesci, "\n"))
	return komorka
}

// ooxmlTekstAkapitu składa treść akapitu z jego fragmentów, schodząc w głąb
// odsyłacza, kontrolki treści i zmiany śledzonej, które fragmenty opakowują.
func ooxmlTekstAkapitu(akapit *ooxmlWezel) string {
	var budowa strings.Builder
	for _, dziecko := range akapit.Dzieci {
		switch dziecko.Nazwa.Local {
		case "r":
			budowa.WriteString(ooxmlTekstFragmentu(dziecko))
		case "hyperlink", "ins", "smartTag", "sdt":
			budowa.WriteString(ooxmlTekstAkapitu(dziecko))
		}
	}
	return budowa.String()
}

// ooxmlTekstFragmentu składa treść jednego fragmentu. `w:tab` staje
// tabulatorem, `w:br` znakiem końca wiersza, `w:delText` nie wchodzi wcale —
// to treść usunięta zmianą śledzoną, której w dokumencie już nie ma.
func ooxmlTekstFragmentu(fragment *ooxmlWezel) string {
	var budowa strings.Builder
	for _, dziecko := range fragment.Dzieci {
		switch dziecko.Nazwa.Local {
		case "t":
			budowa.WriteString(dziecko.Tekst)
		case "tab":
			budowa.WriteString("\t")
		case "br", "cr":
			budowa.WriteString("\n")
		case "noBreakHyphen":
			budowa.WriteString("‑")
		case "softHyphen":
			budowa.WriteString("­")
		case "sym":
			// Symbol specjalny stoi kodem szesnastkowym w atrybucie `w:char`.
			if kod, err := strconv.ParseInt(dziecko.atrybut("char"), 16, 32); err == nil {
				budowa.WriteRune(rune(kod))
			}
		}
	}
	return budowa.String()
}

// ── Odczyt postaci ─────────────────────────────────────────────────────────

// ooxmlCzytajPostacZnaku składa postać znaku z `w:rPr`: krój, stopień,
// pogrubienie, kursywę, podkreślenie, barwę i pozostałe przełączniki
// formatowania.
func ooxmlCzytajPostacZnaku(wezel *ooxmlWezel) *shared.StudioCharacterFormat {
	if wezel == nil {
		return nil
	}
	postac := shared.StudioCharacterFormat{}
	pelna := false

	if kroje := wezel.dziecko("rFonts"); kroje != nil {
		// Krój bierze się z `w:ascii`, a gdy go nie ma — z `w:hAnsi`; oba
		// opisują pismo łacińskie.
		krój := kroje.atrybut("ascii")
		if krój == "" {
			krój = kroje.atrybut("hAnsi")
		}
		if krój != "" {
			postac.FontFamily, pelna = wejscieWskaznikTekstu(krój), true
		}
	}
	if stopien, jest := wezel.sciezka("sz").atrybutCalkowity("val"); jest {
		postac.FontSizePt, pelna = wejscieWskaznikRzeczywisty(ooxmlPolpunktyNaPunkty(stopien)), true
	}
	if wezel := wezel.dziecko("b"); wezel != nil {
		postac.Bold, pelna = wejscieWskaznikLogiczny(wezel.wlaczonyPrzelacznik()), true
	}
	if wezel := wezel.dziecko("i"); wezel != nil {
		postac.Italic, pelna = wejscieWskaznikLogiczny(wezel.wlaczonyPrzelacznik()), true
	}
	if wezel := wezel.dziecko("strike"); wezel != nil {
		postac.Strikethrough, pelna = wejscieWskaznikLogiczny(wezel.wlaczonyPrzelacznik()), true
	}
	if wezel := wezel.dziecko("smallCaps"); wezel != nil {
		postac.SmallCaps, pelna = wejscieWskaznikLogiczny(wezel.wlaczonyPrzelacznik()), true
	}
	if wezel := wezel.dziecko("caps"); wezel != nil {
		postac.AllCaps, pelna = wejscieWskaznikLogiczny(wezel.wlaczonyPrzelacznik()), true
	}
	if podkreslenie := wezel.dziecko("u"); podkreslenie != nil {
		postac.Underline, pelna = ooxmlPodkreslenie(podkreslenie.atrybut("val")), true
	}
	if pion := wezel.dziecko("vertAlign"); pion != nil {
		switch strings.ToLower(pion.atrybut("val")) {
		case "superscript":
			postac.Superscript, pelna = wejscieWskaznikLogiczny(true), true
		case "subscript":
			postac.Subscript, pelna = wejscieWskaznikLogiczny(true), true
		}
	}
	if barwa := ooxmlBarwa(wezel.sciezka("color").atrybut("val")); barwa != "" {
		postac.Color, pelna = wejscieWskaznikTekstu(barwa), true
	}
	if wyroznienie := wezel.dziecko("highlight"); wyroznienie != nil {
		if barwa := ooxmlBarwaNazwana(wyroznienie.atrybut("val")); barwa != "" {
			postac.HighlightColor, pelna = wejscieWskaznikTekstu(barwa), true
		}
	}
	if cieniowanie := wezel.dziecko("shd"); cieniowanie != nil && postac.HighlightColor == nil {
		if barwa := ooxmlBarwa(cieniowanie.atrybut("fill")); barwa != "" {
			postac.HighlightColor, pelna = wejscieWskaznikTekstu(barwa), true
		}
	}
	if odstep, jest := wezel.sciezka("spacing").atrybutCalkowity("val"); jest {
		postac.LetterSpacingPt, pelna = wejscieWskaznikRzeczywisty(ooxmlTwipyNaPunkty(odstep)), true
	}
	if styl := wezel.sciezka("rStyle"); styl != nil {
		postac.StyleName, pelna = wejscieWskaznikTekstu(styl.atrybut("val")), true
	}
	if jezyk := wezel.sciezka("lang"); jezyk != nil {
		if wartosc := jezyk.atrybut("val"); wartosc != "" {
			postac.Language, pelna = wejscieWskaznikTekstu(wartosc), true
		}
	}
	if !pelna {
		return nil
	}
	return &postac
}

// ooxmlCzytajPostacAkapitu składa postać akapitu z `w:pPr`: styl, wyrównanie,
// wcięcia, odstępy, tabulatory, obramowanie i przynależność do listy.
func ooxmlCzytajPostacAkapitu(wezel *ooxmlWezel) *shared.StudioParagraphFormat {
	if wezel == nil {
		return nil
	}
	postac := shared.StudioParagraphFormat{}
	if styl := wezel.sciezka("pStyle"); styl != nil {
		postac.StyleName = wejscieWskaznikTekstu(ooxmlNazwaStyluDokumentu(styl.atrybut("val")))
	}
	if wyrownanie := wezel.dziecko("jc"); wyrownanie != nil {
		postac.Align = ooxmlWyrownanie(wyrownanie.atrybut("val"))
	}
	if wciecia := wezel.dziecko("ind"); wciecia != nil {
		if wartosc, jest := wciecia.atrybutCalkowity("left"); jest {
			postac.IndentLeftMm = wejscieWskaznikRzeczywisty(ooxmlTwipyNaMilimetry(wartosc))
		}
		if wartosc, jest := wciecia.atrybutCalkowity("right"); jest {
			postac.IndentRightMm = wejscieWskaznikRzeczywisty(ooxmlTwipyNaMilimetry(wartosc))
		}
		if wartosc, jest := wciecia.atrybutCalkowity("firstLine"); jest {
			postac.FirstLineIndentMm = wejscieWskaznikRzeczywisty(ooxmlTwipyNaMilimetry(wartosc))
		}
		// Wysunięcie (`w:hanging`) jest wcięciem pierwszego wiersza o wartości
		// UJEMNEJ w kontrakcie.
		if wartosc, jest := wciecia.atrybutCalkowity("hanging"); jest {
			postac.FirstLineIndentMm = wejscieWskaznikRzeczywisty(-ooxmlTwipyNaMilimetry(wartosc))
		}
	}
	if odstepy := wezel.dziecko("spacing"); odstepy != nil {
		if wartosc, jest := odstepy.atrybutCalkowity("before"); jest {
			postac.SpaceBeforePt = wejscieWskaznikRzeczywisty(ooxmlTwipyNaPunkty(wartosc))
		}
		if wartosc, jest := odstepy.atrybutCalkowity("after"); jest {
			postac.SpaceAfterPt = wejscieWskaznikRzeczywisty(ooxmlTwipyNaPunkty(wartosc))
		}
		if wartosc, jest := odstepy.atrybutCalkowity("line"); jest {
			zasada, liczba := ooxmlInterlinia(odstepy.atrybut("lineRule"), wartosc)
			postac.LineSpacingRule = zasada
			postac.LineSpacingValue = wejscieWskaznikRzeczywisty(liczba)
		}
	}
	if tabulatory := ooxmlCzytajTabulatory(wezel.dziecko("tabs")); len(tabulatory) > 0 {
		postac.TabStops = tabulatory
	}
	if obramowanie := ooxmlCzytajObramowanie(wezel.dziecko("pBdr")); obramowanie != nil {
		postac.Border = obramowanie
	}
	if cieniowanie := wezel.dziecko("shd"); cieniowanie != nil {
		if barwa := ooxmlBarwa(cieniowanie.atrybut("fill")); barwa != "" {
			postac.ShadingColor = wejscieWskaznikTekstu(barwa)
		}
	}
	if wezel := wezel.dziecko("widowControl"); wezel != nil {
		postac.WidowControl = wejscieWskaznikLogiczny(wezel.wlaczonyPrzelacznik())
	}
	if wezel := wezel.dziecko("keepNext"); wezel != nil {
		postac.KeepWithNext = wejscieWskaznikLogiczny(wezel.wlaczonyPrzelacznik())
	}
	if wezel := wezel.dziecko("keepLines"); wezel != nil {
		postac.KeepLines = wejscieWskaznikLogiczny(wezel.wlaczonyPrzelacznik())
	}
	if wezel := wezel.dziecko("bidi"); wezel != nil {
		postac.RightToLeft = wejscieWskaznikLogiczny(wezel.wlaczonyPrzelacznik())
	}
	if poziom, jest := wezel.sciezka("outlineLvl").atrybutCalkowity("val"); jest {
		// OOXML liczy poziomy konspektu od zera, kontrakt platformy od
		// jednego, stąd przesunięcie.
		postac.OutlineLevel = wejscieWskaznikCalkowity(poziom + 1)
	}
	if numeracja := wezel.dziecko("numPr"); numeracja != nil {
		if kod, jest := numeracja.sciezka("numId").atrybutCalkowity("val"); jest {
			postac.ListId = wejscieWskaznikTekstu(ooxmlNazwaListy(kod))
		}
		if poziom, jest := numeracja.sciezka("ilvl").atrybutCalkowity("val"); jest {
			postac.ListLevel = wejscieWskaznikCalkowity(poziom + 1)
		}
	}
	return &postac
}

// ooxmlCzytajTabulatory składa tabulatory akapitu wraz z ich położeniem,
// rodzajem wyrównania i znakiem wiodącym wypełniającym odstęp do tabulatora.
func ooxmlCzytajTabulatory(wezel *ooxmlWezel) []shared.StudioTabStop {
	if wezel == nil {
		return nil
	}
	tabulatory := make([]shared.StudioTabStop, 0, len(wezel.Dzieci))
	for _, tabulator := range wezel.dzieci("tab") {
		polozenie, jest := tabulator.atrybutCalkowity("pos")
		if !jest {
			continue
		}
		rodzaj := shared.StudioTabKindLeft
		switch strings.ToLower(tabulator.atrybut("val")) {
		case "right":
			rodzaj = shared.StudioTabKindRight
		case "center":
			rodzaj = shared.StudioTabKindCenter
		case "decimal":
			rodzaj = shared.StudioTabKindDecimal
		case "bar":
			rodzaj = shared.StudioTabKindBar
		case "clear":
			continue
		}
		zlozony := shared.StudioTabStop{
			PositionMm: ooxmlTwipyNaMilimetry(polozenie),
			Kind:       shared.StudioTabKind(rodzaj),
		}
		if znak := ooxmlZnakWiodacy(tabulator.atrybut("leader")); znak != nil {
			zlozony.Leader = znak
		}
		tabulatory = append(tabulatory, zlozony)
	}
	if len(tabulatory) == 0 {
		return nil
	}
	return tabulatory
}

// ooxmlCzytajObramowanie składa obramowanie z węzła krawędzi. Kontrakt niesie
// JEDNO obramowanie, a OOXML osobny opis każdej krawędzi — odmiana i grubość
// biorą się z pierwszej krawędzi, która je ma.
func ooxmlCzytajObramowanie(wezel *ooxmlWezel) *shared.StudioBorder {
	if wezel == nil {
		return nil
	}
	obramowanie := shared.StudioBorder{Style: shared.StudioBorderStyleNone}
	nazwy := []string{"top", "bottom", "left", "right"}
	przelaczniki := []**bool{
		&obramowanie.Top, &obramowanie.Bottom, &obramowanie.Left, &obramowanie.Right,
	}
	cokolwiek := false
	for i, nazwa := range nazwy {
		krawedz := wezel.dziecko(nazwa)
		if krawedz == nil {
			continue
		}
		odmiana := ooxmlOdmianaObramowania(krawedz.atrybut("val"))
		if odmiana == shared.StudioBorderStyleNone {
			*przelaczniki[i] = wejscieWskaznikLogiczny(false)
			cokolwiek = true
			continue
		}
		*przelaczniki[i] = wejscieWskaznikLogiczny(true)
		cokolwiek = true
		if obramowanie.Style == shared.StudioBorderStyleNone {
			obramowanie.Style = shared.StudioBorderStyle(odmiana)
		}
		if obramowanie.WidthPt == nil {
			if grubosc, jest := krawedz.atrybutCalkowity("sz"); jest {
				// `w:sz` obramowania liczy się w ósemkach punktu.
				obramowanie.WidthPt = wejscieWskaznikRzeczywisty(float64(grubosc) / 8)
			}
		}
		if obramowanie.Color == nil {
			if barwa := ooxmlBarwa(krawedz.atrybut("color")); barwa != "" {
				obramowanie.Color = wejscieWskaznikTekstu(barwa)
			}
		}
	}
	if !cokolwiek {
		return nil
	}
	return &obramowanie
}

// ooxmlCzytajNastawyStrony składa nastawy strony z `w:sectPr`: rozmiar
// nośnika, orientację, marginesy, kolumny szpaltowe i odbicie marginesów.
func ooxmlCzytajNastawyStrony(wezel *ooxmlWezel) shared.StudioPageSetup {
	nastawy := wejscieDomyslneNastawyStrony("", nil)
	if wezel == nil {
		return nastawy
	}
	if rozmiar := wezel.dziecko("pgSz"); rozmiar != nil {
		szerokosc, maSzerokosc := rozmiar.atrybutCalkowity("w")
		wysokosc, maWysokosc := rozmiar.atrybutCalkowity("h")
		if maSzerokosc && maWysokosc {
			szerokoscMm := ooxmlTwipyNaMilimetry(szerokosc)
			wysokoscMm := ooxmlTwipyNaMilimetry(wysokosc)
			nastawy.WidthMm = wejscieWskaznikRzeczywisty(szerokoscMm)
			nastawy.HeightMm = wejscieWskaznikRzeczywisty(wysokoscMm)
			// Nazwa nośnika bierze się z wymiarów, bo OOXML nazwy nośnika nie
			// niesie wcale.
			if nazwa, jest := wejscieNazwaNosnikaZWymiarow(szerokoscMm, wysokoscMm); jest {
				nastawy.PageSize = wejscieWskaznikTekstu(nazwa)
			} else {
				nastawy.PageSize = wejscieWskaznikTekstu("własny")
				nastawy.PaperKind = wejscieWskaznikRodzajuNosnika(shared.StudioPaperKindCustom)
			}
		}
		if strings.EqualFold(rozmiar.atrybut("orient"), "landscape") {
			nastawy.Orientation = wejscieWskaznikOrientacji(shared.StudioPageOrientationPozioma)
		}
	}
	if marginesy := wezel.dziecko("pgMar"); marginesy != nil {
		if wartosc, jest := marginesy.atrybutCalkowity("top"); jest {
			nastawy.MarginTop = wejscieWskaznikCalkowity(int(ooxmlTwipyNaMilimetry(wartosc) + 0.5))
		}
		if wartosc, jest := marginesy.atrybutCalkowity("bottom"); jest {
			nastawy.MarginBottom = wejscieWskaznikCalkowity(int(ooxmlTwipyNaMilimetry(wartosc) + 0.5))
		}
		if wartosc, jest := marginesy.atrybutCalkowity("left"); jest {
			nastawy.MarginLeft = wejscieWskaznikCalkowity(int(ooxmlTwipyNaMilimetry(wartosc) + 0.5))
		}
		if wartosc, jest := marginesy.atrybutCalkowity("right"); jest {
			nastawy.MarginRight = wejscieWskaznikCalkowity(int(ooxmlTwipyNaMilimetry(wartosc) + 0.5))
		}
		if wartosc, jest := marginesy.atrybutCalkowity("gutter"); jest && wartosc > 0 {
			nastawy.GutterMm = wejscieWskaznikRzeczywisty(ooxmlTwipyNaMilimetry(wartosc))
		}
		nastawy.MarginPreset = wejscieWskaznikTekstu("własne")
	}
	if kolumny := wezel.dziecko("cols"); kolumny != nil {
		if liczba, jest := kolumny.atrybutCalkowity("num"); jest && liczba > 0 {
			nastawy.Columns = wejscieWskaznikCalkowity(liczba)
		}
		if odstep, jest := kolumny.atrybutCalkowity("space"); jest && odstep > 0 {
			nastawy.ColumnGapMm = wejscieWskaznikRzeczywisty(ooxmlTwipyNaMilimetry(odstep))
		}
		if kreska := kolumny.atrybut("sep"); kreska != "" {
			nastawy.ColumnRule = wejscieWskaznikLogiczny(kolumny.dziecko("sep").wlaczonyPrzelacznik() ||
				strings.EqualFold(kreska, "1") || strings.EqualFold(kreska, "true"))
		}
	}
	if odbicie := wezel.dziecko("mirrorMargins"); odbicie != nil {
		nastawy.MirrorMargins = wejscieWskaznikLogiczny(odbicie.wlaczonyPrzelacznik())
	}
	return nastawy
}

// ooxmlPoczatekSekcji rozstrzyga, jak sekcja się zaczyna: ciągiem, nową
// stroną, stroną parzystą, nieparzystą albo nową szpaltą.
func ooxmlPoczatekSekcji(wezel *ooxmlWezel) shared.StudioSectionStart {
	if wezel == nil {
		return shared.StudioSectionStartContinuous
	}
	switch strings.ToLower(wezel.sciezka("type").atrybut("val")) {
	case "nextPage", "nextpage":
		return shared.StudioSectionStartNewPage
	case "evenpage":
		return shared.StudioSectionStartEvenPage
	case "oddpage":
		return shared.StudioSectionStartOddPage
	case "nextcolumn":
		return shared.StudioSectionStartNewColumn
	}
	return shared.StudioSectionStartContinuous
}

// ooxmlCzytajNumeracjeStron składa numerację stron sekcji: numer początkowy
// i format zapisu, gdy `w:sectPr` niesie węzeł `w:pgNumType`.
func ooxmlCzytajNumeracjeStron(wezel *ooxmlWezel) *shared.StudioPageNumbering {
	if wezel == nil {
		return nil
	}
	rodzaj := wezel.dziecko("pgNumType")
	if rodzaj == nil {
		return nil
	}
	numeracja := shared.StudioPageNumbering{Enabled: true}
	if start, jest := rodzaj.atrybutCalkowity("start"); jest {
		numeracja.StartAt = wejscieWskaznikCalkowity(start)
	}
	if format := ooxmlFormatNumeracjiStron(rodzaj.atrybut("fmt")); format != nil {
		numeracja.Format = format
	}
	if numeracja.StartAt == nil && numeracja.Format == nil {
		return nil
	}
	return &numeracja
}

// ooxmlCzytajNaglowkiStopki wczytuje nagłówki i stopki sekcji wraz
// z zasięgiem: domyślnym, pierwszej strony i stron parzystych — w pismach
// urzędowych te trzy bywają różne.
func ooxmlCzytajNaglowkiStopki(wezelSekcji *ooxmlWezel, skladniki map[string][]byte,
	powiazania map[string]string) []shared.StudioHeaderFooter {

	if wezelSekcji == nil {
		return nil
	}
	wynik := make([]shared.StudioHeaderFooter, 0, 6)
	rodzaje := []struct {
		wezel    string
		naglowek bool
	}{{"headerReference", true}, {"footerReference", false}}

	for _, rodzaj := range rodzaje {
		for _, odwolanie := range wezelSekcji.dzieci(rodzaj.wezel) {
			kod := ""
			for _, atrybut := range odwolanie.Atrybuty {
				if atrybut.Name.Local == "id" &&
					strings.Contains(atrybut.Name.Space, "relationships") {
					kod = atrybut.Value
				}
			}
			sciezka := powiazania[kod]
			if sciezka == "" {
				continue
			}
			surowy, jest := skladniki[path.Join("word", sciezka)]
			if !jest {
				surowy, jest = skladniki[sciezka]
			}
			if !jest {
				continue
			}
			drzewo, err := ooxmlRozbierz(surowy)
			if err != nil {
				continue
			}
			tresc := make([]string, 0, len(drzewo.Dzieci))
			for _, akapit := range drzewo.dzieci("p") {
				tresc = append(tresc, ooxmlTekstAkapitu(akapit))
			}
			zlozony := shared.StudioHeaderFooter{
				Scope: ooxmlZasiegNaglowka(odwolanie.atrybut("type")),
			}
			if rodzaj.naglowek {
				zlozony.HeaderText = wejscieWskaznikTekstu(strings.Join(tresc, "\n"))
			} else {
				zlozony.FooterText = wejscieWskaznikTekstu(strings.Join(tresc, "\n"))
			}
			wynik = append(wynik, zlozony)
		}
	}
	if len(wynik) == 0 {
		return nil
	}
	return wynik
}

// ooxmlCzytajPowiazania czyta wykaz powiązań dokumentu z
// `word/_rels/document.xml.rels`: identyfikator powiązania na jego cel.
func ooxmlCzytajPowiazania(surowe []byte) map[string]string {
	powiazania := make(map[string]string)
	if len(surowe) == 0 {
		return powiazania
	}
	drzewo, err := ooxmlRozbierz(surowe)
	if err != nil {
		return powiazania
	}
	for _, wpis := range drzewo.dzieci("Relationship") {
		kod := wpis.atrybut("Id")
		cel := wpis.atrybut("Target")
		if kod == "" || cel == "" {
			continue
		}
		powiazania[kod] = cel
	}
	return powiazania
}

// ooxmlCzytajStyle składa arkusz stylów nazwanych z `word/styles.xml`: postać
// akapitu i znaku, styl nadrzędny oraz styl następny.
func ooxmlCzytajStyle(surowe []byte) []shared.StudioNamedStyle {
	if len(surowe) == 0 {
		return nil
	}
	drzewo, err := ooxmlRozbierz(surowe)
	if err != nil {
		return nil
	}
	arkusz := make([]shared.StudioNamedStyle, 0, len(drzewo.Dzieci))
	for _, wpis := range drzewo.dzieci("style") {
		kod := wpis.atrybut("styleId")
		if kod == "" {
			continue
		}
		nazwaWidoczna := wpis.sciezka("name").atrybut("val")
		styl := shared.StudioNamedStyle{
			Name:        ooxmlNazwaStyluDokumentu(kod),
			DisplayName: wejscieWskaznikTekstu(nazwaWidoczna),
			Kind:        ooxmlRodzajStylu(wpis.atrybut("type")),
			Character:   ooxmlCzytajPostacZnaku(wpis.dziecko("rPr")),
			Paragraph:   ooxmlCzytajPostacAkapitu(wpis.dziecko("pPr")),
		}
		if nadrzedny := wpis.sciezka("basedOn"); nadrzedny != nil {
			styl.BasedOn = wejscieWskaznikTekstu(ooxmlNazwaStyluDokumentu(nadrzedny.atrybut("val")))
		}
		if nastepny := wpis.sciezka("next"); nastepny != nil {
			styl.NextStyle = wejscieWskaznikTekstu(ooxmlNazwaStyluDokumentu(nastepny.atrybut("val")))
		}
		if wpis.dziecko("semiHidden") != nil || strings.EqualFold(wpis.atrybut("default"), "1") {
			styl.Builtin = wejscieWskaznikLogiczny(true)
		}
		arkusz = append(arkusz, styl)
	}
	if len(arkusz) == 0 {
		return nil
	}
	return arkusz
}

// ooxmlCzytajNumeracje składa definicje list z `word/numbering.xml`: poziomy,
// format numeracji lub wypunktowania oraz wcięcia.
func ooxmlCzytajNumeracje(surowe []byte) []shared.StudioListDefinition {
	if len(surowe) == 0 {
		return nil
	}
	drzewo, err := ooxmlRozbierz(surowe)
	if err != nil {
		return nil
	}
	// `w:num` wiąże numerację dokumentu z abstrakcyjną definicją poziomów
	// w `w:abstractNum`.
	abstrakcyjne := make(map[string]*ooxmlWezel)
	for _, wpis := range drzewo.dzieci("abstractNum") {
		abstrakcyjne[wpis.atrybut("abstractNumId")] = wpis
	}
	definicje := make([]shared.StudioListDefinition, 0, len(drzewo.Dzieci))
	for _, wpis := range drzewo.dzieci("num") {
		kod := wpis.atrybut("numId")
		wskazanie := wpis.sciezka("abstractNumId").atrybut("val")
		zrodlo, jest := abstrakcyjne[wskazanie]
		if !jest {
			continue
		}
		definicja := shared.StudioListDefinition{
			Id:   ooxmlNazwaListy(ooxmlLiczbaZNapisu(kod)),
			Kind: shared.StudioListKindNumber,
		}
		wypunktowanie, numeracja := 0, 0
		for _, poziom := range zrodlo.dzieci("lvl") {
			numerPoziomu, _ := poziom.atrybutCalkowity("ilvl")
			format := strings.ToLower(poziom.sciezka("numFmt").atrybut("val"))
			zlozony := shared.StudioListLevel{Level: numerPoziomu + 1}
			if format == "bullet" {
				wypunktowanie++
				zlozony.BulletSource = ooxmlWskaznikZrodlaWypunktowania(shared.StudioBulletSourceCharacter)
				zlozony.BulletCharacter = wejscieWskaznikTekstu(poziom.sciezka("lvlText").atrybut("val"))
			} else {
				numeracja++
				zlozony.NumberFormat = ooxmlFormatNumeracjiListy(format)
				zlozony.Pattern = wejscieWskaznikTekstu(poziom.sciezka("lvlText").atrybut("val"))
			}
			if start, maStart := poziom.sciezka("start").atrybutCalkowity("val"); maStart {
				zlozony.StartAt = wejscieWskaznikCalkowity(start)
			}
			if wciecia := poziom.sciezka("pPr", "ind"); wciecia != nil {
				if wartosc, maWartosc := wciecia.atrybutCalkowity("left"); maWartosc {
					zlozony.IndentMm = wejscieWskaznikRzeczywisty(ooxmlTwipyNaMilimetry(wartosc))
				}
				if wartosc, maWartosc := wciecia.atrybutCalkowity("hanging"); maWartosc {
					zlozony.HangingMm = wejscieWskaznikRzeczywisty(ooxmlTwipyNaMilimetry(wartosc))
				}
			}
			definicja.Levels = append(definicja.Levels, zlozony)
		}
		switch {
		case len(definicja.Levels) > 1:
			definicja.Kind = shared.StudioListKindMultilevel
		case wypunktowanie > numeracja:
			definicja.Kind = shared.StudioListKindBullet
		}
		definicje = append(definicje, definicja)
	}
	if len(definicje) == 0 {
		return nil
	}
	return definicje
}

// ── Przekłady wartości ─────────────────────────────────────────────────────

// ooxmlNazwaStyluDokumentu sprowadza identyfikator stylu OOXML do nazwy, którą
// niesie arkusz stylów tego produktu. Nazwy nagłówków są przekładane, bo model
// i Operator pracują na nazwach pełnych, nie na identyfikatorach Worda.
func ooxmlNazwaStyluDokumentu(kod string) string {
	oczyszczony := strings.TrimSpace(kod)
	if oczyszczony == "" {
		return ""
	}
	maly := strings.ToLower(oczyszczony)
	switch maly {
	case "normal", "standard", "bodytext", "textbody", "default":
		return wejscieStylTekstZasadniczy
	case "quote", "intensequote", "blockquote":
		return wejscieStylCytat
	case "caption":
		return wejscieStylPodpis
	case "footnotetext", "endnotetext":
		return wejscieStylPrzypis
	}
	for poziom := 1; poziom <= 6; poziom++ {
		numer := strconv.Itoa(poziom)
		if maly == "heading"+numer || maly == "heading "+numer || maly == "naglowek"+numer {
			return wejscieStylNaglowkaPoziomu(poziom)
		}
	}
	return oczyszczony
}

// ooxmlNazwaListy składa nazwę definicji listy z numeru OOXML, poprzedzając
// go stałym przedrostkiem rozpoznawanym przy zapisie.
func ooxmlNazwaListy(numer int) string {
	return "lista-" + strconv.Itoa(numer)
}

// ooxmlLiczbaZNapisu czyta liczbę całkowitą z napisu identyfikatora OOXML;
// napis nieliczbowy albo pusty daje zero.
func ooxmlLiczbaZNapisu(tekst string) int {
	liczba, err := strconv.Atoi(strings.TrimSpace(tekst))
	if err != nil {
		return 0
	}
	return liczba
}

// ooxmlRodzajStylu przekłada rodzaj stylu OOXML na rodzaj kontraktu; wartość
// spoza wykazu daje rodzaj akapitowy jako domyślny.
func ooxmlRodzajStylu(rodzaj string) shared.StudioStyleKind {
	switch strings.ToLower(strings.TrimSpace(rodzaj)) {
	case "character":
		return shared.StudioStyleKindCharacter
	case "table":
		return shared.StudioStyleKindTable
	case "numbering":
		return shared.StudioStyleKindList
	}
	return shared.StudioStyleKindParagraph
}

// ooxmlWyrownanie przekłada wyrównanie akapitu OOXML na wyrównanie kontraktu;
// wartość spoza wykazu daje wskaźnik pusty.
func ooxmlWyrownanie(wartosc string) *shared.StudioTextAlign {
	switch strings.ToLower(strings.TrimSpace(wartosc)) {
	case "center":
		return wejscieWskaznikWyrownania(shared.StudioTextAlignCenter)
	case "right", "end":
		return wejscieWskaznikWyrownania(shared.StudioTextAlignRight)
	case "both", "justify", "distribute":
		return wejscieWskaznikWyrownania(shared.StudioTextAlignJustify)
	case "left", "start":
		return wejscieWskaznikWyrownania(shared.StudioTextAlignLeft)
	}
	return nil
}

// ooxmlWyrownaniePionowe przekłada wyrównanie pionowe komórki tabeli; wartość
// spoza wykazu i pusta dają wskaźnik pusty.
func ooxmlWyrownaniePionowe(wartosc string) *shared.StudioVerticalAlign {
	kopia := shared.StudioVerticalAlign(shared.StudioVerticalAlignTop)
	switch strings.ToLower(strings.TrimSpace(wartosc)) {
	case "center":
		kopia = shared.StudioVerticalAlignMiddle
	case "bottom":
		kopia = shared.StudioVerticalAlignBottom
	case "top":
	default:
		return nil
	}
	return &kopia
}

// ooxmlPodkreslenie przekłada odmianę podkreślenia OOXML na odmianę
// kontraktu; wartość spoza wykazu daje podkreślenie pojedyncze.
func ooxmlPodkreslenie(wartosc string) *shared.StudioUnderlineStyle {
	odmiana := shared.StudioUnderlineStyle(shared.StudioUnderlineStyleSingle)
	switch strings.ToLower(strings.TrimSpace(wartosc)) {
	case "none":
		odmiana = shared.StudioUnderlineStyleNone
	case "double":
		odmiana = shared.StudioUnderlineStyleDouble
	case "thick":
		odmiana = shared.StudioUnderlineStyleThick
	case "dotted", "dottedheavy":
		odmiana = shared.StudioUnderlineStyleDotted
	case "dash", "dashed", "dashedheavy", "dashlong":
		odmiana = shared.StudioUnderlineStyleDashed
	case "wave", "wavy", "wavyheavy", "wavydouble":
		odmiana = shared.StudioUnderlineStyleWavy
	case "words":
		odmiana = shared.StudioUnderlineStyleWords
	}
	return &odmiana
}

// ooxmlOdmianaObramowania przekłada odmianę obramowania OOXML na odmianę
// kontraktu; wartość spoza wykazu daje obramowanie pojedyncze.
func ooxmlOdmianaObramowania(wartosc string) shared.StudioBorderStyle {
	switch strings.ToLower(strings.TrimSpace(wartosc)) {
	case "", "none", "nil":
		return shared.StudioBorderStyleNone
	case "double", "doublewave":
		return shared.StudioBorderStyleDouble
	case "thick", "thickthinsmallgap", "triple":
		return shared.StudioBorderStyleThick
	case "dotted", "dotdash", "dotdotdash":
		return shared.StudioBorderStyleDotted
	case "dashed", "dashsmallgap", "dashdotstroked":
		return shared.StudioBorderStyleDashed
	}
	return shared.StudioBorderStyleSingle
}

// ooxmlZnakWiodacy przekłada znak wiodący tabulatora OOXML na znak
// kontraktu; wartość spoza wykazu daje wskaźnik pusty.
func ooxmlZnakWiodacy(wartosc string) *shared.StudioTabLeader {
	znak := shared.StudioTabLeader(shared.StudioTabLeaderNone)
	switch strings.ToLower(strings.TrimSpace(wartosc)) {
	case "dot", "middledot":
		znak = shared.StudioTabLeaderDot
	case "hyphen":
		znak = shared.StudioTabLeaderDash
	case "underscore", "heavy":
		znak = shared.StudioTabLeaderUnderline
	default:
		return nil
	}
	return &znak
}

// ooxmlInterlinia przekłada interlinię OOXML na zasadę i wartość kontraktu.
// `auto` liczy się w dwustu czterdziestych częściach wiersza: 240 znaczy
// jeden wiersz, 360 — półtora, 480 — dwa; pozostałe wartości zostają
// mnożnikiem.
func ooxmlInterlinia(zasada string, wartosc int) (*shared.StudioLineSpacingRule, float64) {
	switch strings.ToLower(strings.TrimSpace(zasada)) {
	case "atleast":
		return wejscieWskaznikZasadyInterlinii(shared.StudioLineSpacingRuleAtLeast),
			ooxmlTwipyNaPunkty(wartosc)
	case "exact", "exactly":
		return wejscieWskaznikZasadyInterlinii(shared.StudioLineSpacingRuleExactly),
			ooxmlTwipyNaPunkty(wartosc)
	}
	mnoznik := float64(wartosc) / 240
	switch wartosc {
	case 240:
		return wejscieWskaznikZasadyInterlinii(shared.StudioLineSpacingRuleSingle), 1
	case 360:
		return wejscieWskaznikZasadyInterlinii(shared.StudioLineSpacingRuleOneAndHalf), 1.5
	case 480:
		return wejscieWskaznikZasadyInterlinii(shared.StudioLineSpacingRuleDouble), 2
	}
	return wejscieWskaznikZasadyInterlinii(shared.StudioLineSpacingRuleMultiple), mnoznik
}

// ooxmlBarwa sprowadza barwę OOXML do zapisu szesnastkowego z krzyżykiem.
// `auto` znaczy „barwa wynikowa składu", więc nie jest barwą i nie wychodzi.
func ooxmlBarwa(wartosc string) string {
	oczyszczona := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(wartosc), "#"))
	if oczyszczona == "" || strings.EqualFold(oczyszczona, "auto") {
		return ""
	}
	if len(oczyszczona) != 6 && len(oczyszczona) != 8 {
		return ""
	}
	if _, err := strconv.ParseUint(oczyszczona, 16, 64); err != nil {
		return ""
	}
	return "#" + strings.ToUpper(oczyszczona[:6])
}

// ooxmlBarwaNazwana przekłada barwę wyróżnienia, którą OOXML podaje nazwą,
// na zapis szesnastkowy; nazwa spoza wykazu wraca do ooxmlBarwa.
func ooxmlBarwaNazwana(wartosc string) string {
	switch strings.ToLower(strings.TrimSpace(wartosc)) {
	case "yellow":
		return "#FFFF00"
	case "green", "brightgreen":
		return "#00FF00"
	case "cyan":
		return "#00FFFF"
	case "magenta":
		return "#FF00FF"
	case "blue":
		return "#0000FF"
	case "red":
		return "#FF0000"
	case "darkblue":
		return "#000080"
	case "darkcyan":
		return "#008080"
	case "darkgreen":
		return "#008000"
	case "darkmagenta":
		return "#800080"
	case "darkred":
		return "#800000"
	case "darkyellow":
		return "#808000"
	case "darkgray", "darkgrey":
		return "#808080"
	case "lightgray", "lightgrey":
		return "#C0C0C0"
	case "black":
		return "#000000"
	case "white":
		return "#FFFFFF"
	case "none":
		return ""
	}
	return ooxmlBarwa(wartosc)
}

// ooxmlWskaznikZrodlaWypunktowania oddaje wskaźnik na przekazaną wartość
// źródła znaku listy, bo kontrakt platformy niesie to pole wskaźnikiem.
func ooxmlWskaznikZrodlaWypunktowania(wartosc shared.StudioBulletSource) *shared.StudioBulletSource {
	kopia := wartosc
	return &kopia
}

// ooxmlFormatNumeracjiListy przekłada format numeracji poziomu listy OOXML
// na format kontraktu; wartość spoza wykazu daje format arabski.
func ooxmlFormatNumeracjiListy(format string) *shared.StudioListNumberFormat {
	wynik := shared.StudioListNumberFormat(shared.StudioListNumberFormatArabic)
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "upperroman":
		wynik = shared.StudioListNumberFormatRomanUpper
	case "lowerroman":
		wynik = shared.StudioListNumberFormatRomanLower
	case "upperletter":
		wynik = shared.StudioListNumberFormatLetterUpper
	case "lowerletter":
		wynik = shared.StudioListNumberFormatLetterLower
	}
	return &wynik
}

// ooxmlFormatNumeracjiStron przekłada format numeru strony OOXML na format
// kontraktu; napis pusty oddaje wskaźnik pusty, nie format domyślny.
func ooxmlFormatNumeracjiStron(format string) *shared.StudioPageNumberFormat {
	wynik := shared.StudioPageNumberFormat(shared.StudioPageNumberFormatArabic)
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "":
		return nil
	case "upperroman":
		wynik = shared.StudioPageNumberFormatRomanUpper
	case "lowerroman":
		wynik = shared.StudioPageNumberFormatRomanLower
	case "upperletter":
		wynik = shared.StudioPageNumberFormatLetterUpper
	case "lowerletter":
		wynik = shared.StudioPageNumberFormatLetterLower
	}
	return &wynik
}

// ooxmlZasiegNaglowka przekłada zasięg nagłówka albo stopki OOXML na zasięg
// kontraktu; wartość spoza wykazu daje zasięg domyślny.
func ooxmlZasiegNaglowka(rodzaj string) shared.StudioHeaderScope {
	switch strings.ToLower(strings.TrimSpace(rodzaj)) {
	case "first":
		return shared.StudioHeaderScopeFirstPage
	case "even":
		return shared.StudioHeaderScopeEvenPages
	}
	return shared.StudioHeaderScopeDefault
}

// ooxmlPierwszeOdwolanieObrazu szuka wgłąb rysunku identyfikatora powiązania
// obrazu (`a:blip r:embed`). Rysunek OOXML jest głęboko zagnieżdżony i ścieżka
// bywa różna dla obrazu w linii i obrazu opływanego, więc szuka się wgłąb.
func ooxmlPierwszeOdwolanieObrazu(wezel *ooxmlWezel) string {
	if wezel == nil {
		return ""
	}
	if wezel.Nazwa.Local == "blip" {
		for _, atrybut := range wezel.Atrybuty {
			if atrybut.Name.Local == "embed" || atrybut.Name.Local == "link" {
				return atrybut.Value
			}
		}
	}
	for _, dziecko := range wezel.Dzieci {
		if kod := ooxmlPierwszeOdwolanieObrazu(dziecko); kod != "" {
			return kod
		}
	}
	return ""
}

// ooxmlOpisObrazu składa opis obrazu: tekst zastępczy z pliku, a gdy go nie ma
// — nazwę składnika archiwum. Puste pole byłoby stratą wiedzy, którą plik nosi.
func ooxmlOpisObrazu(wezel *ooxmlWezel, sciezka string) string {
	if opis := ooxmlSzukajWglabAtrybut(wezel, "docPr", "descr"); opis != "" {
		return opis
	}
	if nazwa := ooxmlSzukajWglabAtrybut(wezel, "docPr", "name"); nazwa != "" {
		return nazwa
	}
	return path.Base(sciezka)
}

// ooxmlSzukajWglabAtrybut szuka wgłąb węzła o wskazanej nazwie i oddaje
// wartość jego atrybutu, schodząc rekurencyjnie po całym poddrzewie.
func ooxmlSzukajWglabAtrybut(wezel *ooxmlWezel, nazwaWezla, nazwaAtrybutu string) string {
	if wezel == nil {
		return ""
	}
	if wezel.Nazwa.Local == nazwaWezla {
		if wartosc := strings.TrimSpace(wezel.atrybut(nazwaAtrybutu)); wartosc != "" {
			return wartosc
		}
	}
	for _, dziecko := range wezel.Dzieci {
		if wartosc := ooxmlSzukajWglabAtrybut(dziecko, nazwaWezla, nazwaAtrybutu); wartosc != "" {
			return wartosc
		}
	}
	return ""
}

// ooxmlWymiaryRysunku czyta wymiary obrazu. OOXML podaje je w EMU
// (angielskiej jednostce metrycznej): 360 000 EMU to jeden centymetr.
func ooxmlWymiaryRysunku(wezel *ooxmlWezel) (float64, float64, bool) {
	rozmiar := ooxmlSzukajWglabWezel(wezel, "ext")
	if rozmiar == nil {
		return 0, 0, false
	}
	szerokosc, maSzerokosc := rozmiar.atrybutCalkowity("cx")
	wysokosc, maWysokosc := rozmiar.atrybutCalkowity("cy")
	if !maSzerokosc || !maWysokosc {
		return 0, 0, false
	}
	return float64(szerokosc) / 36000, float64(wysokosc) / 36000, true
}

// ooxmlSzukajWglabWezel szuka wgłąb pierwszego węzła o wskazanej nazwie,
// schodząc rekurencyjnie po całym poddrzewie węzła wejściowego.
func ooxmlSzukajWglabWezel(wezel *ooxmlWezel, nazwa string) *ooxmlWezel {
	if wezel == nil {
		return nil
	}
	if wezel.Nazwa.Local == nazwa {
		return wezel
	}
	for _, dziecko := range wezel.Dzieci {
		if znaleziony := ooxmlSzukajWglabWezel(dziecko, nazwa); znaleziony != nil {
			return znaleziony
		}
	}
	return nil
}

// ooxmlZdanieBilansu składa zdanie o uczciwym stanie odczytu: liczbę
// odzyskanych akapitów, tabel, stylów, sekcji, obrazów i rzeczy pominiętych.
func ooxmlZdanieBilansu(bilans *shared.StudioImportBalance) string {
	czesci := []string{
		"odczytano " + strconv.Itoa(wartoscCalkowita(bilans.ParagraphsRecovered)) + " akapitów",
		strconv.Itoa(wartoscCalkowita(bilans.TablesRecognized)) + " tabel",
		strconv.Itoa(wartoscCalkowita(bilans.StylesRecovered)) + " stylów nazwanych",
		strconv.Itoa(wartoscCalkowita(bilans.SectionsRecovered)) + " sekcji",
		strconv.Itoa(wartoscCalkowita(bilans.ImagesEmbedded)) + " obrazów",
	}
	zdanie := strings.Join(czesci, ", ")
	if len(bilans.Skipped) > 0 {
		zdanie += "; pominięto " + strconv.Itoa(len(bilans.Skipped)) +
			" rzeczy wymienionych w wykazie"
	}
	return zdanie
}

// wartoscCalkowita oddaje wartość wskaźnika liczby całkowitej albo zero,
// gdy wskaźnik jest pusty, bez wyrzucania panik.
func wartoscCalkowita(wskazanie *int) int {
	if wskazanie == nil {
		return 0
	}
	return *wskazanie
}

// ── Złożenie ────────────────────────────────────────────────────────────────

// wejscieZlozOoxml składa `.docx` albo `.dotx` z postaci dokumentu. Archiwum
// niesie pięć składników: wykaz typów treści, powiązania paczki, dokument
// główny, arkusz stylów i powiązania dokumentu — mniej nie wystarcza do
// otwarcia w edytorze.
func wejscieZlozOoxml(postac *shared.StudioDocumentForm, tresc, tytul string,
	szablon bool) ([]byte, []shared.StudioSkippedItem, error) {

	pominiete := make([]shared.StudioSkippedItem, 0, 4)
	bloki := postac.Blocks
	if len(bloki) == 0 {
		// Postaci nie ma — treść płaska staje akapitami, żeby plik nie
		// wyszedł pusty.
		zastepcza := wejsciePostacZTekstu(postac.DocumentId, tresc)
		bloki = zastepcza.Blocks
		if postac.Styles == nil {
			postac.Styles = zastepcza.Styles
		}
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "dokument nie miał zapisanej postaci; wydanie odtworzyło akapity z treści płaskiej",
		})
	}

	var cialo strings.Builder
	for _, blok := range bloki {
		switch blok.Kind {
		case wejscieRodzajBlokuTabela:
			cialo.WriteString(ooxmlZlozTabele(postac, blok.TableId, &pominiete))
		case wejscieRodzajBlokuObiekt:
			// Obraz osadzony wychodzi akapitem z jego tekstem zastępczym,
			// bez bajtów obrazu z magazynu zasobów.
			cialo.WriteString(ooxmlZlozAkapitObiektu(postac, blok.ObjectId, &pominiete))
		case wejscieRodzajBlokuPodzial:
			cialo.WriteString(`<w:p><w:r><w:br w:type="page"/></w:r></w:p>`)
		default:
			cialo.WriteString(ooxmlZlozAkapit(blok))
		}
	}
	if len(postac.Apparatus) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "aparat dokumentu wyszedł treścią, bez pól odświeżalnych Worda",
			Detail: wejscieWskaznikTekstu("spis treści, przypisy, podpisy i odsyłacze — " +
				strconv.Itoa(len(postac.Apparatus)) + " elementów"),
		})
	}
	if len(postac.Fields) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "pola dokumentu wyszły wartością policzoną, nie polem odświeżalnym",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Fields)) + " pól"),
		})
	}
	// Nagłówek i stopka są OSOBNYMI składnikami archiwum, na które sekcja
	// wskazuje odwołaniem.
	naglowki := ooxmlZlozNaglowkiStopki(postac)
	cialo.WriteString(ooxmlZlozNastawySekcji(postac, naglowki.odwolania))

	dokument := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" +
		`<w:document xmlns:w="` + ooxmlPrzestrzenGlowna + `" xmlns:r="` +
		ooxmlPrzestrzenPowiazania + `"><w:body>` + cialo.String() + `</w:body></w:document>`

	typDokumentu := "application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"
	if szablon {
		typDokumentu = "application/vnd.openxmlformats-officedocument.wordprocessingml.template.main+xml"
	}

	skladniki := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Override PartName="/word/document.xml" ContentType="` + typDokumentu + `"/>` +
			`<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>` +
			naglowki.typyTresci +
			`<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>` +
			`</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rIdDokument" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>` +
			`<Relationship Id="rIdWlasciwosci" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>` +
			`</Relationships>`,
		"word/_rels/document.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rIdStyle" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>` +
			naglowki.powiazania +
			`</Relationships>`,
		ooxmlSkladnikDokumentu: dokument,
		ooxmlSkladnikStylow:    ooxmlZlozArkuszStylow(postac.Styles),
		"docProps/core.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" +
			`<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties"` +
			` xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:title>` +
			ooxmlZabezpiecz(tytul) + `</dc:title></cp:coreProperties>`,
	}

	for nazwa, tresc := range naglowki.skladniki {
		skladniki[nazwa] = tresc
	}

	bajty, err := wejscieZlozArchiwum(skladniki, "")
	if err != nil {
		return nil, pominiete, err
	}
	return bajty, pominiete, nil
}

// wejscieZlozArchiwum zapisuje składniki jako archiwum ZIP. Składnik
// `mimetype` — gdy wskazany — idzie PIERWSZY i BEZ KOMPRESJI, bo tak
// OpenDocument rozpoznaje swoje pliki; OOXML tego składnika nie ma, i wtedy
// wskazanie jest puste.
func wejscieZlozArchiwum(skladniki map[string]string, pierwszy string) ([]byte, error) {
	var bufor bytes.Buffer
	archiwum := zip.NewWriter(&bufor)

	if pierwszy != "" {
		if tresc, jest := skladniki[pierwszy]; jest {
			wpis, err := archiwum.CreateHeader(&zip.FileHeader{
				Name: pierwszy, Method: zip.Store,
			})
			if err != nil {
				return nil, bladStudio(err)
			}
			if _, err := wpis.Write([]byte(tresc)); err != nil {
				return nil, bladStudio(err)
			}
		}
	}

	// Kolejność pozostałych składników jest ustalona sortowaniem, nie
	// przebiegiem po mapie.
	nazwy := make([]string, 0, len(skladniki))
	for nazwa := range skladniki {
		if nazwa == pierwszy {
			continue
		}
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	for _, nazwa := range nazwy {
		wpis, err := archiwum.Create(nazwa)
		if err != nil {
			return nil, bladStudio(err)
		}
		if _, err := wpis.Write([]byte(skladniki[nazwa])); err != nil {
			return nil, bladStudio(err)
		}
	}
	if err := archiwum.Close(); err != nil {
		return nil, bladStudio(err)
	}
	return bufor.Bytes(), nil
}

// ooxmlZlozAkapit składa akapit wraz z postacią akapitu i fragmentami,
// zwracając gotowy węzeł `w:p` w postaci tekstu XML.
func ooxmlZlozAkapit(blok shared.StudioDocumentBlock) string {
	var budowa strings.Builder
	budowa.WriteString("<w:p>")
	budowa.WriteString(ooxmlZlozPostacAkapitu(blok.Paragraph))
	for _, fragment := range blok.Runs {
		budowa.WriteString(ooxmlZlozFragment(fragment))
	}
	budowa.WriteString("</w:p>")
	return budowa.String()
}

// ooxmlZlozFragment składa fragment tekstu. Znak końca wiersza wewnątrz
// fragmentu wychodzi węzłem `w:br`, a tabulator węzłem `w:tab` — wpisane
// wprost, byłyby w OOXML pojedynczą spacją.
func ooxmlZlozFragment(fragment shared.StudioDocumentRun) string {
	var budowa strings.Builder
	budowa.WriteString("<w:r>")
	budowa.WriteString(ooxmlZlozPostacZnaku(fragment.Format))
	for i, wiersz := range strings.Split(fragment.Text, "\n") {
		if i > 0 {
			budowa.WriteString("<w:br/>")
		}
		for j, czesc := range strings.Split(wiersz, "\t") {
			if j > 0 {
				budowa.WriteString("<w:tab/>")
			}
			if czesc == "" {
				continue
			}
			budowa.WriteString(`<w:t xml:space="preserve">` + ooxmlZabezpiecz(czesc) + `</w:t>`)
		}
	}
	budowa.WriteString("</w:r>")
	return budowa.String()
}

// ooxmlZlozPostacZnaku składa węzeł `w:rPr` z postaci znaku kontraktu: styl,
// krój, przełączniki, indeks, stopień, barwę i język.
func ooxmlZlozPostacZnaku(postac *shared.StudioCharacterFormat) string {
	if postac == nil {
		return ""
	}
	var budowa strings.Builder
	budowa.WriteString("<w:rPr>")
	if postac.StyleName != nil && *postac.StyleName != "" {
		budowa.WriteString(`<w:rStyle w:val="` + ooxmlZabezpiecz(*postac.StyleName) + `"/>`)
	}
	if postac.FontFamily != nil && *postac.FontFamily != "" {
		nazwa := ooxmlZabezpiecz(*postac.FontFamily)
		budowa.WriteString(`<w:rFonts w:ascii="` + nazwa + `" w:hAnsi="` + nazwa + `"/>`)
	}
	if postac.Bold != nil {
		budowa.WriteString(ooxmlPrzelacznik("w:b", *postac.Bold))
	}
	if postac.Italic != nil {
		budowa.WriteString(ooxmlPrzelacznik("w:i", *postac.Italic))
	}
	if postac.Strikethrough != nil {
		budowa.WriteString(ooxmlPrzelacznik("w:strike", *postac.Strikethrough))
	}
	if postac.SmallCaps != nil {
		budowa.WriteString(ooxmlPrzelacznik("w:smallCaps", *postac.SmallCaps))
	}
	if postac.AllCaps != nil {
		budowa.WriteString(ooxmlPrzelacznik("w:caps", *postac.AllCaps))
	}
	if postac.Underline != nil {
		budowa.WriteString(`<w:u w:val="` + ooxmlNazwaPodkreslenia(*postac.Underline) + `"/>`)
	}
	switch {
	case postac.Superscript != nil && *postac.Superscript:
		budowa.WriteString(`<w:vertAlign w:val="superscript"/>`)
	case postac.Subscript != nil && *postac.Subscript:
		budowa.WriteString(`<w:vertAlign w:val="subscript"/>`)
	}
	if postac.FontSizePt != nil && *postac.FontSizePt > 0 {
		polpunkty := strconv.Itoa(int(*postac.FontSizePt * 2))
		budowa.WriteString(`<w:sz w:val="` + polpunkty + `"/><w:szCs w:val="` + polpunkty + `"/>`)
	}
	if postac.Color != nil {
		if barwa := ooxmlBezKrzyzyka(*postac.Color); barwa != "" {
			budowa.WriteString(`<w:color w:val="` + barwa + `"/>`)
		}
	}
	if postac.HighlightColor != nil {
		if barwa := ooxmlBezKrzyzyka(*postac.HighlightColor); barwa != "" {
			// Wyróżnienie wychodzi cieniowaniem znaku, nie `w:highlight`,
			// bo paleta wyróżnień jest dowolna.
			budowa.WriteString(`<w:shd w:val="clear" w:color="auto" w:fill="` + barwa + `"/>`)
		}
	}
	if postac.LetterSpacingPt != nil && *postac.LetterSpacingPt != 0 {
		budowa.WriteString(`<w:spacing w:val="` +
			strconv.Itoa(int(*postac.LetterSpacingPt*20)) + `"/>`)
	}
	if postac.Language != nil && *postac.Language != "" {
		budowa.WriteString(`<w:lang w:val="` + ooxmlZabezpiecz(*postac.Language) + `"/>`)
	}
	budowa.WriteString("</w:rPr>")
	return budowa.String()
}

// ooxmlZlozPostacAkapitu składa węzeł `w:pPr` z postaci akapitu kontraktu:
// styl, listę, odstępy, wcięcia, tabulatory i obramowanie.
func ooxmlZlozPostacAkapitu(postac *shared.StudioParagraphFormat) string {
	if postac == nil {
		return ""
	}
	var budowa strings.Builder
	budowa.WriteString("<w:pPr>")
	if postac.StyleName != nil && *postac.StyleName != "" {
		budowa.WriteString(`<w:pStyle w:val="` + ooxmlZabezpiecz(*postac.StyleName) + `"/>`)
	}
	if postac.ListId != nil && *postac.ListId != "" {
		poziom := 0
		if postac.ListLevel != nil && *postac.ListLevel > 0 {
			poziom = *postac.ListLevel - 1
		}
		budowa.WriteString(`<w:numPr><w:ilvl w:val="` + strconv.Itoa(poziom) +
			`"/><w:numId w:val="` +
			strconv.Itoa(ooxmlNumerListy(*postac.ListId)) + `"/></w:numPr>`)
	}
	if postac.KeepWithNext != nil {
		budowa.WriteString(ooxmlPrzelacznik("w:keepNext", *postac.KeepWithNext))
	}
	if postac.KeepLines != nil {
		budowa.WriteString(ooxmlPrzelacznik("w:keepLines", *postac.KeepLines))
	}
	if postac.WidowControl != nil {
		budowa.WriteString(ooxmlPrzelacznik("w:widowControl", *postac.WidowControl))
	}
	if postac.Border != nil {
		budowa.WriteString(ooxmlZlozObramowanie("w:pBdr", postac.Border))
	}
	if postac.ShadingColor != nil {
		if barwa := ooxmlBezKrzyzyka(*postac.ShadingColor); barwa != "" {
			budowa.WriteString(`<w:shd w:val="clear" w:color="auto" w:fill="` + barwa + `"/>`)
		}
	}
	if len(postac.TabStops) > 0 {
		budowa.WriteString("<w:tabs>")
		for _, tabulator := range postac.TabStops {
			budowa.WriteString(`<w:tab w:val="` + ooxmlNazwaTabulatora(tabulator.Kind) +
				`" w:pos="` + strconv.Itoa(ooxmlMilimetryNaTwipy(tabulator.PositionMm)) + `"`)
			if tabulator.Leader != nil {
				if nazwa := ooxmlNazwaZnakuWiodacego(*tabulator.Leader); nazwa != "" {
					budowa.WriteString(` w:leader="` + nazwa + `"`)
				}
			}
			budowa.WriteString("/>")
		}
		budowa.WriteString("</w:tabs>")
	}
	if odstepy := ooxmlZlozOdstepy(postac); odstepy != "" {
		budowa.WriteString(odstepy)
	}
	if wciecia := ooxmlZlozWciecia(postac); wciecia != "" {
		budowa.WriteString(wciecia)
	}
	if postac.Align != nil {
		budowa.WriteString(`<w:jc w:val="` + ooxmlNazwaWyrownania(*postac.Align) + `"/>`)
	}
	if postac.OutlineLevel != nil && *postac.OutlineLevel > 0 {
		budowa.WriteString(`<w:outlineLvl w:val="` +
			strconv.Itoa(*postac.OutlineLevel-1) + `"/>`)
	}
	if postac.RightToLeft != nil && *postac.RightToLeft {
		budowa.WriteString("<w:bidi/>")
	}
	budowa.WriteString("</w:pPr>")
	return budowa.String()
}

// ooxmlZlozOdstepy składa węzeł `w:spacing` akapitu z odstępów przed i po
// oraz z zasady i wartości interlinii kontraktu.
func ooxmlZlozOdstepy(postac *shared.StudioParagraphFormat) string {
	czesci := make([]string, 0, 3)
	if postac.SpaceBeforePt != nil {
		czesci = append(czesci, `w:before="`+strconv.Itoa(int(*postac.SpaceBeforePt*20))+`"`)
	}
	if postac.SpaceAfterPt != nil {
		czesci = append(czesci, `w:after="`+strconv.Itoa(int(*postac.SpaceAfterPt*20))+`"`)
	}
	if postac.LineSpacingRule != nil {
		wartosc := 1.0
		if postac.LineSpacingValue != nil {
			wartosc = *postac.LineSpacingValue
		}
		switch *postac.LineSpacingRule {
		case shared.StudioLineSpacingRuleAtLeast:
			czesci = append(czesci, `w:line="`+strconv.Itoa(int(wartosc*20))+`" w:lineRule="atLeast"`)
		case shared.StudioLineSpacingRuleExactly:
			czesci = append(czesci, `w:line="`+strconv.Itoa(int(wartosc*20))+`" w:lineRule="exact"`)
		case shared.StudioLineSpacingRuleSingle:
			czesci = append(czesci, `w:line="240" w:lineRule="auto"`)
		case shared.StudioLineSpacingRuleOneAndHalf:
			czesci = append(czesci, `w:line="360" w:lineRule="auto"`)
		case shared.StudioLineSpacingRuleDouble:
			czesci = append(czesci, `w:line="480" w:lineRule="auto"`)
		default:
			czesci = append(czesci, `w:line="`+strconv.Itoa(int(wartosc*240))+`" w:lineRule="auto"`)
		}
	}
	if len(czesci) == 0 {
		return ""
	}
	return "<w:spacing " + strings.Join(czesci, " ") + "/>"
}

// ooxmlZlozWciecia składa `w:ind` akapitu. Wcięcie pierwszego wiersza o wartości
// ujemnej wychodzi jako wysunięcie (`w:hanging`) — OOXML nie zna wcięcia
// ujemnego i wartość ujemna w `w:firstLine` zostałaby zignorowana.
func ooxmlZlozWciecia(postac *shared.StudioParagraphFormat) string {
	czesci := make([]string, 0, 3)
	if postac.IndentLeftMm != nil {
		czesci = append(czesci, `w:left="`+strconv.Itoa(ooxmlMilimetryNaTwipy(*postac.IndentLeftMm))+`"`)
	}
	if postac.IndentRightMm != nil {
		czesci = append(czesci, `w:right="`+strconv.Itoa(ooxmlMilimetryNaTwipy(*postac.IndentRightMm))+`"`)
	}
	if postac.FirstLineIndentMm != nil {
		twipy := ooxmlMilimetryNaTwipy(*postac.FirstLineIndentMm)
		if twipy < 0 {
			czesci = append(czesci, `w:hanging="`+strconv.Itoa(-twipy)+`"`)
		} else if twipy > 0 {
			czesci = append(czesci, `w:firstLine="`+strconv.Itoa(twipy)+`"`)
		}
	}
	if len(czesci) == 0 {
		return ""
	}
	return "<w:ind " + strings.Join(czesci, " ") + "/>"
}

// ooxmlZlozObramowanie składa obramowanie akapitu, tabeli albo komórki pod
// wskazaną nazwą węzła, wspólną dla tych trzech miejsc OOXML.
func ooxmlZlozObramowanie(nazwaWezla string, obramowanie *shared.StudioBorder) string {
	if obramowanie == nil || obramowanie.Style == shared.StudioBorderStyleNone {
		return ""
	}
	odmiana := ooxmlNazwaOdmianyObramowania(obramowanie.Style)
	grubosc := 4
	if obramowanie.WidthPt != nil && *obramowanie.WidthPt > 0 {
		grubosc = int(*obramowanie.WidthPt * 8)
	}
	barwa := "auto"
	if obramowanie.Color != nil {
		if zapis := ooxmlBezKrzyzyka(*obramowanie.Color); zapis != "" {
			barwa = zapis
		}
	}
	// Krawędź bez jawnego wskazania jest krawędzią OBECNĄ, nie brakującą.
	krawedzie := []struct {
		nazwa       string
		przelacznik *bool
	}{
		{"top", obramowanie.Top}, {"left", obramowanie.Left},
		{"bottom", obramowanie.Bottom}, {"right", obramowanie.Right},
	}
	var budowa strings.Builder
	budowa.WriteString("<" + nazwaWezla + ">")
	for _, krawedz := range krawedzie {
		if krawedz.przelacznik != nil && !*krawedz.przelacznik {
			budowa.WriteString(`<w:` + krawedz.nazwa + ` w:val="nil"/>`)
			continue
		}
		budowa.WriteString(`<w:` + krawedz.nazwa + ` w:val="` + odmiana +
			`" w:sz="` + strconv.Itoa(grubosc) + `" w:space="0" w:color="` + barwa + `"/>`)
	}
	budowa.WriteString("</" + nazwaWezla + ">")
	return budowa.String()
}

// ooxmlZlozTabele składa tabelę wraz z siatką kolumn, scaleniami poziomymi
// i pionowymi oraz postacią i treścią każdej komórki.
func ooxmlZlozTabele(postac *shared.StudioDocumentForm, kodTabeli *string,
	pominiete *[]shared.StudioSkippedItem) string {

	if kodTabeli == nil {
		return ""
	}
	var tabela *shared.StudioDocumentTable
	for i := range postac.Tables {
		if postac.Tables[i].Id == *kodTabeli {
			tabela = &postac.Tables[i]
			break
		}
	}
	if tabela == nil {
		*pominiete = append(*pominiete, shared.StudioSkippedItem{
			Reason: "blok wskazuje tabelę, której postać dokumentu nie niesie",
			Detail: wejscieWskaznikTekstu(*kodTabeli),
		})
		return ""
	}

	var budowa strings.Builder
	budowa.WriteString("<w:tbl><w:tblPr>")
	if tabela.StyleName != nil && *tabela.StyleName != "" {
		budowa.WriteString(`<w:tblStyle w:val="` + ooxmlZabezpiecz(*tabela.StyleName) + `"/>`)
	}
	if tabela.WidthMm != nil && *tabela.WidthMm > 0 {
		budowa.WriteString(`<w:tblW w:w="` +
			strconv.Itoa(ooxmlMilimetryNaTwipy(*tabela.WidthMm)) + `" w:type="dxa"/>`)
	} else {
		budowa.WriteString(`<w:tblW w:w="0" w:type="auto"/>`)
	}
	budowa.WriteString(ooxmlZlozObramowanie("w:tblBorders", tabela.Border))
	budowa.WriteString("</w:tblPr>")

	budowa.WriteString("<w:tblGrid>")
	szerokosci := tabela.ColumnWidthsMm
	if len(szerokosci) != tabela.Columns && tabela.Columns > 0 {
		rowna := ooxmlSzerokoscObszaruPisania(postac) / float64(tabela.Columns)
		szerokosci = make([]float64, tabela.Columns)
		for i := range szerokosci {
			szerokosci[i] = rowna
		}
	}
	for _, szerokosc := range szerokosci {
		budowa.WriteString(`<w:gridCol w:w="` +
			strconv.Itoa(ooxmlMilimetryNaTwipy(szerokosc)) + `"/>`)
	}
	budowa.WriteString("</w:tblGrid>")

	liczbaNaglowkow := 0
	if tabela.HeaderRows != nil {
		liczbaNaglowkow = *tabela.HeaderRows
	}
	powtarzaj := tabela.RepeatHeader != nil && *tabela.RepeatHeader
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		budowa.WriteString("<w:tr>")
		if powtarzaj && wiersz < liczbaNaglowkow {
			budowa.WriteString("<w:trPr><w:tblHeader/></w:trPr>")
		}
		kolumna := 0
		for kolumna < tabela.Columns {
			komorka := ooxmlSzukajKomorki(tabela, wiersz, kolumna)
			if komorka == nil {
				budowa.WriteString("<w:tc><w:tcPr/><w:p/></w:tc>")
				kolumna++
				continue
			}
			if komorka.Merged != nil && *komorka.Merged {
				// Komórka wchłonięta scaleniem poziomym nie wychodzi wcale;
				// scalenie pionowe wychodzi jako `w:vMerge`.
				if ooxmlScalonaPoziomo(tabela, wiersz, kolumna) {
					kolumna++
					continue
				}
				budowa.WriteString("<w:tc><w:tcPr><w:vMerge/></w:tcPr><w:p/></w:tc>")
				kolumna++
				continue
			}
			rozpietosc := 1
			if komorka.ColumnSpan != nil && *komorka.ColumnSpan > 1 {
				rozpietosc = *komorka.ColumnSpan
			}
			budowa.WriteString(ooxmlZlozKomorke(komorka, rozpietosc))
			kolumna += rozpietosc
		}
		budowa.WriteString("</w:tr>")
	}
	budowa.WriteString("</w:tbl>")
	// Word łączy dwie tabele stojące bezpośrednio po sobie; akapit pusty
	// między nimi jest koniecznością.
	budowa.WriteString("<w:p/>")
	return budowa.String()
}

// ooxmlScalonaPoziomo rozstrzyga, czy komórkę wchłonęło scalenie poziome —
// czyli czy w tym samym wierszu, przed nią, stoi komórka o rozpiętości
// obejmującej jej kolumnę.
func ooxmlScalonaPoziomo(tabela *shared.StudioDocumentTable, wiersz, kolumna int) bool {
	for i := range tabela.Cells {
		komorka := &tabela.Cells[i]
		if komorka.Row != wiersz || komorka.Column >= kolumna {
			continue
		}
		if komorka.ColumnSpan == nil || *komorka.ColumnSpan <= 1 {
			continue
		}
		if komorka.Column+*komorka.ColumnSpan > kolumna {
			return true
		}
	}
	return false
}

// ooxmlSzukajKomorki oddaje komórkę tabeli o wskazanym wierszu i kolumnie
// albo nic, gdy takiej komórki tabela nie niesie.
func ooxmlSzukajKomorki(tabela *shared.StudioDocumentTable, wiersz, kolumna int) *shared.StudioTableCell {
	for i := range tabela.Cells {
		if tabela.Cells[i].Row == wiersz && tabela.Cells[i].Column == kolumna {
			return &tabela.Cells[i]
		}
	}
	return nil
}

// ooxmlZlozKomorke składa komórkę tabeli wraz ze scaleniem, cieniowaniem,
// obramowaniem, wyrównaniem pionowym i treścią akapitów.
func ooxmlZlozKomorke(komorka *shared.StudioTableCell, rozpietosc int) string {
	var budowa strings.Builder
	budowa.WriteString("<w:tc><w:tcPr>")
	if rozpietosc > 1 {
		budowa.WriteString(`<w:gridSpan w:val="` + strconv.Itoa(rozpietosc) + `"/>`)
	}
	if komorka.RowSpan != nil && *komorka.RowSpan > 1 {
		budowa.WriteString(`<w:vMerge w:val="restart"/>`)
	}
	if komorka.ShadingColor != nil {
		if barwa := ooxmlBezKrzyzyka(*komorka.ShadingColor); barwa != "" {
			budowa.WriteString(`<w:shd w:val="clear" w:color="auto" w:fill="` + barwa + `"/>`)
		}
	}
	budowa.WriteString(ooxmlZlozObramowanie("w:tcBorders", komorka.Border))
	if komorka.VerticalAlign != nil {
		nazwa := "top"
		switch *komorka.VerticalAlign {
		case shared.StudioVerticalAlignMiddle:
			nazwa = "center"
		case shared.StudioVerticalAlignBottom:
			nazwa = "bottom"
		}
		budowa.WriteString(`<w:vAlign w:val="` + nazwa + `"/>`)
	}
	budowa.WriteString("</w:tcPr>")

	tresc := ""
	if komorka.Text != nil {
		tresc = *komorka.Text
	}
	wiersze := strings.Split(tresc, "\n")
	for _, wiersz := range wiersze {
		budowa.WriteString(ooxmlZlozAkapit(shared.StudioDocumentBlock{
			Paragraph: komorka.Paragraph,
			Runs:      []shared.StudioDocumentRun{{Text: wiersz, Format: komorka.Character}},
		}))
	}
	budowa.WriteString("</w:tc>")
	return budowa.String()
}

// ooxmlZlozAkapitObiektu składa akapit zastępujący obiekt osadzony jego
// opisem, bo bajtów obiektu składacz do archiwum nie wpisuje.
func ooxmlZlozAkapitObiektu(postac *shared.StudioDocumentForm, kodObiektu *string,
	pominiete *[]shared.StudioSkippedItem) string {

	if kodObiektu == nil {
		return ""
	}
	for i := range postac.Objects {
		obiekt := &postac.Objects[i]
		if obiekt.Id != *kodObiektu {
			continue
		}
		opis := "obiekt osadzony"
		if obiekt.AltText != nil && *obiekt.AltText != "" {
			opis = *obiekt.AltText
		} else if obiekt.Caption != nil && *obiekt.Caption != "" {
			opis = *obiekt.Caption
		}
		*pominiete = append(*pominiete, shared.StudioSkippedItem{
			Reason: "obiekt osadzony wyszedł opisem, bez bajtów obrazu",
			Detail: wejscieWskaznikTekstu(string(obiekt.Kind) + ": " + opis),
		})
		return ooxmlZlozAkapit(shared.StudioDocumentBlock{
			Paragraph: &shared.StudioParagraphFormat{
				StyleName: wejscieWskaznikTekstu(wejscieStylPodpis),
			},
			Runs: []shared.StudioDocumentRun{{Text: "[" + opis + "]"}},
		})
	}
	return ""
}

// ooxmlZlozNastawySekcji składa węzeł `w:sectPr` z nastaw strony dokumentu:
// rozmiar, orientację, marginesy i kolumny szpaltowe.
func ooxmlZlozNastawySekcji(postac *shared.StudioDocumentForm, odwolania string) string {
	nastawy := postac.PageSetup
	if nastawy == nil && len(postac.Sections) > 0 {
		nastawy = postac.Sections[len(postac.Sections)-1].PageSetup
	}
	if nastawy == nil {
		domyslne := wejscieDomyslneNastawyStrony("", nil)
		nastawy = &domyslne
	}

	szerokosc, wysokosc := 210.0, 297.0
	if nastawy.PageSize != nil {
		if wymiar, jest := wejscieWymiaryNosnika(*nastawy.PageSize); jest {
			szerokosc, wysokosc = wymiar.szerokosc, wymiar.wysokosc
		}
	}
	if nastawy.WidthMm != nil && *nastawy.WidthMm > 0 {
		szerokosc = *nastawy.WidthMm
	}
	if nastawy.HeightMm != nil && *nastawy.HeightMm > 0 {
		wysokosc = *nastawy.HeightMm
	}
	orientacja := "portrait"
	if nastawy.Orientation != nil && *nastawy.Orientation == shared.StudioPageOrientationPozioma {
		orientacja = "landscape"
		// Orientacja pozioma wymaga zamiany wymiarów w `w:pgSz`, nie tylko
		// atrybutu orientacji.
		if szerokosc < wysokosc {
			szerokosc, wysokosc = wysokosc, szerokosc
		}
	}

	var budowa strings.Builder
	budowa.WriteString(`<w:sectPr>` + odwolania + `<w:pgSz w:w="` +
		strconv.Itoa(ooxmlMilimetryNaTwipy(szerokosc)) + `" w:h="` +
		strconv.Itoa(ooxmlMilimetryNaTwipy(wysokosc)) + `" w:orient="` + orientacja + `"/>`)
	budowa.WriteString(`<w:pgMar w:top="` + ooxmlMarginesTwipy(nastawy.MarginTop) +
		`" w:right="` + ooxmlMarginesTwipy(nastawy.MarginRight) +
		`" w:bottom="` + ooxmlMarginesTwipy(nastawy.MarginBottom) +
		`" w:left="` + ooxmlMarginesTwipy(nastawy.MarginLeft) + `" w:gutter="` +
		strconv.Itoa(ooxmlMilimetryNaTwipy(wartoscRzeczywista(nastawy.GutterMm))) + `"/>`)
	if nastawy.Columns != nil && *nastawy.Columns > 1 {
		budowa.WriteString(`<w:cols w:num="` + strconv.Itoa(*nastawy.Columns) + `"`)
		if nastawy.ColumnGapMm != nil {
			budowa.WriteString(` w:space="` +
				strconv.Itoa(ooxmlMilimetryNaTwipy(*nastawy.ColumnGapMm)) + `"`)
		}
		budowa.WriteString("/>")
	}
	if nastawy.MirrorMargins != nil && *nastawy.MirrorMargins {
		budowa.WriteString("<w:mirrorMargins/>")
	}
	budowa.WriteString("</w:sectPr>")
	return budowa.String()
}

// ooxmlMarginesTwipy przekłada margines w milimetrach na twipy; brak marginesu
// daje 25 mm, czyli nastawę domyślną pisma.
func ooxmlMarginesTwipy(milimetry *int) string {
	wartosc := 25
	if milimetry != nil {
		wartosc = *milimetry
	}
	return strconv.Itoa(ooxmlMilimetryNaTwipy(float64(wartosc)))
}

// wartoscRzeczywista oddaje wartość wskaźnika liczby rzeczywistej albo zero,
// gdy wskaźnik jest pusty, bez wyrzucania panik.
func wartoscRzeczywista(wskazanie *float64) float64 {
	if wskazanie == nil {
		return 0
	}
	return *wskazanie
}

// ooxmlZlozArkuszStylow składa składnik `word/styles.xml` z wykazu stylów
// nazwanych kontraktu, a przy braku stylów — z arkusza domyślnego.
func ooxmlZlozArkuszStylow(arkusz []shared.StudioNamedStyle) string {
	if len(arkusz) == 0 {
		arkusz = wejscieDomyslnyArkuszStylow()
	}
	var budowa strings.Builder
	budowa.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" +
		`<w:styles xmlns:w="` + ooxmlPrzestrzenGlowna + `">`)
	for _, styl := range arkusz {
		rodzaj := "paragraph"
		switch styl.Kind {
		case shared.StudioStyleKindCharacter:
			rodzaj = "character"
		case shared.StudioStyleKindTable:
			rodzaj = "table"
		case shared.StudioStyleKindList:
			rodzaj = "numbering"
		}
		budowa.WriteString(`<w:style w:type="` + rodzaj + `" w:styleId="` +
			ooxmlZabezpiecz(styl.Name) + `">`)
		nazwaWidoczna := styl.Name
		if styl.DisplayName != nil && *styl.DisplayName != "" {
			nazwaWidoczna = *styl.DisplayName
		}
		budowa.WriteString(`<w:name w:val="` + ooxmlZabezpiecz(nazwaWidoczna) + `"/>`)
		if styl.BasedOn != nil && *styl.BasedOn != "" {
			budowa.WriteString(`<w:basedOn w:val="` + ooxmlZabezpiecz(*styl.BasedOn) + `"/>`)
		}
		if styl.NextStyle != nil && *styl.NextStyle != "" {
			budowa.WriteString(`<w:next w:val="` + ooxmlZabezpiecz(*styl.NextStyle) + `"/>`)
		}
		budowa.WriteString(ooxmlZlozPostacAkapitu(styl.Paragraph))
		budowa.WriteString(ooxmlZlozPostacZnaku(styl.Character))
		budowa.WriteString("</w:style>")
	}
	budowa.WriteString("</w:styles>")
	return budowa.String()
}

// ── Przekłady przy zapisie ─────────────────────────────────────────────────

// ooxmlPrzelacznik składa węzeł przełącznika OOXML. Wyłączenie wychodzi jawnie
// (`w:val="0"`), bo sam brak węzła znaczyłby „bez zmiany wobec stylu", a nie
// „wyłączone" — i pogrubienie odziedziczone po stylu zostałoby na miejscu.
func ooxmlPrzelacznik(nazwa string, wlaczony bool) string {
	if wlaczony {
		return "<" + nazwa + "/>"
	}
	return "<" + nazwa + ` w:val="0"/>`
}

// ooxmlNazwaWyrownania przekłada wyrównanie kontraktu na nazwę OOXML;
// wartość spoza wykazu daje wyrównanie do lewej.
func ooxmlNazwaWyrownania(wyrownanie shared.StudioTextAlign) string {
	switch wyrownanie {
	case shared.StudioTextAlignCenter:
		return "center"
	case shared.StudioTextAlignRight:
		return "right"
	case shared.StudioTextAlignJustify:
		return "both"
	}
	return "left"
}

// ooxmlNazwaPodkreslenia przekłada odmianę podkreślenia kontraktu na nazwę
// OOXML; wartość spoza wykazu daje podkreślenie pojedyncze.
func ooxmlNazwaPodkreslenia(odmiana shared.StudioUnderlineStyle) string {
	switch odmiana {
	case shared.StudioUnderlineStyleNone:
		return "none"
	case shared.StudioUnderlineStyleDouble:
		return "double"
	case shared.StudioUnderlineStyleThick:
		return "thick"
	case shared.StudioUnderlineStyleDotted:
		return "dotted"
	case shared.StudioUnderlineStyleDashed:
		return "dash"
	case shared.StudioUnderlineStyleWavy:
		return "wave"
	case shared.StudioUnderlineStyleWords:
		return "words"
	}
	return "single"
}

// ooxmlNazwaOdmianyObramowania przekłada odmianę obramowania kontraktu na
// nazwę OOXML; wartość spoza wykazu daje obramowanie pojedyncze.
func ooxmlNazwaOdmianyObramowania(odmiana shared.StudioBorderStyle) string {
	switch odmiana {
	case shared.StudioBorderStyleDouble:
		return "double"
	case shared.StudioBorderStyleThick:
		return "thick"
	case shared.StudioBorderStyleDotted:
		return "dotted"
	case shared.StudioBorderStyleDashed:
		return "dashed"
	}
	return "single"
}

// ooxmlNazwaTabulatora przekłada rodzaj tabulatora kontraktu na nazwę OOXML;
// wartość spoza wykazu daje tabulator lewy.
func ooxmlNazwaTabulatora(rodzaj shared.StudioTabKind) string {
	switch rodzaj {
	case shared.StudioTabKindRight:
		return "right"
	case shared.StudioTabKindCenter:
		return "center"
	case shared.StudioTabKindDecimal:
		return "decimal"
	case shared.StudioTabKindBar:
		return "bar"
	}
	return "left"
}

// ooxmlNazwaZnakuWiodacego przekłada znak wiodący kontraktu na nazwę OOXML;
// wartość spoza wykazu daje napis pusty, czyli brak znaku.
func ooxmlNazwaZnakuWiodacego(znak shared.StudioTabLeader) string {
	switch znak {
	case shared.StudioTabLeaderDot:
		return "dot"
	case shared.StudioTabLeaderDash:
		return "hyphen"
	case shared.StudioTabLeaderUnderline:
		return "underscore"
	}
	return ""
}

// ooxmlNumerListy odzyskuje numer OOXML z nazwy definicji listy. Nazwa nie
// pochodząca z odczytu OOXML dostaje numer jeden — lista musi mieć numer,
// a jedynka jest numerem pierwszej listy dokumentu.
func ooxmlNumerListy(nazwa string) int {
	if numer, err := strconv.Atoi(strings.TrimPrefix(nazwa, "lista-")); err == nil && numer > 0 {
		return numer
	}
	return 1
}

// ooxmlBezKrzyzyka sprowadza barwę kontraktu do zapisu OOXML bez krzyżyka;
// barwa nieprawidłowa daje napis pusty.
func ooxmlBezKrzyzyka(barwa string) string {
	oczyszczona := strings.TrimPrefix(strings.TrimSpace(barwa), "#")
	if len(oczyszczona) < 6 {
		return ""
	}
	if _, err := strconv.ParseUint(oczyszczona[:6], 16, 64); err != nil {
		return ""
	}
	return strings.ToUpper(oczyszczona[:6])
}

// ooxmlZabezpiecz zamienia znaki o znaczeniu składniowym XML na encje, żeby
// tekst dowolnej treści trafił do dokumentu bezpiecznie.
func ooxmlZabezpiecz(tekst string) string {
	var bufor bytes.Buffer
	if err := xml.EscapeText(&bufor, []byte(tekst)); err != nil {
		return ""
	}
	return bufor.String()
}

// ── Nagłówek i stopka ───────────────────────────────────────────────────────

// ooxmlNaglowkiArchiwum zbiera to, czego nagłówek i stopka wymagają w
// czterech miejscach archiwum naraz: składniki z ich treścią, powiązania
// dokumentu, wykaz typów treści i odwołania w nastawach sekcji.
type ooxmlNaglowkiArchiwum struct {
	skladniki  map[string]string
	powiazania string
	typyTresci string
	odwolania  string
}

// ooxmlZlozNaglowkiStopki składa nagłówek i stopkę dokumentu. Zasięgi idą
// osobno — domyślny, pierwszej strony i stron parzystych. Zasięg pierwszej
// strony dokłada do sekcji przełącznik `w:titlePg`, bez którego nagłówek
// pierwszej strony się nie pokaże.
func ooxmlZlozNaglowkiStopki(postac *shared.StudioDocumentForm) ooxmlNaglowkiArchiwum {
	wynik := ooxmlNaglowkiArchiwum{skladniki: map[string]string{}}

	wpisy := []shared.StudioHeaderFooter{}
	if len(postac.Sections) > 0 {
		wpisy = postac.Sections[0].HeadersFooters
	}
	if len(wpisy) == 0 && postac.PageSetup != nil {
		// Starsza droga kontraktu: nagłówek i stopka jednym polem nastaw
		// strony całego dokumentu.
		wpis := shared.StudioHeaderFooter{Scope: shared.StudioHeaderScopeDefault}
		if postac.PageSetup.Header != nil && *postac.PageSetup.Header != "" {
			wpis.HeaderText = postac.PageSetup.Header
		}
		if postac.PageSetup.Footer != nil && *postac.PageSetup.Footer != "" {
			wpis.FooterText = postac.PageSetup.Footer
		}
		if wpis.HeaderText != nil || wpis.FooterText != nil {
			wpisy = append(wpisy, wpis)
		}
	}
	if len(wpisy) == 0 {
		return wynik
	}

	numer := 0
	maPierwsza := false
	for _, wpis := range wpisy {
		rodzajOdwolania := "default"
		switch wpis.Scope {
		case shared.StudioHeaderScopeFirstPage:
			rodzajOdwolania = "first"
			maPierwsza = true
		case shared.StudioHeaderScopeEvenPages:
			rodzajOdwolania = "even"
		}
		for _, czesc := range []struct {
			tekst     *string
			wezel     string
			nazwa     string
			typ       string
			relacja   string
			odwolanie string
		}{
			{wpis.HeaderText, "w:hdr", "header", "header", "header", "headerReference"},
			{wpis.FooterText, "w:ftr", "footer", "footer", "footer", "footerReference"},
		} {
			if czesc.tekst == nil || strings.TrimSpace(*czesc.tekst) == "" {
				continue
			}
			numer++
			plik := czesc.nazwa + strconv.Itoa(numer) + ".xml"
			odwolanieKodu := "rIdNaglowek" + strconv.Itoa(numer)

			var akapity strings.Builder
			for _, wiersz := range strings.Split(*czesc.tekst, "\n") {
				akapity.WriteString(`<w:p><w:r><w:t xml:space="preserve">` +
					ooxmlZabezpiecz(wiersz) + `</w:t></w:r></w:p>`)
			}
			wynik.skladniki["word/"+plik] = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
				"\n" + `<` + czesc.wezel + ` xmlns:w="` + ooxmlPrzestrzenGlowna + `">` +
				akapity.String() + `</` + czesc.wezel + `>`
			wynik.powiazania += `<Relationship Id="` + odwolanieKodu +
				`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/` +
				czesc.relacja + `" Target="` + plik + `"/>`
			wynik.typyTresci += `<Override PartName="/word/` + plik +
				`" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.` +
				czesc.typ + `+xml"/>`
			wynik.odwolania += `<w:` + czesc.odwolanie + ` w:type="` + rodzajOdwolania +
				`" r:id="` + odwolanieKodu + `"/>`
		}
	}
	if maPierwsza {
		wynik.odwolania += `<w:titlePg/>`
	}
	return wynik
}
