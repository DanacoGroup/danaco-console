// Odpowiedzialność pliku: doprowadzenie do zapytania kanału dwóch rzeczy, bez
// których model by ich nie zobaczył — katalogu roboczego sesji i konfiguracji
// mostów MCP okna rozmowy.
//
// To dwa różne byty i dwie różne drogi. Katalog roboczy mówi, gdzie model
// zostawia własne pliki, i pochodzi z rozstrzygnięcia dwóch kluczy katalogu
// ustawień. Nadania dostępu mówią, do czego model sięga, i pochodzą ze zbioru
// nadań tego jednego okna. Ani jedno nie wynika z drugiego, więc ani jedno nie
// jest liczone z drugiego.
//
// Brak ustalenia katalogu zostawia domyślne zachowanie kanału (pierwszy katalog
// roboczy okna), a brak nadań zostawia proces bez przełącznika `--mcp-config`.
// Rozmowa toczy się w obu przypadkach.
package core

import (
	"context"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
)

// ZKatalogiemRoboczym wpina ustalacz katalogu sesji. Bez niego katalogiem
// startowym procesu zostaje pierwszy katalog roboczy okna.
func (a *adapterRozmowy) ZKatalogiemRoboczym(k *KatalogRoboczy) *adapterRozmowy {
	a.katalog = k
	return a
}

// ZMostami wpina składacz konfiguracji mostów MCP okna rozmowy.
func (a *adapterRozmowy) ZMostami(m *mostyOkna) *adapterRozmowy {
	a.mosty = m
	return a
}

// ZParametramiWykonania wpina odczyt nakładu rozumowania i modelu zapasowego
// z konfiguracji. Adapter bez tego portu prowadzi turę na ustawieniach kanału.
func (a *adapterRozmowy) ZParametramiWykonania(p ParametryWykonania) *adapterRozmowy {
	a.wykonanie = p
	return a
}

// uzupelnijSrodowisko dokłada do zapytania katalog sesji, konfigurację mostów
// oraz parametry wykonania z konfiguracji. Osobno od zapytanieKanalu, bo tamta
// funkcja jest czysta — bierze okno i dwie wiadomości, a te wartości wymagają
// odczytu konfiguracji i bazy.
func (a *adapterRozmowy) uzupelnijSrodowisko(ctx context.Context, okno session.Okno,
	zapytanie *models.Zapytanie) {

	if a == nil || zapytanie == nil {
		return
	}
	zapytanie.KatalogSesji = a.katalogSesji(okno)
	if a.mosty != nil {
		// Konfiguracja MCP okna niesie dwie różne rzeczy, które jadą jednym
		// przełącznikiem: mosty do maszyn (z nadań okna) i narzędzia sterujące
		// platformą. Pierwsze przysługują z nadania, drugie każdemu oknu, które
		// w ogóle rozmawia — inaczej model nie otworzyłby modułu bez wglądu
		// w serwer, co popychałoby Operatora do rozdawania dostępu, którego nikt
		// nie potrzebuje. Granica uprawnień siedzi wewnątrz wykazu narzędzi:
		// zapisy zastrzeżone są poza nim strukturalnie.
		//
		// Zestaw narzędzi tury składa się tutaj, nie przy starcie procesu. Wynik
		// jedynego składacza konfiguracji przechodzi przez dopisanie zestawu
		// (adapter_rozmowa_zestaw.go): podstawa z definicji eksperta plus doraźne
		// dołożenia sesji. Drugiej konfiguracji nie ma; zestaw niezawężony nie
		// dokłada nic i tura jedzie pełnym wykazem kontraktu.
		zapytanie.KonfiguracjaMCP = a.zKonfiguracjaZestawu(ctx, okno,
			a.mosty.tekstZNarzedziami(ctx, okno.Id))
	}
	a.uzupelnijWykonanie(ctx, okno, zapytanie)
}

// uzupelnijWykonanie przenosi ustawienia Operatora do zapytania kanału.
//
// Wartość pusta nie nadpisuje tego, co przyszło z okna albo z wiersza rejestru:
// „bez wskazania na żadnym poziomie" znaczy „zostaw decyzję kanałowi", a nie
// „wyczyść". Dlatego przypisanie jest warunkowe, nie bezwarunkowe.
func (a *adapterRozmowy) uzupelnijWykonanie(ctx context.Context, okno session.Okno,
	zapytanie *models.Zapytanie) {

	if a.wykonanie == nil {
		return
	}
	wykonanie := a.wykonanie.Ustal(ctx, okno)
	if wykonanie.Naklad != "" {
		zapytanie.NakladRozumowania = wykonanie.Naklad
	}
	if wykonanie.ModelZapasowy != "" {
		zapytanie.ModelZapasowy = wykonanie.ModelZapasowy
	}
	if wykonanie.PulapKosztuUSD > 0 {
		zapytanie.PulapKosztuUSD = wykonanie.PulapKosztuUSD
	}
}

// katalogSesji ustala katalog roboczy tej sesji w kontekście zasięgu okna.
// Poziom okna jest najwęższy, więc wskazanie okna wygrywa z każdym
// szerszym zapisem Operatora.
func (a *adapterRozmowy) katalogSesji(okno session.Okno) string {
	if a.katalog == nil {
		return ""
	}
	return a.katalog.Ustal(konfig.Kontekst{Okno: okno.Id}, okno.IdSesji).Sciezka
}
