package transport

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/konfiguracja"
)

const (
	// SciezkaGniazdaDomyslna to ścieżka HTTP, pod którą klient nawiązuje kanał
	// WebSocket. Reszta ścieżek trafia do plików klienta.
	SciezkaGniazdaDomyslna = "/ws"
	// pojemnoscKolejkiDomyslna to liczba komunikatów oczekujących na zapis do
	// jednego gniazda. Bufor chroni rdzeń przed zablokowaniem na wolnym
	// urządzeniu; przepełnienie kończy pojedynczą wysyłkę, nie sesję.
	pojemnoscKolejkiDomyslna = 256
	// czasZamknieciaDomyslny ogranicza czas oczekiwania na zamknięcie nasłuchu przy zatrzymaniu tego serwera.
	czasZamknieciaDomyslny = 5 * time.Second
	// KontoDomyslne obowiązuje, dopóki urządzenie nie wskaże konta.
	// Uwierzytelnianie jest jedyną kontrolą dostępu i w fazie budowy nie działa,
	// więc brak konta nie może wstrzymać połączenia.
	KontoDomyslne = "lokalne"
	// ParametrKonta nazywa parametr zapytania i nagłówek, którymi urządzenie
	// wskazuje konto przy nawiązaniu.
	ParametrKonta = "konto"
	// NaglowekKonta jest nagłówkową postacią ParametrKonta — dla klientów, które
	// nie mogą dopisać parametru do adresu.
	NaglowekKonta = "X-Danaco-Konto"
	// ParametrSekretu i NaglowekSekretu niosą sekret nawiązania, którym powłoka
	// przedstawia się przed uaktualnieniem gniazda. Dwie drogi z tego samego
	// powodu co przy koncie: nie każdy klient dopisze nagłówek do adresu gniazda.
	ParametrSekretu = "sekret"
	NaglowekSekretu = "X-Danaco-Sekret"
	// ZmiennaZniesieniaBramki nazywa zmienną środowiska zdejmującą wymóg
	// logowania. Zmienna, nie nastawa: nastawę zmienia komenda z gniazda,
	// a zmienną wskazuje ten, kto uruchamia proces na maszynie.
	ZmiennaZniesieniaBramki = "DANACO_BRAMKA_ZNIESIONA"
	// ZmiennaSekretuNawiazania nazywa zmienną środowiska z sekretem nawiązania;
	// bez niej sprawdzenie sekretu nie obowiązuje, bo nie ma z czym porównywać.
	ZmiennaSekretuNawiazania = "DANACO_SEKRET_NAWIAZANIA"
	// adresDomyslny wiąże nasłuch z pętlą zwrotną, ponieważ domyślna wartość ma być bezpieczna, a szeroka ma być wyborem.
	adresDomyslny = "127.0.0.1"
)

// Ustawienia to komplet nastaw warstwy transportu. Każde pole ma wartość
// domyślną: brak nastawy nigdy nie wstrzymuje startu.
type Ustawienia struct {
	// Adres wskazuje interfejs nasłuchu; pusty bierze adres domyślny, czyli pętlę zwrotną maszyny.
	Adres string
	// WszystkieInterfejsy przywraca dawne zachowanie pustego adresu: nasłuch na wszystkich interfejsach.
	WszystkieInterfejsy bool
	// PochodzeniaDozwolone to wykaz wzorców nagłówka Origin przyjmowanych przy nawiązaniu gniazda.
	PochodzeniaDozwolone []string
	// WymogLogowania jest rozstrzygnięciem Operatora nad dopuszczeniem bramki; wskaźnik niesie trzy stany zamiast dwóch.
	WymogLogowania *bool
	// BramkaZniesiona zdejmuje wymóg logowania i jest jedyną drogą jego zdjęcia.
	// Wartość bierze się ze zmiennej środowiska albo z przełącznika wiersza
	// poleceń — nigdy z nastawy, bo nastawę zmienia komenda z gniazda.
	BramkaZniesiona bool
	// SekretNawiazania jest sekretem, którym powłoka przedstawia się przy
	// uaktualnieniu gniazda. Pusty znaczy brak sprawdzenia: nagłówek Origin nie
	// obowiązuje klienta, który go nie wysyła, więc dopóki powłoka sekretu nie
	// wystawia, dopóty gniazdo stoi otworem dla procesów tej maszyny.
	SekretNawiazania string
	// CertyfikatTLS i KluczTLS wskazują parę plików warstwy TLS; wskazanie obu włącza szyfrowany nasłuch.
	CertyfikatTLS string
	KluczTLS      string
	// Port nasłuchu rdzenia; zero oznacza port domyślny konfiguracji tego serwera transportu.
	Port int
	// PortDowolny każe systemowi wskazać wolny port; rzeczywisty adres odczytuje metoda Adres serwera.
	PortDowolny bool
	// SciezkaGniazda to ścieżka HTTP kanału WebSocket.
	SciezkaGniazda string
	// KatalogKlienta wskazuje pakiet interfejsu; pusty albo zły katalog nie wstrzymuje nasłuchu.
	KatalogKlienta string
	// PojemnoscKolejki to bufor wyjściowy jednego połączenia.
	PojemnoscKolejki int
	// CzasZamkniecia ogranicza łagodne zamknięcie nasłuchu.
	CzasZamkniecia time.Duration
	// Dziennik przyjmuje zapisy diagnostyczne. Nil kieruje je do kosza.
	Dziennik *log.Logger
}

