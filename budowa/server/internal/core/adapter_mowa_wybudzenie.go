// Plik obsługuje `speech.wake.get`, `speech.wake.set`, `speech.listen.start` i `speech.listen.stop`: nastawę wybudzania oraz nasłuch ciągły rdzenia jako umowę z oknem, a nie dostęp do jego mikrofonu.
package core

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/mowa"
	"danacoconsole/shared"
)

const (
	// Klucze są nastawami konfiguracji, nie polami kontraktu: kontrakt niesie `WakeWordConfig`
	// w całości, a konfiguracja rozstrzyga jej pola osobno po poziomach.
	kluczFrazyWybudzajacej = "mowa_fraza_wybudzajaca"
	kluczTrybuNasluchu     = "mowa_tryb_nasluchu"
	kluczProguDetekcji     = "mowa_prog_detekcji"
	kluczOdszumiania       = "mowa_odszumianie"

	przedrostekNasluchu = "nasluch-"
)

// Rejestr jest w pamięci, nie w bazie: nasłuch żyje tyle, co połączenie karty,
// a wiersz w bazie przeżyłby restart rdzenia i obiecywał nasłuch, którego nie ma.
type rejestrNasluchow struct {
	zamek    sync.Mutex
	wgOkna   map[string]nasluchMowy
	wgKlucza map[string]string
}

type nasluchMowy struct {
	Kod    string
	Okno   string
	Sesja  string
	Tryb   shared.ListenMode
	Jezyk  string
	Ruszyl int64
}

func nowyRejestrNasluchow() *rejestrNasluchow {
	return &rejestrNasluchow{
		wgOkna:   map[string]nasluchMowy{},
		wgKlucza: map[string]string{},
	}
}

func (r *rejestrNasluchow) zaloz(n nasluchMowy) {
	r.zamek.Lock()
	defer r.zamek.Unlock()
	if poprzedni, jest := r.wgOkna[n.Okno]; jest {
		delete(r.wgKlucza, poprzedni.Kod)
	}
	r.wgOkna[n.Okno] = n
	r.wgKlucza[n.Kod] = n.Okno
}

// Zatrzymanie nasłuchu, którego nie ma, nie jest błędem — tak stanowi kontrakt.
func (r *rejestrNasluchow) zdejmij(kod, okno string) bool {
	r.zamek.Lock()
	defer r.zamek.Unlock()
	if kod != "" {
		oknoNasluchu, jest := r.wgKlucza[kod]
		if !jest {
			return false
		}
		delete(r.wgKlucza, kod)
		delete(r.wgOkna, oknoNasluchu)
		return true
	}
	if okno == "" {
		return false
	}
	nasluch, jest := r.wgOkna[okno]
	if !jest {
		return false
	}
	delete(r.wgKlucza, nasluch.Kod)
	delete(r.wgOkna, okno)
	return true
}

func (r *rejestrNasluchow) nasluchOkna(okno string) (nasluchMowy, bool) {
	r.zamek.Lock()
	defer r.zamek.Unlock()
	nasluch, jest := r.wgOkna[okno]
	return nasluch, jest
}

