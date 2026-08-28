// Odpowiedzialność pliku: złożenie portu podagentów — zbudowanie adaptera z kompletu wiązań, każdej zależności dołożonej osobno i znoszącej się przy pustej wartości.
package core

import (
	"context"
	"log"
	"sync"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/podagenci"
)

// adapterPodagentow wypełnia port Podagenci na czterech zależnościach: repozytorium podagentów, repozytorium okien, repozytorium biegów oraz adapter kolejek jako jedyny wykonawca pracy.
type adapterPodagentow struct {
	repozytorium dane.RepozytoriumPodagentow
	okna         dane.RepozytoriumOkien
	biegi        dane.RepozytoriumBiegow
	kolejki      *adapterKolejek
	// zycie jest kontekstem rdzenia: praca podagenta przeżywa odpowiedź na subagent.spawn.
	zycie    context.Context
	dziennik *log.Logger
	// nadzor ocenia żywotność procesu okna; pusty znosi się sam, bez zdania o procesie orkiestratora.
	nadzor podagenci.OcenaProcesu
	// uruchomienie znaczy to uruchomienie rdzenia; puste znosi się samo, bez oznaczenia prowadzenia.
	uruchomienie string
	// drogaNarzedzia pilnuje, żeby meldunek o narzędziu padł raz na proces, nie raz na powołanie.
	drogaNarzedzia sync.Once
	// prace trzyma odwołania pracy podagentów trwających, kluczowane kodem podagenta, nie oknem.
	muPrace sync.Mutex
	prace   map[string]context.CancelFunc
	// rozgloszenie niesie subagent.changed do paneli; puste znosi się samo, bez zdarzeń.
	rozgloszenie *emiter
	// straz czyta zakres eksperta okna wykonawcy; pusta znaczy nie wpięto, obowiązuje granica.
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

// zlozPodagentow wypełnia port Podagenci gotowymi bytami składu portów, dokładając pięć wiązań osobno, każde znoszące się przy pustej zależności.
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

// ZRozgloszeniem dokłada nadajnik zdarzeń subagent.changed. Nadajnik stoi przy adapterze, bo podagent zmienia stan głównie poza żądaniem, w pracy puszczonej w tle.
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

// zWykonaniem wpina silnik kolejek i kontekst życia rdzenia, jedyną drogę wykonania podagenta, oraz ustawia adapter jako ujście wyniku silnika, przez które wraca treść tury.
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

// zZywotnoscia nadaje znacznik uruchomienia i zamyka podagentów porzuconych przez rdzeń poprzedni; woła się raz, przy montażu, zanim ktokolwiek powoła podagenta.
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

// przyjmijWynikPozycji jest ujściem wyniku silnika kolejek: utrwala zebraną treść tury na wierszu podagenta związanym z pozycją, kontekstem życia rdzenia, nie tury.
func (a *adapterPodagentow) przyjmijWynikPozycji(ctx context.Context, pozycjaID int64, tresc string) {
	if a.zycie != nil {
		ctx = a.zycie
	}
	if err := a.repozytorium.ZapiszWynikPozycji(ctx, pozycjaID, tresc); err != nil {
		a.zapisz("podagenci: wynik pozycji %d nie został utrwalony: %v", pozycjaID, err)
	}
}

// ZDziennikiem dokłada dziennik rdzenia do adaptera podagentów, znosząc się sam przy dzienniku pustym.
func (a *adapterPodagentow) ZDziennikiem(dziennik *log.Logger) *adapterPodagentow {
	a.dziennik = dziennik
	return a
}

// zapisz nanosi wiersz do dziennika rdzenia adaptera podagentów, nie robiąc niczego przy dzienniku pustym.
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
