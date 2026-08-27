// Odpowiedzialność pliku: JEDEN PRZEŁĄCZNIK pokaż wszystko, co zrobił
// model — wykaz zmian wykonawcy, skakanie po nich i cofnięcie
// Z ZACHOWANIEM zmian Operatora.
package core

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZmianyModelu obsługuje studio.model.changes.list, oddając wykaz
// wszystkich zmian wykonawcy wraz z licznikiem otwartych.
func (a *adapterStudia) ZmianyModelu(ctx context.Context,
	z shared.StudioModelChangesListRequest) (shared.StudioModelChangesListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioModelChangesListResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioModelChangesListResponse{}, err
	}

	zRozstrzygnietymi := z.IncludeDecided != nil && *z.IncludeDecided
	tylkoPostac := z.OnlyFormChanges != nil && *z.OnlyFormChanges

	wybrane, err := a.zmianyModeluWybierz(ctx, skladnica, dokument.ID,
		zRozstrzygnietymi, tylkoPostac, z.AgentId, z.SubagentId)
	if err != nil {
		return shared.StudioModelChangesListResponse{}, err
	}

	zestawienie := shared.StudioModelChangeSummary{
		DocumentId: dokument.Kod,
		Changes:    []shared.StudioTrackedChange{},
		Actions:    []shared.StudioDocumentAction{},
	}
	liczniki := map[string]*shared.StudioAgentChangeCount{}
	kolejnosc := []string{}
	for _, zmiana := range wybrane {
		zestawienie.Changes = append(zestawienie.Changes, zmianyModeluZloz(zmiana))
		zestawienie.Total++
		if zmiana.Rodzaj == shared.StudioChangeKindFormatowanie {
			zestawienie.FormChanges++
		} else {
			zestawienie.ContentChanges++
		}
		if zmiana.Decyzja == "oczekuje" {
			zestawienie.OpenChanges++
		}

		klucz := wartoscTekstu(zmiana.AutorAgentKod) + "\x00" + wartoscTekstu(zmiana.AutorPodagentKod)
		pozycja, jest := liczniki[klucz]
		if !jest {
			zero := 0
			pozycja = &shared.StudioAgentChangeCount{
				Actor:       zmianyModeluAktor(zmiana),
				OpenChanges: &zero,
			}
			liczniki[klucz] = pozycja
			kolejnosc = append(kolejnosc, klucz)
		}
		pozycja.Total++
		if zmiana.Rodzaj == shared.StudioChangeKindFormatowanie {
			pozycja.FormChanges++
		} else {
			pozycja.ContentChanges++
		}
		if zmiana.Decyzja == "oczekuje" {
			*pozycja.OpenChanges++
		}
	}
	for _, klucz := range kolejnosc {
		zestawienie.ByAgent = append(zestawienie.ByAgent, *liczniki[klucz])
	}

	// Czynności dziennika autora model — po nich widać zmiany POSTACI bez
	// śladu w literach.
	czynnosci, err := skladnica.CzynnosciDokumentu(ctx, dokument.ID)
	if err != nil {
		return shared.StudioModelChangesListResponse{}, bladStudio(err)
	}
	for _, czynnosc := range czynnosci {
		if czynnosc.AutorRodzaj != string(shared.StudioAuthorModel) {
			continue
		}
		if !zRozstrzygnietymi && czynnosc.Stan != string(shared.StudioActionStateActive) {
			continue
		}
		if kontrolaTekstNiepusty(z.AgentId) &&
			wartoscTekstu(czynnosc.AutorAgentKod) != *z.AgentId {

			continue
		}
		if kontrolaTekstNiepusty(z.SubagentId) &&
			wartoscTekstu(czynnosc.AutorPodagentKod) != *z.SubagentId {

			continue
		}
		zestawienie.Actions = append(zestawienie.Actions, dziennikZlozCzynnosc(czynnosc))
	}

	// Spięcia wykonawców o ten sam fragment: Operator ma widzieć, czyja
	// zmiana została odłożona.
	spiecia, err := skladnica.SpieciaWykonawcow(ctx, dokument.ID)
	if err != nil {
		return shared.StudioModelChangesListResponse{}, bladStudio(err)
	}
	for _, spiecie := range spiecia {
		zestawienie.Conflicts = append(zestawienie.Conflicts, agenciZlozSpiecie(spiecie))
	}
	return shared.StudioModelChangesListResponse{Summary: zestawienie}, nil
}

