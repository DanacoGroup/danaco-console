package core

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Skutek aparatu dokumentu: czy spis treści zgadza się z nagłówkami, czy przypis
// przenumerowuje się po wstawieniu przypisu przed nim i czy znacznik
// nieświeżości mówi prawdę.
//
// Szkody, które ten plik ma wykluczyć:
//  1. spis treści oddany jako odświeżony, a niosący nagłówki sprzed zmiany;
//  2. numeracja przypisów nadawana w kolejności ZAPISU, po której przypis
//     wstawiony w środek dokumentu kłamie do końca życia pisma;
//  3. znacznik nieświeżości trzymany w drzewie postaci — drzewo zapisuje się bez
//     aparatu, więc znacznik ginąłby przy pierwszym zapisie;
//  4. odwołanie do elementu, którego dokument nie ma, przyjęte jako założone.

// aparatUprzazSprawdzianu składa adapter modułu i dokument z nagłówkami.
func aparatUprzazSprawdzianu(t *testing.T) (*adapterStudia, context.Context, string) {
	t.Helper()
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-aparat",
		strings.Join([]string{
			"Rozdział pierwszy",
			"Treść rozdziału pierwszego z powołaniem.",
			"Rozdział drugi",
			"Treść rozdziału drugiego.",
		}, "\n"))
	return nowyAdapterStudia(zmontowany.dane.Studio), zycie, dokument.Id
}

// aparatNadajNaglowki oznacza wskazane akapity jako nagłówki wskazanego poziomu.
//
// Idzie to drogą zapisu postaci dokumentu, czyli tą samą, którą przestawia
// nagłówki okno — sprawdzian nie ma prawa wpisywać tego do bazy z boku, bo
// mierzyłby wtedy własny zapis, nie zachowanie rdzenia.
func aparatNadajNaglowki(t *testing.T, adapter *adapterStudia, zycie context.Context,
	dokument string, poziomy map[int]int) {
	t.Helper()

	postac, err := adapter.PostacDokumentu(zycie, shared.StudioDocumentFormGetRequest{
		DocumentId: dokument,
	})
	if err != nil {
		t.Fatalf("odczyt postaci dokumentu odmówił: %v", err)
	}
	bloki := postac.Form.Blocks
	akapit := 0
	for i := range bloki {
		if !postacBlokNiesieTekst(bloki[i]) {
			continue
		}
		wskazanie := akapit
		akapit++
		poziom, jest := poziomy[wskazanie]
		if !jest {
			continue
		}
		bloki[i].Kind = blokPostaciNaglowek
		bloki[i].Paragraph = postacScalAkapit(bloki[i].Paragraph,
			shared.StudioParagraphFormat{OutlineLevel: wskaznik(poziom)})
	}
	zapis, err := json.Marshal(shared.StudioDocumentForm{DocumentId: dokument, Blocks: bloki})
	if err != nil {
		t.Fatalf("nie można złożyć postaci do zapisu: %v", err)
	}
	if _, err := adapter.ZapiszPostacDokumentu(zycie, shared.StudioDocumentFormSaveRequest{
		DocumentId: dokument, Form: zapis,
	}); err != nil {
		t.Fatalf("zapis postaci dokumentu odmówił: %v", err)
	}
}

// aparatElementSprawdzianu odczytuje element aparatu osobnym wywołaniem wykazu.
func aparatElementSprawdzianu(t *testing.T, adapter *adapterStudia, zycie context.Context,
	dokument, element string) shared.StudioApparatusItem {
	t.Helper()

	wykaz, err := adapter.WykazAparatu(zycie, shared.StudioApparatusListRequest{
		DocumentId: dokument,
	})
	if err != nil {
		t.Fatalf("wykaz aparatu odmówił: %v", err)
	}
	for _, pozycja := range wykaz.Items {
		if pozycja.Id == element {
			return pozycja
		}
	}
	t.Fatalf("wykaz aparatu nie niesie elementu %s", element)
	return shared.StudioApparatusItem{}
}

