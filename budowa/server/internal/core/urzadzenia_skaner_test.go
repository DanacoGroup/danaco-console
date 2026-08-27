// Sprawdzian odmów warstwy urządzeń przy braku programu samej warstwy: droga
// SANE (`scanimage`) ma odmawiać tak samo, jak droga WIA (`pwsh`) — zdaniem
// nazywającym brak, naprawę i drogę obejścia `studio.ingest.queue.add`.
//
// Sprawdzian nie uruchamia ani jednego programu. Brak jest wymuszony odcięciem
// PATH, nie stanem maszyny — na stanowisku deweloperskim `scanimage` bywa
// zainstalowany, więc sprawdzian liczący na jego nieobecność kłamałby tam,
// gdzie się go uruchamia. Uruchamiacz-atrapa pilnuje drugiej strony pomiaru:
// odmowa braku ma zapaść PRZED startem procesu, więc dojście do uruchamiacza
// jest usterką mierzonej drogi, nie zdarzeniem.
package core

import (
	"context"
	"runtime"
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// uruchamiaczNieruszany jest atrapą portu session.Uruchamiacz dla pomiarów
// odmów zapadających przed startem procesu.
type uruchamiaczNieruszany struct{ t *testing.T }

func (u uruchamiaczNieruszany) UruchomProces(session.Okno,
	session.Polecenie) (session.UchwytProcesu, error) {
	u.t.Fatal("mimo braku programu warstwy uruchomiono proces — odmowa miała " +
		"zapaść przed startem")
	return nil, nil
}

// odetnijProgramyWarstw wymusza brak programów obu warstw skanera i potwierdza,
// że wymuszenie zadziałało — pomiar na maszynie, na której program dalej widać,
// mierzyłby inną drogę i meldował ją jako wynik.
func odetnijProgramyWarstw(t *testing.T) {
	t.Helper()
	t.Setenv("PATH", t.TempDir())
	if zewnetrzne.Stoi(narzedzieSkanera) {
		t.Fatalf("pomiar niewykonany: program %s widać mimo odcięcia PATH",
			narzedzieSkanera.Program)
	}
	if zewnetrzne.Stoi(narzedzieSkaneraWia) {
		t.Fatalf("pomiar niewykonany: program %s widać mimo odcięcia PATH",
			narzedzieSkaneraWia.Program)
	}
}

// czlonyOdmowyBrakuSane to trzy rzeczy, które odmowa braku programu SANE musi
// nazwać naraz — czym rdzeń skanuje, czym brak naprawić i którą drogą materiał
// wchodzi mimo braku. Ten sam komplet niesie droga WIA w `bladWarstwyWia`.
func czlonyOdmowyBrakuSane() []string {
	return []string{narzedzieSkanera.Program, narzedzieSkanera.Pakiet,
		"studio.ingest.queue.add"}
}

func TestWykazSkanerowNazywaBrakProgramuSaneIDrogeObejscia(t *testing.T) {
	if runtime.GOOS != systemLinux {
		t.Skipf("pomiar niewykonany: droga SANE dotyczy Linuksa, a maszyna to %s",
			runtime.GOOS)
	}
	odetnijProgramyWarstw(t)
	a := &adapterStudia{uruchamiacz: uruchamiaczNieruszany{t: t}}

	urzadzenia, err := a.wykazSkanerow(context.Background())
	if err == nil {
		t.Fatalf("brak programu warstwy oddany jako wykaz (%d pozycji) — to "+
			"twierdzenie „szukałem i nic nie ma”, którego rdzeń nie ma prawa "+
			"postawić", len(urzadzenia))
	}
	odmowa := protocol.BladZeZrodla(shared.ErrorCodeInternalError, err)
	if odmowa.Code != shared.ErrorCodeChannelUnavailable {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q — brak programu warstwy "+
			"jest zapleczem niedostępnym, tak samo jak na drodze WIA; treść: %s",
			odmowa.Code, shared.ErrorCodeChannelUnavailable, odmowa.Message)
	}
	for _, czlon := range czlonyOdmowyBrakuSane() {
		if !strings.Contains(odmowa.Message, czlon) {
			t.Errorf("odmowa wykazu nie niesie członu %q: %s", czlon, odmowa.Message)
		}
	}
}

