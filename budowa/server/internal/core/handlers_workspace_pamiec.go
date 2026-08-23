// Odpowiedzialność pliku: wpięcie pięciu komend rodziny `memory.*`.
//
// Plik stoi osobno od `handlers_workspace.go`, który wpina sześć komend obszaru
// `workspace.*`. Rodzina `memory.*` weszła do kontraktu osobno i osobno się
// wpina, choć jedzie na tej samej maszynerii pamięci. Port poniżej jest
// rozszerzeniem portu PrzestrzenRobocza, nie drugim portem: pamięć ma w rdzeniu
// jednego właściciela.
//
// Jedna zmiana wpisu rozgłasza dwa zdarzenia i nie jest to powtórzenie:
// `workspace.project.changed` niesie projekt i mówi oknu Workspace, że coś się
// w nim ruszyło, a `memory.changed` niesie sam wpis wraz z rodzajem zmiany
// i zasila okno Context Memory bez odpytywania `memory.list`. To dwa różne byty
// tej samej czynności — tym samym wzorem, co `queue.changed` obok
// `automation.execution.status`.
//
// `memory.toggle` nie rozgłasza niczego: przestawia konfigurację karty sesji,
// nie stan wpisu ani projektu, a zdarzenia dla tego bytu kontrakt nie ma.
//
// Rodzaj zmiany idzie za tym, co się stało. `memory.set` bez `entryId` zakłada
// wpis (`created`), ze wskazaniem — zmienia go (`updated`); `memory.detach`
// zwęża zasięg wpisu, który żyje dalej (`updated`); `memory.delete` kasuje
// (`deleted`). Odesłanie wszystkiego jako `updated` kazałoby klientowi zgadywać,
// czy wpis dopisać do wykazu, czy z niego zdjąć.
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
	// PrzestawWylaczeniePamieci obsługuje `memory.disable.set` — założenie
	// wyłączenia i jego zniesienie idą jedną komendą, bo odwracalność jednym
	// ruchem jest wymogiem produktu.
	PrzestawWylaczeniePamieci(ctx context.Context,
		z shared.MemoryDisableSetRequest) (shared.MemoryDisableSetResponse, error)

	// WpisPamieciKontraktu oddaje wpis w kształcie kontraktu. Potrzebują go oba
	// rozgłoszenia po `memory.delete`, którego żądanie niesie sam identyfikator:
	// projekt dla `workspace.project.changed`, cały wpis dla `memory.changed`.
	WpisPamieciKontraktu(ctx context.Context, identyfikator string) (shared.WorkspaceMemoryEntry, error)
}

// zarejestrujPamiec wpina pięć komend rodziny `memory.*`.
func zarejestrujPamiec(r *Rejestr, m PamiecPrzestrzeni, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandMemoryDelete,
		obsluz(func(ctx context.Context, z shared.MemoryDeleteRequest) (shared.MemoryDeleteResponse, error) {
			// Wpis odczytuje się przed usunięciem — po nim już go nie ma
			// i nie ma jak powiedzieć ani czyja pamięć się zmieniła, ani co
			// zniknęło. Nieudany odczyt kończy wyłącznie rozgłoszenie, nie
			// komendę.
			wpis, _ := m.WpisPamieciKontraktu(ctx, z.EntryId)
			odpowiedz, err := m.UsunWpisPamieciKomenda(ctx, z)
			if err == nil && wpis.Id != "" {
				rozglosProjekt(ctx, m, e, wpis.ProjectId)
				e.wpisPamieci(shared.ChangeKindDeleted, wpis)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandMemoryDetach,
		obsluz(func(ctx context.Context, z shared.MemoryDetachRequest) (shared.MemoryDetachResponse, error) {
			odpowiedz, err := m.OdepnijWpisPamieci(ctx, z)
			// Bez odpięcia nie ma zmiany, więc nie ma czego rozgłaszać.
			if err == nil && odpowiedz.Detached {
				rozglosProjekt(ctx, m, e, odpowiedz.Entry.ProjectId)
				e.wpisPamieci(shared.ChangeKindUpdated, odpowiedz.Entry)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandMemoryList, obsluz(m.WpisyPamieciZasiegu))

	r.Zarejestruj(shared.CommandMemorySet,
		obsluz(func(ctx context.Context, z shared.MemorySetRequest) (shared.MemorySetResponse, error) {
			odpowiedz, err := m.ZapiszPamiec(ctx, z)
			if err == nil {
				rozglosProjekt(ctx, m, e, odpowiedz.Entry.ProjectId)
				e.wpisPamieci(rodzajZapisuPamieci(z.EntryId), odpowiedz.Entry)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandMemoryToggle, obsluz(m.PrzestawPamiecSesji))

	r.Zarejestruj(shared.CommandMemoryDisableList, obsluz(m.WylaczeniaPamieciZasiegu))

	r.Zarejestruj(shared.CommandMemoryDisableSet,
		obsluz(func(ctx context.Context, z shared.MemoryDisableSetRequest) (shared.MemoryDisableSetResponse, error) {
			odpowiedz, err := m.PrzestawWylaczeniePamieci(ctx, z)
			// Bez zmiany wykazu nie ma czego rozgłaszać: wyłączenie powtórzone
			// nie jest zdarzeniem. Rodzaj zmiany idzie za żądaniem — wyłączenie
			// zakłada wiersz, zniesienie go kasuje — żeby powłoka nie musiała
			// zgadywać, czy pozycję dopisać do wykazu, czy z niego zdjąć.
			if err == nil && odpowiedz.Changed {
				e.wylaczeniePamieci(rodzajZmianyWylaczenia(z.Disabled), odpowiedz.Disable)
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
func (e *emiter) wylaczeniePamieci(zmiana shared.ChangeKind, wylaczenie shared.MemoryDisable) {
	e.wyslij(shared.EventMemoryDisableChanged, "",
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
func (e *emiter) wpisPamieci(zmiana shared.ChangeKind, wpis shared.WorkspaceMemoryEntry) {
	e.wyslij(shared.EventMemoryChanged, "", shared.MemoryChangedEvent{Change: zmiana, Entry: wpis})
}
