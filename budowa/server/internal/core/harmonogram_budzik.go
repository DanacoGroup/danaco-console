// Plik wpina budzik harmonogramu — długożyjący komponent, który obserwuje wyliczone terminy
// automatyk i o wskazanej godzinie odpala ich uruchomienie tą samą drogą co Operator.
package core

import (
	"context"
	"log"
	"time"

	"danacoconsole/server/internal/dane"
)

// interwalBudzika wyznacza takt sprawdzania terminów. Minuta odpowiada
// rozdzielczości zapisu cron (pole minuty) — częstsze pytanie nie trafiłoby
// w żaden nowy termin, rzadsze spóźniałoby uruchomienia.
const interwalBudzika = time.Minute

// budzikHarmonogramu odpala automatyki o wyliczonej godzinie. Stoi na tym samym
// adapterze Automations, którym jedzie moduł — nie zna drugiego repozytorium ani
// drugiego silnika kolejek.
type budzikHarmonogramu struct {
	automatyki *adapterAutomatyk
	dziennik   *log.Logger
}

// nowyBudzikHarmonogramu wiąże budzik z adapterem Automations i dziennikiem
// rdzenia. Dziennik pusty wyłącza wyłącznie ślad, nie budzik.
func nowyBudzikHarmonogramu(automatyki *adapterAutomatyk, dziennik *log.Logger) *budzikHarmonogramu {
	return &budzikHarmonogramu{automatyki: automatyki, dziennik: dziennik}
}

// Uruchom wpina budzik w cykl życia rdzenia: zakłada wątek, który tyka aż do
// zamknięcia kontekstu życia. Budzik bez adaptera Automations nie ma czego
// odpalać i nie startuje wątku.
func (b *budzikHarmonogramu) Uruchom(zycie context.Context) {
	if b == nil || b.automatyki == nil || b.automatyki.repozytorium == nil {
		return
	}
	go b.petla(zycie)
}

// petla tyka co takt i wyzwala harmonogramy należne. Pierwszy przebieg idzie od
// razu po starcie, żeby terminy przeoczone w czasie postoju rdzenia odpaliły
// się bez czekania na pełny takt.
func (b *budzikHarmonogramu) petla(zycie context.Context) {
	zegar := time.NewTicker(interwalBudzika)
	defer zegar.Stop()
	b.wyzwolNalezne(zycie, time.Now().UTC())
	for {
		select {
		case <-zycie.Done():
			return
		case teraz := <-zegar.C:
			b.wyzwolNalezne(zycie, teraz.UTC())
		}
	}
}

// wyzwolNalezne odpala wszystkie harmonogramy czynne, których termin już minął w tym takcie zegara budzika.
func (b *budzikHarmonogramu) wyzwolNalezne(ctx context.Context, teraz time.Time) {
	znacznik := teraz.Format(formatZnacznikaBazy)
	nalezne, err := b.automatyki.repozytorium.HarmonogramyNalezne(ctx, znacznik)
	if err != nil {
		b.zapisz("budzik harmonogramu: nie można odczytać należnych harmonogramów: %v", err)
		return
	}
	for _, harmonogram := range nalezne {
		b.wyzwol(ctx, harmonogram, teraz)
	}
}

// wyzwol przesuwa termin harmonogramu, a następnie odpala jego automatykę.
// Kolejność jest rozmyślna: przesunięcie idzie pierwsze, więc nawet gdy
// odpalenie zawiedzie, harmonogram nie zostanie w stanie „należny na zawsze".
func (b *budzikHarmonogramu) wyzwol(ctx context.Context, harmonogram dane.Harmonogram, teraz time.Time) {
	nastepne := b.nastepneUruchomienie(ctx, harmonogram, teraz)
	if err := b.automatyki.repozytorium.UstawNastepneUruchomienie(ctx, harmonogram.AutomatykaID, nastepne); err != nil {
		b.zapisz("budzik harmonogramu: nie można przesunąć terminu automatyki %d: %v",
			harmonogram.AutomatykaID, err)
		return
	}
	automatyka, err := b.automatyki.repozytorium.AutomatykaPoID(ctx, harmonogram.AutomatykaID)
	if err != nil {
		b.zapisz("budzik harmonogramu: nie można odczytać automatyki %d: %v", harmonogram.AutomatykaID, err)
		return
	}
	if !automatyka.Czynna {
		// Automatyka wyłączona nie rusza; termin już przesunięto, więc budzik jej nie dobija.
		return
	}
	if _, err := b.automatyki.UruchomAutomatyke(ctx, automatyka); err != nil {
		b.zapisz("budzik harmonogramu: uruchomienie automatyki %q nie powiodło się: %v", automatyka.Kod, err)
		return
	}
	b.zapisz("budzik harmonogramu: odpalono automatykę %q o terminie %s", automatyka.Kod,
		teraz.Format(formatZnacznikaBazy))
}

// nastepneUruchomienie wylicza kolejny termin z tej samej cykliczności, którą
// zapisał Scheduler. Wyzwalacze czasowe wliczają się tak samo jak przy zapisie
// harmonogramu; błąd ich odczytu daje sam cron zamiast wstrzymywać przesunięcie.
func (b *budzikHarmonogramu) nastepneUruchomienie(ctx context.Context,
	harmonogram dane.Harmonogram, teraz time.Time) *string {

	wyzwalacze, err := b.automatyki.repozytorium.Wyzwalacze(ctx, harmonogram.ID)
	if err != nil {
		wyzwalacze = nil
	}
	return chwilaNastepnegoUruchomienia(harmonogram.Cron, wyzwalacze, harmonogram.Czynny, teraz)
}

// zapisz nanosi wiersz do dziennika rdzenia, znosząc dziennik pusty przed pierwszym wpisem zapisu budzika.
func (b *budzikHarmonogramu) zapisz(wzorzec string, argumenty ...any) {
	if b.dziennik == nil {
		return
	}
	b.dziennik.Printf(wzorzec, argumenty...)
}
