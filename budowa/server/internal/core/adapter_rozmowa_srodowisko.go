// Plik doprowadza do zapytania kanału katalog roboczy sesji i konfigurację
// mostów MCP okna rozmowy, dwa różne byty pochodzące z osobnych dróg
// rozstrzygania.
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

// ZMostami wpina składacz konfiguracji mostów MCP okna rozmowy do zapytania
// kierowanego do kanału modelu.
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
// oraz parametry wykonania z konfiguracji, osobno od czystej funkcji
// zapytanieKanalu.
func (a *adapterRozmowy) uzupelnijSrodowisko(ctx context.Context, okno session.Okno,
	zapytanie *models.Zapytanie) {

	if a == nil || zapytanie == nil {
		return
	}
	zapytanie.KatalogSesji = a.katalogSesji(okno)
	if a.mosty != nil {
		// Konfiguracja MCP okna niesie mosty z nadań i narzędzia platformy pod jednym przełącznikiem.
		zapytanie.KonfiguracjaMCP = a.zKonfiguracjaZestawu(ctx, okno,
			a.mosty.tekstZNarzedziami(ctx, okno.Id))
	}
	a.uzupelnijWykonanie(ctx, okno, zapytanie)
}

// uzupelnijWykonanie przenosi ustawienia uczestnika rozmowy do zapytania
// kanału warunkowo: wartość pusta nie nadpisuje tego, co przyszło z okna albo
// z wiersza rejestru.
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
