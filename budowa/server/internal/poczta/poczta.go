// Pakiet poczta jest klientem cudzej skrzynki i niczym więcej. Nie nasłuchuje na
// żadnym porcie, nie zakłada kont, nie kolejkuje poczty i nie pośredniczy przez
// żadną infrastrukturę Danaco: otwiera gniazdo do serwera, który już stoi, mówi
// IMAP-em i SMTP-em tyle, ile trzeba, i się rozłącza.
//
// Wybór bibliotek:
//
// Odbiór — `github.com/emersion/go-imap/v2` (z `go-sasl` jako jej zależnością).
// IMAP-a nie ma w bibliotece standardowej Go. Własny parser protokołu musiałby
// unieść literały (`{123}` z kontynuacją serwera), zbiory sekwencji i UID-ów,
// rekurencyjny BODYSTRUCTURE, odpowiedzi niezamówione przychodzące w środku
// cudzej komendy, kodowanie modified-UTF-7 w nazwach folderów oraz rozbieżności
// IMAP4rev1 wobec rev2. Wybrana biblioteka ma klienta i serwer w jednym drzewie,
// więc jest testowana z obu stron, i mówi obydwoma wydaniami protokołu. Cena
// wyboru: gałąź v2 jest rozwojowa (`v2.0.0-beta.8`). Gałąź v1 jest zamrożona,
// a jej API oddaje wyniki kanałami, co przy wzorcu „jedno wywołanie, jedna
// odpowiedź" wymusza gorszą obsługę błędów. POP3, którym jedzie starsza poczta
// rdzenia (`dane/poczta.go`), nie zna folderów, szkiców ani oznaczeń i nie umie
// APPEND, więc `mail.folder.list`, `mail.draft.save` i `mail.message.flag` nie
// miałyby czym zadziałać.
//
// Rozbiór i składanie MIME — `github.com/emersion/go-message` (z `.../mail`).
// `net/mail` z biblioteki standardowej czyta nagłówki i na tym kończy: nie
// schodzi w zagnieżdżony `multipart/*`, nie odkodowuje
// `Content-Transfer-Encoding`, nie rozpoznaje `Content-Disposition`. Ręczny
// rozbiór przez `mime/multipart` i `mime/quotedprintable` byłby przepisaniem tej
// samej rekursji, a `go-message` i tak wchodzi jako zależność `go-imap/v2`.
//
// Wysyłka — `net/smtp` z biblioteki standardowej, bez nowej zależności. Klient
// SMTP to EHLO, STARTTLS, AUTH oraz MAIL/RCPT/DATA i stdlib robi komplet;
// `github.com/emersion/go-smtp` wnosi ponad to serwer SMTP, którego ten pakiet
// nie buduje.
//
// Pakiet nie zna bazy, sejfu ani kontraktu i nie loguje niczego. Sekret dostaje
// argumentem, trzyma go w polu i oddaje wyłącznie serwerowi, do którego się
// loguje. Przekład na typy kontraktu robi adapter rdzenia.
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

// Protokół skrzynki przychodzącej — wartości takie, jak w kontrakcie
// (`MailProtocol`). Pakiet mówi o nich tekstem, bo kontraktu nie importuje.
//
// Kontrakt zna trzy wartości (`imap`, `jmap`, `pop3`), a ten pakiet mówi jednym
// protokołem — IMAP-em — i rozpoznaje drugi (POP3) wyłącznie po to, żeby
// wyczytać go z nastaw klienta poczty na urządzeniu. Stałej na JMAP tu nie ma,
// bo nie miałaby ani jednego wołacza i wyglądałaby w kodzie jak zdolność wpięta.
// Skrzynkę opisaną JMAP-em `Polacz` odrzuca zdaniem nazywającym brak wprost.
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
	// SzyfrujOdbior i SzyfrujWysylke rozstrzygają, czy gniazdo idzie przez TLS
	// od pierwszego bajtu. Ten pakiet ich NIE WYPROWADZA z portu ani z niczego
	// innego — dostaje je gotowe z wiersza skrzynki i wykonuje. Rozstrzyga je
	// raz, przy podpięciu, adapter rdzenia (`core/adapter_modul_poczta.go`),
	// żeby prawda o tym, jak rdzeń rozmawia z tą skrzynką, leżała w jednym
	// miejscu i przeżywała restart.
	SzyfrujOdbior  bool
	SzyfrujWysylke bool
	// WeryfikujTLS jest jawną decyzją Operatora zapisaną przy skrzynce
	// w kolumnie `tls_weryfikacja`, nie tylną furtką.
	WeryfikujTLS bool
}

