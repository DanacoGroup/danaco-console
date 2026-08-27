// Odpowiedzialność pliku: wyrys kompozycji Design Board do jednego pliku —
// design.board.export; warsztat obrazu leży osobno.
package core

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"os"
	"strings"

	"golang.org/x/image/draw"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// domyslnyBokWyrysuDesignu jest bokiem płótna dla kompozycji, której
// warstwy nie niosą ani położenia, ani wymiarów.
const domyslnyBokWyrysuDesignu = 1024

// WyrysujKompozycje wyrysowuje całą kompozycję albo wskazany obszar do
// jednego pliku — obsługuje design.board.export.
func (a *adapterDesignu) WyrysujKompozycje(ctx context.Context,
	z shared.DesignBoardExportRequest) (shared.DesignBoardExportResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignBoardExportResponse{}, bladWskazaniaDesignu(
			"komenda design.board.export bez wskazania kompozycji")
	}
	format := normalizujFormatWydaniaDesignu(z.Format)
	switch format {
	case "png", "pdf", "svg":
	default:
		return shared.DesignBoardExportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.board.export z formatem %q: wyrys kompozycji wychodzi jako png, pdf albo svg",
			z.Format))
	}
	if err := sprawdzSkaleWydaniaDesignu("design.board.export", z.Scale); err != nil {
		return shared.DesignBoardExportResponse{}, err
	}
	if err := sprawdzObszarWyrysuDesignu(z.Region); err != nil {
		return shared.DesignBoardExportResponse{}, err
	}

	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignBoardExportResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	warstwy, err := a.repozytorium.Warstwy(ctx, kompozycja.ID)
	if err != nil {
		return shared.DesignBoardExportResponse{}, bladDesignu(err)
	}

	kafle, pominietych, err := a.kafleWyrysuDesignu(ctx, warstwy)
	if err != nil {
		return shared.DesignBoardExportResponse{}, err
	}
	if len(kafle) == 0 {
		return shared.DesignBoardExportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"kompozycja %s nie ma ani jednej warstwy z bajtami do wyrysowania "+
				"(warstw pominiętych: %d) — rdzeń odmawia zamiast oddać pusty prostokąt",
			kompozycja.Kod, pominietych))
	}

	obszar := obszarWyrysuDesignu(kafle, z.Region)
	skala := 1.0
	if z.Scale != nil && *z.Scale > 0 {
		skala = *z.Scale
	}
	nazwa := nazwaWyrysuKompozycjiDesignu(kompozycja, format)

	if format == "svg" {
		tresc := zlozWyrysSvgDesignu(kafle, obszar, skala)
		return shared.DesignBoardExportResponse{
			ContentBase64: wBaza64Designu([]byte(tresc)),
			FileName:      nazwa,
			MediaType:     typTresciWydaniaDesignu("svg"),
		}, nil
	}

	plotno := zlozWyrysRastrowyDesignu(kafle, obszar, skala)
	bajty, typTresci, err := zakodujObrazDesignu(plotno, format, nil)
	if err != nil {
		return shared.DesignBoardExportResponse{}, bladWydaniaDesignu(err.Error())
	}
	return shared.DesignBoardExportResponse{
		ContentBase64: wBaza64Designu(bajty),
		FileName:      nazwa,
		MediaType:     typTresci,
	}, nil
}

// kafelWyrysuDesignu to jedna warstwa gotowa do wyrysowania: jej obraz
// wraz z prostokątem, jaki zajmuje w jednostkach kompozycji.
type kafelWyrysuDesignu struct {
	obraz     image.Image
	bajty     []byte
	x         float64
	y         float64
	szerokosc float64
	wysokosc  float64
}

