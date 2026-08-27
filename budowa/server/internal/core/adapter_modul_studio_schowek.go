// Odpowiedzialność pliku: schowek dostępny modelowi tak samo jak operatorowi —
// odłożenie fragmentu (`studio.clipboard.copy`), wklejenie
// (`studio.clipboard.paste`) oraz wykaz pochodzenia fragmentów dokumentu
// (`studio.provenance.list`).
package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// SkopiujDoSchowka obsługuje `studio.clipboard.copy`: odkłada fragment treści
// do wspólnej historii schowka platformy, wraz z postacią źródła, i wycina go
// z dokumentu, gdy żądanie o to poprosi.
func (a *adapterStudia) SkopiujDoSchowka(ctx context.Context,
	z shared.StudioClipboardCopyRequest) (shared.StudioClipboardCopyResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioClipboardCopyResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioClipboardCopyResponse{}, err
	}
	if z.RangeEnd < z.RangeStart {
		return shared.StudioClipboardCopyResponse{}, bladWskazaniaStudio(
			"odłożenie fragmentu o zakresie odwróconym: koniec " + strconv.Itoa(z.RangeEnd) +
				" leży przed początkiem " + strconv.Itoa(z.RangeStart))
	}
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	})
	od, do := postacZakres(&z.RangeStart, &z.RangeEnd, postacDlugosc(&stan.forma))
	tekst := znakowanieTekstZakresu(&stan.forma, od, do)

	// Sama postać, bez treści — to jest malarz formatów i wołamy go, a nie
	// zakładamy drugiego.
	if z.FormatOnly != nil && *z.FormatOnly {
		prawda := true
		zabrana, err := a.ZabierzPostac(ctx, shared.StudioFormatPainterCopyRequest{
			DocumentId: stan.dokument.Kod, RangeStart: &od, RangeEnd: &do,
			IncludeParagraph: &prawda,
		})
		if err != nil {
			return shared.StudioClipboardCopyResponse{}, err
		}
		return shared.StudioClipboardCopyResponse{
			ClipboardEntryId: zabrana.ClipId,
			Text:             "",
		}, nil
	}

	if strings.TrimSpace(tekst) == "" {
		return shared.StudioClipboardCopyResponse{}, bladWskazaniaStudio(
			"fragment od znaku " + strconv.Itoa(od) + " do " + strconv.Itoa(do) +
				" nie niesie treści — odłożenie pustki do schowka nie jest odłożeniem")
	}

	okno := stan.dokument.Okno
	// Postać źródła jedzie z treścią, migracja 390 — inaczej „zachowaj postać
	// źródła” nie ma co zachować.
	postacZrodla, err := schowekPostacZakresu(&stan.forma, od, do)
	if err != nil {
		return shared.StudioClipboardCopyResponse{}, err
	}
	wpis, _, err := skladnica.DopiszWpisSchowkaStudia(ctx, dane.WpisSchowka{
		Kod:          nowyIdentyfikator(przedrostekWpisuSchowka),
		Rodzaj:       shared.ClipboardEntryKindText,
		Tresc:        tekst,
		OknoZrodlowe: &okno,
		Utworzono:    time.Now().UnixMilli(),
		PostacJSON:   postacZrodla,
	})
	if err != nil {
		return shared.StudioClipboardCopyResponse{}, bladStudio(err)
	}

	odpowiedz := shared.StudioClipboardCopyResponse{
		ClipboardEntryId: wpis.Kod,
		Text:             tekst,
	}

	// Wycięcie zmienia treść, więc odkłada zmianę śledzoną i wpis dziennika, by
	// dało się je cofnąć.

	// Blokada sprawdza się tu: zapora rejestru tej komendy nie pilnuje, bo samo
	// skopiowanie jest odczytem.
	if z.Cut != nil && *z.Cut {
		blokady, err := skladnica.BlokadyFragmentow(ctx, stan.dokument.ID)
		if err != nil {
			return shared.StudioClipboardCopyResponse{}, bladStudio(err)
		}
		if blokada, trafiona := blokadaNaZakresie(blokady, od, do, wykonawca); trafiona {
			return shared.StudioClipboardCopyResponse{},
				kontrolaBladBlokady(blokada, wykonawca.nazwaWykonawcy())
		}
		postacRozetnij(&stan.forma, od, do)
		postacZamienTresc(&stan.forma, od, do, "", nil, nil)
		pusty := ""
		zmiana, err := a.postacOdlozZmiane(ctx, stan, wykonawca.Rodzaj,
			shared.StudioChangeKindUsuniecie, od, do, &tekst, &pusty)
		if err != nil {
			return shared.StudioClipboardCopyResponse{}, err
		}
		if err := a.postacZapisz(ctx, stan); err != nil {
			return shared.StudioClipboardCopyResponse{}, err
		}
		drzewo, err := dziennikZapisDrzewa(stan.forma)
		if err != nil {
			return shared.StudioClipboardCopyResponse{}, err
		}
		var kodZmiany *string
		if zmiana != nil {
			kodZmiany = &zmiana.Id
		}
		czynnosc, err := a.dziennikOdlozCzynnosc(ctx, stan.dokument, wykonawca,
			shared.StudioActionKindTextEdit,
			"wycięcie fragmentu od znaku "+strconv.Itoa(od)+" do "+strconv.Itoa(do)+
				" do schowka; wyciął: "+wykonawca.nazwaWykonawcy(),
			&od, &do, drzewo, drzewo, kodZmiany)
		if err != nil {
			return shared.StudioClipboardCopyResponse{}, err
		}
		if err := a.zmianyModeluStempluj(ctx, zmiana, wykonawca, czynnosc); err != nil {
			return shared.StudioClipboardCopyResponse{}, err
		}
		forma := stan.forma
		odpowiedz.Form = &forma
		odpowiedz.ActionId = czynnosc
	}
	return odpowiedz, nil
}

