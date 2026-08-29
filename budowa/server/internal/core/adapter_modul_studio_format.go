// Postać znaku i postać akapitu na wskazanym fragmencie, czyszczenie
// formatowania, wielkość liter, malarz formatów, zaznaczanie wedle podobnego
// formatowania oraz znajdź i zamień z postacią; malarz kopiuje postać wierszem
// tabeli, nie pamięcią procesu.
package core

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"time"
	"unicode"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// postacTrwanieMalarza mówi, jak długo obowiązuje zabrana postać: godzina jest
// miarą jednej pracy nad pismem, nie magazynem na później.
const postacTrwanieMalarza = time.Hour

// MalarzSkladnicaStudia jest kontraktem tabeli malarza formatów (migracja
// 368); brak tabeli jest brakiem montażu rdzenia i mówi to wprost.
type MalarzSkladnicaStudia interface {
	ZapiszPostacMalarza(ctx context.Context,
		postac dane.PostacMalarzaStudia) (dane.PostacMalarzaStudia, error)
	PostacMalarza(ctx context.Context, kod string) (dane.PostacMalarzaStudia, error)
	NajswiezszaPostacMalarza(ctx context.Context, okno string) (dane.PostacMalarzaStudia, error)
	SprzatnijPostacieMalarza(ctx context.Context) (int, error)
}

// malarzSkladnica oddaje tabelę malarza formatów, odmawiając wprost, gdy
// rdzeń jej nie zmontował przy starcie.
func (a *adapterStudia) malarzSkladnica() (MalarzSkladnicaStudia, error) {
	if a == nil || a.repozytorium == nil {
		return nil, postacBladZaplecza(
			"repozytorium Studia nie zostało podane przy montażu serwera")
	}
	skladnica, jest := a.repozytorium.(MalarzSkladnicaStudia)
	if !jest {
		return nil, postacBladZaplecza("repozytorium Studia nie niesie tabeli malarza " +
			"formatów z migracji 368 — zabranej postaci nie ma gdzie odłożyć")
	}
	return skladnica, nil
}

// postacZabrana to postać zabrana malarzem, złożona z wiersza albo do wiersza
// tabeli malarza formatów.
type postacZabrana struct {
	znak   *shared.StudioCharacterFormat
	akapit *shared.StudioParagraphFormat
}

// postacZabranaDoWiersza przekłada zabraną postać na wiersz warstwy danych;
// postać jedzie zapisem JSON, nie osobną kolumną na każdą cechę.
func postacZabranaDoWiersza(kod, okno, kodDokumentu string,
	zabrana postacZabrana) (dane.PostacMalarzaStudia, error) {

	wiersz := dane.PostacMalarzaStudia{Kod: kod, Okno: okno}
	if strings.TrimSpace(kodDokumentu) != "" {
		wiersz.DokumentKod = postacWskaznikTekstu(kodDokumentu)
	}
	if zabrana.znak != nil {
		zapis, err := json.Marshal(zabrana.znak)
		if err != nil {
			return dane.PostacMalarzaStudia{}, postacBladZaplecza(
				"zabranej postaci znaku nie da się zapisać: " + err.Error())
		}
		wiersz.PostacZnakuJSON = postacWskaznikTekstu(string(zapis))
	}
	if zabrana.akapit != nil {
		zapis, err := json.Marshal(zabrana.akapit)
		if err != nil {
			return dane.PostacMalarzaStudia{}, postacBladZaplecza(
				"zabranej postaci akapitu nie da się zapisać: " + err.Error())
		}
		wiersz.PostacAkapituJSON = postacWskaznikTekstu(string(zapis))
	}
	wygasa := time.Now().UTC().Add(postacTrwanieMalarza).Format("2006-01-02T15:04:05.000Z")
	wiersz.Wygasa = &wygasa
	return wiersz, nil
}

// postacZabranaZWiersza składa zabraną postać z wiersza warstwy danych.
//
// Zapis nieczytelny jest tu USTERKĄ ZAPLECZA, nie brakiem postaci: wiersz stoi,
// więc Operator postać zabrał, a odmowa „nie masz zabranej postaci" byłaby
// nieprawdą o tym, co zrobił.
func postacZabranaZWiersza(wiersz dane.PostacMalarzaStudia) (postacZabrana, error) {
	var zabrana postacZabrana
	if wiersz.PostacZnakuJSON != nil && *wiersz.PostacZnakuJSON != "" {
		var znak shared.StudioCharacterFormat
		if err := json.Unmarshal([]byte(*wiersz.PostacZnakuJSON), &znak); err != nil {
			return postacZabrana{}, postacBladZaplecza("zabrana postać znaku " +
				wiersz.Kod + " jest nieczytelna: " + err.Error())
		}
		zabrana.znak = &znak
	}
	if wiersz.PostacAkapituJSON != nil && *wiersz.PostacAkapituJSON != "" {
		var akapit shared.StudioParagraphFormat
		if err := json.Unmarshal([]byte(*wiersz.PostacAkapituJSON), &akapit); err != nil {
			return postacZabrana{}, postacBladZaplecza("zabrana postać akapitu " +
				wiersz.Kod + " jest nieczytelna: " + err.Error())
		}
		zabrana.akapit = &akapit
	}
	return zabrana, nil
}

// ── Postać znaku ────────────────────────────────────────────────────────────

