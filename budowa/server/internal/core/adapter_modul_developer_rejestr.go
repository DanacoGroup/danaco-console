// Odpowiedzialność pliku: stan żywy modułu Developer — przebiegi budowania
// biegnące w tej chwili wraz z uchwytami do ich drzew procesów i narastającym
// ogonem logu. Na okno przypada jeden przebieg naraz, a log przycina rdzeń,
// nie baza.
package core

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// wierszyOgonaLogu jest liczbą wierszy zachowywanych w pamięci i w dzienniku dla każdego przebiegu budowania okna.
const wierszyOgonaLogu = 500

// najwiecejWierszyTestow jest granicą zbioru wierszy niosących wynik testu.
// Granica jest wysoka z zamysłem: repozytorium z kilkoma tysiącami testów ma
// oddać wynik każdego z nich, a jeden taki wiersz to kilkadziesiąt bajtów.
const najwiecejWierszyTestow = 20000

// przebiegBudowania jest jednym uruchomieniem zadania budowania wraz
// z uchwytami, bez których nie da się go przerwać.
type przebiegBudowania struct {
	kod       string
	oknoKod   string
	idSesji   string
	zadanie   string
	argumenty []string

	mu   sync.Mutex
	stan shared.BuildStatus
	// ogonPrzyciety mówi, czy z początku logu coś już wypadło, poza zachowaną końcówkę.
	ogonPrzyciety bool
	// wierszeTestow zbiera wiersze wyniku testu, bo ogon logu ich nie zachowa przy wielu testach.
	wierszeTestow []string
	kodWyjscia    *int
	zgloszenia    []shared.BuildProblem
	ogon          []string
	uruchomiono   time.Time
	zakonczono    time.Time

	uchwyt session.UchwytProcesu
	drzewo *session.DrzewoProcesu
	// domkniete pilnuje, żeby obserwator zakończenia zamknął kanał raz.
	domkniete atomic.Bool
	koniec    chan struct{}
}

// Dopisz dokłada wiersz do ogona logu przebiegu i rozpoznaje w nim zgłoszenie budowania albo wynik testu.
func (p *przebiegBudowania) Dopisz(wiersz string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ogon = append(p.ogon, wiersz)
	if len(p.ogon) > wierszyOgonaLogu {
		p.ogon = p.ogon[len(p.ogon)-wierszyOgonaLogu:]
		p.ogonPrzyciety = true
	}
	if czyWierszTestu(wiersz) && len(p.wierszeTestow) < najwiecejWierszyTestow {
		p.wierszeTestow = append(p.wierszeTestow, wiersz)
	}
	if zgloszenie, jest := rozpoznajZgloszenie(wiersz); jest {
		p.zgloszenia = append(p.zgloszenia, zgloszenie)
	}
}

// Ogon oddaje zachowaną końcówkę logu przebiegu budowania, złożoną w jeden napis do pokazania w oknie.
func (p *przebiegBudowania) Ogon() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return strings.Join(p.ogon, "\n")
}

// OgonPrzyciety mówi, czy zachowana końcówka logu jest krótsza od całości zapisanego przebiegu budowania.
func (p *przebiegBudowania) OgonPrzyciety() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ogonPrzyciety
}

// WierszeTestow oddaje zebrane wiersze wyniku testów i pokrycia zebrane z przebiegu budowania danego okna platformy.
func (p *przebiegBudowania) WierszeTestow() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string(nil), p.wierszeTestow...)
}

// Migawka oddaje stan przebiegu w postaci odpornej na równoległą zmianę. Stan
// czytają trzy wątki naraz: obsługiwacz komendy, pompa wyjścia i obserwator
// zakończenia.
func (p *przebiegBudowania) Migawka() (shared.BuildStatus, *int, []shared.BuildProblem, time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()
	zgloszenia := append([]shared.BuildProblem(nil), p.zgloszenia...)
	if p.kodWyjscia == nil {
		return p.stan, nil, zgloszenia, p.zakonczono
	}
	kod := *p.kodWyjscia
	return p.stan, &kod, zgloszenia, p.zakonczono
}

