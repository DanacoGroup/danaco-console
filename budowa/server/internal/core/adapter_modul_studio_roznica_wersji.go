// Odpowiedzialność pliku: RÓŻNICA DWÓCH DOWOLNYCH WERSJI liczona na POSTACI
// (`studio.diff.form.compare`) oraz PRZENIESIENIE POJEDYNCZEGO FRAGMENTU ze
// wskazanej wersji do stanu bieżącego (`studio.diff.hunk.apply`).
//
// ── Dlaczego różnica postaci jest osobną komendą ─────────────────────────────
// `studio.diff.compare` liczy różnicę TREŚCI wierszami i to jest właściwe dla
// czerwonego i zielonego. Zmiana kroju albo wcięcia nie rusza ani jednej litery,
// więc tamten rachunek o niej MILCZY — a Właściciel żąda wprost, żeby była
// widoczna jako zmiana. Dlatego postać porównuje się cechą po cesze, obszar po
// obszarze, i oddaje wykaz z brzmieniem „przed" i „po".
//
// ── Dlaczego przeniesienie fragmentu, a nie samo patrzenie ──────────────────
// „Praktyczny sens tego widoku, nie samo patrzenie" — słowa Właściciela.
// Fragmenty numeruje `policzFragmentyRoznicy` — TEN SAM rachunek, który Operator
// widzi w oknie różnicy, więc numer fragmentu w żądaniu znaczy dokładnie ten
// fragment, na który Operator patrzył. Drugi rachunek fragmentów ponumerowałby
// je inaczej i „przenieś fragment trzeci" znaczyłoby co innego dla okna i dla
// rdzenia.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Nazwy obszarów różnicy postaci i rodzajów wpisu. Pełne nazwy, widoczne dla
// Operatora — bez kodów i numeracji.
const (
	roznicaRodzajDodane    = "dodane"
	roznicaRodzajUsuniete  = "usuniete"
	roznicaRodzajZmienione = "zmienione"

	roznicaObszarNastawStrony = "nastawy strony"
	roznicaObszarStylu        = "styl nazwany"
	roznicaObszarSekcji       = "sekcja"
	roznicaObszarObiektu      = "obiekt osadzony"
	roznicaObszarPola         = "pole dokumentu"
	roznicaObszarAparatu      = "aparat dokumentu"
	roznicaObszarPostaciZnaku = "postać znaku"
	roznicaObszarAkapitu      = "postać akapitu"
)

// PorownajPostac obsługuje `studio.diff.form.compare`.
func (a *adapterStudia) PorownajPostac(ctx context.Context,
	z shared.StudioDiffFormCompareRequest) (shared.StudioDiffFormCompareResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDiffFormCompareResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioDiffFormCompareResponse{}, err
	}

	odniesienie, err := a.roznicaPostacStrony(ctx, skladnica, stan, z.BaseVersionId, true)
	if err != nil {
		return shared.StudioDiffFormCompareResponse{}, err
	}
	porownywana, err := a.roznicaPostacStrony(ctx, skladnica, stan, z.TargetVersionId, false)
	if err != nil {
		return shared.StudioDiffFormCompareResponse{}, err
	}

	obszar := strings.TrimSpace(wartoscTekstu(z.Area))
	wpisy := roznicaZlozWpisy(odniesienie, porownywana, obszar)
	odpowiedz := shared.StudioDiffFormCompareResponse{Entries: wpisy}
	for _, wpis := range wpisy {
		switch wpis.Kind {
		case roznicaRodzajDodane:
			odpowiedz.Added++
		case roznicaRodzajUsuniete:
			odpowiedz.Removed++
		default:
			odpowiedz.Changed++
		}
	}
	return odpowiedz, nil
}

