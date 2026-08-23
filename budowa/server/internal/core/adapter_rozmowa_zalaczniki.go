// Doprowadzenie załączników wiadomości do modelu.
//
// Odwołania z `MessageSendRequest.attachments` idą dwiema drogami: do wiersza
// wiadomości (`adapter_rozmowa.go`) oraz — tutaj — do treści zapytania kanału.
// Drugiej drogi nie da się ominąć, bo `zapytanieKanalu` niesie samą treść,
// a kanał CLI nie ma osobnego pola na załącznik.
//
// Kontrakt nazywa elementy listy „odwołaniami do załączników" i nie zawęża ich
// kształtu. Rdzeń rozstrzyga trzy przypadki:
//   - URI danych (`data:image/png;base64,…`) niesie bajty ze sobą: lądują
//     w magazynie treści rdzenia, a do modelu idzie ścieżka bloba. Tak posyła
//     adnotacje moduł Browser (`client/src/moduly/browser/warstwa-adnotacji.ts`);
//   - ścieżka bezwzględna istniejącego pliku idzie dalej bez kopiowania — tego
//     samego rodzaju odwołanie zapisuje moduł Library
//     (`adapter_modul_library_magazyn.go`);
//   - pozostałe odwołania wracają nazwane jako niedoręczone, bez zgadywania.
//
// Ładunek base64 nie wchodzi do treści: model nie zdekoduje kilobajtów napisu,
// a okno kontekstu za nie zapłaci. Wchodzi wyłącznie ścieżka, po którą model
// sięga narzędziem odczytu pliku.
//
// Plik nie zmienia kształtu kontraktu i nie sprząta magazynu: blob adresowany
// sumą kontrolną żyje tak samo długo jak wiersz wiadomości, który go wymienia.
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
	// przedrostekURIDanych rozpoznaje odwołanie niosące bajty ze sobą.
	przedrostekURIDanych = "data:"
	// znacznikBase64URI oddziela nagłówek URI danych od ładunku base64.
	znacznikBase64URI = ";base64,"
)

// rozszerzeniaZalacznikow przekłada typ nośny URI danych na rozszerzenie pliku.
//
// Narzędzia programu `claude` rozpoznają rodzaj pliku po nazwie, więc obraz bez
// rozszerzenia zostaje odczytany jako bajty tekstu. Nazwą pliku pozostaje suma
// kontrolna (jedna treść = jeden blob), a rozszerzenie tylko ją domyka.
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

// magazynZalacznikow odkłada bajty załącznika pod katalogiem danych rdzenia.
// Buduje się na tym samym typie, co magazyn treści biblioteki — zapis atomowy
// pod nazwą będącą sumą kontrolną jest tam już rozstrzygnięty
// (`adapter_modul_library_magazyn.go`) — ale nad własnym korzeniem.
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

// rozwiazZalaczniki przekłada odwołania żądania na ścieżki na nośniku.
//
// Kolejność sprawdzeń idzie od kształtu najpewniejszego: URI danych rozpoznaje
// się po przedrostku bez dotykania dysku, ścieżkę — dopiero po nieudanym
// dopasowaniu URI i wyłącznie po potwierdzeniu, że plik istnieje. Ścieżka
// niesprawdzona byłaby dla modelu obietnicą pliku, którego nie ma.
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

// rozwiazZalacznik rozstrzyga pojedyncze odwołanie.
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
			Powod:     "ścieżka wskazana załącznikiem nie prowadzi do pliku na nośniku rdzenia",
		}
	}
	// Rdzeń nie zna magazynu, który tłumaczyłby nieprzejrzysty identyfikator na
	// treść, więc odwołanie wraca nazwane jako niedoręczone zamiast zostać po
	// cichu pominięte.
	return zalacznikTury{
		Odwolanie: odwolanie,
		Powod: "odwołanie nie jest ani URI danych, ani ścieżką bezwzględną — " +
			"rdzeń nie ma magazynu, który rozwiązałby je do treści",
	}
}

// zURIDanych materializuje bajty niesione w odwołaniu i oddaje ścieżkę bloba.
func zURIDanych(magazyn *magazynTresciBiblioteki, odwolanie string) zalacznikTury {
	skrot := skrocOdwolanie(odwolanie)
	miejsce := strings.Index(odwolanie, znacznikBase64URI)
	if miejsce < 0 {
		return zalacznikTury{Odwolanie: skrot,
			Powod: "URI danych bez ładunku base64 — rdzeń nie odkłada treści, której w odwołaniu nie ma"}
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
		// Magazyn wpina montaż. Jego brak nie przerywa tury, ale zostaje
		// nazwany w treści — inaczej załącznik znikałby bez śladu.
		return zalacznikTury{Odwolanie: skrot,
			Powod: "magazyn załączników rozmowy nie jest wpięty — bajty nie mają gdzie się położyć"}
	}
	suma := sha256.Sum256(bajty)
	nazwa := hex.EncodeToString(suma[:]) + rozszerzenieZalacznika(typNosny)
	sciezka, err := magazyn.Zapisz(bajty, nazwa)
	if err != nil {
		return zalacznikTury{Odwolanie: skrot, Powod: "zapis załącznika w magazynie rdzenia nie powiódł się: " + err.Error()}
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

// wplecZalaczniki dokłada do treści pytania nazwany blok odwołań.
//
// Blok idzie w treści, a nie osobnym polem, bo kanał główny przyjmuje od rdzenia
// dokładnie jedno wejście rozmowy: `injection.Zapytanie.Tekst`, wysyłane na
// stdin jako JSON-lines (`injection/ustawienia.go`, `injection/przebieg.go`).
// Ścieżka wpleciona w treść jest więc jedyną drogą odwołania do modelu, a droga
// ta jest skuteczna: model czyta plik narzędziem odczytu.
//
// Blok jest oddzielony i podpisany, żeby model odróżnił go od słów operatora.
// Odwołania niedoręczone stoją w tym samym bloku, aby model wiedział, że coś
// pokazano i czego nie dostał.
func wplecZalaczniki(tresc string, zalaczniki []zalacznikTury) string {
	if len(zalaczniki) == 0 {
		return tresc
	}
	var blok strings.Builder
	blok.WriteString("[Załączniki wiadomości Operatora]\n")
	blok.WriteString("Poniższe pliki leżą na nośniku tego rdzenia. Sięgnij po nie narzędziem odczytu pliku, jeżeli są potrzebne do odpowiedzi.\n")
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
