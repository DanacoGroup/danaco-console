package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
)

// obsluz buduje obsługiwacza komendy z czynności domeny: odczytuje ładunek żądania kontraktu, wywołuje czynność i zamienia wynik lub błąd na odpowiedź protokołu.
func obsluz[Z any, W any](czynnosc func(context.Context, Z) (W, error)) Obsluga {
	return func(ctx context.Context, z protocol.Request) protocol.Odpowiedz {
		var zadanie Z
		if err := z.LadunekDo(&zadanie); err != nil {
			return bladNiepoprawnegoLadunku(err)
		}
		// Zgodność z kontraktem sprawdza się po odczytaniu ładunku, przed wywołaniem czynności domeny.
		if err := sprawdzZadanieWobecKontraktu(ctx, z.Komenda, z.Ladunek, zadanie); err != nil {
			return porazka(err)
		}
		// Tożsamość żądania jedzie kontekstem, bo ładunek jej nie niesie; pole id mieszka w kopercie.
		rezultat, err := czynnosc(protocol.ZIdZadania(ctx, z.Id), zadanie)
		if err != nil {
			return porazka(err)
		}
		return wynik(rezultat)
	}
}
