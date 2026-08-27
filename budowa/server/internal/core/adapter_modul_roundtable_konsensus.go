// Modul Roundtable — stanowisko koncowe redagowane przez Operatora: zapis
// stanowiska, zdanie odrebne, wykaz wersji redakcji oraz przekazanie
// stanowiska do modulu docelowego.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przedrostki bytow stanowiska nadawane przez rdzen przy zapisie kazdego
// nowego wpisu w tym repozytorium.
const (
	przedrostekWersjiStanowiska = "wersja-"
	przedrostekZdaniaOdrebnego  = "odreb-"
	przedrostekPrzekazania      = "przekaz-"
)

// ZapiszStanowisko zapisuje treść stanowiska nadaną przez Operatora wraz
// z zapisem decyzji: kontekstem, rozważanymi wariantami i konsekwencjami.
func (a *adapterDebaty) ZapiszStanowisko(ctx context.Context,
	z shared.RoundtableConsensusSetRequest) (shared.RoundtableConsensusSetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	tresc := strings.TrimSpace(z.Content)
	if okno == "" {
		return shared.RoundtableConsensusSetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if tresc == "" {
		return shared.RoundtableConsensusSetResponse{},
			bladWskazaniaDebaty("stanowisko bez treści — pusta redakcja skasowałaby zapis debaty")
	}

	// Tury wskazane w zadaniu sprawdza sie przed zapisem, zeby stanowisko nie objelo tury z cudzego okna.
	kody := make([]string, 0, len(z.TurnIds))
	for _, kod := range z.TurnIds {
		przyciety := strings.TrimSpace(kod)
		if przyciety == "" {
			continue
		}
		tura, err := a.repozytorium.Tura(ctx, przyciety)
		if err != nil {
			return shared.RoundtableConsensusSetResponse{}, bladNieznanejTury(przyciety, err)
		}
		if tura.Okno != okno {
			return shared.RoundtableConsensusSetResponse{},
				bladWskazaniaDebaty("tura " + przyciety + " nie należy do okna " + okno)
		}
		kody = append(kody, przyciety)
	}

	// Redakcja stanowiska juz istniejacego zachowuje jego identyfikator: zdania odrebne wskazuja je kodem.
	kod := nowyIdentyfikator(przedrostekStanowiska)
	if poprzednie, err := a.repozytorium.Stanowisko(ctx, okno, ""); err == nil {
		kod = poprzednie.Kod
	}

	zaakceptowane := z.Accepted != nil && *z.Accepted
	stanowisko, err := a.repozytorium.RedagujStanowisko(ctx, dane.StanowiskoDebaty{
		Kod: kod, Okno: okno, Tura: "", Tresc: wskaznikTekstu(tresc),
		Zaakceptowane: zaakceptowane, Kontekst: z.Context, Warianty: z.Options,
		Konsekwencje: z.Consequences, Tury: strings.Join(kody, "\n"),
	})
	if err != nil {
		return shared.RoundtableConsensusSetResponse{}, bladDebaty(err)
	}
	// Wersja odklada sie po zapisie, z numerem nadanym przez baze, a nie przez rdzen.
	if err := a.repozytorium.ZapiszWersjeStanowiskaDebaty(ctx, dane.WersjaStanowiskaDebaty{
		Kod: nowyIdentyfikator(przedrostekWersjiStanowiska), Stanowisko: stanowisko.Kod,
		Wersja: stanowisko.Wersja, Tresc: tresc,
	}); err != nil {
		return shared.RoundtableConsensusSetResponse{}, bladDebaty(err)
	}

	wynik, err := a.stanowiskoZDodatkami(ctx, stanowisko)
	if err != nil {
		return shared.RoundtableConsensusSetResponse{}, err
	}
	return shared.RoundtableConsensusSetResponse{Consensus: wynik}, nil
}

// ZapiszZdanieOdrebne utrwala i podpisuje zdanie odrebne uczestnika, ktory nie
// dolaczyl do konsensusu debaty.
func (a *adapterDebaty) ZapiszZdanieOdrebne(ctx context.Context,
	z shared.RoundtableConsensusMinoritySetRequest) (shared.RoundtableConsensusMinoritySetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kodStanowiska := strings.TrimSpace(z.ConsensusId)
	kodUczestnika := strings.TrimSpace(z.ParticipantId)
	tresc := strings.TrimSpace(z.Content)
	if okno == "" {
		return shared.RoundtableConsensusMinoritySetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kodStanowiska == "" {
		return shared.RoundtableConsensusMinoritySetResponse{},
			bladWskazaniaDebaty("zdanie odrębne bez wskazania stanowiska")
	}
	if kodUczestnika == "" {
		return shared.RoundtableConsensusMinoritySetResponse{},
			bladWskazaniaDebaty("zdanie odrębne bez podpisu uczestnika")
	}
	if tresc == "" {
		return shared.RoundtableConsensusMinoritySetResponse{},
			bladWskazaniaDebaty("zdanie odrębne bez treści")
	}

	stanowisko, err := a.repozytorium.StanowiskoPoKodzie(ctx, kodStanowiska)
	if err != nil {
		return shared.RoundtableConsensusMinoritySetResponse{},
			bladNieznanegoStanowiska(kodStanowiska, err)
	}
	if stanowisko.Okno != okno {
		return shared.RoundtableConsensusMinoritySetResponse{},
			bladWskazaniaDebaty("stanowisko " + kodStanowiska + " nie należy do okna " + okno)
	}
	uczestnik, err := a.repozytorium.Uczestnik(ctx, kodUczestnika)
	if err != nil {
		return shared.RoundtableConsensusMinoritySetResponse{},
			bladNieznanegoUczestnika(kodUczestnika, err)
	}
	if uczestnik.Okno != okno {
		return shared.RoundtableConsensusMinoritySetResponse{},
			bladWskazaniaDebaty("uczestnik " + kodUczestnika + " nie należy do okna " + okno)
	}

	zdanie, err := a.repozytorium.ZapiszZdanieOdrebneDebaty(ctx, dane.ZdanieOdrebneDebaty{
		Kod: nowyIdentyfikator(przedrostekZdaniaOdrebnego), Okno: okno,
		Stanowisko: kodStanowiska, Uczestnik: kodUczestnika, Tresc: tresc,
	})
	if err != nil {
		return shared.RoundtableConsensusMinoritySetResponse{}, bladDebaty(err)
	}
	return shared.RoundtableConsensusMinoritySetResponse{
		Minority: zdanieOdrebneKontraktu(zdanie),
	}, nil
}

// WersjeStanowiska oddaje kolejne redakcje wskazanego stanowiska tej debaty
// od najstarszej do najnowszej.
func (a *adapterDebaty) WersjeStanowiska(ctx context.Context,
	z shared.RoundtableConsensusVersionListRequest) (shared.RoundtableConsensusVersionListResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableConsensusVersionListResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	kod := strings.TrimSpace(wartoscTekstu(z.ConsensusId))
	if kod == "" {
		stanowisko, err := a.repozytorium.Stanowisko(ctx, okno, "")
		if err != nil {
			return shared.RoundtableConsensusVersionListResponse{},
				bladNieznanegoStanowiska(okno, err)
		}
		kod = stanowisko.Kod
	}

	wersje, err := a.repozytorium.WersjeStanowiskaDebaty(ctx, kod)
	if err != nil {
		return shared.RoundtableConsensusVersionListResponse{}, bladDebaty(err)
	}
	wykaz := make([]shared.RoundtableConsensusVersion, 0, len(wersje))
	for _, wersja := range wersje {
		wykaz = append(wykaz, shared.RoundtableConsensusVersion{
			Id: wersja.Kod, ConsensusId: wersja.Stanowisko, Version: wersja.Wersja,
			Content: wersja.Tresc, CreatedAt: chwilaBazy(wersja.Utworzono),
		})
	}
	return shared.RoundtableConsensusVersionListResponse{Versions: wykaz}, nil
}

// PrzekazStanowisko wydaje stanowisko jako artefakt dla modulu docelowego.
// Artefakt da sie otworzyc niezaleznie od tego, czy debata jeszcze istnieje;
// transkrypt dolacza sie na zadanie.
func (a *adapterDebaty) PrzekazStanowisko(ctx context.Context,
	z shared.RoundtableConsensusHandoffRequest) (shared.RoundtableConsensusHandoffResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kodStanowiska := strings.TrimSpace(z.ConsensusId)
	modul := strings.TrimSpace(string(z.Target))
	if okno == "" {
		return shared.RoundtableConsensusHandoffResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kodStanowiska == "" {
		return shared.RoundtableConsensusHandoffResponse{},
			bladWskazaniaDebaty("przekazanie bez wskazania stanowiska")
	}
	switch modul {
	case shared.RoundtableHandoffTargetStudio, shared.RoundtableHandoffTargetResearch,
		shared.RoundtableHandoffTargetLibrary, shared.RoundtableHandoffTargetAutomations:
	default:
		return shared.RoundtableConsensusHandoffResponse{},
			bladWskazaniaDebaty("moduł docelowy " + modul + " nie jest modułem znanym kontraktowi")
	}

	stanowisko, err := a.repozytorium.StanowiskoPoKodzie(ctx, kodStanowiska)
	if err != nil {
		return shared.RoundtableConsensusHandoffResponse{},
			bladNieznanegoStanowiska(kodStanowiska, err)
	}
	if stanowisko.Okno != okno {
		return shared.RoundtableConsensusHandoffResponse{},
			bladWskazaniaDebaty("stanowisko " + kodStanowiska + " nie należy do okna " + okno)
	}

	tresc := zapisDecyzjiDebaty(stanowisko)
	if z.IncludeTranscript != nil && *z.IncludeTranscript {
		transkrypt, err := a.transkryptMarkdown(ctx, okno, "", false)
		if err != nil {
			return shared.RoundtableConsensusHandoffResponse{}, err
		}
		tresc += "\n\n---\n\n" + transkrypt
	}
	artefakt, err := a.wydajArtefaktDebaty(ctx, okno, rodzajArtefaktuTranskryptu,
		shared.RoundtableTranscriptFormatMarkdown, []byte(tresc), 0)
	if err != nil {
		return shared.RoundtableConsensusHandoffResponse{}, err
	}

	przekazanie := dane.PrzekazanieDebaty{
		Kod: nowyIdentyfikator(przedrostekPrzekazania), Okno: okno,
		Stanowisko: kodStanowiska, Modul: modul, Artefakt: artefakt.Kod,
	}
	if err := a.repozytorium.ZapiszPrzekazanieDebaty(ctx, przekazanie); err != nil {
		return shared.RoundtableConsensusHandoffResponse{}, bladDebaty(err)
	}
	kodArtefaktu := artefakt.Kod
	return shared.RoundtableConsensusHandoffResponse{
		Handoff: shared.RoundtableHandoff{
			Id: przekazanie.Kod, WindowId: okno, ConsensusId: kodStanowiska,
			Target: shared.RoundtableHandoffTarget(modul), ArtifactId: &kodArtefaktu,
			CreatedAt: chwilaBazy(artefakt.Utworzono),
		},
	}, nil
}

// stanowiskoZDodatkami dopina do stanowiska zdania odrębne oraz punkty zgody
// i sporu — Consensus Panel pokazuje je razem z treścią.
func (a *adapterDebaty) stanowiskoZDodatkami(ctx context.Context,
	stanowisko dane.StanowiskoDebaty) (shared.RoundtableConsensus, error) {

	wynik := stanowiskoKontraktu(stanowisko, nil)

	zdania, err := a.repozytorium.ZdaniaOdrebneDebaty(ctx, stanowisko.Kod)
	if err != nil {
		return shared.RoundtableConsensus{}, bladDebaty(err)
	}
	for _, zdanie := range zdania {
		wynik.Minority = append(wynik.Minority, zdanieOdrebneKontraktu(zdanie))
	}

	zgodnosc, err := a.Zgodnosc(ctx, shared.RoundtableAgreementGetRequest{WindowId: stanowisko.Okno})
	if err != nil {
		return shared.RoundtableConsensus{}, err
	}
	for _, punkt := range zgodnosc.Points {
		if punkt.Agreed {
			wynik.AgreementPoints = append(wynik.AgreementPoints, punkt)
			continue
		}
		wynik.DisputePoints = append(wynik.DisputePoints, punkt)
	}

	// Poparcie wazone liczy sie z wag uczestnikow, ktorzy nie podpisali zdania odrebnego.
	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, stanowisko.Okno)
	if err != nil {
		return shared.RoundtableConsensus{}, bladDebaty(err)
	}
	odrebni := make(map[string]struct{}, len(zdania))
	for _, zdanie := range zdania {
		odrebni[zdanie.Uczestnik] = struct{}{}
	}
	suma, popierajace := 0.0, 0.0
	for _, uczestnik := range uczestnicy {
		suma += uczestnik.Waga
		if _, odrebny := odrebni[uczestnik.Kod]; !odrebny {
			popierajace += uczestnik.Waga
		}
	}
	if suma > 0 {
		poparcie := popierajace / suma
		wynik.Support = &poparcie
	}
	return wynik, nil
}

// zdanieOdrebneKontraktu przeklada zdanie odrebne zapisane w repozytorium na
// byt oddawany w kontrakcie.
func zdanieOdrebneKontraktu(z dane.ZdanieOdrebneDebaty) shared.RoundtableMinorityReport {
	return shared.RoundtableMinorityReport{
		Id: z.Kod, WindowId: z.Okno, ConsensusId: z.Stanowisko, ParticipantId: z.Uczestnik,
		Content: z.Tresc, CreatedAt: chwilaBazy(z.Utworzono),
	}
}

// zapisDecyzjiDebaty składa stanowisko w postać zapisu decyzji: kontekst,
// rozważane warianty, decyzja, konsekwencje. Części pustych nie wypisuje —
// nagłówek bez treści byłby obietnicą, której zapis nie spełnia.
func zapisDecyzjiDebaty(s dane.StanowiskoDebaty) string {
	czesci := make([]string, 0, 5)
	czesci = append(czesci, "# Stanowisko końcowe debaty")
	if s.Kontekst != nil && strings.TrimSpace(*s.Kontekst) != "" {
		czesci = append(czesci, "## Kontekst\n\n"+strings.TrimSpace(*s.Kontekst))
	}
	if s.Warianty != nil && strings.TrimSpace(*s.Warianty) != "" {
		czesci = append(czesci, "## Rozważane warianty\n\n"+strings.TrimSpace(*s.Warianty))
	}
	if s.Tresc != nil && strings.TrimSpace(*s.Tresc) != "" {
		czesci = append(czesci, "## Decyzja\n\n"+strings.TrimSpace(*s.Tresc))
	}
	if s.Konsekwencje != nil && strings.TrimSpace(*s.Konsekwencje) != "" {
		czesci = append(czesci, "## Konsekwencje\n\n"+strings.TrimSpace(*s.Konsekwencje))
	}
	return strings.Join(czesci, "\n\n")
}
