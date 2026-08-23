// Moduł Apps — Architecture Designer poza samą definicją układu: walidacja
// architektury, historia wersji, adnotacje projektowe i eksport diagramu.
//
// Obsługiwane komendy: `apps.architecture.validate`,
// `apps.architecture.version.list`, `apps.architecture.annotation.save`,
// `apps.architecture.export`.
//
// WALIDACJA LICZY, NIE PRZECHOWUJE. Migracja 051 zostawiła w architekturze
// kolumnę `zastrzezenia_walidacji`, a `apps.architecture.define` przepisywał ją
// bez zmiany, bo nie miał z czego liczyć nowych. `apps.architecture.validate`
// jest tym miejscem, które liczy: przechodzi komponenty i graf zależności
// i wykrywa cztery rzeczy, o których mówi opracowanie („komponenty bez połączeń
// i brakujące zależności"):
//   - komponent bez ani jednej krawędzi (ostrzeżenie — jest na kanwie, a nic go
//     nie dotyczy),
//   - krawędź wskazującą komponent spoza układu (zastrzeżenie poważne — układ
//     odwołuje się do czegoś, czego nie ma),
//   - cykl w grafie zależności (zastrzeżenie poważne),
//   - układ bez ani jednego komponentu (spostrzeżenie).
//
// Policzone zastrzeżenia wracają do kolumny, więc kolejny `apps.architecture.get`
// pokazuje je bez powtarzania rachunku, a `apps.architecture.define` — który je
// przepisuje bez zmiany — nie kasuje pracy walidatora.
//
// ŻADNE ZASTRZEŻENIE NICZEGO NIE BLOKUJE. Kontrakt mówi to wprost przy
// `AppValidationIssue` („OSTRZEZENIE, nie brama"), a opracowanie powtarza przy
// narzędziu walidacji („nieblokujące dalszej pracy”). Walidacja jest komendą
// odczytu z zapisem wyniku, nie warunkiem zapisu układu.
//
// EKSPORT WYTWARZA PLIK, NIE OPIS PLIKU. Cztery formaty kontraktu powstają
// bibliotekami wkompilowanymi w binarium — SVG, Mermaid i Markdown są tekstem
// składanym tutaj, PNG rysuje `image/png` ze stdlib wraz z czcionką rastrową
// `x/image/font/basicfont`. Żaden nie woła programu z zewnątrz, więc eksport
// działa na instalce, która niesie sam rdzeń.
package core

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekAdnotacjiApp znakuje identyfikatory notatek projektowych.
const przedrostekAdnotacjiApp = "adn-"

// Kody zastrzeżeń walidacji układu. Kody są nazwami własnymi walidatora, nie
// numeracją: kontrakt każe podać `code`, a klient rozpoznaje po nim rodzaj
// zastrzeżenia bez czytania treści.
const (
	kodZastrzezeniaOsamotnionyApp = "componentWithoutDependencies"
	kodZastrzezeniaNieznanyApp    = "dependencyOnUnknownComponent"
	kodZastrzezeniaCyklApp        = "dependencyCycle"
	kodZastrzezeniaPustyApp       = "architectureWithoutComponents"
)