// WklejZeSchowka obsługuje `studio.clipboard.paste`: wnosi treść do dokumentu
// sposobem `keepFormat`, `plainText` albo `mergeFormat`, wybieranym w żądaniu;
// brak wskazania znaczy zachowanie postaci źródła.
func (a *adapterStudia) WklejZeSchowka(ctx context.Context,
	z shared.StudioClipboardPasteRequest) (shared.StudioClipboardPasteResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioClipboardPasteResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioClipboardPasteResponse{}, err
	}
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	})

	sposob := shared.StudioPasteMode(shared.StudioPasteModeKeepFormat)
	if z.Mode != nil && strings.TrimSpace(string(*z.Mode)) != "" {
		sposob = *z.Mode
		if err := schowekSprawdzSposob(sposob); err != nil {
			return shared.StudioClipboardPasteResponse{}, err
		}
	}

	// Treść wklejana: wprost z żądania, ze wpisu wskazanego albo najświeższego;
	// schowek pusty jest odmową.
	tekst := ""
	zeSchowka := ""
	// postacZrodla jest postacią fragmentu odłożonego; treść podana wprost
	// w żądaniu postaci nie ma.
	var postacZrodla *schowekPostacFragmentu
	switch {
	case z.Text != nil:
		tekst = *z.Text
	case kontrolaTekstNiepusty(z.ClipboardEntryId):
		wpis, err := skladnica.WpisSchowkaStudia(ctx, strings.TrimSpace(*z.ClipboardEntryId))
		if err != nil {
			return shared.StudioClipboardPasteResponse{}, dziennikBrakWiersza(err,
				"historia schowka nie zna wpisu "+*z.ClipboardEntryId)
		}
		if wpis.Rodzaj != shared.ClipboardEntryKindText {
			return shared.StudioClipboardPasteResponse{}, bladWskazaniaStudio("wpis schowka " +
				wpis.Kod + " jest rodzaju „" + wpis.Rodzaj + "”, a nie tekstem — obraz " +
				"i plik wnosi się do dokumentu przez rodzinę studio.insert.*")
		}
		tekst, zeSchowka = wpis.Tresc, wpis.Kod
		postacZrodla, err = schowekPostacZWpisu(wpis)
		if err != nil {
			return shared.StudioClipboardPasteResponse{}, err
		}
	default:
		wpis, err := skladnica.NajswiezszyWpisSchowkaStudia(ctx)
		if err != nil {
			return shared.StudioClipboardPasteResponse{}, dziennikBrakWiersza(err,
				"historia schowka jest pusta — nie ma czego wkleić. Odłóż fragment "+
					"przez studio.clipboard.copy albo podaj treść wprost")
		}
		tekst, zeSchowka = wpis.Tresc, wpis.Kod
		postacZrodla, err = schowekPostacZWpisu(wpis)
		if err != nil {
			return shared.StudioClipboardPasteResponse{}, err
		}
	}
	if tekst == "" {
		return shared.StudioClipboardPasteResponse{}, bladWskazaniaStudio(
			"wklejenie pustej treści nie jest wklejeniem")
	}

	dlugosc := postacDlugosc(&stan.forma)
	miejsce := z.Offset
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > dlugosc {
		miejsce = dlugosc
	}
	od, do := miejsce, miejsce
	if z.ReplaceRangeStart != nil && z.ReplaceRangeEnd != nil {
		od, do = postacZakres(z.ReplaceRangeStart, z.ReplaceRangeEnd, dlugosc)
	}

	// Blokada sprawdza się tu, przed dotknięciem treści: żądanie niesie miejsce
	// wklejenia pod inną nazwą.
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, wykonawca.Rodzaj)
	if len(odcinki) == 0 || len(pominiete) > 0 && od == do {
		return shared.StudioClipboardPasteResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodePermissionDenied, "moduł Studio: "+wykonawca.nazwaWykonawcy()+
				" nie może wkleić w miejsce "+strconv.Itoa(miejsce)+" — stoi na nim blokada "+
				postacNazwaBlokad(pominiete)+". Fragment wymagający zmiany wskazuje się "+
				"propozycją na marginesie (studio.markup.add o rodzaju suggestion)."))
	}

	zastane := znakowanieTekstZakresu(&stan.forma, od, do)

	// Sposób wklejenia rozstrzyga, jaką postać dostaje treść wniesiona.

	// postacZeZrodla mówi, czy postać naprawdę przyszła ze źródła, od czego
	// zależy treść bilansu.
	var postacDocelowa *shared.StudioCharacterFormat
	var akapitZrodla *shared.StudioParagraphFormat
	postacZeZrodla := false
	postacMiejsca := func() *shared.StudioCharacterFormat {
		postacRozetnij(&stan.forma, od, do)
		if wskazania := postacFragmentyZakresu(&stan.forma, od, do); len(wskazania) > 0 {
			znak := postacZnakSkuteczny(&stan.forma,
				stan.forma.Blocks[wskazania[0][0]].Runs[wskazania[0][1]].Format)
			return &znak
		}
		return nil
	}
	switch sposob {
	case shared.StudioPasteModeMergeFormat:
		postacDocelowa = postacMiejsca()
	case shared.StudioPasteModePlainText:
		// Czysty tekst wchodzi bez postaci własnej: run bez formatu bierze postać
		// z arkusza stylów akapitu.
		postacDocelowa = nil
	default:
		// Zachowanie postaci źródła: postać jedzie razem z wpisem schowka
		// (migracja 390).

		// Wpis bez postaci — odłożony przed dobudową albo innym oknem platformy —
		// wraca do postaci miejsca.
		if postacZrodla != nil && postacZrodla.Znak != nil {
			postacRozetnij(&stan.forma, od, do)
			postacDocelowa = postacZrodla.Znak
			akapitZrodla = postacZrodla.Akapit
			postacZeZrodla = true
		} else {
			postacDocelowa = postacMiejsca()
		}
	}

	postacRozetnij(&stan.forma, od, do)
	postacZamienTresc(&stan.forma, od, do, tekst, postacDocelowa, nil)
	koniec := od + len([]rune(tekst))

	// Postać akapitu źródła nakłada się osobno, bo `postacZamienTresc` niesie
	// postać znaku, nie akapitu.

	// Nakłada się wyłącznie przy „zachowaj postać źródła” i wyłącznie na bloki
	// objęte wklejeniem.
	if akapitZrodla != nil {
		for _, wskazanie := range postacBlokiZakresu(&stan.forma, od, koniec) {
			stan.forma.Blocks[wskazanie].Paragraph = postacKopiaAkapitu(akapitZrodla)
		}
	}

	zmiana, err := a.postacOdlozZmiane(ctx, stan, wykonawca.Rodzaj,
		shared.StudioChangeKindWstawienie, od, koniec, &zastane, &tekst)
	if err != nil {
		return shared.StudioClipboardPasteResponse{}, err
	}
	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioClipboardPasteResponse{}, err
	}
	drzewo, err := dziennikZapisDrzewa(stan.forma)
	if err != nil {
		return shared.StudioClipboardPasteResponse{}, err
	}
	var kodZmiany *string
	if zmiana != nil {
		kodZmiany = &zmiana.Id
	}
	czynnosc, err := a.dziennikOdlozCzynnosc(ctx, stan.dokument, wykonawca,
		shared.StudioActionKindTextEdit,
		"wklejenie ze schowka w miejsce "+strconv.Itoa(od)+" sposobem „"+string(sposob)+
			"”; wkleił: "+wykonawca.nazwaWykonawcy(),
		&od, &koniec, drzewo, drzewo, kodZmiany)
	if err != nil {
		return shared.StudioClipboardPasteResponse{}, err
	}
	if err := a.zmianyModeluStempluj(ctx, zmiana, wykonawca, czynnosc); err != nil {
		return shared.StudioClipboardPasteResponse{}, err
	}

	// Pochodzenie: fragment wniesiony ze schowka niesie zapis, skąd jest.
	if zeSchowka != "" {
		if _, err := skladnica.ZapiszPochodzenieFragmentu(ctx,
			dane.PochodzenieFragmentuStudia{
				Kod:             nowyIdentyfikator(przedrostekPochodzeniaStudia),
				DokumentID:      stan.dokument.ID,
				Rodzaj:          shared.StudioProvenanceKindClipboard,
				ZakresOd:        int64(od),
				ZakresDo:        int64(koniec),
				WersjaZrodla:    &zeSchowka,
				AutorRodzaj:     string(wykonawca.Rodzaj),
				AutorAgentKod:   wykonawca.AgentKod,
				AutorAgentNazwa: wykonawca.AgentNazwa,
			}); err != nil {

			return shared.StudioClipboardPasteResponse{}, bladStudio(err)
		}
	}

	bilans := shared.StudioActionBalance{
		Applied: 1, SkippedCount: len(pominiete), Skipped: pominiete,
	}
	uwaga := "Wklejono " + strconv.Itoa(len([]rune(tekst))) + " znaków sposobem „" +
		string(sposob) + "”."
	if sposob == shared.StudioPasteModeKeepFormat {
		switch {
		case postacZeZrodla && akapitZrodla != nil:
			uwaga += " Zachowana została postać ŹRÓDŁA — postać znaku i postać akapitu " +
				"fragmentu odłożonego do schowka."
		case postacZeZrodla:
			uwaga += " Zachowana została postać znaku ze ŹRÓDŁA odłożonego do schowka."
		default:
			uwaga += " Wpis schowka nie niesie postaci źródła — odłożyło go okno, " +
				"które postaci nie zapisuje, albo odłożono go przed dobudową kolumny " +
				"postaci. Zachowaną postacią jest więc postać miejsca wklejenia; postać " +
				"fragmentu źródłowego przenosi malarz formatów " +
				"(studio.format.painter.copy i .apply)."
		}
	}
	bilans.Note = &uwaga
	if bilans.Skipped == nil {
		bilans.Skipped = []shared.StudioSkippedItem{}
	}
	return shared.StudioClipboardPasteResponse{
		Form:     stan.forma,
		Balance:  bilans,
		Change:   zmiana,
		ActionId: czynnosc,
	}, nil
}

