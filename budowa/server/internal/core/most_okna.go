// Odpowiedzialność pliku: złożenie konfiguracji mostów MCP dla jednego okna
// rozmowy — od wierszy nadań i punktów do gotowego tekstu `mcpServers`.
//
// Tu domyka się dostęp. Punkt dostępu i nadanie bez tego pliku byłyby wpisami
// w oknie konfiguracji i niczym więcej: model nie dostawałby ani jednego mostu.
// Ścieżka jest jedna — nadania okna → punkty → generator wpisów
// (`most_mcp.go`) → plik `--mcp-config` procesu modelu.
//
// Słownictwo trybu jest danymi. Argument, jakim dany most nazywa tryb
// nadania, pochodzi z wiersza `argument_trybu_mostu` tego punktu. Rdzeń nie zna
// ani jednego takiego słowa; brak wiersza znaczy uruchomienie bez argumentu,
// czyli tryb odczytu — wariant bezpieczniejszy.
//
// Błąd odczytu nie zerwie tury. Okno bez mostów rozmawia dalej, tyle że bez
// wglądu w maszyny. Odmowa rozmowy z powodu niedostępnego katalogu dostępów
// byłaby bramą, której dokumentacja nie stawia.
package core

import (
	"context"
	"log"
	"sync"

	"danacoconsole/server/internal/dane"
)

// mostyOkna składa konfigurację mostów z repozytoriów warstwy danych.
//
// Dziennik i rejestr zgłoszonych braków służą wyłącznie wpisowi serwera
// narzędzi (`most_narzedzi.go`): odmowa dołożenia tego wpisu ma być
// powiedziana, a nie przemilczana, i ma być powiedziana raz na powód.
type mostyOkna struct {
	nadania        dane.RepozytoriumNadan
	punkty         dane.RepozytoriumPunktowDostepu
	okna           dane.RepozytoriumOkien
	dziennik       *log.Logger
	zgloszoneBraki sync.Map
}

// ZDziennikiem wskazuje dziennik rdzenia, do którego idą meldunki o braku
// serwera narzędzi. Bez niego składacz działa tak samo, tylko milczy — dlatego
// jest osobnym krokiem, a nie argumentem konstruktora.
func (m *mostyOkna) ZDziennikiem(dziennik *log.Logger) *mostyOkna {
	if m != nil {
		m.dziennik = dziennik
	}
	return m
}

// noweMostyOkna wiąże składacz z trzema repozytoriami. Nadanie niesie numery
// wierszy okna i punktu, a rdzeń posługuje się ich identyfikatorami trwałymi.
func noweMostyOkna(nadania dane.RepozytoriumNadan, punkty dane.RepozytoriumPunktowDostepu,
	okna dane.RepozytoriumOkien) *mostyOkna {

	return &mostyOkna{nadania: nadania, punkty: punkty, okna: okna}
}

// Tekst zwraca konfigurację MCP okna w postaci tekstu JSON. Okno bez nadań
// czynnych daje napis pusty — wywołujący nie zapisuje wtedy pliku i nie podaje
// przełącznika `--mcp-config`.
func (m *mostyOkna) Tekst(ctx context.Context, idOkna string) string {
	nadania := m.nadaniaMostow(ctx, idOkna)
	if len(nadania) == 0 {
		return ""
	}
	tekst, err := TekstKonfiguracjiMostu(nadania)
	if err != nil {
		return ""
	}
	return tekst
}

// nadaniaMostow czyta zbiór nadań okna i wiąże każde z jego punktem. Punkty
// rodzaju localDirectory i wpisy wygaszone odsiewa sam generator, więc tutaj
// przechodzi wszystko, co ma wiersz punktu.
func (m *mostyOkna) nadaniaMostow(ctx context.Context, idOkna string) []NadanieMostu {
	if m == nil || m.nadania == nil || m.punkty == nil || m.okna == nil || idOkna == "" {
		return nil
	}
	okno, err := m.okna.PoIdentyfikatorze(ctx, idOkna)
	if err != nil {
		return nil
	}
	wiersze, err := m.nadania.ListaOkna(ctx, okno.ID, true)
	if err != nil {
		return nil
	}
	punkty := map[int64]dane.PunktDostepu{}
	nadania := make([]NadanieMostu, 0, len(wiersze))
	for _, wiersz := range wiersze {
		punkt, jest := m.punkt(ctx, punkty, wiersz.PunktDostepuID)
		if !jest {
			continue
		}
		nadania = append(nadania, NadanieMostu{
			Punkt:         punktKontraktu(punkt),
			Nadanie:       nadanieKontraktu(wiersz, punkt.Kod, idOkna),
			ArgumentTrybu: punkt.ArgumentTrybu(wiersz.Tryb),
		})
	}
	return nadania
}

// punkt odczytuje wiersz punktu, korzystając z pamięci podręcznej przebiegu:
// okno bywa związane z kilkoma nadaniami tego samego mostu (odczyt i zapis
// osobno), więc bez niej ten sam wiersz szedłby z bazy wielokrotnie.
func (m *mostyOkna) punkt(ctx context.Context, pamiec map[int64]dane.PunktDostepu,
	punktID int64) (dane.PunktDostepu, bool) {

	if punkt, jest := pamiec[punktID]; jest {
		return punkt, punkt.ID != 0
	}
	punkt, err := m.punkty.Pobierz(ctx, punktID)
	if err != nil {
		pamiec[punktID] = dane.PunktDostepu{}
		return dane.PunktDostepu{}, false
	}
	pamiec[punktID] = punkt
	return punkt, true
}
