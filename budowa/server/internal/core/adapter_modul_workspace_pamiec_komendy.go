// Rodzina komend memory.* obsługuje wpis pamięci w tabeli wpis_pamieci_projektu
// oraz konfigurację dostępu karty sesji do poziomów pamięci w tabeli
// konfiguracja_pamieci_sesji.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// repozytoriumWpisowPamieci rozszerza repozytorium przestrzeni roboczej
// o czynności, których okno Context Memory nie potrzebuje: odczyt jednego
// wpisu, jego usunięcie i przestawienie zasięgu.
type repozytoriumWpisowPamieci interface {
	WpisPamieciPoIdentyfikatorze(ctx context.Context, identyfikator string) (dane.WpisPamieciProjektu, bool, error)
	UsunWpisPamieci(ctx context.Context, identyfikator string) (bool, error)
	PrzestawZasiegWpisuPamieci(ctx context.Context, identyfikator string,
		poziom shared.ConfigScope, kluczZasiegu string) (dane.WpisPamieciProjektu, error)
}

// adapterPamieciPrzestrzeni wypełnia port PamiecPrzestrzeni: adapter modułu
// Workspace rozszerzony o pamięć sesji i wykaz kart sesji.
type adapterPamieciPrzestrzeni struct {
	*adapterPrzestrzeniRoboczej
	// wpisy zerowe daje odmowę z kodem internal_error.
	wpisy  repozytoriumWpisowPamieci
	pamiec dane.RepozytoriumPamieci
	sesje  dane.RepozytoriumSesji
	// wylaczenia jest magazynem wyłączeń pamięci; zerowe daje odmowę także
	// w memory.list.
	wylaczenia dane.RepozytoriumWylaczenPamieci
}

// ZPamiecia rozszerza adapter modułu Workspace o rodzinę `memory.*`. Pamięć
// karty sesji i wykaz kart to dwa repozytoria, których okno Context Memory nie
// potrzebuje — `memory.toggle` bez nich nie ma czego przestawić.
func (a *adapterPrzestrzeniRoboczej) ZPamiecia(pamiec dane.RepozytoriumPamieci,
	sesje dane.RepozytoriumSesji) *adapterPamieciPrzestrzeni {

	rozszerzony := &adapterPamieciPrzestrzeni{adapterPrzestrzeniRoboczej: a, pamiec: pamiec, sesje: sesje}
	if wpisy, ok := a.repozytorium.(repozytoriumWpisowPamieci); ok {
		rozszerzony.wpisy = wpisy
	}
	return rozszerzony
}

// ── memory.delete ────────────────────────────────────────────────────────────

// UsunWpisPamieciKomenda usuwa wpis pamięci wskazany identyfikatorem.
//
// Wpis nieistniejący daje odmowę z kodem `not_found` i nazwą wpisu, a nie
// odpowiedź `deleted: false`.
func (a *adapterPamieciPrzestrzeni) UsunWpisPamieciKomenda(ctx context.Context,
	z shared.MemoryDeleteRequest) (shared.MemoryDeleteResponse, error) {

	wpis, err := a.wpisZadania(ctx, z.EntryId)
	if err != nil {
		return shared.MemoryDeleteResponse{}, err
	}
	usuniety, err := a.wpisy.UsunWpisPamieci(ctx, wpis.Identyfikator)
	if err != nil {
		return shared.MemoryDeleteResponse{}, err
	}
	if !usuniety {
		return shared.MemoryDeleteResponse{}, bladBrakuWpisuPamieci(z.EntryId)
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, wpis.ProjektID); err != nil {
		return shared.MemoryDeleteResponse{}, err
	}
	return shared.MemoryDeleteResponse{Deleted: true}, nil
}

