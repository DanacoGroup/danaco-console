// Plik obsługuje siedem komend rodziny `isolation.*` stojących ponad
// pojedynczym przełącznikiem: pięć komend profilu, przełączenie warstwy oraz
// podgląd polityki obowiązującej.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// przedrostekProfiluIzolacji poprzedza identyfikator profilu zakładanego przez
// rdzeń, gdy żądanie zapisu nie wskazuje profilu istniejącego.
const przedrostekProfiluIzolacji = "profil-izolacji-"

// ZapiszProfilIzolacji zakłada profil albo zmienia istniejący.
//
// Żądanie z `profileId` jest zmianą profilu istniejącego, a nie założeniem bytu
// pod narzuconym identyfikatorem: profil, którego nie ma, wraca odmową
// `not_found`, zamiast po cichu powstać.
func (a *adapterIzolacji) ZapiszProfilIzolacji(ctx context.Context,
	z shared.IsolationProfileSaveRequest) (shared.IsolationProfileSaveResponse, error) {

	if a.profile == nil {
		return shared.IsolationProfileSaveResponse{}, bladPodlozaIzolacji("magazyn profili izolacji")
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.IsolationProfileSaveResponse{}, bladZadaniaIzolacji("profil bez nazwy")
	}
	kod := wartoscTekstu(z.ProfileId)
	if kod != "" {
		if _, err := a.profilZadania(ctx, kod); err != nil {
			return shared.IsolationProfileSaveResponse{}, err
		}
	} else {
		kod = nowyIdentyfikator(przedrostekProfiluIzolacji)
	}
	przelaczniki, err := przelacznikiProfilu(z.ContextSwitches, z.TechnicalSwitches)
	if err != nil {
		return shared.IsolationProfileSaveResponse{}, err
	}
	zapisany, err := a.profile.ZapiszProfil(ctx, dane.ProfilIzolacji{
		Kod:          kod,
		Nazwa:        nazwa,
		Opis:         strings.TrimSpace(wartoscTekstu(z.Description)),
		Przelaczniki: przelaczniki,
	})
	if err != nil {
		return shared.IsolationProfileSaveResponse{}, err
	}
	a.rozglosProfilIzolacji(ctx, rodzajZapisuProfilu(z.ProfileId), zapisany.Kod)
	return shared.IsolationProfileSaveResponse{Profile: a.profilKontraktu(zapisany)}, nil
}

// ProfileIzolacji zwraca profile izolacji, a przy wskazanym poziomie — profil
// przypisany do wskazanego bytu, wykaz zawężony do zera albo jednego elementu.
func (a *adapterIzolacji) ProfileIzolacji(ctx context.Context,
	z shared.IsolationProfileListRequest) (shared.IsolationProfileListResponse, error) {

	if a.profile == nil {
		return shared.IsolationProfileListResponse{}, bladPodlozaIzolacji("magazyn profili izolacji")
	}
	if z.Scope != nil && *z.Scope != "" {
		adres, err := a.adresIzolacji(*z.Scope, z.ScopeId, nil)
		if err != nil {
			return shared.IsolationProfileListResponse{}, err
		}
		return a.profilePoziomu(ctx, adres)
	}
	zapisane, err := a.profile.Profile(ctx)
	if err != nil {
		return shared.IsolationProfileListResponse{}, err
	}
	profile := make([]shared.IsolationProfile, 0, len(zapisane))
	for _, profil := range zapisane {
		profile = append(profile, a.profilKontraktu(profil))
	}
	return shared.IsolationProfileListResponse{Profiles: profile}, nil
}

