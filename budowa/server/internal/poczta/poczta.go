// Pakiet poczta jest klientem cudzej skrzynki i niczym więcej: otwiera
// gniazdo do serwera, który już stoi, mówi IMAP-em i SMTP-em tyle, ile
// trzeba, i się rozłącza.
package poczta

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

// Protokół skrzynki przychodzącej — wartości jak w kontrakcie. Pakiet mówi
// jednym protokołem, IMAP-em, rozpoznając drugi (POP3) wyłącznie po to, żeby
// wyczytać go z nastaw klienta poczty na urządzeniu.
const (
	ProtokolImap = "imap"
	ProtokolPop3 = "pop3"
)

// Nastawy opisują JEDNĄ skrzynkę Operatora — tyle, ile trzeba, żeby się do niej
// zalogować. Sekret jest tu polem przelotnym: przychodzi z sejfu rdzenia tuż
// przed połączeniem i ginie razem z tą strukturą.
type Nastawy struct {
	Adres            string
	NazwaWyswietlana string
	Uzytkownik       string
	Sekret           string
	Protokol         string
	HostOdbioru      string
	PortOdbioru      int
	HostWysylki      string
	PortWysylki      int
	// SzyfrujOdbior i SzyfrujWysylke: ten pakiet dostaje je gotowe.
	SzyfrujOdbior  bool
	SzyfrujWysylke bool
	// WeryfikujTLS jest jawną decyzją Operatora zapisaną przy skrzynce.
	WeryfikujTLS bool
}

// Naglowek to wiadomość BEZ treści — tyle, ile potrzeba, żeby ją rozpoznać
// w wykazie. Ciągnięcie treści całej skrzynki przy każdym pytaniu byłoby
// pobraniem archiwum, więc `Wykaz` zatrzymuje się tutaj.
type Naglowek struct {
	// Identyfikator jest parą folder:UID, bo UID jest unikalny w folderze.
	Identyfikator  string
	Folder         string
	Od             string
	Do             []string
	Kopia          []string
	Temat          string
	Chwila         time.Time
	Zapowiedz      string
	Nieprzeczytana bool
	Watek          string
	// NazwyZalacznikow, nie same załączniki: bajty ciągnie dopiero `Pobierz`.
	NazwyZalacznikow []string
}

// Zalacznik niesie bajty, które magazyn rdzenia bierze i odkłada pod sumą
// kontrolną sha256, wraz z nazwą i typem MIME.
type Zalacznik struct {
	Nazwa     string
	TypTresci string
	Bajty     []byte
}

// List to nagłówek wraz z treścią i załącznikami — pełny wynik `Pobierz`,
// w odróżnieniu od samego nagłówka zwracanego przy wykazie.
type List struct {
	Naglowek
	Tresc      string
	Zalaczniki []Zalacznik
}

// Zawezenie odwzorowuje pola `mail.message.list`. Wszystkie są opcjonalne:
// zawężenie puste znaczy „cały folder", a nie „nic".
type Zawezenie struct {
	Folder              string
	Fraza               string
	Nadawca             string
	Od                  time.Time
	TylkoNieprzeczytane bool
	Granica             int
}

// Wychodzacy to list DO NADANIA — wspólny kształt dla szkicu i wysyłki, bo
// jedno i drugie składa dokładnie ten sam dokument MIME. Różni je
// wyłącznie to, dokąd ten dokument trafia: do folderu szkiców albo do SMTP.
type Wychodzacy struct {
	Od               string
	NazwaWyswietlana string
	Do               []string
	Kopia            []string
	Temat            string
	Tresc            string
	WOdpowiedziNa    string
	Zalaczniki       []Zalacznik
}

// FolderSzkicow i FolderWyslanych to nazwy domyślne. Serwer bywa nazywa je
// inaczej, więc `odnajdzFolder` pyta o folder ze znacznikiem specjalnym
// najpierw, schodząc na te nazwy dopiero gdy znacznika nie oddał.
const (
	FolderOdebranych = "INBOX"
	FolderSzkicow    = "Drafts"
	FolderWyslanych  = "Sent"
)

// Klient jest otwartym połączeniem do skrzynki Operatora. Nie jest pulą i nie
// jest długowieczny: rdzeń otwiera go na czas jednej komendy i zamyka.
type Klient struct {
	imap    *imapclient.Client
	nastawy Nastawy
}

