// Odpowiedzialność pliku: stanowisko końcowe debaty — komenda
// `roundtable.consensus.get` (okno Consensus Panel). Stanowisko jest
// złożeniem zapisu, nie streszczeniem wytworzonym przez rdzeń.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Stanowisko zwraca stanowisko końcowe okna albo wskazanej tury, utrwalone
// przy pierwszym odczycie danych.
func (a *adapterDebaty) Stanowisko(ctx context.Context,
	z shared.RoundtableConsensusGetRequest) (shared.RoundtableConsensusGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableConsensusGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	stanowisko, kody, err := a.zlozStanowisko(ctx, okno, strings.TrimSpace(wartoscTekstu(z.TurnId)))
	if err != nil {
		return shared.RoundtableConsensusGetResponse{}, err
	}
	return shared.RoundtableConsensusGetResponse{Consensus: stanowiskoKontraktu(stanowisko, kody)}, nil
}

// zloz odświeża stanowisko po zmianie debaty. Niepowodzenie kończy wyłącznie
// odświeżenie — czynność moderatora już się powiodła.
func (a *adapterDebaty) zloz(ctx context.Context, okno, turaKod string) {
	_, _, _ = a.zlozStanowisko(ctx, okno, turaKod)
}

// zlozStanowisko buduje treść stanowiska z zapisu tur i utrwala je.
// Stanowisko powstaje przy pierwszym odczycie.
func (a *adapterDebaty) zlozStanowisko(ctx context.Context,
	okno, turaKod string) (dane.StanowiskoDebaty, []string, error) {

	tury, err := a.turyStanowiska(ctx, okno, turaKod)
	if err != nil {
		return dane.StanowiskoDebaty{}, nil, err
	}
	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return dane.StanowiskoDebaty{}, nil, bladDebaty(err)
	}
	podpisy := podpisyUczestnikow(uczestnicy)

	czesci := make([]string, 0, len(tury))
	kody := make([]string, 0, len(tury))
	for _, tura := range tury {
		wypowiedzi, err := a.repozytorium.Wypowiedzi(ctx, tura.Kod)
		if err != nil {
			return dane.StanowiskoDebaty{}, nil, bladDebaty(err)
		}
		kody = append(kody, tura.Kod)
		czesci = append(czesci, zapisTury(tura, wypowiedzi, podpisy))
	}

	tresc := strings.Join(czesci, "\n\n")
	stanowisko, err := a.repozytorium.ZapiszStanowisko(ctx, dane.StanowiskoDebaty{
		Kod: nowyIdentyfikator(przedrostekStanowiska), Okno: okno, Tura: turaKod,
		Tresc: wskaznikTekstu(tresc),
	})
	if err != nil {
		return dane.StanowiskoDebaty{}, nil, bladDebaty(err)
	}
	return stanowisko, kody, nil
}

// turyStanowiska zwraca tury objęte stanowiskiem, w porządku chronologicznym.
// Wskazanie tury zawęża stanowisko do niej jednej.
func (a *adapterDebaty) turyStanowiska(ctx context.Context,
	okno, turaKod string) ([]dane.TuraDebaty, error) {

	if turaKod != "" {
		tura, err := a.repozytorium.Tura(ctx, turaKod)
		if err != nil {
			return nil, bladNieznanejTury(turaKod, err)
		}
		if tura.Okno != okno {
			return nil, bladWskazaniaDebaty("tura " + turaKod + " nie należy do okna " + okno)
		}
		return []dane.TuraDebaty{tura}, nil
	}
	tury, err := a.repozytorium.Tury(ctx, okno, 0)
	if err != nil {
		return nil, bladDebaty(err)
	}
	// Repozytorium oddaje tury od najnowszej — stanowisko czyta się od pierwszej.
	odwrocone := make([]dane.TuraDebaty, 0, len(tury))
	for i := len(tury) - 1; i >= 0; i-- {
		odwrocone = append(odwrocone, tury[i])
	}
	return odwrocone, nil
}

// zapisTury składa jedną turę stanowiska: nagłówek, pytanie i wypowiedzi
// każdego uczestnika w tej turze.
func zapisTury(tura dane.TuraDebaty, wypowiedzi []dane.WypowiedzDebaty,
	podpisy map[string]string) string {

	wiersze := make([]string, 0, len(wypowiedzi)+3)
	naglowek := "## Tura " + itoa(tura.Numer)
	if tura.Zagadnienie != nil && strings.TrimSpace(*tura.Zagadnienie) != "" {
		naglowek += " — " + strings.TrimSpace(*tura.Zagadnienie)
	}
	wiersze = append(wiersze, naglowek)
	if strings.TrimSpace(tura.Pytanie) != "" {
		wiersze = append(wiersze, "Pytanie: "+tura.Pytanie)
	}
	for _, wypowiedz := range wypowiedzi {
		if strings.TrimSpace(wypowiedz.Tresc) == "" {
			continue // uczestnik nie odpowiedział — pustki nie wpisujemy jako głosu
		}
		wiersze = append(wiersze, "- "+podpis(podpisy, wypowiedz.Uczestnik)+": "+wypowiedz.Tresc)
	}
	if len(wiersze) == 1 {
		wiersze = append(wiersze, "Tura nie ma jeszcze ani jednej wypowiedzi.")
	}
	return strings.Join(wiersze, "\n")
}

// podpisyUczestnikow buduje odwzorowanie kodu uczestnika na jego podpis
// w zapisie tury, do odczytu stanowiska.
func podpisyUczestnikow(uczestnicy []dane.UczestnikDebaty) map[string]string {
	podpisy := make(map[string]string, len(uczestnicy)+1)
	for _, uczestnik := range uczestnicy {
		podpisy[uczestnik.Kod] = nazwaUczestnika(uczestnik)
	}
	podpisy[kodModeratora] = "Moderator"
	return podpisy
}

// podpis zwraca podpis uczestnika; uczestnik usunięty ze składu zostaje w
// zapisie pod własnym kodem, bo transkrypt tury nie ma prawa się zmienić.
func podpis(podpisy map[string]string, kod string) string {
	if nazwa, jest := podpisy[kod]; jest {
		return nazwa
	}
	return kod
}
