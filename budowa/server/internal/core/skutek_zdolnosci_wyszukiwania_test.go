// Sprawdziany skutku dwóch zdolności rodziny `knowledge.*`: przesiewu wyników
// i osi obrazu.
//
// Część z nich żąda wag na dysku i te są pominięte tam, gdzie wag nie ma.
// Pominięcie jest tu jedyną uczciwą odpowiedzią: wagi ważą łącznie blisko
// cztery gigabajty, więc sprawdzian, który by je pobierał, zamieniałby bieg
// sprawdzianów w pobieranie modeli — a sprawdzian, który by ich nie potrzebował,
// mierzyłby atrapę i milczałby dokładnie wtedy, gdy zdolność przestanie działać.
// Katalog wag wskazuje zmienna środowiska `DANACO_MODELE`; nazwy podkatalogów
// są nazwami zdolności, nie wydawców modeli.
//
// Sprawdziany odmowy wag NIE żądają: brak silnika ma być odpowiedzią nazywającą
// brak na każdej maszynie, więc mierzy się go tam, gdzie modelu nie ma z samego
// założenia.
package core

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/store"
	"danacoconsole/server/internal/wiedza"
	"danacoconsole/shared"
)

// zmiennaKatalogModeli wskazuje katalog, w którym leżą wagi modeli tej rodziny.
const zmiennaKatalogModeli = "DANACO_MODELE"

// katalogModeliSprawdzianu oddaje katalog wag albo pomija sprawdzian.
//
// Pominięcie nazywa, czego brakuje i jak to wskazać — pominięcie milczące
// wyglądałoby w wyniku biegu tak samo jak sprawdzian zdany.
func katalogModeliSprawdzianu(t *testing.T, podkatalog string) string {
	t.Helper()

	korzen := strings.TrimSpace(os.Getenv(zmiennaKatalogModeli))
	if korzen == "" {
		t.Skipf("zmienna %s nie wskazuje katalogu wag — bez wag nie ma czym liczyć; "+
			"wskaż katalog niosący podkatalogi `embedder`, `%s` i `%s`",
			zmiennaKatalogModeli, "reranker", "clip")
	}
	katalog := filepath.Join(korzen, podkatalog)
	if _, err := os.Stat(filepath.Join(katalog, "model.safetensors")); err != nil {
		if _, blad := os.Stat(filepath.Join(katalog, "onnx", "model.onnx")); blad != nil {
			t.Skipf("w katalogu %s nie ma wag modelu (%v)", katalog, err)
		}
	}
	return katalog
}

// ustawWiedzy zapisuje nastawę zasięgu globalnego wprost w tabeli ustawień.
//
// Drogą Operatora byłaby komenda `config.set`, ale ta sprawdza klucz wobec
// katalogu ustawień, a wiersze katalogu zakłada migracja nastaw — plik
// z pakietu, którego ten teren nie rusza. Rozstrzyganie nastawy wiersza katalogu
// nie wymaga (`konfig/rozstrzyganie.go` czyta zapisy przed definicjami), więc
// zdolność działa, a niedostępna jest wyłącznie droga jej ustawienia z okna
// konfiguracji. Sprawdzian zapisuje więc to, co zapisałaby komenda, i mierzy to,
// co mierzyć ma — zamiast milczeć o zdolności, dopóki nie powstanie wiersz
// katalogu.
func ustawWiedzy(t *testing.T, katalogDanych, klucz, wartosc string) {
	t.Helper()

	baza, err := store.Otworz(filepath.Join(katalogDanych, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do zapisu nastawy %s: %v", klucz, err)
	}
	defer baza.Zamknij()

	wynik, err := baza.DB.Exec(`
		INSERT INTO ustawienie (poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz,
		                        wartosc, rodzaj_wartosci)
		VALUES ((SELECT id FROM poziom_zasiegu WHERE kod = 'globalny'), '', 'platform', '',
		        ?, ?, 'tekst')
		ON CONFLICT(poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz) DO UPDATE SET
		    wartosc = excluded.wartosc`,
		klucz, wartosc)
	if err != nil {
		t.Fatalf("zapis nastawy %s nie powiódł się: %v", klucz, err)
	}
	if ile, _ := wynik.RowsAffected(); ile == 0 {
		t.Fatalf("zapis nastawy %s nie dotknął ani jednego wiersza — poziom `globalny` "+
			"nie istnieje w tej bazie", klucz)
	}
}

// wgrajDokument wnosi do biblioteki jeden dokument tekstowy.
func wgrajDokument(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	nazwa, tresc string) {

	t.Helper()
	var wgrany shared.LibraryFileUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileUpload,
		shared.LibraryFileUploadRequest{
			Name:          nazwa,
			MimeType:      wskaznik("text/plain"),
			ContentBase64: wskaznik(wBase64([]byte(tresc))),
		}, &wgrany)
}

