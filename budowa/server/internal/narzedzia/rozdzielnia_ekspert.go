// Odpowiedzialność pliku: droga zasięgu eksperta wewnątrz rozdzielni —
// odczyt definicji, złożenie podzbioru, pomiar i meldunek.
//
// Osobno od `rozdzielnia.go`, bo droga zwykła nie ma prawa się o to potknąć:
// bez `--ekspert` ta cała maszyneria nie rusza ani razu i rozdzielnia zachowuje
// się co do znaku tak, jak bez niej.
//
// Dlaczego podzbiór zapamiętuje się, a niepowodzenie nie.
// Definicję czyta się przy pierwszym `tools/list`, bo gniazdo do rdzenia jest
// leniwe z zamysłem (`polaczenie.go`). Odczyt udany zapamiętuje się na czas
// życia procesu: ekspert okna nie zmienia się w trakcie tury, a pytanie rdzenia
// przy każdym wykazie byłoby ruchem bez treści. Odczyt nieudany nie zapamiętuje
// się nigdy — rdzeń bywa niegotowy w chwili startu procesu modelu, więc każde
// następne `tools/list` próbuje od nowa. Zapamiętana porażka zamieniłaby jedno
// nieudane połączenie w oknie bez doboru narzędzi na całą turę.
//
// Zawężenie obowiązuje także przy wywołaniu.
// Wykaz zawężony, z którego nadal da się wywołać wszystko, jest zawężeniem
// pozornym — oszczędza żetony i nie zmienia niczego więcej. Dlatego `Wywolaj`
// pyta o przynależność do podzbioru i odmawia nazwie spoza niego, mówiąc wprost,
// który ekspert jej nie niesie. Odmowa jest treścią dla modelu, nie usterką
// procesu.
package narzedzia

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// doborEksperta prowadzi stan doboru narzędzi jednego okna: kod z przełącznika,
// zapamiętany wynik złożenia i dziennik, na który idą meldunki.
type doborEksperta struct {
	kod string
	// dolozenia są nazwami dołożonymi doraźnie w sesji, spoza definicji eksperta
	// (`--dolozenia`, `injection/zestaw_narzedzi.go`).
	dolozenia []string
	dziennik  *log.Logger
	zamek     sync.Mutex
	// zlozony jest zapamiętanym wynikiem udanego odczytu definicji; nil znaczy
	// „jeszcze nie czytano albo odczyt się nie udał".
	zlozony *WykazEksperta
	// zameldowane pilnuje, żeby ten sam powód braku zawężenia nie szedł do
	// dziennika przy każdym `tools/list`: powtarzany meldunek zasłania resztę
	// dziennika i sam przestaje być czytany (tak samo rozstrzyga `core/
	// most_narzedzi.go`).
	zameldowane map[string]bool
}

// wykaz zwraca wykaz podawany modelowi w zasięgu eksperta.
//
// `zrodlo` jest wykazem, który dostałoby to okno bez eksperta. Każda droga —
// udana i nieudana — kończy się meldunkiem z liczbami, bo cena zestawu jest
// tym, czego Operator ma się dowiedzieć przed limitem, a nie po nim.
func (d *doborEksperta) wykaz(kontekst context.Context, rdzen Rdzen, zrodlo []Narzedzie) []Narzedzie {
	d.zamek.Lock()
	defer d.zamek.Unlock()

	if d.zlozony != nil {
		return d.zlozony.Pozycje
	}
	definicja, err := OdczytajEksperta(kontekst, rdzen, d.kod)
	if err != nil {
		wynik := ZDolozeniami(WykazBezEksperta(d.kod, err, zrodlo), d.dolozenia, zrodlo)
		d.zamelduj(wynik)
		return wynik.Pozycje // bez zapamiętania: następne `tools/list` próbuje od nowa
	}
	wynik := ZDolozeniami(ZlozWykazEksperta(definicja, zrodlo), d.dolozenia, zrodlo)
	d.zlozony = &wynik
	d.zamelduj(wynik)
	return wynik.Pozycje
}

// wolno rozstrzyga, czy nazwa narzędzia mieści się w podzbiorze eksperta.
//
// Przed pierwszym udanym odczytem i przy wykazie niezawężonym wolno wszystko:
// zasięgiem jest wtedy wykaz okna, więc zawężać wywołań nie ma czym. Prawda
// druga mówi, czy odpowiedź opiera się na znanym podzbiorze — wywołujący
// potrzebuje jej, żeby odmowa nazwała powód.
func (d *doborEksperta) wolno(nazwa string) (dopuszczone, znany bool) {
	d.zamek.Lock()
	defer d.zamek.Unlock()

	if d.zlozony == nil || !d.zlozony.Zawezony {
		return true, false
	}
	for _, pozycja := range d.zlozony.Pozycje {
		if pozycja.Nazwa == nazwa {
			return true, true
		}
	}
	return false, true
}

// odmowaPozaDoborem buduje treść odmowy dla narzędzia spoza podzbioru eksperta.
//
// Mówi wprost, że rzecz jest w doborze, a nie w granicy uprawnień modelu
// (`zastrzezenia.go`): to są dwie różne odmowy o dwóch różnych naprawach —
// tamtą naprawia się kontraktem, tę zmianą definicji eksperta.
func (d *doborEksperta) odmowaPozaDoborem(nazwa string) error {
	return fmt.Errorf("narzędzie %q stoi poza doborem eksperta %q: okno pracuje w zasięgu eksperta,"+
		" a jego skillIds i connectorIds tego narzędzia ani jego grupy nie wskazują."+
		" Komplet narzędzi tego okna zwraca tools/list; rozszerza go Operator w Agent Builderze,"+
		" nie model w rozmowie", nazwa, d.kod)
}

// zamelduj kładzie na dziennik cenę zestawu i powód braku zawężenia.
//
// Meldunek idzie zawsze przy pierwszym złożeniu — także przy zawężeniu udanym,
// bo liczba pozycji i bajtów jest tym, po co licznik powstał. Powtórka tego
// samego powodu jest tłumiona, sam pomiar nie: pomiar pada raz, bo złożenie
// udane zapamiętuje się na czas życia procesu.
func (d *doborEksperta) zamelduj(wynik WykazEksperta) {
	if d.dziennik == nil {
		return
	}
	pomiar := Zmierz(wynik.Pozycje)
	if wynik.Zawezony {
		d.dziennik.Printf("dobór eksperta %q: %d narzędzi w %d grupach, %d bajtów wykazu tools/list",
			d.kod, pomiar.Pozycji, len(pomiar.Grupy), pomiar.Bajtow)
	} else if d.raz(wynik.Powod) {
		d.dziennik.Printf("BRAK ZAWĘŻENIA: %s. Model dostaje %d narzędzi i %d bajtów wykazu tools/list,"+
			" czyli tyle, ile dostawał bez eksperta", wynik.Powod, pomiar.Pozycji, pomiar.Bajtow)
	}
	if len(wynik.Nierozpoznane) > 0 && d.raz("nierozpoznane") {
		d.dziennik.Printf("dobór eksperta %q: %d kodów nierozpoznanych (nie nazywają ani narzędzia,"+
			" ani grupy) — %v", d.kod, len(wynik.Nierozpoznane), wynik.Nierozpoznane)
	}
}

// raz mówi, czy ten meldunek jeszcze nie padł. Wołane pod zamkiem `wykaz`.
func (d *doborEksperta) raz(klucz string) bool {
	if d.zameldowane == nil {
		d.zameldowane = map[string]bool{}
	}
	if d.zameldowane[klucz] {
		return false
	}
	d.zameldowane[klucz] = true
	return true
}
