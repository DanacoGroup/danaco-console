// Silnik wykonania przebiegu wdrożenia modułu Apps: przesuwa przebieg założony
// przez `apps.deployment.run` (`adapter_modul_aplikacje_wdrozenie.go`) ze stanu
// `pending` przez `running` do `succeeded` albo `failed` i rozgłasza każde
// przejście zdarzeniem `apps.build.changed`.
//
// Silnik pracuje poza żądaniem. Komenda `apps.deployment.run` kończy się, gdy
// przebieg ruszy, a nie gdy się skończy — rozłączenie klienta w połowie nie
// przerywa wdrożenia. Dlatego bieg idzie własną gorutyną i własnym kontekstem
// (`context.Background`), a nie kontekstem komendy.
//
// Stan końcowy wynika z wykonanej pracy. Krokiem wdrożenia jest sprawdzenie,
// czy jest co wdrożyć:
//   - wdrożenie w przód udaje się, gdy przestrzeń robocza okna niesie choć jeden
//     plik; pusta przestrzeń kończy się `failed`;
//   - cofnięcie udaje się, gdy wdrożenie docelowe kiedykolwiek weszło w
//     `succeeded`; cofnięcie do przebiegu, który nigdy się nie powiódł, kończy
//     się `failed`.
//
// Rdzeń nie hostuje produktu, więc wdrożenie w przód zostawia pole `Url` puste —
// postawienie serwera produktu wymagałoby uruchamiacza procesu wpiętego w
// adapter. Cofnięcie dziedziczy `Url` wprost z wdrożenia docelowego.
//
// PRZEBIEG ZOSTAWIA PO SOBIE ARTEFAKT I DZIENNIK. Udane wdrożenie w przód pakuje
// przestrzeń roboczą okna w archiwum `zip` i kładzie je w magazynie treści
// rdzenia, a wiersz `artefakt_apps` wskazuje ten plik wraz z rozmiarem i sumą
// kontrolną — to on jest wejściem `apps.package.build` i pozycją
// `apps.artifact.list`. Każdy krok przebiegu dopisuje wiersz do
// `wiersz_dziennika_apps`, skąd czyta go `apps.deployment.log.read`. Bez tych
// dwóch rzeczy trzy komendy rodziny meldowałyby pustkę przy przebiegu, który
// naprawdę się odbył.
package core

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekArtefaktuApp znakuje identyfikatory artefaktów budowania.
const przedrostekArtefaktuApp = "art-"

// errBrakMagazynuArtefaktuApp nazywa brak katalogu danych przy składaniu
// artefaktu. Osobny błąd, a nie odmowa komendy: przebieg wdrożenia już się
// zakończył i jego stanu ten brak nie zmienia.
var errBrakMagazynuArtefaktuApp = errors.New(
	"moduł Apps: rdzeń zmontowano bez katalogu danych, więc artefakt nie ma gdzie powstać")

// uruchomWdrozenie oddaje przebieg silnikowi wykonania. Bieg idzie osobną
// gorutyną: przejście przez stany dzieje się po odesłaniu odpowiedzi komendy.
func (a *adapterAplikacji) uruchomWdrozenie(wdrozenie dane.WdrozenieApp) {
	go a.wykonajWdrozenie(wdrozenie)
}

