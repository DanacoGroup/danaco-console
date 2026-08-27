// Rachunek postaci dokumentu Studio: wczytanie drzewa, przeliczenie zakresów
// w runach, rozcięcie fragmentów, przesunięcie zakotwiczeń, bilans blokad
// i zmiana śledzona, wraz z odczytem i zapisem postaci i treści fragmentu.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów postaci dokumentu. Byt nadany przez rdzeń
// wychodzi kontraktem pod własnym identyfikatorem, nie pod kluczem wiersza.
const (
	przedrostekBlokuPostaci   = "studio-blok-"
	przedrostekSekcjiPostaci  = "studio-sek-"
	przedrostekTabeliPostaci  = "studio-tab-"
	przedrostekObiektuPostaci = "studio-obi-"
	przedrostekListyPostaci   = "studio-lis-"
	przedrostekAparatuPostaci = "studio-apa-"
	przedrostekPolaPostaci    = "studio-pol-"
	przedrostekZmianyPostaci  = "studio-zm-"
	przedrostekMalarzaPostaci = "studio-mal-"
)

// Rodzaje bloków drzewa. Kontrakt trzyma rodzaj bloku napisem, bo wykaz jest
// otwarty — bloków dokłada się wraz z rodzajami obiektów.
const (
	blokPostaciAkapit   = "akapit"
	blokPostaciNaglowek = "naglowek"
	blokPostaciTabela   = "tabela"
	blokPostaciObiekt   = "obiekt"
	blokPostaciPodzial  = "podzial"
)

// stanPostaci to postać dokumentu wczytana do pracy: wiersz dokumentu, drzewo
// uzgodnione z treścią oraz treść sprzed czynności.
type stanPostaci struct {
	dokument dane.DokumentStudia
	forma    shared.StudioDocumentForm
	// tekstPrzed niesie treść sprzed czynności dla pola `before` zmiany śledzonej.
	tekstPrzed string
	// drzewoPrzed niesie postać sprzed czynności dla cofnięcia zmiany POSTACI, nie treści, w dzienniku.
	drzewoPrzed string
	// czynnosc to wpis dziennika: czynność cofa się nim pojedynczo, a odpowiedź niesie go jako `actionId`.
	czynnosc *string
	// opisCzynnosci trafia do dziennika jako zdanie dla Operatora.
	opisCzynnosci string
}

// ── Wczytanie i zapis ───────────────────────────────────────────────────────

// postacWczytaj wczytuje dokument wraz z postacią. Dokument bez zapisanej
// postaci dostaje ją założoną z treści wraz z arkuszem stylów fabrycznym —
// pierwsza czynność na postaci jest normalną drogą, nie usterką.
func (a *adapterStudia) postacWczytaj(ctx context.Context, kod string) (*stanPostaci, error) {
	if strings.TrimSpace(kod) == "" {
		return nil, bladWskazaniaStudio("czynność na postaci bez wskazania dokumentu")
	}
	dokument, err := a.repozytorium.Dokument(ctx, kod)
	if err != nil {
		return nil, bladNieznanegoDokumentu(kod, err)
	}

	stan := &stanPostaci{dokument: dokument}
	stan.forma.DocumentId = dokument.Kod
	stan.tekstPrzed = wartoscTekstu(dokument.Tresc)

	wiersz, err := a.repozytorium.PostacDokumentu(ctx, dokument.ID)
	switch {
	case err == nil:
		if bladOdczytu := json.Unmarshal([]byte(wiersz.PostacJSON), &stan.forma); bladOdczytu != nil {
			return nil, postacBladZaplecza("drzewo postaci dokumentu " + kod +
				" jest nieczytelne: " + bladOdczytu.Error())
		}
		stan.forma.DocumentId = dokument.Kod
		wersja := wiersz.WersjaPostaci
		stan.forma.Revision = &wersja
		if wiersz.NastawyStronyJSON != nil && *wiersz.NastawyStronyJSON != "" {
			var nastawy shared.StudioPageSetup
			if json.Unmarshal([]byte(*wiersz.NastawyStronyJSON), &nastawy) == nil {
				stan.forma.PageSetup = &nastawy
			}
		}
	case errors.Is(err, dane.ErrBrakWiersza):
		stan.forma.PageSetup = postacNastawyDomyslne()
	default:
		return nil, bladStudio(err)
	}

	style, err := a.repozytorium.StyleNazwane(ctx, dokument.ID, "")
	if err != nil {
		return nil, bladStudio(err)
	}
	if len(style) == 0 {
		if err := a.postacZasiejStyleFabryczne(ctx, dokument.ID); err != nil {
			return nil, err
		}
		if style, err = a.repozytorium.StyleNazwane(ctx, dokument.ID, ""); err != nil {
			return nil, bladStudio(err)
		}
	}
	stan.forma.Styles = postacZlozStyle(style)

	sekcje, err := a.repozytorium.Sekcje(ctx, dokument.ID)
	if err != nil {
		return nil, bladStudio(err)
	}
	stan.forma.Sections = postacZlozSekcje(sekcje)

	if err := a.postacWczytajWiersze(ctx, stan); err != nil {
		return nil, err
	}

	postacUzgodnijZTrescia(&stan.forma, stan.tekstPrzed)
	postacPrzeliczZakresy(&stan.forma)
	postacPrzeliczUzycieStylow(&stan.forma)

	// Drzewo sprzed czynności zapisuje się tutaj, zanim je zmieni czynność, dla stanu przed w dzienniku.
	if zapis, err := dziennikZapisDrzewa(stan.forma); err == nil && zapis != nil {
		stan.drzewoPrzed = *zapis
	}
	return stan, nil
}

