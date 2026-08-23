// Komendy `terminal.watch.start`, `terminal.watch.stop` i `terminal.watch.list`
// — obserwacje plików uruchamiające polecenie karty przy ich zmianie.
//
// ── Czego brakowało ─────────────────────────────────────────────────────────
// Wyzwalacz plikowy istniał w kontrakcie wyłącznie dla automatyk, czyli poza
// powłoką. Nie dało się powiedzieć „po każdej zmianie w tym katalogu zbuduj
// projekt W TEJ karcie, w jej katalogu, jej powłoką i jej środowiskiem”.
//
// ── Dlaczego przegląd, a nie zdarzenia systemu plików ───────────────────────
// Zdarzenia jądra (inotify, ReadDirectoryChangesW) wymagałyby nowej zależności
// modułowej i osobnej implementacji na każdy system. Przegląd po czasach zmiany
// jest wkompilowany w bibliotekę standardową, zachowuje się jednakowo wszędzie
// i jest dokładnie tak dokładny, jak trzeba: obserwacja i tak TŁUMI powtórzenia
// (`debounceMs`), więc rozdzielczość poniżej tłumienia byłaby wyrzucona przez
// samo tłumienie. Cena — przegląd katalogu co odstęp — jest znikoma wobec
// polecenia, które ten przegląd wyzwala.
//
// ── Skutkiem obserwacji jest PROCES, nie zapis ──────────────────────────────
// Wyzwolenie idzie tą samą drogą co `terminal.command.exec`: brama trybu
// uprawnień okna, egzekutor izolacji, port `session.Uruchamiacz`, wpis w Process
// Monitorze i strumień do Output Console. Obserwacja nie jest drugą drogą
// uruchamiania procesów — jest wyzwalaczem tej jedynej.
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

// przedrostekObserwacji znakuje identyfikator obserwacji plików.
const przedrostekObserwacji = "twch-"

// odstepPrzegladu jest rytmem przeglądania obserwowanych ścieżek.
const odstepPrzegladu = 500 * time.Millisecond

// tlumienieDomyslne obowiązuje, gdy żądanie nie poda własnego. Pół sekundy
// scala zapis pliku wykonany przez edytor w kilku ruchach (plik tymczasowy,
// zamiana nazwy) w jedno wyzwolenie.
const tlumienieDomyslne = 500 * time.Millisecond

// granicaPrzegladu ogranicza liczbę plików objętych jednym przeglądem.
// Obserwacja założona na korzeniu wielkiego drzewa ma tłumić samą siebie, a nie
// zająć maszynę chodzeniem po katalogach.
const granicaPrzegladu = 20000

// obserwacjaZywa jest jedną biegnącą obserwacją wraz z drogą jej zatrzymania.
type obserwacjaZywa struct {
	kod       string
	zatrzymaj context.CancelFunc
	koniec    chan struct{}
}

// rejestrObserwacji trzyma obserwacje czynne jednego biegu rdzenia.
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

// zatrzymajWszystkie kończy obserwacje przy zatrzymaniu rdzenia.
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

// ZalozObserwacje obsługuje `terminal.watch.start`.
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
	// Bramę trybu uprawnień sprawdzamy TERAZ, a nie dopiero przy wyzwoleniu:
	// obserwacja założona w oknie, które i tak nie może uruchomić procesu, byłaby
	// obietnicą bez pokrycia, a odmowa przyszłaby po pierwszej zmianie pliku,
	// czyli w chwili, w której nikt jej nie czyta.
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

	// Pętla przeglądu żyje poza żądaniem: obserwacja przeżywa komendę, która ją
	// założyła, i biegnie do zatrzymania albo do końca pracy rdzenia.
	zycie, zatrzymaj := context.WithCancel(context.Background())
	zywa := &obserwacjaZywa{kod: wiersz.Kod, zatrzymaj: zatrzymaj, koniec: make(chan struct{})}
	a.obserwacje.zapisz(zywa)
	go a.przegladajObserwacje(zycie, zywa, wiersz, karta, korzen)

	zapisana, err := dziennik.Obserwacja(ctx, wiersz.Kod)
	if err != nil {
		return shared.TerminalWatchStartResponse{}, err
	}
	return shared.TerminalWatchStartResponse{Watch: obserwacjaKontraktu(zapisana)}, nil
}

