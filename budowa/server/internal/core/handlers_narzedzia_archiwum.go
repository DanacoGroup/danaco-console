// Plik definiuje port i wpina dwie komendy rodziny archive.* — spakowanie
// pracy w archiwum i jego rozpakowanie, dostępne zarazem jako narzędzia
// modelu.
package core

import (
	"context"

	"danacoconsole/shared"
)

// NarzedziaArchiwum jest portem rodziny `archive.*`. Mówi wyłącznie
// typami kontraktu; `7z`, jego wiersz poleceń, kwarantanna rozpakowania i blob
// wyniku leżą po drugiej stronie adaptera.
type NarzedziaArchiwum interface {
	// Spakuj obsługuje `archive.pack` — złożenie pracy w jeden plik.
	Spakuj(ctx context.Context, z shared.ArchivePackRequest) (shared.ArchivePackResponse, error)
	// Rozpakuj obsługuje `archive.unpack` — otwarcie archiwum.
	Rozpakuj(ctx context.Context, z shared.ArchiveUnpackRequest) (shared.ArchiveUnpackResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ NarzedziaArchiwum = (*adapterNarzedziArchiwum)(nil)

// zarejestrujNarzedziaArchiwum wpina dwie komendy rodziny archive.* na
// porcie NarzedziaArchiwum, pomijając rejestrację bez portu.
func zarejestrujNarzedziaArchiwum(r *Rejestr, n NarzedziaArchiwum) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandArchivePack, obsluz(n.Spakuj))
	r.Zarejestruj(shared.CommandArchiveUnpack, obsluz(n.Rozpakuj))
}
