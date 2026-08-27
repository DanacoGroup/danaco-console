// Zgodność rdzenia z kontraktem: kontrakt jest jedynym źródłem prawdy nazw,
// a jedynym pomiarem mówiącym prawdę o obsłudze komendy jest zmontowany
// rejestr, nie grep źródeł.
package core

import (
	"context"
	"sort"
	"strings"
	"testing"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// granicaKomendySprawdzianu jest granicą czasu JEDNEGO wywołania komendy przez
// uprząż sprawdzianu. Wypada wyłącznie wtedy, gdy rdzeń zwisł — nigdy wtedy,
// gdy komenda po prostu długo pracuje.
const granicaKomendySprawdzianu = granicaWykazuUrzadzen + 15*time.Second

// komendyBezObslugiwacza wylicza komendy kontraktu, których rdzeń dziś nie
// obsługuje. Wykaz jest zaporą, nie zgodą, i stoi dziś PUSTY: każda komenda
// kontraktu ma w rdzeniu obsługiwacza.
var komendyBezObslugiwacza = []shared.MessageType{}

// TestRejestrPokrywaKomendyKontraktu sprawdza, że każda komenda kontraktu ma
// w rdzeniu obsługiwacza — poza wyliczonymi wprost powyżej.
func TestRejestrPokrywaKomendyKontraktu(t *testing.T) {
	zmontowany, _ := zmontujDoSprawdzenia(t)
	obslugiwane := zbiorNazw(zmontowany.Rdzen.rejestr.Nazwy())
	uznane := zbiorNazw(komendyBezObslugiwacza)

	komendy := shared.WszystkieKomendy()
	brakujace := make([]string, 0)
	for _, komenda := range komendy {
		if !obslugiwane[komenda] && !uznane[komenda] {
			brakujace = append(brakujace, komenda)
		}
	}
	sort.Strings(brakujace)

	if len(brakujace) > 0 {
		t.Errorf("komendy kontraktu bez obsługiwacza w rdzeniu: %d z %d\n%s",
			len(brakujace), len(komendy), strings.Join(brakujace, "\n"))
	}

	// Druga strona zapory: wiersz wykazu, który przestał opisywać brak.
	nadmiarowe := make([]string, 0)
	for _, komenda := range komendyBezObslugiwacza {
		if obslugiwane[komenda] {
			nadmiarowe = append(nadmiarowe, komenda)
		}
	}
	sort.Strings(nadmiarowe)

	if len(nadmiarowe) > 0 {
		t.Errorf("komendy uznane za nieobsłużone mają już obsługiwacza — zdejmij je z wykazu komendyBezObslugiwacza: %d\n%s",
			len(nadmiarowe), strings.Join(nadmiarowe, "\n"))
	}
}

// TestRejestrNieMaNazwSpozaKontraktu pilnuje drugiej strony tej samej zgodności:
// rdzeń nie obsługuje nazwy, której kontrakt nie zna.
func TestRejestrNieMaNazwSpozaKontraktu(t *testing.T) {
	zmontowany, _ := zmontujDoSprawdzenia(t)
	kontraktowe := zbiorNazw(shared.WszystkieKomendy())

	obce := make([]string, 0)
	for _, nazwa := range zmontowany.Rdzen.rejestr.Nazwy() {
		if !kontraktowe[nazwa] {
			obce = append(obce, nazwa)
		}
	}
	sort.Strings(obce)

	if len(obce) > 0 {
		t.Errorf("rdzeń obsługuje nazwy spoza kontraktu: %d\n%s",
			len(obce), strings.Join(obce, "\n"))
	}
}

// TestPowitanieOddajeWykazZRejestru sprawdza obietnicę z komentarza powitania:
// klient dostaje wykaz komend rzeczywiście obsługiwanych, nie wykaz z kontraktu.
func TestPowitanieOddajeWykazZRejestru(t *testing.T) {
	zmontowany, zycie := zmontujDoSprawdzenia(t)

	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandConnectionHello, powitanieSprawdzianu())
	if odpowiedz.Error != nil {
		t.Fatalf("powitanie odmówiło: %+v", *odpowiedz.Error)
	}

	var tresc shared.ConnectionHelloResponse
	if err := protocol.LadunekDo(odpowiedz, &tresc); err != nil {
		t.Fatalf("nieczytelny ładunek powitania: %v", err)
	}
	if tresc.ProtocolVersion != shared.ProtocolVersion {
		t.Errorf("powitanie podaje wersję protokołu %q, kontrakt mówi %q",
			tresc.ProtocolVersion, shared.ProtocolVersion)
	}
	if got, want := len(tresc.Commands), zmontowany.Rdzen.rejestr.Liczba(); got != want {
		t.Errorf("powitanie oddaje %d komend, rejestr ma %d", got, want)
	}
	wRejestrze := zbiorNazw(zmontowany.Rdzen.rejestr.Nazwy())
	for _, komenda := range tresc.Commands {
		if !wRejestrze[komenda] {
			t.Errorf("powitanie zapowiada komendę %q, której rejestr nie ma", komenda)
		}
	}
}

