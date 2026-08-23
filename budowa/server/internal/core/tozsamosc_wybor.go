// Odpowiedzialność pliku: rozstrzygnięcie, która treść kategorii obowiązuje,
// gdy zapisano ją na kilku osiach naraz.
//
// Reguła pierwszeństwa: konto bije model, a model bije platformę.
//
// Uzasadnienie płynie z kontraktu. Struktura Account niesie pole `defaultModel`
// — to konto wskazuje model, którym pracuje, a nie model konto. Oś konta jest
// więc bytem węższym: jedno konto obsługuje jeden zestaw modeli, jeden model
// bywa obsługiwany przez wiele kont. Do tego konto jest nośnikiem
// uwierzytelnienia i katalogu konfiguracji kanału głównego
// (`CLAUDE_CONFIG_DIR` per okno), czyli tym bytem, przy którym Operator zakłada
// odrębną tożsamość roboczą. Porządek osi powtarza zasadę „węższy wygrywa”,
// tyle że na osi prostopadłej do poziomu zasięgu.
//
// Brak zapisu to co innego niż zapis pusty. Wygrywa zapis czynny osi najwęższej
// — także wtedy, gdy jego treść jest pusta. To jedyna droga, by Operator
// wyciszył kategorię dla jednego konta, nie ruszając platformy, i jest to ta
// sama zasada, którą stosuje rezolwer konfiguracji: brak wiersza znaczy „nie
// ustawiono”, wiersz z wartością pustą znaczy „ustawiono pustą”. Zapis
// nieczynny jest jak brak wiersza — przepuszcza oś szerszą.
package core

import (
	"sort"
	"strings"

	"danacoconsole/shared"
)

// wybranaTresc wiąże kategorię katalogu z treścią, która dla niej obowiązuje.
type wybranaTresc struct {
	Kategoria shared.IdentityCategory
	Dokument  shared.IdentityDocument
	Ma        bool
}

// pierwszenstwoOsi porządkuje osie od najszerszej do najwęższej. Oś nieznana
// katalogowi dostaje pierwszeństwo najniższe — nie wywraca wyboru, po prostu
// nie wygrywa z osią rozpoznaną.
func pierwszenstwoOsi(os shared.ConfigAxis) int {
	switch os {
	case shared.ConfigAxisAccount:
		return 3
	case shared.ConfigAxisModel:
		return 2
	case shared.ConfigAxisPlatform:
		return 1
	default:
		return 0
	}
}

// wybierzTresci przypisuje każdej kategorii katalogu treść obowiązującą.
// Kolejność wyniku jest kolejnością składania nakładki: warstwy według
// krytyczności, a wewnątrz warstwy porządek katalogu.
func wybierzTresci(kategorie []shared.IdentityCategory,
	dokumenty []shared.IdentityDocument) []wybranaTresc {

	najlepsze := map[string]shared.IdentityDocument{}
	for _, dokument := range dokumenty {
		if !dokument.Enabled {
			continue
		}
		obecny, jest := najlepsze[dokument.CategoryId]
		if jest && pierwszenstwoOsi(obecny.Axis) >= pierwszenstwoOsi(dokument.Axis) {
			continue
		}
		najlepsze[dokument.CategoryId] = dokument
	}

	wybrane := make([]wybranaTresc, 0, len(kategorie))
	for _, kategoria := range kategorie {
		if !kategoria.Enabled {
			continue
		}
		dokument, jest := najlepsze[kategoria.Id]
		wybrane = append(wybrane, wybranaTresc{Kategoria: kategoria, Dokument: dokument, Ma: jest})
	}
	uporzadkujTresci(wybrane)
	return wybrane
}

// uporzadkujTresci układa kategorie wg krytyczności warstw, dalej wg porządku
// katalogu, a przy równości wg kodu kategorii. Trzeci klucz nie jest ozdobą:
// bez niego dwie kategorie o tej samej kolejności dawałyby prompt zależny od
// kolejności odczytu, a ta sama konfiguracja ma dawać bajtowo ten sam prompt.
func uporzadkujTresci(wybrane []wybranaTresc) {
	sort.SliceStable(wybrane, func(i, j int) bool {
		lewa, prawa := wybrane[i].Kategoria, wybrane[j].Kategoria
		if porzadekWarstwy(lewa.Layer) != porzadekWarstwy(prawa.Layer) {
			return porzadekWarstwy(lewa.Layer) < porzadekWarstwy(prawa.Layer)
		}
		if lewa.Order != prawa.Order {
			return lewa.Order < prawa.Order
		}
		return lewa.Id < prawa.Id
	})
}

// Tresc zwraca treść wybraną dla kategorii, przyciętą z białych znaków
// brzegowych. Przycięcie jest warunkiem powtarzalności odcisku: ta sama treść
// zapisana z innym zakończeniem wiersza nie ma prawa dać innego promptu.
func (w wybranaTresc) Tresc() string {
	if !w.Ma {
		return ""
	}
	return strings.TrimSpace(w.Dokument.Content)
}

// Wnosi mówi, czy kategoria wchodzi do promptu. Kategoria bez zapisu i kategoria
// wyciszona zapisem pustym nie wnoszą nic — i jedna, i druga liczy się do
// wykazu braków, gdy jest obowiązkowa.
func (w wybranaTresc) Wnosi() bool {
	return w.Tresc() != ""
}
