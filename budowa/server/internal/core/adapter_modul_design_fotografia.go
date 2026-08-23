// Odpowiedzialność pliku: trzynaście czynności warsztatu fotografii modułu
// Design — `design.photo.crop`, `.transform`, `.resample`, `.upscale`,
// `.enhance`, `.color.correct`, `.filter.apply`, `.retouch`, `.inpaint`,
// `.expand`, `.background.remove`, `.select.object`, `.mask.set`. Wsad, nastawy,
// łańcuch, metadane, kompozycja warstw i obrysowanie leżą
// w `adapter_modul_design_fotografia_wsad.go`; rachunek na pikselach
// w `_fotografia_rachunek.go` i `_fotografia_maski.go`.
//
// ── Obróbka zakłada WARIANT, oryginał zostaje ───────────────────────────────
// Każda czynność zapisuje NOWY zasób i wskazuje źródło polem
// `variantOfAssetId` (kolumna istniała od migracji 048). Oryginał nie jest nigdy
// nadpisywany — dzięki temu łańcuch edycji da się przejść wstecz do zdjęcia,
// które Operator wniósł, a „cofnij" nie wymaga historii operacji w pamięci okna.
// Obok wariantu powstaje wiersz czynności (migracja 346) z nastawami, którymi
// poszła — bez niego `design.photo.history.get` pokazywałby wykaz obrazków bez
// słowa o tym, co je różni.
//
// ── Wynik wychodzi jako PNG ─────────────────────────────────────────────────
// Bezstratnie i z kanałem krycia. JPEG traciłby jakość przy każdym ogniwie
// łańcucha (dziesięć korekcji to dziesięć kompresji) i nie uniósłby
// przezroczystości po odcięciu tła. Format wyjścia Operator zmienia dopiero przy
// wydaniu (`design.asset.export`), gdzie strata jest jednorazowa i świadoma.
//
// ── `computedBy` mówi PRAWDĘ o drodze rachunku ──────────────────────────────
// Cztery czynności (`upscale`, `background.remove`, `inpaint`, `expand`) mają
// DWIE drogi i obie stoją:
//
//   - `channelId` WSKAZANY — liczy kanał obrazowy. Zdjęcie Operatora jedzie do
//     punktu końcowego edycji jako obraz WEJŚCIOWY (`models.ObrazWejsciowy`
//     w roli materiału), a domalowanie i rozszerzenie kadru dokładają maskę
//     w roli maski. Odpowiedź oddaje `computedBy: kanalModelu`.
//   - `channelId` POMINIĘTY — liczy rachunek wkompilowany i odpowiedź oddaje
//     `computedBy: rachunekRdzenia`. Tak mówi kontrakt tych czterech pól i tak
//     ma być: zejście na „pierwszy czynny kanał" wysyłałoby materiał Operatora
//     do dostawcy, o którego nie prosił.
//
// Droga jest zapisana w łańcuchu edycji (`czynnosc_fotografii_design.
// policzone_przez`), więc po tygodniu da się powiedzieć, którą powstał
// konkretny wariant.
//
// WYMIAR wyniku jest zawsze ten, który czynność obiecała: `factor: 4` znaczy
// cztery razy, a `expand` o sto punktów znaczy sto — także na drodze kanału,
// który oddaje obraz o rozmiarze ze swojej nastawy. Sprowadzenie do zamówionego
// wymiaru robi rachunek rdzenia i to jedyna rzecz, którą on na tej drodze robi
// (`sprowadzWynikKanaluFotografiiDesignu`); treść pozostaje kanału i o niej mówi
// `computedBy`.
//
// Niepowodzenie wskazanego kanału jest ODMOWĄ, nie cichym zejściem na rachunek:
// Operator prosił o drogę neuronową, a obraz policzony inaczej byłby odpowiedzią
// na inne żądanie.
package core

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"strings"

	"github.com/disintegration/imaging"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// formatWynikuFotografiiDesignu jest formatem wszystkich wyników warsztatu —
// powód stoi w nagłówku pliku.
const formatWynikuFotografiiDesignu = "png"

// Kadruj kadruje zdjęcie — obsługuje `design.photo.crop`.
func (a *adapterDesignu) Kadruj(ctx context.Context,
	z shared.DesignPhotoCropRequest) (shared.DesignPhotoCropResponse, error) {

	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.crop", z.AssetId)
	if err != nil {
		return shared.DesignPhotoCropResponse{}, err
	}
	wynik, err := przytnijObrazFotografiiDesignu(obraz, z)
	if err != nil {
		return shared.DesignPhotoCropResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.crop: " + err.Error())
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoCrop, "kadr", wynik, z, nil)
	if err != nil {
		return shared.DesignPhotoCropResponse{}, err
	}
	granice := wynik.Bounds()
	return shared.DesignPhotoCropResponse{
		Asset: zasob, Width: granice.Dx(), Height: granice.Dy(),
	}, nil
}

