// Odpowiedzialność pliku: siedem czynności warsztatu fotografii pracujących nad
// WIELOMA zasobami albo nad wiedzą o zasobie — `design.photo.layer.composite`,
// `.batch.apply`, `.preset.save`, `.preset.list`, `.vectorize`, `.metadata.get`,
// `.history.get`. Czynności nad jednym obrazem leżą
// w `adapter_modul_design_fotografia.go`.
//
// ── Wsad powtarza CZYNNOŚCI, nie kopiuje wyniku ─────────────────────────────
// `design.photo.batch.apply` puszcza ten sam zestaw czynności na każdym zasobie
// osobno, przez te same uchwyty, którymi jadą czynności pojedyncze. Nie ma tu
// drugiej drogi rachunku: gdyby wsad liczył po swojemu, wynik wsadowy różniłby
// się od pojedynczego i Operator nie mógłby zaufać żadnemu z nich.
//
// ── Zasób, którego nie udało się przetworzyć, WRACA w bilansie ──────────────
// `failedAssetIds` niesie zasób wraz z powodem. Wsad na stu zdjęciach, z których
// trzy padły, wygląda bez tego pola jak wsad kompletny, a brak wyjdzie na jaw
// dopiero przy przeglądaniu wyników.
//
// ── Metadane są POMIAREM z pliku, nie echem wiersza ─────────────────────────
// `design.photo.metadata.get` czyta nagłówek pliku leżącego w magazynie: wymiary,
// format, model barwny, obecność kanału krycia, rozdzielczość i pola EXIF. Wiersz
// bazy niesie tylko to, co zmierzono przy wniesieniu — a plik może być jedyną
// prawdą o tym, co Operator naprawdę ma.
package core

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// przedrostekNastawyFotografiiDesign znakuje identyfikatory zewnętrzne nastaw.
	przedrostekNastawyFotografiiDesign = "nastawa-foto-"

	// granicaZasobowWsaduDesignu chroni wsad przed żądaniem na dziesięciu
	// tysiącach zasobów: każdy przechodzi pełny rachunek i zakłada nowe bajty
	// w magazynie.
	granicaZasobowWsaduDesignu = 200

	// granicaCzynnosciWsaduDesignu chroni przed zestawem stukrokowym: każdy krok
	// zakłada wariant, więc sto kroków na stu zasobach to dziesięć tysięcy plików.
	granicaCzynnosciWsaduDesignu = 24

	// granicaOgniwLancuchaDesignu jest górną granicą wykazu łańcucha edycji.
	granicaOgniwLancuchaDesignu = 500

	// domyslnaLiczbaBarwWektoryzacjiDesignu jest liczbą barw, na które rdzeń
	// rozkłada obraz przy obrysowaniu konturów.
	domyslnaLiczbaBarwWektoryzacjiDesignu = 6
)

// ZlozWarstwyFotografii składa warstwy rastrowe — obsługuje
// `design.photo.layer.composite`.
func (a *adapterDesignu) ZlozWarstwyFotografii(ctx context.Context,
	z shared.DesignPhotoLayerCompositeRequest) (shared.DesignPhotoLayerCompositeResponse, error) {

	if len(z.Layers) == 0 {
		return shared.DesignPhotoLayerCompositeResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.layer.composite bez ani jednej warstwy: kompozycja pusta nie " +
				"ma czego pokazać")
	}
	for numer, warstwa := range z.Layers {
		if warstwa.BlendMode != nil {
			if err := sprawdzWyliczenieDesignu("design.photo.layer.composite",
				fmt.Sprintf("layers[%d].blendMode", numer), *warstwa.BlendMode,
				shared.WartosciDesignPhotoBlendMode()); err != nil {
				return shared.DesignPhotoLayerCompositeResponse{}, err
			}
		}
		if warstwa.Opacity != nil && (*warstwa.Opacity < 0 || *warstwa.Opacity > 1) {
			return shared.DesignPhotoLayerCompositeResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"warstwa numer %d z kryciem %v: krycie idzie od 0 do 1", numer+1, *warstwa.Opacity))
		}
	}

	// Pierwsza warstwa rozstrzyga o rozmiarze płótna, gdy żądanie go nie podało:
	// jest spodem kompozycji, więc jej wymiar jest wymiarem obrazu.
	pierwszy, wierszPierwszej, err := a.obrazZasobuPoKodzieDesignu(ctx,
		"design.photo.layer.composite", z.Layers[0].AssetId)
	if err != nil {
		return shared.DesignPhotoLayerCompositeResponse{}, err
	}
	szerokosc, wysokosc := pierwszy.Bounds().Dx(), pierwszy.Bounds().Dy()
	if z.Width != nil && *z.Width > 0 {
		szerokosc = *z.Width
	}
	if z.Height != nil && *z.Height > 0 {
		wysokosc = *z.Height
	}
	if err := sprawdzRozmiarFotografiiDesignu(szerokosc, wysokosc); err != nil {
		return shared.DesignPhotoLayerCompositeResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.layer.composite: " + err.Error())
	}

	plotno := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	zlozonych := 0
	pominiete := []string{}
	for _, warstwa := range z.Layers {
		obraz, _, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.layer.composite",
			warstwa.AssetId)
		if err != nil {
			// Warstwa nie do odczytania NIE kończy kompozycji: wynik z pozostałych
			// warstw jest lepszy niż odmowa, a brak wraca w bilansie.
			pominiete = append(pominiete, fmt.Sprintf("%s (%s)",
				strings.TrimSpace(warstwa.AssetId), err.Error()))
			continue
		}
		obszar := obszarWarstwyFotografiiDesignu(warstwa, obraz, szerokosc, wysokosc)
		krycie := 1.0
		if warstwa.Opacity != nil {
			krycie = *warstwa.Opacity
		}
		tryb := shared.DesignPhotoBlendMode(shared.DesignPhotoBlendModeNormal)
		if warstwa.BlendMode != nil {
			tryb = *warstwa.BlendMode
		}
		var maska image.Image
		if warstwa.MaskAssetId != nil && strings.TrimSpace(*warstwa.MaskAssetId) != "" {
			odczytana, _, err := a.obrazZasobuPoKodzieDesignu(ctx,
				"design.photo.layer.composite", *warstwa.MaskAssetId)
			if err != nil {
				pominiete = append(pominiete, fmt.Sprintf("maska warstwy %s (%s)",
					strings.TrimSpace(warstwa.AssetId), err.Error()))
			} else {
				maska = odczytana
			}
		}
		zlozWarstwyFotografiiDesignu(plotno, obraz, obszar, krycie, tryb, maska)
		zlozonych++
	}
	if zlozonych == 0 {
		return shared.DesignPhotoLayerCompositeResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"żadna z %d warstw nie weszła: %s", len(z.Layers), strings.Join(pominiete, "; ")))
	}

	zasob, err := a.zapiszWariantFotografiiDesignu(ctx, wierszPierwszej, z.WindowId,
		shared.CommandDesignPhotoLayerComposite, "kompozycja warstw", plotno, z, nil)
	if err != nil {
		return shared.DesignPhotoLayerCompositeResponse{}, err
	}
	return shared.DesignPhotoLayerCompositeResponse{
		Asset: zasob, Composited: zlozonych,
		SkippedLayerAssetIds: uporzadkujBilansDesignu(pominiete),
	}, nil
}

