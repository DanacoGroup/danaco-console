package core

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"math"
	"strings"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font/sfnt"

	"danacoconsole/shared"
)

// Skutek warsztatów WEKTORA, IKON, MAKIETY i BARWY modułu Design.
//
// ── Po co ten plik istnieje ─────────────────────────────────────────────────
// Sprawdziany tych czterech obszarów były przez jedną turę sprawdzianami
// KOMPILACJI: wołały komendę i patrzyły, czy odpowiedź jest odpowiedzią.
// Sprawdzian tego rodzaju przechodzi także wtedy, gdy operacja logiczna oddaje
// pierwszy kształt zamiast sumy, krój niesie zero glifów, a schemat wychodzi
// pustym płótnem — bo wszystko to są odpowiedzi udane.
//
// Dlatego każdy sprawdzian tego pliku pyta o LICZBĘ zmierzoną w wyniku:
//   - operacja logiczna — o współrzędne węzłów i o to, które punkty płaszczyzny
//     do kształtu należą (sprawdzenie przynależności, nie oglądanie prostokąta
//     otaczającego);
//   - krój ikonowy — o liczbę glifów ODCZYTANĄ z pliku TTF przez czytnik krojów
//     wraz z liczbą krzywych w konturach;
//   - schemat — o liczbę ścieżek w dokumencie SVG;
//   - barwa — o współczynnik kontrastu przeliczony niezależnie wzorem WCAG
//     i o kąty odcieni harmonii;
//   - układ makiety — o położenia warstw ODCZYTANE z bazy po zapisie;
//   - odszumienie — o to, czy krawędź została krawędzią, a płaski obszar
//     wygładzeniem.
//
// Pierwszy przebieg tego pliku wykrył defekt zastany: tabela `maxp` składanego
// kroju miała 36 bajtów wobec 32 wymaganych przez format, więc żaden czytnik
// krojów nie wczytywał pliku, choć rdzeń oddawał go bez odmowy. To jest miara
// wartości sprawdzianu skutku — sprawdzian kompilacji świecił nad tym zielono.

// ── Warsztat wektorowy: operacje logiczne na węzłach ────────────────────────

// zalozPlanszeWektorowaSprawdzianu zakłada pustą kompozycję, na której stają
// kształty warsztatu wektorowego.
func zalozPlanszeWektorowaSprawdzianu(t *testing.T, zmontowany *Zmontowany,
	zycie context.Context, okno string) string {

	t.Helper()

	var wynik shared.DesignBoardUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardUpdate,
		shared.DesignBoardUpdateRequest{
			WindowId: okno,
			Name:     wskaznik("plansza sprawdzianu wektora"),
		}, &wynik)
	return wynik.Board.Id
}

// dolozProstokatSprawdzianu wstawia prostokąt jako ścieżkę o węzłach i oddaje
// jej identyfikator.
func dolozProstokatSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	plansza string, x, y, szerokosc, wysokosc float64) string {

	t.Helper()

	var wynik shared.DesignVectorShapeAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignVectorShapeAdd,
		shared.DesignVectorShapeAddRequest{
			BoardId: plansza, Kind: shared.DesignShapeKindRectangle,
			X: x, Y: y, Width: szerokosc, Height: wysokosc,
		}, &wynik)
	return wynik.Path.Id
}

// czyPunktWSciezceDesignu rozstrzyga, czy punkt leży w wielokącie węzłów —
// regułą parzystości przecięć.
//
// To jest sedno pomiaru operacji logicznej: prostokąt otaczający wyniku bywa
// identyczny dla sumy i dla pierwszego z kształtów, a przynależność punktu
// rozróżnia je bez pudła.
func czyPunktWSciezceDesignu(wezly []shared.DesignVectorNode, x, y float64) bool {
	wewnatrz := false
	ile := len(wezly)
	for numer := 0; numer < ile; numer++ {
		biezacy := wezly[numer]
		poprzedni := wezly[(numer-1+ile)%ile]
		if (biezacy.Y > y) == (poprzedni.Y > y) {
			continue
		}
		przeciecie := (poprzedni.X-biezacy.X)*(y-biezacy.Y)/(poprzedni.Y-biezacy.Y) + biezacy.X
		if x < przeciecie {
			wewnatrz = !wewnatrz
		}
	}
	return wewnatrz
}

// prostokatOtaczajacyWezlowDesignu oddaje krańce wykazu węzłów.
func prostokatOtaczajacyWezlowDesignu(wezly []shared.DesignVectorNode) (float64, float64,
	float64, float64) {

	najmniejszyX, najwiekszyX := math.MaxFloat64, -math.MaxFloat64
	najmniejszyY, najwiekszyY := math.MaxFloat64, -math.MaxFloat64
	for _, wezel := range wezly {
		najmniejszyX = math.Min(najmniejszyX, wezel.X)
		najwiekszyX = math.Max(najwiekszyX, wezel.X)
		najmniejszyY = math.Min(najmniejszyY, wezel.Y)
		najwiekszyY = math.Max(najwiekszyY, wezel.Y)
	}
	return najmniejszyX, najmniejszyY, najwiekszyX, najwiekszyY
}

