// Sprawdziany rachunku cofania z dziennika czynności wykluczają cofnięcie ze środka jako przywrócenie wersji, cofnięcie po cichu z zależnością i cofnięcie postaci, które nie rusza liter.
package core

import (
	"encoding/json"
	"strings"
	"testing"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// dziennikBlokDoSprawdzenia składa blok o wskazanym brzmieniu i kroju, gotowy do wpisania w dokument sprawdzianu.
func dziennikBlokDoSprawdzenia(kod, tekst, kroj string) shared.StudioDocumentBlock {
	return shared.StudioDocumentBlock{
		Id:   kod,
		Kind: "akapit",
		Runs: []shared.StudioDocumentRun{{
			Text:   tekst,
			Format: &shared.StudioCharacterFormat{FontFamily: &kroj},
		}},
	}
}

// TestDziennikCofnieciePostaciNieRuszaLiter jest sednem wymagania Właściciela:
// pomyłkowa zmiana kroju cofa się tak samo jak skasowany akapit, choć nie ruszyła
// ani jednej litery.
func TestDziennikCofnieciePostaciNieRuszaLiter(t *testing.T) {
	przed := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok-pierwszy", "Podstawa prawna zamówienia.", "Times"),
	}}
	po := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok-pierwszy", "Podstawa prawna zamówienia.", "Comic Sans"),
	}}
	biezaca := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok-pierwszy", "Podstawa prawna zamówienia.", "Comic Sans"),
	}}

	dziennikNalozBloki(&biezaca, &po, &przed)

	if postacTekstFormy(&biezaca) != "Podstawa prawna zamówienia." {
		t.Fatalf("cofnięcie zmiany kroju ruszyło treść: %q", postacTekstFormy(&biezaca))
	}
	kroj := biezaca.Blocks[0].Runs[0].Format.FontFamily
	if kroj == nil || *kroj != "Times" {
		t.Errorf("krój nie wrócił do brzmienia sprzed czynności: %+v", kroj)
	}
}

// TestDziennikCofnieciePozostawiaPozniejszaPrace mierzy różnicę między
// cofnięciem POJEDYNCZEJ czynności a przywróceniem wersji: akapit dopisany PO
// cofanej czynności ma zostać w dokumencie.
func TestDziennikCofnieciePozostawiaPozniejszaPrace(t *testing.T) {
	przed := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok-pierwszy", "Wstęp pisma.", "Times"),
	}}
	po := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok-pierwszy", "Wstęp pisma.", "Times"),
		dziennikBlokDoSprawdzenia("blok-modelu", "Akapit dopisany przez model.", "Times"),
	}}
	// Stan bieżący: po czynności modelu Operator dopisał jeszcze jeden akapit.
	biezaca := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok-pierwszy", "Wstęp pisma.", "Times"),
		dziennikBlokDoSprawdzenia("blok-modelu", "Akapit dopisany przez model.", "Times"),
		dziennikBlokDoSprawdzenia("blok-operatora", "Akapit dopisany przez Operatora.", "Times"),
	}}

	dziennikNalozBloki(&biezaca, &po, &przed)

	tresc := postacTekstFormy(&biezaca)
	if tresc == "" {
		t.Fatalf("cofnięcie zabrało całą treść dokumentu")
	}
	for _, blok := range biezaca.Blocks {
		if blok.Id == "blok-modelu" {
			t.Errorf("akapit wniesiony cofaną czynnością został w dokumencie: %q", tresc)
		}
	}
	zostal := false
	for _, blok := range biezaca.Blocks {
		if blok.Id == "blok-operatora" {
			zostal = true
		}
	}
	if !zostal {
		t.Errorf("cofnięcie zabrało pracę Operatora naniesioną PO cofanej czynności — "+
			"to jest przywrócenie wersji, nie cofnięcie pojedynczej czynności: %q", tresc)
	}
}

