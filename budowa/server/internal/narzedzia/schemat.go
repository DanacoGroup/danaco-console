// Plik przekłada schemat wejścia z kontraktu na schemat wejścia protokołu MCP
// w postaci JSON Schema, płytko, bez powtarzania rozstrzygnięć typów.
package narzedzia

import "danacoconsole/shared"

// Nazwy pól schematu wejścia. Stoją osobno, bo padają w kilku miejscach pliku,
// a literał powtórzony jest literałem, który da się rozjechać.
const (
	poleTypu        = "type"
	poleWlasciwosci = "properties"
	poleWymaganych  = "required"
	poleOpisu       = "description"
	poleWartosci    = "enum"
	poleElementu    = "items"
	typObiektu      = "object"
	typTablicy      = "array"
)

// schematWejscia składa schemat wejścia narzędzia z pól treści żądania komendy.
// Komenda bez pól daje schemat obiektu bez właściwości — narzędzie wywoływane
// pustą treścią jest narzędziem poprawnym, nie brakiem schematu.
func schematWejscia(parametry []shared.ToolParameter) map[string]any {
	wlasciwosci := make(map[string]any, len(parametry))
	wymagane := make([]string, 0, len(parametry))
	for _, parametr := range parametry {
		wlasciwosci[parametr.Name] = wlasciwoscSchematu(parametr)
		if parametr.Required {
			wymagane = append(wymagane, parametr.Name)
		}
	}
	schemat := map[string]any{
		poleTypu:        typObiektu,
		poleWlasciwosci: wlasciwosci,
	}
	// Puste required różni się od jego braku dla części klientów MCP, więc pole
	// idzie tylko z treścią.
	if len(wymagane) > 0 {
		schemat[poleWymaganych] = wymagane
	}
	return schemat
}

// wlasciwoscSchematu opisuje jedno pole treści żądania.
//
// Wyliczenie przy polu tablicowym dotyczy elementu, nie tablicy (tak stanowi
// kontrakt), dlatego dopuszczalne wartości schodzą wtedy do `items`.
func wlasciwoscSchematu(parametr shared.ToolParameter) map[string]any {
	wlasciwosc := map[string]any{
		poleTypu:  parametr.Type,
		poleOpisu: parametr.Description,
	}
	if parametr.Type != typTablicy {
		if len(parametr.Enum) > 0 {
			wlasciwosc[poleWartosci] = parametr.Enum
		}
		return wlasciwosc
	}
	element := map[string]any{poleTypu: parametr.Items}
	if len(parametr.Enum) > 0 {
		element[poleWartosci] = parametr.Enum
	}
	wlasciwosc[poleElementu] = element
	return wlasciwosc
}
