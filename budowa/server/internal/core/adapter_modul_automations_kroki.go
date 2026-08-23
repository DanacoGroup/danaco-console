// Odpowiedzialność pliku: przekład kroków automatyki między kontraktem
// a warstwą danych oraz utrzymanie układu zależności przy zapisie definicji
// (okno Workflow Builder).
//
// Zależność ma jedno źródło. Kontrakt niesie ten sam łuk grafu dwa razy:
// polem `AutomationStep.dependsOn` i strukturą `AutomationDependency` komendy
// orkiestratora. Prawdą jest tabela `zaleznosc_kroku_automatyki`, a `dependsOn`
// powstaje z niej przy odczycie — drugiego zapisu tej samej rzeczy nie ma.
//
// Zapis definicji nie kasuje pracy orkiestratora. Workflow Builder oddaje kroki,
// nie układ. Gdyby zapis podmieniał zależności na same `dependsOn`, ręcznie
// ustawione łuki równoległe i warunkowe znikałyby przy każdym zapisie nazwy
// kroku. Zapis scala więc dwa zbiory: zastane łuki, których oba końce nadal
// istnieją, oraz łuki wynikające z `dependsOn`.
package core

import (
	"encoding/json"
	"strconv"

	"context"
	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// zapiszKroki podmienia kroki automatyki i uzgadnia z nimi układ zależności.
func (a *adapterAutomatyk) zapiszKroki(ctx context.Context, automatykaID int64,
	kroki []shared.AutomationStep) error {

	zastane, err := a.repozytorium.Zaleznosci(ctx, automatykaID)
	if err != nil {
		return err
	}
	wiersze := make([]dane.KrokAutomatyki, 0, len(kroki))
	for numer, krok := range kroki {
		wiersze = append(wiersze, wierszKroku(krok, numer))
	}
	if err := a.repozytorium.ZapiszKroki(ctx, automatykaID, wiersze); err != nil {
		return err
	}
	return a.repozytorium.ZapiszZaleznosci(ctx, automatykaID,
		scalZaleznosci(zastane, kroki, wiersze))
}

// wierszKroku przekłada krok kontraktu na wiersz. Krok bez identyfikatora
// dostaje go od rdzenia — Workflow Builder pozwala dołożyć krok, zanim Operator
// go nazwie, a bez identyfikatora nie dałoby się ustalić zależności.
func wierszKroku(krok shared.AutomationStep, numer int) dane.KrokAutomatyki {
	kod := krok.Id
	if kod == "" {
		kod = przedrostekKroku + strconv.Itoa(numer+1)
	}
	rodzaj := string(krok.Kind)
	if rodzaj == "" {
		rodzaj = shared.AutomationStepKindCommand
	}
	kolejnosc := numer + 1
	if krok.Order != nil && *krok.Order > 0 {
		kolejnosc = *krok.Order
	}
	wiersz := dane.KrokAutomatyki{
		Kod: kod, Nazwa: krok.Name, Rodzaj: rodzaj, Komenda: krok.Command,
		Warunek: krok.Condition, Kolejnosc: kolejnosc,
	}
	if len(krok.Params) > 0 {
		parametry := string(krok.Params)
		wiersz.Parametry = &parametry
	}
	// Odwołania do skarbca zapisują się wraz z krokiem. Bez nich
	// `automation.secret.remove` nie miałby jak nazwać kroków, które straciły
	// pokrycie — skarbiec zdejmowałby klucz w ciszy, a Operator dowiadywałby
	// się o skutku dopiero z nieudanego przebiegu.
	if len(krok.SecretRefs) > 0 {
		wiersz.OdwolaniaSekretow = zapisStrukturalny(krok.SecretRefs)
	}
	return wiersz
}

// scalZaleznosci składa układ po zapisie kroków: zastane łuki o obu końcach
// nadal istniejących plus łuki wynikające z `dependsOn`. Łuk zastany zachowuje
// swój rodzaj — sekwencyjność wynikająca z `dependsOn` nie nadpisuje ustawienia
// równoległego ani warunkowego wprowadzonego w Orchestratorze.
func scalZaleznosci(zastane []dane.ZaleznoscKroku, kroki []shared.AutomationStep,
	wiersze []dane.KrokAutomatyki) []dane.ZaleznoscKroku {

	istniejace := make(map[string]bool, len(wiersze))
	for _, wiersz := range wiersze {
		istniejace[wiersz.Kod] = true
	}
	wynik := make([]dane.ZaleznoscKroku, 0, len(zastane)+len(kroki))
	zajete := map[string]bool{}
	for _, zaleznosc := range zastane {
		if !istniejace[zaleznosc.KrokZ] || !istniejace[zaleznosc.KrokDo] {
			continue
		}
		wynik = append(wynik, zaleznosc)
		zajete[kluczLuku(zaleznosc.KrokZ, zaleznosc.KrokDo)] = true
	}
	for numer, krok := range kroki {
		docelowy := wiersze[numer].Kod
		for _, poprzednik := range krok.DependsOn {
			klucz := kluczLuku(poprzednik, docelowy)
			if poprzednik == docelowy || zajete[klucz] || !istniejace[poprzednik] {
				continue
			}
			wynik = append(wynik, dane.ZaleznoscKroku{
				KrokZ: poprzednik, KrokDo: docelowy,
				Rodzaj: shared.AutomationDependencyKindSequential,
			})
			zajete[klucz] = true
		}
	}
	return wynik
}

// kluczLuku znakuje parę kroków, żeby ten sam łuk nie wszedł do układu dwa razy.
func kluczLuku(z, do string) string { return z + "\x00" + do }

// krokiKontraktu składa kroki automatyki wraz z polem `dependsOn` odtworzonym
// z układu zależności.
func (a *adapterAutomatyk) krokiKontraktu(ctx context.Context, automatykaID int64) ([]shared.AutomationStep, error) {
	wiersze, err := a.repozytorium.Kroki(ctx, automatykaID)
	if err != nil {
		return nil, err
	}
	zaleznosci, err := a.repozytorium.Zaleznosci(ctx, automatykaID)
	if err != nil {
		return nil, err
	}
	poprzednicy := map[string][]string{}
	for _, zaleznosc := range zaleznosci {
		poprzednicy[zaleznosc.KrokDo] = append(poprzednicy[zaleznosc.KrokDo], zaleznosc.KrokZ)
	}
	// Notatka kroku leży w adnotacjach kanwy, nie w wierszu kroku: zapis
	// definicji podmienia wiersze kroków w całości, więc kolumna kasowałaby się
	// przy każdej zmianie nazwy. Nieudany odczyt adnotacji zostawia kroki bez
	// notatek zamiast wywracać odczyt definicji.
	notatki := map[string]*string{}
	if adnotacje, err := a.repozytorium.AdnotacjeKrokow(ctx, automatykaID); err == nil {
		for _, adnotacja := range adnotacje {
			notatki[adnotacja.KrokKod] = adnotacja.Notatka
		}
	}
	kroki := make([]shared.AutomationStep, 0, len(wiersze))
	for _, wiersz := range wiersze {
		krok := krokKontraktu(wiersz, poprzednicy[wiersz.Kod])
		krok.Note = notatki[wiersz.Kod]
		kroki = append(kroki, krok)
	}
	return kroki, nil
}

// krokKontraktu przekłada wiersz kroku na krok kontraktu.
func krokKontraktu(wiersz dane.KrokAutomatyki, poprzednicy []string) shared.AutomationStep {
	kolejnosc := wiersz.Kolejnosc
	krok := shared.AutomationStep{
		Id: wiersz.Kod, Name: wiersz.Nazwa, Kind: shared.AutomationStepKind(wiersz.Rodzaj),
		Command: wiersz.Komenda, DependsOn: poprzednicy, Condition: wiersz.Warunek,
		Order: &kolejnosc,
	}
	if wiersz.OdwolaniaSekretow != nil && *wiersz.OdwolaniaSekretow != "" {
		krok.SecretRefs = kodyKrokowZZapisu(*wiersz.OdwolaniaSekretow)
	}
	if wiersz.Parametry != nil && *wiersz.Parametry != "" {
		krok.Params = json.RawMessage(*wiersz.Parametry)
	}
	return krok
}
