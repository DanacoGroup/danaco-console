// Zbiorcze wyjście modułu Terminal: dziennik wierszy ze wszystkich otwartych
// kart powłoki oraz ewidencja okien zapisanych na ten strumień. Na tym stoi
// komenda `terminal.output.stream`
// (`adapter_modul_terminal_wyjscie_komenda.go`).
//
// Byt osobny od `adapter_modul_terminal_strumien.go`: tamten rozsyła wyjście
// jednego procesu do okna, w którym stoi karta (`WindowId: proces.oknoKod`),
// a tutaj zbiorcze wyjście wszystkich kart odbiera okno wskazane w żądaniu,
// wraz z ogonem historii. Okno Output Console nie musi być oknem karty.
//
// Dziennik nie czyta potoków procesów po raz drugi: owija nadajnik, którym
// istniejąca pompa już wysyła fragmenty (`nadajnikZDziennikiem`). Fragment idzie
// swoją dotychczasową drogą nietknięty, a do dziennika i do okna obserwatora
// wchodzi jego kopia. Dwie pompy na jeden potok czytałyby sobie nawzajem bajty.
//
// Fragment nie jest wierszem: pompa czyta bloki po 32 KiB, a kontrakt niesie
// `TerminalOutputLine`, czyli wiersz. Blok bywa urwany w połowie wiersza, więc
// resztka czeka tu na ciąg dalszy, osobno dla wyjścia zwykłego i diagnostycznego
// (dwie gorutyny, dwa niezależne strumienie). Resztka bez znaku końca linii
// dłuższa niż `maksDlugoscWierszaWyjscia` idzie do dziennika jako wiersz mimo
// wszystko, bo pasek postępu nie kończy się nigdy, a zużycie pamięci ma pozostać
// ograniczone.
//
// Wierszy wyjścia nie ma gdzie zapisać: schemat bazy ma tabele `terminal_karta`
// i `terminal_proces`, żadnej tabeli wyjścia. Dziennik jest więc pierścieniem
// w pamięci, tak samo jak pamięciowa warstwa historii procesów
// (`pojemnoscHistorii` w `adapter_modul_terminal_rejestr.go`). Po restarcie ogon
// jest pusty.
package core

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// pojemnoscDziennikaWyjscia ogranicza liczbę wierszy trzymanych w pamięci.
// Powyżej progu wypada wiersz najstarszy.
const pojemnoscDziennikaWyjscia = 5000

// domyslnyOgonWyjscia obowiązuje, gdy żądanie nie poda własnej liczby wierszy.
const domyslnyOgonWyjscia = 200

// maksDlugoscWierszaWyjscia domyka wiersz, którego proces nie domyka sam.
const maksDlugoscWierszaWyjscia = 8 * 1024

// przypisanieWyjscia wiąże proces z kartą, oknem i sesją. Zapamiętujemy je przy
// pierwszym fragmencie, bo proces zakończony wypada kiedyś z rejestru żywego
// (`rejestrTerminala.Przytnij`), a jego ostatnie wiersze mają wtedy nadal
// wskazywać kartę, z której wyszły.
type przypisanieWyjscia struct {
	kartaKod string
	oknoKod  string
	idSesji  string
}

// wpisWyjscia jest wierszem dziennika: byt kontraktu wraz z dwoma polami,
// których kontrakt w wierszu nie niesie, a po których zawęża żądanie — sesją
// i oknem karty.
type wpisWyjscia struct {
	wiersz  shared.TerminalOutputLine
	idSesji string
	oknoKod string
}

// filtrWyjscia zawęża zbiorcze wyjście do kart wskazanych w żądaniu. Pole puste
// nie zawęża niczego — strumień bez wskazań jest strumieniem pełnym.
type filtrWyjscia struct {
	idSesji string
	karty   map[string]struct{}
}

