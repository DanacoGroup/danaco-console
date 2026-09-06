// Wstawienie zlecenia, priorytet i skierowanie kolejki automatyki.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Priorytet i skierowanie wymagają wskazania pozycji: kolejność dotyczy zleceń.
func (a *adapterAutomatyk) ulozPozycje(ctx context.Context, z shared.AutomationQueueActionRequest) error {
	if err := a.dopiszZleceniaAutomatyki(ctx, z); err != nil {
		return err
	}
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

// Zasilenie pomija kolejkę niepustą, więc bez tej drogi zlecenie nie trafiłoby
// do przebiegu już rozpoczętego.
func (a *adapterAutomatyk) dopiszZleceniaAutomatyki(ctx context.Context,
	z shared.AutomationQueueActionRequest) error {

	if z.Action != shared.QueueActionEnqueue || wartoscTekstu(z.WorkflowId) == "" {
		return nil
	}
	if a.kolejki == nil {
		return errBrakSilnikaKolejek
	}
	kolejkaID, err := strconv.ParseInt(z.QueueId, 10, 64)
	if err != nil {
		return bladWskazaniaAutomatyki("kolejka o nieznanym identyfikatorze: " + z.QueueId)
	}
	automatyka, err := a.automatykaZadania(ctx, z.WorkflowId)
	if err != nil || automatyka == nil {
		return err
	}
	kroki, err := a.zleceniaDoWstawienia(ctx, kolejkaID, *automatyka)
	if err != nil {
		return err
	}
	if len(kroki) == 0 {
		return errZleceniaJuzWKolejce
	}
	for _, pozycja := range kroki {
		if _, err := a.kolejki.repozytorium.DodajPozycje(ctx, pozycja); err != nil {
			return err
		}
	}
	return nil
}

// Rozpoznanie po tytule zlecenia: powtórzone wstawienie ma dołożyć brakujące
// kroki, a nie podwoić przebieg.
func (a *adapterAutomatyk) zleceniaDoWstawienia(ctx context.Context, kolejkaID int64,
	automatyka dane.Automatyka) ([]dane.Pozycja, error) {

	zastane, err := a.kolejki.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return nil, err
	}
	kroki, err := a.repozytorium.Kroki(ctx, automatyka.ID)
	if err != nil {
		return nil, err
	}
	zaleznosci, err := a.repozytorium.Zaleznosci(ctx, automatyka.ID)
	if err != nil {
		return nil, err
	}
	obecne := make(map[string]struct{}, len(zastane))
	for _, pozycja := range zastane {
		obecne[pozycja.Tytul] = struct{}{}
	}
	brakujace := make([]dane.Pozycja, 0, len(kroki))
	for _, krok := range ulozoneKroki(kroki, zaleznosci) {
		pozycja := pozycjaZKroku(kolejkaID, krok)
		if _, jest := obecne[pozycja.Tytul]; jest {
			continue
		}
		brakujace = append(brakujace, pozycja)
	}
	return brakujace, nil
}

// Brak wskazania jest błędem żądania — rdzeń nie zgaduje pozycji.
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

// Odmowy nazwane raz, żeby powód nie brzmiał w rdzeniu na dwa sposoby.
var (
	errBrakSilnikaKolejek = bladNiedostepnegoSilnika(
		"silnik kolejek nie jest wpięty — działania na kolejce automatyki nie ma kto wykonać")

	errPozycjaNiewskazana = bladWskazaniaAutomatyki(
		"zmiana priorytetu i skierowanie wymagają wskazania zlecenia (itemId)")

	errNieznanaKolejkaDocelowa = bladWskazaniaAutomatyki(
		"kolejka docelowa o nieznanym identyfikatorze")

	errZleceniaJuzWKolejce = bladWskazaniaAutomatyki(
		"wszystkie kroki wskazanej automatyki są już w kolejce — nie ma czego wstawić")
)

// Kod channel_unavailable, bo kontrakt nie ma kodu domena niewpięta.
func bladNiedostepnegoSilnika(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Automations: "+powod))
}
