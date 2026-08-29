// Plik obsługuje rodzinę komend `control.*`, którą Operator odbiera i oddaje
// prowadzenie zlecenia pętli koordynator-wykonawca: stan bieżący żyje w
// pamięci procesu, ślad przejęcia idzie do dziennika akcji okna.
package core

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Sterujacy nazywa rękę prowadzącą zlecenie. Kontrakt nie ma wyliczenia
// sterującego, więc katalog wartości mieszka w rdzeniu.
type Sterujacy string

const (
	// SterujacyKoordynator — zlecenie prowadzi pętla koordynator-wykonawca,
	// stan domyślny: brak przejęcia ma wyglądać jak Koordynator, nie jak brak
	// wiedzy.
	SterujacyKoordynator Sterujacy = "coordinator"
	// SterujacyOperator — sterowanie trzyma człowiek, który przejął zlecenie
	// komendą `control.takeover` i zwraca je komendą `control.release`.
	SterujacyOperator Sterujacy = "operator"
)

// SterZlecenia jest odpowiedzią na pytanie „kto prowadzi to zlecenie teraz”,
// niesie okno, sterującego, chwilę przejęcia i klienta, który przejął.
type SterZlecenia struct {
	// Okno koordynatora, którego zlecenia dotyczy ster.
	Okno string
	// Sterujacy — Koordynator albo Operator.
	Sterujacy Sterujacy
	// PrzejeteO — chwila przejęcia w milisekundach epoki; zero, gdy prowadzi
	// Koordynator.
	PrzejeteO int64
	// PrzejalKlient — identyfikator klienta Operatora, który przejął, pusty gdy nierozstrzygalny.
	PrzejalKlient string
}

// sterKoordynatora składa stan domyślny: zlecenie w rękach pętli
// koordynator-wykonawca, bez przejęcia przez Operatora.
func sterKoordynatora(okno string) SterZlecenia {
	return SterZlecenia{Okno: okno, Sterujacy: SterujacyKoordynator}
}

// rejestrSteru pamięta bieżącego sterującego każdego zlecenia z osobna, bytem
// równoległym do rejestru biegów, ale odpowiadającym na pytanie z zewnątrz,
// w dowolnej chwili.
type rejestrSteru struct {
	mu    sync.RWMutex
	stany map[string]SterZlecenia
}

// nowyRejestrSteru zakłada pusty rejestr steru dla wszystkich zleceń procesu
// rdzenia, gotowy do zapisu i odczytu przez adapter przejęcia.
func nowyRejestrSteru() *rejestrSteru {
	return &rejestrSteru{stany: map[string]SterZlecenia{}}
}

// ustaw zapisuje bieżącego sterującego zlecenia wskazanego oknem koordynatora
// w rejestrze steru procesu.
func (r *rejestrSteru) ustaw(s SterZlecenia) {
	if r == nil || s.Okno == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stany[s.Okno] = s
}

// zdejmij kasuje wpis rejestru steru — zlecenie wraca pod prowadzenie
// Koordynatora po oddaniu przez Operatora.
func (r *rejestrSteru) zdejmij(okno string) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.stany, okno)
}

