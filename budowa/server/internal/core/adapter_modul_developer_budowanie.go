// Plik obsługuje `developer.build.run`: uruchomienie i przerwanie budowania okna Build Output. Komenda kończy się, gdy budowanie ruszy, a nie gdy się skończy; stan i log dochodzą zdarzeniem `developer.build.changed`.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Budowanie obsługuje `developer.build.run`, rozbierając zadanie na program i wiodące parametry przed uruchomieniem.
func (a *adapterDevelopera) Budowanie(ctx context.Context,
	z shared.DeveloperBuildRunRequest) (shared.DeveloperBuildRunResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperBuildRunResponse{}, err
	}
	if z.Stop != nil && *z.Stop {
		return a.przerwijBudowanie(ctx, okno.Id)
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "uruchomienie budowania"); err != nil {
		return shared.DeveloperBuildRunResponse{}, err
	}

	program, wiodace, err := rozbierzZadanie(z.Task)
	if err != nil {
		return shared.DeveloperBuildRunResponse{}, err
	}
	argumenty := append(wiodace, parametryZadania(z.Arguments)...)
	polecenie, err := a.polecenieDopuszczoneDevelopera(ctx, okno, program, argumenty)
	if err != nil {
		return shared.DeveloperBuildRunResponse{}, err
	}

	przebieg, err := a.uruchomBudowanie(ctx, okno, polecenie, strings.TrimSpace(z.Task), argumenty)
	if err != nil {
		return shared.DeveloperBuildRunResponse{}, err
	}
	return shared.DeveloperBuildRunResponse{Build: budowanieKontraktu(przebieg)}, nil
}

// przerwijBudowanie kończy przebieg czynny w oknie. Okno bez czynnego przebiegu nie jest błędem żądania: Operator mógł nacisnąć „zatrzymaj" sekundę po zakończeniu; odpowiedzią jest wtedy ostatni znany przebieg z dziennika.
func (a *adapterDevelopera) przerwijBudowanie(ctx context.Context,
	oknoKod string) (shared.DeveloperBuildRunResponse, error) {

	if przebieg, jest := a.rejestr.Przebieg(oknoKod); jest {
		if err := przebieg.Przerwij(); err != nil {
			return shared.DeveloperBuildRunResponse{}, bladWykonaniaDevelopera(
				"nie można przerwać budowania " + przebieg.kod + ": " + err.Error())
		}
		// Obserwator zakończenia domknie przebieg i rozgłosi zmianę; odpowiedź niesie stan z chwili polecenia.
		przebieg.CzekajNaKoniec(czasNaDomknieciePrzebiegu)
		return shared.DeveloperBuildRunResponse{Build: budowanieKontraktu(przebieg)}, nil
	}

	ostatni, err := a.ostatniPrzebiegOkna(ctx, oknoKod)
	if err != nil {
		return shared.DeveloperBuildRunResponse{}, err
	}
	return shared.DeveloperBuildRunResponse{Build: ostatni}, nil
}

// ostatniPrzebiegOkna czyta ostatni przebieg okna z dziennika budowań, gdy przebieg czynny już nie istnieje.
func (a *adapterDevelopera) ostatniPrzebiegOkna(ctx context.Context,
	oknoKod string) (shared.DeveloperBuild, error) {

	if a.repozytorium == nil {
		return shared.DeveloperBuild{}, bladZasobuDevelopera(
			"w oknie " + oknoKod + " nie biegnie żadne budowanie, a serwer nie ma dziennika przebiegów")
	}
	wiersz, err := a.repozytorium.OstatniPrzebieg(ctx, oknoKod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.DeveloperBuild{}, bladZasobuDevelopera(
			"w oknie " + oknoKod + " nie biegnie żadne budowanie i nie ma go w dzienniku")
	}
	if err != nil {
		return shared.DeveloperBuild{}, bladWykonaniaDevelopera(
			"nie można odczytać dziennika budowań okna " + oknoKod + ": " + err.Error())
	}
	return budowanieZDziennika(wiersz), nil
}

// rozbierzZadanie dzieli zadanie na program i jego wiodące parametry według pierwszego słowa wiersza polecenia.
func rozbierzZadanie(zadanie string) (string, []string, error) {
	czesci := strings.Fields(strings.TrimSpace(zadanie))
	if len(czesci) == 0 {
		return "", nil, bladZadaniaDevelopera("budowanie wymaga wskazania zadania")
	}
	return czesci[0], czesci[1:], nil
}

// parametryZadania odsiewa puste parametry żądania, zanim trafią do wywołania egzekutora izolacji okna.
func parametryZadania(zadane []string) []string {
	parametry := make([]string, 0, len(zadane))
	for _, parametr := range zadane {
		if tresc := strings.TrimSpace(parametr); tresc != "" {
			parametry = append(parametry, tresc)
		}
	}
	return parametry
}
