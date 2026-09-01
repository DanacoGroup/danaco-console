// Obszar warsztatu i wdrozen modulu Apps: apps.workspace.update i
// apps.deployment.run. Architektura lezy w osobnym pliku tego samego adaptera,
// a silnik wykonania przebiegu wdrozenia w kolejnym.
package core

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZaktualizujPrzestrzen obsługuje `apps.workspace.update`: zapisuje bieżącą
// treść pliku warstwy frontendu albo backendu produktu w oknie Apps.
func (a *adapterAplikacji) ZaktualizujPrzestrzen(ctx context.Context,
	z shared.AppsWorkspaceUpdateRequest) (shared.AppsWorkspaceUpdateResponse, error) {

	oknoKod := strings.TrimSpace(z.WindowId)
	if oknoKod == "" {
		return shared.AppsWorkspaceUpdateResponse{}, bladWskazaniaAplikacji(
			"apps.workspace.update wymaga okna")
	}
	if strings.TrimSpace(z.Path) == "" {
		return shared.AppsWorkspaceUpdateResponse{}, bladWskazaniaAplikacji(
			"apps.workspace.update wymaga ścieżki pliku")
	}
	if err := sprawdzWarstweWarsztatu(z.Layer); err != nil {
		return shared.AppsWorkspaceUpdateResponse{}, err
	}

	// Rodzaj zmiany rozstrzyga sie przed zapisem, bo zapis jest UPSERT-em i po fakcie nie widac wyniku.
	var zmiana shared.ChangeKind = shared.ChangeKindUpdated
	if _, err := a.repozytorium.PlikWarsztatu(ctx, oknoKod, z.Layer, z.Path); err != nil {
		if !errors.Is(err, dane.ErrBrakWiersza) {
			return shared.AppsWorkspaceUpdateResponse{}, bladAplikacji(err)
		}
		zmiana = shared.ChangeKindCreated
	}

	plik := dane.PlikWarsztatu{
		Okno:        oknoKod,
		Warstwa:     z.Layer,
		Sciezka:     z.Path,
		Tresc:       z.Content,
		KomponentID: z.ComponentId,
	}
	zapisany, err := a.repozytorium.ZapiszPlikWarsztatu(ctx, plik)
	if err != nil {
		return shared.AppsWorkspaceUpdateResponse{}, bladAplikacji(err)
	}
	// Rozgloszenie idzie po udanym zapisie, zeby zdarzenie opisywalo stan, ktory naprawde wszedl do bazy.
	plikKontraktu := plikWarsztatuKontraktu(zapisany)
	a.rozglosWarsztat(ctx, zmiana, oknoKod, zapisany.Warstwa, plikKontraktu)

	return shared.AppsWorkspaceUpdateResponse{
		Layer: zapisany.Warstwa,
		File:  plikKontraktu,
	}, nil
}

