// Plik wpina porty i rejestrację trzech rodzin przekrojowych obsługi tekstu: historii schowka,
// słownika skrótów oraz nazwanych kontekstów pamięci i zasad retencji.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Schowek jest portem rodziny clipboard.* obsługującym odczyt i zapis historii schowka tego konta użytkownika.
type Schowek interface {
	WykazSchowka(ctx context.Context, z shared.ClipboardListRequest) (shared.ClipboardListResponse, error)
	DopiszDoSchowka(ctx context.Context, z shared.ClipboardPushRequest) (shared.ClipboardPushResponse, error)
	PrzypnijWpisSchowka(ctx context.Context, z shared.ClipboardPinRequest) (shared.ClipboardPinResponse, error)
	UsunZeSchowka(ctx context.Context, z shared.ClipboardDeleteRequest) (shared.ClipboardDeleteResponse, error)
}

// SkrotyTekstowe jest portem rodziny snippet.* obsługującym słownik skrótów tekstowych tego konta użytkownika.
type SkrotyTekstowe interface {
	WykazSkrotow(ctx context.Context, z shared.SnippetListRequest) (shared.SnippetListResponse, error)
	ZapiszSkrot(ctx context.Context, z shared.SnippetSetRequest) (shared.SnippetSetResponse, error)
	UsunSkrot(ctx context.Context, z shared.SnippetDeleteRequest) (shared.SnippetDeleteResponse, error)
}

// KontekstyPamieci jest portem sześciu komend rodziny `memory.*`, które nie
// dotyczą pojedynczego wpisu pamięci.
type KontekstyPamieci interface {
	WykazKontekstow(ctx context.Context, z shared.MemoryContextListRequest) (shared.MemoryContextListResponse, error)
	ZapiszKontekst(ctx context.Context, z shared.MemoryContextSaveRequest) (shared.MemoryContextSaveResponse, error)
	UaktywnijKontekst(ctx context.Context, z shared.MemoryContextActivateRequest) (shared.MemoryContextActivateResponse, error)
	UsunKontekst(ctx context.Context, z shared.MemoryContextDeleteRequest) (shared.MemoryContextDeleteResponse, error)
	ZasadyRetencji(ctx context.Context, z shared.MemoryRetentionGetRequest) (shared.MemoryRetentionGetResponse, error)
	ZapiszZasadeRetencji(ctx context.Context, z shared.MemoryRetentionSetRequest) (shared.MemoryRetentionSetResponse, error)
}

// Adaptery wypełniają wszystkie trzy porty w całości, bez żadnej metody pozostawionej niezaimplementowanej.
var (
	_ Schowek          = (*adapterSchowka)(nil)
	_ SkrotyTekstowe   = (*adapterSkrotowTekstowych)(nil)
	_ KontekstyPamieci = (*adapterKontekstowPamieci)(nil)
)

// zarejestrujSchowek wpina cztery komendy rodziny clipboard.* obsługujące historię schowka tego konta.
func zarejestrujSchowek(r *Rejestr, s Schowek) {
	if r == nil || s == nil {
		return
	}
	r.Zarejestruj(shared.CommandClipboardList, obsluz(s.WykazSchowka))
	r.Zarejestruj(shared.CommandClipboardPush, obsluz(s.DopiszDoSchowka))
	r.Zarejestruj(shared.CommandClipboardPin, obsluz(s.PrzypnijWpisSchowka))
	r.Zarejestruj(shared.CommandClipboardDelete, obsluz(s.UsunZeSchowka))
}

// zarejestrujSkrotyTekstowe wpina trzy komendy rodziny snippet.* obsługujące słownik skrótów tekstowych konta.
func zarejestrujSkrotyTekstowe(r *Rejestr, s SkrotyTekstowe) {
	if r == nil || s == nil {
		return
	}
	r.Zarejestruj(shared.CommandSnippetList, obsluz(s.WykazSkrotow))
	r.Zarejestruj(shared.CommandSnippetSet, obsluz(s.ZapiszSkrot))
	r.Zarejestruj(shared.CommandSnippetDelete, obsluz(s.UsunSkrot))
}

// zarejestrujKontekstyPamieci wpina sześć komend rodziny memory.* spoza obszaru pojedynczego
// wpisu, obsługujących konteksty pamięci i zasady retencji.
func zarejestrujKontekstyPamieci(r *Rejestr, k KontekstyPamieci) {
	if r == nil || k == nil {
		return
	}
	r.Zarejestruj(shared.CommandMemoryContextList, obsluz(k.WykazKontekstow))
	r.Zarejestruj(shared.CommandMemoryContextSave, obsluz(k.ZapiszKontekst))
	r.Zarejestruj(shared.CommandMemoryContextActivate, obsluz(k.UaktywnijKontekst))
	r.Zarejestruj(shared.CommandMemoryContextDelete, obsluz(k.UsunKontekst))
	r.Zarejestruj(shared.CommandMemoryRetentionGet, obsluz(k.ZasadyRetencji))
	r.Zarejestruj(shared.CommandMemoryRetentionSet, obsluz(k.ZapiszZasadeRetencji))
}
