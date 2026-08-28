// Porównanie wektorów: zamiana wektora na bajty i z powrotem oraz wyłonienie
// fragmentów najbliższych pytaniu. Podobieństwo liczy się w Go, a nie w bazie,
// bo sterownik SQLite użyty w produkcie nie wystawia drogi ładowania
// rozszerzeń natywnych.
package wiedza

import (
	"encoding/binary"
	"errors"
	"math"
	"sort"
)

// rozmiarLiczby to długość jednej współrzędnej w zapisie bajtowym. Cztery
// bajty, nie osiem: model oddaje liczby o pojedynczej precyzji.
const rozmiarLiczby = 4

// Znormalizuj dzieli wektor przez jego długość i oddaje kopię. Wektor zerowy
// wraca bez zmiany, a nie przez dzielenie przez zero.
func Znormalizuj(wektor []float32) []float32 {
	var suma float64
	for _, wspolrzedna := range wektor {
		suma += float64(wspolrzedna) * float64(wspolrzedna)
	}
	dlugosc := math.Sqrt(suma)
	kopia := make([]float32, len(wektor))
	if dlugosc == 0 {
		copy(kopia, wektor)
		return kopia
	}
	for i, wspolrzedna := range wektor {
		kopia[i] = float32(float64(wspolrzedna) / dlugosc)
	}
	return kopia
}

// NaBajty zapisuje wektor w postaci, w której leży w bazie: ciąg liczb
// pojedynczej precyzji, porządek bajtów mało-końcowy, ustalony jawnie.
func NaBajty(wektor []float32) []byte {
	bajty := make([]byte, len(wektor)*rozmiarLiczby)
	for i, wspolrzedna := range wektor {
		binary.LittleEndian.PutUint32(bajty[i*rozmiarLiczby:], math.Float32bits(wspolrzedna))
	}
	return bajty
}

// ZBajtow odczytuje wektor zapisany przez NaBajty.
//
// Długość niepodzielna przez rozmiar liczby jest odmową, a nie odczytem
// obciętym: wiersz uszkodzony ma się zgłosić, a nie udawać krótszego wektora,
// który porówna się ze wszystkim i z niczym.
func ZBajtow(bajty []byte) ([]float32, error) {
	if len(bajty)%rozmiarLiczby != 0 {
		return nil, errors.New("wskaźnik znaczenia: wiersz wskaźnika niesie " +
			liczba(len(bajty)) + " bajtów, czyli liczbę niepodzielną przez " +
			liczba(rozmiarLiczby) + " — wektor jest uszkodzony; " +
			"naprawa: przebudować wskaźnik (`knowledge.index` z `rebuild`)")
	}
	wektor := make([]float32, len(bajty)/rozmiarLiczby)
	for i := range wektor {
		wektor[i] = math.Float32frombits(binary.LittleEndian.Uint32(bajty[i*rozmiarLiczby:]))
	}
	return wektor, nil
}

// blisko liczy podobieństwo dwóch wektorów znormalizowanych. Wektory różnej
// długości dają zero, a nie panikę ani porównanie części wspólnej.
func blisko(pierwszy, drugi []float32) float32 {
	if len(pierwszy) != len(drugi) {
		return 0
	}
	var suma float32
	for i := range pierwszy {
		suma += pierwszy[i] * drugi[i]
	}
	return suma
}

// Trafienie wiąże pozycję wskaźnika z jej podobieństwem do pytania, w skali
// kosinusa wektorów znormalizowanych.
type Trafienie struct {
	Pozycja    Pozycja
	Podobienst float32
}

// Najblizsze przegląda wszystkie pozycje i oddaje `ile` najbliższych pytaniu,
// od najbliższej. Pozycje o podobieństwie niedodatnim odpadają.
func Najblizsze(pozycje []Pozycja, pytanie []float32, ile int) []Trafienie {
	trafienia := make([]Trafienie, 0, len(pozycje))
	for _, pozycja := range pozycje {
		podobienstwo := blisko(pozycja.Wektor, pytanie)
		if podobienstwo <= 0 {
			continue
		}
		trafienia = append(trafienia, Trafienie{Pozycja: pozycja, Podobienst: podobienstwo})
	}
	posortujMalejaco(trafienia)
	if ile > 0 && len(trafienia) > ile {
		trafienia = trafienia[:ile]
	}
	return trafienia
}

// posortujMalejaco układa trafienia od najbliższego, zachowując kolejność
// wejściową przy równej ocenie, żeby to samo pytanie dawało zawsze ten sam
// wynik.
func posortujMalejaco(trafienia []Trafienie) {
	sort.SliceStable(trafienia, func(i, j int) bool {
		return trafienia[i].Podobienst > trafienia[j].Podobienst
	})
}

// WSetnych przelicza kosinus na pole `score` kontraktu w setnych. Granice
// domknięte, bo kosinus bywa minimalnie większy od jedynki przez błąd
// zaokrągleń pojedynczej precyzji.
func WSetnych(podobienstwo float32) int {
	setne := int(math.Round(float64(podobienstwo) * 100))
	if setne < 0 {
		return 0
	}
	if setne > 100 {
		return 100
	}
	return setne
}
