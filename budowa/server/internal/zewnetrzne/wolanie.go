// Pakiet zewnetrzne jest jedną drogą, którą rdzeń woła binarium arsenału:
// ImageMagick, libvips, ffmpeg, tesseract, pandoc, poppler, 7z i każde
// następne. Nie wie, co znaczy wyjście programu: oddaje bajty i kod wyjścia.
package zewnetrzne

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
)

// Narzedzie opisuje jedno binarium arsenału.
//
// `Nazwa` jest nazwą dla CZŁOWIEKA i wchodzi do treści odmowy — Operator ma
// przeczytać „brakuje ImageMagick (convert)", a nie samą ścieżkę programu.
type Narzedzie struct {
	// Nazwa czytelna, na przykład "ImageMagick".
	Nazwa string
	// Program to nazwa binarium w PATH albo ścieżka bezwzględna.
	Program string
	// Pakiet podpowiada, czym je dociągnąć. Puste pomija podpowiedź.
	Pakiet string
}

// Wynik niesie surowy rezultat jednego uruchomienia. Dwa strumienie
// oddzielnie: na wyjściu stoi wynik pracy, na diagnostyce ostrzeżenia
// programu.
type Wynik struct {
	Wyjscie     []byte
	Diagnostyka string
}

// BrakNarzedzia mówi, że binarium nie stoi na maszynie. Osobny typ, bo to
// odmowa innej klasy: brak, który Operator usuwa jedną instalacją, nie
// usterka rdzenia.
type BrakNarzedzia struct {
	Narzedzie Narzedzie
}

func (b *BrakNarzedzia) Error() string {
	zdanie := "arsenał: nie ma na tej maszynie programu " + b.Narzedzie.Nazwa +
		" (" + b.Narzedzie.Program + "), a rdzeń tej czynności nie wykona bez niego"
	if b.Narzedzie.Pakiet != "" {
		zdanie += "; naprawa: zainstalować pakiet " + b.Narzedzie.Pakiet
	}
	return zdanie
}

// Stoi sprawdza obecność binarium — w pakiecie produktu albo na ścieżce systemu.
//
// Sprawdzenie jest tanie i idzie przed uruchomieniem, bo odmowa „nie ma czym"
// jest dla Operatora czymś innym niż „program wystartował i się wywrócił".
func Stoi(n Narzedzie) bool {
	_, jest := Odnajdz(n)
	return jest
}

// Wolaj przeprowadza jedno uruchomienie narzędzia przez port session.Uruchamiacz.
// Limit czasu jest obowiązkowy. Środowisko nie jest dziedziczone, bo
// narzędzia arsenału są binariami samodzielnymi.
func Wolaj(ctx context.Context, u session.Uruchamiacz, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	n Narzedzie, argumenty []string, katalog string, limit time.Duration) (Wynik, error) {

	if u == nil {
		return Wynik{}, errors.New("arsenał: rdzeń nie ma uruchamiacza procesów — " +
			"narzędzia zewnętrzne nie mają czym wystartować; " +
			"naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia")
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

	// Do uruchomienia idzie ścieżka odnaleziona, nie sama nazwa programu.
	polecenie := session.Polecenie{
		Program:   sciezka,
		Argumenty: argumenty,
		Katalog:   katalog,
	}
	dopuszczone, err := session.SprawdzPolecenie(zasady, obszar, polecenie)
	if err != nil {
		return Wynik{}, err
	}

	uchwyt, err := u.UruchomProces(okno, dopuszczone)
	if err != nil {
		return Wynik{}, errors.New("arsenał: nie można uruchomić " + n.Nazwa +
			" (" + n.Program + ") w " + dopuszczone.Katalog + ": " + err.Error())
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		// Proces biegnie, a uchwytu drzewa nie ma — trzeba go ubić od razu.
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		return Wynik{}, errors.New("arsenał: nie można objąć drzewa procesu " +
			n.Nazwa + ": " + err.Error())
	}
	defer drzewo.Zwolnij()

	return zbierz(ctx, n, uchwyt, drzewo, limit)
}

// zbierz prowadzi uruchomiony proces do końca: pompuje oba strumienie, czeka
// z granicą czasu i ubija całe drzewo, gdy granica albo rdzeń każą przerwać.
// Pompy ruszają przed czekaniem, żeby uniknąć zakleszczenia na buforze potoku.
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
		powod = "żądanie przerwane przez rdzeń"
		_ = drzewo.Ubij()
		bladZakonczenia = <-zakonczenie
	}

	// Odbiór z obu pomp po zakończeniu i bezwarunkowo, także przy przerwaniu.
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

// opisDiagnostyki dokłada do odmowy to, co program powiedział o sobie sam.
// Diagnostyka bywa długa, więc bierze się jej początek.
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

// czytajCalosc zbiera cały strumień procesu. Błąd odczytu oddaje to, co zdążyło
// przyjść — tak samo jak w silniku mowy i w adapterze Developera.
func czytajCalosc(zrodlo io.Reader) []byte {
	if zrodlo == nil {
		return nil
	}
	bajty, _ := io.ReadAll(zrodlo)
	return bajty
}
