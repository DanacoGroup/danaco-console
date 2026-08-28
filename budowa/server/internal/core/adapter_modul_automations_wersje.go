// Plik obsługuje panel wersji Workflow Buildera: migawkę definicji przy każdym
// zapisie, wykaz wersji, porównanie strukturalne dwóch wersji, przywrócenie
// wcześniejszej, publikację, udostępnienie w organizacji, etykiety oraz
// przebieg próbny definicji.
package core

import (
	"context"
	"encoding/json"
	"sort"

	"danacoconsole/shared"
)

// Symulacja przeprowadza przebieg próbny definicji bez wpięcia produkcyjnego,
// sprawdzając spójność kroków, nie wywołując przy tym żadnego z ich efektów.
func (a *adapterAutomatyk) Symulacja(ctx context.Context,
	z shared.AutomationWorkflowSimulateRequest) (shared.AutomationWorkflowSimulateResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationWorkflowSimulateResponse{}, err
	}
	kroki := z.Steps
	if kroki == nil {
		if kroki, err = a.krokiKontraktu(ctx, wiersz.ID); err != nil {
			return shared.AutomationWorkflowSimulateResponse{}, bladAutomatyki(err)
		}
	}
	wyniki, zastrzezenia, powodzenie := przebiegProbny(kroki, wartoscTekstu(z.StopAtStepId))
	a.zapisAudytu(ctx, &wiersz.ID, "symulacja definicji", map[string]any{
		"krokow": len(kroki), "powodzenie": powodzenie,
	})
	odpowiedz := shared.AutomationWorkflowSimulateResponse{Results: wyniki, Succeeded: powodzenie}
	if len(zastrzezenia) > 0 {
		odpowiedz.Issues = zastrzezenia
	}
	return odpowiedz, nil
}

// przebiegProbny ocenia każdy krok po kolei i zatrzymuje się na kroku
// wskazanym. Krok bez pokrycia kończy się stanem nieudanym i zastrzeżeniem —
// to jest dokładnie to, co przebieg próbny ma znaleźć przed wpięciem
// produkcyjnym.
func przebiegProbny(kroki []shared.AutomationStep,
	krokStopu string) ([]shared.AutomationStepResult, []string, bool) {

	wyniki := make([]shared.AutomationStepResult, 0, len(kroki))
	zastrzezenia := []string{}
	powodzenie := true
	for _, krok := range kroki {
		powod := brakPokryciaKroku(krok)
		wynik := shared.AutomationStepResult{
			StepId: krok.Id, Status: shared.AutomationExecutionStatusSucceeded,
		}
		if powod != "" {
			wynik.Status = shared.AutomationExecutionStatusFailed
			komunikat := powod
			wynik.ErrorMessage = &komunikat
			zastrzezenia = append(zastrzezenia, krok.Id+": "+powod)
			powodzenie = false
		}
		wyniki = append(wyniki, wynik)
		if krokStopu != "" && krok.Id == krokStopu {
			break
		}
	}
	return wyniki, zastrzezenia, powodzenie
}

// brakPokryciaKroku nazywa to, czego krokowi brakuje do wykonania — pusty
// napis znaczy krok kompletny, gotowy do wpięcia produkcyjnego bez dalszej
// poprawki.
func brakPokryciaKroku(krok shared.AutomationStep) string {
	switch krok.Kind {
	case shared.AutomationStepKindCommand, shared.AutomationStepKindExtension:
		if wartoscTekstu(krok.Command) == "" {
			return "krok nie wskazuje komendy do wykonania"
		}
	case shared.AutomationStepKindCondition:
		if wartoscTekstu(krok.Condition) == "" {
			return "krok warunku nie niesie warunku"
		}
	case shared.AutomationStepKindSubflow:
		if len(krok.Params) == 0 {
			return "krok podprzepływu nie wskazuje automatyki wywoływanej"
		}
	}
	if krok.Id == "" {
		return "krok bez identyfikatora nie może wejść w układ zależności"
	}
	return ""
}

