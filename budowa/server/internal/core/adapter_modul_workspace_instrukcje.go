// Odpowiedzialność pliku: instrukcje systemowe projektu — okno Instructions
// Panel modułu Workspace.
//
// Instrukcje są ustawieniem ośmiu poziomów zasięgu, więc zapisują się do
// tabeli `ustawienie` pod kluczem `workspace.instrukcje`, a warstwę
// obowiązującą wskazuje pakiet `internal/konfig` — ten sam, który rozstrzyga
// resztę konfiguracji. Drugiej tabeli instrukcji i drugiego porządku poziomów
// nie ma.
//
// Kontrakt daje Instructions Panel wyłącznie komendę zapisu, dlatego jej wynik
// niesie warstwę obowiązującą, a nie echo żądania: zapis na poziomie projektu
// bywa przykryty zapisem karty sesji albo okna, a panel ma pokazać treść,
// poziom i byt poziomu, które obowiązują po zapisie.
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

// ZapiszInstrukcje zapisuje instrukcje na wskazanym poziomie zasięgu i zwraca
// warstwę obowiązującą po zapisie.
//
// Brak wskazania poziomu znaczy poziom projektu — instrukcje obowiązują
// domyślnie w zakresie projektu.
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
	// Historia instrukcji rośnie przy zapisie, a nie przy odczycie panelu
	// „Wersje": wersja nieodłożona w chwili zmiany nie da się odtworzyć później
	// z niczego. Niepowodzenie odłożenia nie unieważnia zapisu, który już
	// osiadł — dlatego wynik nie wraca odmową.
	_, _ = a.repozytorium.ZapiszWersjeInstrukcjiWorkspace(ctx, dane.WersjaInstrukcjiWorkspace{
		ProjektID: projekt.ID, Identyfikator: nowyIdentyfikator("wsiv-"),
		Tresc: tresc, Odcisk: odciskTresci(tresc), Poziom: poziom, KluczZasiegu: bytPoziomu,
	})
	a.odnotujZdarzenieWorkspace(ctx, projekt.ID, shared.WorkspaceActivityKindInstructionsChanged,
		shared.ChangeKindUpdated, shared.WorkspaceEntityKindInstructions, projekt.Kod,
		"instrukcje projektu zapisane na poziomie "+string(poziom))
	return shared.WorkspaceInstructionsSetResponse{
		Instructions: a.instrukcjeObowiazujace(projekt.Kod, poziom, bytPoziomu, tresc),
	}, nil
}

// adresInstrukcji wskazuje miejsce zapisu: poziom zasięgu wraz z bytem tego
// poziomu. Poziom pominięty znaczy projekt, a byt pominięty na poziomie
// projektu — projekt z żądania. Poziom globalny bytu nie ma.
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

// instrukcjeObowiazujace rozstrzyga warstwę instrukcji obowiązującą w projekcie.
//
// Kontekst rozstrzygania niesie poziom projektu, a gdy zapis poszedł na poziom
// węższy (karta sesji, rola, okno) — także jego byt. Bez tego rozstrzyganie
// przeszłoby obok warstwy właśnie zapisanej i panel pokazałby jako
// obowiązującą warstwę szerszą.
//
// Brak rozstrzygacza nie jest błędem: wynik schodzi wtedy na treść zapisaną,
// bo zapis już się powiódł.
func (a *adapterPrzestrzeniRoboczej) instrukcjeObowiazujace(idProjektu string,
	poziom shared.ConfigScope, bytPoziomu, zapisana string) shared.WorkspaceInstructions {

	instrukcje := shared.WorkspaceInstructions{
		ProjectId: idProjektu, Content: zapisana, Scope: poziom,
		UpdatedAt: time.Now().UnixMilli(),
	}
	if bytPoziomu != "" {
		instrukcje.ScopeId = &bytPoziomu
	}
	if a.rozstrzygacz != nil {
		wynik := a.rozstrzygacz.Rozstrzygnij(kontekstInstrukcji(idProjektu, poziom, bytPoziomu),
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

// kontekstInstrukcji buduje kontekst rozstrzygania: projekt zawsze, a byt
// poziomu węższego wtedy, gdy zapis właśnie na nim osiadł.
func kontekstInstrukcji(idProjektu string, poziom shared.ConfigScope, bytPoziomu string) konfig.Kontekst {
	kontekst := konfig.Kontekst{Projekt: idProjektu}
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

// odciskTresci liczy skrót SHA-256 instrukcji. Służy rozpoznaniu zmiany treści
// przez interfejs — niczego nie dopuszcza i niczego nie blokuje.
func odciskTresci(tresc string) string {
	if tresc == "" {
		return ""
	}
	suma := sha256.Sum256([]byte(tresc))
	return hex.EncodeToString(suma[:])
}
