// Odpowiedzialność pliku: pamięć projektu — okno Context Memory modułu
// Workspace. Propozycja modelu to wpis o pochodzeniu `model`.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekWpisuPamieci znakuje identyfikator wpisu nadany przez rdzeń,
// przy zapisie nowego wpisu pamięci.
const przedrostekWpisuPamieci = "pam-"

// ZapiszWpisPamieci zakłada wpis pamięci albo zmienia wpis wskazany
// identyfikatorem, z pochodzeniem podanym w żądaniu.
func (a *adapterPrzestrzeniRoboczej) ZapiszWpisPamieci(ctx context.Context,
	z shared.WorkspaceContextSetRequest) (shared.WorkspaceContextSetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceContextSetResponse{}, err
	}
	if z.Content == "" {
		return shared.WorkspaceContextSetResponse{}, bladProjektu("wpis pamięci bez treści")
	}
	wpis := dane.WpisPamieciProjektu{
		ProjektID:     projekt.ID,
		Identyfikator: identyfikatorWpisuPamieci(z.EntryId),
		Tresc:         z.Content,
		Przypiety:     z.Pinned != nil && *z.Pinned,
		Pochodzenie:   pochodzenieWpisu(z.Origin),
		Poziom:        zasiegWpisu(z.Scope),
		KluczZasiegu:  projekt.Kod,
	}
	// Byt zasięgu ma sens tylko dla poziomu projektu; poziom szerszy
	// obowiązuje ponad projektami.
	if wpis.Poziom != shared.ConfigScopeProject {
		wpis.KluczZasiegu = ""
	}
	zapisany, err := a.repozytorium.ZapiszWpisPamieci(ctx, wpis)
	if err != nil {
		return shared.WorkspaceContextSetResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.WorkspaceContextSetResponse{}, err
	}
	return shared.WorkspaceContextSetResponse{Entry: wpisPamieciKontraktu(zapisany)}, nil
}

// WpisyPamieci zwraca pamięć projektu, a na żądanie także ustalenia wspólne
// zapisane na poziomie szerszym niż projekt.
func (a *adapterPrzestrzeniRoboczej) WpisyPamieci(ctx context.Context,
	z shared.WorkspaceContextGetRequest) (shared.WorkspaceContextGetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceContextGetResponse{}, err
	}
	granica := 0
	if z.Limit != nil {
		granica = *z.Limit
	}
	wiersze, err := a.repozytorium.WpisyPamieci(ctx, projekt.ID, granica)
	if err != nil {
		return shared.WorkspaceContextGetResponse{}, err
	}
	if z.IncludeShared != nil && *z.IncludeShared {
		wspolne, err := a.repozytorium.WpisyPamieciWspoldzielone(ctx, projekt.ID, granica)
		if err != nil {
			return shared.WorkspaceContextGetResponse{}, err
		}
		wiersze = append(wiersze, wspolne...)
	}
	wpisy := make([]shared.WorkspaceMemoryEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wpisy = append(wpisy, wpisPamieciKontraktu(wiersz))
	}
	return shared.WorkspaceContextGetResponse{Entries: wpisy}, nil
}

// identyfikatorWpisuPamieci zwraca identyfikator wpisu zmienianego albo nadaje
// nowy. Brak wskazania znaczy wpis nowy, nie błąd żądania.
func identyfikatorWpisuPamieci(wskazany *string) string {
	if wskazany != nil && *wskazany != "" {
		return *wskazany
	}
	return nowyIdentyfikator(przedrostekWpisuPamieci)
}

// pochodzenieWpisu przekłada pole opcjonalne żądania. Brak wskazania znaczy
// ustalenie Operatora — propozycję modelu zgłasza wprost pochodzenie `model`.
func pochodzenieWpisu(wskazane *shared.MemoryEntryOrigin) shared.MemoryEntryOrigin {
	if wskazane != nil && *wskazane != "" {
		return *wskazane
	}
	return shared.MemoryEntryOriginOperator
}

// zasiegWpisu przekłada pole opcjonalne żądania. Brak wskazania znaczy zasięg
// projektu — pamięć jest domyślnie odrębna dla każdego projektu.
func zasiegWpisu(wskazany *shared.ConfigScope) shared.ConfigScope {
	if wskazany != nil && *wskazany != "" {
		return *wskazany
	}
	return shared.ConfigScopeProject
}

// wpisPamieciKontraktu przekłada wiersz pamięci z bazy danych na byt
// kontraktu zwracany wołającemu oknu.
func wpisPamieciKontraktu(w dane.WpisPamieciProjektu) shared.WorkspaceMemoryEntry {
	wpis := shared.WorkspaceMemoryEntry{
		Id:        w.Identyfikator,
		ProjectId: w.ProjektKod,
		Content:   w.Tresc,
		Pinned:    &w.Przypiety,
		Origin:    w.Pochodzenie,
		Scope:     w.Poziom,
		CreatedAt: chwilaBazy(w.Utworzono),
		UpdatedAt: chwilaBazy(w.Zaktualizowano),
	}
	if w.KluczZasiegu != "" {
		byt := w.KluczZasiegu
		wpis.ScopeId = &byt
	}
	return wpis
}