// SprawdzArchitekture obsługuje `apps.architecture.validate`.
func (a *adapterAplikacji) SprawdzArchitekture(ctx context.Context,
	z shared.AppsArchitectureValidateRequest) (shared.AppsArchitectureValidateResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.architecture.validate")
	if err != nil {
		return shared.AppsArchitectureValidateResponse{}, err
	}
	wiersz, err := a.architekturaZadania(ctx, okno, z.ArchitectureId)
	if err != nil {
		return shared.AppsArchitectureValidateResponse{}, err
	}

	komponenty, err := a.repozytorium.Komponenty(ctx, wiersz.ID)
	if err != nil {
		return shared.AppsArchitectureValidateResponse{}, bladAplikacji(err)
	}
	zaleznosci, err := a.repozytorium.ZaleznosciKomponentow(ctx, wiersz.ID)
	if err != nil {
		return shared.AppsArchitectureValidateResponse{}, bladAplikacji(err)
	}

	zastrzezenia := policzZastrzezeniaUkladuApp(komponenty, zaleznosci)

	// Wynik wraca do kolumny, żeby `apps.architecture.get` pokazywał go bez
	// powtarzania rachunku. Nieudany zapis nie gubi wyniku: zastrzeżenia i tak
	// jadą w odpowiedzi, bo to one są odpowiedzią na tę komendę.
	tresci := make([]string, 0, len(zastrzezenia))
	for _, zastrzezenie := range zastrzezenia {
		tresci = append(tresci, zastrzezenie.Code+": "+zastrzezenie.Message)
	}
	komponentyZapisu, zaleznosciZapisu := rozlozKomponentyWierszy(komponenty, zaleznosci)
	if _, err := a.repozytorium.ZapiszArchitekture(ctx, dane.ArchitekturaApp{
		Kod: wiersz.Kod, Okno: wiersz.Okno, Nazwa: wiersz.Nazwa, Szablon: wiersz.Szablon,
		ZastrzezeniaWalidacji: tresci,
		RoznicaWersji: wskaznikNapisuApp("walidacja układu: zastrzeżeń " +
			strconv.Itoa(len(zastrzezenia))),
	}, komponentyZapisu, zaleznosciZapisu); err != nil {
		return shared.AppsArchitectureValidateResponse{}, bladAplikacji(err)
	}

	return shared.AppsArchitectureValidateResponse{
		Issues:      zastrzezenia,
		ValidatedAt: time.Now().UTC().UnixMilli(),
	}, nil
}

// WypiszWersjeArchitektury obsługuje `apps.architecture.version.list`.
func (a *adapterAplikacji) WypiszWersjeArchitektury(ctx context.Context,
	z shared.AppsArchitectureVersionListRequest) (shared.AppsArchitectureVersionListResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.architecture.version.list")
	if err != nil {
		return shared.AppsArchitectureVersionListResponse{}, err
	}
	wiersz, err := a.architekturaZadania(ctx, okno, nil)
	if err != nil {
		return shared.AppsArchitectureVersionListResponse{}, err
	}
	wiersze, err := a.repozytorium.WersjeArchitekturyApp(ctx, wiersz.ID)
	if err != nil {
		return shared.AppsArchitectureVersionListResponse{}, bladAplikacji(err)
	}

	razem := len(wiersze)
	if z.Limit != nil && *z.Limit > 0 && *z.Limit < razem {
		wiersze = wiersze[:*z.Limit]
	}
	wersje := make([]shared.AppArchitectureVersion, 0, len(wiersze))
	for _, wersja := range wiersze {
		wersje = append(wersje, shared.AppArchitectureVersion{
			Version:        wersja.Wersja,
			ArchitectureId: wersja.KodArchitektury,
			ComponentCount: wersja.LiczbaKomponentow,
			DiffSummary:    wersja.Roznica,
			CreatedAt:      chwilaBazy(wersja.Utworzono),
		})
	}
	return shared.AppsArchitectureVersionListResponse{Versions: wersje, Total: razem}, nil
}

