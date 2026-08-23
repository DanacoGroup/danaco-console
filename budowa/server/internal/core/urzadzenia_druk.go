// Odpowiedzialność pliku: warstwa druku lokalnego rozdzielona po systemie —
// wykaz drukarek systemowych i WYSŁANIE gotowego wydania na jedną z nich.
//
// ── Czym to się różni od `design.print.export` ──────────────────────────────
// `design.print.export` kończy pracę PLIKIEM: wydaje PDF, TIFF albo EPS z
// przestrzenią barw, spadami i znacznikami cięcia, i kładzie go w magazynie jako
// zasób. To jest wydanie do drukarni. Tego pliku nikt jednak nie wydrukował na
// drukarce stojącej obok Operatora — i to jest dziura, którą zamyka ten plik:
// zasób wydany przez `design.print.export` (albo dowolny inny plik widziany
// przez rdzeń) idzie tu na kolejkę druku systemu.
//
// ── KTÓREJ KOMENDY KONTRAKTU BRAKUJE ────────────────────────────────────────
// Kontrakt (`shared/contract.json`) nie ma komendy wysłania na drukarkę ani
// komendy wykazu drukarek. Rodzina `design.print.*` ma nastawy profilu, kontrolę
// przeddrukową, wydanie i podział wielkoformatowy — i na tym się kończy.
// Nazwy nie wymyślam: warstwa stoi tu gotowa i czeka na dwie komendy, które
// kontrakt musi wnieść (wykaz drukarek oraz zlecenie druku wraz z jego stanem).
// Do tego czasu funkcje tego pliku są wystawione poza pakiet, żeby adapter
// modułu Design mógł je wziąć jedną linią w dniu, w którym komendy powstaną —
// bez przepisywania warstwy.
//
// ── Rozdzielenie po systemie: `runtime.GOOS`, nie warunek budowy ────────────
// Powód ten sam, co przy skanerze (`urzadzenia_skaner.go`): `go build ./...` na
// Linuksie nie skompilowałby gałęzi Windows ani razu, więc zepsułaby się
// niezauważona aż do wydania instalki natywnej.
//
// Na Linuksie drogą jest CUPS (`lpstat`, `lp`) — ta sama, którą druku używa cały
// system. Na Windowsie drogą jest PowerShell (`pwsh`): `Get-Printer` oddaje
// wykaz, a wysłanie idzie przez .NET albo przez czasownik `PrintTo` powłoki
// systemu. Program zewnętrzny idzie WYŁĄCZNIE przez `zewnetrzne.Wolaj`.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// granicaDrukuLokalnego jest granicą czasu jednego zlecenia. Wysłanie na kolejkę
// jest szybkie — długo trwa sam druk, a ten dzieje się już po stronie systemu i
// rdzeń na niego nie czeka.
const granicaDrukuLokalnego = 90 * time.Second

// granicaWykazuDrukarek jest granicą czasu odpytania systemu o drukarki.
const granicaWykazuDrukarek = 30 * time.Second

// Narzędzia warstwy druku. `lp` i `lpstat` NIE stoją jeszcze w wykazie
// zależności rdzenia (`zaleznosci_zewnetrzne.go`) — wpis do wykazu należy do
// tego pliku dopiero wtedy, gdy warstwę zawoła komenda kontraktu, bo sonda
// startowa ma mówić o brakach czynności, które Operator może wykonać.
var (
	narzedzieDrukuCups = zewnetrzne.Narzedzie{
		Nazwa: "CUPS (lp)", Program: "lp", Pakiet: "cups-client",
	}
	narzedzieWykazuDrukarekCups = zewnetrzne.Narzedzie{
		Nazwa: "CUPS (lpstat)", Program: "lpstat", Pakiet: "cups-client",
	}
	narzedzieDrukuWindows = zewnetrzne.Narzedzie{
		Nazwa:   "PowerShell 7 (droga do drukarki Windows)",
		Program: "pwsh",
		Pakiet:  "PowerShell 7 (winget install Microsoft.PowerShell)",
	}
)