// postacZapisz utrwala drzewo, nastawy strony i treść wyprowadzoną z drzewa.
// Treść idzie razem z postacią jednym przebiegiem — inaczej zapis postaci
// zostawiłby dokument, którego treść mówi co innego niż jego bloki.
func (a *adapterStudia) postacZapisz(ctx context.Context, stan *stanPostaci) error {
	postacPrzeliczZakresy(&stan.forma)
	postacPrzeliczUzycieStylow(&stan.forma)

	// Styl nazwany i sekcja mają własne wiersze będące prawdą; z drzewa idą wycięte przy zapisie.
	doZapisu := stan.forma
	doZapisu.Styles = nil
	doZapisu.Sections = nil
	doZapisu.Objects = nil
	doZapisu.Apparatus = nil
	doZapisu.Fields = nil
	doZapisu.Locks = nil
	doZapisu.Revision = nil
	doZapisu.UpdatedAt = nil
	drzewo, err := json.Marshal(doZapisu)
	if err != nil {
		return postacBladZaplecza("drzewa postaci nie da się zapisać: " + err.Error())
	}

	wiersz := dane.PostacDokumentuStudia{DokumentID: stan.dokument.ID, PostacJSON: string(drzewo)}
	if stan.forma.PageSetup != nil {
		nastawy, err := json.Marshal(stan.forma.PageSetup)
		if err != nil {
			return postacBladZaplecza("nastaw strony nie da się zapisać: " + err.Error())
		}
		zapis := string(nastawy)
		wiersz.NastawyStronyJSON = &zapis
	}
	zapisana, err := a.repozytorium.ZapiszPostacDokumentu(ctx, wiersz)
	if err != nil {
		return bladStudio(err)
	}
	wersja := zapisana.WersjaPostaci
	stan.forma.Revision = &wersja
	chwila := chwilaBazy(zapisana.Zaktualizowano)
	stan.forma.UpdatedAt = &chwila

	tresc := postacTekstFormy(&stan.forma)
	stan.dokument.Tresc = &tresc
	dokument, err := a.repozytorium.ZapiszDokument(ctx, stan.dokument)
	if err != nil {
		return bladStudio(err)
	}
	stan.dokument = dokument
	return nil
}

// postacZasiejStyleFabryczne zakłada arkusz stylów dokumentu: nagłówki sześciu
// poziomów, tekst zasadniczy, cytat, podpis i przypis. Dokument bez arkusza nie
// miałby czego stosować, a wykaz fabryczny jest tym, co Operator zna z pakietu
// biurowego.
func (a *adapterStudia) postacZasiejStyleFabryczne(ctx context.Context, dokumentID int64) error {
	for _, styl := range postacStyleFabryczne() {
		wiersz, err := postacStylDoWiersza(dokumentID, styl)
		if err != nil {
			return err
		}
		if _, err := a.repozytorium.ZapiszStylNazwany(ctx, wiersz); err != nil {
			return bladStudio(err)
		}
	}
	return nil
}

// postacWczytajWiersze dokłada do postaci to, co warstwa danych trzyma
// wierszami, a nie w drzewie: obiekty, aparat, pola i blokady. Po tych bytach
// się PYTA („które przypisy są nieświeże", „co stoi pod blokadą"), więc mają
// wiersze i to one są prawdą.
func (a *adapterStudia) postacWczytajWiersze(ctx context.Context, stan *stanPostaci) error {
	obiekty, err := a.repozytorium.ObiektyDokumentu(ctx, stan.dokument.ID, "")
	if err != nil {
		return bladStudio(err)
	}
	stan.forma.Objects = postacZlozObiekty(obiekty)

	aparat, err := a.repozytorium.ElementyAparatu(ctx, stan.dokument.ID, "")
	if err != nil {
		return bladStudio(err)
	}
	stan.forma.Apparatus = postacZlozAparat(aparat)

	pola, err := a.repozytorium.PolaDokumentu(ctx, stan.dokument.ID, "")
	if err != nil {
		return bladStudio(err)
	}
	stan.forma.Fields = postacZlozPola(pola)

	// Blokady fragmentów stoją w tabeli obszaru kontroli pracy; postać ich nie zakłada, tylko pokazuje.
	kontrola, err := a.kontrolaSkladnica()
	if err != nil {
		return err
	}
	blokady, err := kontrola.BlokadyFragmentow(ctx, stan.dokument.ID)
	if err != nil {
		return bladStudio(err)
	}
	stan.forma.Locks = postacZlozBlokady(stan.dokument.Kod, blokady)
	return nil
}

// ── Uzgodnienie drzewa z treścią ────────────────────────────────────────────

// postacUzgodnijZTrescia prowadzi bloki tekstowe za treścią dokumentu. Bloki
// nietekstowe — tabela, obiekt, podział — nie zajmują ani jednego znaku treści:
// ich miejsce niesie kolejność bloków, nie zakres.
func postacUzgodnijZTrescia(forma *shared.StudioDocumentForm, tresc string) {
	if postacMaBlokiTekstowe(forma) && postacTekstFormy(forma) == tresc {
		return
	}
	linie := strings.Split(tresc, "\n")

	nowe := make([]shared.StudioDocumentBlock, 0, len(forma.Blocks)+len(linie))
	nastepna := 0
	for _, blok := range forma.Blocks {
		if !postacBlokNiesieTekst(blok) {
			nowe = append(nowe, blok)
			continue
		}
		if nastepna >= len(linie) {
			// Blok bez wiersza treści przestaje istnieć, inaczej dałby akapit widmo, którego Operator nie widzi.
			continue
		}
		blok.Runs = postacRunyZWiersza(linie[nastepna], blok.Runs)
		nowe = append(nowe, blok)
		nastepna++
	}
	for ; nastepna < len(linie); nastepna++ {
		nowe = append(nowe, shared.StudioDocumentBlock{
			Id:   nowyIdentyfikator(przedrostekBlokuPostaci),
			Kind: blokPostaciAkapit,
			Runs: postacRunyZWiersza(linie[nastepna], nil),
		})
	}
	forma.Blocks = nowe
}

// postacRunyZWiersza zakłada fragmenty bloku z wiersza treści, przejmując
// postać znaku fragmentu pierwszego — zmiana treści akapitu nie ma prawa
// zdejmować jego kroju.
func postacRunyZWiersza(wiersz string,
	zastane []shared.StudioDocumentRun) []shared.StudioDocumentRun {

	var postac *shared.StudioCharacterFormat
	var autor *shared.StudioAuthor
	if len(zastane) > 0 {
		postac = postacKopiaZnaku(zastane[0].Format)
		autor = zastane[0].AuthoredBy
	}
	return []shared.StudioDocumentRun{{Text: wiersz, Format: postac, AuthoredBy: autor}}
}

