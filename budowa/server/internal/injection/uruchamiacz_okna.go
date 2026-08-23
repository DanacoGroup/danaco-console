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

// Uruchamianie procesu okna komunikacji po stronie kanału.
//
// Pakiet session nie startuje procesów — zna wyłącznie port session.Uruchamiacz
// i wypełnia go ten plik. Dzięki temu w drzewie jest jedna droga uruchomienia
// procesu modelu (rozruch.go) i jedna droga ubijania jego drzewa potomstwa
// (session.PrzejmijDrzewo).
//
// Kierunek zależności jest ten sam, który zapowiada session/przejecie.go:
// warstwa kanału sięga po część sesyjną, nigdy odwrotnie.
//
// Zasięg wykonania jest czytany i ma skutek. Okno niesie wybór
// Operatora w polu SrodowiskoWykonania i to tutaj — w jedynym spawnerze
// platformy — ten wybór rozstrzyga, gdzie proces rusza:
//
//	core   — host rdzenia; proces rusza tutaj i to jest wykonanie zgodne z wyborem;
//	remote — host zdalny; proces jedzie torem SSH pakietu internal/zdalne:
//	         host wskazuje ustawienie `host_wykonania`, zgodę per
//	         host trzyma tabela `host_zdalny`, a każde brakujące
//	         ogniwo drogi jest osobną, nazwaną odmową — nie cichym startem
//	         na maszynie rdzenia, bo to byłaby praca w innym miejscu, niż
//	         wskazał Operator;
//	local  — urządzenie Operatora; toru zwrotnego do urządzenia w drzewie nie ma
//	         i nie domknie go ta warstwa: powłoka natywna wystawia interfejsowi
//	         trzy polecenia bez uruchamiania procesów (desktop/src-tauri,
//	         invoke_handler), a kontrakt nie ma kanału, którym rdzeń prowadziłby
//	         strumienie procesu na kliencie — tor zwrotny wymaga nowych poleceń
//	         powłoki i nowych komend kontraktu (contract.* poza tym pakietem).
//	         Zasięg obsługuje więc host rdzenia i idzie o tym wpis do dziennika;
//	         wybór jest honorowany dosłownie dopóty, dopóki rdzeń stoi na
//	         urządzeniu Operatora — a tak stoi dziś każda instalacja lokalna.
//
// Wartość pusta i wartość spoza wyliczenia nie zatrzymują pracy — schodzą na
// zachowanie dotychczasowe (host rdzenia), ale zostawiają ślad w dzienniku,
// bo brak wskazania ma dawać pracę, nie odmowę.

// uruchamiaczOkien wypełnia port session.Uruchamiacz.
type uruchamiaczOkien struct{}

// UruchamiaczOkien zwraca uruchamiacz procesów okien dla warstwy sesji.
// Wynik podaje się nadzorcy przy montażu: session.ZUruchamiaczem(...).
func UruchamiaczOkien() session.Uruchamiacz {
	return uruchamiaczOkien{}
}

// UruchomProces startuje proces okna i oddaje jego uchwyt sesji. Proces rusza
// jako korzeń własnego drzewa, bo zaraz po starcie obejmie go uchwyt systemowy
// warstwy sesji; bez tego wnuki procesu przeżyłyby zamknięcie okna.
//
// Przed startem rozstrzygany jest zasięg wykonania okna: proces, którego nie da
// się uruchomić tam, gdzie wskazał Operator, nie rusza tutaj po cichu.
func (u uruchamiaczOkien) UruchomProces(o session.Okno, p session.Polecenie) (session.UchwytProcesu, error) {
	rozruch, err := rozruchWedlugZasiegu(o, p)
	if err != nil {
		return nil, err
	}
	return Wystartuj(context.Background(), rozruch)
}

// rozruchWedlugZasiegu odpowiada na jedno pytanie: jaki proces uruchomić na
// hoście rdzenia, żeby praca działa się tam, gdzie wskazał Operator. Dla
// zasięgu `core` (i dróg schodzących na niego) jest to sam proces okna; dla
// zasięgu `remote` — proces transportu SSH, którego strumienie są strumieniami
// procesu na hoście zdalnym.
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

// rozruchZdalny prowadzi polecenie okna torem SSH pakietu zdalne. Katalog
// i środowisko jadą w komendzie zdalnej; proces transportu dziedziczy
// środowisko rdzenia, bo ssh potrzebuje własnej konfiguracji (klucze, agent).
// Odmowa toru wraca do wołającego z powodem — uruchomienie na hoście rdzenia
// byłoby pracą w innym miejscu, niż wskazał Operator.
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

// odnotowaneZasiegi pamięta adnotacje już opisane w dzienniku. Moduł Terminal
// startuje proces przy każdym poleceniu, więc wpis przy każdym starcie
// utopiłby dziennik w powtórzeniach i wyszłoby z tego to samo, co z ciszy:
// nikt by tego nie czytał. Zmiana zasięgu okna — a przy torze zdalnym także
// zmiana wyniku (odmowa kontra tor) — daje nowy klucz, więc kolejny obrót
// sprawy znów zostawia ślad.
var odnotowaneZasiegi sync.Map

// odnotujZasiegRaz zapisuje adnotację raz na trójkę okno–zasięg–wynik.
func odnotujZasiegRaz(o session.Okno, wynik string, wzor string, argumenty ...any) {
	klucz := o.Id + "\x00" + string(o.SrodowiskoWykonania) + "\x00" + wynik
	if _, byl := odnotowaneZasiegi.LoadOrStore(klucz, struct{}{}); byl {
		return
	}
	log.Printf(wzor, argumenty...)
}
