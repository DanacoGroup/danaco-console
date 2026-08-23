// Odpowiedzialność pliku: przekład między wierszami obszaru Diagnostics
// a bytami kontraktu oraz zdejmowanie wskaźników z pól opcjonalnych żądań.
//
// Pole, którego wiersz nie niesie, wychodzi puste i nie dostaje wartości
// zastępczej: w oknie diagnostycznym zero znaczy co innego niż brak danych.
package core

import (
	"encoding/json"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// wpisyDziennikaKontraktu przekłada wiersze dziennika na wpisy Logs Viewer.
func wpisyDziennikaKontraktu(wiersze []dane.WpisDiagnostyki) []shared.LogEntry {
	wpisy := make([]shared.LogEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wpis := shared.LogEntry{
			Id: wiersz.Kod, Timestamp: wiersz.Chwila, Level: wiersz.Poziom,
			Source: wiersz.Zrodlo, Message: wiersz.Tresc,
			SessionId: wiersz.SesjaKod, WindowId: wiersz.OknoKod, ProcessId: wiersz.ProcesKod,
		}
		// Licznik wychodzi tylko wtedy, gdy naprawdę doszło do scalenia —
		// „1" przy każdym wpisie sugerowałoby deduplikację, której nie było.
		if wiersz.Powtorzenia > 1 {
			wpis.RepeatCount = wskaznikLiczby(wiersz.Powtorzenia)
		}
		wpisy = append(wpisy, wpis)
	}
	return wpisy
}

// bledyKontraktu przekłada wiersze błędów na byty Errors Panel.
func bledyKontraktu(wiersze []dane.BladDiagnostyczny) []shared.DiagnosticError {
	bledy := make([]shared.DiagnosticError, 0, len(wiersze))
	for _, wiersz := range wiersze {
		bledy = append(bledy, bladKontraktu(wiersz))
	}
	return bledy
}

// bladKontraktu przekłada jeden wiersz błędu.
func bladKontraktu(wiersz dane.BladDiagnostyczny) shared.DiagnosticError {
	blad := shared.DiagnosticError{
		Id: wiersz.Kod, Fingerprint: wiersz.Odcisk, Message: wiersz.Tresc,
		Source: wiersz.Zrodlo, ErrorCode: wskaznikNiepusty(string(wiersz.KodBledu)),
		Status: wiersz.Stan, Priority: wiersz.Priorytet,
		Occurrences: wskaznikLiczby(wiersz.Wystapienia), Note: wiersz.Notatka,
		FirstSeenAt: wiersz.Pierwsze, LastSeenAt: wiersz.Ostatnie,
	}
	if wiersz.Kontekst != nil && *wiersz.Kontekst != "" {
		blad.Context = json.RawMessage(*wiersz.Kontekst)
	}
	// CommitId i DeploymentId zostają puste: rdzeń nie wiąże odmowy komendy
	// z zatwierdzeniem repozytorium ani z wdrożeniem.
	return blad
}

// analizaKontraktu przekłada migawkę wraz z kodami rekomendacji z niej powstałych.
func analizaKontraktu(wiersz dane.AnalizaDiagnostyczna, rekomendacje []string) shared.DiagnosticAnalysis {
	return shared.DiagnosticAnalysis{
		Id: wiersz.Kod, WindowId: wiersz.OknoKod,
		FromTime: wiersz.ZakresOd, ToTime: wiersz.ZakresDo,
		Summary:  wiersz.Podsumowanie,
		ErrorIds: kodyBledowAnalizy(wiersz.Bledy), RecommendationIds: rekomendacje,
		ComparedAnalysisId: wiersz.PorownanaKod, CreatedAt: wiersz.Utworzono,
	}
}

// rekomendacjeKontraktu przekłada wiersze rekomendacji.
func rekomendacjeKontraktu(wiersze []dane.RekomendacjaDiagnostyczna) []shared.DiagnosticRecommendation {
	rekomendacje := make([]shared.DiagnosticRecommendation, 0, len(wiersze))
	for _, wiersz := range wiersze {
		rekomendacje = append(rekomendacje, shared.DiagnosticRecommendation{
			Id: wiersz.Kod, AnalysisId: wiersz.AnalizaKod, Title: wiersz.Tytul,
			Detail: wiersz.Szczegol, Priority: wiersz.Priorytet, Status: wiersz.Stan,
			TargetPath: wiersz.Sciezka, Patch: wiersz.Poprawka, CreatedAt: wiersz.Utworzono,
		})
	}
	return rekomendacje
}

// kodyBledowAnalizy odczytuje tablicę kodów zapisaną w kolumnie migawki.
// Kolumna nieczytelna daje wykaz pusty — migawka ma się pokazać, a nie zniknąć
// z powodu jednej kolumny.
func kodyBledowAnalizy(zapis *string) []string {
	if zapis == nil || *zapis == "" {
		return nil
	}
	var kody []string
	if err := json.Unmarshal([]byte(*zapis), &kody); err != nil {
		return nil
	}
	return kody
}

// wartoscChwili zdejmuje wskaźnik z granicy czasu; brak znaczy „bez granicy”.
func wartoscChwili(wskazanie *int64) int64 {
	if wskazanie == nil {
		return 0
	}
	return *wskazanie
}

// wartoscPrawdy zdejmuje wskaźnik z przełącznika; brak znaczy „wyłączony”.
func wartoscPrawdy(wskazanie *bool) bool {
	return wskazanie != nil && *wskazanie
}

// wartoscPoziomu zdejmuje wskaźnik z poziomu wpisu; brak nie zawęża wyniku.
func wartoscPoziomu(wskazanie *shared.LogLevel) shared.LogLevel {
	if wskazanie == nil {
		return ""
	}
	return *wskazanie
}

// wartoscStanuBledu zdejmuje wskaźnik ze stanu błędu; brak nie zawęża wyniku.
func wartoscStanuBledu(wskazanie *shared.DiagnosticErrorStatus) shared.DiagnosticErrorStatus {
	if wskazanie == nil {
		return ""
	}
	return *wskazanie
}

// wartoscPriorytetu zdejmuje wskaźnik z priorytetu; brak nie zawęża wyniku.
func wartoscPriorytetu(wskazanie *shared.DiagnosticPriority) shared.DiagnosticPriority {
	if wskazanie == nil {
		return ""
	}
	return *wskazanie
}

// wartoscStanuRekomendacji zdejmuje wskaźnik ze stanu rekomendacji.
func wartoscStanuRekomendacji(wskazanie *shared.RecommendationStatus) shared.RecommendationStatus {
	if wskazanie == nil {
		return ""
	}
	return *wskazanie
}

// wskaznikLiczby zakłada wskaźnik na liczbę — pole opcjonalne kontraktu niesie
// wtedy wartość wprost, a nie brak wartości.
func wskaznikLiczby(wartosc int) *int {
	kopia := wartosc
	return &kopia
}

// wskaznikPrawdy zakłada wskaźnik na rozstrzygnięcie logiczne.
func wskaznikPrawdy(wartosc bool) *bool {
	kopia := wartosc
	return &kopia
}

// wskaznikNiepusty zakłada wskaźnik wyłącznie dla napisu niepustego; pusty
// zostaje brakiem wartości, a nie napisem zerowej długości.
func wskaznikNiepusty(napis string) *string {
	if napis == "" {
		return nil
	}
	kopia := napis
	return &kopia
}
