// Odpowiedzialność pliku: panel „Zmienne ▼” Workflow Buildera (zmienne
// przepływu i mapowanie danych między krokami), notatka przy kroku i układ
// węzłów na kanwie.
//
// Układ kanwy jest zapisem, nie wyliczeniem. Bez niego położenie węzłów
// liczyłoby się z układu zależności przy każdym otwarciu okna, a Operator
// zastawałby kanwę ułożoną od nowa — praca nad rozmieszczeniem procesu
// nie przeżyłaby zamknięcia karty.
//
// Zapis zmiennych oddaje zastrzeżenia, lecz niczego nie odmawia. Zmienna
// nieużywana i mapowanie do kroku nieistniejącego są ostrzeżeniem na kanwie,
// tak samo jak reszta walidacji definicji: zapis pozostaje możliwy, bo Operator
// buduje proces etapami i połowa definicji jest stanem poprawnym.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// UstawZmienne zapisuje zmienne przepływu i mapowanie danych między krokami.
func (a *adapterAutomatyk) UstawZmienne(ctx context.Context,
	z shared.AutomationWorkflowVariablesSetRequest) (shared.AutomationWorkflowVariablesSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationWorkflowVariablesSetResponse{}, err
	}
	if err := a.repozytorium.ZapiszZmienneAutomatyki(ctx, wiersz.ID,
		wierszeZmiennych(z.Variables)); err != nil {
		return shared.AutomationWorkflowVariablesSetResponse{}, bladAutomatyki(err)
	}
	// Pole `mappings` nieobecne zostawia mapowania zastane. Zapis samych
	// zmiennych nie ma kasować pracy nad przepływem danych — tak samo jak zapis
	// definicji bez pola `steps` nie kasuje kroków.
	if z.Mappings != nil {
		if err := a.repozytorium.ZapiszMapowaniaAutomatyki(ctx, wiersz.ID,
			wierszeMapowan(z.Mappings)); err != nil {
			return shared.AutomationWorkflowVariablesSetResponse{}, bladAutomatyki(err)
		}
	}
	zmienne, mapowania, err := a.zmienneIMapowania(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationWorkflowVariablesSetResponse{}, bladAutomatyki(err)
	}
	kroki, err := a.krokiKontraktu(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationWorkflowVariablesSetResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, &wiersz.ID, "zapis zmiennych przepływu",
		map[string]any{"zmiennych": len(zmienne), "mapowan": len(mapowania)})

	odpowiedz := shared.AutomationWorkflowVariablesSetResponse{
		Variables: zmienne, Mappings: mapowania,
	}
	if zastrzezenia := zastrzezeniaPrzeplywu(zmienne, mapowania, kroki); len(zastrzezenia) > 0 {
		odpowiedz.Issues = zastrzezenia
	}
	return odpowiedz, nil
}

// zmienneIMapowania odczytuje oba zbiory po zapisie i przekłada je na kontrakt.
func (a *adapterAutomatyk) zmienneIMapowania(ctx context.Context,
	automatykaID int64) ([]shared.AutomationVariable, []shared.AutomationDataMapping, error) {

	wierszeZmiennych, err := a.repozytorium.ZmienneAutomatyki(ctx, automatykaID)
	if err != nil {
		return nil, nil, err
	}
	wierszeMapowan, err := a.repozytorium.MapowaniaAutomatyki(ctx, automatykaID)
	if err != nil {
		return nil, nil, err
	}
	zmienne := make([]shared.AutomationVariable, 0, len(wierszeZmiennych))
	for _, wiersz := range wierszeZmiennych {
		zmienna := shared.AutomationVariable{
			Name: wiersz.Nazwa, Kind: wiersz.Rodzaj, SecretRef: wiersz.OdwolanieSekretu,
		}
		if wiersz.WartoscDomyslna != nil {
			zmienna.DefaultValue = []byte(*wiersz.WartoscDomyslna)
		}
		zmienne = append(zmienne, zmienna)
	}
	mapowania := make([]shared.AutomationDataMapping, 0, len(wierszeMapowan))
	for _, wiersz := range wierszeMapowan {
		mapowania = append(mapowania, shared.AutomationDataMapping{
			FromStepId: wiersz.KrokZ, FromPath: wiersz.SciezkaZ,
			ToStepId: wiersz.KrokDo, ToField: wiersz.PoleDo, Template: wiersz.Szablon,
		})
	}
	return zmienne, mapowania, nil
}

// zastrzezeniaPrzeplywu nazywa to, co w przepływie danych nie trzyma się kupy:
// mapowanie wskazujące krok, którego nie ma, oraz zmienną, po którą nie sięga
// żaden krok ani żadne mapowanie.
func zastrzezeniaPrzeplywu(zmienne []shared.AutomationVariable,
	mapowania []shared.AutomationDataMapping, kroki []shared.AutomationStep) []string {

	istniejace := make(map[string]bool, len(kroki))
	for _, krok := range kroki {
		istniejace[krok.Id] = true
	}
	zastrzezenia := []string{}
	for _, mapowanie := range mapowania {
		if !istniejace[mapowanie.FromStepId] {
			zastrzezenia = append(zastrzezenia,
				"mapowanie wskazuje krok źródłowy, którego nie ma: "+mapowanie.FromStepId)
		}
		if !istniejace[mapowanie.ToStepId] {
			zastrzezenia = append(zastrzezenia,
				"mapowanie wskazuje krok docelowy, którego nie ma: "+mapowanie.ToStepId)
		}
	}
	for _, zmienna := range zmienne {
		if !czyZmiennaUzywana(zmienna.Name, mapowania, kroki) {
			zastrzezenia = append(zastrzezenia, "zmienna nieużywana: "+zmienna.Name)
		}
	}
	return zastrzezenia
}