// UstawPostacZnaku nakłada postać znaku na fragment
// (`studio.format.character.set`), z opcją powiększenia względem stopnia skutecznego.
func (a *adapterStudia) UstawPostacZnaku(ctx context.Context,
	z shared.StudioFormatCharacterSetRequest) (shared.StudioFormatCharacterSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatCharacterSetResponse{}, err
	}
	autor := postacAutor(z.Author)
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))

	zmiana := shared.StudioCharacterFormat{
		FontFamily: z.FontFamily, FontSizePt: z.FontSizePt, Bold: z.Bold, Italic: z.Italic,
		Underline: z.Underline, Strikethrough: z.Strikethrough, Superscript: z.Superscript,
		Subscript: z.Subscript, Color: z.Color, HighlightColor: z.HighlightColor,
		LetterSpacingPt: z.LetterSpacingPt, SmallCaps: z.SmallCaps, AllCaps: z.AllCaps,
		Effect: z.Effect, Language: z.Language,
	}
	if zmiana == (shared.StudioCharacterFormat{}) && z.FontSizeStepPt == nil {
		return shared.StudioFormatCharacterSetResponse{}, bladWskazaniaStudio(
			"ustawienie postaci znaku bez ani jednej cechy do ustawienia")
	}

	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)
	bilans := shared.StudioActionBalance{Skipped: pominiete}
	for _, odcinek := range odcinki {
		postacRozetnij(&stan.forma, odcinek[0], odcinek[1])
		for _, wskazanie := range postacFragmentyZakresu(&stan.forma, odcinek[0], odcinek[1]) {
			run := &stan.forma.Blocks[wskazanie[0]].Runs[wskazanie[1]]
			nowa := zmiana
			if z.FontSizeStepPt != nil {
				// Powiększenie liczy się od stopnia skutecznego — „powiększ o 2” powiększa, nie ustawia na 2.
				skuteczna := postacZnakSkutecznyWBloku(&stan.forma,
					&stan.forma.Blocks[wskazanie[0]], run.Format)
				podstawa := 11.0
				if skuteczna.FontSizePt != nil {
					podstawa = *skuteczna.FontSizePt
				}
				nowyStopien := podstawa + *z.FontSizeStepPt
				if nowyStopien < 1 {
					nowyStopien = 1
				}
				if nowyStopien > 999 {
					nowyStopien = 999
				}
				nowa.FontSizePt = postacWskaznikMiary(nowyStopien)
			}
			run.Format = postacScalZnak(run.Format, nowa)
			bilans.Applied++
		}
	}
	if bilans.Applied == 0 {
		return shared.StudioFormatCharacterSetResponse{}, bladWskazaniaStudio(
			"postać znaku nie miała na czym stanąć: zakres " + postacZapisZakresu(od, do) +
				" nie obejmuje ani jednego fragmentu wolnego od blokady")
	}
	bilans.Note = postacWskaznikTekstu("postać znaku nałożona " + postacZapisZakresu(od, do))

	forma, bilansGotowy, zmianaSledzona, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindFormatChange, od, do, bilans)
	if err != nil {
		return shared.StudioFormatCharacterSetResponse{}, err
	}
	return shared.StudioFormatCharacterSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmianaSledzona,
	}, nil
}

// PostacZnaku oddaje postać znaku obowiązującą na fragmencie
// (`studio.format.character.get`) wraz z wykazem cech niejednolitych.
func (a *adapterStudia) PostacZnaku(ctx context.Context,
	z shared.StudioFormatCharacterGetRequest) (shared.StudioFormatCharacterGetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatCharacterGetResponse{}, err
	}
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	postacRozetnij(&stan.forma, od, do)

	runy := make([]shared.StudioDocumentRun, 0, 8)
	for _, wskazanie := range postacFragmentyZakresu(&stan.forma, od, do) {
		run := stan.forma.Blocks[wskazanie[0]].Runs[wskazanie[1]]
		skuteczna := postacZnakSkutecznyWBloku(&stan.forma,
			&stan.forma.Blocks[wskazanie[0]], run.Format)
		run.Format = &skuteczna
		runy = append(runy, run)
	}
	odpowiedz := shared.StudioFormatCharacterGetResponse{
		Runs:        runy,
		MixedFields: postacPolaNiejednolite(runy),
	}
	if len(runy) > 0 && runy[0].Format != nil {
		odpowiedz.Character = *runy[0].Format
	}
	return odpowiedz, nil
}

// ── Postać akapitu ──────────────────────────────────────────────────────────

// UstawPostacAkapitu nakłada postać akapitu (`studio.format.paragraph.set`)
// na wszystkie akapity zakresu.
func (a *adapterStudia) UstawPostacAkapitu(ctx context.Context,
	z shared.StudioFormatParagraphSetRequest) (shared.StudioFormatParagraphSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatParagraphSetResponse{}, err
	}
	autor := postacAutor(z.Author)
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))

	zmiana := shared.StudioParagraphFormat{
		Align: z.Align, FirstLineIndentMm: z.FirstLineIndentMm,
		IndentLeftMm: z.IndentLeftMm, IndentRightMm: z.IndentRightMm,
		SpaceBeforePt: z.SpaceBeforePt, SpaceAfterPt: z.SpaceAfterPt,
		LineSpacingRule: z.LineSpacingRule, LineSpacingValue: z.LineSpacingValue,
		ShadingColor: z.ShadingColor, WidowControl: z.WidowControl,
		KeepWithNext: z.KeepWithNext, KeepLines: z.KeepLines,
		OutlineLevel: z.OutlineLevel, RightToLeft: z.RightToLeft,
	}
	if len(z.TabStops) > 0 {
		var tabulatory []shared.StudioTabStop
		if err := json.Unmarshal(z.TabStops, &tabulatory); err != nil {
			return shared.StudioFormatParagraphSetResponse{}, bladWskazaniaStudio(
				"wykaz tabulatorów jest nieczytelny: " + err.Error())
		}
		zmiana.TabStops = tabulatory
	}
	if len(z.Border) > 0 {
		var obramowanie shared.StudioBorder
		if err := json.Unmarshal(z.Border, &obramowanie); err != nil {
			return shared.StudioFormatParagraphSetResponse{}, bladWskazaniaStudio(
				"obramowanie akapitu jest nieczytelne: " + err.Error())
		}
		zmiana.Border = &obramowanie
	}
	if postacAkapitPusty(zmiana) && z.IndentStepMm == nil {
		return shared.StudioFormatParagraphSetResponse{}, bladWskazaniaStudio(
			"ustawienie postaci akapitu bez ani jednej cechy do ustawienia")
	}

	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)
	bilans := shared.StudioActionBalance{Skipped: pominiete}
	for _, odcinek := range odcinki {
		for _, wskazanie := range postacBlokiZakresu(&stan.forma, odcinek[0], odcinek[1]) {
			blok := &stan.forma.Blocks[wskazanie]
			nowa := zmiana
			if z.IndentStepMm != nil {
				// „Zwiększ wcięcie” liczy od wcięcia skutecznego i nie schodzi pod zero.
				skuteczna := postacAkapitSkuteczny(&stan.forma, *blok)
				podstawa := 0.0
				if skuteczna.IndentLeftMm != nil {
					podstawa = *skuteczna.IndentLeftMm
				}
				nowe := podstawa + *z.IndentStepMm
				if nowe < 0 {
					nowe = 0
				}
				nowa.IndentLeftMm = postacWskaznikMiary(nowe)
			}
			blok.Paragraph = postacScalAkapit(blok.Paragraph, nowa)
			if nowa.OutlineLevel != nil {
				if *nowa.OutlineLevel > 0 {
					blok.Kind = blokPostaciNaglowek
				} else {
					blok.Kind = blokPostaciAkapit
				}
			}
			bilans.Applied++
		}
	}
	if bilans.Applied == 0 {
		return shared.StudioFormatParagraphSetResponse{}, bladWskazaniaStudio(
			"postać akapitu nie miała na czym stanąć: zakres " + postacZapisZakresu(od, do) +
				" nie obejmuje ani jednego akapitu wolnego od blokady")
	}
	bilans.Note = postacWskaznikTekstu("postać akapitu nałożona " + postacZapisZakresu(od, do))

	forma, bilansGotowy, zmianaSledzona, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindFormatChange, od, do, bilans)
	if err != nil {
		return shared.StudioFormatParagraphSetResponse{}, err
	}
	return shared.StudioFormatParagraphSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmianaSledzona,
	}, nil
}

