package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// pomijBezProgramu pomija sprawdzian, nazywając program, którego na maszynie
// zabrakło, żeby wynik pominięcia był czytelny w dzienniku.
func pomijBezProgramu(t *testing.T, nazwa, program string) {
	t.Helper()
	if _, err := exec.LookPath(program); err != nil {
		t.Skipf("pomiar niewykonany: na tej maszynie nie ma programu %s (%s) — %v",
			nazwa, program, err)
	}
}

// pomijBezWydaniaJavy pomija sprawdzian, gdy na maszynie nie ma Javy albo
// brakuje wydania (archiwum) programu wymaganego do pracy.
func pomijBezWydaniaJavy(t *testing.T, nazwa, archiwum string) {
	t.Helper()
	pomijBezProgramu(t, nazwa, "java")
	if archiwum == "" {
		t.Skipf("pomiar niewykonany: na tej maszynie nie ma wydania %s", nazwa)
	}
}

// TestSkladPdfOddajeDokumentDoOdczytania wykazuje drogę składu: markdown
// wychodzi PDF-em, z którego da się odczytać zdanie, które do niego weszło.
// Odczyt idzie drugą komendą i innym programem (poppler), więc mierzy plik,
// a nie własną pamięć.
func TestSkladPdfOddajeDokumentDoOdczytania(t *testing.T) {
	pomijBezProgramu(t, narzedziePandoc.Nazwa, narzedziePandoc.Program)
	pomijBezProgramu(t, narzedzieTypst.Nazwa, narzedzieTypst.Program)
	pomijBezProgramu(t, narzedziePdfDoTekstu.Nazwa, narzedziePdfDoTekstu.Program)
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const zdanie = "Zażółć gęślą jaźń — protokół odbioru terenu."
	var wynik shared.DocumentConvertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDocumentConvert,
		shared.DocumentConvertRequest{
			Content:    wskaznik("# Protokół\n\n" + zdanie + "\n"),
			FromFormat: wskaznik("markdown"),
			ToFormat:   "pdf",
		}, &wynik)

	if wynik.SizeBytes == 0 {
		t.Fatal("skład oddał dokument o zerowym rozmiarze")
	}
	if wynik.Asset.Uri == nil || *wynik.Asset.Uri == "" {
		t.Fatal("wynik składu nie niesie odwołania do bajtów")
	}
	bajty, err := os.ReadFile(*wynik.Asset.Uri)
	if err != nil {
		t.Fatalf("pod odwołaniem wyniku nie ma pliku: %v", err)
	}
	if !strings.HasPrefix(string(bajty), "%PDF") {
		t.Fatalf("wynik nie jest dokumentem PDF — pierwsze bajty: %q", pierwszeBajty(bajty))
	}
	// Nazwa silnika stoi w metadanych dokumentu (pole `Creator`) — dowód, że
	// złożył go silnik składu.
	if !strings.Contains(strings.ToLower(string(bajty)), narzedzieTypst.Program) {
		t.Fatal("dokument nie niesie w metadanych nazwy silnika składu — " +
			"PDF powstał inną drogą niż typst")
	}

	// Odczyt idzie ścieżką bloku, nie identyfikatorem: żądanie bez `windowId`
	// nie zakłada wiersza zasobu.
	var odczyt shared.DocumentTextExtractResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{SourcePath: wynik.Asset.Uri}, &odczyt)
	if !strings.Contains(bezZlamanWiersza(odczyt.Text), bezZlamanWiersza(zdanie)) {
		t.Fatalf("ze złożonego dokumentu nie da się odczytać zdania, które do niego weszło\n"+
			" oczekiwano fragmentu: %q\n odczytano: %q", zdanie, odczyt.Text)
	}
	if odczyt.UsedOcr {
		t.Fatal("dokument złożony ma warstwę tekstową, a odczyt poszedł rozpoznaniem pisma")
	}
}