func (a *adapterMowy) NastawaWybudzania(ctx context.Context,
	z shared.SpeechWakeGetRequest) (shared.SpeechWakeGetResponse, error) {

	nastawa := a.odczytajNastaweWybudzania(ctx, z.Scope, z.ScopeId)
	gotowosc, err := a.Gotowosc(ctx, shared.SpeechAvailabilityGetRequest{})
	if err != nil {
		// Niepowodzenie sprawdzenia silnika nie dowodzi niewykonalności — oddaje powód, nie rozstrzygnięcie.
		powod := "nie udało się sprawdzić silnika mowy: " + err.Error()
		return shared.SpeechWakeGetResponse{Config: nastawa, Available: false, Reason: &powod}, nil
	}

	if !gotowosc.Available {
		powod := "wybudzanie frazą wymaga silnika rozpoznawania mowy, a ten nie jest gotowy"
		if gotowosc.Reason != nil && strings.TrimSpace(*gotowosc.Reason) != "" {
			powod += " — " + *gotowosc.Reason
		}
		return shared.SpeechWakeGetResponse{Config: nastawa, Available: false, Reason: &powod}, nil
	}
	if strings.TrimSpace(nastawa.Phrase) == "" && nastawa.Mode == shared.ListenModeWakeWord {
		powod := "tryb nasłuchu jest ustawiony na frazę wybudzającą, ale sama fraza jest pusta; " +
			"naprawa: wpisać frazę komendą speech.wake.set albo przestawić tryb"
		return shared.SpeechWakeGetResponse{Config: nastawa, Available: false, Reason: &powod}, nil
	}
	return shared.SpeechWakeGetResponse{Config: nastawa, Available: true}, nil
}

func (a *adapterMowy) ZapiszNastaweWybudzania(ctx context.Context,
	z shared.SpeechWakeSetRequest) (shared.SpeechWakeSetResponse, error) {

	if a.konfiguracja == nil {
		return shared.SpeechWakeSetResponse{}, bladZapleczaNagran(
			"serwer nie ma wpiętego magazynu konfiguracji — nastawy wybudzania nie ma gdzie zapisać")
	}
	poziom := shared.ConfigScopeGlobal
	if z.Scope != nil && strings.TrimSpace(string(*z.Scope)) != "" {
		poziom = string(*z.Scope)
	}
	if err := sprawdzPoziomKontekstu(shared.ConfigScope(poziom)); err != nil {
		return shared.SpeechWakeSetResponse{}, err
	}
	kluczZasiegu := strings.TrimSpace(wartoscTekstu(z.ScopeId))

	if z.Mode != nil {
		if err := sprawdzTrybNasluchu(*z.Mode); err != nil {
			return shared.SpeechWakeSetResponse{}, err
		}
		if err := a.zapiszNastaweMowy(ctx, poziom, kluczZasiegu, kluczTrybuNasluchu,
			string(*z.Mode)); err != nil {
			return shared.SpeechWakeSetResponse{}, err
		}
	}
	if z.Phrase != nil {
		if err := a.zapiszNastaweMowy(ctx, poziom, kluczZasiegu, kluczFrazyWybudzajacej,
			strings.TrimSpace(*z.Phrase)); err != nil {
			return shared.SpeechWakeSetResponse{}, err
		}
	}
	if z.VadThreshold != nil {
		if *z.VadThreshold < 0 || *z.VadThreshold > 100 {
			return shared.SpeechWakeSetResponse{}, bladWskazaniaNagrania(
				"próg detekcji mowy podaje się w setnych z zakresu 0–100; zero znaczy próg domyślny silnika")
		}
		if err := a.zapiszNastaweMowy(ctx, poziom, kluczZasiegu, kluczProguDetekcji,
			strconv.Itoa(*z.VadThreshold)); err != nil {
			return shared.SpeechWakeSetResponse{}, err
		}
	}
	if z.NoiseSuppression != nil {
		wartosc := "0"
		if *z.NoiseSuppression {
			wartosc = "1"
		}
		if err := a.zapiszNastaweMowy(ctx, poziom, kluczZasiegu, kluczOdszumiania, wartosc); err != nil {
			return shared.SpeechWakeSetResponse{}, err
		}
	}

	return shared.SpeechWakeSetResponse{
		Config: a.odczytajNastaweWybudzania(ctx, z.Scope, z.ScopeId),
	}, nil
}

