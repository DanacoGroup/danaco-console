// Wpięcie dwóch komend zakresu narzędzi profilu asystenta —
// `tools.scope.list` i `tools.scope.set` — wraz z portem ZakresyNarzedzi
// i strażą, która ich zapis czyta przy wykonaniu.
//
// Port osobny od NarzedziaSesji: tamten opisuje DOŁOŻENIE narzędzia na czas
// sesji i wykaz po ukośniku, ten — zakres i limit wywołań właściwy profilowi.
// Dwa różne pytania nad tym samym katalogiem.
//
// Zdarzenia nie ma: kontrakt nie daje rodzinie `tools.scope.*` żadnego
// zdarzenia, więc rdzeń go nie wymyśla. Okno odświeża wykaz po odpowiedzi.
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

// zarejestrujZakresyNarzedzi wpina dwa uchwyty rodziny.
func zarejestrujZakresyNarzedzi(r *Rejestr, zn ZakresyNarzedzi) {
	if r == nil || zn == nil {
		return
	}
	r.Zarejestruj(shared.CommandToolsScopeList, obsluz(zn.WykazZakresow))
	r.Zarejestruj(shared.CommandToolsScopeSet, obsluz(zn.ZapiszZakres))
}

// StrazZakresowNarzedzi rozstrzyga, czy wywołanie ręki modelu mieści się
// w zakresie zapisanym dla profilu asystenta, i odnotowuje je w rachunku limitu.
//
// Port stoi po stronie odbiorcy — pyta go rdzeń przed skierowaniem komendy do
// obsługiwacza (`rdzen.go`). Straż niewpięta nie zmienia niczego: stanem
// wyjściowym platformy jest pełny dostęp bez granicy.
type StrazZakresowNarzedzi interface {
	// SprawdzWywolanie zwraca odmowę, gdy Operator zawęził zakres tak, że
	// wywołanie się w nim nie mieści. Nil znaczy „wolno" — i to jest odpowiedź
	// zwykła, bo wiersz zakresu ma mniejszość pozycji katalogu.
	SprawdzWywolanie(ctx context.Context, komenda shared.MessageType, idSesji string) error
}

var _ StrazZakresowNarzedzi = (*adapterZakresowNarzedzi)(nil)

// SprawdzWywolanie czyta zakres pozycji i odmawia, gdy Operator ją wyłączył
// albo wyczerpał się limit wywołań w oknie czasu. Wywołanie dopuszczone jest
// odnotowywane — bez tego zapisu limit byłby liczbą, której nikt nie zużywa.
//
// Pozycja bez wiersza zakresu przechodzi bez zapytania i bez rachunku: stanem
// wyjściowym jest pełny dostęp, a rachunek prowadzony dla wszystkich pozycji
// dopisywałby wiersz przy każdym wywołaniu narzędzia w produkcie.
func (a *adapterZakresowNarzedzi) SprawdzWywolanie(ctx context.Context,
	komenda shared.MessageType, idSesji string) error {

	if a == nil || a.zakresy == nil {
		return nil
	}
	profil, err := a.kodProfilu(ctx, nil)
	if err != nil {
		// Brak profilu domyślnego nie jest zawężeniem. Platforma bez ani jednego
		// profilu ma pracować dalej — tak samo, jak pracuje bez warstwy promptu.
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
	// Rachunek prowadzimy wyłącznie dla pozycji objętych zakresem — i wyłącznie
	// po dopuszczeniu. Zliczanie odmów kazałoby limitowi rosnąć od prób, które
	// niczego nie wykonały.
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
