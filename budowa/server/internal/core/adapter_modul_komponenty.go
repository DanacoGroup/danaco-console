// Rodzina component.*: kafle Strefy 2 Strony głównej, których targetId wskazuje byt
// magazynu modułowego. Komenda create zakłada ten byt; update i delete magazynu nie tykają.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const przedrostekKomponentu = "komp-"

type adapterKomponentow struct {
	rejestr    dane.RepozytoriumKomponentow
	projekty   dane.RepozytoriumPrzestrzeniRoboczej
	agenci     dane.RepozytoriumAgentow
	automatyki dane.RepozytoriumAutomatyk
	teraz      func() int64
}

func nowyAdapterKomponentow(rejestr dane.RepozytoriumKomponentow) *adapterKomponentow {
	return &adapterKomponentow{
		rejestr: rejestr,
		teraz:   func() int64 { return time.Now().UTC().UnixMilli() },
	}
}

// Profile asystenta nie mają w warstwie danych drogi zapisu, więc nie ma ich w sygnaturze.
func (a *adapterKomponentow) ZMagazynamiModulow(projekty dane.RepozytoriumPrzestrzeniRoboczej,
	agenci dane.RepozytoriumAgentow, automatyki dane.RepozytoriumAutomatyk) *adapterKomponentow {

	a.projekty = projekty
	a.agenci = agenci
	a.automatyki = automatyki
	return a
}

func (a *adapterKomponentow) Wykaz(ctx context.Context,
	z shared.ComponentListRequest) (shared.ComponentListResponse, error) {

	if err := a.sprawdzRejestr(); err != nil {
		return shared.ComponentListResponse{}, err
	}
	filtr := dane.FiltrKomponentow{DolaczWylaczone: wartoscLogiczna(z.IncludeDisabled)}
	if z.Kind != nil {
		rodzaj, err := rodzajKomponentu(string(*z.Kind))
		if err != nil {
			return shared.ComponentListResponse{}, err
		}
		filtr.Rodzaj = rodzaj
	}
	wiersze, err := a.rejestr.Komponenty(ctx, filtr)
	if err != nil {
		return shared.ComponentListResponse{}, bladKomponentu(err)
	}
	komponenty := make([]shared.Component, 0, len(wiersze))
	for _, wiersz := range wiersze {
		komponenty = append(komponenty, komponentKontraktu(wiersz))
	}
	return shared.ComponentListResponse{Components: komponenty}, nil
}

func (a *adapterKomponentow) Utworz(ctx context.Context,
	z shared.ComponentCreateRequest) (shared.ComponentCreateResponse, error) {

	if err := a.sprawdzRejestr(); err != nil {
		return shared.ComponentCreateResponse{}, err
	}
	rodzaj, err := rodzajKomponentu(string(z.Kind))
	if err != nil {
		return shared.ComponentCreateResponse{}, err
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.ComponentCreateResponse{}, bladWskazaniaKomponentu("żądanie bez nazwy komponentu")
	}
	bytDocelowy, err := a.zalozBytModulowy(ctx, rodzaj, nazwa, z.Description)
	if err != nil {
		return shared.ComponentCreateResponse{}, err
	}
	chwila := a.teraz()
	zalozony, err := a.rejestr.ZalozKomponent(ctx, dane.Komponent{
		Kod:            nowyIdentyfikator(przedrostekKomponentu),
		Rodzaj:         rodzaj,
		Nazwa:          nazwa,
		Opis:           z.Description,
		BytDocelowy:    &bytDocelowy,
		Czynny:         true,
		Konfiguracja:   string(z.Config),
		Utworzono:      chwila,
		Zaktualizowano: chwila,
	})
	if err != nil {
		return shared.ComponentCreateResponse{}, bladKomponentu(err)
	}
	return shared.ComponentCreateResponse{Component: komponentKontraktu(zalozony)}, nil
}

