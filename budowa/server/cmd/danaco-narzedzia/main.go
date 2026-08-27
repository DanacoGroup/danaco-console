// Punkt wejścia serwera narzędzi modelu: wyłącznie kompozycja — odczyt
// przełączników, gniazdo do rdzenia, rozdzielnia w zasięgu okna, protokół MCP
// po strumieniach procesu.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"danacoconsole/server/cmd/danaco-narzedzia/stdio"
	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/narzedzia"
)

// Rozdzielnia jest katalogiem narzędzi warstwy protokołu. Sprawdzenie stoi przy
// kompilacji, żeby rozjazd sygnatur zatrzymał budowanie, a nie rozmowę.
var _ stdio.Katalog = (*narzedzia.Rozdzielnia)(nil)

func main() {
	dziennik := log.New(os.Stderr, narzedzia.NazwaBinarium+" ", log.LstdFlags)

	okno := flag.String(narzedzia.PrzelacznikOkna, "",
		"okno rozmowy, w którego zasięgu pracuje serwer narzędzi")
	adres := flag.String(narzedzia.PrzelacznikRdzenia, narzedzia.AdresRdzenia(),
		"adres gniazda WebSocket rdzenia")
	// Zasięg jest rolą okna, nie prośbą modelu; przełącznik czyta się raz, przy uruchomieniu.
	zasieg := flag.String(narzedzia.PrzelacznikZasiegu, string(narzedzia.ZasiegOkna),
		"rola okna rozstrzygająca zasięg narzędzi: okno, klawiatura albo ekspert")
	// Kod eksperta nałożonego na okno; przełącznik osobny od zasięgu, bo niesie wartość, nie nazwę roli.
	ekspert := flag.String(narzedzia.PrzelacznikEksperta, "",
		"kod eksperta nałożonego na okno; zawęża wykaz narzędzi do jego skillIds i connectorIds")
	// Doraźne dołożenia sesji; przełącznik nieznany kończy proces, więc odczyt musi stać przed dołożeniem.
	dolozenia := flag.String(injection.PrzelacznikDolozen, "",
		"narzędzia dołożone doraźnie w sesji, po przecinku; dokładają się do wykazu eksperta")
	flag.Parse()

	rola := rolaOkna(dziennik, *zasieg, *ekspert)

	// Brak okna nie zatrzymuje serwera: narzędzia działają dalej, bez uzupełniania zasięgu.
	if *okno == "" {
		dziennik.Printf("uruchomienie bez --%s: zasięg okna nie będzie uzupełniany", narzedzia.PrzelacznikOkna)
	}

	kontekst, zatrzymaj := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer zatrzymaj()

	polaczenie := narzedzia.Polacz(*adres, *okno, rola)
	defer polaczenie.Zamknij()

	pomiar := narzedzia.Zmierz(narzedzia.WykazZasiegu(rola))
	dziennik.Printf("start: okno=%q zasięg=%s ekspert=%q rdzeń=%s narzędzi=%d grup=%d bajtów=%d",
		*okno, rola, *ekspert, *adres, pomiar.Pozycji, len(pomiar.Grupy), pomiar.Bajtow)
	if rola == narzedzia.ZasiegEksperta && *ekspert != "" {
		// Liczby wyżej dotyczą zestawu niezawężonego; cenę eksperta melduje dobór przy pierwszym wykazie.
		dziennik.Printf("dobór eksperta %q: liczby wyżej dotyczą zestawu PRZED zawężeniem —"+
			" definicja czytana z rdzenia przy pierwszym tools/list", *ekspert)
	}

	dolozone := narzedzia.RozbijDolozenia(*dolozenia)
	if len(dolozone) > 0 && rola != narzedzia.ZasiegEksperta {
		// Dołożenie ma co dołożyć wyłącznie do wykazu zawężonego; okno bez eksperta ma już pełny wykaz.
		dziennik.Printf("dołożenia sesji %v bez zawężenia (zasięg %s): wykaz nie jest zawężony,"+
			" więc dołożone narzędzia już w nim stoją — dołożenie nie zmienia tury", dolozone, rola)
	}

	rozdzielnia := narzedzia.NowaRozdzielnia(polaczenie, *okno, rola).
		ZEkspertem(*ekspert, dziennik).
		ZDolozeniamiSesji(dolozone)
	if err := stdio.Obsluguj(kontekst, rozdzielnia, os.Stdin, os.Stdout); err != nil {
		dziennik.Printf("zatrzymanie: %v", err)
	}
	dziennik.Print("zatrzymanie: wejście wyczerpane")
}

// rolaOkna rozstrzyga zasięg z dwóch przełączników naraz i melduje ich
// rozjazd; żadna z dwóch stron rozjazdu nie przechodzi w milczeniu.
func rolaOkna(dziennik *log.Logger, zasieg, ekspert string) narzedzia.Zasieg {
	rola := narzedzia.RozpoznajZasieg(zasieg)
	if ekspert != "" && rola != narzedzia.ZasiegEksperta {
		dziennik.Printf("przełącznik --%s niesie kod %q przy zasięgu %s — zasięg podniesiony do %s",
			narzedzia.PrzelacznikEksperta, ekspert, rola, narzedzia.ZasiegEksperta)
		return narzedzia.ZasiegEksperta
	}
	if ekspert == "" && rola == narzedzia.ZasiegEksperta {
		dziennik.Printf("zasięg %s bez przełącznika --%s: nie ma o kogo zapytać rdzenia,"+
			" wykaz NIE ZOSTANIE zawężony — zasięg schodzi do %s",
			narzedzia.ZasiegEksperta, narzedzia.PrzelacznikEksperta, narzedzia.ZasiegOkna)
		return narzedzia.ZasiegOkna
	}
	return rola
}
