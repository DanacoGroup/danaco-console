// Plik zawiera sprawdziany obszaru list, mierzące rachunek poziomów i wcięć, na którym stoi
// wygląd punktu na kartce.
package core

import (
	"strings"
	"testing"

	"danacoconsole/shared"
)

// TestListyPoziomyPowstajaZNastawami mierzy, że poziom listy nie jest pustą pozycją bez własnej nastawy.
func TestListyPoziomyPowstajaZNastawami(t *testing.T) {
	wypunktowanie := shared.StudioListDefinition{
		Id: "studio-lis-jedna", Kind: shared.StudioListKindBullet,
	}
	listyZapewnijPoziomy(&wypunktowanie, 3)
	if len(wypunktowanie.Levels) != 3 {
		t.Fatalf("zapewnienie poziomu trzeciego dało %d poziomów", len(wypunktowanie.Levels))
	}
	for numer, poziom := range wypunktowanie.Levels {
		if poziom.Level != numer+1 {
			t.Errorf("poziom na miejscu %d ma numer %d", numer, poziom.Level)
		}
		if poziom.BulletCharacter == nil || *poziom.BulletCharacter == "" {
			t.Errorf("poziom %d wypunktowania nie ma znaku", poziom.Level)
		}
		if poziom.IndentMm == nil {
			t.Errorf("poziom %d nie ma wcięcia", poziom.Level)
		}
	}
	// Wcięcie rośnie z poziomem — bez tego lista wielopoziomowa wyglądałaby jak
	// jednopoziomowa.
	if *wypunktowanie.Levels[0].IndentMm >= *wypunktowanie.Levels[2].IndentMm {
		t.Error("wcięcie nie rośnie wraz z poziomem listy")
	}
	// Znaki poziomów różnią się — inaczej Operator nie odróżniłby zagnieżdżenia.
	if *wypunktowanie.Levels[0].BulletCharacter == *wypunktowanie.Levels[1].BulletCharacter {
		t.Error("dwa pierwsze poziomy wypunktowania mają ten sam znak")
	}

	// Zapewnienie powtórne nie mnoży poziomów.
	listyZapewnijPoziomy(&wypunktowanie, 2)
	if len(wypunktowanie.Levels) != 3 {
		t.Errorf("zapewnienie powtórne dało %d poziomów", len(wypunktowanie.Levels))
	}
}

// TestListyWzorNumeracjiPrawniczej mierzy wymaganie wprost wymienione przez
// Właściciela: numeracja prawnicza wielopoziomowa typu 1.1.2.
func TestListyWzorNumeracjiPrawniczej(t *testing.T) {
	lista := shared.StudioListDefinition{
		Id: "studio-lis-prawnicza", Kind: shared.StudioListKindMultilevel,
	}
	listyZapewnijPoziomy(&lista, 3)
	trzeci := listyPoziom(&lista, 3)
	if trzeci == nil || trzeci.Pattern == nil {
		t.Fatal("poziom trzeci listy wielopoziomowej nie ma wzoru numeru")
	}
	if *trzeci.Pattern != "%1.%2.%3." {
		t.Errorf("wzór poziomu trzeciego brzmi %q, a ma składać się z numerów "+
			"wszystkich poziomów nadrzędnych", *trzeci.Pattern)
	}
	pierwszy := listyPoziom(&lista, 1)
	if pierwszy == nil || pierwszy.Pattern == nil || *pierwszy.Pattern != "%1." {
		t.Error("wzór poziomu pierwszego nie jest samym numerem tego poziomu")
	}
	if pierwszy.StartAt == nil || *pierwszy.StartAt != 1 {
		t.Error("poziom numerowany nie ma punktu startu")
	}
}

// TestListyWciecieSchodziDoAkapitu mierzy przejście wcięcia poziomu na postać
// akapitu — bez niego linijka nie ma czym pokazać znaczników.
func TestListyWciecieSchodziDoAkapitu(t *testing.T) {
	poziom := shared.StudioListLevel{
		Level:     2,
		IndentMm:  postacWskaznikMiary(12.7),
		HangingMm: postacWskaznikMiary(6.35),
	}
	zmiana := listyWciecieAkapitu(&poziom)
	if zmiana.IndentLeftMm == nil || *zmiana.IndentLeftMm != 12.7 {
		t.Error("wcięcie poziomu nie zeszło do wcięcia lewego akapitu")
	}
	// Wysunięcie pierwszego wiersza jest ujemne — na tym stoi znak punktu wiszący z lewej strony.
	if zmiana.FirstLineIndentMm == nil || *zmiana.FirstLineIndentMm != -6.35 {
		t.Errorf("wysunięcie pierwszego wiersza wynosi %v, a ma być ujemne",
			zmiana.FirstLineIndentMm)
	}
	// Poziom niepodany nie zmienia niczego — czynność poza listą nie rusza wcięcia w ciszy.
	if pusta := listyWciecieAkapitu(nil); pusta.IndentLeftMm != nil {
		t.Error("brak poziomu ustawił wcięcie akapitu")
	}
}

