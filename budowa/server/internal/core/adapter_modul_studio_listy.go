// Odpowiedzialność pliku: listy dokumentu — wypunktowanie, numeracja, listy
// wielopoziomowe, wznowienie numeracji, własny znak wypunktowania oraz wcięcia
// i odstępy poziomów.
//
// ── Gdzie stoi prawda o liście ───────────────────────────────────────────────
// Definicja listy — jej rodzaj, poziomy, znaki wypunktowania, formaty numeracji
// i wcięcia — stoi w drzewie postaci (`forma.Lists`). Akapit należący do listy
// trzyma sam KOD listy i numer poziomu, a nie kopię jej nastaw. Jest to ta sama
// zasada, na której stoi styl nazwany: zmiana znaku wypunktowania poziomu ma
// przestawić wszystkie punkty tego poziomu jednym ruchem, a nie wymagać
// przejścia po akapitach. Gdyby akapit trzymał kopię, lista o czterdziestu
// punktach wymagałaby czterdziestu poprawek.
//
// ── Dlaczego wcięcie idzie do akapitu, choć poziom je już niesie ─────────────
// Poziom listy niesie wcięcie wzorcowe, ale wcięcie skuteczne akapitu musi być
// widoczne w jego postaci — inaczej linijka nie miałaby czym pokazać znacznika
// wcięcia, a Właściciel wymaga wcięć przestawianych chwytem na linijce. Dlatego
// zastosowanie listy USTAWIA wcięcie akapitu wedle poziomu, a późniejsza zmiana
// wcięcia chwytem jest zmianą akapitu i poziomu nie rusza. Tak samo działa
// pakiet biurowy.
//
// ── Dlaczego wznowienie numeracji zakłada listę nową ─────────────────────────
// Kontrakt niesie punkt startu na liście i na poziomie, ale nie w akapicie —
// i słusznie, bo akapit nie jest miejscem na nastawę listy. Wznowienie numeracji
// od wskazanego miejsca jest więc ROZDZIELENIEM listy: punkty od tego miejsca
// w dół przechodzą do listy nowej o tych samych poziomach i własnym punkcie
// startu. Skutek jest dokładnie ten, którego Operator oczekuje („tu zacznij od
// jedynki"), i jest sprawdzalny w bazie — a nie udawany polem, którego kontrakt
// nie ma.
package core

import (
	"sort"
	"strings"

	"context"

	"danacoconsole/shared"
)

// listyGlebokoscMaksymalna to najgłębszy poziom listy. Dziewięć poziomów niesie
// pakiet biurowy i tyle wystarcza numeracji prawniczej wielopoziomowej
// (1.1.2 i głębiej).
const listyGlebokoscMaksymalna = 9

// listySzerokoscPoziomuMm to odstęp jednego poziomu listy. Jedna czwarta cala
// to wcięcie, które Operator zna z pakietu biurowego.
const listySzerokoscPoziomuMm = 6.35

// listyZnakiPoziomow wylicza znaki wypunktowania kolejnych poziomów. Wykaz jest
// tym, co Operator widzi w pakiecie biurowym, i powtarza się cyklicznie na
// poziomach głębszych.
var listyZnakiPoziomow = []string{"•", "◦", "▪", "‣", "·"}

// listyRodzajZnany sprawdza rodzaj listy wobec kontraktu i oddaje odmowę
// nazywającą wykaz, a nie samo „nie".
func listyRodzajZnany(rodzaj shared.StudioListKind) error {
	for _, wartosc := range shared.WartosciStudioListKind() {
		if wartosc == rodzaj {
			return nil
		}
	}
	nazwy := make([]string, 0, 4)
	for _, wartosc := range shared.WartosciStudioListKind() {
		nazwy = append(nazwy, string(wartosc))
	}
	return bladWskazaniaStudio("rodzaju listy „" + string(rodzaj) +
		"” kontrakt nie zna; rodzaje: " + strings.Join(nazwy, ", "))
}

// listyFormatZnany sprawdza format numeracji wobec kontraktu.
func listyFormatZnany(format shared.StudioListNumberFormat) error {
	for _, wartosc := range shared.WartosciStudioListNumberFormat() {
		if wartosc == format {
			return nil
		}
	}
	nazwy := make([]string, 0, 6)
	for _, wartosc := range shared.WartosciStudioListNumberFormat() {
		nazwy = append(nazwy, string(wartosc))
	}
	return bladWskazaniaStudio("formatu numeracji „" + string(format) +
		"” kontrakt nie zna; formaty: " + strings.Join(nazwy, ", "))
}

// listyZrodloZnane sprawdza źródło znaku wypunktowania wobec kontraktu.
func listyZrodloZnane(zrodlo shared.StudioBulletSource) error {
	for _, wartosc := range shared.WartosciStudioBulletSource() {
		if wartosc == zrodlo {
			return nil
		}
	}
	nazwy := make([]string, 0, 4)
	for _, wartosc := range shared.WartosciStudioBulletSource() {
		nazwy = append(nazwy, string(wartosc))
	}
	return bladWskazaniaStudio("źródła znaku wypunktowania „" + string(zrodlo) +
		"” kontrakt nie zna; źródła: " + strings.Join(nazwy, ", "))
}

