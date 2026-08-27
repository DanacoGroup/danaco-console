package core

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/transport"
	"danacoconsole/shared"
)

// Sprawdziany SKUTKU blokady fragmentu — fragmentu, którego model nie tknie.
//
// ── Dlaczego ten plik mierzy bazą, a nie odpowiedzią ────────────────────────
// Blokada jest zabezpieczeniem PRZED MODELEM. Zabezpieczenie, o którym wiadomo
// tylko tyle, że odpowiedziało odmową, jest zabezpieczeniem niesprawdzonym:
// rdzeń mógł odmówić i zapisać zmianę mimo to, albo odmówić po zapisie. Dlatego
// każdy sprawdzian tego pliku czyta treść dokumentu i wiersze blokad WPROST
// z bazy, osobnym połączeniem, i pyta o świat, a nie o kopertę.
//
// Szkody, które ten plik ma wykluczyć — wszystkie cztery wystąpiły w produktach
// tej klasy:
//  1. blokada pilnowana wyłącznie w oknie: model woła komendy rdzenia tak samo
//     jak klient, więc omija ją bez wysiłku;
//  2. CICHA bezczynność zamiast odmowy — najgorsza możliwa odpowiedź, bo
//     Operator myśli, że model wykonał polecenie;
//  3. odmowa CAŁOŚCI przy zmianie obejmującej blokadę częściowo — zamiana
//     w całym dokumencie przepada z powodu jednego zablokowanego cytatu;
//  4. blokada skierowana przeciw Operatorowi — właściciel dokumentu nie może
//     zmienić własnego fragmentu.
//
// ── Dlaczego ręka modelu bierze się tu z GNIAZDA ────────────────────────────
// Sprawdzian, który podpisuje żądanie polem `author: model`, mierzy drogę
// uprzejmą: model, który tego pola nie poda, nie zostanie sprawdzony. Dlatego
// te sprawdziany stawiają tożsamość gniazda serwera narzędzi
// (`transport.RodzajNarzedzi`) i pola `author` NIE podają. Mierzą więc to, czego
// model nie może o sobie zataić.

const (
	blokadaOknoSprawdzianu = "okno-blokad"
	// Treść z tą samą frazą DWA razy: raz w miejscu, które zostanie zablokowane,
	// raz poza blokadą. Na tym stoi pomiar bilansu — jedna zamiana wchodzi,
	// druga nie, i odpowiedź musi powiedzieć która.
	blokadaTrescSprawdzianu = "Umowa numer 17/2026 zawarta w Warszawie.\n" +
		"Podstawa prawna: Umowa numer 17/2026.\n" +
		"Uwagi redakcyjne do Umowa numer 17/2026 bez znaczenia prawnego."
)

// blokadaUprzazSprawdzianu trzyma rdzeń, kontekst życia i połączenie pomiarowe.
type blokadaUprzazSprawdzianu struct {
	zmontowany *Zmontowany
	zycie      context.Context
	baza       *sql.DB
	dokument   string
}

// blokadaZmontuj składa rdzeń nad świeżą bazą, zakłada dokument i wpisuje mu
// treść — tą samą drogą, którą robi to okno pracy z dokumentem.
func blokadaZmontuj(t *testing.T) *blokadaUprzazSprawdzianu {
	t.Helper()

	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza, err := sql.Open("sqlite", filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do pomiaru niezależnego: %v", err)
	}
	t.Cleanup(func() { _ = baza.Close() })

	uprzaz := &blokadaUprzazSprawdzianu{zmontowany: zmontowany, zycie: zycie, baza: baza}

	var otwarcie shared.StudioDocumentOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentOpen,
		shared.StudioDocumentOpenRequest{WindowId: blokadaOknoSprawdzianu}, &otwarcie)
	uprzaz.dokument = otwarcie.Document.Id

	var zapis shared.StudioDocumentSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId:    uprzaz.dokument,
			Content:       blokadaTrescSprawdzianu,
			CreateVersion: blokadaWskaznikPrawdy(),
		}, &zapis)

	return uprzaz
}

