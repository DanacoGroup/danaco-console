// Odpowiedzialność pliku: ŻYWOTNOŚĆ REALNA — połączenie stanu podagenta
// z rejestrem procesów sesji (`session.RejestrProcesow`).
//
// CO REJESTR POTRAFI. `RejestrProcesow.Przejmij(idOkna, pid)` obejmuje proces
// tury uchwytem systemowym (grupa procesów na Uniksie, Job Object na Windows),
// a dogląd (`proces.dogladaj` → `ubicie_unix.go: syscall.Kill(pid, 0)` co
// 200 ms) utrzymuje odpowiedź `Zyje()` zgodną z prawdą systemu — nie z polem
// w pamięci. To jest żywotność realna i ten plik z niej wyłącznie CZYTA.
//
// CZEGO DROGA DZIŚ NIE MA:
//
//  1. Tura pozycji kolejki (a praca podagenta jest pozycją kolejki) jedzie
//     `kolejka_wykonawca.go` z PUSTYMI zasięgami:
//     `adapter_kolejki.rozwiazKanalPozycji` zwraca `models.Zasiegi{}`,
//     więc haczyk `przejmowanieProcesow.Haczyk("")` jest bezczynny i proces
//     wykonujący pozycję NIGDY nie trafia do rejestru.
//  2. Rejestr kluczuje procesy IDENTYFIKATOREM OKNA i trzyma jeden wpis na
//     okno; `Przejmij` UBIJA wpis poprzedni. Podagentów bywa piętnastu pod
//     jednym oknem i pracują RÓWNOLEGLE z turą własnego okna wykonawcy —
//     zarejestrowanie ich procesów pod oknem wykonawcy ubijałoby nawzajem
//     turę orkiestratora i tury podagentów. Brakuje klucza drobniejszego niż
//     okno (pozycja kolejki / podagent) i tego ten pakiet NIE obchodzi bokiem,
//     bo drugi rejestr procesów byłby drugą prawdą o procesach.
//
// CO WOBEC TEGO JEST POŁĄCZONE. Żywy jest mierzalnie proces ORKIESTRATORA —
// okna wykonawcy, które podagentów powołało: jego turę startuje `message.send`,
// zasięg okna jest wtedy wypełniony i `Przejmij` wpisuje proces do rejestru.
// Ocena niżej mówi więc prawdę o oknie prowadzącym podagentów, a o procesie
// samej pozycji mówi `BezWpisu` — i to zdanie jest prawdziwe, nie zastępcze.
//
// ŻYWOTNOŚĆ PO AWARII I PO RESTARCIE — druga połowa tego pliku. Dogląd wyżej
// mówi o procesie rdzenia, KTÓRY STOI. Gdy rdzeń padnie, nie mówi nic i nie ma
// komu mówić — dlatego pytanie „co się dzieje z podagentem po awarii" ma
// odpowiedź w bazie (`store/migracja_100_zywotnosc_podagentow.sql`), nie
// w rejestrze procesów. Znacznik
// uruchomienia i sprzątanie sierot stoją niżej.
package podagenci

import (
	"context"
	"fmt"
	"os"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
)

// stanProcesu jest trójstanem żywotności — trzeci stan istnieje NAPRAWDĘ:
// okno bez wpisu w rejestrze to nie to samo co okno z procesem zakończonym,
// a zlanie obu w „nie żyje" byłoby orzeczeniem z ciszy.
type stanProcesu int

const (
	// procesBezWpisu — rejestr nie ma wpisu dla tego okna. Dla procesu pozycji
	// podagenta jest to dziś stan JEDYNY (patrz nagłówek).
	procesBezWpisu stanProcesu = iota
	// procesZywy — proces okna jest w rejestrze i dogląd potwierdza życie.
	procesZywy
	// procesZakonczony — proces okna jest w rejestrze, ale już nie pracuje.
	procesZakonczony
)

// Opis oddaje stan zdaniem do dziennika i meldunku — po polsku, wprost.
func (s stanProcesu) Opis() string {
	switch s {
	case procesZywy:
		return "proces żywy (dogląd rejestru procesów potwierdza)"
	case procesZakonczony:
		return "proces zakończony (wpis w rejestrze bez życia)"
	default:
		return "bez wpisu w rejestrze procesów"
	}
}

