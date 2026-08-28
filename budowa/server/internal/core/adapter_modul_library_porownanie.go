// Moduł Library — porównanie treści: library.diff.compare, dwóch zasobów albo dwóch wersji jednego zasobu, po tekście dla dokumentów binarnych.
package core

import (
	"bytes"
	"context"
	"os"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// granicaWierszyPorownania chroni odpowiedź przed rozrostem: dokument o stu
// tysiącach wierszy dałby różnicę, której żadne okno nie pokaże.
const granicaWierszyPorownania = 5000

// granicaWydobyciaTekstu — wydobycie warstwy tekstowej idzie procesem
// zewnętrznym i musi mieć granicę czasu.
const granicaWydobyciaTekstu = 60 * time.Second

// narzedzieTekstuBiblioteki opisuje binarium arsenału wydobywające tekst z dokumentu PDF programem pdftotext pakietu poppler-utils.
var narzedzieTekstuBiblioteki = zewnetrzne.Narzedzie{
	Nazwa: "Poppler (pdftotext)", Program: "pdftotext", Pakiet: "poppler-utils",
}

// PorownajTresci obsługuje library.diff.compare, porównując treść dwóch zasobów albo dwóch wersji jednego zasobu.
func (a *adapterBiblioteki) PorownajTresci(ctx context.Context,
	z shared.LibraryDiffCompareRequest) (shared.LibraryDiffCompareResponse, error) {

	lewy, err := a.plik(ctx, z.LeftFileId)
	if err != nil {
		return shared.LibraryDiffCompareResponse{}, err
	}

	var prawy dane.PlikBiblioteki
	if z.RightFileId != nil && strings.TrimSpace(*z.RightFileId) != "" {
		prawy, err = a.plik(ctx, *z.RightFileId)
		if err != nil {
			return shared.LibraryDiffCompareResponse{}, err
		}
	} else {
		prawy = lewy
	}

	lewaTresc, lewaNazwa, lewaBinarna, err := a.stronaPorownania(ctx, lewy, z.LeftVersionId)
	if err != nil {
		return shared.LibraryDiffCompareResponse{}, err
	}
	prawaTresc, prawaNazwa, prawaBinarna, err := a.stronaPorownania(ctx, prawy, z.RightVersionId)
	if err != nil {
		return shared.LibraryDiffCompareResponse{}, err
	}
	if lewy.Kod == prawy.Kod && lewaNazwa == prawaNazwa {
		return shared.LibraryDiffCompareResponse{}, bladWskazaniaBiblioteki(
			"porównanie wskazuje dwa razy tę samą treść — wskaż drugi zasób albo drugą wersję")
	}

	otoczenie := 0
	if z.ContextLines != nil && *z.ContextLines > 0 {
		otoczenie = *z.ContextLines
	}
	roznice := roznicaWierszyBiblioteki(podzielNaWierszeBiblioteki(lewaTresc), podzielNaWierszeBiblioteki(prawaTresc), otoczenie)

	a.odnotuj(ctx, shared.LibraryAuditActionAccess, wskazanieBiblioteki(lewy.Kod),
		"porównanie z "+prawaNazwa)
	return shared.LibraryDiffCompareResponse{Diff: shared.LibraryDiff{
		LeftLabel: lewaNazwa, RightLabel: prawaNazwa,
		Identical: len(roznice) == 0, ComparedAsText: lewaBinarna || prawaBinarna,
		Hunks: roznice,
	}}, nil
}

// stronaPorownania oddaje treść jednej strony porównania wraz z jej nazwą
// i informacją, czy treść wydobyto z dokumentu binarnego.
func (a *adapterBiblioteki) stronaPorownania(ctx context.Context, zasob dane.PlikBiblioteki,
	kodWersji *string) (string, string, bool, error) {

	odwolanie := zasob.TrescOdwolanie
	nazwa := zasob.Nazwa
	if kodWersji != nil && strings.TrimSpace(*kodWersji) != "" {
		wersje, err := a.repozytorium.Wersje(ctx, zasob.ID)
		if err != nil {
			return "", "", false, bladBiblioteki(err)
		}
		znaleziona := false
		for _, wersja := range wersje {
			if wersja.Kod != strings.TrimSpace(*kodWersji) {
				continue
			}
			odwolanie, znaleziona = wersja.TrescOdwolanie, true
			nazwa = zasob.Nazwa + " — wersja " + wersja.Kod
			break
		}
		if !znaleziona {
			return "", "", false, bladWskazaniaBiblioteki(
				"wersja " + *kodWersji + " nie należy do zasobu " + zasob.Kod)
		}
	}
	if odwolanie == nil || *odwolanie == "" {
		return "", "", false, bladBrakuTresciBiblioteki(zasob.Kod)
	}
	bajty, err := os.ReadFile(*odwolanie)
	if err != nil {
		return "", "", false, bladBrakuTresciBiblioteki(zasob.Kod)
	}
	tresc, wydobyta := a.trescPorownywalnaBiblioteki(bajty, *odwolanie)
	return tresc, nazwa, wydobyta, nil
}

// trescPorownywalnaBiblioteki oddaje tekst do porównania i mówi, czy trzeba było go wydobyć z dokumentu binarnego programem pdftotext.
func (a *adapterBiblioteki) trescPorownywalnaBiblioteki(bajty []byte, sciezka string) (string, bool) {
	if bytes.HasPrefix(bajty, []byte("%PDF-")) {
		return a.tekstDokumentu(sciezka), true
	}
	// Treść z bajtem zerowym nie jest tekstem; odpowiedź mówi, że wydobyć się nie dało.
	if bytes.IndexByte(bajty, 0) >= 0 {
		return "", true
	}
	return string(bajty), false
}

// tekstDokumentu wyciąga warstwę tekstową dokumentu PDF programem pdftotext, jedyną dozwoloną drogą wołania programów zewnętrznych.
func (a *adapterBiblioteki) tekstDokumentu(sciezka string) string {
	if a.uruchamiacz == nil {
		return ""
	}
	ctx, przerwij := context.WithTimeout(context.Background(), granicaWydobyciaTekstu)
	defer przerwij()

	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}
	// Wynik idzie na wyjście standardowe („-" jako plik docelowy), więc plik
	// pośredni nie powstaje.
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzieTekstuBiblioteki, []string{"-layout", "-enc", "UTF-8", sciezka, "-"},
		"", granicaWydobyciaTekstu)
	if err != nil {
		return ""
	}
	return string(wynik.Wyjscie)
}

