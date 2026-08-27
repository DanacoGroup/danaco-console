package core

import (
	"database/sql"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Sprawdziany skutku odwracalnego dziennika i przełącznika zmian wykonawców.

// dziennikSkutekTresc jest treścią o trzech osobnych akapitach, po jednym na
// każdą z trzech czynności cofania nie po kolei w tym pliku.
const dziennikSkutekTresc = "Akapit pierwszy bez zmian.\n" +
	"Akapit drugi do poprawy.\n" +
	"Akapit trzeci bez zmian."

// dziennikSkutekUprzaz składa rdzeń, dokument i połączenie pomiarowe do bazy,
// wspólne dla sprawdzianów tego pliku.
type dziennikSkutekUprzaz struct {
	*blokadaUprzazSprawdzianu
}

// dziennikSkutekZmontuj zakłada dokument o znanej treści `dziennikSkutekTresc`
// i zapisuje jego pierwszą wersję.
func dziennikSkutekZmontuj(t *testing.T) *dziennikSkutekUprzaz {
	t.Helper()

	uprzaz := blokadaZmontuj(t)
	var zapis shared.StudioDocumentSaveResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: uprzaz.dokument, Content: dziennikSkutekTresc,
			CreateVersion: blokadaWskaznikPrawdy(),
		}, &zapis)
	return &dziennikSkutekUprzaz{blokadaUprzazSprawdzianu: uprzaz}
}

// dziennikSkutekZmienTekst wykonuje zmianę treści ręką wykonawcy o wskazanym
// kodzie i oddaje kod czynności odłożonej w dzienniku.
func (u *dziennikSkutekUprzaz) dziennikSkutekZmienTekst(t *testing.T,
	od, do int, tekst, agent string) {
	t.Helper()

	kodAgenta := agent
	odpowiedz := u.blokadaWykonajJakoModel(t, shared.CommandStudioTextEdit,
		shared.StudioTextEditRequest{
			DocumentId: u.dokument, RangeStart: od, RangeEnd: do,
			Text: tekst, AgentId: &kodAgenta, AgentName: &kodAgenta,
		})
	if odpowiedz.Error != nil {
		t.Fatalf("zmiana treści ręką %s odmówiła: kod=%s treść=%s",
			agent, odpowiedz.Error.Code, odpowiedz.Error.Message)
	}
}

// dziennikSkutekWpisy czyta wpisy dziennika wprost z bazy, od najstarszego,
// niezależnie od odpowiedzi rdzenia.
func (u *dziennikSkutekUprzaz) dziennikSkutekWpisy(t *testing.T) []dziennikSkutekWpis {
	t.Helper()

	wiersze, err := u.baza.Query(
		`SELECT c.identyfikator_zewnetrzny, c.rodzaj, c.autor_rodzaj,
		        COALESCE(c.autor_agent_kod, ''), c.stan, c.kolejnosc
		 FROM czynnosc_dokumentu_studio c
		 JOIN dokument_studio d ON d.id = c.dokument_id
		 WHERE d.identyfikator_zewnetrzny = ?
		 ORDER BY c.kolejnosc`, u.dokument)
	if err != nil {
		t.Fatalf("nie można odczytać dziennika z bazy: %v", err)
	}
	defer wiersze.Close()

	wpisy := []dziennikSkutekWpis{}
	for wiersze.Next() {
		var wpis dziennikSkutekWpis
		if err := wiersze.Scan(&wpis.Kod, &wpis.Rodzaj, &wpis.AutorRodzaj,
			&wpis.AgentKod, &wpis.Stan, &wpis.Kolejnosc); err != nil {
			t.Fatalf("nieczytelny wiersz dziennika: %v", err)
		}
		wpisy = append(wpisy, wpis)
	}
	if err := wiersze.Err(); err != nil {
		t.Fatalf("przerwany odczyt dziennika: %v", err)
	}
	return wpisy
}