// ZapiszAdnotacje obsługuje `apps.architecture.annotation.save`. Notatka
// przypina się do komponentu albo do zależności — obu naraz nie, bo wskazywałaby
// dwa różne miejsca kanwy. Pół krawędzi (jeden koniec zależności) też jest
// odmową: nie ma czego podświetlić.
func (a *adapterAplikacji) ZapiszAdnotacje(ctx context.Context,
	z shared.AppsArchitectureAnnotationSaveRequest) (shared.AppsArchitectureAnnotationSaveResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.architecture.annotation.save")
	if err != nil {
		return shared.AppsArchitectureAnnotationSaveResponse{}, err
	}
	if strings.TrimSpace(z.Text) == "" {
		return shared.AppsArchitectureAnnotationSaveResponse{}, bladWskazaniaAplikacji(
			"apps.architecture.annotation.save wymaga treści notatki")
	}
	komponent := wskaznikNapisuApp(strings.TrimSpace(wartoscTekstu(z.ComponentId)))
	zaleznoscZ := wskaznikNapisuApp(strings.TrimSpace(wartoscTekstu(z.DependencyFrom)))
	zaleznoscDo := wskaznikNapisuApp(strings.TrimSpace(wartoscTekstu(z.DependencyTo)))
	if komponent != nil && (zaleznoscZ != nil || zaleznoscDo != nil) {
		return shared.AppsArchitectureAnnotationSaveResponse{}, bladWskazaniaAplikacji(
			"notatka wskazuje naraz komponent i zależność — przypięcie jest jedno albo żadne")
	}
	if (zaleznoscZ == nil) != (zaleznoscDo == nil) {
		return shared.AppsArchitectureAnnotationSaveResponse{}, bladWskazaniaAplikacji(
			"notatka przypięta do zależności wymaga obu jej końców (dependencyFrom i dependencyTo)")
	}

	wiersz, err := a.architekturaZadania(ctx, okno, nil)
	if err != nil {
		return shared.AppsArchitectureAnnotationSaveResponse{}, err
	}

	kod := strings.TrimSpace(wartoscTekstu(z.AnnotationId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekAdnotacjiApp)
	} else if zastana, err := a.repozytorium.AdnotacjaApp(ctx, kod); err != nil {
		return shared.AppsArchitectureAnnotationSaveResponse{},
			bladNieznanegoBytuApp("adnotacja architektury", kod, err)
	} else if zastana.ArchitekturaID != wiersz.ID {
		return shared.AppsArchitectureAnnotationSaveResponse{}, bladWskazaniaAplikacji(
			"adnotacja " + kod + " należy do innej architektury niż architektura okna " + okno)
	}

	zapisana, err := a.repozytorium.ZapiszAdnotacjeApp(ctx, dane.AdnotacjaArchitekturyApp{
		Kod: kod, ArchitekturaID: wiersz.ID, KomponentKod: komponent,
		ZaleznoscZ: zaleznoscZ, ZaleznoscDo: zaleznoscDo, Tresc: z.Text,
	})
	if err != nil {
		return shared.AppsArchitectureAnnotationSaveResponse{}, bladAplikacji(err)
	}
	return shared.AppsArchitectureAnnotationSaveResponse{
		Annotation: shared.AppAnnotation{
			Id: zapisana.Kod, ComponentId: zapisana.KomponentKod,
			DependencyFrom: zapisana.ZaleznoscZ, DependencyTo: zapisana.ZaleznoscDo,
			Text: zapisana.Tresc, UpdatedAt: chwilaBazy(zapisana.Zaktualizowano),
		},
	}, nil
}

// WyeksportujArchitekture obsługuje `apps.architecture.export`: składa plik
// w jednym z czterech formatów kontraktu, kładzie go w magazynie treści rdzenia
// i oddaje odwołanie wraz z rozmiarem, który plik naprawdę ma.
func (a *adapterAplikacji) WyeksportujArchitekture(ctx context.Context,
	z shared.AppsArchitectureExportRequest) (shared.AppsArchitectureExportResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.architecture.export")
	if err != nil {
		return shared.AppsArchitectureExportResponse{}, err
	}
	if err := sprawdzFormatEksportuApp(z.Format); err != nil {
		return shared.AppsArchitectureExportResponse{}, err
	}
	if a.magazyn == nil {
		return shared.AppsArchitectureExportResponse{}, bladBrakuMagazynuApp("eksport diagramu")
	}

	wiersz, err := a.architekturaZadania(ctx, okno, nil)
	if err != nil {
		return shared.AppsArchitectureExportResponse{}, err
	}
	komponenty, err := a.repozytorium.Komponenty(ctx, wiersz.ID)
	if err != nil {
		return shared.AppsArchitectureExportResponse{}, bladAplikacji(err)
	}
	zaleznosci, err := a.repozytorium.ZaleznosciKomponentow(ctx, wiersz.ID)
	if err != nil {
		return shared.AppsArchitectureExportResponse{}, bladAplikacji(err)
	}

	bajty, err := zlozEksportUkladuApp(z.Format, wiersz, komponenty, zaleznosci)
	if err != nil {
		return shared.AppsArchitectureExportResponse{}, bladAplikacji(err)
	}
	odwolanie, rozmiar, err := a.wniesDoMagazynuApp(bajty)
	if err != nil {
		return shared.AppsArchitectureExportResponse{}, bladAplikacji(err)
	}
	return shared.AppsArchitectureExportResponse{ArtifactRef: odwolanie, SizeBytes: rozmiar}, nil
}

