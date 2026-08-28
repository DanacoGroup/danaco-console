// Rodzina `isolation.*` obsługuje adapter okna konfiguracji punktów izolacji
// i pięć komend macierzy przełączników: `scope.list`, `context.get/set`
// i `technical.get/set`.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Wartości punktów izolacji spisane ze stałych nieeksportowanych rejestru
// definicji; zapis musi nieść dokładnie tę wartość, którą czyta egzekutor.
const (
	wartoscKontekstOdrebny = "odrebna"
	wartoscZakresWylaczony = "wylaczony"
)

// punktIzolacji wiąże wartość wyliczenia kontraktu z kluczem ustawienia oraz
// z rodzajem punktu, kontekstowym albo technicznym.
type punktIzolacji struct {
	klucz string
	// kontekstowy odróżnia wymiar kontekstu od zakresu technicznego.
	kontekstowy bool
}

// punktyKontekstu wylicza trzy wymiary izolacji kontekstu w kolejności
// kontraktu. Kolejność jest stała, więc dwa wywołania dają ten sam wynik.
var punktyKontekstu = []struct {
	rodzaj shared.IsolationContextKind
	punkt  punktIzolacji
}{
	{shared.IsolationContextKindHistory, punktIzolacji{konfig.KluczIzolacjaHistoria, true}},
	{shared.IsolationContextKindMemory, punktIzolacji{konfig.KluczIzolacjaPamiec, true}},
	{shared.IsolationContextKindContext, punktIzolacji{konfig.KluczIzolacjaKontekst, true}},
}

// punktyTechniczne wylicza osiem zakresów technicznych izolacji w kolejności
// kontraktu, wiążąc każdy z kluczem ustawienia właściwym temu zakresowi.
var punktyTechniczne = []struct {
	zakres shared.IsolationTechnicalScope
	punkt  punktIzolacji
}{
	{shared.IsolationTechnicalScopeWorkingDirectory, punktIzolacji{konfig.KluczIzolacjaKatalogRoboczy, false}},
	{shared.IsolationTechnicalScopeProcessEnvironment, punktIzolacji{konfig.KluczIzolacjaSrodowiskoProcesu, false}},
	{shared.IsolationTechnicalScopeModelDataDirectory, punktIzolacji{konfig.KluczIzolacjaKatalogDanych, false}},
	{shared.IsolationTechnicalScopeNetworkAccess, punktIzolacji{konfig.KluczIzolacjaDostepSieciowy, false}},
	{shared.IsolationTechnicalScopeFileAccess, punktIzolacji{konfig.KluczIzolacjaPliki, false}},
	{shared.IsolationTechnicalScopeAccountToken, punktIzolacji{konfig.KluczIzolacjaKontoIToken, false}},
	{shared.IsolationTechnicalScopeProcessModel, punktIzolacji{konfig.KluczIzolacjaModelProcesu, false}},
	{shared.IsolationTechnicalScopeExecutionServer, punktIzolacji{konfig.KluczIzolacjaSerwerWykonania, false}},
}

// adapterIzolacji wypełnia port Izolacja trzema zależnościami: ustawienia
// trzymają wartości punktów, profile trzymają szablony i wybór warstwy,
// a rozstrzygacz odpowiada na pytanie o politykę obowiązującą po zasięgu.
type adapterIzolacji struct {
	ustawienia   dane.RepozytoriumKonfiguracjiOsi
	profile      dane.RepozytoriumIzolacji
	rozstrzygacz *konfig.Rozstrzygacz
	rozgloszenie *emiter
}

// nowyAdapterIzolacji wiąże port z zestawem repozytoriów i rozstrzygaczem.
// Zestaw niewpięty daje adapter odmawiający z kodem `internal_error`.
func nowyAdapterIzolacji(repozytoria *dane.Zestaw, rozstrzygacz *konfig.Rozstrzygacz) *adapterIzolacji {
	adapter := &adapterIzolacji{rozstrzygacz: rozstrzygacz}
	if repozytoria != nil {
		adapter.ustawienia = repozytoria.KonfiguracjaOsi
		adapter.profile = repozytoria.ProfileIzolacji()
	}
	return adapter
}

