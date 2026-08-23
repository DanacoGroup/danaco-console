// Rozgłoszenia rodziny `isolation.*`: dwa zdarzenia kontraktu,
// `isolation.profile.changed` i `isolation.policy.changed`, wraz z regułą, które
// wywołanie które z nich wysyła.
//
// Zdarzenia opisują dwa różne byty. `isolation.profile.changed` dotyczy
// szablonu: profil powstał, zmienił się albo zniknął z katalogu; zapisanie
// profilu nie zmienia izolacji (`adapter_modul_isolation_profile.go`), więc samo
// to zdarzenie nie znaczy zmiany warunków pracy. `isolation.policy.changed`
// dotyczy polityki obowiązującej konkretne okno i jest sygnałem do przerysowania.
//
// Zdarzenie nie zastępuje `config.changed`: zapis punktu izolacji nadal rozgłasza
// `config.changed` z wierszem tabeli `ustawienie`, bo tak czyta go panel
// konfiguracji. `isolation.policy.changed` nie niesie wartości, tylko wskazanie
// okna, które ma przeliczyć swoją politykę.
//
// Ładunek kontraktu ma `windowId` jako pole wymagane, więc zmiana na poziomie
// szerszym niż okno — globalnym, środowiska, karty sesji — nie ma czego w nim
// postawić. Rdzeń nie rozgłasza wtedy zdarzenia z pustym oknem: takie zmiany
// zostają przy `config.changed`, bo rdzeń nie umie rozwinąć poziomu szerszego na
// listę objętych nim okien.
package core

import (
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// rozglosProfilIzolacji rozgłasza zmianę profilu izolacji.
func (a *adapterIzolacji) rozglosProfilIzolacji(zmiana shared.ChangeKind, kodProfilu string) {
	if a.rozgloszenie == nil || kodProfilu == "" {
		return
	}
	a.rozgloszenie.wyslij(shared.EventIsolationProfileChanged, "",
		shared.IsolationProfileChangedEvent{Change: zmiana, ProfileId: kodProfilu})
}

// rozglosPolitykeAdresu rozgłasza zmianę polityki obowiązującej okno, gdy adres
// zapisu okna dotyczy. Adres spod innego poziomu nie rozgłasza niczego —
// powody w nagłówku pliku.
func (a *adapterIzolacji) rozglosPolitykeAdresu(adres konfig.Adres) {
	if adres.Poziom != shared.ConfigScopeWindow {
		return
	}
	a.rozglosPolitykeOkna(adres.KluczZasiegu)
}

// rozglosPolitykeOkna rozgłasza `isolation.policy.changed` dla jednego okna.
// Rodzaj zmiany jest zawsze `updated`: polityka okna nie powstaje ani nie znika,
// okno ma ją zawsze, choćby złożoną z samych wartości domyślnych.
func (a *adapterIzolacji) rozglosPolitykeOkna(idOkna string) {
	if a.rozgloszenie == nil || idOkna == "" {
		return
	}
	a.rozgloszenie.wyslij(shared.EventIsolationPolicyChanged, "",
		shared.IsolationPolicyChangedEvent{Change: shared.ChangeKindUpdated, WindowId: idOkna})
}

// rozglosPolitykeWyboruWarstwy rozgłasza zmianę polityki po `isolation.layer.set`.
//
// Przełączenie warstwy zmienia politykę obowiązującą, choć nie zmienia ani
// jednej wartości: warstwa rozstrzyga, spod którego adresu wartości są czytane
// (`warstwaAdresu`), więc okno po przełączeniu pracuje na innym zestawie. Wybór
// zapisany dla karty sesji nie ma okna do wskazania i nie rozgłasza się — tak
// samo jak każdy poziom szerszy niż okno.
func (a *adapterIzolacji) rozglosPolitykeWyboruWarstwy(poziom shared.ConfigScope, byt string) {
	if poziom != shared.ConfigScopeWindow {
		return
	}
	a.rozglosPolitykeOkna(byt)
}

// rodzajZapisuProfilu rozstrzyga, czy `isolation.profile.save` założył profil,
// czy zmienił zastany. Rozstrzyga to samo żądanie: puste `profileId` znaczy
// profil nowy — tak mówi kontrakt i tak zachowuje się adapter.
func rodzajZapisuProfilu(idProfilu *string) shared.ChangeKind {
	if idProfilu == nil || *idProfilu == "" {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindUpdated
}