// TestOperacjeLogiczneNaWezlachDajaZmierzonyKsztalt mierzy WYNIK sumy, części
// wspólnej i różnicy dwóch prostokątów: prostokąt otaczający oraz przynależność
// trzech punktów rozstrzygających. Operacja oddająca pierwszy kształt zamiast
// wyniku przechodzi sprawdzian kompilacji i wywala się tutaj.
func TestOperacjeLogiczneNaWezlachDajaZmierzonyKsztalt(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	plansza := zalozPlanszeWektorowaSprawdzianu(t, zmontowany, zycie, "okno-wektora")

	// Dwa prostokąty 40×40 zachodzące na siebie narożami: część wspólna to
	// kwadrat 20×20 od (20;20) do (40;40).
	przypadki := []struct {
		operacja                    shared.DesignBooleanOp
		odX, odY, doX, doY          float64
		wPierwszym, wCzesciWspolnej bool
		wDrugim                     bool
	}{
		{shared.DesignBooleanOpUnion, 0, 0, 60, 60, true, true, true},
		{shared.DesignBooleanOpIntersect, 20, 20, 40, 40, false, true, false},
		{shared.DesignBooleanOpSubtract, 0, 0, 40, 40, true, false, false},
	}
	for _, przypadek := range przypadki {
		pierwszy := dolozProstokatSprawdzianu(t, zmontowany, zycie, plansza, 0, 0, 40, 40)
		drugi := dolozProstokatSprawdzianu(t, zmontowany, zycie, plansza, 20, 20, 40, 40)

		var wynik shared.DesignVectorBooleanResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandDesignVectorBoolean,
			shared.DesignVectorBooleanRequest{
				PathIds:   []string{pierwszy, drugi},
				Operation: przypadek.operacja,
			}, &wynik)

		// Ścieżki źródłowe mają zniknąć i wrócić w bilansie — inaczej okno
		// pokazywałoby kształty, których w bazie już nie ma.
		if len(wynik.RemovedPathIds) != 2 {
			t.Errorf("operacja %s usunęła %d ścieżek źródłowych, a były dwie",
				przypadek.operacja, len(wynik.RemovedPathIds))
		}

		// Węzły czytamy Z BAZY, nie z odpowiedzi: wynik operacji ma dać się ciągnąć
		// piórem dalej, a to znaczy, że musi w bazie leżeć.
		var wykaz shared.DesignVectorPathListResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandDesignVectorPathList,
			shared.DesignVectorPathListRequest{BoardId: plansza}, &wykaz)
		if len(wykaz.Paths) != 1 {
			t.Fatalf("po operacji %s na planszy leży %d ścieżek, a ma leżeć jedna (wynik)",
				przypadek.operacja, len(wykaz.Paths))
		}
		wezly := wykaz.Paths[0].Nodes
		if len(wezly) < 4 {
			t.Fatalf("wynik operacji %s ma %d węzłów — kształt o mniej niż czterech węzłach nie "+
				"jest wynikiem operacji na dwóch prostokątach", przypadek.operacja, len(wezly))
		}

		odX, odY, doX, doY := prostokatOtaczajacyWezlowDesignu(wezly)
		const tolerancja = 0.6
		if math.Abs(odX-przypadek.odX) > tolerancja || math.Abs(odY-przypadek.odY) > tolerancja ||
			math.Abs(doX-przypadek.doX) > tolerancja || math.Abs(doY-przypadek.doY) > tolerancja {

			t.Errorf("operacja %s dała kształt od (%.2f;%.2f) do (%.2f;%.2f), a ma dać "+
				"od (%.0f;%.0f) do (%.0f;%.0f)", przypadek.operacja, odX, odY, doX, doY,
				przypadek.odX, przypadek.odY, przypadek.doX, przypadek.doY)
		}

		// Trzy punkty rozstrzygające: wyłącznie w pierwszym (10;10), w obu (30;30),
		// wyłącznie w drugim (50;50).
		sprawdz := func(x, y float64, oczekiwane bool, gdzie string) {
			if czyPunktWSciezceDesignu(wezly, x, y) != oczekiwane {
				t.Errorf("operacja %s: punkt (%.0f;%.0f) %s w wyniku, a jest odwrotnie",
					przypadek.operacja, x, y, gdzie)
			}
		}
		sprawdz(10, 10, przypadek.wPierwszym, "leży wyłącznie w pierwszym prostokącie i ma być")
		sprawdz(30, 30, przypadek.wCzesciWspolnej, "leży w obu prostokątach i ma być")
		sprawdz(50, 50, przypadek.wDrugim, "leży wyłącznie w drugim prostokącie i ma być")

		wykonajUdana(t, zmontowany, zycie, shared.CommandDesignVectorPathRemove,
			shared.DesignVectorPathRemoveRequest{PathId: wykaz.Paths[0].Id},
			&shared.DesignVectorPathRemoveResponse{})
	}
}

// ── Krój ikonowy: liczba glifów w pliku TTF ─────────────────────────────────

