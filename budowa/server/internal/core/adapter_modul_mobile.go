// Plik obsługuje rodzinę `mobile.*`: stan platformy, wykaz procesów
// i sterowanie procesem widziane z urządzenia mobilnego, czytane z tego
// samego rejestru telemetrii postępu, co Process Monitor.
package core

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// adapterMobilny wypełnia port mobilnosc: adapter monitora (nawigacja wraz
// z rejestrem telemetrii) rozszerzony o dwa porty czynności.
type adapterMobilny struct {
	*adapterMonitora
	// rozmowa zatrzymuje turę okna, tym samym portem, którym jedzie `message.stop`.
	rozmowa Rozmowa
	// kolejki wykonują działania na kolejce tym samym portem, którym jedzie `queue.action`.
	kolejki Kolejki
}

// ZWarstwaMobilna rozszerza adapter monitora o rodzinę `mobile.*`.
// Wywoływać po `ZTelemetriaProcesow`: warstwa mobilna czyta procesy z rejestru
// telemetrii, a nie z własnego magazynu.
func (a *adapterMonitora) ZWarstwaMobilna(rozmowa Rozmowa, kolejki Kolejki) *adapterMobilny {
	return &adapterMobilny{adapterMonitora: a, rozmowa: rozmowa, kolejki: kolejki}
}

// StanMobilny obsługuje `mobile.status.get`: liczy sesje czynne, okna otwarte,
// procesy w biegu i kolejki czynne, odmawiając zamiast zwracać zero, gdy
// któregoś rejestru nie ma z czego policzyć.
func (a *adapterMobilny) StanMobilny(ctx context.Context,
	z shared.MobileStatusGetRequest) (shared.MobileStatusGetResponse, error) {

	odpisy, err := a.odpisyProcesow()
	if err != nil {
		return shared.MobileStatusGetResponse{}, err
	}
	if a == nil || a.zestaw == nil || a.zestaw.Okna == nil || a.zestaw.Kolejki == nil {
		return shared.MobileStatusGetResponse{}, bladWarstwyMobilnej(
			"repozytoria okien i kolejek nie są wpięte — stanu platformy nie ma z czego policzyć")
	}
	sesje, err := a.sesjeCzynne(ctx)
	if err != nil {
		return shared.MobileStatusGetResponse{}, err
	}
	okna, err := a.zestaw.Okna.LiczbaOtwartych(ctx)
	if err != nil {
		return shared.MobileStatusGetResponse{}, err
	}
	kolejki, err := a.zestaw.Kolejki.LiczbaCzynnych(ctx)
	if err != nil {
		return shared.MobileStatusGetResponse{}, err
	}
	wBiegu := 0
	for _, odpis := range odpisy {
		if odpis.Stan == shared.ProgressStatusRunning {
			wBiegu++
		}
	}
	stan := shared.MobileStatus{
		SessionCount:        len(sesje),
		WindowCount:         okna,
		RunningProcessCount: wBiegu,
		QueueCount:          &kolejki,
		UpdatedAt:           time.Now().UTC().UnixMilli(),
	}
	urzadzenie := strings.TrimSpace(wartoscTekstu(z.DeviceId))
	if urzadzenie != "" {
		wskazane := urzadzenie
		stan.DeviceId = &wskazane
	}
	sparowane, err := a.czySparowane(ctx, urzadzenie)
	if err != nil {
		return shared.MobileStatusGetResponse{}, err
	}
	stan.Paired = sparowane
	return shared.MobileStatusGetResponse{Status: stan}, nil
}

// czySparowane rozstrzyga pole `paired` z katalogu urządzeń: fałsz dla
// urządzenia nieznanego katalogowi albo niewskazanego, odmowa wyłącznie przy
// usterce odczytu katalogu.
func (a *adapterMobilny) czySparowane(ctx context.Context, urzadzenie string) (bool, error) {
	if urzadzenie == "" {
		return false, nil
	}
	if a.zestaw == nil || a.zestaw.Urzadzenia == nil {
		return false, bladWarstwyMobilnej(
			"katalog urządzeń nie jest wpięty — parowania urządzenia " + urzadzenie +
				" nie ma skąd odczytać")
	}
	if numer, err := strconv.ParseInt(urzadzenie, 10, 64); err == nil {
		wiersz, err := a.zestaw.Urzadzenia.Pobierz(ctx, numer)
		switch {
		case err == nil:
			return wiersz.Zaufane, nil
		case errors.Is(err, dane.ErrBrakWiersza):
			return false, nil
		default:
			return false, err
		}
	}
	wiersz, err := a.zestaw.Urzadzenia.PoIdentyfikatorze(ctx, urzadzenie)
	switch {
	case err == nil:
		return wiersz.Zaufane, nil
	case errors.Is(err, dane.ErrBrakWiersza):
		return false, nil
	default:
		return false, err
	}
}

