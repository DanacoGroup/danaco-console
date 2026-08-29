// Odpowiedzialność pliku: arkusz stylów nazwanych dokumentu — wykaz,
// zakładanie i zmiana, stosowanie do zakresu, usunięcie z przeniesieniem miejsc
// użycia, dziedziczenie po stylu nadrzędnym i postać skuteczna.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// postacGlebokoscDziedziczenia zatrzymuje łańcuch stylów nadrzędnych, żeby
// arkusz z cyklem dziedziczenia nie zawiesił rachunku w pętli bez ogranicznika.
const postacGlebokoscDziedziczenia = 10

// postacStyleFabryczne oddaje arkusz stylów, który dostaje każdy dokument:
// nagłówki sześciu poziomów, tekst zasadniczy, cytat, podpis i przypis. Wykaz
// jest tym, co Operator zna z pakietu biurowego — nazwy pełne, bez ani jednego
// wymyślonego kodu.
func postacStyleFabryczne() []shared.StudioNamedStyle {
	zasadniczy := shared.StudioNamedStyle{
		Name: "Tekst zasadniczy", Kind: shared.StudioStyleKindParagraph,
		Builtin: postacWskaznikPrawdy(true),
		Character: &shared.StudioCharacterFormat{
			FontFamily: postacWskaznikTekstu("Liberation Serif"),
			FontSizePt: postacWskaznikMiary(11),
		},
		Paragraph: &shared.StudioParagraphFormat{
			Align:            postacWskaznikWyrownania(shared.StudioTextAlignJustify),
			SpaceAfterPt:     postacWskaznikMiary(6),
			LineSpacingRule:  postacWskaznikInterlinii(shared.StudioLineSpacingRuleSingle),
			LineSpacingValue: postacWskaznikMiary(1),
			WidowControl:     postacWskaznikPrawdy(true),
			OutlineLevel:     postacWskaznikLiczby(0),
		},
	}
	style := []shared.StudioNamedStyle{zasadniczy}

	// Stopnie nagłówków maleją wraz z poziomem, wzorem pakietu biurowego.
	stopnie := []float64{20, 17, 14, 12.5, 11.5, 11}
	for poziom := 1; poziom <= 6; poziom++ {
		nazwa := "Nagłówek poziomu " + string(rune('0'+poziom))
		style = append(style, shared.StudioNamedStyle{
			Name: nazwa, Kind: shared.StudioStyleKindParagraph,
			BasedOn: postacWskaznikTekstu(zasadniczy.Name),
			// Po nagłówku Operator pisze tekst, nie następny nagłówek.
			NextStyle: postacWskaznikTekstu(zasadniczy.Name),
			Builtin:   postacWskaznikPrawdy(true),
			Character: &shared.StudioCharacterFormat{
				FontFamily: postacWskaznikTekstu("Liberation Sans"),
				FontSizePt: postacWskaznikMiary(stopnie[poziom-1]),
				Bold:       postacWskaznikPrawdy(true),
			},
			Paragraph: &shared.StudioParagraphFormat{
				Align:         postacWskaznikWyrownania(shared.StudioTextAlignLeft),
				SpaceBeforePt: postacWskaznikMiary(12),
				SpaceAfterPt:  postacWskaznikMiary(6),
				KeepWithNext:  postacWskaznikPrawdy(true),
				KeepLines:     postacWskaznikPrawdy(true),
				OutlineLevel:  postacWskaznikLiczby(poziom),
			},
		})
	}

	style = append(style,
		shared.StudioNamedStyle{
			Name: "Cytat", Kind: shared.StudioStyleKindParagraph,
			BasedOn: postacWskaznikTekstu(zasadniczy.Name),
			Builtin: postacWskaznikPrawdy(true),
			Character: &shared.StudioCharacterFormat{
				Italic: postacWskaznikPrawdy(true),
			},
			Paragraph: &shared.StudioParagraphFormat{
				IndentLeftMm:  postacWskaznikMiary(12),
				IndentRightMm: postacWskaznikMiary(12),
				SpaceBeforePt: postacWskaznikMiary(6),
				SpaceAfterPt:  postacWskaznikMiary(6),
			},
		},
		shared.StudioNamedStyle{
			Name: "Podpis", Kind: shared.StudioStyleKindParagraph,
			BasedOn: postacWskaznikTekstu(zasadniczy.Name),
			Builtin: postacWskaznikPrawdy(true),
			Character: &shared.StudioCharacterFormat{
				FontSizePt: postacWskaznikMiary(9),
				Italic:     postacWskaznikPrawdy(true),
			},
			Paragraph: &shared.StudioParagraphFormat{
				Align:        postacWskaznikWyrownania(shared.StudioTextAlignCenter),
				SpaceAfterPt: postacWskaznikMiary(10),
			},
		},
		shared.StudioNamedStyle{
			Name: "Przypis", Kind: shared.StudioStyleKindParagraph,
			BasedOn: postacWskaznikTekstu(zasadniczy.Name),
			Builtin: postacWskaznikPrawdy(true),
			Character: &shared.StudioCharacterFormat{
				FontSizePt: postacWskaznikMiary(9),
			},
			Paragraph: &shared.StudioParagraphFormat{
				Align:        postacWskaznikWyrownania(shared.StudioTextAlignLeft),
				SpaceAfterPt: postacWskaznikMiary(0),
			},
		},
		shared.StudioNamedStyle{
			Name: "Wyróżnienie", Kind: shared.StudioStyleKindCharacter,
			Builtin: postacWskaznikPrawdy(true),
			Character: &shared.StudioCharacterFormat{
				Bold: postacWskaznikPrawdy(true),
			},
		},
		shared.StudioNamedStyle{
			Name: "Tabela zwykła", Kind: shared.StudioStyleKindTable,
			Builtin: postacWskaznikPrawdy(true),
		},
	)
	return style
}