// granicaKomendyZWagami — ile czasu wolno zająć komendzie, która wczytuje wagi.
//
// Uprząż sprawdzianów skutku daje komendzie piętnaście sekund i dla komend
// odpowiadających z bazy jest to granica słuszna. Tu jest za krótka i nie z
// powodu obciążenia maszyny: samo wczytanie wag rzędu dwóch gigabajtów do
// pamięci trwa dłużej, a dzieje się na każde wołanie (patrz nagłówek
// `wiedza/silnik.go`). Sprawdzian, który by tę granicę przyjął, mierzyłby czas
// wczytania modelu zamiast wyniku, jaki model daje.
const granicaKomendyZWagami = 10 * time.Minute

// wykonajZWagami wykonuje komendę sięgającą po model i przerywa sprawdzian,
// gdy rdzeń odmówił.
func wykonajZWagami(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	komenda shared.MessageType, ladunek any, wynik any) {

	t.Helper()
	koperta, err := protocol.NowaKoperta(komenda, "sprawdzian", "", ladunek)
	if err != nil {
		t.Fatalf("nie można złożyć koperty %s: %v", komenda, err)
	}
	ctx, przerwij := context.WithTimeout(zycie, granicaKomendyZWagami)
	defer przerwij()

	odpowiedz := zmontowany.Rdzen.Wykonaj(ctx, koperta)
	if odpowiedz.Error != nil {
		t.Fatalf("komenda %s odmówiła: kod=%s treść=%s", komenda,
			odpowiedz.Error.Code, odpowiedz.Error.Message)
	}
	if wynik == nil {
		return
	}
	if err := protocol.LadunekDo(odpowiedz, wynik); err != nil {
		t.Fatalf("nieczytelny ładunek odpowiedzi %s: %v", komenda, err)
	}
}

