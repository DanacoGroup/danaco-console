package core

import (
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// STRAŻE BRAMY KONTRAKTU.
//
// Rdzeń wykonuje żądanie dopiero po sprawdzeniu go wobec kontraktu. Sprawdziany
// w tym pliku pilnują dwóch rzeczy naraz, bo obie znoszą się nawzajem:
//
//  1. żądanie niezgodne z kontraktem ma WRÓCIĆ ODMOWĄ, a nie powodzeniem
//     z wartością dobraną przez rdzeń;
//  2. odmowa ma NAZWAĆ pole, o które idzie — odmowa mówiąca „coś jest nie tak"
//     zostawia wołającego dokładnie tam, gdzie zostawiało go milczenie.
//
// Ładunki są tu wypisane mapą, nie strukturą kontraktu: struktura Go niesie
// pole wymagane zawsze (znacznik bez `omitempty`), więc żądania BEZ pola nie da
// się nią złożyć. Żądanie niepełne przychodzi z drutu i tylko mapą da się je
// odtworzyć.

// odmowaNazywa przerywa sprawdzian, gdy odmowa nie ma spodziewanego kodu albo
// nie nazywa wszystkich oczekiwanych członów treści.
func odmowaNazywa(t *testing.T, blad shared.ErrorInfo, kod shared.ErrorCode, czlony ...string) {
	t.Helper()

	if blad.Code != kod {
		t.Errorf("odmowa ma kod %s, oczekiwano %s; treść: %s", blad.Code, kod, blad.Message)
	}
	for _, czlon := range czlony {
		if !strings.Contains(blad.Message, czlon) {
			t.Errorf("odmowa nie nazywa %q; treść: %s", czlon, blad.Message)
		}
	}
}

// TestZadanieBezPolWymaganychWracaOdmowa pilnuje granicy, na której rdzeń
// przyjmował punkt dostępu bez rodzaju, korzeni i trybu domyślnego, po czym
// zakładał punkt rodzaju `mcpBridge` — rodzaju, o który wołający nie prosił.
func TestZadanieBezPolWymaganychWracaOdmowa(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandAccessPointAdd,
		map[string]any{"name": "punkt bez pól wymaganych"})

	odmowaNazywa(t, blad, shared.ErrorCodeValidationFailed, "kind", "roots", "defaultMode")

	// Odmowa ma zostawić świat nietknięty: punkt założony „przy okazji" byłby
	// tą samą szkodą co powodzenie.
	var wykaz shared.AccessPointListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAccessPointList,
		shared.AccessPointListRequest{}, &wykaz)
	for _, punkt := range wykaz.Points {
		if punkt.Name == "punkt bez pól wymaganych" {
			t.Fatalf("odmowa zostawiła po sobie punkt dostępu %s rodzaju %s", punkt.Id, punkt.Kind)
		}
	}
}

// TestWartoscSpozaWyliczeniaWracaOdmowa pilnuje drugiej połowy bramy: wartość
// spoza wyliczenia kontraktu schodziła dotąd na wartość domyślną, a wołający
// dostawał powodzenie i tryb, którego nie ustawił.
func TestWartoscSpozaWyliczeniaWracaOdmowa(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	sesja := zalozSesjeSprawdzianu(t, zmontowany, zycie)
	okno := zalozOknoSprawdzianu(t, zmontowany, zycie, sesja, "kanal-sprawdzianu")

	blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandWindowUpdate,
		map[string]any{"windowId": okno, "permissionMode": "tryb-spoza-kontraktu"})

	odmowaNazywa(t, blad, shared.ErrorCodeValidationFailed,
		"permissionMode", "tryb-spoza-kontraktu")
}

// TestWartoscWyliczeniaZKontraktuPrzechodzi pilnuje, żeby brama nie zamieniła
// się w zaporę: wartość należąca do wyliczenia ma przejść bez przeszkody.
func TestWartoscWyliczeniaZKontraktuPrzechodzi(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	sesja := zalozSesjeSprawdzianu(t, zmontowany, zycie)
	okno := zalozOknoSprawdzianu(t, zmontowany, zycie, sesja, "kanal-sprawdzianu")

	var zmienione shared.WindowUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWindowUpdate,
		map[string]any{"windowId": okno, "permissionMode": string(shared.PermissionModePlan)},
		&zmienione)

	if zmienione.Window.PermissionMode != shared.PermissionModePlan {
		t.Errorf("okno po zmianie ma tryb %s, wskazano %s",
			zmienione.Window.PermissionMode, shared.PermissionModePlan)
	}
}

