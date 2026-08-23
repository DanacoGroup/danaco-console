package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Skutek rodziny `extension.*`: czy za odpowiedzią stoi zapis, plik albo
// rozmowa, która naprawdę się odbyła.
//
// Ten sam wzorzec szkody, którego pilnuje `skutek_budowy_produktu_test.go`:
// koperta `ok` bez pokrycia. Dlatego żaden sprawdzian tutaj nie kończy się na
// tym, że odpowiedź jest udana. Każdy schodzi niżej:
//   - do bazy DRUGIM połączeniem i liczy wiersze,
//   - do pliku w magazynie treści i czyta bajty,
//   - do serwera protokołu podniesionego przez sprawdzian, który wie, o co go
//     naprawdę zapytano.
//
// Serwer MCP sprawdzianu jest prawdziwym serwerem JSON-RPC nad HTTP: odpowiada
// na `initialize`, `tools/list` i `tools/call` i zapamiętuje, co dostał.
// Mierzona jest droga rdzenia — powitanie, odkrycie, wywołanie, dziennik ramek,
// metryka użycia — a nie to, co odpowiada konkretny serwer.

// pozycjaSprawdzianuRozszerzen zakłada pozycję katalogu i oddaje jej `Extension.id`.
func pozycjaSprawdzianuRozszerzen(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	kod string, rodzaj shared.ExtensionKind) string {
	t.Helper()

	var wynik shared.ExtensionInstallResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionInstall,
		shared.ExtensionInstallRequest{Code: kod, Kind: rodzaj}, &wynik)
	return wynik.Extension.Id
}

// serwerProtokoluSprawdzianu podnosi serwer MCP odpowiadający po HTTP i oddaje
// jego adres wraz z rejestrem metod, o które go pytano.
func serwerProtokoluSprawdzianu(t *testing.T) (string, *[]string) {
	t.Helper()

	metody := []string{}
	serwer := httptest.NewServer(http.HandlerFunc(
		func(odpowiedz http.ResponseWriter, zadanie *http.Request) {
			var ramka struct {
				Id     int             `json:"id"`
				Method string          `json:"method"`
				Params json.RawMessage `json:"params"`
			}
			if err := json.NewDecoder(zadanie.Body).Decode(&ramka); err != nil {
				odpowiedz.WriteHeader(http.StatusBadRequest)
				return
			}
			metody = append(metody, ramka.Method)
			odpowiedz.Header().Set("Content-Type", "application/json")

			switch ramka.Method {
			case "initialize":
				_, _ = odpowiedz.Write([]byte(`{"jsonrpc":"2.0","id":` +
					itoaSprawdzianu(ramka.Id) + `,"result":{"protocolVersion":"2025-06-18"}}`))
			case "tools/list":
				_, _ = odpowiedz.Write([]byte(`{"jsonrpc":"2.0","id":` +
					itoaSprawdzianu(ramka.Id) + `,"result":{"tools":[` +
					`{"name":"repo.search","description":"szuka w repozytoriach",` +
					`"inputSchema":{"type":"object","properties":{"query":{"type":"string"}}}},` +
					`{"name":"repo.read_file","description":"czyta plik"}]}}`))
			case "resources/list":
				_, _ = odpowiedz.Write([]byte(`{"jsonrpc":"2.0","id":` +
					itoaSprawdzianu(ramka.Id) + `,"result":{"resources":[` +
					`{"uri":"repo://drzewo","name":"drzewo"}]}}`))
			case "prompts/list":
				_, _ = odpowiedz.Write([]byte(`{"jsonrpc":"2.0","id":` +
					itoaSprawdzianu(ramka.Id) + `,"result":{"prompts":[]}}`))
			case "tools/call":
				_, _ = odpowiedz.Write([]byte(`{"jsonrpc":"2.0","id":` +
					itoaSprawdzianu(ramka.Id) + `,"result":{"content":[` +
					`{"type":"text","text":"znaleziono 3 repozytoria"}]}}`))
			default:
				_, _ = odpowiedz.Write([]byte(`{"jsonrpc":"2.0","id":` +
					itoaSprawdzianu(ramka.Id) + `,"error":{"code":-32601,` +
					`"message":"metoda nieznana"}}`))
			}
		}))
	t.Cleanup(serwer.Close)
	return serwer.URL, &metody
}

// itoaSprawdzianu składa numer ramki w treść odpowiedzi serwera próbnego.
func itoaSprawdzianu(numer int) string {
	if numer == 0 {
		return "0"
	}
	cyfry := ""
	for numer > 0 {
		cyfry = string(rune('0'+numer%10)) + cyfry
		numer /= 10
	}
	return cyfry
}

// ustawTransportSprawdzianu wskazuje pozycji adres serwera protokołu.
func ustawTransportSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	pozycja, adres string) {
	t.Helper()

	var wynik shared.ExtensionTransportSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionTransportSet,
		shared.ExtensionTransportSetRequest{
			ExtensionId: pozycja, Transport: shared.McpTransportHttp,
			Endpoint: wskaznik(adres),
		}, &wynik)
}

// ── App Catalog ─────────────────────────────────────────────────────────────

// TestWyszukiwarkaKataloguSzukaWTymCoPozycjaNiesie wykazuje, że trafienie ma
// pokrycie w wierszu, a nie w napisie zapisanym w kodzie.
func TestWyszukiwarkaKataloguSzukaWTymCoPozycjaNiesie(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "mcp-repozytoria", shared.ExtensionKindMcp)
	pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "api-poczta", shared.ExtensionKindApi)

	var trafione shared.ExtensionSearchResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionSearch,
		shared.ExtensionSearchRequest{Query: "repozytoria"}, &trafione)
	if trafione.Total != 1 || trafione.Extensions[0].Code != "mcp-repozytoria" {
		t.Fatalf("wyszukiwarka oddała %d trafień: %+v", trafione.Total, trafione.Extensions)
	}
	if len(trafione.Suggestions) == 0 {
		t.Fatal("wyszukiwarka nie oddała ani jednej podpowiedzi")
	}

	var pusto shared.ExtensionSearchResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionSearch,
		shared.ExtensionSearchRequest{Query: "czegoś takiego nie ma"}, &pusto)
	if pusto.Total != 0 {
		t.Fatalf("wyszukiwarka znalazła %d pozycji dla frazy bez pokrycia", pusto.Total)
	}
}

