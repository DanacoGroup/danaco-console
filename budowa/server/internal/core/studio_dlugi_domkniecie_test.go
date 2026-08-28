package core

import (
	"context"
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ── Punkt pierwszy: paginacja z nastaw sekcji ───────────────────────────────

// TestPaginacjaLiczySieZNastawStrony wykazuje, że rachunek stron idzie z nastaw
// dokumentu, a nie z zaszytej kartki A4: ta sama treść na kartce A5 zajmuje
// więcej stron niż na A4, bo kartka jest mniejsza.
func TestPaginacjaLiczySieZNastawStrony(t *testing.T) {
	akapity := make([]string, 0, 80)
	for i := 0; i < 80; i++ {
		akapity = append(akapity,
			"Akapit o treści dostatecznie długiej, żeby zająć wiersz na kartce.")
	}
	tresc := strings.Join(akapity, "\n")

	a4 := shared.StudioPageSetup{PageSize: postacWskaznikTekstu("A4")}
	a5 := shared.StudioPageSetup{PageSize: postacWskaznikTekstu("A5")}
	pozioma := shared.StudioPageSetup{
		PageSize:    postacWskaznikTekstu("A4"),
		Orientation: postacWskaznikOrientacjiStrony(shared.StudioPageOrientationPozioma),
	}

	_, stronA4, err := aparatStronyAkapitow(tresc, &shared.StudioDocumentForm{PageSetup: &a4})
	if err != nil {
		t.Fatalf("rachunek stron A4 odmówił: %v", err)
	}
	_, stronA5, err := aparatStronyAkapitow(tresc, &shared.StudioDocumentForm{PageSetup: &a5})
	if err != nil {
		t.Fatalf("rachunek stron A5 odmówił: %v", err)
	}
	_, stronPoziomo, err := aparatStronyAkapitow(tresc,
		&shared.StudioDocumentForm{PageSetup: &pozioma})
	if err != nil {
		t.Fatalf("rachunek stron A4 poziomo odmówił: %v", err)
	}

	if stronA5 <= stronA4 {
		t.Errorf("ta sama treść zajmuje na A5 %d stron, a na A4 %d — kartka mniejsza "+
			"musi dać stron więcej, inaczej nastawy strony nie dochodzą do rachunku",
			stronA5, stronA4)
	}
	// Miarą orientacji jest różnica liczby stron, bo kierunek przewagi zależy
	// od treści akapitów.
	if stronPoziomo == stronA4 {
		t.Errorf("ta sama treść zajmuje poziomo i pionowo tyle samo stron (%d) — "+
			"orientacja nie dochodzi do rachunku stron", stronA4)
	}

	// Marginesy szerokie zabierają miejsce, więc stron ma być więcej niż przy
	// marginesach wąskich.
	waskie := shared.StudioPageSetup{
		PageSize: postacWskaznikTekstu("A4"),
		// Marginesy w kontrakcie są w milimetrach.
		MarginTop: postacWskaznikLiczby(10), MarginBottom: postacWskaznikLiczby(10),
		MarginLeft: postacWskaznikLiczby(10), MarginRight: postacWskaznikLiczby(10),
	}
	szerokie := shared.StudioPageSetup{
		PageSize:  postacWskaznikTekstu("A4"),
		MarginTop: postacWskaznikLiczby(50), MarginBottom: postacWskaznikLiczby(50),
		MarginLeft: postacWskaznikLiczby(50), MarginRight: postacWskaznikLiczby(50),
	}
	_, stronWaskie, err := aparatStronyAkapitow(tresc,
		&shared.StudioDocumentForm{PageSetup: &waskie})
	if err != nil {
		t.Fatalf("rachunek stron przy marginesach wąskich odmówił: %v", err)
	}
	_, stronSzerokie, err := aparatStronyAkapitow(tresc,
		&shared.StudioDocumentForm{PageSetup: &szerokie})
	if err != nil {
		t.Fatalf("rachunek stron przy marginesach szerokich odmówił: %v", err)
	}
	if stronSzerokie <= stronWaskie {
		t.Errorf("marginesy szerokie dają %d stron, a wąskie %d — margines nie dochodzi "+
			"do rachunku stron", stronSzerokie, stronWaskie)
	}
}

// TestPaginacjaLiczySekcjeOsobno wykazuje, że sekcja o własnym nośniku liczy się
// własną kartką, a sekcja od nowej strony przerywa rachunek wierszy.
func TestPaginacjaLiczySekcjeOsobno(t *testing.T) {
	pierwszy := "Akapit pierwszej sekcji."
	drugi := "Akapit drugiej sekcji."
	tresc := pierwszy + "\n" + drugi

	nowaStrona := shared.StudioSectionStart(shared.StudioSectionStartNewPage)
	ciagla := shared.StudioSectionStart(shared.StudioSectionStartContinuous)
	nastawyDrugiej := shared.StudioPageSetup{PageSize: postacWskaznikTekstu("A5")}

	formaZLamaniem := shared.StudioDocumentForm{
		PageSetup: &shared.StudioPageSetup{PageSize: postacWskaznikTekstu("A4")},
		Sections: []shared.StudioSection{
			{Id: "sek-1", Index: 0, RangeStart: 0, RangeEnd: len([]rune(pierwszy)) + 1},
			{Id: "sek-2", Index: 1, RangeStart: len([]rune(pierwszy)) + 1,
				RangeEnd: len([]rune(tresc)) + 1, Start: &nowaStrona,
				PageSetup: &nastawyDrugiej},
		},
	}
	strony, stron, err := aparatStronyAkapitow(tresc, &formaZLamaniem)
	if err != nil {
		t.Fatalf("rachunek stron z sekcjami odmówił: %v", err)
	}
	if len(strony) != 2 {
		t.Fatalf("rachunek oddał %d numerów stron dla dwóch akapitów", len(strony))
	}
	if strony[0] != 1 {
		t.Errorf("pierwszy akapit stoi na stronie %d, a ma stać na pierwszej", strony[0])
	}
	if strony[1] != 2 {
		t.Errorf("akapit sekcji rozpoczynanej od nowej strony stoi na stronie %d, "+
			"a ma stać na drugiej", strony[1])
	}
	if stron != 2 {
		t.Errorf("dokument z sekcją od nowej strony ma %d stron, a ma mieć dwie", stron)
	}

	// Sekcja ciągła łamania NIE przerywa — jej sensem jest ciągłość.
	formaCiagla := formaZLamaniem
	formaCiagla.Sections = []shared.StudioSection{
		formaZLamaniem.Sections[0],
		{Id: "sek-2", Index: 1, RangeStart: len([]rune(pierwszy)) + 1,
			RangeEnd: len([]rune(tresc)) + 1, Start: &ciagla, PageSetup: &nastawyDrugiej},
	}
	stronyCiagle, stronCiagle, err := aparatStronyAkapitow(tresc, &formaCiagla)
	if err != nil {
		t.Fatalf("rachunek stron z sekcją ciągłą odmówił: %v", err)
	}
	if stronyCiagle[1] != 1 || stronCiagle != 1 {
		t.Errorf("sekcja ciągła przerwała stronę: akapit na stronie %d, stron %d",
			stronyCiagle[1], stronCiagle)
	}
}

// TestGeometriaNieZjadaKartkiMarginesami pilnuje odmowy nieoczywistej: nastawa,
// która zabiera całą kartkę marginesami, nie ma prawa doprowadzić rdzenia do
// dzielenia przez zero ani do kolumny o zerowej szerokości.
func TestGeometriaNieZjadaKartkiMarginesami(t *testing.T) {
	nastawy := shared.StudioPageSetup{
		PageSize:  postacWskaznikTekstu("A6"),
		MarginTop: postacWskaznikLiczby(300), MarginBottom: postacWskaznikLiczby(300),
		MarginLeft: postacWskaznikLiczby(300), MarginRight: postacWskaznikLiczby(300),
	}
	kartka := geometriaZNastawStrony(&nastawy)
	if kartka.szerokoscKolumny() < 1 {
		t.Errorf("kolumna tekstu ma szerokość %d — rachunek łamania nie ma na czym stanąć",
			kartka.szerokoscKolumny())
	}
	if kartka.wierszyNaStrone() < 1 {
		t.Errorf("na stronie mieści się %d wierszy — rachunek stron nigdy się nie skończy",
			kartka.wierszyNaStrone())
	}
}

// ── Punkt drugi: aparat dokumentu cofa się dziennikiem ─────────────────────

// TestAparatCofaSieDziennikiem wykazuje, że cofnięcie czynności zdejmuje
// element aparatu wniesiony tą czynnością z drzewa postaci i z wiersza aparatu,
// bo jeden bez drugiego zostawia dokument niespójny.
func TestAparatCofaSieDziennikiem(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-aparat-cofania",
		"Zdanie pierwsze.\nZdanie drugie.")
	adapter := nowyAdapterStudia(zmontowany.dane.Studio)

	przypis := shared.StudioApparatusKind(shared.StudioApparatusKindFootnote)
	wstawiony, err := adapter.WstawElementAparatu(zycie,
		shared.StudioApparatusInsertRequest{
			DocumentId: dokument.Id,
			Kind:       przypis,
			Offset:     postacWskaznikLiczby(16),
			Text:       postacWskaznikTekstu("Powołanie na źródło."),
		})
	if err != nil {
		t.Fatalf("wstawienie przypisu odmówiło: %v", err)
	}
	if wstawiony.ActionId == nil || *wstawiony.ActionId == "" {
		t.Fatal("wstawienie przypisu nie oddało wpisu dziennika — nie ma czego cofnąć")
	}

	przed, err := adapter.WykazAparatu(zycie,
		shared.StudioApparatusListRequest{DocumentId: dokument.Id})
	if err != nil {
		t.Fatalf("wykaz aparatu odmówił: %v", err)
	}
	if len(przed.Items) == 0 {
		t.Fatal("aparat dokumentu nie niesie wstawionego przypisu")
	}

	cofniete, err := adapter.CofnijCzynnosc(zycie, shared.StudioJournalRevertRequest{
		DocumentId: dokument.Id,
		ActionIds:  []string{*wstawiony.ActionId},
	})
	if err != nil {
		t.Fatalf("cofnięcie czynności odmówiło: %v", err)
	}
	if len(cofniete.Reverted) != 1 {
		t.Fatalf("cofnięcie objęło %d czynności, a wskazano jedną", len(cofniete.Reverted))
	}

	// Bilans nie ma prawa nazywać braku, którego już nie ma.
	for _, pominiecie := range cofniete.Balance.Skipped {
		if strings.Contains(pominiecie.Reason, "aparat") {
			t.Errorf("cofnięcie nadal nazywa aparat jako niecofnięty: %s", pominiecie.Reason)
		}
	}
	if len(cofniete.Form.Apparatus) != 0 {
		t.Errorf("drzewo postaci po cofnięciu niesie %d elementów aparatu, a ma nie nieść "+
			"ani jednego", len(cofniete.Form.Apparatus))
	}

	// Wiersze: odczyt osobnym wywołaniem, bo Operator otworzy dokument na nowo.
	po, err := adapter.WykazAparatu(zycie,
		shared.StudioApparatusListRequest{DocumentId: dokument.Id})
	if err != nil {
		t.Fatalf("wykaz aparatu po cofnięciu odmówił: %v", err)
	}
	if len(po.Items) != 0 {
		t.Errorf("wiersze aparatu po cofnięciu niosą %d elementów — cofnięcie zdjęło "+
			"element z drzewa, ale nie z bazy", len(po.Items))
	}
}

