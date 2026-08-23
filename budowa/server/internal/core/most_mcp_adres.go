package core

import (
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// AdresPolaczeniaMostu to rozłożony adres maszyny: konto, nazwa hosta, port
// i odwołanie do klucza. Odwołanie, nie treść klucza — poświadczenie zostaje
// w sejfie i nigdy nie przechodzi przez tę strukturę.
type AdresPolaczeniaMostu struct {
	Uzytkownik string
	Host       string
	Port       string
	Klucz      string
}

// AdresMostu wyprowadza adres połączenia z punktu dostępu.
//
// Pierwszeństwo ma pole endpoint, bo tylko ono niesie konto i port; pole host
// uzupełnia nazwę maszyny, gdy endpoint jej nie podał. Brak konta i brak portu
// schodzą na wartości domyślne mostu, nie na odmowę złożenia wpisu.
func AdresMostu(punkt shared.AccessPoint) AdresPolaczeniaMostu {
	adres := AdresPolaczeniaMostu{
		Uzytkownik: UzytkownikMostuDomyslny,
		Port:       PortMostuDomyslny,
		Host:       tekstPunktu(punkt.Host),
		Klucz:      tekstPunktu(punkt.CredentialRef),
	}
	rozlozAdres(&adres, tekstPunktu(punkt.Endpoint))
	if adres.Host == "" {
		adres.Host = tekstPunktu(punkt.Host)
	}
	if adres.Host == "" {
		adres.Host = identyfikatorMaszyny(punkt)
	}
	return adres
}

// rozlozAdres rozbiera zapis [ssh://][konto@]host[:port] na części. Zapis pusty
// zostawia adres bez zmian.
func rozlozAdres(adres *AdresPolaczeniaMostu, endpoint string) {
	reszta := strings.TrimSpace(endpoint)
	reszta = strings.TrimPrefix(reszta, "ssh://")
	if reszta == "" {
		return
	}
	if podzial := strings.LastIndex(reszta, "@"); podzial >= 0 {
		if konto := strings.TrimSpace(reszta[:podzial]); konto != "" {
			adres.Uzytkownik = konto
		}
		reszta = reszta[podzial+1:]
	}
	if podzial := strings.LastIndex(reszta, ":"); podzial >= 0 {
		if port := strings.TrimSpace(reszta[podzial+1:]); portPoprawny(port) {
			adres.Port = port
			reszta = reszta[:podzial]
		}
	}
	if host := strings.TrimSpace(reszta); host != "" {
		adres.Host = host
	}
}

// portPoprawny odpowiada, czy tekst jest numerem portu. Zapis niebędący
// numerem nie jest portem — bywa częścią adresu IPv6 albo ścieżki.
func portPoprawny(tekst string) bool {
	numer, err := strconv.Atoi(tekst)
	return err == nil && numer > 0 && numer <= 65535
}

// KluczWpisuMostu zwraca klucz pozycji `mcpServers` dla punktu dostępu:
// stały prefiks mostu konsoli i identyfikator maszyny.
func KluczWpisuMostu(punkt shared.AccessPoint) string {
	return PrefiksMostuKonsoli + identyfikatorMaszyny(punkt)
}

// identyfikatorMaszyny wyprowadza człon rozróżniający klucz wpisu. Pierwszeństwo
// ma nazwa mostu po stronie klienta, bo to ona jest nazwą własną maszyny
// w konfiguracji Operatora; dalej idzie nazwa hosta, a na końcu identyfikator
// wiersza punktu dostępu — zawsze niepusty.
//
// Nazwa mostu bywa podana w całości, razem z prefiksem; prefiks jest wtedy
// obcinany, żeby klucz nie urósł do postaci powtórzonej.
func identyfikatorMaszyny(punkt shared.AccessPoint) string {
	kandydaci := []string{
		strings.TrimPrefix(tekstPunktu(punkt.BridgeName), PrefiksMostuKonsoli),
		tekstPunktu(punkt.Host),
		punkt.Id,
	}
	for _, kandydat := range kandydaci {
		if nazwa := nazwaKluczaMostu(kandydat); nazwa != "" {
			return nazwa
		}
	}
	return NazwaMaszynyZastepcza
}

// nazwaKluczaMostu sprowadza tekst do postaci bezpiecznej dla klucza wpisu:
// małe litery, cyfry i kreska. Klucz trafia do konfiguracji MCP i do nazw
// narzędzi widocznych dla modelu, więc znak spoza tego zbioru jest zamieniany
// na kreskę, a nie przenoszony.
func nazwaKluczaMostu(tekst string) string {
	budowana := strings.Builder{}
	for _, znak := range strings.ToLower(strings.TrimSpace(tekst)) {
		switch {
		case znak >= 'a' && znak <= 'z', znak >= '0' && znak <= '9':
			budowana.WriteRune(znak)
		default:
			budowana.WriteRune('-')
		}
	}
	return strings.Trim(scalKreski(budowana.String()), "-")
}

// scalKreski zamienia ciąg kresek na jedną kreskę.
func scalKreski(tekst string) string {
	for strings.Contains(tekst, "--") {
		tekst = strings.ReplaceAll(tekst, "--", "-")
	}
	return tekst
}

// tekstPunktu zwraca przyciętą treść pola opcjonalnego kontraktu. Odczyt
// wskaźnika należy do wartoscTekstu z adaptera ustawień — tu dokładane jest
// wyłącznie przycięcie białych znaków.
func tekstPunktu(wskaznik *string) string {
	return strings.TrimSpace(wartoscTekstu(wskaznik))
}
