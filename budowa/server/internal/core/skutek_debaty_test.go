package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"danacoconsole/shared"
)

// Skutek modułu Roundtable: czy za odpowiedzią stoi zapis, a nie sama koperta.
//
// Wzorzec szkody, którego pilnuje ten plik, wystąpił w tym produkcie: komenda
// meldowała `status: ok` z wykazem, za którym nie stał ani jeden bajt. Dlatego
// żaden sprawdzian tutaj nie kończy się na tym, że odpowiedź jest udana. Każdy
// schodzi do bazy DRUGIM połączeniem — otwartym niezależnie od rdzenia — i liczy
// wiersze albo czyta bajty z magazynu.
//
// Kanałem uczestników jest kanał echo: odsyła treść zapytania porcjami, tak jak
// zrobiłby to model, bez sieci, konta i klucza. Mierzona jest droga rdzenia —
// rejestr, tura, zapis, analiza, głosowanie, wydanie — a nie to, co odpowiada
// konkretny dostawca.

// oknoSprawdzianuDebaty — okno, w którym stoją wszystkie sprawdziany tego pliku.
const oknoSprawdzianuDebaty = "roundtable.debate-panel/sprawdzian"

// bazaSprawdzianu otwiera drugie połączenie z bazą rdzenia. To ono jest
// miarą: liczba wierszy odczytana tędy nie pochodzi od tego samego kodu, który
// je zapisywał.
func bazaSprawdzianu(t *testing.T, katalog string) *sql.DB {
	t.Helper()

	baza, err := sql.Open("sqlite", filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do pomiaru: %v", err)
	}
	t.Cleanup(func() { _ = baza.Close() })
	return baza
}

// wierszy liczy wiersze spełniające warunek — jedna miara na wszystkie tabele.
func wierszy(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) int {
	t.Helper()

	var ile int
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&ile); err != nil {
		t.Fatalf("nie można policzyć wierszy (%s): %v", zapytanie, err)
	}
	return ile
}

// tekstZBazy odczytuje jedną wartość tekstową.
func tekstZBazy(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) string {
	t.Helper()

	var wartosc string
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&wartosc); err != nil {
		t.Fatalf("nie można odczytać wartości (%s): %v", zapytanie, err)
	}
	return wartosc
}

// wpiszKanalEcho zakłada kanał odsyłający treść zapytania i oddaje jego kod.
func wpiszKanalEcho(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	nazwa string) string {
	t.Helper()

	var wynik shared.ChannelAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandChannelAdd, shared.ChannelAddRequest{
		Name:   nazwa,
		Kind:   "api",
		Model:  wskaznik("model-debaty"),
		Config: jsonSurowy(t, map[string]any{"adapter": "echo", "porcja": "64"}),
	}, &wynik)

	if wynik.Channel.Id == "" {
		t.Fatal("channel.add nie oddał identyfikatora kanału")
	}
	return wynik.Channel.Id
}

// dodajUczestnika wpisuje uczestnika debaty i oddaje jego kod.
func dodajUczestnika(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	kanal, persona string) string {
	t.Helper()

	var wynik shared.RoundtableModelAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableModelAdd,
		shared.RoundtableModelAddRequest{
			WindowId: oknoSprawdzianuDebaty, ChannelId: kanal,
			PersonaName: wskaznik(persona),
		}, &wynik)
	return wynik.Participant.Id
}

// przeprowadzTure uruchamia turę i czeka, aż wypowiedzi znajdą się w bazie.
//
// Czekanie idzie po bazie, nie po zegarze: tura biegnie w gorutynie rdzenia
// (odpowiedź komendy wraca przed wypowiedziami), a uśpienie na stałą liczbę
// milisekund byłoby sprawdzianem szybkości maszyny, nie skutku.
func przeprowadzTure(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	baza *sql.DB, pytanie string, ilu int) string {
	t.Helper()

	var wynik shared.RoundtableDebateStartResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableDebateStart,
		shared.RoundtableDebateStartRequest{WindowId: oknoSprawdzianuDebaty, Question: pytanie},
		&wynik)

	const prob = 400
	for proba := 0; proba < prob; proba++ {
		ile := wierszy(t, baza,
			`SELECT COUNT(*) FROM debata_wypowiedz w
			   JOIN debata_tura t ON t.id = w.tura_id
			  WHERE t.identyfikator_zewnetrzny = ? AND w.tresc <> ''`, wynik.Turn.Id)
		if ile >= ilu {
			return wynik.Turn.Id
		}
		odczekajChwile()
	}
	t.Fatalf("tura %s nie zebrała %d wypowiedzi z treścią", wynik.Turn.Id, ilu)
	return ""
}

// zlozDebateSprawdzianu składa rdzeń, kanał, dwóch uczestników i jedną turę.
func zlozDebateSprawdzianu(t *testing.T) (*Zmontowany, context.Context, string, *sql.DB, string) {
	t.Helper()

	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)
	kanal := wpiszKanalEcho(t, zmontowany, zycie, "kanał debaty")
	dodajUczestnika(t, zmontowany, zycie, kanal, "Zwolennik")
	dodajUczestnika(t, zmontowany, zycie, kanal, "Oponent")
	tura := przeprowadzTure(t, zmontowany, zycie, baza,
		"Czy wdrożenie ma ruszyć w tym kwartale?", 2)
	return zmontowany, zycie, katalog, baza, tura
}

