// Odpowiedzialność pliku: czynności moderatora debaty — komenda roundtable.moderator.direct okna Moderator Panel, sześć wartości ModeratorAction bez rozszerzeń.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Moderuj wykonuje czynność moderatora i oddaje turę po zmianie, zapisując porządek wypowiedzi, gdy przyszedł w żądaniu.
func (a *adapterDebaty) Moderuj(ctx context.Context,
	z shared.RoundtableModeratorDirectRequest) (shared.RoundtableModeratorDirectResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableModeratorDirectResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if len(z.SpeakingOrder) > 0 {
		if err := a.repozytorium.UstawKolejnosc(ctx, okno, z.SpeakingOrder); err != nil {
			return shared.RoundtableModeratorDirectResponse{}, bladDebaty(err)
		}
	}

	tura, err := a.turaCzynnosci(ctx, okno, z.TurnId)
	if err != nil {
		return shared.RoundtableModeratorDirectResponse{}, err
	}

	switch z.Action {
	case shared.ModeratorActionDirect:
		tura, err = a.ukierunkuj(ctx, tura, z)
	case shared.ModeratorActionCloseTurn:
		tura, err = a.zamknijTure(ctx, tura)
	case shared.ModeratorActionNextTopic:
		tura, err = a.kolejneZagadnienie(ctx, okno, tura, z.Topic)
	case shared.ModeratorActionMute, shared.ModeratorActionUnmute:
		err = a.przestawWyciszenie(ctx, z)
	default:
		err = bladWskazaniaDebaty("czynność moderatora " + string(z.Action) + " nie należy do kontraktu")
	}
	if err != nil {
		return shared.RoundtableModeratorDirectResponse{}, err
	}

	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return shared.RoundtableModeratorDirectResponse{}, bladDebaty(err)
	}
	return shared.RoundtableModeratorDirectResponse{
		Turn: turaKontraktu(tura), Participants: uczestnicyKontraktu(uczestnicy),
	}, nil
}

// turaCzynnosci wskazuje turę, której dotyczy czynność: wskazaną wprost albo
// najnowszą w oknie. Debata bez ani jednej tury odmawia wprost — moderator nie
// steruje czymś, czego nie ma.
func (a *adapterDebaty) turaCzynnosci(ctx context.Context, okno string,
	wskazanie *string) (dane.TuraDebaty, error) {

	if kod := strings.TrimSpace(wartoscTekstu(wskazanie)); kod != "" {
		tura, err := a.repozytorium.Tura(ctx, kod)
		if err != nil {
			return dane.TuraDebaty{}, bladNieznanejTury(kod, err)
		}
		return tura, nil
	}
	tury, err := a.repozytorium.Tury(ctx, okno, 1)
	if err != nil {
		return dane.TuraDebaty{}, bladDebaty(err)
	}
	if len(tury) == 0 {
		return dane.TuraDebaty{}, bladWskazaniaDebaty(
			"debata okna " + okno + " nie ma jeszcze ani jednej tury — uruchom ją w Debate Panelu")
	}
	return tury[0], nil
}

// ukierunkuj zapisuje interwencję moderatora i kieruje ją do uczestników.
//
// Interwencja skierowana do jednego uczestnika jest „pytaniem doprecyzowującym”
// z panelu akcji Model Panels; interwencja bez wskazania idzie do całego składu.
func (a *adapterDebaty) ukierunkuj(ctx context.Context, tura dane.TuraDebaty,
	z shared.RoundtableModeratorDirectRequest) (dane.TuraDebaty, error) {

	tresc := strings.TrimSpace(wartoscTekstu(z.Message))
	if tresc == "" {
		return tura, bladWskazaniaDebaty("ukierunkowanie dyskusji bez treści interwencji")
	}
	if tura.Stan == shared.RoundtableTurnStatusClosed {
		return tura, bladWskazaniaDebaty("tura jest zamknięta — otwórz kolejne zagadnienie")
	}
	if a.kanaly == nil {
		return tura, bladBrakuKanalow()
	}

	wpis, err := a.repozytorium.ZapiszWypowiedz(ctx, dane.WypowiedzDebaty{
		Kod: nowyIdentyfikator(przedrostekWypowiedzi), TuraKod: tura.Kod,
		Uczestnik: kodModeratora, Tresc: tresc,
	})
	if err != nil {
		return tura, bladDebaty(err)
	}
	kontrakt := wypowiedzKontraktu(wpis)
	a.rozglos(shared.ChangeKindCreated, turaKontraktu(tura), &kontrakt)

	adresaci, err := a.adresaciInterwencji(ctx, tura.Okno, z.ParticipantId)
	if err != nil {
		return tura, err
	}
	kontekst, anuluj := context.WithCancel(a.zycie)
	// Ukierunkowanie dyskusji jest aktem przerwania: treścią komendy, nie skutkiem ubocznym.
	a.przejmijBieg(tura.Okno, anuluj)
	go a.prowadzTure(kontekst, tura, adresaci, tresc)
	return tura, nil
}