// postacBlokNiesieTekst mówi, czy blok wchodzi do liniowego strumienia treści:
// tabela, obiekt i podział nie niosą w nim ani jednego znaku.
func postacBlokNiesieTekst(blok shared.StudioDocumentBlock) bool {
	switch blok.Kind {
	case blokPostaciTabela, blokPostaciObiekt, blokPostaciPodzial:
		return false
	default:
		return true
	}
}

func postacMaBlokiTekstowe(forma *shared.StudioDocumentForm) bool {
	for _, blok := range forma.Blocks {
		if postacBlokNiesieTekst(blok) {
			return true
		}
	}
	return false
}

// postacTekstFormy składa treść dokumentu z bloków tekstowych, pomijając
// tabele, obiekty i podziały, które strumienia treści nie niosą.
func postacTekstFormy(forma *shared.StudioDocumentForm) string {
	czesci := make([]string, 0, len(forma.Blocks))
	for _, blok := range forma.Blocks {
		if !postacBlokNiesieTekst(blok) {
			continue
		}
		czesci = append(czesci, postacTekstBloku(blok))
	}
	return strings.Join(czesci, "\n")
}

func postacTekstBloku(blok shared.StudioDocumentBlock) string {
	var budowa strings.Builder
	for _, run := range blok.Runs {
		budowa.WriteString(run.Text)
	}
	return budowa.String()
}

// postacPrzeliczZakresy nadaje blokom i fragmentom zakresy w znakach.
// Rachunek jest jeden i stoi tutaj: dwa liczenia tego samego zakresu rozjadą
// się przy pierwszej poprawce.
func postacPrzeliczZakresy(forma *shared.StudioDocumentForm) {
	polozenie := 0
	pierwszyTekstowy := true
	for i := range forma.Blocks {
		blok := &forma.Blocks[i]
		if blok.Id == "" {
			blok.Id = nowyIdentyfikator(przedrostekBlokuPostaci)
		}
		if !postacBlokNiesieTekst(*blok) {
			poczatek, koniec := polozenie, polozenie
			blok.RangeStart, blok.RangeEnd = &poczatek, &koniec
			continue
		}
		if !pierwszyTekstowy {
			polozenie++ // znak podziału wiersza między akapitami
		}
		pierwszyTekstowy = false
		poczatek := polozenie
		for j := range blok.Runs {
			run := &blok.Runs[j]
			start := polozenie
			polozenie += len([]rune(run.Text))
			koniec := polozenie
			run.RangeStart, run.RangeEnd = &start, &koniec
		}
		koniec := polozenie
		blok.RangeStart, blok.RangeEnd = &poczatek, &koniec
	}
}

// postacDlugosc oddaje długość treści dokumentu w znakach — miarę, którą
// posługuje się każdy zakres czynności na postaci.
func postacDlugosc(forma *shared.StudioDocumentForm) int {
	return len([]rune(postacTekstFormy(forma)))
}

// ── Zakres zaznaczenia ──────────────────────────────────────────────────────

// postacZakres ustala zakres czynności. Brak wskazania znaczy CAŁY dokument —
// tak mówi kontrakt każdej z tych komend. Zakres odwrócony prostuje się, bo
// zaznaczenie ciągnięte od prawej do lewej jest zwykłym ruchem myszką, nie
// błędem Operatora.
func postacZakres(od, do *int, dlugosc int) (int, int) {
	poczatek, koniec := 0, dlugosc
	if od != nil {
		poczatek = *od
	}
	if do != nil {
		koniec = *do
	}
	if poczatek > koniec {
		poczatek, koniec = koniec, poczatek
	}
	if poczatek < 0 {
		poczatek = 0
	}
	if koniec > dlugosc {
		koniec = dlugosc
	}
	if poczatek > dlugosc {
		poczatek = dlugosc
	}
	return poczatek, koniec
}

// postacRozetnij rozcina fragmenty na granicach zakresu, żeby postać dała się
// nałożyć DOKŁADNIE na zaznaczenie, a nie na cały fragment, który zaznaczenie
// przecina.
func postacRozetnij(forma *shared.StudioDocumentForm, granice ...int) {
	postacPrzeliczZakresy(forma)
	for i := range forma.Blocks {
		blok := &forma.Blocks[i]
		if !postacBlokNiesieTekst(*blok) {
			continue
		}
		nowe := make([]shared.StudioDocumentRun, 0, len(blok.Runs)+len(granice))
		for _, run := range blok.Runs {
			znaki := []rune(run.Text)
			start := 0
			if run.RangeStart != nil {
				start = *run.RangeStart
			}
			koniec := start + len(znaki)
			ciecia := []int{start, koniec}
			for _, granica := range granice {
				if granica > start && granica < koniec {
					ciecia = append(ciecia, granica)
				}
			}
			sort.Ints(ciecia)
			for j := 0; j+1 < len(ciecia); j++ {
				if ciecia[j] == ciecia[j+1] {
					continue
				}
				czesc := run
				czesc.Text = string(znaki[ciecia[j]-start : ciecia[j+1]-start])
				// Postać znaku idzie kopią, nie wskaźnikiem dzielonym — inaczej pogrubienie połowy pogrubiłoby całość.
				czesc.Format = postacKopiaZnaku(run.Format)
				nowe = append(nowe, czesc)
			}
		}
		if len(nowe) == 0 {
			nowe = append(nowe, shared.StudioDocumentRun{Text: ""})
		}
		blok.Runs = nowe
	}
	postacPrzeliczZakresy(forma)
}

