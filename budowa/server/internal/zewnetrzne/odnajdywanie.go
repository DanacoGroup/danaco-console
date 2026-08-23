package zewnetrzne

// Odpowiedzialność pliku: jedna droga odnajdywania binarium arsenału — najpierw
// w pakiecie produktu, dopiero potem na ścieżce wyszukiwania systemu.
//
// ── Po co pakiet przed ścieżką ────────────────────────────────────────────────
// Produkt jedzie do Operatora jako instalka, a nie jako lista rzeczy do
// doinstalowania. Program dołożony do pakietu (`pomocniki/<program>`) ma być
// znaleziony niezależnie od tego, co Operator ma w systemie — i ma mieć
// pierwszeństwo, bo to jego wersję sprawdzono przed wydaniem. Wersja zastana na
// maszynie bywa starsza, nowsza albo okrojona; zgodność z nią nie jest niczyją
// obietnicą.
//
// Gdy w pakiecie nie ma nic, zostaje ścieżka systemu — tak pracuje maszyna
// deweloperska i tak działa produkt na maszynie, gdzie Operator ma własną
// instalację programu.
//
// ── Skąd wiadomo, gdzie jest pakiet ───────────────────────────────────────────
// Pakiet zna dwa układy — płaski `pomocniki/<program>` oraz własny katalog
// `pomocniki/<program>/<program>` dla programów niosących własne środowisko.
// Powód rozdzielenia stoi przy `miejscaPakietu`, bo tam składa się ścieżki.
//
// Rdzeń i katalog `pomocniki` wychodzą z jednego pakowania i stoją obok siebie,
// więc pierwszym miejscem jest katalog binarium, które właśnie pracuje. Drugim —
// ta sama ścieżka względna liczona od katalogu bieżącego: obsługuje uruchomienie
// z `go run`, gdzie binarium leży w katalogu tymczasowym, a obok niego nie ma
// nic z produktu. Sekwencja jest ta sama, którą od dawna stosuje pomocnik mowy
// (`mowa/pomocnik.go`) — nie druga jej odmiana.

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// katalogPomocnikow to nazwa katalogu, do którego pakowanie wnosi programy
// towarzyszące. Ta sama nazwa stoi w `tauri.conf.json` jako cel zasobów.
const katalogPomocnikow = "pomocniki"

// Odnajdz wskazuje plik wykonywalny narzędzia oraz mówi, czy w ogóle go widać.
//
// Kolejność jest rozstrzygnięciem, nie wygodą: wskazanie bezwzględne Operatora,
// potem pakiet produktu, na końcu ścieżka systemu.
func Odnajdz(n Narzedzie) (string, bool) {
	program := strings.TrimSpace(n.Program)
	if program == "" {
		return "", false
	}

	// Ścieżka wskazana wprost jest wskazaniem Operatora albo montażu i nie
	// podlega szukaniu — sprawdzamy wyłącznie, czy plik nadaje się do
	// uruchomienia.
	if filepath.IsAbs(program) {
		if wykonywalny(program) {
			return program, true
		}
		return program, false
	}

	for _, miejsce := range miejscaPakietu(program) {
		if wykonywalny(miejsce) {
			return miejsce, true
		}
	}

	if zeSciezki, err := exec.LookPath(program); err == nil {
		return zeSciezki, true
	}
	return program, false
}

// miejscaPakietu składa ścieżki, pod którymi program mógłby leżeć w pakiecie.
// Na Windowsie plik wykonywalny nosi rozszerzenie `.exe`, a nazwa programu
// w deklaracji narzędzia go nie niesie — pakowanie nie zmienia nazw, więc
// rozszerzenie dokłada się tutaj.
//
// ── Dlaczego dwa układy, a nie jeden ─────────────────────────────────────────
// Układ płaski (`pomocniki/<program>`) jest podstawowy i wystarcza programom,
// które są jednym plikiem albo niosą obok siebie kilka własnych bibliotek.
// Nie wystarcza jednak programom, które przynoszą CAŁE WŁASNE ŚRODOWISKO:
// PowerShell to kilkaset zestawów .NET wraz z podkatalogiem `Modules`,
// LibreOffice ma własne drzewo `program`, Chromium — swoje pliki wydania,
// a `rembg` cały osadzony Python. Windows szuka bibliotek w katalogu pliku
// wykonywalnego, więc zsypanie ich wszystkich do jednego `pomocniki` znaczy
// kolizję nazw bibliotek: kilka pakietów wnosi plik o tej samej nazwie i innej
// zawartości, a wygrywa ten, który skopiowano później. Taka wygrana jest
// przypadkiem, nie rozstrzygnięciem, i objawia się dopiero u Operatora.
//
// Dlatego dochodzi drugi układ — `pomocniki/<program>/<program>`, czyli własny
// katalog o nazwie programu, z nietkniętym układem wydania wewnątrz. Program
// widzi swoje biblioteki obok siebie, a żaden inny pakiet mu w to nie wchodzi.
//
// Kolejność ma znaczenie: płaski idzie PIERWSZY, żeby dotychczasowe pakiety
// zachowały się dokładnie tak jak przedtem, a katalog własny był drogą dla
// tych, których płasko położyć się nie da.
func miejscaPakietu(program string) []string {
	nazwy := []string{program}
	if runtime.GOOS == "windows" && !strings.EqualFold(filepath.Ext(program), ".exe") {
		nazwy = append(nazwy, program+".exe")
	}

	var korzenie []string
	if biezace, err := os.Executable(); err == nil {
		korzenie = append(korzenie, filepath.Dir(biezace))
	}
	if katalog, err := os.Getwd(); err == nil {
		korzenie = append(korzenie, katalog)
	}

	// Nazwa katalogu własnego jest nazwą programu BEZ rozszerzenia — katalog nie
	// jest plikiem wykonywalnym i nie nosi `.exe`, choć plik w jego środku nosi.
	katalogWlasny := strings.TrimSuffix(program, filepath.Ext(program))
	if !strings.EqualFold(filepath.Ext(program), ".exe") {
		katalogWlasny = program
	}

	var miejsca []string
	for _, korzen := range korzenie {
		for _, nazwa := range nazwy {
			miejsca = append(miejsca, filepath.Join(korzen, katalogPomocnikow, nazwa))
		}
		for _, nazwa := range nazwy {
			miejsca = append(miejsca, filepath.Join(korzen, katalogPomocnikow, katalogWlasny, nazwa))
		}
	}
	return miejsca
}

// wykonywalny odpowiada, czy pod ścieżką stoi plik zwykły, który system zgodzi
// się uruchomić. Katalog o nazwie programu nie jest programem.
func wykonywalny(sciezka string) bool {
	opis, err := os.Stat(sciezka)
	if err != nil || opis.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return opis.Mode().Perm()&0o111 != 0
}
