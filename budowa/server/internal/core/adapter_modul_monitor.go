// Plik obsługuje rodzinę komend `monitor.*` na rejestrze procesów rdzenia
// trzymanym w pamięci: odczyt stanu bieżącego (`monitor.status`) oraz zapis
// okna na telemetrię procesów (`monitor.subscribe`).
package core

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// adapterMonitora wypełnia port monitorProcesow: adapter nawigacji rozszerzony
// o rejestr telemetrii postępu i pamięć okien obserwujących procesy.
type adapterMonitora struct {
	*adapterNawigacji
	// telemetria jest jedynym rejestrem procesów; zerowa znaczy odmowę,
	// nie pusty wykaz.
	telemetria *telemetriaPostepu
	obserwacje *pamiecObserwatorowProcesow
	// zegar pamięta chwilę stwierdzenia stanu procesu, którego chwili
	// zmiany nie zapisał nikt inny.
	zegar *zegarOdpisowProcesow
}

// ZTelemetriaProcesow rozszerza adapter nawigacji o rodzinę `monitor.*`.
// Wywoływać PO `ZObecnoscia`: chwilę ostatniego zgłoszenia telemetrii zna
// pamięć czynności rejestru obecności, a licznik obiegów — rejestr biegów.
func (a *adapterNawigacji) ZTelemetriaProcesow(telemetria *telemetriaPostepu) *adapterMonitora {
	return &adapterMonitora{
		adapterNawigacji: a,
		telemetria:       telemetria,
		obserwacje:       nowaPamiecObserwatorowProcesow(),
		zegar:            nowyZegarOdpisowProcesow(),
	}
}

// ── monitor.status ──────────────────────────────────────────────────────────

// StanProcesow obsługuje `monitor.status`: oddaje stan bieżący procesów
// spełniających warunki żądania, bez zakładania obserwacji. Proces wskazany
// wprost, którego rejestr nie zna, wraca odmową `not_found` z jego nazwą.
func (a *adapterMonitora) StanProcesow(ctx context.Context,
	z shared.MonitorStatusRequest) (shared.MonitorStatusResponse, error) {

	odpisy, err := a.odpisyProcesow()
	if err != nil {
		return shared.MonitorStatusResponse{}, err
	}
	idProcesu := wartoscTekstu(z.ProcessId)
	sito := sitoProcesow{
		idProcesu: idProcesu,
		idOkna:    wartoscTekstu(z.WindowId),
		idSesji:   wartoscTekstu(z.SessionId),
	}
	stany := a.stanyKontraktu(odpisy, sito, 0)
	if idProcesu != "" && len(stany) == 0 {
		return shared.MonitorStatusResponse{}, bladBrakuProcesuTelemetrii(idProcesu)
	}
	return shared.MonitorStatusResponse{Statuses: stany}, nil
}

// ── monitor.subscribe ───────────────────────────────────────────────────────

// ObserwujProcesy obsługuje `monitor.subscribe`: zapisuje okno na telemetrię
// wskazanych procesów i oddaje ich stan bieżący. Pole `windowId` jest oknem
// obserwatora, a nie sitem procesów — sitem są `processIds` i `sessionId`.
func (a *adapterMonitora) ObserwujProcesy(ctx context.Context,
	z shared.MonitorSubscribeRequest) (shared.MonitorSubscribeResponse, error) {

	odpisy, err := a.odpisyProcesow()
	if err != nil {
		return shared.MonitorSubscribeResponse{}, err
	}
	znane := map[string]struct{}{}
	for _, odpis := range odpisy {
		znane[odpis.Id] = struct{}{}
	}
	// Rejestr telemetrii nie wykreśla, identyfikator nieznany znaczy
	// nieistniejący, nie proces domknięty.
	for _, idProcesu := range z.ProcessIds {
		if idProcesu == "" {
			return shared.MonitorSubscribeResponse{}, bladZadaniaMonitora(
				"obserwacja procesu bez wskazania jego identyfikatora")
		}
		if _, jest := znane[idProcesu]; !jest {
			return shared.MonitorSubscribeResponse{}, bladBrakuProcesuTelemetrii(idProcesu)
		}
	}
	sito := sitoProcesow{procesy: z.ProcessIds, idSesji: wartoscTekstu(z.SessionId)}
	stany := a.stanyKontraktu(odpisy, sito, wartoscLiczby(z.Limit))
	zapisana := a.obserwacje.Zapamietaj(wartoscTekstu(z.WindowId), z.ProcessIds)
	return shared.MonitorSubscribeResponse{Statuses: stany, Subscribed: zapisana}, nil
}

