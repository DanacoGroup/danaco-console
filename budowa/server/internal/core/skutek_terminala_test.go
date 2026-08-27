package core

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	_ "modernc.org/sqlite"

	"danacoconsole/shared"
)

// Sprawdziany SKUTKU modułu Terminal.
//
// Różnica wobec sprawdzianu koperty jest tu istotą rzeczy. Koperta ze stanem
// `ok` i pustym wynikiem jest kopertą udaną i zarazem kłamiącą — a klient czyta
// kopertę, nie komentarz w kodzie. Każdy sprawdzian poniżej mierzy więc świat
// NIEZALEŻNIE od odpowiedzi rdzenia:
//
//   - wpis książki hostów, pozycja biblioteki i jej wersje — WŁASNYM zapytaniem
//     SQL do bazy, nie ponownym pytaniem tej samej komendy;
//   - klucz SSH — plikiem na dysku, jego prawami i tym, czy `x/crypto/ssh`
//     potrafi go odczytać;
//   - odczyt pliku — treścią, którą sprawdzian sam wcześniej zapisał;
//   - wstrzymanie procesu — stanem procesu w `/proc`, czyli u systemu, a nie
//     w polu odpowiedzi;
//   - obserwacja plików — PLIKIEM, który powstał, bo wyzwolone polecenie
//     naprawdę się wykonało;
//   - tunel — stanem końcowym procesu `ssh` odczytanym z bazy.

// oknoTerminalaSprawdzianu zakłada sesję i okno, w którym pracują karty
// sprawdzianu, i oddaje identyfikator okna wraz z jego katalogiem roboczym.
//
// Okno jest tu nieodzowne, nie ozdobne: z niego biorą się tryb uprawnień
// i obszar izolacji, a moduł odmawia każdej czynności oknu, którego nie zna.
func oknoTerminalaSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	katalogRoboczy string) string {

	t.Helper()

	var sesja shared.SessionCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandSessionCreate,
		shared.SessionCreateRequest{Title: wskaznikTekstu("sprawdzian terminala")}, &sesja)

	var okno shared.WindowCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWindowCreate, shared.WindowCreateRequest{
		SessionId:      sesja.Session.Id,
		ModuleId:       "terminal",
		ModelChannelId: "kanal-sprawdzianu",
		WorkingDirs:    []string{katalogRoboczy},
		ExecutionEnv:   shared.ExecutionEnvLocal,
		// Tryb `auto` przepuszcza uruchomienie procesu zleconego przez Operatora
		// i przez model — sprawdzian mierzy skutek czynności, nie bramę.
		PermissionMode: shared.PermissionModeAuto,
		WindowRole:     shared.WindowRoleStandalone,
	}, &okno)

	return okno.Window.Id
}

// kartaSprawdzianu otwiera kartę powłoki bash w podanym katalogu.
func kartaSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	oknoKod, katalog string) string {

	t.Helper()
	var karta shared.TerminalSessionOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalSessionOpen,
		shared.TerminalSessionOpenRequest{
			WindowId:   oknoKod,
			Shell:      shared.TerminalShellBash,
			WorkingDir: wskaznikTekstu(katalog),
		}, &karta)
	if karta.Session.Id == "" {
		t.Fatal("otwarcie karty oddało kartę bez identyfikatora")
	}
	return karta.Session.Id
}

// bazaSprawdzianuTerminala otwiera WŁASNE połączenie z plikiem bazy rdzenia.
//
// Własne, a nie uchwyt rdzenia — i to jest sedno sprawdzianu skutku: pomiar ma
// iść drogą, której mierzony kod nie kontroluje. Zapytanie zadane tym
// połączeniem widzi to, co naprawdę zostało zatwierdzone w pliku.
func bazaSprawdzianuTerminala(t *testing.T, katalog string) *sql.DB {
	t.Helper()
	baza, err := sql.Open("sqlite", filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do pomiaru: %v", err)
	}
	t.Cleanup(func() { _ = baza.Close() })
	return baza
}