// TestOdczytSiegaFormatuSpozaSlownika wykazuje drogę Tiki: wiadomość poczty nie
// stoi w słowniku formatów rdzenia i przed tą drogą wracała odmową, a jej treść
// jest tekstem, który model ma prawo przeczytać.
func TestOdczytSiegaFormatuSpozaSlownika(t *testing.T) {
	pomijBezWydaniaJavy(t, narzedzieTiki.Nazwa, sciezkaArchiwumTikiDoPomiaru())
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const tresc = "Zażółć gęślą jaźń w treści wiadomości."
	sciezka := filepath.Join(t.TempDir(), "wiadomosc.eml")
	wiadomosc := "From: nadawca@example.org\r\n" +
		"To: odbiorca@example.org\r\n" +
		"Subject: Protokół odbioru\r\n" +
		"Date: Mon, 01 Sep 2025 10:00:00 +0200\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" + tresc + "\r\n"
	if err := os.WriteFile(sciezka, []byte(wiadomosc), 0o600); err != nil {
		t.Fatalf("nie można założyć materiału sprawdzianu: %v", err)
	}

	// Format `eml` nie stoi w słowniku rdzenia — gdyby stanął, sprawdzian
	// mierzyłby inną drogę.
	if _, stoi := formatyDokumentu["eml"]; stoi {
		t.Fatal("format eml wszedł do słownika rdzenia — sprawdzian mierzy już inną drogę")
	}

	var odczyt shared.DocumentTextExtractResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{SourcePath: wskaznik(sciezka)}, &odczyt)

	if !strings.Contains(odczyt.Text, tresc) {
		t.Fatalf("odczyt nie oddał treści wiadomości\n oczekiwano fragmentu: %q\n odczytano: %q",
			tresc, odczyt.Text)
	}
	if odczyt.UsedOcr {
		t.Fatal("treść wiadomości jest tekstem, a odczyt zgłosił rozpoznanie pisma")
	}
}

// TestOdczytPozaSlownikiemNazywaBrakWydania wykazuje odmowę: gdy Java stoi, a
// wydania Tiki nie ma, odmowa ma nazwać brakujące ARCHIWUM i drogę naprawy —
// bo naprawą nie jest tu instalacja Javy.
func TestOdczytPozaSlownikiemNazywaBrakWydania(t *testing.T) {
	pomijBezProgramu(t, narzedzieTiki.Nazwa, narzedzieTiki.Program)
	t.Setenv(zmiennaTiki, t.TempDir())
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	sciezka := filepath.Join(t.TempDir(), "wiadomosc.eml")
	if err := os.WriteFile(sciezka, []byte("Subject: x\r\n\r\ntreść\r\n"), 0o600); err != nil {
		t.Fatalf("nie można założyć materiału sprawdzianu: %v", err)
	}

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{SourcePath: wskaznik(sciezka)})

	if odmowa.Code != shared.ErrorCodeChannelUnavailable {
		t.Fatalf("brak wydania programu ma nieść kod channel_unavailable, a niesie %s: %s",
			odmowa.Code, odmowa.Message)
	}
	if !strings.Contains(odmowa.Message, "Tiki") || !strings.Contains(odmowa.Message, zmiennaTiki) {
		t.Fatalf("odmowa nie nazywa brakującego wydania ani drogi wskazania: %q", odmowa.Message)
	}
}

// TestKorektaSiegaSlownikaJezyka wykazuje drogę LanguageToola: literówka, której
// żadna reguła napisowa nie widzi, wraca ustaleniem rodzaju „ortografia" wraz
// z propozycją, a propozycja realnie poprawia treść panelu w bazie.
func TestKorektaSiegaSlownikaJezyka(t *testing.T) {
	pomijBezWydaniaJavy(t, narzedzieLanguageToola.Nazwa, archiwumKorekty())
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-slownika", "Zdanie źródłowe.")
	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-slownika", "pl-PL",
		"Protokół zawiera jeden bledem zapisany wyraz.")

	var korekta shared.TranslateProofreadRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateProofreadRun,
		shared.TranslateProofreadRunRequest{
			PanelId: kodPanelu,
			Checks:  []shared.ProofreadCheckKind{shared.ProofreadCheckKindSpelling},
		}, &korekta)

	ustalenie, jest := ustalenieORodzaju(korekta.Findings, shared.ProofreadCheckKindSpelling)
	if !jest {
		t.Fatalf("korekta nie zgłosiła ustalenia ortograficznego wobec wyrazu spoza słownika; "+
			"ustalenia: %s", opisUstalen(korekta.Findings))
	}
	if ustalenie.Severity != shared.ProofreadSeverityError {
		t.Errorf("wyraz spoza słownika ma wagę %q, a jest faktem o słowie, nie opinią reguły",
			ustalenie.Severity)
	}
	if ustalenie.Suggestion == nil || *ustalenie.Suggestion == "" {
		t.Fatal("ustalenie ortograficzne nie niesie propozycji poprawki")
	}
	if strings.Contains(*ustalenie.Suggestion, "bledem") {
		t.Fatalf("propozycja poprawki nadal niesie wyraz spoza słownika: %q", *ustalenie.Suggestion)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateProofreadApply,
		shared.TranslateProofreadApplyRequest{
			PanelId: kodPanelu, FindingIds: []string{ustalenie.Id},
		}, nil)

	wBazie := napisZBazy(t, db,
		`SELECT tresc FROM panel_tlumaczenia WHERE identyfikator_zewnetrzny = ?`, kodPanelu)
	if strings.Contains(wBazie, "bledem") {
		t.Fatalf("po zastosowaniu propozycji treść panelu w bazie nadal niesie literówkę: %q", wBazie)
	}
}