// czyZmiennaUzywana szuka nazwy zmiennej w mapowaniach i w parametrach kroków.
// Wyszukanie idzie po zawartości, bo zmienną przywołuje się w szablonie pola
// albo w ładunku kroku, a rdzeń nie narzuca składni przywołania.
func czyZmiennaUzywana(nazwa string, mapowania []shared.AutomationDataMapping,
	kroki []shared.AutomationStep) bool {

	if nazwa == "" {
		return true
	}
	for _, mapowanie := range mapowania {
		if mapowanie.ToField == nazwa || zawieraNazwe(wartoscTekstu(mapowanie.Template), nazwa) {
			return true
		}
	}
	for _, krok := range kroki {
		if zawieraNazwe(string(krok.Params), nazwa) ||
			zawieraNazwe(wartoscTekstu(krok.Condition), nazwa) {
			return true
		}
	}
	return false
}

// zawieraNazwe mówi, czy tekst przywołuje nazwę zmiennej.
func zawieraNazwe(tekst, nazwa string) bool {
	if tekst == "" {
		return false
	}
	for numer := 0; numer+len(nazwa) <= len(tekst); numer++ {
		if tekst[numer:numer+len(nazwa)] == nazwa {
			return true
		}
	}
	return false
}

// wierszeZmiennych przekłada zmienne kontraktu na wiersze.
func wierszeZmiennych(zmienne []shared.AutomationVariable) []dane.ZmiennaAutomatyki {
	wiersze := make([]dane.ZmiennaAutomatyki, 0, len(zmienne))
	for _, zmienna := range zmienne {
		wiersz := dane.ZmiennaAutomatyki{
			Nazwa: zmienna.Name, Rodzaj: zmienna.Kind, OdwolanieSekretu: zmienna.SecretRef,
		}
		if len(zmienna.DefaultValue) > 0 {
			wartosc := string(zmienna.DefaultValue)
			wiersz.WartoscDomyslna = &wartosc
		}
		wiersze = append(wiersze, wiersz)
	}
	return wiersze
}

// wierszeMapowan przekłada mapowania kontraktu na wiersze.
func wierszeMapowan(mapowania []shared.AutomationDataMapping) []dane.MapowanieDanych {
	wiersze := make([]dane.MapowanieDanych, 0, len(mapowania))
	for _, mapowanie := range mapowania {
		wiersze = append(wiersze, dane.MapowanieDanych{
			KrokZ: mapowanie.FromStepId, SciezkaZ: mapowanie.FromPath,
			KrokDo: mapowanie.ToStepId, PoleDo: mapowanie.ToField, Szablon: mapowanie.Template,
		})
	}
	return wiersze
}

// UstawNotatkeKroku zapisuje notatkę opisową przy kroku; treść pusta ją zdejmuje.
func (a *adapterAutomatyk) UstawNotatkeKroku(ctx context.Context,
	z shared.AutomationStepNoteSetRequest) (shared.AutomationStepNoteSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationStepNoteSetResponse{}, err
	}
	if z.StepId == "" {
		return shared.AutomationStepNoteSetResponse{},
			bladWskazaniaAutomatyki("notatka wymaga wskazania kroku (stepId)")
	}
	var notatka *string
	if z.Note != "" {
		tresc := z.Note
		notatka = &tresc
	}
	if err := a.repozytorium.UstawNotatkeKroku(ctx, wiersz.ID, z.StepId, notatka); err != nil {
		return shared.AutomationStepNoteSetResponse{}, bladAutomatyki(err)
	}
	kroki, err := a.krokiKontraktu(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationStepNoteSetResponse{}, bladAutomatyki(err)
	}
	for _, krok := range kroki {
		if krok.Id == z.StepId {
			return shared.AutomationStepNoteSetResponse{Step: krok}, nil
		}
	}
	// Notatka przy kroku, którego w definicji nie ma, zostaje zapisana — krok
	// przywrócony z wersji wcześniejszej zastanie ją na miejscu. Odpowiedź musi
	// jednak nieść krok, więc brak kroku nazywa się wprost.
	return shared.AutomationStepNoteSetResponse{},
		bladNieznanegoBytuAutomatyki(dane.ErrBrakWiersza, "krok nie istnieje w definicji: ", z.StepId)
}

// UstawUkladKanwy zapisuje położenia węzłów kroków na kanwie.
func (a *adapterAutomatyk) UstawUkladKanwy(ctx context.Context,
	z shared.AutomationStepLayoutSetRequest) (shared.AutomationStepLayoutSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationStepLayoutSetResponse{}, err
	}
	polozenia := make([]dane.AdnotacjaKroku, 0, len(z.Positions))
	for _, polozenie := range z.Positions {
		polozenia = append(polozenia, dane.AdnotacjaKroku{
			KrokKod: polozenie.StepId, X: polozenie.X, Y: polozenie.Y,
		})
	}
	if err := a.repozytorium.ZapiszPolozeniaKrokow(ctx, wiersz.ID, polozenia); err != nil {
		return shared.AutomationStepLayoutSetResponse{}, bladAutomatyki(err)
	}
	adnotacje, err := a.repozytorium.AdnotacjeKrokow(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationStepLayoutSetResponse{}, bladAutomatyki(err)
	}
	zapisane := make([]shared.AutomationStepPosition, 0, len(adnotacje))
	for _, adnotacja := range adnotacje {
		zapisane = append(zapisane, shared.AutomationStepPosition{
			StepId: adnotacja.KrokKod, X: adnotacja.X, Y: adnotacja.Y,
		})
	}
	return shared.AutomationStepLayoutSetResponse{Positions: zapisane}, nil
}
