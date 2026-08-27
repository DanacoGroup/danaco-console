// Plik obsługuje osiem czynności ikon i krojów modułu design — `design.icon.*`,
// `design.favicon.build` i `design.font.*`. Katalog wzorów leży
// w `adapter_modul_design_ikony_katalog.go`, kroje w `adapter_modul_design_kroje.go`.
package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/tdewolff/canvas"
	"golang.org/x/image/font/sfnt"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// przedrostekIkonyDesign znakuje identyfikatory zewnętrzne ikon własnych, obok
	// przedrostków innych bytów warsztatu.
	przedrostekIkonyDesign = "ikona-"

	// domyslnyLimitWykazuIkonDesignu jest górną granicą wykazu ikon, gdy żądanie
	// granicy nie podało jej wprost.
	domyslnyLimitWykazuIkonDesignu = 50

	// granicaOsadzeniaKrojuDesignu jest wielkością kroju, powyżej której podgląd
	// nie osadza go w regule `@font-face`. Krój pełny CJK ma kilkanaście
	// megabajtów, a odpowiedź komendy nie jest miejscem na tyle bajtów.
	granicaOsadzeniaKrojuDesignu = 512 << 10
)

// rozmiaryFaviconyDesignu to komplet rozmiarów przyjęty dla ikony witryny.
// Wykaz jest tym, co naprawdę czytają systemy i przeglądarki: 16 i 32 dla karty,
// 48 dla pulpitu Windows, 180 dla ekranu iOS, 192 i 512 dla manifestu aplikacji
// sieciowej.
var rozmiaryFaviconyDesignu = []int{16, 32, 48, 64, 128, 180, 192, 256, 512}

// SzukajIkon zwraca wzory katalogu wkompilowanego spełniające warunki —
// obsługuje `design.icon.library.search`. Wykaz obejmuje wyłącznie wzory
// rdzenia; żądanie nie niesie okna, więc ikon własnych nie da się zawęzić
// do właściwego okna.
func (a *adapterDesignu) SzukajIkon(_ context.Context,
	z shared.DesignIconLibrarySearchRequest) (shared.DesignIconLibrarySearchResponse, error) {

	limit := domyslnyLimitWykazuIkonDesignu
	if z.Limit != nil {
		if *z.Limit < 1 {
			return shared.DesignIconLibrarySearchResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.icon.library.search z granicą %d: wykaz o niedodatniej długości "+
					"nie jest wykazem", *z.Limit))
		}
		limit = *z.Limit
	}
	if z.Set != nil && strings.TrimSpace(*z.Set) != "" &&
		!strings.EqualFold(strings.TrimSpace(*z.Set), zestawIkonKataloguDesignu) {

		return shared.DesignIconLibrarySearchResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.icon.library.search z zestawem %q, którego rdzeń nie ma; zestawy "+
				"w katalogu rdzenia: %s", *z.Set, strings.Join(zestawyIkonDesignu(), ", ")))
	}

	fraza := ""
	if z.Query != nil {
		fraza = strings.ToLower(strings.TrimSpace(*z.Query))
	}
	spelniajace := []shared.DesignIcon{}
	for _, wzor := range katalogIkonDesignu {
		if fraza != "" && !pasujeWzorDoFrazyDesignu(wzor, fraza) {
			continue
		}
		spelniajace = append(spelniajace, ikonaKataloguKontraktuDesignu(wzor,
			siatkaIkonyKataloguDesignu, gruboscObrysuIkonyDomyslna))
	}
	// total niesie liczbę spełniających warunki, nie długość strony — po to, żeby podnieść granicę.
	razem := len(spelniajace)
	if len(spelniajace) > limit {
		spelniajace = spelniajace[:limit]
	}
	return shared.DesignIconLibrarySearchResponse{
		Icons: spelniajace, Total: razem, Sets: zestawyIkonDesignu(),
	}, nil
}

// pasujeWzorDoFrazyDesignu rozstrzyga, czy wzór spełnia frazę wyszukiwania —
// po nazwie i po etykietach.
func pasujeWzorDoFrazyDesignu(wzor wzorIkonyDesignu, fraza string) bool {
	if strings.Contains(strings.ToLower(wzor.Nazwa), fraza) {
		return true
	}
	for _, etykieta := range wzor.Etykiety {
		if strings.Contains(strings.ToLower(etykieta), fraza) {
			return true
		}
	}
	return false
}