// DrukarkaSystemowa opisuje jedną drukarkę widzianą przez system.
type DrukarkaSystemowa struct {
	// Nazwa jest nazwą kolejki systemu i jednocześnie wskazaniem przy zleceniu.
	Nazwa string
	// Opis jest tym, co system mówi o urządzeniu; bywa pusty.
	Opis string
	// Domyslna mówi, czy system drukuje tam bez wskazania.
	Domyslna bool
	// Gotowa mówi, czy kolejka przyjmuje zlecenia. Wstrzymana kolejka przyjmie
	// plik i nie wydrukuje go — Operator ma to wiedzieć przed wysłaniem.
	Gotowa bool
}

// ZlecenieDruku niesie jedno wysłanie na drukarkę.
type ZlecenieDruku struct {
	// Plik jest ścieżką widzianą przez rdzeń — na przykład bajtami zasobu
	// wydanego przez `design.print.export`.
	Plik string
	// Drukarka jest nazwą z wykazu. Puste bierze drukarkę domyślną systemu.
	Drukarka string
	// Kopie mniejsze od 1 znaczą jedną kopię.
	Kopie int
	// Dwustronnie zamawia druk po obu stronach kartki, gdy urządzenie to umie.
	Dwustronnie bool
	// Tytul jest nazwą zlecenia w kolejce systemu. Puste bierze nazwę pliku.
	Tytul string
}

// WarstwaDruku jest wejściem do drukarki systemowej. Struktura, a nie zbiór
// funkcji z siedmioma parametrami: uruchamiacz, zasady izolacji i obszar roboczy
// idą razem w każdym wywołaniu, więc rozdzielanie ich przy każdym wołaniu
// zaprasza do pominięcia jednego z nich.
type WarstwaDruku struct {
	uruchamiacz session.Uruchamiacz
	okno        session.Okno
	zasady      session.Zasady
	obszar      session.Obszar
}

// NowaWarstwaDruku składa warstwę z kompletu, który rdzeń już ma w adapterze.
// Wystawiona poza pakiet po to, żeby adapter modułu Design wziął ją jedną linią,
// gdy kontrakt wniesie komendy druku (patrz nagłówek pliku).
func NowaWarstwaDruku(uruchamiacz session.Uruchamiacz, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar) *WarstwaDruku {

	return &WarstwaDruku{uruchamiacz: uruchamiacz, okno: okno, zasady: zasady, obszar: obszar}
}

// Drukarki oddaje wykaz drukarek systemu drogą właściwą dla systemu.
//
// Pusty wykaz oddaje TYLKO wtedy, gdy system odpowiedział i nie zgłosił żadnej
// kolejki. Brak drogi — brak CUPS, brak `pwsh`, system nieznany — jest odmową
// nazywającą brak, nie pustką udającą „nie ma drukarek".
func (w *WarstwaDruku) Drukarki(ctx context.Context) ([]DrukarkaSystemowa, error) {
	switch runtime.GOOS {
	case systemLinux:
		wyjscie, err := w.wolaj(ctx, narzedzieWykazuDrukarekCups,
			[]string{"-p", "-d"}, granicaWykazuDrukarek)
		if err != nil {
			return nil, bladWarstwyDruku(err)
		}
		return odczytajDrukarkiCups(string(wyjscie)), nil
	case systemWindows:
		wyjscie, err := w.wolaj(ctx, narzedzieDrukuWindows,
			skryptPowerShell(skryptWykazuDrukarek), granicaWykazuDrukarek)
		if err != nil {
			return nil, bladWarstwyDruku(err)
		}
		return odczytajDrukarkiWindows(string(wyjscie))
	default:
		return nil, odmowaDrukuNaTymSystemie()
	}
}

// Wyslij kładzie plik na kolejce druku systemu i oddaje zdanie o tym, co system
// przyjął. Rdzeń nie czeka na wydruk — czeka na przyjęcie zlecenia, bo tylko to
// jest w jego zasięgu.
func (w *WarstwaDruku) Wyslij(ctx context.Context, z ZlecenieDruku) (string, error) {
	plik := strings.TrimSpace(z.Plik)
	if plik == "" {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"druk lokalny: wysłanie na drukarkę bez wskazania pliku"))
	}
	kopie := z.Kopie
	if kopie < 1 {
		kopie = 1
	}
	tytul := strings.TrimSpace(z.Tytul)
	if tytul == "" {
		tytul = filepath.Base(plik)
	}

	switch runtime.GOOS {
	case systemLinux:
		return w.wyslijCups(ctx, plik, tytul, kopie, z)
	case systemWindows:
		return w.wyslijWindows(ctx, plik, tytul, kopie, z)
	default:
		return "", odmowaDrukuNaTymSystemie()
	}
}

