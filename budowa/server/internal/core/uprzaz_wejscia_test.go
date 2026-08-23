package core

import (
	"bufio"
	"context"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/store"
	"danacoconsole/shared"
)

// Uprząż drogi wejścia do aplikacji.
//
// Droga wejścia — rejestracja, potwierdzenie adresu, odzyskanie konta i wykaz
// urządzeń — różni się od reszty rdzenia jedną rzeczą: jej skutek wychodzi poza
// proces. Konto zapisuje się w bazie, a droga potwierdzenia idzie LISTEM i bez
// tego listu Operator nie wejdzie nigdy. Sprawdzian, który mierzy wyłącznie
// kopertę odpowiedzi, przepuściłby oba te braki: `registered: true` wygląda tak
// samo, gdy list poszedł, i gdy przepadł.
//
// Stąd dwie rzeczy, których nie ma uprząż zgodności z kontraktem
// (`uprzaz_test.go`):
//
//   - baza zwracana wołającemu, bo dowodem założenia konta jest wiersz
//     w `konto_wlasciciela`, a nie zdanie w odpowiedzi;
//   - odbiornik SMTP na `127.0.0.1`, bo dowodem wysłania listu jest list.
//
// Odbiornik jest prawdziwym gniazdem na porcie efemerycznym, nie zaślepką
// podstawioną w miejsce `nadajnik.Wyslij`. Zaślepka sprawdzałaby, czy rdzeń woła
// funkcję; gniazdo sprawdza, czy list DOSZEDŁ — i pozwala przeczytać jego treść,
// czyli rozstrzygnąć, że niesie drogę i nie niesie hasła.
//
// Sieci sprawdzian nie dotyka: nasłuch stoi na pętli zwrotnej, a nastawy idą
// z `SzyfrujStartTLS: false`, więc rozmowa nie próbuje ani podnieść TLS, ani
// wyjść poza maszynę.

// Tożsamość konta zakładanego przez sprawdziany. Jedna dla wszystkich, żeby
// niepowodzenie mówiło o zachowaniu rdzenia, a nie o tym, który sprawdzian
// wpisał jaki adres.
const (
	loginSprawdzianu = "wlasciciel"
	adresSprawdzianu = "wlasciciel@danaco.sprawdzian"
	hasloPierwsze    = "haslo-pierwsze-4471"
	hasloDrugie      = "haslo-drugie-9028"
)

// uprzazWejscia trzyma wszystko, czym mierzy się skutek drogi wejścia: rdzeń,
// jego kontekst życia, bazę pod nim i skrzynkę, do której idą listy.
//
// `poczta` jest zerowa w obu trybach bez działającego odbiornika — sprawdzian,
// który po niej sięgnie, ma stanąć od razu, a nie mierzyć pustą skrzynkę.
type uprzazWejscia struct {
	rdzen  *Zmontowany
	zycie  context.Context
	baza   *store.Baza
	poczta *odbiornikSMTP
}

// trybPoczty opisuje stan konta nadawczego platformy. Trzy stany, bo trzy są
// prawdziwe i różnią się skutkiem — a nie dlatego, że sprawdzianom tak wygodnie.
type trybPoczty int

const (
	// pocztaBrak — Operator nie wpisał serwera poczty wychodzącej.
	pocztaBrak trybPoczty = iota
	// pocztaDziala — odbiornik na pętli zwrotnej przyjmuje listy.
	pocztaDziala
	// pocztaNieosiagalna — serwer wskazany, ale nikt na nim nie słucha. Stan
	// codzienny: przekaźnik zatrzymany, zapora, literówka w nazwie hosta.
	// Nastawy wyglądają wtedy poprawnie i brak wychodzi dopiero przy nadawaniu.
	pocztaNieosiagalna
)

