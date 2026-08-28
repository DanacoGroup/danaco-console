// Sprawdza skutek modułu Translate: kontrolę wierności przekładu oraz
// tłumaczenie zwrotne, wykonywane przez prawdziwy kanał rdzenia wskazujący
// punkt końcowy modelu podniesiony na czas sprawdzianu w tym samym procesie.
package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// znacznikPoleceniaZwrotnego jest fragmentem polecenia, po którym punkt końcowy
// sprawdzianu poznaje, że pytany jest o tłumaczenie zwrotne, a nie o przekład;
// pochodzi z `polecenieTlumaczeniaZwrotnego`.
const znacznikPoleceniaZwrotnego = "KONTROLA WIERNOŚCI"

// serwerModelu podnosi punkt końcowy modelu tekstowego. Odpowiedź składa
// funkcja dostająca treść polecenia, którą rdzeń wysłał — dzięki temu jeden
// serwer obsługuje i przekład, i tłumaczenie zwrotne, a sprawdzian nie musi
// zgadywać kolejności wywołań.
func serwerModelu(t *testing.T, odpowiedz func(polecenie string) string) *httptest.Server {
	t.Helper()

	serwer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		surowe, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []any{map[string]any{
				"message": map[string]any{"content": odpowiedz(string(surowe))},
			}},
		})
	}))
	t.Cleanup(serwer.Close)
	return serwer
}

// wpiszKanalTekstowy zakłada wiersz rejestru kanałów wskazujący podany adres
// i oddaje identyfikator kanału. Kanał jedzie bez strumienia, bo punkt końcowy
// sprawdzianu oddaje jeden dokument JSON.
func wpiszKanalTekstowy(t *testing.T, zmontowany *Zmontowany, zycie context.Context, adres string) string {
	t.Helper()

	var wynik shared.ChannelAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandChannelAdd, shared.ChannelAddRequest{
		Name:  "kanał tekstowy sprawdzianu",
		Kind:  "api",
		Model: wskaznik("model-sprawdzianu"),
		Config: jsonSurowy(t, map[string]any{
			"base_url": adres,
			"strumien": "false",
		}),
	}, &wynik)

	if wynik.Channel.Id == "" {
		t.Fatal("channel.add nie oddał identyfikatora kanału")
	}
	return wynik.Channel.Id
}

// panelOkna odczytuje panel okna tłumaczenia bez wołania modelu.
// `translate.source.set` oddaje panele okna wraz z ich treścią, a wywołanie go
// tą samą treścią źródłową niczego nie zmienia — jest to więc odczyt, nie zmiana.
func panelOkna(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno, tekstZrodlowy, kodPanelu string) shared.TranslationPanel {
	t.Helper()

	var stan shared.TranslateSourceSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSourceSet,
		shared.TranslateSourceSetRequest{
			WindowId:  okno,
			Text:      tekstZrodlowy,
			Resegment: wskaznik(false),
		}, &stan)

	for _, panel := range stan.Panels {
		if panel.Id == kodPanelu {
			return panel
		}
	}
	t.Fatalf("okno %s nie zna panelu %s", okno, kodPanelu)
	return shared.TranslationPanel{}
}

// Teksty sprawdzianu. Źródło niesie trzy rzeczy, których przekład nie może
// zgubić: liczbę, walutę i znacznik podstawienia. Wierny przekład niesie je
// wszystkie, niewierny żaden; różnica długości obu tekstów nie wpływa na wynik.
const (
	zrodloKontroli  = "Zamówienie na 1500 zł obejmuje {liczbaSztuk} sztuk towaru."
	przekladWierny  = "The order for 1500 zł covers {liczbaSztuk} units of goods."
	przekladZgubion = "The order covers a number of units of goods."
)

// zalozPanelPrzekladu składa okno tłumaczenia z tekstem źródłowym i panel
// docelowy wypełniony przekładem modelu, po czym oddaje kod panelu.
func zalozPanelPrzekladu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno, jezykZrodlowy, przeklad, kanal string) string {
	t.Helper()

	zadanie := shared.TranslateSourceSetRequest{WindowId: okno, Text: zrodloKontroli}
	if jezykZrodlowy != "" {
		zadanie.SourceLanguage = &jezykZrodlowy
	}
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSourceSet, zadanie, nil)

	var dodany shared.TranslateTargetAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateTargetAdd,
		shared.TranslateTargetAddRequest{WindowId: okno, Language: "angielski", ChannelId: &kanal}, &dodany)

	if dodany.Panel.Text == nil || *dodany.Panel.Text != przeklad {
		t.Fatalf("panel dostał treść %v, model oddał %q", dodany.Panel.Text, przeklad)
	}
	return dodany.Panel.Id
}

