// Plik obsługuje pracę wielu wykonawców nad jednym dokumentem: zajęcie
// i zwolnienie fragmentu, wykaz zajęć, bilans spięć oraz nastawy pętli
// wykonawczej i pracy wielu agentów.
package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// agenciWaznoscDomyslna to czas, po którym zajęcie fragmentu wygasa samo, gdy
// wykonawca nie podał własnego.
const agenciWaznoscDomyslna = 90

// ZajmijFragment obsługuje studio.agents.claim: zajmuje fragment dokumentu dla
// wykonawcy na czas wskazany albo na czas domyślny.
func (a *adapterStudia) ZajmijFragment(ctx context.Context,
	z shared.StudioAgentsClaimRequest) (shared.StudioAgentsClaimResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioAgentsClaimResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioAgentsClaimResponse{}, err
	}
	if z.RangeEnd < z.RangeStart {
		return shared.StudioAgentsClaimResponse{}, bladWskazaniaStudio(
			"zajęcie o zakresie odwróconym: koniec " + strconv.Itoa(z.RangeEnd) +
				" leży przed początkiem " + strconv.Itoa(z.RangeStart))
	}
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{
		AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	})
	if !wykonawca.czyWykonawca() {
		return shared.StudioAgentsClaimResponse{}, bladWskazaniaStudio(
			"zajęcie fragmentu bez wskazania wykonawcy — zajęcie istnieje po to, " +
				"żeby dwóch wykonawców nie pisało po tym samym akapicie, więc musi " +
				"wiedzieć, czyje jest. Operator zajęcia nie potrzebuje: fragment, " +
				"który ma zostać nietknięty, oznacza blokadą (studio.lock.add)")
	}

	// Zajęcia wygasłe schodzą przed sprawdzeniem, inaczej agent ubity w pół
	// pracy trzymałby fragment.
	if _, err := skladnica.PrzemiecZajeciaFragmentow(ctx, dokument.ID); err != nil {
		return shared.StudioAgentsClaimResponse{}, bladStudio(err)
	}
	zajecia, err := skladnica.ZajeciaFragmentow(ctx, dokument.ID)
	if err != nil {
		return shared.StudioAgentsClaimResponse{}, bladStudio(err)
	}
	for _, zajete := range zajecia {
		if zajete.Wygasle {
			continue
		}
		if agenciToSamRek(zajete, wykonawca) {
			continue
		}
		if !kontrolaZakresyStykaja(z.RangeStart, z.RangeEnd,
			int(zajete.ZakresOd), int(zajete.ZakresDo)) {

			continue
		}
		trzymajacy := agenciZlozZajecie(zajete)
		powod := "fragment od znaku " + strconv.Itoa(int(zajete.ZakresOd)) + " do " +
			strconv.Itoa(int(zajete.ZakresDo)) + " trzyma " + agenciNazwaZajecia(zajete) +
			"; zajęcie wygasa samo, a do tego czasu ten fragment należy do niego"
		return shared.StudioAgentsClaimResponse{
			Claimed: false, HeldBy: &trzymajacy, RefusalReason: &powod,
		}, nil
	}

	// Wykonawca dodatkowy nie wchodzi ponad granicę Operatora — granica
	// przemilczana byłaby pozorna.
	nastawy, err := a.agenciNastawy(ctx, nil, &dokument.Kod, nil)
	if err != nil {
		return shared.StudioAgentsClaimResponse{}, err
	}
	if !nastawy.MultiAgentEnabled {
		czynni := agenciCzynniWykonawcy(zajecia)
		if len(czynni) > 0 && !czynni[agenciKluczReki(wykonawca)] {
			return shared.StudioAgentsClaimResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodeConflict, "moduł Studio: nad dokumentem "+
					dokument.Kod+" pracuje już inny wykonawca, a praca wielu agentów naraz "+
					"jest wyłączona. Włącza ją jawne ustawienie Operatora "+
					"(studio.agents.settings.set, multiAgentEnabled)."))
		}
	}
	if nastawy.MaxConcurrentAgents != nil && *nastawy.MaxConcurrentAgents > 0 {
		czynni := agenciCzynniWykonawcy(zajecia)
		if !czynni[agenciKluczReki(wykonawca)] && len(czynni) >= *nastawy.MaxConcurrentAgents {
			return shared.StudioAgentsClaimResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodeConflict, "moduł Studio: nad dokumentem "+
					dokument.Kod+" pracuje już "+strconv.Itoa(len(czynni))+" wykonawców, "+
					"a Operator ustawił granicę "+strconv.Itoa(*nastawy.MaxConcurrentAgents)+
					". Zwolnij zajęcie (studio.agents.release) albo podnieś granicę."))
		}
	}

	waznosc := int64(agenciWaznoscDomyslna)
	if z.TtlSeconds != nil && *z.TtlSeconds > 0 {
		waznosc = int64(*z.TtlSeconds)
	}
	zapisane, err := skladnica.ZapiszZajecieFragmentu(ctx, dokument.ID,
		dane.ZajecieFragmentuStudia{
			Kod:             nowyIdentyfikator(przedrostekZajeciaStudia),
			WykonawcaRodzaj: string(shared.StudioAuthorModel),
			AgentKod:        wykonawca.AgentKod,
			AgentNazwa:      wykonawca.AgentNazwa,
			PodagentKod:     wykonawca.PodagentKod,
			Stan:            string(shared.StudioAgentSlotStateWorking),
			ZakresOd:        int64(z.RangeStart),
			ZakresDo:        int64(z.RangeEnd),
			ZadanieKod:      z.TaskId,
		}, waznosc)
	if err != nil {
		return shared.StudioAgentsClaimResponse{}, bladStudio(err)
	}
	zlozone := agenciZlozZajecie(zapisane)
	return shared.StudioAgentsClaimResponse{Claimed: true, Slot: &zlozone}, nil
}