// architekturaZadania rozwiązuje architekturę żądania: wskazaną wprost albo
// bieżącą architekturę okna. Brak architektury jest odmową `not_found`, bo
// wszystkie cztery komendy tego pliku pracują NA układzie — walidacja układu,
// którego nie ma, nie miałaby czego zwalidować, a pusta odpowiedź udana
// kazałaby oknu zgadywać, czy układ jest bez zastrzeżeń, czy go nie ma.
func (a *adapterAplikacji) architekturaZadania(ctx context.Context, okno string,
	wskazanie *string) (dane.ArchitekturaApp, error) {

	kod := strings.TrimSpace(wartoscTekstu(wskazanie))
	if kod != "" {
		wiersz, err := a.repozytorium.Architektura(ctx, kod)
		if err != nil {
			return dane.ArchitekturaApp{}, bladNieznanejArchitektury(kod, err)
		}
		if wiersz.Okno != okno {
			return dane.ArchitekturaApp{}, bladWskazaniaAplikacji(
				"architektura " + kod + " należy do okna " + wiersz.Okno +
					", a żądanie przyszło z okna " + okno)
		}
		return wiersz, nil
	}
	wiersz, err := a.repozytorium.ArchitekturaOkna(ctx, okno)
	if err != nil {
		return dane.ArchitekturaApp{}, bladNieznanegoBytuApp("architektura okna", okno, err)
	}
	return wiersz, nil
}

// wniesDoMagazynuApp kładzie bajty wytworu w magazynie treści rdzenia i oddaje
// odwołanie względne katalogu danych wraz z rozmiarem. Rozmiar bierze się
// z długości zapisanych bajtów, nie z zamiaru — odwołanie i liczba mają
// opisywać ten sam plik.
func (a *adapterAplikacji) wniesDoMagazynuApp(bajty []byte) (string, int64, error) {
	if len(bajty) == 0 {
		return "", 0, fmt.Errorf("moduł Apps: wytwór o zerowej długości nie ma czego opisywać")
	}
	suma := sumaTresciApp(bajty)
	sciezka, err := a.magazyn.Zapisz(bajty, suma)
	if err != nil {
		return "", 0, err
	}
	return odwolanieWytworuApp(sciezka), int64(len(bajty)), nil
}

