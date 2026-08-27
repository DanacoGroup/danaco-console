// Odpowiedzialność pliku: reguły alarmowania Execution Monitora, budżety
// czasu przebiegu i kroku, skarbiec poświadczeń oraz odczyt dziennika
// audytu. Skarbiec ma jedną zasadę: wartość poświadczenia nie wraca nigdy
// w odpowiedzi.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// SejfPoswiadczenAutomatyki jest portem magazynu wartości poświadczeń
// leżącego POZA bazą. Wypełnia go `dane.SejfPlikowy` — ten sam byt, który
// trzyma sekrety kont i punktów dostępu.
type SejfPoswiadczenAutomatyki interface {
	Zapisz(ctx context.Context, byt, poswiadczenie string) (string, error)
	Usun(ctx context.Context, byt string) error
}

// UstawReguleAlarmowania zapisuje warunek i kanały powiadomień reguły; brak
// identyfikatora zakłada nową regułę, istniejący nadpisuje starą.
func (a *adapterAutomatyk) UstawReguleAlarmowania(ctx context.Context,
	z shared.AutomationAlertRuleSetRequest) (shared.AutomationAlertRuleSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationAlertRuleSetResponse{}, err
	}
	kanaly := zapisStrukturalny(z.Channels)
	if kanaly == nil {
		pusty := "[]"
		kanaly = &pusty
	}
	kod := wartoscTekstu(z.RuleId)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekReguly)
	}
	zapisana, err := a.repozytorium.ZapiszRegulealarmowania(ctx, dane.RegulaAlarmowania{
		Kod: kod, AutomatykaID: wiersz.ID, Wyzwalacz: string(z.Trigger),
		Warunek: z.Condition, Kanaly: *kanaly, Czynna: czyCzynna(z.Enabled),
	})
	if err != nil {
		return shared.AutomationAlertRuleSetResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, &wiersz.ID, "zapis reguły alarmowania",
		map[string]any{"regula": kod, "wyzwalacz": string(z.Trigger)})
	return shared.AutomationAlertRuleSetResponse{Rule: regulaAlarmowaniaAutomatyki(zapisana)}, nil
}

// WykazRegulAlarmowania oddaje reguły automatyki albo, bez wskazanej
// automatyki, komplet reguł należących do Operatora.
func (a *adapterAutomatyk) WykazRegulAlarmowania(ctx context.Context,
	z shared.AutomationAlertRuleListRequest) (shared.AutomationAlertRuleListResponse, error) {

	var automatykaID int64
	if kod := wartoscTekstu(z.WorkflowId); kod != "" {
		wiersz, err := a.wiersz(ctx, kod)
		if err != nil {
			return shared.AutomationAlertRuleListResponse{}, err
		}
		automatykaID = wiersz.ID
	}
	wiersze, err := a.repozytorium.RegulyAlarmowania(ctx, automatykaID)
	if err != nil {
		return shared.AutomationAlertRuleListResponse{}, bladAutomatyki(err)
	}
	reguly := make([]shared.AutomationAlertRule, 0, len(wiersze))
	for _, wiersz := range wiersze {
		reguly = append(reguly, regulaAlarmowaniaAutomatyki(wiersz))
	}
	return shared.AutomationAlertRuleListResponse{Rules: reguly}, nil
}

// UstawBudzetyPrzebiegu zapisuje maksymalny czas przebiegu i kroku. Zero znaczy
// „bez granicy” — tak mówi kontrakt.
func (a *adapterAutomatyk) UstawBudzetyPrzebiegu(ctx context.Context,
	z shared.AutomationExecutionBudgetSetRequest) (shared.AutomationExecutionBudgetSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationExecutionBudgetSetResponse{}, err
	}
	err = a.repozytorium.UstawBudzetyAutomatyki(ctx, wiersz.ID,
		wartoscLiczby(z.RunBudgetSeconds), wartoscLiczby(z.StepBudgetSeconds), z.AlertRuleId)
	if err != nil {
		return shared.AutomationExecutionBudgetSetResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, &wiersz.ID, "zapis budżetów czasu", map[string]any{
		"przebieg": wartoscLiczby(z.RunBudgetSeconds), "krok": wartoscLiczby(z.StepBudgetSeconds),
	})
	automatyka, err := a.automatykaPoZmianie(ctx, wiersz.Kod)
	if err != nil {
		return shared.AutomationExecutionBudgetSetResponse{}, err
	}
	return shared.AutomationExecutionBudgetSetResponse{Workflow: automatyka}, nil
}

