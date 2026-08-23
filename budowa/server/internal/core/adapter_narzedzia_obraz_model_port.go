// Port i wpięcie dwóch komend, które uruchamiają model nad obrazem —
// `image.upscale` i `image.background.remove`. Adapter leży
// w `adapter_narzedzia_obraz_model.go`, silniki i składanie ich wiersza
// poleceń w `adapter_narzedzia_obraz_model_silniki.go`. Obie są zarazem
// narzędziami modelu: `danaco_image_upscale` i
// `danaco_image_background_remove` w wykazie narzędzi.
//
// Port jest osobny od `NarzedziaObrazu` mimo wspólnego przedrostka `image.`.
// Tamten stoi wyłącznie na ImageMagicku i wywraca się na braku jednego
// binarium; ten stoi na dwóch silnikach neuronowych, z których każdy może być
// nieobecny osobno. Wpięcie tych metod do tamtego interfejsu związałoby
// dostępność sześciu komend w jedno „jest albo nie ma" — maszyna
// z ImageMagickiem i bez Real-ESRGAN-a straciłaby retusz razem
// z powiększaniem.
//
// Ta para nie ma zdarzeń, więc port nie bierze nadajnika: wytworzony zasób
// jest zasobem modułu Design i to jego rodzina zdarzeń o zasobach mówi.
//
// Port niewypełniony nie rejestruje niczego — obie komendy odpowiedzą wtedy
// `image.unknown`, a pozostałe domeny pracują bez zmian.
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
	// RozlozNaWarstwy obsługuje `image.layers.split` — rozkład na obiekty.
	// Należy tutaj, a nie do portu `NarzedziaObrazu`, bo stoi na tej samej
	// sieci segmentującej co wycięcie tła i znika razem z nią.
	RozlozNaWarstwy(ctx context.Context, z shared.ImageLayersSplitRequest) (shared.ImageLayersSplitResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie.
var _ NarzedziaObrazuModelu = (*adapterNarzedziObrazuModelu)(nil)

// zarejestrujNarzedziaObrazuModelu wpina trzy komendy stojące na sieciach.
func zarejestrujNarzedziaObrazuModelu(r *Rejestr, n NarzedziaObrazuModelu) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandImageUpscale, obsluz(n.Powieksz))
	r.Zarejestruj(shared.CommandImageBackgroundRemove, obsluz(n.UsunTlo))
	r.Zarejestruj(shared.CommandImageLayersSplit, obsluz(n.RozlozNaWarstwy))
}
