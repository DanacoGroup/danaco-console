// Odpowiedzialność pliku: gałęzie dokumentu Studia (rodzina F4–F5 opracowania)
// oraz odwołanie do wersji (J5) — `studio.branch.create`, `studio.branch.list`,
// `studio.branch.merge`, `studio.version.reference.create`. Metody dopisują się
// na `adapterStudia` zadeklarowanym w `adapter_modul_studio.go`.
//
// ── Gałąź nie jest kopią dokumentu ──────────────────────────────────────────
// Rozgałęzienie zakłada WIERSZ gałęzi i JEDNĄ wersję startową na niej, a nie
// drugi dokument. Drugi dokument oderwałby wariant od historii, z której
// wyrósł: nie dałoby się już powiedzieć, od czego wariant odszedł, a scalanie
// nie miałoby wspólnego przodka, na którym stoi cała ta rodzina.
//
// ── Scalanie oddaje konflikt, a nie rozstrzyga go domysłem ──────────────────
// Trójstronne porównanie (przodek, gałąź scalana, gałąź docelowa) rozstrzyga
// samo tam, gdzie zmieniła jedna strona. Tam, gdzie zmieniły obie i zmieniły
// inaczej, rdzeń NIE wybiera. Wybór „bierzemy docelową" byłby cichym
// skasowaniem cudzej redakcji — dlatego konflikt wraca kontraktem i czeka na
// rozstrzygnięcie Operatora.
package core

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZalozGalaz obsługuje `studio.branch.create`.
//
// Wersja startowa zostaje punktem odniesienia gałęzi na stałe — to ona będzie
// wspólnym przodkiem przy scalaniu. Obok niej powstaje pierwsza wersja gałęzi
// z treścią przepisaną z punktu startowego: bez niej gałąź świeżo założona nie
// miałaby czoła, a `document.save` dopisywałby wersje do pnia zamiast do niej.
func (a *adapterStudia) ZalozGalaz(ctx context.Context,
	z shared.StudioBranchCreateRequest) (shared.StudioBranchCreateResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" || strings.TrimSpace(z.FromVersionId) == "" ||
		strings.TrimSpace(z.Name) == "" {
		return shared.StudioBranchCreateResponse{}, bladWskazaniaStudio(
			"założenie gałęzi wymaga dokumentu, wersji startowej i nazwy")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioBranchCreateResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	startowa, err := a.repozytorium.Wersja(ctx, z.FromVersionId)
	if err != nil {
		return shared.StudioBranchCreateResponse{}, bladWskazaniaStronyPorownania(z.FromVersionId, err)
	}
	if startowa.DokumentKod != dokument.Kod {
		return shared.StudioBranchCreateResponse{}, bladWskazaniaStudio(
			"wersja " + z.FromVersionId + " należy do dokumentu " + startowa.DokumentKod +
				", a gałąź zakłada się w dokumencie " + dokument.Kod)
	}

	kodGalezi := nowyIdentyfikator(przedrostekGaleziStudia)
	tresc, err := a.trescZOdwolania(startowa.Tresc, startowa.TrescOdwolanie)
	if err != nil {
		return shared.StudioBranchCreateResponse{}, err
	}
	czolo, err := a.repozytorium.ZapiszWersje(ctx, dokument.ID, dane.WersjaDokumentu{
		Kod:          nowyIdentyfikator(przedrostekWersjiStudio),
		Etykieta:     wskaznikTekstu("gałąź " + z.Name),
		Podsumowanie: wskaznikTekstu("punkt startowy gałęzi: wersja " + startowa.Kod),
		Tresc:        &tresc,
		GalazKod:     &kodGalezi,
	})
	if err != nil {
		return shared.StudioBranchCreateResponse{}, bladStudio(err)
	}

	galaz, err := a.repozytorium.ZapiszGalaz(ctx, dokument.ID, dane.GalazStudia{
		Kod:              kodGalezi,
		Nazwa:            strings.TrimSpace(z.Name),
		WersjaStartowaID: startowa.Kod,
		WersjaBiezacaID:  &czolo.Kod,
	})
	if err != nil {
		return shared.StudioBranchCreateResponse{}, bladStudio(err)
	}

	// Gałąź świeżo założona staje się treścią bieżącą edytora — tak opisuje to
	// przycisk „Rozgałęź" w opracowaniu (rozdz. 3.6): tworzy gałąź I OTWIERA ją
	// jako bieżącą treść. Gałąź, która powstaje niewidoczna, byłaby gałęzią,
	// o której Operator dowiaduje się dopiero z wykazu.
	dokument.Tresc = &tresc
	dokument.TrescOdwolanie = nil
	dokument.WersjaBiezacaKod = &czolo.Kod
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return shared.StudioBranchCreateResponse{}, bladStudio(err)
	}

	return shared.StudioBranchCreateResponse{
		Branch:   zlozGalazStudia(galaz),
		Document: a.zlozDokument(zapisany),
	}, nil
}

