package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// StanOkna zwraca stan okna komunikacji: parametry wykonania, stan procesu, miarę historii
// i — na żądanie — konfigurację efektywną. Okno nieznane rejestrom daje odpowiedź pustą
// ze stanem pending.
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
		// Poziom zasięgu najwęższy: brak wskazania zwraca politykę efektywną okna z jej źródłem.
		wpisy, err := a.ustawienia.Odczytaj(ctx, shared.ConfigGetRequest{ScopeId: &idOkna})
		if err != nil {
			return shared.WindowStateGetResponse{}, err
		}
		wynik.Config = wpisy.Entries
	}
	// Bieg naprawczy wychodzi tylko dla okna koordynatora, nie dla samodzielnego i wykonawczego.
	if bieg, jest := a.biegi.Stan(z.WindowId); jest {
		wynik.Loop = &bieg
	}
	return wynik, nil
}

// oknoStanu zwraca okno żywe z rejestru nadzorcy, a gdy tego okna tam nie ma, zwraca okno
// utrwalone wierszem bazy.
func (a *adapterNawigacji) oknoStanu(ctx context.Context, idOkna string) (shared.Window, error) {
	if a.nadzorca != nil {
		if okno, err := a.nadzorca.Rejestr().Okno(idOkna); err == nil {
			return oknoKontraktu(okno), nil
		}
	}
	okno, _, err := oknoUtrwalone(ctx, a.zestaw, idOkna)
	return okno, err
}

// stanProcesu oddaje stan okna tą samą wykładnią, którą niesie zdarzenie
// window.state.changed; rdzeń bez rejestru obecności zna wyłącznie procesy okna.
func (a *adapterNawigacji) stanProcesu(idOkna string) shared.ProgressStatus {
	var czynnosc *pamiecCzynnosci
	if a.obecnosc != nil {
		czynnosc = a.obecnosc.czynnosc
	}
	return stanProcesuOkna(a.nadzorca, czynnosc, idOkna)
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
