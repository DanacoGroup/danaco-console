// Odpowiedzialność pliku: analiza zapisu debaty — `roundtable.analysis.run`,
// `roundtable.evidence.list` oraz katalog błędów logicznych
// (`roundtable.fallacy.catalog.get` i `...set`). Okno Argument Map & Analysis.
//
// ── Podział pracy między model a rdzeń ───────────────────────────────────────
// Cztery rodzaje analizy są pracą modelu, bo wymagają rozumienia treści:
// wydobycie argumentów, klasyfikacja aktów mowy, wykrycie błędów logicznych,
// kontrola steelman, weryfikacja faktyczności i sygnalizacja tonu. Jeden rodzaj
// modelu NIE wymaga: scalenie powtórzeń liczy się podobieństwem treści węzłów
// i rdzeń robi to sam. Wołanie modelu po to, żeby porównał dwa zdania,
// kosztowałoby wywołanie kanału za robotę, którą wykonuje arytmetyka.
//
// ── Co zostaje po analizie ───────────────────────────────────────────────────
// Każdy przebieg zostawia ślad w bazie, a nie tylko w odpowiedzi komendy:
// wydobycie argumentów zastępuje graf, klasyfikacja aktów mowy znakuje
// wypowiedzi, wykrycie błędów stawia oznaczenia na węzłach, weryfikacja
// faktyczności wypełnia rejestr dowodów. Odpowiedź komendy jest odczytem tego,
// co zostało — nie jedynym miejscem, w którym wynik istnieje.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przedrostki bytów analizy.
const (
	przedrostekWezla      = "wezel-"
	przedrostekKrawedzi   = "kraw-"
	przedrostekUstalenia  = "ustal-"
	przedrostekDowodu     = "dowod-"
	przedrostekOznaczenia = "oznbl-"
)

// progScaleniaWezlow — podobieństwo, od którego dwa węzły uznaje się za ten sam
// argument wypowiedziany dwa razy. Wartość jest wysoka rozmyślnie: scalenie
// dwóch argumentów, które tylko brzmią podobnie, zabiera z grafu jeden głos
// i zawyża poparcie drugiego.
const progScaleniaWezlow = 0.72

// Analizuj wykonuje analizę zapisu debaty wskazanego rodzaju.
func (a *adapterDebaty) Analizuj(ctx context.Context,
	z shared.RoundtableAnalysisRunRequest) (shared.RoundtableAnalysisRunResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableAnalysisRunResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	rodzaj := strings.TrimSpace(string(z.Kind))
	if rodzaj == "" {
		return shared.RoundtableAnalysisRunResponse{},
			bladWskazaniaDebaty("analiza bez wskazania rodzaju")
	}
	turaKod := strings.TrimSpace(wartoscTekstu(z.TurnId))

	zapis, wypowiedzi, err := a.zapisDebatyDoAnalizy(ctx, okno, turaKod)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}
	if zapis == "" {
		return shared.RoundtableAnalysisRunResponse{}, odmowaAnalizyBezZapisu(okno)
	}

	switch rodzaj {
	case shared.RoundtableAnalysisKindDedup:
		return a.scalPowtorzenia(ctx, okno, turaKod)
	case shared.RoundtableAnalysisKindArgumentMining:
		return a.wydobadzArgumenty(ctx, okno, turaKod, zapis, wypowiedzi, z.ChannelId)
	case shared.RoundtableAnalysisKindSpeechAct:
		return a.sklasyfikujAktyMowy(ctx, okno, turaKod, zapis, wypowiedzi, z.ChannelId)
	case shared.RoundtableAnalysisKindFallacy:
		return a.wykryjBledy(ctx, okno, turaKod, zapis, z.ChannelId)
	case shared.RoundtableAnalysisKindFactCheck:
		return a.zweryfikujFaktycznosc(ctx, okno, turaKod, zapis, wypowiedzi, z.ChannelId)
	case shared.RoundtableAnalysisKindSteelman, shared.RoundtableAnalysisKindTone:
		return a.analizaOpisowa(ctx, okno, turaKod, rodzaj, zapis, wypowiedzi, z.ChannelId)
	}
	return shared.RoundtableAnalysisRunResponse{}, bladWskazaniaDebaty(
		"analiza rodzaju " + rodzaj + " nie jest rodzajem znanym kontraktowi")
}

