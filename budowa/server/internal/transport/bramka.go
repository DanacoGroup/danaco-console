package transport

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Straż wejścia rdzenia. Bez niej klient bez tokenu wykonuje dowolną komendę,
// w tym `terminal.command.exec`, która przez tor zdalny (`zdalne.Zasil`) potrafi
// sięgnąć po SSH na maszynę Operatora. Zgoda w tabeli `host_zdalny` jest wydana
// hostowi, nie wołającemu, więc wystawiony rdzeń oddawałby cudzemu połączeniu VPS
// i komputer w biurze.
//
// To nie jest bramkowanie uprawnień — to wskazanie miejsca, w którym logowanie do
// aplikacji ma skutek dla rdzenia. Straż stąd:
//
//   - nie zna pojęcia uprawnienia, roli, zakresu ani modułu — pyta o jedno:
//     czy to gniazdo przeszło przez bramkę;
//   - nie pyta o to ani razu na pętli zwrotnej — Operator na własnej maszynie
//     nie zobaczy jednego nowego pytania, bo straż jest tam wyłączona w całości;
//   - po przejściu bramki milczy do końca życia połączenia.
//
// Odmowa opisuje brak, nie zakaz. Nie „nie wolno ci tej komendy", tylko „to
// połączenie nie przeszło przez bramkę" — bo dokładnie to jest faktem, a droga
// naprawy (zaloguj się) stoi w tym samym zdaniu.

// komendyWejscia to jedyne komendy wykonywane przez połączenie, które bramki
// jeszcze nie przeszło. Wykaz jest wyczerpujący i wynika z jednego pytania:
// czego nie da się pominąć, żeby móc się zalogować.
//
//   - `connection.hello` — tędy token wchodzi do rdzenia (transportem
//     jest jedno gniazdo, więc nie ma nagłówka na każdym żądaniu);
//   - `auth.login` — tędy token powstaje;
//   - `auth.register` — bez tego rdzeń bez założonej bramki byłby zamknięty
//     na klucz, którego nikt jeszcze nie wykuł;
//   - `auth.verify` — rejestracja jest dwukrokowa i to TEN krok wydaje token.
//     Bez niego `auth.register` zakłada konto niepotwierdzone, którego już nic
//     nie potwierdzi: rejestracja drugi raz oddaje `conflict`, a logowanie
//     odmawia zdaniem o oczekiwaniu na potwierdzenie adresu;
//   - `auth.recover` i `auth.reset` — dwa kroki odzyskania konta. Naciska je
//     ten, kto hasła nie pamięta, czyli z definicji przez bramkę nie przejdzie;
//     odbite zamieniają zapomniane hasło w koniec instalacji;
//   - `auth.token.refresh` — przedłużenie sesji zapisanej na maszynie. Token
//     przedstawiony w powitaniu wiąże gniazdo i wtedy przedłużenie przechodzi
//     samo, ale token wygasły gniazda nie wiąże — i wtedy Operator ma zobaczyć
//     odmowę rdzenia „sesja wygasła", a nie odmowę straży, która o sesji nic
//     nie mówi.
//
// Wykaz nie jest furtką: żadna z tych komend nie wykonuje pracy Operatora, nie
// sięga po pliki, sieć ani powłokę — wszystkie dotykają wyłącznie bramki.
// Zgadywanie po nich jest ograniczone dwiema rzeczami: dławikiem prób wejścia
// (`adapter_modul_auth_dlawik.go`) i tym, że droga potwierdzenia jest losowa,
// jednorazowa i wygasa po godzinie (`adapter_modul_auth_konto.go`).
//
// Nazwy biorą się ze stałych kontraktu, nie z literałów: zmiana nazwy
// komendy w kontrakcie ma wywrócić kompilację, a nie po cichu zamknąć wejście.
var komendyWejscia = map[shared.MessageType]struct{}{
	shared.CommandConnectionHello:  {},
	shared.CommandAuthLogin:        {},
	shared.CommandAuthRegister:     {},
	shared.CommandAuthVerify:       {},
	shared.CommandAuthRecover:      {},
	shared.CommandAuthReset:        {},
	shared.CommandAuthTokenRefresh: {},
}

