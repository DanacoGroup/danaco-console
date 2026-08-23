package core

import (
	"strings"
	"testing"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Sprawdziany pracy wielu wykonawców nad jednym dokumentem.
//
// Szkody, które ten plik ma wykluczyć:
//  1. zajęcie WYGASŁE pokazywane jako czynne — wykonawca ubity w pół pracy
//     trzymałby wtedy akapit na zawsze, a drugi dostawałby odmowę bez powodu;
//  2. odmowa mówiąca „zajęte" i nic więcej — wykonawca nie wie, czy czekać,
//     czy odstąpić, a Operator nie wie, kto mu blokuje pismo;
//  3. nastawa spięcia spoza słownika, przez którą bilans spięcia przestaje
//     znaczyć to, co mówi.

// TestAgenciZajecieWygasleNieJestCzynne mierzy pierwszą szkodę.
func TestAgenciZajecieWygasleNieJestCzynne(t *testing.T) {
	kod := "agent-redaktor"
	wygasle := dane.ZajecieFragmentuStudia{
		Kod: "studio-zaj-wygasle", DokumentKod: "studio-dok-pierwszy",
		WykonawcaRodzaj: string(shared.StudioAuthorModel), AgentKod: &kod,
		Stan: string(shared.StudioAgentSlotStateWorking), ZakresOd: 10, ZakresDo: 40,
		Zajeto: "2026-08-17T09:00:00.000Z", Wygasle: true,
	}
	zlozone := agenciZlozZajecie(wygasle)
	if zlozone.State != shared.StudioAgentSlotStateIdle {
		t.Errorf("zajęcie wygasłe wyszło jako %q — wykaz twierdzi, że wykonawca wciąż "+
			"pracuje", zlozone.State)
	}

	czynne := wygasle
	czynne.Wygasle = false
	if agenciZlozZajecie(czynne).State != shared.StudioAgentSlotStateWorking {
		t.Errorf("zajęcie czynne wyszło jako bezczynne")
	}

	// Licznik pracujących nie liczy wygasłych.
	czynni := agenciCzynniWykonawcy([]dane.ZajecieFragmentuStudia{wygasle, czynne})
	if len(czynni) != 1 {
		t.Errorf("licznik pracujących wliczył zajęcie wygasłe: %d", len(czynni))
	}
}

// TestAgenciOdmowaNazywaWykonawce mierzy drugą szkodę.
func TestAgenciOdmowaNazywaWykonawce(t *testing.T) {
	kod, nazwa := "agent-redaktor", "Redaktor pisma"
	zNazwa := dane.ZajecieFragmentuStudia{AgentKod: &kod, AgentNazwa: &nazwa}
	if opis := agenciNazwaZajecia(zNazwa); !strings.Contains(opis, nazwa) ||
		!strings.Contains(opis, kod) {

		t.Errorf("nazwa trzymającego nie niesie ani nazwy, ani kodu: %q", opis)
	}
	// Nazwa widoczna stoi PRZED kodem: Operator czyta nazwę, nie identyfikator.
	opis := agenciNazwaZajecia(zNazwa)
	if strings.Index(opis, nazwa) > strings.Index(opis, kod) {
		t.Errorf("kod stanął przed nazwą widoczną: %q", opis)
	}
	bezNiczego := dane.ZajecieFragmentuStudia{}
	if opis := agenciNazwaZajecia(bezNiczego); opis == "" {
		t.Errorf("wykonawca nienazwany dał pustą nazwę — odmowa nie powiedziałaby nic")
	}
}

// TestAgenciZajecieTejSamejRekiNieZderzaSieZSoba mierzy, że wykonawca
// przedłużający własne zajęcie nie dostaje odmowy od siebie samego.
func TestAgenciZajecieTejSamejRekiNieZderzaSieZSoba(t *testing.T) {
	kod := "agent-redaktor"
	inny := "agent-korektor"
	zajecie := dane.ZajecieFragmentuStudia{AgentKod: &kod}
	if !agenciToSamRek(zajecie, kontrolaWykonawca{
		Rodzaj: shared.StudioAuthorModel, AgentKod: &kod,
	}) {
		t.Errorf("wykonawca nie rozpoznał własnego zajęcia")
	}
	if agenciToSamRek(zajecie, kontrolaWykonawca{
		Rodzaj: shared.StudioAuthorModel, AgentKod: &inny,
	}) {
		t.Errorf("cudze zajęcie zostało uznane za własne — dwóch wykonawców pisałoby " +
			"po tym samym akapicie")
	}
}

// TestAgenciNastawaSpieciaPozaSlownikiemOdmawia mierzy trzecią szkodę.
func TestAgenciNastawaSpieciaPozaSlownikiemOdmawia(t *testing.T) {
	if err := agenciSprawdzNastaweSpiecia("nadpisz"); err == nil {
		t.Errorf("nastawa spoza wyliczenia kontraktu przeszła")
	} else if !strings.Contains(err.Error(), "queue") {
		t.Errorf("odmowa nie mówi, co wolno: %v", err)
	}
	for _, nastawa := range shared.WartosciStudioAgentConflictPolicy() {
		if err := agenciSprawdzNastaweSpiecia(nastawa); err != nil {
			t.Errorf("nastawa %q z wyliczenia została odbita: %v", nastawa, err)
		}
	}
}

// TestAgenciSpiecieNiesieOdlozoneBrzmienie mierzy, że praca wykonawcy, którego
// zmiana została odłożona, NIE PRZEPADA — Operator może wnieść ją sam.
func TestAgenciSpiecieNiesieOdlozoneBrzmienie(t *testing.T) {
	odlozone := "Brzmienie, które nie weszło."
	weszla, odlozona := "agent-redaktor", "agent-korektor"
	spiecie := agenciZlozSpiecie(dane.SpiecieWykonawcowStudia{
		Kod: "studio-spc-pierwsze", DokumentKod: "studio-dok-pierwszy",
		ZakresOd: 10, ZakresDo: 40,
		WeszlaRodzaj: string(shared.StudioAuthorModel), WeszlaAgentKod: &weszla,
		OdlozonaRodzaj: string(shared.StudioAuthorModel), OdlozonaAgentKod: &odlozona,
		Nastawa:           shared.StudioAgentConflictPolicyQueue,
		Powod:             "fragment trzymał w tym czasie inny wykonawca",
		BrzmienieOdlozone: &odlozone,
	})
	if spiecie.DeferredText == nil || *spiecie.DeferredText != odlozone {
		t.Fatalf("odłożone brzmienie przepadło — praca wykonawcy zniknęła bez śladu")
	}
	if spiecie.AppliedActor.AgentId == nil || *spiecie.AppliedActor.AgentId != weszla {
		t.Errorf("spięcie nie mówi, czyja zmiana weszła")
	}
	if spiecie.DeferredActor.AgentId == nil || *spiecie.DeferredActor.AgentId != odlozona {
		t.Errorf("spięcie nie mówi, czyja zmiana została odłożona")
	}
	if spiecie.Reason == "" {
		t.Errorf("spięcie bez powodu jest ciszą, nie bilansem")
	}
}

// TestAgenciNastawyDomyslnieWylaczone pilnuje rozstrzygnięcia, że pętla
// wykonawcza i praca wielu agentów są DOMYŚLNIE WYŁĄCZONE.
func TestAgenciNastawyDomyslnieWylaczone(t *testing.T) {
	nastawy := shared.StudioAgentSettings{}
	if nastawy.ExecutionLoopEnabled || nastawy.MultiAgentEnabled {
		t.Errorf("nastawy w postaci zerowej są włączone — narzędzie ma być domyślnie " +
			"wyłączone i włączane jawnym ustawieniem Operatora")
	}
	// Brak granicy liczby wykonawców zostaje BRAKIEM, nie zerem: zero znaczyłoby
	// „ani jeden wykonawca".
	if nastawy.MaxConcurrentAgents != nil {
		t.Errorf("brak granicy wyszedł jako wartość: %+v", nastawy.MaxConcurrentAgents)
	}
}