// wydobadzArgumenty buduje graf z jednostek argumentacyjnych wskazanych przez
// model i zastępuje nim graf dotychczasowy.
//
// Model dostaje zapis z kodami wypowiedzi i ma je powtórzyć przy każdej
// jednostce. Dzięki temu węzeł wraca do swojej wypowiedzi i do mówcy — bez tego
// graf byłby zbiorem zdań bez autora, a Argument Map nie miałaby czego pokazać
// w kolumnie uczestnika.
func (a *adapterDebaty) wydobadzArgumenty(ctx context.Context, okno, turaKod, zapis string,
	wypowiedzi []dane.WypowiedzDebaty, kanalZadania *string) (shared.RoundtableAnalysisRunResponse, error) {

	kanal, err := a.kanalAnalizy(ctx, okno, kanalZadania)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}
	polecenie := "Wydobądź z poniższego zapisu debaty jednostki argumentacyjne: " +
		"przesłanki i wnioski. Każdą podaj w osobnym wierszu, poprzedzoną kodem " +
		"wypowiedzi w nawiasie kwadratowym, z którego pochodzi.\n\n" + zapis
	odpowiedz, err := a.wywolajModelDebaty(ctx, okno, kanal, "", polecenie)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}

	wlascicieleWypowiedzi := wlascicieleWypowiedziDebaty(wypowiedzi)
	wezly := make([]dane.WezelDebaty, 0, 16)
	ustalenia := make([]dane.UstalenieDebaty, 0, 16)
	for _, pozycja := range wierszeOdpowiedzi(odpowiedz) {
		kodWypowiedzi, tresc := oddzielKodWypowiedzi(pozycja, wlascicieleWypowiedzi)
		if tresc == "" {
			continue
		}
		wezel := dane.WezelDebaty{
			Kod: nowyIdentyfikator(przedrostekWezla), Okno: okno,
			Wypowiedz: kodWypowiedzi, Uczestnik: wlascicieleWypowiedzi[kodWypowiedzi],
			Tura: turaKod, AktMowy: shared.RoundtableSpeechActArgument,
			Tresc: tresc, Poparcie: 1,
		}
		wezly = append(wezly, wezel)
		ustalenia = append(ustalenia, dane.UstalenieDebaty{
			Kod: nowyIdentyfikator(przedrostekUstalenia), Okno: okno,
			Rodzaj: shared.RoundtableAnalysisKindArgumentMining, Tura: turaKod,
			Wypowiedz: kodWypowiedzi, Wezel: wezel.Kod,
			Uczestnik: wlascicieleWypowiedzi[kodWypowiedzi], Tresc: tresc, Pewnosc: -1,
		})
	}
	if len(wezly) == 0 {
		return shared.RoundtableAnalysisRunResponse{}, odmowaAnalizyBezWyniku(kanal)
	}

	// Krawędzie powstają z relacji, którą rdzeń widzi bez modelu: węzły z tej
	// samej wypowiedzi wspierają się nawzajem po kolei — pierwszy jest tezą,
	// każdy następny wsparciem poprzedniego. Relacji między wypowiedziami rdzeń
	// nie zgaduje; od tego jest kontrola steelman i wykrywanie sporu.
	krawedzie := krawedzieWJednejWypowiedzi(okno, wezly)
	if err := a.repozytorium.ZastapGrafDebaty(ctx, okno, turaKod, wezly, krawedzie); err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	if err := a.repozytorium.ZastapUstaleniaDebaty(ctx, okno,
		shared.RoundtableAnalysisKindArgumentMining, turaKod, ustalenia); err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}

	graf, err := a.zlozGraf(ctx, okno, turaKod, false)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}
	return shared.RoundtableAnalysisRunResponse{
		Findings: ustaleniaKontraktu(ustalenia), Graph: &graf,
	}, nil
}