// TestKrojIkonowyMaZmierzonaLiczbeGlifowWPliku składa krój z ikon katalogu
// i mierzy PLIK czytnikiem krojów: liczbę glifów, przypisanie punktów kodowych,
// obecność konturów i obecność krzywych kwadratowych.
//
// Pole `included` odpowiedzi mogło powstać z długości wykazu żądania — plik
// mówi, ile glifów naprawdę w nim leży.
func TestKrojIkonowyMaZmierzonaLiczbeGlifowWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	var wyszukanie shared.DesignIconLibrarySearchResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignIconLibrarySearch,
		shared.DesignIconLibrarySearchRequest{Limit: wskaznik(6)}, &wyszukanie)
	if len(wyszukanie.Icons) < 6 {
		t.Fatalf("katalog wkompilowany oddał %d ikon, a sprawdzian potrzebuje sześciu",
			len(wyszukanie.Icons))
	}
	kody := []string{}
	for _, ikona := range wyszukanie.Icons {
		kody = append(kody, ikona.Id)
	}

	var pakiet shared.DesignIconSpriteBuildResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignIconSpriteBuild,
		shared.DesignIconSpriteBuildRequest{
			IconIds: kody, Kind: shared.DesignIconSpriteKindWebfont,
			Name: wskaznik("ikony sprawdzianu"), WindowId: wskaznik("okno-ikon"),
		}, &pakiet)
	if pakiet.Included != len(kody) {
		t.Errorf("odpowiedź mówi o %d ikonach w pakiecie, a wskazano %d",
			pakiet.Included, len(kody))
	}
	if pakiet.Asset.Uri == nil {
		t.Fatal("pakiet kroju bez odwołania do treści — nie ma czego zmierzyć")
	}
	bajty := bajtyPodOdwolaniem(t, katalog, *pakiet.Asset.Uri)

	// Czytnik krojów jest tu MIARĄ, nie ozdobą: jeżeli plik ma tabelę o złej
	// długości albo katalog tabel nie zgadza się z treścią, `sfnt.Parse` odmawia
	// tak samo jak przeglądarka.
	kroj, err := sfnt.Parse(bajty)
	if err != nil {
		t.Fatalf("złożony krój nie daje się wczytać czytnikiem krojów: %v — plik w tym stanie "+
			"nie otworzy się ani w przeglądarce, ani w systemie", err)
	}
	// Glifów jest tyle, ile ikon, plus `.notdef` wymagany przez format.
	if kroj.NumGlyphs() != len(kody)+1 {
		t.Errorf("plik kroju niesie %d glifów, a sześć ikon plus glif `.notdef` daje %d",
			kroj.NumGlyphs(), len(kody)+1)
	}

	var bufor sfnt.Buffer
	krzywych, odcinkow := 0, 0
	for numer := range kody {
		znak := rune(poczatekObszaruPrywatnegoDesignu + numer)
		indeks, err := kroj.GlyphIndex(&bufor, znak)
		if err != nil {
			t.Fatalf("odczyt przypisania znaku %U nie powiódł się: %v", znak, err)
		}
		// Zero znaczy `.notdef`: znak, którego krój nie zna. Ikona bez punktu
		// kodowego jest ikoną, do której arkusz stylów nie ma jak sięgnąć.
		if indeks == 0 {
			t.Fatalf("krój nie przypisuje znaku %U do żadnego glifu, a ikona %s ma tam stać",
				znak, kody[numer])
		}
		odcinki, err := kroj.LoadGlyph(&bufor, indeks, jednostekNaFiretDesignu, nil)
		if err != nil {
			t.Fatalf("glif ikony %s nie daje się odczytać: %v", kody[numer], err)
		}
		if len(odcinki) == 0 {
			t.Errorf("glif ikony %s nie ma ani jednego odcinka konturu — w kroju stoi puste "+
				"miejsce, a strona pokaże nic", kody[numer])
		}
		for _, odcinek := range odcinki {
			switch odcinek.Op {
			case sfnt.SegmentOpQuadTo:
				krzywych++
			case sfnt.SegmentOpLineTo:
				odcinkow++
			}
		}
	}
	// Kontury ikon to obrysy kresek z zaokrąglonymi końcami i złączeniami, więc
	// krzywe w nich BYĆ MUSZĄ. Zero krzywych znaczy, że kontur został spłaszczony
	// do łamanej i plik jest kilka razy większy, niż powinien.
	if krzywych == 0 {
		t.Errorf("glify kroju nie niosą ani jednej krzywej kwadratowej (odcinków: %d) — "+
			"kontury zostały spłaszczone do łamanej", odcinkow)
	}
}

// ── Schemat: liczba ścieżek w dokumencie SVG ────────────────────────────────

// TestSchematMaZmierzonaLiczbeSciezekWPliku mierzy liczbę ścieżek w wydanym
// dokumencie SVG oraz bilans węzłów nieumieszczonych. Schemat wydany jako białe
// płótno jest odpowiedzią udaną i plikiem bez treści.
func TestSchematMaZmierzonaLiczbeSciezekWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	wezly := []shared.DesignDiagramNode{
		{Id: "korzen", Label: "Zamówienie"},
		{Id: "sprawdzenie", Label: "Sprawdzenie", ParentId: wskaznik("korzen")},
		{Id: "wydanie", Label: "Wydanie", ParentId: wskaznik("korzen")},
		{Id: "wysylka", Label: "Wysyłka", ParentId: wskaznik("wydanie")},
	}
	var wynik shared.DesignDiagramRenderResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignDiagramRender,
		shared.DesignDiagramRenderRequest{
			Kind: shared.DesignDiagramKindFlow, Nodes: wezly,
			Edges:  []shared.DesignDiagramEdge{{FromId: "korzen", ToId: "sprawdzenie"}},
			Format: wskaznik("svg"), WindowId: wskaznik("okno-schematu"),
		}, &wynik)
	if len(wynik.UnplacedNodeIds) != 0 {
		t.Errorf("schemat o czterech węzłach spójnych zgłasza nieumieszczone: %v",
			wynik.UnplacedNodeIds)
	}
	if wynik.Asset.Uri == nil {
		t.Fatal("schemat bez odwołania do treści")
	}
	dokument := string(bajtyPodOdwolaniem(t, katalog, *wynik.Asset.Uri))
	if !strings.Contains(dokument, "<svg") {
		t.Fatalf("wydany plik nie jest dokumentem SVG; początek: %s",
			dokument[:min(120, len(dokument))])
	}
	sciezek := strings.Count(dokument, "<path")
	// Cztery węzły, jedno połączenie, tło i podpisy konturami — każdy z tych
	// elementów jest ścieżką. Granica pięciu jest ostrożna i nadal wyłapuje
	// płótno puste albo schemat z jednym węzłem.
	if sciezek < 5 {
		t.Errorf("dokument SVG schematu ma %d ścieżek, a cztery węzły z podpisami dają "+
			"co najmniej pięć", sciezek)
	}

	// Węzeł wskazujący nadrzędnego, którego w wykazie nie ma, MA wrócić
	// w bilansie: postawienie go w warstwie zerowej udawałoby, że jest korzeniem.
	var zBrakiem shared.DesignDiagramRenderResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignDiagramRender,
		shared.DesignDiagramRenderRequest{
			Kind: shared.DesignDiagramKindFlow,
			Nodes: append(wezly, shared.DesignDiagramNode{
				Id: "sierota", Label: "Bez rodzica", ParentId: wskaznik("nie-ma-takiego"),
			}),
			Format: wskaznik("svg"), WindowId: wskaznik("okno-schematu"),
		}, &zBrakiem)
	if len(zBrakiem.UnplacedNodeIds) != 1 || zBrakiem.UnplacedNodeIds[0] != "sierota" {
		t.Errorf("węzeł o nieistniejącym nadrzędnym wrócił w bilansie jako %v, a ma wrócić "+
			"jako [sierota]", zBrakiem.UnplacedNodeIds)
	}
}

