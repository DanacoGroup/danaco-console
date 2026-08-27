package core

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/store"
	"danacoconsole/shared"
)

// Uprząż sprawdzianów skutku pyta, czy po udanej odpowiedzi w świecie naprawdę coś zostało.

// zmontujDoPomiaruSkutku składa rdzeń nad świeżą bazą i oddaje go wraz
// z kontekstem życia oraz katalogiem danych, w którym leżą magazyny treści.
func zmontujDoPomiaruSkutku(t *testing.T) (*Zmontowany, context.Context, string) {
	t.Helper()

	katalog := t.TempDir()
	baza, err := store.Otworz(filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy sprawdzianu: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })

	zycie, zakoncz := context.WithCancel(context.Background())
	t.Cleanup(zakoncz)

	ustawienia := konfiguracja.Domyslna()
	ustawienia.KatalogDanych = katalog
	ustawienia.KatalogKlienta = ""
	ustawienia.KatalogProfili = ""
	ustawienia.Port = 0

	zmontowany, err := Zmontuj(zycie, Montaz{
		Konfiguracja: ustawienia,
		Baza:         baza,
		Dziennik:     dziennikNiemy(),
	})
	if err != nil {
		t.Fatalf("montaż rdzenia nie powiódł się: %v", err)
	}
	t.Cleanup(zmontowany.Zamknij)

	return zmontowany, zycie, katalog
}

// wykonajUdana wywołuje komendę i przerywa sprawdzian, gdy rdzeń odmówił.
// Sprawdziany skutku mierzą to, co zostaje po odpowiedzi udanej, więc odmowa
// w miejscu, gdzie oczekiwano powodzenia, nie ma prawa iść dalej jako pusta
// treść.
func wykonajUdana(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	komenda shared.MessageType, ladunek any, wynik any) {
	t.Helper()

	odpowiedz := wykonajKomende(t, zmontowany, zycie, komenda, ladunek)
	if odpowiedz.Error != nil {
		t.Fatalf("komenda %s odmówiła: kod=%s treść=%s", komenda,
			odpowiedz.Error.Code, odpowiedz.Error.Message)
	}
	if odpowiedz.Status == nil || *odpowiedz.Status != shared.EnvelopeStatusOk {
		t.Fatalf("komenda %s nie oddała stanu udanego: %+v", komenda, odpowiedz.Status)
	}
	if wynik == nil {
		return
	}
	if err := protocol.LadunekDo(odpowiedz, wynik); err != nil {
		t.Fatalf("nieczytelny ładunek odpowiedzi %s: %v", komenda, err)
	}
}

// wykonajOdmowna wywołuje komendę, która ma odmówić, i oddaje jej odmowę.
// Odmowa nazwana jest tu wynikiem oczekiwanym: rdzeń, który zamiast niej oddaje
// `ok`, melduje skutek, którego nie ma.
func wykonajOdmowna(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	komenda shared.MessageType, ladunek any) protocol.Blad {
	t.Helper()

	odpowiedz := wykonajKomende(t, zmontowany, zycie, komenda, ladunek)
	if odpowiedz.Error == nil {
		t.Fatalf("komenda %s zameldowała powodzenie tam, gdzie skutku nie ma; początek ładunku: %s",
			komenda, poczatekLadunku(odpowiedz.Payload))
	}
	if odpowiedz.Status == nil || *odpowiedz.Status != shared.EnvelopeStatusError {
		t.Fatalf("odmowa %s bez stanu błędu — klient stanie w wiecznym ładowaniu: %+v",
			komenda, odpowiedz.Status)
	}
	return *odpowiedz.Error
}

// poczatekLadunku przycina ładunek do wielkości czytelnej w dzienniku sprawdzianu. Ładunki niosące treść strony albo obraz idą w megabajtach, a niepowodzenie ma być czytelne, więc do rozpoznania wystarcza początek.
func poczatekLadunku(ladunek []byte) string {
	const granica = 512
	if len(ladunek) <= granica {
		return string(ladunek)
	}
	return string(ladunek[:granica]) + "… (łącznie " + strconv.Itoa(len(ladunek)) + " bajtów)"
}

