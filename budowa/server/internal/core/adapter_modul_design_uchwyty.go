// Odpowiedzialność pliku: wpięcie ośmiu komend obszaru `design.*` — modułu
// Design (Prompt Builder, Assets Panel, Design Board) — i rozgłoszenie
// `design.asset.changed` po tych z nich, które zasób zmieniają.
//
// `design.asset.tag.set` domyka lukę odczytu: `design.asset.list` zawęża wykaz
// polem `tags`, a bez tej komendy nie byłoby czym etykiet nadać; uchwyt leży
// w `adapter_modul_design_etykiety.go`.
//
// Zdarzenie `created` ma dwóch nadawców: generowanie i wniesienie. Obie drogi
// odkładają bajty w magazynie rdzenia pod sumą kontrolną, zanim powstanie choć
// jeden wiersz (`adapter_modul_design_generowanie.go`), a odmowa nie rozgłasza
// niczego — opakowanie milczy przy błędzie, więc droga bez bajtów jest też
// drogą bez zdarzenia.
//
// Generowanie rozgłasza tyle zdarzeń, ile założyło zasobów. Wariantów bywa
// kilka, a `design.asset.changed` niesie jeden zasób — jedno zdarzenie na cały
// zbiór opisywałoby powstanie jednego z nich i przemilczało resztę.
package core

import (
	"context"
	"strings"

	"danacoconsole/shared"
)

// Design jest portem modułu Design.
type Design interface {
	GenerujZasob(ctx context.Context, z shared.DesignAssetGenerateRequest) (shared.DesignAssetGenerateResponse, error)
	Zasoby(ctx context.Context, z shared.DesignAssetListRequest) (shared.DesignAssetListResponse, error)
	ZapiszKompozycje(ctx context.Context, z shared.DesignBoardUpdateRequest) (shared.DesignBoardUpdateResponse, error)
	UstawEtykietyZasobu(ctx context.Context, z shared.DesignAssetTagSetRequest) (shared.DesignAssetTagSetResponse, error)
	WniesZasob(ctx context.Context, z shared.DesignAssetUploadRequest) (shared.DesignAssetUploadResponse, error)
	UstawUlubionyZasob(ctx context.Context, z shared.DesignAssetFavoriteSetRequest) (shared.DesignAssetFavoriteSetResponse, error)
	// UsunZasob oddaje trzy wartości, jako jedyna w tym porcie. Odpowiedź
	// kontraktu niesie samo `removed`, a zdarzenie `deleted` musi nieść cały
	// usunięty zasób — po usunięciu wiersza nie ma go już skąd odczytać, więc
	// adapter podaje go obok odpowiedzi. Rozgłaszanie z wnętrza adaptera byłoby
	// drugą drogą do szyny zdarzeń (emiter należy do rejestru, nie do modułu).
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

	// --- odczyty dla szyny zdarzeń ---
	// Te dwie metody nie obsługują żadnej komendy. Istnieją, bo zdarzenie
	// `design.board.changed` niesie CAŁĄ kompozycję, a komendy, które ją zmieniają
	// (układ automatyczny, przeliczenie więzi, instancja komponentu, obie drogi
	// makiety), oddają w odpowiedzi wyłącznie swój wycinek. Rejestr musi więc
	// odczytać stan po zapisie, a bazę widzi wyłącznie adapter.
	//
	// Brak wiersza oddają jako `false`, nie jako błąd: zmiana się już udała, więc
	// niepowodzenie odczytu ma zamknąć usta szynie zdarzeń, a nie unieważnić
	// komendę Operatora.
	KompozycjaZdarzenia(ctx context.Context, kod string) (shared.DesignBoard, bool)
	KompozycjaRamkiZdarzenia(ctx context.Context, ramka string) (shared.DesignBoard, bool)
}

