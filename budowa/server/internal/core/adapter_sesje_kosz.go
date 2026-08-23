// Odpowiedzialność pliku: powrót sesji z kosza w torze komendy
// `session.restore`.
//
// Kontrakt mówi o tej komendzie „przywraca sesję do historii bieżącej" —
// i dokładnie to robi powrót z kosza; osobnej komendy kosz nie dostaje.
// Sesja spoza kosza przechodzi tędy bez śladu i wraca z archiwum zwykłą
// drogą stanu.
//
// Powrót ma dwie części, bo usunięcie miało dwie: `session.delete` zdjęło
// znacznikiem wiersz z wykazów oraz wyprowadziło sesję z rejestru żywego
// (zwolnijZasobySesji). Przywrócenie czyści znacznik i wnosi sesję z powrotem
// do rejestru — bez tego `session.list`, czytający rejestr, dalej by jej
// nie widział, a „przywrócona" byłaby słowem bez skutku.
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

// wniesPrzywroconeDoRejestru odtwarza w rejestrze żywym sesje wyjęte z kosza.
// Idzie PO przestawieniu stanu w bazie — wiersz wraca do rejestru już jako
// czynny, a nie w stanie sprzed usunięcia. Niepowodzenie jednej sesji nie
// wstrzymuje pozostałych; ślad zostaje w dzienniku rdzenia.
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