// StanBramki jest rozszerzeniem nieobowiązkowym interfejsu Rdzen — tą samą
// drogą, którą transport pyta rdzeń o obserwację połączeń (ObserwatorPolaczen).
// Transport nadal nie zna rdzenia: zna pytanie i kształt odpowiedzi.
//
// Odpowiedź brzmi „czy połączenie o tym identyfikatorze jest związane z ważną
// sesją bramki". Więź jest jedna i mieszka w rdzeniu (`core/wiez_polaczenia.go`);
// transport nie zakłada drugiej i nie trzyma własnego stanu uwierzytelnienia,
// bo dwa źródła prawdy o jednej rzeczy rozjeżdżają się zawsze.
type StanBramki interface {
	PolaczenieZwiazane(id string) bool
}

// straznikBramki rozstrzyga, czy żądanie z tego gniazda wolno oddać rdzeniowi.
//
// Wymóg bierze się z dwóch rzeczy, nie z jednej:
//
//  1. adres nasłuchu — poza pętlą zwrotną wymóg obowiązuje sam z siebie. Adres
//     jest faktem, a nie nastawą, więc wystawienia nie da się zrobić „przez
//     zapomnienie";
//  2. jawne wskazanie Operatora — dźwignia, którą wymóg włącza się także na
//     pętli zwrotnej (kto pracuje na wspólnej maszynie, ma czym się zamknąć)
//     albo znosi się przy nasłuchu szerszym. Zniesienie jest dozwolone, ale
//     nigdy ciche: dziennik mówi wtedy wprost, co stoi otworem
//     (`wystawienie.go`).
//
// Nastawy nie ma → rozstrzyga adres. To jest cała reguła.
type straznikBramki struct {
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

// straznik składa straż z ustawień. Pętla zwrotna bez wskazania Operatora daje
// straż wyłączoną.
func (u Ustawienia) straznik() straznikBramki {
	return straznikBramki{wymagana: wymogLogowania(u.Adres, u.WymogLogowania)}
}

// przepusc mówi, czy żądanie idzie dalej do rdzenia.
//
// Trzy wyjścia na „tak" i jedno na „nie". Straż wyłączona przepuszcza wszystko;
// komenda wejścia przechodzi zawsze; związane gniazdo przechodzi zawsze. Zostaje
// jeden przypadek: wymóg obowiązuje, komenda spoza wejścia, gniazdo
// nieprzedstawione.
//
// Pytanie jest o gniazdo, nie o tożsamość wołającego — i dlatego przeżyje zmianę
// modelu bramki. Dziś sesja bramki nie ma właściciela (Operator jest bezimienny,
// wzorzec Danaco HUB); gdy dostanie konto i wiele urządzeń z osobnymi tokenami,
// straż nie wymaga ani jednej zmiany: nadal pyta „czy to gniazdo przeszło
// bramkę", a odpowiedź nadal daje więź.
//
// Rdzeń nieznający rozszerzenia nie przepuszcza — to jedyne miejsce w tym
// pakiecie, gdzie brak czegoś zamyka drogę zamiast ją otwierać. Rdzeń, który nie
// umie odpowiedzieć „kto woła", przy obowiązującym wymogu oddawałby komendy
// komukolwiek. Bez wymogu to rozstrzygnięcie nie ma jak zadziałać, bo straż jest
// wtedy wyłączona wcześniej.
func (s straznikBramki) przepusc(rdzen Rdzen, komenda shared.MessageType, ujscie Ujscie) bool {
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

// odmowaBezBramki składa jedyną odmowę tej straży.
//
// Kod jest jeden: `not_authenticated` z kontraktu. Nie `permission_denied`
// i nie `validation_failed` — bo brakuje nie uprawnienia i nie pola w żądaniu,
// tylko przejścia przez bramkę. Klient rozpoznaje ten kod i otwiera okno
// logowania zamiast pokazywać Operatorowi błąd komendy.
func odmowaBezBramki(zadanie protocol.Request) protocol.Koperta {
	blad := protocol.NowyBlad(shared.ErrorCodeNotAuthenticated,
		"to połączenie nie przeszło przez bramkę — zaloguj się; "+
			"na tym nasłuchu obowiązuje wymóg logowania, więc gniazdo wykonuje komendy dopiero po "+
			"przedstawieniu tokenu sesji w connection.hello. Token wydaje auth.login, "+
			"a przy pierwszym uruchomieniu auth.verify — czyli potwierdzenie adresu drogą z listu, "+
			"nie samo auth.register. Po przejściu bramki nie ma już ani jednego pytania")
	return protocol.KopertaBledu(zadanie.Koperta(), blad)
}
