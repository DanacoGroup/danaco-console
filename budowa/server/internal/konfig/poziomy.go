// Pakiet konfig rozstrzyga ustawienia platformy na dziewięciu poziomach
// i wskazuje, skąd pochodzi wartość. Brak ustawienia nie jest blokadą, lecz
// sięgnięciem po wartość szerszego poziomu, a w ostateczności po wartość
// domyślną rejestru definicji.
package konfig

import "danacoconsole/shared"

// Poziom nazywa poziom zasięgu konfiguracji. Jest aliasem typu kontraktu, więc
// pakiet nie tworzy drugiej definicji tego samego pojęcia. Wartości
// pochodzą wyłącznie ze stałych pakietu shared — tu nie ma ani jednego literału
// nazwy poziomu.
type Poziom = shared.ConfigScope

// PoziomAplikacji jest poziomem najszerszym — obejmuje nastawy samego programu,
// nie treści w nim prowadzonej. Nazwany osobno, bo definicje kluczy wskazują go
// po nazwie, a literału nazwy poziomu w tym pakiecie nie ma.
const PoziomAplikacji Poziom = shared.ConfigScopeApplication

// PoziomBrak oznacza brak poziomu źródłowego: wartość nie została zapisana na
// żadnym poziomie i pochodzi z rejestru definicji.
const PoziomBrak Poziom = ""

// poziomyOdNajwezszego wylicza dziewięć poziomów w kolejności rozstrzygania,
// od najwęższego okna komunikacji do najszerszej aplikacji. Poziom aplikacja
// opisuje program, nie prowadzoną treść, więc zapis na każdym z ośmiu
// poziomów treści wygrywa nad nim.
var poziomyOdNajwezszego = []Poziom{
	shared.ConfigScopeWindow,
	shared.ConfigScopeRole,
	shared.ConfigScopeSession,
	shared.ConfigScopeProject,
	shared.ConfigScopeModulePair,
	shared.ConfigScopeModule,
	shared.ConfigScopeEnvironment,
	shared.ConfigScopeGlobal,
	shared.ConfigScopeApplication,
}

// pierwszenstwaPoziomow służy wyłącznie jako zbiór poziomów znanych dla
// funkcji Znany; kolejność rozstrzygania niesie sama kolejność wykazu
// poziomyOdNajwezszego, a wartości liczbowe nie są porównywane z niczym
// innym.
var pierwszenstwaPoziomow = zbudujPierwszenstwa()

func zbudujPierwszenstwa() map[Poziom]int {
	wynik := make(map[Poziom]int, len(poziomyOdNajwezszego))
	for indeks, poziom := range poziomyOdNajwezszego {
		wynik[poziom] = len(poziomyOdNajwezszego) - indeks
	}
	return wynik
}

// Znany odpowiada, czy poziom należy do dziewięciu poziomów kontraktu,
// wykorzystując zbiór zbudowany z wykazu poziomów.
func Znany(poziom Poziom) bool {
	_, jest := pierwszenstwaPoziomow[poziom]
	return jest
}
