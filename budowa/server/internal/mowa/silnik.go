// Odpowiedzialność pliku: złożenie silnika mowy w jeden byt i jego
// zależności; jest jedno miejsce, w którym wiadomo, czy mowa działa i jak ją
// rozpoznać.
package mowa

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
)

// przedrostekTranskrypcji rozpoczyna identyfikator wpisu dziennika
// transkrypcji, żeby dało się go odróżnić od innych bytów rdzenia.
const przedrostekTranskrypcji = "mowa-"

// Silnik rozpoznaje mowę pomocnikiem lokalnym; ustawienia trzymane są
// w silniku, a nie odczytywane przy każdym zleceniu z bazy.
type Silnik struct {
	// uruchamiacz — jedyna droga startu procesu w drzewie.
	uruchamiacz session.Uruchamiacz
	// ustawienia — komplet nastaw obowiązujący dla zleceń tego silnika.
	ustawienia Ustawienia
	// dziennik — trwały ślad transkrypcji; nil znaczy pracę bez śladu.
	dziennik Dziennik
	// teraz oddaje czas w milisekundach epoki; wydzielone, bo jeden byt ma
	// mieć jeden zegar.
	teraz func() int64
}

// NowySilnik zakłada silnik na uruchamiaczu procesów; uruchamiacz jest
// jedyną zależnością obowiązkową, bez niego nie ma czym wywołać pomocnika.
func NowySilnik(uruchamiacz session.Uruchamiacz) *Silnik {
	return &Silnik{
		uruchamiacz: uruchamiacz,
		ustawienia:  UstawieniaDomyslne(),
		teraz:       func() int64 { return time.Now().UTC().UnixMilli() },
	}
}

// ZUstawieniami oddaje silnikowi komplet nastaw złożony z konfiguracji
// zasięgu, gotowy do użycia bez dalszego odczytu.
func (s *Silnik) ZUstawieniami(u Ustawienia) *Silnik {
	s.ustawienia = u
	return s
}

// ZDziennikiem oddaje silnikowi trwałość; zależność opcjonalna, bo silnik
// bez dziennika rozpoznaje mowę tak samo, traci wyłącznie ślad.
func (s *Silnik) ZDziennikiem(d Dziennik) *Silnik {
	s.dziennik = d
	return s
}

// Ustawienia oddaje nastawy, którymi silnik dziś pracuje, w kształcie
// złożonym przy zakładaniu albo późniejszej podmianie.
func (s *Silnik) Ustawienia() Ustawienia {
	return s.ustawienia
}

// katalogPracy wskazuje katalog uruchomienia pomocnika; pusty jest
// odpowiedzią poprawną, nie brakiem, bo brama izolacji uzupełnia go sama.
func (s *Silnik) katalogPracy(obszar session.Obszar) string {
	return strings.TrimSpace(obszar.KatalogRoboczy)
}

// odnotujGotowa zapisuje w dzienniku transkrypcję, która się udała; wynik
// zapisu jest pomijany, bo transkrypcja już jest w ręku wołającego.
func (s *Silnik) odnotujGotowa(ctx context.Context, z Zlecenie, t Transkrypcja) {
	s.zapiszWpis(ctx, WpisTranskrypcji{
		OknoId:           z.Okno.Id,
		NagranieOdnosnik: z.Odnosnik,
		Model:            t.Model,
		Jezyk:            t.Jezyk,
		Znakow:           int64(t.Znakow),
		TrwanieMs:        t.TrwanieMs,
		Stan:             stanZapisu(t),
	})
}

// stanZapisu przekłada stan transkrypcji na stan wiersza dziennika; dwa
// stany udane, nie jeden, bo nagranie bez mowy jest wynikiem, nie odmową.
func stanZapisu(t Transkrypcja) string {
	if t.Stan == StanBezMowy {
		return StanZapisuBezMowy
	}
	return StanGotowa
}

// odnotujOdmowe zapisuje w dzienniku odmowę i oddaje ją wołającemu bez
// zmiany, bo odmowa bez śladu jest nie do zdiagnozowania.
func (s *Silnik) odnotujOdmowe(ctx context.Context, z Zlecenie,
	model, jezyk string, powod error) error {

	s.zapiszWpis(ctx, WpisTranskrypcji{
		OknoId:           z.Okno.Id,
		NagranieOdnosnik: z.Odnosnik,
		Model:            model,
		Jezyk:            jezyk,
		Stan:             StanOdmowa,
		Powod:            powod.Error(),
	})
	return powod
}

// zapiszWpis dokłada identyfikator i czas, po czym oddaje wpis dziennikowi;
// zlecenie bez odnośnika nagrania nie jest zapisywane.
func (s *Silnik) zapiszWpis(ctx context.Context, wpis WpisTranskrypcji) {
	if s.dziennik == nil || strings.TrimSpace(wpis.NagranieOdnosnik) == "" {
		return
	}
	wpis.Identyfikator = nowyIdentyfikator()
	wpis.Utworzono = s.teraz()
	_, _ = s.dziennik.Zapisz(ctx, wpis)
}

// nowyIdentyfikator nadaje wpisowi dziennika tożsamość widoczną poza bazą;
// losowość z crypto/rand, a nie licznik, bo wpisy powstają współbieżnie.
func nowyIdentyfikator() string {
	var bajty [8]byte
	if _, err := rand.Read(bajty[:]); err != nil {
		return przedrostekTranskrypcji + time.Now().UTC().Format("20060102150405.000000000")
	}
	return przedrostekTranskrypcji + hex.EncodeToString(bajty[:])
}
