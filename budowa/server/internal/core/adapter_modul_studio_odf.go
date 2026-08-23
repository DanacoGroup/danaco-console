// Odpowiedzialność pliku: rachunek OpenDocument modułu Studio — odczyt `.odt`
// i `.ott` do postaci dokumentu oraz złożenie jednego i drugiego z powrotem.
//
// ── Czym to liczone ─────────────────────────────────────────────────────────
// `.odt` jest archiwum ZIP z `content.xml` i `styles.xml` w środku, więc rachunek
// stoi na `archive/zip` i `encoding/xml` z biblioteki wzorcowej — dokładnie tak
// samo jak rachunek OOXML, i tym samym drzewem węzłów (`ooxmlWezel`,
// `ooxmlRozbierz`). Drugiego rozbioru XML nie zakładam: postać w ODF leży
// w węzłach obok treści, czyli w tym samym układzie, dla którego to drzewo
// powstało. Bez LibreOffice, bez `pandoc`, bez biblioteki obcej.
//
// ── Gdzie ODF różni się od OOXML i co z tego wynika ─────────────────────────
// Trzy różnice rozstrzygają kształt tego pliku:
//
//  1. POSTAĆ BEZPOŚREDNIA NIE ISTNIEJE. W Wordzie pogrubienie jednego wyrazu
//     stoi wprost przy fragmencie (`w:rPr`). W ODF każde odstępstwo od stylu
//     nazwanego musi być STYLEM AUTOMATYCZNYM o własnej nazwie, wymienionym
//     w `office:automatic-styles`. Dlatego odczyt trzyma mapę stylów
//     automatycznych, a zapis je WYTWARZA. Nazwy tych stylów są nazwami
//     technicznymi formatu pliku (tak samo nazywa je LibreOffice), a nie
//     kodem wymyślonym dla Operatora — Operator ich nigdy nie widzi.
//
//  2. MIARY IDĄ JEDNOSTKĄ W NAPISIE. `fo:page-width="21cm"`,
//     `fo:font-size="12pt"`, `fo:margin-left="0.5in"`. Przeliczenie stoi
//     w jednym miejscu (`odfDlugoscNaMilimetry`, `odfDlugoscNaPunkty`), bo
//     rozsypane po dwudziestu miejscach rozjeżdża się przy pierwszej poprawce.
//
//  3. NAGŁÓWEK I STOPKA WISZĄ NA STRONIE WZORCOWEJ. `style:master-page` wskazuje
//     `style:page-layout` z nastawami strony i niesie `style:header` oraz
//     `style:footer`. Sekcje ODF nie mają własnych nastaw strony w tym sensie,
//     w którym mają je sekcje OOXML — dokument wniesiony z `.odt` dostaje więc
//     jedną sekcję, a nie sekcje udawane. Bilans mówi to wprost.
//
// ── Rodzaj archiwum: `mimetype` musi być pierwszy i nieskompresowany ────────
// Tak stanowi norma OpenDocument i tak to sprawdzają czytniki. Składanie
// archiwum robi `wejscieZlozArchiwum`, któremu podaje się nazwę składnika
// pierwszego właśnie z tego powodu.
package core

import (
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// Nazwy składników archiwum OpenDocument, których dotyka ten rachunek.
const (
	odfSkladnikRodzaju   = "mimetype"
	odfSkladnikTresci    = "content.xml"
	odfSkladnikStylow    = "styles.xml"
	odfSkladnikManifestu = "META-INF/manifest.xml"
	odfKatalogNosnikow   = "Pictures/"

	odfRodzajDokumentu = "application/vnd.oasis.opendocument.text"
	odfRodzajSzablonu  = "application/vnd.oasis.opendocument.text-template"

	// Nazwa strony wzorcowej dokumentu składanego. „Standard" jest nazwą, którą
	// niesie każdy dokument OpenDocument — czytnik szuka jej wprost.
	odfNazwaStronyWzorcowej = "Standard"
	odfNazwaUkladuStrony    = "uklad-strony"
)

// ── Miary ───────────────────────────────────────────────────────────────────

// odfDlugoscNaMilimetry przelicza długość ODF (`21cm`, `1in`, `12pt`, `5mm`) na
// milimetry. Napis bez rozpoznanej jednostki nie daje zera po cichu — daje
// fałsz, żeby wołający zostawił nastawę domyślną, a nie margines zerowy.
func odfDlugoscNaMilimetry(wartosc string) (float64, bool) {
	liczba, jednostka, jest := odfRozbierzDlugosc(wartosc)
	if !jest {
		return 0, false
	}
	switch jednostka {
	case "mm":
		return liczba, true
	case "cm":
		return liczba * 10, true
	case "in":
		return liczba * 25.4, true
	case "pt":
		return liczba / 72 * 25.4, true
	case "pc":
		return liczba * 12 / 72 * 25.4, true
	case "px":
		// Piksel ODF liczy się przy 96 punktach na cal — tak go liczy
		// LibreOffice i tak wychodzi zgodnie z tym, co Operator widział.
		return liczba / 96 * 25.4, true
	}
	return 0, false
}

// odfDlugoscNaPunkty przelicza długość ODF na punkty typograficzne.
func odfDlugoscNaPunkty(wartosc string) (float64, bool) {
	milimetry, jest := odfDlugoscNaMilimetry(wartosc)
	if !jest {
		return 0, false
	}
	return milimetry / 25.4 * 72, true
}

// odfRozbierzDlugosc rozdziela napis na liczbę i jednostkę.
func odfRozbierzDlugosc(wartosc string) (float64, string, bool) {
	czysta := strings.TrimSpace(strings.ToLower(wartosc))
	if czysta == "" {
		return 0, "", false
	}
	granica := len(czysta)
	for i, znak := range czysta {
		if (znak >= '0' && znak <= '9') || znak == '.' || znak == '-' || znak == '+' {
			continue
		}
		granica = i
		break
	}
	liczba, err := strconv.ParseFloat(strings.TrimSpace(czysta[:granica]), 64)
	if err != nil {
		return 0, "", false
	}
	return liczba, strings.TrimSpace(czysta[granica:]), true
}

// odfMilimetryNaZapis składa napis długości w milimetrach z jednostką.
func odfMilimetryNaZapis(milimetry float64) string {
	return strconv.FormatFloat(milimetry, 'f', 3, 64) + "mm"
}

// odfPunktyNaZapis składa napis stopnia pisma albo odstępu w punktach.
func odfPunktyNaZapis(punkty float64) string {
	return strconv.FormatFloat(punkty, 'f', 2, 64) + "pt"
}

// ── Odczyt ──────────────────────────────────────────────────────────────────

// odfStanOdczytu zbiera to, co rozbiór treści musi mieć pod ręką: składaną
// postać, bilans, składniki archiwum oraz mapy stylów automatycznych.
type odfStanOdczytu struct {
	postac    *shared.StudioDocumentForm
	bilans    *shared.StudioImportBalance
	skladniki map[string][]byte
	// styleZnaku i styleAkapitu wiążą nazwę stylu automatycznego z postacią,
	// którą ten styl niesie. Bez tych map pogrubienie wyrazu w pliku Operatora
	// przepadłoby, choć plik je niesie — tylko nie przy fragmencie.
	styleZnaku    map[string]*shared.StudioCharacterFormat
	styleAkapitu  map[string]*shared.StudioParagraphFormat
	stylNadrzedny map[string]string
	// szerokosciKolumn wiąże nazwę stylu automatycznego kolumny tabeli z jej
	// szerokością. Bez tego tabela po wniesieniu ma szerokości zerowe, a to
	// jest dokładnie ten wynik, który sprawdzian odcinka postaci zakazuje.
	szerokosciKolumn map[string]float64
	biezacaSekcja    string
	poziomListy      int
	listaBiezaca     string
}

// wejscieCzytajOdf rozbiera `.odt` albo `.ott` na postać dokumentu, treść płaską
// i bilans wniesienia.
//
// Bilans jest obowiązkowy, nie ozdobny: plik Operatora niesie rzeczy, których
// ten rachunek nie odczytuje (pola obliczane, wykresy, ramki rysunkowe,
// obiekty osadzone innych programów), a przemilczenie ich zamieniłoby dokument
// okaleczony w dokument „wczytany bez uwag".
func wejscieCzytajOdf(kodDokumentu string, bajty []byte,
	format shared.StudioImportFormat) (shared.StudioDocumentForm, string,
	shared.StudioImportBalance, error) {

	bilans := shared.StudioImportBalance{Format: format}
	skladniki, err := wejscieOtworzArchiwum(bajty)
	if err != nil {
		return shared.StudioDocumentForm{}, "", bilans, err
	}
	surowaTresc, jest := skladniki[odfSkladnikTresci]
	if !jest {
		return shared.StudioDocumentForm{}, "", bilans, bladWskazaniaStudio(
			"archiwum nie niesie składnika " + odfSkladnikTresci +
				" — to nie jest dokument OpenDocument (odt, ott)")
	}
	drzewoTresci, err := ooxmlRozbierz(surowaTresc)
	if err != nil {
		return shared.StudioDocumentForm{}, "", bilans, bladWskazaniaStudio(err.Error())
	}

	postac := shared.StudioDocumentForm{DocumentId: kodDokumentu}
	stan := &odfStanOdczytu{
		postac:           &postac,
		bilans:           &bilans,
		skladniki:        skladniki,
		styleZnaku:       map[string]*shared.StudioCharacterFormat{},
		styleAkapitu:     map[string]*shared.StudioParagraphFormat{},
		stylNadrzedny:    map[string]string{},
		szerokosciKolumn: map[string]float64{},
	}

	// Arkusz stylów wchodzi PIERWSZY: akapity odwołują się do stylów nazwą,
	// więc nazwa stylu w bloku ma sens tylko wtedy, gdy styl już istnieje.
	if surowyStylow, jest := skladniki[odfSkladnikStylow]; jest {
		if drzewoStylow, err := ooxmlRozbierz(surowyStylow); err == nil {
			postac.Styles = stan.czytajStyleNazwane(drzewoStylow)
			stan.czytajStyleAutomatyczne(drzewoStylow)
			postac.PageSetup = stan.czytajNastawyStrony(drzewoStylow)
		} else {
			stan.pomin("arkusz stylów pliku jest nieczytelny", err.Error())
		}
	}
	if len(postac.Styles) == 0 {
		postac.Styles = wejscieDomyslnyArkuszStylow()
		stan.pomin("plik nie niesie arkusza stylów nazwanych",
			"dokument dostał arkusz domyślny platformy")
	}
	bilans.StylesRecovered = wejscieWskaznikCalkowity(len(postac.Styles))
	stan.czytajStyleAutomatyczne(drzewoTresci)

	if postac.PageSetup == nil {
		nastawy := wejscieDomyslneNastawyStrony("", nil)
		postac.PageSetup = &nastawy
	}

	cialo := drzewoTresci.sciezka("body", "text")
	if cialo == nil {
		return shared.StudioDocumentForm{}, "", bilans, bladWskazaniaStudio(
			"dokument OpenDocument bez ciała (office:body/office:text) — plik jest uszkodzony")
	}

	stan.zalozSekcje()
	stan.czytajCialo(cialo)

	bilans.SectionsRecovered = wejscieWskaznikCalkowity(len(postac.Sections))
	bilans.ParagraphsRecovered = wejscieWskaznikCalkowity(len(postac.Blocks))
	bilans.TablesRecognized = wejscieWskaznikCalkowity(len(postac.Tables))

	tresc := wejscieTrescZPostaci(&postac)
	wejscieUzycieStylow(&postac)
	bilans.Note = wejscieWskaznikTekstu(ooxmlZdanieBilansu(&bilans))
	return postac, tresc, bilans, nil
}

// pomin dokłada pozycję do wykazu tego, czego rachunek nie odzyskał.
func (s *odfStanOdczytu) pomin(powod, szczegol string) {
	s.bilans.Skipped = append(s.bilans.Skipped, shared.StudioSkippedItem{
		Reason: powod, Detail: wejscieWskaznikTekstu(szczegol),
	})
}

// zalozSekcje zakłada jedyną sekcję dokumentu wniesionego z ODF.
//
// Jedna, a nie kilka: nastawy strony w ODF wiszą na stronie wzorcowej, nie na
// sekcji, więc rozdzielenie dokumentu na sekcje o własnych nastawach byłoby
// wymyśleniem podziału, którego plik nie niesie.
func (s *odfStanOdczytu) zalozSekcje() {
	sekcja := shared.StudioSection{
		Id:        nowyIdentyfikator(przedrostekSekcjiStudia),
		Index:     0,
		Start:     wejscieWskaznikPoczatkuSekcji(shared.StudioSectionStartContinuous),
		PageSetup: s.postac.PageSetup,
	}
	if naglowki := s.czytajNaglowkiStopki(); len(naglowki) > 0 {
		sekcja.HeadersFooters = naglowki
	}
	s.postac.Sections = append(s.postac.Sections, sekcja)
	s.biezacaSekcja = sekcja.Id
}

// czytajCialo przechodzi ciało dokumentu blok po bloku w kolejności czytania.
func (s *odfStanOdczytu) czytajCialo(cialo *ooxmlWezel) {
	for _, wezel := range cialo.Dzieci {
		switch wezel.Nazwa.Local {
		case "p":
			s.czytajAkapit(wezel, 0)
		case "h":
			s.czytajAkapit(wezel, odfPoziomNaglowka(wezel))
		case "list":
			s.czytajListe(wezel)
		case "table":
			s.czytajTabele(wezel)
		case "section", "text-box", "frame":
			// Sekcja nazwana ODF i ramka tekstowa opakowują zwykłą treść.
			// Schodzimy do wnętrza, zamiast pomijać — pominięcie zgubiłoby
			// akapity, które w edytorze Operatora widać normalnie.
			if wezel.Nazwa.Local == "frame" {
				s.czytajObraz(wezel)
			}
			s.czytajCialo(wezel)
		case "sequence-decls", "variable-decls", "user-field-decls":
			// Deklaracje pól obliczanych. Wartości pól stoją w treści akapitów
			// i tam się odczytują; sama deklaracja nie jest treścią.
		case "soft-page-break":
			// Podział strony policzony przez edytor Operatora, nie postawiony
			// przez niego. Wniesienie go jako podziału twardego zmieniłoby
			// dokument.
		default:
			if strings.TrimSpace(wezel.Tekst) != "" {
				s.pomin("węzeł treści, którego rachunek nie odczytuje",
					wezel.Nazwa.Local)
			}
		}
	}
}

// odfPoziomNaglowka odczytuje poziom konspektu nagłówka (`text:outline-level`).
func odfPoziomNaglowka(wezel *ooxmlWezel) int {
	if poziom, jest := wezel.atrybutCalkowity("outline-level"); jest && poziom > 0 {
		if poziom > 6 {
			return 6
		}
		return poziom
	}
	return 1
}

// czytajAkapit składa blok akapitu wraz z postacią akapitu i fragmentami.
func (s *odfStanOdczytu) czytajAkapit(wezel *ooxmlWezel, poziomNaglowka int) {
	nazwaStylu := wezel.atrybut("style-name")
	postacAkapitu := s.postacAkapituStylu(nazwaStylu)
	if postacAkapitu == nil {
		postacAkapitu = &shared.StudioParagraphFormat{}
	}
	if stylNazwany := s.nazwaStyluNazwanego(nazwaStylu); stylNazwany != "" {
		postacAkapitu.StyleName = wejscieWskaznikTekstu(stylNazwany)
	}
	rodzaj := wejscieRodzajBlokuAkapit
	if poziomNaglowka > 0 {
		postacAkapitu.OutlineLevel = wejscieWskaznikCalkowity(poziomNaglowka)
		if postacAkapitu.StyleName == nil {
			postacAkapitu.StyleName = wejscieWskaznikTekstu(
				wejscieStylNaglowkaPoziomu(poziomNaglowka))
		}
		rodzaj = wejscieRodzajBlokuAkapit
	}
	if s.listaBiezaca != "" {
		postacAkapitu.ListId = wejscieWskaznikTekstu(s.listaBiezaca)
		postacAkapitu.ListLevel = wejscieWskaznikCalkowity(s.poziomListy)
	}

	blok := shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      rodzaj,
		SectionId: wejscieWskaznikTekstu(s.biezacaSekcja),
		Paragraph: postacAkapitu,
		Runs:      s.czytajFragmenty(wezel, nil),
	}
	s.postac.Blocks = append(s.postac.Blocks, blok)
}