// postacAkapitPusty mówi, czy żądanie nie niesie ani jednej cechy akapitu do
// ustawienia na fragmencie.
func postacAkapitPusty(zmiana shared.StudioParagraphFormat) bool {
	return zmiana.Align == nil && zmiana.FirstLineIndentMm == nil && zmiana.IndentLeftMm == nil &&
		zmiana.IndentRightMm == nil && zmiana.SpaceBeforePt == nil && zmiana.SpaceAfterPt == nil &&
		zmiana.LineSpacingRule == nil && zmiana.LineSpacingValue == nil && zmiana.TabStops == nil &&
		zmiana.Border == nil && zmiana.ShadingColor == nil && zmiana.WidowControl == nil &&
		zmiana.KeepWithNext == nil && zmiana.KeepLines == nil && zmiana.OutlineLevel == nil &&
		zmiana.RightToLeft == nil
}

// PostacAkapitu oddaje postać akapitu fragmentu (`studio.format.paragraph.get`)
// wraz z cechami niejednolitymi.
func (a *adapterStudia) PostacAkapitu(ctx context.Context,
	z shared.StudioFormatParagraphGetRequest) (shared.StudioFormatParagraphGetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatParagraphGetResponse{}, err
	}
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	bloki := postacBlokiZakresu(&stan.forma, od, do)
	if len(bloki) == 0 {
		return shared.StudioFormatParagraphGetResponse{}, bladWskazaniaStudio(
			"zakres " + postacZapisZakresu(od, do) + " nie obejmuje ani jednego akapitu")
	}
	pierwszy := postacAkapitSkuteczny(&stan.forma, stan.forma.Blocks[bloki[0]])

	niejednolite := make([]string, 0, 6)
	for _, wskazanie := range bloki[1:] {
		biezacy := postacAkapitSkuteczny(&stan.forma, stan.forma.Blocks[wskazanie])
		for nazwa, rozne := range map[string]bool{
			"align":             !postacToSamoWyrownanie(pierwszy.Align, biezacy.Align),
			"indentLeftMm":      !postacTaSamaMiara(pierwszy.IndentLeftMm, biezacy.IndentLeftMm),
			"firstLineIndentMm": !postacTaSamaMiara(pierwszy.FirstLineIndentMm, biezacy.FirstLineIndentMm),
			"spaceBeforePt":     !postacTaSamaMiara(pierwszy.SpaceBeforePt, biezacy.SpaceBeforePt),
			"spaceAfterPt":      !postacTaSamaMiara(pierwszy.SpaceAfterPt, biezacy.SpaceAfterPt),
			"styleName":         !postacTenSamTekst(pierwszy.StyleName, biezacy.StyleName),
		} {
			if rozne && !postacJuzJest(niejednolite, nazwa) {
				niejednolite = append(niejednolite, nazwa)
			}
		}
	}
	return shared.StudioFormatParagraphGetResponse{
		Paragraph: pierwszy, MixedFields: niejednolite,
	}, nil
}

func postacToSamoWyrownanie(pierwsze, drugie *shared.StudioTextAlign) bool {
	if pierwsze == nil || drugie == nil {
		return pierwsze == nil && drugie == nil
	}
	return *pierwsze == *drugie
}

func postacJuzJest(wykaz []string, nazwa string) bool {
	for _, pozycja := range wykaz {
		if pozycja == nazwa {
			return true
		}
	}
	return false
}

// ── Czyszczenie i wielkość liter ────────────────────────────────────────────

