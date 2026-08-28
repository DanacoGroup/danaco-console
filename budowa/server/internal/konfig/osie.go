package konfig

import "danacoconsole/shared"

// Os nazywa oś rozstrzygania ustawienia i jest aliasem typu kontraktu. Oś
// jest prostopadła do poziomu zasięgu: poziom określa zasięg obowiązywania
// wartości, a oś określa jej adresata — platformę, model albo konto.
type Os = shared.ConfigAxis

// Wartości osi pochodzą ze stałych kontraktu, obejmujące oś platformy, oś
// modelu i oś konta; pakiet nie powtarza literału nazwy osi.
const (
	OsPlatformy Os = shared.ConfigAxisPlatform
	OsModelu    Os = shared.ConfigAxisModel
	OsKonta     Os = shared.ConfigAxisAccount
)

// osieOdNajwezszej wylicza osie w kolejności rozstrzygania w ramach jednego
// poziomu zasięgu: konto jest bytem konkretnym, model klasą, platforma tłem.
// Poziom rozstrzyga zawsze pierwszy, oś dopiero w jego ramach.
var osieOdNajwezszej = []Os{OsKonta, OsModelu, OsPlatformy}

// pierwszenstwaOsi odwzorowuje oś na kolumnę `os_zasiegu.pierwszenstwo`:
// 3 dla konta, 1 dla platformy.
var pierwszenstwaOsi = zbudujPierwszenstwaOsi()

func zbudujPierwszenstwaOsi() map[Os]int {
	wynik := make(map[Os]int, len(osieOdNajwezszej))
	for indeks, os := range osieOdNajwezszej {
		wynik[os] = len(osieOdNajwezszej) - indeks
	}
	return wynik
}

// ZnanaOs odpowiada, czy oś należy do trzech osi kontraktu; oś pusta jest
// znana i oznacza oś platformy.
func ZnanaOs(os Os) bool {
	_, jest := pierwszenstwaOsi[OsLubPlatforma(os)]
	return jest
}

// OsLubPlatforma zwraca oś wskazaną, a dla wskazania pustego — oś platformy.
// Brak osi nigdy nie jest błędem: znaczy „wartość obowiązuje tło platformy".
func OsLubPlatforma(os Os) Os {
	if os == "" {
		return OsPlatformy
	}
	return os
}