// TestBrakPolaWymaganegoNazywaPoleNieUsterkeWewnetrzna pilnuje czterech dróg,
// na których brak wartości wracał jako `internal_error`: wołający dostawał tekst
// zapytania SQL albo zdanie o kolumnie bazy zamiast nazwy pola, którego nie
// wypełnił, a kod błędu kazał mu ponowić żądanie, które nie ma prawa się udać.
func TestBrakPolaWymaganegoNazywaPoleNieUsterkeWewnetrzna(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	przypadki := []struct {
		nazwa   string
		komenda shared.MessageType
		ladunek map[string]any
		czlony  []string
	}{
		{"kanał bez nazwy i rodzaju", shared.CommandChannelAdd,
			map[string]any{}, []string{"name", "kind"}},
		{"zmiana kanału bez wskazania kanału", shared.CommandChannelUpdate,
			map[string]any{}, []string{"channelId"}},
		{"źródło bez okna i adresu", shared.CommandBrowserSourceAdd,
			map[string]any{}, []string{"windowId", "url"}},
		{"notatka bez okna i treści", shared.CommandBrowserNoteAdd,
			map[string]any{}, []string{"windowId", "content"}},
	}
	for _, przypadek := range przypadki {
		t.Run(przypadek.nazwa, func(t *testing.T) {
			blad := wykonajOdmowna(t, zmontowany, zycie, przypadek.komenda, przypadek.ladunek)
			odmowaNazywa(t, blad, shared.ErrorCodeValidationFailed, przypadek.czlony...)
		})
	}
}

// TestPustaWartoscPolaWymaganegoNazywaPole pilnuje tej samej granicy od drugiej
// strony. Pole obecne, lecz puste, przechodzi bramę kontraktu — kontrakt żąda
// obecności pola, nie jego niepustości — i rozstrzyga o nim dziedzina. Bez tego
// sprawdzenia pusty rodzaj kanału dojeżdżał do więzu schematu i wracał treścią
// zapytania SQL.
func TestPustaWartoscPolaWymaganegoNazywaPole(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	przypadki := []struct {
		nazwa   string
		komenda shared.MessageType
		ladunek any
		czlony  []string
	}{
		{"kanał z pustą nazwą i rodzajem", shared.CommandChannelAdd,
			shared.ChannelAddRequest{Name: "", Kind: ""}, []string{"name", "kind"}},
		{"zmiana kanału z pustym wskazaniem", shared.CommandChannelUpdate,
			shared.ChannelUpdateRequest{ChannelId: ""}, []string{"channelId"}},
		{"źródło z pustym adresem", shared.CommandBrowserSourceAdd,
			shared.BrowserSourceAddRequest{WindowId: "", Url: ""}, []string{"windowId", "url"}},
		{"notatka z pustą treścią", shared.CommandBrowserNoteAdd,
			shared.BrowserNoteAddRequest{WindowId: "", Content: ""}, []string{"windowId", "content"}},
	}
	for _, przypadek := range przypadki {
		t.Run(przypadek.nazwa, func(t *testing.T) {
			blad := wykonajOdmowna(t, zmontowany, zycie, przypadek.komenda, przypadek.ladunek)
			odmowaNazywa(t, blad, shared.ErrorCodeValidationFailed, przypadek.czlony...)
		})
	}
}

// ── WYJĄTEK POWITANIA ────────────────────────────────────────────────────────

// TestPowitanieNiepelnePrzechodziBrameIOddajeWersjeProtokolu pilnuje jedynego
// wyjątku spod bramy (rejestr decyzji, pozycja 10).
//
// CZYM SIĘ TO ŁAMIE. Brama objęła `connection.hello`, którego trzy pola
// kontrakt oznacza jako wymagane. Klient sprzed wprowadzenia pola `clientId`
// dostawał więc odmowę zamiast wersji protokołu — a powitanie jest jedynym
// miejscem, z którego klient tę wersję czyta, czyli jedynym, w którym rozpoznaje,
// że jest starszy niż rdzeń. Im starszy klient, tym pewniej nie dowiadywał się,
// dlaczego został odrzucony.
//
// Ładunek idzie mapą pustą, bo struktura kontraktu niesie pola wymagane zawsze
// i żądania BEZ nich nie da się nią złożyć — a niepełne powitanie przychodzi
// właśnie z drutu.
func TestPowitanieNiepelnePrzechodziBrameIOddajeWersjeProtokolu(t *testing.T) {
	zmontowany, zycie := zmontujDoSprawdzenia(t)

	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandConnectionHello,
		map[string]any{})
	if odpowiedz.Error != nil {
		t.Fatalf("powitanie niepełne odmówiło: kod=%s treść=%s",
			odpowiedz.Error.Code, odpowiedz.Error.Message)
	}

	var tresc shared.ConnectionHelloResponse
	if err := protocol.LadunekDo(odpowiedz, &tresc); err != nil {
		t.Fatalf("nieczytelny ładunek powitania: %v", err)
	}
	if tresc.ProtocolVersion != shared.ProtocolVersion {
		t.Errorf("powitanie niepełne podaje wersję protokołu %q, kontrakt mówi %q",
			tresc.ProtocolVersion, shared.ProtocolVersion)
	}
	if tresc.ServerVersion != WersjaRdzenia {
		t.Errorf("powitanie niepełne podaje wersję rdzenia %q, rdzeń ma %q",
			tresc.ServerVersion, WersjaRdzenia)
	}
}

