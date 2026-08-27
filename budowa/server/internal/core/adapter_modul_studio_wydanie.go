// Plik obsługuje wydawanie pracy Studia na zewnątrz: eksport repozytorium,
// eksport paczki redakcyjnej i eksport raportu różnicy wersji. Metody
// dopisują się na typie adapterStudia zadeklarowanym w pliku
// adapter_modul_studio.go.
package core

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// wpisArchiwumStudia to jedna pozycja archiwum: nazwa wewnątrz paczki, jej
// bajty oraz chwila zapisu przenoszona do znacznika czasu wpisu archiwum.
type wpisArchiwumStudia struct {
	nazwa  string
	bajty  []byte
	chwila time.Time
}

// manifestWydaniaStudia opisuje zawartość archiwum słowami, których nie da się
// odczytać z samych nazw plików.
type manifestWydaniaStudia struct {
	Dokument   string                 `json:"dokument"`
	Tytul      string                 `json:"tytul,omitempty"`
	Format     string                 `json:"format"`
	Wydano     string                 `json:"wydano"`
	Wersje     []manifestWersjiStudia `json:"wersje"`
	Zawartosc  []string               `json:"zawartosc"`
	Adnotacje  int                    `json:"adnotacje,omitempty"`
	Komentarze int                    `json:"komentarze,omitempty"`
}

// manifestWersjiStudia opisuje jedną wersję odłożoną w archiwum: jej kod,
// plik, etykietę, autora, chwilę utworzenia, przynależność do kamienia
// milowego i gałęzi.
type manifestWersjiStudia struct {
	Wersja       string `json:"wersja"`
	Plik         string `json:"plik"`
	Etykieta     string `json:"etykieta,omitempty"`
	Autor        string `json:"autor,omitempty"`
	Utworzono    string `json:"utworzono"`
	KamienMilowy bool   `json:"kamienMilowy,omitempty"`
	Galaz        string `json:"galaz,omitempty"`
}

// WydajRepozytorium obsługuje studio.repository.export. Wskazanie wersji
// zawęża wydanie, a brak wskazania bierze całą historię, tak jak opisuje to
// przycisk „Eksportuj historię".
func (a *adapterStudia) WydajRepozytorium(ctx context.Context,
	z shared.StudioRepositoryExportRequest) (shared.StudioRepositoryExportResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" {
		return shared.StudioRepositoryExportResponse{}, bladWskazaniaStudio(
			"wydanie repozytorium bez wskazania dokumentu")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioRepositoryExportResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	wersje, err := a.wersjeDoWydania(ctx, dokument, z.VersionIds)
	if err != nil {
		return shared.StudioRepositoryExportResponse{}, err
	}

	wpisy, manifest, err := a.wpisyHistoriiStudia(ctx, dokument, wersje)
	if err != nil {
		return shared.StudioRepositoryExportResponse{}, err
	}
	manifest.Zawartosc = nazwyWpisowStudia(wpisy)
	wpisManifestu, err := wpisManifestuStudia(manifest)
	if err != nil {
		return shared.StudioRepositoryExportResponse{}, err
	}
	wpisy = append([]wpisArchiwumStudia{wpisManifestu}, wpisy...)

	format := formatArchiwumStudia(z.Format)
	bajty, err := spakujStudia(wpisy, format)
	if err != nil {
		return shared.StudioRepositoryExportResponse{}, err
	}
	zasob, err := a.odlozTrescStudia(ctx, bajty,
		"historia-"+dokument.Kod+"."+format, format, wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioRepositoryExportResponse{}, err
	}
	return shared.StudioRepositoryExportResponse{
		Asset: zasob, Entries: len(wersje), SizeBytes: len(bajty),
	}, nil
}