// odmowaDrukuNaTymSystemie nazywa brak drogi druku na systemie, którego rdzeń
// nie obsługuje.
func odmowaDrukuNaTymSystemie() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"druk lokalny: rdzeń nie ma warstwy druku na systemie "+runtime.GOOS+
			". Warstwy, które zna, to CUPS (`lp`, `lpstat`) na Linuksie i PowerShell "+
			"(`pwsh`) na Windowsie. Droga, która działa: wydać materiał plikiem przez "+
			"`design.print.export` i wydrukować go programem systemu."))
}

// ── Droga CUPS ──────────────────────────────────────────────────────────────

// wyslijCups zamawia druk przez `lp`. Nastawy idą opcjami CUPS, bo to ta sama
// droga, którą drukuje reszta systemu — własne przygotowanie strumienia byłoby
// drugą implementacją sterownika obok tej, którą maszyna już ma.
func (w *WarstwaDruku) wyslijCups(ctx context.Context, plik, tytul string, kopie int,
	z ZlecenieDruku) (string, error) {

	argumenty := []string{"-n", strconv.Itoa(kopie), "-t", tytul}
	if drukarka := strings.TrimSpace(z.Drukarka); drukarka != "" {
		argumenty = append(argumenty, "-d", drukarka)
	}
	if z.Dwustronnie {
		argumenty = append(argumenty, "-o", "sides=two-sided-long-edge")
	}
	argumenty = append(argumenty, plik)

	wyjscie, err := w.wolaj(ctx, narzedzieDrukuCups, argumenty, granicaDrukuLokalnego)
	if err != nil {
		return "", bladWarstwyDruku(err)
	}
	// `lp` oddaje wiersz postaci: `request id is HP-42 (1 file(s))`. Zdanie idzie
	// do Operatora bez przekładu — to jest identyfikator, którym system druku
	// nazywa jego zlecenie, i po nim je odnajdzie.
	potwierdzenie := strings.TrimSpace(string(wyjscie))
	if potwierdzenie == "" {
		potwierdzenie = "system druku przyjął zlecenie bez identyfikatora"
	}
	return potwierdzenie, nil
}

// odczytajDrukarkiCups czyta wyjście `lpstat -p -d`. Wiersze mają postać:
// `printer HP-42 is idle.  enabled since ...` oraz `system default destination: HP-42`.
func odczytajDrukarkiCups(wyjscie string) []DrukarkaSystemowa {
	drukarki := []DrukarkaSystemowa{}
	domyslna := ""
	for _, wiersz := range strings.Split(wyjscie, "\n") {
		wiersz = strings.TrimSpace(wiersz)
		switch {
		case strings.HasPrefix(wiersz, "printer "):
			reszta := strings.TrimSpace(strings.TrimPrefix(wiersz, "printer "))
			nazwa := reszta
			opis := ""
			if odstep := strings.Index(reszta, " "); odstep > 0 {
				nazwa = reszta[:odstep]
				opis = strings.TrimSpace(strings.TrimSuffix(reszta[odstep+1:], "."))
			}
			drukarki = append(drukarki, DrukarkaSystemowa{
				Nazwa: nazwa,
				Opis:  opis,
				// Kolejka wstrzymana mówi o sobie „disabled"; każdy inny stan
				// (`idle`, `printing`) przyjmuje zlecenia.
				Gotowa: !strings.Contains(opis, "disabled"),
			})
		case strings.HasPrefix(wiersz, "system default destination:"):
			domyslna = strings.TrimSpace(
				strings.TrimPrefix(wiersz, "system default destination:"))
		}
	}
	for numer := range drukarki {
		if drukarki[numer].Nazwa == domyslna {
			drukarki[numer].Domyslna = true
		}
	}
	return drukarki
}

// ── Droga Windows ───────────────────────────────────────────────────────────

