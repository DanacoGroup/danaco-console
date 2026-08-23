// Odpowiedzialność pliku: zamiana PDF na dokument EDYTOWALNY modułu Studio —
// odzyskanie tekstu, akapitów, układów tabelarycznych i obrazów na tyle, na ile
// PDF je niesie, wraz z BILANSEM tego, co odzyskane, a co nie.
//
// ── Dlaczego to jest odtworzenie, nie odczyt ────────────────────────────────
// PDF nie niesie struktury akapitu ani tabeli. Niesie rozkaz „postaw ten napis
// w tym miejscu strony". Akapit, wiersz tabeli i kolumna są WNIOSKAMI z układu
// napisów, a nie zapisem w pliku. Dlatego:
//
//   - wynik nazywa się odzyskaniem i idzie z bilansem;
//   - bilans mówi, ile stron miało warstwę tekstową, ile akapitów odtworzono,
//     ile układów tabelarycznych rozpoznano, a ile POZOSTAŁO nierozpoznanych;
//   - PDF ze samych skanów NIE UDAJE konwersji — wchodzi do kolejki rozpoznania
//     pisma, a odpowiedź mówi to wprost i oddaje pozycję tej kolejki.
//
// Zlecenie stanowi to wprost: „konwersja PDF → dokument edytowalny bez bilansu
// odzyskanego jest brakiem, nie skrótem".
//
// ── Czym liczone ────────────────────────────────────────────────────────────
// `pdfcpu`, wkompilowany w binarium rdzenia. Programów zewnętrznych zdjętych
// z rdzenia nie wołamy — wykaz zdjętych trzyma zapora rdzenia
// (`zapora_warsztatu_pdf_test.go`) i to ona jest jego jedynym miejscem. Arsenał
// jest wkompilowany, a proces potomny byłby tu regresem: jedno binarium serwera
// zamieniłoby się w dwa, doszedłby koszt uruchomienia i rozjazd wersji.
// Rozpoznanie pisma ze skanu jest osobną drogą (kolejka wczytywania modułu),
// bo wymaga Tesseractu, czyli składnika PAKIETU SERWERA — i tak też jest
// zgłoszone, jako brak w wykazie zależności pakietu, nie jako brak funkcji.
package core

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"danacoconsole/shared"
)

// Granice odzyskania. Chronią rdzeń przed dokumentem, którego rozbiór
// zająłby pamięć maszyny, i mówią wprost, gdzie stoi kres.
const (
	// wejscieNajwiecejStronPdf jest kresem stron branych pod rozbiór za jednym
	// razem. Dokument dłuższy wchodzi w całości, ale nadwyżka stron wychodzi
	// w wykazie pominiętych — Operator ma wiedzieć, że reszta czeka.
	wejscieNajwiecejStronPdf = 500
	// wejscieNajwiecejObrazowPdf ogranicza liczbę obrazów osadzanych w wyniku.
	wejscieNajwiecejObrazowPdf = 200
	// wejscieProgWarstwyTekstowej to najmniejsza liczba znaków, po której strona
	// uznaje się za niosącą warstwę tekstową. Jedna litera z sygnatury drukarki
	// nie czyni ze skanu dokumentu tekstowego.
	wejscieProgWarstwyTekstowej = 16
)

// wejscieOdzyskaniePdf zbiera wynik odzyskania: postać dokumentu, treść płaską,
// bilans oraz wykaz obrazów gotowych do odłożenia w magazynie zasobów.
type wejscieOdzyskaniePdf struct {
	Postac shared.StudioDocumentForm
	Tresc  string
	Bilans shared.StudioImportBalance
	// Obrazy niosą bajty obrazów wyjętych ze stron. Odłożenie ich do magazynu
	// robi wołający, bo tylko on ma kontekst żądania i okno; ten rachunek nie
	// sięga do magazynu, żeby dało się go sprawdzić bez bazy.
	Obrazy []wejscieObrazPdf
	// SamSkan mówi, że dokument nie ma warstwy tekstowej i konwersji udawać
	// nie wolno — wołający kieruje go na rozpoznanie pisma.
	SamSkan bool
}

// wejscieObrazPdf niesie obraz wyjęty ze strony PDF.
type wejscieObrazPdf struct {
	// ObiektKod wiąże obraz z obiektem postaci, do którego trafi kod zasobu po
	// odłożeniu bajtów w magazynie.
	ObiektKod   string
	Nazwa       string
	Format      string
	Bajty       []byte
	Szerokosc   int
	Wysokosc    int
	NumerStrony int
}

