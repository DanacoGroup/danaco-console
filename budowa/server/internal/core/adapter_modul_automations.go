// Moduł Automations obsługuje definicję automatyki (Workflow Builder) i jej wykaz; harmonogram, układ zależności, kolejka i przebiegi mają własne pliki adaptera. Kroki wykonuje jeden silnik kolejek nad adapterem kolejek.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów modułu. Byt nadany przez rdzeń wychodzi
// kontraktem pod własnym identyfikatorem, nie pod kluczem wiersza.
const (
	przedrostekAutomatyki   = "automat-"
	przedrostekKroku        = "krok-"
	przedrostekHarmonogramu = "harm-"
	przedrostekWyzwalacza   = "wyzw-"
	przedrostekPrzebiegu    = "przebieg-"
)

// adapterAutomatyk wypełnia port Automatyki. Zależności są trzy: repozytorium
// modułu, adapter kolejek (jedyny wykonawca kroków) i pamięć obserwatorów
// przebiegów, przez którą Execution Monitor przypina się do telemetrii.
type adapterAutomatyk struct {
	repozytorium dane.RepozytoriumAutomatyk
	kolejki      *adapterKolejek
	obserwatorzy *pamiecObserwatorowPrzebiegow
	// uklad niesie dopełnienia układu zależności; nil znaczy, że cztery komendy odmawiają.
	uklad dane.RepozytoriumUkladuOrkiestracji
	// sejf jest magazynem wartości poświadczeń; nil znaczy, że `automation.secret.set` odmawia.
	sejf SejfPoswiadczenAutomatyki
	// okna służą wyłącznie przełożeniu roli środowiska MultitaskingAI na okno przy spięciu.
	okna dane.RepozytoriumOkien
}

// nowyAdapterAutomatyk wiąże port Automatyki z repozytorium modułu, oddając gotowy adapter do dalszego wpięcia zależności.
func nowyAdapterAutomatyk(repozytorium dane.RepozytoriumAutomatyk) *adapterAutomatyk {
	return &adapterAutomatyk{
		repozytorium: repozytorium,
		obserwatorzy: nowaPamiecObserwatorowPrzebiegow(),
	}
}

// ZKolejkami podpina adapter kolejek. Bez niego moduł buduje i pokazuje
// definicje, lecz `automation.queue.action` odmawia wprost — brak wykonawcy nie
// może udawać wykonania.
func (a *adapterAutomatyk) ZKolejkami(kolejki *adapterKolejek) *adapterAutomatyk {
	a.kolejki = kolejki
	return a
}

// ZSejfem podpina magazyn wartości poświadczeń — ten sam sejf plikowy, którym
// jadą sekrety kont i punktów dostępu. Jedna instancja pod jednym zamkiem: dwa
// sejfy nad tym samym plikiem ścigałyby się o zapis.
func (a *adapterAutomatyk) ZSejfem(sejf SejfPoswiadczenAutomatyki) *adapterAutomatyk {
	a.sejf = sejf
	return a
}

// Zapisz zapisuje definicję automatyki wraz z krokami. Brak `workflowId` zakłada automatykę nową; wskazanie zmienia zastaną i podnosi numer wersji. Kroki podmieniają się w całości, gdy pole `steps` przyszło; pole nieobecne zostawia kroki nietknięte.
func (a *adapterAutomatyk) Zapisz(ctx context.Context,
	z shared.AutomationWorkflowSaveRequest) (shared.AutomationWorkflowSaveResponse, error) {

	kod := wartoscTekstu(z.WorkflowId)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekAutomatyki)
	}
	zapisana, err := a.repozytorium.ZapiszAutomatyke(ctx, dane.Automatyka{
		Kod: kod, Nazwa: z.Name, Opis: z.Description, Czynna: czyCzynna(z.Enabled),
	})
	if err != nil {
		return shared.AutomationWorkflowSaveResponse{}, bladAutomatyki(err)
	}
	if z.Steps != nil {
		if err := a.zapiszKroki(ctx, zapisana.ID, z.Steps); err != nil {
			return shared.AutomationWorkflowSaveResponse{}, bladAutomatyki(err)
		}
	}
	automatyka, err := a.zloz(ctx, zapisana)
	if err != nil {
		return shared.AutomationWorkflowSaveResponse{}, bladAutomatyki(err)
	}
	// Migawka wersji idzie PO złożeniu automatyki, bo zapisuje stan naprawdę stojący w bazie.
	a.odlozWersje(ctx, zapisana.ID, zapisana.Wersja, automatyka.Steps)
	a.zapisAudytu(ctx, &zapisana.ID, "zapis definicji automatyki",
		map[string]any{"wersja": zapisana.Wersja, "krokow": len(automatyka.Steps)})
	return shared.AutomationWorkflowSaveResponse{Workflow: automatyka}, nil
}

// odlozWersje zapisuje migawkę definicji. Nieudany zapis migawki nie wywraca zapisu definicji: definicja już stoi w bazie, a usterka historii jest usterką panelu „Wersje", nie usterką zapisu.
func (a *adapterAutomatyk) odlozWersje(ctx context.Context, automatykaID int64, wersja int,
	kroki []shared.AutomationStep) {

	zapis := zapisStrukturalny(kroki)
	if zapis == nil {
		return
	}
	_ = a.repozytorium.ZapiszWersjeAutomatyki(ctx, dane.WersjaAutomatyki{
		AutomatykaID: automatykaID, Wersja: wersja, Kroki: *zapis,
	})
}

