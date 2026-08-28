// Odpowiedzialność pliku: materiały projektu — wydobycie tekstu z pliku oraz wykrywanie i scalanie duplikatów treści w bibliotece workspace.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// WydobycieTekstuDokumentu jest portem warsztatu dokumentów widzianym przez
// moduł Workspace. Jedna metoda, bo jedna czynność: przeczytać plik.
type WydobycieTekstuDokumentu interface {
	WyciagnijTekst(ctx context.Context,
		z shared.DocumentTextExtractRequest) (shared.DocumentTextExtractResponse, error)
}

// ZDokumentamiWorkspace podpina warsztat dokumentów. Bez niego wydobycie tekstu
// odmawia zdaniem nazywającym brak ogniwa, a reszta modułu pracuje dalej.
func (a *adapterPrzestrzeniRoboczej) ZDokumentamiWorkspace(warsztat WydobycieTekstuDokumentu) *adapterPrzestrzeniRoboczej {
	a.dokumenty = warsztat
	return a
}

// WydobadzTekst obsługuje workspace.library.text.extract, wydobywając tekst warstwą tekstową dokumentu albo rozpoznaniem pisma z pikseli.
func (a *adapterPrzestrzeniRoboczej) WydobadzTekst(ctx context.Context,
	z shared.WorkspaceLibraryTextExtractRequest) (shared.WorkspaceLibraryTextExtractResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceLibraryTextExtractResponse{}, err
	}
	if strings.TrimSpace(z.FileId) == "" {
		return shared.WorkspaceLibraryTextExtractResponse{},
			bladProjektu("wydobycie bez wskazania pliku")
	}
	sciezka, err := a.sciezkaPlikuProjektuWorkspace(projekt.Kod, z.FileId)
	if err != nil {
		return shared.WorkspaceLibraryTextExtractResponse{}, err
	}
	// Gotowy wyciąg wraca bez powtórnego czytania: rozpoznanie pisma ze skanu trwa minuty.
	if z.Force == nil || !*z.Force {
		if gotowy, err := a.repozytorium.WyciagWorkspace(ctx, projekt.ID, z.FileId); err == nil {
			return shared.WorkspaceLibraryTextExtractResponse{
				Extraction: wyciagKontraktuWorkspace(gotowy),
			}, nil
		}
	}
	if a.dokumenty == nil {
		return shared.WorkspaceLibraryTextExtractResponse{},
			bladProjektu("wydobycie tekstu nie ma warsztatu dokumentów — rdzeń zmontowano bez tego ogniwa")
	}
	sposob := sposobWydobyciaWorkspace(z.Method, sciezka)
	jezyki := z.Languages
	if len(jezyki) == 0 {
		jezyki = []string{"pol", "eng"}
	}
	zadanie := shared.DocumentTextExtractRequest{
		SourcePath: &sciezka,
		Language:   tekstOpcjonalnyWorkspace(strings.Join(jezyki, "+")),
	}
	if sposob == shared.WorkspaceExtractionMethodOcr {
		wymus := true
		zadanie.ForceOcr = &wymus
	}
	odpowiedz, err := a.dokumenty.WyciagnijTekst(ctx, zadanie)
	if err != nil {
		return shared.WorkspaceLibraryTextExtractResponse{}, err
	}
	if odpowiedz.UsedOcr {
		sposob = shared.WorkspaceExtractionMethodOcr
	} else if sposob == shared.WorkspaceExtractionMethodOcr {
		sposob = sposobWydobyciaWorkspace(nil, sciezka)
	}
	uzyteJezyki := []string{}
	if sposob == shared.WorkspaceExtractionMethodOcr {
		uzyteJezyki = jezyki
	}
	wyciag := dane.WyciagWorkspace{
		ProjektID: projekt.ID, Plik: z.FileId, Sposob: sposob, Tresc: odpowiedz.Text,
		LiczbaZnakow: len([]rune(odpowiedz.Text)), Jezyki: uzyteJezyki,
	}
	if err := a.repozytorium.ZapiszWyciagWorkspace(ctx, wyciag); err != nil {
		return shared.WorkspaceLibraryTextExtractResponse{}, err
	}
	zapisany, err := a.repozytorium.WyciagWorkspace(ctx, projekt.ID, z.FileId)
	if err != nil {
		return shared.WorkspaceLibraryTextExtractResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindFileChanged,
		shared.ChangeKindUpdated, shared.WorkspaceEntityKindFile, z.FileId,
		"wydobyto treść pliku "+z.FileId+" sposobem "+string(sposob))

	return shared.WorkspaceLibraryTextExtractResponse{
		Extraction: wyciagKontraktuWorkspace(zapisany),
	}, nil
}

