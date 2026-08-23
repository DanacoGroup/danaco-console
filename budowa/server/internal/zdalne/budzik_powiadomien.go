// Odpowiedzialność pliku: pętla tykająca kolejki powiadomień. Bez niej
// „ponowienie i wygaśnięcie" byłoby trzema kolumnami, których nikt nigdy nie
// rusza — czyli atrapą wymagania, a nie jego spełnieniem.
//
// Wzorzec wzięty z działającej pętli rdzenia (core/harmonogram_budzik.go):
// ticker, `select` na kontekście życia, pierwszy przebieg od razu po starcie
// i awaria jednego przebiegu, która nie zatrzymuje pętli. Pierwszy przebieg od
// razu ma znaczenie akurat tutaj: powiadomienia zgłoszone tuż przed postojem
// rdzenia mają dolecieć zaraz po jego powrocie, a nie po pełnym takcie.
package zdalne

import (
	"context"
	"log"
	"time"
)

// InterwalBudzikaPowiadomien jest taktem pętli. Trzydzieści sekund to ten sam
// rząd wielkości, co najkrótszy odstęp ponowienia — takt gęstszy nie przyspiesza
// niczego (i tak czeka na `nastepna_proba`), a rzadszy opóźniałby pierwsze
// podejście o więcej, niż wynosi cała jego zwłoka.
const InterwalBudzikaPowiadomien = 30 * time.Second

// BudzikPowiadomien tyka kolejką w cyklu życia rdzenia.
type BudzikPowiadomien struct {
	dziennik *log.Logger
	interwal time.Duration
}

// NowyBudzikPowiadomien składa budzik. Dziennik pusty wyłącza wyłącznie ślad,
// nie budzik.
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

// przebieg wykonuje jeden takt: najpierw gasi przeterminowane, potem doręcza
// należne.
//
// Kolejność jest rozmyślna. Wygaszanie idzie pierwsze, żeby przeterminowane
// powiadomienie nie zdążyło polecieć w tym samym takcie, w którym straciło
// ważność. Takt to zawsze jakiś kawałek czasu; gdyby wysyłka szła pierwsza,
// budzik zabrzmiałby o sprawie, o której sam za chwilę orzeka, że jest
// nieaktualna.
//
// Wygaszanie idzie także wtedy, gdy nadajnika nie ma — sprawa nieaktualna jest
// nieaktualna niezależnie od tego, czy było komu ją zanieść.
func (b *BudzikPowiadomien) przebieg(ctx context.Context, teraz time.Time) {
	if wygaslo, err := Wygas(ctx); err != nil {
		b.zapisz("budzik powiadomień: wygaszanie przeterminowanych nie powiodło się: %v", err)
	} else if wygaslo > 0 {
		b.zapisz("budzik powiadomień: wygasło %d powiadomień, którym minął termin ważności", wygaslo)
	}

	podsumowanie, err := Wyslij(ctx, teraz)
	if err != nil {
		// Odmowa przebiegu jest zdaniem trójczęściowym (patrz Wyslij) i ma
		// trafić do dziennika w całości — brak nadajnika ma być widoczny,
		// a nie odgadywany z tego, że nic nie dolatuje.
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