type dziennikSkutekWpis struct {
	Kod         string
	Rodzaj      string
	AutorRodzaj string
	AgentKod    string
	Stan        string
	Kolejnosc   int64
}

// ── Dziennik: cofnięcie pojedyncze i NIE PO KOLEI ───────────────────────────

// TestDziennikSkutekCofniecieZeSrodkaZostawiaPozniejsza mierzy sedno wymagania:
// cofnięcie czynności ze ŚRODKA dziennika zdejmuje wyłącznie to, co zrobiła ta
// jedna czynność, a praca późniejsza zostaje.
func TestDziennikSkutekCofniecieZeSrodkaZostawiaPozniejsza(t *testing.T) {
	uprzaz := dziennikSkutekZmontuj(t)

	// Trzy czynności na trzech akapitach, od końca do początku.
	trzeci := strings.Index(dziennikSkutekTresc, "Akapit trzeci bez zmian.")
	uprzaz.dziennikSkutekZmienTekst(t, trzeci, trzeci+len("Akapit trzeci bez zmian."),
		"Akapit trzeci PO ZMIANIE.", "agent-pierwszy")

	drugi := strings.Index(dziennikSkutekTresc, "Akapit drugi do poprawy.")
	uprzaz.dziennikSkutekZmienTekst(t, drugi, drugi+len("Akapit drugi do poprawy."),
		"Akapit drugi POPRAWIONY.", "agent-pierwszy")

	pierwszy := strings.Index(dziennikSkutekTresc, "Akapit pierwszy bez zmian.")
	uprzaz.dziennikSkutekZmienTekst(t, pierwszy, pierwszy+len("Akapit pierwszy bez zmian."),
		"Akapit pierwszy PO ZMIANIE.", "agent-pierwszy")

	wpisy := uprzaz.dziennikSkutekWpisy(t)
	if len(wpisy) < 3 {
		t.Fatalf("dziennik nie odłożył wpisu dla każdej czynności: stoi %d, miało 3", len(wpisy))
	}

	// Czynność środkowa, druga z trzech — sprawdzian nie cofa ostatniej.
	srodkowa := wpisy[1].Kod

	var cofniecie shared.StudioJournalRevertResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioJournalRevert,
		shared.StudioJournalRevertRequest{
			DocumentId: uprzaz.dokument, ActionIds: []string{srodkowa},
		}, &cofniecie)
	if len(cofniecie.Reverted) != 1 || cofniecie.Reverted[0] != srodkowa {
		t.Fatalf("cofnięto co innego niż wskazano: %+v", cofniecie.Reverted)
	}

	// Pomiar niezależny: akapit środkowy wrócił, dwa pozostałe zostały zmienione.
	po := postacTekstFormy(&cofniecie.Form)
	if !strings.Contains(po, "Akapit drugi do poprawy.") {
		t.Errorf("cofnięcie czynności środkowej nie przywróciło jej akapitu; treść: %q", po)
	}
	if !strings.Contains(po, "Akapit pierwszy PO ZMIANIE.") {
		t.Errorf("cofnięcie ze środka zabrało pracę PÓŹNIEJSZĄ — jest przywróceniem "+
			"wersji, nie cofnięciem pojedynczym; treść: %q", po)
	}
	if !strings.Contains(po, "Akapit trzeci PO ZMIANIE.") {
		t.Errorf("cofnięcie ze środka zabrało pracę WCZEŚNIEJSZĄ; treść: %q", po)
	}

	// Stan wpisu w bazie: cofnięty, nie usunięty — dziennik nic nie kasuje.
	stanWpisu := uprzaz.dziennikSkutekStanWpisu(t, srodkowa)
	if stanWpisu != string(shared.StudioActionStateReverted) {
		t.Errorf("wpis dziennika nie stoi w stanie „reverted” po cofnięciu: %q", stanWpisu)
	}
}

