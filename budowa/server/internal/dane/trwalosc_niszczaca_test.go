// Sprawdzian obejmuje trzy drogi kasujące dane Operatora bez jego wskazania:
// kaskadę schematu, kosz sesji i egzekucję retencji historii.
package dane

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"danacoconsole/server/internal/store"
	"danacoconsole/shared"
)

// drzewoSprawdzianu zakłada bazę wraz z jednym łańcuchem karta → sesja → okno
// i zwraca repozytoria oraz identyfikatory.
type drzewoSprawdzianu struct {
	zestaw   *Zestaw
	baza     *store.Baza
	kartaID  int64
	sesjaID  int64
	oknoID   int64
	oknoKod  string
	sesjaKod string
}

// zalozDrzewo buduje łańcuch na świeżej bazie. Środowiska, moduły i kanały
// modelu przychodzą z migracji, więc sprawdzian ich nie zakłada — bierze
// pierwszy wiersz każdego.
func zalozDrzewo(t *testing.T, przyrostek string) *drzewoSprawdzianu {
	t.Helper()
	ctx := context.Background()

	baza, err := store.Otworz(filepath.Join(t.TempDir(), "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })

	zestaw, err := Otworz(ctx, baza)
	if err != nil {
		t.Fatalf("nie można otworzyć repozytoriów: %v", err)
	}
	t.Cleanup(func() { _ = zestaw.Zamknij() })

	srodowisko, err := zestaw.Srodowiska.Pierwsze(ctx)
	if err != nil {
		t.Fatalf("nie można odczytać środowiska: %v", err)
	}
	// Zero w miejscu konta znaczy kartę zastaną, sprzed rozdzielenia kont
	// (migracja 407): przebieg nie stawia kontekstu z kontem Operatora.
	kartaID, err := zestaw.KartySesji.Zapewnij(ctx, srodowisko.ID, 0, "karta "+przyrostek)
	if err != nil {
		t.Fatalf("nie można założyć karty sesji: %v", err)
	}

	sesjaKod := "sesja-" + przyrostek
	sesjaID, err := zestaw.Sesje.Utworz(ctx, Sesja{
		KartaSesjiID:            kartaID,
		Tytul:                   "Sesja sprawdzianu " + przyrostek,
		Stan:                    shared.SessionStatusActive,
		IdentyfikatorZewnetrzny: &sesjaKod,
	})
	if err != nil {
		t.Fatalf("nie można założyć sesji: %v", err)
	}

	var modulID, kanalID int64
	if err := baza.DB.QueryRow("SELECT id FROM modul ORDER BY id LIMIT 1").Scan(&modulID); err != nil {
		t.Fatalf("nie można odczytać modułu: %v", err)
	}
	if err := baza.DB.QueryRow("SELECT id FROM kanal_modelu ORDER BY id LIMIT 1").Scan(&kanalID); err != nil {
		t.Fatalf("nie można odczytać kanału modelu: %v", err)
	}

	oknoKod := "okno-" + przyrostek
	oknoID, err := zestaw.Okna.Utworz(ctx, Okno{
		SesjaID:                 sesjaID,
		ModulID:                 modulID,
		KanalModeluID:           kanalID,
		SrodowiskoWykonania:     shared.ExecutionEnvLocal,
		TrybUprawnien:           shared.PermissionModeManual,
		RolaOkna:                shared.WindowRoleStandalone,
		Stan:                    shared.WindowStatusOpen,
		IdentyfikatorZewnetrzny: &oknoKod,
	})
	if err != nil {
		t.Fatalf("nie można założyć okna: %v", err)
	}

	return &drzewoSprawdzianu{
		zestaw: zestaw, baza: baza,
		kartaID: kartaID, sesjaID: sesjaID, oknoID: oknoID,
		oknoKod: oknoKod, sesjaKod: sesjaKod,
	}
}

// dopiszWiadomosci dokłada oknu wskazaną liczbę wypowiedzi testowych,
// zwracając identyfikatory zapisanych wierszy w kolejności dopisania.
func (d *drzewoSprawdzianu) dopiszWiadomosci(t *testing.T, ile int) []int64 {
	t.Helper()
	ctx := context.Background()

	identyfikatory := make([]int64, 0, ile)
	for i := 0; i < ile; i++ {
		tresc := fmt.Sprintf("wypowiedź %d", i+1)
		kod := fmt.Sprintf("%s-wiadomosc-%d", d.oknoKod, i+1)
		id, err := d.zestaw.Wiadomosci.Dopisz(ctx, Wiadomosc{
			OknoID:                  d.oknoID,
			Rola:                    shared.MessageRoleUser,
			RodzajTresci:            shared.ChunkKindText,
			Stan:                    shared.MessageStatusComplete,
			Tresc:                   &tresc,
			IdentyfikatorZewnetrzny: &kod,
		})
		if err != nil {
			t.Fatalf("nie można dopisać wiadomości: %v", err)
		}
		identyfikatory = append(identyfikatory, id)
	}
	return identyfikatory
}

// policz zwraca liczbę wierszy tabeli spełniających warunek podanego
// zapytania SQL, przerywając sprawdzian przy błędzie odczytu.
func (d *drzewoSprawdzianu) policz(t *testing.T, zapytanie string, argumenty ...any) int {
	t.Helper()
	var liczba int
	if err := d.baza.DB.QueryRow(zapytanie, argumenty...).Scan(&liczba); err != nil {
		t.Fatalf("nie można policzyć wierszy (%s): %v", zapytanie, err)
	}
	return liczba
}

// TestKaskadaZabieraOknaIWiadomosci sprawdza więz schematu, na którym stoi każde
// kasowanie sesji. Kaskada niedziałająca zostawia w bazie okna bez sesji
// i wypowiedzi bez okna — nikt ich nie zobaczy i nikt ich nie sprzątnie.
func TestKaskadaZabieraOknaIWiadomosci(t *testing.T) {
	drzewo := zalozDrzewo(t, "kaskada")
	drzewo.dopiszWiadomosci(t, 3)
	ctx := context.Background()

	if liczba := drzewo.policz(t, "SELECT count(*) FROM wiadomosc WHERE okno_komunikacji_id = ?",
		drzewo.oknoID); liczba != 3 {
		t.Fatalf("przed kasowaniem okno ma %d wypowiedzi, oczekiwane 3", liczba)
	}

	if err := drzewo.zestaw.Sesje.Usun(ctx, drzewo.sesjaID); err != nil {
		t.Fatalf("nie można skasować sesji: %v", err)
	}

	if liczba := drzewo.policz(t, "SELECT count(*) FROM okno_komunikacji WHERE sesja_id = ?",
		drzewo.sesjaID); liczba != 0 {
		t.Errorf("po skasowaniu sesji zostało %d okien-sierot", liczba)
	}
	if liczba := drzewo.policz(t, "SELECT count(*) FROM wiadomosc WHERE okno_komunikacji_id = ?",
		drzewo.oknoID); liczba != 0 {
		t.Errorf("po skasowaniu sesji zostało %d wypowiedzi-sierot", liczba)
	}
	if err := drzewo.baza.SprawdzSpojnosc(); err != nil {
		t.Errorf("baza po kasowaniu sesji jest niespójna: %v", err)
	}
}

// TestKaskadaKartySesjiZabieraCaleDrzewo sprawdza szczebel wyżej: karta sesji
// jest korzeniem łańcucha, więc jej skasowanie ma zabrać wszystko pod nią.
func TestKaskadaKartySesjiZabieraCaleDrzewo(t *testing.T) {
	drzewo := zalozDrzewo(t, "karta")
	drzewo.dopiszWiadomosci(t, 2)

	if _, err := drzewo.baza.DB.Exec("DELETE FROM karta_sesji WHERE id = ?", drzewo.kartaID); err != nil {
		t.Fatalf("nie można skasować karty sesji: %v", err)
	}

	for _, przypadek := range []struct {
		nazwa     string
		zapytanie string
		argument  any
	}{
		{"sesje", "SELECT count(*) FROM sesja WHERE karta_sesji_id = ?", drzewo.kartaID},
		{"okna", "SELECT count(*) FROM okno_komunikacji WHERE sesja_id = ?", drzewo.sesjaID},
		{"wypowiedzi", "SELECT count(*) FROM wiadomosc WHERE okno_komunikacji_id = ?", drzewo.oknoID},
	} {
		if liczba := drzewo.policz(t, przypadek.zapytanie, przypadek.argument); liczba != 0 {
			t.Errorf("po skasowaniu karty sesji zostało %d wierszy: %s", liczba, przypadek.nazwa)
		}
	}
	if err := drzewo.baza.SprawdzSpojnosc(); err != nil {
		t.Errorf("baza po kasowaniu karty sesji jest niespójna: %v", err)
	}
}

// TestKoszUkrywaSesjeIOddajeJaZPowrotem sprawdza obieg pomyłki: sesja wrzucona
// do kosza znika z wykazu żywych, ale nie z bazy, i wraca w całości.
func TestKoszUkrywaSesjeIOddajeJaZPowrotem(t *testing.T) {
	drzewo := zalozDrzewo(t, "kosz")
	drzewo.dopiszWiadomosci(t, 2)
	ctx := context.Background()

	if err := drzewo.zestaw.KoszSesji.PrzeniesDoKosza(ctx, drzewo.sesjaID); err != nil {
		t.Fatalf("nie można przenieść sesji do kosza: %v", err)
	}

	zywe, err := drzewo.zestaw.Sesje.Lista(ctx, drzewo.kartaID)
	if err != nil {
		t.Fatalf("nie można odczytać wykazu sesji: %v", err)
	}
	for _, sesja := range zywe {
		if sesja.ID == drzewo.sesjaID {
			t.Error("sesja z kosza stoi w wykazie sesji żywych")
		}
	}

	wKoszu, err := drzewo.zestaw.KoszSesji.Lista(ctx)
	if err != nil {
		t.Fatalf("nie można odczytać kosza: %v", err)
	}
	if len(wKoszu) != 1 || wKoszu[0].ID != drzewo.sesjaID {
		t.Fatalf("kosz niesie %+v", wKoszu)
	}

	// Wypowiedzi zostają — kosz nie jest kasowaniem.
	if liczba := drzewo.policz(t, "SELECT count(*) FROM wiadomosc WHERE okno_komunikacji_id = ?",
		drzewo.oknoID); liczba != 2 {
		t.Errorf("kosz zabrał wypowiedzi: zostało %d z 2", liczba)
	}

	przywrocona, err := drzewo.zestaw.KoszSesji.Przywroc(ctx, drzewo.sesjaID)
	if err != nil {
		t.Fatalf("nie można przywrócić sesji: %v", err)
	}
	if !przywrocona {
		t.Error("przywrócenie sesji z kosza nie zgłosiło skutku")
	}

	zywePo, err := drzewo.zestaw.Sesje.Lista(ctx, drzewo.kartaID)
	if err != nil {
		t.Fatalf("nie można odczytać wykazu sesji: %v", err)
	}
	if len(zywePo) != 1 || zywePo[0].ID != drzewo.sesjaID {
		t.Errorf("sesja nie wróciła do wykazu żywych: %+v", zywePo)
	}
}

// TestPowtorneWrzuceniePrzywrocenieNieUdajeSkutku pilnuje rozróżnienia opisanego
// w repozytorium: powtórka nie jest trafieniem. Wywołujący nie może wziąć drugiej
// próby za skutek.
func TestPowtorneWrzuceniePrzywrocenieNieUdajeSkutku(t *testing.T) {
	drzewo := zalozDrzewo(t, "powtorka")
	ctx := context.Background()

	if err := drzewo.zestaw.KoszSesji.PrzeniesDoKosza(ctx, drzewo.sesjaID); err != nil {
		t.Fatalf("pierwsze wrzucenie nie powiodło się: %v", err)
	}
	if err := drzewo.zestaw.KoszSesji.PrzeniesDoKosza(ctx, drzewo.sesjaID); err == nil {
		t.Error("powtórne wrzucenie do kosza zgłosiło powodzenie")
	}

	if _, err := drzewo.zestaw.KoszSesji.Przywroc(ctx, drzewo.sesjaID); err != nil {
		t.Fatalf("przywrócenie nie powiodło się: %v", err)
	}
	przywrocona, err := drzewo.zestaw.KoszSesji.Przywroc(ctx, drzewo.sesjaID)
	if err != nil {
		t.Errorf("powtórne przywrócenie zgłosiło błąd: %v", err)
	}
	if przywrocona {
		t.Error("powtórne przywrócenie zgłosiło skutek, którego nie było")
	}
}

// TestCzyszczenieKoszaTykaWylacznieWiersziPoTerminie sprawdza, że czyszczenie
// kosza usuwa wyłącznie wiersze po terminie, a sesja świeżo wrzucona przeżywa.
func TestCzyszczenieKoszaTykaWylacznieWiersziPoTerminie(t *testing.T) {
	drzewo := zalozDrzewo(t, "termin")
	drzewo.dopiszWiadomosci(t, 2)
	ctx := context.Background()

	// Druga sesja w tej samej karcie — wrzucona do kosza świeżo.
	swiezyKod := "sesja-swieza"
	swiezaID, err := drzewo.zestaw.Sesje.Utworz(ctx, Sesja{
		KartaSesjiID:            drzewo.kartaID,
		Tytul:                   "Sesja świeżo wrzucona",
		Stan:                    shared.SessionStatusActive,
		IdentyfikatorZewnetrzny: &swiezyKod,
	})
	if err != nil {
		t.Fatalf("nie można założyć drugiej sesji: %v", err)
	}

	for _, id := range []int64{drzewo.sesjaID, swiezaID} {
		if err := drzewo.zestaw.KoszSesji.PrzeniesDoKosza(ctx, id); err != nil {
			t.Fatalf("nie można przenieść sesji %d do kosza: %v", id, err)
		}
	}

	// Pierwszej sesji cofa się znacznik: wiersz leży w koszu dłużej niż termin.
	dawno := time.Now().UTC().AddDate(0, 0, -30).Format("2006-01-02T15:04:05.000Z")
	if _, err := drzewo.baza.DB.Exec("UPDATE sesja SET usunieto_o = ? WHERE id = ?",
		dawno, drzewo.sesjaID); err != nil {
		t.Fatalf("nie można cofnąć znacznika kosza: %v", err)
	}

	granica := time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02T15:04:05.000Z")
	usuniete, err := drzewo.zestaw.KoszSesji.UsunPrzeterminowane(ctx, granica)
	if err != nil {
		t.Fatalf("czyszczenie kosza nie powiodło się: %v", err)
	}
	if usuniete != 1 {
		t.Errorf("czyszczenie skasowało %d sesji, oczekiwana jedna", usuniete)
	}

	if liczba := drzewo.policz(t, "SELECT count(*) FROM sesja WHERE id = ?",
		drzewo.sesjaID); liczba != 0 {
		t.Error("sesja po terminie przeżyła czyszczenie")
	}
	if liczba := drzewo.policz(t, "SELECT count(*) FROM sesja WHERE id = ?",
		swiezaID); liczba != 1 {
		t.Error("sesja świeżo wrzucona do kosza została skasowana przed terminem")
	}
	if liczba := drzewo.policz(t, "SELECT count(*) FROM wiadomosc WHERE okno_komunikacji_id = ?",
		drzewo.oknoID); liczba != 0 {
		t.Errorf("po trwałym skasowaniu sesji zostało %d wypowiedzi-sierot", liczba)
	}
	if err := drzewo.baza.SprawdzSpojnosc(); err != nil {
		t.Errorf("baza po czyszczeniu kosza jest niespójna: %v", err)
	}
}

// TestCzyszczenieKoszaNieTykaSesjiZywych pilnuje najgorszego możliwego wyniku
// tej drogi: sesja, której nikt nie wyrzucił, ma przeżyć każde czyszczenie.
func TestCzyszczenieKoszaNieTykaSesjiZywych(t *testing.T) {
	drzewo := zalozDrzewo(t, "zywa")
	ctx := context.Background()

	// Granica w przyszłości — czyszczenie najbardziej zachłanne, jakie da się
	// wywołać.
	granica := time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02T15:04:05.000Z")
	usuniete, err := drzewo.zestaw.KoszSesji.UsunPrzeterminowane(ctx, granica)
	if err != nil {
		t.Fatalf("czyszczenie kosza nie powiodło się: %v", err)
	}
	if usuniete != 0 {
		t.Errorf("czyszczenie skasowało %d sesji żywych", usuniete)
	}
	if liczba := drzewo.policz(t, "SELECT count(*) FROM sesja WHERE id = ?",
		drzewo.sesjaID); liczba != 1 {
		t.Error("sesja żywa zniknęła po czyszczeniu kosza")
	}
}

// TestRetencjaLiczbaPozycjiZostawiaNajnowsze sprawdza przycinanie liczbą.
// Zasada kasuje sama, bez wskazania Operatora, więc pomyłka o jeden zabiera
// wypowiedź, której nikt nie kazał usuwać.
func TestRetencjaLiczbaPozycjiZostawiaNajnowsze(t *testing.T) {
	drzewo := zalozDrzewo(t, "retencja-liczba")
	drzewo.dopiszWiadomosci(t, 5)
	ctx := context.Background()

	trzymane := 2
	if _, err := drzewo.zestaw.Historia.ZapiszZasade(ctx, ZasadaPrzechowywania{
		Zakres: "window", ZakresKod: drzewo.oknoKod, PozycjeTrzymane: &trzymane,
	}); err != nil {
		t.Fatalf("nie można zapisać zasady: %v", err)
	}

	usuniete, err := drzewo.zestaw.Historia.Egzekwuj(ctx, drzewo.oknoKod)
	if err != nil {
		t.Fatalf("egzekucja zasady nie powiodła się: %v", err)
	}
	if usuniete != 3 {
		t.Errorf("zasada przycięła %d pozycji, oczekiwane 3", usuniete)
	}

	zostalo, err := drzewo.zestaw.Historia.Policz(ctx, drzewo.oknoKod)
	if err != nil {
		t.Fatalf("nie można policzyć historii: %v", err)
	}
	if zostalo != trzymane {
		t.Errorf("po przycięciu zostało %d pozycji, oczekiwane %d", zostalo, trzymane)
	}

	// Zostać mają najnowsze; przycięcie od złej strony zabrałoby bieżącą pracę.
	pozycje, err := drzewo.zestaw.Historia.Pozycje(ctx, drzewo.oknoKod, 0, 0)
	if err != nil {
		t.Fatalf("nie można odczytać historii: %v", err)
	}
	if len(pozycje) != trzymane {
		t.Fatalf("odczyt oddał %d pozycji", len(pozycje))
	}
}

// TestRetencjaBezZasadyNiczegoNieTyka pilnuje wartości domyślnej: brak zasady
// znaczy brak przycinania, a nie przycinanie do zera.
func TestRetencjaBezZasadyNiczegoNieTyka(t *testing.T) {
	drzewo := zalozDrzewo(t, "retencja-brak")
	drzewo.dopiszWiadomosci(t, 4)
	ctx := context.Background()

	if _, _, err := func() (ZasadaPrzechowywania, bool, error) {
		return drzewo.zestaw.Historia.ZasadaOkna(ctx, drzewo.oknoKod)
	}(); err != nil {
		t.Fatalf("nie można rozstrzygnąć zasady: %v", err)
	}

	usuniete, err := drzewo.zestaw.Historia.Egzekwuj(ctx, drzewo.oknoKod)
	if err != nil {
		t.Fatalf("egzekucja bez zasady nie powiodła się: %v", err)
	}
	if usuniete != 0 {
		t.Errorf("egzekucja bez zasady przycięła %d pozycji", usuniete)
	}
	zostalo, err := drzewo.zestaw.Historia.Policz(ctx, drzewo.oknoKod)
	if err != nil {
		t.Fatalf("nie można policzyć historii: %v", err)
	}
	if zostalo != 4 {
		t.Errorf("po egzekucji bez zasady zostało %d pozycji z 4", zostalo)
	}
}

// TestZasadaOknaBijeSesyjnaISesyjnaGlobalna sprawdza rozstrzyganie zakresu.
// Zasady się nie sumują — obowiązuje jedna, najbliższa oknu. Suma dawałaby
// wynik, którego Operator nie przewidziałby z żadnego pojedynczego ekranu.
func TestZasadaOknaBijeSesyjnaISesyjnaGlobalna(t *testing.T) {
	drzewo := zalozDrzewo(t, "zakresy")
	ctx := context.Background()

	globalna, sesyjna, oknowa := 9, 5, 1
	zapisz := func(zakres, kod string, pozycje int) {
		t.Helper()
		if _, err := drzewo.zestaw.Historia.ZapiszZasade(ctx, ZasadaPrzechowywania{
			Zakres: zakres, ZakresKod: kod, PozycjeTrzymane: &pozycje,
		}); err != nil {
			t.Fatalf("nie można zapisać zasady zakresu %s: %v", zakres, err)
		}
	}

	zapisz("global", "", globalna)
	zasada, jest, err := drzewo.zestaw.Historia.ZasadaOkna(ctx, drzewo.oknoKod)
	if err != nil || !jest {
		t.Fatalf("zasada globalna nie została rozstrzygnięta: jest=%t err=%v", jest, err)
	}
	if zasada.Zakres != "global" {
		t.Errorf("rozstrzygnięto zakres %q, oczekiwany globalny", zasada.Zakres)
	}

	zapisz("session", drzewo.sesjaKod, sesyjna)
	zasada, _, err = drzewo.zestaw.Historia.ZasadaOkna(ctx, drzewo.oknoKod)
	if err != nil {
		t.Fatalf("nie można rozstrzygnąć zasady: %v", err)
	}
	if zasada.Zakres != "session" {
		t.Errorf("zasada sesyjna nie przebiła globalnej: rozstrzygnięto %q", zasada.Zakres)
	}

	zapisz("window", drzewo.oknoKod, oknowa)
	zasada, _, err = drzewo.zestaw.Historia.ZasadaOkna(ctx, drzewo.oknoKod)
	if err != nil {
		t.Fatalf("nie można rozstrzygnąć zasady: %v", err)
	}
	if zasada.Zakres != "window" {
		t.Errorf("zasada okna nie przebiła sesyjnej: rozstrzygnięto %q", zasada.Zakres)
	}
	if zasada.PozycjeTrzymane == nil || *zasada.PozycjeTrzymane != oknowa {
		t.Errorf("zasada okna niesie próg %v, oczekiwany %d", zasada.PozycjeTrzymane, oknowa)
	}
}

// TestZasadaDlaBytuKtoregoNieMaJestOdmowa pilnuje bramy zapisu: zasada dla okna
// albo sesji, której nie ma, byłaby nastawą bez skutku — Operator zobaczyłby ją
// w oknie i uznał, że coś robi.
func TestZasadaDlaBytuKtoregoNieMaJestOdmowa(t *testing.T) {
	drzewo := zalozDrzewo(t, "byt")
	ctx := context.Background()

	for _, przypadek := range []struct{ zakres, kod string }{
		{"window", "okno-ktorego-nie-ma"},
		{"session", "sesja-ktorej-nie-ma"},
	} {
		t.Run(przypadek.zakres, func(t *testing.T) {
			istnieje, err := drzewo.zestaw.Historia.IstniejeByt(ctx, przypadek.zakres, przypadek.kod)
			if err != nil {
				t.Fatalf("nie można sprawdzić bytu: %v", err)
			}
			if istnieje {
				t.Errorf("byt %s/%s uznany za istniejący", przypadek.zakres, przypadek.kod)
			}
		})
	}

	istnieje, err := drzewo.zestaw.Historia.IstniejeByt(ctx, "global", "")
	if err != nil {
		t.Fatalf("nie można sprawdzić bytu globalnego: %v", err)
	}
	if !istnieje {
		t.Error("zakres globalny uznany za nieistniejący — nie wskazuje bytu, więc jest zawsze")
	}
}

// TestEgzekucjaTykaWylacznieWskazanegoOkna sprawdza szczelność przycinania:
// zasada okna nie może zabrać historii okna sąsiedniego.
func TestEgzekucjaTykaWylacznieWskazanegoOkna(t *testing.T) {
	drzewo := zalozDrzewo(t, "szczelnosc")
	drzewo.dopiszWiadomosci(t, 4)
	ctx := context.Background()

	// Drugie okno tej samej sesji, z własną historią.
	var modulID, kanalID int64
	if err := drzewo.baza.DB.QueryRow("SELECT id FROM modul ORDER BY id LIMIT 1").Scan(&modulID); err != nil {
		t.Fatalf("nie można odczytać modułu: %v", err)
	}
	if err := drzewo.baza.DB.QueryRow("SELECT id FROM kanal_modelu ORDER BY id LIMIT 1").Scan(&kanalID); err != nil {
		t.Fatalf("nie można odczytać kanału modelu: %v", err)
	}
	sasiedniKod := "okno-sasiednie"
	sasiedniID, err := drzewo.zestaw.Okna.Utworz(ctx, Okno{
		SesjaID: drzewo.sesjaID, ModulID: modulID, KanalModeluID: kanalID,
		SrodowiskoWykonania:     shared.ExecutionEnvLocal,
		TrybUprawnien:           shared.PermissionModeManual,
		RolaOkna:                shared.WindowRoleStandalone,
		Stan:                    shared.WindowStatusOpen,
		IdentyfikatorZewnetrzny: &sasiedniKod,
	})
	if err != nil {
		t.Fatalf("nie można założyć okna sąsiedniego: %v", err)
	}
	tresc := "wypowiedź okna sąsiedniego"
	kodWiadomosci := "sasiednie-wiadomosc"
	if _, err := drzewo.zestaw.Wiadomosci.Dopisz(ctx, Wiadomosc{
		OknoID: sasiedniID, Rola: shared.MessageRoleUser,
		RodzajTresci: shared.ChunkKindText, Stan: shared.MessageStatusComplete,
		Tresc: &tresc, IdentyfikatorZewnetrzny: &kodWiadomosci,
	}); err != nil {
		t.Fatalf("nie można dopisać wiadomości okna sąsiedniego: %v", err)
	}

	trzymane := 1
	if _, err := drzewo.zestaw.Historia.ZapiszZasade(ctx, ZasadaPrzechowywania{
		Zakres: "window", ZakresKod: drzewo.oknoKod, PozycjeTrzymane: &trzymane,
	}); err != nil {
		t.Fatalf("nie można zapisać zasady: %v", err)
	}
	if _, err := drzewo.zestaw.Historia.Egzekwuj(ctx, drzewo.oknoKod); err != nil {
		t.Fatalf("egzekucja nie powiodła się: %v", err)
	}

	sasiednie, err := drzewo.zestaw.Historia.Policz(ctx, sasiedniKod)
	if err != nil {
		t.Fatalf("nie można policzyć historii okna sąsiedniego: %v", err)
	}
	if sasiednie != 1 {
		t.Errorf("egzekucja zasady jednego okna zabrała historię drugiego: zostało %d z 1", sasiednie)
	}
}