// CzyscPostac zdejmuje formatowanie z fragmentu (`studio.format.clear`); brak
// wskazania znaczy oboje, a nazwa stylu zostaje.
func (a *adapterStudia) CzyscPostac(ctx context.Context,
	z shared.StudioFormatClearRequest) (shared.StudioFormatClearResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatClearResponse{}, err
	}
	autor := postacAutor(z.Author)
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	znak := z.Character == nil || *z.Character
	akapit := z.Paragraph == nil || *z.Paragraph
	if z.Character != nil && z.Paragraph == nil {
		akapit = false
	}
	if z.Paragraph != nil && z.Character == nil {
		znak = false
	}
	if !znak && !akapit {
		return shared.StudioFormatClearResponse{}, bladWskazaniaStudio(
			"czyszczenie formatowania, które nie czyści ani znaku, ani akapitu")
	}

	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)
	bilans := shared.StudioActionBalance{Skipped: pominiete}
	for _, odcinek := range odcinki {
		postacRozetnij(&stan.forma, odcinek[0], odcinek[1])
		if znak {
			for _, wskazanie := range postacFragmentyZakresu(&stan.forma, odcinek[0], odcinek[1]) {
				run := &stan.forma.Blocks[wskazanie[0]].Runs[wskazanie[1]]
				var nazwaStylu *string
				if run.Format != nil {
					nazwaStylu = run.Format.StyleName
				}
				if nazwaStylu == nil {
					run.Format = nil
				} else {
					run.Format = &shared.StudioCharacterFormat{StyleName: nazwaStylu}
				}
				bilans.Applied++
			}
		}
		if akapit {
			for _, wskazanie := range postacBlokiZakresu(&stan.forma, odcinek[0], odcinek[1]) {
				blok := &stan.forma.Blocks[wskazanie]
				var nazwaStylu *string
				if blok.Paragraph != nil {
					nazwaStylu = blok.Paragraph.StyleName
				}
				if nazwaStylu == nil {
					blok.Paragraph = nil
				} else {
					blok.Paragraph = &shared.StudioParagraphFormat{StyleName: nazwaStylu}
				}
				bilans.Applied++
			}
		}
	}
	if bilans.Applied == 0 {
		return shared.StudioFormatClearResponse{}, bladWskazaniaStudio(
			"czyszczenie formatowania nie miało czego wyczyścić w zakresie " +
				postacZapisZakresu(od, do))
	}
	bilans.Note = postacWskaznikTekstu("formatowanie zdjęte " + postacZapisZakresu(od, do) +
		"; przypisanie do stylu nazwanego zostało")

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindFormatChange, od, do, bilans)
	if err != nil {
		return shared.StudioFormatClearResponse{}, err
	}
	return shared.StudioFormatClearResponse{Form: forma, Balance: bilansGotowy, Change: zmiana}, nil
}

// UstawWielkoscLiter przekłada wielkość liter fragmentu
// (`studio.format.case.set`); to zmiana treści, nie postaci, więc idzie
// drogą zamiany treści.
func (a *adapterStudia) UstawWielkoscLiter(ctx context.Context,
	z shared.StudioFormatCaseSetRequest) (shared.StudioFormatCaseSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatCaseSetResponse{}, err
	}
	autor := postacAutor(z.Author)
	znaki := []rune(postacTekstFormy(&stan.forma))
	od, do := postacZakres(z.RangeStart, z.RangeEnd, len(znaki))
	if od == do {
		return shared.StudioFormatCaseSetResponse{}, bladWskazaniaStudio(
			"zmiana wielkości liter na pustym zaznaczeniu — nie ma czego zamienić")
	}

	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)
	bilans := shared.StudioActionBalance{Skipped: pominiete}
	// Odcinki idą od końca, żeby zmiana długości nie unieważniła zakresów jeszcze nieprzetworzonych.
	for i := len(odcinki) - 1; i >= 0; i-- {
		odcinek := odcinki[i]
		biezace := []rune(postacTekstFormy(&stan.forma))
		if odcinek[1] > len(biezace) {
			continue
		}
		zrodlo := string(biezace[odcinek[0]:odcinek[1]])
		przelozone := postacPrzelozWielkosc(zrodlo, z.Transform)
		if przelozone == zrodlo {
			continue
		}
		postacZamienTresc(&stan.forma, odcinek[0], odcinek[1], przelozone, nil, &autor)
		bilans.Applied++
	}
	if bilans.Applied == 0 {
		return shared.StudioFormatCaseSetResponse{}, bladWskazaniaStudio(
			"zakres " + postacZapisZakresu(od, do) + " jest już w tej wielkości liter " +
				"albo w całości zablokowany: " + postacNazwaBlokad(pominiete))
	}
	bilans.Note = postacWskaznikTekstu("wielkość liter przełożona " + postacZapisZakresu(od, do))

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindWstawienie, shared.StudioActionKindTextEdit, od, do, bilans)
	if err != nil {
		return shared.StudioFormatCaseSetResponse{}, err
	}
	return shared.StudioFormatCaseSetResponse{Form: forma, Balance: bilansGotowy, Change: zmiana}, nil
}

// postacPrzelozWielkosc przekłada wielkość liter wedle wskazanej odmiany:
// wielka, mała, przełącznik, każde słowo albo zdanie.
func postacPrzelozWielkosc(zrodlo string, odmiana shared.StudioCaseTransform) string {
	switch odmiana {
	case shared.StudioCaseTransformUpper:
		return strings.ToUpper(zrodlo)
	case shared.StudioCaseTransformLower:
		return strings.ToLower(zrodlo)
	case shared.StudioCaseTransformToggle:
		return strings.Map(func(litera rune) rune {
			if unicode.IsUpper(litera) {
				return unicode.ToLower(litera)
			}
			return unicode.ToUpper(litera)
		}, zrodlo)
	case shared.StudioCaseTransformCapitalize:
		// Każde słowo wielką literą; granicą słowa jest wszystko, co nie jest literą ani cyfrą.
		wynik := make([]rune, 0, len(zrodlo))
		poczatek := true
		for _, litera := range zrodlo {
			if poczatek && unicode.IsLetter(litera) {
				wynik = append(wynik, unicode.ToUpper(litera))
				poczatek = false
				continue
			}
			wynik = append(wynik, unicode.ToLower(litera))
			poczatek = !unicode.IsLetter(litera) && !unicode.IsDigit(litera)
		}
		return string(wynik)
	case shared.StudioCaseTransformSentence:
		// Pierwsza litera zdania wielka, reszta mała. Zdanie kończy kropka,
		// znak zapytania albo wykrzyknik.
		wynik := make([]rune, 0, len(zrodlo))
		poczatek := true
		for _, litera := range zrodlo {
			switch {
			case poczatek && unicode.IsLetter(litera):
				wynik = append(wynik, unicode.ToUpper(litera))
				poczatek = false
			default:
				wynik = append(wynik, unicode.ToLower(litera))
			}
			if litera == '.' || litera == '!' || litera == '?' || litera == '\n' {
				poczatek = true
			}
		}
		return string(wynik)
	default:
		return zrodlo
	}
}

