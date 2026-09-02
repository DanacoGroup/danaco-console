// Obsługa komend memory.context i memory.retention: kontekst jest zestawem wskazań
// na wpisy pamięci, nie ich właścicielem, więc jego usunięcie nie kasuje wpisu.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	przedrostekKontekstuPamieci = "kontekst-pamieci-"
	przedrostekZasadyRetencji   = "retencja-"
)

type adapterKontekstowPamieci struct {
	repozytorium dane.RepozytoriumKontekstowPamieci
}

func nowyAdapterKontekstowPamieci(
	repozytorium dane.RepozytoriumKontekstowPamieci) *adapterKontekstowPamieci {

	return &adapterKontekstowPamieci{repozytorium: repozytorium}
}

func (a *adapterKontekstowPamieci) WykazKontekstow(ctx context.Context,
	z shared.MemoryContextListRequest) (shared.MemoryContextListResponse, error) {

	if a.repozytorium == nil {
		return shared.MemoryContextListResponse{}, bladZapleczaKontekstow()
	}
	zWylaczonymi := z.IncludeDisabled != nil && *z.IncludeDisabled
	konteksty, err := a.repozytorium.KontekstyPamieci(ctx, wartoscTekstu(z.ProfileId), zWylaczonymi)
	if err != nil {
		return shared.MemoryContextListResponse{}, bladMagazynuKontekstow(err)
	}
	wykaz := make([]shared.MemoryContext, 0, len(konteksty))
	for _, kontekst := range konteksty {
		wykaz = append(wykaz, kontekstPamieciKontraktu(kontekst))
	}

	odpowiedz := shared.MemoryContextListResponse{Contexts: wykaz}
	if karta := strings.TrimSpace(wartoscTekstu(z.SessionId)); karta != "" {
		czynny, err := a.repozytorium.CzynnyKontekstPamieci(ctx, karta)
		if err != nil {
			return shared.MemoryContextListResponse{}, bladMagazynuKontekstow(err)
		}
		if czynny != "" {
			odpowiedz.ActiveId = &czynny
		}
	}
	return odpowiedz, nil
}

func (a *adapterKontekstowPamieci) ZapiszKontekst(ctx context.Context,
	z shared.MemoryContextSaveRequest) (shared.MemoryContextSaveResponse, error) {

	if a.repozytorium == nil {
		return shared.MemoryContextSaveResponse{}, bladZapleczaKontekstow()
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.MemoryContextSaveResponse{}, bladWskazaniaKontekstu(
			"kontekst bez nazwy — Operator nie miałby czego wybrać z selektora")
	}
	for _, poziom := range z.Levels {
		if err := sprawdzPoziomKontekstu(poziom); err != nil {
			return shared.MemoryContextSaveResponse{}, err
		}
	}

	teraz := time.Now().UnixMilli()
	kod := strings.TrimSpace(wartoscTekstu(z.ContextId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekKontekstuPamieci)
	}
	kontekst := dane.KontekstPamieci{
		Kod:             kod,
		Nazwa:           strings.TrimSpace(z.Name),
		Opis:            z.Description,
		ProfilKod:       strings.TrimSpace(wartoscTekstu(z.ProfileId)),
		PoziomyJSON:     zapisPoziomowKontekstu(z.Levels),
		WpisyJSON:       zapisPolSzablonu(z.EntryIds),
		PromptSystemowy: z.SystemPrompt,
		Czynny:          z.Enabled == nil || *z.Enabled,
		Utworzono:       teraz,
		Zaktualizowano:  teraz,
	}
	zapisany, err := a.repozytorium.ZapiszKontekstPamieci(ctx, kontekst)
	if err != nil {
		return shared.MemoryContextSaveResponse{}, bladMagazynuKontekstow(err)
	}
	return shared.MemoryContextSaveResponse{Context: kontekstPamieciKontraktu(zapisany)}, nil
}