// postacFragmentyZakresu oddaje wskazania (blok, fragment) leżące w zakresie.
// Zakres pusty oddaje fragment stojący w punkcie wstawienia — inaczej
// pogrubienie przy pustym zaznaczeniu byłoby ciszą.
func postacFragmentyZakresu(forma *shared.StudioDocumentForm, od, do int) [][2]int {
	wskazania := make([][2]int, 0, 8)
	for i := range forma.Blocks {
		blok := &forma.Blocks[i]
		if !postacBlokNiesieTekst(*blok) {
			continue
		}
		for j := range blok.Runs {
			run := &blok.Runs[j]
			start, koniec := 0, 0
			if run.RangeStart != nil {
				start = *run.RangeStart
			}
			if run.RangeEnd != nil {
				koniec = *run.RangeEnd
			}
			if od == do {
				if start <= od && od <= koniec {
					return append(wskazania, [2]int{i, j})
				}
				continue
			}
			if start >= od && koniec <= do && koniec > start {
				wskazania = append(wskazania, [2]int{i, j})
			}
		}
	}
	return wskazania
}

// postacBlokiZakresu oddaje wskazania bloków, które zakres obejmuje choćby
// częściowo. Postać akapitu jest cechą CAŁEGO akapitu, więc zaznaczenie jednego
// słowa zmienia wyrównanie całego akapitu — tak samo jak w pakiecie biurowym.
func postacBlokiZakresu(forma *shared.StudioDocumentForm, od, do int) []int {
	wskazania := make([]int, 0, 4)
	for i := range forma.Blocks {
		blok := &forma.Blocks[i]
		if !postacBlokNiesieTekst(*blok) {
			continue
		}
		start, koniec := 0, 0
		if blok.RangeStart != nil {
			start = *blok.RangeStart
		}
		if blok.RangeEnd != nil {
			koniec = *blok.RangeEnd
		}
		if start <= do && koniec >= od {
			wskazania = append(wskazania, i)
		}
	}
	return wskazania
}

// ── Elementy treści: rozkład, zamiana, złożenie ─────────────────────────────

// elementPostaci to jedno miejsce w zakresie dokumentu: albo znak wraz ze swoją
// postacią, albo granica akapitu wraz z postacią akapitu, który za nią stoi.
type elementPostaci struct {
	znak    rune
	postac  *shared.StudioCharacterFormat
	autor   *shared.StudioAuthor
	granica bool
	// Cechy akapitu ZA granicą — przy zlaniu akapitów odchodzą razem z granicą.
	akapit *shared.StudioParagraphFormat
	rodzaj string
	sekcja *string
	kod    string
}

// postacNaElementy rozkłada bloki tekstowe na elementy i oddaje cechy akapitu
// pierwszego, których żadna granica nie niesie.
func postacNaElementy(forma *shared.StudioDocumentForm) ([]elementPostaci, shared.StudioDocumentBlock) {
	elementy := make([]elementPostaci, 0, 256)
	var pierwszy shared.StudioDocumentBlock
	mamPierwszy := false
	for _, blok := range forma.Blocks {
		if !postacBlokNiesieTekst(blok) {
			continue
		}
		if !mamPierwszy {
			pierwszy = blok
			pierwszy.Runs = nil
			mamPierwszy = true
		} else {
			elementy = append(elementy, elementPostaci{
				granica: true, akapit: blok.Paragraph, rodzaj: blok.Kind,
				sekcja: blok.SectionId, kod: blok.Id,
			})
		}
		for _, run := range blok.Runs {
			for _, znak := range run.Text {
				elementy = append(elementy, elementPostaci{
					znak: znak, postac: run.Format, autor: run.AuthoredBy,
				})
			}
		}
	}
	if !mamPierwszy {
		pierwszy = shared.StudioDocumentBlock{
			Id:   nowyIdentyfikator(przedrostekBlokuPostaci),
			Kind: blokPostaciAkapit,
		}
	}
	return elementy, pierwszy
}

// postacZElementow składa bloki tekstowe z elementów, zbierając ciągi znaków
// o tej samej postaci w jeden fragment — inaczej drzewo rosłoby do jednego
// fragmentu na literę.
func postacZElementow(elementy []elementPostaci,
	pierwszy shared.StudioDocumentBlock) []shared.StudioDocumentBlock {

	bloki := []shared.StudioDocumentBlock{pierwszy}
	for _, element := range elementy {
		if element.granica {
			kod := element.kod
			if kod == "" {
				kod = nowyIdentyfikator(przedrostekBlokuPostaci)
			}
			rodzaj := element.rodzaj
			if rodzaj == "" {
				rodzaj = blokPostaciAkapit
			}
			bloki = append(bloki, shared.StudioDocumentBlock{
				Id: kod, Kind: rodzaj, SectionId: element.sekcja, Paragraph: element.akapit,
			})
			continue
		}
		blok := &bloki[len(bloki)-1]
		if len(blok.Runs) > 0 &&
			postacTaSamaPostacZnaku(blok.Runs[len(blok.Runs)-1].Format, element.postac) &&
			postacTenSamAutor(blok.Runs[len(blok.Runs)-1].AuthoredBy, element.autor) {

			blok.Runs[len(blok.Runs)-1].Text += string(element.znak)
			continue
		}
		blok.Runs = append(blok.Runs, shared.StudioDocumentRun{
			Text: string(element.znak), Format: element.postac, AuthoredBy: element.autor,
		})
	}
	for i := range bloki {
		if bloki[i].Runs == nil {
			bloki[i].Runs = []shared.StudioDocumentRun{{Text: ""}}
		}
	}
	return bloki
}

// postacZamienTresc zamienia zakres treści na nowe brzmienie i oddaje różnicę
// długości. Bloki nietekstowe zostają na swoich miejscach — tabela nie znika,
// bo Operator przepisał akapit nad nią.
func postacZamienTresc(forma *shared.StudioDocumentForm, od, do int, tekst string,
	znak *shared.StudioCharacterFormat, autor *shared.StudioAuthor) int {

	elementy, pierwszy := postacNaElementy(forma)
	if od > len(elementy) {
		od = len(elementy)
	}
	if do > len(elementy) {
		do = len(elementy)
	}
	nowe := make([]elementPostaci, 0, len([]rune(tekst)))
	for _, litera := range tekst {
		if litera == '\n' {
			nowe = append(nowe, elementPostaci{
				granica: true, rodzaj: blokPostaciAkapit,
				kod: nowyIdentyfikator(przedrostekBlokuPostaci),
			})
			continue
		}
		nowe = append(nowe, elementPostaci{znak: litera, postac: postacKopiaZnaku(znak), autor: autor})
	}
	roznica := len(nowe) - (do - od)

	wynik := make([]elementPostaci, 0, len(elementy)+roznica)
	wynik = append(wynik, elementy[:od]...)
	wynik = append(wynik, nowe...)
	wynik = append(wynik, elementy[do:]...)

	tekstowe := postacZElementow(wynik, pierwszy)
	postacWstawBlokiTekstowe(forma, tekstowe)
	postacPrzesunZakotwiczenia(forma, od, do, roznica)
	postacPrzeliczZakresy(forma)
	return roznica
}