// ── Malarz formatów ─────────────────────────────────────────────────────────

// ZabierzPostac zabiera postać fragmentu (`studio.format.painter.copy`)
// i odkłada ją wierszem malarza formatów.
func (a *adapterStudia) ZabierzPostac(ctx context.Context,
	z shared.StudioFormatPainterCopyRequest) (shared.StudioFormatPainterCopyResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatPainterCopyResponse{}, err
	}
	skladnica, err := a.malarzSkladnica()
	if err != nil {
		return shared.StudioFormatPainterCopyResponse{}, err
	}
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	postacRozetnij(&stan.forma, od, do)

	wskazania := postacFragmentyZakresu(&stan.forma, od, do)
	if len(wskazania) == 0 {
		return shared.StudioFormatPainterCopyResponse{}, bladWskazaniaStudio(
			"malarz formatów nie miał skąd zabrać postaci: zakres " +
				postacZapisZakresu(od, do) + " nie obejmuje ani jednego fragmentu")
	}
	znak := postacZnakSkutecznyWBloku(&stan.forma, &stan.forma.Blocks[wskazania[0][0]],
		stan.forma.Blocks[wskazania[0][0]].Runs[wskazania[0][1]].Format)
	zabrana := postacZabrana{znak: &znak}

	odpowiedz := shared.StudioFormatPainterCopyResponse{
		ClipId: nowyIdentyfikator(przedrostekMalarzaPostaci), Character: &znak,
	}
	// Postać akapitu idzie z postacią znaku na żądanie — malarz na słowie maluje krój, nie wyrównanie.
	if z.IncludeParagraph != nil && *z.IncludeParagraph {
		if bloki := postacBlokiZakresu(&stan.forma, od, do); len(bloki) > 0 {
			akapit := postacAkapitSkuteczny(&stan.forma, stan.forma.Blocks[bloki[0]])
			zabrana.akapit = &akapit
			odpowiedz.Paragraph = &akapit
		}
	}
	// Wpisy wygasłe schodzą przy zabraniu, nie przy odczycie; sprzątanie zabrania nie unieważnia.
	_, _ = skladnica.SprzatnijPostacieMalarza(ctx)

	// Oknem wpisu jest okno dokumentu — wykaz „nanieś ostatnio zabrane” rozdziela dwóch Operatorów.
	wiersz, err := postacZabranaDoWiersza(odpowiedz.ClipId, stan.dokument.Okno,
		stan.dokument.Kod, zabrana)
	if err != nil {
		return shared.StudioFormatPainterCopyResponse{}, err
	}
	if _, err := skladnica.ZapiszPostacMalarza(ctx, wiersz); err != nil {
		return shared.StudioFormatPainterCopyResponse{}, bladStudio(err)
	}
	return odpowiedz, nil
}

// PolozPostac kładzie zabraną postać na fragmencie
// (`studio.format.painter.apply`) wedle wskazanego uchwytu.
func (a *adapterStudia) PolozPostac(ctx context.Context,
	z shared.StudioFormatPainterApplyRequest) (shared.StudioFormatPainterApplyResponse, error) {

	skladnica, err := a.malarzSkladnica()
	if err != nil {
		return shared.StudioFormatPainterApplyResponse{}, err
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatPainterApplyResponse{}, err
	}

	// Uchwyt podany wskazuje wprost, czego położyć; pusty znaczy „połóż to, co ostatnio zabrałem w oknie”.
	wiersz, err := malarzWiersz(ctx, skladnica, strings.TrimSpace(z.ClipId), stan.dokument.Okno)
	if err != nil {
		return shared.StudioFormatPainterApplyResponse{}, err
	}
	zabrana, err := postacZabranaZWiersza(wiersz)
	if err != nil {
		return shared.StudioFormatPainterApplyResponse{}, err
	}
	autor := postacAutor(z.Author)
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)

	zAkapitem := zabrana.akapit != nil && (z.IncludeParagraph == nil || *z.IncludeParagraph)
	bilans := shared.StudioActionBalance{Skipped: pominiete}
	for _, odcinek := range odcinki {
		postacRozetnij(&stan.forma, odcinek[0], odcinek[1])
		if zabrana.znak != nil {
			for _, wskazanie := range postacFragmentyZakresu(&stan.forma, odcinek[0], odcinek[1]) {
				run := &stan.forma.Blocks[wskazanie[0]].Runs[wskazanie[1]]
				run.Format = postacKopiaZnaku(zabrana.znak)
				bilans.Applied++
			}
		}
		if zAkapitem {
			for _, wskazanie := range postacBlokiZakresu(&stan.forma, odcinek[0], odcinek[1]) {
				stan.forma.Blocks[wskazanie].Paragraph = postacKopiaAkapitu(zabrana.akapit)
				bilans.Applied++
			}
		}
	}
	if bilans.Applied == 0 {
		return shared.StudioFormatPainterApplyResponse{}, bladWskazaniaStudio(
			"malarz formatów nie miał na czym stanąć w zakresie " + postacZapisZakresu(od, do))
	}
	bilans.Note = postacWskaznikTekstu("zabrana postać położona " + postacZapisZakresu(od, do))

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindFormatChange, od, do, bilans)
	if err != nil {
		return shared.StudioFormatPainterApplyResponse{}, err
	}
	return shared.StudioFormatPainterApplyResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
	}, nil
}