// czytajFragmenty składa fragmenty tekstu akapitu wraz z postacią znaku.
//
// Postać odziedziczona jedzie w dół drzewa: `text:span` w `text:span` jest
// w ODF zwykły, a postać zewnętrzna obowiązuje wewnątrz — inaczej pogrubiony
// akapit z jednym wyrazem w kursywie oddałby ten wyraz bez pogrubienia.
func (s *odfStanOdczytu) czytajFragmenty(wezel *ooxmlWezel,
	odziedziczona *shared.StudioCharacterFormat) []shared.StudioDocumentRun {

	fragmenty := make([]shared.StudioDocumentRun, 0, len(wezel.Dzieci)+1)
	if tekst := wezel.Tekst; tekst != "" {
		fragmenty = append(fragmenty, shared.StudioDocumentRun{
			Text: tekst, Format: odfKopiaPostaciZnaku(odziedziczona),
		})
	}
	for _, dziecko := range wezel.Dzieci {
		switch dziecko.Nazwa.Local {
		case "span":
			postac := odfZlaczPostacZnaku(odziedziczona,
				s.postacZnakuStylu(dziecko.atrybut("style-name")))
			fragmenty = append(fragmenty, s.czytajFragmenty(dziecko, postac)...)
		case "a":
			// Odsyłacz: treść należy do dokumentu, adres do aparatu.
			wewnetrzne := s.czytajFragmenty(dziecko, odziedziczona)
			fragmenty = append(fragmenty, wewnetrzne...)
			s.czytajOdsylacz(dziecko, wewnetrzne)
		case "s":
			// `text:s` to ciąg odstępów zwiniętych do jednego węzła.
			ile := 1
			if liczba, jest := dziecko.atrybutCalkowity("c"); jest && liczba > 0 {
				ile = liczba
			}
			fragmenty = append(fragmenty, shared.StudioDocumentRun{
				Text: strings.Repeat(" ", ile), Format: odfKopiaPostaciZnaku(odziedziczona),
			})
		case "tab":
			fragmenty = append(fragmenty, shared.StudioDocumentRun{
				Text: "\t", Format: odfKopiaPostaciZnaku(odziedziczona),
			})
		case "line-break":
			fragmenty = append(fragmenty, shared.StudioDocumentRun{
				Text: "\n", Format: odfKopiaPostaciZnaku(odziedziczona),
			})
		case "frame":
			s.czytajObraz(dziecko)
		case "note":
			s.czytajPrzypis(dziecko)
		case "bookmark", "bookmark-start", "reference-mark", "reference-mark-start":
			s.czytajZakladke(dziecko)
		default:
			// Pole obliczane (`text:date`, `text:page-number`, `text:variable-set`
			// i pokrewne) niesie WARTOŚĆ w treści węzła. Wartość wchodzi do
			// dokumentu, bo Operator ją widział; wyrażenie, które ją policzyło,
			// idzie do wykazu pominiętych, bo tego rachunek nie odtwarza.
			if tekst := odfTekstWglab(dziecko); tekst != "" {
				fragmenty = append(fragmenty, shared.StudioDocumentRun{
					Text: tekst, Format: odfKopiaPostaciZnaku(odziedziczona),
				})
				s.pomin("pole obliczane wniesione jako wartość, bez wyrażenia",
					dziecko.Nazwa.Local)
			}
		}
		if dziecko.Tekst != "" && dziecko.Nazwa.Local == "" {
			fragmenty = append(fragmenty, shared.StudioDocumentRun{
				Text: dziecko.Tekst, Format: odfKopiaPostaciZnaku(odziedziczona),
			})
		}
	}
	return fragmenty
}

// odfTekstWglab zbiera treść tekstową węzła wraz z dziećmi.
func odfTekstWglab(wezel *ooxmlWezel) string {
	var zbior strings.Builder
	zbior.WriteString(wezel.Tekst)
	for _, dziecko := range wezel.Dzieci {
		zbior.WriteString(odfTekstWglab(dziecko))
	}
	return zbior.String()
}

// czytajListe wchodzi w listę ODF: `text:list` niesie `text:list-item`, a w nich
// stoją akapity. Poziom listy liczy się zagnieżdżeniem — tak samo jak w edytorze.
func (s *odfStanOdczytu) czytajListe(wezel *ooxmlWezel) {
	poprzedniaLista, poprzedniPoziom := s.listaBiezaca, s.poziomListy
	if s.listaBiezaca == "" {
		s.listaBiezaca = s.zalozListe(wezel)
	}
	s.poziomListy++

	for _, dziecko := range wezel.Dzieci {
		switch dziecko.Nazwa.Local {
		case "list-item", "list-header":
			s.czytajCialo(dziecko)
		case "list":
			s.czytajListe(dziecko)
		}
	}
	s.listaBiezaca, s.poziomListy = poprzedniaLista, poprzedniPoziom
}

// zalozListe dokłada definicję listy do postaci dokumentu.
//
// Rodzaj listy rozstrzyga styl listy: styl o nazwie mówiącej o numeracji daje
// listę numerowaną, pozostałe — punktowaną. To jest odtworzenie, nie odczyt,
// bo pełna definicja stylu listy stoi w `text:list-style`, którego szczeble
// rachunek odczytuje wyłącznie w tym zakresie; wykaz pominiętych to mówi.
func (s *odfStanOdczytu) zalozListe(wezel *ooxmlWezel) string {
	rodzaj := shared.StudioListKind(shared.StudioListKindBullet)
	nazwaStylu := strings.ToLower(wezel.atrybut("style-name"))
	if strings.Contains(nazwaStylu, "number") || strings.Contains(nazwaStylu, "numer") {
		rodzaj = shared.StudioListKindNumber
	}
	lista := shared.StudioListDefinition{
		Id:   nowyIdentyfikator(przedrostekListyWejscia),
		Kind: rodzaj,
	}
	s.postac.Lists = append(s.postac.Lists, lista)
	s.pomin("szczeble stylu listy odtworzone, nie odczytane",
		"rachunek bierze rodzaj listy ze stylu; wcięcia i znaki szczebli "+
			"zostają domyślne platformy")
	return lista.Id
}

// przedrostekListyWejscia znakuje definicje list odczytane z pliku.
const przedrostekListyWejscia = "studio-list-"