// sklasyfikujAktyMowy nadaje wypowiedziom rodzaj aktu mowy i znakuje nimi zapis.
func (a *adapterDebaty) sklasyfikujAktyMowy(ctx context.Context, okno, turaKod, zapis string,
	wypowiedzi []dane.WypowiedzDebaty, kanalZadania *string) (shared.RoundtableAnalysisRunResponse, error) {

	kanal, err := a.kanalAnalizy(ctx, okno, kanalZadania)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}
	polecenie := "Sklasyfikuj każdą wypowiedź poniższego zapisu debaty. " +
		"W osobnym wierszu podaj kod wypowiedzi w nawiasie kwadratowym, a po nim " +
		"jedno słowo: claim, argument, counterArgument, question, concession albo " +
		"rebuttal.\n\n" + zapis
	odpowiedz, err := a.wywolajModelDebaty(ctx, okno, kanal, "", polecenie)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}

	wlasciciele := wlascicieleWypowiedziDebaty(wypowiedzi)
	ustalenia := make([]dane.UstalenieDebaty, 0, len(wypowiedzi))
	for _, pozycja := range wierszeOdpowiedzi(odpowiedz) {
		kodWypowiedzi, ogon := oddzielKodWypowiedzi(pozycja, wlasciciele)
		if kodWypowiedzi == "" {
			continue
		}
		akt := rozpoznajAktMowy(ogon)
		if akt != "" {
			// Znakowanie idzie do zapisu wypowiedzi, nie tylko do ustalenia:
			// Debate Panel pokazuje etykietę przy wypowiedzi, a nie w osobnym
			// wykazie ustaleń.
			_ = a.repozytorium.OznaczWypowiedz(ctx, kodWypowiedzi, akt, -1)
		}
		ustalenia = append(ustalenia, dane.UstalenieDebaty{
			Kod: nowyIdentyfikator(przedrostekUstalenia), Okno: okno,
			Rodzaj: shared.RoundtableAnalysisKindSpeechAct, Tura: turaKod,
			Wypowiedz: kodWypowiedzi, Uczestnik: wlasciciele[kodWypowiedzi],
			Tresc: ogon, Pewnosc: -1,
		})
	}
	if len(ustalenia) == 0 {
		return shared.RoundtableAnalysisRunResponse{}, odmowaAnalizyBezWyniku(kanal)
	}
	if err := a.repozytorium.ZastapUstaleniaDebaty(ctx, okno,
		shared.RoundtableAnalysisKindSpeechAct, turaKod, ustalenia); err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	return shared.RoundtableAnalysisRunResponse{Findings: ustaleniaKontraktu(ustalenia)}, nil
}

// wykryjBledy stawia oznaczenia błędów logicznych na węzłach grafu.
//
// Zakres wykrywania bierze się z katalogu okna: błąd wyłączony
// (`roundtable.fallacy.catalog.set`) nie jedzie w poleceniu do modelu i nie
// zostaje oznaczony, choćby model go nazwał. Wyłączenie, które nie wyłącza,
// byłoby ustawieniem bez skutku.
func (a *adapterDebaty) wykryjBledy(ctx context.Context, okno, turaKod, zapis string,
	kanalZadania *string) (shared.RoundtableAnalysisRunResponse, error) {

	kanal, err := a.kanalAnalizy(ctx, okno, kanalZadania)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}
	katalog, err := a.repozytorium.KatalogBledowDebaty(ctx, okno)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	wykrywane := make([]dane.DefinicjaBleduDebaty, 0, len(katalog))
	nazwy := make([]string, 0, len(katalog))
	for _, pozycja := range katalog {
		if !pozycja.Wykrywany {
			continue
		}
		wykrywane = append(wykrywane, pozycja)
		nazwy = append(nazwy, pozycja.Kod+" — "+pozycja.Nazwa)
	}
	if len(wykrywane) == 0 {
		return shared.RoundtableAnalysisRunResponse{}, odmowaPustegoKatalogu(okno)
	}

	wezly, err := a.repozytorium.WezlyDebaty(ctx, okno, turaKod, false)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}

	polecenie := "Wskaż w poniższym zapisie debaty chwyty erystyczne i błędy logiczne. " +
		"Wykrywaj wyłącznie te z wykazu:\n" + strings.Join(nazwy, "\n") +
		"\n\nKażde rozpoznanie podaj w osobnym wierszu: kod wypowiedzi w nawiasie " +
		"kwadratowym, kod błędu, uzasadnienie.\n\n" + zapis
	odpowiedz, err := a.wywolajModelDebaty(ctx, okno, kanal, "", polecenie)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}

	poKodzie := make(map[string]dane.DefinicjaBleduDebaty, len(wykrywane))
	for _, pozycja := range wykrywane {
		poKodzie[strings.ToLower(pozycja.Kod)] = pozycja
	}
	wezelWypowiedzi := make(map[string]string, len(wezly))
	for _, wezel := range wezly {
		if wezel.Wypowiedz != "" {
			wezelWypowiedzi[wezel.Wypowiedz] = wezel.Kod
		}
	}

	ustalenia := make([]dane.UstalenieDebaty, 0, 8)
	for _, pozycja := range wierszeOdpowiedzi(odpowiedz) {
		kodWypowiedzi, ogon := oddzielKodWypowiedzi(pozycja, nil)
		definicja, jest := rozpoznajBlad(ogon, poKodzie)
		if !jest {
			continue
		}
		kodWezla := wezelWypowiedzi[kodWypowiedzi]
		if kodWezla != "" {
			_ = a.repozytorium.ZapiszOznaczenieBleduDebaty(ctx, dane.OznaczenieBleduDebaty{
				Kod: nowyIdentyfikator(przedrostekOznaczenia), Okno: okno, Wezel: kodWezla,
				KodBledu: definicja.Kod, Nazwa: definicja.Nazwa, Uzasadnienie: ogon, Pewnosc: -1,
			})
		}
		ustalenia = append(ustalenia, dane.UstalenieDebaty{
			Kod: nowyIdentyfikator(przedrostekUstalenia), Okno: okno,
			Rodzaj: shared.RoundtableAnalysisKindFallacy, Tura: turaKod,
			Wypowiedz: kodWypowiedzi, Wezel: kodWezla,
			Tresc: definicja.Nazwa + ": " + ogon, Pewnosc: -1,
		})
	}
	if len(ustalenia) == 0 {
		// Brak rozpoznań jest wynikiem, nie usterką: debata bez chwytów
		// erystycznych jest debatą poprawną. Ustalenie zbiorcze mówi to wprost,
		// zamiast oddawać pustkę nie do odróżnienia od analizy, która nie ruszyła.
		ustalenia = append(ustalenia, dane.UstalenieDebaty{
			Kod: nowyIdentyfikator(przedrostekUstalenia), Okno: okno,
			Rodzaj: shared.RoundtableAnalysisKindFallacy, Tura: turaKod,
			Tresc:   "W zapisie nie rozpoznano żadnego z " + itoa(len(wykrywane)) + " wykrywanych błędów.",
			Pewnosc: -1,
		})
	}
	if err := a.repozytorium.ZastapUstaleniaDebaty(ctx, okno,
		shared.RoundtableAnalysisKindFallacy, turaKod, ustalenia); err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	graf, err := a.zlozGraf(ctx, okno, turaKod, false)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}
	return shared.RoundtableAnalysisRunResponse{
		Findings: ustaleniaKontraktu(ustalenia), Graph: &graf,
	}, nil
}