// TestZmianaUczestnikaZostajeWBazie wykazuje skutek `roundtable.model.update`:
// waga, rola i oznaczenie kluczowego leżą w wierszu składu, a nie tylko
// w odpowiedzi komendy.
func TestZmianaUczestnikaZostajeWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)
	kanal := wpiszKanalEcho(t, zmontowany, zycie, "kanał składu")
	uczestnik := dodajUczestnika(t, zmontowany, zycie, kanal, "Ekspert")

	rola := shared.RoundtableParticipantRole(shared.RoundtableParticipantRoleJury)
	var wynik shared.RoundtableModelUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableModelUpdate,
		shared.RoundtableModelUpdateRequest{
			WindowId: oknoSprawdzianuDebaty, ParticipantId: uczestnik,
			Role: &rola, Weight: wskaznik(2.5), Key: wskaznik(true),
			RoleDescription: wskaznik("ocenia argumenty"),
		}, &wynik)

	var waga float64
	var kluczowy int
	var zapisanaRola string
	if err := baza.QueryRow(
		`SELECT waga, kluczowy, rola FROM debata_uczestnik WHERE identyfikator_zewnetrzny = ?`,
		uczestnik).Scan(&waga, &kluczowy, &zapisanaRola); err != nil {
		t.Fatalf("nie można odczytać uczestnika z bazy: %v", err)
	}
	if waga != 2.5 || kluczowy != 1 || zapisanaRola != shared.RoundtableParticipantRoleJury {
		t.Fatalf("w bazie leży waga=%v kluczowy=%d rola=%q, a zapisano 2,5 / kluczowy / jury",
			waga, kluczowy, zapisanaRola)
	}

	// Usunięcie ma zdjąć wiersz, a nie tylko zniknąć z odpowiedzi.
	var poUsunieciu shared.RoundtableModelRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableModelRemove,
		shared.RoundtableModelRemoveRequest{
			WindowId: oknoSprawdzianuDebaty, ParticipantId: uczestnik,
		}, &poUsunieciu)
	if ile := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_uczestnik WHERE identyfikator_zewnetrzny = ?`,
		uczestnik); ile != 0 {
		t.Fatalf("po usunięciu uczestnika w bazie zostało %d wierszy", ile)
	}
}

// TestZespolZapisujeSkladIWnosiGoDoOkna wykazuje, że zespół jest kopią składu:
// wiersze zespołu leżą w bazie, a wniesienie go do innego okna daje uczestników
// tego okna.
func TestZespolZapisujeSkladIWnosiGoDoOkna(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)
	kanal := wpiszKanalEcho(t, zmontowany, zycie, "kanał zespołu")
	dodajUczestnika(t, zmontowany, zycie, kanal, "Pierwszy")
	dodajUczestnika(t, zmontowany, zycie, kanal, "Drugi")

	var zapis shared.RoundtableTeamSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableTeamSave,
		shared.RoundtableTeamSaveRequest{WindowId: oknoSprawdzianuDebaty, Name: "Panel ekspercki"},
		&zapis)

	if ile := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_zespol_uczestnik u
		   JOIN debata_zespol z ON z.id = u.zespol_id
		  WHERE z.identyfikator_zewnetrzny = ?`, zapis.Team.Id); ile != 2 {
		t.Fatalf("zespół zapisał %d uczestników, a skład liczył dwóch", ile)
	}

	const drugieOkno = "roundtable.model-panels/drugie"
	var wniesienie shared.RoundtableTeamApplyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableTeamApply,
		shared.RoundtableTeamApplyRequest{WindowId: drugieOkno, TeamId: zapis.Team.Id},
		&wniesienie)

	if ile := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_uczestnik WHERE okno = ?`, drugieOkno); ile != 2 {
		t.Fatalf("po wniesieniu zespołu w oknie %s stoi %d uczestników, a zespół miał dwóch",
			drugieOkno, ile)
	}
}

// TestBibliotekaRolNieJestPusta wykazuje, że role wnosi migracja, a komenda je
// naprawdę czyta — wykaz pusty byłby komendą działającą i bezużyteczną.
func TestBibliotekaRolNieJestPusta(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var wynik shared.RoundtableRoleListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableRoleList,
		shared.RoundtableRoleListRequest{}, &wynik)

	wBazie := wierszy(t, baza, `SELECT COUNT(*) FROM debata_rola`)
	if wBazie == 0 {
		t.Fatal("migracja nie wniosła ani jednej roli debaty")
	}
	if len(wynik.Roles) != wBazie {
		t.Fatalf("komenda oddała %d ról, a w bazie leży %d", len(wynik.Roles), wBazie)
	}
	for _, rola := range wynik.Roles {
		if strings.TrimSpace(rola.SystemPrompt) == "" {
			t.Fatalf("rola %q nie niesie promptu systemowego — nie da się jej nadać uczestnikowi",
				rola.Name)
		}
	}
}

// TestStanDebatyOddajeZapisZBazy wykazuje, że odczyt stanu wraca z zapisu,
// a nie z pamięci połączenia: liczba wypowiedzi w migawce zgadza się z liczbą
// wierszy w bazie.
func TestStanDebatyOddajeZapisZBazy(t *testing.T) {
	zmontowany, zycie, _, baza, tura := zlozDebateSprawdzianu(t)

	var wynik shared.RoundtableDebateGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableDebateGet,
		shared.RoundtableDebateGetRequest{WindowId: oknoSprawdzianuDebaty}, &wynik)

	wBazie := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_wypowiedz w
		   JOIN debata_tura t ON t.id = w.tura_id
		  WHERE t.okno = ?`, oknoSprawdzianuDebaty)
	if len(wynik.Snapshot.Statements) != wBazie {
		t.Fatalf("migawka niesie %d wypowiedzi, a w bazie leży %d",
			len(wynik.Snapshot.Statements), wBazie)
	}
	if len(wynik.Snapshot.Participants) != 2 {
		t.Fatalf("migawka niesie %d uczestników, a skład liczył dwóch",
			len(wynik.Snapshot.Participants))
	}
	if len(wynik.Snapshot.Turns) == 0 || wynik.Snapshot.Turns[0].Id != tura {
		t.Fatal("migawka nie niesie tury, która właśnie przebiegła")
	}
}

