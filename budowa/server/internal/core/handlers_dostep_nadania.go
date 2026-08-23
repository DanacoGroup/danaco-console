// Odpowiedzialność pliku: wypełnienie portu NadaniaDostepu zbiorem nadań okna
// rozmowy.
//
// Zbiór, nie pojedyncze nadanie. Każda czynność zapisu oddaje nie tylko wiersz
// zmieniony, ale komplet nadań okna po zmianie. Kolejność i oznaczenie głównego
// są własnością zbioru, więc dopisanie jednego nadania przestawia pozostałe;
// klient, który dostałby sam zmieniony wiersz, pokazałby zbiór nieprawdziwy.
//
// Granica korzeni. Zawężenie korzeni poza obszar punktu jest odmową
// merytoryczną (dane.ErrPozaKorzeniami), nie awarią zapisu — nadanie dostępu
// szerszego, niż punkt obiecuje, byłoby obejściem granicy uprawnień.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// przedrostekNadania znakuje identyfikator nadania nadany przez rdzeń.
const przedrostekNadania = "nd-"

// Zgodność adaptera z portem sprawdzana jest przy kompilacji.
var _ NadaniaDostepu = (*adapterNadanDostepu)(nil)

// adapterNadanDostepu wypełnia port NadaniaDostepu tabelą `nadanie_dostepu`.
type adapterNadanDostepu struct {
	nadania dane.RepozytoriumNadan
	okna    dane.RepozytoriumOkien
	punkty  dane.RepozytoriumPunktowDostepu
}

// nowyAdapterNadanDostepu wiąże port z trzema repozytoriami: nadania niosą
// numery wierszy okna i punktu, a kontrakt — ich identyfikatory trwałe.
func nowyAdapterNadanDostepu(nadania dane.RepozytoriumNadan, okna dane.RepozytoriumOkien,
	punkty dane.RepozytoriumPunktowDostepu) *adapterNadanDostepu {

	return &adapterNadanDostepu{nadania: nadania, okna: okna, punkty: punkty}
}

// Dodaj nadaje oknu dostęp do punktu i oddaje zbiór nadań okna po zapisie.
func (a *adapterNadanDostepu) Dodaj(ctx context.Context,
	z shared.AccessGrantAddRequest) (shared.AccessGrantAddResponse, error) {

	if !a.gotowy() {
		return shared.AccessGrantAddResponse{}, bladBrakuKatalogu("nadań dostępu")
	}
	okno, err := a.okna.PoIdentyfikatorze(ctx, z.WindowId)
	if err != nil {
		return shared.AccessGrantAddResponse{}, bladWskazania(err, "okno komunikacji", z.WindowId)
	}
	punkt, err := a.punkty.PoKodzie(ctx, z.AccessPointId)
	if err != nil {
		return shared.AccessGrantAddResponse{}, bladWskazania(err, "punkt dostępu", z.AccessPointId)
	}
	identyfikator := nowyIdentyfikator(przedrostekNadania)
	nadanie := dane.Nadanie{
		OknoKomunikacjiID: okno.ID, PunktDostepuID: punkt.ID, Tryb: z.Mode,
		Korzenie: z.Roots, Kolejnosc: liczbaLub(z.Order, 0),
		Glowne: z.Primary != nil && *z.Primary, Aktywne: true,
		IdentyfikatorZewnetrzny: &identyfikator,
	}
	if _, err := a.nadania.Dodaj(ctx, nadanie); err != nil {
		return shared.AccessGrantAddResponse{}, bladNadania(err)
	}
	zbior, err := a.zbiorOkna(ctx, okno, false)
	if err != nil {
		return shared.AccessGrantAddResponse{}, err
	}
	return shared.AccessGrantAddResponse{Grant: zeZbioru(zbior, identyfikator), Grants: zbior}, nil
}