// wejscieCzytajPdf odzyskuje dokument edytowalny z PDF.
func wejscieCzytajPdf(kodDokumentu string, bajty []byte, zakresStron string,
	odzyskacTabele, osadzacObrazy bool) (wejscieOdzyskaniePdf, error) {

	wynik := wejscieOdzyskaniePdf{
		Bilans: shared.StudioImportBalance{Format: shared.StudioImportFormatPdf},
	}
	if len(bajty) < 5 || string(bajty[:5]) != "%PDF-" {
		return wynik, bladWskazaniaStudio(
			"wskazany plik nie jest dokumentem PDF — nie zaczyna się znacznikiem %PDF-")
	}

	nastawy := nastawyPdf()
	stron, err := api.PageCount(bytes.NewReader(bajty), nastawy)
	if err != nil {
		return wynik, bladWskazaniaStudio("dokumentu PDF nie da się odczytać: " + err.Error())
	}
	wynik.Bilans.Pages = wejscieWskaznikCalkowity(stron)

	var wybraneStrony []string
	if czysty := strings.TrimSpace(zakresStron); czysty != "" {
		wybraneStrony = []string{czysty}
	}

	// Odczyt strumieni treści stron. Każda strona daje swoje wiersze — akapity
	// składa się z wierszy, a nie ze stron, bo akapit potrafi przechodzić przez
	// granicę strony.
	wierszeStron := make([][]string, 0, stron)
	err = api.ExtractContent(bytes.NewReader(bajty), wybraneStrony,
		func(strumien io.Reader, _ int) error {
			tresc, err := io.ReadAll(strumien)
			if err != nil {
				return err
			}
			wierszeStron = append(wierszeStron, wejscieWierszeStronyPdf(string(tresc)))
			return nil
		}, nastawy)
	if err != nil {
		return wynik, bladWskazaniaStudio(
			"nie można sięgnąć po treść dokumentu PDF: " + err.Error())
	}

	zeTekstem, bezTekstu := 0, 0
	for _, wiersze := range wierszeStron {
		if wejscieDlugoscWierszy(wiersze) >= wejscieProgWarstwyTekstowej {
			zeTekstem++
			continue
		}
		bezTekstu++
	}
	wynik.Bilans.PagesWithText = wejscieWskaznikCalkowity(zeTekstem)
	wynik.Bilans.PagesWithoutText = wejscieWskaznikCalkowity(bezTekstu)

	if zeTekstem == 0 {
		// PDF ze samych skanów. Konwersji nie udajemy: dokument bez warstwy
		// tekstowej po „konwersji" byłby pustą kartką, o której Operator
		// pomyślałby, że taki jest jego plik.
		wynik.SamSkan = true
		wynik.Bilans.NeedsTextRecognition = wejscieWskaznikLogiczny(true)
		wynik.Bilans.Note = wejscieWskaznikTekstu("dokument PDF nie ma warstwy tekstowej na " +
			"żadnej z " + strconv.Itoa(stron) + " stron — to jest skan. Odzyskania tekstu " +
			"nie udaję: plik idzie na rozpoznanie pisma, a dokument edytowalny powstanie " +
			"z jego wyniku")
		return wynik, nil
	}
	if bezTekstu > 0 {
		wynik.Bilans.Skipped = append(wynik.Bilans.Skipped, shared.StudioSkippedItem{
			Reason: "strony bez warstwy tekstowej — ich treść nie weszła do dokumentu",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(bezTekstu) + " z " +
				strconv.Itoa(stron) + " stron; naprawa: skierować je na rozpoznanie pisma"),
		})
		wynik.Bilans.NeedsTextRecognition = wejscieWskaznikLogiczny(true)
	}
	if stron > wejscieNajwiecejStronPdf {
		wynik.Bilans.Skipped = append(wynik.Bilans.Skipped, shared.StudioSkippedItem{
			Reason: "dokument jest dłuższy niż granica jednego odzyskania",
			Detail: wejscieWskaznikTekstu("granica to " +
				strconv.Itoa(wejscieNajwiecejStronPdf) + " stron; naprawa: wnosić zakresami " +
				"stron polem żądania"),
		})
	}

	wszystkieWiersze := make([]string, 0, 256)
	for _, wiersze := range wierszeStron {
		wszystkieWiersze = append(wszystkieWiersze, wiersze...)
		// Granica strony jest wierszem pustym: bez niej ostatni akapit strony
		// zlałby się z pierwszym akapitem strony następnej.
		wszystkieWiersze = append(wszystkieWiersze, "")
	}

	postac, rozpoznaneTabele, nierozpoznane := wejsciePostacZWierszyPdf(
		kodDokumentu, wszystkieWiersze, odzyskacTabele)
	wynik.Postac = postac
	wynik.Bilans.ParagraphsRecovered = wejscieWskaznikCalkowity(len(postac.Blocks))
	wynik.Bilans.TablesRecognized = wejscieWskaznikCalkowity(rozpoznaneTabele)
	wynik.Bilans.TablesMissed = wejscieWskaznikCalkowity(nierozpoznane)
	wynik.Bilans.StylesRecovered = wejscieWskaznikCalkowity(len(postac.Styles))
	wynik.Bilans.SectionsRecovered = wejscieWskaznikCalkowity(len(postac.Sections))

	if nierozpoznane > 0 {
		wynik.Bilans.Skipped = append(wynik.Bilans.Skipped, shared.StudioSkippedItem{
			Reason: "układy tabelaryczne, których nie dało się rozpoznać jako tabel",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(nierozpoznane) +
				" — weszły jako akapity z zachowanym rozkładem odstępów"),
		})
	}
	wynik.Bilans.Skipped = append(wynik.Bilans.Skipped, shared.StudioSkippedItem{
		Reason: "postać znaku i akapitu nie weszła — PDF nie niesie stylów, tylko rozkazy " +
			"rysowania napisów",
		Detail: wejscieWskaznikTekstu("dokument dostał arkusz stylów domyślny platformy; " +
			"nagłówki rozpoznane po układzie wiersza, nie po stylu"),
	})

	if osadzacObrazy {
		obrazy, pominieteObrazy := wejscieObrazyPdf(bajty, wybraneStrony, nastawy)
		wynik.Obrazy = obrazy
		wynik.Bilans.ImagesEmbedded = wejscieWskaznikCalkowity(len(obrazy))
		wynik.Bilans.ImagesSkipped = wejscieWskaznikCalkowity(pominieteObrazy)
		if pominieteObrazy > 0 {
			wynik.Bilans.Skipped = append(wynik.Bilans.Skipped, shared.StudioSkippedItem{
				Reason: "obrazy, których nie dało się wyjąć ze stron",
				Detail: wejscieWskaznikTekstu(strconv.Itoa(pominieteObrazy) +
					" — zapis obrazu jest odmianą, której rdzeń nie rozpakowuje"),
			})
		}
		wejscieDolozObrazyDoPostaci(&wynik)
	} else {
		wynik.Bilans.Skipped = append(wynik.Bilans.Skipped, shared.StudioSkippedItem{
			Reason: "obrazy pominięte na wyraźne życzenie żądania",
		})
	}

	wynik.Tresc = wejscieTrescZPostaci(&wynik.Postac)
	wejscieUzycieStylow(&wynik.Postac)
	wynik.Bilans.Note = wejscieWskaznikTekstu(wejscieZdanieOdzyskaniaPdf(&wynik.Bilans))
	return wynik, nil
}