// przepuszcza sprawdza wiersz obiema warunkami filtru.
func (f filtrWyjscia) przepuszcza(w wpisWyjscia) bool {
	if len(f.karty) > 0 {
		if _, jest := f.karty[w.wiersz.TerminalSessionId]; !jest {
			return false
		}
	}
	return f.idSesji == "" || w.idSesji == f.idSesji
}

// obserwacjaWyjscia jest jednym oknem zapisanym na zbiorcze wyjście.
type obserwacjaWyjscia struct {
	oknoKod string
	idSesji string
	filtr   filtrWyjscia
	// numer prowadzi ciąg numerów fragmentów tego okna. Ciąg jest własny, bo
	// okno obserwatora dostaje wyjście wielu procesów naraz, a kontrakt
	// numeruje fragmenty w obrębie jednego strumienia.
	numer atomic.Int64
}

// dziennikWyjscia trzyma ogon zbiorczego wyjścia i wykaz okien obserwujących.
type dziennikWyjscia struct {
	mu sync.Mutex
	// wpisy jest pierścieniem wierszy w kolejności wypisania.
	wpisy []wpisWyjscia
	// niedokonczone trzyma resztkę wiersza po kluczu proces+kanał.
	niedokonczone map[string]string
	// znane pamięta przypisanie procesu do karty, okna i sesji.
	znane map[string]przypisanieWyjscia
	// obserwatorzy wiąże okno odbierające z jego zawężeniem.
	obserwatorzy map[string]*obserwacjaWyjscia

	// rejestr służy wyłącznie odczytaniu przypisania procesu przy pierwszym
	// fragmencie. Dziennik niczego w nim nie zmienia.
	rejestr *rejestrTerminala
	// dalej jest nadajnikiem, którym idzie kopia wiersza do okna obserwatora.
	dalej Nadajnik
}

// nowyDziennikWyjscia zakłada pusty dziennik zbiorczego wyjścia.
func nowyDziennikWyjscia(rejestr *rejestrTerminala) *dziennikWyjscia {
	return &dziennikWyjscia{
		niedokonczone: make(map[string]string),
		znane:         make(map[string]przypisanieWyjscia),
		obserwatorzy:  make(map[string]*obserwacjaWyjscia),
		rejestr:       rejestr,
	}
}