// ProcesyMobilne obsługuje `mobile.process.list`: oddaje procesy platformy
// spełniające warunki żądania, odmawiając przy stanie spoza słownika zamiast
// zwracać pusty wykaz.
func (a *adapterMobilny) ProcesyMobilne(ctx context.Context,
	z shared.MobileProcessListRequest) (shared.MobileProcessListResponse, error) {

	if z.Status != nil {
		if _, znany := shared.WartosciBazyProgressStatus[*z.Status]; !znany {
			return shared.MobileProcessListResponse{}, bladZadaniaMobilnego(
				"stan " + strconv.Quote(string(*z.Status)) + " nie należy do słownika ProgressStatus")
		}
	}
	odpisy, err := a.odpisyProcesow()
	if err != nil {
		return shared.MobileProcessListResponse{}, err
	}
	klucze := a.telemetria.kluczeProcesow()
	sito := sitoProcesow{idSesji: wartoscTekstu(z.SessionId)}
	// Chwila zmiany porządkuje wykaz, nie wychodząc na zewnątrz jako pole kontraktu.
	chwile := map[string]int64{}
	procesy := make([]shared.MobileProcess, 0, len(odpisy))
	for _, odpis := range odpisy {
		if !sito.przepuszcza(odpis) {
			continue
		}
		if z.Status != nil && odpis.Stan != *z.Status {
			continue
		}
		chwile[odpis.Id] = a.chwilaZmiany(odpis)
		procesy = append(procesy, a.procesMobilny(ctx, odpis, klucze[odpis.Id]))
	}
	// Najświeższa zmiana stoi pierwsza; identyfikator rozstrzyga remisy tej samej milisekundy.
	sort.SliceStable(procesy, func(i, j int) bool {
		if chwile[procesy[i].Id] != chwile[procesy[j].Id] {
			return chwile[procesy[i].Id] > chwile[procesy[j].Id]
		}
		return procesy[i].Id < procesy[j].Id
	})
	return shared.MobileProcessListResponse{Processes: procesy}, nil
}

// SterujProcesemMobilnym obsługuje `mobile.process.control`: rozstrzyga drogę
// sterowania po rodzaju procesu i po niej oddaje stan procesu odczytany
// ponownie z rejestru telemetrii.
func (a *adapterMobilny) SterujProcesemMobilnym(ctx context.Context,
	z shared.MobileProcessControlRequest) (shared.MobileProcessControlResponse, error) {

	idProcesu := strings.TrimSpace(z.ProcessId)
	if idProcesu == "" {
		return shared.MobileProcessControlResponse{}, bladZadaniaMobilnego(
			"sterowanie procesem bez wskazania jego identyfikatora")
	}
	if !znaneSterowanieMobilne(z.Control) {
		return shared.MobileProcessControlResponse{}, bladZadaniaMobilnego(
			"sterowanie " + strconv.Quote(string(z.Control)) +
				" nie należy do słownika MobileProcessControl")
	}
	odpisy, err := a.odpisyProcesow()
	if err != nil {
		return shared.MobileProcessControlResponse{}, err
	}
	proces, jest := procesPoIdentyfikatorze(odpisy, idProcesu)
	if !jest {
		return shared.MobileProcessControlResponse{}, bladBrakuProcesuMobilnego(idProcesu)
	}
	// Proces domknięty nie ma czym sterować, więc żądanie takie odmawia.
	if czyStanKoncowy(proces.Stan) {
		return shared.MobileProcessControlResponse{}, bladSporuMobilnego(
			"proces " + idProcesu + " jest już zakończony (stan " + string(proces.Stan) +
				") — nie ma czym sterować")
	}
	klucz := a.telemetria.kluczeProcesow()[idProcesu]
	if numer, jestKolejka := numerKolejkiZKlucza(klucz); jestKolejka {
		err = a.sterujKolejka(ctx, idProcesu, numer, z.Control)
	} else {
		err = a.sterujTura(ctx, idProcesu, proces, z.Control)
	}
	if err != nil {
		return shared.MobileProcessControlResponse{}, err
	}
	poSterowaniu, err := a.odpisyProcesow()
	if err != nil {
		return shared.MobileProcessControlResponse{}, err
	}
	odpis, jest := procesPoIdentyfikatorze(poSterowaniu, idProcesu)
	if !jest {
		return shared.MobileProcessControlResponse{}, bladBrakuProcesuMobilnego(idProcesu)
	}
	return shared.MobileProcessControlResponse{
		Process: a.procesMobilny(ctx, odpis, klucz),
	}, nil
}

