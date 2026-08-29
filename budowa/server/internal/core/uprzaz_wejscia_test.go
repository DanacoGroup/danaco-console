// Plik mierzy uprząż drogi wejścia: rejestrację, potwierdzenie adresu listem,
// odzyskanie konta i wykaz urządzeń — jedyną drogę rdzenia, której skutek
// wychodzi poza proces, do bazy i do prawdziwej skrzynki SMTP.
package core

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"mime/quotedprintable"
	"net"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/store"
	"danacoconsole/shared"
)

// Tożsamość konta zakładanego przez sprawdziany, jedna dla wszystkich, żeby
// niepowodzenie mówiło o zachowaniu rdzenia, a nie o tym, który sprawdzian
// wpisał jaki adres.
const (
	loginSprawdzianu = "wlasciciel"
	adresSprawdzianu = "wlasciciel@danaco.sprawdzian"
	hasloPierwsze    = "haslo-pierwsze-4471"
	hasloDrugie      = "haslo-drugie-9028"
)

// uprzazWejscia trzyma wszystko, czym mierzy się skutek drogi wejścia: rdzeń,
// jego kontekst życia, bazę pod nim i skrzynkę, do której idą listy. Poczta
// jest zerowa bez działającego odbiornika, sprawdzian ma wtedy stanąć od razu.
type uprzazWejscia struct {
	rdzen  *Zmontowany
	zycie  context.Context
	baza   *store.Baza
	poczta *odbiornikSMTP
	// katalog jest katalogiem danych rdzenia, gdzie leży sejf ze znacznikiem bramki bez poczty.
	katalog string
}

// trybPoczty opisuje stan konta nadawczego platformy: trzy stany, bo trzy są
// prawdziwe i różnią się skutkiem, a nie dlatego, że sprawdzianom tak wygodnie.
type trybPoczty int

const (
	// pocztaBrak — nie wpisano serwera poczty wychodzącej, konto nadawcze
	// platformy pozostaje puste i droga potwierdzenia nie ma jak wyjść z rdzenia.
	pocztaBrak trybPoczty = iota
	// pocztaDziala — odbiornik na pętli zwrotnej przyjmuje listy i droga
	// potwierdzenia dociera do skrzynki tak, jak dotarłaby do adresata.
	pocztaDziala
	// pocztaNieosiagalna — serwer wskazany, ale nikt na nim nie słucha. Stan
	// codzienny: przekaźnik zatrzymany, zapora, literówka w nazwie hosta.
	pocztaNieosiagalna
)

// zmontujDrogeWejscia składa rdzeń nad świeżą bazą, z nastawami nadajnika
// ustawianymi zawsze, także na pusto, żeby wynik nie zależał od środowiska.
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

	// Nastawy nadajnika ustawiane są zawsze, także na pusto, by wynik nie
	// zależał od środowiska maszyny.
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
	// Rozmowa idzie otwartym tekstem do przekaźnika na tej samej maszynie,
	// bez próby podniesienia TLS.
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

	return uprzazWejscia{rdzen: zmontowany, zycie: zycie, baza: baza, poczta: poczta,
		katalog: katalog}
}

// ── odbiornik listów ─────────────────────────────────────────────────────────

// listOdebrany to jeden list, który naprawdę przeszedł przez gniazdo: koperta
// zwrotna, odbiorca i cały dokument z nagłówkami odebrany przez testowy serwer.
type listOdebrany struct {
	Nadawca  string
	Odbiorca string
	Dokument string
}

// odbiornikSMTP jest serwerem poczty na czas sprawdzianu, prowadzącym rozmowę
// w zakresie, którego używa nadajnik: powitanie, EHLO, koperta, odbiorca,
// treść, koniec, bez AUTH i STARTTLS.
type odbiornikSMTP struct {
	nasluch net.Listener
	zamek   sync.Mutex
	listy   []listOdebrany
}

// podnieOdbiornikSMTP otwiera gniazdo na porcie efemerycznym pętli zwrotnej
// i zaczyna przyjmować rozmowy SMTP. Gniazdo zamyka się po sprawdzianie samo.
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
// a nasłuchu już nie ma.
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

// Host oddaje adres pętli zwrotnej, pod którym stoi odbiornik testowego
// serwera SMTP podniesiony na czas sprawdzianu.
func (o *odbiornikSMTP) Host() string {
	return o.nasluch.Addr().(*net.TCPAddr).IP.String()
}

// Port oddaje numer portu efemerycznego nadanego przez system operacyjny
// odbiornikowi testowego serwera SMTP.
func (o *odbiornikSMTP) Port() int {
	return o.nasluch.Addr().(*net.TCPAddr).Port
}

// Listy oddaje kopię wykazu listów odebranych do tej chwili przez odbiornik
// testowego serwera SMTP na maszynie.
func (o *odbiornikSMTP) Listy() []listOdebrany {
	o.zamek.Lock()
	defer o.zamek.Unlock()
	kopia := make([]listOdebrany, len(o.listy))
	copy(kopia, o.listy)
	return kopia
}

// Ostatni oddaje list odebrany jako ostatni i przerywa sprawdzian, gdy nie
// przyszedł żaden. Brak listu jest tu niepowodzeniem, nie pustym wynikiem.
func (o *odbiornikSMTP) Ostatni(t *testing.T) listOdebrany {
	t.Helper()

	listy := o.Listy()
	if len(listy) == 0 {
		t.Fatal("do skrzynki nie przyszedł żaden list — droga wejścia nie została wysłana")
	}
	return listy[len(listy)-1]
}

