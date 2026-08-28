// Sprawdziany tego pliku mierzą skutek komendy modułu Browser niezależnie od jej
// odpowiedzi: przez drugie połączenie do bazy rdzenia albo przez bajty magazynu
// treści pod odwołaniem.
package core

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/server/internal/store"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// bazaSprawdzianuPrzegladarki otwiera drugie połączenie do bazy rdzenia.
// Pomiar idzie własnym połączeniem, a nie przez rdzeń: gdyby wykaz komendy
// i pomiar czytały tę samą warstwę, sprawdzian potwierdzałby sam siebie.
func bazaSprawdzianuPrzegladarki(t *testing.T, katalogDanych string) *sql.DB {
	t.Helper()

	baza, err := store.Otworz(filepath.Join(katalogDanych, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do pomiaru: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })
	return baza.DB
}

// policzWierszePrzegladania liczy wiersze zwrócone wskazanym zapytaniem SQL i przerywa
// sprawdzian niepowodzeniem, gdy zapytanie się nie wykona.
func policzWierszePrzegladania(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) int {
	t.Helper()

	var liczba int
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&liczba); err != nil {
		t.Fatalf("pomiar w bazie nie powiódł się (%s): %v", zapytanie, err)
	}
	return liczba
}

// chromiumStoi mówi, czy na tej maszynie jest zainstalowany program Chromium,
// którego wymagają sprawdziany zależne od uruchomionej strony.
func chromiumStoi() bool {
	return zewnetrzne.Stoi(narzedzieChromium())
}

// stronaZTytulem składa prosty dokument HTML o zadanym tytule i treści akapitu,
// używany jako strona serwowana testowemu serwerowi HTTP.
func stronaZTytulem(tytul, tresc string) string {
	return "<html><head><title>" + tytul + "</title></head><body><p id=\"tresc\">" +
		tresc + "</p></body></html>"
}

// TestZrzutStronyNiesieBajtyObrazuAOdczytOddajeTeSameBajty sprawdza zrzut ekranu
// miarą pliku w magazynie: obrazu PNG o dodatnich wymiarach, a nie samego pola
// odwołania w odpowiedzi komendy.
func TestZrzutStronyNiesieBajtyObrazuAOdczytOddajeTeSameBajty(t *testing.T) {
	if !chromiumStoi() {
		t.Skip("na tej maszynie nie ma Chromium — zrzut strony nie ma czym powstać")
	}
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	serwer := serwerTresci(t, stronaZTytulem("Strona zrzutu", "Treść widoczna na zrzucie"))

	var przejscie shared.BrowserNavigateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: "okno-zrzutu", Url: serwer.URL}, &przejscie)

	var zrzut shared.BrowserScreenshotCaptureResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserScreenshotCapture,
		shared.BrowserScreenshotCaptureRequest{
			WindowId: "okno-zrzutu", Mode: shared.BrowserScreenshotModeFullPage,
		}, &zrzut)

	if zrzut.Screenshot.Ref == "" {
		t.Fatal("zrzut wrócił bez odwołania — nie ma czego czytać")
	}
	bajty := bajtyPodOdwolaniem(t, katalog, zrzut.Screenshot.Ref)
	obraz, err := png.DecodeConfig(bytes.NewReader(bajty))
	if err != nil {
		t.Fatalf("treść pod odwołaniem zrzutu nie jest obrazem PNG: %v", err)
	}
	if obraz.Width <= 0 || obraz.Height <= 0 {
		t.Fatalf("zrzut ma wymiary %dx%d — obraz bez powierzchni", obraz.Width, obraz.Height)
	}

	var odczyt shared.BrowserSnapshotScreenshotGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserSnapshotScreenshotGet,
		shared.BrowserSnapshotScreenshotGetRequest{
			ScreenshotRef: wskaznik(zrzut.Screenshot.Ref), IncludeContent: wskaznik(true),
		}, &odczyt)
	if odczyt.Screenshot.ContentBase64 == nil {
		t.Fatal("odczyt zrzutu nie oddał treści, choć o nią poproszono")
	}
	oddane, err := base64.StdEncoding.DecodeString(*odczyt.Screenshot.ContentBase64)
	if err != nil {
		t.Fatalf("treść zrzutu nie jest zapisem base64: %v", err)
	}
	if !bytes.Equal(oddane, bajty) {
		t.Fatalf("odczyt oddał %d bajtów, a w magazynie leży %d — to nie ten sam zrzut",
			len(oddane), len(bajty))
	}
}

