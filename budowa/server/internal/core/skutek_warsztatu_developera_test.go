package core

// Sprawdziany SKUTKU warsztatu modułu Developer — rodzin `developer.*`
// dobudowanych ponad Code Editor, Project Tree, Git Panel i Build Output.
//
// ── Czego te sprawdziany NIE robią ──────────────────────────────────────────
// Nie sprawdzają, czy odpowiedź jest odpowiedzią. Koperta ze stanem `ok`
// i pustym wynikiem jest kopertą udaną i zarazem kłamiącą — a to jest wzorzec
// szkody, który w tym produkcie już wystąpił. Dlatego każdy sprawdzian tego
// pliku po udanej odpowiedzi SCHODZI NIŻEJ i mierzy niezależnie: własnym
// zapytaniem SQL do tej samej bazy albo odczytem pliku z dysku.
//
// Zapytanie idzie osobnym połączeniem do pliku bazy, a nie przez repozytorium
// rdzenia. Gdyby szło przez repozytorium, sprawdzian mierzyłby to samo, czym
// mierzy się rdzeń — i wspólna usterka odczytu zostałaby niewidoczna po obu
// stronach.
//
// Uprząż jest własna, a nie wspólna z `uprzaz_skutku_test.go`: tamta montuje
// cały rdzeń i wchodzi kopertami protokołu, a moduł Developer potrzebuje OKNA
// z katalogiem roboczym, którego kontrakt komendy nie zakłada. Adapter składany
// wprost daje dokładnie tę jedną rzecz, a mierzone i tak jest to, co zostało
// w bazie i na dysku.

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/store"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// warsztatDevelopera niesie złożony moduł wraz z drogami pomiaru skutku.
type warsztatDevelopera struct {
	adapter    *adapterDevelopera
	okno       session.Okno
	repozytnik string
	sciezkaBaz string
}