// UstawIkone zapisuje ikonę własną na wskazanym oknie, wraz z zastrzeżeniami
// siatki — obsługuje `design.icon.set`.
func (a *adapterDesignu) UstawIkone(ctx context.Context,
	z shared.DesignIconSetRequest) (shared.DesignIconSetResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignIconSetResponse{}, bladWskazaniaDesignu(
			"komenda design.icon.set bez wskazania okna")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignIconSetResponse{}, bladWskazaniaDesignu(
			"komenda design.icon.set bez nazwy ikony")
	}
	if strings.TrimSpace(z.Svg) == "" {
		return shared.DesignIconSetResponse{}, bladWskazaniaDesignu(
			"komenda design.icon.set bez treści SVG: ikona bez rysunku nie ma czego pokazać")
	}
	if z.GridSize != nil && *z.GridSize <= 0 {
		return shared.DesignIconSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.icon.set z siatką %d: siatka o niedodatnim boku nie istnieje", *z.GridSize))
	}
	if z.StrokeWidth != nil && *z.StrokeWidth < 0 {
		return shared.DesignIconSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.icon.set z grubością obrysu %v: grubość ujemna nie istnieje",
			*z.StrokeWidth))
	}

	siatka := siatkaIkonyKataloguDesignu
	if z.GridSize != nil {
		siatka = *z.GridSize
	}
	grubosc := gruboscObrysuIkonyDomyslna
	if z.StrokeWidth != nil && *z.StrokeWidth > 0 {
		grubosc = *z.StrokeWidth
	}

	kod := nowyIdentyfikator(przedrostekIkonyDesign)
	if z.IconId != nil && strings.TrimSpace(*z.IconId) != "" {
		kod = strings.TrimSpace(*z.IconId)
		if strings.HasPrefix(kod, zestawIkonKataloguDesignu+"/") {
			return shared.DesignIconSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ikona %s należy do katalogu WKOMPILOWANEGO w binarium rdzenia i nie daje się "+
					"nadpisać — inaczej dwie instalacje produktu miałyby dwa różne katalogi pod tą "+
					"samą nazwą; pomiń pole iconId, żeby założyć własną ikonę z tego wzoru", kod))
		}
		zastana, err := a.repozytorium.IkonaDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignIconSetResponse{}, bladNieznanejIkonyDesignu(kod, err)
		}
		if zastana.Okno != z.WindowId {
			return shared.DesignIconSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ikona %s należy do okna %s, a komenda design.icon.set przyszła z okna %s",
				kod, zastana.Okno, z.WindowId))
		}
	}

	siatkaKolumna := int64(siatka)
	zapisana, err := a.repozytorium.ZapiszIkoneDesignu(ctx, dane.IkonaDesignu{
		Kod:           kod,
		Okno:          z.WindowId,
		Nazwa:         strings.TrimSpace(z.Name),
		Zestaw:        nil,
		SVG:           strings.TrimSpace(z.Svg),
		Siatka:        &siatkaKolumna,
		GruboscObrysu: &grubosc,
		Etykiety:      z.Tags,
	})
	if err != nil {
		return shared.DesignIconSetResponse{}, bladDesignu(err)
	}
	return shared.DesignIconSetResponse{
		Icon:         ikonaWlasnaKontraktuDesignu(zapisana),
		GridWarnings: uporzadkujBilansDesignu(zastrzezeniaSiatkiIkonyDesignu(z.Svg, siatka, grubosc)),
	}, nil
}

// zastrzezeniaSiatkiIkonyDesignu wylicza miejsca, w których ikona nie trzyma
// siatki. To nie jest odmowa: ikona spoza siatki nadal jest ikoną, którą wolno
// zapisać. Jest bilansem — rozjazd siatki i grubości widać dopiero wtedy,
// gdy się go wypisze.
func zastrzezeniaSiatkiIkonyDesignu(dokument string, siatka int, grubosc float64) []string {
	zastrzezenia := []string{}
	if !strings.Contains(dokument, "viewBox") {
		zastrzezenia = append(zastrzezenia,
			"ikona nie ma pola widoku (viewBox) — bez niego nie da się jej przeskalować bez "+
				"utraty proporcji")
	} else if pole := polePoWidokuDesignu(dokument); pole != "" {
		oczekiwane := fmt.Sprintf("0 0 %d %d", siatka, siatka)
		if strings.Join(strings.Fields(pole), " ") != oczekiwane {
			zastrzezenia = append(zastrzezenia, fmt.Sprintf(
				"pole widoku ikony to %q, a siatka zadeklarowana to %d — ikona wyjdzie w innym "+
					"rozmiarze niż reszta zestawu", pole, siatka))
		}
	}
	if zadeklarowana := atrybutDokumentuDesignu(dokument, "stroke-width"); zadeklarowana != "" {
		if zadeklarowana != fmt.Sprintf("%g", grubosc) {
			zastrzezenia = append(zastrzezenia, fmt.Sprintf(
				"treść ikony niesie grubość obrysu %q, a pole strokeWidth mówi %g — w zestawie "+
					"obowiązuje jedna grubość", zadeklarowana, grubosc))
		}
	}
	sciezki := sciezkiZDokumentuSvgDesignu(dokument)
	if len(sciezki) == 0 {
		zastrzezenia = append(zastrzezenia,
			"ikona nie niesie ani jednej ścieżki (atrybutu d) — do pakietu i do kroju "+
				"ikonowego wejdzie jako glif pusty")
	}
	for numer, zapis := range sciezki {
		if poza := wspolrzednePozaSiatkaDesignu(zapis, siatka); poza > 0 {
			zastrzezenia = append(zastrzezenia, fmt.Sprintf(
				"ścieżka numer %d wychodzi poza siatkę %d (współrzędnych poza polem: %d)",
				numer+1, siatka, poza))
		}
	}
	return zastrzezenia
}

// polePoWidokuDesignu wyciąga treść atrybutu viewBox dokumentu SVG ikony, bez
// przetwarzania jej wartości.
func polePoWidokuDesignu(dokument string) string {
	return atrybutDokumentuDesignu(dokument, "viewBox")
}

// atrybutDokumentuDesignu wyciąga treść pierwszego wystąpienia atrybutu
// o wskazanej nazwie w dokumencie.
func atrybutDokumentuDesignu(dokument, nazwa string) string {
	znacznik := nazwa + `="`
	poczatek := strings.Index(dokument, znacznik)
	if poczatek < 0 {
		return ""
	}
	reszta := dokument[poczatek+len(znacznik):]
	koniec := strings.IndexByte(reszta, '"')
	if koniec < 0 {
		return ""
	}
	return strings.TrimSpace(reszta[:koniec])
}

// wspolrzednePozaSiatkaDesignu liczy współrzędne wychodzące poza siatkę,
// rachunkiem przez bibliotekę ścieżek, nie szukaniem liczb w napisie: zapis SVG
// ma współrzędne względne, skróty i łuki, a prostokąt otaczający mówi prawdę
// o każdej postaci zapisu.
func wspolrzednePozaSiatkaDesignu(zapis string, siatka int) int {
	sciezka, err := canvas.ParseSVGPath(zapis)
	if err != nil {
		// Ścieżki nieczytelnej nie liczy się jako wychodzącej poza siatkę — to inne zastrzeżenie.
		return 0
	}
	granice := sciezka.Bounds()
	poza := 0
	if granice.X0 < -0.01 {
		poza++
	}
	if granice.Y0 < -0.01 {
		poza++
	}
	if granice.X1 > float64(siatka)+0.01 {
		poza++
	}
	if granice.Y1 > float64(siatka)+0.01 {
		poza++
	}
	return poza
}

