// Stan repozytorium okna Git Panel doczytywany po każdej czynności, bo kontrakt
// daje jedną komendę bez osobnego odczytu stanu.
package core

import (
	"context"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

const granicaOdczytuStanu = 20 * time.Second

type stanGita struct {
	galaz         string
	zatwierdzenie string
	zmienione     []string
	konflikty     []string
	podsumowanie  string
}

// Czynność, która się wykonała, nie może zostać unieważniona przez nieudany odczyt stanu.
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

func (a *adapterDevelopera) odczytGita(ctx context.Context, okno session.Okno,
	argumenty ...string) (string, bool) {

	wynik, err := a.uruchomGit(ctx, okno, argumenty, granicaOdczytuStanu)
	if err != nil || !wynik.udane {
		return "", false
	}
	return wynik.tresc, true
}

// Konflikt rozpoznaje się po kodach obu stron indeksu: `UU`, `AA`, `DD` i pary z literą `U`.
func rozbierzStatus(wyjscie string) (zmienione, konflikty []string) {
	zmienione = make([]string, 0, 16)
	konflikty = make([]string, 0, 4)
	for _, wiersz := range strings.Split(wyjscie, "\n") {
		wiersz = strings.TrimRight(wiersz, "\r")
		if len(wiersz) < 4 {
			continue
		}
		kod, sciezka := wiersz[:2], strings.TrimSpace(wiersz[3:])
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

func czyKonflikt(kod string) bool {
	if len(kod) < 2 {
		return false
	}
	if strings.ContainsRune(kod, 'U') {
		return true
	}
	return kod == "AA" || kod == "DD"
}

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

func pierwszyWiersz(tresc string) string {
	tresc = strings.TrimSpace(tresc)
	if miejsce := strings.IndexAny(tresc, "\r\n"); miejsce >= 0 {
		return strings.TrimSpace(tresc[:miejsce])
	}
	return tresc
}

func konfigKontekstOkna(ctx context.Context, okno session.Okno) konfig.Kontekst {
	return konfig.Kontekst{Okno: okno.Id, KontoOperatora: dane.KontoOperatora(ctx)}
}

// Kod `permission_denied` znaczy zatrzymanie przez izolację; inna usterka idzie dalej bez zmiany.
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