// TestSpisTresciZgadzaSieZNaglowkami jest wymaganiem zlecenia wprost: spis po
// odświeżeniu zgadza się z nagłówkami dokumentu.
func TestSpisTresciZgadzaSieZNaglowkami(t *testing.T) {
	adapter, zycie, dokument := aparatUprzazSprawdzianu(t)
	aparatNadajNaglowki(t, adapter, zycie, dokument, map[int]int{0: 1, 2: 1})

	spis, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindToc, Offset: wskaznik(0),
	})
	if err != nil {
		t.Fatalf("założenie spisu treści odmówiło: %v", err)
	}
	// Spis świeżo założony jest NIEŚWIEŻY: nie zbierał jeszcze niczego i mówi to
	// o sobie, zamiast udawać gotowy.
	if spis.Item.Stale == nil || !*spis.Item.Stale {
		t.Error("spis treści świeżo założony nie jest oznaczony jako wymagający odświeżenia")
	}

	if _, err := adapter.OdswiezAparat(zycie, shared.StudioApparatusRefreshRequest{
		DocumentId: dokument,
	}); err != nil {
		t.Fatalf("odświeżenie aparatu odmówiło: %v", err)
	}

	odswiezony := aparatElementSprawdzianu(t, adapter, zycie, dokument, spis.Item.Id)
	if odswiezony.Stale == nil || *odswiezony.Stale {
		t.Error("spis treści po odświeżeniu nadal stoi jako nieświeży")
	}
	if len(odswiezony.Entries) != 2 {
		t.Fatalf("spis treści niesie %d pozycji, oczekiwano dwóch: %+v",
			len(odswiezony.Entries), odswiezony.Entries)
	}
	if odswiezony.Entries[0].Text != "Rozdział pierwszy" ||
		odswiezony.Entries[1].Text != "Rozdział drugi" {

		t.Errorf("spis treści nie zgadza się z nagłówkami dokumentu: %+v",
			odswiezony.Entries)
	}
	for numer, pozycja := range odswiezony.Entries {
		if pozycja.PageNumber == nil || *pozycja.PageNumber < 1 {
			t.Errorf("pozycja %d spisu treści nie ma numeru strony", numer)
		}
		if pozycja.AnchorOffset == nil {
			t.Errorf("pozycja %d spisu treści nie ma miejsca w treści", numer)
		}
	}
}

// TestPrzypisPrzenumerowujeSiePoWstawieniuPrzedNim jest wymaganiem zlecenia
// wprost: przypis wstawiony PRZED innym przestawia numerację obu.
func TestPrzypisPrzenumerowujeSiePoWstawieniuPrzedNim(t *testing.T) {
	adapter, zycie, dokument := aparatUprzazSprawdzianu(t)

	dalszy, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindFootnote,
		Offset: wskaznik(40), Text: wskaznik("Przypis do rozdziału drugiego."),
	})
	if err != nil {
		t.Fatalf("założenie przypisu dalszego odmówiło: %v", err)
	}
	if _, err := adapter.OdswiezAparat(zycie, shared.StudioApparatusRefreshRequest{
		DocumentId: dokument, Kind: wskaznik(shared.StudioApparatusKind(
			shared.StudioApparatusKindFootnote)),
	}); err != nil {
		t.Fatalf("odświeżenie przypisów odmówiło: %v", err)
	}
	pierwszyStan := aparatElementSprawdzianu(t, adapter, zycie, dokument, dalszy.Item.Id)
	if pierwszyStan.Number == nil || *pierwszyStan.Number != "1" {
		t.Fatalf("jedyny przypis niesie numer %v, oczekiwano 1", pierwszyStan.Number)
	}

	wczesniejszy, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindFootnote,
		Offset: wskaznik(5), Text: wskaznik("Przypis do rozdziału pierwszego."),
	})
	if err != nil {
		t.Fatalf("założenie przypisu wcześniejszego odmówiło: %v", err)
	}
	// Wstawienie przypisu przed innym unieważnia numerację — i mówi to
	// znacznikiem, nie ciszą.
	dalszyPoWstawieniu := aparatElementSprawdzianu(t, adapter, zycie, dokument, dalszy.Item.Id)
	if dalszyPoWstawieniu.Stale == nil || !*dalszyPoWstawieniu.Stale {
		t.Error("przypis dalszy nie został oznaczony jako nieświeży po wstawieniu przypisu " +
			"przed nim")
	}

	if _, err := adapter.OdswiezAparat(zycie, shared.StudioApparatusRefreshRequest{
		DocumentId: dokument,
	}); err != nil {
		t.Fatalf("odświeżenie aparatu odmówiło: %v", err)
	}
	nowyWczesniejszy := aparatElementSprawdzianu(t, adapter, zycie, dokument, wczesniejszy.Item.Id)
	nowyDalszy := aparatElementSprawdzianu(t, adapter, zycie, dokument, dalszy.Item.Id)
	if nowyWczesniejszy.Number == nil || *nowyWczesniejszy.Number != "1" {
		t.Errorf("przypis wcześniejszy niesie numer %v, oczekiwano 1", nowyWczesniejszy.Number)
	}
	if nowyDalszy.Number == nil || *nowyDalszy.Number != "2" {
		t.Errorf("przypis dalszy niesie numer %v, oczekiwano 2 — numeracja nie przestawiła "+
			"się po wstawieniu przypisu przed nim", nowyDalszy.Number)
	}
}

