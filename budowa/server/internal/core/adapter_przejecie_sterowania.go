// Przejęcie bezpośredniego sterowania — rodzina komend `control.*`, którą
// Operator odbiera prowadzenie zlecenia pętli koordynator–wykonawca i oddaje je
// z powrotem.
//
// Stan bieżący („kto steruje teraz") żyje w pamięci procesu, bo dotyczy biegu
// żywego, a bieg żywy nie przeżywa restartu rdzenia. Ślad („kto i kiedy przejął")
// jest faktem historycznym i idzie do dziennika akcji okna `log_akcji_okna`.
// Po restarcie zlecenie wraca pod Koordynatora, a ślad przejęcia zostaje czytelny.
//
// Plik nie prowadzi biegu: nie liczy obiegów, nie wykrywa braku postępu, nie
// rozpoczyna tur — robi to wyłącznie pętla. Nie zapisuje też do `pozycja_kolejki`:
// wznowienie drogą `queue.action resume` woła `krok()`, a ten przesuwa pozycję po
// mapie `krokNaprzod`, czyli pracę przerwaną w stanie „wykonywana" traktuje jak
// skończoną.
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

// Sterujacy nazywa rękę prowadzącą zlecenie.
//
// Kontrakt nie ma wyliczenia sterującego, więc katalog wartości mieszka w rdzeniu.
type Sterujacy string

const (
	// SterujacyKoordynator — zlecenie prowadzi pętla koordynator–wykonawca.
	// Stan domyślny: brak przejęcia ma wyglądać jak Koordynator, nie jak brak
	// wiedzy.
	SterujacyKoordynator Sterujacy = "coordinator"
	// SterujacyOperator — sterowanie trzyma człowiek.
	SterujacyOperator Sterujacy = "operator"
)

// SterZlecenia jest odpowiedzią na pytanie „kto prowadzi to zlecenie TERAZ".
type SterZlecenia struct {
	// Okno koordynatora, którego zlecenia dotyczy ster.
	Okno string
	// Sterujacy — Koordynator albo Operator.
	Sterujacy Sterujacy
	// PrzejeteO — chwila przejęcia w milisekundach epoki; zero, gdy prowadzi
	// Koordynator.
	PrzejeteO int64
	// PrzejalKlient — identyfikator klienta Operatora, który przejął. Pusty
	// znaczy „rdzeń nie potrafił tego rozstrzygnąć", tak samo jak `actorClientId`
	// w zdarzeniach.
	PrzejalKlient string
}

// sterKoordynatora składa stan domyślny — zlecenie w rękach pętli.
func sterKoordynatora(okno string) SterZlecenia {
	return SterZlecenia{Okno: okno, Sterujacy: SterujacyKoordynator}
}

// rejestrSteru pamięta bieżącego sterującego każdego zlecenia z osobna.
//
// Byt równoległy do `rejestrBiegow` i z tego samego powodu: pętla zna swój stan
// w chwili obiegu i wydaje go jednorazowo, a pytanie „kto steruje tym zleceniem"
// pada z zewnątrz, w dowolnej chwili. Rejestr nie prowadzi biegu i nie zatrzymuje go.
type rejestrSteru struct {
	mu    sync.RWMutex
	stany map[string]SterZlecenia
}

// nowyRejestrSteru zakłada pusty rejestr steru.
func nowyRejestrSteru() *rejestrSteru {
	return &rejestrSteru{stany: map[string]SterZlecenia{}}
}

// ustaw zapisuje bieżącego sterującego.
func (r *rejestrSteru) ustaw(s SterZlecenia) {
	if r == nil || s.Okno == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stany[s.Okno] = s
}

// zdejmij kasuje wpis — zlecenie wraca pod Koordynatora.
func (r *rejestrSteru) zdejmij(okno string) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.stany, okno)
}

