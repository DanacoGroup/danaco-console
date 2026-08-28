// Odpowiedzialność pliku: moduł Studio — byty, które NIE należą do jednego
// dokumentu: operacje, łańcuchy, szablony i profile wydania; dwie cechy
// dokumentu.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	przedrostekOperacjiStudia  = "studio-op-"
	przedrostekLancuchaStudia  = "studio-lan-"
	przedrostekProfiluStudia   = "studio-prof-"
	przedrostekPrzebieguStudia = "studio-bieg-"
)

// ── Operacje własne Tools Panel ─────────────────────────────────────────────

// ZapiszOperacje obsługuje studio.operation.save, zakładając albo
// zmieniając zapisany prompt Operatora w danym zasięgu.
func (a *adapterStudia) ZapiszOperacje(ctx context.Context,
	z shared.StudioOperationSaveRequest) (shared.StudioOperationSaveResponse, error) {

	if strings.TrimSpace(z.Name) == "" || strings.TrimSpace(z.Prompt) == "" {
		return shared.StudioOperationSaveResponse{},
			bladWskazaniaStudio("operacja własna wymaga nazwy i treści promptu")
	}
	kod := wartoscTekstu(z.OperationId)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekOperacjiStudia)
	}
	zasieg, zasiegID := zasiegKatalogowyStudia(z.Scope, z.ScopeId, "global")

	zapisany, err := a.repozytorium.ZapiszOperacje(ctx, dane.WpisKatalogowyStudia{
		Kod: kod, Nazwa: z.Name, Rodzaj: z.Category, Ladunek: z.Prompt,
		Zasieg: zasieg, ZasiegID: zasiegID,
	})
	if err != nil {
		return shared.StudioOperationSaveResponse{}, bladStudio(err)
	}
	return shared.StudioOperationSaveResponse{Operation: złóżOperacje(zapisany)}, nil
}

// Operacje obsługuje studio.operation.list, oddając wykaz zapisanych
// operacji własnych widocznych w danym zasięgu.
func (a *adapterStudia) Operacje(ctx context.Context,
	z shared.StudioOperationListRequest) (shared.StudioOperationListResponse, error) {

	zasieg, zasiegID := zasiegKatalogowyStudia(z.Scope, z.ScopeId, "global")
	wiersze, err := a.repozytorium.Operacje(ctx, zasieg, zasiegID)
	if err != nil {
		return shared.StudioOperationListResponse{}, bladStudio(err)
	}
	kategoria := wartoscTekstu(z.Category)
	operacje := []shared.StudioOperation{}
	for _, wiersz := range wiersze {
		if kategoria != "" && wiersz.Rodzaj != kategoria {
			continue
		}
		operacje = append(operacje, złóżOperacje(wiersz))
	}
	return shared.StudioOperationListResponse{Operations: operacje}, nil
}

// UsunOperacje obsługuje studio.operation.delete; usunięcie bytu, którego
// nie ma, wraca odmową, a nie polem logicznym.
func (a *adapterStudia) UsunOperacje(ctx context.Context,
	z shared.StudioOperationDeleteRequest) (shared.StudioOperationDeleteResponse, error) {

	if z.OperationId == "" {
		return shared.StudioOperationDeleteResponse{},
			bladWskazaniaStudio("usunięcie bez wskazania operacji")
	}
	usunieto, err := a.repozytorium.UsunOperacje(ctx, z.OperationId)
	if err != nil {
		return shared.StudioOperationDeleteResponse{}, bladStudio(err)
	}
	if !usunieto {
		return shared.StudioOperationDeleteResponse{},
			bladBrakuStudio("operacja własna nie istnieje: " + z.OperationId)
	}
	return shared.StudioOperationDeleteResponse{Deleted: true}, nil
}

// ── Łańcuchy operacji ───────────────────────────────────────────────────────