// czytajTabele składa tabelę dokumentu wraz z komórkami, scaleniami
// i szerokościami kolumn.
func (s *odfStanOdczytu) czytajTabele(wezel *ooxmlWezel) {
	tabela := shared.StudioDocumentTable{
		Id: nowyIdentyfikator(przedrostekTabeliStudia),
	}
	if nazwa := s.nazwaStyluNazwanego(wezel.atrybut("style-name")); nazwa != "" {
		tabela.StyleName = wejscieWskaznikTekstu(nazwa)
	}

	szerokosci := []float64{}
	for _, kolumna := range odfWezlyWglab(wezel, "table-column") {
		ile := 1
		if liczba, jest := kolumna.atrybutCalkowity("number-columns-repeated"); jest && liczba > 0 {
			ile = liczba
		}
		szerokosc := s.szerokoscKolumny(kolumna.atrybut("style-name"))
		for i := 0; i < ile; i++ {
			szerokosci = append(szerokosci, szerokosc)
		}
	}

	wiersze := odfWierszeTabeli(wezel)
	kolumn := len(szerokosci)
	for numerWiersza, wiersz := range wiersze {
		numerKolumny := 0
		if numerWiersza == 0 && odfWNaglowku(wezel, wiersz) {
			tabela.HeaderRows = wejscieWskaznikCalkowity(1)
			tabela.RepeatHeader = wejscieWskaznikLogiczny(true)
		}
		// przykrytychDoPominiecia liczy komórki `table:covered-table-cell`, które
		// należą do scalenia POZIOMEGO policzonego już przez rozpiętość komórki
		// scalającej.
		//
		// OpenDocument zapisuje scalenie poziome DWA RAZY: raz jako
		// `table:number-columns-spanned` komórki scalającej, raz jako
		// `table:covered-table-cell` na każdej z przykrytych kolumn — jest ich
		// dokładnie o jedną mniej niż rozpiętość. Liczenie obu naraz dawało
		// tabelę trzykolumnową jako czterokolumnową, a kolumna czwarta wychodziła
		// bez szerokości, czyli w dokumencie niewidoczna.
		//
		// Licznik, a nie granica kolumny: po przejściu komórki scalającej numer
		// kolumny stoi już ZA scaleniem, więc porównanie z granicą nie odróżnia
		// przykrycia poziomego od pionowego. Liczba przykrytych komórek do
		// pominięcia jest jednoznaczna.
		przykrytychDoPominiecia := 0
		for _, wezelKomorki := range wiersz.Dzieci {
			switch wezelKomorki.Nazwa.Local {
			case "table-cell":
				komorka := s.czytajKomorke(wezelKomorki, numerWiersza, numerKolumny)
				tabela.Cells = append(tabela.Cells, komorka)
				rozpietosc := 1
				if komorka.ColumnSpan != nil && *komorka.ColumnSpan > 1 {
					rozpietosc = *komorka.ColumnSpan
				}
				// Kolumny przykryte scaleniem poziomym dostają własne komórki
				// oznaczone jako scalone — tak samo, jak robi to rozbiór OOXML.
				// Jedna umowa dla obu formatów, bo postać dokumentu jest jedna:
				// bez tego tabela z `.odt` i ta sama tabela z `.docx` miałyby
				// różną liczbę komórek.
				for przesuniecie := 1; przesuniecie < rozpietosc; przesuniecie++ {
					tabela.Cells = append(tabela.Cells, shared.StudioTableCell{
						Row: numerWiersza, Column: numerKolumny + przesuniecie,
						Merged: wejscieWskaznikLogiczny(true),
					})
				}
				if rozpietosc > 1 {
					przykrytychDoPominiecia += rozpietosc - 1
				}
				numerKolumny += rozpietosc
			case "covered-table-cell":
				// Komórka przykryta scaleniem POZIOMYM tego wiersza jest już
				// policzona wyżej i kolumny nie zajmuje po raz drugi.
				if przykrytychDoPominiecia > 0 {
					przykrytychDoPominiecia--
					continue
				}
				// Zostaje przykrycie scaleniem PIONOWYM z wiersza wcześniejszego.
				// To jest osobna kolumna i musi wejść do wykazu, bo inaczej
				// wiersz miałby mniej komórek niż tabela kolumn.
				tabela.Cells = append(tabela.Cells, shared.StudioTableCell{
					Row: numerWiersza, Column: numerKolumny,
					Merged: wejscieWskaznikLogiczny(true),
				})
				numerKolumny++
			}
		}
		if numerKolumny > kolumn {
			kolumn = numerKolumny
		}
	}
	tabela.Rows = len(wiersze)
	tabela.Columns = kolumn
	if len(szerokosci) > 0 {
		tabela.ColumnWidthsMm = szerokosci
		suma := 0.0
		for _, szerokosc := range szerokosci {
			suma += szerokosc
		}
		if suma > 0 {
			tabela.WidthMm = wejscieWskaznikRzeczywisty(suma)
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

// odfWierszeTabeli wybiera wiersze tabeli wraz z wierszami nagłówka.
func odfWierszeTabeli(tabela *ooxmlWezel) []*ooxmlWezel {
	wiersze := []*ooxmlWezel{}
	for _, dziecko := range tabela.Dzieci {
		switch dziecko.Nazwa.Local {
		case "table-row":
			wiersze = append(wiersze, dziecko)
		case "table-header-rows", "table-rows", "table-row-group":
			wiersze = append(wiersze, odfWierszeTabeli(dziecko)...)
		}
	}
	return wiersze
}

// odfWNaglowku mówi, czy wiersz stoi w grupie wierszy nagłówkowych.
func odfWNaglowku(tabela, wiersz *ooxmlWezel) bool {
	for _, dziecko := range tabela.Dzieci {
		if dziecko.Nazwa.Local != "table-header-rows" {
			continue
		}
		for _, wnetrze := range dziecko.Dzieci {
			if wnetrze == wiersz {
				return true
			}
		}
	}
	return false
}

// czytajKomorke składa komórkę tabeli wraz z treścią i scaleniem.
func (s *odfStanOdczytu) czytajKomorke(wezel *ooxmlWezel,
	wiersz, kolumna int) shared.StudioTableCell {

	komorka := shared.StudioTableCell{Row: wiersz, Column: kolumna}
	if rozpietosc, jest := wezel.atrybutCalkowity("number-columns-spanned"); jest && rozpietosc > 1 {
		komorka.ColumnSpan = wejscieWskaznikCalkowity(rozpietosc)
	}
	if rozpietosc, jest := wezel.atrybutCalkowity("number-rows-spanned"); jest && rozpietosc > 1 {
		komorka.RowSpan = wejscieWskaznikCalkowity(rozpietosc)
	}
	akapity := []string{}
	for _, dziecko := range wezel.Dzieci {
		if dziecko.Nazwa.Local == "p" || dziecko.Nazwa.Local == "h" {
			akapity = append(akapity, odfTekstWglab(dziecko))
		}
	}
	if len(akapity) > 0 {
		komorka.Text = wejscieWskaznikTekstu(strings.Join(akapity, "\n"))
	}
	if postac := s.postacAkapituStylu(wezel.atrybut("style-name")); postac != nil {
		komorka.Paragraph = postac
	}
	return komorka
}

// czytajObraz dokłada obraz osadzony do wykazu obiektów i liczy go w bilansie.
//
// Bajty obrazu zostają w archiwum: odłożenie ich do magazynu zasobów wymaga
// kontekstu żądania, którego ten rachunek nie ma. Wołający odkłada je po
// odczycie, a tutaj powstaje obiekt wraz ze wskazaniem składnika.
func (s *odfStanOdczytu) czytajObraz(wezel *ooxmlWezel) {
	obraz := ooxmlSzukajWglabWezel(wezel, "image")
	if obraz == nil {
		if wnetrze := ooxmlSzukajWglabWezel(wezel, "text-box"); wnetrze != nil {
			// Ramka z tekstem, nie z obrazem — jej treść czyta `czytajCialo`.
			return
		}
		s.pomin("ramka rysunkowa bez obrazu", wezel.atrybut("name"))
		s.bilans.ImagesSkipped = wejscieWskaznikCalkowity(
			wartoscCalkowita(s.bilans.ImagesSkipped) + 1)
		return
	}
	sciezka := ""
	for _, atrybut := range obraz.Atrybuty {
		if atrybut.Name.Local == "href" {
			sciezka = atrybut.Value
		}
	}
	if sciezka == "" {
		s.pomin("obraz osadzony bez wskazania składnika archiwum", wezel.atrybut("name"))
		s.bilans.ImagesSkipped = wejscieWskaznikCalkowity(
			wartoscCalkowita(s.bilans.ImagesSkipped) + 1)
		return
	}
	obiekt := shared.StudioDocumentObject{
		Id:      nowyIdentyfikator(przedrostekObiektuStudia),
		Kind:    shared.StudioObjectKindImage,
		Source:  wejscieWskaznikZrodlaObiektu(shared.StudioObjectSourceFile),
		AltText: wejscieWskaznikTekstu(odfOpisObrazu(wezel, sciezka)),
	}
	if szerokosc, jest := odfDlugoscNaMilimetry(wezel.atrybut("width")); jest {
		obiekt.WidthMm = wejscieWskaznikRzeczywisty(szerokosc)
	}
	if wysokosc, jest := odfDlugoscNaMilimetry(wezel.atrybut("height")); jest {
		obiekt.HeightMm = wejscieWskaznikRzeczywisty(wysokosc)
	}
	s.postac.Objects = append(s.postac.Objects, obiekt)
	s.postac.Blocks = append(s.postac.Blocks, shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuObiekt,
		SectionId: wejscieWskaznikTekstu(s.biezacaSekcja),
		ObjectId:  wejscieWskaznikTekstu(obiekt.Id),
	})
	s.bilans.ImagesEmbedded = wejscieWskaznikCalkowity(
		wartoscCalkowita(s.bilans.ImagesEmbedded) + 1)
}

// odfOpisObrazu składa tekst zastępczy obrazu: opis z pliku, a gdy go nie ma —
// wskazanie składnika archiwum, z którego obraz pochodzi.
func odfOpisObrazu(ramka *ooxmlWezel, sciezka string) string {
	if opis := ooxmlSzukajWglabWezel(ramka, "desc"); opis != nil {
		if tekst := strings.TrimSpace(odfTekstWglab(opis)); tekst != "" {
			return tekst
		}
	}
	if nazwa := strings.TrimSpace(ramka.atrybut("name")); nazwa != "" {
		return nazwa
	}
	return "obraz ze składnika " + sciezka
}

// czytajPrzypis dokłada przypis do aparatu dokumentu.
func (s *odfStanOdczytu) czytajPrzypis(wezel *ooxmlWezel) {
	rodzaj := shared.StudioApparatusKind(shared.StudioApparatusKindFootnote)
	if strings.Contains(wezel.atrybut("note-class"), "endnote") {
		rodzaj = shared.StudioApparatusKindEndnote
	}
	element := shared.StudioApparatusItem{
		Id:   nowyIdentyfikator(przedrostekAparatuWejscia),
		Kind: rodzaj,
	}
	if znacznik := ooxmlSzukajWglabWezel(wezel, "note-citation"); znacznik != nil {
		element.Number = wejscieWskaznikTekstu(strings.TrimSpace(odfTekstWglab(znacznik)))
	}
	if cialo := ooxmlSzukajWglabWezel(wezel, "note-body"); cialo != nil {
		element.Text = wejscieWskaznikTekstu(strings.TrimSpace(odfTekstWglab(cialo)))
	}
	s.postac.Apparatus = append(s.postac.Apparatus, element)
}

// czytajZakladke dokłada zakładkę do aparatu dokumentu.
func (s *odfStanOdczytu) czytajZakladke(wezel *ooxmlWezel) {
	nazwa := strings.TrimSpace(wezel.atrybut("name"))
	if nazwa == "" {
		return
	}
	s.postac.Apparatus = append(s.postac.Apparatus, shared.StudioApparatusItem{
		Id:    nowyIdentyfikator(przedrostekAparatuWejscia),
		Kind:  shared.StudioApparatusKindBookmark,
		Label: wejscieWskaznikTekstu(nazwa),
	})
}

// czytajOdsylacz dokłada odsyłacz do aparatu dokumentu.
func (s *odfStanOdczytu) czytajOdsylacz(wezel *ooxmlWezel,
	fragmenty []shared.StudioDocumentRun) {

	adres := ""
	for _, atrybut := range wezel.Atrybuty {
		if atrybut.Name.Local == "href" {
			adres = atrybut.Value
		}
	}
	if strings.TrimSpace(adres) == "" {
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

// ── Style ───────────────────────────────────────────────────────────────────

// czytajStyleNazwane składa arkusz stylów nazwanych z `office:styles`.
func (s *odfStanOdczytu) czytajStyleNazwane(drzewo *ooxmlWezel) []shared.StudioNamedStyle {
	wykaz := drzewo.dziecko("styles")
	if wykaz == nil {
		return nil
	}
	style := make([]shared.StudioNamedStyle, 0, len(wykaz.Dzieci))
	for _, wezel := range wykaz.dzieci("style") {
		nazwa := strings.TrimSpace(wezel.atrybut("name"))
		if nazwa == "" {
			continue
		}
		rodzina := wezel.atrybut("family")
		styl := shared.StudioNamedStyle{
			Name:        nazwa,
			DisplayName: odfWskaznikNiepusty(wezel.atrybut("display-name")),
			Kind:        odfRodzajStylu(rodzina),
			BasedOn:     odfWskaznikNiepusty(wezel.atrybut("parent-style-name")),
			NextStyle:   odfWskaznikNiepusty(wezel.atrybut("next-style-name")),
			Character:   odfCzytajPostacZnaku(wezel.dziecko("text-properties")),
			Paragraph:   odfCzytajPostacAkapitu(wezel.dziecko("paragraph-properties")),
		}
		if styl.Paragraph != nil {
			if poziom, jest := wezel.atrybutCalkowity("default-outline-level"); jest && poziom > 0 {
				styl.Paragraph.OutlineLevel = wejscieWskaznikCalkowity(poziom)
			}
		}
		style = append(style, styl)
	}
	return style
}

// czytajStyleAutomatyczne zbiera style automatyczne — czyli postać bezpośrednią
// pliku ODF — do dwóch map: postaci znaku i postaci akapitu.
func (s *odfStanOdczytu) czytajStyleAutomatyczne(drzewo *ooxmlWezel) {
	wykaz := drzewo.dziecko("automatic-styles")
	if wykaz == nil {
		return
	}
	for _, wezel := range wykaz.dzieci("style") {
		nazwa := strings.TrimSpace(wezel.atrybut("name"))
		if nazwa == "" {
			continue
		}
		if rodzic := strings.TrimSpace(wezel.atrybut("parent-style-name")); rodzic != "" {
			s.stylNadrzedny[nazwa] = rodzic
		}
		if postac := odfCzytajPostacZnaku(wezel.dziecko("text-properties")); postac != nil {
			s.styleZnaku[nazwa] = postac
		}
		if postac := odfCzytajPostacAkapitu(wezel.dziecko("paragraph-properties")); postac != nil {
			s.styleAkapitu[nazwa] = postac
		}
		if cechy := wezel.dziecko("table-column-properties"); cechy != nil {
			if szerokosc, jest := odfDlugoscNaMilimetry(cechy.atrybut("column-width")); jest {
				s.szerokosciKolumn[nazwa] = szerokosc
			}
		}
	}
}

// postacZnakuStylu oddaje postać znaku stylu automatycznego albo nazwanego.
func (s *odfStanOdczytu) postacZnakuStylu(nazwa string) *shared.StudioCharacterFormat {
	nazwa = strings.TrimSpace(nazwa)
	if nazwa == "" {
		return nil
	}
	if postac, jest := s.styleZnaku[nazwa]; jest {
		return odfKopiaPostaciZnaku(postac)
	}
	for _, styl := range s.postac.Styles {
		if styl.Name == nazwa {
			return odfKopiaPostaciZnaku(styl.Character)
		}
	}
	return nil
}

// postacAkapituStylu oddaje postać akapitu stylu automatycznego.
func (s *odfStanOdczytu) postacAkapituStylu(nazwa string) *shared.StudioParagraphFormat {
	nazwa = strings.TrimSpace(nazwa)
	if nazwa == "" {
		return nil
	}
	if postac, jest := s.styleAkapitu[nazwa]; jest {
		kopia := *postac
		return &kopia
	}
	return nil
}

// nazwaStyluNazwanego rozstrzyga, którym stylem NAZWANYM akapit się posługuje.
// Styl automatyczny nazwany nie jest — niesie za to wskazanie stylu nadrzędnego,
// a to właśnie ten styl Operator wybrał w edytorze.
func (s *odfStanOdczytu) nazwaStyluNazwanego(nazwa string) string {
	nazwa = strings.TrimSpace(nazwa)
	if nazwa == "" {
		return ""
	}
	if rodzic, jest := s.stylNadrzedny[nazwa]; jest {
		return rodzic
	}
	for _, styl := range s.postac.Styles {
		if styl.Name == nazwa {
			return nazwa
		}
	}
	return ""
}

// szerokoscKolumny odczytuje szerokość kolumny tabeli ze stylu automatycznego.
func (s *odfStanOdczytu) szerokoscKolumny(nazwaStylu string) float64 {
	nazwa := strings.TrimSpace(nazwaStylu)
	if nazwa == "" {
		return 0
	}
	// Szerokość kolumny stoi w `style:table-column-properties` stylu
	// automatycznego — zebrana przy czytaniu stylów, nie liczona tutaj drugi raz.
	return s.szerokosciKolumn[nazwa]
}

// czytajNastawyStrony odczytuje nastawy strony z układu strony wzorcowej.
func (s *odfStanOdczytu) czytajNastawyStrony(drzewo *ooxmlWezel) *shared.StudioPageSetup {
	wykaz := drzewo.dziecko("automatic-styles")
	if wykaz == nil {
		return nil
	}
	uklad := (*ooxmlWezel)(nil)
	for _, wezel := range wykaz.dzieci("page-layout") {
		uklad = wezel
		break
	}
	if uklad == nil {
		return nil
	}
	cechy := uklad.dziecko("page-layout-properties")
	if cechy == nil {
		return nil
	}
	nastawy := wejscieDomyslneNastawyStrony("", nil)
	szerokosc, maSzerokosc := odfDlugoscNaMilimetry(cechy.atrybut("page-width"))
	wysokosc, maWysokosc := odfDlugoscNaMilimetry(cechy.atrybut("page-height"))
	if maSzerokosc && maWysokosc {
		nastawy.WidthMm = wejscieWskaznikRzeczywisty(szerokosc)
		nastawy.HeightMm = wejscieWskaznikRzeczywisty(wysokosc)
		if nazwa, jest := wejscieNazwaNosnikaZWymiarow(szerokosc, wysokosc); jest {
			nastawy.PageSize = wejscieWskaznikTekstu(nazwa)
		}
	}
	for _, margines := range []struct {
		atrybut string
		cel     **int
	}{
		{"margin-top", &nastawy.MarginTop},
		{"margin-bottom", &nastawy.MarginBottom},
		{"margin-left", &nastawy.MarginLeft},
		{"margin-right", &nastawy.MarginRight},
	} {
		if wartosc, jest := odfDlugoscNaMilimetry(cechy.atrybut(margines.atrybut)); jest {
			*margines.cel = wejscieWskaznikCalkowity(int(wartosc + 0.5))
		}
	}
	if kierunek := strings.ToLower(cechy.atrybut("print-orientation")); kierunek != "" {
		if kierunek == "landscape" {
			nastawy.Orientation = wejscieWskaznikOrientacji(shared.StudioPageOrientationPozioma)
		} else {
			nastawy.Orientation = wejscieWskaznikOrientacji(shared.StudioPageOrientationPionowa)
		}
	}
	if kolumny := uklad.dziecko("columns"); kolumny != nil {
		if ile, jest := kolumny.atrybutCalkowity("column-count"); jest && ile > 0 {
			nastawy.Columns = wejscieWskaznikCalkowity(ile)
		}
	}
	return &nastawy
}

// czytajNaglowkiStopki odczytuje nagłówek i stopkę ze strony wzorcowej.
//
// ODF trzyma nagłówek strony pierwszej i stron lewych osobnymi węzłami
// (`style:header-first`, `style:header-left`) — każdy wchodzi jako własny zasięg,
// bo Właściciel wymaga nagłówka osobnego dla pierwszej strony i stron parzystych.
func (s *odfStanOdczytu) czytajNaglowkiStopki() []shared.StudioHeaderFooter {
	surowe, jest := s.skladniki[odfSkladnikStylow]
	if !jest {
		return nil
	}
	drzewo, err := ooxmlRozbierz(surowe)
	if err != nil {
		return nil
	}
	wzorce := drzewo.dziecko("master-styles")
	if wzorce == nil {
		return nil
	}
	naglowki := []shared.StudioHeaderFooter{}
	for _, strona := range wzorce.dzieci("master-page") {
		for _, para := range []struct {
			naglowek string
			stopka   string
			zasieg   shared.StudioHeaderScope
		}{
			{"header", "footer", shared.StudioHeaderScopeDefault},
			{"header-first", "footer-first", shared.StudioHeaderScopeFirstPage},
			{"header-left", "footer-left", shared.StudioHeaderScopeEvenPages},
		} {
			wpis := shared.StudioHeaderFooter{Scope: para.zasieg}
			maCos := false
			if wezel := strona.dziecko(para.naglowek); wezel != nil {
				if tekst := strings.TrimSpace(odfTekstWglab(wezel)); tekst != "" {
					wpis.HeaderText = wejscieWskaznikTekstu(tekst)
					maCos = true
				}
			}
			if wezel := strona.dziecko(para.stopka); wezel != nil {
				if tekst := strings.TrimSpace(odfTekstWglab(wezel)); tekst != "" {
					wpis.FooterText = wejscieWskaznikTekstu(tekst)
					maCos = true
				}
			}
			if maCos {
				naglowki = append(naglowki, wpis)
			}
		}
	}
	return naglowki
}

// odfCzytajPostacZnaku składa postać znaku z `style:text-properties`.
func odfCzytajPostacZnaku(wezel *ooxmlWezel) *shared.StudioCharacterFormat {
	if wezel == nil {
		return nil
	}
	postac := shared.StudioCharacterFormat{}
	niepusta := false

	if krojPisma := odfPierwszyNiepusty(wezel, "font-name", "font-family"); krojPisma != "" {
		postac.FontFamily = wejscieWskaznikTekstu(strings.Trim(krojPisma, "'\""))
		niepusta = true
	}
	if stopien, jest := odfDlugoscNaPunkty(wezel.atrybut("font-size")); jest {
		postac.FontSizePt = wejscieWskaznikRzeczywisty(stopien)
		niepusta = true
	}
	if waga := strings.ToLower(wezel.atrybut("font-weight")); waga != "" {
		gruby := waga == "bold" || waga == "600" || waga == "700" || waga == "800" || waga == "900"
		postac.Bold = wejscieWskaznikLogiczny(gruby)
		niepusta = true
	}
	if odmiana := strings.ToLower(wezel.atrybut("font-style")); odmiana != "" {
		postac.Italic = wejscieWskaznikLogiczny(odmiana == "italic" || odmiana == "oblique")
		niepusta = true
	}
	if odmiana := strings.ToLower(wezel.atrybut("text-underline-style")); odmiana != "" {
		postac.Underline = odfPodkreslenie(odmiana)
		niepusta = true
	}
	if odmiana := strings.ToLower(wezel.atrybut("text-line-through-style")); odmiana != "" {
		postac.Strikethrough = wejscieWskaznikLogiczny(odmiana != "none")
		niepusta = true
	}
	if pozycja := strings.ToLower(wezel.atrybut("text-position")); pozycja != "" {
		switch {
		case strings.HasPrefix(pozycja, "super"), strings.HasPrefix(pozycja, "33"):
			postac.Superscript = wejscieWskaznikLogiczny(true)
			niepusta = true
		case strings.HasPrefix(pozycja, "sub"), strings.HasPrefix(pozycja, "-"):
			postac.Subscript = wejscieWskaznikLogiczny(true)
			niepusta = true
		}
	}
	if barwa := ooxmlBezKrzyzyka(wezel.atrybut("color")); barwa != "" {
		postac.Color = wejscieWskaznikTekstu(barwa)
		niepusta = true
	}
	if barwa := ooxmlBezKrzyzyka(wezel.atrybut("background-color")); barwa != "" &&
		strings.ToLower(barwa) != "transparent" {
		postac.HighlightColor = wejscieWskaznikTekstu(barwa)
		niepusta = true
	}
	if odstep, jest := odfDlugoscNaPunkty(wezel.atrybut("letter-spacing")); jest {
		postac.LetterSpacingPt = wejscieWskaznikRzeczywisty(odstep)
		niepusta = true
	}
	switch strings.ToLower(wezel.atrybut("font-variant")) {
	case "small-caps":
		postac.SmallCaps = wejscieWskaznikLogiczny(true)
		niepusta = true
	}
	switch strings.ToLower(wezel.atrybut("text-transform")) {
	case "uppercase":
		postac.AllCaps = wejscieWskaznikLogiczny(true)
		niepusta = true
	}
	if jezyk := strings.TrimSpace(wezel.atrybut("language")); jezyk != "" {
		if kraj := strings.TrimSpace(wezel.atrybut("country")); kraj != "" {
			jezyk += "-" + kraj
		}
		postac.Language = wejscieWskaznikTekstu(jezyk)
		niepusta = true
	}
	if !niepusta {
		return nil
	}
	return &postac
}

// odfCzytajPostacAkapitu składa postać akapitu z `style:paragraph-properties`.
func odfCzytajPostacAkapitu(wezel *ooxmlWezel) *shared.StudioParagraphFormat {
	if wezel == nil {
		return nil
	}
	postac := shared.StudioParagraphFormat{}
	niepusta := false

	if wyrownanie := odfWyrownanie(wezel.atrybut("text-align")); wyrownanie != nil {
		postac.Align = wyrownanie
		niepusta = true
	}
	for _, miara := range []struct {
		atrybut string
		cel     **float64
	}{
		{"margin-left", &postac.IndentLeftMm},
		{"margin-right", &postac.IndentRightMm},
		{"text-indent", &postac.FirstLineIndentMm},
	} {
		if wartosc, jest := odfDlugoscNaMilimetry(wezel.atrybut(miara.atrybut)); jest {
			*miara.cel = wejscieWskaznikRzeczywisty(wartosc)
			niepusta = true
		}
	}
	for _, odstep := range []struct {
		atrybut string
		cel     **float64
	}{
		{"margin-top", &postac.SpaceBeforePt},
		{"margin-bottom", &postac.SpaceAfterPt},
	} {
		if wartosc, jest := odfDlugoscNaPunkty(wezel.atrybut(odstep.atrybut)); jest {
			*odstep.cel = wejscieWskaznikRzeczywisty(wartosc)
			niepusta = true
		}
	}
	if interlinia := strings.TrimSpace(wezel.atrybut("line-height")); interlinia != "" {
		if strings.HasSuffix(interlinia, "%") {
			if procent, err := strconv.ParseFloat(strings.TrimSuffix(interlinia, "%"), 64); err == nil {
				postac.LineSpacingRule = wejscieWskaznikZasadyInterlinii(
					shared.StudioLineSpacingRuleMultiple)
				postac.LineSpacingValue = wejscieWskaznikRzeczywisty(procent / 100)
				niepusta = true
			}
		} else if punkty, jest := odfDlugoscNaPunkty(interlinia); jest {
			postac.LineSpacingRule = wejscieWskaznikZasadyInterlinii(
				shared.StudioLineSpacingRuleExactly)
			postac.LineSpacingValue = wejscieWskaznikRzeczywisty(punkty)
			niepusta = true
		}
	}
	if barwa := ooxmlBezKrzyzyka(wezel.atrybut("background-color")); barwa != "" &&
		strings.ToLower(barwa) != "transparent" {
		postac.ShadingColor = wejscieWskaznikTekstu(barwa)
		niepusta = true
	}
	if razem := strings.ToLower(wezel.atrybut("keep-with-next")); razem != "" {
		postac.KeepWithNext = wejscieWskaznikLogiczny(razem == "always")
		niepusta = true
	}
	if razem := strings.ToLower(wezel.atrybut("keep-together")); razem != "" {
		postac.KeepLines = wejscieWskaznikLogiczny(razem == "always")
		niepusta = true
	}
	if wdowy := wezel.atrybut("widows"); wdowy != "" {
		postac.WidowControl = wejscieWskaznikLogiczny(wdowy != "0")
		niepusta = true
	}
	if obramowanie := odfObramowanie(wezel.atrybut("border")); obramowanie != nil {
		postac.Border = obramowanie
		niepusta = true
	}
	if tabulatory := odfCzytajTabulatory(wezel.dziecko("tab-stops")); len(tabulatory) > 0 {
		postac.TabStops = tabulatory
		niepusta = true
	}
	if !niepusta {
		return nil
	}
	return &postac
}

// odfCzytajTabulatory składa tabulatory akapitu wraz z rodzajem i znakiem
// wiodącym — bez nich chwyty linijki nie mają czego przenieść.
func odfCzytajTabulatory(wezel *ooxmlWezel) []shared.StudioTabStop {
	if wezel == nil {
		return nil
	}
	tabulatory := []shared.StudioTabStop{}
	for _, dziecko := range wezel.dzieci("tab-stop") {
		pozycja, jest := odfDlugoscNaMilimetry(dziecko.atrybut("position"))
		if !jest {
			continue
		}
		tabulator := shared.StudioTabStop{
			PositionMm: pozycja,
			Kind:       shared.StudioTabKind(shared.StudioTabKindLeft),
		}
		switch strings.ToLower(dziecko.atrybut("type")) {
		case "right":
			tabulator.Kind = shared.StudioTabKindRight
		case "center":
			tabulator.Kind = shared.StudioTabKindCenter
		case "char":
			tabulator.Kind = shared.StudioTabKindDecimal
		}
		switch strings.TrimSpace(dziecko.atrybut("leader-text")) {
		case ".":
			tabulator.Leader = ooxmlZnakWiodacy("dot")
		case "-":
			tabulator.Leader = ooxmlZnakWiodacy("hyphen")
		case "_":
			tabulator.Leader = ooxmlZnakWiodacy("underscore")
		}
		tabulatory = append(tabulatory, tabulator)
	}
	return tabulatory
}

// odfObramowanie odczytuje obramowanie zapisane skrótem `0.5pt solid #000000`.
func odfObramowanie(wartosc string) *shared.StudioBorder {
	czysta := strings.TrimSpace(wartosc)
	if czysta == "" || strings.EqualFold(czysta, "none") {
		return nil
	}
	obramowanie := shared.StudioBorder{Style: shared.StudioBorderStyleSingle}
	for _, czlon := range strings.Fields(czysta) {
		switch {
		case strings.HasPrefix(czlon, "#"):
			obramowanie.Color = wejscieWskaznikTekstu(ooxmlBezKrzyzyka(czlon))
		case czlon == "solid", czlon == "double", czlon == "dashed", czlon == "dotted":
			obramowanie.Style = ooxmlOdmianaObramowania(czlon)
		default:
			if punkty, jest := odfDlugoscNaPunkty(czlon); jest {
				obramowanie.WidthPt = wejscieWskaznikRzeczywisty(punkty)
			}
		}
	}
	return &obramowanie
}

// odfRodzajStylu przekłada rodzinę stylu ODF na rodzaj stylu kontraktu.
func odfRodzajStylu(rodzina string) shared.StudioStyleKind {
	switch strings.ToLower(strings.TrimSpace(rodzina)) {
	case "text":
		return shared.StudioStyleKindCharacter
	case "table":
		return shared.StudioStyleKindTable
	case "list":
		return shared.StudioStyleKindList
	default:
		return shared.StudioStyleKindParagraph
	}
}

// odfWyrownanie przekłada wyrównanie ODF na wyrównanie kontraktu.
func odfWyrownanie(wartosc string) *shared.StudioTextAlign {
	switch strings.ToLower(strings.TrimSpace(wartosc)) {
	case "start", "left":
		return wejscieWskaznikWyrownania(shared.StudioTextAlignLeft)
	case "end", "right":
		return wejscieWskaznikWyrownania(shared.StudioTextAlignRight)
	case "center":
		return wejscieWskaznikWyrownania(shared.StudioTextAlignCenter)
	case "justify":
		return wejscieWskaznikWyrownania(shared.StudioTextAlignJustify)
	}
	return nil
}

// odfPodkreslenie przekłada odmianę podkreślenia ODF na odmianę kontraktu.
func odfPodkreslenie(wartosc string) *shared.StudioUnderlineStyle {
	switch wartosc {
	case "none":
		return ooxmlPodkreslenie("none")
	case "dash", "long-dash", "dot-dash", "dot-dot-dash":
		return ooxmlPodkreslenie("dash")
	case "dotted":
		return ooxmlPodkreslenie("dotted")
	case "wave":
		return ooxmlPodkreslenie("wave")
	default:
		return ooxmlPodkreslenie("single")
	}
}

// odfWskaznikNiepusty oddaje wskaźnik na napis albo nic, gdy napis jest pusty.
func odfWskaznikNiepusty(wartosc string) *string {
	czysta := strings.TrimSpace(wartosc)
	if czysta == "" {
		return nil
	}
	return &czysta
}

// odfPierwszyNiepusty oddaje pierwszy niepusty atrybut z wykazu.
func odfPierwszyNiepusty(wezel *ooxmlWezel, nazwy ...string) string {
	for _, nazwa := range nazwy {
		if wartosc := strings.TrimSpace(wezel.atrybut(nazwa)); wartosc != "" {
			return wartosc
		}
	}
	return ""
}

// odfKopiaPostaciZnaku oddaje kopię postaci znaku — postać wspólna dla wielu
// fragmentów, podana wskaźnikiem, dałaby jedną zmianę widoczną we wszystkich.
func odfKopiaPostaciZnaku(postac *shared.StudioCharacterFormat) *shared.StudioCharacterFormat {
	if postac == nil {
		return nil
	}
	kopia := *postac
	return &kopia
}

// odfZlaczPostacZnaku nakłada postać wewnętrzną na odziedziczoną.
func odfZlaczPostacZnaku(odziedziczona,
	wewnetrzna *shared.StudioCharacterFormat) *shared.StudioCharacterFormat {

	if wewnetrzna == nil {
		return odfKopiaPostaciZnaku(odziedziczona)
	}
	if odziedziczona == nil {
		return odfKopiaPostaciZnaku(wewnetrzna)
	}
	zlaczona := *odziedziczona
	if wewnetrzna.FontFamily != nil {
		zlaczona.FontFamily = wewnetrzna.FontFamily
	}
	if wewnetrzna.FontSizePt != nil {
		zlaczona.FontSizePt = wewnetrzna.FontSizePt
	}
	if wewnetrzna.Bold != nil {
		zlaczona.Bold = wewnetrzna.Bold
	}
	if wewnetrzna.Italic != nil {
		zlaczona.Italic = wewnetrzna.Italic
	}
	if wewnetrzna.Underline != nil {
		zlaczona.Underline = wewnetrzna.Underline
	}
	if wewnetrzna.Strikethrough != nil {
		zlaczona.Strikethrough = wewnetrzna.Strikethrough
	}
	if wewnetrzna.Superscript != nil {
		zlaczona.Superscript = wewnetrzna.Superscript
	}
	if wewnetrzna.Subscript != nil {
		zlaczona.Subscript = wewnetrzna.Subscript
	}
	if wewnetrzna.Color != nil {
		zlaczona.Color = wewnetrzna.Color
	}
	if wewnetrzna.HighlightColor != nil {
		zlaczona.HighlightColor = wewnetrzna.HighlightColor
	}
	if wewnetrzna.LetterSpacingPt != nil {
		zlaczona.LetterSpacingPt = wewnetrzna.LetterSpacingPt
	}
	if wewnetrzna.SmallCaps != nil {
		zlaczona.SmallCaps = wewnetrzna.SmallCaps
	}
	if wewnetrzna.AllCaps != nil {
		zlaczona.AllCaps = wewnetrzna.AllCaps
	}
	if wewnetrzna.Language != nil {
		zlaczona.Language = wewnetrzna.Language
	}
	return &zlaczona
}

// odfWezlyWglab zbiera węzły o wskazanej nazwie z całego poddrzewa.
func odfWezlyWglab(wezel *ooxmlWezel, nazwa string) []*ooxmlWezel {
	znalezione := []*ooxmlWezel{}
	for _, dziecko := range wezel.Dzieci {
		if dziecko.Nazwa.Local == nazwa {
			znalezione = append(znalezione, dziecko)
		}
		znalezione = append(znalezione, odfWezlyWglab(dziecko, nazwa)...)
	}
	return znalezione
}

// ── Złożenie ────────────────────────────────────────────────────────────────

// odfSkladacz zbiera to, co złożenie ODF musi liczyć po drodze: treść ciała,
// style automatyczne wytworzone dla postaci bezpośredniej i wykaz cech, których
// format nie niesie.
//
// Styl automatyczny musi mieć NAZWĘ, bo tak stanowi OpenDocument — postać
// bezpośrednia w tym formacie nie istnieje. Nazwy są techniczne i Operator ich
// nie widzi; ten sam zabieg wykonuje LibreOffice przy każdym zapisie.
type odfSkladacz struct {
	cialo      strings.Builder
	styleAuto  strings.Builder
	pominiete  *[]shared.StudioSkippedItem
	licznikZ   int
	licznikA   int
	licznikTab int
}

// nazwaStyluZnaku wytwarza styl automatyczny postaci znaku i oddaje jego nazwę.
func (s *odfSkladacz) nazwaStyluZnaku(postac *shared.StudioCharacterFormat) string {
	if postac == nil {
		return ""
	}
	cechy := odfZlozCechyZnaku(postac)
	if cechy == "" {
		return ""
	}
	s.licznikZ++
	nazwa := "znak-" + strconv.Itoa(s.licznikZ)
	s.styleAuto.WriteString(`<style:style style:name="` + nazwa +
		`" style:family="text">` + cechy + `</style:style>`)
	return nazwa
}

// nazwaStyluAkapitu wytwarza styl automatyczny postaci akapitu i oddaje jego
// nazwę. Styl nadrzędny to styl NAZWANY akapitu — dzięki temu zmiana stylu
// nazwanego w edytorze Operatora dalej przestawia wszystkie akapity, które go
// używają, mimo że każdy ma własny styl automatyczny.
func (s *odfSkladacz) nazwaStyluAkapitu(postac *shared.StudioParagraphFormat) string {
	stylNazwany := odfNazwaStyluAkapitu(postac)
	if postac == nil {
		return stylNazwany
	}
	cechy := odfZlozCechyAkapitu(postac)
	if cechy == "" {
		return stylNazwany
	}
	s.licznikA++
	nazwa := "akapit-" + strconv.Itoa(s.licznikA)
	rodzic := ""
	if stylNazwany != "" {
		rodzic = ` style:parent-style-name="` + ooxmlZabezpiecz(stylNazwany) + `"`
	}
	s.styleAuto.WriteString(`<style:style style:name="` + nazwa +
		`" style:family="paragraph"` + rodzic + `>` + cechy + `</style:style>`)
	return nazwa
}

// odfNazwaStyluAkapitu oddaje nazwę stylu nazwanego akapitu — z postaci akapitu
// albo z poziomu konspektu.
func odfNazwaStyluAkapitu(postac *shared.StudioParagraphFormat) string {
	if postac == nil {
		return ""
	}
	if postac.StyleName != nil && strings.TrimSpace(*postac.StyleName) != "" {
		return strings.TrimSpace(*postac.StyleName)
	}
	if postac.OutlineLevel != nil && *postac.OutlineLevel > 0 {
		return wejscieStylNaglowkaPoziomu(*postac.OutlineLevel)
	}
	return ""
}

// wejscieZlozOdf składa `.odt` albo `.ott` z postaci dokumentu.
//
// Archiwum niesie cztery składniki: rodzaj dokumentu (`mimetype`, pierwszy
// i nieściśnięty), manifest, treść i arkusz stylów wraz ze stroną wzorcową.
// Mniej nie wystarcza — czytnik odmawia otwarcia archiwum bez manifestu,
// a nastawy strony bez strony wzorcowej nie mają na czym wisieć.
func wejscieZlozOdf(postac *shared.StudioDocumentForm, tresc, tytul string,
	szablon bool) ([]byte, []shared.StudioSkippedItem, error) {

	pominiete := make([]shared.StudioSkippedItem, 0, 4)
	bloki := postac.Blocks
	if len(bloki) == 0 {
		// Postaci nie ma — treść płaska staje akapitami, żeby plik nie wyszedł
		// pusty. To jest odtworzenie, nie odczyt, i tak wychodzi w bilansie.
		zastepcza := wejsciePostacZTekstu(postac.DocumentId, tresc)
		bloki = zastepcza.Blocks
		if postac.Styles == nil {
			postac.Styles = zastepcza.Styles
		}
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "dokument nie miał zapisanej postaci; wydanie odtworzyło akapity z treści płaskiej",
		})
	}

	skladacz := &odfSkladacz{pominiete: &pominiete}
	for _, blok := range bloki {
		switch blok.Kind {
		case wejscieRodzajBlokuTabela:
			skladacz.zlozTabele(postac, blok.TableId)
		case wejscieRodzajBlokuObiekt:
			skladacz.zlozObiekt(postac, blok.ObjectId)
		case wejscieRodzajBlokuPodzial:
			// Podział strony w ODF jest cechą stylu akapitu, nie znakiem
			// w treści — akapit pusty z rozkazem podziału jest tym, co robi
			// każdy edytor OpenDocument.
			skladacz.licznikA++
			nazwa := "akapit-podzial-" + strconv.Itoa(skladacz.licznikA)
			skladacz.styleAuto.WriteString(`<style:style style:name="` + nazwa +
				`" style:family="paragraph"><style:paragraph-properties ` +
				`fo:break-before="page"/></style:style>`)
			skladacz.cialo.WriteString(`<text:p text:style-name="` + nazwa + `"/>`)
		default:
			skladacz.zlozAkapit(blok)
		}
	}

	if len(postac.Apparatus) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "aparat dokumentu wyszedł treścią, bez pól odświeżalnych czytnika",
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

	naglowek := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	przestrzenie := ` xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0"` +
		` xmlns:style="urn:oasis:names:tc:opendocument:xmlns:style:1.0"` +
		` xmlns:text="urn:oasis:names:tc:opendocument:xmlns:text:1.0"` +
		` xmlns:table="urn:oasis:names:tc:opendocument:xmlns:table:1.0"` +
		` xmlns:draw="urn:oasis:names:tc:opendocument:xmlns:drawing:1.0"` +
		` xmlns:fo="urn:oasis:names:tc:opendocument:xmlns:xsl-fo-compatible:1.0"` +
		` xmlns:xlink="http://www.w3.org/1999/xlink"` +
		` xmlns:svg="urn:oasis:names:tc:opendocument:xmlns:svg-compatible:1.0"` +
		` xmlns:meta="urn:oasis:names:tc:opendocument:xmlns:meta:1.0"` +
		` xmlns:dc="http://purl.org/dc/elements/1.1/" office:version="1.3"`

	rodzaj := odfRodzajDokumentu
	if szablon {
		rodzaj = odfRodzajSzablonu
	}

	skladniki := map[string]string{
		odfSkladnikRodzaju: rodzaj,
		odfSkladnikManifestu: naglowek +
			`<manifest:manifest xmlns:manifest="urn:oasis:names:tc:opendocument:xmlns:manifest:1.0"` +
			` manifest:version="1.3">` +
			`<manifest:file-entry manifest:full-path="/" manifest:media-type="` + rodzaj + `"/>` +
			`<manifest:file-entry manifest:full-path="content.xml" manifest:media-type="text/xml"/>` +
			`<manifest:file-entry manifest:full-path="styles.xml" manifest:media-type="text/xml"/>` +
			`<manifest:file-entry manifest:full-path="meta.xml" manifest:media-type="text/xml"/>` +
			`</manifest:manifest>`,
		odfSkladnikTresci: naglowek + `<office:document-content` + przestrzenie + `>` +
			`<office:automatic-styles>` + skladacz.styleAuto.String() + `</office:automatic-styles>` +
			`<office:body><office:text>` + skladacz.cialo.String() +
			`</office:text></office:body></office:document-content>`,
		odfSkladnikStylow: naglowek + `<office:document-styles` + przestrzenie + `>` +
			`<office:styles>` + odfZlozArkuszStylow(postac.Styles) + `</office:styles>` +
			`<office:automatic-styles>` + odfZlozUkladStrony(postac.PageSetup) +
			`</office:automatic-styles>` +
			`<office:master-styles>` + odfZlozStroneWzorcowa(postac) +
			`</office:master-styles></office:document-styles>`,
		"meta.xml": naglowek + `<office:document-meta` + przestrzenie + `>` +
			`<office:meta><dc:title>` + ooxmlZabezpiecz(tytul) +
			`</dc:title></office:meta></office:document-meta>`,
	}

	bajty, err := wejscieZlozArchiwum(skladniki, odfSkladnikRodzaju)
	if err != nil {
		return nil, pominiete, err
	}
	return bajty, pominiete, nil
}

// zlozAkapit składa akapit wraz z fragmentami. Nagłówek wychodzi węzłem
// `text:h` z poziomem konspektu, a nie akapitem o nazwie stylu — inaczej spis
// treści w czytniku Operatora nie miałby czego zebrać.
func (s *odfSkladacz) zlozAkapit(blok shared.StudioDocumentBlock) {
	nazwaStylu := s.nazwaStyluAkapitu(blok.Paragraph)
	atrybutStylu := ""
	if nazwaStylu != "" {
		atrybutStylu = ` text:style-name="` + ooxmlZabezpiecz(nazwaStylu) + `"`
	}

	poziom := 0
	if blok.Paragraph != nil && blok.Paragraph.OutlineLevel != nil {
		poziom = *blok.Paragraph.OutlineLevel
	}
	wezel := "text:p"
	atrybutPoziomu := ""
	if poziom > 0 {
		wezel = "text:h"
		atrybutPoziomu = ` text:outline-level="` + strconv.Itoa(poziom) + `"`
	}

	s.cialo.WriteString("<" + wezel + atrybutStylu + atrybutPoziomu + ">")
	for _, fragment := range blok.Runs {
		s.zlozFragment(fragment)
	}
	s.cialo.WriteString("</" + wezel + ">")
}

// zlozFragment składa fragment tekstu wraz z postacią znaku. Odstępy wielokrotne
// i tabulatory idą węzłami ODF, bo zwykłe odstępy w XML-u zwijają się do jednego
// i Operator dostałby tekst poprzestawiany.
func (s *odfSkladacz) zlozFragment(fragment shared.StudioDocumentRun) {
	if fragment.Text == "" {
		return
	}
	nazwaStylu := s.nazwaStyluZnaku(fragment.Format)
	if nazwaStylu != "" {
		s.cialo.WriteString(`<text:span text:style-name="` +
			ooxmlZabezpiecz(nazwaStylu) + `">`)
	}
	s.cialo.WriteString(odfZlozTekst(fragment.Text))
	if nazwaStylu != "" {
		s.cialo.WriteString(`</text:span>`)
	}
}

// odfZlozTekst zapisuje tekst z zachowaniem tabulatorów, przejść do nowego
// wiersza i odstępów wielokrotnych.
func odfZlozTekst(tekst string) string {
	var budowa strings.Builder
	odstepy := 0
	domknijOdstepy := func() {
		if odstepy <= 1 {
			if odstepy == 1 {
				budowa.WriteString(" ")
			}
			odstepy = 0
			return
		}
		budowa.WriteString(` <text:s text:c="` + strconv.Itoa(odstepy-1) + `"/>`)
		odstepy = 0
	}
	for _, znak := range tekst {
		switch znak {
		case ' ':
			odstepy++
		case '\t':
			domknijOdstepy()
			budowa.WriteString(`<text:tab/>`)
		case '\n':
			domknijOdstepy()
			budowa.WriteString(`<text:line-break/>`)
		case '\r':
			// Powrót karetki bez przejścia wiersza nie jest treścią.
		default:
			domknijOdstepy()
			budowa.WriteString(ooxmlZabezpiecz(string(znak)))
		}
	}
	domknijOdstepy()
	return budowa.String()
}

// zlozTabele składa tabelę wraz z kolumnami, scaleniami i wierszem nagłówkowym.
func (s *odfSkladacz) zlozTabele(postac *shared.StudioDocumentForm, kodTabeli *string) {
	tabela := (*shared.StudioDocumentTable)(nil)
	if kodTabeli != nil {
		for i := range postac.Tables {
			if postac.Tables[i].Id == *kodTabeli {
				tabela = &postac.Tables[i]
				break
			}
		}
	}
	if tabela == nil {
		*s.pominiete = append(*s.pominiete, shared.StudioSkippedItem{
			Reason: "blok tabeli wskazuje tabelę, której postać dokumentu nie niesie",
			Detail: kodTabeli,
		})
		return
	}

	s.licznikTab++
	nazwaTabeli := "Tabela" + strconv.Itoa(s.licznikTab)
	s.cialo.WriteString(`<table:table table:name="` + nazwaTabeli + `">`)
	for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
		atrybut := ""
		if kolumna < len(tabela.ColumnWidthsMm) && tabela.ColumnWidthsMm[kolumna] > 0 {
			s.licznikA++
			nazwa := "kolumna-" + strconv.Itoa(s.licznikA)
			s.styleAuto.WriteString(`<style:style style:name="` + nazwa +
				`" style:family="table-column"><style:table-column-properties ` +
				`style:column-width="` + odfMilimetryNaZapis(tabela.ColumnWidthsMm[kolumna]) +
				`"/></style:style>`)
			atrybut = ` table:style-name="` + nazwa + `"`
		}
		s.cialo.WriteString(`<table:table-column` + atrybut + `/>`)
	}

	wierszyNaglowka := 0
	if tabela.HeaderRows != nil && *tabela.HeaderRows > 0 {
		wierszyNaglowka = *tabela.HeaderRows
	}
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		if wiersz == 0 && wierszyNaglowka > 0 {
			s.cialo.WriteString(`<table:table-header-rows>`)
		}
		s.cialo.WriteString(`<table:table-row>`)
		kolumna := 0
		for kolumna < tabela.Columns {
			komorka := ooxmlSzukajKomorki(tabela, wiersz, kolumna)
			if komorka == nil {
				s.cialo.WriteString(`<table:table-cell><text:p/></table:table-cell>`)
				kolumna++
				continue
			}
			if komorka.Merged != nil && *komorka.Merged {
				s.cialo.WriteString(`<table:covered-table-cell/>`)
				kolumna++
				continue
			}
			rozpietoscKolumn, rozpietoscWierszy := 1, 1
			if komorka.ColumnSpan != nil && *komorka.ColumnSpan > 1 {
				rozpietoscKolumn = *komorka.ColumnSpan
			}
			if komorka.RowSpan != nil && *komorka.RowSpan > 1 {
				rozpietoscWierszy = *komorka.RowSpan
			}
			atrybuty := ""
			if rozpietoscKolumn > 1 {
				atrybuty += ` table:number-columns-spanned="` +
					strconv.Itoa(rozpietoscKolumn) + `"`
			}
			if rozpietoscWierszy > 1 {
				atrybuty += ` table:number-rows-spanned="` +
					strconv.Itoa(rozpietoscWierszy) + `"`
			}
			s.cialo.WriteString(`<table:table-cell office:value-type="string"` + atrybuty + `>`)
			for _, akapit := range strings.Split(wartoscTekstu(komorka.Text), "\n") {
				s.cialo.WriteString(`<text:p>` + odfZlozTekst(akapit) + `</text:p>`)
			}
			s.cialo.WriteString(`</table:table-cell>`)
			kolumna += rozpietoscKolumn
		}
		s.cialo.WriteString(`</table:table-row>`)
		if wiersz+1 == wierszyNaglowka && wierszyNaglowka > 0 {
			s.cialo.WriteString(`</table:table-header-rows>`)
		}
	}
	s.cialo.WriteString(`</table:table>`)

	if tabela.Caption != nil && strings.TrimSpace(*tabela.Caption) != "" {
		s.cialo.WriteString(`<text:p>` + odfZlozTekst(*tabela.Caption) + `</text:p>`)
	}
}

