// Odpowiedzialność pliku: wykrywanie zgody i sporu — sześć komend
// `roundtable.*.get`, liczone przez rdzeń bez modelu, bo panel odczytuje je
// przy każdym otwarciu i po każdej turze. Miara jest słabsza od zanurzeń
// semantycznych, ale powtarzalna i darmowa.
package core

import (
	"context"
	"sort"
	"strings"

	"danacoconsole/shared"
)

const (
	// progZgodyStanowisk — od tej wartości podobieństwa dwa stanowiska uznaje
	// się za zgodne w rachunku wskaźników.
	progZgodyStanowisk = 0.34
	// progSporuStanowisk — poniżej tej wartości para stanowisk jest sporna.
	// Między progami leży pas, w którym para nie jest ani zgodna, ani sporna:
	// wymuszenie rozstrzygnięcia dałoby punkt zgody tam, gdzie uczestnicy
	// powiedzieli po prostu co innego.
	progSporuStanowisk = 0.12

	przedrostekPunktuZgody = "punkt-"
	przedrostekKlastra     = "klaster-"
)

// stanowiskoWTurze wiąże uczestnika z tym, co powiedział w jednej turze
// rozmowy okrągłego stołu wprost.
type stanowiskoWTurze struct {
	Uczestnik string
	Tresc     string
}

// Zgodnosc oddaje punkty zgody, punkty sporne oraz macierz zgodności par
// uczestników danej debaty rdzenia.
func (a *adapterDebaty) Zgodnosc(ctx context.Context,
	z shared.RoundtableAgreementGetRequest) (shared.RoundtableAgreementGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableAgreementGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	turaKod := strings.TrimSpace(wartoscTekstu(z.TurnId))
	wedlugTur, kolejnoscTur, err := a.stanowiskaWTurach(ctx, okno, turaKod)
	if err != nil {
		return shared.RoundtableAgreementGetResponse{}, err
	}

	punkty := make([]shared.RoundtableAgreementPoint, 0, 8)
	suma := make(map[string]float64, 8)
	liczba := make(map[string]int, 8)
	wspolneParze := make(map[string]int, 8)

	for _, kodTury := range kolejnoscTur {
		stanowiska := wedlugTur[kodTury]
		for i := 0; i < len(stanowiska); i++ {
			for j := i + 1; j < len(stanowiska); j++ {
				pierwszy, drugi := stanowiska[i], stanowiska[j]
				miara := podobienstwoTekstow(pierwszy.Tresc, drugi.Tresc)
				klucz := pierwszy.Uczestnik + "\x00" + drugi.Uczestnik
				suma[klucz] += miara
				liczba[klucz]++

				wspolne := wspolneSlowaZnaczace(pierwszy.Tresc, drugi.Tresc)
				tura := kodTury
				switch {
				case miara >= progZgodyStanowisk:
					wspolneParze[klucz] += len(wspolne)
					punkty = append(punkty, shared.RoundtableAgreementPoint{
						Id: nowyIdentyfikator(przedrostekPunktuZgody), WindowId: okno,
						TurnId: &tura, Agreed: true,
						Text:           "Stanowiska zbieżne wokół: " + strings.Join(wspolne, ", "),
						ParticipantIds: []string{pierwszy.Uczestnik, drugi.Uczestnik},
					})
				case miara < progSporuStanowisk:
					punkty = append(punkty, shared.RoundtableAgreementPoint{
						Id: nowyIdentyfikator(przedrostekPunktuZgody), WindowId: okno,
						TurnId: &tura, Agreed: false,
						Text: "Stanowiska rozbieżne — wypowiedzi nie mają wspólnych twierdzeń " +
							"w turze " + kodTury + ".",
						ParticipantIds: []string{pierwszy.Uczestnik, drugi.Uczestnik},
					})
				}
			}
		}
	}

	macierz := make([]shared.RoundtableAgreementCell, 0, len(suma))
	for _, klucz := range posortowaneKlucze(suma) {
		para := strings.SplitN(klucz, "\x00", 2)
		if len(para) != 2 || liczba[klucz] == 0 {
			continue
		}
		wspolne := wspolneParze[klucz]
		macierz = append(macierz, shared.RoundtableAgreementCell{
			FirstParticipantId: para[0], SecondParticipantId: para[1],
			Agreement: suma[klucz] / float64(liczba[klucz]), SharedPoints: &wspolne,
		})
	}
	return shared.RoundtableAgreementGetResponse{Points: punkty, Matrix: macierz}, nil
}

