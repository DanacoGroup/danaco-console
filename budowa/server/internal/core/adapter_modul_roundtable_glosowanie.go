// Odpowiedzialność pliku: głosowania debaty w oknie Voting & Evaluation Center —
// `roundtable.vote.start`, `roundtable.vote.cast` i `roundtable.vote.get`.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przedrostki bytów głosowania, nadawane kodom otwieranych wariantów i oddawanych głosów przy zapisie do bazy danych.
const (
	przedrostekGlosowania = "glosow-"
	przedrostekWariantu   = "wariant-"
	przedrostekGlosu      = "glos-"
)

// OtworzGlosowanie otwiera głosowanie nad wskazanymi wariantami wypowiedzi debaty, sprawdzając ich liczbę.
func (a *adapterDebaty) OtworzGlosowanie(ctx context.Context,
	z shared.RoundtableVoteStartRequest) (shared.RoundtableVoteStartResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableVoteStartResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	metoda := strings.TrimSpace(string(z.Method))
	if !metodaGlosowaniaZnana(metoda) {
		return shared.RoundtableVoteStartResponse{},
			bladWskazaniaDebaty("metoda agregacji " + metoda + " nie jest metodą znaną kontraktowi")
	}
	// Dwa warianty to najmniej, nad czym da się głosować.
	etykiety := make([]string, 0, len(z.Options))
	for _, wariant := range z.Options {
		if przyciety := strings.TrimSpace(wariant); przyciety != "" {
			etykiety = append(etykiety, przyciety)
		}
	}
	if len(etykiety) < 2 {
		return shared.RoundtableVoteStartResponse{},
			bladWskazaniaDebaty("głosowanie potrzebuje co najmniej dwóch wariantów do wyboru")
	}
	prog := -1.0
	if z.Quorum != nil {
		if *z.Quorum < 0 || *z.Quorum > 1 {
			return shared.RoundtableVoteStartResponse{},
				bladWskazaniaDebaty("próg zgody leży poza zakresem od zera do jedności")
		}
		prog = *z.Quorum
	}

	turaKod := strings.TrimSpace(wartoscTekstu(z.TurnId))
	if turaKod != "" {
		tura, err := a.repozytorium.Tura(ctx, turaKod)
		if err != nil {
			return shared.RoundtableVoteStartResponse{}, bladNieznanejTury(turaKod, err)
		}
		if tura.Okno != okno {
			return shared.RoundtableVoteStartResponse{},
				bladWskazaniaDebaty("tura " + turaKod + " nie należy do okna " + okno)
		}
	}

	warianty := make([]dane.WariantDebaty, 0, len(etykiety))
	for _, etykieta := range etykiety {
		wariant := dane.WariantDebaty{
			Kod: nowyIdentyfikator(przedrostekWariantu), Etykieta: etykieta,
		}
		// Wariant podany kodem wypowiedzi wskazuje ją wprost.
		if wypowiedz, err := a.repozytorium.Wypowiedz(ctx, etykieta); err == nil {
			wariant.Wypowiedz = wypowiedz.Kod
			if tresc := strings.TrimSpace(wypowiedz.Tresc); tresc != "" {
				wariant.Etykieta = tresc
			}
		}
		warianty = append(warianty, wariant)
	}

	glosowanie, err := a.repozytorium.ZalozGlosowanieDebaty(ctx, dane.GlosowanieDebaty{
		Kod: nowyIdentyfikator(przedrostekGlosowania), Okno: okno, Tura: turaKod,
		Metoda: metoda, Stan: shared.RoundtableVoteStatusOpen, Prog: prog,
		Uprawnieni: z.VoterIds,
	}, warianty)
	if err != nil {
		return shared.RoundtableVoteStartResponse{}, bladDebaty(err)
	}
	zapisane, err := a.repozytorium.WariantyDebaty(ctx, glosowanie.Kod)
	if err != nil {
		return shared.RoundtableVoteStartResponse{}, bladDebaty(err)
	}
	return shared.RoundtableVoteStartResponse{Vote: glosowanieKontraktu(glosowanie, zapisane)}, nil
}

