// Odpowiedzialność pliku: rachunki wspólne obszaru strony i sekcji — wykaz
// nośnika, wymiary użytkowe po marginesach, scalanie nastaw strony, nagłówków
// i znaku wodnego, przekład sekcji na wiersz danych oraz bilans przeliczenia
// układu po zmianie nośnika.
package core

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// stronaNazwaSekcjiPierwszej to nazwa sekcji zakładanej dokumentowi, który
// swojej nie ma. Nazwa jest pełna i widoczna dla Operatora — żadnego kodu.
const stronaNazwaSekcjiPierwszej = "Sekcja pierwsza"

// stronaNosniki oddaje wykaz formatów nośnika Studia, wzięty z jedynego
// wykazu nośników druku wspólnego dla całego rdzenia, łącznie z kopertami
// C4, C5 i C6.
func stronaNosniki() []shared.StudioPaperFormat {
	return nosnikiDrukuJakoFormatyStudia()
}

// stronaNosnik znajduje nośnik po nazwie, nie bacząc na wielkość liter —
// Operator wpisuje „a4" tak samo często jak „A4".
func stronaNosnik(nazwa string) (shared.StudioPaperFormat, bool) {
	szukana := strings.ToUpper(strings.TrimSpace(nazwa))
	if szukana == "" {
		return shared.StudioPaperFormat{}, false
	}
	for _, nosnik := range stronaNosniki() {
		if strings.ToUpper(nosnik.Name) == szukana {
			return nosnik, true
		}
	}
	return shared.StudioPaperFormat{}, false
}

// stronaNazwyNosnikow oddaje nazwy nośników do treści odmowy. Odmowa bez wykazu
// zostawiałaby Operatora bez drogi wyjścia.
func stronaNazwyNosnikow() []string {
	wykaz := stronaNosniki()
	nazwy := make([]string, 0, len(wykaz))
	for _, nosnik := range wykaz {
		nazwy = append(nazwy, nosnik.Name)
	}
	return nazwy
}

// stronaNastawaMarginesow oddaje marginesy nastawy gotowej w milimetrach
// w kolejności górny, dolny, lewy, prawy. Wartości są tymi, które Operator zna
// z pakietu biurowego.
func stronaNastawaMarginesow(nazwa string) ([4]float64, bool) {
	switch strings.ToLower(strings.TrimSpace(nazwa)) {
	case "wąskie", "waskie":
		return [4]float64{12.7, 12.7, 12.7, 12.7}, true
	case "normalne":
		return [4]float64{25, 25, 25, 25}, true
	case "szerokie":
		return [4]float64{25.4, 25.4, 50.8, 50.8}, true
	default:
		return [4]float64{}, false
	}
}

// stronaWymiary oddaje szerokość i wysokość nośnika nastaw, już po obrocie
// wynikającym z orientacji. Rachunek jest jeden i stoi tutaj — podgląd wydruku,
// bilans przeliczenia i nadruk koperty muszą liczyć te same milimetry.
func stronaWymiary(nastawy *shared.StudioPageSetup) (float64, float64) {
	szerokosc, wysokosc := 210.0, 297.0
	if nastawy == nil {
		return szerokosc, wysokosc
	}
	if nastawy.PageSize != nil {
		if nosnik, jest := stronaNosnik(*nastawy.PageSize); jest {
			szerokosc, wysokosc = nosnik.WidthMm, nosnik.HeightMm
		}
	}
	// Wymiar własny przebija nazwę nośnika.
	if nastawy.WidthMm != nil && *nastawy.WidthMm > 0 {
		szerokosc = *nastawy.WidthMm
	}
	if nastawy.HeightMm != nil && *nastawy.HeightMm > 0 {
		wysokosc = *nastawy.HeightMm
	}
	if nastawy.Orientation != nil && *nastawy.Orientation == shared.StudioPageOrientationPozioma {
		szerokosc, wysokosc = wysokosc, szerokosc
	}
	return szerokosc, wysokosc
}