// WydajPaczke obsługuje studio.package.export. Paczka jest przekazaniem: obok
// dokumentu finalnego niesie to, co odbiorca musi mieć, żeby zrozumieć, jak
// dokument powstał — historię, raport zmian i adnotacje.
func (a *adapterStudia) WydajPaczke(ctx context.Context,
	z shared.StudioPackageExportRequest) (shared.StudioPackageExportResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" {
		return shared.StudioPackageExportResponse{}, bladWskazaniaStudio(
			"paczka redakcyjna bez wskazania dokumentu")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPackageExportResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	tresc, err := a.trescWersjiStudia(ctx, z.VersionId, dokument)
	if err != nil {
		return shared.StudioPackageExportResponse{}, err
	}

	format := strings.TrimSpace(wartoscTekstu(z.DocumentFormat))
	if format == "" {
		format = string(dokument.Format)
	}
	// Kartka wynika z formatu dokumentu: A4 dla dokumentu A5 wydałoby co
	// innego, niż widać w oknie.
	finalny, err := a.postacDokumentuStudia(tresc, format, dokument,
		a.geometriaDokumentuStudia(ctx, dokument.Kod))
	if err != nil {
		return shared.StudioPackageExportResponse{}, err
	}
	chwila := time.Now().UTC()
	wpisy := []wpisArchiwumStudia{{
		nazwa: "dokument." + format, bajty: finalny, chwila: chwila,
	}}

	manifest := manifestWydaniaStudia{
		Dokument: dokument.Kod, Format: format,
		Wydano: chwila.Format(time.RFC3339), Wersje: []manifestWersjiStudia{},
	}
	if dokument.Tytul != nil {
		manifest.Tytul = *dokument.Tytul
	}

	if wlaczone(z.IncludeHistory) {
		wersje, err := a.repozytorium.Wersje(ctx, dokument.ID)
		if err != nil {
			return shared.StudioPackageExportResponse{}, bladStudio(err)
		}
		wpisyHistorii, manifestHistorii, err := a.wpisyHistoriiStudia(ctx, dokument, wersje)
		if err != nil {
			return shared.StudioPackageExportResponse{}, err
		}
		wpisy = append(wpisy, wpisyHistorii...)
		manifest.Wersje = manifestHistorii.Wersje
	}

	if wlaczone(z.IncludeDiffReport) {
		raport, err := a.raportZmianPaczkiStudia(ctx, dokument, z.VersionId, tresc)
		if err != nil {
			return shared.StudioPackageExportResponse{}, err
		}
		wpisy = append(wpisy, wpisArchiwumStudia{
			nazwa: "raport-zmian.md", bajty: []byte(raport), chwila: chwila,
		})
	}

	if wlaczone(z.IncludeAnnotations) {
		adnotacje, err := a.adnotacjeIKomentarzeStudia(ctx, dokument.ID)
		if err != nil {
			return shared.StudioPackageExportResponse{}, err
		}
		bajty, err := json.MarshalIndent(zlozAdnotacjeWydaniaStudia(adnotacje), "", "  ")
		if err != nil {
			return shared.StudioPackageExportResponse{}, bladStudio(err)
		}
		wpisy = append(wpisy, wpisArchiwumStudia{
			nazwa: "adnotacje.json", bajty: bajty, chwila: chwila,
		})
		for _, wpis := range adnotacje {
			if wpis.Rodzaj == "adnotacja" {
				manifest.Adnotacje++
			} else {
				manifest.Komentarze++
			}
		}
	}

	manifest.Zawartosc = nazwyWpisowStudia(wpisy)
	wpisManifestu, err := wpisManifestuStudia(manifest)
	if err != nil {
		return shared.StudioPackageExportResponse{}, err
	}
	wpisy = append([]wpisArchiwumStudia{wpisManifestu}, wpisy...)

	bajty, err := spakujStudia(wpisy, "zip")
	if err != nil {
		return shared.StudioPackageExportResponse{}, err
	}
	zasob, err := a.odlozTrescStudia(ctx, bajty,
		"paczka-"+dokument.Kod+".zip", "zip", wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioPackageExportResponse{}, err
	}
	return shared.StudioPackageExportResponse{
		Asset: zasob, Entries: len(wpisy), SizeBytes: len(bajty),
	}, nil
}