// PochodzenieFragmentow obsługuje `studio.provenance.list`: zwraca wykaz
// zapisanych pochodzeń fragmentów dokumentu, przefiltrowany opcjonalnie
// rodzajem pochodzenia i zakresem znaków.
func (a *adapterStudia) PochodzenieFragmentow(ctx context.Context,
	z shared.StudioProvenanceListRequest) (shared.StudioProvenanceListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioProvenanceListResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioProvenanceListResponse{}, err
	}
	wiersze, err := skladnica.PochodzenieFragmentow(ctx, dokument.ID)
	if err != nil {
		return shared.StudioProvenanceListResponse{}, bladStudio(err)
	}
	odpowiedz := shared.StudioProvenanceListResponse{Entries: []shared.StudioProvenance{}}
	for _, wiersz := range wiersze {
		if z.Kind != nil && *z.Kind != "" && wiersz.Rodzaj != string(*z.Kind) {
			continue
		}
		if z.RangeStart != nil && z.RangeEnd != nil &&
			!kontrolaZakresyStykaja(*z.RangeStart, *z.RangeEnd,
				int(wiersz.ZakresOd), int(wiersz.ZakresDo)) {

			continue
		}
		odpowiedz.Entries = append(odpowiedz.Entries,
			schowekZlozPochodzenie(dokument.Kod, wiersz))
	}
	return odpowiedz, nil
}

