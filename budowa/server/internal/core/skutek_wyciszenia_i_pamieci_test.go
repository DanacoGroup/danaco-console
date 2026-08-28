package core

import (
	"testing"
	"time"

	"danacoconsole/shared"
)

// Każdy sprawdzian pyta wprost tabele `wylaczenie_pamieci`, `wyciszenie_nakladki`
// i `sygnal_nakladki`.

// TestWylaczeniePamieciNieKasujeTresci wykazuje różnicę wobec `memory.delete`:
// wpis wyłączony wypada z kontekstu, jego treść zostaje nietknięta w tabeli
// pamięci, a odpowiedź nazywa zasięg, który go wyłączył.
func TestWylaczeniePamieciNieKasujeTresci(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	const projekt = "proj-wylaczenia"
	const tresc = "Nadawca podpisuje listy imieniem Operatora"

	var zapis shared.MemorySetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemorySet,
		shared.MemorySetRequest{ProjectId: wskaz(projekt), Content: tresc}, &zapis)
	if zapis.Entry.Id == "" {
		t.Fatal("memory.set nie oddał wpisu — nie ma czego wyłączać")
	}

	// Przed wyłączeniem wpis wchodzi do kontekstu.
	var przed shared.MemoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryList,
		shared.MemoryListRequest{ProjectId: wskaz(projekt)}, &przed)
	if len(przed.Entries) != 1 {
		t.Fatalf("przed wyłączeniem pamięć projektu ma %d wpisów, a założono jeden", len(przed.Entries))
	}
	if len(przed.DisabledEntries) != 0 {
		t.Fatalf("bez wyłączeń wykaz wstrzymanych ma %d pozycji", len(przed.DisabledEntries))
	}

	var wylaczenie shared.MemoryDisableSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryDisableSet,
		shared.MemoryDisableSetRequest{
			EntryId:  wskaz(zapis.Entry.Id),
			Scope:    shared.ConfigScopeModule,
			ScopeId:  wskaz("mod_developer"),
			Disabled: true,
		}, &wylaczenie)
	if !wylaczenie.Changed {
		t.Fatal("pierwsze wyłączenie nie zameldowało zmiany wykazu")
	}

	// Świadek niezależny numer jeden: wiersz wyłączenia wraz z zasięgiem.
	var zasieg, bytZasiegu string
	if err := baza.QueryRow(
		`SELECT z.kod, w.klucz_zasiegu
		   FROM wylaczenie_pamieci w
		   JOIN poziom_zasiegu z ON z.id = w.poziom_zasiegu_id
		  WHERE w.identyfikator_zewnetrzny = ?`, wylaczenie.Disable.Id).
		Scan(&zasieg, &bytZasiegu); err != nil {
		t.Fatalf("po wyłączeniu nie ma wiersza w `wylaczenie_pamieci` — cisza bez zapisu: %v", err)
	}
	if zasieg != "modul" || bytZasiegu != "mod_developer" {
		t.Fatalf("wiersz wyłączenia niesie zasięg %q/%q, a wyłączono moduł mod_developer",
			zasieg, bytZasiegu)
	}

	// Świadek niezależny numer dwa — najważniejszy: treść wpisu ŻYJE.
	var trescPoWylaczeniu string
	if err := baza.QueryRow(
		`SELECT tresc FROM wpis_pamieci_projektu WHERE identyfikator_zewnetrzny = ?`,
		zapis.Entry.Id).Scan(&trescPoWylaczeniu); err != nil {
		t.Fatalf("wyłączenie skasowało wiersz pamięci — to jest `memory.delete`, nie wyłączenie: %v", err)
	}
	if trescPoWylaczeniu != tresc {
		t.Fatalf("treść wpisu po wyłączeniu brzmi %q, a zapisano %q", trescPoWylaczeniu, tresc)
	}

	// Odczyt pamięci: wpis nie wchodzi do kontekstu, ale jest nazwany wraz
	// z zasięgiem, który go wyłączył.
	var po shared.MemoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryList,
		shared.MemoryListRequest{ProjectId: wskaz(projekt)}, &po)
	if len(po.Entries) != 0 {
		t.Fatalf("wpis wyłączony nadal wchodzi do kontekstu: %d pozycji", len(po.Entries))
	}
	if len(po.DisabledEntries) != 1 {
		t.Fatalf("wykaz wstrzymanych ma %d pozycji — cisza bez powodu jest gorsza od wyłączenia",
			len(po.DisabledEntries))
	}
	wstrzymany := po.DisabledEntries[0]
	if wstrzymany.EntryId != zapis.Entry.Id {
		t.Fatalf("wykaz wstrzymanych wskazuje wpis %q, a wyłączono %q",
			wstrzymany.EntryId, zapis.Entry.Id)
	}
	if wstrzymany.Scope != shared.ConfigScopeModule {
		t.Fatalf("wykaz wstrzymanych nazywa zasięg %q, a wyłączono w module", wstrzymany.Scope)
	}
	if wstrzymany.DisableId != wylaczenie.Disable.Id {
		t.Fatal("wykaz wstrzymanych nie podaje wyłączenia, którym Operator znosi ciszę jednym ruchem")
	}

	// Zniesienie jednym ruchem — tą samą komendą, polem `disabled`.
	var zniesienie shared.MemoryDisableSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryDisableSet,
		shared.MemoryDisableSetRequest{
			DisableId: wskaz(wylaczenie.Disable.Id),
			EntryId:   wskaz(zapis.Entry.Id),
			Scope:     shared.ConfigScopeModule,
			ScopeId:   wskaz("mod_developer"),
			Disabled:  false,
		}, &zniesienie)
	if !zniesienie.Changed || len(zniesienie.Disables) != 0 {
		t.Fatalf("zniesienie nie opróżniło wykazu wyłączeń: changed=%v, wykaz=%d",
			zniesienie.Changed, len(zniesienie.Disables))
	}

	var zostalo int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM wylaczenie_pamieci`).Scan(&zostalo); err != nil {
		t.Fatalf("nie można policzyć wyłączeń: %v", err)
	}
	if zostalo != 0 {
		t.Fatalf("po zniesieniu w tabeli stoi %d wyłączeń", zostalo)
	}

	var wrocil shared.MemoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryList,
		shared.MemoryListRequest{ProjectId: wskaz(projekt)}, &wrocil)
	if len(wrocil.Entries) != 1 || wrocil.Entries[0].Content != tresc {
		t.Fatalf("po zniesieniu wyłączenia wpis nie wrócił w całości: %d pozycji", len(wrocil.Entries))
	}
}

// TestWylaczeniePoziomuPamieciObejmujeWpisyPoziomu wykazuje drugi byt wyłączany:
// cały poziom pamięci, bez wskazania wpisu. Bez tego „wyłącz pamięć projektu"
// nie miałoby gdzie się zapisać.
func TestWylaczeniePoziomuPamieciObejmujeWpisyPoziomu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const projekt = "proj-poziom"

	var zapis shared.MemorySetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemorySet,
		shared.MemorySetRequest{
			ProjectId: wskaz(projekt),
			Content:   "Ustalenie poziomu projektu",
			Scope:     wskazZasieg(shared.ConfigScopeProject),
		}, &zapis)

	var wylaczenie shared.MemoryDisableSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryDisableSet,
		shared.MemoryDisableSetRequest{
			Level:    wskazPoziom(shared.MemoryLevelProject),
			Scope:    shared.ConfigScopeGlobal,
			Disabled: true,
		}, &wylaczenie)
	if wylaczenie.Disable.Level == nil || *wylaczenie.Disable.Level != shared.MemoryLevelProject {
		t.Fatal("wyłączenie poziomu nie oddało poziomu, którego dotyczy")
	}

	var po shared.MemoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryList,
		shared.MemoryListRequest{ProjectId: wskaz(projekt)}, &po)
	if len(po.Entries) != 0 || len(po.DisabledEntries) != 1 {
		t.Fatalf("wyłączenie poziomu projektu nie objęło wpisu poziomu: czynne=%d, wstrzymane=%d",
			len(po.Entries), len(po.DisabledEntries))
	}
	if po.DisabledEntries[0].Level == nil {
		t.Fatal("wykaz wstrzymanych nie mówi, że wyłączono cały poziom, a nie sam wpis")
	}

	// Odczyt wyłączeń zawężony poziomem — droga okna konfiguracji.
	var wykaz shared.MemoryDisableListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandMemoryDisableList,
		shared.MemoryDisableListRequest{Level: wskazPoziom(shared.MemoryLevelProject)}, &wykaz)
	if len(wykaz.Disables) != 1 {
		t.Fatalf("odczyt wyłączeń poziomu oddał %d pozycji", len(wykaz.Disables))
	}
}

// TestWylaczeniePamieciOdmawiaNazywajacBrak wykazuje, że odmowy nazywają brak
// po właściwej stronie, a nie kończą się pustą odpowiedzią udaną.
func TestWylaczeniePamieciOdmawiaNazywajacBrak(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	bezBytu := wykonajOdmowna(t, zmontowany, zycie, shared.CommandMemoryDisableSet,
		shared.MemoryDisableSetRequest{Scope: shared.ConfigScopeGlobal, Disabled: true})
	if bezBytu.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("wyłączenie bez bytu odmówiło kodem %q", bezBytu.Code)
	}

	bezWyciszenia := wykonajOdmowna(t, zmontowany, zycie, shared.CommandMemoryDisableSet,
		shared.MemoryDisableSetRequest{
			Level:    wskazPoziom(shared.MemoryLevelGlobal),
			Scope:    shared.ConfigScopeGlobal,
			Disabled: false,
		})
	if bezWyciszenia.Code != shared.ErrorCodeNotFound {
		t.Fatalf("zniesienie wyłączenia, którego nie ma, odmówiło kodem %q", bezWyciszenia.Code)
	}
}

// TestWyciszenieNakladkiZostawiaWierszIWstrzymujeSygnal wykazuje, że wyciszenie
// jest bytem rdzenia, a nie stanem jednego okna: ma wiersz, wraca odczytem
// i naprawdę wstrzymuje sygnał klasy zdarzeń.
func TestWyciszenieNakladkiZostawiaWierszIWstrzymujeSygnal(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	// Sygnał klasy „wynik kontroli jakości", nienależnej do telemetrii rdzenia.
	var sygnal shared.AodSignalReportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodSignalReport,
		shared.AodSignalReportRequest{
			EventClass: shared.AodEventClassQualityControlResult,
			Text:       "Kontrola jakości odrzuciła wynik po raz drugi",
			ModuleId:   wskaz("mod_studio"),
		}, &sygnal)
	if sygnal.Suppressed {
		t.Fatal("sygnał wpadł w wyciszenie, choć żadnego nie ma")
	}

	var klasa, tresc string
	if err := baza.QueryRow(
		`SELECT klasa_zdarzen, tresc FROM sygnal_nakladki WHERE identyfikator_zewnetrzny = ?`,
		sygnal.Signal.Id).Scan(&klasa, &tresc); err != nil {
		t.Fatalf("zgłoszony sygnał nie ma wiersza — nośnika nie ma: %v", err)
	}
	if klasa != shared.AodEventClassQualityControlResult {
		t.Fatalf("wiersz sygnału niesie klasę %q", klasa)
	}

	var wyciszenie shared.AodMuteSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodMuteSet,
		shared.AodMuteSetRequest{
			Kind:       wskazRodzajWyciszenia(shared.AodMuteKindEventClass),
			EventClass: wskazKlase(shared.AodEventClassQualityControlResult),
			Muted:      true,
		}, &wyciszenie)
	if !wyciszenie.Changed || len(wyciszenie.Mutes) != 1 {
		t.Fatalf("wyciszenie klasy nie weszło do wykazu: changed=%v, wykaz=%d",
			wyciszenie.Changed, len(wyciszenie.Mutes))
	}

	var rodzaj, zakres, klasaWiersza string
	if err := baza.QueryRow(
		`SELECT rodzaj, zakres, klasa_zdarzen FROM wyciszenie_nakladki
		  WHERE identyfikator_zewnetrzny = ?`, wyciszenie.Mutes[0].Id).
		Scan(&rodzaj, &zakres, &klasaWiersza); err != nil {
		t.Fatalf("wyciszenie nie ma wiersza — zostało stanem jednego okna: %v", err)
	}
	if rodzaj != shared.AodMuteKindEventClass || klasaWiersza != shared.AodEventClassQualityControlResult {
		t.Fatalf("wiersz wyciszenia niesie rodzaj %q i klasę %q", rodzaj, klasaWiersza)
	}

	// Odczyt drugą drogą — tą, którą po połączeniu idzie druga powłoka Operatora.
	var odczyt shared.AodMuteGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodMuteGet,
		shared.AodMuteGetRequest{}, &odczyt)
	if len(odczyt.Mutes) != 1 {
		t.Fatalf("aod.mute.get oddał %d wyciszeń — wyciszenie nie przechodzi między powłokami",
			len(odczyt.Mutes))
	}

	// Skutek wyciszenia: sygnał tej klasy nie ujawnia się, ale jest NAZWANY.
	var wykaz shared.AodSignalListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodSignalList,
		shared.AodSignalListRequest{}, &wykaz)
	if len(wykaz.Signals) != 0 {
		t.Fatalf("sygnał wyciszonej klasy nadal się ujawnia: %d pozycji", len(wykaz.Signals))
	}
	if wykaz.SuppressedCount != 1 || len(wykaz.Mutes) != 1 {
		t.Fatalf("wykaz nie nazywa ciszy: wstrzymanych=%d, wyciszeń=%d",
			wykaz.SuppressedCount, len(wykaz.Mutes))
	}

	// Sygnał przy czynnym wyciszeniu nadal się odkłada — wyciszenie wstrzymuje
	// ujawnienie, nie zapis.
	var drugi shared.AodSignalReportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodSignalReport,
		shared.AodSignalReportRequest{
			EventClass: shared.AodEventClassQualityControlResult,
			Text:       "Kontrola jakości odrzuciła wynik po raz trzeci",
		}, &drugi)
	if !drugi.Suppressed || drugi.Mute == nil {
		t.Fatal("odpowiedź nie mówi, że sygnał wpadł w wyciszenie i w które")
	}
	var odlozonych int
	if err := baza.QueryRow(`SELECT COUNT(*) FROM sygnal_nakladki`).Scan(&odlozonych); err != nil {
		t.Fatalf("nie można policzyć sygnałów: %v", err)
	}
	if odlozonych != 2 {
		t.Fatalf("w tabeli sygnałów stoi %d wierszy — wyciszenie skasowało zapis, a miało "+
			"wstrzymać ujawnienie", odlozonych)
	}

	// Zniesienie jednym ruchem, bez znajomości identyfikatora: rodzajem i klasą.
	var zniesienie shared.AodMuteSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodMuteSet,
		shared.AodMuteSetRequest{
			Kind:       wskazRodzajWyciszenia(shared.AodMuteKindEventClass),
			EventClass: wskazKlase(shared.AodEventClassQualityControlResult),
			Muted:      false,
		}, &zniesienie)
	if !zniesienie.Changed || len(zniesienie.Mutes) != 0 {
		t.Fatalf("zniesienie nie opróżniło wykazu: changed=%v, wykaz=%d",
			zniesienie.Changed, len(zniesienie.Mutes))
	}

	var poZniesieniu shared.AodSignalListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodSignalList,
		shared.AodSignalListRequest{}, &poZniesieniu)
	if len(poZniesieniu.Signals) != 2 || poZniesieniu.SuppressedCount != 0 {
		t.Fatalf("po zniesieniu wyciszenia sygnały nie wróciły: %d pozycji, wstrzymanych %d",
			len(poZniesieniu.Signals), poZniesieniu.SuppressedCount)
	}
}

// TestWyciszenieCzasoweMaKoniecIObowiazujeWszedzie wykazuje dwie rzeczy naraz:
// wyciszenie czasowe bez chwili końca jest odmawiane, a wyciszenie czasowe
// obejmuje sygnał każdej klasy.
func TestWyciszenieCzasoweMaKoniecIObowiazujeWszedzie(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	bezKonca := wykonajOdmowna(t, zmontowany, zycie, shared.CommandAodMuteSet,
		shared.AodMuteSetRequest{
			Kind:  wskazRodzajWyciszenia(shared.AodMuteKindTimed),
			Muted: true,
		})
	if bezKonca.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("wyciszenie czasowe bez końca odmówiło kodem %q", bezKonca.Code)
	}

	koniec := time.Now().UTC().Add(time.Hour).UnixMilli()
	var wyciszenie shared.AodMuteSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodMuteSet,
		shared.AodMuteSetRequest{
			Kind:   wskazRodzajWyciszenia(shared.AodMuteKindTimed),
			EndsAt: &koniec,
			Muted:  true,
		}, &wyciszenie)
	if wyciszenie.Mutes[0].EndsAt == nil {
		t.Fatal("wyciszenie czasowe wróciło bez chwili końca — Operator nie wie, do kiedy milczy")
	}

	// Drugi czas ZASTĘPUJE pierwszy: dwa „do kiedy" naraz nie są odpowiedzią.
	drugiKoniec := time.Now().UTC().Add(2 * time.Hour).UnixMilli()
	var drugie shared.AodMuteSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodMuteSet,
		shared.AodMuteSetRequest{
			Kind:   wskazRodzajWyciszenia(shared.AodMuteKindTimed),
			EndsAt: &drugiKoniec,
			Muted:  true,
		}, &drugie)
	if len(drugie.Mutes) != 1 {
		t.Fatalf("po drugim czasie w wykazie stoi %d wyciszeń czasowych", len(drugie.Mutes))
	}

	var sygnal shared.AodSignalReportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodSignalReport,
		shared.AodSignalReportRequest{
			EventClass: shared.AodEventClassSchedule,
			Text:       "Automatyka uruchomiła się z harmonogramu",
		}, &sygnal)
	if !sygnal.Suppressed {
		t.Fatal("wyciszenie czasowe nie objęło sygnału klasy harmonogramu, choć obejmuje wszystko")
	}
}

// TestWyciszenieKontekstoweNieMilczyBezPodstawy wykazuje granicę wyciszenia
// kontekstowego: obejmuje wskazany moduł, nie sąsiedni, a sygnału bez modułu nie
// wstrzymuje — cisza bez podstawy jest gorsza od ujawnienia.
func TestWyciszenieKontekstoweNieMilczyBezPodstawy(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	bezBytu := wykonajOdmowna(t, zmontowany, zycie, shared.CommandAodMuteSet,
		shared.AodMuteSetRequest{
			Kind:  wskazRodzajWyciszenia(shared.AodMuteKindContextual),
			Scope: wskazZakresWyciszenia(shared.AodMuteScopeModule),
			Muted: true,
		})
	if bezBytu.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("wyciszenie kontekstowe bez bytu odmówiło kodem %q", bezBytu.Code)
	}

	var wyciszenie shared.AodMuteSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAodMuteSet,
		shared.AodMuteSetRequest{
			Kind:      wskazRodzajWyciszenia(shared.AodMuteKindContextual),
			Scope:     wskazZakresWyciszenia(shared.AodMuteScopeModule),
			ScopeId:   wskaz("mod_studio"),
			ScopeName: wskaz("Studio"),
			Muted:     true,
		}, &wyciszenie)
	if wyciszenie.Mutes[0].ScopeName == nil || *wyciszenie.Mutes[0].ScopeName != "Studio" {
		t.Fatal("wyciszenie nie niesie nazwy bytu — Operator nie dowie się, CO milczy")
	}

	wykazSygnalow := []struct {
		modul          string
		wstrzymany     bool
		nazwaPrzypadku string
	}{
		{modul: "mod_studio", wstrzymany: true, nazwaPrzypadku: "wyciszony moduł"},
		{modul: "mod_terminal", wstrzymany: false, nazwaPrzypadku: "moduł sąsiedni"},
		{modul: "", wstrzymany: false, nazwaPrzypadku: "sygnał bez modułu"},
	}
	for _, przypadek := range wykazSygnalow {
		zadanie := shared.AodSignalReportRequest{
			EventClass: shared.AodEventClassOperatorWorkContext,
			Text:       "Powtarzalność czynności Operatora",
		}
		if przypadek.modul != "" {
			zadanie.ModuleId = wskaz(przypadek.modul)
		}
		var odpowiedz shared.AodSignalReportResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandAodSignalReport, zadanie, &odpowiedz)
		if odpowiedz.Suppressed != przypadek.wstrzymany {
			t.Fatalf("%s: wstrzymanie oddane jako %v, oczekiwano %v",
				przypadek.nazwaPrzypadku, odpowiedz.Suppressed, przypadek.wstrzymany)
		}
	}
}

// TestKatalogUstawienNiesiePamiecIRegulyWyciszania wykazuje, że okno konfiguracji
// ma czym sterować: pozycje stoją w katalogu ustawień i idą rodziną `config.*`,
// bez drugiej drogi komend.
func TestKatalogUstawienNiesiePamiecIRegulyWyciszania(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var kategorie shared.SettingsCategoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandSettingsCategoryList,
		shared.SettingsCategoryListRequest{}, &kategorie)
	kody := map[string]bool{}
	for _, kategoria := range kategorie.Categories {
		kody[kategoria.Id] = true
	}
	for _, wymagana := range []string{"pamiec", "aod"} {
		if !kody[wymagana] {
			t.Fatalf("katalog kategorii nie niesie zakresu %q — okno konfiguracji nie ma "+
				"czym sterować", wymagana)
		}
	}

	var definicje shared.SettingsDefinitionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandSettingsDefinitionList,
		shared.SettingsDefinitionListRequest{}, &definicje)
	klucze := map[string]shared.SettingDefinition{}
	for _, definicja := range definicje.Definitions {
		klucze[definicja.Key] = definicja
	}
	for _, klucz := range []string{"memory.level", "memory.enabled", "memory.session.detached",
		"aod.mute.rules"} {
		definicja, jest := klucze[klucz]
		if !jest {
			t.Fatalf("katalog ustawień nie niesie pozycji %q", klucz)
		}
		if definicja.Description == nil || *definicja.Description == "" {
			t.Fatalf("pozycja %q nie ma objaśnienia kontekstowego [?]", klucz)
		}
	}

	// Wyłączenie pamięci ma się dać zapisać w zasięgu modułu i karty sesji.
	zasiegi := map[shared.ConfigScope]bool{}
	for _, zasieg := range klucze["memory.enabled"].AllowedScopes {
		zasiegi[zasieg] = true
	}
	if !zasiegi[shared.ConfigScopeModule] || !zasiegi[shared.ConfigScopeSession] {
		t.Fatalf("pozycja memory.enabled nie daje zasięgu modułu ani karty sesji: %v",
			klucze["memory.enabled"].AllowedScopes)
	}
}

// ── wskazania pól opcjonalnych ─────────────────────────────────────────────

func wskaz(wartosc string) *string { return &wartosc }

func wskazZasieg(zasieg shared.ConfigScope) *shared.ConfigScope { return &zasieg }

func wskazPoziom(poziom shared.MemoryLevel) *shared.MemoryLevel { return &poziom }

func wskazRodzajWyciszenia(rodzaj shared.AodMuteKind) *shared.AodMuteKind { return &rodzaj }

func wskazZakresWyciszenia(zakres shared.AodMuteScope) *shared.AodMuteScope { return &zakres }

func wskazKlase(klasa shared.AodEventClass) *shared.AodEventClass { return &klasa }
