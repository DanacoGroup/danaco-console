// Plik niesie pomiar miar reguł alertu i ewaluację reguł na tych pomiarach; definicje, rejestr i przekłady na kontrakt leżą w adapter_alerty.go.
package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// kluczPulapuKosztu jest nastawą, wobec której liczy się miara budżetu. Ta sama
// nastawa wstrzymuje turę w warstwie kanału (`injection/pulap.go`) — dwa
// odczyty jednej wartości, nie dwie wartości.
const kluczPulapuKosztu = "pulap_kosztu_usd"

const wymiarZuzyciaAlertow = "channel"

func (a *adapterAlertow) sprawdzMiareReguly(miara shared.AlertMetric) error {
	switch miara {
	case shared.AlertMetricErrorCount:
		if a.diagnostyka == nil {
			return bladWskazaniaAlertu("miara errorCount liczy się z dziennika błędów serwera, " +
				"a serwer nie ma go wpiętego — naprawa: podpiąć repozytorium diagnostyki " +
				"przy składaniu serwera")
		}
		return nil
	case shared.AlertMetricErrorRate, shared.AlertMetricCost, shared.AlertMetricTokens,
		shared.AlertMetricCallLatency:
		if a.prowenancja == nil {
			return bladWskazaniaAlertu("miara " + string(miara) + " liczy się ze śladu wywołań " +
				"modelu, a serwer nie ma go wpiętego — naprawa: podpiąć repozytorium " +
				"prowenancji przy składaniu serwera")
		}
		return nil
	case shared.AlertMetricBudgetPercent:
		if a.prowenancja == nil || a.rozstrzygacz == nil {
			return bladWskazaniaAlertu("miara budgetPercent liczy koszt wobec nastawy " +
				kluczPulapuKosztu + ", więc potrzebuje śladu wywołań i rozstrzygacza " +
				"konfiguracji — serwer nie ma wpiętego przynajmniej jednego z nich")
		}
		return nil
	case shared.AlertMetricProbeFailure, shared.AlertMetricProcessFailure:
		return nil
	default:
		return bladWskazaniaAlertu("nie znam miary „" + string(miara) +
			"” — serwer zna: errorCount, errorRate, cost, tokens, budgetPercent, " +
			"callLatency, probeFailure, processFailure")
	}
}

func (a *adapterAlertow) przelicz(ctx context.Context) {
	reguly, err := a.repozytorium.Reguly(ctx, "", "", true, 0)
	if err != nil {
		return
	}
	teraz := terazWMilisekundachAlertow()
	for _, regula := range reguly {
		a.przeliczRegule(ctx, regula, teraz)
	}
}

func (a *adapterAlertow) przeliczRegule(ctx context.Context, regula dane.RegulaAlertu, teraz int64) {
	if regula.WyciszonaDo != nil && *regula.WyciszonaDo > teraz {
		return
	}
	if regula.OstatnieWyzwolenie != nil && teraz-*regula.OstatnieWyzwolenie < regula.OknoMs {
		return
	}

	wartosc, zmierzono := a.zmierzMiare(ctx, regula, teraz)
	if !zmierzono {
		return
	}
	if !przekroczonoProg(wartosc, regula) {
		return
	}

	komunikat := "miara " + regula.Miara + " wyniosła " + zapisLiczbyMiary(wartosc) +
		" w oknie " + zapisLiczbyMiary(float64(regula.OknoMs)/1000) + " s"
	if regula.Prog != nil {
		komunikat += " przy progu " + zapisLiczbyMiary(*regula.Prog)
	}

	wyzwolenie := dane.WyzwolenieAlertu{
		Kod:                   nowyIdentyfikator(przedrostekWyzwoleniaAlertu),
		RegulaKod:             regula.Kod,
		Stan:                  shared.AlertTriggerStatusFiring,
		Waga:                  regula.Waga,
		Miara:                 regula.Miara,
		WartoscObserwowana:    wartosc,
		Prog:                  regula.Prog,
		Komunikat:             komunikat,
		Wyzwolono:             teraz,
		KanalyDostarczoneJSON: regula.KanalyJSON,
	}
	zapisane, err := a.repozytorium.ZapiszWyzwolenie(ctx, wyzwolenie)
	if err != nil {
		return
	}
	a.rozglosWyzwolenie(ctx, zapisane, regula)
	a.wniesDoCentrum(ctx, zapisane)
}

func (a *adapterAlertow) wniesDoCentrum(ctx context.Context, w dane.WyzwolenieAlertu) {
	if a.centrum == nil {
		return
	}
	_, _ = a.centrum.Zglos(ctx, ZgloszenieCentrum{
		Klasa: shared.NotificationClassBlad,
		Tresc: w.Komunikat,
	})
}

