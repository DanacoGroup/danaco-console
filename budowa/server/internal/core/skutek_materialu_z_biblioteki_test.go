package core

import (
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Sprawdzian mierzy skutek: assetId sięga obu magazynów rdzenia, Design i Library.

// TestWydobycieTekstuSiegaPlikuBiblioteki wykazuje drogę, której wcześniej nie
// było: identyfikator pliku biblioteki jako materiał komendy `document.text.extract`.
func TestWydobycieTekstuSiegaPlikuBiblioteki(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	tresc := "materiał wniesiony do biblioteki, nie do zasobów designu"
	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "material.txt", tresc)

	var wynik shared.DocumentTextExtractResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{AssetId: wskaznik(plik.Id)}, &wynik)

	if strings.TrimSpace(wynik.Text) != tresc {
		t.Fatalf("wydobyty tekst nie jest treścią pliku biblioteki\n oczekiwano: %q\n otrzymano:  %q",
			tresc, wynik.Text)
	}
	if wynik.UsedOcr {
		t.Fatalf("plik tekstowy nie wymaga rozpoznania pisma, a narzędzie je zgłosiło")
	}
}

// TestWydobycieTekstuOdmawiaKodowiSpozaObuMagazynow pilnuje, żeby dołożenie
// drugiego magazynu nie zamieniło odmowy w ciszę: identyfikator nieznany obu
// magazynom ma dalej wracać odmową, i to odmową nazywającą oba miejsca.
func TestWydobycieTekstuOdmawiaKodowiSpozaObuMagazynow(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDocumentTextExtract,
		shared.DocumentTextExtractRequest{AssetId: wskaznik("kod-spoza-magazynow")})

	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Fatalf("odmowa ma nieść kod not_found, a niesie %s", odmowa.Code)
	}
	if !strings.Contains(odmowa.Message, "biblioteki") {
		t.Fatalf("odmowa nie nazywa drugiego magazynu — czytelnik nie wie, gdzie szukano: %q",
			odmowa.Message)
	}
}
