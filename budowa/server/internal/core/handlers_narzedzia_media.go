// Odpowiedzialność pliku: wpięcie dwóch komend rodziny `media.*` — pomiaru
// materiału dźwiękowego albo filmowego i jego przetworzenia. Adapter wraz
// z rozstrzygnięciami leży w `adapter_narzedzia_media*.go`, a jedyna droga
// wołania binarium w `server/internal/zewnetrzne/wolanie.go`.
//
// Dwie komendy, ani jednej więcej. Kontrakt niesie w tej rodzinie
// `media.inspect` i `media.transcode` — obie wpisane jako narzędzia modelu
// (`danaco_media_inspect`, `danaco_media_transcode` w wykazie narzędzi
// kontraktu). Trzeciej nie ma i rdzeń jej nie wymyśli: nazwa spoza kontraktu
// byłaby nazwą, której nie zna nikt poza rdzeniem.
//
// Zdarzeń rodzina nie ma. Kontrakt nie zna `media.changed`, więc żadna z komend
// niczego nie rozgłasza i port nie bierze nadajnika — mimo że `media.transcode`
// zakłada zasób. Rozgłoszenie własnego zdarzenia dałoby klientowi kopertę,
// której nie zna jego strona kontraktu.
//
// Port niewypełniony nie rejestruje niczego: obie komendy odpowiedzą wtedy
// `media.unknown`, a pozostałe domeny pracują bez zmian.
package core

import (
	"context"

	"danacoconsole/shared"
)

// NarzedziaMedia jest portem rodziny `media.*`. Mówi wyłącznie
// typami kontraktu; `ffprobe`, `ffmpeg`, ich argumenty i magazyn bajtów leżą
// po drugiej stronie adaptera.
type NarzedziaMedia interface {
	// Zbadaj obsługuje `media.inspect`.
	Zbadaj(ctx context.Context, z shared.MediaInspectRequest) (shared.MediaInspectResponse, error)
	// Przetworz obsługuje `media.transcode`.
	Przetworz(ctx context.Context, z shared.MediaTranscodeRequest) (shared.MediaTranscodeResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ NarzedziaMedia = (*adapterNarzedziMediow)(nil)

// zarejestrujNarzedziaMediow wpina dwie komendy rodziny `media.*`.
func zarejestrujNarzedziaMediow(r *Rejestr, n NarzedziaMedia) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandMediaInspect, obsluz(n.Zbadaj))
	r.Zarejestruj(shared.CommandMediaTranscode, obsluz(n.Przetworz))
}