// zlozWarsztatDevelopera składa adapter modułu nad świeżą bazą i świeżym oknem
// z katalogiem roboczym.
func zlozWarsztatDevelopera(t *testing.T) warsztatDevelopera {
	t.Helper()

	katalog := t.TempDir()
	sciezkaBazy := filepath.Join(katalog, "dane.sqlite")
	baza, err := store.Otworz(sciezkaBazy)
	if err != nil {
		t.Fatalf("nie można otworzyć bazy sprawdzianu: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })
	if err := baza.Migruj(); err != nil {
		t.Fatalf("migracje nie doszły do skutku: %v", err)
	}
	repozytoria, err := dane.Otworz(context.Background(), baza)
	if err != nil {
		t.Fatalf("nie można złożyć repozytoriów: %v", err)
	}

	roboczy := filepath.Join(katalog, "repozytorium")
	if err := os.MkdirAll(roboczy, 0o755); err != nil {
		t.Fatalf("nie można założyć katalogu roboczego: %v", err)
	}

	okna := session.NowyRejestr()
	sesja := okna.ZalozSesje("sprawdzian warsztatu", "")
	okno, err := okna.OtworzOkno(sesja.Id, session.Ustawienia{
		Modul:           "developer",
		KatalogiRobocze: []string{roboczy},
		TrybUprawnien:   shared.PermissionModeAcceptEdits,
	})
	if err != nil {
		t.Fatalf("nie można otworzyć okna modułu: %v", err)
	}

	adapter := nowyAdapterDevelopera(okna, nil).ZTrwaloscia(repozytoria.Developer)
	return warsztatDevelopera{
		adapter:    adapter,
		okno:       okno,
		repozytnik: roboczy,
		sciezkaBaz: sciezkaBazy,
	}
}

// zapiszPlikRoboczy zakłada plik w katalogu roboczym okna i oddaje jego ścieżkę.
func (w warsztatDevelopera) zapiszPlikRoboczy(t *testing.T, nazwa, tresc string) string {
	t.Helper()
	sciezka := filepath.Join(w.repozytnik, nazwa)
	if err := os.MkdirAll(filepath.Dir(sciezka), 0o755); err != nil {
		t.Fatalf("nie można założyć katalogu pliku %s: %v", nazwa, err)
	}
	if err := os.WriteFile(sciezka, []byte(tresc), 0o644); err != nil {
		t.Fatalf("nie można zapisać pliku %s: %v", nazwa, err)
	}
	return sciezka
}

// pomiarWBazieDevelopera otwiera bazę OSOBNYM połączeniem i wykonuje zapytanie
// sprawdzianu. To jest droga niezależna od repozytorium rdzenia — pomiar ma
// mierzyć bazę, a nie powtarzać odczyt, którym rdzeń właśnie zapisał.
func pomiarWBazieDevelopera(t *testing.T, sciezka, zapytanie string, argumenty ...any) *sql.Row {
	t.Helper()
	baza, err := sql.Open("sqlite", filepath.ToSlash(sciezka)+"?_pragma=busy_timeout(5000)")
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do pomiaru: %v", err)
	}
	t.Cleanup(func() { _ = baza.Close() })
	return baza.QueryRow(zapytanie, argumenty...)
}

// liczbaWBazieDevelopera oddaje jedną liczbę zmierzoną w bazie.
func liczbaWBazieDevelopera(t *testing.T, sciezka, zapytanie string, argumenty ...any) int64 {
	t.Helper()
	var liczba int64
	if err := pomiarWBazieDevelopera(t, sciezka, zapytanie, argumenty...).Scan(&liczba); err != nil {
		t.Fatalf("pomiar w bazie nie doszedł do skutku (%s): %v", zapytanie, err)
	}
	return liczba
}

// ── Wersje pliku ────────────────────────────────────────────────────────────

// TestSkutekPrzywroceniaWersjiPliku sprawdza, że przywrócenie wersji naprawdę
// zmienia PLIK NA DYSKU i że stan sprzed przywrócenia zostaje w bazie jako
// osobna migawka — bez niej powrót do wersji kasowałby pracę bezpowrotnie.
func TestSkutekPrzywroceniaWersjiPliku(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)
	sciezka := warsztat.zapiszPlikRoboczy(t, "kod.go", "wersja druga")

	// Migawkę stanu pierwszego zakłada się tą samą drogą, którą zakłada ją
	// zapis pliku z `createVersion`.
	kodWersji := nowyIdentyfikator(przedrostekWersjiPliku)
	if err := warsztat.adapter.repozytorium.ZapiszWersje(context.Background(), dane.WersjaPliku{
		Kod:     kodWersji,
		OknoKod: warsztat.okno.Id,
		Sciezka: sciezka,
		Tresc:   "wersja pierwsza",
		Rozmiar: int64(len("wersja pierwsza")),
	}); err != nil {
		t.Fatalf("nie można założyć migawki sprawdzianu: %v", err)
	}

	wykaz, err := warsztat.adapter.WykazWersjiPliku(context.Background(),
		shared.DeveloperFileVersionListRequest{WindowId: warsztat.okno.Id, Path: sciezka})
	if err != nil {
		t.Fatalf("wykaz wersji odmówił: %v", err)
	}
	if len(wykaz.Versions) != 1 || wykaz.Versions[0].Id != kodWersji {
		t.Fatalf("wykaz wersji nie oddał założonej migawki: %+v", wykaz.Versions)
	}

	if _, err := warsztat.adapter.PrzywrocWersjePliku(context.Background(),
		shared.DeveloperFileVersionRestoreRequest{
			WindowId:  warsztat.okno.Id,
			VersionId: kodWersji,
		}); err != nil {
		t.Fatalf("przywrócenie wersji odmówiło: %v", err)
	}

	// Pomiar pierwszy: plik na dysku.
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać przywróconego pliku: %v", err)
	}
	if string(bajty) != "wersja pierwsza" {
		t.Fatalf("plik na dysku nie niesie przywróconej treści; zastano %q", string(bajty))
	}

	// Pomiar drugi: własne zapytanie do bazy — migawka stanu sprzed
	// przywrócenia ma tam być, inaczej „wersja druga” przepadła na zawsze.
	ile := liczbaWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT COUNT(*) FROM developer_wersja_pliku WHERE okno_kod = ? AND tresc = ?`,
		warsztat.okno.Id, "wersja druga")
	if ile != 1 {
		t.Fatalf("stanu sprzed przywrócenia nie odłożono: w bazie %d migawek treści „wersja druga”", ile)
	}
}

// ── Punkty przerwania ───────────────────────────────────────────────────────

// TestSkutekPunktuPrzerwania sprawdza, że punkt postawiony w oknie zostaje
// w bazie i że jego zdjęcie naprawdę usuwa wiersz. Punkt widoczny w odpowiedzi,
// lecz nieobecny w bazie, zniknąłby przy kolejnym otwarciu okna.
func TestSkutekPunktuPrzerwania(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)
	sciezka := warsztat.zapiszPlikRoboczy(t, "usluga.go", "package usluga\n\nfunc Licz() int {\n\treturn 1\n}\n")

	warunek := "wynik > 3"
	odpowiedz, err := warsztat.adapter.UstawPunktPrzerwania(context.Background(),
		shared.DeveloperBreakpointSetRequest{
			WindowId:  warsztat.okno.Id,
			Path:      sciezka,
			Line:      4,
			Condition: &warunek,
		})
	if err != nil {
		t.Fatalf("ustawienie punktu przerwania odmówiło: %v", err)
	}
	if len(odpowiedz.Breakpoints) != 1 {
		t.Fatalf("odpowiedź nie niesie postawionego punktu: %+v", odpowiedz.Breakpoints)
	}
	// Rodzaj wywiedziony z warunku, a nie wpisany na siłę: punkt z warunkiem
	// jest punktem warunkowym niezależnie od tego, czy okno to nazwało.
	if odpowiedz.Breakpoints[0].Kind != shared.BreakpointKindConditional {
		t.Fatalf("punkt z warunkiem nie został rozpoznany jako warunkowy: %s",
			odpowiedz.Breakpoints[0].Kind)
	}

	var zapisanyWarunek string
	if err := pomiarWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT warunek FROM developer_punkt_przerwania
		 WHERE okno_kod = ? AND sciezka = ? AND wiersz = 4`,
		warsztat.okno.Id, sciezka).Scan(&zapisanyWarunek); err != nil {
		t.Fatalf("punktu przerwania nie ma w bazie: %v", err)
	}
	if zapisanyWarunek != warunek {
		t.Fatalf("baza niesie inny warunek punktu: %q", zapisanyWarunek)
	}

	zdejmij := true
	if _, err := warsztat.adapter.UstawPunktPrzerwania(context.Background(),
		shared.DeveloperBreakpointSetRequest{
			WindowId: warsztat.okno.Id,
			Path:     sciezka,
			Line:     4,
			Remove:   &zdejmij,
		}); err != nil {
		t.Fatalf("zdjęcie punktu przerwania odmówiło: %v", err)
	}
	ile := liczbaWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT COUNT(*) FROM developer_punkt_przerwania WHERE okno_kod = ?`, warsztat.okno.Id)
	if ile != 0 {
		t.Fatalf("po zdjęciu punktu w bazie zostało %d wierszy", ile)
	}
}

// ── Kolekcje zapytań i import OpenAPI ───────────────────────────────────────

// TestSkutekKolekcjiApi sprawdza, że zapisana kolekcja leży w bazie z treścią
// zapytań, a nie samą nazwą.
func TestSkutekKolekcjiApi(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)

	zapisana, err := warsztat.adapter.ZapiszKolekcjeApi(context.Background(),
		shared.DeveloperApiCollectionSaveRequest{
			WindowId: warsztat.okno.Id,
			Name:     "Usługa rozliczeń",
			Requests: json.RawMessage(`[{"name":"Lista faktur","method":"GET","url":"{{baseUrl}}/faktury"}]`),
			Environments: json.RawMessage(
				`{"przejsciowe":{"baseUrl":"https://przejsciowe.example"}}`),
		})
	if err != nil {
		t.Fatalf("zapis kolekcji odmówił: %v", err)
	}
	if zapisana.Collection.Id == "" {
		t.Fatalf("kolekcja wróciła bez identyfikatora")
	}

	var zapytania, srodowiska string
	if err := pomiarWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT zapytania, srodowiska FROM developer_kolekcja_api WHERE kod = ?`,
		zapisana.Collection.Id).Scan(&zapytania, &srodowiska); err != nil {
		t.Fatalf("kolekcji nie ma w bazie: %v", err)
	}
	if !strings.Contains(zapytania, "Lista faktur") {
		t.Fatalf("baza nie niesie treści zapytań kolekcji: %s", zapytania)
	}
	if !strings.Contains(srodowiska, "przejsciowe.example") {
		t.Fatalf("baza nie niesie środowisk kolekcji: %s", srodowiska)
	}

	wykaz, err := warsztat.adapter.WykazKolekcjiApi(context.Background(),
		shared.DeveloperApiCollectionListRequest{WindowId: warsztat.okno.Id})
	if err != nil {
		t.Fatalf("wykaz kolekcji odmówił: %v", err)
	}
	if len(wykaz.Collections) != 1 {
		t.Fatalf("wykaz kolekcji oddał %d pozycji zamiast jednej", len(wykaz.Collections))
	}
}