// Galezie obsługuje `studio.branch.list`.
func (a *adapterStudia) Galezie(ctx context.Context,
	z shared.StudioBranchListRequest) (shared.StudioBranchListResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" {
		return shared.StudioBranchListResponse{}, bladWskazaniaStudio(
			"wykaz gałęzi bez wskazania dokumentu")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioBranchListResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	wiersze, err := a.repozytorium.Galezie(ctx, dokument.ID)
	if err != nil {
		return shared.StudioBranchListResponse{}, bladStudio(err)
	}
	galezie := make([]shared.StudioBranch, 0, len(wiersze))
	for _, wiersz := range wiersze {
		galezie = append(galezie, zlozGalazStudia(wiersz))
	}
	return shared.StudioBranchListResponse{Branches: galezie}, nil
}

// ScalGalezie obsługuje `studio.branch.merge`.
//
// Wspólnym przodkiem jest wersja startowa gałęzi SCALANEJ: od niej wariant
// odszedł, więc to ona mówi, co w każdej z dwóch treści jest zmianą, a co
// stanem zastanym. Bez przodka porównanie dwóch czół dałoby konflikt na każdym
// fragmencie, który zmieniła tylko jedna strona.
func (a *adapterStudia) ScalGalezie(ctx context.Context,
	z shared.StudioBranchMergeRequest) (shared.StudioBranchMergeResponse, error) {

	if strings.TrimSpace(z.SourceBranchId) == "" || strings.TrimSpace(z.TargetBranchId) == "" {
		return shared.StudioBranchMergeResponse{}, bladWskazaniaStudio(
			"scalenie wymaga gałęzi scalanej i gałęzi docelowej")
	}
	if z.SourceBranchId == z.TargetBranchId {
		return shared.StudioBranchMergeResponse{}, bladWskazaniaStudio(
			"gałąź scalana i docelowa to ta sama gałąź " + z.SourceBranchId)
	}
	scalana, err := a.galazStudia(ctx, z.SourceBranchId)
	if err != nil {
		return shared.StudioBranchMergeResponse{}, err
	}
	docelowa, err := a.galazStudia(ctx, z.TargetBranchId)
	if err != nil {
		return shared.StudioBranchMergeResponse{}, err
	}
	if scalana.DokumentKod != docelowa.DokumentKod {
		return shared.StudioBranchMergeResponse{}, bladWskazaniaStudio(
			"gałęzie należą do różnych dokumentów: " + scalana.DokumentKod + " i " + docelowa.DokumentKod)
	}
	dokument, err := a.repozytorium.Dokument(ctx, scalana.DokumentKod)
	if err != nil {
		return shared.StudioBranchMergeResponse{}, bladNieznanegoDokumentu(scalana.DokumentKod, err)
	}

	przodek, err := a.trescWersjiPoKodzie(ctx, scalana.WersjaStartowaID)
	if err != nil {
		return shared.StudioBranchMergeResponse{}, err
	}
	trescScalanej, err := a.trescCzolaGalezi(ctx, scalana)
	if err != nil {
		return shared.StudioBranchMergeResponse{}, err
	}
	trescDocelowej, err := a.trescCzolaGalezi(ctx, docelowa)
	if err != nil {
		return shared.StudioBranchMergeResponse{}, err
	}

	scalenie := scalTrojstronnieStudia(przodek, trescScalanej, trescDocelowej, z.Resolutions)
	if len(scalenie.konflikty) > 0 {
		return shared.StudioBranchMergeResponse{Merged: false, Conflicts: scalenie.konflikty}, nil
	}

	wynik := scalenie.tresc
	nowa, err := a.repozytorium.ZapiszWersje(ctx, dokument.ID, dane.WersjaDokumentu{
		Kod:          nowyIdentyfikator(przedrostekWersjiStudio),
		Etykieta:     wskaznikTekstu("scalenie gałęzi " + scalana.Nazwa + " do " + docelowa.Nazwa),
		Podsumowanie: wskaznikTekstu("scalono gałąź " + scalana.Kod + " do " + docelowa.Kod),
		Tresc:        &wynik,
		GalazKod:     &docelowa.Kod,
	})
	if err != nil {
		return shared.StudioBranchMergeResponse{}, bladStudio(err)
	}

	docelowa.WersjaBiezacaID = &nowa.Kod
	if _, err := a.repozytorium.ZapiszGalaz(ctx, dokument.ID, docelowa); err != nil {
		return shared.StudioBranchMergeResponse{}, bladStudio(err)
	}
	scalana.Scalona = true
	if _, err := a.repozytorium.ZapiszGalaz(ctx, dokument.ID, scalana); err != nil {
		return shared.StudioBranchMergeResponse{}, bladStudio(err)
	}

	dokument.Tresc = &wynik
	dokument.TrescOdwolanie = nil
	dokument.WersjaBiezacaKod = &nowa.Kod
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return shared.StudioBranchMergeResponse{}, bladStudio(err)
	}
	zlozony := a.zlozDokument(zapisany)
	return shared.StudioBranchMergeResponse{Merged: true, Document: &zlozony}, nil
}

