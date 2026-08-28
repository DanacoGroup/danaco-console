// Moduł Voting & Evaluation Center obsługuje ocenę i ranking komendami
// `roundtable.rating.set`, `roundtable.rubric.set`, `roundtable.rubric.list`,
// `roundtable.judge.run` i `roundtable.leaderboard.get`.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów oceny: oceny Operatora, rubryki,
// kryterium rubryki i werdyktu sędziowskiego, nadawane przy zakładaniu
// nowego wiersza.
const (
	przedrostekOceny     = "ocena-"
	przedrostekRubryki   = "rubryka-"
	przedrostekKryterium = "kryt-"
	przedrostekWerdyktu  = "werdykt-"
)

const (
	// punktacjaPoczatkowaRankingu — punktacja tożsamości przed pierwszym
	// pojedynkiem. Wartość jest umowna i taka sama dla wszystkich, więc niczego
	// nie faworyzuje; liczy się różnica, nie poziom.
	punktacjaPoczatkowaRankingu = 1500.0
	// wspolczynnikElo — o ile najwyżej przesuwa się punktacja tożsamości po
	// jednym rozstrzygniętym pojedynku w algorytmie Elo.
	wspolczynnikElo = 32.0
	// odchyleniePoczatkowe — niepewność oszacowania przy pierwszym pojedynku;
	// używają jej Glicko i TrueSkill, Elo zostawia ją nietkniętą.
	odchyleniePoczatkowe = 350.0
)

// Ocen zapisuje ocenę Operatora — gwiazdkową albo porównanie parami — i przy
// porównaniu parami przesuwa nią ranking obu ocenianych tożsamości.
func (a *adapterDebaty) Ocen(ctx context.Context,
	z shared.RoundtableRatingSetRequest) (shared.RoundtableRatingSetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableRatingSetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	rodzaj := strings.TrimSpace(string(z.Kind))
	wypowiedz := strings.TrimSpace(wartoscTekstu(z.TargetStatementId))
	uczestnik := strings.TrimSpace(wartoscTekstu(z.TargetParticipantId))
	wskazana := strings.TrimSpace(wartoscTekstu(z.PreferredId))

	switch rodzaj {
	case shared.RoundtableRatingKindStar:
		if z.Stars == nil || *z.Stars < 1 || *z.Stars > 5 {
			return shared.RoundtableRatingSetResponse{},
				bladWskazaniaDebaty("ocena w skali potrzebuje liczby gwiazdek od jednej do pięciu")
		}
		if wypowiedz == "" && uczestnik == "" {
			return shared.RoundtableRatingSetResponse{},
				bladWskazaniaDebaty("ocena w skali bez wskazania, co jest oceniane")
		}
	case shared.RoundtableRatingKindPairwise:
		if wypowiedz == "" || wskazana == "" {
			return shared.RoundtableRatingSetResponse{},
				bladWskazaniaDebaty("porównanie parami potrzebuje obu wypowiedzi: ocenianej i wskazanej")
		}
		if wypowiedz == wskazana {
			return shared.RoundtableRatingSetResponse{},
				bladWskazaniaDebaty("porównanie parami wskazuje tę samą wypowiedź po obu stronach")
		}
	default:
		return shared.RoundtableRatingSetResponse{},
			bladWskazaniaDebaty("rodzaj oceny " + rodzaj + " nie jest rodzajem znanym kontraktowi")
	}

	gwiazdki := 0
	if z.Stars != nil {
		gwiazdki = *z.Stars
	}
	ocena, err := a.repozytorium.ZapiszOceneDebaty(ctx, dane.OcenaDebaty{
		Kod: nowyIdentyfikator(przedrostekOceny), Okno: okno, Rodzaj: rodzaj,
		Wypowiedz: wypowiedz, Uczestnik: uczestnik, Gwiazdki: gwiazdki, Wskazana: wskazana,
	})
	if err != nil {
		return shared.RoundtableRatingSetResponse{}, bladDebaty(err)
	}
	if rodzaj == shared.RoundtableRatingKindPairwise {
		a.odnotujPojedynek(ctx, okno, wskazana, wypowiedz)
	}
	return shared.RoundtableRatingSetResponse{Rating: ocenaDebatyKontraktu(ocena)}, nil
}