// obszarWarstwyFotografiiDesignu rozstrzyga prostokąt, który warstwa zajmuje na
// płótnie. Brak wymiarów bierze wymiar obrazu; brak położenia bierze naroże.
func obszarWarstwyFotografiiDesignu(warstwa shared.DesignPhotoLayer, obraz image.Image,
	szerokoscPlotna, wysokoscPlotna int) image.Rectangle {

	x, y := 0, 0
	if warstwa.X != nil {
		x = int(*warstwa.X)
	}
	if warstwa.Y != nil {
		y = int(*warstwa.Y)
	}
	szerokosc, wysokosc := obraz.Bounds().Dx(), obraz.Bounds().Dy()
	if warstwa.Width != nil && *warstwa.Width > 0 {
		szerokosc = int(*warstwa.Width)
	}
	if warstwa.Height != nil && *warstwa.Height > 0 {
		wysokosc = int(*warstwa.Height)
	}
	return image.Rect(x, y, x+szerokosc, y+wysokosc).
		Intersect(image.Rect(0, 0, szerokoscPlotna, wysokoscPlotna))
}

// PuscWsadFotografii powtarza zestaw czynności na wielu zasobach — obsługuje
// `design.photo.batch.apply`.
func (a *adapterDesignu) PuscWsadFotografii(ctx context.Context,
	z shared.DesignPhotoBatchApplyRequest) (shared.DesignPhotoBatchApplyResponse, error) {

	if len(z.AssetIds) == 0 {
		return shared.DesignPhotoBatchApplyResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.batch.apply bez ani jednego zasobu")
	}
	if len(z.AssetIds) > granicaZasobowWsaduDesignu {
		return shared.DesignPhotoBatchApplyResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"wsad na %d zasobach przekracza granicę %d — każdy zasób przechodzi pełny rachunek "+
				"i zakłada nowe bajty w magazynie", len(z.AssetIds), granicaZasobowWsaduDesignu))
	}

	czynnosci := z.Operations
	if len(czynnosci) == 0 {
		if z.PresetId == nil || strings.TrimSpace(*z.PresetId) == "" {
			return shared.DesignPhotoBatchApplyResponse{}, bladWskazaniaDesignu(
				"komenda design.photo.batch.apply bez czynności i bez nastawy: rdzeń nie zgaduje, " +
					"co ma z tymi zasobami zrobić")
		}
		nastawa, err := a.repozytorium.NastawaFotografiiDesignuPoKodzie(ctx,
			strings.TrimSpace(*z.PresetId))
		if err != nil {
			return shared.DesignPhotoBatchApplyResponse{}, bladNieznanejNastawyFotografiiDesignu(
				*z.PresetId, err)
		}
		czynnosci, err = czynnosciZZapisuFotografiiDesignu(nastawa.CzynnosciJSON)
		if err != nil {
			return shared.DesignPhotoBatchApplyResponse{}, bladWydaniaDesignu(err.Error())
		}
	}
	if len(czynnosci) == 0 {
		return shared.DesignPhotoBatchApplyResponse{}, bladWskazaniaDesignu(
			"wskazana nastawa nie ma ani jednej czynności")
	}
	if len(czynnosci) > granicaCzynnosciWsaduDesignu {
		return shared.DesignPhotoBatchApplyResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"zestaw %d czynności przekracza granicę %d — każda czynność zakłada wariant, więc "+
				"zestaw razy zasoby to liczba plików w magazynie",
			len(czynnosci), granicaCzynnosciWsaduDesignu))
	}
	for numer, czynnosc := range czynnosci {
		if !czyKomendaWsaduFotografiiDesignu(czynnosc.Command) {
			return shared.DesignPhotoBatchApplyResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"czynność numer %d wskazuje komendę %q, której wsad nie puszcza; wsad przyjmuje: %s",
				numer+1, czynnosc.Command, strings.Join(komendyWsaduFotografiiDesignu(), ", ")))
		}
	}

	wyniki := []shared.DesignAsset{}
	nieudane := []string{}
	for _, kod := range z.AssetIds {
		kod = strings.TrimSpace(kod)
		biezacy := kod
		udalo := true
		for numer, czynnosc := range czynnosci {
			nastepny, err := a.puscCzynnoscFotografiiDesignu(ctx, biezacy, czynnosc, z.WindowId)
			if err != nil {
				nieudane = append(nieudane, fmt.Sprintf("%s (czynność %d, %s: %s)",
					kod, numer+1, czynnosc.Command, err.Error()))
				udalo = false
				break
			}
			biezacy = nastepny.Id
			// Wynik ostatniej czynności jest wynikiem wsadu dla tego zasobu; wyniki
			// pośrednie zostają w magazynie jako ogniwa łańcucha edycji i nie wchodzą
			// do odpowiedzi — inaczej wykaz `assets` mieszałby ogniwa z wynikami.
			if numer == len(czynnosci)-1 {
				wyniki = append(wyniki, nastepny)
			}
		}
		_ = udalo
	}
	if len(wyniki) == 0 {
		return shared.DesignPhotoBatchApplyResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"żaden z %d zasobów nie przeszedł wsadu do końca: %s",
			len(z.AssetIds), strings.Join(nieudane, "; ")))
	}
	return shared.DesignPhotoBatchApplyResponse{
		Assets: wyniki, Applied: len(wyniki), FailedAssetIds: uporzadkujBilansDesignu(nieudane),
	}, nil
}

