// Odpowiedzialność pliku: więź między połączeniem a sesją bramki — to, co
// pozwala rdzeniowi wiedzieć, kto stoi po drugiej stronie gniazda.
//
// `auth.login` zakłada sesję i oddaje token, a token potrzebuje gdzie
// zamieszkać: transportem jest jedno gniazdo WebSocket, nie seria żądań HTTP,
// więc nie ma ani ciasteczka, ani nagłówka na każdym żądaniu. Bez tej więzi
// `connection.hello` nie zna tokenu, a `auth.password.reset` unieważniałby
// wszystkie sesje, bo bieżącej nie dałoby się wskazać.
//
// To nie jest bramka i nie wolno jej tu zbudować: jedynym miejscem kontroli
// jest logowanie do aplikacji, później zero blokad. Więź niczego nie sprawdza
// przed komendą, niczego nie odrzuca i nie zna pojęcia zakresu. Odpowiada na
// jedno pytanie — „z którą sesją bramki związane jest to połączenie" — i
// odpowiedź służy dokładnie dwóm rzeczom: powitaniu (czy klient ma pokazać
// okno logowania) i wyłączeniu bieżącej sesji ze zmiany hasła.
//
// Więź żyje w pamięci, nie w bazie. Wiąże byt nietrwały (połączenie, które
// znika z rozłączeniem) z bytem trwałym (wiersz `sesja_bramki`). Zapis do bazy
// byłby drugą prawdą o czymś, co i tak nie przeżywa restartu — a wzorzec
// rejestru pamięciowego odbudowywanego z wierszy rdzeń już ma.
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

// kluczPolaczenia niesie tożsamość połączenia nadaną przez transport.
//
// Wpis jest jeden. Niesie cały komplet (`transport.Tozsamosc`): identyfikator
// gniazda, identyfikator klienta, rodzaj programu po drugiej stronie i rolę
// jego okna. Drugiego wpisu obok nie ma — dwa wpisy o jednym połączeniu
// rozjeżdżają się zawsze, a rozjazd tożsamości znaczyłby, że Operator widzi na
// ekranie rękę nie tę, co trzeba.
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

// wiezBramki trzyma przypisania połączenie → skrót tokenu sesji bramki.
//
// Skrót, nie token. Token surowy zna klient i zna go rdzeń przez jedną chwilę
// przy zakładaniu sesji; nigdzie indziej nie ma prawa leżeć. Do obu
// zastosowań więzi skrót wystarcza: baza rozpoznaje sesję po skrócie i po
// skrócie wyłącza ją ze zmiany hasła.
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

// Skrot oddaje skrót tokenu związanego z połączeniem. Pusty wynik znaczy
// połączenie niezwiązane.
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

// kluczSesjiBiezacej niesie skrót tokenu sesji, z której przyszło żądanie.
const kluczSesjiBiezacej kluczKontekstu = "danaco:sesja-bramki"

// zSesjaBiezaca dokłada do kontekstu skrót sesji bramki wołającego.
//
// Po co kontekst, a nie pole żądania. Kształt żądania dyktuje kontrakt i tylko
// kontrakt; `auth.password.reset` nie niesie tokenu i nieść go nie
// będzie, bo klient już raz go oddał w powitaniu. Adapter musi jednak wiedzieć,
// której sesji nie unieważniać. Kontekst przenosi fakt o wywołaniu — nie treść
// żądania — i to jest właściwe miejsce na taką wiedzę.
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

// PolaczenieZwiazane odpowiada na jedno pytanie straży transportu (interfejs
// `transport.StanBramki`): czy to gniazdo przeszło przez bramkę.
//
// Metoda stoi tutaj, przy więzi, a nie przy adapterze transportu, bo cała jej
// treść to odczyt więzi i nic ponadto. Odpowiedź jest wyprowadzona z tego
// samego stanu, z którego bierze się `authenticated` w powitaniu — więc rdzeń
// nie może odpowiedzieć klientowi „nie jesteś zalogowany", a straży „jest
// związany" (dwie prawdy o jednej rzeczy).
//
// To nie jest sprawdzenie uprawnienia. Nie ma tu nazwy komendy, roli ani
// zakresu i mieć nie będzie: więź odpowiada „kto", nigdy „czy wolno".
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
