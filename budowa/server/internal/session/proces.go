package session

import (
	"sync"
	"time"
)

// Proces jest procesem jednego okna komunikacji, a sesja prowadzi tyle procesów, ile ma otwartych okien, i wszystkie biegną równolegle.
type Proces struct {
	// IdOkna — klucz procesu w rejestrze procesów.
	IdOkna string
	// IdProcesu — identyfikator nadany przez rdzeń (telemetria).
	IdProcesu string
	// Uruchomiono — chwila objęcia procesu uchwytem sesji.
	Uruchomiono time.Time

	drzewo *DrzewoProcesu

	mu         sync.Mutex
	zakonczony bool
}

// Ubij kończy proces okna wraz z całym drzewem potomstwa i oddaje uchwyty
// systemowe. Ubicie jest jednorazowe: drugie wywołanie nie sięga po zwolniony
// już uchwyt.
func (p *Proces) Ubij() error {
	p.mu.Lock()
	if p.zakonczony {
		p.mu.Unlock()
		return nil
	}
	p.zakonczony = true
	p.mu.Unlock()

	err := p.drzewo.Ubij()
	// Uchwyt zadania i uchwyt procesu oddajemy zaraz po ubiciu, aby nie zalegał w rdzeniu do końca pracy.
	p.drzewo.Zwolnij()
	return err
}

// Metoda dogladaj czeka na faktyczne zakończenie procesu okna i domyka jego cykl życia, oznaczając proces jako zakończony oraz oddając uchwyty systemowe.
func (p *Proces) dogladaj() {
	p.drzewo.Czekaj()

	p.mu.Lock()
	juz := p.zakonczony
	p.zakonczony = true
	p.mu.Unlock()
	if juz {
		// Ubicie przeszło bramkę pierwsze i uchwyty już oddało; wychodzi obserwator poprzedniej tury.
		return
	}
	p.drzewo.Zwolnij()
}

// Metoda Zyje mówi, czy proces okna nadal pracuje, sprawdzając bieżący stan uchwytu systemowego procesu.
func (p *Proces) Zyje() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return !p.zakonczony
}