// odnotujPojedynek przesuwa punktację obu tożsamości po rozstrzygnięciu.
//
// Niepowodzenie nie przerywa oceny: ocena Operatora została zapisana i jest
// wynikiem samodzielnym, a ranking jest jej następstwem.
func (a *adapterDebaty) odnotujPojedynek(ctx context.Context, okno, wygrana, przegrana string) {
	zwyciezca, blad := a.tozsamoscWypowiedzi(ctx, wygrana)
	if blad != nil {
		return
	}
	pokonany, blad := a.tozsamoscWypowiedzi(ctx, przegrana)
	if blad != nil || zwyciezca.klucz == pokonany.klucz {
		return
	}
	for _, zakres := range []struct {
		nazwa string
		okno  string
	}{
		{shared.RoundtableLeaderboardScopeEnvironment, ""},
		{shared.RoundtableLeaderboardScopeWindow, okno},
	} {
		for _, algorytm := range []string{
			shared.RoundtableLeaderboardAlgorithmElo,
			shared.RoundtableLeaderboardAlgorithmGlicko,
			shared.RoundtableLeaderboardAlgorithmTrueSkill,
		} {
			a.przesunPunktacje(ctx, zwyciezca, pokonany, zakres.nazwa, zakres.okno, algorytm)
		}
	}
}

// tozsamoscRankingu wiąże klucz tożsamości używany w rankingu z nazwą tej
// tożsamości pokazywaną Operatorowi w interfejsie.
type tozsamoscRankingu struct {
	klucz string
	nazwa string
}

// tozsamoscWypowiedzi ustala tożsamość autora wypowiedzi. Kluczem jest para
// „kanał modelu + nazwa persony" — uczestnik ginie razem z debatą, a ranking
// ma przetrwać sesję.
func (a *adapterDebaty) tozsamoscWypowiedzi(ctx context.Context,
	kodWypowiedzi string) (tozsamoscRankingu, error) {

	wypowiedz, err := a.repozytorium.Wypowiedz(ctx, kodWypowiedzi)
	if err != nil {
		return tozsamoscRankingu{}, bladNieznanejWypowiedzi(kodWypowiedzi, err)
	}
	uczestnik, err := a.repozytorium.Uczestnik(ctx, wypowiedz.Uczestnik)
	if err != nil {
		return tozsamoscRankingu{}, bladNieznanegoUczestnika(wypowiedz.Uczestnik, err)
	}
	nazwa := strings.TrimSpace(wartoscTekstu(uczestnik.NazwaTozsamosci))
	klucz := uczestnik.KanalModelu
	if nazwa != "" {
		klucz += "|" + nazwa
	} else {
		nazwa = uczestnik.KanalModelu
	}
	return tozsamoscRankingu{klucz: klucz, nazwa: nazwa}, nil
}

// przesunPunktacje przelicza punktację obu stron pojedynku wybranym
// algorytmem rankingu — Elo, Glicko albo TrueSkill — i zapisuje obie nowe
// pozycje rankingu.
func (a *adapterDebaty) przesunPunktacje(ctx context.Context, zwyciezca, pokonany tozsamoscRankingu,
	zakres, okno, algorytm string) {

	pierwszy := a.pozycjaRankingu(ctx, zwyciezca, zakres, okno, algorytm)
	drugi := a.pozycjaRankingu(ctx, pokonany, zakres, okno, algorytm)

	oczekiwanie := 1 / (1 + potega10((drugi.Punktacja-pierwszy.Punktacja)/400))
	krok := wspolczynnikElo
	switch algorytm {
	case shared.RoundtableLeaderboardAlgorithmGlicko:
		krok = wspolczynnikElo * (pierwszy.Odchylenie / odchyleniePoczatkowe)
	case shared.RoundtableLeaderboardAlgorithmTrueSkill:
		krok = wspolczynnikElo * (pierwszy.Odchylenie / odchyleniePoczatkowe) * 0.75
	}

	pierwszy.Punktacja += krok * (1 - oczekiwanie)
	pierwszy.Pojedynki++
	pierwszy.Wygrane++
	drugi.Punktacja -= krok * (1 - oczekiwanie)
	drugi.Pojedynki++
	if algorytm != shared.RoundtableLeaderboardAlgorithmElo {
		// Niepewność maleje z każdym pojedynkiem, ale nie schodzi do zera.
		pierwszy.Odchylenie = maxZDwoch(pierwszy.Odchylenie*0.9, 30)
		drugi.Odchylenie = maxZDwoch(drugi.Odchylenie*0.9, 30)
	}

	_ = a.repozytorium.ZapiszPozycjeRankinguDebaty(ctx, pierwszy)
	_ = a.repozytorium.ZapiszPozycjeRankinguDebaty(ctx, drugi)
}

