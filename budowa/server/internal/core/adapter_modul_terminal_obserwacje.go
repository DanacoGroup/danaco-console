// Komendy terminal.watch.start, terminal.watch.stop i terminal.watch.list:
// obserwacje plików uruchamiające polecenie karty terminala przy zmianie ścieżek.
package core

import (
	"context"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const przedrostekObserwacji = "twch-"

const odstepPrzegladu = 500 * time.Millisecond

// Pół sekundy scala zapis pliku wykonany przez edytor w kilku ruchach w jedno wyzwolenie.
const tlumienieDomyslne = 500 * time.Millisecond

// Obserwacja na korzeniu wielkiego drzewa ma tłumić samą siebie, nie zająć maszyny chodzeniem po katalogach.
const granicaPrzegladu = 20000

type obserwacjaZywa struct {
	kod       string
	zatrzymaj context.CancelFunc
	koniec    chan struct{}
}

type rejestrObserwacji struct {
	mu    sync.Mutex
	wpisy map[string]*obserwacjaZywa
}

func nowyRejestrObserwacji() *rejestrObserwacji {
	return &rejestrObserwacji{wpisy: make(map[string]*obserwacjaZywa)}
}

func (r *rejestrObserwacji) zapisz(o *obserwacjaZywa) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.wpisy[o.kod] = o
}

func (r *rejestrObserwacji) wez(kod string) (*obserwacjaZywa, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, jest := r.wpisy[kod]
	return o, jest
}

func (r *rejestrObserwacji) zdejmij(kod string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.wpisy, kod)
}

func (r *rejestrObserwacji) zatrzymajWszystkie() {
	r.mu.Lock()
	wykaz := make([]*obserwacjaZywa, 0, len(r.wpisy))
	for _, wpis := range r.wpisy {
		wykaz = append(wykaz, wpis)
	}
	r.mu.Unlock()
	for _, wpis := range wykaz {
		wpis.zatrzymaj()
	}
}

func (a *adapterTerminala) ZalozObserwacje(ctx context.Context,
	z shared.TerminalWatchStartRequest) (shared.TerminalWatchStartResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalWatchStartResponse{}, err
	}
	karta, err := a.kartaZadania(z.SessionId)
	if err != nil {
		return shared.TerminalWatchStartResponse{}, err
	}
	wzorzec := strings.TrimSpace(z.Pattern)
	polecenie := strings.TrimSpace(z.Command)
	if wzorzec == "" || polecenie == "" {
		return shared.TerminalWatchStartResponse{}, bladZadaniaTerminala(
			"obserwacja wymaga wzorca ścieżek i polecenia uruchamianego po zmianie")
	}
	okno, err := a.oknoWykonania(karta.oknoKod)
	if err != nil {
		return shared.TerminalWatchStartResponse{}, err
	}
	// Uprawnienie sprawdza się przy zakładaniu, nie wyzwoleniu, by odmowa nie przyszła bez czytelnika.
	if err := sprawdzUprawnienie(okno.TrybUprawnien, shared.ProcessInitiatorOperator); err != nil {
		return shared.TerminalWatchStartResponse{}, err
	}
	korzen, err := a.sciezkaWObszarze(karta, okno, korzenWzorca(wzorzec))
	if err != nil {
		return shared.TerminalWatchStartResponse{}, err
	}

	wiersz := dane.ObserwacjaTerminala{
		Kod:           nowyIdentyfikator(przedrostekObserwacji),
		OknoKod:       karta.oknoKod,
		KartaKod:      karta.kod,
		Wzorzec:       wzorzec,
		Polecenie:     polecenie,
		Rekurencyjnie: z.Recursive != nil && *z.Recursive,
		Stan:          shared.TerminalWatchStatusActive,
	}
	if z.DebounceMs != nil && *z.DebounceMs >= 0 {
		tlumienie := int64(*z.DebounceMs)
		wiersz.Tlumienie = &tlumienie
	}
	if err := dziennik.ZapiszObserwacje(ctx, wiersz); err != nil {
		return shared.TerminalWatchStartResponse{}, err
	}

	// Pętla przeglądu żyje poza żądaniem, aż do zatrzymania albo końca pracy
	// rdzenia; konto zamawiającego jedzie z nią (decyzja 34).
	zycie, zatrzymaj := context.WithCancel(dane.ZKontemOperatora(context.Background(), dane.KontoOperatora(ctx)))
	zywa := &obserwacjaZywa{kod: wiersz.Kod, zatrzymaj: zatrzymaj, koniec: make(chan struct{})}
	a.obserwacje.zapisz(zywa)
	go a.przegladajObserwacje(zycie, zywa, wiersz, karta, korzen)

	zapisana, err := dziennik.Obserwacja(ctx, wiersz.Kod)
	if err != nil {
		return shared.TerminalWatchStartResponse{}, err
	}
	return shared.TerminalWatchStartResponse{Watch: obserwacjaKontraktu(zapisana)}, nil
}