// dziennikSkutekStanWpisu czyta stan wpisu dziennika z bazy, wskazanego kodem
// czynności, wprost, niezależnie od odpowiedzi rdzenia.
func (u *dziennikSkutekUprzaz) dziennikSkutekStanWpisu(t *testing.T, kod string) string {
	t.Helper()

	var stan string
	err := u.baza.QueryRow(
		`SELECT stan FROM czynnosc_dokumentu_studio WHERE identyfikator_zewnetrzny = ?`,
		kod).Scan(&stan)
	if err != nil {
		t.Fatalf("nie można odczytać stanu wpisu %s: %v", kod, err)
	}
	return stan
}

// TestDziennikSkutekZaleznoscOdmawiaZamiastPsuc mierzy, że cofnięcie czynności,
// na której stoi późniejsza, ODMAWIA i nazywa zależność — a nie zostawia
// dokumentu w stanie niespójnym.
func TestDziennikSkutekZaleznoscOdmawiaZamiastPsuc(t *testing.T) {
	uprzaz := dziennikSkutekZmontuj(t)

	// Dwie czynności na tym samym fragmencie: druga stoi na pierwszej.
	drugi := strings.Index(dziennikSkutekTresc, "Akapit drugi do poprawy.")
	koniec := drugi + len("Akapit drugi do poprawy.")
	uprzaz.dziennikSkutekZmienTekst(t, drugi, koniec, "Akapit drugi krok pierwszy.", "agent-pierwszy")
	uprzaz.dziennikSkutekZmienTekst(t, drugi, drugi+len("Akapit drugi krok pierwszy."),
		"Akapit drugi krok drugi.", "agent-pierwszy")

	wpisy := uprzaz.dziennikSkutekWpisy(t)
	if len(wpisy) < 2 {
		t.Fatalf("dziennik nie odłożył dwóch czynności: stoi %d", len(wpisy))
	}
	podstawa := wpisy[0].Kod

	trescPrzed := uprzaz.blokadaTrescZBazy(t)
	odmowa := wykonajOdmowna(t, uprzaz.zmontowany, uprzaz.zycie,
		shared.CommandStudioJournalRevert, shared.StudioJournalRevertRequest{
			DocumentId: uprzaz.dokument, ActionIds: []string{podstawa},
		})
	if !strings.Contains(odmowa.Message, wpisy[1].Kod) {
		t.Errorf("odmowa nie NAZYWA czynności, która stoi na cofanej: %q", odmowa.Message)
	}
	if !strings.Contains(odmowa.Message, "cofnij najpierw") {
		t.Errorf("odmowa nie mówi Operatorowi, co zrobić: %q", odmowa.Message)
	}

	// Pomiar niezależny: dokument nietknięty i wpis dalej czynny.
	if po := uprzaz.blokadaTrescZBazy(t); po != trescPrzed {
		t.Errorf("odmowa wróciła, a treść się zmieniła:\nprzed %q\npo    %q", trescPrzed, po)
	}
	if stan := uprzaz.dziennikSkutekStanWpisu(t, podstawa); stan != string(shared.StudioActionStateActive) {
		t.Errorf("wpis przeszedł w stan %q mimo odmowy", stan)
	}

	// Cofnięcie OBU razem przechodzi: zależność cofana razem z podstawą
	// przeszkodą nie jest.
	var razem shared.StudioJournalRevertResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioJournalRevert,
		shared.StudioJournalRevertRequest{
			DocumentId: uprzaz.dokument, ActionIds: []string{podstawa, wpisy[1].Kod},
		}, &razem)
	if len(razem.Reverted) != 2 {
		t.Errorf("cofnięcie obu czynności razem nie doszło do skutku: %+v", razem.Reverted)
	}
}

// ── Podświetlenie zmian wykonawców ──────────────────────────────────────────

