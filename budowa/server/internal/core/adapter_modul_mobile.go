// Odpowiedzialność pliku: rodzina `mobile.*` — mobilne centrum dowodzenia.
// Trzy komendy: stan platformy (`mobile.status.get`), wykaz procesów
// (`mobile.process.list`) i sterowanie procesem (`mobile.process.control`).
//
// Proces tej rodziny nie jest tym samym bytem, co `terminal_proces`.
// Rozstrzyga kontrakt: pole `MobileProcess.status` ma typ `ProgressStatus` —
// słownik telemetrii postępu — podczas gdy kolumna `terminal_proces.stan`
// niesie CHECK(stan IN ('running','finished','failed','stopped')), czyli
// słownik `TerminalProcessStatus`. Dwa różne słowniki to dwa różne byty.
// Proces mobilny jest więc tym samym procesem, który pokazuje Process Monitor
// warstwy wspólnej: wpisem rejestru telemetrii postępu.
//
// Tabeli pod tym nie ma i mieć nie będzie. Rejestr procesów żyje w pamięci
// i droga odczytu działa bez tabeli, więc warstwa mobilna nie zakłada ani
// drugiej tabeli procesów, ani drugiego rejestru — czyta ten sam
// `telemetriaPostepu`, który czyta monitor i który rozgłasza
// `progress.changed`. Plik jest fasadą, nie magazynem.
//
// Mobilne centrum dowodzenia pyta o stan całej platformy — sesje czynne, okna,
// procesy, kolejki — czyli o warstwę wspólną, którą obsługuje port Nawigacja.
// Adapter mobilny osadza adapter monitora (nawigacja wraz z telemetrią),
// zamiast zakładać port równoległy nad tym samym rejestrem, i dokłada dwa porty
// czynności, których warstwa wspólna sama nie ma: rozmowę (zatrzymanie tury
// okna) i kolejki (działania na kolejce).
//
// Sterowanie nie ma własnego silnika. `mobile.process.control` niczego nie
// wykonuje sam: proces kolejki idzie tą samą drogą, co `queue.action`, a tura
// okna tą samą, co `message.stop` — razem z ich telemetrią i rozgłoszeniami.
// Gdyby warstwa mobilna zatrzymywała turę własnym wywołaniem, stan procesu
// zmieniłby się bez zdarzenia postępu i pulpit pokazywałby proces w biegu,
// którego już nie ma.
//
// Czego rdzeń nie umie, tego ta rodzina nie udaje: tura okna nie ma w rdzeniu
// ani wstrzymania, ani wznowienia, ani ponownego uruchomienia — sterowanie takie
// odmawia głośno, nazywając brakujący byt, zamiast oddać proces „po sterowaniu",
// którym nikt nie sterował. Tak samo `modify`: korekta zlecenia zmienia jego
// treść, a żądanie sterowania niesie wyłącznie proces, sterowanie i urządzenie —
// więc korekta odmawia i wskazuje drogę, którą treść zlecenia naprawdę się
// poprawia (polecenie w oknie rozmowy).
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
	// rozmowa zatrzymuje turę okna. Ten sam port, którym jedzie `message.stop`,
	// więc zatrzymanie idzie razem z telemetrią zatrzymania (telemetria_tury.go).
	rozmowa Rozmowa
	// kolejki wykonują działania na kolejce. Ten sam port, którym jedzie
	// `queue.action`, więc stan kolejki i jej telemetria zmieniają się raz.
	kolejki Kolejki
}

// ZWarstwaMobilna rozszerza adapter monitora o rodzinę `mobile.*`.
// Wywoływać po `ZTelemetriaProcesow`: warstwa mobilna czyta procesy z rejestru
// telemetrii, a nie z własnego magazynu.
func (a *adapterMonitora) ZWarstwaMobilna(rozmowa Rozmowa, kolejki Kolejki) *adapterMobilny {
	return &adapterMobilny{adapterMonitora: a, rozmowa: rozmowa, kolejki: kolejki}
}

// ── mobile.status.get ───────────────────────────────────────────────────────