// ── Punkt piąty: porządek alfabetyczny pisma polskiego ─────────────────────

// TestSortowanieZnaPismoPolskie wykazuje kolację: „ł" stoi za „l" i przed „m",
// a nie za „z", gdzie leży w Unikodzie.
func TestSortowanieZnaPismoPolskie(t *testing.T) {
	zestawiacz := tabelaZestawiaczPolski()

	// Wielkie „Ł" ma w Unikodzie punkt kodowy wyższy od „Z", stąd znaczenie
	// kolacji polskiej.
	if !tabelaMniejszyTekst(zestawiacz, "Łukasiewicz", "Zawadzki") {
		t.Error("„Łukasiewicz” nie wyprzedza „Zawadzkiego” — porządek nie zna pisma polskiego")
	}
	if !tabelaMniejszyTekst(zestawiacz, "lampa", "łaska") {
		t.Error("„lampa” nie wyprzedza „łaski”")
	}
	if !tabelaMniejszyTekst(zestawiacz, "łaska", "mapa") {
		t.Error("„łaska” nie wyprzedza „mapy”")
	}
	if !tabelaMniejszyTekst(zestawiacz, "ćma", "dom") {
		t.Error("„ćma” nie wyprzedza „domu” — wielkość liter nie ma rozstrzygać")
	}
	if !tabelaMniejszyTekst(zestawiacz, "Ćma", "dom") {
		t.Error("„Ćma” nie wyprzedza „domu” — wielkość liter nie ma rozstrzygać")
	}
	// Ogonek rozstrzyga, gdy litery bazowe są te same — „laska" i „łaska" stoją
	// osobno.
	if tabelaMniejszyTekst(zestawiacz, "łaska", "laska") {
		t.Error("„łaska” wyprzedza „laskę” — znak diakrytyczny idzie po literze bazowej")
	}
}

