package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// UruchomAnalize zestawia błędy z zakresu czasu w jedną migawkę i wyprowadza
// z nich rekomendacje.
//
// Każda rekomendacja wskazuje błąd, z którego powstała (kolumna `blad_kod`), bo
// Recommendations Panel pozwala przejść z zalecenia do faktu źródłowego.
//
// Zakres pusty jest stanem poprawnym: analiza powstaje z zerem błędów i mówi to
// w podsumowaniu. Odmowa zlałaby brak znalezisk z niepowodzeniem sprawdzenia.
func (a *adapterDiagnostyki) UruchomAnalize(ctx context.Context,
	z shared.DiagnosticsAnalyzeRunRequest) (shared.DiagnosticsAnalyzeRunResponse, error) {

	if a.repozytorium == nil {
		return shared.DiagnosticsAnalyzeRunResponse{}, bladBrakuTrwalosci("analiza diagnostyczna")
	}
	bledy, err := a.repozytorium.Bledy(ctx, dane.FiltrBledow{
		Od: wartoscChwili(z.FromTime),
		Do: wartoscChwili(z.ToTime),
	})
	if err != nil {
		return shared.DiagnosticsAnalyzeRunResponse{}, bladDiagnostyki(err)
	}

	analiza := a.migawkaAnalizy(z, bledy)
	rekomendacje := rekomendacjeZBledow(analiza.Kod, bledy)
	if err := a.repozytorium.ZapiszAnalize(ctx, analiza, rekomendacje); err != nil {
		return shared.DiagnosticsAnalyzeRunResponse{}, bladDiagnostyki(err)
	}

	wynik := analizaKontraktu(analiza, kodyRekomendacji(rekomendacje))
	if a.zmiana != nil {
		a.zmiana(shared.ChangeKindCreated, wynik)
	}
	return shared.DiagnosticsAnalyzeRunResponse{Analysis: wynik}, nil
}

// WykazRekomendacji zwraca rekomendacje spełniające warunki.
func (a *adapterDiagnostyki) WykazRekomendacji(ctx context.Context,
	z shared.DiagnosticsRecommendationListRequest) (shared.DiagnosticsRecommendationListResponse, error) {

	if a.repozytorium == nil {
		return shared.DiagnosticsRecommendationListResponse{}, bladBrakuTrwalosci("wykaz rekomendacji")
	}
	wiersze, err := a.repozytorium.Rekomendacje(ctx, dane.FiltrRekomendacji{
		AnalizaKod: wartoscTekstu(z.AnalysisId),
		Stan:       wartoscStanuRekomendacji(z.Status),
		Priorytet:  wartoscPriorytetu(z.Priority),
		Granica:    wartoscLiczby(z.Limit),
	})
	if err != nil {
		return shared.DiagnosticsRecommendationListResponse{}, bladDiagnostyki(err)
	}
	return shared.DiagnosticsRecommendationListResponse{
		Recommendations: rekomendacjeKontraktu(wiersze),
	}, nil
}

// migawkaAnalizy składa wiersz analizy z zakresu żądania i znalezionych błędów.
//
// Kody błędów idą do kolumny tablicą JSON, bo migawka ma przetrwać zmianę stanu
// samych błędów: analiza sprzed tygodnia ma pokazywać to, co widziała wtedy,
// a nie to, co widać dziś.
func (a *adapterDiagnostyki) migawkaAnalizy(z shared.DiagnosticsAnalyzeRunRequest,
	bledy []dane.BladDiagnostyczny) dane.AnalizaDiagnostyczna {

	kody := make([]string, 0, len(bledy))
	for _, blad := range bledy {
		kody = append(kody, blad.Kod)
	}
	zapis, err := json.Marshal(kody)
	if err != nil {
		zapis = []byte("[]")
	}
	tresc := string(zapis)
	podsumowanie := podsumowanieAnalizy(bledy, a.odrzucone.Load(), a.niezapisane.Load())

	return dane.AnalizaDiagnostyczna{
		Kod:          nowyIdentyfikator(przedrostekAnalizy),
		OknoKod:      z.WindowId,
		ZakresOd:     z.FromTime,
		ZakresDo:     z.ToTime,
		Podsumowanie: &podsumowanie,
		Bledy:        &tresc,
		Utworzono:    time.Now().UnixMilli(),
	}
}

// podsumowanieAnalizy opisuje materiał, na którym analiza stanęła.
//
// Wpisy odrzucone przez kolejkę i te, których nie przyjęła baza, znaczą materiał
// niekompletny, więc podsumowanie wymienia je wprost.
func podsumowanieAnalizy(bledy []dane.BladDiagnostyczny, odrzucone, niezapisane int64) string {
	otwarte := 0
	for _, blad := range bledy {
		if blad.Stan == shared.DiagnosticErrorStatusNew {
			otwarte++
		}
	}
	podsumowanie := fmt.Sprintf("błędów w zakresie: %d, w tym otwartych: %d", len(bledy), otwarte)
	if odrzucone > 0 || niezapisane > 0 {
		podsumowanie += fmt.Sprintf(
			"; materiał niekompletny — wpisów odrzuconych: %d, niezapisanych: %d",
			odrzucone, niezapisane)
	}
	return podsumowanie
}

// rekomendacjeZBledow wyprowadza zalecenia z błędów objętych analizą.
//
// Jedna rekomendacja na błąd, z zachowaniem jego priorytetu i wskazaniem jego
// kodu. Rdzeń podaje sam fakt i jego wagę, bez treści poprawki: rozpoznanie
// przyczyny odbywa się w oknie Diagnostics.
func rekomendacjeZBledow(kodAnalizy string, bledy []dane.BladDiagnostyczny) []dane.RekomendacjaDiagnostyczna {
	rekomendacje := make([]dane.RekomendacjaDiagnostyczna, 0, len(bledy))
	teraz := time.Now().UnixMilli()
	for _, blad := range bledy {
		kod := blad.Kod
		szczegol := fmt.Sprintf("wystąpień: %d, kod błędu: %s", blad.Wystapienia, blad.KodBledu)
		rekomendacje = append(rekomendacje, dane.RekomendacjaDiagnostyczna{
			Kod:        nowyIdentyfikator(przedrostekRekomendacji),
			AnalizaKod: kodAnalizy,
			BladKod:    &kod,
			Tytul:      blad.Tresc,
			Szczegol:   &szczegol,
			Priorytet:  blad.Priorytet,
			Stan:       shared.RecommendationStatusProposed,
			Utworzono:  teraz,
		})
	}
	return rekomendacje
}

// kodyRekomendacji wyjmuje kody do migawki analizy.
func kodyRekomendacji(rekomendacje []dane.RekomendacjaDiagnostyczna) []string {
	kody := make([]string, 0, len(rekomendacje))
	for _, rekomendacja := range rekomendacje {
		kody = append(kody, rekomendacja.Kod)
	}
	return kody
}