// TestWatekBocznyZakladaTureZOdpowiedzia wykazuje skutek
// `roundtable.debate.followup`: w bazie stoi nowa tura wskazująca turę
// nadrzędną, a w niej wypowiedź adresata z treścią.
func TestWatekBocznyZakladaTureZOdpowiedzia(t *testing.T) {
	zmontowany, zycie, _, baza, tura := zlozDebateSprawdzianu(t)

	var skladu shared.RoundtableModelListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableModelList,
		shared.RoundtableModelListRequest{WindowId: oknoSprawdzianuDebaty}, &skladu)

	var wynik shared.RoundtableDebateFollowupResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableDebateFollowup,
		shared.RoundtableDebateFollowupRequest{
			WindowId: oknoSprawdzianuDebaty, ParticipantId: skladu.Participants[0].Id,
			Question: "Na czym opierasz koszt wdrożenia?", TurnId: &tura,
		}, &wynik)

	if strings.TrimSpace(wynik.Statement.Content) == "" {
		t.Fatal("wątek boczny oddał wypowiedź bez treści")
	}
	if ile := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_tura WHERE okno = ? AND tura_nadrzedna = ?`,
		oknoSprawdzianuDebaty, tura); ile != 1 {
		t.Fatalf("w bazie stoi %d tur pobocznych przy turze %s, a założono jedną", ile, tura)
	}
	tresc := tekstZBazy(t, baza,
		`SELECT tresc FROM debata_wypowiedz WHERE identyfikator_zewnetrzny = ?`,
		wynik.Statement.Id)
	if strings.TrimSpace(tresc) == "" {
		t.Fatal("wypowiedź wątku bocznego leży w bazie bez treści")
	}
}

// TestRegeneracjaPodnosiNumerRedakcji wykazuje, że powtórzenie zastępuje
// wypowiedź, a nie dopisuje drugiej: liczba wierszy zostaje, numer redakcji
// rośnie.
func TestRegeneracjaPodnosiNumerRedakcji(t *testing.T) {
	zmontowany, zycie, _, baza, tura := zlozDebateSprawdzianu(t)

	var stan shared.RoundtableDebateGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableDebateGet,
		shared.RoundtableDebateGetRequest{WindowId: oknoSprawdzianuDebaty, TurnId: &tura}, &stan)
	przed := len(stan.Snapshot.Statements)
	kod := stan.Snapshot.Statements[0].Id

	var wynik shared.RoundtableStatementRegenerateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableStatementRegenerate,
		shared.RoundtableStatementRegenerateRequest{
			WindowId: oknoSprawdzianuDebaty, StatementId: kod,
		}, &wynik)

	var redakcja int
	if err := baza.QueryRow(
		`SELECT redakcja FROM debata_wypowiedz WHERE identyfikator_zewnetrzny = ?`,
		kod).Scan(&redakcja); err != nil {
		t.Fatalf("nie można odczytać redakcji wypowiedzi: %v", err)
	}
	if redakcja != 2 {
		t.Fatalf("po regeneracji numer redakcji wynosi %d, a miał wynieść 2", redakcja)
	}
	po := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_wypowiedz w JOIN debata_tura t ON t.id = w.tura_id
		  WHERE t.identyfikator_zewnetrzny = ?`, tura)
	if po != przed {
		t.Fatalf("regeneracja zmieniła liczbę wypowiedzi z %d na %d — miała zastąpić, nie dopisać",
			przed, po)
	}
}