// UtworzOdwolanieWersji obsługuje `studio.version.reference.create`.
//
// Odwołanie jest adresem W OBRĘBIE PLATFORMY, nie odnośnikiem sieciowym: rdzeń
// nie wystawia treści dokumentu na zewnątrz i nie ma jak zapewnić, że adres
// wyprowadzony na świat byłby czytelny wyłącznie dla uprawnionych.
func (a *adapterStudia) UtworzOdwolanieWersji(ctx context.Context,
	z shared.StudioVersionReferenceCreateRequest) (shared.StudioVersionReferenceCreateResponse, error) {

	if strings.TrimSpace(z.VersionId) == "" {
		return shared.StudioVersionReferenceCreateResponse{}, bladWskazaniaStudio(
			"odwołanie do wersji bez wskazania wersji")
	}
	wersja, err := a.repozytorium.Wersja(ctx, strings.TrimSpace(z.VersionId))
	if err != nil {
		return shared.StudioVersionReferenceCreateResponse{}, bladWskazaniaStronyPorownania(z.VersionId, err)
	}

	zasieg := shared.ConfigScope(shared.ConfigScopeSession)
	if z.Scope != nil && strings.TrimSpace(string(*z.Scope)) != "" {
		zasieg = *z.Scope
	}
	var wygasa *string
	if z.ExpiresAt != nil {
		chwila := time.UnixMilli(*z.ExpiresAt).UTC().Format(time.RFC3339Nano)
		wygasa = &chwila
	}

	odwolanie, err := a.repozytorium.ZapiszOdwolanieWersji(ctx, dane.OdwolanieWersji{
		Kod:       nowyIdentyfikator(przedrostekOdwolaniaStudia),
		WersjaKod: wersja.Kod,
		Zasieg:    string(zasieg),
		Wygasa:    wygasa,
	})
	if err != nil {
		return shared.StudioVersionReferenceCreateResponse{}, bladStudio(err)
	}

	return shared.StudioVersionReferenceCreateResponse{
		Reference: adresOdwolaniaStudia(wersja.DokumentKod, wersja.Kod, odwolanie.Kod),
		Scope:     zasieg,
	}, nil
}

// adresOdwolaniaStudia składa adres wersji w obrębie platformy. Niesie oba
// człony — dokument i wersję — bo odbiorca odwołania ma otworzyć dokument
// USTAWIONY na tej wersji, a nie szukać, do którego dokumentu wersja należy.
func adresOdwolaniaStudia(kodDokumentu, kodWersji, kodOdwolania string) string {
	return "danaco://studio/" + kodDokumentu + "/wersja/" + kodWersji + "?odwolanie=" + kodOdwolania
}

// galazStudia czyta gałąź i odróżnia jej brak od usterki odczytu.
func (a *adapterStudia) galazStudia(ctx context.Context, kod string) (dane.GalazStudia, error) {
	galaz, err := a.repozytorium.Galaz(ctx, strings.TrimSpace(kod))
	if err == nil {
		return galaz, nil
	}
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.GalazStudia{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: gałąź "+kod+" nie istnieje"))
	}
	return dane.GalazStudia{}, bladStudio(err)
}

