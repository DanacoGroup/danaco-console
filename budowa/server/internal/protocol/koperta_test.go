// Koperta jest jedynym kształtem komunikatu w obie strony, więc jej
// sprawdziany są sprawdzianami całej warstwy styku klienta z rdzeniem: pola
// opcjonalne znikają, gdy są puste.
package protocol

import (
	"encoding/json"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// TestKopertaWracaZDrutuBezZmiany sprawdza obieg zamknięty: złożenie, zapis do
// bajtów i odczyt z powrotem.
func TestKopertaWracaZDrutuBezZmiany(t *testing.T) {
	type tresc struct {
		Tytul string `json:"title"`
		Ile   int    `json:"count"`
	}

	zlozona, err := NowaKoperta(shared.CommandSessionCreate, "jeden", "sesja-pierwsza",
		tresc{Tytul: "Praca nad ofertą", Ile: 3})
	if err != nil {
		t.Fatalf("nie można złożyć koperty: %v", err)
	}

	bajty, err := Zakoduj(zlozona)
	if err != nil {
		t.Fatalf("nie można zakodować koperty: %v", err)
	}
	odczytana, err := Odkoduj(bajty)
	if err != nil {
		t.Fatalf("nie można odkodować koperty: %v", err)
	}

	if odczytana.Type != zlozona.Type {
		t.Errorf("typ zmienił się w obiegu: %q → %q", zlozona.Type, odczytana.Type)
	}
	if odczytana.Id != zlozona.Id {
		t.Errorf("identyfikator zmienił się w obiegu: %q → %q", zlozona.Id, odczytana.Id)
	}
	if IdSesji(odczytana) != "sesja-pierwsza" {
		t.Errorf("sesja zmieniła się w obiegu: %q", IdSesji(odczytana))
	}

	var wrocona tresc
	if err := LadunekDo(odczytana, &wrocona); err != nil {
		t.Fatalf("nie można odczytać ładunku: %v", err)
	}
	if wrocona.Tytul != "Praca nad ofertą" || wrocona.Ile != 3 {
		t.Errorf("ładunek zmienił się w obiegu: %+v", wrocona)
	}
}

// TestPolaOpcjonalneZnikajaGdyPuste pilnuje znacznika omitempty kontraktu.
// Pole `sessionId` wysłane jako `null` znaczy dla klienta co innego niż pole
// nieobecne, a powitanie idzie właśnie bez sesji.
func TestPolaOpcjonalneZnikajaGdyPuste(t *testing.T) {
	koperta, err := NowaKoperta(shared.CommandConnectionHello, "jeden", "", nil)
	if err != nil {
		t.Fatalf("nie można złożyć koperty: %v", err)
	}
	bajty, err := Zakoduj(koperta)
	if err != nil {
		t.Fatalf("nie można zakodować koperty: %v", err)
	}

	zapis := string(bajty)
	for _, pole := range []string{"sessionId", "payload", "status", "error", "seq", "done"} {
		if strings.Contains(zapis, `"`+pole+`"`) {
			t.Errorf("koperta bez treści niesie puste pole %q: %s", pole, zapis)
		}
	}
	if !strings.Contains(zapis, `"type"`) || !strings.Contains(zapis, `"id"`) {
		t.Errorf("koperta zgubiła pole wymagane: %s", zapis)
	}
}

// TestKopertaBezTypuJestBledemStrukturalnym oddziela dwa różne niepowodzenia.
// Brak pola `type` to komunikat niepoprawny — czym innym jest typ, którego
// rdzeń nie zna, bo ten wraca zdarzeniem `*.unknown` bez błędu dekodowania.
func TestKopertaBezTypuJestBledemStrukturalnym(t *testing.T) {
	if _, err := Odkoduj([]byte(`{"id":"jeden"}`)); err == nil {
		t.Error("koperta bez pola type została przyjęta")
	}
	if _, err := Odkoduj([]byte(`{"type":"nazwa.spoza.kontraktu","id":"jeden"}`)); err != nil {
		t.Errorf("typ nieznany został potraktowany jak komunikat niepoprawny: %v", err)
	}
}

// TestOdkodowanieOdrzucaSmieci sprawdza wejście, którym idzie każdy bajt
// z gniazda WebSocket tego rdzenia.
func TestOdkodowanieOdrzucaSmieci(t *testing.T) {
	przypadki := map[string][]byte{
		"pusty ciąg":              {},
		"tekst spoza JSON":        []byte("halo"),
		"tablica zamiast obiektu": []byte(`[1,2,3]`),
		"urwany dokument":         []byte(`{"type":"connection.hello"`),
	}
	for nazwa, bajty := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			if _, err := Odkoduj(bajty); err == nil {
				t.Error("bajty niebędące kopertą zostały przyjęte")
			}
		})
	}
}