// stronaSzerokoscUzytkowa liczy szerokość kolumny tekstu: nośnik bez marginesów,
// bez marginesu na oprawę i podzielony na kolumny wraz z odstępami.
func stronaSzerokoscUzytkowa(nastawy *shared.StudioPageSetup) float64 {
	szerokosc, _ := stronaWymiary(nastawy)
	if nastawy == nil {
		return szerokosc
	}
	if nastawy.MarginLeft != nil {
		szerokosc -= float64(*nastawy.MarginLeft)
	}
	if nastawy.MarginRight != nil {
		szerokosc -= float64(*nastawy.MarginRight)
	}
	if nastawy.GutterMm != nil {
		szerokosc -= *nastawy.GutterMm
	}
	kolumny := 1
	if nastawy.Columns != nil && *nastawy.Columns > 1 {
		kolumny = *nastawy.Columns
	}
	if kolumny > 1 {
		odstep := 0.0
		if nastawy.ColumnGapMm != nil {
			odstep = *nastawy.ColumnGapMm
		}
		szerokosc = (szerokosc - odstep*float64(kolumny-1)) / float64(kolumny)
	}
	if szerokosc < 0 {
		szerokosc = 0
	}
	return szerokosc
}

// stronaScalNastawy wnosi do nastaw strony wyłącznie pola podane w żądaniu —
// tą samą zasadą, którą trzyma cały obszar postaci. Zmiana orientacji nie ma
// prawa zdjąć marginesu na oprawę.
func stronaScalNastawy(zastane *shared.StudioPageSetup,
	z shared.StudioPageSetupSetRequest) (shared.StudioPageSetup, int, error) {

	wynik := shared.StudioPageSetup{}
	if zastane != nil {
		wynik = *zastane
		if zastane.Envelope != nil {
			koperta := *zastane.Envelope
			wynik.Envelope = &koperta
		}
	}
	zmian := 0

	if z.PaperName != nil && strings.TrimSpace(*z.PaperName) != "" {
		nosnik, jest := stronaNosnik(*z.PaperName)
		if !jest {
			return shared.StudioPageSetup{}, 0, bladWskazaniaStudio(
				"nośnika „" + strings.TrimSpace(*z.PaperName) + "” rdzeń nie zna. " +
					"Nośniki znane: " + strings.Join(stronaNazwyNosnikow(), ", ") +
					". Format własny podaje się polami widthMm i heightMm")
		}
		wynik.PageSize = postacWskaznikTekstu(nosnik.Name)
		rodzaj := nosnik.Kind
		wynik.PaperKind = &rodzaj
		// Nazwa z wykazu zdejmuje wymiar własny nastawy zastanej.
		wynik.WidthMm, wynik.HeightMm = nil, nil
		zmian++
	}
	if z.WidthMm != nil || z.HeightMm != nil {
		if z.WidthMm != nil {
			if *z.WidthMm <= 0 {
				return shared.StudioPageSetup{}, 0, bladWskazaniaStudio(
					"szerokość nośnika własnego musi być większa od zera")
			}
			wynik.WidthMm = postacWskaznikMiary(*z.WidthMm)
		}
		if z.HeightMm != nil {
			if *z.HeightMm <= 0 {
				return shared.StudioPageSetup{}, 0, bladWskazaniaStudio(
					"wysokość nośnika własnego musi być większa od zera")
			}
			wynik.HeightMm = postacWskaznikMiary(*z.HeightMm)
		}
		if z.PaperName == nil {
			rodzaj := shared.StudioPaperKind(shared.StudioPaperKindCustom)
			wynik.PaperKind = &rodzaj
			wynik.PageSize = postacWskaznikTekstu("własny")
		}
		zmian++
	}
	if z.Orientation != nil && strings.TrimSpace(string(*z.Orientation)) != "" {
		kierunek := *z.Orientation
		wynik.Orientation = &kierunek
		zmian++
	}
	if z.MarginPreset != nil && strings.TrimSpace(*z.MarginPreset) != "" {
		marginesy, jest := stronaNastawaMarginesow(*z.MarginPreset)
		if !jest {
			return shared.StudioPageSetup{}, 0, bladWskazaniaStudio(
				"nastawy marginesów „" + strings.TrimSpace(*z.MarginPreset) +
					"” rdzeń nie zna. Nastawy gotowe: wąskie, normalne, szerokie; " +
					"marginesy własne podaje się polami marginTopMm, marginBottomMm, " +
					"marginLeftMm i marginRightMm")
		}
		wynik.MarginTop = postacWskaznikLiczby(int(marginesy[0] + 0.5))
		wynik.MarginBottom = postacWskaznikLiczby(int(marginesy[1] + 0.5))
		wynik.MarginLeft = postacWskaznikLiczby(int(marginesy[2] + 0.5))
		wynik.MarginRight = postacWskaznikLiczby(int(marginesy[3] + 0.5))
		wynik.MarginPreset = postacWskaznikTekstu(strings.ToLower(strings.TrimSpace(*z.MarginPreset)))
		zmian++
	}
	// Margines podany wprost przebija nastawę gotową i odznacza się jako własny.
	wlasnyMargines := false
	for _, para := range []struct {
		wartosc *float64
		zapis   **int
	}{
		{z.MarginTopMm, &wynik.MarginTop},
		{z.MarginBottomMm, &wynik.MarginBottom},
		{z.MarginLeftMm, &wynik.MarginLeft},
		{z.MarginRightMm, &wynik.MarginRight},
	} {
		if para.wartosc == nil {
			continue
		}
		if *para.wartosc < 0 {
			return shared.StudioPageSetup{}, 0, bladWskazaniaStudio(
				"margines ujemny nie jest marginesem — podaj wartość nie mniejszą od zera")
		}
		*para.zapis = postacWskaznikLiczby(int(*para.wartosc + 0.5))
		wlasnyMargines = true
		zmian++
	}
	if wlasnyMargines && z.MarginPreset == nil {
		wynik.MarginPreset = postacWskaznikTekstu("własne")
	}
	if z.GutterMm != nil {
		if *z.GutterMm < 0 {
			return shared.StudioPageSetup{}, 0, bladWskazaniaStudio(
				"margines na oprawę nie może być ujemny")
		}
		wynik.GutterMm = postacWskaznikMiary(*z.GutterMm)
		zmian++
	}
	if z.MirrorMargins != nil {
		wynik.MirrorMargins = postacWskaznikPrawdy(*z.MirrorMargins)
		zmian++
	}
	if z.Columns != nil {
		if *z.Columns < 1 {
			return shared.StudioPageSetup{}, 0, bladWskazaniaStudio(
				"kolumn musi być co najmniej jedna")
		}
		wynik.Columns = postacWskaznikLiczby(*z.Columns)
		zmian++
	}
	if z.ColumnGapMm != nil {
		if *z.ColumnGapMm < 0 {
			return shared.StudioPageSetup{}, 0, bladWskazaniaStudio(
				"odstęp między kolumnami nie może być ujemny")
		}
		wynik.ColumnGapMm = postacWskaznikMiary(*z.ColumnGapMm)
		zmian++
	}
	if z.ColumnRule != nil {
		wynik.ColumnRule = postacWskaznikPrawdy(*z.ColumnRule)
		zmian++
	}
	return wynik, zmian, nil
}

