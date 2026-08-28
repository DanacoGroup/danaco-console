// Plik obsługuje kosz sesji na utrwalaczu stanów: przeniesienie do kosza po session.delete, przywrócenie po session.restore i czyszczenie po terminie przy starcie rdzenia. Usunięcie sesji jest odwracalne, nie kasuje wiersza od ręki.
package core

import (
	"context"
	"log"
	"time"

	"danacoconsole/server/internal/dane"
)

// terminKoszaDni mówi, ile dni sesja leży w koszu, zanim czyszczenie startowe
// skasuje ją trwale wraz z całym zapisem.
const terminKoszaDni = 30

// formatChwiliKosza odpowiada wyrażeniu strftime schematu bazy dla kolumny usunieto_o wiersza kosza sesji.
const formatChwiliKosza = "2006-01-02T15:04:05.999Z"

// PrzeniesSesjeDoKosza stawia znacznik kosza na wierszu sesji. Brak wiersza
// wraca jako dane.ErrBrakWiersza — wywołujący rozstrzyga, czy to pominięcie
// (usuwanie zbiorcze), czy błąd; sesja bez wiersza nie ma czego tracić.
func (u *utrwalaczStanow) PrzeniesSesjeDoKosza(ctx context.Context, identyfikator string) error {
	if u == nil || u.sesje == nil || u.kosz == nil {
		return nil
	}
	wiersz, err := u.sesje.PoIdentyfikatorze(ctx, identyfikator)
	if err != nil {
		return err
	}
	if err := u.kosz.PrzeniesDoKosza(ctx, wiersz.ID); err != nil {
		return err
	}
	if u.dziennik != nil {
		u.dziennik.Printf("sesja %s przeniesiona do kosza (czyszczenie trwałe po %d dniach)",
			identyfikator, terminKoszaDni)
	}
	return nil
}

// PrzywrocSesjeZKosza czyści znacznik kosza. Fałsz mówi, że sesja w koszu nie
// leżała — przywracanie z archiwum idzie wtedy zwykłą drogą stanu.
func (u *utrwalaczStanow) PrzywrocSesjeZKosza(ctx context.Context, identyfikator string) (bool, error) {
	if u == nil || u.sesje == nil || u.kosz == nil {
		return false, nil
	}
	wiersz, err := u.sesje.PoIdentyfikatorze(ctx, identyfikator)
	if err != nil {
		return false, err
	}
	przywrocona, err := u.kosz.Przywroc(ctx, wiersz.ID)
	if err != nil {
		return false, err
	}
	if przywrocona && u.dziennik != nil {
		u.dziennik.Printf("sesja %s przywrócona z kosza wraz z całym zapisem", identyfikator)
	}
	return przywrocona, nil
}

// usunSesjePoTerminie kasuje trwale sesje leżące w koszu dłużej niż termin, wraz z całym ich zapisem danych.
func usunSesjePoTerminie(kontekst context.Context, repozytoria *dane.Zestaw, dziennik *log.Logger) {
	if repozytoria == nil || repozytoria.KoszSesji == nil {
		return
	}
	granica := time.Now().UTC().AddDate(0, 0, -terminKoszaDni).Format(formatChwiliKosza)
	usuniete, err := repozytoria.KoszSesji.UsunPrzeterminowane(kontekst, granica)
	if err != nil {
		if dziennik != nil {
			dziennik.Printf("czyszczenie kosza sesji: %v", err)
		}
		return
	}
	if usuniete > 0 && dziennik != nil {
		dziennik.Printf("kosz sesji: usunięto trwale %d sesji po terminie %d dni",
			usuniete, terminKoszaDni)
	}
}
