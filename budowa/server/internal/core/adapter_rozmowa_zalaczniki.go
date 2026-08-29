// Plik doprowadza załączniki wiadomości do modelu: URI danych ląduje w magazynie treści, ścieżka istniejącego pliku idzie dalej bez kopiowania, a pozostałe odwołania wracają nazwane jako niedoręczone.
package core

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// podkatalogZalacznikowRozmowy oddziela bloby rozmowy od magazynu modułu
	// Library. Ten sam katalog danych i ten sam sposób zapisu, ale własny
	// korzeń: załącznik rozmowy nie jest pozycją biblioteki i nie pojawia się
	// w jej magazynie jako treść bez wiersza.
	podkatalogZalacznikowRozmowy = "zalaczniki"
	podkatalogRozmowy            = "rozmowa"
	// przedrostekURIDanych rozpoznaje odwołanie niosące bajty ze sobą wprost, bez dotykania dysku rdzenia.
	przedrostekURIDanych = "data:"
	// znacznikBase64URI oddziela nagłówek URI danych od ładunku base64 przy rozbiorze odwołania w tym pliku.
	znacznikBase64URI = ";base64,"
)

// rozszerzeniaZalacznikow przekłada typ nośny URI danych na rozszerzenie pliku; narzędzia programu claude rozpoznają rodzaj pliku po nazwie, więc obraz bez rozszerzenia zostaje odczytany jako bajty tekstu.
var rozszerzeniaZalacznikow = map[string]string{
	"image/png":        ".png",
	"image/jpeg":       ".jpg",
	"image/webp":       ".webp",
	"image/gif":        ".gif",
	"image/svg+xml":    ".svg",
	"application/pdf":  ".pdf",
	"text/plain":       ".txt",
	"text/markdown":    ".md",
	"application/json": ".json",
}

// magazynZalacznikow odkłada bajty załącznika pod katalogiem danych rdzenia, budując się na tym samym typie co magazyn treści biblioteki, ale nad własnym korzeniem.
func magazynZalacznikow(katalogDanych string) *magazynTresciBiblioteki {
	if strings.TrimSpace(katalogDanych) == "" {
		return nil
	}
	return &magazynTresciBiblioteki{
		katalog: filepath.Join(katalogDanych, podkatalogRozmowy, podkatalogZalacznikowRozmowy),
	}
}

// zalacznikTury jest jednym odwołaniem po rozstrzygnięciu.
//
// Pole `Sciezka` puste znaczy odwołanie nierozwiązane — wtedy `Powod` podaje
// przyczynę, a treść wiadomości niesie to zdanie do modelu i do operatora.
type zalacznikTury struct {
	Odwolanie string
	Sciezka   string
	Powod     string
}

// rozwiazZalaczniki przekłada odwołania żądania na ścieżki na nośniku, sprawdzając kolejno od kształtu najpewniejszego: URI danych po przedrostku, ścieżkę dopiero po potwierdzeniu, że plik istnieje.
func rozwiazZalaczniki(magazyn *magazynTresciBiblioteki, odwolania []string) []zalacznikTury {
	if len(odwolania) == 0 {
		return nil
	}
	wynik := make([]zalacznikTury, 0, len(odwolania))
	for _, odwolanie := range odwolania {
		przyciete := strings.TrimSpace(odwolanie)
		if przyciete == "" {
			continue
		}
		wynik = append(wynik, rozwiazZalacznik(magazyn, przyciete))
	}
	return wynik
}

// rozwiazZalacznik rozstrzyga pojedyncze odwołanie, sprawdzając kolejno kształt URI danych i ścieżkę pliku na nośniku.
func rozwiazZalacznik(magazyn *magazynTresciBiblioteki, odwolanie string) zalacznikTury {
	if strings.HasPrefix(odwolanie, przedrostekURIDanych) {
		return zURIDanych(magazyn, odwolanie)
	}
	if filepath.IsAbs(odwolanie) {
		if stan, err := os.Stat(odwolanie); err == nil && !stan.IsDir() {
			return zalacznikTury{Odwolanie: odwolanie, Sciezka: odwolanie}
		}
		return zalacznikTury{
			Odwolanie: odwolanie,
			Powod:     "ścieżka wskazana załącznikiem nie prowadzi do pliku na nośniku serwera",
		}
	}
	// Rdzeń nie zna magazynu tłumaczącego identyfikator na treść, więc odwołanie wraca jako niedoręczone.
	return zalacznikTury{
		Odwolanie: odwolanie,
		Powod: "odwołanie nie jest ani URI danych, ani ścieżką bezwzględną — " +
			"serwer nie ma magazynu, który rozwiązałby je do treści",
	}
}

