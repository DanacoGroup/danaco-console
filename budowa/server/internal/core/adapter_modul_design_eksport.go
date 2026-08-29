// Odpowiedzialność pliku: wydanie zasobu Operatorowi jako pliku i wydanie
// partii zasobów w komplecie skal. Warsztat obrazu leży w
// `adapter_modul_design_obrazy.go`, bo dzieli go z wyrysem kompozycji.
package core

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// WydajZasob wydaje zasób w formacie i skali wskazanych żądaniem — obsługuje
// `design.asset.export`, oddając bajty gotowe do zapisania poza produktem.
func (a *adapterDesignu) WydajZasob(ctx context.Context,
	z shared.DesignAssetExportRequest) (shared.DesignAssetExportResponse, error) {

	if strings.TrimSpace(z.AssetId) == "" {
		return shared.DesignAssetExportResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.export bez wskazania zasobu")
	}
	if strings.TrimSpace(z.Format) == "" {
		return shared.DesignAssetExportResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.export bez formatu wydania; formaty wydania: " + formatyWydaniaDesignu)
	}
	if err := sprawdzFormatWydaniaDesignu("design.asset.export", z.Format); err != nil {
		return shared.DesignAssetExportResponse{}, err
	}
	if err := sprawdzSkaleWydaniaDesignu("design.asset.export", z.Scale); err != nil {
		return shared.DesignAssetExportResponse{}, err
	}

	zasob, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(z.AssetId))
	if err != nil {
		return shared.DesignAssetExportResponse{}, bladNieznanegoZasobuDesignu(z.AssetId, err)
	}

	wydanie, err := wydanieZasobuDesignu(zasob, z.Format, z.Scale, z.Quality)
	if err != nil {
		return shared.DesignAssetExportResponse{}, bladWydaniaDesignu(err.Error())
	}
	return shared.DesignAssetExportResponse{
		ContentBase64: wydanie.ContentBase64,
		FileName:      wydanie.FileName,
		MediaType:     wydanie.MediaType,
		SizeBytes:     wydanie.SizeBytes,
	}, nil
}

// WydajZasobyPartia wydaje wiele zasobów naraz, każdy w komplecie wskazanych
// skal — obsługuje `design.asset.export.batch`. Pole `archive` żądania nie
// jest tu obsługiwane spakowaniem, bo wynik niesie wykaz wydań osobnych.
func (a *adapterDesignu) WydajZasobyPartia(ctx context.Context,
	z shared.DesignAssetExportBatchRequest) (shared.DesignAssetExportBatchResponse, error) {

	if len(z.AssetIds) == 0 {
		return shared.DesignAssetExportBatchResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.export.batch bez wskazania zasobów")
	}
	if strings.TrimSpace(z.Format) == "" {
		return shared.DesignAssetExportBatchResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.export.batch bez formatu wydania; formaty wydania: " +
				formatyWydaniaDesignu)
	}
	if err := sprawdzFormatWydaniaDesignu("design.asset.export.batch", z.Format); err != nil {
		return shared.DesignAssetExportBatchResponse{}, err
	}
	// Skale sprawdza się wszystkie przed pierwszym odczytem zasobu.
	for _, skala := range z.Scales {
		wartosc := skala
		if err := sprawdzSkaleWydaniaDesignu("design.asset.export.batch", &wartosc); err != nil {
			return shared.DesignAssetExportBatchResponse{}, err
		}
	}

	skale := z.Scales
	if len(skale) == 0 {
		// Brak wskazania bierze skalę naturalną — jedno wydanie na zasób.
		skale = []float64{1}
	}

	wydania := []shared.DesignExportFile{}
	odrzucone := []shared.DesignExportRejection{}
	wydanych := 0

	for _, kod := range z.AssetIds {
		kod = strings.TrimSpace(kod)
		if kod == "" {
			continue
		}
		zasob, err := a.repozytorium.Zasob(ctx, kod)
		if err != nil {
			odrzucone = append(odrzucone, shared.DesignExportRejection{
				AssetId: kod,
				Reason:  "zasobu nie ma w magazynie serwera",
			})
			continue
		}
		// Zasób wydaje się we wszystkich skalach albo w żadnej.
		czesc := make([]shared.DesignExportFile, 0, len(skale))
		powod := ""
		for _, skala := range skale {
			wartosc := skala
			wydanie, err := wydanieZasobuDesignu(zasob, z.Format, &wartosc, z.Quality)
			if err != nil {
				powod = err.Error()
				break
			}
			wydanie.Scale = &wartosc
			czesc = append(czesc, wydanie)
		}
		if powod != "" {
			odrzucone = append(odrzucone, shared.DesignExportRejection{AssetId: kod, Reason: powod})
			continue
		}
		wydania = append(wydania, czesc...)
		wydanych++
	}

	odpowiedz := shared.DesignAssetExportBatchResponse{Files: wydania, Exported: wydanych}
	if len(odrzucone) > 0 {
		odpowiedz.Rejected = odrzucone
	}
	return odpowiedz, nil
}

