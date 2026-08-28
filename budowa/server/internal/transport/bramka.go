package transport

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Dopuszczenie do rdzenia rozstrzyga się tam, gdzie logowanie do aplikacji ma skutek dla rdzenia i jego połączeń.

// komendyWejscia to jedyne komendy wykonywane przez połączenie, które jeszcze nie przeszło przez bramkę wejścia rdzenia.
var komendyWejscia = map[shared.MessageType]struct{}{
	shared.CommandConnectionHello:  {},
	shared.CommandAuthLogin:        {},
	shared.CommandAuthRegister:     {},
	shared.CommandAuthVerify:       {},
	shared.CommandAuthRecover:      {},
	shared.CommandAuthReset:        {},
	shared.CommandAuthTokenRefresh: {},
}

// StanBramki jest rozszerzeniem nieobowiązkowym interfejsu Rdzen, pytającym, czy połączenie ma ważną sesję bramki.
type StanBramki interface {
	PolaczenieZwiazane(id string) bool
}

// dopuszczenieBramki rozstrzyga, czy żądanie z danego gniazda wolno oddać rdzeniowi, na podstawie adresu nasłuchu i wskazania Operatora.
type dopuszczenieBramki struct {
	// wymagana mówi, że gniazdo musi się przedstawić, zanim cokolwiek wykona.
	wymagana bool
}

// wymogLogowania składa obie przesłanki w jedno zdanie. Wskazanie Operatora
// (nastawa niezerowa) wygrywa z adresem w OBIE strony — bo to jest dźwignia,
// a nie podpowiedź; brak wskazania oddaje głos adresowi nasłuchu.
func wymogLogowania(adres string, nastawa *bool) bool {
	if nastawa != nil {
		return *nastawa
	}
	return !petlaZwrotna(adres)
}

// Metoda dopuszczenie składa regułę z ustawień okna; pętla zwrotna bez wskazania Operatora daje regułę wyłączoną.
func (u Ustawienia) dopuszczenie() dopuszczenieBramki {
	return dopuszczenieBramki{wymagana: wymogLogowania(u.Adres, u.WymogLogowania)}
}

// Metoda przepusc mówi, czy żądanie idzie dalej do rdzenia, pytając wyłącznie o to, czy dane gniazdo przeszło przez bramkę.
func (s dopuszczenieBramki) przepusc(rdzen Rdzen, komenda shared.MessageType, ujscie Ujscie) bool {
	if !s.wymagana {
		return true
	}
	if _, wejscie := komendyWejscia[komenda]; wejscie {
		return true
	}
	stan, zna := rdzen.(StanBramki)
	if !zna || ujscie == nil {
		return false
	}
	return stan.PolaczenieZwiazane(ujscie.Id())
}

// Funkcja odmowaBezBramki składa jedyną odmowę tego dopuszczenia, kodem oznaczającym brak uwierzytelnienia z kontraktu.
func odmowaBezBramki(zadanie protocol.Request) protocol.Koperta {
	blad := protocol.NowyBlad(shared.ErrorCodeNotAuthenticated,
		"to połączenie nie przeszło przez bramkę — zaloguj się; "+
			"na tym nasłuchu obowiązuje wymóg logowania, więc gniazdo wykonuje komendy dopiero po "+
			"przedstawieniu tokenu sesji w connection.hello. Token wydaje auth.login, "+
			"a przy pierwszym uruchomieniu auth.verify — czyli potwierdzenie adresu drogą z listu, "+
			"nie samo auth.register. Po przejściu bramki nie ma już ani jednego pytania")
	return protocol.KopertaBledu(zadanie.Koperta(), blad)
}
