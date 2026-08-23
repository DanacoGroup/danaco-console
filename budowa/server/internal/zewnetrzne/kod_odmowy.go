// Maszynowa połowa komunikatu odmowy — kod błędu kontraktu, który jedzie do
// Operatora obok zdania.
//
// Zdanie dla człowieka (`uwierzytelnienie.go`, `brak_poswiadczenia.go`) to
// dopiero połowa tego, co Operator czyta. Drugą połowę — kod kontraktu i
// wynikającą z niego ponawialność — klient dokleja do wpisu rozmowy sam. Rdzeń
// domyka turę jednym kodem zapasowym dla każdej odmowy kanału
// (`channel_unavailable`), a katalog kontraktu (`shared.KodyPonawialne`) uznaje
// ten kod za ponawialny. Bez kodu własnego odwołany albo wygasły klucz
// przychodziłby więc podpisany „ponowienie ma sens", choć zdanie obok mówi
// „wymień klucz".
//
// `protocol.BladZeZrodla` najpierw pyta `errors.As`, czy błąd niesie już swój
// kod, i dopiero gdy nie niesie — nakłada kod zapasowy wołającego. Obie odmowy
// tego pakietu niosą kod przez `Unwrap`, więc kod zapasowy nie wygrywa i żaden
// plik poza tym pakietem nie musi się zmieniać.
//
// Pakiet wybiera wyłącznie kod. Ponawialność bierze się z katalogu kontraktu
// przez `protocol.NowyBlad`; ustawiana tutaj byłaby drugą prawdą o ponawianiu
// obok kontraktowej.
package zewnetrzne

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// KodOdmowy dobiera kod kontraktu do rozpoznanego rodzaju odmowy dostawcy.
//
// Przy każdym kodzie stoi zdanie, którym opisuje go `shared/contract.go`:
//
//   - 401 → `not_authenticated` („Brak uwierzytelnienia"). Dostawca nie uznał
//     klucza. Nieponawialny — i tak ma być: ten sam klucz odbije się tak samo.
//   - 403 → `permission_denied` („Uwierzytelniony, lecz bez uprawnienia do
//     czynnosci"). Dokładnie to powiedział dostawca: klucz zna, czynności nie da.
//   - 402 → `permission_denied`. Konto jest uwierzytelnione, a mimo to nie ma
//     prawa do czynności — tyle że powodem jest rozliczenie, nie zasięg klucza.
//     Powód nazywa zdanie; kod ma tu powiedzieć jedno: ponawianie nie pomoże.
//   - 429 → `rate_limited` („Ograniczenie tempa po stronie kanalu lub rdzenia").
//     Ponawialny, i słusznie — to jedyna z tych odmów, która mija sama.
//   - 5xx → `channel_unavailable` („Kanal modelu niedostepny... tylko biezacego
//     wywolania"). Awaria po stronie dostawcy jest właśnie chwilowa.
//
// Pozostałe rodzaje zostają przy `channel_unavailable`. Odrzucone żądanie
// (400, zły model, złe ciało) nie jest wprawdzie chwilowe, ale jego naprawa leży
// w parametrach wiersza kanału, a `validation_failed` mówiłby wołaczom, że to
// ich żądanie było niezgodne z kontraktem rdzenia.
func KodOdmowy(rodzaj string) shared.ErrorCode {
	switch rodzaj {
	case OdmowaUwierzytelnienia:
		return shared.ErrorCodeNotAuthenticated
	case OdmowaUprawnienia, OdmowaPlatnosci:
		return shared.ErrorCodePermissionDenied
	case OdmowaNatezenia:
		return shared.ErrorCodeRateLimited
	default:
		return shared.ErrorCodeChannelUnavailable
	}
}

// KodBraku dobiera kod kontraktu do braku poświadczenia stwierdzonego przed
// wysyłką.
//
// Niewpięty sejf jest błędem rdzenia, nie Operatora, i kod mówi to samo, co
// zdanie: `internal_error`. Gdyby szedł jako `not_authenticated`, Operator
// dostałby maszynowe potwierdzenie, że zawinił jego klucz — a klucza nikt tam
// jeszcze nie czytał, bo nie było czym. Pozostałe powody to rzeczywiście brak
// uwierzytelnienia i żadne ponowienie ich nie usunie.
func KodBraku(powod string) shared.ErrorCode {
	if powod == PowodSejfNiewpiety {
		return shared.ErrorCodeInternalError
	}
	return shared.ErrorCodeNotAuthenticated
}

// Blad podaje odmowę dostawcy w kształcie kontraktu: kod maszynowy i to samo
// zdanie, które czyta człowiek. Jedna treść w obu połowach — rozjazd między
// nimi byłby drugą prawdą o tej samej odmowie.
func (o *OdmowaKanaluZewnetrznego) Blad() protocol.Blad {
	return protocol.NowyBlad(KodOdmowy(o.Rodzaj), o.Error())
}

// Unwrap wystawia błąd kontraktu w łańcuchu, żeby `protocol.BladZeZrodla`
// znalazł go przez `errors.As` i nie nakładał kodu zapasowego wołającego.
//
// Na tym opiera się przeniesienie kodu do Operatora: bez tego ogniwa kod
// musiałby wybierać `core/strumien_odpowiedzi.go` i każde inne miejsce, które
// odmowę kanału przenosi do protokołu. `Error()` pozostaje zdaniem dla
// człowieka, a `errors.As(err, &OdmowaKanaluZewnetrznego{})` działa dalej —
// ogniwo dokłada się za typem, nie zamiast niego.
func (o *OdmowaKanaluZewnetrznego) Unwrap() error {
	return protocol.JakoError(o.Blad())
}

// Blad podaje brak poświadczenia w kształcie kontraktu.
func (b *BrakPoswiadczenia) Blad() protocol.Blad {
	return protocol.NowyBlad(KodBraku(b.Powod), b.Error())
}

// Unwrap wystawia błąd kontraktu w łańcuchu — tak samo jak przy odmowie
// dostawcy i z tego samego powodu.
func (b *BrakPoswiadczenia) Unwrap() error {
	return protocol.JakoError(b.Blad())
}