// ── Postać fragmentu w schowku ──────────────────────────────────────────────

// schowekPostacFragmentu to postać fragmentu odłożonego do schowka — postać
// znaku i postać akapitu początku zakresu, ta sama, jaką przenosi malarz
// formatów.
type schowekPostacFragmentu struct {
	Znak   *shared.StudioCharacterFormat `json:"znak,omitempty"`
	Akapit *shared.StudioParagraphFormat `json:"akapit,omitempty"`
}

// schowekPostacZakresu zbiera postać początku zakresu do zapisu przy wpisie
// schowka. Zakres bez ani jednego fragmentu daje brak postaci, a nie odmowę:
// treść może pochodzić z drzewa, w którym postaci jeszcze nie ma.
func schowekPostacZakresu(forma *shared.StudioDocumentForm, od, do int) (*string, error) {
	postac := schowekPostacFragmentu{}
	if wskazania := postacFragmentyZakresu(forma, od, do); len(wskazania) > 0 {
		znak := postacZnakSkuteczny(forma,
			forma.Blocks[wskazania[0][0]].Runs[wskazania[0][1]].Format)
		postac.Znak = &znak
	}
	if bloki := postacBlokiZakresu(forma, od, do); len(bloki) > 0 {
		akapit := postacAkapitSkuteczny(forma, forma.Blocks[bloki[0]])
		postac.Akapit = &akapit
	}
	if postac.Znak == nil && postac.Akapit == nil {
		return nil, nil
	}
	zapis, err := json.Marshal(postac)
	if err != nil {
		return nil, postacBladZaplecza(
			"postaci odkładanego fragmentu nie da się zapisać: " + err.Error())
	}
	return postacWskaznikTekstu(string(zapis)), nil
}

