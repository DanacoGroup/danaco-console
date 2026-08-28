// Plik przygotowuje moduł Terminal przy montażu rdzenia: osierocenie procesów zostawionych przez
// poprzedni bieg i odtworzenie kart powłok z dziennika.
package core

import (
	"context"
	"log"
)

// przygotujTerminal odtwarza stan modułu Terminal z dziennika rdzenia przy montażu tej platformy konta.
func przygotujTerminal(kontekst context.Context, terminal *adapterTerminala, dziennik *log.Logger) {
	if terminal == nil {
		return
	}
	if err := terminal.Przygotuj(kontekst); err != nil && dziennik != nil {
		dziennik.Printf("moduł Terminal: nie można odtworzyć kart i dziennika procesów: %v", err)
	}
}
