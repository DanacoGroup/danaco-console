// Odpowiedzialność pliku: rachunek na ścieżkach wektorowych modułu Design —
// przekład węzłów kontraktu na ścieżkę biblioteki i z powrotem, serializacja
// wypełnienia i obrysu do zapisu w bazie, złożenie dokumentu SVG oraz wydanie
// PDF i EPS. Czynności kontraktu stoją w `adapter_modul_design_wektor.go`;
// tutaj leży sam rachunek, żeby miał jedno miejsce i jedną prawdę.
//
// ── Rachunek jest wkompilowany, nie wołany ──────────────────────────────────
// Operacje logiczne na ścieżkach, wydanie PDF i wydanie EPS idą przez
// `tdewolff/canvas` — bibliotekę Go wkompilowaną w binarium serwera. Nie ma tu
// ani jednego uruchomienia procesu i mieć nie będzie: programy do obrysowywania
// konturów, rysowania wektorowego i rasteryzacji języka strony leżą poza
// instalką Operatora, więc funkcja od nich zależna byłaby u niego odmową, nie
// funkcją.
//
// ── Węzeł jest bytem produktu, ścieżka biblioteki tylko rachunkiem ──────────
// Kontrakt niesie węzły z uchwytami (`DesignVectorNode`) i to one leżą w bazie.
// Ścieżka biblioteki powstaje na czas rachunku i ginie po nim; wynik wraca
// znowu jako węzły, żeby Operator mógł go dalej ciągnąć piórem. Zapisanie
// wyniku jako gotowego napisu SVG odebrałoby mu edycję — kształt przestałby
// mieć węzły, a zostałby obrazkiem.
//
// ── Uchwyt jest ODSUNIĘCIEM, nie punktem ────────────────────────────────────
// `handleInX`/`handleOutX` kontraktu są odsunięciami od węzła (tak opisuje je
// treść pola). Punkt sterujący krzywej to więc węzeł plus odsunięcie. Odczyt
// odsunięcia jako współrzędnej bezwzględnej przesuwałby krzywe ku początkowi
// układu przy każdym przejściu przez bazę.
package core

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/tdewolff/canvas"

	"danacoconsole/shared"
)

const (
	// tolerancjaSciezkiDesignu jest dopuszczalnym odejściem od kształtu
	// pierwotnego przy zamianie łuku na krzywe i przy spłaszczaniu. Wartość jest
	// w jednostkach kompozycji (pikselach kanwy) — dziesiąta część jednostki
	// jest niewidoczna na ekranie i w druku.
	tolerancjaSciezkiDesignu = 0.1

	// domyslnaPrecyzjaSciezkiDesignu jest liczbą miejsc po przecinku
	// współrzędnych przy czyszczeniu zapisu (`design.vector.optimize`).
	domyslnaPrecyzjaSciezkiDesignu = 2

	// granicaPrecyzjiSciezkiDesignu jest kresem sensu: powyżej dwunastu miejsc
	// po przecinku zapis liczby zmiennoprzecinkowej przestaje nieść informację,
	// a rośnie w bajtach.
	granicaPrecyzjiSciezkiDesignu = 12
)

// wezlyZeZapisuDesignu rozkłada zapis JSON węzłów z kolumny na węzły kontraktu.
func wezlyZeZapisuDesignu(zapis string) ([]shared.DesignVectorNode, error) {
	var wezly []shared.DesignVectorNode
	if err := json.Unmarshal([]byte(zapis), &wezly); err != nil {
		return nil, fmt.Errorf("zapis węzłów ścieżki w bazie nie jest wykazem węzłów: %w", err)
	}
	return wezly, nil
}

// zapisWezlowDesignu składa zapis JSON węzłów do kolumny.
func zapisWezlowDesignu(wezly []shared.DesignVectorNode) (string, error) {
	bajty, err := json.Marshal(wezly)
	if err != nil {
		return "", fmt.Errorf("nie można złożyć zapisu węzłów ścieżki: %w", err)
	}
	return string(bajty), nil
}