// WydajRaportRoznicy obsługuje studio.diff.report.export. Raport jest
// dokumentem: niesie fragmenty różnicy, statystykę i adnotacje przypisane do
// fragmentów, gdy Operator ich nie wyłączył.
func (a *adapterStudia) WydajRaportRoznicy(ctx context.Context,
	z shared.StudioDiffReportExportRequest) (shared.StudioDiffReportExportResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" || strings.TrimSpace(z.BaseVersionId) == "" {
		return shared.StudioDiffReportExportResponse{}, bladWskazaniaStudio(
			"raport różnicy wymaga dokumentu i wersji odniesienia")
	}
	if strings.TrimSpace(z.Format) == "" {
		return shared.StudioDiffReportExportResponse{}, bladWskazaniaStudio(
			"raport różnicy bez wskazania formatu")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDiffReportExportResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}

	baza, err := a.trescWersjiPoKodzie(ctx, z.BaseVersionId)
	if err != nil {
		return shared.StudioDiffReportExportResponse{}, err
	}
	cel, err := a.trescWersjiStudia(ctx, z.TargetVersionId, dokument)
	if err != nil {
		return shared.StudioDiffReportExportResponse{}, err
	}

	var adnotacje []dane.KomentarzStudia
	if wlaczone(z.IncludeAnnotations) {
		adnotacje, err = a.repozytorium.Komentarze(ctx, dokument.ID, "adnotacja")
		if err != nil {
			return shared.StudioDiffReportExportResponse{}, bladStudio(err)
		}
	}

	nazwaCelu := "treść bieżąca"
	if kod := strings.TrimSpace(wartoscTekstu(z.TargetVersionId)); kod != "" {
		nazwaCelu = kod
	}
	tekst := trescRaportuRoznicyStudia(dokument, z.BaseVersionId, nazwaCelu,
		policzFragmentyRoznicy(baza, cel), adnotacje)

	format := strings.ToLower(strings.TrimSpace(z.Format))
	// Raport zmian jest sprawozdaniem, nie kopią dokumentu, więc idzie
	// kartką domyślną.
	bajty, err := a.postacDokumentuStudia(tekst, format, dokument, geometriaDomyslnaStudia())
	if err != nil {
		return shared.StudioDiffReportExportResponse{}, err
	}
	zasob, err := a.odlozTrescStudia(ctx, bajty,
		"raport-zmian-"+dokument.Kod+"."+format, format, wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioDiffReportExportResponse{}, err
	}
	return shared.StudioDiffReportExportResponse{Asset: zasob, SizeBytes: len(bajty)}, nil
}

// ── Składanie zawartości ────────────────────────────────────────────────────

// wersjeDoWydania zawęża historię do wersji jawnie wskazanych przez
// wołającego albo oddaje ją całą, gdy wskazania brak.
func (a *adapterStudia) wersjeDoWydania(ctx context.Context, dokument dane.DokumentStudia,
	wskazane []string) ([]dane.WersjaDokumentu, error) {

	wszystkie, err := a.repozytorium.Wersje(ctx, dokument.ID)
	if err != nil {
		return nil, bladStudio(err)
	}
	if len(wskazane) == 0 {
		return wszystkie, nil
	}
	poKodzie := map[string]dane.WersjaDokumentu{}
	for _, wersja := range wszystkie {
		poKodzie[wersja.Kod] = wersja
	}
	wybrane := make([]dane.WersjaDokumentu, 0, len(wskazane))
	for _, kod := range wskazane {
		wersja, jest := poKodzie[strings.TrimSpace(kod)]
		if !jest {
			return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
				"moduł Studio: wersja "+kod+" nie należy do dokumentu "+dokument.Kod))
		}
		wybrane = append(wybrane, wersja)
	}
	return wybrane, nil
}

