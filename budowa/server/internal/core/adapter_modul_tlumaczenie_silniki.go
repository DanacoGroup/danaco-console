// Silniki przekladu i przebieg pakietowy modulu Translate: profile silnikow,
// porownanie kanalow, polityka pivota oraz uruchomienie zlecenia pakietowego
// translate.batch.run.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// przedrostekProfiluSilnika znakuje identyfikator profilu silnika nadany
	// przez rdzen przy jego zapisie.
	przedrostekProfiluSilnika = "sil-"
	// przedrostekZleceniaPakietu znakuje identyfikator przebiegu pakietowego
	// nadany przez rdzen przy zapisie.
	przedrostekZleceniaPakietu = "pak-"
)

// WykazProfiliSilnikow obsluguje komende translate.engine.profile.list, oddajac
// wykaz profili zapisanych w repozytorium.
func (a *adapterTlumaczenia) WykazProfiliSilnikow(ctx context.Context,
	z shared.TranslateEngineProfileListRequest) (shared.TranslateEngineProfileListResponse, error) {

	zasieg := ""
	if z.Scope != nil {
		zasieg = string(*z.Scope)
	}
	profile, err := a.repozytorium.ProfileSilnikow(ctx, zasieg, napisZeWskaznika(z.ScopeId))
	if err != nil {
		return shared.TranslateEngineProfileListResponse{}, bladTlumaczenia(err)
	}
	wykaz := make([]shared.EngineProfile, 0, len(profile))
	for _, profil := range profile {
		wykaz = append(wykaz, zlozProfilSilnika(profil))
	}
	return shared.TranslateEngineProfileListResponse{Profiles: wykaz}, nil
}