// GenerujIkony zakłada ikony dla wskazanych pojęć — obsługuje
// `design.icon.generate`. Pojęcie trafiające w katalog wkompilowany zakłada
// ikonę własną z tego wzoru, pojęcie nietrafiające idzie do kanału modelu.
func (a *adapterDesignu) GenerujIkony(ctx context.Context,
	z shared.DesignIconGenerateRequest) (shared.DesignIconGenerateResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignIconGenerateResponse{}, bladWskazaniaDesignu(
			"komenda design.icon.generate bez wskazania okna")
	}
	if len(z.Concepts) == 0 {
		return shared.DesignIconGenerateResponse{}, bladWskazaniaDesignu(
			"komenda design.icon.generate bez ani jednego pojęcia: rdzeń nie wymyśla, jakie ikony " +
				"Operator potrzebuje")
	}
	siatka := siatkaIkonyKataloguDesignu
	if z.GridSize != nil && *z.GridSize > 0 {
		siatka = *z.GridSize
	}
	grubosc := gruboscObrysuIkonyDomyslna
	if z.StrokeWidth != nil && *z.StrokeWidth > 0 {
		grubosc = *z.StrokeWidth
	}
	// Ikona wzorcowa narzuca styl: jej siatka i grubość biją nastawy żądania, po to jest wskazana.
	if z.StyleReferenceIconId != nil && strings.TrimSpace(*z.StyleReferenceIconId) != "" {
		wzorcowa, err := a.ikonaPoKodzieDesignu(ctx, strings.TrimSpace(*z.StyleReferenceIconId))
		if err != nil {
			return shared.DesignIconGenerateResponse{}, err
		}
		if wzorcowa.GridSize != nil && *wzorcowa.GridSize > 0 {
			siatka = *wzorcowa.GridSize
		}
		if wzorcowa.StrokeWidth != nil && *wzorcowa.StrokeWidth > 0 {
			grubosc = *wzorcowa.StrokeWidth
		}
	}

	ikony := []shared.DesignIcon{}
	nieudane := []string{}
	for _, pojecie := range z.Concepts {
		pojecie = strings.TrimSpace(pojecie)
		if pojecie == "" {
			continue
		}
		dokument := ""
		if wzor, jest := wzorDlaPojeciaDesignu(pojecie); jest {
			dokument = svgIkonyKataloguDesignu(wzor, siatka, grubosc)
		} else if z.ChannelId != nil && strings.TrimSpace(*z.ChannelId) != "" {
			polecenie := fmt.Sprintf(
				"Narysuj ikonę pojęcia %q jako dokument SVG na siatce %d×%d, konturem o grubości "+
					"%g, bez wypełnienia, bez tekstu i bez komentarza. Oddaj sam dokument SVG.",
				pojecie, siatka, siatka, grubosc)
			if odpowiedz, err := a.zapytajModelDesignu(ctx, z.WindowId,
				strings.TrimSpace(*z.ChannelId), polecenie); err == nil {
				dokument = wyciagnijDokumentSvgDesignu(odpowiedz)
			}
		}
		if dokument == "" {
			nieudane = append(nieudane, pojecie)
			continue
		}
		siatkaKolumna := int64(siatka)
		zapisana, err := a.repozytorium.ZapiszIkoneDesignu(ctx, dane.IkonaDesignu{
			Kod:           nowyIdentyfikator(przedrostekIkonyDesign),
			Okno:          z.WindowId,
			Nazwa:         pojecie,
			SVG:           dokument,
			Siatka:        &siatkaKolumna,
			GruboscObrysu: &grubosc,
			Etykiety:      []string{pojecie},
		})
		if err != nil {
			return shared.DesignIconGenerateResponse{}, bladDesignu(err)
		}
		ikony = append(ikony, ikonaWlasnaKontraktuDesignu(zapisana))
	}

	if len(ikony) == 0 {
		return shared.DesignIconGenerateResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"żadne z pojęć (%s) nie trafiło w katalog wkompilowany rdzenia, a kanał modelu nie "+
				"był wskazany albo nie oddał dokumentu SVG — rdzeń nie zakłada ikon pustych; "+
				"naprawa: użyć nazw z katalogu (design.icon.library.search) albo wskazać "+
				"channelId kanału tekstowego", strings.Join(nieudane, ", ")))
	}
	return shared.DesignIconGenerateResponse{
		Icons: ikony, FailedConcepts: uporzadkujBilansDesignu(nieudane),
	}, nil
}

// wyciagnijDokumentSvgDesignu wyjmuje dokument SVG z odpowiedzi modelu. Model
// bywa rozmowny: odpowiedź niesie zdanie wstępne, blok kodu i podsumowanie.
// Brany jest zakres od `<svg` do `</svg>` i tylko go.
func wyciagnijDokumentSvgDesignu(odpowiedz string) string {
	poczatek := strings.Index(odpowiedz, "<svg")
	if poczatek < 0 {
		return ""
	}
	koniec := strings.LastIndex(odpowiedz, "</svg>")
	if koniec < poczatek {
		return ""
	}
	return strings.TrimSpace(odpowiedz[poczatek : koniec+len("</svg>")])
}