// Funkcja Domyslne zwraca ustawienia obowiązujące bez jawnego wskazania przez operatora tej samej maszyny.
func Domyslne() Ustawienia {
	return Ustawienia{
		Port:             konfiguracja.PortDomyslny,
		SciezkaGniazda:   SciezkaGniazdaDomyslna,
		PojemnoscKolejki: pojemnoscKolejkiDomyslna,
		CzasZamkniecia:   czasZamknieciaDomyslny,
	}
}

// Metoda zNormalizowane uzupełnia wszystkie pola puste tych ustawień wartościami domyślnymi tego transportu.
func (u Ustawienia) zNormalizowane() Ustawienia {
	d := Domyslne()
	if u.Port <= 0 && !u.PortDowolny {
		u.Port = d.Port
	}
	if u.SciezkaGniazda == "" {
		u.SciezkaGniazda = d.SciezkaGniazda
	}
	if u.PojemnoscKolejki <= 0 {
		u.PojemnoscKolejki = d.PojemnoscKolejki
	}
	if u.CzasZamkniecia <= 0 {
		u.CzasZamkniecia = d.CzasZamkniecia
	}
	if u.Dziennik == nil {
		u.Dziennik = log.New(io.Discard, "", 0)
	}
	// Adres pusty przestaje znaczyć wszystkie interfejsy, zaczyna znaczyć brak wskazania, czyli pętlę.
	if strings.TrimSpace(u.Adres) == "" && !u.WszystkieInterfejsy {
		u.Adres = adresDomyslny
	}
	// Zniesienie bramki i sekret nawiązania czyta transport wprost ze środowiska
	// procesu, a nie z nastaw bazy — nastawa przychodzi komendą z gniazda,
	// czyli z tej samej strony, której obie te dźwignie dotyczą.
	if !u.BramkaZniesiona {
		u.BramkaZniesiona = zniesienieZeSrodowiska(os.Getenv(ZmiennaZniesieniaBramki), u.Dziennik)
	}
	if strings.TrimSpace(u.SekretNawiazania) == "" {
		u.SekretNawiazania = strings.TrimSpace(os.Getenv(ZmiennaSekretuNawiazania))
	}
	// Rozpoznanie wystawienia stoi tutaj, bo tędy przechodzi każdy serwer i przechodzi dokładnie raz.
	ostrzezJezeliWystawiony(u)
	return u
}

// zniesienieZeSrodowiska czyta zmienną zdejmującą bramkę. Wartość nieczytelna
// zostaje zniesieniem nieudzielonym: transport nie ma jak przerwać startu z tego
// miejsca, a wybór między „zatrzymaj" a „zdejmij bramkę" rozstrzyga się na
// korzyść wymogu.
func zniesienieZeSrodowiska(tekst string, dziennik *log.Logger) bool {
	tekst = strings.TrimSpace(tekst)
	if tekst == "" {
		return false
	}
	zniesiona, err := strconv.ParseBool(tekst)
	if err != nil {
		if dziennik != nil {
			dziennik.Printf("transport: %s=%q nie jest wartością logiczną (true|false) — bramka zostaje",
				ZmiennaZniesieniaBramki, tekst)
		}
		return false
	}
	return zniesiona
}

// zTLS mówi, czy nasłuch ma iść warstwą szyfrowaną. Para kompletna znaczy tak;
// para pusta — nie. Para niekompletna nie jest tu rozstrzygana, bo to błąd
// wskazania i rozstrzyga go sprawdzTLS przed startem.
func (u Ustawienia) zTLS() bool {
	return strings.TrimSpace(u.CertyfikatTLS) != "" && strings.TrimSpace(u.KluczTLS) != ""
}

// Metoda sprawdzTLS odmawia startu przy wskazaniu połowicznym pary TLS i nazywa brakującą jej połowę pary.
func (u Ustawienia) sprawdzTLS() error {
	certyfikat := strings.TrimSpace(u.CertyfikatTLS)
	klucz := strings.TrimSpace(u.KluczTLS)
	switch {
	case certyfikat != "" && klucz == "":
		return errors.New("transport: wskazano certyfikat TLS bez klucza; " +
			"warstwa TLS wymaga obu plików, a praca otwartym tekstem po wskazaniu certyfikatu byłaby cichym zejściem poniżej wskazania")
	case klucz != "" && certyfikat == "":
		return errors.New("transport: wskazano klucz TLS bez certyfikatu; " +
			"warstwa TLS wymaga obu plików")
	default:
		return nil
	}
}

// opisWarstwy nazywa warstwę nasłuchu w linii dziennika. Czytający dziennik ma
// widzieć, czy połączenie idzie otwartym tekstem, bez wnioskowania z braku.
func (u Ustawienia) opisWarstwy() string {
	if u.zTLS() {
		return "TLS"
	}
	return "otwarty tekst"
}

// Metoda adresNasluchu składa adres nasłuchu przekazywany dalej funkcji sieciowej otwierającej gniazdo.
func (u Ustawienia) adresNasluchu() string {
	port := u.Port
	if u.PortDowolny {
		port = 0
	}
	return net.JoinHostPort(u.Adres, strconv.Itoa(port))
}