// PrzeskocDoZmianyModelu obsługuje studio.model.changes.navigate; kolejność
// jest kolejnością W TREŚCI, nie w czasie.
func (a *adapterStudia) PrzeskocDoZmianyModelu(ctx context.Context,
	z shared.StudioModelChangesNavigateRequest) (shared.StudioModelChangesNavigateResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioModelChangesNavigateResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioModelChangesNavigateResponse{}, err
	}
	kierunek := strings.TrimSpace(strings.ToLower(z.Direction))
	if kierunek != "next" && kierunek != "previous" {
		return shared.StudioModelChangesNavigateResponse{}, bladWskazaniaStudio(
			"kierunek skoku „" + z.Direction + "” nie jest ani „next”, ani „previous”")
	}
	tylkoPostac := z.OnlyFormChanges != nil && *z.OnlyFormChanges
	wybrane, err := a.zmianyModeluWybierz(ctx, skladnica, dokument.ID, false, tylkoPostac,
		z.AgentId, z.SubagentId)
	if err != nil {
		return shared.StudioModelChangesNavigateResponse{}, err
	}
	odpowiedz := shared.StudioModelChangesNavigateResponse{Total: len(wybrane)}
	if len(wybrane) == 0 {
		// Wykaz pusty NIE jest brakiem funkcji: to prawdziwa odpowiedź, gdy
		// model nie pracował w dokumencie.
		return odpowiedz, nil
	}
	od := 0
	if z.FromOffset != nil {
		od = *z.FromOffset
	}
	if kierunek == "next" {
		for numer, zmiana := range wybrane {
			if int(zmiana.ZakresOd) > od {
				zlozona := zmianyModeluZloz(zmiana)
				wskazanie := numer + 1
				odpowiedz.Change, odpowiedz.Index = &zlozona, &wskazanie
				return odpowiedz, nil
			}
		}
		return odpowiedz, nil
	}
	for numer := len(wybrane) - 1; numer >= 0; numer-- {
		if int(wybrane[numer].ZakresOd) < od {
			zlozona := zmianyModeluZloz(wybrane[numer])
			wskazanie := numer + 1
			odpowiedz.Change, odpowiedz.Index = &zlozona, &wskazanie
			return odpowiedz, nil
		}
	}
	return odpowiedz, nil
}