// ZRozgloszeniem dokłada nadajnik zdarzeń. Zapis punktu izolacji zmienia wiersz
// tabeli `ustawienie`, więc idzie tym samym zdarzeniem `config.changed`, co
// zapis rodziny `config.*`. Nadajnik niepodłączony nie wstrzymuje zapisu.
func (a *adapterIzolacji) ZRozgloszeniem(nadajnik Nadajnik) *adapterIzolacji {
	a.rozgloszenie = nowyEmiter(nadajnik)
	return a
}

// ── isolation.scope.list ─────────────────────────────────────────────────────

// PoziomyZasiegu obsługuje `isolation.scope.list` i zwraca osiem poziomów
// zasięgu w kolejności rozstrzygania, wziętych z tabeli `poziom_zasiegu`.
func (a *adapterIzolacji) PoziomyZasiegu(ctx context.Context,
	_ shared.IsolationScopeListRequest) (shared.IsolationScopeListResponse, error) {

	if a.profile == nil {
		return shared.IsolationScopeListResponse{}, bladPodlozaIzolacji("słownik poziomów zasięgu")
	}
	poziomy, err := a.profile.PoziomyZasiegu(ctx)
	if err != nil {
		return shared.IsolationScopeListResponse{}, err
	}
	if len(poziomy) == 0 {
		return shared.IsolationScopeListResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"izolacja: słownik poziomów zasięgu jest pusty — baza nie niesie wierszy poziom_zasiegu"))
	}
	najwezszy := poziomy[len(poziomy)-1].Pierwszenstwo
	wykaz := make([]shared.IsolationScopeLevel, 0, len(poziomy))
	for _, poziom := range poziomy {
		czyNajwezszy := poziom.Pierwszenstwo == najwezszy
		wykaz = append(wykaz, shared.IsolationScopeLevel{
			Scope: poziom.Poziom,
			Label: poziom.Nazwa,
			Order: poziom.Pierwszenstwo,
			// `Narrowest` wraca też wtedy, gdy jest fałszem — to jest odpowiedź.
			Narrowest: wskaznikPrawdy(czyNajwezszy),
		})
	}
	return shared.IsolationScopeListResponse{Scopes: wykaz}, nil
}

// ── isolation.context.get / isolation.context.set ────────────────────────────

// IzolacjaKontekstu obsługuje `isolation.context.get` i zwraca trzy
// przełączniki izolacji kontekstu spod adresu, zawsze w wykazie pełnym.
func (a *adapterIzolacji) IzolacjaKontekstu(ctx context.Context,
	z shared.IsolationContextGetRequest) (shared.IsolationContextGetResponse, error) {

	adres, err := a.adresIzolacji(z.Scope, z.ScopeId, z.Layer)
	if err != nil {
		return shared.IsolationContextGetResponse{}, err
	}
	zapisane, err := a.zapisyAdresu(ctx, adres)
	if err != nil {
		return shared.IsolationContextGetResponse{}, err
	}
	return shared.IsolationContextGetResponse{Switches: a.przelacznikiKontekstu(zapisane)}, nil
}

// ZapiszIzolacjeKontekstu obsługuje `isolation.context.set` i zapisuje
// wskazane przełączniki kontekstu pod adresem, oddając stan po zapisie
// odczytany z bazy, nie odbity z żądania.
func (a *adapterIzolacji) ZapiszIzolacjeKontekstu(ctx context.Context,
	z shared.IsolationContextSetRequest) (shared.IsolationContextSetResponse, error) {

	adres, err := a.adresIzolacji(z.Scope, z.ScopeId, z.Layer)
	if err != nil {
		return shared.IsolationContextSetResponse{}, err
	}
	if len(z.Switches) == 0 {
		return shared.IsolationContextSetResponse{}, bladZadaniaIzolacji(
			"zapis izolacji kontekstu bez ani jednego przełącznika")
	}
	for _, przelacznik := range z.Switches {
		punkt, znany := punktKontekstu(przelacznik.Kind)
		if !znany {
			return shared.IsolationContextSetResponse{}, bladNieznanegoPunktu(string(przelacznik.Kind))
		}
		if err := a.zapiszPunkt(ctx, adres, punkt, przelacznik.Isolated); err != nil {
			return shared.IsolationContextSetResponse{}, err
		}
	}
	zapisane, err := a.zapisyAdresu(ctx, adres)
	if err != nil {
		return shared.IsolationContextSetResponse{}, err
	}
	a.rozglosPolitykeAdresu(adres)
	return shared.IsolationContextSetResponse{Switches: a.przelacznikiKontekstu(zapisane)}, nil
}