// TestKorektaPisowniSchodziNaSlownikSystemu wykazuje drogę zapasową: bez
// wydania LanguageToola pisownię prowadzi hunspell, a ustalenie nazywa
// słownik, którym mierzono. Nieobecność LanguageToola jest tu wywołana
// pustym katalogiem w zmiennej.
func TestKorektaPisowniSchodziNaSlownikSystemu(t *testing.T) {
	pomijBezProgramu(t, narzedzieHunspella.Nazwa, narzedzieHunspella.Program)
	t.Setenv(zmiennaLanguageToola, t.TempDir())
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	if archiwumKorekty() != "" {
		t.Fatalf("sprawdzian miał wywołać nieobecność LanguageToola, a archiwum wciąż stoi: %s",
			archiwumKorekty())
	}

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-zapasowe", "Zdanie źródłowe.")
	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-zapasowe", "pl-PL",
		"Protokół zawiera jeden bledem zapisany wyraz.")

	var korekta shared.TranslateProofreadRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateProofreadRun,
		shared.TranslateProofreadRunRequest{
			PanelId: kodPanelu,
			Checks:  []shared.ProofreadCheckKind{shared.ProofreadCheckKindSpelling},
		}, &korekta)

	ustalenie, jest := ustalenieORodzaju(korekta.Findings, shared.ProofreadCheckKindSpelling)
	if !jest {
		t.Fatalf("droga zapasowa nie zgłosiła ustalenia ortograficznego; ustalenia: %s",
			opisUstalen(korekta.Findings))
	}
	if !strings.Contains(ustalenie.Detail, "słowniku") {
		t.Fatalf("ustalenie nie nazywa słownika, którym mierzono: %q", ustalenie.Detail)
	}
}

// TestKorektaProzyZglaszaPowtorzenie wykazuje drogę vale: powtórzenie wyrazu
// przez granicę zdania, którego reguła wbudowana nie widzi, wraca ustaleniem
// rodzaju „styl". Materiał stawia powtórzenie w jednym wierszu bez kropki
// między wyrazami.
func TestKorektaProzyZglaszaPowtorzenie(t *testing.T) {
	pomijBezProgramu(t, narzedzieVale.Nazwa, narzedzieVale.Program)
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-prozy", "Zdanie źródłowe.")
	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-prozy", "nieznany-tej-maszynie",
		"Raport opisuje stan stan terenu.")

	var korekta shared.TranslateProofreadRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateProofreadRun,
		shared.TranslateProofreadRunRequest{
			PanelId: kodPanelu,
			Checks:  []shared.ProofreadCheckKind{shared.ProofreadCheckKindStyle},
		}, &korekta)

	ustalenie, jest := ustalenieORodzaju(korekta.Findings, shared.ProofreadCheckKindStyle)
	if !jest {
		t.Fatalf("analiza prozy nie zgłosiła powtórzenia wyrazu; ustalenia: %s",
			opisUstalen(korekta.Findings))
	}
	if !strings.Contains(strings.ToLower(ustalenie.Detail), "vale.") {
		t.Fatalf("ustalenie nie nazywa reguły, która je zgłosiła: %q", ustalenie.Detail)
	}
}