// StanMobilny obsługuje `mobile.status.get`: oddaje stan platformy widziany
// z urządzenia mobilnego.
//
// Każda liczba ma jedno źródło i jest nim źródło już istniejące:
//   - sesje czynne — ta sama metoda, którą liczy je Strona główna (`sesjeCzynne`),
//     czyli rejestr żywy uzupełniony wierszami sesji nieodtworzonych;
//   - okna komunikacji — wiersze o stanie `otwarte` (repozytorium okien);
//   - procesy w biegu — wpisy rejestru telemetrii o stanie `running`;
//   - kolejki czynne — wiersze kolejek poza stanem końcowym (repozytorium kolejek).
//
// Rejestr telemetrii niewpięty albo repozytorium niewpięte nie dają „zera
// procesów" ani „zera kolejek" — dają odmowę z nazwą brakującego bytu. Zero
// jest odpowiedzią prawdziwą wyłącznie wtedy, gdy naprawdę policzono.
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

// czySparowane rozstrzyga pole `paired` z katalogu urządzeń (tabela
// `urzadzenie`).
//
// Żądanie bez urządzenia nie jest urządzeniem niesparowanym — jest brakiem
// wskazania, a wtedy nie ma czego sprawdzać i pole niesie fałsz zgodnie ze swoim
// brzmieniem („czy urządzenie jest sparowane").
//
// Urządzenie nieznane katalogowi to fałsz, nie odmowa: `paired: false` jest
// odpowiedzią pełną i prawdziwą, bo maszyna spoza katalogu sparowana nie jest.
// Odmowa idzie wyłącznie wtedy, gdy katalogu nie ma jak zapytać (usterka
// odczytu) — wtedy fałsz byłby zgadywaniem.
//
// Katalog zna dwa klucze i oba są jego własnymi kluczami, nie drugą prawdą:
// numer wiersza (tak wskazuje urządzenie punkt dostępu, `handlers_dostep_przeklad.go`)
// oraz identyfikator sprzętowy (tak rozpoznaje maszynę rdzeń, `urzadzenia.go`).
// Klient mobilny może podać jeden albo drugi, więc sprawdzamy oba — najpierw
// numer, bo tylko on ma kształt rozstrzygalny bez zapytania.
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

// ── mobile.process.list ─────────────────────────────────────────────────────

// ProcesyMobilne obsługuje `mobile.process.list`: oddaje procesy platformy
// spełniające warunki żądania.
//
// Urządzenie nie zawęża wykazu i nie jest to przeoczenie. Proces należy do
// platformy — do okna albo do kolejki — a nie do telefonu, który o niego pyta;
// rejestr telemetrii nie zna i nie ma jak znać urządzenia pytającego. Pole
// `deviceId` żądania mówi, kto pyta, a nie czego dotyczy pytanie.
//
// Stan spoza słownika jest odmową, nie pustym wykazem. Sito po stanie nieznanym
// przepuściłoby zero procesów i wyglądałoby to na „nic nie pracuje", zamiast na
// „takiego stanu nie ma w kontrakcie".
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
	// Chwila zmiany nie wychodzi na zewnątrz — kontrakt `MobileProcess` nie ma
	// na nią pola — ale porządkuje wykaz, żeby mobilne centrum dowodzenia
	// czytało procesy w tej samej kolejności, co Process Monitor.
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
	// Najświeższa zmiana stoi pierwsza. Identyfikator rozstrzyga remisy, żeby
	// wykaz nie przestawiał się przy dwóch procesach zgłoszonych w tej samej
	// milisekundzie.
	sort.SliceStable(procesy, func(i, j int) bool {
		if chwile[procesy[i].Id] != chwile[procesy[j].Id] {
			return chwile[procesy[i].Id] > chwile[procesy[j].Id]
		}
		return procesy[i].Id < procesy[j].Id
	})
	return shared.MobileProcessListResponse{Processes: procesy}, nil
}

// ── mobile.process.control ──────────────────────────────────────────────────