// podzielNaWierszeBiblioteki rozbija treść na wiersze z granicą liczby, zapobiegając nadmiernemu rozrostowi porównania.
func podzielNaWierszeBiblioteki(tresc string) []string {
	wiersze := strings.Split(strings.ReplaceAll(tresc, "\r\n", "\n"), "\n")
	if len(wiersze) > granicaWierszyPorownania {
		wiersze = wiersze[:granicaWierszyPorownania]
	}
	return wiersze
}

// roznicaWierszyBiblioteki składa różnice metodą najdłuższego wspólnego podciągu wierszy, tym samym sposobem co porównanie w module Studio.
func roznicaWierszyBiblioteki(lewe, prawe []string, otoczenie int) []shared.LibraryDiffHunk {
	dlugosci := make([][]int, len(lewe)+1)
	for indeks := range dlugosci {
		dlugosci[indeks] = make([]int, len(prawe)+1)
	}
	for i := len(lewe) - 1; i >= 0; i-- {
		for j := len(prawe) - 1; j >= 0; j-- {
			if lewe[i] == prawe[j] {
				dlugosci[i][j] = dlugosci[i+1][j+1] + 1
				continue
			}
			if dlugosci[i+1][j] >= dlugosci[i][j+1] {
				dlugosci[i][j] = dlugosci[i+1][j]
			} else {
				dlugosci[i][j] = dlugosci[i][j+1]
			}
		}
	}

	roznice := []shared.LibraryDiffHunk{}
	i, j := 0, 0
	for i < len(lewe) && j < len(prawe) {
		if lewe[i] == prawe[j] {
			i, j = i+1, j+1
			continue
		}
		poczatekLewy, poczatekPrawy := i, j
		usuniete, dodane := []string{}, []string{}
		for i < len(lewe) && j < len(prawe) && lewe[i] != prawe[j] {
			if dlugosci[i+1][j] >= dlugosci[i][j+1] {
				usuniete = append(usuniete, lewe[i])
				i++
			} else {
				dodane = append(dodane, prawe[j])
				j++
			}
		}
		roznice = append(roznice, zlozRoznicaBiblioteki(poczatekLewy, i, poczatekPrawy, j,
			usuniete, dodane, lewe, prawe, otoczenie))
	}
	if i < len(lewe) {
		roznice = append(roznice, shared.LibraryDiffHunk{
			Kind: shared.LibraryDiffKindRemoved, LeftFrom: i + 1, LeftTo: len(lewe),
			RightFrom: j, RightTo: j, Text: strings.Join(lewe[i:], "\n"),
		})
	}
	if j < len(prawe) {
		roznice = append(roznice, shared.LibraryDiffHunk{
			Kind: shared.LibraryDiffKindAdded, LeftFrom: i, LeftTo: i,
			RightFrom: j + 1, RightTo: len(prawe), Text: strings.Join(prawe[j:], "\n"),
		})
	}
	return roznice
}

// zlozRoznicaBiblioteki składa jedną różnicę wraz z wierszami otoczenia, ułatwiającymi odczytanie zmiany w kontekście.
func zlozRoznicaBiblioteki(poczatekLewy, koniecLewy, poczatekPrawy, koniecPrawy int,
	usuniete, dodane, lewe, prawe []string, otoczenie int) shared.LibraryDiffHunk {

	rodzaj := shared.LibraryDiffKind(shared.LibraryDiffKindChanged)
	tresc := strings.Join(dodane, "\n")
	switch {
	case len(dodane) == 0:
		rodzaj, tresc = shared.LibraryDiffKindRemoved, strings.Join(usuniete, "\n")
	case len(usuniete) == 0:
		rodzaj = shared.LibraryDiffKindAdded
	default:
		tresc = strings.Join(usuniete, "\n") + "\n→\n" + strings.Join(dodane, "\n")
	}
	if otoczenie > 0 {
		przed := poczatekLewy - otoczenie
		if przed < 0 {
			przed = 0
		}
		po := koniecPrawy + otoczenie
		if po > len(prawe) {
			po = len(prawe)
		}
		otoczka := strings.Join(lewe[przed:poczatekLewy], "\n")
		if otoczka != "" {
			tresc = otoczka + "\n" + tresc
		}
		ogon := strings.Join(prawe[koniecPrawy:po], "\n")
		if ogon != "" {
			tresc = tresc + "\n" + ogon
		}
	}
	return shared.LibraryDiffHunk{
		Kind: rodzaj, LeftFrom: poczatekLewy + 1, LeftTo: koniecLewy,
		RightFrom: poczatekPrawy + 1, RightTo: koniecPrawy, Text: tresc,
	}
}
