// Moduł Automations obsługuje definicję automatyki (Workflow Builder) i jej wykaz; harmonogram, układ zależności, kolejka i przebiegi mają własne pliki adaptera. Kroki wykonuje jeden silnik kolejek nad adapterem kolejek.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/nadajnik"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	przedrostekAutomatyki   = "automat-"
	przedrostekKroku        = "krok-"
	przedrostekHarmonogramu = "harm-"
	przedrostekWyzwalacza   = "wyzw-"
	przedrostekPrzebiegu    = "przebieg-"
)

type adapterAutomatyk struct {
	repozytorium dane.RepozytoriumAutomatyk
	kolejki      *adapterKolejek
	obserwatorzy *pamiecObserwatorowPrzebiegow
	uklad        dane.RepozytoriumUkladuOrkiestracji
	sejf         SejfPoswiadczenAutomatyki
	okna         dane.RepozytoriumOkien

	// konta daje adres, na który idzie list o zakończonym przebiegu.
	konta dane.RepozytoriumKontaWlasciciela
	// nastawyListow oddaje konto nadawcze platformy; puste znaczy, że listów
	// nie ma czym wysłać i przebieg kończy się bez powiadomienia.
	nastawyListow func(context.Context) nadajnik.Nastawy
	// adresKonsoli to publiczny adres, pod który kieruje odsyłacz z listu.
	adresKonsoli string
}

// ZPocztaPrzebiegow wpina drogi listu o zakończonym przebiegu: konto adresata,
// konto nadawcze i adres Konsoli.
func (a *adapterAutomatyk) ZPocztaPrzebiegow(konta dane.RepozytoriumKontaWlasciciela,
	nastawy func(context.Context) nadajnik.Nastawy, adresKonsoli string) *adapterAutomatyk {

	a.konta, a.nastawyListow, a.adresKonsoli = konta, nastawy, adresKonsoli
	return a
}

func nowyAdapterAutomatyk(repozytorium dane.RepozytoriumAutomatyk) *adapterAutomatyk {
	return &adapterAutomatyk{
		repozytorium: repozytorium,
		obserwatorzy: nowaPamiecObserwatorowPrzebiegow(),
	}
}

func (a *adapterAutomatyk) ZKolejkami(kolejki *adapterKolejek) *adapterAutomatyk {
	a.kolejki = kolejki
	return a
}

// Jedna instancja sejfu pod jednym zamkiem: dwa sejfy nad tym samym plikiem ścigałyby się o zapis.
func (a *adapterAutomatyk) ZSejfem(sejf SejfPoswiadczenAutomatyki) *adapterAutomatyk {
	a.sejf = sejf
	return a
}

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

// Nieudany zapis migawki nie wywraca zapisu definicji: definicja już stoi w bazie.
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

func (a *adapterAutomatyk) zloz(ctx context.Context, wiersz dane.Automatyka) (shared.AutomationWorkflow, error) {
	kroki, err := a.krokiKontraktu(ctx, wiersz.ID)
	if err != nil {
		return shared.AutomationWorkflow{}, err
	}
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

// Brak wskazania znaczy automatykę czynną: Operator, który ją zapisał, chce nią pracować.
func czyCzynna(wskazanie *bool) bool {
	return wskazanie == nil || *wskazanie
}

func wartoscLiczby(wskazanie *int) int {
	if wskazanie == nil {
		return 0
	}
	return *wskazanie
}

func bladAutomatyki(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeConflict, err))
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

func bladWskazaniaAutomatyki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Automations: "+powod))
}

func bladNieznanejAutomatyki(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Automations: automatyka nie istnieje: "+kod))
	}
	return bladAutomatyki(err)
}