// TestPrzesiewUkladaOdpowiedzInaczejNizPierwszyPrzebieg mierzy to, po co
// przesiew istnieje: tę samą treść i to samo pytanie raz bez niego, raz z nim.
// Odpowiedź identyczna w obu przebiegach znaczyłaby, że drugi model niczego nie
// wnosi — i wtedy sprawdzian ma upaść, choć obie komendy odpowiedziały.
func TestPrzesiewUkladaOdpowiedzInaczejNizPierwszyPrzebieg(t *testing.T) {
	katalogOsadzarki := katalogModeliSprawdzianu(t, "embedder")
	katalogPrzesiewu := katalogModeliSprawdzianu(t, "reranker")

	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	ustawWiedzy(t, katalogDanych, wiedza.KluczKatalogModeli, katalogOsadzarki)
	ustawWiedzy(t, katalogDanych, wiedza.KluczKatalogPrzesiewu, katalogPrzesiewu)

	// Cztery dokumenty o jednym temacie. Pierwszy przebieg widzi w nich podobne
	// rozłożenie znaczeń i wynosi wysoko notatkę, która pytanie POWTARZA, oraz
	// politykę, która o awarii mówi. Dopiero czytanie pytania razem z fragmentem
	// rozstrzyga, że procedura — jedyna, która na pytanie ODPOWIADA, i to bez ani
	// jednego wspólnego z nim wyrazu poza „usługą" — należy wyżej niż polityka.
	wgrajDokument(t, zmontowany, zycie, "notatka-z-pytaniem.txt",
		"Jak przywrócić usługę po awarii serwera? Pytanie wraca po każdej awarii "+
			"serwera i wciąż nie mamy na nie spisanej odpowiedzi.")
	wgrajDokument(t, zmontowany, zycie, "polityka.txt",
		"Polityka dostępności usług przewiduje, że awarie serwera są zgłaszane "+
			"do rejestru zdarzeń, a ich liczba jest raportowana co miesiąc.")
	wgrajDokument(t, zmontowany, zycie, "procedura.txt",
		"Kolejność czynności po zatrzymaniu maszyny: najpierw sprawdź zasilanie "+
			"węzła, potem uruchom go ponownie, następnie odtwórz ostatnią kopię "+
			"zapasową bazy, a na koniec potwierdź odpowiedź usługi z drugiego węzła.")
	wgrajDokument(t, zmontowany, zycie, "harmonogram.txt",
		"Harmonogram przeglądów serwerów: przegląd zasilania w styczniu, "+
			"przegląd kopii zapasowych w lutym, przegląd sieci w marcu.")

	var wskaznikWiedzy shared.KnowledgeIndexResponse
	wykonajZWagami(t, zmontowany, zycie, shared.CommandKnowledgeIndex,
		shared.KnowledgeIndexRequest{Scope: wskaznik(shared.KnowledgeScope(
			shared.KnowledgeScopeLibrary))}, &wskaznikWiedzy)
	if wskaznikWiedzy.Indexed == 0 {
		t.Fatalf("wskaźnik nie objął ani jednej pozycji — nie ma czego przesiewać "+
			"(model: %v)", wskaznikWiedzy.Model)
	}
	t.Logf("wskaźnik: wniesione=%d wszystkie=%d", wskaznikWiedzy.Indexed, wskaznikWiedzy.Total)

	const pytanie = "Co zrobić krok po kroku, żeby przywrócić usługę po awarii serwera?"

	var bezPrzesiewu shared.KnowledgeSearchResponse
	wykonajZWagami(t, zmontowany, zycie, shared.CommandKnowledgeSearch,
		shared.KnowledgeSearchRequest{Query: pytanie, Limit: wskaznik(4)}, &bezPrzesiewu)

	var zPrzesiewem shared.KnowledgeSearchResponse
	wykonajZWagami(t, zmontowany, zycie, shared.CommandKnowledgeSearch,
		shared.KnowledgeSearchRequest{
			Query: pytanie, Limit: wskaznik(4),
			Rerank: wskaznik(true), RerankCandidates: wskaznik(4),
		}, &zPrzesiewem)

	kolejnoscBez := kolejnoscZrodel(bezPrzesiewu.Results)
	kolejnoscZ := kolejnoscZrodel(zPrzesiewem.Results)
	t.Logf("bez przesiewu: %s", opisKolejnosci(bezPrzesiewu.Results))
	t.Logf("z przesiewem:  %s", opisKolejnosci(zPrzesiewem.Results))

	if len(bezPrzesiewu.Results) == 0 || len(zPrzesiewem.Results) == 0 {
		t.Fatal("któryś z przebiegów oddał pustkę — pustka nie jest wynikiem porównania")
	}
	if zPrzesiewem.Reranked == nil || !*zPrzesiewem.Reranked {
		t.Fatal("odpowiedź z przesiewem nie oznajmia przesiewu polem `reranked`")
	}
	if bezPrzesiewu.Reranked != nil && *bezPrzesiewu.Reranked {
		t.Fatal("odpowiedź bez przesiewu oznajmia przesiew, którego nie było")
	}
	// Porównanie idzie po samych źródłach, nie po trafnościach. Trafność zmienia
	// się z definicji — po przesiewie jest oceną innego modelu — więc porównanie
	// obejmujące ją orzekałoby „przesiew coś zmienił" nawet wtedy, gdyby oddał
	// dokładnie tę samą kolejność.
	if kolejnoscBez == kolejnoscZ {
		t.Fatalf("obie odpowiedzi mają tę samą kolejność źródeł — przesiew niczego "+
			"nie przestawił: %s", kolejnoscZ)
	}
	if miejsce(kolejnoscZ, "procedura.txt") > miejsce(kolejnoscZ, "polityka.txt") {
		t.Fatalf("po przesiewie procedura odpowiadająca na pytanie stoi za polityką, "+
			"która o nim tylko wspomina: %s", kolejnoscZ)
	}
	if miejsce(kolejnoscBez, "procedura.txt") < miejsce(kolejnoscBez, "polityka.txt") {
		t.Fatalf("pierwszy przebieg sam ustawił procedurę przed polityką, więc "+
			"sprawdzian nie mierzy tego, co przesiew wnosi: %s", kolejnoscBez)
	}
}

// kolejnoscZrodel składa same źródła odpowiedzi w jeden napis — to on
// rozstrzyga, czy przesiew przestawił kolejność.
func kolejnoscZrodel(trafienia []shared.KnowledgeHit) string {
	czlony := make([]string, 0, len(trafienia))
	for _, trafienie := range trafienia {
		czlony = append(czlony, trafienie.Source)
	}
	return strings.Join(czlony, " ")
}