// zlozObiekt składa akapit obiektu osadzonego.
//
// Bajty obrazu leżą w magazynie zasobów rdzenia, a wpisanie ich do archiwum
// wymagałoby odczytu magazynu, którego składacz nie ma. Obraz wychodzi więc
// akapitem z tekstem zastępczym, a strata jest NAZWANA — nie przemilczana.
func (s *odfSkladacz) zlozObiekt(postac *shared.StudioDocumentForm, kodObiektu *string) {
	obiekt := (*shared.StudioDocumentObject)(nil)
	if kodObiektu != nil {
		for i := range postac.Objects {
			if postac.Objects[i].Id == *kodObiektu {
				obiekt = &postac.Objects[i]
				break
			}
		}
	}
	if obiekt == nil {
		*s.pominiete = append(*s.pominiete, shared.StudioSkippedItem{
			Reason: "blok obiektu wskazuje obiekt, którego postać dokumentu nie niesie",
			Detail: kodObiektu,
		})
		return
	}
	opis := wartoscTekstu(obiekt.AltText)
	if strings.TrimSpace(opis) == "" {
		opis = "obiekt osadzony"
	}
	*s.pominiete = append(*s.pominiete, shared.StudioSkippedItem{
		Reason: "obraz wyszedł tekstem zastępczym, bez bajtów obrazu",
		Detail: wejscieWskaznikTekstu(opis),
	})
	s.cialo.WriteString(`<text:p>` + odfZlozTekst(opis) + `</text:p>`)
	if obiekt.Caption != nil && strings.TrimSpace(*obiekt.Caption) != "" {
		s.cialo.WriteString(`<text:p>` + odfZlozTekst(*obiekt.Caption) + `</text:p>`)
	}
}