// ── isolation.technical.get / isolation.technical.set ────────────────────────

// IzolacjaTechniczna zwraca osiem przełączników zakresów technicznych spod
// adresu, na tych samych zasadach co wykaz kontekstu.
func (a *adapterIzolacji) IzolacjaTechniczna(ctx context.Context,
	z shared.IsolationTechnicalGetRequest) (shared.IsolationTechnicalGetResponse, error) {

	adres, err := a.adresIzolacji(z.Scope, z.ScopeId, z.Layer)
	if err != nil {
		return shared.IsolationTechnicalGetResponse{}, err
	}
	zapisane, err := a.zapisyAdresu(ctx, adres)
	if err != nil {
		return shared.IsolationTechnicalGetResponse{}, err
	}
	return shared.IsolationTechnicalGetResponse{Switches: a.przelacznikiTechniczne(zapisane)}, nil
}

// ZapiszIzolacjeTechniczna obsługuje `isolation.technical.set` i zapisuje
// wskazane zakresy techniczne pod adresem, na tych samych zasadach co zapis
// kontekstu.
func (a *adapterIzolacji) ZapiszIzolacjeTechniczna(ctx context.Context,
	z shared.IsolationTechnicalSetRequest) (shared.IsolationTechnicalSetResponse, error) {

	adres, err := a.adresIzolacji(z.Scope, z.ScopeId, z.Layer)
	if err != nil {
		return shared.IsolationTechnicalSetResponse{}, err
	}
	if len(z.Switches) == 0 {
		return shared.IsolationTechnicalSetResponse{}, bladZadaniaIzolacji(
			"zapis izolacji technicznej bez ani jednego przełącznika")
	}
	for _, przelacznik := range z.Switches {
		punkt, znany := punktTechniczny(przelacznik.Scope)
		if !znany {
			return shared.IsolationTechnicalSetResponse{}, bladNieznanegoPunktu(string(przelacznik.Scope))
		}
		if err := a.zapiszPunkt(ctx, adres, punkt, przelacznik.Isolated); err != nil {
			return shared.IsolationTechnicalSetResponse{}, err
		}
	}
	zapisane, err := a.zapisyAdresu(ctx, adres)
	if err != nil {
		return shared.IsolationTechnicalSetResponse{}, err
	}
	a.rozglosPolitykeAdresu(adres)
	return shared.IsolationTechnicalSetResponse{Switches: a.przelacznikiTechniczne(zapisane)}, nil
}

// ── adres, odczyt i zapis punktu ─────────────────────────────────────────────

// adresIzolacji składa adres zapisu z pól żądania i sprawdza go w całości.
// Byt poziomu jest wymagany poza poziomem globalnym; warstwę rozstrzyga
// `warstwaAdresu`.
func (a *adapterIzolacji) adresIzolacji(poziom shared.ConfigScope, bytPoziomu *string,
	warstwa *shared.IsolationLayer) (konfig.Adres, error) {

	if !konfig.Znany(poziom) {
		return konfig.Adres{}, bladZadaniaIzolacji("poziom zasięgu " + string(poziom) +
			" nie należy do ośmiu poziomów kontraktu")
	}
	byt := wartoscTekstu(bytPoziomu)
	if poziom == shared.ConfigScopeGlobal {
		byt = ""
	} else if byt == "" {
		return konfig.Adres{}, bladZadaniaIzolacji("poziom zasięgu " + string(poziom) +
			" wymaga wskazania bytu — zapis bez niego nie doszedłby do rozstrzygania")
	}
	if _, err := warstwaAdresu(poziom, warstwa); err != nil {
		return konfig.Adres{}, err
	}
	return konfig.Adres{Poziom: poziom, KluczZasiegu: byt, Os: konfig.OsPlatformy}, nil
}

// warstwaAdresu rozstrzyga warstwę żądania — domyślną albo sesyjną — i sprawdza
// jej zgodność z poziomem zasięgu, bo warstwa sesyjna obowiązuje wyłącznie na
// poziomie karty sesji.
func warstwaAdresu(poziom shared.ConfigScope, warstwa *shared.IsolationLayer) (shared.IsolationLayer, error) {
	wybrana, err := warstwaZadania(warstwa)
	if err != nil {
		return "", err
	}
	if wybrana == shared.IsolationLayerSession && poziom != shared.ConfigScopeSession {
		return "", bladZadaniaIzolacji("warstwa karty sesji obowiązuje na poziomie " +
			shared.ConfigScopeSession + ", nie na poziomie " + string(poziom))
	}
	return wybrana, nil
}