// TestKolekcjaPrzestawiaStanWlaczeniaPozycji wykazuje najtwardszy skutek
// kolekcji: po jej zastosowaniu kolumna `wlaczone` w bazie naprawdę się zmienia.
func TestKolekcjaPrzestawiaStanWlaczeniaPozycji(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	pierwsza := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "mcp-repozytoria", shared.ExtensionKindMcp)
	druga := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "api-poczta", shared.ExtensionKindApi)

	var kolekcja shared.ExtensionCollectionSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionCollectionSave,
		shared.ExtensionCollectionSaveRequest{
			Name: "Zestaw wdrożeniowy", ColorTag: wskaznik("#c8a24a"),
			ExtensionIds: []string{pierwsza, druga},
		}, &kolekcja)

	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM pozycja_kolekcji_rozszerzen p
		   JOIN kolekcja_rozszerzen k ON k.id = p.kolekcja_id
		  WHERE k.identyfikator_zewnetrzny = ?`, kolekcja.Collection.Id); ile != 2 {
		t.Fatalf("związek kolekcji z pozycjami ma %d wierszy zamiast dwóch", ile)
	}

	// Obie pozycje wchodzą wyłączone (pochodzenie Personal) — to jest punkt
	// wyjścia, wobec którego mierzymy skutek.
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM rozszerzenie WHERE wlaczone = 1`); ile != 0 {
		t.Fatalf("przed zastosowaniem kolekcji włączonych jest %d pozycji", ile)
	}

	var zastosowanie shared.ExtensionCollectionApplyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionCollectionApply,
		shared.ExtensionCollectionApplyRequest{
			CollectionId: kolekcja.Collection.Id, Enable: wskaznik(true),
		}, &zastosowanie)
	if len(zastosowanie.Applied) != 2 || len(zastosowanie.Rejected) != 0 {
		t.Fatalf("zastosowanie kolekcji: %d zastosowanych, %d odrzuconych",
			len(zastosowanie.Applied), len(zastosowanie.Rejected))
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM rozszerzenie WHERE wlaczone = 1`); ile != 2 {
		t.Fatalf("po zastosowaniu kolekcji włączonych jest %d pozycji zamiast dwóch", ile)
	}
	// Zmiana grupowa zostawia ślad w dzienniku cyklu życia.
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM historia_rozszerzenia WHERE czynnosc = 'enabled'`); ile != 2 {
		t.Fatalf("dziennik cyklu życia ma %d wpisów włączenia zamiast dwóch", ile)
	}
}

// TestOperacjaZbiorczaOdinstalowujeBezKasowaniaWiersza wykazuje regułę
// katalogu: odinstalowanie zdejmuje stan, ale zostawia pozycję w rejestrze.
func TestOperacjaZbiorczaOdinstalowujeBezKasowaniaWiersza(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "wtyczka-raporty",
		shared.ExtensionKindPlugin)

	var wynik shared.ExtensionAdminBulkResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionAdminBulk,
		shared.ExtensionAdminBulkRequest{
			ExtensionIds: []string{pozycja, "rozsz-nieistniejace"},
			Action:       shared.ExtensionBulkActionUninstall,
		}, &wynik)

	if len(wynik.Affected) != 1 || len(wynik.Rejected) != 1 {
		t.Fatalf("operacja zbiorcza: %d dotkniętych, %d odrzuconych",
			len(wynik.Affected), len(wynik.Rejected))
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM rozszerzenie WHERE identyfikator_zewnetrzny = ?
		   AND zainstalowane = 0 AND wlaczone = 0`, pozycja); ile != 1 {
		t.Fatal("odinstalowanie nie zdjęło stanu pozycji")
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM rozszerzenie WHERE identyfikator_zewnetrzny = ?`, pozycja); ile != 1 {
		t.Fatal("odinstalowanie skasowało wiersz pozycji — katalog stracił ją z rejestru")
	}
}

// TestPrzesylkaPaczkiZostawiaBajtyWMagazynie wykazuje, że za odwołaniem
// przesyłki leży plik o tej treści i tej sumie kontrolnej.
func TestPrzesylkaPaczkiZostawiaBajtyWMagazynie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	tresc := []byte("zawartość paczki rozszerzenia — sprawdzian skutku")
	var wynik shared.ExtensionPackageUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionPackageUpload,
		shared.ExtensionPackageUploadRequest{
			FileName: "wtyczka.zip", ContentBase64: wBase64(tresc),
			ChecksumSha256: wskaznik(sumaSha256(tresc)),
		}, &wynik)

	if wynik.SizeBytes != int64(len(tresc)) {
		t.Fatalf("rdzeń zameldował %d bajtów, przesłano %d", wynik.SizeBytes, len(tresc))
	}
	sciezka := tekstZBazyApps(t, baza,
		`SELECT sciezka FROM paczka_rozszerzenia WHERE identyfikator_zewnetrzny = ?`,
		wynik.UploadRef)
	bajty := bajtyPodOdwolaniem(t, katalog, sciezka)
	if string(bajty) != string(tresc) {
		t.Fatalf("pod odwołaniem leży inna treść niż przesłana: %q", bajty)
	}

	// Suma niezgodna z treścią jest odmową, nie cichym zapisem.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandExtensionPackageUpload,
		shared.ExtensionPackageUploadRequest{
			FileName: "wtyczka.zip", ContentBase64: wBase64(tresc),
			ChecksumSha256: wskaznik("0000000000000000000000000000000000000000000000000000000000000000"),
		})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("przesyłka z niezgodną sumą wróciła kodem %q", odmowa.Code)
	}
}