// TestKontrolaJakosciZglaszaToCzegoPrzekladNieOddal dowodzi, że kontrola ma co
// porównywać: przekład, który zgubił liczbę, walutę i znacznik podstawienia,
// wychodzi z zastrzeżeniem o każdej z tych rzeczy.
func TestKontrolaJakosciZglaszaToCzegoPrzekladNieOddal(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	serwer := serwerModelu(t, func(string) string { return przekladZgubion })
	kanal := wpiszKanalTekstowy(t, zmontowany, zycie, serwer.URL)
	panel := zalozPanelPrzekladu(t, zmontowany, zycie, "okno-zgubione", "polski", przekladZgubion, kanal)

	var kontrola shared.TranslateQualityCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateQualityCheck,
		shared.TranslateQualityCheckRequest{PanelId: panel}, &kontrola)

	if kontrola.PanelId != panel {
		t.Errorf("kontrola opisuje panel %q, pytano o %q", kontrola.PanelId, panel)
	}
	oczekiwane := []shared.TranslationIssueKind{
		shared.TranslationIssueKindNumber,
		shared.TranslationIssueKindCurrency,
		shared.TranslationIssueKindPlaceholder,
	}
	for _, rodzaj := range oczekiwane {
		if !zawieraRodzaj(kontrola.Issues, rodzaj) {
			t.Errorf("kontrola nie zgłosiła rodzaju %q, choć przekład zgubił tę rzecz; wykaz: %v",
				rodzaj, kontrola.Issues)
		}
	}
}

// TestKontrolaJakosciMilczyNadPrzeklademWiernym pilnuje drugiej strony tej samej
// miary. Kontrola, która zgłasza zastrzeżenia zawsze, jest tak samo bezużyteczna
// jak ta, która nie zgłasza ich nigdy — Operator nauczy się ją ignorować.
func TestKontrolaJakosciMilczyNadPrzeklademWiernym(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	serwer := serwerModelu(t, func(string) string { return przekladWierny })
	kanal := wpiszKanalTekstowy(t, zmontowany, zycie, serwer.URL)
	panel := zalozPanelPrzekladu(t, zmontowany, zycie, "okno-wierne", "polski", przekladWierny, kanal)

	var kontrola shared.TranslateQualityCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateQualityCheck,
		shared.TranslateQualityCheckRequest{PanelId: panel}, &kontrola)

	if len(kontrola.Issues) != 0 {
		t.Errorf("kontrola zgłosiła %d zastrzeżeń nad przekładem niosącym liczbę, walutę i znacznik: %v",
			len(kontrola.Issues), kontrola.Issues)
	}
}

// TestTlumaczenieZwrotneNieNadpisujeTresciPanelu sprawdza, że po przebiegu
// istnieją dwa różne teksty: panel w języku docelowym i tłumaczenie zwrotne
// w źródłowym. Panel przepisany wynikiem kontroli byłby porównaniem tekstu
// z samym sobą.
func TestTlumaczenieZwrotneNieNadpisujeTresciPanelu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const zwrotne = "Zamówienie obejmuje pewną liczbę sztuk towaru."
	serwer := serwerModelu(t, func(polecenie string) string {
		if strings.Contains(polecenie, znacznikPoleceniaZwrotnego) {
			return zwrotne
		}
		return przekladZgubion
	})
	kanal := wpiszKanalTekstowy(t, zmontowany, zycie, serwer.URL)
	panel := zalozPanelPrzekladu(t, zmontowany, zycie, "okno-zwrotne", "polski", przekladZgubion, kanal)

	var wynik shared.TranslateBacktranslationRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateBacktranslationRun,
		shared.TranslateBacktranslationRunRequest{PanelId: panel, ChannelId: &kanal}, &wynik)

	if wynik.PanelId != panel {
		t.Errorf("wynik opisuje panel %q, pytano o %q", wynik.PanelId, panel)
	}
	if wynik.Text != zwrotne {
		t.Errorf("tłumaczenie zwrotne niesie %q, model oddał %q", wynik.Text, zwrotne)
	}

	po := panelOkna(t, zmontowany, zycie, "okno-zwrotne", zrodloKontroli, panel)
	if po.Text == nil {
		t.Fatal("panel stracił treść po przebiegu kontroli wierności")
	}
	if *po.Text == wynik.Text {
		t.Fatal("panel został przepisany tłumaczeniem zwrotnym — kontrola porównuje tekst sam ze sobą")
	}
	if *po.Text != przekladZgubion {
		t.Errorf("panel niesie po przebiegu %q, przekładem było %q", *po.Text, przekladZgubion)
	}
}

