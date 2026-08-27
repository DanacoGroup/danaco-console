// Odpowiedzialność pliku: odczyt i zapis formatów dokumentu w module
// Translate. Cała wiedza o tym, jak z pliku wyjąć akapity i jak złożyć je
// z powrotem, mieszka tutaj; formaty czyta się bibliotekami Go, bez programu
// obcego poza rozpoznaniem pisma.
package core

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// formatDokumentuZeSciezki rozpoznaje format po końcówce nazwy. Rdzeń nie
// zgaduje formatu po zawartości: plik bez rozszerzenia jest brakiem wskazania,
// a wskazanie z żądania (`format`) i tak ma pierwszeństwo.
func formatDokumentuZeSciezki(sciezka string) (shared.TranslationDocumentFormat, bool) {
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".docx":
		return shared.TranslationDocumentFormatDocx, true
	case ".pdf":
		return shared.TranslationDocumentFormatPdf, true
	case ".pptx":
		return shared.TranslationDocumentFormatPptx, true
	case ".xlsx":
		return shared.TranslationDocumentFormatXlsx, true
	case ".odt":
		return shared.TranslationDocumentFormatOdt, true
	case ".md", ".markdown":
		return shared.TranslationDocumentFormatMarkdown, true
	case ".html", ".htm":
		return shared.TranslationDocumentFormatHtml, true
	}
	return "", false
}

// segmentyZDokumentu wyjmuje z pliku akapity wraz z miejscem w strukturze.
// Liczba stron wraca osobno — zna ją wyłącznie PDF; pozostałe formaty oddają
// zero, bo strona jest w nich wynikiem składu, a nie własnością pliku.
func segmentyZDokumentu(sciezka string,
	format shared.TranslationDocumentFormat) ([]dane.SegmentDokumentu, int64, error) {

	bajty, err := os.ReadFile(filepath.Clean(sciezka))
	if err != nil {
		return nil, 0, bladPlikuTlumaczenia(sciezka, err)
	}

	switch format {
	case shared.TranslationDocumentFormatMarkdown, shared.TranslationDocumentFormatHtml:
		tekst := string(bajty)
		if format == shared.TranslationDocumentFormatHtml {
			tekst = tekstZHtml(tekst)
		}
		return segmentyZAkapitow(rozdzielAkapity(tekst), "akapit"), 0, nil

	case shared.TranslationDocumentFormatDocx:
		akapity, err := akapityZArchiwum(bajty, "word/document.xml", "w:p")
		if err != nil {
			return nil, 0, err
		}
		return segmentyZAkapitow(akapity, "w:p"), 0, nil

	case shared.TranslationDocumentFormatOdt:
		akapity, err := akapityZArchiwum(bajty, "content.xml", "text:p")
		if err != nil {
			return nil, 0, err
		}
		return segmentyZAkapitow(akapity, "text:p"), 0, nil

	case shared.TranslationDocumentFormatPptx:
		akapity, err := akapityZeSlajdow(bajty)
		if err != nil {
			return nil, 0, err
		}
		return segmentyZAkapitow(akapity, "a:p"), 0, nil

	case shared.TranslationDocumentFormatXlsx:
		akapity, err := akapityZArchiwum(bajty, "xl/sharedStrings.xml", "si")
		if err != nil {
			return nil, 0, err
		}
		return segmentyZAkapitow(akapity, "si"), 0, nil

	case shared.TranslationDocumentFormatPdf:
		akapity, stron, err := akapityZPdf(bajty)
		if err != nil {
			return nil, 0, err
		}
		return segmentyZAkapitow(akapity, "strona"), stron, nil
	}
	return nil, 0, bladWskazaniaTlumaczenia("nieznany format dokumentu: " + string(format))
}

