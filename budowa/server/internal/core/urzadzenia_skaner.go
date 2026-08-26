// Odpowiedzialność pliku: warstwa urządzeń wejściowych (skanera) rozdzielona po
// systemie operacyjnym — wykaz urządzeń i pobranie obrazu drogą właściwą dla
// maszyny, na której stoi rdzeń.
//
// ── Dlaczego JEDEN plik z rozstrzygnięciem po `runtime.GOOS`, a nie `_windows.go` ──
// Rozdzielenie warunkiem budowy (`urzadzenia_skaner_windows.go` /
// `_linux.go`) wygląda porządniej, ale ma tu jedną wadę rozstrzygającą:
// `go build ./...` na maszynie budowy (Linux) NIE SKOMPILOWAŁBY drogi WIA ani
// razu. Droga Windows przestałaby się kompilować przy pierwszej zmianie
// sąsiedniego pliku i nikt by tego nie zobaczył aż do wydania instalki
// natywnej — czyli dokładnie w chwili, w której nie ma już czasu na naprawę.
// Rozstrzygnięcie po `runtime.GOOS` trzyma oba warianty pod jednym sprawdzianem
// kompilacji, a kosztem jest kilka bajtów martwego kodu w wydaniu. To ta sama
// droga, którą rdzeń rozróżnia system już dziś (`urzadzenia.go` linia 51,
// `adapter_modul_developer_okno.go` linia 141) — nie druga jej odmiana.
//
// Warunek budowy będzie właściwy dopiero wtedy, gdy droga Windows sięgnie po
// bibliotekę wołającą COM z Go. Dziś woła `pwsh`, więc nie ma czego chronić
// warunkiem budowy: kod jest zwykłym napisem i kompiluje się wszędzie.
//
// ── Dlaczego WIA przez `pwsh`, a nie TWAIN ─────────────────────────────────
// TWAIN wymaga okna i pętli komunikatów — z procesu serwera bez pulpitu nie
// wystartuje. WIA jest warstwą systemową Windows dostępną przez COM, a jedyną
// drogą do COM, którą rdzeń ma bez wkompilowanej biblioteki, jest PowerShell.
// `pwsh` stoi już w wykazie zależności rdzenia (`zaleznosci_zewnetrzne.go`,
// pozycja `narzedziePowerShell`), więc sonda startowa mówi o jego braku sama.
//
// Zasada bezwzględna: program zewnętrzny idzie WYŁĄCZNIE przez
// `zewnetrzne.Wolaj` — nigdy `exec.Command`. Pilnuje tego
// `zapora_procesow_rdzenia_test.go`.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// Nazwy systemów, po których rozdziela się warstwa urządzeń. Stałe zamiast
// napisów w gałęziach, bo ten sam warunek pada w kilku miejscach pliku, a
// literówka w jednym z nich dałaby cichą odmowę na właściwym systemie.
const (
	systemWindows = "windows"
	systemLinux   = "linux"
)

// granicaSkanowaniaUrzadzenia jest granicą czasu JEDNEGO przebiegu skanera.
// Osobna od granicy arsenału (2 minuty), bo skan strony A4 przy 600 punktach na
// cal na wolnym urządzeniu USB trwa dłużej niż rozpakowanie archiwum, a granica
// urwana w połowie zostawia plik obrazu bez końca.
const granicaSkanowaniaUrzadzenia = 5 * time.Minute

// granicaWykazuUrzadzen jest granicą czasu samego wykazu. Krótka celowo:
// odpytanie warstwy skanera to czynność natychmiastowa, a gdy nie jest —
// urządzenie zwisło i Operator ma to usłyszeć, zanim odejdzie od okna.
const granicaWykazuUrzadzen = 45 * time.Second

// rozdzielczoscSkanowaniaDomyslna to punkty na cal brane, gdy żądanie nie mówi
// nic. 300 jest progiem czytelności dla rozpoznania pisma — niżej Tesseract
// gubi litery, wyżej rośnie czas skanu bez zysku dla tekstu.
const rozdzielczoscSkanowaniaDomyslna = 300

