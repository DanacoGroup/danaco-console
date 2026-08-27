// Odpowiedzialność pliku: pięć czynności barwy modułu Design liczonych bez
// dotykania obrazu — paleta z harmonii, pomiar kontrastu, gradient,
// przeliczenie zapisu i badanie dostępności zestawu żetonów, rodzina
// `design.color.*`.
package core

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/lucasb-eyer/go-colorful"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// domyslnaLiczbaBarwPaletyDesignu jest liczbą barw palety, gdy Operator jej
	// nie podał. Sześć barw wystarcza na parę tła i tekstu plus akcenty, a wykaz
	// dłuższy przestaje być paletą i staje się skalą.
	domyslnaLiczbaBarwPaletyDesignu = 6

	// granicaBarwPaletyDesignu chroni odpowiedź przed skalą tysiącbarwną, której
	// nikt w oknie nie obejrzy.
	granicaBarwPaletyDesignu = 64

	// przedrostekGradientuDesign znakuje identyfikatory zewnętrzne gradientów,
	// odróżniając je od innych bytów modułu Design.
	przedrostekGradientuDesign = "gradient-"
)

// GenerujPalete liczy paletę z barwy wiodącej i reguły harmonii — obsługuje
// `design.color.palette.generate`.
func (a *adapterDesignu) GenerujPalete(_ context.Context,
	z shared.DesignColorPaletteGenerateRequest) (shared.DesignColorPaletteGenerateResponse, error) {

	if err := sprawdzWyliczenieDesignu("design.color.palette.generate", "harmony", z.Harmony,
		shared.WartosciDesignColorHarmony()); err != nil {
		return shared.DesignColorPaletteGenerateResponse{}, err
	}
	wiodaca, err := rozpoznajBarweDesignu(z.BaseColor)
	if err != nil {
		return shared.DesignColorPaletteGenerateResponse{}, bladWskazaniaDesignu(
			"komenda design.color.palette.generate: " + err.Error())
	}
	ile := domyslnaLiczbaBarwPaletyDesignu
	if z.Count != nil {
		if *z.Count < 1 || *z.Count > granicaBarwPaletyDesignu {
			return shared.DesignColorPaletteGenerateResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.color.palette.generate z liczbą barw %d: rdzeń liczy palety od 1 "+
					"do %d barw", *z.Count, granicaBarwPaletyDesignu))
		}
		ile = *z.Count
	}

	barwy := paletaHarmoniiDesignu(wiodaca, z.Harmony, ile)
	wynik := make([]shared.DesignPaletteColor, 0, len(barwy))
	for numer, barwa := range barwy {
		rola := rolaBarwyPaletyDesignu(numer, len(barwy))
		wpis := shared.DesignPaletteColor{Hex: barwa.Clamped().Hex(), Role: &rola}
		// Nazwa wchodzi tylko przy trafieniu dokładnym — „najbliższa” mówi
		// o barwie, której nie ma w palecie.
		if nazwa := nazwaBarwyDesignu(wpis.Hex); nazwa != nil {
			wpis.Name = nazwa
		}
		wynik = append(wynik, wpis)
	}
	return shared.DesignColorPaletteGenerateResponse{
		Colors: wynik, BaseColor: wiodaca.Clamped().Hex(),
	}, nil
}