// wejscieZdanieOdzyskaniaPdf składa zdanie o uczciwym stanie wyniku.
func wejscieZdanieOdzyskaniaPdf(bilans *shared.StudioImportBalance) string {
	return "odzyskanie z PDF jest ODTWORZENIEM, nie odczytem: " +
		strconv.Itoa(wartoscCalkowita(bilans.PagesWithText)) + " z " +
		strconv.Itoa(wartoscCalkowita(bilans.Pages)) + " stron miało warstwę tekstową, " +
		"odtworzono " + strconv.Itoa(wartoscCalkowita(bilans.ParagraphsRecovered)) +
		" akapitów, rozpoznano " + strconv.Itoa(wartoscCalkowita(bilans.TablesRecognized)) +
		" tabel przy " + strconv.Itoa(wartoscCalkowita(bilans.TablesMissed)) +
		" układach nierozpoznanych, osadzono " +
		strconv.Itoa(wartoscCalkowita(bilans.ImagesEmbedded)) + " obrazów"
}

// wejscieDlugoscWierszy liczy znaki treści strony.
func wejscieDlugoscWierszy(wiersze []string) int {
	dlugosc := 0
	for _, wiersz := range wiersze {
		dlugosc += len([]rune(strings.TrimSpace(wiersz)))
	}
	return dlugosc
}

// ── Wiersze ze strumienia treści strony ─────────────────────────────────────

