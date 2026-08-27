// Plik obsługuje translate.proofread.run, translate.proofread.apply i
// translate.consistency.check: korektę języka regułami wbudowanymi i słownikami
// zewnętrznymi oraz sprawdzenie spójności segmentów i terminów między panelami
// okna.
package core

import (
	"context"
	"sort"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekUstaleniaKorekty znakuje identyfikator ustalenia korekty: kod
// z tym przedrostkiem odróżnia zapis korekty od innych bytów bazy, gdy
// proofread.apply adresuje go w żądaniu.
const przedrostekUstaleniaKorekty = "kor-"

// dlugieZdanie jest progiem czytelności: zdanie dłuższe od tylu znaków czyta
// się z wysiłkiem w każdym języku, na który ten moduł tłumaczy.
const dlugieZdanie = 120

// SprawdzKorekte obsługuje `translate.proofread.run`. Zapisuje ustalenia
// w bazie, bo `proofread.apply` adresuje je identyfikatorem w osobnym
// wywołaniu.
func (a *adapterTlumaczenia) SprawdzKorekte(ctx context.Context,
	z shared.TranslateProofreadRunRequest) (shared.TranslateProofreadRunResponse, error) {

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateProofreadRunResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	tresc, err := trescPanelu(panel)
	if err != nil {
		return shared.TranslateProofreadRunResponse{}, err
	}

	zadane := map[shared.ProofreadCheckKind]struct{}{}
	for _, rodzaj := range z.Checks {
		zadane[rodzaj] = struct{}{}
	}
	chce := func(rodzaj shared.ProofreadCheckKind) bool {
		if len(zadane) == 0 {
			return true
		}
		_, jest := zadane[rodzaj]
		return jest
	}

	teraz := time.Now().UnixMilli()
	ustalenia := []dane.UstalenieKorekty{}
	dolóz := func(rodzaj shared.ProofreadCheckKind, waga shared.ProofreadSeverity,
		segment, szczegol, propozycja string) {
		if !chce(rodzaj) {
			return
		}
		ustalenia = append(ustalenia, dane.UstalenieKorekty{
			Kod:        nowyIdentyfikator(przedrostekUstaleniaKorekty),
			PanelID:    panel.ID,
			Rodzaj:     string(rodzaj),
			Waga:       string(waga),
			Segment:    wskaznikNapisu(segment),
			Szczegol:   szczegol,
			Propozycja: wskaznikNapisu(propozycja),
			Utworzono:  teraz,
		})
	}

	if strings.Contains(tresc, "  ") {
		dolóz(shared.ProofreadCheckKindTypography, shared.ProofreadSeverityWarning, "",
			"treść ma podwójne odstępy",
			strings.Join(strings.Fields(tresc), " "))
	}
	for _, znak := range []string{" ,", " .", " ;", " :", " !", " ?"} {
		if strings.Contains(tresc, znak) {
			dolóz(shared.ProofreadCheckKindPunctuation, shared.ProofreadSeverityError, "",
				"odstęp przed znakiem interpunkcyjnym „"+strings.TrimSpace(znak)+"”",
				strings.ReplaceAll(tresc, znak, strings.TrimSpace(znak)))
		}
	}
	if strings.Contains(tresc, "...") {
		dolóz(shared.ProofreadCheckKindTypography, shared.ProofreadSeverityHint, "",
			"wielokropek złożony z trzech kropek zamiast znaku wielokropka",
			strings.ReplaceAll(tresc, "...", "…"))
	}
	if strings.Contains(tresc, "\"") {
		dolóz(shared.ProofreadCheckKindTypography, shared.ProofreadSeverityHint, "",
			"cudzysłów prosty zamiast cudzysłowu drukarskiego rynku docelowego",
			zamienCudzyslowy(tresc))
	}
	for _, zdanie := range podzielNaZdania(tresc) {
		if liczbaZnakow(zdanie) > dlugieZdanie {
			dolóz(shared.ProofreadCheckKindReadability, shared.ProofreadSeverityHint, zdanie,
				"zdanie dłuższe niż 120 znaków — czytelność spada", "")
		}
		if slowoPowtorzone(zdanie) != "" {
			dolóz(shared.ProofreadCheckKindGrammar, shared.ProofreadSeverityWarning, zdanie,
				"słowo „"+slowoPowtorzone(zdanie)+"” powtórzone bezpośrednio po sobie", "")
		}
	}

	// Silniki zewnętrzne uzupełniają reguły wbudowane zakresem gramatyki
	// i pisowni, nie dublują ich.
	zewnetrzneUstalenia, err := a.ustaleniaSilnikow(ctx, panel.Jezyk, tresc)
	if err != nil {
		return shared.TranslateProofreadRunResponse{}, err
	}
	for _, ustalenie := range zewnetrzneUstalenia {
		dolóz(ustalenie.rodzaj, ustalenie.waga, ustalenie.segment,
			ustalenie.szczegol, ustalenie.propozycja)
	}

	if err := a.repozytorium.ZapiszUstaleniaKorekty(ctx, panel.ID, ustalenia); err != nil {
		return shared.TranslateProofreadRunResponse{}, bladTlumaczenia(err)
	}
	zapisane, err := a.repozytorium.UstaleniaKorekty(ctx, panel.ID)
	if err != nil {
		return shared.TranslateProofreadRunResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslateProofreadRunResponse{
		PanelId:     panel.Kod,
		Findings:    zlozUstaleniaKorekty(zapisane),
		Readability: miaryCzytelnosci(tresc),
	}, nil
}

// zamienCudzyslowy zamienia cudzysłowy proste na drukarskie parami: pierwszy
// otwierający, drugi zamykający.
func zamienCudzyslowy(tekst string) string {
	var b strings.Builder
	otwarty := false
	for _, znak := range tekst {
		if znak != '"' {
			b.WriteRune(znak)
			continue
		}
		if otwarty {
			b.WriteString("”")
		} else {
			b.WriteString("„")
		}
		otwarty = !otwarty
	}
	return b.String()
}

// slowoPowtorzone wskazuje pierwsze słowo powtórzone bezpośrednio po sobie —
// usterka, którą oko przeskakuje, a reguła wyłapuje bez wątpliwości.
func slowoPowtorzone(zdanie string) string {
	slowa := strings.Fields(strings.ToLower(zdanie))
	for i := 1; i < len(slowa); i++ {
		if slowa[i] == slowa[i-1] && liczbaZnakow(slowa[i]) >= najmniejszaDlugoscKandydata {
			return slowa[i]
		}
	}
	return ""
}

// miaryCzytelnosci liczy trzy miary, których żadna nie wymaga słownika:
// średnią długość zdania, średnią długość słowa i liczbę zdań.
func miaryCzytelnosci(tresc string) []shared.ReadabilityScore {
	zdania := podzielNaZdania(tresc)
	slowa := strings.Fields(tresc)
	if len(zdania) == 0 || len(slowa) == 0 {
		return []shared.ReadabilityScore{}
	}
	znakiSlow := 0
	for _, slowo := range slowa {
		znakiSlow += liczbaZnakow(slowo)
	}
	sredniaZdania := float64(len(slowa)) / float64(len(zdania))
	sredniaSlowa := float64(znakiSlow) / float64(len(slowa))
	wyjasnienie := "zdania w normie"
	if sredniaZdania > 25 {
		wyjasnienie = "zdania długie — tekst czyta się z wysiłkiem"
	}
	return []shared.ReadabilityScore{
		{Metric: "slow-na-zdanie", Value: sredniaZdania, Interpretation: &wyjasnienie},
		{Metric: "znakow-na-slowo", Value: sredniaSlowa},
		{Metric: "zdan", Value: float64(len(zdania))},
	}
}

// zlozUstaleniaKorekty przekłada wiersze ustaleń na byty kontraktu. Ustalenia
// rozstrzygnięte nie wychodzą — Operator odpowiedział na nie raz.
func zlozUstaleniaKorekty(ustalenia []dane.UstalenieKorekty) []shared.ProofreadFinding {
	wykaz := []shared.ProofreadFinding{}
	for _, ustalenie := range ustalenia {
		if ustalenie.Zastosowano != nil || ustalenie.Odrzucono != nil {
			continue
		}
		wykaz = append(wykaz, shared.ProofreadFinding{
			Id:         ustalenie.Kod,
			Kind:       shared.ProofreadCheckKind(ustalenie.Rodzaj),
			Severity:   wagaOdmowyKontroli(ustalenie.Waga),
			Segment:    ustalenie.Segment,
			Detail:     ustalenie.Szczegol,
			Suggestion: ustalenie.Propozycja,
			CreatedAt:  ustalenie.Utworzono,
		})
	}
	return wykaz
}

// ZastosujKorekte obsługuje `translate.proofread.apply`. Wstawia propozycje
// wskazanych ustaleń w treść panelu — albo, gdy Operator zażądał oddalenia,
// znakuje je jako odrzucone bez ruszania treści.
func (a *adapterTlumaczenia) ZastosujKorekte(ctx context.Context,
	z shared.TranslateProofreadApplyRequest) (shared.TranslateProofreadApplyResponse, error) {

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateProofreadApplyResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	if len(z.FindingIds) == 0 {
		return shared.TranslateProofreadApplyResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.proofread.apply bez wskazania ustaleń")
	}
	odrzuca := z.Dismiss != nil && *z.Dismiss

	tresc, err := trescPanelu(panel)
	if err != nil {
		return shared.TranslateProofreadApplyResponse{}, err
	}

	zastosowane := 0
	for _, kod := range z.FindingIds {
		ustalenie, err := a.repozytorium.UstalenieKorektyPoKodzie(ctx, kod)
		if err != nil {
			return shared.TranslateProofreadApplyResponse{}, bladWskazaniaTlumaczenia(
				"nie ma ustalenia korekty o identyfikatorze " + kod)
		}
		if ustalenie.PanelID != panel.ID {
			return shared.TranslateProofreadApplyResponse{}, bladWskazaniaTlumaczenia(
				"ustalenie " + kod + " dotyczy innego panelu niż wskazany")
		}
		if odrzuca {
			if err := a.repozytorium.RozstrzygnijUstalenieKorekty(ctx, kod, true); err != nil {
				return shared.TranslateProofreadApplyResponse{}, bladTlumaczenia(err)
			}
			zastosowane++
			continue
		}
		if ustalenie.Propozycja == nil || strings.TrimSpace(*ustalenie.Propozycja) == "" {
			return shared.TranslateProofreadApplyResponse{}, bladWskazaniaTlumaczenia(
				"ustalenie " + kod + " nie niesie propozycji poprawki — nie ma czego wstawić")
		}
		// Propozycja to cała treść po poprawce, nie łata; ustalenia
		// działają na treści już zmienionej.
		tresc = *ustalenie.Propozycja
		if err := a.repozytorium.RozstrzygnijUstalenieKorekty(ctx, kod, false); err != nil {
			return shared.TranslateProofreadApplyResponse{}, bladTlumaczenia(err)
		}
		zastosowane++
	}

	if !odrzuca {
		zmieniony, err := a.repozytorium.UstawTlumaczenie(ctx, panel.Kod, &tresc, nil)
		if err != nil {
			return shared.TranslateProofreadApplyResponse{}, bladTlumaczenia(err)
		}
		a.rozglosZmianePanelu(shared.ChangeKindUpdated, zmieniony)
		return shared.TranslateProofreadApplyResponse{
			Panel:        zlozPanelTlumaczenia(zmieniony),
			AppliedCount: zastosowane,
		}, nil
	}

	poZmianie, err := a.repozytorium.Panel(ctx, panel.Kod)
	if err != nil {
		return shared.TranslateProofreadApplyResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateProofreadApplyResponse{
		Panel:        zlozPanelTlumaczenia(poZmianie),
		AppliedCount: zastosowane,
	}, nil
}

// SprawdzSpojnosc obsługuje `translate.consistency.check`: wykrywa segmenty
// przełożone niejednolicie w obrębie języka oraz terminy odbiegające od
// słownika Operatora.
func (a *adapterTlumaczenia) SprawdzSpojnosc(ctx context.Context,
	z shared.TranslateConsistencyCheckRequest) (shared.TranslateConsistencyCheckResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateConsistencyCheckResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateConsistencyCheckResponse{}, bladTlumaczenia(err)
	}
	wskazany := napisZeWskaznika(z.PanelId)

	ustalenia := []shared.ConsistencyFinding{}

	// Pary bierzemy z pamięci tłumaczeń, bo tam leży treść zatwierdzona,
	// nie bieżący szkic panelu.
	for _, panel := range panele {
		if wskazany != "" && panel.Kod != wskazany {
			continue
		}
		pary, _, err := a.repozytorium.WpisyPamieci(ctx,
			dane.FiltrPamieciTlumaczen{Jezyk: panel.Jezyk})
		if err != nil {
			return shared.TranslateConsistencyCheckResponse{}, bladTlumaczenia(err)
		}
		warianty := map[string]map[string]struct{}{}
		for _, para := range pary {
			if warianty[para.SegmentZrodlowy] == nil {
				warianty[para.SegmentZrodlowy] = map[string]struct{}{}
			}
			warianty[para.SegmentZrodlowy][para.SegmentDocelowy] = struct{}{}
		}
		for zrodlo, zbior := range warianty {
			if len(zbior) < 2 {
				continue
			}
			ustalenia = append(ustalenia, shared.ConsistencyFinding{
				Kind:       shared.ConsistencyFindingKindSegment,
				SourceText: zrodlo,
				Variants:   uporzadkowaneWarianty(zbior),
				PanelIds:   []string{panel.Kod},
			})
		}
	}

	// Niespójność terminu: treść panelu odbiega od odpowiednika ustalonego
	// w słowniku Operatora.
	terminy, err := a.repozytorium.Terminy(ctx)
	if err != nil {
		return shared.TranslateConsistencyCheckResponse{}, bladTlumaczenia(err)
	}
	for _, termin := range terminy {
		if termin.Cel == nil || strings.TrimSpace(*termin.Cel) == "" {
			continue
		}
		panelePominiete := []string{}
		for _, panel := range panele {
			if panel.Jezyk != termin.Jezyk || panel.Tresc == nil {
				continue
			}
			if wskazany != "" && panel.Kod != wskazany {
				continue
			}
			if strings.Contains(*panel.Tresc, *termin.Cel) {
				continue
			}
			// Termin źródłowy pozostał w przekładzie zamiast odpowiednika
			// ze słownika.
			if strings.Contains(strings.ToLower(*panel.Tresc), strings.ToLower(termin.Zrodlo)) {
				panelePominiete = append(panelePominiete, panel.Kod)
			}
		}
		if len(panelePominiete) == 0 {
			continue
		}
		ustalenia = append(ustalenia, shared.ConsistencyFinding{
			Kind:       shared.ConsistencyFindingKindTerm,
			SourceText: termin.Zrodlo,
			Variants:   []string{termin.Zrodlo, *termin.Cel},
			PanelIds:   panelePominiete,
		})
	}

	sort.Slice(ustalenia, func(i, j int) bool {
		return ustalenia[i].SourceText < ustalenia[j].SourceText
	})
	return shared.TranslateConsistencyCheckResponse{Findings: ustalenia}, nil
}

// uporzadkowaneWarianty daje wykaz powtarzalny — zbiór w Go nie ma kolejności,
// a odpowiedź komendy ma być ta sama przy tym samym stanie danych.
func uporzadkowaneWarianty(zbior map[string]struct{}) []string {
	wykaz := make([]string, 0, len(zbior))
	for wariant := range zbior {
		wykaz = append(wykaz, wariant)
	}
	sort.Strings(wykaz)
	return wykaz
}