// TestSkutekImportuOpenapi sprawdza, że import kontraktu wytwarza ZAPYTANIA,
// a nie samą pustą kolekcję o ładnej nazwie.
func TestSkutekImportuOpenapi(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)
	warsztat.zapiszPlikRoboczy(t, "otwarte.yaml", `openapi: 3.0.0
info:
  title: Usługa magazynu
  version: "1.0"
servers:
  - url: https://magazyn.example/api
paths:
  /towary:
    get:
      operationId: wykazTowarow
      responses:
        "200":
          description: wykaz
    post:
      operationId: dodajTowar
      responses:
        "201":
          description: dodano
  /towary/{kod}:
    get:
      operationId: towar
      parameters:
        - name: kod
          in: path
          required: true
          schema:
            type: string
      responses:
        "200":
          description: towar
`)

	sciezka := "otwarte.yaml"
	odpowiedz, err := warsztat.adapter.ImportujOpenapi(context.Background(),
		shared.DeveloperApiOpenapiImportRequest{
			WindowId: warsztat.okno.Id,
			Path:     &sciezka,
		})
	if err != nil {
		t.Fatalf("import kontraktu odmówił: %v", err)
	}
	if odpowiedz.RequestCount != 3 {
		t.Fatalf("import oddał %d zapytań zamiast trzech opisanych w kontrakcie",
			odpowiedz.RequestCount)
	}
	if odpowiedz.Collection.Name != "Usługa magazynu" {
		t.Fatalf("kolekcja nie wzięła nazwy z tytułu kontraktu: %q", odpowiedz.Collection.Name)
	}

	var zapytania, srodowiska string
	if err := pomiarWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT zapytania, srodowiska FROM developer_kolekcja_api WHERE kod = ?`,
		odpowiedz.Collection.Id).Scan(&zapytania, &srodowiska); err != nil {
		t.Fatalf("kolekcji z importu nie ma w bazie: %v", err)
	}
	for _, oczekiwane := range []string{"wykazTowarow", "dodajTowar", "towar"} {
		if !strings.Contains(zapytania, oczekiwane) {
			t.Fatalf("kolekcja w bazie nie niesie zapytania %s: %s", oczekiwane, zapytania)
		}
	}
	if !strings.Contains(srodowiska, "https://magazyn.example/api") {
		t.Fatalf("serwer kontraktu nie wszedł do środowiska kolekcji: %s", srodowiska)
	}
}

// ── Konsola danych ──────────────────────────────────────────────────────────

// TestSkutekKonsoliDanych sprawdza całą drogę Data Console na bazie SQLite
// założonej w katalogu roboczym: opis połączenia trafia do bazy produktu, schemat
// czyta się z bazy Operatora, polecenie zmieniające naprawdę wstawia wiersz,
// a nastawa „tylko do odczytu” naprawdę odmawia.
func TestSkutekKonsoliDanych(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)

	// Baza Operatora zakładana jest wprost, żeby sprawdzian nie zależał od
	// żadnej innej komendy modułu.
	sciezkaOperatora := filepath.Join(warsztat.repozytnik, "magazyn.sqlite")
	operator, err := sql.Open("sqlite", filepath.ToSlash(sciezkaOperatora))
	if err != nil {
		t.Fatalf("nie można założyć bazy Operatora: %v", err)
	}
	defer operator.Close()
	if _, err := operator.Exec(
		`CREATE TABLE towar (id INTEGER PRIMARY KEY, nazwa TEXT NOT NULL, ilosc INTEGER)`); err != nil {
		t.Fatalf("nie można założyć tabeli sprawdzianu: %v", err)
	}

	nie := false
	polaczenie, err := warsztat.adapter.UstawPolaczenieDanych(context.Background(),
		shared.DeveloperDataConnectionSetRequest{
			WindowId: warsztat.okno.Id,
			Name:     "magazyn",
			Engine:   shared.DataEngineSqlite,
			Database: "magazyn.sqlite",
			ReadOnly: &nie,
		})
	if err != nil {
		t.Fatalf("opis połączenia odmówił: %v", err)
	}

	var baza string
	if err := pomiarWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT baza FROM developer_polaczenie_danych WHERE kod = ?`,
		polaczenie.Connection.Id).Scan(&baza); err != nil {
		t.Fatalf("połączenia nie ma w bazie produktu: %v", err)
	}
	if baza != sciezkaOperatora {
		t.Fatalf("baza zapisała ścieżkę %q zamiast rozstrzygniętej w obszarze okna %q",
			baza, sciezkaOperatora)
	}

	schemat, err := warsztat.adapter.SchematDanych(context.Background(),
		shared.DeveloperDataSchemaGetRequest{ConnectionId: polaczenie.Connection.Id})
	if err != nil {
		t.Fatalf("odczyt schematu odmówił: %v", err)
	}
	kolumny := map[string]bool{}
	for _, wezel := range schemat.Nodes {
		if wezel.Kind == shared.SchemaNodeKindColumn {
			kolumny[wezel.Name] = true
		}
	}
	for _, oczekiwana := range []string{"id", "nazwa", "ilosc"} {
		if !kolumny[oczekiwana] {
			t.Fatalf("schemat nie wykazał kolumny %s; zastano: %+v", oczekiwana, schemat.Nodes)
		}
	}

	if _, err := warsztat.adapter.WykonajZapytanieDanych(context.Background(),
		shared.DeveloperDataQueryRunRequest{
			ConnectionId: polaczenie.Connection.Id,
			Sql:          `INSERT INTO towar (nazwa, ilosc) VALUES ('młotek', 7)`,
		}); err != nil {
		t.Fatalf("polecenie wstawiające odmówiło: %v", err)
	}

	// Pomiar niezależny: wiersz ma być w bazie OPERATORA, a nie w odpowiedzi.
	var ilosc int64
	if err := operator.QueryRow(`SELECT ilosc FROM towar WHERE nazwa = 'młotek'`).
		Scan(&ilosc); err != nil {
		t.Fatalf("wstawionego wiersza nie ma w bazie Operatora: %v", err)
	}
	if ilosc != 7 {
		t.Fatalf("wiersz w bazie Operatora niesie ilość %d zamiast 7", ilosc)
	}

	odczyt, err := warsztat.adapter.WykonajZapytanieDanych(context.Background(),
		shared.DeveloperDataQueryRunRequest{
			ConnectionId: polaczenie.Connection.Id,
			Sql:          `SELECT nazwa, ilosc FROM towar`,
		})
	if err != nil {
		t.Fatalf("odczyt odmówił: %v", err)
	}
	if odczyt.Result.RowCount != 1 || len(odczyt.Result.Columns) != 2 {
		t.Fatalf("siatka wyników nie niesie wiersza: %+v", odczyt.Result)
	}
	if !strings.Contains(string(odczyt.Result.Rows), "młotek") {
		t.Fatalf("siatka wyników nie niesie wartości: %s", string(odczyt.Result.Rows))
	}

	// Brama „tylko do odczytu” ma zatrzymać zmianę PRZED silnikiem.
	tak := true
	if _, err := warsztat.adapter.UstawPolaczenieDanych(context.Background(),
		shared.DeveloperDataConnectionSetRequest{
			WindowId:     warsztat.okno.Id,
			ConnectionId: &polaczenie.Connection.Id,
			Name:         "magazyn",
			Engine:       shared.DataEngineSqlite,
			Database:     "magazyn.sqlite",
			ReadOnly:     &tak,
		}); err != nil {
		t.Fatalf("przestawienie połączenia odmówiło: %v", err)
	}
	if _, err := warsztat.adapter.WykonajZapytanieDanych(context.Background(),
		shared.DeveloperDataQueryRunRequest{
			ConnectionId: polaczenie.Connection.Id,
			Sql:          `DELETE FROM towar`,
		}); err == nil {
		t.Fatal("połączenie tylko do odczytu przepuściło polecenie kasujące")
	}
	var ile int64
	if err := operator.QueryRow(`SELECT COUNT(*) FROM towar`).Scan(&ile); err != nil {
		t.Fatalf("nie można zmierzyć bazy Operatora: %v", err)
	}
	if ile != 1 {
		t.Fatalf("mimo odmowy w bazie Operatora zostało %d wierszy zamiast jednego", ile)
	}
}