// odfZlozCechyZnaku składa `style:text-properties` z postaci znaku.
func odfZlozCechyZnaku(postac *shared.StudioCharacterFormat) string {
	if postac == nil {
		return ""
	}
	atrybuty := ""
	if postac.FontFamily != nil && *postac.FontFamily != "" {
		atrybuty += ` style:font-name="` + ooxmlZabezpiecz(*postac.FontFamily) +
			`" fo:font-family="` + ooxmlZabezpiecz(*postac.FontFamily) + `"`
	}
	if postac.FontSizePt != nil {
		atrybuty += ` fo:font-size="` + odfPunktyNaZapis(*postac.FontSizePt) + `"`
	}
	if postac.Bold != nil {
		waga := "normal"
		if *postac.Bold {
			waga = "bold"
		}
		atrybuty += ` fo:font-weight="` + waga + `"`
	}
	if postac.Italic != nil {
		odmiana := "normal"
		if *postac.Italic {
			odmiana = "italic"
		}
		atrybuty += ` fo:font-style="` + odmiana + `"`
	}
	if postac.Underline != nil {
		odmiana := "solid"
		if *postac.Underline == shared.StudioUnderlineStyleNone {
			odmiana = "none"
		}
		atrybuty += ` style:text-underline-style="` + odmiana +
			`" style:text-underline-width="auto" style:text-underline-color="font-color"`
	}
	if postac.Strikethrough != nil {
		odmiana := "none"
		if *postac.Strikethrough {
			odmiana = "solid"
		}
		atrybuty += ` style:text-line-through-style="` + odmiana + `"`
	}
	if postac.Superscript != nil && *postac.Superscript {
		atrybuty += ` style:text-position="super 58%"`
	}
	if postac.Subscript != nil && *postac.Subscript {
		atrybuty += ` style:text-position="sub 58%"`
	}
	if postac.Color != nil && *postac.Color != "" {
		atrybuty += ` fo:color="#` + ooxmlBezKrzyzyka(*postac.Color) + `"`
	}
	if postac.HighlightColor != nil && *postac.HighlightColor != "" {
		atrybuty += ` fo:background-color="#` + ooxmlBezKrzyzyka(*postac.HighlightColor) + `"`
	}
	if postac.LetterSpacingPt != nil {
		atrybuty += ` fo:letter-spacing="` + odfPunktyNaZapis(*postac.LetterSpacingPt) + `"`
	}
	if postac.SmallCaps != nil && *postac.SmallCaps {
		atrybuty += ` fo:font-variant="small-caps"`
	}
	if postac.AllCaps != nil && *postac.AllCaps {
		atrybuty += ` fo:text-transform="uppercase"`
	}
	if atrybuty == "" {
		return ""
	}
	return `<style:text-properties` + atrybuty + `/>`
}

