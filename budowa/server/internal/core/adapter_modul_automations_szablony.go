// Odpowiedzialność pliku: biblioteka szablonów przepływów — zapis definicji
// jako szablonu wraz z parametrami, wykaz biblioteki i założenie automatyki z
// szablonu. Podstawienie idzie po nazwie parametru w zapisie strukturalnym
// kroków.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ZapiszSzablon zapisuje definicję automatyki jako szablon przepływu wraz z
// parametrami wykrytymi w zapisie kroków.
func (a *adapterAutomatyk) ZapiszSzablon(ctx context.Context,
	z shared.AutomationTemplateSaveRequest) (shared.AutomationTemplateSaveResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationTemplateSaveResponse{}, err
	}
	kroki, err := a.krokiKontraktu(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationTemplateSaveResponse{}, bladAutomatyki(err)
	}
	zapisKrokow := zapisStrukturalny(kroki)
	if zapisKrokow == nil {
		return shared.AutomationTemplateSaveResponse{},
			bladWskazaniaAutomatyki("kroków automatyki nie da się zapisać jako szablonu")
	}
	nazwa := z.Name
	if nazwa == "" {
		nazwa = wiersz.Nazwa
	}
	zapisany, err := a.repozytorium.ZapiszSzablonAutomatyki(ctx, dane.SzablonAutomatyki{
		Kod: nowyIdentyfikator(przedrostekSzablonuAutomatyki), Nazwa: nazwa,
		Opis: z.Description, Kroki: *zapisKrokow,
	}, wierszeParametrow(z.Parameters))
	if err != nil {
		return shared.AutomationTemplateSaveResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, &wiersz.ID, "zapis szablonu przepływu",
		map[string]any{"szablon": zapisany.Kod, "nazwa": nazwa})
	szablon, err := a.szablonKontraktu(ctx, zapisany)
	if err != nil {
		return shared.AutomationTemplateSaveResponse{}, bladAutomatyki(err)
	}
	return shared.AutomationTemplateSaveResponse{Template: szablon}, nil
}

// WykazSzablonow oddaje bibliotekę szablonów przepływów w kolejności nazw,
// wraz z ich parametrami zapisu.
func (a *adapterAutomatyk) WykazSzablonow(ctx context.Context,
	z shared.AutomationTemplateListRequest) (shared.AutomationTemplateListResponse, error) {

	wiersze, err := a.repozytorium.SzablonyAutomatyki(ctx, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.AutomationTemplateListResponse{}, bladAutomatyki(err)
	}
	szablony := make([]shared.AutomationTemplate, 0, len(wiersze))
	for _, wiersz := range wiersze {
		szablon, err := a.szablonKontraktu(ctx, wiersz)
		if err != nil {
			return shared.AutomationTemplateListResponse{}, bladAutomatyki(err)
		}
		szablony = append(szablony, szablon)
	}
	return shared.AutomationTemplateListResponse{Templates: szablony}, nil
}

// ZastosujSzablon zakłada automatykę z szablonu, podstawiając wartości jego
// parametrów; brakujące parametry wracają nazwane, nie wstrzymują zapisu.
func (a *adapterAutomatyk) ZastosujSzablon(ctx context.Context,
	z shared.AutomationTemplateApplyRequest) (shared.AutomationTemplateApplyResponse, error) {

	if z.TemplateId == "" {
		return shared.AutomationTemplateApplyResponse{},
			bladWskazaniaAutomatyki("zastosowanie szablonu bez wskazania szablonu")
	}
	wiersz, err := a.repozytorium.SzablonAutomatyki(ctx, z.TemplateId)
	if err != nil {
		return shared.AutomationTemplateApplyResponse{},
			bladNieznanegoBytuAutomatyki(err, "szablon nie istnieje: ", z.TemplateId)
	}
	parametry, err := a.repozytorium.ParametrySzablonuAutomatyki(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationTemplateApplyResponse{}, bladAutomatyki(err)
	}
	wartosci := wartosciParametrow(z.Values)
	zapis, brakujace := podstawParametry(wiersz.Kroki, parametry, wartosci)

	nazwa := z.Name
	if nazwa == "" {
		nazwa = wiersz.Nazwa
	}
	odpowiedz, err := a.Zapisz(ctx, shared.AutomationWorkflowSaveRequest{
		Name: nazwa, Description: wiersz.Opis, Steps: krokiZMigawki(zapis),
	})
	if err != nil {
		return shared.AutomationTemplateApplyResponse{}, err
	}
	wynik := shared.AutomationTemplateApplyResponse{Workflow: odpowiedz.Workflow}
	if len(brakujace) > 0 {
		wynik.MissingParameters = brakujace
	}
	return wynik, nil
}