// zapisWypelnieniaDesignu i zapisObrysuDesignu składają zapis JSON wypełnienia
// i obrysu. Brak wskazania zostaje brakiem: ścieżka bez wypełnienia jest
// ścieżką bez wypełnienia, a nie ścieżką wypełnioną czernią domyślną.
func zapisWypelnieniaDesignu(wypelnienie *shared.DesignFill) (*string, error) {
	if wypelnienie == nil {
		return nil, nil
	}
	bajty, err := json.Marshal(wypelnienie)
	if err != nil {
		return nil, fmt.Errorf("nie można złożyć zapisu wypełnienia ścieżki: %w", err)
	}
	zapis := string(bajty)
	return &zapis, nil
}

func zapisObrysuDesignu(obrys *shared.DesignStroke) (*string, error) {
	if obrys == nil {
		return nil, nil
	}
	bajty, err := json.Marshal(obrys)
	if err != nil {
		return nil, fmt.Errorf("nie można złożyć zapisu obrysu ścieżki: %w", err)
	}
	zapis := string(bajty)
	return &zapis, nil
}

// wypelnienieZeZapisuDesignu i obrysZeZapisuDesignu czytają zapis z kolumny.
// Zapis nieczytelny jest tu brakiem, nie odmową: wypełnienie jest cechą
// wtórną, a odmowa odczytu ścieżki z powodu jej barwy odebrałaby Operatorowi
// kształt, który narysował.
func wypelnienieZeZapisuDesignu(zapis *string) *shared.DesignFill {
	if zapis == nil || strings.TrimSpace(*zapis) == "" {
		return nil
	}
	var wypelnienie shared.DesignFill
	if err := json.Unmarshal([]byte(*zapis), &wypelnienie); err != nil {
		return nil
	}
	return &wypelnienie
}

func obrysZeZapisuDesignu(zapis *string) *shared.DesignStroke {
	if zapis == nil || strings.TrimSpace(*zapis) == "" {
		return nil
	}
	var obrys shared.DesignStroke
	if err := json.Unmarshal([]byte(*zapis), &obrys); err != nil {
		return nil
	}
	return &obrys
}

// sciezkaBibliotekiDesignu składa ścieżkę biblioteki z węzłów kontraktu.
//
// Odcinek między dwoma węzłami jest krzywą sześcienną, gdy którykolwiek z nich
// niesie uchwyt po tej stronie odcinka; inaczej jest odcinkiem prostym. Brak
// uchwytu z jednej strony bierze punkt sterujący na samym węźle — tak działa
// pióro w każdym programie wektorowym: węzeł narożny z jednej strony
// i wygładzony z drugiej daje krzywą, która z jednej strony wchodzi prosto.
func sciezkaBibliotekiDesignu(wezly []shared.DesignVectorNode, zamknieta bool) *canvas.Path {
	sciezka := &canvas.Path{}
	if len(wezly) == 0 {
		return sciezka
	}
	sciezka.MoveTo(wezly[0].X, wezly[0].Y)
	for numer := 1; numer < len(wezly); numer++ {
		dopiszOdcinekSciezkiDesignu(sciezka, wezly[numer-1], wezly[numer])
	}
	if zamknieta && len(wezly) > 2 {
		dopiszOdcinekSciezkiDesignu(sciezka, wezly[len(wezly)-1], wezly[0])
		sciezka.Close()
	}
	return sciezka
}

// dopiszOdcinekSciezkiDesignu dokłada jeden odcinek między parą węzłów.
func dopiszOdcinekSciezkiDesignu(sciezka *canvas.Path, od, do shared.DesignVectorNode) {
	wyjscieX, wyjscieY, maWyjscie := uchwytWyjsciaDesignu(od)
	wejscieX, wejscieY, maWejscie := uchwytWejsciaDesignu(do)
	if !maWyjscie && !maWejscie {
		sciezka.LineTo(do.X, do.Y)
		return
	}
	pierwszyX, pierwszyY := od.X, od.Y
	if maWyjscie {
		pierwszyX, pierwszyY = od.X+wyjscieX, od.Y+wyjscieY
	}
	drugiX, drugiY := do.X, do.Y
	if maWejscie {
		drugiX, drugiY = do.X+wejscieX, do.Y+wejscieY
	}
	sciezka.CubeTo(pierwszyX, pierwszyY, drugiX, drugiY, do.X, do.Y)
}

