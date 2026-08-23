// Odpowiedzialność pliku: wniesienie do edytora plików, które nie są archiwum
// biurowym — tekstu czystego, markdown, RTF i HTML — wraz z ROZPOZNANIEM ZAPISU
// ZNAKÓW, oraz wykaz nośników, na którym stoją nastawy strony tego odcinka.
//
// ── Dlaczego rozpoznanie zapisu znaków jest tu osobną pracą ────────────────
// Pliki Operatora bywają starsze niż UTF-8. Pismo urzędowe z lat, w których
// pisano je w Windows-1250, wczytane jako UTF-8 daje albo błąd, albo tekst
// z krzaczkami w miejscu „ą", „ę", „ł". Krzaczki są tu gorsze niż odmowa: model
// przeczyta je jako słowa i zacznie na nich pracować. Dlatego zapis znaków
// rozpoznaje się PRZED rozbiorem treści, a rozpoznanie wychodzi kontraktem
// w bilansie — Operator ma wiedzieć, jak rdzeń odczytał jego plik.
//
// Rozpoznanie idzie `golang.org/x/net/html/charset` i `golang.org/x/text/encoding`
// — obie biblioteki z rodziny wzorcowej Go, obie już w drzewie zależności.
// Kolejność rozstrzygania: znacznik kolejności bajtów, potem deklaracja
// w treści (nagłówek XML, `meta charset` HTML), potem sprawdzenie poprawności
// UTF-8, a na końcu miara rozkładu bajtów rozstrzygająca między stronami
// kodowymi używanymi w polskich dokumentach.
package core

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"danacoconsole/shared"
)

// Nazwy zapisów znaków, które ten rachunek rozpoznaje i wymienia w bilansie.
const (
	wejscieZapisUtf8         = "UTF-8"
	wejscieZapisUtf16LE      = "UTF-16LE"
	wejscieZapisUtf16BE      = "UTF-16BE"
	wejscieZapisWindows1250  = "windows-1250"
	wejscieZapisIso88592     = "ISO-8859-2"
	wejscieZapisWindows1252  = "windows-1252"
	wejscieZapisNierozpoznny = "nierozpoznany"
)

// wejscieRozstrzygaczZnakowXml podaje rozbiórowi XML strumień przełożony na
// UTF-8. Bez niego `encoding/xml` odmawia na deklaracji `encoding="windows-1250"`,
// choć plik jest poprawny — a takie pliki w archiwach biurowych Operatora są.
func wejscieRozstrzygaczZnakowXml(nazwa string, strumien io.Reader) (io.Reader, error) {
	return charset.NewReaderLabel(nazwa, strumien)
}

// wejscieRozpoznajZapisZnakow przekłada bajty pliku na napis UTF-8 i oddaje
// nazwę rozpoznanego zapisu.
//
// Wskazanie Operatora ma pierwszeństwo nad rozpoznaniem: kto wie, w czym jest
// jego plik, wie to lepiej niż miara rozkładu bajtów. Wskazanie nieznanej nazwy
// jest jednak odmową, nie cichym zejściem na UTF-8 — inaczej Operator dostałby
// krzaczki i myślał, że jego wskazanie zadziałało.
func wejscieRozpoznajZapisZnakow(bajty []byte, wskazanie *string) (string, string, error) {
	if wskazanie != nil && strings.TrimSpace(*wskazanie) != "" {
		nazwa := strings.TrimSpace(*wskazanie)
		przeklad, _ := charset.Lookup(nazwa)
		if przeklad == nil {
			return "", "", bladWskazaniaStudio("zapis znaków " + nazwa +
				" nie jest znany rdzeniowi; naprawa: podać nazwę z rejestru IANA " +
				"(na przykład UTF-8, windows-1250, ISO-8859-2) albo nie podawać " +
				"jej wcale i zdać się na rozpoznanie")
		}
		tekst, err := wejscieZastosujZapis(bajty, przeklad)
		if err != nil {
			return "", "", err
		}
		return tekst, nazwa, nil
	}

	// Znacznik kolejności bajtów jest rozstrzygający i nie wymaga miary.
	switch {
	case bytes.HasPrefix(bajty, []byte{0xEF, 0xBB, 0xBF}):
		return string(bajty[3:]), wejscieZapisUtf8, nil
	case bytes.HasPrefix(bajty, []byte{0xFF, 0xFE}):
		tekst, err := wejscieZastosujZapis(bajty, wejscieZapisPoNazwie(wejscieZapisUtf16LE))
		return tekst, wejscieZapisUtf16LE, err
	case bytes.HasPrefix(bajty, []byte{0xFE, 0xFF}):
		tekst, err := wejscieZastosujZapis(bajty, wejscieZapisPoNazwie(wejscieZapisUtf16BE))
		return tekst, wejscieZapisUtf16BE, err
	}

	// Deklaracja w treści — nagłówek XML albo `meta charset` HTML. Biblioteka
	// czyta jedno i drugie; wynik bierzemy tylko wtedy, gdy nazwał coś innego
	// niż UTF-8, bo domyślnie zwraca właśnie UTF-8 i to nie jest rozpoznanie.
	if _, nazwa, pewne := charset.DetermineEncoding(bajty, ""); pewne &&
		!strings.EqualFold(nazwa, "utf-8") && !strings.EqualFold(nazwa, "windows-1252") {

		if przeklad, _ := charset.Lookup(nazwa); przeklad != nil {
			tekst, err := wejscieZastosujZapis(bajty, przeklad)
			if err == nil {
				return tekst, nazwa, nil
			}
		}
	}

	// Poprawny UTF-8 jest UTF-8. Napis, który przechodzi sprawdzenie ciągów
	// wielobajtowych, nie jest stroną kodową jednobajtową przez przypadek —
	// prawdopodobieństwo takiego zbiegu jest znikome dla tekstu z diakrytyką.
	if utf8.Valid(bajty) {
		return string(bajty), wejscieZapisUtf8, nil
	}

	// Zostają strony kodowe jednobajtowe. Rozstrzyga miara: która z nich daje
	// więcej liter polskich, ta jest właściwa. To jest ROZPOZNANIE PO ROZKŁADZIE,
	// nie odczyt deklaracji — i tak wychodzi w bilansie, żeby Operator wiedział,
	// że rdzeń zgadywał, a nie czytał.
	najlepsza, najlepszyTekst, najlepszaMiara := wejscieZapisNierozpoznny, "", -1
	for _, nazwa := range []string{
		wejscieZapisWindows1250, wejscieZapisIso88592, wejscieZapisWindows1252,
	} {
		przeklad := wejscieZapisPoNazwie(nazwa)
		if przeklad == nil {
			continue
		}
		tekst, err := wejscieZastosujZapis(bajty, przeklad)
		if err != nil {
			continue
		}
		miara := wejscieMiaraPolskosci(tekst)
		if miara > najlepszaMiara {
			najlepsza, najlepszyTekst, najlepszaMiara = nazwa, tekst, miara
		}
	}
	if najlepszyTekst == "" {
		return "", "", bladWskazaniaStudio(
			"pliku nie da się odczytać ani jako UTF-8, ani jako strony kodowej " +
				"windows-1250, ISO-8859-2 czy windows-1252; naprawa: wskazać zapis " +
				"znaków polem encoding")
	}
	return najlepszyTekst, najlepsza, nil
}

// wejscieZapisPoNazwie oddaje przekład zapisu znaków po nazwie.
func wejscieZapisPoNazwie(nazwa string) encoding.Encoding {
	switch nazwa {
	case wejscieZapisWindows1250:
		return charmap.Windows1250
	case wejscieZapisIso88592:
		return charmap.ISO8859_2
	case wejscieZapisWindows1252:
		return charmap.Windows1252
	}
	przeklad, _ := charset.Lookup(nazwa)
	return przeklad
}

