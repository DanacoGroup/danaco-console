// Odpowiedzialność pliku: wyszukiwanie w projekcie (`workspace.search.project`),
// oś czasu aktywności (`workspace.activity.list`) i komentarze
// (`workspace.comment.add`, `workspace.comment.list`, `workspace.comment.delete`).
//
// ── Czego to wyszukiwanie nie zastępuje ────────────────────────────────────
// `library.file.search` przeszukuje bibliotekę centralną i wyłącznie pliki. Ta
// komenda przeszukuje JEDEN projekt i więcej niż pliki: zadania, notatki,
// ustalenia pamięci i instrukcje. Wyszukiwanie po ZNACZENIU prowadzi rodzina
// `knowledge.*` — tutaj idzie dopasowanie po słowach.
//
// Treść pliku wchodzi do wyszukiwania przez wskaźnik zbudowany komendą
// `workspace.library.text.extract`; plik bez wyciągu jest dopasowywany po samej
// nazwie. To jest różnica widoczna dla Operatora, więc trafienie z wyciągu
// niesie fragment treści, a trafienie po nazwie — nie.
package core

import (
	"context"
	"errors"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// SzukajWProjekcie obsługuje `workspace.search.project`.
func (a *adapterPrzestrzeniRoboczej) SzukajWProjekcie(ctx context.Context,
	z shared.WorkspaceSearchProjectRequest) (shared.WorkspaceSearchProjectResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceSearchProjectResponse{}, err
	}
	fraza := strings.ToLower(strings.TrimSpace(z.Query))
	if fraza == "" {
		return shared.WorkspaceSearchProjectResponse{}, bladProjektu("wyszukiwanie bez frazy")
	}
	dozwolone := map[shared.WorkspaceEntityKind]bool{}
	for _, rodzaj := range z.Kinds {
		dozwolone[rodzaj] = true
	}
	bierze := func(rodzaj shared.WorkspaceEntityKind) bool {
		return len(dozwolone) == 0 || dozwolone[rodzaj]
	}
	trafienia := []shared.WorkspaceSearchHit{}

	if bierze(shared.WorkspaceEntityKindTask) {
		zadania, err := a.repozytorium.ZadaniaWorkspace(ctx, projekt.ID)
		if err != nil {
			return shared.WorkspaceSearchProjectResponse{}, err
		}
		for _, zadanie := range zadania {
			if trafienie, jest := dopasujWorkspace(fraza, zadanie.Identyfikator,
				shared.WorkspaceEntityKindTask, zadanie.Tytul, zadanie.Opis,
				chwilaBazy(zadanie.Zaktualizowano)); jest {
				trafienia = append(trafienia, trafienie)
			}
		}
	}
	if bierze(shared.WorkspaceEntityKindNote) {
		notatki, err := a.repozytorium.NotatkiWorkspace(ctx, projekt.ID)
		if err != nil {
			return shared.WorkspaceSearchProjectResponse{}, err
		}
		for _, notatka := range notatki {
			if trafienie, jest := dopasujWorkspace(fraza, notatka.Identyfikator,
				shared.WorkspaceEntityKindNote, notatka.Tytul, notatka.Tresc,
				chwilaBazy(notatka.Zaktualizowano)); jest {
				trafienia = append(trafienia, trafienie)
			}
		}
	}
	if bierze(shared.WorkspaceEntityKindMemoryEntry) {
		wpisy, err := a.repozytorium.WpisyPamieci(ctx, projekt.ID, 0)
		if err != nil {
			return shared.WorkspaceSearchProjectResponse{}, err
		}
		for _, wpis := range wpisy {
			if trafienie, jest := dopasujWorkspace(fraza, wpis.Identyfikator,
				shared.WorkspaceEntityKindMemoryEntry, skrocOpisWorkspace(wpis.Tresc, 60),
				wpis.Tresc, chwilaBazy(wpis.Zaktualizowano)); jest {
				trafienia = append(trafienia, trafienie)
			}
		}
	}
	if bierze(shared.WorkspaceEntityKindFile) {
		wyciagi := map[string]string{}
		wykazWyciagow, err := a.repozytorium.WyciagiWorkspace(ctx, projekt.ID)
		if err != nil {
			return shared.WorkspaceSearchProjectResponse{}, err
		}
		for _, wyciag := range wykazWyciagow {
			wyciagi[wyciag.Plik] = wyciag.Tresc
		}
		for _, plik := range a.plikiProjektu(projekt.Kod, "") {
			if trafienie, jest := dopasujWorkspace(fraza, plik.Id,
				shared.WorkspaceEntityKindFile, plik.Name, wyciagi[plik.Id],
				plik.UpdatedAt); jest {
				trafienia = append(trafienia, trafienie)
			}
		}
	}
	if bierze(shared.WorkspaceEntityKindInstructions) {
		wersje, err := a.repozytorium.WersjeInstrukcjiWorkspace(ctx, projekt.ID)
		if err != nil {
			return shared.WorkspaceSearchProjectResponse{}, err
		}
		if len(wersje) > 0 {
			biezaca := wersje[0]
			if trafienie, jest := dopasujWorkspace(fraza, biezaca.Identyfikator,
				shared.WorkspaceEntityKindInstructions, "Instrukcje projektu",
				biezaca.Tresc, chwilaBazy(biezaca.Utworzono)); jest {
				trafienia = append(trafienia, trafienie)
			}
		}
	}
	if bierze(shared.WorkspaceEntityKindComment) {
		komentarze, err := a.repozytorium.KomentarzeWorkspace(ctx, projekt.ID)
		if err != nil {
			return shared.WorkspaceSearchProjectResponse{}, err
		}
		for _, komentarz := range komentarze {
			if trafienie, jest := dopasujWorkspace(fraza, komentarz.Identyfikator,
				shared.WorkspaceEntityKindComment, skrocOpisWorkspace(komentarz.Tresc, 60),
				komentarz.Tresc, chwilaBazy(komentarz.Zaktualizowano)); jest {
				trafienia = append(trafienia, trafienie)
			}
		}
	}
	sort.SliceStable(trafienia, func(i, j int) bool {
		lewy, prawy := float64(0), float64(0)
		if trafienia[i].Score != nil {
			lewy = *trafienia[i].Score
		}
		if trafienia[j].Score != nil {
			prawy = *trafienia[j].Score
		}
		return lewy > prawy
	})
	wszystkich := len(trafienia)
	return shared.WorkspaceSearchProjectResponse{
		Hits: przytnijWykazWorkspace(trafienia, z.Limit, z.Offset), Total: wszystkich,
	}, nil
}