// podstawParametry wstawia wartości parametrów w zapis kroków i nazywa te, dla
// których wartości nie podano. Parametr z wartością domyślną nie jest brakujący
// — domyślna jest wartością, a nie jej brakiem.
func podstawParametry(zapisKrokow string, parametry []dane.ParametrSzablonu,
	wartosci map[string]string) (string, []string) {

	brakujace := []string{}
	wynik := zapisKrokow
	for _, parametr := range parametry {
		wartosc, podano := wartosci[parametr.Nazwa]
		if !podano && parametr.WartoscDomyslna != nil {
			wartosc, podano = *parametr.WartoscDomyslna, true
		}
		if !podano {
			brakujace = append(brakujace, parametr.Nazwa)
			continue
		}
		wynik = strings.ReplaceAll(wynik, "{{"+parametr.Nazwa+"}}", oczyszczonaWartosc(wartosc))
	}
	return wynik, brakujace
}

// oczyszczonaWartosc przygotowuje wartość do wstawienia w zapis strukturalny:
// zdejmuje cudzysłowy zapisu tekstowego i znaki, które rozbiłyby zapis. Bez
// tego wartość z cudzysłowem zamieniałaby poprawną definicję w zapis nieczytelny.
func oczyszczonaWartosc(wartosc string) string {
	tresc, err := json.Marshal(wartosc)
	if err != nil {
		return ""
	}
	return strings.Trim(string(tresc), `"`)
}

// wartosciParametrow rozbiera zapis strukturalny wartości na mapę napisów.
// Wartość złożona (wykaz, zapis zagnieżdżony) wraca w postaci zapisu — tak samo
// wchodzi w ładunek kroku.
func wartosciParametrow(zapis json.RawMessage) map[string]string {
	wartosci := map[string]string{}
	if len(zapis) == 0 {
		return wartosci
	}
	surowe := map[string]json.RawMessage{}
	if err := json.Unmarshal(zapis, &surowe); err != nil {
		return wartosci
	}
	for nazwa, wartosc := range surowe {
		tekst := ""
		if err := json.Unmarshal(wartosc, &tekst); err == nil {
			wartosci[nazwa] = tekst
			continue
		}
		wartosci[nazwa] = string(wartosc)
	}
	return wartosci
}

// szablonKontraktu składa szablon kontraktu wraz z jego parametrami,
// przekładając wiersz biblioteki na byt widoczny wołającemu.
func (a *adapterAutomatyk) szablonKontraktu(ctx context.Context,
	wiersz dane.SzablonAutomatyki) (shared.AutomationTemplate, error) {

	parametry, err := a.repozytorium.ParametrySzablonuAutomatyki(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationTemplate{}, err
	}
	szablon := shared.AutomationTemplate{
		Id: wiersz.Kod, Name: wiersz.Nazwa, Description: wiersz.Opis,
		Steps: krokiZMigawki(wiersz.Kroki), UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}
	if len(parametry) > 0 {
		szablon.Parameters = parametryKontraktu(parametry)
	}
	return szablon, nil
}

// parametryKontraktu przekłada wiersze parametrów szablonu z bazy na byty
// kontraktu, w kolejności ich zapisu.
func parametryKontraktu(wiersze []dane.ParametrSzablonu) []shared.AutomationTemplateParameter {
	parametry := make([]shared.AutomationTemplateParameter, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wymagany := wiersz.Wymagany
		parametr := shared.AutomationTemplateParameter{
			Name: wiersz.Nazwa, Label: wiersz.Etykieta, Required: &wymagany,
		}
		if wiersz.WartoscDomyslna != nil {
			parametr.DefaultValue = []byte(*wiersz.WartoscDomyslna)
		}
		parametry = append(parametry, parametr)
	}
	return parametry
}

// wierszeParametrow przekłada parametry kontraktu przyjęte żądaniem na wiersze
// zapisywane wraz z szablonem.
func wierszeParametrow(parametry []shared.AutomationTemplateParameter) []dane.ParametrSzablonu {
	wiersze := make([]dane.ParametrSzablonu, 0, len(parametry))
	for _, parametr := range parametry {
		wiersz := dane.ParametrSzablonu{
			Nazwa: parametr.Name, Etykieta: parametr.Label,
			Wymagany: parametr.Required != nil && *parametr.Required,
		}
		if len(parametr.DefaultValue) > 0 {
			wartosc := string(parametr.DefaultValue)
			wiersz.WartoscDomyslna = &wartosc
		}
		wiersze = append(wiersze, wiersz)
	}
	return wiersze
}