// zmontujDrogeWejscia składa rdzeń nad świeżą bazą.
func zmontujDrogeWejscia(t *testing.T, tryb trybPoczty) uprzazWejscia {
	t.Helper()

	katalog := t.TempDir()
	baza, err := store.Otworz(filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy sprawdzianu: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })

	zycie, zakoncz := context.WithCancel(context.Background())
	t.Cleanup(zakoncz)

	ustawienia := konfiguracja.Domyslna()
	ustawienia.KatalogDanych = katalog
	ustawienia.KatalogKlienta = ""
	ustawienia.KatalogProfili = ""
	ustawienia.Port = 0

	// Nastawy nadajnika ustawiane są zawsze — także na pusto. Poleganie na
	// wartości domyślnej wiązałoby wynik sprawdzianu ze zmiennymi środowiska
	// maszyny, na której biegnie.
	var poczta *odbiornikSMTP
	switch tryb {
	case pocztaDziala:
		poczta = podnieOdbiornikSMTP(t)
		ustawienia.NadawcaHost = poczta.Host()
		ustawienia.NadawcaPort = poczta.Port()
		ustawienia.NadawcaAdres = "platforma@danaco.sprawdzian"
		ustawienia.NadawcaNazwa = "Danaco Console"
	case pocztaNieosiagalna:
		ustawienia.NadawcaHost = "127.0.0.1"
		ustawienia.NadawcaPort = portBezNasluchu(t)
		ustawienia.NadawcaAdres = "platforma@danaco.sprawdzian"
		ustawienia.NadawcaNazwa = "Danaco Console"
	default:
		ustawienia.NadawcaHost = ""
		ustawienia.NadawcaPort = 0
		ustawienia.NadawcaAdres = ""
		ustawienia.NadawcaNazwa = ""
	}
	// Rozmowa idzie otwartym tekstem do przekaźnika na tej samej maszynie —
	// dokładnie ten przypadek nastawa dopuszcza. Bez tego klient SMTP próbowałby
	// podnieść połączenie do TLS.
	ustawienia.NadawcaStartTLS = wskaznik(false)

	zmontowany, err := Zmontuj(zycie, Montaz{
		Konfiguracja: ustawienia,
		Baza:         baza,
		Dziennik:     dziennikNiemy(),
	})
	if err != nil {
		t.Fatalf("montaż rdzenia nie powiódł się: %v", err)
	}
	t.Cleanup(zmontowany.Zamknij)

	return uprzazWejscia{rdzen: zmontowany, zycie: zycie, baza: baza, poczta: poczta}
}

// ── odbiornik listów ─────────────────────────────────────────────────────────

// listOdebrany to jeden list, który naprawdę przeszedł przez gniazdo: koperta
// zwrotna, odbiorca i cały dokument z nagłówkami.
type listOdebrany struct {
	Nadawca  string
	Odbiorca string
	Dokument string
}

// odbiornikSMTP jest serwerem poczty na czas sprawdzianu.
//
// Rozmowę prowadzi w zakresie, którego używa nadajnik: powitanie, EHLO, koperta,
// odbiorca, treść, koniec. Ani AUTH, ani STARTTLS nie są ogłaszane — nadajnik
// pomija oba, gdy serwer o nich nie mówi, więc rozmowa nie wymaga poświadczenia
// ani certyfikatu.
type odbiornikSMTP struct {
	nasluch net.Listener
	zamek   sync.Mutex
	listy   []listOdebrany
}

// podnieOdbiornikSMTP otwiera gniazdo na porcie efemerycznym pętli zwrotnej
// i zaczyna przyjmować rozmowy. Gniazdo zamyka się po sprawdzianie samo.
func podnieOdbiornikSMTP(t *testing.T) *odbiornikSMTP {
	t.Helper()

	nasluch, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("nie można podnieść odbiornika listów: %v", err)
	}
	odbiornik := &odbiornikSMTP{nasluch: nasluch}
	t.Cleanup(func() { _ = nasluch.Close() })

	go func() {
		for {
			polaczenie, err := nasluch.Accept()
			if err != nil {
				return
			}
			go odbiornik.rozmawiaj(polaczenie)
		}
	}()
	return odbiornik
}