// TestMonitorMierzyZmianeStronyAOdniesienieLezyWMagazynie sprawdza, czy monitor
// naprawdę mierzy: odniesienie ma leżeć w magazynie jako bajty, a sprawdzenie po
// zmianie ma oddać niezerową różnicę wierszy.
func TestMonitorMierzyZmianeStronyAOdniesienieLezyWMagazynie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuPrzegladarki(t, katalog)

	tresc := stronaZTytulem("Cennik", "Cena: 100 zł")
	serwer := serwerStrony(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(tresc))
	})

	var zalozenie shared.BrowserMonitorAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserMonitorAdd,
		shared.BrowserMonitorAddRequest{WindowId: "okno-monitora", Url: serwer.URL}, &zalozenie)

	if zalozenie.Monitor.BaselineRef == nil {
		t.Fatal("monitor powstał bez odniesienia — nie ma czego porównywać")
	}
	odniesienie := bajtyPodOdwolaniem(t, katalog, *zalozenie.Monitor.BaselineRef)
	if !strings.Contains(string(odniesienie), "Cena: 100") {
		t.Fatalf("odniesienie monitora nie niesie treści pilnowanej strony: %q", string(odniesienie))
	}

	var bezZmiany shared.BrowserMonitorCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserMonitorCheck,
		shared.BrowserMonitorCheckRequest{MonitorId: zalozenie.Monitor.Id}, &bezZmiany)
	if bezZmiany.Changed {
		t.Fatal("monitor zameldował zmianę na stronie, która się nie zmieniła")
	}

	tresc = stronaZTytulem("Cennik", "Cena: 145 zł")
	var poZmianie shared.BrowserMonitorCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserMonitorCheck,
		shared.BrowserMonitorCheckRequest{
			MonitorId: zalozenie.Monitor.Id, UpdateBaseline: wskaznik(true),
		}, &poZmianie)
	if !poZmianie.Changed {
		t.Fatal("monitor nie zauważył zmiany ceny na pilnowanej stronie")
	}
	if poZmianie.Diff == nil || poZmianie.Diff.AddedLines+poZmianie.Diff.RemovedLines == 0 {
		t.Fatalf("monitor oddał zmianę bez miary różnicy: %+v", poZmianie.Diff)
	}
	if poZmianie.Monitor.Status != shared.BrowserMonitorStatusChanged {
		t.Fatalf("stan monitora po zmianie to %q, oczekiwano %q",
			poZmianie.Monitor.Status, shared.BrowserMonitorStatusChanged)
	}

	// Miara niezależna: wiersz monitora ma czas sprawdzenia i zmiany, a odniesienie
	// wskazuje treść nową.
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM monitor_przegladania
		 WHERE identyfikator_zewnetrzny = ? AND sprawdzono IS NOT NULL AND zmieniono IS NOT NULL`,
		zalozenie.Monitor.Id); liczba != 1 {
		t.Fatalf("baza nie odnotowała sprawdzenia i zmiany monitora (wierszy: %d)", liczba)
	}
	nowe := bajtyPodOdwolaniem(t, katalog, *poZmianie.Monitor.BaselineRef)
	if !strings.Contains(string(nowe), "145") {
		t.Fatalf("odniesienie po sprawdzeniu nie zostało przesunięte: %q", string(nowe))
	}
}

// TestSubskrypcjaKanaluOdkladaWpisyWBazie sprawdza, czy subskrypcja naprawdę
// odpytała kanał: wpisy mają leżeć w tabeli, a nie tylko w odpowiedzi komendy.
func TestSubskrypcjaKanaluOdkladaWpisyWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuPrzegladarki(t, katalog)

	kanal := `<?xml version="1.0"?><rss version="2.0"><channel>
		<title>Kanał branżowy</title>
		<item><title>Pierwszy wpis</title><link>https://przyklad.test/1</link>
			<description>Treść pierwszego</description><pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate></item>
		<item><title>Drugi wpis</title><link>https://przyklad.test/2</link>
			<description>Treść drugiego</description></item>
	</channel></rss>`
	serwer := serwerStrony(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(kanal))
	})

	var subskrypcja shared.BrowserFeedSubscribeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserFeedSubscribe,
		shared.BrowserFeedSubscribeRequest{WindowId: "okno-kanalow", Url: serwer.URL}, &subskrypcja)

	if subskrypcja.Feed.Format == nil || *subskrypcja.Feed.Format != shared.BrowserFeedFormatRss {
		t.Fatalf("kanał rozpoznany jako %v, oczekiwano rss", subskrypcja.Feed.Format)
	}
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM wpis_kanalu_przegladania WHERE kanal_zewnetrzny_id = ?`,
		subskrypcja.Feed.Id); liczba != 2 {
		t.Fatalf("w bazie leży %d wpisów kanału, a kanał niósł 2", liczba)
	}

	var wykaz shared.BrowserFeedListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserFeedList,
		shared.BrowserFeedListRequest{
			WindowId: wskaznik("okno-kanalow"), IncludeEntries: wskaznik(true),
		}, &wykaz)
	if len(wykaz.Feeds) != 1 || len(wykaz.Feeds[0].Entries) != 2 {
		t.Fatalf("wykaz kanałów nie oddał wpisów: %+v", wykaz.Feeds)
	}

	// Zdjęcie kanału zabiera wpisy kaskadą schematu bazy, nie sprzątaniem — miarą
	// jest stan bazy.
	var zdjecie shared.BrowserFeedRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserFeedRemove,
		shared.BrowserFeedRemoveRequest{FeedId: subskrypcja.Feed.Id}, &zdjecie)
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM wpis_kanalu_przegladania WHERE kanal_zewnetrzny_id = ?`,
		subskrypcja.Feed.Id); liczba != 0 {
		t.Fatalf("po zdjęciu kanału w bazie zostało %d jego wpisów", liczba)
	}
}

// TestKartyIPrzestrzenRoboczaZostajaWBazie sprawdza rząd kart: karta otwarta ma
// mieć wiersz i migawkę odwiedzonej strony, karta zamknięta ma zniknąć z wykazu,
// a przestrzeń robocza ma pamiętać skład liczony z bazy.
func TestKartyIPrzestrzenRoboczaZostajaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuPrzegladarki(t, katalog)
	serwer := serwerTresci(t, stronaZTytulem("Karta pierwsza", "Treść karty"))

	var pierwsza shared.BrowserTabOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserTabOpen,
		shared.BrowserTabOpenRequest{WindowId: "okno-kart", Url: wskaznik(serwer.URL)}, &pierwsza)
	if pierwsza.Snapshot == nil || !strings.Contains(pierwsza.Snapshot.Url, serwer.URL) {
		t.Fatalf("karta otwarta z adresem nie przyniosła migawki strony: %+v", pierwsza.Snapshot)
	}
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM migawka_strony WHERE okno = ?`, "okno-kart"); liczba == 0 {
		t.Fatal("otwarcie karty nie odłożyło migawki strony w bazie")
	}

	var druga shared.BrowserTabOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserTabOpen,
		shared.BrowserTabOpenRequest{WindowId: "okno-kart", Background: wskaznik(true)}, &druga)

	var przestrzen shared.BrowserWorkspaceSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserWorkspaceSave,
		shared.BrowserWorkspaceSaveRequest{WindowId: "okno-kart", Name: "Badanie rynku"}, &przestrzen)
	if przestrzen.Workspace.TabCount != 2 {
		t.Fatalf("przestrzeń zapamiętała %d kart, a w oknie były 2", przestrzen.Workspace.TabCount)
	}

	var zamkniecie shared.BrowserTabCloseResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserTabClose,
		shared.BrowserTabCloseRequest{TabId: druga.Tab.Id}, &zamkniecie)
	if !zamkniecie.Closed {
		t.Fatal("zamknięcie karty zameldowało brak skutku")
	}
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM karta_przegladania WHERE okno = ? AND zamknieta = 0`,
		"okno-kart"); liczba != 1 {
		t.Fatalf("po zamknięciu karty w oknie zostało %d otwartych, oczekiwano 1", liczba)
	}

	var wykaz shared.BrowserTabListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserTabList,
		shared.BrowserTabListRequest{WindowId: "okno-kart"}, &wykaz)
	if len(wykaz.Tabs) != 1 || wykaz.Tabs[0].Id != pierwsza.Tab.Id {
		t.Fatalf("wykaz kart nie zgadza się ze stanem bazy: %+v", wykaz.Tabs)
	}

	var otwarcie shared.BrowserWorkspaceOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserWorkspaceOpen,
		shared.BrowserWorkspaceOpenRequest{
			WorkspaceId: przestrzen.Workspace.Id, WindowId: "okno-kart",
		}, &otwarcie)
	if len(otwarcie.Tabs) != 2 {
		t.Fatalf("przywrócenie przestrzeni oddało %d kart, zapisano 2", len(otwarcie.Tabs))
	}
}

// TestPobranieZasobuNieStronowegoOdkladaBajty sprawdza drogę, którą powstają
// pobrania: wejście pod adres zasobu, którego nie da się pokazać jako strony.
// Miarą jest wiersz pobrania i plik w magazynie, nie treść odmowy.
func TestPobranieZasobuNieStronowegoOdkladaBajty(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuPrzegladarki(t, katalog)

	zawartosc := obrazPNG(t, 40, 30)
	serwer := serwerStrony(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(zawartosc)
	})

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: "okno-pobran", Url: serwer.URL + "/raport.png"})
	if !strings.Contains(odmowa.Message, "pobranie") {
		t.Fatalf("odmowa nie nazwała skutku, który zaszedł: %q", odmowa.Message)
	}

	var wykaz shared.BrowserDownloadListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserDownloadList,
		shared.BrowserDownloadListRequest{WindowId: wskaznik("okno-pobran")}, &wykaz)
	if len(wykaz.Downloads) != 1 {
		t.Fatalf("menedżer pobrań ma %d pozycji, oczekiwano 1", len(wykaz.Downloads))
	}
	pobranie := wykaz.Downloads[0]
	if pobranie.Status != shared.BrowserDownloadStatusCompleted || pobranie.TargetPath == nil {
		t.Fatalf("pobranie nie doszło do skutku: %+v", pobranie)
	}
	bajty := bajtyPodOdwolaniem(t, katalog, *pobranie.TargetPath)
	if !bytes.Equal(bajty, zawartosc) {
		t.Fatalf("w magazynie leży %d bajtów, serwer oddał %d", len(bajty), len(zawartosc))
	}
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM pobranie_przegladania WHERE okno = ? AND stan = 'completed'`,
		"okno-pobran"); liczba != 1 {
		t.Fatalf("baza nie odnotowała ukończonego pobrania (wierszy: %d)", liczba)
	}

	var sterowanie shared.BrowserDownloadControlResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserDownloadControl,
		shared.BrowserDownloadControlRequest{
			DownloadId: pobranie.Id, Action: shared.BrowserDownloadActionCancel,
		}, &sterowanie)
	if sterowanie.Download.Status != shared.BrowserDownloadStatusCancelled {
		t.Fatalf("przerwanie pobrania nie zmieniło jego stanu: %+v", sterowanie.Download)
	}
}