// Przeksztalc obraca, odbija i koryguje perspektywę — obsługuje
// `design.photo.transform`.
func (a *adapterDesignu) Przeksztalc(ctx context.Context,
	z shared.DesignPhotoTransformRequest) (shared.DesignPhotoTransformResponse, error) {

	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.transform", z.AssetId)
	if err != nil {
		return shared.DesignPhotoTransformResponse{}, err
	}
	// Żądanie bez ani jednej nastawy jest odmową, nie kopią: zasób oddany jako
	// „przekształcony" bez przekształcenia zajmowałby miejsce w magazynie
	// i w łańcuchu edycji, nie różniąc się od źródła niczym.
	if (z.RotateDeg == nil || *z.RotateDeg == 0) &&
		(z.FlipHorizontal == nil || !*z.FlipHorizontal) &&
		(z.FlipVertical == nil || !*z.FlipVertical) &&
		(z.LensDistortion == nil || *z.LensDistortion == 0) && len(z.Perspective) == 0 {

		return shared.DesignPhotoTransformResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.transform bez ani jednej nastawy: obrót zerowy, brak odbić, " +
				"brak korekcji perspektywy i brak korekcji obiektywu dałyby kopię źródła")
	}
	wynik, err := przeksztalcObrazFotografiiDesignu(obraz, z)
	if err != nil {
		return shared.DesignPhotoTransformResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.transform: " + err.Error())
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoTransform, "przekształcenie", wynik, z, nil)
	if err != nil {
		return shared.DesignPhotoTransformResponse{}, err
	}
	granice := wynik.Bounds()
	return shared.DesignPhotoTransformResponse{
		Asset: zasob, Width: granice.Dx(), Height: granice.Dy(),
	}, nil
}

// PrzeliczRozdzielczosc przelicza rozdzielczość — obsługuje
// `design.photo.resample`.
func (a *adapterDesignu) PrzeliczRozdzielczosc(ctx context.Context,
	z shared.DesignPhotoResampleRequest) (shared.DesignPhotoResampleResponse, error) {

	filtr := shared.DesignPhotoResampleFilter(shared.DesignPhotoResampleFilterLanczos)
	if z.Filter != nil {
		if err := sprawdzWyliczenieDesignu("design.photo.resample", "filter", *z.Filter,
			shared.WartosciDesignPhotoResampleFilter()); err != nil {
			return shared.DesignPhotoResampleResponse{}, err
		}
		filtr = *z.Filter
	}
	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.resample", z.AssetId)
	if err != nil {
		return shared.DesignPhotoResampleResponse{}, err
	}
	szerokosc, wysokosc := 0, 0
	if z.Width != nil {
		szerokosc = *z.Width
	}
	if z.Height != nil {
		wysokosc = *z.Height
	}
	wynik, err := przeliczRozdzielczoscFotografiiDesignu(obraz, szerokosc, wysokosc, filtr,
		z.KeepAspectRatio == nil || *z.KeepAspectRatio)
	if err != nil {
		return shared.DesignPhotoResampleResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.resample: " + err.Error())
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoResample, "rozdzielczość", wynik, z, nil)
	if err != nil {
		return shared.DesignPhotoResampleResponse{}, err
	}
	granice := wynik.Bounds()
	return shared.DesignPhotoResampleResponse{
		Asset: zasob, Width: granice.Dx(), Height: granice.Dy(),
	}, nil
}

// Powieksz powiększa zdjęcie — obsługuje `design.photo.upscale`.
func (a *adapterDesignu) Powieksz(ctx context.Context,
	z shared.DesignPhotoUpscaleRequest) (shared.DesignPhotoUpscaleResponse, error) {

	if z.Factor != 2 && z.Factor != 4 && z.Factor != 8 {
		return shared.DesignPhotoUpscaleResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.photo.upscale z krotnością %d: rdzeń powiększa 2, 4 albo 8 razy — "+
				"krotność pośrednia daje ten sam wynik co przeliczenie rozdzielczości "+
				"(design.photo.resample)", z.Factor))
	}
	if err := a.sprawdzKanalObrazowyFotografiiDesignu(z.ChannelId); err != nil {
		return shared.DesignPhotoUpscaleResponse{}, err
	}
	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.upscale", z.AssetId)
	if err != nil {
		return shared.DesignPhotoUpscaleResponse{}, err
	}
	droga := shared.DesignPhotoComputeRoute(shared.DesignPhotoComputeRouteRachunekRdzenia)
	var wynik image.Image
	if czyDrogaKanaluFotografiiDesignu(z.ChannelId) {
		zrodloweGranice := obraz.Bounds()
		wynik, err = a.obrazKanalemFotografiiDesignu(ctx, shared.CommandDesignPhotoUpscale,
			z.ChannelId, z.WindowId, fmt.Sprintf(
				"powiększ to zdjęcie %d razy, zachowując treść i wyostrzając szczegóły",
				z.Factor), obraz, nil,
			zrodloweGranice.Dx()*z.Factor, zrodloweGranice.Dy()*z.Factor)
		if err != nil {
			return shared.DesignPhotoUpscaleResponse{}, err
		}
		droga = shared.DesignPhotoComputeRouteKanalModelu
	} else {
		wynik, err = powiekszRachunkiemDesignu(obraz, z.Factor,
			z.SharpenAfter == nil || *z.SharpenAfter)
		if err != nil {
			return shared.DesignPhotoUpscaleResponse{}, bladWskazaniaDesignu(
				"komenda design.photo.upscale: " + err.Error())
		}
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoUpscale, fmt.Sprintf("powiększenie ×%d", z.Factor),
		wynik, z, &droga)
	if err != nil {
		return shared.DesignPhotoUpscaleResponse{}, err
	}
	granice := wynik.Bounds()
	return shared.DesignPhotoUpscaleResponse{
		Asset: zasob, ComputedBy: droga, Width: granice.Dx(), Height: granice.Dy(),
	}, nil
}

