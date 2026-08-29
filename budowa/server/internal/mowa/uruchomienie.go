// Odpowiedzialność pliku: przeprowadzenie pomocnika transkrypcji przez
// jedyną dozwoloną drogę startu procesu i zebranie tego, co powiedział.
package mowa

import (
	"context"
	"errors"
	"io"
	"time"

	"danacoconsole/server/internal/session"
)

// Wynik niesie surowy rezultat jednego uruchomienia pomocnika; dwa
// strumienie oddzielnie, a nie sklejone w jeden napis dla Operatora.
type Wynik struct {
	// Wyjscie — bajty ze standardowego wyjścia pomocnika.
	Wyjscie []byte
	// Diagnostyka — tekst z wyjścia diagnostycznego pomocnika.
	Diagnostyka string
}

// Uruchom przeprowadza pomocnika przez port session.Uruchamiacz; limit
// czasu jest obowiązkowy, bo pomocnik bez granicy nie kończy się nigdy.
func Uruchom(ctx context.Context, u session.Uruchamiacz, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	p Pomocnik, argumenty []string, katalog string, limit time.Duration) (Wynik, error) {

	if u == nil {
		return Wynik{}, errors.New("silnik mowy: serwer nie ma uruchamiacza procesów" +
			" — pomocnik transkrypcji nie ma czym wystartować;" +
			" naprawa: podpiąć warstwę kanału (injection) przy składaniu serwera")
	}
	if p.Program == "" || p.Skrypt == "" {
		return Wynik{}, &BrakPomocnika{Powod: errors.New(
			"para interpreter+skrypt niekompletna w chwili uruchomienia")}
	}
	if limit <= 0 {
		return Wynik{}, errors.New("silnik mowy: uruchomienie bez granicy czasu odrzucone" +
			" — transkrypcja bez granicy wisiałaby bez końca;" +
			" naprawa: podać dodatnią granicę czasu przy wywołaniu")
	}

	polecenie := session.Polecenie{
		Program:   p.Program,
		Argumenty: argumenty,
		Katalog:   katalog,
		// Pomocnik potrzebuje środowiska systemu, żeby interpreter odnalazł
		// własne biblioteki.
		DziedziczSrodowisko: true,
	}
	dopuszczone, err := session.SprawdzPolecenie(zasady, obszar, polecenie)
	if err != nil {
		return Wynik{}, err
	}

	uchwyt, err := u.UruchomProces(okno, dopuszczone)
	if err != nil {
		return Wynik{}, errors.New("silnik mowy: nie można uruchomić pomocnika " +
			p.Program + " w " + dopuszczone.Katalog + ": " + err.Error())
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		// Proces już biegnie, a uchwytu drzewa nie ma — zostawienie go byłoby
		// sierotą.
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		return Wynik{}, errors.New("silnik mowy: nie można objąć drzewa procesu pomocnika: " + err.Error())
	}
	defer drzewo.Zwolnij()

	return zbierz(ctx, uchwyt, drzewo, limit)
}

// zbierz prowadzi uruchomiony proces do końca: pompuje oba strumienie,
// czeka z granicą czasu i ubija całe drzewo, gdy granica każe przerwać.
func zbierz(ctx context.Context, uchwyt session.UchwytProcesu, drzewo *session.DrzewoProcesu,
	limit time.Duration) (Wynik, error) {

	wyjscie := make(chan []byte, 1)
	diagnostyka := make(chan []byte, 1)
	go func() { wyjscie <- czytajCalosc(uchwyt.Wyjscie()) }()
	go func() { diagnostyka <- czytajCalosc(uchwyt.Diagnostyka()) }()

	zakonczenie := make(chan error, 1)
	go func() { zakonczenie <- uchwyt.Czekaj() }()

	zegar := time.NewTimer(limit)
	defer zegar.Stop()

	var powod string
	var bladZakonczenia error
	select {
	case bladZakonczenia = <-zakonczenie:
	case <-zegar.C:
		powod = "przekroczona granica czasu " + limit.String()
		_ = drzewo.Ubij()
		bladZakonczenia = <-zakonczenie
	case <-ctx.Done():
		powod = "żądanie przerwane przez serwer"
		_ = drzewo.Ubij()
		bladZakonczenia = <-zakonczenie
	}

	// Odbiór z obu pomp następuje po zakończeniu procesu bezwarunkowo, także
	// przy przerwaniu.
	wynik := Wynik{Wyjscie: <-wyjscie, Diagnostyka: string(<-diagnostyka)}

	if powod != "" {
		// Wynik przerwany jest niepełny, więc nie idzie jako udany; diagnostyka
		// jedzie razem z odmową.
		return wynik, errors.New("silnik mowy: transkrypcja przerwana — " + powod +
			"; naprawa: podnieść granicę czasu albo podać krótsze nagranie")
	}
	if bladZakonczenia != nil {
		return wynik, errors.New("silnik mowy: pomocnik zakończył się niepowodzeniem: " +
			bladZakonczenia.Error())
	}
	return wynik, nil
}

// czytajCalosc zbiera cały strumień procesu; wyjście pomocnika to jedna
// transkrypcja, więc idzie w całości do wyniku, nie przyrostowo.
func czytajCalosc(zrodlo io.Reader) []byte {
	if zrodlo == nil {
		return nil
	}
	bajty, _ := io.ReadAll(zrodlo)
	return bajty
}
