package core

import (
	"strings"
	"testing"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Sprawdziany znakowania: słownik rodzajów, autor znakowania i liczenie brzmienia w znakach.

// TestZnakowanieRodzajPozaSlownikiemOdmawia pilnuje słownika tabeli rodzajów znakowania zgodnego z kontraktem.
func TestZnakowanieRodzajPozaSlownikiemOdmawia(t *testing.T) {
	if err := znakowanieSprawdzRodzaj("wyroznienie"); err == nil {
		t.Errorf("rodzaj spoza wyliczenia kontraktu przeszedł — tabela odrzuci go " +
			"dopiero przy zapisie")
	} else if !strings.Contains(err.Error(), "highlight") {
		t.Errorf("odmowa nie mówi, co wolno: %v", err)
	}
	for _, rodzaj := range shared.WartosciStudioMarkupKind() {
		if err := znakowanieSprawdzRodzaj(rodzaj); err != nil {
			t.Errorf("rodzaj %q z wyliczenia kontraktu został odbity: %v", rodzaj, err)
		}
	}
}

// TestZnakowanieBrzmienieZastaneLiczySieWZnakach mierzy trzecią szkodę: wycinek
// zakresu idzie w znakach, nie w bajtach.
func TestZnakowanieBrzmienieZastaneLiczySieWZnakach(t *testing.T) {
	forma := shared.StudioDocumentForm{Blocks: []shared.StudioDocumentBlock{
		dziennikBlokDoSprawdzenia("blok", "Zażółć gęślą jaźń.", "Times"),
	}}
	postacPrzeliczZakresy(&forma)

	// Pierwsze sześć ZNAKÓW to „Zażółć"; po bajtach byłoby to „Zażó" z ogonkiem
	// rozciętym w środku.
	wycinek := znakowanieTekstZakresu(&forma, 0, 6)
	if wycinek != "Zażółć" {
		t.Errorf("wycinek zakresu policzony po bajtach, nie po znakach: %q", wycinek)
	}
	// Zakres spoza treści nie wywraca rachunku i nie zgaduje — oddaje pustkę.
	if wycinek := znakowanieTekstZakresu(&forma, 900, 950); wycinek != "" {
		t.Errorf("zakres spoza treści oddał brzmienie: %q", wycinek)
	}
}

// TestZnakowanieAutorIdzieZWiersza mierzy drugą szkodę: znakowanie modelu wychodzi
// kontraktem jako `model`, wraz z tożsamością konkretnego wykonawcy.
func TestZnakowanieAutorIdzieZWiersza(t *testing.T) {
	kod, nazwa := "agent-redaktor", "Redaktor pisma"
	brzmienie := "Proponowane brzmienie zdania."
	wiersz := dane.ZnakowanieStudia{
		Kod:                  "studio-znk-pierwsze",
		DokumentKod:          "studio-dok-pierwszy",
		Rodzaj:               string(shared.StudioMarkupKindSuggestion),
		AutorRodzaj:          string(shared.StudioAuthorModel),
		AutorAgentKod:        &kod,
		AutorAgentNazwa:      &nazwa,
		ZakresOd:             10,
		ZakresDo:             40,
		BrzmienieProponowane: &brzmienie,
		Stan:                 string(shared.StudioMarkupStateOpen),
	}
	zlozone := znakowanieZlozKontrakt(wiersz)
	if zlozone.Author != shared.StudioAuthorModel {
		t.Errorf("autor znakowania zgubił się w przekładzie: %q", zlozone.Author)
	}
	if zlozone.AuthorAgentId == nil || *zlozone.AuthorAgentId != kod {
		t.Errorf("tożsamość wykonawcy nie wyszła kontraktem — dwóch agentów pokaże " +
			"się jako jeden")
	}
	if zlozone.SuggestedText == nil || *zlozone.SuggestedText != brzmienie {
		t.Errorf("propozycja bez brzmienia jest komentarzem, nie propozycją")
	}
}

// TestZnakowanieOpisCzynnosciNazywaRece mierzy, że wpis dziennika mówi Operatorowi,
// CO cofa i CZYJA to była ręka — a nie każe odczytywać tego z rodzaju wpisu.
func TestZnakowanieOpisCzynnosciNazywaRece(t *testing.T) {
	kod := "agent-redaktor"
	wykonawca := kontrolaWykonawca{Rodzaj: shared.StudioAuthorModel, AgentKod: &kod}
	opis := znakowanieOpisCzynnosci(shared.StudioMarkupKindSuggestion, wykonawca, 10, 40)
	if !strings.Contains(opis, "propozycja") {
		t.Errorf("opis czynności nie nazywa czynności: %q", opis)
	}
	if !strings.Contains(opis, kod) {
		t.Errorf("opis czynności nie nazywa ręki: %q", opis)
	}

	opisOperatora := znakowanieOpisCzynnosci(shared.StudioMarkupKindHighlight,
		kontrolaWykonawca{Rodzaj: shared.StudioAuthorUzytkownik}, 0, 5)
	if !strings.Contains(opisOperatora, "Operator") {
		t.Errorf("czynność Operatora nie jest podpisana Operatorem: %q", opisOperatora)
	}
}

// TestZnakowanieRodzajFabrycznyNiesieLicznikUzyc mierzy, że wykaz rodzajów mówi
// prawdę o użyciu — bez tego Operator usunąłby rodzaj, na który wskazują
// znakowania.
func TestZnakowanieRodzajFabrycznyNiesieLicznikUzyc(t *testing.T) {
	barwa := "#ffcc00"
	zlozony := znakowanieZlozRodzaj(dane.RodzajZnacznikaStudia{
		Nazwa: "do sprawdzenia", NazwaWidoczna: "Do sprawdzenia",
		Barwa: &barwa, Fabryczny: true, IleUzyc: 7,
	})
	if zlozony.Builtin == nil || !*zlozony.Builtin {
		t.Errorf("fabryczność rodzaju zgubiła się w przekładzie")
	}
	if zlozony.UsageCount == nil || *zlozony.UsageCount != 7 {
		t.Errorf("licznik użyć nie wyszedł kontraktem: %+v", zlozony.UsageCount)
	}
	if zlozony.Color == nil || *zlozony.Color != barwa {
		t.Errorf("barwa rodzaju zgubiła się w przekładzie")
	}
}