// narzedzieSkaneraWia jest tym samym programem, co `narzedziePowerShell`
// (`pwsh`), ale pod nazwą mówiącą o TEJ czynności. Operator, któremu odmówiono
// skanowania, ma przeczytać, czego brakuje do skanowania — nie do analizy
// skryptów. Wykaz zależności pilnuje programu, nie nazwy czytelnej, więc ta
// deklaracja nie jest drugim wpisem obok tamtego.
var narzedzieSkaneraWia = zewnetrzne.Narzedzie{
	Nazwa:   "PowerShell 7 (droga WIA do skanera Windows)",
	Program: "pwsh",
	Pakiet:  "PowerShell 7 (winget install Microsoft.PowerShell)",
}

// ── Wykaz urządzeń ──────────────────────────────────────────────────────────

// wykazSkanerow oddaje urządzenia wejściowe widziane przez rdzeń — drogą
// właściwą dla systemu. Pusty wykaz oddaje TYLKO wtedy, gdy warstwa skanera
// odpowiedziała i nic nie zgłosiła; brak warstwy jest odmową, nie pustką.
func (a *adapterStudia) wykazSkanerow(ctx context.Context) ([]shared.StudioInputDevice, error) {
	switch runtime.GOOS {
	case systemLinux:
		wyjscie, err := a.wolajUrzadzenie(ctx, narzedzieSkanera, []string{"-L"},
			granicaWykazuUrzadzen)
		if err != nil {
			return nil, err
		}
		return odczytajUrzadzenia(string(wyjscie)), nil
	case systemWindows:
		wyjscie, err := a.wolajUrzadzenie(ctx, narzedzieSkaneraWia,
			skryptPowerShell(skryptWykazuWia), granicaWykazuUrzadzen)
		if err != nil {
			return nil, bladWarstwyWia(err)
		}
		return odczytajUrzadzeniaWia(string(wyjscie))
	default:
		return nil, odmowaSkaneraNaTymSystemie()
	}
}

// odmowaSkaneraNaTymSystemie nazywa brak drogi na systemie, dla którego rdzeń
// warstwy skanera nie ma. Nie pusty wykaz: pusty wykaz znaczy „szukałem i nic
// nie ma", a tu rdzeń nie ma czym szukać.
func odmowaSkaneraNaTymSystemie() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Studio: rdzeń nie ma warstwy skanera na systemie "+runtime.GOOS+
			". Warstwy, które zna, to SANE (`scanimage`) na Linuksie i WIA przez "+
			"PowerShell (`pwsh`) na Windowsie — na tym systemie nie ma ani jednej "+
			"z nich. Droga, która działa: zeskanuj materiał programem systemu i "+
			"dołóż plik do kolejki przez `studio.ingest.queue.add`."))
}

// odczytajUrzadzeniaWia czyta wykaz oddany przez skrypt WIA. Skrypt oddaje JSON,
// a nie wiersze do parsowania — wiersze WIA nie mają ustalonego kształtu, więc
// każde ich czytanie byłoby zgadywaniem.
func odczytajUrzadzeniaWia(wyjscie string) ([]shared.StudioInputDevice, error) {
	tresc := strings.TrimSpace(wyjscie)
	if tresc == "" {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: warstwa WIA nie oddała żadnej odpowiedzi na pytanie o "+
				"urządzenia. To nie znaczy „nie ma skanerów” — znaczy, że skrypt "+
				"PowerShell nie doszedł do końca. Naprawa: sprawdzić, czy `pwsh` "+
				"startuje na tej maszynie i czy usługa Windows Image Acquisition (stisvc) "+
				"jest uruchomiona."))
	}

	var zapis []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(tresc), &zapis); err != nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: odpowiedzi warstwy WIA nie da się odczytać jako wykazu "+
				"urządzeń ("+err.Error()+"). Rdzeń NIE oddaje w tym miejscu pustego "+
				"wykazu, bo nie wie, czy skanera nie ma, czy tylko nie zrozumiał "+
				"odpowiedzi. Treść odpowiedzi: "+skrocDoPodgladu(tresc)))
	}

	urzadzenia := []shared.StudioInputDevice{}
	for _, wiersz := range zapis {
		kod := strings.TrimSpace(wiersz.Id)
		if kod == "" {
			continue
		}
		nazwa := strings.TrimSpace(wiersz.Name)
		if nazwa == "" {
			nazwa = kod
		}
		urzadzenia = append(urzadzenia, shared.StudioInputDevice{
			Id: kod, Name: nazwa, Kind: shared.StudioInputDeviceKindSkaner,
		})
	}
	return urzadzenia, nil
}

