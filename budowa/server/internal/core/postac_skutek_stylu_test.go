package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// postacOknoSprawdzianuStylu jest oknem, w którego imieniu idą żądania tego
// pliku — wszystkie czynności i odczyty dzielą jedną tożsamość okna.
const postacOknoSprawdzianuStylu = "okno-postac-styl"

// postacTrescSprawdzianu to trzy akapity o różnej długości. Trzy, nie dwa: przy
// dwóch nie da się odróżnić „zmieniło wszystkie" od „zmieniło pierwszy
// i ostatni".
var postacTrescSprawdzianu = strings.Join([]string{
	"Rozdział pierwszy",
	"Treść rozdziału pierwszego, dość długa, żeby dała się zaznaczyć fragmentami.",
	"Rozdział drugi",
	"Treść rozdziału drugiego.",
	"Rozdział trzeci",
}, "\n")

// postacWierszStylu odczytuje wiersz arkusza stylów drugim połączeniem do
// bazy, niezależnym od rdzenia, i oddaje osobno postać znaku oraz postać
// akapitu.
func postacWierszStylu(t *testing.T, oboczne *sql.DB, dokument, nazwa string) (string, string) {
	t.Helper()

	var znak, akapit sql.NullString
	err := oboczne.QueryRow(`SELECT s.postac_znaku_json, s.postac_akapitu_json
	                         FROM styl_nazwany_studio s
	                         JOIN dokument_studio d ON d.id = s.dokument_id
	                         WHERE d.identyfikator_zewnetrzny = ? AND s.nazwa = ?`,
		dokument, nazwa).Scan(&znak, &akapit)
	if err != nil {
		t.Fatalf("styl %q nie ma wiersza w bazie: %v", nazwa, err)
	}
	return znak.String, akapit.String
}

// postacWierszDrzewa odczytuje drzewo postaci dokumentu drugim połączeniem do
// bazy i rozbiera zapisany zapis JSON z powrotem na strukturę kontraktu.
func postacWierszDrzewa(t *testing.T, oboczne *sql.DB, dokument string) shared.StudioDocumentForm {
	t.Helper()

	var drzewo string
	err := oboczne.QueryRow(`SELECT p.postac_json
	                         FROM postac_dokumentu_studio p
	                         JOIN dokument_studio d ON d.id = p.dokument_id
	                         WHERE d.identyfikator_zewnetrzny = ?`, dokument).Scan(&drzewo)
	if err != nil {
		t.Fatalf("dokument %q nie ma zapisanej postaci w bazie: %v", dokument, err)
	}
	var forma shared.StudioDocumentForm
	if err := json.Unmarshal([]byte(drzewo), &forma); err != nil {
		t.Fatalf("drzewo postaci w bazie jest nieczytelne: %v", err)
	}
	return forma
}

// postacZmianySledzoneWBazie odczytuje autorów zmian śledzonych własnym
// zapytaniem SQL, w kolejności wpisu, z pominięciem odpowiedzi komendy rdzenia.
func postacZmianySledzoneWBazie(t *testing.T, oboczne *sql.DB, dokument string) []string {
	t.Helper()

	wiersze, err := oboczne.Query(`SELECT z.autor
	                               FROM zmiana_sledzona_studio z
	                               JOIN dokument_studio d ON d.id = z.dokument_id
	                               WHERE d.identyfikator_zewnetrzny = ?
	                               ORDER BY z.id`, dokument)
	if err != nil {
		t.Fatalf("nie można odczytać zmian śledzonych: %v", err)
	}
	defer wiersze.Close()

	autorzy := []string{}
	for wiersze.Next() {
		var autor string
		if err := wiersze.Scan(&autor); err != nil {
			t.Fatalf("nieczytelny wiersz zmiany śledzonej: %v", err)
		}
		autorzy = append(autorzy, autor)
	}
	if err := wiersze.Err(); err != nil {
		t.Fatalf("przerwany odczyt zmian śledzonych: %v", err)
	}
	return autorzy
}

// postacZakresAkapitu ustala zakres akapitu w znakach treści sprawdzianu.
// Liczony z treści, nie wpisany liczbą: liczba wpisana ręcznie rozjedzie się
// przy pierwszej poprawce brzmienia akapitu.
func postacZakresAkapitu(t *testing.T, tresc string, numer int) (int, int) {
	t.Helper()

	polozenie := 0
	for i, wiersz := range strings.Split(tresc, "\n") {
		dlugosc := len([]rune(wiersz))
		if i == numer {
			return polozenie, polozenie + dlugosc
		}
		polozenie += dlugosc + 1
	}
	t.Fatalf("treść sprawdzianu nie ma akapitu numer %d", numer)
	return 0, 0
}

