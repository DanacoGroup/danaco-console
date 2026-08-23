package core

import "danacoconsole/server/internal/konfig"

// DefinicjeKataloguRoboczego zwraca definicje dwóch ustawień katalogu roboczego
// gotowe do dołożenia do rejestru definicji rezolwera.
//
// Definicje powstają tutaj, nie w pakiecie konfig, z powodu kierunku zależności:
// konfig nie zna pakietu core, a klucze i wartości domyślne należą do modułu
// katalogu roboczego. Dołożenie ich do rejestru jest jednym wywołaniem
// konfig.Rejestr.Dodaj przy montażu — rejestr jest zbiorem otwartym, więc nowa
// pozycja okna konfiguracji nie zmienia rozstrzygania.
//
// Wartość domyślna podstawy jest pusta z zamysłem: miejsce instalacji aplikacji
// głównej ustala się w chwili startu procesu i nie da się go zapisać stałą
// tekstową. Pustkę czyta PodstawaLubInstalacja i zamienia na miejsce
// instalacji.
func DefinicjeKataloguRoboczego() []konfig.Definicja {
	return []konfig.Definicja{
		{
			Klucz: KluczKatalogRoboczyPodstawa, Domyslna: "", Rodzaj: konfig.RodzajTekst,
			Objasnienie: "Katalog, w którym powstają katalogi sesyjne i pliki robocze modelu. " +
				"Brak wskazania znaczy miejsce instalacji aplikacji głównej. " +
				"To nie jest dostęp — dostęp mówi, do czego model ma wgląd.",
		},
		{
			Klucz: KluczKatalogRoboczyWzorzecSesji, Domyslna: WzorzecSesjiDomyslny, Rodzaj: konfig.RodzajTekst,
			Objasnienie: "Wzorzec nazwy katalogu jednej sesji, liczony względem podstawy. " +
				"Znacznik " + ZnacznikIdentyfikatoraSesji + " jest miejscem podstawienia " +
				"identyfikatora sesji. Wzorzec wyprowadzający poza podstawę jest pomijany.",
		},
	}
}