// wpisyHistoriiStudia składa pliki wersji wraz z opisem każdej z nich do
// wpisania w manifeście archiwum.
func (a *adapterStudia) wpisyHistoriiStudia(ctx context.Context, dokument dane.DokumentStudia,
	wersje []dane.WersjaDokumentu) ([]wpisArchiwumStudia, manifestWydaniaStudia, error) {

	manifest := manifestWydaniaStudia{
		Dokument: dokument.Kod, Format: string(dokument.Format),
		Wydano: time.Now().UTC().Format(time.RFC3339),
		Wersje: make([]manifestWersjiStudia, 0, len(wersje)),
	}
	if dokument.Tytul != nil {
		manifest.Tytul = *dokument.Tytul
	}

	wpisy := make([]wpisArchiwumStudia, 0, len(wersje))
	for _, wersja := range wersje {
		tresc, err := a.trescZOdwolania(wersja.Tresc, wersja.TrescOdwolanie)
		if err != nil {
			return nil, manifestWydaniaStudia{}, err
		}
		nazwa := "wersje/" + wersja.Kod + ".txt"
		wpisy = append(wpisy, wpisArchiwumStudia{
			nazwa: nazwa, bajty: []byte(tresc), chwila: chwilaWpisuStudia(wersja.Utworzono),
		})
		opis := manifestWersjiStudia{
			Wersja: wersja.Kod, Plik: nazwa, Utworzono: wersja.Utworzono,
			KamienMilowy: wersja.KamienMilowy,
		}
		if wersja.Etykieta != nil {
			opis.Etykieta = *wersja.Etykieta
		}
		if wersja.Autor != nil {
			opis.Autor = *wersja.Autor
		}
		if wersja.GalazKod != nil {
			opis.Galaz = *wersja.GalazKod
		}
		manifest.Wersje = append(manifest.Wersje, opis)
	}
	return wpisy, manifest, nil
}

// raportZmianPaczkiStudia składa raport zmian dla paczki: wersja finalna wobec
// najstarszej zapisanej. Historia bez ani jednej wersji daje raport, który to
// mówi wprost — plik pusty w paczce wyglądałby na usterkę pakowania.
func (a *adapterStudia) raportZmianPaczkiStudia(ctx context.Context, dokument dane.DokumentStudia,
	kodFinalnej *string, trescFinalna string) (string, error) {

	wersje, err := a.repozytorium.Wersje(ctx, dokument.ID)
	if err != nil {
		return "", bladStudio(err)
	}
	if len(wersje) == 0 {
		return "# Raport zmian\n\nDokument " + dokument.Kod +
			" nie ma ani jednej zapisanej wersji — nie ma czego porównać.\n", nil
	}
	// Historia idzie od najnowszej, więc odniesieniem jest ostatni wiersz.
	najstarsza := wersje[len(wersje)-1]
	baza, err := a.trescZOdwolania(najstarsza.Tresc, najstarsza.TrescOdwolanie)
	if err != nil {
		return "", err
	}
	nazwaCelu := "treść bieżąca"
	if kod := strings.TrimSpace(wartoscTekstu(kodFinalnej)); kod != "" {
		nazwaCelu = kod
	}
	adnotacje, err := a.repozytorium.Komentarze(ctx, dokument.ID, "adnotacja")
	if err != nil {
		return "", bladStudio(err)
	}
	return trescRaportuRoznicyStudia(dokument, najstarsza.Kod, nazwaCelu,
		policzFragmentyRoznicy(baza, trescFinalna), adnotacje), nil
}