// wejscieWierszeStronyPdf wyjmuje wiersze tekstu ze strumienia treści strony.
//
// Rozbiór idzie po operatorach zapisu napisu (`Tj`, `TJ`, `'`, `"`) i po
// operatorach przesunięcia wiersza (`Td`, `TD`, `T*`, `Tm`): przesunięcie
// zaczyna wiersz nowy. To jest właśnie ten wniosek z układu, o którym mówi
// nagłówek pliku — PDF nie mówi „nowy wiersz", mówi „przesuń kursor".
//
// Napisy sześciowartościowe (`<0041>` w miejsce `(A)`) też wchodzą: pliki
// składane z kroju osadzonego zapisują tekst właśnie tak, a pominięcie ich
// dałoby dokument pusty przy pliku, który tekst niesie.
func wejscieWierszeStronyPdf(strumien string) []string {
	wiersze := make([]string, 0, 64)
	var wiersz strings.Builder
	var napis strings.Builder
	var rozkaz strings.Builder

	wNapisie, wSzesnastkowym := false, false
	poziom := 0

	domknijWiersz := func() {
		tekst := strings.TrimRight(wiersz.String(), " ")
		wiersz.Reset()
		if strings.TrimSpace(tekst) == "" {
			// Wiersz pusty zostaje pusty, ale nie mnoży się: dwa puste pod rząd
			// to jeden odstęp akapitowy.
			if len(wiersze) > 0 && wiersze[len(wiersze)-1] != "" {
				wiersze = append(wiersze, "")
			}
			return
		}
		wiersze = append(wiersze, tekst)
	}

	for i := 0; i < len(strumien); i++ {
		znak := strumien[i]

		if wNapisie {
			switch znak {
			case '\\':
				if i+1 < len(strumien) {
					i++
					switch strumien[i] {
					case 'n':
						napis.WriteByte('\n')
					case 't':
						napis.WriteByte('\t')
					case 'r':
					default:
						napis.WriteByte(strumien[i])
					}
				}
			case '(':
				poziom++
				napis.WriteByte(znak)
			case ')':
				if poziom > 0 {
					poziom--
					napis.WriteByte(znak)
					continue
				}
				wNapisie = false
				wiersz.WriteString(napis.String())
				napis.Reset()
			default:
				napis.WriteByte(znak)
			}
			continue
		}

		if wSzesnastkowym {
			if znak == '>' {
				wSzesnastkowym = false
				wiersz.WriteString(wejscieNapisSzesnastkowyPdf(napis.String()))
				napis.Reset()
				continue
			}
			napis.WriteByte(znak)
			continue
		}

		switch {
		case znak == '(':
			wNapisie = true
			rozkaz.Reset()
		case znak == '<' && i+1 < len(strumien) && strumien[i+1] != '<':
			wSzesnastkowym = true
			rozkaz.Reset()
		case znak == ' ' || znak == '\n' || znak == '\r' || znak == '\t' ||
			znak == '[' || znak == ']':
			if wejscieRozkazNowegoWierszaPdf(rozkaz.String()) {
				domknijWiersz()
			}
			rozkaz.Reset()
		default:
			rozkaz.WriteByte(znak)
			if rozkaz.Len() > 8 {
				rozkaz.Reset()
			}
		}
	}
	if wejscieRozkazNowegoWierszaPdf(rozkaz.String()) {
		domknijWiersz()
	}
	domknijWiersz()

	// Wiersze pojedynczych napisów potrafią przyjść w kawałkach — łączymy je
	// bez rozdzielania słów, bo `TJ` rozbija wyraz na kilka napisów po to, żeby
	// dosunąć odstępy między literami.
	sprzatniete := make([]string, 0, len(wiersze))
	for _, wiersz := range wiersze {
		sprzatniete = append(sprzatniete, strings.TrimRight(wiersz, " \t"))
	}
	return sprzatniete
}

// wejscieRozkazNowegoWierszaPdf mówi, czy rozkaz przesuwa kursor do nowego
// wiersza. `Td`, `TD` i `Tm` przesuwają punkt wpisywania, `T*` przechodzi
// o wiersz niżej, a `ET` domyka blok tekstu.
func wejscieRozkazNowegoWierszaPdf(rozkaz string) bool {
	switch strings.TrimSpace(rozkaz) {
	case "Td", "TD", "T*", "Tm", "ET", "'", "\"":
		return true
	}
	return false
}