// ZapiszPoswiadczenie umieszcza wartość w sejfie i zapisuje jej referencję.
// Kolejność jest rozmyślna: najpierw sejf, potem baza. Referencja wskazująca
// wpis, którego w sejfie nie ma, byłaby obietnicą bez pokrycia.
func (a *adapterAutomatyk) ZapiszPoswiadczenie(ctx context.Context,
	z shared.AutomationSecretSetRequest) (shared.AutomationSecretSetResponse, error) {

	if a.sejf == nil {
		return shared.AutomationSecretSetResponse{}, errBrakSkarbca
	}
	if z.Name == "" {
		return shared.AutomationSecretSetResponse{},
			bladWskazaniaAutomatyki("poświadczenie bez nazwy nie ma jak zostać przywołane z kroku")
	}
	byt := nowyIdentyfikator(przedrostekPoswiadczenia)
	odwolanie, err := a.sejf.Zapisz(ctx, byt, z.Value)
	if err != nil {
		return shared.AutomationSecretSetResponse{}, bladAutomatyki(err)
	}
	poswiadczenie := dane.PoswiadczenieAutomatyki{
		Odwolanie: odwolanie, Nazwa: z.Name, ZasiegID: z.ScopeId,
	}
	if z.Scope != nil {
		zasieg := string(*z.Scope)
		poswiadczenie.Zasieg = &zasieg
	}
	if err := a.repozytorium.ZapiszPoswiadczenieAutomatyki(ctx, poswiadczenie); err != nil {
		// Baza odmówiła, więc osierocony wpis sejfu zostaje zdjęty od razu.
		_ = a.sejf.Usun(ctx, byt)
		return shared.AutomationSecretSetResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, nil, "zapis poświadczenia w skarbcu", map[string]any{"nazwa": z.Name})

	zapisane, err := a.repozytorium.PoswiadczenieAutomatykiPoOdwolaniu(ctx, odwolanie)
	if err != nil {
		return shared.AutomationSecretSetResponse{}, bladAutomatyki(err)
	}
	return shared.AutomationSecretSetResponse{Secret: poswiadczenieKontraktu(zapisane)}, nil
}

// WykazPoswiadczen oddaje referencje poświadczeń dostępne krokom automatyki.
// Wartości poświadczeń nie wracają nigdy w odpowiedzi.
func (a *adapterAutomatyk) WykazPoswiadczen(ctx context.Context,
	z shared.AutomationSecretListRequest) (shared.AutomationSecretListResponse, error) {

	zasieg := ""
	if z.Scope != nil {
		zasieg = string(*z.Scope)
	}
	wiersze, err := a.repozytorium.PoswiadczeniaAutomatyki(ctx, zasieg, wartoscTekstu(z.ScopeId))
	if err != nil {
		return shared.AutomationSecretListResponse{}, bladAutomatyki(err)
	}
	referencje := make([]shared.AutomationSecretRef, 0, len(wiersze))
	for _, wiersz := range wiersze {
		referencje = append(referencje, poswiadczenieKontraktu(wiersz))
	}
	return shared.AutomationSecretListResponse{Secrets: referencje}, nil
}

// UsunPoswiadczenie zdejmuje referencję i wartość, nazywając kroki, które
// straciły pokrycie. Usunięcie nie jest wstrzymywane tym, że kroki
// poświadczenie przywołują.
func (a *adapterAutomatyk) UsunPoswiadczenie(ctx context.Context,
	z shared.AutomationSecretRemoveRequest) (shared.AutomationSecretRemoveResponse, error) {

	if z.SecretRef == "" {
		return shared.AutomationSecretRemoveResponse{},
			bladWskazaniaAutomatyki("usunięcie poświadczenia bez wskazania referencji")
	}
	przywolujace, err := a.krokiPrzywolujace(ctx, z.SecretRef)
	if err != nil {
		return shared.AutomationSecretRemoveResponse{}, err
	}
	usuniete, err := a.repozytorium.UsunPoswiadczenieAutomatyki(ctx, z.SecretRef)
	if err != nil {
		return shared.AutomationSecretRemoveResponse{}, bladAutomatyki(err)
	}
	if usuniete && a.sejf != nil {
		// Sejf kluczuje bytem, a referencja niesie przedrostek `sejf:`, który
		// rozbiera wołający.
		_ = a.sejf.Usun(ctx, bytSejfuZOdwolania(z.SecretRef))
	}
	a.zapisAudytu(ctx, nil, "usunięcie poświadczenia ze skarbca",
		map[string]any{"odwolanie": z.SecretRef, "usuniete": usuniete})

	odpowiedz := shared.AutomationSecretRemoveResponse{Removed: usuniete}
	if len(przywolujace) > 0 {
		odpowiedz.ReferencingStepIds = przywolujace
	}
	return odpowiedz, nil
}

// krokiPrzywolujace szuka kroków, które przywołują referencję. Krok przywołuje
// poświadczenie polem `secretRefs` albo zmienną przepływu wskazującą sekret,
// więc oba miejsca są przeszukiwane.
func (a *adapterAutomatyk) krokiPrzywolujace(ctx context.Context, odwolanie string) ([]string, error) {
	automatyki, err := a.repozytorium.Automatyki(ctx, false, 0)
	if err != nil {
		return nil, bladAutomatyki(err)
	}
	kody := []string{}
	for _, automatyka := range automatyki {
		kroki, err := a.krokiKontraktu(ctx, automatyka.ID)
		if err != nil {
			continue
		}
		nazwyZmiennych := nazwySekretowychZmiennych(ctx, a, automatyka.ID, odwolanie)
		for _, krok := range kroki {
			if krokPrzywolujeSekret(krok, odwolanie, nazwyZmiennych) {
				kody = append(kody, krok.Id)
			}
		}
	}
	return kody, nil
}