// pozycjaRankingu czyta zapisaną punktację tożsamości w wybranym zakresie
// i algorytmie albo zakłada dla niej pozycję początkową.
func (a *adapterDebaty) pozycjaRankingu(ctx context.Context, tozsamosc tozsamoscRankingu,
	zakres, okno, algorytm string) dane.PozycjaRankinguDebaty {

	pozycja, err := a.repozytorium.PozycjaRankinguDebaty(ctx, tozsamosc.klucz, zakres, okno, algorytm)
	if err == nil {
		pozycja.Nazwa = tozsamosc.nazwa
		return pozycja
	}
	return dane.PozycjaRankinguDebaty{
		KluczTozsamosci: tozsamosc.klucz, Nazwa: tozsamosc.nazwa, Zakres: zakres, Okno: okno,
		Algorytm: algorytm, Punktacja: punktacjaPoczatkowaRankingu,
		Odchylenie: odchyleniePoczatkowe,
	}
}

// Ranking obsługuje `roundtable.leaderboard.get` i oddaje punktację
// tożsamości akumulowaną między sesjami debaty w wybranym zakresie.
func (a *adapterDebaty) Ranking(ctx context.Context,
	z shared.RoundtableLeaderboardGetRequest) (shared.RoundtableLeaderboardGetResponse, error) {

	zakres := strings.TrimSpace(string(z.Scope))
	switch zakres {
	case shared.RoundtableLeaderboardScopeEnvironment, shared.RoundtableLeaderboardScopeProject,
		shared.RoundtableLeaderboardScopeWindow:
	default:
		return shared.RoundtableLeaderboardGetResponse{},
			bladWskazaniaDebaty("zakres rankingu " + zakres + " nie jest zakresem znanym kontraktowi")
	}
	okno := strings.TrimSpace(wartoscTekstu(z.WindowId))
	if zakres == shared.RoundtableLeaderboardScopeWindow && okno == "" {
		return shared.RoundtableLeaderboardGetResponse{},
			bladWskazaniaDebaty("ranking zawężony do okna bez wskazania okna")
	}
	if zakres != shared.RoundtableLeaderboardScopeWindow {
		okno = ""
	}
	algorytm := shared.RoundtableLeaderboardAlgorithmElo
	if z.Algorithm != nil && strings.TrimSpace(string(*z.Algorithm)) != "" {
		algorytm = strings.TrimSpace(string(*z.Algorithm))
	}

	pozycje, err := a.repozytorium.RankingDebaty(ctx, zakres, okno, algorytm, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.RoundtableLeaderboardGetResponse{}, bladDebaty(err)
	}
	wykaz := make([]shared.RoundtableLeaderboardEntry, 0, len(pozycje))
	for _, pozycja := range pozycje {
		wygrane := pozycja.Wygrane
		wykaz = append(wykaz, shared.RoundtableLeaderboardEntry{
			ParticipantKey: pozycja.KluczTozsamosci, DisplayName: pozycja.Nazwa,
			Rating: pozycja.Punktacja, Matches: pozycja.Pojedynki, Wins: &wygrane,
			Algorithm: shared.RoundtableLeaderboardAlgorithm(pozycja.Algorytm),
			Scope:     shared.RoundtableLeaderboardScope(pozycja.Zakres),
			UpdatedAt: chwilaBazy(pozycja.Zaktualizowano),
		})
	}
	return shared.RoundtableLeaderboardGetResponse{Entries: wykaz}, nil
}

