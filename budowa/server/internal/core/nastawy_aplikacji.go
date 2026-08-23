// Odpowiedzialność pliku: droga odczytu nastaw poziomu `aplikacja`
// (`migracja_119_zasieg_aplikacja.sql`) — tych, które opisują sam program,
// a nie treść w nim prowadzoną.
//
// Poza tą drogą wymóg logowania daje się ustawić wyłącznie przy starcie
// rdzenia: przełącznikiem wiersza poleceń (`--wymog-logowania`) albo zmienną
// środowiska (`DANACO_WYMOG_LOGOWANIA`). Poziom zasięgu, na którym wolno zapisać
// tę nastawę komendą `config.set`, wnosi wymieniona migracja; drogę odczytu
// wnosi ten plik.
//
// ── odczyt jest przy nawiązaniu połączenia ─────────────────────────────────
// Warstwa nasłuchu składa straż bramki raz na połączenie (`transport/
// nawiazanie.go` — `s.ustawienia.straznik()`), a nie raz na bieg rdzenia.
// Chwilą, w której nastawa ma znaczenie, jest więc chwila nawiązania — i tam
// rdzeń ją czyta, przy powitaniu `connection.hello`. Zmiana zapisana komendą
// `config.set` obowiązuje od następnego połączenia, bez restartu.
//
// ── drugiego mechanizmu nastaw tu nie ma ───────────────────────────────────
// Nie ma tu pamięci podręcznej, własnego pliku, własnej tabeli ani własnego
// stanu. Jest jedno wywołanie tego samego rozstrzygacza, którym idzie każde inne
// ustawienie platformy, po klucz z tego samego rejestru definicji. Wartość
// mieszka w tabeli `ustawienie` i widzi ją `config.get` tak samo jak każdą inną.
//
// ── i nie jest to bramka ───────────────────────────────────────────────────
// Odczyt niczego nie odmawia i nikogo nie zatrzymuje. Powitanie oddaje wynik
// w polu `loginRequired`, żeby klient wiedział, czy pokazać okno logowania,
// zamiast wyprowadzać to z odmowy.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/konfig"
)

// NastawyAplikacji jest portem odczytu nastaw poziomu `aplikacja`. Rozszerzenie
// nieobowiązkowe portu Ustawienia — tą samą drogą, którą rdzeń pyta ten sam
// adapter o prowenancję konfiguracji (`ProwenancjaKonfiguracji`).
type NastawyAplikacji interface {
	// WymogLogowania zwraca nastawę wraz z informacją, czy w ogóle jest
	// wskazana. Trzy stany, nie dwa: druga wartość fałszywa znaczy „Operator nie
	// wskazał nic”, co jest czymś innym niż „wskazał: nie” i oddaje głos adresowi
	// nasłuchu (transport/bramka.go).
	WymogLogowania(ctx context.Context) (bool, bool)
}

// WymogLogowania wypełnia port NastawyAplikacji na adapterze ustawień.
//
// Kontekst rozstrzygania jest pusty z zamysłem. Nastawa poziomu `aplikacja`
// nie ma bytu —
// programu nie ma czym zawęzić — a rozstrzyganie i tak przejdzie po ośmiu
// węższych poziomach, na których tego klucza nikt nie zapisze, bo definicja
// dopuszcza wyłącznie poziom aplikacji.
func (a *adapterUstawienOsi) WymogLogowania(ctx context.Context) (bool, bool) {
	if a == nil || a.rozstrzygacz == nil {
		return false, false
	}
	// Rozstrzygacz czyta źródło przy każdym wywołaniu, więc wartość jest ta,
	// która stoi w bazie w tej chwili — o to chodzi w drodze „na żywo”.
	_ = ctx
	wynik := a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{}, konfig.KluczWymogLogowania)
	return wartoscWymoguLogowania(wynik.Wartosc)
}

// wartoscWymoguLogowania przekłada zapis kolumny na parę (wartość, wskazana).
// Zapis pusty znaczy brak wskazania; zapis nieczytelny również — bo nastawa,
// której nie da się odczytać, nie jest wskazaniem, a zgadywanie jej sensu byłoby
// rozstrzyganiem za Operatora.
func wartoscWymoguLogowania(zapis string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(zapis)) {
	case "true", "1", "tak":
		return true, true
	case "false", "0", "nie":
		return false, true
	default:
		return false, false
	}
}
