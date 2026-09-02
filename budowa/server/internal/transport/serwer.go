package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// Bez granicy proces przysyłający nagłówki po bajcie trzyma gniazdo bez końca.
	czasNaglowkaZadania = 10 * time.Second
	// Gniazda WebSocket granica nie dotyczy: bezczynne gniazdo rozstrzyga ping.
	czasBezczynnosciHttp = 2 * time.Minute
)

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

func Nowy(ustawienia Ustawienia) *Serwer {
	return &Serwer{
		ustawienia:    ustawienia.zNormalizowane(),
		rejestrKomend: protocol.NowyRejestrKomend(shared.WszystkieKomendy()...),
		polaczenia:    nowyRejestrPolaczen(),
		zakonczone:    make(chan struct{}),
	}
}

// Bez podłączonego rdzenia komendy dostają odpowiedź `*.unknown`, a połączenia żyją.
func (s *Serwer) PodlaczRdzen(rdzen Rdzen) {
	s.zamekRdzenia.Lock()
	s.rdzen = rdzen
	s.zamekRdzenia.Unlock()
}

func (s *Serwer) rdzenPodlaczony() Rdzen {
	s.zamekRdzenia.RLock()
	defer s.zamekRdzenia.RUnlock()
	return s.rdzen
}

// Błąd wraca wyłącznie z zajęcia portu; usterki połączeń nie zatrzymują serwera.
func (s *Serwer) Uruchom(kontekst context.Context) error {
	s.kontekst, s.zakoncz = context.WithCancel(kontekst)

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
	// Czas odczytu całego żądania granicy nie ma: gniazdo czyta przez całą sesję.
	s.serwerHttp = &http.Server{
		Handler:           s.trasy(),
		ReadHeaderTimeout: czasNaglowkaZadania,
		IdleTimeout:       czasBezczynnosciHttp,
	}

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

func (s *Serwer) obsluguj(nasluch net.Listener) error {
	if s.ustawienia.zTLS() {
		return s.serwerHttp.ServeTLS(nasluch, s.ustawienia.CertyfikatTLS, s.ustawienia.KluczTLS)
	}
	return s.serwerHttp.Serve(nasluch)
}

func (s *Serwer) Sluchaj(kontekst context.Context) error {
	if err := s.Uruchom(kontekst); err != nil {
		return err
	}
	<-s.kontekst.Done()
	return s.Zamknij()
}

func (s *Serwer) trasy() http.Handler {
	trasy := http.NewServeMux()
	trasy.HandleFunc(s.ustawienia.SciezkaGniazda, s.nawiaz)
	trasy.Handle("/", uchwytStatyki(s.ustawienia.KatalogKlienta, s.ustawienia.Dziennik))
	return trasy
}

func (s *Serwer) Adres() string {
	if s.nasluch == nil {
		return s.ustawienia.adresNasluchu()
	}
	return s.nasluch.Addr().String()
}

// Najpierw zamyka się wejście, potem czeka na pracę w biegu.
func (s *Serwer) Zamknij() error {
	if s.zakoncz != nil {
		s.zakoncz()
	}
	s.polaczenia.zamknijWszystkie("zatrzymanie serwera")

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
