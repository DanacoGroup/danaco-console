// Odpowiedzialność pliku: cztery czynności, które łączy to, że pracują na
// TREŚCI dokumentu, a nie na jego postaci — `studio.batch.run`,
// `studio.asset.embed`, `studio.search.semantic` i `studio.diff.source`.
package core

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// UruchomWsad obsługuje `studio.batch.run` — wykonuje seryjne polecenia
// studia nad wskazanym oknem produktu.
func (a *adapterStudia) UruchomWsad(ctx context.Context,
	z shared.StudioBatchRunRequest) (shared.StudioBatchRunResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" || strings.TrimSpace(z.ActionId) == "" {
		return shared.StudioBatchRunResponse{}, bladWskazaniaStudio(
			"wsad wymaga okna i pozycji rejestru akcji")
	}
	if len(z.DocumentIds) == 0 {
		return shared.StudioBatchRunResponse{}, bladWskazaniaStudio(
			"wsad bez ani jednego dokumentu")
	}

	kodPrzebiegu := nowyIdentyfikator(przedrostekWsaduStudia)
	odrzucone := []shared.StudioBatchRejection{}
	pozycje := make([]dane.PozycjaWsaduStudia, 0, len(z.DocumentIds))
	przyjete := 0

	for _, kodDokumentu := range z.DocumentIds {
		kod := strings.TrimSpace(kodDokumentu)
		if kod == "" {
			continue
		}
		odpowiedz, err := a.OperacjaKontekstowa(ctx, shared.StudioContextualOpRequest{
			WindowId: z.WindowId, DocumentId: kod, ActionId: z.ActionId,
		})
		if err != nil {
			powod := err.Error()
			odrzucone = append(odrzucone, shared.StudioBatchRejection{DocumentId: kod, Reason: powod})
			pozycje = append(pozycje, dane.PozycjaWsaduStudia{
				DokumentKod: kod, Stan: "odrzucona", Powod: &powod,
			})
			continue
		}
		przyjete++
		pozycje = append(pozycje, dane.PozycjaWsaduStudia{
			DokumentKod: kod, Stan: "przyjeta", PropozycjaKod: odpowiedz.ProposalId,
		})
	}

	przebieg := dane.PrzebiegWsaduStudia{
		Kod: kodPrzebiegu, Okno: z.WindowId, AkcjaID: z.ActionId,
		Przyjete: przyjete, Odrzucone: len(odrzucone),
	}
	if len(z.Params) > 0 {
		parametry := string(z.Params)
		przebieg.ParametryJSON = &parametry
	}
	if _, err := a.repozytorium.ZapiszPrzebiegWsadu(ctx, przebieg, pozycje); err != nil {
		return shared.StudioBatchRunResponse{}, bladStudio(err)
	}

	return shared.StudioBatchRunResponse{
		RunId: kodPrzebiegu, Accepted: przyjete, Rejected: odrzucone,
	}, nil
}