// TestDziennikCofnieciePrzywracaBlokZdjety mierzy drugi kierunek: blok, który
// cofana czynność USUNĘŁA, wraca na swoje miejsce.
func TestDziennikCofnieciePrzywracaBlokZdjety(t *testing.T) {
	przed := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok-pierwszy", "Wstęp pisma.", "Times"),
		dziennikBlokDoSprawdzenia("blok-usuniety", "Akapit skasowany przez model.", "Times"),
	}}
	po := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok-pierwszy", "Wstęp pisma.", "Times"),
	}}
	biezaca := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok-pierwszy", "Wstęp pisma.", "Times"),
	}}

	dziennikNalozBloki(&biezaca, &po, &przed)

	if len(biezaca.Blocks) != 2 {
		t.Fatalf("blok zdjęty czynnością nie wrócił: %d bloków", len(biezaca.Blocks))
	}
	if postacTekstBloku(biezaca.Blocks[1]) != "Akapit skasowany przez model." {
		t.Errorf("blok wrócił, ale nie w tym brzmieniu: %q",
			postacTekstBloku(biezaca.Blocks[1]))
	}
}

// TestDziennikZaleznoscZatrzymujeCofniecie mierzy, że czynność będąca podstawą
// późniejszej zostaje NAZWANA, a nie cofnięta po cichu.
func TestDziennikZaleznoscZatrzymujeCofniecie(t *testing.T) {
	czynnosc := dane.CzynnoscDokumentuStudia{
		Kod:           "studio-czn-podstawa",
		Opis:          "wstawienie tabeli kosztów",
		StojaceNaNiej: []string{"studio-czn-wiersze"},
	}
	stan := map[string]string{
		"studio-czn-podstawa": string(shared.StudioActionStateActive),
		"studio-czn-wiersze":  string(shared.StudioActionStateActive),
	}

	przeszkoda := dziennikStojaceNaNiej(czynnosc, nil, stan, map[string]bool{})
	if przeszkoda != "studio-czn-wiersze" {
		t.Fatalf("zależność nie została rozpoznana: %q", przeszkoda)
	}
	pozycja := dziennikPominiecieZaleznosci(czynnosc, przeszkoda)
	if pozycja.Detail == nil {
		t.Fatalf("odmowa nie nazywa zależności")
	}
	if !dziennikZawieraNapis(*pozycja.Detail, "studio-czn-wiersze") {
		t.Errorf("treść odmowy nie mówi, KTÓRA czynność stoi na cofanej: %q", *pozycja.Detail)
	}

	// Ta sama zależność cofana RAZEM z podstawą przeszkodą już nie jest.
	razem := map[string]bool{"studio-czn-wiersze": true}
	if przeszkoda := dziennikStojaceNaNiej(czynnosc, nil, stan, razem); przeszkoda != "" {
		t.Errorf("czynność cofana razem z podstawą została uznana za przeszkodę: %q", przeszkoda)
	}
	// Zależność już cofnięta też przeszkodą nie jest.
	stan["studio-czn-wiersze"] = string(shared.StudioActionStateReverted)
	if przeszkoda := dziennikStojaceNaNiej(czynnosc, nil, stan, map[string]bool{}); przeszkoda != "" {
		t.Errorf("czynność już cofnięta została uznana za przeszkodę: %q", przeszkoda)
	}
}

// TestDziennikPonowienieWymagaStojacejPodstawy mierzy drugi kierunek zależności:
// czynności nie da się ponowić, dopóki cofnięta stoi ta, na której ona stoi.
func TestDziennikPonowienieWymagaStojacejPodstawy(t *testing.T) {
	czynnosc := dane.CzynnoscDokumentuStudia{
		Kod:          "studio-czn-wiersze",
		Opis:         "wstawienie wierszy tabeli",
		PodstawyKody: []string{"studio-czn-podstawa"},
	}
	stan := map[string]string{
		"studio-czn-podstawa": string(shared.StudioActionStateReverted),
	}
	if przeszkoda := dziennikPodstawaCofnieta(czynnosc, stan, map[string]bool{}); przeszkoda !=
		"studio-czn-podstawa" {

		t.Fatalf("cofnięta podstawa nie zatrzymała ponowienia: %q", przeszkoda)
	}
	stan["studio-czn-podstawa"] = string(shared.StudioActionStateActive)
	if przeszkoda := dziennikPodstawaCofnieta(czynnosc, stan, map[string]bool{}); przeszkoda != "" {
		t.Errorf("stojąca podstawa została uznana za przeszkodę: %q", przeszkoda)
	}
}