// TestWytworSesjiNiesieBajtyAZakladkaZostajeWBazie łączy dwie rodziny, których
// skutek mierzy się dwiema drogami: wytwór — bajtami w magazynie, zakładka —
// wierszem w tabeli.
func TestWytworSesjiNiesieBajtyAZakladkaZostajeWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuPrzegladarki(t, katalog)

	tresc := []byte("nazwa;cena\nProdukt;145\n")
	var wytwor shared.BrowserArtifactAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserArtifactAdd,
		shared.BrowserArtifactAddRequest{
			WindowId: "okno-materialu", Kind: shared.BrowserArtifactKindExtraction,
			Title: wskaznik("Tabela cen"), ContentBase64: wskaznik(wBase64(tresc)),
			MimeType: wskaznik("text/csv"),
		}, &wytwor)

	bajty := bajtyPodOdwolaniem(t, katalog, wytwor.Artifact.ContentRef)
	if !bytes.Equal(bajty, tresc) {
		t.Fatalf("wytwór wskazuje %d bajtów, wniesiono %d", len(bajty), len(tresc))
	}
	if wytwor.Artifact.SizeBytes == nil || *wytwor.Artifact.SizeBytes != int64(len(tresc)) {
		t.Fatalf("wytwór podaje rozmiar %v, a treść ma %d bajtów", wytwor.Artifact.SizeBytes, len(tresc))
	}

	var zakladka shared.BrowserBookmarkAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserBookmarkAdd,
		shared.BrowserBookmarkAddRequest{
			WindowId: "okno-materialu", Url: "https://przyklad.test/cennik",
			Title: wskaznik("Cennik"), Folder: wskaznik("Zakupy"),
			Tags: []string{"ceny", "dostawca"},
		}, &zakladka)
	if len(zakladka.Bookmark.Tags) != 2 {
		t.Fatalf("zakładka wróciła z %d etykietami, podano 2", len(zakladka.Bookmark.Tags))
	}
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM zakladka_przegladania WHERE okno = ? AND folder = 'Zakupy'`,
		"okno-materialu"); liczba != 1 {
		t.Fatalf("zakładki nie ma w bazie (wierszy: %d)", liczba)
	}

	var wykaz shared.BrowserBookmarkListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserBookmarkList,
		shared.BrowserBookmarkListRequest{
			WindowId: wskaznik("okno-materialu"), Query: wskaznik("Cennik"),
		}, &wykaz)
	if len(wykaz.Bookmarks) != 1 {
		t.Fatalf("wyszukanie zakładki oddało %d pozycji", len(wykaz.Bookmarks))
	}

	var zdjecie shared.BrowserBookmarkRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserBookmarkRemove,
		shared.BrowserBookmarkRemoveRequest{BookmarkId: zakladka.Bookmark.Id}, &zdjecie)
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM zakladka_przegladania WHERE identyfikator_zewnetrzny = ?`,
		zakladka.Bookmark.Id); liczba != 0 {
		t.Fatal("zakładka po usunięciu nadal leży w bazie")
	}
}