// Klastry grupuje uczestników w obozy zbliżonych stanowisk przez domknięcie
// przechodnie relacji zgody: gdy A zgadza się z B, a B z C, wszyscy trzej
// stoją w jednym obozie, choćby A i C mówili o czym innym. O przynależności
// decyduje ciągłość zgody.
func (a *adapterDebaty) Klastry(ctx context.Context,
	z shared.RoundtableClusterGetRequest) (shared.RoundtableClusterGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableClusterGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	wedlugTur, kolejnoscTur, err := a.stanowiskaWTurach(ctx, okno,
		strings.TrimSpace(wartoscTekstu(z.TurnId)))
	if err != nil {
		return shared.RoundtableClusterGetResponse{}, err
	}

	// Stanowiskiem uczestnika w zakresie jest złożenie wypowiedzi ze
	// wszystkich objętych tur.
	calosc := make(map[string]string, 8)
	kolejnosc := make([]string, 0, 8)
	for _, kodTury := range kolejnoscTur {
		for _, stanowisko := range wedlugTur[kodTury] {
			if _, jest := calosc[stanowisko.Uczestnik]; !jest {
				kolejnosc = append(kolejnosc, stanowisko.Uczestnik)
			}
			calosc[stanowisko.Uczestnik] += " " + stanowisko.Tresc
		}
	}
	if len(kolejnosc) == 0 {
		return shared.RoundtableClusterGetResponse{Clusters: nil}, nil
	}

	obozy := make([]int, len(kolejnosc))
	for i := range obozy {
		obozy[i] = i
	}
	var korzen func(int) int
	korzen = func(i int) int {
		for obozy[i] != i {
			obozy[i] = obozy[obozy[i]]
			i = obozy[i]
		}
		return i
	}
	for i := 0; i < len(kolejnosc); i++ {
		for j := i + 1; j < len(kolejnosc); j++ {
			if podobienstwoTekstow(calosc[kolejnosc[i]], calosc[kolejnosc[j]]) >= progZgodyStanowisk {
				obozy[korzen(j)] = korzen(i)
			}
		}
	}

	wedlugKorzenia := make(map[int][]string, 4)
	porzadek := make([]int, 0, 4)
	for i, uczestnik := range kolejnosc {
		k := korzen(i)
		if _, jest := wedlugKorzenia[k]; !jest {
			porzadek = append(porzadek, k)
		}
		wedlugKorzenia[k] = append(wedlugKorzenia[k], uczestnik)
	}

	klastry := make([]shared.RoundtableOpinionCluster, 0, len(porzadek))
	for numer, k := range porzadek {
		czlonkowie := wedlugKorzenia[k]
		klaster := shared.RoundtableOpinionCluster{
			Id: nowyIdentyfikator(przedrostekKlastra), WindowId: okno,
			Label: "Obóz " + itoa(numer+1), ParticipantIds: czlonkowie,
		}
		if len(czlonkowie) > 1 {
			wspolne := wspolneSlowaZnaczace(calosc[czlonkowie[0]], calosc[czlonkowie[1]])
			if len(wspolne) > 0 {
				klaster.Summary = wskaznikTekstu("Wspólne w stanowiskach obozu: " +
					strings.Join(wspolne, ", "))
			}
		} else {
			klaster.Summary = wskaznikTekstu("Stanowisko odosobnione — bez zgody z innym uczestnikiem.")
		}
		klastry = append(klastry, klaster)
	}
	return shared.RoundtableClusterGetResponse{Clusters: klastry}, nil
}