// ZapiszLancuch obsługuje studio.chain.save, zakładając albo zmieniając
// łańcuch operacji widoczny w danym zasięgu.
func (a *adapterStudia) ZapiszLancuch(ctx context.Context,
	z shared.StudioChainSaveRequest) (shared.StudioChainSaveResponse, error) {

	if strings.TrimSpace(z.Name) == "" {
		return shared.StudioChainSaveResponse{}, bladWskazaniaStudio("łańcuch bez nazwy")
	}
	if len(z.Steps) == 0 {
		return shared.StudioChainSaveResponse{},
			bladWskazaniaStudio("łańcuch bez kroków — sekwencja pusta niczego nie wykona")
	}
	kroki, err := json.Marshal(z.Steps)
	if err != nil {
		return shared.StudioChainSaveResponse{}, bladStudio(err)
	}
	kod := wartoscTekstu(z.ChainId)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekLancuchaStudia)
	}
	zasieg, zasiegID := zasiegKatalogowyStudia(z.Scope, z.ScopeId, "project")

	zapisany, err := a.repozytorium.ZapiszLancuch(ctx, dane.WpisKatalogowyStudia{
		Kod: kod, Nazwa: z.Name, Ladunek: string(kroki), Zasieg: zasieg, ZasiegID: zasiegID,
	})
	if err != nil {
		return shared.StudioChainSaveResponse{}, bladStudio(err)
	}
	lancuch, err := złóżLancuch(zapisany)
	if err != nil {
		return shared.StudioChainSaveResponse{}, err
	}
	return shared.StudioChainSaveResponse{Chain: lancuch}, nil
}

// Lancuchy obsługuje studio.chain.list, oddając wykaz łańcuchów operacji
// widocznych zespołowi w danym zasięgu.
func (a *adapterStudia) Lancuchy(ctx context.Context,
	z shared.StudioChainListRequest) (shared.StudioChainListResponse, error) {

	zasieg, zasiegID := zasiegKatalogowyStudia(z.Scope, z.ScopeId, "project")
	wiersze, err := a.repozytorium.Lancuchy(ctx, zasieg, zasiegID)
	if err != nil {
		return shared.StudioChainListResponse{}, bladStudio(err)
	}
	lancuchy := []shared.StudioChain{}
	for _, wiersz := range wiersze {
		lancuch, err := złóżLancuch(wiersz)
		if err != nil {
			return shared.StudioChainListResponse{}, err
		}
		lancuchy = append(lancuchy, lancuch)
	}
	return shared.StudioChainListResponse{Chains: lancuchy}, nil
}

// UruchomLancuch obsługuje studio.chain.run, oddając identyfikator
// przebiegu wraz z liczbą kroków, bez wykonania.
func (a *adapterStudia) UruchomLancuch(ctx context.Context,
	z shared.StudioChainRunRequest) (shared.StudioChainRunResponse, error) {

	if z.WindowId == "" {
		return shared.StudioChainRunResponse{}, bladWskazaniaStudio("uruchomienie łańcucha bez okna")
	}
	if _, err := a.dokumentDoCzynnosci(ctx, z.DocumentId); err != nil {
		return shared.StudioChainRunResponse{}, err
	}
	if z.ChainId == "" {
		return shared.StudioChainRunResponse{},
			bladWskazaniaStudio("uruchomienie bez wskazania łańcucha")
	}
	wiersz, err := a.repozytorium.Lancuch(ctx, z.ChainId)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.StudioChainRunResponse{},
				bladBrakuStudio("łańcuch nie istnieje: " + z.ChainId)
		}
		return shared.StudioChainRunResponse{}, bladStudio(err)
	}
	lancuch, err := złóżLancuch(wiersz)
	if err != nil {
		return shared.StudioChainRunResponse{}, err
	}
	// Zakres selection bez granic zaznaczenia jest żądaniem sprzecznym: nie
	// ma czego wyciąć z treści.
	if z.Scope == shared.StudioOperationScopeSelection && (z.SelectionStart == nil || z.SelectionEnd == nil) {
		return shared.StudioChainRunResponse{},
			bladWskazaniaStudio("zakres „zaznaczenie” bez granic zaznaczenia")
	}
	return shared.StudioChainRunResponse{
		RunId: nowyIdentyfikator(przedrostekPrzebieguStudia),
		Steps: len(lancuch.Steps),
	}, nil
}

