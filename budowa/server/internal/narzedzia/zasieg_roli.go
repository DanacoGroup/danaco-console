// Zasięg narzędzi wynikający z roli okna.
//
// Serwer narzędzi należy do jednego okna, a wpis `danaco` w konfiguracji MCP
// powstaje osobno dla każdego okna (`wpiecie.go`), więc zasięg narzędzi wiąże
// się z oknem.
//
// Wykaz narzędzi kontraktu nie niesie `config.session.set` ani
// `model.channel.set`: model nie przestawia sobie własnego wyposażenia. Dla okna
// roboczego jest to właściwe. Okno asystenta ustawia jednak konfigurację
// zlecenia i wybiera kanał modelu — nie na sobie, lecz na oknie docelowym,
// w którym ma pracować model wykonujący zlecenie.
//
// Zasięg jest funkcją roli okna, nie prośby modelu. Ten plik definiuje dwie
// role, trzecią dokłada `zasieg_eksperta.go`:
//
//	ZasiegOkna       — okno robocze: dokładnie wykaz kontraktu, ani jednej
//	                   pozycji więcej. Zasięg domyślny.
//	ZasiegKlawiatury — okno asystenta: wykaz kontraktu powiększony o komendy
//	                   nastawiania cudzego okna wymienione niżej.
//
// Rozszerzenie wynika z roli okna zapisanej w rdzeniu i wchodzi do wpisu MCP
// w chwili jego składania, czyli zanim proces modelu wystartuje. Model nie ma
// czym o rozszerzenie poprosić: nie ma narzędzia zmieniającego zasięg,
// a przełącznik `--zasieg` czyta się raz, przy uruchomieniu serwera, z wpisu
// ułożonego przez rdzeń.
//
// Rozszerzenie nie tyka komend zastrzeżonych Operatorowi (punkty dostępu,
// nadania, konta, tożsamość) ani warstwy połączenia klienta — te zostają poza
// wykazem w każdym zasięgu, co sprawdza strukturalnie `zastrzezenia.go`.
//
// Wykaz rozszerzenia jest polityką, nie kopią danych kontraktu. Kontrakt nie zna
// pojęcia okna asystenta: sekcja `narzedzia` jest listą płaską, bez wymiaru
// roli, więc nie ma w niej danych, z których dałoby się ten podzbiór wyprowadzić.
// Dlatego wykaz jest wymieniony z nazwy, opatrzony powodem i sprawdzany przy
// budowie zasięgu: nazwa, która przestanie być komendą kontraktu albo wejdzie do
// wykazu narzędzi, znika z rozszerzenia zamiast po cichu wisieć.
//
// Schemat wejścia powstaje z wygenerowanej struktury żądania (`schemat_roli.go`),
// bo dla komendy spoza sekcji `narzedzia` kontrakt nie wylicza `ToolParameter`.
// Pola mają zatem typy, ale nie mają opisów ani wykazu wartości wyliczeń.
package narzedzia

import (
	"reflect"
	"strings"

	"danacoconsole/shared"
)

// Zasieg nazywa rolę okna, w którego imieniu pracuje serwer narzędzi.
type Zasieg string

const (
	// ZasiegOkna jest zasięgiem okna roboczego: sam wykaz kontraktu.
	ZasiegOkna Zasieg = "okno"
	// ZasiegKlawiatury jest zasięgiem okna asystenta — klawiatury Operatora.
	ZasiegKlawiatury Zasieg = "klawiatura"
	// PrzelacznikZasiegu jest nazwą przełącznika niosącego zasięg do serwera.
	PrzelacznikZasiegu = "zasieg"
)

// RozpoznajZasieg przekłada napis przełącznika na zasięg. Napis nieznany daje
// zasięg okna roboczego, nie błąd: zasięg węższy jest bezpiecznym domyślnym,
// a serwer ma pracować dalej. Rozszerzenia nie dostaje się przez pomyłkę
// w napisie.
//
// Rozpoznawane są trzy zasięgi; trzeci — ekspercki — definiuje
// `zasieg_eksperta.go`. Dla `core/sprawca.go`, który pyta wyłącznie
// o klawiaturę, zasięg eksperta spada na model roboczy.
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

// pozycjaRoli jest jedną komendą dokładaną przez rolę okna.
//
// `Zastosowanie` jest zdaniem polityki, nie kontraktu: mówi modelowi asystenta,
// po co sięga po komendę, której zwykłe okno nie ma. Kontrakt takiego zdania nie
// niesie, bo nie zna roli, dla której ono powstaje.
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

// narzedziaRoli składa pozycje wykazu dokładane przez rolę okna.
//
// Komenda musi istnieć w kontrakcie i stać poza wykazem narzędzi. Pierwszy
// warunek chroni przed nazwą, która z kontraktu wypadła; drugi — przed podwójną
// pozycją tego samego narzędzia, gdyby kontrakt kiedyś wciągnął tę komendę do
// wykazu sam. Oba liczy `komendyPozaWykazem()`, czyli ta sama różnica, którą
// posługuje się odmowa — jedno źródło rozstrzygnięcia, co stoi poza wykazem.
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

// komendaPoNazwie oddaje komendę stojącą pod nazwą narzędzia spoza wykazu.
func komendaPoNazwie(nazwa string) shared.MessageType {
	return komendyPozaWykazem()[nazwa]
}
