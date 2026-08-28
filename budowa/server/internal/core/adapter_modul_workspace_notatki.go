// Odpowiedzialność pliku: notatki i strony wiki projektu, ich zapis, odczyt,
// wykaz, usunięcie, drzewo stron oraz odnośniki wsteczne liczone przy zapisie
// strony, nie przy jej odczycie.
package core

import (
	"context"
	"errors"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ZapiszNotatke obsługuje `workspace.note.save`: zapisuje treść strony, jej
// tytuł i miejsce w drzewie stron, a przy okazji przelicza odnośniki wsteczne.
func (a *adapterPrzestrzeniRoboczej) ZapiszNotatke(ctx context.Context,
	z shared.WorkspaceNoteSaveRequest) (shared.WorkspaceNoteSaveResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceNoteSaveResponse{}, err
	}
	if strings.TrimSpace(z.Title) == "" {
		return shared.WorkspaceNoteSaveResponse{}, bladProjektu("notatka bez tytułu")
	}
	identyfikator := wartoscTekstuWorkspace(z.NoteId)
	zmiana := shared.ChangeKind(shared.ChangeKindUpdated)
	if identyfikator == "" {
		identyfikator, zmiana = nowyIdentyfikator("wsnt-"), shared.ChangeKindCreated
	}
	istniejace, err := a.repozytorium.NotatkiWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceNoteSaveResponse{}, err
	}
	nazwy := nazwyOdnosnikowWorkspace(z.Content)
	odnosniki, brakujace := odnosnikiNotatkiWorkspace(projekt.ID, identyfikator, z.Content,
		nazwy, istniejace)

	notatka := dane.NotatkaWorkspace{
		ProjektID: projekt.ID, Identyfikator: identyfikator,
		Tytul: strings.TrimSpace(z.Title), Tresc: z.Content,
		NotatkaNadrzedna: wartoscTekstuWorkspace(z.ParentNoteId),
		Etykiety:         z.Tags, Naglowki: naglowkiTresciWorkspace(z.Content),
	}
	zapisana, err := a.repozytorium.ZapiszNotatkeWorkspace(ctx, notatka, odnosniki)
	if err != nil {
		return shared.WorkspaceNoteSaveResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.WorkspaceNoteSaveResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindNoteChanged, zmiana,
		shared.WorkspaceEntityKindNote, zapisana.Identyfikator,
		"zapisano notatkę „"+zapisana.Tytul+"”")

	return shared.WorkspaceNoteSaveResponse{
		Note:         notatkaKontraktuWorkspace(zapisana, istniejace),
		LinkedNames:  nazwy,
		MissingNames: brakujace,
	}, nil
}

// Notatka obsługuje `workspace.note.get`: oddaje jedną stronę wiki wraz z jej
// nagłówkami i ścieżką w drzewie stron projektu.
func (a *adapterPrzestrzeniRoboczej) Notatka(ctx context.Context,
	z shared.WorkspaceNoteGetRequest) (shared.WorkspaceNoteGetResponse, error) {

	notatka, err := a.notatkaWskazanaWorkspace(ctx, z.NoteId)
	if err != nil {
		return shared.WorkspaceNoteGetResponse{}, err
	}
	rodzenstwo, err := a.repozytorium.NotatkiWorkspace(ctx, notatka.ProjektID)
	if err != nil {
		return shared.WorkspaceNoteGetResponse{}, err
	}
	return shared.WorkspaceNoteGetResponse{Note: notatkaKontraktuWorkspace(notatka, rodzenstwo)}, nil
}