// ZatrzymajObserwacje obsługuje `terminal.watch.stop`. Polecenie już uruchomione
// biegnie dalej — jest procesem rejestru rdzenia i kończy się własną drogą.
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

// WykazObserwacji obsługuje `terminal.watch.list`.
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

// przegladajObserwacje prowadzi pętlę przeglądu jednej obserwacji.
//
// Pierwszy przegląd wyłącznie ZAPAMIĘTUJE stan i niczego nie wyzwala: gdyby
// wyzwalał, założenie obserwacji uruchamiałoby polecenie natychmiast, na plikach
// zastanych, a Operator prosił o reakcję na ZMIANĘ.
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

// wyzwolObserwacje uruchamia polecenie obserwacji w karcie.
//
// Niepowodzenie uruchomienia przestawia obserwację na `failed` wraz z powodem
// i KOŃCZY pętlę. Obserwacja, która przy każdej zmianie próbuje uruchomić
// polecenie odrzucone przez izolację, robiłaby to bez końca i bez skutku.
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
			_ = a.repozytorium.ZmienStanObserwacji(context.Background(), wiersz.Kod,
				shared.TerminalWatchStatusFailed, err.Error())
		}
		if zywa, jest := a.obserwacje.wez(wiersz.Kod); jest {
			zywa.zatrzymaj()
		}
		return
	}
	if a.repozytorium != nil {
		_ = a.repozytorium.OdnotujWyzwolenie(context.Background(), wiersz.Kod)
	}
}

// migawkaSciezek zbiera czasy zmiany i rozmiary plików objętych obserwacją.
//
// Kluczem porównania jest para (czas zmiany, rozmiar), nie sama data: zapis,
// który zmienia treść w tej samej sekundzie, zmienia zwykle rozmiar, a zapis
// zmieniający ani jednego, ani drugiego nie zmienia też pliku w żaden sposób
// widoczny dla polecenia, które ma się po nim wykonać.
func (a *adapterTerminala) migawkaSciezek(wiersz dane.ObserwacjaTerminala, korzen string) map[string]string {
	migawka := make(map[string]string, 64)
	wzorzec := filepath.Base(wiersz.Wzorzec)
	_ = filepath.WalkDir(korzen, func(sciezka string, wpis fs.DirEntry, err error) error {
		if err != nil {
			// Katalog, do którego nie ma dostępu, pomija się zamiast przerywać
			// całą obserwację: drzewo projektu bywa niejednorodne.
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

// migawkiRozne porównuje dwie migawki przeglądu.
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

// korzenWzorca wyławia ze wzorca część katalogową — tę, od której zaczyna się
// przegląd. Wzorzec bez katalogu obserwuje katalog roboczy karty.
func korzenWzorca(wzorzec string) string {
	oczyszczony := filepath.FromSlash(strings.TrimSpace(wzorzec))
	katalog := filepath.Dir(oczyszczony)
	if katalog == "." || katalog == "" {
		return "."
	}
	// Katalog z gwiazdką nie jest katalogiem, tylko wzorcem katalogów: przegląd
	// zaczyna się wtedy od jego części stałej.
	for strings.ContainsAny(filepath.Base(katalog), "*?[") {
		rodzic := filepath.Dir(katalog)
		if rodzic == katalog {
			return "."
		}
		katalog = rodzic
	}
	return katalog
}

// pasujeWzorzecObserwacji sprawdza nazwę pliku wzorcem powłoki. Wzorzec pusty
// albo `*` przepuszcza wszystko.
func pasujeWzorzecObserwacji(wzorzec, nazwa string) bool {
	if wzorzec == "" || wzorzec == "*" || wzorzec == "." {
		return true
	}
	pasuje, err := filepath.Match(wzorzec, nazwa)
	// Wzorzec niepoprawny przepuszcza wszystko zamiast nic: obserwacja ma
	// wyzwalać, a nie milczeć o własnej wadzie przez brak wyzwoleń.
	return err != nil || pasuje
}

// obserwacjaKontraktu przekłada wiersz obserwacji na byt kontraktu.
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
