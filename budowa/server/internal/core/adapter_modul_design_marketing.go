// Odpowiedzialność pliku: cztery czynności marketingowe modułu Design —
// komplet rozmiarów kampanii, bazy zdjęciowe i makieta produktowa. Dostawcy
// baz leżą w `adapter_modul_design_bazy_zdjeciowe.go`.
package core

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ZbudujKompletKampanii wydaje kompozycję w każdym zamówionym rozmiarze —
// obsługuje `design.campaign.set.build`.
func (a *adapterDesignu) ZbudujKompletKampanii(ctx context.Context,
	z shared.DesignCampaignSetBuildRequest) (shared.DesignCampaignSetBuildResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignCampaignSetBuildResponse{}, bladWskazaniaDesignu(
			"komenda design.campaign.set.build bez wskazania kompozycji")
	}
	if len(z.Sizes) == 0 {
		return shared.DesignCampaignSetBuildResponse{}, bladWskazaniaDesignu(
			"komenda design.campaign.set.build bez ani jednego rozmiaru: serwer nie wymyśla, " +
				"w jakich formatach ma wyjść kampania")
	}
	format := "png"
	if z.Format != nil && strings.TrimSpace(*z.Format) != "" {
		format = normalizujFormatWydaniaDesignu(*z.Format)
		switch format {
		case "png", "jpeg", "pdf", "ico":
		default:
			return shared.DesignCampaignSetBuildResponse{}, odmowaFormatuWydaniaDesignu(
				"design.campaign.set.build", *z.Format)
		}
	}
	for numer, rozmiar := range z.Sizes {
		if strings.TrimSpace(rozmiar.Name) == "" {
			return shared.DesignCampaignSetBuildResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"rozmiar numer %d bez nazwy: bez niej nie da się powiedzieć, który materiał "+
					"się nie udał", numer+1))
		}
		if rozmiar.Width <= 0 || rozmiar.Height <= 0 {
			return shared.DesignCampaignSetBuildResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"rozmiar %s ma bok %v×%v: materiał o niedodatnim boku nie istnieje",
				rozmiar.Name, rozmiar.Width, rozmiar.Height))
		}
	}

	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignCampaignSetBuildResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	warstwy, err := a.repozytorium.Warstwy(ctx, kompozycja.ID)
	if err != nil {
		return shared.DesignCampaignSetBuildResponse{}, bladDesignu(err)
	}
	kafle, pominietych, err := a.kafleWyrysuDesignu(ctx, warstwy)
	if err != nil {
		return shared.DesignCampaignSetBuildResponse{}, err
	}
	if len(kafle) == 0 {
		return shared.DesignCampaignSetBuildResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"kompozycja %s nie ma ani jednej warstwy z bajtami (warstw pominiętych: %d) — "+
				"komplet kampanii wyszedłby jako zestaw pustych prostokątów",
			kompozycja.Kod, pominietych))
	}
	zrodlowy := obszarWyrysuDesignu(kafle, nil)
	okno := oknoWytworu(z.WindowId, kompozycja.Okno)

	wydane := make([]shared.DesignCampaignSize, 0, len(z.Sizes))
	nieudane := []string{}
	for _, rozmiar := range z.Sizes {
		// Skala bierze mniejszy ze współczynników, więc materiał wchodzi cały.
		skala := mniejszaDesignu(rozmiar.Width/zrodlowy.Width, rozmiar.Height/zrodlowy.Height)
		if skala <= 0 {
			nieudane = append(nieudane, rozmiar.Name)
			continue
		}
		plotno := zlozWyrysRastrowyDesignu(kafle, zrodlowy, skala)
		// Płótno docelowe ma dokładny rozmiar zamówiony, wyrys ląduje w środku.
		docelowe := image.NewRGBA(image.Rect(0, 0,
			int(rozmiar.Width+0.5), int(rozmiar.Height+0.5)))
		granice := plotno.Bounds()
		odsuniecieX := (docelowe.Bounds().Dx() - granice.Dx()) / 2
		odsuniecieY := (docelowe.Bounds().Dy() - granice.Dy()) / 2
		draw.Draw(docelowe, granice.Add(image.Pt(odsuniecieX, odsuniecieY)),
			plotno, granice.Min, draw.Over)

		bajty, _, err := zakodujObrazDesignu(docelowe, format, nil)
		if err != nil {
			nieudane = append(nieudane, rozmiar.Name)
			continue
		}
		nazwa := fmt.Sprintf("%s — %s (%.0f×%.0f)", nazwaKompozycjiDoPlikuDesignu(kompozycja),
			rozmiar.Name, rozmiar.Width, rozmiar.Height)
		zasob, err := a.zalozZasobZBajtowDesignu(ctx, okno, nazwa,
			shared.DesignAssetKindImage, format, bajty)
		if err != nil {
			nieudane = append(nieudane, rozmiar.Name)
			continue
		}
		kod := zasob.Kod
		wydane = append(wydane, shared.DesignCampaignSize{
			Name: rozmiar.Name, Width: rozmiar.Width, Height: rozmiar.Height, AssetId: &kod,
		})
	}
	if len(wydane) == 0 {
		return shared.DesignCampaignSetBuildResponse{}, bladWydaniaDesignu(fmt.Sprintf(
			"nie udało się wydać ani jednego z %d zamówionych rozmiarów", len(z.Sizes)))
	}
	return shared.DesignCampaignSetBuildResponse{
		Sizes: wydane, Exported: len(wydane), FailedSizes: uporzadkujBilansDesignu(nieudane),
	}, nil
}

