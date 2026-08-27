// Plik sprawdza nastawy widoku, różnicę postaci i schowek: nastawę spoza
// wyliczenia kontraktu, milczącą zmianę kroju bez zmiany liter i kolejność
// różnicy niezależną od przebiegu mapy.
package core

import (
	"strings"
	"testing"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// TestWidokNastawaPozaSlownikiemOdmawia mierzy nastawę widoku spoza wyliczenia
// kontraktu, którą tabela odrzuciłaby dopiero przy zapisie, po fakcie.
func TestWidokNastawaPozaSlownikiemOdmawia(t *testing.T) {
	dozwolone := widokNazwy(shared.WartosciStudioSurfaceMode())
	if err := widokSprawdzWartosc("tryb powierzchni", "kolumny", dozwolone); err == nil {
		t.Errorf("tryb powierzchni spoza wyliczenia przeszedł")
	} else if !strings.Contains(err.Error(), "tabs") {
		t.Errorf("odmowa nie mówi, co wolno: %v", err)
	}
	for _, tryb := range shared.WartosciStudioSurfaceMode() {
		if err := widokSprawdzWartosc("tryb powierzchni", string(tryb), dozwolone); err != nil {
			t.Errorf("tryb %q z wyliczenia został odbity: %v", tryb, err)
		}
	}
}

// TestWidokNastawyWychodzaWCalosci mierzy, że wybór Operatora jest PAMIĘTANY:
// każda nastawa wychodzi kontraktem, nie tylko ta ostatnio zmieniona.
func TestWidokNastawyWychodzaWCalosci(t *testing.T) {
	nastawy := widokZlozNastawy(dane.NastawaWidokuStudia{
		Okno:                     "okno-pierwsze",
		TrybPowierzchni:          shared.StudioSurfaceModeSplit,
		KierunekPodzialu:         shared.StudioSplitOrientationHorizontal,
		GranicaPodzialu:          0.35,
		TrybWidoku:               shared.StudioViewModePrintPreview,
		SkalaProcent:             140,
		SkalaNastawa:             shared.StudioZoomPresetCustom,
		LinijkiWidoczne:          true,
		LinijkaJednostka:         shared.StudioRulerUnitInch,
		GranicaMarginesu:         false,
		ZnakiFormatowania:        true,
		StronWRzedzie:            2,
		WidokRozkladowki:         true,
		Przewijanie:              shared.StudioScrollModePage,
		PodswietlenieZmianModelu: true,
		PrzybornikWidoczny:       false,
	}, dane.DokumentStudia{Kod: "studio-dok-pierwszy"})

	if nastawy.SurfaceMode == nil || *nastawy.SurfaceMode != shared.StudioSurfaceModeSplit {
		t.Errorf("tryb powierzchni nie wyszedł kontraktem: %+v", nastawy.SurfaceMode)
	}
	if nastawy.SplitRatio == nil || *nastawy.SplitRatio != 0.35 {
		t.Errorf("położenie granicy podziału nie wyszło kontraktem: %+v", nastawy.SplitRatio)
	}
	if nastawy.ZoomPercent == nil || *nastawy.ZoomPercent != 140 {
		t.Errorf("skala widoku nie wyszła kontraktem — nie jest pamiętana przy dokumencie")
	}
	if nastawy.PagesPerRow == nil || *nastawy.PagesPerRow != 2 {
		t.Errorf("widok dwóch stron obok siebie nie wyszedł kontraktem")
	}
	if nastawy.ModelChangesHighlighted == nil || !*nastawy.ModelChangesHighlighted {
		t.Errorf("przełącznik podświetlenia zmian modelu nie wyszedł kontraktem")
	}
	if nastawy.ToolboxVisible == nil || *nastawy.ToolboxVisible {
		t.Errorf("przybornik zwinięty przez Operatora wyszedł jako rozwinięty")
	}
	if nastawy.DocumentId == nil || *nastawy.DocumentId != "studio-dok-pierwszy" {
		t.Errorf("nastawa nie mówi, przy którym dokumencie stoi")
	}
	if nastawy.WindowId == nil || *nastawy.WindowId != "okno-pierwsze" {
		t.Errorf("nastawa nie mówi, przy którym oknie stoi")
	}
}

// TestWidokRoznicaPostaciNieMilczyOKroju mierzy, że różnica dwóch wersji
// dokumentu liczy się także na postaci, nie tylko na literach, a więc nazywa
// zmianę kroju jako zmianę.
func TestWidokRoznicaPostaciNieMilczyOKroju(t *testing.T) {
	odniesienie := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok", "Podstawa prawna zamówienia.", "Times"),
	}}
	porownywana := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok", "Podstawa prawna zamówienia.", "Comic Sans"),
	}}

	// Rachunek treści milczy o zmianie kroju, dlatego różnica postaci jest osobna.
	fragmenty := policzFragmentyRoznicy(postacTekstFormy(&odniesienie),
		postacTekstFormy(&porownywana))
	for _, fragment := range fragmenty {
		if fragment.Kind != shared.DiffHunkKindContext {
			t.Fatalf("rachunek treści zobaczył zmianę kroju — sprawdzian mierzy nie to, " +
				"co miał mierzyć")
		}
	}

	wpisy := roznicaZlozWpisy(&odniesienie, &porownywana, "")
	if len(wpisy) == 0 {
		t.Fatalf("różnica postaci przemilczała zmianę kroju")
	}
	nazwany := false
	for _, wpis := range wpisy {
		if wpis.Area == roznicaObszarPostaciZnaku {
			nazwany = true
			if wpis.Before == nil || wpis.After == nil {
				t.Errorf("wpis różnicy postaci nie niesie brzmienia „przed” i „po”")
			}
		}
	}
	if !nazwany {
		t.Errorf("różnica nie nazwała obszaru postaci znaku: %+v", wpisy)
	}
}

