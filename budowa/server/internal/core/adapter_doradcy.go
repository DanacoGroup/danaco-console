// Odpowiedzialność pliku: adapter doradcy — wykonanie konsultacji po tej samej
// drodze, którą rdzeń woła każdy inny model.
//
// O tym, który model radzi, rozstrzyga sufit siły, a nie ten plik: model
// z własnej inicjatywy konsultuje wyłącznie model równy sobie albo słabszy.
// Reguła i jej powód stoją w `podagenci/doradca_wybor.go` — adapter jej nie
// powtarza, żeby nie było dwóch miejsc, w których wolno ją poluzować. Komenda
// `advisor.consult`, jedyny wołacz konsultacji, stoi
// w `adapter_doradcy_konsultacja.go`.
//
// Adapter nie zna SQL-a ani protokołu dostawcy. Bierze rejestr kanałów modelu
// i dziennik konsultacji pakietu `podagenci`, a pojęcie doradcy — kto nim może
// być, jak brzmi pytanie, czym jest rada — zostaje po stronie pakietu
// pojęciowego.
//
// Trzy własności pojęcia widać tu jako trzy czynności:
//
//	Jawność. Wszystkie fragmenty strumienia doradcy idą do ujścia wołającego
//	bez zmiany, a na koniec dokładany jest blok `Rada.Jawnie` — Operator widzi
//	pytanie, doradcę i radę w tym samym oknie, w którym pracuje agent.
//	Adapter nie oddaje samej treści rady bez wskazania doradcy.
//
//	Prowenancja. Opis wywołania składa kanał, tak jak przy każdym innym
//	zapytaniu (`models.NadajProwenancje`). Adapter go nie wytwarza — przejmuje
//	fragment `provenance` przelatujący strumieniem i zapisuje ten sam napis
//	w dzienniku. Drugiej prowenancji nie ma.
//
//	Konsultacja, nie delegacja. Zwrócona `Rada` nie zmienia niczego w stanie
//	rdzenia: nie zakłada pozycji kolejki, nie startuje tury, nie zapisuje
//	wiadomości agenta. Wołający dostaje radę i sam rozstrzyga, co z nią zrobi.
//
// Odmowa też trafia do dziennika. Konsultacja, która się nie odbyła, jest
// zdarzeniem, o które Operator zapyta jako pierwsze — wpis o stanie `odmowa`
// niesie powód.
package core

