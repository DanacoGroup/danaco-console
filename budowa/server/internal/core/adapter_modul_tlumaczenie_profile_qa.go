// Odpowiedzialność pliku: profile kontroli jakości translate.qa.profile.* i obieg zatwierdzeń panelu translate.approval.* modułu Translate.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// przedrostekProfiluQa znakuje identyfikator profilu kontroli jakości, przypisywany przy zapisie w bazie.
	przedrostekProfiluQa = "qap-"
	// przedrostekZatwierdzenia znakuje identyfikator kroku obiegu zatwierdzeń panelu tłumaczenia, przypisywany przy zapisie.
	przedrostekZatwierdzenia = "zat-"
	// autorNieznanyObiegu wchodzi, gdy wywołanie nie niesie tożsamości, zamiast wpisu konta z żądania klienta.
	autorNieznanyObiegu = "operator nieustalony"
)

// WykazProfiliQa obsługuje translate.qa.profile.list, zwracając profile kontroli jakości zapisane dla panelu.
func (a *adapterTlumaczenia) WykazProfiliQa(ctx context.Context,
	z shared.TranslateQaProfileListRequest) (shared.TranslateQaProfileListResponse, error) {

	zasieg := ""
	if z.Scope != nil {
		zasieg = string(*z.Scope)
	}
	profile, err := a.repozytorium.ProfileQa(ctx, zasieg, napisZeWskaznika(z.ScopeId))
	if err != nil {
		return shared.TranslateQaProfileListResponse{}, bladTlumaczenia(err)
	}
	wykaz := make([]shared.QaProfile, 0, len(profile))
	for _, profil := range profile {
		wykaz = append(wykaz, zlozProfilQa(profil))
	}
	return shared.TranslateQaProfileListResponse{Profiles: wykaz}, nil
}