// wejscieZastosujZapis przekłada bajty na UTF-8 wskazanym zapisem.
func wejscieZastosujZapis(bajty []byte, przeklad encoding.Encoding) (string, error) {
	if przeklad == nil {
		return string(bajty), nil
	}
	wynik, err := io.ReadAll(transform.NewReader(bytes.NewReader(bajty), przeklad.NewDecoder()))
	if err != nil {
		return "", bladWskazaniaStudio("nie można przełożyć pliku na UTF-8: " + err.Error())
	}
	return string(wynik), nil
}

// wejscieMiaraPolskosci liczy litery, które w polskim tekście występują,
// a w tekście odczytanym niewłaściwą stroną kodową zamieniają się w znaki
// sterujące albo w symbole. Miara jest miarą, nie dowodem — dlatego nazwa
// rozpoznanego zapisu idzie do bilansu.
func wejscieMiaraPolskosci(tekst string) int {
	const polskie = "ąćęłńóśźżĄĆĘŁŃÓŚŹŻ"
	miara := 0
	for _, znak := range tekst {
		switch {
		case strings.ContainsRune(polskie, znak):
			miara += 3
		case znak < 32 && znak != '\n' && znak != '\r' && znak != '\t':
			// Znak sterujący w tekście jest oznaką odczytu niewłaściwym zapisem.
			miara -= 5
		case znak == 0xFFFD:
			miara -= 10
		}
	}
	return miara
}

// ── Markdown ────────────────────────────────────────────────────────────────

// wejsciePostacZMarkdown składa postać dokumentu z markdown.
//
// Rozbiór jest własny, wierszowy, a nie przez `goldmark`: potrzebny jest tu
// przekład na STYLE NAZWANE i bloki dokumentu, a nie na HTML, który goldmark
// oddaje. Przekład HTML→postać dokumentu przez drugą drogę byłby dłuższy
// i gubiłby to samo. Rozpoznaje: nagłówki krzyżykami, cytat blokowy, listę
// wypunktowaną i numerowaną, tabelę kreskami, blok kodu i akapit.
func wejsciePostacZMarkdown(kodDokumentu, tresc string) (shared.StudioDocumentForm, []shared.StudioSkippedItem) {
	postac := wejscieNowaPostac(kodDokumentu, "", nil)
	postac.Blocks = nil
	sekcja := postac.Sections[0].Id
	pominiete := make([]shared.StudioSkippedItem, 0, 2)

	wiersze := strings.Split(strings.ReplaceAll(
		strings.ReplaceAll(tresc, "\r\n", "\n"), "\r", "\n"), "\n")

	listaWypunktowana := shared.StudioListDefinition{
		Id:   "lista-wypunktowana",
		Kind: shared.StudioListKindBullet,
		Levels: []shared.StudioListLevel{{
			Level:           1,
			BulletSource:    ooxmlWskaznikZrodlaWypunktowania(shared.StudioBulletSourceCharacter),
			BulletCharacter: wejscieWskaznikTekstu("•"),
			IndentMm:        wejscieWskaznikRzeczywisty(10),
		}},
	}
	listaNumerowana := shared.StudioListDefinition{
		Id:   "lista-numerowana",
		Kind: shared.StudioListKindNumber,
		Levels: []shared.StudioListLevel{{
			Level:        1,
			NumberFormat: ooxmlFormatNumeracjiListy("decimal"),
			Pattern:      wejscieWskaznikTekstu("%1."),
			IndentMm:     wejscieWskaznikRzeczywisty(10),
		}},
	}
	uzytaWypunktowana, uzytaNumerowana := false, false

	akapit := make([]string, 0, 4)
	domknijAkapit := func() {
		if len(akapit) == 0 {
			return
		}
		blok := wejscieBlokAkapitu(sekcja, strings.Join(akapit, " "), wejscieStylTekstZasadniczy, nil)
		blok.Runs = wejscieFragmentyZMarkdown(strings.Join(akapit, " "))
		postac.Blocks = append(postac.Blocks, blok)
		akapit = akapit[:0]
	}

	for numer := 0; numer < len(wiersze); numer++ {
		wiersz := wiersze[numer]
		przyciety := strings.TrimSpace(wiersz)

		switch {
		case przyciety == "":
			domknijAkapit()

		case strings.HasPrefix(przyciety, "```"), strings.HasPrefix(przyciety, "~~~"):
			// Blok kodu wychodzi akapitami o kroju stałej szerokości. Styl
			// „kod" nie stoi w arkuszu domyślnym, więc krój idzie postacią
			// znaku wprost — to strata nazwana, nie przemilczana.
			domknijAkapit()
			ogranicznik := przyciety[:3]
			numer++
			kod := make([]string, 0, 8)
			for numer < len(wiersze) && !strings.HasPrefix(strings.TrimSpace(wiersze[numer]), ogranicznik) {
				kod = append(kod, wiersze[numer])
				numer++
			}
			for _, wierszKodu := range kod {
				postac.Blocks = append(postac.Blocks, wejscieBlokAkapitu(sekcja, wierszKodu,
					wejscieStylTekstZasadniczy, &shared.StudioCharacterFormat{
						FontFamily: wejscieWskaznikTekstu("Courier New"),
						FontSizePt: wejscieWskaznikRzeczywisty(10),
					}))
			}
			pominiete = append(pominiete, shared.StudioSkippedItem{
				Reason: "blok kodu wszedł akapitami o kroju stałej szerokości — " +
					"arkusz stylów dokumentu nie ma stylu nazwanego dla kodu",
			})

		case strings.HasPrefix(przyciety, "#"):
			domknijAkapit()
			poziom := 0
			for poziom < len(przyciety) && przyciety[poziom] == '#' {
				poziom++
			}
			nazwaStylu := wejscieStylNaglowkaPoziomu(poziom)
			tekstNaglowka := strings.TrimSpace(strings.TrimLeft(przyciety, "#"))
			blok := wejscieBlokAkapitu(sekcja, tekstNaglowka, nazwaStylu, nil)
			blok.Paragraph.OutlineLevel = wejscieWskaznikCalkowity(poziom)
			postac.Blocks = append(postac.Blocks, blok)

		case strings.HasPrefix(przyciety, ">"):
			domknijAkapit()
			postac.Blocks = append(postac.Blocks, wejscieBlokAkapitu(sekcja,
				strings.TrimSpace(strings.TrimPrefix(przyciety, ">")), wejscieStylCytat, nil))

		case przyciety == "---", przyciety == "***", przyciety == "___":
			domknijAkapit()
			postac.Blocks = append(postac.Blocks, shared.StudioDocumentBlock{
				Id:        nowyIdentyfikator(przedrostekBlokuStudia),
				Kind:      wejscieRodzajBlokuPodzial,
				SectionId: wejscieWskaznikTekstu(sekcja),
				BreakKind: wejscieWskaznikRodzajuPodzialu(shared.StudioBreakKindPage),
			})

		case wejscieCzyWierszTabeli(przyciety) && numer+1 < len(wiersze) &&
			wejscieCzyWierszRozdzielajacyTabeli(wiersze[numer+1]):

			domknijAkapit()
			tabela, zuzyte := wejsciePostacTabeliZMarkdown(wiersze[numer:], &postac)
			if tabela != nil {
				postac.Tables = append(postac.Tables, *tabela)
				postac.Blocks = append(postac.Blocks, shared.StudioDocumentBlock{
					Id:        nowyIdentyfikator(przedrostekBlokuStudia),
					Kind:      wejscieRodzajBlokuTabela,
					SectionId: wejscieWskaznikTekstu(sekcja),
					TableId:   wejscieWskaznikTekstu(tabela.Id),
				})
				numer += zuzyte - 1
			}

		case strings.HasPrefix(przyciety, "- "), strings.HasPrefix(przyciety, "* "),
			strings.HasPrefix(przyciety, "+ "):

			domknijAkapit()
			uzytaWypunktowana = true
			poziom := 1 + (len(wiersz)-len(strings.TrimLeft(wiersz, " \t")))/2
			blok := wejscieBlokAkapitu(sekcja, strings.TrimSpace(przyciety[2:]),
				wejscieStylTekstZasadniczy, nil)
			blok.Runs = wejscieFragmentyZMarkdown(strings.TrimSpace(przyciety[2:]))
			blok.Paragraph.ListId = wejscieWskaznikTekstu(listaWypunktowana.Id)
			blok.Paragraph.ListLevel = wejscieWskaznikCalkowity(poziom)
			postac.Blocks = append(postac.Blocks, blok)

		case wejscieCzyWierszNumerowany(przyciety):
			domknijAkapit()
			uzytaNumerowana = true
			kropka := strings.Index(przyciety, ".")
			blok := wejscieBlokAkapitu(sekcja, strings.TrimSpace(przyciety[kropka+1:]),
				wejscieStylTekstZasadniczy, nil)
			blok.Runs = wejscieFragmentyZMarkdown(strings.TrimSpace(przyciety[kropka+1:]))
			blok.Paragraph.ListId = wejscieWskaznikTekstu(listaNumerowana.Id)
			blok.Paragraph.ListLevel = wejscieWskaznikCalkowity(1)
			postac.Blocks = append(postac.Blocks, blok)

		default:
			akapit = append(akapit, przyciety)
		}
	}
	domknijAkapit()

	if uzytaWypunktowana {
		postac.Lists = append(postac.Lists, listaWypunktowana)
	}
	if uzytaNumerowana {
		postac.Lists = append(postac.Lists, listaNumerowana)
	}
	if len(postac.Blocks) == 0 {
		postac.Blocks = wejscieNowaPostac(kodDokumentu, "", nil).Blocks
	}
	return postac, pominiete
}

