// Skutek modułu Browser: czy migawka niesie koniec strony, a przy stronie
// ponad granicą — czy rdzeń odmówił i nie zostawił po sobie migawki uciętej.
package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// znacznikKonca stoi na samym końcu ciała strony sprawdzianu. Migawka, która go
// nie niesie, jest migawką początku strony, choć podaje się za migawkę strony.
const znacznikKonca = "OSTATNIE-SLOWO-STRONY"

// stronaSprawdzianu składa dokument HTML o zadanej objętości wypełniacza,
// zakończony znacznikiem końca. Wypełniacz idzie w komentarzu HTML, więc nie
// wchodzi do tekstu renderowanego i nie zaciemnia porównania.
func stronaSprawdzianu(wypelniacz int) string {
	var b strings.Builder
	b.WriteString("<html><head><title>Strona sprawdzianu</title></head><body>")
	b.WriteString("<p>Pierwszy akapit strony.</p>")
	b.WriteString("<script>var pominac = 1;</script>")
	b.WriteString("<!--")
	b.WriteString(strings.Repeat("x", wypelniacz))
	b.WriteString("-->")
	b.WriteString("<p>" + znacznikKonca + "</p></body></html>")
	return b.String()
}

// serwerStrony podnosi witrynę sprawdzianu oddającą zadaną treść ciała strony
// pod wskazaną ścieżką, do wywołania przez sprawdziany migawki.
func serwerStrony(t *testing.T, obsluga http.HandlerFunc) *httptest.Server {
	t.Helper()

	serwer := httptest.NewServer(obsluga)
	t.Cleanup(serwer.Close)
	return serwer
}

// serwerTresci podnosi witrynę oddającą jeden dokument HTML o zadanej treści,
// do sprawdzianów, które nie potrzebują znacznika końca strony.
func serwerTresci(t *testing.T, tresc string) *httptest.Server {
	t.Helper()

	return serwerStrony(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(tresc))
	})
}

// TestMigawkaNiesieKoniecPobranejStronyANieJejPoczatek mierzy, że tekst migawki
// niesie ostatnie zdanie strony, a nie urywa się tam, gdzie skończył się bufor.
func TestMigawkaNiesieKoniecPobranejStronyANieJejPoczatek(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	tresc := stronaSprawdzianu(64 << 10)
	serwer := serwerTresci(t, tresc)

	var przejscie shared.BrowserNavigateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: "okno-przegladarki", Url: serwer.URL}, &przejscie)

	var migawka shared.BrowserSnapshotGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserSnapshotGet,
		shared.BrowserSnapshotGetRequest{WindowId: "okno-przegladarki"}, &migawka)

	if migawka.Snapshot.Id != przejscie.Snapshot.Id {
		t.Errorf("snapshot.get oddał migawkę %q, przejście założyło %q",
			migawka.Snapshot.Id, przejscie.Snapshot.Id)
	}
	if migawka.Snapshot.Url != serwer.URL {
		t.Errorf("migawka opisuje adres %q, pobierano %q", migawka.Snapshot.Url, serwer.URL)
	}
	if migawka.Snapshot.Title == nil || *migawka.Snapshot.Title != "Strona sprawdzianu" {
		t.Errorf("migawka niesie tytuł %v, strona ma „Strona sprawdzianu”", migawka.Snapshot.Title)
	}
	if migawka.Snapshot.Text == nil {
		t.Fatal("migawka bez tekstu renderowanego — model dostałby pustkę jako stan strony")
	}
	if !strings.Contains(*migawka.Snapshot.Text, znacznikKonca) {
		t.Errorf("tekst migawki nie sięga końca strony (%d znaków) — ogryzek podany jako całość",
			len(*migawka.Snapshot.Text))
	}
	if !strings.Contains(*migawka.Snapshot.Text, "Pierwszy akapit strony.") {
		t.Error("tekst migawki nie niesie początku strony")
	}
	// Zawartość skryptu nie jest treścią widzianą przez czytelnika.

	// Nie ma prawa wejść do tekstu renderowanego.
	if strings.Contains(*migawka.Snapshot.Text, "var pominac") {
		t.Error("tekst migawki niesie zawartość skryptu zamiast treści strony")
	}
}

// TestStronaPonadGranicaRozmiaruNieZostawiaMigawkiOgryzka mierzy samą szkodę:
// strona większa niż granica rdzenia ma skończyć się odmową, a okno ma zostać
// bez migawki.
func TestStronaPonadGranicaRozmiaruNieZostawiaMigawkiOgryzka(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	// O jeden bufor ponad granicę — granica ma być granicą, nie sugestią.
	serwer := serwerTresci(t, stronaSprawdzianu(limitTresci+(1<<10)))

	wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: "okno-za-duzej", Url: serwer.URL})

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserSnapshotGet,
		shared.BrowserSnapshotGetRequest{WindowId: "okno-za-duzej"})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa migawki niesie kod %q, okno bez migawki to %q",
			odmowa.Code, shared.ErrorCodeNotFound)
	}
}

// TestStronaDeklarujacaRozmiarPonadGranicaNieJestPobierana pilnuje tej samej
// granicy po drugiej stronie: witryna zapowiada rozmiar większy niż granica,
// więc rdzeń nie ma po co ciągnąć bajta.
func TestStronaDeklarujacaRozmiarPonadGranicaNieJestPobierana(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	serwer := serwerStrony(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprint(limitTresci+1))
		// Ciało krótsze niż zapowiedź: rdzeń ma odmówić na podstawie zapowiedzi,
		// zanim zacznie czytać.
		_, _ = w.Write([]byte(stronaSprawdzianu(0)))
	})

	wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: "okno-zapowiedzi", Url: serwer.URL})

	wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserSnapshotGet,
		shared.BrowserSnapshotGetRequest{WindowId: "okno-zapowiedzi"})
}