// PopraweJakosc poprawia jakość zdjęcia — obsługuje `design.photo.enhance`.
func (a *adapterDesignu) PopraweJakosc(ctx context.Context,
	z shared.DesignPhotoEnhanceRequest) (shared.DesignPhotoEnhanceResponse, error) {

	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.enhance", z.AssetId)
	if err != nil {
		return shared.DesignPhotoEnhanceResponse{}, err
	}
	wynik, kroki := popraweJakoscFotografiiDesignu(obraz, z)
	if len(kroki) == 0 {
		return shared.DesignPhotoEnhanceResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.enhance z samymi wartościami wyłączonymi: żaden krok nie " +
				"wszedłby, a zasób oddany jako poprawiony bez poprawy byłby kopią źródła")
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoEnhance, "poprawa jakości", wynik, z, nil)
	if err != nil {
		return shared.DesignPhotoEnhanceResponse{}, err
	}
	return shared.DesignPhotoEnhanceResponse{Asset: zasob, AppliedSteps: kroki}, nil
}

// SkorygujBarwe koryguje barwę zdjęcia — obsługuje `design.photo.color.correct`.
//
// Pole `histogramShift` jest POMIAREM: rdzeń liczy średnią jasność źródła
// i wyniku i oddaje różnicę. Sprawdzian skutku ma po czym poznać, że korekcja
// naprawdę ruszyła piksele, a nie tylko zwróciła `status: ok`.
func (a *adapterDesignu) SkorygujBarwe(ctx context.Context,
	z shared.DesignPhotoColorCorrectRequest) (shared.DesignPhotoColorCorrectResponse, error) {

	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.color.correct", z.AssetId)
	if err != nil {
		return shared.DesignPhotoColorCorrectResponse{}, err
	}
	wynik, err := skorygujBarweFotografiiDesignu(obraz, z)
	if err != nil {
		return shared.DesignPhotoColorCorrectResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.color.correct: " + err.Error())
	}
	przesuniecie := sredniaJasnoscDesignu(wynik) - sredniaJasnoscDesignu(obraz)
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoColorCorrect, "korekcja barwy", wynik, z, nil)
	if err != nil {
		return shared.DesignPhotoColorCorrectResponse{}, err
	}
	return shared.DesignPhotoColorCorrectResponse{
		Asset: zasob, HistogramShift: przesuniecie,
	}, nil
}

// sredniaJasnoscDesignu liczy średnią jasność obrazu w skali 0–255 — miara
// pomiaru przesunięcia histogramu.
func sredniaJasnoscDesignu(obraz image.Image) float64 {
	granice := obraz.Bounds()
	suma, punktow := 0.0, 0
	// Próbkowanie co czwarty punkt w obu osiach: średnia z co szesnastego punktu
	// różni się od pełnej o setne części jednostki, a rachunek jest szesnaście
	// razy tańszy. Pomiar ma być pomiarem, nie kosztem.
	for y := granice.Min.Y; y < granice.Max.Y; y += 4 {
		for x := granice.Min.X; x < granice.Max.X; x += 4 {
			r, g, b, _ := obraz.At(x, y).RGBA()
			suma += float64(int(r>>8)*299+int(g>>8)*587+int(b>>8)*114) / 1000
			punktow++
		}
	}
	if punktow == 0 {
		return 0
	}
	return suma / float64(punktow)
}

// NalozFiltr nakłada filtr obrazu — obsługuje `design.photo.filter.apply`.
func (a *adapterDesignu) NalozFiltr(ctx context.Context,
	z shared.DesignPhotoFilterApplyRequest) (shared.DesignPhotoFilterApplyResponse, error) {

	if err := sprawdzWyliczenieDesignu("design.photo.filter.apply", "filter", z.Filter,
		shared.WartosciDesignPhotoFilter()); err != nil {
		return shared.DesignPhotoFilterApplyResponse{}, err
	}
	sila := 0.5
	if z.Amount != nil {
		if *z.Amount < 0 || *z.Amount > 1 {
			return shared.DesignPhotoFilterApplyResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.photo.filter.apply z siłą %v: siła filtru idzie od 0 do 1", *z.Amount))
		}
		sila = *z.Amount
	}
	if sila == 0 {
		return shared.DesignPhotoFilterApplyResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.filter.apply z siłą zero: filtr o zerowej sile dałby kopię źródła")
	}
	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.filter.apply", z.AssetId)
	if err != nil {
		return shared.DesignPhotoFilterApplyResponse{}, err
	}
	wynik := nalozFiltrFotografiiDesignu(obraz, z.Filter, sila)
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoFilterApply, "filtr "+string(z.Filter), wynik, z, nil)
	if err != nil {
		return shared.DesignPhotoFilterApplyResponse{}, err
	}
	return shared.DesignPhotoFilterApplyResponse{Asset: zasob}, nil
}