// segmentyZAkapitow nadaje akapitom kolejne numery i ścieżkę węzła w
// strukturze dokumentu, potrzebną przy składaniu wyniku.
func segmentyZAkapitow(akapity []string, wezel string) []dane.SegmentDokumentu {
	segmenty := make([]dane.SegmentDokumentu, 0, len(akapity))
	for numer, akapit := range akapity {
		tresc := strings.TrimSpace(akapit)
		if tresc == "" {
			continue
		}
		sciezka := wezel + "[" + strconv.Itoa(numer) + "]"
		segmenty = append(segmenty, dane.SegmentDokumentu{
			Kolejnosc:    int64(len(segmenty)),
			Tresc:        tresc,
			SciezkaWezla: &sciezka,
		})
	}
	return segmenty
}

// rozdzielAkapity dzieli tekst na akapity po pustej linii — jedyny podział,
// który w tekście płaskim jest widoczny i powtarzalny.
func rozdzielAkapity(tekst string) []string {
	znormalizowany := strings.ReplaceAll(tekst, "\r\n", "\n")
	return strings.Split(znormalizowany, "\n\n")
}

// tekstZHtml zdejmuje znaczniki, zostawiając treść. Prosty rozbiór strumieniem
// `encoding/xml` odmówiłby na pierwszym niezamkniętym `<br>`, a takie HTML-e
// przychodzą do tłumaczenia najczęściej.
func tekstZHtml(tresc string) string {
	var b strings.Builder
	wZnaczniku := false
	wSkrypcie := false
	for i := 0; i < len(tresc); i++ {
		if !wZnaczniku && strings.HasPrefix(strings.ToLower(tresc[i:]), "<script") {
			wSkrypcie = true
		}
		if wSkrypcie && strings.HasPrefix(strings.ToLower(tresc[i:]), "</script>") {
			wSkrypcie = false
			i += len("</script>") - 1
			continue
		}
		switch {
		case tresc[i] == '<':
			wZnaczniku = true
			// Znacznik blokowy jest granicą akapitu — inaczej cała strona
			// zlałaby się w jeden segment.
			if znacznikBlokowy(tresc[i:]) {
				b.WriteString("\n\n")
			}
		case tresc[i] == '>':
			wZnaczniku = false
		case !wZnaczniku && !wSkrypcie:
			b.WriteByte(tresc[i])
		}
	}
	return b.String()
}

// znacznikBlokowy mówi, czy znacznik zaczynający się w tym miejscu zamyka
// akapit, czyli jest granicą podziału tekstu.
func znacznikBlokowy(fragment string) bool {
	maly := strings.ToLower(fragment)
	for _, znacznik := range []string{"<p", "</p", "<div", "</div", "<br", "<li", "</li",
		"<h1", "<h2", "<h3", "<h4", "<tr", "</tr", "<td"} {
		if strings.HasPrefix(maly, znacznik) {
			return true
		}
	}
	return false
}

// akapityZArchiwum wyjmuje treść wskazanych węzłów z jednego pliku wewnątrz
// archiwum pakietu biurowego.
func akapityZArchiwum(bajty []byte, wpis, wezel string) ([]string, error) {
	archiwum, err := zip.NewReader(bytes.NewReader(bajty), int64(len(bajty)))
	if err != nil {
		return nil, bladWskazaniaTlumaczenia(
			"plik nie jest archiwum pakietu biurowego: " + err.Error())
	}
	for _, plik := range archiwum.File {
		if plik.Name != wpis {
			continue
		}
		strumien, err := plik.Open()
		if err != nil {
			return nil, bladTlumaczenia(err)
		}
		defer strumien.Close()
		tresc, err := io.ReadAll(strumien)
		if err != nil {
			return nil, bladTlumaczenia(err)
		}
		return akapityZXml(tresc, wezel), nil
	}
	return nil, bladWskazaniaTlumaczenia(
		"archiwum nie ma wpisu " + wpis + " — to nie jest dokument tego rodzaju")
}