// CofnijZmianyModelu obsługuje studio.model.changes.revert, cofając
// zmiany wykonawcy Z ZACHOWANIEM zmian Operatora.
func (a *adapterStudia) CofnijZmianyModelu(ctx context.Context,
	z shared.StudioModelChangesRevertRequest) (shared.StudioModelChangesRevertResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioModelChangesRevertResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioModelChangesRevertResponse{}, err
	}
	wszystko := z.All != nil && *z.All
	if !wszystko && len(z.ChangeIds) == 0 {
		return shared.StudioModelChangesRevertResponse{}, bladWskazaniaStudio(
			"cofnięcie zmian modelu bez wskazania zmian i bez „all” — cofnięcie " +
				"wszystkiego jest decyzją Operatora i musi być wypowiedziane wprost")
	}

	wybrane, err := a.zmianyModeluWybierz(ctx, skladnica, stan.dokument.ID, false, false,
		z.AgentId, z.SubagentId)
	if err != nil {
		return shared.StudioModelChangesRevertResponse{}, err
	}
	if !wszystko {
		objete := map[string]bool{}
		for _, kod := range z.ChangeIds {
			objete[strings.TrimSpace(kod)] = true
		}
		zawezone := make([]dane.ZmianaWykonawcyStudia, 0, len(objete))
		znalezione := map[string]bool{}
		for _, zmiana := range wybrane {
			if objete[zmiana.Kod] {
				zawezone = append(zawezone, zmiana)
				znalezione[zmiana.Kod] = true
			}
		}
		for kod := range objete {
			if !znalezione[kod] {
				return shared.StudioModelChangesRevertResponse{}, bladBrakuStudio(
					"dokument " + stan.dokument.Kod + " nie ma nierozstrzygniętej zmiany " +
						"modelu " + kod)
			}
		}
		wybrane = zawezone
	}
	if len(wybrane) == 0 {
		return shared.StudioModelChangesRevertResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict, "moduł Studio: w dokumencie "+stan.dokument.Kod+
				" nie stoi ani jedna nierozstrzygnięta zmiana modelu — nie ma czego cofnąć"))
	}

	// Kopia przed czynnością nieodwracalną, wymagana dla przyjęcia
	// i cofnięcia wszystkich zmian modelu.
	odpowiedz := shared.StudioModelChangesRevertResponse{}
	if z.CreateBackup == nil || *z.CreateBackup {
		kopia, err := a.kopiaPrzedCzynnoscia(ctx, stan.dokument,
			shared.StudioBackupReasonBeforeIrreversible)
		if err != nil {
			return shared.StudioModelChangesRevertResponse{}, err
		}
		odpowiedz.BackupId = &kopia.Kod
	}

	// Od końca dokumentu — zakresy liczone są w treści sprzed cofnięcia.
	sort.Slice(wybrane, func(i, j int) bool { return wybrane[i].ZakresOd > wybrane[j].ZakresOd })

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	czynnosciDoCofniecia := map[string]bool{}
	for _, zmiana := range wybrane {
		if zmiana.Rodzaj == shared.StudioChangeKindFormatowanie {
			// Zmiana POSTACI nie ma czego zdjąć z liter — cofa się ją wpisem
			// dziennika sprzed czynności.
			if zmiana.CzynnoscKod != nil && *zmiana.CzynnoscKod != "" {
				czynnosciDoCofniecia[*zmiana.CzynnoscKod] = true
			} else {
				bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
					Reason: "zmiana postaci bez wpisu dziennika",
					Detail: kontrolaWskaznikTekstu("zmiana " + zmiana.Kod + " dotyczy postaci, " +
						"a nie wskazuje wpisu dziennika, w którym stoi drzewo sprzed czynności " +
						"— bez niego cofnięcie wymagałoby zgadywania, co model przestawił"),
				})
			}
		} else {
			od, do := int(zmiana.ZakresOd), int(zmiana.ZakresDo)
			przed := wartoscTekstu(zmiana.TrescPrzed)
			postacRozetnij(&stan.forma, od, do)
			postacZamienTresc(&stan.forma, od, do, przed, nil, nil)
			bilans.Applied++
		}
		if _, err := a.repozytorium.RozstrzygnijZmianeSledzona(ctx, zmiana.Kod,
			"odrzucona"); err != nil {

			return shared.StudioModelChangesRevertResponse{}, bladStudio(err)
		}
	}
	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioModelChangesRevertResponse{}, err
	}

	// Czynności postaci cofa dziennik — tą samą drogą, którą cofa je
	// Operator pojedynczo.
	if len(czynnosciDoCofniecia) > 0 {
		kody := make([]string, 0, len(czynnosciDoCofniecia))
		for kod := range czynnosciDoCofniecia {
			kody = append(kody, kod)
		}
		nieprawda := false
		cofniete, err := a.CofnijCzynnosc(ctx, shared.StudioJournalRevertRequest{
			DocumentId: stan.dokument.Kod, ActionIds: kody, CreateVersion: &nieprawda,
		})
		if err != nil {
			// Odmowa dziennika NIE przewraca cofnięcia: zmiany treści już
			// weszły, milczenie byłoby nieprawdą.
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "czynności postaci nie dały się cofnąć",
				Detail: kontrolaWskaznikTekstu(err.Error()),
			})
		} else {
			bilans.Applied += cofniete.Balance.Applied
			bilans.Skipped = append(bilans.Skipped, cofniete.Balance.Skipped...)
			if stan, err = a.postacWczytaj(ctx, z.DocumentId); err != nil {
				return shared.StudioModelChangesRevertResponse{}, err
			}
		}
	}

	// Ile zmian OPERATORA zachowano — miara, że cofnięcie nie było
	// przywróceniem wersji.
	zachowane := 0
	wszystkieZmiany, err := skladnica.ZmianyWykonawcow(ctx, stan.dokument.ID)
	if err != nil {
		return shared.StudioModelChangesRevertResponse{}, bladStudio(err)
	}
	for _, zmiana := range wszystkieZmiany {
		if zmiana.Autor == string(shared.StudioAuthorUzytkownik) && zmiana.Decyzja == "oczekuje" {
			zachowane++
		}
	}

	dokument := stan.dokument
	if z.CreateVersion == nil || *z.CreateVersion {
		if dokument, err = a.zalozWersjeDokumentu(ctx, dokument, wartoscTekstu(dokument.Tresc),
			shared.StudioAuthorUzytkownik, nil); err != nil {

			return shared.StudioModelChangesRevertResponse{}, err
		}
	}
	bilans.SkippedCount = len(bilans.Skipped)
	bilans.Note = kontrolaWskaznikTekstu("Cofnięto " + strconv.Itoa(bilans.Applied) +
		" zmian(y) wykonawcy; zmian Operatora oczekujących w dokumencie zostało " +
		strconv.Itoa(zachowane) + " — cofnięcie nie było przywróceniem wersji.")
	return shared.StudioModelChangesRevertResponse{
		RevertedCount:       bilans.Applied,
		KeptOperatorChanges: zachowane,
		Document:            a.zlozDokument(dokument),
		Form:                stan.forma,
		Balance:             bilans,
		BackupId:            odpowiedz.BackupId,
	}, nil
}

