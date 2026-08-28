// Plik rejestruje siedem komend rodziny `aod.*`, obsługujących nakładkę stałej
// obecności asystenta przez jeden wspólny port telemetrii, sesji, rozmowy,
// modułu Assistant i katalogu akcji.
package core

import (
	"context"

	"danacoconsole/shared"
)

// NakladkaAod jest portem rodziny `aod.*`. Mówi wyłącznie typami
// kontraktu; procesy telemetrii, okna nadzorcy i wiersze katalogu akcji leżą po
// drugiej stronie adaptera.
type NakladkaAod interface {
	// StanNakladki obsługuje `aod.status.get`.
	StanNakladki(ctx context.Context, z shared.AodStatusGetRequest) (shared.AodStatusGetResponse, error)
	// KontekstNakladki obsługuje `aod.context.get`.
	KontekstNakladki(ctx context.Context, z shared.AodContextGetRequest) (shared.AodContextGetResponse, error)
	// Podpowiedzi obsługuje `aod.suggestion`.
	Podpowiedzi(ctx context.Context, z shared.AodSuggestionRequest) (shared.AodSuggestionResponse, error)
	// PrzypnijObserwacje obsługuje `aod.observe.attach`.
	PrzypnijObserwacje(ctx context.Context, z shared.AodObserveAttachRequest) (shared.AodObserveAttachResponse, error)
	// OdepnijObserwacje obsługuje `aod.observe.detach`.
	OdepnijObserwacje(ctx context.Context, z shared.AodObserveDetachRequest) (shared.AodObserveDetachResponse, error)
	// PolecenieGlosoweNakladki obsługuje `aod.voice.command`.
	PolecenieGlosoweNakladki(ctx context.Context, z shared.AodVoiceCommandRequest) (shared.AodVoiceCommandResponse, error)

	// WyslijZNakladki obsługuje `aod.chat.send` i oddaje przyjętą wiadomość obok
	// odpowiedzi kontraktu.
	WyslijZNakladki(ctx context.Context,
		z shared.AodChatSendRequest) (shared.AodChatSendResponse, shared.Message, error)

	// WyciszeniaNakladki obsługuje `aod.mute.get`.
	WyciszeniaNakladki(ctx context.Context,
		z shared.AodMuteGetRequest) (shared.AodMuteGetResponse, error)
	// PrzestawWyciszenieNakladki obsługuje `aod.mute.set` — jedna komenda ustawia
	// i znosi wyciszenie.
	PrzestawWyciszenieNakladki(ctx context.Context,
		z shared.AodMuteSetRequest) (shared.AodMuteSetResponse, shared.AodMute, error)
	// ZglosSygnalNakladki obsługuje `aod.signal.report`.
	ZglosSygnalNakladki(ctx context.Context,
		z shared.AodSignalReportRequest) (shared.AodSignalReportResponse, error)
	// SygnalyNakladki obsługuje `aod.signal.list`.
	SygnalyNakladki(ctx context.Context,
		z shared.AodSignalListRequest) (shared.AodSignalListResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ NakladkaAod = (*adapterNakladkiAod)(nil)

// zarejestrujNakladkeAod wpina siedem komend rodziny `aod.*` w router, wiążąc
// każdą z metodą portu NakladkaAod.
func zarejestrujNakladkeAod(r *Rejestr, n NakladkaAod, e *emiter) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandAodStatusGet, obsluz(n.StanNakladki))
	r.Zarejestruj(shared.CommandAodContextGet, obsluz(n.KontekstNakladki))
	r.Zarejestruj(shared.CommandAodSuggestion, obsluz(n.Podpowiedzi))
	r.Zarejestruj(shared.CommandAodObserveAttach, obsluz(n.PrzypnijObserwacje))
	r.Zarejestruj(shared.CommandAodObserveDetach, obsluz(n.OdepnijObserwacje))
	r.Zarejestruj(shared.CommandAodVoiceCommand, obsluz(n.PolecenieGlosoweNakladki))

	r.Zarejestruj(shared.CommandAodChatSend,
		obsluz(func(ctx context.Context, z shared.AodChatSendRequest) (shared.AodChatSendResponse, error) {
			odpowiedz, wiadomosc, err := n.WyslijZNakladki(ctx, z)
			if err == nil {
				e.wiadomosc(ctx, shared.ChangeKindCreated, wiadomosc)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAodMuteGet, obsluz(n.WyciszeniaNakladki))
	r.Zarejestruj(shared.CommandAodSignalReport, obsluz(n.ZglosSygnalNakladki))
	r.Zarejestruj(shared.CommandAodSignalList, obsluz(n.SygnalyNakladki))

	r.Zarejestruj(shared.CommandAodMuteSet,
		obsluz(func(ctx context.Context, z shared.AodMuteSetRequest) (shared.AodMuteSetResponse, error) {
			odpowiedz, wyciszenie, err := n.PrzestawWyciszenieNakladki(ctx, z)
			// Rozgłoszenie odróżnia zdarzenie realne od powtórzonego bez zmiany,
			// zapobiegając rozjazdowi powłok.
			if err == nil && odpowiedz.Changed {
				e.wyciszenieNakladki(rodzajZmianyWyciszenia(z.Muted), wyciszenie, odpowiedz.Mutes)
			}
			return odpowiedz, err
		}))
}

// rodzajZmianyWyciszenia rozstrzyga, czy wyciszenie powstało, czy zniknęło.
// Rozstrzyga to samo żądanie: `muted` równe prawdzie wycisza, równe fałszowi
// wyciszenie znosi.
func rodzajZmianyWyciszenia(wyciszone bool) shared.ChangeKind {
	if wyciszone {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindDeleted
}

// wyciszenieNakladki rozgłasza `aod.mute.changed`. Zdarzenie idzie bez wskazania
// sesji, bo wyciszenie obowiązuje wszystkie powłoki Operatora, a nie kartę,
// z której przyszło — po to jest bytem rdzenia.
func (e *emiter) wyciszenieNakladki(zmiana shared.ChangeKind, wyciszenie shared.AodMute,
	wykaz []shared.AodMute) {

	e.wyslij(shared.EventAodMuteChanged, "", shared.AodMuteChangedEvent{
		Change: zmiana,
		Mute:   wyciszenie,
		Mutes:  wykaz,
	})
}
