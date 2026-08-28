// Plik wpina port dwóch komend uruchamiających model nad obrazem, image.upscale i image.background.remove, osobno od portu ImageMagicka.
package core

import (
	"context"

	"danacoconsole/shared"
)

// NarzedziaObrazuModelu jest portem pary komend stojących na sieciach
// neuronowych. Mówi wyłącznie typami kontraktu; Real-ESRGAN, rembg,
// ich wagi i pliki pośrednie leżą po drugiej stronie adaptera.
type NarzedziaObrazuModelu interface {
	// Powieksz obsługuje `image.upscale` — superrozdzielczość.
	Powieksz(ctx context.Context, z shared.ImageUpscaleRequest) (shared.ImageUpscaleResponse, error)
	// UsunTlo obsługuje `image.background.remove` — wycięcie obiektu z tła.
	UsunTlo(ctx context.Context, z shared.ImageBackgroundRemoveRequest) (shared.ImageBackgroundRemoveResponse, error)
	// RozlozNaWarstwy obsługuje image.layers.split: stoi na sieci segmentującej, wspólnej z wycięciem tła.
	RozlozNaWarstwy(ctx context.Context, z shared.ImageLayersSplitRequest) (shared.ImageLayersSplitResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie.
var _ NarzedziaObrazuModelu = (*adapterNarzedziObrazuModelu)(nil)

// zarejestrujNarzedziaObrazuModelu wpina trzy komendy stojące na sieciach neuronowych do rejestru komend rdzenia.
func zarejestrujNarzedziaObrazuModelu(r *Rejestr, n NarzedziaObrazuModelu) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandImageUpscale, obsluz(n.Powieksz))
	r.Zarejestruj(shared.CommandImageBackgroundRemove, obsluz(n.UsunTlo))
	r.Zarejestruj(shared.CommandImageLayersSplit, obsluz(n.RozlozNaWarstwy))
}