func (a *adapterKomponentow) Zmien(ctx context.Context,
	z shared.ComponentUpdateRequest) (shared.ComponentUpdateResponse, error) {

	if err := a.sprawdzRejestr(); err != nil {
		return shared.ComponentUpdateResponse{}, err
	}
	if strings.TrimSpace(z.ComponentId) == "" {
		return shared.ComponentUpdateResponse{}, bladWskazaniaKomponentu("żądanie bez komponentu")
	}
	if z.Name != nil && strings.TrimSpace(*z.Name) == "" {
		return shared.ComponentUpdateResponse{}, bladWskazaniaKomponentu("nowa nazwa komponentu jest pusta")
	}
	zmiana := dane.ZmianaKomponentu{Nazwa: z.Name, Opis: z.Description, Czynny: z.Enabled}
	if len(z.Config) > 0 {
		konfiguracja := string(z.Config)
		zmiana.Konfiguracja = &konfiguracja
	}
	zmieniony, err := a.rejestr.ZmienKomponent(ctx, z.ComponentId, zmiana, a.teraz())
	if err != nil {
		return shared.ComponentUpdateResponse{}, bladKomponentu(err)
	}
	return shared.ComponentUpdateResponse{Component: komponentKontraktu(zmieniony)}, nil
}

// Komponent, którego nie ma, jest odmową z powodem, nie odpowiedzią deleted: false.
func (a *adapterKomponentow) Usun(ctx context.Context,
	z shared.ComponentDeleteRequest) (shared.ComponentDeleteResponse, error) {

	if err := a.sprawdzRejestr(); err != nil {
		return shared.ComponentDeleteResponse{}, err
	}
	if strings.TrimSpace(z.ComponentId) == "" {
		return shared.ComponentDeleteResponse{}, bladWskazaniaKomponentu("żądanie bez komponentu")
	}
	usuniety, err := a.rejestr.UsunKomponent(ctx, z.ComponentId)
	if err != nil {
		return shared.ComponentDeleteResponse{}, bladKomponentu(err)
	}
	if !usuniety {
		return shared.ComponentDeleteResponse{}, bladNieznanegoKomponentu(z.ComponentId)
	}
	return shared.ComponentDeleteResponse{Deleted: true}, nil
}

// component.assign przyjmuje pięć poziomów zasięgu; pozostałe idą odmową validation_failed.
func (a *adapterKomponentow) Przypisz(ctx context.Context,
	z shared.ComponentAssignRequest) (shared.ComponentAssignResponse, error) {

	if err := a.sprawdzRejestr(); err != nil {
		return shared.ComponentAssignResponse{}, err
	}
	if strings.TrimSpace(z.ComponentId) == "" {
		return shared.ComponentAssignResponse{}, bladWskazaniaKomponentu("żądanie bez komponentu")
	}
	poziom, err := poziomPrzypisania(z.Scope)
	if err != nil {
		return shared.ComponentAssignResponse{}, err
	}
	klucz := strings.TrimSpace(wartoscTekstu(z.ScopeId))
	if z.Scope == shared.ConfigScopeGlobal {
		if klucz != "" {
			return shared.ComponentAssignResponse{}, bladWskazaniaKomponentu(
				"poziom global nie przyjmuje identyfikatora bytu poziomu")
		}
	} else if klucz == "" {
		return shared.ComponentAssignResponse{}, bladWskazaniaKomponentu(
			"poziom " + string(z.Scope) + " wymaga identyfikatora bytu poziomu")
	}
	przypisany, doszlo, err := a.rejestr.PrzypiszKomponent(ctx, z.ComponentId, poziom, klucz, a.teraz())
	if err != nil {
		return shared.ComponentAssignResponse{}, bladKomponentu(err)
	}
	return shared.ComponentAssignResponse{
		Component: komponentKontraktu(przypisany),
		Assigned:  doszlo,
	}, nil
}