// Wyretuszuj retuszuje obszary — obsługuje `design.photo.retouch`.
func (a *adapterDesignu) Wyretuszuj(ctx context.Context,
	z shared.DesignPhotoRetouchRequest) (shared.DesignPhotoRetouchResponse, error) {

	if len(z.Regions) == 0 {
		return shared.DesignPhotoRetouchResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.retouch bez ani jednego obszaru: rdzeń nie zgaduje, co na " +
				"zdjęciu jest niedoskonałością")
	}
	tryb := shared.DesignPhotoRetouchMode(shared.DesignPhotoRetouchModeHeal)
	if z.Mode != nil {
		if err := sprawdzWyliczenieDesignu("design.photo.retouch", "mode", *z.Mode,
			shared.WartosciDesignPhotoRetouchMode()); err != nil {
			return shared.DesignPhotoRetouchResponse{}, err
		}
		tryb = *z.Mode
	}
	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.retouch", z.AssetId)
	if err != nil {
		return shared.DesignPhotoRetouchResponse{}, err
	}
	wynik, weszlo, pominiete := wyretuszujObszaryDesignu(obraz, z.Regions, tryb, z.Source)
	if weszlo == 0 {
		return shared.DesignPhotoRetouchResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"żaden z %d obszarów nie wszedł: %s", len(z.Regions), strings.Join(pominiete, "; ")))
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoRetouch, "retusz", wynik, z, nil)
	if err != nil {
		return shared.DesignPhotoRetouchResponse{}, err
	}
	return shared.DesignPhotoRetouchResponse{
		Asset: zasob, RegionsApplied: weszlo, RegionsSkipped: uporzadkujBilansDesignu(pominiete),
	}, nil
}

// Domaluj domalowuje obszar z maski — obsługuje `design.photo.inpaint`.
func (a *adapterDesignu) Domaluj(ctx context.Context,
	z shared.DesignPhotoInpaintRequest) (shared.DesignPhotoInpaintResponse, error) {

	maZaMaske := z.MaskAssetId != nil && strings.TrimSpace(*z.MaskAssetId) != ""
	if !maZaMaske && len(z.Regions) == 0 {
		return shared.DesignPhotoInpaintResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.inpaint bez maski i bez obszarów: rdzeń nie zgaduje, który " +
				"fragment zdjęcia domalować")
	}
	if err := a.sprawdzKanalObrazowyFotografiiDesignu(z.ChannelId); err != nil {
		return shared.DesignPhotoInpaintResponse{}, err
	}
	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.inpaint", z.AssetId)
	if err != nil {
		return shared.DesignPhotoInpaintResponse{}, err
	}
	granice := obraz.Bounds()

	naleziObszar := maskaZObszarowDesignu(z.Regions)
	if maZaMaske {
		maska, _, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.inpaint", *z.MaskAssetId)
		if err != nil {
			return shared.DesignPhotoInpaintResponse{}, err
		}
		naleziObszar = maskaZObrazuDesignu(maska, granice.Dx(), granice.Dy(), false)
	}

	droga := shared.DesignPhotoComputeRoute(shared.DesignPhotoComputeRouteRachunekRdzenia)
	var wynik image.Image
	if czyDrogaKanaluFotografiiDesignu(z.ChannelId) {
		// Maska jedzie do kanału jako OBRAZ, bo tak przyjmuje ją punkt końcowy
		// edycji. Rdzeń ma już maskę policzoną (z zasobu albo z obszarów), więc
		// zamienia ją w obraz jednokanałowy tą samą drogą, którą zapisuje wyniki.
		polecenie := "domaluj obszar wskazany maską tak, żeby wtopił się w otoczenie"
		if z.Prompt != nil && strings.TrimSpace(*z.Prompt) != "" {
			polecenie = strings.TrimSpace(*z.Prompt)
		}
		wynik, err = a.obrazKanalemFotografiiDesignu(ctx, shared.CommandDesignPhotoInpaint,
			z.ChannelId, z.WindowId, polecenie, obraz,
			obrazMaskiDesignu(naleziObszar, granice.Dx(), granice.Dy()),
			granice.Dx(), granice.Dy())
		if err != nil {
			return shared.DesignPhotoInpaintResponse{}, err
		}
		droga = shared.DesignPhotoComputeRouteKanalModelu
	} else {
		wynik = domalujObszarRachunkiemDesignu(obraz, naleziObszar, granice.Dx()+granice.Dy())
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoInpaint, "domalowanie", wynik, z, &droga)
	if err != nil {
		return shared.DesignPhotoInpaintResponse{}, err
	}
	return shared.DesignPhotoInpaintResponse{Asset: zasob, ComputedBy: droga}, nil
}

