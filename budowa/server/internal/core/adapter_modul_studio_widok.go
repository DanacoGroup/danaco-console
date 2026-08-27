// Odpowiedzialność pliku: NASTAWY WIDOKU okna pracy z dokumentem —
// studio.view.get i studio.view.set, jako jawne, odwracalne ustawienie
// Operatora.
package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// NastawyWidoku obsługuje studio.view.get, oddając stan zapisany w bazie
// dla pary okno-dokument, bez odczytu z klienta.
func (a *adapterStudia) NastawyWidoku(ctx context.Context,
	z shared.StudioViewGetRequest) (shared.StudioViewGetResponse, error) {

	nastawa, dokument, err := a.widokNastawa(ctx, z.WindowId, z.DocumentId)
	if err != nil {
		return shared.StudioViewGetResponse{}, err
	}
	return shared.StudioViewGetResponse{
		Settings: widokZlozNastawy(nastawa, dokument),
	}, nil
}

// UstawWidok obsługuje studio.view.set; pola pominięte zostają w brzmieniu
// zastanym, żądanie nie ma zwijać linijek niepytanych.
func (a *adapterStudia) UstawWidok(ctx context.Context,
	z shared.StudioViewSetRequest) (shared.StudioViewSetResponse, error) {

	nastawa, dokument, err := a.widokNastawa(ctx, z.WindowId, z.DocumentId)
	if err != nil {
		return shared.StudioViewSetResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioViewSetResponse{}, err
	}

	zmienione := 0
	if z.SurfaceMode != nil && *z.SurfaceMode != "" {
		if err := widokSprawdzWartosc("tryb powierzchni", string(*z.SurfaceMode),
			widokNazwy(shared.WartosciStudioSurfaceMode())); err != nil {

			return shared.StudioViewSetResponse{}, err
		}
		nastawa.TrybPowierzchni = string(*z.SurfaceMode)
		zmienione++
	}
	if z.SplitOrientation != nil && *z.SplitOrientation != "" {
		if err := widokSprawdzWartosc("kierunek podziału", string(*z.SplitOrientation),
			widokNazwy(shared.WartosciStudioSplitOrientation())); err != nil {

			return shared.StudioViewSetResponse{}, err
		}
		nastawa.KierunekPodzialu = string(*z.SplitOrientation)
		zmienione++
	}
	if z.SplitRatio != nil {
		if *z.SplitRatio <= 0 || *z.SplitRatio >= 1 {
			return shared.StudioViewSetResponse{}, bladWskazaniaStudio(
				"położenie granicy podziału musi leżeć między zerem a jednością — " +
					"granica na krawędzi schowałaby jeden z dokumentów, a to nie jest podział")
		}
		nastawa.GranicaPodzialu = *z.SplitRatio
		zmienione++
	}
	if z.ViewMode != nil && *z.ViewMode != "" {
		if err := widokSprawdzWartosc("tryb widoku", string(*z.ViewMode),
			widokNazwy(shared.WartosciStudioViewMode())); err != nil {

			return shared.StudioViewSetResponse{}, err
		}
		nastawa.TrybWidoku = string(*z.ViewMode)
		zmienione++
	}
	if z.ZoomPercent != nil {
		if *z.ZoomPercent < 10 || *z.ZoomPercent > 1000 {
			return shared.StudioViewSetResponse{}, bladWskazaniaStudio(
				"skala widoku " + strconv.Itoa(*z.ZoomPercent) + "% wypada poza zakresem " +
					"od 10% do 1000% — poza nim dokument przestaje być czytelny")
		}
		nastawa.SkalaProcent = int64(*z.ZoomPercent)
		zmienione++
	}
	if z.ZoomPreset != nil && *z.ZoomPreset != "" {
		if err := widokSprawdzWartosc("nastawa skali", string(*z.ZoomPreset),
			widokNazwy(shared.WartosciStudioZoomPreset())); err != nil {

			return shared.StudioViewSetResponse{}, err
		}
		nastawa.SkalaNastawa = string(*z.ZoomPreset)
		zmienione++
	}
	if z.RulersVisible != nil {
		nastawa.LinijkiWidoczne = *z.RulersVisible
		zmienione++
	}
	if z.RulerUnit != nil && *z.RulerUnit != "" {
		if err := widokSprawdzWartosc("jednostka linijki", string(*z.RulerUnit),
			widokNazwy(shared.WartosciStudioRulerUnit())); err != nil {

			return shared.StudioViewSetResponse{}, err
		}
		nastawa.LinijkaJednostka = string(*z.RulerUnit)
		zmienione++
	}
	if z.MarginGuides != nil {
		nastawa.GranicaMarginesu = *z.MarginGuides
		zmienione++
	}
	if z.FormattingMarks != nil {
		nastawa.ZnakiFormatowania = *z.FormattingMarks
		zmienione++
	}
	if z.PagesPerRow != nil {
		if *z.PagesPerRow < 1 || *z.PagesPerRow > 8 {
			return shared.StudioViewSetResponse{}, bladWskazaniaStudio(
				"liczba stron w rzędzie musi leżeć między jedną a ośmioma — zero stron " +
					"nie jest widokiem, a więcej niż osiem przestaje być czytelne")
		}
		nastawa.StronWRzedzie = int64(*z.PagesPerRow)
		zmienione++
	}
	if z.SpreadView != nil {
		nastawa.WidokRozkladowki = *z.SpreadView
		zmienione++
	}
	if z.ScrollMode != nil && *z.ScrollMode != "" {
		if err := widokSprawdzWartosc("sposób przewijania", string(*z.ScrollMode),
			widokNazwy(shared.WartosciStudioScrollMode())); err != nil {

			return shared.StudioViewSetResponse{}, err
		}
		nastawa.Przewijanie = string(*z.ScrollMode)
		zmienione++
	}
	if z.ModelChangesHighlighted != nil {
		nastawa.PodswietlenieZmianModelu = *z.ModelChangesHighlighted
		zmienione++
	}
	if z.ToolboxVisible != nil {
		nastawa.PrzybornikWidoczny = *z.ToolboxVisible
		zmienione++
	}
	if zmienione == 0 {
		return shared.StudioViewSetResponse{}, bladWskazaniaStudio(
			"żądanie nie niesie ani jednej nastawy widoku do przestawienia — " +
				"odpowiedź pomyślna kazałaby czytać to jako zmianę, której nie było")
	}

	if err := skladnica.ZapiszNastaweWidoku(ctx, nastawa); err != nil {
		return shared.StudioViewSetResponse{}, bladStudio(err)
	}
	// Odpowiedź czyta stan z bazy: ma mówić, jak jest, a nie powtarzać żądanie.
	zapisana, dokumentPo, err := a.widokNastawa(ctx, z.WindowId, z.DocumentId)
	if err != nil {
		return shared.StudioViewSetResponse{}, err
	}
	_ = dokument
	return shared.StudioViewSetResponse{
		Settings: widokZlozNastawy(zapisana, dokumentPo),
	}, nil
}

