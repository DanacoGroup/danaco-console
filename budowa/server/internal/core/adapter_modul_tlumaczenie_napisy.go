// Odpowiedzialność pliku: napisy i treść mówiona modułu Translate —
// `translate.subtitle.import`, `.export`, `.timing.check`
// oraz `translate.dubbing.script.build`.
//
// Cztery formaty napisów kontraktu czyta i pisze ten rdzeń sam: SRT i WebVTT są
// tekstowe, TTML jest XML-em, a EBU STL — zapisem dwójkowym o stałej ramce
// (blok nagłówkowy GSI liczący 1024 bajty i bloki tekstowe TTI po 128 bajtów).
// Wszystkie cztery idą bibliotekami wkompilowanymi; żaden nie woła programu.
//
// Kwestie są trwałe (tabela `kwestia_napisow`, migracja 165), bo bez nich
// `subtitle.timing.check` i `dubbing.script.build` nie miałyby czego mierzyć:
// taktowanie jest własnością materiału, a nie tekstu panelu.
package core

import (
	"context"
	"encoding/binary"
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// znakowNaSekundeDomyslnie jest progiem czytelności napisów przyjętym
	// w branży: powyżej siedemnastu znaków na sekundę widz nie nadąża.
	znakowNaSekundeDomyslnie = 17
	// dlugoscLiniiDomyslnie to najdłuższa linia napisu mieszcząca się w kadrze.
	dlugoscLiniiDomyslnie = 42
	// najkrotszeWyswietlenieMs — poniżej sekundy napis miga, zamiast być
	// przeczytany.
	najkrotszeWyswietlenieMs = 1000
	// znakowNaSekundeMowy jest tempem mowy lektorskiej przyjętym do wyliczenia
	// długości kwestii dubbingowej, gdy nie ma taktowania napisów.
	znakowNaSekundeMowy = 14
)