// uchwytWyjsciaDesignu i uchwytWejsciaDesignu oddają odsunięcie uchwytu wraz
// z rozstrzygnięciem, czy węzeł go w ogóle ma. Uchwyt zerowy jest uchwytem
// wskazanym wprost i nie jest tym samym co uchwyt nieustawiony, ale dla
// rachunku daje ten sam punkt sterujący — rozróżnienie ma znaczenie przy
// odczycie, nie przy rysowaniu.
func uchwytWyjsciaDesignu(w shared.DesignVectorNode) (float64, float64, bool) {
	if w.HandleOutX == nil && w.HandleOutY == nil {
		return 0, 0, false
	}
	x, y := 0.0, 0.0
	if w.HandleOutX != nil {
		x = *w.HandleOutX
	}
	if w.HandleOutY != nil {
		y = *w.HandleOutY
	}
	return x, y, true
}

func uchwytWejsciaDesignu(w shared.DesignVectorNode) (float64, float64, bool) {
	if w.HandleInX == nil && w.HandleInY == nil {
		return 0, 0, false
	}
	x, y := 0.0, 0.0
	if w.HandleInX != nil {
		x = *w.HandleInX
	}
	if w.HandleInY != nil {
		y = *w.HandleInY
	}
	return x, y, true
}

// wezlyZeSciezkiBibliotekiDesignu rozkłada ścieżkę biblioteki na węzły
// kontraktu wraz z rozstrzygnięciem, czy jest zamknięta.
//
// Łuki wchodzą jako krzywe (`ReplaceArcs`), bo kontrakt nie ma węzła
// łukowego — a łuk zamilczany zgubiłby kawałek kształtu. Wielościeżkowy wynik
// operacji logicznej daje jeden wykaz węzłów: kontrakt niesie jedną ścieżkę
// wynikową (`DesignVectorBooleanResponse.Path`), więc rozdzielone kawałki idą
// po sobie, a nie giną.
func wezlyZeSciezkiBibliotekiDesignu(sciezka *canvas.Path) ([]shared.DesignVectorNode, bool) {
	if sciezka == nil || sciezka.Empty() {
		return nil, false
	}
	wezly := []shared.DesignVectorNode{}
	zamknieta := false
	skaner := sciezka.ReplaceArcs().Scanner()
	for skaner.Scan() {
		switch skaner.Cmd() {
		case canvas.MoveToCmd:
			koniec := skaner.End()
			wezly = append(wezly, shared.DesignVectorNode{
				X: koniec.X, Y: koniec.Y, Kind: shared.DesignVectorNodeKindCorner,
			})
		case canvas.LineToCmd:
			koniec := skaner.End()
			wezly = append(wezly, shared.DesignVectorNode{
				X: koniec.X, Y: koniec.Y, Kind: shared.DesignVectorNodeKindCorner,
			})
		case canvas.QuadToCmd:
			// Krzywa kwadratowa idzie do postaci sześciennej, bo węzeł kontraktu
			// ma dwa uchwyty, nie jeden wspólny punkt sterujący. Przeliczenie
			// jest dokładne: dwie trzecie drogi od każdego końca do punktu
			// sterującego dają tę samą krzywą.
			poczatek, sterujacy, koniec := skaner.Start(), skaner.CP1(), skaner.End()
			pierwszy := canvas.Point{
				X: poczatek.X + 2.0/3.0*(sterujacy.X-poczatek.X),
				Y: poczatek.Y + 2.0/3.0*(sterujacy.Y-poczatek.Y),
			}
			drugi := canvas.Point{
				X: koniec.X + 2.0/3.0*(sterujacy.X-koniec.X),
				Y: koniec.Y + 2.0/3.0*(sterujacy.Y-koniec.Y),
			}
			wezly = dopiszWezelKrzywejDesignu(wezly, poczatek, pierwszy, drugi, koniec)
		case canvas.CubeToCmd:
			wezly = dopiszWezelKrzywejDesignu(wezly,
				skaner.Start(), skaner.CP1(), skaner.CP2(), skaner.End())
		case canvas.CloseCmd:
			zamknieta = true
			// Domknięcie wracające do punktu startowego nie dokłada węzła —
			// węzeł w tym samym miejscu, co pierwszy, byłby drugim zapisem
			// jednego narożnika.
			koniec := skaner.End()
			if len(wezly) > 0 && !bliskoDesignu(wezly[0].X, koniec.X) {
				wezly = append(wezly, shared.DesignVectorNode{
					X: koniec.X, Y: koniec.Y, Kind: shared.DesignVectorNodeKindCorner,
				})
			}
		}
	}
	if len(wezly) > 1 && zamknieta && bliskoDesignu(wezly[0].X, wezly[len(wezly)-1].X) &&
		bliskoDesignu(wezly[0].Y, wezly[len(wezly)-1].Y) {

		// Ostatni węzeł pokrywający się z pierwszym przy ścieżce zamkniętej jest
		// zbędny: domknięcie samo prowadzi z ostatniego do pierwszego. Uchwyt
		// wchodzący zostaje jednak przeniesiony, bo opisuje krzywiznę domknięcia.
		ostatni := wezly[len(wezly)-1]
		wezly[0].HandleInX, wezly[0].HandleInY = ostatni.HandleInX, ostatni.HandleInY
		wezly = wezly[:len(wezly)-1]
	}
	return wezly, zamknieta
}