// ZwolnijFragment obsługuje studio.agents.release: zwalnia fragment trzymany
// przez wskazanego wykonawcę.
func (a *adapterStudia) ZwolnijFragment(ctx context.Context,
	z shared.StudioAgentsReleaseRequest) (shared.StudioAgentsReleaseResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioAgentsReleaseResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioAgentsReleaseResponse{}, err
	}
	wszystkie := z.All != nil && *z.All
	if !kontrolaTekstNiepusty(z.SlotId) && !wszystkie {
		return shared.StudioAgentsReleaseResponse{}, bladWskazaniaStudio(
			"zwolnienie bez wskazania zajęcia i bez „all” — zwolnienie wszystkich " +
				"zajęć wykonawcy jest decyzją i musi być wypowiedziane wprost")
	}

	zwolnionych := int64(0)
	if kontrolaTekstNiepusty(z.SlotId) {
		zwolnione, err := skladnica.ZwolnijZajecieFragmentu(ctx, strings.TrimSpace(*z.SlotId))
		if err != nil {
			return shared.StudioAgentsReleaseResponse{}, bladStudio(err)
		}
		if !zwolnione {
			return shared.StudioAgentsReleaseResponse{},
				bladBrakuStudio("zajęcie nie istnieje albo już wygasło: " + *z.SlotId)
		}
		zwolnionych = 1
	} else {
		if zwolnionych, err = skladnica.ZwolnijZajeciaWykonawcy(ctx, dokument.ID,
			wartoscTekstu(z.AgentId), wartoscTekstu(z.SubagentId)); err != nil {

			return shared.StudioAgentsReleaseResponse{}, bladStudio(err)
		}
	}

	zajecia, err := skladnica.ZajeciaFragmentow(ctx, dokument.ID)
	if err != nil {
		return shared.StudioAgentsReleaseResponse{}, bladStudio(err)
	}
	wykaz := make([]shared.StudioAgentSlot, 0, len(zajecia))
	for _, zajete := range zajecia {
		wykaz = append(wykaz, agenciZlozZajecie(zajete))
	}
	return shared.StudioAgentsReleaseResponse{Released: int(zwolnionych), Slots: wykaz}, nil
}

