// Plik trzyma więź między połączeniem a sesją bramki: co pozwala rdzeniowi wiedzieć, kto stoi po drugiej stronie gniazda. Nie jest to bramka — jedynym miejscem kontroli jest logowanie, później zero blokad. Więź żyje w pamięci, nie w bazie.
package core

import (
	"context"
	"sync"
	"time"

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

// zrodloZadania oddaje adres urządzenia, z którego przyszło żądanie. Pustka
// znaczy kontekst bez gniazda — tak wygląda praca procesu bez zamawiającego.
func zrodloZadania(ctx context.Context) string {
	return tozsamoscZKontekstu(ctx).AdresZrodlowy
}

// wiezBramki trzyma przypisania połączenie-skrót tokenu sesji bramki. Skrót, nie token: token surowy zna klient i zna go rdzeń przez jedną chwilę przy zakładaniu sesji. Do obu zastosowań więzi skrót wystarcza.
type wiezBramki struct {
	zamek        sync.RWMutex
	poPolaczeniu map[string]wpisWiezi
}

// wpisWiezi niesie sesję gniazda wraz z chwilą jej wygaśnięcia. Chwila jest tu
// po to, żeby sesja, o której już wiadomo, że minęła, nie kosztowała zapytania
// do bazy przy każdej komendzie tego gniazda. Zero znaczy „jeszcze nieodczytana".
type wpisWiezi struct {
	skrot  string
	wygasa int64
}

// nowaWiezBramki zakłada pustą więź. Rdzeń bez wpiętej bramki dostaje ją tak
// samo — pusta więź odpowiada „nie wiadomo" i nic się nie psuje.
func nowaWiezBramki() *wiezBramki {
	return &wiezBramki{poPolaczeniu: make(map[string]wpisWiezi)}
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
	w.poPolaczeniu[polaczenie] = wpisWiezi{skrot: skrot}
}

// Skrot oddaje skrót tokenu związanego z połączeniem; pusty wynik znaczy połączenie niezwiązane z sesją.
func (w *wiezBramki) Skrot(polaczenie string) string {
	wpis := w.wpis(polaczenie)
	return wpis.skrot
}

// wpis oddaje komplet zapamiętany o gnieździe: sesję i znaną chwilę jej wygaśnięcia.
func (w *wiezBramki) wpis(polaczenie string) wpisWiezi {
	if w == nil || polaczenie == "" {
		return wpisWiezi{}
	}
	w.zamek.RLock()
	defer w.zamek.RUnlock()
	return w.poPolaczeniu[polaczenie]
}

// zapamietajWygasniecie dopisuje do wpisu chwilę wygaśnięcia odczytaną z wiersza
// sesji. Wpis zdjęty w międzyczasie nie wraca — gniazdo rozłączone ma zostać
// rozłączone.
func (w *wiezBramki) zapamietajWygasniecie(polaczenie string, wygasa int64) {
	if w == nil || polaczenie == "" || wygasa <= 0 {
		return
	}
	w.zamek.Lock()
	defer w.zamek.Unlock()
	wpis, stoi := w.poPolaczeniu[polaczenie]
	if !stoi {
		return
	}
	wpis.wygasa = wygasa
	w.poPolaczeniu[polaczenie] = wpis
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

// RozpoznanieWaznosciSesji oddaje chwilę wygaśnięcia sesji bramki i to, czy
// sesja nadal nadaje. Jest rozszerzeniem nieobowiązkowym rozpoznania konta:
// rdzeń bez trwałości uwierzytelnienia nie ma czego czytać.
type RozpoznanieWaznosciSesji interface {
	WaznoscSesjiBramki(ctx context.Context, skrotTokenu string) (int64, bool, error)
}

// czasSprawdzeniaWaznosci ogranicza odczyt wiersza sesji przy pytaniu bramki.
// Odczyt idzie z pętli odbioru gniazda, więc zapytanie wstrzymane na zamku bazy
// wstrzymałoby całe połączenie, a nie jedną komendę.
const czasSprawdzeniaWaznosci = 2 * time.Second

/*
PolaczenieZwiazane odpowiada na jedno pytanie straży transportu: czy gniazdo
przeszło przez bramkę i czy jego sesja nadal nadaje. Sama więź na to nie
odpowiada: wpis w mapie powstaje przy powitaniu i przeżyłby unieważnienie sesji
zrobione gdzie indziej — `auth.reset`, `auth.password.reset` i `device.revoke`
meldują wtedy Operatorowi liczbę odebranych urządzeń, która nie byłaby prawdą.

Czytany jest wiersz sesji, a nie jej migawka: unieważnienie ma odciąć gniazdo
natychmiast, nie przy najbliższym rozłączeniu. Zapytania nie ma tylko wtedy, gdy
chwila wygaśnięcia zapamiętana w więzi już minęła.
*/
func (w wejscieTransportu) PolaczenieZwiazane(id string) bool {
	if w.rdzen == nil {
		return false
	}
	wpis := w.rdzen.wiez.wpis(id)
	if wpis.skrot == "" {
		return false
	}
	if wpis.wygasa > 0 && wpis.wygasa <= time.Now().UnixMilli() {
		w.rdzen.wiez.Rozwiaz(id)
		return false
	}
	waznosc, umie := w.rdzen.konta.(RozpoznanieWaznosciSesji)
	if !umie {
		// Rdzeń bez trwałości uwierzytelnienia nie ma wiersza sesji do
		// przeczytania; więź jest wtedy całą wiedzą, jaka o gnieździe jest.
		return true
	}
	ctx, koniec := context.WithTimeout(context.Background(), czasSprawdzeniaWaznosci)
	defer koniec()
	wygasa, wazna, err := waznosc.WaznoscSesjiBramki(ctx, wpis.skrot)
	if err != nil {
		// Nierozstrzygnięta ważność zamyka drogę: przepuszczenie żądania byłoby
		// tu wpuszczeniem sesji, o której nic nie wiadomo.
		if w.rdzen.dziennik != nil {
			w.rdzen.dziennik.Printf("core: ważność sesji gniazda %s nierozstrzygnięta: %v", id, err)
		}
		return false
	}
	if !wazna {
		w.rdzen.wiez.Rozwiaz(id)
		return false
	}
	w.rdzen.wiez.zapamietajWygasniecie(id, wygasa)
	return true
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
