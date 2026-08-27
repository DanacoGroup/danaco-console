// Plik definiuje byty i kontrakt obszaru terminala: kartę powłoki oraz wpis
// dziennika procesów, wraz z ich polami i repozytorium.
package dane

import (
	"context"

	"danacoconsole/shared"
)

// KartaTerminala to wiersz tabeli terminal_karta, niosący profil powłoki
// jednej karty okna zarządzania kartami terminala.
type KartaTerminala struct {
	Kod            string
	OknoKod        string
	Powloka        shared.TerminalShell
	Tytul          *string
	KatalogRoboczy *string
	Stan           shared.TerminalSessionStatus
	Utworzono      string
	// CelZdalny to adres powłoki zdalnej albo alias konfiguracji powłoki
	// rdzenia, nie poświadczenie.
	CelZdalny string
	// PortZdalny bierze port domyślny protokołu, gdy jest pusty.
	PortZdalny *int64
	// HostKod wskazuje wpis książki hostów, z którego karta wzięła adres.
	HostKod *string
}

// ProcesTerminala to wiersz tabeli `terminal_proces` — jeden przebieg
// polecenia wraz z jego inicjatorem i wynikiem.
type ProcesTerminala struct {
	Kod          string
	KartaKod     *string
	OknoKod      string
	Pid          *int64
	PidNadrzedny *int64
	Polecenie    string
	Inicjator    shared.ProcessInitiator
	Stan         shared.TerminalProcessStatus
	KodWyjscia   *int64
	Uruchomiono  string
	Zakonczono   *string
}

// FiltrProcesow zawęża dziennik procesów. Pole puste nie zawęża niczego —
// wykaz bez wskazań jest wykazem pełnym.
type FiltrProcesow struct {
	OknoKod   string
	KartaKod  string
	Stan      shared.TerminalProcessStatus
	Inicjator shared.ProcessInitiator
}

// RepozytoriumTerminala jest kontraktem obszaru terminala, osadzającym
// również kontrakt wyposażenia modułu: książki hostów, biblioteki skryptów,
// kluczy, tuneli i obserwacji.
type RepozytoriumTerminala interface {
	RepozytoriumWyposazeniaTerminala

	ZapiszKarte(ctx context.Context, karta KartaTerminala) error
	Karta(ctx context.Context, kod string) (KartaTerminala, error)
	Karty(ctx context.Context) ([]KartaTerminala, error)

	ZapiszProces(ctx context.Context, proces ProcesTerminala) error
	ZakonczProces(ctx context.Context, kod string, stan shared.TerminalProcessStatus, kodWyjscia *int64) error
	Proces(ctx context.Context, kod string) (ProcesTerminala, error)
	Procesy(ctx context.Context, filtr FiltrProcesow) ([]ProcesTerminala, error)
	// OsierociProcesy przestawia procesy działające z poprzedniego biegu
	// rdzenia na zakończone.
	OsierociProcesy(ctx context.Context) (int64, error)
}