// Polacz otwiera połączenie do skrzynki i loguje się poświadczeniem z nastaw.
// Każda odmowa nazywa swój brak osobnym zdaniem — pięć różnych braków, nie
// jedno wspólne "nie udało się połączyć".
func Polacz(n Nastawy) (*Klient, error) {
	if strings.TrimSpace(n.HostOdbioru) == "" {
		return nil, fmt.Errorf("skrzynka %s nie ma wskazanego serwera poczty przychodzącej — "+
			"rdzeń nie zgaduje hosta dostawcy", n.Adres)
	}
	if strings.TrimSpace(n.Sekret) == "" {
		return nil, fmt.Errorf("skrzynka %s nie ma poświadczenia w sejfie rdzenia — "+
			"podepnij ją ponownie komendą mail.account.add, podając hasło albo token", n.Adres)
	}
	if protokol(n) != ProtokolImap {
		return nil, fmt.Errorf("skrzynka %s jest opisana protokołem %q, a rdzeń mówi dziś wyłącznie IMAP-em — "+
			"JMAP i POP3 są w kontrakcie, lecz w rdzeniu ich nie ma", n.Adres, protokol(n))
	}

	// net.JoinHostPort, bo adres IPv6 niesie własne dwukropki.
	adres := net.JoinHostPort(n.HostOdbioru, strconv.Itoa(port(n.PortOdbioru, n.SzyfrujOdbior, 993, 143)))
	opcje := &imapclient.Options{TLSConfig: &tls.Config{
		ServerName: n.HostOdbioru,
		// Weryfikacja łańcucha wyłącza się wyłącznie decyzją Operatora.
		InsecureSkipVerify: !n.WeryfikujTLS,
	}}

	var (
		polaczenie *imapclient.Client
		err        error
	)
	if n.SzyfrujOdbior {
		polaczenie, err = imapclient.DialTLS(adres, opcje)
	} else {
		polaczenie, err = polaczJawnymGniazdem(adres, opcje)
	}
	if err != nil {
		return nil, fmt.Errorf("serwer poczty %s nie odpowiedział: %w", adres, err)
	}

	uzytkownik := n.Uzytkownik
	if strings.TrimSpace(uzytkownik) == "" {
		uzytkownik = n.Adres
	}
	if err := polaczenie.Login(uzytkownik, n.Sekret).Wait(); err != nil {
		_ = polaczenie.Close()
		return nil, fmt.Errorf("serwer poczty %s odrzucił poświadczenie użytkownika %s: %w",
			adres, uzytkownik, err)
	}
	return &Klient{imap: polaczenie, nastawy: n}, nil
}

// polaczJawnymGniazdem otwiera połączenie na porcie nieszyfrowanym i podnosi je
// do TLS-a, gdy serwer to potrafi. Kolejność idzie od mocniejszej drogi do
// słabszej, ze zejściem wyłącznie przy odmowie protokołu, nigdy przy wadzie
// certyfikatu.
func polaczJawnymGniazdem(adres string, opcje *imapclient.Options) (*imapclient.Client, error) {
	polaczenie, err := imapclient.DialStartTLS(adres, opcje)
	if err == nil {
		return polaczenie, nil
	}
	var odmowaSerwera *imap.Error
	if !errors.As(err, &odmowaSerwera) {
		return nil, err
	}
	// Serwer nie zna STARTTLS — gniazdo jawne jest tu jedyną drogą.
	return imapclient.DialInsecure(adres, opcje)
}

// Zamknij kończy rozmowę z serwerem. Nieudane LOGOUT nie jest błędem komendy —
// praca została wykonana przed nim, a gniazdo i tak się zamyka.
func (k *Klient) Zamknij() {
	if k == nil || k.imap == nil {
		return
	}
	_ = k.imap.Logout().Wait()
	_ = k.imap.Close()
}

// protokol oddaje protokół nastaw, biorąc IMAP jako domyślny — kontrakt mówi
// wprost „brak bierze imap".
func protokol(n Nastawy) string {
	p := strings.ToLower(strings.TrimSpace(n.Protokol))
	if p == "" {
		return ProtokolImap
	}
	return p
}

// port oddaje port wskazany, a przy jego braku — domyślny dla trybu szyfrowania.
// Zero znaczy „nie wskazano", nie „port zerowy".
func port(wskazany int, szyfrowany bool, zTLS, bezTLS int) int {
	if wskazany > 0 {
		return wskazany
	}
	if szyfrowany {
		return zTLS
	}
	return bezTLS
}