// listyDefinicja znajduje definicję listy po kodzie.
func listyDefinicja(forma *shared.StudioDocumentForm, kod string) *shared.StudioListDefinition {
	szukany := strings.TrimSpace(kod)
	if szukany == "" {
		return nil
	}
	for i := range forma.Lists {
		if forma.Lists[i].Id == szukany {
			return &forma.Lists[i]
		}
	}
	return nil
}

// listyPoziomWzorcowy składa poziom listy o nastawach domyślnych — wcięcie
// rosnące z poziomem, znak wypunktowania z wykazu cyklicznego, numeracja arabska
// dla listy numerowanej.
func listyPoziomWzorcowy(rodzaj shared.StudioListKind, poziom int) shared.StudioListLevel {
	wynik := shared.StudioListLevel{
		Level:     poziom,
		IndentMm:  postacWskaznikMiary(listySzerokoscPoziomuMm * float64(poziom)),
		HangingMm: postacWskaznikMiary(listySzerokoscPoziomuMm),
		Align:     postacWskaznikWyrownania(shared.StudioTextAlignLeft),
	}
	switch rodzaj {
	case shared.StudioListKindBullet:
		zrodlo := shared.StudioBulletSource(shared.StudioBulletSourceCharacter)
		wynik.BulletSource = &zrodlo
		wynik.BulletCharacter = postacWskaznikTekstu(
			listyZnakiPoziomow[(poziom-1)%len(listyZnakiPoziomow)])
	case shared.StudioListKindNumber:
		format := shared.StudioListNumberFormat(shared.StudioListNumberFormatArabic)
		wynik.NumberFormat = &format
		wynik.Pattern = postacWskaznikTekstu("%" + postacZapisLiczby(poziom) + ".")
		wynik.StartAt = postacWskaznikLiczby(1)
	case shared.StudioListKindMultilevel:
		format := shared.StudioListNumberFormat(shared.StudioListNumberFormatArabic)
		wynik.NumberFormat = &format
		// Numeracja prawnicza wielopoziomowa: 1., 1.1., 1.1.2. — wzór składa się
		// z numerów wszystkich poziomów nadrzędnych, bo tak ją czyta Operator
		// w pismie z podstawami prawnymi.
		czesci := make([]string, 0, poziom)
		for i := 1; i <= poziom; i++ {
			czesci = append(czesci, "%"+postacZapisLiczby(i))
		}
		wynik.Pattern = postacWskaznikTekstu(strings.Join(czesci, ".") + ".")
		wynik.StartAt = postacWskaznikLiczby(1)
	}
	return wynik
}

// listyZapewnijPoziomy dokłada definicji listy poziomy aż do wskazanego. Poziom,
// na którym akapit staje, musi mieć swoje nastawy — inaczej okno nie wiedziałoby,
// jakim znakiem punkt narysować.
func listyZapewnijPoziomy(definicja *shared.StudioListDefinition, poziom int) {
	maPoziom := func(numer int) bool {
		for _, zastany := range definicja.Levels {
			if zastany.Level == numer {
				return true
			}
		}
		return false
	}
	for numer := 1; numer <= poziom; numer++ {
		if maPoziom(numer) {
			continue
		}
		definicja.Levels = append(definicja.Levels, listyPoziomWzorcowy(definicja.Kind, numer))
	}
	sort.SliceStable(definicja.Levels, func(i, j int) bool {
		return definicja.Levels[i].Level < definicja.Levels[j].Level
	})
}

// listyPoziom znajduje poziom definicji.
func listyPoziom(definicja *shared.StudioListDefinition, poziom int) *shared.StudioListLevel {
	for i := range definicja.Levels {
		if definicja.Levels[i].Level == poziom {
			return &definicja.Levels[i]
		}
	}
	return nil
}

// listyWciecieAkapitu przenosi wcięcie poziomu na postać akapitu, żeby linijka
// miała czym pokazać znaczniki wcięcia. Wysunięcie pierwszego wiersza jest
// ujemne — na tym stoi wygląd punktu, którego znak wisi po lewej stronie tekstu.
func listyWciecieAkapitu(poziom *shared.StudioListLevel) shared.StudioParagraphFormat {
	zmiana := shared.StudioParagraphFormat{}
	if poziom == nil {
		return zmiana
	}
	wciecie := listySzerokoscPoziomuMm * float64(poziom.Level)
	if poziom.IndentMm != nil {
		wciecie = *poziom.IndentMm
	}
	wysuniecie := listySzerokoscPoziomuMm
	if poziom.HangingMm != nil {
		wysuniecie = *poziom.HangingMm
	}
	zmiana.IndentLeftMm = postacWskaznikMiary(wciecie)
	zmiana.FirstLineIndentMm = postacWskaznikMiary(-wysuniecie)
	return zmiana
}