// UstawProfilSilnika obsluguje komende translate.engine.profile.set, zapisujaca
// profil silnika w repozytorium.
func (a *adapterTlumaczenia) UstawProfilSilnika(ctx context.Context,
	z shared.TranslateEngineProfileSetRequest) (shared.TranslateEngineProfileSetResponse, error) {

	if strings.TrimSpace(z.Name) == "" {
		return shared.TranslateEngineProfileSetResponse{}, bladWskazaniaTlumaczenia(
			"profil silnika bez nazwy")
	}
	kod := napisZeWskaznika(z.ProfileId)
	if strings.TrimSpace(kod) == "" {
		kod = nowyIdentyfikator(przedrostekProfiluSilnika)
	}
	profil := dane.ProfilSilnika{
		Kod:         kod,
		Nazwa:       z.Name,
		Dziedzina:   z.Domain,
		Zasieg:      string(shared.ConfigScopeGlobal),
		ZasiegID:    z.ScopeId,
		Kanaly:      z.ChannelIds,
		Temperatura: z.Temperature,
	}
	if z.Scope != nil {
		profil.Zasieg = string(*z.Scope)
	}
	if z.Adaptive != nil {
		profil.Adaptacyjny = *z.Adaptive
	}
	if z.MemoryScope != nil {
		zasieg := string(*z.MemoryScope)
		profil.ZasiegPamieci = &zasieg
	}
	zapisany, err := a.repozytorium.ZapiszProfilSilnika(ctx, profil)
	if err != nil {
		return shared.TranslateEngineProfileSetResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateEngineProfileSetResponse{Profile: zlozProfilSilnika(zapisany)}, nil
}

// zlozProfilSilnika przeklada wiersz profilu z repozytorium na byt kontraktu
// oddawany w odpowiedzi komendy.
func zlozProfilSilnika(profil dane.ProfilSilnika) shared.EngineProfile {
	byt := shared.EngineProfile{
		Id:          profil.Kod,
		Name:        profil.Nazwa,
		Domain:      profil.Dziedzina,
		Scope:       shared.ConfigScope(profil.Zasieg),
		ScopeId:     profil.ZasiegID,
		ChannelIds:  profil.Kanaly,
		Adaptive:    profil.Adaptacyjny,
		Temperature: profil.Temperatura,
		UpdatedAt:   profil.Zaktualizowano,
	}
	if byt.ChannelIds == nil {
		byt.ChannelIds = []string{}
	}
	if profil.ZasiegPamieci != nil {
		zasieg := shared.TranslationMemoryScope(*profil.ZasiegPamieci)
		byt.MemoryScope = &zasieg
	}
	return byt
}

// PorownajSilniki obsługuje `translate.engine.compare`. Każdy wskazany kanał
// dostaje ten sam segment i to samo polecenie; różnice między wariantami są
// wtedy różnicami między modelami, a nie między poleceniami.
func (a *adapterTlumaczenia) PorownajSilniki(ctx context.Context,
	z shared.TranslateEngineCompareRequest) (shared.TranslateEngineCompareResponse, error) {

	if len(z.ChannelIds) == 0 {
		return shared.TranslateEngineCompareResponse{}, bladWskazaniaTlumaczenia(
			"porównanie bez wskazania kanałów nie ma czego z czym zestawić")
	}
	if strings.TrimSpace(z.Segment) == "" {
		return shared.TranslateEngineCompareResponse{}, bladWskazaniaTlumaczenia(
			"porównanie bez segmentu — nie ma czego tłumaczyć")
	}
	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateEngineCompareResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	okno, err := a.repozytorium.Okno(ctx, panel.OknoKod)
	if err != nil {
		return shared.TranslateEngineCompareResponse{}, bladNieznanegoOkna(panel.OknoKod, err)
	}

	// Trasa przekladu bierze sie z polityki pivota: jezyk pomiedzy idzie przez jezyk posredni.
	trasa := a.trasaPrzekladu(ctx, jezykZrodlowyZadania(okno.JezykZrodlowy), panel.Jezyk)

	warianty := make([]shared.EngineVariant, 0, len(z.ChannelIds))
	for _, kanal := range z.ChannelIds {
		wiazania, err := a.wiazaniaSlownikaDlaJezyka(ctx, panel.Jezyk)
		if err != nil {
			return shared.TranslateEngineCompareResponse{}, err
		}
		polecenie := poleceniePrzekladu(panel.Jezyk, panel.Ton, wiazania, z.Segment)
		przeklad, err := a.zapytajModel(ctx, okno.Kod, kanal, polecenie)
		if err != nil {
			return shared.TranslateEngineCompareResponse{}, err
		}
		if strings.TrimSpace(przeklad) == "" {
			return shared.TranslateEngineCompareResponse{}, bladTlumaczenia(errPustyPrzeklad)
		}
		wariant := shared.EngineVariant{ChannelId: kanal, Text: strings.TrimSpace(przeklad)}
		if trasa != nil {
			wariant.Route = trasa
		}
		warianty = append(warianty, wariant)
	}
	return shared.TranslateEngineCompareResponse{Variants: warianty}, nil
}

// trasaPrzekladu składa trasę wedle polityki pivota zasięgu globalnego. Brak
// polityki znaczy trasę wprost, więc oddaje nic — trasa bez języka pośredniego
// nie niesie żadnej informacji ponad to, co i tak widać w panelu.
func (a *adapterTlumaczenia) trasaPrzekladu(ctx context.Context,
	jezykZrodlowy, jezykDocelowy string) *shared.PivotRoute {

	polityka, err := a.repozytorium.PolitykaPivotaZasiegu(ctx, string(shared.ConfigScopeGlobal), "")
	if err != nil {
		return nil
	}
	for _, para := range polityka.Pary {
		if para.JezykZrodla == jezykZrodlowy && para.JezykCelu == jezykDocelowy {
			pivot := para.JezykPivota
			return &shared.PivotRoute{
				SourceLanguage: jezykZrodlowy,
				PivotLanguage:  &pivot,
				TargetLanguage: jezykDocelowy,
			}
		}
	}
	if polityka.JezykDomyslny != nil && strings.TrimSpace(*polityka.JezykDomyslny) != "" {
		return &shared.PivotRoute{
			SourceLanguage: jezykZrodlowy,
			PivotLanguage:  polityka.JezykDomyslny,
			TargetLanguage: jezykDocelowy,
		}
	}
	return nil
}

// PolitykaPivota obsluguje komende translate.pivot.policy.get, oddajaca polityke
// pivota zapisana dla zasiegu.
func (a *adapterTlumaczenia) PolitykaPivota(ctx context.Context,
	z shared.TranslatePivotPolicyGetRequest) (shared.TranslatePivotPolicyGetResponse, error) {

	zasieg, zasiegID := zasiegZadania(z.Scope, z.ScopeId)
	polityka, err := a.repozytorium.PolitykaPivotaZasiegu(ctx, zasieg, zasiegID)
	if err != nil {
		// Polityki nieustawionej nie udajemy odmowa: pusta polityka znaczy przeklad wprost.
		return shared.TranslatePivotPolicyGetResponse{Policy: shared.PivotPolicy{
			Scope: shared.ConfigScope(zasieg),
			Pairs: []shared.PivotPair{},
		}}, nil
	}
	return shared.TranslatePivotPolicyGetResponse{Policy: zlozPolitykePivota(polityka)}, nil
}

// UstawPolitykePivota obsluguje komende translate.pivot.policy.set, zapisujaca
// polityke pivota dla zasiegu.
func (a *adapterTlumaczenia) UstawPolitykePivota(ctx context.Context,
	z shared.TranslatePivotPolicySetRequest) (shared.TranslatePivotPolicySetResponse, error) {

	zasieg, zasiegID := zasiegZadania(z.Scope, z.ScopeId)
	polityka := dane.PolitykaPivota{
		Zasieg:        zasieg,
		ZasiegID:      zasiegID,
		JezykDomyslny: z.DefaultPivot,
	}
	for _, para := range z.Pairs {
		if strings.TrimSpace(para.PivotLanguage) == "" {
			return shared.TranslatePivotPolicySetResponse{}, bladWskazaniaTlumaczenia(
				"para pivota bez języka pośredniego nie jest parą pivota")
		}
		polityka.Pary = append(polityka.Pary, dane.ParaPivota{
			JezykZrodla: para.SourceLanguage,
			JezykCelu:   para.TargetLanguage,
			JezykPivota: para.PivotLanguage,
		})
	}
	zapisana, err := a.repozytorium.ZapiszPolitykePivota(ctx, polityka)
	if err != nil {
		return shared.TranslatePivotPolicySetResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslatePivotPolicySetResponse{Policy: zlozPolitykePivota(zapisana)}, nil
}

// zasiegZadania sprowadza wskazanie zasięgu do pary wartości schematu. Zasięg
// nieokreślony jest zasięgiem globalnym — najszerszym, na którym nastawa
// obowiązuje wszędzie tam, gdzie nikt nie zapisał węższej.
func zasiegZadania(zasieg *shared.ConfigScope, zasiegID *string) (string, string) {
	nazwa := string(shared.ConfigScopeGlobal)
	if zasieg != nil && strings.TrimSpace(string(*zasieg)) != "" {
		nazwa = string(*zasieg)
	}
	return nazwa, napisZeWskaznika(zasiegID)
}

// zlozPolitykePivota przeklada wiersze polityki z repozytorium na byt kontraktu
// oddawany w odpowiedzi.
func zlozPolitykePivota(polityka dane.PolitykaPivota) shared.PivotPolicy {
	pary := make([]shared.PivotPair, 0, len(polityka.Pary))
	for _, para := range polityka.Pary {
		pary = append(pary, shared.PivotPair{
			SourceLanguage: para.JezykZrodla,
			TargetLanguage: para.JezykCelu,
			PivotLanguage:  para.JezykPivota,
		})
	}
	byt := shared.PivotPolicy{
		Scope:        shared.ConfigScope(polityka.Zasieg),
		DefaultPivot: polityka.JezykDomyslny,
		Pairs:        pary,
		UpdatedAt:    polityka.Zaktualizowano,
	}
	if polityka.ZasiegID != "" {
		zasiegID := polityka.ZasiegID
		byt.ScopeId = &zasiegID
	}
	return byt
}

// UruchomPakiet obsluguje komende translate.batch.run, wykonujaca zlecenie
// pakietowe wiersz po wierszu.
func (a *adapterTlumaczenia) UruchomPakiet(ctx context.Context,
	z shared.TranslateBatchRunRequest) (shared.TranslateBatchRunResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateBatchRunResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	if len(z.Operations) == 0 {
		return shared.TranslateBatchRunResponse{}, bladWskazaniaTlumaczenia(
			"przebieg pakietowy bez wskazania operacji nie zrobiłby niczego")
	}
	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateBatchRunResponse{}, bladTlumaczenia(err)
	}
	wskazane := map[string]struct{}{}
	for _, kod := range z.PanelIds {
		wskazane[kod] = struct{}{}
	}

	pozycje := []dane.PozycjaPakietuTlumaczenia{}
	for _, panel := range panele {
		if len(wskazane) > 0 {
			if _, jest := wskazane[panel.Kod]; !jest {
				continue
			}
		}
		for _, operacja := range z.Operations {
			stan, szczegol := a.wykonajOperacjePakietu(ctx, okno, panel, operacja, z.ChannelId)
			pozycje = append(pozycje, dane.PozycjaPakietuTlumaczenia{
				PanelKod: panel.Kod,
				Operacja: string(operacja),
				Stan:     stan,
				Szczegol: wskaznikNapisu(szczegol),
			})
		}
	}
	if len(pozycje) == 0 {
		return shared.TranslateBatchRunResponse{}, bladWskazaniaTlumaczenia(
			"okno nie ma paneli objętych wskazaniem — przebieg nie miałby ani jednej pozycji")
	}

	kod := nowyIdentyfikator(przedrostekZleceniaPakietu)
	if err := a.repozytorium.ZalozZleceniePakietu(ctx, kod, okno.ID, pozycje); err != nil {
		return shared.TranslateBatchRunResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateBatchRunResponse{QueueId: kod, ItemCount: len(pozycje)}, nil
}

// wykonajOperacjePakietu wykonuje jedną pozycję przebiegu i oddaje jej stan
// wraz ze szczegółem. Niepowodzenie pozycji nie przewraca całego przebiegu:
// pakiet ma dojść do końca i pokazać, co się udało, a co nie.
func (a *adapterTlumaczenia) wykonajOperacjePakietu(ctx context.Context,
	okno dane.OknoTlumaczenia, panel dane.PanelTlumaczenia,
	operacja shared.BatchOperationKind, kanal *string) (string, string) {

	switch operacja {
	case shared.BatchOperationKindTranslate:
		if okno.TekstZrodlowy == nil || strings.TrimSpace(*okno.TekstZrodlowy) == "" {
			return "error", "okno nie ma tekstu źródłowego"
		}
		przeklad, err := a.przetlumaczModelem(ctx, okno.Kod, panel.Jezyk, panel.Ton,
			*okno.TekstZrodlowy, kanal)
		if err != nil {
			return "error", err.Error()
		}
		zmieniony, err := a.repozytorium.UstawTlumaczenie(ctx, panel.Kod, &przeklad, nil)
		if err != nil {
			return "error", err.Error()
		}
		a.rozglosZmianePanelu(shared.ChangeKindUpdated, zmieniony)
		return "done", ""

	case shared.BatchOperationKindQualityCheck:
		if panel.Tresc == nil {
			return "error", "panel nie ma treści do kontroli"
		}
		tekstZrodlowy := ""
		if okno.TekstZrodlowy != nil {
			tekstZrodlowy = *okno.TekstZrodlowy
		}
		niezgodnosci := zbadajPanel(panel, tekstZrodlowy)
		if err := a.repozytorium.ZapiszNiezgodnosci(ctx, panel.ID, niezgodnosci); err != nil {
			return "error", err.Error()
		}
		return "done", ""

	case shared.BatchOperationKindProofread:
		if _, err := a.SprawdzKorekte(ctx, shared.TranslateProofreadRunRequest{
			PanelId: panel.Kod}); err != nil {
			return "error", err.Error()
		}
		return "done", ""

	case shared.BatchOperationKindExport:
		if _, err := a.EksportujPanel(ctx, shared.TranslatePanelExportRequest{
			PanelId: panel.Kod, Format: shared.ExportFormatMarkdown}); err != nil {
			return "error", err.Error()
		}
		return "done", ""
	}
	return "error", "nieznana operacja pakietu: " + string(operacja)
}