// Notatki obsługuje `workspace.note.list`: wykazuje strony wiki projektu bez
// pełnej treści, do zbudowania listy albo panelu nawigacji.
func (a *adapterPrzestrzeniRoboczej) Notatki(ctx context.Context,
	z shared.WorkspaceNoteListRequest) (shared.WorkspaceNoteListResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceNoteListResponse{}, err
	}
	wszystkie, err := a.repozytorium.NotatkiWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceNoteListResponse{}, err
	}
	fraza := strings.ToLower(strings.TrimSpace(wartoscTekstuWorkspace(z.Query)))
	dobrane := []shared.WorkspaceNote{}
	for _, notatka := range wszystkie {
		if fraza != "" && !strings.Contains(strings.ToLower(notatka.Tytul), fraza) {
			continue
		}
		if z.Tag != nil && *z.Tag != "" && !zawieraTekstWorkspace(notatka.Etykiety, *z.Tag) {
			continue
		}
		if z.ParentNoteId != nil && *z.ParentNoteId != "" &&
			notatka.NotatkaNadrzedna != *z.ParentNoteId {
			continue
		}
		dobrane = append(dobrane, notatkaKontraktuWorkspace(notatka, wszystkie))
	}
	wszystkich := len(dobrane)
	return shared.WorkspaceNoteListResponse{
		Notes: przytnijWykazWorkspace(dobrane, z.Limit, z.Offset), Total: wszystkich,
	}, nil
}

// UsunNotatke obsługuje `workspace.note.delete`.
//
// Bez `withChildren` strony podrzędne przechodzą pod stronę nadrzędną
// usuwanej — wiki nie gubi wtedy gałęzi razem z jej korzeniem.
func (a *adapterPrzestrzeniRoboczej) UsunNotatke(ctx context.Context,
	z shared.WorkspaceNoteDeleteRequest) (shared.WorkspaceNoteDeleteResponse, error) {

	notatka, err := a.repozytorium.NotatkaWorkspace(ctx, z.NoteId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.WorkspaceNoteDeleteResponse{Deleted: false}, nil
	}
	if err != nil {
		return shared.WorkspaceNoteDeleteResponse{}, err
	}
	wszystkie, err := a.repozytorium.NotatkiWorkspace(ctx, notatka.ProjektID)
	if err != nil {
		return shared.WorkspaceNoteDeleteResponse{}, err
	}
	usuwane := []string{notatka.Identyfikator}
	if z.WithChildren != nil && *z.WithChildren {
		usuwane = append(usuwane, podrzedneNotatkiWorkspace(wszystkie, notatka.Identyfikator)...)
	}
	if _, err := a.repozytorium.UsunNotatkiWorkspace(ctx, usuwane, notatka.NotatkaNadrzedna); err != nil {
		return shared.WorkspaceNoteDeleteResponse{}, err
	}
	odnosniki, err := a.repozytorium.OdnosnikiWorkspace(ctx, notatka.ProjektID)
	if err != nil {
		return shared.WorkspaceNoteDeleteResponse{}, err
	}
	osierocone := 0
	for _, odnosnik := range odnosniki {
		if odnosnik.NotatkaDocelowa == "" && strings.EqualFold(odnosnik.NazwaDocelowa, notatka.Tytul) {
			osierocone++
		}
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, notatka.ProjektID); err != nil {
		return shared.WorkspaceNoteDeleteResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, notatka.ProjektID, shared.WorkspaceActivityKindNoteChanged,
		shared.ChangeKindDeleted, shared.WorkspaceEntityKindNote, notatka.Identyfikator,
		"usunięto notatkę „"+notatka.Tytul+"”")

	return shared.WorkspaceNoteDeleteResponse{
		Deleted: true, DeletedNoteIds: usuwane, OrphanedBacklinkCount: &osierocone,
	}, nil
}

