// Odpowiedzialność pliku: operacja kontekstowa Tools Panel (`studio.contextual.op`)
// i porównanie Diff/Grep Panel (`studio.diff.compare`) — druga połowa portu Studio.
package core

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// przedrostekPropozycjiStudio nadaje identyfikator zewnętrzny propozycji zmiany.
// Przedrostki dokumentu i wersji stoją w `adapter_modul_studio.go`; propozycja
// jest bytem tego pliku, więc przedrostek nadaje ten plik.
const przedrostekPropozycjiStudio = "studio-prop-"

// przedrostekZmianyModeluStudia nadaje identyfikator zewnętrzny zmianie
// śledzonej odłożonej przez operację kontekstową; przedrostek nadaje ten
// plik, bo tutaj ona powstaje.
const przedrostekZmianyModeluStudia = "studio-zm-"

// OperacjaKontekstowa wykonuje operację kontekstową silnikiem modelu, odkłada
// jej wynik jako propozycję zmiany oraz wpisuje go do treści dokumentu jako
// zmianę śledzoną autora `model`.
func (a *adapterStudia) OperacjaKontekstowa(ctx context.Context,
	z shared.StudioContextualOpRequest) (shared.StudioContextualOpResponse, error) {

	if z.WindowId == "" || z.DocumentId == "" || z.ActionId == "" {
		return shared.StudioContextualOpResponse{}, bladWskazaniaStudio(
			"operacja kontekstowa wymaga okna, dokumentu i pozycji rejestru akcji")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioContextualOpResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	if a.kanaly == nil || a.okna == nil {
		return shared.StudioContextualOpResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable,
			"moduł Studio: rdzeń złożony bez rejestru kanałów modelu — operacja "+
				z.ActionId+" nie ma czym się wykonać"))
	}
	okno, err := a.okna.Okno(z.WindowId)
	if err != nil {
		return shared.StudioContextualOpResponse{}, bladSesji(err)
	}
	if okno.KanalModelu == "" {
		return shared.StudioContextualOpResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable,
			"moduł Studio: okno "+z.WindowId+" nie ma wskazanego kanału modelu"))
	}

	var wynik strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			wynik.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	zapytanie := models.Zapytanie{
		Zasiegi:   models.Zasiegi{Okno: z.WindowId},
		Wiadomosc: z.ActionId,
		Tresc:     trescOperacjiStudia(z, dokument),
		Kanal:     okno.KanalModelu,
	}
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		return shared.StudioContextualOpResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable,
			"moduł Studio: operacja "+z.ActionId+" nie doszła do skutku: "+err.Error()))
	}
	tresc := wynik.String()
	if tresc == "" {
		return shared.StudioContextualOpResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable,
			"moduł Studio: silnik modelu nie oddał ani jednego fragmentu treści dla "+z.ActionId))
	}

	// Propozycja powstaje po wykonaniu, nie przed nim, i niesie wynik gotowy.
	propozycja, err := a.repozytorium.ZapiszPropozycje(ctx, dokument.ID, dane.PropozycjaZmiany{
		IdentyfikatorZewnetrzny: nowyIdentyfikator(przedrostekPropozycjiStudio),
		AkcjaID:                 z.ActionId,
		TrescWyniku:             &tresc,
	})
	if err != nil {
		return shared.StudioContextualOpResponse{}, bladStudio(err)
	}
	kod := propozycja.IdentyfikatorZewnetrzny

	if err := a.odlozWynikModeluStudia(ctx, dokument, z, tresc, kod); err != nil {
		return shared.StudioContextualOpResponse{}, err
	}
	return shared.StudioContextualOpResponse{ProposalId: &kod, ResultText: &tresc}, nil
}

