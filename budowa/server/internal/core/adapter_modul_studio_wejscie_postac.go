// Odpowiedzialność pliku: jedno miejsce, w którym odcinek wejścia i wyjścia
// modułu Studio odkłada i czyta POSTAĆ dokumentu — arkusz stylów, nastawy
// strony, sekcje, bloki, tabele, obiekty, listy, aparat i pola.
//
// ── Dlaczego postać ma własny magazyn ───────────────────────────────────────
// `studio.document.save` do dziś przyjmował `documentId`, `content` i `title`,
// więc postać dokumentu ginęła przy każdym zapisie: model widział tekst, nie
// widział kroju ani tabeli. Kontrakt dostał pole `form`, a to jest droga, którą
// ta postać dojeżdża do rdzenia i wraca z niego nieuszkodzona.
//
// ── Gdzie postać leży — rozstrzygnięte, po wniesieniu warstwy danych ────────
// Wcześniejsza sesja tego odcinka nie miała jeszcze `dane/studio_postac_*.go`
// (migracje 361-362) i odkładała postać ładunkiem JSON w tabeli katalogowej
// modułu, pod zasięgiem własnym. Ta warstwa danych STOI, więc obie metody niżej
// zostały przełożone na nią — czyli na `postacWczytaj` i `postacZapisz` obszaru
// postaci (`adapter_modul_studio_postac.go`).
//
// Powód jest ten, którego zlecenie nie odpuszcza: dwa magazyny jednej postaci to
// dwie prawdy o tym samym dokumencie. Odcinek kontroli pracy woła
// `wejsciePostacDokumentu` przy zakładaniu kopii zapasowej — czytając dawny
// magazyn, dostawał postać PUSTĄ dla dokumentu, który postać ma, i kopia
// zapasowa niosłaby dokument bez formatowania. Po przełożeniu obie drogi czytają
// i piszą to samo miejsce.
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
	// wejscieZasiegPostaci trzyma wiersze postaci osobno od profili wydania.
	wejscieZasiegPostaci = "postac-dokumentu"
	// wejscieZasiegSzablonu trzyma postać wzorcową szablonu pisma.
	wejscieZasiegSzablonu = "postac-szablonu"
	// wejscieZasiegPochodzenia trzyma zapisy pochodzenia fragmentów dokumentu.
	wejscieZasiegPochodzenia = "pochodzenie-dokumentu"

	// Przedrostki kluczy wiersza. Klucz jest złożony z przedrostka i kodu
	// dokumentu albo szablonu, więc upsert trafia zawsze w ten sam wiersz.
	wejscieKluczPostaci     = "studio-postac-"
	wejscieKluczSzablonu    = "studio-szpostac-"
	wejscieKluczPochodzenia = "studio-poch-"

	// Przedrostki identyfikatorów bytów nadawanych przez ten odcinek.
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

// wejscieTeraz oddaje chwilę w milisekundach epoki.
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

// wejscieWskaznikCalkowity oddaje wskaźnik na liczbę całkowitą.
func wejscieWskaznikCalkowity(wartosc int) *int {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikDlugi oddaje wskaźnik na liczbę 64-bitową.
func wejscieWskaznikDlugi(wartosc int64) *int64 {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikRzeczywisty oddaje wskaźnik na liczbę rzeczywistą.
func wejscieWskaznikRzeczywisty(wartosc float64) *float64 {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikLogiczny oddaje wskaźnik na wartość logiczną.
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
		// Dokumentu, którego nie ma, nie udajemy postacią pustą: to zamieniłoby
		// brak dokumentu w dokument bez formatowania, a Operator uznałby, że
		// nigdy go nie miał.
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
			"rdzeń złożony bez repozytorium Studia — postaci dokumentu nie ma gdzie odłożyć")
	}
	dokument, err := a.repozytorium.Dokument(ctx, kodDokumentu)
	if err != nil {
		return shared.StudioDocumentForm{}, bladNieznanegoDokumentu(kodDokumentu, err)
	}
	postac.DocumentId = kodDokumentu

	// Numer porządkowy podbija BAZA przy zapisie (migracja 361), nie ten kod:
	// dwa zapisy z tym samym numerem byłyby możliwe, gdyby liczył go rdzeń,
	// a okno nie miałoby po czym poznać, że trzyma stan przestarzały.
	stan := &stanPostaci{dokument: dokument, forma: postac}
	stan.tekstPrzed = wartoscTekstu(dokument.Tresc)
	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioDocumentForm{}, err
	}
	return stan.forma, nil
}

