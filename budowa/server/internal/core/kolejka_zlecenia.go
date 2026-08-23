package core

import (
	"bytes"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Odczyt zleceń początkowych kolejki z ładunku komendy queue.create.
//
// Kontrakt niesie pole `items` jako surowy JSON
// (QueueCreateRequest.Items = json.RawMessage) i nie opisuje pozycji kolejki
// strukturą. Dopóki tak jest, rdzeń przyjmuje trzy zapisy tej samej rzeczy —
// odmowa z powodu nawiasu byłaby odmową z powodu braku kontraktu, nie z powodu
// błędu Operatora:
//
//	[{"title":"…","content":"…"}, …]  — zlecenie z tytułem i treścią
//	["…", …]                          — sama treść zlecenia
//	{"items":[…]}                     — wykaz w opakowaniu
//
// Brak pola, pole puste i `null` dają kolejkę bez zleceń — kolejka pusta jest
// stanem poprawnym. Ładunek nieczytelny jako JSON jest błędem
// żądania: to nie brak danych, tylko dane uszkodzone.

// dlugoscTytuluZTresci ogranicza tytuł wyprowadzony z treści zlecenia. Kolumna
// `pozycja_kolejki.tytul` jest etykietą pozycji w Mission Control, nie kopią
// polecenia.
const dlugoscTytuluZTresci = 80

// zlecenieKolejki to jedno zlecenie w postaci przyjmowanej z ładunku komendy.
type zlecenieKolejki struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// odczytajZlecenia przekłada ładunek `items` na pozycje kolejki gotowe do
// zapisania. Kolejność wykazu jest kolejnością wykonania.
func odczytajZlecenia(ladunek json.RawMessage) ([]dane.Pozycja, error) {
	elementy, err := elementyZlecen(ladunek)
	if err != nil {
		return nil, err
	}
	pozycje := make([]dane.Pozycja, 0, len(elementy))
	for _, element := range elementy {
		zlecenie, err := odczytajZlecenie(element)
		if err != nil {
			return nil, err
		}
		if pozycja, jest := pozycjaZeZlecenia(zlecenie); jest {
			pozycje = append(pozycje, pozycja)
		}
	}
	return pozycje, nil
}

// elementyZlecen rozpakowuje ładunek do wykazu elementów. Obiekt z polem
// `items` jest opakowaniem wykazu; obiekt bez tego pola jest pojedynczym
// zleceniem zapisanym bez nawiasu tablicy.
func elementyZlecen(ladunek json.RawMessage) ([]json.RawMessage, error) {
	tresc := bytes.TrimSpace(ladunek)
	if len(tresc) == 0 || string(tresc) == "null" {
		return nil, nil
	}
	if tresc[0] == '{' {
		var opakowanie struct {
			Items json.RawMessage `json:"items"`
		}
		if err := json.Unmarshal(tresc, &opakowanie); err != nil {
			return nil, bladZlecen(err)
		}
		if len(bytes.TrimSpace(opakowanie.Items)) == 0 {
			return []json.RawMessage{tresc}, nil
		}
		return elementyZlecen(opakowanie.Items)
	}
	var wykaz []json.RawMessage
	if err := json.Unmarshal(tresc, &wykaz); err != nil {
		return nil, bladZlecen(err)
	}
	return wykaz, nil
}

// odczytajZlecenie czyta jeden element wykazu: napis niesie samą treść
// zlecenia, obiekt niesie tytuł i treść.
func odczytajZlecenie(element json.RawMessage) (zlecenieKolejki, error) {
	tresc := bytes.TrimSpace(element)
	if len(tresc) == 0 {
		return zlecenieKolejki{}, nil
	}
	if tresc[0] == '"' {
		var napis string
		if err := json.Unmarshal(tresc, &napis); err != nil {
			return zlecenieKolejki{}, bladZlecen(err)
		}
		return zlecenieKolejki{Content: napis}, nil
	}
	var zlecenie zlecenieKolejki
	if err := json.Unmarshal(tresc, &zlecenie); err != nil {
		return zlecenieKolejki{}, bladZlecen(err)
	}
	return zlecenie, nil
}

// pozycjaZeZlecenia składa wiersz pozycji. Zlecenie bez tytułu bierze go
// z pierwszego wiersza treści, a zlecenie puste jest pomijane — pusty wpis
// w wykazie nie ma czego wykonać i nie jest powodem odmowy.
func pozycjaZeZlecenia(zlecenie zlecenieKolejki) (dane.Pozycja, bool) {
	tytul := strings.TrimSpace(zlecenie.Title)
	tresc := strings.TrimSpace(zlecenie.Content)
	if tytul == "" {
		tytul = tytulZTresci(tresc)
	}
	if tytul == "" {
		return dane.Pozycja{}, false
	}
	pozycja := dane.Pozycja{Tytul: tytul, Stan: stanPozycjiOczekuje}
	if tresc != "" {
		pozycja.TrescZlecenia = &tresc
	}
	return pozycja, true
}

// tytulZTresci wyprowadza etykietę pozycji z pierwszego wiersza polecenia.
func tytulZTresci(tresc string) string {
	wiersz, _, _ := strings.Cut(tresc, "\n")
	wiersz = strings.TrimSpace(wiersz)
	znaki := []rune(wiersz)
	if len(znaki) > dlugoscTytuluZTresci {
		return strings.TrimSpace(string(znaki[:dlugoscTytuluZTresci]))
	}
	return wiersz
}

// bladZlecen nazywa odmowę odczytu wykazu zleceń.
func bladZlecen(err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"nieczytelny wykaz zleceń kolejki: "+err.Error()))
}
