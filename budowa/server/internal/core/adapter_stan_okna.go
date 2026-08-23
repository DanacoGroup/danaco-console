package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// StanOkna zwraca stan okna komunikacji: parametry wykonania, stan procesu,
// miarę historii i — na żądanie — konfigurację efektywną.
//
// Odpowiedź składa się z trzech warstw, bo z trzech warstw składa się samo okno:
// parametry wykonania zna rejestr nadzorcy (a po restarcie rdzenia wiersz),
// stan procesu zna rejestr procesów, historię zna baza. Okno nieznane wszystkim
// trzem daje odpowiedź pustą ze stanem `pending` — pytanie o stan nie ma prawa
// zerwać niczego.
func (a *adapterNawigacji) StanOkna(ctx context.Context, z shared.WindowStateGetRequest) (shared.WindowStateGetResponse, error) {
	wynik := shared.WindowStateGetResponse{
		ProcessStatus: a.stanProcesu(z.WindowId),
		Streaming:     a.strumien != nil && a.strumien(z.WindowId),
	}
	okno, err := a.oknoStanu(ctx, z.WindowId)
	if err != nil {
		return shared.WindowStateGetResponse{}, err
	}
	wynik.Window = okno

	liczba, ostatnia, err := a.historiaOkna(ctx, z.WindowId)
	if err != nil {
		return shared.WindowStateGetResponse{}, err
	}
	wynik.MessageCount, wynik.LastMessageId = liczba, ostatnia

	if z.IncludeConfig != nil && *z.IncludeConfig && a.ustawienia != nil {
		idOkna := z.WindowId
		// Poziom zasięgu najwęższy: bez wskazania poziomu wraca polityka
		// efektywna okna, czyli wartość obowiązująca wraz z jej źródłem.
		wpisy, err := a.ustawienia.Odczytaj(ctx, shared.ConfigGetRequest{ScopeId: &idOkna})
		if err != nil {
			return shared.WindowStateGetResponse{}, err
		}
		wynik.Config = wpisy.Entries
	}
	// Bieg naprawczy wychodzi wyłącznie dla okna koordynatora; okno samodzielne
	// i wykonawcze nie prowadzą pętli.
	if bieg, jest := a.biegi.Stan(z.WindowId); jest {
		wynik.Loop = &bieg
	}
	return wynik, nil
}

// oknoStanu zwraca okno żywe z rejestru nadzorcy, a gdy tego nie ma —
// utrwalone wierszem.
func (a *adapterNawigacji) oknoStanu(ctx context.Context, idOkna string) (shared.Window, error) {
	if a.nadzorca != nil {
		if okno, err := a.nadzorca.Rejestr().Okno(idOkna); err == nil {
			return oknoKontraktu(okno), nil
		}
	}
	okno, _, err := oknoUtrwalone(ctx, a.zestaw, idOkna)
	return okno, err
}

// stanProcesu przekłada obecność procesu okna na stan telemetrii postępu,
// odpowiednik kolumny `proces_sesji.stan`. Okno bez procesu oczekuje,
// proces żywy pracuje, proces zamknięty jest zatrzymany.
func (a *adapterNawigacji) stanProcesu(idOkna string) shared.ProgressStatus {
	if a.nadzorca == nil {
		return shared.ProgressStatusPending
	}
	proces, jest := a.nadzorca.Procesy().Proces(idOkna)
	switch {
	case !jest:
		return shared.ProgressStatusPending
	case proces.Zyje():
		return shared.ProgressStatusRunning
	default:
		return shared.ProgressStatusStopped
	}
}

// historiaOkna mierzy zapisaną rozmowę okna: liczbę wiadomości i identyfikator
// ostatniej. Okno bez wiersza albo bez historii daje zero i brak wskazania —
// nowe okno nie jest oknem błędnym.
func (a *adapterNawigacji) historiaOkna(ctx context.Context, idOkna string) (int, *string, error) {
	wiersz, err := a.zestaw.Okna.PoIdentyfikatorze(ctx, idOkna)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return 0, nil, nil
	}
	if err != nil {
		return 0, nil, err
	}
	wiadomosci, err := a.zestaw.Wiadomosci.ListaOkna(ctx, wiersz.ID, 0)
	if err != nil {
		return 0, nil, err
	}
	if len(wiadomosci) == 0 {
		return 0, nil, nil
	}
	ostatnia := wiadomosci[len(wiadomosci)-1]
	identyfikator := identyfikatorWiersza(ostatnia.IdentyfikatorZewnetrzny, ostatnia.ID)
	return len(wiadomosci), &identyfikator, nil
}
