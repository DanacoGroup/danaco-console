package core

import (
	"sync"
	"time"
)

// wieziKlientow pamięta, co klient ma otwarte: kartę sesji w ognisku i okno
// ogniskowane wewnątrz tej karty.
//
// Ognisko żyje w rdzeniu, a nie w bazie. Jest właściwością klienta, nie
// konta (kontrakt komendy `session.focus`), i trwa tyle, co połączenie. Schemat
// nie ma dla niego kolumny: tabela `polaczenie` zna stan łącza, nie kartę
// w ognisku. Więź żyje więc w rdzeniu i kończy się wraz z jego procesem —
// sesja i jej procesy trwają dalej, a klient po restarcie rdzenia wskazuje
// kartę ponownie. Sesji, okien ani słowników ten rejestr nie
// przechowuje; mają własne repozytoria i własnego nadzorcę.
type wieziKlientow struct {
	mu    sync.RWMutex
	wpisy map[string]wiezKlienta
}

// wiezKlienta jest stanem jednego klienta.
type wiezKlienta struct {
	idSesji string
	idOkna  string
	ognisko time.Time
}

// zmianaOgniska niesie wynik przeniesienia ogniska: stan po zmianie, kartę
// tracącą ognisko i chwilę przeniesienia.
type zmianaOgniska struct {
	IdSesji    string
	IdOkna     string
	Poprzednia string
	Chwila     time.Time
}

// noweWieziKlientow zakłada pusty rejestr więzi.
func noweWieziKlientow() *wieziKlientow {
	return &wieziKlientow{wpisy: map[string]wiezKlienta{}}
}

// Ogniskuj przenosi ognisko klienta na wskazaną kartę sesji. Okno niewskazane
// zostawia okno dotychczasowe, o ile ognisko nie przechodzi na inną kartę —
// okno poprzedniej karty nie jest oknem karty nowej.
func (w *wieziKlientow) Ogniskuj(idKlienta, idSesji string, idOkna *string) zmianaOgniska {
	if w == nil {
		return zmianaOgniska{IdSesji: idSesji, Chwila: time.Now().UTC()}
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	poprzedni := w.wpisy[idKlienta]
	wpis := wiezKlienta{idSesji: idSesji, idOkna: poprzedni.idOkna, ognisko: time.Now().UTC()}
	if poprzedni.idSesji != idSesji {
		wpis.idOkna = ""
	}
	if idOkna != nil && *idOkna != "" {
		wpis.idOkna = *idOkna
	}
	w.wpisy[idKlienta] = wpis

	zmiana := zmianaOgniska{IdSesji: wpis.idSesji, IdOkna: wpis.idOkna, Chwila: wpis.ognisko}
	if poprzedni.idSesji != idSesji {
		zmiana.Poprzednia = poprzedni.idSesji
	}
	return zmiana
}

// Powiaz zapisuje sesję, z którą klient związał połączenie. Powiązanie nadaje
// też ognisko, gdy klient jeszcze żadnego nie ma — powrót do sesji jest
// powrotem do jej karty. Klient ogniskujący wcześniej inną kartę zachowuje
// swoje ognisko: powiązanie odtwarza stan, nie przestawia widoku.
func (w *wieziKlientow) Powiaz(idKlienta, idSesji string) {
	if w == nil || idSesji == "" {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	wpis := w.wpisy[idKlienta]
	if wpis.idSesji == "" {
		wpis.idSesji, wpis.ognisko = idSesji, time.Now().UTC()
	}
	w.wpisy[idKlienta] = wpis
}

// Ognisko zwraca kartę sesji i okno w ognisku klienta. Klient nieznany daje
// wartości puste — pierwsze wejście nie ma poprzednika.
func (w *wieziKlientow) Ognisko(idKlienta string) (string, string) {
	if w == nil {
		return "", ""
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	wpis := w.wpisy[idKlienta]
	return wpis.idSesji, wpis.idOkna
}