import (
	"context"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/podagenci"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Adapter wypełnia port rodziny `advisor.*` w całości (handlers_doradcy.go).
var _ Doradcy = (*adapterDoradcow)(nil)

// przedrostekKonsultacji znakuje tożsamość zewnętrzną wpisu konsultacji.
const przedrostekKonsultacji = "kons-"

// adapterDoradcow wykonuje konsultacje. Dziennik bywa pusty — konsultacja
// odbywa się wtedy tak samo i zostaje jawna w strumieniu, traci wyłącznie ślad
// w bazie.
type adapterDoradcow struct {
	kanaly   *models.Rejestr
	dziennik podagenci.Dziennik
	// okna jest rejestrem okien sesji i stoi tu po jedno: kanał pytającego
	// bierze się z okna, nigdy z żądania. Bez tego rejestru komenda musiałaby
	// wierzyć modelowi na słowo, kim jest — a wtedy sufit siły dałoby się
	// obejść jednym polem żądania.
	okna *session.Rejestr
	// nadajnik jest drogą, którą rada trafia do okna, a nie tylko do modelu.
	// Bez niego jawność konsultacji kończy się na buforze, którego nikt nie
	// czyta — patrz `adapter_doradcy_konsultacja.go`.
	nadajnik Nadajnik
}

// nowyAdapterDoradcow wiąże adapter z rejestrem kanałów modelu.
func nowyAdapterDoradcow(kanaly *models.Rejestr) *adapterDoradcow {
	return &adapterDoradcow{kanaly: kanaly}
}

// ZDziennikiem wpina trwałość konsultacji. Zwraca adapter, żeby montaż wiązał
// zależność w łańcuchu — wzorem `ZKanalami` pozostałych modułów.
func (a *adapterDoradcow) ZDziennikiem(d podagenci.Dziennik) *adapterDoradcow {
	a.dziennik = d
	return a
}

// ZOknami wpina rejestr okien sesji — źródło kanału pytającego. Rejestr niewpięty nie gasi adaptera: `Skonsultuj` woła się nadal
// tą samą drogą, natomiast komenda `advisor.consult` odmawia, bo okna nie ma
// z czego odczytać.
func (a *adapterDoradcow) ZOknami(r *session.Rejestr) *adapterDoradcow {
	a.okna = r
	return a
}

// ZNadajnikiem wpina szynę zdarzeń: fragmenty konsultacji jadą nią do okna
// kopertami `stream.chunk`. Nadajnik niewpięty nie gasi konsultacji —
// gasi wyłącznie jej widoczność w oknie.
func (a *adapterDoradcow) ZNadajnikiem(n Nadajnik) *adapterDoradcow {
	a.nadajnik = n
	return a
}

// Doradcy oddaje kanały, których wolno zapytać — do pokazania Operatorowi
// i do odpowiedzi na pytanie „kto tu może radzić".
func (a *adapterDoradcow) Doradcy() []podagenci.Kandydat {
	if a.kanaly == nil {
		return nil
	}
	return podagenci.Kandydaci(a.kanaly.Wykaz())
}

// Skonsultuj wykonuje jedną konsultację: wybiera doradcę, wysyła pytanie jego
// kanałem, przepuszcza cały strumień doradcy do ujścia wołającego, dokłada blok
// jawności i zapisuje ślad razem z prowenancją wywołania.
//
// Ujście pochodzi od wołającego i jest zwykle tym samym ujściem, którym płynie
// odpowiedź agenta — dlatego rada widoczna jest tam, gdzie pracuje agent, a nie
// w osobnym, cichym kanale.
func (a *adapterDoradcow) Skonsultuj(ctx context.Context, pytanie podagenci.Pytanie,
	ujscie models.Ujscie) (podagenci.Rada, error) {

	if err := pytanie.Sprawdz(); err != nil {
		return podagenci.Rada{}, bladPytaniaDoradcy(err.Error())
	}
	if a.kanaly == nil {
		return podagenci.Rada{}, bladBrakuKanalowDoradcy()
	}
	// Sufit siły stoi w pakiecie pojęciowym (`podagenci/doradca_wybor.go`),
	// a nie tutaj — adapter go wykonuje, nie ma własnej wersji. Odmowa
	// wraca wołającemu nazwana i jednocześnie idzie do dziennika jako wpis
	// `odmowa`, bo pytanie „dlaczego rdzeń nie zapytał mocniejszego modelu"
	// pada po fakcie i musi mieć odpowiedź w bazie.
	doradca, err := pytanie.Doradca(a.kanaly.Wykaz())
	if err != nil {
		a.zapiszOdmowe(ctx, pytanie, podagenci.Kandydat{Kanal: "-", Model: "-"}, err.Error())
		return podagenci.Rada{}, bladDoboruDoradcy(err)
	}

	zapytanie := models.Zapytanie{
		Zasiegi:   models.Zasiegi{Okno: pytanie.Okno},
		Wiadomosc: pytanie.Wiadomosc,
		Tresc:     pytanie.Tresc(),
		Kanal:     doradca.Kanal,
	}
	zbierak := &zbierakRadyDoradcy{dalej: ujscie}
	if err := a.kanaly.Wyslij(ctx, zapytanie, zbierak); err != nil {
		a.zapiszOdmowe(ctx, pytanie, doradca, err.Error())
		return podagenci.Rada{}, bladWywolaniaDoradcy(err)
	}

	rada := podagenci.ZlozRade(doradca, pytanie, zbierak.tekst.String())
	// Blok jawności idzie strumieniem po odpowiedzi doradcy: Operator widzi
	// najpierw to, co doradca powiedział, a potem podpis mówiący, że to była
	// rada cudza i niewiążąca.
	_ = models.NadajTekst(ctx, ujscie, zapytanie, rada.Jawnie(pytanie))
	return rada, a.zapiszRade(ctx, pytanie, rada, zbierak.prowenancja)
}

// zapiszRade utrwala odbytą konsultację razem z prowenancją wywołania.
func (a *adapterDoradcow) zapiszRade(ctx context.Context, pytanie podagenci.Pytanie,
	rada podagenci.Rada, prowenancja string) error {

	if a.dziennik == nil {
		return nil
	}
	_, err := a.dziennik.Zapisz(ctx, podagenci.Konsultacja{
		Identyfikator: nowyIdentyfikator(przedrostekKonsultacji),
		OknoId:        pytanie.Okno,
		Pytajacy:      kanalPytajacego(pytanie),
		DoradcaKanal:  rada.Doradca.Kanal,
		DoradcaModel:  rada.Doradca.Model,
		Pytanie:       pytanie.Tresc(),
		Rada:          radaDoZapisu(rada.Tresc),
		Skrot:         rada.Skrot,
		Prowenancja:   prowenancja,
		Stan:          podagenci.StanRady,
		Utworzono:     rada.Chwila.UnixMilli(),
	})
	if err != nil {
		return bladSladuDoradcy(err)
	}
	return nil
}

// zapiszOdmowe utrwala konsultację, która się nie odbyła. Niepowodzenie zapisu
// śladu nie przykrywa powodu odmowy — wołający i tak dostaje odmowę pierwotną,
// więc ta funkcja niczego nie zwraca.
func (a *adapterDoradcow) zapiszOdmowe(ctx context.Context, pytanie podagenci.Pytanie,
	doradca podagenci.Kandydat, powod string) {

	if a.dziennik == nil {
		return
	}
	_, _ = a.dziennik.Zapisz(ctx, podagenci.Konsultacja{
		Identyfikator: nowyIdentyfikator(przedrostekKonsultacji),
		OknoId:        pytanie.Okno,
		Pytajacy:      kanalPytajacego(pytanie),
		DoradcaKanal:  doradca.Kanal,
		DoradcaModel:  doradca.Model,
		Pytanie:       pytanie.Tresc(),
		Stan:          podagenci.StanOdmowyRady,
		Powod:         powod,
		Utworzono:     time.Now().UTC().UnixMilli(),
	})
}

// kanalPytajacego nazywa kanał agenta zadającego pytanie. Kolumna jest
// wymagana, a pustka nie jest tu prawdą o zdarzeniu — pytanie zawsze ktoś
// zadał, choćby jego kanał nie był opisany w rejestrze.
func kanalPytajacego(p podagenci.Pytanie) string {
	if kanal := strings.TrimSpace(p.Pytajacy); kanal != "" {
		return kanal
	}
	return "-"
}

// radaDoZapisu zabezpiecza więz stanu z treścią: doradca, który odpowiedział
// samą ciszą, dałby wpis „rada bez rady" odrzucany przez warunek CHECK tabeli
// dziennika konsultacji.
func radaDoZapisu(tresc string) string {
	if strings.TrimSpace(tresc) == "" {
		return "(doradca odpowiedział pustą treścią)"
	}
	return tresc
}

// zbierakRadyDoradcy przepuszcza strumień doradcy do ujścia wołającego
// i po drodze zapamiętuje dwie rzeczy: tekst rady oraz opis wywołania.
//
// Przepuszcza wszystkie fragmenty, nie tylko tekst. Gdyby zatrzymywał je
// u siebie, konsultacja byłaby niewidoczna do chwili jej zakończenia, a rada
// pokazana dopiero jako gotowy napis — czyli tak samo jak odpowiedź własna
// agenta, wbrew jawności konsultacji.
type zbierakRadyDoradcy struct {
	dalej       models.Ujscie
	tekst       strings.Builder
	prowenancja string
}

// Fragment odbiera jeden fragment strumienia doradcy.
func (z *zbierakRadyDoradcy) Fragment(ctx context.Context, f models.Fragment) error {
	switch f.Kind {
	case shared.ChunkKindText:
		z.tekst.WriteString(models.TrescFragmentu(f))
	case shared.ChunkKindProvenance:
		if len(f.Data) > 0 {
			z.prowenancja = string(f.Data)
		}
	}
	if z.dalej == nil {
		return nil
	}
	return z.dalej.Fragment(ctx, f)
}

// bladOknaDoradcy — konsultacji nie da się przypisać oknu, więc nie da się
// ustalić kanału pytającego. Odmowa jedzie z powodem: cicha konsultacja
// „od nikogo" znosiłaby sufit siły.
func bladOknaDoradcy(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "doradca: "+powod))
}