// PrzeniesFragmentRoznicy obsługuje `studio.diff.hunk.apply`.
func (a *adapterStudia) PrzeniesFragmentRoznicy(ctx context.Context,
	z shared.StudioDiffHunkApplyRequest) (shared.StudioDiffHunkApplyResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDiffHunkApplyResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioDiffHunkApplyResponse{}, err
	}
	if strings.TrimSpace(z.SourceVersionId) == "" {
		return shared.StudioDiffHunkApplyResponse{}, bladWskazaniaStudio(
			"przeniesienie fragmentu bez wskazania wersji źródłowej")
	}
	if z.HunkIndex == nil && (z.RangeStart == nil || z.RangeEnd == nil) {
		return shared.StudioDiffHunkApplyResponse{}, bladWskazaniaStudio(
			"przeniesienie fragmentu bez wskazania ani numeru fragmentu różnicy, " +
				"ani zakresu w wersji źródłowej — bez tego rdzeń nie wie, co przenieść")
	}
	wersja, err := skladnica.WersjaSzeregu(ctx, strings.TrimSpace(z.SourceVersionId))
	if err != nil {
		return shared.StudioDiffHunkApplyResponse{}, dziennikBrakWiersza(err,
			"wersja nie istnieje: "+z.SourceVersionId)
	}
	if wersja.DokumentKod != stan.dokument.Kod {
		return shared.StudioDiffHunkApplyResponse{}, bladWskazaniaStudio("wersja " +
			z.SourceVersionId + " należy do dokumentu " + wersja.DokumentKod + ", nie do " +
			stan.dokument.Kod)
	}
	if wersja.Tresc == nil {
		return shared.StudioDiffHunkApplyResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict, "moduł Studio: wersja "+wersja.Kod+" nie niesie "+
				"treści — nie ma z niej czego przenieść"))
	}
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	})

	biezaca := postacTekstFormy(&stan.forma)
	zrodlowa := *wersja.Tresc

	var nowaTresc string
	var od, do int
	if z.HunkIndex != nil {
		// Fragment wskazany numerem: bierze się go z rachunku, który Operator
		// widział, a nie z drugiego.
		zlozona, err := zlozTrescZFragmentow(biezaca, zrodlowa, []int{*z.HunkIndex})
		if err != nil {
			return shared.StudioDiffHunkApplyResponse{}, err
		}
		nowaTresc = zlozona
		od, do = roznicaZakresRozbieznosci(biezaca, nowaTresc)
	} else {
		// Fragment wskazany zakresem W WERSJI ŹRÓDŁOWEJ. Zakres liczy się
		// w ZNAKACH, tak jak nazywa go kontrakt.
		znakiZrodla := []rune(zrodlowa)
		odZrodla, doZrodla, poprawny := kontrolaZakresWTresci(znakiZrodla, *z.RangeStart, *z.RangeEnd)
		if !poprawny {
			return shared.StudioDiffHunkApplyResponse{}, bladWskazaniaStudio(
				"zakres od " + strconv.Itoa(*z.RangeStart) + " do " + strconv.Itoa(*z.RangeEnd) +
					" nie mieści się w treści wersji " + wersja.Kod)
		}
		brzmienie := string(znakiZrodla[odZrodla:doZrodla])
		znakiBiezacej := []rune(biezaca)
		od, do, poprawny = kontrolaZakresWTresci(znakiBiezacej, *z.RangeStart, *z.RangeEnd)
		if !poprawny {
			return shared.StudioDiffHunkApplyResponse{}, bladWskazaniaStudio(
				"zakres od " + strconv.Itoa(*z.RangeStart) + " do " + strconv.Itoa(*z.RangeEnd) +
					" nie mieści się w treści bieżącej dokumentu " + stan.dokument.Kod)
		}
		nowaTresc = string(znakiBiezacej[:od]) + brzmienie + string(znakiBiezacej[do:])
	}
	if nowaTresc == biezaca {
		return shared.StudioDiffHunkApplyResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict, "moduł Studio: wskazany fragment wersji "+wersja.Kod+
				" brzmi tak samo jak w dokumencie — przeniesienie nie zmieniłoby ani "+
				"jednej litery, więc nie jest przeniesieniem"))
	}

	// Blokada PRZED dotknięciem treści. Rachunek uzgodnienia jest ten sam, którym
	// jedzie zapora rejestru — fragment pod blokadą zostaje w brzmieniu zastanym,
	// a odpowiedź niesie bilans pominięć.
	blokady, err := skladnica.BlokadyFragmentow(ctx, stan.dokument.ID)
	if err != nil {
		return shared.StudioDiffHunkApplyResponse{}, bladStudio(err)
	}
	uzgodnienie := blokadaUzgodnijTresc(biezaca, nowaTresc, blokady, wykonawca)
	if uzgodnienie.CalkiemWBlokadzie {
		return shared.StudioDiffHunkApplyResponse{},
			kontrolaBladBlokady(uzgodnienie.PierwszaBlokada, wykonawca.nazwaWykonawcy())
	}
	przed := biezaca
	postacUzgodnijZTrescia(&stan.forma, uzgodnienie.Tresc)

	// Postać fragmentu jedzie razem z treścią, gdy Operator o to poprosi (brak
	// znaczy tak): przeniesienie samego brzmienia zostawiłoby akapit w kroju
	// bieżącym, a wtedy „przeniosłem fragment ze starej wersji" byłoby półprawdą.
	if (z.IncludeForm == nil || *z.IncludeForm) && wersja.PostacJSON != nil &&
		strings.TrimSpace(*wersja.PostacJSON) != "" {

		var postacZrodla shared.StudioDocumentForm
		if err := json.Unmarshal([]byte(*wersja.PostacJSON), &postacZrodla); err != nil {
			return shared.StudioDiffHunkApplyResponse{}, kontrolaBladZaplecza(
				"wersja " + wersja.Kod + " niesie nieczytelną postać: " + err.Error())
		}
		roznicaPrzeniesPostacZakresu(&stan.forma, &postacZrodla, od, do)
	}

	zmiana, err := a.postacOdlozZmiane(ctx, stan, wykonawca.Rodzaj,
		shared.StudioChangeKindWstawienie, od, do, &przed, &uzgodnienie.Tresc)
	if err != nil {
		return shared.StudioDiffHunkApplyResponse{}, err
	}
	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioDiffHunkApplyResponse{}, err
	}
	drzewo, err := dziennikZapisDrzewa(stan.forma)
	if err != nil {
		return shared.StudioDiffHunkApplyResponse{}, err
	}
	var kodZmiany *string
	if zmiana != nil {
		kodZmiany = &zmiana.Id
	}
	czynnosc, err := a.dziennikOdlozCzynnosc(ctx, stan.dokument, wykonawca,
		shared.StudioActionKindTextEdit,
		"przeniesienie fragmentu z wersji "+wersja.Kod+" do stanu bieżącego; przeniósł: "+
			wykonawca.nazwaWykonawcy(), &od, &do, drzewo, drzewo, kodZmiany)
	if err != nil {
		return shared.StudioDiffHunkApplyResponse{}, err
	}
	if err := a.zmianyModeluStempluj(ctx, zmiana, wykonawca, czynnosc); err != nil {
		return shared.StudioDiffHunkApplyResponse{}, err
	}

	bilans := blokadaBilans(uzgodnienie)
	if bilans.Applied == 0 {
		bilans.Applied = 1
	}
	if bilans.Skipped == nil {
		bilans.Skipped = []shared.StudioSkippedItem{}
	}
	if bilans.Note == nil {
		bilans.Note = kontrolaWskaznikTekstu("Fragment z wersji " + wersja.Kod +
			" wszedł do stanu bieżącego; wersje nowsze zostają, więc przeniesienie " +
			"da się cofnąć przez studio.journal.revert.")
	}
	return shared.StudioDiffHunkApplyResponse{
		Form:     stan.forma,
		Balance:  bilans,
		Change:   zmiana,
		ActionId: czynnosc,
		Document: a.zlozDokument(stan.dokument),
	}, nil
}