// wejscieWskaznikRodzajuPodzialu oddaje wskaźnik na rodzaj podziału.
func wejscieWskaznikRodzajuPodzialu(wartosc shared.StudioBreakKind) *shared.StudioBreakKind {
	kopia := wartosc
	return &kopia
}

// wejscieFragmentyZMarkdown rozbija wiersz markdown na fragmenty o postaci
// znaku: pogrubienie, kursywę i kod w linii. Znaczników nieobsługiwanych nie
// zjada — zostają w treści jako znaki, bo zjedzenie ich byłoby utratą tekstu.
func wejscieFragmentyZMarkdown(wiersz string) []shared.StudioDocumentRun {
	fragmenty := make([]shared.StudioDocumentRun, 0, 4)
	znaki := []rune(wiersz)
	zwykly := make([]rune, 0, len(znaki))

	domknij := func() {
		if len(zwykly) == 0 {
			return
		}
		fragmenty = append(fragmenty, shared.StudioDocumentRun{Text: string(zwykly)})
		zwykly = zwykly[:0]
	}

	for i := 0; i < len(znaki); i++ {
		// Pogrubienie: dwa znaki gwiazdki albo podkreślenia.
		if i+1 < len(znaki) && (znaki[i] == '*' || znaki[i] == '_') && znaki[i+1] == znaki[i] {
			if koniec := wejscieSzukajZamkniecia(znaki, i+2, string([]rune{znaki[i], znaki[i]})); koniec > 0 {
				domknij()
				fragmenty = append(fragmenty, shared.StudioDocumentRun{
					Text:   string(znaki[i+2 : koniec]),
					Format: &shared.StudioCharacterFormat{Bold: wejscieWskaznikLogiczny(true)},
				})
				i = koniec + 1
				continue
			}
		}
		// Kursywa: jeden znak gwiazdki albo podkreślenia.
		if znaki[i] == '*' || znaki[i] == '_' {
			if koniec := wejscieSzukajZamkniecia(znaki, i+1, string(znaki[i])); koniec > 0 {
				domknij()
				fragmenty = append(fragmenty, shared.StudioDocumentRun{
					Text:   string(znaki[i+1 : koniec]),
					Format: &shared.StudioCharacterFormat{Italic: wejscieWskaznikLogiczny(true)},
				})
				i = koniec
				continue
			}
		}
		// Kod w linii: znak wstecznego apostrofu.
		if znaki[i] == '`' {
			if koniec := wejscieSzukajZamkniecia(znaki, i+1, "`"); koniec > 0 {
				domknij()
				fragmenty = append(fragmenty, shared.StudioDocumentRun{
					Text: string(znaki[i+1 : koniec]),
					Format: &shared.StudioCharacterFormat{
						FontFamily: wejscieWskaznikTekstu("Courier New"),
					},
				})
				i = koniec
				continue
			}
		}
		zwykly = append(zwykly, znaki[i])
	}
	domknij()
	if len(fragmenty) == 0 {
		return []shared.StudioDocumentRun{{Text: wiersz}}
	}
	return fragmenty
}

// wejscieSzukajZamkniecia szuka znacznika zamykającego od wskazanego miejsca.
// Zwraca położenie jego początku albo liczbę ujemną, gdy znacznika nie ma —
// wtedy znacznik otwierający nie był znacznikiem, tylko zwykłym znakiem.
func wejscieSzukajZamkniecia(znaki []rune, od int, znacznik string) int {
	szukany := []rune(znacznik)
	for i := od; i+len(szukany) <= len(znaki); i++ {
		zgodne := true
		for j, znak := range szukany {
			if znaki[i+j] != znak {
				zgodne = false
				break
			}
		}
		if zgodne && i > od {
			return i
		}
	}
	return -1
}

// wejscieCzyWierszTabeli rozstrzyga, czy wiersz jest wierszem tabeli markdown.
func wejscieCzyWierszTabeli(wiersz string) bool {
	return strings.Contains(wiersz, "|") && strings.Count(wiersz, "|") >= 2
}

// wejscieCzyWierszRozdzielajacyTabeli rozstrzyga, czy wiersz jest wierszem
// kresek rozdzielającym nagłówek tabeli od jej treści.
func wejscieCzyWierszRozdzielajacyTabeli(wiersz string) bool {
	przyciety := strings.TrimSpace(wiersz)
	if !strings.Contains(przyciety, "|") || !strings.Contains(przyciety, "-") {
		return false
	}
	for _, znak := range przyciety {
		switch znak {
		case '|', '-', ':', ' ', '\t':
		default:
			return false
		}
	}
	return true
}

// wejscieCzyWierszNumerowany rozstrzyga, czy wiersz zaczyna pozycję listy
// numerowanej („1. treść").
func wejscieCzyWierszNumerowany(wiersz string) bool {
	kropka := strings.Index(wiersz, ".")
	if kropka <= 0 || kropka+1 >= len(wiersz) || wiersz[kropka+1] != ' ' {
		return false
	}
	_, err := strconv.Atoi(wiersz[:kropka])
	return err == nil
}