// wykonajWdrozenie przesuwa przebieg z `pending` przez `running` do stanu
// końcowego i rozgłasza każde przejście. Kontekst jest własny, nie komendy —
// wdrożenie przeżywa rozłączenie klienta (patrz nagłówek pliku).
func (a *adapterAplikacji) wykonajWdrozenie(wdrozenie dane.WdrozenieApp) {
	ctx := context.Background()

	wdrozenie.Stan = shared.AppDeployStatusRunning
	wdrozenie = a.zapiszIRozglosWdrozenie(ctx, wdrozenie)
	a.dopiszDziennikApp(ctx, wdrozenie.OknoKod, &wdrozenie.Kod, nil,
		"przebieg "+wdrozenie.Kod+" ruszył na środowisko "+string(wdrozenie.Srodowisko)+
			" strategią "+string(wdrozenie.Strategia))

	stan, adres, powod := a.krokWdrozenia(ctx, wdrozenie)
	wdrozenie.Stan = stan
	if adres != nil && wdrozenie.Adres == nil {
		wdrozenie.Adres = adres
	}
	if powod != "" {
		odwolanie := powod
		wdrozenie.LogOdwolanie = &odwolanie
	}
	zakonczono := time.Now().UTC().Format(formatZnacznikaBazy)
	wdrozenie.Zakonczono = &zakonczono

	// Artefakt powstaje wyłącznie po udanym wdrożeniu w przód: cofnięcie
	// przywraca wersję już zapakowaną, a przebieg nieudany nie ma czego wydać.
	if stan == shared.AppDeployStatusSucceeded && wdrozenie.CofnieteDoKodu == nil {
		if opis, err := a.zlozArtefaktPrzebieguApp(ctx, wdrozenie); err != nil {
			a.dopiszDziennikApp(ctx, wdrozenie.OknoKod, &wdrozenie.Kod, nil,
				"artefaktu przebiegu nie udało się złożyć: "+err.Error())
		} else {
			a.dopiszDziennikApp(ctx, wdrozenie.OknoKod, &wdrozenie.Kod, nil, opis)
		}
	}

	a.zapiszIRozglosWdrozenie(ctx, wdrozenie)
	if powod != "" {
		a.dopiszDziennikApp(ctx, wdrozenie.OknoKod, &wdrozenie.Kod, nil,
			"przebieg "+wdrozenie.Kod+" zakończony stanem "+string(stan)+": "+powod)
		return
	}
	a.dopiszDziennikApp(ctx, wdrozenie.OknoKod, &wdrozenie.Kod, nil,
		"przebieg "+wdrozenie.Kod+" zakończony stanem "+string(stan))
}

// zlozArtefaktPrzebieguApp pakuje przestrzeń roboczą okna w archiwum i zapisuje
// je w magazynie treści rdzenia wraz z wierszem artefaktu. Zwraca zdanie do
// dziennika albo powód niepowodzenia.
//
// Brak magazynu nie przewraca przebiegu: wdrożenie już się udało, a artefakt
// jest jego wynikiem ubocznym. Dziennik mówi wtedy wprost, czego zabrakło —
// milczenie kazałoby szukać artefaktu, którego nikt nie miał gdzie odłożyć.
func (a *adapterAplikacji) zlozArtefaktPrzebieguApp(ctx context.Context,
	wdrozenie dane.WdrozenieApp) (string, error) {

	if a.magazyn == nil {
		return "", errBrakMagazynuArtefaktuApp
	}
	pliki, err := a.repozytorium.PlikiWarsztatu(ctx, wdrozenie.OknoKod)
	if err != nil {
		return "", err
	}
	// Nazwa wpisu niesie warstwę, bo ten sam plik może stać w obu warstwach
	// warsztatu pod tą samą ścieżką — klucz naturalny to (okno, warstwa, ścieżka).
	wpisy := map[string][]byte{}
	nazwy := make([]string, 0, len(pliki))
	for _, plik := range pliki {
		nazwa := string(plik.Warstwa) + "/" + strings.TrimPrefix(plik.Sciezka, "/")
		wpisy[nazwa] = []byte(plik.Tresc)
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)

	archiwum, err := archiwumZipApp(nazwy, wpisy)
	if err != nil {
		return "", err
	}
	odwolanie, rozmiar, err := a.wniesDoMagazynuApp(archiwum)
	if err != nil {
		return "", err
	}
	suma := sumaTresciApp(archiwum)
	kod := nowyIdentyfikator(przedrostekArtefaktuApp)
	if _, err := a.repozytorium.ZalozArtefaktApp(ctx, dane.ArtefaktApp{
		Kod: kod, Okno: wdrozenie.OknoKod, WdrozenieKod: &wdrozenie.Kod,
		Rodzaj: shared.AppArtifactKindBundle, Sciezka: odwolanie,
		Rozmiar: &rozmiar, SumaKontrolna: &suma,
	}); err != nil {
		return "", err
	}
	return "artefakt " + kod + " złożony z " + strconv.Itoa(len(nazwy)) + " plików warsztatu (" +
		strconv.FormatInt(rozmiar, 10) + " bajtów, sha256 " + suma + ")", nil
}