// malarzWiersz odnajduje wiersz zabranej postaci: po uchwycie albo — gdy
// uchwytu nie podano — najświeższy niewygasły wpis okna. Odmowa nazywa powód
// prawdziwy, wspólny dla wpisu wygasłego i nieistniejącego.
func malarzWiersz(ctx context.Context, skladnica MalarzSkladnicaStudia,
	uchwyt, okno string) (dane.PostacMalarzaStudia, error) {

	if uchwyt != "" {
		wiersz, err := skladnica.PostacMalarza(ctx, uchwyt)
		if err != nil {
			return dane.PostacMalarzaStudia{}, bladWskazaniaStudio("malarz formatów nie ma " +
				"zabranej postaci o uchwycie " + uchwyt + " — wpisu nie ma albo wygasł " +
				"(zabrana postać obowiązuje godzinę). Zabierz postać ponownie przez " +
				"studio.format.painter.copy")
		}
		return wiersz, nil
	}
	wiersz, err := skladnica.NajswiezszaPostacMalarza(ctx, okno)
	if err != nil {
		return dane.PostacMalarzaStudia{}, bladWskazaniaStudio("położenie postaci bez " +
			"wskazania uchwytu, a w tym oknie nie ma ani jednej zabranej postaci " +
			"obowiązującej. Zabierz postać przez studio.format.painter.copy")
	}
	return wiersz, nil
}

// ── Zaznaczanie wedle podobnego formatowania ────────────────────────────────

// ZaznaczPodobne oddaje fragmenty o postaci podobnej do wzorca
// (`studio.format.similar.select`): stylu nazwanego albo postaci fragmentu
// wskazanego zakresem.
func (a *adapterStudia) ZaznaczPodobne(ctx context.Context,
	z shared.StudioFormatSimilarSelectRequest) (shared.StudioFormatSimilarSelectResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatSimilarSelectResponse{}, err
	}
	poZnaku := z.MatchCharacter == nil || *z.MatchCharacter
	poAkapicie := z.MatchParagraph != nil && *z.MatchParagraph

	var wzorZnaku *shared.StudioCharacterFormat
	var wzorAkapitu *shared.StudioParagraphFormat
	nazwaStylu := ""

	if z.StyleName != nil && strings.TrimSpace(*z.StyleName) != "" {
		nazwaStylu = strings.TrimSpace(*z.StyleName)
		if postacStyl(&stan.forma, nazwaStylu) == nil {
			return shared.StudioFormatSimilarSelectResponse{}, bladWskazaniaStudio(
				"styl „" + nazwaStylu + "” nie istnieje w arkuszu tego dokumentu")
		}
	} else {
		od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
		postacRozetnij(&stan.forma, od, do)
		wskazania := postacFragmentyZakresu(&stan.forma, od, do)
		if len(wskazania) == 0 {
			return shared.StudioFormatSimilarSelectResponse{}, bladWskazaniaStudio(
				"zaznaczanie wedle podobnego formatowania bez wzorca: podaj styl nazwany " +
					"albo zakres fragmentu, którego postać ma być wzorcem")
		}
		znak := postacZnakSkutecznyWBloku(&stan.forma, &stan.forma.Blocks[wskazania[0][0]],
			stan.forma.Blocks[wskazania[0][0]].Runs[wskazania[0][1]].Format)
		wzorZnaku = &znak
		if poAkapicie {
			akapit := postacAkapitSkuteczny(&stan.forma, stan.forma.Blocks[wskazania[0][0]])
			wzorAkapitu = &akapit
		}
	}

	postacPrzeliczZakresy(&stan.forma)
	zgodne := make([]shared.StudioDocumentRun, 0, 16)
	for i := range stan.forma.Blocks {
		blok := &stan.forma.Blocks[i]
		if !postacBlokNiesieTekst(*blok) {
			continue
		}
		if poAkapicie && wzorAkapitu != nil {
			biezacy := postacAkapitSkuteczny(&stan.forma, *blok)
			if !postacAkapitPodobny(*wzorAkapitu, biezacy) {
				continue
			}
		}
		for j := range blok.Runs {
			run := blok.Runs[j]
			if run.Text == "" {
				continue
			}
			skuteczna := postacZnakSkutecznyWBloku(&stan.forma, blok, run.Format)
			switch {
			case nazwaStylu != "":
				wlasny := ""
				if run.Format != nil && run.Format.StyleName != nil {
					wlasny = *run.Format.StyleName
				}
				akapitowy := ""
				if blok.Paragraph != nil && blok.Paragraph.StyleName != nil {
					akapitowy = *blok.Paragraph.StyleName
				}
				if wlasny != nazwaStylu && akapitowy != nazwaStylu {
					continue
				}
			case poZnaku && wzorZnaku != nil:
				if !postacZnakPodobny(*wzorZnaku, skuteczna) {
					continue
				}
			}
			run.Format = &skuteczna
			zgodne = append(zgodne, run)
		}
	}
	if len(zgodne) == 0 {
		return shared.StudioFormatSimilarSelectResponse{}, bladWskazaniaStudio(
			"w dokumencie nie ma ani jednego fragmentu o tej postaci")
	}
	return shared.StudioFormatSimilarSelectResponse{Matches: zgodne, Count: len(zgodne)}, nil
}

// postacZnakPodobny porównuje postać znaku wedle cech, które Operator widzi
// gołym okiem. Barwa wyróżnienia i język nie liczą się do podobieństwa —
// „to samo formatowanie" znaczy dla Operatora ten sam krój, stopień i grubość.
func postacZnakPodobny(wzor, biezaca shared.StudioCharacterFormat) bool {
	return postacTenSamTekst(wzor.FontFamily, biezaca.FontFamily) &&
		postacTaSamaMiara(wzor.FontSizePt, biezaca.FontSizePt) &&
		postacTaSamaPrawda(wzor.Bold, biezaca.Bold) &&
		postacTaSamaPrawda(wzor.Italic, biezaca.Italic) &&
		postacTenSamTekst(wzor.Color, biezaca.Color)
}

// postacAkapitPodobny porównuje postać akapitu wedle cech widocznych: wyrównania,
// wcięcia lewego i wcięcia pierwszego wiersza.
func postacAkapitPodobny(wzor, biezacy shared.StudioParagraphFormat) bool {
	return postacToSamoWyrownanie(wzor.Align, biezacy.Align) &&
		postacTaSamaMiara(wzor.IndentLeftMm, biezacy.IndentLeftMm) &&
		postacTaSamaMiara(wzor.FirstLineIndentMm, biezacy.FirstLineIndentMm)
}