// TestKomendaSpozaKontraktuWracaJakoNieznana pilnuje ścieżki opisanej
// w rejestrze rdzenia: typ spoza kontraktu dostaje zdarzenie `*.unknown`
// swojego obszaru wraz ze stanem błędu i kodem `not_found`.
func TestKomendaSpozaKontraktuWracaJakoNieznana(t *testing.T) {
	zmontowany, zycie := zmontujDoSprawdzenia(t)

	odpowiedz := wykonajKomende(t, zmontowany, zycie, "session.wymyslona", struct{}{})
	if oczekiwane := shared.ZdarzenieNieznanej("session.wymyslona"); odpowiedz.Type != oczekiwane {
		t.Errorf("odpowiedź ma typ %q, kontrakt wskazuje %q", odpowiedz.Type, oczekiwane)
	}
	if odpowiedz.Status == nil || *odpowiedz.Status != shared.EnvelopeStatusError {
		t.Errorf("odmowa bez stanu błędu — klient stanie w wiecznym ładowaniu: %+v", odpowiedz.Status)
	}
	if odpowiedz.Error == nil {
		t.Fatal("odmowa bez pola error")
	}
	if odpowiedz.Error.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q", odpowiedz.Error.Code, shared.ErrorCodeNotFound)
	}

	// Po odmowie rdzeń ma pracować dalej.
	dalsza := wykonajKomende(t, zmontowany, zycie, shared.CommandConnectionHello, powitanieSprawdzianu())
	if dalsza.Error != nil {
		t.Errorf("po nieznanej komendzie powitanie odmówiło: %+v", *dalsza.Error)
	}
}

// powitanieSprawdzianu składa powitanie kompletne wobec kontraktu. Powitanie
// jest tu narzędziem, nie przedmiotem pomiaru.
func powitanieSprawdzianu() shared.ConnectionHelloRequest {
	return shared.ConnectionHelloRequest{
		ClientId:        "sprawdzian",
		ClientVersion:   "0",
		ProtocolVersion: shared.ProtocolVersion,
	}
}

// TestKomunikatNieczytelnyNieZrywaRdzenia podaje na wejście bajty, które nie są
// kopertą. Wejściem jest tu ta sama droga, którą wchodzi transport
// (WykonajSurowe), bo to jedyne miejsce, gdzie rdzeń widzi bajty z gniazda.
func TestKomunikatNieczytelnyNieZrywaRdzenia(t *testing.T) {
	zmontowany, zycie := zmontujDoSprawdzenia(t)

	przypadki := map[string][]byte{
		"bajty spoza JSON":    []byte("to nie jest koperta"),
		"pusty dokument":      []byte("{}"),
		"koperta bez typu":    []byte(`{"id":"jeden","payload":{}}`),
		"ładunek innej maści": []byte(`{"type":"connection.hello","id":"jeden","payload":[1,2,3]}`),
	}
	for nazwa, bajty := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			odpowiedz := zmontowany.Rdzen.WykonajSurowe(zycie, bajty)
			if len(odpowiedz) == 0 {
				t.Fatal("rdzeń nie oddał żadnej odpowiedzi — klient zostaje bez rozstrzygnięcia")
			}
			koperta, err := protocol.Odkoduj(odpowiedz)
			if err != nil {
				t.Fatalf("odpowiedź rdzenia sama jest nieczytelna: %v", err)
			}
			if koperta.Status == nil {
				t.Error("odpowiedź bez pola status")
			}
		})
	}
}

// TestKazdaKomendaZnosiPustyLadunek wywołuje wszystkie zarejestrowane komendy
// z ładunkiem pustym. Sprawdzian ocenia, że obsługiwacz nie przerywa wykonania
// i że odmowa jest odmową nazwaną: koperta ze stanem i kodem kontraktu.
func TestKazdaKomendaZnosiPustyLadunek(t *testing.T) {
	zmontowany, zycie := zmontujDoSprawdzenia(t)

	for _, komenda := range zmontowany.Rdzen.rejestr.Nazwy() {
		t.Run(komenda, func(t *testing.T) {
			ctx, przerwij := context.WithTimeout(zycie, granicaKomendySprawdzianu)
			defer przerwij()

			odpowiedz := zmontowany.Rdzen.Wykonaj(ctx, protocol.Koperta{
				Type: komenda,
				Id:   "sprawdzian",
			})
			if odpowiedz.Type == "" {
				t.Fatal("odpowiedź bez typu")
			}
			if odpowiedz.Status == nil {
				t.Fatal("odpowiedź bez stanu — obietnica wywołania nie zostaje rozstrzygnięta")
			}
			if *odpowiedz.Status != shared.EnvelopeStatusError {
				return
			}
			if odpowiedz.Error == nil {
				t.Fatal("stan błędu bez pola error")
			}
			if _, znany := shared.KodyPonawialne[odpowiedz.Error.Code]; !znany {
				t.Errorf("odmowa niesie kod spoza kontraktu: %q", odpowiedz.Error.Code)
			}
			if strings.TrimSpace(odpowiedz.Error.Message) == "" {
				t.Error("odmowa bez treści — Operator nie dowie się, czego brakuje")
			}
		})
	}
}

// wykonajKomende składa kopertę i przepuszcza ją przez rdzeń tą samą drogą,
// którą wchodzi transport, wraz z drogą zwrotną odpowiedzi do wołającego.
func wykonajKomende(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	komenda shared.MessageType, ladunek any) protocol.Koperta {
	t.Helper()

	koperta, err := protocol.NowaKoperta(komenda, "sprawdzian", "", ladunek)
	if err != nil {
		t.Fatalf("nie można złożyć koperty %s: %v", komenda, err)
	}
	ctx, przerwij := context.WithTimeout(zycie, granicaKomendySprawdzianu)
	defer przerwij()
	return zmontowany.Rdzen.Wykonaj(ctx, koperta)
}

// zbiorNazw zamienia wykaz nazw na zbiór do sprawdzeń przynależności, żeby
// porównanie dwóch wykazów nie zależało od kolejności ich elementów.
func zbiorNazw(nazwy []shared.MessageType) map[shared.MessageType]bool {
	zbior := make(map[shared.MessageType]bool, len(nazwy))
	for _, nazwa := range nazwy {
		zbior[nazwa] = true
	}
	return zbior
}