// TestObrobkaWstepnaProstujeSkosPrzedRozpoznaniem wykazuje drogę unpapera:
// nastawy pozycji przestały być zapisem bez skutku — skan pochylony przechodzi
// przez prostowanie, a rozpoznanie oddaje słowa materiału.
func TestObrobkaWstepnaProstujeSkosPrzedRozpoznaniem(t *testing.T) {
	pomijBezProgramu(t, narzedzieCzyszczeniaSkanu.Nazwa, narzedzieCzyszczeniaSkanu.Program)
	pomijBezProgramu(t, narzedzieRozpoznaniaStudia.Nazwa, narzedzieRozpoznaniaStudia.Program)
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const tresc = "PROTOKOL ODBIORU"
	sciezka := skanPochylony(t, tresc)
	pozycja := dolozPozycje(t, zmontowany, zycie, "okno-obrobki", sciezka)

	// Przebieg pierwszy bez obróbki jest odniesieniem — bez niego drugi
	// przebieg nie dowiódłby zmiany.
	var bezObrobki shared.StudioIngestRecognizeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestRecognize,
		shared.StudioIngestRecognizeRequest{
			ItemId:   pozycja.Id,
			Settings: &shared.StudioRecognitionSettings{Languages: []string{"pol"}},
		}, &bezObrobki)
	if bezObrobki.Item.Text != nil && strings.Contains(
		strings.ToUpper(bezZlamanWiersza(*bezObrobki.Item.Text)), "PROTOKOL") {
		t.Skip("pomiar bezprzedmiotowy: rozpoznanie czyta materiał pochylony bez obróbki, " +
			"więc ten sprawdzian nie odróżniłby drogi z unpaperem od drogi bez niego")
	}

	var rozpoznanie shared.StudioIngestRecognizeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestRecognize,
		shared.StudioIngestRecognizeRequest{
			ItemId: pozycja.Id,
			Settings: &shared.StudioRecognitionSettings{
				Languages: []string{"pol"},
				Deskew:    wskaznik(true),
			},
		}, &rozpoznanie)

	if rozpoznanie.Item.Text == nil || strings.TrimSpace(*rozpoznanie.Item.Text) == "" {
		t.Fatal("rozpoznanie po obróbce wstępnej nie oddało ani jednego znaku")
	}
	odczytane := strings.ToUpper(bezZlamanWiersza(*rozpoznanie.Item.Text))
	if !strings.Contains(odczytane, "PROTOKOL") {
		t.Fatalf("tekst rozpoznany po obróbce nie niesie słowa z materiału\n"+
			" materiał: %q\n odczytano: %q", tresc, *rozpoznanie.Item.Text)
	}
	if len(rozpoznanie.Words) == 0 {
		t.Fatal("rozpoznanie nie oddało warstwy słów — korekta rozpoznania nie ma czego poprawiać")
	}
}

// TestWyciagnijTekstProstujeSkosPrzedRozpoznaniem wykazuje drogę unpapera
// z komendy document.text.extract: pole preprocess kontraktu przestało być
// zapisem bez skutku — skan pochylony przechodzi przez prostowanie, a
// rozpoznanie oddaje słowo materiału.
func TestWyciagnijTekstProstujeSkosPrzedRozpoznaniem(t *testing.T) {
	pomijBezProgramu(t, narzedzieCzyszczeniaSkanu.Nazwa, narzedzieCzyszczeniaSkanu.Program)
	pomijBezProgramu(t, narzedzieTesseract.Nazwa, narzedzieTesseract.Program)
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const tresc = "PROTOKOL ODBIORU"
	sciezka := skanPochylony(t, tresc)

	// Przebieg pierwszy bez obróbki jest odniesieniem — bez niego drugi
	// przebieg nie dowiódłby zmiany. Rozpoznanie materiału pochylonego może
	// tu odmówić pustym odczytem albo oddać tekst bez szukanego słowa; oba
	// wyniki są odniesieniem, tylko odczyt SŁOWA czyniłby sprawdzian bezprzedmiotowym.
	odpowiedzOdniesienia := wykonajKomende(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{SourcePath: wskaznik(sciezka), Language: wskaznik("pol")})
	if odpowiedzOdniesienia.Error == nil {
		var bezObrobki shared.DocumentTextExtractResponse
		if err := protocol.LadunekDo(odpowiedzOdniesienia, &bezObrobki); err != nil {
			t.Fatalf("nieczytelny ładunek odpowiedzi odniesienia: %v", err)
		}
		if strings.Contains(strings.ToUpper(bezZlamanWiersza(bezObrobki.Text)), "PROTOKOL") {
			t.Skip("pomiar bezprzedmiotowy: rozpoznanie czyta materiał pochylony bez obróbki, " +
				"więc ten sprawdzian nie odróżniłby drogi z unpaperem od drogi bez niego")
		}
	}

	var poObrobce shared.DocumentTextExtractResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{
			SourcePath: wskaznik(sciezka), Language: wskaznik("pol"), Preprocess: wskaznik(true),
		}, &poObrobce)

	odczytane := strings.ToUpper(bezZlamanWiersza(poObrobce.Text))
	if !strings.Contains(odczytane, "PROTOKOL") {
		t.Fatalf("tekst rozpoznany po obróbce nie niesie słowa z materiału\n"+
			" materiał: %q\n po obróbce: %q", tresc, poObrobce.Text)
	}
}