// postacWstawBlokiTekstowe podmienia bloki tekstowe drzewa, zostawiając bloki
// nietekstowe na swoich miejscach względem kolejności czytania.
func postacWstawBlokiTekstowe(forma *shared.StudioDocumentForm,
	tekstowe []shared.StudioDocumentBlock) {

	nietekstowe := make([]shared.StudioDocumentBlock, 0, len(forma.Blocks))
	// Blok nietekstowy zapamiętuje, po ilu blokach tekstowych stał, i wraca na to samo miejsce.
	miejsca := make([]int, 0, len(forma.Blocks))
	ile := 0
	for _, blok := range forma.Blocks {
		if postacBlokNiesieTekst(blok) {
			ile++
			continue
		}
		nietekstowe = append(nietekstowe, blok)
		miejsca = append(miejsca, ile)
	}

	nowe := make([]shared.StudioDocumentBlock, 0, len(tekstowe)+len(nietekstowe))
	nastepny := 0
	for i, blok := range tekstowe {
		for nastepny < len(miejsca) && miejsca[nastepny] == i {
			nowe = append(nowe, nietekstowe[nastepny])
			nastepny++
		}
		nowe = append(nowe, blok)
	}
	for ; nastepny < len(nietekstowe); nastepny++ {
		nowe = append(nowe, nietekstowe[nastepny])
	}
	forma.Blocks = nowe
}

// ── Blokady fragmentów ──────────────────────────────────────────────────────

// postacOdcinkiDozwolone dzieli zakres czynności na odcinki wolne od blokad
// i oddaje bilans pominięć. Zasięg `model` zatrzymuje wyłącznie czynność
// autora `model`, zasięg `everyone` — każdą.
func postacOdcinkiDozwolone(forma *shared.StudioDocumentForm, od, do int,
	autor shared.StudioAuthor) ([][2]int, []shared.StudioSkippedItem) {

	blokady := make([]shared.StudioFragmentLock, 0, len(forma.Locks))
	for _, blokada := range forma.Locks {
		if blokada.Scope == shared.StudioLockScopeEveryone || autor == shared.StudioAuthorModel {
			blokady = append(blokady, blokada)
		}
	}
	if len(blokady) == 0 {
		return [][2]int{{od, do}}, nil
	}
	sort.Slice(blokady, func(i, j int) bool { return blokady[i].RangeStart < blokady[j].RangeStart })

	odcinki := make([][2]int, 0, len(blokady)+1)
	pominiete := make([]shared.StudioSkippedItem, 0, len(blokady))
	biezacy := od
	for _, blokada := range blokady {
		if blokada.RangeEnd <= od || blokada.RangeStart >= do {
			continue
		}
		poczatek, koniec := blokada.RangeStart, blokada.RangeEnd
		if poczatek < od {
			poczatek = od
		}
		if koniec > do {
			koniec = do
		}
		if poczatek > biezacy {
			odcinki = append(odcinki, [2]int{biezacy, poczatek})
		}
		nazwa, kod := blokada.Name, blokada.Id
		poczatekPominiecia, koniecPominiecia := poczatek, koniec
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "fragment zablokowany",
			Detail: postacWskaznikTekstu("blokada „" + nazwa +
				"” nie przepuszcza tej czynności"),
			LockId:     &kod,
			LockName:   &nazwa,
			RangeStart: &poczatekPominiecia,
			RangeEnd:   &koniecPominiecia,
		})
		if koniec > biezacy {
			biezacy = koniec
		}
	}
	if biezacy < do {
		odcinki = append(odcinki, [2]int{biezacy, do})
	}
	if len(pominiete) == 0 {
		return [][2]int{{od, do}}, nil
	}
	return odcinki, pominiete
}

// postacBlokadyZakresu oddaje blokady przecinające zakres — okno ma pokazać
// Operatorowi, czego model nie tknie, PRZED próbą, a nie po odmowie.
func postacBlokadyZakresu(forma *shared.StudioDocumentForm, od, do int) []shared.StudioFragmentLock {
	wybrane := make([]shared.StudioFragmentLock, 0, len(forma.Locks))
	for _, blokada := range forma.Locks {
		if blokada.RangeStart < do && blokada.RangeEnd > od {
			wybrane = append(wybrane, blokada)
		}
	}
	return wybrane
}

// postacNazwaBlokad zbiera nazwy blokad, które czynność zatrzymały. Odmowa
// nazwana, nie cicha: Operator ma wiedzieć, KTÓRA blokada go zatrzymała.
func postacNazwaBlokad(pominiete []shared.StudioSkippedItem) string {
	nazwy := make([]string, 0, len(pominiete))
	for _, pozycja := range pominiete {
		if pozycja.LockName != nil {
			nazwy = append(nazwy, "„"+*pozycja.LockName+"”")
		}
	}
	if len(nazwy) == 0 {
		return "blokada bez nazwy"
	}
	return strings.Join(nazwy, ", ")
}

// ── Zmiana śledzona ─────────────────────────────────────────────────────────

// postacAutor rozstrzyga autora czynności. Brak wskazania znaczy Operatora —
// narzędzie modelu podaje autora wprost i to jest jedyna droga, którą w rdzeniu
// staje się zmiana autora `model`.
func postacAutor(wskazanie *shared.StudioAuthor) shared.StudioAuthor {
	if wskazanie != nil && *wskazanie == shared.StudioAuthorModel {
		return shared.StudioAuthorModel
	}
	return shared.StudioAuthorUzytkownik
}