// TestSortowanieTabeliUkladaNazwiskaPoPolsku mierzy tę samą rzecz na TABELI,
// czyli tam, gdzie Operator ją widzi.
func TestSortowanieTabeliUkladaNazwiskaPoPolsku(t *testing.T) {
	nazwiska := []string{"Zawadzki", "Łukasiewicz", "Śliwa", "Adamczyk", "Żuk", "Ćwik"}
	tabela := shared.StudioDocumentTable{
		Id: "tab-nazwiska", Rows: len(nazwiska), Columns: 1,
	}
	tabelaSiatkaPelna(&tabela)
	for i, nazwisko := range nazwiska {
		komorka := tabelaKomorka(&tabela, i, 0)
		komorka.Text = postacWskaznikTekstu(nazwisko)
	}

	if _, err := tabelaSortujWiersze(&tabela, 0, false, false, false); err != nil {
		t.Fatalf("sortowanie tabeli odmówiło: %v", err)
	}
	ulozone := make([]string, 0, len(nazwiska))
	for i := 0; i < tabela.Rows; i++ {
		komorka := tabelaKomorka(&tabela, i, 0)
		if komorka != nil && komorka.Text != nil {
			ulozone = append(ulozone, *komorka.Text)
		}
	}
	oczekiwane := []string{"Adamczyk", "Ćwik", "Łukasiewicz", "Śliwa", "Zawadzki", "Żuk"}
	if strings.Join(ulozone, ",") != strings.Join(oczekiwane, ",") {
		t.Errorf("tabela ułożona jako %v, a porządkiem alfabetycznym pisma polskiego "+
			"jest %v", ulozone, oczekiwane)
	}
}

