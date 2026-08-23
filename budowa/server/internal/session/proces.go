package session

import (
	"sync"
	"time"
)

// Proces jest procesem JEDNEGO okna komunikacji. Sesja prowadzi tyle
// procesów, ile ma otwartych okien, i wszystkie biegną równolegle.
//
// Proces NIE powstaje w tym pakiecie ANI NIE JEST STĄD URUCHAMIANY. Startuje go
// warstwa kanału własną drogą (injection), a sesja obejmuje proces już biegnący
// przez RejestrProcesow.Przejmij. Do sesji należy wyłącznie to, czego
// kanał nie umie: objęcie całego drzewa potomstwa jednym uchwytem systemowym —
// Job Object na Windows, grupa procesów na systemach uniksowych. Ubicie okna
// kończy zatem także wnuki, bez `taskkill` i bez innej zależności od narzędzi
// systemu.
type Proces struct {
	// IdOkna — klucz procesu w rejestrze procesów.
	IdOkna string
	// IdProcesu — identyfikator nadany przez rdzeń (telemetria).
	IdProcesu string
	// Uruchomiono — chwila objęcia procesu uchwytem sesji.
	Uruchomiono time.Time

	drzewo *DrzewoProcesu

	mu         sync.Mutex
	zakonczony bool
}

// Ubij kończy proces okna wraz z całym drzewem potomstwa i oddaje uchwyty
// systemowe. Ubicie jest jednorazowe: drugie wywołanie nie sięga po zwolniony
// już uchwyt.
func (p *Proces) Ubij() error {
	p.mu.Lock()
	if p.zakonczony {
		p.mu.Unlock()
		return nil
	}
	p.zakonczony = true
	p.mu.Unlock()

	err := p.drzewo.Ubij()
	// Uchwyt zadania i uchwyt procesu oddajemy zaraz po ubiciu. Bez tego każde
	// zamknięte okno zostawiałoby w rdzeniu uchwyt systemowy aż do końca pracy
	// procesu rdzenia.
	p.drzewo.Zwolnij()
	return err
}

// dogladaj czeka na faktyczne zakończenie procesu okna i domyka jego cykl
// życia: oznacza proces jako zakończony i oddaje uchwyty systemowe.
//
// Bez tego doglądu `zakonczony` miałby jednego pisarza — Ubij — a do Ubij
// prowadzą wyłącznie dwie drogi: zatrzymanie okna i przejęcie procesu następnej
// tury. Proces, który kończy się SAM, nie wyzwala żadnej z nich. Okno meldowało
// więc „Running" (stanProcesu → ProgressStatusRunning, Biegnace) od samoistnego
// wyjścia aż do najbliższej tury albo zamknięcia okna, a przez ten sam czas
// wisiał uchwyt Job Object i uchwyt os.Process.
//
// BRAMKA. Zwolnij ma teraz dwóch wołających — Ubij i ten dogląd — a
// drzewoProcesow.zwolnij NIE jest współbieżnie idempotentne: dwa CloseHandle na
// tym samym uchwycie zamykają uchwyt, który Windows zdążył już nadać ponownie
// czemu innemu. Bezpieczeństwo stoi WYŁĄCZNIE na odczycie-i-zapisie
// `zakonczony` pod tym samym p.mu, co w Ubij: kto zastanie false, ten jeden
// przechodzi dalej. Tego NIE WOLNO uprościć do `if p.Zyje() { ... }` — to byłoby
// check-then-act i obaj wołający mogliby wejść.
func (p *Proces) dogladaj() {
	p.drzewo.Czekaj()

	p.mu.Lock()
	juz := p.zakonczony
	p.zakonczony = true
	p.mu.Unlock()
	if juz {
		// Ubicie przeszło bramkę pierwsze i uchwyty już oddało. Tak wychodzi
		// obserwator poprzedniej tury, obudzony przez Przejmij.
		return
	}
	p.drzewo.Zwolnij()
}

// Zyje mówi, czy proces okna nadal pracuje.
func (p *Proces) Zyje() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return !p.zakonczony
}