// odfZlozCechyAkapitu składa `style:paragraph-properties` z postaci akapitu.
func odfZlozCechyAkapitu(postac *shared.StudioParagraphFormat) string {
	if postac == nil {
		return ""
	}
	atrybuty := ""
	if postac.Align != nil {
		atrybuty += ` fo:text-align="` + odfNazwaWyrownania(*postac.Align) + `"`
	}
	if postac.IndentLeftMm != nil {
		atrybuty += ` fo:margin-left="` + odfMilimetryNaZapis(*postac.IndentLeftMm) + `"`
	}
	if postac.IndentRightMm != nil {
		atrybuty += ` fo:margin-right="` + odfMilimetryNaZapis(*postac.IndentRightMm) + `"`
	}
	if postac.FirstLineIndentMm != nil {
		atrybuty += ` fo:text-indent="` + odfMilimetryNaZapis(*postac.FirstLineIndentMm) + `"`
	}
	if postac.SpaceBeforePt != nil {
		atrybuty += ` fo:margin-top="` + odfPunktyNaZapis(*postac.SpaceBeforePt) + `"`
	}
	if postac.SpaceAfterPt != nil {
		atrybuty += ` fo:margin-bottom="` + odfPunktyNaZapis(*postac.SpaceAfterPt) + `"`
	}
	if postac.LineSpacingRule != nil && postac.LineSpacingValue != nil {
		switch *postac.LineSpacingRule {
		case shared.StudioLineSpacingRuleMultiple:
			atrybuty += ` fo:line-height="` +
				strconv.FormatFloat(*postac.LineSpacingValue*100, 'f', 0, 64) + `%"`
		default:
			atrybuty += ` fo:line-height="` + odfPunktyNaZapis(*postac.LineSpacingValue) + `"`
		}
	}
	if postac.ShadingColor != nil && *postac.ShadingColor != "" {
		atrybuty += ` fo:background-color="#` + ooxmlBezKrzyzyka(*postac.ShadingColor) + `"`
	}
	if postac.KeepWithNext != nil {
		wartosc := "auto"
		if *postac.KeepWithNext {
			wartosc = "always"
		}
		atrybuty += ` fo:keep-with-next="` + wartosc + `"`
	}
	if postac.KeepLines != nil {
		wartosc := "auto"
		if *postac.KeepLines {
			wartosc = "always"
		}
		atrybuty += ` fo:keep-together="` + wartosc + `"`
	}
	if postac.Border != nil {
		atrybuty += ` fo:border="` + odfZlozObramowanie(postac.Border) + `"`
	}

	wnetrze := ""
	if len(postac.TabStops) > 0 {
		wnetrze += `<style:tab-stops>`
		for _, tabulator := range postac.TabStops {
			rodzaj := "left"
			switch tabulator.Kind {
			case shared.StudioTabKindRight:
				rodzaj = "right"
			case shared.StudioTabKindCenter:
				rodzaj = "center"
			case shared.StudioTabKindDecimal:
				rodzaj = "char"
			}
			wiodacy := ""
			if tabulator.Leader != nil {
				switch *tabulator.Leader {
				case shared.StudioTabLeaderDot:
					wiodacy = ` style:leader-text="." style:leader-style="dotted"`
				case shared.StudioTabLeaderDash:
					wiodacy = ` style:leader-text="-" style:leader-style="solid"`
				case shared.StudioTabLeaderUnderline:
					wiodacy = ` style:leader-text="_" style:leader-style="solid"`
				}
			}
			wnetrze += `<style:tab-stop style:position="` +
				odfMilimetryNaZapis(tabulator.PositionMm) + `" style:type="` + rodzaj +
				`"` + wiodacy + `/>`
		}
		wnetrze += `</style:tab-stops>`
	}

	if atrybuty == "" && wnetrze == "" {
		return ""
	}
	if wnetrze == "" {
		return `<style:paragraph-properties` + atrybuty + `/>`
	}
	return `<style:paragraph-properties` + atrybuty + `>` + wnetrze +
		`</style:paragraph-properties>`
}