// TestKolejkaCzytaniaZestawyIWatkiZapisujaOznaczenia sprawdza to, co dotąd żyło
// wyłącznie w kliencie i ginęło przy przeładowaniu karty: klasyfikację notatki,
// jej wątek i przypięcie oraz przynależność źródła do zestawu.
func TestKolejkaCzytaniaZestawyIWatkiZapisujaOznaczenia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuPrzegladarki(t, katalog)
	const okno = "okno-porzadku"

	var zrodlo shared.BrowserSourceAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserSourceAdd,
		shared.BrowserSourceAddRequest{
			WindowId: okno, Url: "https://przyklad.test/oferta", Title: wskaznik("Oferta"),
		}, &zrodlo)

	var notatka shared.BrowserNoteAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNoteAdd,
		shared.BrowserNoteAddRequest{
			WindowId: okno, Content: "Ceny wzrosły o 8%", SourceId: wskaznik(zrodlo.Source.Id),
		}, &notatka)

	var zestaw shared.BrowserSourceGroupSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserSourceGroupSet,
		shared.BrowserSourceGroupSetRequest{
			WindowId: okno, Name: "Dostawcy", SourceIds: []string{zrodlo.Source.Id},
		}, &zestaw)
	if len(zestaw.Group.SourceIds) != 1 {
		t.Fatalf("zestaw wrócił ze składem %v, przypisano jedno źródło", zestaw.Group.SourceIds)
	}
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM zrodlo_przegladania WHERE identyfikator_zewnetrzny = ? AND grupa = ?`,
		zrodlo.Source.Id, zestaw.Group.Id); liczba != 1 {
		t.Fatal("przynależność źródła do zestawu nie została zapisana w bazie")
	}

	var watek shared.BrowserNoteThreadSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNoteThreadSet,
		shared.BrowserNoteThreadSetRequest{
			WindowId: okno, Name: "Wnioski", NoteIds: []string{notatka.Note.Id},
		}, &watek)

	var zmiana shared.BrowserNoteUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNoteUpdate,
		shared.BrowserNoteUpdateRequest{
			NoteId:         notatka.Note.Id,
			Classification: wskaznik(shared.BrowserNoteClassification(shared.BrowserNoteClassificationConclusion)),
			Pinned:         wskaznik(true),
		}, &zmiana)
	if zmiana.Note.Classification == nil || string(*zmiana.Note.Classification) != shared.BrowserNoteClassificationConclusion {
		t.Fatalf("klasyfikacja notatki nie została zapisana: %+v", zmiana.Note.Classification)
	}
	if zmiana.Note.Content != "Ceny wzrosły o 8%" {
		t.Fatalf("zmiana klasyfikacji nadpisała treść notatki: %q", zmiana.Note.Content)
	}
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM notatka_przegladania
		 WHERE identyfikator_zewnetrzny = ? AND klasyfikacja = 'conclusion' AND przypieta = 1 AND watek = ?`,
		notatka.Note.Id, watek.Thread.Id); liczba != 1 {
		t.Fatal("oznaczenia notatki nie leżą w bazie — po przeładowaniu okna zniknęłyby")
	}

	var kolejka shared.BrowserReadlistAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserReadlistAdd,
		shared.BrowserReadlistAddRequest{
			WindowId: okno, Url: "https://przyklad.test/artykul", Title: wskaznik("Artykuł"),
		}, &kolejka)
	var oznaczenie shared.BrowserReadlistRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserReadlistRemove,
		shared.BrowserReadlistRemoveRequest{ItemId: kolejka.Item.Id, MarkRead: wskaznik(true)}, &oznaczenie)
	if oznaczenie.Item == nil || !oznaczenie.Item.Read {
		t.Fatalf("pozycja kolejki nie została oznaczona jako przeczytana: %+v", oznaczenie.Item)
	}
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM pozycja_czytania_przegladania WHERE identyfikator_zewnetrzny = ? AND przeczytana = 1`,
		kolejka.Item.Id); liczba != 1 {
		t.Fatal("oznaczenie przeczytania nie zostało zapisane")
	}

	var usuniecie shared.BrowserSourceRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserSourceRemove,
		shared.BrowserSourceRemoveRequest{SourceId: zrodlo.Source.Id}, &usuniecie)
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM zrodlo_przegladania WHERE identyfikator_zewnetrzny = ?`,
		zrodlo.Source.Id); liczba != 0 {
		t.Fatal("źródło po usunięciu nadal leży w bazie")
	}
}

