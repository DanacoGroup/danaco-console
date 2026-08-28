// Maszynowa połowa komunikatu odmowy — kod błędu kontraktu, który jedzie do
// Operatora obok zdania. Pakiet wybiera wyłącznie kod; ponawialność bierze
// się z katalogu kontraktu.
package zewnetrzne

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// KodOdmowy dobiera kod kontraktu do rozpoznanego rodzaju odmowy dostawcy,
// każdy opisany osobnym zdaniem w `shared/contract.go`. Pozostałe rodzaje
// zostają przy `channel_unavailable`.
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
// wysyłką. Niewpięty sejf jest błędem rdzenia, nie Operatora, i kod mówi to
// samo, co zdanie.
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
func (o *OdmowaKanaluZewnetrznego) Unwrap() error {
	return protocol.JakoError(o.Blad())
}

// Blad podaje brak poświadczenia w kształcie kontraktu, gotowym do
// przeniesienia do Operatora produktu.
func (b *BrakPoswiadczenia) Blad() protocol.Blad {
	return protocol.NowyBlad(KodBraku(b.Powod), b.Error())
}

// Unwrap wystawia błąd kontraktu w łańcuchu — tak samo jak przy odmowie
// dostawcy i z tego samego powodu.
func (b *BrakPoswiadczenia) Unwrap() error {
	return protocol.JakoError(b.Blad())
}
