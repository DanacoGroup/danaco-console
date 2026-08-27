package core

import (
	"database/sql"
	"testing"

	"danacoconsole/shared"
)

// Skutek zakresów narzędzi: czy tools.scope.set zapisuje wiersz, który realnie powstrzymuje wywołanie.

// profilZakresuSprawdzianu zakłada profil asystenta wprost w bazie i oddaje jego
// kod. Komendy `assistant.profile.*` w kontrakcie nie ma, więc profil nie ma jak
// powstać drogą komendy — a zakres bez profilu nie ma do czego przylgnąć.
func profilZakresuSprawdzianu(t *testing.T, baza *sql.DB, kod string, domyslny bool) string {
	t.Helper()

	znacznik := 0
	if domyslny {
		znacznik = 1
	}
	if _, err := baza.Exec(
		`INSERT INTO profil_asystenta (identyfikator_zewnetrzny, nazwa, domyslny,
		     utworzono, zaktualizowano)
		 VALUES (?, ?, ?, 0, 0)`, kod, "Profil sprawdzianu", znacznik); err != nil {
		t.Fatalf("nie można założyć profilu asystenta: %v", err)
	}
	return kod
}

// TestSkutekZakresuNarzedziaOdmawiaWywolania mierzy obie rzeczy naraz: zapis
// zakresu w bazie i odmowę, którą ten zapis wywołuje przy wykonaniu.
func TestSkutekZakresuNarzedziaOdmawiaWywolania(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	profil := profilZakresuSprawdzianu(t, baza, "profil-sprawdzianu", true)

	nazwa := zrodloPlatformy + ":" + string(shared.CommandSessionList)
	wylaczona := false
	var zapis shared.ToolsScopeSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandToolsScopeSet,
		shared.ToolsScopeSetRequest{ProfileId: profil, ToolName: nazwa, Enabled: &wylaczona}, &zapis)

	if zapis.Scope.Enabled {
		t.Fatal("odpowiedź oddaje pozycję jako dostępną mimo wyłączenia")
	}
	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM narzedzie_zakres_profilu z JOIN profil_asystenta p ON p.id = z.profil_id
		  WHERE p.identyfikator_zewnetrzny = ? AND z.nazwa_pelna = ? AND z.dostepne = 0`,
		profil, nazwa) != 1 {
		t.Fatal("zakres nie doszedł do bazy")
	}

	// Straż jest tą samą, którą pyta rdzeń przed skierowaniem komendy ręki modelu do obsługiwacza.
	straz := NowyPortZakresowNarzedzi(zmontowany.dane.ZakresyNarzedzi, zmontowany.dane.Asystent)
	if err := straz.SprawdzWywolanie(zycie, shared.CommandSessionList, ""); err == nil {
		t.Fatal("straż przepuściła wywołanie pozycji wyłączonej — zakres jest napisem, nie regułą")
	}
	// Pozycja bez wiersza zakresu przechodzi: stanem wyjściowym jest pełny dostęp bez granicy.
	if err := straz.SprawdzWywolanie(zycie, shared.CommandAgentList, ""); err != nil {
		t.Fatalf("straż odmówiła pozycji, dla której Operator niczego nie zapisał: %v", err)
	}
}

// TestSkutekLimituWywolanNarzedzia mierzy granicę wywołań w oknie czasu:
// zużycie ma rosnąć wierszami rachunku, a po wyczerpaniu granicy wywołanie ma
// dostać odmowę.
func TestSkutekLimituWywolanNarzedzia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	profil := profilZakresuSprawdzianu(t, baza, "profil-limitu", true)

	nazwa := zrodloPlatformy + ":" + string(shared.CommandSessionList)
	dwa := 2
	wykonajUdana(t, zmontowany, zycie, shared.CommandToolsScopeSet,
		shared.ToolsScopeSetRequest{ProfileId: profil, ToolName: nazwa, CallLimit: &dwa}, nil)

	straz := NowyPortZakresowNarzedzi(zmontowany.dane.ZakresyNarzedzi, zmontowany.dane.Asystent)
	for numer := 1; numer <= 2; numer++ {
		if err := straz.SprawdzWywolanie(zycie, shared.CommandSessionList, ""); err != nil {
			t.Fatalf("wywołanie %d z dwóch dozwolonych dostało odmowę: %v", numer, err)
		}
	}
	if err := straz.SprawdzWywolanie(zycie, shared.CommandSessionList, ""); err == nil {
		t.Fatal("trzecie wywołanie przy granicy dwóch przeszło — limit jest liczbą bez skutku")
	}

	// Rachunek leży w bazie, nie w pamięci procesu: po restarcie rdzenia granica
	// ma obowiązywać dalej.
	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM narzedzie_wywolanie w JOIN profil_asystenta p ON p.id = w.profil_id
		  WHERE p.identyfikator_zewnetrzny = ? AND w.nazwa_pelna = ?`, profil, nazwa) != 2 {
		t.Fatal("rachunek wywołań nie doszedł do bazy")
	}

	// Wykaz pokazuje zużycie — bez niego okno nie wie, ile z granicy zostało.
	var wykaz shared.ToolsScopeListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandToolsScopeList,
		shared.ToolsScopeListRequest{ProfileId: &profil}, &wykaz)
	if wykaz.Total != 1 || len(wykaz.Scopes) != 1 {
		t.Fatalf("wykaz zakresów liczy %d pozycji zamiast jednej", wykaz.Total)
	}
	if wykaz.Scopes[0].CallsUsed == nil || *wykaz.Scopes[0].CallsUsed != 2 {
		t.Fatalf("wykaz nie niesie zużycia limitu: %v", wykaz.Scopes[0].CallsUsed)
	}
	if wykaz.Scopes[0].ShortName != string(shared.CommandSessionList) {
		t.Errorf("wykaz niesie nazwę skróconą %q", wykaz.Scopes[0].ShortName)
	}
}

// TestSkutekZakresuNarzedziaProfilDomyslny sprawdza, że żądanie bez wskazania
// profilu bierze profil domyślny, a wskazanie nieznane wraca odmową — zakres
// zapisany pod profilem, którego nie ma, nie obowiązywałby nigdy.
func TestSkutekZakresuNarzedziaProfilDomyslny(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	profilZakresuSprawdzianu(t, baza, "profil-domyslny", true)

	var wykaz shared.ToolsScopeListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandToolsScopeList,
		shared.ToolsScopeListRequest{}, &wykaz)
	if wykaz.Total != 0 {
		t.Fatalf("profil bez zapisanych zakresów oddał %d pozycji", wykaz.Total)
	}

	nieznany := "profil-ktorego-nie-ma"
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandToolsScopeSet,
		shared.ToolsScopeSetRequest{ProfileId: nieznany, ToolName: "danaco:session.list"})
	if odmowa.Code == "" {
		t.Error("zakres pod nieznanym profilem zapisał się bez odmowy")
	}
}
