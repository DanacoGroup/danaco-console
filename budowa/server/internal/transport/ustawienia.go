package transport

import (
	"errors"
	"io"
	"log"
	"net"
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
	// czasZamknieciaDomyslny ogranicza oczekiwanie na zamknięcie nasłuchu.
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
	// adresDomyslny wiąże nasłuch z pętlą zwrotną. Pusty adres znaczy dla
	// net.Listen wszystkie interfejsy, więc domyślną wartością musi być pętla
	// zwrotna: domyślna ma być bezpieczna, a szeroka ma być wyborem. Wyjście poza
	// pętlę zwrotną jest osiągalne jednym polem (WszystkieInterfejsy albo Adres
	// wprost) i nadal ostrzega. Stała jest wewnętrzna: wołający wskazują adres,
	// nie sięgają po domyślny.
	adresDomyslny = "127.0.0.1"
)

// Ustawienia to komplet nastaw warstwy transportu. Każde pole ma wartość
// domyślną: brak nastawy nigdy nie wstrzymuje startu.
type Ustawienia struct {
	// Adres wskazuje interfejs nasłuchu. Pusty bierze adres domyślny, czyli pętlę
	// zwrotną; wystawienie szersze wskazuje się wprost albo polem
	// WszystkieInterfejsy. Rdzeń zmienia umiejscowienie w czasie, więc
	// droga do wystawienia zostaje otwarta — zmienia się wyłącznie to, co dzieje
	// się bez wskazania.
	Adres string
	// WszystkieInterfejsy przywraca dawne zachowanie pustego adresu: nasłuch na
	// wszystkich interfejsach maszyny. Osobne pole, a nie pusty napis, bo
	// „nie wskazałem" i „chcę wszędzie" to dwa różne zdania i mają wyglądać
	// różnie w miejscu wywołania.
	WszystkieInterfejsy bool
	// PochodzeniaDozwolone to wykaz wzorców nagłówka Origin przyjmowanych przy
	// nawiązaniu gniazda. Pusty wykaz bierze pochodzenia własne (patrz
	// pochodzeniaWlasne w nawiazanie.go).
	PochodzeniaDozwolone []string
	// WymogLogowania jest dźwignią Operatora nad strażą bramki (`bramka.go`).
	// Trzy stany, nie dwa, i dlatego wskaźnik:
	//   - nil  — Operator nie wskazał nic, rozstrzyga adres nasłuchu (pętla
	//            zwrotna: bez wymogu; szerzej: z wymogiem);
	//   - true — wymóg obowiązuje także na pętli zwrotnej;
	//   - false — wymóg zniesiony także przy nasłuchu szerszym; wolno, ale
	//            dziennik mówi o tym wprost, bo wtedy maszyny Operatora
	//            (tor zdalny) stoją otworem.
	// Wartość logiczna zamiast wskaźnika kasowałaby różnicę między „nie
	// wskazałem" a „wskazałem: nie" — a to jest tu cała różnica.
	WymogLogowania *bool
	// CertyfikatTLS i KluczTLS wskazują parę plików warstwy TLS. Wskazanie obu
	// przełącza nasłuch na https/wss; brak obu zostawia otwarty tekst i — poza
	// pętlą zwrotną — ostrzeżenie w dzienniku. Wskazanie jednego z dwóch jest
	// błędem konfiguracji i zatrzymuje start, bo cicha praca otwartym tekstem
	// przy wskazanym certyfikacie byłaby najgorszym z możliwych wyników.
	CertyfikatTLS string
	KluczTLS      string
	// Port nasłuchu rdzenia. Zero oznacza port domyślny konfiguracji; wartość
	// ujemna nie występuje, bo konfiguracja sprawdza zakres.
	Port int
	// PortDowolny każe systemowi wskazać wolny port (nasłuch na porcie 0).
	// Rzeczywisty adres odczytuje się metodą Adres serwera.
	PortDowolny bool
	// SciezkaGniazda to ścieżka HTTP kanału WebSocket.
	SciezkaGniazda string
	// KatalogKlienta wskazuje pakiet interfejsu (client/dist). Pusty albo
	// nieistniejący katalog nie wstrzymuje nasłuchu — gniazdo działa bez
	// plików statycznych.
	KatalogKlienta string
	// PojemnoscKolejki to bufor wyjściowy jednego połączenia.
	PojemnoscKolejki int
	// CzasZamkniecia ogranicza łagodne zamknięcie nasłuchu.
	CzasZamkniecia time.Duration
	// Dziennik przyjmuje zapisy diagnostyczne. Nil kieruje je do kosza.
	Dziennik *log.Logger
}

// Domyslne zwraca ustawienia obowiązujące bez wskazania Operatora.
func Domyslne() Ustawienia {
	return Ustawienia{
		Port:             konfiguracja.PortDomyslny,
		SciezkaGniazda:   SciezkaGniazdaDomyslna,
		PojemnoscKolejki: pojemnoscKolejkiDomyslna,
		CzasZamkniecia:   czasZamknieciaDomyslny,
	}
}

// zNormalizowane uzupełnia pola puste wartościami domyślnymi.
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
	// Adres pusty przestaje znaczyć „wszystkie interfejsy" i zaczyna znaczyć
	// „nie wskazano" — a wtedy obowiązuje pętla zwrotna. Kto chce szerzej, mówi
	// to wprost jednym z dwóch pól.
	if strings.TrimSpace(u.Adres) == "" && !u.WszystkieInterfejsy {
		u.Adres = adresDomyslny
	}
	// Rozpoznanie wystawienia stoi tutaj, bo tędy przechodzi każdy serwer i
	// przechodzi dokładnie raz: zNormalizowane woła wyłącznie Nowy, a Serwera
	// nie da się zbudować inaczej (pola nieeksportowane). Wpięcie przy samym
	// net.Listen byłoby bliżej faktu, ale adresNasluchu wołane jest kilka razy
	// i ostrzeżenie by się dublowało. Ostrzeżenie idzie po ustaleniu dziennika,
	// żeby brak dziennika kierował je do kosza, a nie gubił wywołania.
	ostrzezJezeliWystawiony(u)
	return u
}

// zTLS mówi, czy nasłuch ma iść warstwą szyfrowaną. Para kompletna znaczy tak;
// para pusta — nie. Para niekompletna nie jest tu rozstrzygana, bo to błąd
// wskazania i rozstrzyga go sprawdzTLS przed startem.
func (u Ustawienia) zTLS() bool {
	return strings.TrimSpace(u.CertyfikatTLS) != "" && strings.TrimSpace(u.KluczTLS) != ""
}

// sprawdzTLS odmawia startu przy wskazaniu połowicznym i nazywa brakującą połowę.
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

// adresNasluchu składa adres przekazywany funkcji net.Listen.
func (u Ustawienia) adresNasluchu() string {
	port := u.Port
	if u.PortDowolny {
		port = 0
	}
	return net.JoinHostPort(u.Adres, strconv.Itoa(port))
}