// TestSkutekMigracjiDanych sprawdza, że migracja naprawdę zmienia SCHEMAT bazy
// Operatora i że drugie uruchomienie jej nie powtarza.
func TestSkutekMigracjiDanych(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)

	sciezkaOperatora := filepath.Join(warsztat.repozytnik, "sklep.sqlite")
	operator, err := sql.Open("sqlite", filepath.ToSlash(sciezkaOperatora))
	if err != nil {
		t.Fatalf("nie można założyć bazy Operatora: %v", err)
	}
	defer operator.Close()
	if _, err := operator.Exec(`CREATE TABLE zaczep (id INTEGER PRIMARY KEY)`); err != nil {
		t.Fatalf("nie można założyć bazy Operatora: %v", err)
	}

	warsztat.zapiszPlikRoboczy(t, "migracje/001_klient.sql",
		"CREATE TABLE klient (id INTEGER PRIMARY KEY, nazwa TEXT NOT NULL);")
	warsztat.zapiszPlikRoboczy(t, "migracje/002_zamowienie.sql",
		"CREATE TABLE zamowienie (id INTEGER PRIMARY KEY, klient_id INTEGER NOT NULL);")

	nie := false
	polaczenie, err := warsztat.adapter.UstawPolaczenieDanych(context.Background(),
		shared.DeveloperDataConnectionSetRequest{
			WindowId: warsztat.okno.Id,
			Name:     "sklep",
			Engine:   shared.DataEngineSqlite,
			Database: "sklep.sqlite",
			ReadOnly: &nie,
		})
	if err != nil {
		t.Fatalf("opis połączenia odmówił: %v", err)
	}

	katalog := "migracje"
	pierwsze, err := warsztat.adapter.UruchomMigracjeDanych(context.Background(),
		shared.DeveloperDataMigrationRunRequest{
			ConnectionId: polaczenie.Connection.Id,
			WindowId:     warsztat.okno.Id,
			Direction:    "up",
			Target:       &katalog,
		})
	if err != nil {
		t.Fatalf("migracja odmówiła: %v", err)
	}
	if len(pierwsze.Applied) != 2 {
		t.Fatalf("migracja zameldowała %d zastosowanych kroków zamiast dwóch: %+v",
			len(pierwsze.Applied), pierwsze)
	}

	// Pomiar: tabele mają NAPRAWDĘ istnieć w bazie Operatora.
	for _, tabela := range []string{"klient", "zamowienie"} {
		var nazwa string
		if err := operator.QueryRow(
			`SELECT name FROM sqlite_master WHERE type='table' AND name = ?`, tabela).
			Scan(&nazwa); err != nil {
			t.Fatalf("migracja zameldowała skutek, a tabeli %s nie ma w bazie: %v", tabela, err)
		}
	}

	drugie, err := warsztat.adapter.UruchomMigracjeDanych(context.Background(),
		shared.DeveloperDataMigrationRunRequest{
			ConnectionId: polaczenie.Connection.Id,
			WindowId:     warsztat.okno.Id,
			Direction:    "up",
			Target:       &katalog,
		})
	if err != nil {
		t.Fatalf("powtórna migracja odmówiła: %v", err)
	}
	if len(drugie.Applied) != 0 {
		t.Fatalf("powtórna migracja zastosowała kroki po raz drugi: %+v", drugie.Applied)
	}
}

// ── Wynik testów i pokrycie ─────────────────────────────────────────────────

