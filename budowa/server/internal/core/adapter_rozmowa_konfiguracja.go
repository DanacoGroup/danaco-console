// Doprowadzenie obowiązującej konfiguracji sesji do zbudowanego wywołania modelu.
// Jednolity model konfiguracji `shared.SessionConfig` zapisuje komenda
// `config.session.set`; ten plik jest jego czytelnikiem na drodze tury.
//
// Obszary tłumaczą się na te same pola `models.Zapytanie`, którymi jedzie reszta
// wywołania, a stamtąd — przez `adapter_kanal_cli.go` — na wejście warstwy
// injection (`argumenty.go`, `proces.go`): model i konto na wybór wywołania,
// tools i permissions na napis `--settings`, environment i provider na zmienne
// środowiska, mcp na osobne `--mcp-config`. Zmiana obszaru w oknie konfiguracji
// zmienia więc zbudowane wywołanie modelu.
//
// Brak czytelnika albo błąd odczytu zostawia turę na wartościach okna i wiersza
// rejestru: konfiguracja sesji dokłada rozstrzygnięcia, nie odbiera dawnych.
//
// Przełożone są obszary: model, account, permissions, tools, hooks, skills
// (samo wyłączenie obszaru), environment, provider (fragment) oraz mcp. Obszary
// i fragmenty obszarów nieprzełożone wymienia blok na końcu tego pliku.
package core