// trescRaportuRoznicyStudia składa raport w Markdownie — postaci wspólnej dla
// wszystkich formatów wyjściowych. PDF i DOCX powstają z tego samego tekstu,
// więc raport wydany dwa razy w dwóch formatach mówi to samo.
func trescRaportuRoznicyStudia(dokument dane.DokumentStudia, kodBazy, kodCelu string,
	fragmenty []shared.StudioDiffHunk, adnotacje []dane.KomentarzStudia) string {

	poFragmencie := map[int64][]dane.KomentarzStudia{}
	for _, adnotacja := range adnotacje {
		if adnotacja.FragmentNumer != nil {
			poFragmencie[*adnotacja.FragmentNumer] = append(poFragmencie[*adnotacja.FragmentNumer], adnotacja)
		}
	}

	dodane, usuniete, zmienione := 0, 0, 0
	for _, fragment := range fragmenty {
		switch fragment.Kind {
		case shared.DiffHunkKindAdded:
			dodane++
		case shared.DiffHunkKindRemoved:
			usuniete++
		case shared.DiffHunkKindChanged:
			zmienione++
		}
	}

	var raport strings.Builder
	tytul := dokument.Kod
	if dokument.Tytul != nil && strings.TrimSpace(*dokument.Tytul) != "" {
		tytul = *dokument.Tytul
	}
	raport.WriteString("# Raport zmian — " + tytul + "\n\n")
	raport.WriteString("Dokument: " + dokument.Kod + "\n")
	raport.WriteString("Wersja odniesienia: " + kodBazy + "\n")
	raport.WriteString("Wersja porownywana: " + kodCelu + "\n")
	raport.WriteString("Wydano: " + time.Now().UTC().Format(time.RFC3339) + "\n\n")
	raport.WriteString("## Statystyka\n\n")
	raport.WriteString("Fragmentow dodanych: " + strconv.Itoa(dodane) + "\n")
	raport.WriteString("Fragmentow usunietych: " + strconv.Itoa(usuniete) + "\n")
	raport.WriteString("Fragmentow zmienionych: " + strconv.Itoa(zmienione) + "\n\n")
	raport.WriteString("## Fragmenty\n\n")

	for _, fragment := range fragmenty {
		if fragment.Kind == shared.DiffHunkKindContext {
			continue
		}
		raport.WriteString("### Fragment " + strconv.Itoa(fragment.Index) +
			" — " + string(fragment.Kind) + "\n\n")
		if fragment.StartLine != nil && fragment.EndLine != nil {
			raport.WriteString("Wiersze " + strconv.Itoa(*fragment.StartLine) + "–" +
				strconv.Itoa(*fragment.EndLine) + "\n\n")
		}
		if fragment.Before != nil {
			raport.WriteString("Przed:\n\n" + *fragment.Before + "\n\n")
		}
		if fragment.After != nil {
			raport.WriteString("Po:\n\n" + *fragment.After + "\n\n")
		}
		for _, adnotacja := range poFragmencie[int64(fragment.Index)] {
			raport.WriteString("Adnotacja (" + adnotacja.Autor + "): " + adnotacja.Tresc + "\n\n")
		}
	}
	if dodane+usuniete+zmienione == 0 {
		raport.WriteString("Obie wersje maja te sama tresc — nie ma czego pokazac.\n")
	}
	return raport.String()
}

// adnotacjaWydaniaStudia to postać adnotacji odkładana w paczce, niezależna
// od struktury warstwy danych.
type adnotacjaWydaniaStudia struct {
	Kod        string `json:"kod"`
	Rodzaj     string `json:"rodzaj"`
	Autor      string `json:"autor"`
	Tresc      string `json:"tresc"`
	Fragment   *int64 `json:"fragment,omitempty"`
	ZakresOd   *int64 `json:"zakresOd,omitempty"`
	ZakresDo   *int64 `json:"zakresDo,omitempty"`
	Rozwiazany bool   `json:"rozwiazany"`
	Utworzono  string `json:"utworzono"`
}

// zlozAdnotacjeWydaniaStudia przekłada wiersze warstwy danych na postać
// adnotacji zapisywaną w archiwum.
func zlozAdnotacjeWydaniaStudia(wiersze []dane.KomentarzStudia) []adnotacjaWydaniaStudia {
	lista := make([]adnotacjaWydaniaStudia, 0, len(wiersze))
	for _, wiersz := range wiersze {
		lista = append(lista, adnotacjaWydaniaStudia{
			Kod: wiersz.Kod, Rodzaj: wiersz.Rodzaj, Autor: wiersz.Autor, Tresc: wiersz.Tresc,
			Fragment: wiersz.FragmentNumer, ZakresOd: wiersz.ZakresOd, ZakresDo: wiersz.ZakresDo,
			Rozwiazany: wiersz.Rozwiazany, Utworzono: wiersz.Utworzono,
		})
	}
	return lista
}