// blokadaWskaznikPrawdy oddaje wskaźnik na prawdę — pola kontraktu są
// nieobowiązkowe, więc trzeba adresu, nie wartości.
func blokadaWskaznikPrawdy() *bool {
	wartosc := true
	return &wartosc
}

// blokadaZalozBlokade zakłada blokadę na wskazanym fragmencie ręką Operatora
// i oddaje jej identyfikator.
func (u *blokadaUprzazSprawdzianu) blokadaZalozBlokade(t *testing.T, od, do int,
	nazwa, powod string) string {
	t.Helper()

	var odpowiedz shared.StudioLockAddResponse
	wykonajUdana(t, u.zmontowany, u.zycie, shared.CommandStudioLockAdd,
		shared.StudioLockAddRequest{
			DocumentId: u.dokument, RangeStart: od, RangeEnd: do,
			Name: nazwa, Reason: &powod,
		}, &odpowiedz)
	if odpowiedz.Lock.Id == "" {
		t.Fatal("blokada założona bez identyfikatora — nie ma czym jej potem wskazać")
	}
	return odpowiedz.Lock.Id
}

// blokadaWykonajJakoModel wywołuje komendę z tożsamością gniazda serwera
// narzędzi modelu. Pola `author` NIE dokłada z zamysłem: mierzymy drogę, której
// model nie może o sobie zataić.
func (u *blokadaUprzazSprawdzianu) blokadaWykonajJakoModel(t *testing.T,
	komenda shared.MessageType, ladunek any) protocol.Koperta {
	t.Helper()

	koperta, err := protocol.NowaKoperta(komenda, "sprawdzian-modelu", "", ladunek)
	if err != nil {
		t.Fatalf("nie można złożyć koperty %s: %v", komenda, err)
	}
	ctx, przerwij := context.WithTimeout(u.zycie, granicaSprawdzianuKomendy(komenda))
	defer przerwij()
	// Tak wygląda gniazdo serwera narzędzi modelu i nic innego tak nie wygląda
	// (`core/sprawca.go`): rodzaj `narzedzia` przedstawia się przy nawiązaniu.
	ctx = zPolaczeniem(ctx, transport.Tozsamosc{
		IdPolaczenia: "gniazdo-narzedzi-sprawdzianu",
		IdKlienta:    "serwer-narzedzi",
		Rodzaj:       transport.RodzajNarzedzi,
		IdOkna:       blokadaOknoSprawdzianu,
	})
	return u.zmontowany.Rdzen.Wykonaj(ctx, koperta)
}

// blokadaTrescZBazy czyta treść dokumentu WPROST z bazy — pomiar niezależny od
// tego, co powiedziała odpowiedź.
func (u *blokadaUprzazSprawdzianu) blokadaTrescZBazy(t *testing.T) string {
	t.Helper()

	var tresc sql.NullString
	err := u.baza.QueryRow(
		`SELECT tresc FROM dokument_studio WHERE identyfikator_zewnetrzny = ?`,
		u.dokument).Scan(&tresc)
	if err != nil {
		t.Fatalf("nie można odczytać treści dokumentu z bazy: %v", err)
	}
	return tresc.String
}

// blokadaIleBlokadWBazie liczy wiersze blokad dokumentu.
func (u *blokadaUprzazSprawdzianu) blokadaIleBlokadWBazie(t *testing.T) int {
	t.Helper()

	var liczba int
	err := u.baza.QueryRow(`SELECT COUNT(*) FROM blokada_fragmentu_studio b
	                        JOIN dokument_studio d ON d.id = b.dokument_id
	                        WHERE d.identyfikator_zewnetrzny = ?`, u.dokument).Scan(&liczba)
	if err != nil {
		t.Fatalf("nie można policzyć blokad w bazie: %v", err)
	}
	return liczba
}

// ── Sprawdzian pierwszy: odmowa NAZWANA ─────────────────────────────────────

