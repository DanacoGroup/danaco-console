// Repozytorium wystawia kontrakt trwałości bytów wyposażenia modułu Terminal:
// książki hostów, biblioteki skryptów, wykazu kluczy SSH, przekierowań portów
// i obserwacji plików.
package dane

import (
	"context"

	"danacoconsole/shared"
)

// HostTerminala to wiersz tabeli `terminal_host`, reprezentujący jeden wpis
// książki hostów modułu Terminal.
type HostTerminala struct {
	Kod             string
	Nazwa           string
	Cel             string
	Port            *int64
	Grupa           string
	KatalogRoboczy  string
	KluczKod        *string
	HostPosredniKod *string
	Notatka         string
	Utworzono       string
	Zaktualizowano  string
}

// FiltrHostowTerminala zawęża wykaz zwracanych hostów książki; pole
// pozostawione puste nie zawęża niczego.
type FiltrHostowTerminala struct {
	Grupa string
	// Fraza szukana w nazwie, adresie celu i notatce, bez rozróżniania
	// wielkości liter.
	Fraza string
}

// SkryptTerminala to wiersz tabeli `terminal_skrypt`, jedna pozycja biblioteki
// skryptów w brzmieniu bieżącym.
type SkryptTerminala struct {
	Kod            string
	Nazwa          string
	Rodzaj         shared.TerminalScriptKind
	Powloka        shared.TerminalShell
	Tresc          string
	Znaczniki      string
	Alias          *string
	Wersja         int64
	Uruchomiono    *string
	Utworzono      string
	Zaktualizowano string
}

// FiltrSkryptowTerminala zawęża wykaz zwracanych pozycji biblioteki skryptów;
// pole puste nie zawęża niczego.
type FiltrSkryptowTerminala struct {
	Rodzaj shared.TerminalScriptKind
	// Znacznik szukany w rozdzielanym przecinkiem wykazie znaczników pozycji.
	Znacznik string
	// Fraza szukana w nazwie i treści.
	Fraza string
}

// KluczTerminala to wiersz tabeli `terminal_klucz`, reprezentujący jeden
// klucz SSH wykazu modułu Terminal.
type KluczTerminala struct {
	Kod        string
	Nazwa      string
	Rodzaj     shared.TerminalKeyType
	Odcisk     string
	KluczJawny string
	Sciezka    string
	Haslo      bool
	Utworzono  string
}

// TunelTerminala to wiersz tabeli `terminal_tunel`, reprezentujący jedno
// przekierowanie portu terminala.
type TunelTerminala struct {
	Kod          string
	OknoKod      string
	Rodzaj       shared.TerminalTunnelKind
	HostKod      *string
	Cel          string
	PortLokalny  *int64
	HostDocelowy string
	PortDocelowy *int64
	Stan         shared.TerminalTunnelStatus
	Powod        string
	Zalozono     string
	Zamknieto    *string
}

// FiltrTuneliTerminala zawęża wykaz zwracanych przekierowań portów terminala;
// pole puste nie zawęża niczego.
type FiltrTuneliTerminala struct {
	OknoKod string
	Stan    shared.TerminalTunnelStatus
}

// ObserwacjaTerminala to wiersz tabeli `terminal_obserwacja` — jedna obserwacja
// plików wraz z licznikiem wyzwoleń.
type ObserwacjaTerminala struct {
	Kod           string
	OknoKod       string
	KartaKod      string
	Wzorzec       string
	Polecenie     string
	Tlumienie     *int64
	Rekurencyjnie bool
	Stan          shared.TerminalWatchStatus
	Licznik       int64
	Wyzwolono     *string
	Powod         string
	Zalozono      string
}

// FiltrObserwacjiTerminala zawęża wykaz zwracanych obserwacji plików
// terminala; pole puste nie zawęża niczego.
type FiltrObserwacjiTerminala struct {
	OknoKod string
	Stan    shared.TerminalWatchStatus
}

// RepozytoriumWyposazeniaTerminala jest kontraktem trwałości wyposażenia
// modułu Terminal. Osadza go `RepozytoriumTerminala`, więc rdzeń dostaje jedno
// pole zestawu, a nie pięć.
type RepozytoriumWyposazeniaTerminala interface {
	// ── Książka hostów ──────────────────────────────────────────────────────
	ZapiszHosta(ctx context.Context, host HostTerminala) error
	Host(ctx context.Context, kod string) (HostTerminala, error)
	Hosty(ctx context.Context, filtr FiltrHostowTerminala) ([]HostTerminala, error)
	UsunHosta(ctx context.Context, kod string) (bool, error)
	// OdepnijKlucz zdejmuje wskazanie klucza z wpisów, które go używały, i oddaje
	// ich kody.
	OdepnijKlucz(ctx context.Context, kluczKod string) ([]string, error)

	// ZapiszSkrypt zakłada pozycję biblioteki albo dokłada kolejną wersję
	// i oddaje nadany numer.
	ZapiszSkrypt(ctx context.Context, skrypt SkryptTerminala) (int64, bool, error)
	Skrypt(ctx context.Context, kod string) (SkryptTerminala, error)
	Skrypty(ctx context.Context, filtr FiltrSkryptowTerminala) ([]SkryptTerminala, error)
	UsunSkrypt(ctx context.Context, kod string) (bool, error)

	// ── Wykaz kluczy SSH ────────────────────────────────────────────────────
	ZapiszKlucz(ctx context.Context, klucz KluczTerminala) error
	Klucz(ctx context.Context, kod string) (KluczTerminala, error)
	Klucze(ctx context.Context) ([]KluczTerminala, error)
	UsunKlucz(ctx context.Context, kod string) (bool, error)

	// ── Przekierowania portów ───────────────────────────────────────────────
	ZapiszTunel(ctx context.Context, tunel TunelTerminala) error
	Tunel(ctx context.Context, kod string) (TunelTerminala, error)
	Tunele(ctx context.Context, filtr FiltrTuneliTerminala) ([]TunelTerminala, error)
	ZmienStanTunelu(ctx context.Context, kod string, stan shared.TerminalTunnelStatus,
		powod string, zamkniety bool) error
	// OsierocTunele przestawia tunele stanu `active` z poprzedniego biegu
	// rdzenia na `inactive`.
	OsierocTunele(ctx context.Context) (int64, error)

	// ── Obserwacje plików ───────────────────────────────────────────────────
	ZapiszObserwacje(ctx context.Context, obserwacja ObserwacjaTerminala) error
	Obserwacja(ctx context.Context, kod string) (ObserwacjaTerminala, error)
	Obserwacje(ctx context.Context, filtr FiltrObserwacjiTerminala) ([]ObserwacjaTerminala, error)
	ZmienStanObserwacji(ctx context.Context, kod string, stan shared.TerminalWatchStatus, powod string) error
	// OdnotujWyzwolenie podnosi licznik wyzwoleń i znaczy chwilę ostatniego.
	OdnotujWyzwolenie(ctx context.Context, kod string) error
	// OsierocObserwacje przestawia obserwacje poprzedniego biegu na `stopped`.
	OsierocObserwacje(ctx context.Context) (int64, error)
}
