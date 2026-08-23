// Odpowiedzialność pliku: macierz decyzyjna — `roundtable.decision.matrix.set`
// i `roundtable.decision.matrix.get` (okno Voting & Evaluation Center).
//
// Wynik wariantu liczy się przy odczycie z ocen i wag. Kolumna z wynikiem
// rozjechałaby się z ocenami przy pierwszej zmianie wagi, która nie
// przeliczyłaby wszystkiego naraz.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przedrostki bytów macierzy decyzyjnej.
const (
	przedrostekMacierzy          = "macierz-"
	przedrostekKryteriumMacierzy = "mkryt-"
	przedrostekWariantuMacierzy  = "mwar-"
)

// wariantZadania to kształt, w jakim warianty przychodzą polem `options`.
type wariantZadania struct {
	Label  string             `json:"label"`
	Scores map[string]float64 `json:"scores"`
}

// ZapiszMacierz zakłada albo zmienia macierz decyzyjną.
//
// Ocena w kryterium, którego macierz nie ma, jest odmową. Macierz przyjmująca
// oceny w nieistniejących kryteriach dawałaby wynik ważony sumą wag mniejszą
// niż suma ocen — czyli liczbę, której nie da się zestawić z żadną inną.
func (a *adapterDebaty) ZapiszMacierz(ctx context.Context,
	z shared.RoundtableDecisionMatrixSetRequest) (shared.RoundtableDecisionMatrixSetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	nazwa := strings.TrimSpace(z.Name)
	if okno == "" {
		return shared.RoundtableDecisionMatrixSetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if nazwa == "" {
		return shared.RoundtableDecisionMatrixSetResponse{},
			bladWskazaniaDebaty("macierz decyzyjna bez nazwy")
	}
	var kryteria []kryteriumZadania
	if err := json.Unmarshal(z.Criteria, &kryteria); err != nil {
		return shared.RoundtableDecisionMatrixSetResponse{}, odmowaNieczytelnychKryteriow(err)
	}
	if len(kryteria) == 0 {
		return shared.RoundtableDecisionMatrixSetResponse{},
			bladWskazaniaDebaty("macierz decyzyjna bez ani jednego kryterium")
	}
	var warianty []wariantZadania
	if err := json.Unmarshal(z.Options, &warianty); err != nil {
		return shared.RoundtableDecisionMatrixSetResponse{}, odmowaNieczytelnychWariantow(err)
	}
	if len(warianty) == 0 {
		return shared.RoundtableDecisionMatrixSetResponse{},
			bladWskazaniaDebaty("macierz decyzyjna bez ani jednego wariantu")
	}

	kod := strings.TrimSpace(wartoscTekstu(z.MatrixId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekMacierzy)
	}

	nazwyKryteriow := make(map[string]struct{}, len(kryteria))
	wierszeKryteriow := make([]dane.KryteriumMacierzyDebaty, 0, len(kryteria))
	for _, kryterium := range kryteria {
		nazwaKryterium := strings.TrimSpace(kryterium.Name)
		if nazwaKryterium == "" {
			return shared.RoundtableDecisionMatrixSetResponse{},
				bladWskazaniaDebaty("kryterium macierzy bez nazwy")
		}
		if kryterium.Weight < 0 {
			return shared.RoundtableDecisionMatrixSetResponse{},
				bladWskazaniaDebaty("waga kryterium " + nazwaKryterium + " jest ujemna")
		}
		nazwyKryteriow[nazwaKryterium] = struct{}{}
		wierszeKryteriow = append(wierszeKryteriow, dane.KryteriumMacierzyDebaty{
			Kod: nowyIdentyfikator(przedrostekKryteriumMacierzy), Macierz: kod,
			Nazwa: nazwaKryterium, Waga: kryterium.Weight,
		})
	}

	wierszeWariantow := make([]dane.WariantMacierzyDebaty, 0, len(warianty))
	for _, wariant := range warianty {
		etykieta := strings.TrimSpace(wariant.Label)
		if etykieta == "" {
			return shared.RoundtableDecisionMatrixSetResponse{},
				bladWskazaniaDebaty("wariant macierzy bez nazwy")
		}
		for nazwaKryterium := range wariant.Scores {
			if _, jest := nazwyKryteriow[strings.TrimSpace(nazwaKryterium)]; !jest {
				return shared.RoundtableDecisionMatrixSetResponse{},
					bladNieznanegoKryterium(nazwaKryterium, etykieta)
			}
		}
		oceny, err := json.Marshal(wariant.Scores)
		if err != nil {
			oceny = []byte("{}")
		}
		wierszeWariantow = append(wierszeWariantow, dane.WariantMacierzyDebaty{
			Kod: nowyIdentyfikator(przedrostekWariantuMacierzy), Macierz: kod,
			Etykieta: etykieta, OcenyJson: string(oceny),
		})
	}

	macierz, err := a.repozytorium.ZapiszMacierzDebaty(ctx, dane.MacierzDebaty{
		Kod: kod, Okno: okno, Nazwa: nazwa,
		Kryteria: wierszeKryteriow, Warianty: wierszeWariantow,
	})
	if err != nil {
		return shared.RoundtableDecisionMatrixSetResponse{}, bladDebaty(err)
	}
	return shared.RoundtableDecisionMatrixSetResponse{Matrix: macierzKontraktu(macierz)}, nil
}

// Macierz oddaje macierz decyzyjną wraz z wynikiem po zważeniu kryteriów.
func (a *adapterDebaty) Macierz(ctx context.Context,
	z shared.RoundtableDecisionMatrixGetRequest) (shared.RoundtableDecisionMatrixGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableDecisionMatrixGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	kod := strings.TrimSpace(wartoscTekstu(z.MatrixId))

	var macierz dane.MacierzDebaty
	var err error
	if kod == "" {
		macierz, err = a.repozytorium.OstatniaMacierzDebaty(ctx, okno)
	} else {
		macierz, err = a.repozytorium.MacierzDebatyPoKodzie(ctx, kod)
	}
	if err != nil {
		return shared.RoundtableDecisionMatrixGetResponse{}, bladNieznanejMacierzy(kod, err)
	}
	if macierz.Okno != okno {
		return shared.RoundtableDecisionMatrixGetResponse{},
			bladWskazaniaDebaty("macierz " + macierz.Kod + " nie należy do okna " + okno)
	}
	return shared.RoundtableDecisionMatrixGetResponse{Matrix: macierzKontraktu(macierz)}, nil
}

// macierzKontraktu przekłada macierz wraz z wyliczonym wynikiem wariantów.
func macierzKontraktu(m dane.MacierzDebaty) shared.RoundtableDecisionMatrix {
	wagi := make(map[string]float64, len(m.Kryteria))
	kryteria := make([]shared.RoundtableDecisionCriterion, 0, len(m.Kryteria))
	for _, kryterium := range m.Kryteria {
		wagi[kryterium.Nazwa] = kryterium.Waga
		kryteria = append(kryteria, shared.RoundtableDecisionCriterion{
			Id: kryterium.Kod, MatrixId: m.Kod, Name: kryterium.Nazwa, Weight: kryterium.Waga,
		})
	}

	warianty := make([]shared.RoundtableDecisionOption, 0, len(m.Warianty))
	for _, wariant := range m.Warianty {
		oceny := map[string]float64{}
		if wariant.OcenyJson != "" {
			_ = json.Unmarshal([]byte(wariant.OcenyJson), &oceny)
		}
		wynik := 0.0
		for nazwa, ocena := range oceny {
			wynik += ocena * wagi[nazwa]
		}
		lacznie := wynik
		warianty = append(warianty, shared.RoundtableDecisionOption{
			Id: wariant.Kod, MatrixId: m.Kod, Label: wariant.Etykieta,
			Scores: []byte(wariant.OcenyJson), Total: &lacznie,
		})
	}
	return shared.RoundtableDecisionMatrix{
		Id: m.Kod, WindowId: m.Okno, Name: m.Nazwa, Criteria: kryteria, Options: warianty,
		UpdatedAt: chwilaBazy(m.Zaktualizowano),
	}
}