// wejsciePostacTabeliZMarkdown składa tabelę z wierszy markdown i oddaje, ile
// wierszy zużyła. Wyrównanie kolumn bierze się z wiersza kresek (`:---`, `---:`,
// `:---:`) i wchodzi do postaci akapitu komórek — inaczej byłaby to wiedza
// z pliku wyrzucona przy wczytaniu.
func wejsciePostacTabeliZMarkdown(wiersze []string,
	postac *shared.StudioDocumentForm) (*shared.StudioDocumentTable, int) {

	if len(wiersze) < 2 {
		return nil, 0
	}
	naglowek := wejscieKomorkiWierszaMarkdown(wiersze[0])
	wyrownania := wejscieWyrownaniaKolumnMarkdown(wiersze[1])
	tresc := make([][]string, 0, len(wiersze))
	zuzyte := 2
	for i := 2; i < len(wiersze); i++ {
		if !wejscieCzyWierszTabeli(strings.TrimSpace(wiersze[i])) {
			break
		}
		tresc = append(tresc, wejscieKomorkiWierszaMarkdown(wiersze[i]))
		zuzyte++
	}

	kolumny := len(naglowek)
	for _, wiersz := range tresc {
		if len(wiersz) > kolumny {
			kolumny = len(wiersz)
		}
	}
	if kolumny == 0 {
		return nil, 0
	}

	tabela := shared.StudioDocumentTable{
		Id:           nowyIdentyfikator(przedrostekTabeliStudia),
		Rows:         1 + len(tresc),
		Columns:      kolumny,
		HeaderRows:   wejscieWskaznikCalkowity(1),
		RepeatHeader: wejscieWskaznikLogiczny(true),
		Border: &shared.StudioBorder{
			Style:   shared.StudioBorderStyleSingle,
			WidthPt: wejscieWskaznikRzeczywisty(0.5),
		},
	}
	// Szerokości kolumn są POLICZONE z obszaru pisania, nie zerowe: tabela
	// markdown szerokości nie niesie, a zero po zapisie do docx dałoby kolumny
	// niewidoczne.
	rowna := ooxmlSzerokoscObszaruPisania(postac) / float64(kolumny)
	for i := 0; i < kolumny; i++ {
		tabela.ColumnWidthsMm = append(tabela.ColumnWidthsMm, rowna)
	}

	wszystkie := append([][]string{naglowek}, tresc...)
	for numerWiersza, wiersz := range wszystkie {
		for kolumna := 0; kolumna < kolumny; kolumna++ {
			zawartosc := ""
			if kolumna < len(wiersz) {
				zawartosc = wiersz[kolumna]
			}
			komorka := shared.StudioTableCell{
				Row: numerWiersza, Column: kolumna,
				Text: wejscieWskaznikTekstu(zawartosc),
			}
			if kolumna < len(wyrownania) && wyrownania[kolumna] != nil {
				komorka.Paragraph = &shared.StudioParagraphFormat{Align: wyrownania[kolumna]}
			}
			if numerWiersza == 0 {
				komorka.Character = &shared.StudioCharacterFormat{
					Bold: wejscieWskaznikLogiczny(true),
				}
			}
			tabela.Cells = append(tabela.Cells, komorka)
		}
	}
	return &tabela, zuzyte
}

// wejscieKomorkiWierszaMarkdown rozbija wiersz tabeli na komórki.
func wejscieKomorkiWierszaMarkdown(wiersz string) []string {
	przyciety := strings.TrimSpace(wiersz)
	przyciety = strings.TrimPrefix(przyciety, "|")
	przyciety = strings.TrimSuffix(przyciety, "|")
	czesci := strings.Split(przyciety, "|")
	komorki := make([]string, 0, len(czesci))
	for _, czesc := range czesci {
		komorki = append(komorki, strings.TrimSpace(czesc))
	}
	return komorki
}

// wejscieWyrownaniaKolumnMarkdown czyta wyrównania kolumn z wiersza kresek.
func wejscieWyrownaniaKolumnMarkdown(wiersz string) []*shared.StudioTextAlign {
	komorki := wejscieKomorkiWierszaMarkdown(wiersz)
	wyrownania := make([]*shared.StudioTextAlign, 0, len(komorki))
	for _, komorka := range komorki {
		lewa := strings.HasPrefix(komorka, ":")
		prawa := strings.HasSuffix(komorka, ":")
		switch {
		case lewa && prawa:
			wyrownania = append(wyrownania, wejscieWskaznikWyrownania(shared.StudioTextAlignCenter))
		case prawa:
			wyrownania = append(wyrownania, wejscieWskaznikWyrownania(shared.StudioTextAlignRight))
		case lewa:
			wyrownania = append(wyrownania, wejscieWskaznikWyrownania(shared.StudioTextAlignLeft))
		default:
			wyrownania = append(wyrownania, nil)
		}
	}
	return wyrownania
}

// ── HTML ────────────────────────────────────────────────────────────────────

// wejsciePostacZHtml składa postać dokumentu z HTML.
//
// Rozbiór idzie `golang.org/x/net/html` — rozbiorem zgodnym z zachowaniem
// przeglądarki, a nie wyrażeniem regularnym po znacznikach. HTML Operatora
// bywa niepoprawny (niezamknięte znaczniki, atrybuty bez cudzysłowów), a ten
// rozbiór to znosi tak samo jak przeglądarka. Wyrażenie regularne po
// znacznikach rozsypałoby się na pierwszym takim pliku.
func wejsciePostacZHtml(kodDokumentu, tresc string) (shared.StudioDocumentForm, []shared.StudioSkippedItem) {
	pominiete := make([]shared.StudioSkippedItem, 0, 2)
	drzewo, err := html.Parse(strings.NewReader(tresc))
	if err != nil {
		return wejsciePostacZTekstu(kodDokumentu, tekstZeStronyStudia(tresc)),
			append(pominiete, shared.StudioSkippedItem{
				Reason: "HTML nie dał się rozebrać; treść weszła jako tekst płaski",
				Detail: wejscieWskaznikTekstu(err.Error()),
			})
	}
	postac := wejscieNowaPostac(kodDokumentu, "", nil)
	postac.Blocks = nil
	stan := &wejscieStanHtml{
		postac:    &postac,
		sekcja:    postac.Sections[0].Id,
		pominiete: &pominiete,
	}
	stan.przejdz(drzewo, shared.StudioCharacterFormat{})
	stan.domknijAkapit(wejscieStylTekstZasadniczy, 0)
	if len(postac.Blocks) == 0 {
		postac.Blocks = wejscieNowaPostac(kodDokumentu, "", nil).Blocks
	}
	return postac, pominiete
}

// wejscieStanHtml zbiera stan rozbioru HTML: składany akapit i wykaz pominięć.
type wejscieStanHtml struct {
	postac    *shared.StudioDocumentForm
	sekcja    string
	pominiete *[]shared.StudioSkippedItem
	fragmenty []shared.StudioDocumentRun
}