// listyBlokiListy oddaje wskazania bloków należących do wskazanej listy,
// w kolejności czytania.
func listyBlokiListy(forma *shared.StudioDocumentForm, kod string) []int {
	wskazania := make([]int, 0, 8)
	for i := range forma.Blocks {
		blok := &forma.Blocks[i]
		if blok.Paragraph == nil || blok.Paragraph.ListId == nil {
			continue
		}
		if *blok.Paragraph.ListId == kod {
			wskazania = append(wskazania, i)
		}
	}
	return wskazania
}

// ── Czynności ───────────────────────────────────────────────────────────────

// ZastosujListe zakłada wypunktowanie, numerację albo listę wielopoziomową na
// wskazanym fragmencie (`studio.list.apply`).
//
// Rodzaj `none` ZDEJMUJE listę — i jest czynnością prawdziwą, nie brakiem.
// Zdjęcie listy z fragmentu, który do żadnej listy nie należy, wraca odmową
// nazwaną: cisza kazałaby Operatorowi sądzić, że lista została zdjęta, choć jej
// tam nigdy nie było.
func (a *adapterStudia) ZastosujListe(ctx context.Context,
	z shared.StudioListApplyRequest) (shared.StudioListApplyResponse, error) {

	if err := listyRodzajZnany(z.Kind); err != nil {
		return shared.StudioListApplyResponse{}, err
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioListApplyResponse{}, err
	}
	if z.NumberFormat != nil && strings.TrimSpace(string(*z.NumberFormat)) != "" {
		if err := listyFormatZnany(*z.NumberFormat); err != nil {
			return shared.StudioListApplyResponse{}, err
		}
	}
	poziom := 1
	if z.Level != nil {
		poziom = *z.Level
	}
	if poziom < 1 || poziom > listyGlebokoscMaksymalna {
		return shared.StudioListApplyResponse{}, bladWskazaniaStudio(
			"poziom listy liczy się od jednego do " +
				postacZapisLiczby(listyGlebokoscMaksymalna) + "; podano " +
				postacZapisLiczby(poziom))
	}
	autor := postacAutor(z.Author)
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)
	bilans := shared.StudioActionBalance{Skipped: pominiete}

	if z.Kind == shared.StudioListKindNone {
		return a.listyZdejmij(ctx, stan, od, do, odcinki, bilans, autor)
	}

	// Lista wskazana kodem musi istnieć; brak wskazania zakłada listę nową.
	// Wskazanie kodu, którego dokument nie ma, jest odmową nazwaną — cicha zamiana
	// go na listę nową dałaby dwie listy tam, gdzie Operator chciał jednej.
	var definicja *shared.StudioListDefinition
	nowa := false
	if z.ListId != nil && strings.TrimSpace(*z.ListId) != "" {
		definicja = listyDefinicja(&stan.forma, *z.ListId)
		if definicja == nil {
			return shared.StudioListApplyResponse{}, bladWskazaniaStudio(
				"listy „" + strings.TrimSpace(*z.ListId) + "” dokument " +
					stan.dokument.Kod + " nie ma; listy dokumentu stoją w jego postaci, " +
					"którą oddaje komenda studio.document.form.get")
		}
		if definicja.Kind != z.Kind {
			definicja.Kind = z.Kind
			bilans.Applied++
		}
	} else {
		stan.forma.Lists = append(stan.forma.Lists, shared.StudioListDefinition{
			Id:   nowyIdentyfikator(przedrostekListyPostaci),
			Kind: z.Kind,
		})
		definicja = &stan.forma.Lists[len(stan.forma.Lists)-1]
		nowa = true
	}
	if z.StartAt != nil {
		if *z.StartAt < 0 {
			return shared.StudioListApplyResponse{}, bladWskazaniaStudio(
				"numeracja listy nie zaczyna się od liczby ujemnej")
		}
		definicja.StartAt = postacWskaznikLiczby(*z.StartAt)
	}
	listyZapewnijPoziomy(definicja, poziom)
	if z.NumberFormat != nil && strings.TrimSpace(string(*z.NumberFormat)) != "" {
		if wskazany := listyPoziom(definicja, poziom); wskazany != nil {
			format := *z.NumberFormat
			wskazany.NumberFormat = &format
		}
	}
	kod := definicja.Id
	nastawyPoziomu := listyPoziom(definicja, poziom)
	wciecie := listyWciecieAkapitu(nastawyPoziomu)

	for _, odcinek := range odcinki {
		for _, wskazanie := range postacBlokiZakresu(&stan.forma, odcinek[0], odcinek[1]) {
			blok := &stan.forma.Blocks[wskazanie]
			zmiana := wciecie
			zmiana.ListId = postacWskaznikTekstu(kod)
			zmiana.ListLevel = postacWskaznikLiczby(poziom)
			blok.Paragraph = postacScalAkapit(blok.Paragraph, zmiana)
			bilans.Applied++
		}
	}
	if bilans.Applied == 0 {
		if nowa {
			// Lista założona, a nie zastosowana, byłaby definicją bez ani jednego
			// punktu — czyli zapisem, którego Operator nigdzie nie zobaczy.
			stan.forma.Lists = stan.forma.Lists[:len(stan.forma.Lists)-1]
		}
		return shared.StudioListApplyResponse{}, bladWskazaniaStudio(
			"lista nie miała na czym stanąć: zakres " + postacZapisZakresu(od, do) +
				" nie obejmuje ani jednego akapitu wolnego od blokady")
	}
	stan.opisCzynnosci = "lista „" + string(z.Kind) + "” zastosowana " +
		postacZapisZakresu(od, do) + " na poziomie " + postacZapisLiczby(poziom)
	bilans.Note = postacWskaznikTekstu(stan.opisCzynnosci)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindListChange, od, do, bilans)
	if err != nil {
		return shared.StudioListApplyResponse{}, err
	}
	wynikowa := listyDefinicja(&forma, kod)
	if wynikowa == nil {
		return shared.StudioListApplyResponse{}, postacBladZaplecza(
			"lista " + kod + " zastosowana, ale nie wróciła z postaci dokumentu")
	}
	return shared.StudioListApplyResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, List: *wynikowa,
	}, nil
}