// RozszerzKadr rozszerza kadr poza pierwotną ramkę — obsługuje
// `design.photo.expand`.
func (a *adapterDesignu) RozszerzKadr(ctx context.Context,
	z shared.DesignPhotoExpandRequest) (shared.DesignPhotoExpandResponse, error) {

	lewa, prawa, gora, dol := 0, 0, 0, 0
	if z.Left != nil {
		lewa = *z.Left
	}
	if z.Right != nil {
		prawa = *z.Right
	}
	if z.Top != nil {
		gora = *z.Top
	}
	if z.Bottom != nil {
		dol = *z.Bottom
	}
	if lewa < 0 || prawa < 0 || gora < 0 || dol < 0 {
		return shared.DesignPhotoExpandResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.expand z ujemnym przyrostem: rozszerzenie kadru nie obcina " +
				"obrazu; obcięcie zamawia się komendą design.photo.crop")
	}
	if lewa+prawa+gora+dol == 0 {
		return shared.DesignPhotoExpandResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.expand z zerowym przyrostem po każdej stronie dałaby kopię źródła")
	}
	if err := a.sprawdzKanalObrazowyFotografiiDesignu(z.ChannelId); err != nil {
		return shared.DesignPhotoExpandResponse{}, err
	}
	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.expand", z.AssetId)
	if err != nil {
		return shared.DesignPhotoExpandResponse{}, err
	}
	granice := obraz.Bounds()
	if err := sprawdzRozmiarFotografiiDesignu(granice.Dx()+lewa+prawa,
		granice.Dy()+gora+dol); err != nil {
		return shared.DesignPhotoExpandResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.expand: " + err.Error())
	}

	droga := shared.DesignPhotoComputeRoute(shared.DesignPhotoComputeRouteRachunekRdzenia)
	var wynik image.Image
	if czyDrogaKanaluFotografiiDesignu(z.ChannelId) {
		// Materiałem jest płótno POWIĘKSZONE z oryginałem w środku i pustymi
		// marginesami, a maską — same marginesy. Wysłanie samego oryginału kazałoby
		// modelowi domyślać się, gdzie ma dorysować i ile; płótno mówi to wymiarem.
		polecenie := "dorysuj brakujące części obrazu w obszarze wskazanym maską, " +
			"kontynuując treść zdjęcia"
		if z.Prompt != nil && strings.TrimSpace(*z.Prompt) != "" {
			polecenie = strings.TrimSpace(*z.Prompt)
		}
		material, maska := plotnoRozszerzeniaDesignu(obraz, lewa, prawa, gora, dol)
		wynik, err = a.obrazKanalemFotografiiDesignu(ctx, shared.CommandDesignPhotoExpand,
			z.ChannelId, z.WindowId, polecenie, material, maska,
			granice.Dx()+lewa+prawa, granice.Dy()+gora+dol)
		if err != nil {
			return shared.DesignPhotoExpandResponse{}, err
		}
		droga = shared.DesignPhotoComputeRouteKanalModelu
	} else {
		wynik = rozszerzKadrRachunkiemDesignu(obraz, lewa, prawa, gora, dol)
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoExpand, "rozszerzenie kadru", wynik, z, &droga)
	if err != nil {
		return shared.DesignPhotoExpandResponse{}, err
	}
	wynikoweGranice := wynik.Bounds()
	return shared.DesignPhotoExpandResponse{
		Asset: zasob, ComputedBy: droga,
		Width: wynikoweGranice.Dx(), Height: wynikoweGranice.Dy(),
	}, nil
}

// OdetnijTlo odcina tło — obsługuje `design.photo.background.remove`.
func (a *adapterDesignu) OdetnijTlo(ctx context.Context,
	z shared.DesignPhotoBackgroundRemoveRequest) (shared.DesignPhotoBackgroundRemoveResponse, error) {

	tolerancja := domyslnaTolerancjaTlaDesignu
	if z.Tolerance != nil {
		if *z.Tolerance < 0 || *z.Tolerance > 1 {
			return shared.DesignPhotoBackgroundRemoveResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.photo.background.remove z tolerancją %v: tolerancja idzie od 0 do 1",
				*z.Tolerance))
		}
		tolerancja = *z.Tolerance
	}
	if err := a.sprawdzKanalObrazowyFotografiiDesignu(z.ChannelId); err != nil {
		return shared.DesignPhotoBackgroundRemoveResponse{}, err
	}
	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx,
		"design.photo.background.remove", z.AssetId)
	if err != nil {
		return shared.DesignPhotoBackgroundRemoveResponse{}, err
	}

	droga := shared.DesignPhotoComputeRoute(shared.DesignPhotoComputeRouteRachunekRdzenia)
	var wynik image.Image
	udzial := 0.0
	if czyDrogaKanaluFotografiiDesignu(z.ChannelId) {
		granice := obraz.Bounds()
		wynik, err = a.obrazKanalemFotografiiDesignu(ctx,
			shared.CommandDesignPhotoBackgroundRemove, z.ChannelId, z.WindowId,
			"odetnij tło tego zdjęcia i oddaj sam przedmiot na tle przezroczystym",
			obraz, nil, granice.Dx(), granice.Dy())
		if err != nil {
			return shared.DesignPhotoBackgroundRemoveResponse{}, err
		}
		// Udział punktów przezroczystych liczy się z PLIKU kanału, nie z założenia:
		// obraz bez ani jednego punktu przezroczystego nie jest odcięciem tła i pola
		// `hasAlpha: true` nie wolno nad nim postawić.
		udzial = udzialPrzezroczystosciDesignu(wynik)
		if udzial == 0 {
			return shared.DesignPhotoBackgroundRemoveResponse{}, bladWydaniaDesignu(fmt.Sprintf(
				"komenda design.photo.background.remove: kanał %s oddał obraz bez ani jednego "+
					"punktu przezroczystego — to nie jest odcięcie tła, a odpowiedź `hasAlpha: true` "+
					"nad takim plikiem byłaby nieprawdą; naprawa: sprawdzić, czy punkt końcowy kanału "+
					"oddaje PNG z kanałem krycia, albo pominąć channelId i policzyć rachunkiem rdzenia",
				strings.TrimSpace(*z.ChannelId)))
		}
		droga = shared.DesignPhotoComputeRouteKanalModelu
	} else {
		wynik, udzial = odetnijTloRachunkiemDesignu(obraz, tolerancja)
		if udzial == 0 {
			return shared.DesignPhotoBackgroundRemoveResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"przy tolerancji %.2f rdzeń nie uznał ANI JEDNEGO punktu za tło — zasób oddany "+
					"jako „bez tła\" byłby kopią źródła; naprawa: podnieść tolerancję albo wskazać "+
					"zdjęcie na jednolitym tle", tolerancja))
		}
		if udzial > 0.98 {
			return shared.DesignPhotoBackgroundRemoveResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"przy tolerancji %.2f tłem okazało się %.0f%% obrazu — wynik byłby niemal pustym "+
					"płótnem; naprawa: obniżyć tolerancję", tolerancja, udzial*100))
		}
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoBackgroundRemove, "bez tła", wynik, z, &droga)
	if err != nil {
		return shared.DesignPhotoBackgroundRemoveResponse{}, err
	}
	return shared.DesignPhotoBackgroundRemoveResponse{
		Asset: zasob, ComputedBy: droga, HasAlpha: true, TransparentShare: udzial,
	}, nil
}

