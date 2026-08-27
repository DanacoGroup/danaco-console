// Punkt wejścia rdzenia Danaco Console: wyłącznie kompozycja — wczytanie
// konfiguracji, przygotowanie katalogu danych, otwarcie trwałości, złożenie
// rdzenia i praca torów aż do sygnału zatrzymania.
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"danacoconsole/server/cmd/danaco-console/uruchomienie"
	"danacoconsole/server/internal/core"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/store"
	"danacoconsole/server/internal/zdalne"
)

func main() {
	// Dziennik idzie na wyjście diagnostyczne, bo wyjście standardowe niesie kontrakt roli agent.
	dziennik := log.New(os.Stderr, "danaco-console ", log.LstdFlags)

	// Tryb wypisania wykazu zależności stoi przed odczytem konfiguracji, bo wykaz nie potrzebuje bazy.
	if core.ZadanoWykazZaleznosci(os.Args[1:]) {
		if err := core.WypiszWykazZaleznosci(os.Stdout); err != nil {
			dziennik.Fatalf("wykaz zależności: %v", err)
		}
		return
	}
	if core.ZadanoWykazMowy(os.Args[1:]) {
		if err := core.WypiszWykazMowy(os.Stdout); err != nil {
			dziennik.Fatalf("wykaz mowy: %v", err)
		}
		return
	}

	kon, err := konfiguracja.Wczytaj(os.Args[1:], os.Getenv)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		dziennik.Fatalf("konfiguracja: %v", err)
	}
	if err := konfiguracja.PrzygotujKatalogDanych(kon.KatalogDanych); err != nil {
		dziennik.Fatalf("katalog danych: %v", err)
	}

	baza, err := store.Otworz(kon.SciezkaBazy())
	if err != nil {
		dziennik.Fatalf("trwałość: %v", err)
	}
	defer baza.Zamknij()

	// Kontrola spójności zaraz po otwarciu: uszkodzona baza nie może nieść rdzenia, więc zatrzymuje start.
	if err := baza.SprawdzSpojnosc(); err != nil {
		dziennik.Fatalf("spójność bazy: %v", err)
	}

	// Wpięcie toru do hosta zdalnego następuje tutaj, bo plik SQLite ma jedną pulę połączeń w procesie.
	zdalne.Zasil(baza.DB)

	kontekst, zatrzymaj := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer zatrzymaj()

	zmontowany, err := core.Zmontuj(kontekst, core.Montaz{
		Konfiguracja:   kon,
		Baza:           baza,
		Dziennik:       dziennik,
		KatalogKlienta: kon.KatalogKlienta,
		KatalogProfili: kon.KatalogProfili,
	})
	if err != nil {
		dziennik.Fatalf("montaż rdzenia: %v", err)
	}
	defer zmontowany.Zamknij()

	dziennik.Printf("start: %s baza=%s", kon.Opis(), baza.Sciezka)
	otoczenie := uruchomienie.Otoczenie{Wejscie: os.Stdin, Wyjscie: os.Stdout, Dziennik: dziennik}
	if err := uruchomienie.Wedlug(kontekst, kon.Rola, zmontowany.Rdzen, otoczenie); err != nil {
		dziennik.Printf("zatrzymanie: %v", err)
	}
	dziennik.Print("zatrzymanie: rdzeń zamknięty")
}