// TestUsunieciePrzypisuPrzenumerowujePozostale mierzy drugą stronę tej samej
// zasady: po usunięciu przypisu numery nie mają dziur.
func TestUsunieciePrzypisuPrzenumerowujePozostale(t *testing.T) {
	adapter, zycie, dokument := aparatUprzazSprawdzianu(t)

	kody := make([]string, 0, 3)
	for _, miejsce := range []int{5, 20, 45} {
		element, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
			DocumentId: dokument, Kind: shared.StudioApparatusKindFootnote,
			Offset: wskaznik(miejsce), Text: wskaznik("Przypis."),
		})
		if err != nil {
			t.Fatalf("założenie przypisu odmówiło: %v", err)
		}
		kody = append(kody, element.Item.Id)
	}
	if _, err := adapter.UsunElementAparatuDokumentu(zycie,
		shared.StudioApparatusRemoveRequest{DocumentId: dokument, ItemId: kody[1]}); err != nil {

		t.Fatalf("usunięcie przypisu odmówiło: %v", err)
	}
	if _, err := adapter.OdswiezAparat(zycie, shared.StudioApparatusRefreshRequest{
		DocumentId: dokument,
	}); err != nil {
		t.Fatalf("odświeżenie aparatu odmówiło: %v", err)
	}
	pierwszy := aparatElementSprawdzianu(t, adapter, zycie, dokument, kody[0])
	trzeci := aparatElementSprawdzianu(t, adapter, zycie, dokument, kody[2])
	if pierwszy.Number == nil || *pierwszy.Number != "1" {
		t.Errorf("przypis pierwszy niesie numer %v, oczekiwano 1", pierwszy.Number)
	}
	if trzeci.Number == nil || *trzeci.Number != "2" {
		t.Errorf("po usunięciu przypisu środkowego trzeci niesie numer %v, oczekiwano 2 — "+
			"numeracja ma nie mieć dziur", trzeci.Number)
	}
}

