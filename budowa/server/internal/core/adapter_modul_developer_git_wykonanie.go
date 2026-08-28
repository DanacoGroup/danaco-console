// Odpowiedzialność pliku: uruchomienie gita przez port warstwy kanału
// session.Uruchamiacz i złożenie wyniku czynności wraz ze stanem repozytorium.
// Konflikt scalenia czy brak gałęzi są wynikami czynności, nie awariami.
package core

import (
	"context"
	"io"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// programGita jest jedynym programem, który ten moduł uruchamia dla Git
// Panelu, portem session.Uruchamiacz.
const programGita = "git"

// wykonajGit uruchamia jedną czynność repozytorium i składa jej wynik
// wraz ze stanem repozytorium po wykonaniu.
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

// wynikGita niesie surowy rezultat jednego uruchomienia gita: powodzenie,
// treść wyjścia i powód niepowodzenia.
type wynikGita struct {
	udane bool
	tresc string
	powod string
}

// uruchomGit startuje gita, zbiera całe wyjście i czeka na zakończenie,
// synchronicznie. Granica czasu pilnuje, żeby git czekający na hasło do
// repozytorium zdalnego nie zatrzymał obsługiwacza na zawsze.
func (a *adapterDevelopera) uruchomGit(ctx context.Context, okno session.Okno,
	argumenty []string, granica time.Duration) (wynikGita, error) {

	if a.uruchamiacz == nil {
		return wynikGita{}, bladWykonaniaDevelopera("rdzeń nie ma uruchamiacza procesów")
	}
	polecenie, err := a.polecenieDopuszczoneDevelopera(okno, programGita, argumenty)
	if err != nil {
		return wynikGita{}, err
	}
	uchwyt, err := a.uruchamiacz.UruchomProces(okno, polecenie)
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
		powod = "żądanie przerwane przez rdzeń"
		_ = drzewo.Ubij()
		bladZakonczenia = <-zakonczenie
	}

	tresc := <-zebrane + <-zebrane
	return wynikGita{udane: bladZakonczenia == nil && powod == "", tresc: tresc, powod: powod}, nil
}

// polecenieDopuszczoneDevelopera składa polecenie i przepuszcza je przez
// egzekutor izolacji okna. Naruszenie izolacji wraca jako odmowa uprawnienia,
// nie jako usterka wykonania.
func (a *adapterDevelopera) polecenieDopuszczoneDevelopera(okno session.Okno, program string,
	argumenty []string) (session.Polecenie, error) {

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
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfigKontekstOkna(okno))
	}
	dopuszczone, err := session.SprawdzPolecenie(zasady, a.obszarDevelopera(okno), polecenie)
	if err != nil {
		return session.Polecenie{}, odmowaIzolacjiDevelopera(err)
	}
	return dopuszczone, nil
}

// czytajDoKonca zbiera cały strumień procesu. Wyjście gita jest krótkie, więc
// idzie w całości do wyniku czynności. Błąd odczytu zostawia to, co zdążyło
// przyjść: urwane wyjście gita niesie więcej niż jego brak.
func czytajDoKonca(zrodlo io.Reader) string {
	if zrodlo == nil {
		return ""
	}
	bajty, _ := io.ReadAll(zrodlo)
	return string(bajty)
}
