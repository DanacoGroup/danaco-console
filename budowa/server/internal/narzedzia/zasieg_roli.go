// Zasięg narzędzi wynikający z roli okna: okno robocze dostaje sam wykaz
// kontraktu, okno asystenta dostaje wykaz powiększony o komendy nastawiania
// okna docelowego.
package narzedzia

import (
	"reflect"
	"strings"

	"danacoconsole/shared"
)

// Zasieg nazywa rolę danego okna tej rozmowy, w którego imieniu pracuje serwer
// narzędzi, jedną z trzech.
type Zasieg string

const (
	// ZasiegOkna jest zasięgiem okna roboczego: sam wykaz kontraktu, ani jednej
	// pozycji więcej, zasięg domyślny.
	ZasiegOkna Zasieg = "okno"
	// ZasiegKlawiatury jest zasięgiem okna asystenta, czyli klawiatury
	// Operatora, powiększonym o komendy roli.
	ZasiegKlawiatury Zasieg = "klawiatura"
	// PrzelacznikZasiegu jest nazwą przełącznika uruchomieniowego, który niesie
	// zasięg do tego serwera narzędzi.
	PrzelacznikZasiegu = "zasieg"
)

// RozpoznajZasieg przekłada napis przełącznika na zasięg; napis nieznany daje
// zasięg okna roboczego, nie błąd.
func RozpoznajZasieg(napis string) Zasieg {
	switch Zasieg(strings.TrimSpace(napis)) {
	case ZasiegKlawiatury:
		return ZasiegKlawiatury
	case ZasiegEksperta:
		return ZasiegEksperta
	default:
		return ZasiegOkna
	}
}

// ArgumentyZasiegu zwraca argumenty uruchomienia dopisywane do wpisu `danaco`
// dla wskazanego zasięgu. Zasięg okna roboczego nie dopisuje nic — wpis
// dotychczasowy zostaje wpisem dotychczasowym, co do znaku.
func ArgumentyZasiegu(zasieg Zasieg) []string {
	if zasieg != ZasiegKlawiatury {
		return nil
	}
	return []string{"--" + PrzelacznikZasiegu, string(ZasiegKlawiatury)}
}

// pozycjaRoli jest jedną komendą dokładaną przez rolę okna; Zastosowanie jest
// zdaniem polityki, nie kontraktu.
type pozycjaRoli struct {
	Komenda      shared.MessageType
	Zastosowanie string
	Zadanie      any
}

// rozszerzenieKlawiatury wymienia komendy, które okno asystenta ma ponad wykaz
// kontraktu. Obie wykonuje się na oknie docelowym, nie na oknie asystenta.
var rozszerzenieKlawiatury = []pozycjaRoli{
	{
		Komenda: shared.CommandModelChannelSet,
		Zastosowanie: "Uzyj, aby wybrac kanal modelu dla OKNA DOCELOWEGO, w ktorym ma pracowac model " +
			"wykonujacy zlecenie. Wskaz windowId okna docelowego — nigdy okna asystenta",
		Zadanie: shared.ModelChannelSetRequest{},
	},
	{
		Komenda: shared.CommandConfigSessionSet,
		Zastosowanie: "Uzyj, aby ustawic konfiguracje zlecenia na OKNIE DOCELOWYM przed wyslaniem promptu " +
			"— scope window i scopeId okna docelowego. Obszar jest jednostka zapisu (pole areas)",
		Zadanie: shared.ConfigSessionSetRequest{},
	},
}

// narzedziaRoli składa pozycje wykazu dokładane przez rolę okna; komenda musi
// istnieć w kontrakcie i stać poza wykazem narzędzi.
func narzedziaRoli(zasieg Zasieg) []Narzedzie {
	if zasieg != ZasiegKlawiatury {
		return nil
	}
	poza := komendyPozaWykazem()
	nazwy := map[shared.MessageType]string{}
	for nazwa, komenda := range poza {
		nazwy[komenda] = nazwa
	}
	wykaz := make([]Narzedzie, 0, len(rozszerzenieKlawiatury))
	for _, pozycja := range rozszerzenieKlawiatury {
		nazwa, stoiPoza := nazwy[pozycja.Komenda]
		if !stoiPoza {
			continue
		}
		wykaz = append(wykaz, Narzedzie{
			Nazwa:   nazwa,
			Opis:    pozycja.Zastosowanie,
			Schemat: schematWejscia(polaZadania(reflect.TypeOf(pozycja.Zadanie))),
			Grupa:   grupaKomendy(pozycja.Komenda),
		})
	}
	return wykaz
}

// komendaRoli odszukuje komendę narzędzia dołożonego przez rolę okna. Zwraca
// fałsz dla każdej nazwy spoza rozszerzenia TEGO zasięgu — łącznie z nazwami
// rozszerzeń innych ról, gdyby kiedyś powstały.
func komendaRoli(nazwa string, zasieg Zasieg) (shared.MessageType, bool) {
	for _, pozycja := range narzedziaRoli(zasieg) {
		if pozycja.Nazwa == nazwa {
			return komendaPoNazwie(nazwa), true
		}
	}
	return "", false
}

// komendaPoNazwie oddaje komendę stojącą pod nazwą danego narzędzia spoza
// wykazu narzędzi tego kontraktu.
func komendaPoNazwie(nazwa string) shared.MessageType {
	return komendyPozaWykazem()[nazwa]
}