// UruchomWdrozenie obsluguje apps.deployment.run: zaklada przebieg wdrozenia
// albo cofniecia w stanie pending i oddaje go silnikowi wykonania, ktory
// dokonczy przejscie stanu poza wykonaniem komendy.
func (a *adapterAplikacji) UruchomWdrozenie(ctx context.Context,
	z shared.AppsDeploymentRunRequest) (shared.AppsDeploymentRunResponse, error) {

	oknoKod := strings.TrimSpace(z.WindowId)
	if oknoKod == "" {
		return shared.AppsDeploymentRunResponse{}, bladWskazaniaAplikacji(
			"apps.deployment.run wymaga okna")
	}
	if strings.TrimSpace(string(z.Environment)) == "" {
		return shared.AppsDeploymentRunResponse{}, bladWskazaniaAplikacji(
			"apps.deployment.run wymaga środowiska wdrożenia")
	}
	if err := sprawdzSrodowiskoWdrozenia(z.Environment); err != nil {
		return shared.AppsDeploymentRunResponse{}, err
	}
	// Okno musi istniec, zanim powstanie wiersz przebiegu, bo wdrozenie bierze tresc z jego przestrzeni.
	if err := a.sprawdzOknoWdrozenia(oknoKod); err != nil {
		return shared.AppsDeploymentRunResponse{}, err
	}

	wdrozenie := dane.WdrozenieApp{
		Kod:            nowyIdentyfikator(przedrostekWdrozeniaApp),
		OknoKod:        oknoKod,
		Srodowisko:     z.Environment,
		Stan:           shared.AppDeployStatusPending,
		Wersja:         z.Version,
		NotatkiWydania: z.ReleaseNotes,
	}
	if z.Strategy != nil {
		if err := sprawdzStrategieWdrozenia(*z.Strategy); err != nil {
			return shared.AppsDeploymentRunResponse{}, err
		}
		wdrozenie.Strategia = *z.Strategy
	} else {
		wdrozenie.Strategia = shared.AppDeployStrategyImmediate
	}

	if z.RollbackToDeploymentId != nil {
		cel, err := a.repozytorium.Wdrozenie(ctx, *z.RollbackToDeploymentId)
		if err != nil {
			return shared.AppsDeploymentRunResponse{}, bladNieznanegoWdrozeniaApp(*z.RollbackToDeploymentId, err)
		}
		if cel.Stan == shared.AppDeployStatusPending || cel.Stan == shared.AppDeployStatusRunning {
			// Cel cofniecia jeszcze nie ma werdyktu, wiec zadanie odrzuca sie od razu, a nie po wyscigu silnika.
			return shared.AppsDeploymentRunResponse{}, bladWskazaniaAplikacji(
				"wdrożenie docelowe cofnięcia " + *z.RollbackToDeploymentId +
					" jeszcze trwa (stan: " + string(cel.Stan) + ") — poczekaj na jego koniec")
		}
		// Cofnac da sie tylko wdrozenie tego samego okna i tego samego srodowiska.
		if cel.OknoKod != oknoKod {
			return shared.AppsDeploymentRunResponse{}, bladWskazaniaAplikacji(
				"wdrożenie docelowe cofnięcia " + *z.RollbackToDeploymentId + " należy do okna " +
					cel.OknoKod + ", a żądanie przyszło z okna " + oknoKod +
					" — cofnąć można wyłącznie wdrożenie tego samego okna")
		}
		if cel.Srodowisko != z.Environment {
			return shared.AppsDeploymentRunResponse{}, bladWskazaniaAplikacji(
				"wdrożenie docelowe cofnięcia " + *z.RollbackToDeploymentId + " poszło na środowisko " +
					string(cel.Srodowisko) + ", a żądanie dotyczy środowiska " + string(z.Environment) +
					" — przywrócić da się tylko stan tego samego środowiska")
		}
		wdrozenie.CofnieteDoKodu = z.RollbackToDeploymentId
		if wdrozenie.Wersja == nil {
			wdrozenie.Wersja = cel.Wersja
		}
	}

	zapisane, err := a.repozytorium.ZapiszWdrozenie(ctx, wdrozenie)
	if err != nil {
		return shared.AppsDeploymentRunResponse{}, bladAplikacji(err)
	}
	// Przebieg powstal w pending — rozglos jego zalozenie i oddaj silnikowi wykonania.
	a.rozglosWdrozenie(ctx, shared.ChangeKindCreated, zapisane)
	a.uruchomWdrozenie(ctx, zapisane)
	return shared.AppsDeploymentRunResponse{Deployment: wdrozenieKontraktu(zapisane)}, nil
}

// sprawdzOknoWdrozenia odmawia wdrozenia w oknie, ktorego rdzen nie zna, kodem
// not_found. Bez wpietego rejestru okien sprawdzian milczy, zamiast odmawiac
// z braku zaleznosci.
func (a *adapterAplikacji) sprawdzOknoWdrozenia(oknoKod string) error {
	if a.okna == nil {
		return nil
	}
	if _, err := a.okna.Okno(oknoKod); err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Apps: okno "+oknoKod+" nie istnieje — wdrożenie bierze treść z przestrzeni"+
				" roboczej okna, więc nie ma go gdzie uruchomić"))
	}
	return nil
}

