// Plik wpina wejście do egzekutora izolacji od strony rdzenia: rozstrzygnięcie jedenastu punktów
// dla kontekstu zasięgu oraz złożenie obszaru własnego okna komunikacji.
package core

import (
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
)

// nazwaKataloguDanychModelu jest nazwą podkatalogu, w którym leżą dane pomocnicze kanału modelu
// wydzielone dla jednego okna: konfiguracja kanału, dane tymczasowe, ustawienia dostawcy.
const nazwaKataloguDanychModelu = "dane-modelu"

// ZasadyIzolacji rozstrzyga jedenaście punktów izolacji obowiązujących w kontekście zasięgu i
// sprowadza je do postaci wykonawczej, bez rozstrzygacza dając stan wyjściowy platformy.
func ZasadyIzolacji(rozstrzygacz *konfig.Rozstrzygacz, kontekst konfig.Kontekst) session.Zasady {
	return session.ZasadyZPolityki(rozstrzygacz.PolitykaEfektywna(kontekst))
}

// ObszarOkna składa obszar własny okna z ustalonego katalogu roboczego, dokładając wyłącznie
// podkatalog danych modelu do katalogu ustalonego po ośmiu poziomach zasięgu.
func ObszarOkna(ustalenie UstalenieKatalogu, idOkna string) session.Obszar {
	sciezka := strings.TrimSpace(ustalenie.Sciezka)
	obszar := session.Obszar{IdOkna: idOkna}
	if sciezka == "" {
		return obszar
	}
	obszar.KatalogRoboczy = sciezka
	obszar.KatalogDanych = filepath.Join(sciezka, nazwaKataloguDanychModelu)
	return obszar
}