// portBezNasluchu oddaje numer portu, na którym na pewno nikt nie słucha:
// gniazdo podnosi się i od razu zamyka, więc system zdążył go przydzielić,
// a nasłuchu już nie ma. Wpisanie liczby na sztywno wiązałoby sprawdzian
// z maszyną, na której akurat biegnie.
func portBezNasluchu(t *testing.T) int {
	t.Helper()

	nasluch, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("nie można wybrać portu bez nasłuchu: %v", err)
	}
	port := nasluch.Addr().(*net.TCPAddr).Port
	if err := nasluch.Close(); err != nil {
		t.Fatalf("nie można zwolnić portu bez nasłuchu: %v", err)
	}
	return port
}

// Host oddaje adres, pod którym stoi odbiornik.
func (o *odbiornikSMTP) Host() string {
	return o.nasluch.Addr().(*net.TCPAddr).IP.String()
}

// Port oddaje port nadany przez system.
func (o *odbiornikSMTP) Port() int {
	return o.nasluch.Addr().(*net.TCPAddr).Port
}

// Listy oddaje kopię wykazu listów odebranych do tej chwili.
func (o *odbiornikSMTP) Listy() []listOdebrany {
	o.zamek.Lock()
	defer o.zamek.Unlock()
	kopia := make([]listOdebrany, len(o.listy))
	copy(kopia, o.listy)
	return kopia
}

// Ostatni oddaje list odebrany jako ostatni i przerywa sprawdzian, gdy nie
// przyszedł żaden. Brak listu jest tu niepowodzeniem, nie pustym wynikiem:
// czynność, która obiecała list, a go nie wysłała, zostawia Operatora przed
// kontem bez drogi wejścia.
func (o *odbiornikSMTP) Ostatni(t *testing.T) listOdebrany {
	t.Helper()

	listy := o.Listy()
	if len(listy) == 0 {
		t.Fatal("do skrzynki nie przyszedł żaden list — droga wejścia nie została wysłana")
	}
	return listy[len(listy)-1]
}

// Wyczysc opróżnia skrzynkę, żeby kolejny pomiar liczył listy od zera.
func (o *odbiornikSMTP) Wyczysc() {
	o.zamek.Lock()
	defer o.zamek.Unlock()
	o.listy = nil
}

// rozmawiaj prowadzi jedną rozmowę SMTP i odkłada odebrany list.
func (o *odbiornikSMTP) rozmawiaj(polaczenie net.Conn) {
	defer func() { _ = polaczenie.Close() }()
	_ = polaczenie.SetDeadline(time.Now().Add(30 * time.Second))

	czytnik := bufio.NewReader(polaczenie)
	odpowiedz := func(tresc string) bool {
		_, err := polaczenie.Write([]byte(tresc + "\r\n"))
		return err == nil
	}
	if !odpowiedz("220 odbiornik-sprawdzianu ESMTP") {
		return
	}

	var nadawca, odbiorca string
	for {
		linia, err := czytnik.ReadString('\n')
		if err != nil {
			return
		}
		polecenie := strings.TrimSpace(linia)
		wielkimi := strings.ToUpper(polecenie)
		switch {
		case strings.HasPrefix(wielkimi, "EHLO"):
			// Ani AUTH, ani STARTTLS: rozmowa ma być najprostsza z możliwych.
			if !odpowiedz("250-odbiornik-sprawdzianu\r\n250 8BITMIME") {
				return
			}
		case strings.HasPrefix(wielkimi, "HELO"):
			if !odpowiedz("250 odbiornik-sprawdzianu") {
				return
			}
		case strings.HasPrefix(wielkimi, "MAIL FROM"):
			nadawca = adresZPolecenia(polecenie)
			if !odpowiedz("250 koperta zwrotna przyjęta") {
				return
			}
		case strings.HasPrefix(wielkimi, "RCPT TO"):
			odbiorca = adresZPolecenia(polecenie)
			if !odpowiedz("250 odbiorca przyjęty") {
				return
			}
		case wielkimi == "DATA":
			if !odpowiedz("354 dawaj treść, zakończ kropką") {
				return
			}
			dokument, err := czytajTrescListu(czytnik)
			if err != nil {
				return
			}
			o.zamek.Lock()
			o.listy = append(o.listy, listOdebrany{
				Nadawca: nadawca, Odbiorca: odbiorca, Dokument: dokument,
			})
			o.zamek.Unlock()
			if !odpowiedz("250 list przyjęty") {
				return
			}
		case wielkimi == "QUIT":
			odpowiedz("221 do widzenia")
			return
		case wielkimi == "RSET":
			nadawca, odbiorca = "", ""
			if !odpowiedz("250 wyzerowano") {
				return
			}
		default:
			if !odpowiedz("250 przyjęto") {
				return
			}
		}
	}
}

