// Plik rozstrzyga, która treść kategorii obowiązuje, gdy zapisano ją na kilku osiach naraz. Reguła pierwszeństwa: konto bije model, model bije platformę. Wygrywa zapis czynny osi najwęższej, także wtedy, gdy jego treść jest pusta.
package core

import (
	"sort"
	"strings"

	"danacoconsole/shared"
)

// wybranaTresc wiąże kategorię katalogu z treścią, która dla niej obowiązuje, i informacją, czy w ogóle istnieje.
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

// uporzadkujTresci układa kategorie wg krytyczności warstw, dalej wg porządku katalogu, a przy równości wg kodu kategorii, żeby ta sama konfiguracja dawała bajtowo ten sam prompt, niezależnie od kolejności odczytu.
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