// wejscieNapisSzesnastkowyPdf przekłada napis zapisany szesnastkowo na tekst.
//
// Zapis dwubajtowy (`<00410042>`) jest w plikach z krojem osadzonym normą.
// Bajt wyższy zerowy znaczy zwykły znak ASCII; wartość wyższą przepuszczamy
// jako punkt kodowy, bo to najlepsze, co da się zrobić bez mapy kroju.
func wejscieNapisSzesnastkowyPdf(zapis string) string {
	czysty := strings.Map(func(znak rune) rune {
		switch {
		case znak >= '0' && znak <= '9', znak >= 'a' && znak <= 'f', znak >= 'A' && znak <= 'F':
			return znak
		}
		return -1
	}, zapis)
	if len(czysty)%2 == 1 {
		czysty += "0"
	}
	var budowa strings.Builder
	if len(czysty)%4 == 0 && len(czysty) >= 4 {
		for i := 0; i+4 <= len(czysty); i += 4 {
			wartosc, err := strconv.ParseUint(czysty[i:i+4], 16, 32)
			if err != nil {
				continue
			}
			if wartosc == 0 {
				continue
			}
			budowa.WriteRune(rune(wartosc))
		}
		return budowa.String()
	}
	for i := 0; i+2 <= len(czysty); i += 2 {
		wartosc, err := strconv.ParseUint(czysty[i:i+2], 16, 16)
		if err != nil || wartosc == 0 {
			continue
		}
		budowa.WriteRune(rune(wartosc))
	}
	return budowa.String()
}

// ── Odtworzenie postaci z wierszy ───────────────────────────────────────────

// wejsciePostacZWierszyPdf składa postać dokumentu z wierszy odzyskanych ze
// stron: akapity, nagłówki i tabele.
//
// Oddaje też liczbę tabel rozpoznanych i liczbę układów tabelarycznych, których
// rozpoznać się nie dało — bilans stoi na tych dwóch liczbach.
func wejsciePostacZWierszyPdf(kodDokumentu string, wiersze []string,
	odzyskacTabele bool) (shared.StudioDocumentForm, int, int) {

	postac := shared.StudioDocumentForm{DocumentId: kodDokumentu}
	nastawy := wejscieDomyslneNastawyStrony("A4", nil)
	postac.PageSetup = &nastawy
	postac.Styles = wejscieDomyslnyArkuszStylow()

	sekcja := shared.StudioSection{
		Id:        nowyIdentyfikator(przedrostekSekcjiStudia),
		Index:     0,
		Start:     wejscieWskaznikPoczatkuSekcji(shared.StudioSectionStartContinuous),
		PageSetup: &nastawy,
	}
	postac.Sections = []shared.StudioSection{sekcja}

	rozpoznane, nierozpoznane := 0, 0
	akapit := []string{}

	domknijAkapit := func() {
		if len(akapit) == 0 {
			return
		}
		tresc := strings.Join(akapit, " ")
		akapit = akapit[:0]
		tresc = strings.Join(strings.Fields(tresc), " ")
		if tresc == "" {
			return
		}
		styl := wejscieStylTekstZasadniczy
		poziom := 0
		if poziomNaglowka := wejsciePoziomNaglowkaPdf(tresc); poziomNaglowka > 0 {
			poziom = poziomNaglowka
			styl = wejscieStylNaglowkaPoziomu(poziomNaglowka)
		}
		blok := wejscieBlokAkapituPdf(sekcja.Id, tresc, styl, poziom)
		postac.Blocks = append(postac.Blocks, blok)
	}

	for i := 0; i < len(wiersze); i++ {
		wiersz := wiersze[i]
		if strings.TrimSpace(wiersz) == "" {
			domknijAkapit()
			continue
		}

		if wejscieWierszTabelarycznyPdf(wiersz) {
			blok := wejscieZbierzBlokTabelaryczny(wiersze, i)
			if len(blok) >= 2 {
				domknijAkapit()
				if odzyskacTabele {
					if tabela, jest := wejscieTabelaZBlokuPdf(blok); jest {
						postac.Tables = append(postac.Tables, tabela)
						postac.Blocks = append(postac.Blocks, shared.StudioDocumentBlock{
							Id:        nowyIdentyfikator(przedrostekBlokuStudia),
							Kind:      wejscieRodzajBlokuTabela,
							SectionId: wejscieWskaznikTekstu(sekcja.Id),
							TableId:   wejscieWskaznikTekstu(tabela.Id),
						})
						rozpoznane++
						i += len(blok) - 1
						continue
					}
				}
				// Układu nie dało się rozpoznać jako tabeli (albo Operator
				// odzyskiwania tabel nie chciał): wiersze wchodzą akapitami
				// z zachowanym rozkładem odstępów, a bilans to liczy.
				nierozpoznane++
				for _, wierszUkladu := range blok {
					postac.Blocks = append(postac.Blocks, wejscieBlokAkapituPdf(
						sekcja.Id, wierszUkladu, wejscieStylTekstZasadniczy, 0))
				}
				i += len(blok) - 1
				continue
			}
		}

		akapit = append(akapit, strings.TrimSpace(wiersz))
		// Wiersz kończący zdanie domyka akapit: PDF nie mówi, gdzie akapit się
		// kończy, a zdanie zamknięte kropką i krótszy wiersz są najmocniejszą
		// przesłanką, jaką układ napisów daje.
		if wejscieWierszZamykaAkapitPdf(wiersz, wiersze, i) {
			domknijAkapit()
		}
	}
	domknijAkapit()

	if len(postac.Blocks) == 0 {
		postac.Blocks = append(postac.Blocks, wejscieBlokAkapituPdf(
			sekcja.Id, "", wejscieStylTekstZasadniczy, 0))
	}
	wejsciePrzeliczSekcje(&postac)
	return postac, rozpoznane, nierozpoznane
}