// stronaBilansUkladu przelicza układ pod nowy nośnik zamiast obcinać treść:
// tabela i obraz szerszy niż kolumna tekstu schodzą do szerokości użytkowej
// z zachowaniem proporcji, a każde przeliczenie wchodzi do bilansu.
func stronaBilansUkladu(forma *shared.StudioDocumentForm,
	nastawy *shared.StudioPageSetup) []shared.StudioSkippedItem {

	uzytkowa := stronaSzerokoscUzytkowa(nastawy)
	if uzytkowa <= 0 {
		return []shared.StudioSkippedItem{{
			Reason: "marginesy szersze niż nośnik",
			Detail: postacWskaznikTekstu("po odjęciu marginesów i marginesu na oprawę " +
				"na treść nie zostaje ani milimetr — zmniejsz marginesy albo weź nośnik większy"),
		}}
	}
	pominiete := make([]shared.StudioSkippedItem, 0, 4)
	for i := range forma.Tables {
		tabela := &forma.Tables[i]
		szerokosc := 0.0
		if tabela.WidthMm != nil {
			szerokosc = *tabela.WidthMm
		}
		// Szerokość prawdziwa tabeli jest większą z dwóch: podanej wprost i sumy
		// szerokości kolumn.
		if suma := sumaMiarStrony(tabela.ColumnWidthsMm); suma > szerokosc {
			szerokosc = suma
		}
		if szerokosc <= uzytkowa+0.5 {
			continue
		}
		// Kolumny schodzą w tej samej proporcji, w jakiej stały, zamiast obcinać
		// się na krawędzi nośnika.
		wspolczynnik := uzytkowa / szerokosc
		for numer := range tabela.ColumnWidthsMm {
			tabela.ColumnWidthsMm[numer] *= wspolczynnik
		}
		tabela.WidthMm = postacWskaznikMiary(uzytkowa)
		if len(tabela.ColumnWidthsMm) == 0 && tabela.Columns > 0 {
			// Tabela bez policzonych kolumn dostaje siatkę równą.
			tabela.ColumnWidthsMm = make([]float64, tabela.Columns)
			for numer := range tabela.ColumnWidthsMm {
				tabela.ColumnWidthsMm[numer] = uzytkowa / float64(tabela.Columns)
			}
		}
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "tabela szersza niż kolumna tekstu nowego nośnika",
			Detail: postacWskaznikTekstu("tabela " + tabela.Id + " miała " +
				stronaZapisMiary(szerokosc) + " mm szerokości, a kolumna tekstu ma " +
				stronaZapisMiary(uzytkowa) + " mm; szerokości kolumn zostały " +
				"przeliczone w tej samej proporcji — sprawdź, czy treść komórek " +
				"nadal się w nich mieści"),
		})
	}
	for i := range forma.Objects {
		obiekt := &forma.Objects[i]
		if obiekt.WidthMm == nil || *obiekt.WidthMm <= uzytkowa+0.5 {
			continue
		}
		bylo := *obiekt.WidthMm
		wspolczynnik := uzytkowa / bylo
		obiekt.WidthMm = postacWskaznikMiary(uzytkowa)
		// Wysokość idzie za szerokością: obraz przeskalowany w jednej osi
		// przestaje być tym obrazem.
		if obiekt.HeightMm != nil {
			obiekt.HeightMm = postacWskaznikMiary(*obiekt.HeightMm * wspolczynnik)
		}
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "obiekt szerszy niż kolumna tekstu nowego nośnika",
			Detail: postacWskaznikTekstu("obiekt " + obiekt.Id + " miał " +
				stronaZapisMiary(bylo) + " mm szerokości, a kolumna tekstu ma " +
				stronaZapisMiary(uzytkowa) + " mm; rozmiar został zmniejszony " +
				"z zachowaniem proporcji"),
			RangeStart: obiekt.AnchorOffset,
			RangeEnd:   obiekt.AnchorOffset,
		})
	}
	return pominiete
}