// sterujKolejka przekłada sterowanie mobilne na działanie kolejki i wykonuje je
// portem kolejek — tym samym, którym jedzie `queue.action`.
func (a *adapterMobilny) sterujKolejka(ctx context.Context, idProcesu string, numer int64,
	sterowanie shared.MobileProcessControl) error {

	if a.kolejki == nil {
		return bladWarstwyMobilnej("port kolejek nie jest wpięty — procesem kolejki " +
			idProcesu + " nie ma czym sterować")
	}
	var dzialanie shared.QueueAction
	switch sterowanie {
	case shared.MobileProcessControlPause:
		dzialanie = shared.QueueActionPause
	case shared.MobileProcessControlResume, shared.MobileProcessControlApprove:
		dzialanie = shared.QueueActionResume
	case shared.MobileProcessControlStop:
		dzialanie = shared.QueueActionStop
	case shared.MobileProcessControlRestart:
		dzialanie = shared.QueueActionStart
	case shared.MobileProcessControlModify:
		return bladWarstwyMobilnej("proces " + idProcesu + ": korekta zlecenia zmienia jego treść, " +
			"a sterowanie procesem treści poprawionej nie niesie — zlecenie kolejki poprawia się " +
			"poleceniem w oknie rozmowy, nie sterowaniem procesem")
	default:
		// Sterowanie zgadnięte poza znanym słownikiem byłoby czynnością nieżądaną, więc miejsce to odmawia.
		return bladWarstwyMobilnej("proces " + idProcesu + ": sterowania " +
			strconv.Quote(string(sterowanie)) + " nie ma na co przełożyć w słowniku działań kolejki")
	}
	_, err := a.kolejki.Wykonaj(ctx, shared.QueueActionRequest{
		QueueId: strconv.FormatInt(numer, 10),
		Action:  dzialanie,
	})
	return err
}

// sterujTura wykonuje sterowanie procesem tury okna. Rdzeń zna dla tury jedno
// sterowanie — zatrzymanie `message.stop` — więc każde inne żądanie sterowania
// odmawia, nazywając brakujący byt.
func (a *adapterMobilny) sterujTura(ctx context.Context, idProcesu string, proces procesPostepu,
	sterowanie shared.MobileProcessControl) error {

	if sterowanie != shared.MobileProcessControlStop {
		return bladWarstwyMobilnej("proces " + idProcesu + ": tura okna nie ma w rdzeniu sterowania " +
			strconv.Quote(string(sterowanie)) + " — rdzeń zna dla niej wyłącznie zatrzymanie")
	}
	if a.rozmowa == nil {
		return bladWarstwyMobilnej("port rozmowy nie jest wpięty — tury okna " +
			idProcesu + " nie ma czym zatrzymać")
	}
	if proces.IdOkna == "" {
		return bladWarstwyMobilnej("proces " + idProcesu +
			" nie wskazuje okna — zatrzymanie tury idzie przez okno i nie ma bez niego adresu")
	}
	odpowiedz, err := a.rozmowa.Zatrzymaj(ctx, shared.MessageStopRequest{WindowId: proces.IdOkna})
	if err != nil {
		return err
	}
	if !odpowiedz.Stopped {
		return bladSporuMobilnego("proces " + idProcesu + ": w oknie " + proces.IdOkna +
			" nie biegła tura — nie było czego zatrzymać")
	}
	return nil
}

// procesMobilny składa jeden byt `MobileProcess` z odpisu rejestru telemetrii,
// ustalając jego nazwę widoczną i, gdy rejestr zna liczbę etapów, stopień
// ukończenia.
func (a *adapterMobilny) procesMobilny(ctx context.Context, p procesPostepu, klucz string) shared.MobileProcess {
	proces := shared.MobileProcess{
		Id:     p.Id,
		Label:  a.nazwaProcesu(ctx, p, klucz),
		Status: p.Stan,
	}
	if p.IdSesji != "" {
		sesja := p.IdSesji
		proces.SessionId = &sesja
	}
	if p.IdOkna != "" {
		okno := p.IdOkna
		proces.WindowId = &okno
	}
	// Stopnia ukończenia przy nieznanej liczbie etapów rdzeń nie zmyśla, tak jak telemetria i monitor.
	if p.Etapow > 0 || p.Stan == shared.ProgressStatusDone {
		ukonczenie := int(stopienUkonczenia(p))
		proces.Completion = &ukonczenie
	}
	return proces
}