// ZajeciaWykonawcow obsługuje studio.agents.slots.list. Wykazu nie przemiata:
// zajęcie wygasłe ma dać się zobaczyć na żądanie Operatora, bo po nim widać
// agenta ubitego w pół pracy.
func (a *adapterStudia) ZajeciaWykonawcow(ctx context.Context,
	z shared.StudioAgentsSlotsListRequest) (shared.StudioAgentsSlotsListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioAgentsSlotsListResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioAgentsSlotsListResponse{}, err
	}
	zajecia, err := skladnica.ZajeciaFragmentow(ctx, dokument.ID)
	if err != nil {
		return shared.StudioAgentsSlotsListResponse{}, bladStudio(err)
	}
	zWygaslymi := z.IncludeExpired != nil && *z.IncludeExpired

	odpowiedz := shared.StudioAgentsSlotsListResponse{Slots: []shared.StudioAgentSlot{}}
	czynni := map[string]bool{}
	for _, zajete := range zajecia {
		if zajete.Wygasle && !zWygaslymi {
			continue
		}
		if z.State != nil && *z.State != "" && zajete.Stan != string(*z.State) {
			continue
		}
		odpowiedz.Slots = append(odpowiedz.Slots, agenciZlozZajecie(zajete))
		if !zajete.Wygasle {
			czynni[wartoscTekstu(zajete.AgentKod)+"\x00"+wartoscTekstu(zajete.PodagentKod)] = true
		}
	}
	odpowiedz.ActiveAgents = len(czynni)
	return odpowiedz, nil
}

// SpieciaWykonawcowDokumentu obsługuje studio.agents.conflicts.list. Rdzeń
// niesie prawdę o tym, co się stało: czyja zmiana weszła, czyja została
// odłożona i dlaczego.
func (a *adapterStudia) SpieciaWykonawcowDokumentu(ctx context.Context,
	z shared.StudioAgentsConflictsListRequest) (shared.StudioAgentsConflictsListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioAgentsConflictsListResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioAgentsConflictsListResponse{}, err
	}
	wiersze, err := skladnica.SpieciaWykonawcow(ctx, dokument.ID)
	if err != nil {
		return shared.StudioAgentsConflictsListResponse{}, bladStudio(err)
	}
	odpowiedz := shared.StudioAgentsConflictsListResponse{
		Conflicts: []shared.StudioAgentConflict{},
	}
	for _, wiersz := range wiersze {
		if !wiersz.Domkniete {
			odpowiedz.OpenCount++
		}
		if z.Resolved != nil && *z.Resolved != wiersz.Domkniete {
			continue
		}
		if kontrolaTekstNiepusty(z.AgentId) &&
			wartoscTekstu(wiersz.WeszlaAgentKod) != *z.AgentId &&
			wartoscTekstu(wiersz.OdlozonaAgentKod) != *z.AgentId {

			continue
		}
		if z.Limit != nil && *z.Limit > 0 && len(odpowiedz.Conflicts) >= *z.Limit {
			continue
		}
		odpowiedz.Conflicts = append(odpowiedz.Conflicts, agenciZlozSpiecie(wiersz))
	}
	return odpowiedz, nil
}

// NastawyWykonawcow obsługuje studio.agents.settings.get: oddaje nastawy
// pętli wykonawczej obowiązujące w danym kontekście.
func (a *adapterStudia) NastawyWykonawcow(ctx context.Context,
	z shared.StudioAgentsSettingsGetRequest) (shared.StudioAgentsSettingsGetResponse, error) {

	nastawy, err := a.agenciNastawy(ctx, z.WindowId, z.DocumentId, z.SessionId)
	if err != nil {
		return shared.StudioAgentsSettingsGetResponse{}, err
	}
	return shared.StudioAgentsSettingsGetResponse{Settings: nastawy}, nil
}