// profilePoziomu zwraca profil przypisany pod adresem: pustą listę albo
// dokładnie jeden element wykazu.
func (a *adapterIzolacji) profilePoziomu(ctx context.Context,
	adres konfig.Adres) (shared.IsolationProfileListResponse, error) {

	kod, jest, err := a.profile.ProfilPoziomu(ctx, adres.Poziom, adres.KluczZasiegu)
	if err != nil {
		return shared.IsolationProfileListResponse{}, err
	}
	profile := make([]shared.IsolationProfile, 0, 1)
	if jest {
		profil, err := a.profilZadania(ctx, kod)
		if err != nil {
			return shared.IsolationProfileListResponse{}, err
		}
		profile = append(profile, a.profilKontraktu(profil))
	}
	return shared.IsolationProfileListResponse{Profiles: profile}, nil
}

// WczytajProfilIzolacji oddaje profil do panelu wskazany jego
// identyfikatorem, niczego nie przypisując i niczego nie zapisując.
func (a *adapterIzolacji) WczytajProfilIzolacji(ctx context.Context,
	z shared.IsolationProfileLoadRequest) (shared.IsolationProfileLoadResponse, error) {

	profil, err := a.profilZadania(ctx, z.ProfileId)
	if err != nil {
		return shared.IsolationProfileLoadResponse{}, err
	}
	return shared.IsolationProfileLoadResponse{Profile: a.profilKontraktu(profil)}, nil
}

// PrzypiszProfilIzolacji przepisuje przełączniki profilu pod wskazany adres
// i odnotowuje, z którego profilu polityka tego bytu pochodzi, odmawiając
// dla profilu bez przełączników.
func (a *adapterIzolacji) PrzypiszProfilIzolacji(ctx context.Context,
	z shared.IsolationProfileAssignRequest) (shared.IsolationProfileAssignResponse, error) {

	profil, err := a.profilZadania(ctx, z.ProfileId)
	if err != nil {
		return shared.IsolationProfileAssignResponse{}, err
	}
	adres, err := a.adresIzolacji(z.Scope, z.ScopeId, z.Layer)
	if err != nil {
		return shared.IsolationProfileAssignResponse{}, err
	}
	if len(profil.Przelaczniki) == 0 {
		return shared.IsolationProfileAssignResponse{}, bladZadaniaIzolacji(
			"profil " + profil.Kod + " nie niesie ani jednego przełącznika — przypisanie nie zmieniłoby niczego")
	}
	for _, przelacznik := range profil.Przelaczniki {
		punkt, znany := punktPoKluczu(przelacznik.Klucz)
		if !znany {
			return shared.IsolationProfileAssignResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodeInternalError, "izolacja: profil "+profil.Kod+" niesie punkt "+
					przelacznik.Klucz+" spoza wykazu jedenastu punktów izolacji"))
		}
		if err := a.zapiszPunkt(ctx, adres, punkt, odcietyPunkt(punkt, przelacznik.Wartosc)); err != nil {
			return shared.IsolationProfileAssignResponse{}, err
		}
	}
	if err := a.profile.PrzypiszProfil(ctx, profil.Kod, adres.Poziom, adres.KluczZasiegu); err != nil {
		return shared.IsolationProfileAssignResponse{}, err
	}
	warstwa, err := warstwaAdresu(adres.Poziom, z.Layer)
	if err != nil {
		return shared.IsolationProfileAssignResponse{}, err
	}
	polityka, err := a.politykaAdresu(ctx, adres, warstwa)
	if err != nil {
		return shared.IsolationProfileAssignResponse{}, err
	}
	// Przypisanie zmienia politykę bytu, nie sam profil, więc idzie rozgłoszenie polityki.
	a.rozglosPolitykeAdresu(ctx, adres)
	return shared.IsolationProfileAssignResponse{Policy: polityka}, nil
}

