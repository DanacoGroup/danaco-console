package uruchomienie

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
)

// zakonczenieWiersza oddziela żądania na wejściu i odpowiedzi na wyjściu toru
// wykonawczego. Jeden wiersz to jedna koperta kontraktu.
const zakonczenieWiersza = '\n'

// odczyt niesie jeden wiersz wejścia strumienia albo przyczynę, dla której
// odczyt tego wiersza się przerwał.
type odczyt struct {
	dane []byte
	err  error
}

// TorWykonawczy przyjmuje żądania kontraktu wierszami strumienia wejścia
// i odsyła odpowiedzi wierszami strumienia wyjścia; wraca po wyczerpaniu
// wejścia albo po zamknięciu kontekstu.
func TorWykonawczy(kontekst context.Context, rdzen Rdzen, o Otoczenie) error {
	if rdzen == nil {
		<-kontekst.Done()
		return nil
	}
	strumien := czytajWiersze(kontekst, o.wejscie())
	wyjscie := o.wyjscie()
	for {
		select {
		case <-kontekst.Done():
			return nil
		case przyjety, otwarty := <-strumien:
			if !otwarty {
				return nil
			}
			if przyjety.err != nil {
				return fmt.Errorf("tor wykonawczy: odczyt żądania: %w", przyjety.err)
			}
			if err := obsluzWiersz(kontekst, rdzen, przyjety.dane, wyjscie, o); err != nil {
				return err
			}
		}
	}
}

// obsluzWiersz wykonuje jedno żądanie i odsyła odpowiedź. Wiersz pusty jest
// pomijany, a odpowiedź niemożliwa do zakodowania odnotowana w dzienniku —
// jedno wadliwe żądanie nie zabiera ze sobą całego toru.
func obsluzWiersz(kontekst context.Context, rdzen Rdzen, wiersz []byte, wyjscie io.Writer, o Otoczenie) error {
	wiersz = bytes.TrimSpace(wiersz)
	if len(wiersz) == 0 {
		return nil
	}
	odpowiedz := rdzen.WykonajSurowe(kontekst, wiersz)
	if len(odpowiedz) == 0 {
		o.zapisz("tor wykonawczy: żądanie bez odpowiedzi do odesłania")
		return nil
	}
	if _, err := wyjscie.Write(append(odpowiedz, zakonczenieWiersza)); err != nil {
		return fmt.Errorf("tor wykonawczy: zapis odpowiedzi: %w", err)
	}
	return nil
}

// czytajWiersze prowadzi odczyt w osobnym wątku, żeby zamknięcie kontekstu
// zatrzymywało tor także wtedy, gdy strumień wejścia milczy.
func czytajWiersze(kontekst context.Context, wejscie io.Reader) <-chan odczyt {
	strumien := make(chan odczyt)
	go func() {
		defer close(strumien)
		czytnik := bufio.NewReader(wejscie)
		for {
			wiersz, err := czytnik.ReadBytes(zakonczenieWiersza)
			if len(wiersz) > 0 && !podaj(kontekst, strumien, odczyt{dane: wiersz}) {
				return
			}
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				podaj(kontekst, strumien, odczyt{err: err})
				return
			}
		}
	}()
	return strumien
}

// podaj przekazuje odczyt dalej albo kończy pracę czytnika, gdy odbiorca
// zniknął wraz z zamknięciem kontekstu.
func podaj(kontekst context.Context, strumien chan<- odczyt, co odczyt) bool {
	select {
	case strumien <- co:
		return true
	case <-kontekst.Done():
		return false
	}
}