func (a *adapterMowy) UruchomNasluch(ctx context.Context,
	z shared.SpeechListenStartRequest) (shared.SpeechListenStartResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.SpeechListenStartResponse{}, bladWskazaniaNagrania(
			"nasłuch bez wskazania okna — serwer nie wiedziałby, czyje odcinki rozpoznaje")
	}
	if a.nasluchy == nil {
		return shared.SpeechListenStartResponse{}, bladZapleczaNagran(
			"serwer nie ma rejestru nasłuchów")
	}

	nastawa := a.odczytajNastaweWybudzania(ctx, nil, nil)
	tryb := nastawa.Mode
	if z.Mode != nil && strings.TrimSpace(string(*z.Mode)) != "" {
		if err := sprawdzTrybNasluchu(*z.Mode); err != nil {
			return shared.SpeechListenStartResponse{}, err
		}
		tryb = *z.Mode
	}

	gotowosc, err := a.Gotowosc(ctx, shared.SpeechAvailabilityGetRequest{})
	if err != nil || !gotowosc.Available {
		powod := "nasłuch ciągły rozpoznaje każdy odcinek silnikiem mowy, a ten nie jest gotowy"
		if err != nil {
			powod += " — " + err.Error()
		} else if gotowosc.Reason != nil && strings.TrimSpace(*gotowosc.Reason) != "" {
			powod += " — " + *gotowosc.Reason
		}
		return shared.SpeechListenStartResponse{
			Listening: false, ListenerId: "", Reason: &powod,
		}, nil
	}

	nasluch := nasluchMowy{
		Kod:    nowyIdentyfikator(przedrostekNasluchu),
		Okno:   strings.TrimSpace(z.WindowId),
		Sesja:  strings.TrimSpace(wartoscTekstu(z.SessionId)),
		Tryb:   tryb,
		Jezyk:  strings.TrimSpace(wartoscTekstu(z.Language)),
		Ruszyl: time.Now().UnixMilli(),
	}
	a.nasluchy.zaloz(nasluch)
	return shared.SpeechListenStartResponse{Listening: true, ListenerId: nasluch.Kod}, nil
}

func (a *adapterMowy) ZatrzymajNasluch(_ context.Context,
	z shared.SpeechListenStopRequest) (shared.SpeechListenStopResponse, error) {

	if a.nasluchy == nil {
		return shared.SpeechListenStopResponse{Stopped: false}, nil
	}
	zdjeto := a.nasluchy.zdejmij(strings.TrimSpace(wartoscTekstu(z.ListenerId)),
		strings.TrimSpace(wartoscTekstu(z.WindowId)))
	return shared.SpeechListenStopResponse{Stopped: zdjeto}, nil
}

func (a *adapterMowy) ogloszOdcinekNasluchu(ctx context.Context, okno, odnosnik string) {
	if a.nasluchy == nil || a.nadajnik == nil || strings.TrimSpace(okno) == "" {
		return
	}
	nasluch, jest := a.nasluchy.nasluchOkna(strings.TrimSpace(okno))
	if !jest {
		return
	}

	zadanie := shared.SpeechTranscribeRequest{AudioRef: odnosnik}
	if nasluch.Jezyk != "" {
		jezyk := nasluch.Jezyk
		zadanie.Language = &jezyk
	}
	wynik, err := a.Przepisz(ctx, zadanie)
	if err != nil {
		return
	}
	tekst := strings.TrimSpace(wynik.Transcript)
	if tekst == "" {
		return
	}

	a.nadajnik.czesciowaTranskrypcja(ctx, nasluch, tekst, wynik.Confidence)

	// Fraza wybudzająca jest rozpoznawana na tekście, nie osobnym modelem słowa kluczowego.
	nastawa := a.odczytajNastaweWybudzania(ctx, nil, nil)
	fraza := strings.TrimSpace(nastawa.Phrase)
	if fraza == "" || nasluch.Tryb == shared.ListenModePushToTalk {
		return
	}
	if strings.Contains(strings.ToLower(tekst), strings.ToLower(fraza)) {
		a.nadajnik.wykrycieFrazyWybudzajacej(ctx, nasluch, fraza)
	}
}