// UsunProfilIzolacji kasuje profil wraz z jego przełącznikami i przypisaniami,
// nie cofając wartości już przepisanych do poziomów przez wcześniejsze
// przypisanie.
func (a *adapterIzolacji) UsunProfilIzolacji(ctx context.Context,
	z shared.IsolationProfileDeleteRequest) (shared.IsolationProfileDeleteResponse, error) {

	if _, err := a.profilZadania(ctx, z.ProfileId); err != nil {
		return shared.IsolationProfileDeleteResponse{}, err
	}
	usuniety, err := a.profile.UsunProfil(ctx, z.ProfileId)
	if err != nil {
		return shared.IsolationProfileDeleteResponse{}, err
	}
	if !usuniety {
		return shared.IsolationProfileDeleteResponse{}, bladBrakuProfilu(z.ProfileId)
	}
	a.rozglosProfilIzolacji(ctx, shared.ChangeKindDeleted, z.ProfileId)
	return shared.IsolationProfileDeleteResponse{Deleted: true}, nil
}

// UstawWarstweIzolacji zapisuje warstwę, którą panel izolacji pokazuje dla
// wskazanej karty sesji albo okna komunikacji, nie zmieniając żadnej
// wartości izolacji.
func (a *adapterIzolacji) UstawWarstweIzolacji(ctx context.Context,
	z shared.IsolationLayerSetRequest) (shared.IsolationLayerSetResponse, error) {

	if a.profile == nil {
		return shared.IsolationLayerSetResponse{}, bladPodlozaIzolacji("magazyn warstwy izolacji")
	}
	warstwa, err := warstwaZadania(&z.Layer)
	if err != nil {
		return shared.IsolationLayerSetResponse{}, err
	}
	poziom, byt, err := bytWyboruWarstwy(z.SessionId, z.WindowId)
	if err != nil {
		return shared.IsolationLayerSetResponse{}, err
	}
	if err := a.profile.ZapiszWarstwe(ctx, poziom, byt, string(warstwa)); err != nil {
		return shared.IsolationLayerSetResponse{}, err
	}
	a.rozglosPolitykeWyboruWarstwy(ctx, poziom, byt)
	return shared.IsolationLayerSetResponse{Layer: warstwa}, nil
}

// bytWyboruWarstwy ustala byt, do którego należy wybór warstwy: kartę sesji, gdy
// wskazana, a w przeciwnym razie okno komunikacji.
func bytWyboruWarstwy(idSesji, idOkna *string) (shared.ConfigScope, string, error) {
	if karta := wartoscTekstu(idSesji); karta != "" {
		return shared.ConfigScopeSession, karta, nil
	}
	if okno := wartoscTekstu(idOkna); okno != "" {
		return shared.ConfigScopeWindow, okno, nil
	}
	return "", "", bladZadaniaIzolacji("przełączenie warstwy bez wskazania karty sesji ani okna komunikacji")
}

// PodgladPolitykiIzolacji zwraca politykę obowiązującą po rozstrzygnięciu
// ośmiu poziomów zasięgu i wskazanej warstwy, niczego nie zapisując.
func (a *adapterIzolacji) PodgladPolitykiIzolacji(ctx context.Context,
	z shared.IsolationPolicyPreviewRequest) (shared.IsolationPolicyPreviewResponse, error) {

	poziom, byt := poziomPodgladu(z)
	adres, err := a.adresIzolacji(poziom, &byt, z.Layer)
	if err != nil {
		return shared.IsolationPolicyPreviewResponse{}, err
	}
	warstwa, err := a.warstwaPodgladu(ctx, z, adres)
	if err != nil {
		return shared.IsolationPolicyPreviewResponse{}, err
	}
	polityka, err := a.politykaAdresu(ctx, adres, warstwa)
	if err != nil {
		return shared.IsolationPolicyPreviewResponse{}, err
	}
	return shared.IsolationPolicyPreviewResponse{Policy: polityka}, nil
}

// poziomPodgladu ustala poziom i byt podglądu z pól żądania, sięgając po
// okno komunikacji przed kartą sesji.
func poziomPodgladu(z shared.IsolationPolicyPreviewRequest) (shared.ConfigScope, string) {
	if z.Scope != nil && *z.Scope != "" {
		return *z.Scope, wartoscTekstu(z.ScopeId)
	}
	if okno := wartoscTekstu(z.WindowId); okno != "" {
		return shared.ConfigScopeWindow, okno
	}
	if karta := wartoscTekstu(z.SessionId); karta != "" {
		return shared.ConfigScopeSession, karta
	}
	return shared.ConfigScopeGlobal, ""
}