// TestWykresMaZmierzonaLiczbeSciezekWPliku mierzy to samo co sprawdzian schematu
// i po tym samym płótnie: wykres słupkowy wydany do SVG ma mieć w pliku ścieżkę
// na każdy słupek, osie i podpisy.
//
// Sprawdzian stoi obok schematu, bo obie czynności idą jednym płótnem i jednym
// wydawcą — a wydawca SVG przerywa wykonanie przy styl obrysu bez zakończenia
// kreski i właśnie tak przerywał je przy osiach każdego wykresu.
func TestWykresMaZmierzonaLiczbeSciezekWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	var wynik shared.DesignChartRenderResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignChartRender,
		shared.DesignChartRenderRequest{
			Kind:       shared.DesignChartKindBar,
			Categories: []string{"styczeń", "luty", "marzec"},
			Series: []shared.DesignChartSeries{
				{Name: "sprzedaż", Values: []float64{12, 18, 9}},
				{Name: "zwroty", Values: []float64{2, 3, 1}},
			},
			Title:    wskaznik("Kwartał"),
			Format:   wskaznik("svg"),
			WindowId: wskaznik("okno-wykresu"),
		}, &wynik)
	if wynik.Asset.Uri == nil {
		t.Fatal("wykres bez odwołania do treści")
	}
	dokument := string(bajtyPodOdwolaniem(t, katalog, *wynik.Asset.Uri))
	if !strings.Contains(dokument, "<svg") {
		t.Fatal("wydany plik wykresu nie jest dokumentem SVG")
	}
	// Sześć słupków, osie, tło, tytuł, legenda i podpisy kategorii — ścieżek jest
	// znacznie więcej niż osiem, a granica ośmiu wyłapuje płótno puste.
	if sciezek := strings.Count(dokument, "<path"); sciezek < 8 {
		t.Errorf("dokument SVG wykresu ma %d ścieżek, a dwie serie po trzy wartości z osiami "+
			"i podpisami dają co najmniej osiem", sciezek)
	}
	// Osie są kreską — a kreska bez wskazanego zakończenia przerywała wydawcę SVG.
	if !strings.Contains(dokument, "stroke-linecap") {
		t.Error("dokument SVG wykresu nie niesie zakończenia kreski; styl obrysu bez niego " +
			"przerywa wydawcę SVG i wykres nie powstaje wcale")
	}
}

// ── Makieta: układ automatyczny stawia warstwy w policzalnych miejscach ─────