// policzZastrzezeniaUkladuApp jest walidatorem układu — patrz czoło pliku.
func policzZastrzezeniaUkladuApp(komponenty []dane.KomponentArchitektury,
	zaleznosci []dane.ZaleznoscKomponentu) []shared.AppValidationIssue {

	zastrzezenia := []shared.AppValidationIssue{}
	if len(komponenty) == 0 {
		return append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityInfo,
			Code:     kodZastrzezeniaPustyApp,
			Message:  "układ nie ma ani jednego komponentu",
		})
	}

	znane := map[string]struct{}{}
	for _, komponent := range komponenty {
		znane[komponent.KodZewnetrzny] = struct{}{}
	}

	polaczone := map[string]struct{}{}
	nastepnicy := map[string][]string{}
	for _, zaleznosc := range zaleznosci {
		polaczone[zaleznosc.KomponentZ] = struct{}{}
		polaczone[zaleznosc.KomponentDo] = struct{}{}
		nastepnicy[zaleznosc.KomponentZ] = append(nastepnicy[zaleznosc.KomponentZ], zaleznosc.KomponentDo)

		for _, koniec := range []string{zaleznosc.KomponentZ, zaleznosc.KomponentDo} {
			if _, jest := znane[koniec]; jest {
				continue
			}
			kopia := koniec
			zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
				Severity:    shared.AppValidationSeverityError,
				Code:        kodZastrzezeniaNieznanyApp,
				Message:     "zależność wskazuje komponent spoza układu: " + koniec,
				ComponentId: &kopia,
			})
		}
	}

	for _, komponent := range komponenty {
		if _, jest := polaczone[komponent.KodZewnetrzny]; jest {
			continue
		}
		kopia := komponent.KodZewnetrzny
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity:    shared.AppValidationSeverityWarning,
			Code:        kodZastrzezeniaOsamotnionyApp,
			Message:     "komponent " + komponent.Nazwa + " nie ma ani jednej zależności",
			ComponentId: &kopia,
		})
	}

	for _, wierzcholek := range cykleUkladuApp(znane, nastepnicy) {
		kopia := wierzcholek
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity:    shared.AppValidationSeverityError,
			Code:        kodZastrzezeniaCyklApp,
			Message:     "komponent " + wierzcholek + " leży w cyklu zależności",
			ComponentId: &kopia,
		})
	}
	return zastrzezenia
}

// cykleUkladuApp zwraca komponenty leżące w cyklu zależności, uporządkowane.
// Obieg jest zwykłym przejściem w głąb z trzema kolorami — układ produktu ma
// dziesiątki węzłów, nie miliony, więc prostota bije tu spryt.
func cykleUkladuApp(znane map[string]struct{}, nastepnicy map[string][]string) []string {
	const (
		bialy  = 0
		szary  = 1
		czarny = 2
	)
	kolor := map[string]int{}
	wCyklu := map[string]struct{}{}

	var idz func(string)
	idz = func(wierzcholek string) {
		kolor[wierzcholek] = szary
		for _, nastepny := range nastepnicy[wierzcholek] {
			if _, jest := znane[nastepny]; !jest {
				continue
			}
			switch kolor[nastepny] {
			case bialy:
				idz(nastepny)
			case szary:
				wCyklu[nastepny] = struct{}{}
				wCyklu[wierzcholek] = struct{}{}
			}
		}
		kolor[wierzcholek] = czarny
	}

	poczatki := make([]string, 0, len(znane))
	for wierzcholek := range znane {
		poczatki = append(poczatki, wierzcholek)
	}
	sort.Strings(poczatki)
	for _, wierzcholek := range poczatki {
		if kolor[wierzcholek] == bialy {
			idz(wierzcholek)
		}
	}

	lista := make([]string, 0, len(wCyklu))
	for wierzcholek := range wCyklu {
		lista = append(lista, wierzcholek)
	}
	sort.Strings(lista)
	return lista
}

// rozlozKomponentyWierszy oddaje wiersze komponentów i zależności w kształcie,
// którego oczekuje zapis architektury. Walidacja przepisuje układ bez zmiany —
// zmienia wyłącznie kolumnę zastrzeżeń — więc musi podać go z powrotem
// w komplecie, bo zapis jest wymianą („usuń, wstaw od nowa").
func rozlozKomponentyWierszy(komponenty []dane.KomponentArchitektury,
	zaleznosci []dane.ZaleznoscKomponentu) ([]dane.KomponentArchitektury, []dane.ZaleznoscKomponentu) {

	wiersze := make([]dane.KomponentArchitektury, 0, len(komponenty))
	for _, komponent := range komponenty {
		komponent.ID = 0
		wiersze = append(wiersze, komponent)
	}
	if zaleznosci == nil {
		zaleznosci = []dane.ZaleznoscKomponentu{}
	}
	return wiersze, zaleznosci
}