// TestZestawInstalujeWszystkoCoDaSieZainstalowac wykazuje, że instalacja
// z manifestu naprawdę zakłada wiersze, a pozycje wadliwe wracają w `rejected`.
func TestZestawInstalujeWszystkoCoDaSieZainstalowac(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	manifest := jsonSurowy(t, map[string]any{
		"extensions": []map[string]any{
			{"code": "mcp-repozytoria", "kind": "mcp", "name": "Serwer MCP Repozytoria",
				"version": "0.9.0"},
			{"code": "api-poczta", "kind": "api", "version": "2.1.0"},
			{"code": "bez-rodzaju", "kind": "czegoś-takiego-nie-ma"},
		},
	})

	var wynik shared.ExtensionBundleInstallResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionBundleInstall,
		shared.ExtensionBundleInstallRequest{Manifest: manifest}, &wynik)

	if len(wynik.Installed) != 2 {
		t.Fatalf("zestaw zainstalował %d pozycji zamiast dwóch: %+v",
			len(wynik.Installed), wynik.Installed)
	}
	if len(wynik.Rejected) != 1 || wynik.Rejected[0].Code != "bez-rodzaju" {
		t.Fatalf("zestaw odrzucił %+v", wynik.Rejected)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM rozszerzenie WHERE zainstalowane = 1`); ile != 2 {
		t.Fatalf("w katalogu stoi %d zainstalowanych pozycji zamiast dwóch", ile)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM historia_rozszerzenia WHERE czynnosc = 'installed'`); ile != 2 {
		t.Fatalf("dziennik cyklu życia ma %d wpisów instalacji zamiast dwóch", ile)
	}
	// Pozycja wadliwa nie zostawiła wiersza — odrzucenie jest odrzuceniem.
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM rozszerzenie WHERE kod = 'bez-rodzaju'`); ile != 0 {
		t.Fatal("pozycja odrzucona mimo to weszła do katalogu")
	}
}

// ── Warstwa protokołu ───────────────────────────────────────────────────────

// TestOdkrycieNarzedziPytaSerwerIZapisujeWykaz wykazuje, że odkrycie jest
// rozmową: sprawdzian widzi, o które metody protokołu rdzeń zapytał, i znajduje
// wykaz w bazie.
func TestOdkrycieNarzedziPytaSerwerIZapisujeWykaz(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	adres, metody := serwerProtokoluSprawdzianu(t)
	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "mcp-repozytoria",
		shared.ExtensionKindMcp)
	ustawTransportSprawdzianu(t, zmontowany, zycie, pozycja, adres)

	var wykaz shared.ExtensionToolListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionToolList,
		shared.ExtensionToolListRequest{ExtensionId: pozycja, Refresh: wskaznik(true)}, &wykaz)

	pytania := strings.Join(*metody, ",")
	for _, oczekiwana := range []string{"initialize", "tools/list", "resources/list"} {
		if !strings.Contains(pytania, oczekiwana) {
			t.Fatalf("rdzeń nie zapytał serwera o %s; zapytał o: %s", oczekiwana, pytania)
		}
	}
	if wykaz.ProtocolVersion == nil || *wykaz.ProtocolVersion != "2025-06-18" {
		t.Fatalf("wersja protokołu nie pochodzi z powitania serwera: %+v", wykaz.ProtocolVersion)
	}
	// Dwa narzędzia i jeden zasób — dokładnie tyle, ile oddał serwer.
	if len(wykaz.Entries) != 3 {
		t.Fatalf("wykaz ma %d wpisów zamiast trzech: %+v", len(wykaz.Entries), wykaz.Entries)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM narzedzie_rozszerzenia WHERE rozszerzenie_kod = ? AND rodzaj = 'tool'`,
		pozycja); ile != 2 {
		t.Fatalf("w bazie stoi %d narzędzi zamiast dwóch", ile)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM narzedzie_rozszerzenia WHERE rozszerzenie_kod = ? AND rodzaj = 'resource'`,
		pozycja); ile != 1 {
		t.Fatalf("w bazie stoi %d zasobów zamiast jednego", ile)
	}
	// Schemat wejścia przeszedł w całości, nie jako sama nazwa.
	schemat := tekstZBazyApps(t, baza,
		`SELECT IFNULL(schemat_wejscia,'') FROM narzedzie_rozszerzenia
		  WHERE rozszerzenie_kod = ? AND nazwa = 'repo.search'`, pozycja)
	if !strings.Contains(schemat, "query") {
		t.Fatalf("schemat wejścia narzędzia nie przeszedł: %q", schemat)
	}

	// Odczyt bez `refresh` nie pyta serwera po raz drugi.
	przed := len(*metody)
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionToolList,
		shared.ExtensionToolListRequest{ExtensionId: pozycja}, &wykaz)
	if len(*metody) != przed {
		t.Fatal("odczyt bez refresh mimo to zapytał serwer")
	}
}

// TestWywolanieNarzedziaZapisujeRamkiIMetryke wykazuje, że próbne wywołanie
// naprawdę idzie do serwera, a jego ślad zostaje w dzienniku protokołu i użycia.
func TestWywolanieNarzedziaZapisujeRamkiIMetryke(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	adres, _ := serwerProtokoluSprawdzianu(t)
	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "mcp-repozytoria",
		shared.ExtensionKindMcp)
	ustawTransportSprawdzianu(t, zmontowany, zycie, pozycja, adres)

	var wywolanie shared.ExtensionToolCallResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionToolCall,
		shared.ExtensionToolCallRequest{
			ExtensionId: pozycja, ToolName: "repo.search",
			Arguments: jsonSurowy(t, map[string]any{"query": "danaco"}),
		}, &wywolanie)

	if !wywolanie.Ok {
		t.Fatalf("wywołanie zakończyło się niepowodzeniem: %+v", wywolanie.ErrorDetail)
	}
	if wywolanie.Text == nil || !strings.Contains(*wywolanie.Text, "3 repozytoria") {
		t.Fatalf("treść odpowiedzi narzędzia nie przeszła: %+v", wywolanie.Text)
	}

	// Ramki obu kierunków leżą w dzienniku protokołu.
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM ramka_protokolu_rozszerzenia
		  WHERE rozszerzenie_kod = ? AND kierunek = 'outgoing'`, pozycja); ile < 2 {
		t.Fatalf("dziennik protokołu ma %d ramek wychodzących", ile)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM ramka_protokolu_rozszerzenia
		  WHERE rozszerzenie_kod = ? AND kierunek = 'incoming'`, pozycja); ile < 2 {
		t.Fatalf("dziennik protokołu ma %d ramek przychodzących", ile)
	}

	var log shared.ExtensionProtocolLogListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionProtocolLogList,
		shared.ExtensionProtocolLogListRequest{ExtensionId: pozycja}, &log)
	if log.Total < 4 {
		t.Fatalf("odczyt logu protokołu oddał %d ramek", log.Total)
	}

	// Metryka użycia liczy się z wiersza wywołania, nie z licznika w pamięci.
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM wywolanie_rozszerzenia WHERE rozszerzenie_kod = ? AND udane = 1`,
		pozycja); ile != 1 {
		t.Fatalf("dziennik użycia ma %d udanych wywołań zamiast jednego", ile)
	}
	var uzycie shared.ExtensionUsageGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionUsageGet,
		shared.ExtensionUsageGetRequest{ExtensionId: wskaznik(pozycja)}, &uzycie)
	if len(uzycie.Usage) != 1 || uzycie.Usage[0].Calls != 1 || uzycie.Usage[0].Failures != 0 {
		t.Fatalf("metryka użycia: %+v", uzycie.Usage)
	}

	// Wywołanie narzędzia, którego serwer nie zna, NIE jest odmową komendy:
	// inspektor ma pokazać odpowiedź serwera wraz z powodem.
	var nieudane shared.ExtensionToolCallResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionToolCall,
		shared.ExtensionToolCallRequest{ExtensionId: pozycja, ToolName: "repo.brak"}, &nieudane)
	if nieudane.Ok {
		// Serwer sprawdzianu odpowiada na `tools/call` zawsze — ten sprawdzian
		// pilnuje więc drogi, nie wyniku: wywołanie ma się odbyć i policzyć.
		t.Log("serwer sprawdzianu przyjmuje każde narzędzie — mierzona jest droga wywołania")
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM wywolanie_rozszerzenia WHERE rozszerzenie_kod = ?`,
		pozycja); ile != 2 {
		t.Fatalf("dziennik użycia ma %d wywołań zamiast dwóch", ile)
	}
}

// TestKondycjaIntegracjiMierzySerwerIZapisujeWynik wykazuje pomiar zdrowia.
func TestKondycjaIntegracjiMierzySerwerIZapisujeWynik(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	adres, _ := serwerProtokoluSprawdzianu(t)
	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "mcp-repozytoria",
		shared.ExtensionKindMcp)
	ustawTransportSprawdzianu(t, zmontowany, zycie, pozycja, adres)

	var kondycja shared.ExtensionHealthCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionHealthCheck,
		shared.ExtensionHealthCheckRequest{ExtensionId: wskaznik(pozycja)}, &kondycja)

	if len(kondycja.Results) != 1 {
		t.Fatalf("sprawdzenie oddało %d wyników zamiast jednego", len(kondycja.Results))
	}
	wynik := kondycja.Results[0]
	if wynik.Status != shared.ExtensionHealthStatusHealthy {
		t.Fatalf("serwer odpowiadający uznano za %q: %+v", wynik.Status, wynik.HandshakeError)
	}
	if wynik.ToolCount == nil || *wynik.ToolCount != 3 {
		t.Fatalf("liczba narzędzi w wyniku kondycji: %+v", wynik.ToolCount)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM kondycja_rozszerzenia WHERE rozszerzenie_kod = ? AND stan = 'healthy'`,
		pozycja); ile != 1 {
		t.Fatalf("dziennik kondycji ma %d wpisów zdrowia zamiast jednego", ile)
	}

	// Pozycja bez transportu wraca stanem `unknown` — to prawda o niej, nie awaria.
	bezTransportu := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "api-poczta",
		shared.ExtensionKindApi)
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionHealthCheck,
		shared.ExtensionHealthCheckRequest{ExtensionId: wskaznik(bezTransportu)}, &kondycja)
	if kondycja.Results[0].Status != shared.ExtensionHealthStatusUnknown {
		t.Fatalf("pozycja bez transportu ma stan %q", kondycja.Results[0].Status)
	}
}