// komendyWsaduFotografiiDesignu wylicza komendy, które wsad umie puścić.
//
// Wykaz jest wąski z zamysłu: wsad puszcza czynności działające na JEDNYM
// obrazie i niepotrzebujące wskazania drugiego zasobu. Kompozycja warstw, maska
// i zaznaczenie obiektu wymagają wskazań, których wsad nie ma skąd wziąć dla
// każdego zasobu osobno — przepuszczenie ich dałoby ten sam obszar albo tę samą
// maskę nałożoną na sto różnych zdjęć.
func komendyWsaduFotografiiDesignu() []string {
	return []string{
		shared.CommandDesignPhotoCrop,
		shared.CommandDesignPhotoTransform,
		shared.CommandDesignPhotoResample,
		shared.CommandDesignPhotoUpscale,
		shared.CommandDesignPhotoEnhance,
		shared.CommandDesignPhotoColorCorrect,
		shared.CommandDesignPhotoFilterApply,
		shared.CommandDesignPhotoBackgroundRemove,
	}
}

// czyKomendaWsaduFotografiiDesignu rozstrzyga, czy wsad puszcza tę komendę.
func czyKomendaWsaduFotografiiDesignu(komenda string) bool {
	for _, znana := range komendyWsaduFotografiiDesignu() {
		if znana == komenda {
			return true
		}
	}
	return false
}