func sumaMiarStrony(miary []float64) float64 {
	suma := 0.0
	for _, miara := range miary {
		suma += miara
	}
	return suma
}

// stronaZapisMiary zapisuje milimetry z jednym miejscem po przecinku — tyle,
// ile Operator widzi na linijce, i ani cyfry więcej.
func stronaZapisMiary(wartosc float64) string {
	calosc := int(wartosc)
	dziesiate := int((wartosc-float64(calosc))*10 + 0.5)
	if dziesiate >= 10 {
		calosc++
		dziesiate = 0
	}
	if dziesiate == 0 {
		return postacZapisLiczby(calosc)
	}
	return postacZapisLiczby(calosc) + "," + postacZapisLiczby(dziesiate)
}

// ── Sekcje ──────────────────────────────────────────────────────────────────

// stronaSekcjaPoKodzie znajduje sekcję postaci po identyfikatorze, albo
// oddaje wskaźnik pusty, gdy dokument sekcji o takim wskazaniu nie niesie.
func stronaSekcjaPoKodzie(forma *shared.StudioDocumentForm, kod string) *shared.StudioSection {
	szukany := strings.TrimSpace(kod)
	for i := range forma.Sections {
		if forma.Sections[i].Id == szukany {
			return &forma.Sections[i]
		}
	}
	return nil
}

