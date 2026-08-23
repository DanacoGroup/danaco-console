package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"danacoconsole/shared"
)

// Definicja jest wierszem rejestru kanałów — odwzorowaniem rekordu tabeli
// kanal_modelu. Rejestr powstaje z takich wierszy w czasie działania, więc nowy
// kanał znaczy nowy wiersz danych, nie nowy typ w kodzie.
//
// PoswiadczenieOdwolanie jest odwołaniem do danych dostępowych — nazwą zmiennej
// środowiskowej albo pozycji magazynu. Sekret nie żyje ani w bazie, ani
// w repozytorium.
type Definicja struct {
	Id                     int64
	Kod                    string
	Nazwa                  string
	Dostawca               string
	Model                  string
	Rodzaj                 string
	KontoId                int64
	PoswiadczenieOdwolanie string
	Parametry              map[string]any
	Multimodalny           bool
	Aktywny                bool
	Kolejnosc              int
	Utworzono              string
}

// KluczAdaptera wskazuje fabrykę, która zbuduje adapter dla tego wiersza.
// Rozstrzyga parametr "adapter", a w jego braku kolumna rodzaj_kanalu. Obie
// wartości są danymi, więc kanał na nowym adapterze powstaje wpisem do tabeli.
func (d Definicja) KluczAdaptera() string {
	if wskazany := strings.TrimSpace(d.Parametr("adapter")); wskazany != "" {
		return wskazany
	}
	return strings.TrimSpace(d.Rodzaj)
}

// Parametr odczytuje parametr kanału jako tekst. Brak parametru daje pusty
// napis — wywołujący dobiera wtedy wartość domyślną, zamiast odmawiać pracy.
func (d Definicja) Parametr(klucz string) string {
	wartosc, jest := d.Parametry[klucz]
	if !jest || wartosc == nil {
		return ""
	}
	switch typowa := wartosc.(type) {
	case string:
		return typowa
	case bool:
		return strconv.FormatBool(typowa)
	case float64:
		return strconv.FormatFloat(typowa, 'f', -1, 64)
	case json.Number:
		return typowa.String()
	default:
		surowe, err := json.Marshal(typowa)
		if err != nil {
			return ""
		}
		return string(surowe)
	}
}

// ParametrJest odczytuje parametr wraz z informacją, czy wiersz go zawiera.
// Odróżnia parametr nieustawiony od ustawionego pustą wartością — pusty
// przedrostek klucza jest świadomym ustawieniem, nie brakiem.
func (d Definicja) ParametrJest(klucz string) (string, bool) {
	if _, jest := d.Parametry[klucz]; !jest {
		return "", false
	}
	return d.Parametr(klucz), true
}

// ParametrLub odczytuje parametr, a przy jego braku zwraca wartość domyślną.
func (d Definicja) ParametrLub(klucz, domyslna string) string {
	if wartosc := strings.TrimSpace(d.Parametr(klucz)); wartosc != "" {
		return wartosc
	}
	return domyslna
}

// Identyfikator zwraca identyfikator kanału w postaci używanej przez kontrakt.
func (d Definicja) Identyfikator() string {
	return strconv.FormatInt(d.Id, 10)
}

// identyfikatorKontraktu zwraca identyfikator, którym klient posługuje się w
// kolejnych komendach (channel.update, channel.remove, window.create).
//
// Musi to być kod kanału, nie numer wiersza. Rejestr indeksuje kanał pod
// obydwoma kluczami (kluczeKanalu), więc odczyt działa tak czy inaczej — ale
// warstwa danych zna wyłącznie kod, więc numer wiersza dałby klientowi
// identyfikator, którym nie da się nic zrobić.
//
// Numer wiersza zostaje wartością zapasową dla wiersza bez kodu.
func (d Definicja) identyfikatorKontraktu() string {
	if kod := strings.TrimSpace(d.Kod); kod != "" {
		return kod
	}
	return d.Identyfikator()
}

// Kontrakt odwzorowuje wiersz rejestru na kształt kontraktu (shared.Channel),
// który idzie do klienta w odpowiedzi channel.list. Parametry jadą bez zmian —
// zawierają wyłącznie odwołania do danych dostępowych.
func (d Definicja) Kontrakt() shared.Channel {
	kanal := shared.Channel{
		Id:      d.identyfikatorKontraktu(),
		Name:    d.Nazwa,
		Kind:    d.Rodzaj,
		Enabled: d.Aktywny,
	}
	if strings.TrimSpace(d.Model) != "" {
		model := d.Model
		kanal.Model = &model
	}
	if len(d.Parametry) > 0 {
		if surowe, err := json.Marshal(d.Parametry); err == nil {
			kanal.Config = surowe
		}
	}
	kanal.CreatedAt = czasWMilisekundach(d.Utworzono)
	return kanal
}

// czasWMilisekundach przenosi znacznik czasu bazy (ISO 8601, strefa UTC) na
// milisekundy epoki, których używa kontrakt. Znacznik nieczytelny daje zero —
// data pomocnicza nie może wywrócić odpowiedzi.
func czasWMilisekundach(znacznik string) int64 {
	chwila, err := time.Parse(time.RFC3339, strings.TrimSpace(znacznik))
	if err != nil {
		return 0
	}
	return chwila.UnixMilli()
}