// postacOdlozZmiane rejestruje zmianę śledzoną czynności na postaci. Czynność
// modelu odkłada się zawsze; czynność operatora — wyłącznie gdy śledzenie
// zmian dokumentu jest włączone.
func (a *adapterStudia) postacOdlozZmiane(ctx context.Context, stan *stanPostaci,
	autor shared.StudioAuthor, rodzaj shared.StudioChangeKind, od, do int,
	przed, po *string) (*shared.StudioTrackedChange, error) {

	if autor != shared.StudioAuthorModel {
		czynne, err := a.repozytorium.Sledzenie(ctx, stan.dokument.Kod)
		if err != nil {
			return nil, bladStudio(err)
		}
		if !czynne {
			return nil, nil
		}
	}
	zapisana, err := a.repozytorium.ZapiszZmianeSledzona(ctx, stan.dokument.ID, dane.ZmianaSledzona{
		Kod:         nowyIdentyfikator(przedrostekZmianyPostaci),
		DokumentKod: stan.dokument.Kod,
		Rodzaj:      string(rodzaj),
		Autor:       string(autor),
		ZakresOd:    int64(od),
		ZakresDo:    int64(do),
		TrescPrzed:  przed,
		TrescPo:     po,
		Decyzja:     string(shared.StudioChangeDecisionOczekuje),
	})
	if err != nil {
		return nil, bladStudio(err)
	}
	zmiana := złóżZmianeSledzona(zapisana)
	return &zmiana, nil
}

// postacZakoncz domyka czynność zmieniającą postać: zapisuje drzewo, odkłada
// zmianę śledzoną i domyka bilans — jedna droga wyjścia dla wszystkich
// czynności postaci. Rodzaj zmiany śledzonej i rodzaj czynności dziennika są
// dwoma osobnymi słownikami.
func (a *adapterStudia) postacZakoncz(ctx context.Context, stan *stanPostaci,
	autor shared.StudioAuthor, rodzaj shared.StudioChangeKind,
	czynnosc shared.StudioActionKind, od, do int,
	bilans shared.StudioActionBalance) (shared.StudioDocumentForm, shared.StudioActionBalance,
	*shared.StudioTrackedChange, error) {

	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioDocumentForm{}, shared.StudioActionBalance{}, nil, err
	}
	if bilans.Skipped == nil {
		bilans.Skipped = []shared.StudioSkippedItem{}
	}
	bilans.SkippedCount = len(bilans.Skipped)

	var zmiana *shared.StudioTrackedChange
	if bilans.Applied > 0 {
		przed, po := stan.tekstPrzed, wartoscTekstu(stan.dokument.Tresc)
		wykonawca := postacWykonawca(ctx, autor)
		var err error
		zmiana, err = a.postacOdlozZmiane(ctx, stan, autor, rodzaj, od, do, &przed, &po)
		if err != nil {
			return shared.StudioDocumentForm{}, shared.StudioActionBalance{}, nil, err
		}
		if err := a.postacOdlozCzynnosc(ctx, stan, wykonawca, czynnosc, od, do,
			zmiana, bilans); err != nil {

			return shared.StudioDocumentForm{}, shared.StudioActionBalance{}, nil, err
		}
		// Stempel wykonawcy idzie po odłożeniu czynności, bo niesie także jej kod dla wiązania.
		if zmiana != nil && wykonawca.czyWykonawca() {
			if err := a.zmianyModeluStempluj(ctx, zmiana, wykonawca, stan.czynnosc); err != nil {
				return shared.StudioDocumentForm{}, shared.StudioActionBalance{}, nil, err
			}
		}
	}
	return stan.forma, bilans, zmiana, nil
}

// postacWykonawca rozstrzyga, kto wykonuje czynność na postaci, wraz
// z tożsamością agenta, o ile żądanie ją podało. Podpis wchodzi z kontekstu;
// rodzaj autora rozstrzyga się zasadą „szerszy wygrywa".
func postacWykonawca(ctx context.Context, autor shared.StudioAuthor) kontrolaWykonawca {
	wykonawca, jest := kontrolaWykonawcaZKontekstu(ctx)
	if !jest {
		return kontrolaWykonawca{Rodzaj: autor}
	}
	if autor == shared.StudioAuthorModel {
		wykonawca.Rodzaj = shared.StudioAuthorModel
	}
	return wykonawca
}

// postacOdlozCzynnosc dopisuje czynność do odwracalnego dziennika dokumentu,
// wspólnego dla całego modułu. Stan sprzed czynności idzie całym drzewem
// postaci, nie samą treścią.
func (a *adapterStudia) postacOdlozCzynnosc(ctx context.Context, stan *stanPostaci,
	wykonawca kontrolaWykonawca, rodzaj shared.StudioActionKind, od, do int,
	zmiana *shared.StudioTrackedChange, bilans shared.StudioActionBalance) error {

	opis := stan.opisCzynnosci
	if opis == "" && bilans.Note != nil {
		opis = *bilans.Note
	}
	if opis == "" {
		opis = "czynność na postaci dokumentu " + postacZapisZakresu(od, do)
	}
	poZmianie, err := json.Marshal(stan.forma)
	if err != nil {
		return postacBladZaplecza("stanu po czynności nie da się zapisać w dzienniku: " + err.Error())
	}
	// Przedrostek kodu czynności jest jeden dla całego dziennika modułu, wspólny z kontrolą pracy.
	wpis := dane.CzynnoscDokumentuStudia{
		Kod:         nowyIdentyfikator(przedrostekCzynnosciStudia),
		DokumentKod: stan.dokument.Kod,
		Rodzaj:      string(rodzaj),
		AutorRodzaj: string(wykonawca.Rodzaj),
		// Kod agenta rozróżnia dwóch wykonawców pracujących naraz nad jednym pismem, rodzaj autora nie.
		AutorAgentKod:    wykonawca.AgentKod,
		AutorAgentNazwa:  wykonawca.AgentNazwa,
		AutorAgentWersja: wykonawca.AgentWersja,
		AutorPodagentKod: wykonawca.PodagentKod,
		Opis:             opis,
		ZakresOd:         postacWskaznikLiczby64(int64(od)),
		ZakresDo:         postacWskaznikLiczby64(int64(do)),
		StanPo:           postacWskaznikTekstu(string(poZmianie)),
		// Stan wpisu jest słownikiem tabeli dziennika: `active` albo `reverted`.
		Stan: "active",
	}
	if stan.drzewoPrzed != "" {
		wpis.StanPrzed = postacWskaznikTekstu(stan.drzewoPrzed)
	}
	if zmiana != nil {
		wpis.ZmianaSledzonaK = postacWskaznikTekstu(zmiana.Id)
	}
	kontrola, err := a.kontrolaSkladnica()
	if err != nil {
		return err
	}
	zapisana, err := kontrola.ZapiszCzynnoscDokumentu(ctx, stan.dokument.ID, wpis)
	if err != nil {
		return bladStudio(err)
	}
	stan.czynnosc = postacWskaznikTekstu(zapisana.Kod)
	return nil
}