// puscCzynnoscFotografiiDesignu wykonuje jedną czynność wsadu na wskazanym
// zasobie, przez ten sam uchwyt, którym jedzie czynność pojedyncza.
//
// Nastawy czynności są przekładane na żądanie komendy przez rozbiór JSON razem
// z dołożonym wskazaniem zasobu. Dzięki temu nastawa zapisana w oknie i nastawa
// puszczona wsadem to DOKŁADNIE ten sam zestaw pól — nie ma tu drugiego zapisu
// nastaw, który mógłby się rozjechać z pierwszym.
func (a *adapterDesignu) puscCzynnoscFotografiiDesignu(ctx context.Context, zasob string,
	czynnosc shared.DesignPhotoOperation, okno *string) (shared.DesignAsset, error) {

	nastawy := map[string]any{}
	if len(czynnosc.Settings) > 0 {
		if err := json.Unmarshal(czynnosc.Settings, &nastawy); err != nil {
			return shared.DesignAsset{}, fmt.Errorf(
				"nastawy czynności nie są obiektem JSON: %w", err)
		}
	}
	nastawy["assetId"] = zasob
	if okno != nil && strings.TrimSpace(*okno) != "" {
		nastawy["windowId"] = strings.TrimSpace(*okno)
	}
	zapis, err := json.Marshal(nastawy)
	if err != nil {
		return shared.DesignAsset{}, fmt.Errorf("nie można złożyć żądania czynności: %w", err)
	}

	switch czynnosc.Command {
	case shared.CommandDesignPhotoCrop:
		var zadanie shared.DesignPhotoCropRequest
		if err := json.Unmarshal(zapis, &zadanie); err != nil {
			return shared.DesignAsset{}, err
		}
		odpowiedz, err := a.Kadruj(ctx, zadanie)
		return odpowiedz.Asset, err
	case shared.CommandDesignPhotoTransform:
		var zadanie shared.DesignPhotoTransformRequest
		if err := json.Unmarshal(zapis, &zadanie); err != nil {
			return shared.DesignAsset{}, err
		}
		odpowiedz, err := a.Przeksztalc(ctx, zadanie)
		return odpowiedz.Asset, err
	case shared.CommandDesignPhotoResample:
		var zadanie shared.DesignPhotoResampleRequest
		if err := json.Unmarshal(zapis, &zadanie); err != nil {
			return shared.DesignAsset{}, err
		}
		odpowiedz, err := a.PrzeliczRozdzielczosc(ctx, zadanie)
		return odpowiedz.Asset, err
	case shared.CommandDesignPhotoUpscale:
		var zadanie shared.DesignPhotoUpscaleRequest
		if err := json.Unmarshal(zapis, &zadanie); err != nil {
			return shared.DesignAsset{}, err
		}
		odpowiedz, err := a.Powieksz(ctx, zadanie)
		return odpowiedz.Asset, err
	case shared.CommandDesignPhotoEnhance:
		var zadanie shared.DesignPhotoEnhanceRequest
		if err := json.Unmarshal(zapis, &zadanie); err != nil {
			return shared.DesignAsset{}, err
		}
		odpowiedz, err := a.PopraweJakosc(ctx, zadanie)
		return odpowiedz.Asset, err
	case shared.CommandDesignPhotoColorCorrect:
		var zadanie shared.DesignPhotoColorCorrectRequest
		if err := json.Unmarshal(zapis, &zadanie); err != nil {
			return shared.DesignAsset{}, err
		}
		odpowiedz, err := a.SkorygujBarwe(ctx, zadanie)
		return odpowiedz.Asset, err
	case shared.CommandDesignPhotoFilterApply:
		var zadanie shared.DesignPhotoFilterApplyRequest
		if err := json.Unmarshal(zapis, &zadanie); err != nil {
			return shared.DesignAsset{}, err
		}
		odpowiedz, err := a.NalozFiltr(ctx, zadanie)
		return odpowiedz.Asset, err
	case shared.CommandDesignPhotoBackgroundRemove:
		var zadanie shared.DesignPhotoBackgroundRemoveRequest
		if err := json.Unmarshal(zapis, &zadanie); err != nil {
			return shared.DesignAsset{}, err
		}
		odpowiedz, err := a.OdetnijTlo(ctx, zadanie)
		return odpowiedz.Asset, err
	}
	return shared.DesignAsset{}, fmt.Errorf("wsad nie puszcza komendy %q", czynnosc.Command)
}

// ZapiszNastaweFotografii zapisuje zestaw czynności pod nazwą — obsługuje
// `design.photo.preset.save`.
func (a *adapterDesignu) ZapiszNastaweFotografii(ctx context.Context,
	z shared.DesignPhotoPresetSaveRequest) (shared.DesignPhotoPresetSaveResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignPhotoPresetSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.preset.save bez wskazania okna")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignPhotoPresetSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.preset.save bez nazwy nastawy")
	}
	if len(z.Operations) == 0 {
		return shared.DesignPhotoPresetSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.preset.save bez ani jednej czynności: nastawa pusta nie " +
				"opisuje żadnej pracy")
	}
	for numer, czynnosc := range z.Operations {
		if !czyKomendaWsaduFotografiiDesignu(czynnosc.Command) {
			return shared.DesignPhotoPresetSaveResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"czynność numer %d wskazuje komendę %q, której wsad nie puszcza — nastawa "+
					"z taką czynnością byłaby nastawą nie do użycia; komendy: %s",
				numer+1, czynnosc.Command, strings.Join(komendyWsaduFotografiiDesignu(), ", ")))
		}
	}

	kod := nowyIdentyfikator(przedrostekNastawyFotografiiDesign)
	if z.PresetId != nil && strings.TrimSpace(*z.PresetId) != "" {
		kod = strings.TrimSpace(*z.PresetId)
		zastana, err := a.repozytorium.NastawaFotografiiDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignPhotoPresetSaveResponse{}, bladNieznanejNastawyFotografiiDesignu(
				kod, err)
		}
		if zastana.Okno != z.WindowId {
			return shared.DesignPhotoPresetSaveResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"nastawa %s należy do okna %s, a komenda przyszła z okna %s",
				kod, zastana.Okno, z.WindowId))
		}
	}

	zapis, err := json.Marshal(z.Operations)
	if err != nil {
		return shared.DesignPhotoPresetSaveResponse{}, bladWydaniaDesignu(
			"nie można złożyć zapisu czynności nastawy: " + err.Error())
	}
	zapisana, err := a.repozytorium.ZapiszNastaweFotografiiDesignu(ctx,
		dane.NastawaFotografiiDesignu{
			Kod: kod, Okno: z.WindowId, Nazwa: strings.TrimSpace(z.Name),
			CzynnosciJSON: string(zapis),
		})
	if err != nil {
		return shared.DesignPhotoPresetSaveResponse{}, bladDesignu(err)
	}
	nastawa, err := nastawaKontraktuFotografiiDesignu(zapisana)
	if err != nil {
		return shared.DesignPhotoPresetSaveResponse{}, bladWydaniaDesignu(err.Error())
	}
	return shared.DesignPhotoPresetSaveResponse{Preset: nastawa}, nil
}