// TestOdpowiedzBezTresciNieZakladaMigawkiPrzejscia pilnuje odpowiedzi, które
// treści strony nie niosą wcale. Każda z nich kończyłaby się `ok` i migawką bez
// tytułu, tekstu i HTML-a, a historia nawigacji zapisałaby przejście, którego
// nie było.
func TestOdpowiedzBezTresciNieZakladaMigawkiPrzejscia(t *testing.T) {
	przypadki := map[string]http.HandlerFunc{
		"stan bez ciała": func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
		"przekierowanie bez wskazania dokąd": func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusFound)
		},
		"strona nieodnaleziona": func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "nie ma", http.StatusNotFound)
		},
		"treść nietekstowa": func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write([]byte{0x00, 0x01, 0x02, 0x03})
		},
	}

	for nazwa, obsluga := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
			serwer := serwerStrony(t, obsluga)

			wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserNavigate,
				shared.BrowserNavigateRequest{WindowId: "okno-bez-tresci", Url: serwer.URL})

			odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserSnapshotGet,
				shared.BrowserSnapshotGetRequest{WindowId: "okno-bez-tresci"})
			if odmowa.Code != shared.ErrorCodeNotFound {
				t.Errorf("po nieudanym przejściu okno ma migawkę albo inny powód odmowy: %q", odmowa.Code)
			}
		})
	}
}

// TestMigawkaOddajeZrodloStronyDopieroNaZadanie pilnuje pola, które wychodzi
// warunkowo: wskazanie flagi `includeHtml` ma oddać źródło całe, aż po
// zamknięcie dokumentu.
func TestMigawkaOddajeZrodloStronyDopieroNaZadanie(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	tresc := stronaSprawdzianu(4 << 10)
	serwer := serwerTresci(t, tresc)

	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: "okno-zrodla", Url: serwer.URL}, nil)

	var bezZrodla shared.BrowserSnapshotGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserSnapshotGet,
		shared.BrowserSnapshotGetRequest{WindowId: "okno-zrodla"}, &bezZrodla)
	if bezZrodla.Snapshot.Html != nil {
		t.Error("migawka oddała źródło strony, choć o nie nie proszono")
	}

	var zeZrodlem shared.BrowserSnapshotGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserSnapshotGet,
		shared.BrowserSnapshotGetRequest{
			WindowId:          "okno-zrodla",
			IncludeHtml:       wskaznik(true),
			IncludeScreenshot: wskaznik(true),
		}, &zeZrodlem)

	if zeZrodlem.Snapshot.Html == nil {
		t.Fatal("migawka nie oddała źródła strony mimo wskazania includeHtml")
	}
	if *zeZrodlem.Snapshot.Html != tresc {
		t.Errorf("źródło w migawce ma %d bajtów, strona miała %d — dokument jest niepełny",
			len(*zeZrodlem.Snapshot.Html), len(tresc))
	}
	// Zrzutu ekranu rdzeń nie robi, brak silnika przeglądarki. Pole ma zostać puste.

	// Odsyłacz wskazujący nic byłby migawką kłamiącą o tym, co ma.
	if zeZrodlem.Snapshot.ScreenshotRef != nil {
		t.Errorf("migawka niesie odsyłacz do zrzutu ekranu (%q), którego rdzeń nie wykonuje",
			*zeZrodlem.Snapshot.ScreenshotRef)
	}
}

// TestSnapshotGetOddajeMigawkeNajswiezszegoPrzejscia pilnuje wskaźnika czasu:
// `browser.snapshot.get` pyta o okno, nie o kod migawki, więc ma oddać stan
// bieżący.
func TestSnapshotGetOddajeMigawkeNajswiezszegoPrzejscia(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	pierwsza := serwerTresci(t, "<html><head><title>Strona pierwsza</title></head><body>"+
		"<p>treść pierwsza</p></body></html>")
	druga := serwerTresci(t, "<html><head><title>Strona druga</title></head><body>"+
		"<p>treść druga</p></body></html>")

	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: "okno-historii", Url: pierwsza.URL}, nil)
	var drugie shared.BrowserNavigateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: "okno-historii", Url: druga.URL}, &drugie)

	var migawka shared.BrowserSnapshotGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserSnapshotGet,
		shared.BrowserSnapshotGetRequest{WindowId: "okno-historii"}, &migawka)

	if migawka.Snapshot.Id != drugie.Snapshot.Id {
		t.Errorf("snapshot.get oddał migawkę %q, ostatnie przejście założyło %q",
			migawka.Snapshot.Id, drugie.Snapshot.Id)
	}
	if migawka.Snapshot.Url != druga.URL {
		t.Errorf("migawka opisuje adres %q, okno stoi na %q", migawka.Snapshot.Url, druga.URL)
	}
	if migawka.Snapshot.Text == nil || !strings.Contains(*migawka.Snapshot.Text, "treść druga") {
		t.Errorf("migawka niesie %v zamiast treści strony, na której stoi okno", migawka.Snapshot.Text)
	}
	// Migawka poprzedniego przejścia zostaje w historii, ale nie jest stanem bieżącym.

	// Jej treść nie ma prawa wyjść jako odpowiedź na pytanie o okno.
	if migawka.Snapshot.Text != nil && strings.Contains(*migawka.Snapshot.Text, "treść pierwsza") {
		t.Error("migawka bieżąca niesie treść strony poprzedniej")
	}
}
