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

// wymiarZuzyciaAlertow jest osią, po której grupuje się ślad wywołań na
// potrzeby miar. Kanał, bo każde wywołanie ma kanał; sumy po wszystkich
// wierszach dają wartość globalną, a rozbicie zostaje na przyszłe zawężenie
// reguły do zasięgu.
const wymiarZuzyciaAlertow = "channel"

// sprawdzMiareReguly odbija miarę spoza wyliczenia kontraktu oraz miarę, której źródła rdzeń nie ma wpiętego — reguła bez źródła przechodziłaby ewaluację w ciszy i wyglądała na czynną.
func (a *adapterAlertow) sprawdzMiareReguly(miara shared.AlertMetric) error {
	switch miara {
	case shared.AlertMetricErrorCount:
		if a.diagnostyka == nil {
			return bladWskazaniaAlertu("miara errorCount liczy się z dziennika błędów rdzenia, " +
				"a rdzeń nie ma go wpiętego — naprawa: podpiąć repozytorium diagnostyki " +
				"przy składaniu rdzenia")
		}
		return nil
	case shared.AlertMetricErrorRate, shared.AlertMetricCost, shared.AlertMetricTokens,
		shared.AlertMetricCallLatency:
		if a.prowenancja == nil {
			return bladWskazaniaAlertu("miara " + string(miara) + " liczy się ze śladu wywołań " +
				"modelu, a rdzeń nie ma go wpiętego — naprawa: podpiąć repozytorium " +
				"prowenancji przy składaniu rdzenia")
		}
		return nil
	case shared.AlertMetricBudgetPercent:
		if a.prowenancja == nil || a.rozstrzygacz == nil {
			return bladWskazaniaAlertu("miara budgetPercent liczy koszt wobec nastawy " +
				kluczPulapuKosztu + ", więc potrzebuje śladu wywołań i rozstrzygacza " +
				"konfiguracji — rdzeń nie ma wpiętego przynajmniej jednego z nich")
		}
		return nil
	case shared.AlertMetricProbeFailure, shared.AlertMetricProcessFailure:
		return nil
	default:
		return bladWskazaniaAlertu("nie znam miary „" + string(miara) +
			"” — rdzeń zna: errorCount, errorRate, cost, tokens, budgetPercent, " +
			"callLatency, probeFailure, processFailure")
	}
}

// przelicz ewaluuje wszystkie reguły czynne i zapisuje wyzwolenia tych, które przekroczyły próg; metoda nie zwraca błędu, bo niepowodzenie nie ma prawa odebrać wykazu wyzwoleń wcześniejszych.
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

// przeliczRegule wykonuje jeden pomiar i, gdy próg został przekroczony, zapisuje wyzwolenie; wyciszenie i powtórzenie sprawdza się przed pomiarem, żeby reguła nie zawołała dwa razy za to samo.
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
	a.rozglosWyzwolenie(zapisane, regula)
	a.wniesDoCentrum(ctx, zapisane)
}

// wniesDoCentrum zapisuje wyzwolony alert w rejestrze centrum powiadomień jedną klasą blad; niepowodzenie zapisu nie przerywa niczego, bo alarm już został zapisany i rozgłoszony.
func (a *adapterAlertow) wniesDoCentrum(ctx context.Context, w dane.WyzwolenieAlertu) {
	if a.centrum == nil {
		return
	}
	_, _ = a.centrum.Zglos(ctx, ZgloszenieCentrum{
		Klasa: shared.NotificationClassBlad,
		Tresc: w.Komunikat,
	})
}

// rozglosWyzwolenie nadaje zdarzenie `alert.triggered`. Zdarzenie jest drogą
// alertu do okien — bez niego Operator dowiedziałby się o wyzwoleniu dopiero
// przy następnym otwarciu wykazu.
func (a *adapterAlertow) rozglosWyzwolenie(w dane.WyzwolenieAlertu, regula dane.RegulaAlertu) {
	if a.nadajnik == nil {
		return
	}
	a.nadajnik.wyzwolenieAlertu(wyzwolenieKontraktu(w), regulaKontraktu(regula))
}

// zmierzMiare wykonuje pomiar jednej miary w oknie czasu reguły. Drugi zwracany
// wynik mówi, czy pomiar się odbył — miara, której nie zmierzono, nie wyzwala
// niczego, bo nie ma wartości do porównania z progiem.
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

// zmierzZeSladu liczy miary pochodzące ze śladu wywołań modelu.
//
// Jeden odczyt na miarę, nie pięć: sumy po wymiarze niosą komplet potrzebnych
// składników, a drugie zapytanie o ten sam okres dałoby liczby z innej chwili.
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

// pulapKosztu odczytuje nastawę pułapu kosztu z konfiguracji obowiązującej ocenę reguły alertu tej maszyny.
func (a *adapterAlertow) pulapKosztu(_ context.Context) float64 {
	if a.rozstrzygacz == nil {
		return 0
	}
	wynik := a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{}, kluczPulapuKosztu)
	pulap, err := strconv.ParseFloat(strings.TrimSpace(wynik.Wartosc), 64)
	if err != nil {
		return 0
	}
	return pulap
}

// przekroczonoProg rozstrzyga, czy zmierzona wartość wyzwala regułę; brak progu znaczy każda wartość niezerowa, a brak porównania znaczy więcej niż.
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
