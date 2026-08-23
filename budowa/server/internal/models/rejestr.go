package models

import (
	"context"
	"strings"
	"sync"

	"danacoconsole/shared"
)

// Pominiety opisuje wiersz rejestru, dla którego adapter nie powstał. Rejestr
// takiego wiersza nie przemilcza i nie przerywa przez niego budowy — kanał
// nieudany jest jawny, pozostałe pracują.
type Pominiety struct {
	Kod     string `json:"code"`
	Adapter string `json:"adapter"`
	Powod   string `json:"reason"`
}

// Rejestr trzyma kanały modelu zbudowane w czasie działania z wierszy tabeli
// kanal_modelu. Żaden kanał nie jest wpisany w kod na sztywno: dopisanie kanału
// to dopisanie wiersza i odświeżenie rejestru.
type Rejestr struct {
	mu        sync.RWMutex
	zrodlo    ZrodloDefinicji
	fabryki   Fabryki
	kanaly    map[string]Kanal
	definicje map[string]Definicja
	wykaz     []Definicja
	pominiete []Pominiety
}

// NowyRejestr składa pusty rejestr nad źródłem wierszy i zestawem fabryk.
// Kanały powstają dopiero przy Odswiez — konstruktor niczego nie odpytuje.
func NowyRejestr(zrodlo ZrodloDefinicji, fabryki Fabryki) *Rejestr {
	pelne := Fabryki{}
	for klucz, fabryka := range fabryki {
		pelne[klucz] = fabryka
	}
	return &Rejestr{
		zrodlo:    zrodlo,
		fabryki:   pelne,
		kanaly:    map[string]Kanal{},
		definicje: map[string]Definicja{},
	}
}

// UstawFabryke dokłada fabrykę adaptera pod kluczem danych. Tą drogą pakiet
// internal/injection wnosi kanał główny CLI, nie zmieniając ani rejestru,
// ani wierszy w bazie.
func (r *Rejestr) UstawFabryke(klucz string, fabryka Fabryka) {
	if r == nil || strings.TrimSpace(klucz) == "" || fabryka == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fabryki[klucz] = fabryka
}

// Odswiez czyta wiersze na nowo i przebudowuje kanały. Wywoływany przy starcie
// oraz po każdej zmianie rejestru kanałów (komendy channel.*). Kanał, którego
// wiersz nie zmienił się, zachowuje swoją instancję — odświeżenie konfiguracji
// nie przerywa biegnącego strumienia innego kanału.
func (r *Rejestr) Odswiez(ctx context.Context) error {
	definicje, err := r.zrodlo.Definicje(ctx)
	if err != nil {
		return err
	}

	r.mu.Lock()
	poprzednie := r.kanaly
	poprzednieDef := r.definicje
	// Kopia zestawu fabryk: budowa kanałów biegnie bez zamka, a UstawFabryke
	// może w tym czasie dopisać fabrykę do mapy rejestru.
	fabryki := make(Fabryki, len(r.fabryki))
	for klucz, fabryka := range r.fabryki {
		fabryki[klucz] = fabryka
	}
	r.mu.Unlock()

	kanaly := map[string]Kanal{}
	nowe := map[string]Definicja{}
	wykaz := make([]Definicja, 0, len(definicje))
	pominiete := make([]Pominiety, 0)
	zachowane := map[Kanal]bool{}

	for _, d := range definicje {
		wykaz = append(wykaz, d)
		if !d.Aktywny {
			continue
		}
		kanal, zachowany, powod := zbudujKanal(d, fabryki, poprzednie, poprzednieDef)
		if kanal == nil {
			pominiete = append(pominiete, Pominiety{Kod: d.Kod, Adapter: d.KluczAdaptera(), Powod: powod})
			continue
		}
		if zachowany {
			zachowane[kanal] = true
		}
		for _, klucz := range kluczeKanalu(d) {
			kanaly[klucz] = kanal
			nowe[klucz] = d
		}
	}

	r.mu.Lock()
	r.kanaly, r.definicje, r.wykaz, r.pominiete = kanaly, nowe, wykaz, pominiete
	r.mu.Unlock()

	zamknijNieuzywane(poprzednie, zachowane)
	return nil
}

// Kanal zwraca kanał po identyfikatorze wiersza albo po kodzie kanału.
func (r *Rejestr) Kanal(klucz string) (Kanal, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	kanal, jest := r.kanaly[strings.TrimSpace(klucz)]
	return kanal, jest
}

// Definicja zwraca wiersz rejestru czynnego kanału.
func (r *Rejestr) Definicja(klucz string) (Definicja, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, jest := r.definicje[strings.TrimSpace(klucz)]
	return d, jest
}

// Wykaz zwraca wszystkie wiersze rejestru w kolejności odczytu, także nieczynne.
func (r *Rejestr) Wykaz() []Definicja {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Definicja(nil), r.wykaz...)
}

// Kontrakt zwraca wykaz kanałów w kształcie kontraktu, gotowy dla odpowiedzi
// channel.list. Ograniczenie do czynnych rozstrzyga wywołujący.
//
// „Czynny" to nie to samo co „włączony w wierszu". Wiersz może mieć aktywny = 1,
// a mimo to nie mieć zbudowanego adaptera — bo jego rodzaj nie ma fabryki albo
// fabryka odmówiła budowy (taki wiersz trafia do Pominiete). Kanał bez adaptera
// nie jest gotowy do pracy i przy tylkoCzynne nie pokazuje się jako czynny:
// inaczej okno wskazałoby kanał, który przy pierwszej turze odmówi.
func (r *Rejestr) Kontrakt(tylkoCzynne bool) []shared.Channel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	kanaly := make([]shared.Channel, 0, len(r.wykaz))
	for _, d := range r.wykaz {
		if tylkoCzynne {
			if !d.Aktywny {
				continue
			}
			if _, gotowy := r.kanaly[d.Identyfikator()]; !gotowy {
				continue
			}
		}
		kanaly = append(kanaly, d.Kontrakt())
	}
	return kanaly
}

// Pominiete zwraca wiersze, dla których adapter nie powstał.
func (r *Rejestr) Pominiete() []Pominiety {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Pominiety(nil), r.pominiete...)
}

// Zamknij zwalnia zasoby wszystkich kanałów rejestru.
func (r *Rejestr) Zamknij() error {
	r.mu.Lock()
	kanaly := r.kanaly
	r.kanaly, r.definicje = map[string]Kanal{}, map[string]Definicja{}
	r.mu.Unlock()
	zamknijNieuzywane(kanaly, map[Kanal]bool{})
	return nil
}
