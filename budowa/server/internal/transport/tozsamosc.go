// Odpowiedzialność pliku: tożsamość połączenia — komplet faktów o tym, kto
// stoi po drugiej stronie gniazda, zebranych w jednym miejscu i wnoszonych do
// rdzenia jedną drogą.
//
// Jeden mechanizm, nie drugi obok.
// Tożsamość mieszka przy połączeniu — tam, gdzie już mieszka konto — i wchodzi
// do kontekstu żądania tą samą jedną drogą, którą wchodził identyfikator
// gniazda (`wejscieTransportu.Obsluz`). Nie powstaje ani drugi wpis kontekstu,
// ani drugi rejestr, ani pole w kopercie kontraktu: koperta jest zamrożona, a
// tożsamość i tak nie jest własnością komunikatu, tylko własnością łącza.
//
// Skąd się biorą fakty — dwa źródła, jeden skład:
//  1. Nawiązanie. Serwer narzędzi modelu przedstawia się parametrami zapytania
//     przy zestawianiu gniazda — tą samą drogą, którą urządzenie od zawsze
//     wskazuje konto (`kontoZadania`). Musi tak, bo powitania nie wysyła: jest
//     klientem wołającym komendy, a nie oknem interfejsu.
//  2. Powitanie. Okno interfejsu niesie `clientId` w `connection.hello` —
//     kontrakt tak mówi od pierwszego wydania i pole to już tam jest. Transport
//     odczytuje je z ładunku powitania i dokłada do tożsamości połączenia.
//
// Transport niczego nie rozstrzyga. Nie zna pojęcia `ActorKind`, nie wie, co to
// operator ani asystent, i nie ma w tym pliku ani jednej wartości wyliczenia
// kontraktu. Niesie fakty; rozstrzygnięcie „czyja to ręka" należy do rdzenia
// (`core/sprawca.go`), bo tylko rdzeń zna rolę okna.
//
// To nie jest uprawnienie. Tożsamość nie rozstrzyga ani razu, czy coś wolno —
// tym zajmuje się wyłącznie straż bramki (`bramka.go`) i pyta o zupełnie co
// innego. Parametr podany przez wołającego jest tu opisem, więc jego podrobienie
// niczego nie otwiera; gdyby cokolwiek od niego zależało, byłby bramką.
package transport

import (
	"net/http"
	"strings"
)

const (
	// ParametrKlienta nazywa parametr zapytania i nagłówek, którymi klient
	// przedstawia swój identyfikator już przy nawiązaniu. Okno interfejsu go nie
	// używa — jemu wystarczy powitanie — ale klient bez powitania (serwer
	// narzędzi modelu) innej drogi nie ma.
	ParametrKlienta = "klient"
	// NaglowekKlienta jest nagłówkową postacią ParametrKlienta, wzorem konta.
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
	// RodzajNarzedzi oznacza gniazdo serwera narzędzi modelu. Stała stoi tutaj,
	// a nie w pakiecie `narzedzia`, żeby obie strony rozmowy — ta, która napis
	// wysyła, i ta, która go czyta — brały go z jednego miejsca.
	RodzajNarzedzi = "narzedzia"
)

// Tozsamosc jest kompletem faktów o drugiej stronie gniazda.
//
// Każde pole może być puste i pustka jest odpowiedzią, nie usterką: „nie
// wiadomo" to stan zwykły dla gniazda, które jeszcze się nie przedstawiło.
// Rdzeń ma to milczenie przenieść dalej — kontrakt mówi o polu
// `actor` wprost: „brak znaczy, że rdzeń nie potrafił tego rozstrzygnąć".
type Tozsamosc struct {
	// IdPolaczenia jest identyfikatorem gniazda nadanym przy nawiązaniu.
	IdPolaczenia string
	// IdKlienta jest identyfikatorem klienta: z powitania albo z nawiązania.
	IdKlienta string
	// Rodzaj mówi, czym jest program po drugiej stronie. Puste znaczy klienta
	// nieprzedstawionego rodzajem — czyli okno interfejsu.
	Rodzaj string
	// Zasieg jest rolą okna serwera narzędzi. Puste dla klienta, który rolą się
	// nie przedstawił.
	Zasieg string
	// IdOkna jest oknem rozmowy serwera narzędzi.
	IdOkna string
}

// Narzedzia mówi, czy gniazdo należy do serwera narzędzi modelu.
func (t Tozsamosc) Narzedzia() bool {
	return t.Rodzaj == RodzajNarzedzi
}

// tozsamoscZadania czyta tożsamość przedstawioną przy nawiązaniu.
//
// Parametr zapytania stoi przed nagłówkiem — tak samo jak przy koncie — bo
// klient WebSocket w przeglądarce nagłówków ustawić nie potrafi, a klient
// biblioteczny potrafi obu dróg. Brak wskazania nie odrzuca nawiązania:
// gniazdo bez tożsamości pracuje dalej, tyle że jego zdarzenia pójdą bez
// sprawcy.
func tozsamoscZadania(id string, r *http.Request) Tozsamosc {
	tozsamosc := Tozsamosc{IdPolaczenia: id}
	if r == nil {
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

// ParametryTozsamosci składa parametry zapytania, którymi klient przedstawia
// się przy nawiązaniu.
//
// Funkcja stoi tutaj, przy odczycie, a nie po stronie klienta, bo napis
// parametru ma mieć jedno źródło dla piszącego i czytającego. Klient
// (`narzedzia/adres.go`) podaje wartości; nazw nie zna i nie powtarza.
// Wartości puste nie wchodzą — parametr pusty i nieobecny mają znaczyć to samo.
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
