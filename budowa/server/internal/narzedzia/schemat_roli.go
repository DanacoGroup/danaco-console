// Plik składa schemat wejścia narzędzia dołożonego przez rolę okna, odbiciem
// struktury żądania wygenerowanej z kontraktu, drogą osobną od schemat.go.
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

// parametrPola składa jeden parametr schematu wejścia narzędzia z typu Go
// danego pola struktury żądania.
func parametrPola(typ reflect.Type, nazwa string, wymagane bool) shared.ToolParameter {
	rodzaj := typSchematu(typ)
	parametr := shared.ToolParameter{Name: nazwa, Type: rodzaj, Required: wymagane}
	if rodzaj == typTablicy {
		parametr.Items = typSchematu(rozpakuj(typ).Elem())
	}
	return parametr
}

// typSchematu przekłada typ Go na typ zapisu schematu, zgodny z tym, który
// kontrakt wylicza dla narzędzi wykazu.
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