// listyZdejmij zdejmuje listę z akapitów zakresu — droga rodzaju `none`.
func (a *adapterStudia) listyZdejmij(ctx context.Context, stan *stanPostaci,
	od, do int, odcinki [][2]int, bilans shared.StudioActionBalance,
	autor shared.StudioAuthor) (shared.StudioListApplyResponse, error) {

	zdjeta := ""
	for _, odcinek := range odcinki {
		for _, wskazanie := range postacBlokiZakresu(&stan.forma, odcinek[0], odcinek[1]) {
			blok := &stan.forma.Blocks[wskazanie]
			if blok.Paragraph == nil || blok.Paragraph.ListId == nil {
				continue
			}
			if zdjeta == "" {
				zdjeta = *blok.Paragraph.ListId
			}
			blok.Paragraph.ListId = nil
			blok.Paragraph.ListLevel = nil
			// Wcięcie wniesione listą schodzi razem z nią, bo było jej częścią.
			// Zostawienie go dałoby akapit bez punktu, ale wciśnięty tam, gdzie
			// stał punkt.
			blok.Paragraph.IndentLeftMm = postacWskaznikMiary(0)
			blok.Paragraph.FirstLineIndentMm = postacWskaznikMiary(0)
			bilans.Applied++
		}
	}
	if bilans.Applied == 0 {
		return shared.StudioListApplyResponse{}, bladWskazaniaStudio(
			"nie ma czego zdjąć: zakres " + postacZapisZakresu(od, do) +
				" nie obejmuje ani jednego akapitu należącego do listy")
	}
	stan.opisCzynnosci = "lista zdjęta " + postacZapisZakresu(od, do)
	bilans.Note = postacWskaznikTekstu(stan.opisCzynnosci)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindListChange, od, do, bilans)
	if err != nil {
		return shared.StudioListApplyResponse{}, err
	}
	// Odpowiedź niesie definicję listy, z której punkty zeszły — Operator ma
	// wiedzieć, czego zdjęcie dotyczyło, a nie dostać pole puste.
	wynikowa := shared.StudioListDefinition{Id: zdjeta, Kind: shared.StudioListKindNone}
	if zastana := listyDefinicja(&forma, zdjeta); zastana != nil {
		wynikowa = *zastana
		wynikowa.Kind = shared.StudioListKindNone
	}
	return shared.StudioListApplyResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, List: wynikowa,
	}, nil
}

