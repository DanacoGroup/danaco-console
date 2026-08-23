// Odpowiedzialność pliku: przygotowanie modułu Terminal przy montażu rdzenia —
// osierocenie procesów zostawionych przez poprzedni bieg i odtworzenie kart
// powłok z dziennika.
//
// Krok jest osobną funkcją, nie wierszem w Zmontuj, bo ma własną regułę
// niepowodzenia: nieudane przygotowanie nie zatrzymuje startu rdzenia. Moduł
// rusza wtedy z pustym stanem żywym, a Operator dostaje wiadomość w dzienniku —
// terminal bez historii jest gorszy od terminala z historią, lecz nieuruchomiony
// serwer jest gorszy od obu.
package core

import (
	"context"
	"log"
)

// przygotujTerminal odtwarza stan modułu Terminal z dziennika.
func przygotujTerminal(kontekst context.Context, terminal *adapterTerminala, dziennik *log.Logger) {
	if terminal == nil {
		return
	}
	if err := terminal.Przygotuj(kontekst); err != nil && dziennik != nil {
		dziennik.Printf("moduł Terminal: nie można odtworzyć kart i dziennika procesów: %v", err)
	}
}