// TestWyciagnijTekstRozpoznajeDwaJezykiZlozoneZnakiemPlus wykazuje drogę
// produkcyjną wielojęzycznego rozpoznania: dwa języki skrótami ISO 639-1
// („pl", „en") mają dać Tesseractowi wykaz, który ten rozumie („pol+eng"),
// nie surowy człon, którego dane językowe tej maszyny nie niosą.
func TestWyciagnijTekstRozpoznajeDwaJezykiZlozoneZnakiemPlus(t *testing.T) {
	pomijBezProgramu(t, narzedzieTesseract.Nazwa, narzedzieTesseract.Program)
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	sciezka := kartkaTekstu(t, "PROTOKOL ODBIORU")
	var odczyt shared.DocumentTextExtractResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{SourcePath: wskaznik(sciezka), Language: wskaznik("pl+en")},
		&odczyt)

	odczytane := strings.ToUpper(bezZlamanWiersza(odczyt.Text))
	if !strings.Contains(odczytane, "PROTOKOL") {
		t.Fatalf("rozpoznanie dwoma językami (pl+en) nie oddało słowa z materiału\n odczytano: %q",
			odczyt.Text)
	}
}

// TestWyciagnijTekstOdmawiaJezykaNieniesionegoPrzezTesseracta wykazuje
// odmowę nazwaną, gdy żądanie wskaże język, którego danych tej maszyny nie
// niosą — cichej próby rozpoznania w innym języku niż zamówiony rdzeń nie
// dopuszcza.
func TestWyciagnijTekstOdmawiaJezykaNieniesionegoPrzezTesseracta(t *testing.T) {
	pomijBezProgramu(t, narzedzieTesseract.Nazwa, narzedzieTesseract.Program)
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const jezykNieniesiony = "xx-jezyk-ktorego-nie-ma"
	sciezka := kartkaTekstu(t, "PROTOKOL ODBIORU")
	blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{
			SourcePath: wskaznik(sciezka), Language: wskaznik(jezykNieniesiony),
		})
	if blad.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("odmowa języka nieniesionego niesie kod %q, oczekiwano %q",
			blad.Code, shared.ErrorCodeValidationFailed)
	}
	if !strings.Contains(blad.Message, jezykNieniesiony) {
		t.Fatalf("odmowa nie nazywa języka, którego maszyna nie niesie: %q", blad.Message)
	}
}

// TestWyciagnijTekstOdmawiaObrobkiWstepnejBezUnpapera pilnuje, żeby brak
// programu na maszynie dał odmowę nazwaną, nie cichy odczyt bez obróbki.
func TestWyciagnijTekstOdmawiaObrobkiWstepnejBezUnpapera(t *testing.T) {
	if zewnetrzne.Stoi(narzedzieCzyszczeniaSkanu) {
		t.Skip("pomiar niewykonany: unpaper jest na tej maszynie, sprawdzian mierzy jego brak")
	}
	pomijBezProgramu(t, "ImageMagick", "magick")
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	sciezka := skanPochylony(t, "MATERIAL")
	blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{
			SourcePath: wskaznik(sciezka), Preprocess: wskaznik(true),
		})
	if blad.Code != shared.ErrorCodeChannelUnavailable {
		t.Fatalf("odmowa braku unpapera niesie kod %q, oczekiwano %q",
			blad.Code, shared.ErrorCodeChannelUnavailable)
	}
	if !strings.Contains(strings.ToLower(blad.Message), "unpaper") {
		t.Fatalf("treść odmowy nie nazywa brakującego programu: %q", blad.Message)
	}
}