// UstawNumeracjeListy ustawia format numeracji poziomu listy
// (`studio.list.numbering.set`).
//
// Format prawniczy wielopoziomowy (1.1.2) bierze się z pola `pattern`, a nie
// z domysłu: wzór numeru jest jawny, żeby pismo z podstawami prawnymi dało się
// ponumerować dokładnie tak, jak wymaga wzór urzędowy.
func (a *adapterStudia) UstawNumeracjeListy(ctx context.Context,
	z shared.StudioListNumberingSetRequest) (shared.StudioListNumberingSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioListNumberingSetResponse{}, err
	}
	definicja, err := listyDefinicjaZadania(&stan.forma, z.ListId, stan.dokument.Kod)
	if err != nil {
		return shared.StudioListNumberingSetResponse{}, err
	}
	if z.Level < 1 || z.Level > listyGlebokoscMaksymalna {
		return shared.StudioListNumberingSetResponse{}, bladWskazaniaStudio(
			"poziom listy liczy się od jednego do " +
				postacZapisLiczby(listyGlebokoscMaksymalna) + "; podano " +
				postacZapisLiczby(z.Level))
	}
	if z.NumberFormat != nil && strings.TrimSpace(string(*z.NumberFormat)) != "" {
		if err := listyFormatZnany(*z.NumberFormat); err != nil {
			return shared.StudioListNumberingSetResponse{}, err
		}
	}
	autor := postacAutor(z.Author)
	listyZapewnijPoziomy(definicja, z.Level)
	poziom := listyPoziom(definicja, z.Level)
	if poziom == nil {
		return shared.StudioListNumberingSetResponse{}, postacBladZaplecza(
			"poziom " + postacZapisLiczby(z.Level) + " listy " + definicja.Id +
				" nie powstał, choć został zapewniony")
	}

	zmian := 0
	if z.NumberFormat != nil && strings.TrimSpace(string(*z.NumberFormat)) != "" {
		format := *z.NumberFormat
		poziom.NumberFormat = &format
		// Numeracja zdejmuje znak wypunktowania tego poziomu: punkt nie ma
		// jednocześnie kropki i numeru.
		poziom.BulletSource, poziom.BulletCharacter, poziom.BulletAssetId = nil, nil, nil
		zmian++
	}
	if z.Pattern != nil {
		if strings.TrimSpace(*z.Pattern) == "" {
			poziom.Pattern = nil
		} else {
			poziom.Pattern = postacWskaznikTekstu(*z.Pattern)
		}
		zmian++
	}
	if z.StartAt != nil {
		if *z.StartAt < 0 {
			return shared.StudioListNumberingSetResponse{}, bladWskazaniaStudio(
				"numeracja poziomu listy nie zaczyna się od liczby ujemnej")
		}
		poziom.StartAt = postacWskaznikLiczby(*z.StartAt)
		zmian++
	}
	if zmian == 0 {
		return shared.StudioListNumberingSetResponse{}, bladWskazaniaStudio(
			"numeracja poziomu bez ani jednej rzeczy do ustawienia — żądanie nie niosło " +
				"formatu, wzoru numeru ani punktu startu")
	}
	if definicja.Kind == shared.StudioListKindBullet || definicja.Kind == shared.StudioListKindNone {
		// Lista, której poziom dostał numerację, przestaje być wypunktowaniem.
		// Zostawienie rodzaju `bullet` dałoby definicję kłamiącą o sobie samej.
		definicja.Kind = shared.StudioListKindNumber
	}

	kod := definicja.Id
	wskazania := listyBlokiListy(&stan.forma, kod)
	stan.opisCzynnosci = "numeracja poziomu " + postacZapisLiczby(z.Level) + " listy " +
		kod + " ustawiona; przestawiła " +
		postacLiczebnik(len(wskazania), "punkt", "punkty", "punktów")
	bilans := shared.StudioActionBalance{
		Applied: zmian + len(wskazania),
		Skipped: []shared.StudioSkippedItem{},
		Note:    postacWskaznikTekstu(stan.opisCzynnosci),
	}
	od, do := listyZakresListy(&stan.forma, wskazania)
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindListChange, od, do, bilans)
	if err != nil {
		return shared.StudioListNumberingSetResponse{}, err
	}
	wynikowa := listyDefinicja(&forma, kod)
	if wynikowa == nil {
		return shared.StudioListNumberingSetResponse{}, postacBladZaplecza(
			"lista " + kod + " zmieniona, ale nie wróciła z postaci dokumentu")
	}
	return shared.StudioListNumberingSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, List: *wynikowa,
	}, nil
}