// schowekPostacZWpisu odczytuje postać źródła z wpisu historii schowka.
//
// Zapis nieczytelny jest USTERKĄ ZAPLECZA, nie brakiem postaci: kolumna niesie
// treść, więc postać ktoś zapisał, a ciche wklejenie postacią miejsca byłoby
// przemilczeniem usterki.
func schowekPostacZWpisu(wpis dane.WpisSchowka) (*schowekPostacFragmentu, error) {
	if wpis.PostacJSON == nil || strings.TrimSpace(*wpis.PostacJSON) == "" {
		return nil, nil
	}
	var postac schowekPostacFragmentu
	if err := json.Unmarshal([]byte(*wpis.PostacJSON), &postac); err != nil {
		return nil, postacBladZaplecza("postać źródła wpisu schowka " + wpis.Kod +
			" jest nieczytelna: " + err.Error())
	}
	return &postac, nil
}

// ── Wspólne ─────────────────────────────────────────────────────────────────

// schowekSprawdzSposob odbija sposób wklejenia spoza wyliczenia kontraktu,
// oddając odmowę z wykazem sposobów dozwolonych zamiast przepuszczać
// nieznaną wartość dalej.
func schowekSprawdzSposob(sposob shared.StudioPasteMode) error {
	for _, dozwolony := range shared.WartosciStudioPasteMode() {
		if sposob == dozwolony {
			return nil
		}
	}
	nazwy := make([]string, 0, 3)
	for _, dozwolony := range shared.WartosciStudioPasteMode() {
		nazwy = append(nazwy, string(dozwolony))
	}
	return bladWskazaniaStudio("sposób wklejenia „" + string(sposob) +
		"” nie jest znany — wolno: " + strings.Join(nazwy, ", "))
}

// schowekZlozPochodzenie składa zapis pochodzenia kontraktu z wiersza bazy,
// dodając do niego kod dokumentu, którego wiersz nie niesie sam.
func schowekZlozPochodzenie(kodDokumentu string,
	wiersz dane.PochodzenieFragmentuStudia) shared.StudioProvenance {

	zapis := shared.StudioProvenance{
		Id:            wiersz.Kod,
		DocumentId:    kodDokumentu,
		Kind:          shared.StudioProvenanceKind(wiersz.Rodzaj),
		RangeStart:    int(wiersz.ZakresOd),
		RangeEnd:      int(wiersz.ZakresDo),
		SourceUrl:     wiersz.AdresZrodla,
		LibraryFileId: wiersz.BibliotekaPlikKod,
		SourceVersion: wiersz.WersjaZrodla,
		SourceTitle:   wiersz.TytulZrodla,
	}
	if wiersz.Siegnieto != nil && *wiersz.Siegnieto != "" {
		chwila := chwilaBazy(*wiersz.Siegnieto)
		zapis.RetrievedAt = &chwila
	}
	if wiersz.AutorRodzaj != "" {
		autor := shared.StudioAuthor(wiersz.AutorRodzaj)
		zapis.Author = &autor
	}
	return zapis
}
