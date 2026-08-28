// Plik trzyma więź między połączeniem a sesją bramki: co pozwala rdzeniowi wiedzieć, kto stoi po drugiej stronie gniazda. Nie jest to bramka — jedynym miejscem kontroli jest logowanie, później zero blokad. Więź żyje w pamięci, nie w bazie.
package core

import (
	"context"
	"sync"

	"danacoconsole/server/internal/transport"
)

// kluczKontekstu jest typem własnym, żeby wartość w kontekście nie mogła się
// zderzyć z wartością innego pakietu. Napis jako klucz byłby zderzeniem
// czekającym na okazję.
type kluczKontekstu string

// kluczPolaczenia niesie tożsamość połączenia nadaną przez transport. Wpis jest jeden, niesie cały komplet: identyfikator gniazda, identyfikator klienta, rodzaj programu po drugiej stronie i rolę jego okna. Drugiego wpisu obok nie ma.
const kluczPolaczenia kluczKontekstu = "danaco:polaczenie"

// zPolaczeniem dokłada do kontekstu tożsamość połączenia, z którego przyszło
// żądanie. Wpina to `wejscieTransportu.Obsluz` — jedyne miejsce, przez które
// przechodzi każde żądanie z gniazda.
func zPolaczeniem(ctx context.Context, tozsamosc transport.Tozsamosc) context.Context {
	if ctx == nil || tozsamosc.IdPolaczenia == "" {
		return ctx
	}
	return context.WithValue(ctx, kluczPolaczenia, tozsamosc)
}

// tozsamoscZKontekstu odczytuje komplet faktów o wołającym. Wartość zerowa
// znaczy „żądanie nie przyszło z gniazda" — tak wygląda wywołanie z testu,
// z sondy stdio i z biegu wewnętrznego rdzenia. To nie jest błąd.
func tozsamoscZKontekstu(ctx context.Context) transport.Tozsamosc {
	if ctx == nil {
		return transport.Tozsamosc{}
	}
	tozsamosc, jest := ctx.Value(kluczPolaczenia).(transport.Tozsamosc)
	if !jest {
		return transport.Tozsamosc{}
	}
	return tozsamosc
}

// polaczenieZKontekstu odczytuje sam identyfikator połączenia — tyle, ile
// potrzebuje więź bramki. Pusty wynik znaczy żądanie spoza gniazda.
func polaczenieZKontekstu(ctx context.Context) string {
	return tozsamoscZKontekstu(ctx).IdPolaczenia
}

// wiezBramki trzyma przypisania połączenie-skrót tokenu sesji bramki. Skrót, nie token: token surowy zna klient i zna go rdzeń przez jedną chwilę przy zakładaniu sesji. Do obu zastosowań więzi skrót wystarcza.
type wiezBramki struct {
	zamek        sync.RWMutex
	poPolaczeniu map[string]string
}

// nowaWiezBramki zakłada pustą więź. Rdzeń bez wpiętej bramki dostaje ją tak
// samo — pusta więź odpowiada „nie wiadomo" i nic się nie psuje.
func nowaWiezBramki() *wiezBramki {
	return &wiezBramki{poPolaczeniu: make(map[string]string)}
}

// Zwiaz przypisuje połączeniu sesję bramki. Powtórzone wiązanie nadpisuje
// poprzednie: Operator, który zalogował się drugi raz na tym samym połączeniu,
// pracuje na sesji nowszej, nie na dwóch naraz.
func (w *wiezBramki) Zwiaz(polaczenie, skrot string) {
	if w == nil || polaczenie == "" || skrot == "" {
		return
	}
	w.zamek.Lock()
	defer w.zamek.Unlock()
	w.poPolaczeniu[polaczenie] = skrot
}

// Skrot oddaje skrót tokenu związanego z połączeniem; pusty wynik znaczy połączenie niezwiązane z sesją.
func (w *wiezBramki) Skrot(polaczenie string) string {
	if w == nil || polaczenie == "" {
		return ""
	}
	w.zamek.RLock()
	defer w.zamek.RUnlock()
	return w.poPolaczeniu[polaczenie]
}

// SkrotKontekstu jest skrótem myślowym dla dwóch wywołań, które i tak zawsze
// idą razem: wyjmij połączenie z kontekstu, oddaj jego skrót.
func (w *wiezBramki) SkrotKontekstu(ctx context.Context) string {
	return w.Skrot(polaczenieZKontekstu(ctx))
}

// kluczSesjiBiezacej niesie skrót tokenu sesji, z której przyszło żądanie, wpisany do kontekstu wywołania.
const kluczSesjiBiezacej kluczKontekstu = "danaco:sesja-bramki"

// zSesjaBiezaca dokłada do kontekstu skrót sesji bramki wołającego. Kontrakt nie niesie tokenu w polu żądania, bo klient już raz go oddał w powitaniu, a adapter musi wiedzieć, której sesji nie unieważniać. Kontekst przenosi fakt o wywołaniu.
func zSesjaBiezaca(ctx context.Context, skrot string) context.Context {
	if ctx == nil || skrot == "" {
		return ctx
	}
	return context.WithValue(ctx, kluczSesjiBiezacej, skrot)
}

// sesjaBiezacaZKontekstu oddaje skrót sesji wołającego. Pusty wynik znaczy
// „nie wiadomo, z której sesji przyszło żądanie" i prowadzi do
// unieważnienia wszystkich sesji.
func sesjaBiezacaZKontekstu(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	skrot, jest := ctx.Value(kluczSesjiBiezacej).(string)
	if !jest {
		return ""
	}
	return skrot
}

// PolaczenieZwiazane odpowiada na jedno pytanie straży transportu: czy gniazdo przeszło przez bramkę. Odpowiedź pochodzi z tego samego stanu, co authenticated w powitaniu. To nie jest sprawdzenie uprawnienia — więź odpowiada kto, nigdy czy wolno.
func (w wejscieTransportu) PolaczenieZwiazane(id string) bool {
	if w.rdzen == nil {
		return false
	}
	return w.rdzen.wiez.Skrot(id) != ""
}

// Rozwiaz zdejmuje przypisanie po rozłączeniu urządzenia. Bez tego mapa rosłaby
// przez całe życie procesu o jeden wpis na każde nawiązanie — więź pamięciowa
// bytu nietrwałego musi umieć zapomnieć.
func (w *wiezBramki) Rozwiaz(polaczenie string) {
	if w == nil || polaczenie == "" {
		return
	}
	w.zamek.Lock()
	defer w.zamek.Unlock()
	delete(w.poPolaczeniu, polaczenie)
}
