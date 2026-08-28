// Plik zajmuje się powrotem sesji z kosza w torze komendy przywracania: sesja
// spoza kosza przechodzi tędy bez śladu i wraca z archiwum zwykłą drogą stanu.
package core

import (
	"context"
	"log"
)

// wyjmijZKosza czyści znaczniki kosza wskazanych sesji i zwraca te, które
// naprawdę w koszu leżały. Wskazanie spoza kosza ani wskazanie bez wiersza
// nie jest błędem — czynność zbiorcza nie pada przez jedną pozycję.
func (a *adapterSesji) wyjmijZKosza(ctx context.Context, wskazania []string) []string {
	wyjete := make([]string, 0, len(wskazania))
	for _, identyfikator := range wskazania {
		if a.trwalosc == nil {
			break
		}
		przywrocona, err := a.trwalosc.PrzywrocSesjeZKosza(ctx, identyfikator)
		if err != nil || !przywrocona {
			continue
		}
		wyjete = append(wyjete, identyfikator)
	}
	return wyjete
}

// wniesPrzywroconeDoRejestru odtwarza w rejestrze żywym sesje wyjęte z kosza,
// idąc po przestawieniu stanu w bazie; niepowodzenie jednej sesji nie
// wstrzymuje pozostałych.
func (a *adapterSesji) wniesPrzywroconeDoRejestru(ctx context.Context, identyfikatory []string) {
	for _, identyfikator := range identyfikatory {
		odtworzSesjePoIdentyfikatorze(ctx, a.zestaw, a.nadzorca, a.dziennikKosza(), identyfikator)
	}
}

// dziennikKosza podaje dziennik rdzenia dla śladu powrotu. Utrwalacz go ma;
// adapter własnego nie prowadzi.
func (a *adapterSesji) dziennikKosza() *log.Logger {
	if a.trwalosc == nil {
		return nil
	}
	return a.trwalosc.dziennik
}