// UstawNastawyWykonawcow obsługuje studio.agents.settings.set: zapisuje
// nastawę na wskazanym poziomie i rozgłasza jej zmianę.
func (a *adapterStudia) UstawNastawyWykonawcow(ctx context.Context,
	z shared.StudioAgentsSettingsSetRequest) (shared.StudioAgentsSettingsSetResponse, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioAgentsSettingsSetResponse{}, err
	}
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{})
	if wykonawca.czyWykonawca() {
		return shared.StudioAgentsSettingsSetResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodePermissionDenied, "moduł Studio: nastawy pętli wykonawczej "+
				"i pracy wielu agentów przestawia WYŁĄCZNIE Operator — "+
				wykonawca.nazwaWykonawcy()+" nie podnosi granicy, w której sam pracuje"))
	}

	zasieg := shared.ConfigScope(shared.ConfigScopeWindow)
	if z.Scope != nil && strings.TrimSpace(string(*z.Scope)) != "" {
		zasieg = *z.Scope
	}
	bytZasiegu := wartoscTekstu(z.ScopeId)
	if zasieg != shared.ConfigScopeGlobal && bytZasiegu == "" {
		return shared.StudioAgentsSettingsSetResponse{}, bladWskazaniaStudio(
			"nastawa na poziomie „" + string(zasieg) + "” bez wskazania bytu tego " +
				"poziomu — nastawa musi wiedzieć, czyja jest")
	}

	zapisy := []struct {
		klucz   string
		wartosc *string
	}{
		{kluczPetliWykonawczej, agenciZapisLogiczny(z.ExecutionLoopEnabled)},
		{kluczPracyWieluAgentow, agenciZapisLogiczny(z.MultiAgentEnabled)},
		{kluczNajwiecejAgentow, agenciZapisLiczby(z.MaxConcurrentAgents)},
		{kluczObiegowPetli, agenciZapisLiczby(z.LoopMaxIterations)},
		{kluczProguBezPostepu, agenciZapisLiczby(z.LoopNoProgressThreshold)},
		{kluczWymogZajecia, agenciZapisLogiczny(z.RequireFragmentClaim)},
	}
	if z.ConflictPolicy != nil && *z.ConflictPolicy != "" {
		if err := agenciSprawdzNastaweSpiecia(*z.ConflictPolicy); err != nil {
			return shared.StudioAgentsSettingsSetResponse{}, err
		}
		wartosc := string(*z.ConflictPolicy)
		zapisy = append(zapisy, struct {
			klucz   string
			wartosc *string
		}{kluczNastawySpiecia, &wartosc})
	}

	zmienione := 0
	for _, zapis := range zapisy {
		if zapis.wartosc == nil {
			continue
		}
		if err := skladnica.ZapiszNastaweZasieguStudia(ctx, dane.Ustawienie{
			Poziom: zasieg, KluczZasiegu: bytZasiegu, Klucz: zapis.klucz,
			Wartosc: zapis.wartosc, RodzajWartosci: string(konfig.RodzajTekst),
		}); err != nil {
			return shared.StudioAgentsSettingsSetResponse{}, bladStudio(err)
		}
		zmienione++
	}
	if zmienione == 0 {
		return shared.StudioAgentsSettingsSetResponse{}, bladWskazaniaStudio(
			"żądanie nie niesie ani jednej nastawy do przestawienia — odpowiedź " +
				"pomyślna kazałaby czytać to jako zmianę, której nie było")
	}
	// Rozgłoszenie nastawy dla okien nasłuchujących; rozstrzygacz opcjonalny,
	// brak go nie wywraca zapisu.
	if a.rozstrzygacz != nil {
		a.rozstrzygacz.Oglos(kluczPetliWykonawczej, kluczPracyWieluAgentow,
			kluczNajwiecejAgentow, kluczNastawySpiecia, kluczObiegowPetli,
			kluczProguBezPostepu, kluczWymogZajecia)
	}

	nastawy, err := a.agenciNastawyZasiegu(ctx, zasieg, bytZasiegu)
	if err != nil {
		return shared.StudioAgentsSettingsSetResponse{}, err
	}
	return shared.StudioAgentsSettingsSetResponse{Settings: nastawy}, nil
}

// ── Wspólne ─────────────────────────────────────────────────────────────────