// TestSpisTabelNieswiezyPoZmianiePodpisuTabeli mierzy powiązanie dwóch obszarów:
// zmiana podpisu tabeli unieważnia spis tabel, a znacznik przeżywa zapis.
func TestSpisTabelNieswiezyPoZmianiePodpisuTabeli(t *testing.T) {
	adapter, zycie, dokument := aparatUprzazSprawdzianu(t)

	tabela, err := adapter.WstawTabele(zycie, shared.StudioTableInsertRequest{
		DocumentId: dokument, Offset: 0, Rows: 2, Columns: 2,
	})
	if err != nil {
		t.Fatalf("wstawienie tabeli odmówiło: %v", err)
	}
	spis, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindTableIndex, Offset: wskaznik(0),
	})
	if err != nil {
		t.Fatalf("założenie spisu tabel odmówiło: %v", err)
	}
	if _, err := adapter.OdswiezAparat(zycie, shared.StudioApparatusRefreshRequest{
		DocumentId: dokument,
	}); err != nil {
		t.Fatalf("odświeżenie aparatu odmówiło: %v", err)
	}
	swiezy := aparatElementSprawdzianu(t, adapter, zycie, dokument, spis.Item.Id)
	if swiezy.Stale == nil || *swiezy.Stale {
		t.Fatal("spis tabel po odświeżeniu stoi jako nieświeży")
	}
	if len(swiezy.Entries) != 1 {
		t.Errorf("spis tabel niesie %d pozycji, oczekiwano jednej: %+v", len(swiezy.Entries),
			swiezy.Entries)
	}

	if _, err := adapter.UstawPostacTabeli(zycie, shared.StudioTableFormatSetRequest{
		DocumentId: dokument, TableId: tabela.Table.Id,
		Caption: wskaznik("Wykaz stron postępowania"),
	}); err != nil {
		t.Fatalf("zmiana podpisu tabeli odmówiła: %v", err)
	}
	nieswiezy := aparatElementSprawdzianu(t, adapter, zycie, dokument, spis.Item.Id)
	if nieswiezy.Stale == nil || !*nieswiezy.Stale {
		t.Error("zmiana podpisu tabeli nie zapaliła znacznika nieświeżości spisu tabel — " +
			"spis pokazywałby stan sprzed zmiany bez ostrzeżenia")
	}

	if _, err := adapter.OdswiezAparat(zycie, shared.StudioApparatusRefreshRequest{
		DocumentId: dokument, ItemId: &spis.Item.Id,
	}); err != nil {
		t.Fatalf("odświeżenie spisu tabel odmówiło: %v", err)
	}
	poOdswiezeniu := aparatElementSprawdzianu(t, adapter, zycie, dokument, spis.Item.Id)
	if len(poOdswiezeniu.Entries) != 1 ||
		!strings.Contains(poOdswiezeniu.Entries[0].Text, "Wykaz stron postępowania") {

		t.Errorf("spis tabel po odświeżeniu nie niesie nowego podpisu: %+v",
			poOdswiezeniu.Entries)
	}
}

// TestBibliografiaSkladaSieZPowolan mierzy bibliografię wraz z numeracją
// powołań: numer w treści i pozycja wykazu mają być tym samym źródłem.
func TestBibliografiaSkladaSieZPowolan(t *testing.T) {
	adapter, zycie, dokument := aparatUprzazSprawdzianu(t)

	if _, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindCitation, Offset: wskaznik(30),
		CitationKey: wskaznik("kodeks"), SourceTitle: wskaznik("Kodeks postępowania"),
		SourceAuthor: wskaznik("Ustawodawca"), SourceYear: wskaznik("1964"),
	}); err != nil {
		t.Fatalf("założenie powołania odmówiło: %v", err)
	}
	powolanieWczesniejsze, err := adapter.WstawElementAparatu(zycie,
		shared.StudioApparatusInsertRequest{
			DocumentId: dokument, Kind: shared.StudioApparatusKindCitation,
			Offset:      wskaznik(5),
			CitationKey: wskaznik("komentarz"),
			SourceTitle: wskaznik("Komentarz do kodeksu"),
		})
	if err != nil {
		t.Fatalf("założenie powołania wcześniejszego odmówiło: %v", err)
	}
	wykazZrodel, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindBibliography,
		Offset: wskaznik(0),
	})
	if err != nil {
		t.Fatalf("założenie bibliografii odmówiło: %v", err)
	}
	if _, err := adapter.OdswiezAparat(zycie, shared.StudioApparatusRefreshRequest{
		DocumentId: dokument,
	}); err != nil {
		t.Fatalf("odświeżenie aparatu odmówiło: %v", err)
	}

	bibliografia := aparatElementSprawdzianu(t, adapter, zycie, dokument, wykazZrodel.Item.Id)
	if len(bibliografia.Entries) != 2 {
		t.Fatalf("bibliografia niesie %d pozycji, oczekiwano dwóch: %+v",
			len(bibliografia.Entries), bibliografia.Entries)
	}
	if !strings.Contains(bibliografia.Entries[0].Text, "Komentarz do kodeksu") {
		t.Errorf("pozycja pierwsza bibliografii to %q, a pierwsze powołanie w treści "+
			"dotyczy komentarza", bibliografia.Entries[0].Text)
	}
	powolanie := aparatElementSprawdzianu(t, adapter, zycie, dokument,
		powolanieWczesniejsze.Item.Id)
	if powolanie.Number == nil || *powolanie.Number != "1" {
		t.Errorf("powołanie wcześniejsze niesie numer %v, oczekiwano 1 — numer w treści "+
			"i pozycja wykazu mają wskazywać to samo źródło", powolanie.Number)
	}
}

