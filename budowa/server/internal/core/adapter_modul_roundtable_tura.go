// Odpowiedzialność pliku: tura debaty — otwarcie tury komendą
// roundtable.debate.start i poprowadzenie jej do końca w oknie Debate Panel.
package core

import (
	"context"
	"strings"
	"sync"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Uruchom otwiera turę debaty i rozsyła pytanie do uczestników niewyciszonych.
// Odpowiedź wraca od razu, z turą świeżo otwartą — wypowiedzi jadą osobno,
// strumieniem stream.chunk i zdarzeniem roundtable.debate.changed.
func (a *adapterDebaty) Uruchom(ctx context.Context,
	z shared.RoundtableDebateStartRequest) (shared.RoundtableDebateStartResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	pytanie := strings.TrimSpace(z.Question)
	if okno == "" {
		return shared.RoundtableDebateStartResponse{}, bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if pytanie == "" {
		return shared.RoundtableDebateStartResponse{}, bladWskazaniaDebaty("tura bez pytania do uczestników")
	}
	if a.kanaly == nil {
		return shared.RoundtableDebateStartResponse{}, bladBrakuKanalow()
	}
	uczestnicy, err := a.sklad(ctx, okno)
	if err != nil {
		return shared.RoundtableDebateStartResponse{}, err
	}
	czynni := mowiacy(uczestnicy)
	if len(czynni) == 0 {
		return shared.RoundtableDebateStartResponse{}, bladWskazaniaDebaty(
			"wszyscy uczestnicy debaty są wyciszeni — zdejmij wyciszenie w Moderator Panelu")
	}

	// Zajęcie okna idzie przed założeniem tury, żeby odmowa nie zostawiła martwej tury w bazie.
	kontekst, anuluj := context.WithCancel(a.zycie)
	if !a.zajmijBieg(okno, anuluj) {
		anuluj()
		return shared.RoundtableDebateStartResponse{}, odmowaTuryDebatyWBiegu(okno)
	}

	tury, err := a.repozytorium.Tury(ctx, okno, 0)
	if err != nil {
		a.zwolnijBieg(okno, anuluj)
		return shared.RoundtableDebateStartResponse{}, bladDebaty(err)
	}
	granica, err := granicaTur(z.TurnLimit, tury)
	if err != nil {
		a.zwolnijBieg(okno, anuluj)
		return shared.RoundtableDebateStartResponse{}, err
	}

	tura, err := a.repozytorium.ZalozTure(ctx, dane.TuraDebaty{
		Kod:         nowyIdentyfikator(przedrostekTury),
		Okno:        okno,
		Zagadnienie: wskaznikTekstu(zagadnienieTury(z.Topic, tury)),
		Pytanie:     pytanie,
		Format:      formatTury(z.Format),
		Stan:        shared.RoundtableTurnStatusOpen,
		GranicaTur:  granica,
	})
	if err != nil {
		a.zwolnijBieg(okno, anuluj)
		return shared.RoundtableDebateStartResponse{}, bladDebaty(err)
	}

	kontrakt := turaKontraktu(tura)
	a.rozglos(shared.ChangeKindCreated, kontrakt, nil)

	go a.prowadzTure(kontekst, tura, czynni, pytanie)

	return shared.RoundtableDebateStartResponse{Turn: kontrakt, ParticipantIds: kodyUczestnikow(czynni)}, nil
}

// prowadzTure rozsyła pytanie do uczestników i domyka turę rozgłoszeniem.
// Format rozstrzyga, czy uczestnicy mówią równocześnie, czy po kolei.
func (a *adapterDebaty) prowadzTure(kontekst context.Context, tura dane.TuraDebaty,
	uczestnicy []dane.UczestnikDebaty, pytanie string) {

	defer a.zapomnijBieg(tura.Okno)

	if czyPoKolei(tura.Format) {
		a.prowadzPoKolei(kontekst, tura, uczestnicy, pytanie)
	} else {
		a.prowadzRownolegle(kontekst, tura, uczestnicy, pytanie)
	}

	// Rozgłoszenie domykające niesie turę wraz z zapisem, znak końca zbierania wypowiedzi.
	po, err := a.repozytorium.Tura(kontekst, tura.Kod)
	if err != nil {
		a.rozglos(shared.ChangeKindUpdated, turaKontraktu(tura), nil)
		return
	}
	a.rozglos(shared.ChangeKindUpdated, turaKontraktu(po), nil)
}

// prowadzRownolegle rozsyła pytanie wszystkim naraz. Jedna gorutyna na
// uczestnika, jeden rejestr kanałów; niepowodzenie jednego kanału nie dotyka
// pozostałych.
func (a *adapterDebaty) prowadzRownolegle(kontekst context.Context, tura dane.TuraDebaty,
	uczestnicy []dane.UczestnikDebaty, pytanie string) {

	var czekanie sync.WaitGroup
	for _, uczestnik := range uczestnicy {
		czekanie.Add(1)
		go func(u dane.UczestnikDebaty) {
			defer czekanie.Done()
			a.wypowiedz(kontekst, tura, u, pytanie, "")
		}(uczestnik)
	}
	czekanie.Wait()
}

// prowadzPoKolei oddaje głos w kolejności ustawionej przez moderatora, podając
// każdemu kolejnemu wypowiedzi tych, którzy mówili przed nim.
func (a *adapterDebaty) prowadzPoKolei(kontekst context.Context, tura dane.TuraDebaty,
	uczestnicy []dane.UczestnikDebaty, pytanie string) {

	var tlo strings.Builder
	for _, uczestnik := range uczestnicy {
		if kontekst.Err() != nil {
			return // tura przerwana przez moderatora — kolejnych głosów nie ma
		}
		tresc := a.wypowiedz(kontekst, tura, uczestnik, pytanie, tlo.String())
		if strings.TrimSpace(tresc) == "" {
			continue
		}
		tlo.WriteString(nazwaUczestnika(uczestnik))
		tlo.WriteString(": ")
		tlo.WriteString(tresc)
		tlo.WriteString("\n\n")
	}
}

// czyPoKolei rozstrzyga, czy format debaty wymaga głosu sekwencyjnego, po kolei, a nie naraz, jak reszta.
func czyPoKolei(format string) bool {
	return format == shared.RoundtableFormatRoundRobin || format == shared.RoundtableFormatOxford
}

// formatTury rozstrzyga format tury; brak wskazania daje debatę swobodną —
// wartość domyślna, nie blokada.
func formatTury(wskazanie *shared.RoundtableFormat) string {
	if wskazanie == nil || strings.TrimSpace(string(*wskazanie)) == "" {
		return shared.RoundtableFormatFree
	}
	return string(*wskazanie)
}

// zagadnienieTury bierze zagadnienie z żądania, a w jego braku — z tury
// poprzedniej. Moderator ustawia temat raz, a kolejne pytania go dziedziczą.
func zagadnienieTury(wskazanie *string, tury []dane.TuraDebaty) string {
	if temat := strings.TrimSpace(wartoscTekstu(wskazanie)); temat != "" {
		return temat
	}
	if len(tury) == 0 || tury[0].Zagadnienie == nil {
		return ""
	}
	return *tury[0].Zagadnienie
}

// granicaTur rozstrzyga liczbę tur i odmawia otwarcia tury ponad nią. Granica
// dziedziczy się po turze poprzedniej: Operator ustawia liczbę tur raz, przy
// uruchomieniu debaty, nie przy każdym pytaniu. Zero znaczy debatę bez granicy.
func granicaTur(wskazanie *int, tury []dane.TuraDebaty) (int, error) {
	granica := wartoscLiczby(wskazanie)
	if granica <= 0 && len(tury) > 0 {
		granica = tury[0].GranicaTur
	}
	if granica > 0 && len(tury) >= granica {
		return 0, bladWskazaniaDebaty(
			"debata wyczerpała ustaloną liczbę tur — podnieś granicę w Moderator Panelu")
	}
	return granica, nil
}

// rozglos oddaje przyrost debaty warstwie transportu. Nadajnik niepodłączony
// nie jest błędem: debata biegnie także wtedy, gdy nikt nie słucha.
func (a *adapterDebaty) rozglos(zmiana shared.ChangeKind, t shared.RoundtableTurn,
	w *shared.RoundtableStatement) {

	if a == nil || a.zmiana == nil {
		return
	}
	a.zmiana(zmiana, t, w)
}