// skryptWykazuDrukarek pyta system o kolejki druku. `Get-Printer` stoi w module
// PrintManagement systemu Windows, a nie w PowerShellu — brak modułu jest
// rozpoznany osobno, bo naprawia się go inaczej niż brak samego `pwsh`.
const skryptWykazuDrukarek = `$ErrorActionPreference='Stop';` +
	`try{$d=Get-Printer}catch{Write-Output '[]';` +
	`[Console]::Error.WriteLine('BRAK-MODULU-DRUKU: '+$_.Exception.Message);exit 3};` +
	`$domyslna='';try{$domyslna=(Get-CimInstance -ClassName Win32_Printer | ` +
	`Where-Object {$_.Default -eq $true} | Select-Object -First 1).Name}catch{};` +
	`$w=New-Object System.Collections.ArrayList;` +
	`foreach($p in $d){[void]$w.Add([pscustomobject]@{nazwa=$p.Name;` +
	`opis=[string]$p.DriverName;domyslna=($p.Name -eq $domyslna);` +
	// Kolejka w stanie błędu przyjmie plik i go nie wydrukuje — Operator ma to
	// wiedzieć z wykazu, nie z pustej tacy.
	`gotowa=([string]$p.PrinterStatus -ne 'Error')})};` +
	`ConvertTo-Json -InputObject @($w) -Compress`

// odczytajDrukarkiWindows czyta wykaz oddany przez PowerShell jako JSON.
func odczytajDrukarkiWindows(wyjscie string) ([]DrukarkaSystemowa, error) {
	tresc := strings.TrimSpace(wyjscie)
	if tresc == "" {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"druk lokalny: system nie odpowiedział na pytanie o drukarki. To nie znaczy "+
				"„nie ma drukarek” — znaczy, że skrypt PowerShell nie doszedł do końca. "+
				"Naprawa: sprawdzić, czy `pwsh` startuje i czy usługa bufora wydruku "+
				"(spooler) jest uruchomiona."))
	}
	var zapis []struct {
		Nazwa    string `json:"nazwa"`
		Opis     string `json:"opis"`
		Domyslna bool   `json:"domyslna"`
		Gotowa   bool   `json:"gotowa"`
	}
	if err := json.Unmarshal([]byte(tresc), &zapis); err != nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"druk lokalny: odpowiedzi systemu o drukarkach nie da się odczytać ("+
				err.Error()+"). Treść odpowiedzi: "+skrocDoPodgladu(tresc)))
	}
	drukarki := make([]DrukarkaSystemowa, 0, len(zapis))
	for _, wiersz := range zapis {
		nazwa := strings.TrimSpace(wiersz.Nazwa)
		if nazwa == "" {
			continue
		}
		drukarki = append(drukarki, DrukarkaSystemowa{
			Nazwa: nazwa, Opis: strings.TrimSpace(wiersz.Opis),
			Domyslna: wiersz.Domyslna, Gotowa: wiersz.Gotowa,
		})
	}
	return drukarki, nil
}

// wyslijWindows zamawia druk przez PowerShell. Droga zależy od RODZAJU pliku, bo
// Windows nie ma jednej: obraz drukuje .NET (`System.Drawing.Printing`), tekst
// idzie przez `Out-Printer`, a dokument złożony (PDF, PostScript) potrzebuje
// programu, który go rozumie — i dlatego idzie czasownikiem `PrintTo` powłoki
// systemu. Gdy tego czasownika nikt nie zarejestrował, odmowa mówi to wprost
// zamiast milczeć: plik wysłany w nicość wygląda jak wydruk, który się nie
// pojawił.
func (w *WarstwaDruku) wyslijWindows(ctx context.Context, plik, tytul string, kopie int,
	z ZlecenieDruku) (string, error) {

	skrypt := skryptDrukuWindows(plik, tytul, strings.TrimSpace(z.Drukarka), kopie, z.Dwustronnie)
	wyjscie, err := w.wolaj(ctx, narzedzieDrukuWindows, skryptPowerShell(skrypt),
		granicaDrukuLokalnego)
	if err != nil {
		return "", bladWarstwyDruku(err)
	}
	var zapis struct {
		Zlecenie string `json:"zlecenie"`
		Blad     string `json:"blad"`
	}
	tresc := strings.TrimSpace(string(wyjscie))
	if tresc == "" {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"druk lokalny: system nie odpowiedział na zlecenie druku. Naprawa: sprawdzić, "+
				"czy usługa bufora wydruku (spooler) jest uruchomiona."))
	}
	if err := json.Unmarshal([]byte(tresc), &zapis); err != nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"druk lokalny: odpowiedzi systemu na zlecenie druku nie da się odczytać ("+
				err.Error()+"). Treść odpowiedzi: "+skrocDoPodgladu(tresc)))
	}
	if zapis.Blad != "" {
		return "", odmowaDrukuWindows(zapis.Blad, plik)
	}
	if strings.TrimSpace(zapis.Zlecenie) == "" {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"druk lokalny: system nie potwierdził przyjęcia zlecenia i nie nazwał powodu"))
	}
	return zapis.Zlecenie, nil
}

