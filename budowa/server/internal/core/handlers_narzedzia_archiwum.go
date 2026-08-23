// Odpowiedzialność pliku: port i wpięcie dwóch komend rodziny `archive.*` —
// `archive.pack` i `archive.unpack`. Adapter leży w
// `adapter_narzedzia_archiwum.go`, pakowanie i rozpakowanie w plikach
// czynności, wyrok o zawartości archiwum w pliku spisu.
//
// Tyle niesie kontrakt w tej rodzinie (`shared/contract.go`:
// CommandArchivePack, CommandArchiveUnpack) i obie są zarazem narzędziami
// modelu — `danaco_archive_pack` i `danaco_archive_unpack` w wykazie narzędzi.
// Trzecia pozycja o tym przedrostku, `archive.unknown`, jest odpowiedzią na
// komendę nieznaną obszaru, a nie komendą: nie ma pary żądanie/wynik i
// w rejestrze komend się nie zjawia.
//
// Zdarzeń rodzina nie ma, więc port nie bierze nadajnika. Archiwum wytworzone
// przez `archive.pack` jest zasobem modułu Design i to jego rodzina zdarzeń
// o zasobach mówi.
//
// Port niewypełniony nie rejestruje niczego: obie komendy odpowiedzą wtedy
// `archive.unknown`, a pozostałe domeny pracują bez zmian. Rdzeń niczym nie
// warunkuje startu.
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

// zarejestrujNarzedziaArchiwum wpina dwie komendy rodziny `archive.*`.
func zarejestrujNarzedziaArchiwum(r *Rejestr, n NarzedziaArchiwum) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandArchivePack, obsluz(n.Spakuj))
	r.Zarejestruj(shared.CommandArchiveUnpack, obsluz(n.Rozpakuj))
}
