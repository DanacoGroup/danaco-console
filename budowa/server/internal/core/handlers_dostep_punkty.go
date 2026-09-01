// Plik wypełnia port PunktyDostepu katalogiem punktów z bazy, obsługując
// założenie, wykaz i zmianę punktu konfiguracji platformy.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekPunktuDostepu znakuje kod punktu nadany przez rdzeń. Kod jest
// identyfikatorem punktu w kontrakcie i w nadaniach, więc musi być trwały.
const przedrostekPunktuDostepu = "pd-"

// Zgodność adaptera z portem PunktyDostepu sprawdzana jest przy kompilacji,
// bez osobnego testu zgodności typów.
var _ PunktyDostepu = (*adapterPunktowDostepu)(nil)

// adapterPunktowDostepu wypełnia port PunktyDostepu tabelą `punkt_dostepu`,
// wraz z sejfem poświadczeń i próbą osiągalności.
type adapterPunktowDostepu struct {
	repozytorium dane.RepozytoriumPunktowDostepu
	sejf         SejfPoswiadczen
	// sprawdzenie bada osiągalność punktu, polem zamiast funkcją wolnostojącą,
	// by test mógł je podstawić.
	sprawdzenie ProbaPunktu
}

// nowyAdapterPunktowDostepu wiąże port z repozytorium katalogu i wbudowaną
// próbą sprawdzenia korzeni lokalnych.
func nowyAdapterPunktowDostepu(repozytorium dane.RepozytoriumPunktowDostepu) *adapterPunktowDostepu {
	return &adapterPunktowDostepu{repozytorium: repozytorium, sprawdzenie: probaKorzeniLokalnych}
}

// ZSejfem wpina magazyn sekretów. Bez niego punkt powstaje bez odwołania do
// poświadczenia, a podany sekret jest odrzucany w tym samym wywołaniu.
func (a *adapterPunktowDostepu) ZSejfem(sejf SejfPoswiadczen) *adapterPunktowDostepu {
	a.sejf = sejf
	return a
}

// Dodaj zakłada punkt dostępu wraz z jego korzeniami, odrzucając żądanie,
// które kształtowi punktu nie odpowiada.
func (a *adapterPunktowDostepu) Dodaj(ctx context.Context,
	z shared.AccessPointAddRequest) (shared.AccessPointAddResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AccessPointAddResponse{}, bladBrakuKatalogu("punktów dostępu")
	}
	if err := sprawdzKsztaltPunktu(z); err != nil {
		return shared.AccessPointAddResponse{}, err
	}
	urzadzenie, err := urzadzenieWiersza(z.DeviceId)
	if err != nil {
		return shared.AccessPointAddResponse{}, err
	}
	kod := nowyIdentyfikator(przedrostekPunktuDostepu)
	odwolanie, err := odwolaniePoswiadczenia(ctx, a.sejf, kod, z.Credential)
	if err != nil {
		return shared.AccessPointAddResponse{}, err
	}
	punkt := dane.PunktDostepu{
		Kod: kod, Nazwa: z.Name, Opis: wartoscTekstu(z.Description), Rodzaj: z.Kind,
		UrzadzenieID: urzadzenie, NazwaMostu: wartoscTekstu(z.BridgeName),
		PoswiadczenieOdwolanie: odwolanie, TrybDomyslny: z.DefaultMode,
		Stan: shared.AccessPointStatusUnknown, Aktywny: z.Enabled == nil || *z.Enabled,
		Korzenie: z.Roots,
	}
	rozlozAdresPunktu(&punkt, z.Host, z.Endpoint)
	if _, err := a.repozytorium.Dodaj(ctx, punkt); err != nil {
		// Punkt nie powstał, więc nikt już nie sięgnie po jego sekret; w sejfie
		// zostałby wpisem bez wiersza, którego nie zdejmie żadna komenda.
		if odwolanie != nil {
			usunPoswiadczenie(ctx, a.sejf, kod)
		}
		return shared.AccessPointAddResponse{}, odmowaUrzadzeniaPunktu(err)
	}
	zapisany, err := a.repozytorium.PoKodzie(ctx, kod)
	if err != nil {
		return shared.AccessPointAddResponse{}, err
	}
	return shared.AccessPointAddResponse{Point: punktKontraktu(zapisany)}, nil
}

