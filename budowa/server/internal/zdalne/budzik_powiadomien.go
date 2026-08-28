// Pakiet zdalne niesie pętlę tykającą kolejki powiadomień, wzorowaną na już działającej pętli tego rdzenia.
package zdalne

import (
	"context"
	"log"
	"time"
)

// InterwalBudzikaPowiadomien jest taktem pętli i wynosi trzydzieści sekund, ten sam rząd wielkości ponowienia.
const InterwalBudzikaPowiadomien = 30 * time.Second

// BudzikPowiadomien tyka kolejką powiadomień w całym cyklu życia rdzenia tej samej aplikacji serwerowej.
type BudzikPowiadomien struct {
	dziennik *log.Logger
	interwal time.Duration
}

// Funkcja NowyBudzikPowiadomien składa budzik powiadomień, gotowy do uruchomienia w cyklu życia rdzenia.
func NowyBudzikPowiadomien(dziennik *log.Logger) *BudzikPowiadomien {
	return &BudzikPowiadomien{dziennik: dziennik, interwal: InterwalBudzikaPowiadomien}
}

// Uruchom wpina budzik w cykl życia rdzenia.
//
// Wpięcie: jedna linia w kompozycji, obok `zdalne.Zasil(baza.DB)`. Budzik bez
// wpięcia nie tyka i nic tego nie ukrywa:
// powiadomienia stoją w kolejce ze stanem `oczekuje`, a nie znikają.
func (b *BudzikPowiadomien) Uruchom(zycie context.Context) {
	if b == nil {
		return
	}
	go b.petla(zycie)
}

func (b *BudzikPowiadomien) petla(zycie context.Context) {
	zegar := time.NewTicker(b.interwal)
	defer zegar.Stop()
	b.przebieg(zycie, time.Now().UTC())
	for {
		select {
		case <-zycie.Done():
			return
		case teraz := <-zegar.C:
			b.przebieg(zycie, teraz.UTC())
		}
	}
}

// Metoda przebieg wykonuje jeden takt: najpierw gasi przeterminowane powiadomienia, potem doręcza należne.
func (b *BudzikPowiadomien) przebieg(ctx context.Context, teraz time.Time) {
	if wygaslo, err := Wygas(ctx); err != nil {
		b.zapisz("budzik powiadomień: wygaszanie przeterminowanych nie powiodło się: %v", err)
	} else if wygaslo > 0 {
		b.zapisz("budzik powiadomień: wygasło %d powiadomień, którym minął termin ważności", wygaslo)
	}

	podsumowanie, err := Wyslij(ctx, teraz)
	if err != nil {
		// Odmowa przebiegu jest zdaniem trójczęściowym i ma trafić do dziennika w całości.
		b.zapisz("budzik powiadomień: %v", err)
		return
	}
	if podsumowanie.Rozpatrzonych == 0 {
		return
	}
	b.zapisz("budzik powiadomień: rozpatrzono %d, doręczono %d, pominięto po decyzji %d, "+
		"odłożono bez odbiorcy %d, wygasło w takcie %d",
		podsumowanie.Rozpatrzonych, podsumowanie.Doreczonych, podsumowanie.Pominietych,
		podsumowanie.Przelozonych, podsumowanie.Wygaslych)
}

func (b *BudzikPowiadomien) zapisz(wzor string, argumenty ...any) {
	if b.dziennik == nil {
		return
	}
	b.dziennik.Printf(wzor, argumenty...)
}
