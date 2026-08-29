// Odpowiedzialność pliku: drobne rachunki wspólne dla wszystkich obszarów
// postaci dokumentu — kopie struktur postaci, ich scalanie, porównanie,
// wskaźniki na wartości i zapis zakresu w treści odmowy.
package core

import (
	"encoding/json"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// postacBladZaplecza nazywa brak po stronie montażu rdzenia — nie winę
// Operatora i nie usterkę przemijającą.
func postacBladZaplecza(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Studio, postać dokumentu: "+powod))
}

// postacWskaznikTekstu oddaje wskaźnik na napis — kontrakt trzyma pola
// nieobowiązkowe wskaźnikami, a zdania odmowy i noty bilansu są takimi polami.
func postacWskaznikTekstu(wartosc string) *string {
	if wartosc == "" {
		return nil
	}
	kopia := wartosc
	return &kopia
}

// postacWskaznikLiczby oddaje wskaźnik na kopię liczby całkowitej, tą samą
// zasadą co `postacWskaznikTekstu`: pole nieobowiązkowe kontraktu.
func postacWskaznikLiczby(wartosc int) *int {
	kopia := wartosc
	return &kopia
}

// postacWskaznikMiary oddaje wskaźnik na kopię liczby zmiennoprzecinkowej,
// tą samą zasadą co `postacWskaznikLiczby`.
func postacWskaznikMiary(wartosc float64) *float64 {
	kopia := wartosc
	return &kopia
}

// postacWskaznikPrawdy oddaje wskaźnik na kopię wartości logicznej, tą samą
// zasadą co `postacWskaznikLiczby`.
func postacWskaznikPrawdy(wartosc bool) *bool {
	kopia := wartosc
	return &kopia
}

// postacZapisLiczby zapisuje liczbę całkowitą napisem, do zdań odmowy i not
// bilansu czytelnych bez zaglądania do kodu.
func postacZapisLiczby(wartosc int) string {
	return strconv.Itoa(wartosc)
}

// postacZapisZakresu opisuje zakres w znakach zdaniem, które Operator rozumie
// bez zaglądania do kodu źródłowego.
func postacZapisZakresu(od, do int) string {
	if od == do {
		return "w miejscu znaku " + strconv.Itoa(od)
	}
	return "od znaku " + strconv.Itoa(od) + " do znaku " + strconv.Itoa(do)
}

// postacKopiaZnaku oddaje kopię postaci znaku. Bez kopii dwa fragmenty
// pracowałyby na jednej strukturze i pogrubienie jednego pogrubiałoby drugi.
func postacKopiaZnaku(wzor *shared.StudioCharacterFormat) *shared.StudioCharacterFormat {
	if wzor == nil {
		return nil
	}
	kopia := *wzor
	return &kopia
}

// postacKopiaAkapitu oddaje kopię postaci akapitu wraz z osobną tablicą
// tabulatorów — tablica wspólna dawałaby ten sam błąd co postać wspólna.
func postacKopiaAkapitu(wzor *shared.StudioParagraphFormat) *shared.StudioParagraphFormat {
	if wzor == nil {
		return nil
	}
	kopia := *wzor
	if wzor.TabStops != nil {
		kopia.TabStops = append([]shared.StudioTabStop(nil), wzor.TabStops...)
	}
	if wzor.Border != nil {
		obramowanie := *wzor.Border
		kopia.Border = &obramowanie
	}
	return &kopia
}

// postacTaSamaPostacZnaku porównuje dwie postacie znaku. Służy zbieraniu
// ciągów znaków w jeden fragment: bez tego drzewo rosłoby do jednego fragmentu
// na literę i zapis dokumentu byłby wielokrotnie większy od jego treści.
func postacTaSamaPostacZnaku(pierwsza, druga *shared.StudioCharacterFormat) bool {
	if pierwsza == nil || druga == nil {
		return pierwsza == nil && druga == nil
	}
	return *pierwsza == *druga
}

// postacTenSamAutor porównuje autorów fragmentu, wskaźniki puste licząc
// jako autora tego samego, a nie jako autorów różnych.
func postacTenSamAutor(pierwszy, drugi *shared.StudioAuthor) bool {
	if pierwszy == nil || drugi == nil {
		return pierwszy == nil && drugi == nil
	}
	return *pierwszy == *drugi
}

// postacScalZnak wnosi do postaci znaku wyłącznie pola podane w żądaniu,
// pozostałe biorąc z postaci zastanej.
func postacScalZnak(zastana *shared.StudioCharacterFormat,
	zmiana shared.StudioCharacterFormat) *shared.StudioCharacterFormat {

	wynik := shared.StudioCharacterFormat{}
	if zastana != nil {
		wynik = *zastana
	}
	if zmiana.FontFamily != nil {
		wynik.FontFamily = zmiana.FontFamily
	}
	if zmiana.FontSizePt != nil {
		wynik.FontSizePt = zmiana.FontSizePt
	}
	if zmiana.Bold != nil {
		wynik.Bold = zmiana.Bold
	}
	if zmiana.Italic != nil {
		wynik.Italic = zmiana.Italic
	}
	if zmiana.Underline != nil {
		wynik.Underline = zmiana.Underline
	}
	if zmiana.Strikethrough != nil {
		wynik.Strikethrough = zmiana.Strikethrough
	}
	if zmiana.Superscript != nil {
		wynik.Superscript = zmiana.Superscript
		// Indeks górny i dolny wykluczają się wzajemnie.
		if *zmiana.Superscript {
			wynik.Subscript = postacWskaznikPrawdy(false)
		}
	}
	if zmiana.Subscript != nil {
		wynik.Subscript = zmiana.Subscript
		if *zmiana.Subscript {
			wynik.Superscript = postacWskaznikPrawdy(false)
		}
	}
	if zmiana.Color != nil {
		wynik.Color = zmiana.Color
	}
	if zmiana.HighlightColor != nil {
		// Barwa pusta ZDEJMUJE wyróżnienie — tak mówi kontrakt tego pola.
		if *zmiana.HighlightColor == "" {
			wynik.HighlightColor = nil
		} else {
			wynik.HighlightColor = zmiana.HighlightColor
		}
	}
	if zmiana.LetterSpacingPt != nil {
		wynik.LetterSpacingPt = zmiana.LetterSpacingPt
	}
	if zmiana.SmallCaps != nil {
		wynik.SmallCaps = zmiana.SmallCaps
	}
	if zmiana.AllCaps != nil {
		wynik.AllCaps = zmiana.AllCaps
	}
	if zmiana.Effect != nil {
		wynik.Effect = zmiana.Effect
	}
	if zmiana.StyleName != nil {
		wynik.StyleName = zmiana.StyleName
	}
	if zmiana.Language != nil {
		wynik.Language = zmiana.Language
	}
	return &wynik
}

// postacScalAkapit wnosi do postaci akapitu wyłącznie pola podane w żądaniu,
// tą samą zasadą co `postacScalZnak`.
func postacScalAkapit(zastana *shared.StudioParagraphFormat,
	zmiana shared.StudioParagraphFormat) *shared.StudioParagraphFormat {

	wynik := shared.StudioParagraphFormat{}
	if zastana != nil {
		wynik = *postacKopiaAkapitu(zastana)
	}
	if zmiana.Align != nil {
		wynik.Align = zmiana.Align
	}
	if zmiana.FirstLineIndentMm != nil {
		wynik.FirstLineIndentMm = zmiana.FirstLineIndentMm
	}
	if zmiana.IndentLeftMm != nil {
		wynik.IndentLeftMm = zmiana.IndentLeftMm
	}
	if zmiana.IndentRightMm != nil {
		wynik.IndentRightMm = zmiana.IndentRightMm
	}
	if zmiana.SpaceBeforePt != nil {
		wynik.SpaceBeforePt = zmiana.SpaceBeforePt
	}
	if zmiana.SpaceAfterPt != nil {
		wynik.SpaceAfterPt = zmiana.SpaceAfterPt
	}
	if zmiana.LineSpacingRule != nil {
		wynik.LineSpacingRule = zmiana.LineSpacingRule
	}
	if zmiana.LineSpacingValue != nil {
		wynik.LineSpacingValue = zmiana.LineSpacingValue
	}
	if zmiana.TabStops != nil {
		wynik.TabStops = append([]shared.StudioTabStop(nil), zmiana.TabStops...)
	}
	if zmiana.Border != nil {
		wynik.Border = zmiana.Border
	}
	if zmiana.ShadingColor != nil {
		if *zmiana.ShadingColor == "" {
			wynik.ShadingColor = nil
		} else {
			wynik.ShadingColor = zmiana.ShadingColor
		}
	}
	if zmiana.WidowControl != nil {
		wynik.WidowControl = zmiana.WidowControl
	}
	if zmiana.KeepWithNext != nil {
		wynik.KeepWithNext = zmiana.KeepWithNext
	}
	if zmiana.KeepLines != nil {
		wynik.KeepLines = zmiana.KeepLines
	}
	if zmiana.OutlineLevel != nil {
		wynik.OutlineLevel = zmiana.OutlineLevel
	}
	if zmiana.StyleName != nil {
		wynik.StyleName = zmiana.StyleName
	}
	if zmiana.ListId != nil {
		wynik.ListId = zmiana.ListId
	}
	if zmiana.ListLevel != nil {
		wynik.ListLevel = zmiana.ListLevel
	}
	if zmiana.RightToLeft != nil {
		wynik.RightToLeft = zmiana.RightToLeft
	}
	return &wynik
}

// postacPolaNiejednolite nazywa pola, które we fragmencie mają więcej niż jedną
// wartość. Okno pokazuje je jako nastawę nieokreśloną — kłamstwem byłoby
// pokazanie pierwszej napotkanej wartości jako wartości całego zaznaczenia.
func postacPolaNiejednolite(runy []shared.StudioDocumentRun) []string {
	if len(runy) < 2 {
		return []string{}
	}
	wzor := shared.StudioCharacterFormat{}
	if runy[0].Format != nil {
		wzor = *runy[0].Format
	}
	znaleziono := map[string]bool{}
	for _, run := range runy[1:] {
		biezaca := shared.StudioCharacterFormat{}
		if run.Format != nil {
			biezaca = *run.Format
		}
		if !postacTaSamaMiara(wzor.FontSizePt, biezaca.FontSizePt) {
			znaleziono["fontSizePt"] = true
		}
		if !postacTenSamTekst(wzor.FontFamily, biezaca.FontFamily) {
			znaleziono["fontFamily"] = true
		}
		if !postacTaSamaPrawda(wzor.Bold, biezaca.Bold) {
			znaleziono["bold"] = true
		}
		if !postacTaSamaPrawda(wzor.Italic, biezaca.Italic) {
			znaleziono["italic"] = true
		}
		if !postacTaSamaPrawda(wzor.Strikethrough, biezaca.Strikethrough) {
			znaleziono["strikethrough"] = true
		}
		if !postacTenSamTekst(wzor.Color, biezaca.Color) {
			znaleziono["color"] = true
		}
		if !postacTenSamTekst(wzor.HighlightColor, biezaca.HighlightColor) {
			znaleziono["highlightColor"] = true
		}
		if !postacTenSamTekst(wzor.StyleName, biezaca.StyleName) {
			znaleziono["styleName"] = true
		}
	}
	rozne := make([]string, 0, len(znaleziono))
	for _, nazwa := range []string{"fontFamily", "fontSizePt", "bold", "italic",
		"strikethrough", "color", "highlightColor", "styleName"} {
		if znaleziono[nazwa] {
			rozne = append(rozne, nazwa)
		}
	}
	return rozne
}

func postacTenSamTekst(pierwszy, drugi *string) bool {
	if pierwszy == nil || drugi == nil {
		return pierwszy == nil && drugi == nil
	}
	return *pierwszy == *drugi
}

func postacTaSamaMiara(pierwsza, druga *float64) bool {
	if pierwsza == nil || druga == nil {
		return pierwsza == nil && druga == nil
	}
	return *pierwsza == *druga
}

func postacTaSamaPrawda(pierwsza, druga *bool) bool {
	if pierwsza == nil || druga == nil {
		return pierwsza == nil && druga == nil
	}
	return *pierwsza == *druga
}

// postacWskaznikLiczby64 oddaje wskaźnik na liczbę sześćdziesięcioczterobitową —
// warstwa danych trzyma zakresy w znakach właśnie tym typem.
func postacWskaznikLiczby64(wartosc int64) *int64 {
	kopia := wartosc
	return &kopia
}

// Dostęp do warstwy danych
// postacSkladnica oddaje tabele obszaru postaci dokumentu przez kontrakt
// `dane.RepozytoriumPostaciStudia`, a nie przez całe repozytorium.
func (a *adapterStudia) postacSkladnica() (dane.RepozytoriumPostaciStudia, error) {
	if a == nil || a.repozytorium == nil {
		return nil, postacBladZaplecza("repozytorium Studia nie zostało podane przy montażu serwera")
	}
	return a.repozytorium, nil
}

// Nastawy strony domyślne
// postacNastawyDomyslne oddaje nastawy strony dokumentu, który swoich nie ma.
// Nastawa jedzie z jednego miejsca — `wejscieDomyslneNastawyStrony`.
func postacNastawyDomyslne() *shared.StudioPageSetup {
	nastawy := wejscieDomyslneNastawyStrony("A4", nil)
	return &nastawy
}

// Składanie postaci z wierszy warstwy danych
// postacZlozSekcje składa sekcje dokumentu z wierszy, kolumną nadpisując pole
// JSON, gdy obie niosą tę samą wartość — kolumna jest prawdą, pole JSON tylko
// tym, na co kolumny nie ma.
func postacZlozSekcje(wiersze []dane.SekcjaDokumentuStudia) []shared.StudioSection {
	sekcje := make([]shared.StudioSection, 0, len(wiersze))
	for _, wiersz := range wiersze {
		sekcja := shared.StudioSection{
			Id:         wiersz.Kod,
			Index:      int(wiersz.Kolejnosc),
			Title:      wiersz.Tytul,
			RangeStart: int(wiersz.ZakresOd),
			RangeEnd:   int(wiersz.ZakresDo),
		}
		if wiersz.Rozpoczecie != "" {
			rozpoczecie := shared.StudioSectionStart(wiersz.Rozpoczecie)
			sekcja.Start = &rozpoczecie
		}
		if wiersz.NastawyStronyJSON != nil && *wiersz.NastawyStronyJSON != "" {
			var nastawy shared.StudioPageSetup
			if json.Unmarshal([]byte(*wiersz.NastawyStronyJSON), &nastawy) == nil {
				sekcja.PageSetup = &nastawy
			}
		}
		if wiersz.NaglowkiJSON != nil && *wiersz.NaglowkiJSON != "" {
			var naglowki []shared.StudioHeaderFooter
			if json.Unmarshal([]byte(*wiersz.NaglowkiJSON), &naglowki) == nil {
				sekcja.HeadersFooters = naglowki
			}
		}
		if wiersz.NumeracjaJSON != nil && *wiersz.NumeracjaJSON != "" {
			var numeracja shared.StudioPageNumbering
			if json.Unmarshal([]byte(*wiersz.NumeracjaJSON), &numeracja) == nil {
				sekcja.Numbering = &numeracja
			}
		}
		if wiersz.ZnakWodnyJSON != nil && *wiersz.ZnakWodnyJSON != "" {
			var znak shared.StudioWatermark
			if json.Unmarshal([]byte(*wiersz.ZnakWodnyJSON), &znak) == nil {
				sekcja.Watermark = &znak
			}
		}
		sekcje = append(sekcje, sekcja)
	}
	return sekcje
}

// postacZlozObiekty składa obiekty osadzone z wierszy, tą samą zasadą co
// `postacZlozSekcje`: kolumna nadpisuje pole JSON.
func postacZlozObiekty(wiersze []dane.ObiektDokumentuStudia) []shared.StudioDocumentObject {
	obiekty := make([]shared.StudioDocumentObject, 0, len(wiersze))
	for _, wiersz := range wiersze {
		obiekt := shared.StudioDocumentObject{}
		if wiersz.PostacJSON != nil && *wiersz.PostacJSON != "" {
			_ = json.Unmarshal([]byte(*wiersz.PostacJSON), &obiekt)
		}
		obiekt.Id = wiersz.Kod
		obiekt.Kind = shared.StudioObjectKind(wiersz.Rodzaj)
		obiekt.AssetId = wiersz.ZasobKod
		obiekt.DesignNodeId = wiersz.DesignWezelKod
		obiekt.LibraryFileId = wiersz.BibliotekaPlikKod
		obiekt.SourceUrl = wiersz.AdresZrodla
		obiekt.AltText = wiersz.TekstZastepczy
		obiekt.AnchorOffset = postacWskaznikLiczby(int(wiersz.ZakotwiczeniePozycja))
		obiekt.ZOrder = postacWskaznikLiczby(int(wiersz.Warstwa))
		if wiersz.Zrodlo != nil && *wiersz.Zrodlo != "" {
			zrodlo := shared.StudioObjectSource(*wiersz.Zrodlo)
			obiekt.Source = &zrodlo
		}
		if wiersz.Zakotwiczenie != "" {
			zakotwiczenie := shared.StudioAnchorKind(wiersz.Zakotwiczenie)
			obiekt.Anchor = &zakotwiczenie
		}
		obiekty = append(obiekty, obiekt)
	}
	return obiekty
}

// postacZlozAparat składa aparat dokumentu z wierszy, tą samą zasadą co
// `postacZlozSekcje`: kolumna nadpisuje pole JSON.
func postacZlozAparat(wiersze []dane.ElementAparatuStudia) []shared.StudioApparatusItem {
	aparat := make([]shared.StudioApparatusItem, 0, len(wiersze))
	for _, wiersz := range wiersze {
		element := shared.StudioApparatusItem{}
		if wiersz.DaneJSON != nil && *wiersz.DaneJSON != "" {
			_ = json.Unmarshal([]byte(*wiersz.DaneJSON), &element)
		}
		element.Id = wiersz.Kod
		element.Kind = shared.StudioApparatusKind(wiersz.Rodzaj)
		element.AnchorStart = postacWskaznikLiczby(int(wiersz.KotwicaOd))
		element.AnchorEnd = postacWskaznikLiczby(int(wiersz.KotwicaDo))
		element.Number = wiersz.Numer
		element.Label = wiersz.Etykieta
		element.Text = wiersz.Tresc
		element.TargetId = wiersz.CelKod
		element.TargetUrl = wiersz.CelAdres
		element.Stale = postacWskaznikPrawdy(wiersz.Nieswiezy)
		aparat = append(aparat, element)
	}
	return aparat
}

// postacZlozPola składa pola dokumentu z wierszy, tą samą zasadą co
// `postacZlozSekcje`: kolumna nadpisuje pole JSON.
func postacZlozPola(wiersze []dane.PoleDokumentuStudia) []shared.StudioDocumentField {
	pola := make([]shared.StudioDocumentField, 0, len(wiersze))
	for _, wiersz := range wiersze {
		pola = append(pola, shared.StudioDocumentField{
			Id:           wiersz.Kod,
			Kind:         shared.StudioFieldKind(wiersz.Rodzaj),
			AnchorOffset: postacWskaznikLiczby(int(wiersz.Kotwica)),
			Format:       wiersz.Format,
			Expression:   wiersz.Wyrazenie,
			PropertyName: wiersz.NazwaWlasciwosci,
			Value:        wiersz.Wartosc,
			Stale:        postacWskaznikPrawdy(wiersz.Nieswieze),
		})
	}
	return pola
}

// postacZlozBlokady składa blokady fragmentów widziane przez postać dokumentu,
// przez `blokadaZlozKontrakt` z obszaru kontroli pracy.
func postacZlozBlokady(dokumentKod string,
	wiersze []dane.BlokadaFragmentuStudia) []shared.StudioFragmentLock {

	blokady := make([]shared.StudioFragmentLock, 0, len(wiersze))
	for _, wiersz := range wiersze {
		blokada := blokadaZlozKontrakt(wiersz)
		if blokada.DocumentId == "" {
			blokada.DocumentId = dokumentKod
		}
		blokady = append(blokady, blokada)
	}
	return blokady
}