// NastawyFotografii zwraca nastawy okna — obsługuje `design.photo.preset.list`.
func (a *adapterDesignu) NastawyFotografii(ctx context.Context,
	z shared.DesignPhotoPresetListRequest) (shared.DesignPhotoPresetListResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignPhotoPresetListResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.preset.list bez wskazania okna")
	}
	wiersze, err := a.repozytorium.NastawyFotografiiDesignu(ctx, z.WindowId)
	if err != nil {
		return shared.DesignPhotoPresetListResponse{}, bladDesignu(err)
	}
	nastawy := make([]shared.DesignPhotoPreset, 0, len(wiersze))
	for _, wiersz := range wiersze {
		nastawa, err := nastawaKontraktuFotografiiDesignu(wiersz)
		if err != nil {
			return shared.DesignPhotoPresetListResponse{}, bladWydaniaDesignu(err.Error())
		}
		nastawy = append(nastawy, nastawa)
	}
	return shared.DesignPhotoPresetListResponse{Presets: nastawy, Total: len(nastawy)}, nil
}

// nastawaKontraktuFotografiiDesignu składa `DesignPhotoPreset` kontraktu
// z wiersza.
func nastawaKontraktuFotografiiDesignu(
	wiersz dane.NastawaFotografiiDesignu) (shared.DesignPhotoPreset, error) {

	czynnosci, err := czynnosciZZapisuFotografiiDesignu(wiersz.CzynnosciJSON)
	if err != nil {
		return shared.DesignPhotoPreset{}, err
	}
	chwila := chwilaBazy(wiersz.Zaktualizowano)
	return shared.DesignPhotoPreset{
		Id: wiersz.Kod, WindowId: wiersz.Okno, Name: wiersz.Nazwa,
		Operations: czynnosci, UpdatedAt: &chwila,
	}, nil
}

// czynnosciZZapisuFotografiiDesignu rozkłada zapis czynności z kolumny.
func czynnosciZZapisuFotografiiDesignu(zapis string) ([]shared.DesignPhotoOperation, error) {
	var czynnosci []shared.DesignPhotoOperation
	if err := json.Unmarshal([]byte(zapis), &czynnosci); err != nil {
		return nil, fmt.Errorf(
			"zapis czynności nastawy w bazie nie jest wykazem czynności: %w", err)
	}
	return czynnosci, nil
}

// ObrysujKontury zamienia raster w rysunek wektorowy — obsługuje
// `design.photo.vectorize`.
func (a *adapterDesignu) ObrysujKontury(ctx context.Context,
	z shared.DesignPhotoVectorizeRequest) (shared.DesignPhotoVectorizeResponse, error) {

	barw := domyslnaLiczbaBarwWektoryzacjiDesignu
	if z.Colors != nil {
		if *z.Colors < 2 || *z.Colors > 64 {
			return shared.DesignPhotoVectorizeResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.photo.vectorize z %d barwami: rdzeń rozkłada obraz na od 2 do 64 "+
					"barw — jedna barwa nie daje rysunku, a powyżej sześćdziesięciu czterech "+
					"rysunek przestaje być rysunkiem", *z.Colors))
		}
		barw = *z.Colors
	}
	prog := 0.0
	if z.Threshold != nil {
		if *z.Threshold < 0 || *z.Threshold > 1 {
			return shared.DesignPhotoVectorizeResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.photo.vectorize z progiem %v: próg idzie od 0 do 1", *z.Threshold))
		}
		prog = *z.Threshold
	}
	wygladzenie := 0.0
	if z.Smoothing != nil {
		if *z.Smoothing < 0 {
			return shared.DesignPhotoVectorizeResponse{}, bladWskazaniaDesignu(
				"komenda design.photo.vectorize z ujemnym wygładzeniem")
		}
		wygladzenie = *z.Smoothing
	}

	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.photo.vectorize", z.AssetId)
	if err != nil {
		return shared.DesignPhotoVectorizeResponse{}, err
	}
	dokument, sciezek, odrzuconych, err := obrysujKonturyDesignu(obraz, barw, prog, wygladzenie)
	if err != nil {
		return shared.DesignPhotoVectorizeResponse{}, bladWskazaniaDesignu(
			"komenda design.photo.vectorize: " + err.Error())
	}

	// Wynik jest WEKTOREM, więc nie idzie drogą wariantu rastrowego: format i
	// rodzaj zasobu są inne, a `zapiszWariantFotografiiDesignu` koduje PNG.
	nazwa := nazwaZasobuDesignu(zrodlo) + " — obrysowanie konturów"
	zapisany, err := a.zalozZasobZBajtowDesignu(ctx, oknoWytworu(z.WindowId, zrodlo.Okno),
		nazwa, shared.DesignAssetKindVector, "svg", []byte(dokument))
	if err != nil {
		return shared.DesignPhotoVectorizeResponse{}, err
	}
	zapisany.WariantZasobuID = &zrodlo.Kod
	zapisany, err = a.repozytorium.ZapiszZasob(ctx, zapisany)
	if err != nil {
		return shared.DesignPhotoVectorizeResponse{}, bladDesignu(err)
	}
	nastawy, _ := json.Marshal(z)
	tresc := string(nastawy)
	_ = a.repozytorium.ZapiszCzynnoscFotografiiDesignu(ctx, dane.CzynnoscFotografiiDesignu{
		ZasobID: zapisany.ID, ZasobZrodlaID: &zrodlo.ID,
		Komenda: shared.CommandDesignPhotoVectorize, NastawyJSON: &tresc,
	})
	return shared.DesignPhotoVectorizeResponse{
		Asset: zasobWytworzonyKontraktu(zapisany), PathCount: sciezek,
		DroppedRegions: &odrzuconych,
	}, nil
}

