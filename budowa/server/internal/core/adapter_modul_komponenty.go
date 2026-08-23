// Odpowiedzialność pliku: adapter rodziny `component.*` — komponentów własnych
// Strefy 2 Strony głównej. Uchwyty i port stoją w `handlers_komponenty.go`.
//
// Rodzina stoi na własnym rejestrze komponentów, który wskazuje byt modułowy:
// komponent własny jest kaflem Strefy 2, a `Component.targetId` niesie
// identyfikator bytu magazynu modułowego, który ten kafel reprezentuje. Czystej
// fasady nad samymi magazynami modułowymi zbudować się nie da: `config` nie ma
// kolumny w żadnej z tabel modułowych, para (poziom zasięgu, klucz) z
// `component.assign` nie ma gdzie usiąść, a `component.list` nie miałby
// wspólnego porządku wyświetlania.
//
// `create` sięga do magazynu modułowego, `update` i `delete` — nie:
//
//   - `component.create` zakłada komponent w magazynie właściwym jego rodzajowi,
//     a żądanie nie niesie `targetId`, więc byt docelowy powstaje tutaj: projekt
//     przez `ZapewnijProjekt`, ekspert przez `Dodaj`, automatyka przez
//     `ZapiszAutomatyke`.
//   - `component.update` zmienia kafel, nie byt magazynu. Nazwa eksperta ma swoją
//     komendę (`agent.update`), definicja automatyki swoją
//     (`automation.workflow.save`) — pisanie tych kolumn także stąd dałoby dwóch
//     pisarzy jednej kolumny.
//   - `component.delete` zdejmuje kafel; byt modułowy zostaje pod swoim
//     identyfikatorem. `DELETE` istnieje zresztą wyłącznie dla tabeli `agent`
//     (`dane/agenci_zapis.go`), a `projekt` i `automatyka` nie mają go w warstwie
//     danych.
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

// przedrostekKomponentu nadaje identyfikator kaflowi Strefy 2. Byt docelowy
// dostaje przedrostek właściwy swojemu magazynowi, a nie ten: ekspert `ag-`
// (przedrostekEksperta), automatyka `automat-` (przedrostekAutomatyki), projekt
// `prj-` (przedrostekProjektu). Stałe modułów są używane tu wprost, żeby klucz
// założony przez `component.create` wyglądał dokładnie tak jak klucz założony
// komendą własną modułu.
const przedrostekKomponentu = "komp-"

// adapterKomponentow wypełnia port Komponenty. Rejestr jest zależnością
// obowiązkową; trzy repozytoria modułowe są zależnościami `component.create`
// i mogą być puste — wtedy założenie komponentu tego rodzaju odmawia z powodem,
// zamiast zakładać kafel wskazujący na nic.
type adapterKomponentow struct {
	rejestr    dane.RepozytoriumKomponentow
	projekty   dane.RepozytoriumPrzestrzeniRoboczej
	agenci     dane.RepozytoriumAgentow
	automatyki dane.RepozytoriumAutomatyk
	// teraz oddaje czas w milisekundach epoki. Wydzielone w pole, bo obie
	// kolumny czasu niosą wartość kontraktu wprost i muszą pochodzić z jednego
	// zegara.
	teraz func() int64
}

// nowyAdapterKomponentow wiąże adapter z rejestrem komponentów.
func nowyAdapterKomponentow(rejestr dane.RepozytoriumKomponentow) *adapterKomponentow {
	return &adapterKomponentow{
		rejestr: rejestr,
		teraz:   func() int64 { return time.Now().UTC().UnixMilli() },
	}
}

// ZMagazynamiModulow oddaje adapterowi trzy magazyny, w których `component.create`
// zakłada byt docelowy. Profile asystenta nie mają w warstwie danych drogi
// zapisu, więc nie ma ich też w tej sygnaturze.
func (a *adapterKomponentow) ZMagazynamiModulow(projekty dane.RepozytoriumPrzestrzeniRoboczej,
	agenci dane.RepozytoriumAgentow, automatyki dane.RepozytoriumAutomatyk) *adapterKomponentow {

	a.projekty = projekty
	a.agenci = agenci
	a.automatyki = automatyki
	return a
}

// Wykaz obsługuje `component.list`. Brak zawężenia rodzajem oddaje wszystkie
// cztery rodzaje; `includeDisabled` pominięte znaczy „tylko czynne”.
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

// Utworz obsługuje `component.create`: zakłada byt w magazynie właściwym
// rodzajowi, a potem kafel wskazujący na ten byt.
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

// Zmien obsługuje `component.update`. Zmienia wyłącznie kafel — patrz nagłówek.
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