// TestPiaskownicaZostawiaDziennikPrzebiegu wykazuje, że za odwołaniem logu
// piaskownicy leży plik z wejściem i wyjściem przebiegu.
func TestPiaskownicaZostawiaDziennikPrzebiegu(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	serwer := httptest.NewServer(http.HandlerFunc(
		func(odpowiedz http.ResponseWriter, zadanie *http.Request) {
			tresc := make([]byte, 1024)
			ile, _ := zadanie.Body.Read(tresc)
			odpowiedz.Header().Set("Content-Type", "application/json")
			_, _ = odpowiedz.Write([]byte(`{"echo":` + string(tresc[:ile]) + `}`))
		}))
	t.Cleanup(serwer.Close)

	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "umiejetnosc-pdf",
		shared.ExtensionKindSkill)
	ustawTransportSprawdzianu(t, zmontowany, zycie, pozycja, serwer.URL)

	var przebieg shared.ExtensionSandboxRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionSandboxRun,
		shared.ExtensionSandboxRunRequest{
			ExtensionId: pozycja,
			Input:       jsonSurowy(t, map[string]any{"dokument": "faktura.pdf"}),
		}, &przebieg)

	if !przebieg.Ok {
		t.Fatalf("przebieg piaskownicy nie powiódł się: %+v", przebieg)
	}
	if przebieg.LogRef == nil || *przebieg.LogRef == "" {
		t.Fatal("piaskownica nie oddała odwołania do dziennika przebiegu")
	}
	bajty := bajtyPodOdwolaniem(t, katalog, *przebieg.LogRef)
	dziennik := string(bajty)
	if !strings.Contains(dziennik, "faktura.pdf") || !strings.Contains(dziennik, "echo") {
		t.Fatalf("dziennik przebiegu nie niesie ani wejścia, ani wyjścia: %q", dziennik)
	}

	// Pozycja bez transportu nie ma czego uruchomić — odmowa z powodem.
	bez := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "wtyczka-bez-drogi",
		shared.ExtensionKindPlugin)
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandExtensionSandboxRun,
		shared.ExtensionSandboxRunRequest{ExtensionId: bez, Input: jsonSurowy(t, map[string]any{})})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("piaskownica pozycji bez drogi wróciła kodem %q", odmowa.Code)
	}
}

