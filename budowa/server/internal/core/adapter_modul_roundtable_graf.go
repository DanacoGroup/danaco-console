// Odpowiedzialność pliku: graf argumentów — `roundtable.argument.list`,
// `roundtable.argument.pin` i `roundtable.argument.export` (okno Argument Map
// & Analysis).
//
// Graf czyta się z bazy, a nie liczy przy odczycie. Powód stoi przy migracji
// 192: oznaczenie węzła jako kluczowego stawia Operator, a oznaczenie na węźle
// wyliczanym w locie znikałoby razem z jego identyfikatorem przy następnym
// otwarciu panelu.
package core

import (
	"context"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Graf oddaje graf argumentów okna albo wybranej tury.
func (a *adapterDebaty) Graf(ctx context.Context,
	z shared.RoundtableArgumentListRequest) (shared.RoundtableArgumentListResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableArgumentListResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	tylkoKluczowe := z.PinnedOnly != nil && *z.PinnedOnly
	graf, err := a.zlozGraf(ctx, okno, strings.TrimSpace(wartoscTekstu(z.TurnId)), tylkoKluczowe)
	if err != nil {
		return shared.RoundtableArgumentListResponse{}, err
	}
	return shared.RoundtableArgumentListResponse{Graph: graf}, nil
}

// zlozGraf czyta węzły, krawędzie i oznaczenia błędów, i składa z nich graf
// kontraktu wraz z rozpoznanym punktem spornym.
func (a *adapterDebaty) zlozGraf(ctx context.Context, okno, turaKod string,
	tylkoKluczowe bool) (shared.RoundtableArgumentGraph, error) {

	wezly, err := a.repozytorium.WezlyDebaty(ctx, okno, turaKod, tylkoKluczowe)
	if err != nil {
		return shared.RoundtableArgumentGraph{}, bladDebaty(err)
	}
	krawedzie, err := a.repozytorium.KrawedzieDebaty(ctx, okno)
	if err != nil {
		return shared.RoundtableArgumentGraph{}, bladDebaty(err)
	}
	oznaczenia, err := a.repozytorium.OznaczeniaBledowDebaty(ctx, okno)
	if err != nil {
		return shared.RoundtableArgumentGraph{}, bladDebaty(err)
	}

	wedlugWezla := make(map[string][]shared.RoundtableFallacyMark, len(oznaczenia))
	for _, oznaczenie := range oznaczenia {
		znacznik := shared.RoundtableFallacyMark{
			Id: oznaczenie.Kod, NodeId: oznaczenie.Wezel, Code: oznaczenie.KodBledu,
			Name: oznaczenie.Nazwa, Rationale: oznaczenie.Uzasadnienie,
		}
		if oznaczenie.Pewnosc >= 0 {
			pewnosc := oznaczenie.Pewnosc
			znacznik.Confidence = &pewnosc
		}
		wedlugWezla[oznaczenie.Wezel] = append(wedlugWezla[oznaczenie.Wezel], znacznik)
	}

	zywe := make(map[string]struct{}, len(wezly))
	tury := make([]string, 0, 4)
	widzianeTury := make(map[string]struct{}, 4)
	wykazWezlow := make([]shared.RoundtableArgumentNode, 0, len(wezly))
	for _, wezel := range wezly {
		zywe[wezel.Kod] = struct{}{}
		if wezel.Tura != "" {
			if _, jest := widzianeTury[wezel.Tura]; !jest {
				widzianeTury[wezel.Tura] = struct{}{}
				tury = append(tury, wezel.Tura)
			}
		}
		wykazWezlow = append(wykazWezlow, wezelKontraktu(wezel, wedlugWezla[wezel.Kod]))
	}

	wykazKrawedzi := make([]shared.RoundtableArgumentEdge, 0, len(krawedzie))
	for _, krawedz := range krawedzie {
		if _, zrodlo := zywe[krawedz.WezelZrodlowy]; !zrodlo {
			continue
		}
		if _, cel := zywe[krawedz.WezelWskazany]; !cel {
			continue
		}
		pozycja := shared.RoundtableArgumentEdge{
			Id: krawedz.Kod, FromNodeId: krawedz.WezelZrodlowy,
			ToNodeId: krawedz.WezelWskazany,
			Relation: shared.RoundtableArgumentRelation(krawedz.Relacja),
		}
		if krawedz.Pewnosc >= 0 {
			pewnosc := krawedz.Pewnosc
			pozycja.Confidence = &pewnosc
		}
		wykazKrawedzi = append(wykazKrawedzi, pozycja)
	}

	graf := shared.RoundtableArgumentGraph{
		WindowId: okno, TurnIds: tury, Nodes: wykazWezlow, Edges: wykazKrawedzi,
	}
	if punkt := punktSpornyGrafu(okno, wezly, krawedzie); punkt != nil {
		graf.Crux = punkt
	}
	return graf, nil
}

// Przypnij oznacza węzeł jako argument kluczowy albo zdejmuje oznaczenie.
func (a *adapterDebaty) Przypnij(ctx context.Context,
	z shared.RoundtableArgumentPinRequest) (shared.RoundtableArgumentPinResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kod := strings.TrimSpace(z.NodeId)
	if okno == "" {
		return shared.RoundtableArgumentPinResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kod == "" {
		return shared.RoundtableArgumentPinResponse{},
			bladWskazaniaDebaty("oznaczenie bez wskazania węzła grafu")
	}
	wezel, err := a.repozytorium.WezelDebatyPoKodzie(ctx, kod)
	if err != nil {
		return shared.RoundtableArgumentPinResponse{}, bladNieznanegoWezla(kod, err)
	}
	if wezel.Okno != okno {
		return shared.RoundtableArgumentPinResponse{},
			bladWskazaniaDebaty("węzeł " + kod + " nie należy do okna " + okno)
	}
	if err := a.repozytorium.OznaczWezelDebaty(ctx, kod, z.Pinned); err != nil {
		return shared.RoundtableArgumentPinResponse{}, bladNieznanegoWezla(kod, err)
	}
	po, err := a.repozytorium.WezelDebatyPoKodzie(ctx, kod)
	if err != nil {
		return shared.RoundtableArgumentPinResponse{}, bladNieznanegoWezla(kod, err)
	}
	return shared.RoundtableArgumentPinResponse{Node: wezelKontraktu(po, nil)}, nil
}

// punktSpornyGrafu wskazuje węzeł, którego rozstrzygnięcie zmienia najwięcej.
//
// Miarą jest liczba wymierzonych w węzeł podważeń wraz z poparciem, jakie sam
// zebrał: węzeł, który wielu popiera i wielu atakuje, jest osią sporu. Węzeł
// bez ani jednego podważenia punktem spornym nie jest — nikt się z nim nie
// spiera — więc graf bez podważeń oddaje brak, a nie pierwszy węzeł z brzegu.
func punktSpornyGrafu(okno string, wezly []dane.WezelDebaty,
	krawedzie []dane.KrawedzDebaty) *shared.RoundtableCrux {

	podwazenia := make(map[string]int, len(wezly))
	for _, krawedz := range krawedzie {
		if krawedz.Relacja == shared.RoundtableArgumentRelationAttacks {
			podwazenia[krawedz.WezelWskazany]++
		}
	}
	najlepszy, najlepszaSila := dane.WezelDebaty{}, 0
	for _, wezel := range wezly {
		sila := podwazenia[wezel.Kod]
		if sila == 0 {
			continue
		}
		sila += wezel.Poparcie
		if sila > najlepszaSila {
			najlepszy, najlepszaSila = wezel, sila
		}
	}
	if najlepszaSila == 0 {
		return nil
	}
	kod := najlepszy.Kod
	wplyw := float64(podwazenia[najlepszy.Kod]) / float64(najlepszaSila)
	return &shared.RoundtableCrux{
		Id: najlepszy.Kod, WindowId: okno, NodeId: &kod, Text: najlepszy.Tresc,
		Rationale: "Argument podważony " + itoa(podwazenia[najlepszy.Kod]) +
			" raz(y) przy poparciu " + itoa(najlepszy.Poparcie) +
			" uczestnika(-ów) — jego rozstrzygnięcie przestawia wnioski obu stron.",
		Impact: &wplyw,
	}
}

// wezelKontraktu przekłada węzeł grafu na byt kontraktu.
func wezelKontraktu(w dane.WezelDebaty,
	bledy []shared.RoundtableFallacyMark) shared.RoundtableArgumentNode {

	poparcie, kluczowy := w.Poparcie, w.Kluczowy
	wezel := shared.RoundtableArgumentNode{
		Id: w.Kod, WindowId: w.Okno, SpeechAct: shared.RoundtableSpeechAct(w.AktMowy),
		Text: w.Tresc, Support: &poparcie, Pinned: &kluczowy, Fallacies: bledy,
	}
	if w.Wypowiedz != "" {
		kod := w.Wypowiedz
		wezel.StatementId = &kod
	}
	if w.Uczestnik != "" {
		kod := w.Uczestnik
		wezel.ParticipantId = &kod
	}
	return wezel
}

// posortowaneKlucze oddaje klucze mapy w porządku ustalonym. Wyniki panelu mają
// być te same przy dwóch kolejnych odczytach, a przebieg po mapie w Go jest
// losowy z założenia.
func posortowaneKlucze[T any](mapa map[string]T) []string {
	klucze := make([]string, 0, len(mapa))
	for klucz := range mapa {
		klucze = append(klucze, klucz)
	}
	sort.Strings(klucze)
	return klucze
}