// dopiszWezelKrzywejDesignu dokłada węzeł końcowy krzywej i uzupełnia uchwyt
// wyjścia węzła poprzedniego. Rodzaj węzła poprzedniego przelicza się na nowo:
// węzeł, który ma teraz uchwyt po obu stronach, przestał być narożnikiem.
func dopiszWezelKrzywejDesignu(wezly []shared.DesignVectorNode,
	poczatek, pierwszy, drugi, koniec canvas.Point) []shared.DesignVectorNode {

	if len(wezly) == 0 {
		wezly = append(wezly, shared.DesignVectorNode{
			X: poczatek.X, Y: poczatek.Y, Kind: shared.DesignVectorNodeKindCorner,
		})
	}
	poprzedni := &wezly[len(wezly)-1]
	if !bliskoDesignu(pierwszy.X, poprzedni.X) || !bliskoDesignu(pierwszy.Y, poprzedni.Y) {
		odsuniecieX, odsuniecieY := pierwszy.X-poprzedni.X, pierwszy.Y-poprzedni.Y
		poprzedni.HandleOutX, poprzedni.HandleOutY = &odsuniecieX, &odsuniecieY
		poprzedni.Kind = rodzajWezlaDesignu(*poprzedni)
	}
	nowy := shared.DesignVectorNode{X: koniec.X, Y: koniec.Y}
	if !bliskoDesignu(drugi.X, koniec.X) || !bliskoDesignu(drugi.Y, koniec.Y) {
		odsuniecieX, odsuniecieY := drugi.X-koniec.X, drugi.Y-koniec.Y
		nowy.HandleInX, nowy.HandleInY = &odsuniecieX, &odsuniecieY
	}
	nowy.Kind = rodzajWezlaDesignu(nowy)
	return append(wezly, nowy)
}

// rodzajWezlaDesignu rozstrzyga rodzaj węzła z jego uchwytów. Uchwyty
// przeciwległe i równej długości to węzeł symetryczny; przeciwległe różnej
// długości — wygładzony; wszystko inne — narożnik.
func rodzajWezlaDesignu(w shared.DesignVectorNode) shared.DesignVectorNodeKind {
	wejscieX, wejscieY, maWejscie := uchwytWejsciaDesignu(w)
	wyjscieX, wyjscieY, maWyjscie := uchwytWyjsciaDesignu(w)
	if !maWejscie || !maWyjscie {
		return shared.DesignVectorNodeKindCorner
	}
	dlugoscWejscia := math.Hypot(wejscieX, wejscieY)
	dlugoscWyjscia := math.Hypot(wyjscieX, wyjscieY)
	if dlugoscWejscia == 0 || dlugoscWyjscia == 0 {
		return shared.DesignVectorNodeKindCorner
	}
	// Iloczyn wektorowy bliski zeru znaczy uchwyty na jednej prostej; iloczyn
	// skalarny ujemny — po przeciwnych stronach węzła.
	wektorowy := wejscieX*wyjscieY - wejscieY*wyjscieX
	skalarny := wejscieX*wyjscieX + wejscieY*wyjscieY
	wspolliniowe := math.Abs(wektorowy) <= tolerancjaSciezkiDesignu*dlugoscWejscia*dlugoscWyjscia
	if !wspolliniowe || skalarny >= 0 {
		return shared.DesignVectorNodeKindCorner
	}
	if math.Abs(dlugoscWejscia-dlugoscWyjscia) <= tolerancjaSciezkiDesignu {
		return shared.DesignVectorNodeKindSymmetric
	}
	return shared.DesignVectorNodeKindSmooth
}