// warstwaZadania rozpoznaje warstwę podaną w żądaniu. Warstwa pominięta znaczy
// domyślną platformy; wartość spoza wyliczenia kontraktu jest odmową, nie cichym
// zejściem na domyślną — inaczej literówka wyglądałaby na zapis wykonany pod
// wskazanym adresem.
func warstwaZadania(warstwa *shared.IsolationLayer) (shared.IsolationLayer, error) {
	if warstwa == nil || *warstwa == "" {
		return shared.IsolationLayerDefault, nil
	}
	switch *warstwa {
	case shared.IsolationLayerDefault, shared.IsolationLayerSession:
		return *warstwa, nil
	default:
		return "", bladZadaniaIzolacji("warstwa izolacji " + string(*warstwa) + " nie należy do kontraktu")
	}
}

// zapisyAdresu czyta wartości punktów izolacji zapisane pod adresem. Klucz bez
// wiersza po prostu nie wchodzi do mapy — brak zapisu znaczy wartość domyślną.
func (a *adapterIzolacji) zapisyAdresu(ctx context.Context, adres konfig.Adres) (map[string]string, error) {
	if a.ustawienia == nil {
		return nil, bladPodlozaIzolacji("magazyn ustawień")
	}
	wiersze, err := a.ustawienia.ListaOsi(ctx, adres.Poziom, adres.KluczZasiegu, adres.Os, adres.KluczOsi)
	if err != nil {
		return nil, err
	}
	zapisane := make(map[string]string, len(wiersze))
	for _, wiersz := range wiersze {
		zapisane[wiersz.Klucz] = wartoscTekstu(wiersz.Wartosc)
	}
	return zapisane, nil
}

// zapiszPunkt utrwala jeden punkt izolacji pod adresem i rozgłasza zmianę
// zdarzeniem `config.changed`, gdy nadajnik jest podłączony.
func (a *adapterIzolacji) zapiszPunkt(ctx context.Context, adres konfig.Adres,
	punkt punktIzolacji, odciety bool) error {

	if a.ustawienia == nil {
		return bladPodlozaIzolacji("magazyn ustawień")
	}
	wartosc := wartoscPunktu(punkt, odciety)
	ustawienie := dane.Ustawienie{
		Poziom: adres.Poziom, KluczZasiegu: adres.KluczZasiegu,
		Os: adres.Os, KluczOsi: adres.KluczOsi,
		Klucz: punkt.klucz, Wartosc: &wartosc, RodzajWartosci: string(konfig.RodzajTekst),
	}
	if err := a.ustawienia.Ustaw(ctx, ustawienie); err != nil {
		return err
	}
	if a.rozgloszenie != nil {
		a.rozgloszenie.ustawienie(shared.ChangeKindUpdated, wpisOsi(ustawienie))
	}
	return nil
}

// wartoscPunktu przekłada rozstrzygnięcie kontraktu (`isolated`) na wartość
// ustawienia — dokładnie tę, którą czyta `session.ZasadyZPolityki`.
func wartoscPunktu(punkt punktIzolacji, odciety bool) string {
	if punkt.kontekstowy {
		if odciety {
			return wartoscKontekstOdrebny
		}
		return konfig.IzolacjaWspoldzielona
	}
	if odciety {
		return konfig.IzolacjaWlaczona
	}
	return wartoscZakresWylaczony
}

// odcietyPunkt odpowiada, czy punkt jest odcięty, regułą przepisaną
// z egzekutora co do znaku, osobną dla wymiaru kontekstu i dla zakresu
// technicznego.
func odcietyPunkt(punkt punktIzolacji, wartosc string) bool {
	if punkt.kontekstowy {
		return wartosc != konfig.IzolacjaWspoldzielona
	}
	return wartosc == konfig.IzolacjaWlaczona
}

