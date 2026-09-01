package core

import (
	"context"
	"testing"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// portDolozenPanikujacy udaje adapter, który wywraca się w torze tury. Atrapa jest jedyną drogą pomiaru: usterkę adaptera trzeba wywołać umyślnie, a osłona ma trzymać każdą usterkę toru tury, nie jedną konkretną.
type portDolozenPanikujacy struct{}

func (portDolozenPanikujacy) DolozeniaNarzedzi(context.Context, string) ([]string, error) {
	panic("umyślna usterka adaptera w torze tury")
}

// TestTuraNieGasiRdzeniaPrzyUsterceAdaptera dowodzi, że panika w gorutynie tury kończy turę, a nie proces. Bez osłony panika ginie razem z całym rdzeniem, sesją, kolejką i połączeniem. Sprawdzian mierzy przeżycie procesu, nie treść odmowy.
func TestTuraNieGasiRdzeniaPrzyUsterceAdaptera(t *testing.T) {
	rozmowa := &adapterRozmowy{}
	rozmowa.ZDolozeniamiSesji(portDolozenPanikujacy{})

	domkniete := make(chan struct{})
	go func() {
		defer close(domkniete)
		// Tor tury wywołany wprost: Wyslij powołuje tę samą gorutynę, lecz wymaga kanału modelu i wpisu.
		okno := session.Okno{Id: "okn-oslona", IdSesji: "ses-oslona"}
		pytanie := shared.Message{Id: "msg-pyt", WindowId: okno.Id, SessionId: okno.IdSesji}
		odpowiedz := shared.Message{Id: "msg-odp", WindowId: okno.Id, SessionId: okno.IdSesji}
		rozmowa.prowadzTure(context.Background(), okno, pytanie, odpowiedz, "zad-oslona", nil)
	}()

	select {
	case <-domkniete:
		// Gorutyna wróciła — panika została przechwycona.
	case <-time.After(5 * time.Second):
		t.Fatal("tor tury nie domknął się w czasie — osłona albo zawiesza, albo nie wraca")
	}
}
