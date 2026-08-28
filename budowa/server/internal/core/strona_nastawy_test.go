package core

import (
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Sprawdziany obszaru strony i sekcji mierzą skutek na postaci dokumentu.

// TestStronaWykazNosnikowNiesieKoperty mierzy uzupełnienie wykazu rdzenia
// o koperty DL, C4, C5 i C6, każda oznaczona rodzajem i bez powtórzeń nazwy.
func TestStronaWykazNosnikowNiesieKoperty(t *testing.T) {
	wykaz := stronaNosniki()
	wymagane := map[string]bool{
		"A3": false, "A4": false, "A5": false, "Letter": false, "Legal": false,
		"Tabloid": false, "DL": false, "C4": false, "C5": false, "C6": false,
	}
	for _, nosnik := range wykaz {
		if _, jest := wymagane[nosnik.Name]; jest {
			wymagane[nosnik.Name] = true
		}
	}
	for nazwa, jest := range wymagane {
		if !jest {
			t.Errorf("wykaz nośników nie niesie formatu %s", nazwa)
		}
	}

	// Koperty muszą być oznaczone rodzajem, bo od tego zależy, czy nadruk
	// koperty w ogóle się wykona.
	for _, nazwa := range []string{"DL", "C4", "C5", "C6"} {
		nosnik, jest := stronaNosnik(nazwa)
		if !jest {
			t.Fatalf("nośnika %s nie ma w wykazie", nazwa)
		}
		if nosnik.Kind != shared.StudioPaperKindEnvelope {
			t.Errorf("nośnik %s ma rodzaj %q, a jest kopertą", nazwa, nosnik.Kind)
		}
	}

	// Wymiary koperty C5 są normą ISO 269, którą sprawdzian ma prawo pilnować.
	c5, _ := stronaNosnik("c5")
	if c5.WidthMm != 162 || c5.HeightMm != 229 {
		t.Errorf("koperta C5 ma wymiary %v na %v mm, a normą jest 162 na 229",
			c5.WidthMm, c5.HeightMm)
	}
	// Żadna nazwa nie powtarza się dwa razy w wykazie nośników.
	widziane := map[string]int{}
	for _, nosnik := range wykaz {
		widziane[strings.ToUpper(nosnik.Name)]++
	}
	for nazwa, ile := range widziane {
		if ile > 1 {
			t.Errorf("nośnik %s stoi w wykazie %d razy", nazwa, ile)
		}
	}
}

// TestStronaOrientacjaZamieniaWymiary pilnuje, żeby zmiana orientacji na
// poziomą naprawdę zamieniała szerokość z wysokością strony.
func TestStronaOrientacjaZamieniaWymiary(t *testing.T) {
	pionowa := shared.StudioPageSetup{
		PageSize:    postacWskaznikTekstu("A4"),
		Orientation: postacWskaznikOrientacjiStrony(shared.StudioPageOrientationPionowa),
	}
	szerokosc, wysokosc := stronaWymiary(&pionowa)
	if szerokosc != 210 || wysokosc != 297 {
		t.Fatalf("A4 pionowa ma wymiary %v na %v mm, a ma mieć 210 na 297", szerokosc, wysokosc)
	}
	pozioma := pionowa
	pozioma.Orientation = postacWskaznikOrientacjiStrony(shared.StudioPageOrientationPozioma)
	szerokosc, wysokosc = stronaWymiary(&pozioma)
	if szerokosc != 297 || wysokosc != 210 {
		t.Errorf("A4 pozioma ma wymiary %v na %v mm, a ma mieć 297 na 210", szerokosc, wysokosc)
	}
}

// TestStronaScalenieNastawNieZdejmujePozostalych mierzy zasadę scalania: pole
// podane wchodzi, pole niepodane zostaje.
func TestStronaScalenieNastawNieZdejmujePozostalych(t *testing.T) {
	zastane := shared.StudioPageSetup{
		PageSize:    postacWskaznikTekstu("A4"),
		Orientation: postacWskaznikOrientacjiStrony(shared.StudioPageOrientationPionowa),
		Columns:     postacWskaznikLiczby(2),
		GutterMm:    postacWskaznikMiary(10),
		MarginTop:   postacWskaznikLiczby(25),
	}
	nastawy, zmian, err := stronaScalNastawy(&zastane, shared.StudioPageSetupSetRequest{
		PaperName:   postacWskaznikTekstu("A3"),
		Orientation: postacWskaznikOrientacjiStrony(shared.StudioPageOrientationPozioma),
	})
	if err != nil {
		t.Fatalf("scalenie nastaw odmówiło: %v", err)
	}
	if zmian != 2 {
		t.Errorf("scalenie policzyło %d zmian, a podano nośnik i orientację", zmian)
	}
	if nastawy.Columns == nil || *nastawy.Columns != 2 {
		t.Error("zmiana nośnika zdjęła kolumny")
	}
	if nastawy.GutterMm == nil || *nastawy.GutterMm != 10 {
		t.Error("zmiana nośnika zdjęła margines na oprawę")
	}
	if nastawy.MarginTop == nil || *nastawy.MarginTop != 25 {
		t.Error("zmiana nośnika zdjęła margines górny")
	}
	if nastawy.PageSize == nil || *nastawy.PageSize != "A3" {
		t.Error("nośnik nie przestawił się na A3")
	}

	// Nośnik nieznany jest odmową nazwaną wraz z wykazem — nie ciszą.
	if _, _, err := stronaScalNastawy(&zastane, shared.StudioPageSetupSetRequest{
		PaperName: postacWskaznikTekstu("A9"),
	}); err == nil {
		t.Error("nośnik spoza wykazu przeszedł bez odmowy")
	} else if !strings.Contains(err.Error(), "A4") {
		t.Errorf("odmowa nie nazywa wykazu nośników: %v", err)
	}

	// Wymiar własny przebija nazwę i przestawia rodzaj nośnika na własny.
	wlasny, _, err := stronaScalNastawy(nil, shared.StudioPageSetupSetRequest{
		WidthMm:  postacWskaznikMiary(500),
		HeightMm: postacWskaznikMiary(700),
	})
	if err != nil {
		t.Fatalf("nośnik własny odmówił: %v", err)
	}
	if wlasny.PaperKind == nil || *wlasny.PaperKind != shared.StudioPaperKindCustom {
		t.Error("nośnik podany wymiarami nie dostał rodzaju własnego")
	}
	if szerokosc, _ := stronaWymiary(&wlasny); szerokosc != 500 {
		t.Errorf("nośnik własny ma szerokość %v, a podano 500", szerokosc)
	}
}

// TestStronaSzerokoscUzytkowaOdejmujeMarginesyIKolumny mierzy rachunek, na
// którym stoi bilans przeliczenia układu.
func TestStronaSzerokoscUzytkowaOdejmujeMarginesyIKolumny(t *testing.T) {
	nastawy := shared.StudioPageSetup{
		PageSize:    postacWskaznikTekstu("A4"),
		MarginLeft:  postacWskaznikLiczby(25),
		MarginRight: postacWskaznikLiczby(25),
	}
	if szerokosc := stronaSzerokoscUzytkowa(&nastawy); szerokosc != 160 {
		t.Errorf("kolumna tekstu A4 z marginesami 25 mm ma %v mm, a ma mieć 160", szerokosc)
	}
	nastawy.GutterMm = postacWskaznikMiary(10)
	if szerokosc := stronaSzerokoscUzytkowa(&nastawy); szerokosc != 150 {
		t.Errorf("margines na oprawę nie zszedł z kolumny tekstu: %v mm", szerokosc)
	}
	nastawy.Columns = postacWskaznikLiczby(2)
	nastawy.ColumnGapMm = postacWskaznikMiary(10)
	if szerokosc := stronaSzerokoscUzytkowa(&nastawy); szerokosc != 70 {
		t.Errorf("dwie kolumny z odstępem 10 mm dają %v mm, a mają dawać 70", szerokosc)
	}
}

// TestStronaBilansNazywaTabeleSzerszaNizNosnik pilnuje zakazu ciszy: format
// uboższy niż dokument jest normalną sytuacją, przemilczenie straty nie jest.
func TestStronaBilansNazywaTabeleSzerszaNizNosnik(t *testing.T) {
	forma := shared.StudioDocumentForm{
		Tables: []shared.StudioDocumentTable{{
			Id: "studio-tab-szeroka", Rows: 2, Columns: 3,
			ColumnWidthsMm: []float64{100, 100, 100},
		}},
		Objects: []shared.StudioDocumentObject{{
			Id: "studio-obi-obraz", WidthMm: postacWskaznikMiary(400),
		}},
	}
	nastawy := shared.StudioPageSetup{
		PageSize:    postacWskaznikTekstu("A4"),
		MarginLeft:  postacWskaznikLiczby(25),
		MarginRight: postacWskaznikLiczby(25),
	}
	pominiete := stronaBilansUkladu(&forma, &nastawy)
	if len(pominiete) != 2 {
		t.Fatalf("bilans nazwał %d pozycji, a nie mieszczą się tabela i obraz: %+v",
			len(pominiete), pominiete)
	}
	if pominiete[0].Detail == nil ||
		!strings.Contains(*pominiete[0].Detail, "studio-tab-szeroka") {
		t.Error("bilans nie nazywa tabeli, która się nie zmieściła")
	}

	// Na nośniku, który treść mieści, bilans milczy — i to też jest prawdą,
	// której sprawdzian ma pilnować.
	szeroki := shared.StudioPageSetup{
		PageSize:    postacWskaznikTekstu("A2"),
		Orientation: postacWskaznikOrientacjiStrony(shared.StudioPageOrientationPozioma),
		MarginLeft:  postacWskaznikLiczby(10),
		MarginRight: postacWskaznikLiczby(10),
	}
	if pominiete := stronaBilansUkladu(&forma, &szeroki); len(pominiete) != 0 {
		t.Errorf("bilans zgłasza pominięcia na nośniku, który wszystko mieści: %+v", pominiete)
	}

	// Marginesy szersze niż nośnik są nazwanym brakiem, nie zerową kolumną.
	ciasny := shared.StudioPageSetup{
		PageSize:    postacWskaznikTekstu("A6"),
		MarginLeft:  postacWskaznikLiczby(80),
		MarginRight: postacWskaznikLiczby(80),
	}
	if pominiete := stronaBilansUkladu(&forma, &ciasny); len(pominiete) == 0 {
		t.Error("marginesy szersze niż nośnik przeszły bez nazwania")
	}
}

// TestStronaNaglowkiSaOsobneWedleZasiegu jest sednem tego odcinka: nagłówek
// pierwszej strony, stron zwykłych i stron parzystych to TRZY osobne nagłówki
// jednej sekcji.
func TestStronaNaglowkiSaOsobneWedleZasiegu(t *testing.T) {
	wykaz, zmian := stronaScalNaglowek(nil, shared.StudioPageHeaderfooterSetRequest{
		HeaderText: postacWskaznikTekstu("Nagłówek stron zwykłych"),
		FooterText: postacWskaznikTekstu("Stopka stron zwykłych"),
	}, shared.StudioHeaderScopeDefault)
	if zmian != 2 || len(wykaz) != 1 {
		t.Fatalf("pierwszy nagłówek: %d zmian, %d pozycji", zmian, len(wykaz))
	}

	wykaz, zmian = stronaScalNaglowek(wykaz, shared.StudioPageHeaderfooterSetRequest{
		HeaderText: postacWskaznikTekstu("Bez nagłówka na pierwszej"),
	}, shared.StudioHeaderScopeFirstPage)
	if zmian != 1 {
		t.Fatalf("nagłówek pierwszej strony: %d zmian", zmian)
	}
	if len(wykaz) != 2 {
		t.Fatalf("nagłówek pierwszej strony nadpisał nagłówek stron zwykłych: %+v", wykaz)
	}
	zwykle := stronaNaglowekZasiegu(wykaz, shared.StudioHeaderScopeDefault)
	if zwykle == nil || zwykle.HeaderText == nil ||
		*zwykle.HeaderText != "Nagłówek stron zwykłych" {
		t.Error("nagłówek stron zwykłych zginął po ustawieniu pierwszej strony")
	}
	if zwykle.FooterText == nil || *zwykle.FooterText != "Stopka stron zwykłych" {
		t.Error("stopka stron zwykłych zginęła po ustawieniu nagłówka pierwszej strony")
	}

	// Napis pusty zdejmuje nagłówek — sprawdzian bierze adres zmiennej wprost.
	pusty := ""
	wykaz, zmian = stronaScalNaglowek(wykaz, shared.StudioPageHeaderfooterSetRequest{
		HeaderText: &pusty,
	}, shared.StudioHeaderScopeFirstPage)
	if zmian != 1 {
		t.Fatalf("zdjęcie nagłówka: %d zmian", zmian)
	}
	pierwsza := stronaNaglowekZasiegu(wykaz, shared.StudioHeaderScopeFirstPage)
	if pierwsza == nil || pierwsza.HeaderText != nil {
		t.Error("napis pusty nie zdjął nagłówka pierwszej strony")
	}

	// Zasięg spoza kontraktu jest odmową nazwaną.
	if _, err := stronaZasiegNaglowka("ostatniaStrona"); err == nil {
		t.Error("zasięg spoza kontraktu przeszedł bez odmowy")
	}
	// Zasięg pusty znaczy strony zwykłe — to jest nastawa, nie brak.
	if zasieg, err := stronaZasiegNaglowka(""); err != nil ||
		zasieg != shared.StudioHeaderScopeDefault {
		t.Errorf("zasięg pusty dał %q, błąd %v", zasieg, err)
	}
}

// stronaNaglowekZasiegu jest pomocnikiem sprawdzianu — wyszukuje nagłówek
// wskazanego zasięgu w wykazie i oddaje wskaźnik na niego albo nil.
func stronaNaglowekZasiegu(wykaz []shared.StudioHeaderFooter,
	zasieg shared.StudioHeaderScope) *shared.StudioHeaderFooter {

	for i := range wykaz {
		if wykaz[i].Scope == zasieg {
			return &wykaz[i]
		}
	}
	return nil
}

// TestStronaNaglowekZNastawStarychNieGinie mierzy przejście: dokument założony
// przed tą dobudową ma nagłówek w jednym polu nastaw strony i nie wolno mu go
// zgubić.
func TestStronaNaglowekZNastawStarychNieGinie(t *testing.T) {
	nastawy := shared.StudioPageSetup{
		Header: postacWskaznikTekstu("Pismo urzędowe"),
		Footer: postacWskaznikTekstu("strona"),
	}
	wykaz := stronaNaglowkiZNastaw(&nastawy)
	if len(wykaz) != 1 {
		t.Fatalf("nagłówek ze starych nastaw dał %d pozycji", len(wykaz))
	}
	if wykaz[0].Scope != shared.StudioHeaderScopeDefault {
		t.Errorf("nagłówek ze starych nastaw dostał zasięg %q", wykaz[0].Scope)
	}
	if wykaz[0].HeaderText == nil || *wykaz[0].HeaderText != "Pismo urzędowe" {
		t.Error("treść nagłówka ze starych nastaw nie przeszła")
	}
	if stronaNaglowkiZNastaw(&shared.StudioPageSetup{}) != nil {
		t.Error("nastawy bez nagłówka dały nagłówek widmo")
	}
}

// TestStronaZnakWodnyWymagaTegoCzymMaStanac pilnuje, żeby znak wodny nie był
// zapisem, którego nie da się narysować.
func TestStronaZnakWodnyWymagaTegoCzymMaStanac(t *testing.T) {
	if _, err := stronaScalZnakWodny(nil, shared.StudioPageWatermarkSetRequest{
		Kind: shared.StudioWatermarkKindText,
	}); err == nil {
		t.Error("znak wodny napisowy bez napisu przeszedł bez odmowy")
	}
	if _, err := stronaScalZnakWodny(nil, shared.StudioPageWatermarkSetRequest{
		Kind: shared.StudioWatermarkKindImage,
	}); err == nil {
		t.Error("znak wodny obrazowy bez zasobu przeszedł bez odmowy")
	}
	znak, err := stronaScalZnakWodny(nil, shared.StudioPageWatermarkSetRequest{
		Kind: shared.StudioWatermarkKindText,
		Text: postacWskaznikTekstu("KOPIA"),
	})
	if err != nil {
		t.Fatalf("znak wodny napisowy odmówił: %v", err)
	}
	if znak.Opacity == nil || znak.FontSizePt == nil {
		t.Error("znak wodny nie dostał krycia ani stopnia pisma, więc nie ma czym stanąć")
	}
	if _, err := stronaScalZnakWodny(nil, shared.StudioPageWatermarkSetRequest{
		Kind:    shared.StudioWatermarkKindText,
		Text:    postacWskaznikTekstu("KOPIA"),
		Opacity: postacWskaznikMiary(4),
	}); err == nil {
		t.Error("krycie poza zakresem od zera do jednego przeszło bez odmowy")
	}
}

// TestStronaNumeracjaSprawdzaFormat pilnuje, żeby numeracja nie przyjęła
// formatu, którego kontrakt nie zna.
func TestStronaNumeracjaSprawdzaFormat(t *testing.T) {
	numeracja, err := stronaScalNumeracje(nil, shared.StudioPageNumberingSetRequest{
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("numeracja odmówiła: %v", err)
	}
	if numeracja.Format == nil || *numeracja.Format != shared.StudioPageNumberFormatArabic {
		t.Error("numeracja bez formatu nie dostała cyfr arabskich")
	}
	zmyslony := shared.StudioPageNumberFormat("gwiazdki")
	if _, err := stronaScalNumeracje(nil, shared.StudioPageNumberingSetRequest{
		Enabled: true, Format: &zmyslony,
	}); err == nil {
		t.Error("format numeru spoza kontraktu przeszedł bez odmowy")
	}
}

// TestStronaPodzialRozdzielaAkapit mierzy skutek na drzewie postaci: podział
// w środku akapitu ma go rozdzielić, a na granicy — nie ruszać.
func TestStronaPodzialRozdzielaAkapit(t *testing.T) {
	forma := shared.StudioDocumentForm{
		Blocks: []shared.StudioDocumentBlock{
			{Id: "studio-blok-jeden", Kind: blokPostaciAkapit,
				Runs: []shared.StudioDocumentRun{{Text: "Zażółć gęślą jaźń"}}},
			{Id: "studio-blok-dwa", Kind: blokPostaciAkapit,
				Runs: []shared.StudioDocumentRun{{Text: "Drugi akapit"}}},
		},
	}
	postacPrzeliczZakresy(&forma)
	trescPrzed := postacTekstFormy(&forma)

	// Podział po szóstym znaku, dziewiątym bajcie — „ż" i „ó" są dwubajtowe.
	wskazanie, roznica := stronaRozdzielAkapit(&forma, 6)
	if roznica != 1 {
		t.Errorf("rozdzielenie akapitu dało różnicę %d, a dokłada jeden znak podziału", roznica)
	}
	if len(forma.Blocks) != 3 {
		t.Fatalf("rozdzielenie dało %d bloków, a ma dać trzy", len(forma.Blocks))
	}
	if postacTekstBloku(forma.Blocks[0]) != "Zażółć" {
		t.Errorf("pierwsza połowa akapitu brzmi %q, a ma brzmieć „Zażółć”",
			postacTekstBloku(forma.Blocks[0]))
	}
	if postacTekstBloku(forma.Blocks[1]) != " gęślą jaźń" {
		t.Errorf("druga połowa akapitu brzmi %q", postacTekstBloku(forma.Blocks[1]))
	}
	if wskazanie != 1 {
		t.Errorf("podział ma stanąć przed blokiem 1, a wskazano %d", wskazanie)
	}
	// Litery nie zginęły ani nie doszły — poza jednym znakiem podziału wiersza.
	if len([]rune(postacTekstFormy(&forma))) != len([]rune(trescPrzed))+1 {
		t.Errorf("rozdzielenie zmieniło treść: %q wobec %q",
			postacTekstFormy(&forma), trescPrzed)
	}

	// Podział na granicy akapitu niczego nie rozdziela.
	ile := len(forma.Blocks)
	if _, roznica := stronaRozdzielAkapit(&forma, 0); roznica != 0 {
		t.Error("podział na początku akapitu rozdzielił akapit")
	}
	if len(forma.Blocks) != ile {
		t.Error("podział na granicy dołożył blok")
	}
}

// TestStronaTabulatorNieMnozySieNaTymSamymPolozeniu mierzy zachowanie linijki:
// powtórne kliknięcie w to samo miejsce przestawia tabulator, a nie stawia
// drugiego obok.
func TestStronaTabulatorNieMnozySieNaTymSamymPolozeniu(t *testing.T) {
	wykaz, zmieniono := stronaTabulatory(nil, 40, shared.StudioTabKindLeft, nil, false)
	if !zmieniono || len(wykaz) != 1 {
		t.Fatalf("pierwszy tabulator: zmieniono %v, pozycji %d", zmieniono, len(wykaz))
	}
	kropki := shared.StudioTabLeader(shared.StudioTabLeaderDot)
	wykaz, _ = stronaTabulatory(wykaz, 40.05, shared.StudioTabKindRight, &kropki, false)
	if len(wykaz) != 1 {
		t.Fatalf("tabulator na tym samym położeniu postawił drugi: %+v", wykaz)
	}
	if wykaz[0].Kind != shared.StudioTabKindRight {
		t.Errorf("tabulator nie przestawił rodzaju: %q", wykaz[0].Kind)
	}
	if wykaz[0].Leader == nil || *wykaz[0].Leader != shared.StudioTabLeaderDot {
		t.Error("tabulator nie przyjął znaku wiodącego")
	}

	// Tabulatory stoją po położeniu rosnąco, w kolejności na kartce.
	wykaz, _ = stronaTabulatory(wykaz, 20, shared.StudioTabKindDecimal, nil, false)
	if len(wykaz) != 2 || wykaz[0].PositionMm != 20 {
		t.Errorf("tabulatory nie są uporządkowane po położeniu: %+v", wykaz)
	}

	// Zdjęcie tabulatora, którego nie ma, oddaje fałsz, nie cichą zgodę.
	if _, zmieniono := stronaTabulatory(wykaz, 99, shared.StudioTabKindLeft, nil, true); zmieniono {
		t.Error("zdjęcie tabulatora nieistniejącego wróciło jako wykonane")
	}
	wykaz, zmieniono = stronaTabulatory(wykaz, 20, shared.StudioTabKindLeft, nil, true)
	if !zmieniono || len(wykaz) != 1 {
		t.Errorf("zdjęcie tabulatora: zmieniono %v, pozostało %d", zmieniono, len(wykaz))
	}
}

// TestStronaNastawySekcjiNakladajaSieNaDokument mierzy dziedziczenie: sekcja
// nadpisuje wyłącznie to, o czym mówi.
func TestStronaNastawySekcjiNakladajaSieNaDokument(t *testing.T) {
	dokumentu := shared.StudioPageSetup{
		PageSize:    postacWskaznikTekstu("A4"),
		Orientation: postacWskaznikOrientacjiStrony(shared.StudioPageOrientationPionowa),
		MarginTop:   postacWskaznikLiczby(25),
		Columns:     postacWskaznikLiczby(1),
	}
	sekcji := shared.StudioPageSetup{
		PageSize:    postacWskaznikTekstu("A3"),
		Orientation: postacWskaznikOrientacjiStrony(shared.StudioPageOrientationPozioma),
	}
	skuteczne := stronaNastawySkuteczne(&dokumentu, &sekcji)
	if skuteczne.PageSize == nil || *skuteczne.PageSize != "A3" {
		t.Error("nośnik sekcji nie przebił nośnika dokumentu")
	}
	if skuteczne.MarginTop == nil || *skuteczne.MarginTop != 25 {
		t.Error("sekcja bez marginesu zgubiła margines dokumentu")
	}
	szerokosc, wysokosc := stronaWymiary(&skuteczne)
	if szerokosc != 420 || wysokosc != 297 {
		t.Errorf("sekcja pozioma A3 ma wymiary %v na %v mm, a ma mieć 420 na 297",
			szerokosc, wysokosc)
	}
	// Sekcja bez własnych nastaw nie ma prawa niczego zmienić.
	if skuteczne := stronaNastawySkuteczne(&dokumentu, nil); skuteczne.PageSize == nil ||
		*skuteczne.PageSize != "A4" {
		t.Error("sekcja bez nastaw zmieniła nośnik dokumentu")
	}
}

// postacWskaznikOrientacjiStrony oddaje wskaźnik na orientację — pola
// nieobowiązkowe kontraktu są wskaźnikami, a adresu stałej wziąć nie można.
func postacWskaznikOrientacjiStrony(wartosc shared.StudioPageOrientation) *shared.StudioPageOrientation {
	kopia := wartosc
	return &kopia
}