// ── Punkt trzeci: malarz formatów w tabeli ──────────────────────────────────

// TestMalarzFormatowPrzezywaPrzeladowanieRdzenia wykazuje, że postać zabrana
// malarzem formatów przeżywa przeładowanie rdzenia: adapter złożony od nowa nie
// dzieli zmiennych pakietu z poprzednim, a położenie udaje się tylko, gdy
// postać leży w wierszu.
func TestMalarzFormatowPrzezywaPrzeladowanieRdzenia(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-malarza",
		"Zdanie wzorcowe.\nZdanie do umalowania.")

	zabierajacy := nowyAdapterStudia(zmontowany.dane.Studio)
	prawda := true
	zabrana, err := zabierajacy.ZabierzPostac(zycie, shared.StudioFormatPainterCopyRequest{
		DocumentId:       dokument.Id,
		RangeStart:       postacWskaznikLiczby(0),
		RangeEnd:         postacWskaznikLiczby(16),
		IncludeParagraph: &prawda,
	})
	if err != nil {
		t.Fatalf("zabranie postaci odmówiło: %v", err)
	}
	if zabrana.ClipId == "" {
		t.Fatal("zabranie postaci nie oddało uchwytu")
	}

	// Rdzeń wstaje od nowa: nowy adapter nad tą samą bazą.
	kladacy := nowyAdapterStudia(zmontowany.dane.Studio)
	polozona, err := kladacy.PolozPostac(zycie, shared.StudioFormatPainterApplyRequest{
		DocumentId: dokument.Id,
		ClipId:     zabrana.ClipId,
		RangeStart: postacWskaznikLiczby(17),
		RangeEnd:   postacWskaznikLiczby(37),
	})
	if err != nil {
		t.Fatalf("położenie postaci po przeładowaniu rdzenia odmówiło: %v — postać "+
			"zabrana malarzem nie przeżyła, choć migracja 368 ma na nią tabelę", err)
	}
	if polozona.Balance.Applied == 0 {
		t.Error("położenie postaci nie objęło ani jednego fragmentu")
	}
}