// TestUkladAutomatycznyStawiaWarstwyWZmierzonychMiejscach liczy położenia
// warstw w głowie i porównuje z tym, co po zapisie leży W BAZIE. Układ, który
// zwraca warstwy przeliczone w pamięci i nie zapisuje ich, przechodzi sprawdzian
// kompilacji i gubi pracę Operatora przy odświeżeniu okna.
func TestUkladAutomatycznyStawiaWarstwyWZmierzonychMiejscach(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var plansza shared.DesignBoardUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardUpdate,
		shared.DesignBoardUpdateRequest{
			WindowId: "okno-makiety",
			Name:     wskaznik("makieta sprawdzianu"),
			Layers: []shared.DesignBoardLayer{
				{Width: wskaznik(50.0), Height: wskaznik(10.0)},
				{Width: wskaznik(50.0), Height: wskaznik(20.0)},
				{Width: wskaznik(50.0), Height: wskaznik(30.0)},
			},
		}, &plansza)
	if len(plansza.Board.Layers) != 3 {
		t.Fatalf("plansza sprawdzianu ma %d warstw, a założono trzy", len(plansza.Board.Layers))
	}
	kody := []string{}
	for _, warstwa := range plansza.Board.Layers {
		kody = append(kody, warstwa.Id)
	}

	var ramka shared.DesignFrameSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignFrameSet,
		shared.DesignFrameSetRequest{
			BoardId: plansza.Board.Id, Name: "ramka sprawdzianu",
			Width: wskaznik(200.0), Height: wskaznik(300.0),
			X: wskaznik(0.0), Y: wskaznik(0.0),
		}, &ramka)

	var uklad shared.DesignLayoutAutoResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignLayoutAuto,
		shared.DesignLayoutAutoRequest{
			FrameId:  ramka.Frame.Id,
			LayerIds: kody,
			Layout: shared.DesignAutoLayout{
				Direction:   shared.DesignLayoutDirectionVertical,
				Gap:         wskaznik(5.0),
				PaddingTop:  wskaznik(8.0),
				PaddingLeft: wskaznik(4.0),
			},
		}, &uklad)

	// Trzy warstwy o wysokościach 10, 20 i 30 z odstępem 5 i odstępem wewnętrznym
	// 8 od góry stają na 8, 23 i 48. Treść ma wtedy 50 szerokości i 70 wysokości.
	oczekiwaneY := map[string]float64{kody[0]: 8, kody[1]: 23, kody[2]: 48}
	if uklad.ContentWidth == nil || math.Abs(*uklad.ContentWidth-50) > 0.001 {
		t.Errorf("szerokość treści po ułożeniu to %v, a trzy warstwy o szerokości 50 dają 50",
			uklad.ContentWidth)
	}
	if uklad.ContentHeight == nil || math.Abs(*uklad.ContentHeight-70) > 0.001 {
		t.Errorf("wysokość treści po ułożeniu to %v, a 10+5+20+5+30 daje 70", uklad.ContentHeight)
	}

	// Pomiar niezależny: stan po zapisie czytamy drugą komendą, nie z odpowiedzi
	// układu.
	var wykaz shared.DesignBoardListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardList,
		shared.DesignBoardListRequest{WindowId: "okno-makiety"}, &wykaz)
	zmierzonych := 0
	for _, tablica := range wykaz.Boards {
		if tablica.Id != plansza.Board.Id {
			continue
		}
		for _, warstwa := range tablica.Layers {
			oczekiwane, jest := oczekiwaneY[warstwa.Id]
			if !jest {
				continue
			}
			zmierzonych++
			if warstwa.Y == nil || math.Abs(*warstwa.Y-oczekiwane) > 0.001 {
				t.Errorf("warstwa %s leży w bazie na wysokości %v, a układ stawia ją na %.0f",
					warstwa.Id, warstwa.Y, oczekiwane)
			}
			if warstwa.X == nil || math.Abs(*warstwa.X-4) > 0.001 {
				t.Errorf("warstwa %s leży w bazie na %v w poziomie, a odstęp wewnętrzny od lewej "+
					"to 4", warstwa.Id, warstwa.X)
			}
		}
	}
	if zmierzonych != 3 {
		t.Errorf("w bazie znalazły się %d z trzech ułożonych warstw", zmierzonych)
	}
}

// ── Barwa: kontrast wzorem WCAG i kąty harmonii ─────────────────────────────

// jasnoscWzglednaWcagSprawdzianu liczy jasność względną barwy wzorem WCAG 2.1 —
// niezależnie od rdzenia, żeby pomiar nie brał liczby z tego samego rachunku,
// który mierzy.
func jasnoscWzglednaWcagSprawdzianu(barwa colorful.Color) float64 {
	skladowa := func(wartosc float64) float64 {
		if wartosc <= 0.03928 {
			return wartosc / 12.92
		}
		return math.Pow((wartosc+0.055)/1.055, 2.4)
	}
	return 0.2126*skladowa(barwa.R) + 0.7152*skladowa(barwa.G) + 0.0722*skladowa(barwa.B)
}

// TestKontrastBarwJestZmierzonyWzoremNiezaleznym sprawdza współczynnik kontrastu
// dwóch par o znanym wyniku oraz zgodność progów AA i AAA z tym współczynnikiem.
func TestKontrastBarwJestZmierzonyWzoremNiezaleznym(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	pary := []struct{ pierwszoplanowa, tlo string }{
		{"#000000", "#ffffff"},
		{"#767676", "#ffffff"},
		{"#1f6feb", "#ffffff"},
	}
	for _, para := range pary {
		var wynik shared.DesignColorContrastCheckResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandDesignColorContrastCheck,
			shared.DesignColorContrastCheckRequest{
				Foreground: para.pierwszoplanowa, Background: para.tlo,
			}, &wynik)

		pierwsza, err := colorful.Hex(para.pierwszoplanowa)
		if err != nil {
			t.Fatalf("barwa sprawdzianu %s nieczytelna: %v", para.pierwszoplanowa, err)
		}
		druga, err := colorful.Hex(para.tlo)
		if err != nil {
			t.Fatalf("barwa sprawdzianu %s nieczytelna: %v", para.tlo, err)
		}
		jasniejsza := math.Max(jasnoscWzglednaWcagSprawdzianu(pierwsza),
			jasnoscWzglednaWcagSprawdzianu(druga))
		ciemniejsza := math.Min(jasnoscWzglednaWcagSprawdzianu(pierwsza),
			jasnoscWzglednaWcagSprawdzianu(druga))
		zmierzony := (jasniejsza + 0.05) / (ciemniejsza + 0.05)

		if math.Abs(wynik.Result.Ratio-zmierzony) > 0.05 {
			t.Errorf("rdzeń podaje kontrast %.3f dla pary %s na %s, a wzór WCAG daje %.3f",
				wynik.Result.Ratio, para.pierwszoplanowa, para.tlo, zmierzony)
		}
		// Próg AA dla tekstu zwykłego to 4,5, AAA to 7 — orzeczenie ma się zgadzać
		// z liczbą, którą sam rdzeń podał. Rozjazd znaczy, że Operator dostaje
		// zielone światło do pary, która progu nie spełnia.
		if wynik.Result.PassesAA != (wynik.Result.Ratio >= 4.5) {
			t.Errorf("para %s na %s: kontrast %.3f, a orzeczenie AA to %v",
				para.pierwszoplanowa, para.tlo, wynik.Result.Ratio, wynik.Result.PassesAA)
		}
		if wynik.Result.PassesAAA != (wynik.Result.Ratio >= 7) {
			t.Errorf("para %s na %s: kontrast %.3f, a orzeczenie AAA to %v",
				para.pierwszoplanowa, para.tlo, wynik.Result.Ratio, wynik.Result.PassesAAA)
		}
	}

	// Czerń na białym ma kontrast 21 — najwyższy, jaki wzór daje. Liczba inna
	// znaczy, że rachunek nie jest rachunkiem WCAG.
	var skrajny shared.DesignColorContrastCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignColorContrastCheck,
		shared.DesignColorContrastCheckRequest{Foreground: "#000000", Background: "#ffffff"},
		&skrajny)
	if math.Abs(skrajny.Result.Ratio-21) > 0.05 {
		t.Errorf("kontrast czerni na białym wyszedł %.3f, a wzór WCAG daje dokładnie 21",
			skrajny.Result.Ratio)
	}
}