// przejdz przechodzi drzewo HTML, składając bloki dokumentu.
//
// Postać znaku PŁYNIE W DÓŁ drzewa: `<b><i>tekst</i></b>` daje fragment
// pogrubiony i pochylony naraz, bo każdy poziom dokłada swoją cechę do postaci
// odziedziczonej. Postać liczona osobno na każdym poziomie zgubiłaby zewnętrzną.
func (s *wejscieStanHtml) przejdz(wezel *html.Node, odziedziczona shared.StudioCharacterFormat) {
	if wezel.Type == html.TextNode {
		tekst := strings.ReplaceAll(wezel.Data, "\n", " ")
		if strings.TrimSpace(tekst) == "" {
			// Sam odstęp między znacznikami blokowymi nie jest treścią, ale
			// odstęp między znacznikami tekstowymi jest — inaczej „<b>a</b> <i>b</i>"
			// dałoby „ab".
			if tekst != "" && len(s.fragmenty) > 0 {
				s.fragmenty = append(s.fragmenty, shared.StudioDocumentRun{Text: " "})
			}
			return
		}
		postac := odziedziczona
		s.fragmenty = append(s.fragmenty, shared.StudioDocumentRun{
			Text: tekst, Format: wejsciePostacZnakuAlboNic(postac),
		})
		return
	}
	if wezel.Type != html.ElementNode {
		for dziecko := wezel.FirstChild; dziecko != nil; dziecko = dziecko.NextSibling {
			s.przejdz(dziecko, odziedziczona)
		}
		return
	}

	switch wezel.Data {
	case "script", "style", "noscript", "head", "meta", "link":
		// Skrypt i arkusz nie są treścią dokumentu.
		return

	case "h1", "h2", "h3", "h4", "h5", "h6":
		s.domknijAkapit(wejscieStylTekstZasadniczy, 0)
		poziom, _ := strconv.Atoi(wezel.Data[1:])
		s.przejdzDzieci(wezel, odziedziczona)
		s.domknijAkapit(wejscieStylNaglowkaPoziomu(poziom), poziom)
		return

	case "p", "div", "section", "article", "li", "dd", "dt":
		s.domknijAkapit(wejscieStylTekstZasadniczy, 0)
		s.przejdzDzieci(wezel, odziedziczona)
		s.domknijAkapit(wejscieStylTekstZasadniczy, 0)
		return

	case "blockquote":
		s.domknijAkapit(wejscieStylTekstZasadniczy, 0)
		s.przejdzDzieci(wezel, odziedziczona)
		s.domknijAkapit(wejscieStylCytat, 0)
		return

	case "figcaption", "caption":
		s.domknijAkapit(wejscieStylTekstZasadniczy, 0)
		s.przejdzDzieci(wezel, odziedziczona)
		s.domknijAkapit(wejscieStylPodpis, 0)
		return

	case "br":
		s.fragmenty = append(s.fragmenty, shared.StudioDocumentRun{Text: "\n"})
		return

	case "hr":
		s.domknijAkapit(wejscieStylTekstZasadniczy, 0)
		s.postac.Blocks = append(s.postac.Blocks, shared.StudioDocumentBlock{
			Id:        nowyIdentyfikator(przedrostekBlokuStudia),
			Kind:      wejscieRodzajBlokuPodzial,
			SectionId: wejscieWskaznikTekstu(s.sekcja),
			BreakKind: wejscieWskaznikRodzajuPodzialu(shared.StudioBreakKindPage),
		})
		return

	case "table":
		s.domknijAkapit(wejscieStylTekstZasadniczy, 0)
		s.czytajTabeleHtml(wezel)
		return

	case "img":
		s.czytajObrazHtml(wezel)
		return

	case "a":
		adres := wejscieAtrybutHtml(wezel, "href")
		poczatek := len(s.fragmenty)
		postac := odziedziczona
		postac.Color = wejscieWskaznikTekstu("#0563C1")
		postac.Underline = ooxmlPodkreslenie("single")
		s.przejdzDzieci(wezel, postac)
		if adres != "" && len(s.fragmenty) > poczatek {
			etykieta := make([]string, 0, len(s.fragmenty)-poczatek)
			for _, fragment := range s.fragmenty[poczatek:] {
				etykieta = append(etykieta, fragment.Text)
			}
			s.postac.Apparatus = append(s.postac.Apparatus, shared.StudioApparatusItem{
				Id:        nowyIdentyfikator(przedrostekAparatuWejscia),
				Kind:      shared.StudioApparatusKindHyperlink,
				Label:     wejscieWskaznikTekstu(strings.Join(etykieta, "")),
				TargetUrl: wejscieWskaznikTekstu(adres),
			})
		}
		return
	}

	// Znaczniki, które dokładają wyłącznie postać znaku.
	postac := odziedziczona
	switch wezel.Data {
	case "b", "strong":
		postac.Bold = wejscieWskaznikLogiczny(true)
	case "i", "em", "cite":
		postac.Italic = wejscieWskaznikLogiczny(true)
	case "u", "ins":
		postac.Underline = ooxmlPodkreslenie("single")
	case "s", "strike", "del":
		postac.Strikethrough = wejscieWskaznikLogiczny(true)
	case "sup":
		postac.Superscript = wejscieWskaznikLogiczny(true)
	case "sub":
		postac.Subscript = wejscieWskaznikLogiczny(true)
	case "mark":
		postac.HighlightColor = wejscieWskaznikTekstu("#FFFF00")
	case "small":
		postac.FontSizePt = wejscieWskaznikRzeczywisty(10)
	case "code", "kbd", "samp", "pre", "tt":
		postac.FontFamily = wejscieWskaznikTekstu("Courier New")
	}
	if styl := wejscieAtrybutHtml(wezel, "style"); styl != "" {
		wejscieDopiszPostacZeStyluHtml(&postac, styl)
	}
	s.przejdzDzieci(wezel, postac)
}

// przejdzDzieci przechodzi dzieci węzła.
func (s *wejscieStanHtml) przejdzDzieci(wezel *html.Node, postac shared.StudioCharacterFormat) {
	for dziecko := wezel.FirstChild; dziecko != nil; dziecko = dziecko.NextSibling {
		s.przejdz(dziecko, postac)
	}
}

// domknijAkapit odkłada zebrane fragmenty jako blok akapitu.
func (s *wejscieStanHtml) domknijAkapit(styl string, poziomKonspektu int) {
	if len(s.fragmenty) == 0 {
		return
	}
	// Fragment, który po złożeniu jest samym odstępem, nie jest akapitem.
	pelny := ""
	for _, fragment := range s.fragmenty {
		pelny += fragment.Text
	}
	if strings.TrimSpace(pelny) == "" {
		s.fragmenty = nil
		return
	}
	blok := shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuAkapit,
		SectionId: wejscieWskaznikTekstu(s.sekcja),
		Paragraph: &shared.StudioParagraphFormat{
			StyleName: wejscieWskaznikTekstu(styl),
		},
		Runs: s.fragmenty,
	}
	if poziomKonspektu > 0 {
		blok.Paragraph.OutlineLevel = wejscieWskaznikCalkowity(poziomKonspektu)
	}
	s.postac.Blocks = append(s.postac.Blocks, blok)
	s.fragmenty = nil
}

// czytajTabeleHtml składa tabelę z HTML wraz ze scaleniami.
func (s *wejscieStanHtml) czytajTabeleHtml(wezel *html.Node) {
	wiersze := wejscieWezlyHtml(wezel, "tr")
	if len(wiersze) == 0 {
		return
	}
	tabela := shared.StudioDocumentTable{
		Id:   nowyIdentyfikator(przedrostekTabeliStudia),
		Rows: len(wiersze),
		Border: &shared.StudioBorder{
			Style:   shared.StudioBorderStyleSingle,
			WidthPt: wejscieWskaznikRzeczywisty(0.5),
		},
	}
	naglowki := 0
	for numerWiersza, wiersz := range wiersze {
		kolumna := 0
		wierszNaglowkowy := true
		for komorkaHtml := wiersz.FirstChild; komorkaHtml != nil; komorkaHtml = komorkaHtml.NextSibling {
			if komorkaHtml.Type != html.ElementNode ||
				(komorkaHtml.Data != "td" && komorkaHtml.Data != "th") {
				continue
			}
			if komorkaHtml.Data != "th" {
				wierszNaglowkowy = false
			}
			rozpietoscKolumn := wejscieLiczbaAtrybutuHtml(komorkaHtml, "colspan", 1)
			rozpietoscWierszy := wejscieLiczbaAtrybutuHtml(komorkaHtml, "rowspan", 1)
			komorka := shared.StudioTableCell{
				Row: numerWiersza, Column: kolumna,
				Text: wejscieWskaznikTekstu(strings.TrimSpace(wejscieTekstHtml(komorkaHtml))),
			}
			if rozpietoscKolumn > 1 {
				komorka.ColumnSpan = wejscieWskaznikCalkowity(rozpietoscKolumn)
			}
			if rozpietoscWierszy > 1 {
				komorka.RowSpan = wejscieWskaznikCalkowity(rozpietoscWierszy)
			}
			if komorkaHtml.Data == "th" {
				komorka.Character = &shared.StudioCharacterFormat{
					Bold: wejscieWskaznikLogiczny(true),
				}
			}
			tabela.Cells = append(tabela.Cells, komorka)
			for przesuniecie := 1; przesuniecie < rozpietoscKolumn; przesuniecie++ {
				tabela.Cells = append(tabela.Cells, shared.StudioTableCell{
					Row: numerWiersza, Column: kolumna + przesuniecie,
					Merged: wejscieWskaznikLogiczny(true),
				})
			}
			kolumna += rozpietoscKolumn
		}
		if kolumna > tabela.Columns {
			tabela.Columns = kolumna
		}
		if wierszNaglowkowy && kolumna > 0 && numerWiersza == naglowki {
			naglowki++
		}
	}
	if tabela.Columns == 0 {
		return
	}
	if naglowki > 0 {
		tabela.HeaderRows = wejscieWskaznikCalkowity(naglowki)
		tabela.RepeatHeader = wejscieWskaznikLogiczny(true)
	}
	rowna := ooxmlSzerokoscObszaruPisania(s.postac) / float64(tabela.Columns)
	for i := 0; i < tabela.Columns; i++ {
		tabela.ColumnWidthsMm = append(tabela.ColumnWidthsMm, rowna)
	}
	s.postac.Tables = append(s.postac.Tables, tabela)
	s.postac.Blocks = append(s.postac.Blocks, shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuTabela,
		SectionId: wejscieWskaznikTekstu(s.sekcja),
		TableId:   wejscieWskaznikTekstu(tabela.Id),
	})
}