// stan zwraca bieżącego sterującego. Brak wpisu daje fałsz, a nie ster pusty —
// wołający sam rozstrzyga, czy „nikt nie przejął” znaczy dla niego Koordynatora.
func (r *rejestrSteru) stan(okno string) (SterZlecenia, bool) {
	if r == nil || okno == "" {
		return SterZlecenia{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, jest := r.stany[okno]
	return s, jest
}

// zachowaj zostawia wyłącznie stery okien wciąż otwartych i zwraca liczbę
// wykreślonych, wzorowane na sprzątaniu rejestru biegów tej samej pętli.
func (r *rejestrSteru) zachowaj(okna map[string]struct{}) int {
	if r == nil || len(okna) == 0 {
		return 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	wykreslone := 0
	for okno := range r.stany {
		if _, zyje := okna[okno]; zyje {
			continue
		}
		delete(r.stany, okno)
		wykreslone++
	}
	return wykreslone
}

// sterBiegow jest jedynym rejestrem steru w procesie rdzenia: bytem
// pakietowym, nie polem struktury, bo odpowiedź „kto prowadzi to zlecenie”
// jest jedna dla całego procesu.
var sterBiegow = nowyRejestrSteru()

// zachowajStery wykreśla stery okien już zamkniętych, idąc od strony okien
// żywych, bo pętla koordynator-wykonawca o zamknięciu okna nie mówi.
func zachowajStery(okna map[string]struct{}) int {
	return sterBiegow.zachowaj(okna)
}

// adapterPrzejeciaSterowania obsługuje rodzinę `control.*`: rejestr okien
// mówi, czy okno prowadzi zlecenie, pętla wykonuje zatrzymanie i wznowienie.
type adapterPrzejeciaSterowania struct {
	okna         *session.Rejestr
	petla        *session.Petla
	repozytorium dane.RepozytoriumPrzekazan
	ster         *rejestrSteru
}

// nowyAdapterPrzejeciaSterowania składa adapter czwartej drogi interwencji,
// ze wspólnym dla procesu rejestrem steru.
func nowyAdapterPrzejeciaSterowania(okna *session.Rejestr, petla *session.Petla,
	repozytorium dane.RepozytoriumPrzekazan) *adapterPrzejeciaSterowania {
	return &adapterPrzejeciaSterowania{
		okna: okna, petla: petla, repozytorium: repozytorium, ster: sterBiegow,
	}
}

// Przejmij wykonuje `control.takeover`: Operator bierze zlecenie w swoje ręce,
// a kolejność zapisu śladu, rejestru i zatrzymania pętli jest treścią, nie
// stylem.
func (a *adapterPrzejeciaSterowania) Przejmij(ctx context.Context,
	z ZadaniePrzejeciaSterowania) (WynikPrzejeciaSterowania, error) {
	if err := a.sprawdzZlecenie(z.WindowId); err != nil {
		return WynikPrzejeciaSterowania{}, err
	}
	_, klient := sprawca(ctx)
	zapis := ZapisSterowania{
		WindowId:   z.WindowId,
		Controller: SterujacyOperator,
		At:         time.Now().UTC().UnixMilli(),
		Actor:      klient,
		Reason:     z.Reason,
	}
	if err := a.zapiszSlad(ctx, dane.AkcjaPrzejeciaSterowania, zapis, a.petla.Stan(z.WindowId)); err != nil {
		return WynikPrzejeciaSterowania{}, err
	}

	a.ster.ustaw(SterZlecenia{
		Okno: z.WindowId, Sterujacy: SterujacyOperator,
		PrzejeteO: zapis.At, PrzejalKlient: tekstWskaznika(klient),
	})
	stan := a.petla.Zatrzymaj(z.WindowId)
	return WynikPrzejeciaSterowania{Loop: biegKontraktu(stan), Handover: zapis}, nil
}

// Oddaj wykonuje `control.release`: Operator zwraca zlecenie Koordynatorowi,
// a oddanie zlecenia, którego nikt nie przejął, jest odmawiane nazwanym
// błędem.
func (a *adapterPrzejeciaSterowania) Oddaj(ctx context.Context,
	z ZadanieOddaniaSterowania) (WynikOddaniaSterowania, error) {
	if err := a.sprawdzZlecenie(z.WindowId); err != nil {
		return WynikOddaniaSterowania{}, err
	}
	if biezacy, jest := a.ster.stan(z.WindowId); !jest || biezacy.Sterujacy != SterujacyOperator {
		return WynikOddaniaSterowania{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
			"przejęcie sterowania: zlecenia okna "+z.WindowId+
				" nie przejął nikt — nie ma czego oddać Koordynatorowi"))
	}
	_, klient := sprawca(ctx)
	zapis := ZapisSterowania{
		WindowId:   z.WindowId,
		Controller: SterujacyKoordynator,
		At:         time.Now().UTC().UnixMilli(),
		Actor:      klient,
		Reason:     z.Note,
	}
	if err := a.zapiszSlad(ctx, dane.AkcjaOddaniaSterowania, zapis, a.petla.Stan(z.WindowId)); err != nil {
		return WynikOddaniaSterowania{}, err
	}

	a.ster.zdejmij(z.WindowId)
	stan := a.petla.Wznow(z.WindowId)
	return WynikOddaniaSterowania{Loop: biegKontraktu(stan), Handover: zapis}, nil
}

// Ster wykonuje `control.get`: kto prowadzi zlecenie i kto je dotąd
// przejmował. Okno bez ani jednego przejęcia oddaje wartość Koordynatora
// i wykaz pusty.
func (a *adapterPrzejeciaSterowania) Ster(ctx context.Context,
	z ZadanieOdczytuSterowania) (WynikOdczytuSterowania, error) {
	if err := a.sprawdzZlecenie(z.WindowId); err != nil {
		return WynikOdczytuSterowania{}, err
	}
	biezacy, jest := a.ster.stan(z.WindowId)
	if !jest {
		biezacy = sterKoordynatora(z.WindowId)
	}
	historia, err := a.historia(ctx, z.WindowId, wartoscWskaznika(z.Limit))
	if err != nil {
		return WynikOdczytuSterowania{}, err
	}
	return WynikOdczytuSterowania{
		Loop:       biegKontraktu(a.petla.Stan(z.WindowId)),
		Controller: biezacy.Sterujacy,
		History:    historia,
	}, nil
}

// sprawdzZlecenie odmawia wszystkiemu, co nie jest zleceniem prowadzonym,
// sprawdzając żywą rolę okna z rejestru nadzorcy, a nie kolumnę w bazie.
func (a *adapterPrzejeciaSterowania) sprawdzZlecenie(idOkna string) error {
	if idOkna == "" {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"przejęcie sterowania: żądanie bez identyfikatora okna"))
	}
	if a.petla == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: serwer nie niesie pętli koordynator–wykonawca — "+
				"biegu nie ma czym zatrzymać ani wznowić"))
	}
	if a.repozytorium == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: serwer nie niesie dziennika akcji okna — "+
				"przejęcia nie ma gdzie odnotować"))
	}
	if a.okna == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: serwer nie niesie rejestru okien — "+
				"nie ma jak rozstrzygnąć, czy okno prowadzi zlecenie"))
	}
	okno, err := a.okna.Okno(idOkna)
	if err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"przejęcie sterowania: okno "+idOkna+" nie żyje na tym serwerze"))
	}
	if !okno.CzyKoordynator() {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"przejęcie sterowania: okno "+idOkna+" nie prowadzi zlecenia (rola "+
				string(okno.RolaOkna)+") — przejmuje się bieg Koordynatora, nie pojedyncze okno"))
	}
	return nil
}