// warstwaPodgladu zwraca warstwę wskazaną w żądaniu, a bez wskazania — warstwę
// wybraną wcześniej dla tego bytu. Brak wyboru daje warstwę domyślną.
func (a *adapterIzolacji) warstwaPodgladu(ctx context.Context,
	z shared.IsolationPolicyPreviewRequest, adres konfig.Adres) (shared.IsolationLayer, error) {

	if z.Layer != nil && *z.Layer != "" {
		return warstwaAdresu(adres.Poziom, z.Layer)
	}
	if a.profile == nil {
		return shared.IsolationLayerDefault, nil
	}
	zapisana, jest, err := a.profile.Warstwa(ctx, adres.Poziom, adres.KluczZasiegu)
	if err != nil {
		return "", err
	}
	if !jest {
		return shared.IsolationLayerDefault, nil
	}
	wybrana := shared.IsolationLayer(zapisana)
	return warstwaZadania(&wybrana)
}

// politykaAdresu liczy politykę izolacji obowiązującą dla bytu spod adresu,
// odmawiając odpowiedzi, gdy zapisów źródła nie udało się odczytać.
func (a *adapterIzolacji) politykaAdresu(ctx context.Context, adres konfig.Adres,
	warstwa shared.IsolationLayer) (shared.IsolationPolicy, error) {

	polityka := a.rozstrzygacz.PolitykaEfektywna(kontekstPoziomu(adres.Poziom, adres.KluczZasiegu))
	if polityka.BladZrodla != nil {
		return shared.IsolationPolicy{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"izolacja: nie da się odczytać zapisów polityki: "+polityka.BladZrodla.Error()))
	}
	wynik := shared.IsolationPolicy{
		Scope:             adres.Poziom,
		Layer:             warstwa,
		ContextSwitches:   przelacznikiKontekstuPolityki(polityka),
		TechnicalSwitches: przelacznikiTechnicznePolityki(polityka),
		Origin:            pochodzeniePolityki(polityka),
	}
	if adres.KluczZasiegu != "" {
		byt := adres.KluczZasiegu
		wynik.ScopeId = &byt
	}
	if a.profile != nil {
		kod, jest, err := a.profile.ProfilPoziomu(ctx, adres.Poziom, adres.KluczZasiegu)
		if err != nil {
			return shared.IsolationPolicy{}, err
		}
		if jest {
			profil := kod
			wynik.ProfileId = &profil
		}
	}
	return wynik, nil
}

// kontekstPoziomu składa kontekst rozstrzygania z jednym bytem — tym, dla
// którego liczona jest polityka. Poziomy szersze zostają puste i rozstrzyganie
// schodzi na wartości globalne albo domyślne.
func kontekstPoziomu(poziom shared.ConfigScope, byt string) konfig.Kontekst {
	kontekst := konfig.Kontekst{}
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
	return kontekst
}

// przelacznikiKontekstuPolityki przekłada trzy wymiary kontekstu polityki na
// przełączniki kształtu kontraktu.
func przelacznikiKontekstuPolityki(polityka konfig.Polityka) []shared.IsolationSwitch {
	przelaczniki := make([]shared.IsolationSwitch, 0, len(punktyKontekstu))
	for _, pozycja := range punktyKontekstu {
		wartosc, objasnienie := wynikPunktu(polityka, pozycja.punkt)
		przelaczniki = append(przelaczniki, shared.IsolationSwitch{
			Kind:        pozycja.rodzaj,
			Isolated:    odcietyPunkt(pozycja.punkt, wartosc),
			Explanation: objasnienie,
		})
	}
	return przelaczniki
}

