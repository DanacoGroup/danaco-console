package core

import "danacoconsole/shared"

// zarejestrujAkcje wpina katalog akcji sterowany danymi: panel akcji i narzędzia kanału
// modelu czerpią z jednego rejestru budowanego z wierszy tabeli, nie z listy zaszytej
// w kodzie.
func zarejestrujAkcje(r *Rejestr, akcje Akcje, nazwa shared.MessageType) {
	if r == nil || akcje == nil || nazwa == "" {
		return
	}
	r.Zarejestruj(nazwa, obsluz(akcje.Wykaz))
}