func postacWskaznikWyrownania(wartosc shared.StudioTextAlign) *shared.StudioTextAlign {
	kopia := wartosc
	return &kopia
}

func postacWskaznikInterlinii(wartosc shared.StudioLineSpacingRule) *shared.StudioLineSpacingRule {
	kopia := wartosc
	return &kopia
}

// postacStylDoWiersza przekłada styl kontraktu na wiersz warstwy danych,
// serializując postać znaku i akapitu do zapisu JSON kolumny.
func postacStylDoWiersza(dokumentID int64,
	styl shared.StudioNamedStyle) (dane.StylNazwanyStudia, error) {

	wiersz := dane.StylNazwanyStudia{
		DokumentID:    dokumentID,
		Nazwa:         strings.TrimSpace(styl.Name),
		NazwaWidoczna: styl.DisplayName,
		Rodzaj:        string(styl.Kind),
		StylNadrzedny: styl.BasedOn,
		StylNastepny:  styl.NextStyle,
		Fabryczny:     styl.Builtin != nil && *styl.Builtin,
	}
	if wiersz.Rodzaj == "" {
		wiersz.Rodzaj = string(shared.StudioStyleKindParagraph)
	}
	if styl.Character != nil {
		zapis, err := json.Marshal(styl.Character)
		if err != nil {
			return dane.StylNazwanyStudia{}, postacBladZaplecza(
				"postaci znaku stylu " + styl.Name + " nie da się zapisać: " + err.Error())
		}
		wiersz.PostacZnakuJSON = postacWskaznikTekstu(string(zapis))
	}
	if styl.Paragraph != nil {
		zapis, err := json.Marshal(styl.Paragraph)
		if err != nil {
			return dane.StylNazwanyStudia{}, postacBladZaplecza(
				"postaci akapitu stylu " + styl.Name + " nie da się zapisać: " + err.Error())
		}
		wiersz.PostacAkapituJSON = postacWskaznikTekstu(string(zapis))
	}
	return wiersz, nil
}

