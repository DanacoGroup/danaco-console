// Odpowiedzialność pliku: wykaz podagentów (`subagent.list`) i zbieranie ich
// wyników (`subagent.result.collect`) — dwa odczyty tego samego wiersza, dwa
// różne pytania panelu.
//
// Wykaz pusty jest poprawną odpowiedzią, nie odmową: okno, które nikogo nie
// powołało, dostaje wykaz pusty i pusty panel. Dlatego zawężenia są warunkami
// zapytania, a nie warunkami wstępnymi żądania.
//
// Pole `waitForAll` oznacza czekanie rzeczywiste: odpowiedź pyta bazę w takcie,
// aż wszyscy objęci zbieraniem wejdą w stan końcowy albo aż zerwie się kontekst
// żądania. Odpowiedź natychmiastowa z `complete=false` pomijałaby to pole,
// a górnego limitu czekania adapter nie narzuca — czekanie kończy zamknięcie
// żądania przez klienta.
//
// Domyślnie rodzic zbiera po drodze. Podagenta powołuje model w trakcie tury
// rodzica, więc odpowiedź domyślnie blokująca zatrzymywałaby turę orkiestratora
// na cudzej pracy i unieważniała sens tła. Czekanie na wszystkich jest wyborem
// jawnym polem `waitForAll`. Czekania na pierwszego kontrakt nie zna (nie ma
// pola `waitForAny`) i ten adapter go nie wymyśla; rodzic osiąga to samo,
// zbierając po drodze i czytając pole `complete` oraz stany pozycji.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// taktZbierania wyznacza, jak często zbieranie pyta bazę o stan podagentów.
// Ćwierć sekundy jest niżej od czasu tury modelu i wyżej od kosztu jednego
// odczytu wykazu — częstsze pytanie nie zastałoby nowego stanu.
const taktZbierania = 250 * time.Millisecond

// stanyKoncowePodagenta wylicza stany, po których podagent nie wraca do pracy
// sam. Powrót jest możliwy wyłącznie biegiem naprawczym pozycji, wywołanym
// przez Operatora — ten nie ma limitu obiegów.
var stanyKoncowePodagenta = map[string]struct{}{
	dane.StanPodagentaUkonczony: {}, dane.StanPodagentaBledny: {},
	dane.StanPodagentaZatrzymany: {},
}

// Wykaz oddaje podagentów okna wykonawcy albo karty sesji.
//
// Oba wskazania idą do zapytania jednocześnie: podane naraz zawężają wykaz
// podwójnie, co nie jest sprzecznością, tylko węższym pytaniem. Kontrakt
// opisuje `sessionId` jako kartę sesji braną pod uwagę, gdy okna nie wskazano.
func (a *adapterPodagentow) Wykaz(ctx context.Context,
	z shared.SubagentListRequest) (shared.SubagentListResponse, error) {

	filtr := dane.FiltrPodagentow{
		OknoKod:  wartoscTekstu(z.WindowId),
		SesjaKod: wartoscTekstu(z.SessionId),
	}
	if z.Status != nil {
		filtr.Stan = string(*z.Status)
	}
	wiersze, err := a.repozytorium.Podagenci(ctx, filtr)
	if err != nil {
		return shared.SubagentListResponse{}, bladPodagentow(err)
	}
	return shared.SubagentListResponse{Subagents: podagenciKontraktu(wiersze)}, nil
}

// ZbierzWyniki oddaje podagentów wraz z wynikami i orzeka, czy wszyscy objęci
// zbieraniem skończyli pracę.
func (a *adapterPodagentow) ZbierzWyniki(ctx context.Context,
	z shared.SubagentResultCollectRequest) (shared.SubagentResultCollectResponse, error) {

	wiersze, err := a.objeciZbieraniem(ctx, z)
	if err != nil {
		return shared.SubagentResultCollectResponse{}, err
	}
	if z.WaitForAll != nil && *z.WaitForAll {
		wiersze, err = a.poczekajNaKomplet(ctx, z, wiersze)
		if err != nil {
			return shared.SubagentResultCollectResponse{}, err
		}
	}
	return shared.SubagentResultCollectResponse{
		Subagents: podagenciKontraktu(wiersze),
		Complete:  czyKomplet(wiersze),
	}, nil
}

// objeciZbieraniem dobiera podagentów, których wyniki są zbierane.
//
// Pusta lista wskazań znaczy komplet — tak mówi kontrakt. Komplet zawęża się
// wtedy oknem, jeśli okno wskazano; żądanie bez okna i bez wskazania podagentów
// zbiera wszystkich.
func (a *adapterPodagentow) objeciZbieraniem(ctx context.Context,
	z shared.SubagentResultCollectRequest) ([]dane.Podagent, error) {

	kody := niepusteKody(z.SubagentIds)
	if len(kody) == 0 {
		wiersze, err := a.repozytorium.Podagenci(ctx,
			dane.FiltrPodagentow{OknoKod: wartoscTekstu(z.WindowId)})
		if err != nil {
			return nil, bladPodagentow(err)
		}
		return wiersze, nil
	}
	wiersze, err := a.repozytorium.PodagenciPoKodach(ctx, kody)
	if err != nil {
		return nil, bladPodagentow(err)
	}
	// Wskazanie, któremu nie odpowiada żaden wiersz, jest pomyłką co do bytu,
	// nie pustym wynikiem: pusty wykaz czyta się jako „ci podagenci nic nie
	// oddali", a nie jako „takich podagentów nie ma".
	if len(wiersze) == 0 {
		return nil, bladNieznanychPodagentow(kody)
	}
	return wiersze, nil
}

// poczekajNaKomplet pyta bazę w takcie, aż wszyscy objęci zbieraniem wejdą
// w stan końcowy. Zerwany kontekst żądania kończy czekanie i oddaje stan
// zastany, zamiast błędu.
func (a *adapterPodagentow) poczekajNaKomplet(ctx context.Context,
	z shared.SubagentResultCollectRequest, zastane []dane.Podagent) ([]dane.Podagent, error) {

	zegar := time.NewTicker(taktZbierania)
	defer zegar.Stop()
	wiersze := zastane
	for !czyKomplet(wiersze) {
		select {
		case <-ctx.Done():
			return wiersze, nil
		case <-zegar.C:
			swiezsze, err := a.objeciZbieraniem(ctx, z)
			if err != nil {
				return nil, err
			}
			wiersze = swiezsze
		}
	}
	return wiersze, nil
}

// czyKomplet orzeka, czy wszyscy podagenci wykazu skończyli pracę. Wykaz pusty
// jest kompletny — nie ma na kogo czekać.
func czyKomplet(wiersze []dane.Podagent) bool {
	for _, podagent := range wiersze {
		if _, koncowy := stanyKoncowePodagenta[podagent.Stan]; !koncowy {
			return false
		}
	}
	return true
}

// niepusteKody odsiewa wskazania puste — puste wskazanie nie zawęża zbierania.
func niepusteKody(kody []string) []string {
	wybrane := make([]string, 0, len(kody))
	for _, kod := range kody {
		if przyciety := strings.TrimSpace(kod); przyciety != "" {
			wybrane = append(wybrane, przyciety)
		}
	}
	return wybrane
}
