// Odpowiedzialność pliku: pamięć okien obserwujących przebiegi automatyk —
// zaplecze pola `subscribed` komendy `automation.execution.subscribe`.
//
// Rejestr nie służy rozsyłaniu zdarzeń: `automation.execution.status` dociera
// do wszystkich połączeń konta i rejestr niczego w tym nie zmienia. Rozwiązuje
// inną rzecz — telemetrię postępu. Proces kolejki przypina się do okna
// (`opisProcesu.IdOkna`), a kolejka wykonująca automatykę powstaje ze strony
// głównej, więc nie ma okna rozmowy, pod którym miałaby się zgłaszać. Oknem,
// które tę pracę obserwuje, jest Execution Monitor; rejestr jest jedynym
// miejscem, z którego rdzeń może się tego dowiedzieć.
//
// Dlatego `subscribed` mówi prawdę: zapisanie okna jest czynnością o skutku,
// a komenda wywołana bez `windowId` jest zwykłym odczytem i oddaje `false`.
package core

import "sync"

// pamiecObserwatorowPrzebiegow wiąże okno z automatyką, której przebiegi
// obserwuje. Pusty kod automatyki znaczy obserwację wszystkich przebiegów.
type pamiecObserwatorowPrzebiegow struct {
	mu      sync.RWMutex
	zakresy map[string]string
}

// nowaPamiecObserwatorowPrzebiegow zakłada pusty rejestr obserwacji.
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
// wszystkie przebiegi.
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
// kolejki. Dzięki temu telemetria postępu kolejki wykonującej automatykę
// zgłasza się pod oknem, które ją obserwuje, a `Queue.windowIds` nazywa je
// wprost. Kolejka bez obserwatora zostaje jak była — powiązanie jest dodatkiem,
// nie warunkiem.
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