// przelacznikiKontekstu składa trzy przełączniki kontekstu z zapisów adresu,
// każdy wraz z objaśnieniem punktu wziętym z rejestru definicji.
func (a *adapterIzolacji) przelacznikiKontekstu(zapisane map[string]string) []shared.IsolationSwitch {
	przelaczniki := make([]shared.IsolationSwitch, 0, len(punktyKontekstu))
	for _, pozycja := range punktyKontekstu {
		przelaczniki = append(przelaczniki, shared.IsolationSwitch{
			Kind:        pozycja.rodzaj,
			Isolated:    odcietyPunkt(pozycja.punkt, a.wartoscObowiazujaca(pozycja.punkt, zapisane)),
			Explanation: a.objasnieniePunktu(pozycja.punkt),
		})
	}
	return przelaczniki
}

// przelacznikiTechniczne składa osiem przełączników technicznych z zapisów
// adresu, każdy wraz z objaśnieniem punktu wziętym z rejestru definicji.
func (a *adapterIzolacji) przelacznikiTechniczne(zapisane map[string]string) []shared.IsolationTechnicalSwitch {
	przelaczniki := make([]shared.IsolationTechnicalSwitch, 0, len(punktyTechniczne))
	for _, pozycja := range punktyTechniczne {
		przelaczniki = append(przelaczniki, shared.IsolationTechnicalSwitch{
			Scope:       pozycja.zakres,
			Isolated:    odcietyPunkt(pozycja.punkt, a.wartoscObowiazujaca(pozycja.punkt, zapisane)),
			Explanation: a.objasnieniePunktu(pozycja.punkt),
		})
	}
	return przelaczniki
}

// wartoscObowiazujaca zwraca wartość zapisaną pod adresem, a przy jej braku —
// wartość domyślną z rejestru definicji.
func (a *adapterIzolacji) wartoscObowiazujaca(punkt punktIzolacji, zapisane map[string]string) string {
	if wartosc, jest := zapisane[punkt.klucz]; jest && wartosc != "" {
		return wartosc
	}
	definicja, _ := a.rozstrzygacz.Rejestr().Definicja(punkt.klucz)
	return definicja.Domyslna
}

// objasnieniePunktu bierze opis kontrolki [?] z rejestru definicji. Rdzeń nie ma
// własnego zdania o punkcie izolacji — objaśnienie ma jedno źródło.
func (a *adapterIzolacji) objasnieniePunktu(punkt punktIzolacji) *string {
	definicja, jest := a.rozstrzygacz.Rejestr().Definicja(punkt.klucz)
	if !jest || definicja.Objasnienie == "" {
		return nil
	}
	objasnienie := definicja.Objasnienie
	return &objasnienie
}

// punktKontekstu rozpoznaje wymiar kontekstu podany w żądaniu i oddaje jego
// punkt izolacji albo informację, że wymiar nie jest znany.
func punktKontekstu(rodzaj shared.IsolationContextKind) (punktIzolacji, bool) {
	for _, pozycja := range punktyKontekstu {
		if pozycja.rodzaj == rodzaj {
			return pozycja.punkt, true
		}
	}
	return punktIzolacji{}, false
}

// punktTechniczny rozpoznaje zakres techniczny podany w żądaniu i oddaje jego
// punkt izolacji albo informację, że zakres nie jest znany.
func punktTechniczny(zakres shared.IsolationTechnicalScope) (punktIzolacji, bool) {
	for _, pozycja := range punktyTechniczne {
		if pozycja.zakres == zakres {
			return pozycja.punkt, true
		}
	}
	return punktIzolacji{}, false
}

// bladZadaniaIzolacji odmawia żądania niezgodnego z kontraktem albo ze
// schematem. To odmowa jednego wywołania, nie awaria rdzenia.
func bladZadaniaIzolacji(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"izolacja: "+powod))
}

// bladNieznanegoPunktu odmawia zapisu punktu izolacji spoza wyliczeń
// kontraktu, wskazanego rodzajem albo zakresem, którego rdzeń nie zna.
func bladNieznanegoPunktu(nazwa string) error {
	return bladZadaniaIzolacji("punkt izolacji " + nazwa + " nie należy do kontraktu")
}

// bladPodlozaIzolacji odmawia obsługi, gdy zależność portu nie została wpięta.
// Brak podłoża jest awarią montażu i ma być słyszalny, a nie zamieniony w pustą
// odpowiedź.
func bladPodlozaIzolacji(czego string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"izolacja: "+czego+" nie został wpięty do rdzenia"))
}