// adnotacjeIKomentarzeStudia zbiera OBA rodzaje wpisów redakcyjnych. Warstwa
// danych rozdziela je zapytaniem po rodzaju, a paczka przekazania niesie
// jedno i drugie: odbiorca ma zobaczyć całą rozmowę wokół dokumentu, nie jej
// połowę.
func (a *adapterStudia) adnotacjeIKomentarzeStudia(ctx context.Context,
	dokumentID int64) ([]dane.KomentarzStudia, error) {

	wszystkie := []dane.KomentarzStudia{}
	for _, rodzaj := range []string{"komentarz", "adnotacja"} {
		wiersze, err := a.repozytorium.Komentarze(ctx, dokumentID, rodzaj)
		if err != nil {
			return nil, bladStudio(err)
		}
		wszystkie = append(wszystkie, wiersze...)
	}
	return wszystkie, nil
}

// postacDokumentuStudia oddaje treść w formacie docelowym.
//
// Format nieznany jest odmową, nie zejściem na tekst: Operator, który poprosił
// o DOCX, a dostał plik tekstowy z rozszerzeniem `.docx`, dowiedziałby się
// o tym dopiero przy otwieraniu u odbiorcy.
func (a *adapterStudia) postacDokumentuStudia(tresc, format string,
	dokument dane.DokumentStudia, kartka geometriaStronyStudia) ([]byte, error) {

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "txt", "markdown", "md":
		return []byte(tresc), nil
	case "docx":
		return docxZTekstuStudia(tresc, dokument)
	case "pdf":
		naglowek := dokument.Kod
		if dokument.Tytul != nil && strings.TrimSpace(*dokument.Tytul) != "" {
			naglowek = *dokument.Tytul
		}
		karty, err := wyrysujStronyStudia(tresc,
			nastawyWyrysuStudia{naglowek: naglowek, geometria: kartka})
		if err != nil {
			return nil, err
		}
		strony := make([][]byte, 0, len(karty))
		for _, karta := range karty {
			bajty, err := pngZeStronyStudia(karta.obraz)
			if err != nil {
				return nil, err
			}
			strony = append(strony, bajty)
		}
		return pdfZeStronStudia(strony)
	default:
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"moduł Studio: format "+format+" nie jest formatem wydania; "+
				"dopuszczalne: pdf, docx, markdown, txt"))
	}
}

// ── Archiwa i dokument biurowy ──────────────────────────────────────────────

// formatArchiwumStudia rozstrzyga postać archiwum spośród wartości
// dopuszczalnych. Brak wskazania bierze zip, tak jak mówi kontrakt.
func formatArchiwumStudia(wskazanie *string) string {
	format := strings.ToLower(strings.TrimSpace(wartoscTekstu(wskazanie)))
	if format == "tar" {
		return "tar"
	}
	return "zip"
}

// spakujStudia składa archiwum wskazanej postaci, wybierając między zapisem
// w formacie ZIP a zapisem w formacie TAR.
func spakujStudia(wpisy []wpisArchiwumStudia, format string) ([]byte, error) {
	if format == "tar" {
		return tarStudia(wpisy)
	}
	return zipStudia(wpisy)
}

// zipStudia składa archiwum w formacie ZIP ze wskazanych wpisów, zachowując
// dla każdego znacznik czasu zapisu.
func zipStudia(wpisy []wpisArchiwumStudia) ([]byte, error) {
	var bufor bytes.Buffer
	pakowacz := zip.NewWriter(&bufor)
	for _, wpis := range wpisy {
		naglowek := &zip.FileHeader{Name: wpis.nazwa, Method: zip.Deflate}
		if !wpis.chwila.IsZero() {
			naglowek.Modified = wpis.chwila
		}
		ujscie, err := pakowacz.CreateHeader(naglowek)
		if err != nil {
			return nil, bladStudio(err)
		}
		if _, err := ujscie.Write(wpis.bajty); err != nil {
			return nil, bladStudio(err)
		}
	}
	if err := pakowacz.Close(); err != nil {
		return nil, bladStudio(err)
	}
	return bufor.Bytes(), nil
}