// SzukajWBazachZdjeciowych pyta dostawców o zasoby spełniające frazę —
// obsługuje `design.stock.search`.
func (a *adapterDesignu) SzukajWBazachZdjeciowych(ctx context.Context,
	z shared.DesignStockSearchRequest) (shared.DesignStockSearchResponse, error) {

	if strings.TrimSpace(z.Query) == "" {
		return shared.DesignStockSearchResponse{}, bladWskazaniaDesignu(
			"komenda design.stock.search bez frazy wyszukiwania")
	}
	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignStockSearchResponse{}, bladWskazaniaDesignu(
			"komenda design.stock.search bez wskazania okna")
	}
	limit := domyslnyLimitWynikowDostawcyDesignu
	if z.Limit != nil {
		if *z.Limit < 1 || *z.Limit > granicaWynikowDostawcyDesignu {
			return shared.DesignStockSearchResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.stock.search z granicą %d: serwer pyta o od 1 do %d zasobów na "+
					"dostawcę", *z.Limit, granicaWynikowDostawcyDesignu))
		}
		limit = *z.Limit
	}

	pytani := dostawcyZdjecDesignu()
	if z.Provider != nil && strings.TrimSpace(*z.Provider) != "" {
		dostawca, znany := dostawcaZdjecPoNazwieDesignu(*z.Provider)
		if !znany {
			return shared.DesignStockSearchResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.stock.search z dostawcą %q, którego serwer nie zna; dostawcy: %s",
				*z.Provider, strings.Join(nazwyDostawcowZdjecDesignu(), ", ")))
		}
		pytani = []dostawcaZdjecDesignu{dostawca}
	}

	klient := klientDostawcyZdjecDesignu()
	zapytani := []string{}
	nieudani := []string{}
	zasoby := []shared.DesignStockAsset{}
	for _, dostawca := range pytani {
		zapytani = append(zapytani, dostawca.Nazwa)
		klucz := ""
		if dostawca.WymagaKlucza {
			klucz = a.kluczDostawcyZdjecDesignu(ctx, dostawca.Nazwa)
			if klucz == "" {
				// Brak klucza nie jest odmową całej komendy, tylko tego dostawcy.
				nieudani = append(nieudani, fmt.Sprintf(
					"%s (brak klucza w sejfie pod bytem %s%s)",
					dostawca.Nazwa, przedrostekBytuSejfuZdjecDesignu, dostawca.Nazwa))
				continue
			}
		}
		wyniki, err := dostawca.Szukaj(ctx, klient, klucz, strings.TrimSpace(z.Query), limit)
		if err != nil {
			nieudani = append(nieudani, fmt.Sprintf("%s (%s)", dostawca.Nazwa, err.Error()))
			continue
		}
		for _, wynik := range wyniki {
			if strings.TrimSpace(wynik.Identyfikator) == "" {
				continue
			}
			zasoby = append(zasoby, zasobBazyKontraktuDesignu(dostawca.Nazwa, wynik))
		}
	}

	// Odmowa należy się wyłącznie sytuacji, w której nikt nie odpowiedział.
	if len(zasoby) == 0 && len(nieudani) == len(zapytani) {
		return shared.DesignStockSearchResponse{}, bladDesignu(fmt.Errorf(
			"żaden z zapytanych dostawców nie odpowiedział: %s", strings.Join(nieudani, "; ")))
	}
	return shared.DesignStockSearchResponse{
		Assets: zasoby, ProvidersQueried: zapytani,
		ProvidersFailed: uporzadkujBilansDesignu(nieudani),
	}, nil
}

