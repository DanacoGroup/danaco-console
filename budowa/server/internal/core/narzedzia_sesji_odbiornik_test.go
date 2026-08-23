package core

import (
	"context"
	"testing"
)

// Straż odbiornika zerowego w porcie doraźnych dołożeń narzędzi.
//
// Adapter wchodzi do składacza zestawu tury przez interfejs, a wskaźnik zerowy
// schowany w interfejsie przechodzi porównanie `== nil` u wołającego — więc
// sprawdzenie po tamtej stronie go nie zatrzyma. Bez straży pierwsze
// `message.send` w takim montażu zabijało CAŁY proces rdzenia panicą w gorutynie
// tury, a Operator tracił sesję, kolejkę i połączenie naraz.
//
// Sprawdzian pilnuje zamiany paniki na odmowę, a nie samego montażu: montaż
// dziś wpina adapter poprawnie, ale to jest stan do popsucia jednym pominiętym
// ogniwem, i wtedy skutkiem ma być zdanie, nie zgaszony rdzeń.
func TestPortDolozenNieGasiRdzeniaPrzyOdbiornikuZerowym(t *testing.T) {
	var zerowy *adapterNarzedziSesji

	// Wejście przez interfejs, nie przez wskaźnik — dokładnie tą drogą woła
	// składacz zestawu tury. Wywołanie na wskaźniku pominęłoby badany przypadek.
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

// Ta sama straż na drugim wejściu portu: czynność Operatora, nie tura modelu.
func TestCzynnosciDolozenNieGasiaRdzeniaPrzyOdbiornikuZerowym(t *testing.T) {
	var zerowy *adapterNarzedziSesji

	if _, err := zerowy.sesjaZadania(context.Background(), "sesja-dowolna"); err == nil {
		t.Error("przekład sesji na odbiorniku zerowym nie odmówił")
	}
}