// skryptDrukuWindows składa skrypt wysłania. Nastawy liczbowe wchodzą liczbami,
// wartości tekstowe — napisem dosłownym PowerShella (`napisPowerShell`), bo tylko
// on nie wypuszcza wartości Operatora w polecenie.
func skryptDrukuWindows(plik, tytul, drukarka string, kopie int, dwustronnie bool) string {
	var s strings.Builder
	s.WriteString(`$ErrorActionPreference='Stop';`)
	s.WriteString(`$wynik=[pscustomobject]@{zlecenie='';blad=''};`)
	s.WriteString(`$plik=` + napisPowerShell(plik) + `;`)
	s.WriteString(`$tytul=` + napisPowerShell(tytul) + `;`)
	s.WriteString(`$drukarka=` + napisPowerShell(drukarka) + `;`)
	s.WriteString(fmt.Sprintf(`$kopie=%d;`, kopie))
	s.WriteString(`if(-not (Test-Path -LiteralPath $plik)){` +
		`$wynik.blad='BRAK-PLIKU: '+$plik;ConvertTo-Json -InputObject $wynik -Compress;exit 0};`)
	// Drukarka wskazana musi istnieć. Zejście na domyślną przy literówce
	// wydrukowałoby materiał na innym urządzeniu — cicho i nieodwracalnie.
	s.WriteString(`if($drukarka -ne ''){try{$null=Get-Printer -Name $drukarka}catch{` +
		`$wynik.blad='BRAK-DRUKARKI: '+$drukarka;` +
		`ConvertTo-Json -InputObject $wynik -Compress;exit 0}};`)
	s.WriteString(`$rozsz=[System.IO.Path]::GetExtension($plik).ToLower();`)
	s.WriteString(`try{if($rozsz -in '.txt','.log','.csv','.md'){` +
		`if($drukarka -ne ''){Get-Content -LiteralPath $plik | Out-Printer -Name $drukarka}` +
		`else{Get-Content -LiteralPath $plik | Out-Printer};` +
		`$wynik.zlecenie='zlecenie tekstowe przyjęte przez bufor wydruku: '+$tytul}` +
		`elseif($rozsz -in '.png','.jpg','.jpeg','.bmp','.gif','.tif','.tiff'){` +
		`Add-Type -AssemblyName System.Drawing;` +
		`for($k=1;$k -le $kopie;$k++){` +
		`$obraz=[System.Drawing.Image]::FromFile($plik);` +
		`$dok=New-Object System.Drawing.Printing.PrintDocument;` +
		`if($drukarka -ne ''){$dok.PrinterSettings.PrinterName=$drukarka};`)
	if dwustronnie {
		s.WriteString(`try{$dok.PrinterSettings.Duplex=` +
			`[System.Drawing.Printing.Duplex]::Vertical}catch{};`)
	}
	s.WriteString(`$dok.DocumentName=$tytul;` +
		`$dok.add_PrintPage({param($nadawca,$zdarzenie);` +
		`$zdarzenie.Graphics.DrawImage($obraz,$zdarzenie.MarginBounds)});` +
		`$dok.Print();$dok.Dispose();$obraz.Dispose()};` +
		`$wynik.zlecenie='zlecenie obrazu przyjęte przez bufor wydruku: '+$tytul}` +
		`else{` +
		// Dokument złożony: czasownik `PrintTo` z drukarką wskazaną, a bez
		// wskazania — `Print` na urządzeniu domyślnym systemu.
		`$czasownik=$(if($drukarka -ne ''){'PrintTo'}else{'Print'});` +
		`if($drukarka -ne ''){Start-Process -FilePath $plik -Verb $czasownik ` +
		`-ArgumentList $drukarka -PassThru -WindowStyle Hidden | Out-Null}` +
		`else{Start-Process -FilePath $plik -Verb $czasownik ` +
		`-PassThru -WindowStyle Hidden | Out-Null};` +
		`$wynik.zlecenie='zlecenie oddane powłoce systemu czasownikiem '+$czasownik+': '+$tytul}}` +
		`catch{$wynik.blad='DRUK-NIEUDANY: '+$_.Exception.Message};`)
	s.WriteString(`ConvertTo-Json -InputObject $wynik -Compress`)
	return s.String()
}