// OsadzZasob obsługuje `studio.asset.embed`. Osadzenie wstawia odwołanie do
// zasobu, nie jego bajty, zapisem Markdown.
func (a *adapterStudia) OsadzZasob(ctx context.Context,
	z shared.StudioAssetEmbedRequest) (shared.StudioAssetEmbedResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" || strings.TrimSpace(z.AssetId) == "" {
		return shared.StudioAssetEmbedResponse{}, bladWskazaniaStudio(
			"osadzenie wymaga dokumentu i zasobu")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioAssetEmbedResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	if a.zasoby == nil {
		return shared.StudioAssetEmbedResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"moduł Studio: osadzenie zasobu nie ma drogi — serwer złożony bez repozytorium "+
				"zasobów; naprawa: podpiąć magazyn zasobów przy składaniu serwera"))
	}
	zasob, err := a.zasoby.Zasob(ctx, strings.TrimSpace(z.AssetId))
	if err != nil {
		return shared.StudioAssetEmbedResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "moduł Studio: zasobu "+z.AssetId+" nie ma w magazynie serwera"))
	}

	tresc, err := a.trescZOdwolania(dokument.Tresc, dokument.TrescOdwolanie)
	if err != nil {
		return shared.StudioAssetEmbedResponse{}, err
	}
	wstawka := wstawkaZasobuStudia(zasob, z.AltText, z.Caption)

	// Położenie liczy się w znakach, a nie w bajtach.
	runy := []rune(tresc)
	miejsce := len(runy)
	if z.Position != nil {
		miejsce = *z.Position
		if miejsce < 0 || miejsce > len(runy) {
			return shared.StudioAssetEmbedResponse{}, bladWskazaniaStudio(
				"położenie wstawienia " + strconv.Itoa(miejsce) + " wychodzi poza treść dokumentu " +
					"o długości " + strconv.Itoa(len(runy)) + " znaków")
		}
	}
	nowa := string(runy[:miejsce]) + wstawka + string(runy[miejsce:])

	dokument.Tresc = &nowa
	dokument.TrescOdwolanie = nil
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return shared.StudioAssetEmbedResponse{}, bladStudio(err)
	}
	return shared.StudioAssetEmbedResponse{Document: a.zlozDokument(zapisany)}, nil
}

// wstawkaZasobuStudia składa zapis osadzenia. Podpis idzie osobnym wierszem pod
// obrazem, bo tytuł w nawiasie okrągłym pokazuje się dopiero po najechaniu, a
// podpis ma być widoczny zawsze.
func wstawkaZasobuStudia(zasob dane.ZasobDesignu, tekstAlternatywny, podpis *string) string {
	opis := zasob.Kod
	if tekst := strings.TrimSpace(wartoscTekstu(tekstAlternatywny)); tekst != "" {
		opis = tekst
	} else if zasob.Nazwa != nil && strings.TrimSpace(*zasob.Nazwa) != "" {
		opis = *zasob.Nazwa
	}
	// Odwołanie w treści wskazuje zasób, nie ścieżkę na dysku rdzenia.
	wstawka := "\n\n![" + opis + "](danaco://zasob/" + zasob.Kod + ")\n"
	if tekst := strings.TrimSpace(wartoscTekstu(podpis)); tekst != "" {
		wstawka += "\n" + tekst + "\n"
	}
	return wstawka
}

// PorownajZeZrodlem obsługuje `studio.diff.source`. Materiał wejściowy
// wskazuje się plikiem repozytorium albo zasobem magazynu.
func (a *adapterStudia) PorownajZeZrodlem(ctx context.Context,
	z shared.StudioDiffSourceRequest) (shared.StudioDiffSourceResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" {
		return shared.StudioDiffSourceResponse{}, bladWskazaniaStudio(
			"zestawienie ze źródłem bez wskazania dokumentu")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDiffSourceResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	robocza, err := a.trescWersjiStudia(ctx, z.VersionId, dokument)
	if err != nil {
		return shared.StudioDiffSourceResponse{}, err
	}

	zrodlo, odczytane := a.trescMaterialuWejsciowegoStudia(ctx, z, dokument)
	if !odczytane {
		return shared.StudioDiffSourceResponse{
			Hunks: []shared.StudioDiffHunk{}, SourceResolved: false,
		}, nil
	}
	return shared.StudioDiffSourceResponse{
		Hunks: policzFragmentyRoznicy(zrodlo, robocza), SourceResolved: true,
	}, nil
}

