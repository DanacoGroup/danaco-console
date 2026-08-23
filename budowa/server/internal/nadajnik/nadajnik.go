// Odpowiedzialność pakietu: listy, które pisze SAMA APLIKACJA — potwierdzenie
// adresu przy rejestracji i droga odzyskania konta. Dwa listy, żadnych innych.
//
// ── CZYM TEN NADAJNIK RÓŻNI SIĘ OD MODUŁU POCZTY ────────────────────────────
// Moduł poczty (`internal/poczta`) jest klientem SKRZYNKI OPERATORA: czyta jego
// listy, wysyła w jego imieniu i odkłada kopię w jego folderze wysłanych. Ten
// pakiet nadaje w imieniu PLATFORMY, do Operatora, i kopii nigdzie nie odkłada —
// list systemowy w folderze „wysłane" Operatora byłby śladem czynności, której
// on nie wykonał.
//
// Stąd osobne poświadczenie i osobny host: konto nadawcze platformy nie jest
// skrzynką Operatora i nie wolno ich mieszać. Gdyby aplikacja pisała jego
// kontem, utrata dostępu do skrzynki odcinałaby drogę odzyskania konta —
// czyli dokładnie wtedy, gdy jest potrzebna.
//
// ── ANI JEDNEGO SEKRETU W LIŚCIE ────────────────────────────────────────────
// List niesie DROGĘ potwierdzenia, nie hasło. Hasła platforma nie zna
// w postaci jawnej i nigdy go nie odsyła.
//
// SMTP jedzie biblioteką standardową — bez nowej zależności w go.mod.
package nadajnik

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// limitRozmowy zamyka rozmowę, która utknęła. Bez niego nadanie do serwera,
// który przyjął gniazdo i zamilkł, wisiałoby aż do zamknięcia procesu.
const limitRozmowy = 60 * time.Second

// Nastawy opisują konto nadawcze platformy.
//
// `Sekret` przychodzi z sejfu poświadczeń i pakiet traktuje go jak nieprzezroczysty
// napis: nie zapisuje go, nie loguje i nie umieszcza w treści listu.
type Nastawy struct {
	Host             string
	Port             int
	Uzytkownik       string
	Sekret           string
	Adres            string
	NazwaWyswietlana string
	// SzyfrujStartTLS włącza podniesienie połączenia do TLS. Wyłączenie ma sens
	// wyłącznie dla przekaźnika na tej samej maszynie; w sieci jest błędem.
	SzyfrujStartTLS bool
	// WeryfikujTLS wyłącza się WYŁĄCZNIE wtedy, gdy Operator tak zapisał —
	// nigdy samo, nigdy „bo localhost".
	WeryfikujTLS bool
}

// Gotowe mówi, czy nadajnik ma czym nadać.
//
// Rozstrzygnięcie jest tutaj, a nie u wołającego, bo brak konta nadawczego to
// nie jest usterka do zgłoszenia w połowie rejestracji — to stan, o którym
// warstwa wyżej musi wiedzieć, ZANIM założy konto i obieca list.
func (n Nastawy) Gotowe() bool {
	return strings.TrimSpace(n.Host) != "" && strings.TrimSpace(n.Adres) != ""
}

// Brak nazywa, czego nadajnikowi brakuje. Każdy brak osobnym zdaniem, bo każdy
// Operator naprawia inaczej: jeden wskazaniem serwera, drugi adresem nadawcy.
func (n Nastawy) Brak() error {
	if strings.TrimSpace(n.Host) == "" {
		return fmt.Errorf("konto nadawcze platformy nie ma wskazanego serwera poczty wychodzącej — " +
			"bez niego platforma nie wyśle ani potwierdzenia adresu, ani drogi odzyskania konta")
	}
	if strings.TrimSpace(n.Adres) == "" {
		return fmt.Errorf("konto nadawcze platformy nie ma adresu nadawcy — " +
			"serwer odrzuci list bez koperty zwrotnej")
	}
	return nil
}

// List jest jednym z dwóch listów systemowych: temat i treść tekstowa.
//
// Załączników nie ma i nie będzie: list systemowy niesie jedno zdanie i jedną
// drogę, a załącznik w liście o odzyskaniu konta jest wzorcem, po którym
// rozpoznaje się podszycie.
type List struct {
	Do    string
	Temat string
	Tresc string
}

// Wyslij nadaje list i oddaje chwilę nadania.
//
// Nieudane nadanie NIE jest ciszą: wraca błędem nazywającym, na czym rozmowa
// stanęła. Rejestracja, która obiecała list i go nie wysłała, zostawiłaby
// Operatora przed kontem, do którego nie ma jak wejść.
func Wyslij(n Nastawy, l List) (time.Time, error) {
	if err := n.Brak(); err != nil {
		return time.Time{}, err
	}
	if strings.TrimSpace(l.Do) == "" {
		return time.Time{}, fmt.Errorf("list bez odbiorcy — nie ma dokąd go nadać")
	}

	dokument := zloz(n, l)
	if err := nadaj(n, l.Do, dokument); err != nil {
		return time.Time{}, err
	}
	return time.Now(), nil
}

