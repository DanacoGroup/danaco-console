package core

import (
	"sync"

	"danacoconsole/shared"
)

// telemetriaPostepu jest jedynym producentem zdarzenia progress.changed. Producent nie zna ani tury modelu, ani kolejki, wyłącznie proces wskazany kluczem, więc jeden byt telemetryczny obsługuje wszystkie źródła postępu bez drugiej implementacji.
type telemetriaPostepu struct {
	emiter *emiter

	mu      sync.Mutex
	procesy map[string]*procesPostepu
}

// przedrostekProcesu znakuje identyfikator procesu telemetrii. Proces jest
// bytem rdzenia, nie odpowiedzią na żądanie, więc niesie własny identyfikator.
const przedrostekProcesu = "proc-"

// Nazwy etapów wychodzące polem stepLabel kontraktu. Są opisem dla Operatora,
// więc brzmią po polsku i mówią o rzeczywistym punkcie pracy.
const (
	etapPrzyjecieTury      = "przyjęcie wiadomości"
	etapStartTury          = "start tury"
	etapRozumowanie        = "tok rozumowania"
	etapFragment           = "fragment odpowiedzi"
	etapNarzedzie          = "wywołanie narzędzia"
	etapWynikNarzedzia     = "wynik narzędzia"
	etapKoniecTury         = "koniec tury"
	etapZatrzymanie        = "zatrzymanie tury"
	etapZalozenieKolejki   = "założenie kolejki"
	etapBiegNaprawczy      = "bieg naprawczy"
	etapDzialanieNaKolejce = "działanie na kolejce"
)

// opisProcesu wskazuje, czego dotyczy zgłoszenie. Klucz rozróżnia procesy:
// tura okna idzie pod identyfikatorem okna, kolejka pod własnym kluczem.
// Pola puste i zerowe nie nadpisują tego, co telemetria już wie.
type opisProcesu struct {
	Klucz   string
	IdOkna  string
	IdSesji string
	Etap    int
	Etapow  int
}

// nowaTelemetriePostepu zakłada producenta nad nadajnikiem transportu, gotowego do otwierania i przesuwania procesów.
func nowaTelemetriePostepu(nadajnik Nadajnik) *telemetriaPostepu {
	return &telemetriaPostepu{emiter: nowyEmiter(nadajnik), procesy: map[string]*procesPostepu{}}
}

// Zacznij otwiera nowy proces i rozgłasza jego pierwszy etap. Poprzedni proces
// tego samego klucza ustępuje — jedno okno prowadzi jedną turę naraz.
func (t *telemetriaPostepu) Zacznij(o opisProcesu, nazwaEtapu string) {
	if t == nil || o.Klucz == "" {
		return
	}
	t.mu.Lock()
	proces := t.otworz(o)
	proces.Stan = shared.ProgressStatusRunning
	proces.Nazwa = nazwaEtapu
	nanies(proces, o)
	odpis := *proces
	t.mu.Unlock()
	t.rozglos(odpis)
}

// Krok przesuwa proces o jeden etap. Zgłoszenie kroku dla procesu nieznanego
// albo domkniętego otwiera proces nowy — telemetria zaczyna wtedy liczyć od
// zastanego punktu, zamiast milczeć.
func (t *telemetriaPostepu) Krok(o opisProcesu, nazwaEtapu string) {
	if t == nil || o.Klucz == "" {
		return
	}
	t.mu.Lock()
	proces, nowy := t.czynny(o)
	if !nowy && o.Etap == 0 {
		proces.Etap++
	}
	proces.Stan = shared.ProgressStatusRunning
	proces.Nazwa = nazwaEtapu
	nanies(proces, o)
	odpis := *proces
	t.mu.Unlock()
	t.rozglos(odpis)
}

// Stan nadaje procesowi stan wskazany wprost: wstrzymanie, zatrzymanie, domknięcie albo niepowodzenie, przesuwając licznik tak jak Krok, chyba że zgłaszający podał etap sam. Powtórne domknięcie już zamkniętego procesu nie rozgłasza niczego.
func (t *telemetriaPostepu) Stan(o opisProcesu, stan shared.ProgressStatus, nazwaEtapu string) {
	if t == nil || o.Klucz == "" {
		return
	}
	t.mu.Lock()
	if zastany, jest := t.procesy[o.Klucz]; jest && zastany.Zamkniety && czyStanKoncowy(stan) {
		t.mu.Unlock()
		return
	}
	proces, nowy := t.czynny(o)
	if !nowy && o.Etap == 0 {
		proces.Etap++
	}
	proces.Stan = stan
	proces.Nazwa = nazwaEtapu
	proces.Zamkniety = czyStanKoncowy(stan)
	nanies(proces, o)
	odpis := *proces
	t.mu.Unlock()
	t.rozglos(odpis)
}

// czynny zwraca proces gotowy do zmiany: zastany i niedomknięty albo świeżo
// otwarty. Druga wartość mówi, czy proces powstał w tym wywołaniu.
// Wywoływać wyłącznie pod zamkiem.
func (t *telemetriaPostepu) czynny(o opisProcesu) (*procesPostepu, bool) {
	if proces, jest := t.procesy[o.Klucz]; jest && !proces.Zamkniety {
		return proces, false
	}
	return t.otworz(o), true
}

// otworz zakłada proces pod kluczem opisu, nadając mu nowy identyfikator. Wywoływać wyłącznie pod zamkiem.
func (t *telemetriaPostepu) otworz(o opisProcesu) *procesPostepu {
	proces := &procesPostepu{
		Id:      nowyIdentyfikator(przedrostekProcesu),
		IdOkna:  o.IdOkna,
		IdSesji: o.IdSesji,
		Etap:    1,
		Stan:    shared.ProgressStatusPending,
	}
	t.procesy[o.Klucz] = proces
	return proces
}

// rozglos oddaje ładunek telemetrii emiterowi zdarzeń, wysyłając go do kanału sesji, której dotyczy proces.
func (t *telemetriaPostepu) rozglos(p procesPostepu) {
	t.emiter.postep(p.IdSesji, p.ladunek())
}