func (a *adapterKontekstowPamieci) UaktywnijKontekst(ctx context.Context,
	z shared.MemoryContextActivateRequest) (shared.MemoryContextActivateResponse, error) {

	if a.repozytorium == nil {
		return shared.MemoryContextActivateResponse{}, bladZapleczaKontekstow()
	}
	if strings.TrimSpace(z.SessionId) == "" {
		return shared.MemoryContextActivateResponse{}, bladWskazaniaKontekstu(
			"aktywacja bez wskazania karty sesji — kontekst czynny należy do karty, nie do serwera")
	}
	kontekst, err := a.repozytorium.KontekstPamieciPoKodzie(ctx, strings.TrimSpace(z.ContextId))
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.MemoryContextActivateResponse{}, bladNieznanegoKontekstuPamieci(z.ContextId)
		}
		return shared.MemoryContextActivateResponse{}, bladMagazynuKontekstow(err)
	}
	if !kontekst.Czynny {
		// Kontekst wyłączony da się uaktywnić, ale odpowiedź musi to powiedzieć, nie udawać aktywacji.
		return shared.MemoryContextActivateResponse{}, bladWskazaniaKontekstu(
			"kontekst „" + kontekst.Nazwa + "” jest wyłączony — najpierw włącz go " +
				"komendą memory.context.save (pole enabled)")
	}

	err = a.repozytorium.UaktywnijKontekstPamieci(ctx, strings.TrimSpace(z.SessionId),
		kontekst.Kod, time.Now().UnixMilli())
	if err != nil {
		return shared.MemoryContextActivateResponse{}, bladMagazynuKontekstow(err)
	}
	pozycja := kontekstPamieciKontraktu(kontekst)
	return shared.MemoryContextActivateResponse{
		Context: pozycja, Levels: pozycja.Levels, Activated: true,
	}, nil
}

func (a *adapterKontekstowPamieci) UsunKontekst(ctx context.Context,
	z shared.MemoryContextDeleteRequest) (shared.MemoryContextDeleteResponse, error) {

	if a.repozytorium == nil {
		return shared.MemoryContextDeleteResponse{}, bladZapleczaKontekstow()
	}
	usuniety, err := a.repozytorium.UsunKontekstPamieci(ctx, strings.TrimSpace(z.ContextId))
	if err != nil {
		return shared.MemoryContextDeleteResponse{}, bladMagazynuKontekstow(err)
	}
	return shared.MemoryContextDeleteResponse{Deleted: usuniety}, nil
}

func (a *adapterKontekstowPamieci) ZasadyRetencji(ctx context.Context,
	z shared.MemoryRetentionGetRequest) (shared.MemoryRetentionGetResponse, error) {

	if a.repozytorium == nil {
		return shared.MemoryRetentionGetResponse{}, bladZapleczaKontekstow()
	}
	zasieg := ""
	if z.Scope != nil {
		zasieg = string(*z.Scope)
	}
	zasady, err := a.repozytorium.ZasadyRetencjiPamieci(ctx, zasieg,
		wartoscTekstu(z.ScopeId), wartoscTekstu(z.ProfileId))
	if err != nil {
		return shared.MemoryRetentionGetResponse{}, bladMagazynuKontekstow(err)
	}
	wykaz := make([]shared.MemoryRetentionPolicy, 0, len(zasady))
	for _, zasada := range zasady {
		wykaz = append(wykaz, zasadaRetencjiKontraktu(zasada))
	}
	return shared.MemoryRetentionGetResponse{Policies: wykaz}, nil
}