// bliskoDesignu rozstrzyga równość współrzędnych z tolerancją rachunku
// zmiennoprzecinkowego. Porównanie dokładne dawałoby uchwyty o długości 1e-16
// przy każdym przejściu przez bibliotekę.
func bliskoDesignu(pierwsza, druga float64) bool {
	return math.Abs(pierwsza-druga) <= tolerancjaSciezkiDesignu/10
}

// zapisSvgSciezkiDesignu składa atrybut `d` ścieżki SVG z węzłów kontraktu —
// przez ścieżkę biblioteki, żeby zapis SVG i rachunek logiczny opisywały
// dokładnie ten sam kształt.
func zapisSvgSciezkiDesignu(wezly []shared.DesignVectorNode, zamknieta bool) string {
	return sciezkaBibliotekiDesignu(wezly, zamknieta).ToSVG()
}

// przytnijPrecyzjeWezlowDesignu zaokrągla współrzędne i uchwyty do zadanej
// liczby miejsc po przecinku — rachunek czyszczenia `design.vector.optimize`.
//
// Czyszczenie nie usuwa węzłów. Węzeł, którego Operator postawił, jest jego
// rozstrzygnięciem o kształcie; ubytek bajtów bierze się z krótszego zapisu
// liczb, a nie z gubienia jego pracy.
func przytnijPrecyzjeWezlowDesignu(wezly []shared.DesignVectorNode,
	miejsca int) []shared.DesignVectorNode {

	mnoznik := math.Pow(10, float64(miejsca))
	zaokraglij := func(wartosc float64) float64 {
		return math.Round(wartosc*mnoznik) / mnoznik
	}
	zaokraglijWskaznik := func(wartosc *float64) *float64 {
		if wartosc == nil {
			return nil
		}
		nowa := zaokraglij(*wartosc)
		return &nowa
	}
	wynik := make([]shared.DesignVectorNode, 0, len(wezly))
	for _, wezel := range wezly {
		wezel.X = zaokraglij(wezel.X)
		wezel.Y = zaokraglij(wezel.Y)
		wezel.HandleInX = zaokraglijWskaznik(wezel.HandleInX)
		wezel.HandleInY = zaokraglijWskaznik(wezel.HandleInY)
		wezel.HandleOutX = zaokraglijWskaznik(wezel.HandleOutX)
		wezel.HandleOutY = zaokraglijWskaznik(wezel.HandleOutY)
		wynik = append(wynik, wezel)
	}
	return wynik
}

