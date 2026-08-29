// Plik odkłada i czyta postać dokumentu modułu Studio: arkusz stylów, nastawy strony, sekcje, bloki, tabele, obiekty, listy, aparat i pola, jedynym miejscem, przez które postać dojeżdża do rdzenia i wraca z niego nieuszkodzona.
package core

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// wejscieZasiegPostaci trzyma wiersze postaci osobno od profili wydania, żeby profil nie nadpisywał postaci przy zapisie.
	wejscieZasiegPostaci = "postac-dokumentu"
	// wejscieZasiegSzablonu trzyma postać wzorcową szablonu pisma, osobno od postaci dokumentów zakładanych z niego.
	wejscieZasiegSzablonu = "postac-szablonu"
	// wejscieZasiegPochodzenia trzyma zapisy pochodzenia fragmentów dokumentu, osobno od samej postaci dokumentu.
	wejscieZasiegPochodzenia = "pochodzenie-dokumentu"

	// Przedrostki kluczy wiersza. Klucz jest złożony z przedrostka i kodu
	// dokumentu albo szablonu, więc upsert trafia zawsze w ten sam wiersz.
	wejscieKluczPostaci     = "studio-postac-"
	wejscieKluczSzablonu    = "studio-szpostac-"
	wejscieKluczPochodzenia = "studio-poch-"

	// Przedrostki identyfikatorów bytów nadawanych przez ten odcinek: pochodzenia, obiektu, bloku, sekcji, tabeli, szablonu i czynności.
	przedrostekPochodzeniaStudia = "studio-poch-"
	przedrostekObiektuStudia     = "studio-obj-"
	przedrostekBlokuStudia       = "studio-blok-"
	przedrostekSekcjiStudia      = "studio-sek-"
	przedrostekTabeliStudia      = "studio-tab-"
	przedrostekSzablonuStudia    = "studio-szab-"
	przedrostekCzynnosciWejscia  = "studio-czyn-"
)

// wejscieNazwyStylowDomyslnych wymienia arkusz stylów nowego dokumentu.
// Nazwy są nazwami pełnymi, nie kodami — numeracja wymyślona jest w tym
// produkcie zakazana, a `naglowek-1` jest nazwą poziomu, nie kodem.
const (
	wejscieStylTekstZasadniczy = "tekst zasadniczy"
	wejscieStylNaglowek1       = "naglowek-1"
	wejscieStylNaglowek2       = "naglowek-2"
	wejscieStylNaglowek3       = "naglowek-3"
	wejscieStylNaglowek4       = "naglowek-4"
	wejscieStylNaglowek5       = "naglowek-5"
	wejscieStylNaglowek6       = "naglowek-6"
	wejscieStylCytat           = "cytat"
	wejscieStylPodpis          = "podpis"
	wejscieStylPrzypis         = "przypis"
)

// wejscieTeraz oddaje chwilę bieżącą w milisekundach epoki, jednostce, którą kontrakt niesie w polach czasu.
func wejscieTeraz() int64 {
	return time.Now().UnixMilli()
}

// wejscieWskaznikTekstu oddaje wskaźnik na napis albo nic, gdy napis pusty.
// Pole nieobowiązkowe wypełnione pustym napisem byłoby twierdzeniem „wartość
// jest i jest pusta" w miejscu, gdzie prawdą jest „wartości nie ma".
func wejscieWskaznikTekstu(tekst string) *string {
	if strings.TrimSpace(tekst) == "" {
		return nil
	}
	kopia := tekst
	return &kopia
}

// wejscieWskaznikCalkowity oddaje wskaźnik na liczbę całkowitą, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikCalkowity(wartosc int) *int {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikDlugi oddaje wskaźnik na liczbę 64-bitową, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikDlugi(wartosc int64) *int64 {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikRzeczywisty oddaje wskaźnik na liczbę rzeczywistą, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikRzeczywisty(wartosc float64) *float64 {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikLogiczny oddaje wskaźnik na wartość logiczną, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikLogiczny(wartosc bool) *bool {
	kopia := wartosc
	return &kopia
}