// postacZlozStyle składa arkusz stylów z wierszy warstwy danych, odczytując
// z powrotem to, co postacStylDoWiersza zapisał.
func postacZlozStyle(wiersze []dane.StylNazwanyStudia) []shared.StudioNamedStyle {
	style := make([]shared.StudioNamedStyle, 0, len(wiersze))
	for _, wiersz := range wiersze {
		styl := shared.StudioNamedStyle{
			Name:        wiersz.Nazwa,
			DisplayName: wiersz.NazwaWidoczna,
			Kind:        shared.StudioStyleKind(wiersz.Rodzaj),
			BasedOn:     wiersz.StylNadrzedny,
			NextStyle:   wiersz.StylNastepny,
			Builtin:     postacWskaznikPrawdy(wiersz.Fabryczny),
			UsageCount:  postacWskaznikLiczby(0),
		}
		if wiersz.PostacZnakuJSON != nil && *wiersz.PostacZnakuJSON != "" {
			var znak shared.StudioCharacterFormat
			if json.Unmarshal([]byte(*wiersz.PostacZnakuJSON), &znak) == nil {
				styl.Character = &znak
			}
		}
		if wiersz.PostacAkapituJSON != nil && *wiersz.PostacAkapituJSON != "" {
			var akapit shared.StudioParagraphFormat
			if json.Unmarshal([]byte(*wiersz.PostacAkapituJSON), &akapit) == nil {
				styl.Paragraph = &akapit
			}
		}
		style = append(style, styl)
	}
	return style
}

// postacStyl znajduje styl w arkuszu po dokładnej nazwie, albo oddaje
// wskaźnik pusty, gdy arkusz stylu o takiej nazwie nie niesie.
func postacStyl(forma *shared.StudioDocumentForm, nazwa string) *shared.StudioNamedStyle {
	for i := range forma.Styles {
		if forma.Styles[i].Name == nazwa {
			return &forma.Styles[i]
		}
	}
	return nil
}

// postacStylFabryczny mówi, czy styl o tej nazwie jest fabryczny, oddając
// fałsz, gdy arkusz stylu o takiej nazwie nie niesie.
func postacStylFabryczny(style []shared.StudioNamedStyle, nazwa string) bool {
	for _, styl := range style {
		if styl.Name == nazwa {
			return styl.Builtin != nil && *styl.Builtin
		}
	}
	return false
}

// postacZnakStylu liczy postać znaku stylu wraz z dziedziczeniem po łańcuchu
// stylów nadrzędnych. Styl własny dokłada do nadrzędnego, nie zaczyna od zera.
func postacZnakStylu(forma *shared.StudioDocumentForm, nazwa string) *shared.StudioCharacterFormat {
	lancuch := postacLancuchStylow(forma, nazwa)
	var wynik *shared.StudioCharacterFormat
	for _, styl := range lancuch {
		if styl.Character != nil {
			wynik = postacScalZnak(wynik, *styl.Character)
		}
	}
	return wynik
}

// postacAkapitStylu liczy postać akapitu stylu wraz z dziedziczeniem po
// łańcuchu stylów nadrzędnych, stylem własnym dokładanym do nadrzędnego.
func postacAkapitStylu(forma *shared.StudioDocumentForm, nazwa string) *shared.StudioParagraphFormat {
	lancuch := postacLancuchStylow(forma, nazwa)
	var wynik *shared.StudioParagraphFormat
	for _, styl := range lancuch {
		if styl.Paragraph != nil {
			wynik = postacScalAkapit(wynik, *styl.Paragraph)
		}
	}
	return wynik
}