// postacStopienZnaku odczytuje stopień pisma osobnym wywołaniem odczytu
// postaci, a nie z odpowiedzi czynności, która ten stopień ustawiła.
func postacStopienZnaku(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	dokument string, od, do int) float64 {
	t.Helper()

	var odczyt shared.StudioFormatCharacterGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioFormatCharacterGet,
		shared.StudioFormatCharacterGetRequest{
			DocumentId: dokument, RangeStart: &od, RangeEnd: &do,
		}, &odczyt)
	if odczyt.Character.FontSizePt == nil {
		t.Fatalf("postać znaku zakresu %d-%d nie niesie stopnia pisma", od, do)
	}
	return *odczyt.Character.FontSizePt
}

// TestPostacStyluPrzestawiaWszystkieMiejscaUzycia wykazuje, że zmiana stopnia
// pisma w stylu nazwanym przestawia stopień we wszystkich trzech akapitach,
// które go używają, bez kopiowania postaci stylu na fragmenty.
func TestPostacStyluPrzestawiaWszystkieMiejscaUzycia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneStudia(t, katalog)
	dokument := dokumentZTrescia(t, zmontowany, zycie,
		postacOknoSprawdzianuStylu, postacTrescSprawdzianu).Id

	const nazwaStylu = "Nagłówek rozdziału pisma"
	var zalozony shared.StudioStyleSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioStyleSave,
		shared.StudioStyleSaveRequest{
			DocumentId: dokument,
			Name:       nazwaStylu,
			BasedOn:    wskaznik("Tekst zasadniczy"),
			Character:  json.RawMessage(`{"fontSizePt":14,"bold":true}`),
		}, &zalozony)

	// Trzy akapity nagłówkowe: pierwszy, trzeci i piąty wiersz treści.
	naglowki := []int{0, 2, 4}
	for _, numer := range naglowki {
		od, do := postacZakresAkapitu(t, postacTrescSprawdzianu, numer)
		wykonajUdana(t, zmontowany, zycie, shared.CommandStudioStyleApply,
			shared.StudioStyleApplyRequest{
				DocumentId: dokument, Name: nazwaStylu, RangeStart: &od, RangeEnd: &do,
			}, nil)
	}

	for _, numer := range naglowki {
		od, do := postacZakresAkapitu(t, postacTrescSprawdzianu, numer)
		if stopien := postacStopienZnaku(t, zmontowany, zycie, dokument, od, do); stopien != 14 {
			t.Fatalf("akapit %d przed zmianą stylu ma stopień %v, oczekiwano 14 ze stylu",
				numer, stopien)
		}
	}

	// Sedno sprawdzianu: JEDNA zmiana stylu.
	var zmieniony shared.StudioStyleSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioStyleSave,
		shared.StudioStyleSaveRequest{
			DocumentId: dokument,
			Name:       nazwaStylu,
			Character:  json.RawMessage(`{"fontSizePt":22}`),
		}, &zmieniony)

	for _, numer := range naglowki {
		od, do := postacZakresAkapitu(t, postacTrescSprawdzianu, numer)
		stopien := postacStopienZnaku(t, zmontowany, zycie, dokument, od, do)
		if stopien != 22 {
			t.Errorf("akapit %d po zmianie stylu ma stopień %v, oczekiwano 22 — zmiana "+
				"stylu nazwanego NIE przestawiła wszystkich miejsc użycia", numer, stopien)
		}
	}

	// Akapit bez tego stylu ma zostać nietknięty — „wszystkie miejsca użycia”
	// to nie „cały dokument”.
	odTresci, doTresci := postacZakresAkapitu(t, postacTrescSprawdzianu, 1)
	if stopien := postacStopienZnaku(t, zmontowany, zycie, dokument, odTresci, doTresci); stopien == 22 {
		t.Errorf("akapit treści zasadniczej dostał stopień 22, choć stylu %q nie używa",
			nazwaStylu)
	}

	// Miara niezależna pierwsza: wiersz stylu w bazie niesie nowy stopień.
	znak, _ := postacWierszStylu(t, oboczne, dokument, nazwaStylu)
	var postacStylu shared.StudioCharacterFormat
	if err := json.Unmarshal([]byte(znak), &postacStylu); err != nil {
		t.Fatalf("postać znaku stylu w bazie jest nieczytelna: %v", err)
	}
	if postacStylu.FontSizePt == nil || *postacStylu.FontSizePt != 22 {
		t.Errorf("wiersz stylu w bazie niesie stopień %v, oczekiwano 22", postacStylu.FontSizePt)
	}
	// Pogrubienie z zapisu pierwszego ma przeżyć zapis drugi, bo zmiana pola
	// nie podmienia całej postaci.
	if postacStylu.Bold == nil || !*postacStylu.Bold {
		t.Errorf("zmiana stopnia zdjęła pogrubienie stylu — zapis stylu podmienia " +
			"całą postać, zamiast wnosić do niej wskazane pola")
	}

	// Fragmenty w drzewie postaci nie niosą stopnia pisma — dowód, że styl
	// działa odwołaniem, a nie kopią.
	forma := postacWierszDrzewa(t, oboczne, dokument)
	for _, blok := range forma.Blocks {
		if blok.Paragraph == nil || blok.Paragraph.StyleName == nil ||
			*blok.Paragraph.StyleName != nazwaStylu {
			continue
		}
		for _, run := range blok.Runs {
			if run.Format != nil && run.Format.FontSizePt != nil {
				t.Errorf("fragment %q niesie stopień własny %v — stosowanie stylu "+
					"skopiowało jego postać na fragment, więc zmiana stylu nie miałaby "+
					"czego przestawić", run.Text, *run.Format.FontSizePt)
			}
		}
	}

	// Bilans ma powiedzieć, ilu miejsc zmiana dotknęła — cisza po zmianie stylu
	// jest tu szkodą osobną.
	if zmieniony.Balance.Applied < len(naglowki) {
		t.Errorf("bilans zmiany stylu mówi o %d miejscach, a stylu używa %d akapitów",
			zmieniony.Balance.Applied, len(naglowki))
	}
	var wykaz shared.StudioStyleListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioStyleList,
		shared.StudioStyleListRequest{DocumentId: dokument}, &wykaz)
	znaleziony := false
	for _, styl := range wykaz.Styles {
		if styl.Name != nazwaStylu {
			continue
		}
		znaleziony = true
		if styl.UsageCount == nil || *styl.UsageCount != len(naglowki) {
			t.Errorf("wykaz stylów mówi o %v miejscach użycia, oczekiwano %d",
				styl.UsageCount, len(naglowki))
		}
	}
	if !znaleziony {
		t.Fatalf("styl %q nie wrócił w wykazie stylów dokumentu", nazwaStylu)
	}
}

