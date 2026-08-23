// Odpowiedzialność pliku: sprawdzenie, czy silnik mowy ma czym pracować,
// i odczyt odpowiedzi, którą pomocnik na to pytanie daje.
//
// Sprawdzenie poprzedza mikrofon: klient pyta o gotowość, zanim narysuje
// przycisk nagrywania, zamiast tłumaczyć jego milczenie po nieudanej próbie.
//
// Brak silnika nie jest awarią i nie wraca błędem. Pomocnik pytany o gotowość
// kończy się powodzeniem także wtedy, gdy silnika nie ma — oddaje wówczas
// `gotowy` fałszywe i powód. Błędem jest dopiero brak odpowiedzi pomocnika:
// nie dało się go uruchomić albo odpowiedział czymś, co nie jest jego
// odpowiedzią. Pierwsze naprawia się instalacją, drugie — zgłoszeniem usterki.
//
// Dwa braki są rozróżniane osobno, bo mają dwie różne naprawy: nie ma czym
// uruchomić skryptu (interpreter) kontra skrypt się uruchomił, lecz nie zastał
// silnika rozpoznawania.
package mowa

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
)

// granicaSprawdzenia to granica czasu pytania o gotowość.
//
// Pytanie o gotowość nie ładuje wag i nie dekoduje dźwięku, więc pomocnik
// odpowiada w ułamku sekundy. Dłuższe milczenie oznacza interpreter, który
// utknął; granica zamienia je na odmowę zamiast zawieszonego ekranu przed
// pierwszym nagraniem.
const granicaSprawdzenia = 20 * time.Second

// przelacznikWersji to tryb pomocnika, w którym pyta się go o samą gotowość.
const przelacznikWersji = "--wersja"

// Dostepnosc jest odpowiedzią pomocnika na pytanie, czy ma czym rozpoznać mowę.
//
// Pola przepisane co do znaku z odpowiedzi pomocnika
// (pomocniki/transkrypcja/transkrypcja.py, tryb `--wersja`). Klucze są polskie,
// bo pomocnik jest częścią tego produktu, a produkt jest polskojęzyczny — to nie
// jest nazewnictwo kontraktu, którego stałe zostają angielskie.
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

// Dostepnosc pyta pomocnika o gotowość i oddaje jego odpowiedź.
//
// Odmowa uruchomienia pomocnika jest zamieniana na odpowiedź `Gotowy` fałszywe
// z powodem, a nie na błąd. Brak interpretera to ten sam rodzaj wiadomości dla
// Operatora, co brak biblioteki — jedno pytanie daje jedną odpowiedź,
// niezależnie od tego, na którym ogniwie łańcuch się urwał. Błędem zostaje
// wyłącznie odpowiedź nieczytelna: pomocnik odezwał się czymś, co nie jest jego
// odpowiedzią.
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

// odczytajDostepnosc rozbiera odpowiedź pomocnika.
//
// Odpowiedź czytana jest z wyjścia, nigdy z diagnostyki: pomocnik pisze na
// diagnostykę ostrzeżenia bibliotek, których nie kontroluje, a na wyjście
// wyłącznie swój JSON. Wyjście puste znaczy, że pomocnik nie doszedł do
// wypisania odpowiedzi — wtedy jedyną wiadomością jest diagnostyka i to ona
// idzie w błędzie zamiast zdania o niepoprawnym JSON-ie, które niczego nie
// tłumaczy.
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
			"; naprawa: zgłosić usterkę pomocnika — rdzeń oczekuje odpowiedzi w JSON"}
	}

	// Odmowa bez powodu zostawiłaby Operatora z „nie da się” bez zdania, co z tym
	// zrobić. Powód zastępczy mówi wprost, że pomocnik go nie podał.
	if !odpowiedz.Gotowy && strings.TrimSpace(odpowiedz.Powod) == "" {
		odpowiedz.Powod = "pomocnik transkrypcji zgłosił brak gotowości i nie podał powodu" +
			ogonDiagnostyki(wynik.Diagnostyka)
	}
	return odpowiedz, nil
}

// powodZUruchomienia składa powód odmowy z błędu uruchomienia i diagnostyki.
//
// Diagnostyka dokładana jest do treści błędu, bo w niej stoi zwykle jedyne
// zdanie mówiące, co poszło nie tak po stronie interpretera — samo „nie można
// uruchomić” nie wskazuje naprawy.
func powodZUruchomienia(err error, wynik Wynik) string {
	return err.Error() + ogonDiagnostyki(wynik.Diagnostyka)
}

// ogonDiagnostyki dokłada wyjście diagnostyczne, gdy jest czym dołożyć.
func ogonDiagnostyki(diagnostyka string) string {
	tresc := strings.TrimSpace(diagnostyka)
	if tresc == "" {
		return ""
	}
	return " — pomocnik powiedział: " + tresc
}