// czytajTrescListu zbiera dokument aż do samotnej kropki i zdejmuje wypełnienie
// kropką z początku wiersza (`..` → `.`), które wkłada nadawca.
func czytajTrescListu(czytnik *bufio.Reader) (string, error) {
	var dokument strings.Builder
	for {
		linia, err := czytnik.ReadString('\n')
		if err != nil {
			return "", err
		}
		bez := strings.TrimRight(linia, "\r\n")
		if bez == "." {
			return dokument.String(), nil
		}
		if strings.HasPrefix(bez, "..") {
			bez = bez[1:]
		}
		dokument.WriteString(bez)
		dokument.WriteString("\n")
	}
}

// adresZPolecenia wyjmuje adres z `MAIL FROM:<adres>` i `RCPT TO:<adres>`.
func adresZPolecenia(polecenie string) string {
	poczatek := strings.Index(polecenie, "<")
	koniec := strings.Index(polecenie, ">")
	if poczatek < 0 || koniec < poczatek {
		return ""
	}
	return polecenie[poczatek+1 : koniec]
}

// ── odczyt drogi z listu ─────────────────────────────────────────────────────

// naglowekDrogi jest wierszem, po którym w obu listach systemowych stoi droga
// potwierdzenia. Jeden napis dla obu celów, bo oba listy składa jedna funkcja.
const naglowekDrogi = "Droga potwierdzenia:"

// drogaZListu wyjmuje z listu materiał, który Operator ma wpisać w oknie.
//
// Wyjmowanie idzie po treści, a nie po wartości podpatrzonej w bazie — bo
// sprawdzianem jest właśnie to, czy droga DOSZŁA do skrzynki. Skrót z bazy
// przeszedłby także wtedy, gdyby list wyszedł pusty.
func drogaZListu(t *testing.T, list listOdebrany) string {
	t.Helper()

	wiersze := strings.Split(list.Dokument, "\n")
	for numer, wiersz := range wiersze {
		if strings.TrimSpace(wiersz) != naglowekDrogi {
			continue
		}
		for _, dalszy := range wiersze[numer+1:] {
			if droga := strings.TrimSpace(dalszy); droga != "" {
				return droga
			}
		}
	}
	t.Fatalf("list nie niesie drogi potwierdzenia; dokument:\n%s", list.Dokument)
	return ""
}

// ── pomiar stanu bazy ────────────────────────────────────────────────────────

// liczbaWierszy odpowiada na pytanie, na które odpowiedź komendy odpowiedzieć
// nie może: czy w bazie naprawdę coś zostało.
func liczbaWierszy(t *testing.T, u uprzazWejscia, zapytanie string, argumenty ...any) int {
	t.Helper()

	var ile int
	if err := u.baza.DB.QueryRowContext(u.zycie, zapytanie, argumenty...).Scan(&ile); err != nil {
		t.Fatalf("nie można policzyć wierszy (%s): %v", zapytanie, err)
	}
	return ile
}

// ── czynności powtarzane ─────────────────────────────────────────────────────

// zarejestrujWlasciciela wykonuje `auth.register` i zwraca drogę potwierdzenia
// z listu, który po niej przyszedł.
func zarejestrujWlasciciela(t *testing.T, u uprzazWejscia) string {
	t.Helper()

	var odpowiedz shared.AuthRegisterResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRegister, shared.AuthRegisterRequest{
		Login:      loginSprawdzianu,
		Email:      adresSprawdzianu,
		Password:   hasloPierwsze,
		DeviceName: wskaznik("maszyna sprawdzianu"),
	}, &odpowiedz)
	if !odpowiedz.Registered {
		t.Fatal("rejestracja oddała stan udany i registered=false naraz")
	}
	return drogaZListu(t, u.poczta.Ostatni(t))
}
