// Obszar warsztatu i wdrożeń modułu Apps: `apps.workspace.update` i
// `apps.deployment.run`. Architektura leży w `adapter_modul_aplikacje.go` (ten
// sam typ `adapterAplikacji`), a bieg przebiegu wdrożenia — silnik wykonania i
// obserwator zakończenia — w `adapter_modul_aplikacje_wdrozenie_bieg.go`.
//
// Przebieg wdrożenia realnie przechodzi stany. `apps.deployment.run` zakłada
// przebieg w stanie `pending`, odsyła go operatorowi i uruchamia silnik
// wykonania, który przesuwa go przez `running` do `succeeded` albo `failed` na
// podstawie wykonanej pracy. Krokiem wdrożenia jest sprawdzenie, czy jest co
// wdrożyć: dla wdrożenia w przód — treść przestrzeni roboczej okna (warsztat
// frontendu i backendu), dla cofnięcia — czy wdrożenie docelowe kiedykolwiek
// się powiodło. Pusta przestrzeń albo cofnięcie do przebiegu, który nigdy nie
// wszedł w `succeeded`, kończy się `failed`. Rdzeń nie hostuje produktu, więc
// pole `Url` zostaje puste.
//
// Warsztat Apps to UPSERT po kluczu naturalnym, nie historia wersji:
// `dane.PlikWarsztatu` nie ma odpowiednika `createVersion` modułu Developer —
// `apps.workspace.update` nadpisuje stan bieżący pliku warstwy i oddaje go z
// powrotem (`File` w odpowiedzi, bez pola wersji osobnej od repozytorium).
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

	// Rodzaj zmiany rozstrzyga się przed zapisem, bo zapis jest UPSERT-em i po
	// fakcie nie widać, czy wiersz powstał, czy został nadpisany. Brak pliku
	// nie jest tu odmową — to zwykłe „ten plik dopiero powstaje".
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
	// Rozgłoszenie idzie po udanym zapisie: zdarzenie ma opisywać stan, który
	// naprawdę wszedł do bazy, a nie zamiar. Usterka zapisu kończy komendę
	// odmową i nikogo nie zawiadamia o zmianie, której nie było.
	plikKontraktu := plikWarsztatuKontraktu(zapisany)
	a.rozglosWarsztat(zmiana, oknoKod, zapisany.Warstwa, plikKontraktu)

	return shared.AppsWorkspaceUpdateResponse{
		Layer: zapisany.Warstwa,
		File:  plikKontraktu,
	}, nil
}

// UruchomWdrozenie obsługuje `apps.deployment.run`: zakłada przebieg wdrożenia
// albo cofnięcia w stanie `pending`, odsyła go operatorowi i oddaje silnikowi
// wykonania, który dokończy przejście stanu poza wykonaniem komendy (patrz
// nagłówek pliku i `adapter_modul_aplikacje_wdrozenie_bieg.go`). Komenda kończy
// się, gdy przebieg ruszy, nie gdy się skończy — stan końcowy dochodzi
// zdarzeniem `apps.build.changed`, tak jak `developer.build.run` oddaje wynik
// zdarzeniem `developer.build.changed`.
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
	// Okno musi istnieć, zanim powstanie wiersz przebiegu: wdrożenie bierze
	// treść z przestrzeni roboczej okna żądania, a okna, którego nie ma, nie ma
	// z czego wdrażać. Bez tego sprawdzianu dowolny łańcuch jako `windowId`
	// zakłada wiersz w dzienniku wdrożeń i kończy się powodem o pustym
	// warsztacie zamiast o braku okna.
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
			// Cel cofnięcia jeszcze nie ma werdyktu. Gdyby przebieg mimo to
			// ruszył, jego wynik zależałby od tego, kto pierwszy dobiegnie —
			// silnik odczytuje stan celu dopiero w swojej gorutynie, więc to
			// samo żądanie kończyłoby się raz `succeeded`, raz `failed`.
			// Żądanie jest więc odrzucane od razu i z powodem.
			return shared.AppsDeploymentRunResponse{}, bladWskazaniaAplikacji(
				"wdrożenie docelowe cofnięcia " + *z.RollbackToDeploymentId +
					" jeszcze trwa (stan: " + string(cel.Stan) + ") — poczekaj na jego koniec")
		}
		// Cofnąć da się tylko wdrożenie tego samego okna i tego samego
		// środowiska. Krok cofnięcia (`krokCofniecia`) patrzy wyłącznie na stan
		// celu, więc jedynym miejscem, w którym da się sprawdzić, czyje to
		// wdrożenie, jest przyjęcie żądania. Bez tych dwóch warunków okno bez
		// jednego pliku warsztatu dostaje `succeeded` na production tylko
		// dlatego, że inne okno wdrożyło się kiedyś na dev.
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
	// Przebieg powstał w `pending` — rozgłoś jego założenie i oddaj silnikowi
	// wykonania, który przesunie go przez `running` do stanu końcowego.
	a.rozglosWdrozenie(shared.ChangeKindCreated, zapisane)
	a.uruchomWdrozenie(zapisane)
	return shared.AppsDeploymentRunResponse{Deployment: wdrozenieKontraktu(zapisane)}, nil
}