// plikWarsztatuKontraktu skada plik warsztatu w ksztalcie kontraktu
// DeveloperFile, bez pola wersji: warsztat Apps nie ma historii wersji,
// w odroznieniu od Code Editor.
func plikWarsztatuKontraktu(plik dane.PlikWarsztatu) shared.DeveloperFile {
	tresc := plik.Tresc
	rozmiar := plik.Rozmiar
	return shared.DeveloperFile{
		Path:      plik.Sciezka,
		Content:   &tresc,
		Language:  jezykPliku(plik.Sciezka),
		SizeBytes: &rozmiar,
		UpdatedAt: chwilaBazy(plik.Zaktualizowano),
	}
}

// wdrozenieKontraktu skada przebieg wdrozenia w ksztalcie kontraktu
// AppDeployment; pole Status wraca dokladnie tak, jak zapisane w repozytorium.
func wdrozenieKontraktu(wdrozenie dane.WdrozenieApp) shared.AppDeployment {
	deployment := shared.AppDeployment{
		Id:               wdrozenie.Kod,
		WindowId:         wdrozenie.OknoKod,
		Environment:      wdrozenie.Srodowisko,
		Strategy:         wdrozenie.Strategia,
		Status:           wdrozenie.Stan,
		Version:          wdrozenie.Wersja,
		ReleaseNotes:     wdrozenie.NotatkiWydania,
		Url:              wdrozenie.Adres,
		LogRef:           wdrozenie.LogOdwolanie,
		RolledBackFromId: wdrozenie.CofnieteDoKodu,
		StartedAt:        chwilaBazy(wdrozenie.Rozpoczeto),
	}
	if wdrozenie.Zakonczono != nil {
		koniec := chwilaBazy(*wdrozenie.Zakonczono)
		deployment.FinishedAt = &koniec
	}
	return deployment
}

// Sprawdzenie warstwy i srodowiska zamienia usterke zapisu na czytelna odmowe zadania.

// sprawdzWarstweWarsztatu dopuszcza wylacznie warstwy znane kontraktowi produktu: frontend albo backend Apps.
func sprawdzWarstweWarsztatu(warstwa shared.AppWorkspaceLayer) error {
	if warstwa == shared.AppWorkspaceLayerFrontend || warstwa == shared.AppWorkspaceLayerBackend {
		return nil
	}
	// Tresc nie nazywa komendy: ten sam sprawdzian obsluguje zapis i zawezenie odczytu.
	return bladWskazaniaAplikacji("nieznana warstwa warsztatu " +
		strconv.Quote(string(warstwa)) + " — dopuszczalne: frontend, backend")
}

// sprawdzSrodowiskoWdrozenia dopuszcza wylacznie srodowiska znane kontraktowi: dev, staging albo production.
func sprawdzSrodowiskoWdrozenia(srodowisko shared.AppDeployEnvironment) error {
	switch srodowisko {
	case shared.AppDeployEnvironmentDev, shared.AppDeployEnvironmentStaging,
		shared.AppDeployEnvironmentProduction:
		return nil
	}
	// Bez nazwy komendy w tresci: sprawdzian obsluguje i zapis, i zawezenie odczytu.
	return bladWskazaniaAplikacji("nieznane środowisko wdrożenia " +
		strconv.Quote(string(srodowisko)) + " — dopuszczalne: dev, staging, production")
}

// sprawdzStrategieWdrozenia dopuszcza wylacznie strategie kontraktu; przebieg
// dla kazdej biegnie tak samo, bo rdzen nie hostuje produktu.
func sprawdzStrategieWdrozenia(strategia shared.AppDeployStrategy) error {
	switch strategia {
	case shared.AppDeployStrategyImmediate, shared.AppDeployStrategyStaged,
		shared.AppDeployStrategyBlueGreen:
		return nil
	}
	return bladWskazaniaAplikacji("apps.deployment.run: nieznana strategia wdrożenia " +
		strconv.Quote(string(strategia)) + " — dopuszczalne: immediate, staged, blueGreen")
}

// bladNieznanegoWdrozeniaApp odróżnia „wdrożenia nie ma" (cel cofnięcia poza
// dziennikiem) od usterki odczytu.
func bladNieznanegoWdrozeniaApp(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Apps: wdrożenie do cofnięcia nie istnieje: "+kod))
	}
	return bladAplikacji(err)
}