// TestBlokadaOdmawiaModelowiINazywaFragment mierzy sedno wymagania: czynność
// modelu godząca w zablokowany fragment wraca błędem, który mówi, KTÓRY
// fragment i JAKA blokada ją zatrzymała — a treść dokumentu zostaje nietknięta.
func TestBlokadaOdmawiaModelowiINazywaFragment(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	// Blokada obejmuje drugi wiersz — podstawę prawną, która ma zostać dosłownie.
	poczatek := strings.Index(blokadaTrescSprawdzianu, "Podstawa prawna")
	koniec := poczatek + len("Podstawa prawna: Umowa numer 17/2026.")
	uprzaz.blokadaZalozBlokade(t, poczatek, koniec,
		"podstawa prawna", "podstawa prawna — nie zmieniać")

	trescPrzed := uprzaz.blokadaTrescZBazy(t)

	odpowiedz := uprzaz.blokadaWykonajJakoModel(t, shared.CommandStudioTextEdit,
		shared.StudioTextEditRequest{
			DocumentId: uprzaz.dokument,
			RangeStart: poczatek, RangeEnd: koniec,
			Text: "Podstawa prawna: brak.",
		})

	if odpowiedz.Error == nil {
		t.Fatalf("model zmienił zablokowany fragment bez odmowy; ładunek: %s",
			poczatekLadunku(odpowiedz.Payload))
	}
	// Odmowa NAZWANA, nie „nie wolno": treść musi nieść nazwę blokady i zakres.
	tresc := odpowiedz.Error.Message
	if !strings.Contains(tresc, "podstawa prawna") {
		t.Errorf("odmowa nie nazywa blokady, która zatrzymała czynność: %q", tresc)
	}
	if !strings.Contains(tresc, "nie zmieniać") {
		t.Errorf("odmowa nie niesie powodu blokady: %q", tresc)
	}
	if !strings.Contains(tresc, "od znaku") {
		t.Errorf("odmowa nie nazywa fragmentu, którego dotyczy: %q", tresc)
	}

	// Pomiar niezależny: treść w bazie musi być ta sama.
	if po := uprzaz.blokadaTrescZBazy(t); po != trescPrzed {
		t.Errorf("odmowa wróciła, a treść w bazie się zmieniła:\nprzed: %q\npo:    %q",
			trescPrzed, po)
	}
}

// ── Sprawdzian drugi: bilans, nie odmowa całości ────────────────────────────

// TestBlokadaZamianaWCalymDokumencieOddajeBilans mierzy trzecią odpowiedź
// zderzenia: zmiana obejmująca blokadę CZĘŚCIOWO wykonuje się POZA blokadą
// i oddaje bilans. Odmowa całości byłaby nieproporcjonalna, przemilczenie
// pominięcia — zakazane.
func TestBlokadaZamianaWCalymDokumencieOddajeBilans(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	poczatek := strings.Index(blokadaTrescSprawdzianu, "Podstawa prawna")
	koniec := poczatek + len("Podstawa prawna: Umowa numer 17/2026.")
	uprzaz.blokadaZalozBlokade(t, poczatek, koniec,
		"podstawa prawna", "podstawa prawna — nie zmieniać")

	// Zamiana w CAŁYM dokumencie: bez wskazania zakresu. Trafia w trzy miejsca,
	// z których jedno stoi pod blokadą.
	nowaTresc := strings.ReplaceAll(blokadaTrescSprawdzianu,
		"Umowa numer 17/2026", "Umowa numer 18/2026")
	odpowiedz := uprzaz.blokadaWykonajJakoModel(t, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: uprzaz.dokument,
			Content:    nowaTresc,
		})

	if odpowiedz.Error != nil {
		t.Fatalf("zamiana obejmująca blokadę CZĘŚCIOWO odmówiła całości — "+
			"odmowa nieproporcjonalna: kod=%s treść=%s",
			odpowiedz.Error.Code, odpowiedz.Error.Message)
	}

	// Bilans musi stać w odpowiedzi i nazywać blokadę.
	bilans := blokadaBilansZOdpowiedzi(t, odpowiedz)
	if bilans.SkippedCount == 0 {
		t.Fatal("bilans nie zgłasza ani jednego pominięcia, choć zamiana trafiła w blokadę")
	}
	nazwana := false
	for _, pominiecie := range bilans.Skipped {
		if pominiecie.LockName != nil && *pominiecie.LockName == "podstawa prawna" {
			nazwana = true
		}
	}
	if !nazwana {
		t.Errorf("bilans nie mówi, PRZEZ KTÓRĄ blokadę pominięto fragment: %+v", bilans.Skipped)
	}
	if bilans.Applied == 0 {
		t.Error("bilans nie zgłasza ani jednej zmiany wniesionej — " +
			"zamiana poza blokadą miała wejść")
	}

	// Pomiar niezależny: fragment zablokowany dosłownie ten sam, reszta zmieniona.
	po := uprzaz.blokadaTrescZBazy(t)
	if !strings.Contains(po, "Podstawa prawna: Umowa numer 17/2026.") {
		t.Errorf("zablokowany fragment został zmieniony wbrew blokadzie; treść: %q", po)
	}
	if !strings.Contains(po, "Umowa numer 18/2026 zawarta w Warszawie") {
		t.Errorf("zmiana POZA blokadą nie weszła — odmowa całości podana jako bilans; treść: %q", po)
	}
}