// ── Strony porównania ───────────────────────────────────────────────────────

// roznicaPostacStrony bierze postać jednej strony porównania.
//
// Brak wskazania wersji znaczy dla ODNIESIENIA wersję założycielską, a dla
// strony PORÓWNYWANEJ stan bieżący — tak stanowi kontrakt i tak czyta to widok
// różnicy: „co się zmieniło od stanu pierwotnego do teraz".
func (a *adapterStudia) roznicaPostacStrony(ctx context.Context, skladnica KontrolaPracyStudia,
	stan *stanPostaci, wskazanie *string, odniesienie bool) (*shared.StudioDocumentForm, error) {

	if !kontrolaTekstNiepusty(wskazanie) {
		if !odniesienie {
			forma := stan.forma
			return &forma, nil
		}
		zalozycielska, err := skladnica.WersjaZalozycielska(ctx, stan.dokument.ID)
		if err != nil {
			if errors.Is(err, dane.ErrBrakWiersza) {
				return nil, bladBrakuStudio("dokument " + stan.dokument.Kod + " nie ma wersji " +
					"założycielskiej, a wersji odniesienia nie wskazano — porównanie nie ma " +
					"z czym zestawić stanu bieżącego")
			}
			return nil, bladStudio(err)
		}
		return roznicaPostacWersji(zalozycielska)
	}
	wersja, err := skladnica.WersjaSzeregu(ctx, strings.TrimSpace(*wskazanie))
	if err != nil {
		return nil, dziennikBrakWiersza(err, "wersja nie istnieje: "+*wskazanie)
	}
	if wersja.DokumentKod != stan.dokument.Kod {
		return nil, bladWskazaniaStudio("wersja " + *wskazanie + " należy do dokumentu " +
			wersja.DokumentKod + ", nie do " + stan.dokument.Kod)
	}
	return roznicaPostacWersji(wersja)
}