// MetadaneZasobu oddaje zmierzone właściwości zasobu — obsługuje
// `design.photo.metadata.get`.
func (a *adapterDesignu) MetadaneZasobu(ctx context.Context,
	z shared.DesignPhotoMetadataGetRequest) (shared.DesignPhotoMetadataGetResponse, error) {

	zasob, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(z.AssetId))
	if err != nil {
		return shared.DesignPhotoMetadataGetResponse{}, bladNieznanegoZasobuDesignu(z.AssetId, err)
	}
	if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
		return shared.DesignPhotoMetadataGetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"zasób %s nie ma treści w magazynie — nie ma czego zmierzyć", zasob.Kod))
	}
	bajty, err := os.ReadFile(*zasob.URI)
	if err != nil {
		return shared.DesignPhotoMetadataGetResponse{}, bladDesignu(fmt.Errorf(
			"treści zasobu %s nie ma pod jego odwołaniem w magazynie: %w", zasob.Kod, err))
	}

	metadane := shared.DesignPhotoMetadata{}
	wielkosc := len(bajty)
	metadane.SizeBytes = &wielkosc

	// Wymiary, format i model barwny mierzymy z NAGŁÓWKA — `image.DecodeConfig`
	// nie rozpakowuje obrazu, więc pomiar kosztuje kilkadziesiąt bajtów odczytu.
	if opis, format, err := image.DecodeConfig(bytes.NewReader(bajty)); err == nil {
		szerokosc, wysokosc := opis.Width, opis.Height
		metadane.Width, metadane.Height = &szerokosc, &wysokosc
		nazwaFormatu := format
		metadane.Format = &nazwaFormatu
		model := nazwaModeluBarwnegoDesignu(opis)
		metadane.ColorModel = &model
		zKryciem := czyModelZKryciemDesignu(opis)
		metadane.HasAlpha = &zKryciem
	} else if zasob.Format != nil {
		// Plik, którego rdzeń nie rozkłada (SVG, dokument), nadal ma format
		// zmierzony przy wniesieniu. Pusto tam, gdzie czegoś nie wiemy, ale nie
		// zapominamy tego, co już wiedzieliśmy.
		metadane.Format = zasob.Format
	}

	exif, pol := odczytajExifDesignu(bajty)
	if pol > 0 {
		metadane.Exif = exif
	}
	metadane.ExifFieldsRead = &pol
	if rozdzielczosc, jest := rozdzielczoscZExifDesignu(exif); jest {
		metadane.Dpi = &rozdzielczosc
	}
	return shared.DesignPhotoMetadataGetResponse{Metadata: metadane}, nil
}

// nazwaModeluBarwnegoDesignu nazywa model barwny odczytany z nagłówka.
//
// Rozpoznanie idzie po TOŻSAMOŚCI modelu biblioteki standardowej, nie po nazwie
// typu: modele są wartościami jednostkowymi pakietu `image/color`, więc
// porównanie jest dokładne i nie łamie się przy zmianie nazw wewnętrznych.
func nazwaModeluBarwnegoDesignu(opis image.Config) string {
	switch opis.ColorModel {
	case color.GrayModel, color.Gray16Model:
		return "szarosc"
	case color.CMYKModel:
		return "cmyk"
	case color.NRGBAModel, color.RGBAModel, color.NRGBA64Model, color.RGBA64Model:
		return "rgb-z-kryciem"
	case color.YCbCrModel:
		return "ycbcr"
	}
	return "rgb"
}

// czyModelZKryciemDesignu rozstrzyga, czy nagłówek pliku zapowiada kanał krycia.
//
// To POMIAR zapowiedzi, nie treści: plik PNG z kanałem krycia wypełnionym
// wszędzie pełną wartością nadal jest plikiem z kanałem krycia, i tak ma być
// powiedziane. Czy krycie naprawdę jest gdzieś częściowe, mówi pole
// `transparentShare` przy odcięciu tła.
func czyModelZKryciemDesignu(opis image.Config) bool {
	switch opis.ColorModel {
	case color.NRGBAModel, color.RGBAModel, color.NRGBA64Model, color.RGBA64Model,
		color.AlphaModel, color.Alpha16Model:
		return true
	}
	return false
}