// UstawPunktatorListy ustawia znak wypunktowania poziomu
// (`studio.list.bullet.set`).
//
// Znak bierze się ze znaku gotowego, dowolnego symbolu, ikony albo obrazu
// własnego — wszystkie cztery źródła stoją w kontrakcie i wszystkie tu działają.
// Ikona i obraz idą ZASOBEM (`bulletAssetId`), a nie wklejonym rysunkiem: ikony
// wystawia moduł Design (`design.icon.library.search`), a zasoby magazyn rdzenia,
// i drugiego rachunku ikony Studio nie zakłada.
func (a *adapterStudia) UstawPunktatorListy(ctx context.Context,
	z shared.StudioListBulletSetRequest) (shared.StudioListBulletSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioListBulletSetResponse{}, err
	}
	definicja, err := listyDefinicjaZadania(&stan.forma, z.ListId, stan.dokument.Kod)
	if err != nil {
		return shared.StudioListBulletSetResponse{}, err
	}
	if z.Level < 1 || z.Level > listyGlebokoscMaksymalna {
		return shared.StudioListBulletSetResponse{}, bladWskazaniaStudio(
			"poziom listy liczy się od jednego do " +
				postacZapisLiczby(listyGlebokoscMaksymalna) + "; podano " +
				postacZapisLiczby(z.Level))
	}
	if z.BulletSource != nil && strings.TrimSpace(string(*z.BulletSource)) != "" {
		if err := listyZrodloZnane(*z.BulletSource); err != nil {
			return shared.StudioListBulletSetResponse{}, err
		}
	}
	autor := postacAutor(z.Author)
	listyZapewnijPoziomy(definicja, z.Level)
	poziom := listyPoziom(definicja, z.Level)
	if poziom == nil {
		return shared.StudioListBulletSetResponse{}, postacBladZaplecza(
			"poziom " + postacZapisLiczby(z.Level) + " listy " + definicja.Id +
				" nie powstał, choć został zapewniony")
	}

	zmian := 0
	if z.BulletSource != nil && strings.TrimSpace(string(*z.BulletSource)) != "" {
		zrodlo := *z.BulletSource
		poziom.BulletSource = &zrodlo
		zmian++
	}
	if z.BulletCharacter != nil {
		if strings.TrimSpace(*z.BulletCharacter) == "" {
			poziom.BulletCharacter = nil
		} else {
			poziom.BulletCharacter = postacWskaznikTekstu(*z.BulletCharacter)
		}
		zmian++
	}
	if z.BulletAssetId != nil {
		if strings.TrimSpace(*z.BulletAssetId) == "" {
			poziom.BulletAssetId = nil
		} else {
			poziom.BulletAssetId = postacWskaznikTekstu(strings.TrimSpace(*z.BulletAssetId))
		}
		zmian++
	}
	if z.IndentMm != nil {
		if *z.IndentMm < 0 {
			return shared.StudioListBulletSetResponse{}, bladWskazaniaStudio(
				"wcięcie poziomu listy liczy się od lewego marginesu i nie może być ujemne")
		}
		poziom.IndentMm = postacWskaznikMiary(*z.IndentMm)
		zmian++
	}
	if z.HangingMm != nil {
		if *z.HangingMm < 0 {
			return shared.StudioListBulletSetResponse{}, bladWskazaniaStudio(
				"odstęp znaku od tekstu nie może być ujemny")
		}
		poziom.HangingMm = postacWskaznikMiary(*z.HangingMm)
		zmian++
	}
	if z.Align != nil && strings.TrimSpace(string(*z.Align)) != "" {
		wyrownanie := *z.Align
		poziom.Align = &wyrownanie
		zmian++
	}
	if zmian == 0 {
		return shared.StudioListBulletSetResponse{}, bladWskazaniaStudio(
			"znak wypunktowania bez ani jednej rzeczy do ustawienia — żądanie nie niosło " +
				"źródła, znaku, zasobu, wcięcia, odstępu ani wyrównania")
	}

	// Źródło niepodane wynika ze wskazanego znaku albo zasobu — Operator, który
	// podał sam znak, nie musi jeszcze raz mówić, że to znak.
	if poziom.BulletSource == nil {
		zrodlo := shared.StudioBulletSource(shared.StudioBulletSourceCharacter)
		if poziom.BulletAssetId != nil {
			zrodlo = shared.StudioBulletSource(shared.StudioBulletSourceImage)
		}
		poziom.BulletSource = &zrodlo
	}
	switch *poziom.BulletSource {
	case shared.StudioBulletSourceCharacter, shared.StudioBulletSourceSymbol:
		if poziom.BulletCharacter == nil {
			return shared.StudioListBulletSetResponse{}, bladWskazaniaStudio(
				"wypunktowanie ze źródła „" + string(*poziom.BulletSource) +
					"” bez znaku — podaj pole bulletCharacter")
		}
	case shared.StudioBulletSourceIcon, shared.StudioBulletSourceImage:
		if poziom.BulletAssetId == nil {
			return shared.StudioListBulletSetResponse{}, bladWskazaniaStudio(
				"wypunktowanie ze źródła „" + string(*poziom.BulletSource) +
					"” bez zasobu — podaj pole bulletAssetId. Ikony wystawia moduł " +
					"Design komendą design.icon.library.search, obrazy własne magazyn zasobów rdzenia")
		}
	}
	// Znak wypunktowania zdejmuje numerację tego poziomu — punkt nie ma
	// jednocześnie kropki i numeru.
	poziom.NumberFormat, poziom.Pattern = nil, nil
	if definicja.Kind == shared.StudioListKindNumber || definicja.Kind == shared.StudioListKindNone {
		definicja.Kind = shared.StudioListKindBullet
	}

	kod := definicja.Id
	wskazania := listyBlokiListy(&stan.forma, kod)
	wciecie := listyWciecieAkapitu(poziom)
	przestawione := 0
	for _, wskazanie := range wskazania {
		blok := &stan.forma.Blocks[wskazanie]
		if blok.Paragraph == nil || blok.Paragraph.ListLevel == nil ||
			*blok.Paragraph.ListLevel != z.Level {
			continue
		}
		blok.Paragraph = postacScalAkapit(blok.Paragraph, wciecie)
		przestawione++
	}
	stan.opisCzynnosci = "znak wypunktowania poziomu " + postacZapisLiczby(z.Level) +
		" listy " + kod + " ustawiony; przestawił " +
		postacLiczebnik(przestawione, "punkt", "punkty", "punktów")
	bilans := shared.StudioActionBalance{
		Applied: zmian + przestawione,
		Skipped: []shared.StudioSkippedItem{},
		Note:    postacWskaznikTekstu(stan.opisCzynnosci),
	}
	od, do := listyZakresListy(&stan.forma, wskazania)
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindListChange, od, do, bilans)
	if err != nil {
		return shared.StudioListBulletSetResponse{}, err
	}
	wynikowa := listyDefinicja(&forma, kod)
	if wynikowa == nil {
		return shared.StudioListBulletSetResponse{}, postacBladZaplecza(
			"lista " + kod + " zmieniona, ale nie wróciła z postaci dokumentu")
	}
	return shared.StudioListBulletSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, List: *wynikowa,
	}, nil
}