// Wykaz zwraca nadania okna w kolejności zbioru.
func (a *adapterNadanDostepu) Wykaz(ctx context.Context,
	z shared.AccessGrantListRequest) (shared.AccessGrantListResponse, error) {

	pusty := shared.AccessGrantListResponse{Grants: []shared.AccessGrant{}}
	if !a.gotowy() || z.WindowId == nil || *z.WindowId == "" {
		return pusty, nil
	}
	okno, err := a.okna.PoIdentyfikatorze(ctx, *z.WindowId)
	if brakWiersza(err) {
		return pusty, nil
	}
	if err != nil {
		return shared.AccessGrantListResponse{}, err
	}
	zbior, err := a.zbiorOkna(ctx, okno, z.EnabledOnly != nil && *z.EnabledOnly)
	if err != nil {
		return shared.AccessGrantListResponse{}, err
	}
	nadania := make([]shared.AccessGrant, 0, len(zbior))
	for _, nadanie := range zbior {
		if z.AccessPointId != nil && *z.AccessPointId != "" && nadanie.AccessPointId != *z.AccessPointId {
			continue
		}
		nadania = append(nadania, nadanie)
	}
	return shared.AccessGrantListResponse{Grants: nadania}, nil
}

// Zmien zapisuje tryb, korzenie, kolejność albo oznaczenie głównego.
func (a *adapterNadanDostepu) Zmien(ctx context.Context,
	z shared.AccessGrantUpdateRequest) (shared.AccessGrantUpdateResponse, error) {

	if !a.gotowy() {
		return shared.AccessGrantUpdateResponse{}, bladBrakuKatalogu("nadań dostępu")
	}
	nadanie, err := a.nadania.PoIdentyfikatorze(ctx, z.GrantId)
	if err != nil {
		return shared.AccessGrantUpdateResponse{}, bladWskazania(err, "nadanie dostępu", z.GrantId)
	}
	zastosujZmianeNadania(&nadanie, z)
	if err := a.nadania.Aktualizuj(ctx, nadanie); err != nil {
		return shared.AccessGrantUpdateResponse{}, bladNadania(err)
	}
	if z.Primary != nil && *z.Primary {
		if err := a.nadania.OznaczGlowne(ctx, nadanie.ID); err != nil {
			return shared.AccessGrantUpdateResponse{}, err
		}
	}
	okno, err := a.okna.Pobierz(ctx, nadanie.OknoKomunikacjiID)
	if err != nil {
		return shared.AccessGrantUpdateResponse{}, err
	}
	zbior, err := a.zbiorOkna(ctx, okno, false)
	if err != nil {
		return shared.AccessGrantUpdateResponse{}, err
	}
	return shared.AccessGrantUpdateResponse{Grant: zeZbioru(zbior, z.GrantId), Grants: zbior}, nil
}

// Usun odbiera oknu nadanie. Nadanie nieznane nie jest błędem — wynik mówi
// wtedy, że nic nie odebrano.
func (a *adapterNadanDostepu) Usun(ctx context.Context,
	z shared.AccessGrantRemoveRequest) (shared.AccessGrantRemoveResponse, error) {

	pusty := shared.AccessGrantRemoveResponse{Removed: false, Grants: []shared.AccessGrant{}}
	if !a.gotowy() {
		return pusty, nil
	}
	nadanie, err := a.nadania.PoIdentyfikatorze(ctx, z.GrantId)
	if brakWiersza(err) {
		return pusty, nil
	}
	if err != nil {
		return shared.AccessGrantRemoveResponse{}, err
	}
	if err := a.nadania.Usun(ctx, nadanie.ID); err != nil {
		if brakWiersza(err) {
			return pusty, nil
		}
		return shared.AccessGrantRemoveResponse{}, err
	}
	okno, err := a.okna.Pobierz(ctx, nadanie.OknoKomunikacjiID)
	if err != nil {
		return shared.AccessGrantRemoveResponse{Removed: true, Grants: []shared.AccessGrant{}}, nil
	}
	zbior, err := a.zbiorOkna(ctx, okno, false)
	if err != nil {
		return shared.AccessGrantRemoveResponse{}, err
	}
	return shared.AccessGrantRemoveResponse{Removed: true, Grants: zbior}, nil
}

// gotowy odpowiada, czy adapter ma komplet repozytoriów potrzebnych do zapisu.
func (a *adapterNadanDostepu) gotowy() bool {
	return a != nil && a.nadania != nil && a.okna != nil && a.punkty != nil
}

// bladNadania przekłada odmowę granicy korzeni na kod kontraktu. Pozostałe
// błędy przechodzą nietknięte — są awarią zapisu, nie pomyłką wskazania.
func bladNadania(err error) error {
	if errors.Is(err, dane.ErrPozaKorzeniami) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, err.Error()))
	}
	return err
}