// ── Pobranie obrazu ─────────────────────────────────────────────────────────

// zamowienieSkanu niesie nastawy jednego pobrania. Struktura kontraktu nie
// jedzie w głąb, bo warstwa urządzeń ma być czytelna bez kontraktu — i bo tryb
// barwny trzeba przełożyć na dwa różne słowniki (SANE i WIA).
type zamowienieSkanu struct {
	// Urzadzenie jest identyfikatorem z wykazu. Puste znaczy urządzenie
	// domyślne warstwy, a nie „dowolne".
	Urzadzenie string
	// Rozdzielczosc w punktach na cal.
	Rozdzielczosc int
	// TrybBarwny jest słowem Operatora („color", „gray", „bw"). Puste zostawia
	// nastawę urządzenia nietkniętą.
	TrybBarwny string
	// Stron mówi, ile kartek pobrać z podajnika. Mniej niż 1 znaczy jedną.
	Stron int
}

// skanujUrzadzenie pobiera obraz (albo obrazy) i oddaje ścieżki plików zapisanych
// na dysku rdzenia. Zapis idzie do katalogu roboczego okna, bo tam sięga
// izolacja — plik poza obszarem byłby materiałem, którego kolejka nie ma prawa
// przeczytać.
func (a *adapterStudia) skanujUrzadzenie(ctx context.Context,
	zamowienie zamowienieSkanu) ([]string, error) {

	katalog, err := a.katalogSkanow()
	if err != nil {
		return nil, err
	}
	stron := zamowienie.Stron
	if stron < 1 {
		stron = 1
	}
	rozdzielczosc := zamowienie.Rozdzielczosc
	if rozdzielczosc < 1 {
		rozdzielczosc = rozdzielczoscSkanowaniaDomyslna
	}

	switch runtime.GOOS {
	case systemLinux:
		return a.skanujSane(ctx, katalog, zamowienie, rozdzielczosc, stron)
	case systemWindows:
		return a.skanujWia(ctx, katalog, zamowienie, rozdzielczosc, stron)
	default:
		return nil, odmowaSkaneraNaTymSystemie()
	}
}

// skanujSane pobiera obrazy przez SANE. Obraz przychodzi WYJŚCIEM programu, nie
// plikiem: `scanimage --format=png` pisze PNG na wyjście standardowe, a rdzeń
// zapisuje bajty sam. Przekierowanie do pliku wymagałoby powłoki, a powłoka jest
// drugą drogą uruchomienia procesu obok `zewnetrzne.Wolaj`.
func (a *adapterStudia) skanujSane(ctx context.Context, katalog string,
	zamowienie zamowienieSkanu, rozdzielczosc, stron int) ([]string, error) {

	sciezki := make([]string, 0, stron)
	for numer := 1; numer <= stron; numer++ {
		argumenty := []string{"--format=png", "--resolution", strconv.Itoa(rozdzielczosc)}
		if kod := strings.TrimSpace(zamowienie.Urzadzenie); kod != "" {
			argumenty = append(argumenty, "-d", kod)
		}
		if tryb := trybBarwnySane(zamowienie.TrybBarwny); tryb != "" {
			argumenty = append(argumenty, "--mode", tryb)
		}

		wyjscie, err := a.wolajUrzadzenie(ctx, narzedzieSkanera, argumenty,
			granicaSkanowaniaUrzadzenia)
		if err != nil {
			if numer > 1 {
				// Podajnik pustego arkusza nie poda i SANE odmawia — strony
				// pobrane przed tym momentem są prawdziwe i wracają do kolejki.
				return sciezki, nil
			}
			return nil, a.odmowaSkanuSane(ctx, err)
		}
		if len(wyjscie) == 0 {
			return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Studio: SANE zakończył skanowanie bez ani jednego bajtu obrazu. "+
					"Naprawa: sprawdzić `scanimage -L`, czy urządzenie odpowiada, i czy "+
					"proces rdzenia ma prawo do jego węzła (grupa `scanner`)."))
		}
		sciezka, err := zapiszSkan(katalog, numer, wyjscie)
		if err != nil {
			return nil, err
		}
		sciezki = append(sciezki, sciezka)
	}
	return sciezki, nil
}