// DrzewoNotatek obsługuje `workspace.note.tree.get`: składa strony wiki
// projektu w drzewo według relacji strona nadrzędna i strona podrzędna.
func (a *adapterPrzestrzeniRoboczej) DrzewoNotatek(ctx context.Context,
	z shared.WorkspaceNoteTreeGetRequest) (shared.WorkspaceNoteTreeGetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceNoteTreeGetResponse{}, err
	}
	notatki, err := a.repozytorium.NotatkiWorkspace(ctx, projekt.ID)
	if err != nil {
		return shared.WorkspaceNoteTreeGetResponse{}, err
	}
	podrzedne := map[string][]dane.NotatkaWorkspace{}
	znane := map[string]bool{}
	for _, notatka := range notatki {
		znane[notatka.Identyfikator] = true
	}
	for _, notatka := range notatki {
		rodzic := notatka.NotatkaNadrzedna
		// Strona bez istniejącego rodzica wchodzi do korzenia, gałąź osierocona
		// pozostaje widoczna.
		if rodzic != "" && !znane[rodzic] {
			rodzic = ""
		}
		podrzedne[rodzic] = append(podrzedne[rodzic], notatka)
	}
	wezly := []shared.WorkspaceNoteNode{}
	var obejdz func(rodzic string)
	obejdz = func(rodzic string) {
		dzieci := podrzedne[rodzic]
		sort.SliceStable(dzieci, func(i, j int) bool { return dzieci[i].Tytul < dzieci[j].Tytul })
		for numer, notatka := range dzieci {
			wezel := shared.WorkspaceNoteNode{
				NoteId: notatka.Identyfikator, Title: notatka.Tytul, Order: numer,
				ChildCount: len(podrzedne[notatka.Identyfikator]),
			}
			if notatka.NotatkaNadrzedna != "" && znane[notatka.NotatkaNadrzedna] {
				nadrzedna := notatka.NotatkaNadrzedna
				wezel.ParentNoteId = &nadrzedna
			}
			wezly = append(wezly, wezel)
			obejdz(notatka.Identyfikator)
		}
	}
	obejdz("")
	return shared.WorkspaceNoteTreeGetResponse{Nodes: wezly}, nil
}

// OdnosnikiWsteczne obsługuje `workspace.note.backlink.list`: wykazuje strony
// wskazujące daną stronę odnośnikiem, wraz z fragmentem treści wokół odnośnika.
func (a *adapterPrzestrzeniRoboczej) OdnosnikiWsteczne(ctx context.Context,
	z shared.WorkspaceNoteBacklinkListRequest) (shared.WorkspaceNoteBacklinkListResponse, error) {

	notatka, err := a.notatkaWskazanaWorkspace(ctx, z.NoteId)
	if err != nil {
		return shared.WorkspaceNoteBacklinkListResponse{}, err
	}
	odnosniki, err := a.repozytorium.OdnosnikiWorkspace(ctx, notatka.ProjektID)
	if err != nil {
		return shared.WorkspaceNoteBacklinkListResponse{}, err
	}
	notatki, err := a.repozytorium.NotatkiWorkspace(ctx, notatka.ProjektID)
	if err != nil {
		return shared.WorkspaceNoteBacklinkListResponse{}, err
	}
	tytuly := map[string]string{}
	for _, pozycja := range notatki {
		tytuly[pozycja.Identyfikator] = pozycja.Tytul
	}
	wykaz := []shared.WorkspaceBacklink{}
	for _, odnosnik := range odnosniki {
		wskazuje := odnosnik.NotatkaDocelowa == notatka.Identyfikator ||
			(odnosnik.NotatkaDocelowa == "" && strings.EqualFold(odnosnik.NazwaDocelowa, notatka.Tytul))
		if !wskazuje || odnosnik.NotatkaZrodlowa == notatka.Identyfikator {
			continue
		}
		wsteczny := shared.WorkspaceBacklink{
			SourceNoteId: odnosnik.NotatkaZrodlowa,
			SourceTitle:  tytuly[odnosnik.NotatkaZrodlowa],
			TargetName:   odnosnik.NazwaDocelowa,
			Kind:         odnosnik.Rodzaj,
		}
		if odnosnik.NotatkaDocelowa != "" {
			cel := odnosnik.NotatkaDocelowa
			wsteczny.TargetNoteId = &cel
		}
		if odnosnik.Kontekst != "" {
			kontekst := odnosnik.Kontekst
			wsteczny.Context = &kontekst
		}
		wykaz = append(wykaz, wsteczny)
	}
	wszystkich := len(wykaz)
	return shared.WorkspaceNoteBacklinkListResponse{
		Backlinks: przytnijWykazWorkspace(wykaz, z.Limit, z.Offset), Total: wszystkich,
	}, nil
}

