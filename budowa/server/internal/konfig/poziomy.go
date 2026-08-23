// Pakiet konfig rozstrzyga ustawienia platformy na dziewięciu poziomach zasięgu
// i dla każdej wartości wskazuje, skąd pochodzi.
//
// Zasada nadrzędna pakietu: brak ustawienia nie jest blokadą, lecz sięgnięciem
// po wartość szerszego poziomu, a w ostateczności po wartość domyślną z rejestru
// definicji. Żadna ścieżka pakietu nie odmawia rozstrzygnięcia — także błąd
// źródła danych kończy się polityką domyślną.
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

// poziomyOdNajwezszego wylicza dziewięć poziomów w kolejności rozstrzygania:
// od najwęższego (okno komunikacji) do najszerszego (aplikacja).
// Kolejność jest odwrotnością kolumny poziom_zasiegu.pierwszenstwo.
//
// Aplikacja stoi na końcu, czyli przegrywa ze wszystkim. Poziom ten opisuje sam
// program — wymóg logowania, postać nasłuchu — a nie treść w nim prowadzoną,
// więc zapis na którymkolwiek z ośmiu poziomów treści wygrywa nad poziomem
// aplikacji.
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

// pierwszenstwaPoziomow służy wyłącznie za zbiór poziomów znanych (funkcja
// Znany niżej) — wartości liczbowe nie są tu porównywane z niczym, bo kolejność
// rozstrzygania niesie sama kolejność wykazu wyżej. Kolumna
// poziom_zasiegu.pierwszenstwo mieszka w bazie i jest jej własną numeracją
// (aplikacja 0, globalny 1, ..., okno 8); powielanie jej tutaj byłoby drugą
// prawdą o tej samej kolejności.
var pierwszenstwaPoziomow = zbudujPierwszenstwa()

func zbudujPierwszenstwa() map[Poziom]int {
	wynik := make(map[Poziom]int, len(poziomyOdNajwezszego))
	for indeks, poziom := range poziomyOdNajwezszego {
		wynik[poziom] = len(poziomyOdNajwezszego) - indeks
	}
	return wynik
}

// Znany odpowiada, czy poziom należy do dziewięciu poziomów kontraktu.
func Znany(poziom Poziom) bool {
	_, jest := pierwszenstwaPoziomow[poziom]
	return jest
}