// TestListyRodzajIFormatSprawdzajaKontrakt pilnuje, żeby wartość spoza kontraktu
// nie przeszła w ciszy, a odmowa nazywała wykaz.
func TestListyRodzajIFormatSprawdzajaKontrakt(t *testing.T) {
	if err := listyRodzajZnany(shared.StudioListKindBullet); err != nil {
		t.Errorf("rodzaj z kontraktu odmówił: %v", err)
	}
	err := listyRodzajZnany(shared.StudioListKind("kropkowana"))
	if err == nil {
		t.Fatal("rodzaj spoza kontraktu przeszedł bez odmowy")
	}
	if !strings.Contains(err.Error(), string(shared.StudioListKindMultilevel)) {
		t.Errorf("odmowa nie nazywa wykazu rodzajów: %v", err)
	}
	if err := listyFormatZnany(shared.StudioListNumberFormat("gwiazdki")); err == nil {
		t.Error("format numeracji spoza kontraktu przeszedł bez odmowy")
	}
	if err := listyFormatZnany(shared.StudioListNumberFormatLegal); err != nil {
		t.Errorf("format prawniczy odmówił: %v", err)
	}
	if err := listyZrodloZnane(shared.StudioBulletSource("naklejka")); err == nil {
		t.Error("źródło znaku wypunktowania spoza kontraktu przeszło bez odmowy")
	}
	for _, zrodlo := range shared.WartosciStudioBulletSource() {
		if err := listyZrodloZnane(zrodlo); err != nil {
			t.Errorf("źródło %q z kontraktu odmówiło: %v", zrodlo, err)
		}
	}
}

// TestListyBlokiIZakresListy mierzy rachunek, z którego bierze się zakres wpisu
// dziennika: zmiana nastaw listy musi wiedzieć, których znaków dotknęła.
func TestListyBlokiIZakresListy(t *testing.T) {
	forma := shared.StudioDocumentForm{
		Blocks: []shared.StudioDocumentBlock{
			{Id: "studio-blok-poza", Kind: blokPostaciAkapit,
				Runs: []shared.StudioDocumentRun{{Text: "Akapit poza listą"}}},
			{Id: "studio-blok-punkt", Kind: blokPostaciAkapit,
				Paragraph: &shared.StudioParagraphFormat{
					ListId:    postacWskaznikTekstu("studio-lis-jedna"),
					ListLevel: postacWskaznikLiczby(1),
				},
				Runs: []shared.StudioDocumentRun{{Text: "Punkt pierwszy"}}},
			{Id: "studio-blok-punkt-dwa", Kind: blokPostaciAkapit,
				Paragraph: &shared.StudioParagraphFormat{
					ListId:    postacWskaznikTekstu("studio-lis-jedna"),
					ListLevel: postacWskaznikLiczby(2),
				},
				Runs: []shared.StudioDocumentRun{{Text: "Punkt drugi"}}},
		},
	}
	postacPrzeliczZakresy(&forma)

	wskazania := listyBlokiListy(&forma, "studio-lis-jedna")
	if len(wskazania) != 2 {
		t.Fatalf("lista ma %d punktów, a ma mieć dwa", len(wskazania))
	}
	od, do := listyZakresListy(&forma, wskazania)
	if od != *forma.Blocks[1].RangeStart || do != *forma.Blocks[2].RangeEnd {
		t.Errorf("zakres listy wyszedł %d-%d, a punkty stoją %d-%d",
			od, do, *forma.Blocks[1].RangeStart, *forma.Blocks[2].RangeEnd)
	}
	// Akapit poza listą nie wchodzi do zakresu, inaczej wpis dziennika mówiłby o zmianie, której nie było.
	if od == 0 {
		t.Error("zakres listy objął akapit stojący poza listą")
	}
	if len(listyBlokiListy(&forma, "studio-lis-nieistniejaca")) != 0 {
		t.Error("lista nieistniejąca ma punkty")
	}
}

// TestListyDefinicjaZadaniaOdmawiaNazwane pilnuje, żeby wskazanie listy, której
// dokument nie ma, było odmową nazywającą wykaz — nie cichym założeniem drugiej.
func TestListyDefinicjaZadaniaOdmawiaNazwane(t *testing.T) {
	forma := shared.StudioDocumentForm{}
	_, err := listyDefinicjaZadania(&forma, "studio-lis-widmo", "dokument-1")
	if err == nil {
		t.Fatal("wskazanie listy w dokumencie bez list przeszło bez odmowy")
	}
	if !strings.Contains(err.Error(), "studio.list.apply") {
		t.Errorf("odmowa nie wskazuje drogi wyjścia: %v", err)
	}

	forma.Lists = []shared.StudioListDefinition{{Id: "studio-lis-jedna"}}
	_, err = listyDefinicjaZadania(&forma, "studio-lis-widmo", "dokument-1")
	if err == nil || !strings.Contains(err.Error(), "studio-lis-jedna") {
		t.Errorf("odmowa nie nazywa list, które dokument ma: %v", err)
	}
	if _, err := listyDefinicjaZadania(&forma, "  ", "dokument-1"); err == nil {
		t.Error("czynność na liście bez wskazania listy przeszła bez odmowy")
	}
	definicja, err := listyDefinicjaZadania(&forma, "studio-lis-jedna", "dokument-1")
	if err != nil || definicja == nil || definicja.Id != "studio-lis-jedna" {
		t.Errorf("lista istniejąca nie wróciła: %+v, błąd %v", definicja, err)
	}
}