// ksztaltWezlamiDesignu składa węzły kształtu podstawowego — obsługuje
// `design.vector.shape.add`.
//
// Kształt powstaje OD RAZU jako węzły ścieżki, a nie jako osobny byt do
// późniejszej zamiany: prostokąt dorysowany na kanwie ma dać się natychmiast
// ciągnąć piórem za narożnik, bez komendy „zamień w ścieżkę", której kontrakt
// nie ma i mieć nie będzie.
func ksztaltWezlamiDesignu(z shared.DesignVectorShapeAddRequest) ([]shared.DesignVectorNode, bool) {
	switch z.Kind {
	case shared.DesignShapeKindRectangle:
		if z.CornerRadius != nil && *z.CornerRadius > 0 {
			return prostokatZaokraglonyDesignu(z, *z.CornerRadius), true
		}
		return []shared.DesignVectorNode{
			{X: z.X, Y: z.Y, Kind: shared.DesignVectorNodeKindCorner},
			{X: z.X + z.Width, Y: z.Y, Kind: shared.DesignVectorNodeKindCorner},
			{X: z.X + z.Width, Y: z.Y + z.Height, Kind: shared.DesignVectorNodeKindCorner},
			{X: z.X, Y: z.Y + z.Height, Kind: shared.DesignVectorNodeKindCorner},
		}, true
	case shared.DesignShapeKindEllipse:
		return elipsaWezlamiDesignu(z.X, z.Y, z.Width, z.Height), true
	case shared.DesignShapeKindPolygon:
		return wielokatWezlamiDesignu(z, liczbaWierzcholkowDesignu(z.Points, 3), 1), true
	case shared.DesignShapeKindStar:
		ramiona := liczbaWierzcholkowDesignu(z.Points, 5)
		udzialWewnetrzny := 0.5
		if z.InnerRadius != nil && *z.InnerRadius > 0 {
			// Promień wewnętrzny podaje się w jednostkach kompozycji, a rachunek
			// gwiazdy potrzebuje udziału względem promienia zewnętrznego.
			zewnetrzny := mniejszaDesignu(z.Width, z.Height) / 2
			if zewnetrzny > 0 {
				udzialWewnetrzny = przytnijUlamekDesignu(*z.InnerRadius / zewnetrzny)
			}
		}
		return gwiazdaWezlamiDesignu(z, ramiona, udzialWewnetrzny), true
	case shared.DesignShapeKindLine:
		return []shared.DesignVectorNode{
			{X: z.X, Y: z.Y, Kind: shared.DesignVectorNodeKindCorner},
			{X: z.X + z.Width, Y: z.Y + z.Height, Kind: shared.DesignVectorNodeKindCorner},
		}, false
	}
	return nil, false
}

// liczbaWierzcholkowDesignu przycina liczbę wierzchołków do sensu: wielokąt
// o dwóch wierzchołkach jest odcinkiem, a nie wielokątem.
func liczbaWierzcholkowDesignu(wskazana *int, najmniej int) int {
	if wskazana == nil || *wskazana < najmniej {
		return najmniej
	}
	if *wskazana > 512 {
		return 512
	}
	return *wskazana
}

// wspolczynnikLukuDesignu jest długością uchwytu przybliżającego kwadrans łuku
// krzywą sześcienną. Wartość jest znana z rachunku: 4/3·(√2−1).
const wspolczynnikLukuDesignu = 0.5522847498307936

// elipsaWezlamiDesignu składa elipsę z czterech węzłów gładkich — kwadrans na
// węzeł, uchwyty w stycznych.
func elipsaWezlamiDesignu(x, y, szerokosc, wysokosc float64) []shared.DesignVectorNode {
	promienX, promienY := szerokosc/2, wysokosc/2
	srodekX, srodekY := x+promienX, y+promienY
	uchwytX := promienX * wspolczynnikLukuDesignu
	uchwytY := promienY * wspolczynnikLukuDesignu
	wskaz := func(wartosc float64) *float64 { return &wartosc }
	return []shared.DesignVectorNode{
		{X: srodekX, Y: srodekY - promienY, Kind: shared.DesignVectorNodeKindSymmetric,
			HandleInX: wskaz(-uchwytX), HandleInY: wskaz(0),
			HandleOutX: wskaz(uchwytX), HandleOutY: wskaz(0)},
		{X: srodekX + promienX, Y: srodekY, Kind: shared.DesignVectorNodeKindSymmetric,
			HandleInX: wskaz(0), HandleInY: wskaz(-uchwytY),
			HandleOutX: wskaz(0), HandleOutY: wskaz(uchwytY)},
		{X: srodekX, Y: srodekY + promienY, Kind: shared.DesignVectorNodeKindSymmetric,
			HandleInX: wskaz(uchwytX), HandleInY: wskaz(0),
			HandleOutX: wskaz(-uchwytX), HandleOutY: wskaz(0)},
		{X: srodekX - promienX, Y: srodekY, Kind: shared.DesignVectorNodeKindSymmetric,
			HandleInX: wskaz(0), HandleInY: wskaz(uchwytY),
			HandleOutX: wskaz(0), HandleOutY: wskaz(-uchwytY)},
	}
}

