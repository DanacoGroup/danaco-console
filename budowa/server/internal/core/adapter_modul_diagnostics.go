// Odpowiedzialność pliku: moduł Diagnostics — wypełnienie portu Diagnostyka komendami obszaru diagnostics.* oraz podłączenie źródeł faktów: dziennik rdzenia i odmowy wykonania komend.
package core

import (
	"context"
	"sync"
	"sync/atomic"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

var _ Diagnostyka = (*adapterDiagnostyki)(nil)

// Przedrostki identyfikatorów bytów modułu Diagnostics, po jednym na każdy rodzaj bytu trwałego tego dziennika.
const (
	przedrostekWpisuDziennika = "log-"
	przedrostekBleduDiagnozy  = "blad-"
	przedrostekAnalizy        = "analiza-"
	przedrostekRekomendacji   = "rekom-"
)

// pojemnoscKolejkiWpisow ogranicza zaległość zapisu dziennika. Wartość jest kompromisem: bufor za mały gubiłby wpisy przy nagłym wysypie, za duży trzymałby w pamięci dziennik, którego i tak nikt nie zdąży odczytać.
const pojemnoscKolejkiWpisow = 1024

// adapterDiagnostyki wypełnia port Diagnostyka i trzyma stan pisarza dziennika rdzenia tej platformy Danaco.
type adapterDiagnostyki struct {
	repozytorium dane.RepozytoriumDiagnostyki

	// wpisy przyjmuje linie dziennika bez blokowania piszącego, zapis do bazy idzie osobną goroutine.
	wpisy chan dane.WpisDiagnostyki
	// odrzucone liczy wpisy, których kolejka nie przyjęła, do podsumowania analizy.
	odrzucone atomic.Int64
	// niezapisane liczy wpisy, których nie przyjęła baza, osobno od odrzuconych.
	niezapisane atomic.Int64

	zamkniecie sync.Once
	koniec     chan struct{}

	// przedrostek i flagi opisują format dziennika, potrzebny do zdjęcia nagłówka linii.
	przedrostek string
	flagi       int

	// zmiana rozgłasza diagnostics.analysis.changed do klientów; podpina ją obsługiwacz komend.
	zmiana func(context.Context, shared.ChangeKind, shared.DiagnosticAnalysis)
}

// nowyAdapterDiagnostyki wiąże port z repozytorium modułu i uruchamia pisarza dziennika. Bez repozytorium moduł nie ma czego przeszukiwać i komendy odmawiają, zamiast oddawać pusty wykaz.
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

// PodepnijRozgloszenie wypełnia port: adapter zapamiętuje drogę do rozgłoszenia zdarzenia zmiany analizy.
func (a *adapterDiagnostyki) PodepnijRozgloszenie(rozglos func(context.Context, shared.ChangeKind, shared.DiagnosticAnalysis)) {
	a.zmiana = rozglos
}

// Zamknij domyka pisarza dziennika. Wpisy już zakolejkowane zostają zapisane, bo ostatnie linie przed zatrzymaniem rdzenia niosą najwięcej dla diagnozy.
func (a *adapterDiagnostyki) Zamknij() {
	if a == nil {
		return
	}
	a.zamkniecie.Do(func() { close(a.wpisy) })
	<-a.koniec
}

// bladDiagnostyki nazywa usterkę warstwy danych modułu kodem błędu wewnętrznego kontraktu tej platformy.
func bladDiagnostyki(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladBrakuTrwalosci odmawia wykonania komendy, gdy moduł nie ma repozytorium. Pusty wykaz znaczyłby, że nic się nie wydarzyło, a tego adapter bez dziennika nie ma jak stwierdzić.
func bladBrakuTrwalosci(czynnosc string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Diagnostics: "+czynnosc+" wymaga dziennika, a serwer nie ma podłączonego repozytorium"))
}

// bladWskazaniaDiagnostyki nazywa brak danych w żądaniu — to błąd wywołującego, nie usterka rdzenia, więc kod odmowy jest inny.
func bladWskazaniaDiagnostyki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Diagnostics: "+powod))
}
