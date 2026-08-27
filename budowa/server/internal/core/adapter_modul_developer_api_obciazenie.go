// Odpowiedzialność pliku: przebieg obciążeniowy punktu końcowego —
// `developer.api.load.run`.
//
// ── Jedno żądanie a kształt wyniku ───────────────────────────────────────────
// `developer.api.request` strzela JEDNYM żądaniem i oddaje jego status, czas
// i rozmiar. To wystarcza, żeby sprawdzić, czy punkt końcowy odpowiada i co
// odpowiada; nie wystarcza, żeby powiedzieć o nim cokolwiek pod obciążeniem.
// Jeden pomiar nie ma percentyla ani przepustowości — a właśnie ogon rozkładu,
// nie średnia, rozstrzyga o tym, czy usługa jest do użycia.
//
// ── Dlaczego program, a nie pętla po `net/http` ──────────────────────────────
// Rzetelny przebieg obciążeniowy to nie pętla wywołań: trzeba utrzymać zadaną
// liczbę połączeń równolegle, zbierać histogram czasów bez wpływu na pomiar
// i policzyć percentyle z pełnego rozkładu, nie z próbki. Napisane od nowa
// mierzyłoby w dużej mierze samo siebie.
//
// ── Wybór: autocannon, nie k6 ────────────────────────────────────────────────
// Oba programy stoją na maszynie i oba umieją percentyle. Rozstrzygnęły trzy
// rzeczy:
//
//  1. K6 opisuje przebieg SKRYPTEM w JavaScripcie, a nie parametrami. Wpięcie
//     go tutaj znaczyłoby albo kontrakt niosący program do wykonania — inna
//     i znacznie szersza powierzchnia niż adres z parametrami — albo skrypt
//     składany przez rdzeń, czyli generowanie cudzego języka.
//  2. Jedyna droga rdzenia do procesu (`zewnetrzne.Wolaj`) zbiera WYJŚCIE
//     programu. Autocannon oddaje cały wynik na wyjście (`-j`); k6 pisze
//     podsumowanie maszynowe do PLIKU, więc wymagałby pisania i odczytu plików
//     pośrednich, których ta droga nie obsługuje.
//  3. Kształt, o który pyta kontrakt — percentyle czasu, żądania na sekundę,
//     bajty na sekundę, rozbicie po kodach stanu — autocannon oddaje wprost.
//
// Siłą k6 są przebiegi narastające, progi i scenariusze. Są nieosiągalne bez
// skryptu, a skryptu nikt tu nie zamawiał.
//
// ── Zero żądań nie jest wynikiem zerowym ─────────────────────────────────────
// Program kończy się powodzeniem także wtedy, gdy ani jedno żądanie nie doszło
// do skutku: punkt końcowy milczy, a wynik niesie same zera obok licznika
// błędów. Podanie takiego wyniku jako pomiaru byłoby brakiem pomiaru w przebraniu
// — „zero żądań na sekundę" czyta się jak usługa skrajnie wolna, a nie jak
// usługa, której nie ma. Dlatego przebieg bez ani jednej odpowiedzi wraca
// odmową nazywającą liczbę błędów.
//
// Odpowiedzi spoza klasy 2xx to co innego: punkt końcowy odpowiedział, więc
// pomiar SIĘ ODBYŁ. Wynik wychodzi wraz z licznikiem `non2xx` i rozbiciem po
// kodach — wołający ma zobaczyć, że mierzył ścieżkę błędu, a nie zgadywać.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// narzedzieAutocannon opisuje program przebiegu obciążeniowego. Deklaracja stoi
// przy miejscu użycia; wykaz zależności odwołuje się do niej, zamiast powtarzać
// nazwę po raz drugi.
var narzedzieAutocannon = zewnetrzne.Narzedzie{
	Nazwa: "autocannon", Program: "autocannon", Pakiet: "npm i -g autocannon"}

