package core

import (
	"context"
	"log"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// utrwalaczStanow zapisuje w bazie zmianę stanu sesji i okna komunikacji.
//
// Zamknięcie okna, zamknięcie sesji i jej usunięcie to czynności rejestru
// nadzorcy — pakiet sesji jest ich właścicielem i o bazie nie wie. Bez tego
// utrwalacza stan `zamkniete`/`zakonczona` żyłby wyłącznie w pamięci procesu:
// historia wiadomości byłaby trwała, a sesja po restarcie wracałaby jako
// czynna. Utrwalacz domyka ten rozjazd w jednym miejscu, zamiast powtarzać
// zapis w każdym adapterze z osobna.
//
// Wiersz odnajduje po identyfikatorze rdzenia (`identyfikator_zewnetrzny`).
// Brak wiersza nie jest usterką: sesja bez ani jednej utrwalonej wiadomości
// nie ma jeszcze wiersza, a zamknięcie i tak ma się odbyć.
type utrwalaczStanow struct {
	zycie context.Context
	sesje dane.RepozytoriumSesji
	okna  dane.RepozytoriumOkien
	// kosz to odwrotna strona tabeli sesji: tam trafia wiersz po
	// `session.delete` i stamtąd wraca po `session.restore`. Czynności kosza
	// mieszkają w `trwalosc_kosza.go` — ten plik zna wyłącznie pole.
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

// StanOkna zapisuje stan okna komunikacji.
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

// StanSesji zapisuje stan sesji.
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

// StanOkien zapisuje jeden stan dla wielu okien — zamknięcie sesji zamyka
// wszystkie jej okna naraz.
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

// Utrwalacz nie kasuje wierszy sesji. `session.delete` przenosi wiersz do kosza
// (`trwalosc_kosza.go`), a fizyczny DELETE wykonuje wyłącznie czyszczenie
// startowe po terminie, które sprząta również bloki wiadomości bez klucza
// obcego.
//
// Katalog roboczy sesji zostaje nietknięty w obu fazach: pliki, które model
// zostawił, nie należą do bazy i nie znikają razem z wierszem.
