// Odpowiedzialność pliku: pamięć okien obserwujących przebiegi automatyk —
// zaplecze pola `subscribed` komendy `automation.execution.subscribe`.
// Rejestr nie służy rozsyłaniu zdarzeń, rozwiązuje telemetrię postępu.
package core

import "sync"

// pamiecObserwatorowPrzebiegow wiąże okno z automatyką, której przebiegi
// obserwuje. Pusty kod automatyki znaczy obserwację wszystkich przebiegów.
type pamiecObserwatorowPrzebiegow struct {
	mu      sync.RWMutex
	zakresy map[string]string
}

// nowaPamiecObserwatorowPrzebiegow zakłada pusty rejestr obserwacji, gotowy
// do zapisów `Zapamietaj` i odczytów `Okna`.
func nowaPamiecObserwatorowPrzebiegow() *pamiecObserwatorowPrzebiegow {
	return &pamiecObserwatorowPrzebiegow{zakresy: map[string]string{}}
}

// Zapamietaj zapisuje obserwację okna i mówi, czy została założona. Żądanie bez
// okna nie zakłada obserwacji — i nie udaje, że założyło.
func (p *pamiecObserwatorowPrzebiegow) Zapamietaj(idOkna, kodAutomatyki string) bool {
	if p == nil || idOkna == "" {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.zakresy[idOkna] = kodAutomatyki
	return true
}

// Okna zwraca okna obserwujące wskazaną automatykę wraz z oknami obserwującymi
// wszystkie przebiegi naraz, bez rozróżnienia automatyki.
func (p *pamiecObserwatorowPrzebiegow) Okna(kodAutomatyki string) []string {
	if p == nil {
		return nil
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	okna := []string{}
	for idOkna, zakres := range p.zakresy {
		if zakres == "" || zakres == kodAutomatyki {
			okna = append(okna, idOkna)
		}
	}
	return okna
}

// przypniOknaObserwatorow wpisuje okna Execution Monitora do pamięci powiązań
// kolejki. Dzięki temu telemetria postępu zgłasza się pod oknem, które ją
// obserwuje. Kolejka bez obserwatora zostaje jak była.
func (a *adapterAutomatyk) przypniOknaObserwatorow(kolejkaID int64, kodAutomatyki string) {
	if a.kolejki == nil {
		return
	}
	okna := a.obserwatorzy.Okna(kodAutomatyki)
	if len(okna) == 0 {
		return
	}
	idSesji, _ := a.kolejki.sesjeKolejek.Odczytaj(kolejkaID)
	a.kolejki.sesjeKolejek.Zapamietaj(kolejkaID, idSesji, okna)
}
