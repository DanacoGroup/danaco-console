// Plik wpina osiem komend rodziny speech.* — sprawdzenie gotowości silnika
// mowy, transkrypcję, obsługę nagrania, nastawę wybudzania i nasłuch ciągły,
// na jednym porcie i jednym adapterze.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Mowa jest portem rodziny `speech.*`. Mówi wyłącznie typami
// kontraktu; pomocnik Pythona, jego wyjście i wiersze dziennika transkrypcji
// leżą po drugiej stronie adaptera.
type Mowa interface {
	// Gotowosc obsługuje `speech.availability.get`.
	Gotowosc(ctx context.Context, z shared.SpeechAvailabilityGetRequest) (shared.SpeechAvailabilityGetResponse, error)
	// Przepisz obsługuje `speech.transcribe`.
	Przepisz(ctx context.Context, z shared.SpeechTranscribeRequest) (shared.SpeechTranscribeResponse, error)
	// PrzyjmijNagranie obsługuje `speech.audio.upload`.
	PrzyjmijNagranie(ctx context.Context, z shared.SpeechAudioUploadRequest) (shared.SpeechAudioUploadResponse, error)
	// OddajNagranie obsługuje `speech.audio.fetch`.
	OddajNagranie(ctx context.Context, z shared.SpeechAudioFetchRequest) (shared.SpeechAudioFetchResponse, error)
	// NastawaWybudzania obsługuje `speech.wake.get`.
	NastawaWybudzania(ctx context.Context, z shared.SpeechWakeGetRequest) (shared.SpeechWakeGetResponse, error)
	// ZapiszNastaweWybudzania obsługuje `speech.wake.set`.
	ZapiszNastaweWybudzania(ctx context.Context, z shared.SpeechWakeSetRequest) (shared.SpeechWakeSetResponse, error)
	// UruchomNasluch obsługuje `speech.listen.start`.
	UruchomNasluch(ctx context.Context, z shared.SpeechListenStartRequest) (shared.SpeechListenStartResponse, error)
	// ZatrzymajNasluch obsługuje `speech.listen.stop`.
	ZatrzymajNasluch(ctx context.Context, z shared.SpeechListenStopRequest) (shared.SpeechListenStopResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ Mowa = (*adapterMowy)(nil)

// zarejestrujMowe wpina osiem komend rodziny speech.* na porcie Mowa,
// pomijając rejestrację, gdy port albo rejestr są puste.
func zarejestrujMowe(r *Rejestr, m Mowa) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandSpeechAvailabilityGet, obsluz(m.Gotowosc))
	r.Zarejestruj(shared.CommandSpeechTranscribe, obsluz(m.Przepisz))
	r.Zarejestruj(shared.CommandSpeechAudioUpload, obsluz(m.PrzyjmijNagranie))
	r.Zarejestruj(shared.CommandSpeechAudioFetch, obsluz(m.OddajNagranie))
	r.Zarejestruj(shared.CommandSpeechWakeGet, obsluz(m.NastawaWybudzania))
	r.Zarejestruj(shared.CommandSpeechWakeSet, obsluz(m.ZapiszNastaweWybudzania))
	r.Zarejestruj(shared.CommandSpeechListenStart, obsluz(m.UruchomNasluch))
	r.Zarejestruj(shared.CommandSpeechListenStop, obsluz(m.ZatrzymajNasluch))
}
