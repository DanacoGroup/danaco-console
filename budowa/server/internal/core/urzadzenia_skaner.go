// Warstwa urządzeń wejściowych skanera rozdzielona po systemie operacyjnym pod
// runtime.GOOS w jednym pliku: wykaz urządzeń i pobranie obrazu drogą SANE na
// Linuksie i WIA przez PowerShell na Windowsie.
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

// granicaSkanowaniaUrzadzenia jest granicą czasu jednego przebiegu skanera,
// osobną od granicy arsenału, bo skan strony A4 na wolnym urządzeniu USB trwa
// dłużej niż rozpakowanie archiwum.
const granicaSkanowaniaUrzadzenia = 5 * time.Minute

// granicaWykazuUrzadzen jest granicą czasu samego wykazu. Krótka celowo:
// odpytanie warstwy skanera to czynność natychmiastowa, a gdy nie jest —
// urządzenie zwisło i Operator ma to usłyszeć, zanim odejdzie od okna.
const granicaWykazuUrzadzen = 45 * time.Second

// rozdzielczoscSkanowaniaDomyslna to punkty na cal brane, gdy żądanie nie mówi
// nic. 300 jest progiem czytelności dla rozpoznania pisma — niżej Tesseract
// gubi litery, wyżej rośnie czas skanu bez zysku dla tekstu.
const rozdzielczoscSkanowaniaDomyslna = 300

// narzedzieSkaneraWia jest tym samym programem co narzedziePowerShell (pwsh),
// ale pod nazwą mówiącą o tej czynności, żeby odmowa nazywała brak do
// skanowania, nie do analizy skryptów.
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
			return nil, bladWarstwySane(err)
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
		"moduł Studio: serwer nie ma warstwy skanera na systemie "+runtime.GOOS+
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
				"urządzeń ("+err.Error()+"). Serwer NIE oddaje w tym miejscu pustego "+
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
	// TrybBarwny nazywa tryb koloru słowem produktowym; puste zostawia nastawę
	// urządzenia nietkniętą.
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

// skanujSane pobiera obrazy przez SANE: obraz przychodzi wyjściem programu,
// nie plikiem, bo scanimage pisze PNG na wyjście standardowe, a rdzeń zapisuje
// bajty sam.
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
				// Podajnik pustego arkusza nie poda i SANE odmawia — strony pobrane
				// wcześniej są prawdziwe.
				return sciezki, nil
			}
			return nil, a.odmowaSkanuSane(ctx, err)
		}
		if len(wyjscie) == 0 {
			return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
				"moduł Studio: SANE zakończył skanowanie bez ani jednego bajtu obrazu. "+
					"Naprawa: sprawdzić `scanimage -L`, czy urządzenie odpowiada, i czy "+
					"proces serwera ma prawo do jego węzła (grupa `scanner`)."))
		}
		sciezka, err := zapiszSkan(katalog, numer, wyjscie)
		if err != nil {
			return nil, err
		}
		sciezki = append(sciezki, sciezka)
	}
	return sciezki, nil
}

// odmowaSkanuSane rozstrzyga, czy skan nie doszedł do skutku z braku urządzenia
// czy z innej przyczyny, pytaniem o wykaz, nie czytaniem diagnostyki programu —
// scanimage kończy się tym samym kodem przy każdej przyczynie.
func (a *adapterStudia) odmowaSkanuSane(ctx context.Context, pierwotna error) error {
	// Brak samego programu rozstrzyga się bez pytania o wykaz: wykaz idzie tym
	// samym programem.
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(pierwotna, &brak) {
		return bladWarstwySane(pierwotna)
	}
	urzadzenia, blad := a.wykazSkanerow(ctx)
	if blad != nil || len(urzadzenia) > 0 {
		return pierwotna
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Studio: warstwa SANE (`scanimage`) odpowiedziała i nie widzi ani jednego "+
			"skanera. To nie jest awaria drogi ani usterka serwera, lecz brak urządzenia "+
			"po stronie maszyny Operatora. Naprawa: podłączyć skaner, włączyć go i sprawdzić "+
			"wykaz komendą `studio.ingest.device.list`; jeżeli urządzenie tam jest, procesowi "+
			"serwera brakuje prawa do jego węzła (grupa `scanner`). Droga, która działa bez "+
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

// skanujWia pobiera obrazy przez WIA: obraz nie wraca wyjściem procesu, bo
// PowerShell psułby bajty znakowaniem, więc skrypt zapisuje pliki metodą
// SaveFile i oddaje ich ścieżki jako JSON.
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
				"obrazu, nie nazywając powodu. Serwer nie oddaje tu pustej kolejki, bo "+
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

// skryptPowerShell składa wywołanie pwsh z jednym skryptem, tymi samymi
// przełącznikami co moduł Terminal: bez logo, bez profilu i bez interakcji.
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

// skryptSkanuWia składa skrypt pobrania stron. Nastawy wchodzą liczbami
// i ścieżką katalogu, nigdy tekstem od zewnątrz — identyfikator urządzenia
// idzie w apostrofach z podwojonym apostrofem.
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
	// Nastawy idą pojedynczo, każda w swoim try: nastawa nieznana ma oddać
	// obraz, nie odmówić skanu.
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
				"którą serwer sięga po skaner. Naprawa: uruchomić usługę Windows Image "+
				"Acquisition (stisvc) poleceniem `Start-Service stisvc`. "+
				"Diagnostyka warstwy: "+tresc))
	case strings.HasPrefix(tresc, "BRAK-URZADZENIA"):
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: "+tresc+". Serwer pytał warstwę WIA i dostał odpowiedź — to "+
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

