package core

import (
	"context"
	"testing"
)

// Plik sprawdza straż odbiornika zerowego w porcie doraźnych dołożeń narzędzi: adapter wchodzi
// przez interfejs, gdzie wskaźnik zerowy przechodzi porównanie z nil u wołającego.
func TestPortDolozenNieGasiRdzeniaPrzyOdbiornikuZerowym(t *testing.T) {
	var zerowy *adapterNarzedziSesji

	// Wejście przez interfejs, nie przez wskaźnik — dokładnie tą drogą woła składacz zestawu tury.
	var port DolozeniaNarzedziSesji = zerowy
	if port == nil {
		t.Fatal("wskaźnik zerowy w interfejsie porównał się do nil — " +
			"sprawdzian przestał badać przypadek, dla którego powstał")
	}

	nazwy, err := port.DolozeniaNarzedzi(context.Background(), "sesja-dowolna")
	if err == nil {
		t.Fatal("port z odbiornikiem zerowym nie odmówił — cisza tutaj znaczy, " +
			"że tura pojedzie bez dołożeń i nikt się o tym nie dowie")
	}
	if nazwy != nil {
		t.Errorf("odmowa oddała wykaz %v — odmowa nie ma prawa nieść wyniku", nazwy)
	}
}

// Ta sama straż działa na drugim wejściu portu: czynność Operatora, nie tura modelu w tej samej rozmowie.
func TestCzynnosciDolozenNieGasiaRdzeniaPrzyOdbiornikuZerowym(t *testing.T) {
	var zerowy *adapterNarzedziSesji

	if _, err := zerowy.sesjaZadania(context.Background(), "sesja-dowolna"); err == nil {
		t.Error("przekład sesji na odbiorniku zerowym nie odmówił")
	}
}