// TestMakroZapisujeKrokiAGraniceObowiazujaPoOdczycie sprawdza nagrywarkę makra
// (kroki mają być w bazie w kolejności zapisu) i granice Wykonawcy (odczyt ma
// oddać to, co ustawiono, a przy braku ustawienia — wartość domyślną rdzenia,
// nie odmowę).
func TestMakroZapisujeKrokiAGraniceObowiazujaPoOdczycie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuPrzegladarki(t, katalog)
	const okno = "okno-automatyzacji"

	var start shared.BrowserMacroRecordResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserMacroRecord,
		shared.BrowserMacroRecordRequest{
			WindowId: okno, Action: shared.BrowserMacroActionStart, Name: wskaznik("Zbiórka cenników"),
		}, &start)
	if !start.Recording {
		t.Fatal("nagrywarka nie zameldowała nagrywania po starcie")
	}

	for _, adres := range []string{"https://przyklad.test/a", "https://przyklad.test/b"} {
		var krok shared.BrowserMacroRecordResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserMacroRecord,
			shared.BrowserMacroRecordRequest{
				WindowId: okno, Action: shared.BrowserMacroActionStep, MacroId: wskaznik(start.Macro.Id),
				Step: &shared.AutomationStep{
					Id: adres, Kind: shared.AutomationStepKindCommand, Name: wskaznik(adres),
				},
			}, &krok)
	}

	var koniec shared.BrowserMacroRecordResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserMacroRecord,
		shared.BrowserMacroRecordRequest{
			WindowId: okno, Action: shared.BrowserMacroActionStop, MacroId: wskaznik(start.Macro.Id),
		}, &koniec)
	if koniec.Recording {
		t.Fatal("nagrywarka nadal melduje nagrywanie po jego zakończeniu")
	}
	if len(koniec.Macro.Steps) != 2 {
		t.Fatalf("makro wróciło z %d krokami, doklejono 2", len(koniec.Macro.Steps))
	}

	var kroki string
	if err := baza.QueryRow(`SELECT IFNULL(kroki_json,'') FROM makro_przegladania
		WHERE identyfikator_zewnetrzny = ?`, start.Macro.Id).Scan(&kroki); err != nil {
		t.Fatalf("nie można odczytać kroków makra z bazy: %v", err)
	}
	if !strings.Contains(kroki, "przyklad.test/a") || !strings.Contains(kroki, "przyklad.test/b") {
		t.Fatalf("kroki makra nie leżą w bazie: %q", kroki)
	}

	var domyslne shared.BrowserExecutorLimitsGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserExecutorLimitsGet,
		shared.BrowserExecutorLimitsGetRequest{WindowId: wskaznik(okno)}, &domyslne)
	if domyslne.Limits.MaxSteps <= 0 {
		t.Fatalf("odczyt granic bez ustawienia oddał %d kroków — pętla bez granicy", domyslne.Limits.MaxSteps)
	}

	var ustawienie shared.BrowserExecutorLimitsSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserExecutorLimitsSet,
		shared.BrowserExecutorLimitsSetRequest{
			WindowId: wskaznik(okno), MaxSteps: wskaznik(12),
			AllowedDomains: []string{"przyklad.test"},
		}, &ustawienie)

	var odczyt shared.BrowserExecutorLimitsGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserExecutorLimitsGet,
		shared.BrowserExecutorLimitsGetRequest{WindowId: wskaznik(okno)}, &odczyt)
	if odczyt.Limits.MaxSteps != 12 || len(odczyt.Limits.AllowedDomains) != 1 {
		t.Fatalf("granice po zapisie: %+v", odczyt.Limits)
	}
	if liczba := policzWierszePrzegladania(t, baza,
		`SELECT COUNT(*) FROM granica_wykonawcy_przegladania WHERE zasieg_id = ? AND max_krokow = 12`,
		okno); liczba != 1 {
		t.Fatal("granice Wykonawcy nie zostały zapisane w bazie")
	}
}