// Wykaz zwraca automatyki dostępne Operatorowi, każdą razem z krokami i ich
// zależnościami — wykaz Workflow Buildera pokazuje strukturę, nie same nazwy.
func (a *adapterAutomatyk) Wykaz(ctx context.Context,
	z shared.AutomationWorkflowListRequest) (shared.AutomationWorkflowListResponse, error) {

	tylkoCzynne := z.EnabledOnly != nil && *z.EnabledOnly
	wiersze, err := a.repozytorium.Automatyki(ctx, tylkoCzynne, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.AutomationWorkflowListResponse{}, bladAutomatyki(err)
	}
	automatyki := make([]shared.AutomationWorkflow, 0, len(wiersze))
	for _, wiersz := range wiersze {
		automatyka, err := a.zloz(ctx, wiersz)
		if err != nil {
			return shared.AutomationWorkflowListResponse{}, bladAutomatyki(err)
		}
		automatyki = append(automatyki, automatyka)
	}
	return shared.AutomationWorkflowListResponse{Workflows: automatyki}, nil
}

// Automatyka oddaje definicję kontraktu po kodzie. Służy oknom, które znają
// wyłącznie identyfikator automatyki — Scheduler, Orchestrator, Execution Monitor.
func (a *adapterAutomatyk) Automatyka(ctx context.Context, kod string) (shared.AutomationWorkflow, error) {
	wiersz, err := a.wiersz(ctx, kod)
	if err != nil {
		return shared.AutomationWorkflow{}, err
	}
	automatyka, err := a.zloz(ctx, wiersz)
	if err != nil {
		return shared.AutomationWorkflow{}, bladAutomatyki(err)
	}
	return automatyka, nil
}

// wiersz odnajduje automatykę po kodzie i nazywa brak wprost. Automatyka
// nieznana jest błędem żądania, nie awarią rdzenia — Operator wskazał byt,
// którego nie ma.
func (a *adapterAutomatyk) wiersz(ctx context.Context, kod string) (dane.Automatyka, error) {
	if kod == "" {
		return dane.Automatyka{}, bladWskazaniaAutomatyki("komenda bez wskazania automatyki")
	}
	wiersz, err := a.repozytorium.Automatyka(ctx, kod)
	if err != nil {
		return dane.Automatyka{}, bladNieznanejAutomatyki(kod, err)
	}
	return wiersz, nil
}

// zloz składa automatykę kontraktu z wiersza wraz z krokami i zależnościami odczytanymi z repozytorium.
func (a *adapterAutomatyk) zloz(ctx context.Context, wiersz dane.Automatyka) (shared.AutomationWorkflow, error) {
	kroki, err := a.krokiKontraktu(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationWorkflow{}, err
	}
	// Etykiety idą razem z definicją; nieudany odczyt etykiet daje wykaz pusty, nie usterkę.
	etykiety, err := a.repozytorium.EtykietyAutomatyki(ctx, wiersz.ID)
	if err != nil {
		etykiety = nil
	}
	// Wersja oddawana kontraktem jest wersją WYKONYWANĄ, nie roboczą.
	wersja := wiersz.Wersja
	if wiersz.WersjaOpublikowana != nil {
		wersja = *wiersz.WersjaOpublikowana
	}
	return shared.AutomationWorkflow{
		Id: wiersz.Kod, Name: wiersz.Nazwa, Description: wiersz.Opis,
		Steps: kroki, Enabled: wiersz.Czynna, Version: &wersja, Tags: etykiety,
		CreatedAt: chwilaBazy(wiersz.Utworzono), UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}, nil
}

// czyCzynna rozstrzyga przełącznik „automatyka czynna”. Brak wskazania znaczy
// automatykę czynną: Operator, który właśnie ją zapisał, chce nią pracować.
func czyCzynna(wskazanie *bool) bool {
	return wskazanie == nil || *wskazanie
}

// wartoscLiczby zdejmuje wskaźnik z pola opcjonalnego; brak znaczy zero, a zero
// znaczy „bez wskazania granicy”.
func wartoscLiczby(wskazanie *int) int {
	if wskazanie == nil {
		return 0
	}
	return *wskazanie
}

// bladAutomatyki znakuje usterkę kodem kontraktu, żeby okno modułu pokazało powód. Błąd, któremu kod już nadano, przechodzi tędy bez zmiany; dopiero usterka bez kodu staje się usterką wewnętrzną rdzenia.
func bladAutomatyki(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaAutomatyki nazywa brak danych w żądaniu — to błąd Operatora,
// nie rdzenia, więc kod odmowy jest inny niż przy usterce.
func bladWskazaniaAutomatyki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Automations: "+powod))
}

// bladNieznanejAutomatyki odróżnia „automatyki nie ma” od „odczyt się nie
// powiódł”. Okno pokazuje wtedy inny komunikat i inaczej podpowiada Operatorowi.
func bladNieznanejAutomatyki(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Automations: automatyka nie istnieje: "+kod))
	}
	return bladAutomatyki(err)
}