// WpisPamieciKontraktu zwraca wskazany wpis w kształcie kontraktu. Służy
// dwóm rozgłoszeniom po memory.delete, którym odczyt musi poprzedzić usunięcie.
func (a *adapterPamieciPrzestrzeni) WpisPamieciKontraktu(ctx context.Context,
	identyfikator string) (shared.WorkspaceMemoryEntry, error) {

	wpis, err := a.wpisZadania(ctx, identyfikator)
	if err != nil {
		return shared.WorkspaceMemoryEntry{}, err
	}
	return wpisPamieciKontraktu(wpis), nil
}

// ── memory.detach ────────────────────────────────────────────────────────────

// OdepnijWpisPamieci zwęża zasięg wpisu pamięci ze wspólnego z powrotem do
// projektu, który go niesie, zostawiając treść wpisu nietkniętą. Odpięcia od
// projektu, który wpisu nie niesie, komenda nie wykonuje i odmawia z kodem
// conflict.
func (a *adapterPamieciPrzestrzeni) OdepnijWpisPamieci(ctx context.Context,
	z shared.MemoryDetachRequest) (shared.MemoryDetachResponse, error) {

	wpis, err := a.wpisZadania(ctx, z.EntryId)
	if err != nil {
		return shared.MemoryDetachResponse{}, err
	}
	if err := a.sprawdzProjektOdpiecia(ctx, wpis, z); err != nil {
		return shared.MemoryDetachResponse{}, err
	}
	// Zasięg równy projektowi albo węższy nie ma czego zdejmować.
	if !zasiegSzerszyNizProjekt(wpis.Poziom) {
		return shared.MemoryDetachResponse{Entry: wpisPamieciKontraktu(wpis), Detached: false}, nil
	}
	odpiety, err := a.wpisy.PrzestawZasiegWpisuPamieci(ctx, wpis.Identyfikator,
		shared.ConfigScopeProject, wpis.ProjektKod)
	if err != nil {
		return shared.MemoryDetachResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, wpis.ProjektID); err != nil {
		return shared.MemoryDetachResponse{}, err
	}
	return shared.MemoryDetachResponse{Entry: wpisPamieciKontraktu(odpiety), Detached: true}, nil
}

// sprawdzProjektOdpiecia odmawia odpięcia wpisu od projektu, który go nie
// niesie. Schemat nie zna wyłączenia ustalenia wspólnego w pojedynczym
// projekcie, więc takiego żądania nie da się wykonać.
func (a *adapterPamieciPrzestrzeni) sprawdzProjektOdpiecia(ctx context.Context,
	wpis dane.WpisPamieciProjektu, z shared.MemoryDetachRequest) error {

	wskazany, err := a.kodProjektuZadania(ctx, z.ProjectId, z.SessionId)
	if err != nil || wskazany == "" {
		return err
	}
	if wskazany != wpis.ProjektKod {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
			"moduł Workspace: wpis pamięci "+wpis.Identyfikator+" należy do projektu "+
				wpis.ProjektKod+"; odpięcie zwęża zasięg samego wpisu, więc nie ma jak odpiąć go "+
				"od projektu "+wskazany+". Wstrzymanie tego wpisu w projekcie "+wskazany+
				" robi memory.disable.set — bez ruszania treści"))
	}
	return nil
}

// zasiegSzerszyNizProjekt mówi, czy wpis obowiązuje ponad swoim projektem.
// Poziomy szersze wymienione są wprost, bo pierwszeństwo poziomów prowadzi baza.
func zasiegSzerszyNizProjekt(poziom shared.ConfigScope) bool {
	switch poziom {
	case shared.ConfigScopeGlobal, shared.ConfigScopeEnvironment,
		shared.ConfigScopeModule, shared.ConfigScopeModulePair:
		return true
	default:
		return false
	}
}

// ── memory.list ──────────────────────────────────────────────────────────────

