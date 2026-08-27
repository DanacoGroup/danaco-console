// Porównanie wektorów: zamiana wektora na bajty i z powrotem oraz wyłonienie
// fragmentów najbliższych pytaniu.
//
// Podobieństwo liczy się w Go, a nie w bazie. Produkt wozi sterownik
// `modernc.org/sqlite`, czyli SQLite przepisany na Go, bez zależności od C.
// Rozszerzenie wektorowe `sqlite-vec` wydawane jest jako biblioteka natywna
// (`.so`, `.dylib`, `.dll`) ładowana przez `sqlite3_load_extension`, a sterownik
// nie wystawia drogi włączenia ładowania rozszerzeń — `load_extension('vec0')`
// kończy się `SQL logic error: not authorized (1)`. Nie miałby też jak jej
// wystawić: kod przepisany na Go nie wciągnie do siebie natywnej biblioteki C.
// Drogą alternatywną byłaby zamiana sterownika na wariant z CGO w całym
// produkcie.
//
// Zostaje przegląd zupełny: wektory leżą w tabeli, a podobieństwo liczy się w Go
// dla każdego wiersza. Przy 768 wymiarach i jednym rdzeniu przegląd stu tysięcy
// wektorów zajmuje rząd sześćdziesięciu milisekund, a biblioteka rzędu tysiąca
// dokumentów daje kilkanaście tysięcy fragmentów, czyli kilkanaście milisekund na
// zapytanie — wobec sekund, które zajmuje wczytanie wag przez pomocnika (patrz
// nagłówek `silnik.go`). Indeks przybliżony zaczyna się opłacać dopiero powyżej
// setek tysięcy pozycji.
//
// Wektory wchodzą do bazy znormalizowane. Podobieństwo kosinusowe to iloczyn
// skalarny podzielony przez iloczyn długości, a długość wektora dokumentu nie
// zmienia się między zapytaniami; wektor podzielony przez własną długość już
// przy zapisie sprawia, że porównanie jest samym iloczynem skalarnym, a wynik
// wprost kosinusem.
package wiedza

import (
	"encoding/binary"
	"errors"
	"math"
	"sort"
)

// rozmiarLiczby to długość jednej współrzędnej w zapisie bajtowym.
//
// Cztery bajty, nie osiem: model oddaje liczby o pojedynczej precyzji, więc
// osiem bajtów zapisywałoby zera po przecinku, których w wektorze nie ma —
// dwukrotność miejsca i dwukrotność odczytu z dysku za zero dokładności.
const rozmiarLiczby = 4

// Znormalizuj dzieli wektor przez jego długość i oddaje kopię.
//
// Wektor zerowy wraca bez zmiany, a nie przez dzielenie przez zero: jego
// iloczyn skalarny z czymkolwiek wynosi zero, więc nigdy nie wygra rankingu,
// podczas gdy wektor z NaN psułby porównania.
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
// pojedynczej precyzji, porządek bajtów mało-końcowy.
//
// Porządek ustalony jawnie, a nie odziedziczony po maszynie: baza bywa
// przenoszona między maszynami razem z katalogiem danych, a wektor odczytany
// w odwrotnym porządku bajtów to nie gorsze trafienia, tylko liczby bez żadnego
// związku z tekstem.
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

// blisko liczy podobieństwo dwóch wektorów znormalizowanych.
//
// Wektory różnej długości dają zero, a nie panikę ani porównanie części
// wspólnej: różna długość znaczy różne modele, więc te dwa wektory nie leżą
// w jednej przestrzeni i porównanie ich części byłoby liczbą udającą wynik.
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

// Trafienie wiąże pozycję wskaźnika z jej podobieństwem do pytania.
type Trafienie struct {
	Pozycja    Pozycja
	Podobienst float32
}

// Najblizsze przegląda wszystkie pozycje i oddaje `ile` najbliższych pytaniu,
// od najbliższej. Pozycje o podobieństwie niedodatnim odpadają: kosinus
// niedodatni znaczy, że fragment nie ma z pytaniem nic wspólnego, a oddanie go
// jako trafienia dałoby cytat do zbudowania odpowiedzi z niczego.
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
// wejściową przy równej ocenie.
//
// Sortowanie stabilne, żeby dwa fragmenty o równej ocenie wracały zawsze w tej
// samej kolejności — inaczej to samo pytanie zadane dwa razy dawałoby dwie różne
// odpowiedzi bez żadnej zmiany w wiedzy. Wspólne dla obu przebiegów: pierwszy
// układa po kosinusie, drugi po ocenie kodera (`przesiew.go`), a wymóg
// powtarzalności jest ten sam.
func posortujMalejaco(trafienia []Trafienie) {
	sort.SliceStable(trafienia, func(i, j int) bool {
		return trafienia[i].Podobienst > trafienia[j].Podobienst
	})
}

// WSetnych przelicza kosinus na pole `score` kontraktu („trafnosc w setnych —
// 100 znaczy najblizszy mozliwy"). Zaokrąglenie w górę od połowy, granice
// domknięte — kosinus bywa minimalnie większy od jedynki przez błąd
// zaokrągleń pojedynczej precyzji, a `score` większe od stu byłoby liczbą spoza
// zakresu obiecanego przez kontrakt.
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
