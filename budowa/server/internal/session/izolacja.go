// Egzekutor izolacji obejmuje jedenaście punktów izolacji w postaci
// wykonawczej oraz wspólne pojęcia: naruszenie i przynależność ścieżki do
// obszaru.
package session

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"

	"danacoconsole/server/internal/konfig"
)

// ErrIzolacja jest wspólnym korzeniem naruszeń izolacji — warstwa wyżej
// rozpoznaje przyczynę przez errors.Is, nie przez treść komunikatu.
var ErrIzolacja = errors.New("session: naruszenie izolacji")

// naruszenie mówi, który punkt izolacji został naruszony i czym dokładnie,
// treścią gotową dla Operatora.
type naruszenie struct {
	// Klucz punktu izolacji ze stałych pakietu konfig.
	Klucz string
	// Powod niesie treść naruszenia gotową dla Operatora.
	Powod string
}

// Error składa czytelny komunikat tego naruszenia z klucza punktu izolacji
// oraz z powodu tego naruszenia.
func (n naruszenie) Error() string {
	return "session: izolacja " + n.Klucz + ": " + n.Powod
}

// Unwrap wiąże to naruszenie ze wspólnym korzeniem błędów izolacji całego
// tego egzekutora tej platformy.
func (n naruszenie) Unwrap() error { return ErrIzolacja }

// NoweNaruszenie składa naruszenie punktu izolacji. Egzekutor plików, sieci
// i kontekstu żyje w rdzeniu, a naruszenie ma jedną definicję.
func NoweNaruszenie(klucz, powod string) error {
	return naruszenie{Klucz: klucz, Powod: powod}
}

// Zasady to jedenaście punktów izolacji sprowadzonych do postaci wykonawczej:
// trzy wymiary kontekstu (prawda znaczy „odrębny") i osiem zakresów technicznych
// (prawda znaczy „włączony").
type Zasady struct {
	// HistoriaOdrebna — zapis wymiany wiadomości nie schodzi się między oknami.
	HistoriaOdrebna bool
	// PamiecOdrebna — pamięć długoterminowa jednego zasięgu niewidoczna w innym.
	PamiecOdrebna bool
	// KontekstOdrebny — bieżący stan roboczy nie przenosi się między oknami.
	KontekstOdrebny bool

	// KatalogRoboczy — proces okna pracuje we własnym katalogu.
	KatalogRoboczy bool
	// SrodowiskoProcesu — proces okna dostaje własny zestaw zmiennych.
	SrodowiskoProcesu bool
	// KatalogDanychModelu — kanał modelu ma własny katalog danych.
	KatalogDanychModelu bool
	// DostepSieciowy — połączenia wychodzące ograniczone do jawnie dozwolonych.
	DostepSieciowy bool
	// Pliki — dostęp do plików wyłącznie w ścieżkach jawnie dozwolonych.
	Pliki bool
	// KontoIToken — okno korzysta z własnych danych dostępowych.
	KontoIToken bool
	// ModelProcesu — okno ma własną instancję procesu wykonawczego.
	ModelProcesu bool
	// SerwerWykonania — okno korzysta z dedykowanego serwera wykonania.
	SerwerWykonania bool
}

// ZasadyZPolityki czyta jedenaście punktów izolacji z polityki efektywnej.
// Klucz nierozstrzygnięty schodzi na stan wyjściowy platformy: kontekst
// odrębny, żaden zakres techniczny niewłączony.
func ZasadyZPolityki(polityka konfig.Polityka) Zasady {
	return Zasady{
		HistoriaOdrebna:     odrebna(polityka, konfig.KluczIzolacjaHistoria),
		PamiecOdrebna:       odrebna(polityka, konfig.KluczIzolacjaPamiec),
		KontekstOdrebny:     odrebna(polityka, konfig.KluczIzolacjaKontekst),
		KatalogRoboczy:      wlaczony(polityka, konfig.KluczIzolacjaKatalogRoboczy),
		SrodowiskoProcesu:   wlaczony(polityka, konfig.KluczIzolacjaSrodowiskoProcesu),
		KatalogDanychModelu: wlaczony(polityka, konfig.KluczIzolacjaKatalogDanych),
		DostepSieciowy:      wlaczony(polityka, konfig.KluczIzolacjaDostepSieciowy),
		Pliki:               wlaczony(polityka, konfig.KluczIzolacjaPliki),
		KontoIToken:         wlaczony(polityka, konfig.KluczIzolacjaKontoIToken),
		ModelProcesu:        wlaczony(polityka, konfig.KluczIzolacjaModelProcesu),
		SerwerWykonania:     wlaczony(polityka, konfig.KluczIzolacjaSerwerWykonania),
	}
}

// wlaczony odpowiada, czy zakres techniczny został włączony. Wartość spoza
// słownika nie włącza zakresu — stanem wyjściowym jest wyłączenie.
func wlaczony(polityka konfig.Polityka, klucz string) bool {
	pozycja, jest := polityka.Pozycja(klucz)
	return jest && strings.TrimSpace(pozycja.Wartosc) == konfig.IzolacjaWlaczona
}

// odrebna odpowiada, czy wymiar kontekstu pozostaje odrębny. Współdzielenie
// jest decyzją zapisaną wprost; każda inna wartość zostawia wymiar odrębny.
func odrebna(polityka konfig.Polityka, klucz string) bool {
	pozycja, jest := polityka.Pozycja(klucz)
	if !jest {
		return true
	}
	return strings.TrimSpace(pozycja.Wartosc) != konfig.IzolacjaWspoldzielona
}

// SciezkaWewnatrz odpowiada, czy ścieżka leży w korzeniu albo jest samym
// korzeniem, porównaniem po członach ścieżki.
func SciezkaWewnatrz(korzen, sciezka string) bool {
	korzenNorm := normalizujSciezke(korzen)
	sciezkaNorm := normalizujSciezke(sciezka)
	if korzenNorm == "" || sciezkaNorm == "" {
		return false
	}
	wzgledna, err := filepath.Rel(korzenNorm, sciezkaNorm)
	if err != nil {
		return false
	}
	if wzgledna == "." {
		return true
	}
	return wzgledna != ".." && !strings.HasPrefix(wzgledna, ".."+string(filepath.Separator))
}

// normalizujSciezke sprowadza ścieżkę do postaci porównywalnej: bezwzględnej,
// oczyszczonej z członów `.` i `..`, a na Windowsie także z wielkości liter.
func normalizujSciezke(sciezka string) string {
	przycieta := strings.TrimSpace(sciezka)
	if przycieta == "" {
		return ""
	}
	if bezwzgledna, err := filepath.Abs(przycieta); err == nil {
		przycieta = bezwzgledna
	}
	przycieta = filepath.Clean(przycieta)
	if runtime.GOOS == "windows" {
		return strings.ToLower(przycieta)
	}
	return przycieta
}