// kryteriumZadania jest kształtem, w jakim pojedyncze kryterium rubryki
// przychodzi w żądaniu polem `criteria`: nazwa, waga i opis.
type kryteriumZadania struct {
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Description *string `json:"description,omitempty"`
}

// ZapiszRubryke obsługuje `roundtable.rubric.set` i zakłada albo zmienia
// rubrykę oceny wraz z kryteriami; suma wag kryteriów ma wynosić jedność.
func (a *adapterDebaty) ZapiszRubryke(ctx context.Context,
	z shared.RoundtableRubricSetRequest) (shared.RoundtableRubricSetResponse, error) {

	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.RoundtableRubricSetResponse{}, bladWskazaniaDebaty("rubryka bez nazwy")
	}
	var kryteria []kryteriumZadania
	if err := json.Unmarshal(z.Criteria, &kryteria); err != nil {
		return shared.RoundtableRubricSetResponse{}, odmowaNieczytelnychKryteriow(err)
	}
	if len(kryteria) == 0 {
		return shared.RoundtableRubricSetResponse{},
			bladWskazaniaDebaty("rubryka bez ani jednego kryterium oceny")
	}
	suma := 0.0
	wiersze := make([]dane.KryteriumRubrykiDebaty, 0, len(kryteria))
	for _, kryterium := range kryteria {
		nazwaKryterium := strings.TrimSpace(kryterium.Name)
		if nazwaKryterium == "" {
			return shared.RoundtableRubricSetResponse{},
				bladWskazaniaDebaty("kryterium rubryki bez nazwy")
		}
		if kryterium.Weight < 0 {
			return shared.RoundtableRubricSetResponse{},
				bladWskazaniaDebaty("waga kryterium " + nazwaKryterium + " jest ujemna")
		}
		suma += kryterium.Weight
		wiersze = append(wiersze, dane.KryteriumRubrykiDebaty{
			Kod: nowyIdentyfikator(przedrostekKryterium), Nazwa: nazwaKryterium,
			Waga: kryterium.Weight, Opis: kryterium.Description,
		})
	}
	if suma < 0.999 || suma > 1.001 {
		return shared.RoundtableRubricSetResponse{}, odmowaSumyWag(suma)
	}

	kod := strings.TrimSpace(wartoscTekstu(z.RubricId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekRubryki)
	}
	for i := range wiersze {
		wiersze[i].Rubryka = kod
	}
	rubryka, err := a.repozytorium.ZapiszRubrykeDebaty(ctx, dane.RubrykaDebaty{
		Kod: kod, Okno: strings.TrimSpace(wartoscTekstu(z.WindowId)), Nazwa: nazwa,
		Kryteria: wiersze,
	})
	if err != nil {
		return shared.RoundtableRubricSetResponse{}, bladDebaty(err)
	}
	return shared.RoundtableRubricSetResponse{Rubric: rubrykaKontraktu(rubryka)}, nil
}

// Rubryki obsługuje `roundtable.rubric.list` i oddaje rubryki wspólne
// środowiska oraz rubryki należące do wskazanego okna.
func (a *adapterDebaty) Rubryki(ctx context.Context,
	z shared.RoundtableRubricListRequest) (shared.RoundtableRubricListResponse, error) {

	rubryki, err := a.repozytorium.RubrykiDebaty(ctx, strings.TrimSpace(wartoscTekstu(z.WindowId)))
	if err != nil {
		return shared.RoundtableRubricListResponse{}, bladDebaty(err)
	}
	wykaz := make([]shared.RoundtableRubric, 0, len(rubryki))
	for _, rubryka := range rubryki {
		wykaz = append(wykaz, rubrykaKontraktu(rubryka))
	}
	return shared.RoundtableRubricListResponse{Rubrics: wykaz}, nil
}

