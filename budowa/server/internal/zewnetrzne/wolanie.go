// Pakiet zewnetrzne jest JEDNĄ drogą, którą rdzeń woła binarium arsenału:
// ImageMagick, libvips, ffmpeg, tesseract, pandoc, poppler, 7z i każde następne.
//
// Po co powstał. Produkt daje modelowi arsenał czynności — obraz, dźwięk,
// dokument, archiwum — a większość z nich to opakowanie dojrzałego programu,
// nie pisanie go od nowa w Go. Bez tego pakietu każda rodzina narzędzi
// zbudowałaby własne uruchamianie procesu, własny limit czasu i własne
// sprzątanie potomstwa — cztery prawdy o jednej rzeczy, a przy czwartej ktoś
// zapomniałby ubić drzewo.
//
// Sekwencja jest przepisana z `mowa/uruchomienie.go`. W całym drzewie stoi
// dokładnie jedno `exec.Command`
// (`injection/rozruch.go`); wszystko inne idzie portem `session.Uruchamiacz`,
// przez tę samą bramę izolacji okna i to samo obejmowanie potomstwa. Każde
// odstępstwo od tej sekwencji jest wyciekiem procesu albo uchwytu — a binarium
// przetwarzające film rozgałęzia wątki tak samo jak dekoder mowy.
//
// Czego ten pakiet nie robi. Nie wie, co znaczy wyjście programu: nie rozpoznaje
// formatów, nie czyta obrazów i nie tłumaczy komunikatów ImageMagicka na polski.
// Oddaje bajty i kod wyjścia; rozpoznanie należy do rodziny narzędzi, która zna
// kształt odpowiedzi swojego programu.
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
	// Nazwa czytelna, np. "ImageMagick".
	Nazwa string
	// Program to nazwa binarium w PATH albo ścieżka bezwzględna.
	Program string
	// Pakiet podpowiada, czym je dociągnąć — wchodzi do treści odmowy, żeby
	// Operator nie musiał szukać. Puste pomija podpowiedź.
	Pakiet string
}

// Wynik niesie surowy rezultat jednego uruchomienia.
//
// Dwa strumienie oddzielnie, tak samo jak w silniku mowy: na wyjściu stoi
// wynik pracy (bywa nim binarna treść obrazu), na diagnostyce ostrzeżenia
// programu. Sklejenie ich wstawiłoby ostrzeżenie w środek pliku PNG.
type Wynik struct {
	Wyjscie     []byte
	Diagnostyka string
}

// BrakNarzedzia mówi, że binarium nie stoi na maszynie.
//
// Osobny typ, bo to odmowa innej klasy. „Nie ma czym" jest BRAKIEM, który
// Operator usuwa jedną instalacją, a nie usterką rdzenia — i rodzina narzędzi
// ma go odróżnić, żeby powiedzieć Operatorowi, co dociągnąć.
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
//
// Limit czasu jest obowiązkowy — tak samo jak w silniku mowy. Transkodowanie
// filmu bez granicy potrafi zająć maszynę na godziny, a program, który utknął
// na uszkodzonym pliku, nie kończy się nigdy. Limit niedodatni to odmowa,
// nie „bez granicy": wołający, który granicy nie podał, o niej zapomniał.
//
// Środowisko nie jest dziedziczone. Silnik mowy dziedziczy je, bo interpreter
// Pythona musi odnaleźć własne biblioteki; narzędzia arsenału są binariami
// samodzielnymi i nie mają powodu widzieć zmiennych rdzenia. Gdy okno ma
// włączony punkt izolacji środowiska, brama niżej i tak rozstrzyga — ale
// domyślna wartość ma być węższa, nie szersza.
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

	// Do uruchomienia idzie ścieżka odnaleziona, nie sama nazwa: program
	// dołożony do pakietu leży poza ścieżką wyszukiwania systemu i po nazwie
	// nie wystartowałby.
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
		// Proces biegnie, a uchwytu drzewa nie ma — zostawienie go tak znaczyłoby
		// sierotę poza rejestrem rdzenia.
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
//
// Pompy ruszają przed czekaniem. Bufor potoku ma kilkadziesiąt kilobajtów;
// program piszący obraz na wyjście zablokowałby się na zapisie, a rdzeń czekałby
// na koniec procesu, który czeka na rdzeń. Zakleszczenie kończy się dopiero
// granicą czasu i wygląda jak zawieszone narzędzie.
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

	// Odbiór z obu pomp PO zakończeniu i BEZWARUNKOWO, także przy przerwaniu:
	// gorutyny zostałyby inaczej zawieszone na zapisie do kanału.
	wynik := Wynik{Wyjscie: <-wyjscie, Diagnostyka: string(<-diagnostyka)}

	if powod != "" {
		// Wynik przerwany jest niepełny, więc NIE WOLNO oddać go jako udanego.
		// Plik obrazu urwany w połowie jest gorszy niż jego brak.
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
//
// Bez tego członu Operator czyta „narzędzie zakończyło się niepowodzeniem"
// i nie wie nic. Diagnostyka bywa długa, więc bierzemy jej początek — pierwsze
// zdanie programu prawie zawsze niesie powód, a reszta jest śladem stosu.
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
