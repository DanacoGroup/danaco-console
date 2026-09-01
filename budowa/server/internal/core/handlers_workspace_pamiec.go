// Plik wpina pięć komend rodziny memory.* jako rozszerzenie portu PrzestrzenRobocza, bo pamięć
// ma w rdzeniu jednego właściciela wspólnego z obszarem workspace.*.
package core

import (
	"context"

	"danacoconsole/shared"
)

// PamiecPrzestrzeni jest portem rodziny `memory.*` — portem modułu
// Workspace rozszerzonym o pięć komend pamięci.
type PamiecPrzestrzeni interface {
	PrzestrzenRobocza

	// UsunWpisPamieciKomenda obsługuje `memory.delete`.
	UsunWpisPamieciKomenda(ctx context.Context, z shared.MemoryDeleteRequest) (shared.MemoryDeleteResponse, error)
	// OdepnijWpisPamieci obsługuje `memory.detach`.
	OdepnijWpisPamieci(ctx context.Context, z shared.MemoryDetachRequest) (shared.MemoryDetachResponse, error)
	// WpisyPamieciZasiegu obsługuje `memory.list`.
	WpisyPamieciZasiegu(ctx context.Context, z shared.MemoryListRequest) (shared.MemoryListResponse, error)
	// ZapiszPamiec obsługuje `memory.set`.
	ZapiszPamiec(ctx context.Context, z shared.MemorySetRequest) (shared.MemorySetResponse, error)
	// PrzestawPamiecSesji obsługuje `memory.toggle`.
	PrzestawPamiecSesji(ctx context.Context, z shared.MemoryToggleRequest) (shared.MemoryToggleResponse, error)
	// WylaczeniaPamieciZasiegu obsługuje `memory.disable.list`.
	WylaczeniaPamieciZasiegu(ctx context.Context,
		z shared.MemoryDisableListRequest) (shared.MemoryDisableListResponse, error)
	// PrzestawWylaczeniePamieci obsługuje założenie i zniesienie wyłączenia jedną komendą odwracalną.
	PrzestawWylaczeniePamieci(ctx context.Context,
		z shared.MemoryDisableSetRequest) (shared.MemoryDisableSetResponse, error)

	// WpisPamieciKontraktu oddaje wpis w kształcie kontraktu, potrzebny obu rozgłoszeniom po usunięciu go.
	WpisPamieciKontraktu(ctx context.Context, identyfikator string) (shared.WorkspaceMemoryEntry, error)
}

// zarejestrujPamiec wpina pięć komend rodziny memory.* w rejestrze rdzenia tej platformy konta użytkownika.
func zarejestrujPamiec(r *Rejestr, m PamiecPrzestrzeni, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandMemoryDelete,
		obsluz(func(ctx context.Context, z shared.MemoryDeleteRequest) (shared.MemoryDeleteResponse, error) {
			// Wpis odczytuje się przed usunięciem, bo po nim go nie ma; nieudany odczyt kończy rozgłoszenie.
			wpis, _ := m.WpisPamieciKontraktu(ctx, z.EntryId)
			odpowiedz, err := m.UsunWpisPamieciKomenda(ctx, z)
			if err == nil && wpis.Id != "" {
				rozglosProjekt(ctx, m, e, wpis.ProjectId)
				e.wpisPamieci(ctx, shared.ChangeKindDeleted, wpis)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandMemoryDetach,
		obsluz(func(ctx context.Context, z shared.MemoryDetachRequest) (shared.MemoryDetachResponse, error) {
			odpowiedz, err := m.OdepnijWpisPamieci(ctx, z)
			// Bez odpięcia nie ma zmiany, więc nie ma czego rozgłaszać.
			if err == nil && odpowiedz.Detached {
				rozglosProjekt(ctx, m, e, odpowiedz.Entry.ProjectId)
				e.wpisPamieci(ctx, shared.ChangeKindUpdated, odpowiedz.Entry)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandMemoryList, obsluz(m.WpisyPamieciZasiegu))

	r.Zarejestruj(shared.CommandMemorySet,
		obsluz(func(ctx context.Context, z shared.MemorySetRequest) (shared.MemorySetResponse, error) {
			odpowiedz, err := m.ZapiszPamiec(ctx, z)
			if err == nil {
				rozglosProjekt(ctx, m, e, odpowiedz.Entry.ProjectId)
				e.wpisPamieci(ctx, rodzajZapisuPamieci(z.EntryId), odpowiedz.Entry)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandMemoryToggle, obsluz(m.PrzestawPamiecSesji))

	r.Zarejestruj(shared.CommandMemoryDisableList, obsluz(m.WylaczeniaPamieciZasiegu))

	r.Zarejestruj(shared.CommandMemoryDisableSet,
		obsluz(func(ctx context.Context, z shared.MemoryDisableSetRequest) (shared.MemoryDisableSetResponse, error) {
			odpowiedz, err := m.PrzestawWylaczeniePamieci(ctx, z)
			// Bez zmiany wykazu nie ma czego rozgłaszać; wyłączenie zakłada wiersz, zniesienie go kasuje.
			if err == nil && odpowiedz.Changed {
				e.wylaczeniePamieci(ctx, rodzajZmianyWylaczenia(z.Disabled), odpowiedz.Disable)
			}
			return odpowiedz, err
		}))
}

// rodzajZmianyWylaczenia rozstrzyga, czy wyłączenie powstało, czy zniknęło.
// Rozstrzyga to samo żądanie: `disabled` równe prawdzie zakłada wyłączenie,
// równe fałszowi je znosi.
func rodzajZmianyWylaczenia(wylaczone bool) shared.ChangeKind {
	if wylaczone {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindDeleted
}

// wylaczeniePamieci rozgłasza `memory.disable.changed`. Zdarzenie idzie bez
// wskazania sesji, tak samo jak `memory.changed`: wyłączenie obowiązuje zasięg,
// który sam nazywa, a nie kartę, z której przyszło.
func (e *emiter) wylaczeniePamieci(ctx context.Context, zmiana shared.ChangeKind, wylaczenie shared.MemoryDisable) {
	e.wyslijDoKonta(ctx, shared.EventMemoryDisableChanged, "",
		shared.MemoryDisableChangedEvent{Change: zmiana, Disable: wylaczenie})
}

// rodzajZapisuPamieci rozstrzyga, czy `memory.set` założył wpis, czy zmienił
// zastany. Rozstrzyga to samo żądanie: puste `entryId` znaczy wpis nowy — tak
// mówi opis komendy w kontrakcie i tak zachowuje się adapter.
func rodzajZapisuPamieci(idWpisu *string) shared.ChangeKind {
	if idWpisu == nil || *idWpisu == "" {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindUpdated
}

// wpisPamieci rozgłasza `memory.changed`. Wpis pamięci należy do projektu, nie
// do karty sesji — projekt bywa otwarty w wielu kartach naraz — więc zdarzenie
// idzie bez wskazania sesji, tak samo jak `workspace.project.changed`.
func (e *emiter) wpisPamieci(ctx context.Context, zmiana shared.ChangeKind, wpis shared.WorkspaceMemoryEntry) {
	e.wyslijDoKonta(ctx, shared.EventMemoryChanged, "", shared.MemoryChangedEvent{Change: zmiana, Entry: wpis})
}
