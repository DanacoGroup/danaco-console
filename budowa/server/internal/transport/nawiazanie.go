package transport

import (
	"crypto/subtle"
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

const (
	// limitPolaczen jest górną granicą gniazd stojących naraz. Każde gniazdo to
	// dwie gorutyny, bufor wyjściowy i ramka do 16 MiB w odczycie, więc rejestr
	// bez granicy jest drogą wyczerpania pamięci maszyny z jednego procesu.
	limitPolaczen = 512
	// limitPolaczenNiezwiazanych ogranicza gniazda, które nie przeszły przez
	// bramkę. Obowiązuje wyłącznie tam, gdzie bramka stoi: bez wymogu logowania
	// żadne gniazdo nie jest związane i granica odcięłaby pracę własną.
	limitPolaczenNiezwiazanych = 64
)

// Metoda nawiaz przyjmuje żądanie uaktualnienia do WebSocket i prowadzi całe życie połączenia od rejestracji do wykreślenia.
func (s *Serwer) nawiaz(w http.ResponseWriter, r *http.Request) {
	// Sekret i granica rejestru rozstrzygają się przed uaktualnieniem gniazda:
	// po uaktualnieniu odmowa nie ma już postaci kodu HTTP.
	if !s.sekretZgodny(r) {
		s.ustawienia.Dziennik.Printf("transport: nawiązanie z %s odrzucone — sekret nawiązania niezgodny", adresZdalny(r))
		http.Error(w, "sekret nawiązania niezgodny", http.StatusForbidden)
		return
	}
	if powod, pelno := s.brakMiejscaWRejestrze(); pelno {
		s.ustawienia.Dziennik.Printf("transport: nawiązanie z %s odrzucone — %s", adresZdalny(r), powod)
		http.Error(w, powod, http.StatusServiceUnavailable)
		return
	}
	// Poświadczenie serwera narzędzi rozstrzyga się tu, przed uaktualnieniem:
	// rodzaj i zasięg z zapytania wchodzą do tożsamości dopiero po sprawdzeniu.
	poswiadczone, odmowa := s.poswiadczenieNarzedziZgodne(r)
	if odmowa != "" {
		s.ustawienia.Dziennik.Printf("transport: nawiązanie z %s odrzucone — %s", adresZdalny(r), odmowa)
		http.Error(w, odmowa, http.StatusForbidden)
		return
	}
	gniazdo, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: s.pochodzeniaDozwolone(),
	})
	if err != nil {
		s.ustawienia.Dziennik.Printf("transport: nawiązanie odrzucone: %v", err)
		return
	}
	gniazdo.SetReadLimit(LimitOdczytu)

	id := s.nastepnyId()
	// Konto nadaje rdzeń po przejściu bramki; przy nawiązaniu każde gniazdo stoi na koncie domyślnym.
	polaczenie := nowePolaczenie(s.kontekst, id, KontoDomyslne, tozsamoscZadania(id, r, poswiadczone), gniazdo, s.ustawienia.PojemnoscKolejki, s.ustawienia.Dziennik)
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
	go polaczenie.petlaPingu()
	polaczenie.petlaOdbioru(s.kontekst, s.rdzenPodlaczony, s.rejestrKomend, &s.pracaRdzenia, s.ustawienia.dopuszczenie())
}

// sekretZgodny sprawdza sekret nawiązania podawany przez powłokę. Sekret
// niewskazany znaczy sprawdzenie niepodniesione — nie ma z czym porównywać.
// Porównanie idzie czasem stałym, bo różnica czasu odpowiedzi wystarcza do
// odgadnięcia sekretu bajt po bajcie.
func (s *Serwer) sekretZgodny(r *http.Request) bool {
	oczekiwany := s.ustawienia.SekretNawiazania
	if oczekiwany == "" {
		return true
	}
	podany := r.URL.Query().Get(ParametrSekretu)
	if podany == "" {
		podany = r.Header.Get(NaglowekSekretu)
	}
	return subtle.ConstantTimeCompare([]byte(podany), []byte(oczekiwany)) == 1
}