// wejscieBlokAkapituPdf składa blok akapitu odzyskanego z PDF wraz z poziomem
// konspektu. Poziom jest tu istotny, nie ozdobny: bez niego spis treści po
// wniesieniu nie miałby czego zebrać, choć nagłówki w dokumencie stoją.
func wejscieBlokAkapituPdf(kodSekcji, tresc, styl string,
	poziom int) shared.StudioDocumentBlock {

	blok := wejscieBlokAkapitu(kodSekcji, tresc, styl, nil)
	if poziom > 0 && blok.Paragraph != nil {
		blok.Paragraph.OutlineLevel = wejscieWskaznikCalkowity(poziom)
	}
	return blok
}

// wejscieWierszZamykaAkapitPdf rozstrzyga, czy wiersz domyka akapit.
func wejscieWierszZamykaAkapitPdf(wiersz string, wiersze []string, wskazanie int) bool {
	czysty := strings.TrimSpace(wiersz)
	if czysty == "" {
		return true
	}
	ostatni := czysty[len(czysty)-1]
	if ostatni != '.' && ostatni != '!' && ostatni != '?' && ostatni != ':' {
		return false
	}
	if wskazanie+1 >= len(wiersze) {
		return true
	}
	nastepny := strings.TrimSpace(wiersze[wskazanie+1])
	if nastepny == "" {
		return true
	}
	// Zdanie zamknięte, a wiersz następny zaczyna się wielką literą albo
	// cyfrą — akapit nowy. Wiersz następny małą literą znaczy zdanie łamane
	// przez skrót („art. 5 ust. 1"), a nie akapit nowy.
	pierwszy := []rune(nastepny)[0]
	return pierwszy >= 'A' && pierwszy <= 'Z' || pierwszy >= '0' && pierwszy <= '9' ||
		strings.ContainsRune("ĄĆĘŁŃÓŚŹŻ", pierwszy)
}

// wejsciePoziomNaglowkaPdf rozpoznaje nagłówek po układzie wiersza: krótki,
// bez kropki na końcu, pisany wersalikami albo poprzedzony numeracją własną
// dokumentu Operatora.
//
// To jest odtworzenie po układzie, nie odczyt stylu — i tak wychodzi w bilansie.
// Numeracja, o którą tu chodzi, jest numeracją PISMA OPERATORA (rozdział 1,
// paragraf 3), a nie kodem wymyślonym przez rdzeń.
func wejsciePoziomNaglowkaPdf(tresc string) int {
	czysta := strings.TrimSpace(tresc)
	if czysta == "" || len([]rune(czysta)) > 120 {
		return 0
	}
	ostatni := czysta[len(czysta)-1]
	if ostatni == '.' || ostatni == ',' || ostatni == ';' {
		return 0
	}
	bezZnakow := strings.Map(func(znak rune) rune {
		if znak >= 'a' && znak <= 'z' || strings.ContainsRune("ąćęłńóśźż", znak) {
			return znak
		}
		return -1
	}, czysta)
	if bezZnakow == "" && len([]rune(czysta)) > 3 {
		// Sam zapis wersalikami — nagłówek poziomu pierwszego.
		return 1
	}
	slowa := strings.Fields(czysta)
	if len(slowa) == 0 {
		return 0
	}
	pierwsze := strings.TrimRight(slowa[0], ".")
	kropki := strings.Count(pierwsze, ".")
	if _, err := strconv.Atoi(strings.ReplaceAll(pierwsze, ".", "")); err == nil && kropki <= 5 {
		poziom := kropki + 1
		if poziom > 6 {
			poziom = 6
		}
		return poziom
	}
	return 0
}

