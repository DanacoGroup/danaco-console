// Pakiet zewnetrzne jest jedną drogą, którą rdzeń woła binaria arsenału; oddaje bajty i kod wyjścia, nie rozbiera wyjścia programu.
package zewnetrzne

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
)

// Nazwa jest nazwą dla Operatora i wchodzi do treści odmowy.
type Narzedzie struct {
	Nazwa   string
	Program string
	Pakiet  string
}

type Wynik struct {
	Wyjscie     []byte
	Diagnostyka string
}

type BrakNarzedzia struct {
	Narzedzie Narzedzie
}

func (b *BrakNarzedzia) Error() string {
	zdanie := "arsenał: nie ma na tej maszynie programu " + b.Narzedzie.Nazwa +
		" (" + b.Narzedzie.Program + "), a serwer tej czynności nie wykona bez niego"
	/* Odmowa nazywa brakujący pakiet, ale nie wydaje polecenia instalacji:
	   arsenał stoi po stronie serwera i jego niekompletność jest usterką
	   wdrożenia, nie zadaniem dla czytającego tę odmowę. */
	if b.Narzedzie.Pakiet != "" {
		zdanie += "; instalacja serwera jest niepełna — brakuje pakietu " + b.Narzedzie.Pakiet
	}
	return zdanie
}

func Stoi(n Narzedzie) bool {
	_, jest := Odnajdz(n)
	return jest
}

// Środowisko nie jest dziedziczone: binaria arsenału są samodzielne.
func Wolaj(ctx context.Context, u session.Uruchamiacz, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	n Narzedzie, argumenty []string, katalog string, limit time.Duration) (Wynik, error) {

	if u == nil {
		return Wynik{}, errors.New("arsenał: serwer nie ma uruchamiacza procesów — " +
			"narzędzia zewnętrzne nie mają czym wystartować; " +
			"naprawa: podpiąć warstwę kanału (injection) przy składaniu serwera")
	}
	sciezka, jest := Odnajdz(n)
	if !jest {
		return Wynik{}, &BrakNarzedzia{Narzedzie: n}
	}
	if limit <= 0 {
		return Wynik{}, errors.New("arsenał: uruchomienie " + n.Nazwa +
			" bez granicy czasu odrzucone — przetwarzanie bez granicy wisiałoby bez końca; " +
			"naprawa: podać dodatnią granicę czasu przy wywołaniu")
	}

	polecenie := session.Polecenie{
		Program:   sciezka,
		Argumenty: argumenty,
		Katalog:   katalog,
	}
	dopuszczone, err := session.SprawdzPolecenie(zasady, obszar, polecenie)
	if err != nil {
		return Wynik{}, err
	}

	uchwyt, err := u.UruchomProces(ctx, okno, dopuszczone)
	if err != nil {
		return Wynik{}, errors.New("arsenał: nie można uruchomić " + n.Nazwa +
			" (" + n.Program + ") w " + dopuszczone.Katalog + ": " + err.Error())
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		return Wynik{}, errors.New("arsenał: nie można objąć drzewa procesu " +
			n.Nazwa + ": " + err.Error())
	}
	defer drzewo.Zwolnij()

	return zbierz(ctx, n, uchwyt, drzewo, limit)
}

// Pompy ruszają przed czekaniem, inaczej bufor potoku zakleszcza proces.
func zbierz(ctx context.Context, n Narzedzie, uchwyt session.UchwytProcesu,
	drzewo *session.DrzewoProcesu, limit time.Duration) (Wynik, error) {

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

	wynik := Wynik{Wyjscie: <-wyjscie, Diagnostyka: string(<-diagnostyka)}

	if powod != "" {
		// Wynik przerwany jest niepełny, więc nie wolno oddać go jako udanego.
		return wynik, errors.New("arsenał: " + n.Nazwa + " przerwany — " + powod +
			"; naprawa: podnieść granicę czasu albo podać mniejszy materiał")
	}
	if bladZakonczenia != nil {
		return wynik, errors.New("arsenał: " + n.Nazwa +
			" zakończył się niepowodzeniem: " + bladZakonczenia.Error() +
			opisDiagnostyki(wynik.Diagnostyka))
	}
	return wynik, nil
}

func opisDiagnostyki(diagnostyka string) string {
	tresc := strings.TrimSpace(diagnostyka)
	if tresc == "" {
		return ""
	}
	if len(tresc) > 400 {
		tresc = tresc[:400] + "…"
	}
	return "; program powiedział: " + tresc
}

func czytajCalosc(zrodlo io.Reader) []byte {
	if zrodlo == nil {
		return nil
	}
	bajty, _ := io.ReadAll(zrodlo)
	return bajty
}