// Usun obsługuje `component.delete`. Zdejmuje kafel; bytu magazynu modułowego
// nie tyka — patrz nagłówek. Komponent, którego nie ma, jest odmową z powodem,
// nie odpowiedzią `deleted: false` udającą wykonaną czynność.
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

// Przypisz obsługuje `component.assign`: zapamiętuje na wierszu komponentu, na
// którym poziomie zasięgu ten komponent obowiązuje. Nic ponad to — przypisanie
// nie przenosi bytu modułowego, nie nadaje uprawnień i nie włącza komponentu do
// żadnej pętli wykonania; żądanie niesie tylko komponent, poziom i identyfikator
// bytu poziomu, a zapis tej pary jest jedyną czynnością, którą te pola opisują.
//
// Z wartości `ConfigScope` przyjmowanych jest pięć, które wylicza opis komendy:
// global, environment, project, session i window. Pozostałe — `module`,
// `modulePair`, `role` i `application` — idą odmową `validation_failed`.
//
// `assigned: false` znaczy powtórzenie tego samego przypisania — nic nie doszło
// do skutku i wynik mówi to wprost.
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

// zalozBytModulowy zakłada byt w magazynie właściwym rodzajowi i oddaje jego
// identyfikator zewnętrzny. Rodzaj `assistant` odmawia — patrz niżej.
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
		// Warstwa danych oddaje profil asystenta wyłącznie do odczytu
		// (`dane/asystent_profil.go`) — nie ma drogi zapisu, którą dałoby się
		// tu założyć byt docelowy. Kafel bez `targetId` byłby sukcesem bez
		// skutku, więc rodzaj odmawia z powodem.
		return "", bladBrakuMagazynuKomponentu(rodzaj,
			"warstwa danych nie ma zapisu profili asystenta")
	}
	return "", bladNieznanegoRodzajuKomponentu(rodzaj)
}

// komponentKontraktu przekłada wiersz rejestru na kształt kontraktu.
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

// rodzajKomponentu przepuszcza wyłącznie cztery wartości `ComponentKind`.
// Kolumna `komponent.rodzaj` niesie wartość kontraktu wprost, więc przekładu
// tu nie ma — jest sprawdzenie.
func rodzajKomponentu(rodzaj string) (string, error) {
	switch rodzaj {
	case shared.ComponentKindAutomations, shared.ComponentKindAgents,
		shared.ComponentKindWorkspace, shared.ComponentKindAssistant:
		return rodzaj, nil
	}
	return "", bladNieznanegoRodzajuKomponentu(rodzaj)
}

// poziomPrzypisania przekłada poziom kontraktu na kod kolumny `poziom_zasiegu.kod`
// zastanym słownikiem platformy — i przepuszcza wyłącznie pięć poziomów, które
// wylicza opis komendy (patrz komentarz przy Przypisz).
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

// sprawdzRejestr odmawia czynności, gdy rejestr komponentów nie został wpięty.
// Bez tego sprawdzenia `component.list` oddałby pusty wykaz, a
// `component.delete` — `deleted: false`; jedno i drugie wyglądałoby jak stan
// platformy, a nie jak niezłożony port. Kod `internal_error`: żądanie jest
// poprawne, wina leży po stronie montażu.
func (a *adapterKomponentow) sprawdzRejestr() error {
	if a == nil || a.rejestr == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"komponenty własne: rejestr komponentów nie jest wpięty"))
	}
	return nil
}

// bladKomponentu przekłada niepowodzenie warstwy danych na odmowę kontraktu.
// Brak wiersza jest `not_found`, reszta — `internal_error`.
func bladKomponentu(err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"komponenty własne: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"komponenty własne: "+err.Error()))
}

// bladWskazaniaKomponentu odmawia żądaniu niezgodnemu z kontraktem.
func bladWskazaniaKomponentu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"komponenty własne: "+powod))
}

// bladNieznanegoRodzajuKomponentu odmawia rodzajowi spoza wyliczenia kontraktu.
func bladNieznanegoRodzajuKomponentu(rodzaj string) error {
	return bladWskazaniaKomponentu("rodzaj komponentu " + rodzaj + " nie należy do kontraktu")
}

// bladNieznanegoKomponentu odmawia czynności na kaflu, którego nie ma.
func bladNieznanegoKomponentu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"komponenty własne: komponent "+kod+" nie istnieje"))
}

// bladBrakuMagazynuKomponentu odmawia założeniu komponentu rodzaju, którego
// magazynu nie ma albo nie wpięto. Kod `conflict`: żądanie jest poprawne, lecz
// stan platformy wyklucza czynność.
func bladBrakuMagazynuKomponentu(rodzaj, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"komponenty własne: komponentu rodzaju "+rodzaj+" nie da się założyć — "+powod))
}
