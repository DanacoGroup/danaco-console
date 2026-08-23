// Odpowiedzialność pliku: obszar modułu Terminal — karta powłoki (tabela
// `terminal_karta`) i dziennik procesów (tabela `terminal_proces`) — byty
// obszaru i jego kontrakt. Odczyt leży w `terminal_odczyt.go`,
// zapis w `terminal_zapis.go`.
//
// Repozytorium nie prowadzi procesów. Uchwyt do biegnącego procesu i jego
// drzewa potomstwa ma wyłącznie rdzeń — tylko on potrafi proces ubić. Tutaj
// zapisuje się to, co po procesie zostaje: polecenie, inicjator, kod wyjścia
// i czasy. Dzięki temu Process Monitor pokazuje także procesy zakończone,
// których stan żywy rdzenia już nie trzyma.
package dane

import (
	"context"

	"danacoconsole/shared"
)

// KartaTerminala to wiersz tabeli `terminal_karta` — profil powłoki jednej
// karty okna Terminal Tabs.
type KartaTerminala struct {
	Kod            string
	OknoKod        string
	Powloka        shared.TerminalShell
	Tytul          *string
	KatalogRoboczy *string
	Stan           shared.TerminalSessionStatus
	Utworzono      string
	// CelZdalny to adres powłoki zdalnej w postaci `użytkownik@host` albo alias
	// konfiguracji OpenSSH maszyny rdzenia. Poświadczeniem nie jest — jest tym
	// samym, co widnieje w wykazie książki hostów — więc ma kolumnę, inaczej niż
	// zmienne środowiska karty (migracja 251).
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

// RepozytoriumTerminala jest kontraktem obszaru Terminal. Osadza kontrakt
// wyposażenia modułu (`terminal_wyposazenie.go`) — książki hostów, biblioteki
// skryptów, wykazu kluczy, tuneli i obserwacji — bo wszystkie te byty należą do
// jednego modułu i jednego adaptera rdzenia.
type RepozytoriumTerminala interface {
	RepozytoriumWyposazeniaTerminala

	ZapiszKarte(ctx context.Context, karta KartaTerminala) error
	Karta(ctx context.Context, kod string) (KartaTerminala, error)
	Karty(ctx context.Context) ([]KartaTerminala, error)

	ZapiszProces(ctx context.Context, proces ProcesTerminala) error
	ZakonczProces(ctx context.Context, kod string, stan shared.TerminalProcessStatus, kodWyjscia *int64) error
	Proces(ctx context.Context, kod string) (ProcesTerminala, error)
	Procesy(ctx context.Context, filtr FiltrProcesow) ([]ProcesTerminala, error)
	// OsierociProcesy przestawia procesy zostawione w stanie `running` przez
	// poprzedni bieg rdzenia na `stopped`. Rdzeń po restarcie nie ma do nich
	// uchwytu, więc wykazywanie ich jako czynnych byłoby nieprawdą.
	OsierociProcesy(ctx context.Context) (int64, error)
}