// TestKsiazkaHostowZostajeWBazie sprawdza, że zapis wpisu naprawdę odkłada
// wiersz, a usunięcie naprawdę go zdejmuje. Pomiar idzie własnym zapytaniem do
// tabeli `terminal_host` — pytanie komendy `terminal.host.list` o skutek komendy
// `terminal.host.save` mierzyłoby zgodność rdzenia z samym sobą.
func TestKsiazkaHostowZostajeWBazie(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuTerminala(t, katalogDanych)

	port := 2222
	var zapisany shared.TerminalHostSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalHostSave,
		shared.TerminalHostSaveRequest{Host: shared.TerminalHost{
			Name:   "maszyna wydania",
			Target: "operator@danaco-system.example",
			Port:   &port,
			Group:  wskaznikTekstu("wydanie"),
			Note:   wskaznikTekstu("wpis sprawdzianu"),
		}}, &zapisany)

	if !zapisany.Created {
		t.Error("pierwszy zapis wpisu nie zameldował powstania wpisu")
	}
	var nazwa, cel string
	var zapisanyPort int
	wiersz := baza.QueryRow(
		`SELECT nazwa, cel, port FROM terminal_host WHERE kod = ?`, zapisany.Host.Id)
	if err := wiersz.Scan(&nazwa, &cel, &zapisanyPort); err != nil {
		t.Fatalf("wpisu hosta nie ma w bazie po udanym zapisie: %v", err)
	}
	if nazwa != "maszyna wydania" || cel != "operator@danaco-system.example" || zapisanyPort != 2222 {
		t.Errorf("wiersz hosta niesie inną treść niż zapisana: nazwa=%q cel=%q port=%d",
			nazwa, cel, zapisanyPort)
	}

	// Zapis pod istniejącym identyfikatorem ma PODMIENIĆ wpis, a nie założyć drugi.
	var podmieniony shared.TerminalHostSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalHostSave,
		shared.TerminalHostSaveRequest{Host: shared.TerminalHost{
			Id:     zapisany.Host.Id,
			Name:   "maszyna wydania (po zmianie)",
			Target: "operator@danaco-system.example",
		}}, &podmieniony)
	if podmieniony.Created {
		t.Error("podmiana istniejącego wpisu zameldowała powstanie nowego")
	}
	var ile int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM terminal_host`).Scan(&ile); err != nil {
		t.Fatalf("nie można policzyć wpisów książki hostów: %v", err)
	}
	if ile != 1 {
		t.Errorf("po podmianie wpisu w książce hostów leży %d wierszy zamiast jednego", ile)
	}

	var usuniety shared.TerminalHostRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalHostRemove,
		shared.TerminalHostRemoveRequest{HostId: zapisany.Host.Id}, &usuniety)
	if !usuniety.Removed {
		t.Error("usunięcie istniejącego wpisu zameldowało, że wpisu nie było")
	}
	if err := baza.QueryRow(`SELECT COUNT(*) FROM terminal_host`).Scan(&ile); err != nil {
		t.Fatalf("nie można policzyć wpisów książki hostów: %v", err)
	}
	if ile != 0 {
		t.Errorf("po usunięciu wpisu w książce hostów został %d wiersz", ile)
	}
}

// TestBibliotekaSkryptowTrzymaWersje sprawdza, że drugi zapis tej samej pozycji
// zakłada KOLEJNĄ WERSJĘ, a nie nadpisuje poprzedniej treści, oraz że usunięcie
// pozycji zabiera wszystkie jej wersje.
func TestBibliotekaSkryptowTrzymaWersje(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuTerminala(t, katalogDanych)

	var pierwsza shared.TerminalScriptSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalScriptSave,
		shared.TerminalScriptSaveRequest{Script: shared.TerminalScript{
			Name:    "wydanie",
			Kind:    shared.TerminalScriptKindScript,
			Shell:   shared.TerminalShellBash,
			Content: "echo pierwsza",
		}}, &pierwsza)
	if !pierwsza.Created || pierwsza.Script.Version != 1 {
		t.Fatalf("pierwszy zapis oddał created=%v wersja=%d, oczekiwano true/1",
			pierwsza.Created, pierwsza.Script.Version)
	}

	var druga shared.TerminalScriptSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalScriptSave,
		shared.TerminalScriptSaveRequest{Script: shared.TerminalScript{
			Id:      pierwsza.Script.Id,
			Name:    "wydanie",
			Kind:    shared.TerminalScriptKindScript,
			Shell:   shared.TerminalShellBash,
			Content: "echo druga",
		}}, &druga)
	if druga.Created || druga.Script.Version != 2 {
		t.Fatalf("drugi zapis oddał created=%v wersja=%d, oczekiwano false/2",
			druga.Created, druga.Script.Version)
	}

	// Pomiar niezależny: obie wersje mają leżeć w tabeli wersji, każda ze swoją
	// treścią. Sam numer w odpowiedzi nie dowodzi, że poprzednia treść przetrwała.
	wiersze, err := baza.Query(
		`SELECT wersja, tresc FROM terminal_skrypt_wersja WHERE skrypt_kod = ? ORDER BY wersja`,
		pierwsza.Script.Id)
	if err != nil {
		t.Fatalf("nie można odczytać wersji pozycji biblioteki: %v", err)
	}
	defer wiersze.Close()
	tresci := map[int64]string{}
	for wiersze.Next() {
		var wersja int64
		var tresc string
		if err := wiersze.Scan(&wersja, &tresc); err != nil {
			t.Fatalf("nieczytelny wiersz wersji: %v", err)
		}
		tresci[wersja] = tresc
	}
	if tresci[1] != "echo pierwsza" || tresci[2] != "echo druga" {
		t.Errorf("wersje pozycji nie niosą obu brzmień treści: %#v", tresci)
	}

	var usunieta shared.TerminalScriptRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalScriptRemove,
		shared.TerminalScriptRemoveRequest{ScriptId: pierwsza.Script.Id}, &usunieta)
	if !usunieta.Removed {
		t.Error("usunięcie istniejącej pozycji zameldowało, że pozycji nie było")
	}
	var wersjiPoUsunieciu int
	if err := baza.QueryRow(
		`SELECT COUNT(*) FROM terminal_skrypt_wersja WHERE skrypt_kod = ?`,
		pierwsza.Script.Id).Scan(&wersjiPoUsunieciu); err != nil {
		t.Fatalf("nie można policzyć wersji po usunięciu: %v", err)
	}
	if wersjiPoUsunieciu != 0 {
		t.Errorf("po usunięciu pozycji zostało %d jej wersji", wersjiPoUsunieciu)
	}
}

// TestKluczSSHPowstajeNaDysku sprawdza, że wytworzenie klucza kładzie na dysku
// PRAWDZIWY klucz: plik daje się odczytać biblioteką SSH, jego odcisk zgadza się
// z odciskiem w odpowiedzi, a prawa pliku są takie, że `ssh` go przyjmie.
func TestKluczSSHPowstajeNaDysku(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuTerminala(t, katalogDanych)

	var wytworzony shared.TerminalKeyGenerateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalKeyGenerate,
		shared.TerminalKeyGenerateRequest{
			Name:    "klucz wydania",
			KeyType: shared.TerminalKeyTypeEd25519,
			Comment: wskaznikTekstu("danaco@sprawdzian"),
		}, &wytworzony)

	if wytworzony.Key.Path == nil || *wytworzony.Key.Path == "" {
		t.Fatal("wytworzony klucz nie ma ścieżki, więc nie da się go użyć")
	}
	sciezka := *wytworzony.Key.Path
	opis, err := os.Stat(sciezka)
	if err != nil {
		t.Fatalf("pliku klucza nie ma na dysku po udanym wytworzeniu: %v", err)
	}
	if prawa := opis.Mode().Perm(); prawa != 0o600 {
		t.Errorf("klucz prywatny ma prawa %o; `ssh` odmawia użycia klucza szerzej otwartego niż 0600", prawa)
	}
	tresc, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać pliku klucza: %v", err)
	}
	prywatny, err := ssh.ParseRawPrivateKey(tresc)
	if err != nil {
		t.Fatalf("plik zapisany jako klucz prywatny nie jest kluczem prywatnym: %v", err)
	}
	podpisujacy, err := ssh.NewSignerFromKey(prywatny)
	if err != nil {
		t.Fatalf("z zapisanego klucza nie da się złożyć podpisującego: %v", err)
	}
	if odcisk := ssh.FingerprintSHA256(podpisujacy.PublicKey()); odcisk != wytworzony.Key.Fingerprint {
		t.Errorf("odcisk w odpowiedzi (%s) nie opisuje klucza leżącego na dysku (%s)",
			wytworzony.Key.Fingerprint, odcisk)
	}

	// Wpis książki hostów wskazujący ten klucz ma po zdjęciu klucza wrócić do
	// klucza domyślnego — i ma zostać wymieniony w odpowiedzi.
	var host shared.TerminalHostSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalHostSave,
		shared.TerminalHostSaveRequest{Host: shared.TerminalHost{
			Name: "maszyna z kluczem", Target: "operator@example", KeyId: &wytworzony.Key.Id,
		}}, &host)

	prawda := true
	var zdjety shared.TerminalKeyRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalKeyRemove,
		shared.TerminalKeyRemoveRequest{KeyId: wytworzony.Key.Id, DeleteFiles: &prawda}, &zdjety)

	if !zdjety.Removed {
		t.Error("zdjęcie istniejącego klucza zameldowało, że klucza nie było")
	}
	if len(zdjety.DetachedHostIds) != 1 || zdjety.DetachedHostIds[0] != host.Host.Id {
		t.Errorf("odpowiedź nie wymienia wpisu, który stracił wskazanie klucza: %#v",
			zdjety.DetachedHostIds)
	}
	var wskazanie *string
	if err := baza.QueryRow(
		`SELECT klucz_kod FROM terminal_host WHERE kod = ?`, host.Host.Id).Scan(&wskazanie); err != nil {
		t.Fatalf("nie można odczytać wpisu hosta po zdjęciu klucza: %v", err)
	}
	if wskazanie != nil {
		t.Errorf("wpis hosta wciąż wskazuje zdjęty klucz (%q)", *wskazanie)
	}
	if _, err := os.Stat(sciezka); !os.IsNotExist(err) {
		t.Error("plik klucza został na dysku mimo prośby o jego usunięcie")
	}
}

// TestOdczytPlikuOddajeTrescZDysku sprawdza, że `terminal.file.read` czyta
// PLIK — treść, którą sprawdzian sam zapisał — i że ogon przycina go od końca.
func TestOdczytPlikuOddajeTrescZDysku(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	roboczy := t.TempDir()
	oknoKod := oknoTerminalaSprawdzianu(t, zmontowany, zycie, roboczy)
	kartaKod := kartaSprawdzianu(t, zmontowany, zycie, oknoKod, roboczy)

	tresc := "pierwszy\ndrugi\ntrzeci\n"
	if err := os.WriteFile(filepath.Join(roboczy, "manifest.txt"), []byte(tresc), 0o644); err != nil {
		t.Fatalf("nie można przygotować pliku sprawdzianu: %v", err)
	}

	var odczyt shared.TerminalFileReadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalFileRead,
		shared.TerminalFileReadRequest{SessionId: kartaKod, Path: "manifest.txt"}, &odczyt)

	if odczyt.Content != tresc {
		t.Errorf("odczyt oddał treść %q, a na dysku leży %q", odczyt.Content, tresc)
	}
	if odczyt.SizeBytes == nil || *odczyt.SizeBytes != int64(len(tresc)) {
		t.Errorf("odczyt nie podał rozmiaru pliku albo podał zły: %#v", odczyt.SizeBytes)
	}

	ogon := 2
	var przyciety shared.TerminalFileReadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalFileRead,
		shared.TerminalFileReadRequest{SessionId: kartaKod, Path: "manifest.txt", Tail: &ogon}, &przyciety)
	if strings.Contains(przyciety.Content, "pierwszy") {
		t.Errorf("ogon dwóch wierszy oddał także wiersz pierwszy: %q", przyciety.Content)
	}
	if !przyciety.Truncated {
		t.Error("odczyt przycięty ogonem nie zameldował przycięcia")
	}

	// Plik, którego nie ma, ma dać odmowę nazwaną, a nie pustą treść: pusty
	// napis znaczyłby „plik jest pusty”, czyli nieprawdę.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandTerminalFileRead,
		shared.TerminalFileReadRequest{SessionId: kartaKod, Path: "nie-ma-takiego.txt"})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("odczyt nieistniejącego pliku odmówił kodem %s zamiast not_found", odmowa.Code)
	}
}

// TestZamkniecieKartyZdejmujeJaZWykazu sprawdza, że zamknięcie karty zostaje
// w bazie i że wykaz kart przestaje ją pokazywać, dopóki się o nią wprost nie
// zapyta.
func TestZamkniecieKartyZdejmujeJaZWykazu(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuTerminala(t, katalogDanych)
	roboczy := t.TempDir()
	oknoKod := oknoTerminalaSprawdzianu(t, zmontowany, zycie, roboczy)
	kartaKod := kartaSprawdzianu(t, zmontowany, zycie, oknoKod, roboczy)

	var wykaz shared.TerminalSessionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalSessionList,
		shared.TerminalSessionListRequest{WindowId: &oknoKod}, &wykaz)
	if wykaz.Total != 1 || wykaz.Sessions[0].Id != kartaKod {
		t.Fatalf("wykaz kart nie pokazuje karty otwartej przed chwilą: %#v", wykaz)
	}

	var zamknieta shared.TerminalSessionCloseResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalSessionClose,
		shared.TerminalSessionCloseRequest{SessionId: kartaKod}, &zamknieta)
	if zamknieta.Session.Status != shared.TerminalSessionStatusExited {
		t.Errorf("karta po zamknięciu ma stan %q", zamknieta.Session.Status)
	}

	var stan string
	if err := baza.QueryRow(
		`SELECT stan FROM terminal_karta WHERE kod = ?`, kartaKod).Scan(&stan); err != nil {
		t.Fatalf("nie można odczytać karty z bazy: %v", err)
	}
	if stan != string(shared.TerminalSessionStatusExited) {
		t.Errorf("w bazie karta ma stan %q, choć rdzeń zameldował jej zamknięcie", stan)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalSessionList,
		shared.TerminalSessionListRequest{WindowId: &oknoKod}, &wykaz)
	if wykaz.Total != 0 {
		t.Errorf("wykaz kart czynnych wciąż pokazuje kartę zamkniętą: %#v", wykaz.Sessions)
	}

	prawda := true
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalSessionList,
		shared.TerminalSessionListRequest{WindowId: &oknoKod, IncludeExited: &prawda}, &wykaz)
	if wykaz.Total != 1 {
		t.Errorf("wykaz z kartami zakończonymi nie pokazuje karty zamkniętej: %#v", wykaz)
	}

	// Druga próba zamknięcia ma odmówić konfliktem, a nie zameldować powodzenia
	// po raz drugi.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandTerminalSessionClose,
		shared.TerminalSessionCloseRequest{SessionId: kartaKod})
	if odmowa.Code != shared.ErrorCodeConflict {
		t.Errorf("powtórne zamknięcie karty odmówiło kodem %s zamiast conflict", odmowa.Code)
	}
}

// TestAnalizaSkryptuWidziBladIFormatuje sprawdza, że analiza naprawdę uruchamia
// program i wraca z jego uwagami, a nie z pustym wykazem podanym jako „treść bez
// zastrzeżeń”.
func TestAnalizaSkryptuWidziBladIFormatuje(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	prawda := true
	var wynik shared.TerminalScriptLintResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalScriptLint,
		shared.TerminalScriptLintRequest{
			// `$y` nigdzie nie jest ustawiona — ShellCheck ma to zgłosić (SC2154).
			Content: "#!/usr/bin/env bash\nx=1\nif [ $x == 1 ]; then echo \"$y\"; fi\n",
			Shell:   shared.TerminalShellBash,
			Format:  &prawda,
		}, &wynik)

	if !wynik.AnalyzerAvailable {
		t.Skip("ShellCheck nie stoi na tej maszynie — analizy nie ma czym zmierzyć")
	}
	if len(wynik.Findings) == 0 {
		t.Fatal("analiza treści z niezdefiniowaną zmienną oddała pusty wykaz uwag, " +
			"czyli zameldowała treść bez zastrzeżeń")
	}
	znaleziona := false
	for _, uwaga := range wynik.Findings {
		if uwaga.Line <= 0 {
			t.Errorf("uwaga bez numeru wiersza jest uwagą, której nie da się wskazać: %#v", uwaga)
		}
		if uwaga.Rule != nil && strings.Contains(*uwaga.Rule, "2154") {
			znaleziona = true
		}
	}
	if !znaleziona {
		t.Errorf("analiza nie zgłosiła użycia niezdefiniowanej zmiennej: %#v", wynik.Findings)
	}
	if wynik.Formatted == nil || strings.TrimSpace(*wynik.Formatted) == "" {
		t.Error("prośba o formatowanie nie oddała treści sformatowanej")
	}

	// Powłoka bez analizatora ma powiedzieć to WPROST, a nie oddać pusty wykaz
	// uwag, który czyta się jak brak zastrzeżeń.
	var bezAnalizatora shared.TerminalScriptLintResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalScriptLint,
		shared.TerminalScriptLintRequest{Content: "dir", Shell: shared.TerminalShellCmd},
		&bezAnalizatora)
	if bezAnalizatora.AnalyzerAvailable {
		t.Error("rdzeń zameldował dostępny analizator dla powłoki, dla której go nie ma")
	}
	if strings.TrimSpace(bezAnalizatora.Analyzer) == "" {
		t.Error("odpowiedź o braku analizatora nie mówi, czego brakuje")
	}
}

// TestObserwacjaUruchamiaPolecenieNaZmianie jest sprawdzianem skutku
// najostrzejszym w tym pliku: mierzy PLIK, który powstał, bo obserwacja naprawdę
// wyzwoliła polecenie, które naprawdę się wykonało. Wpis w bazie o wyzwoleniu
// jest tu sprawdzeniem drugim, nie pierwszym.
func TestObserwacjaUruchamiaPolecenieNaZmianie(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuTerminala(t, katalogDanych)
	roboczy := t.TempDir()
	oknoKod := oknoTerminalaSprawdzianu(t, zmontowany, zycie, roboczy)
	kartaKod := kartaSprawdzianu(t, zmontowany, zycie, oknoKod, roboczy)

	zrodlo := filepath.Join(roboczy, "zrodlo.txt")
	if err := os.WriteFile(zrodlo, []byte("przed\n"), 0o644); err != nil {
		t.Fatalf("nie można przygotować pliku obserwowanego: %v", err)
	}
	swiadek := filepath.Join(roboczy, "swiadek.txt")

	tlumienie := 0
	var zalozona shared.TerminalWatchStartResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalWatchStart,
		shared.TerminalWatchStartRequest{
			SessionId:  kartaKod,
			Pattern:    "*.txt",
			Command:    "echo wyzwolone > " + swiadek,
			DebounceMs: &tlumienie,
		}, &zalozona)
	if zalozona.Watch.Status != shared.TerminalWatchStatusActive {
		t.Fatalf("obserwacja po założeniu ma stan %q", zalozona.Watch.Status)
	}

	// Zmiana musi nastąpić PO pierwszym przeglądzie, bo pierwszy przegląd
	// wyłącznie zapamiętuje stan zastany.
	time.Sleep(odstepPrzegladu + 200*time.Millisecond)
	if err := os.WriteFile(zrodlo, []byte("po zmianie, dłuższa treść\n"), 0o644); err != nil {
		t.Fatalf("nie można zmienić pliku obserwowanego: %v", err)
	}

	if !doczekajPliku(swiadek, 15*time.Second) {
		t.Fatal("obserwacja nie uruchomiła polecenia: plik, który polecenie miało " +
			"utworzyć, nie powstał")
	}
	tresc, err := os.ReadFile(swiadek)
	if err != nil || !strings.Contains(string(tresc), "wyzwolone") {
		t.Errorf("plik świadka powstał, ale nie niesie wyjścia polecenia: %q, %v", tresc, err)
	}

	var licznik int
	if err := baza.QueryRow(
		`SELECT licznik FROM terminal_obserwacja WHERE kod = ?`, zalozona.Watch.Id).Scan(&licznik); err != nil {
		t.Fatalf("nie można odczytać obserwacji z bazy: %v", err)
	}
	if licznik < 1 {
		t.Errorf("licznik wyzwoleń obserwacji stoi na %d, choć polecenie się wykonało", licznik)
	}

	var zatrzymana shared.TerminalWatchStopResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalWatchStop,
		shared.TerminalWatchStopRequest{WatchId: zalozona.Watch.Id}, &zatrzymana)
	if zatrzymana.Watch.Status != shared.TerminalWatchStatusStopped {
		t.Errorf("obserwacja po zatrzymaniu ma stan %q", zatrzymana.Watch.Status)
	}

	// Po zatrzymaniu kolejna zmiana nie ma prawa nic uruchomić.
	if err := os.Remove(swiadek); err != nil {
		t.Fatalf("nie można sprzątnąć pliku świadka: %v", err)
	}
	if err := os.WriteFile(zrodlo, []byte("zmiana po zatrzymaniu obserwacji\n"), 0o644); err != nil {
		t.Fatalf("nie można zmienić pliku obserwowanego: %v", err)
	}
	if doczekajPliku(swiadek, 3*time.Second) {
		t.Error("obserwacja zatrzymana wciąż uruchamia polecenie")
	}
}

// TestTunelNiedostepnegoCeluKonczySieNiepowodzeniem sprawdza, że stan tunelu
// bierze się z PROCESU, a nie z zapisu: `ssh` do celu, którego nie ma, kończy
// się, a rdzeń zapisuje `failed` wraz z powodem. Tunel udany wymagałby serwera
// SSH, którego sprawdzian nie stawia — ale to właśnie ta ścieżka rozstrzyga,
// czy rdzeń w ogóle patrzy na wynik procesu.
func TestTunelNiedostepnegoCeluKonczySieNiepowodzeniem(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuTerminala(t, katalogDanych)
	roboczy := t.TempDir()
	oknoKod := oknoTerminalaSprawdzianu(t, zmontowany, zycie, roboczy)

	cel := "operator@127.0.0.1"
	portDocelowy := 80
	// Port 1 na pętli zwrotnej nie ma nasłuchu, więc połączenie SSH odpada
	// natychmiast — bez czekania na czas sieci.
	portCelu := 1
	var otwarty shared.TerminalTunnelOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalTunnelOpen,
		shared.TerminalTunnelOpenRequest{
			WindowId:     oknoKod,
			Kind:         shared.TerminalTunnelKindLocal,
			RemoteTarget: &cel,
			RemoteHost:   wskaznikTekstu("127.0.0.1"),
			RemotePort:   &portDocelowy,
			LocalPort:    &portCelu,
		}, &otwarty)

	if otwarty.Tunnel.Status != shared.TerminalTunnelStatusFailed {
		t.Fatalf("tunel do celu bez nasłuchu ma stan %q zamiast failed — rdzeń nie patrzy "+
			"na wynik procesu ssh", otwarty.Tunnel.Status)
	}
	if otwarty.Tunnel.ErrorMessage == nil || strings.TrimSpace(*otwarty.Tunnel.ErrorMessage) == "" {
		t.Error("tunel nieudany nie niesie powodu niepowodzenia")
	}
	var stan, powod string
	if err := baza.QueryRow(
		`SELECT stan, powod FROM terminal_tunel WHERE kod = ?`, otwarty.Tunnel.Id).
		Scan(&stan, &powod); err != nil {
		t.Fatalf("nie można odczytać tunelu z bazy: %v", err)
	}
	if stan != string(shared.TerminalTunnelStatusFailed) || strings.TrimSpace(powod) == "" {
		t.Errorf("w bazie tunel ma stan %q i powód %q", stan, powod)
	}

	var wykaz shared.TerminalTunnelListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalTunnelList,
		shared.TerminalTunnelListRequest{WindowId: &oknoKod}, &wykaz)
	if wykaz.Total != 1 || wykaz.Tunnels[0].Id != otwarty.Tunnel.Id {
		t.Errorf("wykaz tuneli nie pokazuje założonego tunelu: %#v", wykaz)
	}

	var zamkniety shared.TerminalTunnelCloseResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalTunnelClose,
		shared.TerminalTunnelCloseRequest{TunnelId: otwarty.Tunnel.Id}, &zamkniety)
	if zamkniety.Tunnel.ClosedAt == nil {
		t.Error("tunel po zamknięciu nie ma chwili zamknięcia")
	}
}

// TestKartaPowlokiUrzadzeniowejMaWierszWBazie pilnuje szkody, która byłaby cicha
// w najgorszy możliwy sposób.
//
// Rdzeń nauczył się czterech powłok sięgających poza jego maszynę — kontenera,
// poda, konsoli szeregowej i sesji Telnet — a warunek CHECK kolumny `powloka`
// wymieniał sześć wartości z migracji 041. Karta takiego rodzaju powstawała
// wtedy w pamięci i DZIAŁAŁA, ale jej zapis odbijał się od warunku, a zapis
// karty z zamysłu nie wywraca czynności. Skutek: karta znika po restarcie
// rdzenia i nic tego nie zapowiada. Sprawdzian mierzy więc WIERSZ, nie odpowiedź.
func TestKartaPowlokiUrzadzeniowejMaWierszWBazie(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuTerminala(t, katalogDanych)
	roboczy := t.TempDir()
	oknoKod := oknoTerminalaSprawdzianu(t, zmontowany, zycie, roboczy)

	for _, przypadek := range []struct {
		powloka shared.TerminalShell
		zadanie shared.TerminalSessionOpenRequest
	}{
		{shared.TerminalShellContainer, shared.TerminalSessionOpenRequest{
			ContainerRef: &shared.TerminalContainerRef{ContainerId: wskaznikTekstu("wydanie-1")},
		}},
		{shared.TerminalShellPod, shared.TerminalSessionOpenRequest{
			ContainerRef: &shared.TerminalContainerRef{PodName: wskaznikTekstu("wydanie-pod")},
		}},
		{shared.TerminalShellSerial, shared.TerminalSessionOpenRequest{
			SerialDevice: wskaznikTekstu("/dev/ttyUSB0"),
		}},
		{shared.TerminalShellTelnet, shared.TerminalSessionOpenRequest{
			RemoteTarget: wskaznikTekstu("192.0.2.10"),
		}},
	} {
		zadanie := przypadek.zadanie
		zadanie.WindowId = oknoKod
		zadanie.Shell = przypadek.powloka

		var karta shared.TerminalSessionOpenResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalSessionOpen, zadanie, &karta)

		var powloka string
		if err := baza.QueryRow(`SELECT powloka FROM terminal_karta WHERE kod = ?`, karta.Session.Id).
			Scan(&powloka); err != nil {
			t.Fatalf("karty powłoki %q nie ma w bazie po udanym otwarciu — zapis odbił się po cichu: %v",
				przypadek.powloka, err)
		}
		if powloka != string(przypadek.powloka) {
			t.Errorf("wiersz karty niesie powłokę %q zamiast %q", powloka, przypadek.powloka)
		}
	}
}

// doczekajPliku czeka na pojawienie się pliku nie dłużej niż podany czas.
func doczekajPliku(sciezka string, najdluzej time.Duration) bool {
	koniec := time.Now().Add(najdluzej)
	for time.Now().Before(koniec) {
		if _, err := os.Stat(sciezka); err == nil {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// TestAnalizaSkryptuPythonaIdzieRuffem sprawdza, że karta `python` dostaje
// analizę REGUŁ, a nie samo orzeczenie o składni.
//
// Treść jest składniowo poprawna, więc orzeczenie o składni oddałoby pusty wykaz
// uwag — czyli zdanie „treść bez zastrzeżeń" o treści, która zastrzeżenia ma.
// Sprawdzian mierzy więc konkretną regułę (`F401`, import nieużywany), a nie
// samą liczbę uwag.
func TestAnalizaSkryptuPythonaIdzieRuffem(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	prawda := true
	var wynik shared.TerminalScriptLintResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalScriptLint,
		shared.TerminalScriptLintRequest{
			Content: "import os, sys\ndef f( x ):\n    return x\n",
			Shell:   shared.TerminalShellPython,
			Format:  &prawda,
		}, &wynik)

	if !wynik.AnalyzerAvailable {
		t.Skipf("program analizy Pythona nie stoi na tej maszynie (%s) — "+
			"analizy nie ma czym zmierzyć", wynik.Analyzer)
	}
	if !strings.Contains(strings.ToLower(wynik.Analyzer), narzedzieRuff.Program) {
		t.Skipf("analizę Pythona prowadzi tu %s, a nie %s — reguł nie ma czym zmierzyć",
			wynik.Analyzer, narzedzieRuff.Nazwa)
	}

	if len(wynik.Findings) == 0 {
		t.Fatal("analiza treści z nieużywanymi importami oddała pusty wykaz uwag, " +
			"czyli zameldowała treść bez zastrzeżeń")
	}
	regulaZnaleziona := false
	for _, uwaga := range wynik.Findings {
		if uwaga.Line <= 0 {
			t.Errorf("uwaga bez numeru wiersza jest uwagą, której nie da się wskazać: %#v", uwaga)
		}
		if uwaga.Rule != nil && *uwaga.Rule == "F401" {
			regulaZnaleziona = true
		}
	}
	if !regulaZnaleziona {
		t.Errorf("analiza nie zgłosiła nieużywanego importu (F401): %#v", wynik.Findings)
	}

	// Formatowanie ma oddać treść RÓŻNĄ od wejściowej — inaczej pole obiecuje
	// pracę, której nie wykonano.
	if wynik.Formatted == nil || strings.TrimSpace(*wynik.Formatted) == "" {
		t.Fatal("prośba o formatowanie nie oddała treści sformatowanej")
	}
	if strings.Contains(*wynik.Formatted, "def f( x ):") {
		t.Errorf("treść oddana jako sformatowana została nietknięta:\n%s", *wynik.Formatted)
	}
}

// TestAnalizaSkryptuPythonaZgłaszaBladSkladni pilnuje, że przejście na analizę
// regułami nie odebrało orzeczenia o składni: błąd składni ma wracać dalej.
func TestAnalizaSkryptuPythonaZglaszaBladSkladni(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var wynik shared.TerminalScriptLintResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTerminalScriptLint,
		shared.TerminalScriptLintRequest{
			Content: "def f(\n",
			Shell:   shared.TerminalShellPython,
		}, &wynik)

	if !wynik.AnalyzerAvailable {
		t.Skipf("program analizy Pythona nie stoi na tej maszynie (%s)", wynik.Analyzer)
	}
	if len(wynik.Findings) == 0 {
		t.Fatal("analiza treści o zepsutej składni oddała pusty wykaz uwag")
	}
}
