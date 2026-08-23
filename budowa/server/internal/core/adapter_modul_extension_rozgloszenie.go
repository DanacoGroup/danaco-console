// Odpowiedzialność pliku: rozgłoszenie rodziny `extension.*` — jedno zdarzenie
// kontraktu, `extension.changed`, wraz z regułą, która komenda niesie który
// rodzaj zmiany. Podział wobec `adapter_modul_extension.go` idzie wzdłuż
// odpowiedzialności, wzorem pary `adapter_modul_isolation.go` /
// `adapter_modul_isolation_rozgloszenie.go`.
//
// Stan pozycji katalogu zmieniają trzy komendy z pięciu:
//
//	`extension.install`   → `created`
//	`extension.configure` → `updated`
//	`extension.toggle`    → `updated`
//	`extension.uninstall` → `deleted`
//
// `extension.list` niczego nie zmienia i niczego nie rozgłasza.
//
// Przywrócenie pozycji odinstalowanej idzie jako `created`, nie `updated`.
// Wiersz stoi w katalogu przez cały czas (odinstalowanie zdejmuje znaczniki,
// nie wiersz — nagłówek `adapter_modul_extension.go`), ale zdarzenie opisuje
// pozycję katalogu, nie wiersz tabeli. Skoro odinstalowanie mówi `deleted`, to
// instalacja tej samej pozycji musi mówić `created` — inaczej para komend
// byłaby niesymetryczna i klient, który po `deleted` zdjął pozycję z widoku, po
// `updated` nie miałby czego zaktualizować.
//
// `extension.uninstall` nie ma kłopotu z ładunkiem sprzed usunięcia. Wzorzec
// `memory.changed` czyta wpis przed skasowaniem, bo `memory.delete` kasuje
// wiersz i po nim nie ma czego włożyć w ładunek. Tutaj kasowania nie ma:
// `ZmienRozszerzenie` oddaje wiersz po zdjęciu znaczników i tego żąda kontrakt
// — pole `extension` opisane jest jako „Rozszerzenie po zmianie". Odczyt
// uprzedni dałby ładunek nieprawdziwy: pozycję ze znacznikami
// `installed`/`enabled` jeszcze ustawionymi, przy `change: deleted`.
//
// Rozgłoszenie idzie wyłącznie po udanym zapisie i nigdy po odmowie; brak
// nadajnika nie wywraca komendy, bo `emiter.wyslij` na nilu milczy.
//
// Dopóki `montaz_porty.go` składa port bez wywołania `ZRozgloszeniem`,
// `a.rozgloszenie` zostaje nilem i wszystkie cztery wywołania niżej milkną.
// Wstawka montażu stoi w nagłówku `handlers_extension.go`.
package core

import "danacoconsole/shared"

// ZRozgloszeniem dokłada nadajnik zdarzeń `extension.changed`.
// Nadajnik niepodłączony nie wstrzymuje żadnej z pięciu komend.
func (a *adapterRozszerzen) ZRozgloszeniem(nadajnik Nadajnik) *adapterRozszerzen {
	a.rozgloszenie = nowyEmiter(nadajnik)
	return a
}

// rozglosRozszerzenie rozgłasza `extension.changed` po udanej zmianie pozycji
// katalogu. Zdarzenie idzie bez wskazania sesji: katalog rozszerzeń stoi poziom
// wyżej niż karta sesji i wyżej niż ekspert (nagłówek
// `adapter_modul_extension.go`), więc zmiana dotyczy każdego połączenia konta.
func (a *adapterRozszerzen) rozglosRozszerzenie(zmiana shared.ChangeKind, rozszerzenie shared.Extension) {
	if a == nil || a.rozgloszenie == nil || rozszerzenie.Id == "" {
		return
	}
	a.rozgloszenie.wyslij(shared.EventExtensionChanged, "",
		shared.ExtensionChangedEvent{Change: zmiana, Extension: rozszerzenie})
}