// stan zwraca bieżącego sterującego. Brak wpisu daje fałsz, a nie ster pusty —
// wołający sam rozstrzyga, czy „nikt nie przejął" znaczy dla niego Koordynatora.
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
// wykreślonych. Wykaz pusty niczego nie kasuje — brak wiedzy o oknach nie jest
// wiedzą o ich zamknięciu. Wzorowane na `rejestrBiegow.Zachowaj`.
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

// sterBiegow jest jedynym rejestrem steru w procesie rdzenia.
//
// Byt pakietowy, a nie pole struktury: odpis licznika obiegów staje się
// `LoopState` kontraktu w wolnej funkcji `biegKontraktu(s session.StanObiegu)`
// z `core/stan_obiegu.go`, której jedynym argumentem jest ten odpis. Żeby pola
// `controller`, `takenOverAt` i `takenOverBy` mogły wyjść obiema rodzinami
// zdarzeń niosącymi `LoopState`, funkcja musi sięgnąć po ster bez zmiany swojej
// sygnatury. Klucz rejestru jest niepowtarzalny w procesie (identyfikator
// zewnętrzny okna), więc jeden rejestr na proces jest poprawny.
var sterBiegow = nowyRejestrSteru()

// sterZlecenia odpowiada, kto prowadzi zlecenie wskazanego koordynatora.
//
// Woła ją przekład licznika obiegów na `LoopState` w `core/stan_obiegu.go`.
// Fałsz znaczy „nikt nie przejmował" i po stronie przekładu schodzi na
// `coordinator`, nigdy na pustkę.
func sterZlecenia(idKoordynatora string) (SterZlecenia, bool) {
	return sterBiegow.stan(idKoordynatora)
}

// zachowajStery wykreśla stery okien już zamkniętych. Sprzątanie idzie od strony
// okien żywych, bo pętla o zamknięciu okna nie mówi.
func zachowajStery(okna map[string]struct{}) int {
	return sterBiegow.zachowaj(okna)
}

// adapterPrzejeciaSterowania obsługuje rodzinę `control.*`.
//
// Rejestr okien mówi, czy okno prowadzi zlecenie; pętla wykonuje zatrzymanie
// i wznowienie; repozytorium przekazań pisze i czyta ślad w dzienniku akcji okna.
type adapterPrzejeciaSterowania struct {
	okna         *session.Rejestr
	petla        *session.Petla
	repozytorium dane.RepozytoriumPrzekazan
	ster         *rejestrSteru
}

// nowyAdapterPrzejeciaSterowania składa adapter czwartej drogi interwencji.
// Rejestr steru jest wspólny dla procesu — patrz komentarz przy `sterBiegow`.
func nowyAdapterPrzejeciaSterowania(okna *session.Rejestr, petla *session.Petla,
	repozytorium dane.RepozytoriumPrzekazan) *adapterPrzejeciaSterowania {
	return &adapterPrzejeciaSterowania{
		okna: okna, petla: petla, repozytorium: repozytorium, ster: sterBiegow,
	}
}

// Przejmij wykonuje `control.takeover`: Operator bierze zlecenie w swoje ręce.
//
// Kolejność czynności jest treścią, nie stylem:
//  1. ślad idzie do dziennika pierwszy, bo tylko on może się nie udać —
//     przejęcie bez zapisu zostałoby bez świadka;
//  2. ster wchodzi do rejestru przed zatrzymaniem, żeby rozgłoszenie stanu biegu
//     wywołane zatrzymaniem niosło już nowego sterującego;
//  3. `petla.Zatrzymaj` staje na końcu i nie kasuje niczego z dorobku
//     Koordynatora — Operator wchodzi tam, gdzie proces stoi.
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

// Oddaj wykonuje `control.release`: Operator zwraca zlecenie Koordynatorowi.
//
// `Petla.Wznow` kasuje wyłącznie licznik braku postępu i ostatni odcisk
// strumienia: historia obiegów zostaje, więc Koordynator podejmuje bieg, a nie
// zaczyna go od nowa.
//
// Oddanie zlecenia, którego nikt nie przejął, jest odmawiane nazwanym błędem.
// Cicha zgoda wznowiłaby bieg, którego Operator nie zatrzymywał — czyli zmiana
// stanu, o którą nikt nie prosił.
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

