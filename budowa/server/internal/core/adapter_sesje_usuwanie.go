package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Usun wyprowadza wskazane sesje z historii do kosza — jedyną drogą utraty danych sesji
// w produkcie. Zdarzenie jest odwracalne w oknie terminu komendą session.restore.
func (a *adapterSesji) Usun(ctx context.Context, z shared.SessionDeleteRequest) (shared.SessionDeleteResponse, error) {
	if !z.Confirm {
		return shared.SessionDeleteResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeValidationFailed,
			"usunięcie trwałe wymaga potwierdzenia — pole confirm nie zostało ustawione"))
	}
	wynik := shared.SessionDeleteResponse{DeletedIds: []string{}, MissingIds: []string{}}
	for _, identyfikator := range z.SessionIds {
		a.zwolnijZasobySesji(identyfikator)
		if err := a.usunZapisSesji(ctx, identyfikator); err != nil {
			if errors.Is(err, dane.ErrBrakWiersza) {
				wynik.MissingIds = append(wynik.MissingIds, identyfikator)
				continue
			}
			return shared.SessionDeleteResponse{}, err
		}
		wynik.DeletedIds = append(wynik.DeletedIds, identyfikator)
	}
	wynik.DeletedCount = len(wynik.DeletedIds)
	return wynik, nil
}

// zwolnijZasobySesji zatrzymuje procesy okien i zdejmuje sesję z rejestru żywego przed
// skasowaniem zapisu, aby nie zostawić biegnącego procesu bez sesji. Niepowodzenie nie
// wstrzymuje usuwania z historii.
func (a *adapterSesji) zwolnijZasobySesji(identyfikator string) {
	okna, err := a.nadzorca.Rejestr().UsunSesje(identyfikator)
	if err != nil {
		return
	}
	_ = a.nadzorca.Procesy().ZatrzymajOkna(identyfikatoryOkienSesji(okna))
}

// usunZapisSesji przenosi wiersz sesji do kosza — zapis zostaje w całości, znika wyłącznie
// z wykazów. Brak utrwalacza znaczy rdzeń bez bazy, co nie jest błędem.
func (a *adapterSesji) usunZapisSesji(ctx context.Context, identyfikator string) error {
	if a.trwalosc == nil {
		return nil
	}
	return a.trwalosc.PrzeniesSesjeDoKosza(ctx, identyfikator)
}

// ZapewnieniemSesji nazywa czynność utrwalenia sesji zaraz po jej założeniu.
// Port jest wąski celowo — adapter sesji nie potrzebuje całego utrwalacza
// rozmowy, tylko tej jednej gwarancji.
type ZapewnienieSesji interface {
	ZapewnijSesje(kontekst context.Context, idSesji, tytul, projekt string) (int64, error)
	// ZapewnijSesjeSrodowiska odkłada sesję w karcie środowiska wejścia Operatora.
	ZapewnijSesjeSrodowiska(kontekst context.Context, idSesji, tytul, projekt, kodSrodowiska string) (int64, error)
}

// ZZapewnieniem wpina utrwalanie sesji przy zakładaniu. Bez tego portu sesja materializuje się
// w bazie dopiero przy pierwszej wiadomości, a nieużyta ginie po restarcie.
func (a *adapterSesji) ZZapewnieniem(z ZapewnienieSesji) *adapterSesji {
	a.zapewnienie = z
	return a
}

// utrwalZalozona zapisuje świeżo założoną sesję. Niepowodzenie nie przerywa
// zakładania: sesja żyje w rejestrze i pracuje, tylko nie przetrwa restartu —
// odmowa założenia byłaby dla Operatora gorsza.
func (a *adapterSesji) utrwalZalozona(kontekst context.Context, sesja session.Sesja) {
	if a.zapewnienie == nil {
		return
	}
	if _, err := a.zapewnienie.ZapewnijSesjeSrodowiska(kontekst, sesja.Id, sesja.Tytul,
		sesja.IdProjektu, sesja.KodSrodowiska); err != nil {
		_ = err
	}
}