// odmowaSkanuSane rozstrzyga, czy skan nie doszedł do skutku z BRAKU URZĄDZENIA,
// czy z innej przyczyny — i dopiero wtedy nazywa brak.
//
// Bez tego rozstrzygnięcia obie sytuacje wychodziły jednym zdaniem: `scanimage`
// kończy się kodem 1 przy każdej przyczynie, a odmowa arsenału przekłada kod
// niezerowy na usterkę wewnętrzną wraz ze zrzutem procesu. Operator, który po
// prostu nie ma podłączonego skanera, dostawał więc `internal_error` — kod
// mówiący „usterka rdzenia, zgłoś ją" — zamiast zdania o tym, czego brakuje
// i co z tym zrobić. Droga WIA rozróżnia te dwie rzeczy od początku
// (`odmowaWia`, przypadek `BRAK-URZADZENIA`) i to samo należy się drodze SANE:
// stan maszyny jest ten sam, więc i odpowiedź ma być ta sama.
//
// Rozstrzyga PYTANIEM O WYKAZ, nie czytaniem diagnostyki programu. Zdanie, które
// `scanimage` mówi o sobie, jest napisem obcego programu — rdzeń nie ma prawa
// opierać kodu odmowy na tym, że napis nie zmieni się przy następnym wydaniu.
// Wykaz jest drogą własną rdzenia i odpowiada wprost na pytanie, które tu padło.
//
// Wykaz pusty znaczy „nie ma czego skanować": to nie awaria drogi, tylko brak
// urządzenia, więc kod jest `not_found`. Wykaz niepusty albo niedostępny
// zostawia odmowę pierwotną — rdzeń nie wie wtedy nic ponad to, co powiedział
// program, a odmowa zgadnięta byłaby gorsza od surowej.
//
// Pytanie idzie WYŁĄCZNIE po nieudanym skanie, więc droga udana nie płaci za nie
// ani jednym wywołaniem.
func (a *adapterStudia) odmowaSkanuSane(ctx context.Context, pierwotna error) error {
	urzadzenia, blad := a.wykazSkanerow(ctx)
	if blad != nil || len(urzadzenia) > 0 {
		return pierwotna
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Studio: warstwa SANE (`scanimage`) odpowiedziała i nie widzi ani jednego "+
			"skanera. To nie jest awaria drogi ani usterka rdzenia, lecz brak urządzenia "+
			"po stronie maszyny Operatora. Naprawa: podłączyć skaner, włączyć go i sprawdzić "+
			"wykaz komendą `studio.ingest.device.list`; jeżeli urządzenie tam jest, procesowi "+
			"rdzenia brakuje prawa do jego węzła (grupa `scanner`). Droga, która działa bez "+
			"skanera: zeskanuj materiał programem systemu i dołóż plik komendą "+
			"`studio.ingest.queue.add`. Diagnostyka warstwy: "+pierwotna.Error()))
}

