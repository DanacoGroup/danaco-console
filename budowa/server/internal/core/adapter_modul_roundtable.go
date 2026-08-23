// Odpowiedzialność pliku: moduł Roundtable — skład debaty (okno Model Panels)
// i wiązanie adaptera z rejestrem kanałów. Tura i jej wykonanie leżą
// w `adapter_modul_roundtable_tura.go`, czynności moderatora
// w `adapter_modul_roundtable_moderator.go`, stanowisko końcowe
// w `adapter_modul_roundtable_stanowisko.go`, przekład bytów
// w `adapter_modul_roundtable_przeklad.go`.
//
// To jedyny moduł, w którym jedno okno rozmawia z wieloma kanałami naraz.
// Rdzeń to wspiera bez zmian: `models.Rejestr.Wyslij` jest bezpieczny do
// równoległego wywołania (rejestr trzyma kanały pod RWMutex, a każdy adapter
// dostaje własne zapytanie i własne ujście), więc N uczestników to N gorutyn
// nad jednym rejestrem. Drugiego rejestru kanałów moduł nie zakłada.
//
// Ten sam kanał może wystąpić dwukrotnie pod odrębnymi tożsamościami. Tożsamość
// uczestnika jedzie do kanału warstwą `models.Nakladka.ProfilRoli` — tą samą,
// którą okno rozmowy niesie profil roli. Dzięki temu dwaj uczestnicy na kanale
// `claude-cli` różnią się promptem systemowym, a nie kodem kanału.
package core