func (a *adapterKontekstowPamieci) ZapiszZasadeRetencji(ctx context.Context,
	z shared.MemoryRetentionSetRequest) (shared.MemoryRetentionSetResponse, error) {

	if a.repozytorium == nil {
		return shared.MemoryRetentionSetResponse{}, bladZapleczaKontekstow()
	}
	zasieg := shared.ConfigScopeGlobal
	if z.Scope != nil && strings.TrimSpace(string(*z.Scope)) != "" {
		zasieg = string(*z.Scope)
	}
	if err := sprawdzPoziomKontekstu(shared.ConfigScope(zasieg)); err != nil {
		return shared.MemoryRetentionSetResponse{}, err
	}
	if z.TtlDays != nil && *z.TtlDays < 0 {
		return shared.MemoryRetentionSetResponse{}, bladWskazaniaKontekstu(
			"liczba dni wygasania nie może być ujemna; zero znaczy pamięć trwałą")
	}

	teraz := time.Now()
	zasada := dane.ZasadaRetencjiPamieci{
		Kod:              nowyIdentyfikator(przedrostekZasadyRetencji),
		Zasieg:           zasieg,
		ZasiegKod:        strings.TrimSpace(wartoscTekstu(z.ScopeId)),
		ProfilKod:        strings.TrimSpace(wartoscTekstu(z.ProfileId)),
		DniWygasania:     wartoscLiczby(z.TtlDays),
		WrazliweDomyslne: z.SensitiveDefault != nil && *z.SensitiveDefault,
		WzorceJSON:       zapisPolSzablonu(z.NeverStorePatterns),
		Czynna:           z.Enabled == nil || *z.Enabled,
		Zaktualizowano:   teraz.UnixMilli(),
	}
	zapisana, err := a.repozytorium.ZapiszZasadeRetencjiPamieci(ctx, zasada)
	if err != nil {
		return shared.MemoryRetentionSetResponse{}, bladMagazynuKontekstow(err)
	}

	// Wpisy zastane liczy się po zapisie: to dokładnie te, których zasada dotknie przy wygaszaniu.
	dotkniete, err := a.repozytorium.LiczbaWpisowPamieciProfilu(ctx,
		teraz.UTC().Format("2006-01-02T15:04:05.000Z"))
	if err != nil {
		return shared.MemoryRetentionSetResponse{}, bladMagazynuKontekstow(err)
	}
	return shared.MemoryRetentionSetResponse{
		Policy: zasadaRetencjiKontraktu(zapisana), AffectedEntries: dotkniete,
	}, nil
}

func zapisPoziomowKontekstu(poziomy []shared.ConfigScope) string {
	if len(poziomy) == 0 {
		return "[]"
	}
	bajty, err := json.Marshal(poziomy)
	if err != nil {
		return "[]"
	}
	return string(bajty)
}

func odczytPoziomowKontekstu(zapis string) []shared.ConfigScope {
	if strings.TrimSpace(zapis) == "" {
		return nil
	}
	var poziomy []shared.ConfigScope
	if err := json.Unmarshal([]byte(zapis), &poziomy); err != nil {
		return nil
	}
	return poziomy
}

func sprawdzPoziomKontekstu(poziom shared.ConfigScope) error {
	for _, znany := range shared.WartosciConfigScope() {
		if poziom == znany {
			return nil
		}
	}
	return bladWskazaniaKontekstu("nie znam poziomu zasięgu „" + string(poziom) + "”")
}

func kontekstPamieciKontraktu(k dane.KontekstPamieci) shared.MemoryContext {
	kontekst := shared.MemoryContext{
		Id: k.Kod, Name: k.Nazwa, Description: k.Opis,
		Levels: odczytPoziomowKontekstu(k.PoziomyJSON), EntryIds: odczytPolSzablonu(k.WpisyJSON),
		SystemPrompt: k.PromptSystemowy, Enabled: k.Czynny,
		CreatedAt: k.Utworzono, UpdatedAt: k.Zaktualizowano,
	}
	if strings.TrimSpace(k.ProfilKod) != "" {
		profil := k.ProfilKod
		kontekst.ProfileId = &profil
	}
	return kontekst
}

func zasadaRetencjiKontraktu(z dane.ZasadaRetencjiPamieci) shared.MemoryRetentionPolicy {
	zasada := shared.MemoryRetentionPolicy{
		Id: z.Kod, Scope: shared.ConfigScope(z.Zasieg), TtlDays: z.DniWygasania,
		SensitiveDefault: z.WrazliweDomyslne, NeverStorePatterns: odczytPolSzablonu(z.WzorceJSON),
		Enabled: z.Czynna, UpdatedAt: z.Zaktualizowano,
	}
	if strings.TrimSpace(z.ZasiegKod) != "" {
		kod := z.ZasiegKod
		zasada.ScopeId = &kod
	}
	if strings.TrimSpace(z.ProfilKod) != "" {
		profil := z.ProfilKod
		zasada.ProfileId = &profil
	}
	return zasada
}

func bladZapleczaKontekstow() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"konteksty pamięci: serwer nie ma wpiętego magazynu — naprawa: podpiąć "+
			"repozytorium kontekstów pamięci przy składaniu serwera"))
}

func bladWskazaniaKontekstu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"konteksty pamięci: "+powod))
}

func bladNieznanegoKontekstuPamieci(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"konteksty pamięci: nie ma kontekstu o identyfikatorze "+strings.TrimSpace(kod)))
}

func bladMagazynuKontekstow(err error) error {
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
			"konteksty pamięci: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"konteksty pamięci: "+err.Error()))
}
