// Odpowiedzialność pliku: moduł Studio — to, co ktoś o dokumencie POWIEDZIAŁ
// (komentarze redakcyjne, adnotacje przy różnicach) i to, co w nim ZMIENIŁ
// (śledzenie zmian, decyzja o propozycji).
package core

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	przedrostekKomentarzaStudia = "studio-kom-"
	przedrostekAdnotacjiStudia  = "studio-adn-"
	rodzajKomentarza            = "komentarz"
	rodzajAdnotacji             = "adnotacja"
)

// DodajKomentarz obsługuje `studio.comment.add`. Autor komentarza pochodzi
// z żądania i z gniazda uruchomieniowego, nie z zaszytej wartości domyślnej.
func (a *adapterStudia) DodajKomentarz(ctx context.Context,
	z shared.StudioCommentAddRequest) (shared.StudioCommentAddResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioCommentAddResponse{}, err
	}
	if z.Body == "" {
		return shared.StudioCommentAddResponse{}, bladWskazaniaStudio("komentarz bez treści")
	}

	// Autor zaszyty jako Operator przemilczałby komentarze modelu w widoku.
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	})
	wiersz := dane.KomentarzStudia{
		Kod:               nowyIdentyfikator(przedrostekKomentarzaStudia),
		Rodzaj:            rodzajKomentarza,
		WersjaKod:         z.VersionId,
		WatekNadrzednyKod: z.ParentCommentId,
		Autor:             string(wykonawca.Rodzaj),
		Tresc:             z.Body,
	}
	wiersz.ZakresOd = liczbaZeWskaznika(z.SelectionStart)
	wiersz.ZakresDo = liczbaZeWskaznika(z.SelectionEnd)

	zapisany, err := a.repozytorium.ZapiszKomentarz(ctx, dokument.ID, wiersz)
	if err != nil {
		return shared.StudioCommentAddResponse{}, bladStudio(err)
	}
	return shared.StudioCommentAddResponse{Comment: złóżKomentarz(zapisany)}, nil
}

// Komentarze obsługuje `studio.comment.list`.
//
// Wątki rozwiązane są domyślnie POZA wykazem: kontrakt mówi „brak znaczy nie",
// a redakcja zamknięta nie ma zaśmiecać marginesu. Operator sięga po nie jawnie.
func (a *adapterStudia) Komentarze(ctx context.Context,
	z shared.StudioCommentListRequest) (shared.StudioCommentListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioCommentListResponse{}, err
	}
	wiersze, err := a.repozytorium.Komentarze(ctx, dokument.ID, rodzajKomentarza)
	if err != nil {
		return shared.StudioCommentListResponse{}, bladStudio(err)
	}
	zRozwiazanymi := z.IncludeResolved != nil && *z.IncludeResolved
	wersja := wartoscTekstu(z.VersionId)

	komentarze := []shared.StudioComment{}
	for _, wiersz := range wiersze {
		if wiersz.Rozwiazany && !zRozwiazanymi {
			continue
		}
		if wersja != "" && wartoscTekstu(wiersz.WersjaKod) != wersja {
			continue
		}
		komentarze = append(komentarze, złóżKomentarz(wiersz))
	}
	return shared.StudioCommentListResponse{Comments: komentarze}, nil
}

// RozstrzygnijKomentarz obsługuje `studio.comment.resolve`: zapisuje stan
// rozwiązania wątku komentarza i oddaje jego postać po zapisie.
func (a *adapterStudia) RozstrzygnijKomentarz(ctx context.Context,
	z shared.StudioCommentResolveRequest) (shared.StudioCommentResolveResponse, error) {

	if z.CommentId == "" {
		return shared.StudioCommentResolveResponse{},
			bladWskazaniaStudio("rozstrzygnięcie bez wskazania komentarza")
	}
	zapisany, err := a.repozytorium.RozstrzygnijKomentarz(ctx, z.CommentId, z.Resolved)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.StudioCommentResolveResponse{},
				bladBrakuStudio("komentarz nie istnieje: " + z.CommentId)
		}
		return shared.StudioCommentResolveResponse{}, bladStudio(err)
	}
	return shared.StudioCommentResolveResponse{Comment: złóżKomentarz(zapisany)}, nil
}