// trescCzolaGalezi oddaje treść wersji bieżącej gałęzi; gałąź bez czoła czyta
// się treścią swojego punktu startowego.
func (a *adapterStudia) trescCzolaGalezi(ctx context.Context, galaz dane.GalazStudia) (string, error) {
	kod := galaz.WersjaStartowaID
	if galaz.WersjaBiezacaID != nil && strings.TrimSpace(*galaz.WersjaBiezacaID) != "" {
		kod = *galaz.WersjaBiezacaID
	}
	return a.trescWersjiPoKodzie(ctx, kod)
}

// trescWersjiPoKodzie czyta treść wersji wskazanej kodem.
func (a *adapterStudia) trescWersjiPoKodzie(ctx context.Context, kod string) (string, error) {
	wersja, err := a.repozytorium.Wersja(ctx, kod)
	if err != nil {
		return "", bladWskazaniaStronyPorownania(kod, err)
	}
	return a.trescZOdwolania(wersja.Tresc, wersja.TrescOdwolanie)
}

// zlozGalazStudia składa gałąź kontraktu z wiersza warstwy danych.
func zlozGalazStudia(wiersz dane.GalazStudia) shared.StudioBranch {
	return shared.StudioBranch{
		Id:            wiersz.Kod,
		DocumentId:    wiersz.DokumentKod,
		Name:          wiersz.Nazwa,
		FromVersionId: wiersz.WersjaStartowaID,
		HeadVersionId: wiersz.WersjaBiezacaID,
		Merged:        wiersz.Scalona,
		CreatedAt:     chwilaBazy(wiersz.Utworzono),
	}
}

// ── Scalanie trójstronne ────────────────────────────────────────────────────

// wynikScaleniaStudia niesie albo treść scaloną, albo komplet konfliktów.
// Nigdy oba naraz: treść złożona z konfliktem nierozstrzygniętym byłaby treścią,
// w której jeden z wariantów przepadł bez śladu.
type wynikScaleniaStudia struct {
	tresc     string
	konflikty []shared.StudioMergeConflict
}

// edycjaScaleniaStudia to jedna zmiana jednej strony wobec wspólnego przodka:
// zakres wierszy przodka i treść, która ma go zastąpić.
type edycjaScaleniaStudia struct {
	od, do int
	tresc  []string
	// zeScalanej odróżnia stronę zmiany. Konflikt powstaje wtedy i tylko wtedy,
	// gdy ten sam zakres przodka ruszyły OBIE strony — jedna strona zmieniająca
	// dwa sąsiadujące fragmenty konfliktem nie jest.
	zeScalanej bool
}