// stronaSekcjaZadania rozstrzyga, na której sekcji czynność stoi. Wskazanie
// sekcji, której dokument nie ma, jest odmową nazwaną; brak wskazania znaczy
// sekcję pierwszą, zakładaną na całą treść, gdy dokument sekcji jeszcze nie ma.
func (a *adapterStudia) stronaSekcjaZadania(ctx context.Context, stan *stanPostaci,
	kod *string) (*shared.StudioSection, error) {

	if kod != nil && strings.TrimSpace(*kod) != "" {
		sekcja := stronaSekcjaPoKodzie(&stan.forma, *kod)
		if sekcja == nil {
			return nil, bladWskazaniaStudio("sekcji „" + strings.TrimSpace(*kod) +
				"” dokument " + stan.dokument.Kod + " nie ma; wykaz sekcji oddaje " +
				"komenda studio.section.list")
		}
		return sekcja, nil
	}
	if len(stan.forma.Sections) > 0 {
		sort.SliceStable(stan.forma.Sections, func(i, j int) bool {
			return stan.forma.Sections[i].Index < stan.forma.Sections[j].Index
		})
		return &stan.forma.Sections[0], nil
	}
	sekcja := shared.StudioSection{
		Id:         nowyIdentyfikator(przedrostekSekcjiPostaci),
		Index:      0,
		Title:      postacWskaznikTekstu(stronaNazwaSekcjiPierwszej),
		RangeStart: 0,
		RangeEnd:   postacDlugosc(&stan.forma),
	}
	rozpoczecie := shared.StudioSectionStart(shared.StudioSectionStartContinuous)
	sekcja.Start = &rozpoczecie
	if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, sekcja); err != nil {
		return nil, err
	}
	stan.forma.Sections = append(stan.forma.Sections, sekcja)
	return &stan.forma.Sections[len(stan.forma.Sections)-1], nil
}

// stronaZapiszSekcje utrwala sekcję wierszem warstwy danych, a nie miejscem
// w drzewie postaci, bo po sekcji się pyta, więc zapis postaci wycina sekcje
// z drzewa.
func (a *adapterStudia) stronaZapiszSekcje(ctx context.Context, dokumentID int64,
	sekcja shared.StudioSection) error {

	skladnica, err := a.postacSkladnica()
	if err != nil {
		return err
	}
	wiersz := dane.SekcjaDokumentuStudia{
		Kod:        sekcja.Id,
		DokumentID: dokumentID,
		Kolejnosc:  int64(sekcja.Index),
		Tytul:      sekcja.Title,
		ZakresOd:   int64(sekcja.RangeStart),
		ZakresDo:   int64(sekcja.RangeEnd),
	}
	if sekcja.Start != nil && strings.TrimSpace(string(*sekcja.Start)) != "" {
		wiersz.Rozpoczecie = string(*sekcja.Start)
	}
	if sekcja.PageSetup != nil {
		zapis, err := json.Marshal(sekcja.PageSetup)
		if err != nil {
			return postacBladZaplecza("nastaw strony sekcji nie da się zapisać: " + err.Error())
		}
		wiersz.NastawyStronyJSON = postacWskaznikTekstu(string(zapis))
	}
	if len(sekcja.HeadersFooters) > 0 {
		zapis, err := json.Marshal(sekcja.HeadersFooters)
		if err != nil {
			return postacBladZaplecza("nagłówków sekcji nie da się zapisać: " + err.Error())
		}
		wiersz.NaglowkiJSON = postacWskaznikTekstu(string(zapis))
	}
	if sekcja.Numbering != nil {
		zapis, err := json.Marshal(sekcja.Numbering)
		if err != nil {
			return postacBladZaplecza("numeracji stron sekcji nie da się zapisać: " + err.Error())
		}
		wiersz.NumeracjaJSON = postacWskaznikTekstu(string(zapis))
	}
	if sekcja.Watermark != nil {
		zapis, err := json.Marshal(sekcja.Watermark)
		if err != nil {
			return postacBladZaplecza("znaku wodnego sekcji nie da się zapisać: " + err.Error())
		}
		wiersz.ZnakWodnyJSON = postacWskaznikTekstu(string(zapis))
	}
	if _, err := skladnica.ZapiszSekcje(ctx, wiersz); err != nil {
		return bladStudio(err)
	}
	return nil
}