// czytajObrazHtml zakłada obiekt obrazu wskazanego znacznikiem.
//
// Bajtów obrazu ten rachunek NIE pobiera: pobranie z sieci należy do czynności
// wniesienia ze strony (`studio.insert.from.web`), która ma na to kontekst
// żądania i granicę czasu. Tutaj powstaje obiekt wraz z adresem — i to jest
// prawda o tym, co rdzeń wie.
func (s *wejscieStanHtml) czytajObrazHtml(wezel *html.Node) {
	adres := wejscieAtrybutHtml(wezel, "src")
	if adres == "" {
		return
	}
	obiekt := shared.StudioDocumentObject{
		Id:        nowyIdentyfikator(przedrostekObiektuStudia),
		Kind:      shared.StudioObjectKindImage,
		Source:    wejscieWskaznikZrodlaObiektu(shared.StudioObjectSourceWeb),
		SourceUrl: wejscieWskaznikTekstu(adres),
		AltText:   wejscieWskaznikTekstu(wejscieAtrybutHtml(wezel, "alt")),
	}
	if szerokosc := wejscieLiczbaAtrybutuHtml(wezel, "width", 0); szerokosc > 0 {
		// Wymiary HTML są w punktach obrazu; jeden punkt obrazu to 25,4/96 mm.
		obiekt.WidthMm = wejscieWskaznikRzeczywisty(float64(szerokosc) * 25.4 / 96)
	}
	if wysokosc := wejscieLiczbaAtrybutuHtml(wezel, "height", 0); wysokosc > 0 {
		obiekt.HeightMm = wejscieWskaznikRzeczywisty(float64(wysokosc) * 25.4 / 96)
	}
	s.postac.Objects = append(s.postac.Objects, obiekt)
	s.postac.Blocks = append(s.postac.Blocks, shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuObiekt,
		SectionId: wejscieWskaznikTekstu(s.sekcja),
		ObjectId:  wejscieWskaznikTekstu(obiekt.Id),
	})
	*s.pominiete = append(*s.pominiete, shared.StudioSkippedItem{
		Reason: "obraz wszedł jako obiekt z adresem źródła; bajtów rdzeń nie pobierał",
		Detail: wejscieWskaznikTekstu(adres),
	})
}

// wejsciePostacZnakuAlboNic oddaje postać znaku albo nic, gdy postać jest pusta.
func wejsciePostacZnakuAlboNic(postac shared.StudioCharacterFormat) *shared.StudioCharacterFormat {
	if postac.FontFamily == nil && postac.FontSizePt == nil && postac.Bold == nil &&
		postac.Italic == nil && postac.Underline == nil && postac.Strikethrough == nil &&
		postac.Superscript == nil && postac.Subscript == nil && postac.Color == nil &&
		postac.HighlightColor == nil && postac.SmallCaps == nil && postac.AllCaps == nil {
		return nil
	}
	kopia := postac
	return &kopia
}

// wejscieDopiszPostacZeStyluHtml czyta z atrybutu stylu te cechy, które kontrakt
// niesie. Nastaw nieznanych nie zgaduje — styl HTML ma setki cech, a udawanie,
// że rdzeń rozumie je wszystkie, byłoby nieprawdą.
func wejscieDopiszPostacZeStyluHtml(postac *shared.StudioCharacterFormat, styl string) {
	for _, nastawa := range strings.Split(styl, ";") {
		czesci := strings.SplitN(nastawa, ":", 2)
		if len(czesci) != 2 {
			continue
		}
		nazwa := strings.ToLower(strings.TrimSpace(czesci[0]))
		wartosc := strings.TrimSpace(czesci[1])
		switch nazwa {
		case "color":
			if barwa := wejscieBarwaZeStyluHtml(wartosc); barwa != "" {
				postac.Color = wejscieWskaznikTekstu(barwa)
			}
		case "background-color", "background":
			if barwa := wejscieBarwaZeStyluHtml(wartosc); barwa != "" {
				postac.HighlightColor = wejscieWskaznikTekstu(barwa)
			}
		case "font-weight":
			if wartosc == "bold" || wartosc == "bolder" || wartosc >= "600" {
				postac.Bold = wejscieWskaznikLogiczny(true)
			}
		case "font-style":
			if wartosc == "italic" || wartosc == "oblique" {
				postac.Italic = wejscieWskaznikLogiczny(true)
			}
		case "font-family":
			krój := strings.Trim(strings.SplitN(wartosc, ",", 2)[0], ` "'`)
			if krój != "" {
				postac.FontFamily = wejscieWskaznikTekstu(krój)
			}
		case "font-size":
			if stopien, jest := wejscieStopienZeStyluHtml(wartosc); jest {
				postac.FontSizePt = wejscieWskaznikRzeczywisty(stopien)
			}
		case "text-decoration", "text-decoration-line":
			if strings.Contains(wartosc, "underline") {
				postac.Underline = ooxmlPodkreslenie("single")
			}
			if strings.Contains(wartosc, "line-through") {
				postac.Strikethrough = wejscieWskaznikLogiczny(true)
			}
		}
	}
}

// wejscieBarwaZeStyluHtml czyta barwę zapisaną szesnastkowo. Nazw barw i zapisu
// `rgb()` nie tłumaczy — to byłby drugi wykaz barw w rdzeniu.
func wejscieBarwaZeStyluHtml(wartosc string) string {
	oczyszczona := strings.TrimSpace(wartosc)
	if !strings.HasPrefix(oczyszczona, "#") {
		return ""
	}
	zapis := strings.TrimPrefix(oczyszczona, "#")
	if len(zapis) == 3 {
		// Zapis skrócony rozwija się przez podwojenie każdego znaku.
		zapis = string([]byte{zapis[0], zapis[0], zapis[1], zapis[1], zapis[2], zapis[2]})
	}
	if len(zapis) < 6 {
		return ""
	}
	if _, err := strconv.ParseUint(zapis[:6], 16, 64); err != nil {
		return ""
	}
	return "#" + strings.ToUpper(zapis[:6])
}

// wejscieStopienZeStyluHtml czyta stopień pisma podany w punktach albo
// w punktach obrazu.
func wejscieStopienZeStyluHtml(wartosc string) (float64, bool) {
	oczyszczona := strings.ToLower(strings.TrimSpace(wartosc))
	switch {
	case strings.HasSuffix(oczyszczona, "pt"):
		if liczba, err := strconv.ParseFloat(strings.TrimSuffix(oczyszczona, "pt"), 64); err == nil {
			return liczba, true
		}
	case strings.HasSuffix(oczyszczona, "px"):
		if liczba, err := strconv.ParseFloat(strings.TrimSuffix(oczyszczona, "px"), 64); err == nil {
			// Punkt obrazu to trzy czwarte punktu typograficznego.
			return liczba * 0.75, true
		}
	}
	return 0, false
}

// wejscieAtrybutHtml oddaje atrybut węzła HTML.
func wejscieAtrybutHtml(wezel *html.Node, nazwa string) string {
	for _, atrybut := range wezel.Attr {
		if strings.EqualFold(atrybut.Key, nazwa) {
			return strings.TrimSpace(atrybut.Val)
		}
	}
	return ""
}