// notatkaWskazanaWorkspace odnajduje notatkę wskazaną żądaniem, zgłaszając
// odmowę, gdy identyfikator nie wskazuje istniejącej strony projektu.
func (a *adapterPrzestrzeniRoboczej) notatkaWskazanaWorkspace(ctx context.Context,
	identyfikator string) (dane.NotatkaWorkspace, error) {

	if strings.TrimSpace(identyfikator) == "" {
		return dane.NotatkaWorkspace{}, bladProjektu("wywołanie bez wskazania notatki")
	}
	notatka, err := a.repozytorium.NotatkaWorkspace(ctx, identyfikator)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.NotatkaWorkspace{}, bladProjektu("notatki " + identyfikator +
			" nie ma w projekcie")
	}
	return notatka, err
}

// podrzedneNotatkiWorkspace zbiera poddrzewo stron podrzędnych wskazanej
// strony, przechodząc rekurencyjnie kolejne poziomy relacji rodzic i dziecko.
func podrzedneNotatkiWorkspace(wszystkie []dane.NotatkaWorkspace, korzen string) []string {
	podrzedne := []string{}
	kolejka := []string{korzen}
	for len(kolejka) > 0 {
		biezaca := kolejka[0]
		kolejka = kolejka[1:]
		for _, notatka := range wszystkie {
			if notatka.NotatkaNadrzedna != biezaca || notatka.Identyfikator == korzen {
				continue
			}
			podrzedne = append(podrzedne, notatka.Identyfikator)
			kolejka = append(kolejka, notatka.Identyfikator)
		}
	}
	return podrzedne
}

// nazwyOdnosnikowWorkspace wyjmuje z treści nazwy wskazane zapisem `[[nazwa]]`,
// każdą raz i w kolejności wystąpienia.
func nazwyOdnosnikowWorkspace(tresc string) []string {
	nazwy := []string{}
	widziane := map[string]bool{}
	reszta := tresc
	for {
		otwarcie := strings.Index(reszta, "[[")
		if otwarcie < 0 {
			break
		}
		reszta = reszta[otwarcie+2:]
		zamkniecie := strings.Index(reszta, "]]")
		if zamkniecie < 0 {
			break
		}
		nazwa := strings.TrimSpace(reszta[:zamkniecie])
		reszta = reszta[zamkniecie+2:]
		if nazwa == "" || widziane[strings.ToLower(nazwa)] {
			continue
		}
		widziane[strings.ToLower(nazwa)] = true
		nazwy = append(nazwy, nazwa)
	}
	return nazwy
}

// odnosnikiNotatkiWorkspace wiąże nazwy z treści ze stronami projektu. Drugi
// wynik niesie nazwy bez strony — to nie usterka, lecz wiki w budowie.
func odnosnikiNotatkiWorkspace(projektID int64, zrodlo, tresc string, nazwy []string,
	strony []dane.NotatkaWorkspace) ([]dane.OdnosnikWorkspace, []string) {

	wedlugTytulu := map[string]string{}
	for _, strona := range strony {
		wedlugTytulu[strings.ToLower(strona.Tytul)] = strona.Identyfikator
	}
	odnosniki := make([]dane.OdnosnikWorkspace, 0, len(nazwy))
	brakujace := []string{}
	for _, nazwa := range nazwy {
		cel := wedlugTytulu[strings.ToLower(nazwa)]
		if cel == "" {
			brakujace = append(brakujace, nazwa)
		}
		odnosniki = append(odnosniki, dane.OdnosnikWorkspace{
			ProjektID: projektID, NotatkaZrodlowa: zrodlo, NazwaDocelowa: nazwa,
			NotatkaDocelowa: cel, Rodzaj: shared.WorkspaceGraphEdgeKindWikilink,
			Kontekst: kontekstOdnosnikaWorkspace(tresc, nazwa),
		})
	}
	return odnosniki, brakujace
}