// ── Znajdź i zamień z postacią ──────────────────────────────────────────────

// ZamienZPostacia szuka i zamienia treść oraz postać (`studio.format.replace`):
// brzmienie, postać znalezionego fragmentu i styl nazwany, także bez tekstu.
func (a *adapterStudia) ZamienZPostacia(ctx context.Context,
	z shared.StudioFormatReplaceRequest) (shared.StudioFormatReplaceResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFormatReplaceResponse{}, err
	}
	autor := postacAutor(z.Author)

	szukanyTekst := ""
	if z.FindText != nil {
		szukanyTekst = *z.FindText
	}
	var szukanaPostac *shared.StudioCharacterFormat
	if len(z.FindFormat) > 0 {
		var postac shared.StudioCharacterFormat
		if err := json.Unmarshal(z.FindFormat, &postac); err != nil {
			return shared.StudioFormatReplaceResponse{}, bladWskazaniaStudio(
				"postać szukana jest nieczytelna: " + err.Error())
		}
		szukanaPostac = &postac
	}
	szukanyStyl := ""
	if z.FindStyleName != nil {
		szukanyStyl = strings.TrimSpace(*z.FindStyleName)
	}
	if szukanyTekst == "" && szukanaPostac == nil && szukanyStyl == "" {
		return shared.StudioFormatReplaceResponse{}, bladWskazaniaStudio(
			"zamiana bez wskazania, czego szukać: podaj brzmienie, postać albo styl nazwany")
	}

	var nowaPostac *shared.StudioCharacterFormat
	if len(z.ReplaceFormat) > 0 {
		var postac shared.StudioCharacterFormat
		if err := json.Unmarshal(z.ReplaceFormat, &postac); err != nil {
			return shared.StudioFormatReplaceResponse{}, bladWskazaniaStudio(
				"postać zastępcza jest nieczytelna: " + err.Error())
		}
		nowaPostac = &postac
	}
	nowyStyl := ""
	if z.ReplaceStyleName != nil {
		nowyStyl = strings.TrimSpace(*z.ReplaceStyleName)
		if nowyStyl != "" && postacStyl(&stan.forma, nowyStyl) == nil {
			return shared.StudioFormatReplaceResponse{}, bladWskazaniaStudio(
				"styl zastępczy „" + nowyStyl + "” nie istnieje w arkuszu tego dokumentu")
		}
	}
	if z.ReplaceText == nil && nowaPostac == nil && nowyStyl == "" {
		return shared.StudioFormatReplaceResponse{}, bladWskazaniaStudio(
			"zamiana bez wskazania, na co zamienić: podaj brzmienie, postać albo styl nazwany")
	}

	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	wszystkie := z.ReplaceAll == nil || *z.ReplaceAll

	trafienia, err := a.postacZnajdzTrafienia(&stan.forma, szukanyTekst, szukanaPostac,
		szukanyStyl, od, do, z)
	if err != nil {
		return shared.StudioFormatReplaceResponse{}, err
	}
	if len(trafienia) == 0 {
		return shared.StudioFormatReplaceResponse{}, bladWskazaniaStudio(
			"w zakresie " + postacZapisZakresu(od, do) + " nie ma ani jednego trafienia")
	}
	if !wszystkie {
		trafienia = trafienia[:1]
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	zamienione := 0
	// Trafienia idą od końca — zamiana zmieniająca długość nie unieważni trafień nieprzetworzonych.
	for i := len(trafienia) - 1; i >= 0; i-- {
		trafienie := trafienia[i]
		odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, trafienie[0], trafienie[1], autor)
		if len(pominiete) > 0 {
			bilans.Skipped = append(bilans.Skipped, pominiete...)
			continue
		}
		_ = odcinki
		if z.ReplaceText != nil {
			postacRozetnij(&stan.forma, trafienie[0], trafienie[1])
			var przejeta *shared.StudioCharacterFormat
			if wskazania := postacFragmentyZakresu(&stan.forma, trafienie[0], trafienie[1]); len(wskazania) > 0 {
				przejeta = postacKopiaZnaku(
					stan.forma.Blocks[wskazania[0][0]].Runs[wskazania[0][1]].Format)
			}
			roznica := postacZamienTresc(&stan.forma, trafienie[0], trafienie[1],
				*z.ReplaceText, przejeta, &autor)
			trafienie[1] += roznica
		}
		if nowaPostac != nil || nowyStyl != "" {
			postacRozetnij(&stan.forma, trafienie[0], trafienie[1])
			for _, wskazanie := range postacFragmentyZakresu(&stan.forma, trafienie[0], trafienie[1]) {
				run := &stan.forma.Blocks[wskazanie[0]].Runs[wskazanie[1]]
				zmiana := shared.StudioCharacterFormat{}
				if nowaPostac != nil {
					zmiana = *nowaPostac
				}
				if nowyStyl != "" {
					zmiana.StyleName = postacWskaznikTekstu(nowyStyl)
				}
				run.Format = postacScalZnak(run.Format, zmiana)
			}
		}
		zamienione++
		bilans.Applied++
	}
	if zamienione == 0 {
		return shared.StudioFormatReplaceResponse{}, bladWskazaniaStudio(
			"wszystkie trafienia w zakresie " + postacZapisZakresu(od, do) +
				" stoją pod blokadą: " + postacNazwaBlokad(bilans.Skipped))
	}
	bilans.Note = postacWskaznikTekstu("trafień " + postacZapisLiczby(len(trafienia)) +
		", zamienionych " + postacZapisLiczby(zamienione))

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindWstawienie, shared.StudioActionKindTextEdit, od, do, bilans)
	if err != nil {
		return shared.StudioFormatReplaceResponse{}, err
	}
	return shared.StudioFormatReplaceResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		Matches: len(trafienia), Replaced: zamienione,
	}, nil
}