// TestHarmonieBarwMajaZmierzoneKatyOdcieni mierzy KĄTY odcieni palety: triada ma
// trzy odcienie co 120 stopni, dopełnienie — dwa co 180. Paleta zwracająca trzy
// razy barwę wiodącą jest odpowiedzią udaną i nie jest harmonią.
func TestHarmonieBarwMajaZmierzoneKatyOdcieni(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	przypadki := []struct {
		harmonia shared.DesignColorHarmony
		ile      int
		obroty   []float64
	}{
		{shared.DesignColorHarmonyTriad, 3, []float64{0, 120, 240}},
		{shared.DesignColorHarmonyComplementary, 2, []float64{0, 180}},
		{shared.DesignColorHarmonyTetrad, 4, []float64{0, 90, 180, 270}},
	}
	for _, przypadek := range przypadki {
		// Barwa wiodąca jest wybrana z zamysłem: jej nasycenie jest takie, że po
		// obrocie odcienia o dowolny kąt barwa nadal MIEŚCI SIĘ w zakresie
		// wyświetlacza. Barwa nasycona mocno wychodziłaby przy części obrotów poza
		// zakres, przycięcie przesuwałoby odcień o kilkanaście stopni i sprawdzian
		// mierzyłby zakres wyświetlacza zamiast reguły harmonii.
		var wynik shared.DesignColorPaletteGenerateResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandDesignColorPaletteGenerate,
			shared.DesignColorPaletteGenerateRequest{
				BaseColor: "#8f7f9f", Harmony: przypadek.harmonia,
				Count: wskaznik(przypadek.ile),
			}, &wynik)
		if len(wynik.Colors) != przypadek.ile {
			t.Fatalf("harmonia %s o %d barwach oddała %d",
				przypadek.harmonia, przypadek.ile, len(wynik.Colors))
		}
		wiodaca, err := colorful.Hex(wynik.BaseColor)
		if err != nil {
			t.Fatalf("barwa wiodąca %q nieczytelna: %v", wynik.BaseColor, err)
		}
		odcienWiodacy, _, _ := wiodaca.Hcl()
		for numer, wpis := range wynik.Colors {
			barwa, err := colorful.Hex(wpis.Hex)
			if err != nil {
				t.Fatalf("barwa palety %q nieczytelna: %v", wpis.Hex, err)
			}
			odcien, _, _ := barwa.Hcl()
			oczekiwany := math.Mod(odcienWiodacy+przypadek.obroty[numer]+360, 360)
			roznica := math.Abs(math.Mod(odcien-oczekiwany+540, 360) - 180)
			// Trzy stopnie tolerancji — tyle, ile bierze zaokrąglenie barwy do zapisu
			// szesnastkowego i powrót z niego. Obrót o 90 stopni ma się od tego różnić
			// bez wątpliwości.
			if roznica > 3 {
				t.Errorf("harmonia %s, barwa %d (%s): odcień %.1f stopni, a obrót o %.0f od "+
					"wiodącej (%.1f) daje %.1f", przypadek.harmonia, numer+1, wpis.Hex,
					odcien, przypadek.obroty[numer], odcienWiodacy, oczekiwany)
			}
			if wpis.Role == nil || strings.TrimSpace(*wpis.Role) == "" {
				t.Errorf("barwa %d harmonii %s nie ma nazwanej roli w palecie",
					numer+1, przypadek.harmonia)
			}
		}
	}
}

// ── Obrysowanie konturów: pole `smoothing` ma skutek w pliku ────────────────