// przelacznikiTechnicznePolityki przekłada osiem zakresów technicznych
// polityki na przełączniki kształtu kontraktu.
func przelacznikiTechnicznePolityki(polityka konfig.Polityka) []shared.IsolationTechnicalSwitch {
	przelaczniki := make([]shared.IsolationTechnicalSwitch, 0, len(punktyTechniczne))
	for _, pozycja := range punktyTechniczne {
		wartosc, objasnienie := wynikPunktu(polityka, pozycja.punkt)
		przelaczniki = append(przelaczniki, shared.IsolationTechnicalSwitch{
			Scope:       pozycja.zakres,
			Isolated:    odcietyPunkt(pozycja.punkt, wartosc),
			Explanation: objasnienie,
		})
	}
	return przelaczniki
}

// wynikPunktu wyjmuje z polityki wartość jednego punktu izolacji wraz z jego
// objaśnieniem, gdy ono istnieje.
func wynikPunktu(polityka konfig.Polityka, punkt punktIzolacji) (string, *string) {
	pozycja, jest := polityka.Pozycja(punkt.klucz)
	if !jest {
		return "", nil
	}
	if pozycja.Objasnienie == "" {
		return pozycja.Wartosc, nil
	}
	objasnienie := pozycja.Objasnienie
	return pozycja.Wartosc, &objasnienie
}

// pochodzeniePolityki wskazuje poziom, z którego polityka została
// odziedziczona, pozostawiając pole puste, gdy punkty pochodzą z poziomów
// różnych.
func pochodzeniePolityki(polityka konfig.Polityka) *shared.ConfigScope {
	var znaleziony konfig.Poziom
	for _, pozycja := range punktyKontekstu {
		poziom, pochodzenie := polityka.Zrodlo(pozycja.punkt.klucz)
		if pochodzenie != konfig.PochodzenieZapis {
			continue
		}
		if znaleziony != konfig.PoziomBrak && znaleziony != poziom {
			return nil
		}
		znaleziony = poziom
	}
	for _, pozycja := range punktyTechniczne {
		poziom, pochodzenie := polityka.Zrodlo(pozycja.punkt.klucz)
		if pochodzenie != konfig.PochodzenieZapis {
			continue
		}
		if znaleziony != konfig.PoziomBrak && znaleziony != poziom {
			return nil
		}
		znaleziony = poziom
	}
	if znaleziony == konfig.PoziomBrak {
		return nil
	}
	poziom := shared.ConfigScope(znaleziony)
	return &poziom
}

// profilZadania odczytuje profil wskazany identyfikatorem. Brak wskazania jest
// żądaniem niezgodnym z kontraktem, brak wiersza — bytem, którego nie ma.
func (a *adapterIzolacji) profilZadania(ctx context.Context, kod string) (dane.ProfilIzolacji, error) {
	if a.profile == nil {
		return dane.ProfilIzolacji{}, bladPodlozaIzolacji("magazyn profili izolacji")
	}
	if strings.TrimSpace(kod) == "" {
		return dane.ProfilIzolacji{}, bladZadaniaIzolacji("komenda profilu bez wskazania profilu")
	}
	profil, err := a.profile.Profil(ctx, kod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.ProfilIzolacji{}, bladBrakuProfilu(kod)
	}
	if err != nil {
		return dane.ProfilIzolacji{}, err
	}
	return profil, nil
}

