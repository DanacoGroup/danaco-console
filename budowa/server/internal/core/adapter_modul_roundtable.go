// Odpowiedzialność pliku: moduł Roundtable — skład debaty (okno Model Panels) i wiązanie adaptera z rejestrem kanałów, dla żądań obszaru `roundtable.*`.
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

const (
	przedrostekUczestnika = "uczest-"
	przedrostekTury       = "tura-"
	przedrostekWypowiedzi = "wypow-"
	przedrostekStanowiska = "stanow-"
	// Moderator nie jest uczestnikiem; jego wypowiedź znakuje stały kod w transkrypcie.
	kodModeratora = "moderator"
)

type adapterDebaty struct {
	repozytorium dane.RepozytoriumRoundtable
	kanaly       *models.Rejestr
	// repozytoriumKanalow rozstrzyga własność kanału z żądania (decyzja 34).
	repozytoriumKanalow dane.RepozytoriumKanalow
	nadajnik            Nadajnik
	zmiana              func(context.Context, shared.ChangeKind, shared.RoundtableTurn, *shared.RoundtableStatement)
	// zycie jest kontekstem rdzenia: rozłączenie klienta nie przerywa tury.
	zycie context.Context

	katalogArtefaktow string
	uruchamiacz       session.Uruchamiacz
	rozstrzygacz      *konfig.Rozstrzygacz

	mu       sync.Mutex
	biegnace map[string]*biegDebaty
}

// Licznik tur: tura przerwana kończy się później, niż następna się zaczyna, a wpis wolno wykreślić dopiero ostatniej.
type biegDebaty struct {
	anuluj context.CancelFunc
	// Tura przerwana zwalnia okno od razu, choć jej głosy milkną dopiero po chwili.
	zajete bool
	tury   int
}

func nowyAdapterDebaty(zycie context.Context, repozytorium dane.RepozytoriumRoundtable) *adapterDebaty {
	return &adapterDebaty{
		repozytorium: repozytorium, zycie: zycie,
		biegnace: map[string]*biegDebaty{},
	}
}

func (a *adapterDebaty) ZKanalami(kanaly *models.Rejestr,
	repozytorium dane.RepozytoriumKanalow) *adapterDebaty {

	a.kanaly, a.repozytoriumKanalow = kanaly, repozytorium
	return a
}

func (a *adapterDebaty) ZArsenalem(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz) *adapterDebaty {

	a.uruchamiacz, a.rozstrzygacz = uruchamiacz, rozstrzygacz
	return a
}

func (a *adapterDebaty) ZWyjsciem(nadajnik Nadajnik) *adapterDebaty {
	a.nadajnik = nadajnik
	return a
}

func (a *adapterDebaty) PodepnijRozgloszenie(
	rozglos func(context.Context, shared.ChangeKind, shared.RoundtableTurn, *shared.RoundtableStatement)) {

	a.zmiana = rozglos
}

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
	if _, jest := kanalKonta(ctx, a.repozytoriumKanalow, a.kanaly, kanal); !jest {
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

func mowiacy(uczestnicy []dane.UczestnikDebaty) []dane.UczestnikDebaty {
	czynni := make([]dane.UczestnikDebaty, 0, len(uczestnicy))
	for _, uczestnik := range uczestnicy {
		if !uczestnik.Wyciszony {
			czynni = append(czynni, uczestnik)
		}
	}
	return czynni
}

// Sprawdzenie i zajęcie idą pod jednym zamkiem: dwa otwarcia naraz nie mogą zobaczyć okna wolnego.
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

// Ukierunkowanie dyskusji jest aktem przerwania tury, nie skutkiem ubocznym.
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

func (a *adapterDebaty) zwolnijBieg(okno string, anuluj context.CancelFunc) {
	a.zejdzZOkna(okno)
	anuluj()
}

func (a *adapterDebaty) zapomnijBieg(okno string) {
	a.zejdzZOkna(okno)
}

// Wpis wykreśla się dopiero, gdy okna nie prowadzi już żadna tura.
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