// ZaznaczObiekt zaznacza obiekt i oddaje maskę — obsługuje
// `design.photo.select.object`.
func (a *adapterDesignu) ZaznaczObiekt(ctx context.Context,
	z shared.DesignPhotoSelectObjectRequest) (shared.DesignPhotoSelectObjectResponse, error) {

	tolerancja := domyslnaTolerancjaTlaDesignu
	if z.Tolerance != nil {
		if *z.Tolerance < 0 || *z.Tolerance > 1 {
			return shared.DesignPhotoSelectObjectResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.photo.select.object z tolerancją %v: tolerancja idzie od 0 do 1",
				*z.Tolerance))
		}
		tolerancja = *z.Tolerance
	}
	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx,
		"design.photo.select.object", z.AssetId)
	if err != nil {
		return shared.DesignPhotoSelectObjectResponse{}, err
	}
	maska, udzial, prostokat, err := zaznaczObiektDesignu(obraz,
		int(z.Point.X), int(z.Point.Y), tolerancja)
	if err != nil {
		return shared.DesignPhotoSelectObjectResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.select.object: " + err.Error())
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoSelectObject, "maska zaznaczenia", maska, z, nil)
	if err != nil {
		return shared.DesignPhotoSelectObjectResponse{}, err
	}
	obszar := shared.DesignPhotoRegion{
		X: float64(prostokat.Min.X), Y: float64(prostokat.Min.Y),
		Width: float64(prostokat.Dx()), Height: float64(prostokat.Dy()),
	}
	return shared.DesignPhotoSelectObjectResponse{
		Mask: zasob, Coverage: udzial, Bounds: &obszar,
	}, nil
}

// UstawMaske zakłada maskę nieniszczącą — obsługuje `design.photo.mask.set`.
//
// „Nieniszcząca" znaczy tu dokładnie to: ŹRÓDŁO zostaje nietknięte, a maska
// wchodzi do wariantu i do łańcucha edycji razem ze wskazaniem, którym zasobem
// jest. Operator zdejmuje ją wracając do źródła, którego nikt nie nadpisał.
func (a *adapterDesignu) UstawMaske(ctx context.Context,
	z shared.DesignPhotoMaskSetRequest) (shared.DesignPhotoMaskSetResponse, error) {

	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.mask.set", z.AssetId)
	if err != nil {
		return shared.DesignPhotoMaskSetResponse{}, err
	}
	maska, _, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.mask.set", z.MaskAssetId)
	if err != nil {
		return shared.DesignPhotoMaskSetResponse{}, err
	}
	if z.FeatherPx != nil && *z.FeatherPx < 0 {
		return shared.DesignPhotoMaskSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.photo.mask.set z miękkością %v: miękkość ujemna nie istnieje",
			*z.FeatherPx))
	}

	granice := obraz.Bounds()
	odwroc := z.Invert != nil && *z.Invert
	miekkosc := 0.0
	if z.FeatherPx != nil {
		miekkosc = *z.FeatherPx
	}
	wynik := nalozMaskeNaObrazDesignu(obraz, maska, odwroc, miekkosc)
	if granice.Dx() < 1 || granice.Dy() < 1 {
		return shared.DesignPhotoMaskSetResponse{}, bladWskazaniaDesignu(
			"zasób o zerowym boku nie ma czego zamaskować")
	}
	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, zrodlo, z.WindowId,
		shared.CommandDesignPhotoMaskSet, "maska", wynik, z, nil)
	if err != nil {
		return shared.DesignPhotoMaskSetResponse{}, err
	}
	return shared.DesignPhotoMaskSetResponse{
		Asset: zasob, MaskAssetId: strings.TrimSpace(z.MaskAssetId),
	}, nil
}