// przelacznikiProfilu przekłada przełączniki żądania na wiersze profilu. Wartość
// zapisywana jest dokładnie tą, którą czyta rozstrzygacz, więc przypisanie
// profilu jest przepisaniem wiersza, bez ani jednego przekładu po drodze.
func przelacznikiProfilu(kontekst []shared.IsolationSwitch,
	techniczne []shared.IsolationTechnicalSwitch) ([]dane.PrzelacznikProfiluIzolacji, error) {

	przelaczniki := make([]dane.PrzelacznikProfiluIzolacji, 0, len(kontekst)+len(techniczne))
	for _, przelacznik := range kontekst {
		punkt, znany := punktKontekstu(przelacznik.Kind)
		if !znany {
			return nil, bladNieznanegoPunktu(string(przelacznik.Kind))
		}
		przelaczniki = append(przelaczniki, dane.PrzelacznikProfiluIzolacji{
			Klucz: punkt.klucz, Wartosc: wartoscPunktu(punkt, przelacznik.Isolated),
		})
	}
	for _, przelacznik := range techniczne {
		punkt, znany := punktTechniczny(przelacznik.Scope)
		if !znany {
			return nil, bladNieznanegoPunktu(string(przelacznik.Scope))
		}
		przelaczniki = append(przelaczniki, dane.PrzelacznikProfiluIzolacji{
			Klucz: punkt.klucz, Wartosc: wartoscPunktu(punkt, przelacznik.Isolated),
		})
	}
	return przelaczniki, nil
}

// profilKontraktu przekłada wiersze profilu na kształt kontraktu, zostawiając
// punkt nieujęty w profilu pominiętym, a nie dopowiedzianym wartością
// domyślną.
func (a *adapterIzolacji) profilKontraktu(profil dane.ProfilIzolacji) shared.IsolationProfile {
	kontekst := make([]shared.IsolationSwitch, 0, len(punktyKontekstu))
	techniczne := make([]shared.IsolationTechnicalSwitch, 0, len(punktyTechniczne))
	zapisane := make(map[string]string, len(profil.Przelaczniki))
	for _, przelacznik := range profil.Przelaczniki {
		zapisane[przelacznik.Klucz] = przelacznik.Wartosc
	}
	for _, pozycja := range punktyKontekstu {
		wartosc, jest := zapisane[pozycja.punkt.klucz]
		if !jest {
			continue
		}
		kontekst = append(kontekst, shared.IsolationSwitch{
			Kind:        pozycja.rodzaj,
			Isolated:    odcietyPunkt(pozycja.punkt, wartosc),
			Explanation: a.objasnieniePunktu(pozycja.punkt),
		})
	}
	for _, pozycja := range punktyTechniczne {
		wartosc, jest := zapisane[pozycja.punkt.klucz]
		if !jest {
			continue
		}
		techniczne = append(techniczne, shared.IsolationTechnicalSwitch{
			Scope:       pozycja.zakres,
			Isolated:    odcietyPunkt(pozycja.punkt, wartosc),
			Explanation: a.objasnieniePunktu(pozycja.punkt),
		})
	}
	wynik := shared.IsolationProfile{
		Id:                profil.Kod,
		Name:              profil.Nazwa,
		ContextSwitches:   kontekst,
		TechnicalSwitches: techniczne,
	}
	if profil.Opis != "" {
		opis := profil.Opis
		wynik.Description = &opis
	}
	if chwila, czytelna := chwilaZapisu(profil.Utworzono); czytelna {
		wynik.CreatedAt = chwila
	}
	if chwila, czytelna := chwilaZapisu(profil.Zaktualizowano); czytelna {
		wynik.UpdatedAt = chwila
	}
	return wynik
}

// punktPoKluczu rozpoznaje punkt izolacji po kluczu ustawienia, przeszukując
// kontekstowe i techniczne punkty.
func punktPoKluczu(klucz string) (punktIzolacji, bool) {
	for _, pozycja := range punktyKontekstu {
		if pozycja.punkt.klucz == klucz {
			return pozycja.punkt, true
		}
	}
	for _, pozycja := range punktyTechniczne {
		if pozycja.punkt.klucz == klucz {
			return pozycja.punkt, true
		}
	}
	return punktIzolacji{}, false
}

// bladBrakuProfilu składa odmowę wskazującą profil, którego magazyn profili
// izolacji w ogóle nie posiada.
func bladBrakuProfilu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"izolacja: profil "+kod+" nie istnieje"))
}