// ── Wspólne ─────────────────────────────────────────────────────────────────

// widokNastawa odczytuje kolumny widoku wiersza nastaw dla pary
// okno-dokument albo dla samego okna, zakładając wiersz, gdy trzeba.
func (a *adapterStudia) widokNastawa(ctx context.Context, oknoZadania,
	dokumentZadania *string) (dane.NastawaWidokuStudia, dane.DokumentStudia, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return dane.NastawaWidokuStudia{}, dane.DokumentStudia{}, err
	}
	var dokument dane.DokumentStudia
	var wskazanie *int64
	okno := strings.TrimSpace(wartoscTekstu(oknoZadania))
	if kontrolaTekstNiepusty(dokumentZadania) {
		if dokument, err = a.dokumentDoCzynnosci(ctx, *dokumentZadania); err != nil {
			return dane.NastawaWidokuStudia{}, dane.DokumentStudia{}, err
		}
		wskazanie = &dokument.ID
		if okno == "" {
			okno = dokument.Okno
		}
	}
	if okno == "" {
		return dane.NastawaWidokuStudia{}, dane.DokumentStudia{}, bladWskazaniaStudio(
			"nastawy widoku bez wskazania okna ani dokumentu — nastawa musi wiedzieć, " +
				"czyja jest: wiersz bez dokumentu jest nastawą okna, wiersz z dokumentem " +
				"nastawą tego dokumentu")
	}
	// Zakłada wiersz, jeśli go nie ma — i dopiero potem czyta kolumny widoku.
	if _, err := skladnica.NastawaPracy(ctx, okno, wskazanie); err != nil {
		return dane.NastawaWidokuStudia{}, dane.DokumentStudia{}, bladStudio(err)
	}
	nastawa, err := skladnica.NastawaWidoku(ctx, okno, wskazanie)
	if err != nil {
		return dane.NastawaWidokuStudia{}, dane.DokumentStudia{}, bladStudio(err)
	}
	return nastawa, dokument, nil
}

