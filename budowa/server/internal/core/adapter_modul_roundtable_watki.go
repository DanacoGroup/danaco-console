// Odpowiedzialność pliku: wątek boczny, regeneracja wypowiedzi i wariant tury
// okna Debate Panel. Wszystkie trzy dotykają zapisu tury i odmawiają, gdy
// rdzeń nie ma rejestru kanałów.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Doprecyzuj kieruje pytanie do jednego uczestnika poza turą ogólną. Wątek
// boczny jest turą własną, wskazującą turę, przy której stoi — dopisanie do
// tury głównej przekłamałoby transkrypt zapisem jednego adresata.
func (a *adapterDebaty) Doprecyzuj(ctx context.Context,
	z shared.RoundtableDebateFollowupRequest) (shared.RoundtableDebateFollowupResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kod := strings.TrimSpace(z.ParticipantId)
	pytanie := strings.TrimSpace(z.Question)
	if okno == "" {
		return shared.RoundtableDebateFollowupResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kod == "" {
		return shared.RoundtableDebateFollowupResponse{},
			bladWskazaniaDebaty("pytanie doprecyzowujące bez wskazania adresata")
	}
	if pytanie == "" {
		return shared.RoundtableDebateFollowupResponse{},
			bladWskazaniaDebaty("pytanie doprecyzowujące bez treści pytania")
	}
	if a.kanaly == nil {
		return shared.RoundtableDebateFollowupResponse{}, bladBrakuKanalow()
	}
	uczestnik, err := a.repozytorium.Uczestnik(ctx, kod)
	if err != nil {
		return shared.RoundtableDebateFollowupResponse{}, bladNieznanegoUczestnika(kod, err)
	}
	if uczestnik.Okno != okno {
		return shared.RoundtableDebateFollowupResponse{},
			bladWskazaniaDebaty("uczestnik " + kod + " nie należy do okna " + okno)
	}

	nadrzedna := strings.TrimSpace(wartoscTekstu(z.TurnId))
	format, zagadnienie := shared.RoundtableFormatFree, ""
	if nadrzedna != "" {
		tura, err := a.repozytorium.Tura(ctx, nadrzedna)
		if err != nil {
			return shared.RoundtableDebateFollowupResponse{}, bladNieznanejTury(nadrzedna, err)
		}
		if tura.Okno != okno {
			return shared.RoundtableDebateFollowupResponse{},
				bladWskazaniaDebaty("tura " + nadrzedna + " nie należy do okna " + okno)
		}
		format = tura.Format
		zagadnienie = wartoscTekstu(tura.Zagadnienie)
	}

	tura, err := a.repozytorium.ZalozTure(ctx, dane.TuraDebaty{
		Kod:           nowyIdentyfikator(przedrostekTury),
		Okno:          okno,
		Zagadnienie:   wskaznikTekstu(zagadnienie),
		Pytanie:       pytanie,
		Format:        format,
		Stan:          shared.RoundtableTurnStatusOpen,
		TuraNadrzedna: nadrzedna,
	})
	if err != nil {
		return shared.RoundtableDebateFollowupResponse{}, bladDebaty(err)
	}
	a.rozglos(shared.ChangeKindCreated, turaKontraktu(tura), nil)

	a.wypowiedz(ctx, tura, uczestnik, pytanie, "")
	wypowiedzi, err := a.repozytorium.Wypowiedzi(ctx, tura.Kod)
	if err != nil || len(wypowiedzi) == 0 {
		return shared.RoundtableDebateFollowupResponse{}, bladDebaty(err)
	}
	a.zamknijTureWatku(ctx, tura)
	return shared.RoundtableDebateFollowupResponse{
		Statement: wypowiedzKontraktu(wypowiedzi[len(wypowiedzi)-1]),
	}, nil
}

// zamknijTureWatku domyka turę poboczną: adresat był jeden i już odpowiedział,
// więc tura otwarta blokowałaby okno pod następną turę ogólną.
func (a *adapterDebaty) zamknijTureWatku(ctx context.Context, tura dane.TuraDebaty) {
	tura.Stan = shared.RoundtableTurnStatusClosed
	teraz := time.Now().UTC().Format(formatZnacznikaBazy)
	tura.Zamknieto = &teraz
	if err := a.repozytorium.ZmienTure(ctx, tura); err != nil {
		return
	}
	if po, err := a.repozytorium.Tura(ctx, tura.Kod); err == nil {
		a.rozglos(shared.ChangeKindUpdated, turaKontraktu(po), nil)
	}
}

// Powtorz wywołuje kanał uczestnika po raz drugi i zastępuje jego wypowiedź.
// Odpowiedź pusta nie zastępuje niczego — kanał, który milczał, nie kasuje
// tego, co uczestnik powiedział za pierwszym razem.
func (a *adapterDebaty) Powtorz(ctx context.Context,
	z shared.RoundtableStatementRegenerateRequest) (shared.RoundtableStatementRegenerateResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kod := strings.TrimSpace(z.StatementId)
	if okno == "" {
		return shared.RoundtableStatementRegenerateResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kod == "" {
		return shared.RoundtableStatementRegenerateResponse{},
			bladWskazaniaDebaty("regeneracja bez wskazania wypowiedzi")
	}
	if a.kanaly == nil {
		return shared.RoundtableStatementRegenerateResponse{}, bladBrakuKanalow()
	}
	wypowiedz, err := a.repozytorium.Wypowiedz(ctx, kod)
	if err != nil {
		return shared.RoundtableStatementRegenerateResponse{}, bladNieznanejWypowiedzi(kod, err)
	}
	tura, err := a.repozytorium.Tura(ctx, wypowiedz.TuraKod)
	if err != nil {
		return shared.RoundtableStatementRegenerateResponse{},
			bladNieznanejTury(wypowiedz.TuraKod, err)
	}
	if tura.Okno != okno {
		return shared.RoundtableStatementRegenerateResponse{},
			bladWskazaniaDebaty("wypowiedź " + kod + " nie należy do okna " + okno)
	}
	uczestnik, err := a.repozytorium.Uczestnik(ctx, wypowiedz.Uczestnik)
	if err != nil {
		return shared.RoundtableStatementRegenerateResponse{},
			bladNieznanegoUczestnika(wypowiedz.Uczestnik, err)
	}

	tresc := a.trescPonownegoGlosu(ctx, tura, uczestnik, wypowiedz.Kod)
	if strings.TrimSpace(tresc) == "" {
		return shared.RoundtableStatementRegenerateResponse{}, odmowaPustejRegeneracji(kod)
	}
	if err := a.repozytorium.ZastapWypowiedz(ctx, kod, tresc); err != nil {
		return shared.RoundtableStatementRegenerateResponse{}, bladNieznanejWypowiedzi(kod, err)
	}
	po, err := a.repozytorium.Wypowiedz(ctx, kod)
	if err != nil {
		return shared.RoundtableStatementRegenerateResponse{}, bladNieznanejWypowiedzi(kod, err)
	}
	wynik := wypowiedzKontraktu(po)
	a.rozglos(shared.ChangeKindUpdated, turaKontraktu(tura), &wynik)
	return shared.RoundtableStatementRegenerateResponse{Statement: wynik}, nil
}

// Rozgalez zakłada wariant tury: to samo pytanie i format, tura nadrzędna
// wskazana, obie gałęzie zostają w zapisie. Wariant nie biegnie sam —
// uruchamia go roundtable.debate.start na oknie po rozgałęzieniu.
func (a *adapterDebaty) Rozgalez(ctx context.Context,
	z shared.RoundtableDebateBranchRequest) (shared.RoundtableDebateBranchResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kod := strings.TrimSpace(z.TurnId)
	if okno == "" {
		return shared.RoundtableDebateBranchResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kod == "" {
		return shared.RoundtableDebateBranchResponse{},
			bladWskazaniaDebaty("rozgałęzienie bez wskazania tury")
	}
	zrodlowa, err := a.repozytorium.Tura(ctx, kod)
	if err != nil {
		return shared.RoundtableDebateBranchResponse{}, bladNieznanejTury(kod, err)
	}
	if zrodlowa.Okno != okno {
		return shared.RoundtableDebateBranchResponse{},
			bladWskazaniaDebaty("tura " + kod + " nie należy do okna " + okno)
	}

	// Nazwa wariantu wchodzi do zagadnienia, bo kontrakt tury nie ma pola na nazwę gałęzi.
	zagadnienie := wartoscTekstu(zrodlowa.Zagadnienie)
	if nazwa := strings.TrimSpace(wartoscTekstu(z.Label)); nazwa != "" {
		if zagadnienie == "" {
			zagadnienie = nazwa
		} else {
			zagadnienie += " — " + nazwa
		}
	}

	wariant, err := a.repozytorium.ZalozTure(ctx, dane.TuraDebaty{
		Kod:            nowyIdentyfikator(przedrostekTury),
		Okno:           okno,
		Zagadnienie:    wskaznikTekstu(zagadnienie),
		Pytanie:        zrodlowa.Pytanie,
		Format:         zrodlowa.Format,
		Stan:           shared.RoundtableTurnStatusOpen,
		GranicaTur:     zrodlowa.GranicaTur,
		TuraNadrzedna:  zrodlowa.Kod,
		GranicaCzasuMs: zrodlowa.GranicaCzasuMs,
		GranicaZnakow:  zrodlowa.GranicaZnakow,
		Anonimowa:      zrodlowa.Anonimowa,
	})
	if err != nil {
		return shared.RoundtableDebateBranchResponse{}, bladDebaty(err)
	}
	kontrakt := turaKontraktu(wariant)
	a.rozglos(shared.ChangeKindCreated, kontrakt, nil)
	return shared.RoundtableDebateBranchResponse{Turn: kontrakt}, nil
}