// Duplikaty obsługuje workspace.library.duplicate.list, wskazując grupy plików o identycznej treści wedlug sumy kontrolnej.
func (a *adapterPrzestrzeniRoboczej) Duplikaty(ctx context.Context,
	z shared.WorkspaceLibraryDuplicateListRequest) (shared.WorkspaceLibraryDuplicateListResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceLibraryDuplicateListResponse{}, err
	}
	granica := int64(0)
	if z.MinSizeBytes != nil {
		granica = *z.MinSizeBytes
	}
	pliki := a.plikiProjektu(projekt.Kod, "")
	wedlugSumy := map[string][]shared.LibraryFile{}
	rozmiary := map[string]int64{}
	for _, plik := range pliki {
		if plik.SizeBytes == nil || *plik.SizeBytes < granica {
			continue
		}
		suma, err := sumaPlikuWorkspace(a.katalogProjektu(projekt.Kod), plik.Id)
		if err != nil {
			continue
		}
		wpisany := plik
		wpisany.Checksum = &suma
		wedlugSumy[suma] = append(wedlugSumy[suma], wpisany)
		rozmiary[suma] = *plik.SizeBytes
	}
	grupy := []shared.WorkspaceDuplicateGroup{}
	doOdzyskania := int64(0)
	for suma, wykaz := range wedlugSumy {
		if len(wykaz) < 2 {
			continue
		}
		sort.SliceStable(wykaz, func(i, j int) bool { return wykaz[i].CreatedAt < wykaz[j].CreatedAt })
		zachowany := wykaz[0].Id
		grupy = append(grupy, shared.WorkspaceDuplicateGroup{
			Checksum: suma, SizeBytes: rozmiary[suma], Files: wykaz,
			SuggestedKeepFileId: &zachowany,
		})
		doOdzyskania += rozmiary[suma] * int64(len(wykaz)-1)
	}
	sort.SliceStable(grupy, func(i, j int) bool { return grupy[i].Checksum < grupy[j].Checksum })
	return shared.WorkspaceLibraryDuplicateListResponse{
		Groups: grupy, Total: len(grupy), ReclaimableBytes: &doOdzyskania,
	}, nil
}

// ScalDuplikaty obsługuje `workspace.library.duplicate.merge`.
//
// Scalenie usuwa pliki wskazane jako scalane i zdejmuje ich wyciągi treści ze
// wskaźnika wyszukiwania — zostawiony wyciąg wskazywałby plik, którego nie ma.
func (a *adapterPrzestrzeniRoboczej) ScalDuplikaty(ctx context.Context,
	z shared.WorkspaceLibraryDuplicateMergeRequest) (shared.WorkspaceLibraryDuplicateMergeResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceLibraryDuplicateMergeResponse{}, err
	}
	if strings.TrimSpace(z.KeepFileId) == "" || len(z.MergedFileIds) == 0 {
		return shared.WorkspaceLibraryDuplicateMergeResponse{},
			bladProjektu("scalenie bez wskazania pliku zachowywanego albo plików scalanych")
	}
	korzen := a.katalogProjektu(projekt.Kod)
	if _, err := a.sciezkaPlikuProjektuWorkspace(projekt.Kod, z.KeepFileId); err != nil {
		return shared.WorkspaceLibraryDuplicateMergeResponse{}, err
	}
	sumaZachowanego, err := sumaPlikuWorkspace(korzen, z.KeepFileId)
	if err != nil {
		return shared.WorkspaceLibraryDuplicateMergeResponse{},
			bladProjektu("pliku zachowywanego nie da się odczytać: " + err.Error())
	}
	scalone, odzyskane := 0, int64(0)
	for _, kod := range z.MergedFileIds {
		if kod == z.KeepFileId {
			continue
		}
		sciezka, err := a.sciezkaPlikuProjektuWorkspace(projekt.Kod, kod)
		if err != nil {
			continue
		}
		suma, err := sumaPlikuWorkspace(korzen, kod)
		if err != nil || suma != sumaZachowanego {
			// Plik o innej treści nie jest duplikatem; jego usunięcie byłoby utratą materiału.
			continue
		}
		opis, err := os.Stat(sciezka)
		if err != nil {
			continue
		}
		if err := os.Remove(sciezka); err != nil {
			continue
		}
		_ = a.repozytorium.UsunWyciagWorkspace(ctx, projekt.ID, kod)
		odzyskane += opis.Size()
		scalone++
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.WorkspaceLibraryDuplicateMergeResponse{}, err
	}
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindFileChanged,
		shared.ChangeKindDeleted, shared.WorkspaceEntityKindFile, z.KeepFileId,
		"scalono duplikaty przy pliku "+z.KeepFileId)

	zachowany := shared.LibraryFile{}
	for _, plik := range a.plikiProjektu(projekt.Kod, "") {
		if plik.Id == z.KeepFileId {
			plik.Checksum = &sumaZachowanego
			zachowany = plik
		}
	}
	return shared.WorkspaceLibraryDuplicateMergeResponse{
		File: zachowany, Merged: scalone, ReclaimedBytes: &odzyskane,
	}, nil
}