// widokNazwy przekłada wykaz wartości wyliczenia na napisy do treści
// odmowy, czytelne dla Operatora czytającego komunikat.
func widokNazwy[T ~string](wartosci []T) []string {
	nazwy := make([]string, 0, len(wartosci))
	for _, wartosc := range wartosci {
		nazwy = append(nazwy, string(wartosc))
	}
	return nazwy
}

// widokSprawdzWartosc odbija wartość spoza wyliczenia kontraktu, zanim
// trafi do zapisu w tabeli nastaw widoku.
func widokSprawdzWartosc(nazwaPola, wartosc string, dozwolone []string) error {
	for _, dozwolona := range dozwolone {
		if wartosc == dozwolona {
			return nil
		}
	}
	return bladWskazaniaStudio(nazwaPola + " „" + wartosc + "” nie jest znany — wolno: " +
		strings.Join(dozwolone, ", "))
}

// widokZlozNastawy składa nastawy widoku kontraktu z wiersza warstwy
// danych, w kształcie oczekiwanym przez odpowiedź.
func widokZlozNastawy(wiersz dane.NastawaWidokuStudia,
	dokument dane.DokumentStudia) shared.StudioViewSettings {

	trybPowierzchni := shared.StudioSurfaceMode(wiersz.TrybPowierzchni)
	kierunek := shared.StudioSplitOrientation(wiersz.KierunekPodzialu)
	granica := wiersz.GranicaPodzialu
	trybWidoku := shared.StudioViewMode(wiersz.TrybWidoku)
	skala := int(wiersz.SkalaProcent)
	nastawaSkali := shared.StudioZoomPreset(wiersz.SkalaNastawa)
	jednostka := shared.StudioRulerUnit(wiersz.LinijkaJednostka)
	stron := int(wiersz.StronWRzedzie)
	przewijanie := shared.StudioScrollMode(wiersz.Przewijanie)

	nastawy := shared.StudioViewSettings{
		SurfaceMode:             &trybPowierzchni,
		SplitOrientation:        &kierunek,
		SplitRatio:              &granica,
		ViewMode:                &trybWidoku,
		ZoomPercent:             &skala,
		ZoomPreset:              &nastawaSkali,
		RulersVisible:           wskaznikLogiczny(wiersz.LinijkiWidoczne),
		RulerUnit:               &jednostka,
		MarginGuides:            wskaznikLogiczny(wiersz.GranicaMarginesu),
		FormattingMarks:         wskaznikLogiczny(wiersz.ZnakiFormatowania),
		PagesPerRow:             &stron,
		SpreadView:              wskaznikLogiczny(wiersz.WidokRozkladowki),
		ScrollMode:              &przewijanie,
		ModelChangesHighlighted: wskaznikLogiczny(wiersz.PodswietlenieZmianModelu),
		ToolboxVisible:          wskaznikLogiczny(wiersz.PrzybornikWidoczny),
	}
	if wiersz.Okno != "" {
		okno := wiersz.Okno
		nastawy.WindowId = &okno
	}
	if dokument.Kod != "" {
		kod := dokument.Kod
		nastawy.DocumentId = &kod
	}
	return nastawy
}