// paletaHarmoniiDesignu liczy barwy palety wedle reguły harmonii: reguła daje
// odcienie, a liczba barw żądania rozciąga je jasnością, więc ośmiobarwna
// triada dostaje osiem barw, nie trzy i pustkę.
func paletaHarmoniiDesignu(wiodaca colorful.Color, harmonia shared.DesignColorHarmony,
	ile int) []colorful.Color {

	odcien, nasycenie, jasnosc := wiodaca.Hcl()
	obroty := obrotyHarmoniiDesignu(harmonia)

	barwy := make([]colorful.Color, 0, ile)
	for numer := 0; numer < ile; numer++ {
		obrot := obroty[numer%len(obroty)]
		// Jasność zmienia się dopiero po wyczerpaniu odcieni reguły harmonii.

		// Krok jasności jest ułamkiem, żeby paleta nie wyszła poza skalę.
		okrag := numer / len(obroty)
		przesuniecieJasnosci := 0.0
		if okrag > 0 {
			kierunek := 1.0
			if okrag%2 == 0 {
				kierunek = -1.0
			}
			przesuniecieJasnosci = kierunek * 0.12 * float64((okrag+1)/2)
		}
		nowaJasnosc := przytnijUlamekDesignu(jasnosc + przesuniecieJasnosci)
		nowyOdcien := znormalizujKatDesignu(odcien + obrot)
		nasycenieBarwy := nasycenie
		if harmonia == shared.DesignColorHarmonyMono && okrag == 0 && numer > 0 {
			// Harmonia monochromatyczna nie ma innych odcieni — różnicuje ją
			// nasycenie i jasność.
			nasycenieBarwy = przytnijUlamekDesignu(nasycenie * (1 - 0.15*float64(numer)))
			nowaJasnosc = przytnijUlamekDesignu(jasnosc + 0.1*float64(numer) - 0.2)
		}
		barwy = append(barwy, colorful.Hcl(nowyOdcien, nasycenieBarwy, nowaJasnosc).Clamped())
	}
	return barwy
}

// obrotyHarmoniiDesignu oddaje obroty odcienia dla reguły harmonii. Wykaz
// zaczyna się od zera, bo barwa wiodąca należy do własnej palety.
func obrotyHarmoniiDesignu(harmonia shared.DesignColorHarmony) []float64 {
	switch harmonia {
	case shared.DesignColorHarmonyMono:
		return []float64{0}
	case shared.DesignColorHarmonyAnalogous:
		return []float64{0, 30, -30, 60, -60}
	case shared.DesignColorHarmonyComplementary:
		return []float64{0, 180}
	case shared.DesignColorHarmonySplitComplementary:
		return []float64{0, 150, 210}
	case shared.DesignColorHarmonyTriad:
		return []float64{0, 120, 240}
	case shared.DesignColorHarmonyTetrad:
		return []float64{0, 90, 180, 270}
	}
	return []float64{0}
}

// rolaBarwyPaletyDesignu nazywa rolę barwy w systemie. Nazwy są opisem miejsca
// w palecie, a nie obietnicą użycia: pierwsza barwa jest wiodącą, bo od niej
// paleta powstała.
func rolaBarwyPaletyDesignu(numer, ile int) string {
	switch {
	case numer == 0:
		return "wiodąca"
	case numer == 1 && ile > 2:
		return "wtórna"
	case numer == 2 && ile > 3:
		return "akcent"
	}
	return fmt.Sprintf("uzupełniająca %d", numer-2)
}

// nazwaBarwyDesignu oddaje nazwę próbki przy trafieniu dokładnym w zapis
// szesnastkowy, a wskaźnik pusty, gdy próbka nie ma nazwy.
func nazwaBarwyDesignu(hex string) *string {
	for nazwa, wzor := range barwyNazwaneDesignu {
		if wzor == hex {
			wartosc := nazwa
			return &wartosc
		}
	}
	return nil
}

// SprawdzKontrast mierzy kontrast pary barw wedle WCAG 2.1 — obsługuje
// `design.color.contrast.check`, uwzględniając rozmiar i grubość pisma
// przy progu.
func (a *adapterDesignu) SprawdzKontrast(_ context.Context,
	z shared.DesignColorContrastCheckRequest) (shared.DesignColorContrastCheckResponse, error) {

	pierwszoplanowa, err := rozpoznajBarweDesignu(z.Foreground)
	if err != nil {
		return shared.DesignColorContrastCheckResponse{}, bladWskazaniaDesignu(
			"komenda design.color.contrast.check, barwa pierwszoplanowa: " + err.Error())
	}
	tlo, err := rozpoznajBarweDesignu(z.Background)
	if err != nil {
		return shared.DesignColorContrastCheckResponse{}, bladWskazaniaDesignu(
			"komenda design.color.contrast.check, barwa tła: " + err.Error())
	}
	if z.FontSize != nil && *z.FontSize < 0 {
		return shared.DesignColorContrastCheckResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.color.contrast.check z rozmiarem pisma %v: rozmiar ujemny nie istnieje",
			*z.FontSize))
	}
	return shared.DesignColorContrastCheckResponse{
		Result: wynikKontrastuDesignu(z.Foreground, z.Background, pierwszoplanowa, tlo,
			z.FontSize, z.Bold),
	}, nil
}