// stronaOdswierzSekcje wczytuje sekcje z warstwy danych na powrót do postaci —
// po zapisie odpowiedź ma nieść to, co stoi w bazie, a nie to, co rdzeń zamierzał
// tam zapisać.
func (a *adapterStudia) stronaOdswierzSekcje(ctx context.Context, stan *stanPostaci) error {
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return err
	}
	wiersze, err := skladnica.Sekcje(ctx, stan.dokument.ID)
	if err != nil {
		return bladStudio(err)
	}
	stan.forma.Sections = postacZlozSekcje(wiersze)
	return nil
}

// ── Nagłówek i stopka ───────────────────────────────────────────────────────

// stronaZasiegNaglowka sprawdza zasięg nagłówka. Zasięg jest tu rzeczą
// rozstrzygającą — całe to zlecenie istnieje między innymi dlatego, że dokument
// miał JEDNO pole nagłówka na wszystko.
func stronaZasiegNaglowka(zasieg shared.StudioHeaderScope) (shared.StudioHeaderScope, error) {
	if strings.TrimSpace(string(zasieg)) == "" {
		return shared.StudioHeaderScope(shared.StudioHeaderScopeDefault), nil
	}
	for _, znany := range shared.WartosciStudioHeaderScope() {
		if znany == zasieg {
			return zasieg, nil
		}
	}
	nazwy := make([]string, 0, 3)
	for _, znany := range shared.WartosciStudioHeaderScope() {
		nazwy = append(nazwy, string(znany))
	}
	return "", bladWskazaniaStudio("zasięgu nagłówka „" + string(zasieg) +
		"” kontrakt nie zna; zasięgi: " + strings.Join(nazwy, ", "))
}

// stronaScalNaglowek wnosi do nagłówka zasięgu wyłącznie pola podane w żądaniu
// i oddaje wykaz nagłówków sekcji wraz ze zmienionym.
func stronaScalNaglowek(zastane []shared.StudioHeaderFooter,
	z shared.StudioPageHeaderfooterSetRequest,
	zasieg shared.StudioHeaderScope) ([]shared.StudioHeaderFooter, int) {

	wykaz := append([]shared.StudioHeaderFooter(nil), zastane...)
	wskazanie := -1
	for i := range wykaz {
		if wykaz[i].Scope == zasieg {
			wskazanie = i
			break
		}
	}
	if wskazanie < 0 {
		wykaz = append(wykaz, shared.StudioHeaderFooter{Scope: zasieg})
		wskazanie = len(wykaz) - 1
	}
	naglowek := &wykaz[wskazanie]
	zmian := 0
	if z.HeaderText != nil {
		// Napis pusty zdejmuje nagłówek zamiast zostawiać stan poprzedni.
		if strings.TrimSpace(*z.HeaderText) == "" {
			naglowek.HeaderText = nil
		} else {
			naglowek.HeaderText = postacWskaznikTekstu(*z.HeaderText)
		}
		zmian++
	}
	if z.FooterText != nil {
		if strings.TrimSpace(*z.FooterText) == "" {
			naglowek.FooterText = nil
		} else {
			naglowek.FooterText = postacWskaznikTekstu(*z.FooterText)
		}
		zmian++
	}
	if z.HeaderDistanceMm != nil {
		naglowek.HeaderDistanceMm = postacWskaznikMiary(*z.HeaderDistanceMm)
		zmian++
	}
	if z.FooterDistanceMm != nil {
		naglowek.FooterDistanceMm = postacWskaznikMiary(*z.FooterDistanceMm)
		zmian++
	}
	if z.LinkedToPrevious != nil {
		naglowek.LinkedToPrevious = postacWskaznikPrawdy(*z.LinkedToPrevious)
		zmian++
	}
	sort.SliceStable(wykaz, func(i, j int) bool {
		return stronaKolejnoscZasiegu(wykaz[i].Scope) < stronaKolejnoscZasiegu(wykaz[j].Scope)
	})
	return wykaz, zmian
}

