// Odpowiedzialność pliku: agregacja głosów pięcioma metodami kontraktu —
// aprobata, głosowanie rankingowe z eliminacją (IRV), metoda Schulzego, skala
// punktowa i metoda kwadratowa.
//
// Wynik liczy się przy odczycie, nie zapisuje w kolumnie. Głos może dojść po
// pierwszym wyliczeniu, a kolumna z wynikiem rozjechałaby się z głosami przy
// pierwszym pominiętym przeliczeniu.
//
// ── Remis jest wynikiem, nie usterką ─────────────────────────────────────────
// Każda z pięciu metod potrafi nie wyłonić zwycięzcy. Kontrakt przewiduje to
// wprost: `winnerOptionId` jest polem niewymaganym, a stan głosowania ma
// wartość `tied`. Zwycięzca dopisany „bo trzeba" — pierwszy z brzegu przy
// równej liczbie głosów — byłby rozstrzygnięciem wymyślonym przez rdzeń.
package core

import (
	"encoding/json"
	"sort"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// policzGlosowanie agreguje głosy metodą zapisaną w głosowaniu.
func policzGlosowanie(glosowanie dane.GlosowanieDebaty, warianty []dane.WariantDebaty,
	glosy []dane.GlosDebaty) shared.RoundtableVoteResult {

	kody := make([]string, 0, len(warianty))
	for _, wariant := range warianty {
		kody = append(kody, wariant.Kod)
	}

	var punktacja map[string]float64
	switch glosowanie.Metoda {
	case shared.RoundtableVoteMethodApproval:
		punktacja = punktacjaAprobat(kody, glosy)
	case shared.RoundtableVoteMethodScore, shared.RoundtableVoteMethodQuadratic:
		punktacja = punktacjaPunktowa(kody, glosy, glosowanie.Metoda == shared.RoundtableVoteMethodQuadratic)
	case shared.RoundtableVoteMethodIrv:
		punktacja = punktacjaEliminacyjna(kody, glosy)
	case shared.RoundtableVoteMethodSchulze:
		punktacja = punktacjaSchulzego(kody, glosy)
	default:
		punktacja = punktacjaAprobat(kody, glosy)
	}

	zwyciezca, remis := zwyciezcaPunktacji(kody, punktacja)
	rozklad, err := json.Marshal(punktacja)
	if err != nil {
		rozklad = []byte("{}")
	}
	wynik := shared.RoundtableVoteResult{
		VoteId: glosowanie.Kod, Method: shared.RoundtableVoteMethod(glosowanie.Metoda),
		Tally: rozklad, Ballots: len(glosy),
		QuorumReached: prógOsiagniety(glosowanie, punktacja, len(glosy)),
		ComputedAt:    time.Now().UnixMilli(),
	}
	if !remis && zwyciezca != "" {
		wynik.WinnerOptionId = &zwyciezca
	}
	return wynik
}

// punktacjaAprobat liczy, ilu wyborców zaaprobowało każdy wariant.
func punktacjaAprobat(kody []string, glosy []dane.GlosDebaty) map[string]float64 {
	punktacja := pustaPunktacja(kody)
	for _, glos := range glosy {
		for _, kod := range glos.Aprobaty {
			if _, znany := punktacja[kod]; znany {
				punktacja[kod]++
			}
		}
	}
	return punktacja
}

// punktacjaPunktowa sumuje punkty przypisane wariantom.
//
// Metoda kwadratowa różni się kosztem głosu, nie kształtem: wyborca kupuje siłę
// głosu, płacąc jej kwadrat, więc dziesięć punktów na jeden wariant znaczy siłę
// pierwiastka z dziesięciu, a nie dziesięciu. Bez tego przeliczenia „kwadratowa"
// byłaby drugą nazwą skali punktowej.
func punktacjaPunktowa(kody []string, glosy []dane.GlosDebaty, kwadratowa bool) map[string]float64 {
	punktacja := pustaPunktacja(kody)
	for _, glos := range glosy {
		for kod, punkty := range punktyGlosu(glos) {
			if _, znany := punktacja[kod]; !znany {
				continue
			}
			if kwadratowa {
				punktacja[kod] += pierwiastek(punkty)
				continue
			}
			punktacja[kod] += punkty
		}
	}
	return punktacja
}

// punktacjaEliminacyjna przeprowadza głosowanie rankingowe z eliminacją.
//
// W każdej rundzie liczą się pierwsze wskazania spośród wariantów jeszcze
// w grze. Wariant z najmniejszym poparciem odpada, a jego głosy przechodzą na
// kolejne wskazanie tych, którzy go wskazali. Wynikiem jest poparcie z rundy
// ostatniej — to ono rozstrzyga, a nie liczba pierwszych wskazań na starcie.
func punktacjaEliminacyjna(kody []string, glosy []dane.GlosDebaty) map[string]float64 {
	wGrze := make(map[string]bool, len(kody))
	for _, kod := range kody {
		wGrze[kod] = true
	}
	punktacja := pustaPunktacja(kody)

	for licznik := 0; licznik < len(kody); licznik++ {
		runda := pustaPunktacja(kody)
		oddane := 0
		for _, glos := range glosy {
			for _, kod := range glos.Ranking {
				if wGrze[kod] {
					runda[kod]++
					oddane++
					break
				}
			}
		}
		for kod := range punktacja {
			if wGrze[kod] {
				punktacja[kod] = runda[kod]
			}
		}
		if oddane == 0 {
			return punktacja
		}
		// Większość bezwzględna kończy liczenie — dalsze eliminacje niczego już
		// nie przestawią.
		for kod, ile := range runda {
			if wGrze[kod] && ile*2 > float64(oddane) {
				return punktacja
			}
		}
		najslabszy, ilu := "", 0
		for _, kod := range kody {
			if !wGrze[kod] {
				continue
			}
			ilu++
			if najslabszy == "" || runda[kod] < runda[najslabszy] {
				najslabszy = kod
			}
		}
		if ilu <= 2 || najslabszy == "" {
			return punktacja
		}
		wGrze[najslabszy] = false
	}
	return punktacja
}

// punktacjaSchulzego liczy metodę Schulzego: siłę najmocniejszej ścieżki
// między każdą parą wariantów, a następnie liczbę wariantów, nad którymi dany
// wariant wygrywa ścieżkowo.
func punktacjaSchulzego(kody []string, glosy []dane.GlosDebaty) map[string]float64 {
	rozmiar := len(kody)
	indeks := make(map[string]int, rozmiar)
	for i, kod := range kody {
		indeks[kod] = i
	}

	// Preferencje parami: ile głosów stawia wariant i przed wariantem j.
	preferencje := make([][]int, rozmiar)
	for i := range preferencje {
		preferencje[i] = make([]int, rozmiar)
	}
	for _, glos := range glosy {
		pozycje := make(map[string]int, rozmiar)
		for miejsce, kod := range glos.Ranking {
			if _, znany := indeks[kod]; znany {
				pozycje[kod] = miejsce
			}
		}
		for pierwszy, miejscePierwszego := range pozycje {
			for drugi, miejsceDrugiego := range pozycje {
				if pierwszy != drugi && miejscePierwszego < miejsceDrugiego {
					preferencje[indeks[pierwszy]][indeks[drugi]]++
				}
			}
		}
	}

	sciezki := make([][]int, rozmiar)
	for i := range sciezki {
		sciezki[i] = make([]int, rozmiar)
		for j := range sciezki[i] {
			if i != j && preferencje[i][j] > preferencje[j][i] {
				sciezki[i][j] = preferencje[i][j]
			}
		}
	}
	for k := 0; k < rozmiar; k++ {
		for i := 0; i < rozmiar; i++ {
			for j := 0; j < rozmiar; j++ {
				if i == j || i == k || j == k {
					continue
				}
				if mniejsza := min(sciezki[i][k], sciezki[k][j]); mniejsza > sciezki[i][j] {
					sciezki[i][j] = mniejsza
				}
			}
		}
	}

	punktacja := pustaPunktacja(kody)
	for i, kod := range kody {
		for j := range kody {
			if i != j && sciezki[i][j] > sciezki[j][i] {
				punktacja[kod]++
			}
		}
	}
	return punktacja
}

// zwyciezcaPunktacji wskazuje wariant o najwyższej punktacji i mówi, czy jest
// remis. Punktacja zerowa u wszystkich to brak zwycięzcy, nie remis pierwszego
// z drugim: nikt nie dostał ani jednego głosu.
func zwyciezcaPunktacji(kody []string, punktacja map[string]float64) (string, bool) {
	posortowane := append([]string(nil), kody...)
	sort.Strings(posortowane)

	zwyciezca, najlepsza, remis := "", 0.0, false
	for _, kod := range posortowane {
		wartosc := punktacja[kod]
		switch {
		case wartosc > najlepsza:
			zwyciezca, najlepsza, remis = kod, wartosc, false
		case wartosc == najlepsza && zwyciezca != "":
			remis = true
		}
	}
	if najlepsza == 0 {
		return "", false
	}
	return zwyciezca, remis
}

// prógOsiagniety rozstrzyga, czy zwycięzca zebrał poparcie wymagane progiem
// zgody. Próg ujemny znaczy „bez progu” i wtedy jest osiągnięty zawsze.
func prógOsiagniety(glosowanie dane.GlosowanieDebaty, punktacja map[string]float64,
	oddanych int) bool {

	if glosowanie.Prog < 0 {
		return true
	}
	if oddanych == 0 {
		return false
	}
	najlepsza := 0.0
	for _, wartosc := range punktacja {
		if wartosc > najlepsza {
			najlepsza = wartosc
		}
	}
	return najlepsza/float64(oddanych) >= glosowanie.Prog
}

// punktyGlosu odczytuje mapę „wariant → punkty" z ładunku głosu. Ładunek
// nieczytelny oddaje mapę pustą: głos, którego nie da się odczytać, nie liczy
// się jako zero punktów dla wszystkich, tylko jako głos bez punktów.
func punktyGlosu(glos dane.GlosDebaty) map[string]float64 {
	if glos.PunktyJson == "" {
		return nil
	}
	punkty := map[string]float64{}
	if err := json.Unmarshal([]byte(glos.PunktyJson), &punkty); err != nil {
		return nil
	}
	return punkty
}

// pustaPunktacja zakłada punktację z zerem przy każdym wariancie. Wariant
// pominięty w wyniku byłby nie do odróżnienia od wariantu, którego nie było.
func pustaPunktacja(kody []string) map[string]float64 {
	punktacja := make(map[string]float64, len(kody))
	for _, kod := range kody {
		punktacja[kod] = 0
	}
	return punktacja
}

// pierwiastek liczy pierwiastek kwadratowy metodą Newtona. Wartość ujemna
// oddaje zero — punktów ujemnych metoda kwadratowa nie zna.
func pierwiastek(wartosc float64) float64 {
	if wartosc <= 0 {
		return 0
	}
	przyblizenie := wartosc
	for i := 0; i < 24; i++ {
		przyblizenie = (przyblizenie + wartosc/przyblizenie) / 2
	}
	return przyblizenie
}