// OcenaProcesu odpowiada na pytanie o żywotność procesu wskazanego okna.
// Funkcja, nie interfejs: rejestr procesów oddaje typ konkretny
// (`*session.Proces`), a wąska funkcja pozwala testom podstawić ocenę bez
// budowania rejestru.
type OcenaProcesu func(idOkna string) stanProcesu

// ZSesji wiąże ocenę z żywym rejestrem procesów nadzorcy sesji. Rejestr pusty
// albo niewpięty daje ocenę mówiącą zawsze `procesBezWpisu` — brak nadzorcy
// nie wymyśla żywotności i nie zatrzymuje niczego.
func ZSesji(procesy *session.RejestrProcesow) OcenaProcesu {
	return func(idOkna string) stanProcesu {
		if procesy == nil || idOkna == "" {
			return procesBezWpisu
		}
		proces, jest := procesy.Proces(idOkna)
		if !jest || proces == nil {
			return procesBezWpisu
		}
		if proces.Zyje() {
			return procesZywy
		}
		return procesZakonczony
	}
}

// wyjasnienieOsierocenia trafia w pole `wynik` sieroty, ale WYŁĄCZNIE gdy jest
// ono puste (COALESCE w zapytaniu): praca oddana przed awarią jest ważniejsza
// niż wyjaśnienie, dlaczego się urwała. Zdanie jest po polsku, bo czyta je
// Operator w panelu zadań w tle, a nie maszyna.
const wyjasnienieOsierocenia = "Praca przerwana zatrzymaniem rdzenia — " +
	"proces wykonujący zadanie nie istnieje po restarcie. Powołaj podagenta na nowo."

// ZnacznikUruchomienia nadaje znacznik bieżącemu uruchomieniu rdzenia.
//
// SKŁADA SIĘ Z DWÓCH RZECZY, BO ŻADNA SAMA NIE WYSTARCZA: numer procesu jest
// w systemie powtarzalny (po restarcie maszyny ten sam PID wraca), a czas sam
// nie odróżnia dwóch rdzeni wstałych w tej samej milisekundzie. Razem są
// jednoznaczne w praktyce, a jednoznaczności absolutnej ten znacznik nie
// potrzebuje: rozstrzyga wyłącznie pytanie „czy to nadal ja".
//
// Znacznik zakłada się RAZ na proces i podaje dalej wartością — losowania po
// drodze nie ma, więc nikt nie osieroci sam siebie.
func ZnacznikUruchomienia() string {
	return fmt.Sprintf("rdzen-%d-%s", os.Getpid(), time.Now().UTC().Format("20060102T150405.000Z"))
}

// PosprzatajPoRestarcie zamyka podagentów porzuconych przez uruchomienia
// wcześniejsze i oddaje wykaz zamkniętych — do meldunku w dzienniku.
//
// WOŁA SIĘ RAZ, PRZY STARCIE, PRZED PIERWSZYM POWOŁANIEM. Wywołanie późniejsze
// zamknęłoby pracę powołaną przez ten sam rdzeń, gdyby jej oznaczenie
// prowadzenia jeszcze nie doszło.
//
// TRWAŁOŚĆ PUSTA ZNOSI SIĘ SAMA: rdzeń bez repozytorium podagentów
// startuje, a nie odmawia startu — po prostu nie ma czego sprzątać.
func PosprzatajPoRestarcie(ctx context.Context, trwalosc dane.ZywotnoscPodagentow,
	uruchomienie string) ([]dane.Podagent, error) {

	if trwalosc == nil || uruchomienie == "" {
		return nil, nil
	}
	osieroceni, err := trwalosc.ZamknijOsierocone(ctx, uruchomienie, wyjasnienieOsierocenia)
	if err != nil {
		return nil, fmt.Errorf("podagenci: sprzątanie po restarcie nie doszło do skutku: %w", err)
	}
	return osieroceni, nil
}

// ZdanieOSprzataniu składa meldunek startowy. Zero sierot też jest meldunkiem —
// cisza po sprzątaniu nie odróżnia „nie było czego" od „nie zawołano".
func ZdanieOSprzataniu(osieroceni []dane.Podagent) string {
	if len(osieroceni) == 0 {
		return "sprzątanie po restarcie: podagentów osieroconych brak"
	}
	return fmt.Sprintf("sprzątanie po restarcie: zamknięto %d podagentów osieroconych "+
		"(praca przerwana zatrzymaniem poprzedniego rdzenia)", len(osieroceni))
}
