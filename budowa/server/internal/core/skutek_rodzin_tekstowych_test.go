package core

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Skutek rodzin obsługi tekstu i mowy poza jednym oknem: historii schowka,
// słownika skrótów, kontekstów pamięci, nagrań mowy, wywoływacza poleceń,
// wyróżnienia wpisu dziennika i pomiaru zajętości okna kontekstu.
//
// Wzorzec sprawdzianu jest ten sam co przy warsztacie PDF: żaden nie kończy się
// na odczytaniu odpowiedzi. Każdy schodzi własnym zapytaniem SQL do tabeli albo
// otwiera plik na dysku i mierzy go niezależnie od tego, co komenda
// zameldowała.

// TestHistoriaSchowkaZostajeWTabeli wykazuje, że treść oddana rdzeniowi
// naprawdę przeżywa kartę: leży w tabeli `wpis_schowka` wraz z odciskiem.
func TestHistoriaSchowkaZostajeWTabeli(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	const tresc = "ustalenia ze spotkania: termin 30.09, właściciel — Dariusz"

	var zapis shared.ClipboardPushResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandClipboardPush,
		shared.ClipboardPushRequest{Content: tresc}, &zapis)
	if zapis.AlreadyPresent {
		t.Fatal("pierwszy zapis treści zameldował powtórzenie")
	}

	var wTabeli, odcisk string
	if err := baza.QueryRow(
		`SELECT tresc, odcisk FROM wpis_schowka WHERE identyfikator_zewnetrzny = ?`,
		zapis.Entry.Id).Scan(&wTabeli, &odcisk); err != nil {
		t.Fatalf("wpis schowka nie zostawił wiersza: %v", err)
	}
	if wTabeli != tresc {
		t.Fatalf("w tabeli leży inna treść niż oddana: %q", wTabeli)
	}
	if odcisk == "" {
		t.Fatal("wiersz nie ma odcisku treści — powtórzenie mnożyłoby wpisy")
	}

	// Powtórzenie tej samej treści ma podnieść wpis zastany, a nie założyć drugi.
	var powtorzenie shared.ClipboardPushResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandClipboardPush,
		shared.ClipboardPushRequest{Content: tresc}, &powtorzenie)
	if !powtorzenie.AlreadyPresent {
		t.Fatal("powtórzenie treści nie zostało rozpoznane")
	}
	if wierszy(t, baza, `SELECT COUNT(*) FROM wpis_schowka`) != 1 {
		t.Fatal("powtórzenie treści zwielokrotniło wpisy historii")
	}
}

// TestCzyszczenieHistoriiOmijaPrzypiete wykazuje, że przypięcie naprawdę chroni
// wpis: po wyczyszczeniu historii wiersz przypięty zostaje w tabeli.
func TestCzyszczenieHistoriiOmijaPrzypiete(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var przypinany, zwykly shared.ClipboardPushResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandClipboardPush,
		shared.ClipboardPushRequest{Content: "numer konta"}, &przypinany)
	wykonajUdana(t, zmontowany, zycie, shared.CommandClipboardPush,
		shared.ClipboardPushRequest{Content: "przypadkowa kopia"}, &zwykly)

	wykonajUdana(t, zmontowany, zycie, shared.CommandClipboardPin,
		shared.ClipboardPinRequest{EntryId: przypinany.Entry.Id, Pinned: true}, nil)

	var czyszczenie shared.ClipboardDeleteResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandClipboardDelete,
		shared.ClipboardDeleteRequest{}, &czyszczenie)
	if czyszczenie.Deleted != 1 {
		t.Fatalf("czyszczenie skasowało %d wpisów, a nieprzypięty był jeden", czyszczenie.Deleted)
	}

	if wierszy(t, baza, `SELECT COUNT(*) FROM wpis_schowka WHERE przypiety = 1`) != 1 {
		t.Fatal("wpis przypięty nie przeżył czyszczenia historii")
	}
	if wierszy(t, baza, `SELECT COUNT(*) FROM wpis_schowka WHERE przypiety = 0`) != 0 {
		t.Fatal("wpis nieprzypięty przeżył czyszczenie historii")
	}
}