// TestSkutekWynikuTestowIPokrycia sprawdza, że rozbiór wyjścia przebiegu zapisuje
// wyniki i pokrycie do bazy, a komendy odczytu oddają to, co tam leży.
//
// Sprawdzian wchodzi drogą, którą wchodzi rdzeń: zakłada przebieg, dopisuje mu
// wiersze wyjścia tak, jak robi to pompa logu, i domyka pomiar. Mierzy potem
// bazę własnym zapytaniem — bo to ona jest jedynym miejscem, z którego wynik
// testu da się odczytać po zamknięciu okna.
func TestSkutekWynikuTestowIPokrycia(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)

	profil := warsztat.zapiszPlikRoboczy(t, "pokrycie.out", `mode: set
danacoconsole/liczydlo/suma.go:10.20,12.3 2 1
danacoconsole/liczydlo/suma.go:14.20,16.3 2 0
danacoconsole/liczydlo/roznica.go:8.20,9.10 1 1
`)

	przebieg := &przebiegBudowania{
		kod:       nowyIdentyfikator(przedrostekBudowania),
		oknoKod:   warsztat.okno.Id,
		zadanie:   "go test",
		argumenty: []string{"-coverprofile=" + filepath.Base(profil), "./..."},
		stan:      shared.BuildStatusRunning,
		koniec:    make(chan struct{}),
	}
	for _, wiersz := range []string{
		`{"Action":"run","Package":"danacoconsole/liczydlo","Test":"TestSuma"}`,
		`{"Action":"pass","Package":"danacoconsole/liczydlo","Test":"TestSuma","Elapsed":0.012}`,
		`{"Action":"output","Package":"danacoconsole/liczydlo","Test":"TestRoznica","Output":"    roznica_test.go:31: zastano 4, oczekiwano 2\n"}`,
		`{"Action":"fail","Package":"danacoconsole/liczydlo","Test":"TestRoznica","Elapsed":0.004}`,
		`{"Action":"skip","Package":"danacoconsole/liczydlo","Test":"TestIloraz","Elapsed":0}`,
		"kompilacja zakończona",
	} {
		przebieg.Dopisz(wiersz)
	}
	if err := warsztat.adapter.repozytorium.ZapiszPrzebieg(context.Background(),
		dane.PrzebiegBudowania{
			Kod:     przebieg.kod,
			OknoKod: warsztat.okno.Id,
			Zadanie: "go test",
			Stan:    shared.BuildStatusRunning,
		}); err != nil {
		t.Fatalf("nie można założyć przebiegu: %v", err)
	}
	warsztat.adapter.odlozPomiarPrzebiegu(przebieg)

	// Pomiar pierwszy: baza, własnym zapytaniem.
	nieprzeszly := liczbaWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT COUNT(*) FROM developer_wynik_testu WHERE budowanie_kod = ? AND stan = 'failed'`,
		przebieg.kod)
	if nieprzeszly != 1 {
		t.Fatalf("baza niesie %d testów nieudanych zamiast jednego", nieprzeszly)
	}
	wszystkie := liczbaWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT COUNT(*) FROM developer_wynik_testu WHERE budowanie_kod = ?`, przebieg.kod)
	if wszystkie != 3 {
		t.Fatalf("baza niesie %d wyników testów zamiast trzech", wszystkie)
	}

	// Pomiar drugi: odpowiedź komendy ma zgadzać się z bazą.
	wynik, err := warsztat.adapter.WynikTestow(context.Background(),
		shared.DeveloperTestResultGetRequest{BuildId: przebieg.kod})
	if err != nil {
		t.Fatalf("odczyt wyniku testów odmówił: %v", err)
	}
	if wynik.Passed != 1 || wynik.Failed != 1 || wynik.Skipped != 1 {
		t.Fatalf("podsumowanie testów nie zgadza się z wyjściem: %+v", wynik)
	}
	znalezione := false
	for _, pozycja := range wynik.Results {
		if pozycja.Name != "TestRoznica" {
			continue
		}
		znalezione = true
		if pozycja.Line == nil || *pozycja.Line != 31 {
			t.Fatalf("miejsce niepowodzenia nie zostało odczytane z treści: %+v", pozycja.Line)
		}
		if pozycja.Message == nil || !strings.Contains(*pozycja.Message, "oczekiwano 2") {
			t.Fatalf("treść niepowodzenia nie została zachowana: %+v", pozycja.Message)
		}
	}
	if !znalezione {
		t.Fatal("wykaz wyników nie zawiera testu, który nie przeszedł")
	}

	// Pokrycie: profil na dysku ma dać pomiar co do instrukcji i wiersza.
	pokrycie, err := warsztat.adapter.Pokrycie(context.Background(),
		shared.DeveloperCoverageGetRequest{BuildId: przebieg.kod})
	if err != nil {
		t.Fatalf("odczyt pokrycia odmówił: %v", err)
	}
	if len(pokrycie.Files) != 2 {
		t.Fatalf("pokrycie objęło %d plików zamiast dwóch z profilu: %+v",
			len(pokrycie.Files), pokrycie.Files)
	}
	// Instrukcji jest 5, pokrytych 3 — sześćdziesiąt procent.
	if pokrycie.Percent != 60 {
		t.Fatalf("pokrycie zbiorcze wyszło %d%% zamiast 60%%", pokrycie.Percent)
	}
	for _, plik := range pokrycie.Files {
		if !strings.HasSuffix(plik.Path, "suma.go") {
			continue
		}
		if len(plik.UncoveredLines) != 3 {
			t.Fatalf("wiersze bez pokrycia pliku suma.go: %+v", plik.UncoveredLines)
		}
	}
}

// TestSkutekWykazuIloguBudowan sprawdza, że historia przebiegów i log czytają
// dziennik, a nie oddają pustych wykazów.
func TestSkutekWykazuIloguBudowan(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)

	kod := nowyIdentyfikator(przedrostekBudowania)
	log := "krok pierwszy\nkrok drugi\nkrok trzeci"
	if err := warsztat.adapter.repozytorium.ZapiszPrzebieg(context.Background(),
		dane.PrzebiegBudowania{
			Kod:     kod,
			OknoKod: warsztat.okno.Id,
			Zadanie: "go build ./...",
			Stan:    shared.BuildStatusRunning,
		}); err != nil {
		t.Fatalf("nie można założyć przebiegu: %v", err)
	}
	if err := warsztat.adapter.repozytorium.ZakonczPrzebieg(context.Background(), kod,
		shared.BuildStatusSucceeded, nil, log); err != nil {
		t.Fatalf("nie można domknąć przebiegu: %v", err)
	}

	wykaz, err := warsztat.adapter.WykazBudowan(context.Background(),
		shared.DeveloperBuildListRequest{WindowId: warsztat.okno.Id})
	if err != nil {
		t.Fatalf("wykaz budowań odmówił: %v", err)
	}
	if len(wykaz.Builds) != 1 || wykaz.Builds[0].Id != kod {
		t.Fatalf("wykaz budowań nie oddał założonego przebiegu: %+v", wykaz.Builds)
	}
	if wykaz.Builds[0].Status != shared.BuildStatusSucceeded {
		t.Fatalf("wykaz oddał stan %s zamiast domkniętego", wykaz.Builds[0].Status)
	}

	ogon := 2
	wynik, err := warsztat.adapter.LogBudowania(context.Background(),
		shared.DeveloperBuildLogGetRequest{BuildId: kod, Tail: &ogon})
	if err != nil {
		t.Fatalf("odczyt logu odmówił: %v", err)
	}
	if len(wynik.Lines) != 2 || wynik.Lines[0] != "krok drugi" {
		t.Fatalf("ogon logu nie zgadza się z dziennikiem: %+v", wynik.Lines)
	}
	if !wynik.Truncated {
		t.Fatal("odpowiedź nie mówi o przycięciu, choć oddała mniej wierszy niż dziennik")
	}
}

// ── Zależności i skanowanie ─────────────────────────────────────────────────

// TestSkutekWykazuZaleznosci sprawdza, że wykaz powstaje z odczytu manifestu,
// a nie z pustej listy — i że rozróżnia zależność bezpośrednią od pośredniej.
func TestSkutekWykazuZaleznosci(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)
	warsztat.zapiszPlikRoboczy(t, "go.mod", `module przyklad

go 1.26

require (
	github.com/coder/websocket v1.8.15
	github.com/go-git/go-git/v5 v5.19.2
	golang.org/x/sys v0.47.0 // indirect
)
`)

	odpowiedz, err := warsztat.adapter.WykazZaleznosci(context.Background(),
		shared.DeveloperDependencyListRequest{WindowId: warsztat.okno.Id})
	if err != nil {
		t.Fatalf("wykaz zależności odmówił: %v", err)
	}
	if len(odpowiedz.Dependencies) != 3 {
		t.Fatalf("wykaz oddał %d zależności zamiast trzech z manifestu: %+v",
			len(odpowiedz.Dependencies), odpowiedz.Dependencies)
	}
	posrednie := 0
	for _, pozycja := range odpowiedz.Dependencies {
		if !pozycja.Direct {
			posrednie++
		}
	}
	if posrednie != 1 {
		t.Fatalf("rozpoznano %d zależności pośrednich zamiast jednej", posrednie)
	}
	if !strings.HasSuffix(odpowiedz.Manifest, "go.mod") {
		t.Fatalf("odpowiedź nie nazywa manifestu, z którego powstała: %q", odpowiedz.Manifest)
	}
}

