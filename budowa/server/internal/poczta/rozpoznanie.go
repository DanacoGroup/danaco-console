// Rozpoznanie skrzynek na urządzeniu — obsługa `mail.account.discover`.
// Rdzeń czyta nastawy skonfigurowanego klienta poczty, żeby nie kazać
// Operatorowi przepisywać ręcznie hosta, portu i nazwy użytkownika.
package poczta

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Rozpoznana to skrzynka odczytana z urządzenia. Osobny typ, nie `Nastawy`:
// rozpoznanie nie niesie sekretu, więc pola na sekret tu nie ma.
type Rozpoznana struct {
	Zrodlo           string
	Adres            string
	NazwaWyswietlana string
	Uzytkownik       string
	Protokol         string
	HostOdbioru      string
	PortOdbioru      int
	HostWysylki      string
	PortWysylki      int
}

// Rozpoznaj przechodzi znane miejsca nastaw na urządzeniu i oddaje to, co
// znalazł. Pusty wynik nie jest błędem: urządzenie bez klienta poczty jest
// urządzeniem normalnym, a odmowa kazałaby Operatorowi naprawiać brak, który
// brakiem nie jest.
func Rozpoznaj(katalogDomowy string) []Rozpoznana {
	if strings.TrimSpace(katalogDomowy) == "" {
		return nil
	}
	znalezione := zeSrodowiska()
	znalezione = append(znalezione, zThunderbirda(katalogDomowy)...)
	znalezione = append(znalezione, zEvolution(katalogDomowy)...)
	znalezione = append(znalezione, zMutta(katalogDomowy)...)
	return odsiejPowtorki(znalezione)
}

// wzorzecPrefsThunderbird wyłuskuje pary z `user_pref("klucz", wartość);`.
// Plik `prefs.js` jest kodem JavaScript, ale jego treść to w praktyce sam
// wykaz takich wywołań, czytany wyrażeniem, a nie interpreterem.
var wzorzecPrefsThunderbird = regexp.MustCompile(`user_pref\("([^"]+)",\s*(.*?)\);`)

// zThunderbirda czyta `prefs.js` każdego profilu Thunderbirda. Składanie
// konta jest dwustopniowe: konto wskazuje serwer i tożsamość osobnymi
// kluczami, odczyt idzie po numerach serwerów.
func zThunderbirda(katalogDomowy string) []Rozpoznana {
	var wynik []Rozpoznana
	for _, korzen := range []string{
		filepath.Join(katalogDomowy, ".thunderbird"),
		filepath.Join(katalogDomowy, ".mozilla-thunderbird"),
		filepath.Join(katalogDomowy, "snap", "thunderbird", "common", ".thunderbird"),
	} {
		profile, err := os.ReadDir(korzen)
		if err != nil {
			continue
		}
		for _, profil := range profile {
			if !profil.IsDir() {
				continue
			}
			nastawy := wczytajPrefs(filepath.Join(korzen, profil.Name(), "prefs.js"))
			if len(nastawy) == 0 {
				continue
			}
			wynik = append(wynik, kontaThunderbirda(nastawy)...)
		}
	}
	return wynik
}

// wczytajPrefs zamienia `prefs.js` na mapę klucz-wartość bez cudzysłowów,
// czytając linie postaci user_pref("klucz", "wartość").
func wczytajPrefs(sciezka string) map[string]string {
	plik, err := os.Open(sciezka)
	if err != nil {
		return nil
	}
	defer plik.Close()
	nastawy := map[string]string{}
	czytnik := bufio.NewScanner(plik)
	czytnik.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for czytnik.Scan() {
		for _, para := range wzorzecPrefsThunderbird.FindAllStringSubmatch(czytnik.Text(), -1) {
			nastawy[para[1]] = strings.Trim(strings.TrimSpace(para[2]), `"`)
		}
	}
	return nastawy
}

// kontaThunderbirda składa konta z płaskiej mapy nastaw profilu.
//
// Klucze idą posortowane, bo kolejność przechodzenia mapy w Go jest losowa —
// bez sortowania ten sam profil oddawałby konta w innej kolejności przy każdym
// wywołaniu.
func kontaThunderbirda(nastawy map[string]string) []Rozpoznana {
	var wynik []Rozpoznana
	for _, klucz := range posortowaneKlucze(nastawy) {
		host := nastawy[klucz]
		numer, jest := strings.CutPrefix(klucz, "mail.server.server")
		if !jest || !strings.HasSuffix(numer, ".hostname") {
			continue
		}
		numer = strings.TrimSuffix(numer, ".hostname")
		przedrostek := "mail.server.server" + numer + "."
		rodzaj := strings.ToLower(nastawy[przedrostek+"type"])
		if rodzaj != ProtokolImap && rodzaj != ProtokolPop3 {
			// `none` znaczy w Thunderbirdzie folder lokalny, nie skrzynkę.
			continue
		}
		konto := Rozpoznana{
			Zrodlo:      "thunderbird",
			Protokol:    rodzaj,
			HostOdbioru: host,
			PortOdbioru: liczba(nastawy[przedrostek+"port"]),
			Uzytkownik:  nastawy[przedrostek+"userName"],
			Adres:       nastawy[przedrostek+"userName"],
		}
		dopiszTozsamoscThunderbirda(&konto, nastawy)
		wynik = append(wynik, konto)
	}
	return wynik
}