// WpisyPamieciZasiegu zwraca wpisy pamięci widoczne w zasięgu żądania. Wpis
// wyłączony nie wchodzi do wykazu czynnych, lecz idzie osobnym wykazem
// disabledEntries wraz z zasięgiem wyłączenia, a jego treść zostaje nietknięta.
func (a *adapterPamieciPrzestrzeni) WpisyPamieciZasiegu(ctx context.Context,
	z shared.MemoryListRequest) (shared.MemoryListResponse, error) {

	projekt, err := a.projektZadania(ctx, z.ProjectId, z.SessionId)
	if err != nil {
		return shared.MemoryListResponse{}, err
	}
	wylaczone, err := a.sitoWylaczen(ctx)
	if err != nil {
		return shared.MemoryListResponse{}, err
	}
	granica := 0
	if z.Limit != nil {
		granica = *z.Limit
	}
	// Przy sicie granica idzie na koniec, nie do bazy.
	sito := sitoZasieguWpisow(z.Scope, z.ScopeId)
	granicaBazy := granica
	if sito != nil {
		granicaBazy = 0
	}
	wiersze, err := a.repozytorium.WpisyPamieci(ctx, projekt.ID, granicaBazy)
	if err != nil {
		return shared.MemoryListResponse{}, err
	}
	if z.IncludeShared != nil && *z.IncludeShared {
		wspolne, err := a.repozytorium.WpisyPamieciWspoldzielone(ctx, projekt.ID, granicaBazy)
		if err != nil {
			return shared.MemoryListResponse{}, err
		}
		wiersze = append(wiersze, wspolne...)
	}
	wpisy := make([]shared.WorkspaceMemoryEntry, 0, len(wiersze))
	wstrzymane := []shared.MemoryDisabledEntry{}
	for _, wiersz := range wiersze {
		if sito != nil && !sito(wiersz) {
			continue
		}
		// Wyłączenie odsiewa się przed granicą; wpis wstrzymany musi zostać
		// nazwany, nie zniknąć bez śladu.
		if wylaczone != nil {
			if opis, wylaczony := wylaczone(wiersz); wylaczony {
				wstrzymane = append(wstrzymane, opis)
				continue
			}
		}
		if granica > 0 && len(wpisy) == granica {
			break
		}
		wpisy = append(wpisy, wpisPamieciKontraktu(wiersz))
	}
	return shared.MemoryListResponse{Entries: wpisy, DisabledEntries: wstrzymane}, nil
}

// sitoZasieguWpisow składa sito wpisów z pól opcjonalnych żądania. Zwraca nil,
// gdy żądanie o żadne zawężenie nie prosi — wtedy wykaz idzie w całości.
func sitoZasieguWpisow(poziom *shared.ConfigScope,
	byt *string) func(dane.WpisPamieciProjektu) bool {

	zadanyPoziom := ""
	if poziom != nil {
		zadanyPoziom = string(*poziom)
	}
	zadanyByt := ""
	if byt != nil {
		zadanyByt = *byt
	}
	if zadanyPoziom == "" && zadanyByt == "" {
		return nil
	}
	return func(wpis dane.WpisPamieciProjektu) bool {
		if zadanyPoziom != "" && string(wpis.Poziom) != zadanyPoziom {
			return false
		}
		return zadanyByt == "" || wpis.KluczZasiegu == zadanyByt
	}
}

// ── memory.set ───────────────────────────────────────────────────────────────