// TestWariantTuryWskazujeTureRozgaleziana wykazuje, że obie gałęzie zostają
// w zapisie.
func TestWariantTuryWskazujeTureRozgaleziana(t *testing.T) {
	zmontowany, zycie, _, baza, tura := zlozDebateSprawdzianu(t)

	var wynik shared.RoundtableDebateBranchResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableDebateBranch,
		shared.RoundtableDebateBranchRequest{
			WindowId: oknoSprawdzianuDebaty, TurnId: tura, Label: wskaznik("wariant ostrożny"),
		}, &wynik)

	nadrzedna := tekstZBazy(t, baza,
		`SELECT tura_nadrzedna FROM debata_tura WHERE identyfikator_zewnetrzny = ?`,
		wynik.Turn.Id)
	if nadrzedna != tura {
		t.Fatalf("wariant wskazuje turę %q, a rozgałęziano %q", nadrzedna, tura)
	}
	if ile := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_tura WHERE identyfikator_zewnetrzny = ?`, tura); ile != 1 {
		t.Fatal("rozgałęzienie usunęło turę źródłową — obie gałęzie mają zostać")
	}
}

// TestAnalizaArgumentowZapisujeGrafWBazie wykazuje skutek analizy: po
// `roundtable.analysis.run` w bazie leżą węzły grafu i ustalenia, a nie sama
// odpowiedź komendy.
func TestAnalizaArgumentowZapisujeGrafWBazie(t *testing.T) {
	zmontowany, zycie, _, baza, _ := zlozDebateSprawdzianu(t)

	var wynik shared.RoundtableAnalysisRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableAnalysisRun,
		shared.RoundtableAnalysisRunRequest{
			WindowId: oknoSprawdzianuDebaty, Kind: shared.RoundtableAnalysisKindArgumentMining,
		}, &wynik)

	if len(wynik.Findings) == 0 {
		t.Fatal("analiza oddała pusty wykaz ustaleń")
	}
	wezlow := wierszy(t, baza, `SELECT COUNT(*) FROM debata_wezel WHERE okno = ?`,
		oknoSprawdzianuDebaty)
	if wezlow == 0 {
		t.Fatal("po wydobyciu argumentów w bazie nie leży ani jeden węzeł grafu")
	}
	ustalen := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_ustalenie WHERE okno = ? AND rodzaj = ?`,
		oknoSprawdzianuDebaty, shared.RoundtableAnalysisKindArgumentMining)
	if ustalen != len(wynik.Findings) {
		t.Fatalf("komenda oddała %d ustaleń, a w bazie leży %d", len(wynik.Findings), ustalen)
	}
	if wynik.Graph == nil || len(wynik.Graph.Nodes) != wezlow {
		t.Fatal("odpowiedź nie niesie grafu zgodnego z zapisem w bazie")
	}

	// Oznaczenie węzła jako kluczowego ma zostać w bazie — to wybór Operatora.
	var oznaczenie shared.RoundtableArgumentPinResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableArgumentPin,
		shared.RoundtableArgumentPinRequest{
			WindowId: oknoSprawdzianuDebaty, NodeId: wynik.Graph.Nodes[0].Id, Pinned: true,
		}, &oznaczenie)
	if ile := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_wezel WHERE identyfikator_zewnetrzny = ? AND kluczowy = 1`,
		wynik.Graph.Nodes[0].Id); ile != 1 {
		t.Fatal("oznaczenie argumentu kluczowego nie zostało w bazie")
	}
}

// TestKatalogBledowZawezaSieDoWskazanychKodow wykazuje skutek
// `roundtable.fallacy.catalog.set`: zawężenie leży w bazie i widać je w odczycie.
func TestKatalogBledowZawezaSieDoWskazanychKodow(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var przed shared.RoundtableFallacyCatalogGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableFallacyCatalogGet,
		shared.RoundtableFallacyCatalogGetRequest{}, &przed)
	if len(przed.Definitions) == 0 {
		t.Fatal("migracja nie wniosła ani jednej pozycji katalogu błędów")
	}
	for _, pozycja := range przed.Definitions {
		if !pozycja.Enabled {
			t.Fatalf("okno bez zawężenia ma wykrywać wszystko, a %s jest wyłączony", pozycja.Code)
		}
	}

	var po shared.RoundtableFallacyCatalogSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableFallacyCatalogSet,
		shared.RoundtableFallacyCatalogSetRequest{
			WindowId: oknoSprawdzianuDebaty, Codes: []string{"adHominem", "strawman"},
		}, &po)

	if ile := wierszy(t, baza, `SELECT COUNT(*) FROM debata_katalog_okna WHERE okno = ?`,
		oknoSprawdzianuDebaty); ile != 2 {
		t.Fatalf("w bazie leży %d wierszy zawężenia katalogu, a wskazano dwa kody", ile)
	}
	wykrywanych := 0
	for _, pozycja := range po.Definitions {
		if pozycja.Enabled {
			wykrywanych++
		}
	}
	if wykrywanych != 2 {
		t.Fatalf("po zawężeniu okno wykrywa %d błędów, a wskazano dwa", wykrywanych)
	}

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandRoundtableFallacyCatalogSet,
		shared.RoundtableFallacyCatalogSetRequest{
			WindowId: oknoSprawdzianuDebaty, Codes: []string{"bledNieistniejacy"},
		})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("kod spoza katalogu ma być odmową validation_failed, a jest %s", odmowa.Code)
	}
}

// TestGlosowanieAprobacyjneWylaniaZwyciezceZOddanychGlosow wykazuje skutek
// głosowania: głosy leżą w bazie, a wynik agregacji wynika z nich, nie z życzeń.
func TestGlosowanieAprobacyjneWylaniaZwyciezceZOddanychGlosow(t *testing.T) {
	zmontowany, zycie, _, baza, _ := zlozDebateSprawdzianu(t)

	var otwarcie shared.RoundtableVoteStartResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableVoteStart,
		shared.RoundtableVoteStartRequest{
			WindowId: oknoSprawdzianuDebaty,
			Method:   shared.RoundtableVoteMethodApproval,
			Options:  []string{"Ruszamy w tym kwartale", "Przekładamy na następny"},
		}, &otwarcie)
	if len(otwarcie.Vote.Options) != 2 {
		t.Fatalf("głosowanie otwarto z %d wariantami, a podano dwa", len(otwarcie.Vote.Options))
	}
	pierwszy := otwarcie.Vote.Options[0].Id

	for _, wyborca := range []string{"operator", "jury-a", "jury-b"} {
		var glos shared.RoundtableVoteCastResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableVoteCast,
			shared.RoundtableVoteCastRequest{
				WindowId: oknoSprawdzianuDebaty, VoteId: otwarcie.Vote.Id,
				VoterId: wyborca, Approvals: []string{pierwszy},
			}, &glos)
	}

	if ile := wierszy(t, baza, `SELECT COUNT(*) FROM debata_glos WHERE glosowanie = ?`,
		otwarcie.Vote.Id); ile != 3 {
		t.Fatalf("w bazie leży %d głosów, a oddano trzy", ile)
	}

	var wynik shared.RoundtableVoteGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableVoteGet,
		shared.RoundtableVoteGetRequest{
			WindowId: oknoSprawdzianuDebaty, VoteId: &otwarcie.Vote.Id,
		}, &wynik)

	if wynik.Result == nil {
		t.Fatal("głosowanie z trzema głosami nie oddało wyniku")
	}
	if wynik.Result.Ballots != 3 {
		t.Fatalf("wynik liczy %d głosów, a oddano trzy", wynik.Result.Ballots)
	}
	if wynik.Result.WinnerOptionId == nil || *wynik.Result.WinnerOptionId != pierwszy {
		t.Fatal("wariant z trzema aprobatami nie został zwycięzcą")
	}
	rozklad := map[string]float64{}
	if err := json.Unmarshal(wynik.Result.Tally, &rozklad); err != nil {
		t.Fatalf("rozkładu głosów nie da się odczytać: %v", err)
	}
	if rozklad[pierwszy] != 3 {
		t.Fatalf("rozkład przypisuje zwycięzcy %v aprobat, a padły trzy", rozklad[pierwszy])
	}

	// Powtórny głos tego samego wyborcy zastępuje poprzedni, a nie dokłada się.
	var powtorny shared.RoundtableVoteCastResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableVoteCast,
		shared.RoundtableVoteCastRequest{
			WindowId: oknoSprawdzianuDebaty, VoteId: otwarcie.Vote.Id,
			VoterId: "operator", Approvals: []string{otwarcie.Vote.Options[1].Id},
		}, &powtorny)
	if ile := wierszy(t, baza, `SELECT COUNT(*) FROM debata_glos WHERE glosowanie = ?`,
		otwarcie.Vote.Id); ile != 3 {
		t.Fatalf("po zmianie zdania w bazie leży %d głosów, a wyborców było trzech", ile)
	}
}

// TestMetodaEliminacyjnaLiczyRundy wykazuje, że IRV rozstrzyga przeniesieniem
// głosów, a nie samą liczbą pierwszych wskazań.
func TestMetodaEliminacyjnaLiczyRundy(t *testing.T) {
	zmontowany, zycie, _, _, _ := zlozDebateSprawdzianu(t)

	var otwarcie shared.RoundtableVoteStartResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableVoteStart,
		shared.RoundtableVoteStartRequest{
			WindowId: oknoSprawdzianuDebaty, Method: shared.RoundtableVoteMethodIrv,
			Options: []string{"A", "B", "C"},
		}, &otwarcie)
	a := otwarcie.Vote.Options[0].Id
	b := otwarcie.Vote.Options[1].Id
	c := otwarcie.Vote.Options[2].Id

	// Dwa głosy na A, dwa na B, jeden na C z drugim wskazaniem na B.
	// Po odpadnięciu C wygrywa B — mimo że w pierwszej rundzie był remis.
	rankingi := map[string][]string{
		"w1": {a, c, b}, "w2": {a, c, b},
		"w3": {b, c, a}, "w4": {b, c, a},
		"w5": {c, b, a},
	}
	for wyborca, ranking := range rankingi {
		var glos shared.RoundtableVoteCastResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableVoteCast,
			shared.RoundtableVoteCastRequest{
				WindowId: oknoSprawdzianuDebaty, VoteId: otwarcie.Vote.Id,
				VoterId: wyborca, Ranking: ranking,
			}, &glos)
	}

	var wynik shared.RoundtableVoteGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableVoteGet,
		shared.RoundtableVoteGetRequest{
			WindowId: oknoSprawdzianuDebaty, VoteId: &otwarcie.Vote.Id,
		}, &wynik)
	if wynik.Result == nil || wynik.Result.WinnerOptionId == nil {
		t.Fatal("głosowanie rankingowe nie wyłoniło zwycięzcy mimo pięciu głosów")
	}
	if *wynik.Result.WinnerOptionId != b {
		t.Fatalf("po eliminacji wygrał wariant %s, a głosy przenoszą się na %s",
			*wynik.Result.WinnerOptionId, b)
	}
}

// TestOcenaParamiPrzesuwaRanking wykazuje, że wskazanie wypowiedzi bardziej
// przekonującej zostawia ślad w rankingu, a nie tylko w wykazie ocen.
func TestOcenaParamiPrzesuwaRanking(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)
	kanal := wpiszKanalEcho(t, zmontowany, zycie, "kanał oceny")
	dodajUczestnika(t, zmontowany, zycie, kanal, "Zwolennik")
	dodajUczestnika(t, zmontowany, zycie, kanal, "Oponent")
	tura := przeprowadzTure(t, zmontowany, zycie, baza, "Czy wariant A jest lepszy?", 2)

	var stan shared.RoundtableDebateGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableDebateGet,
		shared.RoundtableDebateGetRequest{WindowId: oknoSprawdzianuDebaty, TurnId: &tura}, &stan)

	var ocena shared.RoundtableRatingSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableRatingSet,
		shared.RoundtableRatingSetRequest{
			WindowId: oknoSprawdzianuDebaty, Kind: shared.RoundtableRatingKindPairwise,
			TargetStatementId: wskaznik(stan.Snapshot.Statements[1].Id),
			PreferredId:       wskaznik(stan.Snapshot.Statements[0].Id),
		}, &ocena)

	if ile := wierszy(t, baza, `SELECT COUNT(*) FROM debata_ocena WHERE okno = ?`,
		oknoSprawdzianuDebaty); ile != 1 {
		t.Fatalf("w bazie leży %d ocen, a postawiono jedną", ile)
	}
	if ile := wierszy(t, baza,
		`SELECT COUNT(*) FROM debata_ranking WHERE algorytm = 'elo' AND pojedynki > 0`); ile < 2 {
		t.Fatalf("po pojedynku w rankingu Elo stoi %d pozycji z pojedynkami, a miały być dwie", ile)
	}

	var ranking shared.RoundtableLeaderboardGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableLeaderboardGet,
		shared.RoundtableLeaderboardGetRequest{
			Scope: shared.RoundtableLeaderboardScopeEnvironment,
		}, &ranking)
	if len(ranking.Entries) < 2 {
		t.Fatalf("ranking oddał %d pozycji, a w pojedynku stanęły dwie tożsamości",
			len(ranking.Entries))
	}
	if ranking.Entries[0].Rating <= ranking.Entries[len(ranking.Entries)-1].Rating {
		t.Fatal("po rozstrzygniętym pojedynku punktacje obu tożsamości są równe")
	}
}

// TestRubrykaOdmawiaWagomNiesumujacymSieDoJednosci pilnuje odmowy nazwanej:
// rubryka o innej sumie wag dawałaby wynik nieporównywalny z żadnym innym.
func TestRubrykaOdmawiaWagomNiesumujacymSieDoJednosci(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandRoundtableRubricSet,
		shared.RoundtableRubricSetRequest{
			Name: "Rzetelność i wykonalność",
			Criteria: jsonSurowy(t, []map[string]any{
				{"name": "rzetelność", "weight": 0.4},
				{"name": "wykonalność", "weight": 0.4},
			}),
		})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("odmowa ma nieść kod validation_failed, a niesie %s", odmowa.Code)
	}
	if !strings.Contains(odmowa.Message, "jedności") {
		t.Fatalf("odmowa nie mówi, że wagi mają sumować się do jedności: %q", odmowa.Message)
	}
}

// TestWerdyktSedziegoZostajeWBazie wykazuje skutek `roundtable.judge.run`:
// werdykty leżą w bazie wraz z uzasadnieniem wypowiedzianym przez sędziego.
func TestWerdyktSedziegoZostajeWBazie(t *testing.T) {
	zmontowany, zycie, _, baza, tura := zlozDebateSprawdzianu(t)

	var rubryka shared.RoundtableRubricSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableRubricSet,
		shared.RoundtableRubricSetRequest{
			Name: "Rzetelność, wykonalność, ryzyko",
			Criteria: jsonSurowy(t, []map[string]any{
				{"name": "rzetelność", "weight": 0.4},
				{"name": "wykonalność", "weight": 0.3},
				{"name": "ryzyko", "weight": 0.3},
			}),
		}, &rubryka)
	if len(rubryka.Rubric.Criteria) != 3 {
		t.Fatalf("rubryka zapisała %d kryteriów, a podano trzy", len(rubryka.Rubric.Criteria))
	}

	var skladu shared.RoundtableModelListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableModelList,
		shared.RoundtableModelListRequest{WindowId: oknoSprawdzianuDebaty}, &skladu)

	var werdykty shared.RoundtableJudgeRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableJudgeRun,
		shared.RoundtableJudgeRunRequest{
			WindowId: oknoSprawdzianuDebaty, RubricId: rubryka.Rubric.Id,
			JudgeParticipantIds: []string{skladu.Participants[0].Id}, TurnId: &tura,
		}, &werdykty)

	if len(werdykty.Judgements) == 0 {
		t.Fatal("ocena sędziowska nie wystawiła ani jednego werdyktu")
	}
	wBazie := wierszy(t, baza, `SELECT COUNT(*) FROM debata_werdykt WHERE okno = ?`,
		oknoSprawdzianuDebaty)
	if wBazie != len(werdykty.Judgements) {
		t.Fatalf("komenda oddała %d werdyktów, a w bazie leży %d",
			len(werdykty.Judgements), wBazie)
	}
	for _, werdykt := range werdykty.Judgements {
		if strings.TrimSpace(werdykt.Rationale) == "" {
			t.Fatal("werdykt bez uzasadnienia — rdzeń nie ma prawa ocenić za sędziego")
		}
	}
}

// TestMacierzDecyzyjnaLiczyWynikWazony wykazuje, że wynik wariantu wynika
// z ocen i wag, a nie z pola przepisanego z żądania.
func TestMacierzDecyzyjnaLiczyWynikWazony(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)

	var zapis shared.RoundtableDecisionMatrixSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableDecisionMatrixSet,
		shared.RoundtableDecisionMatrixSetRequest{
			WindowId: oknoSprawdzianuDebaty, Name: "Wybór wariantu wdrożenia",
			Criteria: jsonSurowy(t, []map[string]any{
				{"name": "koszt", "weight": 0.5},
				{"name": "ryzyko", "weight": 0.5},
			}),
			Options: jsonSurowy(t, []map[string]any{
				{"label": "Wariant A", "scores": map[string]float64{"koszt": 8, "ryzyko": 6}},
				{"label": "Wariant B", "scores": map[string]float64{"koszt": 4, "ryzyko": 4}},
			}),
		}, &zapis)

	if ile := wierszy(t, baza, `SELECT COUNT(*) FROM debata_macierz_wariant WHERE macierz = ?`,
		zapis.Matrix.Id); ile != 2 {
		t.Fatalf("w bazie leży %d wariantów macierzy, a podano dwa", ile)
	}
	if len(zapis.Matrix.Options) != 2 {
		t.Fatalf("macierz oddała %d wariantów", len(zapis.Matrix.Options))
	}
	if zapis.Matrix.Options[0].Total == nil || *zapis.Matrix.Options[0].Total != 7 {
		t.Fatalf("wariant A ma wynik ważony 0,5·8 + 0,5·6 = 7, a oddano %v",
			zapis.Matrix.Options[0].Total)
	}
	if zapis.Matrix.Options[1].Total == nil || *zapis.Matrix.Options[1].Total != 4 {
		t.Fatalf("wariant B ma wynik ważony 4, a oddano %v", zapis.Matrix.Options[1].Total)
	}

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandRoundtableDecisionMatrixSet,
		shared.RoundtableDecisionMatrixSetRequest{
			WindowId: oknoSprawdzianuDebaty, Name: "Macierz z obcym kryterium",
			Criteria: jsonSurowy(t, []map[string]any{{"name": "koszt", "weight": 1}}),
			Options: jsonSurowy(t, []map[string]any{
				{"label": "Wariant", "scores": map[string]float64{"czas": 5}},
			}),
		})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("ocena w nieistniejącym kryterium ma być odmową, a kod to %s", odmowa.Code)
	}
}

// TestStanowiskoRedagowaneNieWracaDoZapisuTur wykazuje najważniejszy skutek
// Consensus Panelu: treść nadana przez Operatora zostaje, a kolejny odczyt jej
// nie nadpisuje złożeniem z tur.
func TestStanowiskoRedagowaneNieWracaDoZapisuTur(t *testing.T) {
	zmontowany, zycie, _, baza, _ := zlozDebateSprawdzianu(t)

	const tresc = "Wdrożenie rusza w tym kwartale, z rezerwą na testy wydajności."
	var zapis shared.RoundtableConsensusSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConsensusSet,
		shared.RoundtableConsensusSetRequest{
			WindowId: oknoSprawdzianuDebaty, Content: tresc,
			Context:  wskaznik("Debata dwóch person nad terminem wdrożenia."),
			Accepted: wskaznik(true),
		}, &zapis)

	// Odczyt stanowiska po redakcji nie ma prawa jej nadpisać.
	var odczyt shared.RoundtableConsensusGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConsensusGet,
		shared.RoundtableConsensusGetRequest{WindowId: oknoSprawdzianuDebaty}, &odczyt)

	wBazie := tekstZBazy(t, baza,
		`SELECT tresc FROM debata_stanowisko WHERE okno = ? AND tura = ''`,
		oknoSprawdzianuDebaty)
	if wBazie != tresc {
		t.Fatalf("w bazie leży treść %q, a Operator zapisał %q", wBazie, tresc)
	}
	if odczyt.Consensus.Content == nil || *odczyt.Consensus.Content != tresc {
		t.Fatal("odczyt stanowiska nadpisał redakcję Operatora złożeniem z tur")
	}

	// Druga redakcja podnosi wersję i odkłada osobny wpis.
	const poprawiona = tresc + " Termin potwierdzono po turze drugiej."
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConsensusSet,
		shared.RoundtableConsensusSetRequest{WindowId: oknoSprawdzianuDebaty, Content: poprawiona},
		&zapis)

	var wersje shared.RoundtableConsensusVersionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConsensusVersionList,
		shared.RoundtableConsensusVersionListRequest{WindowId: oknoSprawdzianuDebaty}, &wersje)
	if len(wersje.Versions) < 2 {
		t.Fatalf("po dwóch redakcjach jest %d wersji stanowiska", len(wersje.Versions))
	}
	if wersje.Versions[len(wersje.Versions)-1].Content != poprawiona {
		t.Fatal("ostatnia wersja nie niesie treści z drugiej redakcji")
	}
}

// TestZdanieOdrebneObnizaPoparcieStanowiska wykazuje, że podpis pod zdaniem
// odrębnym zmienia liczbę, a nie tylko dokłada wiersz.
func TestZdanieOdrebneObnizaPoparcieStanowiska(t *testing.T) {
	zmontowany, zycie, _, baza, _ := zlozDebateSprawdzianu(t)

	var zapis shared.RoundtableConsensusSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConsensusSet,
		shared.RoundtableConsensusSetRequest{
			WindowId: oknoSprawdzianuDebaty, Content: "Wdrożenie rusza w tym kwartale.",
		}, &zapis)
	if zapis.Consensus.Support == nil || *zapis.Consensus.Support != 1 {
		t.Fatal("stanowisko bez zdań odrębnych ma poparcie pełne")
	}

	var skladu shared.RoundtableModelListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableModelList,
		shared.RoundtableModelListRequest{WindowId: oknoSprawdzianuDebaty}, &skladu)

	var zdanie shared.RoundtableConsensusMinoritySetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConsensusMinoritySet,
		shared.RoundtableConsensusMinoritySetRequest{
			WindowId: oknoSprawdzianuDebaty, ConsensusId: zapis.Consensus.Id,
			ParticipantId: skladu.Participants[1].Id,
			Content:       "Termin jest nierealny bez testów wydajności.",
		}, &zdanie)

	if ile := wierszy(t, baza, `SELECT COUNT(*) FROM debata_zdanie_odrebne WHERE stanowisko = ?`,
		zapis.Consensus.Id); ile != 1 {
		t.Fatalf("w bazie leży %d zdań odrębnych, a podpisano jedno", ile)
	}

	var po shared.RoundtableConsensusSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConsensusSet,
		shared.RoundtableConsensusSetRequest{
			WindowId: oknoSprawdzianuDebaty, Content: "Wdrożenie rusza w tym kwartale.",
		}, &po)
	if po.Consensus.Support == nil || *po.Consensus.Support >= 1 {
		t.Fatalf("po zdaniu odrębnym poparcie wynosi %v, a miało spaść poniżej jedności",
			po.Consensus.Support)
	}
	if len(po.Consensus.Minority) != 1 {
		t.Fatal("stanowisko nie niesie podpisanego zdania odrębnego")
	}
}

// TestTranskryptMarkdownLezyWMagazynie wykazuje skutek wydania: pod odwołaniem
// leży plik, a w nim wypowiedzi debaty.
func TestTranskryptMarkdownLezyWMagazynie(t *testing.T) {
	zmontowany, zycie, katalog, _, _ := zlozDebateSprawdzianu(t)

	var wynik shared.RoundtableTranscriptExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableTranscriptExport,
		shared.RoundtableTranscriptExportRequest{
			WindowId: oknoSprawdzianuDebaty, Format: shared.RoundtableTranscriptFormatMarkdown,
		}, &wynik)

	if wynik.Uri == nil {
		t.Fatal("wydanie transkryptu nie oddało odwołania do treści")
	}
	bajty := bajtyPodOdwolaniem(t, katalog, *wynik.Uri)
	if !strings.Contains(string(bajty), "Zwolennik") {
		t.Fatalf("transkrypt nie niesie podpisu uczestnika; początek: %.200s", bajty)
	}
	if !strings.Contains(string(bajty), "Tura 1") {
		t.Fatal("transkrypt nie niesie nagłówka tury")
	}
}

// TestTranskryptPdfMaStronyDokumentu wykazuje, że wydanie PDF jest dokumentem,
// a nie tekstem z rozszerzeniem: strony liczy się z bajtów pliku.
func TestTranskryptPdfMaStronyDokumentu(t *testing.T) {
	zmontowany, zycie, katalog, _, _ := zlozDebateSprawdzianu(t)

	var wynik shared.RoundtableTranscriptExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableTranscriptExport,
		shared.RoundtableTranscriptExportRequest{
			WindowId: oknoSprawdzianuDebaty, Format: shared.RoundtableTranscriptFormatPdf,
		}, &wynik)

	if wynik.Uri == nil {
		t.Fatal("wydanie dokumentu nie oddało odwołania do treści")
	}
	sciezka := sciezkaWMagazynie(katalog, *wynik.Uri)
	if strony := stronWyniku(t, sciezka); strony < 1 {
		t.Fatalf("dokument transkryptu ma %d stron", strony)
	}
}

// TestWydanieGrafuDajePlikWKazdymFormacie wykazuje, że każdy z sześciu formatów
// zostawia w magazynie plik z treścią właściwą dla formatu.
func TestWydanieGrafuDajePlikWKazdymFormacie(t *testing.T) {
	zmontowany, zycie, katalog, _, _ := zlozDebateSprawdzianu(t)

	var analiza shared.RoundtableAnalysisRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableAnalysisRun,
		shared.RoundtableAnalysisRunRequest{
			WindowId: oknoSprawdzianuDebaty, Kind: shared.RoundtableAnalysisKindArgumentMining,
		}, &analiza)

	znaczniki := map[shared.RoundtableArgumentFormat]string{
		shared.RoundtableArgumentFormatDot:     "digraph",
		shared.RoundtableArgumentFormatGraphml: "graphml",
		shared.RoundtableArgumentFormatArgdown: "[",
		shared.RoundtableArgumentFormatAif:     "nodeID",
		shared.RoundtableArgumentFormatSvg:     "<svg",
		shared.RoundtableArgumentFormatPng:     "\x89PNG",
	}
	for format, znacznik := range znaczniki {
		var wynik shared.RoundtableArgumentExportResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableArgumentExport,
			shared.RoundtableArgumentExportRequest{
				WindowId: oknoSprawdzianuDebaty, Format: format,
			}, &wynik)
		if wynik.Uri == nil {
			t.Fatalf("wydanie grafu w formacie %s nie oddało odwołania", format)
		}
		bajty := bajtyPodOdwolaniem(t, katalog, *wynik.Uri)
		if !strings.Contains(string(bajty), znacznik) {
			t.Fatalf("plik grafu w formacie %s nie niesie znacznika %q", format, znacznik)
		}
	}
}

// TestZgodaKlastryIZbieznoscLiczaSieZZapisu wykazuje, że wskaźniki panelu biorą
// się z wypowiedzi, a nie z wartości zastępczych.
func TestZgodaKlastryIZbieznoscLiczaSieZZapisu(t *testing.T) {
	zmontowany, zycie, _, _, _ := zlozDebateSprawdzianu(t)

	var zgodnosc shared.RoundtableAgreementGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableAgreementGet,
		shared.RoundtableAgreementGetRequest{WindowId: oknoSprawdzianuDebaty}, &zgodnosc)
	if len(zgodnosc.Matrix) == 0 {
		t.Fatal("macierz zgodności jest pusta mimo dwóch stanowisk w turze")
	}
	for _, komorka := range zgodnosc.Matrix {
		if komorka.Agreement < 0 || komorka.Agreement > 1 {
			t.Fatalf("zgodność pary leży poza zakresem od zera do jedności: %v", komorka.Agreement)
		}
	}

	var klastry shared.RoundtableClusterGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableClusterGet,
		shared.RoundtableClusterGetRequest{WindowId: oknoSprawdzianuDebaty}, &klastry)
	wKlastrach := 0
	for _, klaster := range klastry.Clusters {
		wKlastrach += len(klaster.ParticipantIds)
	}
	if wKlastrach != 2 {
		t.Fatalf("klastry obejmują %d uczestników, a w debacie stanęło dwóch", wKlastrach)
	}

	var zbieznosc shared.RoundtableConvergenceGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConvergenceGet,
		shared.RoundtableConvergenceGetRequest{WindowId: oknoSprawdzianuDebaty}, &zbieznosc)
	if len(zbieznosc.Points) != 1 {
		t.Fatalf("zbieżność zmierzono w %d turach, a przebiegła jedna", len(zbieznosc.Points))
	}
	if zbieznosc.Points[0].Convergence < 0 || zbieznosc.Points[0].Convergence > 1 {
		t.Fatalf("zbieżność leży poza zakresem: %v", zbieznosc.Points[0].Convergence)
	}
}

// TestSzablonModeracjiZapisujeFormatIKolejnosc wykazuje, że szablon zapamiętuje
// przebieg debaty, a nie samą nazwę.
func TestSzablonModeracjiZapisujeFormatIKolejnosc(t *testing.T) {
	zmontowany, zycie, _, baza, _ := zlozDebateSprawdzianu(t)

	var zapis shared.RoundtableModerationTemplateSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableModerationTemplateSave,
		shared.RoundtableModerationTemplateSaveRequest{
			WindowId: oknoSprawdzianuDebaty, Name: "Panel dwóch person",
		}, &zapis)

	kolejnosc := tekstZBazy(t, baza,
		`SELECT kolejnosc_glosu FROM debata_szablon WHERE identyfikator_zewnetrzny = ?`,
		zapis.Template.Id)
	if len(strings.Split(strings.TrimSpace(kolejnosc), "\n")) != 2 {
		t.Fatalf("szablon zapamiętał kolejność %q, a w składzie stało dwóch uczestników", kolejnosc)
	}

	var wykaz shared.RoundtableModerationTemplateListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableModerationTemplateList,
		shared.RoundtableModerationTemplateListRequest{}, &wykaz)
	if len(wykaz.Templates) != 1 {
		t.Fatalf("wykaz szablonów niesie %d pozycji, a zapisano jedną", len(wykaz.Templates))
	}
	if len(wykaz.Templates[0].SpeakingOrder) != 2 {
		t.Fatal("szablon w wykazie zgubił kolejność głosu")
	}
}

// TestPrzekazanieStanowiskaZakladaArtefakt wykazuje, że przekazanie do modułu
// docelowego zostawia plik, a nie samą adnotację.
func TestPrzekazanieStanowiskaZakladaArtefakt(t *testing.T) {
	zmontowany, zycie, katalog, baza, _ := zlozDebateSprawdzianu(t)

	var zapis shared.RoundtableConsensusSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConsensusSet,
		shared.RoundtableConsensusSetRequest{
			WindowId: oknoSprawdzianuDebaty, Content: "Wdrożenie rusza w tym kwartale.",
			Consequences: wskaznik("Zespół testowy pracuje równolegle."),
		}, &zapis)

	var przekazanie shared.RoundtableConsensusHandoffResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableConsensusHandoff,
		shared.RoundtableConsensusHandoffRequest{
			WindowId: oknoSprawdzianuDebaty, ConsensusId: zapis.Consensus.Id,
			Target: shared.RoundtableHandoffTargetStudio, IncludeTranscript: wskaznik(true),
		}, &przekazanie)

	if przekazanie.Handoff.ArtifactId == nil {
		t.Fatal("przekazanie nie założyło artefaktu")
	}
	odwolanie := tekstZBazy(t, baza,
		`SELECT odwolanie FROM debata_artefakt WHERE identyfikator_zewnetrzny = ?`,
		*przekazanie.Handoff.ArtifactId)
	bajty := bajtyPodOdwolaniem(t, katalog, odwolanie)
	if !strings.Contains(string(bajty), "Wdrożenie rusza") {
		t.Fatal("artefakt przekazania nie niesie treści stanowiska")
	}
	if !strings.Contains(string(bajty), "Konsekwencje") {
		t.Fatal("artefakt przekazania nie niesie zapisu decyzji")
	}
	if !strings.Contains(string(bajty), "Tura 1") {
		t.Fatal("artefakt nie niesie dołączonego transkryptu")
	}
}

// TestOdsluchDebatyDajeNagranieOMierzalnejDlugosci wykazuje skutek syntezy
// mowy: pod odwołaniem leży plik WAV, którego długość zgadza się z liczbą
// próbek zapisanych w jego własnym nagłówku.
//
// Sprawdzian pomija się, gdy na maszynie nie ma silnika mowy: mierzona jest
// droga rdzenia, a nie obecność pakietu, a brak programu ma być widoczny
// w sondzie zależności, nie jako czerwony sprawdzian.
func TestOdsluchDebatyDajeNagranieOMierzalnejDlugosci(t *testing.T) {
	if !czyStoiSilnikMowy() {
		t.Skip("silnik mowy espeak-ng nie stoi na tej maszynie")
	}
	zmontowany, zycie, katalog, _, _ := zlozDebateSprawdzianu(t)

	var wynik shared.RoundtableSpeechSynthesizeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandRoundtableSpeechSynthesize,
		shared.RoundtableSpeechSynthesizeRequest{WindowId: oknoSprawdzianuDebaty}, &wynik)

	if wynik.Uri == nil {
		t.Fatal("odsłuch nie oddał odwołania do nagrania")
	}
	bajty := bajtyPodOdwolaniem(t, katalog, *wynik.Uri)
	if string(bajty[0:4]) != "RIFF" || string(bajty[8:12]) != "WAVE" {
		t.Fatalf("plik odsłuchu nie jest nagraniem WAV; początek: %q", bajty[:12])
	}
	if wynik.DurationMs == nil || *wynik.DurationMs <= 0 {
		t.Fatalf("nagranie melduje długość %v", wynik.DurationMs)
	}
	// Długość liczona z bajtów pliku ma się zgadzać z meldowaną: nagłówek
	// przepisany bez próbek dałby plik poprawny formalnie i pusty w odsłuchu.
	zNaglowka := dlugoscNagraniaMs(len(bajty) - naglowekWav)
	if zNaglowka != *wynik.DurationMs {
		t.Fatalf("z bajtów pliku wychodzi %d ms, a odpowiedź meldowała %d",
			zNaglowka, *wynik.DurationMs)
	}
}

// czyStoiSilnikMowy sprawdza obecność programu bez uruchamiania go.
func czyStoiSilnikMowy() bool {
	for _, pozycja := range ZaleznosciZewnetrzne() {
		if pozycja.Narzedzie.Program == narzedzieSyntezyMowy.Program {
			return pozycja.Stoi
		}
	}
	return false
}

// odczekajChwile oddaje procesor na krótko, czekając na turę biegnącą
// w gorutynie rdzenia. Odstęp jest krótki, a liczba prób skończona: sprawdzian
// ma czekać na skutek, a nie na zegar.
func odczekajChwile() {
	time.Sleep(10 * time.Millisecond)
}