import (
	"context"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// CzytelnikKonfiguracjiSesji odczytuje konfigurację obowiązującą okna. Wypełnia
// go adapterUstawienOsi (KonfiguracjaSesjiOkna); port pozwala adapterowi rozmowy
// nie znać składania obszarów. Osobny od portu KonfiguracjaSesji rodziny
// config.session.* — tamten mówi typami kontraktu, ten oddaje złożony model.
type CzytelnikKonfiguracjiSesji interface {
	KonfiguracjaSesjiOkna(ctx context.Context, idOkna, idSesji string) (shared.SessionConfig, error)
}

// ZKonfiguracjaSesji wpina czytelnika konfiguracji obowiązującej. Adapter bez
// tego portu prowadzi turę na ustawieniach okna i wiersza rejestru — tak jak
// przed wpięciem konfiguracji sesji.
func (a *adapterRozmowy) ZKonfiguracjaSesji(k CzytelnikKonfiguracjiSesji) *adapterRozmowy {
	a.konfiguracja = k
	return a
}

// uzupelnijKonfiguracje odczytuje konfigurację obowiązującą okna i przekłada jej
// obszary na pola zapytania. Osobno od uzupelnijSrodowisko, bo tamto wpina byty
// okna (katalog, mosty, wykonanie), a to — jednolity model konfiguracji sesji.
func (a *adapterRozmowy) uzupelnijKonfiguracje(ctx context.Context, okno session.Okno,
	zapytanie *models.Zapytanie) {

	if a == nil || a.konfiguracja == nil || zapytanie == nil {
		return
	}
	konfiguracja, err := a.konfiguracja.KonfiguracjaSesjiOkna(ctx, okno.Id, okno.IdSesji)
	if err != nil {
		return // brak odczytu zostawia turę na wartościach okna
	}
	przelozModel(konfiguracja, zapytanie)
	przelozKonto(konfiguracja, zapytanie)
	przelozUprawnienia(konfiguracja, zapytanie)
	przelozPowierzchnieProcesu(konfiguracja, zapytanie)
}

// przelozModel przenosi obszar model do wyboru modelu wywołania: wskazanie per
// sesja/okno trafia do przełącznika --model. Wartość pusta nie nadpisuje
// wskazania okna ani wiersza kanału.
func przelozModel(k shared.SessionConfig, z *models.Zapytanie) {
	if k.Model == nil {
		return
	}
	if prowadzacy := wartoscTekstu(k.Model.PrimaryModel); prowadzacy != "" {
		z.Model = prowadzacy
	}
	if zapasowy := wartoscTekstu(k.Model.FallbackModel); zapasowy != "" {
		z.ModelZapasowy = zapasowy
	}
	if k.Model.ReasoningEffort != nil && *k.Model.ReasoningEffort != "" {
		z.NakladRozumowania = string(*k.Model.ReasoningEffort)
	}
}

// przelozKonto przenosi obszar account do wyboru konta wywołania. Konto puste
// zostawia rozstrzygnięcie wierszowi kanału (`kanal_modelu.konto_id`) i puli
// kont — nie czyści wskazania, którego nie ma.
func przelozKonto(k shared.SessionConfig, z *models.Zapytanie) {
	if konto := kontoKonfiguracji(k.Account); konto != "" {
		z.Konto = konto
	}
}

// kontoKonfiguracji rozstrzyga konto z obszaru account wg sposobu wyboru:
// fixed i kindDefault biorą wskazane accountId, pool bierze pierwsze konto puli
// jako punkt wejścia rotacji. Rotacja po wyczerpaniu limitu należy do puli kont
// warstwy injection — tu zapada tylko wejściowe wskazanie.
func kontoKonfiguracji(a *shared.SessionConfigAccount) string {
	if a == nil {
		return ""
	}
	if a.Selection != nil && *a.Selection == shared.AccountSelectionPool {
		if len(a.PoolAccountIds) > 0 {
			return a.PoolAccountIds[0]
		}
	}
	return wartoscTekstu(a.AccountId)
}

// przelozUprawnienia przenosi tryb uprawnień obszaru permissions na tryb
// wywołania (--permission-mode). Reguły allow/deny/ask jadą osobno, plikiem
// ustawień (przelozPowierzchnieProcesu). Tryb pusty zostawia tryb okna.
func przelozUprawnienia(k shared.SessionConfig, z *models.Zapytanie) {
	if k.Permissions == nil || k.Permissions.Mode == nil || *k.Permissions.Mode == "" {
		return
	}
	z.TrybUprawnien = *k.Permissions.Mode
}

// Obszary (i fragmenty obszarów) nieprzełożone na powierzchnię procesu.
// Wymienione tu jawnie, żeby obszar wypełniony, a nieprzełożony, nie uchodził za
// wpięty:
//
//   - memory — cały obszar pozostaje poza powierzchnią. Treść pamięci
//     (InlineMemory) wymaga zapisania pliku pamięci (CLAUDE.md) w katalogu
//     roboczym, a warstwa injection nie zapisuje dziś plików sesji ani nie ma na
//     to pola w Ustawieniach/Zapytaniu. Przełączniki pamięci projektu i
//     użytkownika (ProjectMemoryEnabled, UserMemoryEnabled) oraz ścieżki
//     dodatkowe (AdditionalMemoryPaths) nie mają odpowiednika w pliku ustawień
//     bieżącej powierzchni. Wpięcie wymaga nowego mechanizmu (zapis plików sesji
//     plus pole je niosące) — leży poza tym plikiem.
//   - skills (poza wyłączeniem) — wyłączenie obszaru odmawia narzędzia Skill
//     (dodajRegulyUmiejetnosci). Dopuszczanie imienne (AllowedSkillIds), katalogi
//     wyszukiwania (Directories) i samowykrywanie (AutoDiscovery) nie mają pola
//     ani przełącznika na bieżącej powierzchni — imienny słownik reguł narzędzia
//     Skill nie jest tu potwierdzony, więc jego wpisanie byłoby atrapą.
//   - systemPrompt — tożsamość jedzie osobną drogą nakładki
//     (ZTozsamoscia + injection.Nakladka), nie tym przekładem.
//   - projectContext, conversationContext, workingDirectory, additionalDirectories,
//     inputOutput, runtime, sessionLifecycle — wpina je warstwa okna
//     (KatalogRoboczy, mosty, parametry wykonania) albo nie są wpięte wcale.
//   - provider (poza endpointUrl), model.samplingTemperature, model.betaFeatures —
//     brak odpowiednika na powierzchni CLI albo osobny przełącznik niewpięty.
