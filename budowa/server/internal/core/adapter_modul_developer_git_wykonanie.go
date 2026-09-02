// Odpowiedzialność pliku: uruchomienie gita przez port session.Uruchamiacz i złożenie wyniku czynności ze stanem repozytorium.
package core

import (
	"context"
	"io"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

const programGita = "git"

func (a *adapterDevelopera) wykonajGit(ctx context.Context, okno session.Okno,
	czynnosc shared.GitActionKind, argumenty []string, granica time.Duration) (shared.GitActionResult, error) {

	wyjscie, err := a.uruchomGit(ctx, okno, argumenty, granica)
	if err != nil {
		return shared.GitActionResult{}, err
	}

	wynik := shared.GitActionResult{Action: czynnosc, Succeeded: wyjscie.udane}
	tresc := "$ " + programGita + " " + strings.Join(argumenty, " ") + "\n" + wyjscie.tresc
	if wyjscie.powod != "" {
		tresc += "\n[" + wyjscie.powod + "]"
	}

	// Stan repozytorium odczytuje się po czynności, jedyny moment zgodny z Git Panel.
	stan := a.stanRepozytorium(ctx, okno)
	if stan.galaz != "" {
		wynik.Branch = &stan.galaz
	}
	if stan.zatwierdzenie != "" {
		wynik.CommitId = &stan.zatwierdzenie
	}
	wynik.ChangedPaths = stan.zmienione
	wynik.ConflictPaths = stan.konflikty
	if stan.podsumowanie != "" {
		tresc += "\n" + stan.podsumowanie
	}
	wynik.Output = &tresc
	return wynik, nil
}

type wynikGita struct {
	udane bool
	tresc string
	powod string
}

// Granica czasu: git czekający na hasło repozytorium zdalnego nie może zatrzymać obsługiwacza.
func (a *adapterDevelopera) uruchomGit(ctx context.Context, okno session.Okno,
	argumenty []string, granica time.Duration) (wynikGita, error) {

	if a.uruchamiacz == nil {
		return wynikGita{}, bladWykonaniaDevelopera("serwer nie ma uruchamiacza procesów")
	}
	polecenie, err := a.polecenieDopuszczoneDevelopera(ctx, okno, programGita, argumenty)
	if err != nil {
		return wynikGita{}, err
	}
	uchwyt, err := a.uruchamiacz.UruchomProces(ctx, okno, polecenie)
	if err != nil {
		return wynikGita{}, bladWykonaniaDevelopera(
			"nie można uruchomić polecenia git w " + polecenie.Katalog + ": " + err.Error())
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		return wynikGita{}, bladWykonaniaDevelopera("nie można objąć drzewa procesu git: " + err.Error())
	}
	defer drzewo.Zwolnij()

	zebrane := make(chan string, 2)
	go func() { zebrane <- czytajDoKonca(uchwyt.Wyjscie()) }()
	go func() { zebrane <- czytajDoKonca(uchwyt.Diagnostyka()) }()

	zakonczenie := make(chan error, 1)
	go func() { zakonczenie <- uchwyt.Czekaj() }()

	zegar := time.NewTimer(granica)
	defer zegar.Stop()

	var powod string
	var bladZakonczenia error
	select {
	case bladZakonczenia = <-zakonczenie:
	case <-zegar.C:
		powod = "przekroczona granica czasu " + granica.String()
		_ = drzewo.Ubij()
		bladZakonczenia = <-zakonczenie
	case <-ctx.Done():
		powod = "żądanie przerwane przez serwer"
		_ = drzewo.Ubij()
		bladZakonczenia = <-zakonczenie
	}

	tresc := <-zebrane + <-zebrane
	return wynikGita{udane: bladZakonczenia == nil && powod == "", tresc: tresc, powod: powod}, nil
}

func (a *adapterDevelopera) polecenieDopuszczoneDevelopera(ctx context.Context, okno session.Okno,
	program string, argumenty []string) (session.Polecenie, error) {

	korzenie := a.korzenieOkna(okno)
	if len(korzenie) == 0 {
		return session.Polecenie{}, bladZasobuDevelopera(
			"okno nie ma katalogu roboczego, więc nie ma gdzie uruchomić polecenia " + program)
	}
	polecenie := session.Polecenie{
		Program:             program,
		Argumenty:           argumenty,
		Katalog:             korzenie[0],
		DziedziczSrodowisko: true,
	}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfigKontekstOkna(ctx, okno))
	}
	dopuszczone, err := session.SprawdzPolecenie(zasady, a.obszarDevelopera(okno), polecenie)
	if err != nil {
		return session.Polecenie{}, odmowaIzolacjiDevelopera(err)
	}
	return dopuszczone, nil
}

func czytajDoKonca(zrodlo io.Reader) string {
	if zrodlo == nil {
		return ""
	}
	bajty, _ := io.ReadAll(zrodlo)
	return string(bajty)
}