// TestSkrotTekstowyLezyWSlownikuRdzenia wykazuje, że skrót jest własnością
// rdzenia, a nie jednego okna: wiersz leży w tabeli wraz z polami szablonu.
func TestSkrotTekstowyLezyWSlownikuRdzenia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var zapis shared.SnippetSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandSnippetSet,
		shared.SnippetSetRequest{
			Shortcut:  ";odmowa",
			Content:   "Dziękuję za propozycję. Tym razem nie skorzystamy — {powod}.",
			Variables: []string{"powod"},
		}, &zapis)

	var skrot, tresc, pola string
	if err := baza.QueryRow(
		`SELECT skrot, tresc, pola_json FROM skrot_tekstowy WHERE identyfikator_zewnetrzny = ?`,
		zapis.Snippet.Id).Scan(&skrot, &tresc, &pola); err != nil {
		t.Fatalf("skrót nie zostawił wiersza w słowniku: %v", err)
	}
	if skrot != ";odmowa" || !strings.Contains(tresc, "{powod}") {
		t.Fatalf("wiersz słownika niesie skrót %q i treść %q", skrot, tresc)
	}
	if !strings.Contains(pola, "powod") {
		t.Fatalf("pola szablonu nie doszły do wiersza: %q", pola)
	}

	// Drugi skrót o tej samej frazie w tym samym profilu ma zostać odrzucony:
	// dwa rozwinięcia jednego skrótu rozstrzygałaby kolejność odczytu.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandSnippetSet,
		shared.SnippetSetRequest{Shortcut: ";odmowa", Content: "coś innego"})
	if odmowa.Code != shared.ErrorCodeConflict {
		t.Fatalf("powtórzony skrót dał kod %q zamiast konfliktu", odmowa.Code)
	}
}

// TestUsuniecieKontekstuNieRuszaWpisowPamieci wykazuje najważniejsze
// rozstrzygnięcie tej rodziny: kontekst jest zestawem wskazań, a nie
// właścicielem treści.
func TestUsuniecieKontekstuNieRuszaWpisowPamieci(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var kontekst shared.MemoryContextSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryContextSave,
		shared.MemoryContextSaveRequest{
			Name:         "projekt Atlas",
			Levels:       []shared.ConfigScope{shared.ConfigScopeProject, shared.ConfigScopeSession},
			EntryIds:     []string{"wpis-jeden", "wpis-dwa"},
			SystemPrompt: wskaznik("Mów zwięźle."),
		}, &kontekst)

	var poziomy, wpisy string
	if err := baza.QueryRow(
		`SELECT poziomy_json, wpisy_json FROM kontekst_pamieci WHERE identyfikator_zewnetrzny = ?`,
		kontekst.Context.Id).Scan(&poziomy, &wpisy); err != nil {
		t.Fatalf("kontekst nie zostawił wiersza: %v", err)
	}
	if !strings.Contains(poziomy, "project") || !strings.Contains(wpisy, "wpis-dwa") {
		t.Fatalf("wskazania kontekstu nie doszły do wiersza: %q / %q", poziomy, wpisy)
	}

	var aktywacja shared.MemoryContextActivateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryContextActivate,
		shared.MemoryContextActivateRequest{
			ContextId: kontekst.Context.Id, SessionId: "karta-sprawdzianu",
		}, &aktywacja)
	if !aktywacja.Activated || len(aktywacja.Levels) != 2 {
		t.Fatalf("aktywacja oddała %+v", aktywacja)
	}
	if wierszy(t, baza,
		`SELECT COUNT(*) FROM kontekst_pamieci_czynny WHERE sesja_kod = 'karta-sprawdzianu'`) != 1 {
		t.Fatal("aktywacja nie zostawiła wskazania kontekstu czynnego karty")
	}

	var usuniecie shared.MemoryContextDeleteResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryContextDelete,
		shared.MemoryContextDeleteRequest{ContextId: kontekst.Context.Id}, &usuniecie)
	if !usuniecie.Deleted {
		t.Fatal("usunięcie kontekstu nie doszło do skutku")
	}
	if wierszy(t, baza, `SELECT COUNT(*) FROM kontekst_pamieci`) != 0 {
		t.Fatal("kontekst przeżył usunięcie")
	}
}