// roznicaPostacWersji czyta drzewo postaci odłożone przy wersji.
//
// Wersja bez postaci NIE zamienia się w postać pustą: to pokazałoby cały arkusz
// stylów jako „usunięty", czyli skłamałoby o różnicy. Odmowa nazywa powód.
func roznicaPostacWersji(wersja dane.WersjaSzereguStudia) (*shared.StudioDocumentForm, error) {
	if wersja.PostacJSON == nil || strings.TrimSpace(*wersja.PostacJSON) == "" {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
			"moduł Studio: wersja "+wersja.Kod+" została założona bez drzewa postaci, "+
				"więc różnicy postaci nie ma z czym policzyć. Porównanie treści idzie "+
				"przez studio.diff.compare i działa dla każdej wersji."))
	}
	var forma shared.StudioDocumentForm
	if err := json.Unmarshal([]byte(*wersja.PostacJSON), &forma); err != nil {
		return nil, kontrolaBladZaplecza("wersja " + wersja.Kod +
			" niesie nieczytelną postać: " + err.Error())
	}
	return &forma, nil
}

// ── Rachunek różnicy postaci ────────────────────────────────────────────────

// roznicaZlozWpisy porównuje dwa drzewa postaci cechą po cesze.
func roznicaZlozWpisy(odniesienie, porownywana *shared.StudioDocumentForm,
	obszar string) []shared.StudioFormDiffEntry {

	wpisy := []shared.StudioFormDiffEntry{}
	dotyczy := func(nazwa string) bool {
		return obszar == "" || strings.EqualFold(obszar, nazwa)
	}

	if dotyczy(roznicaObszarNastawStrony) &&
		dziennikRoznePola(odniesienie.PageSetup, porownywana.PageSetup) {

		wpisy = append(wpisy, roznicaWpis(roznicaRodzajZmienione, roznicaObszarNastawStrony,
			"nastawy strony dokumentu różnią się — nośnik, orientacja, marginesy "+
				"albo kolumny", roznicaZapis(odniesienie.PageSetup),
			roznicaZapis(porownywana.PageSetup)))
	}

	if dotyczy(roznicaObszarStylu) {
		wStarej := map[string]shared.StudioNamedStyle{}
		for _, styl := range odniesienie.Styles {
			wStarej[styl.Name] = styl
		}
		wNowej := map[string]shared.StudioNamedStyle{}
		for _, styl := range porownywana.Styles {
			wNowej[styl.Name] = styl
		}
		nazwy := roznicaNazwyScalone(wStarej, wNowej)
		for _, nazwa := range nazwy {
			stary, byl := wStarej[nazwa]
			nowy, jest := wNowej[nazwa]
			switch {
			case byl && !jest:
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajUsuniete, roznicaObszarStylu,
					"styl nazwany „"+nazwa+"” odpadł z arkusza",
					roznicaZapis(stary), nil))
			case !byl && jest:
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajDodane, roznicaObszarStylu,
					"styl nazwany „"+nazwa+"” doszedł do arkusza",
					nil, roznicaZapis(nowy)))
			case dziennikRoznePola(stary, nowy):
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajZmienione, roznicaObszarStylu,
					"styl nazwany „"+nazwa+"” zmienił postać",
					roznicaZapis(stary), roznicaZapis(nowy)))
			}
		}
	}

	if dotyczy(roznicaObszarSekcji) {
		wStarej := map[string]shared.StudioSection{}
		for _, sekcja := range odniesienie.Sections {
			wStarej[sekcja.Id] = sekcja
		}
		wNowej := map[string]shared.StudioSection{}
		for _, sekcja := range porownywana.Sections {
			wNowej[sekcja.Id] = sekcja
		}
		for _, kod := range roznicaNazwyScalone(wStarej, wNowej) {
			stara, byla := wStarej[kod]
			nowa, jest := wNowej[kod]
			switch {
			case byla && !jest:
				wpisy = append(wpisy, roznicaWpisZakresu(roznicaRodzajUsuniete,
					roznicaObszarSekcji, "sekcja odpadła", roznicaZapis(stara), nil,
					stara.RangeStart, stara.RangeEnd))
			case !byla && jest:
				wpisy = append(wpisy, roznicaWpisZakresu(roznicaRodzajDodane,
					roznicaObszarSekcji, "sekcja doszła", nil, roznicaZapis(nowa),
					nowa.RangeStart, nowa.RangeEnd))
			case dziennikRoznePola(stara, nowa):
				wpisy = append(wpisy, roznicaWpisZakresu(roznicaRodzajZmienione,
					roznicaObszarSekcji, "sekcja zmieniła nastawy albo zakres",
					roznicaZapis(stara), roznicaZapis(nowa), nowa.RangeStart, nowa.RangeEnd))
			}
		}
	}

	if dotyczy(roznicaObszarObiektu) {
		wStarej := map[string]shared.StudioDocumentObject{}
		for _, obiekt := range odniesienie.Objects {
			wStarej[obiekt.Id] = obiekt
		}
		wNowej := map[string]shared.StudioDocumentObject{}
		for _, obiekt := range porownywana.Objects {
			wNowej[obiekt.Id] = obiekt
		}
		for _, kod := range roznicaNazwyScalone(wStarej, wNowej) {
			stary, byl := wStarej[kod]
			nowy, jest := wNowej[kod]
			switch {
			case byl && !jest:
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajUsuniete, roznicaObszarObiektu,
					"obiekt osadzony odpadł z dokumentu", roznicaZapis(stary), nil))
			case !byl && jest:
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajDodane, roznicaObszarObiektu,
					"obiekt osadzony doszedł do dokumentu", nil, roznicaZapis(nowy)))
			case dziennikRoznePola(stary, nowy):
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajZmienione, roznicaObszarObiektu,
					"obiekt osadzony zmienił postać albo położenie",
					roznicaZapis(stary), roznicaZapis(nowy)))
			}
		}
	}

	if dotyczy(roznicaObszarPola) {
		wStarej := map[string]shared.StudioDocumentField{}
		for _, pole := range odniesienie.Fields {
			wStarej[pole.Id] = pole
		}
		wNowej := map[string]shared.StudioDocumentField{}
		for _, pole := range porownywana.Fields {
			wNowej[pole.Id] = pole
		}
		for _, kod := range roznicaNazwyScalone(wStarej, wNowej) {
			stare, bylo := wStarej[kod]
			nowe, jest := wNowej[kod]
			switch {
			case bylo && !jest:
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajUsuniete, roznicaObszarPola,
					"pole dokumentu odpadło", roznicaZapis(stare), nil))
			case !bylo && jest:
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajDodane, roznicaObszarPola,
					"pole dokumentu doszło", nil, roznicaZapis(nowe)))
			case dziennikRoznePola(stare, nowe):
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajZmienione, roznicaObszarPola,
					"pole dokumentu zmieniło rodzaj, format albo wartość",
					roznicaZapis(stare), roznicaZapis(nowe)))
			}
		}
	}

	if dotyczy(roznicaObszarAparatu) &&
		dziennikRoznePola(odniesienie.Apparatus, porownywana.Apparatus) {

		wpisy = append(wpisy, roznicaWpis(roznicaRodzajZmienione, roznicaObszarAparatu,
			"aparat dokumentu różni się — spis treści, przypisy, bibliografia, "+
				"odwołania, podpisy, indeks albo zakładki",
			roznicaZapis(odniesienie.Apparatus), roznicaZapis(porownywana.Apparatus)))
	}

	// Postać znaku i akapitu — po blokach. To jest właśnie ta różnica, o której
	// rachunek treści milczy: zmiana kroju nie rusza ani jednej litery.
	if dotyczy(roznicaObszarPostaciZnaku) || dotyczy(roznicaObszarAkapitu) {
		wStarej := dziennikBlokiPoKodzie(odniesienie.Blocks)
		wNowej := dziennikBlokiPoKodzie(porownywana.Blocks)
		for _, kod := range roznicaNazwyScalone(wStarej, wNowej) {
			stary, byl := wStarej[kod]
			nowy, jest := wNowej[kod]
			if !byl || !jest {
				// Blok, który doszedł albo odpadł, jest zmianą TREŚCI i widać go
				// w `studio.diff.compare`. Powtarzanie go tutaj mnożyłoby jedną
				// zmianę na dwa wykazy.
				continue
			}
			if dotyczy(roznicaObszarAkapitu) && dziennikRoznePola(stary.Paragraph, nowy.Paragraph) {
				wpisy = append(wpisy, roznicaWpis(roznicaRodzajZmienione, roznicaObszarAkapitu,
					"akapit zmienił postać — wyrównanie, wcięcie, odstęp, interlinia, "+
						"obramowanie albo cieniowanie",
					roznicaZapis(stary.Paragraph), roznicaZapis(nowy.Paragraph)))
			}
			if dotyczy(roznicaObszarPostaciZnaku) &&
				dziennikRoznePola(roznicaPostacieRunow(stary), roznicaPostacieRunow(nowy)) {

				wpisy = append(wpisy, roznicaWpis(roznicaRodzajZmienione,
					roznicaObszarPostaciZnaku,
					"fragment zmienił postać znaku — krój, stopień, grubość, odmiana, "+
						"podkreślenie, barwa albo wyróżnienie",
					roznicaZapis(roznicaPostacieRunow(stary)),
					roznicaZapis(roznicaPostacieRunow(nowy))))
			}
		}
	}
	return wpisy
}