// zasobBazyKontraktuDesignu składa `DesignStockAsset` kontraktu z zasobu
// dostawcy bazy zdjęciowej, przenosząc tytuł, podgląd i licencję, gdy istnieją.
func zasobBazyKontraktuDesignu(dostawca string,
	wynik zasobDostawcyZdjecDesignu) shared.DesignStockAsset {

	zasob := shared.DesignStockAsset{Id: wynik.Identyfikator, Provider: dostawca}
	if wynik.Tytul != "" {
		tytul := wynik.Tytul
		zasob.Title = &tytul
	}
	if podglad := pierwszyNiepustyDesignu(wynik.Podglad, wynik.Pelny); podglad != "" {
		zasob.PreviewUrl = &podglad
	}
	if wynik.Licencja != "" {
		licencja := wynik.Licencja
		zasob.License = &licencja
	}
	if wynik.Autor != "" {
		autor := wynik.Autor
		zasob.Author = &autor
	}
	return zasob
}

// WciagnijZBazyZdjeciowej wciąga zasób dostawcy do magazynu rdzenia wraz
// z licencją — obsługuje `design.stock.import`. Licencja zapisuje się razem
// z zasobem i jej brak jest odmową.
func (a *adapterDesignu) WciagnijZBazyZdjeciowej(ctx context.Context,
	z shared.DesignStockImportRequest) (shared.DesignStockImportResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignStockImportResponse{}, bladWskazaniaDesignu(
			"komenda design.stock.import bez wskazania okna")
	}
	if strings.TrimSpace(z.ExternalId) == "" {
		return shared.DesignStockImportResponse{}, bladWskazaniaDesignu(
			"komenda design.stock.import bez identyfikatora zasobu u dostawcy")
	}
	dostawca, znany := dostawcaZdjecPoNazwieDesignu(z.Provider)
	if !znany {
		return shared.DesignStockImportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.stock.import z dostawcą %q, którego serwer nie zna; dostawcy: %s",
			z.Provider, strings.Join(nazwyDostawcowZdjecDesignu(), ", ")))
	}
	klucz := ""
	if dostawca.WymagaKlucza {
		klucz = a.kluczDostawcyZdjecDesignu(ctx, dostawca.Nazwa)
		if klucz == "" {
			return shared.DesignStockImportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"dostawca %s wymaga klucza, a w sejfie serwera nie ma wpisu pod bytem %s%s — "+
					"wciągnięcie pojedynczego zasobu nie ma czym zapytać; dostawcy pracujący bez "+
					"klucza: openverse, wikimedia, met, nasa",
				dostawca.Nazwa, przedrostekBytuSejfuZdjecDesignu, dostawca.Nazwa))
		}
	}

	klient := klientDostawcyZdjecDesignu()
	opis, err := dostawca.Pobierz(ctx, klient, klucz, strings.TrimSpace(z.ExternalId))
	if err != nil {
		return shared.DesignStockImportResponse{}, bladDesignu(fmt.Errorf(
			"dostawca %s nie oddał opisu zasobu %s: %w", dostawca.Nazwa, z.ExternalId, err))
	}
	if strings.TrimSpace(opis.Pelny) == "" {
		return shared.DesignStockImportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"dostawca %s nie ma dla zasobu %s adresu pliku — z samego opisu zasób nie powstanie",
			dostawca.Nazwa, z.ExternalId))
	}
	if strings.TrimSpace(opis.Licencja) == "" {
		return shared.DesignStockImportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"dostawca %s nie podał licencji zasobu %s — serwer nie wciąga materiału, o którym "+
				"nikt później nie powie, czy wolno go było użyć",
			dostawca.Nazwa, z.ExternalId))
	}

	bajty, typTresci, err := pobierzBajtyZdjeciaDesignu(ctx, klient, opis.Pelny)
	if err != nil {
		return shared.DesignStockImportResponse{}, bladDesignu(fmt.Errorf(
			"bajtów zasobu %s u dostawcy %s nie da się wciągnąć: %w",
			z.ExternalId, dostawca.Nazwa, err))
	}
	format := "jpeg"
	if zTypu := formatZTypuTresci(typTresci); zTypu != nil {
		format = *zTypu
	}
	nazwa := pierwszyNiepustyDesignu(opis.Tytul, dostawca.Nazwa+" "+opis.Identyfikator)
	zapisany, err := a.zalozZasobZBajtowDesignu(ctx, strings.TrimSpace(z.WindowId), nazwa,
		shared.DesignAssetKindImage, format, bajty)
	if err != nil {
		return shared.DesignStockImportResponse{}, err
	}

	licencja := opis.Licencja
	wiersz := dane.LicencjaZasobuDesignu{
		ZasobID:                zapisany.ID,
		Dostawca:               dostawca.Nazwa,
		IdentyfikatorUDostawcy: opis.Identyfikator,
		Licencja:               &licencja,
	}
	if opis.Autor != "" {
		autor := opis.Autor
		wiersz.Autor = &autor
	}
	if opis.Odsylacz != "" {
		odsylacz := opis.Odsylacz
		wiersz.Odsylacz = &odsylacz
	}
	if err := a.repozytorium.ZapiszLicencjeZasobuDesignu(ctx, wiersz); err != nil {
		// Zasób bez zapisanej licencji jest usterką: wiersz zasobu znika.
		if _, bladUsuniecia := a.repozytorium.UsunZasob(ctx, zapisany.Kod); bladUsuniecia != nil {
			return shared.DesignStockImportResponse{}, bladDesignu(fmt.Errorf(
				"licencji zasobu %s nie udało się zapisać (%w), a wiersza zasobu nie udało się "+
					"wycofać (%v) — zasób bez prowenancji został w bazie i wymaga ręcznego usunięcia",
				zapisany.Kod, err, bladUsuniecia))
		}
		return shared.DesignStockImportResponse{}, bladDesignu(fmt.Errorf(
			"licencji zasobu z bazy %s nie udało się zapisać, więc zasób został wycofany: %w",
			dostawca.Nazwa, err))
	}
	return shared.DesignStockImportResponse{
		Asset: zasobWytworzonyKontraktu(zapisany), License: &licencja,
	}, nil
}