// wejscieLiczbaAtrybutuHtml oddaje atrybut jako liczbę albo wartość domyślną.
func wejscieLiczbaAtrybutuHtml(wezel *html.Node, nazwa string, domyslna int) int {
	tekst := strings.TrimSuffix(wejscieAtrybutHtml(wezel, nazwa), "px")
	liczba, err := strconv.Atoi(strings.TrimSpace(tekst))
	if err != nil || liczba <= 0 {
		return domyslna
	}
	return liczba
}

// wejscieWezlyHtml zbiera wgłąb węzły o wskazanej nazwie.
func wejscieWezlyHtml(wezel *html.Node, nazwa string) []*html.Node {
	wynik := make([]*html.Node, 0, 8)
	var przejdz func(*html.Node)
	przejdz = func(biezacy *html.Node) {
		if biezacy.Type == html.ElementNode && biezacy.Data == nazwa {
			wynik = append(wynik, biezacy)
		}
		for dziecko := biezacy.FirstChild; dziecko != nil; dziecko = dziecko.NextSibling {
			przejdz(dziecko)
		}
	}
	przejdz(wezel)
	return wynik
}

// wejscieTekstHtml składa treść tekstową węzła wraz z dziećmi.
func wejscieTekstHtml(wezel *html.Node) string {
	var budowa strings.Builder
	var przejdz func(*html.Node)
	przejdz = func(biezacy *html.Node) {
		if biezacy.Type == html.TextNode {
			budowa.WriteString(biezacy.Data)
			return
		}
		if biezacy.Type == html.ElementNode &&
			(biezacy.Data == "script" || biezacy.Data == "style") {
			return
		}
		if biezacy.Type == html.ElementNode && biezacy.Data == "br" {
			budowa.WriteString("\n")
		}
		for dziecko := biezacy.FirstChild; dziecko != nil; dziecko = dziecko.NextSibling {
			przejdz(dziecko)
		}
	}
	przejdz(wezel)
	return strings.Join(strings.Fields(budowa.String()), " ")
}

// ── RTF ─────────────────────────────────────────────────────────────────────

// wejsciePostacZRtf składa postać dokumentu z RTF.
//
// RTF jest formatem znakowym o składni własnej — nie XML i nie archiwum.
// Rozbiór jest własny i obejmuje to, co niesie pismo: akapity (`\par`),
// pogrubienie, kursywę, podkreślenie, stopień pisma (`\fsN`, w półpunktach),
// wyrównanie (`\qc`, `\qr`, `\qj`) oraz znaki spoza zakresu jednobajtowego
// zapisane jako `\'hh` i `\uN`.
//
// UCZCIWIE: tabele RTF (`\trowd`) i obrazy (`\pict`) NIE są odzyskiwane. Tabela
// RTF jest ciągiem akapitów z granicami komórek zapisanymi w rozkazach składu,
// więc jej odzyskanie jest odtworzeniem układu, nie odczytem struktury. Wchodzi
// treścią i wychodzi w bilansie jako układ nierozpoznany.
func wejsciePostacZRtf(kodDokumentu, tresc string) (shared.StudioDocumentForm, []shared.StudioSkippedItem) {
	pominiete := make([]shared.StudioSkippedItem, 0, 2)
	postac := wejscieNowaPostac(kodDokumentu, "", nil)
	postac.Blocks = nil
	sekcja := postac.Sections[0].Id

	stan := wejscieStanRtf{}
	stos := make([]wejscieStanRtf, 0, 16)
	fragmenty := make([]shared.StudioDocumentRun, 0, 8)
	biezacy := strings.Builder{}
	wyrownanie := (*shared.StudioTextAlign)(nil)
	tabel := 0
	obrazow := 0
	pomijanaGrupa := 0

	domknijFragment := func() {
		if biezacy.Len() == 0 {
			return
		}
		fragmenty = append(fragmenty, shared.StudioDocumentRun{
			Text: biezacy.String(), Format: stan.postacZnaku(),
		})
		biezacy.Reset()
	}
	domknijAkapit := func() {
		domknijFragment()
		if len(fragmenty) == 0 {
			return
		}
		blok := shared.StudioDocumentBlock{
			Id:        nowyIdentyfikator(przedrostekBlokuStudia),
			Kind:      wejscieRodzajBlokuAkapit,
			SectionId: wejscieWskaznikTekstu(sekcja),
			Paragraph: &shared.StudioParagraphFormat{
				StyleName: wejscieWskaznikTekstu(wejscieStylTekstZasadniczy),
				Align:     wyrownanie,
			},
			Runs: fragmenty,
		}
		postac.Blocks = append(postac.Blocks, blok)
		fragmenty = make([]shared.StudioDocumentRun, 0, 8)
	}

	znaki := []byte(tresc)
	for i := 0; i < len(znaki); i++ {
		switch znaki[i] {
		case '{':
			stos = append(stos, stan)
			if pomijanaGrupa > 0 {
				pomijanaGrupa++
			}
			continue
		case '}':
			if pomijanaGrupa > 0 {
				pomijanaGrupa--
			}
			if len(stos) > 0 {
				domknijFragment()
				stan = stos[len(stos)-1]
				stos = stos[:len(stos)-1]
			}
			continue
		case '\r', '\n':
			continue
		case '\\':
			rozkaz, parametr, maParametr, dalej := wejscieRozkazRtf(znaki, i)
			i = dalej
			if rozkaz == "" {
				continue
			}
			switch rozkaz {
			case "'":
				// Znak zapisany szesnastkowo w bieżącej stronie kodowej. RTF
				// pisma polskiego jedzie zwykle windows-1250, i tak go czytamy —
				// nagłówek `\ansicpgN` bywa nieprawdziwy częściej niż strona
				// kodowa, którą naprawdę użyto.
				if pomijanaGrupa == 0 && maParametr {
					przelozone, err := wejscieZastosujZapis([]byte{byte(parametr)}, charmap.Windows1250)
					if err == nil {
						biezacy.WriteString(przelozone)
					}
				}
			case "u":
				if pomijanaGrupa == 0 && maParametr {
					// `\uN` niesie znak Unicode; wartość ujemna jest zapisem
					// liczby bez znaku, którą RTF wyraża w zakresie ujemnym.
					kod := parametr
					if kod < 0 {
						kod += 65536
					}
					biezacy.WriteRune(rune(kod))
				}
			case "par", "line", "sect":
				if pomijanaGrupa == 0 {
					domknijAkapit()
				}
			case "tab":
				if pomijanaGrupa == 0 {
					biezacy.WriteString("\t")
				}
			case "b", "i", "ul", "ulnone", "strike", "super", "sub", "nosupersub", "plain":
				domknijFragment()
				stan.zastosujPrzelacznik(rozkaz, parametr, maParametr)
			case "fs":
				if maParametr {
					domknijFragment()
					stan.stopien = wejscieWskaznikRzeczywisty(float64(parametr) / 2)
				}
			case "cf":
				// Barwa idzie numerem w tabeli barw dokumentu. Tabeli nie
				// czytamy, więc numeru nie da się przełożyć na barwę i pole
				// zostaje puste — to jest brak, nie zmyślona barwa.
			case "ql":
				wyrownanie = wejscieWskaznikWyrownania(shared.StudioTextAlignLeft)
			case "qc":
				wyrownanie = wejscieWskaznikWyrownania(shared.StudioTextAlignCenter)
			case "qr":
				wyrownanie = wejscieWskaznikWyrownania(shared.StudioTextAlignRight)
			case "qj":
				wyrownanie = wejscieWskaznikWyrownania(shared.StudioTextAlignJustify)
			case "trowd":
				tabel++
			case "pict":
				obrazow++
				pomijanaGrupa = 1
			case "fonttbl", "colortbl", "stylesheet", "info", "listtable",
				"listoverridetable", "generator", "themedata", "datastore",
				"rsidtbl", "xmlnstbl", "latentstyles":
				// Grupy nagłówkowe nie są treścią; ich zawartość się pomija.
				pomijanaGrupa = 1
			}
			continue
		}
		if pomijanaGrupa == 0 {
			biezacy.WriteByte(znaki[i])
		}
	}
	domknijAkapit()

	if tabel > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "tabele RTF weszły treścią akapitów, bez struktury wierszy i kolumn — " +
				"RTF wyraża tabelę rozkazami składu, nie strukturą",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(tabel) + " wierszy tabelarycznych"),
		})
	}
	if obrazow > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "obrazy osadzone w RTF pominięto",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(obrazow) + " obrazów"),
		})
	}
	if len(postac.Blocks) == 0 {
		postac.Blocks = wejscieNowaPostac(kodDokumentu, "", nil).Blocks
	}
	return postac, pominiete
}

