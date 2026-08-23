// Odpowiedzialność pliku: złożenie portu podagentów — zbudowanie adaptera
// z kompletu wiązań i dołożenie każdej zależności osobno.
//
// Podział wobec `adapter_modul_orkiestracja.go` idzie wzdłuż odpowiedzialności:
// tam mieszka powołanie podagentów (`subagent.spawn`) wraz z pracą w tle,
// tutaj wyłącznie budowanie portu i jego wiązania.
//
// Wszystkie dokładki znoszą się same przy pustej zależności.
package core

import (
	"context"
	"log"
	"sync"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/podagenci"
)

// adapterPodagentow wypełnia port Podagenci.
//
// Zależności są cztery i każda ma powód: repozytorium podagentów (trwałość),
// repozytorium okien (kontrakt niesie identyfikator zewnętrzny okna, schemat
// klucz wiersza), repozytorium biegów (podagent powołany w biegu ma do niego
// należeć) i adapter kolejek (jedyny wykonawca pracy).
type adapterPodagentow struct {
	repozytorium dane.RepozytoriumPodagentow
	okna         dane.RepozytoriumOkien
	biegi        dane.RepozytoriumBiegow
	kolejki      *adapterKolejek
	// zycie jest kontekstem rdzenia, nie kontekstem żądania. Praca podagenta
	// przeżywa odpowiedź na `subagent.spawn`, więc kontekst żądania zamknąłby ją
	// natychmiast po odesłaniu wyniku.
	zycie    context.Context
	dziennik *log.Logger
	// nadzor ocenia żywotność procesu okna z rejestru procesów sesji. Pusty znosi
	// się sam: powołanie idzie wtedy bez zdania o procesie orkiestratora, a nie
	// z żywotnością zgadywaną.
	nadzor podagenci.OcenaProcesu
	// uruchomienie to znacznik tego uruchomienia rdzenia. Puste znosi się samo:
	// powołanie idzie bez oznaczenia prowadzenia, a sprzątanie po restarcie nie
	// rusza niczego.
	uruchomienie string
	// drogaNarzedzia pilnuje, żeby meldunek o drodze narzędzia modelu padł raz
	// na proces, nie raz na powołanie — powołań bywa kilkanaście na turę.
	drogaNarzedzia sync.Once
	// prace trzyma odwołania pracy podagentów trwających, kluczowane kodem
	// podagenta, nie oknem. Bez tego wykazu `subagent.stop` nie miałby czego
	// zatrzymać i przepisywałby wyłącznie wiersz
	// (adapter_modul_orkiestracja_zatrzymanie.go).
	muPrace sync.Mutex
	prace   map[string]context.CancelFunc
	// rozgloszenie niesie `subagent.changed` do paneli. Puste znosi się samo:
	// podagenci powstają i pracują tak samo, gdy nikt zdarzeń nie słucha.
	// Treść rozgłoszeń — `adapter_modul_orkiestracja_rozgloszenie.go`.
	rozgloszenie *emiter
	// straz czyta zakres eksperta okna wykonawcy przed powołaniem podagentów.
	// Pusta znaczy „nie wpięto" — obowiązuje wtedy sama granica platformy.
	straz StrazEksperta
}

// Asercja kompilatora: adapter wypełnia port. Bez niej rozjazd podpisu wyszedłby
// dopiero przy montażu.
var _ Podagenci = (*adapterPodagentow)(nil)

// nowyAdapterPodagentow wiąże port z repozytorium podagentów i repozytorium
// okien. Obu wymaga sam przekład kontrakt↔schemat, więc idą argumentami, a nie
// dokładkami.
func nowyAdapterPodagentow(repozytorium dane.RepozytoriumPodagentow,
	okna dane.RepozytoriumOkien) *adapterPodagentow {

	return &adapterPodagentow{repozytorium: repozytorium, okna: okna}
}

// zlozPodagentow wypełnia port Podagenci gotowymi bytami składu portów.
//
// Stoi tu, a nie w montaz_porty.go, z dwóch powodów: montaż przekłada byty na
// porty jednym wierszem na port, a wiązań podagentów jest pięć; dzięki temu
// montaż nie musi też znać pakietu `podagenci` — wiedzę o ocenie żywotności
// trzyma adapter, który jako jedyny jej używa.
//
// Wszystkie dokładki znoszą się same przy pustej zależności:
// bez silnika kolejek podagent zostaje `pending`, bez biegów powstaje poza
// biegiem, bez nadzoru idzie bez zdania o procesie orkiestratora.
func zlozPodagentow(s skladPortow) Podagenci {
	return nowyAdapterPodagentow(s.repozytoria.Podagenci(), s.repozytoria.Okna).
		zBiegami(s.repozytoria.BiegiOrkiestracji()).
		zWykonaniem(s.kolejki, s.zycie).
		zNadzorem(podagenci.ZSesji(s.nadzorca.Procesy())).
		ZDziennikiem(s.montaz.Dziennik).
		ZRozgloszeniem(s.szynaZdarzen).
		ZeStrazaEksperta(s.strazEkspertow).
		zZywotnoscia(s.zycie)
}

