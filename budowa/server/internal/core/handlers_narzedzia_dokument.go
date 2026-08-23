// Odpowiedzialność pliku: wpięcie dwóch komend rodziny `document.*` — zamiany
// formatu dokumentu i odczytu jego treści. Adapter wraz z rozstrzygnięciami
// (skąd biorą się zasady i obszar uruchomienia, którym binarium jedzie która
// droga, dlaczego `usedOcr` mówi prawdę) leży w plikach
// `adapter_narzedzia_dokument*.go`.
//
// Kontrakt niesie w tej rodzinie dwie komendy: `document.convert`
// i `document.text.extract`.
//
// Są to narzędzia modelu, nie panel Operatora — model wykonuje je sam w trakcie
// tury, dlatego żadna nie niesie okna i żadna niczego nie rozgłasza.
//
// Zdarzeń rodzina nie ma: kontrakt nie zna `document.changed`, więc port nie
// bierze nadajnika. Zdarzenie spoza kontraktu dałoby klientowi nazwę, której
// nie zna nikt poza rdzeniem.
//
// Port niewypełniony nie rejestruje niczego: obie komendy odpowiedzą wtedy
// `document.unknown`, a pozostałe domeny pracują bez zmian.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Dokumenty jest portem rodziny `document.*`. Mówi wyłącznie typami
// kontraktu; binaria arsenału, katalogi tymczasowe i wyjścia programów leżą po
// drugiej stronie adaptera.
type Dokumenty interface {
	// Przeksztalc obsługuje `document.convert`.
	Przeksztalc(ctx context.Context, z shared.DocumentConvertRequest) (shared.DocumentConvertResponse, error)
	// WyciagnijTekst obsługuje `document.text.extract`.
	WyciagnijTekst(ctx context.Context, z shared.DocumentTextExtractRequest) (shared.DocumentTextExtractResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u modelu.
var _ Dokumenty = (*adapterNarzedziDokumentu)(nil)

// zarejestrujDokumenty wpina dwie komendy rodziny `document.*`.
func zarejestrujDokumenty(r *Rejestr, d Dokumenty) {
	if r == nil || d == nil {
		return
	}

	r.Zarejestruj(shared.CommandDocumentConvert, obsluz(d.Przeksztalc))
	r.Zarejestruj(shared.CommandDocumentTextExtract, obsluz(d.WyciagnijTekst))
}
