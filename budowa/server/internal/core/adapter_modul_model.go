// Plik obsługuje `model.channel.set`: wybór kanału modelu obsługującego okno komunikacji albo kartę sesji. Rozszerza port okien, nie tworzy drugiego adaptera nad tym samym rejestrem.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// UstawKanalModelu obsługuje `model.channel.set`. Zwraca odpowiedź kontraktu oraz wykaz okien, których kanał został zmieniony, bo żądanie karty sesji dotyka wielu okien naraz.
func (a *adapterOkien) UstawKanalModelu(ctx context.Context,
	z shared.ModelChannelSetRequest) (shared.ModelChannelSetResponse, []shared.Window, error) {

	kanal, err := a.kanalRejestru(ctx, z.ModelChannelId)
	if err != nil {
		return shared.ModelChannelSetResponse{}, nil, err
	}
	idOkna := wartoscTekstu(z.WindowId)
	idSesji := wartoscTekstu(z.SessionId)
	switch {
	case idOkna != "":
		return a.kanalOkna(ctx, kanal, idOkna, z.AgentId)
	case idSesji != "":
		return a.kanalKartySesji(ctx, kanal, idSesji, z.AgentId)
	default:
		return shared.ModelChannelSetResponse{}, nil, bladKanaluModelu(
			"żądanie bez wskazania okna i bez wskazania karty sesji; " +
				"nie ma czemu nadać kanału")
	}
}

// kanalRejestru odnajduje wiersz rejestru kanałów po kodzie z żądania. Kanał jest wybierany z rejestru, nie zakładany; kodu, którego tam nie ma, nie da się nadać oknu.
func (a *adapterOkien) kanalRejestru(ctx context.Context, kod string) (dane.Kanal, error) {
	if kod == "" {
		return dane.Kanal{}, bladKanaluModelu("żądanie bez wskazania kanału modelu")
	}
	if a.kanaly == nil {
		return dane.Kanal{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"rodzina model.*: rejestr kanałów niewpięty, kanału nie da się rozpoznać"))
	}
	kanaly, err := a.kanaly.Lista(ctx, false)
	if err != nil {
		return dane.Kanal{}, err
	}
	for _, kanal := range kanaly {
		if kanal.Kod == kod {
			return kanal, nil
		}
	}
	return dane.Kanal{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"rodzina model.*: kanał modelu "+kod+" nie istnieje w rejestrze"))
}

// kanalOkna nadaje kanał jednemu oknu komunikacji wraz z ekspertem żądania. Ekspert jedzie razem z kanałem, bo kontrakt tak go podaje polem `agentId`.
func (a *adapterOkien) kanalOkna(ctx context.Context, kanal dane.Kanal,
	idOkna string, agent *string) (shared.ModelChannelSetResponse, []shared.Window, error) {

	okno, err := a.nadzorca.Rejestr().ZmienOkno(idOkna,
		session.Zmiana{KanalModelu: &kanal.Kod, Agent: agent})
	if err != nil {
		return shared.ModelChannelSetResponse{}, nil, bladSesji(err)
	}
	if err := a.utrwalKanalOkna(ctx, idOkna, kanal.ID); err != nil {
		return shared.ModelChannelSetResponse{}, nil, err
	}
	if err := a.utrwalAgentaOkna(ctx, idOkna, agent); err != nil {
		return shared.ModelChannelSetResponse{}, nil, err
	}
	zmienione := []shared.Window{oknoKontraktu(okno)}
	return shared.ModelChannelSetResponse{
		Channel: kanalOdpowiedzi(kanal), WindowId: &idOkna,
	}, zmienione, nil
}

// kanalKartySesji nadaje kanał karcie sesji, czyli wszystkim jej otwartym oknom. Kanał naprawdę obsługuje wyłącznie okno, więc nadanie karcie znaczy: nadaj go każdemu oknu, które ta karta prowadzi.
func (a *adapterOkien) kanalKartySesji(ctx context.Context, kanal dane.Kanal,
	idSesji string, agent *string) (shared.ModelChannelSetResponse, []shared.Window, error) {

	okna, err := a.nadzorca.Rejestr().OknaSesji(idSesji)
	if err != nil {
		return shared.ModelChannelSetResponse{}, nil, bladSesji(err)
	}
	zmienione := make([]shared.Window, 0, len(okna))
	for _, okno := range okna {
		if !okno.CzyOtwarte() {
			continue
		}
		zmienione, err = a.dopiszZmianeKanalu(ctx, zmienione, okno.Id, kanal, agent)
		if err != nil {
			return shared.ModelChannelSetResponse{}, nil, err
		}
	}
	// Karta bez otwartego okna nie ma czemu nadać kanału.
	if len(zmienione) == 0 {
		return shared.ModelChannelSetResponse{}, nil, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict,
			"rodzina model.*: karta sesji "+idSesji+" nie prowadzi ani jednego otwartego okna; "+
				"kanał modelu obsługuje okno, nie kartę"))
	}
	return shared.ModelChannelSetResponse{
		Channel: kanalOdpowiedzi(kanal), SessionId: &idSesji,
	}, zmienione, nil
}

// kanalOdpowiedzi składa kanał kontraktu na odpowiedź tej komendy, doklejając czas utworzenia, którego przekład wiersza rejestru nie nanosi.
func kanalOdpowiedzi(k dane.Kanal) shared.Channel {
	kanal := kanalKontraktu(k)
	if chwila, jest := chwilaZapisu(k.Utworzono); jest {
		kanal.CreatedAt = chwila
	}
	return kanal
}

// dopiszZmianeKanalu nanosi kanał i eksperta na jedno okno karty i dopisuje je do wykazu zmienionych. Ekspert obejmuje te same okna co kanał.
func (a *adapterOkien) dopiszZmianeKanalu(ctx context.Context, zmienione []shared.Window,
	idOkna string, kanal dane.Kanal, agent *string) ([]shared.Window, error) {

	okno, err := a.nadzorca.Rejestr().ZmienOkno(idOkna,
		session.Zmiana{KanalModelu: &kanal.Kod, Agent: agent})
	if err != nil {
		return nil, bladSesji(err)
	}
	if err := a.utrwalKanalOkna(ctx, idOkna, kanal.ID); err != nil {
		return nil, err
	}
	if err := a.utrwalAgentaOkna(ctx, idOkna, agent); err != nil {
		return nil, err
	}
	return append(zmienione, oknoKontraktu(okno)), nil
}

// utrwalKanalOkna zapisuje wybór w wierszu okna, żeby przeżył restart rdzenia. Parametry modelu są pamiętane per okno, a rejestr pamięciowy odtwarza się po restarcie z wierszy.
func (a *adapterOkien) utrwalKanalOkna(ctx context.Context, idOkna string, kanalID int64) error {
	if a.trwalosc == nil || a.trwalosc.okna == nil {
		return nil
	}
	wiersz, err := a.trwalosc.okna.PoIdentyfikatorze(ctx, idOkna)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return nil
	}
	if err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"rodzina model.*: nie da się odczytać wiersza okna "+idOkna+": "+err.Error()))
	}
	if wiersz.KanalModeluID == kanalID {
		return nil
	}
	wiersz.KanalModeluID = kanalID
	if err := a.trwalosc.okna.Aktualizuj(ctx, wiersz); err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"rodzina model.*: nie da się zapisać kanału okna "+idOkna+": "+err.Error()))
	}
	return nil
}

// bladKanaluModelu składa odmowę żądania niezgodnego z kontraktem, niosąc kod błędu i nazwę bytu, którego żądanie dotyczyło.
func bladKanaluModelu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"rodzina model.*: "+powod))
}