// ── Przesunięcie zakotwiczeń ────────────────────────────────────────────────

// postacPrzesunZakotwiczenia przesuwa wszystko, co wisi na miejscu w treści,
// po zmianie długości tej treści: obiekty, pola, aparat, sekcje i blokady.
func postacPrzesunZakotwiczenia(forma *shared.StudioDocumentForm, od, do, roznica int) {
	if roznica == 0 {
		return
	}
	przesun := func(miejsce *int) {
		if miejsce == nil {
			return
		}
		switch {
		case *miejsce >= do:
			*miejsce += roznica
		case *miejsce > od:
			*miejsce = od
		}
	}
	for i := range forma.Objects {
		przesun(forma.Objects[i].AnchorOffset)
	}
	for i := range forma.Fields {
		przesun(forma.Fields[i].AnchorOffset)
	}
	for i := range forma.Apparatus {
		przesun(forma.Apparatus[i].AnchorStart)
		przesun(forma.Apparatus[i].AnchorEnd)
		for j := range forma.Apparatus[i].Entries {
			przesun(forma.Apparatus[i].Entries[j].AnchorOffset)
		}
	}
	for i := range forma.Sections {
		poczatek, koniec := forma.Sections[i].RangeStart, forma.Sections[i].RangeEnd
		przesun(&poczatek)
		przesun(&koniec)
		forma.Sections[i].RangeStart, forma.Sections[i].RangeEnd = poczatek, koniec
	}
	for i := range forma.Locks {
		poczatek, koniec := forma.Locks[i].RangeStart, forma.Locks[i].RangeEnd
		przesun(&poczatek)
		przesun(&koniec)
		forma.Locks[i].RangeStart, forma.Locks[i].RangeEnd = poczatek, koniec
	}
}

// ── Czynności podstawowe ────────────────────────────────────────────────────

// PostacDokumentu oddaje pełną postać dokumentu (`studio.document.form.get`),
// z możliwością ograniczenia odpowiedzi do samego zakresu zaznaczenia.
func (a *adapterStudia) PostacDokumentu(ctx context.Context,
	z shared.StudioDocumentFormGetRequest) (shared.StudioDocumentFormGetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDocumentFormGetResponse{}, err
	}
	forma := stan.forma
	if z.IncludeBlocks != nil && !*z.IncludeBlocks {
		forma.Blocks = nil
		return shared.StudioDocumentFormGetResponse{Form: forma}, nil
	}
	if z.RangeStart != nil || z.RangeEnd != nil {
		od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&forma))
		wybrane := make([]shared.StudioDocumentBlock, 0, len(forma.Blocks))
		for _, wskazanie := range postacBlokiZakresu(&forma, od, do) {
			wybrane = append(wybrane, forma.Blocks[wskazanie])
		}
		forma.Blocks = wybrane
	}
	return shared.StudioDocumentFormGetResponse{Form: forma}, nil
}

// ZapiszPostacDokumentu zapisuje dokument wraz z jego postacią
// (`studio.document.form.save`) — to jest droga, którą postać przestaje ginąć
// przy zapisie.
func (a *adapterStudia) ZapiszPostacDokumentu(ctx context.Context,
	z shared.StudioDocumentFormSaveRequest) (shared.StudioDocumentFormSaveResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDocumentFormSaveResponse{}, err
	}
	autor := postacAutor(z.Author)

	var przyslana shared.StudioDocumentForm
	if len(z.Form) > 0 {
		if err := json.Unmarshal(z.Form, &przyslana); err != nil {
			return shared.StudioDocumentFormSaveResponse{}, bladWskazaniaStudio(
				"postać dokumentu w żądaniu jest nieczytelna: " + err.Error())
		}
	}
	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}

	if przyslana.PageSetup != nil {
		stan.forma.PageSetup = przyslana.PageSetup
		bilans.Applied++
	}
	if len(przyslana.Blocks) > 0 {
		stan.forma.Blocks = przyslana.Blocks
		bilans.Applied += len(przyslana.Blocks)
	}
	for _, obszar := range []struct {
		ile   int
		zapis func()
	}{
		{len(przyslana.Tables), func() { stan.forma.Tables = przyslana.Tables }},
		{len(przyslana.Objects), func() { stan.forma.Objects = przyslana.Objects }},
		{len(przyslana.Lists), func() { stan.forma.Lists = przyslana.Lists }},
		{len(przyslana.Apparatus), func() { stan.forma.Apparatus = przyslana.Apparatus }},
		{len(przyslana.Fields), func() { stan.forma.Fields = przyslana.Fields }},
	} {
		if obszar.ile > 0 {
			obszar.zapis()
			bilans.Applied += obszar.ile
		}
	}
	if len(przyslana.Styles) > 0 {
		skladnica, err := a.postacSkladnica()
		if err != nil {
			return shared.StudioDocumentFormSaveResponse{}, err
		}
		for _, styl := range przyslana.Styles {
			if postacStylFabryczny(stan.forma.Styles, styl.Name) {
				bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
					Reason: "styl fabryczny",
					Detail: postacWskaznikTekstu("styl „" + styl.Name +
						"” jest fabryczny; zmiana idzie stylem własnym dziedziczącym po nim"),
				})
				continue
			}
			wiersz, err := postacStylDoWiersza(stan.dokument.ID, styl)
			if err != nil {
				return shared.StudioDocumentFormSaveResponse{}, err
			}
			if _, err := skladnica.ZapiszStylNazwany(ctx, wiersz); err != nil {
				return shared.StudioDocumentFormSaveResponse{}, bladStudio(err)
			}
			bilans.Applied++
		}
		if style, err := skladnica.StyleNazwane(ctx, stan.dokument.ID, ""); err == nil {
			stan.forma.Styles = postacZlozStyle(style)
		}
	}

	if z.Content != nil {
		postacUzgodnijZTrescia(&stan.forma, *z.Content)
		bilans.Applied++
	}
	if z.Title != nil {
		stan.dokument.Tytul = z.Title
		bilans.Applied++
	}
	if bilans.Applied == 0 {
		return shared.StudioDocumentFormSaveResponse{}, bladWskazaniaStudio(
			"zapis postaci bez ani jednej rzeczy do zapisania — żądanie nie niosło " +
				"postaci, treści ani tytułu")
	}
	bilans.SkippedCount = len(bilans.Skipped)
	bilans.Note = postacWskaznikTekstu("postać dokumentu zapisana wraz z treścią")

	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioDocumentFormSaveResponse{}, err
	}

	odpowiedz := shared.StudioDocumentFormSaveResponse{
		Document: a.zlozDokument(stan.dokument),
		Form:     stan.forma,
		Balance:  bilans,
	}
	if z.CreateVersion != nil && *z.CreateVersion {
		dokument, err := a.zalozWersjeDokumentu(ctx, stan.dokument,
			wartoscTekstu(stan.dokument.Tresc), autor, nil)
		if err != nil {
			return shared.StudioDocumentFormSaveResponse{}, err
		}
		stan.dokument = dokument
		odpowiedz.Document = a.zlozDokument(dokument)
		if kod := dokument.WersjaBiezacaKod; kod != nil {
			if wersja, err := a.repozytorium.Wersja(ctx, *kod); err == nil {
				zlozona := a.zlozWersje(wersja)
				odpowiedz.Version = &zlozona
			}
		}
	}
	return odpowiedz, nil
}