// PrzeliczBarwe oddaje barwę we wszystkich przestrzeniach naraz — obsługuje
// `design.color.convert`; przestrzeń wejściową rdzeń rozpoznaje z zapisu,
// a wskazanie `space` zawęża to rozpoznanie.
func (a *adapterDesignu) PrzeliczBarwe(_ context.Context,
	z shared.DesignColorConvertRequest) (shared.DesignColorConvertResponse, error) {

	if strings.TrimSpace(z.Value) == "" {
		return shared.DesignColorConvertResponse{}, bladWskazaniaDesignu(
			"komenda design.color.convert bez zapisu barwy")
	}
	zapis := strings.TrimSpace(z.Value)
	if z.Space != nil {
		if err := sprawdzWyliczenieDesignu("design.color.convert", "space", *z.Space,
			shared.WartosciDesignColorSpace()); err != nil {
			return shared.DesignColorConvertResponse{}, err
		}
		// Zapis bez nawiasu wraz ze wskazaniem przestrzeni funkcyjnej dostaje
		// nawias tutaj.

		// Operator, który podał zapis bez nawiasu, ma dostać barwę, nie odmowę
		// za brak nawiasu.
		if *z.Space != shared.DesignColorSpaceHex && !strings.Contains(zapis, "(") {
			zapis = string(*z.Space) + "(" + zapis + ")"
		}
	}
	barwa, err := rozpoznajBarweDesignu(zapis)
	if err != nil {
		return shared.DesignColorConvertResponse{}, bladWskazaniaDesignu(
			"komenda design.color.convert: " + err.Error())
	}
	return shared.DesignColorConvertResponse{Color: barwaKontraktuDesignu(barwa)}, nil
}