// ── Magazyn postaci ─────────────────────────────────────────────────────────

// wejsciePostacDokumentu czyta postać dokumentu. Dokument bez zapisanej postaci
// nie jest usterką — oddaje postać pustą o właściwym identyfikatorze, żeby
// wołający miał na czym pracować, zamiast rozstrzygać brak drugi raz u siebie.
func (a *adapterStudia) wejsciePostacDokumentu(ctx context.Context,
	kodDokumentu string) (shared.StudioDocumentForm, error) {

	pusta := shared.StudioDocumentForm{DocumentId: kodDokumentu}
	if a.repozytorium == nil {
		return pusta, nil
	}
	stan, err := a.postacWczytaj(ctx, kodDokumentu)
	if err != nil {
		// Dokumentu, którego nie ma, nie udaje się postacią pustą, żeby brak nie wyszedł jako brak stylu
		return pusta, err
	}
	return stan.forma, nil
}

// wejscieZapiszPostac utrwala postać dokumentu i podbija jej numer porządkowy.
//
// Numer porządkowy rośnie w rdzeniu, nie u wołającego: dwie czynności modelu
// wysłane naraz z tym samym numerem nie mają jak się rozminąć, jeśli numer
// nadaje jedno miejsce.
func (a *adapterStudia) wejscieZapiszPostac(ctx context.Context, kodDokumentu string,
	postac shared.StudioDocumentForm) (shared.StudioDocumentForm, error) {

	if a.repozytorium == nil {
		return shared.StudioDocumentForm{}, wejscieBladZaplecza(
			"serwer złożony bez repozytorium Studia — postaci dokumentu nie ma gdzie odłożyć")
	}
	dokument, err := a.repozytorium.Dokument(ctx, kodDokumentu)
	if err != nil {
		return shared.StudioDocumentForm{}, bladNieznanegoDokumentu(kodDokumentu, err)
	}
	postac.DocumentId = kodDokumentu

	// Numer porządkowy podbija baza przy zapisie, nie ten kod, inaczej dwa zapisy mogłyby dzielić numer
	stan := &stanPostaci{dokument: dokument, forma: postac}
	stan.tekstPrzed = wartoscTekstu(dokument.Tresc)
	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioDocumentForm{}, err
	}
	return stan.forma, nil
}

// wejsciePrzyjmijPostacZapisu utrwala postać przyjechaną polem `form` komendy `studio.document.save`: treść z pola `content` niesie litery, postać z pola `form` niesie strukturę, a brak pola `form` znaczy „bez zmiany postaci”.
func (a *adapterStudia) wejsciePrzyjmijPostacZapisu(ctx context.Context,
	dokument dane.DokumentStudia, surowa []byte) (dane.DokumentStudia, error) {

	var postac shared.StudioDocumentForm
	if err := json.Unmarshal(surowa, &postac); err != nil {
		return dokument, bladWskazaniaStudio("pole form nie jest postacią dokumentu " +
			"w kształcie StudioDocumentForm: " + err.Error())
	}
	postac.DocumentId = dokument.Kod

	stan := &stanPostaci{dokument: dokument, forma: postac}
	stan.tekstPrzed = wartoscTekstu(dokument.Tresc)
	stan.opisCzynnosci = "zapis dokumentu wraz z postacią"
	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return dokument, err
	}
	return stan.dokument, nil
}

// ── Postać nowa i domyślna ──────────────────────────────────────────────────

