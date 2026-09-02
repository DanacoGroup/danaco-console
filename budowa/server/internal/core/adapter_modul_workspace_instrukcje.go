// Plik obsługuje instrukcje systemowe projektu dla okna Instructions Panel modułu Workspace. Instrukcje są ustawieniem ośmiu poziomów zasięgu, zapisywanym pod kluczem `workspace.instrukcje`.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// kluczInstrukcjiProjektu jest kluczem ustawienia niosącego instrukcje
// systemowe. Odpowiada kolumnie `ustawienie.klucz`.
const kluczInstrukcjiProjektu = "workspace.instrukcje"

func (a *adapterPrzestrzeniRoboczej) ZapiszInstrukcje(ctx context.Context,
	z shared.WorkspaceInstructionsSetRequest) (shared.WorkspaceInstructionsSetResponse, error) {

	projekt, err := a.projektDlaZapisu(ctx, z.ProjectId)
	if err != nil {
		return shared.WorkspaceInstructionsSetResponse{}, err
	}
	if a.konfiguracja == nil {
		return shared.WorkspaceInstructionsSetResponse{},
			bladProjektu("instrukcje nie mają gdzie zostać zapisane — brak repozytorium konfiguracji")
	}
	poziom, bytPoziomu := adresInstrukcji(z, projekt.Kod)
	tresc := z.Content
	err = a.konfiguracja.Ustaw(ctx, dane.Ustawienie{
		Poziom: poziom, KluczZasiegu: bytPoziomu, Klucz: kluczInstrukcjiProjektu,
		Wartosc: &tresc, RodzajWartosci: string(konfig.RodzajTekst),
	})
	if err != nil {
		return shared.WorkspaceInstructionsSetResponse{}, err
	}
	if err := a.repozytorium.OdnotujCzynnosc(ctx, projekt.ID); err != nil {
		return shared.WorkspaceInstructionsSetResponse{}, err
	}
	// Historia instrukcji rośnie przy zapisie, nie przy odczycie panelu „Wersje".
	_, _ = a.repozytorium.ZapiszWersjeInstrukcjiWorkspace(ctx, dane.WersjaInstrukcjiWorkspace{
		ProjektID: projekt.ID, Identyfikator: nowyIdentyfikator("wsiv-"),
		Tresc: tresc, Odcisk: odciskTresci(tresc), Poziom: poziom, KluczZasiegu: bytPoziomu,
	})
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindInstructionsChanged,
		shared.ChangeKindUpdated, shared.WorkspaceEntityKindInstructions, projekt.Kod,
		"instrukcje projektu zapisane na poziomie "+string(poziom))
	return shared.WorkspaceInstructionsSetResponse{
		Instructions: a.instrukcjeObowiazujace(ctx, projekt.Kod, poziom, bytPoziomu, tresc),
	}, nil
}

func adresInstrukcji(z shared.WorkspaceInstructionsSetRequest, idProjektu string) (shared.ConfigScope, string) {
	poziom := shared.ConfigScope(shared.ConfigScopeProject)
	if z.Scope != nil && *z.Scope != "" {
		poziom = *z.Scope
	}
	if z.ScopeId != nil && *z.ScopeId != "" {
		return poziom, *z.ScopeId
	}
	if poziom == shared.ConfigScopeProject {
		return poziom, idProjektu
	}
	return poziom, ""
}

func (a *adapterPrzestrzeniRoboczej) instrukcjeObowiazujace(ctx context.Context, idProjektu string,
	poziom shared.ConfigScope, bytPoziomu, zapisana string) shared.WorkspaceInstructions {

	instrukcje := shared.WorkspaceInstructions{
		ProjectId: idProjektu, Content: zapisana, Scope: poziom,
		UpdatedAt: time.Now().UnixMilli(),
	}
	if bytPoziomu != "" {
		instrukcje.ScopeId = &bytPoziomu
	}
	if a.rozstrzygacz != nil {
		wynik := a.rozstrzygacz.Rozstrzygnij(kontekstInstrukcji(ctx, idProjektu, poziom, bytPoziomu),
			kluczInstrukcjiProjektu)
		if wynik.Pochodzenie == konfig.PochodzenieZapis {
			instrukcje.Content, instrukcje.Scope = wynik.Wartosc, wynik.Poziom
			instrukcje.ScopeId = nil
			if wynik.KluczZasiegu != "" {
				byt := wynik.KluczZasiegu
				instrukcje.ScopeId = &byt
			}
		}
	}
	if odcisk := odciskTresci(instrukcje.Content); odcisk != "" {
		instrukcje.ContentHash = &odcisk
	}
	return instrukcje
}

func kontekstInstrukcji(ctx context.Context, idProjektu string, poziom shared.ConfigScope,
	bytPoziomu string) konfig.Kontekst {

	kontekst := konfig.Kontekst{Projekt: idProjektu, KontoOperatora: dane.KontoOperatora(ctx)}
	switch poziom {
	case shared.ConfigScopeSession:
		kontekst.KartaSesji = bytPoziomu
	case shared.ConfigScopeRole:
		kontekst.Rola = bytPoziomu
	case shared.ConfigScopeWindow:
		kontekst.Okno = bytPoziomu
	case shared.ConfigScopeModule:
		kontekst.Modul = bytPoziomu
	case shared.ConfigScopeModulePair:
		kontekst.ParaModulow = bytPoziomu
	case shared.ConfigScopeEnvironment:
		kontekst.Srodowisko = bytPoziomu
	}
	return kontekst
}

func odciskTresci(tresc string) string {
	if tresc == "" {
		return ""
	}
	suma := sha256.Sum256([]byte(tresc))
	return hex.EncodeToString(suma[:])
}
