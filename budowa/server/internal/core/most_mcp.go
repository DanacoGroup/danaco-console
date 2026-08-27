// Plik generuje wpis mcpServers dla mostu MCP mcp-danaco-pulpit-console z punktów dostępu rodzaju
// maszyna i nadań okna rozmowy, wraz z listą argumentów połączenia ssh.
package core

import (
	"encoding/json"
	"strings"

	"danacoconsole/shared"
)

// Stałe mostu konsoli odpowiadają dosłownie ustawieniom po stronie serwerowej tego mostu MCP tej platformy.
const (
	// PrefiksMostuKonsoli poprzedza identyfikator maszyny docelowej w kluczu wpisu mapy mcpServers tego mostu.
	PrefiksMostuKonsoli = "mcp-danaco-pulpit-console-"
	// sciezkaUruchomieniaMostu to ścieżka skryptu uruchamiającego most konsoli w trybie stdio na tej maszynie.
	sciezkaUruchomieniaMostu = "/opt/danaco/most-konsoli/uruchom-stdio.sh"
	// zmiennaKorzeniMostu przenosi nazwę zmiennej środowiskowej z korzeniami katalogów, poza które most nie wychodzi.
	zmiennaKorzeniMostu = "DANACO_MOST_KORZENIE"
	// rozdzielnikKorzeniMostu rozdziela korzenie katalogów w wartości zmiennej środowiskowej tego mostu konsoli.
	rozdzielnikKorzeniMostu = ":"
	// RodzajWpisuMostu to rodzaj połączenia MCP używany przez ten most przy uruchomieniu procesu narzędzi konsoli.
	RodzajWpisuMostu = "stdio"
	// poleceniePolaczeniaMostu to program zestawiający połączenie ssh z maszyną docelową tego mostu konsoli.
	poleceniePolaczeniaMostu = "ssh"
	// UzytkownikMostuDomyslny obowiązuje, gdy adres punktu dostępu nie niesie nazwy konta systemowego maszyny.
	UzytkownikMostuDomyslny = "ubuntu"
	// PortMostuDomyslny obowiązuje, gdy adres punktu dostępu nie niesie własnego numeru portu sieciowego maszyny.
	PortMostuDomyslny = "22"
	// NazwaMaszynyZastepcza obejmuje punkt bez znaku dopuszczalnego w kluczu wpisu mapy mcpServers tego mostu konsoli.
	NazwaMaszynyZastepcza = "maszyna"
)

// Słownictwa trybu tu nie ma z zamysłem: jest ono własnością mostu, nie kontraktu.

// WpisMostu to jedna pozycja mapy mcpServers opisująca sposób uruchomienia jednego mostu MCP tej konsoli platformy.
type WpisMostu struct {
	Type    string            `json:"type"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// KonfiguracjaMostu to komplet wpisów mcpServers przekazywany procesowi modelu przy jego uruchomieniu na tej platformie.
type KonfiguracjaMostu struct {
	McpServers map[string]WpisMostu `json:"mcpServers"`
}

// NadanieMostu wiąże punkt dostępu z nadaniem okna rozmowy. Punkt mówi, gdzie
// stoi maszyna; nadanie mówi, w jakim trybie i do których korzeni sięga to
// jedno okno.
type NadanieMostu struct {
	Punkt   shared.AccessPoint
	Nadanie shared.AccessGrant
	// ArgumentTrybu jest słowem, jakim ten most nazywa tryb nadania; napis pusty uruchamia sam skrypt.
	ArgumentTrybu string
}

// zlozKonfiguracjeMostu buduje komplet wpisów mcpServers dla nadań okna, pomijając nadania
// i punkty wygaszone oraz punkty rodzaju innego niż maszyna.
func zlozKonfiguracjeMostu(nadania []NadanieMostu) KonfiguracjaMostu {
	konfiguracja := KonfiguracjaMostu{McpServers: map[string]WpisMostu{}}
	for _, nadanie := range uporzadkowaneNadania(nadania) {
		klucz := kluczUnikalny(KluczWpisuMostu(nadanie.Punkt), konfiguracja.McpServers)
		konfiguracja.McpServers[klucz] = zlozWpisMostu(nadanie)
	}
	return konfiguracja
}

// TekstKonfiguracjiMostu zwraca konfigurację w postaci tekstu JSON gotowego do
// zapisania obok procesu modelu. Klucze mapy wychodzą uporządkowane, więc ten
// sam zbiór nadań zawsze daje ten sam tekst.
func TekstKonfiguracjiMostu(nadania []NadanieMostu) (string, error) {
	tresc, err := json.MarshalIndent(zlozKonfiguracjeMostu(nadania), "", "  ")
	if err != nil {
		return "", err
	}
	return string(tresc), nil
}

// zlozWpisMostu buduje jedną pozycję mapy mcpServers dla nadania okna rozmowy przez ten most MCP platformy.
func zlozWpisMostu(nadanie NadanieMostu) WpisMostu {
	adres := AdresMostu(nadanie.Punkt)
	wpis := WpisMostu{
		Type:    RodzajWpisuMostu,
		Command: poleceniePolaczeniaMostu,
		Args:    argumentyPolaczenia(adres, nadanie.ArgumentTrybu),
	}
	if korzenie := KorzenieNadania(nadanie); len(korzenie) > 0 {
		wpis.Env = map[string]string{
			zmiennaKorzeniMostu: strings.Join(korzenie, rozdzielnikKorzeniMostu),
		}
	}
	return wpis
}

// argumentyPolaczenia składa listę argumentów ssh; odwołanie do klucza wchodzi wyłącznie wtedy,
// gdy punkt je niesie, a tryb pusty daje wywołanie samego skryptu bez argumentu.
func argumentyPolaczenia(adres AdresPolaczeniaMostu, tryb string) []string {
	argumenty := []string{
		"-o", "BatchMode=yes",
		"-o", "StrictHostKeyChecking=accept-new",
		"-p", adres.Port,
	}
	if adres.Klucz != "" {
		argumenty = append(argumenty, "-i", adres.Klucz)
	}
	polecenie := sciezkaUruchomieniaMostu
	if strings.TrimSpace(tryb) != "" {
		polecenie += " " + tryb
	}
	return append(argumenty, adres.Uzytkownik+"@"+adres.Host, polecenie)
}
