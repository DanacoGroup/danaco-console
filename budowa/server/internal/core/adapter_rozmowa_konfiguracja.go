// Plik doprowadza obowiązującą konfigurację sesji do zbudowanego wywołania
// modelu, przekładając obszary konfiguracji na pola zapytania modelu, konta,
// narzędzi, uprawnień, środowiska i dostawcy.
package core

import (
	"context"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// CzytelnikKonfiguracjiSesji odczytuje konfigurację obowiązującą okna, złożoną
// z jej obszarów, tak aby adapter rozmowy nie musiał znać ich składania.
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

// kontoKonfiguracji rozstrzyga konto z obszaru account: sposoby fixed
// i kindDefault biorą wskazane konto, sposób pool bierze pierwsze konto puli
// jako wejście rotacji.
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

// Część obszarów i ich fragmentów nie ma dziś odpowiednika na powierzchni procesu.
