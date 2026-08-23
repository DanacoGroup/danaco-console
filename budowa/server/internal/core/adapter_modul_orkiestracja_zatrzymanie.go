// Odpowiedzialność pliku: zatrzymanie jednego podagenta — komenda
// `subagent.stop`.
//
// `message.stop` zatrzymuje całe okno: odwołuje turę orkiestratora i wstrzymuje
// pętlę naprawczą. Podagentów pod jednym oknem bywa kilkunastu i pracują
// równolegle, więc zatrzymanie pojedynczego wymaga osobnej komendy.
//
// Zatrzymanie jest realne, nie tabelaryczne: przepisanie wiersza na `stopped`
// bez ruszenia procesu dawałoby stan niezgodny z faktami. Zatrzymanie idzie
// odwołaniem kontekstu pracy — `injection/rozruch.go` startuje program przez
// `exec.CommandContext`, więc odwołanie kontekstu kończy proces modelu, a nie
// tylko przestaje go obserwować. Wiersz przechodzi na `stopped` po odwołaniu,
// jako zapis tego, co zaszło.
//
// Podagent już zakończony nie jest błędem — wraca w wykazie `notRunning`, nie
// w odmowie. Rozdział na dwa wykazy mówi, których zatrzymało to wywołanie,
// a którzy skończyli wcześniej sami.
//
// Wskazanie, któremu nie odpowiada żaden wiersz, jest pomyłką co do bytu —
// wzorem `objeciZbieraniem` (adapter_modul_orkiestracja_wykaz.go). Pusty wynik
// czytałoby się jako „nikt nie pracował", a to co innego niż „takich podagentów
// nie ma".
package core