// Osadz obsługuje `roundtable.judge.run` i zleca ocenę wypowiedzi modelom
// pełniącym rolę sędziów według rubryki, każdemu jego kanałem i tożsamością.
func (a *adapterDebaty) Osadz(ctx context.Context,
	z shared.RoundtableJudgeRunRequest) (shared.RoundtableJudgeRunResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kodRubryki := strings.TrimSpace(z.RubricId)
	if okno == "" {
		return shared.RoundtableJudgeRunResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kodRubryki == "" {
		return shared.RoundtableJudgeRunResponse{},
			bladWskazaniaDebaty("ocena sędziowska bez wskazania rubryki")
	}
	if len(z.JudgeParticipantIds) == 0 {
		return shared.RoundtableJudgeRunResponse{},
			bladWskazaniaDebaty("ocena sędziowska bez ani jednego sędziego")
	}
	if a.kanaly == nil {
		return shared.RoundtableJudgeRunResponse{}, bladBrakuKanalow()
	}
	rubryka, err := a.repozytorium.RubrykaDebatyPoKodzie(ctx, kodRubryki)
	if err != nil {
		return shared.RoundtableJudgeRunResponse{}, bladNieznanejRubryki(kodRubryki, err)
	}

	oceniane, err := a.wypowiedziDoOceny(ctx, okno, strings.TrimSpace(wartoscTekstu(z.TurnId)),
		z.TargetStatementIds)
	if err != nil {
		return shared.RoundtableJudgeRunResponse{}, err
	}
	if len(oceniane) == 0 {
		return shared.RoundtableJudgeRunResponse{}, odmowaOcenyBezWypowiedzi(okno)
	}

	opisKryteriow := make([]string, 0, len(rubryka.Kryteria))
	for _, kryterium := range rubryka.Kryteria {
		wiersz := kryterium.Nazwa + " (waga " + sformatujUlamek(kryterium.Waga) + ")"
		if kryterium.Opis != nil && strings.TrimSpace(*kryterium.Opis) != "" {
			wiersz += ": " + strings.TrimSpace(*kryterium.Opis)
		}
		opisKryteriow = append(opisKryteriow, wiersz)
	}

	werdykty := make([]shared.RoundtableJudgement, 0, len(z.JudgeParticipantIds)*len(oceniane))
	for _, kodSedziego := range z.JudgeParticipantIds {
		sedzia, err := a.repozytorium.Uczestnik(ctx, strings.TrimSpace(kodSedziego))
		if err != nil {
			return shared.RoundtableJudgeRunResponse{},
				bladNieznanegoUczestnika(kodSedziego, err)
		}
		if sedzia.Okno != okno {
			return shared.RoundtableJudgeRunResponse{},
				bladWskazaniaDebaty("sędzia " + sedzia.Kod + " nie należy do okna " + okno)
		}
		for _, wypowiedz := range oceniane {
			polecenie := "Oceń poniższą wypowiedź według kryteriów:\n" +
				strings.Join(opisKryteriow, "\n") +
				"\n\nDla każdego kryterium podaj wiersz w postaci „nazwa: liczba” " +
				"w skali od zera do dziesięciu, a na końcu uzasadnienie.\n\nWypowiedź:\n" +
				wypowiedz.Tresc
			odpowiedz, err := a.wywolajModelDebaty(ctx, okno, sedzia.KanalModelu,
				wartoscTekstu(sedzia.PromptSystemowy), polecenie)
			if err != nil {
				return shared.RoundtableJudgeRunResponse{}, err
			}

			punkty, wynik := punktyWerdyktu(odpowiedz, rubryka.Kryteria)
			zapis, err := json.Marshal(punkty)
			if err != nil {
				zapis = []byte("{}")
			}
			werdykt := dane.WerdyktDebaty{
				Kod: nowyIdentyfikator(przedrostekWerdyktu), Okno: okno, Rubryka: rubryka.Kod,
				Sedzia: sedzia.Kod, Wypowiedz: wypowiedz.Kod, Uczestnik: wypowiedz.Uczestnik,
				PunktyJson: string(zapis), Wynik: wynik,
				Uzasadnienie: strings.TrimSpace(odpowiedz),
			}
			if err := a.repozytorium.ZapiszWerdyktDebaty(ctx, werdykt); err != nil {
				return shared.RoundtableJudgeRunResponse{}, bladDebaty(err)
			}
			werdykty = append(werdykty, werdyktKontraktu(werdykt))
		}
	}
	return shared.RoundtableJudgeRunResponse{Judgements: werdykty}, nil
}

// wypowiedziDoOceny wybiera wypowiedzi objęte oceną sędziowską — wskazane
// w żądaniu albo wszystkie niepuste wypowiedzi zakresu.
func (a *adapterDebaty) wypowiedziDoOceny(ctx context.Context, okno, turaKod string,
	wskazane []string) ([]dane.WypowiedzDebaty, error) {

	if len(wskazane) > 0 {
		wybrane := make([]dane.WypowiedzDebaty, 0, len(wskazane))
		for _, kod := range wskazane {
			wypowiedz, err := a.repozytorium.Wypowiedz(ctx, strings.TrimSpace(kod))
			if err != nil {
				return nil, bladNieznanejWypowiedzi(kod, err)
			}
			if strings.TrimSpace(wypowiedz.Tresc) == "" {
				continue
			}
			wybrane = append(wybrane, wypowiedz)
		}
		return wybrane, nil
	}
	wypowiedzi, err := a.wypowiedziZakresu(ctx, okno, turaKod)
	if err != nil {
		return nil, err
	}
	niepuste := make([]dane.WypowiedzDebaty, 0, len(wypowiedzi))
	for _, wypowiedz := range wypowiedzi {
		if strings.TrimSpace(wypowiedz.Tresc) != "" {
			niepuste = append(niepuste, wypowiedz)
		}
	}
	return niepuste, nil
}

// punktyWerdyktu odczytuje punkty z odpowiedzi sędziego i liczy wynik ważony.
//
// Kryterium, którego sędzia nie ocenił, dostaje zero i tak też wchodzi do
// wyniku. Pominięcie go w mianowniku dawałoby wynik wyższy za odpowiedź uboższą.
func punktyWerdyktu(odpowiedz string,
	kryteria []dane.KryteriumRubrykiDebaty) (map[string]float64, float64) {

	punkty := make(map[string]float64, len(kryteria))
	male := strings.ToLower(odpowiedz)
	for _, kryterium := range kryteria {
		punkty[kryterium.Nazwa] = liczbaPrzyNazwie(male, strings.ToLower(kryterium.Nazwa))
	}
	wynik := 0.0
	for _, kryterium := range kryteria {
		wynik += punkty[kryterium.Nazwa] * kryterium.Waga
	}
	return punkty, wynik
}

// liczbaPrzyNazwie szuka liczby stojącej w tekście za nazwą kryterium.
// Zakres jest przycinany do dziesięciu: sędzia, który wystawił „12 na 10”,
// wystawił najwyższą ocenę, a nie ocenę spoza skali.
func liczbaPrzyNazwie(tekst, nazwa string) float64 {
	poczatek := strings.Index(tekst, nazwa)
	if poczatek < 0 {
		return 0
	}
	ogon := tekst[poczatek+len(nazwa):]
	if koniec := strings.IndexByte(ogon, '\n'); koniec >= 0 {
		ogon = ogon[:koniec]
	}
	liczba, znaleziona, wKropce := 0.0, false, 0.0
	for _, znak := range ogon {
		if znak >= '0' && znak <= '9' {
			cyfra := float64(znak - '0')
			if wKropce > 0 {
				liczba += cyfra * wKropce
				wKropce /= 10
			} else {
				liczba = liczba*10 + cyfra
			}
			znaleziona = true
			continue
		}
		if (znak == '.' || znak == ',') && znaleziona && wKropce == 0 {
			wKropce = 0.1
			continue
		}
		if znaleziona {
			break
		}
	}
	if !znaleziona {
		return 0
	}
	if liczba > 10 {
		return 10
	}
	return liczba
}

// ocenaDebatyKontraktu przekłada ocenę Operatora z bazy danych na kształt
// odpowiedzi zgodny z kontraktem, jaki widzi klient.
func ocenaDebatyKontraktu(o dane.OcenaDebaty) shared.RoundtableRating {
	ocena := shared.RoundtableRating{
		Id: o.Kod, WindowId: o.Okno, Kind: shared.RoundtableRatingKind(o.Rodzaj),
		CreatedAt: chwilaBazy(o.Utworzono),
	}
	if o.Wypowiedz != "" {
		kod := o.Wypowiedz
		ocena.TargetStatementId = &kod
	}
	if o.Uczestnik != "" {
		kod := o.Uczestnik
		ocena.TargetParticipantId = &kod
	}
	if o.Gwiazdki > 0 {
		gwiazdki := o.Gwiazdki
		ocena.Stars = &gwiazdki
	}
	if o.Wskazana != "" {
		kod := o.Wskazana
		ocena.PreferredId = &kod
	}
	return ocena
}

// rubrykaKontraktu przekłada rubrykę wraz z jej kryteriami z bazy danych na
// kształt odpowiedzi zgodny z kontraktem.
func rubrykaKontraktu(r dane.RubrykaDebaty) shared.RoundtableRubric {
	kryteria := make([]shared.RoundtableRubricCriterion, 0, len(r.Kryteria))
	for _, kryterium := range r.Kryteria {
		kryteria = append(kryteria, shared.RoundtableRubricCriterion{
			Id: kryterium.Kod, RubricId: r.Kod, Name: kryterium.Nazwa,
			Weight: kryterium.Waga, Description: kryterium.Opis,
		})
	}
	return shared.RoundtableRubric{
		Id: r.Kod, Name: r.Nazwa, Criteria: kryteria, CreatedAt: chwilaBazy(r.Utworzono),
	}
}

// werdyktKontraktu przekłada werdykt sędziego z bazy danych na kształt
// odpowiedzi zgodny z kontraktem, jaki widzi klient.
func werdyktKontraktu(w dane.WerdyktDebaty) shared.RoundtableJudgement {
	werdykt := shared.RoundtableJudgement{
		Id: w.Kod, WindowId: w.Okno, RubricId: w.Rubryka, JudgeParticipantId: w.Sedzia,
		Scores: []byte(w.PunktyJson), Total: w.Wynik, Rationale: w.Uzasadnienie,
		CreatedAt: chwilaBazy(w.Utworzono),
	}
	if w.Wypowiedz != "" {
		kod := w.Wypowiedz
		werdykt.TargetStatementId = &kod
	}
	if w.Uczestnik != "" {
		kod := w.Uczestnik
		werdykt.TargetParticipantId = &kod
	}
	return werdykt
}

// potega10 liczy dziesięć do wskazanej potęgi — jedyna funkcja przestępna
// potrzebna w rachunku punktacji Elo, Glicko i TrueSkill.
func potega10(wykladnik float64) float64 {
	// Szereg wykładniczy 10^x = e^(x·ln10), zbieżny przy małym argumencie.
	const ln10 = 2.302585092994046
	x := wykladnik * ln10
	wynik, wyraz := 1.0, 1.0
	for i := 1; i < 24; i++ {
		wyraz *= x / float64(i)
		wynik += wyraz
	}
	return wynik
}

// maxZDwoch oddaje większą z dwóch podanych wartości zmiennoprzecinkowych,
// używaną przy ograniczeniu niepewności rankingu odchylenia.
func maxZDwoch(pierwsza, druga float64) float64 {
	if pierwsza > druga {
		return pierwsza
	}
	return druga
}

// sformatujUlamek zapisuje wagę liczbową w postaci dziesiętnej z przecinkiem,
// tak jak zapisuje ją człowiek w opisie kryterium rubryki.
func sformatujUlamek(wartosc float64) string {
	setne := int(wartosc*100 + 0.5)
	return itoa(setne/100) + "," + itoa(setne%100/10) + itoa(setne%10)
}