// nazwaProcesu ustala nazwę procesu do wyświetlenia: nazwę kolejki dla
// procesu kolejki albo tytuł okna dla procesu tury, a bez tytułu — opis
// wskazujący byt jednoznacznie.
func (a *adapterMobilny) nazwaProcesu(ctx context.Context, p procesPostepu, klucz string) string {
	if numer, jestKolejka := numerKolejkiZKlucza(klucz); jestKolejka {
		nazwa := ""
		if a.zestaw != nil && a.zestaw.Kolejki != nil {
			if wiersz, err := a.zestaw.Kolejki.PobierzKolejke(ctx, numer); err == nil {
				nazwa = strings.TrimSpace(wiersz.Nazwa)
			}
		}
		if nazwa != "" {
			return nazwa
		}
		return "kolejka " + strconv.FormatInt(numer, 10)
	}
	if p.IdOkna == "" {
		return "proces " + p.Id
	}
	// Tytuł okna czytany jest tą samą drogą i kolejnością, co `window.state.get`.
	if a.nadzorca != nil {
		if okno, err := a.nadzorca.Rejestr().Okno(p.IdOkna); err == nil {
			if tytul := strings.TrimSpace(okno.Tytul); tytul != "" {
				return tytul
			}
		}
	}
	if a.zestaw != nil && a.zestaw.Okna != nil {
		if wiersz, err := a.zestaw.Okna.PoIdentyfikatorze(ctx, p.IdOkna); err == nil && wiersz.Tytul != nil {
			if tytul := strings.TrimSpace(*wiersz.Tytul); tytul != "" {
				return tytul
			}
		}
	}
	return "okno " + p.IdOkna
}

// kluczeProcesow oddaje odwzorowanie identyfikatora procesu na klucz, pod którym
// trzyma go rejestr telemetrii, rozróżniający proces tury okna od procesu
// kolejki.
func (t *telemetriaPostepu) kluczeProcesow() map[string]string {
	if t == nil {
		return map[string]string{}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	klucze := make(map[string]string, len(t.procesy))
	for klucz, proces := range t.procesy {
		if proces == nil {
			continue
		}
		klucze[proces.Id] = klucz
	}
	return klucze
}

// procesPoIdentyfikatorze wyszukuje w wykazie odpisów rejestru telemetrii
// ten jeden, którego identyfikator odpowiada podanemu.
func procesPoIdentyfikatorze(odpisy []procesPostepu, id string) (procesPostepu, bool) {
	for _, odpis := range odpisy {
		if odpis.Id == id {
			return odpis, true
		}
	}
	return procesPostepu{}, false
}

// numerKolejkiZKlucza rozstrzyga, czy klucz rejestru telemetrii wskazuje kolejkę,
// i oddaje jej numer. Klucz nierozpoznany znaczy proces tury okna.
func numerKolejkiZKlucza(klucz string) (int64, bool) {
	if !strings.HasPrefix(klucz, przedrostekProcesuKolejki) {
		return 0, false
	}
	numer, err := strconv.ParseInt(strings.TrimPrefix(klucz, przedrostekProcesuKolejki), 10, 64)
	if err != nil {
		return 0, false
	}
	return numer, true
}

// znaneSterowanieMobilne mówi, czy sterowanie należy do słownika
// `MobileProcessControl` opisanego kontraktem modułu Mobile.
func znaneSterowanieMobilne(sterowanie shared.MobileProcessControl) bool {
	switch sterowanie {
	case shared.MobileProcessControlPause, shared.MobileProcessControlResume,
		shared.MobileProcessControlStop, shared.MobileProcessControlRestart,
		shared.MobileProcessControlApprove, shared.MobileProcessControlModify:
		return true
	default:
		return false
	}
}

// bladBrakuProcesuMobilnego składa odmowę wskazującą proces, którego rejestr
// telemetrii nie zna. Rejestr procesów nie wykreśla, więc identyfikator nieznany
// znaczy proces nieistniejący, a nie proces właśnie domknięty.
func bladBrakuProcesuMobilnego(idProcesu string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"warstwa mobilna: proces "+idProcesu+" nie istnieje w rejestrze telemetrii postępu"))
}

// bladZadaniaMobilnego składa odmowę żądania niezgodnego z kontraktem
// warstwy mobilnej, nazywając powód niezgodności.
func bladZadaniaMobilnego(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"warstwa mobilna: "+powod))
}

// bladSporuMobilnego składa odmowę czynności niewykonalnej w zastanym stanie
// procesu albo tury, nazywając ten stan wprost.
func bladSporuMobilnego(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"warstwa mobilna: "+powod))
}

// bladWarstwyMobilnej składa odmowę bytu, którego rdzeniowi brakuje do
// wykonania żądanej czynności warstwy mobilnej.
func bladWarstwyMobilnej(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"warstwa mobilna: "+powod))
}