// obrazSchodkowyPNG składa obraz z ukośnym schodkiem: obszar o krawędzi łamanej,
// na której wygładzenie MA co zaokrąglić.
func obrazSchodkowyPNG(t *testing.T, bok int) []byte {
	t.Helper()

	plotno := image.NewRGBA(image.Rect(0, 0, bok, bok))
	tlo := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	kształt := color.RGBA{R: 0x10, G: 0x30, B: 0x90, A: 0xff}
	for y := 0; y < bok; y++ {
		for x := 0; x < bok; x++ {
			if x+y < bok {
				plotno.SetRGBA(x, y, kształt)
				continue
			}
			plotno.SetRGBA(x, y, tlo)
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// TestWygladzenieObrysuKonturowMaSkutekWPliku mierzy różnicę między obrysowaniem
// bez wygładzenia i z wygładzeniem: bez niego dokument nie ma ani jednej krzywej,
// z nim — ma. Pole przyjmowane i nieużywane obiecuje Operatorowi czynność,
// której nie ma.
func TestWygladzenieObrysuKonturowMaSkutekWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-wektoryzacji",
		obrazSchodkowyPNG(t, 48))

	dokumentDla := func(wygladzenie *float64) (string, int) {
		var wynik shared.DesignPhotoVectorizeResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoVectorize,
			shared.DesignPhotoVectorizeRequest{
				AssetId: zrodlo.Id, Colors: wskaznik(2), Smoothing: wygladzenie,
			}, &wynik)
		if wynik.PathCount < 1 {
			t.Fatalf("obrysowanie konturów oddało %d ścieżek", wynik.PathCount)
		}
		if wynik.Asset.Uri == nil {
			t.Fatal("wynik obrysowania bez odwołania do treści")
		}
		return string(bajtyPodOdwolaniem(t, katalog, *wynik.Asset.Uri)), wynik.PathCount
	}

	bezWygladzenia, sciezekBez := dokumentDla(nil)
	zWygladzeniem, sciezekZ := dokumentDla(wskaznik(2.0))

	if strings.Contains(bezWygladzenia, "Q") {
		t.Error("obrysowanie BEZ wygładzenia niesie krzywe — obrys bez wskazanego wygładzenia " +
			"ma być pomiarem obrazu, nie kształtem zaokrąglonym o liczbę, której nikt nie podał")
	}
	if !strings.Contains(zWygladzeniem, "Q") {
		t.Error("obrysowanie Z wygładzeniem 2 nie niesie ani jednej krzywej kwadratowej — " +
			"pole smoothing jest przyjmowane i nie ma skutku w pliku")
	}
	if sciezekBez != sciezekZ {
		t.Errorf("wygładzenie zmieniło liczbę ścieżek z %d na %d — ma zaokrąglać naroża, "+
			"a nie gubić obszary", sciezekBez, sciezekZ)
	}
	// Kształt zostaje kształtem: obrys nie wychodzi poza obraz ani się nie zwija.
	for _, dokument := range []string{bezWygladzenia, zWygladzeniem} {
		if !strings.Contains(dokument, `viewBox="0 0 48 48"`) {
			t.Errorf("dokument obrysu nie ma pola widoku obrazu 48×48; początek: %s",
				dokument[:min(200, len(dokument))])
		}
		if !strings.Contains(dokument, "<path") {
			t.Error("dokument obrysu nie ma ani jednej ścieżki")
		}
	}
}

// TestObrysObszaruObchodziDziureWPrzeciwnaStrone mierzy dwie własności obrysu, na
// których stoi poprawność wypełnienia w SVG: obszar z dziurą daje DWA kontury,
// a dziura biegnie w stronę przeciwną do konturu zewnętrznego. Bez tego dziura
// wypełniłaby się barwą obszaru i pierścień wyszedłby kołem.
func TestObrysObszaruObchodziDziureWPrzeciwnaStrone(t *testing.T) {
	// Pierścień 5×5 z dziurą 3×3 w środku, zapisany odcinkami wierszy.
	pierscien := [][3]int{
		{0, 0, 4},
		{1, 0, 0}, {1, 4, 4},
		{2, 0, 0}, {2, 4, 4},
		{3, 0, 0}, {3, 4, 4},
		{4, 0, 4},
	}
	kontury := obrysObszaruDesignu(pierscien)
	if len(kontury) != 2 {
		t.Fatalf("obrys pierścienia dał %d konturów, a pierścień ma zewnętrzny i dziurę",
			len(kontury))
	}
	// Pole ze wzoru na wielokąt (podwojone): znak mówi o kierunku obchodzenia.
	poleZeZnakiem := func(wierzcholki []image.Point) int {
		suma := 0
		for numer := range wierzcholki {
			biezacy := wierzcholki[numer]
			nastepny := wierzcholki[(numer+1)%len(wierzcholki)]
			suma += biezacy.X*nastepny.Y - nastepny.X*biezacy.Y
		}
		return suma
	}
	zewnetrzny := poleZeZnakiem(kontury[0])
	dziura := poleZeZnakiem(kontury[1])
	if zewnetrzny == 0 || dziura == 0 {
		t.Fatalf("kontur o zerowym polu: zewnętrzny %d, dziura %d", zewnetrzny, dziura)
	}
	if (zewnetrzny > 0) == (dziura > 0) {
		t.Errorf("kontur zewnętrzny (pole %d) i dziura (pole %d) biegną w tę samą stronę — "+
			"reguła niezerowa wypełni wtedy dziurę barwą obszaru", zewnetrzny, dziura)
	}
	// Kwadrat 5×5 ma pole 25, dziura 3×3 — pole 9. Podwojone: 50 i 18.
	if zewnetrzny != 50 && zewnetrzny != -50 {
		t.Errorf("kontur zewnętrzny ma podwojone pole %d, a kwadrat 5×5 daje 50", zewnetrzny)
	}
	if dziura != 18 && dziura != -18 {
		t.Errorf("dziura ma podwojone pole %d, a kwadrat 3×3 daje 18", dziura)
	}
}

// ── Odszumienie: krawędź zostaje krawędzią ──────────────────────────────────

