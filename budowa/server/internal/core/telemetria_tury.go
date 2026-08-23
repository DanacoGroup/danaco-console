package core

import (
	"context"

	"danacoconsole/shared"
)

// rozmowaZTelemetria dokłada do portu rozmowy zgłoszenia telemetrii postępu
// z punktów, których nie widać na szynie zdarzeń: przyjęcia wiadomości,
// otwarcia tury i jej zatrzymania.
//
// Owinięcie, nie druga implementacja rozmowy: adapter rozmowy nadal
// robi całą pracę, a telemetria wyłącznie odnotowuje, co się z nim stało.
// Owinięcie jest też jedynym sposobem, żeby producent telemetrii nie musiał
// być wpleciony w każdy obsługiwacz z osobna.
type rozmowaZTelemetria struct {
	rozmowa    Rozmowa
	telemetria *telemetriaPostepu
}

// OwinRozmowe zakłada telemetrię na port rozmowy. Brak telemetrii albo brak
// portu zwraca port bez zmiany — owinięcie jest dodatkiem, nie warunkiem
// pracy rdzenia.
func (t *telemetriaPostepu) OwinRozmowe(rozmowa Rozmowa) Rozmowa {
	if t == nil || rozmowa == nil {
		return rozmowa
	}
	return rozmowaZTelemetria{rozmowa: rozmowa, telemetria: t}
}

// Wyslij otwiera turę okna i zgłasza jej start. Niepowodzenie przyjęcia
// wiadomości jest końcem procesu, bo tura się nie zaczęła.
func (r rozmowaZTelemetria) Wyslij(ctx context.Context, z shared.MessageSendRequest) (shared.MessageSendResponse, error) {
	odpowiedz, err := r.rozmowa.Wyslij(ctx, z)
	if err != nil {
		r.telemetria.Stan(opisTury(z.WindowId, ""), shared.ProgressStatusFailed, etapPrzyjecieTury)
		return odpowiedz, err
	}
	r.telemetria.Zacznij(opisTury(odpowiedz.Message.WindowId, odpowiedz.Message.SessionId), etapStartTury)
	return odpowiedz, nil
}

// Zatrzymaj przerywa turę okna. Zgłoszenie idzie wyłącznie wtedy, gdy tura
// naprawdę biegła — przycisk zatrzymania jest czynny zawsze, ale
// zatrzymanie procesu, którego nie ma, nie jest zdarzeniem postępu.
func (r rozmowaZTelemetria) Zatrzymaj(ctx context.Context, z shared.MessageStopRequest) (shared.MessageStopResponse, error) {
	odpowiedz, err := r.rozmowa.Zatrzymaj(ctx, z)
	if err == nil && odpowiedz.Stopped {
		r.telemetria.Stan(opisTury(z.WindowId, ""), shared.ProgressStatusStopped, etapZatrzymanie)
	}
	return odpowiedz, err
}

// Wykaz przechodzi bez zmiany: odczyt historii nie jest punktem pracy tury.
func (r rozmowaZTelemetria) Wykaz(ctx context.Context, z shared.MessageListRequest) (shared.MessageListResponse, error) {
	return r.rozmowa.Wykaz(ctx, z)
}

// opisTury buduje opis procesu tury. Kluczem jest okno, bo jedno okno prowadzi
// jedną turę naraz, a kolejne okna sesji biegną równolegle.
func opisTury(idOkna, idSesji string) opisProcesu {
	return opisProcesu{Klucz: idOkna, IdOkna: idOkna, IdSesji: idSesji}
}
