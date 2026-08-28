// Plik buduje rejestr definicji z katalogu ustawień zapisanego w bazie,
// zamiast z listy zaszytej w kodzie. Katalog pusty, niedostępny albo
// uszkodzony nie zatrzymuje rozstrzygania; rejestr schodzi wtedy na
// wbudowany zestaw definicji rdzenia.
package konfig

import "danacoconsole/shared"

// zrodloKatalogu dostarcza pozycje katalogu ustawień zapisanego w bazie,
// oddzielając rejestr definicji od bezpośredniego dostępu do połączenia.
type zrodloKatalogu interface {
	Definicje() ([]shared.SettingDefinition, error)
}

// rodzajeKontraktu odwzorowują rodzaj pozycji katalogu na postać wartości
// zapisywanej w kolumnie ustawienie.rodzaj_wartosci; rodzaje prezentacyjne
// są napisami, listy i struktury idą jako JSON, liczby i wartości logiczne
// mają własne rodzaje.
var rodzajeKontraktu = map[shared.SettingValueType]Rodzaj{
	shared.SettingValueTypeString:   RodzajTekst,
	shared.SettingValueTypeText:     RodzajTekst,
	shared.SettingValueTypeEnum:     RodzajTekst,
	shared.SettingValueTypePath:     RodzajTekst,
	shared.SettingValueTypeSecret:   RodzajTekst,
	shared.SettingValueTypeInt:      RodzajLiczba,
	shared.SettingValueTypeFloat:    RodzajLiczba,
	shared.SettingValueTypeBool:     RodzajLogiczna,
	shared.SettingValueTypeJson:     RodzajJSON,
	shared.SettingValueTypeEnumList: RodzajJSON,
	shared.SettingValueTypePathList: RodzajJSON,
}

// rodzajZKontraktu zwraca postać wartości dla rodzaju pozycji katalogu.
// Rodzaj nierozpoznany daje tekst — nie przerywa odczytu katalogu.
func rodzajZKontraktu(rodzaj shared.SettingValueType) Rodzaj {
	if postac, jest := rodzajeKontraktu[rodzaj]; jest {
		return postac
	}
	return RodzajTekst
}

// definicjaZKontraktu zamienia pozycję katalogu na definicję rejestru.
// Wartość domyślna wraca do postaci przechowywanej w kolumnie
// ustawienie.wartosc, więc porównanie z wartością zapisaną jest porównaniem
// dwóch napisów tej samej postaci.
func definicjaZKontraktu(pozycja shared.SettingDefinition) Definicja {
	domyslna, _ := DekodujJSON(pozycja.DefaultValue)
	return Definicja{
		Klucz:            pozycja.Key,
		Domyslna:         domyslna,
		Rodzaj:           rodzajZKontraktu(pozycja.ValueType),
		Objasnienie:      tekstLubPusty(pozycja.Description),
		Kategoria:        pozycja.CategoryId,
		Nazwa:            pozycja.Name,
		DozwoloneZasiegi: pozycja.AllowedScopes,
		DozwoloneOsie:    pozycja.AllowedAxes,
		Kolejnosc:        pozycja.Order,
		Aktywna:          pozycja.Enabled,
	}
}

// RejestrZKatalogu buduje rejestr definicji z katalogu ustawień. Katalog
// pusty albo niedostępny daje rejestr wbudowany rdzenia. Drugi wynik mówi,
// czy rejestr pochodzi z katalogu, i służy diagnostyce startu.
func RejestrZKatalogu(zrodlo zrodloKatalogu) (*Rejestr, bool) {
	if zrodlo == nil {
		return RejestrWbudowany(), false
	}
	pozycje, err := zrodlo.Definicje()
	if err != nil || len(pozycje) == 0 {
		return RejestrWbudowany(), false
	}
	definicje := make([]Definicja, 0, len(pozycje))
	for _, pozycja := range pozycje {
		if pozycja.Key == "" {
			continue
		}
		definicje = append(definicje, definicjaZKontraktu(pozycja))
	}
	if len(definicje) == 0 {
		return RejestrWbudowany(), false
	}
	return NowyRejestr(definicje...), true
}

// tekstLubPusty odczytuje pole opcjonalne kontraktu i zwraca napis pusty
// zamiast wskaźnika pustego, upraszczając dalsze porównania wartości.
func tekstLubPusty(wartosc *string) string {
	if wartosc == nil {
		return ""
	}
	return *wartosc
}
