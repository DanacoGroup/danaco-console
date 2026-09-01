package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// nadajnikZTelemetria czyta szynę zdarzeń i zamienia fragmenty strumienia odpowiedzi na telemetrię postępu w czasie rzeczywistym. Zdarzenie telemetrii wychodzi nadajnikiem opakowanym, więc nie wraca tutaj powtórnie.
type nadajnikZTelemetria struct {
	nadajnik   Nadajnik
	telemetria *telemetriaPostepu
}

// OwinNadajnik zakłada telemetrię na nadajnik transportu. Brak telemetrii albo
// brak nadajnika zwraca nadajnik bez zmiany.
func (t *telemetriaPostepu) OwinNadajnik(nadajnik Nadajnik) Nadajnik {
	if t == nil || nadajnik == nil {
		return nadajnik
	}
	return nadajnikZTelemetria{nadajnik: nadajnik, telemetria: t}
}

// Rozglos przepuszcza komunikat i dopiero potem go odczytuje. Kolejność jest
// celowa: telemetria nie ma prawa opóźnić ani zablokować komunikatu właściwego.
func (n nadajnikZTelemetria) Rozglos(konto string, k protocol.Koperta) {
	n.nadajnik.Rozglos(konto, k)
	n.odnotuj(k)
}

// RozlaczPoUniewaznieniu przepuszcza zrywanie gniazd do nadajnika opakowanego;
// nadajnik bez tej zdolności oddaje zero.
func (n nadajnikZTelemetria) RozlaczPoUniewaznieniu(ctx context.Context, konto, powod string) int {
	rozlaczanie, umie := n.nadajnik.(RozlaczanieSesji)
	if !umie {
		return 0
	}
	return rozlaczanie.RozlaczPoUniewaznieniu(ctx, konto, powod)
}

// odnotuj zgłasza telemetrii punkt pracy wyczytany z komunikatu. Komunikat
// spoza strumienia i fragment nieczytelny przechodzą bez zgłoszenia.
func (n nadajnikZTelemetria) odnotuj(k protocol.Koperta) {
	if k.Type != shared.EventStreamChunk {
		return
	}
	fragment, err := protocol.FragmentZKoperty(k)
	if err != nil || fragment.WindowId == "" {
		return
	}
	opis := opisTury(fragment.WindowId, protocol.IdSesji(k))
	if protocol.Ostatni(k) {
		n.telemetria.Stan(opis, stanPoStrumieniu(fragment.Kind), etapKoniecTury)
		return
	}
	n.telemetria.Krok(opis, etapFragmentu(fragment.Kind))
}

// stanPoStrumieniu rozstrzyga stan procesu po ostatnim fragmencie tury: powodzenie albo niepowodzenie.
func stanPoStrumieniu(rodzaj shared.ChunkKind) shared.ProgressStatus {
	if rodzaj == shared.ChunkKindError {
		return shared.ProgressStatusFailed
	}
	return shared.ProgressStatusDone
}

// etapFragmentu nazywa etap odpowiadający rodzajowi fragmentu. Nazwa jedzie
// polem stepLabel kontraktu, więc Operator widzi, co dzieje się w turze.
func etapFragmentu(rodzaj shared.ChunkKind) string {
	switch rodzaj {
	case shared.ChunkKindToolUse:
		return etapNarzedzie
	case shared.ChunkKindToolResult:
		return etapWynikNarzedzia
	case shared.ChunkKindThinking:
		return etapRozumowanie
	default:
		return etapFragment
	}
}
