// Wpięcie ośmiu komend obszaru `design.*` — modułu Design (Prompt Builder,
// Assets Panel, Design Board) — i rozgłoszenie `design.asset.changed` po tych
// z nich, które zasób zmieniają; pozostałe czynności wpinają funkcje
// pomocnicze niżej w pliku.
package core

import (
	"context"
	"strings"

	"danacoconsole/shared"
)

// Design jest portem modułu Design: interfejsem czynności, które adapter
// warstwy danych musi spełnić dla rejestru komend.
type Design interface {
	GenerujZasob(ctx context.Context, z shared.DesignAssetGenerateRequest) (shared.DesignAssetGenerateResponse, error)
	Zasoby(ctx context.Context, z shared.DesignAssetListRequest) (shared.DesignAssetListResponse, error)
	ZapiszKompozycje(ctx context.Context, z shared.DesignBoardUpdateRequest) (shared.DesignBoardUpdateResponse, error)
	UstawEtykietyZasobu(ctx context.Context, z shared.DesignAssetTagSetRequest) (shared.DesignAssetTagSetResponse, error)
	WniesZasob(ctx context.Context, z shared.DesignAssetUploadRequest) (shared.DesignAssetUploadResponse, error)
	UstawUlubionyZasob(ctx context.Context, z shared.DesignAssetFavoriteSetRequest) (shared.DesignAssetFavoriteSetResponse, error)
	// UsunZasob oddaje trzy wartości, jedyna w tym porcie: odpowiedź, zasób usunięty i błąd.
	UsunZasob(ctx context.Context, z shared.DesignAssetRemoveRequest) (shared.DesignAssetRemoveResponse, shared.DesignAsset, error)
	Kompozycje(ctx context.Context, z shared.DesignBoardListRequest) (shared.DesignBoardListResponse, error)

	// --- treść i wydania zasobu ---
	TrescZasobu(ctx context.Context, z shared.DesignAssetContentGetRequest) (shared.DesignAssetContentGetResponse, error)
	WydajZasob(ctx context.Context, z shared.DesignAssetExportRequest) (shared.DesignAssetExportResponse, error)
	WydajZasobyPartia(ctx context.Context, z shared.DesignAssetExportBatchRequest) (shared.DesignAssetExportBatchResponse, error)

	// --- kolekcje ---
	ZalozKolekcje(ctx context.Context, z shared.DesignCollectionCreateRequest) (shared.DesignCollectionCreateResponse, error)
	PrzypiszDoKolekcji(ctx context.Context, z shared.DesignCollectionAssignRequest) (shared.DesignCollectionAssignResponse, error)
	Kolekcje(ctx context.Context, z shared.DesignCollectionListRequest) (shared.DesignCollectionListResponse, error)

	// --- szablony promptu i historia ---
	ZapiszSzablonPromptu(ctx context.Context, z shared.DesignPromptTemplateSaveRequest) (shared.DesignPromptTemplateSaveResponse, error)
	SzablonyPromptu(ctx context.Context, z shared.DesignPromptTemplateListRequest) (shared.DesignPromptTemplateListResponse, error)
	HistoriaPromptow(ctx context.Context, z shared.DesignPromptHistoryListRequest) (shared.DesignPromptHistoryListResponse, error)

	// --- wersje kompozycji i wyrys ---
	ZapiszWersjeKompozycji(ctx context.Context, z shared.DesignBoardVersionSaveRequest) (shared.DesignBoardVersionSaveResponse, error)
	WersjeKompozycji(ctx context.Context, z shared.DesignBoardVersionListRequest) (shared.DesignBoardVersionListResponse, error)
	PrzywrocWersjeKompozycji(ctx context.Context, z shared.DesignBoardVersionRestoreRequest) (shared.DesignBoardVersionRestoreResponse, error)
	WyrysujKompozycje(ctx context.Context, z shared.DesignBoardExportRequest) (shared.DesignBoardExportResponse, error)

	// --- adnotacje i obecność ---
	UstawAdnotacje(ctx context.Context, z shared.DesignAnnotationSetRequest) (shared.DesignAnnotationSetResponse, error)
	Adnotacje(ctx context.Context, z shared.DesignAnnotationListRequest) (shared.DesignAnnotationListResponse, error)
	ZglosObecnosc(ctx context.Context, z shared.DesignPresenceReportRequest) (shared.DesignPresenceReportResponse, error)

	// --- zestawy żetonów i przewodnik stylu ---
	ZapiszZestawZetonow(ctx context.Context, z shared.DesignTokensetSaveRequest) (shared.DesignTokensetSaveResponse, error)
	ZestawyZetonow(ctx context.Context, z shared.DesignTokensetListRequest) (shared.DesignTokensetListResponse, error)
	WydajZestawZetonow(ctx context.Context, z shared.DesignTokensetExportRequest) (shared.DesignTokensetExportResponse, error)
	WczytajZestawZetonow(ctx context.Context, z shared.DesignTokensetImportRequest) (shared.DesignTokensetImportResponse, error)
	WydajPrzewodnikStylu(ctx context.Context, z shared.DesignStyleguidePublishRequest) (shared.DesignStyleguidePublishResponse, error)

	// --- warsztat wektorowy (adapter_modul_design_wektor.go) ---
	UstawSciezke(ctx context.Context, z shared.DesignVectorPathSetRequest) (shared.DesignVectorPathSetResponse, error)
	SciezkiWektorowe(ctx context.Context, z shared.DesignVectorPathListRequest) (shared.DesignVectorPathListResponse, error)
	UsunSciezke(ctx context.Context, z shared.DesignVectorPathRemoveRequest) (shared.DesignVectorPathRemoveResponse, error)
	DolozKsztalt(ctx context.Context, z shared.DesignVectorShapeAddRequest) (shared.DesignVectorShapeAddResponse, error)
	ZlozSciezkiLogicznie(ctx context.Context, z shared.DesignVectorBooleanRequest) (shared.DesignVectorBooleanResponse, error)
	TekstNaSciezce(ctx context.Context, z shared.DesignVectorTextPathRequest) (shared.DesignVectorTextPathResponse, error)
	OczyscSciezki(ctx context.Context, z shared.DesignVectorOptimizeRequest) (shared.DesignVectorOptimizeResponse, error)
	UstawSymbol(ctx context.Context, z shared.DesignVectorSymbolSetRequest) (shared.DesignVectorSymbolSetResponse, error)
	Symbole(ctx context.Context, z shared.DesignVectorSymbolListRequest) (shared.DesignVectorSymbolListResponse, error)
	WydajWektor(ctx context.Context, z shared.DesignVectorExportRequest) (shared.DesignVectorExportResponse, error)

	// --- warsztat makiety (adapter_modul_design_makiety.go, _makiety_zrzut.go) ---
	UstawRamke(ctx context.Context, z shared.DesignFrameSetRequest) (shared.DesignFrameSetResponse, error)
	Ramki(ctx context.Context, z shared.DesignFrameListRequest) (shared.DesignFrameListResponse, error)
	UsunRamke(ctx context.Context, z shared.DesignFrameRemoveRequest) (shared.DesignFrameRemoveResponse, error)
	UlozAutomatycznie(ctx context.Context, z shared.DesignLayoutAutoRequest) (shared.DesignLayoutAutoResponse, error)
	UstawWiezy(ctx context.Context, z shared.DesignConstraintSetRequest) (shared.DesignConstraintSetResponse, error)
	ZmienRozmiarRamki(ctx context.Context, z shared.DesignFrameResizeApplyRequest) (shared.DesignFrameResizeApplyResponse, error)
	UstawSiatke(ctx context.Context, z shared.DesignGridSetRequest) (shared.DesignGridSetResponse, error)
	ZapiszKomponent(ctx context.Context, z shared.DesignComponentSaveRequest) (shared.DesignComponentSaveResponse, error)
	Komponenty(ctx context.Context, z shared.DesignComponentListRequest) (shared.DesignComponentListResponse, error)
	DolozInstancjeKomponentu(ctx context.Context, z shared.DesignComponentInstanceAddRequest) (shared.DesignComponentInstanceAddResponse, error)
	UstawPolaczeniePrototypu(ctx context.Context, z shared.DesignPrototypeLinkSetRequest) (shared.DesignPrototypeLinkSetResponse, error)
	Prototyp(ctx context.Context, z shared.DesignPrototypeGetRequest) (shared.DesignPrototypeGetResponse, error)
	UsunPolaczeniePrototypu(ctx context.Context, z shared.DesignPrototypeLinkRemoveRequest) (shared.DesignPrototypeLinkRemoveResponse, error)
	GenerujMakiete(ctx context.Context, z shared.DesignMockupGenerateRequest) (shared.DesignMockupGenerateResponse, error)
	WczytajMakiete(ctx context.Context, z shared.DesignMockupImportRequest) (shared.DesignMockupImportResponse, error)

	// --- barwa (adapter_modul_design_kolor.go, _kolor_obraz.go) ---
	GenerujPalete(ctx context.Context, z shared.DesignColorPaletteGenerateRequest) (shared.DesignColorPaletteGenerateResponse, error)
	WyciagnijPalete(ctx context.Context, z shared.DesignColorPaletteExtractRequest) (shared.DesignColorPaletteExtractResponse, error)
	SprawdzKontrast(ctx context.Context, z shared.DesignColorContrastCheckRequest) (shared.DesignColorContrastCheckResponse, error)
	SymulujWidzenie(ctx context.Context, z shared.DesignColorVisionSimulateRequest) (shared.DesignColorVisionSimulateResponse, error)
	UstawGradient(ctx context.Context, z shared.DesignColorGradientSetRequest) (shared.DesignColorGradientSetResponse, error)
	PrzeliczBarwe(ctx context.Context, z shared.DesignColorConvertRequest) (shared.DesignColorConvertResponse, error)
	ZbadajDostepnoscBarw(ctx context.Context, z shared.DesignColorAccessibilityAuditRequest) (shared.DesignColorAccessibilityAuditResponse, error)

	// --- ikony i kroje (adapter_modul_design_ikony.go) ---
	SzukajIkon(ctx context.Context, z shared.DesignIconLibrarySearchRequest) (shared.DesignIconLibrarySearchResponse, error)
	UstawIkone(ctx context.Context, z shared.DesignIconSetRequest) (shared.DesignIconSetResponse, error)
	GenerujIkony(ctx context.Context, z shared.DesignIconGenerateRequest) (shared.DesignIconGenerateResponse, error)
	ZbudujPakietIkon(ctx context.Context, z shared.DesignIconSpriteBuildRequest) (shared.DesignIconSpriteBuildResponse, error)
	ZbudujFavicone(ctx context.Context, z shared.DesignFaviconBuildRequest) (shared.DesignFaviconBuildResponse, error)
	ZaproponujZestawieniaKrojow(ctx context.Context, z shared.DesignFontPairSuggestRequest) (shared.DesignFontPairSuggestResponse, error)
	PodgladKroju(ctx context.Context, z shared.DesignFontPreviewRequest) (shared.DesignFontPreviewResponse, error)
	GlifyKroju(ctx context.Context, z shared.DesignFontGlyphsGetRequest) (shared.DesignFontGlyphsGetResponse, error)

	// --- szablony materiału (adapter_modul_design_szablony_materialu.go) ---
	ZapiszSzablonMaterialu(ctx context.Context, z shared.DesignTemplateSaveRequest) (shared.DesignTemplateSaveResponse, error)
	SzablonyMaterialu(ctx context.Context, z shared.DesignTemplateListRequest) (shared.DesignTemplateListResponse, error)
	ZastosujSzablonMaterialu(ctx context.Context, z shared.DesignTemplateApplyRequest) (shared.DesignTemplateApplyResponse, error)

	// --- marketing i bazy zdjęciowe (adapter_modul_design_marketing.go) ---
	ZbudujKompletKampanii(ctx context.Context, z shared.DesignCampaignSetBuildRequest) (shared.DesignCampaignSetBuildResponse, error)
	SzukajWBazachZdjeciowych(ctx context.Context, z shared.DesignStockSearchRequest) (shared.DesignStockSearchResponse, error)
	WciagnijZBazyZdjeciowej(ctx context.Context, z shared.DesignStockImportRequest) (shared.DesignStockImportResponse, error)
	WyrysujMakieteProduktowa(ctx context.Context, z shared.DesignProductMockupRenderRequest) (shared.DesignProductMockupRenderResponse, error)

	// --- druk, wykres i schemat (adapter_modul_design_druk.go, _wykresy.go) ---
	UstawProfilDruku(ctx context.Context, z shared.DesignPrintProfileSetRequest) (shared.DesignPrintProfileSetResponse, error)
	ProfileDruku(ctx context.Context, z shared.DesignPrintProfileListRequest) (shared.DesignPrintProfileListResponse, error)
	NosnikiDruku(ctx context.Context, z shared.DesignPrintPaperListRequest) (shared.DesignPrintPaperListResponse, error)
	KontrolaPrzeddrukowa(ctx context.Context, z shared.DesignPrintPreflightRequest) (shared.DesignPrintPreflightResponse, error)
	WydajDoDruku(ctx context.Context, z shared.DesignPrintExportRequest) (shared.DesignPrintExportResponse, error)
	PodzielMaterialWielkoformatowy(ctx context.Context, z shared.DesignLargeformatTileRequest) (shared.DesignLargeformatTileResponse, error)
	WyrysujWykres(ctx context.Context, z shared.DesignChartRenderRequest) (shared.DesignChartRenderResponse, error)
	WyrysujSchemat(ctx context.Context, z shared.DesignDiagramRenderRequest) (shared.DesignDiagramRenderResponse, error)

	// --- warsztat fotografii (adapter_modul_design_fotografia.go, _fotografia_wsad.go) ---
	Kadruj(ctx context.Context, z shared.DesignPhotoCropRequest) (shared.DesignPhotoCropResponse, error)
	Przeksztalc(ctx context.Context, z shared.DesignPhotoTransformRequest) (shared.DesignPhotoTransformResponse, error)
	PrzeliczRozdzielczosc(ctx context.Context, z shared.DesignPhotoResampleRequest) (shared.DesignPhotoResampleResponse, error)
	Powieksz(ctx context.Context, z shared.DesignPhotoUpscaleRequest) (shared.DesignPhotoUpscaleResponse, error)
	PopraweJakosc(ctx context.Context, z shared.DesignPhotoEnhanceRequest) (shared.DesignPhotoEnhanceResponse, error)
	SkorygujBarwe(ctx context.Context, z shared.DesignPhotoColorCorrectRequest) (shared.DesignPhotoColorCorrectResponse, error)
	NalozFiltr(ctx context.Context, z shared.DesignPhotoFilterApplyRequest) (shared.DesignPhotoFilterApplyResponse, error)
	Wyretuszuj(ctx context.Context, z shared.DesignPhotoRetouchRequest) (shared.DesignPhotoRetouchResponse, error)
	Domaluj(ctx context.Context, z shared.DesignPhotoInpaintRequest) (shared.DesignPhotoInpaintResponse, error)
	RozszerzKadr(ctx context.Context, z shared.DesignPhotoExpandRequest) (shared.DesignPhotoExpandResponse, error)
	OdetnijTlo(ctx context.Context, z shared.DesignPhotoBackgroundRemoveRequest) (shared.DesignPhotoBackgroundRemoveResponse, error)
	ZaznaczObiekt(ctx context.Context, z shared.DesignPhotoSelectObjectRequest) (shared.DesignPhotoSelectObjectResponse, error)
	UstawMaske(ctx context.Context, z shared.DesignPhotoMaskSetRequest) (shared.DesignPhotoMaskSetResponse, error)
	ZlozWarstwyFotografii(ctx context.Context, z shared.DesignPhotoLayerCompositeRequest) (shared.DesignPhotoLayerCompositeResponse, error)
	PuscWsadFotografii(ctx context.Context, z shared.DesignPhotoBatchApplyRequest) (shared.DesignPhotoBatchApplyResponse, error)
	ZapiszNastaweFotografii(ctx context.Context, z shared.DesignPhotoPresetSaveRequest) (shared.DesignPhotoPresetSaveResponse, error)
	NastawyFotografii(ctx context.Context, z shared.DesignPhotoPresetListRequest) (shared.DesignPhotoPresetListResponse, error)
	ObrysujKontury(ctx context.Context, z shared.DesignPhotoVectorizeRequest) (shared.DesignPhotoVectorizeResponse, error)
	MetadaneZasobu(ctx context.Context, z shared.DesignPhotoMetadataGetRequest) (shared.DesignPhotoMetadataGetResponse, error)
	LancuchEdycji(ctx context.Context, z shared.DesignPhotoHistoryGetRequest) (shared.DesignPhotoHistoryGetResponse, error)

	// --- odczyty dla szyny zdarzeń — te dwie metody nie obsługują komendy; czytają stan po zapisie ---
	KompozycjaZdarzenia(ctx context.Context, kod string) (shared.DesignBoard, bool)
	KompozycjaRamkiZdarzenia(ctx context.Context, ramka string) (shared.DesignBoard, bool)
}