// bladWarstwySane dokłada do odmowy arsenału zdanie o tym, czego ta droga
// wymaga, bo sam brak scanimage nie mówi, że bez niego nie ma na Linuksie
// skanera wcale.
func bladWarstwySane(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: skanowanie na Linuksie idzie warstwą SANE, a serwer sięga po "+
				"nią programem `scanimage` — tego programu na tej maszynie nie ma. Bez "+
				"niego serwer nie ma ŻADNEJ drogi do skanera na Linuksie (WIA na tym "+
				"systemie nie istnieje). Naprawa: zainstalować "+brak.Narzedzie.Pakiet+
				". Droga, która działa bez tego: zeskanuj materiał programem systemu "+
				"i dołóż plik komendą `studio.ingest.queue.add`."))
	}
	return err
}

// bladWarstwyWia dokłada do odmowy arsenału zdanie o tym, czego ta droga
// wymaga, bo sam brak pwsh nie mówi, że bez niego nie ma na Windowsie skanera
// wcale.
func bladWarstwyWia(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Studio: skanowanie na Windowsie idzie warstwą WIA, a serwer sięga po "+
				"nią PowerShellem — programu `pwsh` na tej maszynie nie ma. Bez niego "+
				"serwer nie ma ŻADNEJ drogi do skanera Windows (SANE na tym systemie nie "+
				"istnieje). Naprawa: zainstalować "+brak.Narzedzie.Pakiet+". Droga, "+
				"która działa bez tego: zeskanuj materiał programem systemu i dołóż "+
				"plik komendą `studio.ingest.queue.add`."))
	}
	return err
}

// ── Wspólne dla obu dróg ────────────────────────────────────────────────────

// wolajUrzadzenie jest jedyną drogą warstwy urządzeń do programu zewnętrznego,
// osobną od wolajNarzedzie, bo granicę czasu dobiera wołający: wykaz ma być
// szybki, a skan wolno może trwać minuty.
func (a *adapterStudia) wolajUrzadzenie(ctx context.Context, n zewnetrzne.Narzedzie,
	argumenty []string, granica time.Duration) ([]byte, error) {

	if a.uruchamiacz == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Studio: serwer nie ma uruchamiacza procesów, więc warstwa urządzeń ("+
				n.Nazwa+") nie ma czym wystartować; naprawa: podpiąć warstwę kanału "+
				"przy składaniu serwera"))
	}
	okno, zasady, obszar := a.zasiegStudia()
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar, n,
		argumenty, "", granica)
	if err != nil {
		// Brak programu wychodzi stąd surowy — tylko gałąź drogi wołającej zna
		// zdanie o obejściu.
		var brak *zewnetrzne.BrakNarzedzia
		if errors.As(err, &brak) {
			return wynik.Wyjscie, err
		}
		return wynik.Wyjscie, bladArsenaluStudia(err)
	}
	return wynik.Wyjscie, nil
}

// katalogSkanow wskazuje katalog, w którym wolno położyć pobrany obraz —
// obszar roboczy okna, bo tam sięga izolacja i tam kolejka wczytywania ma
// prawo czytać.
func (a *adapterStudia) katalogSkanow() (string, error) {
	korzen := ""
	if a.katalog != nil {
		korzen = strings.TrimSpace(a.katalog.Ustal(konfig.Kontekst{}, "").Sciezka)
	}
	if korzen == "" {
		tymczasowy, err := os.MkdirTemp("", "danaco-skan-")
		if err != nil {
			return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
				"moduł Studio: nie ma gdzie położyć pobranego skanu — serwer nie ma "+
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