// nalozMaskeNaObrazDesignu przepuszcza obraz przez maskę: jasność maski staje się
// kryciem punktu.
func nalozMaskeNaObrazDesignu(obraz, maska image.Image, odwroc bool,
	miekkosc float64) image.Image {

	granice := obraz.Bounds()
	uzytaMaska := maska
	if miekkosc > 0 {
		// Miękkość krawędzi to rozmycie SAMEJ maski — rozmycie obrazu zmieniłoby
		// jego treść, a nie kształt zaznaczenia.
		uzytaMaska = nalozFiltrFotografiiDesignu(maska, shared.DesignPhotoFilterBlur,
			przytnijUlamekDesignu(miekkosc/8))
	}
	wynik := image.NewNRGBA(image.Rect(0, 0, granice.Dx(), granice.Dy()))
	for y := 0; y < granice.Dy(); y++ {
		for x := 0; x < granice.Dx(); x++ {
			punkt := nrgbaPunktuDesignu(obraz, granice.Min.X+x, granice.Min.Y+y)
			udzial := udzialMaskiDesignu(uzytaMaska, x, y, granice.Dx(), granice.Dy())
			if odwroc {
				udzial = 1 - udzial
			}
			punkt.A = przytnijSkladowaDesignu(float64(punkt.A) * udzial)
			wynik.SetNRGBA(x, y, punkt)
		}
	}
	return wynik
}

// sprawdzKanalObrazowyFotografiiDesignu sprawdza kanał wskazany przy czynności
// o wariancie neuronowym.
//
// Wskazanie kanału, którego nie ma, jest ODMOWĄ, nie ciszą: Operator prosił
// o drogę neuronową i ma prawo wiedzieć, że wskazał kanał nieistniejący.
// Sprawdzenie idzie tą samą drogą, co przy `design.asset.generate`, więc odmowy
// są te same i Operator naprawia je tymi samymi ruchami.
//
// Sprawdzenie zostaje przed odczytem zasobu, choć wywołanie kanału i tak wybrało
// by go ponownie: odmowa za zły kanał ma padać PRZED wczytaniem zdjęcia
// z magazynu, żeby Operator nie czekał na rachunek, którego wynik i tak nie
// wyjdzie.
func (a *adapterDesignu) sprawdzKanalObrazowyFotografiiDesignu(kanal *string) error {
	if kanal == nil || strings.TrimSpace(*kanal) == "" {
		return nil
	}
	if _, err := a.kanalObrazowyZadania(kanal, shared.DesignPrompt{
		Subject: "obróbka zdjęcia warsztatem fotografii",
	}, 1); err != nil {
		return err
	}
	return nil
}

// czyDrogaKanaluFotografiiDesignu rozstrzyga, czy czynność idzie kanałem modelu.
//
// Rozstrzyga WSKAZANIE, nie dostępność: kontrakt mówi o tych czterech polach
// wprost — „gdy stoi, liczy kanał, gdy nie stoi — rachunek wkompilowany".
// Zejście na pierwszy czynny kanał obrazowy przy pominiętym polu wysyłałoby
// zdjęcie Operatora do dostawcy, o którego nie prosił, i kosztowało go pieniądze
// bez jednego słowa.
func czyDrogaKanaluFotografiiDesignu(kanal *string) bool {
	return kanal != nil && strings.TrimSpace(*kanal) != ""
}

// obrazKanalemFotografiiDesignu wysyła materiał (i maskę, gdy czynność ją ma) do
// kanału obrazowego i oddaje obraz, który kanał policzył.
//
// ── Wynik jest sprowadzany do wymiaru, który czynność OBIECAŁA ──────────────
// Punkt końcowy edycji oddaje obraz o rozmiarze ze swojej nastawy, nie o tym,
// który wynika z żądania Operatora: `factor: 4` znaczy dokładnie cztery razy,
// a `expand` o sto punktów z lewej znaczy dokładnie sto. Dlatego wynik kanału
// dochodzi do zamówionego wymiaru przeliczeniem rozdzielczości rdzenia —
// najpierw kadrem do proporcji celu (żeby obraz nie został ściśnięty), potem
// przeliczeniem. TREŚĆ pozostaje kanału i to o niej mówi `computedBy`; wymiar
// jest zamówieniem Operatora i o nim mówią pola `width` i `height`.
//
// ── Niepowodzenie kanału jest ODMOWĄ, nie cichym zejściem na rachunek ───────
// Operator wskazał kanał wprost, więc wynik rachunku wkompilowanego byłby
// odpowiedzią na inne żądanie. Zdanie kanału (brak poświadczenia, odmowa
// dostawcy, odpowiedź bez obrazu) jedzie do niego w treści odmowy, bo kanał wie
// o swoim braku więcej niż ten moduł.
func (a *adapterDesignu) obrazKanalemFotografiiDesignu(ctx context.Context, komenda string,
	kanal, okno *string, polecenie string, material, maska image.Image,
	szerokosc, wysokosc int) (image.Image, error) {

	wiersz, err := a.kanalObrazowyZadania(kanal, shared.DesignPrompt{Subject: polecenie}, 1)
	if err != nil {
		return nil, err
	}

	obrazy := []models.ObrazWejsciowy{}
	bajtyMaterialu, _, err := zakodujObrazDesignu(material, formatWynikuFotografiiDesignu, nil)
	if err != nil {
		return nil, bladWydaniaDesignu(err.Error())
	}
	obrazy = append(obrazy, models.ObrazWejsciowy{
		Rola: models.RolaObrazuMaterial, TypTresci: "image/png",
		Base64: base64.StdEncoding.EncodeToString(bajtyMaterialu), Nazwa: "material.png",
	})
	if maska != nil {
		bajtyMaski, _, err := zakodujObrazDesignu(maska, formatWynikuFotografiiDesignu, nil)
		if err != nil {
			return nil, bladWydaniaDesignu(err.Error())
		}
		obrazy = append(obrazy, models.ObrazWejsciowy{
			Rola: models.RolaObrazuMaska, TypTresci: "image/png",
			Base64: base64.StdEncoding.EncodeToString(bajtyMaski), Nazwa: "maska.png",
		})
	}

	bajty, _, err := a.wytworzObraz(ctx, wiersz, oknoWytworu(okno, ""), polecenie, obrazy)
	if err != nil {
		return nil, err
	}
	wynik, _, err := image.Decode(bytes.NewReader(bajty))
	if err != nil {
		return nil, bladWydaniaDesignu(fmt.Sprintf(
			"komenda %s: kanał %s oddał treść, która nie jest obrazem: %s",
			komenda, wiersz.Identyfikator(), err.Error()))
	}
	return sprowadzWynikKanaluFotografiiDesignu(wynik, szerokosc, wysokosc), nil
}