// odlozWynikModeluStudia wpisuje wynik operacji do treści dokumentu i rejestruje
// go jako zmianę śledzoną autora `model`. Zakres liczony jest w znakach, nie
// w bajtach — tak nazywa go kontrakt i tak czyta go decyzja o zmianach.
func (a *adapterStudia) odlozWynikModeluStudia(ctx context.Context,
	dokument dane.DokumentStudia, z shared.StudioContextualOpRequest,
	wynik, kodPropozycji string) error {

	stara := []rune(wartoscTekstu(dokument.Tresc))
	od, do_ := zakresOperacjiStudia(z, len(stara))
	nowa := string(stara[:od]) + wynik + string(stara[do_:])
	przed := string(stara[od:do_])

	rodzaj := shared.StudioChangeKindWstawienie
	if wynik == "" {
		rodzaj = shared.StudioChangeKindUsuniecie
	}
	if _, err := a.repozytorium.ZapiszZmianeSledzona(ctx, dokument.ID, dane.ZmianaSledzona{
		Kod:        nowyIdentyfikator(przedrostekZmianyModeluStudia),
		Rodzaj:     string(rodzaj),
		Autor:      string(shared.StudioAuthorModel),
		ZakresOd:   int64(od),
		ZakresDo:   int64(od + len([]rune(wynik))),
		TrescPrzed: &przed,
		TrescPo:    &wynik,
		Decyzja:    "oczekuje",
	}); err != nil {
		return bladStudio(err)
	}

	dokument.Tresc = &nowa
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return bladStudio(err)
	}
	if _, err := a.zalozWersjeDokumentu(ctx, zapisany, nowa,
		shared.StudioAuthorModel, &kodPropozycji); err != nil {
		return err
	}
	return nil
}

// zakresOperacjiStudia oddaje zakres, na którym operacja pracowała, w znakach.
// Zakres poza treścią i zakres odwrócony schodzą na cały dokument, a nie na
// odmowę: żądanie już się wykonało i wynik modelu jest w ręku.
func zakresOperacjiStudia(z shared.StudioContextualOpRequest, dlugosc int) (int, int) {
	if z.Scope != shared.StudioOperationScopeSelection ||
		z.SelectionStart == nil || z.SelectionEnd == nil {
		return 0, dlugosc
	}
	od, do_ := *z.SelectionStart, *z.SelectionEnd
	if od < 0 || do_ > dlugosc || od > do_ {
		return 0, dlugosc
	}
	return od, do_
}

// trescOperacjiStudia składa polecenie dla modelu: czynność z rejestru akcji
// wraz z treścią dokumentu, na której ma pracować.
func trescOperacjiStudia(z shared.StudioContextualOpRequest, dokument dane.DokumentStudia) string {
	czesci := []string{"Czynnosc: " + z.ActionId}
	// Parametry operacji jadą do modelu, nie tylko do bazy.
	if len(z.Params) > 0 && strings.TrimSpace(string(z.Params)) != "null" {
		czesci = append(czesci, "Nastawy i polecenie Operatora:", string(z.Params))
	}
	if dokument.Tresc != nil && *dokument.Tresc != "" {
		// Wycinek bierze się po runach, nie po bajtach.
		runy := []rune(*dokument.Tresc)
		if z.SelectionStart != nil && z.SelectionEnd != nil {
			od, do_ := *z.SelectionStart, *z.SelectionEnd
			if od >= 0 && do_ > od && do_ <= len(runy) {
				czesci = append(czesci, "Zaznaczenie:", string(runy[od:do_]))
			}
		}
		czesci = append(czesci, "Dokument:", *dokument.Tresc)
	}
	return strings.Join(czesci, "\n\n")
}

// Porownaj obsługuje `studio.diff.compare`: liczy fragmenty różnicy między
// dwiema treściami i, niezależnie, wyszukuje wzorzec w treści porównywanej
// strony.
func (a *adapterStudia) Porownaj(ctx context.Context,
	z shared.StudioDiffCompareRequest) (shared.StudioDiffCompareResponse, error) {

	if z.DocumentId == "" {
		return shared.StudioDiffCompareResponse{}, bladWskazaniaStudio(
			"porównanie bez wskazania dokumentu")
	}
	if _, err := a.repozytorium.Dokument(ctx, z.DocumentId); err != nil {
		return shared.StudioDiffCompareResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}

	baza, idBazy, err := a.trescStrony(ctx, z.BaseVersionId, nil)
	if err != nil {
		return shared.StudioDiffCompareResponse{}, err
	}
	cel, idCelu, err := a.trescStrony(ctx, z.TargetVersionId, z.ProposalId)
	if err != nil {
		return shared.StudioDiffCompareResponse{}, err
	}

	// Żądanie bez fragmentów i bez trafień wzorca jest żądaniem bez odpowiedzi.
	maFragmenty := idBazy != "" && idCelu != ""
	maWzorzec := z.Pattern != nil && strings.TrimSpace(*z.Pattern) != ""
	maTrafienia := maWzorzec && (idBazy != "" || idCelu != "")
	if !maFragmenty && !maTrafienia {
		return shared.StudioDiffCompareResponse{},
			bladWskazaniaStudio(brakStronPorownania(idBazy, idCelu, maWzorzec))
	}

	odpowiedz := shared.StudioDiffCompareResponse{}
	if maFragmenty {
		odpowiedz.Hunks = policzFragmentyRoznicy(baza, cel)
	}
	if maWzorzec {
		trescDoSzukania, idStrony := cel, idCelu
		if idStrony == "" {
			trescDoSzukania, idStrony = baza, idBazy
		}
		dopasowania, err := szukajWzorca(trescDoSzukania, *z.Pattern, z.Regex, idStrony)
		if err != nil {
			return shared.StudioDiffCompareResponse{}, bladWskazaniaStudio(err.Error())
		}
		odpowiedz.Matches = dopasowania
	}
	return odpowiedz, nil
}

