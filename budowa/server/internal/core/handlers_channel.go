package core

import "danacoconsole/shared"

// zarejestrujKanaly wpina rejestr kanałów modelu sterowany danymi, w którym
// nowy kanał jest nowym wierszem rejestru, a nie nowym typem w kodzie.
func zarejestrujKanaly(r *Rejestr, kanaly Kanaly) {
	if r == nil || kanaly == nil {
		return
	}
	r.Zarejestruj(shared.CommandChannelAdd, obsluz(kanaly.Dodaj))
	r.Zarejestruj(shared.CommandChannelUpdate, obsluz(kanaly.Zmien))
	r.Zarejestruj(shared.CommandChannelRemove, obsluz(kanaly.Usun))
	r.Zarejestruj(shared.CommandChannelList, obsluz(kanaly.Wykaz))
	r.Zarejestruj(shared.CommandChannelCheck, obsluz(kanaly.Sprawdz))
	r.Zarejestruj(shared.CommandChannelCredentialStatus, obsluz(kanaly.StanPoswiadczenia))
}