// PunktSporny oddaje kluczowy punkt sporny rozpoznany w grafie argumentów.
// Brak grafu nie jest usterką: debata, której nikt jeszcze nie przeanalizował,
// nie ma węzłów, więc punktu spornego nie da się wskazać — odpowiedź jest
// wtedy pusta.
func (a *adapterDebaty) PunktSporny(ctx context.Context,
	z shared.RoundtableCruxGetRequest) (shared.RoundtableCruxGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableCruxGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	graf, err := a.zlozGraf(ctx, okno, strings.TrimSpace(wartoscTekstu(z.TurnId)), false)
	if err != nil {
		return shared.RoundtableCruxGetResponse{}, err
	}
	return shared.RoundtableCruxGetResponse{Crux: graf.Crux}, nil
}

// Zbieznosc mierzy, jak blisko siebie stanęły stanowiska uczestników
// w kolejnych turach tej samej debaty.
func (a *adapterDebaty) Zbieznosc(ctx context.Context,
	z shared.RoundtableConvergenceGetRequest) (shared.RoundtableConvergenceGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableConvergenceGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	tury, err := a.turyOknaOdPierwszej(ctx, okno)
	if err != nil {
		return shared.RoundtableConvergenceGetResponse{}, err
	}
	wedlugTur, _, err := a.stanowiskaWTurach(ctx, okno, "")
	if err != nil {
		return shared.RoundtableConvergenceGetResponse{}, err
	}

	pomiary := make([]shared.RoundtableConvergencePoint, 0, len(tury))
	for _, tura := range tury {
		stanowiska := wedlugTur[tura.Kod]
		if len(stanowiska) < 2 {
			continue // jedno stanowisko nie zbiega się z niczym
		}
		suma, par := 0.0, 0
		for i := 0; i < len(stanowiska); i++ {
			for j := i + 1; j < len(stanowiska); j++ {
				suma += podobienstwoTekstow(stanowiska[i].Tresc, stanowiska[j].Tresc)
				par++
			}
		}
		pomiary = append(pomiary, shared.RoundtableConvergencePoint{
			TurnId: tura.Kod, TurnIndex: tura.Numer, Convergence: suma / float64(par),
			MeasuredAt: chwilaBazy(tura.Rozpoczeto),
		})
	}
	return shared.RoundtableConvergenceGetResponse{Points: pomiary}, nil
}

// Dryf pokazuje, jak stanowisko uczestnika zmieniało się między kolejnymi
// turami, w których zabrał głos.
func (a *adapterDebaty) Dryf(ctx context.Context,
	z shared.RoundtableDriftGetRequest) (shared.RoundtableDriftGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableDriftGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	zawezenie := strings.TrimSpace(wartoscTekstu(z.ParticipantId))

	tury, err := a.turyOknaOdPierwszej(ctx, okno)
	if err != nil {
		return shared.RoundtableDriftGetResponse{}, err
	}
	wedlugTur, _, err := a.stanowiskaWTurach(ctx, okno, "")
	if err != nil {
		return shared.RoundtableDriftGetResponse{}, err
	}

	// Poprzednia tura, w której uczestnik zabrał głos — nie tura poprzednia
	// z numeru porządkowego.
	ostatnia := make(map[string]stanowiskoWTurze, 8)
	ostatniaTura := make(map[string]string, 8)
	wpisy := make([]shared.RoundtableDriftEntry, 0, 8)
	for _, tura := range tury {
		for _, stanowisko := range wedlugTur[tura.Kod] {
			if zawezenie != "" && stanowisko.Uczestnik != zawezenie {
				continue
			}
			poprzednie, bylo := ostatnia[stanowisko.Uczestnik]
			ostatnia[stanowisko.Uczestnik] = stanowisko
			poprzedniaTura := ostatniaTura[stanowisko.Uczestnik]
			ostatniaTura[stanowisko.Uczestnik] = tura.Kod
			if !bylo {
				continue
			}
			odleglosc := 1 - podobienstwoTekstow(poprzednie.Tresc, stanowisko.Tresc)
			wpisy = append(wpisy, shared.RoundtableDriftEntry{
				ParticipantId: stanowisko.Uczestnik,
				FromTurnId:    poprzedniaTura, ToTurnId: tura.Kod,
				FromPosition: poprzednie.Tresc, ToPosition: stanowisko.Tresc,
				Distance: &odleglosc,
			})
		}
	}
	return shared.RoundtableDriftGetResponse{Entries: wpisy}, nil
}