// ZbudujPakietIkon składa pakiet ikon jako dokument SVG albo krój webfont —
// obsługuje `design.icon.sprite.build`.
func (a *adapterDesignu) ZbudujPakietIkon(ctx context.Context,
	z shared.DesignIconSpriteBuildRequest) (shared.DesignIconSpriteBuildResponse, error) {

	if len(z.IconIds) == 0 {
		return shared.DesignIconSpriteBuildResponse{}, bladWskazaniaDesignu(
			"komenda design.icon.sprite.build bez ani jednej ikony: pakiet pusty nie ma czego nieść")
	}
	if err := sprawdzWyliczenieDesignu("design.icon.sprite.build", "kind", z.Kind,
		shared.WartosciDesignIconSpriteKind()); err != nil {
		return shared.DesignIconSpriteBuildResponse{}, err
	}

	nazwa := "ikony"
	if z.Name != nil && strings.TrimSpace(*z.Name) != "" {
		nazwa = strings.TrimSpace(*z.Name)
	}
	okno := ""
	ikony := make([]shared.DesignIcon, 0, len(z.IconIds))
	for _, kod := range z.IconIds {
		ikona, err := a.ikonaPoKodzieDesignu(ctx, strings.TrimSpace(kod))
		if err != nil {
			return shared.DesignIconSpriteBuildResponse{}, err
		}
		ikony = append(ikony, ikona)
		// Okno wytworu bierze okno pierwszej ikony WŁASNEJ — wzór rdzenia okna nie
		// ma i nie da rady go podać.
		if okno == "" {
			if wiersz, err := a.repozytorium.IkonaDesignuPoKodzie(ctx, strings.TrimSpace(kod)); err == nil {
				okno = wiersz.Okno
			}
		}
	}
	okno = oknoWytworu(z.WindowId, okno)
	if strings.TrimSpace(okno) == "" {
		return shared.DesignIconSpriteBuildResponse{}, bladWskazaniaDesignu(
			"komenda design.icon.sprite.build bez wskazania okna, a wszystkie ikony pakietu " +
				"pochodzą z katalogu wkompilowanego, który okna nie ma — wskaż windowId, żeby " +
				"pakiet miał gdzie stanąć")
	}

	if z.Kind == shared.DesignIconSpriteKindWebfont {
		doKroju := make([]ikonaDoKrojuDesignu, 0, len(ikony))
		for _, ikona := range ikony {
			dokument := ""
			if ikona.Svg != nil {
				dokument = *ikona.Svg
			}
			wpis := ikonaDoKrojuDesignu{
				Nazwa:   ikona.Name,
				Sciezki: sciezkiZDokumentuSvgDesignu(dokument),
				Grubosc: gruboscObrysuIkonyDomyslna,
			}
			if ikona.StrokeWidth != nil && *ikona.StrokeWidth > 0 {
				wpis.Grubosc = *ikona.StrokeWidth
			}
			if ikona.GridSize != nil && *ikona.GridSize > 0 {
				wpis.Siatka = *ikona.GridSize
			}
			doKroju = append(doKroju, wpis)
		}
		bajty, glifow, err := zlozKrojIkonowyDesignu(nazwa, doKroju)
		if err != nil {
			return shared.DesignIconSpriteBuildResponse{}, bladWskazaniaDesignu(
				"komenda design.icon.sprite.build w postaci webfont: " + err.Error())
		}
		plik := oczyscNazwePlikuDesignu(nazwa) + ".ttf"
		zasob, err := a.zalozZasobZBajtowDesignu(ctx, okno, plik,
			shared.DesignAssetKindDocument, "ttf", bajty)
		if err != nil {
			return shared.DesignIconSpriteBuildResponse{}, err
		}
		return shared.DesignIconSpriteBuildResponse{
			Asset: zasobWytworzonyKontraktu(zasob), FileName: plik, Included: glifow,
		}, nil
	}

	tresc, wlaczonych := zlozPakietSvgIkonDesignu(ikony)
	if wlaczonych == 0 {
		return shared.DesignIconSpriteBuildResponse{}, bladWskazaniaDesignu(
			"żadna ze wskazanych ikon nie niesie ścieżki (atrybutu d) — pakiet SVG wyszedłby " +
				"z pustymi symbolami, czyli plikiem, w którym nie widać nic")
	}
	plik := oczyscNazwePlikuDesignu(nazwa) + ".svg"
	zasob, err := a.zalozZasobZBajtowDesignu(ctx, okno, plik,
		shared.DesignAssetKindVector, "svg", []byte(tresc))
	if err != nil {
		return shared.DesignIconSpriteBuildResponse{}, err
	}
	return shared.DesignIconSpriteBuildResponse{
		Asset: zasobWytworzonyKontraktu(zasob), FileName: plik, Included: wlaczonych,
	}, nil
}

// zlozPakietSvgIkonDesignu składa pakiet ikon jako jeden dokument SVG
// z symbolami — postać, którą strona wciąga raz i używa przez `<use>`.
func zlozPakietSvgIkonDesignu(ikony []shared.DesignIcon) (string, int) {
	var dokument strings.Builder
	dokument.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	dokument.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" style="display:none">` + "\n")
	wlaczonych := 0
	for _, ikona := range ikony {
		dokumentIkony := ""
		if ikona.Svg != nil {
			dokumentIkony = *ikona.Svg
		}
		sciezki := sciezkiZDokumentuSvgDesignu(dokumentIkony)
		if len(sciezki) == 0 {
			continue
		}
		siatka := siatkaIkonyKataloguDesignu
		if ikona.GridSize != nil && *ikona.GridSize > 0 {
			siatka = *ikona.GridSize
		}
		grubosc := gruboscObrysuIkonyDomyslna
		if ikona.StrokeWidth != nil && *ikona.StrokeWidth > 0 {
			grubosc = *ikona.StrokeWidth
		}
		fmt.Fprintf(&dokument,
			`  <symbol id="%s" viewBox="0 0 %d %d" fill="none" stroke="currentColor" `+
				`stroke-width="%g" stroke-linecap="round" stroke-linejoin="round">`+"\n",
			oczyscNazwePlikuDesignu(ikona.Name), siatka, siatka, grubosc)
		for _, sciezka := range sciezki {
			fmt.Fprintf(&dokument, `    <path d="%s"/>`+"\n", sciezka)
		}
		dokument.WriteString("  </symbol>\n")
		wlaczonych++
	}
	dokument.WriteString("</svg>\n")
	return dokument.String(), wlaczonych
}

