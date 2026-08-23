// Odpowiedzialność pliku: rejestrator bloków wiadomości — zapis nietekstowych
// fragmentów strumienia odpowiedzi do tabeli `blok_wiadomosci`
// (`store/migracja_096_tresc_rozmowy.sql`) W TRAKCIE tury.
//
// GDZIE STOI I DLACZEGO TU. Jedynym miejscem, które widzi KAŻDY fragment
// strumienia przed spakowaniem w kopertę, jest ujście tury w adapterze rozmowy
// (adapter_rozmowa.go, prowadzTure). Dziennik rozmowy i utrwalacz widzą dopiero
// domkniętą wiadomość — czyli sam sklejony tekst; wpięcie tam wymagałoby
// przenoszenia bloków przez typ Message, którego kontrakt nie modeluje.
// Rejestrator jest więc portem adaptera, wzorem odbiornika zdarzeń
// wykonawczych: adapter woła jedną metodę, montaż podaje całość.
//
// CO ZAPISUJE. Wszystkie rodzaje fragmentów POZA tekstem: rozumowanie,
// wywołania narzędzi z wynikami, obraz, dźwięk, błąd, prowenancję i metadane
// konta. Tekst pomija świadomie — treść tekstową domyka dziennik rozmowy
// w `wiadomosc.tresc`, a druga kopia byłaby drugą prawdą.
//
// AWARIA ZAPISU NIE PRZERYWA TURY. Okno, którego zapis bloków
// zawiódł, schodzi z rejestracji do końca życia procesu — wzorem degradacji
// dziennika rozmowy: jedno zgłoszenie do dziennika rdzenia, zero powtórek.
package core

import (
	"context"
	"log"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// rejestratorBlokow zapisuje bloki strumienia przez repozytorium warstwy danych.
type rejestratorBlokow struct {
	zycie    context.Context
	bloki    dane.RepozytoriumBlokow
	dziennik *log.Logger

	mu           sync.Mutex
	zdegradowane map[string]struct{}
}

// nowyRejestratorBlokow wiąże rejestrator z repozytorium bloków. Puste
// repozytorium daje rejestrator bierny — rdzeń bez bazy prowadzi rozmowę
// bez zapisu bloków, tak jak bez zapisu wiadomości.
func nowyRejestratorBlokow(zycie context.Context, bloki dane.RepozytoriumBlokow,
	dziennik *log.Logger) *rejestratorBlokow {

	if zycie == nil {
		zycie = context.Background()
	}
	return &rejestratorBlokow{
		zycie: zycie, bloki: bloki, dziennik: dziennik,
		zdegradowane: map[string]struct{}{},
	}
}

// Zanotuj zapisuje jeden fragment strumienia jako blok wiadomości. Fragment
// tekstowy i fragment bez wiadomości przechodzą bez śladu — pierwszy ma swoją
// prawdę w `wiadomosc.tresc`, drugi nie ma czego wskazać.
func (r *rejestratorBlokow) Zanotuj(f models.Fragment) {
	if r == nil || r.bloki == nil {
		return
	}
	if f.Kind == shared.ChunkKindText || f.MessageId == "" {
		return
	}
	if r.czyZdegradowane(f.WindowId) {
		return
	}
	blok := dane.BlokWiadomosci{
		Chwila:       time.Now().UnixMilli(),
		OknoKod:      f.WindowId,
		WiadomoscKod: f.MessageId,
		Rodzaj:       f.Kind,
		Tresc:        models.TrescFragmentu(f),
		Ladunek:      f.Data,
	}
	if err := r.bloki.Dopisz(r.zycie, blok); err != nil {
		r.zdegraduj(f.WindowId, err)
	}
}

// czyZdegradowane mówi, czy okno wypadło z rejestracji bloków.
func (r *rejestratorBlokow) czyZdegradowane(idOkna string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, zdegradowane := r.zdegradowane[idOkna]
	return zdegradowane
}

// zdegraduj wyłącza rejestrację bloków okna i zgłasza to raz do dziennika
// rdzenia — powtarzanie wpisu przy każdym fragmencie zaśmiecałoby ślad.
func (r *rejestratorBlokow) zdegraduj(idOkna string, przyczyna error) {
	r.mu.Lock()
	_, juz := r.zdegradowane[idOkna]
	r.zdegradowane[idOkna] = struct{}{}
	r.mu.Unlock()
	if juz || r.dziennik == nil {
		return
	}
	r.dziennik.Printf("zapis bloków rozmowy: okno %s bez rejestracji bloków do końca procesu: %v",
		idOkna, przyczyna)
}
