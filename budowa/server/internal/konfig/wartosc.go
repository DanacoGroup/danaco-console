package konfig

import (
	"encoding/json"
	"strings"
)

// Rodzaj nazywa postać wartości ustawienia. Wartości odpowiadają dosłownie
// kolumnie ustawienie.rodzaj_wartosci.
type Rodzaj string

// Wartości Rodzaj.
const (
	RodzajTekst    Rodzaj = "tekst"
	RodzajLiczba   Rodzaj = "liczba"
	RodzajLogiczna Rodzaj = "logiczna"
	RodzajJSON     Rodzaj = "json"
)

// rodzajeZnane zamyka zbiór wartości dopuszczonych przez schemat bazy.
var rodzajeZnane = map[Rodzaj]struct{}{
	RodzajTekst:    {},
	RodzajLiczba:   {},
	RodzajLogiczna: {},
	RodzajJSON:     {},
}

// Znany odpowiada, czy rodzaj mieści się w zbiorze dopuszczonym przez schemat.
func (r Rodzaj) Znany() bool {
	_, jest := rodzajeZnane[r]
	return jest
}

// RodzajLubTekst zwraca rodzaj znany, a dla wartości nierozpoznanej — tekst.
// Nierozpoznany rodzaj nie przerywa odczytu.
func RodzajLubTekst(rodzaj Rodzaj) Rodzaj {
	if rodzaj.Znany() {
		return rodzaj
	}
	return RodzajTekst
}

// KodujJSON zamienia wartość zapisaną w bazie na treść pola value koperty
// kontraktu. Liczba, wartość logiczna i JSON idą surowo, jeżeli są poprawnym
// JSON-em; wszystko pozostałe idzie jako napis. Funkcja nigdy nie zawodzi —
// wartość uszkodzona trafia do kontraktu jako napis, nie jako błąd.
func KodujJSON(wartosc string, rodzaj Rodzaj) json.RawMessage {
	switch RodzajLubTekst(rodzaj) {
	case RodzajLiczba, RodzajLogiczna, RodzajJSON:
		przyciete := strings.TrimSpace(wartosc)
		if przyciete != "" && json.Valid([]byte(przyciete)) {
			return json.RawMessage(przyciete)
		}
	}
	surowa, err := json.Marshal(wartosc)
	if err != nil {
		return json.RawMessage(`""`)
	}
	return surowa
}

// DekodujJSON zamienia treść pola value koperty kontraktu na parę
// wartość + rodzaj gotową do zapisu w kolumnach ustawienie.wartosc
// oraz ustawienie.rodzaj_wartosci.
func DekodujJSON(surowa json.RawMessage) (string, Rodzaj) {
	przyciete := strings.TrimSpace(string(surowa))
	if przyciete == "" {
		return "", RodzajTekst
	}
	var napis string
	if err := json.Unmarshal(surowa, &napis); err == nil {
		return napis, RodzajTekst
	}
	var logiczna bool
	if err := json.Unmarshal(surowa, &logiczna); err == nil {
		return przyciete, RodzajLogiczna
	}
	var liczba float64
	if err := json.Unmarshal(surowa, &liczba); err == nil {
		return przyciete, RodzajLiczba
	}
	return przyciete, RodzajJSON
}
