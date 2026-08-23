// Odpowiedzialność pliku: przeprowadzenie pomocnika transkrypcji przez jedyną
// dozwoloną drogę startu procesu i zebranie tego, co powiedział.
//
// Pomocnik startuje tym samym uruchamiaczem, co okno rozmowy, terminal i git:
// w całym drzewie jest dokładnie jedno `exec.Command` (injection/rozruch.go).
// Dzięki temu pomocnik przechodzi przez tę samą bramę izolacji okna i przez to
// samo obejmowanie potomstwa, co pozostałe procesy — a proces Pythona, który
// rozgałęzia własne wątki dekodera, bez objęcia drzewem zostawiałby sieroty.
//
// Sekwencja jest ta sama, co w `core/adapter_modul_terminal_bieg.go`
// i `core/adapter_modul_developer_git_wykonanie.go`: strażnik nil → sprawdzenie
// izolacji → UruchomProces → PrzejmijDrzewo(pid) → defer Zwolnij → pompy obu
// strumieni → select na Czekaj/timer/ctx.Done → Ubij całego drzewa przy
// przekroczeniu. Odstępstwo od niej kończy się wyciekiem procesu albo uchwytu.
//
// Czego ten plik nie robi: nie rozstrzyga, czy transkrypcja się udała. Oddaje
// wyjście i diagnostykę takie, jakie przyszły; rozpoznanie braku silnika czy
// pustej odpowiedzi należy do warstwy, która zna kształt odpowiedzi pomocnika.
package mowa

import (
	"context"
	"errors"
	"io"
	"time"

	"danacoconsole/server/internal/session"
)

// Wynik niesie surowy rezultat jednego uruchomienia pomocnika.
//
// Dwa strumienie oddzielnie, a nie sklejone w jeden napis: na wyjściu stoi
// transkrypcja przeznaczona dla Operatora, a na diagnostyce ostrzeżenia
// Pythona i komunikat nieudanego importu. Sklejenie ich wstawiłoby ostrzeżenie
// biblioteki w środek przepisanego zdania, a rozpoznanie braku silnika
// (BrakSilnika) straciłoby jedyne miejsce, w którym go widać.
type Wynik struct {
	// Wyjscie — bajty ze standardowego wyjścia pomocnika.
	Wyjscie []byte
	// Diagnostyka — tekst z wyjścia diagnostycznego pomocnika.
	Diagnostyka string
}

// Uruchom przeprowadza pomocnika przez port session.Uruchamiacz.
//
// `zasady` i `obszar` są w sygnaturze, bo `session.SprawdzPolecenie(zasady,
// obszar, polecenie)` innego kształtu nie ma, a bez bramy izolacji pomocnik
// startowałby jako jedyny proces w drzewie poza katalogiem roboczym okna
// i z odziedziczonym środowiskiem rdzenia. Wołający ma te dwie wartości: rdzeń
// wylicza je tak samo dla Terminala i Developera (ZasadyIzolacji + obszar okna).
//
// Limit czasu jest obowiązkowy. Transkrypcja godzinnego nagrania na słabym
// procesorze trwa dłużej niż samo nagranie, a pomocnik, który utknął na
// pobieraniu modelu, nie kończy się nigdy. Limit niedodatni to odmowa, nie
// „bez granicy”: domyślenie granicy za wołającego ukryłoby brak wskazania aż do
// pierwszego zawieszonego procesu na maszynie Operatora.
func Uruchom(ctx context.Context, u session.Uruchamiacz, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar,
	p Pomocnik, argumenty []string, katalog string, limit time.Duration) (Wynik, error) {

	if u == nil {
		return Wynik{}, errors.New("silnik mowy: rdzeń nie ma uruchamiacza procesów" +
			" — pomocnik transkrypcji nie ma czym wystartować;" +
			" naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia")
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
		// Pomocnik potrzebuje środowiska systemu, żeby interpreter odnalazł własne
		// biblioteki. Gdy okno ma włączony punkt izolacji środowiska, brama niżej
		// to odrzuci — i tak ma być: rozstrzyga ustawienie Operatora, nie ten plik.
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
		// Proces już biegnie, a uchwytu drzewa nie ma — zostawienie go tak
		// znaczyłoby sierotę poza rejestrem rdzenia.
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		return Wynik{}, errors.New("silnik mowy: nie można objąć drzewa procesu pomocnika: " + err.Error())
	}
	defer drzewo.Zwolnij()

	return zbierz(ctx, uchwyt, drzewo, limit)
}

// zbierz prowadzi uruchomiony proces do końca: pompuje oba strumienie, czeka
// z granicą czasu i ubija całe drzewo, gdy granica albo rdzeń każą przerwać.
//
// Pompy ruszają przed czekaniem, nie po nim. Bufor potoku ma kilkadziesiąt
// kilobajtów; pomocnik piszący dłuższą transkrypcję zablokowałby się na zapisie,
// a rdzeń czekałby na koniec procesu, który czeka na rdzeń — zakleszczenie,
// które kończy się dopiero granicą czasu i wygląda jak zawieszony silnik.
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
		powod = "żądanie przerwane przez rdzeń"
		_ = drzewo.Ubij()
		bladZakonczenia = <-zakonczenie
	}

	// Odbiór z obu pomp następuje po zakończeniu procesu i bezwarunkowo, także
	// przy przerwaniu: gorutyny zostałyby inaczej zawieszone na zapisie do
	// kanału, a to wyciek na każde przekroczenie granicy.
	wynik := Wynik{Wyjscie: <-wyjscie, Diagnostyka: string(<-diagnostyka)}

	if powod != "" {
		// Wynik przerwany jest niepełny, więc nie idzie jako udany. Diagnostyka
		// jedzie razem z odmową, bo zwykle stoi w niej jedyne zdanie mówiące,
		// na czym pomocnik utknął.
		return wynik, errors.New("silnik mowy: transkrypcja przerwana — " + powod +
			"; naprawa: podnieść granicę czasu albo podać krótsze nagranie")
	}
	if bladZakonczenia != nil {
		return wynik, errors.New("silnik mowy: pomocnik zakończył się niepowodzeniem: " +
			bladZakonczenia.Error())
	}
	return wynik, nil
}

// czytajCalosc zbiera cały strumień procesu. Wyjście pomocnika to jedna
// transkrypcja, a nie ciągnący się dziennik, więc idzie w całości do wyniku —
// przyrostowe zdarzenia byłyby tu kosztem bez odbiorcy. Przy błędzie odczytu
// zwracane jest to, co zdążyło przyjść — niepełna transkrypcja niesie więcej
// niż pusty wynik.
func czytajCalosc(zrodlo io.Reader) []byte {
	if zrodlo == nil {
		return nil
	}
	bajty, _ := io.ReadAll(zrodlo)
	return bajty
}
