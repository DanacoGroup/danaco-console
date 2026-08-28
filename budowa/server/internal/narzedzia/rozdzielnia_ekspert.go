// Plik prowadzi drogę zasięgu eksperta wewnątrz rozdzielni: odczyt definicji,
// złożenie podzbioru, pomiar i meldunek; bez trybu eksperta cała maszyneria
// nie rusza.
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
	// dolozenia są nazwami dołożonymi doraźnie w sesji, spoza definicji eksperta.
	dolozenia []string
	dziennik  *log.Logger
	zamek     sync.Mutex
	// zlozony jest zapamiętanym wynikiem udanego odczytu definicji; nil znaczy
	// jeszcze nie czytano.
	zlozony *WykazEksperta
	// zameldowane pilnuje, żeby ten sam powód braku zawężenia nie szedł do
	// dziennika powtórnie.
	zameldowane map[string]bool
}

// wykaz zwraca wykaz podawany modelowi w zasięgu eksperta; każda droga, udana
// i nieudana, kończy się meldunkiem z liczbami.
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

// wolno rozstrzyga, czy nazwa narzędzia mieści się w podzbiorze eksperta;
// przed pierwszym odczytem wolno wszystko.
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

// odmowaPozaDoborem buduje treść odmowy dla narzędzia spoza podzbioru eksperta,
// odróżnioną od odmowy granicą uprawnień modelu.
func (d *doborEksperta) odmowaPozaDoborem(nazwa string) error {
	return fmt.Errorf("narzędzie %q stoi poza doborem eksperta %q: okno pracuje w zasięgu eksperta,"+
		" a jego skillIds i connectorIds tego narzędzia ani jego grupy nie wskazują."+
		" Komplet narzędzi tego okna zwraca tools/list; rozszerza go Operator w Agent Builderze,"+
		" nie model w rozmowie", nazwa, d.kod)
}

// zamelduj kładzie na dziennik cenę zestawu i powód braku zawężenia; meldunek
// idzie zawsze przy pierwszym złożeniu.
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

// raz mówi, czy ten konkretny meldunek dziennika jeszcze w ogóle nie padł;
// wołane tylko pod zamkiem wykaz.
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