// trybBarwnySane przekłada słowo Operatora na tryb SANE. Nierozpoznane słowo
// jedzie do urządzenia bez zmiany — słowniki trybów są backendowe i rdzeń nie
// ma prawa zamieniać nastawy, której nie rozumie, na własną.
func trybBarwnySane(tryb string) string {
	switch strings.ToLower(strings.TrimSpace(tryb)) {
	case "":
		return ""
	case "color", "colour", "kolor", "barwny", "rgb":
		return "Color"
	case "gray", "grey", "grayscale", "greyscale", "szarosci", "szarości":
		return "Gray"
	case "bw", "mono", "lineart", "black-and-white", "czarnobialy":
		return "Lineart"
	default:
		return strings.TrimSpace(tryb)
	}
}

// skanujWia pobiera obrazy przez WIA. Tu obraz NIE wraca wyjściem procesu:
// PowerShell oddaje bajty obrazu jako tekst i po drodze psuje je znakowaniem, a
// obraz zapisany przez WIA metodą `SaveFile` jest tym samym obrazem bez ryzyka.
// Skrypt zapisuje więc pliki do katalogu wskazanego przez rdzeń i oddaje na
// wyjściu ich ścieżki jako JSON.
func (a *adapterStudia) skanujWia(ctx context.Context, katalog string,
	zamowienie zamowienieSkanu, rozdzielczosc, stron int) ([]string, error) {

	skrypt := skryptSkanuWia(katalog, strings.TrimSpace(zamowienie.Urzadzenie),
		rozdzielczosc, trybBarwnyWia(zamowienie.TrybBarwny), stron)

	wyjscie, err := a.wolajUrzadzenie(ctx, narzedzieSkaneraWia, skryptPowerShell(skrypt),
		granicaSkanowaniaUrzadzenia)
	if err != nil {
		return nil, bladWarstwyWia(err)
	}

	var zapis struct {
		Pliki []string `json:"pliki"`
		Blad  string   `json:"blad"`
	}
	tresc := strings.TrimSpace(string(wyjscie))
	if tresc == "" {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: warstwa WIA nie odpowiedziała na żądanie skanowania. "+
				"Naprawa: sprawdzić, czy usługa Windows Image Acquisition (stisvc) "+
				"jest uruchomiona i czy skaner ma zainstalowany sterownik WIA."))
	}
	if err := json.Unmarshal([]byte(tresc), &zapis); err != nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: odpowiedzi warstwy WIA na skanowanie nie da się odczytać ("+
				err.Error()+"). Treść odpowiedzi: "+skrocDoPodgladu(tresc)))
	}
	if zapis.Blad != "" {
		return nil, odmowaWia(zapis.Blad)
	}
	if len(zapis.Pliki) == 0 {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: warstwa WIA zakończyła skanowanie bez ani jednego pliku "+
				"obrazu, nie nazywając powodu. Rdzeń nie oddaje tu pustej kolejki, bo "+
				"nie wie, czy kartki nie było, czy transfer się nie udał."))
	}
	return zapis.Pliki, nil
}

// trybBarwnyWia przekłada słowo Operatora na WIA_IPS_CUR_INTENT: 1 to barwa,
// 2 to szarości, 4 to tekst czarno-biały. Nierozpoznane słowo zostawia nastawę
// urządzenia — zero znaczy „nie dotykaj".
func trybBarwnyWia(tryb string) int {
	switch strings.ToLower(strings.TrimSpace(tryb)) {
	case "color", "colour", "kolor", "barwny", "rgb":
		return 1
	case "gray", "grey", "grayscale", "greyscale", "szarosci", "szarości":
		return 2
	case "bw", "mono", "lineart", "black-and-white", "czarnobialy":
		return 4
	default:
		return 0
	}
}

// ── Skrypty warstwy WIA ─────────────────────────────────────────────────────

// skryptPowerShell składa wywołanie `pwsh` z jednym skryptem. Te same przełączniki,
// którymi woła PowerShell moduł Terminal (`adapter_modul_terminal_analiza.go`):
// bez logo, bez profilu Operatora i bez interakcji — profil maszyny nie ma prawa
// zmienić wyniku czynności rdzenia.
func skryptPowerShell(skrypt string) []string {
	return []string{"-NoLogo", "-NoProfile", "-NonInteractive", "-Command", skrypt}
}