// zarejestrujDesign wpina osiem komend modułu Design bezpośrednio oraz woła
// rejestrację pozostałych obszarów.
func zarejestrujDesign(r *Rejestr, m Design, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandDesignAssetList, obsluz(m.Zasoby))
	// Zapis układu kompozycji rozgłasza zmianę tak samo jak zastosowanie
	// szablonu: bez tego okno, które kompozycji nie zapisywało, nie wie o niej.
	r.Zarejestruj(shared.CommandDesignBoardUpdate,
		obsluz(func(ctx context.Context, z shared.DesignBoardUpdateRequest) (shared.DesignBoardUpdateResponse, error) {
			odpowiedz, err := m.ZapiszKompozycje(ctx, z)
			if err == nil {
				e.kompozycjaDesignu(ctx, shared.ChangeKindUpdated, odpowiedz.Board)
			}
			return odpowiedz, err
		}))
	// Odczyt kompozycji okna, bez opakowania rozgłaszającego — niczego nie zmienia, jak wykaz zasobów.
	r.Zarejestruj(shared.CommandDesignBoardList, obsluz(m.Kompozycje))

	// Zdarzenie idzie po zasobie, nie po komendzie: każdy wariant ma bajty, wiersz i zdarzenie własne.
	r.Zarejestruj(shared.CommandDesignAssetGenerate,
		obsluz(func(ctx context.Context, z shared.DesignAssetGenerateRequest) (shared.DesignAssetGenerateResponse, error) {
			odpowiedz, err := m.GenerujZasob(ctx, z)
			if err == nil {
				for _, zasob := range odpowiedz.Assets {
					e.zasobDesignu(ctx, shared.ChangeKindCreated, zasob)
				}
			}
			return odpowiedz, err
		}))

	// `design.asset.changed` jest jedynym zdarzeniem obszaru; nadanie etykiet daje rodzaj `updated`.
	r.Zarejestruj(shared.CommandDesignAssetTagSet,
		obsluz(func(ctx context.Context, z shared.DesignAssetTagSetRequest) (shared.DesignAssetTagSetResponse, error) {
			odpowiedz, err := m.UstawEtykietyZasobu(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindUpdated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))

	// Wniesienie jest drugim źródłem `created` — obie drogi kończą się bajtami, wierszem i zdarzeniem.
	r.Zarejestruj(shared.CommandDesignAssetUpload,
		obsluz(func(ctx context.Context, z shared.DesignAssetUploadRequest) (shared.DesignAssetUploadResponse, error) {
			odpowiedz, err := m.WniesZasob(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindCreated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))

	// Oznaczenie ulubionego zmienia zasób zastany — `updated`, tak jak etykiety.
	r.Zarejestruj(shared.CommandDesignAssetFavoriteSet,
		obsluz(func(ctx context.Context, z shared.DesignAssetFavoriteSetRequest) (shared.DesignAssetFavoriteSetResponse, error) {
			odpowiedz, err := m.UstawUlubionyZasob(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindUpdated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))

	// Usunięcie jest jedynym nadawcą `deleted`, tylko gdy wiersz naprawdę zniknął — inaczej szyna milczy.
	r.Zarejestruj(shared.CommandDesignAssetRemove,
		obsluz(func(ctx context.Context, z shared.DesignAssetRemoveRequest) (shared.DesignAssetRemoveResponse, error) {
			odpowiedz, usuniety, err := m.UsunZasob(ctx, z)
			if err == nil && odpowiedz.Removed {
				e.zasobDesignu(ctx, shared.ChangeKindDeleted, usuniety)
			}
			return odpowiedz, err
		}))

	zarejestrujDesignWydania(r, m)
	zarejestrujDesignKolekcjeIPrompty(r, m)
	zarejestrujDesignPlansze(r, m, e)
	zarejestrujDesignZetony(r, m)
	zarejestrujDesignWektor(r, m)
	zarejestrujDesignMakiety(r, m, e)
	zarejestrujDesignBarwy(r, m, e)
	zarejestrujDesignIkony(r, m, e)
	zarejestrujDesignSzablonyMaterialu(r, m, e)
	zarejestrujDesignMarketing(r, m, e)
	zarejestrujDesignDruk(r, m, e)
	zarejestrujDesignFotografia(r, m, e)
}

// zarejestrujDesignFotografia wpina dwadzieścia czynności warsztatu fotografii.
// Każda czynność zakładająca wariant zasobu rozgłasza `created`; trzy czynności
// milczą, bo są odczytem albo oddają wykaz bez pojedynczego zasobu.
func zarejestrujDesignFotografia(r *Rejestr, m Design, e *emiter) {
	r.Zarejestruj(shared.CommandDesignPhotoMetadataGet, obsluz(m.MetadaneZasobu))
	r.Zarejestruj(shared.CommandDesignPhotoHistoryGet, obsluz(m.LancuchEdycji))
	r.Zarejestruj(shared.CommandDesignPhotoBatchApply, obsluz(m.PuscWsadFotografii))
	r.Zarejestruj(shared.CommandDesignPhotoPresetSave, obsluz(m.ZapiszNastaweFotografii))
	r.Zarejestruj(shared.CommandDesignPhotoPresetList, obsluz(m.NastawyFotografii))

	r.Zarejestruj(shared.CommandDesignPhotoCrop, zZasobemFotografiiDesignu(e, m.Kadruj,
		func(o shared.DesignPhotoCropResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoTransform, zZasobemFotografiiDesignu(e, m.Przeksztalc,
		func(o shared.DesignPhotoTransformResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoResample, zZasobemFotografiiDesignu(e, m.PrzeliczRozdzielczosc,
		func(o shared.DesignPhotoResampleResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoUpscale, zZasobemFotografiiDesignu(e, m.Powieksz,
		func(o shared.DesignPhotoUpscaleResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoEnhance, zZasobemFotografiiDesignu(e, m.PopraweJakosc,
		func(o shared.DesignPhotoEnhanceResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoColorCorrect, zZasobemFotografiiDesignu(e, m.SkorygujBarwe,
		func(o shared.DesignPhotoColorCorrectResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoFilterApply, zZasobemFotografiiDesignu(e, m.NalozFiltr,
		func(o shared.DesignPhotoFilterApplyResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoRetouch, zZasobemFotografiiDesignu(e, m.Wyretuszuj,
		func(o shared.DesignPhotoRetouchResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoInpaint, zZasobemFotografiiDesignu(e, m.Domaluj,
		func(o shared.DesignPhotoInpaintResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoExpand, zZasobemFotografiiDesignu(e, m.RozszerzKadr,
		func(o shared.DesignPhotoExpandResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoBackgroundRemove, zZasobemFotografiiDesignu(e, m.OdetnijTlo,
		func(o shared.DesignPhotoBackgroundRemoveResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoSelectObject, zZasobemFotografiiDesignu(e, m.ZaznaczObiekt,
		func(o shared.DesignPhotoSelectObjectResponse) shared.DesignAsset { return o.Mask }))
	r.Zarejestruj(shared.CommandDesignPhotoMaskSet, zZasobemFotografiiDesignu(e, m.UstawMaske,
		func(o shared.DesignPhotoMaskSetResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoLayerComposite, zZasobemFotografiiDesignu(e, m.ZlozWarstwyFotografii,
		func(o shared.DesignPhotoLayerCompositeResponse) shared.DesignAsset { return o.Asset }))
	r.Zarejestruj(shared.CommandDesignPhotoVectorize, zZasobemFotografiiDesignu(e, m.ObrysujKontury,
		func(o shared.DesignPhotoVectorizeResponse) shared.DesignAsset { return o.Asset }))
}

// zZasobemFotografiiDesignu składa obsługę jednej czynności warsztatu
// fotografii wraz z rozgłoszeniem zasobu, który z niej powstał — jedna funkcja
// zamiast piętnastu identycznych opakowań.
func zZasobemFotografiiDesignu[Z any, W any](e *emiter,
	czynnosc func(context.Context, Z) (W, error), zasob func(W) shared.DesignAsset) Obsluga {

	return obsluz(func(ctx context.Context, z Z) (W, error) {
		odpowiedz, err := czynnosc(ctx, z)
		if err == nil {
			e.zasobDesignu(ctx, shared.ChangeKindCreated, zasob(odpowiedz))
		}
		return odpowiedz, err
	})
}

// zarejestrujDesignWektor wpina dziesięć czynności warsztatu wektorowego.
// Żadna nie rozgłasza `design.asset.changed`: ścieżka i symbol są bytami
// kompozycji, nie zasobami Assets Panelu.
func zarejestrujDesignWektor(r *Rejestr, m Design) {
	r.Zarejestruj(shared.CommandDesignVectorPathSet, obsluz(m.UstawSciezke))
	r.Zarejestruj(shared.CommandDesignVectorPathList, obsluz(m.SciezkiWektorowe))
	r.Zarejestruj(shared.CommandDesignVectorPathRemove, obsluz(m.UsunSciezke))
	r.Zarejestruj(shared.CommandDesignVectorShapeAdd, obsluz(m.DolozKsztalt))
	r.Zarejestruj(shared.CommandDesignVectorBoolean, obsluz(m.ZlozSciezkiLogicznie))
	r.Zarejestruj(shared.CommandDesignVectorTextPath, obsluz(m.TekstNaSciezce))
	r.Zarejestruj(shared.CommandDesignVectorOptimize, obsluz(m.OczyscSciezki))
	r.Zarejestruj(shared.CommandDesignVectorSymbolSet, obsluz(m.UstawSymbol))
	r.Zarejestruj(shared.CommandDesignVectorSymbolList, obsluz(m.Symbole))
	r.Zarejestruj(shared.CommandDesignVectorExport, obsluz(m.WydajWektor))
}

// zarejestrujDesignMakiety wpina piętnaście czynności warsztatu makiety.
// Rozgłoszenie `design.board.changed` idzie po tych, które zmieniają układ
// warstw kompozycji; ramka, więz, siatka, komponent i przejście prototypu nie
// rozgłaszają go.
func zarejestrujDesignMakiety(r *Rejestr, m Design, e *emiter) {
	r.Zarejestruj(shared.CommandDesignFrameSet, obsluz(m.UstawRamke))
	r.Zarejestruj(shared.CommandDesignFrameList, obsluz(m.Ramki))
	r.Zarejestruj(shared.CommandDesignFrameRemove, obsluz(m.UsunRamke))
	r.Zarejestruj(shared.CommandDesignConstraintSet, obsluz(m.UstawWiezy))
	r.Zarejestruj(shared.CommandDesignGridSet, obsluz(m.UstawSiatke))
	r.Zarejestruj(shared.CommandDesignComponentSave, obsluz(m.ZapiszKomponent))
	r.Zarejestruj(shared.CommandDesignComponentList, obsluz(m.Komponenty))
	r.Zarejestruj(shared.CommandDesignPrototypeLinkSet, obsluz(m.UstawPolaczeniePrototypu))
	r.Zarejestruj(shared.CommandDesignPrototypeGet, obsluz(m.Prototyp))
	r.Zarejestruj(shared.CommandDesignPrototypeLinkRemove, obsluz(m.UsunPolaczeniePrototypu))

	r.Zarejestruj(shared.CommandDesignLayoutAuto,
		obsluz(func(ctx context.Context, z shared.DesignLayoutAutoRequest) (shared.DesignLayoutAutoResponse, error) {
			odpowiedz, err := m.UlozAutomatycznie(ctx, z)
			if err == nil {
				e.kompozycjaRamkiDesignu(ctx, m, z.FrameId)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandDesignFrameResizeApply,
		obsluz(func(ctx context.Context, z shared.DesignFrameResizeApplyRequest) (shared.DesignFrameResizeApplyResponse, error) {
			odpowiedz, err := m.ZmienRozmiarRamki(ctx, z)
			if err == nil {
				e.kompozycjaPoKodzieDesignu(ctx, m, odpowiedz.Frame.BoardId)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandDesignComponentInstanceAdd,
		obsluz(func(ctx context.Context, z shared.DesignComponentInstanceAddRequest) (shared.DesignComponentInstanceAddResponse, error) {
			odpowiedz, err := m.DolozInstancjeKomponentu(ctx, z)
			if err == nil {
				e.kompozycjaPoKodzieDesignu(ctx, m, z.BoardId)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandDesignMockupGenerate,
		obsluz(func(ctx context.Context, z shared.DesignMockupGenerateRequest) (shared.DesignMockupGenerateResponse, error) {
			odpowiedz, err := m.GenerujMakiete(ctx, z)
			if err == nil {
				e.kompozycjaPoKodzieDesignu(ctx, m, z.BoardId)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandDesignMockupImport,
		obsluz(func(ctx context.Context, z shared.DesignMockupImportRequest) (shared.DesignMockupImportResponse, error) {
			odpowiedz, err := m.WczytajMakiete(ctx, z)
			if err == nil {
				e.kompozycjaPoKodzieDesignu(ctx, m, z.BoardId)
			}
			return odpowiedz, err
		}))
}

// zarejestrujDesignBarwy wpina siedem czynności barwy. Zdarzenie idzie po
// jednej: symulacja wady widzenia zakłada zasób w magazynie i rozgłasza
// `created`; pozostałe są rachunkiem.
func zarejestrujDesignBarwy(r *Rejestr, m Design, e *emiter) {
	r.Zarejestruj(shared.CommandDesignColorPaletteGenerate, obsluz(m.GenerujPalete))
	r.Zarejestruj(shared.CommandDesignColorPaletteExtract, obsluz(m.WyciagnijPalete))
	r.Zarejestruj(shared.CommandDesignColorContrastCheck, obsluz(m.SprawdzKontrast))
	r.Zarejestruj(shared.CommandDesignColorConvert, obsluz(m.PrzeliczBarwe))
	r.Zarejestruj(shared.CommandDesignColorAccessibilityAudit, obsluz(m.ZbadajDostepnoscBarw))
	r.Zarejestruj(shared.CommandDesignColorGradientSet, obsluz(m.UstawGradient))

	r.Zarejestruj(shared.CommandDesignColorVisionSimulate,
		obsluz(func(ctx context.Context, z shared.DesignColorVisionSimulateRequest) (shared.DesignColorVisionSimulateResponse, error) {
			odpowiedz, err := m.SymulujWidzenie(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindCreated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))
}

// zarejestrujDesignIkony wpina osiem czynności ikon i krojów. Pakiet ikon
// zakłada zasób magazynu i rozgłasza `created`; ikony własne i favicony milczą,
// bo nie są zasobem Assets Panelu albo nie niosą go w odpowiedzi.
func zarejestrujDesignIkony(r *Rejestr, m Design, e *emiter) {
	r.Zarejestruj(shared.CommandDesignIconLibrarySearch, obsluz(m.SzukajIkon))
	r.Zarejestruj(shared.CommandDesignIconSet, obsluz(m.UstawIkone))
	r.Zarejestruj(shared.CommandDesignIconGenerate, obsluz(m.GenerujIkony))
	r.Zarejestruj(shared.CommandDesignFaviconBuild, obsluz(m.ZbudujFavicone))
	r.Zarejestruj(shared.CommandDesignFontPairSuggest, obsluz(m.ZaproponujZestawieniaKrojow))
	r.Zarejestruj(shared.CommandDesignFontPreview, obsluz(m.PodgladKroju))
	r.Zarejestruj(shared.CommandDesignFontGlyphsGet, obsluz(m.GlifyKroju))

	r.Zarejestruj(shared.CommandDesignIconSpriteBuild,
		obsluz(func(ctx context.Context, z shared.DesignIconSpriteBuildRequest) (shared.DesignIconSpriteBuildResponse, error) {
			odpowiedz, err := m.ZbudujPakietIkon(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindCreated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))
}

// zarejestrujDesignSzablonyMaterialu wpina trzy czynności szablonów materiału.
//
// Zastosowanie szablonu ZAKŁADA ALBO PRZESTAWIA kompozycję, więc rozgłasza
// `design.board.changed`. Zapis i wykaz szablonów milczą — szablon nie jest
// kompozycją ani zasobem.
func zarejestrujDesignSzablonyMaterialu(r *Rejestr, m Design, e *emiter) {
	r.Zarejestruj(shared.CommandDesignTemplateSave, obsluz(m.ZapiszSzablonMaterialu))
	r.Zarejestruj(shared.CommandDesignTemplateList, obsluz(m.SzablonyMaterialu))

	r.Zarejestruj(shared.CommandDesignTemplateApply,
		obsluz(func(ctx context.Context, z shared.DesignTemplateApplyRequest) (shared.DesignTemplateApplyResponse, error) {
			odpowiedz, err := m.ZastosujSzablonMaterialu(ctx, z)
			if err == nil {
				e.kompozycjaDesignu(ctx, shared.ChangeKindUpdated, odpowiedz.Board)
			}
			return odpowiedz, err
		}))
}

// zarejestrujDesignMarketing wpina cztery czynności marketingowe. Wciągnięcie
// z bazy zdjęciowej i makieta produktowa zakładają zasób i rozgłaszają
// `created`; komplet kampanii i wyszukanie milczą.
func zarejestrujDesignMarketing(r *Rejestr, m Design, e *emiter) {
	r.Zarejestruj(shared.CommandDesignCampaignSetBuild, obsluz(m.ZbudujKompletKampanii))
	r.Zarejestruj(shared.CommandDesignStockSearch, obsluz(m.SzukajWBazachZdjeciowych))

	r.Zarejestruj(shared.CommandDesignStockImport,
		obsluz(func(ctx context.Context, z shared.DesignStockImportRequest) (shared.DesignStockImportResponse, error) {
			odpowiedz, err := m.WciagnijZBazyZdjeciowej(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindCreated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandDesignProductMockupRender,
		obsluz(func(ctx context.Context, z shared.DesignProductMockupRenderRequest) (shared.DesignProductMockupRenderResponse, error) {
			odpowiedz, err := m.WyrysujMakieteProduktowa(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindCreated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))
}

// zarejestrujDesignDruk wpina osiem czynności druku, wykresu i schematu.
// Wydanie, wykres i schemat zakładają zasób i rozgłaszają `created`; podział
// na kafle, kontrola i wykazy milczą.
func zarejestrujDesignDruk(r *Rejestr, m Design, e *emiter) {
	r.Zarejestruj(shared.CommandDesignPrintProfileSet, obsluz(m.UstawProfilDruku))
	r.Zarejestruj(shared.CommandDesignPrintProfileList, obsluz(m.ProfileDruku))
	r.Zarejestruj(shared.CommandDesignPrintPaperList, obsluz(m.NosnikiDruku))
	r.Zarejestruj(shared.CommandDesignPrintPreflight, obsluz(m.KontrolaPrzeddrukowa))
	r.Zarejestruj(shared.CommandDesignLargeformatTile, obsluz(m.PodzielMaterialWielkoformatowy))

	r.Zarejestruj(shared.CommandDesignPrintExport,
		obsluz(func(ctx context.Context, z shared.DesignPrintExportRequest) (shared.DesignPrintExportResponse, error) {
			odpowiedz, err := m.WydajDoDruku(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindCreated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandDesignChartRender,
		obsluz(func(ctx context.Context, z shared.DesignChartRenderRequest) (shared.DesignChartRenderResponse, error) {
			odpowiedz, err := m.WyrysujWykres(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindCreated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandDesignDiagramRender,
		obsluz(func(ctx context.Context, z shared.DesignDiagramRenderRequest) (shared.DesignDiagramRenderResponse, error) {
			odpowiedz, err := m.WyrysujSchemat(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindCreated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))
}

// kompozycjaPoKodzieDesignu rozgłasza `design.board.changed` dla kompozycji
// wskazanej kodem, odczytując jej stan po zmianie — zdarzenie ma nieść skutek,
// nie zamiar żądania.
func (e *emiter) kompozycjaPoKodzieDesignu(ctx context.Context, m Design, kod string) {
	if e == nil || strings.TrimSpace(kod) == "" {
		return
	}
	kompozycja, jest := m.KompozycjaZdarzenia(ctx, strings.TrimSpace(kod))
	if !jest {
		return
	}
	e.kompozycjaDesignu(ctx, shared.ChangeKindUpdated, kompozycja)
}

// kompozycjaRamkiDesignu rozgłasza `design.board.changed` dla kompozycji, do
// której należy wskazana ramka. Układ automatyczny wskazuje RAMKĘ, a zdarzenie
// dotyczy kompozycji — przekład idzie przez port, bo tylko adapter widzi bazę.
func (e *emiter) kompozycjaRamkiDesignu(ctx context.Context, m Design, ramka string) {
	if e == nil || strings.TrimSpace(ramka) == "" {
		return
	}
	kompozycja, jest := m.KompozycjaRamkiZdarzenia(ctx, strings.TrimSpace(ramka))
	if !jest {
		return
	}
	e.kompozycjaDesignu(ctx, shared.ChangeKindUpdated, kompozycja)
}

// zarejestrujDesignWydania wpina trzy komendy oddające treść poza rdzeń.
// Żadna z nich niczego nie zmienia, więc żadna nie ma opakowania
// rozgłaszającego — tak samo jak `design.asset.list`.
func zarejestrujDesignWydania(r *Rejestr, m Design) {
	r.Zarejestruj(shared.CommandDesignAssetContentGet, obsluz(m.TrescZasobu))
	r.Zarejestruj(shared.CommandDesignAssetExport, obsluz(m.WydajZasob))
	r.Zarejestruj(shared.CommandDesignAssetExportBatch, obsluz(m.WydajZasobyPartia))
}

// zarejestrujDesignKolekcjeIPrompty wpina kolekcje zasobów oraz szablony
// i historię promptów. Kolekcja i szablon nie rozgłaszają zdarzenia: żaden nie
// jest zasobem, kompozycją ani obecnością.
func zarejestrujDesignKolekcjeIPrompty(r *Rejestr, m Design) {
	r.Zarejestruj(shared.CommandDesignCollectionCreate, obsluz(m.ZalozKolekcje))
	r.Zarejestruj(shared.CommandDesignCollectionAssign, obsluz(m.PrzypiszDoKolekcji))
	r.Zarejestruj(shared.CommandDesignCollectionList, obsluz(m.Kolekcje))

	r.Zarejestruj(shared.CommandDesignPromptTemplateSave, obsluz(m.ZapiszSzablonPromptu))
	r.Zarejestruj(shared.CommandDesignPromptTemplateList, obsluz(m.SzablonyPromptu))
	r.Zarejestruj(shared.CommandDesignPromptHistoryList, obsluz(m.HistoriaPromptow))
}

// zarejestrujDesignPlansze wpina wersjonowanie kompozycji, jej wyrys,
// adnotacje i obecność. Przywrócenie wersji rozgłasza `design.board.changed`;
// zgłoszenie obecności rozgłasza `design.board.presence` z kompletem obecnych.
func zarejestrujDesignPlansze(r *Rejestr, m Design, e *emiter) {
	r.Zarejestruj(shared.CommandDesignBoardVersionSave, obsluz(m.ZapiszWersjeKompozycji))
	r.Zarejestruj(shared.CommandDesignBoardVersionList, obsluz(m.WersjeKompozycji))
	r.Zarejestruj(shared.CommandDesignBoardExport, obsluz(m.WyrysujKompozycje))
	r.Zarejestruj(shared.CommandDesignAnnotationList, obsluz(m.Adnotacje))
	r.Zarejestruj(shared.CommandDesignAnnotationSet, obsluz(m.UstawAdnotacje))

	r.Zarejestruj(shared.CommandDesignBoardVersionRestore,
		obsluz(func(ctx context.Context, z shared.DesignBoardVersionRestoreRequest) (shared.DesignBoardVersionRestoreResponse, error) {
			odpowiedz, err := m.PrzywrocWersjeKompozycji(ctx, z)
			if err == nil {
				e.kompozycjaDesignu(ctx, shared.ChangeKindUpdated, odpowiedz.Board)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandDesignPresenceReport,
		obsluz(func(ctx context.Context, z shared.DesignPresenceReportRequest) (shared.DesignPresenceReportResponse, error) {
			odpowiedz, err := m.ZglosObecnosc(ctx, z)
			if err == nil {
				e.obecnoscDesignu(z.BoardId, odpowiedz.Participants)
			}
			return odpowiedz, err
		}))
}

// zarejestrujDesignZetony wpina zestawy żetonów i przewodnik stylu; żadna
// z tych pięciu czynności nie rozgłasza zdarzenia.
func zarejestrujDesignZetony(r *Rejestr, m Design) {
	r.Zarejestruj(shared.CommandDesignTokensetSave, obsluz(m.ZapiszZestawZetonow))
	r.Zarejestruj(shared.CommandDesignTokensetList, obsluz(m.ZestawyZetonow))
	r.Zarejestruj(shared.CommandDesignTokensetExport, obsluz(m.WydajZestawZetonow))
	r.Zarejestruj(shared.CommandDesignTokensetImport, obsluz(m.WczytajZestawZetonow))
	r.Zarejestruj(shared.CommandDesignStyleguidePublish, obsluz(m.WydajPrzewodnikStylu))
}

// kompozycjaDesignu rozgłasza `design.board.changed`. Kompozycja jest bytem
// okna modułu, nie karty sesji, więc zdarzenie idzie bez wskazania sesji —
// wzorem `zasobDesignu`.
func (e *emiter) kompozycjaDesignu(ctx context.Context, zmiana shared.ChangeKind,
	kompozycja shared.DesignBoard) {

	zdarzenie := shared.DesignBoardChangedEvent{Change: zmiana, Board: kompozycja}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventDesignBoardChanged, "", zdarzenie)
}

// obecnoscDesignu rozgłasza `design.board.presence`. Zdarzenie nie niesie
// sprawcy: sprawcą jest jeden z obecnych i już jest w wykazie pod swoim
// identyfikatorem klienta, a drugie wskazanie tego samego nie dokłada wiedzy.
func (e *emiter) obecnoscDesignu(kompozycja string, obecni []shared.DesignPresence) {
	e.wyslij(shared.EventDesignBoardPresence, "", shared.DesignBoardPresenceEvent{
		BoardId: kompozycja, Participants: obecni,
	})
}

// zasobDesignu rozgłasza `design.asset.changed`. Zasób jest bytem okna modułu
// Design (niesie własne `windowId`), nie karty sesji, więc zdarzenie idzie bez
// wskazania sesji — Assets Panel odbiera je stroną własną, wzorem
// `raportBadania`.
func (e *emiter) zasobDesignu(ctx context.Context, zmiana shared.ChangeKind, zasob shared.DesignAsset) {
	zdarzenie := shared.DesignAssetChangedEvent{Change: zmiana, Asset: zasob}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventDesignAssetChanged, "", zdarzenie)
}