// Ster wykonuje `control.get`: kto prowadzi zlecenie i kto je dotąd przejmował.
//
// Okno bez ani jednego przejęcia oddaje `coordinator` i wykaz pusty.
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

// sprawdzZlecenie odmawia wszystkiemu, co nie jest zleceniem prowadzonym.
//
// Sprawdzana jest żywa rola okna z rejestru nadzorcy, a nie kolumna w bazie —
// przejęcie dotyczy biegu żywego, a więź koordynatora bywa w bazie pusta.
func (a *adapterPrzejeciaSterowania) sprawdzZlecenie(idOkna string) error {
	if idOkna == "" {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"przejęcie sterowania: żądanie bez identyfikatora okna"))
	}
	if a.petla == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: rdzeń nie niesie pętli koordynator–wykonawca — "+
				"biegu nie ma czym zatrzymać ani wznowić"))
	}
	if a.repozytorium == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: rdzeń nie niesie dziennika akcji okna — "+
				"przejęcia nie ma gdzie odnotować"))
	}
	if a.okna == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: rdzeń nie niesie rejestru okien — "+
				"nie ma jak rozstrzygnąć, czy okno prowadzi zlecenie"))
	}
	okno, err := a.okna.Okno(idOkna)
	if err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"przejęcie sterowania: okno "+idOkna+" nie żyje na tym rdzeniu"))
	}
	if !okno.CzyKoordynator() {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"przejęcie sterowania: okno "+idOkna+" nie prowadzi zlecenia (rola "+
				string(okno.RolaOkna)+") — przejmuje się bieg Koordynatora, nie pojedyncze okno"))
	}
	return nil
}

// zapiszSlad odnotowuje przejęcie albo oddanie w dzienniku akcji okna.
//
// Ślad idzie dwiema kolumnami surowego zapisu, dokładnie tak, jak dziennik ich
// używa: `parametry` niosą to, o co poprosił Operator, `wynik` — odpis biegu
// z chwili czynności. Warstwa danych żadnej z nich nie rozbiera.
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

// historia odczytuje ślad sterowania i przekłada go z powrotem na zapisy.
//
// Repozytorium sięgane jest asercją typu po wąski interfejs `dane.RepozytoriumSteru`,
// bo szeroki kontrakt obszaru window.* deklaruje w całości inny plik. Port,
// który tej zdolności nie niesie, dostaje odmowę nazywającą brak — nigdy pusty
// wykaz udający „nikt nie przejmował".
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

// zapisZeSladu odtwarza zapis przejęcia z wiersza dziennika akcji.
//
// Wiersz nieczytelny zostaje w historii z rodzajem odczytanym z kolumny
// `akcja_id`, bo sam fakt „ktoś tu przejmował sterowanie" jest prawdziwy nawet
// wtedy, gdy szczegóły są uszkodzone.
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

// zapisJSON zamienia byt na surowy zapis kolumny dziennika.
func zapisJSON(v any) (*string, error) {
	tresc, err := json.Marshal(v)
	if err != nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"przejęcie sterowania: nie można złożyć zapisu śladu: "+err.Error()))
	}
	zapis := string(tresc)
	return &zapis, nil
}

// tekstWskaznika odczytuje pole opcjonalne jako napis; brak daje pustkę.
func tekstWskaznika(w *string) string {
	if w == nil {
		return ""
	}
	return *w
}

// wartoscWskaznika odczytuje liczbę opcjonalną; brak daje zero, które warstwa
// danych sprowadza do limitu domyślnego.
func wartoscWskaznika(w *int) int {
	if w == nil {
		return 0
	}
	return *w
}