// trescStrony czyta treść jednej strony porównania — wersji albo, gdy podano,
// propozycji zamiast wersji docelowej. Strona pominięta oddaje pusty
// identyfikator.
func (a *adapterStudia) trescStrony(ctx context.Context, kodWersji, kodPropozycji *string) (string, string, error) {
	if kodPropozycji != nil && strings.TrimSpace(*kodPropozycji) != "" {
		propozycja, err := a.repozytorium.Propozycja(ctx, *kodPropozycji)
		if err != nil {
			return "", "", bladWskazaniaStronyPorownania(*kodPropozycji, err)
		}
		return trescBytuPorownania(*kodPropozycji, propozycja.Tresc, propozycja.TrescOdwolanie)
	}
	if kodWersji != nil && strings.TrimSpace(*kodWersji) != "" {
		wersja, err := a.repozytorium.Wersja(ctx, *kodWersji)
		if err != nil {
			return "", "", bladWskazaniaStronyPorownania(*kodWersji, err)
		}
		return trescBytuPorownania(*kodWersji, wersja.Tresc, wersja.TrescOdwolanie)
	}
	return "", "", nil
}

// trescBytuPorownania oddaje treść krótką wprost. Treść obszerną, przechowaną
// poza bazą przez `trescOdwolanie`, adapter Studio nie potrafi tu doczytać —
// odmawia wprost zamiast oddać porównanie połowy treści.
func trescBytuPorownania(kod string, tresc, odwolanie *string) (string, string, error) {
	if tresc != nil {
		return *tresc, kod, nil
	}
	if odwolanie != nil && *odwolanie != "" {
		return "", "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: treść porównania leży poza bazą (odwołanie "+*odwolanie+
				"), rdzeń nie ma tu mechanizmu jej odczytu"))
	}
	return "", kod, nil
}

// brakStronPorownania nazywa to, czego w żądaniu zabrakło, polami kontraktu,
// zamiast mówić ogólnie, że czegoś brakuje.
func brakStronPorownania(idBazy, idCelu string, maWzorzec bool) string {
	if idBazy == "" && idCelu == "" {
		if maWzorzec {
			return "wyszukanie wzorca bez wskazania strony przeszukiwanej — " +
				"wzorzec szuka w wersji albo propozycji, nie w bieżącej treści dokumentu; " +
				"wskaż baseVersionId albo targetVersionId lub proposalId"
		}
		return "porównanie bez wskazania stron — potrzebne są baseVersionId " +
			"oraz targetVersionId albo proposalId"
	}
	if idCelu == "" {
		return "porównanie bez strony porównywanej — wskaż targetVersionId albo proposalId"
	}
	return "porównanie bez strony odniesienia — wskaż baseVersionId"
}

// bladWskazaniaStronyPorownania odróżnia brak wiersza od usterki wewnętrznej
// przy czytaniu strony porównania (wzór: `bladNieznanegoDokumentu`).
func bladWskazaniaStronyPorownania(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: wersja albo propozycja "+kod+" nie istnieje"))
	}
	return bladStudio(err)
}