// WykazWersji oddaje migawki definicji od najnowszej, w liczbie ograniczonej
// wskazaniem żądania, wraz z autorem i chwilą zapisu każdej z nich.
func (a *adapterAutomatyk) WykazWersji(ctx context.Context,
	z shared.AutomationWorkflowVersionListRequest) (shared.AutomationWorkflowVersionListResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationWorkflowVersionListResponse{}, err
	}
	migawki, err := a.repozytorium.WersjeAutomatyki(ctx, wiersz.ID, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.AutomationWorkflowVersionListResponse{}, bladAutomatyki(err)
	}
	wersje := make([]shared.AutomationWorkflowVersion, 0, len(migawki))
	for _, migawka := range migawki {
		wersje = append(wersje, shared.AutomationWorkflowVersion{
			WorkflowId: wiersz.Kod, Version: migawka.Wersja,
			Steps: krokiZMigawki(migawka.Kroki), Published: migawka.Opublikowana,
			AuthorId: migawka.Autor, CreatedAt: chwilaBazy(migawka.Utworzono),
		})
	}
	return shared.AutomationWorkflowVersionListResponse{Versions: wersje}, nil
}

// PrzywrocWersje ustawia migawkę jako definicję bieżącą. Wersja zastana zostaje
// w historii: przywrócenie jest zapisem nowym, nie cofnięciem czasu.
func (a *adapterAutomatyk) PrzywrocWersje(ctx context.Context,
	z shared.AutomationWorkflowVersionRestoreRequest) (shared.AutomationWorkflowVersionRestoreResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationWorkflowVersionRestoreResponse{}, err
	}
	migawka, err := a.repozytorium.WersjaAutomatykiNumer(ctx, wiersz.ID, z.Version)
	if err != nil {
		return shared.AutomationWorkflowVersionRestoreResponse{}, bladNieznanejWersji(z.Version, err)
	}
	kroki := krokiZMigawki(migawka.Kroki)
	odpowiedz, err := a.Zapisz(ctx, shared.AutomationWorkflowSaveRequest{
		WorkflowId: &wiersz.Kod, Name: wiersz.Nazwa, Description: wiersz.Opis,
		Steps: kroki, Enabled: &wiersz.Czynna,
	})
	if err != nil {
		return shared.AutomationWorkflowVersionRestoreResponse{}, err
	}
	a.zapisAudytu(ctx, &wiersz.ID, "przywrócenie wersji definicji",
		map[string]any{"wersja": z.Version})
	return shared.AutomationWorkflowVersionRestoreResponse{Workflow: odpowiedz.Workflow}, nil
}

// PorownajWersje oddaje różnicę strukturalną dwóch migawek, porządkując
// zmiany po identyfikatorze kroku, którego dotyczą.
func (a *adapterAutomatyk) PorownajWersje(ctx context.Context,
	z shared.AutomationWorkflowVersionDiffRequest) (shared.AutomationWorkflowVersionDiffResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationWorkflowVersionDiffResponse{}, err
	}
	wyjsciowa, err := a.repozytorium.WersjaAutomatykiNumer(ctx, wiersz.ID, z.FromVersion)
	if err != nil {
		return shared.AutomationWorkflowVersionDiffResponse{}, bladNieznanejWersji(z.FromVersion, err)
	}
	porownywana, err := a.repozytorium.WersjaAutomatykiNumer(ctx, wiersz.ID, z.ToVersion)
	if err != nil {
		return shared.AutomationWorkflowVersionDiffResponse{}, bladNieznanejWersji(z.ToVersion, err)
	}
	zmiany := roznicaKrokow(krokiZMigawki(wyjsciowa.Kroki), krokiZMigawki(porownywana.Kroki))
	return shared.AutomationWorkflowVersionDiffResponse{Changes: zmiany}, nil
}