func (a *adapterAlertow) rozglosWyzwolenie(ctx context.Context, w dane.WyzwolenieAlertu, regula dane.RegulaAlertu) {
	if a.nadajnik == nil {
		return
	}
	a.nadajnik.wyzwolenieAlertu(ctx, wyzwolenieKontraktu(w), regulaKontraktu(regula))
}

func (a *adapterAlertow) zmierzMiare(ctx context.Context, regula dane.RegulaAlertu,
	teraz int64) (float64, bool) {

	od := teraz - regula.OknoMs

	switch shared.AlertMetric(regula.Miara) {
	case shared.AlertMetricErrorCount:
		if a.diagnostyka == nil {
			return 0, false
		}
		liczba, err := a.diagnostyka.LiczbaBledow(ctx, dane.FiltrBledow{Od: od, Do: teraz})
		if err != nil {
			return 0, false
		}
		return float64(liczba), true

	case shared.AlertMetricProbeFailure:
		liczba, err := a.repozytorium.LiczbaNieudanychPomiarowSond(ctx, od, teraz)
		if err != nil {
			return 0, false
		}
		return float64(liczba), true

	case shared.AlertMetricProcessFailure:
		liczba, err := a.repozytorium.LiczbaNieudanychPozycjiKolejki(ctx, od, teraz)
		if err != nil {
			return 0, false
		}
		return float64(liczba), true

	case shared.AlertMetricErrorRate, shared.AlertMetricCost, shared.AlertMetricTokens,
		shared.AlertMetricCallLatency, shared.AlertMetricBudgetPercent:
		return a.zmierzZeSladu(ctx, shared.AlertMetric(regula.Miara), od, teraz)

	default:
		return 0, false
	}
}

func (a *adapterAlertow) zmierzZeSladu(ctx context.Context, miara shared.AlertMetric,
	od, do int64) (float64, bool) {

	if a.prowenancja == nil {
		return 0, false
	}
	sumy, err := a.prowenancja.Zuzycie(ctx, wymiarZuzyciaAlertow, &od, &do, 0)
	if err != nil {
		return 0, false
	}

	zadania, bledne, zetony := 0, 0, 0
	koszt, opoznienie := 0.0, 0.0
	for _, suma := range sumy {
		zadania += suma.Zadania
		bledne += suma.ZadaniaBledne
		zetony += suma.TokenyRazem
		koszt += suma.Koszt
		opoznienie += float64(suma.SrednieOpoznienie) * float64(suma.Zadania)
	}

	switch miara {
	case shared.AlertMetricCost:
		return koszt, true
	case shared.AlertMetricTokens:
		return float64(zetony), true
	case shared.AlertMetricErrorRate:
		if zadania == 0 {
			// Okres bez wywołań nie ma udziału niepowodzeń; zero tu znaczyłoby udało się, choć nic się nie działo.
			return 0, false
		}
		return float64(bledne) * 100 / float64(zadania), true
	case shared.AlertMetricCallLatency:
		if zadania == 0 {
			return 0, false
		}
		return opoznienie / float64(zadania), true
	case shared.AlertMetricBudgetPercent:
		pulap := a.pulapKosztu(ctx)
		if pulap <= 0 {
			// Pułap zerowy znaczy pułap zniesiony, więc udziału w budżecie nie ma z czego policzyć.
			return 0, false
		}
		return koszt * 100 / pulap, true
	default:
		return 0, false
	}
}

func (a *adapterAlertow) pulapKosztu(ctx context.Context) float64 {
	if a.rozstrzygacz == nil {
		return 0
	}
	wynik := a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{KontoOperatora: dane.KontoOperatora(ctx)}, kluczPulapuKosztu)
	pulap, err := strconv.ParseFloat(strings.TrimSpace(wynik.Wartosc), 64)
	if err != nil {
		return 0
	}
	return pulap
}

func przekroczonoProg(wartosc float64, regula dane.RegulaAlertu) bool {
	if regula.Prog == nil {
		return wartosc > 0
	}
	porownanie := shared.AlertComparisonGreaterThan
	if regula.Porownanie != nil && strings.TrimSpace(*regula.Porownanie) != "" {
		porownanie = *regula.Porownanie
	}
	switch shared.AlertComparison(porownanie) {
	case shared.AlertComparisonGreaterThan:
		return wartosc > *regula.Prog
	case shared.AlertComparisonGreaterOrEqual:
		return wartosc >= *regula.Prog
	case shared.AlertComparisonLessThan:
		return wartosc < *regula.Prog
	case shared.AlertComparisonLessOrEqual:
		return wartosc <= *regula.Prog
	default:
		return false
	}
}