// obrazZKrawedziaIZiarnemPNG składa obraz o dwóch płaskich połowach różniących
// się jasnością, z ziarnem naniesionym na obie. Materiał, na którym da się
// zmierzyć jedno i drugie: czy ziarno zeszło i czy krawędź została.
func obrazZKrawedziaIZiarnemPNG(t *testing.T, bok int) []byte {
	t.Helper()

	plotno := image.NewRGBA(image.Rect(0, 0, bok, bok))
	for y := 0; y < bok; y++ {
		for x := 0; x < bok; x++ {
			podstawa := 40
			if x >= bok/2 {
				podstawa = 200
			}
			// Ziarno naprzemienne: deterministyczne, więc sprawdzian mierzy zawsze
			// ten sam obraz.
			if (x+y)%2 == 0 {
				podstawa += 12
			} else {
				podstawa -= 12
			}
			plotno.SetRGBA(x, y, color.RGBA{
				R: uint8(podstawa), G: uint8(podstawa), B: uint8(podstawa), A: 0xff,
			})
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// odchylenieObszaruSprawdzianu liczy odchylenie jasności w prostokącie obrazu.
func odchylenieObszaruSprawdzianu(obraz image.Image, obszar image.Rectangle) float64 {
	suma, sumaKwadratow, punktow := 0.0, 0.0, 0.0
	for y := obszar.Min.Y; y < obszar.Max.Y; y++ {
		for x := obszar.Min.X; x < obszar.Max.X; x++ {
			r, g, b, _ := obraz.At(x, y).RGBA()
			jasnosc := float64(int(r>>8)*299+int(g>>8)*587+int(b>>8)*114) / 1000
			suma += jasnosc
			sumaKwadratow += jasnosc * jasnosc
			punktow++
		}
	}
	if punktow == 0 {
		return 0
	}
	srednia := suma / punktow
	return math.Sqrt(math.Max(0, sumaKwadratow/punktow-srednia*srednia))
}

// sredniaObszaruSprawdzianu liczy średnią jasność w prostokącie obrazu.
func sredniaObszaruSprawdzianu(obraz image.Image, obszar image.Rectangle) float64 {
	suma, punktow := 0.0, 0.0
	for y := obszar.Min.Y; y < obszar.Max.Y; y++ {
		for x := obszar.Min.X; x < obszar.Max.X; x++ {
			r, g, b, _ := obraz.At(x, y).RGBA()
			suma += float64(int(r>>8)*299+int(g>>8)*587+int(b>>8)*114) / 1000
			punktow++
		}
	}
	if punktow == 0 {
		return 0
	}
	return suma / punktow
}

// TestOdszumienieZdejmujeZiarnoIZostawiaKrawedz mierzy w PLIKU dwie rzeczy
// naraz: odchylenie jasności w płaskim obszarze (ziarno ma zejść) i skok
// jasności na krawędzi (krawędź ma zostać). Rozmycie zdejmuje jedno i drugie —
// filtr bilateralny tylko pierwsze, i to jest cała różnica między nimi.
func TestOdszumienieZdejmujeZiarnoIZostawiaKrawedz(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	const bok = 40
	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-odszumiania",
		obrazZKrawedziaIZiarnemPNG(t, bok))

	var wynik shared.DesignPhotoEnhanceResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoEnhance,
		shared.DesignPhotoEnhanceRequest{
			AssetId:    zrodlo.Id,
			AutoLevels: wskaznik(false),
			Denoise:    wskaznik(1.0),
			Sharpen:    wskaznik(0.0),
		}, &wynik)
	maOdszumienie := false
	for _, krok := range wynik.AppliedSteps {
		if strings.HasPrefix(krok, "odszumienie") {
			maOdszumienie = true
		}
	}
	if !maOdszumienie {
		t.Fatalf("bilans kroków nie wymienia odszumienia: %v", wynik.AppliedSteps)
	}

	przed := obrazZMagazynuFotografii(t, katalog, zrodlo)
	po := obrazZMagazynuFotografii(t, katalog, wynik.Asset)

	// Obszar płaski wzięty z dala od krawędzi — tam ziarno jest jedyną zmiennością.
	plaski := image.Rect(2, 2, bok/2-6, bok-2)
	odchyleniePrzed := odchylenieObszaruSprawdzianu(przed, plaski)
	odchyleniePo := odchylenieObszaruSprawdzianu(po, plaski)
	if odchyleniePrzed <= 1 {
		t.Fatalf("obraz źródłowy ma w płaskim obszarze odchylenie %.2f — sprawdzian nie ma "+
			"czego odszumiać", odchyleniePrzed)
	}
	if odchyleniePo > odchyleniePrzed/2 {
		t.Errorf("po odszumieniu odchylenie w płaskim obszarze to %.2f, a przed było %.2f — "+
			"ziarno nie zeszło nawet o połowę", odchyleniePo, odchyleniePrzed)
	}

	// Krawędź: dwa pasy po jej obu stronach. Skok między nimi ma zostać.
	lewy := image.Rect(bok/2-3, 2, bok/2, bok-2)
	prawy := image.Rect(bok/2, 2, bok/2+3, bok-2)
	skokPrzed := sredniaObszaruSprawdzianu(przed, prawy) - sredniaObszaruSprawdzianu(przed, lewy)
	skokPo := sredniaObszaruSprawdzianu(po, prawy) - sredniaObszaruSprawdzianu(po, lewy)
	if skokPrzed <= 100 {
		t.Fatalf("obraz źródłowy ma na krawędzi skok %.1f — sprawdzian nie ma czego bronić",
			skokPrzed)
	}
	// Osiemdziesiąt procent skoku: rozmycie o promieniu trzech punktów zjada go
	// na tych pasach w dużej części, a filtr bilateralny prawie nie rusza.
	if skokPo < 0.8*skokPrzed {
		t.Errorf("po odszumieniu skok na krawędzi to %.1f, a przed był %.1f — odszumienie "+
			"rozmyło krawędź zamiast ją zachować", skokPo, skokPrzed)
	}
}
