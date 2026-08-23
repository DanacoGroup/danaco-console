// Odpowiedzialność pliku: złożenie konfiguracji sesji obowiązującej — czyli
// rozstrzygnięcie obszarów po ośmiu poziomach zasięgu i trzech osiach
// wraz ze wskazaniem, skąd wzięty jest każdy obszar.
//
// Determinizm. Kolejność adresów bierze pakiet konfig (Kontekst.Adresy) —
// najpierw poziom, w ramach poziomu oś. Wygrywa pierwszy adres, pod którym
// obszar jest zapisany, a obszary wychodzą w kolejności pól kontraktu. Ta sama
// zawartość rejestru daje więc zawsze ten sam wynik; rdzeń nie ma tu ani
// jednego rozstrzygnięcia zależnego od kolejności mapy.
//
// Pierwszeństwo należy do pakietu konfig. Rdzeń nie zna kolejności
// poziomów ani osi i jej nie powtarza; pyta o wykaz adresów i czyta po kolei.
//
// Granica katalogu roboczego. Rozejście „ustawiony vs faktycznie używany”
// składane jest z samej konfiguracji, bez dotykania dysku: odczyt konfiguracji
// obowiązującej ma być powtarzalny i nie ma prawa zakładać katalogów. Sprawdzenie
// dysku należy do core.KatalogRoboczy na drodze uruchomienia tury,
// dlatego pole checkedAt pozostaje zerowe — nic tu nie było sprawdzane.
package core

import (
	"context"
	"encoding/json"
	"time"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// KonfiguracjaObowiazujaca rozstrzyga obszary po wszystkich adresach kontekstu
// i zwraca konfigurację wraz z pochodzeniem każdego obszaru. Brak zapisu na
// każdym z adresów daje obszar pusty o pochodzeniu „wartość domyślna”, nie
// odmowę.
func (a *adapterUstawienOsi) KonfiguracjaObowiazujaca(ctx context.Context,
	z shared.ConfigEffectiveGetRequest) (shared.ConfigEffectiveGetResponse, error) {

	obszary, err := obszaryZadania(z.Areas)
	if err != nil {
		return shared.ConfigEffectiveGetResponse{}, err
	}
	tresci, pochodzenie, err := a.rozstrzygnijObszary(ctx, kontekstKonfiguracjiSesji(z), obszary)
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

// KonfiguracjaSesjiOkna rozstrzyga konfigurację obowiązującą dla okna i jego
// sesji po wszystkich adresach zasięgu. Służy drodze tury: adapter
// rozmowy odczytuje nią obowiązującą konfigurację i tłumaczy jej obszary na
// wejście procesu przez warstwę injection (adapter_rozmowa_konfiguracja.go).
//
// Bez osi modelu i konta. Osie te są dopiero wynikiem tej konfiguracji (obszary
// model i account tłumaczą się na wybór modelu i konta wywołania), więc nie
// mogą wchodzić do jej rozstrzygania — inaczej powstałaby pętla „konto zależy od
// konfiguracji, która zależy od konta”. Brak zapisu na każdym adresie daje
// konfigurację pustą, nie odmowę.
func (a *adapterUstawienOsi) KonfiguracjaSesjiOkna(ctx context.Context,
	idOkna, idSesji string) (shared.SessionConfig, error) {

	if a == nil {
		return shared.SessionConfig{}, nil
	}
	kontekst := konfig.Kontekst{Okno: idOkna, KartaSesji: idSesji}
	tresci, _, err := a.rozstrzygnijObszary(ctx, kontekst, obszaryKontraktu)
	if err != nil {
		return shared.SessionConfig{}, err
	}
	return zlozKonfiguracje(tresci)
}

// rozstrzygnijObszary czyta adresy kontekstu od najwęższego i przypisuje
// obszarowi treść z pierwszego adresu, pod którym obszar jest zapisany.
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

// kontekstKonfiguracjiSesji buduje kontekst rozstrzygania z pól żądania. Okno
// i karta sesji wchodzą wprost, wskazany byt poziomu ląduje na swoim poziomie,
// a oś modelu albo konta na swojej osi. Pole puste znaczy, że poziom albo oś
// nie dotyczy tego wywołania i jest pomijana.
func kontekstKonfiguracjiSesji(z shared.ConfigEffectiveGetRequest) konfig.Kontekst {
	kontekst := konfig.Kontekst{
		Okno:       wartoscTekstu(z.WindowId),
		KartaSesji: wartoscTekstu(z.SessionId),
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

// pochodzeniaObszarow opisuje, skąd wzięty jest każdy żądany obszar. Obszar
// bez zapisu na żadnym adresie pochodzi z wartości domyślnej kontraktu —
// to jest odpowiedź pełna, nie brak odpowiedzi.
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

// rozstrzygniecieKatalogu zestawia katalog ustawiony z katalogiem, który
// wynika z konfiguracji. Katalog wskazany przez Operatora obowiązuje; jego brak
// oddaje pole katalogowi zastępczemu i mówi o tym wprost, zamiast milczeć.
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