// krokWdrozenia wykonuje właściwą pracę przebiegu i zwraca stan końcowy, adres
// produktu (gdy jest prawdziwy) oraz powód niepowodzenia (gdy nastąpiło).
func (a *adapterAplikacji) krokWdrozenia(ctx context.Context,
	wdrozenie dane.WdrozenieApp) (shared.AppDeployStatus, *string, string) {

	if wdrozenie.CofnieteDoKodu != nil {
		return a.krokCofniecia(ctx, *wdrozenie.CofnieteDoKodu)
	}
	return a.krokWdrozeniaWPrzod(ctx, wdrozenie.OknoKod)
}

// krokWdrozeniaWPrzod sprawdza, czy przestrzeń robocza okna niesie coś do
// wdrożenia. Pusta przestrzeń kończy przebieg realnym niepowodzeniem.
func (a *adapterAplikacji) krokWdrozeniaWPrzod(ctx context.Context,
	okno string) (shared.AppDeployStatus, *string, string) {

	pliki, err := a.repozytorium.PlikiWarsztatu(ctx, okno)
	if err != nil {
		return shared.AppDeployStatusFailed, nil,
			"nie można odczytać przestrzeni roboczej okna " + okno + ": " + err.Error()
	}
	if len(pliki) == 0 {
		return shared.AppDeployStatusFailed, nil,
			"przestrzeń robocza okna " + okno + " jest pusta — nie ma czego wdrożyć"
	}
	return shared.AppDeployStatusSucceeded, nil, ""
}

// krokCofniecia sprawdza, czy wdrożenie docelowe kiedykolwiek się powiodło —
// tylko takie da się przywrócić. Adres produktu dziedziczy się z celu, bo jest
// wartością prawdziwą (rdzeń niczego nie hostuje na nowo).
func (a *adapterAplikacji) krokCofniecia(ctx context.Context,
	kodCelu string) (shared.AppDeployStatus, *string, string) {

	cel, err := a.repozytorium.Wdrozenie(ctx, kodCelu)
	if err != nil {
		return shared.AppDeployStatusFailed, nil,
			"nie można odczytać wdrożenia docelowego cofnięcia " + kodCelu + ": " + err.Error()
	}
	// Powód rozróżnia dwa przypadki, bo prowadzą do różnych działań: przy
	// przebiegu trwającym trzeba poczekać, przy padłym — wybrać inny cel.
	// Żądanie cofnięcia do przebiegu bez werdyktu odrzuca już obsługiwacz
	// komendy (`adapter_modul_aplikacje_wdrozenie.go`); tutaj zostaje przypadek,
	// w którym cel zmienił stan między odczytem komendy a odczytem silnika.
	if cel.Stan == shared.AppDeployStatusPending || cel.Stan == shared.AppDeployStatusRunning {
		return shared.AppDeployStatusFailed, nil,
			"wdrożenie docelowe cofnięcia " + kodCelu +
				" jeszcze trwa (stan: " + string(cel.Stan) + ") — nie ma czego przywrócić"
	}
	if cel.Stan != shared.AppDeployStatusSucceeded {
		return shared.AppDeployStatusFailed, nil,
			"wdrożenie docelowe cofnięcia " + kodCelu +
				" nigdy się nie powiodło (stan: " + string(cel.Stan) + ")"
	}
	return shared.AppDeployStatusSucceeded, cel.Adres, ""
}

// zapiszIRozglosWdrozenie utrwala przejście stanu i rozgłasza je. Nieudany
// zapis nie gubi przejścia: przebieg rozgłasza się wtedy w kształcie z pamięci,
// żeby okno zobaczyło stan mimo usterki dziennika.
func (a *adapterAplikacji) zapiszIRozglosWdrozenie(ctx context.Context,
	wdrozenie dane.WdrozenieApp) dane.WdrozenieApp {

	zapisane, err := a.repozytorium.ZapiszWdrozenie(ctx, wdrozenie)
	if err != nil {
		a.rozglosWdrozenie(shared.ChangeKindUpdated, wdrozenie)
		return wdrozenie
	}
	a.rozglosWdrozenie(shared.ChangeKindUpdated, zapisane)
	return zapisane
}

// rozglosWdrozenie oddaje przebieg obsługiwaczowi, który rozsyła
// `apps.build.changed`. Brak podpięcia nie zmienia pracy modułu.
func (a *adapterAplikacji) rozglosWdrozenie(zmiana shared.ChangeKind, wdrozenie dane.WdrozenieApp) {
	if a.przyrostWdrozenia == nil {
		return
	}
	a.przyrostWdrozenia(zmiana, wdrozenieKontraktu(wdrozenie))
}
