// Plik obsługuje siedem komend rodziny aod.* nakładki Always On Display: stan
// nakładki, wysłanie wiadomości, polecenie głosowe, komplet kontekstu,
// podpowiedzi oraz przypięcie i odpięcie obserwowanego procesu.
package core

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// adapterNakladkiAod wypełnia port NakladkaAod, składając odpowiedź kontraktu
// z bytów, które zna telemetria, nadzorca sesji, port rozmowy, moduł
// Assistant i katalog akcji.
type adapterNakladkiAod struct {
	nadzorca   *session.Nadzorca
	telemetria *telemetriaPostepu
	obecnosc   *rejestrObecnosci
	rozmowa    Rozmowa
	asystent   zlecenieGlosoweAsystenta
	akcje      *RejestrAkcji
	komplety   zrodloKompletuOkna
	urzadzenia dane.RepozytoriumUrzadzen
	przypiecia *pamiecPrzypiecAod
	// wyciszenia jest magazynem wyciszeń nakładki i sygnałów klas zdarzeń,
	// trwałym w bazie danych.
	wyciszenia dane.RepozytoriumWyciszenNakladki
}

// zlecenieGlosoweAsystenta jest tą częścią modułu Assistant, której nakładka
// naprawdę używa. Zlecenie asystenta ma w rdzeniu jednego właściciela
// — moduł Assistant — i nakładka go nie dubluje.
type zlecenieGlosoweAsystenta interface {
	PolecenieGlosowe(ctx context.Context,
		z shared.AssistantVoiceCommandRequest) (shared.AssistantVoiceCommandResponse, error)
}

// zrodloKompletuOkna oddaje komplet kontekstu zapisany przy oknie.
// Wypełnia je adapter przenoszenia kontekstu — ten sam, którym jedzie
// `context.transfer`, bo komplet okna jest jeden.
type zrodloKompletuOkna interface {
	KompletOkna(idOkna string, ograniczenieHistorii int) (shared.ContextBundle, error)
}

// nowyAdapterNakladkiAod wiąże nakładkę z nadzorcą sesji, rejestrem telemetrii
// i rejestrem obecności. Pamięć przypięć powstaje od razu — przypięcie procesu
// ma działać także wtedy, gdy pozostałych źródeł nie wpięto.
func nowyAdapterNakladkiAod(nadzorca *session.Nadzorca, telemetria *telemetriaPostepu,
	obecnosc *rejestrObecnosci) *adapterNakladkiAod {

	return &adapterNakladkiAod{
		nadzorca:   nadzorca,
		telemetria: telemetria,
		obecnosc:   obecnosc,
		przypiecia: nowaPamiecPrzypiecAod(),
	}
}

// ZRozmowa wpina port rozmowy — drogę, którą `aod.chat.send` zakłada
// wiadomość w rozmowie prowadzonej przez nakładkę.
func (a *adapterNakladkiAod) ZRozmowa(r Rozmowa) *adapterNakladkiAod {
	a.rozmowa = r
	return a
}

// ZAsystentem wpina moduł Assistant — jedynego właściciela zleceń asystenta,
// z którego nakładka odczytuje polecenie głosowe.
func (a *adapterNakladkiAod) ZAsystentem(z zlecenieGlosoweAsystenta) *adapterNakladkiAod {
	a.asystent = z
	return a
}

// ZPodpowiedziami wpina katalog akcji — źródło pozycji `aod.suggestion`
// pokazywanych w nakładce jako podpowiedzi.
func (a *adapterNakladkiAod) ZPodpowiedziami(r *RejestrAkcji) *adapterNakladkiAod {
	a.akcje = r
	return a
}

// ZKontekstem wpina magazyn kompletu kontekstu okien, z którego nakładka
// odczytuje komplet dla okna żądania.
func (a *adapterNakladkiAod) ZKontekstem(z zrodloKompletuOkna) *adapterNakladkiAod {
	a.komplety = z
	return a
}

// ZUrzadzeniami wpina katalog maszyn. Bez niego nakładka nie ma jak sprawdzić
// urządzenia wskazanego przez żądanie i mówi to wprost, zamiast przyjmować
// identyfikator, którego nikt nie potwierdził.
func (a *adapterNakladkiAod) ZUrzadzeniami(r dane.RepozytoriumUrzadzen) *adapterNakladkiAod {
	a.urzadzenia = r
	return a
}

// ── aod.status.get ──────────────────────────────────────────────────────────

