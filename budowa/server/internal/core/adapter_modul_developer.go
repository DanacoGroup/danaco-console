// Odpowiedzialność pliku: moduł Developer — wypełnienie portu Developer pięcioma
// komendami obszaru `developer.*`. Obszar okna i sprawdzenie ścieżek leżą
// w `_okno.go`, praca z plikiem w `_plik.go`, drzewo w `_drzewo.go`,
// repozytorium w `_git.go` i `_git_wykonanie.go`, budowanie w `_budowanie.go`
// i `_budowanie_bieg.go`, ewidencja przebiegów w `_rejestr.go`, a przekład na
// kontrakt w `_przeklad.go`.
//
// Przed każdą zmianą stoją dwie bramy, w tej kolejności:
//  1. tryb uprawnień okna — czy wolno w ogóle zmienić stan systemu
//     (PermissionMode, `_okno.go`); tryb `plan` wyklucza zapis pliku,
//     czynność repozytorium i uruchomienie budowania;
//  2. obszar okna — czy ścieżka mieści się w katalogach roboczych okna
//     (`_okno.go`), a przy uruchomieniu procesu dodatkowo egzekutor izolacji
//     (session.SprawdzPolecenie nad zasadami z `izolacja.go`).
//
// Rdzeń nie buduje własnego `exec.Cmd`: git i zadanie budowania startuje ten sam
// port session.Uruchamiacz, którym jedzie okno rozmowy i moduł Terminal —
// Terminal jest warstwą wykonawczą budowania, instalacji zależności
// i uruchamiania.
package core

import (
	"context"
	"log"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Zgodność adaptera z portem sprawdzana jest przy kompilacji.
var _ Developer = (*adapterDevelopera)(nil)

const (
	// przedrostekBudowania znakuje identyfikator przebiegu budowania.
	przedrostekBudowania = "build-"
	// przedrostekWersjiPliku znakuje migawkę treści pliku założoną przy zapisie.
	przedrostekWersjiPliku = "fver-"
)

// adapterDevelopera wypełnia port Developer.
type adapterDevelopera struct {
	repozytorium dane.RepozytoriumDevelopera
	rejestr      *rejestrBudowan
	// sesjeDebugowania trzyma biegi debuggera czynne w tej chwili. Stan żywy,
	// nie zapis: sesja ma uchwyt do procesu adaptera i gaśnie razem z nim.
	sesjeDebugowania *rejestrSesjiDebugowania
	// okna daje tryb uprawnień okna, jego sesję i listę katalogów roboczych.
	// Bez niego moduł nie tknie ani jednego pliku: nie wiedziałby, w jakim
	// obszarze wolno mu pracować, a praca „gdziekolwiek" nie jest pracą
	// w warunkach niepełnych danych, tylko wyjściem poza izolację okna.
	okna *session.Rejestr
	// uruchamiacz jest portem warstwy kanału — jedyną drogą startu procesu.
	uruchamiacz session.Uruchamiacz
	// rozstrzygacz i katalog składają zasady izolacji obowiązujące w oknie oraz
	// katalog zastępczy dla okna bez własnej listy katalogów roboczych.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// kanaly są rejestrem kanałów modelu — jedyną drogą operacji kontekstowych
	// (`developer.contextual.op`). Bez niego moduł działa w całości poza tą
	// jedną rodziną, która odpowiada wtedy odmową nazywającą brak.
	kanaly *models.Rejestr
	// przyrost rozgłasza `developer.build.changed`. Podpina go obsługiwacz.
	przyrost func(shared.ChangeKind, shared.DeveloperBuild, string)
}

// nowyAdapterDevelopera wiąże port z rejestrem okien i uruchamiaczem procesów.
func nowyAdapterDevelopera(okna *session.Rejestr, uruchamiacz session.Uruchamiacz) *adapterDevelopera {
	return &adapterDevelopera{
		rejestr:          nowyRejestrBudowan(),
		sesjeDebugowania: nowyRejestrSesjiDebugowania(),
		okna:             okna,
		uruchamiacz:      uruchamiacz,
	}
}

// ZTrwaloscia podpina wersje plików i dziennik budowań
// (`migracja_042_developer.sql`). Bez
// niego moduł pracuje w pamięci jednego biegu rdzenia: zapis pliku
// nadal działa, lecz `createVersion` nie ma gdzie założyć migawki i zapis
// odpowiada odmową zamiast cicho gubić wersję.
func (a *adapterDevelopera) ZTrwaloscia(repozytorium dane.RepozytoriumDevelopera) *adapterDevelopera {
	a.repozytorium = repozytorium
	return a
}

// ZIzolacja podpina rozstrzygacz zasięgu i ustalacz katalogu roboczego — dwa
// źródła, z których powstaje obszar okna egzekwowany przy każdej ścieżce.
func (a *adapterDevelopera) ZIzolacja(rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterDevelopera {
	a.rozstrzygacz, a.katalog = rozstrzygacz, katalog
	return a
}

// ZKanalami wpina rejestr kanałów modelu — drogę operacji kontekstowych paska
// pływającego Code Editora.
func (a *adapterDevelopera) ZKanalami(kanaly *models.Rejestr) *adapterDevelopera {
	a.kanaly = kanaly
	return a
}

// Przygotuj osierocą przebiegi budowania, do których rdzeń stracił uchwyt przy
// restarcie. Wywołuje się to raz, przy montażu.
//
// Bez tego kroku Build Output pokazywałby przebieg oznaczony jako trwający,
// którego nikt już nie prowadzi i którego nie da się przerwać.
func (a *adapterDevelopera) Przygotuj(ctx context.Context) error {
	if a.repozytorium == nil {
		return nil
	}
	_, err := a.repozytorium.OsierocPrzebiegi(ctx)
	return err
}

// Zamknij przerywa przebiegi budowania czynne w chwili zatrzymania rdzenia.
// Budowanie przeżywa rozłączenie klienta, lecz nie przeżywa końca rdzenia:
// bez uchwytu zostałoby sierotą poza rejestrem.
func (a *adapterDevelopera) Zamknij() {
	if a == nil {
		return
	}
	a.rejestr.Zamknij()
	a.sesjeDebugowania.Zamknij()
}

// przygotujDevelopera odtwarza stan modułu Developer przy montażu rdzenia.
// Niepowodzenie nie zatrzymuje startu — idzie do dziennika, bo
// przebieg oznaczony jako trwający wprowadzałby Operatora w błąd.
func przygotujDevelopera(kontekst context.Context, developer *adapterDevelopera, dziennik *log.Logger) {
	if err := developer.Przygotuj(kontekst); err != nil && dziennik != nil {
		dziennik.Printf("moduł Developer: nie można osierocić przebiegów budowania: %v", err)
	}
}