// stronaKolejnoscZasiegu ustala kolejność zasięgów w wykazie: strony zwykłe,
// pierwsza strona, strony parzyste — tak, jak stoją w kontrakcie i jak Operator
// czyta je w oknie nastaw.
func stronaKolejnoscZasiegu(zasieg shared.StudioHeaderScope) int {
	for i, znany := range shared.WartosciStudioHeaderScope() {
		if znany == zasieg {
			return i
		}
	}
	return len(shared.WartosciStudioHeaderScope())
}

// stronaNaglowkiZNastaw przekłada pola `header` i `footer` nastaw strony na
// nagłówek zasięgu zwykłego, żeby wykaz sekcji bez nagłówka własnego domykał
// się z tych pól, zamiast wracać pusty.
func stronaNaglowkiZNastaw(nastawy *shared.StudioPageSetup) []shared.StudioHeaderFooter {
	if nastawy == nil {
		return nil
	}
	maNaglowek := nastawy.Header != nil && strings.TrimSpace(*nastawy.Header) != ""
	maStopke := nastawy.Footer != nil && strings.TrimSpace(*nastawy.Footer) != ""
	if !maNaglowek && !maStopke {
		return nil
	}
	naglowek := shared.StudioHeaderFooter{
		Scope: shared.StudioHeaderScope(shared.StudioHeaderScopeDefault),
	}
	if maNaglowek {
		naglowek.HeaderText = postacWskaznikTekstu(*nastawy.Header)
	}
	if maStopke {
		naglowek.FooterText = postacWskaznikTekstu(*nastawy.Footer)
	}
	return []shared.StudioHeaderFooter{naglowek}
}

// ── Numeracja i znak wodny ──────────────────────────────────────────────────

// stronaScalNumeracje wnosi do numeracji stron wyłącznie pola podane
// w żądaniu, zostawiając pozostałe pola nastawy zastanej bez zmiany.
func stronaScalNumeracje(zastana *shared.StudioPageNumbering,
	z shared.StudioPageNumberingSetRequest) (shared.StudioPageNumbering, error) {

	wynik := shared.StudioPageNumbering{}
	if zastana != nil {
		wynik = *zastana
	}
	wynik.Enabled = z.Enabled
	if z.Format != nil && strings.TrimSpace(string(*z.Format)) != "" {
		znany := false
		for _, wartosc := range shared.WartosciStudioPageNumberFormat() {
			if wartosc == *z.Format {
				znany = true
				break
			}
		}
		if !znany {
			nazwy := make([]string, 0, 5)
			for _, wartosc := range shared.WartosciStudioPageNumberFormat() {
				nazwy = append(nazwy, string(wartosc))
			}
			return shared.StudioPageNumbering{}, bladWskazaniaStudio(
				"formatu numeru strony „" + string(*z.Format) + "” kontrakt nie zna; " +
					"formaty: " + strings.Join(nazwy, ", "))
		}
		format := *z.Format
		wynik.Format = &format
	}
	if z.StartAt != nil {
		if *z.StartAt < 0 {
			return shared.StudioPageNumbering{}, bladWskazaniaStudio(
				"numeracja stron nie zaczyna się od liczby ujemnej")
		}
		wynik.StartAt = postacWskaznikLiczby(*z.StartAt)
	}
	if z.RestartInSection != nil {
		wynik.RestartInSection = postacWskaznikPrawdy(*z.RestartInSection)
	}
	if z.ShowTotal != nil {
		wynik.ShowTotal = postacWskaznikPrawdy(*z.ShowTotal)
	}
	if z.Position != nil && strings.TrimSpace(*z.Position) != "" {
		wynik.Position = postacWskaznikTekstu(strings.TrimSpace(*z.Position))
	}
	if wynik.Format == nil {
		format := shared.StudioPageNumberFormat(shared.StudioPageNumberFormatArabic)
		wynik.Format = &format
	}
	return wynik, nil
}

