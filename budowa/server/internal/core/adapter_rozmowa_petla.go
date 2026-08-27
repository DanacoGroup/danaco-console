package core

import (
	"context"
	"fmt"
	"strings"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Strona rozmowy pętli koordynator–wykonawca.
//
// Pakiet sesji rozstrzyga, kiedy zacząć obieg — pilnuje wybudzeń, licznika
// obiegów i warunku zatrzymania. Ten plik rozstrzyga, jak go zacząć: obieg
// wchodzi do okna koordynatora zwykłą wypowiedzią i rusza tą samą drogą, co
// wiadomość Operatora. Drugiego silnika tury w rdzeniu nie ma.

// RozpocznijObieg podejmuje kolejny obieg koordynatora. Strumień wykonawcy
// wchodzi do okna koordynatora jako wypowiedź, a tura rusza zwykłą drogą Wyslij.
func (a *adapterRozmowy) RozpocznijObieg(o session.Obieg) error {
	_, err := a.Wyslij(a.zycie, shared.MessageSendRequest{
		WindowId: o.IdKoordynatora,
		Content:  zlecenieObiegu(o),
	})
	return err
}

// zlecenieObiegu składa wypowiedź dla koordynatora ze strumienia wykonawcy.
// Starsza część strumienia jedzie liczbami, nie treścią.
func zlecenieObiegu(o session.Obieg) string {
	var tresc strings.Builder
	fmt.Fprintf(&tresc, "Obieg %d — okno wykonawcy %s zakończyło turę (%s).\n",
		o.Numer, o.IdWykonawcy, o.PowodTury)
	if o.Strumien.Zwiniete.Wpisow > 0 {
		fmt.Fprintf(&tresc, "Starsza część strumienia zwinięta: %d wpisów, %d znaków.\n",
			o.Strumien.Zwiniete.Wpisow, o.Strumien.Zwiniete.Znakow)
	}
	for _, wpis := range o.Strumien.Jawne {
		fmt.Fprintf(&tresc, "[%s] %s\n", wpis.Rodzaj, wpis.Tresc)
	}
	return tresc.String()
}

// powodTury nazywa przyczynę zamknięcia tury przekazywaną koordynatorowi.
//
// Warunki i ich kolejność są te same, co w `stanOdpowiedziZeZdarzen`
// (adapter_zdarzenia_zaczepow.go), i tak ma być: pętla poznaje wynik pracy po
// powodzie tury, Operator po stanie wiadomości, a oba wychodzą z jednego
// zamknięcia. Zdarzenie `result` z `is_error` liczy się także wtedy, gdy kanał
// dowiózł turę bez błędu — inaczej tura zamknięta błędem przy sprawnym kanale
// szłaby do pętli jako wynik i bieg ogłaszałby ukończenie, choć wiadomość ma
// stan `error`.
func powodTury(kontekst context.Context, err error, zamkniecie *zamkniecieTury) string {
	switch {
	case kontekst.Err() != nil:
		return "zatrzymanie tury"
	case err != nil:
		return "błąd tury"
	case zamkniecie != nil && zamkniecie.Blad:
		return "błąd tury"
	default:
		return session.PowodWynik
	}
}