// DodajAdnotacje obsługuje `studio.annotation.add`. Autor pochodzi z żądania,
// tą samą zasadą co przy komentarzu: adnotacja modelu przy fragmencie różnicy
// ma być podpisana jako `model`.
func (a *adapterStudia) DodajAdnotacje(ctx context.Context,
	z shared.StudioAnnotationAddRequest) (shared.StudioAnnotationAddResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioAnnotationAddResponse{}, err
	}
	if z.Body == "" {
		return shared.StudioAnnotationAddResponse{}, bladWskazaniaStudio("adnotacja bez treści")
	}
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	})
	fragment := int64(z.HunkIndex)
	wiersz := dane.KomentarzStudia{
		Kod:                 nowyIdentyfikator(przedrostekAdnotacjiStudia),
		Rodzaj:              rodzajAdnotacji,
		Autor:               string(wykonawca.Rodzaj),
		FragmentNumer:       &fragment,
		WersjaOdniesieniaID: z.BaseVersionId,
		WersjaPorownywanaID: z.TargetVersionId,
		PropozycjaID:        z.ProposalId,
		Tresc:               z.Body,
	}
	zapisany, err := a.repozytorium.ZapiszKomentarz(ctx, dokument.ID, wiersz)
	if err != nil {
		return shared.StudioAnnotationAddResponse{}, bladStudio(err)
	}
	return shared.StudioAnnotationAddResponse{Annotation: złóżAdnotacje(zapisany)}, nil
}

// Adnotacje obsługuje `studio.annotation.list`. Para wersji zawęża wykaz, bo
// adnotacja opisuje różnicę, a nie dokument.
func (a *adapterStudia) Adnotacje(ctx context.Context,
	z shared.StudioAnnotationListRequest) (shared.StudioAnnotationListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioAnnotationListResponse{}, err
	}
	wiersze, err := a.repozytorium.Komentarze(ctx, dokument.ID, rodzajAdnotacji)
	if err != nil {
		return shared.StudioAnnotationListResponse{}, bladStudio(err)
	}
	odniesienie, porownywana := wartoscTekstu(z.BaseVersionId), wartoscTekstu(z.TargetVersionId)

	adnotacje := []shared.StudioAnnotation{}
	for _, wiersz := range wiersze {
		if odniesienie != "" && wartoscTekstu(wiersz.WersjaOdniesieniaID) != odniesienie {
			continue
		}
		if porownywana != "" && wartoscTekstu(wiersz.WersjaPorownywanaID) != porownywana {
			continue
		}
		adnotacje = append(adnotacje, złóżAdnotacje(wiersz))
	}
	return shared.StudioAnnotationListResponse{Annotations: adnotacje}, nil
}

// UstawSledzenie obsługuje `studio.tracking.set`: przełącza rejestrowanie
// zmian w dokumencie i oddaje stan czynny po zapisie, odczytany z bazy.
func (a *adapterStudia) UstawSledzenie(ctx context.Context,
	z shared.StudioTrackingSetRequest) (shared.StudioTrackingSetResponse, error) {

	if z.DocumentId == "" {
		return shared.StudioTrackingSetResponse{},
			bladWskazaniaStudio("przestawienie śledzenia bez wskazania dokumentu")
	}
	if err := a.repozytorium.UstawSledzenie(ctx, z.DocumentId, z.Enabled); err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.StudioTrackingSetResponse{},
				bladBrakuStudio("dokument nie istnieje: " + z.DocumentId)
		}
		return shared.StudioTrackingSetResponse{}, bladStudio(err)
	}
	czynne, err := a.repozytorium.Sledzenie(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTrackingSetResponse{}, bladStudio(err)
	}
	return shared.StudioTrackingSetResponse{Enabled: czynne}, nil
}

