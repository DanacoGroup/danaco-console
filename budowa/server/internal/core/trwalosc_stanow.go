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

// Kontekst jest kontekstem rdzenia, nie połączenia: zapis stanu ma dobiec końca
// także wtedy, gdy klient rozłączy się zaraz po komendzie.
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
func (u *utrwalaczStanow) StanOkna(ctx context.Context, idOkna string, stan shared.WindowStatus) {
	if u == nil || u.okna == nil || idOkna == "" {
		return
	}
	zapis := u.kontekst(ctx)
	wiersz, err := u.okna.PoIdentyfikatorze(zapis, idOkna)
	if err != nil {
		u.odnotuj("stan okna %s: %v", idOkna, err)
		return
	}
	if err := u.okna.ZmienStan(zapis, wiersz.ID, stan); err != nil {
		u.odnotuj("zapis stanu okna %s: %v", idOkna, err)
	}
}

// StanSesji zapisuje stan sesji w wierszu odnalezionym po identyfikatorze rdzenia dla tej właśnie sesji.
func (u *utrwalaczStanow) StanSesji(ctx context.Context, idSesji string, stan shared.SessionStatus) {
	if u == nil || u.sesje == nil || idSesji == "" {
		return
	}
	zapis := u.kontekst(ctx)
	wiersz, err := u.sesje.PoIdentyfikatorze(zapis, idSesji)
	if err != nil {
		u.odnotuj("stan sesji %s: %v", idSesji, err)
		return
	}
	if err := u.sesje.ZmienStan(zapis, wiersz.ID, stan); err != nil {
		u.odnotuj("zapis stanu sesji %s: %v", idSesji, err)
	}
}

// StanOkien zapisuje jeden stan dla wielu okien naraz, bo zamknięcie sesji zamyka od razu wszystkie jej okna.
func (u *utrwalaczStanow) StanOkien(ctx context.Context, idOkien []string, stan shared.WindowStatus) {
	for _, idOkna := range idOkien {
		u.StanOkna(ctx, idOkna, stan)
	}
}

// kontekst składa kontekst zapisu: życie rdzenia z kontem zamawiającego
// (decyzja 34). Brak kontekstu życia schodzi na tło.
func (u *utrwalaczStanow) kontekst(ctx context.Context) context.Context {
	zycie := u.zycie
	if zycie == nil {
		zycie = context.Background()
	}
	return dane.ZKontemOperatora(zycie, dane.KontoOperatora(ctx))
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