import (
	"context"
	"strings"
	"sync"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów modułu.
const (
	przedrostekUczestnika = "uczest-"
	przedrostekTury       = "tura-"
	przedrostekWypowiedzi = "wypow-"
	przedrostekStanowiska = "stanow-"
	// kodModeratora znakuje wypowiedź moderatora w zapisie tury. Moderator nie
	// jest uczestnikiem — nie ma kanału ani persony — więc kod jest stały
	// i rozpoznawalny w transkrypcie.
	kodModeratora = "moderator"
)

// adapterDebaty wypełnia port Debata. Zależności są trzy: repozytorium modułu,
// rejestr kanałów modelu (jedyne źródło uczestników) i nadajnik zdarzeń, którym
// idą zarówno przyrosty debaty, jak i strumień wypowiedzi.
type adapterDebaty struct {
	repozytorium dane.RepozytoriumRoundtable
	kanaly       *models.Rejestr
	nadajnik     Nadajnik
	zmiana       func(shared.ChangeKind, shared.RoundtableTurn, *shared.RoundtableStatement)
	// zycie jest kontekstem rdzenia, nie połączenia: rozłączenie klienta nie
	// przerywa rozpoczętej tury.
	zycie context.Context

	// katalogArtefaktow trzyma bajty wydanych transkryptów, grafów i nagrań
	// (patrz `adapter_modul_roundtable_magazyn.go`).
	katalogArtefaktow string
	// uruchamiacz jest jedyną drogą startu programu serwerowego: Pandoc przy
	// dokumencie biurowym, silnik mowy przy odsłuchu. Bez niego obie czynności
	// odmawiają, nazywając brak, zamiast milczeć.
	uruchamiacz session.Uruchamiacz
	// rozstrzygacz składa zasady izolacji egzekwowane przy uruchomieniu — ten
	// sam, którym jadą Terminal, Developer i rodzina narzędzi mediów.
	rozstrzygacz *konfig.Rozstrzygacz

	mu       sync.Mutex
	biegnace map[string]context.CancelFunc
}

// nowyAdapterDebaty wiąże port z repozytorium modułu.
func nowyAdapterDebaty(zycie context.Context, repozytorium dane.RepozytoriumRoundtable) *adapterDebaty {
	return &adapterDebaty{
		repozytorium: repozytorium, zycie: zycie,
		biegnace: map[string]context.CancelFunc{},
	}
}

// ZKanalami podpina rejestr kanałów modelu. Bez niego moduł prowadzi skład,
// tury i stanowisko, lecz uruchomienie tury odmawia wprost — debata bez
// wykonawcy nie ma prawa udawać, że uczestnicy odpowiedzieli.
func (a *adapterDebaty) ZKanalami(kanaly *models.Rejestr) *adapterDebaty {
	a.kanaly = kanaly
	return a
}

// ZArsenalem podpina uruchamiacz procesów i rozstrzygacz zasięgu — dwa
// źródła, bez których nie ruszy ani zamiana transkryptu na dokument biurowy,
// ani synteza mowy. Bez nich reszta modułu działa bez zmian: debata, analiza
// i głosowanie nie wołają ani jednego programu.
func (a *adapterDebaty) ZArsenalem(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz) *adapterDebaty {

	a.uruchamiacz, a.rozstrzygacz = uruchamiacz, rozstrzygacz
	return a
}

// ZWyjsciem podpina nadajnik strumienia wypowiedzi. Bez niego debata biegnie
// dalej, a Model Panels dostają wypowiedź dopiero w całości.
func (a *adapterDebaty) ZWyjsciem(nadajnik Nadajnik) *adapterDebaty {
	a.nadajnik = nadajnik
	return a
}

// PodepnijRozgloszenie wypełnia port: adapter zapamiętuje drogę do zdarzenia.
func (a *adapterDebaty) PodepnijRozgloszenie(
	rozglos func(shared.ChangeKind, shared.RoundtableTurn, *shared.RoundtableStatement)) {

	a.zmiana = rozglos
}

// DodajModel dopisuje uczestnika debaty.
//
// Kanał sprawdza się w rejestrze. Uczestnik na kanale, którego nie ma, milczałby
// w każdej turze, a Operator dowiedziałby się o tym dopiero po zadaniu pytania.
// Odmowa w chwili dodania mówi mu to od razu i wskazuje przyczynę.
func (a *adapterDebaty) DodajModel(ctx context.Context,
	z shared.RoundtableModelAddRequest) (shared.RoundtableModelAddResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kanal := strings.TrimSpace(z.ChannelId)
	if okno == "" {
		return shared.RoundtableModelAddResponse{}, bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kanal == "" {
		return shared.RoundtableModelAddResponse{}, bladWskazaniaDebaty("uczestnik bez wskazania kanału modelu")
	}
	if a.kanaly == nil {
		return shared.RoundtableModelAddResponse{}, bladBrakuKanalow()
	}
	if _, jest := a.kanaly.Kanal(kanal); !jest {
		return shared.RoundtableModelAddResponse{}, bladNieznanegoKanalu(kanal)
	}

	skladu, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return shared.RoundtableModelAddResponse{}, bladDebaty(err)
	}
	uczestnik, err := a.repozytorium.ZapiszUczestnika(ctx, dane.UczestnikDebaty{
		Kod:             nowyIdentyfikator(przedrostekUczestnika),
		Okno:            okno,
		KanalModelu:     kanal,
		NazwaTozsamosci: wskaznikTekstu(strings.TrimSpace(wartoscTekstu(z.PersonaName))),
		PromptSystemowy: wskaznikTekstu(strings.TrimSpace(wartoscTekstu(z.SystemPrompt))),
		Kolejnosc:       len(skladu) + 1,
	})
	if err != nil {
		return shared.RoundtableModelAddResponse{}, bladDebaty(err)
	}
	return shared.RoundtableModelAddResponse{Participant: uczestnikKontraktu(uczestnik)}, nil
}

// sklad zwraca uczestników okna wraz z odmową, gdy debata nie ma jeszcze
// żadnego. Pytanie bez adresata nie jest turą, tylko pomyłką Operatora.
func (a *adapterDebaty) sklad(ctx context.Context, okno string) ([]dane.UczestnikDebaty, error) {
	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return nil, bladDebaty(err)
	}
	if len(uczestnicy) == 0 {
		return nil, bladWskazaniaDebaty(
			"debata okna " + okno + " nie ma ani jednego uczestnika — dodaj model w Model Panels")
	}
	return uczestnicy, nil
}