// kafleWyrysuDesignu zamienia warstwy kompozycji na kafle wyrysu i oddaje
// liczbę warstw pominiętych, wliczoną do bilansu odmowy.
func (a *adapterDesignu) kafleWyrysuDesignu(ctx context.Context,
	warstwy []dane.WarstwaKompozycji) ([]kafelWyrysuDesignu, int, error) {

	kafle := make([]kafelWyrysuDesignu, 0, len(warstwy))
	pominietych := 0
	for _, warstwa := range warstwy {
		if warstwa.ZasobID == nil || strings.TrimSpace(*warstwa.ZasobID) == "" {
			pominietych++
			continue
		}
		zasob, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(*warstwa.ZasobID))
		if err != nil {
			if czyBrakZasobuDesignu(err) {
				pominietych++
				continue
			}
			return nil, 0, bladDesignu(err)
		}
		if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
			pominietych++
			continue
		}
		obraz, err := obrazZasobuDesignu(*zasob.URI)
		if err != nil {
			// Zasób wektorowy albo dokument w warstwie: rdzeń go nie rasteryzuje.
			// To pominięcie, nie awaria.
			pominietych++
			continue
		}
		bajty, err := os.ReadFile(*zasob.URI)
		if err != nil {
			pominietych++
			continue
		}
		granice := obraz.Bounds()
		kafel := kafelWyrysuDesignu{
			obraz:     obraz,
			bajty:     bajty,
			szerokosc: float64(granice.Dx()),
			wysokosc:  float64(granice.Dy()),
		}
		if warstwa.X != nil {
			kafel.x = *warstwa.X
		}
		if warstwa.Y != nil {
			kafel.y = *warstwa.Y
		}
		// Wymiary warstwy biją wymiary obrazu: rozmiar nadany na kanwie bije
		// rozmiar powstania obrazu.
		if warstwa.Szerokosc != nil && *warstwa.Szerokosc > 0 {
			kafel.szerokosc = *warstwa.Szerokosc
		}
		if warstwa.Wysokosc != nil && *warstwa.Wysokosc > 0 {
			kafel.wysokosc = *warstwa.Wysokosc
		}
		kafle = append(kafle, kafel)
	}
	return kafle, pominietych, nil
}

// obszarWyrysuDesignu rozstrzyga, co wchodzi w kadr: obszar wskazany
// żądaniem albo prostokąt obejmujący wszystkie kafle.
func obszarWyrysuDesignu(kafle []kafelWyrysuDesignu,
	obszar *shared.DesignBoardRegion) shared.DesignBoardRegion {

	if obszar != nil {
		return *obszar
	}
	lewa, gora := kafle[0].x, kafle[0].y
	prawa, dol := kafle[0].x+kafle[0].szerokosc, kafle[0].y+kafle[0].wysokosc
	for _, kafel := range kafle[1:] {
		lewa = mniejszaDesignu(lewa, kafel.x)
		gora = mniejszaDesignu(gora, kafel.y)
		prawa = wiekszaDesignu(prawa, kafel.x+kafel.szerokosc)
		dol = wiekszaDesignu(dol, kafel.y+kafel.wysokosc)
	}
	szerokosc, wysokosc := prawa-lewa, dol-gora
	if szerokosc <= 0 {
		szerokosc = domyslnyBokWyrysuDesignu
	}
	if wysokosc <= 0 {
		wysokosc = domyslnyBokWyrysuDesignu
	}
	return shared.DesignBoardRegion{X: lewa, Y: gora, Width: szerokosc, Height: wysokosc}
}

// sprawdzObszarWyrysuDesignu odrzuca obszar o niedodatnim boku PRZED
// odczytem czegokolwiek, bo to pomyłka w żądaniu.
func sprawdzObszarWyrysuDesignu(obszar *shared.DesignBoardRegion) error {
	if obszar == nil {
		return nil
	}
	if obszar.Width <= 0 || obszar.Height <= 0 {
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.board.export z obszarem %v×%v: obszar o niedodatnim boku nie jest kadrem",
			obszar.Width, obszar.Height))
	}
	return nil
}

