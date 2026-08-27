// Odpowiedzialność pliku: wysyłka — jedyna czynność rdzenia, której skutek
// wychodzi poza maszynę Operatora i której nie da się cofnąć. Kopia w `Sent`
// idzie po udanym nadaniu, nigdy przed.
package poczta

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
)

// limitRozmowySMTP zamyka rozmowę, która utknęła. Bez niego wysyłka do serwera,
// który przyjął gniazdo i zamilkł, wisiałaby aż do zamknięcia procesu.
const limitRozmowySMTP = 60 * time.Second

// Wyslij nadaje list i oddaje jego identyfikator wraz z chwilą nadania.
// Identyfikator pochodzi z nagłówka `Message-ID` złożonego dokumentu, nie
// jest wymyślony po fakcie.
func (k *Klient) Wyslij(w Wychodzacy) (string, time.Time, error) {
	adresaci := append(append([]string{}, oczyszczone(w.Do)...), oczyszczone(w.Kopia)...)
	if len(adresaci) == 0 {
		return "", time.Time{}, fmt.Errorf("list bez ani jednego odbiorcy — nie ma dokąd go nadać")
	}
	if strings.TrimSpace(k.nastawy.HostWysylki) == "" {
		return "", time.Time{}, fmt.Errorf("skrzynka %s nie ma wskazanego serwera poczty wychodzącej — "+
			"rdzeń nie zgaduje, przez kogo nadać list", k.nastawy.Adres)
	}

	dokument, err := zlozWychodzacy(w)
	if err != nil {
		return "", time.Time{}, err
	}
	identyfikator := identyfikatorDokumentu(dokument)

	if err := k.nadaj(adresaci, dokument); err != nil {
		return "", time.Time{}, err
	}
	nadano := time.Now()

	// Kopia do wysłanych idzie po nadaniu, ze znacznikiem Seen.
	if folder, err := k.odnajdzFolder(imap.MailboxAttrSent, FolderWyslanych); err == nil {
		_, _ = k.dolozDoFolderu(folder, dokument, []imap.Flag{imap.FlagSeen})
	}
	return identyfikator, nadano, nil
}