// sprowadzWynikKanaluFotografiiDesignu sprowadza obraz oddany przez kanał do
// wymiaru zamówionego przez czynność. Wymiar zgodny zostawia obraz nietknięty —
// przeliczenie „w tę samą rozdzielczość" byłoby stratą jakości bez powodu.
func sprowadzWynikKanaluFotografiiDesignu(wynik image.Image,
	szerokosc, wysokosc int) image.Image {

	if szerokosc < 1 || wysokosc < 1 {
		return wynik
	}
	granice := wynik.Bounds()
	if granice.Dx() == szerokosc && granice.Dy() == wysokosc {
		return wynik
	}
	// Kadr do proporcji celu przed przeliczeniem: bez niego obraz kwadratowy
	// oddany na żądanie o proporcji 4:3 zostałby ściśnięty, a ściśnięte zdjęcie
	// jest zdjęciem zepsutym, nie przeliczonym.
	przyciety := wpiszProporcjeWKadrDesignu(granice, float64(szerokosc)/float64(wysokosc))
	if przyciety.Dx() > 0 && przyciety.Dy() > 0 && przyciety != granice {
		wynik = imaging.Crop(wynik, przyciety)
	}
	return imaging.Resize(wynik, szerokosc, wysokosc, imaging.Lanczos)
}

// zapiszWariantFotografiiDesignu utrwala wynik obróbki jako WARIANT źródła
// i dokłada ogniwo do łańcucha edycji.
//
// Kolejność jest ta sama, co przy wniesieniu i przy generowaniu: najpierw bajty
// w magazynie, potem wiersz zasobu, na końcu ogniwo łańcucha. Ogniwo na końcu, bo
// wskazuje klucz wiersza, którego przed zapisem nie ma.
//
// Niepowodzenie zapisu OGNIWA nie unieważnia komendy: bajty i wiersz są już
// prawdziwe, a odmowa dlatego, że nie dało się zapisać opisu obróbki, zabierałaby
// Operatorowi rzecz, o którą prosił. Zasób wychodzi wtedy bez ogniwa — czyli
// `design.photo.history.get` mówi prawdę o tym, czego rdzeń o nim nie wie.
func (a *adapterDesignu) zapiszWariantFotografiiDesignu(ctx context.Context,
	zrodlo dane.ZasobDesignu, okno *string, komenda, opis string, obraz image.Image,
	nastawy any, droga *shared.DesignPhotoComputeRoute) (shared.DesignAsset, error) {

	bajty, _, err := zakodujObrazDesignu(obraz, formatWynikuFotografiiDesignu, nil)
	if err != nil {
		return shared.DesignAsset{}, bladWydaniaDesignu(err.Error())
	}
	nazwa := nazwaZasobuDesignu(zrodlo) + " — " + opis
	zapisany, err := a.zalozZasobZBajtowDesignu(ctx, oknoWytworu(okno, zrodlo.Okno), nazwa,
		shared.DesignAssetKindImage, formatWynikuFotografiiDesignu, bajty)
	if err != nil {
		return shared.DesignAsset{}, err
	}

	// Wariant wskazuje ŹRÓDŁO — po tym łańcuch edycji da się przejść wstecz do
	// zdjęcia wniesionego przez Operatora.
	zapisany.WariantZasobuID = &zrodlo.Kod
	zapisany, err = a.repozytorium.ZapiszZasob(ctx, zapisany)
	if err != nil {
		return shared.DesignAsset{}, bladDesignu(err)
	}

	ogniwo := dane.CzynnoscFotografiiDesignu{
		ZasobID:       zapisany.ID,
		ZasobZrodlaID: &zrodlo.ID,
		Komenda:       komenda,
	}
	if zapis, err := json.Marshal(nastawy); err == nil {
		tresc := string(zapis)
		ogniwo.NastawyJSON = &tresc
	}
	if droga != nil {
		tresc := string(*droga)
		ogniwo.PoliczonePrzez = &tresc
	}
	_ = a.repozytorium.ZapiszCzynnoscFotografiiDesignu(ctx, ogniwo)

	return zasobWytworzonyKontraktu(zapisany), nil
}
