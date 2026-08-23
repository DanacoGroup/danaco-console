// Odpowiedzialność pliku: złożenie silnika mowy w jeden byt i jego zależności.
//
// Jest jedno miejsce, w którym wiadomo, czy mowa działa i jak ją rozpoznać.
// Odbiorcy dostają silnik wstrzyknięciem i nie budują własnego.
//
// Czego silnik nie robi. Nie syntezuje mowy — `translate.speech.synthesize`
// idzie w drugą stronę (tekst na dźwięk) i potrzebuje innego silnika, którego
// ten pakiet nie wnosi. Nie nagrywa: nagranie powstaje po stronie Operatora,
// a tutaj przychodzi już jako ścieżka pliku. Nie zna kontraktu: mówi własnymi
// typami.
//
// Dźwięk nie opuszcza maszyny Operatora: nie ma tu klienta HTTP ani adresu, pod
// który cokolwiek by poszło. Jedyne wyjście na zewnątrz procesu to uruchomienie
// pomocnika lokalnego.
package mowa

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
)

// przedrostekTranskrypcji rozpoczyna identyfikator wpisu dziennika.
const przedrostekTranskrypcji = "mowa-"

// Silnik rozpoznaje mowę pomocnikiem lokalnym.
//
// Ustawienia trzymane są w silniku, a nie odczytywane przy każdym zleceniu:
// wołający składa je z konfiguracji zasięgu (`config.get`) i podaje gotowe.
// Silnik nie sięga do bazy po nastawy, bo dwie drogi do tej samej wartości
// byłyby dwiema prawdami.
type Silnik struct {
	// uruchamiacz — jedyna droga startu procesu w drzewie.
	uruchamiacz session.Uruchamiacz
	// ustawienia — komplet nastaw obowiązujący dla zleceń tego silnika.
	ustawienia Ustawienia
	// dziennik — trwały ślad transkrypcji; nil znaczy pracę bez śladu.
	dziennik Dziennik
	// teraz oddaje czas w milisekundach epoki. Wydzielone w pole, bo czas
	// wpisu podaje wołający, a jeden byt ma mieć jeden zegar.
	teraz func() int64
}

// NowySilnik zakłada silnik na uruchamiaczu procesów.
//
// Uruchamiacz jest jedyną zależnością obowiązkową, bo bez niego nie ma czym
// wywołać pomocnika. Ustawienia startują wartościami domyślnymi — brak
// wskazania Operatora znaczy wartość domyślną, nie odmowę pracy.
func NowySilnik(uruchamiacz session.Uruchamiacz) *Silnik {
	return &Silnik{
		uruchamiacz: uruchamiacz,
		ustawienia:  UstawieniaDomyslne(),
		teraz:       func() int64 { return time.Now().UTC().UnixMilli() },
	}
}

// ZUstawieniami oddaje silnikowi komplet nastaw złożony z konfiguracji.
func (s *Silnik) ZUstawieniami(u Ustawienia) *Silnik {
	s.ustawienia = u
	return s
}

// ZDziennikiem oddaje silnikowi trwałość.
//
// Zależność opcjonalna. Silnik bez dziennika rozpoznaje mowę tak samo — traci
// wyłącznie ślad, więc brak miejsca na wiersz nie wstrzymuje transkrypcji.
func (s *Silnik) ZDziennikiem(d Dziennik) *Silnik {
	s.dziennik = d
	return s
}

// Ustawienia oddaje nastawy, którymi silnik dziś pracuje.
func (s *Silnik) Ustawienia() Ustawienia {
	return s.ustawienia
}

// katalogPracy wskazuje katalog uruchomienia pomocnika.
//
// Pusty jest odpowiedzią poprawną, nie brakiem: brama izolacji uzupełnia wtedy
// katalog własny okna (`session.SprawdzPolecenie`), a gdy punkt izolacji jest
// wyłączony — proces rusza w katalogu bieżącym rdzenia.
func (s *Silnik) katalogPracy(obszar session.Obszar) string {
	return strings.TrimSpace(obszar.KatalogRoboczy)
}

// odnotujGotowa zapisuje w dzienniku transkrypcję, która się udała.
//
// Wynik zapisu jest pomijany: transkrypcja już powstała i jest w ręku
// wołającego, więc niepowodzenie zapisu śladu nie może jej unieważnić.
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

// stanZapisu przekłada stan transkrypcji na stan wiersza dziennika.
//
// Dwa stany udane, nie jeden: nagranie bez mowy jest wynikiem, a nie odmową,
// więc dziennik ma je odróżnić wprost, zamiast kazać czytającemu wnioskować
// z `znakow = 0` (patrz stałe stanów w transkrypcja.go).
func stanZapisu(t Transkrypcja) string {
	if t.Stan == StanBezMowy {
		return StanZapisuBezMowy
	}
	return StanGotowa
}

// odnotujOdmowe zapisuje w dzienniku odmowę i oddaje ją wołającemu bez zmiany.
//
// Odmowy zapisywane są tak samo jak powodzenia, bo odmowa bez śladu jest nie do
// zdiagnozowania. Błąd wraca nietknięty: jego typ niesie rozróżnienie, po którym
// wołający nada właściwy kod kontraktu, a owinięcie go tutaj to rozróżnienie by
// zatarło.
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

// zapiszWpis dokłada identyfikator i czas, po czym oddaje wpis dziennikowi.
//
// Zlecenie bez odnośnika nagrania nie jest zapisywane: dziennik opisuje, co
// zrobiono z nagraniem, a wiersz o nagraniu, którego nie wskazano, nie odpowiada
// na żadne pytanie — i tak zostałby odrzucony przez sprawdzenie wpisu. Odmowa
// idzie wtedy do wołającego samą drogą błędu, bez wiersza.
func (s *Silnik) zapiszWpis(ctx context.Context, wpis WpisTranskrypcji) {
	if s.dziennik == nil || strings.TrimSpace(wpis.NagranieOdnosnik) == "" {
		return
	}
	wpis.Identyfikator = nowyIdentyfikator()
	wpis.Utworzono = s.teraz()
	_, _ = s.dziennik.Zapisz(ctx, wpis)
}

// nowyIdentyfikator nadaje wpisowi dziennika tożsamość widoczną poza bazą.
//
// Losowość z `crypto/rand`, a nie licznik: wpisy powstają współbieżnie w wielu
// oknach, a licznik wymagałby wspólnego stanu, którego ten pakiet nie ma.
// Błąd źródła losowości nie przerywa transkrypcji — identyfikator opada wtedy
// na znacznik czasu, bo wpis bez tożsamości nie zapisze się wcale, a to gorsza
// strata niż tożsamość mniej odporna na zbieg.
func nowyIdentyfikator() string {
	var bajty [8]byte
	if _, err := rand.Read(bajty[:]); err != nil {
		return przedrostekTranskrypcji + time.Now().UTC().Format("20060102150405.000000000")
	}
	return przedrostekTranskrypcji + hex.EncodeToString(bajty[:])
}