// TestPostacFragmentuZmieniaWylacznieFragment mierzy pracę na fragmentach.
//
// Pogrubienie idzie na SIEDEM znaków w środku akapitu. Wymaganie: pogrubione
// jest dokładnie tych siedem znaków — ani jeden przed, ani jeden po, ani jeden
// znak akapitu sąsiedniego.
func TestPostacFragmentuZmieniaWylacznieFragment(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneStudia(t, katalog)
	dokument := dokumentZTrescia(t, zmontowany, zycie,
		postacOknoSprawdzianuStylu, postacTrescSprawdzianu).Id

	poczatekAkapitu, _ := postacZakresAkapitu(t, postacTrescSprawdzianu, 1)
	od := poczatekAkapitu + 7
	do := od + 7

	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioFormatCharacterSet,
		shared.StudioFormatCharacterSetRequest{
			DocumentId: dokument, RangeStart: &od, RangeEnd: &do,
			Bold: wskaznik(true), Color: wskaznik("#B00020"),
		}, nil)

	// Miara pierwsza: drzewo w bazie. Zbieramy znaki pogrubione i porównujemy
	// ich zakres ze zamówionym.
	forma := postacWierszDrzewa(t, oboczne, dokument)
	polozenie := 0
	pierwszyPogrubiony, ostatniPogrubiony := -1, -1
	pogrubionych := 0
	for numer, blok := range forma.Blocks {
		if numer > 0 {
			polozenie++ // znak podziału wiersza między akapitami
		}
		for _, run := range blok.Runs {
			dlugosc := len([]rune(run.Text))
			pogrubiony := run.Format != nil && run.Format.Bold != nil && *run.Format.Bold
			if pogrubiony && dlugosc > 0 {
				if pierwszyPogrubiony < 0 {
					pierwszyPogrubiony = polozenie
				}
				ostatniPogrubiony = polozenie + dlugosc
				pogrubionych += dlugosc
			}
			polozenie += dlugosc
		}
	}
	if pogrubionych != do-od {
		t.Errorf("pogrubionych znaków %d, zamówiono %d — postać rozlała się poza "+
			"zaznaczenie albo go nie pokryła", pogrubionych, do-od)
	}
	if pierwszyPogrubiony != od || ostatniPogrubiony != do {
		t.Errorf("pogrubienie stoi na zakresie %d-%d, zamówiono %d-%d",
			pierwszyPogrubiony, ostatniPogrubiony, od, do)
	}

	// Znak przed zaznaczeniem i znak po nim nie mogą być pogrubione, mierzone
	// osobnym odczytem postaci.
	postacPogrubienieZakresu := func(poczatek, koniec int) *bool {
		t.Helper()
		var odczyt shared.StudioFormatCharacterGetResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandStudioFormatCharacterGet,
			shared.StudioFormatCharacterGetRequest{
				DocumentId: dokument, RangeStart: &poczatek, RangeEnd: &koniec,
			}, &odczyt)
		return odczyt.Character.Bold
	}
	if pogrubienie := postacPogrubienieZakresu(od, do); pogrubienie == nil || !*pogrubienie {
		t.Errorf("zaznaczony fragment nie jest pogrubiony: %v", pogrubienie)
	}
	if pogrubienie := postacPogrubienieZakresu(od-1, od); pogrubienie != nil && *pogrubienie {
		t.Errorf("znak PRZED zaznaczeniem został pogrubiony — postać rozlała się w lewo")
	}
	if pogrubienie := postacPogrubienieZakresu(do, do+1); pogrubienie != nil && *pogrubienie {
		t.Errorf("znak PO zaznaczeniu został pogrubiony — postać rozlała się w prawo")
	}

	// Akapit sąsiedni ma zostać nietknięty w całości.
	odSasiada, doSasiada := postacZakresAkapitu(t, postacTrescSprawdzianu, 2)
	if pogrubienie := postacPogrubienieZakresu(odSasiada, doSasiada); pogrubienie != nil && *pogrubienie {
		t.Errorf("akapit sąsiedni został pogrubiony — postać przeszła granicę akapitu")
	}

	// Treść nie ma prawa się zmienić: nałożenie postaci nie jest zmianą liter.
	if biezaca := trescDokumentu(t, zmontowany, zycie,
		postacOknoSprawdzianuStylu, dokument); biezaca != postacTrescSprawdzianu {

		t.Errorf("nałożenie postaci zmieniło treść dokumentu:\nbyło:  %q\njest:  %q",
			postacTrescSprawdzianu, biezaca)
	}
}