// policzFragmentyRoznicy liczy fragmenty różnicy dwóch treści wierszami —
// przycina wspólny przedrostek i wspólny sufiks, a to, co zostaje pomiędzy,
// jest jednym fragmentem zmiany.
func policzFragmentyRoznicy(bazowa, docelowa string) []shared.StudioDiffHunk {
	przed := podzielNaWiersze(bazowa)
	po := podzielNaWiersze(docelowa)

	przedrostek := 0
	for przedrostek < len(przed) && przedrostek < len(po) && przed[przedrostek] == po[przedrostek] {
		przedrostek++
	}
	sufiks := 0
	for sufiks < len(przed)-przedrostek && sufiks < len(po)-przedrostek &&
		przed[len(przed)-1-sufiks] == po[len(po)-1-sufiks] {
		sufiks++
	}

	hunki := []shared.StudioDiffHunk{}
	numer := 1
	if przedrostek > 0 {
		hunki = append(hunki, kontekstFragmentu(numer, przed[:przedrostek], 1))
		numer++
	}

	srodekPrzed := przed[przedrostek : len(przed)-sufiks]
	srodekPo := po[przedrostek : len(po)-sufiks]
	if len(srodekPrzed) > 0 || len(srodekPo) > 0 {
		hunki = append(hunki, fragmentZmiany(numer, srodekPrzed, srodekPo, przedrostek+1))
		numer++
	}

	if sufiks > 0 {
		hunki = append(hunki, kontekstFragmentu(numer, przed[len(przed)-sufiks:], len(przed)-sufiks+1))
	}
	return hunki
}

// kontekstFragmentu składa fragment wierszy niezmienionych — treść przed i po
// jest ta sama, więc kontrakt niesie ją raz, w polu `Before`.
func kontekstFragmentu(numer int, wiersze []string, wierszPoczatkowy int) shared.StudioDiffHunk {
	tresc := strings.Join(wiersze, "\n")
	start, koniec := wierszPoczatkowy, wierszPoczatkowy+len(wiersze)-1
	return shared.StudioDiffHunk{
		Index: numer, Kind: shared.DiffHunkKindContext,
		Before: &tresc, StartLine: &start, EndLine: &koniec,
	}
}

// fragmentZmiany składa jeden fragment różnicy między wersją bazową i
// docelową. Rodzaj zależy od tego, która strona ma tu wiersze.
func fragmentZmiany(numer int, przed, po []string, wierszPoczatkowy int) shared.StudioDiffHunk {
	var rodzaj shared.DiffHunkKind = shared.DiffHunkKindChanged
	switch {
	case len(przed) == 0:
		rodzaj = shared.DiffHunkKindAdded
	case len(po) == 0:
		rodzaj = shared.DiffHunkKindRemoved
	}
	fragment := shared.StudioDiffHunk{Index: numer, Kind: rodzaj}
	if len(przed) > 0 {
		tresc := strings.Join(przed, "\n")
		fragment.Before = &tresc
	}
	if len(po) > 0 {
		tresc := strings.Join(po, "\n")
		fragment.After = &tresc
	}
	dlugosc := len(przed)
	if len(po) > dlugosc {
		dlugosc = len(po)
	}
	start, koniec := wierszPoczatkowy, wierszPoczatkowy+dlugosc-1
	fragment.StartLine, fragment.EndLine = &start, &koniec
	return fragment
}

// podzielNaWiersze rozbija treść na wiersze. Treść pusta daje wykaz pusty, nie
// jeden wiersz pusty — inaczej porównanie dwóch treści pustych pokazałoby
// fragment „zmieniono”, którego nie ma.
func podzielNaWiersze(tresc string) []string {
	if tresc == "" {
		return []string{}
	}
	return strings.Split(tresc, "\n")
}

// szukajWzorca przeszukuje treść wiersz po wierszu — dosłownie albo jako
// wyrażenie regularne, zależnie od `regex`.
func szukajWzorca(tresc, wzorzec string, jakoRegex *bool, idStrony string) ([]shared.StudioTextMatch, error) {
	if tresc == "" {
		return nil, nil
	}
	var dopasuj func(string) []int
	if jakoRegex != nil && *jakoRegex {
		wyrazenie, err := regexp.Compile(wzorzec)
		if err != nil {
			return nil, err
		}
		dopasuj = func(wiersz string) []int { return wyrazenie.FindStringIndex(wiersz) }
	} else {
		dopasuj = func(wiersz string) []int {
			pozycja := strings.Index(wiersz, wzorzec)
			if pozycja < 0 {
				return nil
			}
			return []int{pozycja, pozycja + len(wzorzec)}
		}
	}

	var idWskaznik *string
	if idStrony != "" {
		idWskaznik = &idStrony
	}
	dopasowania := []shared.StudioTextMatch{}
	for numer, wiersz := range podzielNaWiersze(tresc) {
		zakres := dopasuj(wiersz)
		if zakres == nil {
			continue
		}
		kolumna := zakres[0]
		dopasowania = append(dopasowania, shared.StudioTextMatch{
			Line: numer + 1, Column: &kolumna, Text: wiersz, VersionId: idWskaznik,
		})
	}
	return dopasowania, nil
}