// agenciNastawy oddaje nastawy obowiązujące w kontekście okna, dokumentu
// i karty sesji, idąc poziomami od węższego do szerszego.
func (a *adapterStudia) agenciNastawy(ctx context.Context,
	okno, dokumentKod, kartaSesji *string) (shared.StudioAgentSettings, error) {

	oknoNastaw := wartoscTekstu(okno)
	if oknoNastaw == "" && kontrolaTekstNiepusty(dokumentKod) {
		dokument, err := a.dokumentDoCzynnosci(ctx, *dokumentKod)
		if err != nil {
			return shared.StudioAgentSettings{}, err
		}
		oknoNastaw = dokument.Okno
	}
	if a.rozstrzygacz != nil {
		kontekst := konfig.Kontekst{Okno: oknoNastaw, KartaSesji: wartoscTekstu(kartaSesji)}
		nastawy := shared.StudioAgentSettings{
			ExecutionLoopEnabled: agenciCzytajLogiczny(a.rozstrzygacz, kontekst,
				kluczPetliWykonawczej),
			MultiAgentEnabled: agenciCzytajLogiczny(a.rozstrzygacz, kontekst,
				kluczPracyWieluAgentow),
			MaxConcurrentAgents: agenciCzytajLiczbe(a.rozstrzygacz, kontekst,
				kluczNajwiecejAgentow),
			LoopMaxIterations: agenciCzytajLiczbe(a.rozstrzygacz, kontekst, kluczObiegowPetli),
			LoopNoProgressThreshold: agenciCzytajLiczbe(a.rozstrzygacz, kontekst,
				kluczProguBezPostepu),
			RequireFragmentClaim: wskaznikLogiczny(agenciCzytajLogiczny(a.rozstrzygacz,
				kontekst, kluczWymogZajecia)),
		}
		wynik := a.rozstrzygacz.Rozstrzygnij(kontekst, kluczNastawySpiecia)
		if wynik.Wartosc != "" {
			nastawa := shared.StudioAgentConflictPolicy(wynik.Wartosc)
			nastawy.ConflictPolicy = &nastawa
		}
		if wynik.Poziom != "" {
			poziom := wynik.Poziom
			nastawy.Scope = &poziom
			if wynik.KluczZasiegu != "" {
				byt := wynik.KluczZasiegu
				nastawy.ScopeId = &byt
			}
		}
		return nastawy, nil
	}
	if oknoNastaw != "" {
		nastawy, err := a.agenciNastawyZasiegu(ctx, shared.ConfigScopeWindow, oknoNastaw)
		if err != nil {
			return shared.StudioAgentSettings{}, err
		}
		if nastawy.Scope != nil {
			return nastawy, nil
		}
	}
	return a.agenciNastawyZasiegu(ctx, shared.ConfigScopeGlobal, "")
}

// agenciNastawyZasiegu czyta nastawy zapisane na wskazanym poziomie zasięgu,
// bez sięgania po poziomy pozostałe.
func (a *adapterStudia) agenciNastawyZasiegu(ctx context.Context, zasieg shared.ConfigScope,
	bytZasiegu string) (shared.StudioAgentSettings, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioAgentSettings{}, err
	}
	nastawy := shared.StudioAgentSettings{}
	cokolwiek := false
	odczyt := func(klucz string) (string, error) {
		ustawienie, jest, err := skladnica.NastawaZasieguStudia(ctx, zasieg, bytZasiegu, klucz)
		if err != nil {
			return "", bladStudio(err)
		}
		if !jest {
			return "", nil
		}
		cokolwiek = true
		return wartoscTekstu(ustawienie.Wartosc), nil
	}
	for _, para := range []struct {
		klucz string
		cel   *bool
	}{
		{kluczPetliWykonawczej, &nastawy.ExecutionLoopEnabled},
		{kluczPracyWieluAgentow, &nastawy.MultiAgentEnabled},
	} {
		wartosc, err := odczyt(para.klucz)
		if err != nil {
			return shared.StudioAgentSettings{}, err
		}
		*para.cel = wartosc == "true"
	}
	for _, para := range []struct {
		klucz string
		cel   **int
	}{
		{kluczNajwiecejAgentow, &nastawy.MaxConcurrentAgents},
		{kluczObiegowPetli, &nastawy.LoopMaxIterations},
		{kluczProguBezPostepu, &nastawy.LoopNoProgressThreshold},
	} {
		wartosc, err := odczyt(para.klucz)
		if err != nil {
			return shared.StudioAgentSettings{}, err
		}
		if liczba, blad := strconv.Atoi(wartosc); blad == nil {
			kopia := liczba
			*para.cel = &kopia
		}
	}
	if wartosc, err := odczyt(kluczWymogZajecia); err != nil {
		return shared.StudioAgentSettings{}, err
	} else if wartosc != "" {
		nastawy.RequireFragmentClaim = wskaznikLogiczny(wartosc == "true")
	}
	if wartosc, err := odczyt(kluczNastawySpiecia); err != nil {
		return shared.StudioAgentSettings{}, err
	} else if wartosc != "" {
		nastawa := shared.StudioAgentConflictPolicy(wartosc)
		nastawy.ConflictPolicy = &nastawa
	}
	if cokolwiek {
		poziom := zasieg
		nastawy.Scope = &poziom
		if bytZasiegu != "" {
			byt := bytZasiegu
			nastawy.ScopeId = &byt
		}
	}
	return nastawy, nil
}

