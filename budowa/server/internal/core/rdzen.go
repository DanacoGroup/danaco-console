package core

import (
	"context"
	"fmt"
	"log"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// WersjaRdzenia jest wersją produktu zgłaszaną klientowi w powitaniu.
// Do pierwszej publikacji obowiązuje v1.0, bez wersjonowania w trakcie budowy.
const WersjaRdzenia = "1.0"

// Rdzen kieruje komendę do obsługiwacza i zwraca kopertę odpowiedzi.
//
// Jest jedynym miejscem, przez które przechodzi każdy komunikat przychodzący.
// Nie zna transportu ani domen — zna rejestr nazw i porty.
type Rdzen struct {
	rejestr  *Rejestr
	komendy  *protocol.RejestrKomend
	nasluch  Nasluch
	dziennik *log.Logger
	// niepowodzenia przyjmuje odmowy wykonania komend, żeby stały się widoczne w module Diagnostics.
	niepowodzenia ObserwatorNiepowodzen
	// straz rozstrzyga, czy wywołanie ręki modelu mieści się w zakresie zapisanym dla profilu asystenta.
	straz StrazZakresowNarzedzi
	// wiez trzyma przypisania połączenie -> sesja bramki i nigdy nie jest zerowa po założeniu.
	wiez *wiezBramki
}

// ZObserwatoremNiepowodzen wpina odbiorcę odmów wykonania komend, dzięki czemu odmowy stają się widoczne poza rdzeniem, w module Diagnostics.
func (r *Rdzen) ZObserwatoremNiepowodzen(o ObserwatorNiepowodzen) *Rdzen {
	r.niepowodzenia = o
	return r
}

// ZeStrazaZakresow wpina straż zakresów narzędzi, która rozstrzyga, czy wywołanie ręki modelu mieści się w zakresie zapisanym dla profilu asystenta.
func (r *Rdzen) ZeStrazaZakresow(s StrazZakresowNarzedzi) *Rdzen {
	r.straz = s
	return r
}

// Wykonaj kieruje jedno żądanie i zwraca kopertę zwrotną.
//
// Zwraca zawsze. Komenda spoza rejestru dostaje `*.unknown`, błąd wykonania
// wraca kodem kontraktu, a usterka obsługiwacza kończy wyłącznie to jedno
// wywołanie.
func (r *Rdzen) Wykonaj(ctx context.Context, zadanie protocol.Koperta) protocol.Koperta {
	return r.WykonajZadanie(ctx, protocol.ZbudujZadanie(zadanie, r.komendy))
}

// WykonajZadanie kieruje żądanie już rozpoznane przez warstwę niższą. Komenda bez obsługiwacza jest rozpoznawana ponownie, żeby odpowiedź „*.unknown" wróciła pod nazwą zdarzenia obszaru, a obie drogi odmowy trafiają jednakowo do Errors Panel.
func (r *Rdzen) WykonajZadanie(ctx context.Context, z protocol.Request) protocol.Koperta {
	ctx = zDziennikiemRdzenia(ctx, r.dziennik)
	if !z.Znana {
		return r.odmowaNieznanej(ctx, z)
	}
	obsluga, jest := r.rejestr.Obsluga(z.Komenda)
	if !jest {
		return r.odmowaNieznanej(ctx, protocol.ZbudujZadanie(z.Koperta(), r.komendy))
	}
	if blad := r.odmowaZakresu(ctx, z); blad != nil {
		r.odnotujNiepowodzenie(ctx, z, blad)
		return protocol.KopertaOdpowiedzi(z.Koperta(), protocol.Odpowiedz{Blad: blad})
	}
	odpowiedz := r.wykonajOdpornie(ctx, obsluga, z)
	// Odmowa staje się faktem po wykonaniu, nie zamiast niego, niezależnie od wyniku zapisu.
	r.odnotujNiepowodzenie(ctx, z, odpowiedz.Blad)
	return protocol.KopertaOdpowiedzi(z.Koperta(), odpowiedz)
}

// WykonajSurowe przyjmuje bajty z gniazda i oddaje bajty do odesłania — jest
// wejściem warstwy transportu. Komunikat nieczytelny nie zrywa połączenia:
// wraca odpowiedź z kodem błędu kontraktu.
func (r *Rdzen) WykonajSurowe(ctx context.Context, dane []byte) []byte {
	odpowiedz := r.kopertaOdpowiedzi(ctx, dane)
	bajty, err := protocol.Zakoduj(odpowiedz)
	if err != nil {
		r.zapisz("kodowanie odpowiedzi %s: %v", odpowiedz.Type, err)
		return nil
	}
	return bajty
}

// Uruchom oddaje sterowanie warstwie nasłuchu i wraca po zamknięciu kontekstu.
// Rdzeń bez podłączonego nasłuchu pracuje dalej — czeka na zatrzymanie zamiast
// przerywać start.
func (r *Rdzen) Uruchom(ctx context.Context) error {
	r.zapisz("rdzeń gotowy: komend=%d", r.rejestr.Liczba())
	r.zglosZaleznosci()
	if r.nasluch == nil {
		<-ctx.Done()
		return nil
	}
	return r.nasluch.Sluchaj(ctx)
}

// kopertaOdpowiedzi odkodowuje komunikat przychodzący z gniazda, a komunikat nieczytelny zwraca jako odpowiedź z kodem błędu.
func (r *Rdzen) kopertaOdpowiedzi(ctx context.Context, dane []byte) protocol.Koperta {
	zadanie, err := protocol.Odkoduj(dane)
	if err != nil {
		r.zapisz("komunikat nieczytelny: %v", err)
		return kopertaNieczytelna(err)
	}
	return r.Wykonaj(ctx, zadanie)
}

// wykonajOdpornie uruchamia obsługiwacza tak, by jego usterka nie zabrała ze
// sobą połączenia ani procesu. Błąd techniczny jest błędem bieżącego
// wywołania — nie blokuje sesji, konta ani kolejnych prób.
func (r *Rdzen) wykonajOdpornie(ctx context.Context, obsluga Obsluga, z protocol.Request) (odpowiedz protocol.Odpowiedz) {
	defer func() {
		if przyczyna := recover(); przyczyna != nil {
			r.zapisz("obsługiwacz %s przerwał wykonanie: %v", z.Komenda, przyczyna)
			odpowiedz = protocol.PorazkaKodem(shared.ErrorCodeInternalError,
				fmt.Sprintf("rdzeń: obsługa %s przerwana", z.Komenda))
		}
	}()
	return obsluga(ctx, z)
}

// zapisz odnotowuje zdarzenie techniczne rdzenia. Dziennik jest opcjonalny —
// jego brak nie zmienia zachowania.
func (r *Rdzen) zapisz(wzorzec string, argumenty ...any) {
	if r == nil || r.dziennik == nil {
		return
	}
	r.dziennik.Printf(wzorzec, argumenty...)
}

// kluczDziennikaRdzenia znakuje dziennik włożony do kontekstu żądania, żeby warstwy niższe mogły go stamtąd odczytać.
type kluczDziennikaRdzenia struct{}

// zDziennikiemRdzenia niesie dziennik do warstw, które stoją na drodze żądania, a rdzenia nie widzą, na przykład do bramy kontraktu, która zapisuje ślad przepuszczonego niepełnego powitania.
func zDziennikiemRdzenia(ctx context.Context, dziennik *log.Logger) context.Context {
	if dziennik == nil {
		return ctx
	}
	return context.WithValue(ctx, kluczDziennikaRdzenia{}, dziennik)
}

// dziennikZKontekstu oddaje dziennik rdzenia albo nic. Nic jest odpowiedzią
// prawidłową: rdzeń złożony bez dziennika pracuje tak samo, tylko milcząco.
func dziennikZKontekstu(ctx context.Context) *log.Logger {
	if ctx == nil {
		return nil
	}
	dziennik, jest := ctx.Value(kluczDziennikaRdzenia{}).(*log.Logger)
	if !jest {
		return nil
	}
	return dziennik
}

// odmowaZakresu pyta straż wyłącznie o rękę modelu. Zakres profilu asystenta nie obejmuje klawiatury Operatora ani pracy własnej rdzenia, więc te drogi przechodzą bez zawężenia.
func (r *Rdzen) odmowaZakresu(ctx context.Context, z protocol.Request) *protocol.Blad {
	if r == nil || r.straz == nil {
		return nil
	}
	rodzaj, _ := sprawca(ctx)
	if rodzaj == nil || *rodzaj != shared.ActorKindModel {
		return nil
	}
	if err := r.straz.SprawdzWywolanie(ctx, z.Komenda, z.Zasieg.Sesja); err != nil {
		blad := protocol.BladZeZrodla(shared.ErrorCodePermissionDenied, err)
		return &blad
	}
	return nil
}
