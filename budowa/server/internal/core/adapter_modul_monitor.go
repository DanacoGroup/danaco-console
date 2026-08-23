// Odpowiedzialność pliku: rodzina `monitor.*` — Process Monitor warstwy
// wspólnej. Dwie komendy: odczyt stanu bieżącego procesów
// (`monitor.status`) i zapis okna na ich telemetrię (`monitor.subscribe`).
//
// Rejestr procesów rdzenia jest w pamięci; tabeli procesów nie ma. Proces okna
// to proces systemowy objęty uchwytem rdzenia i ginie razem z rdzeniem, więc
// wiersz, który przeżyłby restart ze stanem `running`, opisywałby proces
// nieistniejący. Monitor czyta ten sam rejestr, który rozgłasza
// `progress.changed` (`telemetria.go`) — jedno źródło faktu, nie dwa.
//
// Gdy port nie niesie telemetrii, monitor nie ma czego czytać i odmawia kodem
// `internal_error`, nazywając brakujący byt; pusty wykaz procesów byłby wtedy
// odpowiedzią nieprawdziwą. Rejestr wpięty i pusty to co innego — wtedy pusty
// wykaz jest prawdą.
//
// Monitor procesów jest oknem warstwy wspólnej, nie modułu (zdarzenia.go).
// Warstwę wspólną — stronę główną, środowiska, moduły i stan okna operacyjnego
// — obsługuje port Nawigacja, więc monitor osadza jego adapter zamiast zakładać
// port równoległy.
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
	// telemetria jest jedynym rejestrem procesów rdzenia. Zerowa znaczy odmowę
	// z kodem `internal_error`, nie pusty wykaz.
	telemetria *telemetriaPostepu
	obserwacje *pamiecObserwatorowProcesow
	// zegar zapamiętuje chwilę, w której rdzeń STWIERDZIŁ stan procesu, którego
	// chwili zmiany nie zapisał nikt inny. Patrz `chwilaZmiany` niżej.
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
// spełniających warunki żądania, bez zakładania obserwacji.
//
// Proces wskazany wprost, którego rejestr nie zna, jest bytem nieistniejącym:
// zamiast pustego wykazu idzie odmowa `not_found` z nazwą procesu.
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
// wskazanych procesów i oddaje ich stan bieżący.
//
// Dwa pola żądania mówią o dwóch różnych oknach: `windowId` jest oknem
// odbierającym telemetrię (obserwatorem), a nie sitem procesów — sitem są
// `processIds` (pusta lista znaczy komplet) oraz `sessionId`. Tak samo
// rozstrzyga to bliźniacza komenda `automation.execution.subscribe`.
//
// Żądanie bez `windowId` jest zwykłym odczytem — nie ma czego zapisać, więc
// pole `subscribed` niesie `false`. Zdarzenie `progress.changed` i tak dociera
// do wszystkich połączeń konta; zapis mówi rdzeniowi, które okno których
// procesów pilnuje.
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
	// Obserwacja procesu, którego rejestr nie zna, jest obserwacją niczego.
	// Rejestr telemetrii procesów nie wykreśla, więc identyfikator nieznany
	// znaczy identyfikator nieistniejący, a nie proces właśnie domknięty.
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

// odpisyProcesow oddaje odpis rejestru procesów. Rejestr niewpięty odmawia
// głośno — patrz nagłówek pliku.
//
// Sesja procesu bywa uzupełniana z okna i dzieje się to przed sitem, nie po
// nim. Telemetria bywa uboższa od stanu rdzenia: proces otwarty na
// niepowodzeniu przyjęcia wiadomości zna okno, lecz nie zna jeszcze sesji
// (`opisTury(z.WindowId, "")`, telemetria_tury.go). Sito po samym zapisie
// telemetrii oddałoby pustkę przy procesie, który do wskazanej sesji należy.
// Sesję okna rozstrzyga rejestr obecności — tą samą metodą, którą rozstrzyga
// ją nasłuch telemetrii.
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

// Odpisy oddaje odpis rejestru procesów telemetrii postępu.
//
// Odpis, a nie wskaźniki: proces bywa zmieniany w tej samej chwili przez wątek
// tury albo kolejki, a czytelnik ma dostać stan spójny, nie stan w połowie
// zmiany. Metoda mieszka w pliku monitora, bo to monitor jej potrzebuje —
// producent telemetrii nie ma czytelników poza nim.
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