// SterujProcesemMobilnym obsługuje `mobile.process.control`.
//
// Droga sterowania wynika z rodzaju procesu, a rodzaj mówi klucz rejestru
// telemetrii, nie zgadywanie po polach. Klucz z przedrostkiem `kolejka-` niesie
// numer kolejki (telemetria_kolejki.go), każdy inny klucz jest identyfikatorem
// okna (telemetria_tury.go, `opisTury`). Dlatego proces kolejki idzie do portu
// kolejek, a proces tury do portu rozmowy.
//
// Po sterowaniu rejestr czytany jest ponownie: odpowiedź niesie „proces po
// sterowaniu", więc pokazuje stan zastany po czynności, a nie odpis sprzed niej
// z doklejonym stanem, którego nikt nie potwierdził.
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
	// Proces domknięty nie ma czym sterować. Oddanie go „po sterowaniu" bez
	// zmiany byłoby zgłoszeniem czynności, której nie było.
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
//
// Trzy przełożenia są jednoznaczne (pause→pause, resume→resume, stop→stop), bo
// oba słowniki nazywają tę samą czynność. `restart` („Uruchom ponownie") idzie
// na `start`, a nie na `retry`, ponieważ `retry` dotyczy pozycji kolejki
// i bierze `itemId`, a sterowanie mobilne dotyczy procesu jako całości i pozycji
// nie wskazuje.
//
// `approve` („Zatwierdź") idzie na `resume` i nie jest to przełożenie na skróty.
// Przyjęcie wyniku kroku wyraża w rdzeniu wybór działania, a nie osobne
// działanie: krok naprzód (`start`/`resume`) znaczy przyjęcie, `retry` znaczy
// odesłanie do poprawy — tak nazywa to protokół silnika kolejki
// (`kolejka_silnik.go`), i tam też stoi wprost, że działania „zamknij pozycję
// werdyktem" kontrakt nie ma. Pozycja stojąca w `do_weryfikacji` przechodzi tym
// krokiem w `ukonczona` z werdyktem przyjęcia, czyli zatwierdzenie ma skutek
// zapisany, a nie samą etykietę odpowiedzi.
//
// `modify` („Modyfikuj") nie ma tu czego wykonać i mówi to wprost. Korekta
// zlecenia w toku jest zmianą jego treści, a `mobile.process.control` treści
// poprawionej nie niesie — żądanie ma wyłącznie proces, sterowanie
// i urządzenie. Rdzeń nie ma więc czego zapisać, a droga korekty zlecenia
// biegnie poleceniem w oknie rozmowy (`message.send`), nie sterowaniem procesem.
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
		// Słownik urósł, a to miejsce nie — działanie zgadnięte byłoby czynnością
		// wykonaną nie na żądanie Operatora, więc zamiast zgadywać nazywamy brak.
		return bladWarstwyMobilnej("proces " + idProcesu + ": sterowania " +
			strconv.Quote(string(sterowanie)) + " nie ma na co przełożyć w słowniku działań kolejki")
	}
	_, err := a.kolejki.Wykonaj(ctx, shared.QueueActionRequest{
		QueueId: strconv.FormatInt(numer, 10),
		Action:  dzialanie,
	})
	return err
}

// sterujTura wykonuje sterowanie procesem tury okna.
//
// Rdzeń zna dla tury jedno sterowanie — zatrzymanie (`message.stop`).
// Wstrzymania, wznowienia ani ponownego uruchomienia tury w rdzeniu nie ma:
// pętla tury to strumień odpowiedzi modelu prowadzony pod kontekstem
// z odwołaniem, a odwołania nie da się cofnąć. Odmowa nazywa brakujący byt
// zamiast oddawać proces „po sterowaniu", którego nikt nie tknął.
//
// Zatwierdzenia tura też nie ma i nie jest to przeoczenie tego pliku: punktu
// decyzyjnego kontrakt nie zna — nie ma ani komendy zatwierdzenia kroku, ani
// zdarzenia proszącego o decyzję. Zatwierdzać można wyłącznie tam, gdzie stan
// czekający naprawdę stoi zapisany, czyli w pozycji kolejki. Modyfikacji
// dotyczy to samo, co w kolejce: treści poprawionej żądanie nie niesie.
//
// Tura, która nie biegła, jest sporem stanu, nie powodzeniem. `message.stop`
// odpowiada wtedy `stopped: false` i tak mówi prawdę; odpowiedź mobilna nie ma
// takiego pola, więc prawdę może powiedzieć wyłącznie odmową.
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

// ── przekład na byt kontraktu ───────────────────────────────────────────────