// zweryfikujFaktycznosc wypełnia rejestr dowodów: twierdzenie, źródło, adres.
func (a *adapterDebaty) zweryfikujFaktycznosc(ctx context.Context, okno, turaKod, zapis string,
	wypowiedzi []dane.WypowiedzDebaty, kanalZadania *string) (shared.RoundtableAnalysisRunResponse, error) {

	kanal, err := a.kanalAnalizy(ctx, okno, kanalZadania)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}
	polecenie := "Wypisz twierdzenia faktograficzne z poniższego zapisu debaty. " +
		"Każde w osobnym wierszu: kod wypowiedzi w nawiasie kwadratowym, treść " +
		"twierdzenia, a po znaku | źródło, na które uczestnik się powołał. " +
		"Twierdzenie bez źródła zostaw bez znaku |.\n\n" + zapis
	odpowiedz, err := a.wywolajModelDebaty(ctx, okno, kanal, "", polecenie)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}

	wlasciciele := wlascicieleWypowiedziDebaty(wypowiedzi)
	dowody := make([]dane.DowodDebaty, 0, 8)
	ustalenia := make([]dane.UstalenieDebaty, 0, 8)
	for _, pozycja := range wierszeOdpowiedzi(odpowiedz) {
		kodWypowiedzi, ogon := oddzielKodWypowiedzi(pozycja, wlasciciele)
		if strings.TrimSpace(ogon) == "" {
			continue
		}
		twierdzenie, zrodlo := ogon, ""
		if kreska := strings.Index(ogon, "|"); kreska >= 0 {
			twierdzenie = strings.TrimSpace(ogon[:kreska])
			zrodlo = strings.TrimSpace(ogon[kreska+1:])
		}
		if twierdzenie == "" {
			continue
		}
		dowod := dane.DowodDebaty{
			Kod: nowyIdentyfikator(przedrostekDowodu), Okno: okno, Tura: turaKod,
			Wypowiedz: kodWypowiedzi, Twierdzenie: twierdzenie, Poparte: zrodlo != "",
		}
		if zrodlo != "" {
			dowod.Zrodlo = wskaznikTekstu(zrodlo)
			if adres := adresZeZrodla(zrodlo); adres != "" {
				dowod.Adres = wskaznikTekstu(adres)
			}
		}
		dowody = append(dowody, dowod)
		ustalenia = append(ustalenia, dane.UstalenieDebaty{
			Kod: nowyIdentyfikator(przedrostekUstalenia), Okno: okno,
			Rodzaj: shared.RoundtableAnalysisKindFactCheck, Tura: turaKod,
			Wypowiedz: kodWypowiedzi, Uczestnik: wlasciciele[kodWypowiedzi],
			Tresc: twierdzenie, Pewnosc: -1,
		})
	}
	if len(dowody) == 0 {
		return shared.RoundtableAnalysisRunResponse{}, odmowaAnalizyBezWyniku(kanal)
	}
	if err := a.repozytorium.ZastapDowodyDebaty(ctx, okno, turaKod, dowody); err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	if err := a.repozytorium.ZastapUstaleniaDebaty(ctx, okno,
		shared.RoundtableAnalysisKindFactCheck, turaKod, ustalenia); err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	return shared.RoundtableAnalysisRunResponse{Findings: ustaleniaKontraktu(ustalenia)}, nil
}