// ── Szablony dokumentów ─────────────────────────────────────────────────────

// Szablony obsługuje studio.template.list, oddając wykaz szablonów
// dokumentów widocznych w danym zasięgu Operatora.
func (a *adapterStudia) Szablony(ctx context.Context,
	_ shared.StudioTemplateListRequest) (shared.StudioTemplateListResponse, error) {

	wiersze, err := a.repozytorium.Szablony(ctx)
	if err != nil {
		return shared.StudioTemplateListResponse{}, bladStudio(err)
	}
	szablony := []shared.StudioTemplate{}
	for _, wiersz := range wiersze {
		szablon, err := złóżSzablon(wiersz)
		if err != nil {
			return shared.StudioTemplateListResponse{}, err
		}
		szablony = append(szablony, szablon)
	}
	return shared.StudioTemplateListResponse{Templates: szablony}, nil
}

// ZastosujSzablon obsługuje studio.template.apply, podstawiając wartości
// pod znaczniki wprost w treści szablonu.
func (a *adapterStudia) ZastosujSzablon(ctx context.Context,
	z shared.StudioTemplateApplyRequest) (shared.StudioTemplateApplyResponse, error) {

	if z.WindowId == "" {
		return shared.StudioTemplateApplyResponse{}, bladWskazaniaStudio("zastosowanie szablonu bez okna")
	}
	if z.TemplateId == "" {
		return shared.StudioTemplateApplyResponse{},
			bladWskazaniaStudio("zastosowanie bez wskazania szablonu")
	}
	wiersz, err := a.repozytorium.Szablon(ctx, z.TemplateId)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.StudioTemplateApplyResponse{},
				bladBrakuStudio("szablon nie istnieje: " + z.TemplateId)
		}
		return shared.StudioTemplateApplyResponse{}, bladStudio(err)
	}

	tresc := wypelnijSzablon(wiersz.Tresc, z.Values)
	tytul := z.Title
	if tytul == nil {
		nazwa := wiersz.Nazwa
		tytul = &nazwa
	}
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dane.DokumentStudia{
		Kod:    nowyIdentyfikator(przedrostekDokumentuStudio),
		Okno:   z.WindowId,
		Tytul:  tytul,
		Format: shared.StudioDocumentFormat(wiersz.Format),
		Tresc:  &tresc,
	})
	if err != nil {
		return shared.StudioTemplateApplyResponse{}, bladStudio(err)
	}
	zapisany, err = a.zalozWersjeDokumentu(ctx, zapisany, tresc, shared.StudioAuthorUzytkownik, nil)
	if err != nil {
		return shared.StudioTemplateApplyResponse{}, err
	}
	return shared.StudioTemplateApplyResponse{Document: a.zlozDokument(zapisany)}, nil
}

// wypelnijSzablon podstawia wartości pod znaczniki nazwa w treści
// szablonu, zostawiając brakujące pola widoczne.
func wypelnijSzablon(tresc string, wartosci json.RawMessage) string {
	if len(wartosci) == 0 {
		return tresc
	}
	var pola map[string]any
	if err := json.Unmarshal(wartosci, &pola); err != nil {
		return tresc
	}
	for nazwa, wartosc := range pola {
		napis, gotowy := wartosc.(string)
		if !gotowy {
			continue
		}
		tresc = strings.ReplaceAll(tresc, "{{"+nazwa+"}}", napis)
	}
	return tresc
}

// ── Profile wydania ─────────────────────────────────────────────────────────