// trescMaterialuWejsciowegoStudia czyta materiał wejściowy trzema drogami po
// kolei: zasób magazynu, plik biblioteki, plik zapamiętany przy otwarciu
// dokumentu. Nieudany odczyt nie jest odmową — kontrakt ma na to własne pole.
func (a *adapterStudia) trescMaterialuWejsciowegoStudia(ctx context.Context,
	z shared.StudioDiffSourceRequest, dokument dane.DokumentStudia) (string, bool) {

	if kod := strings.TrimSpace(wartoscTekstu(z.SourceAssetId)); kod != "" {
		bajty, err := a.bajtyZasobuStudia(ctx, kod)
		if err != nil {
			return "", false
		}
		return string(bajty), true
	}

	kodPliku := strings.TrimSpace(wartoscTekstu(z.SourceLibraryFileId))
	if kodPliku == "" && dokument.PlikRepozytoriumID != nil {
		kodPliku = strings.TrimSpace(*dokument.PlikRepozytoriumID)
	}
	if kodPliku == "" || a.biblioteka == nil {
		return "", false
	}
	plik, err := a.biblioteka.Plik(ctx, kodPliku)
	if err != nil || plik.TrescOdwolanie == nil {
		return "", false
	}
	tresc, err := a.trescZOdwolania(nil, plik.TrescOdwolanie)
	if err != nil {
		return "", false
	}
	return tresc, true
}

// Wyszukiwanie znaczeniowe liczy bliskość fragmentów dokumentu do zapytania,
// modelem językowym albo rdzeniem. Nazwy dróg, którymi liczy się bliskość,
// wychodzą kontraktem w polu `mode`.
const (
	drogaSemantykiModelem = "kanal-modelu"
	drogaSemantykiMiara   = "miara-serwera"
)

// granicaFragmentowSemantyki chroni kanał modelu przed pytaniem o dokument,
// którego i tak nie zdąży przeczytać w jednym wywołaniu.
const granicaFragmentowSemantyki = 400

// WyszukajZnaczeniowo obsługuje `studio.search.semantic`, drogą kanału modelu
// okna albo miarą arytmetyczną w rdzeniu.
func (a *adapterStudia) WyszukajZnaczeniowo(ctx context.Context,
	z shared.StudioSearchSemanticRequest) (shared.StudioSearchSemanticResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" || strings.TrimSpace(z.Query) == "" {
		return shared.StudioSearchSemanticResponse{}, bladWskazaniaStudio(
			"wyszukiwanie znaczeniowe wymaga dokumentu i zapytania")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioSearchSemanticResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	tresc, err := a.trescWersjiStudia(ctx, z.VersionId, dokument)
	if err != nil {
		return shared.StudioSearchSemanticResponse{}, err
	}

	fragmenty := fragmentyTresciStudia(tresc)
	if len(fragmenty) == 0 {
		return shared.StudioSearchSemanticResponse{
			Matches: []shared.StudioSemanticMatch{}, Mode: drogaSemantykiMiara,
		}, nil
	}
	if len(fragmenty) > granicaFragmentowSemantyki {
		fragmenty = fragmenty[:granicaFragmentowSemantyki]
	}

	droga := drogaSemantykiMiara
	oceny, policzylModel := a.ocenyModeluStudia(ctx, dokument, z.Query, fragmenty)
	if policzylModel {
		droga = drogaSemantykiModelem
	} else {
		oceny = ocenyMiaryStudia(z.Query, fragmenty)
	}

	prog := 0.0
	if z.MinScore != nil {
		prog = *z.MinScore
	}
	trafienia := []shared.StudioSemanticMatch{}
	for i, fragment := range fragmenty {
		if oceny[i] < prog || oceny[i] <= 0 {
			continue
		}
		wiersz := fragment.wiersz
		trafienia = append(trafienia, shared.StudioSemanticMatch{
			Text: fragment.tekst, RangeStart: fragment.od, RangeEnd: fragment.do_,
			Score: oceny[i], Line: &wiersz,
		})
	}
	sort.SliceStable(trafienia, func(i, j int) bool { return trafienia[i].Score > trafienia[j].Score })
	if z.Limit != nil && *z.Limit >= 0 && *z.Limit < len(trafienia) {
		trafienia = trafienia[:*z.Limit]
	}
	return shared.StudioSearchSemanticResponse{Matches: trafienia, Mode: droga}, nil
}