// TestWidokRoznicaPostaciWidziNastawyStrony mierzy, że zmiana nośnika albo
// orientacji strony też nie milczy, lecz wychodzi jako osobny wpis różnicy.
func TestWidokRoznicaPostaciWidziNastawyStrony(t *testing.T) {
	pionowa := shared.StudioPageOrientation(shared.StudioPageOrientationPionowa)
	pozioma := shared.StudioPageOrientation(shared.StudioPageOrientationPozioma)
	odniesienie := shared.StudioDocumentForm{
		PageSetup: &shared.StudioPageSetup{Orientation: &pionowa},
	}
	porownywana := shared.StudioDocumentForm{
		PageSetup: &shared.StudioPageSetup{Orientation: &pozioma},
	}
	wpisy := roznicaZlozWpisy(&odniesienie, &porownywana, "")
	if len(wpisy) != 1 || wpisy[0].Area != roznicaObszarNastawStrony {
		t.Fatalf("zmiana orientacji nie wyszła różnicą postaci: %+v", wpisy)
	}
	if wpisy[0].Kind != roznicaRodzajZmienione {
		t.Errorf("zmiana nastaw strony wyszła jako %q", wpisy[0].Kind)
	}

	// Zawężenie do obszaru działa: pytanie o styl nazwany nie oddaje nastaw strony.
	if wpisy := roznicaZlozWpisy(&odniesienie, &porownywana, roznicaObszarStylu); len(wpisy) != 0 {
		t.Errorf("zawężenie do obszaru stylu oddało wpisy innego obszaru: %+v", wpisy)
	}
}

// TestWidokRoznicaMaUstalonaKolejnosc mierzy trzecią szkodę: ten sam dokument
// oglądany dwa razy daje ten sam wykaz.
func TestWidokRoznicaMaUstalonaKolejnosc(t *testing.T) {
	pierwsza := map[string]int{"trzeci": 3, "pierwszy": 1, "drugi": 2}
	druga := map[string]int{"czwarty": 4}
	wzorzec := roznicaNazwyScalone(pierwsza, druga)
	for przebieg := 0; przebieg < 20; przebieg++ {
		kolejny := roznicaNazwyScalone(pierwsza, druga)
		if len(kolejny) != len(wzorzec) {
			t.Fatalf("wykaz zmienił długość między przebiegami")
		}
		for i := range wzorzec {
			if kolejny[i] != wzorzec[i] {
				t.Fatalf("kolejność wykazu różnicy zmienia się między przebiegami: %v vs %v",
					wzorzec, kolejny)
			}
		}
	}
}

// TestWidokZakresRozbieznosciLiczyWZnakach pilnuje, że zakres przeniesionego
// fragmentu liczy się w znakach: po bajtach rozciąłby polskie litery dwubajtowe.
func TestWidokZakresRozbieznosciLiczyWZnakach(t *testing.T) {
	przed := "Zażółć gęślą jaźń."
	po := "Zażółć gęślą duszę."
	od, do := roznicaZakresRozbieznosci(przed, po)
	znaki := []rune(przed)
	if od < 0 || do > len(znaki) || od > do {
		t.Fatalf("zakres rozbieżności wypadł poza treścią: od %d do %d przy %d znakach",
			od, do, len(znaki))
	}
	if od < 7 {
		t.Errorf("zakres rozbieżności zaczyna się przed pierwszą różnicą — policzono "+
			"go po bajtach: od %d", od)
	}
}

// TestSchowekSposobWklejeniaPozaSlownikiemOdmawia pilnuje, że jawny wybór
// Operatora jest wyborem z wykazu, a nie dowolnym napisem.
func TestSchowekSposobWklejeniaPozaSlownikiemOdmawia(t *testing.T) {
	if err := schowekSprawdzSposob("bezFormatu"); err == nil {
		t.Errorf("sposób wklejenia spoza wyliczenia przeszedł")
	} else if !strings.Contains(err.Error(), "plainText") {
		t.Errorf("odmowa nie mówi, co wolno: %v", err)
	}
	for _, sposob := range shared.WartosciStudioPasteMode() {
		if err := schowekSprawdzSposob(sposob); err != nil {
			t.Errorf("sposób %q z wyliczenia został odbity: %v", sposob, err)
		}
	}
}

// TestSchowekPochodzenieNiesieZrodlo mierzy, że fragment wniesiony ze schowka
// niesie zapis, skąd jest — podstawa pod bibliografię i pod „podobieństwa".
func TestSchowekPochodzenieNiesieZrodlo(t *testing.T) {
	wpis := "schowek-pierwszy"
	zapis := schowekZlozPochodzenie("studio-dok-pierwszy", dane.PochodzenieFragmentuStudia{
		Kod: "studio-poch-pierwsze", Rodzaj: shared.StudioProvenanceKindClipboard,
		ZakresOd: 10, ZakresDo: 40, WersjaZrodla: &wpis,
		AutorRodzaj: string(shared.StudioAuthorModel),
	})
	if zapis.Kind != shared.StudioProvenanceKindClipboard {
		t.Errorf("rodzaj źródła zgubił się w przekładzie: %q", zapis.Kind)
	}
	if zapis.SourceVersion == nil || *zapis.SourceVersion != wpis {
		t.Errorf("zapis pochodzenia nie mówi, z którego wpisu schowka fragment przyszedł")
	}
	if zapis.Author == nil || *zapis.Author != shared.StudioAuthorModel {
		t.Errorf("zapis pochodzenia nie mówi, kto fragment wniósł")
	}
	if zapis.DocumentId != "studio-dok-pierwszy" {
		t.Errorf("zapis pochodzenia zgubił dokument")
	}
}
