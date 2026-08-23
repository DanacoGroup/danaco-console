// Odpowiedzialność pliku: komenda `window.action` — metoda WykonajAkcje na
// *adapterPrzekazaniaOkna. Typ i konstruktor deklaruje
// `adapter_okno_przekazanie.go`; ten plik dokłada wyłącznie metodę obszaru
// dziennika akcji, tak jak `adapter_modul_automations_kroki.go` dokłada
// metody obok typu zadeklarowanego w `adapter_modul_automations.go`.
//
// Metoda nie wykonuje akcji katalogu — sprawdza ją, zapisuje ślad zgłoszenia i
// odmawia wprost. Meldunek o wykonaniu czynności, której nikt nie wykonał, jest
// gorszy od odmowy: Operator odchodzi od ekranu przekonany, że rzecz się stała.
//
// Katalog akcji niesie kod komendy (`dane.Akcja.Komenda`: `home.enter`,
// `config.set`, `session.create` …), więc wykonanie jest kiedyś osiągalne. Nie
// da się go jednak domknąć teraz i uczciwie, bo:
//
//   - obsługiwacze rdzenia (`core.Rejestr`, `wpisy map[MessageType]Obsluga`)
//     przyjmują `protocol.Request` — pełną kopertę z tożsamością klienta,
//     wiązaniem sesji i identyfikatorem żądania. Port `PrzekazanieOkna` dostaje
//     samą treść (`shared.WindowActionRequest`), bez koperty. Sklejenie koperty
//     zastępczej znaczyłoby wykonanie komendy w cudzym albo w żadnym kontekście
//     sesji — to zamiana fałszywego meldunku na fałszywy skutek;
//   - `parameters` komendy `window.action` to surowy JSON o kształcie zależnym
//     od akcji; kontrakt nie mówi, że jest to treść żądania komendy docelowej.
//     Przyjęcie tego założenia byłoby zgadywaniem kontraktu.
//
// Dlatego metoda robi trzy rzeczy, do których ma pełne pokrycie, i ani jednej
// więcej: sprawdza akcję w katalogu, zostawia trwały ślad zgłoszenia w dzienniku
// i ODMAWIA WPROST kodem kontraktu,
// podając Operatorowi kod komendy, którą ma wywołać sam. Odmowa niesie tę samą
// wiedzę co wykonanie, minus nieprawda.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// WykonajAkcje przyjmuje zgłoszenie akcji panelu okna (`window.action`),
// odnotowuje je w dzienniku i odmawia wykonania — rdzeń nie ma wykonawcy akcji
// katalogu (patrz komentarz nagłówkowy pliku). Sukcesu ta metoda nie zwraca w
// żadnej ścieżce; gdy wykonanie dojdzie, sukces pojawi się razem z nim.
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

	// Zapis idzie PRZED odmową wykonania, bo pełni dwie role naraz: jest śladem
	// zgłoszenia (kto, w którym oknie, o co prosił) i jedynym sprawdzeniem
	// istnienia okna — nieznane okno odrzuca repozytorium samo
	// (`dane.ErrBrakWiersza`). Nieudany zapis wychodzi swoją odmową, żeby
	// „okna nie ma" nie zlało się z „rdzeń nie wykonuje akcji".
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
			"window.action: katalog akcji nie jest podpięty do rdzenia; "+
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
	powod := "window.action: rdzeń nie wykonuje akcji katalogu. Akcja " + kodAkcji +
		" wskazuje komendę " + string(komenda) + " — wywołaj ją wprost. " +
		"Zgłoszenie zostało zapisane w dzienniku akcji okna."
	if komenda == "" {
		powod = "window.action: rdzeń nie wykonuje akcji katalogu, a akcja " + kodAkcji +
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

// bladWskazaniaAkcjiOkna nazywa brak danych w żądaniu — to błąd Operatora,
// nie rdzenia.
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