// wejscieNowaPostac składa postać dokumentu pustego: arkusz stylów nazwanych, nastawy strony, jedną sekcję i jeden akapit pusty, gotowy do pisania od pierwszego znaku.
func wejscieNowaPostac(kodDokumentu, nazwaNosnika string,
	orientacja *shared.StudioPageOrientation) shared.StudioDocumentForm {

	nastawy := wejscieDomyslneNastawyStrony(nazwaNosnika, orientacja)
	sekcja := shared.StudioSection{
		Id:         nowyIdentyfikator(przedrostekSekcjiStudia),
		Index:      0,
		RangeStart: 0,
		RangeEnd:   0,
		Start:      wejscieWskaznikPoczatkuSekcji(shared.StudioSectionStartContinuous),
		PageSetup:  &nastawy,
	}
	blok := shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuAkapit,
		SectionId: wejscieWskaznikTekstu(sekcja.Id),
		Paragraph: &shared.StudioParagraphFormat{
			StyleName: wejscieWskaznikTekstu(wejscieStylTekstZasadniczy),
		},
		RangeStart: wejscieWskaznikCalkowity(0),
		RangeEnd:   wejscieWskaznikCalkowity(0),
	}
	return shared.StudioDocumentForm{
		DocumentId: kodDokumentu,
		PageSetup:  &nastawy,
		Styles:     wejscieDomyslnyArkuszStylow(),
		Sections:   []shared.StudioSection{sekcja},
		Blocks:     []shared.StudioDocumentBlock{blok},
		Revision:   wejscieWskaznikDlugi(1),
		UpdatedAt:  wejscieWskaznikDlugi(wejscieTeraz()),
	}
}

// Rodzaje bloku dokumentu. Kontrakt trzyma je napisem, nie wyliczeniem, więc
// nazwy stoją tutaj raz, a nie rozsypane po literałach w dziesięciu miejscach.
const (
	wejscieRodzajBlokuAkapit  = "akapit"
	wejscieRodzajBlokuTabela  = "tabela"
	wejscieRodzajBlokuObiekt  = "obiekt"
	wejscieRodzajBlokuPodzial = "podzial"
)

// wejscieWskaznikPoczatkuSekcji oddaje wskaźnik na sposób rozpoczęcia sekcji, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikPoczatkuSekcji(wartosc shared.StudioSectionStart) *shared.StudioSectionStart {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikWyrownania oddaje wskaźnik na wyrównanie tekstu, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikWyrownania(wartosc shared.StudioTextAlign) *shared.StudioTextAlign {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikOrientacji oddaje wskaźnik na orientację strony, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikOrientacji(wartosc shared.StudioPageOrientation) *shared.StudioPageOrientation {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikRodzajuNosnika oddaje wskaźnik na rodzaj nośnika, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikRodzajuNosnika(wartosc shared.StudioPaperKind) *shared.StudioPaperKind {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikZrodlaObiektu oddaje wskaźnik na rodzaj obiektu, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikZrodlaObiektu(wartosc shared.StudioObjectSource) *shared.StudioObjectSource {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikAutora oddaje wskaźnik na autora czynności, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikAutora(wartosc shared.StudioAuthor) *shared.StudioAuthor {
	kopia := wartosc
	return &kopia
}

// wejscieAutorCzynnosci rozstrzyga autora czynności: brak wskazania znaczy Operator, a wartość spoza wyliczenia jest odmową, nie cichym zejściem na Operatora.
func wejscieAutorCzynnosci(wskazanie *shared.StudioAuthor) (shared.StudioAuthor, error) {
	if wskazanie == nil || strings.TrimSpace(string(*wskazanie)) == "" {
		return shared.StudioAuthorUzytkownik, nil
	}
	switch *wskazanie {
	case shared.StudioAuthorUzytkownik, shared.StudioAuthorModel:
		return *wskazanie, nil
	}
	return "", bladWskazaniaStudio("autor czynności " + string(*wskazanie) +
		" nie jest ani " + string(shared.StudioAuthorUzytkownik) +
		", ani " + string(shared.StudioAuthorModel))
}