// TestZmianyModeluLicznikNieLiczyOperatora mierzy licznik przy przełączniku:
// ma mówić, ile zmian WYKONAWCY stoi w dokumencie, a nie ile zmian jest w ogóle.
func TestZmianyModeluLicznikNieLiczyOperatora(t *testing.T) {
	uprzaz := dziennikSkutekZmontuj(t)

	// Dwie zmiany wykonawcy i jedna Operatora, każda na innym akapicie.
	trzeci := strings.Index(dziennikSkutekTresc, "Akapit trzeci bez zmian.")
	uprzaz.dziennikSkutekZmienTekst(t, trzeci, trzeci+len("Akapit trzeci bez zmian."),
		"Akapit trzeci od modelu.", "agent-pierwszy")
	drugi := strings.Index(dziennikSkutekTresc, "Akapit drugi do poprawy.")
	uprzaz.dziennikSkutekZmienTekst(t, drugi, drugi+len("Akapit drugi do poprawy."),
		"Akapit drugi od modelu.", "agent-drugi")

	pierwszy := strings.Index(dziennikSkutekTresc, "Akapit pierwszy bez zmian.")
	var zmianaOperatora shared.StudioTextEditResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioTextEdit,
		shared.StudioTextEditRequest{
			DocumentId: uprzaz.dokument, RangeStart: pierwszy,
			RangeEnd: pierwszy + len("Akapit pierwszy bez zmian."),
			Text:     "Akapit pierwszy od Operatora.",
		}, &zmianaOperatora)

	var zestawienie shared.StudioModelChangesListResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioModelChangesList,
		shared.StudioModelChangesListRequest{DocumentId: uprzaz.dokument}, &zestawienie)

	if zestawienie.Summary.Total != 2 {
		t.Errorf("licznik zmian modelu mówi %d, a wykonawcy wnieśli 2 zmiany "+
			"(zmiana Operatora nie ma się w nim liczyć)", zestawienie.Summary.Total)
	}

	// Pomiar niezależny: tyle samo w bazie.
	var wBazie int
	err := uprzaz.baza.QueryRow(
		`SELECT COUNT(*) FROM zmiana_sledzona_studio z
		 JOIN dokument_studio d ON d.id = z.dokument_id
		 WHERE d.identyfikator_zewnetrzny = ? AND z.autor = 'model'`,
		uprzaz.dokument).Scan(&wBazie)
	if err != nil {
		t.Fatalf("nie można policzyć zmian modelu w bazie: %v", err)
	}
	if wBazie != 2 {
		t.Errorf("w bazie stoi %d zmian autora model, a wykonawcy wnieśli 2", wBazie)
	}

	// Filtr per wykonawca: dwóch agentów nie może wyglądać jak jeden.
	agent := "agent-drugi"
	var jedenAgent shared.StudioModelChangesListResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioModelChangesList,
		shared.StudioModelChangesListRequest{DocumentId: uprzaz.dokument, AgentId: &agent},
		&jedenAgent)
	if jedenAgent.Summary.Total != 1 {
		t.Errorf("filtr po wykonawcy %q oddał %d zmian, a ten wykonawca wniósł 1 — "+
			"dwóch agentów wygląda jak jeden", agent, jedenAgent.Summary.Total)
	}
}