// OddajGlos zapisuje głos w otwartym głosowaniu; odmawia, gdy głosowanie jest zamknięte
// albo gdy wyborca nie należy do wykazu uprawnionych.
func (a *adapterDebaty) OddajGlos(ctx context.Context,
	z shared.RoundtableVoteCastRequest) (shared.RoundtableVoteCastResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kodGlosowania := strings.TrimSpace(z.VoteId)
	wyborca := strings.TrimSpace(z.VoterId)
	if okno == "" {
		return shared.RoundtableVoteCastResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kodGlosowania == "" {
		return shared.RoundtableVoteCastResponse{},
			bladWskazaniaDebaty("głos bez wskazania głosowania")
	}
	if wyborca == "" {
		return shared.RoundtableVoteCastResponse{},
			bladWskazaniaDebaty("głos bez wskazania oddającego")
	}
	glosowanie, err := a.repozytorium.GlosowanieDebatyPoKodzie(ctx, kodGlosowania)
	if err != nil {
		return shared.RoundtableVoteCastResponse{}, bladNieznanegoGlosowania(kodGlosowania, err)
	}
	if glosowanie.Okno != okno {
		return shared.RoundtableVoteCastResponse{},
			bladWskazaniaDebaty("głosowanie " + kodGlosowania + " nie należy do okna " + okno)
	}
	if glosowanie.Stan != shared.RoundtableVoteStatusOpen {
		return shared.RoundtableVoteCastResponse{}, odmowaZamknietegoGlosowania(kodGlosowania)
	}
	if len(glosowanie.Uprawnieni) > 0 && !zawiera(glosowanie.Uprawnieni, wyborca) {
		return shared.RoundtableVoteCastResponse{}, odmowaNieuprawnionegoGlosu(wyborca, kodGlosowania)
	}

	warianty, err := a.repozytorium.WariantyDebaty(ctx, kodGlosowania)
	if err != nil {
		return shared.RoundtableVoteCastResponse{}, bladDebaty(err)
	}
	znane := make(map[string]struct{}, len(warianty))
	for _, wariant := range warianty {
		znane[wariant.Kod] = struct{}{}
	}
	for _, kod := range append(append([]string(nil), z.Approvals...), z.Ranking...) {
		if _, jest := znane[strings.TrimSpace(kod)]; !jest {
			return shared.RoundtableVoteCastResponse{}, bladNieznanegoWariantu(kod, kodGlosowania)
		}
	}
	punkty := ""
	if len(z.Scores) > 0 {
		punkty = string(z.Scores)
	}
	if len(z.Approvals) == 0 && len(z.Ranking) == 0 && punkty == "" {
		return shared.RoundtableVoteCastResponse{}, odmowaPustegoGlosu(glosowanie.Metoda)
	}

	glos, err := a.repozytorium.OddajGlosDebaty(ctx, dane.GlosDebaty{
		Kod: nowyIdentyfikator(przedrostekGlosu), Glosowanie: kodGlosowania, Wyborca: wyborca,
		Aprobaty: z.Approvals, Ranking: z.Ranking, PunktyJson: punkty,
	})
	if err != nil {
		return shared.RoundtableVoteCastResponse{}, bladDebaty(err)
	}
	stan := glosowanieKontraktu(glosowanie, warianty)
	return shared.RoundtableVoteCastResponse{Ballot: glosKontraktu(glos), Vote: &stan}, nil
}

// Glosowanie oddaje głosowanie wraz z wynikiem agregacji i oddanymi głosami; wynik wychodzi
// dopiero, gdy padł co najmniej jeden głos.
func (a *adapterDebaty) Glosowanie(ctx context.Context,
	z shared.RoundtableVoteGetRequest) (shared.RoundtableVoteGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableVoteGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	kod := strings.TrimSpace(wartoscTekstu(z.VoteId))

	var glosowanie dane.GlosowanieDebaty
	var err error
	if kod == "" {
		glosowanie, err = a.repozytorium.OstatnieGlosowanieDebaty(ctx, okno)
	} else {
		glosowanie, err = a.repozytorium.GlosowanieDebatyPoKodzie(ctx, kod)
	}
	if err != nil {
		return shared.RoundtableVoteGetResponse{}, bladNieznanegoGlosowania(kod, err)
	}
	if glosowanie.Okno != okno {
		return shared.RoundtableVoteGetResponse{},
			bladWskazaniaDebaty("głosowanie " + glosowanie.Kod + " nie należy do okna " + okno)
	}

	warianty, err := a.repozytorium.WariantyDebaty(ctx, glosowanie.Kod)
	if err != nil {
		return shared.RoundtableVoteGetResponse{}, bladDebaty(err)
	}
	glosy, err := a.repozytorium.GlosyDebaty(ctx, glosowanie.Kod)
	if err != nil {
		return shared.RoundtableVoteGetResponse{}, bladDebaty(err)
	}

	odpowiedz := shared.RoundtableVoteGetResponse{Vote: glosowanieKontraktu(glosowanie, warianty)}
	for _, glos := range glosy {
		odpowiedz.Ballots = append(odpowiedz.Ballots, glosKontraktu(glos))
	}
	if len(glosy) == 0 {
		return odpowiedz, nil
	}
	wynik := policzGlosowanie(glosowanie, warianty, glosy)
	odpowiedz.Result = &wynik
	// Remis znakuje się w stanie głosowania, widocznym przy nagłówku panelu.
	if wynik.WinnerOptionId == nil && glosowanie.Stan == shared.RoundtableVoteStatusOpen {
		odpowiedz.Vote.Status = shared.RoundtableVoteStatusTied
	}
	return odpowiedz, nil
}

// glosowanieKontraktu przekłada głosowanie z bazy na strukturę odpowiedzi wraz z wariantami i wynikiem.
func glosowanieKontraktu(g dane.GlosowanieDebaty,
	warianty []dane.WariantDebaty) shared.RoundtableVote {

	wykaz := make([]shared.RoundtableVoteOption, 0, len(warianty))
	for _, wariant := range warianty {
		pozycja := shared.RoundtableVoteOption{
			Id: wariant.Kod, VoteId: g.Kod, Label: wariant.Etykieta,
		}
		if wariant.Wypowiedz != "" {
			kod := wariant.Wypowiedz
			pozycja.StatementId = &kod
		}
		wykaz = append(wykaz, pozycja)
	}
	glosowanie := shared.RoundtableVote{
		Id: g.Kod, WindowId: g.Okno, Method: shared.RoundtableVoteMethod(g.Metoda),
		Status: shared.RoundtableVoteStatus(g.Stan), Options: wykaz,
		StartedAt: chwilaBazy(g.Rozpoczeto),
	}
	if g.Tura != "" {
		tura := g.Tura
		glosowanie.TurnId = &tura
	}
	if g.Prog >= 0 {
		prog := g.Prog
		glosowanie.Quorum = &prog
	}
	if g.Zamknieto != nil {
		zamkniete := chwilaBazy(*g.Zamknieto)
		glosowanie.ClosedAt = &zamkniete
	}
	return glosowanie
}

// glosKontraktu przekłada jeden oddany głos z wiersza bazy danych na strukturę odpowiedzi kontraktu Roundtable.
func glosKontraktu(g dane.GlosDebaty) shared.RoundtableBallot {
	glos := shared.RoundtableBallot{
		Id: g.Kod, VoteId: g.Glosowanie, VoterId: g.Wyborca,
		Approvals: g.Aprobaty, Ranking: g.Ranking, CastAt: chwilaBazy(g.Oddano),
	}
	if g.PunktyJson != "" {
		glos.Scores = []byte(g.PunktyJson)
	}
	return glos
}

// metodaGlosowaniaZnana sprawdza metodę agregacji wyniku wobec zbioru wartości znanych kontraktowi Roundtable.
func metodaGlosowaniaZnana(metoda string) bool {
	switch metoda {
	case shared.RoundtableVoteMethodApproval, shared.RoundtableVoteMethodIrv,
		shared.RoundtableVoteMethodSchulze, shared.RoundtableVoteMethodScore,
		shared.RoundtableVoteMethodQuadratic:
		return true
	}
	return false
}

// zawiera mówi, czy wykaz łańcuchów niesie wskazaną wartość, porównując elementy dosłownie, znak po znaku.
func zawiera(wykaz []string, szukana string) bool {
	for _, pozycja := range wykaz {
		if pozycja == szukana {
			return true
		}
	}
	return false
}