// zloz składa dokument listu: nagłówki i treść rozdzielone pustą linią.
//
// Kodowanie jest jawne (`UTF-8`), bo temat i treść niosą polskie znaki
// diakrytyczne, a serwer bez deklaracji przyjmie je za bajty ósemkowe
// i Operator zobaczy krzaki zamiast zdania.
func zloz(n Nastawy, l List) []byte {
	nadawca := n.Adres
	if nazwa := strings.TrimSpace(n.NazwaWyswietlana); nazwa != "" {
		nadawca = fmt.Sprintf("%s <%s>", nazwa, n.Adres)
	}
	naglowki := []string{
		"From: " + nadawca,
		"To: " + l.Do,
		"Subject: " + l.Temat,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: 8bit",
		"Date: " + time.Now().Format(time.RFC1123Z),
	}
	return []byte(strings.Join(naglowki, "\r\n") + "\r\n\r\n" + l.Tresc + "\r\n")
}

// nadaj prowadzi rozmowę SMTP: połączenie, STARTTLS, uwierzytelnienie, koperta
// i treść.
//
// Uwierzytelnienie jest warunkowe. Serwer dostawcy zawsze go żąda, ale przekaźnik
// na tej samej maszynie często nie ogłasza AUTH wcale — wpychanie mu wtedy
// poświadczenia kończy się odmową przy komendzie, która bez AUTH przeszłaby.
func nadaj(n Nastawy, odbiorca string, dokument []byte) error {
	// net.JoinHostPort, a nie sklejenie z dwukropkiem: adres IPv6 sam niesie
	// dwukropki i sklejenie dałoby adres nie do rozłożenia.
	adres := net.JoinHostPort(n.Host, strconv.Itoa(port(n)))

	polaczenie, err := net.DialTimeout("tcp", adres, limitRozmowy)
	if err != nil {
		return fmt.Errorf("serwer poczty wychodzącej %s jest nieosiągalny: %w", adres, err)
	}
	_ = polaczenie.SetDeadline(time.Now().Add(limitRozmowy))

	rozmowa, err := smtp.NewClient(polaczenie, n.Host)
	if err != nil {
		_ = polaczenie.Close()
		return fmt.Errorf("serwer %s przyjął połączenie, lecz nie przedstawił się protokołem SMTP: %w", adres, err)
	}
	defer func() { _ = rozmowa.Quit() }()

	if n.SzyfrujStartTLS {
		if czy, _ := rozmowa.Extension("STARTTLS"); czy {
			nastawyTLS := &tls.Config{
				ServerName:         n.Host,
				InsecureSkipVerify: !n.WeryfikujTLS,
			}
			if err := rozmowa.StartTLS(nastawyTLS); err != nil {
				return fmt.Errorf("serwer %s odrzucił podniesienie połączenia do TLS: %w", adres, err)
			}
		}
	}

	if czy, mechanizmy := rozmowa.Extension("AUTH"); czy && strings.TrimSpace(n.Sekret) != "" {
		if sposob := sposobUwierzytelnienia(n, mechanizmy); sposob != nil {
			if err := rozmowa.Auth(sposob); err != nil {
				return fmt.Errorf("serwer %s odrzucił poświadczenie konta nadawczego platformy: %w", adres, err)
			}
		}
	}

	if err := rozmowa.Mail(n.Adres); err != nil {
		return fmt.Errorf("serwer %s odrzucił kopertę zwrotną %s: %w", adres, n.Adres, err)
	}
	if err := rozmowa.Rcpt(odbiorca); err != nil {
		return fmt.Errorf("serwer %s odrzucił odbiorcę %s: %w", adres, odbiorca, err)
	}
	strumien, err := rozmowa.Data()
	if err != nil {
		return fmt.Errorf("serwer %s nie przyjął treści listu: %w", adres, err)
	}
	if _, err := strumien.Write(dokument); err != nil {
		_ = strumien.Close()
		return fmt.Errorf("przerwane przesyłanie treści listu do %s: %w", adres, err)
	}
	if err := strumien.Close(); err != nil {
		return fmt.Errorf("serwer %s odrzucił list przy zamknięciu: %w", adres, err)
	}
	return nil
}

// port oddaje port nadawania: zapisany przy koncie albo domyślny dla trybu.
//
// 587 to port zgłaszania listów ze STARTTLS, 25 to przekaźnik bez szyfrowania.
func port(n Nastawy) int {
	if n.Port > 0 {
		return n.Port
	}
	if n.SzyfrujStartTLS {
		return 587
	}
	return 25
}

// sposobUwierzytelnienia dobiera mechanizm do tego, co serwer ogłosił.
//
// PLAIN idzie wyłącznie po podniesieniu połączenia do TLS — poświadczenie
// nadawcy przesłane otwartym tekstem byłoby oddaniem go każdemu po drodze.
func sposobUwierzytelnienia(n Nastawy, mechanizmy string) smtp.Auth {
	if strings.Contains(mechanizmy, "PLAIN") && n.SzyfrujStartTLS {
		return smtp.PlainAuth("", n.Uzytkownik, n.Sekret, n.Host)
	}
	if strings.Contains(mechanizmy, "LOGIN") && n.SzyfrujStartTLS {
		return smtp.PlainAuth("", n.Uzytkownik, n.Sekret, n.Host)
	}
	if strings.Contains(mechanizmy, "CRAM-MD5") {
		return smtp.CRAMMD5Auth(n.Uzytkownik, n.Sekret)
	}
	return nil
}