// roznicaKrokow zestawia dwa komplety kroków po ich identyfikatorach. Kroki
// porządkuje się po identyfikatorze, żeby ta sama para wersji dawała zawsze tę
// samą różnicę — wykaz zależny od kolejności mapy byłby różnicą losową.
func roznicaKrokow(wyjsciowe, porownywane []shared.AutomationStep) []shared.AutomationVersionChange {
	przedZmiana := map[string]shared.AutomationStep{}
	for _, krok := range wyjsciowe {
		przedZmiana[krok.Id] = krok
	}
	poZmianie := map[string]shared.AutomationStep{}
	for _, krok := range porownywane {
		poZmianie[krok.Id] = krok
	}
	kody := []string{}
	for kod := range przedZmiana {
		kody = append(kody, kod)
	}
	for kod := range poZmianie {
		if _, jest := przedZmiana[kod]; !jest {
			kody = append(kody, kod)
		}
	}
	sort.Strings(kody)

	zmiany := []shared.AutomationVersionChange{}
	for _, kod := range kody {
		przed, byl := przedZmiana[kod]
		po, jest := poZmianie[kod]
		switch {
		case !byl:
			zmiany = append(zmiany, shared.AutomationVersionChange{
				StepId: kod, Change: shared.ChangeKindCreated, After: zapisKroku(po),
			})
		case !jest:
			zmiany = append(zmiany, shared.AutomationVersionChange{
				StepId: kod, Change: shared.ChangeKindDeleted, Before: zapisKroku(przed),
			})
		default:
			zmiany = append(zmiany, roznicePol(kod, przed, po)...)
		}
	}
	return zmiany
}

// roznicePol nazywa pola kroku, które się rozeszły. Krok bez różnic nie oddaje
// niczego — zmiana pusta byłaby wierszem mówiącym „nic się nie stało”.
func roznicePol(kod string, przed, po shared.AutomationStep) []shared.AutomationVersionChange {
	zmiany := []shared.AutomationVersionChange{}
	for _, pole := range []struct {
		nazwa     string
		przed, po any
	}{
		{"name", przed.Name, po.Name},
		{"kind", przed.Kind, po.Kind},
		{"command", wartoscTekstu(przed.Command), wartoscTekstu(po.Command)},
		{"condition", wartoscTekstu(przed.Condition), wartoscTekstu(po.Condition)},
		{"params", string(przed.Params), string(po.Params)},
		{"dependsOn", przed.DependsOn, po.DependsOn},
	} {
		if rowneWartosci(pole.przed, pole.po) {
			continue
		}
		nazwa := pole.nazwa
		zmiany = append(zmiany, shared.AutomationVersionChange{
			StepId: kod, Change: shared.ChangeKindUpdated, Field: &nazwa,
			Before: zapisWartosci(pole.przed), After: zapisWartosci(pole.po),
		})
	}
	return zmiany
}

// rowneWartosci porównuje dwie wartości pola przez ich zapis strukturalny —
// jedyny sposób, który tak samo traktuje napis, wyliczenie i wykaz.
func rowneWartosci(pierwsza, druga any) bool {
	return string(zapisWartosci(pierwsza)) == string(zapisWartosci(druga))
}

// zapisWartosci zamienia wartość pola na zapis strukturalny kontraktu, oddając
// brak zamiast błędu, gdy wartość nie daje się zakodować.
func zapisWartosci(wartosc any) json.RawMessage {
	tresc, err := json.Marshal(wartosc)
	if err != nil {
		return nil
	}
	return json.RawMessage(tresc)
}

// zapisKroku zamienia cały krok na zapis strukturalny — nośnik pól `before`
// i `after` przy kroku dodanym albo usuniętym.
func zapisKroku(krok shared.AutomationStep) json.RawMessage {
	return zapisWartosci(krok)
}

// UstawEtykiety podmienia komplet etykiet automatyki na wykaz przekazany
// w żądaniu i odnotowuje zmianę w dzienniku audytu.
func (a *adapterAutomatyk) UstawEtykiety(ctx context.Context,
	z shared.AutomationWorkflowTagSetRequest) (shared.AutomationWorkflowTagSetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationWorkflowTagSetResponse{}, err
	}
	if err := a.repozytorium.UstawEtykietyAutomatyki(ctx, wiersz.ID, z.Tags); err != nil {
		return shared.AutomationWorkflowTagSetResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, &wiersz.ID, "zmiana etykiet automatyki", map[string]any{"etykiety": z.Tags})
	automatyka, err := a.zloz(ctx, wiersz)
	if err != nil {
		return shared.AutomationWorkflowTagSetResponse{}, bladAutomatyki(err)
	}
	return shared.AutomationWorkflowTagSetResponse{Workflow: automatyka}, nil
}

