// Odpowiedzialność pliku: tożsamość żądania niesiona kontekstem — jeden
// identyfikator, którym klient wiąże żądanie, odpowiedź i cały strumień.
//
// Co mówi kontrakt — cytat z pola `Envelope.Id` (shared/contract.go):
//
//	„Identyfikator zadania; odpowiedz i fragmenty strumienia powtarzaja
//	 identyfikator zadania"
//
// Zdanie to nie zostawia wyboru: koperta `stream.chunk` ma nieść ten sam
// identyfikator, co koperta żądania, które turę otworzyło. Nadawca strumienia
// jednak żądania nie widzi — obsługiwacz komendy dostaje wyłącznie rozpakowany
// ładunek (`obsluz` w core/obsluga.go), a identyfikator zostaje w `Request`.
// Bez tego wpisu fragmenty tury i odpowiedź na `message.send` niosłyby różne
// identyfikatory — jedno pole o dwóch znaczeniach, zmieniające się w połowie
// strumienia.
//
// Dlaczego kontekst, a nie nowy parametr. Droga alternatywna to poszerzenie
// podpisu każdej czynności domenowej o identyfikator żądania — sto z górą
// podpisów zmienionych po to, by trzy z nich go użyły. Kontekst niesie już
// zasięg wykonania i odwołanie tury, więc tożsamość żądania jest tu bytem tej
// samej klasy i wchodzi jednym wpisem w `obsluz`, wspólnym dla wszystkich
// komend.
//
// Klucz jest typem prywatnym. Typ nieeksportowany wyklucza kolizję z kluczem
// innego pakietu — nikt spoza `protocol` nie ma jak zapisać ani nadpisać tego
// wpisu inaczej niż funkcją niżej.
//
// Brak wpisu nie jest błędem. Tura powołana poza drogą komendy
// (koordynator pętli, bieg naprawczy, podagent) nie ma żądania i dostaje napis
// pusty; wywołujący podstawia wtedy własną tożsamość zastępczą i mówi o tym
// wprost, zamiast udawać żądanie, którego nie było.
package protocol

import "context"

// kluczIdZadania jest prywatnym kluczem wpisu kontekstu.
type kluczIdZadania struct{}

// ZIdZadania zwraca kontekst niosący identyfikator żądania. Identyfikator pusty
// nie zakłada wpisu — kontekst zostaje nietknięty, więc odczyt niżej odróżnia
// „żądania nie było" od „żądanie miało puste pole".
func ZIdZadania(ctx context.Context, id string) context.Context {
	if ctx == nil || id == "" {
		return ctx
	}
	return context.WithValue(ctx, kluczIdZadania{}, id)
}

// idZadania odczytuje identyfikator żądania z kontekstu. Brak wpisu daje napis
// pusty — pytanie o tożsamość żądania poza drogą komendy jest zasadne i ma
// odpowiedź, a nie awarię.
func idZadania(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, jest := ctx.Value(kluczIdZadania{}).(string)
	if !jest {
		return ""
	}
	return id
}

// TozsamoscStrumienia rozstrzyga identyfikator koperty dla strumienia jednej
// tury: identyfikator żądania, a gdy tury nie otworzyło żadne żądanie —
// tożsamość zastępczą podaną przez wywołującego.
//
// Funkcja istnieje po to, żeby rozstrzygnięcie stało w jednym miejscu. Gdyby
// każdy nadawca strumienia wybierał sam, „ten sam identyfikator przez cały
// strumień" byłby obietnicą powtarzaną w kilku plikach — a obietnica
// powtórzona to obietnica, którą któryś z nich kiedyś złamie.
func TozsamoscStrumienia(ctx context.Context, zastepcza string) string {
	if id := idZadania(ctx); id != "" {
		return id
	}
	return zastepcza
}