// TestZasadaRetencjiNieKasujeWpisowWstecz wykazuje, że zapis zasady obejmuje
// zapisy KOLEJNE: wpisy pamięci zastane zostają nietknięte, a odpowiedź niesie
// ich policzoną liczbę.
func TestZasadaRetencjiNieKasujeWpisowWstecz(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	// Wpis pamięci projektu zastany przed zapisem zasady.
	if _, err := baza.Exec(
		`INSERT INTO projekt (kod, nazwa) VALUES ('projekt-sprawdzianu', 'sprawdzian')`); err != nil {
		t.Fatalf("nie można założyć projektu materiału: %v", err)
	}
	if _, err := baza.Exec(
		`INSERT INTO wpis_pamieci_projektu
		     (projekt_id, identyfikator_zewnetrzny, tresc, poziom_zasiegu_id)
		 SELECT p.id, 'wpis-zastany', 'spotkania we wtorki', z.id
		   FROM projekt p, poziom_zasiegu z
		  WHERE p.kod = 'projekt-sprawdzianu' AND z.kod = 'projekt'`); err != nil {
		t.Fatalf("nie można założyć wpisu pamięci: %v", err)
	}

	var zasada shared.MemoryRetentionSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryRetentionSet,
		shared.MemoryRetentionSetRequest{
			Scope:              wskaznik(shared.ConfigScope(shared.ConfigScopeGlobal)),
			TtlDays:            wskaznik(30),
			NeverStorePatterns: []string{"numer karty"},
		}, &zasada)

	if zasada.AffectedEntries != 1 {
		t.Fatalf("zasada melduje %d wpisów zastanych, a założono jeden",
			zasada.AffectedEntries)
	}
	if wierszy(t, baza, `SELECT COUNT(*) FROM wpis_pamieci_projektu`) != 1 {
		t.Fatal("zapis zasady retencji skasował wpis zastany — zmiana reguły " +
			"zabrała ustalenie, na które Operator się nie umawiał")
	}

	var odczyt shared.MemoryRetentionGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryRetentionGet,
		shared.MemoryRetentionGetRequest{}, &odczyt)
	if len(odczyt.Policies) != 1 || odczyt.Policies[0].TtlDays != 30 {
		t.Fatalf("odczyt zasad oddał %+v", odczyt.Policies)
	}
}

// TestNagranieWraca wykazuje pełną drogę bajtów: wgrane nagranie leży na dysku
// pod odnośnikiem, a odsłuch oddaje DOKŁADNIE te same bajty.
func TestNagranieWraca(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	// Materiał: nagłówek WAV wraz z krótką próbką. Treść jest dowolna — mierzymy
	// drogę bajtów, a nie rozpoznawanie mowy.
	bajty := nagranieWavProbne()

	var przyjecie shared.SpeechAudioUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandSpeechAudioUpload,
		shared.SpeechAudioUploadRequest{
			Audio: base64.StdEncoding.EncodeToString(bajty), ContentType: "audio/wav",
		}, &przyjecie)

	if przyjecie.SizeBytes != len(bajty) {
		t.Fatalf("rdzeń przyjął %d bajtów, a wysłano %d", przyjecie.SizeBytes, len(bajty))
	}

	// Świadek pierwszy: plik na dysku pod odnośnikiem, o tej samej treści.
	naDysku, err := os.ReadFile(przyjecie.AudioRef)
	if err != nil {
		t.Fatalf("odnośnik nagrania nie prowadzi do pliku: %v", err)
	}
	if string(naDysku) != string(bajty) {
		t.Fatal("bajty na dysku różnią się od wysłanych")
	}

	// Świadek drugi: wiersz rejestru, po którym odsłuch rozstrzyga, że wolno
	// oddać te bajty.
	var typTresci string
	var rozmiar int64
	if err := baza.QueryRow(
		`SELECT typ_tresci, rozmiar_bajtow FROM nagranie_mowy WHERE sciezka = ?`,
		przyjecie.AudioRef).Scan(&typTresci, &rozmiar); err != nil {
		t.Fatalf("nagranie nie zostawiło wiersza rejestru: %v", err)
	}
	if typTresci != "audio/wav" || rozmiar != int64(len(bajty)) {
		t.Fatalf("wiersz rejestru niesie %q i %d bajtów", typTresci, rozmiar)
	}

	// Świadek trzeci: odsłuch oddaje te same bajty, a nie pustkę.
	var odsluch shared.SpeechAudioFetchResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandSpeechAudioFetch,
		shared.SpeechAudioFetchRequest{AudioRef: przyjecie.AudioRef}, &odsluch)
	oddane, err := base64.StdEncoding.DecodeString(odsluch.Audio)
	if err != nil {
		t.Fatalf("odsłuch oddał zapis, którego nie da się odczytać: %v", err)
	}
	if string(oddane) != string(bajty) {
		t.Fatal("odsłuch oddał inne bajty niż przyjęte")
	}
}