// prostokatZaokraglonyDesignu składa prostokąt o zaokrąglonych narożach: dwa
// węzły na naroże, uchwyty w stycznych.
func prostokatZaokraglonyDesignu(z shared.DesignVectorShapeAddRequest,
	promien float64) []shared.DesignVectorNode {

	granica := mniejszaDesignu(z.Width, z.Height) / 2
	if promien > granica {
		promien = granica
	}
	uchwyt := promien * wspolczynnikLukuDesignu
	wskaz := func(wartosc float64) *float64 { return &wartosc }
	lewa, gora := z.X, z.Y
	prawa, dol := z.X+z.Width, z.Y+z.Height
	return []shared.DesignVectorNode{
		{X: lewa + promien, Y: gora, Kind: shared.DesignVectorNodeKindSmooth,
			HandleInX: wskaz(-uchwyt), HandleInY: wskaz(0)},
		{X: prawa - promien, Y: gora, Kind: shared.DesignVectorNodeKindSmooth,
			HandleOutX: wskaz(uchwyt), HandleOutY: wskaz(0)},
		{X: prawa, Y: gora + promien, Kind: shared.DesignVectorNodeKindSmooth,
			HandleInX: wskaz(0), HandleInY: wskaz(-uchwyt)},
		{X: prawa, Y: dol - promien, Kind: shared.DesignVectorNodeKindSmooth,
			HandleOutX: wskaz(0), HandleOutY: wskaz(uchwyt)},
		{X: prawa - promien, Y: dol, Kind: shared.DesignVectorNodeKindSmooth,
			HandleInX: wskaz(uchwyt), HandleInY: wskaz(0)},
		{X: lewa + promien, Y: dol, Kind: shared.DesignVectorNodeKindSmooth,
			HandleOutX: wskaz(-uchwyt), HandleOutY: wskaz(0)},
		{X: lewa, Y: dol - promien, Kind: shared.DesignVectorNodeKindSmooth,
			HandleInX: wskaz(0), HandleInY: wskaz(uchwyt)},
		{X: lewa, Y: gora + promien, Kind: shared.DesignVectorNodeKindSmooth,
			HandleOutX: wskaz(0), HandleOutY: wskaz(-uchwyt)},
	}
}

// wielokatWezlamiDesignu składa wielokąt wpisany w prostokąt żądania.
func wielokatWezlamiDesignu(z shared.DesignVectorShapeAddRequest, wierzcholkow int,
	udzialWewnetrzny float64) []shared.DesignVectorNode {

	promienX, promienY := z.Width/2, z.Height/2
	srodekX, srodekY := z.X+promienX, z.Y+promienY
	wezly := make([]shared.DesignVectorNode, 0, wierzcholkow)
	for numer := 0; numer < wierzcholkow; numer++ {
		// Pierwszy wierzchołek na górze: wielokąt obrócony o pół kroku wygląda
		// jak przekrzywiony, a nikt o obrót nie prosił.
		kat := -math.Pi/2 + 2*math.Pi*float64(numer)/float64(wierzcholkow)
		wezly = append(wezly, shared.DesignVectorNode{
			X:    srodekX + promienX*udzialWewnetrzny*math.Cos(kat),
			Y:    srodekY + promienY*udzialWewnetrzny*math.Sin(kat),
			Kind: shared.DesignVectorNodeKindCorner,
		})
	}
	return wezly
}

// gwiazdaWezlamiDesignu składa gwiazdę: wierzchołek zewnętrzny i wewnętrzny na
// przemian.
func gwiazdaWezlamiDesignu(z shared.DesignVectorShapeAddRequest, ramion int,
	udzialWewnetrzny float64) []shared.DesignVectorNode {

	promienX, promienY := z.Width/2, z.Height/2
	srodekX, srodekY := z.X+promienX, z.Y+promienY
	wezly := make([]shared.DesignVectorNode, 0, 2*ramion)
	for numer := 0; numer < 2*ramion; numer++ {
		udzial := 1.0
		if numer%2 == 1 {
			udzial = udzialWewnetrzny
		}
		kat := -math.Pi/2 + math.Pi*float64(numer)/float64(ramion)
		wezly = append(wezly, shared.DesignVectorNode{
			X:    srodekX + promienX*udzial*math.Cos(kat),
			Y:    srodekY + promienY*udzial*math.Sin(kat),
			Kind: shared.DesignVectorNodeKindCorner,
		})
	}
	return wezly
}