// scalTrojstronnieStudia składa treść z przodka i dwóch wariantów.
func scalTrojstronnieStudia(przodek, scalana, docelowa string,
	rozstrzygniecia []shared.StudioMergeResolution) wynikScaleniaStudia {

	wierszePrzodka := podzielNaWiersze(przodek)
	edycje := []edycjaScaleniaStudia{}
	for _, blok := range blokiZmianyStudia(wierszePrzodka, podzielNaWiersze(scalana)) {
		edycje = append(edycje, edycjaScaleniaStudia{
			od: blok.odA, do: blok.doA, tresc: blok.trescB, zeScalanej: true,
		})
	}
	for _, blok := range blokiZmianyStudia(wierszePrzodka, podzielNaWiersze(docelowa)) {
		edycje = append(edycje, edycjaScaleniaStudia{
			od: blok.odA, do: blok.doA, tresc: blok.trescB,
		})
	}
	sort.SliceStable(edycje, func(i, j int) bool {
		if edycje[i].od == edycje[j].od {
			return edycje[i].do < edycje[j].do
		}
		return edycje[i].od < edycje[j].od
	})

	przesuniecia := przesunieciaWierszyStudia(wierszePrzodka)
	poRozstrzygnieciu := map[int]shared.StudioMergeResolution{}
	for _, r := range rozstrzygniecia {
		poRozstrzygnieciu[r.Index] = r
	}

	wynik := []string{}
	konflikty := []shared.StudioMergeConflict{}
	pozycja, i, numerKonfliktu := 0, 0, 0
	for i < len(edycje) {
		grupa := []edycjaScaleniaStudia{edycje[i]}
		poczatek, koniec := edycje[i].od, edycje[i].do
		for i+1 < len(edycje) && edycje[i+1].od < koniec {
			i++
			grupa = append(grupa, edycje[i])
			if edycje[i].do > koniec {
				koniec = edycje[i].do
			}
		}
		i++

		wynik = append(wynik, wierszePrzodka[pozycja:poczatek]...)
		trescScalanej := zastosujGrupeStudia(wierszePrzodka, grupa, poczatek, koniec, true)
		trescDocelowej := zastosujGrupeStudia(wierszePrzodka, grupa, poczatek, koniec, false)

		obieStrony := obieStronyWGrupieStudia(grupa)
		if !obieStrony || trescScalanej == trescDocelowej {
			// Zmieniła jedna strona albo obie zmieniły tak samo — nie ma czego
			// rozstrzygać, więc rdzeń bierze zmianę bez pytania.
			if obieStrony || grupa[0].zeScalanej {
				wynik = append(wynik, podzielNaWiersze(trescScalanej)...)
			} else {
				wynik = append(wynik, podzielNaWiersze(trescDocelowej)...)
			}
			pozycja = koniec
			continue
		}

		numerKonfliktu++
		if rozstrzygniecie, podane := poRozstrzygnieciu[numerKonfliktu]; podane {
			wynik = append(wynik, podzielNaWiersze(trescRozstrzygnieciaStudia(
				rozstrzygniecie, trescScalanej, trescDocelowej))...)
			pozycja = koniec
			continue
		}
		podstawa := strings.Join(wierszePrzodka[poczatek:koniec], "\n")
		konflikty = append(konflikty, shared.StudioMergeConflict{
			Index:      numerKonfliktu,
			RangeStart: przesuniecia[poczatek],
			RangeEnd:   przesuniecia[koniec],
			Source:     trescScalanej,
			Target:     trescDocelowej,
			Base:       &podstawa,
		})
		pozycja = koniec
	}
	wynik = append(wynik, wierszePrzodka[pozycja:]...)

	if len(konflikty) > 0 {
		return wynikScaleniaStudia{konflikty: konflikty}
	}
	return wynikScaleniaStudia{tresc: strings.Join(wynik, "\n")}
}

// obieStronyWGrupieStudia mówi, czy w grupie nakładających się zmian są obie
// strony scalenia.
func obieStronyWGrupieStudia(grupa []edycjaScaleniaStudia) bool {
	scalana, docelowa := false, false
	for _, e := range grupa {
		if e.zeScalanej {
			scalana = true
		} else {
			docelowa = true
		}
	}
	return scalana && docelowa
}

// zastosujGrupeStudia odtwarza, jak wskazana strona widzi zakres przodka:
// nakłada na wiersze przodka wyłącznie własne zmiany z grupy.
func zastosujGrupeStudia(przodek []string, grupa []edycjaScaleniaStudia,
	poczatek, koniec int, zeScalanej bool) string {

	wynik := []string{}
	pozycja := poczatek
	for _, e := range grupa {
		if e.zeScalanej != zeScalanej {
			continue
		}
		if e.od > pozycja {
			wynik = append(wynik, przodek[pozycja:e.od]...)
		}
		wynik = append(wynik, e.tresc...)
		pozycja = e.do
	}
	if pozycja < koniec {
		wynik = append(wynik, przodek[pozycja:koniec]...)
	}
	return strings.Join(wynik, "\n")
}

// trescRozstrzygnieciaStudia oddaje treść wybraną przez Operatora. Strona
// „obie" zestawia warianty jeden pod drugim, bo to jedyne złożenie, które
// niczego nie gubi; strona „własna" bez treści zachowuje się jak wybór gałęzi
// docelowej, bo pusta treść własna nie jest rozstrzygnięciem.
func trescRozstrzygnieciaStudia(r shared.StudioMergeResolution, scalana, docelowa string) string {
	switch r.Side {
	case shared.StudioMergeSideScalana:
		return scalana
	case shared.StudioMergeSideObie:
		return scalana + "\n" + docelowa
	case shared.StudioMergeSideWlasna:
		if r.Text != nil {
			return *r.Text
		}
		return docelowa
	default:
		return docelowa
	}
}