// kontekstOdnosnikaWorkspace wycina fragment treści wokół odnośnika — panel
// „co linkuje tutaj" pokazuje zdanie, w którym odnośnik stoi, a nie samą nazwę
// strony linkującej.
func kontekstOdnosnikaWorkspace(tresc, nazwa string) string {
	const otoczenie = 60
	miejsce := strings.Index(tresc, "[["+nazwa)
	if miejsce < 0 {
		return ""
	}
	poczatek := miejsce - otoczenie
	if poczatek < 0 {
		poczatek = 0
	}
	koniec := miejsce + len(nazwa) + otoczenie
	if koniec > len(tresc) {
		koniec = len(tresc)
	}
	return strings.TrimSpace(strings.ReplaceAll(tresc[poczatek:koniec], "\n", " "))
}

// naglowkiTresciWorkspace zbiera nagłówki Markdowna — materiał spisu treści
// notatki. Zbiera je zapis, a nie odczyt, bo spis czyta się przy każdym
// otwarciu strony, a treść zmienia rzadziej.
func naglowkiTresciWorkspace(tresc string) []string {
	naglowki := []string{}
	wBloku := false
	for _, wiersz := range strings.Split(tresc, "\n") {
		przyciety := strings.TrimSpace(wiersz)
		if strings.HasPrefix(przyciety, "```") {
			wBloku = !wBloku
			continue
		}
		if wBloku || !strings.HasPrefix(przyciety, "#") {
			continue
		}
		naglowek := strings.TrimSpace(strings.TrimLeft(przyciety, "#"))
		if naglowek != "" {
			naglowki = append(naglowki, naglowek)
		}
	}
	return naglowki
}

// notatkaKontraktuWorkspace przekłada wiersz notatki na byt kontraktu wraz ze
// ścieżką w drzewie stron złożoną z tytułów.
func notatkaKontraktuWorkspace(n dane.NotatkaWorkspace,
	wszystkie []dane.NotatkaWorkspace) shared.WorkspaceNote {

	notatka := shared.WorkspaceNote{
		Id: n.Identyfikator, ProjectId: n.ProjektKod, Title: n.Tytul, Content: n.Tresc,
		Tags: n.Etykiety, Headings: n.Naglowki,
		CreatedAt: chwilaBazy(n.Utworzono), UpdatedAt: chwilaBazy(n.Zaktualizowano),
	}
	if n.NotatkaNadrzedna != "" {
		nadrzedna := n.NotatkaNadrzedna
		notatka.ParentNoteId = &nadrzedna
	}
	rodzaj := n.RodzajAutora
	notatka.AuthorKind = &rodzaj
	if n.Autor != "" {
		autor := n.Autor
		notatka.AuthorId = &autor
	}
	if sciezka := sciezkaStronyWorkspace(n, wszystkie); sciezka != "" {
		notatka.Path = &sciezka
	}
	return notatka
}

// sciezkaStronyWorkspace składa ścieżkę strony z tytułów jej przodków, od
// korzenia drzewa stron do samej strony.
func sciezkaStronyWorkspace(n dane.NotatkaWorkspace, wszystkie []dane.NotatkaWorkspace) string {
	wedlugKodu := map[string]dane.NotatkaWorkspace{}
	for _, notatka := range wszystkie {
		wedlugKodu[notatka.Identyfikator] = notatka
	}
	czlony := []string{n.Tytul}
	rodzic := n.NotatkaNadrzedna
	for glebokosc := 0; rodzic != "" && glebokosc <= len(wszystkie); glebokosc++ {
		nadrzedna, jest := wedlugKodu[rodzic]
		if !jest {
			break
		}
		czlony = append([]string{nadrzedna.Tytul}, czlony...)
		rodzic = nadrzedna.NotatkaNadrzedna
	}
	return strings.Join(czlony, " / ")
}