// ZmianySledzone obsługuje `studio.tracking.list`: oddaje wykaz zmian
// śledzonych dokumentu, opcjonalnie zawężony do samych oczekujących.
func (a *adapterStudia) ZmianySledzone(ctx context.Context,
	z shared.StudioTrackingListRequest) (shared.StudioTrackingListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTrackingListResponse{}, err
	}
	wiersze, err := a.repozytorium.ZmianySledzone(ctx, dokument.ID)
	if err != nil {
		return shared.StudioTrackingListResponse{}, bladStudio(err)
	}
	tylkoOczekujace := z.PendingOnly != nil && *z.PendingOnly

	zmiany := []shared.StudioTrackedChange{}
	for _, wiersz := range wiersze {
		if tylkoOczekujace && wiersz.Decyzja != "oczekuje" {
			continue
		}
		zmiany = append(zmiany, złóżZmianeSledzona(wiersz))
	}
	return shared.StudioTrackingListResponse{Changes: zmiany}, nil
}

// RozstrzygnijZmiany obsługuje `studio.tracking.decide`. Decyzja zmienia treść
// dokumentu, nie sam znacznik: przyjęcie wstawienia zostawia tekst wstawiony,
// odrzucenie usuwa go.
func (a *adapterStudia) RozstrzygnijZmiany(ctx context.Context,
	z shared.StudioTrackingDecideRequest) (shared.StudioTrackingDecideResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTrackingDecideResponse{}, err
	}
	if len(z.ChangeIds) == 0 {
		return shared.StudioTrackingDecideResponse{},
			bladWskazaniaStudio("decyzja bez wskazania zmian")
	}

	wszystkie, err := a.repozytorium.ZmianySledzone(ctx, dokument.ID)
	if err != nil {
		return shared.StudioTrackingDecideResponse{}, bladStudio(err)
	}
	objete := map[string]bool{}
	for _, kod := range z.ChangeIds {
		objete[kod] = true
	}

	wybrane := make([]dane.ZmianaSledzona, 0, len(z.ChangeIds))
	for _, zmiana := range wszystkie {
		if objete[zmiana.Kod] && zmiana.Decyzja == "oczekuje" {
			wybrane = append(wybrane, zmiana)
		}
	}
	// Od końca — zakresy liczone są w treści sprzed decyzji.
	sort.Slice(wybrane, func(i, j int) bool { return wybrane[i].ZakresOd > wybrane[j].ZakresOd })

	decyzja := "odrzucona"
	if z.Accept {
		decyzja = "przyjeta"
	}
	// Autor wersji bierze się z gniazda, które decyzję podjęło.
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{})

	tresc := wartoscTekstu(dokument.Tresc)
	rozstrzygniete := 0
	for _, zmiana := range wybrane {
		nowa, zmieniono := zastosujZmianeSledzona(tresc, zmiana, z.Accept)
		if zmieniono {
			tresc = nowa
		}
		zapisana, err := a.repozytorium.RozstrzygnijZmianeSledzona(ctx, zmiana.Kod, decyzja)
		if err != nil {
			return shared.StudioTrackingDecideResponse{}, bladStudio(err)
		}
		if zapisana {
			rozstrzygniete++
		}
	}

	dokument.Tresc = &tresc
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return shared.StudioTrackingDecideResponse{}, bladStudio(err)
	}
	if z.CreateVersion == nil || *z.CreateVersion {
		zapisany, err = a.zalozWersjeDokumentu(ctx, zapisany, tresc, wykonawca.Rodzaj, nil)
		if err != nil {
			return shared.StudioTrackingDecideResponse{}, err
		}
	}

	return shared.StudioTrackingDecideResponse{
		Document: a.zlozDokument(zapisany),
		Decided:  rozstrzygniete,
	}, nil
}