// UstawGradient zakłada gradient na wskazanym celu albo zmienia zastany —
// obsługuje `design.color.gradient.set`; gradient rozstrzyga się celem, nie
// identyfikatorem, więc powtórne wywołanie na tym samym celu zmienia
// gradient zamiast dokładać drugi.
func (a *adapterDesignu) UstawGradient(ctx context.Context,
	z shared.DesignColorGradientSetRequest) (shared.DesignColorGradientSetResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignColorGradientSetResponse{}, bladWskazaniaDesignu(
			"komenda design.color.gradient.set bez wskazania kompozycji")
	}
	if err := sprawdzWyliczenieDesignu("design.color.gradient.set", "gradient.kind",
		z.Gradient.Kind, shared.WartosciDesignGradientKind()); err != nil {
		return shared.DesignColorGradientSetResponse{}, err
	}
	if len(z.Gradient.Stops) < 2 {
		return shared.DesignColorGradientSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.color.gradient.set z %d stopniami: gradient o mniej niż dwóch "+
				"stopniach jest barwą jednolitą, a nie przejściem", len(z.Gradient.Stops)))
	}
	sciezkaCelu, warstwaCelu := "", ""
	if z.PathId != nil {
		sciezkaCelu = strings.TrimSpace(*z.PathId)
	}
	if z.LayerId != nil {
		warstwaCelu = strings.TrimSpace(*z.LayerId)
	}

	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignColorGradientSetResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	// Cel wskazany a nieznany jest odmową: gradient bez kształtu nie ma czego
	// wypełniać.
	if sciezkaCelu != "" {
		wiersz, err := a.repozytorium.SciezkaWektorowaDesignuPoKodzie(ctx, sciezkaCelu)
		if err != nil {
			return shared.DesignColorGradientSetResponse{}, bladNieznanejSciezkiDesignu(sciezkaCelu, err)
		}
		if wiersz.KompozycjaID != kompozycja.ID {
			return shared.DesignColorGradientSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ścieżka %s leży na innej kompozycji niż %s", sciezkaCelu, kompozycja.Kod))
		}
	}
	if warstwaCelu != "" {
		wiersz, err := a.repozytorium.WarstwaKompozycjiDesignuPoKodzie(ctx, warstwaCelu)
		if err != nil {
			return shared.DesignColorGradientSetResponse{}, bladNieznanejWarstwyDesignu(warstwaCelu, err)
		}
		if wiersz.KompozycjaID != kompozycja.ID {
			return shared.DesignColorGradientSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"warstwa %s leży na innej kompozycji niż %s", warstwaCelu, kompozycja.Kod))
		}
	}

	stopnie := make([]dane.StopienGradientuDesignu, 0, len(z.Gradient.Stops))
	for numer, stopien := range z.Gradient.Stops {
		barwa, err := rozpoznajBarweDesignu(stopien.Color)
		if err != nil {
			return shared.DesignColorGradientSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"stopień numer %d gradientu: %s", numer+1, err.Error()))
		}
		if stopien.Offset < 0 || stopien.Offset > 1 {
			return shared.DesignColorGradientSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"stopień numer %d gradientu w położeniu %v: położenie stopnia idzie od 0 do 1",
				numer+1, stopien.Offset))
		}
		// Barwę zapisujemy sprowadzoną do zapisu szesnastkowego: „red” i
		// „#ff0000” są tą samą barwą.

		// Wykaz mieszający zapisy zmuszałby czytelnika do ponownego rozpoznania
		// barwy.
		stopnie = append(stopnie, dane.StopienGradientuDesignu{
			Polozenie: stopien.Offset, Barwa: barwa.Clamped().Hex(), Krycie: stopien.Opacity,
		})
	}
	// Stopnie idą w kolejności położenia, nie żądania: gradient jest
	// przejściem od lewej do prawej.
	sort.SliceStable(stopnie, func(i, j int) bool {
		return stopnie[i].Polozenie < stopnie[j].Polozenie
	})

	kod := nowyIdentyfikator(przedrostekGradientuDesign)
	if zastany, err := a.repozytorium.GradientDesignuPoCelu(ctx, kompozycja.ID,
		sciezkaCelu, warstwaCelu); err == nil {
		kod = zastany.Kod
	} else if !czyBrakZasobuDesignu(err) {
		return shared.DesignColorGradientSetResponse{}, bladDesignu(err)
	}

	zapisany, err := a.repozytorium.ZapiszGradientDesignu(ctx, dane.GradientDesignu{
		Kod:          kod,
		KompozycjaID: kompozycja.ID,
		SciezkaID:    sciezkaCelu,
		WarstwaID:    warstwaCelu,
		Rodzaj:       string(z.Gradient.Kind),
		Kat:          z.Gradient.Angle,
		Stopnie:      stopnie,
	})
	if err != nil {
		return shared.DesignColorGradientSetResponse{}, bladDesignu(err)
	}

	gradient := gradientKontraktuDesignu(zapisany)
	podglad := podgladGradientuSvgDesignu(gradient)
	return shared.DesignColorGradientSetResponse{Gradient: gradient, PreviewSvg: &podglad}, nil
}

// gradientKontraktuDesignu składa `DesignGradient` kontraktu z wiersza bazy,
// przekładając stopnie gradientu do postaci kontraktu.
func gradientKontraktuDesignu(wiersz dane.GradientDesignu) shared.DesignGradient {
	stopnie := make([]shared.DesignGradientStop, 0, len(wiersz.Stopnie))
	for _, stopien := range wiersz.Stopnie {
		stopnie = append(stopnie, shared.DesignGradientStop{
			Offset: stopien.Polozenie, Color: stopien.Barwa, Opacity: stopien.Krycie,
		})
	}
	kod := wiersz.Kod
	return shared.DesignGradient{
		Id: &kod, Kind: shared.DesignGradientKind(wiersz.Rodzaj), Angle: wiersz.Kat,
		Stops: stopnie,
	}
}

