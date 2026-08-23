package core

import "danacoconsole/shared"

// zarejestrujKanaly wpina rejestr kanałów modelu sterowany danymi:
// nowy kanał to nowy wiersz rejestru, nie nowy typ w kodzie ani nowa gałąź
// warunku. Rdzeń nie zna ani jednego rodzaju kanału — rodzaj jest wartością
// danych, którą czyta warstwa modeli.
//
// Kontrakt nie ma zdarzenia zmiany kanału, więc ta domena niczego nie rozgłasza.
// Rdzeń nie dokłada zdarzenia spoza kontraktu.
//
// Parametry kanału niosą wyłącznie odwołania do danych dostępowych, nigdy ich
// treść — pilnuje tego warstwa danych, rdzeń przenosi ładunek bez
// zaglądania do niego.
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