// agenciCzytajLogiczny rozstrzyga nastawę logiczną. Brak zapisu znaczy FAŁSZ —
// pętla wykonawcza i praca wielu agentów są domyślnie wyłączone.
func agenciCzytajLogiczny(rozstrzygacz *konfig.Rozstrzygacz, kontekst konfig.Kontekst,
	klucz string) bool {

	return rozstrzygacz.Rozstrzygnij(kontekst, klucz).Wartosc == "true"
}

// agenciCzytajLiczbe rozstrzyga nastawę liczbową. Brak zapisu zostaje BRAKIEM,
// nie zerem: zero znaczyłoby „ani jeden wykonawca", a brak znaczy „bez granicy
// ustawionej przez Operatora".
func agenciCzytajLiczbe(rozstrzygacz *konfig.Rozstrzygacz, kontekst konfig.Kontekst,
	klucz string) *int {

	wartosc := rozstrzygacz.Rozstrzygnij(kontekst, klucz).Wartosc
	liczba, err := strconv.Atoi(strings.TrimSpace(wartosc))
	if err != nil {
		return nil
	}
	return &liczba
}

// agenciZapisLogiczny składa wartość nastawy logicznej, gotową do zapisu
// w bazowej tabeli ustawień platformy.
func agenciZapisLogiczny(wskazanie *bool) *string {
	if wskazanie == nil {
		return nil
	}
	wartosc := "false"
	if *wskazanie {
		wartosc = "true"
	}
	return &wartosc
}

// agenciZapisLiczby składa wartość nastawy liczbowej, gotową do zapisu
// w bazowej tabeli ustawień platformy.
func agenciZapisLiczby(wskazanie *int) *string {
	if wskazanie == nil {
		return nil
	}
	wartosc := strconv.Itoa(*wskazanie)
	return &wartosc
}

// agenciSprawdzNastaweSpiecia odbija nastawę spięcia spoza wyliczenia
// kontraktu, dopuszczając wyłącznie wartości kontraktowi znane.
func agenciSprawdzNastaweSpiecia(nastawa shared.StudioAgentConflictPolicy) error {
	for _, dozwolona := range shared.WartosciStudioAgentConflictPolicy() {
		if nastawa == dozwolona {
			return nil
		}
	}
	nazwy := make([]string, 0, 3)
	for _, dozwolona := range shared.WartosciStudioAgentConflictPolicy() {
		nazwy = append(nazwy, string(dozwolona))
	}
	return bladWskazaniaStudio("nastawa spięcia „" + string(nastawa) +
		"” nie jest znana — wolno: " + strings.Join(nazwy, ", "))
}

// agenciToSamRek mówi, czy zajęcie należy do tej samej ręki. Wykonawca
// przedłużający własne zajęcie nie zderza się sam ze sobą.
func agenciToSamRek(zajecie dane.ZajecieFragmentuStudia, wykonawca kontrolaWykonawca) bool {
	return wartoscTekstu(zajecie.AgentKod) == wartoscTekstu(wykonawca.AgentKod) &&
		wartoscTekstu(zajecie.PodagentKod) == wartoscTekstu(wykonawca.PodagentKod)
}

// agenciKluczReki składa klucz tożsamości wykonawcy, używany do liczenia,
// ilu wykonawców pracuje naraz.
func agenciKluczReki(wykonawca kontrolaWykonawca) string {
	return wartoscTekstu(wykonawca.AgentKod) + "\x00" + wartoscTekstu(wykonawca.PodagentKod)
}