// zapiszSlad odnotowuje przejęcie albo oddanie w dzienniku akcji okna,
// dwiema kolumnami surowego zapisu, których warstwa danych nie rozbiera.
func (a *adapterPrzejeciaSterowania) zapiszSlad(ctx context.Context, rodzaj string,
	zapis ZapisSterowania, stan session.StanObiegu) error {
	parametry, err := zapisJSON(zapis)
	if err != nil {
		return err
	}
	odpis, err := zapisJSON(biegKontraktu(stan))
	if err != nil {
		return err
	}
	_, err = a.repozytorium.ZapiszAkcje(ctx, dane.AkcjaOkna{
		Okno: zapis.WindowId, AkcjaID: rodzaj, Parametry: parametry, Wynik: odpis,
	})
	if err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: nie można odnotować czynności "+rodzaj+
				" okna "+zapis.WindowId+": "+err.Error()))
	}
	return nil
}

// historia odczytuje ślad sterowania i przekłada go z powrotem na zapisy,
// sięgając repozytorium asercją typu po wąski interfejs
// `dane.RepozytoriumSteru`.
func (a *adapterPrzejeciaSterowania) historia(ctx context.Context, okno string, limit int) ([]ZapisSterowania, error) {
	czytnik, niesie := a.repozytorium.(dane.RepozytoriumSteru)
	if !niesie {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: repozytorium okien nie umie czytać śladu sterowania — "+
				"historii przejęć nie ma skąd wziąć"))
	}
	wiersze, err := czytnik.SladySterowania(ctx, okno, limit)
	if err != nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: nie można odczytać historii okna "+okno+": "+err.Error()))
	}
	historia := make([]ZapisSterowania, 0, len(wiersze))
	for _, wiersz := range wiersze {
		historia = append(historia, zapisZeSladu(wiersz))
	}
	return historia, nil
}

// zapisZeSladu odtwarza zapis przejęcia z wiersza dziennika akcji. Wiersz
// nieczytelny zostaje w historii, bo sam fakt przejęcia jest prawdziwy nawet
// uszkodzony.
func zapisZeSladu(wiersz dane.AkcjaOkna) ZapisSterowania {
	zapis := ZapisSterowania{WindowId: wiersz.Okno, Controller: SterujacyKoordynator}
	if wiersz.AkcjaID == dane.AkcjaPrzejeciaSterowania {
		zapis.Controller = SterujacyOperator
	}
	if wiersz.Parametry == nil {
		return zapis
	}
	var odczytany ZapisSterowania
	if err := json.Unmarshal([]byte(*wiersz.Parametry), &odczytany); err != nil {
		return zapis
	}
	odczytany.WindowId = wiersz.Okno
	odczytany.Controller = zapis.Controller
	return odczytany
}

// zapisJSON zamienia byt przejęcia sterowania na surowy zapis kolumny
// dziennika akcji okna, gotowy do zapisania w repozytorium.
func zapisJSON(v any) (*string, error) {
	tresc, err := json.Marshal(v)
	if err != nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: nie można złożyć zapisu śladu: "+err.Error()))
	}
	zapis := string(tresc)
	return &zapis, nil
}

// tekstWskaznika odczytuje pole opcjonalne surowego zapisu jako napis; brak
// wartości daje pustkę zamiast odmowy.
func tekstWskaznika(w *string) string {
	if w == nil {
		return ""
	}
	return *w
}

// wartoscWskaznika odczytuje liczbę opcjonalną surowego zapisu; brak daje
// zero, które warstwa danych sprowadza do limitu domyślnego.
func wartoscWskaznika(w *int) int {
	if w == nil {
		return 0
	}
	return *w
}
