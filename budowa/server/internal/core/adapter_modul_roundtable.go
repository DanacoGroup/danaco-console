// Odpowiedzialność pliku: moduł Roundtable — skład debaty (okno Model Panels)
// i wiązanie adaptera z rejestrem kanałów, dla żądań obszaru `roundtable.*`.
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

// Przedrostki identyfikatorów bytów modułu, nadawanych przez rdzeń przy
// zakładaniu nowego wiersza debaty albo tury.
const (
	przedrostekUczestnika = "uczest-"
	przedrostekTury       = "tura-"
	przedrostekWypowiedzi = "wypow-"
	przedrostekStanowiska = "stanow-"
	// kodModeratora znakuje wypowiedź moderatora w zapisie tury, stałym kodem
	// rozpoznawalnym w transkrypcie debaty, bo moderator nie jest uczestnikiem.
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
	// zycie jest kontekstem rdzenia: rozłączenie klienta nie przerywa tury.
	zycie context.Context

	// katalogArtefaktow trzyma bajty wydanych transkryptów, grafów i nagrań.
	katalogArtefaktow string
	// uruchamiacz jest jedyną drogą startu programu Pandoc i silnika mowy.
	uruchamiacz session.Uruchamiacz
	// rozstrzygacz składa zasady izolacji egzekwowane przy uruchomieniu.
	rozstrzygacz *konfig.Rozstrzygacz

	mu       sync.Mutex
	biegnace map[string]*biegDebaty
}

// biegDebaty opisuje zajęcie okna pod turę. Licznik tur jest konieczny, bo
// tura przerwana kończy się później, niż następna się zaczyna, a wpis wolno
// wykreślić dopiero tej, która schodzi z okna ostatnia.
type biegDebaty struct {
	// anuluj przerywa turę, która okno zajmuje w tej chwili.
	anuluj context.CancelFunc
	// zajete mówi, czy okno prowadzi turę przyjmującą wypowiedzi. Tura
	// przerwana zwalnia okno od razu, choć jej głosy milkną dopiero po chwili.
	zajete bool
	// tury liczy tury, które okna jeszcze nie opuściły — wraz z przerwanymi.
	tury int
}

// nowyAdapterDebaty wiąże port Debata z repozytorium modułu, jedyną
// zależnością wymaganą konstruktorem.
func nowyAdapterDebaty(zycie context.Context, repozytorium dane.RepozytoriumRoundtable) *adapterDebaty {
	return &adapterDebaty{
		repozytorium: repozytorium, zycie: zycie,
		biegnace: map[string]*biegDebaty{},
	}
}

// ZKanalami podpina rejestr kanałów modelu. Bez niego moduł prowadzi skład,
// tury i stanowisko, lecz uruchomienie tury odmawia wprost — debata bez
// wykonawcy nie ma prawa udawać, że uczestnicy odpowiedzieli.
func (a *adapterDebaty) ZKanalami(kanaly *models.Rejestr) *adapterDebaty {
	a.kanaly = kanaly
	return a
}

// ZArsenalem podpina uruchamiacz procesów i rozstrzygacz zasięgu, bez których
// nie ruszy zamiana transkryptu na dokument biurowy ani synteza mowy.
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

// PodepnijRozgloszenie wypełnia port nadajnika: adapter zapamiętuje drogę
// do rozgłoszenia wypowiedzi debaty w czasie rzeczywistym do zdarzenia.
func (a *adapterDebaty) PodepnijRozgloszenie(
	rozglos func(shared.ChangeKind, shared.RoundtableTurn, *shared.RoundtableStatement)) {

	a.zmiana = rozglos
}

// DodajModel dopisuje uczestnika debaty, sprawdzając od razu kanał w rejestrze,
// żeby odmówić w chwili dodania, a nie dopiero po zadaniu pytania.
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
// Sprawdzenie i zajęcie idą pod jednym zamkiem, żeby dwa otwarcia naraz nie
// zobaczyły obie okna wolnego. Okno zajęte nie jest przerywane po cichu.
func (a *adapterDebaty) zajmijBieg(okno string, anuluj context.CancelFunc) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	bieg, jest := a.biegnace[okno]
	if !jest {
		a.biegnace[okno] = &biegDebaty{anuluj: anuluj, zajete: true, tury: 1}
		return true
	}
	if bieg.zajete {
		return false
	}
	bieg.anuluj, bieg.zajete, bieg.tury = anuluj, true, bieg.tury+1
	return true
}

// przejmijBieg zajmuje okno pod turę moderatora, przerywając turę biegnącą —
// jawnie, w odróżnieniu od `zajmijBieg`, bo ukierunkowanie dyskusji jest
// aktem przerwania, a nie skutkiem ubocznym wysłania czegoś innego.
func (a *adapterDebaty) przejmijBieg(okno string, anuluj context.CancelFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()
	bieg, jest := a.biegnace[okno]
	if !jest {
		a.biegnace[okno] = &biegDebaty{anuluj: anuluj, zajete: true, tury: 1}
		return
	}
	if bieg.zajete {
		bieg.anuluj()
	}
	bieg.anuluj, bieg.zajete, bieg.tury = anuluj, true, bieg.tury+1
}

// zwolnijBieg oddaje okno zajęte pod turę, która ostatecznie nie ruszyła —
// każde wyjście błędem między zajęciem a założeniem tury musi okno oddać,
// inaczej Debate Panel zostawałby zablokowany turą, której nigdy nie było.
func (a *adapterDebaty) zwolnijBieg(okno string, anuluj context.CancelFunc) {
	a.zejdzZOkna(okno)
	anuluj()
}

// zapomnijBieg zdejmuje turę z rejestru biegów po jej zakończeniu, zwalniając
// okno pod kolejne otwarcie.
func (a *adapterDebaty) zapomnijBieg(okno string) {
	a.zejdzZOkna(okno)
}

// zejdzZOkna odlicza turę, która zeszła z okna, i wykreśla wpis dopiero wtedy,
// gdy okna nie prowadzi już żadna. Tura kończąca się nie zna swojego wpisu —
// wpis zastany może należeć do tury, która ją przejęła — więc wykreślenie
// bezwarunkowe zdejmowałoby z okna turę cudzą i wpuszczało drugą obok niej.
func (a *adapterDebaty) zejdzZOkna(okno string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	bieg, jest := a.biegnace[okno]
	if !jest {
		return
	}
	bieg.tury--
	if bieg.tury <= 0 {
		delete(a.biegnace, okno)
	}
}

// PrzerwijBieg przerywa turę okna i oddaje prawdę, gdy było co przerywać.
// Wywołuje to zamknięcie tury przez moderatora: uczestnicy, którzy jeszcze
// mówią, mają przestać, bo tura już nie przyjmuje wypowiedzi. Zatrzymanie jest
// dostępne zawsze. Wpisu nie wykreśla — zdejmuje go tura, która z okna schodzi,
// bo tura przerwana kończy się dopiero po tym wywołaniu.
func (a *adapterDebaty) PrzerwijBieg(okno string) bool {
	a.mu.Lock()
	var anuluj context.CancelFunc
	if bieg, jest := a.biegnace[okno]; jest && bieg.zajete {
		anuluj, bieg.zajete = bieg.anuluj, false
	}
	a.mu.Unlock()
	if anuluj == nil {
		return false
	}
	anuluj()
	return true
}
