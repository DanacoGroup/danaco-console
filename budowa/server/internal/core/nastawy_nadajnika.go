// Odpowiedzialność pliku: droga odczytu konta nadawczego platformy — skrzynki,
// z której rdzeń pisze dwa listy systemowe (potwierdzenie adresu przy
// rejestracji i droga odzyskania konta).
//
// ── dlaczego nie wystarczy start ───────────────────────────────────────────
// Do tej pory konto nadawcze wchodziło wyłącznie przy starcie rdzenia
// (`DANACO_NADAWCA_*`, `konfiguracja/srodowisko.go`). Skutek: pomyłka w adresie
// serwera poczty zamykała rejestrację, a naprawa wymagała zatrzymania rdzenia
// i wiedzy spoza produktu. Nastawy z katalogu (`migracja_372`) dają drogę
// wewnątrz okna Konfiguracji i nie odbierają tamtej: start nadal zasila
// wartości, a zapis w tabeli `ustawienie` je przesłania.
//
// ── nakładka, nie druga prawda ─────────────────────────────────────────────
// Nastawa pusta niczego nie kasuje — pusta znaczy „nie wskazałem”, więc zostaje
// wartość ze startu. Inaczej pierwsze wejście do okna Konfiguracji kasowałoby
// konto nadawcze podane środowiskiem samym faktem, że pola stoją puste.
//
// Wyjątkiem jest szyfrowanie: tam nastawa ma trzy stany tak samo jak wymóg
// logowania (`nastawy_aplikacji.go`) — brak wskazania oddaje głos startowi,
// a wskazanie wygrywa w obie strony, bo zejście do rozmowy otwartym tekstem ma
// być zapisem jawnym.
//
// ── drugiego mechanizmu nastaw tu nie ma ───────────────────────────────────
// Jedno wywołanie tego samego rozstrzygacza, którym idzie każde inne ustawienie
// platformy, po klucze z tego samego rejestru definicji. Bez pamięci podręcznej:
// list idzie rzadko, a nastawa poprawiona po nieudanej próbie ma obowiązywać
// przy próbie następnej, nie po restarcie.
package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/nadajnik"
)

// dwieDrogiKontaNadawczego nazywa obie drogi ustawienia skrzynki nadawczej.
// Zdanie stoi w odmowie rejestracji i odzyskania konta, bo tam odmowa zatrzymuje
// Operatora przed platformą, a odesłanie do dziennika rdzenia prowadzi tam,
// gdzie on nie sięga.
const dwieDrogiKontaNadawczego = "konto nadawcze ustawia się na dwa sposoby: " +
	"zmiennymi środowiska przy starcie rdzenia (DANACO_NADAWCA_HOST, DANACO_NADAWCA_ADRES " +
	"i pozostałe DANACO_NADAWCA_*) albo nastawami mailer.host i mailer.address " +
	"w oknie Konfiguracji, w kategorii „Konto nadawcze platformy”"

// NastawyPlatformy jest portem odczytu pojedynczej nastawy poziomu `aplikacja`.
// Rozszerzenie nieobowiązkowe portu Ustawienia — tą samą drogą, którą rdzeń pyta
// ten sam adapter o wymóg logowania (`NastawyAplikacji`).
type NastawyPlatformy interface {
	// Nastawa zwraca wartość zapisaną pod kluczem; napis pusty znaczy „brak
	// wskazania”, a nie „wartość pusta” — kasowania wartości ze startu ta droga
	// nie zna.
	Nastawa(ctx context.Context, klucz string) string
}

// Nastawa wypełnia port NastawyPlatformy na adapterze ustawień.
//
// Kontekst rozstrzygania jest pusty z zamysłem: nastawy poziomu `aplikacja` nie
// mają bytu — programu nie ma czym zawęzić.
func (a *adapterUstawienOsi) Nastawa(ctx context.Context, klucz string) string {
	if a == nil || a.rozstrzygacz == nil {
		return ""
	}
	_ = ctx
	return strings.TrimSpace(a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{}, klucz).Wartosc)
}

// kontoNadawcze nakłada nastawy Operatora na konto nadawcze podane przy starcie.
//
// Źródło nieobecne (montaż bez tego ogniwa) oddaje konto ze startu bez zmian —
// brak drogi odczytu nie jest wskazaniem Operatora.
func kontoNadawcze(ctx context.Context, startu nadajnik.Nastawy,
	zrodlo NastawyPlatformy) nadajnik.Nastawy {

	if zrodlo == nil {
		return startu
	}
	wynik := startu
	nadpisz(&wynik.Host, zrodlo.Nastawa(ctx, konfig.KluczNadawcaHost))
	nadpisz(&wynik.Adres, zrodlo.Nastawa(ctx, konfig.KluczNadawcaAdres))
	nadpisz(&wynik.NazwaWyswietlana, zrodlo.Nastawa(ctx, konfig.KluczNadawcaNazwa))
	nadpisz(&wynik.Uzytkownik, zrodlo.Nastawa(ctx, konfig.KluczNadawcaUzytkow))
	nadpisz(&wynik.Sekret, zrodlo.Nastawa(ctx, konfig.KluczNadawcaSekret))
	if port, err := strconv.Atoi(zrodlo.Nastawa(ctx, konfig.KluczNadawcaPort)); err == nil && port > 0 {
		wynik.Port = port
	}
	// Odczyt trójstanowy jest jeden na cały rdzeń (`nastawy_aplikacji.go`):
	// „tak”, „nie” i „nie wskazałem”. Druga jego kopia rozjechałaby się przy
	// pierwszym dopisaniu formy zapisu.
	if szyfruj, wskazane := wartoscWymoguLogowania(zrodlo.Nastawa(ctx, konfig.KluczNadawcaStartTLS)); wskazane {
		wynik.SzyfrujStartTLS = szyfruj
	}
	return wynik
}

// nadpisz podmienia wartość wyłącznie wtedy, gdy nastawa cokolwiek mówi.
func nadpisz(cel *string, nastawa string) {
	if strings.TrimSpace(nastawa) != "" {
		*cel = strings.TrimSpace(nastawa)
	}
}