// odmowaDrukuWindows przekłada rozpoznanie skryptu na zdanie wraz z drogą
// naprawy. Każdy brak ma własną drogę — brakującej drukarki nie naprawi się
// instalacją czytnika PDF.
func odmowaDrukuWindows(rozpoznanie, plik string) error {
	tresc := strings.TrimSpace(rozpoznanie)
	switch {
	case strings.HasPrefix(tresc, "BRAK-MODULU-DRUKU"):
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"druk lokalny: na tej maszynie nie ma modułu zarządzania drukiem "+
				"(`Get-Printer` z PrintManagement), którym rdzeń pyta system o kolejki. "+
				"Naprawa: włączyć składnik Windows „Zarządzanie drukowaniem”. "+
				"Diagnostyka: "+tresc))
	case strings.HasPrefix(tresc, "BRAK-PLIKU"):
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"druk lokalny: pliku "+plik+" nie ma na dysku maszyny rdzenia — nie ma czego "+
				"wysłać na drukarkę. Naprawa: wydać materiał komendą `design.print.export` "+
				"i podać ścieżkę bajtów wydanego zasobu."))
	case strings.HasPrefix(tresc, "BRAK-DRUKARKI"):
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"druk lokalny: "+tresc+". Rdzeń NIE zeszedł na drukarkę domyślną — wydruk na "+
				"innym urządzeniu niż wskazane jest szkodą nieodwracalną. Naprawa: "+
				"sprawdzić nazwę w wykazie drukarek systemu."))
	default:
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"druk lokalny: system nie przyjął zlecenia druku: "+tresc+
				". Częsty powód przy dokumencie PDF: żaden program nie ma zarejestrowanego "+
				"czasownika drukowania, więc powłoka systemu nie wie, czym go wydrukować; "+
				"naprawa: zainstalować czytnik PDF obsługujący druk z powłoki."))
	}
}

// ── Wspólne ─────────────────────────────────────────────────────────────────

// wolaj jest jedyną drogą warstwy druku do programu zewnętrznego.
func (w *WarstwaDruku) wolaj(ctx context.Context, n zewnetrzne.Narzedzie,
	argumenty []string, granica time.Duration) ([]byte, error) {

	if w == nil || w.uruchamiacz == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"druk lokalny: rdzeń nie ma uruchamiacza procesów, więc warstwa druku ("+
				n.Nazwa+") nie ma czym wystartować; naprawa: podpiąć warstwę kanału "+
				"przy składaniu rdzenia"))
	}
	wynik, err := zewnetrzne.Wolaj(ctx, w.uruchamiacz, w.okno, w.zasady, w.obszar, n,
		argumenty, "", granica)
	if err != nil {
		return wynik.Wyjscie, err
	}
	return wynik.Wyjscie, nil
}

// bladWarstwyDruku przekłada brak programu na odmowę mówiącą, czego brakuje i co
// bez tego nie działa. Rozstrzygnięcie to samo, co w arsenale modułu Studio:
// brak binarium jest zapleczem niedostępnym i kodem PONAWIALNYM, bo po
// instalacji to samo żądanie przejdzie.
func bladWarstwyDruku(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		zdanie := "druk lokalny: nie ma na tej maszynie programu " + brak.Narzedzie.Nazwa +
			" (" + brak.Narzedzie.Program + "), a bez niego rdzeń nie ma drogi do drukarki " +
			"systemowej — plik wydany do druku pozostanie plikiem"
		if brak.Narzedzie.Pakiet != "" {
			zdanie += "; naprawa: zainstalować " + brak.Narzedzie.Pakiet
		}
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable, zdanie))
	}
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"druk lokalny: "+err.Error()))
}