// procesMobilny składa jeden byt `MobileProcess` z odpisu rejestru telemetrii.
//
// Pole `label` jest wymagane, więc musi nieść coś prawdziwego. Rejestr
// telemetrii zna nazwę etapu, a nie nazwę procesu — monitor zostawia z tego
// powodu swoje pole `label` puste, bo u niego jest opcjonalne. Tutaj pustki
// zostawić nie wolno, a etap w polu nazwy byłby podaniem jednej rzeczy za
// drugą. Nazwą procesu jest więc nazwa bytu, którego proces dotyczy: tytuł okna
// albo nazwa kolejki. Byt bez nazwy własnej wychodzi opisem wskazującym go
// jednoznacznie („okno win-3"), a nie napisem zmyślonym.
//
// Pola `startedAt` nie ma i nie jest to przeoczenie. Rejestr telemetrii nie
// znakuje procesów chwilą uruchomienia — zna wyłącznie stan i etap, a chwilę
// ostatniej zmiany rdzeń odtwarza z pamięci czynności rejestru obecności
// (monitor, `chwilaZmiany`). Chwila ostatniej zmiany to nie chwila startu,
// a pole jest opcjonalne, więc brak jest tu odpowiedzią zgodną z kontraktem.
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
	// Stopnia ukończenia przy nieznanej liczbie etapów rdzeń nie zmyśla — tak
	// samo liczy go telemetria dla `progress.changed` i monitor dla `monitor.status`.
	if p.Etapow > 0 || p.Stan == shared.ProgressStatusDone {
		ukonczenie := int(stopienUkonczenia(p))
		proces.Completion = &ukonczenie
	}
	return proces
}

// nazwaProcesu ustala nazwę procesu do wyświetlenia. Patrz `procesMobilny`.
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
	// Tytuł okna bierzemy tą samą drogą i w tej samej kolejności, co
	// `window.state.get`: najpierw rejestr żywy, potem wiersz. Nieudany odczyt
	// nie jest usterką wykazu — proces zostaje w wykazie pod opisem
	// wskazującym jego okno, a nie znika z niego.
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

// ── odczyt rejestru telemetrii ──────────────────────────────────────────────

// kluczeProcesow oddaje odwzorowanie identyfikatora procesu na klucz, pod którym
// trzyma go rejestr telemetrii. Klucz jest jedynym miejscem, w którym rdzeń
// zapisał rodzaj procesu: okno prowadzi turę pod kluczem okna, kolejka pod
// kluczem `kolejka-<numer>`. Bez tego odwzorowania sterowanie musiałoby rodzaj
// zgadywać z pól odpisu, a proces kolejki bez okien wyglądałby wtedy tak samo,
// jak proces tury, której okno przepadło.
//
// Metoda mieszka w pliku warstwy mobilnej, bo to ona jej potrzebuje — producent
// telemetrii nie ma innych czytelników, tak samo jak `Odpisy` mieszka
// w pliku monitora.
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

// procesPoIdentyfikatorze wyszukuje odpis procesu w wykazie.
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

// znaneSterowanieMobilne mówi, czy sterowanie należy do słownika kontraktu.
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

// ── odmowy ──────────────────────────────────────────────────────────────────

// bladBrakuProcesuMobilnego składa odmowę wskazującą proces, którego rejestr
// telemetrii nie zna. Rejestr procesów nie wykreśla, więc identyfikator nieznany
// znaczy proces nieistniejący, a nie proces właśnie domknięty.
func bladBrakuProcesuMobilnego(idProcesu string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"warstwa mobilna: proces "+idProcesu+" nie istnieje w rejestrze telemetrii postępu"))
}

// bladZadaniaMobilnego składa odmowę żądania niezgodnego z kontraktem.
func bladZadaniaMobilnego(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"warstwa mobilna: "+powod))
}

// bladSporuMobilnego składa odmowę czynności niewykonalnej w zastanym stanie.
func bladSporuMobilnego(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"warstwa mobilna: "+powod))
}

// bladWarstwyMobilnej składa odmowę bytu, którego rdzeniowi brakuje.
func bladWarstwyMobilnej(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"warstwa mobilna: "+powod))
}