// ZbudujFavicone wydaje komplet ikon witryny z zasobu — obsługuje
// `design.favicon.build`. Rozmiar powyżej 256 wychodzi jako PNG, nie jako ICO:
// pole boku w katalogu ikony ma jeden bajt. Wykaz `sizes` odpowiedzi niesie
// rozmiary, które naprawdę powstały.
func (a *adapterDesignu) ZbudujFavicone(ctx context.Context,
	z shared.DesignFaviconBuildRequest) (shared.DesignFaviconBuildResponse, error) {

	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.favicon.build", z.AssetId)
	if err != nil {
		return shared.DesignFaviconBuildResponse{}, err
	}
	rozmiary := rozmiaryFaviconyDesignu
	if len(z.Sizes) > 0 {
		rozmiary = make([]int, 0, len(z.Sizes))
		for _, rozmiar := range z.Sizes {
			if rozmiar < 1 || rozmiar > 1024 {
				return shared.DesignFaviconBuildResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
					"komenda design.favicon.build z rozmiarem %d: rdzeń wydaje ikony witryny "+
						"o boku od 1 do 1024", rozmiar))
			}
			rozmiary = append(rozmiary, rozmiar)
		}
		sort.Ints(rozmiary)
	}
	okno := oknoWytworu(z.WindowId, zrodlo.Okno)

	kody := []string{}
	powstale := []int{}
	granice := obraz.Bounds()
	for _, rozmiar := range rozmiary {
		skala := float64(rozmiar) / float64(granice.Dx())
		if granice.Dy() > granice.Dx() {
			skala = float64(rozmiar) / float64(granice.Dy())
		}
		zmniejszony := przeskalujObrazDesignu(obraz, skala)
		format := "ico"
		if rozmiar > 256 {
			format = "png"
		}
		bajty, _, err := zakodujObrazDesignu(zmniejszony, format, nil)
		if err != nil {
			return shared.DesignFaviconBuildResponse{}, bladWydaniaDesignu(err.Error())
		}
		nazwa := fmt.Sprintf("%s-%dx%d.%s", nazwaZasobuDesignu(zrodlo), rozmiar, rozmiar, format)
		zasob, err := a.zalozZasobZBajtowDesignu(ctx, okno, nazwa,
			shared.DesignAssetKindImage, format, bajty)
		if err != nil {
			return shared.DesignFaviconBuildResponse{}, err
		}
		kody = append(kody, zasob.Kod)
		powstale = append(powstale, rozmiar)
	}
	if len(kody) == 0 {
		return shared.DesignFaviconBuildResponse{}, bladWydaniaDesignu(
			"nie powstała ani jedna ikona witryny")
	}

	odpowiedz := shared.DesignFaviconBuildResponse{AssetIds: kody, Sizes: powstale}
	if z.IncludeManifest != nil && *z.IncludeManifest {
		manifest, err := manifestAplikacjiDesignu(nazwaZasobuDesignu(zrodlo), kody, powstale)
		if err != nil {
			return shared.DesignFaviconBuildResponse{}, bladWydaniaDesignu(err.Error())
		}
		odpowiedz.Manifest = &manifest
	}
	return odpowiedz, nil
}

// manifestAplikacjiDesignu składa treść manifestu aplikacji sieciowej. Odsyłacze
// wskazują identyfikatory zasobów rdzenia, nie ścieżki na dysku: podstawia się
// je pod adresy witryny, a ścieżka magazynu rdzenia nie otworzy się nigdzie
// poza tą maszyną.
func manifestAplikacjiDesignu(nazwa string, kody []string, rozmiary []int) (string, error) {
	type ikonaManifestu struct {
		Src   string `json:"src"`
		Sizes string `json:"sizes"`
		Type  string `json:"type"`
	}
	ikony := make([]ikonaManifestu, 0, len(kody))
	for numer, kod := range kody {
		typ := "image/vnd.microsoft.icon"
		if rozmiary[numer] > 256 {
			typ = "image/png"
		}
		ikony = append(ikony, ikonaManifestu{
			Src:   kod,
			Sizes: fmt.Sprintf("%dx%d", rozmiary[numer], rozmiary[numer]),
			Type:  typ,
		})
	}
	manifest := map[string]any{
		"name":       nazwa,
		"short_name": nazwa,
		"icons":      ikony,
	}
	bajty, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("nie można złożyć manifestu aplikacji sieciowej: %w", err)
	}
	return string(bajty), nil
}

// ZaproponujZestawieniaKrojow oddaje propozycje par krojów — obsługuje
// `design.font.pair.suggest`. Pole `source` mówi, którą drogą propozycje
// powstały: kanałem modelu albo regułą rdzenia, zestawiającą wyłącznie kroje,
// którymi rdzeń naprawdę dysponuje.
func (a *adapterDesignu) ZaproponujZestawieniaKrojow(ctx context.Context,
	z shared.DesignFontPairSuggestRequest) (shared.DesignFontPairSuggestResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignFontPairSuggestResponse{}, bladWskazaniaDesignu(
			"komenda design.font.pair.suggest bez wskazania okna")
	}
	ile := 3
	if z.Count != nil {
		if *z.Count < 1 || *z.Count > 24 {
			return shared.DesignFontPairSuggestResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.font.pair.suggest z liczbą propozycji %d: rdzeń liczy od 1 do 24",
				*z.Count))
		}
		ile = *z.Count
	}

	if z.ChannelId != nil && strings.TrimSpace(*z.ChannelId) != "" {
		polecenie := zlozPolecenieZestawienKrojowDesignu(z, ile)
		if odpowiedz, err := a.zapytajModelDesignu(ctx, z.WindowId,
			strings.TrimSpace(*z.ChannelId), polecenie); err == nil {
			if pary := paryZOdpowiedziModeluDesignu(odpowiedz, ile); len(pary) > 0 {
				return shared.DesignFontPairSuggestResponse{Pairs: pary, Source: "kanał modelu"}, nil
			}
		}
	}

	pary := zestawieniaKrojowRdzeniaDesignu(z, ile)
	if len(pary) == 0 {
		return shared.DesignFontPairSuggestResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"rdzeń ma za mało krojów, żeby złożyć choć jedno zestawienie (krojów wkompilowanych: "+
				"%d, krojów tej maszyny: %d)",
			len(krojeWkompilowaneDesignu), len(przegladajKrojeSerweraDesignu())))
	}
	return shared.DesignFontPairSuggestResponse{Pairs: pary, Source: "reguła rdzenia"}, nil
}

