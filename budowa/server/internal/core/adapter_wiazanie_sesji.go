package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// adapterWiazaniaSesji wypełnia port WiazanieSesji: ognisko karty sesji i powiązanie
// połączenia z sesją trwającą na rdzeniu. Obie czynności należą do klienta, nie do konta.
type adapterWiazaniaSesji struct {
	zestaw   *dane.Zestaw
	nadzorca *session.Nadzorca
	klienci  *wieziKlientow
}

// nowyAdapterWiazaniaSesji wiąże port z repozytoriami, nadzorcą sesji i rejestrem więzi
// klientów urządzeń.
func nowyAdapterWiazaniaSesji(zestaw *dane.Zestaw, nadzorca *session.Nadzorca,
	klienci *wieziKlientow) *adapterWiazaniaSesji {
	return &adapterWiazaniaSesji{zestaw: zestaw, nadzorca: nadzorca, klienci: klienci}
}

// Ogniskuj przenosi ognisko klienta na wskazaną kartę sesji i, gdy wskazano, na okno w jej
// wnętrzu. Istnienia karty nie warunkuje: ognisko jest zapisem tego, co klient ma na wierzchu.
func (a *adapterWiazaniaSesji) Ogniskuj(_ context.Context, z shared.SessionFocusRequest) (shared.SessionFocusResponse, error) {
	zmiana := a.klienci.Ogniskuj(z.ClientId, z.SessionId, z.WindowId)
	wynik := shared.SessionFocusResponse{SessionId: zmiana.IdSesji, FocusedAt: zmiana.Chwila.UnixMilli()}
	if zmiana.IdOkna != "" {
		idOkna := zmiana.IdOkna
		wynik.WindowId = &idOkna
	}
	if zmiana.Poprzednia != "" {
		poprzednia := zmiana.Poprzednia
		wynik.PreviousSessionId = &poprzednia
	}
	return wynik, nil
}

// Powiaz wiąże połączenie z sesją i odtwarza jej okna. Pole resumed odpowiada na pytanie,
// czy sesja trwała na rdzeniu mimo rozłączenia klienta.
func (a *adapterWiazaniaSesji) Powiaz(ctx context.Context, z shared.SessionBindRequest) (shared.SessionBindResponse, error) {
	wynik := shared.SessionBindResponse{Windows: []shared.Window{}}
	if sesja, okna, jest := a.sesjaZywa(z.SessionId); jest {
		wynik.Session, wynik.Bound, wynik.Resumed = sesja, true, true
		return a.zapamietaj(z.ClientId, wynik, okna, z.WindowIds), nil
	}
	sesja, okna, jest, err := a.sesjaUtrwalona(ctx, z.SessionId)
	if err != nil || !jest {
		return wynik, err
	}
	wynik.Session, wynik.Bound = sesja, true
	return a.zapamietaj(z.ClientId, wynik, okna, z.WindowIds), nil
}

// sesjaZywa zwraca sesję rejestru nadzorcy wraz z jej oknami — tę, której
// procesy przetrwały rozłączenie klienta.
func (a *adapterWiazaniaSesji) sesjaZywa(idSesji string) (shared.Session, []shared.Window, bool) {
	if a.nadzorca == nil {
		return shared.Session{}, nil, false
	}
	sesja, err := a.nadzorca.Rejestr().Sesja(idSesji)
	if err != nil {
		return shared.Session{}, nil, false
	}
	okna, err := a.nadzorca.Rejestr().OknaSesji(idSesji)
	if err != nil {
		return shared.Session{}, nil, false
	}
	return sesjaKontraktu(sesja), oknaKontraktu(okna), true
}

// sesjaUtrwalona zwraca sesję z wierszy wraz z jej oknami — drogę po restarcie
// rdzenia, gdy rejestr nadzorcy jest już pusty.
func (a *adapterWiazaniaSesji) sesjaUtrwalona(ctx context.Context, idSesji string) (shared.Session, []shared.Window, bool, error) {
	wiersz, err := a.zestaw.Sesje.PoIdentyfikatorze(ctx, idSesji)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.Session{}, nil, false, nil
	}
	if err != nil {
		return shared.Session{}, nil, false, err
	}
	sesja := sesjaWierszaKontraktu(wiersz)
	okna, err := oknaSesjiUtrwalone(ctx, a.zestaw, wiersz.ID, sesja.Id)
	if err != nil {
		return shared.Session{}, nil, false, err
	}
	return sesja, okna, true, nil
}

// zapamietaj domyka powiązanie: zawęża okna do wskazanych, uzupełnia relację jeden do wielu
// sesji i zapisuje więź klienta z sesją.
func (a *adapterWiazaniaSesji) zapamietaj(idKlienta string, wynik shared.SessionBindResponse,
	okna []shared.Window, wskazane []string) shared.SessionBindResponse {
	wynik.Session.WindowIds = identyfikatoryOkien(okna)
	wynik.Windows = wybraneOkna(okna, wskazane)
	a.klienci.Powiaz(idKlienta, wynik.Session.Id)
	return wynik
}