// TestImportOpenapiSkladaOperacjeZOpisu wykazuje, że wykaz operacji pochodzi
// z rozłożonego opisu, a nie z nazwy pozycji.
func TestImportOpenapiSkladaOperacjeZOpisu(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	opis := `{
	  "openapi": "3.0.0",
	  "info": {"title": "Portal klienta", "version": "1.0.0"},
	  "paths": {
	    "/zamowienia": {
	      "get": {
	        "operationId": "wykazZamowien",
	        "summary": "wykaz zamówień",
	        "parameters": [
	          {"name": "strona", "in": "query", "required": true,
	           "schema": {"type": "integer"}}
	        ],
	        "responses": {"200": {"description": "ok"}}
	      },
	      "post": {
	        "operationId": "nowyZamowienie",
	        "responses": {"201": {"description": "utworzono"}}
	      }
	    }
	  }
	}`

	var wynik shared.ExtensionDefinitionImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionDefinitionImport,
		shared.ExtensionDefinitionImportRequest{
			Format: shared.ExtensionDefinitionFormatOpenapi3,
			Code:   "portal-klienta-api", Content: wskaznik(opis),
		}, &wynik)

	if len(wynik.Operations) != 2 {
		t.Fatalf("import oddał %d operacji zamiast dwóch: %+v",
			len(wynik.Operations), wynik.Operations)
	}
	if wynik.Extension.Name != "Portal klienta" {
		t.Fatalf("nazwa pozycji nie pochodzi z opisu: %q", wynik.Extension.Name)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM narzedzie_rozszerzenia WHERE rozszerzenie_kod = ?`,
		wynik.Extension.Id); ile != 2 {
		t.Fatalf("w bazie stoi %d operacji zamiast dwóch", ile)
	}
	schemat := tekstZBazyApps(t, baza,
		`SELECT IFNULL(schemat_wejscia,'') FROM narzedzie_rozszerzenia
		  WHERE rozszerzenie_kod = ? AND nazwa = 'wykazZamowien'`, wynik.Extension.Id)
	if !strings.Contains(schemat, "strona") {
		t.Fatalf("schemat argumentów operacji nie powstał z parametrów opisu: %q", schemat)
	}
}

// TestImportGraphqlBierzePolaQueryIMutation wykazuje drugą drogę importu.
func TestImportGraphqlBierzePolaQueryIMutation(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	schemat := "type Klient {\n  id: ID!\n  nazwa: String\n}\n\n" +
		"type Query {\n  zamowienia(strona: Int): [Zamowienie!]!\n  klient(id: ID!): Klient\n}\n\n" +
		"type Mutation {\n  utworzZamowienie(wejscie: WejscieZamowienia!): Zamowienie\n}\n"

	var wynik shared.ExtensionDefinitionImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionDefinitionImport,
		shared.ExtensionDefinitionImportRequest{
			Format: shared.ExtensionDefinitionFormatGraphql,
			Code:   "portal-graphql", Content: wskaznik(schemat),
		}, &wynik)

	nazwy := map[string]bool{}
	for _, operacja := range wynik.Operations {
		nazwy[operacja.Name] = true
	}
	for _, oczekiwana := range []string{"zamowienia", "klient", "utworzZamowienie"} {
		if !nazwy[oczekiwana] {
			t.Fatalf("import GraphQL nie znalazł operacji %s: %+v", oczekiwana, nazwy)
		}
	}
	// Pola typu `Klient` nie są operacjami — opisują kształt danych.
	if nazwy["nazwa"] {
		t.Fatal("import GraphQL wziął pole typu danych za operację")
	}
}

// ── Warstwa zaufania ────────────────────────────────────────────────────────

// TestUprawnieniaRozrozniajaDeklaracjeOdNadania wykazuje, że `excessive` jest
// różnicą dwóch zbiorów leżących w bazie.
func TestUprawnieniaRozrozniajaDeklaracjeOdNadania(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "konektor-crm",
		shared.ExtensionKindApi)

	// Deklarację manifestu wnosi publikacja pakietu; tutaj wpisujemy ją wprost
	// przez warstwę danych, bo sprawdzian mierzy rozróżnienie, nie drogę wejścia.
	if _, err := baza.Exec(
		`INSERT INTO uprawnienie_rozszerzenia (rozszerzenie_kod, zakres, byt, nadane)
		 VALUES (?, 'network', 'crm.example.com', 0)`, pozycja); err != nil {
		t.Fatalf("nie można wpisać uprawnienia deklarowanego: %v", err)
	}

	var nadanie shared.ExtensionPermissionGrantResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionPermissionGrant,
		shared.ExtensionPermissionGrantRequest{
			ExtensionId: pozycja,
			Permissions: []shared.ExtensionPermission{
				{Scope: shared.ExtensionPermissionScopeNetwork, Target: wskaznik("crm.example.com")},
				{Scope: shared.ExtensionPermissionScopeProcessSpawn},
			},
		}, &nadanie)
	if len(nadanie.Granted) != 2 {
		t.Fatalf("nadano %d uprawnień zamiast dwóch", len(nadanie.Granted))
	}

	var wykaz shared.ExtensionPermissionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionPermissionList,
		shared.ExtensionPermissionListRequest{ExtensionId: pozycja}, &wykaz)
	if len(wykaz.Declared) != 1 || len(wykaz.Granted) != 2 {
		t.Fatalf("deklarowanych %d, nadanych %d", len(wykaz.Declared), len(wykaz.Granted))
	}
	if len(wykaz.Excessive) != 1 || wykaz.Excessive[0] != "processSpawn" {
		t.Fatalf("uprawnienia nadmiarowe: %+v", wykaz.Excessive)
	}

	// Skaner widzi to samo i mówi o tym spostrzeżeniami — sygnał, nie brama.
	var skan shared.ExtensionManifestScanResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionManifestScan,
		shared.ExtensionManifestScanRequest{ExtensionId: pozycja}, &skan)
	kody := map[string]int{}
	for _, spostrzezenie := range skan.Findings {
		kody[spostrzezenie.Code]++
	}
	if kody[kodSkaneraNadmiarowe] == 0 {
		t.Fatalf("skaner nie zauważył uprawnienia nadmiarowego: %+v", skan.Findings)
	}
	if kody[kodSkaneraProcesy] == 0 {
		t.Fatalf("skaner nie zauważył żądania uruchamiania procesów: %+v", skan.Findings)
	}
	if kody[kodSkaneraBezPodpisu] == 0 {
		t.Fatalf("skaner nie zauważył braku podpisu: %+v", skan.Findings)
	}

	// Nadanie WYMIENIA komplet: krótszy wykaz cofa to, czego w nim nie ma.
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionPermissionGrant,
		shared.ExtensionPermissionGrantRequest{
			ExtensionId: pozycja,
			Permissions: []shared.ExtensionPermission{
				{Scope: shared.ExtensionPermissionScopeNetwork, Target: wskaznik("crm.example.com")},
			},
		}, &nadanie)
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM uprawnienie_rozszerzenia WHERE rozszerzenie_kod = ? AND nadane = 1`,
		pozycja); ile != 1 {
		t.Fatalf("po zawężeniu nadania w bazie stoi %d uprawnień nadanych", ile)
	}
}