// fragmentSemantykiStudia to jeden akapit treści wraz z jego położeniem
// w dokumencie źródłowym studia.
type fragmentSemantykiStudia struct {
	tekst   string
	od, do_ int
	wiersz  int
}

// fragmentyTresciStudia dzieli treść na akapity. Akapit, a nie zdanie: pytanie
// znaczeniowe pada o myśl, a myśl w dokumencie redakcyjnym mieści się
// w akapicie; zdanie wyrwane z akapitu bywa nierozstrzygalne bez sąsiadów.
func fragmentyTresciStudia(tresc string) []fragmentSemantykiStudia {
	fragmenty := []fragmentSemantykiStudia{}
	runy := []rune(tresc)
	poczatek, wiersz, wierszPoczatku := 0, 1, 1
	dodaj := func(koniec int) {
		tekst := strings.TrimSpace(string(runy[poczatek:koniec]))
		if tekst != "" {
			fragmenty = append(fragmenty, fragmentSemantykiStudia{
				tekst: tekst, od: poczatek, do_: koniec, wiersz: wierszPoczatku,
			})
		}
	}
	for i := 0; i < len(runy); i++ {
		if runy[i] != '\n' {
			continue
		}
		wiersz++
		if i+1 < len(runy) && runy[i+1] == '\n' {
			dodaj(i)
			i++
			wiersz++
			poczatek, wierszPoczatku = i+1, wiersz
		}
	}
	dodaj(len(runy))
	return fragmenty
}

// odpowiedzModeluStudia to postać, w której model językowy oddaje ocenę
// dopasowania fragmentów zapytaniu.
type odpowiedzModeluStudia struct {
	Fragmenty []struct {
		Numer int     `json:"numer"`
		Ocena float64 `json:"ocena"`
	} `json:"fragmenty"`
}

// ocenyModeluStudia pyta kanał modelu okna o bliskość każdego fragmentu,
// wracając do miary arytmetycznej, gdy odpowiedzi nie da się wziąć za prawdę.
func (a *adapterStudia) ocenyModeluStudia(ctx context.Context, dokument dane.DokumentStudia,
	zapytanie string, fragmenty []fragmentSemantykiStudia) ([]float64, bool) {

	if a.kanaly == nil || a.okna == nil || strings.TrimSpace(dokument.Okno) == "" {
		return nil, false
	}
	okno, err := a.okna.Okno(dokument.Okno)
	if err != nil || okno.KanalModelu == "" {
		return nil, false
	}

	var polecenie strings.Builder
	polecenie.WriteString("Ocen, jak blisko znaczeniowo kazdy fragment jest zapytaniu.\n")
	polecenie.WriteString("Odpowiedz WYLACZNIE dokumentem JSON postaci ")
	polecenie.WriteString(`{"fragmenty":[{"numer":1,"ocena":0.0}]}`)
	polecenie.WriteString(", gdzie ocena jest liczba od 0 do 1.\n\nZapytanie: ")
	polecenie.WriteString(zapytanie)
	polecenie.WriteString("\n\nFragmenty:\n")
	for i, fragment := range fragmenty {
		polecenie.WriteString(strconv.Itoa(i+1) + ". " + fragment.tekst + "\n")
	}

	var odpowiedz strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			odpowiedz.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	err = a.kanaly.Wyslij(ctx, models.Zapytanie{
		Zasiegi:   models.Zasiegi{Okno: dokument.Okno},
		Wiadomosc: "studio.search.semantic",
		Tresc:     polecenie.String(),
		Kanal:     okno.KanalModelu,
	}, ujscie)
	if err != nil || strings.TrimSpace(odpowiedz.String()) == "" {
		return nil, false
	}

	surowa := wytnijJsonStudia(odpowiedz.String())
	var wynik odpowiedzModeluStudia
	if surowa == "" || json.Unmarshal([]byte(surowa), &wynik) != nil || len(wynik.Fragmenty) == 0 {
		return nil, false
	}
	oceny := make([]float64, len(fragmenty))
	for _, pozycja := range wynik.Fragmenty {
		if pozycja.Numer < 1 || pozycja.Numer > len(fragmenty) {
			continue
		}
		ocena := pozycja.Ocena
		if ocena < 0 {
			ocena = 0
		}
		if ocena > 1 {
			ocena = 1
		}
		oceny[pozycja.Numer-1] = ocena
	}
	return oceny, true
}

