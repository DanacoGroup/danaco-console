// Warstwa danych obsługuje obszar modułu Diagnostics: byty czterech tabel
// diagnostycznych i kontrakt repozytorium, przenoszące wyłącznie fakty
// otrzymane od rdzenia.
package dane

import (
	"context"

	"danacoconsole/shared"
)

// WpisDiagnostyki to wiersz tabeli `diagnostyka_wpis` — jedna linia dziennika
// widoczna w oknie Logs Viewer.
type WpisDiagnostyki struct {
	Kod         string
	Chwila      int64
	Poziom      shared.LogLevel
	Zrodlo      *string
	Tresc       string
	SesjaKod    *string
	OknoKod     *string
	ProcesKod   *string
	Odcisk      string
	Powtorzenia int
}

// FiltrDziennika zawęża dziennik. Pole puste i zero nie zawężają niczego —
// zapytanie bez wskazań jest zapytaniem o cały dziennik.
type FiltrDziennika struct {
	Poziom shared.LogLevel
	Zrodlo string
	Od     int64
	Do     int64
	// Wzorzec dopasowuje się do treści wpisu; dopasowanie regularne wykonuje rdzeń, nie SQLite.
	Wzorzec string
	// Scal włącza grupowanie po odcisku wraz z licznikiem powtórzeń.
	Scal bool
	// Granica ucina wynik. Zero znaczy „granica domyślna repozytorium".
	Granica int
}

// BladDiagnostyczny to wiersz tabeli `diagnostyka_blad` — błąd zgrupowany po
// odcisku wraz z licznikiem wystąpień i kodem błędu kontraktu.
type BladDiagnostyczny struct {
	Kod         string
	Odcisk      string
	Tresc       string
	Zrodlo      *string
	KodBledu    shared.ErrorCode
	Stan        shared.DiagnosticErrorStatus
	Priorytet   shared.DiagnosticPriority
	Wystapienia int
	Notatka     *string
	Kontekst    *string
	Pierwsze    int64
	Ostatnie    int64
}

// FiltrBledow zawęża wykaz błędów diagnostycznych po stanie, priorytecie,
// zakresie czasu i granicy wyniku.
type FiltrBledow struct {
	Stan      shared.DiagnosticErrorStatus
	Priorytet shared.DiagnosticPriority
	Od        int64
	Do        int64
	Granica   int
}

// AnalizaDiagnostyczna to wiersz tabeli `diagnostyka_analiza` — migawka stanu
// systemu z chwili uruchomienia analizy.
type AnalizaDiagnostyczna struct {
	Kod          string
	OknoKod      *string
	ZakresOd     *int64
	ZakresDo     *int64
	Podsumowanie *string
	// Bledy niesie kody błędów objętych analizą, w postaci tablicy JSON.
	Bledy        *string
	PorownanaKod *string
	Utworzono    int64
}

// RekomendacjaDiagnostyczna to wiersz tabeli `diagnostyka_rekomendacja`.
// Pole BladKod wskazuje fakt, z którego rekomendacja wynika — rekomendacja bez
// wskazania faktu nie ma prawa powstać.
type RekomendacjaDiagnostyczna struct {
	Kod        string
	AnalizaKod string
	BladKod    *string
	Tytul      string
	Szczegol   *string
	Priorytet  shared.DiagnosticPriority
	Stan       shared.RecommendationStatus
	Sciezka    *string
	Poprawka   *string
	Utworzono  int64
}

// FiltrRekomendacji zawęża wykaz rekomendacji po kodzie analizy, stanie,
// priorytecie i granicy wyniku.
type FiltrRekomendacji struct {
	AnalizaKod string
	Stan       shared.RecommendationStatus
	Priorytet  shared.DiagnosticPriority
	Granica    int
}

// LicznikPoziomow zlicza wpisy dziennika w rozbiciu na poziomy — podstawa
// zagregowanego stanu pokazywanego w Diagnostics Center.
type LicznikPoziomow map[shared.LogLevel]int

// LicznikStanow zlicza błędy diagnostyczne w rozbiciu na stany, wedle stanu
// każdego zgrupowanego błędu.
type LicznikStanow map[shared.DiagnosticErrorStatus]int

// RepozytoriumDiagnostyki jest kontraktem obszaru Diagnostics: dziennik
// zdarzeń, błędy zgrupowane po odcisku, analizy i rekomendacje.
type RepozytoriumDiagnostyki interface {
	DopiszWpisy(ctx context.Context, wpisy []WpisDiagnostyki) error
	Wpisy(ctx context.Context, filtr FiltrDziennika) ([]WpisDiagnostyki, int, error)
	// ZrodlaWpisow zwraca źródła obecne w dzienniku, do budowy filtra źródła w interfejsie.
	ZrodlaWpisow(ctx context.Context) ([]string, error)
	PoziomyWpisow(ctx context.Context, od, do int64) (LicznikPoziomow, error)

	// ZapiszBlad dopisuje wystąpienie: ten sam odcisk podnosi licznik, inny zakłada nowy wiersz.
	ZapiszBlad(ctx context.Context, blad BladDiagnostyczny) (BladDiagnostyczny, error)
	Bledy(ctx context.Context, filtr FiltrBledow) ([]BladDiagnostyczny, error)
	LiczbaBledow(ctx context.Context, filtr FiltrBledow) (int, error)
	BledyPoKodach(ctx context.Context, kody []string) ([]BladDiagnostyczny, error)
	StanyBledow(ctx context.Context, od, do int64) (LicznikStanow, error)

	ZapiszAnalize(ctx context.Context, analiza AnalizaDiagnostyczna,
		rekomendacje []RekomendacjaDiagnostyczna) error
	Analiza(ctx context.Context, kod string) (AnalizaDiagnostyczna, error)
	Rekomendacje(ctx context.Context, filtr FiltrRekomendacji) ([]RekomendacjaDiagnostyczna, error)
}