// sitoProcesow zbiera zawężenia żądania. Pole puste nie zawęża niczego.
type sitoProcesow struct {
	idProcesu string
	idOkna    string
	idSesji   string
	procesy   []string
}

// przepuszcza mówi, czy proces spełnia warunki żądania.
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
	// Najświeższa zmiana stoi pierwsza — Process Monitor czyta od góry.
	// Identyfikator rozstrzyga remisy, żeby wykaz nie przestawiał się przy
	// dwóch procesach zgłoszonych w tej samej milisekundzie.
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
	// Stopnia ukończenia przy nieznanej liczbie etapów rdzeń nie zmyśla
	// — tak samo liczy go telemetria dla `progress.changed`.
	if p.Etapow > 0 {
		etapow := p.Etapow
		stan.StageCount = &etapow
	}
	if p.Etapow > 0 || p.Stan == shared.ProgressStatusDone {
		ukonczenie := int(stopienUkonczenia(p))
		stan.Completion = &ukonczenie
	}
	// Licznik obiegów naprawczych bez granicy prowadzi rejestr biegów
	// koordynatorów; proces okna, które biegu nie prowadzi, licznika nie ma.
	if p.IdOkna != "" && a.biegi != nil {
		if bieg, jest := a.biegi.Stan(p.IdOkna); jest {
			obiegi := bieg.Loops
			stan.Cycle = &obiegi
		}
	}
	return stan
}

// chwilaZmiany oddaje chwilę ostatniej zmiany procesu w milisekundach epoki.
//
// Źródłem pierwszym jest pamięć czynności rejestru obecności: zapisuje ona
// znacznik przy każdym zgłoszeniu telemetrii okna (stan_sesji_czynnosc.go),
// więc dla procesu okna jest to dokładnie chwila ostatniej zmiany.
//
// Proces bez okna chwili zmiany nie ma gdzie zapisanej: telemetria trzyma etap
// i stan, lecz nie czas, a pamięć czynności jest kluczowana oknem i proces bez
// okna (kolejka założona przed pierwszą wiadomością) do niej nie trafia. Dla
// takiego procesu monitor podaje chwilę, w której rdzeń stwierdził ten stan —
// wartość będącą górnym ograniczeniem chwili zmiany, nie nią samą.
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
// stan procesu. Nie jest to drugi rejestr procesów: nie trzyma ani etapu, ani
// stanu jako faktu, tylko znacznik przypięty do odcisku odpisu. Stan
// niezmieniony znacznika nie przesuwa, więc kolejne odczyty monitora nie
// odmładzają procesu, który stoi.
type zegarOdpisowProcesow struct {
	mu        sync.Mutex
	znaczniki map[string]znacznikOdpisu
	teraz     func() time.Time
	pojemnosc int
}

// znacznikOdpisu wiąże odcisk odpisu z chwilą jego stwierdzenia.
type znacznikOdpisu struct {
	odcisk string
	chwila time.Time
}

// pojemnoscZegaraProcesow ogranicza pamięć zegara. Rejestr telemetrii procesów
// nie wykreśla, więc bez granicy zegar rósłby razem z nim przez cały czas życia
// rdzenia; po przekroczeniu granicy zegar zaczyna od nowa, zamiast puchnąć.
const pojemnoscZegaraProcesow = 4096

// nowyZegarOdpisowProcesow zakłada pusty zegar.
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
// obserwuje. Pusty wykaz procesów znaczy obserwację kompletu — tak mówi
// kontrakt komendy wprost.
//
// Zdarzenie `progress.changed` dociera do wszystkich połączeń konta i rejestr
// niczego w tym nie zmienia. Rejestr istnieje po to, by pole `subscribed`
// niosło prawdę, a rdzeń wiedział, które okno których procesów pilnuje. Ten sam
// wzorzec prowadzi `pamiecObserwatorowPrzebiegow` dla Execution Monitora.
type pamiecObserwatorowProcesow struct {
	mu      sync.RWMutex
	zakresy map[string][]string
}

// nowaPamiecObserwatorowProcesow zakłada pusty rejestr obserwacji.
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
// komplet procesów.
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

// bladBrakuProcesuTelemetrii składa odmowę wskazującą proces, którego rejestr
// telemetrii nie zna.
func bladBrakuProcesuTelemetrii(idProcesu string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"monitor procesów: proces "+idProcesu+" nie istnieje w rejestrze telemetrii postępu"))
}

// bladZadaniaMonitora składa odmowę żądania niezgodnego z kontraktem.
func bladZadaniaMonitora(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"monitor procesów: "+powod))
}