// UstawProfilQa obsługuje translate.qa.profile.set, zapisując albo zmieniając profil kontroli jakości panelu.
func (a *adapterTlumaczenia) UstawProfilQa(ctx context.Context,
	z shared.TranslateQaProfileSetRequest) (shared.TranslateQaProfileSetResponse, error) {

	if strings.TrimSpace(z.Name) == "" {
		return shared.TranslateQaProfileSetResponse{}, bladWskazaniaTlumaczenia(
			"profil kontroli jakości bez nazwy")
	}
	if len(z.Checks) == 0 {
		return shared.TranslateQaProfileSetResponse{}, bladWskazaniaTlumaczenia(
			"profil bez ani jednej kontroli nie sprawdza niczego — wskaż kontrole")
	}
	kod := napisZeWskaznika(z.ProfileId)
	if strings.TrimSpace(kod) == "" {
		kod = nowyIdentyfikator(przedrostekProfiluQa)
	}
	profil := dane.ProfilQa{
		Kod:      kod,
		Nazwa:    z.Name,
		Opis:     z.Description,
		Zasieg:   string(shared.ConfigScopeGlobal),
		ZasiegID: z.ScopeId,
	}
	if z.Scope != nil {
		profil.Zasieg = string(*z.Scope)
	}
	for _, kontrola := range z.Checks {
		profil.Kontrole = append(profil.Kontrole, dane.KontrolaProfiluQa{
			Rodzaj:   string(kontrola.Kind),
			Waga:     string(kontrola.Severity),
			Wlaczona: kontrola.Enabled,
		})
	}
	zapisany, err := a.repozytorium.ZapiszProfilQa(ctx, profil)
	if err != nil {
		return shared.TranslateQaProfileSetResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateQaProfileSetResponse{Profile: zlozProfilQa(zapisany)}, nil
}

// UsunProfilQa obsługuje translate.qa.profile.delete, usuwając zapisany profil kontroli jakości panelu.
func (a *adapterTlumaczenia) UsunProfilQa(ctx context.Context,
	z shared.TranslateQaProfileDeleteRequest) (shared.TranslateQaProfileDeleteResponse, error) {

	if strings.TrimSpace(z.ProfileId) == "" {
		return shared.TranslateQaProfileDeleteResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.qa.profile.delete bez wskazania profilu")
	}
	zeszlo, err := a.repozytorium.UsunProfilQa(ctx, z.ProfileId)
	if err != nil {
		return shared.TranslateQaProfileDeleteResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateQaProfileDeleteResponse{Deleted: zeszlo}, nil
}

// zlozProfilQa przekłada wiersze profilu kontroli jakości na byt kontraktu wymiany z klientem Operatora.
func zlozProfilQa(profil dane.ProfilQa) shared.QaProfile {
	kontrole := make([]shared.QaProfileCheck, 0, len(profil.Kontrole))
	for _, kontrola := range profil.Kontrole {
		kontrole = append(kontrole, shared.QaProfileCheck{
			Kind:     shared.TranslationIssueKind(kontrola.Rodzaj),
			Severity: wagaOdmowyKontroli(kontrola.Waga),
			Enabled:  kontrola.Wlaczona,
		})
	}
	return shared.QaProfile{
		Id:          profil.Kod,
		Name:        profil.Nazwa,
		Description: profil.Opis,
		Scope:       shared.ConfigScope(profil.Zasieg),
		ScopeId:     profil.ZasiegID,
		Checks:      kontrole,
		UpdatedAt:   profil.Zaktualizowano,
	}
}

// UstawZatwierdzenie obsługuje translate.approval.set, zapisując krok obiegu zatwierdzeń wraz z migawką na panelu.
func (a *adapterTlumaczenia) UstawZatwierdzenie(ctx context.Context,
	z shared.TranslateApprovalSetRequest) (shared.TranslateApprovalSetResponse, error) {

	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateApprovalSetResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	if strings.TrimSpace(string(z.Stage)) == "" {
		return shared.TranslateApprovalSetResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.approval.set bez etapu obiegu")
	}
	zapis, err := a.repozytorium.ZapiszZatwierdzenie(ctx, dane.ZatwierdzeniePanelu{
		Kod:     nowyIdentyfikator(przedrostekZatwierdzenia),
		PanelID: panel.ID,
		Etap:    string(z.Stage),
		Autor:   autorObiegu(ctx),
		Uwaga:   z.Note,
	})
	if err != nil {
		return shared.TranslateApprovalSetResponse{}, bladTlumaczenia(err)
	}
	poZmianie, err := a.repozytorium.Panel(ctx, panel.Kod)
	if err != nil {
		return shared.TranslateApprovalSetResponse{}, bladTlumaczenia(err)
	}
	a.rozglosZmianePanelu(ctx, shared.ChangeKindUpdated, poZmianie)

	zapis.PanelKod = panel.Kod
	return shared.TranslateApprovalSetResponse{
		Panel:  zlozPanelTlumaczenia(poZmianie),
		Record: zlozZatwierdzenie(zapis),
	}, nil
}

// WykazZatwierdzen obsługuje `translate.approval.list`. Wskazanie panelu zawęża
// obieg do jednego panelu, wskazanie okna — do wszystkich paneli tego okna.
func (a *adapterTlumaczenia) WykazZatwierdzen(ctx context.Context,
	z shared.TranslateApprovalListRequest) (shared.TranslateApprovalListResponse, error) {

	var panelID, oknoID int64
	if kod := napisZeWskaznika(z.PanelId); strings.TrimSpace(kod) != "" {
		panel, err := a.repozytorium.Panel(ctx, kod)
		if err != nil {
			return shared.TranslateApprovalListResponse{}, bladNieznanegoPanelu(kod, err)
		}
		panelID = panel.ID
	}
	if kod := napisZeWskaznika(z.WindowId); panelID == 0 && strings.TrimSpace(kod) != "" {
		okno, err := a.repozytorium.Okno(ctx, kod)
		if err != nil {
			return shared.TranslateApprovalListResponse{}, bladNieznanegoOkna(kod, err)
		}
		oknoID = okno.ID
	}
	if panelID == 0 && oknoID == 0 {
		return shared.TranslateApprovalListResponse{}, bladWskazaniaTlumaczenia(
			"wykaz obiegu bez wskazania panelu ani okna objąłby całą instalację")
	}
	zapisy, err := a.repozytorium.Zatwierdzenia(ctx, panelID, oknoID)
	if err != nil {
		return shared.TranslateApprovalListResponse{}, bladTlumaczenia(err)
	}
	wykaz := make([]shared.ApprovalRecord, 0, len(zapisy))
	for _, zapis := range zapisy {
		wykaz = append(wykaz, zlozZatwierdzenie(zapis))
	}
	return shared.TranslateApprovalListResponse{Records: wykaz}, nil
}

// autorObiegu bierze tożsamość wołającego z kontekstu połączenia. Wywołanie
// spoza gniazda (sonda, bieg wewnętrzny, sprawdzian) tożsamości nie ma —
// obieg mówi wtedy wprost, że autora nie ustalono, zamiast wpisywać cudze konto.
func autorObiegu(ctx context.Context) string {
	tozsamosc := tozsamoscZKontekstu(ctx)
	if strings.TrimSpace(tozsamosc.IdKlienta) != "" {
		return tozsamosc.IdKlienta
	}
	return autorNieznanyObiegu
}

// zlozZatwierdzenie przekłada wiersz obiegu zatwierdzeń na byt kontraktu wymiany z klientem Operatora.
func zlozZatwierdzenie(zapis dane.ZatwierdzeniePanelu) shared.ApprovalRecord {
	return shared.ApprovalRecord{
		Id:        zapis.Kod,
		PanelId:   zapis.PanelKod,
		Stage:     shared.ApprovalStage(zapis.Etap),
		Author:    zapis.Autor,
		Note:      zapis.Uwaga,
		CreatedAt: zapis.Utworzono,
	}
}