// TestPustyLadunekNieNadpisujeCelu sprawdza obietnicę z opisu LadunekDo: brak
// treści nie jest błędem i nie tknie struktury docelowej.
func TestPustyLadunekNieNadpisujeCelu(t *testing.T) {
	cel := struct {
		Tytul string `json:"title"`
	}{Tytul: "wartość sprzed odczytu"}

	if err := LadunekDo(Koperta{Type: shared.CommandSessionList}, &cel); err != nil {
		t.Fatalf("pusty ładunek zgłosił błąd: %v", err)
	}
	if cel.Tytul != "wartość sprzed odczytu" {
		t.Errorf("pusty ładunek nadpisał cel: %q", cel.Tytul)
	}
}

// TestLadunekInnegoKsztaltuJestBledem pilnuje, by niezgodność kształtu wyszła
// błędem, a nie cichym zerem.
func TestLadunekInnegoKsztaltuJestBledem(t *testing.T) {
	koperta := Koperta{
		Type:    shared.CommandSessionCreate,
		Payload: json.RawMessage(`"napis zamiast obiektu"`),
	}
	var cel struct {
		Tytul string `json:"title"`
	}
	if err := LadunekDo(koperta, &cel); err == nil {
		t.Error("ładunek innego kształtu został przyjęty")
	}
}

// TestPolaStrumieniaCzytaneZKoperty sprawdza odczyt numeru i znacznika końca.
// Kontrakt liczy fragmenty od jedynki, więc zero znaczy „komunikat spoza
// strumienia" i nie wolno go mylić z pierwszym fragmentem.
func TestPolaStrumieniaCzytaneZKoperty(t *testing.T) {
	pusta := Koperta{Type: shared.EventStreamChunk}
	if Numer(pusta) != 0 {
		t.Errorf("koperta spoza strumienia ma numer %d, oczekiwane 0", Numer(pusta))
	}
	if Ostatni(pusta) {
		t.Error("koperta bez znacznika końca została uznana za domykającą")
	}

	pierwszy := Koperta{Type: shared.EventStreamChunk, Seq: wskaznik(1)}
	if Numer(pierwszy) != 1 {
		t.Errorf("pierwszy fragment ma numer %d, oczekiwane 1", Numer(pierwszy))
	}
	if Ostatni(pierwszy) {
		t.Error("fragment bez znacznika końca został uznany za domykający")
	}

	niedomykajacy := Koperta{Type: shared.EventStreamChunk, Seq: wskaznik(2), Done: wskaznik(false)}
	if Ostatni(niedomykajacy) {
		t.Error("znacznik końca równy fałszowi został uznany za domknięcie")
	}

	domykajacy := Koperta{Type: shared.EventStreamChunk, Seq: wskaznik(3), Done: wskaznik(true)}
	if !Ostatni(domykajacy) {
		t.Error("znacznik końca równy prawdzie nie domknął strumienia")
	}
}

// TestLadunekNieDaSieZakodowac sprawdza drogę odmowy przy budowie koperty.
// Kanał jest wartością, której encoding/json nie potrafi zapisać.
func TestLadunekNieDaSieZakodowac(t *testing.T) {
	if _, err := NowaKoperta(shared.CommandSessionCreate, "jeden", "", make(chan int)); err == nil {
		t.Error("ładunek niedający się zakodować został przyjęty")
	}
}