// roznicaNazwyScalone oddaje klucze obu map w jednym, ustalonym porządku.
//
// Porządek jest ustalony (posortowany), bo wykaz różnicy oglądany dwa razy ma
// wyglądać tak samo: kolejność wzięta z przebiegu mapy zmieniałaby się przy
// każdym wywołaniu i Operator nie odnalazłby pozycji, na którą patrzył.
func roznicaNazwyScalone[T any](pierwsza, druga map[string]T) []string {
	zbior := make(map[string]bool, len(pierwsza)+len(druga))
	for klucz := range pierwsza {
		zbior[klucz] = true
	}
	for klucz := range druga {
		zbior[klucz] = true
	}
	nazwy := make([]string, 0, len(zbior))
	for klucz := range zbior {
		nazwy = append(nazwy, klucz)
	}
	sort.Strings(nazwy)
	return nazwy
}

// roznicaPostacieRunow zbiera postacie znaku wszystkich runów bloku.
func roznicaPostacieRunow(blok shared.StudioDocumentBlock) []*shared.StudioCharacterFormat {
	postacie := make([]*shared.StudioCharacterFormat, 0, len(blok.Runs))
	for _, run := range blok.Runs {
		postacie = append(postacie, run.Format)
	}
	return postacie
}

// roznicaZapis składa czytelny zapis stanu cechy do wykazu różnicy.
func roznicaZapis(wartosc any) *string {
	if wartosc == nil {
		return nil
	}
	zapis, err := json.Marshal(wartosc)
	if err != nil {
		return nil
	}
	tekst := string(zapis)
	if tekst == "null" {
		return nil
	}
	return &tekst
}

