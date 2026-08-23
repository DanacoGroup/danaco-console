// Generator wpisu `mcpServers` dla mostu MCP „mcp-danaco-pulpit-console".
//
// Wejściem są punkty dostępu rodzaju maszyna (shared.AccessPointKindMcpBridge)
// wraz z nadaniem dostępu okna rozmowy. Wyjściem jest gotowa konfiguracja MCP
// przekazywana procesowi modelu. Tryb uprawnień mostu bierze się z nadania,
// nie z ustawienia globalnego: nadanie żyje per okno rozmowy.
//
// Cały plik to funkcje czyste: dane wchodzą, tekst wychodzi. Nie ma tu ani
// jednego połączenia sieciowego, ani jednego wywołania ssh — dzięki temu
// kształt wpisu daje się sprawdzić bez maszyny po drugiej stronie.
//
// Kształt wpisu odpowiada skryptowi /opt/danaco/most-konsoli/uruchom-stdio.sh:
// tryb wchodzi argumentem pozycyjnym, a korzenie zmienną DANACO_MOST_KORZENIE
// rozdzieloną dwukropkiem.
package core

import (
	"encoding/json"
	"strings"

	"danacoconsole/shared"
)

// Stałe mostu konsoli. Odpowiadają dosłownie stronie serwerowej mostu.
const (
	// PrefiksMostuKonsoli poprzedza identyfikator maszyny w kluczu wpisu.
	// Nazwa mostu jest stała po stronie serwera, więc klucze wpisów różnią się
	// wyłącznie identyfikatorem maszyny.
	PrefiksMostuKonsoli = "mcp-danaco-pulpit-console-"
	// sciezkaUruchomieniaMostu to skrypt uruchamiający most w trybie stdio.
	sciezkaUruchomieniaMostu = "/opt/danaco/most-konsoli/uruchom-stdio.sh"
	// zmiennaKorzeniMostu przenosi korzenie, poza które most nie wychodzi.
	zmiennaKorzeniMostu = "DANACO_MOST_KORZENIE"
	// rozdzielnikKorzeniMostu rozdziela korzenie w zmiennej środowiskowej.
	rozdzielnikKorzeniMostu = ":"
	// RodzajWpisuMostu to rodzaj połączenia MCP używany przez most.
	RodzajWpisuMostu = "stdio"
	// poleceniePolaczeniaMostu to program zestawiający połączenie.
	poleceniePolaczeniaMostu = "ssh"
	// UzytkownikMostuDomyslny obowiązuje, gdy adres punktu nie niesie nazwy konta.
	UzytkownikMostuDomyslny = "ubuntu"
	// PortMostuDomyslny obowiązuje, gdy adres punktu nie niesie portu.
	PortMostuDomyslny = "22"
	// NazwaMaszynyZastepcza obejmuje punkt, z którego identyfikatora nie
	// zostaje ani jeden znak dopuszczalny w kluczu wpisu.
	NazwaMaszynyZastepcza = "maszyna"
)

/* Słownictwa trybu („odczyt", „zapis") tu nie ma z zamysłem. Jest ono
   własnością konkretnego mostu, a nie typu kontraktu, więc mieszka w tabeli
   `argument_trybu_mostu` i dojeżdża tutaj polem
   NadanieMostu.ArgumentTrybu. Zapisanie go drugi raz w kodzie sprawiłoby, że
   most o innym słownictwie wymagałby zmiany rdzenia zamiast wiersza
*/

// WpisMostu to jedna pozycja mapy `mcpServers`.
type WpisMostu struct {
	Type    string            `json:"type"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// KonfiguracjaMostu to komplet wpisów przekazywany procesowi modelu.
type KonfiguracjaMostu struct {
	McpServers map[string]WpisMostu `json:"mcpServers"`
}

// NadanieMostu wiąże punkt dostępu z nadaniem okna rozmowy. Punkt mówi, gdzie
// stoi maszyna; nadanie mówi, w jakim trybie i do których korzeni sięga to
// jedno okno.
type NadanieMostu struct {
	Punkt   shared.AccessPoint
	Nadanie shared.AccessGrant
	// ArgumentTrybu jest słowem, jakim ten most nazywa tryb nadania. Wypełnia je
	// warstwa danych z tabeli `argument_trybu_mostu` — kontrakt nie
	// niesie słownictwa konkretnego mostu, bo jest ono własnością wiersza,
	// nie typu. Napis pusty znaczy: uruchomić sam skrypt, bez argumentu; most
	// `mcp-danaco-pulpit-console` wchodzi wtedy w odczyt, czyli w wariant
	// bezpieczniejszy.
	ArgumentTrybu string
}

// zlozKonfiguracjeMostu buduje komplet wpisów `mcpServers` dla nadań okna.
//
// Pomijane są nadania i punkty wygaszone oraz punkty rodzaju innego niż
// maszyna — katalog lokalny nie jest mostem i nie ma czego uruchamiać przez
// ssh. Puste wejście daje pustą, lecz poprawną mapę: brak dostępów nie jest
// błędem konfiguracji.
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

// zlozWpisMostu buduje jedną pozycję `mcpServers`.
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

// argumentyPolaczenia składa listę argumentów ssh. Odwołanie do klucza wchodzi
// wyłącznie wtedy, gdy punkt je niesie — brak odwołania zostawia rozstrzygnięcie
// konfiguracji ssh zamiast wstawiać ścieżkę zmyśloną.
//
// Tryb pusty daje wywołanie samego skryptu: most bez argumentu wchodzi w odczyt,
// więc dopisanie tu wartości domyślnej byłoby drugim zapisem tej samej reguły.
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