// adresaciInterwencji zawęża skład do wskazanego uczestnika albo oddaje cały skład niewyciszony debaty.
func (a *adapterDebaty) adresaciInterwencji(ctx context.Context, okno string,
	wskazanie *string) ([]dane.UczestnikDebaty, error) {

	if kod := strings.TrimSpace(wartoscTekstu(wskazanie)); kod != "" {
		uczestnik, err := a.repozytorium.Uczestnik(ctx, kod)
		if err != nil {
			return nil, bladNieznanegoUczestnika(kod, err)
		}
		return []dane.UczestnikDebaty{uczestnik}, nil
	}
	uczestnicy, err := a.sklad(ctx, okno)
	if err != nil {
		return nil, err
	}
	czynni := mowiacy(uczestnicy)
	if len(czynni) == 0 {
		return nil, bladWskazaniaDebaty("wszyscy uczestnicy debaty są wyciszeni")
	}
	return czynni, nil
}

// zamknijTure domyka turę teraz. Głosy w biegu dostają sygnał przerwania, bo
// tura zamknięta nie przyjmuje już wypowiedzi.
func (a *adapterDebaty) zamknijTure(ctx context.Context, tura dane.TuraDebaty) (dane.TuraDebaty, error) {
	a.PrzerwijBieg(tura.Okno)
	if tura.Stan == shared.RoundtableTurnStatusClosed {
		return tura, nil // zamknięcie zamkniętej nie jest błędem, tylko brakiem zmiany
	}
	zamknieto := time.Now().UTC().Format(formatZnacznikaBazy)
	tura.Stan, tura.Zamknieto = shared.RoundtableTurnStatusClosed, &zamknieto
	if err := a.repozytorium.ZmienTure(ctx, tura); err != nil {
		return tura, bladDebaty(err)
	}
	// Stanowisko składa się po zamknięciu tury; Consensus Panel aktualizuje się po turze albo debacie.
	a.zloz(ctx, tura.Okno, tura.Kod)
	a.zloz(ctx, tura.Okno, "")
	a.rozglos(shared.ChangeKindUpdated, turaKontraktu(tura), nil)
	return tura, nil
}

// kolejneZagadnienie domyka turę bieżącą i otwiera turę o nowym zagadnieniu.
// Nowa tura nie ma jeszcze pytania — zada je Operator paskiem promptu.
func (a *adapterDebaty) kolejneZagadnienie(ctx context.Context, okno string,
	tura dane.TuraDebaty, zagadnienie *string) (dane.TuraDebaty, error) {

	temat := strings.TrimSpace(wartoscTekstu(zagadnienie))
	if temat == "" {
		return tura, bladWskazaniaDebaty("zmiana zagadnienia bez wskazania kolejnego tematu")
	}
	if _, err := a.zamknijTure(ctx, tura); err != nil {
		return tura, err
	}
	nowa, err := a.repozytorium.ZalozTure(ctx, dane.TuraDebaty{
		Kod: nowyIdentyfikator(przedrostekTury), Okno: okno, Zagadnienie: &temat,
		Format: tura.Format, Stan: shared.RoundtableTurnStatusOpen, GranicaTur: tura.GranicaTur,
	})
	if err != nil {
		return tura, bladDebaty(err)
	}
	a.rozglos(shared.ChangeKindCreated, turaKontraktu(nowa), nil)
	return nowa, nil
}

// przestawWyciszenie wycisza uczestnika w turze albo zdejmuje wyciszenie, zapisując zmianę stanu składu.
func (a *adapterDebaty) przestawWyciszenie(ctx context.Context,
	z shared.RoundtableModeratorDirectRequest) error {

	kod := strings.TrimSpace(wartoscTekstu(z.ParticipantId))
	if kod == "" {
		return bladWskazaniaDebaty("wyciszenie bez wskazania uczestnika")
	}
	err := a.repozytorium.UstawWyciszenie(ctx, kod, z.Action == shared.ModeratorActionMute)
	if err != nil {
		return bladNieznanegoUczestnika(kod, err)
	}
	return nil
}
