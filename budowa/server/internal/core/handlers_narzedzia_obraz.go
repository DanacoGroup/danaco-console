// Plik wpina sześć komend rodziny image.* — badanie, przekształcanie,
// poprawę, konwersję, złożenie i wektoryzację obrazu, jako narzędzia modelu
// na jednym porcie.
package core

import (
	"context"

	"danacoconsole/shared"
)

// NarzedziaObrazu jest portem rodziny `image.*`. Mówi wyłącznie
// typami kontraktu; ImageMagick, jego wiersz poleceń i blob wyniku leżą po
// drugiej stronie adaptera.
type NarzedziaObrazu interface {
	// Zbadaj obsługuje `image.inspect` — czynność czytająca.
	Zbadaj(ctx context.Context, z shared.ImageInspectRequest) (shared.ImageInspectResponse, error)
	// Przeksztalc obsługuje `image.transform` — geometria obrazu.
	Przeksztalc(ctx context.Context, z shared.ImageTransformRequest) (shared.ImageTransformResponse, error)
	// Popraw obsługuje `image.adjust` — retusz.
	Popraw(ctx context.Context, z shared.ImageAdjustRequest) (shared.ImageAdjustResponse, error)
	// Przekonwertuj obsługuje `image.convert` — format i kompresja.
	Przekonwertuj(ctx context.Context, z shared.ImageConvertRequest) (shared.ImageConvertResponse, error)
	// Zloz obsługuje `image.compose` — obraz na obrazie.
	Zloz(ctx context.Context, z shared.ImageComposeRequest) (shared.ImageComposeResponse, error)
	// Zwektoryzuj obsługuje `image.vectorize` — raster na ścieżki.
	Zwektoryzuj(ctx context.Context, z shared.ImageVectorizeRequest) (shared.ImageVectorizeResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie w produkcie.
var _ NarzedziaObrazu = (*adapterNarzedziObrazu)(nil)

// zarejestrujNarzedziaObrazu wpina sześć komend rodziny image.* na wspólnym
// porcie, bo dzielą rozwiązanie źródła, magazyn wyniku i wykaz zasobów okna.
func zarejestrujNarzedziaObrazu(r *Rejestr, n NarzedziaObrazu) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandImageInspect, obsluz(n.Zbadaj))
	r.Zarejestruj(shared.CommandImageTransform, obsluz(n.Przeksztalc))
	r.Zarejestruj(shared.CommandImageAdjust, obsluz(n.Popraw))
	r.Zarejestruj(shared.CommandImageConvert, obsluz(n.Przekonwertuj))
	r.Zarejestruj(shared.CommandImageCompose, obsluz(n.Zloz))
	r.Zarejestruj(shared.CommandImageVectorize, obsluz(n.Zwektoryzuj))
}