// Wyczysc opróżnia skrzynkę odbiornika, żeby kolejny pomiar liczył listy od
// zera, niezależnie od poprzednich rozmów SMTP.
func (o *odbiornikSMTP) Wyczysc() {
	o.zamek.Lock()
	defer o.zamek.Unlock()
	o.listy = nil
}

// rozmawiaj prowadzi jedną rozmowę SMTP i odkłada odebrany list do skrzynki
// odbiornika testowego, bez AUTH i STARTTLS.
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
// kropką z początku wiersza, które wkłada nadawca zgodnie z protokołem SMTP.
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

// adresZPolecenia wyjmuje adres z poleceń MAIL FROM i RCPT TO odebranej
// rozmowy SMTP prowadzonej z nadajnikiem rdzenia.
func adresZPolecenia(polecenie string) string {
	poczatek := strings.Index(polecenie, "<")
	koniec := strings.Index(polecenie, ">")
	if poczatek < 0 || koniec < poczatek {
		return ""
	}
	return polecenie[poczatek+1 : koniec]
}

// ── odczyt drogi z listu ─────────────────────────────────────────────────────

/*
wzorzecKodu wyjmuje kod z szyny listu. Etykieta zaczyna się od „KOD" i kończy
dwukropkiem, a wartością jest sześć cyfr rozdzielonych spacją co trzy —
tak stawia je każdy z szablonów.

Etykieta jest dopasowywana wzorcem, nie porównywana z napisem, bo każdy list
nazywa kod inaczej: „KOD AKTYWACJI", „KOD LOGOWANIA", „KOD RESETU HASŁA".
*/
var wzorzecKodu = regexp.MustCompile(`(?m)^[ \t]*KOD[^:\n]*:[ \t]*([0-9][0-9 ]*[0-9])[ \t]*\r?$`)

/*
drogaZListu wyjmuje z listu materiał do wpisania w oknie. Wyjmowanie idzie po
treści, a nie po wartości z bazy — sprawdzianem jest właśnie to, czy droga
doszła do skrzynki.

Dokument przechodzi wpierw przez rozkodowanie quoted-printable: część tekstowa
listu jedzie tym kodowaniem, więc polskie znaki w etykiecie stoją w nim jako
`=C5=81`, a wiersze dłuższe niż 76 znaków są łamane znakiem `=` na końcu.

Spacje rozdzielające kod są usuwane: w liście stoją po to, żeby dało się kod
przeczytać z ekranu, a rdzeń porównuje sam ciąg cyfr.
*/
func drogaZListu(t *testing.T, list listOdebrany) string {
	t.Helper()

	rozkodowany, err := io.ReadAll(quotedprintable.NewReader(strings.NewReader(list.Dokument)))
	if err != nil {
		// Rozkodowanie jest ułatwieniem, nie warunkiem: dokument bez kodowania
		// czyta się wprost.
		rozkodowany = []byte(list.Dokument)
	}
	if trafienie := wzorzecKodu.FindSubmatch(rozkodowany); trafienie != nil {
		return strings.ReplaceAll(string(trafienie[1]), " ", "")
	}
	t.Fatalf("list nie niesie kodu; dokument:\n%s", rozkodowany)
	return ""
}

// ── pomiar stanu bazy ────────────────────────────────────────────────────────

// liczbaWierszy odpowiada na pytanie, na które odpowiedź komendy odpowiedzieć
// nie może: czy w bazie naprawdę coś zostało zapisane po wykonaniu drogi.
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
// z listu, który po niej przyszedł do skrzynki testowej.
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

// znacznikWSejfie odpowiada, czy bramka pamięta, że powstała bez poczty, i jaki
// adres wtedy zapamiętała, czytając sejf poświadczeń, gdzie ten stan leży.
func znacznikWSejfie(t *testing.T, u uprzazWejscia) (string, bool) {
	t.Helper()

	return dane.NowySejfPlikowy(u.katalog).Odczytaj(u.zycie, bytZnacznikaBezPoczty)
}

// ustawNadajnik zapisuje konto nadawcze platformy komendą `config.set` na
// zasięgu aplikacji, tą samą drogą, którą idzie zmiana w oknie Konfiguracji.
func ustawNadajnik(t *testing.T, u uprzazWejscia, odbiornik *odbiornikSMTP) {
	t.Helper()

	nastawy := []struct {
		klucz   string
		wartosc string
	}{
		{konfig.KluczNadawcaHost, odbiornik.Host()},
		{konfig.KluczNadawcaPort, strconv.Itoa(odbiornik.Port())},
		{konfig.KluczNadawcaAdres, "platforma@danaco.sprawdzian"},
		{konfig.KluczNadawcaNazwa, "Danaco Console"},
		{konfig.KluczNadawcaStartTLS, "false"},
	}
	for _, nastawa := range nastawy {
		wartosc, err := json.Marshal(nastawa.wartosc)
		if err != nil {
			t.Fatalf("nie można złożyć wartości nastawy %s: %v", nastawa.klucz, err)
		}
		wykonajUdana(t, u.rdzen, u.zycie, shared.CommandConfigSet, shared.ConfigSetRequest{
			Key: nastawa.klucz, Value: wartosc, Scope: shared.ConfigScopeApplication,
		}, nil)
	}
}
