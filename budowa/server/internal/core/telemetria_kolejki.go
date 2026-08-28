package core

import (
	"context"
	"strconv"

	"danacoconsole/shared"
)

// Telemetria postępu kolejki: etapem jest pozycja, biegiem naprawczym powtórzenie pozycji bez limitu.

// przedrostekProcesuKolejki rozdziela klucze procesów kolejek od kluczy
// procesów okien w jednym rejestrze telemetrii.
const przedrostekProcesuKolejki = "kolejka-"

// Słownik stanów pozycji stoi w kolejka_stany.go razem z tabelą przejść silnika wykonania.

// odnotujPostep zgłasza telemetrii stan kolejki po czynności. Brak telemetrii
// nie jest błędem — adapter pracuje wtedy bez rozgłaszania postępu.
func (a *adapterKolejek) odnotujPostep(ctx context.Context, id int64, kolejka shared.Queue, nazwaEtapu string) {
	if a.telemetria == nil {
		return
	}
	etap, etapow := a.etapyKolejki(ctx, id)
	opis := opisProcesu{
		Klucz:   przedrostekProcesuKolejki + strconv.FormatInt(id, 10),
		IdOkna:  pierwszeOkno(kolejka.WindowIds),
		IdSesji: kolejka.SessionId,
		Etap:    etap,
		Etapow:  etapow,
	}
	a.telemetria.Stan(opis, stanPostepuKolejki(kolejka.Status), nazwaEtapu)
}

// etapyKolejki liczy etap bieżący i liczbę etapów z pozycji kolejki. Kolejka
// bez pozycji i błąd odczytu dają liczbę etapów nieznaną — telemetria nie
// zmyśla wtedy stopnia ukończenia.
func (a *adapterKolejek) etapyKolejki(ctx context.Context, id int64) (int, int) {
	pozycje, err := a.repozytorium.ListaPozycji(ctx, id)
	if err != nil || len(pozycje) == 0 {
		return 0, 0
	}
	zakonczone := 0
	for _, pozycja := range pozycje {
		if czyStanKoncowyPozycji(pozycja.Stan) {
			zakonczone++
		}
	}
	if zakonczone >= len(pozycje) {
		return len(pozycje), len(pozycje)
	}
	return zakonczone + 1, len(pozycje)
}

// stanPostepuKolejki przekłada stan kolejki na stan procesu telemetrii. Oba
// słowniki należą do kontraktu — rdzeń tylko je wiąże.
func stanPostepuKolejki(stan shared.QueueStatus) shared.ProgressStatus {
	switch stan {
	case shared.QueueStatusRunning:
		return shared.ProgressStatusRunning
	case shared.QueueStatusPaused:
		return shared.ProgressStatusPaused
	case shared.QueueStatusStopped:
		return shared.ProgressStatusStopped
	case shared.QueueStatusDone:
		return shared.ProgressStatusDone
	default:
		return shared.ProgressStatusPending
	}
}

// pierwszeOkno wskazuje okno, do którego telemetria przypina proces kolejki.
// Kolejka bez okien wychodzi zdarzeniem bez okna — pole jest opcjonalne.
func pierwszeOkno(okna []string) string {
	if len(okna) == 0 {
		return ""
	}
	return okna[0]
}