// OsCzasu obsługuje `workspace.activity.list`.
func (a *adapterPrzestrzeniRoboczej) OsCzasu(ctx context.Context,
	z shared.WorkspaceActivityListRequest) (shared.WorkspaceActivityListResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceActivityListResponse{}, err
	}
	zdarzenia, err := a.repozytorium.ZdarzeniaWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceActivityListResponse{}, err
	}
	dozwolone := map[shared.WorkspaceActivityKind]bool{}
	for _, rodzaj := range z.Kinds {
		dozwolone[rodzaj] = true
	}
	dobrane := []shared.WorkspaceActivityEntry{}
	for _, zdarzenie := range zdarzenia {
		if len(dozwolone) > 0 && !dozwolone[zdarzenie.Rodzaj] {
			continue
		}
		if z.EntityKind != nil && *z.EntityKind != "" && zdarzenie.RodzajBytu != *z.EntityKind {
			continue
		}
		if z.EntityId != nil && *z.EntityId != "" && zdarzenie.Byt != *z.EntityId {
			continue
		}
		chwila := chwilaBazy(zdarzenie.Zaszlo)
		if z.From != nil && chwila < *z.From {
			continue
		}
		if z.To != nil && chwila > *z.To {
			continue
		}
		dobrane = append(dobrane, zdarzenieKontraktuWorkspace(projekt.Kod, zdarzenie, chwila))
	}
	wszystkich := len(dobrane)
	return shared.WorkspaceActivityListResponse{
		Entries: przytnijWykazWorkspace(dobrane, z.Limit, z.Offset), Total: wszystkich,
	}, nil
}