// zlozPolecenieZestawienKrojowDesignu składa polecenie dla kanału modelu. Wykaz
// krojów, jakimi rdzeń dysponuje, wchodzi w treść polecenia — inaczej model
// proponowałby kroje, których u Operatora nie ma.
func zlozPolecenieZestawienKrojowDesignu(z shared.DesignFontPairSuggestRequest, ile int) string {
	var polecenie strings.Builder
	fmt.Fprintf(&polecenie, "Zaproponuj %d zestawień krojów pisma. ", ile)
	if z.Mood != nil && strings.TrimSpace(*z.Mood) != "" {
		fmt.Fprintf(&polecenie, "Charakter: %s. ", strings.TrimSpace(*z.Mood))
	}
	if z.BaseFont != nil && strings.TrimSpace(*z.BaseFont) != "" {
		fmt.Fprintf(&polecenie, "Krój, do którego szukamy pary: %s. ", strings.TrimSpace(*z.BaseFont))
	}
	nazwy := nazwyKrojowDesignu()
	if len(nazwy) > 40 {
		nazwy = nazwy[:40]
	}
	fmt.Fprintf(&polecenie,
		"Wybieraj WYŁĄCZNIE z krojów, jakimi rdzeń dysponuje: %s. "+
			"Każde zestawienie w osobnym wierszu w postaci: nagłówek | tekst | powód.",
		strings.Join(nazwy, ", "))
	return polecenie.String()
}

// paryZOdpowiedziModeluDesignu rozkłada odpowiedź modelu na pary. Wiersz
// wskazujący krój, którego rdzeń nie ma, jest POMIJANY — propozycja nie do
// zobaczenia nie jest propozycją.
func paryZOdpowiedziModeluDesignu(odpowiedz string, ile int) []shared.DesignFontPair {
	pary := []shared.DesignFontPair{}
	for _, wiersz := range strings.Split(odpowiedz, "\n") {
		czesci := strings.Split(wiersz, "|")
		if len(czesci) < 2 {
			continue
		}
		naglowek := strings.TrimSpace(czesci[0])
		tekst := strings.TrimSpace(czesci[1])
		if naglowek == "" || tekst == "" {
			continue
		}
		if _, _, _, err := krojDesignu(naglowek); err != nil {
			continue
		}
		if _, _, _, err := krojDesignu(tekst); err != nil {
			continue
		}
		para := shared.DesignFontPair{Heading: naglowek, Body: tekst}
		if len(czesci) > 2 && strings.TrimSpace(czesci[2]) != "" {
			powod := strings.TrimSpace(czesci[2])
			para.Reason = &powod
		}
		probka := domyslnyTekstProbnyDesignu
		para.SampleText = &probka
		pary = append(pary, para)
		if len(pary) >= ile {
			break
		}
	}
	return pary
}

// zestawieniaKrojowRdzeniaDesignu składa zestawienia regułą rdzenia. Reguła
// jest typograficzna, nie losowa: nagłówek dostaje krój o większej masie albo
// szeryfowy, tekst — o mniejszej masie albo bezszeryfowy, i oba mają być
// krojami różnymi.
func zestawieniaKrojowRdzeniaDesignu(z shared.DesignFontPairSuggestRequest,
	ile int) []shared.DesignFontPair {

	dostepne := nazwyKrojowDesignu()
	if len(dostepne) < 2 {
		return nil
	}
	naglowkowe := []string{}
	tekstowe := []string{}
	for _, nazwa := range dostepne {
		maly := strings.ToLower(nazwa)
		switch {
		case strings.Contains(maly, "bold") || strings.Contains(maly, "demi") ||
			strings.Contains(maly, "medium") || strings.Contains(maly, "black"):
			naglowkowe = append(naglowkowe, nazwa)
		case strings.Contains(maly, "italic") || strings.Contains(maly, "oblique"):
			// Odmiana pochyła nie jest krojem tekstu ciągłego ani nagłówka — pomija się ją.
		default:
			tekstowe = append(tekstowe, nazwa)
		}
	}
	if len(naglowkowe) == 0 {
		naglowkowe = tekstowe
	}
	if len(tekstowe) == 0 {
		tekstowe = naglowkowe
	}

	// Krój wskazany w żądaniu wchodzi jako tekstowy i jest w każdej parze zestawienia.
	podstawowy := ""
	if z.BaseFont != nil && strings.TrimSpace(*z.BaseFont) != "" {
		if _, nazwa, _, err := krojDesignu(*z.BaseFont); err == nil {
			podstawowy = nazwa
		}
	}

	probka := domyslnyTekstProbnyDesignu
	pary := []shared.DesignFontPair{}
	for _, naglowek := range naglowkowe {
		for _, tekst := range tekstowe {
			if podstawowy != "" {
				tekst = podstawowy
			}
			if strings.EqualFold(naglowek, tekst) {
				continue
			}
			powod := zlozPowodZestawieniaDesignu(naglowek, tekst, z.Mood)
			pary = append(pary, shared.DesignFontPair{
				Heading: naglowek, Body: tekst, Reason: &powod, SampleText: &probka,
			})
			if len(pary) >= ile {
				return pary
			}
			if podstawowy != "" {
				break
			}
		}
	}
	return pary
}

