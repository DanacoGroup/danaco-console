// Odpowiedzialność pliku: stan repozytorium okna Git Panel — gałąź, ostatnie
// zatwierdzenie, ścieżki zmienione i ścieżki z konfliktem.
//
// Stan doczytuje się po każdej czynności, bo kontrakt daje Git Panelowi jedną
// komendę — `developer.git.action` — i nie ma w nim osobnej komendy odczytu
// stanu repozytorium, historii ani różnicy. Bez stanu dołączonego do wyniku
// czynności okno po zatwierdzeniu zmian dalej pokazywałoby wykaz sprzed niego.
package core

import (
	"context"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// granicaOdczytuStanu jest krótka: odczyt stanu nie sięga do sieci, a klient
// czeka na niego razem z wynikiem czynności.
const granicaOdczytuStanu = 20 * time.Second

// stanGita niesie migawkę repozytorium po czynności.
type stanGita struct {
	galaz         string
	zatwierdzenie string
	zmienione     []string
	konflikty     []string
	podsumowanie  string
}

// stanRepozytorium czyta gałąź, ostatnie zatwierdzenie i wykaz zmian.
// Niepowodzenie któregokolwiek odczytu zostawia pole puste — czynność, która
// się wykonała, nie ma prawa zostać unieważniona przez nieudany odczyt.
func (a *adapterDevelopera) stanRepozytorium(ctx context.Context, okno session.Okno) stanGita {
	stan := stanGita{}
	if galaz, ok := a.odczytGita(ctx, okno, "rev-parse", "--abbrev-ref", "HEAD"); ok {
		stan.galaz = pierwszyWiersz(galaz)
	}
	if wersja, ok := a.odczytGita(ctx, okno, "rev-parse", "--short", "HEAD"); ok {
		stan.zatwierdzenie = pierwszyWiersz(wersja)
	}
	if wykaz, ok := a.odczytGita(ctx, okno, "status", "--porcelain"); ok {
		stan.zmienione, stan.konflikty = rozbierzStatus(wykaz)
		stan.podsumowanie = podsumowanieStanu(stan)
	}
	return stan
}

// odczytGita uruchamia odczytowe polecenie gita i oddaje jego wyjście.
func (a *adapterDevelopera) odczytGita(ctx context.Context, okno session.Okno,
	argumenty ...string) (string, bool) {

	wynik, err := a.uruchomGit(ctx, okno, argumenty, granicaOdczytuStanu)
	if err != nil || !wynik.udane {
		return "", false
	}
	return wynik.tresc, true
}

// rozbierzStatus czyta wyjście `git status --porcelain` na dwa wykazy.
//
// Konflikt rozpoznajemy po kodach obu stron indeksu — `UU`, `AA`, `DD` oraz
// wszystkich parach z literą `U`. Bez tego rozróżnienia Git Panel pokazywałby
// plik z konfliktem jako zwykłą zmianę, a znaczniki scalenia trafiłyby do
// zatwierdzenia razem z kodem.
func rozbierzStatus(wyjscie string) (zmienione, konflikty []string) {
	zmienione = make([]string, 0, 16)
	konflikty = make([]string, 0, 4)
	for _, wiersz := range strings.Split(wyjscie, "\n") {
		wiersz = strings.TrimRight(wiersz, "\r")
		if len(wiersz) < 4 {
			continue
		}
		kod, sciezka := wiersz[:2], strings.TrimSpace(wiersz[3:])
		// Zmiana nazwy ma postać „stara -> nowa”; do wykazu wchodzi nowa.
		if miejsce := strings.LastIndex(sciezka, " -> "); miejsce >= 0 {
			sciezka = sciezka[miejsce+4:]
		}
		sciezka = strings.Trim(sciezka, `"`)
		if sciezka == "" {
			continue
		}
		if czyKonflikt(kod) {
			konflikty = append(konflikty, sciezka)
			continue
		}
		zmienione = append(zmienione, sciezka)
	}
	return zmienione, konflikty
}

// czyKonflikt rozpoznaje parę kodów statusu oznaczającą scalenie nierozstrzygnięte.
func czyKonflikt(kod string) bool {
	if len(kod) < 2 {
		return false
	}
	if strings.ContainsRune(kod, 'U') {
		return true
	}
	return kod == "AA" || kod == "DD"
}

// podsumowanieStanu składa wiersz stanu doklejany do wyjścia czynności: gałąź,
// ostatnie zatwierdzenie oraz liczby ścieżek zmienionych i konfliktowych.
func podsumowanieStanu(stan stanGita) string {
	czesci := make([]string, 0, 4)
	if stan.galaz != "" {
		czesci = append(czesci, "gałąź "+stan.galaz)
	}
	if stan.zatwierdzenie != "" {
		czesci = append(czesci, "zatwierdzenie "+stan.zatwierdzenie)
	}
	czesci = append(czesci, "zmienionych ścieżek "+itoa(len(stan.zmienione)))
	if len(stan.konflikty) > 0 {
		czesci = append(czesci, "konfliktów "+itoa(len(stan.konflikty)))
	}
	return "[stan repozytorium: " + strings.Join(czesci, ", ") + "]"
}

// pierwszyWiersz odcina wszystko po pierwszym końcu wiersza.
func pierwszyWiersz(tresc string) string {
	tresc = strings.TrimSpace(tresc)
	if miejsce := strings.IndexAny(tresc, "\r\n"); miejsce >= 0 {
		return strings.TrimSpace(tresc[:miejsce])
	}
	return tresc
}

// konfigKontekstOkna składa kontekst rozstrzygania zasad izolacji dla okna.
func konfigKontekstOkna(okno session.Okno) konfig.Kontekst {
	return konfig.Kontekst{Okno: okno.Id}
}

// odmowaIzolacjiDevelopera znakuje naruszenie izolacji kodem uprawnienia.
// Usterka innego rodzaju idzie dalej bez zmiany — kod `permission_denied` ma
// znaczyć zatrzymanie przez izolację, a nie „coś się nie udało”.
func odmowaIzolacjiDevelopera(err error) error {
	if err == nil {
		return nil
	}
	if !errors.Is(err, session.ErrIzolacja) {
		return err
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
		"moduł Developer: "+err.Error()))
}