// dopiszTozsamoscThunderbirda dokłada adres, nazwę wyświetlaną i serwer
// wychodzący z pierwszej tożsamości, jaką profil niesie — "pierwsza" znaczy
// pierwsza po posortowaniu kluczy, bo przechodzenie mapy w Go ma kolejność
// losową.
func dopiszTozsamoscThunderbirda(konto *Rozpoznana, nastawy map[string]string) {
	klucze := posortowaneKlucze(nastawy)
	for _, klucz := range klucze {
		if strings.HasPrefix(klucz, "mail.identity.") && strings.HasSuffix(klucz, ".useremail") {
			if konto.Adres == "" || !strings.Contains(konto.Adres, "@") {
				konto.Adres = nastawy[klucz]
			}
			przedrostek := strings.TrimSuffix(klucz, ".useremail")
			konto.NazwaWyswietlana = nastawy[przedrostek+".fullName"]
			break
		}
	}
	for _, klucz := range klucze {
		if strings.HasPrefix(klucz, "mail.smtpserver.") && strings.HasSuffix(klucz, ".hostname") {
			konto.HostWysylki = nastawy[klucz]
			konto.PortWysylki = liczba(nastawy[strings.TrimSuffix(klucz, ".hostname")+".port"])
			break
		}
	}
}

// posortowaneKlucze oddaje klucze mapy w porządku rosnącym, do powtarzalnego
// przejścia po numerach serwerów i tożsamości.
func posortowaneKlucze(nastawy map[string]string) []string {
	klucze := make([]string, 0, len(nastawy))
	for klucz := range nastawy {
		klucze = append(klucze, klucz)
	}
	sort.Strings(klucze)
	return klucze
}

// zEvolution czyta źródła Evolution z `~/.config/evolution/sources/*.source`.
// Pliki są w formacie INI, więc czyta je ta sama pętla klucz=wartość co mutta.
func zEvolution(katalogDomowy string) []Rozpoznana {
	pliki, err := filepath.Glob(filepath.Join(katalogDomowy, ".config", "evolution", "sources", "*.source"))
	if err != nil {
		return nil
	}
	var wynik []Rozpoznana
	for _, sciezka := range pliki {
		nastawy := wczytajIni(sciezka)
		host := pierwszyNiepusty(nastawy["Host"])
		if host == "" {
			continue
		}
		rodzaj := strings.ToLower(nastawy["BackendName"])
		if rodzaj != ProtokolImap && rodzaj != ProtokolPop3 && rodzaj != "smtp" {
			continue
		}
		wpis := Rozpoznana{
			Zrodlo:           "evolution",
			Protokol:         rodzaj,
			Uzytkownik:       nastawy["User"],
			Adres:            pierwszyNiepusty(nastawy["Address"], nastawy["User"]),
			NazwaWyswietlana: nastawy["DisplayName"],
		}
		if rodzaj == "smtp" {
			wpis.Protokol = ProtokolImap
			wpis.HostWysylki, wpis.PortWysylki = host, liczba(nastawy["Port"])
		} else {
			wpis.HostOdbioru, wpis.PortOdbioru = host, liczba(nastawy["Port"])
		}
		wynik = append(wynik, wpis)
	}
	return wynik
}

// zMutta czyta `.muttrc` i `.config/(neo)mutt/*rc`, biorąc `set folder`
// (skrzynka IMAP), `set smtp_url` i `set from`. Hasła jawne (`set imap_pass=`)
// odczyt pomija.
func zMutta(katalogDomowy string) []Rozpoznana {
	sciezki := []string{filepath.Join(katalogDomowy, ".muttrc"), filepath.Join(katalogDomowy, ".mutt", "muttrc")}
	for _, wzorzec := range []string{"mutt", "neomutt"} {
		dalsze, _ := filepath.Glob(filepath.Join(katalogDomowy, ".config", wzorzec, "*rc"))
		sciezki = append(sciezki, dalsze...)
	}
	var wynik []Rozpoznana
	for _, sciezka := range sciezki {
		tresc, err := os.ReadFile(sciezka)
		if err != nil {
			continue
		}
		wpis := Rozpoznana{Zrodlo: "mutt", Protokol: ProtokolImap}
		for _, linia := range strings.Split(string(tresc), "\n") {
			pola := strings.Fields(strings.TrimSpace(linia))
			if len(pola) < 2 || pola[0] != "set" {
				continue
			}
			klucz, wartosc, _ := strings.Cut(strings.Join(pola[1:], " "), "=")
			wartosc = strings.Trim(strings.TrimSpace(wartosc), `"'`)
			switch strings.TrimSpace(klucz) {
			case "folder", "imap_url":
				wpis.HostOdbioru, wpis.PortOdbioru = zAdresuURL(wartosc, 993)
			case "smtp_url":
				wpis.HostWysylki, wpis.PortWysylki = zAdresuURL(wartosc, 587)
			case "imap_user":
				wpis.Uzytkownik = wartosc
			case "from":
				wpis.Adres = wartosc
			case "realname":
				wpis.NazwaWyswietlana = wartosc
			}
		}
		if wpis.HostOdbioru != "" || wpis.HostWysylki != "" {
			wynik = append(wynik, wpis)
		}
	}
	return wynik
}