// roznicaWpis składa pozycję wykazu różnicy postaci.
func roznicaWpis(rodzaj, obszar, opis string, przed, po *string) shared.StudioFormDiffEntry {
	return shared.StudioFormDiffEntry{
		Kind: rodzaj, Area: obszar, Detail: opis, Before: przed, After: po,
	}
}

// roznicaWpisZakresu składa pozycję wykazu wraz z zakresem w treści.
func roznicaWpisZakresu(rodzaj, obszar, opis string, przed, po *string,
	od, do int) shared.StudioFormDiffEntry {

	wpis := roznicaWpis(rodzaj, obszar, opis, przed, po)
	poczatek, koniec := od, do
	wpis.RangeStart, wpis.RangeEnd = &poczatek, &koniec
	return wpis
}

// roznicaZakresRozbieznosci oddaje zakres znaków, w którym dwie treści się
// rozchodzą. Liczy w ZNAKACH, tak jak nazywa zakresy kontrakt.
func roznicaZakresRozbieznosci(przed, po string) (int, int) {
	znakiPrzed, znakiPo := []rune(przed), []rune(po)
	przedrostek := 0
	for przedrostek < len(znakiPrzed) && przedrostek < len(znakiPo) &&
		znakiPrzed[przedrostek] == znakiPo[przedrostek] {

		przedrostek++
	}
	sufiks := 0
	for sufiks < len(znakiPrzed)-przedrostek && sufiks < len(znakiPo)-przedrostek &&
		znakiPrzed[len(znakiPrzed)-1-sufiks] == znakiPo[len(znakiPo)-1-sufiks] {

		sufiks++
	}
	return przedrostek, len(znakiPrzed) - sufiks
}

