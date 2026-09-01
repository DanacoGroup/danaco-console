// Pakiet transport niesie tożsamość połączenia — komplet faktów o tym, kto stoi po drugiej stronie gniazda, wnoszonych do rdzenia jedną drogą.
package transport

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
)

const (
	// ParametrKlienta nazywa parametr zapytania i nagłówek, którymi klient
	// przedstawia swój identyfikator już przy nawiązaniu. Okno interfejsu go nie
	// używa — jemu wystarczy powitanie — ale klient bez powitania (serwer
	// narzędzi modelu) innej drogi nie ma.
	ParametrKlienta = "klient"
	// NaglowekKlienta jest nagłówkową postacią stałej ParametrKlienta, przyjętą wzorem stałej dotyczącej konta.
	NaglowekKlienta = "X-Danaco-Klient"
	// ParametrRodzaju niesie rodzaj klienta: czym jest program po drugiej
	// stronie gniazda. Wartość znaną transportowi jest jedna — RodzajNarzedzi.
	ParametrRodzaju = "rodzaj"
	// ParametrZasiegu niesie rolę okna, w którego imieniu pracuje serwer
	// narzędzi. Wartości nazywa `narzedzia/zasieg_roli.go`; transport ich nie
	// zna i nie porównuje — przenosi napis, bo import w tę stronę odwróciłby
	// zależność.
	ParametrZasiegu = "zasieg"
	// ParametrOkna niesie okno rozmowy serwera narzędzi. Do rozstrzygnięcia
	// sprawcy niepotrzebne, do zrozumienia dziennika — konieczne.
	ParametrOkna = "okno"
	// ParametrPoswiadczenia niesie poświadczenie serwera narzędzi wydane przez
	// rdzeń przy uruchomieniu. Bez zgodnego poświadczenia rodzaj i zasięg z
	// zapytania nie nadają niczego — nadawałby je wtedy każdy proces maszyny.
	ParametrPoswiadczenia = "poswiadczenie"
	// RodzajNarzedzi oznacza gniazdo serwera narzędzi modelu. Stała stoi tutaj,
	// a nie w pakiecie `narzedzia`, żeby obie strony rozmowy — ta, która napis
	// wysyła, i ta, która go czyta — brały go z jednego miejsca.
	RodzajNarzedzi = "narzedzia"
)

// Tozsamosc jest kompletem faktów o drugiej stronie gniazda; każde pole może być puste, i pustka jest odpowiedzią, nie usterką.
type Tozsamosc struct {
	// IdPolaczenia jest identyfikatorem gniazda nadanym przy nawiązaniu.
	IdPolaczenia string
	// IdKlienta jest identyfikatorem klienta: z powitania albo z nawiązania.
	IdKlienta string
	// Rodzaj mówi, czym jest program po drugiej stronie; pustka znaczy klienta bez rodzaju.
	Rodzaj string
	// Zasieg jest rolą okna serwera narzędzi. Puste dla klienta, który rolą się
	// nie przedstawił.
	Zasieg string
	// IdOkna jest oknem rozmowy serwera narzędzi.
	IdOkna string
	// PoswiadczenieSprawdzone mówi, że transport porównał poświadczenie z
	// nawiązania z poświadczeniem wydanym przez rdzeń. Sam napis tu nie wchodzi:
	// tożsamość jedzie do dziennika i do zdarzeń, a sekret nie ma tam czego szukać.
	PoswiadczenieSprawdzone bool
}

// Narzedzia mówi, czy gniazdo należy do serwera narzędzi modelu. Rodzaj bez
// sprawdzonego poświadczenia nie liczy się — to napis z zapytania.
func (t Tozsamosc) Narzedzia() bool {
	return t.PoswiadczenieSprawdzone && t.Rodzaj == RodzajNarzedzi
}

// tozsamoscZadania czyta tożsamość przedstawioną przy nawiązaniu. Bez
// sprawdzonego poświadczenia zostaje sam identyfikator gniazda: parametry
// nawiązania przedstawiają wyłącznie serwer narzędzi, okno przedstawia się
// powitaniem.
func tozsamoscZadania(id string, r *http.Request, poswiadczone bool) Tozsamosc {
	tozsamosc := Tozsamosc{IdPolaczenia: id, PoswiadczenieSprawdzone: poswiadczone}
	if r == nil || !poswiadczone {
		return tozsamosc
	}
	zapytanie := r.URL.Query()
	tozsamosc.IdKlienta = strings.TrimSpace(zapytanie.Get(ParametrKlienta))
	if tozsamosc.IdKlienta == "" {
		tozsamosc.IdKlienta = strings.TrimSpace(r.Header.Get(NaglowekKlienta))
	}
	tozsamosc.Rodzaj = strings.TrimSpace(zapytanie.Get(ParametrRodzaju))
	tozsamosc.Zasieg = strings.TrimSpace(zapytanie.Get(ParametrZasiegu))
	tozsamosc.IdOkna = strings.TrimSpace(zapytanie.Get(ParametrOkna))
	return tozsamosc
}

// Funkcja ParametryTozsamosci składa parametry zapytania, którymi klient przedstawia się przy nawiązaniu połączenia.
func ParametryTozsamosci(t Tozsamosc) map[string]string {
	parametry := map[string]string{}
	for nazwa, wartosc := range map[string]string{
		ParametrKlienta: t.IdKlienta,
		ParametrRodzaju: t.Rodzaj,
		ParametrZasiegu: t.Zasieg,
		ParametrOkna:    t.IdOkna,
	} {
		if wartosc != "" {
			parametry[nazwa] = wartosc
		}
	}
	return parametry
}

// Poświadczenie serwera narzędzi tego procesu: wydawane raz, przy pierwszym
// pytaniu, i trzymane do końca procesu.
var (
	poswiadczenieRaz     sync.Once
	poswiadczenieProcesu string
)

// PoswiadczenieNarzedzi oddaje poświadczenie serwera narzędzi tego procesu.
// Rdzeń wręcza je serwerowi narzędzi w argumentach uruchomienia, a transport
// porównuje z nim napis z nawiązania. Napis pusty znaczy brak źródła losowego:
// nawiązanie narzędzi jest wtedy odrzucane, nie przepuszczane.
func PoswiadczenieNarzedzi() string {
	poswiadczenieRaz.Do(func() {
		losowe := make([]byte, 32)
		if _, err := rand.Read(losowe); err != nil {
			return
		}
		poswiadczenieProcesu = hex.EncodeToString(losowe)
	})
	return poswiadczenieProcesu
}