// zlozEksportUkladuApp składa bajty wytworu w żądanym formacie.
func zlozEksportUkladuApp(format shared.AppExportFormat, architektura dane.ArchitekturaApp,
	komponenty []dane.KomponentArchitektury, zaleznosci []dane.ZaleznoscKomponentu) ([]byte, error) {

	switch format {
	case shared.AppExportFormatMermaid:
		return []byte(eksportMermaidApp(komponenty, zaleznosci)), nil
	case shared.AppExportFormatMarkdown:
		return []byte(eksportMarkdownApp(architektura, komponenty, zaleznosci)), nil
	case shared.AppExportFormatSvg:
		return []byte(eksportSvgApp(komponenty, zaleznosci)), nil
	case shared.AppExportFormatPng:
		return eksportPngApp(komponenty, zaleznosci)
	}
	return nil, fmt.Errorf("moduł Apps: nieobsłużony format eksportu %q", format)
}

// eksportMermaidApp składa diagram w notacji Mermaid.
func eksportMermaidApp(komponenty []dane.KomponentArchitektury,
	zaleznosci []dane.ZaleznoscKomponentu) string {

	var zapis strings.Builder
	zapis.WriteString("graph TD\n")
	for _, komponent := range komponenty {
		zapis.WriteString("    " + identyfikatorMermaidApp(komponent.KodZewnetrzny) +
			"[" + strconv.Quote(komponent.Nazwa+" ("+komponent.Rodzaj+")") + "]\n")
	}
	for _, zaleznosc := range zaleznosci {
		zapis.WriteString("    " + identyfikatorMermaidApp(zaleznosc.KomponentZ) +
			" --> " + identyfikatorMermaidApp(zaleznosc.KomponentDo) + "\n")
	}
	return zapis.String()
}

// identyfikatorMermaidApp zamienia kod komponentu na identyfikator węzła
// bezpieczny w notacji Mermaid — kod zewnętrzny bywa napisem z myślnikami,
// których notacja w tym miejscu nie przyjmuje.
func identyfikatorMermaidApp(kod string) string {
	var zapis strings.Builder
	for _, znak := range kod {
		if (znak >= 'a' && znak <= 'z') || (znak >= 'A' && znak <= 'Z') ||
			(znak >= '0' && znak <= '9') {
			zapis.WriteRune(znak)
			continue
		}
		zapis.WriteByte('_')
	}
	if zapis.Len() == 0 {
		return "wezel"
	}
	return "w_" + zapis.String()
}

// eksportMarkdownApp składa dokument opisujący komponenty i zależności.
func eksportMarkdownApp(architektura dane.ArchitekturaApp,
	komponenty []dane.KomponentArchitektury, zaleznosci []dane.ZaleznoscKomponentu) string {

	nazwa := architektura.Kod
	if architektura.Nazwa != nil && *architektura.Nazwa != "" {
		nazwa = *architektura.Nazwa
	}
	var zapis strings.Builder
	zapis.WriteString("# Architektura — " + nazwa + "\n\n")
	zapis.WriteString("| | |\n|---|---|\n")
	zapis.WriteString("| Szablon | " + architektura.Szablon + " |\n")
	zapis.WriteString("| Wersja | " + strconv.Itoa(architektura.Wersja) + " |\n")
	zapis.WriteString("| Komponentów | " + strconv.Itoa(len(komponenty)) + " |\n\n")

	zapis.WriteString("## Komponenty\n\n| Komponent | Rodzaj | Stos | Opis |\n|---|---|---|---|\n")
	for _, komponent := range komponenty {
		zapis.WriteString("| " + komponent.Nazwa + " | " + komponent.Rodzaj + " | " +
			wartoscTekstu(komponent.Stos) + " | " + wartoscTekstu(komponent.Opis) + " |\n")
	}

	zapis.WriteString("\n## Zależności\n\n")
	if len(zaleznosci) == 0 {
		zapis.WriteString("Układ nie niesie ani jednej zależności.\n")
		return zapis.String()
	}
	zapis.WriteString("| Z | Do |\n|---|---|\n")
	for _, zaleznosc := range zaleznosci {
		zapis.WriteString("| " + zaleznosc.KomponentZ + " | " + zaleznosc.KomponentDo + " |\n")
	}
	return zapis.String()
}