// wejscieDomyslneNastawyStrony składa nastawy strony nowego dokumentu: A4 i marginesy 25 milimetrów są nastawą domyślną pisma urzędowego, a nazwa nośnika podana przez Operatora ma pierwszeństwo.
func wejscieDomyslneNastawyStrony(nazwaNosnika string,
	orientacja *shared.StudioPageOrientation) shared.StudioPageSetup {

	nazwa := strings.TrimSpace(nazwaNosnika)
	if nazwa == "" {
		nazwa = "A4"
	}
	kierunek := shared.StudioPageOrientation(shared.StudioPageOrientationPionowa)
	if orientacja != nil && strings.TrimSpace(string(*orientacja)) != "" {
		kierunek = *orientacja
	}
	return shared.StudioPageSetup{
		PageSize:     wejscieWskaznikTekstu(nazwa),
		Orientation:  wejscieWskaznikOrientacji(kierunek),
		MarginTop:    wejscieWskaznikCalkowity(25),
		MarginBottom: wejscieWskaznikCalkowity(25),
		MarginLeft:   wejscieWskaznikCalkowity(25),
		MarginRight:  wejscieWskaznikCalkowity(25),
		PaperKind:    wejscieWskaznikRodzajuNosnika(shared.StudioPaperKindSheet),
		MarginPreset: wejscieWskaznikTekstu("normalne"),
		Columns:      wejscieWskaznikCalkowity(1),
	}
}

// wejscieDomyslnyArkuszStylow składa arkusz stylów nazwanych nowego dokumentu: tekst zasadniczy, sześć poziomów nagłówków, cytat, podpis i przypis, gdzie nagłówki dziedziczą postać po tekście zasadniczym.
func wejscieDomyslnyArkuszStylow() []shared.StudioNamedStyle {
	zasadniczy := shared.StudioNamedStyle{
		Name:        wejscieStylTekstZasadniczy,
		DisplayName: wejscieWskaznikTekstu("Tekst zasadniczy"),
		Kind:        shared.StudioStyleKindParagraph,
		Builtin:     wejscieWskaznikLogiczny(true),
		Character: &shared.StudioCharacterFormat{
			FontFamily: wejscieWskaznikTekstu("Times New Roman"),
			FontSizePt: wejscieWskaznikRzeczywisty(12),
			Language:   wejscieWskaznikTekstu("pl-PL"),
		},
		Paragraph: &shared.StudioParagraphFormat{
			Align:            wejscieWskaznikWyrownania(shared.StudioTextAlignJustify),
			LineSpacingRule:  wejscieWskaznikZasadyInterlinii(shared.StudioLineSpacingRuleSingle),
			LineSpacingValue: wejscieWskaznikRzeczywisty(1),
			SpaceAfterPt:     wejscieWskaznikRzeczywisty(6),
			WidowControl:     wejscieWskaznikLogiczny(true),
			OutlineLevel:     wejscieWskaznikCalkowity(0),
		},
	}
	arkusz := []shared.StudioNamedStyle{zasadniczy}

	// Stopnie nagłówków maleją z poziomem, 18 do 11 punktów; poziom konspektu to poziom nagłówka.
	stopnie := []float64{18, 16, 14, 13, 12, 11}
	nazwy := []string{
		wejscieStylNaglowek1, wejscieStylNaglowek2, wejscieStylNaglowek3,
		wejscieStylNaglowek4, wejscieStylNaglowek5, wejscieStylNaglowek6,
	}
	widoczne := []string{
		"Nagłówek poziomu pierwszego", "Nagłówek poziomu drugiego",
		"Nagłówek poziomu trzeciego", "Nagłówek poziomu czwartego",
		"Nagłówek poziomu piątego", "Nagłówek poziomu szóstego",
	}
	for i, nazwa := range nazwy {
		arkusz = append(arkusz, shared.StudioNamedStyle{
			Name:        nazwa,
			DisplayName: wejscieWskaznikTekstu(widoczne[i]),
			Kind:        shared.StudioStyleKindParagraph,
			BasedOn:     wejscieWskaznikTekstu(wejscieStylTekstZasadniczy),
			NextStyle:   wejscieWskaznikTekstu(wejscieStylTekstZasadniczy),
			Builtin:     wejscieWskaznikLogiczny(true),
			Character: &shared.StudioCharacterFormat{
				FontFamily: wejscieWskaznikTekstu("Arial"),
				FontSizePt: wejscieWskaznikRzeczywisty(stopnie[i]),
				Bold:       wejscieWskaznikLogiczny(true),
			},
			Paragraph: &shared.StudioParagraphFormat{
				Align:         wejscieWskaznikWyrownania(shared.StudioTextAlignLeft),
				SpaceBeforePt: wejscieWskaznikRzeczywisty(12),
				SpaceAfterPt:  wejscieWskaznikRzeczywisty(6),
				KeepWithNext:  wejscieWskaznikLogiczny(true),
				OutlineLevel:  wejscieWskaznikCalkowity(i + 1),
			},
		})
	}

	arkusz = append(arkusz,
		shared.StudioNamedStyle{
			Name:        wejscieStylCytat,
			DisplayName: wejscieWskaznikTekstu("Cytat"),
			Kind:        shared.StudioStyleKindParagraph,
			BasedOn:     wejscieWskaznikTekstu(wejscieStylTekstZasadniczy),
			Builtin:     wejscieWskaznikLogiczny(true),
			Character: &shared.StudioCharacterFormat{
				Italic: wejscieWskaznikLogiczny(true),
			},
			Paragraph: &shared.StudioParagraphFormat{
				IndentLeftMm:  wejscieWskaznikRzeczywisty(10),
				IndentRightMm: wejscieWskaznikRzeczywisty(10),
			},
		},
		shared.StudioNamedStyle{
			Name:        wejscieStylPodpis,
			DisplayName: wejscieWskaznikTekstu("Podpis pod ilustracją"),
			Kind:        shared.StudioStyleKindParagraph,
			BasedOn:     wejscieWskaznikTekstu(wejscieStylTekstZasadniczy),
			Builtin:     wejscieWskaznikLogiczny(true),
			Character: &shared.StudioCharacterFormat{
				FontSizePt: wejscieWskaznikRzeczywisty(10),
				Italic:     wejscieWskaznikLogiczny(true),
			},
			Paragraph: &shared.StudioParagraphFormat{
				Align: wejscieWskaznikWyrownania(shared.StudioTextAlignCenter),
			},
		},
		shared.StudioNamedStyle{
			Name:        wejscieStylPrzypis,
			DisplayName: wejscieWskaznikTekstu("Przypis"),
			Kind:        shared.StudioStyleKindParagraph,
			BasedOn:     wejscieWskaznikTekstu(wejscieStylTekstZasadniczy),
			Builtin:     wejscieWskaznikLogiczny(true),
			Character: &shared.StudioCharacterFormat{
				FontSizePt: wejscieWskaznikRzeczywisty(9),
			},
			Paragraph: &shared.StudioParagraphFormat{
				Align:        wejscieWskaznikWyrownania(shared.StudioTextAlignLeft),
				SpaceAfterPt: wejscieWskaznikRzeczywisty(0),
			},
		},
	)
	return arkusz
}