// sciezkaWMagazynie składa ścieżkę na dysku z odwołania wypuszczonego przez
// rdzeń. Odwołanie jest względne wobec katalogu danych i zawsze z ukośnikiem
// `/` (odwolanieMagazynu), więc przełożenie na ścieżkę systemu idzie przez
// `filepath.FromSlash`.
func sciezkaWMagazynie(katalogDanych, odwolanie string) string {
	return filepath.Join(katalogDanych, filepath.FromSlash(odwolanie))
}

// bajtyPodOdwolaniem czyta treść leżącą pod odwołaniem wypuszczonym z rdzenia i przerywa sprawdzian, gdy pliku nie ma albo jest pusty. Plik pusty jest osobnym niepowodzeniem, nie odmianą braku, bo nie ma czego pokazać.
func bajtyPodOdwolaniem(t *testing.T, katalogDanych, odwolanie string) []byte {
	t.Helper()

	if odwolanie == "" {
		t.Fatal("odwołanie jest puste — nie ma czego czytać")
	}
	sciezka := sciezkaWMagazynie(katalogDanych, odwolanie)
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("odwołanie %q nie prowadzi do treści (%s): %v", odwolanie, sciezka, err)
	}
	if len(bajty) == 0 {
		t.Fatalf("odwołanie %q prowadzi do pliku o zerowej długości (%s)", odwolanie, sciezka)
	}
	return bajty
}

// zapiszPlikSprawdzianu kładzie treść na dysku pod wskazaną ścieżką. Służy drogom, w których żądanie niesie sourcePath, oraz sprawdzeniu, czy rdzeń poszedł za nadpisanym plikiem źródłowym, czy zamroził bajty u siebie.
func zapiszPlikSprawdzianu(t *testing.T, sciezka string, tresc []byte) {
	t.Helper()

	if err := os.WriteFile(sciezka, tresc, 0o600); err != nil {
		t.Fatalf("nie można zapisać pliku sprawdzianu %s: %v", sciezka, err)
	}
}

// sumaSha256 liczy sumę kontrolną treści w postaci, w której niosą ją kontrakt
// (`checksum`) i nazwa bloba w magazynie.
func sumaSha256(bajty []byte) string {
	suma := sha256.Sum256(bajty)
	return hex.EncodeToString(suma[:])
}

// obrazPNG składa prawdziwy plik PNG o zadanych wymiarach. Sprawdziany zasobów
// mierzą format i wymiary odczytane z nagłówka utrwalonego pliku, więc treść
// musi być obrazem naprawdę, a nie napisem udającym obraz.
func obrazPNG(t *testing.T, szerokosc, wysokosc int) []byte {
	t.Helper()

	plotno := image.NewRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for x := 0; x < szerokosc; x++ {
		for y := 0; y < wysokosc; y++ {
			plotno.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 0x40, A: 0xff})
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// wBase64 koduje treść tak, jak niosą ją pola contentBase64 kontraktu, standardowym kodowaniem Base64.
func wBase64(bajty []byte) string {
	return base64.StdEncoding.EncodeToString(bajty)
}

// jsonSurowy składa ładunek konfiguracji kanału. `channel.add` przyjmuje
// `config` jako surowy JSON, więc sprawdzian buduje go tą samą drogą, którą
// budowałby go klient.
func jsonSurowy(t *testing.T, wartosc any) json.RawMessage {
	t.Helper()

	surowe, err := json.Marshal(wartosc)
	if err != nil {
		t.Fatalf("nie można złożyć konfiguracji: %v", err)
	}
	return surowe
}

// wskaznik oddaje wskaźnik na wartość — pola opcjonalne kontraktu są
// wskaźnikami, a literału adresu wziąć nie można.
func wskaznik[T any](wartosc T) *T {
	return &wartosc
}