// wydanieZasobuDesignu składa jedno wydanie: czyta bajty zasobu z magazynu,
// skaluje obraz i koduje go w żądanym formacie. SVG przechodzi bez
// rasteryzacji i bez skali.
func wydanieZasobuDesignu(zasob dane.ZasobDesignu, format string,
	skala *float64, jakosc *int) (shared.DesignExportFile, error) {

	if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
		return shared.DesignExportFile{}, fmt.Errorf(
			"zasób %s nie ma odwołania do treści — w magazynie nie leżą jego bajty", zasob.Kod)
	}
	nazwa := nazwaWydaniaDesignu(zasob, format, skala)

	if normalizujFormatWydaniaDesignu(format) == "svg" {
		bajty, err := os.ReadFile(*zasob.URI)
		if err != nil {
			return shared.DesignExportFile{}, fmt.Errorf(
				"odwołanie zasobu %s nie prowadzi do treści w magazynie: %w", zasob.Kod, err)
		}
		if !czyZapisWektorowyDesignu(bajty) {
			return shared.DesignExportFile{}, fmt.Errorf(
				"zasób %s nie jest zapisem wektorowym, a wydanie svg przepuszcza treść bez "+
					"przekształcenia — serwer nie obrysuje bitmapy ścieżkami; zamiana rastra na "+
					"ścieżki jest osobną czynnością (image.vectorize)", zasob.Kod)
		}
		return shared.DesignExportFile{
			AssetId:       zasob.Kod,
			FileName:      nazwa,
			ContentBase64: wBaza64Designu(bajty),
			MediaType:     typTresciWydaniaDesignu(format),
			SizeBytes:     len(bajty),
		}, nil
	}

	obraz, err := obrazZasobuDesignu(*zasob.URI)
	if err != nil {
		return shared.DesignExportFile{}, fmt.Errorf("zasób %s: %w", zasob.Kod, err)
	}
	if skala != nil {
		obraz = przeskalujObrazDesignu(obraz, *skala)
	}
	bajty, _, err := zakodujObrazDesignu(obraz, format, jakosc)
	if err != nil {
		return shared.DesignExportFile{}, fmt.Errorf("zasób %s: %w", zasob.Kod, err)
	}
	return shared.DesignExportFile{
		AssetId:       zasob.Kod,
		FileName:      nazwa,
		ContentBase64: wBaza64Designu(bajty),
		MediaType:     typTresciWydaniaDesignu(format),
		SizeBytes:     len(bajty),
	}, nil
}

// sprawdzFormatWydaniaDesignu odrzuca format, którego rdzeń nie zapisuje,
// przed odczytem czegokolwiek z magazynu. Odmowa wymienia formaty obsługiwane.
func sprawdzFormatWydaniaDesignu(komenda, format string) error {
	switch normalizujFormatWydaniaDesignu(format) {
	case "png", "jpeg", "svg", "pdf", "ico":
		return nil
	}
	return odmowaFormatuWydaniaDesignu(komenda, format)
}

// nazwaWydaniaDesignu składa proponowaną nazwę pliku: nazwę zasobu, gdy jest,
// albo kod zasobu. Wariant gęstości dokleja się z krotności skali (`@2x`).
func nazwaWydaniaDesignu(zasob dane.ZasobDesignu, format string, skala *float64) string {
	rdzen := zasob.Kod
	if zasob.Nazwa != nil && strings.TrimSpace(*zasob.Nazwa) != "" {
		rdzen = strings.TrimSpace(*zasob.Nazwa)
	}
	rdzen = oczyscNazwePlikuDesignu(rdzen)
	if skala != nil && *skala != 1 {
		rdzen += "@" + strconv.FormatFloat(*skala, 'g', -1, 64) + "x"
	}
	return rdzen + "." + rozszerzenieWydaniaDesignu(format)
}

// oczyscNazwePlikuDesignu zdejmuje z nazwy znaki, którymi dałoby się wyjść
// z katalogu zapisu u Operatora, bo klient zapisuje plik pod nazwą wprost.
func oczyscNazwePlikuDesignu(nazwa string) string {
	nazwa = strings.Map(func(znak rune) rune {
		switch znak {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|', 0:
			return '-'
		}
		return znak
	}, nazwa)
	nazwa = strings.TrimLeft(nazwa, ". ")
	nazwa = strings.TrimSpace(nazwa)
	if nazwa == "" {
		return "wydanie"
	}
	return nazwa
}

// czyZapisWektorowyDesignu rozstrzyga, czy bajty są zapisem SVG, po treści,
// nie po kolumnie `format`, sprawdzeniem pierwszych bajtów zasobu.
func czyZapisWektorowyDesignu(bajty []byte) bool {
	granica := len(bajty)
	if granica > 512 {
		granica = 512
	}
	poczatek := strings.ToLower(string(bajty[:granica]))
	return strings.Contains(poczatek, "<svg")
}