// TestMalarzOdmawiaNazywajacBrakUchwytu pilnuje odmowy: uchwyt, którego nie ma,
// nie ma prawa położyć postaci przypadkowej ani przemilczeć braku.
func TestMalarzOdmawiaNazywajacBrakUchwytu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-malarza-braku", "Zdanie.")
	adapter := nowyAdapterStudia(zmontowany.dane.Studio)

	_, err := adapter.PolozPostac(zycie, shared.StudioFormatPainterApplyRequest{
		DocumentId: dokument.Id,
		ClipId:     "studio-mal-ktorego-nie-ma",
		RangeStart: postacWskaznikLiczby(0),
		RangeEnd:   postacWskaznikLiczby(7),
	})
	if err == nil {
		t.Fatal("położenie postaci o uchwycie nieistniejącym zameldowało powodzenie")
	}
	if !strings.Contains(err.Error(), "studio-mal-ktorego-nie-ma") {
		t.Errorf("odmowa nie nazywa uchwytu, o który szło: %v", err)
	}
}

// ── Punkt szósty: postać przez schowek ──────────────────────────────────────

// TestSchowekPrzenosiPostacZrodla wykazuje, że wklejenie „zachowaj postać
// źródła" przenosi postać fragmentu źródłowego, a nie postać miejsca
// wklejenia: fragment dostaje wytłuszczenie, które po wklejeniu stoi w miejscu
// wklejenia.
func TestSchowekPrzenosiPostacZrodla(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	zrodlowy := "Fragment wytluszczony."
	dalszy := "Zdanie zwykle."
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-schowka",
		zrodlowy+"\n"+dalszy)
	adapter := nowyAdapterStudia(zmontowany.dane.Studio)

	// Wytłuszczenie fragmentu źródłowego idzie drogą kontraktu, nie zapisem
	// z boku sprawdzianu.
	prawda := true
	if _, err := adapter.UstawPostacZnaku(zycie, shared.StudioFormatCharacterSetRequest{
		DocumentId: dokument.Id,
		RangeStart: postacWskaznikLiczby(0),
		RangeEnd:   postacWskaznikLiczby(len([]rune(zrodlowy))),
		Bold:       &prawda,
	}); err != nil {
		t.Fatalf("wytłuszczenie fragmentu źródłowego odmówiło: %v", err)
	}

	skopiowane, err := adapter.SkopiujDoSchowka(zycie, shared.StudioClipboardCopyRequest{
		DocumentId: dokument.Id,
		RangeStart: 0,
		RangeEnd:   len([]rune(zrodlowy)),
	})
	if err != nil {
		t.Fatalf("odłożenie fragmentu do schowka odmówiło: %v", err)
	}
	if skopiowane.ClipboardEntryId == "" {
		t.Fatal("odłożenie nie oddało wpisu schowka")
	}

	// Wklejenie na końcu dokumentu, czyli w miejscu BEZ wytłuszczenia.
	sposob := shared.StudioPasteMode(shared.StudioPasteModeKeepFormat)
	miejsce := len([]rune(zrodlowy)) + 1 + len([]rune(dalszy))
	wklejone, err := adapter.WklejZeSchowka(zycie, shared.StudioClipboardPasteRequest{
		DocumentId:       dokument.Id,
		ClipboardEntryId: &skopiowane.ClipboardEntryId,
		Offset:           miejsce,
		Mode:             &sposob,
	})
	if err != nil {
		t.Fatalf("wklejenie ze schowka odmówiło: %v", err)
	}

	wytluszczone := false
	for _, blok := range wklejone.Form.Blocks {
		for _, fragment := range blok.Runs {
			if fragment.RangeStart == nil || *fragment.RangeStart < miejsce {
				continue
			}
			if fragment.Format != nil && fragment.Format.Bold != nil && *fragment.Format.Bold {
				wytluszczone = true
			}
		}
	}
	if !wytluszczone {
		t.Error("wklejenie sposobem „zachowaj postać źródła” nie przeniosło " +
			"wytłuszczenia — wpis schowka nie niesie postaci źródła")
	}
	if wklejone.Balance.Note == nil ||
		!strings.Contains(*wklejone.Balance.Note, "ŹRÓDŁA") {

		t.Errorf("bilans nie mówi, że zachowana została postać źródła: %v",
			wklejone.Balance.Note)
	}
}