// WyrysujMakieteProduktowa nakłada projekt na zdjęcie produktu — obsługuje
// `design.product.mockup.render`.
func (a *adapterDesignu) WyrysujMakieteProduktowa(ctx context.Context,
	z shared.DesignProductMockupRenderRequest) (shared.DesignProductMockupRenderResponse, error) {

	if z.Width <= 0 || z.Height <= 0 {
		return shared.DesignProductMockupRenderResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.product.mockup.render z nałożeniem %v×%v: obszar o niedodatnim boku "+
				"nie jest nałożeniem", z.Width, z.Height))
	}
	if z.Opacity != nil && (*z.Opacity < 0 || *z.Opacity > 1) {
		return shared.DesignProductMockupRenderResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.product.mockup.render z kryciem %v: krycie idzie od 0 do 1", *z.Opacity))
	}
	projekt, wierszProjektu, err := a.obrazZasobuPoKodzieDesignu(ctx,
		"design.product.mockup.render", z.DesignAssetId)
	if err != nil {
		return shared.DesignProductMockupRenderResponse{}, err
	}
	produkt, wierszProduktu, err := a.obrazZasobuPoKodzieDesignu(ctx,
		"design.product.mockup.render", z.ProductAssetId)
	if err != nil {
		return shared.DesignProductMockupRenderResponse{}, err
	}

	wynik := nalozProjektNaProduktDesignu(produkt, projekt, z)
	bajty, _, err := zakodujObrazDesignu(wynik, "png", nil)
	if err != nil {
		return shared.DesignProductMockupRenderResponse{}, bladWydaniaDesignu(err.Error())
	}
	nazwa := nazwaZasobuDesignu(wierszProduktu) + " z " + nazwaZasobuDesignu(wierszProjektu)
	okno := oknoWytworu(z.WindowId, wierszProduktu.Okno)
	zapisany, err := a.zalozZasobZBajtowDesignu(ctx, okno, nazwa,
		shared.DesignAssetKindImage, "png", bajty)
	if err != nil {
		return shared.DesignProductMockupRenderResponse{}, err
	}
	return shared.DesignProductMockupRenderResponse{
		Asset: zasobWytworzonyKontraktu(zapisany),
	}, nil
}