// PrzestawPoziomListy zwiększa albo zmniejsza poziom listy wskazanego fragmentu
// (`studio.list.level.indent`).
//
// Wcięcie i odstęp poziomu idą ZA nim — to jest cała treść tej czynności:
// „zwiększ wcięcie" w pakiecie biurowym zmienia poziom punktu, a nie samo
// wcięcie akapitu. Zmniejszenie poniżej pierwszego poziomu nie zdejmuje listy
// samo z siebie: zdjęcie listy jest osobną czynnością (`studio.list.apply`
// z rodzajem `none`) i cichego zdejmowania tu nie ma.
func (a *adapterStudia) PrzestawPoziomListy(ctx context.Context,
	z shared.StudioListLevelIndentRequest) (shared.StudioListLevelIndentResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioListLevelIndentResponse{}, err
	}
	if z.Step == 0 {
		return shared.StudioListLevelIndentResponse{}, bladWskazaniaStudio(
			"przestawienie poziomu listy o zero poziomów nie jest czynnością — " +
				"podaj krok dodatni, żeby zwiększyć wcięcie, albo ujemny, żeby zmniejszyć")
	}
	autor := postacAutor(z.Author)
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)
	bilans := shared.StudioActionBalance{Skipped: pominiete}

	poza := 0
	naGranicy := 0
	for _, odcinek := range odcinki {
		for _, wskazanie := range postacBlokiZakresu(&stan.forma, odcinek[0], odcinek[1]) {
			blok := &stan.forma.Blocks[wskazanie]
			if blok.Paragraph == nil || blok.Paragraph.ListId == nil {
				poza++
				continue
			}
			zastany := 1
			if blok.Paragraph.ListLevel != nil {
				zastany = *blok.Paragraph.ListLevel
			}
			docelowy := zastany + z.Step
			if docelowy < 1 {
				docelowy = 1
			}
			if docelowy > listyGlebokoscMaksymalna {
				docelowy = listyGlebokoscMaksymalna
			}
			if docelowy == zastany {
				naGranicy++
				continue
			}
			definicja := listyDefinicja(&stan.forma, *blok.Paragraph.ListId)
			if definicja == nil {
				poza++
				continue
			}
			listyZapewnijPoziomy(definicja, docelowy)
			zmiana := listyWciecieAkapitu(listyPoziom(definicja, docelowy))
			zmiana.ListLevel = postacWskaznikLiczby(docelowy)
			blok.Paragraph = postacScalAkapit(blok.Paragraph, zmiana)
			bilans.Applied++
		}
	}
	if poza > 0 {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "akapit poza listą",
			Detail: postacWskaznikTekstu(postacLiczebnik(poza, "akapit", "akapity", "akapitów") +
				" w zakresie nie należy do żadnej listy, więc nie ma poziomu do przestawienia; " +
				"listę zakłada komenda studio.list.apply"),
		})
	}
	if naGranicy > 0 {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "poziom listy na granicy",
			Detail: postacWskaznikTekstu(postacLiczebnik(naGranicy, "punkt", "punkty", "punktów") +
				" stoi już na poziomie krańcowym — pierwszym albo " +
				postacZapisLiczby(listyGlebokoscMaksymalna) + "; głębiej ani wyżej nie ma gdzie"),
		})
	}
	if bilans.Applied == 0 {
		return shared.StudioListLevelIndentResponse{}, bladWskazaniaStudio(
			"poziom listy nie miał na czym stanąć: zakres " + postacZapisZakresu(od, do) +
				" nie obejmuje ani jednego punktu listy, którego poziom da się przestawić")
	}
	kierunek := "zwiększony"
	if z.Step < 0 {
		kierunek = "zmniejszony"
	}
	stan.opisCzynnosci = "poziom listy " + kierunek + " " + postacZapisZakresu(od, do)
	bilans.Note = postacWskaznikTekstu(stan.opisCzynnosci)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindListChange, od, do, bilans)
	if err != nil {
		return shared.StudioListLevelIndentResponse{}, err
	}
	return shared.StudioListLevelIndentResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana, ActionId: stan.czynnosc,
	}, nil
}

