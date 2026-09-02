// Rejestrator bloków wiadomości zapisuje nietekstowe fragmenty strumienia do tabeli blok_wiadomosci w trakcie tury; awaria zapisu zdejmuje okno z rejestracji do końca życia procesu.
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

type rejestratorBlokow struct {
	zycie    context.Context
	bloki    dane.RepozytoriumBlokow
	dziennik *log.Logger

	mu           sync.Mutex
	zdegradowane map[string]struct{}
}

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

// Zanotuj pomija fragment tekstowy, bo jego treść niesie `wiadomosc.tresc`.
// Zapis idzie na kontekście życia procesu z kontem tury (decyzja 34).
func (r *rejestratorBlokow) Zanotuj(ctx context.Context, f models.Fragment) {
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
	zapis := dane.ZKontemOperatora(r.zycie, dane.KontoOperatora(ctx))
	if err := r.bloki.Dopisz(zapis, blok); err != nil {
		r.zdegraduj(f.WindowId, err)
	}
}

func (r *rejestratorBlokow) czyZdegradowane(idOkna string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, zdegradowane := r.zdegradowane[idOkna]
	return zdegradowane
}

// zdegraduj zgłasza okno do dziennika raz; wpis przy każdym fragmencie zaśmiecałby ślad.
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