// ── Punkt czwarty: tożsamość agenta na drodze postaci ──────────────────────

// TestTozsamoscAgentaStemplujeSieNaZmianiePostaci wykazuje, że kod wykonawcy
// dojeżdża do wiersza dziennika i do zmiany śledzonej, bo podpis wchodzi
// kontekstem tą samą drogą, którą wkłada go wpięcie rejestru zapisujące podpis
// wykonawcy studia.
func TestTozsamoscAgentaStemplujeSieNaZmianiePostaci(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-agenta",
		"Zdanie do zmiany postaci.")
	adapter := nowyAdapterStudia(zmontowany.dane.Studio)

	kodAgenta := "agent-redaktor-umow"
	nazwaAgenta := "Redaktor umów"
	model := shared.StudioAuthor(shared.StudioAuthorModel)
	kontekst := kontrolaZapiszWykonawce(zycie, kontrolaWykonawca{
		Rodzaj:     model,
		AgentKod:   &kodAgenta,
		AgentNazwa: &nazwaAgenta,
	})

	prawda := true
	if _, err := adapter.UstawPostacZnaku(kontekst, shared.StudioFormatCharacterSetRequest{
		DocumentId: dokument.Id,
		RangeStart: postacWskaznikLiczby(0),
		RangeEnd:   postacWskaznikLiczby(6),
		Bold:       &prawda,
		Author:     &model,
	}); err != nil {
		t.Fatalf("zmiana postaci przez wykonawcę odmówiła: %v", err)
	}

	// Dziennik: wpis musi nieść kod agenta, nie sam rodzaj autora.
	dziennik, err := adapter.DziennikCzynnosci(zycie,
		shared.StudioJournalListRequest{DocumentId: dokument.Id})
	if err != nil {
		t.Fatalf("wykaz czynności dokumentu odmówił: %v", err)
	}
	stempel := false
	for _, czynnosc := range dziennik.Actions {
		if czynnosc.AuthorAgentId != nil && *czynnosc.AuthorAgentId == kodAgenta {
			stempel = true
			if czynnosc.AuthorAgentName == nil || *czynnosc.AuthorAgentName != nazwaAgenta {
				t.Error("wpis dziennika niesie kod agenta bez jego nazwy")
			}
		}
	}
	if !stempel {
		t.Error("wpis dziennika nie niesie kodu wykonawcy — rozbicie zmian po wykonawcy " +
			"pokaże czynność jako czynność nienazwanego")
	}

	// Zmiany modelu: rozbicie po wykonawcy musi znaleźć tę zmianę pod kodem agenta.
	zmiany, err := adapter.ZmianyModelu(zycie, shared.StudioModelChangesListRequest{
		DocumentId: dokument.Id, AgentId: &kodAgenta,
	})
	if err != nil {
		t.Fatalf("wykaz zmian modelu odmówił: %v", err)
	}
	if len(zmiany.Summary.Changes) == 0 {
		t.Error("zawężenie zmian modelu do kodu agenta nie znalazło ani jednej zmiany — " +
			"tożsamość nie została odbita na zmianie śledzonej")
	}
	// Rozbicie po wykonawcy musi znaleźć agenta nazwanego, a nie jedną pozycję
	// bez nazwy.
	nazwany := false
	for _, wykonawca := range zmiany.Summary.ByAgent {
		if wykonawca.Actor.AgentId != nil && *wykonawca.Actor.AgentId == kodAgenta {
			nazwany = true
		}
	}
	if !nazwany {
		t.Errorf("rozbicie zmian po wykonawcy nie zna agenta %s: %+v",
			kodAgenta, zmiany.Summary.ByAgent)
	}
}