// wejscieStanRtf trzyma postać znaku bieżącej grupy RTF. Grupa dziedziczy
// postać po grupie nadrzędnej i oddaje ją przy zamknięciu — dlatego stan jest
// odkładany na stos, a nie liczony na bieżąco.
type wejscieStanRtf struct {
	pogrubienie   bool
	kursywa       bool
	podkreslenie  bool
	przekreslenie bool
	indeksGorny   bool
	indeksDolny   bool
	stopien       *float64
}

// zastosujPrzelacznik przestawia przełącznik postaci znaku. Rozkaz RTF
// z parametrem zero WYŁĄCZA cechę (`\b0`), a bez parametru ją włącza.
func (s *wejscieStanRtf) zastosujPrzelacznik(rozkaz string, parametr int, maParametr bool) {
	wlaczony := !maParametr || parametr != 0
	switch rozkaz {
	case "b":
		s.pogrubienie = wlaczony
	case "i":
		s.kursywa = wlaczony
	case "ul":
		s.podkreslenie = wlaczony
	case "ulnone":
		s.podkreslenie = false
	case "strike":
		s.przekreslenie = wlaczony
	case "super":
		s.indeksGorny, s.indeksDolny = wlaczony, false
	case "sub":
		s.indeksDolny, s.indeksGorny = wlaczony, false
	case "nosupersub":
		s.indeksGorny, s.indeksDolny = false, false
	case "plain":
		*s = wejscieStanRtf{}
	}
}

// postacZnaku składa postać znaku ze stanu grupy RTF.
func (s wejscieStanRtf) postacZnaku() *shared.StudioCharacterFormat {
	postac := shared.StudioCharacterFormat{FontSizePt: s.stopien}
	if s.pogrubienie {
		postac.Bold = wejscieWskaznikLogiczny(true)
	}
	if s.kursywa {
		postac.Italic = wejscieWskaznikLogiczny(true)
	}
	if s.podkreslenie {
		postac.Underline = ooxmlPodkreslenie("single")
	}
	if s.przekreslenie {
		postac.Strikethrough = wejscieWskaznikLogiczny(true)
	}
	if s.indeksGorny {
		postac.Superscript = wejscieWskaznikLogiczny(true)
	}
	if s.indeksDolny {
		postac.Subscript = wejscieWskaznikLogiczny(true)
	}
	return wejsciePostacZnakuAlboNic(postac)
}

// wejscieRozkazRtf czyta rozkaz RTF od znaku odwrotnego ukośnika. Oddaje nazwę
// rozkazu, parametr liczbowy, informację, czy parametr stał w pliku, oraz
// położenie ostatniego zużytego bajtu.
func wejscieRozkazRtf(znaki []byte, poczatek int) (string, int, bool, int) {
	i := poczatek + 1
	if i >= len(znaki) {
		return "", 0, false, i
	}
	// Znak sterujący zapisany ukośnikiem: `\\`, `\{`, `\}` są znakami treści.
	switch znaki[i] {
	case '\\', '{', '}':
		return "", 0, false, i
	case '\'':
		if i+2 < len(znaki) {
			if wartosc, err := strconv.ParseUint(string(znaki[i+1:i+3]), 16, 16); err == nil {
				return "'", int(wartosc), true, i + 2
			}
		}
		return "'", 0, false, i
	case '*', '~', '-', '_', ':', '|':
		return "", 0, false, i
	}

	nazwa := i
	for i < len(znaki) && ((znaki[i] >= 'a' && znaki[i] <= 'z') || (znaki[i] >= 'A' && znaki[i] <= 'Z')) {
		i++
	}
	rozkaz := string(znaki[nazwa:i])
	if rozkaz == "" {
		return "", 0, false, i
	}

	parametr, maParametr := 0, false
	if i < len(znaki) && (znaki[i] == '-' || (znaki[i] >= '0' && znaki[i] <= '9')) {
		liczba := i
		if znaki[i] == '-' {
			i++
		}
		for i < len(znaki) && znaki[i] >= '0' && znaki[i] <= '9' {
			i++
		}
		if wartosc, err := strconv.Atoi(string(znaki[liczba:i])); err == nil {
			parametr, maParametr = wartosc, true
		}
	}
	// Odstęp po rozkazie należy do rozkazu, nie do treści.
	if i < len(znaki) && znaki[i] == ' ' {
		return rozkaz, parametr, maParametr, i
	}
	return rozkaz, parametr, maParametr, i - 1
}

// ── Wykaz nośników ──────────────────────────────────────────────────────────

// wejscieWymiarNosnika jest wymiarami nośnika w milimetrach.
type wejscieWymiarNosnika struct {
	szerokosc float64
	wysokosc  float64
}

// wejscieWymiaryNosnika oddaje wymiary nośnika po nazwie.
//
// Wykaz nośników NIE JEST tu zakładany od nowa: czytany jest wykaz wspólny
// rdzenia (`wykazNosnikowDruku` z `nosniki_druku_wspolne.go`), bo dwa wykazy
// rozmiarów w jednym produkcie rozjadą się przy pierwszej poprawce i Operator
// dostanie A3 o wymiarach A4. Wcześniej stał tu odczyt wykazu własnego modułu
// Design; wykaz wspólny zastąpił oba i niesie koperty C4, C5 oraz C6 wskazane
// przez Właściciela.
func wejscieWymiaryNosnika(nazwa string) (wejscieWymiarNosnika, bool) {
	szukana := strings.ToLower(strings.TrimSpace(nazwa))
	if szukana == "" {
		return wejscieWymiarNosnika{}, false
	}
	for _, nosnik := range wykazNosnikowDruku {
		if strings.ToLower(nosnik.Nazwa) == szukana {
			return wejscieWymiarNosnika{
				szerokosc: nosnik.SzerokoscMm, wysokosc: nosnik.WysokoscMm,
			}, true
		}
	}
	return wejscieWymiarNosnika{}, false
}

// wejscieNazwaNosnikaZWymiarow szuka nazwy nośnika o wskazanych wymiarach.
// Tolerancja jednego milimetra bierze się z przeliczeń: A4 zapisane w twipach
// i odczytane z powrotem daje 209,97 mm, a to jest ten sam A4.
func wejscieNazwaNosnikaZWymiarow(szerokosc, wysokosc float64) (string, bool) {
	for _, nosnik := range wykazNosnikowDruku {
		zgodnaPionowo := wejscieBliskie(nosnik.SzerokoscMm, szerokosc) &&
			wejscieBliskie(nosnik.WysokoscMm, wysokosc)
		zgodnaPoziomo := wejscieBliskie(nosnik.SzerokoscMm, wysokosc) &&
			wejscieBliskie(nosnik.WysokoscMm, szerokosc)
		if zgodnaPionowo || zgodnaPoziomo {
			return nosnik.Nazwa, true
		}
	}
	return "", false
}

// wejscieBliskie rozstrzyga, czy dwa wymiary są tym samym wymiarem.
func wejscieBliskie(pierwszy, drugi float64) bool {
	roznica := pierwszy - drugi
	if roznica < 0 {
		roznica = -roznica
	}
	return roznica <= 1
}