// tarStudia składa archiwum w formacie TAR ze wskazanych wpisów, zachowując
// dla każdego znacznik czasu zapisu.
func tarStudia(wpisy []wpisArchiwumStudia) ([]byte, error) {
	var bufor bytes.Buffer
	pakowacz := tar.NewWriter(&bufor)
	for _, wpis := range wpisy {
		naglowek := &tar.Header{
			Name: wpis.nazwa, Mode: 0o644, Size: int64(len(wpis.bajty)),
			ModTime: wpis.chwila, Typeflag: tar.TypeReg,
		}
		if naglowek.ModTime.IsZero() {
			naglowek.ModTime = time.Now().UTC()
		}
		if err := pakowacz.WriteHeader(naglowek); err != nil {
			return nil, bladStudio(err)
		}
		if _, err := pakowacz.Write(wpis.bajty); err != nil {
			return nil, bladStudio(err)
		}
	}
	if err := pakowacz.Close(); err != nil {
		return nil, bladStudio(err)
	}
	return bufor.Bytes(), nil
}

// wpisManifestuStudia składa plik manifest.json z opisu wydania i odkłada go
// jako wpis gotowy do archiwum.
func wpisManifestuStudia(manifest manifestWydaniaStudia) (wpisArchiwumStudia, error) {
	bajty, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return wpisArchiwumStudia{}, bladStudio(err)
	}
	return wpisArchiwumStudia{nazwa: "manifest.json", bajty: bajty, chwila: time.Now().UTC()}, nil
}

// nazwyWpisowStudia wymienia całą zawartość archiwum w manifeście, wypisując
// nazwę każdego wpisu z osobna.
func nazwyWpisowStudia(wpisy []wpisArchiwumStudia) []string {
	nazwy := make([]string, 0, len(wpisy))
	for _, wpis := range wpisy {
		nazwy = append(nazwy, wpis.nazwa)
	}
	return nazwy
}

// chwilaWpisuStudia przekłada znacznik czasu z warstwy bazy danych na czas
// wpisu archiwum w strefie UTC.
func chwilaWpisuStudia(znacznik string) time.Time {
	milisekundy := chwilaBazy(znacznik)
	if milisekundy == 0 {
		return time.Now().UTC()
	}
	return time.UnixMilli(milisekundy).UTC()
}

// wlaczone rozstrzyga pole logiczne nieobowiązkowe, którego brak znaczy „tak".
// Trzy pola paczki mają tę samą regułę, więc reguła stoi w jednym miejscu.
func wlaczone(wskazanie *bool) bool {
	return wskazanie == nil || *wskazanie
}

// docxZTekstuStudia składa najprostszy poprawny dokument DOCX z użyciem
// archive/zip biblioteki standardowej, bez stylów i tabel, gotowy do dalszej
// redakcji w dowolnym edytorze.
func docxZTekstuStudia(tresc string, dokument dane.DokumentStudia) ([]byte, error) {
	var akapity strings.Builder
	tytul := dokument.Kod
	if dokument.Tytul != nil && strings.TrimSpace(*dokument.Tytul) != "" {
		tytul = *dokument.Tytul
	}
	for _, wiersz := range strings.Split(tytul+"\n\n"+tresc, "\n") {
		akapity.WriteString(`<w:p><w:r><w:t xml:space="preserve">` +
			zaslonXmlStudia(wiersz) + `</w:t></w:r></w:p>`)
	}

	dokumentXml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		`<w:body>` + akapity.String() + `</w:body></w:document>`

	typyXml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
		`<Default Extension="xml" ContentType="application/xml"/>` +
		`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>` +
		`</Types>`

	powiazaniaXml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>` +
		`</Relationships>`

	chwila := time.Now().UTC()
	return zipStudia([]wpisArchiwumStudia{
		{nazwa: "[Content_Types].xml", bajty: []byte(typyXml), chwila: chwila},
		{nazwa: "_rels/.rels", bajty: []byte(powiazaniaXml), chwila: chwila},
		{nazwa: "word/document.xml", bajty: []byte(dokumentXml), chwila: chwila},
	})
}

// zaslonXmlStudia zasłania znaki, które w XML mają własne znaczenie. Bez tego
// dokument z ostrym nawiasem w treści byłby dokumentem, którego edytor nie
// otworzy.
func zaslonXmlStudia(tekst string) string {
	zamiennik := strings.NewReplacer(
		"&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&apos;")
	return zamiennik.Replace(tekst)
}