// ZapiszPamiec zapisuje ustalenie w pamięci. Puste entryId zakłada wpis,
// podane zmienia istniejący. Wskazany byt zasięgu jest brany wprost; bez
// wskazania bytem poziomu projektu jest projekt, a poziomu szerszego — pustka.
func (a *adapterPamieciPrzestrzeni) ZapiszPamiec(ctx context.Context,
	z shared.MemorySetRequest) (shared.MemorySetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, wartoscTekstu(z.ProjectId))
	if err != nil {
		return shared.MemorySetResponse{}, err
	}
	if z.Content == "" {
		return shared.MemorySetResponse{}, bladProjektu("wpis pamięci bez treści")
	}
	wpis := dane.WpisPamieciProjektu{
		ProjektID:     projekt.ID,
		Identyfikator: identyfikatorWpisuPamieci(z.EntryId),
		Tresc:         z.Content,
		Przypiety:     z.Pinned != nil && *z.Pinned,
		Pochodzenie:   pochodzenieWpisu(z.Origin),
		Poziom:        zasiegWpisu(z.Scope),
		KluczZasiegu:  wartoscTekstu(z.ScopeId),
	}
	if wpis.KluczZasiegu == "" && wpis.Poziom == shared.ConfigScopeProject {
		wpis.KluczZasiegu = projekt.Kod
	}
	zapisany, err := a.repozytorium.ZapiszWpisPamieci(ctx, wpis)
	if err != nil {
		return shared.MemorySetResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.MemorySetResponse{}, err
	}
	return shared.MemorySetResponse{Entry: wpisPamieciKontraktu(zapisany)}, nil
}

// ── wspólne ustalenia żądań ──────────────────────────────────────────────────

// wpisZadania odczytuje wpis wskazany identyfikatorem. Brak wskazania jest
// żądaniem niezgodnym z kontraktem, brak wpisu — bytem, którego nie ma.
func (a *adapterPamieciPrzestrzeni) wpisZadania(ctx context.Context,
	identyfikator string) (dane.WpisPamieciProjektu, error) {

	if identyfikator == "" {
		return dane.WpisPamieciProjektu{}, bladProjektu("komenda pamięci bez wskazania wpisu")
	}
	if a.wpisy == nil {
		return dane.WpisPamieciProjektu{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"moduł Workspace: repozytorium nie niesie czynności wpisu pamięci"))
	}
	wpis, jest, err := a.wpisy.WpisPamieciPoIdentyfikatorze(ctx, identyfikator)
	if err != nil {
		return dane.WpisPamieciProjektu{}, err
	}
	if !jest {
		return dane.WpisPamieciProjektu{}, bladBrakuWpisuPamieci(identyfikator)
	}
	return wpis, nil
}

// projektZadania ustala projekt komendy pamięci: ze wskazania wprost albo
// z projektu karty sesji, gdy projektu nie wskazano.
func (a *adapterPamieciPrzestrzeni) projektZadania(ctx context.Context,
	idProjektu, idSesji *string) (dane.Projekt, error) {

	kod, err := a.kodProjektuZadania(ctx, idProjektu, idSesji)
	if err != nil {
		return dane.Projekt{}, err
	}
	return a.projektDlaZapisu(ctx, kod)
}

// kodProjektuZadania zwraca kod projektu wskazanego żądaniem. Pusty wynik bez
// błędu znaczy, że żądanie nie wskazało ani projektu, ani karty sesji.
func (a *adapterPamieciPrzestrzeni) kodProjektuZadania(ctx context.Context,
	idProjektu, idSesji *string) (string, error) {

	if kod := wartoscTekstu(idProjektu); kod != "" {
		return kod, nil
	}
	nazwaSesji := wartoscTekstu(idSesji)
	if nazwaSesji == "" {
		return "", nil
	}
	if a.sesje == nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Workspace: wykaz kart sesji niewpięty, projektu karty nie da się ustalić"))
	}
	sesja, err := a.sesje.PoIdentyfikatorze(ctx, nazwaSesji)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Workspace: karta sesji "+nazwaSesji+" nie istnieje"))
	}
	if err != nil {
		return "", err
	}
	kod := tekstLubPusty(sesja.Projekt)
	if kod == "" {
		return "", bladProjektu("karta sesji " + nazwaSesji + " nie pracuje w żadnym projekcie")
	}
	return kod, nil
}

// bladBrakuWpisuPamieci składa odmowę z kodem not_found wskazującą, że wpis
// pamięci o podanym identyfikatorze w repozytorium nie istnieje.
func bladBrakuWpisuPamieci(identyfikator string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Workspace: wpis pamięci "+identyfikator+" nie istnieje"))
}