// miejsce oddaje pozycję źródła w kolejności; brak źródła jest pozycją dalszą
// niż każda obecna, więc porównania nie trzeba obwarowywać osobnym warunkiem.
func miejsce(kolejnosc, zrodlo string) int {
	for numer, czlon := range strings.Fields(kolejnosc) {
		if czlon == zrodlo {
			return numer
		}
	}
	return len(kolejnosc) + 1
}

// opisKolejnosci składa kolejność wraz z trafnościami — do przytoczenia
// w wyniku biegu, nie do rozstrzygania.
func opisKolejnosci(trafienia []shared.KnowledgeHit) string {
	czlony := make([]string, 0, len(trafienia))
	for _, trafienie := range trafienia {
		trafnosc := 0
		if trafienie.Score != nil {
			trafnosc = *trafienie.Score
		}
		czlony = append(czlony, trafienie.Source+"="+strconv.Itoa(trafnosc))
	}
	return strings.Join(czlony, " ")
}

// TestOsObrazuOddajeObrazOpisanyZdaniem mierzy trafienie, nie samą odpowiedź:
// wśród trzech obrazów różniących się wyłącznie tym, co na nich widać, zdanie
// ma wskazać ten właściwy.
func TestOsObrazuOddajeObrazOpisanyZdaniem(t *testing.T) {
	katalogObrazu := katalogModeliSprawdzianu(t, "clip")

	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	ustawWiedzy(t, katalogDanych, wiedza.KluczKatalogObrazu, katalogObrazu)

	wgrajObraz(t, zmontowany, zycie, "pierwszy.png", kolo(kolorCzerwony))
	wgrajObraz(t, zmontowany, zycie, "drugi.png", kwadrat(kolorNiebieski))
	wgrajObraz(t, zmontowany, zycie, "trzeci.png", kwadrat(kolorZielony))

	var wynik shared.KnowledgeImageSearchResponse
	wykonajZWagami(t, zmontowany, zycie, shared.CommandKnowledgeImageSearch,
		shared.KnowledgeImageSearchRequest{
			Query: "a red circle on a white background", Limit: wskaznik(3),
		}, &wynik)

	if wynik.Examined == nil || *wynik.Examined != 3 {
		t.Fatalf("do porównania weszło %v obrazów zamiast trzech wgranych", wynik.Examined)
	}
	if len(wynik.Results) == 0 {
		t.Fatal("oś obrazu oddała pustkę przy trzech obrazach w bibliotece")
	}
	for _, trafienie := range wynik.Results {
		trafnosc := 0
		if trafienie.Score != nil {
			trafnosc = *trafienie.Score
		}
		rodzaj := ""
		if trafienie.MimeType != nil {
			rodzaj = *trafienie.MimeType
		}
		t.Logf("obraz %s (%s) trafność=%d", trafienie.Source, rodzaj, trafnosc)
	}
	if wynik.Results[0].Source != "pierwszy.png" {
		t.Fatalf("na czele stoi %q, choć zdanie opisuje czerwone koło z pliku pierwszy.png",
			wynik.Results[0].Source)
	}
	if wynik.Results[0].SourceId == nil || *wynik.Results[0].SourceId == "" {
		t.Fatal("trafienie nie niesie identyfikatora, którym da się po obraz sięgnąć")
	}
}

// TestOsObrazuBezSilnikaOdmawiaNazywajacBrak pilnuje, że brak zaplecza jest
// odmową nazywającą brak i drogę naprawy, a nie usterką wewnętrzną. Wag nie
// żąda: interpreter wskazany nastawą nie istnieje, więc pomocnik nie ruszy na
// żadnej maszynie i odmowa jest ta sama wszędzie.
func TestOsObrazuBezSilnikaOdmawiaNazywajacBrak(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	ustawWiedzy(t, katalog, wiedza.KluczProgram,
		filepath.Join(katalog, "nie-ma-takiego-interpretera"))

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandKnowledgeImageSearch,
		shared.KnowledgeImageSearchRequest{Query: "czerwone koło"})

	t.Logf("odmowa: kod=%s treść=%s", odmowa.Code, odmowa.Message)
	if odmowa.Code != shared.ErrorCodeChannelUnavailable {
		t.Fatalf("brak zaplecza dostał kod %s zamiast %s — usterka wewnętrzna nie mówi, "+
			"czego brakuje", odmowa.Code, shared.ErrorCodeChannelUnavailable)
	}
	for _, czlon := range []string{"oś obrazu", "naprawa:", "pillow"} {
		if !strings.Contains(odmowa.Message, czlon) {
			t.Fatalf("odmowa nie niesie członu %q: %s", czlon, odmowa.Message)
		}
	}
}

