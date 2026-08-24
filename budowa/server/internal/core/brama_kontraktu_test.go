package core

import (
	"strings"
	"testing"

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