// Naglowek to wiadomość BEZ treści — tyle, ile potrzeba, żeby ją rozpoznać
// w wykazie. Ciągnięcie treści całej skrzynki przy każdym pytaniu byłoby
// pobraniem archiwum, więc `Wykaz` zatrzymuje się tutaj.
type Naglowek struct {
	// Identyfikator jest PARĄ „folder:UID" — złożoną w `zlozIdentyfikator`.
	// Kontrakt daje `mail.message.get` i `mail.message.flag` samo `messageId`,
	// bez folderu, a UID w IMAP-ie jest unikalny wyłącznie w obrębie folderu:
	// bez folderu w identyfikatorze nie dałoby się odnaleźć listu, którego
	// wykaz przed chwilą oddał (i szukanie po wszystkich folderach byłoby
	// zgadywaniem, który z kilku listów o tym samym UID jest ten właściwy).
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
	// NazwyZalacznikow, nie same załączniki: wykaz mówi, CO jest w liście,
	// a bajty ciągnie dopiero `Pobierz`.
	NazwyZalacznikow []string
}

// Zalacznik niesie BAJTY — magazyn rdzenia bierze je i odkłada pod sumą sha256.
type Zalacznik struct {
	Nazwa     string
	TypTresci string
	Bajty     []byte
}

// List to nagłówek wraz z treścią i załącznikami — wynik `Pobierz`.
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
// inaczej (`INBOX.Drafts`, `[Gmail]/Wersje robocze`), więc `odnajdzFolder`
// pyta serwer o folder ze znacznikiem specjalnym i schodzi na te nazwy dopiero
// wtedy, gdy serwer znacznika nie oddał.
const (
	FolderOdebranych = "INBOX"
	FolderSzkicow    = "Drafts"
	FolderWyslanych  = "Sent"
)

// Klient jest OTWARTYM połączeniem do skrzynki Operatora. Nie jest pulą i nie
// jest długowieczny: rdzeń otwiera go na czas jednej komendy i zamyka. Skrzynka
// należy do Operatora, więc trzymanie w niej stałego uchwytu (i zajmowanie
// limitu połączeń jego dostawcy) byłoby zabraniem sobie czegoś, czego nam nie
// dano.
type Klient struct {
	imap    *imapclient.Client
	nastawy Nastawy
}

// Polacz otwiera połączenie do skrzynki i loguje się poświadczeniem z nastaw.
//
// Każda odmowa nazywa swój brak osobnym zdaniem. Brak hosta, brak poświadczenia,
// protokół nieobsługiwany, serwer nieosiągalny i odrzucone poświadczenie to pięć
// różnych braków; jedno wspólne „nie udało się połączyć" zostawiłoby Operatora
// bez wskazówki, co ma naprawić.
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

	// net.JoinHostPort — powód ten sam, co przy wysyłce (`smtp.go`): adres IPv6
	// niesie własne dwukropki.
	adres := net.JoinHostPort(n.HostOdbioru, strconv.Itoa(port(n.PortOdbioru, n.SzyfrujOdbior, 993, 143)))
	opcje := &imapclient.Options{TLSConfig: &tls.Config{
		ServerName: n.HostOdbioru,
		// Weryfikacja łańcucha wyłącza się WYŁĄCZNIE wtedy, gdy Operator tak
		// zapisał przy skrzynce — nigdy sama, nigdy „bo localhost".
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
// do TLS-a, gdy serwer to potrafi. Port 143 ze STARTTLS jest u dostawców regułą,
// a `DialInsecure` puściłoby hasło Operatora łączem jawnym także tam, gdzie
// serwer oferuje szyfrowanie.
//
// Biblioteka nie wystawia wariantu „podnieś, jeśli ogłoszone" — ma
// `DialStartTLS` (podnosi zawsze) albo `DialInsecure` (nie podnosi nigdy).
// Kolejność jest więc od mocniejszej drogi do słabszej, ze zejściem wyłącznie
// przy odmowie protokołu i nigdy przy wadzie certyfikatu: odpowiedź `BAD`/`NO`
// na komendę STARTTLS znaczy, że serwer jej nie zna, i wtedy gniazdo jawne jest
// jedyną drogą. Zerwany uścisk TLS albo niezaufany łańcuch znaczą co innego —
// szyfrowanie jest, tylko mu nie ufamy — a ciche zejście na jawne obniżyłoby
// ochronę dokładnie w chwili, w której pojawił się powód do czujności.
func polaczJawnymGniazdem(adres string, opcje *imapclient.Options) (*imapclient.Client, error) {
	polaczenie, err := imapclient.DialStartTLS(adres, opcje)
	if err == nil {
		return polaczenie, nil
	}
	var odmowaSerwera *imap.Error
	if !errors.As(err, &odmowaSerwera) {
		return nil, err
	}
	// Serwer nie zna STARTTLS — gniazdo jawne jest tu jedyną drogą. Tak stoi
	// skrzynka na tej samej maszynie i przekaźnik w sieci Operatora.
	return imapclient.DialInsecure(adres, opcje)
}

// Zamknij kończy rozmowę z serwerem. Nieudane LOGOUT nie jest błędem komendy —
// praca została wykonana przed nim, a gniazdo i tak zamykamy.
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