// TestSkutekSkanuBezpieczenstwa sprawdza, że skan naprawdę znajduje sekret
// i wzorzec podatności w plikach repozytorium, zapisuje znaleziska do bazy
// i oddaje je wykazem — a treści sekretu do bazy NIE wpisuje.
func TestSkutekSkanuBezpieczenstwa(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)
	warsztat.zapiszPlikRoboczy(t, "konfiguracja.go", `package konfiguracja

// Klucz wpisany w kod — to ma znaleźć skan sekretów.
const Klucz = "AKIAIOSFODNN7EXAMPLE"
`)
	warsztat.zapiszPlikRoboczy(t, "zapytanie.go", `package zapytanie

func Szukaj(nazwa string) string {
	return "SELECT * FROM klient WHERE nazwa = '" + nazwa + "'"
}
`)

	skan, err := warsztat.adapter.UruchomSkan(context.Background(),
		shared.DeveloperScanRunRequest{
			WindowId: warsztat.okno.Id,
			Kinds:    []shared.ScanKind{shared.ScanKindSecrets, shared.ScanKindCode},
		})
	if err != nil {
		t.Fatalf("skan odmówił: %v", err)
	}
	if skan.Scan.Status != shared.BuildStatusSucceeded {
		t.Fatalf("skan domknął się stanem %s", skan.Scan.Status)
	}
	if skan.Scan.FindingCount == nil || *skan.Scan.FindingCount == 0 {
		t.Fatal("skan zameldował powodzenie i zero znalezisk tam, gdzie sekret i wzorzec stoją w plikach")
	}

	sekrety := liczbaWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT COUNT(*) FROM developer_znalezisko WHERE skan_kod = ? AND rodzaj = 'secrets'`,
		skan.Scan.Id)
	if sekrety == 0 {
		t.Fatal("skan sekretów nie zapisał ani jednego znaleziska")
	}
	kod := liczbaWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT COUNT(*) FROM developer_znalezisko WHERE skan_kod = ? AND rodzaj = 'code'`,
		skan.Scan.Id)
	if kod == 0 {
		t.Fatal("skan kodu nie zapisał ani jednego znaleziska")
	}

	// Sekret NIE ma prawa wylądować w bazie produktu: baza byłaby wtedy drugim
	// miejscem, w którym ten klucz leży.
	wyciek := liczbaWBazieDevelopera(t, warsztat.sciezkaBaz,
		`SELECT COUNT(*) FROM developer_znalezisko
		 WHERE skan_kod = ? AND (tytul LIKE '%AKIAIOSFODNN7EXAMPLE%'
		                      OR IFNULL(opis,'') LIKE '%AKIAIOSFODNN7EXAMPLE%')`,
		skan.Scan.Id)
	if wyciek != 0 {
		t.Fatal("treść znalezionego sekretu została zapisana do bazy produktu")
	}

	wykaz, err := warsztat.adapter.WykazZnalezisk(context.Background(),
		shared.DeveloperScanResultListRequest{ScanId: &skan.Scan.Id})
	if err != nil {
		t.Fatalf("wykaz znalezisk odmówił: %v", err)
	}
	if len(wykaz.Findings) != int(sekrety+kod) {
		t.Fatalf("wykaz oddał %d znalezisk, a w bazie leży %d",
			len(wykaz.Findings), sekrety+kod)
	}
	// Kolejność wagi: najcięższe u góry, bo tak czyta się wykaz.
	if wykaz.Findings[0].Severity != shared.ProblemSeverityError {
		t.Fatalf("wykaz nie zaczyna się od znaleziska najcięższego: %s",
			wykaz.Findings[0].Severity)
	}
	for _, znalezisko := range wykaz.Findings {
		if znalezisko.Path == nil || *znalezisko.Path == "" {
			t.Fatalf("znalezisko bez wskazania pliku nie prowadzi do niczego: %+v", znalezisko)
		}
	}
}

// ── Sonda programów warsztatu ───────────────────────────────────────────────

// TestSkutekSondyWarsztatu sprawdza, że sonda mierzy STAN SERWERA, a nie oddaje
// wykazu z wpisanymi na stałe odpowiedziami: program obecny ma mieć ścieżkę,
// a program wymyślony ma wrócić jako nieobecny.
func TestSkutekSondyWarsztatu(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)

	odpowiedz, err := warsztat.adapter.SprawdzWarsztat(context.Background(),
		shared.DeveloperToolchainCheckRequest{
			Programs: []string{"gofmt", "danaco-program-ktorego-nie-ma"},
		})
	if err != nil {
		t.Fatalf("sonda warsztatu odmówiła: %v", err)
	}
	if len(odpowiedz.Programs) != 2 {
		t.Fatalf("sonda oddała %d pozycji zamiast dwóch", len(odpowiedz.Programs))
	}
	for _, pozycja := range odpowiedz.Programs {
		switch pozycja.Program {
		case "gofmt":
			// gofmt jedzie z instalacją języka, więc na maszynie budującej
			// produkt stoi zawsze — a gdyby nie stał, sonda ma o tym powiedzieć.
			if pozycja.Present && (pozycja.Path == nil || *pozycja.Path == "") {
				t.Fatal("sonda uznała gofmt za obecny, lecz nie podała jego ścieżki")
			}
		case "danaco-program-ktorego-nie-ma":
			if pozycja.Present {
				t.Fatal("sonda uznała wymyślony program za obecny na serwerze")
			}
		default:
			t.Fatalf("sonda oddała pozycję, o którą nie pytano: %s", pozycja.Program)
		}
	}

	// Pytanie puste znaczy „powiedz o wszystkim, co znasz” — wykaz ma wtedy
	// objąć komplet programów warsztatu.
	pelna, err := warsztat.adapter.SprawdzWarsztat(context.Background(),
		shared.DeveloperToolchainCheckRequest{})
	if err != nil {
		t.Fatalf("sonda pełna odmówiła: %v", err)
	}
	if len(pelna.Programs) != len(narzedziaWarsztatuDevelopera()) {
		t.Fatalf("sonda pełna oddała %d pozycji, a warsztat zna %d programów",
			len(pelna.Programs), len(narzedziaWarsztatuDevelopera()))
	}
}

// ── Kontenery ───────────────────────────────────────────────────────────────

