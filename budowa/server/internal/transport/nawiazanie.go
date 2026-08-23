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

// nawiaz przyjmuje żądanie uaktualnienia do WebSocket i prowadzi całe życie
// połączenia: rejestrację, obie pętle i wykreślenie po rozłączeniu.
//
// Pochodzenie jest sprawdzane i nie jest to bramka kontrolna: wykaz pochodzeń
// nie pyta, kim jest wołający i czego mu wolno — odcina wyłącznie stronę trzecią,
// która namówiła przeglądarkę Operatora, żeby otworzyła gniazdo do jego rdzenia.
// Operator nie widzi tego nigdy; widzi to wyłącznie cudza strona.
//
// Klientem bywa webview powłoki, który przedstawia się rozmaitym Origin —
// pochodzenia własne obejmują `tauri://`, `http://localhost` i `http://127.0.0.1`
// z dowolnym portem, więc powłoka wchodzi bez wskazywania czegokolwiek.
// Wystawienie pod inną domenę dopisuje ją polem PochodzeniaDozwolone.
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
	// Tożsamość jest w linii dziennika, bo od niej zależy sprawca zdarzeń
	// (`core/sprawca.go`). Gdy Operator zobaczy w interfejsie rękę nie tę, co
	// trzeba, pierwszym miejscem do sprawdzenia jest to, czym gniazdo się
	// przedstawiło — a to widać tylko tutaj.
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

// pochodzeniaWlasne to wzorce Origin, którymi przedstawia się własny interfejs
// produktu. Powłoka Tauri podaje `tauri://localhost` (Windows: `https://tauri.localhost`),
// interfejs otwarty w przeglądarce — `http://127.0.0.1:<port>` albo
// `http://localhost:<port>`, a port bywa dowolny (nasłuch potrafi wziąć port
// wskazany przez system). Stąd gwiazdka w porcie, a nie w całym wzorcu.
var pochodzeniaWlasne = []string{
	"tauri://*",
	"https://tauri.localhost",
	"http://localhost:*",
	"http://127.0.0.1:*",
	"https://localhost:*",
	"https://127.0.0.1:*",
	"localhost:*",
	"127.0.0.1:*",
	// Pętla zwrotna ma dwa adresy, nie jeden. Przeglądarka na maszynie
	// z pierwszeństwem IPv6 rozwiązuje `localhost` na `::1` i podaje wtedy
	// pochodzenie `http://[::1]:<port>` — ta sama pętla zwrotna, na której rdzeń
	// nasłuchuje, a bez tych wzorców nawiązanie kończyłoby się odmową 403.
	//
	// Ukośniki odwrotne są tu konieczne, nie ozdobne. Dopasowanie idzie przez
	// `path.Match`, gdzie nawias kwadratowy otwiera klasę znaków; wzorzec
	// z nawiasem gołym jest wzorcem wadliwym, a biblioteka gniazda przerywa
	// wtedy przegląd wykazu błędem — czyli jeden zły wzorzec potrafi odciąć
	// pochodzenia sprawdzane po nim. Ukośnik zdejmuje znakowi znaczenie i nawias
	// wraca do bycia nawiasem.
	`http://\[::1\]:*`,
	`https://\[::1\]:*`,
	`\[::1\]:*`,
}

// pochodzeniaDozwolone składa wykaz wzorców dla biblioteki gniazda.
//
// Wykaz wskazany dopisuje się do własnych, nie zastępuje ich. Wystawienie pod
// domenę nie jest powodem, żeby produkt przestał wpuszczać własną powłokę —
// a taki właśnie byłby skutek zastąpienia.
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

// adresZdalny podaje adres urządzenia po drugiej stronie gniazda.
//
// Po co to w linii dziennika. Telemetria połączeń jedzie dziennikiem rdzenia
// i przez rozgałęzienie trafia do `diagnostyka_wpis` — trwale i z czytelnikiem
// (`diagnostics.log.query`). Adres urządzenia jest jednym z faktów, które ta
// linia ma nieść.
//
// Adres bierzemy z pola żądania HTTP, nie z nagłówków przekazywania (X-Forwarded-For
// i pokrewnych): te podaje strona trzecia i można je napisać dowolnie, a dziennik
// ma nieść fakt gniazda, nie deklarację nadawcy. Brak adresu daje wpis „nieustalony"
// zamiast pustego miejsca — czytający ma widzieć różnicę między „nie wiadomo"
// a „pominięto".
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

// nastepnyId nadaje identyfikator kolejnemu połączeniu.
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

// zawiadomOdlaczono powiadamia rdzeń o rozłączeniu urządzenia.
func (s *Serwer) zawiadomOdlaczono(ujscie Ujscie) {
	if obserwator, zna := s.rdzenPodlaczony().(ObserwatorPolaczen); zna {
		obserwator.Odlaczono(ujscie)
	}
}