// TestPodpisWykonawcyWchodziZLadunkuZadania wykazuje, że wpięcie rejestru
// czyta podpis z ładunku żądania i wkłada go do kontekstu, co odróżnia tę
// miarę od stemplowania zmiany postaci sprawdzanego osobno.
func TestPodpisWykonawcyWchodziZLadunkuZadania(t *testing.T) {
	var widziany kontrolaWykonawca
	var byl bool
	owinieta := podpisemWykonawcy(
		func(ctx context.Context, _ protocol.Request) protocol.Odpowiedz {
			widziany, byl = kontrolaWykonawcaZKontekstu(ctx)
			return protocol.Odpowiedz{Status: shared.EnvelopeStatusOk}
		})

	ladunek := []byte(`{"documentId":"studio-dok-1","agentId":"agent-korekty",` +
		`"agentName":"Korektor","author":"model"}`)
	owinieta(context.Background(), protocol.Request{
		Komenda: shared.CommandStudioFormatCharacterSet, Ladunek: ladunek,
	})

	if !byl {
		t.Fatal("wpięcie rejestru nie włożyło wykonawcy do kontekstu")
	}
	if widziany.AgentKod == nil || *widziany.AgentKod != "agent-korekty" {
		t.Errorf("kod agenta nie doszedł z ładunku: %v", widziany.AgentKod)
	}
	if widziany.AgentNazwa == nil || *widziany.AgentNazwa != "Korektor" {
		t.Errorf("nazwa agenta nie doszła z ładunku: %v", widziany.AgentNazwa)
	}
	if !widziany.czyWykonawca() {
		t.Error("żądanie z kodem agenta nie zostało rozpoznane jako czynność wykonawcy")
	}

	// Ładunek bez podpisu jedzie dalej jako czynność Operatora, a nie odmowa.
	widziany, byl = kontrolaWykonawca{}, false
	owinieta(context.Background(), protocol.Request{
		Komenda: shared.CommandStudioFormatCharacterSet,
		Ladunek: []byte(`{"documentId":"studio-dok-1"}`),
	})
	if !byl {
		t.Fatal("żądanie bez podpisu nie doszło do obsługiwacza z rozstrzygnięciem o autorze")
	}
	if widziany.czyWykonawca() {
		t.Error("żądanie bez podpisu zostało uznane za czynność wykonawcy")
	}

	// Ładunek nieczytelny nie ma prawa zatrzymać komendy — wpięcie dokłada
	// wiedzę, a nie odmawia.
	odpowiedz := owinieta(context.Background(), protocol.Request{
		Komenda: shared.CommandStudioFormatCharacterSet, Ladunek: []byte(`nie-json`),
	})
	if odpowiedz.Status != shared.EnvelopeStatusOk {
		t.Error("ładunek nieczytelny zatrzymał komendę na wpięciu podpisu")
	}
}