// postacZnajdzTrafienia zbiera zakresy pasujące do wskazania: brzmienie, postać
// albo styl nazwany. Trafienia wracają w kolejności czytania.
func (a *adapterStudia) postacZnajdzTrafienia(forma *shared.StudioDocumentForm,
	szukanyTekst string, szukanaPostac *shared.StudioCharacterFormat, szukanyStyl string,
	od, do int, z shared.StudioFormatReplaceRequest) ([][2]int, error) {

	trafienia := make([][2]int, 0, 16)

	if szukanyTekst != "" {
		znaki := []rune(postacTekstFormy(forma))
		if do > len(znaki) {
			do = len(znaki)
		}
		obszar := string(znaki[od:do])
		if z.Regex != nil && *z.Regex {
			wzor := szukanyTekst
			if z.MatchCase == nil || !*z.MatchCase {
				wzor = "(?i)" + wzor
			}
			wyrazenie, err := regexp.Compile(wzor)
			if err != nil {
				return nil, bladWskazaniaStudio(
					"wyrażenie szukane jest niepoprawne: " + err.Error())
			}
			for _, para := range wyrazenie.FindAllStringIndex(obszar, -1) {
				poczatek := od + len([]rune(obszar[:para[0]]))
				koniec := od + len([]rune(obszar[:para[1]]))
				trafienia = append(trafienia, [2]int{poczatek, koniec})
			}
		} else {
			trafienia = append(trafienia,
				postacTrafieniaNapisu(obszar, szukanyTekst, od,
					z.MatchCase != nil && *z.MatchCase,
					z.WholeWord != nil && *z.WholeWord)...)
		}
		return trafienia, nil
	}

	// Wyszukiwanie samą postacią albo stylem idzie po fragmentach — każdy fragment ma jednolitą postać.
	postacPrzeliczZakresy(forma)
	for i := range forma.Blocks {
		blok := &forma.Blocks[i]
		if !postacBlokNiesieTekst(*blok) {
			continue
		}
		for j := range blok.Runs {
			run := &blok.Runs[j]
			if run.Text == "" || run.RangeStart == nil || run.RangeEnd == nil {
				continue
			}
			if *run.RangeStart < od || *run.RangeEnd > do {
				continue
			}
			skuteczna := postacZnakSkutecznyWBloku(forma, blok, run.Format)
			if szukanyStyl != "" {
				wlasny := ""
				if run.Format != nil && run.Format.StyleName != nil {
					wlasny = *run.Format.StyleName
				}
				akapitowy := ""
				if blok.Paragraph != nil && blok.Paragraph.StyleName != nil {
					akapitowy = *blok.Paragraph.StyleName
				}
				if wlasny != szukanyStyl && akapitowy != szukanyStyl {
					continue
				}
			}
			if szukanaPostac != nil && !postacZnakZgodnyZWzorem(*szukanaPostac, skuteczna) {
				continue
			}
			trafienia = append(trafienia, [2]int{*run.RangeStart, *run.RangeEnd})
		}
	}
	return trafienia, nil
}

// postacZnakZgodnyZWzorem sprawdza WYŁĄCZNIE te cechy, które wzorzec niesie.
// Wzorzec „barwa czerwona" ma trafiać we wszystko czerwone, niezależnie od
// kroju — inaczej szukanie postacią byłoby bezużyteczne.
func postacZnakZgodnyZWzorem(wzor, biezaca shared.StudioCharacterFormat) bool {
	if wzor.FontFamily != nil && !postacTenSamTekst(wzor.FontFamily, biezaca.FontFamily) {
		return false
	}
	if wzor.FontSizePt != nil && !postacTaSamaMiara(wzor.FontSizePt, biezaca.FontSizePt) {
		return false
	}
	if wzor.Bold != nil && !postacTaSamaPrawda(wzor.Bold, biezaca.Bold) {
		return false
	}
	if wzor.Italic != nil && !postacTaSamaPrawda(wzor.Italic, biezaca.Italic) {
		return false
	}
	if wzor.Strikethrough != nil && !postacTaSamaPrawda(wzor.Strikethrough, biezaca.Strikethrough) {
		return false
	}
	if wzor.Color != nil && !postacTenSamTekst(wzor.Color, biezaca.Color) {
		return false
	}
	if wzor.HighlightColor != nil &&
		!postacTenSamTekst(wzor.HighlightColor, biezaca.HighlightColor) {
		return false
	}
	if wzor.Underline != nil {
		if biezaca.Underline == nil || *biezaca.Underline != *wzor.Underline {
			return false
		}
	}
	return true
}

// postacTrafieniaNapisu szuka brzmienia w obszarze i oddaje zakresy w znakach
// całego dokumentu, nie samego obszaru.
func postacTrafieniaNapisu(obszar, szukane string, przesuniecie int,
	zWielkoscia, caleSlowo bool) [][2]int {

	znakiObszaru := []rune(obszar)
	znakiSzukane := []rune(szukane)
	porownywalne := znakiObszaru
	wzor := znakiSzukane
	if !zWielkoscia {
		porownywalne = []rune(strings.ToLower(obszar))
		wzor = []rune(strings.ToLower(szukane))
	}
	trafienia := make([][2]int, 0, 8)
	for i := 0; i+len(wzor) <= len(porownywalne); i++ {
		zgodne := true
		for j := range wzor {
			if porownywalne[i+j] != wzor[j] {
				zgodne = false
				break
			}
		}
		if !zgodne {
			continue
		}
		if caleSlowo {
			przed := i == 0 || !postacZnakSlowa(znakiObszaru[i-1])
			po := i+len(wzor) == len(znakiObszaru) || !postacZnakSlowa(znakiObszaru[i+len(wzor)])
			if !przed || !po {
				continue
			}
		}
		trafienia = append(trafienia,
			[2]int{przesuniecie + i, przesuniecie + i + len(wzor)})
		i += len(wzor) - 1 // trafienia nie zachodzą na siebie
	}
	return trafienia
}

func postacZnakSlowa(litera rune) bool {
	return unicode.IsLetter(litera) || unicode.IsDigit(litera) || litera == '_'
}