// TestZmianyModeluCofniecieWszystkiegoZachowujeOperatora mierzy, że cofnięcie
// wszystkiego, co zrobił model, zachowuje zmiany Operatora naniesione w tym
// czasie, zamiast przywracać wersję sprzed nich.
func TestZmianyModeluCofniecieWszystkiegoZachowujeOperatora(t *testing.T) {
	uprzaz := dziennikSkutekZmontuj(t)

	// Śledzenie włączone, żeby zmiany Operatora też odkładały się wierszem.
	var sledzenie shared.StudioTrackingSetResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioTrackingSet,
		shared.StudioTrackingSetRequest{DocumentId: uprzaz.dokument, Enabled: true}, &sledzenie)

	trzeci := strings.Index(dziennikSkutekTresc, "Akapit trzeci bez zmian.")
	uprzaz.dziennikSkutekZmienTekst(t, trzeci, trzeci+len("Akapit trzeci bez zmian."),
		"Akapit trzeci od modelu.", "agent-pierwszy")

	// Zmiana Operatora naniesiona PO pracy modelu — ta ma zostać.
	pierwszy := strings.Index(dziennikSkutekTresc, "Akapit pierwszy bez zmian.")
	var zmianaOperatora shared.StudioTextEditResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioTextEdit,
		shared.StudioTextEditRequest{
			DocumentId: uprzaz.dokument, RangeStart: pierwszy,
			RangeEnd: pierwszy + len("Akapit pierwszy bez zmian."),
			Text:     "Akapit pierwszy OPERATORA.",
		}, &zmianaOperatora)

	wszystko := true
	var cofniecie shared.StudioModelChangesRevertResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioModelChangesRevert,
		shared.StudioModelChangesRevertRequest{DocumentId: uprzaz.dokument, All: &wszystko},
		&cofniecie)

	if cofniecie.RevertedCount == 0 {
		t.Fatal("cofnięcie wszystkiego nie cofnęło ani jednej zmiany modelu")
	}
	if cofniecie.KeptOperatorChanges == 0 {
		t.Error("odpowiedź nie melduje ani jednej zachowanej zmiany Operatora")
	}

	// Pomiar niezależny — i to jest sedno: praca Operatora ZOSTAJE, praca modelu
	// znika.
	po := uprzaz.blokadaTrescZBazy(t)
	if !strings.Contains(po, "Akapit pierwszy OPERATORA.") {
		t.Errorf("cofnięcie zmian modelu SKASOWAŁO pracę Operatora — zrobiono je "+
			"przywróceniem wersji sprzed; treść: %q", po)
	}
	if strings.Contains(po, "Akapit trzeci od modelu.") {
		t.Errorf("zmiana modelu została w dokumencie po cofnięciu wszystkiego; treść: %q", po)
	}

	// Kopia zapasowa przed czynnością nieodwracalną — cofnięcie wszystkiego nią jest.
	var kopie int
	err := uprzaz.baza.QueryRow(
		`SELECT COUNT(*) FROM kopia_zapasowa_studio k
		 JOIN dokument_studio d ON d.id = k.dokument_id
		 WHERE d.identyfikator_zewnetrzny = ?`, uprzaz.dokument).Scan(&kopie)
	if err != nil {
		t.Fatalf("nie można policzyć kopii zapasowych: %v", err)
	}
	if kopie == 0 {
		t.Error("cofnięcie wszystkich zmian modelu nie założyło kopii zapasowej, " +
			"choć jest czynnością nieodwracalną")
	}
}

// TestZmianyModeluCofniecieWybranychNieRuszaReszty mierzy cofnięcie WYBRANYCH:
// Operator odhacza część zmian i cofa TYLKO je.
func TestZmianyModeluCofniecieWybranychNieRuszaReszty(t *testing.T) {
	uprzaz := dziennikSkutekZmontuj(t)

	trzeci := strings.Index(dziennikSkutekTresc, "Akapit trzeci bez zmian.")
	uprzaz.dziennikSkutekZmienTekst(t, trzeci, trzeci+len("Akapit trzeci bez zmian."),
		"Akapit trzeci od modelu.", "agent-pierwszy")
	drugi := strings.Index(dziennikSkutekTresc, "Akapit drugi do poprawy.")
	uprzaz.dziennikSkutekZmienTekst(t, drugi, drugi+len("Akapit drugi do poprawy."),
		"Akapit drugi od modelu.", "agent-pierwszy")

	var zestawienie shared.StudioModelChangesListResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioModelChangesList,
		shared.StudioModelChangesListRequest{DocumentId: uprzaz.dokument}, &zestawienie)
	if len(zestawienie.Summary.Changes) < 2 {
		t.Fatalf("wykaz zmian modelu niesie %d pozycji, miał 2",
			len(zestawienie.Summary.Changes))
	}

	// Cofamy JEDNĄ — tę, która dotyczy akapitu trzeciego.
	var wybrana string
	for _, zmiana := range zestawienie.Summary.Changes {
		if zmiana.After != nil && strings.Contains(*zmiana.After, "trzeci") {
			wybrana = zmiana.Id
		}
	}
	if wybrana == "" {
		t.Fatalf("nie znaleziono zmiany do odhaczenia w wykazie: %+v", zestawienie.Summary.Changes)
	}

	var cofniecie shared.StudioModelChangesRevertResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioModelChangesRevert,
		shared.StudioModelChangesRevertRequest{
			DocumentId: uprzaz.dokument, ChangeIds: []string{wybrana},
		}, &cofniecie)
	if cofniecie.RevertedCount != 1 {
		t.Errorf("cofnięto %d zmian, a odhaczona była jedna", cofniecie.RevertedCount)
	}

	po := uprzaz.blokadaTrescZBazy(t)
	if strings.Contains(po, "Akapit trzeci od modelu.") {
		t.Errorf("zmiana odhaczona nie została cofnięta; treść: %q", po)
	}
	if !strings.Contains(po, "Akapit drugi od modelu.") {
		t.Errorf("cofnięcie WYBRANEJ zmiany ruszyło zmianę nieodhaczoną; treść: %q", po)
	}
}