// Opublikuj rozdziela wersję roboczą od opublikowanej. Brak wskazania publikuje
// wersję bieżącą — tak mówi kontrakt.
func (a *adapterAutomatyk) Opublikuj(ctx context.Context,
	z shared.AutomationWorkflowPublishRequest) (shared.AutomationWorkflowPublishResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationWorkflowPublishResponse{}, err
	}
	numer := wartoscLiczbyLub(z.Version, wiersz.Wersja)
	if _, err := a.repozytorium.WersjaAutomatykiNumer(ctx, wiersz.ID, numer); err != nil {
		return shared.AutomationWorkflowPublishResponse{}, bladNieznanejWersji(numer, err)
	}
	if err := a.repozytorium.OpublikujWersjeAutomatyki(ctx, wiersz.ID, numer); err != nil {
		return shared.AutomationWorkflowPublishResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, &wiersz.ID, "publikacja wersji definicji", map[string]any{"wersja": numer})
	automatyka, err := a.automatykaPoZmianie(ctx, wiersz.Kod)
	if err != nil {
		return shared.AutomationWorkflowPublishResponse{}, err
	}
	return shared.AutomationWorkflowPublishResponse{Workflow: automatyka, PublishedVersion: numer}, nil
}

// Udostepnij przełącza udostępnienie automatyki jako komponentu własnego
// w obrębie organizacji i odnotowuje zmianę w dzienniku audytu.
func (a *adapterAutomatyk) Udostepnij(ctx context.Context,
	z shared.AutomationWorkflowShareRequest) (shared.AutomationWorkflowShareResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationWorkflowShareResponse{}, err
	}
	if err := a.repozytorium.UstawUdostepnienieAutomatyki(ctx, wiersz.ID, z.Shared); err != nil {
		return shared.AutomationWorkflowShareResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, &wiersz.ID, "zmiana udostępnienia automatyki",
		map[string]any{"udostepniona": z.Shared})
	automatyka, err := a.automatykaPoZmianie(ctx, wiersz.Kod)
	if err != nil {
		return shared.AutomationWorkflowShareResponse{}, err
	}
	return shared.AutomationWorkflowShareResponse{Workflow: automatyka}, nil
}

// automatykaPoZmianie oddaje automatykę odczytaną NA NOWO. Zmiany publikacji
// i udostępnienia idą osobnym zapisem kolumn, więc wiersz zastany sprzed
// zmiany niósłby stan sprzed niej.
func (a *adapterAutomatyk) automatykaPoZmianie(ctx context.Context,
	kod string) (shared.AutomationWorkflow, error) {

	wiersz, err := a.wiersz(ctx, kod)
	if err != nil {
		return shared.AutomationWorkflow{}, err
	}
	automatyka, err := a.zloz(ctx, wiersz)
	if err != nil {
		return shared.AutomationWorkflow{}, bladAutomatyki(err)
	}
	return automatyka, nil
}

// krokiZMigawki odczytuje kroki z zapisu strukturalnego wersji. Zapis
// nieczytelny daje wykaz pusty, nie awarię: historia ma się wyświetlić,
// a wersja uszkodzona ma być widoczna jako pusta, nie jako brak całego wykazu.
func krokiZMigawki(zapis string) []shared.AutomationStep {
	kroki := []shared.AutomationStep{}
	if zapis == "" {
		return kroki
	}
	if err := json.Unmarshal([]byte(zapis), &kroki); err != nil {
		return []shared.AutomationStep{}
	}
	return kroki
}

// bladNieznanejWersji odróżnia wersję definicji, której magazyn nie zna, od
// usterki samego odczytu tej wersji.
func bladNieznanejWersji(numer int, err error) error {
	return bladNieznanegoBytuAutomatyki(err, "wersja definicji nie istnieje: ", numer)
}