// blokadaBilansZOdpowiedzi wyjmuje bilans z ładunku odpowiedzi. Odpowiedź bez
// bilansu jest tu niepowodzeniem, nie brakiem: pominięcia nie wolno przemilczeć.
func blokadaBilansZOdpowiedzi(t *testing.T, odpowiedz protocol.Koperta) shared.StudioActionBalance {
	t.Helper()

	var ladunek struct {
		Balance *shared.StudioActionBalance `json:"balance"`
	}
	if err := protocol.LadunekDo(odpowiedz, &ladunek); err != nil {
		t.Fatalf("nieczytelny ładunek odpowiedzi: %v", err)
	}
	if ladunek.Balance == nil {
		t.Fatalf("odpowiedź nie niesie bilansu, choć czynność pominęła fragment zablokowany; "+
			"ładunek: %s", poczatekLadunku(odpowiedz.Payload))
	}
	return *ladunek.Balance
}

// ── Sprawdzian trzeci: Operator bez przeszkód ───────────────────────────────

// TestBlokadaNieZatrzymujeOperatora mierzy, że blokada jest skierowana przeciw
// modelowi, a nie przeciw właścicielowi dokumentu. Blokada działająca także na
// Operatora jest osobnym, jawnym ustawieniem — nie zachowaniem domyślnym.
func TestBlokadaNieZatrzymujeOperatora(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	poczatek := strings.Index(blokadaTrescSprawdzianu, "Podstawa prawna")
	koniec := poczatek + len("Podstawa prawna: Umowa numer 17/2026.")
	uprzaz.blokadaZalozBlokade(t, poczatek, koniec, "podstawa prawna", "nie zmieniać")

	// Ta sama czynność, ten sam zakres, ręka Operatora — czyli bez tożsamości
	// gniazda serwera narzędzi i bez podpisu `author`.
	nowaTresc := strings.Replace(blokadaTrescSprawdzianu,
		"Podstawa prawna: Umowa numer 17/2026.",
		"Podstawa prawna: Umowa numer 19/2026.", 1)
	var zapis shared.StudioDocumentSaveResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: uprzaz.dokument, Content: nowaTresc,
		}, &zapis)

	po := uprzaz.blokadaTrescZBazy(t)
	if !strings.Contains(po, "Umowa numer 19/2026") {
		t.Errorf("blokada zatrzymała OPERATORA — właściciel dokumentu nie zmienił "+
			"własnego fragmentu; treść: %q", po)
	}
}