// TrescFragmentu oddaje treść zakresu wraz z jego postacią (`studio.text.get`),
// rozcinając fragmenty na granicach zakresu, by postać była dokładna.
func (a *adapterStudia) TrescFragmentu(ctx context.Context,
	z shared.StudioTextGetRequest) (shared.StudioTextGetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTextGetResponse{}, err
	}
	znaki := []rune(postacTekstFormy(&stan.forma))
	od, do := postacZakres(z.RangeStart, z.RangeEnd, len(znaki))

	postacRozetnij(&stan.forma, od, do)
	wskazania := postacFragmentyZakresu(&stan.forma, od, do)
	runy := make([]shared.StudioDocumentRun, 0, len(wskazania))
	for _, wskazanie := range wskazania {
		runy = append(runy, stan.forma.Blocks[wskazanie[0]].Runs[wskazanie[1]])
	}

	odpowiedz := shared.StudioTextGetResponse{
		Text:  string(znaki[od:do]),
		Runs:  runy,
		Locks: postacBlokadyZakresu(&stan.forma, od, do),
	}
	bloki := postacBlokiZakresu(&stan.forma, od, do)
	if len(bloki) > 0 {
		akapit := postacAkapitSkuteczny(&stan.forma, stan.forma.Blocks[bloki[0]])
		odpowiedz.Paragraph = &akapit
	}
	if len(runy) > 0 {
		var pierwszy *shared.StudioDocumentBlock
		if len(bloki) > 0 {
			pierwszy = &stan.forma.Blocks[bloki[0]]
		}
		znak := postacZnakSkutecznyWBloku(&stan.forma, pierwszy, runy[0].Format)
		odpowiedz.Character = &znak
	}
	return odpowiedz, nil
}

// ZmienTresc zamienia treść wskazanego fragmentu (`studio.text.edit`). Zamiana
// idzie po elementach, nie po napisie: postać wokół zmiany zostaje, a zakotwiczenia
// przesuwają się o różnicę długości.
func (a *adapterStudia) ZmienTresc(ctx context.Context,
	z shared.StudioTextEditRequest) (shared.StudioTextEditResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTextEditResponse{}, err
	}
	autor := postacAutor(z.Author)
	od, do := postacZakres(&z.RangeStart, &z.RangeEnd, postacDlugosc(&stan.forma))

	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)
	if len(odcinki) == 0 {
		return shared.StudioTextEditResponse{}, bladWskazaniaStudio(
			"zmiana treści zatrzymana w całości przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}
	// Zamiana idzie na odcinku pierwszym wolnym od blokady; zakres poszatkowany nie ma jednego skutku.
	roboczyOd, roboczyDo := odcinki[0][0], odcinki[0][1]

	var przejeta *shared.StudioCharacterFormat
	if z.KeepFormat == nil || *z.KeepFormat {
		postacRozetnij(&stan.forma, roboczyOd, roboczyDo)
		if wskazania := postacFragmentyZakresu(&stan.forma, roboczyOd, roboczyDo); len(wskazania) > 0 {
			przejeta = postacKopiaZnaku(
				stan.forma.Blocks[wskazania[0][0]].Runs[wskazania[0][1]].Format)
		}
	}
	autorWpisu := autor
	postacZamienTresc(&stan.forma, roboczyOd, roboczyDo, z.Text, przejeta, &autorWpisu)

	rodzaj := shared.StudioChangeKind(shared.StudioChangeKindWstawienie)
	if z.Text == "" {
		rodzaj = shared.StudioChangeKind(shared.StudioChangeKindUsuniecie)
	}
	bilans := shared.StudioActionBalance{Applied: 1, Skipped: pominiete}
	bilans.Note = postacWskaznikTekstu("treść zakresu " + postacZapisZakresu(roboczyOd, roboczyDo) +
		" zamieniona")

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor, rodzaj,
		shared.StudioActionKindTextEdit, roboczyOd, roboczyDo, bilans)
	if err != nil {
		return shared.StudioTextEditResponse{}, err
	}
	return shared.StudioTextEditResponse{Form: forma, Balance: bilansGotowy, Change: zmiana}, nil
}