// wejscieWskaznikZasadyInterlinii oddaje wskaźnik na zasadę liczenia interlinii, potrzebny przy polach nieobowiązkowych kontraktu.
func wejscieWskaznikZasadyInterlinii(
	wartosc shared.StudioLineSpacingRule) *shared.StudioLineSpacingRule {

	kopia := wartosc
	return &kopia
}

// ── Postać a treść: jedna prawda o położeniu w znakach ──────────────────────

// wejscieTrescZPostaci składa treść dokumentu z bloków postaci i jednocześnie przelicza zakresy znakowe bloków oraz fragmentów, licząc długość w runach, nie w bajtach, bo kontrakt mówi „w znakach”.
func wejscieTrescZPostaci(postac *shared.StudioDocumentForm) string {
	if postac == nil {
		return ""
	}
	var budowa strings.Builder
	polozenie := 0
	for i := range postac.Blocks {
		blok := &postac.Blocks[i]
		if i > 0 {
			budowa.WriteString("\n")
			polozenie++
		}
		blok.RangeStart = wejscieWskaznikCalkowity(polozenie)

		switch blok.Kind {
		case wejscieRodzajBlokuTabela:
			// Tabela wchodzi do treści jako wiersze rozdzielone tabulatorem, dla zgodności zaznaczenia z treścią.
			tekst := wejscieTekstTabeli(postac, blok.TableId)
			budowa.WriteString(tekst)
			polozenie += len([]rune(tekst))
		case wejscieRodzajBlokuObiekt:
			// Obiekt osadzony nie niesie znaków treści; zakres jest pusty i zakotwiczony w miejscu, gdzie stoi.
		case wejscieRodzajBlokuPodzial:
			// Podział jest cechą składu, nie treścią — znaków nie dokłada.
		default:
			for j := range blok.Runs {
				fragment := &blok.Runs[j]
				fragment.RangeStart = wejscieWskaznikCalkowity(polozenie)
				polozenie += len([]rune(fragment.Text))
				fragment.RangeEnd = wejscieWskaznikCalkowity(polozenie)
				budowa.WriteString(fragment.Text)
			}
		}
		blok.RangeEnd = wejscieWskaznikCalkowity(polozenie)
	}

	// Sekcje obejmują bloki, które do nich należą — granice biorą się z zakresów bloków, nie ze wskazania.
	wejsciePrzeliczSekcje(postac)
	return budowa.String()
}