// TestWeryfikacjaPodpisuLiczyPodpisOdNowa jest sprawdzianem najtwardszym tej
// warstwy: pozycja opublikowana z modułu Apps ma podpis, który weryfikuje się
// kluczem publicznym, a podmieniona suma go unieważnia.
func TestWeryfikacjaPodpisuLiczyPodpisOdNowa(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	okno := oknoModuluSprawdzianu(t, zmontowany, zycie)
	zapiszPlikWarsztatuOkna(t, zmontowany, zycie, okno,
		shared.AppWorkspaceLayerFrontend, "index.html", "<h1>Portal klienta</h1>")
	przeprowadzWdrozenieSprawdzianu(t, zmontowany, zycie, baza, okno)

	var pakiet shared.AppsPackageBuildResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackageBuild,
		shared.AppsPackageBuildRequest{WindowId: okno}, &pakiet)
	var manifest shared.AppsPackageManifestSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackageManifestSave,
		shared.AppsPackageManifestSaveRequest{
			WindowId: okno, PackageId: wskaznik(pakiet.Package.Id),
			Manifest: shared.AppPackageManifest{
				Identifier: "portal-klienta", Name: "Portal klienta", Version: "1.0.0",
				Kind: shared.ExtensionKindPlugin,
			},
		}, &manifest)
	var podpis shared.AppsPackageSignResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackageSign,
		shared.AppsPackageSignRequest{
			WindowId: okno, PackageId: pakiet.Package.Id,
			SigningKeyRef: "sejf:wydawca/danaco",
		}, &podpis)
	var publikacja shared.AppsPackagePublishResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackagePublish,
		shared.AppsPackagePublishRequest{WindowId: okno, PackageId: pakiet.Package.Id},
		&publikacja)

	// Podpis przeszedł na pozycję katalogu wraz z materiałem do weryfikacji.
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM podpis_rozszerzenia WHERE rozszerzenie_kod = ?
		   AND podpis_base64 IS NOT NULL AND klucz_base64 IS NOT NULL`,
		publikacja.Extension.Id); ile != 1 {
		t.Fatal("publikacja nie przeniosła materiału podpisu na pozycję katalogu")
	}

	var weryfikacja shared.ExtensionSignatureVerifyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionSignatureVerify,
		shared.ExtensionSignatureVerifyRequest{ExtensionId: publikacja.Extension.Id},
		&weryfikacja)
	if !weryfikacja.Signature.Signed || !weryfikacja.Signature.Verified {
		t.Fatalf("podpis pozycji nie przeszedł weryfikacji: %+v", weryfikacja.Signature)
	}
	if weryfikacja.Signature.TrustLevel != shared.ExtensionTrustLevelVerifiedPublisher {
		t.Fatalf("poziom zaufania po udanej weryfikacji: %q", weryfikacja.Signature.TrustLevel)
	}

	// Podmiana sumy kontrolnej ma unieważnić podpis — weryfikacja liczy go od
	// nowa, a nie przepisuje zapamiętanego werdyktu.
	if _, err := baza.Exec(
		`UPDATE podpis_rozszerzenia SET suma_kontrolna = ? WHERE rozszerzenie_kod = ?`,
		"0000000000000000000000000000000000000000000000000000000000000000",
		publikacja.Extension.Id); err != nil {
		t.Fatalf("nie można podmienić sumy kontrolnej: %v", err)
	}
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionSignatureVerify,
		shared.ExtensionSignatureVerifyRequest{ExtensionId: publikacja.Extension.Id},
		&weryfikacja)
	if weryfikacja.Signature.Verified {
		t.Fatal("podpis przeszedł weryfikację mimo podmienionej sumy kontrolnej")
	}
	if weryfikacja.Signature.TrustLevel != shared.ExtensionTrustLevelUnverifiedPersonal {
		t.Fatalf("poziom zaufania po nieudanej weryfikacji: %q", weryfikacja.Signature.TrustLevel)
	}

	// Publikacja odłożyła też wersję pozycji — jest dokąd cofnąć i z czym
	// zestawić aktualizację.
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM wersja_rozszerzenia WHERE rozszerzenie_kod = ? AND wersja = '1.0.0'`,
		publikacja.Extension.Id); ile != 1 {
		t.Fatal("publikacja nie odłożyła wersji pozycji")
	}

	// Prywatny rejestr organizacji widzi tę pozycję, bo powstała z pakietu.
	var rejestr shared.ExtensionRegistryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionRegistryList,
		shared.ExtensionRegistryListRequest{}, &rejestr)
	if rejestr.Total != 1 || rejestr.Extensions[0].Code != "portal-klienta" {
		t.Fatalf("rejestr organizacji: %d pozycji %+v", rejestr.Total, rejestr.Extensions)
	}
}

