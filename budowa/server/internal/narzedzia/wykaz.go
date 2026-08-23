// Pakiet narzedzia jest rozdzielnią sterowania platformą przez model:
// narzędzie → komenda kontraktu → rdzeń → wynik.
//
// Wykaz pochodzi z kontraktu, nie z tego pakietu. Ani jedna nazwa narzędzia,
// ani jeden opis i ani jedno pole schematu nie są tu zapisane: wszystko czyta
// się z `shared.NarzedziaModelu()` i `shared.KomendyNarzedzi`, wytworzonych
// z `shared/contract.json`. Dopisanie komendy do sekcji `narzedzia`
// kontraktu powiększa ten serwer bez zmiany choćby jednej linii kodu — i tak
// samo działa w drugą stronę: wykreślenie komendy odbiera modelowi narzędzie.
//
// Drugiego wykazu nie ma z zamysłem: wykaz własny rozjechałby się
// z kontraktem, gdy tylko kontrakt urośnie.
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
	// Grupa mówi, do czego narzędzie służy — obszar nazwy jego komendy
	// (`grupa.go`). Pole DANYCH, nie ozdoba wykazu: te same grupy są potem
	// gałęziami drzewa wyboru u Operatora i jednostką doboru narzędzi eksperta
	// (`ekspert_wykaz.go`).
	Grupa string
}

// Wykaz zwraca komplet narzędzi zadeklarowanych przez kontrakt.
//
// Liczba pozycji nie jest tu zapisana ani sprawdzana: wykaz ma tyle pozycji,
// ile ich niesie kontrakt w chwili budowy. Kolejność zachowuje kolejność
// kontraktu, więc ten sam kontrakt zawsze daje ten sam wykaz.
func Wykaz() []Narzedzie {
	return WykazZasiegu(ZasiegOkna)
}

// WykazZasiegu zwraca wykaz kontraktu POWIĘKSZONY o pozycje, które dokłada rola
// okna (`zasieg_roli.go`). Zasięg okna roboczego oddaje sam wykaz kontraktu,
// bez dokładania pozycji roli.
//
// Powiększenie idzie NA KOŃCU listy, żeby kolejność kontraktu została kolejnością
// kontraktu: ten sam kontrakt daje ten sam wykaz, a rola dokłada, nie przestawia.
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
