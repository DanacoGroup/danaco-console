// Plik wpina dwie komendy rodziny document.* — zamianę formatu dokumentu
// i odczyt jego treści, jako narzędzia modelu wykonywane samodzielnie
// w trakcie tury.
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

// zarejestrujDokumenty wpina dwie komendy rodziny document.* na porcie
// Dokumenty, pomijając rejestrację, gdy port jest pusty.
func zarejestrujDokumenty(r *Rejestr, d Dokumenty) {
	if r == nil || d == nil {
		return
	}

	r.Zarejestruj(shared.CommandDocumentConvert, obsluz(d.Przeksztalc))
	r.Zarejestruj(shared.CommandDocumentTextExtract, obsluz(d.WyciagnijTekst))
}
