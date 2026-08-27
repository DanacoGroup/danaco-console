// Moduł Developer wypełnia port Developer pięciu komendami obszaru `developer.*` wraz z pracą na pliku, drzewie, repozytorium i budowaniem, rozdzielonymi po plikach pomocniczych.
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

// Zgodność adaptera z portem sprawdzana jest przy kompilacji, wprost przez pustą asercję interfejsu Developer.
var _ Developer = (*adapterDevelopera)(nil)

const (
	// przedrostekBudowania znakuje identyfikator przebiegu budowania nadawany przez rdzen przy jego zalozeniu w bazie danych.
	przedrostekBudowania = "build-"
	// przedrostekWersjiPliku znakuje migawkę treści pliku założoną przy zapisie do repozytorium wersji plikow.
	przedrostekWersjiPliku = "fver-"
)

// adapterDevelopera wypełnia port Developer, niosąc repozytorium, rejestr biegów i zależności wpinane osobno.
type adapterDevelopera struct {
	repozytorium dane.RepozytoriumDevelopera
	rejestr      *rejestrBudowan
	// sesjeDebugowania trzyma biegi debuggera czynne w tej chwili; gaśnie razem z adapterem.
	sesjeDebugowania *rejestrSesjiDebugowania
	// okna daje tryb uprawnień okna, sesję i listę katalogów roboczych; bez niego moduł nie tknie pliku.
	okna *session.Rejestr
	// uruchamiacz jest portem warstwy kanału — jedyną drogą startu procesu.
	uruchamiacz session.Uruchamiacz
	// rozstrzygacz i katalog składają zasady izolacji obowiązujące w oknie oraz katalog zastępczy.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// kanaly są rejestrem kanałów modelu, jedyną drogą operacji kontekstowych (`developer.contextual.op`).
	kanaly *models.Rejestr
	// przyrost rozgłasza `developer.build.changed`. Podpina go obsługiwacz.
	przyrost func(shared.ChangeKind, shared.DeveloperBuild, string)
}

// nowyAdapterDevelopera wiąże port z rejestrem okien i uruchamiaczem procesów, oddając adapter gotowy do dalszego wpięcia zależności.
func nowyAdapterDevelopera(okna *session.Rejestr, uruchamiacz session.Uruchamiacz) *adapterDevelopera {
	return &adapterDevelopera{
		rejestr:          nowyRejestrBudowan(),
		sesjeDebugowania: nowyRejestrSesjiDebugowania(),
		okna:             okna,
		uruchamiacz:      uruchamiacz,
	}
}

// ZTrwaloscia podpina wersje plików i dziennik budowań (`migracja_042_developer.sql`). Bez niego moduł pracuje w pamięci jednego biegu rdzenia i `createVersion` odpowiada odmową zamiast cicho gubić wersję.
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

// Przygotuj osierocą przebiegi budowania, do których rdzeń stracił uchwyt przy restarcie. Wywołuje się raz, przy montażu. Bez tego kroku Build Output pokazywałby przebieg, którego nie da się przerwać.
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
