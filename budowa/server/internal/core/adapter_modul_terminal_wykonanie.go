// Komenda `terminal.command.exec` — uruchomienie jednego polecenia w karcie
// terminala. Plik prowadzi trzy bramy: kartę, uprawnienie okna i izolację;
// sam bieg procesu leży w `adapter_modul_terminal_bieg.go`.
//
// Proces przeżywa rozłączenie klienta: gniazdo WebSocket może paść w połowie
// kompilacji, a kompilacja ma dobiec końca. Dlatego obserwator zakończenia
// pracuje we własnej gorutynie i własnym kontekście, a nie w kontekście
// komendy.
package core

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// czasNaDomkniecie jest chwilą, przez którą komenda kończąca proces czeka na
// obserwatora zakończenia. Wartość jest krótka, bo odpowiedź ma wrócić szybko,
// a stan końcowy dochodzi zdarzeniem `terminal.process.changed`.
const czasNaDomkniecie = 750 * time.Millisecond

// granicaCzasuDomyslna obowiązuje, gdy żądanie nie poda własnej. Zero znaczy
// brak granicy, żeby polecenie długie (budowanie, instalacja zależności) nie
// zostało przerwane wartością domyślną.
const granicaCzasuDomyslna = 0

// WykonajPolecenie obsługuje `terminal.command.exec`.
func (a *adapterTerminala) WykonajPolecenie(ctx context.Context,
	z shared.TerminalCommandExecRequest) (shared.TerminalCommandExecResponse, error) {

	karta, err := a.kartaZadania(z.SessionId)
	if err != nil {
		return shared.TerminalCommandExecResponse{}, err
	}
	tresc := strings.TrimSpace(z.Command)
	if tresc == "" {
		return shared.TerminalCommandExecResponse{}, bladZadaniaTerminala("polecenie puste nie ma czego wykonać")
	}
	okno, err := a.oknoWykonania(karta.oknoKod)
	if err != nil {
		return shared.TerminalCommandExecResponse{}, err
	}
	inicjator := inicjatorZadania(z.Initiator)
	if err := sprawdzUprawnienie(okno.TrybUprawnien, inicjator); err != nil {
		return shared.TerminalCommandExecResponse{}, err
	}
	polecenie, err := a.polecenieDopuszczone(karta, okno, tresc)
	if err != nil {
		return shared.TerminalCommandExecResponse{}, err
	}

	proces, err := a.uruchom(ctx, karta, okno, polecenie, tresc, inicjator)
	if err != nil {
		return shared.TerminalCommandExecResponse{}, err
	}
	a.pilnujZakonczenia(proces, granicaCzasu(z.TimeoutMs))
	return shared.TerminalCommandExecResponse{Process: procesKontraktu(proces)}, nil
}

// kartaZadania odnajduje kartę wskazaną w żądaniu.
func (a *adapterTerminala) kartaZadania(idKarty string) (*kartaTerminala, error) {
	kod := strings.TrimSpace(idKarty)
	if kod == "" {
		return nil, bladZadaniaTerminala("wykonanie polecenia wymaga wskazania karty")
	}
	karta, jest := a.rejestr.Karta(kod)
	if !jest {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Terminal: karta "+kod+" nie jest otwarta"))
	}
	return karta, nil
}

// polecenieDopuszczone składa polecenie powłoki i przepuszcza je przez
// egzekutor izolacji. Naruszenie izolacji wraca jako odmowa uprawnienia, a nie
// błąd wewnętrzny, i nazywa punkt izolacji, który zatrzymał uruchomienie.
func (a *adapterTerminala) polecenieDopuszczone(karta *kartaTerminala, okno session.Okno,
	tresc string) (session.Polecenie, error) {

	polecenie, err := polecenieKarty(karta, tresc)
	if err != nil {
		return session.Polecenie{}, bladZadaniaTerminala(err.Error())
	}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{Okno: okno.Id})
	}
	polecenie, err = session.SprawdzPolecenie(zasady, a.obszarOkna(okno), polecenie)
	if err != nil {
		if errors.Is(err, session.ErrIzolacja) {
			return session.Polecenie{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodePermissionDenied, "moduł Terminal: "+err.Error()))
		}
		return session.Polecenie{}, err
	}
	return rozwinSrodowisko(polecenie), nil
}

// rozwinSrodowisko dokłada środowisko rdzenia przed wpisami własnymi polecenia,
// zgodnie z jego polem `DziedziczSrodowisko`. Bez tego powłoka z własnym
// zestawem zmiennych straciłaby PATH i nie znalazłaby żadnego narzędzia.
func rozwinSrodowisko(p session.Polecenie) session.Polecenie {
	if !p.DziedziczSrodowisko || len(p.Srodowisko) == 0 {
		return p
	}
	p.Srodowisko = append(os.Environ(), p.Srodowisko...)
	return p
}
