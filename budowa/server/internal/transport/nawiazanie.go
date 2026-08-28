package transport

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/coder/websocket"
)

// LimitOdczytu podnosi domyślny limit ramki biblioteki (32 KiB) do rozmiaru,
// w którym mieści się wiadomość z obszernym ładunkiem. Limit biblioteczny
// zrywa połączenie, a zerwanie z powodu długości wiadomości byłoby bramą, której
// kontrakt nie przewiduje.
const LimitOdczytu = 16 << 20

// Metoda nawiaz przyjmuje żądanie uaktualnienia do WebSocket i prowadzi całe życie połączenia od rejestracji do wykreślenia.
func (s *Serwer) nawiaz(w http.ResponseWriter, r *http.Request) {
	gniazdo, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: s.pochodzeniaDozwolone(),
	})
	if err != nil {
		s.ustawienia.Dziennik.Printf("transport: nawiązanie odrzucone: %v", err)
		return
	}
	gniazdo.SetReadLimit(LimitOdczytu)

	id := s.nastepnyId()
	polaczenie := nowePolaczenie(s.kontekst, id, kontoZadania(r), tozsamoscZadania(id, r), gniazdo, s.ustawienia.PojemnoscKolejki, s.ustawienia.Dziennik)
	s.polaczenia.dodaj(polaczenie)
	s.zawiadomPrzylaczono(polaczenie)
	// Tożsamość jest w linii dziennika, bo od niej zależy sprawca zdarzeń widoczny w interfejsie.
	tozsamosc := polaczenie.Tozsamosc()
	s.ustawienia.Dziennik.Printf("transport: przyłączenie %s konto=%s klient=%q rodzaj=%q zasięg=%q okno=%q adres=%s (łącznie %d)",
		polaczenie.Id(), polaczenie.Konto(), tozsamosc.IdKlienta, tozsamosc.Rodzaj,
		tozsamosc.Zasieg, tozsamosc.IdOkna, adresZdalny(r), s.polaczenia.liczba())

	defer func() {
		polaczenie.Zamknij("koniec obsługi")
		s.polaczenia.usun(polaczenie.Id())
		s.zawiadomOdlaczono(polaczenie)
	}()

	go polaczenie.petlaWysylki()
	polaczenie.petlaOdbioru(s.kontekst, s.rdzenPodlaczony, s.rejestrKomend, &s.pracaRdzenia, s.ustawienia.straznik())
}

// pochodzeniaWlasne to wzorce Origin, którymi przedstawia się własny interfejs produktu w powłoce i przeglądarce.
var pochodzeniaWlasne = []string{
	"tauri://*",
	"https://tauri.localhost",
	"http://localhost:*",
	"http://127.0.0.1:*",
	"https://localhost:*",
	"https://127.0.0.1:*",
	"localhost:*",
	"127.0.0.1:*",
	// Pętla zwrotna ma dwa adresy; ukośnik zdejmuje nawiasowi znaczenie klasy znaków we wzorcu.
	`http://\[::1\]:*`,
	`https://\[::1\]:*`,
	`\[::1\]:*`,
}

// Metoda pochodzeniaDozwolone składa wykaz wzorców Origin dla biblioteki gniazda z wzorców własnych i wskazanych.
func (s *Serwer) pochodzeniaDozwolone() []string {
	wzorce := make([]string, 0, len(pochodzeniaWlasne)+len(s.ustawienia.PochodzeniaDozwolone))
	wzorce = append(wzorce, pochodzeniaWlasne...)
	for _, wzorzec := range s.ustawienia.PochodzeniaDozwolone {
		if wzorzec = strings.TrimSpace(wzorzec); wzorzec != "" {
			wzorce = append(wzorce, wzorzec)
		}
	}
	return wzorce
}

// Funkcja adresZdalny podaje adres urządzenia po drugiej stronie gniazda, niesiony dalej w linii dziennika rdzenia.
func adresZdalny(r *http.Request) string {
	if r == nil || r.RemoteAddr == "" {
		return "nieustalony"
	}
	return r.RemoteAddr
}

// kontoZadania odczytuje konto urządzenia z parametru zapytania albo nagłówka.
// Brak wskazania daje konto domyślne — połączenie nie jest odrzucane.
func kontoZadania(r *http.Request) string {
	if konto := r.URL.Query().Get(ParametrKonta); konto != "" {
		return konto
	}
	if konto := r.Header.Get(NaglowekKonta); konto != "" {
		return konto
	}
	return KontoDomyslne
}

// Metoda nastepnyId nadaje unikalny identyfikator tekstowy kolejnemu nawiązywanemu połączeniu gniazda WebSocket.
func (s *Serwer) nastepnyId() string {
	return fmt.Sprintf("pol-%d", atomic.AddUint64(&s.licznikPolaczen, 1))
}

// zawiadomPrzylaczono powiadamia rdzeń o nowym urządzeniu, jeżeli rdzeń zna
// rozszerzenie ObserwatorPolaczen. Brak rozszerzenia nie jest błędem.
func (s *Serwer) zawiadomPrzylaczono(ujscie Ujscie) {
	if obserwator, zna := s.rdzenPodlaczony().(ObserwatorPolaczen); zna {
		obserwator.Przylaczono(ujscie)
	}
}

// Metoda zawiadomOdlaczono powiadamia rdzeń o rozłączeniu urządzenia korzystającego uprzednio z tego połączenia.
func (s *Serwer) zawiadomOdlaczono(ujscie Ujscie) {
	if obserwator, zna := s.rdzenPodlaczony().(ObserwatorPolaczen); zna {
		obserwator.Odlaczono(ujscie)
	}
}
