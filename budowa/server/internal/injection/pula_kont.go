package injection

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Konto jest jednym profilem uwierzytelnienia kanału. W bazie i konfiguracji
// żyje wyłącznie ODWOŁANIE — kod konta i ścieżka katalogu; poświadczenia
// zostają w profilu na dysku i nigdy nie przechodzą przez rdzeń.
type Konto struct {
	Kod                 string
	KatalogKonfiguracji string
	// WyczerpaneDo niesie chwilę odnowienia limitu odczytaną z trwałego zapisu
	// katalogu (kolumny konto.stan / konto.wyczerpane_do). Zero znaczy „konto
	// nie było wyczerpane". Dzięki temu wyczerpanie przeżywa restart rdzenia:
	// pula odtwarza pamięć limitu z bazy, zamiast zaczynać od czystej mapy.
	WyczerpaneDo time.Time
}

// PulaKont trzyma kolejność kont i pamięć o wyczerpanych limitach.
// Wyczerpanie jednego konta nie kończy sesji: pula podaje następne,
// a gdy wszystkie są wyczerpane, mówi to wprost zamiast blokować.
type PulaKont struct {
	mu         sync.Mutex
	konta      []Konto
	biezace    int
	wyczerpane map[string]time.Time
	teraz      func() time.Time
	// utrwal zapisuje wyczerpanie konta w trwałym katalogu, żeby przeżyło
	// restart. Nil znaczy pulę bez trwałości (np. z katalogu profili na dysku).
	// Wołany poza zamkiem — zapis do bazy nie może blokować rotacji.
	utrwal func(kod string, doChwili time.Time)
	// zrodlo podaje bieżący wykaz kont z katalogu. Pula sięga po nie na progu
	// tury, dzięki czemu konto dodane komendą account.* wchodzi do rotacji bez
	// restartu (analogicznie do Odswiez rejestru kanałów). Nil znaczy pulę
	// nieodświeżalną.
	zrodlo func() ([]Konto, bool)
}

// NowaPula składa pulę z podanych kont, pomijając wpisy niekompletne. Wpis
// z niezerowym WyczerpaneDo zasila pamięć wyczerpania — pula odtwarza stan
// limitu zapisany w katalogu.
func NowaPula(konta ...Konto) *PulaKont {
	pula := &PulaKont{wyczerpane: map[string]time.Time{}, teraz: time.Now}
	for _, konto := range konta {
		if strings.TrimSpace(konto.Kod) == "" || strings.TrimSpace(konto.KatalogKonfiguracji) == "" {
			continue
		}
		pula.konta = append(pula.konta, konto)
		if !konto.WyczerpaneDo.IsZero() {
			pula.wyczerpane[konto.Kod] = konto.WyczerpaneDo
		}
	}
	return pula
}

// UstawUtrwalanie wpina zapis wyczerpania do trwałego katalogu. Bez niego pula
// pamięta wyczerpanie wyłącznie do restartu.
func (p *PulaKont) UstawUtrwalanie(utrwal func(kod string, doChwili time.Time)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.utrwal = utrwal
}

// UstawZrodlo wpina odczyt bieżącego wykazu kont. Źródło zwraca wykaz i znacznik
// powodzenia: przy niepowodzeniu odczytu (drugi wynik false) pula zachowuje stan
// dotychczasowy zamiast się opróżniać.
func (p *PulaKont) UstawZrodlo(zrodlo func() ([]Konto, bool)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.zrodlo = zrodlo
}

// OdswiezZeZrodla przebudowuje wykaz kont ze źródła, jeśli je wpięto. Wołane na
// progu tury — konto dodane albo usunięte komendą account.* wchodzi do rotacji
// bez restartu. Pamięć wyczerpania kont, które przetrwały, zostaje zachowana.
func (p *PulaKont) OdswiezZeZrodla() {
	p.mu.Lock()
	zrodlo := p.zrodlo
	p.mu.Unlock()
	if zrodlo == nil {
		return
	}
	if konta, ok := zrodlo(); ok {
		p.przeladuj(konta)
	}
}

// przeladuj zastępuje wykaz kont nowym, zachowując pamięć wyczerpania dla kont,
// które nadal istnieją, i dokładając wyczerpanie zapisane w katalogu. Wpisy
// niekompletne są pomijane, tak jak w NowaPula.
func (p *PulaKont) przeladuj(konta []Konto) {
	p.mu.Lock()
	defer p.mu.Unlock()
	stare := p.wyczerpane
	p.konta = p.konta[:0]
	nowe := map[string]time.Time{}
	obecne := map[string]bool{}
	for _, konto := range konta {
		if strings.TrimSpace(konto.Kod) == "" || strings.TrimSpace(konto.KatalogKonfiguracji) == "" {
			continue
		}
		p.konta = append(p.konta, konto)
		obecne[konto.Kod] = true
		if !konto.WyczerpaneDo.IsZero() {
			nowe[konto.Kod] = konto.WyczerpaneDo
		}
	}
	// Pamięć sesji jest źródłem świeższym niż katalog dla kont wciąż obecnych:
	// wyczerpanie rozpoznane w tej turze zapisało się już do bazy, ale odczyt
	// mógł je wyprzedzić. Zachowujemy późniejszą z dwóch chwil.
	for kod, chwila := range stare {
		if !obecne[kod] {
			continue
		}
		if istniejaca, jest := nowe[kod]; !jest || chwila.After(istniejaca) {
			nowe[kod] = chwila
		}
	}
	p.wyczerpane = nowe
	if p.biezace >= len(p.konta) {
		p.biezace = 0
	}
}

