package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// KanalAPI jest generycznym kanałem sieciowym. Nie zna żadnego dostawcy:
// adres, model, odwołanie do klucza, nagłówki, kształt ciała żądania i ścieżka
// do treści w odpowiedzi pochodzą z wiersza rejestru. Nowy dostawca to nowy
// wiersz danych, nie nowy typ w kodzie.
//
// Parametry wiersza (kolumna parametry_json):
//
//	base_url            — pełny adres punktu końcowego (wymagany)
//	strumien            — czy żądać strumienia; domyślnie prawda
//	sciezka_tekstu      — ścieżka do porcji tekstu w zdarzeniu strumienia
//	sciezka_odpowiedzi  — ścieżka do treści w odpowiedzi bez strumienia
//	sciezka_bledu       — ścieżka do komunikatu błędu w odpowiedzi
//	naglowek_klucza     — nazwa nagłówka z danymi dostępowymi
//	przedrostek_klucza  — przedrostek wartości tego nagłówka
//	naglowki            — obiekt dodatkowych nagłówków
//	cialo_dodatkowe     — obiekt scalany z ciałem żądania
//	limit_sekund        — czas oczekiwania na odpowiedź
type KanalAPI struct {
	def    Definicja
	klient *http.Client
}

// NowyKanalAPI buduje kanał sieciowy z wiersza rejestru. Brak adresu jest
// jedynym powodem odmowy budowy — bez niego kanał nie ma dokąd wysłać żądania.
func NowyKanalAPI(d Definicja) (Kanal, error) {
	if strings.TrimSpace(d.Parametr("base_url")) == "" {
		return nil, fmt.Errorf("models: kanał %q bez parametru base_url", d.Kod)
	}
	return &KanalAPI{def: d, klient: &http.Client{Timeout: limitCzasu(d)}}, nil
}

// limitCzasu odczytuje czas oczekiwania na odpowiedź; brak parametru daje
// dziesięć minut, co mieści długie odpowiedzi modeli rozumujących.
func limitCzasu(d Definicja) time.Duration {
	sekundy := 0
	if wskazane := d.Parametr("limit_sekund"); wskazane != "" {
		fmt.Sscanf(wskazane, "%d", &sekundy)
	}
	if sekundy <= 0 {
		return 10 * time.Minute
	}
	return time.Duration(sekundy) * time.Second
}

// Kod zwraca kod kanału z wiersza rejestru.
func (k *KanalAPI) Kod() string {
	return k.def.Kod
}

// Definicja zwraca wiersz rejestru, z którego kanał powstał.
func (k *KanalAPI) Definicja() Definicja {
	return k.def
}

// Wyslij nadaje prowenancję wywołania, wykonuje żądanie i przekazuje odpowiedź
// jako fragmenty tekstu. Klucz dostępowy nie trafia do prowenancji — jedzie tam
// wyłącznie nazwa odwołania.
func (k *KanalAPI) Wyslij(ctx context.Context, z Zapytanie, u Ujscie) error {
	prowenancja := ProwenancjaZapytania(z, k.def)
	prowenancja.Adres = k.def.Parametr("base_url")
	if err := NadajProwenancje(ctx, u, z, prowenancja); err != nil {
		return err
	}
	if k.def.PoswiadczenieOdwolanie != "" {
		konto := MetadaneKonta{Konto: z.WybraneKonto(k.def), Odwolanie: k.def.PoswiadczenieOdwolanie}
		if err := NadajKonto(ctx, u, z, konto); err != nil {
			return err
		}
	}

	zadanie, err := k.zbudujZadanie(ctx, z)
	if err != nil {
		return err
	}
	odpowiedz, err := k.klient.Do(zadanie)
	if err != nil {
		return fmt.Errorf("models: kanał %s: %w", k.def.Kod, err)
	}
	defer odpowiedz.Body.Close()

	if odpowiedz.StatusCode < 200 || odpowiedz.StatusCode > 299 {
		return bladOdpowiedzi(k.def, odpowiedz)
	}
	if !k.strumieniowy() {
		return k.nadajCalosc(ctx, z, u, odpowiedz.Body)
	}
	return k.nadajStrumien(ctx, z, u, odpowiedz.Body)
}

// strumieniowy odpowiada, czy kanał żąda odpowiedzi strumieniem.
func (k *KanalAPI) strumieniowy() bool {
	return k.def.ParametrLub("strumien", "true") != "false"
}

// nadajStrumien przenosi zdarzenia strumienia na fragmenty tekstu.
func (k *KanalAPI) nadajStrumien(ctx context.Context, z Zapytanie, u Ujscie, tresc io.Reader) error {
	sciezka := k.def.ParametrLub("sciezka_tekstu", "choices.0.delta.content")
	return CzytajZdarzenia(tresc, func(zdarzenie []byte) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		porcja, jest := WartoscZeSciezki(zdarzenie, sciezka)
		if !jest {
			return nil
		}
		return NadajTekst(ctx, u, z, porcja)
	})
}

// nadajCalosc przenosi odpowiedź bez strumienia na jeden fragment tekstu.
func (k *KanalAPI) nadajCalosc(ctx context.Context, z Zapytanie, u Ujscie, tresc io.Reader) error {
	surowe, err := io.ReadAll(tresc)
	if err != nil {
		return fmt.Errorf("models: kanał %s: odczyt odpowiedzi: %w", k.def.Kod, err)
	}
	sciezka := k.def.ParametrLub("sciezka_odpowiedzi", "choices.0.message.content")
	odpowiedz, jest := WartoscZeSciezki(surowe, sciezka)
	if !jest {
		return fmt.Errorf("models: kanał %s: odpowiedź bez treści pod ścieżką %s", k.def.Kod, sciezka)
	}
	return NadajTekst(ctx, u, z, odpowiedz)
}

// bladOdpowiedzi składa błąd z kodu stanu i komunikatu dostawcy. Komunikat
// czyta ścieżką z wiersza rejestru, a gdy go tam nie ma — surową treścią.
// Wolna funkcja, nie metoda: tej samej drogi używa kanał obrazowy, a błąd
// dostawcy ma brzmieć tak samo niezależnie od tego, co kanał oddaje.
func bladOdpowiedzi(d Definicja, odpowiedz *http.Response) error {
	surowe, _ := io.ReadAll(io.LimitReader(odpowiedz.Body, 8<<10))
	komunikat := strings.TrimSpace(string(surowe))
	if sciezka := d.Parametr("sciezka_bledu"); sciezka != "" {
		if wyjete, jest := WartoscZeSciezki(json.RawMessage(surowe), sciezka); jest {
			komunikat = wyjete
		}
	}
	return fmt.Errorf("models: kanał %s: odpowiedź %d: %s", d.Kod, odpowiedz.StatusCode, komunikat)
}
