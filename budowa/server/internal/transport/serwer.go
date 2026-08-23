package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Serwer jest nasłuchem rdzenia: przyjmuje połączenia WebSocket, serwuje pliki
// klienta i realizuje interfejs Rozglosnik.
//
// Rdzeń podłącza się po utworzeniu serwera (PodlaczRdzen), a nie przy budowie.
// Tak zerwana jest zależność cykliczna: rdzeń potrzebuje rozgłośni transportu,
// transport potrzebuje rdzenia do obsługi komend, a żaden nie importuje
// pakietu drugiego.
type Serwer struct {
	ustawienia    Ustawienia
	rejestrKomend *protocol.RejestrKomend
	polaczenia    *rejestrPolaczen

	zamekRdzenia    sync.RWMutex
	rdzen           Rdzen
	licznikPolaczen uint64
	pracaRdzenia    sync.WaitGroup

	serwerHttp *http.Server
	nasluch    net.Listener
	kontekst   context.Context
	zakoncz    context.CancelFunc
	zakonczone chan struct{}
}

// Nowy buduje serwer transportu. Zbiór znanych komend pochodzi w całości
// z kontraktu — transport nie zna ani jednego literału nazwy.
func Nowy(ustawienia Ustawienia) *Serwer {
	return &Serwer{
		ustawienia:    ustawienia.zNormalizowane(),
		rejestrKomend: protocol.NowyRejestrKomend(shared.WszystkieKomendy()...),
		polaczenia:    nowyRejestrPolaczen(),
		zakonczone:    make(chan struct{}),
	}
}

// PodlaczRdzen wskazuje realizację obsługi komend. Wywołanie przed startem
// nasłuchu i w jego trakcie jest równoważne: dopóki rdzeń nie jest podłączony,
// komendy dostają odpowiedź `*.unknown`, a połączenia żyją.
func (s *Serwer) PodlaczRdzen(rdzen Rdzen) {
	s.zamekRdzenia.Lock()
	s.rdzen = rdzen
	s.zamekRdzenia.Unlock()
}

// rdzenPodlaczony odczytuje bieżącą realizację obsługi komend.
func (s *Serwer) rdzenPodlaczony() Rdzen {
	s.zamekRdzenia.RLock()
	defer s.zamekRdzenia.RUnlock()
	return s.rdzen
}

// Uruchom otwiera nasłuch i przechodzi do obsługi w tle. Błąd wraca wyłącznie
// z zajęcia portu; dalsze usterki pojedynczych połączeń nie zatrzymują serwera.
func (s *Serwer) Uruchom(kontekst context.Context) error {
	s.kontekst, s.zakoncz = context.WithCancel(kontekst)

	// Niekompletna para TLS zatrzymuje start celowo i jest jedynym miejscem
	// w tym pakiecie, gdzie coś się nie uruchamia. Reguła fail-open
	// mówi, że brak nastawy nie wstrzymuje pracy — tu nastawa nie jest brakiem,
	// tylko połową wskazania. Praca otwartym tekstem przy wskazanym certyfikacie
	// byłaby cichym zejściem poniżej tego, o co poprosił Operator.
	if err := s.ustawienia.sprawdzTLS(); err != nil {
		s.zakoncz()
		return err
	}

	nasluch, err := net.Listen("tcp", s.ustawienia.adresNasluchu())
	if err != nil {
		s.zakoncz()
		return fmt.Errorf("transport: nasłuch %s: %w", s.ustawienia.adresNasluchu(), err)
	}
	s.nasluch = nasluch
	s.serwerHttp = &http.Server{Handler: s.trasy()}

	go func() {
		defer close(s.zakonczone)
		if err := s.obsluguj(nasluch); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.ustawienia.Dziennik.Printf("transport: obsługa nasłuchu zakończona: %v", err)
		}
	}()

	s.ustawienia.Dziennik.Printf("transport: nasłuch %s gniazdo=%s klient=%s warstwa=%s",
		s.Adres(), s.ustawienia.SciezkaGniazda, opisKatalogu(s.ustawienia.KatalogKlienta),
		s.ustawienia.opisWarstwy())
	return nil
}

// obsluguj prowadzi nasłuch warstwą właściwą dla nastaw: TLS, gdy para plików
// jest wskazana, otwartym tekstem w przeciwnym razie. Rozstrzygnięcie stoi tutaj,
// a nie w Uruchom, żeby po starcie została jedna droga wyjścia i jeden zapis
// dziennika.
func (s *Serwer) obsluguj(nasluch net.Listener) error {
	if s.ustawienia.zTLS() {
		return s.serwerHttp.ServeTLS(nasluch, s.ustawienia.CertyfikatTLS, s.ustawienia.KluczTLS)
	}
	return s.serwerHttp.Serve(nasluch)
}

// Sluchaj otwiera nasłuch i oddaje sterowanie dopiero po zamknięciu kontekstu.
// Jest postacią blokującą Uruchom — dla warstwy składającej, która prowadzi
// nasłuch jako jedno zadanie o czasie życia procesu.
func (s *Serwer) Sluchaj(kontekst context.Context) error {
	if err := s.Uruchom(kontekst); err != nil {
		return err
	}
	<-s.kontekst.Done()
	return s.Zamknij()
}

// trasy składa mapę ścieżek: kanał WebSocket oraz pliki klienta pod resztą.
func (s *Serwer) trasy() http.Handler {
	trasy := http.NewServeMux()
	trasy.HandleFunc(s.ustawienia.SciezkaGniazda, s.nawiaz)
	trasy.Handle("/", uchwytStatyki(s.ustawienia.KatalogKlienta, s.ustawienia.Dziennik))
	return trasy
}

// Adres zwraca rzeczywisty adres nasłuchu — po starcie z portem dowolnym jest
// jedynym sposobem poznania numeru portu.
func (s *Serwer) Adres() string {
	if s.nasluch == nil {
		return s.ustawienia.adresNasluchu()
	}
	return s.nasluch.Addr().String()
}

// Zamknij kończy nasłuch, rozłącza urządzenia i czeka na dokończenie obsługi
// żądań już przyjętych. Kolejność jest istotna: najpierw zamyka się wejście,
// potem czeka na pracę w biegu — inaczej wynik pracy nie miałby dokąd wrócić.
func (s *Serwer) Zamknij() error {
	if s.zakoncz != nil {
		s.zakoncz()
	}
	s.polaczenia.zamknijWszystkie("zatrzymanie rdzenia")

	var blad error
	if s.serwerHttp != nil {
		kontekst, koniec := context.WithTimeout(context.Background(), s.ustawienia.CzasZamkniecia)
		defer koniec()
		if err := s.serwerHttp.Shutdown(kontekst); err != nil {
			blad = fmt.Errorf("transport: zamknięcie nasłuchu: %w", err)
		}
		<-s.zakonczone
	}
	s.pracaRdzenia.Wait()
	return blad
}
