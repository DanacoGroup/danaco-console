// Odpowiedzialność pliku: dwie komendy zakresu narzędzi profilu asystenta —
// `tools.scope.list` i `tools.scope.set` — oraz straż, która ich zapis czyta
// przy wykonaniu.
//
// ZAKRES JEST NASTAWĄ ZASIĘGU, NIE BRAMKĄ WBUDOWANĄ. Stanem wyjściowym jest
// pełny dostęp bez limitu: pozycja bez wiersza zakresu jest dostępna i nie ma
// granicy wywołań. Zgodnie z zasadą zero blokad platforma niczego nie zawęża
// z góry.
//
// ALE ZAWĘŻENIE ZAPISANE PRZEZ OPERATORA JEST ZAWĘŻENIEM EGZEKWOWANYM. Wiersz
// z `dostepne = 0` albo z wyczerpanym limitem kończy wywołanie odmową
// w `StrazZakresowNarzedzi` niżej — tej samej, którą pyta rdzeń przed
// skierowaniem komendy do obsługiwacza (`rdzen.go`). Zapis, którego nikt nie
// czyta przy wykonaniu, byłby gorszy niż jego brak: Operator widziałby
// ograniczenie, którego nikt nie pilnuje.
//
// Straż pyta wyłącznie o wywołania RĘKI MODELU. Klawiatura Operatora nie
// podlega zakresowi profilu asystenta: zakres opisuje, jak szeroko działa
// asystent w imieniu Operatora, a nie co wolno samemu Operatorowi
// (`sprawca.go`).
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// oknoLimituWyjsciowe to długość okna czasu limitu przyjmowana, gdy żądanie jej
// nie poda: godzina. Wartość jest nastawą, nie granicą — `tools.scope.set`
// przyjmuje własną.
const oknoLimituWyjsciowe = 3600

// adapterZakresowNarzedzi wypełnia port ZakresyNarzedzi.
type adapterZakresowNarzedzi struct {
	zakresy dane.RepozytoriumZakresowNarzedzi
	profile dane.RepozytoriumAsystenta
}

var _ ZakresyNarzedzi = (*adapterZakresowNarzedzi)(nil)

// NowyPortZakresowNarzedzi wiąże port z repozytorium zakresów i katalogiem
// profili asystenta.
func NowyPortZakresowNarzedzi(zakresy dane.RepozytoriumZakresowNarzedzi,
	profile dane.RepozytoriumAsystenta) *adapterZakresowNarzedzi {

	return &adapterZakresowNarzedzi{zakresy: zakresy, profile: profile}
}

// WykazZakresow oddaje zakresy zapisane dla profilu wraz z zużyciem limitu.
func (a *adapterZakresowNarzedzi) WykazZakresow(ctx context.Context,
	z shared.ToolsScopeListRequest) (shared.ToolsScopeListResponse, error) {

	if a == nil || a.zakresy == nil {
		return shared.ToolsScopeListResponse{Scopes: []shared.ToolScope{}}, nil
	}
	profil, err := a.kodProfilu(ctx, z.ProfileId)
	if err != nil {
		return shared.ToolsScopeListResponse{}, err
	}
	wiersze, err := a.zakresy.ZakresyProfilu(ctx, profil, strings.TrimSpace(wartoscTekstu(z.ToolName)))
	if err != nil {
		return shared.ToolsScopeListResponse{}, bladNosnikaNarzedziSesji(err)
	}
	zuzycie, err := a.zakresy.ZuzycieProfilu(ctx, profil, strings.TrimSpace(wartoscTekstu(z.SessionId)))
	if err != nil {
		return shared.ToolsScopeListResponse{}, bladNosnikaNarzedziSesji(err)
	}
	zakresy := make([]shared.ToolScope, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wiersz.ZuzyteWywolania = zuzycie[wiersz.NazwaPelna]
		zakresy = append(zakresy, zakresNarzedziaKontraktu(wiersz))
	}
	return shared.ToolsScopeListResponse{Scopes: zakresy, Total: len(zakresy)}, nil
}