const (
	// polaczeniaPrzebieguDomyslne jest liczbą połączeń równoległych przy braku
	// wskazania. Dziesięć, bo tyle wystarcza, żeby rozkład czasów miał ogon,
	// a jednocześnie nie zajmuje maszyny rdzenia samym generowaniem ruchu.
	polaczeniaPrzebieguDomyslne = 10
	// najwiecejPolaczenPrzebiegu jest granicą, której żądanie nie przekroczy.
	// Powyżej niej wąskim gardłem przestaje być mierzona usługa, a staje się
	// maszyna rdzenia — i pomiar zaczyna mierzyć siebie.
	najwiecejPolaczenPrzebiegu = 1000
	// czasPrzebieguDomyslny jest czasem trwania przy braku wskazania.
	czasPrzebieguDomyslny = 10
	// najdluzszyPrzebieg jest granicą czasu trwania w sekundach. Przebieg
	// półgodzinny trzyma połączenie klienta i obciąża cudzą usługę dłużej,
	// niż ktokolwiek na wynik czeka.
	najdluzszyPrzebieg = 600
	// granicaZapasuPrzebiegu jest zapasem, o który granica arsenału przewyższa
	// zamówiony czas trwania: program potrzebuje chwili na start, domknięcie
	// połączeń i policzenie rozkładu PO upływie czasu przebiegu.
	granicaZapasuPrzebiegu = 60 * time.Second
)

// WykonajPrzebiegObciazeniowy obsługuje `developer.api.load.run`.
func (a *adapterDevelopera) WykonajPrzebiegObciazeniowy(ctx context.Context,
	z shared.DeveloperApiLoadRunRequest) (shared.DeveloperApiLoadRunResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperApiLoadRunResponse{}, err
	}
	// Przebieg obciążeniowy wysyła tysiące żądań zmieniających stan po drugiej
	// stronie sieci — jest tym samym, czym pojedyncze zapytanie, tylko wielokrotnie.
	// Tryb planistyczny wyklucza zmiany w systemie, więc wyklucza i te.
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "przebieg obciążeniowy punktu końcowego"); err != nil {
		return shared.DeveloperApiLoadRunResponse{}, err
	}
	if a.uruchamiacz == nil {
		return shared.DeveloperApiLoadRunResponse{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Developer: rdzeń nie ma uruchamiacza procesów — przebieg obciążeniowy "+
					"nie ma czym wystartować; naprawa: podpiąć warstwę kanału (injection) "+
					"przy składaniu rdzenia"))
	}

	podstawienia, err := a.podstawieniaSrodowiska(ctx, okno.Id, z.EnvironmentId)
	if err != nil {
		return shared.DeveloperApiLoadRunResponse{}, err
	}
	adres := podstawWSzablonie(strings.TrimSpace(z.Url), podstawienia)
	if adres == "" {
		return shared.DeveloperApiLoadRunResponse{}, bladZadaniaDevelopera(
			"przebieg obciążeniowy wymaga adresu")
	}
	if !strings.HasPrefix(adres, "http://") && !strings.HasPrefix(adres, "https://") {
		return shared.DeveloperApiLoadRunResponse{}, bladZadaniaDevelopera(
			"adres punktu końcowego ma zaczynać się od http:// albo https://; otrzymano " + adres)
	}
	metoda := strings.ToUpper(strings.TrimSpace(wartoscTekstuLubPusta(z.Method)))
	if metoda == "" {
		metoda = http.MethodGet
	}

	polaczenia := wGranicach(wartoscLiczbyLub(z.Connections, polaczeniaPrzebieguDomyslne),
		1, najwiecejPolaczenPrzebiegu)
	sekundy := wGranicach(wartoscLiczbyLub(z.DurationSeconds, czasPrzebieguDomyslny),
		1, najdluzszyPrzebieg)

	argumenty := []string{
		"--connections", strconv.Itoa(polaczenia),
		"--duration", strconv.Itoa(sekundy),
		"--method", metoda,
		// Bez tego wiersza program rysuje pasek postępu na diagnostyce przez cały
		// przebieg — nikt go tu nie ogląda, a bufor potoku ma swoją pojemność.
		"--no-progress",
		// Wynik maszynowy na wyjście — jedyny kształt, który ta droga zbiera.
		"--json",
	}
	for _, naglowek := range naglowkiPrzebiegu(z.Headers, podstawienia) {
		argumenty = append(argumenty, "--headers", naglowek)
	}
	if z.Body != nil && *z.Body != "" {
		argumenty = append(argumenty, "--body", podstawWSzablonie(*z.Body, podstawienia))
	}
	argumenty = append(argumenty, adres)

	granica := time.Duration(sekundy)*time.Second + granicaZapasuPrzebiegu
	obszar := a.obszarDevelopera(okno)
	poczatek := time.Now()
	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieAutocannon, argumenty,
		obszar.KatalogRoboczy, granica)
	trwanie := time.Since(poczatek)
	if err != nil {
		return shared.DeveloperApiLoadRunResponse{}, bladProgramuDevelopera(
			"developer.api.load.run", err, trwanie, granica)
	}

	var odczyt wynikPrzebieguObciazeniowego
	if err := json.Unmarshal(wynik.Wyjscie, &odczyt); err != nil {
		return shared.DeveloperApiLoadRunResponse{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Developer: program przebiegu obciążeniowego nie oddał wyniku dla "+adres+
					", więc przebieg się nie odbył: "+err.Error()))
	}
	if odczyt.Zadania.Razem <= 0 {
		return shared.DeveloperApiLoadRunResponse{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Developer: przebieg obciążeniowy punktu końcowego "+adres+
					" nie zebrał ani jednej odpowiedzi w ciągu "+strconv.Itoa(sekundy)+
					" s (błędów: "+strconv.FormatInt(odczyt.Bledy, 10)+", przekroczeń czasu: "+
					strconv.FormatInt(odczyt.Przekroczenia, 10)+"), więc nie ma czego zmierzyć. "+
					"Wynik zerowy byłby tu odpowiedzią o usłudze, do której nikt się nie dodzwonił"))
	}

	przebieg := odczyt.jakoPrzebieg(adres, metoda, polaczenia, poczatek)
	// Wersja programu wchodzi do wyniku, bo rozkład czasów jest orzeczeniem
	// konkretnego wydania generatora ruchu — wynik bez wersji nie daje się
	// porównać z wynikiem sprzed miesiąca.
	if sciezka, jest := zewnetrzne.Odnajdz(narzedzieAutocannon); jest {
		przebieg.ToolVersion = wersjaProgramuWarsztatu(sciezka)
	}
	return shared.DeveloperApiLoadRunResponse{Run: przebieg}, nil
}