// zarejestrujDesign wpina osiem komend modułu Design.
func zarejestrujDesign(r *Rejestr, m Design, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandDesignAssetList, obsluz(m.Zasoby))
	r.Zarejestruj(shared.CommandDesignBoardUpdate, obsluz(m.ZapiszKompozycje))
	// Odczyt kompozycji okna. Bez opakowania rozgłaszającego — niczego nie
	// zmienia, tak samo jak `design.asset.list`.
	r.Zarejestruj(shared.CommandDesignBoardList, obsluz(m.Kompozycje))

	// Zdarzenie idzie po zasobie, nie po komendzie: każdy wariant dostaje własne
	// bajty w magazynie i własny wiersz, więc każdy ma własne zdarzenie
	// (nagłówek pliku).
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

	// `design.asset.changed` jest jedynym zdarzeniem obszaru; nadanie etykiet
	// zmienia zasób zastany, więc `updated`. `design.asset.list`
	// i `design.board.list` niczego nie zmieniają, a `design.board.update`
	// zmienia kompozycję, dla której kontrakt osobnego zdarzenia nie ma —
	// kompozycji więc nie rozgłaszamy zdarzeniem zasobu, bo odbiorca dostałby
	// zmianę bytu, który się nie zmienił.
	r.Zarejestruj(shared.CommandDesignAssetTagSet,
		obsluz(func(ctx context.Context, z shared.DesignAssetTagSetRequest) (shared.DesignAssetTagSetResponse, error) {
			odpowiedz, err := m.UstawEtykietyZasobu(ctx, z)
			if err == nil {
				e.zasobDesignu(ctx, shared.ChangeKindUpdated, odpowiedz.Asset)
			}
			return odpowiedz, err
		}))

	// Wniesienie jest drugim źródłem `created` — obie drogi zasobu do modułu
	// (wniesiona przez Operatora i wygenerowana kanałem) kończą się tak samo:
	// bajtami w magazynie pod sumą kontrolną, wierszem i zdarzeniem.
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

	// Usunięcie jest jedynym nadawcą `deleted`. Zdarzenie idzie wyłącznie wtedy,
	// gdy wiersz naprawdę zniknął (`Removed`): usunięcie zasobu, którego nie
	// było, kończy się odpowiedzią `removed: false` i milczeniem szyny —
	// rozgłoszenie donosiłoby panelowi o zniknięciu czegoś, czego nie miał.
	// Zasób w zdarzeniu jest tym odczytanym przed usunięciem, bo po nim nie ma
	// już czego czytać.
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
//
// Każda czynność zakładająca WARIANT zasobu rozgłasza `design.asset.changed`
// z rodzajem `created`: Assets Panel ma pokazać nowy wariant natychmiast, bo
// warsztat fotografii jest pracą ciągłą i Operator sięga po poprzedni wynik
// w następnym kroku.
//
// Trzy czynności milczą i każda z innego powodu. `design.photo.metadata.get`
// i `design.photo.history.get` niczego nie zmieniają — są odczytem.
// `design.photo.batch.apply` zakłada po jednym wariancie na zasób i oddaje wykaz;
// jedno zdarzenie opisywałoby jeden z nich i przemilczało resztę, a rozgłoszenie
// dwustu zdarzeń z jednej komendy zalałoby szynę. Panel odświeża się tam
// odpowiedzią komendy, która niesie komplet zasobów.
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
// fotografii wraz z rozgłoszeniem zasobu, który z niej powstał.
//
// Jedna funkcja na piętnaście czynności, a nie piętnaście opakowań: wszystkie
// robią DOKŁADNIE to samo — wołają uchwyt i rozgłaszają jeden zasób
// z odpowiedzi. Piętnaście kopii tego samego opakowania byłoby piętnastoma
// miejscami, w których da się zapomnieć o zdarzeniu.
//
// Wpięcie zostaje przy wołającym, z nazwą komendy wpisaną WPROST: wykaz komend
// obsługiwanych przez rdzeń czyta się z treści plików rdzenia (`shared.Command…`
// obok `Zarejestruj`), a nazwa komendy schowana w zmiennej wypadałaby z tego
// wykazu i komenda wyglądałaby na niewpiętą.
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
//
// Żadna z nich nie rozgłasza `design.asset.changed`: ścieżka i symbol są bytami
// KOMPOZYCJI, a nie zasobami Assets Panelu. Rozgłoszenie zdarzenia zasobu po
// narysowaniu krzywej donosiłoby panelowi o zmianie zasobu, który się nie zmienił.
//
// `design.vector.export` też milczy — oddaje treść w odpowiedzi (`contentBase64`),
// a nie zasób w magazynie, więc nie ma czego rozgłosić. To ta sama zasada, którą
// jadą `design.asset.export` i `design.board.export`.
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
//
// Rozgłoszenie `design.board.changed` idzie po tych, które ZMIENIAJĄ układ
// warstw kompozycji: układ automatyczny, przeliczenie więzi po zmianie rozmiaru
// ramki, dołożenie instancji komponentu oraz obie drogi makiety. Drugie okno nad
// tą samą tablicą ma się dowiedzieć, że warstwy, które widzi, przestały być
// aktualne.
//
// Ramka, więz, siatka, komponent i przejście prototypu NIE rozgłaszają zdarzenia
// kompozycji: obszar `design.*` ma trzy zdarzenia — zasobu, kompozycji
// i obecności — a żaden z tych bytów nie jest kompozycją. Zdarzenie kompozycji po
// zapisie samej ramki mówiłoby, że zmienił się układ warstw, który się nie
// zmienił.
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