// TestSkutekWykazuKontenerow sprawdza rzecz, która w tym produkcie jest
// rozstrzygnięciem, a nie drobiazgiem: brak silnika kontenerów na serwerze ma
// wrócić JAWNIE polem `engineAvailable`, a nie udawać, że kontenerów nie ma.
//
// Sprawdzian nie zakłada, czy silnik na maszynie stoi — sprawdza spójność
// odpowiedzi z tym, co zastała. Wykaz niepusty przy `engineAvailable: false`
// byłby odpowiedzią wewnętrznie sprzeczną.
func TestSkutekWykazuKontenerow(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)

	odpowiedz, err := warsztat.adapter.WykazKontenerow(context.Background(),
		shared.DeveloperContainerListRequest{WindowId: warsztat.okno.Id})
	if err != nil {
		t.Fatalf("wykaz kontenerów odmówił zamiast powiedzieć o braku silnika: %v", err)
	}
	if odpowiedz.Containers == nil {
		t.Fatal("wykaz kontenerów oddał wykaz pusty jako brak wykazu — klient nie odróżni tego od usterki")
	}
	if !odpowiedz.EngineAvailable && len(odpowiedz.Containers) > 0 {
		t.Fatalf("odpowiedź mówi, że silnika nie ma, i zarazem wykazuje %d kontenerów",
			len(odpowiedz.Containers))
	}

	// Czynność na kontenerze przy braku silnika ma ODMÓWIĆ zdaniem nazywającym
	// brak — cisza kazałaby Operatorowi czekać na skutek, którego nie będzie.
	if !odpowiedz.EngineAvailable {
		_, err := warsztat.adapter.CzynnoscKontenera(context.Background(),
			shared.DeveloperContainerActionRequest{
				ContainerId: "dowolny",
				Action:      shared.ContainerActionKindStart,
			})
		if err == nil {
			t.Fatal("czynność na kontenerze zameldowała powodzenie bez silnika kontenerów")
		}
		if !strings.Contains(err.Error(), "silnik") {
			t.Fatalf("odmowa nie nazywa braku silnika: %v", err)
		}
	}
}

// ── Warstwa językowa ────────────────────────────────────────────────────────

// TestSkutekFormatowania sprawdza, że formatowanie z zapisem na dysk naprawdę
// zmienia PLIK, a nie tylko oddaje sformatowaną treść w odpowiedzi.
//
// Sprawdzian pomija się, gdy serwer nie ma formatera plików Go: brak programu
// jest wtedy stanem serwera, a nie usterką modułu — i mówi o tym sonda
// warsztatu, która ma własny sprawdzian.
func TestSkutekFormatowania(t *testing.T) {
	narzedzie, _, jest := formaterPliku("x.go")
	if !jest || !zewnetrzne.Stoi(narzedzie) {
		t.Skipf("serwer nie ma formatera plików Go (%s) — sprawdzian nie ma czym zmierzyć skutku",
			narzedzie.Program)
	}

	warsztat := zlozWarsztatDevelopera(t)
	// Formater jest programem serwera, więc adapter dostaje ten sam uruchamiacz
	// procesów, którym jedzie budowanie i git.
	warsztat.adapter.uruchamiacz = injection.UruchamiaczOkien()

	// Plik celowo źle wcięty i z niepotrzebnym odstępem — po sformatowaniu ma
	// wyglądać inaczej, więc pomiar ma co porównać.
	sciezka := warsztat.zapiszPlikRoboczy(t, "krzywy.go",
		"package krzywy\n\nfunc Suma( a int ,b int ) int {\n\t\t\treturn a+b\n}\n")
	przed, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać pliku sprawdzianu: %v", err)
	}

	zapisz := true
	odpowiedz, err := warsztat.adapter.Formatuj(context.Background(),
		shared.DeveloperFormatRunRequest{
			WindowId:    warsztat.okno.Id,
			Path:        sciezka,
			WriteToDisk: &zapisz,
		})
	if err != nil {
		t.Fatalf("formatowanie odmówiło: %v", err)
	}
	if !odpowiedz.Changed {
		t.Fatal("formatowanie zameldowało brak zmiany na pliku, który wymagał sformatowania")
	}

	// Pomiar: plik NA DYSKU, odczytany niezależnie od odpowiedzi.
	po, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać sformatowanego pliku: %v", err)
	}
	if string(po) == string(przed) {
		t.Fatal("odpowiedź mówi o zmianie, a plik na dysku został nietknięty")
	}
	if !strings.Contains(string(po), "func Suma(a int, b int) int {") {
		t.Fatalf("plik na dysku nie został sformatowany; zastano:\n%s", string(po))
	}
	if odpowiedz.File.Content == nil || *odpowiedz.File.Content != string(po) {
		t.Fatal("treść w odpowiedzi rozjeżdża się z treścią pliku na dysku")
	}

	// Plik roboczy formatowania nie ma prawa zostać w repozytorium Operatora.
	wpisy, err := os.ReadDir(warsztat.repozytnik)
	if err != nil {
		t.Fatalf("nie można przejrzeć katalogu roboczego: %v", err)
	}
	for _, wpis := range wpisy {
		if strings.Contains(wpis.Name(), "danaco-format") {
			t.Fatalf("po formatowaniu został plik roboczy %s", wpis.Name())
		}
	}
}

// ── Zapytanie HTTP ──────────────────────────────────────────────────────────