// analizaOpisowa obsługuje kontrolę steelman i sygnalizację tonu: model opisuje
// zapis, a rdzeń utrwala jego ustalenia bez dopisywania czegokolwiek od siebie.
func (a *adapterDebaty) analizaOpisowa(ctx context.Context, okno, turaKod, rodzaj, zapis string,
	wypowiedzi []dane.WypowiedzDebaty, kanalZadania *string) (shared.RoundtableAnalysisRunResponse, error) {

	kanal, err := a.kanalAnalizy(ctx, okno, kanalZadania)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}
	polecenie := "Oceń, czy kontrargumenty w poniższym zapisie debaty mierzą się " +
		"z najmocniejszą wersją podważanej tezy, czy z wersją osłabioną. " +
		"Każde ustalenie w osobnym wierszu, poprzedzone kodem wypowiedzi " +
		"w nawiasie kwadratowym.\n\n"
	if rodzaj == shared.RoundtableAnalysisKindTone {
		polecenie = "Wskaż w poniższym zapisie debaty przejawy agresji, odejścia od " +
			"tematu i dominacji jednego uczestnika. Każde ustalenie w osobnym " +
			"wierszu, poprzedzone kodem wypowiedzi w nawiasie kwadratowym.\n\n"
	}
	odpowiedz, err := a.wywolajModelDebaty(ctx, okno, kanal, "", polecenie+zapis)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}

	wlasciciele := wlascicieleWypowiedziDebaty(wypowiedzi)
	ustalenia := make([]dane.UstalenieDebaty, 0, 8)
	for _, pozycja := range wierszeOdpowiedzi(odpowiedz) {
		kodWypowiedzi, ogon := oddzielKodWypowiedzi(pozycja, wlasciciele)
		if strings.TrimSpace(ogon) == "" {
			continue
		}
		ustalenia = append(ustalenia, dane.UstalenieDebaty{
			Kod: nowyIdentyfikator(przedrostekUstalenia), Okno: okno, Rodzaj: rodzaj,
			Tura: turaKod, Wypowiedz: kodWypowiedzi, Uczestnik: wlasciciele[kodWypowiedzi],
			Tresc: ogon, Pewnosc: -1,
		})
	}
	if len(ustalenia) == 0 {
		return shared.RoundtableAnalysisRunResponse{}, odmowaAnalizyBezWyniku(kanal)
	}
	if err := a.repozytorium.ZastapUstaleniaDebaty(ctx, okno, rodzaj, turaKod, ustalenia); err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	return shared.RoundtableAnalysisRunResponse{Findings: ustaleniaKontraktu(ustalenia)}, nil
}