// StanNakladki obsługuje `aod.status.get`: oddaje stan nakładki Always On
// Display — urządzenie, kartę i okno pokazywane w nakładce, procesy
// przypięte do obserwacji i liczbę procesów w biegu.
func (a *adapterNakladkiAod) StanNakladki(ctx context.Context,
	z shared.AodStatusGetRequest) (shared.AodStatusGetResponse, error) {

	urzadzenie, err := a.urzadzenieNakladki(ctx, z.DeviceId)
	if err != nil {
		return shared.AodStatusGetResponse{}, err
	}
	if a.telemetria == nil {
		return shared.AodStatusGetResponse{}, bladBrakuSkladnikaNakladki(
			"rejestr telemetrii postępu")
	}
	stan := shared.AodStatus{
		DeviceId:            tekstOpcjonalny(urzadzenie),
		AttachedProcessIds:  a.przypiecia.Wykaz(urzadzenie),
		RunningProcessCount: a.procesyWBiegu(),
		UpdatedAt:           time.Now().UTC().UnixMilli(),
	}
	// Okno i karta w nakładce biorą się z ostatniego punktu pracy; brak
	// pracy zostawia pola puste.
	if okno, jest := a.oknoOstatniejCzynnosci(); jest {
		stan.ActiveWindowId = tekstOpcjonalny(okno.Id)
		stan.ActiveSessionId = tekstOpcjonalny(okno.IdSesji)
		// Moduł idzie tym samym oknem, co karta sesji, zgodnie z bieżącym
		// punktem pracy.
		stan.ModuleId = tekstOpcjonalny(okno.Modul)
	}
	return shared.AodStatusGetResponse{Status: stan}, nil
}

// procesyWBiegu liczy procesy telemetrii o stanie `running`. Kontrakt pyta
// o procesy w biegu, więc `pending`, `paused` i stany końcowe do liczby nie
// wchodzą.
func (a *adapterNakladkiAod) procesyWBiegu() int {
	biegnace := 0
	for _, odpis := range a.telemetria.Odpisy() {
		if odpis.Stan == shared.ProgressStatusRunning {
			biegnace++
		}
	}
	return biegnace
}

// ── aod.observe.attach ──────────────────────────────────────────────────────

// PrzypnijObserwacje obsługuje `aod.observe.attach`: przypina proces do
// obserwacji w nakładce wskazanego urządzenia, odmawiając kodem `not_found`,
// gdy rejestr telemetrii procesu nie zna.
func (a *adapterNakladkiAod) PrzypnijObserwacje(ctx context.Context,
	z shared.AodObserveAttachRequest) (shared.AodObserveAttachResponse, error) {

	urzadzenie, err := a.urzadzenieNakladki(ctx, z.DeviceId)
	if err != nil {
		return shared.AodObserveAttachResponse{}, err
	}
	if z.ProcessId == "" {
		return shared.AodObserveAttachResponse{}, bladZadaniaNakladki(
			"przypięcie obserwacji bez wskazania procesu")
	}
	if a.telemetria == nil {
		return shared.AodObserveAttachResponse{}, bladBrakuSkladnikaNakladki(
			"rejestr telemetrii postępu")
	}
	if !a.czyProcesZnany(z.ProcessId) {
		return shared.AodObserveAttachResponse{}, bladBrakuProcesuTelemetrii(z.ProcessId)
	}
	a.przypiecia.przypnij(urzadzenie, z.ProcessId)
	return shared.AodObserveAttachResponse{
		AttachedProcessIds: a.przypiecia.Wykaz(urzadzenie),
	}, nil
}

// ── aod.observe.detach ──────────────────────────────────────────────────────

// OdepnijObserwacje obsługuje `aod.observe.detach`: zdejmuje proces z
// obserwacji w nakładce wskazanego urządzenia, odmawiając kodem `not_found`,
// gdy przypięcia nie było.
func (a *adapterNakladkiAod) OdepnijObserwacje(ctx context.Context,
	z shared.AodObserveDetachRequest) (shared.AodObserveDetachResponse, error) {

	urzadzenie, err := a.urzadzenieNakladki(ctx, z.DeviceId)
	if err != nil {
		return shared.AodObserveDetachResponse{}, err
	}
	if z.ProcessId == "" {
		return shared.AodObserveDetachResponse{}, bladZadaniaNakladki(
			"odpięcie obserwacji bez wskazania procesu")
	}
	if !a.przypiecia.odepnij(urzadzenie, z.ProcessId) {
		return shared.AodObserveDetachResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "nakładka AOD: proces "+z.ProcessId+
				" nie jest przypięty do obserwacji "+opisNakladki(urzadzenie)))
	}
	return shared.AodObserveDetachResponse{
		AttachedProcessIds: a.przypiecia.Wykaz(urzadzenie),
	}, nil
}

