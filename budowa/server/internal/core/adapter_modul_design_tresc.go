// Odpowiedzialność pliku: oddanie TREŚCI zasobu z magazynu rdzenia
// (`design.asset.content.get`). Metody stoją na `*adapterDesignu`
// (`adapter_modul_design.go`).
//
// Ta komenda domyka lukę, przez którą moduł miał precedens szkody. Pole `uri`
// zasobu jest ścieżką w systemie plików rdzenia — przeglądarka nie wczyta spod
// niego niczego, także wtedy, gdy zasób powstał bez zarzutu. Bez tej drogi
// Assets Panel pokazywał kafelek z odwołaniem, którego nie da się otworzyć,
// i wyglądało to identycznie jak zasób bez bajtów.
//
// Rdzeń nie oddaje bajtów zastępczych ŻADNĄ drogą. Zasób nieznany to odmowa
// `not_found`. Zasób bez odwołania — odmowa nazywająca brak treści. Odwołanie
// prowadzące donikąd — odmowa nazywająca odwołanie. Treść większa niż granica
// wołającego — odmowa PODAJĄCA ZMIERZONĄ WIELKOŚĆ, nigdy treść ucięta: klient,
// który dostałby połowę pliku ze stanem `ok`, zapisałby ją jako plik cały.
//
// Suma kontrolna liczy się z BAJTÓW ODCZYTANYCH, nie z nazwy bloba. Nazwą bloba
// jest wprawdzie suma jego zawartości, więc obie wartości powinny być równe —
// i właśnie dlatego liczymy je osobno: rozjazd znaczy, że plik pod odwołaniem
// przestał być tym, za który się podaje, a przemilczenie tego byłoby oddaniem
// cudzej treści pod nazwą zasobu.
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
			"zasób %s nie ma odwołania do treści — w magazynie rdzenia nie leżą jego bajty; "+
				"rdzeń nie oddaje treści zastępczej", z.AssetId))
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

	// Granica wołającego sprawdza się PRZED złożeniem odpowiedzi, ale PO
	// zmierzeniu treści: odmowa ma podać wielkość rzeczywistą, a nie samą
	// wiadomość o przekroczeniu.
	if z.MaxBytes != nil && *z.MaxBytes > 0 && len(bajty) > *z.MaxBytes {
		return shared.DesignAssetContentGetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"treść zasobu %s waży %d bajtów, a wołający przyjmuje najwyżej %d — "+
				"rdzeń odmawia zamiast oddać treść uciętą; po samo odsyłanie sięgnij "+
				"postacią odpowiedzi \"reference\"", z.AssetId, len(bajty), *z.MaxBytes))
	}

	suma := sha256.Sum256(bajty)
	odpowiedz := shared.DesignAssetContentGetResponse{
		AssetId:   zasob.Kod,
		MediaType: typTresciZasobuDesignu(zasob.Format, bajty),
		SizeBytes: len(bajty),
		Checksum:  hex.EncodeToString(suma[:]),
	}
	// Postać odsyłania oddaje miarę i sumę bez bajtów — kontrakt czyni
	// `contentBase64` niewymaganym właśnie po to.
	if z.Disposition != nil && *z.Disposition == shared.AssetContentDispositionReference {
		return odpowiedz, nil
	}
	tresc := base64.StdEncoding.EncodeToString(bajty)
	odpowiedz.ContentBase64 = &tresc
	return odpowiedz, nil
}

// typTresciZasobuDesignu rozstrzyga typ treści wedle IANA. Pierwszeństwo ma
// format zapisany w wierszu (zmierzony przy wniesieniu z nagłówka pliku);
// dopiero gdy go nie ma, rozstrzyga sam wykrywacz biblioteki standardowej,
// który czyta początkowe bajty.
//
// Zgadywania po nazwie pliku tu nie ma: nazwa jest wolnym tekstem Operatora
// i mówi o pliku dokładnie tyle, ile Operator w nią wpisał.
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
	// DetectContentType dokleja parametr zestawu znaków do typów tekstowych;
	// kontrakt chce samego typu.
	if numer := strings.Index(wykryty, ";"); numer > 0 {
		return strings.TrimSpace(wykryty[:numer])
	}
	return wykryty
}
