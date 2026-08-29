// Plik dokłada do adapterPrzekazaniaOkna metodę WykonajAkcje obsługującą
// komendę window.action, która sprawdza akcję w katalogu, zapisuje ślad
// zgłoszenia w dzienniku i odmawia wykonania kodem kontraktu, nie zwracając
// sukcesu w żadnej ścieżce.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// WykonajAkcje przyjmuje zgłoszenie akcji panelu okna, odnotowuje je
// w dzienniku i odmawia wykonania, ponieważ rdzeń nie ma wykonawcy akcji
// katalogu; sukcesu nie zwraca w żadnej ścieżce.
func (a *adapterPrzekazaniaOkna) WykonajAkcje(ctx context.Context,
	z shared.WindowActionRequest) (shared.WindowActionResponse, error) {

	if z.WindowId == "" || z.ActionId == "" {
		return shared.WindowActionResponse{}, bladWskazaniaAkcjiOkna(
			"komenda window.action bez wskazania okna lub akcji")
	}

	akcja, err := a.akcjaZKatalogu(ctx, z.ActionId)
	if err != nil {
		return shared.WindowActionResponse{}, err
	}

	// Zapis idzie przed odmową, bo jest śladem zgłoszenia i jedynym sprawdzeniem istnienia okna.
	if _, err := a.repozytorium.ZapiszAkcje(ctx, dane.AkcjaOkna{
		Okno:      z.WindowId,
		AkcjaID:   z.ActionId,
		Parametry: surowyZapisDoTekstu(z.Parameters),
	}); err != nil {
		return shared.WindowActionResponse{}, bladDziennikaAkcjiOkna(z.WindowId, err)
	}

	return shared.WindowActionResponse{}, bladBrakuWykonawcyAkcji(z.ActionId, akcja.Komenda)
}

// akcjaZKatalogu rozstrzyga, czy zgłoszona akcja w ogóle istnieje i jest czynna.
// Bez katalogu (montaż bez wiązania) metoda odmawia zamiast zgadywać — brak
// sprawdzenia nie może uchodzić za sprawdzenie udane.
func (a *adapterPrzekazaniaOkna) akcjaZKatalogu(ctx context.Context, kod string) (dane.Akcja, error) {
	if a.akcje == nil {
		return dane.Akcja{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"window.action: katalog akcji nie jest podpięty do serwera; "+
				"nie ma czym sprawdzić akcji "+kod))
	}
	akcja, err := a.akcje.PoKodzie(ctx, kod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.Akcja{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"window.action: akcja "+kod+" nie istnieje w katalogu akcji"))
	}
	if err != nil {
		return dane.Akcja{}, protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
	}
	if !akcja.Aktywna {
		return dane.Akcja{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
			"window.action: akcja "+kod+" jest w katalogu wygaszona"))
	}
	return akcja, nil
}

// bladBrakuWykonawcyAkcji nazywa rzecz po imieniu: zgłoszenie przyjęte i
// zapisane, wykonania nie ma. Komunikat niesie kod komendy z katalogu, żeby
// odmowa była dla Operatora drogą dalej, a nie ślepym końcem.
func bladBrakuWykonawcyAkcji(kodAkcji string, komenda shared.MessageType) error {
	powod := "window.action: serwer nie wykonuje akcji katalogu. Akcja " + kodAkcji +
		" wskazuje komendę " + string(komenda) + " — wywołaj ją wprost. " +
		"Zgłoszenie zostało zapisane w dzienniku akcji okna."
	if komenda == "" {
		powod = "window.action: serwer nie wykonuje akcji katalogu, a akcja " + kodAkcji +
			" nie wskazuje żadnej komendy kontraktu. " +
			"Zgłoszenie zostało zapisane w dzienniku akcji okna."
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict, powod))
}

// surowyZapisDoTekstu przenosi ładunek JSON do wskaźnika tekstu, jak wymaga
// warstwa danych (`AkcjaOkna.Parametry`). Ładunek pusty zostaje brakiem
// wartości, nie wskaźnikiem na pusty łańcuch.
func surowyZapisDoTekstu(ladunek []byte) *string {
	if len(ladunek) == 0 {
		return nil
	}
	tekst := string(ladunek)
	return &tekst
}

// bladWskazaniaAkcjiOkna nazywa brak wymaganych danych w żądaniu jako błąd
// zgłaszającego, a nie usterkę rdzenia.
func bladWskazaniaAkcjiOkna(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, powod))
}

// bladDziennikaAkcjiOkna odróżnia „okno nie istnieje" od usterki zapisu —
// okno pokazuje wtedy inny komunikat.
func bladDziennikaAkcjiOkna(oknoId string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"window.action: okno nie istnieje: "+oknoId))
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}