// postacLancuchStylow oddaje łańcuch dziedziczenia OD NAJSTARSZEGO nadrzędnego
// do stylu wskazanego — w tej kolejności scalanie daje pierwszeństwo stylowi
// najbliższemu, a nie najdalszemu.
func postacLancuchStylow(forma *shared.StudioDocumentForm, nazwa string) []shared.StudioNamedStyle {
	lancuch := make([]shared.StudioNamedStyle, 0, 4)
	odwiedzone := map[string]bool{}
	biezaca := strings.TrimSpace(nazwa)
	for krok := 0; krok < postacGlebokoscDziedziczenia && biezaca != ""; krok++ {
		if odwiedzone[biezaca] {
			break // cykl w arkuszu przejętym z szablonu; łańcuch się kończy
		}
		odwiedzone[biezaca] = true
		styl := postacStyl(forma, biezaca)
		if styl == nil {
			break
		}
		lancuch = append([]shared.StudioNamedStyle{*styl}, lancuch...)
		if styl.BasedOn == nil {
			break
		}
		biezaca = strings.TrimSpace(*styl.BasedOn)
	}
	return lancuch
}

// postacZnakSkuteczny liczy postać znaku, którą Operator naprawdę widzi:
// arkusz stylów pod spodem, postać własna fragmentu na wierzchu.
func postacZnakSkuteczny(forma *shared.StudioDocumentForm,
	wlasna *shared.StudioCharacterFormat) shared.StudioCharacterFormat {

	return postacZnakSkutecznyWBloku(forma, nil, wlasna)
}

// postacZnakSkutecznyWBloku liczy postać znaku fragmentu, scalając trzy
// warstwy w kolejności pierwszeństwa: postać znaku stylu akapitu, postać
// znaku stylu znaku nałożonego na fragment i postać własną fragmentu.
func postacZnakSkutecznyWBloku(forma *shared.StudioDocumentForm,
	blok *shared.StudioDocumentBlock,
	wlasna *shared.StudioCharacterFormat) shared.StudioCharacterFormat {

	nazwaAkapitu := "Tekst zasadniczy"
	if blok != nil && blok.Paragraph != nil && blok.Paragraph.StyleName != nil &&
		*blok.Paragraph.StyleName != "" {

		nazwaAkapitu = *blok.Paragraph.StyleName
	}
	wynik := postacZnakStylu(forma, nazwaAkapitu)
	if wlasna != nil && wlasna.StyleName != nil && *wlasna.StyleName != "" {
		if zeStylu := postacZnakStylu(forma, *wlasna.StyleName); zeStylu != nil {
			wynik = postacScalZnak(wynik, *zeStylu)
		}
	}
	if wlasna != nil {
		wynik = postacScalZnak(wynik, *wlasna)
	}
	if wynik == nil {
		return shared.StudioCharacterFormat{}
	}
	return *wynik
}

// postacAkapitSkuteczny liczy postać akapitu, którą Operator naprawdę widzi,
// scalając styl nazwany bloku z postacią akapitu ustawioną na samym bloku.
func postacAkapitSkuteczny(forma *shared.StudioDocumentForm,
	blok shared.StudioDocumentBlock) shared.StudioParagraphFormat {

	nazwaStylu := "Tekst zasadniczy"
	if blok.Paragraph != nil && blok.Paragraph.StyleName != nil && *blok.Paragraph.StyleName != "" {
		nazwaStylu = *blok.Paragraph.StyleName
	}
	wynik := postacAkapitStylu(forma, nazwaStylu)
	if blok.Paragraph != nil {
		wynik = postacScalAkapit(wynik, *blok.Paragraph)
	}
	if wynik == nil {
		return shared.StudioParagraphFormat{}
	}
	return *wynik
}

// postacPrzeliczUzycieStylow liczy, ile miejsc dokumentu używa każdego stylu.
// Liczba jedzie w odpowiedzi, bo galeria stylów ma pokazywać, czego dokument
// naprawdę używa, a `style.delete` — czego przeniesienie dotknie.
func postacPrzeliczUzycieStylow(forma *shared.StudioDocumentForm) {
	uzycie := map[string]int{}
	for _, blok := range forma.Blocks {
		if blok.Paragraph != nil && blok.Paragraph.StyleName != nil {
			uzycie[*blok.Paragraph.StyleName]++
		}
		for _, run := range blok.Runs {
			if run.Format != nil && run.Format.StyleName != nil {
				uzycie[*run.Format.StyleName]++
			}
		}
	}
	for _, tabela := range forma.Tables {
		if tabela.StyleName != nil {
			uzycie[*tabela.StyleName]++
		}
	}
	for i := range forma.Styles {
		forma.Styles[i].UsageCount = postacWskaznikLiczby(uzycie[forma.Styles[i].Name])
	}
}