// roznicaPrzeniesPostacZakresu przenosi postać znaku i akapitu bloków objętych
// zakresem ze wersji źródłowej do stanu bieżącego.
//
// Tożsamością bloku jest jego identyfikator — blok, którego wersja źródłowa nie
// zna, zostaje w postaci bieżącej. Zgadywanie odpowiedniości bloków byłoby tu
// gorsze niż nieprzeniesienie postaci, bo nałożyłoby krój obcego akapitu.
func roznicaPrzeniesPostacZakresu(biezaca, zrodlowa *shared.StudioDocumentForm, od, do int) {
	wZrodle := dziennikBlokiPoKodzie(zrodlowa.Blocks)
	postacPrzeliczZakresy(biezaca)
	for _, numer := range postacBlokiZakresu(biezaca, od, do) {
		zrodlowy, jest := wZrodle[biezaca.Blocks[numer].Id]
		if !jest {
			continue
		}
		biezaca.Blocks[numer].Paragraph = zrodlowy.Paragraph
		// Postać runów przenosi się WYŁĄCZNIE wtedy, gdy brzmienie bloku jest to
		// samo: run niesie zarówno tekst, jak i postać, a przy różnym brzmieniu
		// podstawienie runów źródłowych podmieniłoby treść pod pozorem postaci.
		if postacTekstBloku(biezaca.Blocks[numer]) == postacTekstBloku(zrodlowy) {
			biezaca.Blocks[numer].Runs = zrodlowy.Runs
		}
	}
}