// TestPostacNarzedziaModeluOdkladaAutoraModel mierzy przełącznik „pokaż
// wszystko, co zrobił model”: ta sama czynność jako narzędzie modelu odkłada
// w bazie zmianę autora `model`, a przy śledzeniu Operatora — zmianę autora
// `uzytkownik`.
func TestPostacNarzedziaModeluOdkladaAutoraModel(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneStudia(t, katalog)
	dokument := dokumentZTrescia(t, zmontowany, zycie,
		postacOknoSprawdzianuStylu, postacTrescSprawdzianu).Id

	od, do := postacZakresAkapitu(t, postacTrescSprawdzianu, 1)
	autorModel := shared.StudioAuthor(shared.StudioAuthorModel)
	var czynnoscModelu shared.StudioFormatCharacterSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioFormatCharacterSet,
		shared.StudioFormatCharacterSetRequest{
			DocumentId: dokument, RangeStart: &od, RangeEnd: &do,
			Italic: wskaznik(true), Author: &autorModel,
		}, &czynnoscModelu)

	if czynnoscModelu.Change == nil {
		t.Fatalf("czynność modelu nie odłożyła zmiany śledzonej — nie ma czego podświetlić")
	}
	if czynnoscModelu.Change.Author != shared.StudioAuthorModel {
		t.Errorf("odpowiedź mówi o autorze %q, oczekiwano %q",
			czynnoscModelu.Change.Author, shared.StudioAuthorModel)
	}
	if czynnoscModelu.Change.Kind != shared.StudioChangeKindFormatowanie {
		t.Errorf("zmiana postaci ma rodzaj %q, oczekiwano %q",
			czynnoscModelu.Change.Kind, shared.StudioChangeKindFormatowanie)
	}

	// Miara niezależna: wiersz w bazie, własnym zapytaniem.
	autorzy := postacZmianySledzoneWBazie(t, oboczne, dokument)
	if len(autorzy) != 1 {
		t.Fatalf("baza niesie %d zmian śledzonych, oczekiwano jednej: %v", len(autorzy), autorzy)
	}
	if autorzy[0] != string(shared.StudioAuthorModel) {
		t.Fatalf("baza zapisała autora %q, a czynność wywołał model — przełącznik "+
			"zmian modelu nie odróżni jej od pracy Operatora", autorzy[0])
	}

	// Czynność Operatora przy śledzeniu włączonym odkłada się jako `uzytkownik`.
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTrackingSet,
		shared.StudioTrackingSetRequest{DocumentId: dokument, Enabled: true}, nil)

	odDrugi, doDrugi := postacZakresAkapitu(t, postacTrescSprawdzianu, 3)
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioFormatCharacterSet,
		shared.StudioFormatCharacterSetRequest{
			DocumentId: dokument, RangeStart: &odDrugi, RangeEnd: &doDrugi,
			Bold: wskaznik(true),
		}, nil)

	autorzy = postacZmianySledzoneWBazie(t, oboczne, dokument)
	if len(autorzy) != 2 {
		t.Fatalf("baza niesie %d zmian śledzonych, oczekiwano dwóch: %v", len(autorzy), autorzy)
	}
	if autorzy[1] != string(shared.StudioAuthorUzytkownik) {
		t.Errorf("czynność Operatora zapisała autora %q, oczekiwano %q — autor nie może "+
			"być zaszyty", autorzy[1], shared.StudioAuthorUzytkownik)
	}

	// Wykaz zmian śledzonych ma oddawać autora tą drogą, którą czyta go okno,
	// dla filtrowania.
	var wykaz shared.StudioTrackingListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTrackingList,
		shared.StudioTrackingListRequest{DocumentId: dokument}, &wykaz)
	zmianModelu := 0
	for _, zmiana := range wykaz.Changes {
		if zmiana.Author == shared.StudioAuthorModel {
			zmianModelu++
		}
	}
	if zmianModelu != 1 {
		t.Errorf("wykaz zmian śledzonych pokazuje %d zmian autora model, oczekiwano jednej",
			zmianModelu)
	}
}

