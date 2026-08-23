package core

import (
	"sync"
	"time"

	"danacoconsole/shared"
)

// pamiecCzynnosci pamięta, co ostatnio działo się w oknach sesji.
//
// Telemetria postępu mówi o procesie i jego etapie, ale niczego nie
// pamięta między zdarzeniami i nie wie, że okna składają się na sesję.
// Kontrolka powrotu pyta odwrotnie: która sesja żyje, które jej okno pracowało
// ostatnio i kiedy. Ta pamięć jest przejściem między jednym a drugim.
//
// Trzyma wyłącznie stan pracy i znacznik czasu. Sesji, okien ani procesów nie
// przechowuje — mają własne rejestry.
type pamiecCzynnosci struct {
	mu    sync.RWMutex
	okna  map[string]czynnoscOkna
	sesje map[string]czynnoscOkna
}

// czynnoscOkna jest ostatnim punktem pracy jednego okna.
type czynnoscOkna struct {
	// IdOkna, którego dotyczy punkt pracy.
	IdOkna string
	// IdSesji okna; puste, gdy telemetria nie zna sesji.
	IdSesji string
	// Stan procesu okna wzięty wprost z telemetrii.
	Stan shared.ProgressStatus
	// Chwila zgłoszenia.
	Chwila time.Time
}

// nowaPamiecCzynnosci zakłada pustą pamięć.
func nowaPamiecCzynnosci() *pamiecCzynnosci {
	return &pamiecCzynnosci{
		okna:  map[string]czynnoscOkna{},
		sesje: map[string]czynnoscOkna{},
	}
}

// Odnotuj zapisuje punkt pracy okna i odpowiada, czy zmiana jest widoczna dla
// Operatora.
//
// Prawdę zwraca wyłącznie zmiana stanu pracy — start tury, jej domknięcie,
// zatrzymanie albo niepowodzenie. Kolejny etap tej samej tury przesuwa
// wyłącznie znacznik czasu, bo inaczej każdy fragment odpowiedzi modelu
// rozgłaszałby zdarzenie o sesji, w której nic się nie zmieniło.
func (p *pamiecCzynnosci) Odnotuj(idSesji, idOkna string, stan shared.ProgressStatus) bool {
	if p == nil || idOkna == "" {
		return false
	}
	wpis := czynnoscOkna{IdOkna: idOkna, IdSesji: idSesji, Stan: stan, Chwila: time.Now().UTC()}

	p.mu.Lock()
	defer p.mu.Unlock()
	poprzedni, znane := p.okna[idOkna]
	if wpis.IdSesji == "" {
		wpis.IdSesji = poprzedni.IdSesji
	}
	p.okna[idOkna] = wpis
	if wpis.IdSesji != "" {
		p.sesje[wpis.IdSesji] = wpis
	}
	return !znane || poprzedni.Stan != wpis.Stan
}

// Sesja zwraca ostatni punkt pracy w sesji: okno, jego stan i chwilę. Sesja bez
// zapisanej czynności daje fałsz — odpis obecności bierze wtedy znacznik zmiany
// samej sesji.
func (p *pamiecCzynnosci) Sesja(idSesji string) (czynnoscOkna, bool) {
	if p == nil || idSesji == "" {
		return czynnoscOkna{}, false
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	wpis, jest := p.sesje[idSesji]
	return wpis, jest
}

// Zachowaj zostawia wyłącznie okna wciąż otwarte i wykreśla ślad pozostałych
// wraz z osieroconą czynnością ich sesji. Wykaz pusty niczego nie kasuje: brak
// wiedzy o oknach nie jest wiedzą o ich zamknięciu.
func (p *pamiecCzynnosci) Zachowaj(okna map[string]struct{}) {
	if p == nil || len(okna) == 0 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for idOkna := range p.okna {
		if _, zyje := okna[idOkna]; !zyje {
			delete(p.okna, idOkna)
		}
	}
	for idSesji, wpis := range p.sesje {
		if _, zyje := okna[wpis.IdOkna]; !zyje {
			delete(p.sesje, idSesji)
		}
	}
}