// skryptWykazuWia pyta warstwę WIA o urządzenia rodzaju skaner (typ 1) i oddaje
// je jako JSON. Brak warstwy WIA jest rozpoznany OSOBNO od braku urządzeń, bo to
// dwie różne odpowiedzi dla Operatora.
const skryptWykazuWia = `$ErrorActionPreference='Stop';` +
	`try{$m=New-Object -ComObject WIA.DeviceManager}catch{Write-Output '[]';` +
	`[Console]::Error.WriteLine('WARSTWA-WIA-NIEDOSTEPNA: '+$_.Exception.Message);exit 3};` +
	`$w=New-Object System.Collections.ArrayList;` +
	`foreach($d in $m.DeviceInfos){if($d.Type -ne 1){continue};$n=$d.DeviceID;` +
	`try{$n=$d.Properties.Item('Name').Value}catch{};` +
	`[void]$w.Add([pscustomobject]@{id=$d.DeviceID;name=$n})};` +
	`ConvertTo-Json -InputObject @($w) -Compress`

// skryptSkanuWia składa skrypt pobrania stron. Nastawy wchodzą liczbami i
// ścieżką katalogu, nigdy tekstem od Operatora — identyfikator urządzenia jest
// jedyną wartością zmienną i idzie w apostrofach z podwojonym apostrofem, bo tak
// PowerShell odczytuje napis dosłowny.
func skryptSkanuWia(katalog, urzadzenie string, rozdzielczosc, intencja, stron int) string {
	var skrypt strings.Builder
	skrypt.WriteString(`$ErrorActionPreference='Stop';`)
	skrypt.WriteString(`$wynik=[pscustomobject]@{pliki=New-Object System.Collections.ArrayList;blad=''};`)
	skrypt.WriteString(`try{$m=New-Object -ComObject WIA.DeviceManager}catch{` +
		`$wynik.blad='WARSTWA-WIA-NIEDOSTEPNA: '+$_.Exception.Message;` +
		`ConvertTo-Json -InputObject $wynik -Compress;exit 0};`)
	skrypt.WriteString(`$szukany=` + napisPowerShell(urzadzenie) + `;`)
	skrypt.WriteString(`$info=$null;foreach($d in $m.DeviceInfos){if($d.Type -ne 1){continue};` +
		`if($szukany -eq '' -or $d.DeviceID -eq $szukany){$info=$d;break}};`)
	skrypt.WriteString(`if($info -eq $null){` +
		`$wynik.blad='BRAK-URZADZENIA: warstwa WIA nie widzi skanera' + ` +
		`$(if($szukany -ne ''){' o identyfikatorze '+$szukany}else{''});` +
		`ConvertTo-Json -InputObject $wynik -Compress;exit 0};`)
	skrypt.WriteString(`try{$dev=$info.Connect()}catch{` +
		`$wynik.blad='BRAK-STEROWNIKA: nie da sie polaczyc ze skanerem: '+$_.Exception.Message;` +
		`ConvertTo-Json -InputObject $wynik -Compress;exit 0};`)
	skrypt.WriteString(`$item=$dev.Items.Item(1);`)
	skrypt.WriteString(fmt.Sprintf(`$dpi=%d;`, rozdzielczosc))
	// Nastawy idą pojedynczo i każda w swoim `try`: backend, który jednej z nich
	// nie zna, ma oddać obraz w nastawie własnej, a nie odmówić całego skanu.
	skrypt.WriteString(`foreach($p in 6147,6148){try{$item.Properties.Item($p).Value=$dpi}catch{}};`)
	if intencja != 0 {
		skrypt.WriteString(fmt.Sprintf(`try{$item.Properties.Item(6146).Value=%d}catch{};`, intencja))
	}
	skrypt.WriteString(`$katalog=` + napisPowerShell(katalog) + `;`)
	skrypt.WriteString(`if(-not (Test-Path -LiteralPath $katalog)){` +
		`[void](New-Item -ItemType Directory -Force -Path $katalog)};`)
	skrypt.WriteString(fmt.Sprintf(`for($i=1;$i -le %d;$i++){`, stron))
	skrypt.WriteString(`$plik=Join-Path $katalog ('skan-'+(Get-Date -Format 'yyyyMMdd-HHmmss')+'-'+$i+'.png');`)
	skrypt.WriteString(`try{$obraz=$item.Transfer('{B96B3CAF-0728-11D3-9D7B-0000F81EF32E}');` +
		`$obraz.SaveFile($plik);[void]$wynik.pliki.Add($plik)}catch{` +
		`if($i -eq 1){$wynik.blad='TRANSFER-NIEUDANY: '+$_.Exception.Message};break}};`)
	skrypt.WriteString(`ConvertTo-Json -InputObject $wynik -Compress`)
	return skrypt.String()
}