// TestBlokadaZasieguEveryoneWiazeTakzeOperatora jest drugą stroną tej samej
// zasady: zasięg `everyone` jest tym OSOBNYM, JAWNYM ustawieniem i wtedy
// blokada wiąże także Operatora.
func TestBlokadaZasieguEveryoneWiazeTakzeOperatora(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	poczatek := strings.Index(blokadaTrescSprawdzianu, "Podstawa prawna")
	koniec := poczatek + len("Podstawa prawna: Umowa numer 17/2026.")
	zasieg := shared.StudioLockScope(shared.StudioLockScopeEveryone)
	powod := "cytat urzędowy"
	var zalozenie shared.StudioLockAddResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioLockAdd,
		shared.StudioLockAddRequest{
			DocumentId: uprzaz.dokument, RangeStart: poczatek, RangeEnd: koniec,
			Name: "cytat", Reason: &powod, Scope: &zasieg,
		}, &zalozenie)

	trescPrzed := uprzaz.blokadaTrescZBazy(t)
	odmowa := wykonajOdmowna(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioTextEdit,
		shared.StudioTextEditRequest{
			DocumentId: uprzaz.dokument, RangeStart: poczatek, RangeEnd: koniec,
			Text: "Podstawa prawna: brak.",
		})
	if !strings.Contains(odmowa.Message, "cytat") {
		t.Errorf("odmowa nie nazywa blokady zasięgu everyone: %q", odmowa.Message)
	}
	if po := uprzaz.blokadaTrescZBazy(t); po != trescPrzed {
		t.Errorf("blokada zasięgu everyone nie zatrzymała Operatora; treść: %q", po)
	}
}

// ── Sprawdzian czwarty: blokada przechodzi przez wersje ─────────────────────

// TestBlokadaPrzechodziPrzezPrzywrocenieWersji mierzy, że przywrócenie
// wcześniejszej wersji NIE GUBI blokad. Blokada wisi przy dokumencie, nie przy
// wersji — ale sprawdzić trzeba świat, nie zamysł.
func TestBlokadaPrzechodziPrzezPrzywrocenieWersji(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	poczatek := strings.Index(blokadaTrescSprawdzianu, "Podstawa prawna")
	koniec := poczatek + len("Podstawa prawna: Umowa numer 17/2026.")
	kodBlokady := uprzaz.blokadaZalozBlokade(t, poczatek, koniec,
		"podstawa prawna", "nie zmieniać")

	// Wersja późniejsza, żeby było do czego wracać.
	var drugiZapis shared.StudioDocumentSaveResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId:    uprzaz.dokument,
			Content:       blokadaTrescSprawdzianu + "\nAneks pierwszy.",
			CreateVersion: blokadaWskaznikPrawdy(),
		}, &drugiZapis)

	var historia shared.StudioVersionSeriesListResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioVersionSeriesList,
		shared.StudioVersionSeriesListRequest{DocumentId: uprzaz.dokument}, &historia)
	if historia.InitialVersionId == nil || *historia.InitialVersionId == "" {
		t.Fatal("wykaz szeregów nie wskazuje wersji założycielskiej — nie ma do czego wrócić")
	}

	var przywrocenie shared.StudioRepositoryRestoreResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioRepositoryRestore,
		shared.StudioRepositoryRestoreRequest{
			DocumentId: uprzaz.dokument, VersionId: *historia.InitialVersionId,
		}, &przywrocenie)

	// Pomiar niezależny: wiersz blokady stoi dalej.
	if ile := uprzaz.blokadaIleBlokadWBazie(t); ile != 1 {
		t.Fatalf("przywrócenie wersji zgubiło blokady: w bazie stoi %d, miała stać 1", ile)
	}
	var wykaz shared.StudioLockListResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioLockList,
		shared.StudioLockListRequest{DocumentId: uprzaz.dokument}, &wykaz)
	znaleziona := false
	for _, blokada := range wykaz.Locks {
		if blokada.Id == kodBlokady {
			znaleziona = true
		}
	}
	if !znaleziona {
		t.Errorf("blokada %s nie wychodzi wykazem po przywróceniu wersji: %+v",
			kodBlokady, wykaz.Locks)
	}

	// I dalej DZIAŁA — blokada, która przetrwała jako wiersz, ale przestała
	// pilnować, jest gorsza niż zgubiona, bo wygląda na czynną.
	odpowiedz := uprzaz.blokadaWykonajJakoModel(t, shared.CommandStudioTextEdit,
		shared.StudioTextEditRequest{
			DocumentId: uprzaz.dokument, RangeStart: poczatek, RangeEnd: koniec,
			Text: "Podstawa prawna: brak.",
		})
	if odpowiedz.Error == nil {
		t.Error("blokada przetrwała przywrócenie wersji jako wiersz, ale przestała pilnować")
	}
}