// zarejestrujDesignBarwy wpina siedem czynności barwy.
//
// Zdarzenie idzie po jednej: symulacja wady widzenia zakłada ZASÓB w magazynie,
// więc rozgłasza `design.asset.changed` z rodzajem `created`. Pozostałe sześć
// niczego w bazie zasobów nie zmienia — paleta, kontrast i przeliczenie są
// rachunkiem, a gradient jest bytem kompozycji, nie zasobem.
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

// zarejestrujDesignIkony wpina osiem czynności ikon i krojów.
//
// Pakiet ikon zakłada zasób magazynu i rozgłasza `created`. Ikony własne
// (`design.icon.set`, `design.icon.generate`) NIE rozgłaszają zdarzenia zasobu:
// ikona jest bytem osobnym od zasobu Assets Panelu — leży we własnej tabeli
// i kontrakt nie ma dla niej zdarzenia. Rozgłoszenie `design.asset.changed`
// z ikoną wymagałoby złożenia `DesignAsset` z bytu, który nim nie jest.
//
// `design.favicon.build` też milczy, choć zakłada zasoby: odpowiedź niesie same
// IDENTYFIKATORY (`assetIds`), a zdarzenie wymaga całego zasobu. Rozgłoszenie
// jednego z dziewięciu opisywałoby powstanie jednego i przemilczało resztę —
// ta sama zasada, którą jedzie generowanie wariantów.
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

// zarejestrujDesignMarketing wpina cztery czynności marketingowe.
//
// Wciągnięcie z bazy zdjęciowej i makieta produktowa zakładają zasób i oddają go
// w całości, więc rozgłaszają `created`. Komplet kampanii oddaje same
// identyfikatory rozmiarów, więc milczy — z tego samego powodu, co
// `design.favicon.build`. Wyszukanie niczego nie zmienia.
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
//
// Wydanie do druku, wykres i schemat zakładają zasób magazynu i oddają go
// w całości, więc rozgłaszają `created`. Podział na kafle oddaje wykaz kafli
// z identyfikatorami, więc milczy — jedno zdarzenie na kilkadziesiąt kafli
// opisywałoby powstanie jednego i przemilczało resztę.
//
// Kontrola przeddrukowa, profile i wykaz nośników niczego nie zmieniają.
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
// wskazanej kodem — odczytując jej stan PO zmianie.
//
// Odczyt jest ceną prawdy o skutku: zdarzenie niesie komplet warstw, a warstwy
// zmieniły się właśnie w tym wywołaniu. Zdarzenie złożone z żądania mówiłoby
// o zamiarze, nie o stanie, i drugie okno narysowałoby układ, którego w bazie nie
// ma. Niepowodzenie odczytu NIE unieważnia komendy — zmiana się już udała, a szyna
// zdarzeń milczy wtedy zamiast rozgłaszać nieprawdę.
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
// i historię promptów.
//
// Kolekcja i szablon nie rozgłaszają zdarzenia: obszar `design.*` ma trzy
// zdarzenia — zasobu, kompozycji i obecności — a kolekcja nie jest żadnym
// z tych bytów. Rozgłoszenie `design.asset.changed` po zmianie kolekcji
// donosiłoby panelowi o zmianie zasobu, który się nie zmienił.
func zarejestrujDesignKolekcjeIPrompty(r *Rejestr, m Design) {
	r.Zarejestruj(shared.CommandDesignCollectionCreate, obsluz(m.ZalozKolekcje))
	r.Zarejestruj(shared.CommandDesignCollectionAssign, obsluz(m.PrzypiszDoKolekcji))
	r.Zarejestruj(shared.CommandDesignCollectionList, obsluz(m.Kolekcje))

	r.Zarejestruj(shared.CommandDesignPromptTemplateSave, obsluz(m.ZapiszSzablonPromptu))
	r.Zarejestruj(shared.CommandDesignPromptTemplateList, obsluz(m.SzablonyPromptu))
	r.Zarejestruj(shared.CommandDesignPromptHistoryList, obsluz(m.HistoriaPromptow))
}

// zarejestrujDesignPlansze wpina wersjonowanie kompozycji, jej wyrys,
// adnotacje i obecność.
//
// Przywrócenie wersji ZMIENIA układ kompozycji, więc rozgłasza
// `design.board.changed` — drugie okno nad tą samą tablicą ma się dowiedzieć,
// że warstwy, które widzi, przestały być aktualne. Zapis wersji i wykaz wersji
// układu nie ruszają, więc milczą.
//
// Zgłoszenie obecności rozgłasza `design.board.presence` — po to jest.
// Zdarzenie niesie komplet obecnych, nie samą zmianę: odbiorca ma narysować
// wszystkie kursory, a nie doliczać stan z ciągu przyrostów, którego początku
// nie widział.
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

// zarejestrujDesignZetony wpina zestawy żetonów i przewodnik stylu.
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