// Obserwuj zapisuje okno na zbiorcze wyjście. Ponowny zapis tego samego okna
// podmienia zawężenie zamiast dokładać drugą obserwację: okno ma jeden widok
// wyjścia niezależnie od liczby otwarć.
func (d *dziennikWyjscia) Obserwuj(oknoKod, idSesji string, filtr filtrWyjscia) {
	if d == nil || oknoKod == "" {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.obserwatorzy[oknoKod] = &obserwacjaWyjscia{oknoKod: oknoKod, idSesji: idSesji, filtr: filtr}
}

// Zapomnij wykreśla obserwacje okien, których warunek już nie obejmuje —
// służy odsianiu okien zamkniętych. Nie ma zdarzenia „okno przestało słuchać”,
// więc porządkowanie odbywa się przy kolejnym żądaniu.
func (d *dziennikWyjscia) Zapomnij(zostaw func(oknoKod string) bool) {
	if d == nil || zostaw == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	for kod := range d.obserwatorzy {
		if !zostaw(kod) {
			delete(d.obserwatorzy, kod)
		}
	}
}

// Ogon zwraca ostatnie `ile` wierszy spełniających filtr, w kolejności
// wypisania. Wynik jest zawsze listą — pusty ogon to pusta lista, nie brak pola.
func (d *dziennikWyjscia) Ogon(filtr filtrWyjscia, ile int) []shared.TerminalOutputLine {
	wynik := make([]shared.TerminalOutputLine, 0, 16)
	if d == nil || ile <= 0 {
		return wynik
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	// Przebieg od końca bierze ostatnie pasujące wiersze; odwrócenie wyniku
	// przywraca kolejność wypisania, której żąda kontrakt.
	for i := len(d.wpisy) - 1; i >= 0 && len(wynik) < ile; i-- {
		if filtr.przepuszcza(d.wpisy[i]) {
			wynik = append(wynik, d.wpisy[i].wiersz)
		}
	}
	for lewy, prawy := 0, len(wynik)-1; lewy < prawy; lewy, prawy = lewy+1, prawy-1 {
		wynik[lewy], wynik[prawy] = wynik[prawy], wynik[lewy]
	}
	return wynik
}

// przyjmij odkłada w dzienniku wiersze niesione przez fragment strumienia.
// Kopertę innego rodzaju niż `stream.chunk` przepuszcza bez śladu.
func (d *dziennikWyjscia) przyjmij(k protocol.Koperta) {
	if d == nil || k.Type != shared.EventStreamChunk {
		return
	}
	fragment, err := protocol.FragmentZKoperty(k)
	if err != nil {
		// Fragment nieczytelny dla dziennika poszedł już swoją drogą do okna
		// karty, więc pominięcie go kosztuje tylko brak wiersza w ogonie.
		return
	}
	kodProcesu := strings.TrimSpace(fragment.MessageId)
	if kodProcesu == "" {
		return
	}
	przypisanie, zna := d.przypisanie(kodProcesu)
	if !zna {
		// Wiersz bez karty nie ma jak wejść do zbiorczego wyjścia: pole
		// `terminalSessionId` jest w kontrakcie wymagane, a wpisanie tam pustki
		// byłoby wskazaniem karty, której nie ma.
		return
	}

	kanal := shared.TerminalOutputChannel(shared.TerminalOutputChannelStdout)
	if fragment.Kind == shared.ChunkKindError {
		kanal = shared.TerminalOutputChannelStderr
	}
	for _, tekst := range d.zloz(kodProcesu, kanal, protocol.Tresc(fragment), protocol.Ostatni(k)) {
		d.odloz(przypisanie, kodProcesu, kanal, tekst)
	}
}

// przypisanie zwraca kartę, okno i sesję procesu — z pamięci dziennika albo,
// przy pierwszym fragmencie, z rejestru żywego.
func (d *dziennikWyjscia) przypisanie(kodProcesu string) (przypisanieWyjscia, bool) {
	d.mu.Lock()
	zapamietane, jest := d.znane[kodProcesu]
	d.mu.Unlock()
	if jest {
		return zapamietane, true
	}
	if d.rejestr == nil {
		return przypisanieWyjscia{}, false
	}
	proces, biegnie := d.rejestr.Proces(kodProcesu)
	if !biegnie || proces.kartaKod == "" {
		return przypisanieWyjscia{}, false
	}
	zapamietane = przypisanieWyjscia{
		kartaKod: proces.kartaKod,
		oknoKod:  proces.oknoKod,
		idSesji:  proces.idSesji,
	}
	d.mu.Lock()
	d.znane[kodProcesu] = zapamietane
	d.mu.Unlock()
	return zapamietane, true
}

// zloz zamienia porcję bajtów na wiersze domknięte. Resztka bez znaku końca
// linii czeka na ciąg dalszy; fragment ostatni domyka ją siłą, bo ciągu dalszego
// już nie będzie.
func (d *dziennikWyjscia) zloz(kodProcesu string, kanal shared.TerminalOutputChannel,
	tresc string, ostatni bool) []string {

	klucz := kodProcesu + "\x00" + string(kanal)
	d.mu.Lock()
	defer d.mu.Unlock()

	bufor := d.niedokonczone[klucz] + tresc
	czesci := strings.Split(bufor, "\n")
	wiersze := make([]string, 0, len(czesci))
	for i := 0; i < len(czesci)-1; i++ {
		wiersze = append(wiersze, strings.TrimSuffix(czesci[i], "\r"))
	}
	reszta := czesci[len(czesci)-1]

	switch {
	case ostatni:
		if reszta != "" {
			wiersze = append(wiersze, strings.TrimSuffix(reszta, "\r"))
		}
		delete(d.niedokonczone, klucz)
		delete(d.znane, kodProcesu)
	case len(reszta) >= maksDlugoscWierszaWyjscia:
		wiersze = append(wiersze, reszta)
		delete(d.niedokonczone, klucz)
	case reszta == "":
		delete(d.niedokonczone, klucz)
	default:
		d.niedokonczone[klucz] = reszta
	}
	return wiersze
}

// odloz wpisuje jeden wiersz do pierścienia i rozsyła go obserwatorom.
func (d *dziennikWyjscia) odloz(przypisanie przypisanieWyjscia, kodProcesu string,
	kanal shared.TerminalOutputChannel, tekst string) {

	proces := kodProcesu
	wpis := wpisWyjscia{
		wiersz: shared.TerminalOutputLine{
			TerminalSessionId: przypisanie.kartaKod,
			ProcessId:         &proces,
			Channel:           kanal,
			Text:              tekst,
			At:                time.Now().UnixMilli(),
		},
		idSesji: przypisanie.idSesji,
		oknoKod: przypisanie.oknoKod,
	}

	d.mu.Lock()
	d.wpisy = append(d.wpisy, wpis)
	if nadmiar := len(d.wpisy) - pojemnoscDziennikaWyjscia; nadmiar > 0 {
		d.wpisy = append(d.wpisy[:0], d.wpisy[nadmiar:]...)
	}
	odbiorcy := make([]*obserwacjaWyjscia, 0, len(d.obserwatorzy))
	for _, obserwacja := range d.obserwatorzy {
		odbiorcy = append(odbiorcy, obserwacja)
	}
	dalej := d.dalej
	d.mu.Unlock()

	if dalej == nil {
		return
	}
	for _, obserwacja := range odbiorcy {
		// Okno karty dostało ten wiersz swoją drogą — drugi raz dostawałoby go
		// podwójnie i Output Console pokazywałaby każdy wiersz dwa razy.
		if obserwacja.oknoKod == wpis.oknoKod || !obserwacja.filtr.przepuszcza(wpis) {
			continue
		}
		d.wyslijDoOkna(dalej, obserwacja, wpis, kanal)
	}
}

// wyslijDoOkna oddaje wiersz oknu obserwatora tą samą drogą, którą idzie całe
// wyjście modułu — zdarzeniem `stream.chunk`. Kontrakt nie ma osobnego
// zdarzenia zbiorczego wyjścia.
func (d *dziennikWyjscia) wyslijDoOkna(dalej Nadajnik, obserwacja *obserwacjaWyjscia,
	wpis wpisWyjscia, kanal shared.TerminalOutputChannel) {

	rodzaj := shared.ChunkKind(shared.ChunkKindText)
	if kanal == shared.TerminalOutputChannelStderr {
		rodzaj = shared.ChunkKindError
	}
	kodProcesu := ""
	if wpis.wiersz.ProcessId != nil {
		kodProcesu = *wpis.wiersz.ProcessId
	}
	fragment := protocol.ChunkTekstu(obserwacja.oknoKod, kodProcesu, wpis.wiersz.Text+"\n")
	fragment.Kind = rodzaj
	koperta, err := protocol.KopertaFragmentu(kodProcesu, obserwacja.idSesji,
		int(obserwacja.numer.Add(1)), false, fragment)
	if err != nil {
		return
	}
	dalej.Rozglos(koperta)
}

// nadajnikZDziennikiem owija nadajnik pompy wyjścia: fragment idzie dalej
// nietknięty, a jego kopia wchodzi do dziennika zbiorczego.
//
// Kolejność jest istotna — najpierw droga do okna karty, potem dziennik — żeby
// składanie wierszy i rozsyłka do obserwatorów nie opóźniały wyjścia karty.
type nadajnikZDziennikiem struct {
	dalej    Nadajnik
	dziennik *dziennikWyjscia
}

// Rozglos wypełnia port Nadajnik.
func (n *nadajnikZDziennikiem) Rozglos(k protocol.Koperta) {
	if n.dalej != nil {
		n.dalej.Rozglos(k)
	}
	n.dziennik.przyjmij(k)
}
