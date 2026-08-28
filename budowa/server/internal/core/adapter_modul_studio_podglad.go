// Plik obsługuje komendy `studio.preview.render` i `studio.diff.visual`:
// dostarcza treść, dobiera nastawy wyrysu i odkłada wynik jako zasoby stron.
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

// WyrenderujPodglad obsługuje `studio.preview.render`: strony powstają zawsze
// jako obrazy, a przy formacie `pdf` z tych samych stron składa się dokument.
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

	// Dokument w formacie docelowym dokłada się na końcu wykazu, obok stron.
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

// PorownajWizualnie obsługuje `studio.diff.visual`: porównuje wyrys obu
// wersji tym samym silnikiem i tymi samymi nastawami, nie ich tekst.
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
	// Strona, której druga wersja nie ma, wychodzi jako zmiana całej strony.
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
		// Ładunek nieczytelny nie unieważnia podglądu, tylko pomija nastawę.
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

// geometriaDokumentuStudia oddaje geometrię wyrysu wyliczoną z nastaw strony
// dokumentu, a przy postaci nieczytelnej albo niezapisanej — geometrię
// domyślną zamiast odmowy.
func (a *adapterStudia) geometriaDokumentuStudia(ctx context.Context,
	kodDokumentu string) geometriaStronyStudia {

	postac, err := a.wejsciePostacDokumentu(ctx, kodDokumentu)
	if err != nil {
		return geometriaDomyslnaStudia()
	}
	return geometriaZNastawStrony(postac.PageSetup)
}

// tytulDokumentuStudia oddaje nazwę, którą dokument nosi w nagłówku strony,
// z kodem dokumentu jako wartością zastępczą, gdy tytułu nie wskazano.
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

// stronaZeWskaznikaStudia sprowadza pole liczbowe nieobowiązkowe do zera,
// gdy wskaźnik jest pusty albo wartość jest ujemna.
func stronaZeWskaznikaStudia(wskazanie *int) int {
	if wskazanie == nil || *wskazanie < 0 {
		return 0
	}
	return *wskazanie
}