// WczytajNapisy obsługuje `translate.subtitle.import`.
func (a *adapterTlumaczenia) WczytajNapisy(ctx context.Context,
	z shared.TranslateSubtitleImportRequest) (shared.TranslateSubtitleImportResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateSubtitleImportResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	sciezka := strings.TrimSpace(z.Path)
	if sciezka == "" {
		return shared.TranslateSubtitleImportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.subtitle.import bez ścieżki pliku napisów")
	}
	format := shared.SubtitleFormat("")
	if z.Format != nil {
		format = *z.Format
	}
	if strings.TrimSpace(string(format)) == "" {
		rozpoznany, jest := formatNapisowZeSciezki(sciezka)
		if !jest {
			return shared.TranslateSubtitleImportResponse{}, bladWskazaniaTlumaczenia(
				"nie da się rozpoznać formatu napisów po końcówce nazwy — wskaż format wprost")
		}
		format = rozpoznany
	}

	bajty, err := os.ReadFile(filepath.Clean(sciezka))
	if err != nil {
		return shared.TranslateSubtitleImportResponse{}, bladPlikuTlumaczenia(sciezka, err)
	}
	kwestie, err := rozbierzNapisy(bajty, format)
	if err != nil {
		return shared.TranslateSubtitleImportResponse{}, err
	}
	if len(kwestie) == 0 {
		return shared.TranslateSubtitleImportResponse{}, bladWskazaniaTlumaczenia(
			"plik " + filepath.Base(sciezka) + " nie niesie ani jednej kwestii")
	}

	if err := a.repozytorium.ZapiszKwestieNapisow(ctx, okno.ID, 0, kwestie); err != nil {
		return shared.TranslateSubtitleImportResponse{}, bladTlumaczenia(err)
	}

	// Treść kwestii wchodzi do okna jako tekst źródłowy — bez tego napisy byłyby
	// wczytane, a tłumaczyć nie byłoby czego.
	tresci := make([]string, 0, len(kwestie))
	for _, kwestia := range kwestie {
		tresci = append(tresci, kwestia.Tresc)
	}
	tekst := strings.Join(tresci, "\n\n")
	liczba := int64(len(kwestie))
	if _, err := a.repozytorium.ZapiszOkno(ctx, dane.OknoTlumaczenia{
		Kod:             okno.Kod,
		TekstZrodlowy:   &tekst,
		JezykZrodlowy:   okno.JezykZrodlowy,
		LiczbaSegmentow: &liczba,
	}); err != nil {
		return shared.TranslateSubtitleImportResponse{}, bladTlumaczenia(err)
	}
	if err := a.repozytorium.UstawSegmentyOkna(ctx, okno.ID, tresci); err != nil {
		return shared.TranslateSubtitleImportResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslateSubtitleImportResponse{
		Cues:          zlozKwestie(kwestie),
		ImportedCount: len(kwestie),
	}, nil
}

// WydajNapisy obsługuje `translate.subtitle.export`. Bierze taktowanie
// z kwestii materiału źródłowego, a treść — z przekładu panelu, kwestia po
// kwestii. Panel bez kwestii własnych nie ma własnego taktowania, więc idzie
// taktowanie źródła: przekład ma trafiać w te same momenty obrazu.
func (a *adapterTlumaczenia) WydajNapisy(ctx context.Context,
	z shared.TranslateSubtitleExportRequest) (shared.TranslateSubtitleExportResponse, error) {

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateSubtitleExportResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	sciezka := strings.TrimSpace(z.Path)
	if sciezka == "" {
		return shared.TranslateSubtitleExportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.subtitle.export bez ścieżki pliku wyniku")
	}
	kwestie, err := a.kwestiePanelu(ctx, panel)
	if err != nil {
		return shared.TranslateSubtitleExportResponse{}, err
	}
	if err := zapiszNapisy(sciezka, z.Format, kwestie); err != nil {
		return shared.TranslateSubtitleExportResponse{}, err
	}
	// Kwestie przekładu zostają przy panelu: kolejne wywołanie ma mierzyć
	// taktowanie tego, co realnie wydano, a nie składać je od nowa.
	if err := a.repozytorium.ZapiszKwestieNapisow(ctx, panel.OknoID, panel.ID, kwestie); err != nil {
		return shared.TranslateSubtitleExportResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateSubtitleExportResponse{
		Path:          sciezka,
		ExportedCount: len(kwestie),
	}, nil
}

// kwestiePanelu składa kwestie przekładu: taktowanie z kwestii panelu, a przy
// ich braku z kwestii źródłowych okna, treść z akapitów panelu.
func (a *adapterTlumaczenia) kwestiePanelu(ctx context.Context,
	panel dane.PanelTlumaczenia) ([]dane.KwestiaNapisow, error) {

	tresc, err := trescPanelu(panel)
	if err != nil {
		return nil, err
	}
	wlasne, err := a.repozytorium.KwestieNapisow(ctx, panel.OknoID, panel.ID)
	if err != nil {
		return nil, bladTlumaczenia(err)
	}
	taktowanie := wlasne
	if len(taktowanie) == 0 {
		zrodlowe, err := a.repozytorium.KwestieNapisow(ctx, panel.OknoID, 0)
		if err != nil {
			return nil, bladTlumaczenia(err)
		}
		taktowanie = zrodlowe
	}
	if len(taktowanie) == 0 {
		return nil, bladWskazaniaTlumaczenia(
			"okno nie ma wczytanych napisów — bez taktowania nie da się wydać pliku napisów")
	}

	akapity := rozdzielAkapity(tresc)
	kwestie := make([]dane.KwestiaNapisow, 0, len(taktowanie))
	for numer, kwestia := range taktowanie {
		nowa := kwestia
		if numer < len(akapity) && strings.TrimSpace(akapity[numer]) != "" {
			nowa.Tresc = strings.TrimSpace(akapity[numer])
		}
		nowa.Kolejnosc = int64(numer)
		kwestie = append(kwestie, nowa)
	}
	return kwestie, nil
}

// SprawdzTaktowanieNapisow obsługuje `translate.subtitle.timing.check`.
func (a *adapterTlumaczenia) SprawdzTaktowanieNapisow(ctx context.Context,
	z shared.TranslateSubtitleTimingCheckRequest) (shared.TranslateSubtitleTimingCheckResponse, error) {

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateSubtitleTimingCheckResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	kwestie, err := a.kwestiePanelu(ctx, panel)
	if err != nil {
		return shared.TranslateSubtitleTimingCheckResponse{}, err
	}

	naSekunde := float64(liczbaCalkowitaZeWskaznika(z.CharsPerSecond))
	if naSekunde <= 0 {
		naSekunde = znakowNaSekundeDomyslnie
	}
	dlugoscLinii := float64(liczbaCalkowitaZeWskaznika(z.LineLength))
	if dlugoscLinii <= 0 {
		dlugoscLinii = dlugoscLiniiDomyslnie
	}
	najkrotsze := float64(liczbaCalkowitaZeWskaznika(z.MinDurationMs))
	if najkrotsze <= 0 {
		najkrotsze = najkrotszeWyswietlenieMs
	}

	zastrzezenia := []shared.SubtitleTimingIssue{}
	for numer, kwestia := range kwestie {
		trwanie := float64(kwestia.KoniecMs - kwestia.PoczatekMs)
		if trwanie <= 0 {
			zastrzezenia = append(zastrzezenia, shared.SubtitleTimingIssue{
				Index: numer, Kind: shared.SubtitleTimingIssueKindMinDuration,
				Value: trwanie, Limit: najkrotsze,
			})
			continue
		}
		if trwanie < najkrotsze {
			zastrzezenia = append(zastrzezenia, shared.SubtitleTimingIssue{
				Index: numer, Kind: shared.SubtitleTimingIssueKindMinDuration,
				Value: trwanie, Limit: najkrotsze,
			})
		}
		tempo := float64(liczbaZnakow(kwestia.Tresc)) / (trwanie / 1000)
		if tempo > naSekunde {
			zastrzezenia = append(zastrzezenia, shared.SubtitleTimingIssue{
				Index: numer, Kind: shared.SubtitleTimingIssueKindCharsPerSecond,
				Value: tempo, Limit: naSekunde,
			})
		}
		for _, linia := range strings.Split(kwestia.Tresc, "\n") {
			if float64(liczbaZnakow(linia)) > dlugoscLinii {
				zastrzezenia = append(zastrzezenia, shared.SubtitleTimingIssue{
					Index: numer, Kind: shared.SubtitleTimingIssueKindLineLength,
					Value: float64(liczbaZnakow(linia)), Limit: dlugoscLinii,
				})
				break
			}
		}
		if numer > 0 && kwestia.PoczatekMs < kwestie[numer-1].KoniecMs {
			zastrzezenia = append(zastrzezenia, shared.SubtitleTimingIssue{
				Index: numer, Kind: shared.SubtitleTimingIssueKindOverlap,
				Value: float64(kwestie[numer-1].KoniecMs - kwestia.PoczatekMs), Limit: 0,
			})
		}
	}
	return shared.TranslateSubtitleTimingCheckResponse{Issues: zastrzezenia}, nil
}

// ZlozScenariuszDubbingu obsługuje `translate.dubbing.script.build`.
// Kwestia dubbingowa niesie długość wypowiedzenia przekładu i długość, w którą
// ma się zmieścić; przekroczenia są policzone, a nie przemilczane.
func (a *adapterTlumaczenia) ZlozScenariuszDubbingu(ctx context.Context,
	z shared.TranslateDubbingScriptBuildRequest) (shared.TranslateDubbingScriptBuildResponse, error) {

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateDubbingScriptBuildResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	zrodlo := string(shared.DubbingDurationSourceSubtitle)
	if z.TargetDurationSource != nil {
		zrodlo = string(*z.TargetDurationSource)
	}

	var kwestie []dane.KwestiaNapisow
	switch shared.DubbingDurationSource(zrodlo) {
	case shared.DubbingDurationSourceSubtitle, shared.DubbingDurationSourceSourceAudio:
		kwestie, err = a.kwestiePanelu(ctx, panel)
		if err != nil {
			return shared.TranslateDubbingScriptBuildResponse{}, err
		}
	case shared.DubbingDurationSourceSpeechRate:
		// Bez taktowania długość docelowa wynika z tempa mowy: tyle czasu, ile
		// lektor potrzebuje na wypowiedzenie tego zdania.
		tresc, err := trescPanelu(panel)
		if err != nil {
			return shared.TranslateDubbingScriptBuildResponse{}, err
		}
		for numer, zdanie := range podzielNaZdania(tresc) {
			trwanie := int64(liczbaZnakow(zdanie)) * 1000 / znakowNaSekundeMowy
			kwestie = append(kwestie, dane.KwestiaNapisow{
				Kolejnosc: int64(numer), Tresc: zdanie, KoniecMs: trwanie,
			})
		}
	default:
		return shared.TranslateDubbingScriptBuildResponse{}, bladWskazaniaTlumaczenia(
			"nieznane źródło długości kwestii: " + zrodlo)
	}
	if len(kwestie) == 0 {
		return shared.TranslateDubbingScriptBuildResponse{}, bladWskazaniaTlumaczenia(
			"panel nie ma z czego złożyć scenariusza — brak kwestii i brak treści")
	}

	linie := make([]shared.DubbingLine, 0, len(kwestie))
	ponadLimit := 0
	for numer, kwestia := range kwestie {
		wypowiedzenie := int64(liczbaZnakow(kwestia.Tresc)) * 1000 / znakowNaSekundeMowy
		linia := shared.DubbingLine{
			Index:      numer,
			Text:       kwestia.Tresc,
			DurationMs: wypowiedzenie,
		}
		if mowca := mowcaKwestii(kwestia, z.SpeakerHints, numer); mowca != "" {
			linia.Speaker = &mowca
		}
		if trwanie := kwestia.KoniecMs - kwestia.PoczatekMs; trwanie > 0 {
			docelowe := trwanie
			linia.TargetDurationMs = &docelowe
			if wypowiedzenie > docelowe {
				ponadLimit++
			}
		}
		linie = append(linie, linia)
	}
	return shared.TranslateDubbingScriptBuildResponse{
		Lines:          linie,
		OverLimitCount: ponadLimit,
	}, nil
}

// mowcaKwestii bierze mówcę z kwestii, a przy jego braku z podpowiedzi żądania
// przypisywanych kolejno po numerze kwestii.
func mowcaKwestii(kwestia dane.KwestiaNapisow, podpowiedzi []string, numer int) string {
	if kwestia.Mowca != nil && strings.TrimSpace(*kwestia.Mowca) != "" {
		return *kwestia.Mowca
	}
	if len(podpowiedzi) == 0 {
		return ""
	}
	return podpowiedzi[numer%len(podpowiedzi)]
}

// zlozKwestie przekłada wiersze kwestii na byty kontraktu.
func zlozKwestie(kwestie []dane.KwestiaNapisow) []shared.SubtitleCue {
	wykaz := make([]shared.SubtitleCue, 0, len(kwestie))
	for _, kwestia := range kwestie {
		wykaz = append(wykaz, shared.SubtitleCue{
			Index:   int(kwestia.Kolejnosc),
			StartMs: kwestia.PoczatekMs,
			EndMs:   kwestia.KoniecMs,
			Text:    kwestia.Tresc,
			Speaker: kwestia.Mowca,
		})
	}
	return wykaz
}

// formatNapisowZeSciezki rozpoznaje format po końcówce nazwy.
func formatNapisowZeSciezki(sciezka string) (shared.SubtitleFormat, bool) {
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".srt":
		return shared.SubtitleFormatSrt, true
	case ".vtt", ".webvtt":
		return shared.SubtitleFormatWebvtt, true
	case ".ttml", ".dfxp", ".xml":
		return shared.SubtitleFormatTtml, true
	case ".stl":
		return shared.SubtitleFormatStl, true
	}
	return "", false
}

// rozbierzNapisy czyta kwestie z pliku wskazanego formatu.
func rozbierzNapisy(bajty []byte, format shared.SubtitleFormat) ([]dane.KwestiaNapisow, error) {
	switch format {
	case shared.SubtitleFormatSrt, shared.SubtitleFormatWebvtt:
		return kwestieZTekstu(string(bajty)), nil
	case shared.SubtitleFormatTtml:
		return kwestieZTtml(bajty)
	case shared.SubtitleFormatStl:
		return kwestieZStl(bajty)
	}
	return nil, bladWskazaniaTlumaczenia("nieznany format napisów: " + string(format))
}

// kwestieZTekstu czyta SRT i WebVTT jednym rozbiorem: oba dzielą kwestie pustą
// linią, oba niosą przedział czasu w linii ze strzałką, oba mają resztę bloku
// jako treść. Różnią się separatorem części ułamkowej (przecinek w SRT, kropka
// w WebVTT) i nagłówkiem `WEBVTT`, który wypada wraz z blokiem bez czasu.
func kwestieZTekstu(tresc string) []dane.KwestiaNapisow {
	bloki := strings.Split(strings.ReplaceAll(tresc, "\r\n", "\n"), "\n\n")
	kwestie := []dane.KwestiaNapisow{}
	for _, blok := range bloki {
		linie := strings.Split(strings.TrimSpace(blok), "\n")
		poczatek, koniec := int64(-1), int64(-1)
		tresci := []string{}
		for _, linia := range linie {
			if strings.Contains(linia, "-->") {
				od, do_, jest := strings.Cut(linia, "-->")
				if !jest {
					continue
				}
				poczatek = czasNapisowMs(od)
				koniec = czasNapisowMs(do_)
				continue
			}
			if poczatek >= 0 {
				tresci = append(tresci, linia)
			}
		}
		if poczatek < 0 || koniec < 0 || len(tresci) == 0 {
			continue
		}
		kwestie = append(kwestie, dane.KwestiaNapisow{
			Kolejnosc:  int64(len(kwestie)),
			PoczatekMs: poczatek,
			KoniecMs:   koniec,
			Tresc:      strings.Join(tresci, "\n"),
		})
	}
	return kwestie
}

// czasNapisowMs czyta znacznik czasu w postaci `hh:mm:ss,mmm` albo
// `hh:mm:ss.mmm`, a także skróconej `mm:ss.mmm`.
func czasNapisowMs(zapis string) int64 {
	oczyszczony := strings.TrimSpace(strings.ReplaceAll(zapis, ",", "."))
	oczyszczony = strings.Fields(oczyszczony + " ")[0]
	czesci := strings.Split(oczyszczony, ":")
	if len(czesci) < 2 {
		return -1
	}
	var razem float64
	for _, czesc := range czesci {
		wartosc, err := strconv.ParseFloat(czesc, 64)
		if err != nil {
			return -1
		}
		razem = razem*60 + wartosc
	}
	return int64(razem * 1000)
}

// zapisCzasuNapisow składa znacznik czasu w postaci wymaganej przez format.
func zapisCzasuNapisow(ms int64, kropka bool) string {
	godziny := ms / 3600000
	minuty := (ms % 3600000) / 60000
	sekundy := (ms % 60000) / 1000
	tysieczne := ms % 1000
	separator := ","
	if kropka {
		separator = "."
	}
	return dwaZnaki(godziny) + ":" + dwaZnaki(minuty) + ":" + dwaZnaki(sekundy) +
		separator + trzyZnaki(tysieczne)
}

// dwaZnaki i trzyZnaki dopełniają liczbę zerami — znacznik czasu ma stałą
// szerokość pól, inaczej odtwarzacz go nie przyjmie.
func dwaZnaki(wartosc int64) string {
	if wartosc < 10 {
		return "0" + strconv.FormatInt(wartosc, 10)
	}
	return strconv.FormatInt(wartosc, 10)
}

func trzyZnaki(wartosc int64) string {
	zapis := strconv.FormatInt(wartosc, 10)
	for len(zapis) < 3 {
		zapis = "0" + zapis
	}
	return zapis
}

// kwestieZTtml czyta kwestie z dokumentu TTML.
func kwestieZTtml(bajty []byte) ([]dane.KwestiaNapisow, error) {
	czytnik := xml.NewDecoder(strings.NewReader(string(bajty)))
	kwestie := []dane.KwestiaNapisow{}
	var biezaca *dane.KwestiaNapisow
	var tresc strings.Builder
	for {
		znacznik, err := czytnik.Token()
		if err != nil {
			break
		}
		switch element := znacznik.(type) {
		case xml.StartElement:
			if element.Name.Local != "p" {
				continue
			}
			kwestia := dane.KwestiaNapisow{Kolejnosc: int64(len(kwestie))}
			for _, cecha := range element.Attr {
				switch cecha.Name.Local {
				case "begin":
					kwestia.PoczatekMs = czasNapisowMs(cecha.Value)
				case "end":
					kwestia.KoniecMs = czasNapisowMs(cecha.Value)
				case "agent", "speaker":
					mowca := cecha.Value
					kwestia.Mowca = &mowca
				}
			}
			biezaca = &kwestia
			tresc.Reset()
		case xml.CharData:
			if biezaca != nil {
				tresc.Write(element)
			}
		case xml.EndElement:
			if element.Name.Local != "p" || biezaca == nil {
				continue
			}
			biezaca.Tresc = strings.TrimSpace(tresc.String())
			if biezaca.Tresc != "" {
				kwestie = append(kwestie, *biezaca)
			}
			biezaca = nil
		}
	}
	return kwestie, nil
}

// rozmiarGsi i rozmiarTti to stałe ramki standardu EBU 3264: blok nagłówkowy
// i blok tekstowy.
const (
	rozmiarGsi = 1024
	rozmiarTti = 128
)

// kwestieZStl czyta plik EBU STL. Czas kodowany jest czterema bajtami
// (godzina, minuta, sekunda, klatka); liczba klatek na sekundę bierze się
// z pola `DFC` nagłówka, a jego brak oznacza dwadzieścia pięć klatek — tak
// stanowi standard dla materiału europejskiego.
func kwestieZStl(bajty []byte) ([]dane.KwestiaNapisow, error) {
	if len(bajty) < rozmiarGsi+rozmiarTti {
		return nil, bladWskazaniaTlumaczenia("plik jest krótszy niż nagłówek EBU STL — to nie są napisy STL")
	}
	klatek := 25.0
	if string(bajty[6:9]) == "STL" {
		if zapis := strings.TrimSpace(string(bajty[11:14])); zapis == "30" {
			klatek = 30
		}
	}
	kwestie := []dane.KwestiaNapisow{}
	for przesuniecie := rozmiarGsi; przesuniecie+rozmiarTti <= len(bajty); przesuniecie += rozmiarTti {
		blok := bajty[przesuniecie : przesuniecie+rozmiarTti]
		poczatek := czasStlMs(blok[5:9], klatek)
		koniec := czasStlMs(blok[9:13], klatek)
		tresc := trescStl(blok[16:])
		if strings.TrimSpace(tresc) == "" {
			continue
		}
		kwestie = append(kwestie, dane.KwestiaNapisow{
			Kolejnosc:  int64(len(kwestie)),
			PoczatekMs: poczatek,
			KoniecMs:   koniec,
			Tresc:      tresc,
		})
	}
	return kwestie, nil
}

// czasStlMs przekłada czwórkę bajtów standardu na milisekundy.
func czasStlMs(bajty []byte, klatek float64) int64 {
	if len(bajty) < 4 {
		return 0
	}
	sekundy := float64(bajty[0])*3600 + float64(bajty[1])*60 + float64(bajty[2])
	return int64((sekundy + float64(bajty[3])/klatek) * 1000)
}

// trescStl czyta pole tekstowe bloku: bajt 0x8A jest końcem linii, bajty od
// 0x80 wzwyż są sterujące barwą i tłem, a 0x8F wypełnia resztę pola.
func trescStl(pole []byte) string {
	var b strings.Builder
	for _, bajt := range pole {
		switch {
		case bajt == 0x8A:
			b.WriteByte('\n')
		case bajt == 0x8F:
			// Wypełniacz do końca pola — dalej nie ma już treści.
			return strings.TrimSpace(b.String())
		case bajt >= 0x20 && bajt < 0x80:
			b.WriteByte(bajt)
		}
	}
	return strings.TrimSpace(b.String())
}

// zapiszNapisy wypisuje kwestie do pliku wskazanego formatu.
func zapiszNapisy(sciezka string, format shared.SubtitleFormat,
	kwestie []dane.KwestiaNapisow) error {

	if err := os.MkdirAll(filepath.Dir(sciezka), 0o700); err != nil {
		return bladPlikuTlumaczenia(sciezka, err)
	}
	sort.SliceStable(kwestie, func(i, j int) bool {
		return kwestie[i].Kolejnosc < kwestie[j].Kolejnosc
	})

	switch format {
	case shared.SubtitleFormatSrt:
		var b strings.Builder
		for numer, kwestia := range kwestie {
			b.WriteString(strconv.Itoa(numer+1) + "\n")
			b.WriteString(zapisCzasuNapisow(kwestia.PoczatekMs, false) + " --> " +
				zapisCzasuNapisow(kwestia.KoniecMs, false) + "\n")
			b.WriteString(kwestia.Tresc + "\n\n")
		}
		return zapiszPlikWyniku(sciezka, []byte(b.String()))

	case shared.SubtitleFormatWebvtt:
		var b strings.Builder
		b.WriteString("WEBVTT\n\n")
		for _, kwestia := range kwestie {
			b.WriteString(zapisCzasuNapisow(kwestia.PoczatekMs, true) + " --> " +
				zapisCzasuNapisow(kwestia.KoniecMs, true) + "\n")
			b.WriteString(kwestia.Tresc + "\n\n")
		}
		return zapiszPlikWyniku(sciezka, []byte(b.String()))

	case shared.SubtitleFormatTtml:
		var b strings.Builder
		b.WriteString(xml.Header)
		b.WriteString(`<tt xmlns="http://www.w3.org/ns/ttml"><body><div>` + "\n")
		for _, kwestia := range kwestie {
			b.WriteString(`  <p begin="` + zapisCzasuNapisow(kwestia.PoczatekMs, true) +
				`" end="` + zapisCzasuNapisow(kwestia.KoniecMs, true) + `">` +
				zabezpieczHtml(strings.ReplaceAll(kwestia.Tresc, "\n", " ")) + "</p>\n")
		}
		b.WriteString("</div></body></tt>\n")
		return zapiszPlikWyniku(sciezka, []byte(b.String()))

	case shared.SubtitleFormatStl:
		return zapiszPlikWyniku(sciezka, zlozStl(kwestie))
	}
	return bladWskazaniaTlumaczenia("nieznany format napisów: " + string(format))
}

// zlozStl składa plik EBU STL: nagłówek GSI i po jednym bloku TTI na kwestię.
// Pola nagłówka, których treść nie wynika z materiału, wypełnia znak odstępu —
// tak każe standard, i tak plik otwierają czytniki napisów nadawczych.
func zlozStl(kwestie []dane.KwestiaNapisow) []byte {
	nagłowek := make([]byte, rozmiarGsi)
	for i := range nagłowek {
		nagłowek[i] = ' '
	}
	copy(nagłowek[0:], "850STL25.01")
	copy(nagłowek[16:], "Danaco Console")
	// Liczba bloków tekstowych i liczba napisów — pola pięcioznakowe.
	liczba := strconv.Itoa(len(kwestie))
	for len(liczba) < 5 {
		liczba = "0" + liczba
	}
	copy(nagłowek[238:], liczba)
	copy(nagłowek[243:], liczba)

	wynik := nagłowek
	for numer, kwestia := range kwestie {
		blok := make([]byte, rozmiarTti)
		for i := range blok {
			blok[i] = 0x8F
		}
		blok[0] = 0 // grupa napisów
		binary.LittleEndian.PutUint16(blok[1:3], uint16(numer+1))
		blok[3] = 0    // numer bloku rozszerzenia
		blok[4] = 0x00 // barwa kanału
		copy(blok[5:9], czasStlZMs(kwestia.PoczatekMs))
		copy(blok[9:13], czasStlZMs(kwestia.KoniecMs))
		blok[13] = 22   // wiersz wyświetlenia
		blok[14] = 0x02 // wyrównanie do środka
		blok[15] = 0x00 // rodzaj bloku: napis zwykły
		tresc := []byte(strings.ReplaceAll(kwestia.Tresc, "\n", "\x8a"))
		if len(tresc) > rozmiarTti-16 {
			tresc = tresc[:rozmiarTti-16]
		}
		copy(blok[16:], tresc)
		wynik = append(wynik, blok...)
	}
	return wynik
}

// czasStlZMs przekłada milisekundy na czwórkę bajtów standardu przy dwudziestu
// pięciu klatkach na sekundę — tej samej liczbie, którą deklaruje nagłówek.
func czasStlZMs(ms int64) []byte {
	sekundy := ms / 1000
	klatki := (ms % 1000) * 25 / 1000
	return []byte{
		byte(sekundy / 3600), byte((sekundy % 3600) / 60), byte(sekundy % 60), byte(klatki),
	}
}
