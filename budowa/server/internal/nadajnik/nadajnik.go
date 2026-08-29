// Pakiet nadajnik wysyła listy pisane przez samą aplikację: potwierdzenie
// adresu przy rejestracji i drogę odzyskania konta, w imieniu platformy,
// osobnym poświadczeniem od skrzynki Operatora.
package nadajnik

import (
	"crypto/tls"
	"encoding/base64"
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

// List jest jednym z dwóch listów systemowych: temat i treść tekstowa.
//
// Załączników nie ma i nie będzie: list systemowy niesie jedno zdanie i jedną
// drogę, a załącznik w liście o odzyskaniu konta jest wzorcem, po którym
// rozpoznaje się podszycie.
type List struct {
	Do    string
	Temat string
	Tresc string
	// TrescHtml jest tą samą wiadomością w postaci graficznej — z papeterią,
	// znakiem marki i typografią produktu. Puste znaczy list wyłącznie tekstowy.
	//
	// Obie postacie idą razem, nigdy sama HTML: klient pocztowy, który grafiki
	// nie pokazuje — a takich jest wiele w ustawieniach domyślnych — dostaje
	// wtedy tekst, nie pustą kartkę z kodem, którego nie widać.
	TrescHtml string
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
		"Date: " + time.Now().Format(time.RFC1123Z),
	}
	if strings.TrimSpace(l.TrescHtml) == "" {
		naglowki = append(naglowki,
			"Content-Type: text/plain; charset=UTF-8",
			"Content-Transfer-Encoding: 8bit")
		return []byte(strings.Join(naglowki, "\r\n") + "\r\n\r\n" + l.Tresc + "\r\n")
	}
	return zlozDwiePostacie(naglowki, l)
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

/*
granicaCzesci rozdziela części listu wieloczęściowego. Wartość jest stała
i nie może wystąpić w treści — obie postacie listu systemowego składa rdzeń,
więc żadna z nich nie niesie napisu przypadkowego.
*/
const granicaCzesci = "danaco-console-granica-czesci-listu"

/*
zlozDwiePostacie składa list `multipart/alternative`: najpierw postać tekstowa,
po niej graficzna.

Kolejność jest wiążąca i wynika z RFC 2046: klient pocztowy pokazuje część
OSTATNIĄ, którą umie wyświetlić. Tekst przed HTML-em znaczy więc „pokaż
papeterię, jeżeli umiesz; jeżeli nie — pokaż tekst". Odwrócenie tej kolejności
zostawiłoby z papeterią wyłącznie tych, którzy jej nie potrzebują.

Postać graficzna idzie kodowaniem base64, bo niesie znak marki wpisany w treść
i długie wiersze stylu; ósemkowe kodowanie łamałoby je na siedemdziesiątym
ósmym znaku i rozbijało dokument.
*/
func zlozDwiePostacie(naglowki []string, l List) []byte {
	naglowki = append(naglowki,
		`Content-Type: multipart/alternative; boundary="`+granicaCzesci+`"`)

	var dokument strings.Builder
	dokument.WriteString(strings.Join(naglowki, "\r\n"))
	dokument.WriteString("\r\n\r\n")

	dokument.WriteString("--" + granicaCzesci + "\r\n")
	dokument.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	dokument.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
	dokument.WriteString(l.Tresc)
	dokument.WriteString("\r\n\r\n")

	dokument.WriteString("--" + granicaCzesci + "\r\n")
	dokument.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	dokument.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	dokument.WriteString(lamany(base64.StdEncoding.EncodeToString([]byte(l.TrescHtml))))
	dokument.WriteString("\r\n\r\n")

	dokument.WriteString("--" + granicaCzesci + "--\r\n")
	return []byte(dokument.String())
}

// lamany łamie ciąg base64 na wiersze po 76 znaków, jak żąda RFC 2045.
// Serwer pocztowy odrzuca wiersze dłuższe niż 998 znaków, a jedna postać
// graficzna listu ma ich kilkadziesiąt tysięcy.
func lamany(zakodowany string) string {
	const dlugoscWiersza = 76
	var wynik strings.Builder
	for i := 0; i < len(zakodowany); i += dlugoscWiersza {
		koniec := i + dlugoscWiersza
		if koniec > len(zakodowany) {
			koniec = len(zakodowany)
		}
		if i > 0 {
			wynik.WriteString("\r\n")
		}
		wynik.WriteString(zakodowany[i:koniec])
	}
	return wynik.String()
}
