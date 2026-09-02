// Plik niesie adapter doradcy: konsultacja idzie tą samą drogą, którą rdzeń
// woła każdy inny model; o doborze doradcy rozstrzyga sufit siły pakietu podagenci.
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

var _ Doradcy = (*adapterDoradcow)(nil)

const przedrostekKonsultacji = "kons-"

type adapterDoradcow struct {
	kanaly   *models.Rejestr
	dziennik podagenci.Dziennik
	okna     *session.Rejestr
	nadajnik Nadajnik
}

func nowyAdapterDoradcow(kanaly *models.Rejestr) *adapterDoradcow {
	return &adapterDoradcow{kanaly: kanaly}
}

func (a *adapterDoradcow) ZDziennikiem(d podagenci.Dziennik) *adapterDoradcow {
	a.dziennik = d
	return a
}

func (a *adapterDoradcow) ZOknami(r *session.Rejestr) *adapterDoradcow {
	a.okna = r
	return a
}

func (a *adapterDoradcow) ZNadajnikiem(n Nadajnik) *adapterDoradcow {
	a.nadajnik = n
	return a
}

func (a *adapterDoradcow) Doradcy() []podagenci.Kandydat {
	if a.kanaly == nil {
		return nil
	}
	return podagenci.Kandydaci(a.kanaly.Wykaz())
}

func (a *adapterDoradcow) Skonsultuj(ctx context.Context, pytanie podagenci.Pytanie,
	ujscie models.Ujscie) (podagenci.Rada, error) {

	if err := pytanie.Sprawdz(); err != nil {
		return podagenci.Rada{}, bladPytaniaDoradcy(err.Error())
	}
	if a.kanaly == nil {
		return podagenci.Rada{}, bladBrakuKanalowDoradcy()
	}
	// Sufit siły rozstrzyga pakiet podagenci; odmowa doboru idzie wołającemu i do dziennika.
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
	// Blok jawności idzie po odpowiedzi doradcy: podpis o cudzej, niewiążącej treści.
	_ = models.NadajTekst(ctx, ujscie, zapytanie, rada.Jawnie(pytanie))
	return rada, a.zapiszRade(ctx, pytanie, rada, zbierak.prowenancja)
}

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

// Niepowodzenie zapisu śladu nie przykrywa powodu odmowy, więc funkcja niczego nie zwraca.
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

// Kolumna pytającego w dzienniku konsultacji jest wymagana; kanał nieopisany dostaje znak „-".
func kanalPytajacego(p podagenci.Pytanie) string {
	if kanal := strings.TrimSpace(p.Pytajacy); kanal != "" {
		return kanal
	}
	return "-"
}

// Warunek CHECK tabeli dziennika konsultacji odrzuca radę z pustą treścią.
func radaDoZapisu(tresc string) string {
	if strings.TrimSpace(tresc) == "" {
		return "(doradca odpowiedział pustą treścią)"
	}
	return tresc
}

type zbierakRadyDoradcy struct {
	dalej       models.Ujscie
	tekst       strings.Builder
	prowenancja string
}

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

func bladOknaDoradcy(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "doradca: "+powod))
}

func bladBrakuKanalowDoradcy() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"doradca: rejestr kanałów modelu nie jest wpięty — konsultacja nie ma czym wołać doradcy"))
}

// Kod odmowy mówi klientowi, czy ponawiać: przyczyna trwała idzie jako validation_failed.
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

func bladPytaniaDoradcy(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"doradca: "+powod))
}

func bladWywolaniaDoradcy(err error) error {
	return protocol.JakoError(models.BladKanalu(err))
}

func bladSladuDoradcy(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}