// wejsciePrzyjmijPostacZapisu utrwala postać przyjechaną polem `form` komendy
// `studio.document.save` i oddaje wiersz dokumentu po zapisie.
//
// ── Dlaczego to stoi tutaj, a nie w obsłudze zapisu ─────────────────────────
// Zapis dokumentu przyjmował do niedawna `documentId`, `content` i `title`, więc
// postać przy każdym zapisie GINĘŁA: model widział tekst, nie widział kroju,
// wcięcia, tabeli ani obrazu. To jest ta jedna dziura, od której zaczęło się
// całe zlecenie. Droga utrwalenia jest tu ta sama, którą jedzie wniesienie
// pliku (`wejscieUtrwalPostac`) — dwie drogi zapisu postaci znaczyłyby dwie
// prawdy o postaci dokumentu.
//
// ── Co jest prawdą, gdy żądanie niesie i treść, i postać ────────────────────
// Treść z pola `content` jest prawdą o LITERACH: to jest to, co Operator ma
// w edytorze. Postać z pola `form` jest prawdą o STRUKTURZE: arkusz stylów,
// nastawy strony, sekcje, tabele, obiekty i aparat. Utrwalenie postaci składa
// treść z bloków (`postacTekstFormy`) i wpisuje ją do wiersza, więc wołający
// nadpisuje ją potem treścią z żądania — inaczej zapis cofałby litery dopisane
// w oknie do stanu, który zna drzewo postaci.
//
// Brak pola `form` NIE JEST tu obsługiwany: kontrakt mówi, że brak znaczy „bez
// zmiany postaci", więc wołający sprawdza obecność pola przed wywołaniem, a nie
// ta metoda po wywołaniu. Postać wyzerowana przy zwykłym zapisie treści byłaby
// tą samą szkodą, którą to pole ma naprawić.
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

// wejscieNowaPostac składa postać dokumentu pustego: arkusz stylów nazwanych,
// nastawy strony, jedną sekcję i jeden akapit pusty gotowy do pisania.
//
// Nowa strona bez akapitu byłaby stroną, na której nie ma gdzie postawić
// kursora — dlatego akapit pusty jest tu treścią, nie ozdobą.
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