func (a *adapterKomponentow) zalozBytModulowy(ctx context.Context, rodzaj, nazwa string,
	opis *string) (string, error) {

	switch rodzaj {
	case shared.ComponentKindWorkspace:
		if a.projekty == nil {
			return "", bladBrakuMagazynuKomponentu(rodzaj, "magazyn projektów nie jest wpięty")
		}
		kod := nowyIdentyfikator(przedrostekProjektu)
		if _, _, err := a.projekty.ZapewnijProjekt(ctx, kod, nazwa); err != nil {
			return "", bladKomponentu(err)
		}
		return kod, nil

	case shared.ComponentKindAgents:
		if a.agenci == nil {
			return "", bladBrakuMagazynuKomponentu(rodzaj, "biblioteka ekspertów nie jest wpięta")
		}
		kod := nowyIdentyfikator(przedrostekEksperta)
		if _, err := a.agenci.Dodaj(ctx, dane.Agent{
			Kod: kod, Nazwa: nazwa, Opis: wartoscTekstu(opis), Aktywny: true,
		}); err != nil {
			return "", bladKomponentu(err)
		}
		return kod, nil

	case shared.ComponentKindAutomations:
		if a.automatyki == nil {
			return "", bladBrakuMagazynuKomponentu(rodzaj, "magazyn automatyk nie jest wpięty")
		}
		kod := nowyIdentyfikator(przedrostekAutomatyki)
		if _, err := a.automatyki.ZapiszAutomatyke(ctx, dane.Automatyka{
			Kod: kod, Nazwa: nazwa, Opis: opis, Czynna: true,
		}); err != nil {
			return "", bladKomponentu(err)
		}
		return kod, nil

	case shared.ComponentKindAssistant:
		// Warstwa danych oddaje profil asystenta wyłącznie do odczytu.
		return "", bladBrakuMagazynuKomponentu(rodzaj,
			"warstwa danych nie ma zapisu profili asystenta")
	}
	return "", bladNieznanegoRodzajuKomponentu(rodzaj)
}

func komponentKontraktu(k dane.Komponent) shared.Component {
	komponent := shared.Component{
		Id:        k.Kod,
		Kind:      shared.ComponentKind(k.Rodzaj),
		Name:      k.Nazwa,
		Enabled:   k.Czynny,
		CreatedAt: k.Utworzono,
		UpdatedAt: k.Zaktualizowano,
	}
	komponent.Description = k.Opis
	komponent.TargetId = k.BytDocelowy
	if k.Konfiguracja != "" {
		komponent.Config = json.RawMessage(k.Konfiguracja)
	}
	return komponent
}

// Kolumna `komponent.rodzaj` niesie wartość kontraktu wprost — sprawdzenie, nie przekład.
func rodzajKomponentu(rodzaj string) (string, error) {
	switch rodzaj {
	case shared.ComponentKindAutomations, shared.ComponentKindAgents,
		shared.ComponentKindWorkspace, shared.ComponentKindAssistant:
		return rodzaj, nil
	}
	return "", bladNieznanegoRodzajuKomponentu(rodzaj)
}

// Poziom kontraktu przekłada się na kod kolumny poziom_zasiegu.kod słownikiem platformy.
func poziomPrzypisania(poziom shared.ConfigScope) (string, error) {
	switch poziom {
	case shared.ConfigScopeGlobal, shared.ConfigScopeEnvironment, shared.ConfigScopeProject,
		shared.ConfigScopeSession, shared.ConfigScopeWindow:
		kod, jest := shared.WartosciBazyConfigScope[poziom]
		if !jest {
			return "", bladWskazaniaKomponentu("poziom zasięgu " + string(poziom) +
				" nie ma odpowiednika w słowniku bazy")
		}
		return kod, nil
	}
	return "", bladWskazaniaKomponentu("poziom zasięgu " + string(poziom) +
		" nie jest wymieniony w opisie component.assign; przyjmowane są global, environment, project, session i window")
}

func (a *adapterKomponentow) sprawdzRejestr() error {
	if a == nil || a.rejestr == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"komponenty własne: rejestr komponentów nie jest wpięty"))
	}
	return nil
}

func bladKomponentu(err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"komponenty własne: "+err.Error()))
	}
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
			"komponenty własne: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"komponenty własne: "+err.Error()))
}

func bladWskazaniaKomponentu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"komponenty własne: "+powod))
}

func bladNieznanegoRodzajuKomponentu(rodzaj string) error {
	return bladWskazaniaKomponentu("rodzaj komponentu " + rodzaj + " nie należy do kontraktu")
}

func bladNieznanegoKomponentu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"komponenty własne: komponent "+kod+" nie istnieje"))
}

// Kod `conflict`: żądanie jest poprawne, lecz stan platformy wyklucza czynność.
func bladBrakuMagazynuKomponentu(rodzaj, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"komponenty własne: komponentu rodzaju "+rodzaj+" nie da się założyć — "+powod))
}