// scalPowtorzenia łączy węzły o zbliżonej treści w jeden, licząc poparcie.
//
// Model nie bierze w tym udziału. Podobieństwo treści jest miarą, nie sądem —
// liczy je rdzeń, a wywołanie kanału po to samo kosztowałoby tyle, co
// wypowiedź uczestnika, i dawałoby wynik niepowtarzalny między przebiegami.
func (a *adapterDebaty) scalPowtorzenia(ctx context.Context,
	okno, turaKod string) (shared.RoundtableAnalysisRunResponse, error) {

	wezly, err := a.repozytorium.WezlyDebaty(ctx, okno, turaKod, false)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	if len(wezly) == 0 {
		return shared.RoundtableAnalysisRunResponse{}, odmowaScaleniaBezGrafu(okno)
	}

	scalone := make([]dane.WezelDebaty, 0, len(wezly))
	ustalenia := make([]dane.UstalenieDebaty, 0, 8)
	pochlonieta := make([]bool, len(wezly))
	for i := range wezly {
		if pochlonieta[i] {
			continue
		}
		wezel := wezly[i]
		mowcy := map[string]struct{}{wezel.Uczestnik: {}}
		for j := i + 1; j < len(wezly); j++ {
			if pochlonieta[j] {
				continue
			}
			if podobienstwoTekstow(wezel.Tresc, wezly[j].Tresc) < progScaleniaWezlow {
				continue
			}
			pochlonieta[j] = true
			mowcy[wezly[j].Uczestnik] = struct{}{}
			ustalenia = append(ustalenia, dane.UstalenieDebaty{
				Kod: nowyIdentyfikator(przedrostekUstalenia), Okno: okno,
				Rodzaj: shared.RoundtableAnalysisKindDedup, Tura: turaKod,
				Wezel: wezel.Kod, Uczestnik: wezly[j].Uczestnik,
				Tresc:   "Argument powtórzony: „" + wezly[j].Tresc + "” scalono z węzłem „" + wezel.Tresc + "”.",
				Pewnosc: podobienstwoTekstow(wezel.Tresc, wezly[j].Tresc),
			})
		}
		wezel.Poparcie = len(mowcy)
		scalone = append(scalone, wezel)
	}
	if len(ustalenia) == 0 {
		ustalenia = append(ustalenia, dane.UstalenieDebaty{
			Kod: nowyIdentyfikator(przedrostekUstalenia), Okno: okno,
			Rodzaj: shared.RoundtableAnalysisKindDedup, Tura: turaKod,
			Tresc:   "Żaden z " + itoa(len(wezly)) + " argumentów nie powtarza innego.",
			Pewnosc: -1,
		})
	}

	krawedzie, err := a.repozytorium.KrawedzieDebaty(ctx, okno)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	if err := a.repozytorium.ZastapGrafDebaty(ctx, okno, turaKod, scalone,
		krawedzieMiedzyWezlami(krawedzie, scalone)); err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	if err := a.repozytorium.ZastapUstaleniaDebaty(ctx, okno,
		shared.RoundtableAnalysisKindDedup, turaKod, ustalenia); err != nil {
		return shared.RoundtableAnalysisRunResponse{}, bladDebaty(err)
	}
	graf, err := a.zlozGraf(ctx, okno, turaKod, false)
	if err != nil {
		return shared.RoundtableAnalysisRunResponse{}, err
	}
	return shared.RoundtableAnalysisRunResponse{
		Findings: ustaleniaKontraktu(ustalenia), Graph: &graf,
	}, nil
}

// Dowody oddaje rejestr dowodów i cytowań.
func (a *adapterDebaty) Dowody(ctx context.Context,
	z shared.RoundtableEvidenceListRequest) (shared.RoundtableEvidenceListResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableEvidenceListResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	tylkoNiepoparte := z.UnsupportedOnly != nil && *z.UnsupportedOnly
	dowody, err := a.repozytorium.DowodyDebaty(ctx, okno,
		strings.TrimSpace(wartoscTekstu(z.TurnId)), tylkoNiepoparte)
	if err != nil {
		return shared.RoundtableEvidenceListResponse{}, bladDebaty(err)
	}
	wykaz := make([]shared.RoundtableEvidence, 0, len(dowody))
	for _, dowod := range dowody {
		wykaz = append(wykaz, shared.RoundtableEvidence{
			Id: dowod.Kod, WindowId: dowod.Okno, StatementId: dowod.Wypowiedz,
			Claim: dowod.Twierdzenie, Source: dowod.Zrodlo, Uri: dowod.Adres,
			Supported: dowod.Poparte,
		})
	}
	return shared.RoundtableEvidenceListResponse{Evidence: wykaz}, nil
}

// KatalogBledow oddaje katalog wraz z zakresem wykrywania w oknie.
func (a *adapterDebaty) KatalogBledow(ctx context.Context,
	z shared.RoundtableFallacyCatalogGetRequest) (shared.RoundtableFallacyCatalogGetResponse, error) {

	katalog, err := a.repozytorium.KatalogBledowDebaty(ctx,
		strings.TrimSpace(wartoscTekstu(z.WindowId)))
	if err != nil {
		return shared.RoundtableFallacyCatalogGetResponse{}, bladDebaty(err)
	}
	return shared.RoundtableFallacyCatalogGetResponse{Definitions: katalogKontraktu(katalog)}, nil
}