// TestWyciagnijTekstOdmawiaObrobkiPrzedRasteryzacjaPdf wykazuje, że odmowa
// braku unpapera na drodze PDF zapada PRZED rasteryzacją stron: gdy brakuje
// obu programów, odmowa nazywa unpaper, a nie poppler — dowód, że rasteryzacja
// w ogóle nie ruszyła, tak jak Studio pyta o program przed pracą.
func TestWyciagnijTekstOdmawiaObrobkiPrzedRasteryzacjaPdf(t *testing.T) {
	pomijBezProgramu(t, "ImageMagick", "magick")

	zastaneCzyszczenie := narzedzieCzyszczeniaSkanu
	narzedzieCzyszczeniaSkanu = zewnetrzne.Narzedzie{
		Nazwa: zastaneCzyszczenie.Nazwa, Program: "danaco-unpaper-ktorego-nie-ma",
		Pakiet: zastaneCzyszczenie.Pakiet,
	}
	t.Cleanup(func() { narzedzieCzyszczeniaSkanu = zastaneCzyszczenie })

	zastanyPoppler := narzedziePdfDoObrazu
	narzedziePdfDoObrazu = zewnetrzne.Narzedzie{
		Nazwa: zastanyPoppler.Nazwa, Program: "danaco-pdftoppm-ktorego-nie-ma",
		Pakiet: zastanyPoppler.Pakiet,
	}
	t.Cleanup(func() { narzedziePdfDoObrazu = zastanyPoppler })

	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	sciezka := pdfObrazowy(t, kartkaTekstu(t, "MATERIAL"))

	blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{
			SourcePath: wskaznik(sciezka), Preprocess: wskaznik(true),
		})
	if blad.Code != shared.ErrorCodeChannelUnavailable {
		t.Fatalf("odmowa braku unpapera niesie kod %q, oczekiwano %q",
			blad.Code, shared.ErrorCodeChannelUnavailable)
	}
	if !strings.Contains(strings.ToLower(blad.Message), "unpaper") {
		t.Fatalf("odmowa nie nazywa unpapera — rasteryzacja PDF-u pobiegła przed sprawdzeniem "+
			"jego obecności: %q", blad.Message)
	}
}

// TestObrobkaWstepnaNieRuszaBezNastaw pilnuje, żeby program nie startował tam,
// gdzie nikt o niego nie prosił: pozycja bez nastaw obróbki idzie do rozpoznania
// wprost.
func TestObrobkaWstepnaNieRuszaBezNastaw(t *testing.T) {
	if czyszczenieZadane(nastawyRozpoznania{}) {
		t.Fatal("nastawy puste uznano za żądanie obróbki wstępnej — program ruszałby bez powodu")
	}
	for nazwa, nastawy := range map[string]nastawyRozpoznania{
		"prostowanie": {Prostowanie: true},
		"odszumianie": {Odszumianie: true},
		"progowanie":  {Progowanie: true},
		"marginesy":   {PrzycinanieMarginesow: true},
		"sam układ":   {Uklad: true},
		"sam język":   {Jezyki: []string{"pol"}},
	} {
		zadane := czyszczenieZadane(nastawy)
		oczekiwane := nazwa != "sam układ" && nazwa != "sam język"
		if zadane != oczekiwane {
			t.Errorf("nastawa %q daje żądanie obróbki %v, oczekiwano %v", nazwa, zadane, oczekiwane)
		}
	}
}

// TestArgumentyObrobkiWylaczajaFiltryNiezamowione pilnuje, że unpaper ma
// filtry włączone domyślnie, a kontrakt wyłączone, więc filtr niezamówiony
// musi zostać wyłączony jawnie, inaczej `deskew` włączałoby po cichu obróbkę,
// o którą nikt nie prosił.
func TestArgumentyObrobkiWylaczajaFiltryNiezamowione(t *testing.T) {
	argumenty := strings.Join(
		argumentyCzyszczenia(nastawyRozpoznania{Prostowanie: true}, "we.ppm", "wy.ppm"), " ")

	for _, wylaczony := range []string{
		"--no-noisefilter", "--no-blurfilter", "--no-grayfilter",
		"--no-border-scan", "--no-border-align", "--no-blackfilter",
	} {
		if !strings.Contains(argumenty, wylaczony) {
			t.Errorf("filtr niezamówiony nie został wyłączony (%s): %s", wylaczony, argumenty)
		}
	}
	if strings.Contains(argumenty, "--no-deskew") {
		t.Errorf("prostowanie zamówione nastawą zostało wyłączone: %s", argumenty)
	}
	if strings.Contains(argumenty, "--type") {
		t.Errorf("progowanie niezamówione nastawą weszło do wiersza wywołania: %s", argumenty)
	}
}