// DolozKomentarz obsługuje `workspace.comment.add`.
func (a *adapterPrzestrzeniRoboczej) DolozKomentarz(ctx context.Context,
	z shared.WorkspaceCommentAddRequest) (shared.WorkspaceCommentAddResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceCommentAddResponse{}, err
	}
	if strings.TrimSpace(z.Content) == "" {
		return shared.WorkspaceCommentAddResponse{}, bladProjektu("komentarz bez treści")
	}
	if strings.TrimSpace(z.TargetId) == "" {
		return shared.WorkspaceCommentAddResponse{}, bladProjektu("komentarz bez wskazania bytu")
	}
	przywolania := przywolaniaTresciWorkspace(z.Content)
	komentarz, err := a.repozytorium.ZapiszKomentarzWorkspace(ctx, dane.KomentarzWorkspace{
		ProjektID: projekt.ID, Identyfikator: nowyIdentyfikator("wscm-"),
		RodzajBytu: z.TargetKind, Byt: z.TargetId, Tresc: z.Content,
		KomentarzNadrzedny: wartoscTekstuWorkspace(z.ParentCommentId),
		Przywolania:        przywolania,
	})
	if err != nil {
		return shared.WorkspaceCommentAddResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.WorkspaceCommentAddResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindCommentAdded,
		shared.ChangeKindCreated, z.TargetKind, z.TargetId,
		"dołożono komentarz przy bycie "+z.TargetId)

	return shared.WorkspaceCommentAddResponse{
		Comment: komentarzKontraktuWorkspace(komentarz), MentionedIds: przywolania,
	}, nil
}

// Komentarze obsługuje `workspace.comment.list`.
func (a *adapterPrzestrzeniRoboczej) Komentarze(ctx context.Context,
	z shared.WorkspaceCommentListRequest) (shared.WorkspaceCommentListResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceCommentListResponse{}, err
	}
	wszystkie, err := a.repozytorium.KomentarzeWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceCommentListResponse{}, err
	}
	dobrane := []shared.WorkspaceComment{}
	for _, komentarz := range wszystkie {
		if z.TargetKind != nil && *z.TargetKind != "" && komentarz.RodzajBytu != *z.TargetKind {
			continue
		}
		if z.TargetId != nil && *z.TargetId != "" && komentarz.Byt != *z.TargetId {
			continue
		}
		dobrane = append(dobrane, komentarzKontraktuWorkspace(komentarz))
	}
	wszystkich := len(dobrane)
	return shared.WorkspaceCommentListResponse{
		Comments: przytnijWykazWorkspace(dobrane, z.Limit, z.Offset), Total: wszystkich,
	}, nil
}

// UsunKomentarz obsługuje `workspace.comment.delete`. Odpowiedzi w wątku idą
// wraz z komentarzem, chyba że żądanie mówi inaczej: odpowiedź bez pytania
// zostaje w wykazie jako zdanie bez kontekstu.
func (a *adapterPrzestrzeniRoboczej) UsunKomentarz(ctx context.Context,
	z shared.WorkspaceCommentDeleteRequest) (shared.WorkspaceCommentDeleteResponse, error) {

	komentarz, err := a.repozytorium.KomentarzWorkspace(ctx, z.CommentId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.WorkspaceCommentDeleteResponse{Deleted: false}, nil
	}
	if err != nil {
		return shared.WorkspaceCommentDeleteResponse{}, err
	}
	usuwane := []string{komentarz.Identyfikator}
	if z.WithReplies == nil || *z.WithReplies {
		wszystkie, err := a.repozytorium.KomentarzeWorkspace(ctx, komentarz.ProjektID)
		if err != nil {
			return shared.WorkspaceCommentDeleteResponse{}, err
		}
		usuwane = append(usuwane, odpowiedziWatkuWorkspace(wszystkie, komentarz.Identyfikator)...)
	}
	usuniete, err := a.repozytorium.UsunKomentarzeWorkspace(ctx, usuwane)
	if err != nil {
		return shared.WorkspaceCommentDeleteResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, komentarz.ProjektID, shared.WorkspaceActivityKindCommentAdded,
		shared.ChangeKindDeleted, komentarz.RodzajBytu, komentarz.Byt,
		"usunięto komentarz przy bycie "+komentarz.Byt)

	return shared.WorkspaceCommentDeleteResponse{
		Deleted: usuniete > 0, DeletedCommentIds: usuwane,
	}, nil
}

// odpowiedziWatkuWorkspace zbiera poddrzewo odpowiedzi komentarza.
func odpowiedziWatkuWorkspace(wszystkie []dane.KomentarzWorkspace, korzen string) []string {
	odpowiedzi := []string{}
	kolejka := []string{korzen}
	for len(kolejka) > 0 {
		biezacy := kolejka[0]
		kolejka = kolejka[1:]
		for _, komentarz := range wszystkie {
			if komentarz.KomentarzNadrzedny != biezacy || komentarz.Identyfikator == korzen {
				continue
			}
			odpowiedzi = append(odpowiedzi, komentarz.Identyfikator)
			kolejka = append(kolejka, komentarz.Identyfikator)
		}
	}
	return odpowiedzi
}

