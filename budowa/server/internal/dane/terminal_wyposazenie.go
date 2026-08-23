// Odpowiedzialność pliku: byty wyposażenia modułu Terminal — książka hostów,
// biblioteka skryptów, wykaz kluczy SSH, przekierowania portów i obserwacje
// plików — wraz z ich kontraktem repozytoryjnym. Odczyt leży
// w `terminal_wyposazenie_odczyt.go`, zapis w `terminal_wyposazenie_zapis.go`.
//
// Dlaczego jeden kontrakt, a nie pięć. Wszystkie te byty należą do jednego
// modułu i jednego adaptera rdzenia; rozbicie ich na pięć interfejsów dałoby
// pięć pól w zestawie repozytoriów i pięć podpięć w montażu, a ani jednej
// nowej granicy. Granica jest tu modułowa — Terminal — i tak ją prowadzimy.
//
// Czego w tych bytach NIE MA: materiału tajnego. Wpis hosta niesie adres i
// wskazanie klucza, wpis klucza — ścieżkę i odcisk. Hasło i klucz prywatny
// zostają na dysku maszyny rdzenia; baza nie jest sejfem.
package dane

import (
	"context"

	"danacoconsole/shared"
)

// HostTerminala to wiersz tabeli `terminal_host` — jeden wpis książki hostów.
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

// FiltrHostowTerminala zawęża książkę hostów. Pole puste nie zawęża niczego.
type FiltrHostowTerminala struct {
	Grupa string
	// Fraza szukana w nazwie, adresie celu i notatce; porównanie bez
	// rozróżniania wielkości liter leży po stronie zapytania.
	Fraza string
}

// SkryptTerminala to wiersz tabeli `terminal_skrypt` — jedna pozycja biblioteki
// w brzmieniu bieżącym.
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

// FiltrSkryptowTerminala zawęża bibliotekę skryptów.
type FiltrSkryptowTerminala struct {
	Rodzaj shared.TerminalScriptKind
	// Znacznik szukany w rozdzielanym przecinkiem wykazie znaczników pozycji.
	Znacznik string
	// Fraza szukana w nazwie i treści.
	Fraza string
}

// KluczTerminala to wiersz tabeli `terminal_klucz` — jeden klucz SSH wykazu.
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

// TunelTerminala to wiersz tabeli `terminal_tunel` — jedno przekierowanie portu.
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

// FiltrTuneliTerminala zawęża wykaz przekierowań portów.
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

// FiltrObserwacjiTerminala zawęża wykaz obserwacji plików.
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
	// OdepnijKlucz zdejmuje wskazanie klucza z wpisów, które go używały,
	// i oddaje ich kody. Wpisy wracają wtedy do klucza domyślnego konfiguracji
	// maszyny rdzenia — kontrakt `terminal.key.remove` wymaga ich wymienienia.
	OdepnijKlucz(ctx context.Context, kluczKod string) ([]string, error)

	// ── Biblioteka skryptów ─────────────────────────────────────────────────
	// ZapiszSkrypt zakłada pozycję albo dokłada jej kolejną wersję. Zwraca
	// nadany numer wersji oraz prawdę, gdy pozycja powstała. Numeru wersji nie
	// przyjmuje od wołającego: nadaje go rdzeń w jednej transakcji z zapisem,
	// żeby dwa równoległe zapisy nie dostały tego samego numeru.
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
	// OsierocTunele przestawia tunele zostawione w stanie `active` przez
	// poprzedni bieg rdzenia na `inactive`; po restarcie żaden z nich nie
	// przenosi już ani jednego bajtu.
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