func TestSkanowanieNazywaBrakProgramuSaneIDrogeObejscia(t *testing.T) {
	if runtime.GOOS != systemLinux {
		t.Skipf("pomiar niewykonany: droga SANE dotyczy Linuksa, a maszyna to %s",
			runtime.GOOS)
	}
	odetnijProgramyWarstw(t)
	// Katalog skanów powstaje przed rozstrzygnięciem drogi, a adapter bez
	// katalogu roboczego schodzi na katalog tymczasowy systemu — ten idzie pod
	// katalog sprawdzianu, żeby pomiar nie zostawił śladu na maszynie.
	t.Setenv("TMPDIR", t.TempDir())
	a := &adapterStudia{uruchamiacz: uruchamiaczNieruszany{t: t}}

	sciezki, err := a.skanujUrzadzenie(context.Background(), zamowienieSkanu{})
	if err == nil {
		t.Fatalf("brak programu warstwy oddany jako udany skan: %v", sciezki)
	}
	odmowa := protocol.BladZeZrodla(shared.ErrorCodeInternalError, err)
	if odmowa.Code != shared.ErrorCodeChannelUnavailable {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q — brak programu warstwy "+
			"jest zapleczem niedostępnym, tak samo jak na drodze WIA; treść: %s",
			odmowa.Code, shared.ErrorCodeChannelUnavailable, odmowa.Message)
	}
	for _, czlon := range czlonyOdmowyBrakuSane() {
		if !strings.Contains(odmowa.Message, czlon) {
			t.Errorf("odmowa skanowania nie niesie członu %q: %s", czlon, odmowa.Message)
		}
	}
}

// TestOdmowaBrakuProgramuJestTaSamaNaObuDrogachSkanera mierzy parytet na
// poziomie mechanizmu, bo gałęzi Windows nie da się wziąć na Linuksie przez
// `runtime.GOOS`: obie odmowy składa dokładnie tak, jak składają je gałęzie
// `wykazSkanerow` — wołaniem warstwy i przekładem braku — i wymaga od nich
// jednego kodu oraz tej samej drogi obejścia.
func TestOdmowaBrakuProgramuJestTaSamaNaObuDrogachSkanera(t *testing.T) {
	odetnijProgramyWarstw(t)
	a := &adapterStudia{uruchamiacz: uruchamiaczNieruszany{t: t}}
	ctx := context.Background()

	_, bladSane := a.wolajUrzadzenie(ctx, narzedzieSkanera, []string{"-L"},
		granicaWykazuUrzadzen)
	if bladSane == nil {
		t.Fatal("pomiar niewykonany: droga SANE nie odmówiła mimo odcięcia programu")
	}
	odmowaSane := protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		bladWarstwySane(bladSane))

	_, bladWia := a.wolajUrzadzenie(ctx, narzedzieSkaneraWia,
		skryptPowerShell(skryptWykazuWia), granicaWykazuUrzadzen)
	if bladWia == nil {
		t.Fatal("pomiar niewykonany: droga WIA nie odmówiła mimo odcięcia programu")
	}
	odmowaWia := protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		bladWarstwyWia(bladWia))

	if odmowaSane.Code != odmowaWia.Code {
		t.Errorf("jeden brak, dwa kody: droga SANE %q, droga WIA %q",
			odmowaSane.Code, odmowaWia.Code)
	}
	if !strings.Contains(odmowaSane.Message, "studio.ingest.queue.add") {
		t.Errorf("droga SANE nie wskazuje drogi obejścia: %s", odmowaSane.Message)
	}
	if !strings.Contains(odmowaWia.Message, "studio.ingest.queue.add") {
		t.Errorf("droga WIA nie wskazuje drogi obejścia: %s", odmowaWia.Message)
	}
}