// postacMiejscaUzyciaStylu zbiera miejsca użycia stylu wraz ze stylami, które po
// nim dziedziczą — zmiana stylu nadrzędnego przestawia też miejsca stylu
// potomnego i bilans musi je policzyć.
func postacMiejscaUzyciaStylu(forma *shared.StudioDocumentForm, nazwa string) int {
	rodzina := map[string]bool{nazwa: true}
	// Dziedziczenie bywa wielopoziomowe, więc rodzina domyka się przebiegami,
	// aż przestanie rosnąć.
	for krok := 0; krok < postacGlebokoscDziedziczenia; krok++ {
		przed := len(rodzina)
		for _, styl := range forma.Styles {
			if styl.BasedOn != nil && rodzina[*styl.BasedOn] {
				rodzina[styl.Name] = true
			}
		}
		if len(rodzina) == przed {
			break
		}
	}
	ile := 0
	for _, styl := range forma.Styles {
		if !rodzina[styl.Name] {
			continue
		}
		if styl.UsageCount != nil {
			ile += *styl.UsageCount
		}
	}
	return ile
}

// ── Czynności ───────────────────────────────────────────────────────────────

// StyleDokumentu oddaje arkusz stylów dokumentu (`studio.style.list`),
// zawężony rodzajem stylu albo wyłączeniem stylów fabrycznych.
func (a *adapterStudia) StyleDokumentu(ctx context.Context,
	z shared.StudioStyleListRequest) (shared.StudioStyleListResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioStyleListResponse{}, err
	}
	wybrane := make([]shared.StudioNamedStyle, 0, len(stan.forma.Styles))
	for _, styl := range stan.forma.Styles {
		if z.Kind != nil && *z.Kind != "" && styl.Kind != *z.Kind {
			continue
		}
		fabryczny := styl.Builtin != nil && *styl.Builtin
		if z.IncludeBuiltin != nil && !*z.IncludeBuiltin && fabryczny {
			continue
		}
		wybrane = append(wybrane, styl)
	}
	return shared.StudioStyleListResponse{Styles: wybrane}, nil
}

