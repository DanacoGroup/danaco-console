// Tożsamość żądania niesiona kontekstem: jeden identyfikator, którym klient
// wiąże żądanie, odpowiedź i cały strumień fragmentów tury.
package protocol

import "context"

// kluczIdZadania jest prywatnym, nieeksportowanym typem klucza wpisu kontekstu
// w tym pakiecie protocol.
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
// tury: identyfikator żądania, a gdy tury nie otworzyło żadne żądanie,
// tożsamość zastępczą podaną przez wywołującego.
func TozsamoscStrumienia(ctx context.Context, zastepcza string) string {
	if id := idZadania(ctx); id != "" {
		return id
	}
	return zastepcza
}