// TestNarzedziaInspekcyjneCzytajaStroneUruchomiona sprawdza silnik przeglądarki:
// drzewo DOM niesie element zbudowany skryptem, konsola niesie komunikat wypisany
// przez stronę, a rejestr sieciowy niesie żądanie dokumentu.
func TestNarzedziaInspekcyjneCzytajaStroneUruchomiona(t *testing.T) {
	if !chromiumStoi() {
		t.Skip("na tej maszynie nie ma Chromium — strony nie ma czym uruchomić")
	}
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	const okno = "okno-inspekcji"

	strona := `<html><head><title>Strona inspekcji</title></head><body>
		<div id="korzen"></div>
		<script>
			const d = document.createElement('p');
			d.id = 'zbudowany';
			d.textContent = 'Treść dopisana skryptem';
			document.getElementById('korzen').appendChild(d);
			console.log('komunikat ze strony inspekcji');
		</script></body></html>`
	serwer := serwerTresci(t, strona)

	var przejscie shared.BrowserNavigateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: okno, Url: serwer.URL}, &przejscie)

	var drzewo shared.BrowserDomInspectResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserDomInspect,
		shared.BrowserDomInspectRequest{
			WindowId: okno, Selector: wskaznik("#korzen"), IncludeStyles: wskaznik(true),
		}, &drzewo)
	znaleziony := false
	for _, wezel := range drzewo.Nodes {
		if wezel.Selector != nil && strings.Contains(*wezel.Selector, "zbudowany") {
			znaleziony = true
			if wezel.ComputedStyles == nil {
				t.Fatal("węzeł wrócił bez stylów, choć o nie poproszono")
			}
		}
	}
	if !znaleziony {
		t.Fatalf("drzewo DOM nie niesie elementu zbudowanego skryptem — to nie jest strona uruchomiona: %+v",
			drzewo.Nodes)
	}

	var konsola shared.BrowserConsoleReadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserConsoleRead,
		shared.BrowserConsoleReadRequest{WindowId: okno}, &konsola)
	slyszany := false
	for _, wpis := range konsola.Entries {
		if strings.Contains(wpis.Text, "komunikat ze strony inspekcji") {
			slyszany = true
		}
	}
	if !slyszany {
		t.Fatalf("konsola nie oddała komunikatu wypisanego przez stronę: %+v", konsola.Entries)
	}

	var siec shared.BrowserNetworkHarResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNetworkHar,
		shared.BrowserNetworkHarRequest{WindowId: okno}, &siec)
	if len(siec.Entries) == 0 || siec.HarRef == nil {
		t.Fatalf("rejestr sieciowy jest pusty albo bez zapisu HAR: %+v", siec)
	}
}