// ── Sprawdzian piąty: blokadę zdejmuje wyłącznie Operator ───────────────────

// TestBlokadaZdejmujeWylacznieOperator mierzy, że model blokady nie zdejmie —
// ani wprost, ani obejściem przez podpisanie się Operatorem.
func TestBlokadaZdejmujeWylacznieOperator(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	poczatek := strings.Index(blokadaTrescSprawdzianu, "Podstawa prawna")
	koniec := poczatek + len("Podstawa prawna: Umowa numer 17/2026.")
	kodBlokady := uprzaz.blokadaZalozBlokade(t, poczatek, koniec, "podstawa prawna", "nie zmieniać")

	// Obejście: gniazdo narzędzi modelu, a w żądaniu podpis „Operator".
	// Podpis jest twierdzeniem modelu o sobie, nie faktem — i nie ma go przepuścić.
	operator := shared.StudioAuthor(shared.StudioAuthorUzytkownik)
	odpowiedz := uprzaz.blokadaWykonajJakoModel(t, shared.CommandStudioLockRemove,
		shared.StudioLockRemoveRequest{
			DocumentId: uprzaz.dokument, LockId: kodBlokady, Author: &operator,
		})
	if odpowiedz.Error == nil {
		t.Fatalf("model zdjął blokadę, podpisując się Operatorem; ładunek: %s",
			poczatekLadunku(odpowiedz.Payload))
	}
	if !strings.Contains(odpowiedz.Error.Message, "WYŁĄCZNIE Operator") {
		t.Errorf("odmowa nie nazywa zasady, która ją wywołała: %q", odpowiedz.Error.Message)
	}
	if ile := uprzaz.blokadaIleBlokadWBazie(t); ile != 1 {
		t.Errorf("blokada zniknęła z bazy mimo odmowy: stoi %d wierszy", ile)
	}

	// Operator zdejmuje ją bez przeszkód.
	var zdjecie shared.StudioLockRemoveResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioLockRemove,
		shared.StudioLockRemoveRequest{DocumentId: uprzaz.dokument, LockId: kodBlokady}, &zdjecie)
	if !zdjecie.Removed {
		t.Error("Operator nie zdjął własnej blokady")
	}
	if ile := uprzaz.blokadaIleBlokadWBazie(t); ile != 0 {
		t.Errorf("blokada zdjęta odpowiedzią została w bazie: stoi %d wierszy", ile)
	}
}

// TestBlokadaZakladaWylacznieOperator mierzy drugą stronę: model nie
// unieruchamia fragmentów dokumentu, bo blokada jest narzędziem Operatora
// PRZECIW wykonawcom, nie odwrotnie.
func TestBlokadaZakladaWylacznieOperator(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	odpowiedz := uprzaz.blokadaWykonajJakoModel(t, shared.CommandStudioLockAdd,
		shared.StudioLockAddRequest{
			DocumentId: uprzaz.dokument, RangeStart: 0, RangeEnd: 10, Name: "moja blokada",
		})
	if odpowiedz.Error == nil {
		t.Fatalf("model założył blokadę fragmentu; ładunek: %s",
			poczatekLadunku(odpowiedz.Payload))
	}
	if ile := uprzaz.blokadaIleBlokadWBazie(t); ile != 0 {
		t.Errorf("blokada modelu weszła do bazy mimo odmowy: stoi %d wierszy", ile)
	}
}
