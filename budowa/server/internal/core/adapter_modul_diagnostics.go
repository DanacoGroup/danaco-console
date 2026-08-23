// Odpowiedzialność pliku: moduł Diagnostics — wypełnienie portu Diagnostyka
// komendami obszaru `diagnostics.*` oraz podłączenie źródeł faktów. Dziennik
// leży w `adapter_modul_diagnostics_dziennik.go`, błędy w
// `adapter_modul_diagnostics_bledy.go`, analiza wraz z rekomendacjami
// w `adapter_modul_diagnostics_analiza.go`, przekład wierszy w
// `adapter_modul_diagnostics_przeklad.go`.
//
// Moduł ma dwa źródła faktów:
//  1. Dziennik rdzenia — adapter jest odbiorcą wyjścia `*log.Logger` rdzenia,
//     więc każda linia, którą rdzeń zapisuje o sobie, staje się wpisem
//     widocznym w Logs Viewer.
//  2. Odmowy wykonania komend — dyspozytor oddaje adapterowi każdą odpowiedź
//     błędną wraz z kodem ze słownika ErrorCode kontraktu, więc Errors Panel
//     pokazuje błędy rzeczywistych tur i operacji.
//
// Adapter nie liczy obciążenia maszyny, nie sprawdza usług zewnętrznych i nie
// ocenia stanu zdrowia systemu — rdzeń takich faktów nie wystawia. Okno pokazuje
// w tych miejscach brak danych.
package core

import (
	"context"
	"sync"
	"sync/atomic"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Zgodność adaptera z portem sprawdzana jest przy kompilacji.
var _ Diagnostyka = (*adapterDiagnostyki)(nil)

// Przedrostki identyfikatorów bytów modułu.
const (
	przedrostekWpisuDziennika = "log-"
	przedrostekBleduDiagnozy  = "blad-"
	przedrostekAnalizy        = "analiza-"
	przedrostekRekomendacji   = "rekom-"
)

// pojemnoscKolejkiWpisow ogranicza zaległość zapisu dziennika. Wartość jest
// kompromisem: bufor za mały gubiłby wpisy przy nagłym wysypie, za duży
// trzymałby w pamięci dziennik, którego i tak nikt nie zdąży odczytać.
const pojemnoscKolejkiWpisow = 1024

// adapterDiagnostyki wypełnia port Diagnostyka.
type adapterDiagnostyki struct {
	repozytorium dane.RepozytoriumDiagnostyki

	// wpisy przyjmuje linie dziennika bez blokowania piszącego. Zapis do bazy
	// idzie osobną goroutine, bo `log.Logger` trzyma przy zapisie własną
	// blokadę — czekanie na dysk pod tą blokadą wstrzymywałoby cały rdzeń.
	wpisy chan dane.WpisDiagnostyki
	// odrzucone liczy wpisy, których kolejka nie przyjęła. Liczba wychodzi do
	// podsumowania analizy, żeby strata dziennika nie została przemilczana.
	odrzucone atomic.Int64
	// niezapisane liczy wpisy, których nie przyjęła baza. Odrębne od
	// odrzuconych, bo przyczyna i zalecane działanie są inne.
	niezapisane atomic.Int64

	zamkniecie sync.Once
	koniec     chan struct{}

	// przedrostek i flagi opisują format dziennika rdzenia. Bez nich nie da się
	// dokładnie zdjąć z linii nagłówka, który `log.Logger` sam dołożył.
	przedrostek string
	flagi       int

	// zmiana rozgłasza `diagnostics.analysis.changed`. Podpina ją obsługiwacz.
	zmiana func(shared.ChangeKind, shared.DiagnosticAnalysis)
}

// nowyAdapterDiagnostyki wiąże port z repozytorium modułu i uruchamia pisarza
// dziennika. Bez repozytorium moduł nie ma czego przeszukiwać i komendy
// odmawiają, zamiast oddawać pusty wykaz.
func nowyAdapterDiagnostyki(ctx context.Context,
	repozytorium dane.RepozytoriumDiagnostyki) *adapterDiagnostyki {

	adapter := &adapterDiagnostyki{
		repozytorium: repozytorium,
		wpisy:        make(chan dane.WpisDiagnostyki, pojemnoscKolejkiWpisow),
		koniec:       make(chan struct{}),
	}
	go adapter.pisz(ctx)
	return adapter
}

// PodepnijRozgloszenie wypełnia port: adapter zapamiętuje drogę do zdarzenia.
func (a *adapterDiagnostyki) PodepnijRozgloszenie(rozglos func(shared.ChangeKind, shared.DiagnosticAnalysis)) {
	a.zmiana = rozglos
}

// Zamknij domyka pisarza dziennika. Wpisy już zakolejkowane zostają zapisane,
// bo ostatnie linie przed zatrzymaniem rdzenia niosą najwięcej dla diagnozy.
func (a *adapterDiagnostyki) Zamknij() {
	if a == nil {
		return
	}
	a.zamkniecie.Do(func() { close(a.wpisy) })
	<-a.koniec
}

// bladDiagnostyki nazywa usterkę warstwy danych modułu.
func bladDiagnostyki(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladBrakuTrwalosci odmawia wykonania komendy, gdy moduł nie ma repozytorium.
// Pusty wykaz znaczyłby, że nic się nie wydarzyło, a tego adapter bez dziennika
// nie ma jak stwierdzić.
func bladBrakuTrwalosci(czynnosc string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Diagnostics: "+czynnosc+" wymaga dziennika, a rdzeń nie ma podłączonego repozytorium"))
}

// bladWskazaniaDiagnostyki nazywa brak danych w żądaniu — to błąd wywołującego,
// nie usterka rdzenia, więc kod odmowy jest inny.
func bladWskazaniaDiagnostyki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Diagnostics: "+powod))
}