// odfZlozObramowanie składa skrót obramowania `0.50pt solid #000000`.
func odfZlozObramowanie(obramowanie *shared.StudioBorder) string {
	if obramowanie == nil {
		return "none"
	}
	if obramowanie.Style == shared.StudioBorderStyleNone {
		return "none"
	}
	szerokosc := 0.5
	if obramowanie.WidthPt != nil && *obramowanie.WidthPt > 0 {
		szerokosc = *obramowanie.WidthPt
	}
	odmiana := "solid"
	switch obramowanie.Style {
	case shared.StudioBorderStyleDouble:
		odmiana = "double"
	case shared.StudioBorderStyleDashed:
		odmiana = "dashed"
	case shared.StudioBorderStyleDotted:
		odmiana = "dotted"
	}
	barwa := "000000"
	if obramowanie.Color != nil && *obramowanie.Color != "" {
		barwa = ooxmlBezKrzyzyka(*obramowanie.Color)
	}
	return odfPunktyNaZapis(szerokosc) + " " + odmiana + " #" + barwa
}

// odfNazwaWyrownania przekłada wyrównanie kontraktu na wyrównanie ODF.
func odfNazwaWyrownania(wyrownanie shared.StudioTextAlign) string {
	switch wyrownanie {
	case shared.StudioTextAlignRight:
		return "end"
	case shared.StudioTextAlignCenter:
		return "center"
	case shared.StudioTextAlignJustify:
		return "justify"
	default:
		return "start"
	}
}

