// Pakiet narzedzia jest rozdzielnią sterowania platformą przez model:
// narzędzie, komenda kontraktu, rdzeń, wynik; wykaz pochodzi z kontraktu, nie
// z tego pakietu.
package narzedzia

import "danacoconsole/shared"

// Narzedzie jest jedną pozycją wykazu podawaną kanałowi modelu: nazwa, opis
// i schemat wejścia. Kształt odpowiada polom `tools/list` protokołu MCP, lecz
// pakiet nie zna samego protokołu — składa dane, nie ramki.
type Narzedzie struct {
	// Nazwa jest nazwą narzędzia widzianą przez model, wprost z kontraktu.
	Nazwa string
	// Opis mówi modelowi, czym narzędzie jest i kiedy po nie sięgnąć.
	Opis string
	// Schemat jest schematem wejścia (JSON Schema) wyprowadzonym z pól treści
	// żądania komendy.
	Schemat map[string]any
	// Grupa mówi, do czego narzędzie służy — obszar nazwy jego komendy, pole
	// danych, nie ozdoba wykazu.
	Grupa string
}

// Wykaz zwraca komplet narzędzi zadeklarowanych przez kontrakt, w kolejności
// kontraktu, bez zapisanej liczby pozycji.
func Wykaz() []Narzedzie {
	return WykazZasiegu(ZasiegOkna)
}

// WykazZasiegu zwraca wykaz kontraktu powiększony na końcu o pozycje, które
// dokłada rola okna, zachowując kolejność kontraktu.
func WykazZasiegu(zasieg Zasieg) []Narzedzie {
	deklaracje := shared.NarzedziaModelu()
	wykaz := make([]Narzedzie, 0, len(deklaracje))
	for _, deklaracja := range deklaracje {
		wykaz = append(wykaz, Narzedzie{
			Nazwa:   deklaracja.Name,
			Opis:    deklaracja.Description,
			Schemat: schematWejscia(deklaracja.Parameters),
			Grupa:   grupaKomendy(deklaracja.Command),
		})
	}
	return append(wykaz, narzedziaRoli(zasieg)...)
}

// deklaracja odszukuje deklarację narzędzia po nazwie. Zwraca fałsz dla nazwy
// spoza wykazu kontraktu — rozstrzygnięcie, co z taką nazwą zrobić, należy do
// rozdzielni, nie do wyszukiwania.
func deklaracja(nazwa string) (shared.ToolDeclaration, bool) {
	for _, pozycja := range shared.NarzedziaModelu() {
		if pozycja.Name == nazwa {
			return pozycja, true
		}
	}
	return shared.ToolDeclaration{}, false
}