// Wykaz zwraca katalog punktów zawężony rodzajem i urządzeniem, a odmowę —
// gdy repozytorium nie jest wpięte, bo pusty wykaz mówiłby o katalogu bez
// punktów, nie o katalogu, którego nie ma.
func (a *adapterPunktowDostepu) Wykaz(ctx context.Context,
	z shared.AccessPointListRequest) (shared.AccessPointListResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AccessPointListResponse{}, bladBrakuKatalogu("punktów dostępu")
	}
	wiersze, err := a.repozytorium.Lista(ctx, z.EnabledOnly != nil && *z.EnabledOnly)
	if err != nil {
		return shared.AccessPointListResponse{}, err
	}
	punkty := make([]shared.AccessPoint, 0, len(wiersze))
	for _, wiersz := range wiersze {
		punkt := punktKontraktu(wiersz)
		if z.Kind != nil && *z.Kind != "" && punkt.Kind != *z.Kind {
			continue
		}
		if z.DeviceId != nil && *z.DeviceId != "" && wartoscTekstu(punkt.DeviceId) != *z.DeviceId {
			continue
		}
		punkty = append(punkty, punkt)
	}
	return shared.AccessPointListResponse{Points: punkty}, nil
}

// Zmien zapisuje zmienione pola punktu. Pole pominięte zostaje bez zmiany —
// żądanie kontraktu opisuje różnicę, nie komplet wiersza.
func (a *adapterPunktowDostepu) Zmien(ctx context.Context,
	z shared.AccessPointUpdateRequest) (shared.AccessPointUpdateResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AccessPointUpdateResponse{}, bladBrakuKatalogu("punktów dostępu")
	}
	punkt, err := a.repozytorium.PoKodzie(ctx, z.AccessPointId)
	if err != nil {
		return shared.AccessPointUpdateResponse{}, bladWskazania(err, "punkt dostępu", z.AccessPointId)
	}
	odwolanie, err := odwolaniePoswiadczenia(ctx, a.sejf, punkt.Kod, z.Credential)
	if err != nil {
		return shared.AccessPointUpdateResponse{}, err
	}
	if err := zastosujZmianePunktu(&punkt, z, odwolanie); err != nil {
		return shared.AccessPointUpdateResponse{}, err
	}
	if err := a.repozytorium.Aktualizuj(ctx, punkt); err != nil {
		return shared.AccessPointUpdateResponse{}, odmowaUrzadzeniaPunktu(err)
	}
	zapisany, err := a.repozytorium.PoKodzie(ctx, punkt.Kod)
	if err != nil {
		return shared.AccessPointUpdateResponse{}, err
	}
	return shared.AccessPointUpdateResponse{Point: punktKontraktu(zapisany)}, nil
}

// zastosujZmianePunktu nakłada na wiersz pola wskazane w żądaniu i zwraca
// odmowę, gdy wskazanie urządzenia nie jest identyfikatorem katalogu.
func zastosujZmianePunktu(punkt *dane.PunktDostepu, z shared.AccessPointUpdateRequest,
	odwolanie *string) error {

	if z.Name != nil {
		punkt.Nazwa = *z.Name
	}
	if z.Description != nil {
		punkt.Opis = *z.Description
	}
	if z.DeviceId != nil {
		urzadzenie, err := urzadzenieWiersza(z.DeviceId)
		if err != nil {
			return err
		}
		punkt.UrzadzenieID = urzadzenie
	}
	if z.BridgeName != nil {
		punkt.NazwaMostu = *z.BridgeName
	}
	if z.Host != nil || z.Endpoint != nil {
		rozlozAdresPunktu(punkt, z.Host, z.Endpoint)
	}
	if z.Roots != nil {
		punkt.Korzenie = z.Roots
	}
	if z.DefaultMode != nil {
		punkt.TrybDomyslny = *z.DefaultMode
	}
	if z.Enabled != nil {
		punkt.Aktywny = *z.Enabled
	}
	if odwolanie != nil {
		punkt.PoswiadczenieOdwolanie = odwolanie
	}
	return nil
}