// TestWykazZaleznosciZnaProgramyTresciPisanej pilnuje, żeby sześć programów tej
// rodziny stało w wykazie zależności — inaczej sonda startowa milczałaby o ich
// braku, a Operator dowiadywałby się o nim dopiero z funkcji, która odmawia.
func TestWykazZaleznosciZnaProgramyTresciPisanej(t *testing.T) {
	wykaz := zaleznosciZewnetrzne()
	if len(wykaz) == 0 {
		t.Fatal("wykaz zależności jest pusty — nie ma czego pilnować")
	}
	zakresy := map[string][]string{}
	for _, pozycja := range wykaz {
		nazwa := pozycja.Narzedzie.Nazwa
		zakresy[nazwa] = append(zakresy[nazwa], pozycja.Zakres)
	}

	for _, narzedzie := range []struct {
		nazwa  string
		zakres string
	}{
		{narzedzieTypst.Nazwa, "skład dokumentu do PDF-u"},
		{narzedzieTiki.Nazwa, "odczyt formatu spoza słownika rdzenia"},
		{narzedzieLanguageToola.Nazwa, "korekta językowa"},
		{narzedzieHunspella.Nazwa, "ortografia drogą zapasową"},
		{narzedzieVale.Nazwa, "styl prozy"},
		{narzedzieCzyszczeniaSkanu.Nazwa, "obróbka wstępna skanu"},
	} {
		wpisy, stoi := zakresy[narzedzie.nazwa]
		if !stoi {
			t.Errorf("wykaz zależności nie zna programu %s (%s); ten obszar nie ma "+
				"biblioteki czysto-Go, więc program jest składnikiem pakietu serwera",
				narzedzie.nazwa, narzedzie.zakres)
			continue
		}
		if len(wpisy) != 1 {
			t.Errorf("program %s stoi w wykazie %d razy — dwa wiersze o jednym zakresie "+
				"rozjadą się przy pierwszej zmianie", narzedzie.nazwa, len(wpisy))
		}
		if strings.TrimSpace(wpisy[0]) == "" {
			t.Errorf("pozycja %s nie mówi, co przestaje działać przy jej braku", narzedzie.nazwa)
		}
	}
}

// TestJezykKorektyNieZmyslaOdmianyKrajowej pilnuje trzystopniowego
// rozstrzygnięcia: oznaczenie idzie bez zmiany, nazwa własna sprowadza się
// do języka podstawowego, a napis nierozpoznany nie idzie do programu wcale.
func TestJezykKorektyNieZmyslaOdmianyKrajowej(t *testing.T) {
	for wejscie, oczekiwane := range map[string]string{
		"pl-PL":                   "pl-PL",
		"en-GB":                   "en-GB",
		"de":                      "de",
		"de-DE-x-simple-language": "de-DE-x-simple-language",
		"polski":                  "pl",
		"Angielski":               "en",
		"english":                 "en",
		"":                        "",
		"język, którego nie ma":   "",
		"francuski":               "",
	} {
		if wynik := jezykKorekty(wejscie); wynik != oczekiwane {
			t.Errorf("język %q sprowadzono do %q, oczekiwano %q", wejscie, wynik, oczekiwane)
		}
	}
}

// ── Pomocnicy sprawdzianów ──────────────────────────────────────────────────

// skanPochylony rysuje kartkę o znanej treści i pochyla ją o dwa stopnie —
// skos, którego oko prawie nie widzi, a Tesseract nie czyta wcale, więc
// sprawdzian mierzy skutek dwoma przebiegami, z obróbką i bez niej.
func skanPochylony(t *testing.T, tresc string) string {
	t.Helper()

	rysownik, err := exec.LookPath("magick")
	if err != nil {
		t.Skipf("pomiar niewykonany: brak programu magick — nie ma czym narysować materiału: %v", err)
	}
	sciezka := filepath.Join(t.TempDir(), "skan.png")
	polecenie := exec.Command(rysownik,
		"-background", "white", "-fill", "black", "-pointsize", "72", "-density", "300",
		"label:"+tresc, "-bordercolor", "white", "-border", "80",
		"-rotate", "2", "-background", "white", "-flatten", sciezka)
	if wyjscie, err := polecenie.CombinedOutput(); err != nil {
		t.Skipf("pomiar niewykonany: nie udało się narysować materiału: %v (%s)", err, wyjscie)
	}
	if opis, err := os.Stat(sciezka); err != nil || opis.Size() == 0 {
		t.Skip("pomiar niewykonany: materiał sprawdzianu nie powstał albo jest pusty")
	}
	return sciezka
}