// akapityZeSlajdow zbiera akapity ze wszystkich slajdów prezentacji,
// zachowując kolejność slajdów w pliku.
func akapityZeSlajdow(bajty []byte) ([]string, error) {
	archiwum, err := zip.NewReader(bytes.NewReader(bajty), int64(len(bajty)))
	if err != nil {
		return nil, bladWskazaniaTlumaczenia("plik nie jest prezentacją PPTX: " + err.Error())
	}
	akapity := []string{}
	for _, plik := range archiwum.File {
		if !strings.HasPrefix(plik.Name, "ppt/slides/slide") ||
			!strings.HasSuffix(plik.Name, ".xml") {
			continue
		}
		strumien, err := plik.Open()
		if err != nil {
			return nil, bladTlumaczenia(err)
		}
		tresc, err := io.ReadAll(strumien)
		strumien.Close()
		if err != nil {
			return nil, bladTlumaczenia(err)
		}
		akapity = append(akapity, akapityZXml(tresc, "a:p")...)
	}
	return akapity, nil
}

// akapityZXml przechodzi dokument XML strumieniem i skleja treść tekstową
// każdego wystąpienia wskazanego węzła. Strumień, nie model dokumentu: plik
// `word/document.xml` długiej umowy ma dziesiątki megabajtów.
func akapityZXml(tresc []byte, wezel string) []string {
	czytnik := xml.NewDecoder(bytes.NewReader(tresc))
	akapity := []string{}
	var biezacy strings.Builder
	glebokosc := 0
	for {
		znacznik, err := czytnik.Token()
		if err != nil {
			break
		}
		switch element := znacznik.(type) {
		case xml.StartElement:
			if nazwaWezla(element.Name) == wezel {
				glebokosc++
			}
		case xml.CharData:
			if glebokosc > 0 {
				biezacy.Write(element)
			}
		case xml.EndElement:
			if nazwaWezla(element.Name) != wezel {
				continue
			}
			glebokosc--
			if glebokosc == 0 {
				akapity = append(akapity, biezacy.String())
				biezacy.Reset()
			}
		}
	}
	return akapity
}

// nazwaWezla składa nazwę z przedrostkiem przestrzeni nazw w postaci, w jakiej
// stoi w pliku (`w:p`, `text:p`, `a:p`). `encoding/xml` rozwija przestrzeń nazw
// do pełnego adresu, więc przedrostka trzeba szukać po nazwie lokalnej i
// przestrzeni.
func nazwaWezla(nazwa xml.Name) string {
	switch {
	case strings.Contains(nazwa.Space, "wordprocessingml"):
		return "w:" + nazwa.Local
	case strings.Contains(nazwa.Space, "opendocument") && strings.Contains(nazwa.Space, "text"):
		return "text:" + nazwa.Local
	case strings.Contains(nazwa.Space, "drawingml"):
		return "a:" + nazwa.Local
	}
	return nazwa.Local
}

// akapityZPdf wyjmuje tekst z warstwy tekstowej dokumentu strona po stronie
// i oddaje liczbę stron. Rozbiór idzie po strumieniu treści z `pdfcpu`:
// operatory pokazania tekstu niosą napisy, reszta strumienia jest składem.
func akapityZPdf(bajty []byte) ([]string, int64, error) {
	nastawy := nastawyPdf()
	stron, err := api.PageCount(bytes.NewReader(bajty), nastawy)
	if err != nil {
		return nil, 0, bladWskazaniaTlumaczenia("nieczytelny dokument PDF: " + err.Error())
	}
	akapity := make([]string, 0, stron)
	err = api.ExtractContent(bytes.NewReader(bajty), nil, func(strumien io.Reader, _ int) error {
		tresc, err := io.ReadAll(strumien)
		if err != nil {
			return err
		}
		akapity = append(akapity, tekstZeStrumieniaPdf(string(tresc)))
		return nil
	}, nastawy)
	if err != nil {
		return nil, 0, bladWskazaniaTlumaczenia("nie można sięgnąć po treść dokumentu PDF: " + err.Error())
	}
	return akapity, int64(stron), nil
}

