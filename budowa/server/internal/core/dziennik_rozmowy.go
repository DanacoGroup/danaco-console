package core

import (
	"context"
	"log"
	"sync"

	"danacoconsole/shared"
)

// dziennikRozmowy prowadzi historię wiadomości okien komunikacji.
//
// Źródłem prawdy jest baza: każda wiadomość Operatora i każda odpowiedź modelu
// idzie do tabeli `wiadomosc` przez utrwalacz warstwy danych, a odczyt historii
// okna sięga po zapisane wiersze. Dzięki temu rozmowa przeżywa restart rdzenia.
//
// Pamięć procesu zostaje jako bufor podręczny o ograniczonej pojemności. Ma dwa
// zadania: odpowiadać bez odpytywania bazy w trakcie trwającej tury i przejąć
// rozmowę, gdy zapis zawiedzie. Awaria trwałości nie przerywa rozmowy — okno
// schodzi na bufor, zdarzenie trafia do dziennika procesu, tura biegnie dalej.
type dziennikRozmowy struct {
	mu           sync.RWMutex
	okna         map[string][]shared.Message
	zdegradowane map[string]struct{}
	limit        int

	zycie    context.Context
	trwalosc utrwalaczRozmowy
	dziennik *log.Logger
}

// utrwalaczRozmowy jest tą częścią warstwy danych, której dziennik naprawdę
// używa. Rdzeń nie zna SQL ani układu tabel.
type utrwalaczRozmowy interface {
	Zapisz(ctx context.Context, w shared.Message) error
	Zmien(ctx context.Context, w shared.Message) error
	Historia(ctx context.Context, idOkna string) ([]shared.Message, error)
}

// limitDziennikaDomyslny ogranicza liczbę wiadomości trzymanych w buforze na
// okno. Ograniczenie chroni pamięć rdzenia przy długiej rozmowie; starsze
// wiadomości wypadają z bufora, ale zostają w bazie.
const limitDziennikaDomyslny = 500

// nowyDziennikRozmowy zakłada dziennik nad utrwalaczem warstwy danych. Pusty
// utrwalacz jest dopuszczalny — dziennik pracuje wtedy wyłącznie na buforze
// i mówi o tym wprost w dzienniku procesu.
func nowyDziennikRozmowy(zycie context.Context, trwalosc utrwalaczRozmowy, dziennik *log.Logger) *dziennikRozmowy {
	if zycie == nil {
		zycie = context.Background()
	}
	return &dziennikRozmowy{
		okna:         map[string][]shared.Message{},
		zdegradowane: map[string]struct{}{},
		limit:        limitDziennikaDomyslny,
		zycie:        zycie,
		trwalosc:     trwalosc,
		dziennik:     dziennik,
	}
}

// Dopisz dokłada wiadomość na koniec historii okna: do bazy i do bufora.
func (d *dziennikRozmowy) Dopisz(w shared.Message) {
	d.mu.Lock()
	wykaz := append(d.okna[w.WindowId], w)
	if len(wykaz) > d.limit {
		wykaz = wykaz[len(wykaz)-d.limit:]
	}
	d.okna[w.WindowId] = wykaz
	d.mu.Unlock()

	d.utrwal(w.WindowId, "zapis wiadomości", func(ctx context.Context) error {
		return d.trwalosc.Zapisz(ctx, w)
	})
}

// Zmien podmienia wiadomość o tym samym identyfikatorze i mówi, czy trafiła
// w bufor. Tak domyka się odpowiedź modelu po zakończeniu strumienia.
func (d *dziennikRozmowy) Zmien(w shared.Message) bool {
	d.mu.Lock()
	trafiona := false
	wykaz := d.okna[w.WindowId]
	for i := range wykaz {
		if wykaz[i].Id == w.Id {
			wykaz[i] = w
			trafiona = true
			break
		}
	}
	d.mu.Unlock()

	d.utrwal(w.WindowId, "zapis odpowiedzi modelu", func(ctx context.Context) error {
		return d.trwalosc.Zmien(ctx, w)
	})
	return trafiona
}

// Wiadomosc odczytuje pojedynczą wiadomość okna z zapisanej historii.
func (d *dziennikRozmowy) Wiadomosc(idOkna, idWiadomosci string) (shared.Message, bool) {
	for _, w := range d.historia(idOkna) {
		if w.Id == idWiadomosci {
			return w, true
		}
	}
	return shared.Message{}, false
}

// Wykaz zwraca wiadomości okna od najstarszej, z opcjonalnym ograniczeniem
// liczby i odczytem sprzed wskazanej wiadomości. Druga wartość mówi, czy przed
// zwróconym wycinkiem zostały wiadomości starsze.
func (d *dziennikRozmowy) Wykaz(idOkna string, przed *string, ograniczenie *int) ([]shared.Message, bool) {
	wykaz := d.historia(idOkna)
	if przed != nil && *przed != "" {
		for i, w := range wykaz {
			if w.Id == *przed {
				wykaz = wykaz[:i]
				break
			}
		}
	}
	if ograniczenie == nil || *ograniczenie <= 0 || *ograniczenie >= len(wykaz) {
		return wykaz, false
	}
	poczatek := len(wykaz) - *ograniczenie
	return wykaz[poczatek:], poczatek > 0
}

// historia zwraca zapisaną rozmowę okna. Sięga do bazy, a gdy okno zeszło na
// bufor albo odczyt zawiódł — do pamięci procesu.
func (d *dziennikRozmowy) historia(idOkna string) []shared.Message {
	if d.trwalosc != nil && !d.czyZdegradowane(idOkna) {
		wykaz, err := d.trwalosc.Historia(d.zycie, idOkna)
		if err == nil {
			return wykaz
		}
		d.zdegraduj(idOkna, "odczyt historii okna", err)
	}
	return d.bufor(idOkna)
}

// bufor zwraca odpis podręcznej historii okna.
func (d *dziennikRozmowy) bufor(idOkna string) []shared.Message {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return append([]shared.Message(nil), d.okna[idOkna]...)
}

// utrwal wykonuje zapis do bazy. Niepowodzenie schodzi na bufor i nie wraca do
// wywołującego — rozmowa ma się toczyć dalej.
func (d *dziennikRozmowy) utrwal(idOkna, czynnosc string, praca func(context.Context) error) {
	if d.trwalosc == nil || d.czyZdegradowane(idOkna) {
		return
	}
	if err := praca(d.zycie); err != nil {
		d.zdegraduj(idOkna, czynnosc, err)
	}
}

// czyZdegradowane mówi, czy okno pracuje już wyłącznie na buforze.
func (d *dziennikRozmowy) czyZdegradowane(idOkna string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, zdegradowane := d.zdegradowane[idOkna]
	return zdegradowane
}

// zdegraduj przenosi okno na bufor i zgłasza zdarzenie do dziennika procesu.
// Zgłoszenie idzie raz na okno — powtarzanie tego samego wpisu przy każdej
// wiadomości zaśmiecałoby dziennik procesu.
func (d *dziennikRozmowy) zdegraduj(idOkna, czynnosc string, przyczyna error) {
	d.mu.Lock()
	_, juz := d.zdegradowane[idOkna]
	d.zdegradowane[idOkna] = struct{}{}
	d.mu.Unlock()
	if juz || d.dziennik == nil {
		return
	}
	d.dziennik.Printf("trwałość rozmowy: okno %s schodzi na bufor pamięci (%s): %v",
		idOkna, czynnosc, przyczyna)
}
