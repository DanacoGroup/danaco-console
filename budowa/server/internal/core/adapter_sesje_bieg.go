package core

import (
	"context"

	"danacoconsole/shared"
)

// ZZatrzymaniemTur wpina przerywanie tur okien. Bez niego zatrzymanie sesji
// odpowiada pustym wykazem — nie ma czego zatrzymać.
func (a *adapterSesji) ZZatrzymaniemTur(przerwij PrzerwanieTury) *adapterSesji {
	a.przerwijTure = przerwij
	return a
}

// Wznow otwiera zamkniętą sesję wraz z jej oknami.
//
// Odwrotność Zamknij: sesja wraca do pracy bieżącej. Stan trafia też do bazy,
// żeby wznowienie przeżyło restart rdzenia — inaczej sesja wróciłaby po starcie
// jako zakończona, choć Operator ją wznowił.
func (a *adapterSesji) Wznow(_ context.Context, z shared.SessionResumeRequest) (shared.SessionResumeResponse, error) {
	sesja, otwarte, err := a.nadzorca.Rejestr().WznowSesje(z.SessionId)
	if err != nil {
		return shared.SessionResumeResponse{}, bladSesji(err)
	}
	a.trwalosc.StanOkien(identyfikatoryOkienSesji(otwarte), shared.WindowStatusOpen)
	a.trwalosc.StanSesji(sesja.Id, sesja.Stan)
	return shared.SessionResumeResponse{
		Session: sesjaKontraktu(sesja),
		Windows: oknaKontraktu(otwarte),
	}, nil
}

// Zatrzymaj przerywa tury biegnące we wszystkich oknach sesji.
//
// Zatrzymanie sesji nie zamyka jej ani okien: Operator chce wstrzymać pracę
// modelu, a nie stracić miejsce, w którym pracuje. Okna zostają otwarte, zapis
// zostaje w całości, wznowienie pracy jest kolejną wiadomością.
//
// Wykaz zwraca okna, w których faktycznie coś przerwano — okno bez tury w biegu
// nie jest błędem i po prostu nie trafia do wykazu.
func (a *adapterSesji) Zatrzymaj(_ context.Context, z shared.SessionStopRequest) (shared.SessionStopResponse, error) {
	wynik := shared.SessionStopResponse{SessionId: z.SessionId, StoppedWindowIds: []string{}}
	if a.przerwijTure == nil {
		return wynik, nil
	}
	okna, err := a.nadzorca.Rejestr().OknaSesji(z.SessionId)
	if err != nil {
		return shared.SessionStopResponse{}, bladSesji(err)
	}
	for _, okno := range okna {
		if a.przerwijTure(okno.Id) {
			wynik.StoppedWindowIds = append(wynik.StoppedWindowIds, okno.Id)
		}
	}
	return wynik, nil
}