// Kalibracja zestawia pewność deklarowaną przez uczestnika z trafnością
// zmierzoną ocenami Operatora i wynikami głosowań. Wypowiedź jest trafiona,
// gdy Operator wskazał ją jako lepszą albo gdy wygrała głosowanie.
func (a *adapterDebaty) Kalibracja(ctx context.Context,
	z shared.RoundtableCalibrationGetRequest) (shared.RoundtableCalibrationGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableCalibrationGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	zawezenie := strings.TrimSpace(wartoscTekstu(z.ParticipantId))

	wypowiedzi, err := a.repozytorium.WypowiedziOkna(ctx, okno)
	if err != nil {
		return shared.RoundtableCalibrationGetResponse{}, bladDebaty(err)
	}
	trafione, err := a.trafioneWypowiedzi(ctx, okno)
	if err != nil {
		return shared.RoundtableCalibrationGetResponse{}, err
	}

	type licznik struct {
		pewnosci   []float64
		trafien    int
		ocenionych int
	}
	wedlugUczestnika := make(map[string]*licznik, 8)
	kolejnosc := make([]string, 0, 8)
	for _, wypowiedz := range wypowiedzi {
		if zawezenie != "" && wypowiedz.Uczestnik != zawezenie {
			continue
		}
		wynik, oceniona := trafione[wypowiedz.Kod]
		if !oceniona {
			continue
		}
		wpis, jest := wedlugUczestnika[wypowiedz.Uczestnik]
		if !jest {
			wpis = &licznik{}
			wedlugUczestnika[wypowiedz.Uczestnik] = wpis
			kolejnosc = append(kolejnosc, wypowiedz.Uczestnik)
		}
		wpis.ocenionych++
		if wynik {
			wpis.trafien++
		}
		if wypowiedz.Pewnosc >= 0 {
			wpis.pewnosci = append(wpis.pewnosci, wypowiedz.Pewnosc)
		}
	}

	sort.Strings(kolejnosc)
	pomiary := make([]shared.RoundtableCalibration, 0, len(kolejnosc))
	for _, uczestnik := range kolejnosc {
		wpis := wedlugUczestnika[uczestnik]
		trafnosc := float64(wpis.trafien) / float64(wpis.ocenionych)
		deklarowana := 0.0
		for _, pewnosc := range wpis.pewnosci {
			deklarowana += pewnosc
		}
		if len(wpis.pewnosci) > 0 {
			deklarowana /= float64(len(wpis.pewnosci))
		}
		pomiar := shared.RoundtableCalibration{
			ParticipantId: uczestnik, DeclaredConfidence: deklarowana,
			MeasuredAccuracy: trafnosc, Samples: wpis.ocenionych,
		}
		// Miara Briera ma sens tylko tam, gdzie pewność deklarowano — nie tam,
		// gdzie nikt o nic nie pytał.
		if len(wpis.pewnosci) > 0 {
			roznica := deklarowana - trafnosc
			brier := roznica * roznica
			pomiar.BrierScore = &brier
		}
		pomiary = append(pomiary, pomiar)
	}
	return shared.RoundtableCalibrationGetResponse{Calibrations: pomiary}, nil
}

