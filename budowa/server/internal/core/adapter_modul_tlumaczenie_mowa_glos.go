// Plik obsługuje moduł Translate: dobór syntezatora i odnalezienie głosu dla panelu. Uruchomienie, plik wyjściowy i przekład odmów należą do adapter_modul_tlumaczenie_mowa_silnik.go.
package core

import (
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/zewnetrzne"
)

// dobierzSyntezator rozstrzyga, którym silnikiem czytać: piper, gdy się da, espeak, gdy trzeba, odmowa dwuczłonowa, gdy nie da się żadnym. Piper ma dwa warunki: binarium stojące bez głosu dla języka panelu nie jest silnikiem gotowym do pracy.
func dobierzSyntezator(jezyk string) (wyborSyntezatora, error) {
	piper := zewnetrzne.Narzedzie{
		Nazwa:   "piper (synteza mowy)",
		Program: programPipera(),
		Pakiet:  "piper-tts (arsenał: /opt/danaco-arsenal/piper)",
	}
	powodPipera := ""
	switch {
	case !zewnetrzne.Stoi(piper):
		powodPipera = "nie ma binarium " + piper.Program +
			" (wskazanie: zmienna " + zmiennaPipera + " albo nazwa piper w PATH, albo " + piperArsenalu + ")"
	default:
		glos, znaleziony := glosPipera(jezyk)
		if znaleziony {
			return wyborSyntezatora{silnik: silnikPiper, narzedzie: piper, glos: glos}, nil
		}
		// Brak głosu przy stojącym binarium ma osobny powód i osobną naprawę.
		powodPipera = "binarium " + piper.Program + " stoi, ale w katalogu głosów " + katalogGlosow() +
			" nie ma głosu dla języka " + jezyk +
			" (naprawa: dołożyć plik <" + skrotJezyka(jezyk) + ">*.onnx albo wskazać katalog zmienną " +
			zmiennaGlosowPiper + ")"
	}

	espeak := zewnetrzne.Narzedzie{
		Nazwa:   "espeak-ng (synteza zapasowa)",
		Program: programEspeaka(),
		Pakiet:  "espeak-ng",
	}
	if zewnetrzne.Stoi(espeak) {
		return wyborSyntezatora{silnik: silnikEspeak, narzedzie: espeak, glos: jezyk}, nil
	}

	// Odmowa wymienia obie przyczyny; sama druga wskazywałaby espeaka, choć głos pipera bywa bliższy.
	return wyborSyntezatora{}, bladBrakuSyntezatora(
		"nie ma czym zsyntezować mowy — żaden z dwóch silników nie jest gotowy. " +
			"Piper (głos dobry): " + powodPipera + ". " +
			"Espeak-ng (głos zapasowy): nie ma binarium " + espeak.Program +
			" (wskazanie: zmienna " + zmiennaEspeaka + " albo nazwa espeak-ng w PATH; naprawa: apt install espeak-ng)")
}

// programPipera wskazuje binarium pipera: ścieżka ze zmiennej środowiska, potem goła nazwa w PATH, na końcu miejsce typowe arsenału. Ścieżkę ze zmiennej bierze się wprost i bez sprawdzania na dysku — o wykonywalności rozstrzyga `Stoi`.
func programPipera() string {
	if wskazane := strings.TrimSpace(os.Getenv(zmiennaPipera)); wskazane != "" {
		return wskazane
	}
	if zewnetrzne.Stoi(zewnetrzne.Narzedzie{Program: silnikPiper}) {
		return silnikPiper
	}
	return piperArsenalu
}

// programEspeaka wskazuje binarium syntezy zapasowej — ta sama trójstopniowa
// zasada co wyżej, bez członu arsenału: espeak-ng instaluje się pakietem
// systemowym i stoi w PATH albo nie ma go wcale.
func programEspeaka() string {
	if wskazane := strings.TrimSpace(os.Getenv(zmiennaEspeaka)); wskazane != "" {
		return wskazane
	}
	return silnikEspeak
}

// katalogGlosow wskazuje katalog z plikami `.onnx` — ze zmiennej środowiska albo katalog arsenału syntezy mowy.
func katalogGlosow() string {
	if wskazany := strings.TrimSpace(os.Getenv(zmiennaGlosowPiper)); wskazany != "" {
		return wskazany
	}
	return glosyPiperaArsenalu
}

// glosPipera odnajduje plik głosu dla języka panelu po przedrostku nazwy pliku według zwyczaju nazewniczego pipera. Język skraca się do członu przed separatorem. Katalog nieczytelny nie jest odmową, tylko brakiem głosu.
func glosPipera(jezyk string) (string, bool) {
	katalog := katalogGlosow()
	skrot := skrotJezyka(jezyk)
	if skrot == "" {
		return "", false
	}
	wpisy, err := os.ReadDir(katalog)
	if err != nil {
		return "", false
	}
	for _, wpis := range wpisy {
		if wpis.IsDir() {
			continue
		}
		nazwa := wpis.Name()
		if !strings.HasSuffix(strings.ToLower(nazwa), ".onnx") {
			continue
		}
		male := strings.ToLower(nazwa)
		if !strings.HasPrefix(male, skrot+"_") && !strings.HasPrefix(male, skrot+"-") {
			continue
		}
		return filepath.Join(katalog, nazwa), true
	}
	return "", false
}

// skrotJezyka sprowadza język panelu do członu, którym zaczynają się nazwy
// głosów. Nazwy pełne („polski", „polish") skrótu nie dają — zostają sobą
// i nie dopasują żadnego pliku; tabelki nazw języków tu nie ma, więc „polski"
// nie zamienia się w `pl`.
func skrotJezyka(jezyk string) string {
	male := strings.ToLower(strings.TrimSpace(jezyk))
	if czlon := strings.FieldsFunc(male, func(r rune) bool { return r == '_' || r == '-' }); len(czlon) > 0 {
		return czlon[0]
	}
	return male
}