// ── Autor na każdej drodze ──────────────────────────────────────────────────

// TestZmianyModeluAutorZapisanyNaKazdejDrodze mierzy założenie przełącznika:
// czynność wykonawcy odkłada ślad podpisany jako wykonawca nawet wtedy, gdy
// żądanie nie niosło pola `author`.
func TestZmianyModeluAutorZapisanyNaKazdejDrodze(t *testing.T) {
	uprzaz := dziennikSkutekZmontuj(t)

	drugi := strings.Index(dziennikSkutekTresc, "Akapit drugi do poprawy.")
	koniec := drugi + len("Akapit drugi do poprawy.")

	// Żądanie bez pola `author` i bez kodu agenta — sam fakt gniazda narzędzi.
	odpowiedz := uprzaz.blokadaWykonajJakoModel(t, shared.CommandStudioTextEdit,
		shared.StudioTextEditRequest{
			DocumentId: uprzaz.dokument, RangeStart: drugi, RangeEnd: koniec,
			Text: "Akapit drugi zmieniony bez podpisu.",
		})
	if odpowiedz.Error != nil {
		t.Fatalf("czynność odmówiła: %s", odpowiedz.Error.Message)
	}

	// Pomiar niezależny: w bazie ślad autora `model`, nie `uzytkownik`.
	var autorZmiany sql.NullString
	err := uprzaz.baza.QueryRow(
		`SELECT z.autor FROM zmiana_sledzona_studio z
		 JOIN dokument_studio d ON d.id = z.dokument_id
		 WHERE d.identyfikator_zewnetrzny = ?
		 ORDER BY z.id DESC LIMIT 1`, uprzaz.dokument).Scan(&autorZmiany)
	if err != nil {
		t.Fatalf("nie można odczytać autora zmiany śledzonej: %v", err)
	}
	if autorZmiany.String != string(shared.StudioAuthorModel) {
		t.Errorf("zmiana wniesiona z gniazda serwera narzędzi zapisała się jako autor %q — "+
			"model omija zapis autora POMINIĘCIEM pola i przez to nie da się jej podświetlić",
			autorZmiany.String)
	}

	var autorCzynnosci sql.NullString
	err = uprzaz.baza.QueryRow(
		`SELECT c.autor_rodzaj FROM czynnosc_dokumentu_studio c
		 JOIN dokument_studio d ON d.id = c.dokument_id
		 WHERE d.identyfikator_zewnetrzny = ?
		 ORDER BY c.kolejnosc DESC LIMIT 1`, uprzaz.dokument).Scan(&autorCzynnosci)
	if err != nil {
		t.Fatalf("nie można odczytać autora czynności dziennika: %v", err)
	}
	if autorCzynnosci.String != string(shared.StudioAuthorModel) {
		t.Errorf("wpis dziennika zapisał autora %q zamiast wykonawcy", autorCzynnosci.String)
	}
}