// ZapiszStyl zakłada styl własny albo zmienia zastany (`studio.style.save`),
// odmawiając nadpisania stylu fabrycznego. Bilans oddaje liczbę miejsc
// dokumentu, które zmiana stylu dotknęła.
func (a *adapterStudia) ZapiszStyl(ctx context.Context,
	z shared.StudioStyleSaveRequest) (shared.StudioStyleSaveResponse, error) {

	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.StudioStyleSaveResponse{}, bladWskazaniaStudio("zapis stylu bez nazwy")
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioStyleSaveResponse{}, err
	}
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioStyleSaveResponse{}, err
	}
	autor := postacAutor(z.Author)

	if postacStylFabryczny(stan.forma.Styles, nazwa) {
		return shared.StudioStyleSaveResponse{}, bladWskazaniaStudio(
			"styl „" + nazwa + "” jest fabryczny i nie da się go nadpisać. " +
				"Załóż styl własny dziedziczący po nim polem basedOn — wtedy zmiana " +
				"obejmie tylko te miejsca, które go używają")
	}

	styl := shared.StudioNamedStyle{Name: nazwa, Kind: shared.StudioStyleKindParagraph}
	if zastany := postacStyl(&stan.forma, nazwa); zastany != nil {
		styl = *zastany
	}
	if z.Kind != nil && *z.Kind != "" {
		styl.Kind = *z.Kind
	}
	if z.DisplayName != nil {
		styl.DisplayName = z.DisplayName
	}
	if z.BasedOn != nil {
		podstawa := strings.TrimSpace(*z.BasedOn)
		if podstawa != "" && postacStyl(&stan.forma, podstawa) == nil {
			return shared.StudioStyleSaveResponse{}, bladWskazaniaStudio(
				"styl nadrzędny „" + podstawa + "” nie istnieje w arkuszu tego dokumentu")
		}
		if podstawa == nazwa {
			return shared.StudioStyleSaveResponse{}, bladWskazaniaStudio(
				"styl „" + nazwa + "” nie może dziedziczyć po sobie samym")
		}
		styl.BasedOn = postacWskaznikTekstu(podstawa)
	}
	if z.NextStyle != nil {
		styl.NextStyle = z.NextStyle
	}
	if len(z.Character) > 0 {
		var znak shared.StudioCharacterFormat
		if err := json.Unmarshal(z.Character, &znak); err != nil {
			return shared.StudioStyleSaveResponse{}, bladWskazaniaStudio(
				"postać znaku stylu jest nieczytelna: " + err.Error())
		}
		styl.Character = postacScalZnak(styl.Character, znak)
	}
	if len(z.Paragraph) > 0 {
		var akapit shared.StudioParagraphFormat
		if err := json.Unmarshal(z.Paragraph, &akapit); err != nil {
			return shared.StudioStyleSaveResponse{}, bladWskazaniaStudio(
				"postać akapitu stylu jest nieczytelna: " + err.Error())
		}
		styl.Paragraph = postacScalAkapit(styl.Paragraph, akapit)
	}
	styl.Builtin = postacWskaznikPrawdy(false)

	wiersz, err := postacStylDoWiersza(stan.dokument.ID, styl)
	if err != nil {
		return shared.StudioStyleSaveResponse{}, err
	}
	if _, err := skladnica.ZapiszStylNazwany(ctx, wiersz); err != nil {
		return shared.StudioStyleSaveResponse{}, bladStudio(err)
	}
	odczytane, err := skladnica.StyleNazwane(ctx, stan.dokument.ID, "")
	if err != nil {
		return shared.StudioStyleSaveResponse{}, bladStudio(err)
	}
	stan.forma.Styles = postacZlozStyle(odczytane)
	postacPrzeliczUzycieStylow(&stan.forma)

	miejsca := postacMiejscaUzyciaStylu(&stan.forma, nazwa)
	bilans := shared.StudioActionBalance{
		Applied: miejsca + 1,
		Skipped: []shared.StudioSkippedItem{},
		Note: postacWskaznikTekstu("styl „" + nazwa + "” zapisany; przestawił " +
			postacLiczebnik(miejsca, "miejsce użycia", "miejsca użycia", "miejsc użycia") +
			" w dokumencie"),
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindStyleChange, 0, postacDlugosc(&stan.forma), bilans)
	if err != nil {
		return shared.StudioStyleSaveResponse{}, err
	}
	zapisany := postacStyl(&forma, nazwa)
	if zapisany == nil {
		return shared.StudioStyleSaveResponse{}, postacBladZaplecza(
			"styl " + nazwa + " zapisany, ale nie wrócił z arkusza")
	}
	return shared.StudioStyleSaveResponse{
		Style: *zapisany, Form: forma, Balance: bilansGotowy, Change: zmiana,
	}, nil
}