// WznowNumeracjeListy wznawia numerację listy od wskazanego miejsca
// (`studio.list.restart`).
//
// Wznowienie jest ROZDZIELENIEM listy: punkty od wskazanego miejsca w dół
// przechodzą do listy nowej o tych samych poziomach i własnym punkcie startu.
// Powód stoi w nagłówku pliku — akapit nie ma w kontrakcie miejsca na punkt
// startu, a udawanie go polem, którego nie ma, byłoby zapisem nie do odczytania.
func (a *adapterStudia) WznowNumeracjeListy(ctx context.Context,
	z shared.StudioListRestartRequest) (shared.StudioListRestartResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioListRestartResponse{}, err
	}
	definicja, err := listyDefinicjaZadania(&stan.forma, z.ListId, stan.dokument.Kod)
	if err != nil {
		return shared.StudioListRestartResponse{}, err
	}
	autor := postacAutor(z.Author)
	miejsce, _ := postacZakres(&z.Offset, &z.Offset, postacDlugosc(&stan.forma))
	poczatek := 1
	if z.StartAt != nil {
		if *z.StartAt < 0 {
			return shared.StudioListRestartResponse{}, bladWskazaniaStudio(
				"numeracja listy nie wznawia się od liczby ujemnej")
		}
		poczatek = *z.StartAt
	}
	kod := definicja.Id

	// Lista nowa niesie poziomy listy zastanej KOPIĄ, nie wskaźnikiem — inaczej
	// zmiana znaku wypunktowania w jednej przestawiłaby drugą i wznowienie
	// przestałoby być rozdzieleniem.
	wznowiona := shared.StudioListDefinition{
		Id:      nowyIdentyfikator(przedrostekListyPostaci),
		Kind:    definicja.Kind,
		StartAt: postacWskaznikLiczby(poczatek),
		Levels:  make([]shared.StudioListLevel, 0, len(definicja.Levels)),
	}
	for _, poziom := range definicja.Levels {
		kopia := poziom
		if kopia.Level == 1 {
			kopia.StartAt = postacWskaznikLiczby(poczatek)
		}
		wznowiona.Levels = append(wznowiona.Levels, kopia)
	}

	przeniesione := 0
	for _, wskazanie := range listyBlokiListy(&stan.forma, kod) {
		blok := &stan.forma.Blocks[wskazanie]
		start := 0
		if blok.RangeStart != nil {
			start = *blok.RangeStart
		}
		koniec := start
		if blok.RangeEnd != nil {
			koniec = *blok.RangeEnd
		}
		if koniec < miejsce {
			continue
		}
		blok.Paragraph = postacScalAkapit(blok.Paragraph,
			shared.StudioParagraphFormat{ListId: postacWskaznikTekstu(wznowiona.Id)})
		przeniesione++
	}
	if przeniesione == 0 {
		return shared.StudioListRestartResponse{}, bladWskazaniaStudio(
			"numeracji listy " + kod + " nie ma od czego wznowić: od miejsca " +
				postacZapisZakresu(miejsce, miejsce) + " w dół nie stoi ani jeden punkt " +
				"tej listy")
	}
	stan.forma.Lists = append(stan.forma.Lists, wznowiona)
	stan.opisCzynnosci = "numeracja listy " + kod + " wznowiona od " +
		postacZapisLiczby(poczatek) + " " + postacZapisZakresu(miejsce, miejsce) +
		"; " + postacLiczebnik(przeniesione, "punkt", "punkty", "punktów") +
		" przeszło do listy " + wznowiona.Id
	bilans := shared.StudioActionBalance{
		Applied: przeniesione,
		Skipped: []shared.StudioSkippedItem{},
		Note:    postacWskaznikTekstu(stan.opisCzynnosci),
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindListChange,
		miejsce, postacDlugosc(&stan.forma), bilans)
	if err != nil {
		return shared.StudioListRestartResponse{}, err
	}
	return shared.StudioListRestartResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana, ActionId: stan.czynnosc,
	}, nil
}

// listyDefinicjaZadania odnajduje listę wskazaną żądaniem i odmawia nazwanie,
// gdy jej nie ma. Jedno miejsce dla trzech czynności, które listy wymagają.
func listyDefinicjaZadania(forma *shared.StudioDocumentForm, kod,
	dokumentKod string) (*shared.StudioListDefinition, error) {

	szukany := strings.TrimSpace(kod)
	if szukany == "" {
		return nil, bladWskazaniaStudio("czynność na liście bez wskazania listy")
	}
	definicja := listyDefinicja(forma, szukany)
	if definicja == nil {
		nazwy := make([]string, 0, len(forma.Lists))
		for _, lista := range forma.Lists {
			nazwy = append(nazwy, lista.Id)
		}
		if len(nazwy) == 0 {
			return nil, bladWskazaniaStudio("dokument " + dokumentKod +
				" nie ma ani jednej listy; listę zakłada komenda studio.list.apply")
		}
		return nil, bladWskazaniaStudio("listy „" + szukany + "” dokument " + dokumentKod +
			" nie ma; listy dokumentu: " + strings.Join(nazwy, ", "))
	}
	return definicja, nil
}

// listyZakresListy oddaje zakres w znakach, który punkty listy obejmują — wpis
// dziennika i zmiana śledzona muszą wiedzieć, czego zmiana dotknęła.
func listyZakresListy(forma *shared.StudioDocumentForm, wskazania []int) (int, int) {
	if len(wskazania) == 0 {
		return 0, postacDlugosc(forma)
	}
	od, do := -1, 0
	for _, wskazanie := range wskazania {
		blok := &forma.Blocks[wskazanie]
		if blok.RangeStart != nil && (od < 0 || *blok.RangeStart < od) {
			od = *blok.RangeStart
		}
		if blok.RangeEnd != nil && *blok.RangeEnd > do {
			do = *blok.RangeEnd
		}
	}
	if od < 0 {
		od = 0
	}
	return od, do
}