// odfZlozArkuszStylow składa arkusz stylów nazwanych.
//
// Dziedziczenie idzie `style:parent-style-name` — to jest sens stylu nadrzędnego
// i to sprawia, że zmiana tekstu zasadniczego przestawia cały dokument jednym
// ruchem także po otwarciu pliku w edytorze Operatora.
func odfZlozArkuszStylow(arkusz []shared.StudioNamedStyle) string {
	var budowa strings.Builder
	for _, styl := range arkusz {
		if strings.TrimSpace(styl.Name) == "" {
			continue
		}
		rodzina := "paragraph"
		switch styl.Kind {
		case shared.StudioStyleKindCharacter:
			rodzina = "text"
		case shared.StudioStyleKindTable:
			rodzina = "table"
		case shared.StudioStyleKindList:
			rodzina = "list"
		}
		budowa.WriteString(`<style:style style:name="` + ooxmlZabezpiecz(styl.Name) +
			`" style:family="` + rodzina + `"`)
		if styl.DisplayName != nil && *styl.DisplayName != "" {
			budowa.WriteString(` style:display-name="` + ooxmlZabezpiecz(*styl.DisplayName) + `"`)
		}
		if styl.BasedOn != nil && *styl.BasedOn != "" {
			budowa.WriteString(` style:parent-style-name="` + ooxmlZabezpiecz(*styl.BasedOn) + `"`)
		}
		if styl.NextStyle != nil && *styl.NextStyle != "" {
			budowa.WriteString(` style:next-style-name="` + ooxmlZabezpiecz(*styl.NextStyle) + `"`)
		}
		if styl.Paragraph != nil && styl.Paragraph.OutlineLevel != nil &&
			*styl.Paragraph.OutlineLevel > 0 {
			budowa.WriteString(` style:default-outline-level="` +
				strconv.Itoa(*styl.Paragraph.OutlineLevel) + `"`)
		}
		budowa.WriteString(`>`)
		budowa.WriteString(odfZlozCechyAkapitu(styl.Paragraph))
		budowa.WriteString(odfZlozCechyZnaku(styl.Character))
		budowa.WriteString(`</style:style>`)
	}
	return budowa.String()
}

// odfZlozUkladStrony składa układ strony z nastaw strony dokumentu.
func odfZlozUkladStrony(nastawy *shared.StudioPageSetup) string {
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
	kierunek := "portrait"
	if nastawy.Orientation != nil && *nastawy.Orientation == shared.StudioPageOrientationPozioma {
		kierunek = "landscape"
		if szerokosc < wysokosc {
			szerokosc, wysokosc = wysokosc, szerokosc
		}
	}
	margines := func(wskazanie *int, domyslny int) string {
		wartosc := domyslny
		if wskazanie != nil {
			wartosc = *wskazanie
		}
		return odfMilimetryNaZapis(float64(wartosc))
	}
	kolumny := ""
	if nastawy.Columns != nil && *nastawy.Columns > 1 {
		kolumny = `<style:columns fo:column-count="` + strconv.Itoa(*nastawy.Columns) +
			`" fo:column-gap="` + odfMilimetryNaZapis(wartoscRzeczywista(nastawy.ColumnGapMm)) + `"/>`
	}
	return `<style:page-layout style:name="` + odfNazwaUkladuStrony + `">` +
		`<style:page-layout-properties fo:page-width="` + odfMilimetryNaZapis(szerokosc) +
		`" fo:page-height="` + odfMilimetryNaZapis(wysokosc) +
		`" style:print-orientation="` + kierunek +
		`" fo:margin-top="` + margines(nastawy.MarginTop, 25) +
		`" fo:margin-bottom="` + margines(nastawy.MarginBottom, 25) +
		`" fo:margin-left="` + margines(nastawy.MarginLeft, 25) +
		`" fo:margin-right="` + margines(nastawy.MarginRight, 25) + `">` +
		kolumny + `</style:page-layout-properties></style:page-layout>`
}

// odfZlozStroneWzorcowa składa stronę wzorcową wraz z nagłówkiem i stopką.
//
// Nagłówek pierwszej strony i stron lewych wychodzą osobnymi węzłami — to jest
// wprost wymaganie Właściciela: nagłówek osobny dla sekcji, pierwszej strony
// i stron parzystych.
func odfZlozStroneWzorcowa(postac *shared.StudioDocumentForm) string {
	naglowki := []shared.StudioHeaderFooter{}
	if len(postac.Sections) > 0 {
		naglowki = postac.Sections[0].HeadersFooters
	}
	if len(naglowki) == 0 && postac.PageSetup != nil {
		// Nastawy strony niosą nagłówek i stopkę jednym polem dla całego
		// dokumentu — starsza droga, którą wciąż jedzie klient.
		wpis := shared.StudioHeaderFooter{Scope: shared.StudioHeaderScopeDefault}
		if postac.PageSetup.Header != nil && *postac.PageSetup.Header != "" {
			wpis.HeaderText = postac.PageSetup.Header
		}
		if postac.PageSetup.Footer != nil && *postac.PageSetup.Footer != "" {
			wpis.FooterText = postac.PageSetup.Footer
		}
		if wpis.HeaderText != nil || wpis.FooterText != nil {
			naglowki = append(naglowki, wpis)
		}
	}

	var budowa strings.Builder
	budowa.WriteString(`<style:master-page style:name="` + odfNazwaStronyWzorcowej +
		`" style:page-layout-name="` + odfNazwaUkladuStrony + `">`)
	for _, wpis := range naglowki {
		nazwaNaglowka, nazwaStopki := "style:header", "style:footer"
		switch wpis.Scope {
		case shared.StudioHeaderScopeFirstPage:
			nazwaNaglowka, nazwaStopki = "style:header-first", "style:footer-first"
		case shared.StudioHeaderScopeEvenPages:
			nazwaNaglowka, nazwaStopki = "style:header-left", "style:footer-left"
		}
		if wpis.HeaderText != nil && *wpis.HeaderText != "" {
			budowa.WriteString(`<` + nazwaNaglowka + `><text:p>` +
				odfZlozTekst(*wpis.HeaderText) + `</text:p></` + nazwaNaglowka + `>`)
		}
		if wpis.FooterText != nil && *wpis.FooterText != "" {
			budowa.WriteString(`<` + nazwaStopki + `><text:p>` +
				odfZlozTekst(*wpis.FooterText) + `</text:p></` + nazwaStopki + `>`)
		}
	}
	budowa.WriteString(`</style:master-page>`)
	return budowa.String()
}