// ZastosujStyl przypisuje styl nazwany do zakresu (`studio.style.apply`):
// styl akapitu idzie na całe akapity, które zaznaczenie obejmuje, a styl
// znaku dokładnie na zaznaczone fragmenty.
func (a *adapterStudia) ZastosujStyl(ctx context.Context,
	z shared.StudioStyleApplyRequest) (shared.StudioStyleApplyResponse, error) {

	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.StudioStyleApplyResponse{}, bladWskazaniaStudio(
			"zastosowanie stylu bez nazwy stylu")
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioStyleApplyResponse{}, err
	}
	styl := postacStyl(&stan.forma, nazwa)
	if styl == nil {
		return shared.StudioStyleApplyResponse{}, bladWskazaniaStudio(
			"styl „" + nazwa + "” nie istnieje w arkuszu tego dokumentu")
	}
	autor := postacAutor(z.Author)
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)

	bilans := shared.StudioActionBalance{Skipped: pominiete}
	for _, odcinek := range odcinki {
		postacRozetnij(&stan.forma, odcinek[0], odcinek[1])
		switch styl.Kind {
		case shared.StudioStyleKindCharacter:
			for _, wskazanie := range postacFragmentyZakresu(&stan.forma, odcinek[0], odcinek[1]) {
				run := &stan.forma.Blocks[wskazanie[0]].Runs[wskazanie[1]]
				run.Format = postacScalZnak(run.Format,
					shared.StudioCharacterFormat{StyleName: postacWskaznikTekstu(nazwa)})
				bilans.Applied++
			}
		default:
			for _, wskazanie := range postacBlokiZakresu(&stan.forma, odcinek[0], odcinek[1]) {
				blok := &stan.forma.Blocks[wskazanie]
				blok.Paragraph = postacScalAkapit(blok.Paragraph,
					shared.StudioParagraphFormat{StyleName: postacWskaznikTekstu(nazwa)})
				// Styl nagłówkowy przestawia rodzaj bloku, nie tylko postać.
				poziom := 0
				if skuteczna := postacAkapitSkuteczny(&stan.forma, *blok); skuteczna.OutlineLevel != nil {
					poziom = *skuteczna.OutlineLevel
				}
				if poziom > 0 {
					blok.Kind = blokPostaciNaglowek
				} else if blok.Kind == blokPostaciNaglowek {
					blok.Kind = blokPostaciAkapit
				}
				bilans.Applied++
			}
		}
	}
	if bilans.Applied == 0 {
		return shared.StudioStyleApplyResponse{}, bladWskazaniaStudio(
			"styl „" + nazwa + "” nie miał na czym stanąć: zakres " +
				postacZapisZakresu(od, do) + " nie obejmuje ani jednego akapitu wolnego " +
				"od blokady")
	}
	bilans.Note = postacWskaznikTekstu("styl „" + nazwa + "” zastosowany " +
		postacZapisZakresu(od, do))

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindStyleChange, od, do, bilans)
	if err != nil {
		return shared.StudioStyleApplyResponse{}, err
	}
	return shared.StudioStyleApplyResponse{Form: forma, Balance: bilansGotowy, Change: zmiana}, nil
}

