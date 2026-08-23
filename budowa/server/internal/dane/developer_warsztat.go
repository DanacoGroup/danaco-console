// Odpowiedzialność pliku: byty warsztatu modułu Developer, które przeżywają
// pojedynczą sesję — punkty przerwania, kolekcje zapytań API, połączenia
// bazodanowe, przebiegi skanowania wraz ze znaleziskami oraz wyniki testów
// i pokrycie kodu przypisane do przebiegu budowania.
//
// Dlaczego to nie leży w pamięci rdzenia. Sesja debugowania, połączenie z bazą
// i biegnący skan to stan ŻYWY — ma uchwyty do procesów i gaśnie razem z nimi;
// takich rzeczy w bazie nie ma i mieć nie powinno. Tutaj leży to, co Operator
// ułożył i czego nie ma prawa stracić przy zamknięciu okna: postawiony punkt
// przerwania, zapisane zapytanie, opisane połączenie, wynik pomiaru.
//
// Odczyt leży w `developer_warsztat_odczyt.go`, zapis w
// `developer_warsztat_zapis.go`; kontrakt obszaru — w `developer.go`.
package dane

import "danacoconsole/shared"

// PunktPrzerwania to wiersz tabeli `developer_punkt_przerwania` — punkt
// postawiony w marginesie Code Editora, należący do okna i pliku, nie do sesji
// debugowania.
type PunktPrzerwania struct {
	Kod            string
	OknoKod        string
	Sciezka        string
	Wiersz         int64
	Rodzaj         shared.BreakpointKind
	Warunek        *string
	WarunekTrafien *string
	Wpis           *string
	Zweryfikowany  bool
}

// KolekcjaApi to wiersz tabeli `developer_kolekcja_api` — zestaw zapytań HTTP
// wraz ze środowiskami, oba w postaci tekstu JSON kontraktu.
type KolekcjaApi struct {
	Kod        string
	OknoKod    string
	Nazwa      string
	Zapytania  string
	Srodowiska *string
	Zmieniono  string
}

// PolaczenieDanych to wiersz tabeli `developer_polaczenie_danych`. Wiersz
// opisuje, do czego się łączyć; hasło stoi w sejfie pod odwołaniem z pola
// `Poswiadczenie`.
type PolaczenieDanych struct {
	Kod           string
	OknoKod       string
	Nazwa         string
	Silnik        shared.DataEngine
	Host          *string
	Port          *int64
	Baza          string
	Uzytkownik    *string
	Poswiadczenie *string
	TylkoOdczyt   bool
}

// PrzebiegSkanu to wiersz tabeli `developer_skan` — nagłówek jednego pomiaru
// bezpieczeństwa i jakości repozytorium.
type PrzebiegSkanu struct {
	Kod         string
	OknoKod     string
	Rodzaje     string
	Stan        shared.BuildStatus
	Znalezisk   *int64
	Uruchomiono string
	Zakonczono  *string
}

// ZnaleziskoSkanu to wiersz tabeli `developer_znalezisko` — jedno spostrzeżenie
// przebiegu skanowania.
type ZnaleziskoSkanu struct {
	Kod           string
	SkanKod       string
	Rodzaj        shared.ScanKind
	Waga          shared.ProblemSeverity
	Tytul         string
	Opis          *string
	Sciezka       *string
	Wiersz        *int64
	Regula        *string
	Cve           *string
	Pakiet        *string
	WersjaNaprawy *string
}

// WynikTestu to wiersz tabeli `developer_wynik_testu` — jeden test rozpoznany
// w wyjściu przebiegu budowania.
type WynikTestu struct {
	BudowanieKod string
	Zestaw       *string
	Nazwa        string
	Stan         shared.TestStatus
	CzasMs       *int64
	Tresc        *string
	Sciezka      *string
	Wiersz       *int64
}

// PokryciePliku to wiersz tabeli `developer_pokrycie` — pomiar pokrycia jednego
// pliku w jednym przebiegu.
type PokryciePliku struct {
	BudowanieKod       string
	Sciezka            string
	Instrukcje         int64
	Pokryte            int64
	Procent            int64
	WierszeBezPokrycia *string
}

// FiltrZnalezisk zawęża wykaz znalezisk. Pole puste znaczy „bez zawężenia”.
type FiltrZnalezisk struct {
	SkanKod string
	OknoKod string
	Rodzaj  string
	Waga    string
	Limit   int
}
