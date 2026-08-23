package stdio

import (
	"context"
	"encoding/json"

	"danacoconsole/server/internal/narzedzia"
)

// Nazwy metod protokołu MCP obsługiwanych przez serwer.
const (
	metodaPowitania = "initialize"
	metodaPingu     = "ping"
	metodaWykazu    = "tools/list"
	metodaWywolania = "tools/call"
)

// wywolanieNarzedzia jest treścią żądania `tools/call`.
type wywolanieNarzedzia struct {
	Nazwa     string         `json:"name"`
	Argumenty map[string]any `json:"arguments"`
}

// odpowiedzNaMetode rozstrzyga jedno wywołanie. Zwraca fałsz dla zawiadomienia,
// na które protokół odpowiedzi nie przewiduje.
func odpowiedzNaMetode(kontekst context.Context, katalog Katalog, z zadanie) (odpowiedz, bool) {
	if len(z.Id) == 0 {
		return odpowiedz{}, false // zawiadomienie: przyjęte i odłożone bez odpowiedzi
	}
	switch z.Metoda {
	case metodaPowitania:
		return wynikiem(z.Id, powitanie()), true
	case metodaPingu:
		return wynikiem(z.Id, map[string]any{}), true
	case metodaWykazu:
		return wynikiem(z.Id, map[string]any{"tools": wykazNarzedzi(kontekst, katalog)}), true
	case metodaWywolania:
		return wywolaj(kontekst, katalog, z), true
	default:
		return bledem(z.Id, kodBrakMetody, "metoda "+z.Metoda+" nie jest obsługiwana"), true
	}
}

// powitanie zwraca odpowiedź na `initialize`: wersję protokołu, zdolności
// serwera i jego nazwę. Serwer ogłasza wyłącznie narzędzia — zasobów ani
// podpowiedzi nie prowadzi, więc ich nie zapowiada.
func powitanie() map[string]any {
	return map[string]any{
		"protocolVersion": wersjaMcp,
		"capabilities":    map[string]any{"tools": map[string]any{}},
		"serverInfo":      map[string]any{"name": nazwaSerwera, "version": wersjaSerwera},
	}
}

// wykazNarzedzi przekłada wykaz rozdzielni na kształt `tools/list`.
//
// Kształt pozycji składa `narzedzia.PozycjeWykazu`, a nie ta funkcja: licznik
// wykazu (`narzedzia/licznik.go`) mierzy wagę tej samej odpowiedzi w bajtach,
// a dwa miejsca składające ją osobno dałyby pomiar czegoś innego, niż dostaje
// model. Tutaj zostaje wyłącznie zamówienie wykazu.
func wykazNarzedzi(kontekst context.Context, katalog Katalog) []map[string]any {
	return narzedzia.PozycjeWykazu(katalog.Narzedzia(kontekst))
}

// wywolaj wykonuje `tools/call`.
//
// Odmowa narzędzia i błąd rdzenia wracają wynikiem oznaczonym jako błędny, a nie
// błędem protokołu: model ma je przeczytać i poprawić wywołanie, zamiast dostać
// usterkę połączenia. Błędem protokołu zostaje wyłącznie treść żądania, której
// nie da się odczytać.
func wywolaj(kontekst context.Context, katalog Katalog, z zadanie) odpowiedz {
	var wywolanie wywolanieNarzedzia
	if err := json.Unmarshal(z.Params, &wywolanie); err != nil {
		return bledem(z.Id, kodZleParametry, "treść wywołania narzędzia nieczytelna: "+err.Error())
	}
	if wywolanie.Nazwa == "" {
		return bledem(z.Id, kodZleParametry, "wywołanie narzędzia bez pola name")
	}
	tresc, err := katalog.Wywolaj(kontekst, wywolanie.Nazwa, wywolanie.Argumenty)
	if err != nil {
		return wynikiem(z.Id, trescWyniku(err.Error(), true))
	}
	return wynikiem(z.Id, trescWyniku(tresc, false))
}

// trescWyniku składa wynik `tools/call` w kształcie protokołu.
func trescWyniku(tekst string, blad bool) map[string]any {
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": tekst}},
		"isError": blad,
	}
}