// nazwySekretowychZmiennych zbiera nazwy zmiennych przepływu wskazujących
// wskazaną referencję — krok przywołuje je po nazwie, nie po referencji.
func nazwySekretowychZmiennych(ctx context.Context, a *adapterAutomatyk, automatykaID int64,
	odwolanie string) []string {

	zmienne, err := a.repozytorium.ZmienneAutomatyki(ctx, automatykaID)
	if err != nil {
		return nil
	}
	nazwy := []string{}
	for _, zmienna := range zmienne {
		if zmienna.OdwolanieSekretu != nil && *zmienna.OdwolanieSekretu == odwolanie {
			nazwy = append(nazwy, zmienna.Nazwa)
		}
	}
	return nazwy
}

// krokPrzywolujeSekret rozstrzyga, czy dany krok stracił pokrycie sekretem
// po usunięciu poświadczenia z magazynu.
func krokPrzywolujeSekret(krok shared.AutomationStep, odwolanie string, nazwy []string) bool {
	for _, referencja := range krok.SecretRefs {
		if referencja == odwolanie {
			return true
		}
	}
	for _, nazwa := range nazwy {
		if zawieraNazwe(string(krok.Params), nazwa) {
			return true
		}
	}
	return false
}

// OdczytajAudyt oddaje dziennik audytu automatyki od najnowszego wpisu,
// w porcjach ograniczonych żądaniem wołającego.
func (a *adapterAutomatyk) OdczytajAudyt(ctx context.Context,
	z shared.AutomationAuditListRequest) (shared.AutomationAuditListResponse, error) {

	var automatykaID int64
	if kod := wartoscTekstu(z.WorkflowId); kod != "" {
		wiersz, err := a.wiersz(ctx, kod)
		if err != nil {
			return shared.AutomationAuditListResponse{}, err
		}
		automatykaID = wiersz.ID
	}
	wiersze, err := a.repozytorium.AudytAutomatyki(ctx, automatykaID,
		znacznikChwiliWskazanejAutomatyzacji(z.FromAt), znacznikChwiliWskazanejAutomatyzacji(z.ToAt),
		wartoscLiczby(z.Limit))
	if err != nil {
		return shared.AutomationAuditListResponse{}, bladAutomatyki(err)
	}
	wpisy := make([]shared.AutomationAuditEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wpis := shared.AutomationAuditEntry{
			Id: wiersz.Kod, WorkflowId: wiersz.AutomatykaKod, ActorId: wiersz.Wykonawca,
			Action: wiersz.Czynnosc, At: chwilaBazy(wiersz.Chwila),
		}
		if wiersz.Szczegoly != nil {
			wpis.Detail = []byte(*wiersz.Szczegoly)
		}
		wpisy = append(wpisy, wpis)
	}
	return shared.AutomationAuditListResponse{Entries: wpisy}, nil
}

// regulaAlarmowaniaAutomatyki przekłada wiersz reguły z bazy danych na byt
// kontraktu zwracany wołającemu.
func regulaAlarmowaniaAutomatyki(wiersz dane.RegulaAlarmowania) shared.AutomationAlertRule {
	return shared.AutomationAlertRule{
		Id: wiersz.Kod, WorkflowId: wiersz.AutomatykaKod,
		Trigger: shared.AutomationAlertTrigger(wiersz.Wyzwalacz), Condition: wiersz.Warunek,
		Channels: kodyKrokowZZapisu(wiersz.Kanaly), Enabled: wiersz.Czynna,
		UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}
}

// poswiadczenieKontraktu przekłada wiersz referencji na byt kontraktu. Pola na
// wartość w tym bycie nie ma i nie będzie.
func poswiadczenieKontraktu(wiersz dane.PoswiadczenieAutomatyki) shared.AutomationSecretRef {
	referencja := shared.AutomationSecretRef{
		Ref: wiersz.Odwolanie, Name: wiersz.Nazwa, ScopeId: wiersz.ZasiegID,
		UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}
	if wiersz.Zasieg != nil {
		zasieg := shared.ConfigScope(*wiersz.Zasieg)
		referencja.Scope = &zasieg
	}
	return referencja
}

// bytSejfuZOdwolania zdejmuje przedrostek `sejf:` z referencji. Sejf kluczuje
// bytem, a referencja jest bytem opatrzonym znakiem pochodzenia.
func bytSejfuZOdwolania(odwolanie string) string {
	const przedrostek = "sejf:"
	if len(odwolanie) > len(przedrostek) && odwolanie[:len(przedrostek)] == przedrostek {
		return odwolanie[len(przedrostek):]
	}
	return odwolanie
}

// errBrakSkarbca nazywa brak magazynu wartości poświadczeń. Kod `channel_unavailable`
// jak przy braku silnika kolejek: nie ma przez co wykonać, a wpięcie sejfu czyni
// żądanie wykonalnym bez zmiany treści.
var errBrakSkarbca = bladNiedostepnegoSilnika(
	"skarbiec poświadczeń nie jest wpięty — wartość nie ma gdzie spocząć poza bazą")