// zlozWyrysRastrowyDesignu składa warstwy na jedno płótno w kolejności
// renderowania, wziętej wprost z kolejności wykazu.
func zlozWyrysRastrowyDesignu(kafle []kafelWyrysuDesignu,
	obszar shared.DesignBoardRegion, skala float64) image.Image {

	szerokosc := int(obszar.Width*skala + 0.5)
	wysokosc := int(obszar.Height*skala + 0.5)
	if szerokosc < 1 {
		szerokosc = 1
	}
	if wysokosc < 1 {
		wysokosc = 1
	}
	plotno := image.NewRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for _, kafel := range kafle {
		lewa := int((kafel.x-obszar.X)*skala + 0.5)
		gora := int((kafel.y-obszar.Y)*skala + 0.5)
		prawa := int((kafel.x-obszar.X+kafel.szerokosc)*skala + 0.5)
		dol := int((kafel.y-obszar.Y+kafel.wysokosc)*skala + 0.5)
		cel := image.Rect(lewa, gora, prawa, dol)
		if cel.Empty() {
			continue
		}
		draw.CatmullRom.Scale(plotno, cel, kafel.obraz, kafel.obraz.Bounds(), draw.Over, nil)
	}
	return plotno
}

// zlozWyrysSvgDesignu składa wyrys wektorowy: jeden element image na
// warstwę, z treścią osadzoną jako dane wprost w dokumencie.
func zlozWyrysSvgDesignu(kafle []kafelWyrysuDesignu,
	obszar shared.DesignBoardRegion, skala float64) string {

	szerokosc := obszar.Width * skala
	wysokosc := obszar.Height * skala

	var dokument strings.Builder
	dokument.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	fmt.Fprintf(&dokument,
		`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" `+
			`width="%g" height="%g" viewBox="0 0 %g %g">`+"\n",
		szerokosc, wysokosc, szerokosc, wysokosc)
	for _, kafel := range kafle {
		fmt.Fprintf(&dokument,
			`  <image x="%g" y="%g" width="%g" height="%g" preserveAspectRatio="none" `+
				`xlink:href="data:%s;base64,%s"/>`+"\n",
			(kafel.x-obszar.X)*skala, (kafel.y-obszar.Y)*skala,
			kafel.szerokosc*skala, kafel.wysokosc*skala,
			typTresciOsadzenegoDesignu(kafel.bajty), wBaza64Designu(kafel.bajty))
	}
	dokument.WriteString("</svg>\n")
	return dokument.String()
}

// typTresciOsadzenegoDesignu rozstrzyga typ treści osadzanej w dokumencie
// SVG z samych bajtów; kolumna format tu nie wystarcza.
func typTresciOsadzenegoDesignu(bajty []byte) string {
	switch {
	case bytes.HasPrefix(bajty, []byte("\x89PNG\r\n\x1a\n")):
		return "image/png"
	case bytes.HasPrefix(bajty, []byte("\xff\xd8\xff")):
		return "image/jpeg"
	case bytes.HasPrefix(bajty, []byte("GIF8")):
		return "image/gif"
	case len(bajty) > 12 && bytes.Equal(bajty[0:4], []byte("RIFF")) &&
		bytes.Equal(bajty[8:12], []byte("WEBP")):
		return "image/webp"
	}
	return "application/octet-stream"
}

// mniejszaDesignu i wiekszaDesignu wybierają skrajną z dwóch liczb; własne
// pomocniki, bo pakiet core przesłania wbudowane min i max.
func mniejszaDesignu(pierwsza, druga float64) float64 {
	if druga < pierwsza {
		return druga
	}
	return pierwsza
}

func wiekszaDesignu(pierwsza, druga float64) float64 {
	if druga > pierwsza {
		return druga
	}
	return pierwsza
}

// nazwaWyrysuKompozycjiDesignu składa proponowaną nazwę pliku wyrysu —
// z nazwy kompozycji, gdy Operator ją nadał, albo z jej kodu.
func nazwaWyrysuKompozycjiDesignu(kompozycja dane.KompozycjaDesignu, format string) string {
	rdzen := kompozycja.Kod
	if kompozycja.Nazwa != nil && strings.TrimSpace(*kompozycja.Nazwa) != "" {
		rdzen = strings.TrimSpace(*kompozycja.Nazwa)
	}
	return oczyscNazwePlikuDesignu(rdzen) + "." + rozszerzenieWydaniaDesignu(format)
}