// TestPostacZmianaNosnikaPrzeliczaUkladIOddajeBilans mierzy przeliczenie
// układu tabeli po zmianie nośnika: tabela na A3 poziomej po zejściu na A5
// pionową dostaje przeliczone szerokości kolumn i bilans nazywający tabelę,
// która się nie zmieściła.
func TestPostacZmianaNosnikaPrzeliczaUkladIOddajeBilans(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneStudia(t, katalog)
	dokument := dokumentZTrescia(t, zmontowany, zycie,
		postacOknoSprawdzianuStylu, postacTrescSprawdzianu).Id

	pozioma := shared.StudioPageOrientation(shared.StudioPageOrientationPozioma)
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPageSetupSet,
		shared.StudioPageSetupSetRequest{
			DocumentId: dokument, PaperName: wskaznik("A3"), Orientation: &pozioma,
		}, nil)

	var wstawiona shared.StudioTableInsertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTableInsert,
		shared.StudioTableInsertRequest{
			DocumentId: dokument, Offset: 0, Rows: 2, Columns: 3,
			WidthMm: postacWskaznikMiary(300),
		}, &wstawiona)
	tabela := wstawiona.Table.Id

	// Sedno: nośnik schodzi na mniejszy i pionowy.
	pionowa := shared.StudioPageOrientation(shared.StudioPageOrientationPionowa)
	var przelozony shared.StudioPageSetupSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPageSetupSet,
		shared.StudioPageSetupSetRequest{
			DocumentId: dokument, PaperName: wskaznik("A5"), Orientation: &pionowa,
		}, &przelozony)

	if przelozony.Balance.SkippedCount == 0 {
		t.Fatalf("zmiana nośnika z A3 poziomej na A5 pionową przeszła bez bilansu — " +
			"tabela na 300 mm nie mieści się na A5, a odpowiedź o tym milczy")
	}
	nazwana := false
	for _, pozycja := range przelozony.Balance.Skipped {
		if pozycja.Detail != nil && strings.Contains(*pozycja.Detail, tabela) {
			nazwana = true
		}
	}
	if !nazwana {
		t.Errorf("bilans nie nazywa tabeli %q, która się nie zmieściła: %+v",
			tabela, przelozony.Balance.Skipped)
	}

	// Miara niezależna pierwsza: nastawy strony odczytane osobnym wywołaniem.
	var nastawy shared.StudioPageSetupGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPageSetupGet,
		shared.StudioPageSetupGetRequest{DocumentId: dokument}, &nastawy)
	if nastawy.PageSetup.PageSize == nil || *nastawy.PageSetup.PageSize != "A5" {
		t.Fatalf("nastawy strony niosą nośnik %v, oczekiwano A5", nastawy.PageSetup.PageSize)
	}
	if nastawy.PageSetup.Orientation == nil ||
		*nastawy.PageSetup.Orientation != shared.StudioPageOrientationPionowa {
		t.Errorf("nastawy strony niosą orientację %v, oczekiwano pionowej",
			nastawy.PageSetup.Orientation)
	}
	// Pole `widthMm` niesie wymiar wyłącznie dla nośnika własnego, więc dla A5
	// wolno mu być puste.
	if nastawy.PageSetup.WidthMm != nil && *nastawy.PageSetup.WidthMm > 149 {
		t.Errorf("nastawy niosą szerokość %v mm, a A5 pionowa ma 148 mm",
			*nastawy.PageSetup.WidthMm)
	}
	// Wymiar nośnika nazwanego czyta się z jednego wykazu nośników rdzenia,
	// nie z osobnego wykazu Studia.
	var nosniki shared.StudioPagePaperListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPagePaperList,
		shared.StudioPagePaperListRequest{}, &nosniki)
	szerokoscA5 := 0.0
	for _, nosnik := range nosniki.Papers {
		if nosnik.Name == "A5" {
			szerokoscA5 = nosnik.WidthMm
		}
	}
	if szerokoscA5 <= 0 || szerokoscA5 > 149 {
		t.Fatalf("wykaz nośników podaje dla A5 szerokość %v mm, oczekiwano 148", szerokoscA5)
	}

	// Wykaz tabel niesie szerokości przeliczone: żadna nie jest zerowa ani
	// szersza od nowego nośnika.
	var wykaz shared.StudioTableListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTableList,
		shared.StudioTableListRequest{DocumentId: dokument, TableId: &tabela}, &wykaz)
	if len(wykaz.Tables) != 1 {
		t.Fatalf("wykaz tabel niesie %d pozycji, oczekiwano jednej", len(wykaz.Tables))
	}
	po := wykaz.Tables[0]
	if len(po.ColumnWidthsMm) != po.Columns {
		t.Fatalf("tabela ma %d kolumn, a szerokości %d", po.Columns, len(po.ColumnWidthsMm))
	}
	suma := 0.0
	for numer, szerokosc := range po.ColumnWidthsMm {
		if szerokosc <= 0 {
			t.Errorf("kolumna %d ma szerokość %v — zmiana nośnika obcięła układ, "+
				"zamiast go przeliczyć", numer, szerokosc)
		}
		suma += szerokosc
	}
	if po.WidthMm == nil {
		t.Fatalf("tabela po zmianie nośnika nie ma szerokości całości")
	}
	if *po.WidthMm > 149 {
		t.Errorf("tabela ma szerokość %v mm, a A5 pionowa ma 148 mm — układ nie został "+
			"przeliczony", *po.WidthMm)
	}
	if tabelaRoznicaMiar(suma, *po.WidthMm) > 0.5 {
		t.Errorf("suma szerokości kolumn %v nie zgadza się ze szerokością tabeli %v",
			suma, *po.WidthMm)
	}

	// Miara niezależna trzecia: nastawy strony leżą w bazie, nie w pamięci.
	var zapis sql.NullString
	err := oboczne.QueryRow(`SELECT p.nastawy_strony_json
	                         FROM postac_dokumentu_studio p
	                         JOIN dokument_studio d ON d.id = p.dokument_id
	                         WHERE d.identyfikator_zewnetrzny = ?`, dokument).Scan(&zapis)
	if err != nil {
		t.Fatalf("nastawy strony nie mają wiersza w bazie: %v", err)
	}
	if !strings.Contains(zapis.String, "A5") {
		t.Errorf("wiersz nastaw strony w bazie nie niesie A5: %s", zapis.String)
	}
}
