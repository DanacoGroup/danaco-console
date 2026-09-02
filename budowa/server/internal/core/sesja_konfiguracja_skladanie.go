// Plik składa konfigurację sesji obowiązującą, rozstrzygając obszary po ośmiu poziomach zasięgu i trzech osiach wraz ze wskazaniem pochodzenia każdego obszaru, w kolejności adresów z pakietu konfig, bez dotykania dysku.
package core

import (
	"context"
	"encoding/json"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

func (a *adapterUstawienOsi) KonfiguracjaObowiazujaca(ctx context.Context,
	z shared.ConfigEffectiveGetRequest) (shared.ConfigEffectiveGetResponse, error) {

	obszary, err := obszaryZadania(z.Areas)
	if err != nil {
		return shared.ConfigEffectiveGetResponse{}, err
	}
	tresci, pochodzenie, err := a.rozstrzygnijObszary(ctx, kontekstKonfiguracjiSesji(ctx, z), obszary)
	if err != nil {
		return shared.ConfigEffectiveGetResponse{}, err
	}
	konfiguracja, err := zlozKonfiguracje(tresci)
	if err != nil {
		return shared.ConfigEffectiveGetResponse{}, err
	}
	obowiazujaca := shared.SessionConfigEffective{
		Config:           konfiguracja,
		Origins:          pochodzeniaObszarow(obszary, pochodzenie),
		WorkingDirectory: rozstrzygniecieKatalogu(konfiguracja),
		ResolvedAt:       time.Now().UnixMilli(),
	}
	if z.IncludeCapabilities != nil && *z.IncludeCapabilities {
		zdolnosci := deklaracjaZdolnosci(nil)
		obowiazujaca.Capabilities = &zdolnosci
	}
	return shared.ConfigEffectiveGetResponse{Effective: obowiazujaca}, nil
}

func (a *adapterUstawienOsi) KonfiguracjaSesjiOkna(ctx context.Context,
	idOkna, idSesji string) (shared.SessionConfig, error) {

	if a == nil {
		return shared.SessionConfig{}, nil
	}
	kontekst := konfig.Kontekst{Okno: idOkna, KartaSesji: idSesji, KontoOperatora: dane.KontoOperatora(ctx)}
	tresci, _, err := a.rozstrzygnijObszary(ctx, kontekst, obszaryKontraktu)
	if err != nil {
		return shared.SessionConfig{}, err
	}
	return zlozKonfiguracje(tresci)
}

func (a *adapterUstawienOsi) rozstrzygnijObszary(ctx context.Context, kontekst konfig.Kontekst,
	obszary []shared.SessionConfigArea) (map[shared.SessionConfigArea]json.RawMessage,
	map[shared.SessionConfigArea]konfig.Adres, error) {

	tresci := map[shared.SessionConfigArea]json.RawMessage{}
	pochodzenie := map[shared.SessionConfigArea]konfig.Adres{}
	for _, adres := range kontekst.Adresy() {
		zapisane, _, err := a.obszaryAdresu(ctx, adres, obszary)
		if err != nil {
			return nil, nil, err
		}
		for _, obszar := range obszary {
			tresc, jest := zapisane[obszar]
			if !jest {
				continue
			}
			if _, rozstrzygniety := tresci[obszar]; rozstrzygniety {
				continue
			}
			tresci[obszar] = tresc
			pochodzenie[obszar] = adres
		}
	}
	return tresci, pochodzenie, nil
}

func kontekstKonfiguracjiSesji(ctx context.Context, z shared.ConfigEffectiveGetRequest) konfig.Kontekst {
	kontekst := konfig.Kontekst{
		Okno:           wartoscTekstu(z.WindowId),
		KartaSesji:     wartoscTekstu(z.SessionId),
		KontoOperatora: dane.KontoOperatora(ctx),
	}
	if z.Scope != nil {
		ustawBytPoziomu(&kontekst, *z.Scope, wartoscTekstu(z.ScopeId))
	}
	os, bytOsi := osZadania(z.Axis, z.AxisId)
	switch konfig.OsLubPlatforma(os) {
	case konfig.OsModelu:
		kontekst.Model = bytOsi
	case konfig.OsKonta:
		kontekst.Konto = bytOsi
	}
	return kontekst
}

// ustawBytPoziomu wpisuje byt na wskazany poziom zasięgu. Jest odwrotnością
// konfig.Kontekst.Adres — pakiet konfig czyta byt poziomu, żądanie go podaje.
func ustawBytPoziomu(kontekst *konfig.Kontekst, poziom shared.ConfigScope, byt string) {
	if byt == "" {
		return
	}
	switch poziom {
	case shared.ConfigScopeEnvironment:
		kontekst.Srodowisko = byt
	case shared.ConfigScopeModule:
		kontekst.Modul = byt
	case shared.ConfigScopeModulePair:
		kontekst.ParaModulow = byt
	case shared.ConfigScopeProject:
		kontekst.Projekt = byt
	case shared.ConfigScopeSession:
		kontekst.KartaSesji = byt
	case shared.ConfigScopeRole:
		kontekst.Rola = byt
	case shared.ConfigScopeWindow:
		kontekst.Okno = byt
	}
}

func pochodzeniaObszarow(obszary []shared.SessionConfigArea,
	pochodzenie map[shared.SessionConfigArea]konfig.Adres) []shared.SessionConfigOrigin {

	wykaz := make([]shared.SessionConfigOrigin, 0, len(obszary))
	for _, obszar := range obszary {
		adres, znalezione := pochodzenie[obszar]
		if !znalezione {
			wykaz = append(wykaz, shared.SessionConfigOrigin{
				Area: obszar, Source: shared.SessionConfigSourceDefault,
			})
			continue
		}
		poziom := adres.Poziom
		wpis := shared.SessionConfigOrigin{
			Area: obszar, Source: shared.SessionConfigSourceOverride,
			Scope: &poziom, ScopeId: wskaznikTekstu(adres.KluczZasiegu),
		}
		if konfig.OsLubPlatforma(adres.Os) != konfig.OsPlatformy {
			os := konfig.OsLubPlatforma(adres.Os)
			wpis.Axis, wpis.AxisId = &os, wskaznikTekstu(adres.KluczOsi)
		}
		wykaz = append(wykaz, wpis)
	}
	return wykaz
}

func rozstrzygniecieKatalogu(k shared.SessionConfig) shared.WorkingDirectoryResolution {
	rozstrzygniecie := shared.WorkingDirectoryResolution{
		Reason: shared.WorkingDirectoryDegradationNone,
	}
	if k.WorkingDirectory == nil {
		return rozstrzygniecie
	}
	zadany := wartoscTekstu(k.WorkingDirectory.RequestedPath)
	zastepczy := wartoscTekstu(k.WorkingDirectory.FallbackPath)
	rozstrzygniecie.RequestedPath = wskaznikTekstu(zadany)
	if zadany != "" {
		rozstrzygniecie.EffectivePath = zadany
		return rozstrzygniecie
	}
	if zastepczy != "" {
		rozstrzygniecie.EffectivePath = zastepczy
		rozstrzygniecie.Reason = shared.WorkingDirectoryDegradationFallbackUsed
		rozstrzygniecie.Detail = wskaznikTekstu(
			"konfiguracja nie wskazuje katalogu roboczego — obowiązuje katalog zastępczy")
	}
	return rozstrzygniecie
}
