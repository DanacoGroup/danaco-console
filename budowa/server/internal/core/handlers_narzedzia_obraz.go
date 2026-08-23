// Odpowiedzialność pliku: port i wpięcie czterech komend rodziny `image.*` —
// `image.inspect`, `image.transform`, `image.adjust`, `image.convert`.
// Adapter wraz z rozstrzygnięciami (skąd bierze się źródło, jak woła się
// binarium, gdzie ląduje wynik) leży w `adapter_narzedzia_obraz.go`, składanie
// argumentów każdej operacji w `adapter_narzedzia_obraz_czynnosci.go`.
//
// Tyle właśnie komend niesie kontrakt w tej rodzinie i wszystkie cztery są
// zarazem narzędziami modelu (`danaco_image_*`). Piąta pozycja o tym
// przedrostku, `image.unknown`, jest odpowiedzią na komendę nieznaną obszaru,
// a nie komendą: nie ma pary żądanie/wynik i w rejestrze się nie zjawia.
//
// Zdarzeń rodzina nie ma, więc port nie bierze nadajnika. Wytworzony zasób jest
// zasobem modułu Design i mówi o nim jego rodzina zdarzeń; `image.changed`
// byłoby nazwą, której klient nie zna.
//
// Port niewypełniony nie rejestruje niczego: cztery komendy odpowiedzą wtedy
// `image.unknown`, a pozostałe domeny pracują bez zmian — rdzeń niczym nie
// warunkuje startu.
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

// zarejestrujNarzedziaObrazu wpina sześć komend rodziny `image.*` stojących na
// tym porcie.
//
// Dwie ostatnie — złożenie i wektoryzacja — nie wołają ani jednego programu:
// pracują bibliotekami wkompilowanymi w binarium rdzenia. Stoją mimo to na tym
// samym porcie, bo dzielą z resztą rodziny wszystko poza sposobem liczenia:
// rozwiązanie źródła (zasób albo ścieżka), magazyn wyniku i wykaz zasobów okna.
// Osobny port dałby drugą prawdę o tym, gdzie rdzeń odkłada bajty obrazu.
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
