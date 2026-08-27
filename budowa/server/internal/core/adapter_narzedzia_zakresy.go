// Plik obsługuje komendy tools.scope.list i tools.scope.set oraz straż, która zapis zakresu czyta przy wykonaniu narzędzia rdzenia.
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

// adapterZakresowNarzedzi wypełnia port ZakresyNarzedzi repozytorium zakresów wraz z katalogiem profili asystenta.
type adapterZakresowNarzedzi struct {
	zakresy dane.RepozytoriumZakresowNarzedzi
	profile dane.RepozytoriumAsystenta
}

var _ ZakresyNarzedzi = (*adapterZakresowNarzedzi)(nil)

// NowyPortZakresowNarzedzi wiąże port zakresów narzędzi z repozytorium zakresów i katalogiem profili asystenta.
func NowyPortZakresowNarzedzi(zakresy dane.RepozytoriumZakresowNarzedzi,
	profile dane.RepozytoriumAsystenta) *adapterZakresowNarzedzi {

	return &adapterZakresowNarzedzi{zakresy: zakresy, profile: profile}
}

// WykazZakresow oddaje zakresy zapisane dla profilu wraz z zużyciem limitu wywołań w bieżącej sesji rozmowy.
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

// ZapiszZakres ustala zakres uprawnień i limit wywołań pozycji katalogu; pole pominięte zostaje bez zmiany.
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

// zakresNarzedziaKontraktu przekłada wiersz zakresu narzędzia na kształt pola zwracanego przez kontrakt rdzenia.
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