func (a *adapterMowy) odczytajNastaweWybudzania(ctx context.Context, poziom *shared.ConfigScope,
	kluczZasiegu *string) shared.WakeWordConfig {

	nastawa := shared.WakeWordConfig{
		Mode:  shared.ListenModePushToTalk,
		Scope: shared.ConfigScopeGlobal,
	}
	if poziom != nil && strings.TrimSpace(string(*poziom)) != "" {
		nastawa.Scope = *poziom
	}
	if a.rozstrzygacz == nil {
		return nastawa
	}
	kontekst := kontekstZasieguMowy(ctx)
	_ = kluczZasiegu

	nastawa.Phrase = strings.TrimSpace(
		a.rozstrzygacz.Rozstrzygnij(kontekst, kluczFrazyWybudzajacej).Wartosc)
	if tryb := strings.TrimSpace(
		a.rozstrzygacz.Rozstrzygnij(kontekst, kluczTrybuNasluchu).Wartosc); tryb != "" {
		nastawa.Mode = shared.ListenMode(tryb)
	}
	if prog, err := strconv.Atoi(strings.TrimSpace(
		a.rozstrzygacz.Rozstrzygnij(kontekst, kluczProguDetekcji).Wartosc)); err == nil {
		nastawa.VadThreshold = prog
	}
	switch strings.ToLower(strings.TrimSpace(
		a.rozstrzygacz.Rozstrzygnij(kontekst, kluczOdszumiania).Wartosc)) {
	case "1", "true", "tak":
		nastawa.NoiseSuppression = true
	}
	nastawa.UpdatedAt = time.Now().UnixMilli()
	return nastawa
}

func (a *adapterMowy) zapiszNastaweMowy(ctx context.Context, poziom, kluczZasiegu,
	klucz, wartosc string) error {

	kopia := wartosc
	err := a.konfiguracja.Ustaw(ctx, dane.Ustawienie{
		Poziom:       shared.ConfigScope(poziom),
		KluczZasiegu: kluczZasiegu,
		Klucz:        klucz,
		Wartosc:      &kopia,
	})
	if err != nil {
		return bladZapleczaNagran("nie można zapisać nastawy " + klucz + ": " + err.Error())
	}
	if a.rozstrzygacz != nil {
		a.rozstrzygacz.Oglos(klucz)
	}
	return nil
}

func sprawdzTrybNasluchu(tryb shared.ListenMode) error {
	switch tryb {
	case shared.ListenModePushToTalk, shared.ListenModeWakeWord, shared.ListenModeContinuous:
		return nil
	default:
		return bladWskazaniaNagrania("nie znam trybu nasłuchu „" + string(tryb) +
			"” — serwer zna: pushToTalk, wakeWord, continuous")
	}
}

func kontekstZasieguMowy(ctx context.Context) konfig.Kontekst {
	return konfig.Kontekst{KontoOperatora: dane.KontoOperatora(ctx)}
}

func (e *emiter) czesciowaTranskrypcja(ctx context.Context, n nasluchMowy, tekst string, pewnosc *int) {
	// `final` jest prawdą, bo odcinek został rozpoznany w całości, a nie jako strumień w locie.
	e.wyslijDoKonta(ctx, shared.EventSpeechListenPartial, n.Sesja, shared.SpeechListenPartialEvent{
		ListenerId: n.Kod, WindowId: n.Okno, Transcript: tekst, Final: true,
		Confidence: pewnosc,
	})
}

func (e *emiter) wykrycieFrazyWybudzajacej(ctx context.Context, n nasluchMowy, fraza string) {
	e.wyslijDoKonta(ctx, shared.EventSpeechWakeDetected, n.Sesja, shared.SpeechWakeDetectedEvent{
		ListenerId: n.Kod, WindowId: n.Okno, Phrase: fraza,
		DetectedAt: time.Now().UnixMilli(),
	})
}

// _ trzyma import pakietu `mowa`: nastawy tego pliku opisują silnik transkrypcji odcinka.
var _ = mowa.FormatyNagran