// Wymiary kanwy eksportu. Jedna kolumna węzłów, bo układ produktu czyta się
// z góry na dół, a rozstawianie grafu w dwóch wymiarach wymagałoby algorytmu
// rozkładu, którego kontrakt nie opisuje i którego nikt nie zamawiał.
const (
	szerokoscKanwyApp = 640
	wysokoscWezlaApp  = 48
	odstepWezlaApp    = 28
	marginesKanwyApp  = 24
	szerokoscWezlaApp = szerokoscKanwyApp - 2*marginesKanwyApp
)

// wysokoscKanwyApp liczy wysokość rysunku dla danej liczby węzłów.
func wysokoscKanwyApp(ile int) int {
	if ile == 0 {
		ile = 1
	}
	return 2*marginesKanwyApp + ile*wysokoscWezlaApp + (ile-1)*odstepWezlaApp
}

// eksportSvgApp rysuje układ jako dokument wektorowy.
func eksportSvgApp(komponenty []dane.KomponentArchitektury,
	zaleznosci []dane.ZaleznoscKomponentu) string {

	wysokosc := wysokoscKanwyApp(len(komponenty))
	pozycje := map[string]int{}
	var zapis strings.Builder
	fmt.Fprintf(&zapis, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" `+
		`viewBox="0 0 %d %d">`+"\n", szerokoscKanwyApp, wysokosc, szerokoscKanwyApp, wysokosc)
	fmt.Fprintf(&zapis, `<rect width="%d" height="%d" fill="#ffffff"/>`+"\n",
		szerokoscKanwyApp, wysokosc)

	for indeks, komponent := range komponenty {
		gora := marginesKanwyApp + indeks*(wysokoscWezlaApp+odstepWezlaApp)
		pozycje[komponent.KodZewnetrzny] = gora
		fmt.Fprintf(&zapis, `<rect x="%d" y="%d" width="%d" height="%d" rx="6" `+
			`fill="#f4f4f5" stroke="#3f3f46"/>`+"\n",
			marginesKanwyApp, gora, szerokoscWezlaApp, wysokoscWezlaApp)
		fmt.Fprintf(&zapis, `<text x="%d" y="%d" font-family="sans-serif" font-size="14" `+
			`fill="#18181b">%s</text>`+"\n",
			marginesKanwyApp+12, gora+29, tekstSvgApp(komponent.Nazwa+" — "+komponent.Rodzaj))
	}

	for _, zaleznosc := range zaleznosci {
		zGora, jestZ := pozycje[zaleznosc.KomponentZ]
		doGora, jestDo := pozycje[zaleznosc.KomponentDo]
		if !jestZ || !jestDo {
			continue
		}
		fmt.Fprintf(&zapis, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3f3f46" `+
			`stroke-width="2"/>`+"\n",
			marginesKanwyApp+szerokoscWezlaApp/2, zGora+wysokoscWezlaApp,
			marginesKanwyApp+szerokoscWezlaApp/2, doGora)
	}
	zapis.WriteString("</svg>\n")
	return zapis.String()
}

// tekstSvgApp zabezpiecza znaki, które w dokumencie XML mają własne znaczenie.
func tekstSvgApp(tekst string) string {
	zamiana := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;",
		`"`, "&quot;", "'", "&apos;")
	return zamiana.Replace(tekst)
}

