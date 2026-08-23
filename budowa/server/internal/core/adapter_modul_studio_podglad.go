// Odpowiedzialność pliku: Preview Window i różnica wizualna —
// `studio.preview.render` oraz `studio.diff.visual`. Silnik wyrysu stoi
// w `adapter_studio_wyrys.go`; ten plik odpowiada wyłącznie za to, skąd wziąć
// treść, jakimi nastawami ją wyrysować i gdzie odłożyć wynik.
//
// ── Podgląd oddaje STRONY, nie tekst ────────────────────────────────────────
// Kontrakt obu komend mówi o zasobach: `pageAssetIds` i `overlayAssetIds`.
// Powód stoi w opracowaniu (rozdz. 3.7): podgląd ma pokazać UKŁAD — typografię,
// paginację, nagłówek i stopkę — a nie ten sam tekst, który Operator widzi
// w edytorze. Strona wychodzi więc obrazem, bo obraz jest jedyną postacią,
// w której układ da się zobaczyć bez drugiego silnika składu po stronie okna.
package core

import (
	"context"
	"encoding/json"
	"image"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ustawieniaProfiluStudia to nastawy strony zapamiętane przy profilu wydania.
// Jadą przez JSON, bo profil jest wpisem katalogowym o jednym polu ładunku —
// tym samym, którym jedzie prompt operacji i kroki łańcucha.
type ustawieniaProfiluStudia struct {
	Naglowek  string `json:"naglowek,omitempty"`
	Stopka    string `json:"stopka,omitempty"`
	ZnakWodny string `json:"znakWodny,omitempty"`
}

// WyrenderujPodglad obsługuje `studio.preview.render`.
//
// Format docelowy rozstrzyga o postaci wyniku dwustopniowo: strony powstają
// zawsze jako obrazy (bo podgląd pokazuje układ), a przy formacie `pdf` z tych
// samych stron składa się dodatkowo dokument. Dwa różne silniki — jeden dla
// podglądu, drugi dla wydania — dawałyby dwa różne układy, a wtedy podgląd
// przestaje być podglądem.
func (a *adapterStudia) WyrenderujPodglad(ctx context.Context,
	z shared.StudioPreviewRenderRequest) (shared.StudioPreviewRenderResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" {
		return shared.StudioPreviewRenderResponse{}, bladWskazaniaStudio(
			"render podglądu bez wskazania dokumentu")
	}
	if strings.TrimSpace(z.Format) == "" {
		return shared.StudioPreviewRenderResponse{}, bladWskazaniaStudio(
			"render podglądu bez wskazania formatu docelowego")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPreviewRenderResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	tresc, err := a.trescWersjiStudia(ctx, z.VersionId, dokument)
	if err != nil {
		return shared.StudioPreviewRenderResponse{}, err
	}

	nastawy, err := a.nastawyPodgladuStudia(ctx, dokument, z)
	if err != nil {
		return shared.StudioPreviewRenderResponse{}, err
	}
	karty, err := wyrysujStronyStudia(tresc, nastawy)
	if err != nil {
		return shared.StudioPreviewRenderResponse{}, err
	}

	okno := wartoscTekstu(z.WindowId)
	kody := make([]string, 0, len(karty))
	strony := make([][]byte, 0, len(karty))
	for _, karta := range karty {
		bajty, err := pngZeStronyStudia(karta.obraz)
		if err != nil {
			return shared.StudioPreviewRenderResponse{}, err
		}
		strony = append(strony, bajty)
		zasob, err := a.odlozObrazStudia(ctx, bajty,
			"podglad-"+dokument.Kod+"-strona-"+strconv.Itoa(karta.numer)+".png", okno,
			karta.obraz.Bounds().Dx(), karta.obraz.Bounds().Dy())
		if err != nil {
			return shared.StudioPreviewRenderResponse{}, err
		}
		kody = append(kody, zasob.Id)
	}

	// Dokument w formacie docelowym powstaje OBOK stron i dokłada się na końcu
	// wykazu — Preview Window pobiera go jako całość, a strony pokazuje jedna
	// po drugiej. Wykaz bez niego zmuszałby okno do drugiego żądania po to
	// samo, co rdzeń ma już złożone.
	if strings.EqualFold(strings.TrimSpace(z.Format), "pdf") {
		bajty, err := pdfZeStronStudia(strony)
		if err != nil {
			return shared.StudioPreviewRenderResponse{}, err
		}
		zasob, err := a.odlozTrescStudia(ctx, bajty, "podglad-"+dokument.Kod+".pdf", "pdf", okno)
		if err != nil {
			return shared.StudioPreviewRenderResponse{}, err
		}
		kody = append(kody, zasob.Id)
	}

	return shared.StudioPreviewRenderResponse{PageAssetIds: kody, Pages: len(karty)}, nil
}

// PorownajWizualnie obsługuje `studio.diff.visual`.
//
// Porównanie idzie po WYRYSIE obu wersji, nie po ich tekście: sedno tej
// czynności jest w tym, żeby zobaczyć zmianę, której różnica tekstowa nie widzi
// — przesunięcie akapitu na następną stronę, zmianę łamania, przestawienie
// nagłówka. Obie strony rysuje ten sam silnik tymi samymi nastawami, więc
// różnica pikseli jest różnicą treści, a nie różnicą sposobu rysowania.
func (a *adapterStudia) PorownajWizualnie(ctx context.Context,
	z shared.StudioDiffVisualRequest) (shared.StudioDiffVisualResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" || strings.TrimSpace(z.BaseVersionId) == "" ||
		strings.TrimSpace(z.TargetVersionId) == "" {
		return shared.StudioDiffVisualResponse{}, bladWskazaniaStudio(
			"porównanie wizualne wymaga dokumentu i obu wersji")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDiffVisualResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	baza, err := a.trescWersjiPoKodzie(ctx, z.BaseVersionId)
	if err != nil {
		return shared.StudioDiffVisualResponse{}, err
	}
	cel, err := a.trescWersjiPoKodzie(ctx, z.TargetVersionId)
	if err != nil {
		return shared.StudioDiffVisualResponse{}, err
	}

	nastawy := nastawyWyrysuStudia{
		naglowek: tytulDokumentuStudia(dokument),
		stronaOd: stronaZeWskaznikaStudia(z.PageFrom), stronaDo: stronaZeWskaznikaStudia(z.PageTo),
		geometria: a.geometriaDokumentuStudia(ctx, dokument.Kod),
	}
	kartyBazy, err := wyrysujStronyStudia(baza, nastawy)
	if err != nil {
		return shared.StudioDiffVisualResponse{}, err
	}
	kartyCelu, err := wyrysujStronyStudia(cel, nastawy)
	if err != nil {
		return shared.StudioDiffVisualResponse{}, err
	}

	okno := wartoscTekstu(z.WindowId)
	obszary := []shared.StudioVisualDiffRegion{}
	nakladki := []string{}
	// Wersje bywają różnej długości. Strona, której druga wersja nie ma, jest
	// zmianą CAŁEJ strony — i tak ją tu widać, bo za brakującą stronę wchodzi
	// pusta biała karta o tych samych wymiarach.
	dluzsza := len(kartyBazy)
	if len(kartyCelu) > dluzsza {
		dluzsza = len(kartyCelu)
	}
	for i := 0; i < dluzsza; i++ {
		stronaBazy := kartaAlboPustaStudia(kartyBazy, i, nastawy.kartka())
		stronaCelu := kartaAlboPustaStudia(kartyCelu, i, nastawy.kartka())
		numer := i + 1
		if i < len(kartyCelu) {
			numer = kartyCelu[i].numer
		} else if i < len(kartyBazy) {
			numer = kartyBazy[i].numer
		}

		dopasowana := dopasujStroneStudia(stronaBazy, stronaCelu.Bounds())
		stroneObszary := obszaryRoznicyStudia(dopasowana, stronaCelu, numer)
		if len(stroneObszary) == 0 {
			continue
		}
		obszary = append(obszary, stroneObszary...)

		bajty, err := pngZeStronyStudia(nakladkaRoznicyStudia(stronaCelu, stroneObszary))
		if err != nil {
			return shared.StudioDiffVisualResponse{}, err
		}
		zasob, err := a.odlozObrazStudia(ctx, bajty,
			"roznica-"+dokument.Kod+"-strona-"+strconv.Itoa(numer)+".png", okno,
			stronaCelu.Bounds().Dx(), stronaCelu.Bounds().Dy())
		if err != nil {
			return shared.StudioDiffVisualResponse{}, err
		}
		nakladki = append(nakladki, zasob.Id)
	}

	return shared.StudioDiffVisualResponse{Regions: obszary, OverlayAssetIds: nakladki}, nil
}

// nastawyPodgladuStudia składa nastawy strony: z profilu wydania, gdy Operator
// go wskazał, a znak wodny — z żądania, bo należy do podglądu, nie do profilu.
func (a *adapterStudia) nastawyPodgladuStudia(ctx context.Context, dokument dane.DokumentStudia,
	z shared.StudioPreviewRenderRequest) (nastawyWyrysuStudia, error) {

	nastawy := nastawyWyrysuStudia{
		naglowek:  tytulDokumentuStudia(dokument),
		znakWodny: strings.TrimSpace(wartoscTekstu(z.Watermark)),
		stronaOd:  stronaZeWskaznikaStudia(z.PageFrom),
		stronaDo:  stronaZeWskaznikaStudia(z.PageTo),
		geometria: a.geometriaDokumentuStudia(ctx, dokument.Kod),
	}
	kodProfilu := strings.TrimSpace(wartoscTekstu(z.ProfileId))
	if kodProfilu == "" {
		return nastawy, nil
	}
	profil, err := a.repozytorium.ProfilWydania(ctx, kodProfilu)
	if err != nil {
		return nastawyWyrysuStudia{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "moduł Studio: profil wydania "+kodProfilu+" nie istnieje"))
	}
	if strings.TrimSpace(profil.Ladunek) == "" {
		return nastawy, nil
	}
	var ustawienia ustawieniaProfiluStudia
	if err := json.Unmarshal([]byte(profil.Ladunek), &ustawienia); err != nil {
		// Ładunek nieczytelny nie unieważnia podglądu: profil jest nastawą
		// strony, a nie treścią. Odmowa byłaby tu odmową pokazania dokumentu
		// z powodu ustawienia nagłówka.
		return nastawy, nil
	}
	if ustawienia.Naglowek != "" {
		nastawy.naglowek = ustawienia.Naglowek
	}
	nastawy.stopka = ustawienia.Stopka
	if nastawy.znakWodny == "" {
		nastawy.znakWodny = ustawienia.ZnakWodny
	}
	return nastawy, nil
}

// geometriaDokumentuStudia oddaje geometrię wyrysu wyliczoną z NASTAW STRONY
// dokumentu, a nie ze stałej A4.
//
// Postać dokumentu nieczytelna albo niezapisana daje geometrię domyślną, a nie
// odmowę: podgląd dokumentu, który nastaw strony jeszcze nie ma, jest normalną
// drogą, a nie usterką.
//
// Nastawy SEKCJI ta droga bierze przez nastawy dokumentu, na które sekcja
// pierwsza się nakłada. Dokument o sekcjach różnych NOŚNIKÓW wychodzi w wyrysie
// jednym rozmiarem — wyrys składa jeden ciąg kartek i drugiego rozmiaru w tym
// samym ciągu nie umie. Rachunek stron aparatu (`aparatStronyAkapitow`) liczy
// za to sekcja po sekcji, bo numer strony musi być prawdziwy nawet wtedy, gdy
// obraz kartki jest przybliżeniem.
func (a *adapterStudia) geometriaDokumentuStudia(ctx context.Context,
	kodDokumentu string) geometriaStronyStudia {

	postac, err := a.wejsciePostacDokumentu(ctx, kodDokumentu)
	if err != nil {
		return geometriaDomyslnaStudia()
	}
	return geometriaZNastawStrony(postac.PageSetup)
}

// tytulDokumentuStudia oddaje nazwę, którą dokument nosi w nagłówku strony.
func tytulDokumentuStudia(dokument dane.DokumentStudia) string {
	if dokument.Tytul != nil && strings.TrimSpace(*dokument.Tytul) != "" {
		return *dokument.Tytul
	}
	return dokument.Kod
}

// kartaAlboPustaStudia oddaje stronę o danym numerze albo — gdy wersja jej nie
// ma — pustą kartę tych samych wymiarów.
func kartaAlboPustaStudia(karty []kartaWyrysuStudia, i int,
	kartka geometriaStronyStudia) *image.RGBA {

	if i < len(karty) {
		return karty[i].obraz
	}
	return pustaStronaStudia(kartka)
}

// stronaZeWskaznikaStudia sprowadza pole liczbowe nieobowiązkowe do zera przy braku.
func stronaZeWskaznikaStudia(wskazanie *int) int {
	if wskazanie == nil || *wskazanie < 0 {
		return 0
	}
	return *wskazanie
}
