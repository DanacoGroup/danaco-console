package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

// WartoscZeSciezki wyjmuje wartość tekstową z dokumentu JSON, idąc ścieżką
// rozdzieloną kropkami. Człon liczbowy oznacza pozycję tablicy, na przykład
// "choices.0.delta.content".
//
// Dzięki temu adapter sieciowy nie zna kształtu odpowiedzi żadnego dostawcy:
// kształt jest parametrem wiersza rejestru, nie warunkiem w kodzie.
func WartoscZeSciezki(dokument json.RawMessage, sciezka string) (string, bool) {
	if len(dokument) == 0 || strings.TrimSpace(sciezka) == "" {
		return "", false
	}
	var wezel any
	if err := json.Unmarshal(dokument, &wezel); err != nil {
		return "", false
	}
	for _, czlon := range strings.Split(sciezka, ".") {
		nastepny, jest := zejdz(wezel, czlon)
		if !jest {
			return "", false
		}
		wezel = nastepny
	}
	return jakoTekst(wezel)
}

// zejdz wykonuje jeden krok ścieżki: po kluczu obiektu albo pozycji tablicy.
func zejdz(wezel any, czlon string) (any, bool) {
	switch typowy := wezel.(type) {
	case map[string]any:
		wartosc, jest := typowy[czlon]
		return wartosc, jest
	case []any:
		pozycja, err := strconv.Atoi(czlon)
		if err != nil || pozycja < 0 || pozycja >= len(typowy) {
			return nil, false
		}
		return typowy[pozycja], true
	default:
		return nil, false
	}
}

// jakoTekst przenosi wartość węzła na tekst. Wartość złożona wraca w zapisie
// JSON, dzięki czemu nic z odpowiedzi nie ginie po drodze.
func jakoTekst(wezel any) (string, bool) {
	switch typowy := wezel.(type) {
	case nil:
		return "", false
	case string:
		return typowy, typowy != ""
	case bool:
		return strconv.FormatBool(typowy), true
	case float64:
		return strconv.FormatFloat(typowy, 'f', -1, 64), true
	default:
		surowe, err := json.Marshal(typowy)
		if err != nil {
			return "", false
		}
		return string(surowe), true
	}
}