// napisPowerShell zamyka wartość w napisie dosłownym PowerShella. Apostrof
// podwaja się — to jedyne znakowanie, jakie napis dosłowny zna, i dlatego
// wartość Operatora nie może z niego wyjść w polecenie.
func napisPowerShell(wartosc string) string {
	return "'" + strings.ReplaceAll(wartosc, "'", "''") + "'"
}

// odmowaWia przekłada rozpoznanie skryptu na zdanie dla Operatora wraz z drogą
// naprawy. Każdy brak ma własną drogę: warstwy WIA nie naprawi się podłączeniem
// skanera, a braku sterownika — uruchomieniem usługi.
func odmowaWia(rozpoznanie string) error {
	tresc := strings.TrimSpace(rozpoznanie)
	switch {
	case strings.HasPrefix(tresc, "WARSTWA-WIA-NIEDOSTEPNA"):
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: na tej maszynie nie ma warstwy WIA (Windows Image Acquisition), "+
				"którą rdzeń sięga po skaner. Naprawa: uruchomić usługę Windows Image "+
				"Acquisition (stisvc) poleceniem `Start-Service stisvc`. "+
				"Diagnostyka warstwy: "+tresc))
	case strings.HasPrefix(tresc, "BRAK-URZADZENIA"):
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: "+tresc+". Rdzeń pytał warstwę WIA i dostał odpowiedź — to "+
				"nie jest awaria drogi, a brak urządzenia. Naprawa: podłączyć skaner, "+
				"włączyć go i sprawdzić wykaz komendą `studio.ingest.device.list`."))
	case strings.HasPrefix(tresc, "BRAK-STEROWNIKA"):
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: skaner jest w wykazie WIA, ale nie da się z nim połączyć — "+
				"to brak albo usterka sterownika WIA urządzenia. Naprawa: doinstalować "+
				"sterownik producenta ze wsparciem WIA (sam sterownik TWAIN nie "+
				"wystarczy). Diagnostyka warstwy: "+tresc))
	case strings.HasPrefix(tresc, "TRANSFER-NIEUDANY"):
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: skaner przyjął nastawy, ale nie oddał obrazu. Naprawa: "+
				"sprawdzić, czy kartka leży na szybie albo w podajniku i czy pokrywa "+
				"jest zamknięta. Diagnostyka warstwy: "+tresc))
	default:
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: warstwa WIA odmówiła skanowania: "+tresc))
	}
}

// bladWarstwyWia dokłada do odmowy arsenału zdanie o tym, CZEGO ta droga
// wymaga. Sam brak `pwsh` opisuje `zewnetrzne.BrakNarzedzia` poprawnie, ale nie
// mówi, że bez niego nie ma na Windowsie skanera wcale — a to jest wiadomość,
// po której Operator wie, co zainstalować.
func bladWarstwyWia(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: skanowanie na Windowsie idzie warstwą WIA, a rdzeń sięga po "+
				"nią PowerShellem — programu `pwsh` na tej maszynie nie ma. Bez niego "+
				"rdzeń nie ma ŻADNEJ drogi do skanera Windows (SANE na tym systemie nie "+
				"istnieje). Naprawa: zainstalować "+brak.Narzedzie.Pakiet+". Droga, "+
				"która działa bez tego: zeskanuj materiał programem systemu i dołóż "+
				"plik komendą `studio.ingest.queue.add`."))
	}
	return err
}

