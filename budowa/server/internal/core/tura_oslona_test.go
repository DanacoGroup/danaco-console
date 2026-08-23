package core

import (
	"context"
	"testing"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// portDolozenPanikujacy udaje adapter, który wywraca się w torze tury.
//
// Atrapa jest tu jedyną drogą pomiaru: usterkę adaptera trzeba wywołać
// umyślnie, a nie czekać, aż wróci ta sama, którą właśnie naprawiono. Osłona ma
// trzymać KAŻDĄ usterkę toru tury, nie tę jedną.
type portDolozenPanikujacy struct{}

func (portDolozenPanikujacy) DolozeniaNarzedzi(context.Context, string) ([]string, error) {
	panic("umyślna usterka adaptera w torze tury")
}

// TestTuraNieGasiRdzeniaPrzyUsterceAdaptera dowodzi, że panika w gorutynie tury
// kończy turę, a nie proces.
//
// Bez osłony panika w gorutynie nie ma kto przechwycić: ginie cały rdzeń wraz
// z sesją, kolejką i połączeniem Operatora. Miarą naprawy jest to, że Operator
// traci jedną odpowiedź zamiast całej pracy — dlatego sprawdzian mierzy PRZEŻYCIE
// procesu, a nie treść odmowy.
//
// Sprawdzian wypada niepomyślnie przez padnięcie całego przebiegu, gdy osłona
// zniknie — i tak ma być: to jest dokładnie ta szkoda, przed którą stoi.
func TestTuraNieGasiRdzeniaPrzyUsterceAdaptera(t *testing.T) {
	rozmowa := &adapterRozmowy{}
	rozmowa.ZDolozeniamiSesji(portDolozenPanikujacy{})

	domkniete := make(chan struct{})
	go func() {
		defer close(domkniete)
		// Tor tury wywołany wprost: `Wyslij` powołuje tę samą gorutynę, lecz
		// wymaga kanału modelu i wpisu rozmowy. Badana jest osłona, a nie droga
		// dojścia do niej.
		okno := session.Okno{Id: "okn-oslona", IdSesji: "ses-oslona"}
		pytanie := shared.Message{Id: "msg-pyt", WindowId: okno.Id, SessionId: okno.IdSesji}
		odpowiedz := shared.Message{Id: "msg-odp", WindowId: okno.Id, SessionId: okno.IdSesji}
		rozmowa.prowadzTure(context.Background(), okno, pytanie, odpowiedz, "zad-oslona")
	}()

	select {
	case <-domkniete:
		// Gorutyna wróciła — panika została przechwycona.
	case <-time.After(5 * time.Second):
		t.Fatal("tor tury nie domknął się w czasie — osłona albo zawiesza, albo nie wraca")
	}
}
