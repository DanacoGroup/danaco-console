// Odpowiedzialność pliku: schemat wejścia narzędzia dołożonego przez rolę okna.
//
// Droga osobna od `schemat.go`. Kontrakt wylicza `ToolParameter` — gotowy typ
// schematu, typ elementu tablicy i komplet wartości wyliczenia — wyłącznie dla
// komend stojących w sekcji `narzedzia` (generator `shared/gen/narzedzia.mjs`).
// Komenda dołożona przez rolę okna z tej sekcji nie pochodzi, więc jej
// `ToolParameter` nie istnieje i nie ma go skąd wziąć.
//
// Źródłem jest wygenerowana struktura żądania, nie ręczny opis. Struktury
// `shared/contract.go` powstają z tego samego `contract.json`, co wykaz
// narzędzi, więc czytanie ich odbiciem trzyma jedno źródło prawdy: zmiana pola
// w kontrakcie zmienia schemat bez dotykania tego pliku.
//
// Czego ta droga nie daje:
//   - opisu pola. Generator Go kładzie opis w komentarzu, a komentarza odbicie
//     nie widzi, więc pole idzie z opisem pustym.
//   - wykazu wartości wyliczenia. Typ nazwany (`ConfigScope`,
//     `SessionConfigArea`) jest w Go zwykłym napisem, a rejestru wartości
//     `shared` nie wystawia. Pole zostaje napisem — model dostaje odmowę
//     rdzenia przy wartości spoza wykazu i poprawia.
//   - struktury zagnieżdżonej. Pole rodzaju obiektu idzie jako `object` bez
//     właściwości: rozwijanie w głąb (`SessionConfig` to dziesiątki pól
//     i dalsze zagnieżdżenia) urosłoby do schematu większego niż całe okno
//     kontekstu, a granicy głębokości kontrakt nie stanowi.
package narzedzia

import (
	"reflect"
	"strings"

	"danacoconsole/shared"
)

// znacznikPominiecia rozpoznaje pole opcjonalne po zapisie tagu JSON. Pole
// wymagane i pole opcjonalne różnią się w wygenerowanej strukturze dwoma
// znakami naraz — wskaźnikiem i `omitempty` — i wystarczy jeden z nich.
const znacznikPominiecia = "omitempty"

// polaZadania przekłada strukturę treści żądania na parametry schematu.
// Typ inny niż struktura daje pustkę: żądanie bez pól jest żądaniem poprawnym,
// a nie brakiem schematu (tak samo rozstrzyga `schematWejscia`).
func polaZadania(typ reflect.Type) []shared.ToolParameter {
	if typ == nil || typ.Kind() != reflect.Struct {
		return nil
	}
	parametry := make([]shared.ToolParameter, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		pole := typ.Field(i)
		if !pole.IsExported() {
			continue
		}
		nazwa, opcjonalne, niesie := nazwaPolaJSON(pole)
		if !niesie {
			continue
		}
		parametry = append(parametry, parametrPola(pole.Type, nazwa, !opcjonalne))
	}
	return parametry
}

// nazwaPolaJSON odczytuje nazwę pola z tagu i mówi, czy pole jest opcjonalne.
// Pole bez tagu albo z nazwą `-` nie wchodzi do schematu — w kopercie kontraktu
// go nie ma, więc nie ma czego modelowi pokazywać.
func nazwaPolaJSON(pole reflect.StructField) (nazwa string, opcjonalne, niesie bool) {
	tag, jest := pole.Tag.Lookup("json")
	if !jest {
		return "", false, false
	}
	czlony := strings.Split(tag, ",")
	nazwa = strings.TrimSpace(czlony[0])
	if nazwa == "" || nazwa == "-" {
		return "", false, false
	}
	opcjonalne = pole.Type.Kind() == reflect.Pointer
	for _, czlon := range czlony[1:] {
		if strings.TrimSpace(czlon) == znacznikPominiecia {
			opcjonalne = true
		}
	}
	return nazwa, opcjonalne, true
}

// parametrPola składa jeden parametr schematu z typu pola.
func parametrPola(typ reflect.Type, nazwa string, wymagane bool) shared.ToolParameter {
	rodzaj := typSchematu(typ)
	parametr := shared.ToolParameter{Name: nazwa, Type: rodzaj, Required: wymagane}
	if rodzaj == typTablicy {
		parametr.Items = typSchematu(rozpakuj(typ).Elem())
	}
	return parametr
}

// typSchematu przekłada typ Go na typ zapisu schematu. Zapis odpowiada temu,
// który kontrakt wylicza dla narzędzi wykazu (`string`, `integer`, `number`,
// `boolean`, `object`, `array`) — dwa różne zapisy tego samego pojęcia byłyby
// dwiema prawdami o schemacie.
func typSchematu(typ reflect.Type) string {
	switch rozpakuj(typ).Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice, reflect.Array:
		return typTablicy
	default:
		return typObiektu
	}
}

// rozpakuj zdejmuje wskaźnik. Opcjonalność niesie już `nazwaPolaJSON`, więc dla
// kształtu schematu wskaźnik nie znaczy nic.
func rozpakuj(typ reflect.Type) reflect.Type {
	for typ != nil && typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}
