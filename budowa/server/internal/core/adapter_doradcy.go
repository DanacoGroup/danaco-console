// Plik niesie adapter doradcy: wykonanie konsultacji tą samą drogą, którą rdzeń woła każdy inny model. O tym, który model radzi, rozstrzyga sufit siły, a nie ten plik.
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

// Adapter wypełnia port rodziny advisor.* w całości, wraz z komendą advisor.consult obsługiwaną w handlers_doradcy.go.
var _ Doradcy = (*adapterDoradcow)(nil)

// przedrostekKonsultacji znakuje tożsamość zewnętrzną każdego wpisu konsultacji w dzienniku doradców rdzenia.
const przedrostekKonsultacji = "kons-"

// adapterDoradcow wykonuje konsultacje. Dziennik bywa pusty — konsultacja
// odbywa się wtedy tak samo i zostaje jawna w strumieniu, traci wyłącznie ślad
// w bazie.
type adapterDoradcow struct {
	kanaly   *models.Rejestr
	dziennik podagenci.Dziennik
	// okna jest rejestrem okien sesji: kanał pytającego bierze się z okna, nigdy z żądania.
	okna *session.Rejestr
	// nadajnik niesie radę do okna, nie tylko do modelu; bez niego jawność ginie w niewidzianym buforze.
	nadajnik Nadajnik
}

// nowyAdapterDoradcow wiąże adapter z rejestrem kanałów modelu, gotowym do wykonania pierwszej konsultacji.
func nowyAdapterDoradcow(kanaly *models.Rejestr) *adapterDoradcow {
	return &adapterDoradcow{kanaly: kanaly}
}

// ZDziennikiem wpina trwałość konsultacji. Zwraca adapter, żeby montaż wiązał
// zależność w łańcuchu — wzorem `ZKanalami` pozostałych modułów.
func (a *adapterDoradcow) ZDziennikiem(d podagenci.Dziennik) *adapterDoradcow {
	a.dziennik = d
	return a
}

// ZOknami wpina rejestr okien sesji — źródło kanału pytającego. Rejestr niewpięty nie gasi adaptera: `Skonsultuj` woła się nadal tą samą drogą, natomiast komenda `advisor.consult` odmawia.
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

// Skonsultuj wykonuje jedną konsultację: wybiera doradcę, wysyła pytanie jego kanałem, przepuszcza strumień doradcy do ujścia wołającego, dokłada blok jawności i zapisuje ślad z prowenancją.
func (a *adapterDoradcow) Skonsultuj(ctx context.Context, pytanie podagenci.Pytanie,
	ujscie models.Ujscie) (podagenci.Rada, error) {

	if err := pytanie.Sprawdz(); err != nil {
		return podagenci.Rada{}, bladPytaniaDoradcy(err.Error())
	}
	if a.kanaly == nil {
		return podagenci.Rada{}, bladBrakuKanalowDoradcy()
	}
	// Sufit siły stoi w pakiecie pojęciowym; adapter go wykonuje, odmowa idzie wołającemu i do dziennika.
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
	// Blok jawności idzie po odpowiedzi doradcy: widać radę, a potem podpis o cudzej, niewiążącej treści.
	_ = models.NadajTekst(ctx, ujscie, zapytanie, rada.Jawnie(pytanie))
	return rada, a.zapiszRade(ctx, pytanie, rada, zbierak.prowenancja)
}

// zapiszRade utrwala odbytą konsultację razem z jej prowenancją wywołania w dzienniku doradców rdzenia.
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

// zbierakRadyDoradcy przepuszcza strumień doradcy do ujścia wołającego i po drodze zapamiętuje tekst rady oraz opis wywołania, nie zatrzymując fragmentów u siebie.
type zbierakRadyDoradcy struct {
	dalej       models.Ujscie
	tekst       strings.Builder
	prowenancja string
}

// Fragment odbiera jeden fragment strumienia doradcy i przepuszcza go dalej, wprost ku ujściu wywołania.
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

// bladBrakuKanalowDoradcy — konsultacja ma iść modelem, a rejestru kanałów nie wpięto do tego adaptera.
func bladBrakuKanalowDoradcy() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"doradca: rejestr kanałów modelu nie jest wpięty — konsultacja nie ma czym wołać doradcy"))
}

// bladDoboruDoradcy — rejestr jest wpięty, lecz doboru nie da się wykonać; kod odmowy mówi, czy ponawiać, więc trwała przyczyna idzie jako validation_failed, nie channel_unavailable.
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

// bladPytaniaDoradcy — pytania nie da się zadać, bo treść albo kanał doradcy są tu bezużyteczne dla wywołania.
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