// ── Układy tabelaryczne ─────────────────────────────────────────────────────

// wejscieWierszTabelarycznyPdf mówi, czy wiersz wygląda na wiersz tabeli:
// niesie co najmniej dwa odstępy szerokie, którymi PDF rozdziela kolumny.
func wejscieWierszTabelarycznyPdf(wiersz string) bool {
	return len(wejscieKolumnyWierszaPdf(wiersz)) >= 2
}

// wejscieKolumnyWierszaPdf rozdziela wiersz na kolumny po odstępach szerokich.
func wejscieKolumnyWierszaPdf(wiersz string) []string {
	czysty := strings.ReplaceAll(wiersz, "\t", "   ")
	czlony := strings.Split(czysty, "  ")
	kolumny := make([]string, 0, len(czlony))
	for _, czlon := range czlony {
		if strings.TrimSpace(czlon) == "" {
			continue
		}
		kolumny = append(kolumny, strings.TrimSpace(czlon))
	}
	if len(kolumny) < 2 {
		return nil
	}
	return kolumny
}

// wejscieZbierzBlokTabelaryczny zbiera ciąg wierszy tabelarycznych.
func wejscieZbierzBlokTabelaryczny(wiersze []string, od int) []string {
	blok := []string{}
	for i := od; i < len(wiersze); i++ {
		if !wejscieWierszTabelarycznyPdf(wiersze[i]) {
			break
		}
		blok = append(blok, wiersze[i])
	}
	return blok
}

// wejscieTabelaZBlokuPdf składa tabelę z bloku wierszy tabelarycznych.
//
// Tabela powstaje wtedy i tylko wtedy, gdy liczba kolumn jest ZGODNA w całym
// bloku albo różni się o jedną. Blok o rozjeżdżającej się liczbie kolumn nie
// jest tabelą, tylko tekstem w kolumnach — i wtedy rachunek go NIE UDAJE, tylko
// oddaje fałsz, a wołający liczy go jako układ nierozpoznany.
func wejscieTabelaZBlokuPdf(blok []string) (shared.StudioDocumentTable, bool) {
	wiersze := make([][]string, 0, len(blok))
	najwiecej := 0
	najmniej := 0
	for _, wiersz := range blok {
		kolumny := wejscieKolumnyWierszaPdf(wiersz)
		if len(kolumny) == 0 {
			return shared.StudioDocumentTable{}, false
		}
		wiersze = append(wiersze, kolumny)
		if len(kolumny) > najwiecej {
			najwiecej = len(kolumny)
		}
		if najmniej == 0 || len(kolumny) < najmniej {
			najmniej = len(kolumny)
		}
	}
	if najwiecej-najmniej > 1 || najwiecej < 2 {
		return shared.StudioDocumentTable{}, false
	}

	tabela := shared.StudioDocumentTable{
		Id:      nowyIdentyfikator(przedrostekTabeliStudia),
		Rows:    len(wiersze),
		Columns: najwiecej,
		// Wiersz pierwszy bierze się za nagłówkowy: w piśmie urzędowym tabela
		// bez nagłówka jest rzadkością, a odzyskanie mówi wprost, że jest
		// odtworzeniem.
		HeaderRows:   wejscieWskaznikCalkowity(1),
		RepeatHeader: wejscieWskaznikLogiczny(true),
	}
	// Szerokości kolumn liczone udziałem najdłuższej treści w kolumnie: tabela
	// z szerokościami zerowymi jest usterką, którą sprawdzian odcinka postaci
	// mierzy wprost.
	najdluzsze := make([]float64, najwiecej)
	for _, kolumny := range wiersze {
		for numer, tresc := range kolumny {
			dlugosc := float64(len([]rune(tresc)))
			if dlugosc > najdluzsze[numer] {
				najdluzsze[numer] = dlugosc
			}
		}
	}
	suma := 0.0
	for _, dlugosc := range najdluzsze {
		suma += dlugosc
	}
	// Obszar pisania A4 przy marginesach 25 mm.
	const szerokoscObszaru = 160.0
	if suma > 0 {
		tabela.ColumnWidthsMm = make([]float64, najwiecej)
		for numer, dlugosc := range najdluzsze {
			tabela.ColumnWidthsMm[numer] = szerokoscObszaru * dlugosc / suma
		}
		tabela.WidthMm = wejscieWskaznikRzeczywisty(szerokoscObszaru)
	}

	for numerWiersza, kolumny := range wiersze {
		for numerKolumny := 0; numerKolumny < najwiecej; numerKolumny++ {
			tresc := ""
			if numerKolumny < len(kolumny) {
				tresc = kolumny[numerKolumny]
			}
			tabela.Cells = append(tabela.Cells, shared.StudioTableCell{
				Row: numerWiersza, Column: numerKolumny,
				Text: wejscieWskaznikTekstu(tresc),
			})
		}
	}
	return tabela, true
}

