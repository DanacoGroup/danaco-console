// Odpowiedzialność pliku: wpięcie dwóch komend rodziny `speech.*` — uczciwego
// sprawdzenia gotowości silnika mowy i zamiany jednego nagrania na tekst.
// Adapter wraz z rozstrzygnięciami — skąd bierze się okno, zasady i obszar oraz
// jak typowane odmowy pakietu `mowa` przekładają się na kody kontraktu — leży
// w `adapter_modul_mowa.go`; sam silnik w `server/internal/mowa`.
//
// Osiem komend. Poza sprawdzeniem gotowości i transkrypcją rodzina niesie
// przyjęcie i oddanie bajtów nagrania (`speech.audio.*`), nastawę wybudzania
// (`speech.wake.*`) oraz nasłuch ciągły (`speech.listen.*`). Dziewiąta pozycja
// o tym przedrostku, `speech.unknown`, jest odpowiedzią na komendę nieznaną
// obszaru, a nie komendą: nie ma pary żądanie/wynik i w rejestrze się nie zjawia.
//
// Wszystkie osiem stoi na jednym porcie, bo stoją na jednym adapterze i na tym
// samym silniku. Rozdzielenie ich na dwa porty rozdzieliłoby też stan, którego
// rozdzielić nie wolno: nasłuch ciągły rozpoznaje odcinki przysłane komendą
// `speech.audio.upload`, więc obie muszą widzieć ten sam rejestr nasłuchów.
//
// Synteza mowy tu nie należy. `translate.speech.synthesize` idzie w drugą stronę
// — tekst na dźwięk — i obsługuje ją port `Tlumaczenie`. Wspólny przedrostek
// „speech" w nazwie to zbieżność słowa, nie jednej maszynerii: rdzeń rozpoznaje
// mowę i nadal jej nie syntezuje.
//
// Zdarzenia rodzina MA, ale wyłącznie w nasłuchu ciągłym: `speech.listen.partial`
// niesie rozpoznany odcinek, a `speech.wake.detected` — wykrycie frazy
// wybudzającej. Rozgłasza je adapter przy przyjęciu odcinka, nadajnikiem
// wpiętym metodą `ZWyjsciem`, a nie obsługa komend: odcinek przychodzi komendą
// `speech.audio.upload`, której odpowiedź mówi o przyjęciu bajtów, a nie
// o tym, co w nich usłyszano. `speech.changed` kontrakt nie zna i rdzeń go nie
// wymyśla.
//
// Port niewypełniony nie rejestruje niczego: obie komendy odpowiedzą wtedy
// `speech.unknown`, a pozostałe domeny pracują bez zmian.
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

// zarejestrujMowe wpina osiem komend rodziny `speech.*`.
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