// ── odczyt rejestru telemetrii ──────────────────────────────────────────────

// odpisyProcesow oddaje odpis rejestru procesów; rejestr niewpięty odmawia
// głośno. Sesja bywa uzupełniana z okna przed sitem, bo telemetria bywa
// uboższa od stanu rdzenia i sesję okna rozstrzyga rejestr obecności.
func (a *adapterMonitora) odpisyProcesow() ([]procesPostepu, error) {
	if a == nil || a.telemetria == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"monitor procesów: rejestr telemetrii postępu niewpięty — stanu procesów nie ma skąd odczytać"))
	}
	odpisy := a.telemetria.Odpisy()
	if a.obecnosc == nil {
		return odpisy, nil
	}
	for i := range odpisy {
		if odpisy[i].IdSesji == "" && odpisy[i].IdOkna != "" {
			odpisy[i].IdSesji = a.obecnosc.sesjaOkna("", odpisy[i].IdOkna)
		}
	}
	return odpisy, nil
}

// Odpisy oddaje odpis, a nie wskaźniki, rejestru procesów telemetrii
// postępu: proces bywa zmieniany w tej samej chwili przez wątek tury albo
// kolejki, a czytelnik ma dostać stan spójny, nie stan w połowie zmiany.
func (t *telemetriaPostepu) Odpisy() []procesPostepu {
	if t == nil {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	wykaz := make([]procesPostepu, 0, len(t.procesy))
	for _, proces := range t.procesy {
		if proces == nil {
			continue
		}
		wykaz = append(wykaz, *proces)
	}
	return wykaz
}

// sitoProcesow zbiera zawężenia żądania `monitor.status` i `monitor.subscribe`
// nałożone na rejestr procesów. Pole puste nie zawęża niczego, a lista
// procesów pusta znaczy zawężenie do kompletu rejestru.
type sitoProcesow struct {
	idProcesu string
	idOkna    string
	idSesji   string
	procesy   []string
}

// przepuszcza mówi, czy proces spełnia warunki żądania: identyfikator
// procesu, identyfikator okna, identyfikator sesji oraz przynależność do
// listy procesów wskazanej sitem.
func (s sitoProcesow) przepuszcza(p procesPostepu) bool {
	if s.idProcesu != "" && p.Id != s.idProcesu {
		return false
	}
	if s.idOkna != "" && p.IdOkna != s.idOkna {
		return false
	}
	if s.idSesji != "" && p.IdSesji != s.idSesji {
		return false
	}
	if len(s.procesy) == 0 {
		return true
	}
	for _, id := range s.procesy {
		if id == p.Id {
			return true
		}
	}
	return false
}

// stanyKontraktu przekłada odpisy rejestru na byty kontraktu, przesiewa je
// i porządkuje. Granica idzie na koniec, po sicie — obcięcie przed sitem
// oddałoby mniej procesów, niż zażądano, i wyglądałoby na koniec wykazu.
func (a *adapterMonitora) stanyKontraktu(odpisy []procesPostepu, sito sitoProcesow,
	granica int) []shared.MonitorStatus {

	stany := make([]shared.MonitorStatus, 0, len(odpisy))
	for _, odpis := range odpisy {
		if !sito.przepuszcza(odpis) {
			continue
		}
		stany = append(stany, a.stanKontraktu(odpis))
	}
	// Najświeższa zmiana stoi pierwsza. Identyfikator rozstrzyga remisy.
	sort.SliceStable(stany, func(i, j int) bool {
		if stany[i].UpdatedAt != stany[j].UpdatedAt {
			return stany[i].UpdatedAt > stany[j].UpdatedAt
		}
		return stany[i].ProcessId < stany[j].ProcessId
	})
	if granica > 0 && len(stany) > granica {
		stany = stany[:granica]
	}
	return stany
}

// stanKontraktu składa jeden byt `MonitorStatus` z odpisu procesu.
//
// Pole `label` („nazwa procesu do wyświetlenia") zostaje puste: rejestr
// telemetrii zna nazwę etapu, a nie nazwę procesu. Pole jest opcjonalne, więc
// jego brak jest zgodny z kontraktem.
func (a *adapterMonitora) stanKontraktu(p procesPostepu) shared.MonitorStatus {
	stan := shared.MonitorStatus{
		ProcessId: p.Id,
		Status:    p.Stan,
		UpdatedAt: a.chwilaZmiany(p),
	}
	if p.IdOkna != "" {
		okno := p.IdOkna
		stan.WindowId = &okno
	}
	if p.IdSesji != "" {
		sesja := p.IdSesji
		stan.SessionId = &sesja
	}
	if p.Nazwa != "" {
		etap := p.Nazwa
		stan.Stage = &etap
	}
	if p.Etap > 0 {
		numer := p.Etap
		stan.StageIndex = &numer
	}
	// Stopnia ukończenia przy nieznanej liczbie etapów rdzeń nie zmyśla.
	if p.Etapow > 0 {
		etapow := p.Etapow
		stan.StageCount = &etapow
	}
	if p.Etapow > 0 || p.Stan == shared.ProgressStatusDone {
		ukonczenie := int(stopienUkonczenia(p))
		stan.Completion = &ukonczenie
	}
	// Licznik obiegów prowadzi rejestr biegów; okno bez biegu licznika nie ma.
	if p.IdOkna != "" && a.biegi != nil {
		if bieg, jest := a.biegi.Stan(p.IdOkna); jest {
			obiegi := bieg.Loops
			stan.Cycle = &obiegi
		}
	}
	return stan
}

// chwilaZmiany oddaje chwilę ostatniej zmiany procesu w milisekundach epoki,
// czytaną z pamięci czynności rejestru obecności. Proces bez okna chwili
// zmiany nie ma gdzie zapisanej, więc oddawana jest chwila stwierdzenia jego
// stanu przez rdzeń.
func (a *adapterMonitora) chwilaZmiany(p procesPostepu) int64 {
	if p.IdOkna != "" && a.obecnosc != nil {
		if wpis, jest := a.obecnosc.czynnosc.Okno(p.IdOkna); jest {
			return wpis.Chwila.UnixMilli()
		}
	}
	return a.zegar.stwierdzono(p)
}

// Okno zwraca ostatni punkt pracy wskazanego okna. Bliźniak metody `Sesja`,
// osobny dlatego, że monitor pyta o okno procesu, a nie o okno sesji.
func (p *pamiecCzynnosci) Okno(idOkna string) (czynnoscOkna, bool) {
	if p == nil || idOkna == "" {
		return czynnoscOkna{}, false
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	wpis, jest := p.okna[idOkna]
	return wpis, jest
}

// ── zegar odpisów procesów bez okna ─────────────────────────────────────────

// zegarOdpisowProcesow pamięta, kiedy rdzeń po raz pierwszy stwierdził dany
// stan procesu; nie trzyma etapu ani stanu jako faktu, tylko znacznik
// przypięty do odcisku odpisu. Stan niezmieniony znacznika nie przesuwa.
type zegarOdpisowProcesow struct {
	mu        sync.Mutex
	znaczniki map[string]znacznikOdpisu
	teraz     func() time.Time
	pojemnosc int
}

// znacznikOdpisu wiąże odcisk odpisu procesu z chwilą jego pierwszego
// stwierdzenia przez zegar odpisów procesów bez okna.
type znacznikOdpisu struct {
	odcisk string
	chwila time.Time
}

// pojemnoscZegaraProcesow ogranicza pamięć zegara. Rejestr telemetrii procesów
// nie wykreśla, więc bez granicy zegar rósłby razem z nim przez cały czas życia
// rdzenia; po przekroczeniu granicy zegar zaczyna od nowa, zamiast puchnąć.
const pojemnoscZegaraProcesow = 4096

// nowyZegarOdpisowProcesow zakłada pusty zegar odpisów procesów bez okna,
// ograniczony pojemnością `pojemnoscZegaraProcesow`.
func nowyZegarOdpisowProcesow() *zegarOdpisowProcesow {
	return &zegarOdpisowProcesow{
		znaczniki: map[string]znacznikOdpisu{},
		teraz:     func() time.Time { return time.Now().UTC() },
		pojemnosc: pojemnoscZegaraProcesow,
	}
}

// stwierdzono oddaje chwilę, w której rdzeń pierwszy raz zobaczył ten stan
// procesu, w milisekundach epoki.
func (z *zegarOdpisowProcesow) stwierdzono(p procesPostepu) int64 {
	if z == nil {
		return 0
	}
	odcisk := odciskOdpisu(p)
	z.mu.Lock()
	defer z.mu.Unlock()
	if wpis, jest := z.znaczniki[p.Id]; jest && wpis.odcisk == odcisk {
		return wpis.chwila.UnixMilli()
	}
	if len(z.znaczniki) >= z.pojemnosc {
		z.znaczniki = map[string]znacznikOdpisu{}
	}
	chwila := z.teraz()
	z.znaczniki[p.Id] = znacznikOdpisu{odcisk: odcisk, chwila: chwila}
	return chwila.UnixMilli()
}

// odciskOdpisu składa odcisk stanu procesu — wszystko, co monitor o nim oddaje
// poza samym znacznikiem czasu.
func odciskOdpisu(p procesPostepu) string {
	return string(p.Stan) + "|" + p.Nazwa + "|" +
		strconv.Itoa(p.Etap) + "|" + strconv.Itoa(p.Etapow) + "|" +
		p.IdOkna + "|" + p.IdSesji
}

// ── pamięć okien obserwujących procesy ──────────────────────────────────────

// pamiecObserwatorowProcesow wiąże okno z procesami, których telemetrię
// obserwuje. Pusty wykaz procesów znaczy obserwację kompletu. Rejestr
// istnieje po to, by pole `subscribed` niosło prawdę, a rdzeń wiedział,
// które okno których procesów pilnuje.
type pamiecObserwatorowProcesow struct {
	mu      sync.RWMutex
	zakresy map[string][]string
}

// nowaPamiecObserwatorowProcesow zakłada pusty rejestr obserwacji okien nad
// procesami, wypełniany zapisami komendy `monitor.subscribe`.
func nowaPamiecObserwatorowProcesow() *pamiecObserwatorowProcesow {
	return &pamiecObserwatorowProcesow{zakresy: map[string][]string{}}
}

// Zapamietaj zapisuje obserwację okna i mówi, czy została założona. Żądanie bez
// okna obserwacji nie zakłada.
func (p *pamiecObserwatorowProcesow) Zapamietaj(idOkna string, procesy []string) bool {
	if p == nil || idOkna == "" {
		return false
	}
	zakres := make([]string, len(procesy))
	copy(zakres, procesy)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.zakresy[idOkna] = zakres
	return true
}

// Okna zwraca okna obserwujące wskazany proces wraz z oknami obserwującymi
// komplet procesów, posortowane po identyfikatorze okna.
func (p *pamiecObserwatorowProcesow) Okna(idProcesu string) []string {
	if p == nil {
		return nil
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	okna := []string{}
	for idOkna, zakres := range p.zakresy {
		if len(zakres) == 0 {
			okna = append(okna, idOkna)
			continue
		}
		for _, id := range zakres {
			if id == idProcesu {
				okna = append(okna, idOkna)
				break
			}
		}
	}
	sort.Strings(okna)
	return okna
}

// ── odmowy ──────────────────────────────────────────────────────────────────

// bladBrakuProcesuTelemetrii składa odmowę kodem `not_found`, wskazującą
// proces, którego rejestr telemetrii postępu nie zna.
func bladBrakuProcesuTelemetrii(idProcesu string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"monitor procesów: proces "+idProcesu+" nie istnieje w rejestrze telemetrii postępu"))
}

// bladZadaniaMonitora składa odmowę kodem `validation_failed` dla żądania
// niezgodnego z kontraktem komendy monitora procesów.
func bladZadaniaMonitora(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"monitor procesów: "+powod))
}
