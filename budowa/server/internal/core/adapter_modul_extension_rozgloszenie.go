// Odpowiedzialność pliku: rozgłoszenie rodziny `extension.*` — jedno zdarzenie
// kontraktu, `extension.changed`, wraz z regułą, która komenda niesie który
// rodzaj zmiany, wyłącznie po udanym zapisie, nigdy po odmowie.
package core

import "danacoconsole/shared"

// ZRozgloszeniem dokłada nadajnik zdarzeń `extension.changed`.
// Nadajnik niepodłączony nie wstrzymuje żadnej z pięciu komend.
func (a *adapterRozszerzen) ZRozgloszeniem(nadajnik Nadajnik) *adapterRozszerzen {
	a.rozgloszenie = nowyEmiter(nadajnik)
	return a
}

// rozglosRozszerzenie rozgłasza `extension.changed` po udanej zmianie pozycji
// katalogu. Zdarzenie idzie bez wskazania sesji: katalog rozszerzeń stoi
// poziom wyżej niż karta sesji, więc zmiana dotyczy każdego połączenia konta.
func (a *adapterRozszerzen) rozglosRozszerzenie(zmiana shared.ChangeKind, rozszerzenie shared.Extension) {
	if a == nil || a.rozgloszenie == nil || rozszerzenie.Id == "" {
		return
	}
	a.rozgloszenie.wyslij(shared.EventExtensionChanged, "",
		shared.ExtensionChangedEvent{Change: zmiana, Extension: rozszerzenie})
}