// UstawKatalogBledow zawęża zakres wykrywania dla okna.
//
// Kod spoza katalogu jest odmową, a nie wpisem cichym: wykaz przyjęty
// z literówką dawałby okno, które „wykrywa” błąd o nazwie, jakiej nie ma.
func (a *adapterDebaty) UstawKatalogBledow(ctx context.Context,
	z shared.RoundtableFallacyCatalogSetRequest) (shared.RoundtableFallacyCatalogSetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableFallacyCatalogSetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	katalog, err := a.repozytorium.KatalogBledowDebaty(ctx, "")
	if err != nil {
		return shared.RoundtableFallacyCatalogSetResponse{}, bladDebaty(err)
	}
	znane := make(map[string]struct{}, len(katalog))
	for _, pozycja := range katalog {
		znane[pozycja.Kod] = struct{}{}
	}
	kody := make([]string, 0, len(z.Codes))
	for _, kod := range z.Codes {
		przyciety := strings.TrimSpace(kod)
		if przyciety == "" {
			continue
		}
		if _, jest := znane[przyciety]; !jest {
			return shared.RoundtableFallacyCatalogSetResponse{}, bladNieznanegoBledu(przyciety)
		}
		kody = append(kody, przyciety)
	}
	if err := a.repozytorium.UstawKatalogBledowDebaty(ctx, okno, kody); err != nil {
		return shared.RoundtableFallacyCatalogSetResponse{}, bladDebaty(err)
	}
	po, err := a.repozytorium.KatalogBledowDebaty(ctx, okno)
	if err != nil {
		return shared.RoundtableFallacyCatalogSetResponse{}, bladDebaty(err)
	}
	return shared.RoundtableFallacyCatalogSetResponse{Definitions: katalogKontraktu(po)}, nil
}

// katalogKontraktu przekłada katalog błędów na byty kontraktu.
func katalogKontraktu(katalog []dane.DefinicjaBleduDebaty) []shared.RoundtableFallacyDefinition {
	wykaz := make([]shared.RoundtableFallacyDefinition, 0, len(katalog))
	for _, pozycja := range katalog {
		wykaz = append(wykaz, shared.RoundtableFallacyDefinition{
			Code: pozycja.Kod, Name: pozycja.Nazwa, Description: pozycja.Opis,
			Enabled: pozycja.Wykrywany,
		})
	}
	return wykaz
}

// ustaleniaKontraktu przekłada ustalenia analizy na byty kontraktu.
func ustaleniaKontraktu(ustalenia []dane.UstalenieDebaty) []shared.RoundtableAnalysisFinding {
	wykaz := make([]shared.RoundtableAnalysisFinding, 0, len(ustalenia))
	for _, ustalenie := range ustalenia {
		pozycja := shared.RoundtableAnalysisFinding{
			Id: ustalenie.Kod, WindowId: ustalenie.Okno,
			Kind: shared.RoundtableAnalysisKind(ustalenie.Rodzaj),
			Text: ustalenie.Tresc, CreatedAt: chwilaBazy(ustalenie.Utworzono),
		}
		if ustalenie.Wypowiedz != "" {
			kod := ustalenie.Wypowiedz
			pozycja.StatementId = &kod
		}
		if ustalenie.Wezel != "" {
			kod := ustalenie.Wezel
			pozycja.NodeId = &kod
		}
		if ustalenie.Uczestnik != "" {
			kod := ustalenie.Uczestnik
			pozycja.ParticipantId = &kod
		}
		if ustalenie.Pewnosc >= 0 {
			pewnosc := ustalenie.Pewnosc
			pozycja.Confidence = &pewnosc
		}
		wykaz = append(wykaz, pozycja)
	}
	return wykaz
}

// wlascicieleWypowiedziDebaty odwzorowuje kod wypowiedzi na kod jej autora.
func wlascicieleWypowiedziDebaty(wypowiedzi []dane.WypowiedzDebaty) map[string]string {
	wlasciciele := make(map[string]string, len(wypowiedzi))
	for _, wypowiedz := range wypowiedzi {
		wlasciciele[wypowiedz.Kod] = wypowiedz.Uczestnik
	}
	return wlasciciele
}