// nadaj prowadzi całą rozmowę SMTP: połączenie, STARTTLS, uwierzytelnienie,
// koperta i treść. Uwierzytelnienie jest warunkowe, bo nie każdy serwer
// ogłasza AUTH.
func (k *Klient) nadaj(adresaci []string, dokument []byte) error {
	// net.JoinHostPort, a nie sklejenie z dwukropkiem, bo adres IPv6 sam
	// niesie dwukropki.
	adres := net.JoinHostPort(k.nastawy.HostWysylki,
		strconv.Itoa(port(k.nastawy.PortWysylki, k.nastawy.SzyfrujWysylke, 465, 587)))
	ustawienia := &tls.Config{
		ServerName:         k.nastawy.HostWysylki,
		InsecureSkipVerify: !k.nastawy.WeryfikujTLS,
	}

	var (
		gniazdo net.Conn
		err     error
	)
	if k.nastawy.SzyfrujWysylke {
		gniazdo, err = tls.DialWithDialer(&net.Dialer{Timeout: limitRozmowySMTP}, "tcp", adres, ustawienia)
	} else {
		gniazdo, err = net.DialTimeout("tcp", adres, limitRozmowySMTP)
	}
	if err != nil {
		return fmt.Errorf("serwer poczty wychodzącej %s nie odpowiedział: %w", adres, err)
	}
	_ = gniazdo.SetDeadline(time.Now().Add(limitRozmowySMTP))

	rozmowa, err := smtp.NewClient(gniazdo, k.nastawy.HostWysylki)
	if err != nil {
		_ = gniazdo.Close()
		return fmt.Errorf("serwer poczty wychodzącej %s nie przywitał się poprawnie: %w", adres, err)
	}
	defer func() { _ = rozmowa.Close() }()

	// STARTTLS podnosi połączenie do szyfrowanego tylko, gdy serwer je ogłasza.
	if !k.nastawy.SzyfrujWysylke {
		if jest, _ := rozmowa.Extension("STARTTLS"); jest {
			if err := rozmowa.StartTLS(ustawienia); err != nil {
				return fmt.Errorf("serwer poczty wychodzącej %s odrzucił szyfrowanie STARTTLS: %w", adres, err)
			}
		}
	}
	if jest, mechanizmy := rozmowa.Extension("AUTH"); jest && strings.TrimSpace(k.nastawy.Sekret) != "" {
		if err := rozmowa.Auth(sposobUwierzytelnienia(k, mechanizmy)); err != nil {
			return fmt.Errorf("serwer poczty wychodzącej %s odrzucił poświadczenie skrzynki %s: %w",
				adres, k.nastawy.Adres, err)
		}
	}

	if err := rozmowa.Mail(k.nastawy.Adres); err != nil {
		return fmt.Errorf("serwer poczty wychodzącej %s nie przyjął nadawcy %s: %w", adres, k.nastawy.Adres, err)
	}
	for _, odbiorca := range adresaci {
		if err := rozmowa.Rcpt(odbiorca); err != nil {
			return fmt.Errorf("serwer poczty wychodzącej %s nie przyjął odbiorcy %s: %w", adres, odbiorca, err)
		}
	}
	strumien, err := rozmowa.Data()
	if err != nil {
		return fmt.Errorf("serwer poczty wychodzącej %s nie przyjął treści listu: %w", adres, err)
	}
	if _, err := strumien.Write(dokument); err != nil {
		return fmt.Errorf("serwer poczty wychodzącej %s przerwał przyjmowanie treści: %w", adres, err)
	}
	// Tu list wychodzi w świat: domknięcie strumienia jest kropką po `DATA`.
	if err := strumien.Close(); err != nil {
		return fmt.Errorf("serwer poczty wychodzącej %s nie potwierdził przyjęcia listu: %w", adres, err)
	}
	return rozmowa.Quit()
}

// sposobUwierzytelnienia wybiera mechanizm SASL zgodny z tym, co ogłosił
// serwer. CRAM-MD5 pierwszy, PLAIN drugi, bo CRAM nie posyła hasła wcale.
func sposobUwierzytelnienia(k *Klient, mechanizmy string) smtp.Auth {
	uzytkownik := k.nastawy.Uzytkownik
	if strings.TrimSpace(uzytkownik) == "" {
		uzytkownik = k.nastawy.Adres
	}
	if strings.Contains(strings.ToUpper(mechanizmy), "CRAM-MD5") {
		return smtp.CRAMMD5Auth(uzytkownik, k.nastawy.Sekret)
	}
	return smtp.PlainAuth("", uzytkownik, k.nastawy.Sekret, k.nastawy.HostWysylki)
}

// identyfikatorDokumentu wyciąga `Message-ID` ze złożonego listu. Nagłówek
// zawsze tam jest — `zlozWychodzacy` nadaje go przed złożeniem — więc pustka
// tutaj znaczyłaby usterkę składania, nie brak w liście.
func identyfikatorDokumentu(dokument []byte) string {
	for _, linia := range strings.Split(string(dokument), "\r\n") {
		if linia == "" {
			break // koniec nagłówków
		}
		if strings.HasPrefix(strings.ToLower(linia), "message-id:") {
			return strings.TrimSpace(linia[len("message-id:"):])
		}
	}
	return ""
}

// oczyszczone odsiewa adresy puste z wykazu, tak dla „do", jak dla „kopia" —
// pusty odbiorca zerwałby komendę RCPT serwera poczty wychodzącej w rozmowie.
func oczyszczone(lista []string) []string {
	wynik := make([]string, 0, len(lista))
	for _, a := range lista {
		if a = strings.TrimSpace(a); a != "" {
			wynik = append(wynik, a)
		}
	}
	return wynik
}
