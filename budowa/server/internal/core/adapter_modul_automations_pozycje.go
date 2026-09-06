// Układanie zlecenia w kolejce: priorytet i skierowanie do innej kolejki (pola
// priority i targetQueueId komendy automation.queue.action) oraz odmowy
// Queue Managera.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ulozPozycje wykonuje dwie czynności Queue Managera wyprzedzające działanie
// silnika: zmianę priorytetu i skierowanie pozycji do innej kolejki. Obie
// wymagają wskazania pozycji — priorytet całej kolejki nie ma znaczenia,
// bo kolejność dotyczy zleceń.
func (a *adapterAutomatyk) ulozPozycje(ctx context.Context, z shared.AutomationQueueActionRequest) error {
	if z.Priority == nil && z.TargetQueueId == nil {
		return nil
	}
	idPozycji, err := numerPozycji(z.ItemId)
	if err != nil {
		return err
	}
	if z.Priority != nil {
		if err := a.repozytorium.UstawKolejnoscPozycji(ctx, idPozycji, *z.Priority); err != nil {
			return err
		}
	}
	if z.TargetQueueId != nil {
		docelowa, err := strconv.ParseInt(*z.TargetQueueId, 10, 64)
		if err != nil {
			return errNieznanaKolejkaDocelowa
		}
		if err := a.repozytorium.PrzeniesPozycje(ctx, idPozycji, docelowa); err != nil {
			return err
		}
	}
	return nil
}

// numerPozycji czyta identyfikator zlecenia. Brak wskazania przy zmianie
// priorytetu albo skierowaniu jest błędem żądania — rdzeń nie zgaduje, której
// pozycji dotyczy czynność.
func numerPozycji(idPozycji *string) (int64, error) {
	if wartoscTekstu(idPozycji) == "" {
		return 0, errPozycjaNiewskazana
	}
	numer, err := strconv.ParseInt(*idPozycji, 10, 64)
	if err != nil {
		return 0, errPozycjaNiewskazana
	}
	return numer, nil
}

// Odmowy Queue Managera nazwane raz, żeby ten sam powód nie brzmiał w rdzeniu
// na dwa sposoby. Każda niesie kod kontraktu, więc przechodzi przez
// `bladAutomatyki` bez zmiany znaczenia.
var (
	errBrakSilnikaKolejek = bladNiedostepnegoSilnika(
		"silnik kolejek nie jest wpięty — działania na kolejce automatyki nie ma kto wykonać")

	errPozycjaNiewskazana = bladWskazaniaAutomatyki(
		"zmiana priorytetu i skierowanie wymagają wskazania zlecenia (itemId)")

	errNieznanaKolejkaDocelowa = bladWskazaniaAutomatyki(
		"kolejka docelowa o nieznanym identyfikatorze")
)

// bladNiedostepnegoSilnika nazywa brak wykonawcy. Kontrakt nie ma kodu domena
// niewpięta, więc odmowa idzie kodem channel_unavailable, jedynym ponawialnym.
func bladNiedostepnegoSilnika(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Automations: "+powod))
}