// ZapiszProfilWydania obsługuje studio.export.profile.save, zakładając
// albo zmieniając profil wydania dokumentu.
func (a *adapterStudia) ZapiszProfilWydania(ctx context.Context,
	z shared.StudioExportProfileSaveRequest) (shared.StudioExportProfileSaveResponse, error) {

	if strings.TrimSpace(z.Name) == "" || strings.TrimSpace(z.Format) == "" {
		return shared.StudioExportProfileSaveResponse{},
			bladWskazaniaStudio("profil wydania wymaga nazwy i formatu docelowego")
	}
	ustawienia := ""
	if z.PageSetup != nil {
		zapis, err := json.Marshal(z.PageSetup)
		if err != nil {
			return shared.StudioExportProfileSaveResponse{}, bladStudio(err)
		}
		ustawienia = string(zapis)
	}
	kod := wartoscTekstu(z.ProfileId)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekProfiluStudia)
	}
	zasieg, zasiegID := zasiegKatalogowyStudia(z.Scope, z.ScopeId, "global")

	zapisany, err := a.repozytorium.ZapiszProfilWydania(ctx, dane.WpisKatalogowyStudia{
		Kod: kod, Nazwa: z.Name, Rodzaj: z.Format, Ladunek: ustawienia,
		Zasieg: zasieg, ZasiegID: zasiegID,
	})
	if err != nil {
		return shared.StudioExportProfileSaveResponse{}, bladStudio(err)
	}
	profil, err := złóżProfil(zapisany)
	if err != nil {
		return shared.StudioExportProfileSaveResponse{}, err
	}
	return shared.StudioExportProfileSaveResponse{Profile: profil}, nil
}

// ProfileWydania obsługuje studio.export.profile.list, oddając wykaz
// profili wydania widocznych w danym zasięgu.
func (a *adapterStudia) ProfileWydania(ctx context.Context,
	z shared.StudioExportProfileListRequest) (shared.StudioExportProfileListResponse, error) {

	zasieg, zasiegID := zasiegKatalogowyStudia(z.Scope, z.ScopeId, "global")
	wiersze, err := a.repozytorium.ProfileWydania(ctx, zasieg, zasiegID)
	if err != nil {
		return shared.StudioExportProfileListResponse{}, bladStudio(err)
	}
	profile := []shared.StudioExportProfile{}
	for _, wiersz := range wiersze {
		profil, err := złóżProfil(wiersz)
		if err != nil {
			return shared.StudioExportProfileListResponse{}, err
		}
		profile = append(profile, profil)
	}
	return shared.StudioExportProfileListResponse{Profiles: profile}, nil
}

// ── Cechy dokumentu i wersji ────────────────────────────────────────────────

// UstawFormatDokumentu obsługuje studio.document.format.set; zamiana
// treści jest tu WYBOREM, nie skutkiem ubocznym.
func (a *adapterStudia) UstawFormatDokumentu(ctx context.Context,
	z shared.StudioDocumentFormatSetRequest) (shared.StudioDocumentFormatSetResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDocumentFormatSetResponse{}, err
	}
	if z.Format == "" {
		return shared.StudioDocumentFormatSetResponse{},
			bladWskazaniaStudio("przestawienie formatu bez wskazania formatu")
	}
	if z.ConvertContent != nil && *z.ConvertContent {
		return shared.StudioDocumentFormatSetResponse{},
			bladWskazaniaStudio("zamiana treści przy przestawieniu formatu nie należy do tej komendy — " +
				"treść między formatami zamienia komenda obszaru dokumentów, a jej wynikiem jest " +
				"NOWY zasób magazynu, nie nadpisany dokument")
	}
	dokument.Format = z.Format
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return shared.StudioDocumentFormatSetResponse{}, bladStudio(err)
	}
	return shared.StudioDocumentFormatSetResponse{Document: a.zlozDokument(zapisany)}, nil
}