// RozstrzygnijPropozycje obsługuje `studio.proposal.decide`. Decyzja zapada
// po stronie rdzenia: propozycja przestaje być nierozstrzygnięta także dla
// innych urządzeń konta.
func (a *adapterStudia) RozstrzygnijPropozycje(ctx context.Context,
	z shared.StudioProposalDecideRequest) (shared.StudioProposalDecideResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioProposalDecideResponse{}, err
	}
	if z.ProposalId == "" {
		return shared.StudioProposalDecideResponse{},
			bladWskazaniaStudio("decyzja bez wskazania propozycji")
	}
	propozycja, err := a.repozytorium.Propozycja(ctx, z.ProposalId)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.StudioProposalDecideResponse{},
				bladBrakuStudio("propozycja nie istnieje: " + z.ProposalId)
		}
		return shared.StudioProposalDecideResponse{}, bladStudio(err)
	}

	if !z.Accept {
		return shared.StudioProposalDecideResponse{Document: a.zlozDokument(dokument)}, nil
	}

	tresc := wartoscTekstu(propozycja.TrescWyniku)
	if tresc == "" {
		tresc = wartoscTekstu(propozycja.Tresc)
	}
	if tresc == "" {
		return shared.StudioProposalDecideResponse{},
			bladWskazaniaStudio("propozycja " + z.ProposalId + " nie niesie treści — " +
				"przyjęcie zostawiłoby dokument bez zmiany, więc nie jest przyjęciem")
	}

	if len(z.HunkIndexes) > 0 {
		wybrane, err := zlozTrescZFragmentow(wartoscTekstu(dokument.Tresc), tresc, z.HunkIndexes)
		if err != nil {
			return shared.StudioProposalDecideResponse{}, err
		}
		tresc = wybrane
	}

	dokument.Tresc = &tresc
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return shared.StudioProposalDecideResponse{}, bladStudio(err)
	}

	odpowiedz := shared.StudioProposalDecideResponse{}
	if z.CreateVersion == nil || *z.CreateVersion {
		zapisany, err = a.zalozWersjeDokumentu(ctx, zapisany, tresc,
			shared.StudioAuthorModel, &z.ProposalId)
		if err != nil {
			return shared.StudioProposalDecideResponse{}, err
		}
		if zapisany.WersjaBiezacaKod != nil {
			wersja, err := a.repozytorium.Wersja(ctx, *zapisany.WersjaBiezacaKod)
			if err == nil {
				zlozona := a.zlozWersje(wersja)
				odpowiedz.Version = &zlozona
			}
		}
	}
	odpowiedz.Document = a.zlozDokument(zapisany)
	return odpowiedz, nil
}

// zlozTrescZFragmentow składa treść wybiórczo: fragmenty wskazane przez
// Operatora bierze ze strony propozycji, pozostałe zostawia tak, jak stoją
// w dokumencie.
func zlozTrescZFragmentow(bazowa, docelowa string, numery []int) (string, error) {
	fragmenty := policzFragmentyRoznicy(bazowa, docelowa)
	istniejace := map[int]bool{}
	for _, fragment := range fragmenty {
		istniejace[fragment.Index] = true
	}
	wskazane := map[int]bool{}
	for _, numer := range numery {
		if !istniejace[numer] {
			return "", bladWskazaniaStudio("różnica nie ma fragmentu o numerze " +
				strconv.Itoa(numer) + " — fragmentów jest " + strconv.Itoa(len(fragmenty)))
		}
		wskazane[numer] = true
	}

	czesci := []string{}
	for _, fragment := range fragmenty {
		strona := fragment.Before
		if fragment.Kind != shared.DiffHunkKindContext && wskazane[fragment.Index] {
			strona = fragment.After
		}
		if strona == nil {
			continue
		}
		czesci = append(czesci, *strona)
	}
	return strings.Join(czesci, "\n"), nil
}

// zastosujZmianeSledzona nakłada decyzję na treść i mówi, czy treść się
// zmieniła. Zakres liczy się w znakach, tak jak nazywa go kontrakt.
func zastosujZmianeSledzona(tresc string, zmiana dane.ZmianaSledzona, przyjeta bool) (string, bool) {
	znaki := []rune(tresc)
	od, do := int(zmiana.ZakresOd), int(zmiana.ZakresDo)
	if od < 0 || do > len(znaki) || od > do {
		// Zakres spoza treści oznacza dokument zmieniony po rejestracji zmiany.
		return tresc, false
	}
	docelowa := wartoscTekstu(zmiana.TrescPo)
	if !przyjeta {
		docelowa = wartoscTekstu(zmiana.TrescPrzed)
	}
	if zmiana.Rodzaj == "usuniecie" {
		docelowa = ""
		if !przyjeta {
			docelowa = wartoscTekstu(zmiana.TrescPrzed)
		}
	}
	return string(znaki[:od]) + docelowa + string(znaki[do:]), true
}