// TestBrakiPowitaniaIdaDoDziennikaRdzenia pilnuje drugiej połowy pozycji 10:
// braki pól są odnotowane, a nie przemilczane.
//
// Dziennik jest jedynym miejscem, do którego mogą dojść: `ConnectionHelloResponse`
// nie ma pola na wykaz braków, a dołożenie takiego pola jest zmianą kontraktu.
// Sprawdzian mierzy więc zapis, a nie treść odpowiedzi — i tym samym pilnuje,
// żeby wyjątek nie zamienił się w milczenie.
func TestBrakiPowitaniaIdaDoDziennikaRdzenia(t *testing.T) {
	zmontowany, zycie, dziennik := zmontujZDziennikiem(t)

	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandConnectionHello,
		map[string]any{"clientId": "sprawdzian"})
	if odpowiedz.Error != nil {
		t.Fatalf("powitanie niepełne odmówiło: kod=%s treść=%s",
			odpowiedz.Error.Code, odpowiedz.Error.Message)
	}

	zapis := dziennik.Tresc()
	for _, czlon := range []string{
		string(shared.CommandConnectionHello), "clientVersion", "protocolVersion",
	} {
		if !strings.Contains(zapis, czlon) {
			t.Errorf("dziennik rdzenia nie nazywa %q; zapis:\n%s", czlon, zapis)
		}
	}
	// Pole podane nie ma prawa trafić do wykazu braków — inaczej wpis myliłby
	// czytającego dziennik co do tego, czego klient nie przysłał.
	if strings.Contains(zapis, "clientId") {
		t.Errorf("dziennik liczy pole podane jako brakujące; zapis:\n%s", zapis)
	}
	// Braki nie mają prawa wyjść do wołającego — kontrakt nie ma na nie pola.
	if strings.Contains(string(odpowiedz.Payload), "clientVersion") {
		t.Errorf("odpowiedź powitania niesie wykaz braków: %s", string(odpowiedz.Payload))
	}
}

// TestPowitaniePelneNieZostawiaSladuWDzienniku pilnuje, żeby wyjątek nie
// zamienił dziennika w zapis każdego nawiązania połączenia. Powitanie pełne
// jest przypadkiem zwykłym i wpis o nim zasypywałby powody, dla których dziennik
// się czyta.
func TestPowitaniePelneNieZostawiaSladuWDzienniku(t *testing.T) {
	zmontowany, zycie, dziennik := zmontujZDziennikiem(t)

	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandConnectionHello,
		powitanieSprawdzianu())
	if odpowiedz.Error != nil {
		t.Fatalf("powitanie pełne odmówiło: kod=%s treść=%s",
			odpowiedz.Error.Code, odpowiedz.Error.Message)
	}
	if zapis := dziennik.Tresc(); strings.Contains(zapis, "brama kontraktu") {
		t.Errorf("powitanie pełne zostawiło wpis bramy w dzienniku:\n%s", zapis)
	}
}

// TestWyjatekObejmujeWylaczniePowitanie pilnuje granicy wyjątku od drugiej
// strony: sąsiadka powitania w tej samej rodzinie kontraktu przechodzi bramę
// bez ustępstw i odmawia, nazywając brakujące pole.
//
// Bez tego sprawdzianu wyjątek dopisany dla powitania mógłby rozlać się na całą
// rodzinę `connection.*` albo na komendy wołane przed zalogowaniem, a nikt by
// tego nie zobaczył — odmowa zniknięta wygląda jak działanie.
func TestWyjatekObejmujeWylaczniePowitanie(t *testing.T) {
	zmontowany, zycie := zmontujDoSprawdzenia(t)

	przypadki := []struct {
		nazwa   string
		komenda shared.MessageType
		czlon   string
	}{
		{"wejście bez metody", shared.CommandAuthLogin, "method"},
		{"strona główna bez klienta", shared.CommandHomeEnter, "clientId"},
		{"rejestracja bez loginu, adresu i hasła", shared.CommandAuthRegister, "login"},
	}
	for _, przypadek := range przypadki {
		t.Run(przypadek.nazwa, func(t *testing.T) {
			blad := wykonajOdmowna(t, zmontowany, zycie, przypadek.komenda, map[string]any{})
			odmowaNazywa(t, blad, shared.ErrorCodeValidationFailed, przypadek.czlon)
		})
	}
}