// stronaScalZnakWodny wnosi do znaku wodnego pola podane w żądaniu i sprawdza,
// czy znak ma czym stanąć: rodzaj `text` bez napisu i `image` bez zasobu byłyby
// znakiem wodnym, którego nie da się narysować.
func stronaScalZnakWodny(zastany *shared.StudioWatermark,
	z shared.StudioPageWatermarkSetRequest) (shared.StudioWatermark, error) {

	wynik := shared.StudioWatermark{}
	if zastany != nil {
		wynik = *zastany
	}
	rodzajZnany := false
	for _, wartosc := range shared.WartosciStudioWatermarkKind() {
		if wartosc == z.Kind {
			rodzajZnany = true
			break
		}
	}
	if !rodzajZnany {
		nazwy := make([]string, 0, 3)
		for _, wartosc := range shared.WartosciStudioWatermarkKind() {
			nazwy = append(nazwy, string(wartosc))
		}
		return shared.StudioWatermark{}, bladWskazaniaStudio(
			"rodzaju znaku wodnego „" + string(z.Kind) + "” kontrakt nie zna; rodzaje: " +
				strings.Join(nazwy, ", "))
	}
	wynik.Kind = z.Kind
	if z.Text != nil {
		wynik.Text = postacWskaznikTekstu(*z.Text)
	}
	if z.AssetId != nil {
		wynik.AssetId = postacWskaznikTekstu(*z.AssetId)
	}
	if z.Opacity != nil {
		if *z.Opacity < 0 || *z.Opacity > 1 {
			return shared.StudioWatermark{}, bladWskazaniaStudio(
				"krycie znaku wodnego jest liczbą od zera do jednego")
		}
		wynik.Opacity = postacWskaznikMiary(*z.Opacity)
	}
	if z.AngleDeg != nil {
		wynik.AngleDeg = postacWskaznikMiary(*z.AngleDeg)
	}
	if z.Color != nil {
		wynik.Color = postacWskaznikTekstu(*z.Color)
	}
	if z.FontSizePt != nil {
		if *z.FontSizePt <= 0 {
			return shared.StudioWatermark{}, bladWskazaniaStudio(
				"stopień pisma znaku wodnego musi być większy od zera")
		}
		wynik.FontSizePt = postacWskaznikMiary(*z.FontSizePt)
	}

	switch wynik.Kind {
	case shared.StudioWatermarkKindText:
		if wynik.Text == nil || strings.TrimSpace(*wynik.Text) == "" {
			return shared.StudioWatermark{}, bladWskazaniaStudio(
				"znak wodny napisowy bez napisu — podaj pole text")
		}
		if wynik.FontSizePt == nil {
			wynik.FontSizePt = postacWskaznikMiary(72)
		}
	case shared.StudioWatermarkKindImage:
		if wynik.AssetId == nil || strings.TrimSpace(*wynik.AssetId) == "" {
			return shared.StudioWatermark{}, bladWskazaniaStudio(
				"znak wodny obrazowy bez zasobu obrazu — podaj pole assetId")
		}
	case shared.StudioWatermarkKindNone:
		// Rodzaj `none` zdejmuje znak wodny; nastawy zostają do przywrócenia.
	}
	if wynik.Opacity == nil {
		wynik.Opacity = postacWskaznikMiary(0.25)
	}
	return wynik, nil
}