// eksportPngApp rysuje ten sam układ jako obraz rastrowy. Czcionka jest
// rastrowa i wkompilowana (`basicfont.Face7x13`) — rdzeń nie sięga po plik
// czcionki z systemu, więc rysunek wychodzi tak samo na każdej maszynie.
func eksportPngApp(komponenty []dane.KomponentArchitektury,
	zaleznosci []dane.ZaleznoscKomponentu) ([]byte, error) {

	wysokosc := wysokoscKanwyApp(len(komponenty))
	plotno := image.NewRGBA(image.Rect(0, 0, szerokoscKanwyApp, wysokosc))
	bialy := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	obrys := color.RGBA{R: 0x3f, G: 0x3f, B: 0x46, A: 0xff}
	wypelnienie := color.RGBA{R: 0xf4, G: 0xf4, B: 0xf5, A: 0xff}
	draw.Draw(plotno, plotno.Bounds(), &image.Uniform{C: bialy}, image.Point{}, draw.Src)

	rysownik := &font.Drawer{Dst: plotno, Src: &image.Uniform{C: obrys}, Face: basicfont.Face7x13}
	pozycje := map[string]int{}

	for indeks, komponent := range komponenty {
		gora := marginesKanwyApp + indeks*(wysokoscWezlaApp+odstepWezlaApp)
		pozycje[komponent.KodZewnetrzny] = gora
		prostokat := image.Rect(marginesKanwyApp, gora,
			marginesKanwyApp+szerokoscWezlaApp, gora+wysokoscWezlaApp)
		draw.Draw(plotno, prostokat, &image.Uniform{C: wypelnienie}, image.Point{}, draw.Src)
		obrysujProstokatApp(plotno, prostokat, obrys)

		rysownik.Dot = fixed.P(marginesKanwyApp+12, gora+29)
		rysownik.DrawString(komponent.Nazwa + " - " + komponent.Rodzaj)
	}

	for _, zaleznosc := range zaleznosci {
		zGora, jestZ := pozycje[zaleznosc.KomponentZ]
		doGora, jestDo := pozycje[zaleznosc.KomponentDo]
		if !jestZ || !jestDo {
			continue
		}
		poczatek, koniec := zGora+wysokoscWezlaApp, doGora
		if poczatek > koniec {
			poczatek, koniec = koniec, poczatek
		}
		srodek := marginesKanwyApp + szerokoscWezlaApp/2
		for y := poczatek; y < koniec; y++ {
			plotno.Set(srodek, y, obrys)
			plotno.Set(srodek+1, y, obrys)
		}
	}

	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		return nil, fmt.Errorf("moduł Apps: nie można złożyć obrazu układu: %w", err)
	}
	return bufor.Bytes(), nil
}

// obrysujProstokatApp rysuje ramkę węzła jednym pikselem grubości.
func obrysujProstokatApp(plotno *image.RGBA, prostokat image.Rectangle, barwa color.RGBA) {
	for x := prostokat.Min.X; x < prostokat.Max.X; x++ {
		plotno.Set(x, prostokat.Min.Y, barwa)
		plotno.Set(x, prostokat.Max.Y-1, barwa)
	}
	for y := prostokat.Min.Y; y < prostokat.Max.Y; y++ {
		plotno.Set(prostokat.Min.X, y, barwa)
		plotno.Set(prostokat.Max.X-1, y, barwa)
	}
}

// sprawdzFormatEksportuApp dopuszcza wyłącznie formaty kontraktu.
func sprawdzFormatEksportuApp(format shared.AppExportFormat) error {
	switch format {
	case shared.AppExportFormatSvg, shared.AppExportFormatPng,
		shared.AppExportFormatMermaid, shared.AppExportFormatMarkdown:
		return nil
	}
	return bladWskazaniaAplikacji("nieznany format eksportu " + strconv.Quote(string(format)) +
		" — dopuszczalne: svg, png, mermaid, markdown")
}
