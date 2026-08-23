package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
)

// obsluz buduje obsługiwacza komendy z czynności domeny.
//
// Cała powtarzalna praca dyspozycji dzieje się tutaj raz: odczytanie ładunku
// w kształcie żądania kontraktu, wywołanie czynności, zamiana wyniku albo błędu
// na odpowiedź protokołu. Dzięki temu pliki handlers_*.go zawierają wyłącznie
// wiązanie nazwy kontraktu z czynnością — bez powielonej obsługi błędów.
//
// Typ żądania Z i typ wyniku W pochodzą z pakietu shared, więc zmiana kontraktu
// przerywa kompilację obsługiwacza zamiast rozjeżdżać się z nim po cichu.
func obsluz[Z any, W any](czynnosc func(context.Context, Z) (W, error)) Obsluga {
	return func(ctx context.Context, z protocol.Request) protocol.Odpowiedz {
		var zadanie Z
		if err := z.LadunekDo(&zadanie); err != nil {
			return bladNiepoprawnegoLadunku(err)
		}
		// Tożsamość żądania jedzie dalej kontekstem, bo ładunek jej nie niesie:
		// pole `id` mieszka w kopercie, a czynność domeny dostaje wyłącznie
		// rozpakowaną treść. Bez tego wpisu nadawca strumienia nie miałby czym
		// powtórzyć identyfikatora zadania w kopertach `stream.chunk`, choć
		// kontrakt każe mu go powtarzać (`protocol/tozsamosc_zadania.go`).
		rezultat, err := czynnosc(protocol.ZIdZadania(ctx, z.Id), zadanie)
		if err != nil {
			return porazka(err)
		}
		return wynik(rezultat)
	}
}