// zlozPowodZestawieniaDesignu nazywa powód zestawienia. Powód jest opisem
// reguły, którą rdzeń zastosował, a nie oceną estetyczną — rdzeń nie ma gustu
// i nie udaje, że ma.
func zlozPowodZestawieniaDesignu(naglowek, tekst string, charakter *string) string {
	powod := fmt.Sprintf("nagłówek %s ma większą masę niż tekst %s, więc para daje hierarchię "+
		"bez zmiany rodziny", naglowek, tekst)
	if charakter != nil && strings.TrimSpace(*charakter) != "" {
		powod += fmt.Sprintf("; charakter wskazany przez Operatora (%s) rdzeń przekazuje dalej "+
			"bez oceny — reguła rdzenia mierzy masę i szeryfy, nie nastrój",
			strings.TrimSpace(*charakter))
	}
	return powod
}

// PodgladKroju składa podgląd kroju jako dokument SVG w kilku rozmiarach —
// obsługuje `design.font.preview`.
func (a *adapterDesignu) PodgladKroju(_ context.Context,
	z shared.DesignFontPreviewRequest) (shared.DesignFontPreviewResponse, error) {

	if strings.TrimSpace(z.FontFamily) == "" {
		return shared.DesignFontPreviewResponse{}, bladWskazaniaDesignu(
			"komenda design.font.preview bez wskazania kroju")
	}
	rozmiary := z.Sizes
	if len(rozmiary) == 0 {
		rozmiary = []float64{12, 16, 24, 36, 48}
	}
	for _, rozmiar := range rozmiary {
		if rozmiar <= 0 || rozmiar > 512 {
			return shared.DesignFontPreviewResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.font.preview z rozmiarem %v: rdzeń składa podgląd od 1 do 512", rozmiar))
		}
	}
	tekst := domyslnyTekstProbnyDesignu
	if z.SampleText != nil && strings.TrimSpace(*z.SampleText) != "" {
		tekst = strings.TrimSpace(*z.SampleText)
	}

	krojWczytany, nazwaKroju, wkompilowany, err := krojDesignu(z.FontFamily)
	if err != nil {
		// Krój, którego rdzeń nie ma, nie jest odmową komendy — pole available mówi wprost, że go nie ma.
		return shared.DesignFontPreviewResponse{
			PreviewSvg: podgladBrakuKrojuDesignu(z.FontFamily, err.Error()),
			Available:  false,
		}, nil
	}

	podglad, err := podgladKrojuSvgDesignu(krojWczytany, nazwaKroju, tekst, rozmiary)
	if err != nil {
		return shared.DesignFontPreviewResponse{}, bladWydaniaDesignu(err.Error())
	}
	glifow := liczbaGlifowKrojuDesignu(krojWczytany)
	odpowiedz := shared.DesignFontPreviewResponse{
		PreviewSvg: podglad, Available: true, GlyphCount: &glifow,
	}
	if regula := regulaOsadzeniaKrojuDesignu(nazwaKroju, wkompilowany); regula != "" {
		odpowiedz.FontFace = &regula
	}
	return odpowiedz, nil
}

// podgladKrojuSvgDesignu składa podgląd jako dokument SVG z prawdziwymi
// konturami glifów. Kontury, nie element `<text>`: dokument z `<text>`
// pokazałby krój wyłącznie tam, gdzie ten krój jest zainstalowany.
func podgladKrojuSvgDesignu(krojWczytany *sfnt.Font, nazwa, tekst string,
	rozmiary []float64) (string, error) {

	szerokosc := 0.0
	wysokosc := 24.0
	wiersze := make([]string, 0, len(rozmiary))
	for _, rozmiar := range rozmiary {
		kontury, err := sciezkaTekstuDesignu(krojWczytany, tekst, rozmiar, 16, wysokosc+rozmiar)
		if err != nil {
			return "", fmt.Errorf("podglądu kroju %s w rozmiarze %v nie da się złożyć: %w",
				nazwa, rozmiar, err)
		}
		granice := kontury.Bounds()
		szerokosc = wiekszaDesignu(szerokosc, granice.X1+16)
		wiersze = append(wiersze, kontury.ToSVG())
		wysokosc += rozmiar * 1.5
	}
	if szerokosc <= 0 {
		szerokosc = 640
	}

	var dokument strings.Builder
	fmt.Fprintf(&dokument,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%g" height="%g" viewBox="0 0 %g %g">`,
		szerokosc, wysokosc, szerokosc, wysokosc)
	fmt.Fprintf(&dokument, `<title>%s</title>`, nazwa)
	for _, wiersz := range wiersze {
		fmt.Fprintf(&dokument, `<path d="%s" fill="currentColor"/>`, wiersz)
	}
	dokument.WriteString(`</svg>`)
	return dokument.String(), nil
}