// ── Wspólne ─────────────────────────────────────────────────────────────────

// zmianyModeluWybierz przesiewa zmiany śledzone dokumentu do zmian
// WYKONAWCY; jedna droga wyboru dla wykazu, skakania i cofania.
func (a *adapterStudia) zmianyModeluWybierz(ctx context.Context, skladnica KontrolaPracyStudia,
	dokumentID int64, zRozstrzygnietymi, tylkoPostac bool,
	agentKod, podagentKod *string) ([]dane.ZmianaWykonawcyStudia, error) {

	wiersze, err := skladnica.ZmianyWykonawcow(ctx, dokumentID)
	if err != nil {
		return nil, bladStudio(err)
	}
	wybrane := make([]dane.ZmianaWykonawcyStudia, 0, len(wiersze))
	for _, zmiana := range wiersze {
		if zmiana.Autor != string(shared.StudioAuthorModel) {
			continue
		}
		if !zRozstrzygnietymi && zmiana.Decyzja != "oczekuje" {
			continue
		}
		if tylkoPostac && zmiana.Rodzaj != shared.StudioChangeKindFormatowanie {
			continue
		}
		if kontrolaTekstNiepusty(agentKod) && wartoscTekstu(zmiana.AutorAgentKod) != *agentKod {
			continue
		}
		if kontrolaTekstNiepusty(podagentKod) &&
			wartoscTekstu(zmiana.AutorPodagentKod) != *podagentKod {

			continue
		}
		wybrane = append(wybrane, zmiana)
	}
	return wybrane, nil
}

// zmianyModeluZloz składa zmianę śledzoną kontraktu z wiersza niosącego
// tożsamość wykonawcy, gotową do odpowiedzi.
func zmianyModeluZloz(wiersz dane.ZmianaWykonawcyStudia) shared.StudioTrackedChange {
	return shared.StudioTrackedChange{
		Id:                 wiersz.Kod,
		DocumentId:         wiersz.DokumentKod,
		Kind:               shared.StudioChangeKind(wiersz.Rodzaj),
		Author:             shared.StudioAuthor(wiersz.Autor),
		RangeStart:         int(wiersz.ZakresOd),
		RangeEnd:           int(wiersz.ZakresDo),
		Before:             wiersz.TrescPrzed,
		After:              wiersz.TrescPo,
		Decision:           shared.StudioChangeDecision(wiersz.Decyzja),
		CreatedAt:          chwilaBazy(wiersz.Utworzono),
		AuthorAgentId:      wiersz.AutorAgentKod,
		AuthorAgentName:    wiersz.AutorAgentNazwa,
		AuthorAgentVersion: wiersz.AutorAgentWersja,
		AuthorSubagentId:   wiersz.AutorPodagentKod,
	}
}

// zmianyModeluAktor składa tożsamość wykonawcy z wiersza zmiany, rozbijając
// ją na kod agenta i podagenta.
func zmianyModeluAktor(wiersz dane.ZmianaWykonawcyStudia) shared.StudioActor {
	return shared.StudioActor{
		Kind:         shared.StudioAuthor(wiersz.Autor),
		AgentId:      wiersz.AutorAgentKod,
		AgentName:    wiersz.AutorAgentNazwa,
		AgentVersion: wiersz.AutorAgentWersja,
		SubagentId:   wiersz.AutorPodagentKod,
	}
}

// zmianyModeluStempluj dopisuje do zmiany śledzonej tożsamość wykonawcy
// i wpis dziennika, którym da się ją cofnąć pojedynczo.
func (a *adapterStudia) zmianyModeluStempluj(ctx context.Context,
	zmiana *shared.StudioTrackedChange, wykonawca kontrolaWykonawca, czynnoscKod *string) error {

	if zmiana == nil {
		return nil
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return err
	}
	if err := skladnica.StemplujTozsamoscZmiany(ctx, zmiana.Id, wykonawca.AgentKod,
		wykonawca.AgentNazwa, wykonawca.AgentWersja, wykonawca.PodagentKod,
		nil, nil, czynnoscKod); err != nil {

		return bladStudio(err)
	}
	zmiana.AuthorAgentId = wykonawca.AgentKod
	zmiana.AuthorAgentName = wykonawca.AgentNazwa
	zmiana.AuthorAgentVersion = wykonawca.AgentWersja
	zmiana.AuthorSubagentId = wykonawca.PodagentKod
	return nil
}
