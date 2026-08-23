package konfig

import "sort"

// Polityka to komplet ustawień obowiązujących w danym kontekście wraz ze
// wskazaniem, skąd pochodzi każda wartość. To jest podgląd polityki efektywnej.
type Polityka struct {
	Kontekst Kontekst
	Pozycje  []Wynik
	// BladZrodla niesie błąd odczytu warstwy trwałości. Polityka jest zwracana
	// mimo błędu — wartości pochodzą wtedy z rejestru definicji.
	BladZrodla error
}

// PolitykaEfektywna rozstrzyga komplet ustawień rejestru na jednym odczycie
// źródła i dokłada klucze zapisane, lecz nieujęte w rejestrze — te ostatnie
// pozostają widoczne, bo odczyt nie jest bramą.
func (r *Rozstrzygacz) PolitykaEfektywna(kontekst Kontekst) Polityka {
	zapisy, blad := r.zapisy(kontekst)
	klucze := append(r.Rejestr().Klucze(), kluczeSpozaRejestru(r.Rejestr(), zapisy)...)

	pozycje := make([]Wynik, 0, len(klucze))
	for _, klucz := range klucze {
		pozycje = append(pozycje, r.rozstrzygnijZZapisow(kontekst, klucz, zapisy))
	}
	return Polityka{Kontekst: kontekst, Pozycje: pozycje, BladZrodla: blad}
}

// kluczeSpozaRejestru zwraca posortowane klucze obecne w zapisach, a nieznane
// rejestrowi definicji.
func kluczeSpozaRejestru(rejestr *Rejestr, zapisy map[kluczWpisu]Wpis) []string {
	zebrane := make(map[string]struct{})
	for _, wpis := range zapisy {
		if _, jest := rejestr.Definicja(wpis.Klucz); !jest {
			zebrane[wpis.Klucz] = struct{}{}
		}
	}
	klucze := make([]string, 0, len(zebrane))
	for klucz := range zebrane {
		klucze = append(klucze, klucz)
	}
	sort.Strings(klucze)
	return klucze
}

// Pozycja zwraca rozstrzygnięcie jednego ustawienia z gotowej polityki.
func (p Polityka) Pozycja(klucz string) (Wynik, bool) {
	for _, pozycja := range p.Pozycje {
		if pozycja.Klucz == klucz {
			return pozycja, true
		}
	}
	return Wynik{}, false
}

// Zrodlo zwraca poziom, z którego pochodzi wartość klucza, oraz pochodzenie.
// To jest odpowiedź na pytanie „skąd wzięła się ta wartość".
func (p Polityka) Zrodlo(klucz string) (Poziom, Pochodzenie) {
	pozycja, jest := p.Pozycja(klucz)
	if !jest {
		return PoziomBrak, PochodzenieNieznane
	}
	return pozycja.Poziom, pozycja.Pochodzenie
}