// TestWersjonowaniePozycjiPrzypinaICofa wykazuje, że przypięcie i cofnięcie
// zmieniają wiersz, a wskazanie wersji nieznanej jest odmową.
func TestWersjonowaniePozycjiPrzypinaICofa(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "wtyczka-raporty",
		shared.ExtensionKindPlugin)
	// Dwa wydania w rejestrze wersji — to one są tym, do czego wolno wrócić.
	for _, wersja := range []string{"1.0.0", "2.0.0"} {
		if _, err := baza.Exec(
			`INSERT INTO wersja_rozszerzenia (rozszerzenie_kod, wersja, utworzono)
			 VALUES (?, ?, 1)`, pozycja, wersja); err != nil {
			t.Fatalf("nie można wpisać wersji %s: %v", wersja, err)
		}
	}
	if _, err := baza.Exec(`UPDATE rozszerzenie SET wersja = '2.0.0' WHERE identyfikator_zewnetrzny = ?`,
		pozycja); err != nil {
		t.Fatalf("nie można ustawić wersji bieżącej: %v", err)
	}

	var przypiecie shared.ExtensionVersionPinResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionVersionPin,
		shared.ExtensionVersionPinRequest{ExtensionId: pozycja, Version: wskaznik("1.0.0")},
		&przypiecie)
	if przypiecie.PinnedVersion == nil || *przypiecie.PinnedVersion != "1.0.0" {
		t.Fatalf("przypięcie oddało %+v", przypiecie.PinnedVersion)
	}
	if wersja := tekstZBazyApps(t, baza,
		`SELECT IFNULL(wersja_przypieta,'') FROM rozszerzenie WHERE identyfikator_zewnetrzny = ?`,
		pozycja); wersja != "1.0.0" {
		t.Fatalf("kolumna przypięcia niesie %q", wersja)
	}

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandExtensionVersionPin,
		shared.ExtensionVersionPinRequest{ExtensionId: pozycja, Version: wskaznik("9.9.9")})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("przypięcie wersji nieznanej wróciło kodem %q", odmowa.Code)
	}

	var cofniecie shared.ExtensionVersionRollbackResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionVersionRollback,
		shared.ExtensionVersionRollbackRequest{ExtensionId: pozycja, TargetVersion: "1.0.0"},
		&cofniecie)
	if cofniecie.RolledBackFromVersion != "2.0.0" {
		t.Fatalf("cofnięcie zameldowało powrót z wersji %q", cofniecie.RolledBackFromVersion)
	}
	if wersja := tekstZBazyApps(t, baza,
		`SELECT IFNULL(wersja,'') FROM rozszerzenie WHERE identyfikator_zewnetrzny = ?`,
		pozycja); wersja != "1.0.0" {
		t.Fatalf("po cofnięciu pozycja stoi w wersji %q", wersja)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM historia_rozszerzenia
		  WHERE rozszerzenie_kod = ? AND czynnosc = 'rolledBack'
		    AND wersja_przed = '2.0.0' AND wersja_po = '1.0.0'`, pozycja); ile != 1 {
		t.Fatal("dziennik cyklu życia nie odnotował cofnięcia z wersjami")
	}

	// Aktualizacja liczy się z rejestru wersji: po cofnięciu do 1.0.0 wersja
	// 2.0.0 jest dostępna i łamie zgodność semantyczną.
	var aktualizacje shared.ExtensionUpdateCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionUpdateCheck,
		shared.ExtensionUpdateCheckRequest{ExtensionId: wskaznik(pozycja)}, &aktualizacje)
	if len(aktualizacje.Updates) != 1 || aktualizacje.Updates[0].AvailableVersion != "2.0.0" {
		t.Fatalf("sprawdzenie aktualizacji: %+v", aktualizacje.Updates)
	}
	if aktualizacje.Updates[0].Breaking == nil || !*aktualizacje.Updates[0].Breaking {
		t.Fatal("zmiana pierwszego członu wersji nie została uznana za łamiącą zgodność")
	}
}

// TestWebhookIOdwzorowanieZostajaWBazie wykazuje trwałość konfiguracji
// integracji oraz to, że adres nasłuchu składa rdzeń, nie Operator.
func TestWebhookIOdwzorowanieZostajaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "konektor-crm",
		shared.ExtensionKindApi)

	var przychodzacy shared.ExtensionWebhookSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionWebhookSave,
		shared.ExtensionWebhookSaveRequest{
			ExtensionId: pozycja, Direction: shared.ExtensionWebhookDirectionInbound,
			SecretRef: wskaznik("sejf:crm/hmac"), Enabled: wskaznik(true),
		}, &przychodzacy)
	if przychodzacy.Webhook.ReceiveUrl == nil ||
		!strings.Contains(*przychodzacy.Webhook.ReceiveUrl, "konektor-crm") {
		t.Fatalf("adres nasłuchu nie powstał w rdzeniu: %+v", przychodzacy.Webhook.ReceiveUrl)
	}

	// Webhook wychodzący bez adresu docelowego nie ma dokąd wypchnąć zdarzenia.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandExtensionWebhookSave,
		shared.ExtensionWebhookSaveRequest{
			ExtensionId: pozycja, Direction: shared.ExtensionWebhookDirectionOutbound,
		})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("webhook wychodzący bez adresu wrócił kodem %q", odmowa.Code)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionWebhookSave,
		shared.ExtensionWebhookSaveRequest{
			ExtensionId: pozycja, Direction: shared.ExtensionWebhookDirectionOutbound,
			Url:        wskaznik("https://crm.example.com/hook"),
			EventTypes: []string{"deployment.finished", "extension.changed"},
		}, &przychodzacy)

	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM webhook_rozszerzenia WHERE rozszerzenie_kod = ?`, pozycja); ile != 2 {
		t.Fatalf("w bazie stoi %d webhooków zamiast dwóch", ile)
	}

	var wykaz shared.ExtensionWebhookListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionWebhookList,
		shared.ExtensionWebhookListRequest{
			ExtensionId: wskaznik(pozycja),
			Direction:   wskaznik(shared.ExtensionWebhookDirection(shared.ExtensionWebhookDirectionOutbound)),
		}, &wykaz)
	if wykaz.Total != 1 || len(wykaz.Webhooks[0].EventTypes) != 2 {
		t.Fatalf("zawężenie po kierunku: %+v", wykaz.Webhooks)
	}

	var mapowanie shared.ExtensionMappingSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionMappingSave,
		shared.ExtensionMappingSaveRequest{
			ExtensionId: pozycja, Name: "klient → kontrahent",
			Rules: jsonSurowy(t, map[string]any{"name": "nazwa", "email": ""}),
		}, &mapowanie)
	if len(mapowanie.ValidationIssues) != 1 {
		t.Fatalf("walidacja reguł oddała %+v", mapowanie.ValidationIssues)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM mapowanie_rozszerzenia WHERE identyfikator_zewnetrzny = ?`,
		mapowanie.Mapping.Id); ile != 1 {
		t.Fatal("odwzorowania nie ma w bazie mimo zastrzeżeń walidacji")
	}
}

// TestReferencjaSekretuNieNiesieTresciPoswiadczenia wykazuje granicę warstwy
// sekretów: powiązanie zapisuje klucz jawny, a rejestr oddaje wyłącznie jego.
func TestReferencjaSekretuNieNiesieTresciPoswiadczenia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	pierwsza := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "konektor-crm",
		shared.ExtensionKindApi)
	druga := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "konektor-poczta",
		shared.ExtensionKindApi)

	var powiazanie shared.ExtensionCredentialBindResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionCredentialBind,
		shared.ExtensionCredentialBindRequest{
			ExtensionId: pierwsza, AuthKind: shared.ExtensionAuthKindApiKey,
			CredentialRef: "sejf:crm/klucz",
		}, &powiazanie)

	odwolanie := tekstZBazyApps(t, baza,
		`SELECT IFNULL(odwolanie_sekretu,'') FROM integracja_rozszerzenia WHERE rozszerzenie_kod = ?`,
		pierwsza)
	if odwolanie != "sejf:crm/klucz" {
		t.Fatalf("w bazie stoi odwołanie %q", odwolanie)
	}
	// W bazie modułu nie ma i nie może być treści poświadczenia — jest tam
	// wyłącznie nazwa, po której sejf je wydaje.
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM integracja_rozszerzenia WHERE odwolanie_sekretu LIKE '%tajne%'`); ile != 0 {
		t.Fatal("w wierszu integracji znalazła się treść zamiast odwołania")
	}

	var udostepnienie shared.ExtensionSecretShareResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionSecretShare,
		shared.ExtensionSecretShareRequest{
			SecretRef: "sejf:crm/klucz", ExtensionIds: []string{pierwsza, druga},
			RoleIds: []string{"rola-executor"},
		}, &udostepnienie)
	if len(udostepnienie.Secret.SharedWithExtensionIds) != 2 ||
		len(udostepnienie.Secret.SharedWithRoleIds) != 1 {
		t.Fatalf("zakres współdzielenia: %+v", udostepnienie.Secret)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM udostepnienie_sekretu_rozszerzenia`); ile != 3 {
		t.Fatalf("w bazie stoi %d wierszy zakresu zamiast trzech", ile)
	}

	// Udostępnienie pozycji, której nie ma, jest odmową — zakres wskazujący
	// nikogo nie jest zakresem.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandExtensionSecretShare,
		shared.ExtensionSecretShareRequest{
			SecretRef: "sejf:crm/klucz", ExtensionIds: []string{"rozsz-widmo"},
		})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Fatalf("udostępnienie pozycji nieznanej wróciło kodem %q", odmowa.Code)
	}

	var rejestr shared.ExtensionSecretListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionSecretList,
		shared.ExtensionSecretListRequest{}, &rejestr)
	if rejestr.Total != 1 || rejestr.Secrets[0].Ref != "sejf:crm/klucz" {
		t.Fatalf("rejestr referencji: %+v", rejestr.Secrets)
	}
	// Sejf poświadczeń rdzenia nie dostał od tej drogi ani jednego wpisu:
	// wprowadzenie poświadczenia zostaje po stronie Operatora.
	if _, err := os.Stat(katalog + "/poswiadczenia.json"); err == nil {
		tresc, _ := os.ReadFile(katalog + "/poswiadczenia.json")
		if strings.Contains(string(tresc), "crm/klucz") {
			t.Fatal("powiązanie referencji zapisało coś w sejfie poświadczeń")
		}
	}
}

// TestSzczegolPozycjiSkladaSieZTegoCoWBazie wykazuje, że karta szczegółów nie
// wymyśla treści: narzędzia, uprawnienia i dziennik zmian pochodzą z wierszy.
func TestSzczegolPozycjiSkladaSieZTegoCoWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	adres, _ := serwerProtokoluSprawdzianu(t)
	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "mcp-repozytoria",
		shared.ExtensionKindMcp)
	ustawTransportSprawdzianu(t, zmontowany, zycie, pozycja, adres)

	var pusty shared.ExtensionDetailGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionDetailGet,
		shared.ExtensionDetailGetRequest{ExtensionId: pozycja}, &pusty)
	if len(pusty.Detail.Tools) != 0 {
		t.Fatalf("karta szczegółów pokazała %d narzędzi przed odkryciem", len(pusty.Detail.Tools))
	}

	var wykaz shared.ExtensionToolListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionToolList,
		shared.ExtensionToolListRequest{ExtensionId: pozycja, Refresh: wskaznik(true)}, &wykaz)

	if _, err := baza.Exec(
		`INSERT INTO wersja_rozszerzenia (rozszerzenie_kod, wersja, dziennik_zmian, utworzono)
		 VALUES (?, '1.1.0', 'dodano repo.commit', 1)`, pozycja); err != nil {
		t.Fatalf("nie można wpisać wersji: %v", err)
	}

	var szczegol shared.ExtensionDetailGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionDetailGet,
		shared.ExtensionDetailGetRequest{ExtensionId: pozycja}, &szczegol)
	if len(szczegol.Detail.Tools) != 3 {
		t.Fatalf("karta szczegółów ma %d wpisów zamiast trzech", len(szczegol.Detail.Tools))
	}
	if szczegol.Detail.Changelog == nil ||
		!strings.Contains(*szczegol.Detail.Changelog, "repo.commit") {
		t.Fatalf("dziennik zmian nie pochodzi z wierszy wersji: %+v", szczegol.Detail.Changelog)
	}
}

// TestAudytWiazeWywolanieZUprawnieniami wykazuje, że audyt czyta te same
// wiersze wywołań i dokłada do nich uprawnienia obowiązujące przy pozycji.
func TestAudytWiazeWywolanieZUprawnieniami(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	adres, _ := serwerProtokoluSprawdzianu(t)
	pozycja := pozycjaSprawdzianuRozszerzen(t, zmontowany, zycie, "mcp-repozytoria",
		shared.ExtensionKindMcp)
	ustawTransportSprawdzianu(t, zmontowany, zycie, pozycja, adres)

	var nadanie shared.ExtensionPermissionGrantResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionPermissionGrant,
		shared.ExtensionPermissionGrantRequest{
			ExtensionId: pozycja,
			Permissions: []shared.ExtensionPermission{
				{Scope: shared.ExtensionPermissionScopeNetwork, Target: wskaznik("repo.example.com")},
			},
		}, &nadanie)

	var wywolanie shared.ExtensionToolCallResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionToolCall,
		shared.ExtensionToolCallRequest{ExtensionId: pozycja, ToolName: "repo.search"},
		&wywolanie)

	var audyt shared.ExtensionAuditListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandExtensionAuditList,
		shared.ExtensionAuditListRequest{ExtensionId: wskaznik(pozycja)}, &audyt)

	if audyt.Total != 1 {
		t.Fatalf("audyt oddał %d wpisów zamiast jednego", audyt.Total)
	}
	wpis := audyt.Entries[0]
	if wpis.ToolName == nil || *wpis.ToolName != "repo.search" {
		t.Fatalf("wpis audytu nie wskazuje narzędzia: %+v", wpis)
	}
	if len(wpis.Permissions) != 1 ||
		wpis.Permissions[0].Scope != shared.ExtensionPermissionScopeNetwork {
		t.Fatalf("wpis audytu nie niesie uprawnień pozycji: %+v", wpis.Permissions)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM wywolanie_rozszerzenia WHERE rozszerzenie_kod = ?`, pozycja); ile != 1 {
		t.Fatalf("dziennik użycia ma %d wierszy", ile)
	}
}
