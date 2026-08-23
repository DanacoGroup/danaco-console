// Odpowiedzialność pliku: tryb uprawnień okna jako brama uruchomienia procesu
// w module Terminal.
//
// Tryb uprawnień rozstrzyga tutaj, a nie tylko w kanale modelu. Wartości
// PermissionMode odpowiadają przełącznikowi `--permission-mode` kanału
// głównego i tam ograniczają model. Terminal uruchamia proces urządzenia
// z pominięciem kanału, więc bez tej samej bramy okno w trybie `plan` („praca
// planistyczna bez zmian w systemie") uruchamiałoby polecenia powłoki.
//
// Rozstrzyga inicjator, nie samo okno. Kontrakt nie ma rundy zgody dla
// terminala: nie istnieje komenda pytająca Operatora „czy uruchomić"
// i czekająca na odpowiedź. Dlatego w trybach wymagających zgody (`manual`,
// `acceptEdits`) proces zlecony przez model zostaje odrzucony — zgody nie ma
// jak uzyskać. Polecenie Operatora przechodzi, bo jego kliknięcie jest tą zgodą.
package core

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// zgodaNaProces mówi, co wolno w danym trybie uprawnień okna.
type zgodaNaProces struct {
	// Operator — czy wolno uruchomić proces zlecony przez Operatora.
	operator bool
	// Model — czy wolno uruchomić proces zlecony przez model.
	model bool
	// Powod tłumaczy odmowę treścią zrozumiałą dla Operatora.
	powod string
}

// zgodyTrybow wiąże sześć trybów uprawnień kontraktu z zakresem uruchomienia.
// Wykaz jest danymi — nowy tryb w kontrakcie to nowa pozycja tutaj,
// a nie nowa gałąź warunku.
var zgodyTrybow = map[shared.PermissionMode]zgodaNaProces{
	shared.PermissionModePlan: {
		operator: false, model: false,
		powod: "okno pracuje w trybie planistycznym (plan), który wyklucza zmiany w systemie",
	},
	shared.PermissionModeManual: {
		operator: true, model: false,
		powod: "tryb ręczny (manual) wymaga zgody przed każdą zmianą, " +
			"a kontrakt nie ma rundy zgody dla terminala",
	},
	shared.PermissionModeAcceptEdits: {
		operator: true, model: false,
		powod: "tryb acceptEdits daje zgodę na zmiany plików, nie na uruchamianie procesów urządzenia",
	},
	shared.PermissionModeAuto:              {operator: true, model: true},
	shared.PermissionModeDontAsk:           {operator: true, model: true},
	shared.PermissionModeBypassPermissions: {operator: true, model: true},
}

// zgodaTrybu zwraca zakres uruchomienia dla trybu okna. Tryb nierozpoznany —
// na przykład z okna zapisanego przez starszą wersję — schodzi na zakres
// najostrożniejszy: Operator tak, model nie.
func zgodaTrybu(tryb shared.PermissionMode) zgodaNaProces {
	if zgoda, jest := zgodyTrybow[tryb]; jest {
		return zgoda
	}
	return zgodaNaProces{
		operator: true, model: false,
		powod: "tryb uprawnień okna (" + string(tryb) + ") nie należy do słownika kontraktu",
	}
}

// sprawdzUprawnienie rozstrzyga, czy okno o danym trybie może uruchomić proces
// zlecony przez wskazanego inicjatora.
func sprawdzUprawnienie(tryb shared.PermissionMode, inicjator shared.ProcessInitiator) error {
	zgoda := zgodaTrybu(tryb)
	wolno := zgoda.operator
	kto := "Operatora"
	if inicjator == shared.ProcessInitiatorModel {
		wolno, kto = zgoda.model, "model"
	}
	if wolno {
		return nil
	}
	powod := zgoda.powod
	if powod == "" {
		powod = "tryb uprawnień okna wyklucza tę czynność"
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
		"moduł Terminal: proces zlecony przez "+kto+" nie ruszy — "+powod))
}
