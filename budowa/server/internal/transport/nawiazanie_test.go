package transport

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

const (
	sekretSprawdzianu        = "sekret-powloki"
	poswiadczenieSprawdzianu = "poswiadczenie-rdzenia"
)

// Poświadczenie domyślne jest losowane raz na proces, więc obie wartości stoją tu wprost.
func podniesZSekretem(t *testing.T) string {
	t.Helper()

	ustawienia := Domyslne()
	ustawienia.Adres = "127.0.0.1"
	ustawienia.PortDowolny = true
	ustawienia.Dziennik = log.New(io.Discard, "", 0)
	ustawienia.SekretNawiazania = sekretSprawdzianu
	ustawienia.PoswiadczenieNarzedzi = poswiadczenieSprawdzianu

	serwer := Nowy(ustawienia)
	serwer.PodlaczRdzen(&rdzenAtrapa{})

	zycie, zakoncz := context.WithCancel(context.Background())
	t.Cleanup(zakoncz)
	if err := serwer.Uruchom(zycie); err != nil {
		t.Fatalf("nie można uruchomić nasłuchu: %v", err)
	}
	t.Cleanup(func() { _ = serwer.Zamknij() })

	return "ws://" + serwer.Adres() + SciezkaGniazdaDomyslna
}

func adresNarzedzi(adres, poswiadczenie string) string {
	zapytanie := url.Values{}
	zapytanie.Set(ParametrRodzaju, RodzajNarzedzi)
	zapytanie.Set(ParametrPoswiadczenia, poswiadczenie)
	return adres + "?" + zapytanie.Encode()
}

// Kod 403 nie rozróżnia sekretu nawiązania od poświadczenia narzędzi.
func sprawdzOdmowe(t *testing.T, odpowiedz *http.Response, powod string) {
	t.Helper()
	if odpowiedz == nil {
		t.Fatalf("odmowa bez odpowiedzi HTTP; oczekiwany powód %q", powod)
	}
	defer func() { _ = odpowiedz.Body.Close() }()
	if odpowiedz.StatusCode != http.StatusForbidden {
		t.Fatalf("odmowa kodem %d zamiast %d", odpowiedz.StatusCode, http.StatusForbidden)
	}
	tresc, err := io.ReadAll(odpowiedz.Body)
	if err != nil {
		t.Fatalf("nieczytelna treść odmowy: %v", err)
	}
	if !strings.Contains(string(tresc), powod) {
		t.Fatalf("odmowa %q zamiast %q", strings.TrimSpace(string(tresc)), powod)
	}
}

func TestPoswiadczenieNarzedziWchodziMimoSekretuNawiazania(t *testing.T) {
	adres := podniesZSekretem(t)

	ctx, przerwij := context.WithTimeout(context.Background(), 5*time.Second)
	defer przerwij()
	gniazdo, _, err := websocket.Dial(ctx, adresNarzedzi(adres, poswiadczenieSprawdzianu), nil)
	if err != nil {
		t.Fatalf("gniazdo serwera narzędzi odrzucone mimo zgodnego poświadczenia: %v", err)
	}
	_ = gniazdo.CloseNow()
}

func TestNiezgodnePoswiadczenieOdrzucaMimoZgodnegoSekretu(t *testing.T) {
	adres := podniesZSekretem(t)

	zapytanie := url.Values{}
	zapytanie.Set(ParametrRodzaju, RodzajNarzedzi)
	zapytanie.Set(ParametrPoswiadczenia, "poswiadczenie-obce")
	zapytanie.Set(ParametrSekretu, sekretSprawdzianu)

	ctx, przerwij := context.WithTimeout(context.Background(), 5*time.Second)
	defer przerwij()
	gniazdo, odpowiedz, err := websocket.Dial(ctx, adres+"?"+zapytanie.Encode(), nil)
	if err == nil {
		_ = gniazdo.CloseNow()
		t.Fatal("gniazdo z niezgodnym poświadczeniem zostało nawiązane")
	}
	sprawdzOdmowe(t, odpowiedz, "poświadczenie serwera narzędzi niezgodne")
}

func TestGniazdoBezPoswiadczeniaWymagaSekretuNawiazania(t *testing.T) {
	adres := podniesZSekretem(t)

	ctx, przerwij := context.WithTimeout(context.Background(), 5*time.Second)
	defer przerwij()
	gniazdo, odpowiedz, err := websocket.Dial(ctx, adres, &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://localhost"}},
	})
	if err == nil {
		_ = gniazdo.CloseNow()
		t.Fatal("gniazdo bez sekretu nawiązania zostało nawiązane")
	}
	sprawdzOdmowe(t, odpowiedz, "sekret nawiązania niezgodny")
}
