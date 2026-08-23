// Odpowiedzialność pliku: obszar modułu Diagnostics — byty czterech tabel
// i kontrakt repozytorium. Odczyt i zapis dziennika leżą
// w `diagnostics_dziennik.go`, błędy w `diagnostics_bledy.go`, analiza wraz
// z rekomendacjami w `diagnostics_analiza.go`.
//
// Repozytorium nie wytwarza faktów. Nie liczy stanu systemu, nie ocenia
// błędów i nie wymyśla rekomendacji — przenosi wyłącznie to, co rdzeń mu podał.
// Diagnostyka, która sama zmyśla liczbę, jest gorsza niż jej brak.
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
	// Wzorzec dopasowuje się do treści wpisu. Dopasowanie regularne rdzeń
	// wykonuje sam — SQLite bez rozszerzenia nie zna operatora REGEXP,
	// a doładowywanie rozszerzenia dla jednego okna byłoby ceną bez pokrycia.
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

// FiltrBledow zawęża wykaz błędów.
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

// FiltrRekomendacji zawęża wykaz rekomendacji.
type FiltrRekomendacji struct {
	AnalizaKod string
	Stan       shared.RecommendationStatus
	Priorytet  shared.DiagnosticPriority
	Granica    int
}

// LicznikPoziomow zlicza wpisy dziennika w rozbiciu na poziomy — podstawa
// zagregowanego stanu pokazywanego w Diagnostics Center.
type LicznikPoziomow map[shared.LogLevel]int

// LicznikStanow zlicza błędy w rozbiciu na stany.
type LicznikStanow map[shared.DiagnosticErrorStatus]int

// RepozytoriumDiagnostyki jest kontraktem obszaru Diagnostics.
type RepozytoriumDiagnostyki interface {
	DopiszWpisy(ctx context.Context, wpisy []WpisDiagnostyki) error
	Wpisy(ctx context.Context, filtr FiltrDziennika) ([]WpisDiagnostyki, int, error)
	// ZrodlaWpisow zwraca źródła obecne w dzienniku — Logs Viewer buduje z nich
	// filtr źródła zamiast zgadywać nazwy.
	ZrodlaWpisow(ctx context.Context) ([]string, error)
	PoziomyWpisow(ctx context.Context, od, do int64) (LicznikPoziomow, error)

	// ZapiszBlad dopisuje wystąpienie: wiersz o tym samym odcisku podnosi
	// licznik i przesuwa chwilę ostatniego wystąpienia, nowy powstaje od zera.
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