// TestAudytDostepnosciNazywaNaruszeniaWrazZWezlemDom sprawdza audyt strony
// z naruszeniami włożonymi celowo: miarą jest treść wykazu, gdzie naruszenie
// obrazka bez tekstu zastępczego niesie selektor wskazujący ten węzeł.
func TestAudytDostepnosciNazywaNaruszeniaWrazZWezlemDom(t *testing.T) {
	if !chromiumStoi() {
		t.Skip("na tej maszynie nie ma Chromium — strony nie ma czym uruchomić")
	}
	if !zewnetrzne.Stoi(narzedziePa11y) {
		t.Skip("na tej maszynie nie ma pa11y — audyt dostępności nie ma czym zbadać strony")
	}
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	const okno = "okno-audytu-dostepnosci"

	strona := `<html><head><title>Strona audytu</title></head><body>
		<img src="obraz.png">
		<input type="text" name="pole-bez-etykiety">
	</body></html>`
	serwer := serwerTresci(t, strona)

	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: okno, Url: serwer.URL}, nil)

	var audyt shared.BrowserAccessibilityAuditResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserAccessibilityAudit,
		shared.BrowserAccessibilityAuditRequest{WindowId: okno}, &audyt)

	if audyt.ErrorCount == 0 || len(audyt.Issues) == 0 {
		t.Fatalf("strona z naruszeniami włożonymi celowo wróciła bez naruszeń: %+v", audyt)
	}
	if audyt.Standard != shared.BrowserAccessibilityStandardWcag2aa {
		t.Fatalf("audyt bez wskazania normy miał iść normą wcag2aa, poszedł %q", audyt.Standard)
	}
	if !strings.HasPrefix(audyt.Url, serwer.URL) {
		t.Fatalf("audyt orzekł o adresie %q, a strona sprawdzianu stoi pod %q", audyt.Url, serwer.URL)
	}
	if audyt.ToolVersion == "" {
		t.Fatal("wynik bez wersji programu nie daje się porównać z wynikiem sprzed miesiąca")
	}
	obrazWskazany := false
	for _, zgloszenie := range audyt.Issues {
		if zgloszenie.Code == "" || zgloszenie.Message == "" {
			t.Fatalf("zgłoszenie bez kodu reguły albo bez treści: %+v", zgloszenie)
		}
		if zgloszenie.Selector != nil && strings.Contains(*zgloszenie.Selector, "img") {
			obrazWskazany = true
		}
	}
	if !obrazWskazany {
		t.Fatalf("żadne zgłoszenie nie wskazuje węzła obrazka bez tekstu zastępczego: %+v",
			audyt.Issues)
	}

	// Norma spoza wyliczenia kontraktu wraca odmową wskazania, nie komunikatem
	// o nieznanej normie.
	normaSpoza := shared.BrowserAccessibilityStandard("wcag9zz")
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserAccessibilityAudit,
		shared.BrowserAccessibilityAuditRequest{WindowId: okno, Standard: &normaSpoza})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("norma spoza kontraktu wróciła kodem %q", odmowa.Code)
	}
	if !strings.Contains(odmowa.Message, "nie należy do kontraktu") {
		t.Fatalf("odmowa nie nazywa normy spoza kontraktu: %q", odmowa.Message)
	}
}

