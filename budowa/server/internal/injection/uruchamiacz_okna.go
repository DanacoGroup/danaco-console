// Uruchamianie procesu okna komunikacji po stronie kanału, zgodnie
// z zasięgiem wykonania wskazanym w oknie.
package injection

import (
	"context"
	"fmt"
	"log"
	"sync"

	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zdalne"
	"danacoconsole/shared"
)

// uruchamiaczOkien wypełnia port session.Uruchamiacz, uruchamiając procesy
// okien komunikacji platformy.
type uruchamiaczOkien struct{}

// UruchamiaczOkien zwraca uruchamiacz procesów okien dla warstwy sesji.
// Wynik podaje się nadzorcy przy montażu: session.ZUruchamiaczem(...).
func UruchamiaczOkien() session.Uruchamiacz {
	return uruchamiaczOkien{}
}

// UruchomProces startuje proces okna i oddaje jego uchwyt sesji, jako korzeń
// własnego drzewa procesów.
func (u uruchamiaczOkien) UruchomProces(o session.Okno, p session.Polecenie) (session.UchwytProcesu, error) {
	rozruch, err := rozruchWedlugZasiegu(o, p)
	if err != nil {
		return nil, err
	}
	return Wystartuj(context.Background(), rozruch)
}

// rozruchWedlugZasiegu odpowiada na jedno pytanie: jaki proces uruchomić na
// hoście rdzenia, żeby praca działa się tam, gdzie wskazał Operator.
func rozruchWedlugZasiegu(o session.Okno, p session.Polecenie) (Rozruch, error) {
	switch o.SrodowiskoWykonania {
	case shared.ExecutionEnvCore:
		return rozruchMiejscowy(p), nil

	case shared.ExecutionEnvRemote:
		return rozruchZdalny(o, p)

	case shared.ExecutionEnvLocal:
		odnotujZasiegRaz(o, "", "injection: okno %s ma zasięg wykonania %q; toru zwrotnego "+
			"do urządzenia Operatora w drzewie nie ma (powłoka nie uruchamia procesów, "+
			"kontrakt nie niesie ich strumieni), więc proces rusza na hoście rdzenia — "+
			"zgodnie z wyborem tylko dopóty, dopóki rdzeń stoi na urządzeniu Operatora",
			o.Id, string(o.SrodowiskoWykonania))
		return rozruchMiejscowy(p), nil

	case "":
		odnotujZasiegRaz(o, "", "injection: okno %s nie ma wskazanego zasięgu wykonania; "+
			"proces rusza na hoście rdzenia", o.Id)
		return rozruchMiejscowy(p), nil

	default:
		odnotujZasiegRaz(o, "", "injection: okno %s ma zasięg wykonania %q spoza wyliczenia "+
			"ExecutionEnv; proces rusza na hoście rdzenia jak dla zasięgu %q",
			o.Id, string(o.SrodowiskoWykonania), shared.ExecutionEnvCore)
		return rozruchMiejscowy(p), nil
	}
}

// rozruchMiejscowy przekłada polecenie okna na rozruch na hoście rdzenia —
// zachowanie dotychczasowe, wspólne dla zasięgu core i dróg na niego schodzących.
func rozruchMiejscowy(p session.Polecenie) Rozruch {
	return Rozruch{
		Program:             p.Program,
		Argumenty:           p.Argumenty,
		Katalog:             p.Katalog,
		Srodowisko:          p.Srodowisko,
		Atrybuty:            session.AtrybutyDrzewa(),
		WyjscieBledowOsobno: true,
	}
}

// rozruchZdalny prowadzi polecenie okna torem SSH pakietu zdalne, dziedzicząc
// środowisko rdzenia dla własnej konfiguracji.
func rozruchZdalny(o session.Okno, p session.Polecenie) (Rozruch, error) {
	uruchomienie, err := zdalne.Przeloz(o.Id, zdalne.Polecenie{
		Program:    p.Program,
		Argumenty:  p.Argumenty,
		Katalog:    p.Katalog,
		Srodowisko: p.Srodowisko,
	})
	if err != nil {
		err = fmt.Errorf("injection: okno %s ma zasięg wykonania %q, a tor do hosta "+
			"zdalnego odmówił: %w", o.Id, string(o.SrodowiskoWykonania), err)
		odnotujZasiegRaz(o, "odmowa", "%v", err)
		return Rozruch{}, err
	}
	odnotujZasiegRaz(o, "tor", "injection: okno %s ma zasięg wykonania %q; proces jedzie "+
		"torem SSH na host %s (%s)", o.Id, string(o.SrodowiskoWykonania),
		uruchomienie.Host.Nazwa, uruchomienie.Host.AdresPolaczenia())
	return Rozruch{
		Program:             uruchomienie.Program,
		Argumenty:           uruchomienie.Argumenty,
		Atrybuty:            session.AtrybutyDrzewa(),
		WyjscieBledowOsobno: true,
	}, nil
}

// odnotowaneZasiegi pamięta adnotacje już zapisane w dzienniku, żeby uniknąć
// powtórzeń przy każdym starcie procesu.
var odnotowaneZasiegi sync.Map

// odnotujZasiegRaz zapisuje adnotację w dzienniku raz na trójkę okno, zasięg
// i wynik samego uruchomienia.
func odnotujZasiegRaz(o session.Okno, wynik string, wzor string, argumenty ...any) {
	klucz := o.Id + "\x00" + string(o.SrodowiskoWykonania) + "\x00" + wynik
	if _, byl := odnotowaneZasiegi.LoadOrStore(klucz, struct{}{}); byl {
		return
	}
	log.Printf(wzor, argumenty...)
}