// ZapiszZakres ustala zakres uprawnień i limit wywołań pozycji katalogu.
//
// Pole pominięte zostaje bez zmiany, a przy pierwszym zapisie bierze wartość
// wyjściową: dostępna, bez potwierdzenia, bez granicy. Inaczej zapis samego
// limitu odbierałby pozycję profilowi.
func (a *adapterZakresowNarzedzi) ZapiszZakres(ctx context.Context,
	z shared.ToolsScopeSetRequest) (shared.ToolsScopeSetResponse, error) {

	if a == nil || a.zakresy == nil {
		return shared.ToolsScopeSetResponse{}, bladBrakuKatalogu("zakresów narzędzi")
	}
	profil, err := a.kodProfilu(ctx, &z.ProfileId)
	if err != nil {
		return shared.ToolsScopeSetResponse{}, err
	}
	nazwa := strings.TrimSpace(z.ToolName)
	if nazwa == "" {
		return shared.ToolsScopeSetResponse{}, bladZadaniaNarzedziSesji("zakres bez wskazania pozycji katalogu")
	}
	zakres := dane.ZakresNarzedzia{
		ProfilKod: profil, NazwaPelna: nazwa,
		Dostepne: true, OknoSekund: oknoLimituWyjsciowe,
	}
	zastane, err := a.zakresy.ZakresyProfilu(ctx, profil, nazwa)
	if err != nil {
		return shared.ToolsScopeSetResponse{}, bladNosnikaNarzedziSesji(err)
	}
	if len(zastane) == 1 {
		zakres = zastane[0]
	}
	if z.Enabled != nil {
		zakres.Dostepne = *z.Enabled
	}
	if z.ConfirmRequired != nil {
		zakres.Potwierdzenie = *z.ConfirmRequired
	}
	if z.CallLimit != nil {
		if *z.CallLimit < 0 {
			return shared.ToolsScopeSetResponse{},
				bladZadaniaNarzedziSesji("granica wywołań nie bywa ujemna; zero znaczy bez granicy")
		}
		zakres.LimitWywolan = *z.CallLimit
	}
	if z.CallWindowSeconds != nil {
		if *z.CallWindowSeconds < 1 {
			return shared.ToolsScopeSetResponse{},
				bladZadaniaNarzedziSesji("okno czasu limitu jest liczbą sekund większą od zera")
		}
		zakres.OknoSekund = *z.CallWindowSeconds
	}
	if z.ArgumentAllowList != nil {
		zakres.Dopuszczone = z.ArgumentAllowList
	}
	if z.Note != nil {
		zakres.Uzasadnienie = *z.Note
	}
	zapisany, err := a.zakresy.ZapiszZakresNarzedzia(ctx, zakres)
	if err != nil {
		return shared.ToolsScopeSetResponse{}, bladNosnikaNarzedziSesji(err)
	}
	return shared.ToolsScopeSetResponse{Scope: zakresNarzedziaKontraktu(zapisany)}, nil
}

// kodProfilu rozstrzyga profil żądania. Puste wskazanie bierze profil domyślny;
// wskazanie nieznane wraca odmową, bo zakres zapisany pod profilem, którego nie
// ma, nie obowiązywałby nigdy.
func (a *adapterZakresowNarzedzi) kodProfilu(ctx context.Context, wskazanie *string) (string, error) {
	kod := strings.TrimSpace(wartoscTekstu(wskazanie))
	if a.profile == nil {
		if kod == "" {
			return "", bladBrakuKatalogu("profili asystenta")
		}
		return kod, nil
	}
	if kod == "" {
		profil, err := a.profile.ProfilDomyslny(ctx)
		if err != nil {
			return "", bladWskazania(err, "profil asystenta", "domyślny")
		}
		return profil.Kod, nil
	}
	profil, err := a.profile.Profil(ctx, kod)
	if err != nil {
		return "", bladWskazania(err, "profil asystenta", kod)
	}
	return profil.Kod, nil
}

// zakresNarzedziaKontraktu przekłada wiersz zakresu na kształt kontraktu.
func zakresNarzedziaKontraktu(w dane.ZakresNarzedzia) shared.ToolScope {
	zakres := shared.ToolScope{
		ProfileId:         w.ProfilKod,
		ToolName:          w.NazwaPelna,
		ShortName:         nazwaSkroconaPozycji(w.NazwaPelna),
		Enabled:           w.Dostepne,
		ConfirmRequired:   w.Potwierdzenie,
		CallLimit:         w.LimitWywolan,
		CallWindowSeconds: w.OknoSekund,
		ArgumentAllowList: w.Dopuszczone,
	}
	zuzyte := w.ZuzyteWywolania
	zakres.CallsUsed = &zuzyte
	zakres.Note = wskaznikTekstu(w.Uzasadnienie)
	if chwila, ok := chwilaZapisu(w.Zaktualizowano); ok {
		zakres.UpdatedAt = chwila
	}
	return zakres
}

// nazwaSkroconaPozycji zdejmuje przedrostek źródła. Nazwa bez dwukropka jest
// już skrócona i zostaje sobą.
func nazwaSkroconaPozycji(pelna string) string {
	if _, skrocona, zDwukropkiem := strings.Cut(pelna, ":"); zDwukropkiem {
		return skrocona
	}
	return pelna
}