// czyProcesZnany mówi, czy rejestr telemetrii zna proces o tym
// identyfikatorze, rozstrzygając, czy przypięcie wskazuje byt istniejący.
func (a *adapterNakladkiAod) czyProcesZnany(idProcesu string) bool {
	for _, odpis := range a.telemetria.Odpisy() {
		if odpis.Id == idProcesu {
			return true
		}
	}
	return false
}

// ── ustalenia wspólne żądań ─────────────────────────────────────────────────

// urzadzenieNakladki sprawdza urządzenie wskazane w żądaniu i oddaje klucz
// nakładki; żądanie bez urządzenia daje klucz pusty zamiast odmowy.
func (a *adapterNakladkiAod) urzadzenieNakladki(ctx context.Context,
	wskazane *string) (string, error) {

	kod := wartoscTekstu(wskazane)
	if kod == "" {
		return "", nil
	}
	if a.urzadzenia == nil {
		return "", bladBrakuSkladnikaNakladki("katalog urządzeń")
	}
	numer, err := strconv.ParseInt(kod, 10, 64)
	if err != nil || numer <= 0 {
		return "", bladZadaniaNakladki("identyfikator urządzenia " + strconv.Quote(kod) +
			" nie jest identyfikatorem katalogu urządzeń")
	}
	if _, err := a.urzadzenia.Pobierz(ctx, numer); errors.Is(err, dane.ErrBrakWiersza) {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"nakładka AOD: urządzenie "+kod+" nie istnieje w katalogu maszyn"))
	} else if err != nil {
		return "", protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
	}
	return kod, nil
}

// oknoZadania ustala okno, którego dotyczy żądanie nakładki: wskazane
// wprost, okno karty sesji albo — gdy żądanie nie wskazało niczego — okno
// ostatniego punktu pracy odnotowanego przez telemetrię.
func (a *adapterNakladkiAod) oknoZadania(idOkna, idSesji *string) (session.Okno, error) {
	if a.nadzorca == nil {
		return session.Okno{}, bladBrakuSkladnikaNakladki("nadzorca sesji i okien")
	}
	if id := wartoscTekstu(idOkna); id != "" {
		okno, err := a.nadzorca.Rejestr().Okno(id)
		if err != nil {
			return session.Okno{}, bladSesji(err)
		}
		return okno, nil
	}
	if id := wartoscTekstu(idSesji); id != "" {
		return a.oknoKartySesji(id)
	}
	okno, jest := a.oknoOstatniejCzynnosci()
	if !jest {
		return session.Okno{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"nakładka AOD: żądanie nie wskazało ani okna, ani karty sesji, "+
				"a rdzeń nie odnotował dotąd pracy w żadnym otwartym oknie"))
	}
	return okno, nil
}

// oknoKartySesji wskazuje okno karty sesji: to, w którym praca szła ostatnio,
// a bez takiego zapisu — pierwsze okno otwarte. Karta bez otwartego okna nie ma
// czego pokazać nakładce i mówi to odmową, nie pustką.
func (a *adapterNakladkiAod) oknoKartySesji(idSesji string) (session.Okno, error) {
	if _, err := a.nadzorca.Rejestr().Sesja(idSesji); err != nil {
		return session.Okno{}, bladSesji(err)
	}
	okna, err := a.nadzorca.Rejestr().OknaSesji(idSesji)
	if err != nil {
		return session.Okno{}, bladSesji(err)
	}
	ostatnie := ""
	if a.obecnosc != nil {
		if wpis, jest := a.obecnosc.czynnosc.Sesja(idSesji); jest {
			ostatnie = wpis.IdOkna
		}
	}
	var wybrane session.Okno
	for _, okno := range okna {
		if !okno.CzyOtwarte() {
			continue
		}
		if wybrane.Id == "" || okno.Id == ostatnie {
			wybrane = okno
		}
	}
	if wybrane.Id == "" {
		return session.Okno{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"nakładka AOD: karta sesji "+idSesji+" nie ma otwartego okna komunikacji"))
	}
	return wybrane, nil
}

// oknoOstatniejCzynnosci oddaje okno ostatniego punktu pracy, o ile nadal jest
// otwarte. Ślad okna już zamkniętego nie opisuje niczego, więc nie zastępuje
// odpowiedzi.
func (a *adapterNakladkiAod) oknoOstatniejCzynnosci() (session.Okno, bool) {
	if a.obecnosc == nil || a.nadzorca == nil {
		return session.Okno{}, false
	}
	wpis, jest := a.obecnosc.czynnosc.Ostatnie()
	if !jest {
		return session.Okno{}, false
	}
	okno, err := a.nadzorca.Rejestr().Okno(wpis.IdOkna)
	if err != nil || !okno.CzyOtwarte() {
		return session.Okno{}, false
	}
	return okno, true
}