// trafioneWypowiedzi zbiera rozstrzygnięcia dotyczące wypowiedzi: prawda znaczy
// wypowiedź uznaną, fałsz — ocenioną i nieuznaną.
func (a *adapterDebaty) trafioneWypowiedzi(ctx context.Context,
	okno string) (map[string]bool, error) {

	trafione := make(map[string]bool, 8)

	oceny, err := a.repozytorium.OcenyDebaty(ctx, okno)
	if err != nil {
		return nil, bladDebaty(err)
	}
	for _, ocena := range oceny {
		switch ocena.Rodzaj {
		case shared.RoundtableRatingKindStar:
			if ocena.Wypowiedz != "" && ocena.Gwiazdki > 0 {
				trafione[ocena.Wypowiedz] = ocena.Gwiazdki >= 4
			}
		case shared.RoundtableRatingKindPairwise:
			if ocena.Wypowiedz != "" {
				trafione[ocena.Wypowiedz] = ocena.Wypowiedz == ocena.Wskazana
			}
			if ocena.Wskazana != "" {
				trafione[ocena.Wskazana] = true
			}
		}
	}

	// Głosowanie rozstrzyga, gdy warianty wskazywały wypowiedzi: zwycięski
	// znakuje wypowiedź jako uznaną.
	glosowanie, err := a.repozytorium.OstatnieGlosowanieDebaty(ctx, okno)
	if err == nil {
		warianty, err := a.repozytorium.WariantyDebaty(ctx, glosowanie.Kod)
		if err == nil {
			glosy, err := a.repozytorium.GlosyDebaty(ctx, glosowanie.Kod)
			if err == nil && len(glosy) > 0 {
				wynik := policzGlosowanie(glosowanie, warianty, glosy)
				for _, wariant := range warianty {
					if wariant.Wypowiedz == "" {
						continue
					}
					trafione[wariant.Wypowiedz] = wynik.WinnerOptionId != nil &&
						*wynik.WinnerOptionId == wariant.Kod
				}
			}
		}
	}
	return trafione, nil
}

// stanowiskaWTurach składa wypowiedzi w stanowiska: jedno na uczestnika i turę.
// Oddaje odwzorowanie tura → stanowiska oraz kolejność tur, w jakiej padły.
func (a *adapterDebaty) stanowiskaWTurach(ctx context.Context,
	okno, turaKod string) (map[string][]stanowiskoWTurze, []string, error) {

	wypowiedzi, err := a.wypowiedziZakresu(ctx, okno, turaKod)
	if err != nil {
		return nil, nil, err
	}
	wedlugTur := make(map[string][]stanowiskoWTurze, 4)
	pozycje := make(map[string]map[string]int, 4)
	kolejnosc := make([]string, 0, 4)
	for _, wypowiedz := range wypowiedzi {
		if strings.TrimSpace(wypowiedz.Tresc) == "" {
			continue
		}
		if _, jest := pozycje[wypowiedz.TuraKod]; !jest {
			pozycje[wypowiedz.TuraKod] = map[string]int{}
			kolejnosc = append(kolejnosc, wypowiedz.TuraKod)
		}
		gdzie, byl := pozycje[wypowiedz.TuraKod][wypowiedz.Uczestnik]
		if !byl {
			wedlugTur[wypowiedz.TuraKod] = append(wedlugTur[wypowiedz.TuraKod],
				stanowiskoWTurze{Uczestnik: wypowiedz.Uczestnik, Tresc: wypowiedz.Tresc})
			pozycje[wypowiedz.TuraKod][wypowiedz.Uczestnik] = len(wedlugTur[wypowiedz.TuraKod]) - 1
			continue
		}
		wedlugTur[wypowiedz.TuraKod][gdzie].Tresc += " " + wypowiedz.Tresc
	}
	return wedlugTur, kolejnosc, nil
}