// kartkaTekstu rysuje kartkę WPROST, bez pochylenia — materiał do sprawdzianów,
// którym chodzi o samo rozpoznanie, nie o obróbkę wstępną skosu.
func kartkaTekstu(t *testing.T, tresc string) string {
	t.Helper()

	rysownik, err := exec.LookPath("magick")
	if err != nil {
		t.Skipf("pomiar niewykonany: brak programu magick — nie ma czym narysować materiału: %v", err)
	}
	sciezka := filepath.Join(t.TempDir(), "kartka.png")
	polecenie := exec.Command(rysownik,
		"-background", "white", "-fill", "black", "-pointsize", "72", "-density", "300",
		"label:"+tresc, "-bordercolor", "white", "-border", "80", sciezka)
	if wyjscie, err := polecenie.CombinedOutput(); err != nil {
		t.Skipf("pomiar niewykonany: nie udało się narysować materiału: %v (%s)", err, wyjscie)
	}
	if opis, err := os.Stat(sciezka); err != nil || opis.Size() == 0 {
		t.Skip("pomiar niewykonany: materiał sprawdzianu nie powstał albo jest pusty")
	}
	return sciezka
}

// pdfObrazowy zawija obraz w jednostronicowy PDF bez warstwy tekstowej —
// materiał, który document.text.extract musi odczytać rozpoznaniem pisma,
// nie warstwą zapisanych znaków.
func pdfObrazowy(t *testing.T, sciezkaObrazu string) string {
	t.Helper()

	rysownik, err := exec.LookPath("magick")
	if err != nil {
		t.Skipf("pomiar niewykonany: brak programu magick — nie ma czym złożyć PDF-u: %v", err)
	}
	sciezka := filepath.Join(t.TempDir(), "obrazowy.pdf")
	polecenie := exec.Command(rysownik, sciezkaObrazu, sciezka)
	if wyjscie, err := polecenie.CombinedOutput(); err != nil {
		t.Skipf("pomiar niewykonany: nie udało się złożyć PDF-u: %v (%s)", err, wyjscie)
	}
	if opis, err := os.Stat(sciezka); err != nil || opis.Size() == 0 {
		t.Skip("pomiar niewykonany: PDF sprawdzianu nie powstał albo jest pusty")
	}
	return sciezka
}

// sciezkaArchiwumTikiDoPomiaru oddaje ścieżkę klas Tiki albo pustkę — pomocnik
// pomijania, żeby sprawdzian nie powtarzał składania ścieżki.
func sciezkaArchiwumTikiDoPomiaru() string {
	sciezka, jest := sciezkaKlasTiki()
	if !jest {
		return ""
	}
	return sciezka
}

// ustalenieORodzaju wybiera pierwsze ustalenie wskazanego rodzaju z wykazu
// ustaleń, które sprawdzian dostał od programu korekty.
func ustalenieORodzaju(ustalenia []shared.ProofreadFinding,
	rodzaj shared.ProofreadCheckKind) (shared.ProofreadFinding, bool) {

	for _, ustalenie := range ustalenia {
		if ustalenie.Kind == rodzaj {
			return ustalenie, true
		}
	}
	return shared.ProofreadFinding{}, false
}

// opisUstalen składa czytelny wykaz ustaleń do treści niepowodzenia — bez niego
// sprawdzian mówiłby „nie znalazłem" i nie mówił, co znalazł.
func opisUstalen(ustalenia []shared.ProofreadFinding) string {
	if len(ustalenia) == 0 {
		return "brak"
	}
	czesci := make([]string, 0, len(ustalenia))
	for _, ustalenie := range ustalenia {
		czesci = append(czesci, string(ustalenie.Kind)+": "+ustalenie.Detail)
	}
	return strings.Join(czesci, " | ")
}

// bezZlamanWiersza sprowadza tekst do jednego wiersza z pojedynczymi odstępami.
// Skład i rozpoznanie pisma łamią wiersze tam, gdzie im wypadnie, więc
// porównanie znak w znak mierzyłoby łamanie, a nie treść.
func bezZlamanWiersza(tekst string) string {
	return strings.Join(strings.Fields(tekst), " ")
}

// pierwszeBajty oddaje początek treści do komunikatu o niepowodzeniu, żeby
// dziennik sprawdzianu nie niósł całej długiej treści materiału.
func pierwszeBajty(bajty []byte) string {
	if len(bajty) > 16 {
		bajty = bajty[:16]
	}
	return string(bajty)
}
