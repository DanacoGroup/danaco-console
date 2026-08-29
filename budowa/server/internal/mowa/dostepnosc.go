// Odpowiedzialność pliku: sprawdzenie, czy silnik mowy ma czym pracować,
// i odczyt odpowiedzi, którą pomocnik na to pytanie daje; brak silnika nie
// jest awarią.
package mowa

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
)

// granicaSprawdzenia to granica czasu pytania o gotowość; pytanie nie ładuje
// wag ani nie dekoduje dźwięku, więc dłuższe milczenie oznacza interpreter,
// który utknął.
const granicaSprawdzenia = 20 * time.Second

// przelacznikWersji to tryb pomocnika transkrypcji, w którym pyta się go
// wyłącznie o gotowość, bez ładowania modelu rozpoznawania.
const przelacznikWersji = "--wersja"

// Dostepnosc jest odpowiedzią pomocnika na pytanie, czy ma czym rozpoznać
// mowę; pola przepisane co do znaku z odpowiedzi pomocnika transkrypcji,
// tryb `--wersja`.
type Dostepnosc struct {
	// Gotowy — czy transkrypcja ruszy tu i teraz.
	Gotowy bool `json:"gotowy"`
	// Python — wersja interpretera, który odpowiedział.
	Python string `json:"python"`
	// Silnik — nazwa i wersja biblioteki rozpoznawania; pusta, gdy jej nie ma.
	Silnik string `json:"silnik"`
	// Model — rozmiar modelu zastanego na dysku; pusty, gdy go nie pobrano.
	Model string `json:"model"`
	// Powod — trójczęściowa odmowa; wypełniony wyłącznie przy Gotowy fałszywym.
	Powod string `json:"powod"`
}

// Dostepnosc pyta pomocnika o gotowość i oddaje jego odpowiedź; odmowa
// uruchomienia pomocnika jest zamieniana na odpowiedź `Gotowy` fałszywe
// z powodem, a nie na błąd.
func (s *Silnik) Dostepnosc(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar) (Dostepnosc, error) {

	pomocnik, err := OdnajdzPomocnika(s.ustawienia.Program)
	if err != nil {
		return Dostepnosc{Powod: err.Error()}, nil
	}

	wynik, err := Uruchom(ctx, s.uruchamiacz, okno, zasady, obszar,
		pomocnik, pomocnik.Argumenty(przelacznikWersji), s.katalogPracy(obszar),
		granicaSprawdzenia)
	if err != nil {
		return Dostepnosc{Powod: powodZUruchomienia(err, wynik)}, nil
	}
	return odczytajDostepnosc(wynik)
}

// odczytajDostepnosc rozbiera odpowiedź pomocnika; odpowiedź czytana jest
// z wyjścia, nigdy z diagnostyki, na którą pomocnik pisze ostrzeżenia
// bibliotek.
func odczytajDostepnosc(wynik Wynik) (Dostepnosc, error) {
	tresc := strings.TrimSpace(string(wynik.Wyjscie))
	if tresc == "" {
		return Dostepnosc{}, &BrakSilnika{Powod: "pomocnik transkrypcji nie oddał odpowiedzi" +
			" o gotowości" + ogonDiagnostyki(wynik.Diagnostyka) +
			"; naprawa: uruchomić pomocnika ręcznie z przełącznikiem " + przelacznikWersji}
	}

	var odpowiedz Dostepnosc
	if err := json.Unmarshal([]byte(tresc), &odpowiedz); err != nil {
		return Dostepnosc{}, &BrakSilnika{Powod: "odpowiedź pomocnika transkrypcji" +
			" jest nieczytelna (" + err.Error() + ")" + ogonDiagnostyki(wynik.Diagnostyka) +
			"; naprawa: zgłosić usterkę pomocnika — serwer oczekuje odpowiedzi w JSON"}
	}

	// Powód zastępczy mówi wprost, że pomocnik nie podał powodu odmowy.
	if !odpowiedz.Gotowy && strings.TrimSpace(odpowiedz.Powod) == "" {
		odpowiedz.Powod = "pomocnik transkrypcji zgłosił brak gotowości i nie podał powodu" +
			ogonDiagnostyki(wynik.Diagnostyka)
	}
	return odpowiedz, nil
}

// powodZUruchomienia składa powód odmowy z błędu uruchomienia i diagnostyki,
// w której stoi zwykle jedyne zdanie mówiące, co poszło nie tak.
func powodZUruchomienia(err error, wynik Wynik) string {
	return err.Error() + ogonDiagnostyki(wynik.Diagnostyka)
}

// ogonDiagnostyki dokłada wyjście diagnostyczne pomocnika do treści powodu
// odmowy, gdy jest czym dołożyć.
func ogonDiagnostyki(diagnostyka string) string {
	tresc := strings.TrimSpace(diagnostyka)
	if tresc == "" {
		return ""
	}
	return " — pomocnik powiedział: " + tresc
}