// ZRozgloszeniem dokłada nadajnik zdarzeń `subagent.changed`.
//
// Nadajnik stoi przy adapterze, a nie przy uchwycie komendy, bo podagent
// zmienia stan głównie poza żądaniem: powołanie oddaje go jako `pending`,
// a wejście w `running`, zakończenie i niepowodzenie dzieją się w pracy
// puszczonej w tle (`adapter_modul_orkiestracja.go`, `puscWTle`). Uchwyt
// komendy widziałby wyłącznie pierwszy z tych czterech momentów.
func (a *adapterPodagentow) ZRozgloszeniem(nadajnik Nadajnik) *adapterPodagentow {
	a.rozgloszenie = nowyEmiter(nadajnik)
	return a
}

// zBiegami dokłada repozytorium biegów orkiestracji. Bez niego podagenci
// powstają poza biegiem — tak samo jak podagent powołany w zwykłej rozmowie,
// na co schemat pozwala.
func (a *adapterPodagentow) zBiegami(biegi dane.RepozytoriumBiegow) *adapterPodagentow {
	a.biegi = biegi
	return a
}

// zWykonaniem wpina silnik kolejek i kontekst życia rdzenia — jedyną drogę,
// którą praca podagenta naprawdę się wykonuje. Bez niego podagent zostaje
// w stanie `pending` i to jest stan prawdziwy, nie udawane wykonanie.
//
// Tu domyka się powrót wyniku do rodzica: adapter wpina się jako ujście wyniku
// silnika, więc zebrana treść tury pozycji trafia na wiersz podagenta wskazany
// tą pozycją i `subagent.result.collect` oddaje treść, a nie puste pole.
// Wpięcie idzie stąd, bo tylko ten adapter wie, że pozycja miewa podagenta —
// silnik zna wyłącznie pozycje.
func (a *adapterPodagentow) zWykonaniem(kolejki *adapterKolejek, zycie context.Context) *adapterPodagentow {
	a.kolejki, a.zycie = kolejki, zycie
	if kolejki != nil {
		kolejki.UstawUjscieWyniku(a.przyjmijWynikPozycji)
	}
	return a
}

// zNadzorem wpina ocenę żywotności procesów z rejestru procesów sesji. Ocena
// mówi prawdę o procesie okna — dla podagenta jest to proces jego
// orkiestratora; proces samej pozycji nie ma drogi do rejestru
// (`podagenci/zywotnosc.go`).
func (a *adapterPodagentow) zNadzorem(nadzor podagenci.OcenaProcesu) *adapterPodagentow {
	a.nadzor = nadzor
	return a
}

// zZywotnoscia nadaje znacznik uruchomienia i zamyka podagentów porzuconych
// przez rdzeń poprzedni. Woła się raz, przy montażu, zanim
// ktokolwiek powoła podagenta: wywołanie późniejsze zamknęłoby pracę powołaną
// przez ten sam rdzeń. Sierota to wiersz w stanie `pending`/`running` prowadzony
// przez uruchomienie inne niż bieżące — jego proces zginął razem z rdzeniem,
// więc stan `running` po restarcie kłamie.
func (a *adapterPodagentow) zZywotnoscia(ctx context.Context) *adapterPodagentow {
	a.uruchomienie = podagenci.ZnacznikUruchomienia()
	osieroceni, err := podagenci.PosprzatajPoRestarcie(ctx, a.repozytorium, a.uruchomienie)
	if err != nil {
		a.zapisz("podagenci: %v", err)
		return a
	}
	a.zapisz("podagenci: %s", podagenci.ZdanieOSprzataniu(osieroceni))
	return a
}

// kodyPodagentow wyjmuje identyfikatory zewnętrzne powołanych — tak pyta
// oznaczenie prowadzenia, które zna podagenta po kodzie, nie po kluczu wiersza.
func kodyPodagentow(powolani []dane.Podagent) []string {
	kody := make([]string, 0, len(powolani))
	for _, podagent := range powolani {
		kody = append(kody, podagent.Kod)
	}
	return kody
}

// przyjmijWynikPozycji jest ujściem wyniku silnika kolejek: utrwala zebraną
// treść tury na wierszu podagenta związanym z pozycją. Kontekst bierze z życia
// rdzenia, nie z tury — zapis wyniku ma przeżyć zamknięcie żądania, które turę
// wywołało. Brak trafienia to pozycja spoza podagentów (jeden silnik);
// błąd zapisu idzie do dziennika, bo cichej utraty wyniku nikt by nie zobaczył.
func (a *adapterPodagentow) przyjmijWynikPozycji(ctx context.Context, pozycjaID int64, tresc string) {
	if a.zycie != nil {
		ctx = a.zycie
	}
	if err := a.repozytorium.ZapiszWynikPozycji(ctx, pozycjaID, tresc); err != nil {
		a.zapisz("podagenci: wynik pozycji %d nie został utrwalony: %v", pozycjaID, err)
	}
}

// ZDziennikiem dokłada dziennik rdzenia. Znosi dziennik pusty sam.
func (a *adapterPodagentow) ZDziennikiem(dziennik *log.Logger) *adapterPodagentow {
	a.dziennik = dziennik
	return a
}

// zapisz nanosi wiersz do dziennika rdzenia, znosząc dziennik pusty.
func (a *adapterPodagentow) zapisz(wzorzec string, argumenty ...any) {
	if a == nil || a.dziennik == nil {
		return
	}
	a.dziennik.Printf(wzorzec, argumenty...)
}

// ZeStrazaEksperta wpina straż zakresu eksperta. Bez niej powołanie idzie samą
// granicą platformy — stanem wyjściowym jest pełny dostęp.
func (a *adapterPodagentow) ZeStrazaEksperta(straz StrazEksperta) *adapterPodagentow {
	a.straz = straz
	return a
}
