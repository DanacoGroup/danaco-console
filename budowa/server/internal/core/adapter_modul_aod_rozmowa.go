// Cztery komendy nakładki AOD, które nie dotyczą obserwacji procesów:
// wiadomość, komplet kontekstu, podpowiedzi i polecenie głosowe. Wiadomość
// idzie portem rozmowy, polecenie głosowe idzie modułem Assistant, a
// podpowiedzi pochodzą z katalogu akcji.
package core

import (
	"context"
	"sort"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// ── aod.chat.send ───────────────────────────────────────────────────────────

// WyslijZNakladki obsługuje `aod.chat.send`. Oddaje obok odpowiedzi kontraktu
// całą przyjętą wiadomość — rozgłoszenie `message.changed` należy do uchwytu
// i bez treści wiadomości nie miałoby czego rozgłosić.
func (a *adapterNakladkiAod) WyslijZNakladki(ctx context.Context,
	z shared.AodChatSendRequest) (shared.AodChatSendResponse, shared.Message, error) {

	if z.Text == "" {
		return shared.AodChatSendResponse{}, shared.Message{},
			bladZadaniaNakladki("wiadomość z nakładki bez treści")
	}
	if a.rozmowa == nil {
		return shared.AodChatSendResponse{}, shared.Message{},
			bladBrakuSkladnikaNakladki("port rozmowy")
	}
	okno, err := a.oknoZadania(z.WindowId, z.SessionId)
	if err != nil {
		return shared.AodChatSendResponse{}, shared.Message{}, err
	}
	odpowiedz, err := a.rozmowa.Wyslij(ctx, shared.MessageSendRequest{
		WindowId: okno.Id,
		Content:  z.Text,
	})
	if err != nil {
		return shared.AodChatSendResponse{}, shared.Message{}, err
	}
	// Oddawane jest okno, które wiadomość naprawdę przyjęła.
	idOkna := odpowiedz.Message.WindowId
	if idOkna == "" {
		idOkna = okno.Id
	}
	return shared.AodChatSendResponse{
		MessageId: odpowiedz.Message.Id,
		WindowId:  idOkna,
	}, odpowiedz.Message, nil
}

// ── aod.context.get ─────────────────────────────────────────────────────────

// KontekstNakladki obsługuje `aod.context.get`: oddaje komplet kontekstu okna
// — ten sam, który przenosi `context.transfer` i który zasila Context
// Panel. Magazyn kompletu jest jeden na cały rdzeń; nakładka czyta z niego,
// a nie z drugiego, własnego.
func (a *adapterNakladkiAod) KontekstNakladki(_ context.Context,
	z shared.AodContextGetRequest) (shared.AodContextGetResponse, error) {

	if a.komplety == nil {
		return shared.AodContextGetResponse{}, bladBrakuSkladnikaNakladki(
			"magazyn kompletu kontekstu okien")
	}
	okno, err := a.oknoZadania(z.WindowId, z.SessionId)
	if err != nil {
		return shared.AodContextGetResponse{}, err
	}
	komplet, err := a.komplety.KompletOkna(okno.Id, wartoscLiczby(z.HistoryLimit))
	if err != nil {
		return shared.AodContextGetResponse{}, err
	}
	return shared.AodContextGetResponse{Context: komplet}, nil
}

// KompletOkna oddaje komplet kontekstu wskazanego okna wraz z historią rozmowy
// przyciętą do granicy żądania. Granica niedodatnia nie zawęża niczego.
func (a *adapterPrzenoszenia) KompletOkna(idOkna string,
	ograniczenieHistorii int) (shared.ContextBundle, error) {

	if a == nil || a.nadzorca == nil {
		return shared.ContextBundle{}, bladBrakuSkladnikaNakladki(
			"magazyn kompletu kontekstu okien")
	}
	okno, err := a.nadzorca.Rejestr().Okno(idOkna)
	if err != nil {
		return shared.ContextBundle{}, bladSesji(err)
	}
	komplet := uzupelnijZrodlem(a.magazyn.Odczytaj(okno.Id), okno,
		a.projekt(okno.IdSesji), a.historiaZGranica(okno.Id, ograniczenieHistorii))
	// Granica dotyczy kompletu, nie samego odczytu dziennika.
	komplet.HistoryMessageIds = ogonOdwolan(komplet.HistoryMessageIds, ograniczenieHistorii)
	return komplet, nil
}

// historiaZGranica odczytuje identyfikatory wiadomości okna, najwyżej tyle, ile
// dopuszcza granica; granica niedodatnia nie zawęża wykazu wiadomości.
func (a *adapterPrzenoszenia) historiaZGranica(idOkna string, granica int) []string {
	if a.historia == nil {
		return nil
	}
	var ograniczenie *int
	if granica > 0 {
		ograniczenie = &granica
	}
	wiadomosci, _ := a.historia.Wykaz(idOkna, nil, ograniczenie)
	return identyfikatoryWiadomosci(wiadomosci)
}

// ogonOdwolan zostawia najwyżej tyle ostatnich odwołań, ile dopuszcza granica.
// Historia rośnie w przód, więc obcięcie idzie od początku — nakładce potrzebne
// są wiadomości najświeższe, nie najstarsze.
func ogonOdwolan(odwolania []string, granica int) []string {
	if granica <= 0 || len(odwolania) <= granica {
		return odwolania
	}
	return odwolania[len(odwolania)-granica:]
}

// ── aod.suggestion ──────────────────────────────────────────────────────────

// Podpowiedzi obsługuje aod.suggestion: oddaje podpowiedzi następnego kroku
// w kolejności ważności, od bytu najwęższego (okno) do najszerszego (cała
// platforma), złożone z pozycji katalogu akcji.
func (a *adapterNakladkiAod) Podpowiedzi(ctx context.Context,
	z shared.AodSuggestionRequest) (shared.AodSuggestionResponse, error) {

	if a.akcje == nil {
		return shared.AodSuggestionResponse{}, bladBrakuSkladnikaNakladki("katalog akcji")
	}
	okno, jest, err := a.oknoPodpowiedzi(z)
	if err != nil {
		return shared.AodSuggestionResponse{}, err
	}
	pozycje := a.pozycjeZasiegow(ctx, okno, jest)
	chwila := time.Now().UTC().UnixMilli()
	granica := wartoscLiczby(z.Limit)
	podpowiedzi := make([]shared.AodSuggestion, 0, len(pozycje))
	for _, pozycja := range pozycje {
		if granica > 0 && len(podpowiedzi) == granica {
			break
		}
		podpowiedzi = append(podpowiedzi, podpowiedzZAkcji(pozycja, okno, jest, chwila))
	}
	return shared.AodSuggestionResponse{Suggestions: podpowiedzi}, nil
}

// oknoPodpowiedzi ustala okno, którego dotyczą podpowiedzi. Żądanie wskazujące
// byt nieistniejący jest odmawiane; żądanie niewskazujące niczego odmowy nie
// dostaje — zostają wtedy podpowiedzi całej platformy, bo są prawdziwe
// niezależnie od okna.
func (a *adapterNakladkiAod) oknoPodpowiedzi(
	z shared.AodSuggestionRequest) (session.Okno, bool, error) {

	if idProcesu := wartoscTekstu(z.ProcessId); idProcesu != "" {
		return a.oknoProcesu(idProcesu)
	}
	if wartoscTekstu(z.WindowId) == "" && wartoscTekstu(z.SessionId) == "" {
		okno, jest := a.oknoOstatniejCzynnosci()
		return okno, jest, nil
	}
	okno, err := a.oknoZadania(z.WindowId, z.SessionId)
	if err != nil {
		return session.Okno{}, false, err
	}
	return okno, true, nil
}

// oknoProcesu wskazuje okno procesu obserwowanego. Proces nieznany rejestrowi
// telemetrii jest bytem nieistniejącym — odpowiedź złożona z podpowiedzi
// globalnych byłaby wtedy ciszą udającą wynik.
func (a *adapterNakladkiAod) oknoProcesu(idProcesu string) (session.Okno, bool, error) {
	if a.telemetria == nil {
		return session.Okno{}, false, bladBrakuSkladnikaNakladki("rejestr telemetrii postępu")
	}
	for _, odpis := range a.telemetria.Odpisy() {
		if odpis.Id != idProcesu {
			continue
		}
		if odpis.IdOkna == "" || a.nadzorca == nil {
			// Proces bez okna istnieje; podpowiedzi zostają wtedy platformowe.
			return session.Okno{}, false, nil
		}
		okno, err := a.nadzorca.Rejestr().Okno(odpis.IdOkna)
		if err != nil {
			return session.Okno{}, false, nil
		}
		return okno, true, nil
	}
	return session.Okno{}, false, bladBrakuProcesuTelemetrii(idProcesu)
}

// zasiegPodpowiedzi wskazuje jeden poziom zasięgu wraz z bytem tego poziomu,
// z którego biorą się podpowiedzi.
type zasiegPodpowiedzi struct {
	poziom shared.ConfigScope
	klucz  string
}

// pozycjeZasiegow zbiera pozycje katalogu akcji od bytu najwęższego do
// najszerszego, bez powtórzeń. Kolejność zasięgów jest kolejnością ważności
// podpowiedzi: to, co dotyczy tego okna, stoi przed tym, co dotyczy platformy.
func (a *adapterNakladkiAod) pozycjeZasiegow(ctx context.Context, okno session.Okno,
	jest bool) []AkcjaKatalogu {

	zasiegi := []zasiegPodpowiedzi{}
	if jest {
		zasiegi = append(zasiegi,
			zasiegPodpowiedzi{shared.ConfigScopeWindow, okno.Id},
			zasiegPodpowiedzi{shared.ConfigScopeSession, okno.IdSesji},
			zasiegPodpowiedzi{shared.ConfigScopeModule, okno.Modul})
	}
	zasiegi = append(zasiegi, zasiegPodpowiedzi{shared.ConfigScopeGlobal, ""})

	widziane := map[string]struct{}{}
	zebrane := make([]AkcjaKatalogu, 0)
	for _, zasieg := range zasiegi {
		grupa := a.pozycjeZasiegu(ctx, zasieg.poziom, zasieg.klucz, jest)
		for _, pozycja := range grupa {
			if _, juz := widziane[pozycja.Id]; juz {
				continue
			}
			widziane[pozycja.Id] = struct{}{}
			zebrane = append(zebrane, pozycja)
		}
	}
	return zebrane
}

// pozycjeZasiegu odczytuje czynne pozycje jednego zasięgu i porządkuje je
// kolejnością katalogu, tą samą drogą przez port Akcje, którą czyta
// action.list.
func (a *adapterNakladkiAod) pozycjeZasiegu(ctx context.Context, poziom shared.ConfigScope,
	klucz string, oknoZnane bool) []AkcjaKatalogu {

	if poziom != shared.ConfigScopeGlobal && klucz == "" {
		return nil
	}
	czynne := true
	wynik, err := a.akcje.Wykaz(ctx, ZadanieKatalogAkcji{
		Scope: &poziom, ScopeId: &klucz, EnabledOnly: &czynne,
	})
	if err != nil {
		return nil
	}
	wybrane := make([]AkcjaKatalogu, 0, len(wynik.Actions))
	for _, pozycja := range wynik.Actions {
		if !warunekSpelniony(pozycja, oknoZnane) {
			continue
		}
		wybrane = append(wybrane, pozycja)
	}
	sort.SliceStable(wybrane, func(i, j int) bool {
		if wybrane[i].Order != wybrane[j].Order {
			return wybrane[i].Order < wybrane[j].Order
		}
		return wybrane[i].Id < wybrane[j].Id
	})
	return wybrane
}

// warunekSpelniony mówi, czy podpowiedź da się w ogóle wykonać z tego, co
// nakładka ma w ręku: rozstrzyga wyłącznie warunek okna i karty sesji, resztę
// wymaganych bytów uznaje za niespełnioną.
func warunekSpelniony(pozycja AkcjaKatalogu, oknoZnane bool) bool {
	warunek := wartoscTekstu(pozycja.Requires)
	switch warunek {
	case "":
		return true
	case "okno", "sesja":
		return oknoZnane
	default:
		return false
	}
}

// podpowiedzZAkcji składa jedną podpowiedź z pozycji katalogu akcji. Czas
// powstania jest chwilą złożenia odpowiedzi.
func podpowiedzZAkcji(pozycja AkcjaKatalogu, okno session.Okno, jest bool,
	chwila int64) shared.AodSuggestion {

	podpowiedz := shared.AodSuggestion{
		Id:        pozycja.Id,
		Text:      pozycja.Name,
		CreatedAt: chwila,
	}
	podpowiedz.CommandType = tekstOpcjonalny(pozycja.Command)
	if jest {
		podpowiedz.WindowId = tekstOpcjonalny(okno.Id)
		// Moduł idzie wprost z okna, nie okrężnie przez window.list.
		podpowiedz.ModuleId = tekstOpcjonalny(okno.Modul)
	}
	// Klasy zdarzenia podpowiedź z katalogu akcji nie niesie i nie ma nieść.
	return podpowiedz
}

// ── aod.voice.command ───────────────────────────────────────────────────────

// PolecenieGlosoweNakladki obsługuje aod.voice.command: kieruje polecenie
// wydane w nakładce do modułu Assistant i oddaje założone zlecenie. Żądanie
// bez transkrypcji i bez nagrania odmawia, bo nie ma czego wykonać.
func (a *adapterNakladkiAod) PolecenieGlosoweNakladki(ctx context.Context,
	z shared.AodVoiceCommandRequest) (shared.AodVoiceCommandResponse, error) {

	transkrypcja := wartoscTekstu(z.Transcript)
	nagranie := wartoscTekstu(z.AudioRef)
	if transkrypcja == "" && nagranie == "" {
		return shared.AodVoiceCommandResponse{}, bladZadaniaNakladki(
			"polecenie głosowe bez transkrypcji i bez odnośnika do nagrania")
	}
	if a.asystent == nil {
		return shared.AodVoiceCommandResponse{}, bladBrakuSkladnikaNakladki("moduł Assistant")
	}
	okno, err := a.oknoZadania(z.WindowId, z.SessionId)
	if err != nil {
		return shared.AodVoiceCommandResponse{}, err
	}
	odpowiedz, err := a.asystent.PolecenieGlosowe(ctx, shared.AssistantVoiceCommandRequest{
		WindowId:   okno.Id,
		AudioRef:   z.AudioRef,
		Transcript: z.Transcript,
		Speak:      z.Speak,
	})
	if err != nil {
		return shared.AodVoiceCommandResponse{}, err
	}
	return shared.AodVoiceCommandResponse{
		Transcript: odpowiedz.Transcript,
		Action:     odpowiedz.Action,
	}, nil
}

// Zgodność adaptera przenoszenia z kształtem, którego potrzebuje nakładka,
// sprawdzana jest przy kompilacji — wpięcie w montażu nie rozjedzie się po
// cichu.
var _ zrodloKompletuOkna = (*adapterPrzenoszenia)(nil)
