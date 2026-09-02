// Plik składa drogę odczytu konta nadawczego platformy — skrzynki, z której rdzeń pisze listy
// systemowe potwierdzenia adresu i odzyskania konta, nakładając nastawy Operatora na start.
package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/nadajnik"
)

// dwieDrogiKontaNadawczego nazywa obie drogi ustawienia skrzynki nadawczej.
// Zdanie stoi w odmowie rejestracji i odzyskania konta, bo tam odmowa zatrzymuje
// Operatora przed platformą, a odesłanie do dziennika rdzenia prowadzi tam,
// gdzie on nie sięga.
const dwieDrogiKontaNadawczego = "konto nadawcze ustawia się na dwa sposoby: " +
	"zmiennymi środowiska przy starcie serwera (DANACO_NADAWCA_HOST, DANACO_NADAWCA_ADRES " +
	"i pozostałe DANACO_NADAWCA_*) albo nastawami mailer.host i mailer.address " +
	"w oknie Konfiguracji, w kategorii „Konto nadawcze platformy”"

// NastawyPlatformy jest portem odczytu pojedynczej nastawy poziomu `aplikacja`.
// Rozszerzenie nieobowiązkowe portu Ustawienia — tą samą drogą, którą rdzeń pyta
// ten sam adapter o wymóg logowania (`NastawyAplikacji`).
type NastawyPlatformy interface {
	// Nastawa zwraca wartość zapisaną pod kluczem; napis pusty znaczy brak wskazania, nie wartość pustą.
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
	return strings.TrimSpace(a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{KontoOperatora: dane.KontoOperatora(ctx)}, klucz).Wartosc)
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
	// Odczyt trójstanowy jest jeden na cały rdzeń: tak, nie i nie wskazałem.
	if szyfruj, wskazane := wartoscWymoguLogowania(zrodlo.Nastawa(ctx, konfig.KluczNadawcaStartTLS)); wskazane {
		wynik.SzyfrujStartTLS = szyfruj
	}
	return wynik
}

// nadpisz podmienia wartość docelową wyłącznie wtedy, gdy nastawa Operatora coś wprost oznajmia o niej.
func nadpisz(cel *string, nastawa string) {
	if strings.TrimSpace(nastawa) != "" {
		*cel = strings.TrimSpace(nastawa)
	}
}