import (
	"context"
	"sort"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// wyjasnienieZatrzymania trafia w pole `wynik` podagenta zatrzymanego, ale
// wyłącznie gdy jest ono puste (COALESCE w zapytaniu `ustawStanPodagenta`):
// praca oddana przed zatrzymaniem jest ważniejsza niż wyjaśnienie, dlaczego się
// urwała. Zdanie jest po polsku, bo czyta je Operator w panelu zadań w tle.
const wyjasnienieZatrzymania = "Praca zatrzymana przez Operatora — " +
	"proces wykonujący zadanie został przerwany. Powołaj podagenta na nowo, " +
	"jeśli zadanie ma być dokończone."

// zapamietajPrace zapisuje odwołanie pracy podagenta pod jego kodem.
//
// Kluczem jest kod podagenta, nie okno — inaczej niż w
// `session.RejestrProcesow`, który kluczuje oknem i trzyma jeden wpis na okno
// (`podagenci/zywotnosc.go`). Kilkunastu podagentów pod jednym oknem potrzebuje
// tyluż uchwytów; wpis pod oknem ubijałby je nawzajem razem z turą
// orkiestratora.
//
// Nie jest to drugi rejestr procesów: rejestr sesji mówi o procesach systemu,
// a ten wykaz trzyma odwołania kontekstów pracy — byt, którego tamten rejestr
// nie zna.
func (a *adapterPodagentow) zapamietajPrace(kod string, odwolaj context.CancelFunc) {
	if kod == "" || odwolaj == nil {
		return
	}
	a.muPrace.Lock()
	defer a.muPrace.Unlock()
	if a.prace == nil {
		a.prace = map[string]context.CancelFunc{}
	}
	a.prace[kod] = odwolaj
}

// zapomnijPrace zdejmuje odwołanie po zakończeniu pracy — wykaz ma mówić
// o pracy trwającej, a nie rosnąć z każdym powołaniem od startu rdzenia.
func (a *adapterPodagentow) zapomnijPrace(kod string) {
	a.muPrace.Lock()
	defer a.muPrace.Unlock()
	delete(a.prace, kod)
}

// odwolajPrace odwołuje pracę podagenta i mówi, czy było co odwoływać.
//
// Zdjęcie wpisu idzie pod tym samym zamkiem co odczyt: dwa zatrzymania nadane
// w tej samej chwili mają dać jedno `stopped` i jedno `notRunning`, a nie dwa
// razy to samo.
func (a *adapterPodagentow) odwolajPrace(kod string) bool {
	a.muPrace.Lock()
	odwolaj, biegnie := a.prace[kod]
	delete(a.prace, kod)
	a.muPrace.Unlock()
	if !biegnie || odwolaj == nil {
		return false
	}
	odwolaj()
	return true
}

// Zatrzymaj zatrzymuje wskazanych podagentów okna wykonawcy.
//
// Tury okna nie dotyka — na tym polega różnica wobec `message.stop`:
// orkiestrator pracuje dalej, kiedy jeden z jego podagentów zostaje odwołany.
func (a *adapterPodagentow) Zatrzymaj(ctx context.Context,
	z shared.SubagentStopRequest) (shared.SubagentStopResponse, error) {

	objeci, err := a.objeciZatrzymaniem(ctx, z)
	if err != nil {
		return shared.SubagentStopResponse{}, err
	}

	// Wykazy zakładamy puste, a nie zerowe: kontrakt niesie oba jako wymagane,
	// więc `null` w miejscu wykazu byłby dla klienta czym innym niż wykaz pusty.
	zatrzymani := []string{}
	nieczynni := []string{}
	for _, podagent := range objeci {
		if !a.odwolajPrace(podagent.Kod) {
			nieczynni = append(nieczynni, podagent.Kod)
			continue
		}
		zatrzymani = append(zatrzymani, podagent.Kod)
		wyjasnienie := wyjasnienieZatrzymania
		if err := a.ustawStanPodagenta(ctx, podagent.Kod,
			dane.StanPodagentaZatrzymany, &wyjasnienie); err != nil {
			// Proces już stanął, więc nieudany zapis stanu nie zamienia
			// zatrzymania w odmowę. Ślad zostaje w dzienniku, a wiersz dogoni
			// najbliższe przepisanie stanu pozycji.
			a.zapisz("subagent.stop: podagent %s zatrzymany, ale stanu nie zapisano: %v",
				podagent.Kod, err)
		}
	}
	sort.Strings(zatrzymani)
	sort.Strings(nieczynni)

	// Wykaz idzie z ponownego odczytu, już po zatrzymaniu — odpowiedź niesie
	// stan po zmianie, a nie ten sprzed niej.
	poZatrzymaniu, err := a.repozytorium.PodagenciPoKodach(ctx, kodyPodagentow(objeci))
	if err != nil {
		return shared.SubagentStopResponse{}, bladPodagentow(err)
	}
	return shared.SubagentStopResponse{
		Stopped:    zatrzymani,
		NotRunning: nieczynni,
		Subagents:  podagenciKontraktu(poZatrzymaniu),
	}, nil
}

// objeciZatrzymaniem dobiera podagentów objętych wywołaniem.
//
// Pusta lista wskazań znaczy komplet wskazanego okna — tak samo jak przy
// zbieraniu wyników. Różnica jest jedna i celowa: zatrzymanie bez okna i bez
// wskazania podagentów odmawia. Zbieranie wyników bez zawężenia jest pytaniem,
// zatrzymanie bez zawężenia przerwałoby całą pracę platformy jednym wywołaniem.
func (a *adapterPodagentow) objeciZatrzymaniem(ctx context.Context,
	z shared.SubagentStopRequest) ([]dane.Podagent, error) {

	kody := niepusteKody(z.SubagentIds)
	if len(kody) == 0 {
		okno := wartoscTekstu(z.WindowId)
		if okno == "" {
			return nil, bladWskazaniaPodagenta(
				"zatrzymanie bez wskazania okna wykonawcy i bez wskazania podagentów; " +
					"wskaż okno albo wymień podagentów, których praca ma stanąć")
		}
		wiersze, err := a.repozytorium.Podagenci(ctx, dane.FiltrPodagentow{OknoKod: okno})
		if err != nil {
			return nil, bladPodagentow(err)
		}
		return wiersze, nil
	}
	wiersze, err := a.repozytorium.PodagenciPoKodach(ctx, kody)
	if err != nil {
		return nil, bladPodagentow(err)
	}
	if len(wiersze) == 0 {
		return nil, bladNieznanychPodagentow(kody)
	}
	return wiersze, nil
}