// zalozWersjeDokumentu zakłada wersję i czyni ją bieżącą. Wydzielone, bo trzy
// rodziny robią dokładnie to samo w trzech krokach: wersja, wskaźnik wersji
// bieżącej, ponowny zapis dokumentu.
func (a *adapterStudia) zalozWersjeDokumentu(ctx context.Context, dokument dane.DokumentStudia,
	tresc string, autor shared.StudioAuthor, propozycja *string) (dane.DokumentStudia, error) {

	nazwaAutora := string(autor)
	wersja, err := a.repozytorium.ZapiszWersje(ctx, dokument.ID, dane.WersjaDokumentu{
		Kod:           nowyIdentyfikator(przedrostekWersjiStudio),
		Tresc:         &tresc,
		Autor:         &nazwaAutora,
		PropozycjaKod: propozycja,
	})
	if err != nil {
		return dokument, bladStudio(err)
	}
	dokument.WersjaBiezacaKod = &wersja.Kod
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return dokument, bladStudio(err)
	}
	return zapisany, nil
}

// dokumentDoCzynnosci odczytuje dokument wskazany żądaniem — jedna droga
// sprawdzenia wskazania i jedna treść odmowy dla całej rodziny.
func (a *adapterStudia) dokumentDoCzynnosci(ctx context.Context,
	kod string) (dane.DokumentStudia, error) {

	if kod == "" {
		return dane.DokumentStudia{}, bladWskazaniaStudio("czynność bez wskazania dokumentu")
	}
	dokument, err := a.repozytorium.Dokument(ctx, kod)
	if err != nil {
		return dane.DokumentStudia{}, bladNieznanegoDokumentu(kod, err)
	}
	return dokument, nil
}

func złóżKomentarz(wiersz dane.KomentarzStudia) shared.StudioComment {
	komentarz := shared.StudioComment{
		Id:              wiersz.Kod,
		DocumentId:      wiersz.DokumentKod,
		VersionId:       wiersz.WersjaKod,
		ParentCommentId: wiersz.WatekNadrzednyKod,
		Author:          shared.StudioAuthor(wiersz.Autor),
		Body:            wiersz.Tresc,
		Resolved:        wiersz.Rozwiazany,
		CreatedAt:       chwilaBazy(wiersz.Utworzono),
	}
	if wiersz.ZakresOd != nil {
		od := int(*wiersz.ZakresOd)
		komentarz.SelectionStart = &od
	}
	if wiersz.ZakresDo != nil {
		do := int(*wiersz.ZakresDo)
		komentarz.SelectionEnd = &do
	}
	return komentarz
}

func złóżAdnotacje(wiersz dane.KomentarzStudia) shared.StudioAnnotation {
	adnotacja := shared.StudioAnnotation{
		Id:              wiersz.Kod,
		DocumentId:      wiersz.DokumentKod,
		BaseVersionId:   wiersz.WersjaOdniesieniaID,
		TargetVersionId: wiersz.WersjaPorownywanaID,
		ProposalId:      wiersz.PropozycjaID,
		Author:          shared.StudioAuthor(wiersz.Autor),
		Body:            wiersz.Tresc,
		CreatedAt:       chwilaBazy(wiersz.Utworzono),
	}
	if wiersz.FragmentNumer != nil {
		adnotacja.HunkIndex = int(*wiersz.FragmentNumer)
	}
	return adnotacja
}

func złóżZmianeSledzona(wiersz dane.ZmianaSledzona) shared.StudioTrackedChange {
	return shared.StudioTrackedChange{
		Id:         wiersz.Kod,
		DocumentId: wiersz.DokumentKod,
		Kind:       shared.StudioChangeKind(wiersz.Rodzaj),
		Author:     shared.StudioAuthor(wiersz.Autor),
		RangeStart: int(wiersz.ZakresOd),
		RangeEnd:   int(wiersz.ZakresDo),
		Before:     wiersz.TrescPrzed,
		After:      wiersz.TrescPo,
		Decision:   shared.StudioChangeDecision(wiersz.Decyzja),
		CreatedAt:  chwilaBazy(wiersz.Utworzono),
	}
}

// liczbaZeWskaznika przekłada wskaźnik na liczbę całkowitą bazy danych;
// wskaźnik pusty zostaje wskaźnikiem pustym.
func liczbaZeWskaznika(wartosc *int) *int64 {
	if wartosc == nil {
		return nil
	}
	liczba := int64(*wartosc)
	return &liczba
}