// sprawdzOknoWdrozenia odmawia wdrożenia w oknie, którego rdzeń nie zna.
// Odmowa opisuje brak, więc kodem jest `not_found`, tak samo jak przy
// nieznanym wdrożeniu docelowym cofnięcia.
//
// Bez wpiętego rejestru okien sprawdzian milczy: moduł ma wtedy pracować, a nie
// odmawiać wszystkiego z braku zależności.
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

// plikWarsztatuKontraktu składa plik warsztatu w kształcie `DeveloperFile`
// kontraktu — `apps.workspace.update` opisuje plik dokładnie tak, jak robi to
// Code Editor, ale bez `VersionId`: warsztat Apps nie ma historii wersji
// (patrz nagłówek pliku).
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

// wdrozenieKontraktu składa przebieg wdrożenia w kształcie `AppDeployment`
// kontraktu. Pole `Status` wraca dokładnie tak, jak zapisane — przejścia stanu
// nadaje silnik wykonania w `adapter_modul_aplikacje_wdrozenie_bieg.go`, ta
// funkcja tylko przekłada zastany wiersz na kontrakt.
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

// sprawdzWarstweWarsztatu i sprawdzSrodowiskoWdrozenia nazywają odmowę na
// wysokości żądania. Wykaz dopuszczalnych wartości niesie kontrakt, a warstwa
// danych trzyma go po raz drugi (CHECK w schemacie, `dane/aplikacje_warsztat.go`).
// Bez tych dwóch funkcji literówka w żądaniu wraca jako `internal_error` z
// treścią SQL-a i z `retryable: true`, czyli rdzeń bierze na siebie cudzą
// pomyłkę i zachęca klienta do powtórzenia żądania, które nie ma prawa się udać.

// sprawdzWarstweWarsztatu dopuszcza wyłącznie warstwy kontraktu.
func sprawdzWarstweWarsztatu(warstwa shared.AppWorkspaceLayer) error {
	if warstwa == shared.AppWorkspaceLayerFrontend || warstwa == shared.AppWorkspaceLayerBackend {
		return nil
	}
	// Treść nie nazywa komendy: ten sam sprawdzian obsługuje zapis
	// (`apps.workspace.update`) i zawężenie odczytu (`apps.workspace.list`).
	return bladWskazaniaAplikacji("nieznana warstwa warsztatu " +
		strconv.Quote(string(warstwa)) + " — dopuszczalne: frontend, backend")
}

// sprawdzSrodowiskoWdrozenia dopuszcza wyłącznie środowiska kontraktu.
func sprawdzSrodowiskoWdrozenia(srodowisko shared.AppDeployEnvironment) error {
	switch srodowisko {
	case shared.AppDeployEnvironmentDev, shared.AppDeployEnvironmentStaging,
		shared.AppDeployEnvironmentProduction:
		return nil
	}
	// Bez nazwy komendy w treści: sprawdzian obsługuje i zapis
	// (`apps.deployment.run`), i zawężenie odczytu (`apps.deployment.list`).
	return bladWskazaniaAplikacji("nieznane środowisko wdrożenia " +
		strconv.Quote(string(srodowisko)) + " — dopuszczalne: dev, staging, production")
}

// sprawdzStrategieWdrozenia dopuszcza wyłącznie strategie kontraktu. Rdzeń
// zapisuje strategię, ale przebieg dla każdej biegnie tak samo — rdzeń nie
// hostuje produktu, więc nie ma czego wystawiać etapami ani na dwa kolory.
// Sprawdzian pilnuje tylko, żeby nie przyjąć wartości, której dziennik i tak
// nie zapisze.
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