// oddzielKodWypowiedzi zdejmuje z wiersza wiodący kod w nawiasie kwadratowym.
//
// Kod spoza wykazu jest odrzucany, a wiersz zostaje w całości: model, który
// nawiasu użył do czegoś innego, nie ma prawa przypiąć ustalenia do wypowiedzi,
// której nie ma. Wykaz pusty przepuszcza każdy kod — woła się tak wtedy, gdy
// wołający nie ma wykazu do porównania.
func oddzielKodWypowiedzi(wiersz string, znane map[string]string) (string, string) {
	przyciety := strings.TrimSpace(wiersz)
	if !strings.HasPrefix(przyciety, "[") {
		return "", przyciety
	}
	koniec := strings.Index(przyciety, "]")
	if koniec < 0 {
		return "", przyciety
	}
	kod := strings.TrimSpace(przyciety[1:koniec])
	ogon := strings.TrimSpace(przyciety[koniec+1:])
	ogon = strings.TrimLeft(ogon, ":- \t")
	if znane != nil {
		if _, jest := znane[kod]; !jest {
			return "", przyciety
		}
	}
	return kod, strings.TrimSpace(ogon)
}

// rozpoznajAktMowy szuka w tekście nazwy aktu mowy z kontraktu.
func rozpoznajAktMowy(tekst string) string {
	male := strings.ToLower(tekst)
	// Kolejność ma znaczenie: „counterArgument” zawiera w sobie „argument”,
	// więc dłuższa nazwa musi być sprawdzona pierwsza.
	for _, akt := range []string{
		shared.RoundtableSpeechActCounterArgument,
		shared.RoundtableSpeechActArgument,
		shared.RoundtableSpeechActConcession,
		shared.RoundtableSpeechActRebuttal,
		shared.RoundtableSpeechActQuestion,
		shared.RoundtableSpeechActClaim,
	} {
		if strings.Contains(male, strings.ToLower(akt)) {
			return akt
		}
	}
	return ""
}

// rozpoznajBlad szuka w tekście kodu błędu z wykazu wykrywanych.
func rozpoznajBlad(tekst string,
	poKodzie map[string]dane.DefinicjaBleduDebaty) (dane.DefinicjaBleduDebaty, bool) {

	male := strings.ToLower(tekst)
	for kod, definicja := range poKodzie {
		if strings.Contains(male, kod) || strings.Contains(male, strings.ToLower(definicja.Nazwa)) {
			return definicja, true
		}
	}
	return dane.DefinicjaBleduDebaty{}, false
}

// adresZeZrodla wyciąga adres sieciowy z opisu źródła, gdy w nim jest.
func adresZeZrodla(zrodlo string) string {
	for _, slowo := range strings.Fields(zrodlo) {
		if strings.HasPrefix(slowo, "http://") || strings.HasPrefix(slowo, "https://") {
			return strings.Trim(slowo, ".,;)")
		}
	}
	return ""
}

// krawedzieWJednejWypowiedzi buduje relację wsparcia wewnątrz jednej wypowiedzi:
// pierwszy węzeł jest tezą, każdy następny wspiera poprzedni.
func krawedzieWJednejWypowiedzi(okno string, wezly []dane.WezelDebaty) []dane.KrawedzDebaty {
	krawedzie := make([]dane.KrawedzDebaty, 0, len(wezly))
	for i := 1; i < len(wezly); i++ {
		if wezly[i].Wypowiedz == "" || wezly[i].Wypowiedz != wezly[i-1].Wypowiedz {
			continue
		}
		krawedzie = append(krawedzie, dane.KrawedzDebaty{
			Kod: nowyIdentyfikator(przedrostekKrawedzi), Okno: okno,
			WezelZrodlowy: wezly[i].Kod, WezelWskazany: wezly[i-1].Kod,
			Relacja: shared.RoundtableArgumentRelationSupports, Pewnosc: -1,
		})
	}
	return krawedzie
}

// krawedzieMiedzyWezlami odsiewa krawędzie wskazujące węzły, których po
// scaleniu już nie ma. Krawędź do węzła nieistniejącego rysowałaby w grafie
// strzałkę donikąd.
func krawedzieMiedzyWezlami(krawedzie []dane.KrawedzDebaty,
	wezly []dane.WezelDebaty) []dane.KrawedzDebaty {

	zywe := make(map[string]struct{}, len(wezly))
	for _, wezel := range wezly {
		zywe[wezel.Kod] = struct{}{}
	}
	pozostale := make([]dane.KrawedzDebaty, 0, len(krawedzie))
	for _, krawedz := range krawedzie {
		_, zrodlo := zywe[krawedz.WezelZrodlowy]
		_, cel := zywe[krawedz.WezelWskazany]
		if zrodlo && cel {
			pozostale = append(pozostale, krawedz)
		}
	}
	return pozostale
}