// wytnijJsonStudia wyjmuje dokument JSON z odpowiedzi modelu obudowanej
// zdaniem wstępnym albo ogrodzeniem bloku kodu.
func wytnijJsonStudia(tekst string) string {
	poczatek := strings.Index(tekst, "{")
	koniec := strings.LastIndex(tekst, "}")
	if poczatek < 0 || koniec <= poczatek {
		return ""
	}
	return tekst[poczatek : koniec+1]
}

// slowaNieznaczaceStudia to słowa, które w polszczyźnie występują wszędzie
// i o bliskości znaczeniowej nie mówią nic.
var slowaNieznaczaceStudia = map[string]bool{
	"aby": true, "albo": true, "ale": true, "bez": true, "byc": true, "być": true,
	"czy": true, "dla": true, "gdy": true, "jak": true, "jest": true, "jako": true,
	"kiedy": true, "lub": true, "nad": true, "nie": true, "oraz": true, "pod": true,
	"przez": true, "przy": true, "tak": true, "tego": true, "tej": true, "tym": true,
	"wiec": true, "więc": true, "zeby": true, "żeby": true,
}

// ocenyMiaryStudia liczy bliskość arytmetycznie: zbieżność słów znaczących
// ważona rzadkością słowa w dokumencie.
func ocenyMiaryStudia(zapytanie string, fragmenty []fragmentSemantykiStudia) []float64 {
	slowaZapytania := slowaZnaczaceStudia(zapytanie)
	oceny := make([]float64, len(fragmenty))
	if len(slowaZapytania) == 0 {
		return oceny
	}

	wystapienia := map[string]int{}
	slowaFragmentow := make([]map[string]int, len(fragmenty))
	for i, fragment := range fragmenty {
		policzone := map[string]int{}
		for _, slowo := range slowaZnaczaceStudia(fragment.tekst) {
			policzone[slowo]++
		}
		slowaFragmentow[i] = policzone
		for slowo := range policzone {
			wystapienia[slowo]++
		}
	}

	waga := func(slowo string) float64 {
		return math.Log(1 + float64(len(fragmenty))/float64(1+wystapienia[slowo]))
	}
	wagaZapytania := 0.0
	for _, slowo := range slowaZapytania {
		wagaZapytania += waga(slowo)
	}
	if wagaZapytania == 0 {
		return oceny
	}

	for i, policzone := range slowaFragmentow {
		suma := 0.0
		for _, slowo := range slowaZapytania {
			if policzone[slowo] > 0 {
				suma += waga(slowo)
			}
		}
		oceny[i] = suma / wagaZapytania
	}
	return oceny
}

// slowaZnaczaceStudia rozbija tekst na słowa nadające się do porównania,
// pomijając wielkość liter i znaki niebędące literami ani cyframi.
func slowaZnaczaceStudia(tekst string) []string {
	slowa := []string{}
	for _, slowo := range strings.FieldsFunc(strings.ToLower(tekst), func(znak rune) bool {
		return !unicode.IsLetter(znak) && !unicode.IsDigit(znak)
	}) {
		if len([]rune(slowo)) < 3 || slowaNieznaczaceStudia[slowo] {
			continue
		}
		slowa = append(slowa, slowo)
	}
	return slowa
}