// nalozProjektNaProduktDesignu nakłada projekt na zdjęcie produktu z obrotem
// i kryciem. Zdjęcie produktu wchodzi w całości i bez skalowania: to ono jest
// tłem, a Operator wskazuje położenie nałożenia w jego pikselach.
func nalozProjektNaProduktDesignu(produkt, projekt image.Image,
	z shared.DesignProductMockupRenderRequest) image.Image {

	granicePr := produkt.Bounds()
	plotno := image.NewRGBA(image.Rect(0, 0, granicePr.Dx(), granicePr.Dy()))
	draw.Draw(plotno, plotno.Bounds(), produkt, granicePr.Min, draw.Src)

	graniceProjektu := projekt.Bounds()
	if graniceProjektu.Dx() == 0 || graniceProjektu.Dy() == 0 {
		return plotno
	}

	// Macierz przekształcenia: skalowanie, potem obrót, potem przesunięcie.
	skalaX := z.Width / float64(graniceProjektu.Dx())
	skalaY := z.Height / float64(graniceProjektu.Dy())
	kat := 0.0
	if z.Rotation != nil {
		kat = *z.Rotation * math.Pi / 180
	}
	cos, sin := math.Cos(kat), math.Sin(kat)
	srodekX, srodekY := z.X+z.Width/2, z.Y+z.Height/2
	polowaSzerokosci, polowaWysokosci := z.Width/2, z.Height/2

	// Macierz przejścia z układu źródła do płótna, zapisana wierszami.
	macierz := f64.Aff3{
		skalaX * cos, -skalaY * sin,
		srodekX - (cos*polowaSzerokosci - sin*polowaWysokosci),
		skalaX * sin, skalaY * cos,
		srodekY - (sin*polowaSzerokosci + cos*polowaWysokosci),
	}

	// Krycie częściowe idzie przez maskę jednolitą w nastawach transformatora.
	nastawy := &draw.Options{}
	if z.Opacity != nil && *z.Opacity < 1 {
		nastawy.SrcMask = image.NewUniform(przezroczystoscDesignu(*z.Opacity))
	}
	draw.CatmullRom.Transform(plotno, macierz, projekt, graniceProjektu, draw.Over, nastawy)
	return plotno
}

// przezroczystoscDesignu składa jednolitą maskę krycia. Maska ma jeden kanał
// (alfa), bo krycie nie zmienia barw nałożenia — zmienia jego przezroczystość.
func przezroczystoscDesignu(krycie float64) color.Alpha16 {
	return color.Alpha16{A: uint16(przytnijUlamekDesignu(krycie) * 65535)}
}