// TestOdsluchNieWydajeCudzegoPliku wykazuje granicę odsłuchu: komenda oddaje
// nagrania rdzenia, a nie dowolny plik dysku.
func TestOdsluchNieWydajeCudzegoPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	obcy := katalog + "/../obcy.wav"
	zapiszPlikSprawdzianu(t, obcy, nagranieWavProbne())

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandSpeechAudioFetch,
		shared.SpeechAudioFetchRequest{AudioRef: obcy})
	if odmowa.Code != shared.ErrorCodePermissionDenied {
		t.Fatalf("odsłuch pliku spoza magazynu rdzenia dał kod %q zamiast odmowy uprawnienia",
			odmowa.Code)
	}
}

// TestNastawaWybudzaniaZostajeWKonfiguracji wykazuje, że fraza wybudzająca
// mieszka tam, gdzie każde inne ustawienie — w tabeli `ustawienie`.
func TestNastawaWybudzaniaZostajeWKonfiguracji(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var zapis shared.SpeechWakeSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandSpeechWakeSet,
		shared.SpeechWakeSetRequest{
			Phrase:           wskaznik("Danaco, słuchaj"),
			Mode:             wskaznik(shared.ListenMode(shared.ListenModeWakeWord)),
			VadThreshold:     wskaznik(35),
			NoiseSuppression: wskaznik(true),
		}, &zapis)

	if zapis.Config.Phrase != "Danaco, słuchaj" {
		t.Fatalf("nastawa po zapisie niesie frazę %q", zapis.Config.Phrase)
	}
	if zapis.Config.Mode != shared.ListenModeWakeWord || zapis.Config.VadThreshold != 35 {
		t.Fatalf("nastawa po zapisie niesie tryb %q i próg %d",
			zapis.Config.Mode, zapis.Config.VadThreshold)
	}
	if !zapis.Config.NoiseSuppression {
		t.Fatal("odszumianie nie doszło do nastawy")
	}

	var wartosc string
	if err := baza.QueryRow(
		`SELECT wartosc FROM ustawienie WHERE klucz = ?`, kluczFrazyWybudzajacej).
		Scan(&wartosc); err != nil {
		t.Fatalf("fraza wybudzająca nie zostawiła wiersza ustawienia: %v", err)
	}
	if wartosc != "Danaco, słuchaj" {
		t.Fatalf("w tabeli ustawień leży %q", wartosc)
	}
}

// TestSkrotWywolywaczaZostajeMimoBrakuPowloki wykazuje rozstrzygnięcie
// kontraktu: zapis zostaje nawet wtedy, gdy rejestracji nie ma kto wykonać,
// a odpowiedź mówi o tym wprost zamiast obiecywać skrót, który nikogo nie obudzi.
func TestSkrotWywolywaczaZostajeMimoBrakuPowloki(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var zapis shared.LauncherHotkeySetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLauncherHotkeySet,
		shared.LauncherHotkeySetRequest{Hotkey: "Ctrl+Shift+Space"}, &zapis)

	if zapis.Hotkey != "Ctrl+Shift+Space" {
		t.Fatalf("zapis oddał skrót %q", zapis.Hotkey)
	}
	if zapis.Registered {
		t.Fatal("rdzeń zameldował rejestrację skrótu, choć żadna powłoka jej nie zadeklarowała")
	}
	if zapis.Reason == nil || *zapis.Reason == "" {
		t.Fatal("odpowiedź nie mówi, dlaczego rejestracji nie ma")
	}

	var wartosc string
	if err := baza.QueryRow(`SELECT wartosc FROM ustawienie WHERE klucz = ?`,
		kluczSkrotuWywolywacza).Scan(&wartosc); err != nil {
		t.Fatalf("skrót nie zostawił wiersza ustawienia: %v", err)
	}
	if wartosc != "Ctrl+Shift+Space" {
		t.Fatalf("w tabeli ustawień leży skrót %q", wartosc)
	}

	var odczyt shared.LauncherHotkeyGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLauncherHotkeyGet,
		shared.LauncherHotkeyGetRequest{}, &odczyt)
	if odczyt.Hotkey != "Ctrl+Shift+Space" {
		t.Fatalf("odczyt oddał skrót %q", odczyt.Hotkey)
	}
	if odczyt.Supported {
		t.Fatal("odczyt melduje wsparcie powłoki, której nie ma")
	}

	// Skrót z samych modyfikatorów nie ma prawa przejść: nie da się go wcisnąć.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandLauncherHotkeySet,
		shared.LauncherHotkeySetRequest{Hotkey: "Ctrl+Shift"})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("skrót bez klawisza głównego dał kod %q", odmowa.Code)
	}
}