// ── Obrazy ──────────────────────────────────────────────────────────────────

// wejscieObrazyPdf wyjmuje obrazy osadzone w stronach dokumentu.
func wejscieObrazyPdf(bajty []byte, wybraneStrony []string,
	nastawy *model.Configuration) ([]wejscieObrazPdf, int) {

	obrazy := make([]wejscieObrazPdf, 0, 8)
	pominiete := 0

	err := api.ExtractImages(bytes.NewReader(bajty), wybraneStrony,
		func(obraz model.Image, _ bool, numerStrony int) error {
			if len(obrazy) >= wejscieNajwiecejObrazowPdf {
				pominiete++
				return nil
			}
			if obraz.Reader == nil {
				pominiete++
				return nil
			}
			tresc, err := io.ReadAll(io.LimitReader(obraz.Reader, granicaSkladnikaArchiwum))
			if err != nil || len(tresc) == 0 {
				pominiete++
				return nil
			}
			format := strings.TrimPrefix(strings.ToLower(obraz.FileType), ".")
			if format == "" {
				format = "png"
			}
			obrazy = append(obrazy, wejscieObrazPdf{
				ObiektKod:   nowyIdentyfikator(przedrostekObiektuStudia),
				Nazwa:       obraz.Name,
				Format:      format,
				Bajty:       tresc,
				Szerokosc:   obraz.Width,
				Wysokosc:    obraz.Height,
				NumerStrony: numerStrony,
			})
			return nil
		}, nastawy)
	if err != nil {
		// Odmowa biblioteki na obrazach NIE przewraca odzyskania tekstu:
		// dokument z tekstem, a bez obrazów, jest wynikiem gorszym, ale
		// prawdziwym — i bilans mówi, ile obrazów odpadło.
		pominiete++
	}
	return obrazy, pominiete
}

// wejscieDolozObrazyDoPostaci dokłada obiekty obrazów do postaci dokumentu wraz
// z blokami, którymi stoją w treści.
func wejscieDolozObrazyDoPostaci(wynik *wejscieOdzyskaniePdf) {
	if len(wynik.Obrazy) == 0 {
		return
	}
	sekcja := ""
	if len(wynik.Postac.Sections) > 0 {
		sekcja = wynik.Postac.Sections[0].Id
	}
	for _, obraz := range wynik.Obrazy {
		opis := strings.TrimSpace(obraz.Nazwa)
		if opis == "" {
			opis = "obraz ze strony " + strconv.Itoa(obraz.NumerStrony)
		}
		obiekt := shared.StudioDocumentObject{
			Id:      obraz.ObiektKod,
			Kind:    shared.StudioObjectKindImage,
			Source:  wejscieWskaznikZrodlaObiektu(shared.StudioObjectSourceFile),
			AltText: wejscieWskaznikTekstu(opis),
		}
		if obraz.Szerokosc > 0 && obraz.Wysokosc > 0 {
			// Wymiar w milimetrach przy 96 punktach na cal — tyle, ile PDF
			// niesie bez macierzy przekształcenia strony.
			obiekt.WidthMm = wejscieWskaznikRzeczywisty(float64(obraz.Szerokosc) / 96 * 25.4)
			obiekt.HeightMm = wejscieWskaznikRzeczywisty(float64(obraz.Wysokosc) / 96 * 25.4)
		}
		wynik.Postac.Objects = append(wynik.Postac.Objects, obiekt)
		wynik.Postac.Blocks = append(wynik.Postac.Blocks, shared.StudioDocumentBlock{
			Id:        nowyIdentyfikator(przedrostekBlokuStudia),
			Kind:      wejscieRodzajBlokuObiekt,
			SectionId: wejscieWskaznikTekstu(sekcja),
			ObjectId:  wejscieWskaznikTekstu(obiekt.Id),
		})
	}
}
