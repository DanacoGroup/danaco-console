// Punkt wejścia narzędzia oglądania listów: składa pakietem mail list aktywacji
// danymi przykładowymi i odkłada jego postać graficzną plikiem do obejrzenia
// przeglądarką. Odwołania `cid:` podgląd podmienia na obrazy wpisane w treść.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	netmail "net/mail"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/core/listy"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/mail"
)

func main() {
	sciezka := flag.String("plik", "", "plik, do którego odkłada się podgląd; bez niego treść idzie na wyjście")
	flag.Parse()

	strona, err := podglad()
	if err != nil {
		fmt.Fprintln(os.Stderr, "nie złożono podglądu:", err)
		os.Exit(1)
	}
	if *sciezka == "" {
		fmt.Print(strona)
		return
	}
	if err := os.WriteFile(*sciezka, []byte(strona), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "nie odłożono podglądu:", err)
		os.Exit(1)
	}
}

// danePrzykladowe wypełniają zmienne listu aktywacji wartościami wzorowanymi na
// opracowaniu (rozdz. 6.2). Podgląd nie ma konta ani drogi potwierdzenia, więc
// żadna z tych wartości nie pochodzi z rdzenia.
func danePrzykladowe() map[string]string {
	teraz := time.Now()
	return map[string]string{
		"recipient_address": "a.kowalska@przyklad.pl",
		"support_address":   "support@danaco-group.pl",
		"year":              strconv.Itoa(teraz.Year()),
		"code":              "418 402",
		"expiry_minutes":    "60",
		"requested_at":      teraz.Format("02.01.2006, 15:04 MST"),
		"activation_url":    "https://console.danaco-group.pl/aktywacja",
	}
}

// podglad składa list aktywacji i oddaje jego postać graficzną gotową do
// otwarcia przeglądarką.
func podglad() (string, error) {
	komplet, err := mail.WbudowanyKomplet()
	if err != nil {
		return "", err
	}
	wartosci := danePrzykladowe()
	wiadomosc, err := komplet.Build(mail.KindAccountActivation, mail.Envelope{
		From: netmail.Address{
			Name:    konfiguracja.NadawcaNazwaDomyslna,
			Address: konfiguracja.NadawcaAdresDomyslny,
		},
		To: netmail.Address{Address: wartosci["recipient_address"]},
	}, mail.Content{Values: wartosci}, mail.Logos{
		Light: listy.ZnakJasny, LightName: "danaco.png",
		Dark: listy.ZnakCiemny, DarkName: "danaco-ciemny.png",
	})
	if err != nil {
		return "", err
	}
	return postacGraficzna(wiadomosc.Data)
}

// postacGraficzna wyjmuje ze złożonej wiadomości część `text/html` i podmienia
// w niej odwołania `cid:` na obrazy wpisane w treść.
func postacGraficzna(dane []byte) (string, error) {
	wiadomosc, err := netmail.ReadMessage(bytes.NewReader(dane))
	if err != nil {
		return "", fmt.Errorf("odczyt złożonej wiadomości: %w", err)
	}
	_, parametry, err := mime.ParseMediaType(wiadomosc.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("typ wiadomości: %w", err)
	}

	var tresc string
	var obrazy []wpisanyObraz
	czesci := multipart.NewReader(wiadomosc.Body, parametry["boundary"])
	for {
		czesc, err := czesci.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("część wiadomości: %w", err)
		}
		rodzaj, wewnetrzne, err := mime.ParseMediaType(czesc.Header.Get("Content-Type"))
		if err != nil {
			return "", fmt.Errorf("typ części wiadomości: %w", err)
		}
		if rodzaj == "multipart/alternative" {
			if tresc, err = czescGraficzna(czesc, wewnetrzne["boundary"]); err != nil {
				return "", err
			}
			continue
		}
		obraz, err := wpisany(czesc, rodzaj)
		if err != nil {
			return "", err
		}
		if obraz.odwolanie != "" {
			obrazy = append(obrazy, obraz)
		}
	}
	if tresc == "" {
		return "", fmt.Errorf("złożona wiadomość nie ma części text/html")
	}
	return zObrazamiWTresci(tresc, obrazy), nil
}

// wpisanyObraz wiąże odwołanie `cid:` z obrazem gotowym do wpisania w treść.
type wpisanyObraz struct {
	odwolanie string
	dane      string
}

// wpisany przekłada część obrazu na wpis `data:`. Część niesie base64 łamany
// co 76 znaków, a odwołanie `data:` złamań wiersza nie znosi — stąd zdjęcie
// złamań zamiast odkodowania i zakodowania z powrotem.
func wpisany(czesc *multipart.Part, rodzaj string) (wpisanyObraz, error) {
	odwolanie := strings.Trim(czesc.Header.Get("Content-ID"), "<>")
	if odwolanie == "" {
		return wpisanyObraz{}, nil
	}
	surowe, err := io.ReadAll(czesc)
	if err != nil {
		return wpisanyObraz{}, fmt.Errorf("odczyt obrazu %s: %w", odwolanie, err)
	}
	zwarte := strings.NewReplacer("\r", "", "\n", "").Replace(string(surowe))
	return wpisanyObraz{
		odwolanie: odwolanie,
		dane:      "data:" + rodzaj + ";base64," + zwarte,
	}, nil
}

// zObrazamiWTresci podmienia odwołania na obrazy, od najdłuższego odwołania.
// Kolejność jest wiążąca: identyfikator znaku jasnego jest przedrostkiem
// ciemnego, a podmieniacz dopasowuje wzorce w kolejności argumentów.
func zObrazamiWTresci(tresc string, obrazy []wpisanyObraz) string {
	sort.Slice(obrazy, func(i, j int) bool {
		return len(obrazy[i].odwolanie) > len(obrazy[j].odwolanie)
	})
	pary := make([]string, 0, 2*len(obrazy))
	for _, obraz := range obrazy {
		pary = append(pary, "cid:"+obraz.odwolanie, obraz.dane)
	}
	return strings.NewReplacer(pary...).Replace(tresc)
}

// czescGraficzna wyjmuje postać `text/html` z części `multipart/alternative`.
// Odczyt części zdejmuje kodowanie quoted-printable samoczynnie.
func czescGraficzna(czesc *multipart.Part, granica string) (string, error) {
	postacie := multipart.NewReader(czesc, granica)
	for {
		postac, err := postacie.NextPart()
		if err == io.EOF {
			return "", fmt.Errorf("część multipart/alternative nie niesie postaci text/html")
		}
		if err != nil {
			return "", fmt.Errorf("postać listu: %w", err)
		}
		rodzaj, _, err := mime.ParseMediaType(postac.Header.Get("Content-Type"))
		if err != nil {
			return "", fmt.Errorf("typ postaci listu: %w", err)
		}
		if rodzaj != "text/html" {
			continue
		}
		tresc, err := io.ReadAll(postac)
		if err != nil {
			return "", fmt.Errorf("odczyt postaci graficznej: %w", err)
		}
		return string(tresc), nil
	}
}