// tekstZeStrumieniaPdf wyjmuje napisy ze strumienia treści strony, pomijając
// operatory składu, których tłumaczenie nie dotyczy.
func tekstZeStrumieniaPdf(strumien string) string {
	var wynik strings.Builder
	var napis strings.Builder
	wNapisie := false
	poziom := 0
	for i := 0; i < len(strumien); i++ {
		znak := strumien[i]
		if wNapisie {
			switch znak {
			case '\\':
				if i+1 < len(strumien) {
					i++
					napis.WriteByte(strumien[i])
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
				wynik.WriteString(napis.String())
				napis.Reset()
			default:
				napis.WriteByte(znak)
			}
			continue
		}
		switch {
		case znak == '(':
			wNapisie = true
		case znak == 'T' && i+1 < len(strumien) && (strumien[i+1] == 'j' || strumien[i+1] == 'J'):
			wynik.WriteString(" ")
		case znak == '\n' || znak == '\r':
			wynik.WriteString(" ")
		}
	}
	return strings.Join(strings.Fields(wynik.String()), " ")
}

// zapiszDokumentWyniku składa plik wyniku z akapitów przekładu. Formaty
// tekstowe powstają wprost; DOCX powstaje jako poprawne archiwum OOXML.
// Formaty, których rdzeń nie umie złożyć bez utraty układu, są odmawiane
// po nazwie.
func zapiszDokumentWyniku(sciezka string, format shared.TranslationDocumentFormat,
	akapity []string) error {

	if err := os.MkdirAll(filepath.Dir(sciezka), 0o700); err != nil {
		return bladPlikuTlumaczenia(sciezka, err)
	}
	switch format {
	case shared.TranslationDocumentFormatMarkdown:
		return zapiszPlikWyniku(sciezka, []byte(strings.Join(akapity, "\n\n")+"\n"))
	case shared.TranslationDocumentFormatHtml:
		var b strings.Builder
		b.WriteString("<!doctype html>\n<html>\n<body>\n")
		for _, akapit := range akapity {
			b.WriteString("<p>" + zabezpieczHtml(akapit) + "</p>\n")
		}
		b.WriteString("</body>\n</html>\n")
		return zapiszPlikWyniku(sciezka, []byte(b.String()))
	case shared.TranslationDocumentFormatDocx:
		bajty, err := zlozDocx(akapity)
		if err != nil {
			return err
		}
		return zapiszPlikWyniku(sciezka, bajty)
	}
	return bladWskazaniaTlumaczenia("rdzeń składa wynik w formacie markdown, html albo docx; " +
		"format " + string(format) + " wymagałby odtworzenia składu, którego rdzeń nie prowadzi")
}

// zapiszPlikWyniku odkłada bajty wyniku pod wskazaną ścieżką na dysku,
// tworząc plik wynikowy przekładu.
func zapiszPlikWyniku(sciezka string, bajty []byte) error {
	if err := os.WriteFile(sciezka, bajty, 0o600); err != nil {
		return bladPlikuTlumaczenia(sciezka, err)
	}
	return nil
}

// zabezpieczHtml zamienia znaki o znaczeniu składniowym na encje — inaczej
// przekład z nawiasem ostrym rozsypałby wynikowy dokument.
func zabezpieczHtml(tekst string) string {
	zamiennik := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return zamiennik.Replace(tekst)
}

// zlozDocx składa najmniejszy poprawny dokument OOXML: wykaz typów treści
// i jeden dokument główny z akapitami.
func zlozDocx(akapity []string) ([]byte, error) {
	var bufor bytes.Buffer
	archiwum := zip.NewWriter(&bufor)

	typy := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`
	powiazania := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`

	var dokument strings.Builder
	dokument.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>`)
	for _, akapit := range akapity {
		dokument.WriteString(`<w:p><w:r><w:t xml:space="preserve">` +
			zabezpieczHtml(akapit) + `</w:t></w:r></w:p>`)
	}
	dokument.WriteString(`</w:body></w:document>`)

	for _, wpis := range []struct{ nazwa, tresc string }{
		{"[Content_Types].xml", typy},
		{"_rels/.rels", powiazania},
		{"word/document.xml", dokument.String()},
	} {
		strumien, err := archiwum.Create(wpis.nazwa)
		if err != nil {
			return nil, bladTlumaczenia(err)
		}
		if _, err := strumien.Write([]byte(wpis.tresc)); err != nil {
			return nil, bladTlumaczenia(err)
		}
	}
	if err := archiwum.Close(); err != nil {
		return nil, bladTlumaczenia(err)
	}
	return bufor.Bytes(), nil
}