// ── Wspólne dla obu dróg ────────────────────────────────────────────────────

// wolajUrzadzenie jest jedyną drogą warstwy urządzeń do programu zewnętrznego.
// Osobna od `wolajNarzedzie`, bo granicę czasu dobiera wołający: wykaz ma być
// szybki, a skan wolno może trwać minuty. Odmowy przekłada ta sama funkcja, co
// w arsenale modułu — dwa różne przekłady tych samych braków dałyby dwa różne
// zdania o jednej usterce.
func (a *adapterStudia) wolajUrzadzenie(ctx context.Context, n zewnetrzne.Narzedzie,
	argumenty []string, granica time.Duration) ([]byte, error) {

	if a.uruchamiacz == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: rdzeń nie ma uruchamiacza procesów, więc warstwa urządzeń ("+
				n.Nazwa+") nie ma czym wystartować; naprawa: podpiąć warstwę kanału "+
				"przy składaniu rdzenia"))
	}
	okno, zasady, obszar := a.zasiegStudia()
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar, n,
		argumenty, "", granica)
	if err != nil {
		return wynik.Wyjscie, bladArsenaluStudia(err)
	}
	return wynik.Wyjscie, nil
}

// katalogSkanow wskazuje katalog, w którym wolno położyć pobrany obraz. Miejsce
// jest obszarem roboczym okna, bo tam sięga izolacja i tam kolejka wczytywania
// ma prawo czytać. Rdzeń bez ustalonego obszaru schodzi na katalog tymczasowy
// systemu — skan ma powstać, a nie zniknąć na braku nastawy.
func (a *adapterStudia) katalogSkanow() (string, error) {
	korzen := ""
	if a.katalog != nil {
		korzen = strings.TrimSpace(a.katalog.Ustal(konfig.Kontekst{}, "").Sciezka)
	}
	if korzen == "" {
		tymczasowy, err := os.MkdirTemp("", "danaco-skan-")
		if err != nil {
			return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
				"moduł Studio: nie ma gdzie położyć pobranego skanu — rdzeń nie ma "+
					"ustalonego katalogu roboczego, a katalog tymczasowy systemu nie "+
					"powstał: "+err.Error()))
		}
		return tymczasowy, nil
	}
	katalog := filepath.Join(korzen, "studio-skany")
	if err := os.MkdirAll(katalog, 0o755); err != nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: nie da się założyć katalogu skanów "+katalog+": "+err.Error()))
	}
	return katalog, nil
}

// zapiszSkan kładzie bajty obrazu na dysku i oddaje ścieżkę. Nazwa niesie czas
// i numer strony, żeby dwa skany z jednego okna nie nadpisały się wzajemnie.
func zapiszSkan(katalog string, numer int, bajty []byte) (string, error) {
	nazwa := fmt.Sprintf("skan-%s-%d.png", time.Now().UTC().Format("20060102-150405.000"), numer)
	sciezka := filepath.Join(katalog, nazwa)
	if err := os.WriteFile(sciezka, bajty, 0o644); err != nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: pobrany skan nie dał się zapisać do "+sciezka+": "+err.Error()))
	}
	return sciezka, nil
}

// skrocDoPodgladu przycina obcą odpowiedź do długości, którą da się przeczytać
// w oknie. Odpowiedź obcej warstwy bywa wielolinijkowa i wielokilobajtowa, a
// odmowa ma nazwać brak, nie zalać Operatora zrzutem.
func skrocDoPodgladu(tresc string) string {
	jeden := strings.Join(strings.Fields(tresc), " ")
	if len(jeden) <= 240 {
		return jeden
	}
	return jeden[:240] + "…"
}