// przesunieciaWierszyStudia liczy przesunięcia znakowe początków wierszy wraz
// z przesunięciem końca treści. Konflikt kontraktu podaje zakres w ZNAKACH, bo
// Diff Panel zaznacza go w treści, a nie w numerach wierszy.
func przesunieciaWierszyStudia(wiersze []string) []int {
	przesuniecia := make([]int, len(wiersze)+1)
	suma := 0
	for i, wiersz := range wiersze {
		przesuniecia[i] = suma
		suma += len([]rune(wiersz)) + 1
	}
	if suma > 0 {
		suma--
	}
	przesuniecia[len(wiersze)] = suma
	return przesuniecia
}

// ── Bloki zmiany między dwiema treściami ────────────────────────────────────

// blokZmianyStudia to jeden ciągły fragment, w którym treść B odbiega od A.
type blokZmianyStudia struct {
	odA, doA int
	trescB   []string
}

// granicaTablicyScaleniaStudia chroni pamięć przed porównaniem, które nie
// mieści się w rozsądku: tablica podobieństwa rośnie iloczynem długości.
// Powyżej tej granicy porównanie schodzi na przycinanie wspólnego przedrostka
// i sufiksu — mniej dokładne, ale rozstrzygalne i skończone.
const granicaTablicyScaleniaStudia = 4_000_000

// blokiZmianyStudia wyznacza fragmenty, którymi B różni się od A.
//
// Podstawą jest najdłuższy wspólny podciąg wierszy: to on mówi, co zostało
// nietknięte, a wszystko pomiędzy dwoma kolejnymi wspólnymi wierszami jest
// jednym blokiem zmiany. Dla scalania trójstronnego to rozstrzygnięcie jest
// konieczne — przycinanie przedrostka i sufiksu (którym jedzie
// `studio.diff.compare`) dałoby JEDEN wielki blok na całą treść, a wtedy każde
// scalenie dwóch redakcji tego samego dokumentu byłoby konfliktem.
func blokiZmianyStudia(a, b []string) []blokZmianyStudia {
	if len(a)*len(b) > granicaTablicyScaleniaStudia {
		return blokiZPrzycieciaStudia(a, b)
	}

	// Tablica długości wspólnego podciągu — klasyczne wypełnianie od końca.
	dlugosci := make([][]int, len(a)+1)
	for i := range dlugosci {
		dlugosci[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dlugosci[i][j] = dlugosci[i+1][j+1] + 1
				continue
			}
			if dlugosci[i+1][j] >= dlugosci[i][j+1] {
				dlugosci[i][j] = dlugosci[i+1][j]
			} else {
				dlugosci[i][j] = dlugosci[i][j+1]
			}
		}
	}

	bloki := []blokZmianyStudia{}
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		if i < len(a) && j < len(b) && a[i] == b[j] {
			i, j = i+1, j+1
			continue
		}
		odA, odB := i, j
		for i < len(a) || j < len(b) {
			if i < len(a) && j < len(b) && a[i] == b[j] {
				break
			}
			switch {
			case j < len(b) && (i == len(a) || dlugosci[i][j+1] >= dlugosci[i+1][j]):
				j++
			default:
				i++
			}
		}
		bloki = append(bloki, blokZmianyStudia{odA: odA, doA: i, trescB: append([]string{}, b[odB:j]...)})
	}
	return bloki
}

// blokiZPrzycieciaStudia jest drogą zapasową dla treści zbyt długich na tablicę
// podobieństwa: jeden blok między wspólnym przedrostkiem a wspólnym sufiksem.
func blokiZPrzycieciaStudia(a, b []string) []blokZmianyStudia {
	przedrostek := 0
	for przedrostek < len(a) && przedrostek < len(b) && a[przedrostek] == b[przedrostek] {
		przedrostek++
	}
	sufiks := 0
	for sufiks < len(a)-przedrostek && sufiks < len(b)-przedrostek &&
		a[len(a)-1-sufiks] == b[len(b)-1-sufiks] {
		sufiks++
	}
	if przedrostek == len(a) && przedrostek == len(b) {
		return []blokZmianyStudia{}
	}
	return []blokZmianyStudia{{
		odA:    przedrostek,
		doA:    len(a) - sufiks,
		trescB: append([]string{}, b[przedrostek:len(b)-sufiks]...),
	}}
}
