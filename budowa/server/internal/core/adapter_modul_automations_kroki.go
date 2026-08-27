// Odpowiedzialność pliku: przekład kroków automatyki między kontraktem a warstwą danych oraz utrzymanie układu zależności przy zapisie definicji w Workflow Builder.
package core

import (
	"encoding/json"
	"strconv"

	"context"
	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// zapiszKroki podmienia kroki automatyki i uzgadnia z nimi układ zależności, scalając zastane łuki z tymi wynikającymi z dependsOn.
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

// wierszKroku przekłada krok kontraktu na wiersz; krok bez identyfikatora dostaje go od rdzenia przy zapisie.
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
	// Odwołania do skarbca zapisują się wraz z krokiem: usunięcie sekretu poznaje utratę pokrycia.
	if len(krok.SecretRefs) > 0 {
		wiersz.OdwolaniaSekretow = zapisStrukturalny(krok.SecretRefs)
	}
	return wiersz
}

// scalZaleznosci składa układ po zapisie kroków: zastane łuki o obu końcach nadal istniejących plus łuki wynikające z dependsOn, zachowując rodzaj łuku zastanego.
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

// kluczLuku znakuje parę kroków, żeby ten sam łuk nie wszedł do układu zależności dwukrotnie w zapisie.
func kluczLuku(z, do string) string { return z + "\x00" + do }

// krokiKontraktu składa kroki automatyki wraz z polem dependsOn odtworzonym z układu zależności zapisanego w bazie.
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
	// Notatka kroku leży w adnotacjach kanwy, nie w wierszu: zapis podmienia wiersze w całości.
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

// krokKontraktu przekłada wiersz kroku automatyki na krok kontraktu wymiany z klientem Workflow Builder.
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