func (a *adapterTerminala) ZatrzymajObserwacje(ctx context.Context,
	z shared.TerminalWatchStopRequest) (shared.TerminalWatchStopResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalWatchStopResponse{}, err
	}
	kod := strings.TrimSpace(z.WatchId)
	if kod == "" {
		return shared.TerminalWatchStopResponse{}, bladZadaniaTerminala(
			"zatrzymanie obserwacji wymaga jej wskazania")
	}
	if _, err := dziennik.Obserwacja(ctx, kod); err != nil {
		return shared.TerminalWatchStopResponse{}, bladBrakuZasobuTerminala("obserwacja " + kod)
	}
	if zywa, jest := a.obserwacje.wez(kod); jest {
		zywa.zatrzymaj()
		select {
		case <-zywa.koniec:
		case <-time.After(czasNaDomkniecie):
		}
	}
	if err := dziennik.ZmienStanObserwacji(ctx, kod, shared.TerminalWatchStatusStopped, ""); err != nil {
		return shared.TerminalWatchStopResponse{}, err
	}
	zatrzymana, err := dziennik.Obserwacja(ctx, kod)
	if err != nil {
		return shared.TerminalWatchStopResponse{}, err
	}
	return shared.TerminalWatchStopResponse{Watch: obserwacjaKontraktu(zatrzymana)}, nil
}

func (a *adapterTerminala) WykazObserwacji(ctx context.Context,
	z shared.TerminalWatchListRequest) (shared.TerminalWatchListResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalWatchListResponse{}, err
	}
	filtr := dane.FiltrObserwacjiTerminala{OknoKod: strings.TrimSpace(wartoscTekstu(z.WindowId))}
	if z.Status != nil {
		filtr.Stan = *z.Status
	}
	wiersze, err := dziennik.Obserwacje(ctx, filtr)
	if err != nil {
		return shared.TerminalWatchListResponse{}, err
	}
	wykaz := make([]shared.TerminalWatch, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wykaz = append(wykaz, obserwacjaKontraktu(wiersz))
	}
	return shared.TerminalWatchListResponse{Watches: wykaz, Total: len(wykaz)}, nil
}

// Pierwszy przegląd wyłącznie zapamiętuje stan: obserwacja reaguje na zmianę, nie na pliki zastane.
func (a *adapterTerminala) przegladajObserwacje(zycie context.Context, zywa *obserwacjaZywa,
	wiersz dane.ObserwacjaTerminala, karta *kartaTerminala, korzen string) {

	defer close(zywa.koniec)
	defer a.obserwacje.zdejmij(wiersz.Kod)

	tlumienie := tlumienieDomyslne
	if wiersz.Tlumienie != nil {
		tlumienie = time.Duration(*wiersz.Tlumienie) * time.Millisecond
	}
	poprzedni := a.migawkaSciezek(wiersz, korzen)
	ostatnie := time.Time{}

	zegar := time.NewTicker(odstepPrzegladu)
	defer zegar.Stop()
	for {
		select {
		case <-zycie.Done():
			return
		case <-zegar.C:
		}
		biezacy := a.migawkaSciezek(wiersz, korzen)
		if !migawkiRozne(poprzedni, biezacy) {
			continue
		}
		poprzedni = biezacy
		if time.Since(ostatnie) < tlumienie {
			continue
		}
		ostatnie = time.Now()
		a.wyzwolObserwacje(zycie, wiersz, karta)
	}
}