// TestWyroznienieWpisuDziennikaPrzezywaOdczyt wykazuje, że wyróżnienie ma gdzie
// zamieszkać: kolumna wiersza, a nie pole odpowiedzi znikające po odświeżeniu.
func TestWyroznienieWpisuDziennikaPrzezywaOdczyt(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	if _, err := baza.Exec(
		`INSERT INTO wpis_dziennika_asystenta (kod, okno_kod, rodzaj, tresc, utworzono)
		 VALUES ('wpis-sprawdzianu', 'okno-sprawdzianu', 'note', 'ustalenie', ?)`,
		terazWMilisekundachAlertow()); err != nil {
		t.Fatalf("nie można założyć wpisu dziennika: %v", err)
	}

	var oznaczenie shared.AssistantActivityFlagResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAssistantActivityFlag,
		shared.AssistantActivityFlagRequest{
			EntryId: "wpis-sprawdzianu", Important: true,
			Note: wskaznik("do omówienia na przeglądzie"),
		}, &oznaczenie)
	if oznaczenie.Entry.Important == nil || !*oznaczenie.Entry.Important {
		t.Fatal("odpowiedź nie niesie wyróżnienia")
	}

	var wazny int
	var notatka string
	if err := baza.QueryRow(
		`SELECT wazny, COALESCE(notatka_wyroznienia, '') FROM wpis_dziennika_asystenta
		  WHERE kod = 'wpis-sprawdzianu'`).Scan(&wazny, &notatka); err != nil {
		t.Fatalf("nie można odczytać wpisu po oznaczeniu: %v", err)
	}
	if wazny != 1 || notatka != "do omówienia na przeglądzie" {
		t.Fatalf("wiersz niesie wazny=%d i notatkę %q", wazny, notatka)
	}

	// Zdjęcie wyróżnienia kasuje też powód — powód bez znacznika byłby notatką
	// do wpisu, którego nikt nie wyróżnił.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAssistantActivityFlag,
		shared.AssistantActivityFlagRequest{EntryId: "wpis-sprawdzianu", Important: false}, nil)
	if err := baza.QueryRow(
		`SELECT wazny, COALESCE(notatka_wyroznienia, '') FROM wpis_dziennika_asystenta
		  WHERE kod = 'wpis-sprawdzianu'`).Scan(&wazny, &notatka); err != nil {
		t.Fatalf("nie można odczytać wpisu po zdjęciu wyróżnienia: %v", err)
	}
	if wazny != 0 || notatka != "" {
		t.Fatalf("po zdjęciu wyróżnienia wiersz niesie wazny=%d i notatkę %q", wazny, notatka)
	}
}

// nagranieWavProbne składa najkrótszy poprawny plik WAV: nagłówek RIFF wraz
// z jedną ramką ciszy. Format jest tu istotny — magazyn nagrań przyjmuje
// wyłącznie rozszerzenia, które silnik mowy potrafi otworzyć, więc materiał
// musi być nagraniem naprawdę, a nie napisem udającym nagranie.
func nagranieWavProbne() []byte {
	naglowek := []byte("RIFF")
	naglowek = append(naglowek, 0x2c, 0x00, 0x00, 0x00) // rozmiar pliku - 8
	naglowek = append(naglowek, []byte("WAVEfmt ")...)
	naglowek = append(naglowek, 0x10, 0x00, 0x00, 0x00) // długość bloku formatu
	naglowek = append(naglowek, 0x01, 0x00)             // PCM
	naglowek = append(naglowek, 0x01, 0x00)             // jeden kanał
	naglowek = append(naglowek, 0x40, 0x1f, 0x00, 0x00) // 8000 Hz
	naglowek = append(naglowek, 0x80, 0x3e, 0x00, 0x00) // bajtów na sekundę
	naglowek = append(naglowek, 0x02, 0x00)             // wyrównanie ramki
	naglowek = append(naglowek, 0x10, 0x00)             // 16 bitów na próbkę
	naglowek = append(naglowek, []byte("data")...)
	naglowek = append(naglowek, 0x08, 0x00, 0x00, 0x00) // długość danych
	naglowek = append(naglowek, 0, 0, 0, 0, 0, 0, 0, 0) // cztery ramki ciszy
	return naglowek
}