// KontaZKatalogu odczytuje profile z katalogu wskazanego konfiguracją.
// Podkatalogi o nazwie zaczynającej się od kropki albo podkreślenia są
// zapleczem narzędzia, nie kontami.
func KontaZKatalogu(katalog string) ([]Konto, error) {
	wpisy, err := os.ReadDir(katalog)
	if err != nil {
		return nil, fmt.Errorf("injection: odczyt katalogu profili %s: %w", katalog, err)
	}
	var konta []Konto
	for _, wpis := range wpisy {
		nazwa := wpis.Name()
		if !wpis.IsDir() || strings.HasPrefix(nazwa, "_") || strings.HasPrefix(nazwa, ".") {
			continue
		}
		konta = append(konta, Konto{Kod: nazwa, KatalogKonfiguracji: filepath.Join(katalog, nazwa)})
	}
	sort.Slice(konta, func(i, j int) bool { return konta[i].Kod < konta[j].Kod })
	return konta, nil
}

// Konta zwraca kopię wykazu kont w kolejności rotacji.
func (p *PulaKont) Konta() []Konto {
	p.mu.Lock()
	defer p.mu.Unlock()
	kopia := make([]Konto, len(p.konta))
	copy(kopia, p.konta)
	return kopia
}

// PoKodzie zwraca konto o wskazanym kodzie, nie przestawiając rotacji.
// Droga dla wskazania konta per okno: tura wskazana jedzie dokładnie
// tą tożsamością, a wskaźnik `biezace` rotacji pozostaje nietknięty — dwa
// okna na dwóch kontach nie przestawiają sobie nawzajem puli.
func (p *PulaKont) PoKodzie(kod string) (Konto, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, konto := range p.konta {
		if konto.Kod == kod {
			return konto, true
		}
	}
	return Konto{}, false
}

// Pusta mówi, czy pula nie ma ani jednego konta.
//
// Pula pusta to NIE to samo co pula wyczerpana. Brak kont oznacza, że Operator
// nie wskazał żadnej tożsamości — a wtedy program `claude` ma użyć tożsamości
// otoczenia (własnego logowania na maszynie), bo CLAUDE_CONFIG_DIR jest
// wyłącznie nośnikiem tożsamości, nie warunkiem uruchomienia.
// Odmowa w takiej sytuacji łamie fail-open.
func (p *PulaKont) Pusta() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.konta) == 0
}

// Biezace zwraca konto do użycia. Gdy bieżące jest wyczerpane, pula sama
// przechodzi na kolejne dostępne — wywołujący nie musi o tym wiedzieć.
func (p *PulaKont) Biezace() (Konto, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.dostepneOd(p.biezace)
}

// Wyczerpane odnotowuje wyczerpanie konta i przełącza pulę na następne.
// Zwraca konto, którym praca toczy się dalej. Sesja biegnie nieprzerwanie:
// zmiana konta jest zdarzeniem kanału, nie końcem rozmowy.
func (p *PulaKont) Wyczerpane(kod string, doChwili time.Time) (Konto, bool) {
	p.mu.Lock()
	if doChwili.IsZero() {
		doChwili = p.teraz().Add(time.Hour)
	}
	p.wyczerpane[kod] = doChwili
	for i, konto := range p.konta {
		if konto.Kod == kod {
			p.biezace = (i + 1) % len(p.konta)
			break
		}
	}
	konto, dostepne := p.dostepneOd(p.biezace)
	utrwal := p.utrwal
	p.mu.Unlock()
	// Trwały ślad wyczerpania idzie poza zamkiem: zapis do katalogu (baza) nie
	// może wstrzymywać kolejnych decyzji rotacji. Brak utrwalacza znaczy pulę
	// bez trwałości — pamięć limitu żyje wtedy do restartu.
	if utrwal != nil {
		utrwal(kod, doChwili)
	}
	return konto, dostepne
}

// dostepneOd szuka pierwszego konta bez czynnego wyczerpania, zaczynając od
// wskazanego miejsca. Wywoływane pod zamkniętym zamkiem.
func (p *PulaKont) dostepneOd(od int) (Konto, bool) {
	if len(p.konta) == 0 {
		return Konto{}, false
	}
	teraz := p.teraz()
	for krok := 0; krok < len(p.konta); krok++ {
		i := (od + krok) % len(p.konta)
		konto := p.konta[i]
		doChwili, wyczerpane := p.wyczerpane[konto.Kod]
		if wyczerpane && teraz.Before(doChwili) {
			continue
		}
		if wyczerpane {
			delete(p.wyczerpane, konto.Kod)
		}
		p.biezace = i
		return konto, true
	}
	return Konto{}, false
}