// Domknij zapisuje stan końcowy przebiegu. Zwraca fałsz, gdy przebieg był już
// domknięty — pierwszy prawdziwy wynik nie ma prawa zostać nadpisany przez
// późniejsze przerwanie.
func (p *przebiegBudowania) Domknij(stan shared.BuildStatus, kodWyjscia *int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stan != shared.BuildStatusRunning {
		return false
	}
	p.stan, p.kodWyjscia, p.zakonczono = stan, kodWyjscia, time.Now().UTC()
	return true
}

// Przerwij kończy całe drzewo procesu budowania. Budowanie uruchamia narzędzia,
// które uruchamiają kolejne procesy — przerwanie samego korzenia zostawiłoby
// kompilator przy życiu i przy zajętych plikach.
func (p *przebiegBudowania) Przerwij() error {
	if p.drzewo != nil {
		return p.drzewo.Ubij()
	}
	if p.uchwyt != nil {
		return p.uchwyt.Ubij()
	}
	return nil
}

// Zwolnij oddaje uchwyty systemowe po zakończeniu przebiegu budowania i zamyka kanał jego zakończenia raz.
func (p *przebiegBudowania) Zwolnij() {
	if p.drzewo != nil {
		p.drzewo.Zwolnij()
	}
	if p.domkniete.CompareAndSwap(false, true) {
		close(p.koniec)
	}
}

// rejestrBudowan trzyma przebiegi czynne jednego biegu rdzenia, po jednym przebiegu budowania na każde okno.
type rejestrBudowan struct {
	mu        sync.Mutex
	przebiegi map[string]*przebiegBudowania
}

func nowyRejestrBudowan() *rejestrBudowan {
	return &rejestrBudowan{przebiegi: make(map[string]*przebiegBudowania)}
}

// Zajmij wpisuje przebieg okna do rejestru budowań, o ile to okno nie prowadzi już innego budowania w tle.
func (r *rejestrBudowan) Zajmij(przebieg *przebiegBudowania) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, zajete := r.przebiegi[przebieg.oknoKod]; zajete {
		return false
	}
	r.przebiegi[przebieg.oknoKod] = przebieg
	return true
}

// Przebieg zwraca budowanie czynne w oknie wskazanym jego kodem, jeśli takie akurat w nim trwa teraz naprawdę.
func (r *rejestrBudowan) Przebieg(oknoKod string) (*przebiegBudowania, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	przebieg, jest := r.przebiegi[oknoKod]
	return przebieg, jest
}

// PrzebiegPoKodzie odnajduje czynny przebieg po jego identyfikatorze. Pyta o to
// odczyt logu, który zna przebieg, lecz nie zna okna, w którym on biegnie.
func (r *rejestrBudowan) PrzebiegPoKodzie(kod string) (*przebiegBudowania, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, przebieg := range r.przebiegi {
		if przebieg.kod == kod {
			return przebieg, true
		}
	}
	return nil, false
}

// Zwolnij usuwa przebieg okna, o ile to wciąż ten sam przebieg. Sprawdzenie
// tożsamości chroni przed usunięciem budowania uruchomionego zaraz po
// zakończeniu poprzedniego.
func (r *rejestrBudowan) Zwolnij(przebieg *przebiegBudowania) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if biezacy, jest := r.przebiegi[przebieg.oknoKod]; jest && biezacy == przebieg {
		delete(r.przebiegi, przebieg.oknoKod)
	}
}

// Zamknij przerywa wszystkie przebiegi budowania czynne w chwili zatrzymania rdzenia platformy Danaco.
func (r *rejestrBudowan) Zamknij() {
	if r == nil {
		return
	}
	r.mu.Lock()
	przebiegi := make([]*przebiegBudowania, 0, len(r.przebiegi))
	for _, przebieg := range r.przebiegi {
		przebiegi = append(przebiegi, przebieg)
	}
	r.przebiegi = make(map[string]*przebiegBudowania)
	r.mu.Unlock()

	for _, przebieg := range przebiegi {
		_ = przebieg.Przerwij()
	}
}