// wejscieWskaznikPoczatkuSekcji oddaje wskaźnik na sposób rozpoczęcia sekcji.
func wejscieWskaznikPoczatkuSekcji(wartosc shared.StudioSectionStart) *shared.StudioSectionStart {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikWyrownania oddaje wskaźnik na wyrównanie tekstu.
func wejscieWskaznikWyrownania(wartosc shared.StudioTextAlign) *shared.StudioTextAlign {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikOrientacji oddaje wskaźnik na orientację strony.
func wejscieWskaznikOrientacji(wartosc shared.StudioPageOrientation) *shared.StudioPageOrientation {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikRodzajuNosnika oddaje wskaźnik na rodzaj nośnika.
func wejscieWskaznikRodzajuNosnika(wartosc shared.StudioPaperKind) *shared.StudioPaperKind {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikZrodlaObiektu oddaje wskaźnik na rodzaj obiektu.
func wejscieWskaznikZrodlaObiektu(wartosc shared.StudioObjectSource) *shared.StudioObjectSource {
	kopia := wartosc
	return &kopia
}

// wejscieWskaznikAutora oddaje wskaźnik na autora czynności.
func wejscieWskaznikAutora(wartosc shared.StudioAuthor) *shared.StudioAuthor {
	kopia := wartosc
	return &kopia
}

// wejscieAutorCzynnosci rozstrzyga autora czynności. Brak wskazania znaczy
// Operator — tak stanowi kontrakt każdej komendy tego odcinka. Wartość spoza
// wyliczenia jest pomyłką wołającego i nie schodzi cicho na Operatora: model,
// który podał autora przekręconego, nie ma wyjść z tego jako Operator, bo
// wtedy jego praca zniknęłaby z podświetlenia zmian modelu.
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

// wejscieDomyslneNastawyStrony składa nastawy strony nowego dokumentu.
//
// A4 i marginesy 25/25/25/25 mm są nastawą domyślną pisma urzędowego, nie
// upodobaniem wykonawcy. Nazwa nośnika podana przez Operatora ma pierwszeństwo
// i nie jest sprawdzana wykazem — wykaz nośników stoi dziś w pliku modułu
// Design, do którego temu odcinkowi wchodzić nie wolno (patrz sprawozdanie).
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

// wejscieDomyslnyArkuszStylow składa arkusz stylów nazwanych nowego dokumentu:
// tekst zasadniczy, sześć poziomów nagłówków, cytat, podpis i przypis.
//
// Nagłówki dziedziczą po tekście zasadniczym, a nie powtarzają jego postaci —
// to jest sens stylu nadrzędnego i to sprawia, że zmiana kroju w tekście
// zasadniczym przestawia cały dokument jednym ruchem.
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

	// Stopnie nagłówków maleją z poziomem — 18, 16, 14, 13, 12, 11 punktów.
	// Poziom konspektu równa się poziomowi nagłówka, bo spis treści zbiera się
	// po nim, a nie po nazwie stylu.
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

// wejscieWskaznikZasadyInterlinii oddaje wskaźnik na zasadę liczenia interlinii.
func wejscieWskaznikZasadyInterlinii(
	wartosc shared.StudioLineSpacingRule) *shared.StudioLineSpacingRule {

	kopia := wartosc
	return &kopia
}

// ── Postać a treść: jedna prawda o położeniu w znakach ──────────────────────

// wejscieTrescZPostaci składa treść dokumentu z bloków postaci i JEDNOCZEŚNIE
// przelicza zakresy znakowe bloków oraz fragmentów.
//
// Dwie czynności w jednym przebiegu, bo są jedną czynnością: położenie bloku
// w znakach ma sens wyłącznie wobec treści, którą ten przebieg właśnie składa.
// Liczone osobno rozjechałyby się przy pierwszym akapicie ze znakiem spoza
// zakresu jednobajtowego — dlatego długość mierzy się w RUNACH, a nie
// w bajtach: „ł" zajmuje dwa bajty i jeden znak, a kontrakt mówi „w znakach".
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
			// Tabela wchodzi do treści jako wiersze rozdzielone tabulatorem.
			// Bez tego zaznaczenie fragmentu obejmującego tabelę liczyłoby
			// znaki, których w treści nie ma, i przesunęłoby wszystkie
			// późniejsze zakresy.
			tekst := wejscieTekstTabeli(postac, blok.TableId)
			budowa.WriteString(tekst)
			polozenie += len([]rune(tekst))
		case wejscieRodzajBlokuObiekt:
			// Obiekt osadzony nie niesie znaków treści; jego zakres jest pusty
			// i zakotwiczony w miejscu, w którym stoi.
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

	// Sekcje obejmują bloki, które do nich należą — granice biorą się
	// z policzonych zakresów bloków, nie ze wskazania wołającego.
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

// wejscieTekstKomorki oddaje treść komórki tabeli albo napis pusty.
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

// wejsciePostacZTekstu składa postać dokumentu z treści płaskiej: każdy wiersz
// oddzielony pustym wierszem staje się akapitem tekstu zasadniczego.
//
// To jest droga plików bez własnej postaci (tekst czysty, warstwa tekstowa
// PDF-a) i droga awaryjna dla dokumentu, którego postaci rdzeń jeszcze nie zna.
// Nie udaje odczytu formatowania: oddaje akapity i arkusz stylów domyślny.
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

// wejscieBlokAkapitu składa blok akapitu o wskazanym stylu i treści.
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