// mowiacy zawęża skład do uczestników niewyciszonych w turze. Wyciszenie jest
// czynnością moderatora, więc skład zostaje, a milczy tylko wskazany.
func mowiacy(uczestnicy []dane.UczestnikDebaty) []dane.UczestnikDebaty {
	czynni := make([]dane.UczestnikDebaty, 0, len(uczestnicy))
	for _, uczestnik := range uczestnicy {
		if !uczestnik.Wyciszony {
			czynni = append(czynni, uczestnik)
		}
	}
	return czynni
}

// zajmijBieg zajmuje okno pod nową turę i oddaje prawdę, gdy się to udało.
//
// Sprawdzenie i zajęcie idą pod jednym zamkiem — odczyt osobny od zapisu
// zostawiłby szczelinę, w której dwa otwarcia nadane w tej samej chwili obie
// zobaczyłyby okno wolne.
//
// Okno zajęte nie jest przerywane: ciche odwołanie tury biegnącej kasowałoby
// wypowiedzi uczestników w połowie zdania, bez odmowy i bez śladu w panelu.
// Przerwanie należy do moderatora i idzie jego komendą (zamknięcie tury,
// `PrzerwijBieg`), a nie przy okazji otwarcia następnej.
func (a *adapterDebaty) zajmijBieg(okno string, anuluj context.CancelFunc) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, biegnie := a.biegnace[okno]; biegnie {
		return false
	}
	a.biegnace[okno] = anuluj
	return true
}

// przejmijBieg zajmuje okno pod turę moderatora, przerywając turę biegnącą.
//
// Przerwanie jest tu jawne i na tym polega różnica wobec `zajmijBieg`.
// Ukierunkowanie dyskusji (`roundtable.moderator.direct`) jest aktem
// przerwania: moderator wchodzi uczestnikom w słowo, żeby zawrócić rozmowę,
// i po to tę komendę wywołuje. Nie jest to skutek uboczny wysłania czegoś
// innego, tylko treść samej komendy — dlatego przejęcie zostaje, a nie odmawia.
func (a *adapterDebaty) przejmijBieg(okno string, anuluj context.CancelFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if poprzedni, jest := a.biegnace[okno]; jest {
		poprzedni()
	}
	a.biegnace[okno] = anuluj
}

// zwolnijBieg oddaje okno zajęte pod turę, która ostatecznie nie ruszyła.
//
// Istnieje, bo zajęcie idzie przed założeniem tury (adapter_modul_roundtable_tura.go):
// każde wyjście błędem pomiędzy jednym a drugim musi okno oddać, inaczej Debate
// Panel zostawałby zablokowany turą, której nigdy nie było. Odwołanie zwalnia
// się przy okazji — kontekst bez odbiorcy nie ma po co żyć.
//
// Zdjęcie wpisu jest bezwarunkowe i takie być może: między zajęciem
// a zwolnieniem nie startuje żadna gorutyna, a wpis w rejestrze każe kolejnemu
// otwarciu odmówić, więc odwołanie zdejmowane tu jest zawsze tym samym, które
// zajęło okno. Drogą tą wolno wołać wyłącznie przed startem tury; po starcie
// okno zwalnia `zapomnijBieg` z `defer` w prowadzTure.
func (a *adapterDebaty) zwolnijBieg(okno string, anuluj context.CancelFunc) {
	a.mu.Lock()
	delete(a.biegnace, okno)
	a.mu.Unlock()
	anuluj()
}

// zapomnijBieg zdejmuje turę z rejestru biegów po jej zakończeniu.
func (a *adapterDebaty) zapomnijBieg(okno string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.biegnace, okno)
}

// PrzerwijBieg przerywa turę okna. Wywołuje to zamknięcie tury przez moderatora:
// uczestnicy, którzy jeszcze mówią, mają przestać, bo tura już nie przyjmuje
// wypowiedzi. Zatrzymanie jest dostępne zawsze.
func (a *adapterDebaty) PrzerwijBieg(okno string) bool {
	a.mu.Lock()
	anuluj, biegnie := a.biegnace[okno]
	delete(a.biegnace, okno)
	a.mu.Unlock()
	if biegnie {
		anuluj()
	}
	return biegnie
}