// poswiadczenieNarzedziZgodne sprawdza poświadczenie gniazda, które przedstawia
// się rodzajem albo poświadczeniem. Klient bez jednego i drugiego nie jest
// serwerem narzędzi i przechodzi bez sprawdzenia, ale i bez rodzaju. Porównanie
// idzie czasem stałym z tego samego powodu co przy sekrecie nawiązania.
func (s *Serwer) poswiadczenieNarzedziZgodne(r *http.Request) (sprawdzone bool, odmowa string) {
	zapytanie := r.URL.Query()
	podane := zapytanie.Get(ParametrPoswiadczenia)
	if podane == "" && strings.TrimSpace(zapytanie.Get(ParametrRodzaju)) == "" {
		return false, ""
	}
	oczekiwane := s.ustawienia.PoswiadczenieNarzedzi
	if oczekiwane == "" {
		return false, "poświadczenie serwera narzędzi niewydane"
	}
	if subtle.ConstantTimeCompare([]byte(podane), []byte(oczekiwane)) != 1 {
		return false, "poświadczenie serwera narzędzi niezgodne"
	}
	return true, ""
}

// brakMiejscaWRejestrze nazywa granicę, o którą opiera się nawiązanie, albo
// milczy. Granica gniazd niezwiązanych liczy się tylko przy bramce stojącej,
// bo bez wymogu logowania niezwiązane są wszystkie.
func (s *Serwer) brakMiejscaWRejestrze() (string, bool) {
	stojace := s.polaczenia.wszystkie()
	if len(stojace) >= limitPolaczen {
		return fmt.Sprintf("rejestr połączeń pełny (%d z %d)", len(stojace), limitPolaczen), true
	}
	if !s.ustawienia.dopuszczenie().wymagana {
		return "", false
	}
	stan, zna := s.rdzenPodlaczony().(StanBramki)
	if !zna {
		return "", false
	}
	niezwiazane := 0
	for _, p := range stojace {
		if !stan.PolaczenieZwiazane(p.Id()) {
			niezwiazane++
		}
	}
	if niezwiazane >= limitPolaczenNiezwiazanych {
		return fmt.Sprintf("gniazda przed bramką zajęły granicę (%d z %d)",
			niezwiazane, limitPolaczenNiezwiazanych), true
	}
	return "", false
}

// pochodzeniaWlasne to wzorce Origin, którymi przedstawia się własny interfejs
// produktu w powłoce i przeglądarce. Biblioteka gniazda dopasowuje wyłącznie
// GOSPODARZA nagłówka Origin — schemat zdejmuje przed dopasowaniem — więc
// wzorzec ze schematem nie zgadza się nigdy z niczym.
// Wzorców z portem dowolnym tu nie ma: pod `127.0.0.1:*` mieści się każdy
// nasłuch tej maszyny, w tym podgląd warstwy Apps, którego stronę pisze model —
// skrypt takiej strony otwierałby gniazdo z pełnym wykazem komend. Interfejs
// podawany przez sam rdzeń przechodzi bez wzorca, bo gospodarz nagłówka Origin
// jest wtedy gospodarzem żądania; pochodzenia pracy deweloperskiej (serwer Vite)
// wskazuje wykaz z nastaw.
var pochodzeniaWlasne = []string{
	// Powłoka desktopowa podaje stronę z pakietu: Windows przedstawia ją
	// gospodarzem `tauri.localhost`, pozostałe platformy — `localhost` bez portu.
	"tauri.localhost",
	"localhost",
	"127.0.0.1",
	// Pętla zwrotna ma dwa adresy; ukośnik zdejmuje nawiasowi znaczenie klasy znaków we wzorcu.
	`\[::1\]`,
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