// zeSrodowiska czyta zmienne środowiska procesu. Wchodzą do wyniku jako
// pierwsze, bo są wskazaniem najświeższym: Operator, który je ustawił, zrobił to
// dla tego uruchomienia, a nie kiedyś dla klienta poczty.
func zeSrodowiska() []Rozpoznana {
	host := pierwszyNiepusty(os.Getenv("DANACO_POCZTA_IMAP"), os.Getenv("IMAP_HOST"))
	hostWysylki := pierwszyNiepusty(os.Getenv("DANACO_POCZTA_SMTP"), os.Getenv("SMTP_HOST"))
	adres := pierwszyNiepusty(os.Getenv("DANACO_POCZTA_ADRES"), os.Getenv("EMAIL"))
	if host == "" && hostWysylki == "" {
		return nil
	}
	return []Rozpoznana{{
		Zrodlo:      "srodowisko",
		Protokol:    ProtokolImap,
		Adres:       adres,
		Uzytkownik:  pierwszyNiepusty(os.Getenv("DANACO_POCZTA_UZYTKOWNIK"), adres),
		HostOdbioru: host, PortOdbioru: liczba(os.Getenv("DANACO_POCZTA_IMAP_PORT")),
		HostWysylki: hostWysylki, PortWysylki: liczba(os.Getenv("DANACO_POCZTA_SMTP_PORT")),
	}}
}

// wczytajIni czyta plik klucz=wartość, pomijając sekcje i komentarze. Klucz
// powtórzony wygrywa pierwszym wystąpieniem — sekcje dalsze w plikach Evolution
// opisują zwykle inne role tego samego źródła.
func wczytajIni(sciezka string) map[string]string {
	tresc, err := os.ReadFile(sciezka)
	if err != nil {
		return nil
	}
	nastawy := map[string]string{}
	for _, linia := range strings.Split(string(tresc), "\n") {
		linia = strings.TrimSpace(linia)
		if linia == "" || strings.HasPrefix(linia, "[") || strings.HasPrefix(linia, "#") {
			continue
		}
		klucz, wartosc, jest := strings.Cut(linia, "=")
		if !jest {
			continue
		}
		klucz = strings.TrimSpace(klucz)
		if _, byl := nastawy[klucz]; !byl {
			nastawy[klucz] = strings.TrimSpace(wartosc)
		}
	}
	return nastawy
}

// zAdresuURL wyłuskuje host i port z zapisu w rodzaju `imaps://kto@host:993/`,
// tak jak zapisuje go mutt w kluczu `folder`.
func zAdresuURL(wartosc string, domyslnyPort int) (string, int) {
	if _, reszta, jest := strings.Cut(wartosc, "://"); jest {
		wartosc = reszta
	}
	if _, reszta, jest := strings.Cut(wartosc, "@"); jest {
		wartosc = reszta
	}
	wartosc = strings.TrimSuffix(strings.TrimSpace(wartosc), "/")
	if host, port, jest := strings.Cut(wartosc, ":"); jest {
		if numer := liczba(port); numer > 0 {
			return host, numer
		}
		return host, domyslnyPort
	}
	return wartosc, domyslnyPort
}

// odsiejPowtorki usuwa wpisy o tym samym adresie i serwerze — ten sam profil
// bywa widziany dwoma drogami (np. Thunderbird w katalogu domowym i w snapie).
func odsiejPowtorki(wpisy []Rozpoznana) []Rozpoznana {
	widziane := map[string]bool{}
	wynik := make([]Rozpoznana, 0, len(wpisy))
	for _, w := range wpisy {
		klucz := w.Adres + "|" + w.HostOdbioru + "|" + w.HostWysylki
		if widziane[klucz] {
			continue
		}
		widziane[klucz] = true
		wynik = append(wynik, w)
	}
	return wynik
}

// liczba czyta liczbę dodatnią albo oddaje zero. Zero znaczy „nie wskazano",
// więc wołający bierze wtedy port domyślny protokołu.
func liczba(tekst string) int {
	numer, err := strconv.Atoi(strings.TrimSpace(tekst))
	if err != nil || numer <= 0 {
		return 0
	}
	return numer
}