// UsunStyl usuwa styl własny i przenosi jego miejsca użycia
// (`studio.style.delete`), odmawiając usunięcia stylu fabrycznego. Miejsca
// użycia idą na styl wskazany polem `replaceWith`, a bez niego na tekst
// zasadniczy.
func (a *adapterStudia) UsunStyl(ctx context.Context,
	z shared.StudioStyleDeleteRequest) (shared.StudioStyleDeleteResponse, error) {

	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.StudioStyleDeleteResponse{}, bladWskazaniaStudio("usunięcie stylu bez nazwy")
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioStyleDeleteResponse{}, err
	}
	if postacStyl(&stan.forma, nazwa) == nil {
		return shared.StudioStyleDeleteResponse{}, bladWskazaniaStudio(
			"styl „" + nazwa + "” nie istnieje w arkuszu tego dokumentu")
	}
	if postacStylFabryczny(stan.forma.Styles, nazwa) {
		return shared.StudioStyleDeleteResponse{}, bladWskazaniaStudio(
			"styl „" + nazwa + "” jest fabryczny i nie da się go usunąć — " +
				"arkusz fabryczny jest wiedzą serwera, nie zapisem Operatora")
	}
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioStyleDeleteResponse{}, err
	}

	zamiennik := "Tekst zasadniczy"
	if z.ReplaceWith != nil && strings.TrimSpace(*z.ReplaceWith) != "" {
		zamiennik = strings.TrimSpace(*z.ReplaceWith)
		if postacStyl(&stan.forma, zamiennik) == nil {
			return shared.StudioStyleDeleteResponse{}, bladWskazaniaStudio(
				"styl zastępczy „" + zamiennik + "” nie istnieje w arkuszu tego dokumentu")
		}
	}

	przeniesione := 0
	for i := range stan.forma.Blocks {
		blok := &stan.forma.Blocks[i]
		if blok.Paragraph != nil && blok.Paragraph.StyleName != nil &&
			*blok.Paragraph.StyleName == nazwa {

			blok.Paragraph.StyleName = postacWskaznikTekstu(zamiennik)
			przeniesione++
		}
		for j := range blok.Runs {
			run := &blok.Runs[j]
			if run.Format != nil && run.Format.StyleName != nil && *run.Format.StyleName == nazwa {
				run.Format.StyleName = postacWskaznikTekstu(zamiennik)
				przeniesione++
			}
		}
	}
	for i := range stan.forma.Tables {
		if stan.forma.Tables[i].StyleName != nil && *stan.forma.Tables[i].StyleName == nazwa {
			stan.forma.Tables[i].StyleName = postacWskaznikTekstu(zamiennik)
			przeniesione++
		}
	}
	// Styl potomny traci nadrzędnego; bez przestawienia dziedziczyłby po nazwie,
	// której już nie ma.
	for i := range stan.forma.Styles {
		styl := &stan.forma.Styles[i]
		if styl.BasedOn == nil || *styl.BasedOn != nazwa {
			continue
		}
		styl.BasedOn = postacWskaznikTekstu(zamiennik)
		wiersz, err := postacStylDoWiersza(stan.dokument.ID, *styl)
		if err != nil {
			return shared.StudioStyleDeleteResponse{}, err
		}
		if _, err := skladnica.ZapiszStylNazwany(ctx, wiersz); err != nil {
			return shared.StudioStyleDeleteResponse{}, bladStudio(err)
		}
		przeniesione++
	}

	usuniety, err := skladnica.UsunStylNazwany(ctx, stan.dokument.ID, nazwa)
	if err != nil {
		return shared.StudioStyleDeleteResponse{}, bladStudio(err)
	}
	if !usuniety {
		return shared.StudioStyleDeleteResponse{}, postacBladZaplecza(
			"styl " + nazwa + " nie dał się usunąć, choć nie jest fabryczny")
	}
	odczytane, err := skladnica.StyleNazwane(ctx, stan.dokument.ID, "")
	if err != nil {
		return shared.StudioStyleDeleteResponse{}, bladStudio(err)
	}
	stan.forma.Styles = postacZlozStyle(odczytane)

	bilans := shared.StudioActionBalance{
		Applied: przeniesione + 1,
		Skipped: []shared.StudioSkippedItem{},
		Note: postacWskaznikTekstu("styl „" + nazwa + "” usunięty; " +
			postacLiczebnik(przeniesione, "miejsce użycia", "miejsca użycia", "miejsc użycia") +
			" przeniesione na styl „" + zamiennik + "”"),
	}
	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioStyleDeleteResponse{}, err
	}
	bilans.SkippedCount = 0
	return shared.StudioStyleDeleteResponse{
		Deleted: true, Form: stan.forma, Balance: bilans,
	}, nil
}

// postacLiczebnik odmienia rzeczownik po liczbie wedle zasad polskich. Zdanie
// „przestawił 3 miejsc użycia" jest zdaniem produktu niedokończonego.
func postacLiczebnik(ile int, jeden, kilka, wiele string) string {
	liczba := postacZapisLiczby(ile)
	reszta, dziesiatka := ile%10, ile%100
	switch {
	case ile == 1:
		return liczba + " " + jeden
	case reszta >= 2 && reszta <= 4 && (dziesiatka < 12 || dziesiatka > 14):
		return liczba + " " + kilka
	default:
		return liczba + " " + wiele
	}
}