// Ostatnie oddaje najświeższy punkt pracy spośród wszystkich okien. Bliźniak
// metod `Okno` i `Sesja`, osobny dlatego, że nakładka pyta nie o wskazany byt,
// lecz o to, gdzie praca szła ostatnio.
func (p *pamiecCzynnosci) Ostatnie() (czynnoscOkna, bool) {
	if p == nil {
		return czynnoscOkna{}, false
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	var najswiezszy czynnoscOkna
	for _, wpis := range p.okna {
		// Remis rozstrzyga identyfikator okna, żeby odpowiedź nie zmieniała
		// się między odczytami.
		if najswiezszy.IdOkna == "" || wpis.Chwila.After(najswiezszy.Chwila) ||
			(wpis.Chwila.Equal(najswiezszy.Chwila) && wpis.IdOkna < najswiezszy.IdOkna) {
			najswiezszy = wpis
		}
	}
	return najswiezszy, najswiezszy.IdOkna != ""
}

// ── pamięć przypięć nakładki ────────────────────────────────────────────────

// pamiecPrzypiecAod wiąże nakładkę urządzenia z procesami przypiętymi do
// obserwacji. Klucz pusty jest nakładką żądania bez urządzenia — jedną, tak
// jak jeden jest Operator.
type pamiecPrzypiecAod struct {
	mu       sync.RWMutex
	nakladki map[string]map[string]struct{}
}

// nowaPamiecPrzypiecAod zakłada pustą pamięć przypięć, gotową do przyjęcia
// pierwszego przypięcia procesu.
func nowaPamiecPrzypiecAod() *pamiecPrzypiecAod {
	return &pamiecPrzypiecAod{nakladki: map[string]map[string]struct{}{}}
}

// przypnij dopisuje proces do obserwacji nakładki. Przypięcie powtórzone nie
// zmienia niczego — wykaz jest zbiorem, nie dziennikiem.
func (p *pamiecPrzypiecAod) przypnij(nakladka, idProcesu string) {
	if p == nil || idProcesu == "" {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.nakladki[nakladka] == nil {
		p.nakladki[nakladka] = map[string]struct{}{}
	}
	p.nakladki[nakladka][idProcesu] = struct{}{}
}

// odepnij zdejmuje proces z obserwacji nakładki i mówi, czy było co
// zdejmować, żeby wołający rozróżnił zmianę od braku zmiany.
func (p *pamiecPrzypiecAod) odepnij(nakladka, idProcesu string) bool {
	if p == nil {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	procesy, jest := p.nakladki[nakladka]
	if !jest {
		return false
	}
	if _, przypiety := procesy[idProcesu]; !przypiety {
		return false
	}
	delete(procesy, idProcesu)
	if len(procesy) == 0 {
		delete(p.nakladki, nakladka)
	}
	return true
}

// Wykaz oddaje procesy przypięte do nakładki, uporządkowane po identyfikatorze
// — nakładka ma wyglądać tak samo przy każdym odczycie. Wykaz jest zawsze
// tablicą, nigdy brakiem wartości: pusta obserwacja to fakt, nie cisza.
func (p *pamiecPrzypiecAod) Wykaz(nakladka string) []string {
	procesy := []string{}
	if p == nil {
		return procesy
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	for idProcesu := range p.nakladki[nakladka] {
		procesy = append(procesy, idProcesu)
	}
	sort.Strings(procesy)
	return procesy
}

// ── odmowy ──────────────────────────────────────────────────────────────────

// bladZadaniaNakladki nazywa żądanie niezgodne z kontraktem — błąd
// Operatora, nie rdzenia, i wraca kodem walidacji.
func bladZadaniaNakladki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"nakładka AOD: "+powod))
}

// bladBrakuSkladnikaNakladki nazywa składnik rdzenia, którego nie wpięto.
// Odpowiedź pusta byłaby tu ciszą udającą wynik.
func bladBrakuSkladnikaNakladki(byt string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"nakładka AOD: "+byt+" niewpięty — odpowiedzi nie ma skąd złożyć"))
}

// opisNakladki nazywa nakładkę w treści odmowy: urządzenia albo tę jedyną,
// gdy żądanie urządzenia nie wskazało.
func opisNakladki(nakladka string) string {
	if nakladka == "" {
		return "nakładki bez wskazanego urządzenia"
	}
	return "nakładki urządzenia " + nakladka
}