// LancuchEdycji oddaje łańcuch edycji zasobu — obsługuje
// `design.photo.history.get`.
//
// Łańcuch idzie WSTECZ po wariantach: od wskazanego zasobu do zdjęcia, którego
// nikt nie obrabiał. Każdy krok niesie czynność i jej nastawy odczytane z bazy,
// a nie odtworzone z różnicy obrazów — różnica obrazów nie powiedziałaby, jakimi
// suwakami Operator do niej doszedł.
func (a *adapterDesignu) LancuchEdycji(ctx context.Context,
	z shared.DesignPhotoHistoryGetRequest) (shared.DesignPhotoHistoryGetResponse, error) {

	granica := granicaOgniwLancuchaDesignu
	if z.Limit != nil {
		if *z.Limit < 1 || *z.Limit > granicaOgniwLancuchaDesignu {
			return shared.DesignPhotoHistoryGetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.photo.history.get z granicą %d: rdzeń oddaje od 1 do %d ogniw",
				*z.Limit, granicaOgniwLancuchaDesignu))
		}
		granica = *z.Limit
	}
	zasob, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(z.AssetId))
	if err != nil {
		return shared.DesignPhotoHistoryGetResponse{}, bladNieznanegoZasobuDesignu(z.AssetId, err)
	}

	ogniwa := []shared.DesignPhotoEdit{}
	biezacy := zasob
	// Granica przejścia jest podwójna: liczba ogniw i granica bezpieczeństwa na
	// wypadek pętli wariantów (zasób wskazujący sam siebie jako wariant). Pętla
	// w danych nie ma prawa zawiesić odczytu.
	for krokow := 0; krokow < granicaOgniwLancuchaDesignu && len(ogniwa) < granica; krokow++ {
		czynnosc, err := a.repozytorium.CzynnoscFotografiiDesignuWyniku(ctx, biezacy.ID)
		if err != nil {
			if czyBrakZasobuDesignu(err) {
				break
			}
			return shared.DesignPhotoHistoryGetResponse{}, bladDesignu(err)
		}
		ogniwo := shared.DesignPhotoEdit{
			AssetId: biezacy.Kod, Command: czynnosc.Komenda,
			SourceAssetId: biezacy.WariantZasobuID,
		}
		if czynnosc.NastawyJSON != nil {
			ogniwo.Settings = []byte(*czynnosc.NastawyJSON)
		}
		if czynnosc.PoliczonePrzez != nil {
			droga := shared.DesignPhotoComputeRoute(*czynnosc.PoliczonePrzez)
			ogniwo.ComputedBy = &droga
		}
		chwila := chwilaBazy(czynnosc.Utworzono)
		ogniwo.CreatedAt = &chwila
		ogniwa = append(ogniwa, ogniwo)

		if biezacy.WariantZasobuID == nil || strings.TrimSpace(*biezacy.WariantZasobuID) == "" ||
			*biezacy.WariantZasobuID == biezacy.Kod {
			break
		}
		poprzedni, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(*biezacy.WariantZasobuID))
		if err != nil {
			// Źródło usunięte z panelu przerywa łańcuch, ale nie unieważnia tego, co
			// już zebrano: Operator ma widzieć ogniwa, które ocalały.
			break
		}
		biezacy = poprzedni
	}

	// Ogniwa idą od NAJSTARSZEGO — tak opisuje pole kontrakt i tak czyta się
	// historię pracy. Przejście szło wstecz, więc wykaz odwracamy.
	for lewa, prawa := 0, len(ogniwa)-1; lewa < prawa; lewa, prawa = lewa+1, prawa-1 {
		ogniwa[lewa], ogniwa[prawa] = ogniwa[prawa], ogniwa[lewa]
	}
	if len(ogniwa) == 0 {
		return shared.DesignPhotoHistoryGetResponse{Edits: ogniwa, Total: 0}, nil
	}
	return shared.DesignPhotoHistoryGetResponse{Edits: ogniwa, Total: len(ogniwa)}, nil
}

// odczytajExifDesignu odczytuje pola EXIF z bajtów pliku i oddaje je wraz
// z liczbą pól odczytanych.
//
// Rozbiór jest własny i wąski: rdzeń szuka segmentu APP1 pliku JPEG, czyta
// katalog IFD0 i wyciąga z niego pola o znanych numerach. Biblioteki EXIF w tym
// drzewie nie ma, a pełny rozbiór wszystkich katalogów (IFD1, GPS, Interop,
// MakerNote każdego producenta) jest zadaniem na osobną bibliotekę — nie na plik
// modułu.
//
// Liczba odczytanych pól wraca w odpowiedzi (`exifFieldsRead`): zero znaczy, że
// plik EXIF-u nie ma albo że rdzeń go nie rozłożył, i jest to POWIEDZIANE, a nie
// przemilczane pustym obiektem.
func odczytajExifDesignu(bajty []byte) ([]byte, int) {
	segment := segmentExifDesignu(bajty)
	if len(segment) < 8 {
		return nil, 0
	}
	// Porządek bajtów: „II" znaczy od najmniej znaczącego, „MM" od najbardziej.
	var porzadek binary.ByteOrder
	switch {
	case segment[0] == 'I' && segment[1] == 'I':
		porzadek = binary.LittleEndian
	case segment[0] == 'M' && segment[1] == 'M':
		porzadek = binary.BigEndian
	default:
		return nil, 0
	}
	if porzadek.Uint16(segment[2:4]) != 42 {
		return nil, 0
	}
	przesuniecie := porzadek.Uint32(segment[4:8])
	if int(przesuniecie)+2 > len(segment) {
		return nil, 0
	}
	pol := int(porzadek.Uint16(segment[przesuniecie : przesuniecie+2]))
	wpisy := int(przesuniecie) + 2

	odczytane := map[string]any{}
	for numer := 0; numer < pol; numer++ {
		poczatek := wpisy + numer*12
		if poczatek+12 > len(segment) {
			break
		}
		znacznik := porzadek.Uint16(segment[poczatek : poczatek+2])
		rodzaj := porzadek.Uint16(segment[poczatek+2 : poczatek+4])
		liczba := porzadek.Uint32(segment[poczatek+4 : poczatek+8])
		wartosc := segment[poczatek+8 : poczatek+12]
		nazwa, znany := nazwyExifDesignu[znacznik]
		if !znany {
			continue
		}
		odczyt, jest := wartoscExifDesignu(segment, porzadek, rodzaj, liczba, wartosc)
		if !jest {
			continue
		}
		odczytane[nazwa] = odczyt
	}
	if len(odczytane) == 0 {
		return nil, 0
	}
	zapis, err := json.Marshal(odczytane)
	if err != nil {
		return nil, 0
	}
	return zapis, len(odczytane)
}