// TestOsObrazuBezZdaniaOdmawiaWadaZadania pilnuje, że pytanie puste jest wadą
// żądania, a nie pustą odpowiedzią udającą brak obrazów.
func TestOsObrazuBezZdaniaOdmawiaWadaZadania(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandKnowledgeImageSearch,
		shared.KnowledgeImageSearchRequest{Query: "   "})

	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("zdanie puste dostało kod %s zamiast %s", odmowa.Code,
			shared.ErrorCodeValidationFailed)
	}
	t.Logf("odmowa: %s", odmowa.Message)
}

// TestBramaOdmawiaOsiObrazuBezPolaWymaganego pilnuje, że nowa komenda przechodzi
// przez bramę kontraktu tak samo jak sąsiedzi: treść niepełna dostaje odmowę
// nazywającą brakujące pole, zanim dojdzie do adaptera.
func TestBramaOdmawiaOsiObrazuBezPolaWymaganego(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandKnowledgeImageSearch,
		map[string]any{"limit": 3})

	t.Logf("odmowa bramy: kod=%s treść=%s", odmowa.Code, odmowa.Message)
	if !strings.Contains(odmowa.Message, "query") {
		t.Fatalf("odmowa bramy nie nazywa brakującego pola `query`: %s", odmowa.Message)
	}
}

// Barwy obrazów sprawdzianu. Trzy barwy i dwa kształty wystarczą, żeby zdanie
// miało co rozróżniać, a obraz rysowany w kodzie nie wnosi do repozytorium
// plików binarnych.
var (
	kolorCzerwony  = color.RGBA{R: 220, G: 20, B: 20, A: 255}
	kolorNiebieski = color.RGBA{R: 20, G: 40, B: 220, A: 255}
	kolorZielony   = color.RGBA{R: 20, G: 180, B: 60, A: 255}
)

// bokObrazu — obraz kwadratowy o boku wystarczającym, żeby po sprowadzeniu do
// wejścia modelu kształt pozostał rozpoznawalny.
const bokObrazu = 320

// kolo rysuje wypełnione koło na białym tle.
func kolo(barwa color.RGBA) image.Image {
	plotno := bialePlotno()
	srodek, promien := bokObrazu/2, bokObrazu/3
	for y := 0; y < bokObrazu; y++ {
		for x := 0; x < bokObrazu; x++ {
			dx, dy := x-srodek, y-srodek
			if dx*dx+dy*dy <= promien*promien {
				plotno.Set(x, y, barwa)
			}
		}
	}
	return plotno
}

// kwadrat rysuje wypełniony kwadrat na białym tle.
func kwadrat(barwa color.RGBA) image.Image {
	plotno := bialePlotno()
	odstep := bokObrazu / 5
	for y := odstep; y < bokObrazu-odstep; y++ {
		for x := odstep; x < bokObrazu-odstep; x++ {
			plotno.Set(x, y, barwa)
		}
	}
	return plotno
}

// bialePlotno zakłada białe tło — obraz przezroczysty sprowadzony do trzech
// barw wychodzi czarny i przestaje być tym, co miał przedstawiać.
func bialePlotno() *image.RGBA {
	plotno := image.NewRGBA(image.Rect(0, 0, bokObrazu, bokObrazu))
	bialy := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for y := 0; y < bokObrazu; y++ {
		for x := 0; x < bokObrazu; x++ {
			plotno.Set(x, y, bialy)
		}
	}
	return plotno
}

// wgrajObraz wnosi do biblioteki jeden obraz zapisany w PNG.
func wgrajObraz(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	nazwa string, obraz image.Image) {

	t.Helper()
	bufor := bytes.Buffer{}
	if err := png.Encode(&bufor, obraz); err != nil {
		t.Fatalf("nie da się zapisać obrazu %s: %v", nazwa, err)
	}
	var wgrany shared.LibraryFileUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileUpload,
		shared.LibraryFileUploadRequest{
			Name:          nazwa,
			MimeType:      wskaznik("image/png"),
			ContentBase64: wskaznik(wBase64(bufor.Bytes())),
		}, &wgrany)
}