// TestSkutekZapytaniaApi sprawdza, że komenda naprawdę WYSYŁA zapytanie: mierzy
// je po stronie serwera, który je odebrał, a nie po treści odpowiedzi rdzenia.
// Sprawdzian obejmuje też podstawienie `{{nazwa}}` ze środowiska kolekcji —
// bez niego kolekcja miałaby adres wpisany na stałe.
func TestSkutekZapytaniaApi(t *testing.T) {
	warsztat := zlozWarsztatDevelopera(t)

	odebrane := make(chan string, 4)
	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tresc, _ := io.ReadAll(r.Body)
		odebrane <- r.Method + " " + r.URL.Path + " | " +
			r.Header.Get("X-Znacznik") + " | " + string(tresc)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"stan":"przyjete"}`))
	}))
	defer serwer.Close()

	// Środowisko kolekcji niesie adres serwera — zapytanie wskaże go
	// podstawieniem, tak jak robi to okno API Client.
	if _, err := warsztat.adapter.ZapiszKolekcjeApi(context.Background(),
		shared.DeveloperApiCollectionSaveRequest{
			WindowId:     warsztat.okno.Id,
			Name:         "sprawdzian",
			Requests:     json.RawMessage(`[]`),
			Environments: json.RawMessage(`{"lokalne":{"baseUrl":"` + serwer.URL + `"}}`),
		}); err != nil {
		t.Fatalf("zapis kolekcji odmówił: %v", err)
	}

	srodowisko := "lokalne"
	tresc := `{"nazwa":"młotek"}`
	rodzaj := "json"
	odpowiedz, err := warsztat.adapter.WykonajZapytanieApi(context.Background(),
		shared.DeveloperApiRequestRequest{
			WindowId:      warsztat.okno.Id,
			Method:        "post",
			Url:           "{{baseUrl}}/towary",
			Headers:       json.RawMessage(`{"X-Znacznik":"danaco"}`),
			Body:          &tresc,
			BodyKind:      &rodzaj,
			EnvironmentId: &srodowisko,
		})
	if err != nil {
		t.Fatalf("zapytanie odmówiło: %v", err)
	}

	// Pomiar niezależny: co naprawdę doszło do serwera.
	select {
	case slad := <-odebrane:
		if !strings.HasPrefix(slad, "POST /towary | danaco | ") {
			t.Fatalf("serwer odebrał co innego niż wysłano: %q", slad)
		}
		if !strings.Contains(slad, "młotek") {
			t.Fatalf("treść zapytania nie doszła do serwera: %q", slad)
		}
	default:
		t.Fatal("serwer nie odebrał ani jednego zapytania, choć komenda zameldowała odpowiedź")
	}

	if odpowiedz.Response.Status != http.StatusCreated {
		t.Fatalf("odpowiedź niesie stan %d zamiast 201", odpowiedz.Response.Status)
	}
	if odpowiedz.Response.Body == nil || !strings.Contains(*odpowiedz.Response.Body, "przyjete") {
		t.Fatalf("odpowiedź nie niesie treści serwera: %+v", odpowiedz.Response.Body)
	}
	if len(odpowiedz.Response.Headers) == 0 ||
		!strings.Contains(string(odpowiedz.Response.Headers), "application/json") {
		t.Fatalf("odpowiedź nie niesie nagłówków serwera: %s", string(odpowiedz.Response.Headers))
	}
}

// ── Debugowanie ─────────────────────────────────────────────────────────────

// TestSkutekSesjiDebugowania sprawdza całą drogę Run & Debug na PRAWDZIWYM
// adapterze Delve: sesja startuje, zatrzymuje się na postawionym punkcie,
// oddaje stos wywołań ze zmiennymi i liczy wyrażenie w kontekście ramki.
//
// To jest sprawdzian skutku, a nie koperty: odpowiedź `ok` z pustym stosem
// wywołań byłaby dokładnie tym wzorcem szkody, którego ten plik pilnuje. Pomiar
// idzie po WARTOŚCI zmiennej odczytanej z zatrzymanego procesu — takiej, której
// nie da się oddać bez faktycznego zatrzymania programu.
func TestSkutekSesjiDebugowania(t *testing.T) {
	if testing.Short() {
		t.Skip("sprawdzian buduje i uruchamia program pod debuggerem — pomijany w biegu skróconym")
	}
	if !zewnetrzne.Stoi(narzedzieDelve) {
		t.Skipf("serwer nie ma adaptera Delve — sprawdzian nie ma czym zmierzyć skutku")
	}

	warsztat := zlozWarsztatDevelopera(t)
	warsztat.adapter.uruchamiacz = injection.UruchamiaczOkien()

	warsztat.zapiszPlikRoboczy(t, "go.mod", "module sprawdzian\n\ngo 1.26\n")
	sciezka := warsztat.zapiszPlikRoboczy(t, "main.go", `package main

import "fmt"

func main() {
	suma := 0
	for i := 1; i <= 4; i++ {
		suma += i
	}
	fmt.Println(suma)
}
`)

	// Punkt stoi na wierszu z wypisaniem — w tej chwili pętla już policzyła.
	if _, err := warsztat.adapter.UstawPunktPrzerwania(context.Background(),
		shared.DeveloperBreakpointSetRequest{
			WindowId: warsztat.okno.Id,
			Path:     sciezka,
			Line:     10,
		}); err != nil {
		t.Fatalf("ustawienie punktu przerwania odmówiło: %v", err)
	}

	program := warsztat.repozytnik
	sesja, err := warsztat.adapter.UruchomDebugowanie(context.Background(),
		shared.DeveloperDebugSessionStartRequest{
			WindowId: warsztat.okno.Id,
			Program:  &program,
		})
	if err != nil {
		t.Fatalf("uruchomienie debugowania odmówiło: %v", err)
	}
	t.Cleanup(func() {
		_, _ = warsztat.adapter.SterujDebugowaniem(context.Background(),
			shared.DeveloperDebugSessionControlRequest{
				SessionId: sesja.Session.Id,
				Step:      shared.DebugStepKindStop,
			})
	})

	zatrzymana := czekajNaZatrzymanieSesji(t, warsztat.adapter, sesja.Session.Id)
	if zatrzymana.StoppedReason == nil {
		t.Fatal("sesja zatrzymała się bez podania powodu")
	}

	zakres, err := warsztat.adapter.ZakresDebugowania(context.Background(),
		shared.DeveloperDebugScopeGetRequest{SessionId: sesja.Session.Id})
	if err != nil {
		t.Fatalf("odczyt zakresu odmówił: %v", err)
	}
	if len(zakres.Frames) == 0 {
		t.Fatal("sesja zatrzymana oddała pusty stos wywołań")
	}
	if zakres.Frames[0].Name == "" || zakres.Frames[0].Line == nil {
		t.Fatalf("wierzchnia ramka stosu nie niesie miejsca: %+v", zakres.Frames[0])
	}
	if len(zakres.Scopes) == 0 {
		t.Fatal("sesja zatrzymana nie oddała ani jednego zakresu zmiennych")
	}

	// Pomiar właściwy: wartość policzona przez zatrzymany program.
	wynik, err := warsztat.adapter.ObliczWyrazenie(context.Background(),
		shared.DeveloperDebugEvaluateRequest{
			SessionId:  sesja.Session.Id,
			FrameId:    zakres.Frames[0].Id,
			Expression: "suma",
		})
	if err != nil {
		t.Fatalf("obliczenie wyrażenia odmówiło: %v", err)
	}
	if strings.TrimSpace(wynik.Value) != "10" {
		t.Fatalf("zatrzymany program niesie suma=%q zamiast 10 — sesja nie stoi tam, gdzie punkt",
			wynik.Value)
	}
}

// czekajNaZatrzymanieSesji czeka, aż program dojdzie do punktu przerwania.
//
// Czekanie jest odpytywaniem stanu, a nie uśpieniem na stałą chwilę: budowanie
// programu pod debuggerem trwa raz dłużej, raz krócej, a sprawdzian ma mierzyć
// skutek, nie szybkość maszyny.
func czekajNaZatrzymanieSesji(t *testing.T, adapter *adapterDevelopera,
	kod string) shared.DebugSession {
	t.Helper()

	granica := time.Now().Add(90 * time.Second)
	for time.Now().Before(granica) {
		sesja, jest := adapter.sesjeDebugowania.Sesja(kod)
		if !jest {
			t.Fatalf("sesja debugowania %s zniknęła z rejestru", kod)
		}
		migawka := sesja.Migawka()
		switch migawka.Status {
		case shared.DebugStatusStopped:
			return migawka
		case shared.DebugStatusTerminated:
			t.Fatal("program zakończył się, nie zatrzymawszy się na postawionym punkcie przerwania")
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("program nie zatrzymał się na punkcie przerwania w granicy czasu")
	return shared.DebugSession{}
}