// TestAudytDostepnosciStronyZgaszonejOdmawiaZamiastZeraNaruszen sprawdza, że audyt
// strony zgaszonej po przejściu odmawia zamiast zwrócić zero naruszeń — zero na
// stronie, której nie ma, byłoby brakiem pomiaru podanym jako pomiar.
func TestAudytDostepnosciStronyZgaszonejOdmawiaZamiastZeraNaruszen(t *testing.T) {
	if !chromiumStoi() {
		t.Skip("na tej maszynie nie ma Chromium — strony nie ma czym uruchomić")
	}
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	const okno = "okno-audytu-strony-zgaszonej"

	serwer := serwerTresci(t, stronaZTytulem("Strona gasnąca", "treść przed zgaśnięciem"))
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: okno, Url: serwer.URL}, nil)
	serwer.Close()

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserAccessibilityAudit,
		shared.BrowserAccessibilityAuditRequest{WindowId: okno})
	if !strings.Contains(odmowa.Message, "nie dała się pobrać") {
		t.Fatalf("odmowa nie mówi, że strony nie dało się pobrać: %q", odmowa.Message)
	}
}

// TestAudytDostepnosciBezProgramuOdmawiaNazywajacBrakIDrogeNaprawy zwęża ścieżkę
// wyszukiwania do przeglądarki i mierzy odmowę: ma nazwać brakujący program
// i pakiet naprawy, a nie usterkę wewnętrzną rdzenia.
func TestAudytDostepnosciBezProgramuOdmawiaNazywajacBrakIDrogeNaprawy(t *testing.T) {
	if !chromiumStoi() {
		t.Skip("na tej maszynie nie ma Chromium — strony nie ma czym uruchomić")
	}
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	const okno = "okno-audytu-bez-programu"

	serwer := serwerTresci(t, stronaZTytulem("Strona audytu bez programu", "treść"))
	wykonajUdana(t, zmontowany, zycie, shared.CommandBrowserNavigate,
		shared.BrowserNavigateRequest{WindowId: okno, Url: serwer.URL}, nil)

	// Ścieżka ma tylko przeglądarkę: audyt zatrzymuje się na braku programu,
	// nie przeglądarki.
	przegladarka := narzedzieChromium()
	sciezkaPrzegladarki, jest := zewnetrzne.Odnajdz(przegladarka)
	if !jest {
		t.Fatal("przeglądarka zniknęła między sprawdzeniem a pomiarem")
	}
	katalogSciezki := t.TempDir()
	if err := os.Symlink(sciezkaPrzegladarki,
		filepath.Join(katalogSciezki, przegladarka.Program)); err != nil {
		t.Fatalf("nie można wskazać przeglądarki w zwężonej ścieżce: %v", err)
	}
	t.Setenv("PATH", katalogSciezki)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandBrowserAccessibilityAudit,
		shared.BrowserAccessibilityAuditRequest{WindowId: okno})
	if odmowa.Code != shared.ErrorCodeChannelUnavailable {
		t.Fatalf("brak programu wrócił kodem %q", odmowa.Code)
	}
	if !strings.Contains(odmowa.Message, narzedziePa11y.Program) {
		t.Fatalf("odmowa nie nazywa brakującego programu: %q", odmowa.Message)
	}
	if !strings.Contains(odmowa.Message, narzedziePa11y.Pakiet) {
		t.Fatalf("odmowa nie podaje drogi naprawy: %q", odmowa.Message)
	}
}