// agenciCzynniWykonawcy zbiera tożsamości wykonawców, którzy obecnie trzymają
// fragmenty tego dokumentu.
func agenciCzynniWykonawcy(zajecia []dane.ZajecieFragmentuStudia) map[string]bool {
	czynni := map[string]bool{}
	for _, zajete := range zajecia {
		if zajete.Wygasle {
			continue
		}
		czynni[wartoscTekstu(zajete.AgentKod)+"\x00"+wartoscTekstu(zajete.PodagentKod)] = true
	}
	return czynni
}

// agenciNazwaZajecia nazywa wykonawcę trzymającego fragment — nazwa widoczna
// przed kodem, bo Operator czyta nazwę, a nie identyfikator.
func agenciNazwaZajecia(zajecie dane.ZajecieFragmentuStudia) string {
	if zajecie.AgentNazwa != nil && *zajecie.AgentNazwa != "" {
		if zajecie.AgentKod != nil && *zajecie.AgentKod != "" {
			return *zajecie.AgentNazwa + " (" + *zajecie.AgentKod + ")"
		}
		return *zajecie.AgentNazwa
	}
	if zajecie.AgentKod != nil && *zajecie.AgentKod != "" {
		return "wykonawca " + *zajecie.AgentKod
	}
	if zajecie.PodagentKod != nil && *zajecie.PodagentKod != "" {
		return "podagent " + *zajecie.PodagentKod
	}
	return "wykonawca nienazwany"
}

// agenciZlozZajecie składa zajęcie kontraktu z wiersza warstwy danych.
// Zajęcie wygasłe wychodzi w stanie „idle”, nie „working”, bo pokazanie go
// jako pracującego byłoby nieprawdą.
func agenciZlozZajecie(wiersz dane.ZajecieFragmentuStudia) shared.StudioAgentSlot {
	od, do := int(wiersz.ZakresOd), int(wiersz.ZakresDo)
	zajeto := chwilaBazy(wiersz.Zajeto)
	stan := shared.StudioAgentSlotState(wiersz.Stan)
	if wiersz.Wygasle {
		stan = shared.StudioAgentSlotStateIdle
	}
	zajecie := shared.StudioAgentSlot{
		Id:         wiersz.Kod,
		DocumentId: wiersz.DokumentKod,
		Actor: shared.StudioActor{
			Kind:       shared.StudioAuthor(wiersz.WykonawcaRodzaj),
			AgentId:    wiersz.AgentKod,
			AgentName:  wiersz.AgentNazwa,
			SubagentId: wiersz.PodagentKod,
		},
		State:      stan,
		RangeStart: &od,
		RangeEnd:   &do,
		TaskId:     wiersz.ZadanieKod,
		ClaimedAt:  &zajeto,
	}
	if wiersz.Wygasa != nil && *wiersz.Wygasa != "" {
		wygasa := chwilaBazy(*wiersz.Wygasa)
		zajecie.ExpiresAt = &wygasa
	}
	return zajecie
}

// agenciZlozSpiecie składa spięcie kontraktu z wiersza warstwy danych,
// przenosząc jego pola do postaci odpowiedzi.
func agenciZlozSpiecie(wiersz dane.SpiecieWykonawcowStudia) shared.StudioAgentConflict {
	return shared.StudioAgentConflict{
		Id:         wiersz.Kod,
		DocumentId: wiersz.DokumentKod,
		RangeStart: int(wiersz.ZakresOd),
		RangeEnd:   int(wiersz.ZakresDo),
		AppliedActor: shared.StudioActor{
			Kind:      shared.StudioAuthor(wiersz.WeszlaRodzaj),
			AgentId:   wiersz.WeszlaAgentKod,
			AgentName: wiersz.WeszlaAgentNazwa,
		},
		DeferredActor: shared.StudioActor{
			Kind:      shared.StudioAuthor(wiersz.OdlozonaRodzaj),
			AgentId:   wiersz.OdlozonaAgentKod,
			AgentName: wiersz.OdlozonaAgentNaz,
		},
		Policy:           shared.StudioAgentConflictPolicy(wiersz.Nastawa),
		Reason:           wiersz.Powod,
		DeferredText:     wiersz.BrzmienieOdlozone,
		DeferredChangeId: wiersz.ZmianaSledzonaK,
		DeferredMarkupId: wiersz.ZnakowanieKod,
		Resolved:         wskaznikLogiczny(wiersz.Domkniete),
		CreatedAt:        chwilaBazy(wiersz.Utworzono),
	}
}
