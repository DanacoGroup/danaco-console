package core

import (
	"context"
	"log"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// utrwalaczStanow zapisuje w bazie zmianę stanu sesji i okna komunikacji, domykając rozjazd między pamięcią procesu a trwałością w jednym miejscu. Wiersz odnajduje po identyfikatorze rdzenia; brak wiersza nie jest usterką.
type utrwalaczStanow struct {
	zycie context.Context
	sesje dane.RepozytoriumSesji
	okna  dane.RepozytoriumOkien
	// kosz to odwrotna strona sesji: przyjmuje wiersz po session.delete i zwraca po session.restore.
	kosz     dane.RepozytoriumKoszaSesji
	dziennik *log.Logger
}

// nowyUtrwalaczStanow wiąże utrwalacz z repozytoriami sesji i okien. Kontekst
// jest kontekstem rdzenia, nie połączenia: zapis stanu ma dobiec końca także
// wtedy, gdy klient rozłączy się zaraz po komendzie.
func nowyUtrwalaczStanow(zycie context.Context, repozytoria *dane.Zestaw,
	dziennik *log.Logger) *utrwalaczStanow {

	if repozytoria == nil {
		return nil
	}
	return &utrwalaczStanow{
		zycie: zycie, sesje: repozytoria.Sesje, okna: repozytoria.Okna,
		kosz: repozytoria.KoszSesji, dziennik: dziennik,
	}
}

// StanOkna zapisuje stan okna komunikacji w wierszu odnalezionym po identyfikatorze rdzenia tego okna.
func (u *utrwalaczStanow) StanOkna(idOkna string, stan shared.WindowStatus) {
	if u == nil || u.okna == nil || idOkna == "" {
		return
	}
	wiersz, err := u.okna.PoIdentyfikatorze(u.kontekst(), idOkna)
	if err != nil {
		u.odnotuj("stan okna %s: %v", idOkna, err)
		return
	}
	if err := u.okna.ZmienStan(u.kontekst(), wiersz.ID, stan); err != nil {
		u.odnotuj("zapis stanu okna %s: %v", idOkna, err)
	}
}

// StanSesji zapisuje stan sesji w wierszu odnalezionym po identyfikatorze rdzenia dla tej właśnie sesji.
func (u *utrwalaczStanow) StanSesji(idSesji string, stan shared.SessionStatus) {
	if u == nil || u.sesje == nil || idSesji == "" {
		return
	}
	wiersz, err := u.sesje.PoIdentyfikatorze(u.kontekst(), idSesji)
	if err != nil {
		u.odnotuj("stan sesji %s: %v", idSesji, err)
		return
	}
	if err := u.sesje.ZmienStan(u.kontekst(), wiersz.ID, stan); err != nil {
		u.odnotuj("zapis stanu sesji %s: %v", idSesji, err)
	}
}

// StanOkien zapisuje jeden stan dla wielu okien naraz, bo zamknięcie sesji zamyka od razu wszystkie jej okna.
func (u *utrwalaczStanow) StanOkien(idOkien []string, stan shared.WindowStatus) {
	for _, idOkna := range idOkien {
		u.StanOkna(idOkna, stan)
	}
}

// kontekst zwraca kontekst zapisu. Brak kontekstu życia schodzi na tło —
// utrwalacz ma zapisać stan, a nie odmówić z powodu braku nastawy.
func (u *utrwalaczStanow) kontekst() context.Context {
	if u.zycie == nil {
		return context.Background()
	}
	return u.zycie
}

// odnotuj zapisuje niepowodzenie utrwalenia w dzienniku rdzenia. Niepowodzenie
// zapisu nie unieważnia czynności: okno jest zamknięte w rejestrze niezależnie
// od tego, czy baza przyjęła jego stan.
func (u *utrwalaczStanow) odnotuj(wzor string, argumenty ...any) {
	if u.dziennik == nil {
		return
	}
	u.dziennik.Printf(wzor, argumenty...)
}

// Utrwalacz nie kasuje wierszy sesji: session.delete przenosi wiersz do kosza sesji.