// przywolaniaTresciWorkspace wyjmuje z treści byty przywołane znakiem małpy.
func przywolaniaTresciWorkspace(tresc string) []string {
	przywolania := []string{}
	widziane := map[string]bool{}
	for _, czesc := range strings.Fields(tresc) {
		if !strings.HasPrefix(czesc, "@") || len(czesc) < 2 {
			continue
		}
		nazwa := strings.TrimRight(czesc[1:], ".,;:!?")
		if nazwa == "" || widziane[nazwa] {
			continue
		}
		widziane[nazwa] = true
		przywolania = append(przywolania, nazwa)
	}
	return przywolania
}

// dopasujWorkspace ocenia trafienie po nazwie i treści bytu. Trafienie w nazwę
// waży więcej niż trafienie w treść — Operator szuka najczęściej po nazwie.
func dopasujWorkspace(fraza, kod string, rodzaj shared.WorkspaceEntityKind,
	nazwa, tresc string, zmieniono int64) (shared.WorkspaceSearchHit, bool) {

	wNazwie := strings.Contains(strings.ToLower(nazwa), fraza)
	wTresci := strings.Contains(strings.ToLower(tresc), fraza)
	if !wNazwie && !wTresci {
		return shared.WorkspaceSearchHit{}, false
	}
	miara := 0.5
	if wNazwie {
		miara = 1.0
	}
	trafienie := shared.WorkspaceSearchHit{
		EntityId: kod, Kind: rodzaj, Title: nazwa, Score: &miara,
	}
	if zmieniono != 0 {
		chwila := zmieniono
		trafienie.UpdatedAt = &chwila
	}
	if wTresci {
		if fragment := fragmentTrafieniaWorkspace(tresc, fraza); fragment != "" {
			trafienie.Snippet = &fragment
		}
	}
	return trafienie, true
}

// fragmentTrafieniaWorkspace wycina fragment treści wokół trafienia.
func fragmentTrafieniaWorkspace(tresc, fraza string) string {
	const otoczenie = 80
	miejsce := strings.Index(strings.ToLower(tresc), fraza)
	if miejsce < 0 {
		return ""
	}
	poczatek := miejsce - otoczenie
	if poczatek < 0 {
		poczatek = 0
	}
	koniec := miejsce + len(fraza) + otoczenie
	if koniec > len(tresc) {
		koniec = len(tresc)
	}
	return strings.TrimSpace(strings.ReplaceAll(tresc[poczatek:koniec], "\n", " "))
}

// zdarzenieKontraktuWorkspace przekłada wiersz osi czasu na byt kontraktu.
func zdarzenieKontraktuWorkspace(idProjektu string, z dane.ZdarzenieWorkspace,
	chwila int64) shared.WorkspaceActivityEntry {

	zdarzenie := shared.WorkspaceActivityEntry{
		Id: z.Identyfikator, ProjectId: idProjektu, Kind: z.Rodzaj,
		Summary: z.Opis, OccurredAt: chwila,
	}
	if z.Zmiana != "" {
		zmiana := z.Zmiana
		zdarzenie.Change = &zmiana
	}
	if z.RodzajBytu != "" {
		rodzaj := z.RodzajBytu
		zdarzenie.EntityKind = &rodzaj
	}
	if z.Byt != "" {
		byt := z.Byt
		zdarzenie.EntityId = &byt
	}
	rodzajAutora := z.RodzajAutora
	zdarzenie.ActorKind = &rodzajAutora
	if z.Autor != "" {
		autor := z.Autor
		zdarzenie.ActorId = &autor
	}
	return zdarzenie
}

// komentarzKontraktuWorkspace przekłada wiersz komentarza na byt kontraktu.
func komentarzKontraktuWorkspace(k dane.KomentarzWorkspace) shared.WorkspaceComment {
	komentarz := shared.WorkspaceComment{
		Id: k.Identyfikator, ProjectId: k.ProjektKod, TargetKind: k.RodzajBytu,
		TargetId: k.Byt, Content: k.Tresc, AuthorKind: k.RodzajAutora,
		Mentions:  k.Przywolania,
		CreatedAt: chwilaBazy(k.Utworzono), UpdatedAt: chwilaBazy(k.Zaktualizowano),
	}
	if k.KomentarzNadrzedny != "" {
		nadrzedny := k.KomentarzNadrzedny
		komentarz.ParentCommentId = &nadrzedny
	}
	if k.Autor != "" {
		autor := k.Autor
		komentarz.AuthorId = &autor
	}
	return komentarz
}