// podgladGradientuSvgDesignu składa podgląd gradientu jako dokument SVG:
// podgląd jest prawdziwym gradientem SVG, nie obrazkiem, więc okno pokazuje
// go wprost i operator może go skopiować do arkusza stylów.
func podgladGradientuSvgDesignu(gradient shared.DesignGradient) string {
	var dokument strings.Builder
	dokument.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="240" height="48" ` +
		`viewBox="0 0 240 48">`)
	dokument.WriteString(`<defs>`)
	switch gradient.Kind {
	case shared.DesignGradientKindRadial:
		dokument.WriteString(`<radialGradient id="podglad" cx="0.5" cy="0.5" r="0.5">`)
		dopiszStopnieSvgDesignu(&dokument, gradient.Stops)
		dokument.WriteString(`</radialGradient>`)
	case shared.DesignGradientKindConic:
		// Gradientu stożkowego SVG 1.1 nie ma; podgląd idzie wtedy jako liniowy.

		// Tytuł mówi to wprost — inaczej podgląd pokazywałby coś innego, niż
		// zapisano.
		dokument.WriteString(`<linearGradient id="podglad" x1="0" y1="0" x2="1" y2="0">`)
		dopiszStopnieSvgDesignu(&dokument, gradient.Stops)
		dokument.WriteString(`</linearGradient>`)
	default:
		x1, y1, x2, y2 := punktyGradientuDesignu(gradient.Angle)
		fmt.Fprintf(&dokument, `<linearGradient id="podglad" x1="%g" y1="%g" x2="%g" y2="%g">`,
			x1, y1, x2, y2)
		dopiszStopnieSvgDesignu(&dokument, gradient.Stops)
		dokument.WriteString(`</linearGradient>`)
	}
	dokument.WriteString(`</defs>`)
	if gradient.Kind == shared.DesignGradientKindConic {
		dokument.WriteString(`<title>gradient stożkowy pokazany jako liniowy — ` +
			`SVG 1.1 nie ma gradientu stożkowego</title>`)
	}
	dokument.WriteString(`<rect width="240" height="48" fill="url(#podglad)"/></svg>`)
	return dokument.String()
}

// dopiszStopnieSvgDesignu dokłada stopnie gradientu do dokumentu SVG jako
// elementy `<stop>` z barwą i opcjonalnym kryciem.
func dopiszStopnieSvgDesignu(dokument *strings.Builder, stopnie []shared.DesignGradientStop) {
	for _, stopien := range stopnie {
		fmt.Fprintf(dokument, `<stop offset="%g" stop-color="%s"`, stopien.Offset, stopien.Color)
		if stopien.Opacity != nil {
			fmt.Fprintf(dokument, ` stop-opacity="%g"`, *stopien.Opacity)
		}
		dokument.WriteString(`/>`)
	}
}

// punktyGradientuDesignu przekłada kąt na parę punktów gradientu liniowego SVG.
// Zero stopni znaczy przejście w prawo — tak liczy je CSS i tego Operator
// oczekuje.
func punktyGradientuDesignu(kat *float64) (float64, float64, float64, float64) {
	stopnie := 0.0
	if kat != nil {
		stopnie = znormalizujKatDesignu(*kat)
	}
	radiany := stopnie * math.Pi / 180
	dx, dy := math.Cos(radiany), math.Sin(radiany)
	return 0.5 - dx/2, 0.5 - dy/2, 0.5 + dx/2, 0.5 + dy/2
}

// ZbadajDostepnoscBarw mierzy kontrast wszystkich par barwnych żetonów zestawu —
// obsługuje `design.color.accessibility.audit`.
func (a *adapterDesignu) ZbadajDostepnoscBarw(ctx context.Context,
	z shared.DesignColorAccessibilityAuditRequest) (shared.DesignColorAccessibilityAuditResponse, error) {

	if strings.TrimSpace(z.TokenSetId) == "" {
		return shared.DesignColorAccessibilityAuditResponse{}, bladWskazaniaDesignu(
			"komenda design.color.accessibility.audit bez wskazania zestawu żetonów")
	}
	prog := progKontrastuAA
	if z.MinimumRatio != nil {
		if *z.MinimumRatio < 1 || *z.MinimumRatio > 21 {
			return shared.DesignColorAccessibilityAuditResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.color.accessibility.audit z progiem %v: skala kontrastu WCAG idzie "+
					"od 1 (barwa na samej sobie) do 21 (czerń na bieli)", *z.MinimumRatio))
		}
		prog = *z.MinimumRatio
	}

	zestaw, err := a.repozytorium.ZestawZetonowDesignuPoKodzie(ctx, strings.TrimSpace(z.TokenSetId))
	if err != nil {
		if czyBrakZasobuDesignu(err) {
			return shared.DesignColorAccessibilityAuditResponse{}, bladNieznanegoBytuDesignu(
				"zestawu żetonów " + z.TokenSetId + " nie ma w tym rdzeniu")
		}
		return shared.DesignColorAccessibilityAuditResponse{}, bladDesignu(err)
	}
	zetony, err := a.repozytorium.ZetonyZestawuDesignu(ctx, zestaw.ID)
	if err != nil {
		return shared.DesignColorAccessibilityAuditResponse{}, bladDesignu(err)
	}

	// Żeton, którego wartość nie jest barwą, nie wchodzi do pomiaru — zestaw
	// niesie też wymiary i kroje.

	// Żeton rodzaju „color” o wartości nieczytelnej jest natomiast odmową,
	// nie pominięciem po cichu.
	barwy := make([]shared.DesignPaletteColor, 0, len(zetony))
	for _, zeton := range zetony {
		if zeton.Rodzaj != shared.DesignTokenKindColor {
			continue
		}
		barwa, err := rozpoznajBarweDesignu(zeton.Wartosc)
		if err != nil {
			return shared.DesignColorAccessibilityAuditResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"żeton %s zestawu %s ma rodzaj „color\", ale jego wartość %q nie jest barwą: %s",
				zeton.Nazwa, zestaw.Kod, zeton.Wartosc, err.Error()))
		}
		nazwa := zeton.Nazwa
		barwy = append(barwy, shared.DesignPaletteColor{
			Hex: barwa.Clamped().Hex(), Name: &nazwa,
		})
	}
	if len(barwy) < 2 {
		return shared.DesignColorAccessibilityAuditResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"zestaw %s ma %d żetonów barwnych — badanie dostępności mierzy PARY, więc przy mniej "+
				"niż dwóch barwach nie ma czego mierzyć", zestaw.Kod, len(barwy)))
	}

	sprawdzonych, spelniajacych := 0, 0
	lamiace := []shared.DesignContrastResult{}
	for pierwszy := 0; pierwszy < len(barwy); pierwszy++ {
		for drugi := pierwszy + 1; drugi < len(barwy); drugi++ {
			pierwszaBarwa, _ := rozpoznajBarweDesignu(barwy[pierwszy].Hex)
			drugaBarwa, _ := rozpoznajBarweDesignu(barwy[drugi].Hex)
			wynik := wynikKontrastuDesignu(barwy[pierwszy].Hex, barwy[drugi].Hex,
				pierwszaBarwa, drugaBarwa, nil, nil)
			sprawdzonych++
			if wynik.Ratio >= prog {
				spelniajacych++
				continue
			}
			lamiace = append(lamiace, wynik)
		}
	}
	// Pary łamiące próg idą od najgorszej, żeby operator naprawiał najpierw tę
	// najdalszą od progu.
	sort.SliceStable(lamiace, func(i, j int) bool { return lamiace[i].Ratio < lamiace[j].Ratio })
	return shared.DesignColorAccessibilityAuditResponse{
		Violations: lamiace, Checked: sprawdzonych, Passed: spelniajacych,
	}, nil
}