// UstawEtykieteWersji obsługuje studio.version.label.set; etykieta pusta
// ZDEJMUJE etykietę, zgodnie z opisem kontraktu.
func (a *adapterStudia) UstawEtykieteWersji(ctx context.Context,
	z shared.StudioVersionLabelSetRequest) (shared.StudioVersionLabelSetResponse, error) {

	if z.VersionId == "" {
		return shared.StudioVersionLabelSetResponse{},
			bladWskazaniaStudio("etykietowanie bez wskazania wersji")
	}
	wersja, err := a.repozytorium.Wersja(ctx, z.VersionId)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.StudioVersionLabelSetResponse{},
				bladBrakuStudio("wersja nie istnieje: " + z.VersionId)
		}
		return shared.StudioVersionLabelSetResponse{}, bladStudio(err)
	}

	wersja.Etykieta = nil
	if etykieta := strings.TrimSpace(wartoscTekstu(z.Label)); etykieta != "" {
		wersja.Etykieta = &etykieta
	}
	if z.Milestone != nil {
		wersja.KamienMilowy = *z.Milestone
	}
	zapisana, err := a.repozytorium.ZapiszWersje(ctx, wersja.DokumentID, wersja)
	if err != nil {
		return shared.StudioVersionLabelSetResponse{}, bladStudio(err)
	}
	return shared.StudioVersionLabelSetResponse{Version: a.zlozWersje(zapisana)}, nil
}

// ── Składanie odpowiedzi ────────────────────────────────────────────────────

// zasiegKatalogowyStudia sprowadza zasięg żądania do postaci zapisywanej
// w bazie; zasięg globalny nie ma identyfikatora bytu.
func zasiegKatalogowyStudia(zasieg *shared.ConfigScope, zasiegID *string, domyslny string) (string, *string) {
	nazwa := domyslny
	if zasieg != nil && *zasieg != "" {
		nazwa = string(*zasieg)
	}
	if nazwa == string(shared.ConfigScopeGlobal) {
		return nazwa, nil
	}
	return nazwa, zasiegID
}

func złóżOperacje(wiersz dane.WpisKatalogowyStudia) shared.StudioOperation {
	return shared.StudioOperation{
		Id:       wiersz.Kod,
		Name:     wiersz.Nazwa,
		Category: wiersz.Rodzaj,
		Prompt:   wiersz.Ladunek,
		// Wpis w tej tabeli powstaje wyłącznie z zapisu Operatora, nie
		// z operacji fabrycznych rdzenia.
		Builtin: false,
		Scope:   shared.ConfigScope(wiersz.Zasieg),
	}
}

func złóżLancuch(wiersz dane.WpisKatalogowyStudia) (shared.StudioChain, error) {
	kroki := []shared.StudioChainStep{}
	if wiersz.Ladunek != "" {
		if err := json.Unmarshal([]byte(wiersz.Ladunek), &kroki); err != nil {
			return shared.StudioChain{}, bladStudio(err)
		}
	}
	return shared.StudioChain{
		Id: wiersz.Kod, Name: wiersz.Nazwa, Steps: kroki,
		Scope: shared.ConfigScope(wiersz.Zasieg),
	}, nil
}

func złóżProfil(wiersz dane.WpisKatalogowyStudia) (shared.StudioExportProfile, error) {
	profil := shared.StudioExportProfile{
		Id: wiersz.Kod, Name: wiersz.Nazwa, Format: wiersz.Rodzaj,
		Scope: shared.ConfigScope(wiersz.Zasieg),
	}
	if wiersz.Ladunek != "" {
		var ustawienia shared.StudioPageSetup
		if err := json.Unmarshal([]byte(wiersz.Ladunek), &ustawienia); err != nil {
			return shared.StudioExportProfile{}, bladStudio(err)
		}
		profil.PageSetup = &ustawienia
	}
	return profil, nil
}

func złóżSzablon(wiersz dane.SzablonStudia) (shared.StudioTemplate, error) {
	szablon := shared.StudioTemplate{
		Id: wiersz.Kod, Name: wiersz.Nazwa, Description: wiersz.Opis,
		Format:  shared.StudioDocumentFormat(wiersz.Format),
		Builtin: wiersz.Fabryczny,
	}
	if wiersz.PolaJSON != nil && *wiersz.PolaJSON != "" {
		var pola []shared.StudioTemplateField
		if err := json.Unmarshal([]byte(*wiersz.PolaJSON), &pola); err != nil {
			return shared.StudioTemplate{}, bladStudio(err)
		}
		szablon.Fields = pola
	}
	return szablon, nil
}
