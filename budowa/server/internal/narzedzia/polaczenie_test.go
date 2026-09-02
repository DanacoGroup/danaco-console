// Pomiary gniazda do rdzenia poza jednym wywołaniem.
package narzedzia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const czasNaPong = 2 * time.Second

const (
	czasNaOdpowiedz    = 2 * time.Second
	wywolanWLawinie    = 40
	rozgloszenPoKazdym = 200
)

func zadanieProbne(identyfikator string) protocol.Koperta {
	return protocol.Koperta{
		Type:      shared.CommandConnectionHello,
		Id:        identyfikator,
		Timestamp: protocol.Teraz(),
	}
}

func odpowiedzNaZadanie(kontekst context.Context, gniazdo *websocket.Conn) error {
	_, bajty, err := gniazdo.Read(kontekst)
	if err != nil {
		return err
	}
	var zadanie protocol.Koperta
	if err := json.Unmarshal(bajty, &zadanie); err != nil {
		return err
	}
	stan := shared.EnvelopeStatus(shared.EnvelopeStatusOk)
	dane, err := protocol.Zakoduj(protocol.Koperta{
		Type:      zadanie.Type,
		Id:        zadanie.Id,
		Status:    &stan,
		Timestamp: protocol.Teraz(),
	})
	if err != nil {
		return err
	}
	return gniazdo.Write(kontekst, websocket.MessageText, dane)
}

// Rdzeń rozgłasza do gniazda narzędzi bez odsiewania po rodzaju, a strumień
// odpowiedzi rozgłasza każdy swój fragment osobno.
func lawinaRozgloszen(kontekst context.Context, gniazdo *websocket.Conn, ile int) error {
	for numer := 0; numer < ile; numer++ {
		stan := shared.EnvelopeStatus(shared.EnvelopeStatusOk)
		dane, err := protocol.Zakoduj(protocol.Koperta{
			Type:      shared.EventStreamChunk,
			Id:        fmt.Sprintf("rozgloszenie-%d", numer),
			Status:    &stan,
			Timestamp: protocol.Teraz(),
		})
		if err != nil {
			return err
		}
		if err := gniazdo.Write(kontekst, websocket.MessageText, dane); err != nil {
			return err
		}
	}
	return nil
}

func TestGniazdoOdpowiadaNaPingPoWywolaniu(t *testing.T) {
	wynikPingu := make(chan error, 1)
	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gniazdo, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer gniazdo.CloseNow()
		if err := odpowiedzNaZadanie(r.Context(), gniazdo); err != nil {
			wynikPingu <- err
			return
		}
		// Pong wraca pętlą czytającą, więc ping stoi obok odczytu.
		go func() {
			kontekst, koniec := context.WithTimeout(r.Context(), czasNaPong)
			defer koniec()
			wynikPingu <- gniazdo.Ping(kontekst)
		}()
		for {
			if _, _, err := gniazdo.Read(r.Context()); err != nil {
				return
			}
		}
	}))
	defer serwer.Close()

	polaczenie := Polacz(serwer.URL, "okno-probne", ZasiegOkna, "poswiadczenie-probne")
	defer polaczenie.Zamknij()

	if _, err := polaczenie.Wykonaj(context.Background(), zadanieProbne("narzedzia-1")); err != nil {
		t.Fatalf("wywołanie pierwsze nie doszło do skutku: %v", err)
	}
	select {
	case err := <-wynikPingu:
		if err != nil {
			t.Fatalf("gniazdo bezczynne nie odpowiedziało na ping: %v", err)
		}
	case <-time.After(czasNaPong + time.Second):
		t.Fatal("serwer próbny nie zameldował wyniku pingu")
	}
}

