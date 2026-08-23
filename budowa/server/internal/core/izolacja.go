// Odpowiedzialność pliku: wejście do egzekutora izolacji od strony rdzenia —
// rozstrzygnięcie jedenastu punktów dla kontekstu zasięgu oraz złożenie obszaru
// własnego okna komunikacji.
//
// Rozstrzyganie należy do pakietu konfig i tutaj się go nie powtarza;
// ten plik wyłącznie łączy wynik rozstrzygnięcia z postacią wykonawczą, której
// używają strażnicy plików, sieci, kontekstu, polecenia i przydziału.
package core

import (
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
)

// nazwaKataloguDanychModelu jest nazwą podkatalogu, w którym leżą dane
// pomocnicze kanału modelu wydzielone dla jednego okna: konfiguracja kanału,
// dane tymczasowe, ustawienia dostawcy. Podkatalog leży wewnątrz własnego
// katalogu okna, więc izolacja katalogu danych nie wymaga drugiej podstawy.
const nazwaKataloguDanychModelu = "dane-modelu"

// ZasadyIzolacji rozstrzyga jedenaście punktów izolacji obowiązujących
// w kontekście zasięgu i sprowadza je do postaci wykonawczej. Brak rozstrzygacza
// daje stan wyjściowy platformy — kontekst odrębny, żaden zakres techniczny
// niewłączony — a nie odmowę wykonania.
func ZasadyIzolacji(rozstrzygacz *konfig.Rozstrzygacz, kontekst konfig.Kontekst) session.Zasady {
	return session.ZasadyZPolityki(rozstrzygacz.PolitykaEfektywna(kontekst))
}

// ObszarOkna składa obszar własny okna z ustalonego katalogu roboczego.
// Katalog ustala KatalogRoboczy po ośmiu poziomach zasięgu i tylko on — obszar
// niczego nie wylicza od nowa, dokłada wyłącznie podkatalog danych modelu.
//
// Ustalenie bez ścieżki daje obszar pusty. Obszar pusty przy izolacji włączonej
// jest naruszeniem rozstrzyganym w egzekutorze polecenia, a nie milczącym
// przejściem: okno bez własnego katalogu poszłoby do katalogu wspólnego.
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