// Niepowodzenie przestawia obserwację na failed i kończy pętlę: odrzucone polecenie nie ma próbować bez końca.
func (a *adapterTerminala) wyzwolObserwacje(zycie context.Context,
	wiersz dane.ObserwacjaTerminala, karta *kartaTerminala) {

	okno, err := a.oknoWykonania(karta.oknoKod)
	if err == nil {
		dopuszczone, bladIzolacji := a.polecenieDopuszczone(karta, okno, wiersz.Polecenie)
		if bladIzolacji != nil {
			err = bladIzolacji
		} else {
			proces, bladStartu := a.uruchom(zycie, karta, okno, dopuszczone,
				wiersz.Polecenie, shared.ProcessInitiatorOperator)
			if bladStartu != nil {
				err = bladStartu
			} else {
				a.pilnujZakonczenia(proces, granicaCzasuDomyslna)
			}
		}
	}
	if err != nil {
		if a.repozytorium != nil {
			_ = a.repozytorium.ZmienStanObserwacji(dane.ZKontemOperatora(context.Background(),
				dane.KontoOperatora(zycie)), wiersz.Kod, shared.TerminalWatchStatusFailed, err.Error())
		}
		if zywa, jest := a.obserwacje.wez(wiersz.Kod); jest {
			zywa.zatrzymaj()
		}
		return
	}
	if a.repozytorium != nil {
		_ = a.repozytorium.OdnotujWyzwolenie(dane.ZKontemOperatora(context.Background(),
			dane.KontoOperatora(zycie)), wiersz.Kod)
	}
}

// Para czasu zmiany i rozmiaru jest jedyną zmianą pliku widoczną bez czytania jego treści.
func (a *adapterTerminala) migawkaSciezek(wiersz dane.ObserwacjaTerminala, korzen string) map[string]string {
	migawka := make(map[string]string, 64)
	wzorzec := filepath.Base(wiersz.Wzorzec)
	_ = filepath.WalkDir(korzen, func(sciezka string, wpis fs.DirEntry, err error) error {
		if err != nil {
			// Katalog bez dostępu jest pomijany, nie przerywa obserwacji — drzewo projektu bywa niejednorodne.
			if wpis != nil && wpis.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if len(migawka) >= granicaPrzegladu {
			return filepath.SkipAll
		}
		if wpis.IsDir() {
			if sciezka != korzen && !wiersz.Rekurencyjnie {
				return fs.SkipDir
			}
			return nil
		}
		if !pasujeWzorzecObserwacji(wzorzec, filepath.Base(sciezka)) {
			return nil
		}
		opis, err := wpis.Info()
		if err != nil {
			return nil
		}
		migawka[sciezka] = opis.ModTime().UTC().Format(time.RFC3339Nano) + ":" +
			strconv.FormatInt(opis.Size(), 10)
		return nil
	})
	return migawka
}

func migawkiRozne(przed, po map[string]string) bool {
	if len(przed) != len(po) {
		return true
	}
	for sciezka, znacznik := range po {
		if przed[sciezka] != znacznik {
			return true
		}
	}
	return false
}

func korzenWzorca(wzorzec string) string {
	oczyszczony := filepath.FromSlash(strings.TrimSpace(wzorzec))
	katalog := filepath.Dir(oczyszczony)
	if katalog == "." || katalog == "" {
		return "."
	}
	// Katalog z gwiazdką jest wzorcem, nie katalogiem: przegląd zaczyna się od części stałej.
	for strings.ContainsAny(filepath.Base(katalog), "*?[") {
		rodzic := filepath.Dir(katalog)
		if rodzic == katalog {
			return "."
		}
		katalog = rodzic
	}
	return katalog
}

func pasujeWzorzecObserwacji(wzorzec, nazwa string) bool {
	if wzorzec == "" || wzorzec == "*" || wzorzec == "." {
		return true
	}
	pasuje, err := filepath.Match(wzorzec, nazwa)
	// Wzorzec niepoprawny przepuszcza wszystko zamiast nic, żeby milczał wyzwoleniami, a nie własną wadą.
	return err != nil || pasuje
}

func obserwacjaKontraktu(w dane.ObserwacjaTerminala) shared.TerminalWatch {
	obserwacja := shared.TerminalWatch{
		Id:        w.Kod,
		WindowId:  w.OknoKod,
		SessionId: w.KartaKod,
		Pattern:   w.Wzorzec,
		Command:   w.Polecenie,
		Status:    w.Stan,
		CreatedAt: chwilaBazy(w.Zalozono),
	}
	licznik := int(w.Licznik)
	obserwacja.TriggerCount = &licznik
	rekurencyjnie := w.Rekurencyjnie
	obserwacja.Recursive = &rekurencyjnie
	obserwacja.DebounceMs = wskaznikMalej(w.Tlumienie)
	if w.Wyzwolono != nil {
		chwila := chwilaBazy(*w.Wyzwolono)
		obserwacja.LastTriggeredAt = &chwila
	}
	if w.Powod != "" {
		powod := w.Powod
		obserwacja.ErrorMessage = &powod
	}
	return obserwacja
}