// podgladBrakuKrojuDesignu składa podgląd dla kroju, którego rdzeń nie ma.
// Dokument mówi wprost, czego brakuje — pusty prostokąt wyglądałby jak krój
// o niewidzialnych literach.
func podgladBrakuKrojuDesignu(nazwa, powod string) string {
	return fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="640" height="64" viewBox="0 0 640 64">`+
			`<title>kroju %s rdzeń nie ma</title><desc>%s</desc></svg>`,
		nazwa, powod)
}

// regulaOsadzeniaKrojuDesignu składa gotową regułę `@font-face`. Krój
// wkompilowany albo leżący na tej maszynie osadza się wprost, jako dane
// w regule; krój większy niż granica osadzenia dostaje regułę z `local()`
// i mówi to wprost.
func regulaOsadzeniaKrojuDesignu(nazwa string, wkompilowany bool) string {
	bajty := bajtyKrojuDesignu(nazwa, wkompilowany)
	if len(bajty) == 0 {
		return ""
	}
	if len(bajty) > granicaOsadzeniaKrojuDesignu {
		return fmt.Sprintf(
			"/* krój %s ma %d bajtów — powyżej granicy osadzenia %d, więc reguła wskazuje krój "+
				"zainstalowany; osadzenie tak wielkiego kroju w arkuszu stylów zatrzymałoby "+
				"pierwsze wyrysowanie strony */\n"+
				"@font-face { font-family: %q; src: local(%q); }",
			nazwa, len(bajty), granicaOsadzeniaKrojuDesignu, nazwa, nazwa)
	}
	return fmt.Sprintf(
		"@font-face { font-family: %q; src: url(data:font/ttf;base64,%s) format(\"truetype\"); }",
		nazwa, base64.StdEncoding.EncodeToString(bajty))
}

// bajtyKrojuDesignu oddaje bajty kroju — z binarium wkompilowanego albo
// z pliku znalezionego na tej maszynie.
func bajtyKrojuDesignu(nazwa string, wkompilowany bool) []byte {
	if wkompilowany {
		for nazwaKroju, bajty := range krojeWkompilowaneDesignu {
			if kluczKrojuDesignu(nazwaKroju) == kluczKrojuDesignu(nazwa) {
				return bajty
			}
		}
		return nil
	}
	for nazwaKroju, sciezka := range przegladajKrojeSerweraDesignu() {
		if kluczKrojuDesignu(nazwaKroju) != kluczKrojuDesignu(nazwa) {
			continue
		}
		bajty, err := os.ReadFile(sciezka)
		if err != nil {
			return nil
		}
		return bajty
	}
	return nil
}

// GlifyKroju oddaje glify kroju z zakresu punktów kodowych, wraz z ich
// konturami — obsługuje `design.font.glyphs.get`.
func (a *adapterDesignu) GlifyKroju(_ context.Context,
	z shared.DesignFontGlyphsGetRequest) (shared.DesignFontGlyphsGetResponse, error) {

	if strings.TrimSpace(z.FontFamily) == "" {
		return shared.DesignFontGlyphsGetResponse{}, bladWskazaniaDesignu(
			"komenda design.font.glyphs.get bez wskazania kroju")
	}
	krojWczytany, nazwaKroju, _, err := krojDesignu(z.FontFamily)
	if err != nil {
		return shared.DesignFontGlyphsGetResponse{}, bladWskazaniaDesignu(
			"komenda design.font.glyphs.get: " + err.Error())
	}

	// Zakres domyślny to łacina podstawowa i rozszerzona A — tam leżą najczęściej pytane litery polskie.
	od, do := 0x20, 0x17F
	if z.From != nil {
		od = *z.From
	}
	if z.To != nil {
		do = *z.To
	}
	if od < 0 || do < od {
		return shared.DesignFontGlyphsGetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.font.glyphs.get z zakresem od %d do %d: zakres o końcu przed "+
				"początkiem jest pusty", od, do))
	}
	granica := granicaGlifowOdpowiedziDesignu
	if z.Limit != nil {
		if *z.Limit < 1 || *z.Limit > granicaGlifowOdpowiedziDesignu {
			return shared.DesignFontGlyphsGetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.font.glyphs.get z granicą %d: rdzeń oddaje od 1 do %d glifów, "+
					"bo każdy niesie kontur", *z.Limit, granicaGlifowOdpowiedziDesignu))
		}
		granica = *z.Limit
	}

	glify := glifyKrojuDesignu(krojWczytany, od, do, granica)
	if len(glify) == 0 {
		return shared.DesignFontGlyphsGetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"krój %s nie ma ani jednego glifu w zakresie od %d do %d", nazwaKroju, od, do))
	}
	// total niesie liczbę glifów całego kroju — wykaz z zakresu jest wycinkiem, nie kompletem.
	return shared.DesignFontGlyphsGetResponse{
		Glyphs: glify, Total: liczbaGlifowKrojuDesignu(krojWczytany),
	}, nil
}

// ikonaPoKodzieDesignu odczytuje ikonę po identyfikatorze — własną z bazy albo
// wzór z katalogu wkompilowanego.
func (a *adapterDesignu) ikonaPoKodzieDesignu(ctx context.Context,
	kod string) (shared.DesignIcon, error) {

	if strings.HasPrefix(kod, zestawIkonKataloguDesignu+"/") {
		nazwa := strings.TrimPrefix(kod, zestawIkonKataloguDesignu+"/")
		wzor, jest := ikonaKataloguDesignu(nazwa)
		if !jest {
			return shared.DesignIcon{}, bladNieznanegoBytuDesignu(
				"wzoru ikony " + nazwa + " nie ma w katalogu wkompilowanym rdzenia")
		}
		return ikonaKataloguKontraktuDesignu(wzor, siatkaIkonyKataloguDesignu,
			gruboscObrysuIkonyDomyslna), nil
	}
	wiersz, err := a.repozytorium.IkonaDesignuPoKodzie(ctx, kod)
	if err != nil {
		return shared.DesignIcon{}, bladNieznanejIkonyDesignu(kod, err)
	}
	return ikonaWlasnaKontraktuDesignu(wiersz), nil
}

// ikonaWlasnaKontraktuDesignu składa `DesignIcon` kontraktu z wiersza ikony
// własnej, odczytanego z bazy.
func ikonaWlasnaKontraktuDesignu(wiersz dane.IkonaDesignu) shared.DesignIcon {
	svg := wiersz.SVG
	ikona := shared.DesignIcon{
		Id: wiersz.Kod, Name: wiersz.Nazwa, Set: wiersz.Zestaw, Tags: wiersz.Etykiety,
		Svg: &svg, StrokeWidth: wiersz.GruboscObrysu,
	}
	if wiersz.Siatka != nil {
		siatka := int(*wiersz.Siatka)
		ikona.GridSize = &siatka
	}
	return ikona
}

// bladNieznanejIkonyDesignu nazywa ikonę własną, której rdzeń nie zna, kodem
// błędu not_found, nie awarią.
func bladNieznanejIkonyDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("ikony " + kod + " nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}