// nazwyExifDesignu to numery pól EXIF, które rdzeń rozpoznaje.
//
// Wykaz jest krótki z zamysłu: to pola, o które pyta się przy grafice użytkowej —
// aparat, obiektyw, czas, przysłona, czułość, rozdzielczość i orientacja. Pełny
// wykaz EXIF ma kilkaset pól, a większość z nich nie ma w tym module znaczenia.
var nazwyExifDesignu = map[uint16]string{
	0x010E: "opis",
	0x010F: "producent",
	0x0110: "model",
	0x0112: "orientacja",
	0x011A: "rozdzielczoscX",
	0x011B: "rozdzielczoscY",
	0x0128: "jednostkaRozdzielczosci",
	0x0131: "oprogramowanie",
	0x0132: "chwila",
	0x013B: "autor",
	0x8298: "prawa",
	0x829A: "czasEkspozycji",
	0x829D: "przyslona",
	0x8827: "czulosc",
	0x920A: "ogniskowa",
}

// wartoscExifDesignu odczytuje wartość jednego pola EXIF.
func wartoscExifDesignu(segment []byte, porzadek binary.ByteOrder, rodzaj uint16,
	liczba uint32, wartosc []byte) (any, bool) {

	// Wartości dłuższe niż cztery bajty leżą pod przesunięciem zapisanym w polu.
	dlugosc := map[uint16]uint32{1: 1, 2: 1, 3: 2, 4: 4, 5: 8, 7: 1, 9: 4, 10: 8}[rodzaj]
	if dlugosc == 0 {
		return nil, false
	}
	dane := wartosc
	if dlugosc*liczba > 4 {
		przesuniecie := porzadek.Uint32(wartosc)
		koniec := przesuniecie + dlugosc*liczba
		if int(koniec) > len(segment) {
			return nil, false
		}
		dane = segment[przesuniecie:koniec]
	}
	switch rodzaj {
	case 2:
		return strings.TrimRight(string(dane), "\x00"), true
	case 3:
		if len(dane) < 2 {
			return nil, false
		}
		return int(porzadek.Uint16(dane[:2])), true
	case 4:
		if len(dane) < 4 {
			return nil, false
		}
		return int(porzadek.Uint32(dane[:4])), true
	case 5, 10:
		if len(dane) < 8 {
			return nil, false
		}
		licznik := porzadek.Uint32(dane[:4])
		mianownik := porzadek.Uint32(dane[4:8])
		if mianownik == 0 {
			return nil, false
		}
		return float64(licznik) / float64(mianownik), true
	}
	return nil, false
}

// segmentExifDesignu odnajduje treść segmentu APP1 z podpisem „Exif" w pliku
// JPEG.
func segmentExifDesignu(bajty []byte) []byte {
	if len(bajty) < 4 || bajty[0] != 0xFF || bajty[1] != 0xD8 {
		return nil
	}
	numer := 2
	for numer+4 <= len(bajty) {
		if bajty[numer] != 0xFF {
			numer++
			continue
		}
		znacznik := bajty[numer+1]
		if znacznik == 0xD8 || znacznik == 0x01 || (znacznik >= 0xD0 && znacznik <= 0xD7) {
			numer += 2
			continue
		}
		if znacznik == 0xDA || znacznik == 0xD9 {
			// Początek danych obrazu albo koniec pliku: dalej EXIF-u nie ma.
			return nil
		}
		dlugosc := int(binary.BigEndian.Uint16(bajty[numer+2 : numer+4]))
		if numer+2+dlugosc > len(bajty) {
			return nil
		}
		if znacznik == 0xE1 && dlugosc > 8 {
			tresc := bajty[numer+4 : numer+2+dlugosc]
			if len(tresc) > 6 && string(tresc[:4]) == "Exif" {
				return tresc[6:]
			}
		}
		numer += 2 + dlugosc
	}
	return nil
}

// rozdzielczoscZExifDesignu wyciąga rozdzielczość z odczytanych pól EXIF.
func rozdzielczoscZExifDesignu(zapis []byte) (int, bool) {
	if len(zapis) == 0 {
		return 0, false
	}
	var odczytane map[string]any
	if err := json.Unmarshal(zapis, &odczytane); err != nil {
		return 0, false
	}
	wartosc, jest := odczytane["rozdzielczoscX"]
	if !jest {
		return 0, false
	}
	liczba, jest := wartosc.(float64)
	if !jest || liczba <= 0 {
		return 0, false
	}
	// Jednostka 3 znaczy punkty na centymetr; przeliczamy na cale, bo w tym module
	// wszystko inne liczy się w calach (`milimetryNaCal`).
	if jednostka, jest := odczytane["jednostkaRozdzielczosci"].(float64); jest && jednostka == 3 {
		liczba *= 2.54
	}
	return int(liczba + 0.5), true
}

// bladNieznanejNastawyFotografiiDesignu nazywa nastawę, której rdzeń nie zna.
func bladNieznanejNastawyFotografiiDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("nastawy warsztatu fotografii " + kod +
			" nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}
