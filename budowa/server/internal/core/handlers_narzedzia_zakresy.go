// Plik wpina obsługę komend zakresu narzędzi tools.scope.list oraz tools.scope.set wraz z portem
// ZakresyNarzedzi i strażą sprawdzającą zapis zakresu przy każdym wywołaniu narzędzia.
package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZakresyNarzedzi jest portem zakresów uprawnień i limitów wywołań pozycji
// katalogu dla profilu asystenta.
type ZakresyNarzedzi interface {
	WykazZakresow(ctx context.Context, z shared.ToolsScopeListRequest) (shared.ToolsScopeListResponse, error)
	ZapiszZakres(ctx context.Context, z shared.ToolsScopeSetRequest) (shared.ToolsScopeSetResponse, error)
}

// zarejestrujZakresyNarzedzi wpina obsługę komend wykazu i zapisu zakresu narzędzi w rejestrze rdzenia.
func zarejestrujZakresyNarzedzi(r *Rejestr, zn ZakresyNarzedzi) {
	if r == nil || zn == nil {
		return
	}
	r.Zarejestruj(shared.CommandToolsScopeList, obsluz(zn.WykazZakresow))
	r.Zarejestruj(shared.CommandToolsScopeSet, obsluz(zn.ZapiszZakres))
}

// StrazZakresowNarzedzi rozstrzyga, czy wywołanie modelu mieści się w zakresie zapisanym dla
// profilu asystenta, i odnotowuje dopuszczone wywołanie w rachunku limitu.
type StrazZakresowNarzedzi interface {
	// SprawdzWywolanie zwraca odmowę, gdy zakres jej nie obejmuje; nil oznacza wywołanie dozwolone.
	SprawdzWywolanie(ctx context.Context, komenda shared.MessageType, idSesji string) error
}

var _ StrazZakresowNarzedzi = (*adapterZakresowNarzedzi)(nil)

// SprawdzWywolanie czyta zakres pozycji i odmawia wywołania, gdy profil wyłączył pozycję albo
// wyczerpał limit wywołań w oknie czasu, odnotowując każde wywołanie dopuszczone w rachunku.
func (a *adapterZakresowNarzedzi) SprawdzWywolanie(ctx context.Context,
	komenda shared.MessageType, idSesji string) error {

	if a == nil || a.zakresy == nil {
		return nil
	}
	profil, err := a.kodProfilu(ctx, nil)
	if err != nil {
		// Brak profilu domyślnego nie jest zawężeniem: platforma pracuje dalej bez warstwy promptu.
		return nil
	}
	nazwa := zrodloPlatformy + ":" + string(komenda)
	wiersze, err := a.zakresy.ZakresyProfilu(ctx, profil, nazwa)
	if err != nil || len(wiersze) == 0 {
		return nil
	}
	zakres := wiersze[0]
	if !zakres.Dostepne {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
			"zakres narzędzi: pozycja "+nazwa+" jest wyłączona dla profilu "+profil+
				"; Operator włączy ją w oknie zakresu narzędzi albo wykona czynność sam"))
	}
	if zakres.LimitWywolan > 0 {
		zuzycie, err := a.zakresy.ZuzycieProfilu(ctx, profil, strings.TrimSpace(idSesji))
		if err == nil && zuzycie[nazwa] >= zakres.LimitWywolan {
			return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
				"zakres narzędzi: pozycja "+nazwa+" wyczerpała granicę "+
					strconv.Itoa(zakres.LimitWywolan)+" wywołań w oknie "+
					strconv.Itoa(zakres.OknoSekund)+" sekund; Operator podniesie granicę "+
					"w oknie zakresu narzędzi albo poczeka na przesunięcie okna czasu"))
		}
	}
	// Rachunek liczy wyłącznie wywołania dopuszczone — odmowy nie podbijają limitu profilu.
	if err := a.zakresy.OdnotujWywolanieNarzedzia(ctx, profil, nazwa, strings.TrimSpace(idSesji)); err != nil {
		return nil
	}
	return nil
}

// bladZadaniaNarzedziSesji nazywa żądanie zakresu, którego nie da się spełnić.
// Odmowa jest odpowiedzią o żądaniu, nie o nośniku — stąd kod `invalid_request`,
// a nie `internal_error`.
func bladZadaniaNarzedziSesji(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"zakres narzędzi: "+powod))
}
