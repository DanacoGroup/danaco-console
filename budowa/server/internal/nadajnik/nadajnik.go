// Pakiet nadajnik podaje serwerowi poczty wiadomości złożone przez pakiet
// mail: potwierdzenie adresu przy rejestracji i drogę odzyskania konta,
// w imieniu platformy, osobnym poświadczeniem od skrzynki Operatora.
package nadajnik

import (
	"crypto/tls"
	"fmt"
	"net"
	netmail "net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/mail"
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
	// SzyfrujStartTLS podnosi połączenie do TLS; ma sens tylko dla przekaźnika
	// lokalnego.
	SzyfrujStartTLS bool
	// WeryfikujTLS wyłącza się wyłącznie wtedy, gdy Operator tak zapisał.
	WeryfikujTLS bool
}

// Gotowe mówi, czy nadajnik ma czym nadać do odbiorcy; rozstrzygnięcie zapada
// tutaj, nie u wołającego kodu.
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

// Nadawca oddaje adres nadawcy dla koperty wiadomości składanej pakietem mail.
// Nazwa nieustawiona zostawia sam adres — konto nadawcze bez nazwy widocznej
// nadal nadaje.
func (n Nastawy) Nadawca() netmail.Address {
	return netmail.Address{
		Name:    strings.TrimSpace(n.NazwaWyswietlana),
		Address: n.Adres,
	}
}

// Wyslij nadaje wiadomość złożoną przez pakiet mail i oddaje chwilę nadania.
//
// Nieudane nadanie NIE jest ciszą: wraca błędem nazywającym, na czym rozmowa
// stanęła — rejestracja bez listu zostawia Operatora przed zamkniętym kontem.
func Wyslij(n Nastawy, wiadomosc *mail.Message) (time.Time, error) {
	if err := n.Brak(); err != nil {
		return time.Time{}, err
	}
	if wiadomosc == nil || len(wiadomosc.Data) == 0 {
		return time.Time{}, fmt.Errorf("list bez treści — nie ma czego nadać")
	}
	if strings.TrimSpace(wiadomosc.To) == "" {
		return time.Time{}, fmt.Errorf("list bez odbiorcy — nie ma dokąd go nadać")
	}

	if err := nadaj(n, wiadomosc.To, wiadomosc.Data); err != nil {
		return time.Time{}, err
	}
	return time.Now(), nil
}

// nadaj prowadzi całą rozmowę SMTP: połączenie, ewentualny STARTTLS,
// uwierzytelnienie warunkowe, kopertę i treść listu.
func nadaj(n Nastawy, odbiorca string, dokument []byte) error {
	// adres IPv6 sam niesie dwukropki, więc łączy się przez JoinHostPort, nie
	// przez sklejenie.
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