// bladBrakuKanalowDoradcy — konsultacja ma iść modelem, a rejestru kanałów nie
// wpięto.
func bladBrakuKanalowDoradcy() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"doradca: rejestr kanałów modelu nie jest wpięty — konsultacja nie ma czym wołać doradcy"))
}

// bladDoboruDoradcy — rejestr jest wpięty, lecz doboru nie da się wykonać.
// Odmowa jedzie wołającemu z powodem od doboru, bo „nie skonsultowano" bez
// zdania dlaczego jest ciszą tam, gdzie stała decyzja; podmiana
// odciętego doradcy na innego byłaby tą samą ciszą, tylko z radą w tle.
//
// Kod odmowy mówi, czy ponawiać. `channel_unavailable` jest w kontrakcie
// kodem ponawialnym (`shared.KodyPonawialne`), więc odmowa trwała nie może nim
// jechać — model dostałby polecenie ponowienia rozstrzygnięcia, które się nie
// zmieni. Stąd podział:
//
//	kanału nie ma / wygaszony       → not_found (jak nieistniejące okno)
//	sufit, dopuszczenie, brak siły  → validation_failed — odmowa trwała;
//	                                  zmienia ją Operator wpisem do rejestru,
//	                                  nie ponowienie tego samego żądania
//	pozostałe                       → channel_unavailable
func bladDoboruDoradcy(powod error) error {
	var kod protocol.KodBledu = shared.ErrorCodeChannelUnavailable
	switch {
	case errors.Is(powod, podagenci.ErrDoradcaNieznany):
		kod = shared.ErrorCodeNotFound
	case errors.Is(powod, podagenci.ErrEskalacjaBezWskazania),
		errors.Is(powod, podagenci.ErrDoradcaNiedopuszczony),
		errors.Is(powod, podagenci.ErrSilaNieopisana):
		kod = shared.ErrorCodeValidationFailed
	}
	return protocol.JakoError(protocol.NowyBlad(kod, "doradca: "+powod.Error()))
}

// bladPytaniaDoradcy — pytania nie da się zadać.
func bladPytaniaDoradcy(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"doradca: "+powod))
}

// bladWywolaniaDoradcy — kanał doradcy odmówił albo zawiódł. Kod bierze
// z katalogu kontraktu przez `models.BladKanalu`.
func bladWywolaniaDoradcy(err error) error {
	return protocol.JakoError(models.BladKanalu(err))
}

// bladSladuDoradcy — rada padła i została pokazana, ale ślad się nie zapisał.
// Wołający dowiaduje się o tym wprost: konsultacja bez śladu w dzienniku jest
// niepełna, a milczenie o niezapisanym śladzie byłoby orzeczeniem z ciszy.
func bladSladuDoradcy(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}
