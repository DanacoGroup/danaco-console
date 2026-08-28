// Plik wpina dwie komendy rodziny media.* — pomiar materiału dźwiękowego
// albo filmowego i jego przetworzenie, jako narzędzia modelu wykonywane
// w trakcie tury.
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

// zarejestrujNarzedziaMediow wpina dwie komendy rodziny media.* na porcie
// NarzedziaMedia, pomijając rejestrację, gdy port jest pusty.
func zarejestrujNarzedziaMediow(r *Rejestr, n NarzedziaMedia) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandMediaInspect, obsluz(n.Zbadaj))
	r.Zarejestruj(shared.CommandMediaTranscode, obsluz(n.Przetworz))
}