// sciezkaPlikuProjektuWorkspace rozwiązuje identyfikator pliku (ścieżkę
// względną katalogu projektu) na ścieżkę na dysku, pilnując, żeby wskazanie nie
// wyszło poza katalog projektu.
func (a *adapterPrzestrzeniRoboczej) sciezkaPlikuProjektuWorkspace(kodProjektu,
	plik string) (string, error) {

	korzen := a.katalogProjektu(kodProjektu)
	if korzen == "" {
		return "", bladProjektu("projekt nie ma katalogu roboczego — nie ma czego czytać")
	}
	pelna := filepath.Join(korzen, filepath.FromSlash(plik))
	wzgledna, err := filepath.Rel(korzen, pelna)
	if err != nil || strings.HasPrefix(wzgledna, "..") {
		return "", bladProjektu("wskazanie pliku wychodzi poza katalog projektu")
	}
	if _, err := os.Stat(pelna); err != nil {
		return "", bladProjektu("pliku " + plik + " nie ma w bibliotece projektu")
	}
	return pelna, nil
}

// sumaPlikuWorkspace liczy sumę kontrolną treści pliku projektu biblioteką standardową, bez programu z zewnątrz.
func sumaPlikuWorkspace(korzen, plik string) (string, error) {
	bajty, err := os.ReadFile(filepath.Join(korzen, filepath.FromSlash(plik)))
	if err != nil {
		return "", err
	}
	suma := sha256.Sum256(bajty)
	return hex.EncodeToString(suma[:]), nil
}

// sposobWydobyciaWorkspace ustala sposób wydobycia: wskazanie żądania, a przy jego braku rodzaj treści pliku.
func sposobWydobyciaWorkspace(wskazany *shared.WorkspaceExtractionMethod,
	sciezka string) shared.WorkspaceExtractionMethod {

	if wskazany != nil && *wskazany != "" {
		return *wskazany
	}
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".pdf":
		return shared.WorkspaceExtractionMethodPdf
	case ".docx":
		return shared.WorkspaceExtractionMethodDocx
	case ".png", ".jpg", ".jpeg", ".tif", ".tiff", ".bmp", ".webp":
		return shared.WorkspaceExtractionMethodOcr
	default:
		return shared.WorkspaceExtractionMethodText
	}
}

// wyciagKontraktuWorkspace przekłada wiersz wyciągu tekstu na byt kontraktu wymiany z klientem workspace.
func wyciagKontraktuWorkspace(w dane.WyciagWorkspace) shared.WorkspaceTextExtraction {
	wyciag := shared.WorkspaceTextExtraction{
		FileId: w.Plik, Method: w.Sposob, CharacterCount: w.LiczbaZnakow,
		Languages: w.Jezyki, Indexed: true, ExtractedAt: chwilaBazy(w.Utworzono),
	}
	if poczatek := skrocOpisWorkspace(w.Tresc, 400); poczatek != "" {
		wyciag.Excerpt = &poczatek
	}
	return wyciag
}