// TestOdwolanieDoNieistniejacegoCeluOdmawia pilnuje, żeby odsyłacz w nikąd nie
// wchodził do pisma.
func TestOdwolanieDoNieistniejacegoCeluOdmawia(t *testing.T) {
	adapter, zycie, dokument := aparatUprzazSprawdzianu(t)

	_, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindCrossReference,
		Offset: wskaznik(0), TargetId: wskaznik("studio-apa-nie-ma-takiego"),
	})
	if err == nil {
		t.Fatal("odwołanie do elementu, którego dokument nie ma, wróciło bez odmowy")
	}
	if !strings.Contains(err.Error(), "studio-apa-nie-ma-takiego") {
		t.Errorf("odmowa nie nazywa brakującego celu: %v", err)
	}
}

// TestPrzypisBezBrzmieniaOdmawia pilnuje zakazu odpowiedzi „ok" z pustym
// wynikiem.
func TestPrzypisBezBrzmieniaOdmawia(t *testing.T) {
	adapter, zycie, dokument := aparatUprzazSprawdzianu(t)

	if _, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindFootnote, Offset: wskaznik(0),
	}); err == nil {
		t.Error("przypis bez brzmienia wróciło bez odmowy")
	}
	if _, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindBookmark, Offset: wskaznik(0),
	}); err == nil {
		t.Error("zakładka bez nazwy wróciła bez odmowy")
	}
}

// TestOdswiezenieAparatuBezElementowOdmawia pilnuje, że odświeżenie niczego nie
// wraca jako wykonane.
func TestOdswiezenieAparatuBezElementowOdmawia(t *testing.T) {
	adapter, zycie, dokument := aparatUprzazSprawdzianu(t)

	if _, err := adapter.OdswiezAparat(zycie, shared.StudioApparatusRefreshRequest{
		DocumentId: dokument,
	}); err == nil {
		t.Error("odświeżenie aparatu dokumentu bez ani jednego elementu wróciło bez odmowy")
	}
}

// TestZakladkaNieUdajeOdswiezenia pilnuje uczciwości bilansu: element, którego
// nie ma z czego przeliczyć, jest wymieniony w pominięciach.
func TestZakladkaNieUdajeOdswiezenia(t *testing.T) {
	adapter, zycie, dokument := aparatUprzazSprawdzianu(t)

	if _, err := adapter.WstawElementAparatu(zycie, shared.StudioApparatusInsertRequest{
		DocumentId: dokument, Kind: shared.StudioApparatusKindBookmark, Offset: wskaznik(3),
		Label: wskaznik("Do sprawdzenia"),
	}); err != nil {
		t.Fatalf("założenie zakładki odmówiło: %v", err)
	}
	odswiezony, err := adapter.OdswiezAparat(zycie, shared.StudioApparatusRefreshRequest{
		DocumentId: dokument,
	})
	if err != nil {
		t.Fatalf("odświeżenie aparatu odmówiło: %v", err)
	}
	if odswiezony.Balance.SkippedCount == 0 {
		t.Error("odświeżenie zakładki nie oddało pominięcia — zakładka nie liczy się " +
			"z treści i rdzeń nie ma prawa udawać, że ją przeliczył")
	}
}