// naglowkiPrzebiegu przekłada nagłówki żądania na argumenty programu wraz
// z podstawieniami środowiska. Kolejność jest ustalona, bo mapa Go oddaje wpisy
// losowo, a dwa przebiegi tej samej kolekcji mają iść tym samym wierszem poleceń.
func naglowkiPrzebiegu(surowe json.RawMessage, podstawienia map[string]string) []string {
	if len(surowe) == 0 {
		return nil
	}
	wpisy := map[string]string{}
	if err := json.Unmarshal(surowe, &wpisy); err != nil {
		return nil
	}
	nazwy := make([]string, 0, len(wpisy))
	for nazwa := range wpisy {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	wykaz := make([]string, 0, len(nazwy))
	for _, nazwa := range nazwy {
		wykaz = append(wykaz, nazwa+"="+podstawWSzablonie(wpisy[nazwa], podstawienia))
	}
	return wykaz
}

// wynikPrzebieguObciazeniowego jest tą częścią odpowiedzi programu, którą moduł
// czyta. Program oddaje też rozkłady żądań i przepustowości po percentylach;
// kontrakt niesie z nich wartość średnią, bo percentyl liczby żądań na sekundę
// mówi o próbkowaniu sekundowym, a nie o usłudze.
type wynikPrzebieguObciazeniowego struct {
	CzasTrwania   float64 `json:"duration"`
	Bledy         int64   `json:"errors"`
	Przekroczenia int64   `json:"timeouts"`
	Niepowodzenia int64   `json:"non2xx"`
	Opoznienie    struct {
		Srednia   float64 `json:"average"`
		Odchylnie float64 `json:"stddev"`
		Min       float64 `json:"min"`
		Maks      float64 `json:"max"`
		P50       float64 `json:"p50"`
		P90       float64 `json:"p90"`
		P99       float64 `json:"p99"`
	} `json:"latency"`
	Zadania struct {
		Srednia float64 `json:"average"`
		Razem   int64   `json:"total"`
	} `json:"requests"`
	Przepustowosc struct {
		Srednia float64 `json:"average"`
	} `json:"throughput"`
	Stany map[string]struct {
		Ile int64 `json:"count"`
	} `json:"statusCodeStats"`
}

// jakoPrzebieg składa wynik kontraktu z odpowiedzi programu.
func (w wynikPrzebieguObciazeniowego) jakoPrzebieg(adres, metoda string, polaczenia int,
	poczatek time.Time) shared.ApiLoadRun {

	return shared.ApiLoadRun{
		Url:               adres,
		Method:            metoda,
		Connections:       polaczenia,
		DurationSeconds:   w.CzasTrwania,
		RequestsTotal:     w.Zadania.Razem,
		RequestsPerSecond: w.Zadania.Srednia,
		BytesPerSecond:    w.Przepustowosc.Srednia,
		Latency: shared.ApiLoadLatency{
			AverageMs: w.Opoznienie.Srednia,
			StddevMs:  w.Opoznienie.Odchylnie,
			MinMs:     w.Opoznienie.Min,
			MaxMs:     w.Opoznienie.Maks,
			P50Ms:     w.Opoznienie.P50,
			P90Ms:     w.Opoznienie.P90,
			P99Ms:     w.Opoznienie.P99,
		},
		StatusCounts: rozbicieStanow(w.Stany),
		Non2xx:       w.Niepowodzenia,
		Errors:       w.Bledy,
		Timeouts:     w.Przekroczenia,
		StartedAt:    poczatek.UnixMilli(),
		FinishedAt:   time.Now().UnixMilli(),
	}
}

// rozbicieStanow porządkuje wykaz kodów stanu rosnąco. Kolejność jest ustalona,
// bo mapa Go oddaje wpisy losowo, a dwa przebiegi tej samej usługi mają dawać
// wykaz porównywalny wprost.
func rozbicieStanow(stany map[string]struct {
	Ile int64 `json:"count"`
}) []shared.ApiLoadStatusCount {

	wykaz := make([]shared.ApiLoadStatusCount, 0, len(stany))
	for zapis, wpis := range stany {
		kod, err := strconv.Atoi(strings.TrimSpace(zapis))
		if err != nil {
			// Klucz, który nie jest kodem stanu, pomijamy zamiast wpisywać zero:
			// zero nie jest kodem stanu i wykaz kłamałby o odpowiedzi.
			continue
		}
		wykaz = append(wykaz, shared.ApiLoadStatusCount{Status: kod, Count: wpis.Ile})
	}
	sort.Slice(wykaz, func(i, j int) bool { return wykaz[i].Status < wykaz[j].Status })
	return wykaz
}

// wGranicach sprowadza wartość do przedziału dopuszczalnego.
func wGranicach(wartosc, najmniej, najwiecej int) int {
	if wartosc < najmniej {
		return najmniej
	}
	if wartosc > najwiecej {
		return najwiecej
	}
	return wartosc
}

// bladProgramuDevelopera odróżnia brak programu i przekroczenie granicy czasu od
// usterki rdzenia — tak samo jak robi to moduł Apps przy audycie wydajności.
//
// Przekroczenie rozstrzyga się ZMIERZONYM czasem, nie treścią komunikatu:
// zdanie, którym arsenał opisuje przerwanie, jest napisem i przy następnej
// zmianie brzmiałoby inaczej.
func bladProgramuDevelopera(komenda string, err error, trwanie, granica time.Duration) error {
	if err == nil {
		return nil
	}
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Developer: komenda "+komenda+" nie ma czym obciążyć punktu końcowego: "+
				err.Error()))
	}
	if granica > 0 && trwanie >= granica {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Developer: komenda "+komenda+" przekroczyła granicę czasu "+granica.String()+
				" — to jest przekroczenie granicy, nie usterka rdzenia; naprawa: skrócić czas "+
				"trwania przebiegu polem durationSeconds. Diagnostyka warstwy: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Developer: komenda "+komenda+" nie doszła do skutku: "+err.Error()))
}
