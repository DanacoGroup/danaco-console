// Dopełnia rodzinę memory.* obsługą nazwanych kontekstów pamięci oraz zasad
// retencji i wygaszania: kontekst jest zestawem wskazań na wpisy, nie ich
// właścicielem, więc usunięcie kontekstu nie kasuje żadnego wpisu pamięci.
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

// adapterKontekstowPamieci wypełnia port KontekstyPamieci, wiążąc obsługę komend
// memory.context oraz memory.retention z repozytorium kontekstów pamięci
// przechowywanym w magazynie danych.
type adapterKontekstowPamieci struct {
	repozytorium dane.RepozytoriumKontekstowPamieci
}

// nowyAdapterKontekstowPamieci wiąże port KontekstyPamieci z podanym repozytorium
// magazynu, zwracając gotowy adapter dla montażu rdzenia bez dalszej konfiguracji.
func nowyAdapterKontekstowPamieci(
	repozytorium dane.RepozytoriumKontekstowPamieci) *adapterKontekstowPamieci {

	return &adapterKontekstowPamieci{repozytorium: repozytorium}
}

// WykazKontekstow obsługuje komendę memory.context.list: zwraca konteksty profilu
// wraz z identyfikatorem kontekstu czynnego dla wskazanej karty sesji.
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

// ZapiszKontekst obsługuje komendę memory.context.save: zakłada nowy kontekst albo
// nadpisuje istniejący, sprawdzając wcześniej poprawność każdego wskazanego poziomu zasięgu.
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

// UaktywnijKontekst obsługuje komendę memory.context.activate: ustawia kontekst
// czynny wskazanej karty sesji, odmawiając dla kontekstu wyłączonego albo nieznanego.
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

// UsunKontekst obsługuje komendę memory.context.delete, usuwając kontekst o podanym
// identyfikatorze i zwracając informację, czy wiersz w magazynie w ogóle istniał.
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

// ZasadyRetencji obsługuje komendę memory.retention.get, zwracając zasady retencji
// dopasowane do wskazanego zasięgu, profilu i identyfikatora zasięgu.
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

// ZapiszZasadeRetencji obsługuje komendę memory.retention.set: zapisuje nową zasadę
// retencji i liczy wpisy pamięci, których zasada dotknie przy najbliższym wygaszaniu.
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

// zapisPoziomowKontekstu składa poziomy zasięgu kontekstu w zapis strukturalny kolumny
// magazynu, oddając pustą tablicę tekstową, gdy poziomów nie podano.
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

// odczytPoziomowKontekstu rozkłada zapis strukturalny kolumny magazynu z powrotem na
// poziomy zasięgu, oddając pustą wartość dla zapisu pustego albo niepoprawnego.
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

// sprawdzPoziomKontekstu odbija poziom zasięgu spoza wyliczenia kontraktu, porównując
// podaną wartość z pełnym wykazem poziomów znanych kontraktowi.
func sprawdzPoziomKontekstu(poziom shared.ConfigScope) error {
	for _, znany := range shared.WartosciConfigScope() {
		if poziom == znany {
			return nil
		}
	}
	return bladWskazaniaKontekstu("nie znam poziomu zasięgu „" + string(poziom) + "”")
}

// kontekstPamieciKontraktu przekłada wiersz magazynu na kontekst kontraktu, ustawiając
// identyfikator profilu tylko wtedy, gdy wiersz go niesie.
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

// zasadaRetencjiKontraktu przekłada wiersz magazynu na zasadę retencji kontraktu,
// ustawiając identyfikator zasięgu i profilu tylko wtedy, gdy wiersz je niesie.
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

// bladZapleczaKontekstow nazywa brak magazynu kontekstów po stronie rdzenia i wskazuje
// naprawę: podpięcie repozytorium kontekstów pamięci przy składaniu rdzenia.
func bladZapleczaKontekstow() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"konteksty pamięci: serwer nie ma wpiętego magazynu — naprawa: podpiąć "+
			"repozytorium kontekstów pamięci przy składaniu serwera"))
}

// bladWskazaniaKontekstu nazywa niepoprawne żądanie komendy kontekstów pamięci,
// niosąc w treści błędu powód odmowy podany przez wywołanie.
func bladWskazaniaKontekstu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"konteksty pamięci: "+powod))
}

// bladNieznanegoKontekstuPamieci nazywa wskazanie kontekstu o identyfikatorze, którego
// magazyn kontekstów pamięci nie zawiera.
func bladNieznanegoKontekstuPamieci(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"konteksty pamięci: nie ma kontekstu o identyfikatorze "+strings.TrimSpace(kod)))
}

// bladMagazynuKontekstow nazywa niepowodzenie zapisu albo odczytu w magazynie
// kontekstów pamięci, niosąc w treści błędu przyczynę zgłoszoną przez magazyn.
func bladMagazynuKontekstow(err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"konteksty pamięci: "+err.Error()))
}
