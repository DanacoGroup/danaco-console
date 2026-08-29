// Żywotność realna łączy stan podagenta z rejestrem procesów sesji; żywy jest
// mierzalnie proces okna wykonawcy, które podagentów powołało, nie proces
// samej pozycji.
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
	// procesBezWpisu — rejestr nie ma wpisu dla tego okna; dla procesu pozycji
	// podagenta jest to dziś stan jedyny.
	procesBezWpisu stanProcesu = iota
	// procesZywy — proces okna jest w rejestrze procesów sesji, a dogląd
	// rejestru potwierdza życie procesu.
	procesZywy
	// procesZakonczony — proces okna jest w rejestrze procesów sesji, ale on
	// już nie pracuje wedle doglądu.
	procesZakonczony
)

// Opis oddaje stan procesu zdaniem do dziennika i meldunku dla Operatora, po
// polsku, wprost, bez skrótów.
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

// wyjasnienieOsierocenia trafia w pole wynik sieroty, ale wyłącznie gdy jest
// ono puste, po polsku, bo czyta je Operator.
const wyjasnienieOsierocenia = "Praca przerwana zatrzymaniem serwera — " +
	"proces wykonujący zadanie nie istnieje po restarcie. Powołaj podagenta na nowo."

// ZnacznikUruchomienia nadaje znacznik bieżącemu uruchomieniu rdzenia,
// złożony z numeru procesu i chwili startu.
func ZnacznikUruchomienia() string {
	return fmt.Sprintf("rdzen-%d-%s", os.Getpid(), time.Now().UTC().Format("20060102T150405.000Z"))
}

// PosprzatajPoRestarcie zamyka podagentów porzuconych przez uruchomienia
// wcześniejsze i oddaje wykaz zamkniętych.
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
		"(praca przerwana zatrzymaniem poprzedniego serwera)", len(osieroceni))
}