func TestWykonajPonawiaZadanieNaGniezdzieZastanym(t *testing.T) {
	var polaczen atomic.Int64
	zamkniete := make(chan struct{})
	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gniazdo, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer gniazdo.CloseNow()
		if polaczen.Add(1) == 1 {
			if err := odpowiedzNaZadanie(r.Context(), gniazdo); err != nil {
				return
			}
			gniazdo.CloseNow()
			close(zamkniete)
			return
		}
		for {
			if err := odpowiedzNaZadanie(r.Context(), gniazdo); err != nil {
				return
			}
		}
	}))
	defer serwer.Close()

	polaczenie := Polacz(serwer.URL, "okno-probne", ZasiegOkna, "poswiadczenie-probne")
	defer polaczenie.Zamknij()

	if _, err := polaczenie.Wykonaj(context.Background(), zadanieProbne("narzedzia-1")); err != nil {
		t.Fatalf("wywołanie pierwsze nie doszło do skutku: %v", err)
	}
	<-zamkniete

	odpowiedz, err := polaczenie.Wykonaj(context.Background(), zadanieProbne("narzedzia-2"))
	if err != nil {
		t.Fatalf("wywołanie na gnieździe zastanym przepadło zamiast zostać ponowione: %v", err)
	}
	if odpowiedz.Id != "narzedzia-2" {
		t.Errorf("odpowiedź niesie identyfikator %q, a żądanie miało narzedzia-2", odpowiedz.Id)
	}
	if liczba := polaczen.Load(); liczba != 2 {
		t.Errorf("serwer przyjął %d połączeń — ponowienie miało zestawić drugie", liczba)
	}
}

func TestKontekstZerwanyNieDajePonowienia(t *testing.T) {
	var polaczen atomic.Int64
	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gniazdo, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer gniazdo.CloseNow()
		polaczen.Add(1)
		for {
			if err := odpowiedzNaZadanie(r.Context(), gniazdo); err != nil {
				return
			}
		}
	}))
	defer serwer.Close()

	polaczenie := Polacz(serwer.URL, "okno-probne", ZasiegOkna, "poswiadczenie-probne")
	defer polaczenie.Zamknij()

	if _, err := polaczenie.Wykonaj(context.Background(), zadanieProbne("narzedzia-1")); err != nil {
		t.Fatalf("wywołanie pierwsze nie doszło do skutku: %v", err)
	}
	kontekst, zerwij := context.WithCancel(context.Background())
	zerwij()
	if _, err := polaczenie.Wykonaj(kontekst, zadanieProbne("narzedzia-2")); err == nil {
		t.Fatal("wywołanie z kontekstem zerwanym doszło do skutku")
	}
	if liczba := polaczen.Load(); liczba != 1 {
		t.Errorf("serwer przyjął %d połączeń — kontekst zerwany nie miał zestawiać gniazda od nowa", liczba)
	}
}

func TestOdpowiedzPrzezywaLawineRozgloszen(t *testing.T) {
	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gniazdo, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer gniazdo.CloseNow()
		for {
			if err := odpowiedzNaZadanie(r.Context(), gniazdo); err != nil {
				return
			}
			if err := lawinaRozgloszen(r.Context(), gniazdo, rozgloszenPoKazdym); err != nil {
				return
			}
		}
	}))
	defer serwer.Close()

	polaczenie := Polacz(serwer.URL, "okno-probne", ZasiegOkna, "poswiadczenie-probne")
	defer polaczenie.Zamknij()

	for numer := 1; numer <= wywolanWLawinie; numer++ {
		identyfikator := fmt.Sprintf("narzedzia-%d", numer)
		kontekst, koniec := context.WithTimeout(context.Background(), czasNaOdpowiedz)
		odpowiedz, err := polaczenie.Wykonaj(kontekst, zadanieProbne(identyfikator))
		koniec()
		if err != nil {
			t.Fatalf("wywołanie %d: odpowiedź przepadła w lawinie rozgłoszeń: %v", numer, err)
		}
		if odpowiedz.Id != identyfikator {
			t.Fatalf("wywołanie %d: odpowiedź niesie identyfikator %q", numer, odpowiedz.Id)
		}
	}
}
