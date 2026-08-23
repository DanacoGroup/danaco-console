package konfig

import "danacoconsole/shared"

// Os nazywa oś rozstrzygania ustawienia. Jest aliasem typu kontraktu, więc
// pakiet nie tworzy drugiej definicji tego samego pojęcia.
//
// Oś jest PROSTOPADŁA do poziomu zasięgu: poziom mówi JAK WĄSKO obowiązuje
// wartość (okno → … → globalny), oś mówi DLA CZEGO ona obowiązuje —
// dla platformy, dla modelu albo dla konta. Klucz rozstrzygania jest więc
// złożony: klucz + poziom + byt poziomu + oś + byt osi.
type Os = shared.ConfigAxis

// Wartości osi. Pochodzą ze stałych kontraktu — literału nazwy osi tu nie ma.
const (
	OsPlatformy Os = shared.ConfigAxisPlatform
	OsModelu    Os = shared.ConfigAxisModel
	OsKonta     Os = shared.ConfigAxisAccount
)

// osieOdNajwezszej wylicza osie w kolejności rozstrzygania W RAMACH jednego
// poziomu zasięgu: konto jest bytem konkretnym, model klasą, platforma tłem.
//
// Rozstrzygnięcie pierwszeństwa. Poziom rozstrzyga PIERWSZY, oś dopiero
// w ramach poziomu. Ustawienie zapisane per konto na poziomie globalnym NIE
// bije ustawienia zapisanego na oknie komunikacji: zapis na oknie jest aktem
// najwęższym i najbardziej celowym, a oś opisuje adresata wartości, nie jej
// wagę. Odwrotna kolejność znaczyłaby, że wybór konta unieważnia decyzję
// podjętą wprost w oknie — a to jest odebranie Operatorowi sterowania, nie
// jego rozszerzenie.
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

// ZnanaOs odpowiada, czy oś należy do trzech osi kontraktu. Oś pusta jest znana:
// znaczy platformę.
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