// TestDziennikStanWpisuJestSlownikiemTabeli pilnuje pomyłki, która raz już przeszła kompilację: stan wpisu to wyłącznie active albo reverted, a rodzaj czynności to StudioActionKind, nie StudioChangeKind.
func TestDziennikStanWpisuJestSlownikiemTabeli(t *testing.T) {
	stany := shared.WartosciStudioActionState()
	if len(stany) != 2 {
		t.Fatalf("stanów wpisu dziennika ma być dwa, jest %d: %+v", len(stany), stany)
	}
	if string(stany[0]) != "active" || string(stany[1]) != "reverted" {
		t.Errorf("słownik stanów wpisu rozjechał się z migracją 364: %+v", stany)
	}
	// Słownik tabeli poszerza migracja 378 — liczba rodzajów tutaj i warunek CHECK tam muszą się zgadzać.
	if len(shared.WartosciStudioActionKind()) != 12 {
		t.Errorf("rodzajów czynności dziennika ma być dwanaście, jest %d",
			len(shared.WartosciStudioActionKind()))
	}
	if len(shared.WartosciStudioChangeKind()) != 3 {
		t.Errorf("rodzajów zmiany śledzonej ma być trzy, jest %d",
			len(shared.WartosciStudioChangeKind()))
	}
	// Żaden rodzaj zmiany śledzonej nie może być zarazem rodzajem czynności dziennika.
	for _, zmiana := range shared.WartosciStudioChangeKind() {
		for _, czynnosc := range shared.WartosciStudioActionKind() {
			if string(zmiana) == string(czynnosc) {
				t.Errorf("słowniki zachodzą na siebie wartością %q — pomyłka przejdzie "+
					"kompilację i wywróci się na ograniczeniu tabeli", zmiana)
			}
		}
	}
}

// TestDziennikDrzewoNieczytelneNieJestBrakiem pilnuje, żeby uszkodzony ładunek
// wpisu nie zamieniał się w ciszę: Operator uznałby wtedy, że czynność nie
// miała czego cofać.
func TestDziennikDrzewoNieczytelneNieJestBrakiem(t *testing.T) {
	uszkodzony := "{to nie jest JSON"
	if _, err := dziennikDrzewo(&uszkodzony); err == nil {
		t.Errorf("nieczytelne drzewo postaci przeszło jako brak")
	}
	pusty := "   "
	drzewo, err := dziennikDrzewo(&pusty)
	if err != nil {
		t.Errorf("pusty ładunek nie jest uszkodzeniem, a odmówił: %v", err)
	}
	if drzewo != nil {
		t.Errorf("pusty ładunek dał drzewo zamiast braku")
	}
	poprawny := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok", "Treść.", "Times"),
	}}
	zapis, err := json.Marshal(poprawny)
	if err != nil {
		t.Fatalf("nie da się złożyć ładunku sprawdzianu: %v", err)
	}
	tekst := string(zapis)
	odczytane, err := dziennikDrzewo(&tekst)
	if err != nil || odczytane == nil {
		t.Fatalf("poprawne drzewo nie odczytało się: %v", err)
	}
	if postacTekstFormy(odczytane) != "Treść." {
		t.Errorf("drzewo odczytało się w innym brzmieniu: %q", postacTekstFormy(odczytane))
	}
}

// dziennikZawieraNapis mówi, czy napis niesie podnapis. Własny pomocnik
// z przedrostkiem odcinka, bo sprawdziany tego odcinka pytają o TREŚĆ odmowy,
// a nie o sam fakt odmowy.
func dziennikZawieraNapis(napis, podnapis string) bool {
	return strings.Contains(napis, podnapis)
}