// wejscieTekstTabeli składa treść tabeli: komórki rozdzielone tabulatorem,
// wiersze znakiem końca wiersza. Komórka wchłonięta scaleniem nie wnosi znaków.
func wejscieTekstTabeli(postac *shared.StudioDocumentForm, kodTabeli *string) string {
	if kodTabeli == nil {
		return ""
	}
	for i := range postac.Tables {
		tabela := &postac.Tables[i]
		if tabela.Id != *kodTabeli {
			continue
		}
		wiersze := make([]string, tabela.Rows)
		for w := 0; w < tabela.Rows; w++ {
			komorki := make([]string, 0, tabela.Columns)
			for k := 0; k < tabela.Columns; k++ {
				komorki = append(komorki, wejscieTekstKomorki(tabela, w, k))
			}
			wiersze[w] = strings.Join(komorki, "\t")
		}
		return strings.Join(wiersze, "\n")
	}
	return ""
}

// wejscieTekstKomorki oddaje treść komórki tabeli albo napis pusty, gdy komórka jest wchłonięta scaleniem.
func wejscieTekstKomorki(tabela *shared.StudioDocumentTable, wiersz, kolumna int) string {
	for i := range tabela.Cells {
		komorka := &tabela.Cells[i]
		if komorka.Row != wiersz || komorka.Column != kolumna {
			continue
		}
		if komorka.Merged != nil && *komorka.Merged {
			return ""
		}
		if komorka.Text != nil {
			return *komorka.Text
		}
		return ""
	}
	return ""
}

// wejsciePrzeliczSekcje ustawia granice sekcji na granicach bloków, które do
// nich należą. Sekcja bez ani jednego bloku zostaje z granicami zerowymi —
// istnieje, ale nie obejmuje treści, i tak ma wyjść kontraktem.
func wejsciePrzeliczSekcje(postac *shared.StudioDocumentForm) {
	if len(postac.Sections) == 0 {
		return
	}
	for i := range postac.Sections {
		sekcja := &postac.Sections[i]
		sekcja.Index = i
		poczatek, koniec, znaleziono := 0, 0, false
		for j := range postac.Blocks {
			blok := &postac.Blocks[j]
			if blok.SectionId == nil || *blok.SectionId != sekcja.Id {
				continue
			}
			if blok.RangeStart == nil || blok.RangeEnd == nil {
				continue
			}
			if !znaleziono {
				poczatek, koniec, znaleziono = *blok.RangeStart, *blok.RangeEnd, true
				continue
			}
			if *blok.RangeStart < poczatek {
				poczatek = *blok.RangeStart
			}
			if *blok.RangeEnd > koniec {
				koniec = *blok.RangeEnd
			}
		}
		sekcja.RangeStart, sekcja.RangeEnd = poczatek, koniec
	}
}