// TestDrugiPrzebiegKontroliWiernosciOddajeNowyWynik dowodzi, że wynik jest
// mierzony, a nie zapamiętany: przy zmienionej odpowiedzi modelu drugi
// przebieg ma dać co innego niż pierwszy przebieg.
func TestDrugiPrzebiegKontroliWiernosciOddajeNowyWynik(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const pierwszeZwrotne = "Zamówienie obejmuje pewną liczbę sztuk towaru."
	const drugieZwrotne = "Zlecenie dotyczy nieokreślonej liczby sztuk."

	zwrotne := pierwszeZwrotne
	serwer := serwerModelu(t, func(polecenie string) string {
		if strings.Contains(polecenie, znacznikPoleceniaZwrotnego) {
			return zwrotne
		}
		return przekladZgubion
	})
	kanal := wpiszKanalTekstowy(t, zmontowany, zycie, serwer.URL)
	panel := zalozPanelPrzekladu(t, zmontowany, zycie, "okno-dwa-przebiegi", "polski", przekladZgubion, kanal)

	var pierwszy shared.TranslateBacktranslationRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateBacktranslationRun,
		shared.TranslateBacktranslationRunRequest{PanelId: panel, ChannelId: &kanal}, &pierwszy)

	zwrotne = drugieZwrotne
	var drugi shared.TranslateBacktranslationRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateBacktranslationRun,
		shared.TranslateBacktranslationRunRequest{PanelId: panel, ChannelId: &kanal}, &drugi)

	if pierwszy.Text != pierwszeZwrotne {
		t.Errorf("pierwszy przebieg oddał %q, model mówił %q", pierwszy.Text, pierwszeZwrotne)
	}
	if drugi.Text != drugieZwrotne {
		t.Errorf("drugi przebieg oddał %q, model mówił %q", drugi.Text, drugieZwrotne)
	}
	if pierwszy.Text == drugi.Text {
		t.Error("dwa przebiegi nad różnymi odpowiedziami modelu dały ten sam wynik")
	}

	// Panel przeżył oba przebiegi bez zmiany — kontrola nie tyka tego, co bada.
	po := panelOkna(t, zmontowany, zycie, "okno-dwa-przebiegi", zrodloKontroli, panel)
	if po.Text == nil || *po.Text != przekladZgubion {
		t.Errorf("po dwóch przebiegach panel niesie %v, przekładem było %q", po.Text, przekladZgubion)
	}
}

// TestKontrolaWiernosciBezJezykaZrodlowegoOdmawiaINieTykaPanelu sprawdza, że
// bez języka źródłowego okna komenda odmawia, nazywa brak i zostawia panel
// nietknięty, zamiast zgadywać język przekładu.
func TestKontrolaWiernosciBezJezykaZrodlowegoOdmawiaINieTykaPanelu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	serwer := serwerModelu(t, func(polecenie string) string {
		if strings.Contains(polecenie, znacznikPoleceniaZwrotnego) {
			t.Error("rdzeń wołał model o tłumaczenie zwrotne, choć okno nie ma języka źródłowego")
		}
		return przekladWierny
	})
	kanal := wpiszKanalTekstowy(t, zmontowany, zycie, serwer.URL)
	// Okno bez języka źródłowego — `source.set` niczego nie rozpoznaje,
	// kolumna zostaje pusta.
	panel := zalozPanelPrzekladu(t, zmontowany, zycie, "okno-bez-jezyka", "", przekladWierny, kanal)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandTranslateBacktranslationRun,
		shared.TranslateBacktranslationRunRequest{PanelId: panel, ChannelId: &kanal})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Errorf("odmowa niesie kod %q, brak języka źródłowego to %q",
			odmowa.Code, shared.ErrorCodeValidationFailed)
	}

	po := panelOkna(t, zmontowany, zycie, "okno-bez-jezyka", zrodloKontroli, panel)
	if po.Text == nil || *po.Text != przekladWierny {
		t.Errorf("po odmowie panel niesie %v, przekładem było %q", po.Text, przekladWierny)
	}
}

// TestPustyPrzekladModeluNieZakladaPanelu sprawdza wejście: model, który
// oddał pustkę, nie przetłumaczył niczego, więc komenda ma odmówić, a okno
// ma zostać bez panelu.
func TestPustyPrzekladModeluNieZakladaPanelu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	serwer := serwerModelu(t, func(string) string { return "   " })
	kanal := wpiszKanalTekstowy(t, zmontowany, zycie, serwer.URL)

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSourceSet,
		shared.TranslateSourceSetRequest{
			WindowId: "okno-pustego-przekladu", Text: zrodloKontroli,
			SourceLanguage: wskaznik("polski"),
		}, nil)

	wykonajOdmowna(t, zmontowany, zycie, shared.CommandTranslateTargetAdd,
		shared.TranslateTargetAddRequest{
			WindowId: "okno-pustego-przekladu", Language: "angielski", ChannelId: &kanal,
		})

	var stan shared.TranslateSourceSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSourceSet,
		shared.TranslateSourceSetRequest{
			WindowId: "okno-pustego-przekladu", Text: zrodloKontroli, Resegment: wskaznik(false),
		}, &stan)
	if len(stan.Panels) != 0 {
		t.Errorf("po odmowie okno ma %d paneli — panel bez przekładu został założony", len(stan.Panels))
	}
}

// zawieraRodzaj mówi, czy wykaz zastrzeżeń kontroli niesie zastrzeżenie danego
// rodzaju. `quality.check` oddaje zdania, nie byty, a zdanie zaczyna się od
// nazwy rodzaju (`zdaniaNiezgodnosci`).
func zawieraRodzaj(zastrzezenia []string, rodzaj shared.TranslationIssueKind) bool {
	for _, zdanie := range zastrzezenia {
		if strings.HasPrefix(zdanie, string(rodzaj)) {
			return true
		}
	}
	return false
}