// zURIDanych materializuje bajty niesione w odwołaniu i oddaje ścieżkę bloba zapisanego w magazynie załączników.
func zURIDanych(magazyn *magazynTresciBiblioteki, odwolanie string) zalacznikTury {
	skrot := skrocOdwolanie(odwolanie)
	miejsce := strings.Index(odwolanie, znacznikBase64URI)
	if miejsce < 0 {
		return zalacznikTury{Odwolanie: skrot,
			Powod: "URI danych bez ładunku base64 — serwer nie odkłada treści, której w odwołaniu nie ma"}
	}
	typNosny := odwolanie[len(przedrostekURIDanych):miejsce]
	bajty, err := base64.StdEncoding.DecodeString(odwolanie[miejsce+len(znacznikBase64URI):])
	if err != nil {
		return zalacznikTury{Odwolanie: skrot, Powod: "ładunek base64 załącznika nie daje się zdekodować: " + err.Error()}
	}
	if len(bajty) == 0 {
		return zalacznikTury{Odwolanie: skrot, Powod: "załącznik zerowej długości — nie ma czego pokazać modelowi"}
	}
	if magazyn == nil {
		// Magazyn wpina montaż; jego brak nie przerywa tury, tylko zostaje nazwany w treści odpowiedzi.
		return zalacznikTury{Odwolanie: skrot,
			Powod: "magazyn załączników rozmowy nie jest wpięty — bajty nie mają gdzie się położyć"}
	}
	suma := sha256.Sum256(bajty)
	nazwa := hex.EncodeToString(suma[:]) + rozszerzenieZalacznika(typNosny)
	sciezka, err := magazyn.Zapisz(bajty, nazwa)
	if err != nil {
		return zalacznikTury{Odwolanie: skrot, Powod: "zapis załącznika w magazynie serwera nie powiódł się: " + err.Error()}
	}
	return zalacznikTury{Odwolanie: skrot, Sciezka: sciezka}
}

// rozszerzenieZalacznika dobiera rozszerzenie do typu nośnego. Typ nieznany
// zostaje bez rozszerzenia — zmyślone („.bin", „.png") mówiłoby o rodzaju
// pliku rzecz, której odwołanie nie stwierdza.
func rozszerzenieZalacznika(typNosny string) string {
	czysty := strings.TrimSpace(strings.SplitN(typNosny, ";", 2)[0])
	return rozszerzeniaZalacznikow[strings.ToLower(czysty)]
}

// skrocOdwolanie przycina URI danych do nagłówka. Pełny ładunek base64 nie ma
// czego szukać ani w dzienniku, ani w treści wiadomości.
func skrocOdwolanie(odwolanie string) string {
	if miejsce := strings.Index(odwolanie, znacznikBase64URI); miejsce >= 0 {
		return odwolanie[:miejsce+len(znacznikBase64URI)] + "…"
	}
	if len(odwolanie) > 64 {
		return odwolanie[:64] + "…"
	}
	return odwolanie
}

// wplecZalaczniki dokłada do treści pytania nazwany blok odwołań, bo kanał główny przyjmuje od rdzenia dokładnie jedno wejście rozmowy; blok jest oddzielony i podpisany, żeby model odróżnił go od słów operatora.
func wplecZalaczniki(tresc string, zalaczniki []zalacznikTury) string {
	if len(zalaczniki) == 0 {
		return tresc
	}
	var blok strings.Builder
	blok.WriteString("[Załączniki wiadomości Operatora]\n")
	blok.WriteString("Poniższe pliki leżą na nośniku tego serwera. Sięgnij po nie narzędziem odczytu pliku, jeżeli są potrzebne do odpowiedzi.\n")
	for numer, z := range zalaczniki {
		if z.Sciezka != "" {
			fmt.Fprintf(&blok, "%d. %s\n", numer+1, z.Sciezka)
			continue
		}
		fmt.Fprintf(&blok, "%d. NIEDORĘCZONY (%s) — %s\n", numer+1, z.Odwolanie, z.Powod)
	}
	if strings.TrimSpace(tresc) == "" {
		return strings.TrimRight(blok.String(), "\n")
	}
	return tresc + "\n\n" + strings.TrimRight(blok.String(), "\n")
}
