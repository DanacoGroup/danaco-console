// Odpowiedzialność pliku: oddanie treści zasobu z magazynu rdzenia
// (`design.asset.content.get`). Rdzeń nie oddaje bajtów zastępczych żadną
// drogą, a treść większa niż granica wołającego jest odmową, nie ucięciem.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"

	"danacoconsole/shared"
)

// TrescZasobu oddaje treść zasobu wraz z jej miarą i sumą kontrolną —
// obsługuje `design.asset.content.get`.
func (a *adapterDesignu) TrescZasobu(ctx context.Context,
	z shared.DesignAssetContentGetRequest) (shared.DesignAssetContentGetResponse, error) {

	if strings.TrimSpace(z.AssetId) == "" {
		return shared.DesignAssetContentGetResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.content.get bez wskazania zasobu")
	}
	zasob, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(z.AssetId))
	if err != nil {
		return shared.DesignAssetContentGetResponse{}, bladNieznanegoZasobuDesignu(z.AssetId, err)
	}
	if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
		return shared.DesignAssetContentGetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"zasób %s nie ma odwołania do treści — w magazynie serwera nie leżą jego bajty; "+
				"serwer nie oddaje treści zastępczej", z.AssetId))
	}

	bajty, err := os.ReadFile(*zasob.URI)
	if err != nil {
		return shared.DesignAssetContentGetResponse{}, bladWydaniaDesignu(fmt.Sprintf(
			"odwołanie zasobu %s nie prowadzi do treści w magazynie: %v", z.AssetId, err))
	}
	if len(bajty) == 0 {
		return shared.DesignAssetContentGetResponse{}, bladWydaniaDesignu(fmt.Sprintf(
			"treść zasobu %s ma zerową długość — plik pod odwołaniem istnieje, ale nie ma czego pokazać",
			z.AssetId))
	}

	// Granica sprawdza się po zmierzeniu treści, żeby odmowa podała wielkość.
	if z.MaxBytes != nil && *z.MaxBytes > 0 && len(bajty) > *z.MaxBytes {
		return shared.DesignAssetContentGetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"treść zasobu %s waży %d bajtów, a wołający przyjmuje najwyżej %d — "+
				"serwer odmawia zamiast oddać treść uciętą; po samo odsyłanie sięgnij "+
				"postacią odpowiedzi \"reference\"", z.AssetId, len(bajty), *z.MaxBytes))
	}

	suma := sha256.Sum256(bajty)
	odpowiedz := shared.DesignAssetContentGetResponse{
		AssetId:   zasob.Kod,
		MediaType: typTresciZasobuDesignu(zasob.Format, bajty),
		SizeBytes: len(bajty),
		Checksum:  hex.EncodeToString(suma[:]),
	}
	// Postać odsyłania oddaje miarę i sumę bez bajtów.
	if z.Disposition != nil && *z.Disposition == shared.AssetContentDispositionReference {
		return odpowiedz, nil
	}
	tresc := base64.StdEncoding.EncodeToString(bajty)
	odpowiedz.ContentBase64 = &tresc
	return odpowiedz, nil
}

// typTresciZasobuDesignu rozstrzyga typ treści wedle IANA. Pierwszeństwo ma
// format zapisany w wierszu; dopiero gdy go nie ma, rozstrzyga wykrywacz
// biblioteki standardowej. Zgadywania po nazwie pliku tu nie ma.
func typTresciZasobuDesignu(format *string, bajty []byte) string {
	if format != nil {
		switch normalizujFormatWydaniaDesignu(*format) {
		case "png":
			return "image/png"
		case "jpeg":
			return "image/jpeg"
		case "gif":
			return "image/gif"
		case "webp":
			return "image/webp"
		case "svg":
			return "image/svg+xml"
		case "pdf":
			return "application/pdf"
		case "ico":
			return "image/vnd.microsoft.icon"
		case "bmp":
			return "image/bmp"
		case "tiff":
			return "image/tiff"
		}
	}
	wykryty := http.DetectContentType(bajty)
	if wykryty == "" {
		return "application/octet-stream"
	}
	// DetectContentType dokleja parametr zestawu znaków; kontrakt chce typu
	// samego.
	if numer := strings.Index(wykryty, ";"); numer > 0 {
		return strings.TrimSpace(wykryty[:numer])
	}
	return wykryty
}