// wejscieUzycieStylow liczy, ile bloków używa każdego stylu nazwanego, i wpisuje
// wynik do arkusza. Liczba POLICZONA, nie zadeklarowana — galeria stylów ma
// pokazywać, ile miejsc przestawi zmiana stylu, i ta liczba musi być prawdziwa.
func wejscieUzycieStylow(postac *shared.StudioDocumentForm) {
	if postac == nil {
		return
	}
	licznik := make(map[string]int, len(postac.Styles))
	for i := range postac.Blocks {
		blok := &postac.Blocks[i]
		if blok.Paragraph != nil && blok.Paragraph.StyleName != nil {
			licznik[*blok.Paragraph.StyleName]++
		}
		for j := range blok.Runs {
			if blok.Runs[j].Format != nil && blok.Runs[j].Format.StyleName != nil {
				licznik[*blok.Runs[j].Format.StyleName]++
			}
		}
	}
	for i := range postac.Tables {
		if postac.Tables[i].StyleName != nil {
			licznik[*postac.Tables[i].StyleName]++
		}
	}
	for i := range postac.Styles {
		postac.Styles[i].UsageCount = wejscieWskaznikCalkowity(licznik[postac.Styles[i].Name])
	}
}

// ── Postać z treści płaskiej ────────────────────────────────────────────────

// wejsciePostacZTekstu składa postać dokumentu z treści płaskiej: każdy wiersz oddzielony pustym wierszem staje się akapitem tekstu zasadniczego, drogą awaryjną dla dokumentu bez własnej postaci.
func wejsciePostacZTekstu(kodDokumentu, tresc string) shared.StudioDocumentForm {
	postac := wejscieNowaPostac(kodDokumentu, "", nil)
	akapity := wejscieAkapityZTekstu(tresc)
	if len(akapity) == 0 {
		return postac
	}
	sekcja := postac.Sections[0].Id
	bloki := make([]shared.StudioDocumentBlock, 0, len(akapity))
	for _, akapit := range akapity {
		bloki = append(bloki, wejscieBlokAkapitu(sekcja, akapit, wejscieStylTekstZasadniczy, nil))
	}
	postac.Blocks = bloki
	return postac
}

// wejscieAkapityZTekstu rozbija treść płaską na akapity. Wiersz pusty jest
// granicą akapitu; pojedyncze zawinięcie wiersza wewnątrz akapitu zostaje
// spacją, bo w dokumencie edytowalnym łamanie wiersza liczy skład, nie plik.
func wejscieAkapityZTekstu(tresc string) []string {
	znormalizowana := strings.ReplaceAll(strings.ReplaceAll(tresc, "\r\n", "\n"), "\r", "\n")
	czesci := strings.Split(znormalizowana, "\n\n")
	akapity := make([]string, 0, len(czesci))
	for _, czesc := range czesci {
		zlozony := strings.TrimSpace(strings.Join(strings.Fields(czesc), " "))
		if zlozony == "" {
			continue
		}
		akapity = append(akapity, zlozony)
	}
	return akapity
}

// wejscieBlokAkapitu składa blok akapitu o wskazanym stylu i treści, gotowy do dołożenia do listy bloków postaci.
func wejscieBlokAkapitu(kodSekcji, tresc, styl string,
	postacZnaku *shared.StudioCharacterFormat) shared.StudioDocumentBlock {

	return shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuAkapit,
		SectionId: wejscieWskaznikTekstu(kodSekcji),
		Paragraph: &shared.StudioParagraphFormat{
			StyleName: wejscieWskaznikTekstu(styl),
		},
		Runs: []shared.StudioDocumentRun{{Text: tresc, Format: postacZnaku}},
	}
}

// wejscieStylNaglowkaPoziomu oddaje nazwę stylu nagłówka wskazanego poziomu.
// Poziom spoza zakresu jednego do sześciu schodzi na najbliższy istniejący —
// szósty jest najgłębszym poziomem arkusza i głębszego nie ma czym nazwać.
func wejscieStylNaglowkaPoziomu(poziom int) string {
	switch {
	case poziom <= 1:
		return wejscieStylNaglowek1
	case poziom == 2:
		return wejscieStylNaglowek2
	case poziom == 3:
		return wejscieStylNaglowek3
	case poziom == 4:
		return wejscieStylNaglowek4
	case poziom == 5:
		return wejscieStylNaglowek5
	}
	return wejscieStylNaglowek6
}
