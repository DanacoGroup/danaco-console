// Adapter rodziny extension.* obsługuje katalog rozszerzeń platformy.
// Rozszerzenie jest pozycją katalogu, którą ekspert dopiero bierze, a nie
// mostem MCP ani konektorem eksperta. Instalacja zakłada pozycję jako
// zainstalowaną i niczego nie uruchamia.
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

// przedrostekRozszerzenia nadaje tożsamość wierszowi katalogu (`Extension.id`).
// Kodem pozycji (`Extension.code`) rdzeń nie zarządza — ten przynosi żądanie
// i jest „stały między wydaniami", więc nie wolno go nadpisywać własnym kluczem.
const przedrostekRozszerzenia = "rozsz-"

// adapterRozszerzen wypełnia port Rozszerzenia. Rejestr jest zależnością
// obowiązkową. Katalog punktów dostępu i biblioteka ekspertów obsługują
// pojedyncze pola żądań i mogą być puste — wtedy żądanie niosące to pole
// odmawia z powodem.
type adapterRozszerzen struct {
	rejestr dane.RepozytoriumRozszerzen
	punkty  dane.RepozytoriumPunktowDostepu
	agenci  dane.RepozytoriumAgentow
	// rozgloszenie rozgłasza zmianę katalogu; nil znaczy port bez nadajnika,
	// komendy pracują bez zmian.
	rozgloszenie *emiter
	// teraz oddaje czas w milisekundach epoki, wydzielone w pole, by pochodził
	// z jednego zegara.
	teraz func() int64
	// magazyn trzyma bajty paczek instalacji Personal i logi piaskownicy,
	// wspólnie z wytworami Apps.
	magazyn *magazynTresciBiblioteki
	// katalogDanych jest korzeniem, względem którego liczone są odwołania
	// magazynu.
	katalogDanych string
	// sejf wydaje poświadczenie po kluczu jawnym; rodzina nigdy nie ogląda
	// treści, sprawdza obecność.
	sejf sejfKluczaWydawcy
}

// ZMagazynemRozszerzen wpina katalog danych rdzenia: magazyn bajtów paczek
// i logów piaskownicy oraz sejf poświadczeń. Bez niego przesyłka paczki
// i piaskownica odmawiają z powodem, zamiast meldować wytwór bez bajtów.
func (a *adapterRozszerzen) ZMagazynemRozszerzen(katalogDanych string) *adapterRozszerzen {
	a.katalogDanych = katalogDanych
	a.magazyn = magazynWytworowApp(katalogDanych)
	a.sejf = dane.NowySejfPlikowy(katalogDanych)
	return a
}

// nowyAdapterRozszerzen wiąże nowy adapter z katalogiem rozszerzeń i ustawia
// jego źródło czasu operacji.
func nowyAdapterRozszerzen(rejestr dane.RepozytoriumRozszerzen) *adapterRozszerzen {
	return &adapterRozszerzen{
		rejestr: rejestr,
		teraz:   func() int64 { return time.Now().UTC().UnixMilli() },
	}
}

// ZKatalogiemDostepu oddaje adapterowi katalog punktów dostępu
// (`migracja_013_punkty_dostepu.sql`) i bibliotekę ekspertów
// (`migracja_037_agenci.sql`). Pierwszy sprawdza wskazanie mostu
// w `install` i `configure`, druga — wskazanie eksperta w `list`.
func (a *adapterRozszerzen) ZKatalogiemDostepu(punkty dane.RepozytoriumPunktowDostepu,
	agenci dane.RepozytoriumAgentow) *adapterRozszerzen {

	a.punkty = punkty
	a.agenci = agenci
	return a
}

// ── extension.list ───────────────────────────────────────────────────────────

// Wykaz obsługuje `extension.list`: oddaje katalog w kolejności wyświetlania.
// Pole agentId nie zawęża wykazu, tylko liczy przypisywalność pozycji do
// wskazanego eksperta; adapter sprawdza takiego eksperta i odmawia, gdy nie
// istnieje.
func (a *adapterRozszerzen) Wykaz(ctx context.Context,
	z shared.ExtensionListRequest) (shared.ExtensionListResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionListResponse{}, err
	}
	filtr := dane.FiltrRozszerzen{TylkoZainstalowane: wartoscLogiczna(z.InstalledOnly)}
	if z.Kind != nil {
		rodzaj, err := rodzajRozszerzenia(string(*z.Kind))
		if err != nil {
			return shared.ExtensionListResponse{}, err
		}
		filtr.Rodzaj = rodzaj
	}
	if err := a.sprawdzEksperta(ctx, wartoscTekstu(z.AgentId)); err != nil {
		return shared.ExtensionListResponse{}, err
	}
	wiersze, err := a.rejestr.Rozszerzenia(ctx, filtr)
	if err != nil {
		return shared.ExtensionListResponse{}, bladRozszerzenia(err)
	}
	katalog := make([]shared.Extension, 0, len(wiersze))
	for _, wiersz := range wiersze {
		katalog = append(katalog, rozszerzenieKontraktu(wiersz))
	}
	return shared.ExtensionListResponse{Extensions: katalog}, nil
}

// ── extension.install ────────────────────────────────────────────────────────

// Zainstaluj obsługuje `extension.install`, zakładając pozycję katalogu bez
// sięgania do sieci. Nazwą pozycji jest jej kod. Powtórna instalacja pozycji
// zainstalowanej jest odmową, instalacja odinstalowanej jest jej przywróceniem.
func (a *adapterRozszerzen) Zainstaluj(ctx context.Context,
	z shared.ExtensionInstallRequest) (shared.ExtensionInstallResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionInstallResponse{}, err
	}
	kod := strings.TrimSpace(z.Code)
	if kod == "" {
		return shared.ExtensionInstallResponse{}, bladWskazaniaRozszerzenia("żądanie bez kodu pozycji")
	}
	rodzaj, err := rodzajRozszerzenia(string(z.Kind))
	if err != nil {
		return shared.ExtensionInstallResponse{}, err
	}
	punkt, err := a.punktZadania(ctx, z.AccessPointId)
	if err != nil {
		return shared.ExtensionInstallResponse{}, err
	}
	zrodlo := strings.TrimSpace(wartoscTekstu(z.Source))
	pochodzenie, err := zrodloPochodzenia(z.Origin)
	if err != nil {
		return shared.ExtensionInstallResponse{}, err
	}

	zastane, err := a.rejestr.RozszerzeniePoKodzie(ctx, kod)
	switch {
	case err == nil:
		return a.przywroc(ctx, zastane, rodzaj, punkt, zrodlo, pochodzenie, z.Config)
	case errors.Is(err, dane.ErrBrakWiersza):
	default:
		return shared.ExtensionInstallResponse{}, bladRozszerzenia(err)
	}

	zalozone, err := a.rejestr.ZalozRozszerzenie(ctx, dane.Rozszerzenie{
		Identyfikator: nowyIdentyfikator(przedrostekRozszerzenia),
		Kod:           kod,
		Rodzaj:        rodzaj,
		Nazwa:         kod,
		Zainstalowane: true,
		// Zestaw wbudowany staje włączony jako część funkcjonalności bazowej,
		// rozszerzenie własne wyłączone.
		Wlaczone:          pochodzenie == shared.ExtensionOriginDanaco,
		PunktDostepuID:    punkt,
		ZrodloDeklarowane: zrodlo,
		ZrodloPochodzenia: pochodzenie,
		Konfiguracja:      string(z.Config),
		Zaktualizowano:    a.teraz(),
	})
	if err != nil {
		return shared.ExtensionInstallResponse{}, bladRozszerzenia(err)
	}
	pozycja := rozszerzenieKontraktu(zalozone)
	a.rozglosRozszerzenie(shared.ChangeKindCreated, pozycja)
	return shared.ExtensionInstallResponse{Extension: pozycja}, nil
}

// przywroc obsługuje instalację kodu, który już stoi w katalogu. Rodzaj jest
// częścią tożsamości pozycji, więc ten sam kod zgłoszony pod innym rodzajem
// jest odmową, nie przepisaniem rodzaju.
func (a *adapterRozszerzen) przywroc(ctx context.Context, zastane dane.Rozszerzenie,
	rodzaj string, punkt *int64, zrodlo, pochodzenie string,
	konfiguracja json.RawMessage) (shared.ExtensionInstallResponse, error) {

	if zastane.Rodzaj != rodzaj {
		return shared.ExtensionInstallResponse{}, bladStanuRozszerzenia(
			"pozycja " + zastane.Kod + " stoi w katalogu jako " + zastane.Rodzaj +
				", a żądanie instaluje ją jako " + rodzaj)
	}
	if zastane.Zainstalowane {
		return shared.ExtensionInstallResponse{}, bladStanuRozszerzenia(
			"pozycja " + zastane.Kod + " jest już zainstalowana")
	}
	// Przywrócenie staje w stanie instalacji pierwszej, by odinstalowanie
	// i ponowna nie obchodziły stanu.
	zmiana := dane.ZmianaRozszerzenia{
		Zainstalowane:     znacznikRozszerzenia(true),
		Wlaczone:          znacznikRozszerzenia(pochodzenie == shared.ExtensionOriginDanaco),
		PunktDostepuID:    punkt,
		ZrodloPochodzenia: &pochodzenie,
	}
	if zrodlo != "" {
		zmiana.ZrodloDeklarowane = &zrodlo
	}
	if len(konfiguracja) > 0 {
		tresc := string(konfiguracja)
		zmiana.Konfiguracja = &tresc
	}
	przywrocone, err := a.rejestr.ZmienRozszerzenie(ctx, zastane.Identyfikator, zmiana, a.teraz())
	if err != nil {
		return shared.ExtensionInstallResponse{}, bladRozszerzenia(err)
	}
	pozycja := rozszerzenieKontraktu(przywrocone)
	// Przywrócenie pozycji jest jej powrotem do katalogu, więc rozgłasza się
	// jak utworzenie.
	a.rozglosRozszerzenie(shared.ChangeKindCreated, pozycja)
	return shared.ExtensionInstallResponse{Extension: pozycja}, nil
}

// ── extension.configure ──────────────────────────────────────────────────────

// Skonfiguruj obsługuje `extension.configure`: zapisuje konfigurację pozycji
// i, gdy żądanie je niesie, wskazanie punktu dostępu. Konfiguracja jest
// wymagana i musi być poprawnym zapisem JSON, inaczej żądanie dostaje odmowę.
func (a *adapterRozszerzen) Skonfiguruj(ctx context.Context,
	z shared.ExtensionConfigureRequest) (shared.ExtensionConfigureResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionConfigureResponse{}, err
	}
	identyfikator := strings.TrimSpace(z.ExtensionId)
	if identyfikator == "" {
		return shared.ExtensionConfigureResponse{}, bladWskazaniaRozszerzenia("żądanie bez rozszerzenia")
	}
	if len(z.Config) == 0 {
		return shared.ExtensionConfigureResponse{}, bladWskazaniaRozszerzenia(
			"żądanie bez konfiguracji, a kontrakt wymaga jej wprost")
	}
	if !json.Valid(z.Config) {
		return shared.ExtensionConfigureResponse{}, bladWskazaniaRozszerzenia(
			"konfiguracja rozszerzenia " + identyfikator + " nie jest poprawnym JSON-em")
	}
	punkt, err := a.punktZadania(ctx, z.AccessPointId)
	if err != nil {
		return shared.ExtensionConfigureResponse{}, err
	}
	tresc := string(z.Config)
	zmienione, err := a.rejestr.ZmienRozszerzenie(ctx, identyfikator,
		dane.ZmianaRozszerzenia{Konfiguracja: &tresc, PunktDostepuID: punkt}, a.teraz())
	if err != nil {
		return shared.ExtensionConfigureResponse{}, bladRozszerzenia(err)
	}
	pozycja := rozszerzenieKontraktu(zmienione)
	a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
	return shared.ExtensionConfigureResponse{Extension: pozycja}, nil
}

// ── extension.toggle ─────────────────────────────────────────────────────────

// Przestaw obsługuje `extension.toggle`. Pozycji niezainstalowanej nie da się
// włączyć: włączenie odinstalowanej byłoby znacznikiem bez pokrycia, więc
// odmowa idzie kodem konfliktu stanu.
func (a *adapterRozszerzen) Przestaw(ctx context.Context,
	z shared.ExtensionToggleRequest) (shared.ExtensionToggleResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionToggleResponse{}, err
	}
	identyfikator := strings.TrimSpace(z.ExtensionId)
	if identyfikator == "" {
		return shared.ExtensionToggleResponse{}, bladWskazaniaRozszerzenia("żądanie bez rozszerzenia")
	}
	zastane, err := a.rejestr.Rozszerzenie(ctx, identyfikator)
	if err != nil {
		return shared.ExtensionToggleResponse{}, bladRozszerzenia(err)
	}
	if z.Enabled && !zastane.Zainstalowane {
		return shared.ExtensionToggleResponse{}, bladStanuRozszerzenia(
			"pozycja " + zastane.Kod + " nie jest zainstalowana, więc nie da się jej włączyć")
	}
	przestawione, err := a.rejestr.ZmienRozszerzenie(ctx, identyfikator,
		dane.ZmianaRozszerzenia{Wlaczone: znacznikRozszerzenia(z.Enabled)}, a.teraz())
	if err != nil {
		return shared.ExtensionToggleResponse{}, bladRozszerzenia(err)
	}
	pozycja := rozszerzenieKontraktu(przestawione)
	a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
	return shared.ExtensionToggleResponse{Extension: pozycja}, nil
}

// ── extension.uninstall ──────────────────────────────────────────────────────

// Odinstaluj obsługuje `extension.uninstall`: zdejmuje znaczniki instalacji
// i włączenia, zostawiając pozycję w katalogu. Pozycja, której nie ma, jest
// odmową; pozycja już odinstalowana zwraca wynik bez zmiany.
func (a *adapterRozszerzen) Odinstaluj(ctx context.Context,
	z shared.ExtensionUninstallRequest) (shared.ExtensionUninstallResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionUninstallResponse{}, err
	}
	identyfikator := strings.TrimSpace(z.ExtensionId)
	if identyfikator == "" {
		return shared.ExtensionUninstallResponse{}, bladWskazaniaRozszerzenia("żądanie bez rozszerzenia")
	}
	zastane, err := a.rejestr.Rozszerzenie(ctx, identyfikator)
	if err != nil {
		return shared.ExtensionUninstallResponse{}, bladRozszerzenia(err)
	}
	if !zastane.Zainstalowane {
		return shared.ExtensionUninstallResponse{Uninstalled: false}, nil
	}
	odinstalowane, err := a.rejestr.ZmienRozszerzenie(ctx, identyfikator, dane.ZmianaRozszerzenia{
		Zainstalowane: znacznikRozszerzenia(false),
		Wlaczone:      znacznikRozszerzenia(false),
	}, a.teraz())
	if err != nil {
		return shared.ExtensionUninstallResponse{}, bladRozszerzenia(err)
	}
	// Ładunek niesie wiersz po zdjęciu znaczników, nie ten sprzed, bo kasowania
	// wiersza tu nie ma.
	a.rozglosRozszerzenie(shared.ChangeKindDeleted, rozszerzenieKontraktu(odinstalowane))
	return shared.ExtensionUninstallResponse{Uninstalled: true}, nil
}

// ── wspólne ustalenia żądań ──────────────────────────────────────────────────

// punktZadania przekłada wskazany punkt dostępu na jego klucz w katalogu;
// brak wskazania zostawia kolumnę bez zmiany. Punkt jest sprawdzany
// w katalogu, nie przyjmowany na słowo, bo zapis wskazania bez mostu dałby
// pozycję wskazującą donikąd.
func (a *adapterRozszerzen) punktZadania(ctx context.Context, wskazanie *string) (*int64, error) {
	kod := strings.TrimSpace(wartoscTekstu(wskazanie))
	if kod == "" {
		return nil, nil
	}
	if a.punkty == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"katalog rozszerzeń: katalog punktów dostępu nie jest wpięty, wskazania mostu nie da się sprawdzić"))
	}
	punkt, err := a.punkty.PoKodzie(ctx, kod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"katalog rozszerzeń: punkt dostępu "+kod+" nie istnieje"))
	}
	if err != nil {
		return nil, bladRozszerzenia(err)
	}
	klucz := punkt.ID
	return &klucz, nil
}

// sprawdzEksperta sprawdza wskazanego eksperta z `extension.list`. Wskazanie
// puste jest zgodne z kontraktem i nie robi nic.
func (a *adapterRozszerzen) sprawdzEksperta(ctx context.Context, kod string) error {
	kod = strings.TrimSpace(kod)
	if kod == "" {
		return nil
	}
	if a.agenci == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"katalog rozszerzeń: biblioteka ekspertów nie jest wpięta, wskazania eksperta nie da się sprawdzić"))
	}
	if _, err := a.agenci.PoKodzie(ctx, kod); err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
				"katalog rozszerzeń: ekspert "+kod+" nie istnieje"))
		}
		return bladRozszerzenia(err)
	}
	return nil
}

// rozszerzenieKontraktu przekłada wiersz katalogu bazy na kształt kontraktu
// przesyłany do klienta okna.
func rozszerzenieKontraktu(r dane.Rozszerzenie) shared.Extension {
	rozszerzenie := shared.Extension{
		Id:            r.Identyfikator,
		Code:          r.Kod,
		Name:          r.Nazwa,
		Kind:          shared.ExtensionKind(r.Rodzaj),
		Description:   r.Opis,
		Version:       r.Wersja,
		Installed:     r.Zainstalowane,
		Enabled:       r.Wlaczone,
		Origin:        pochodzenieKontraktu(r.ZrodloPochodzenia),
		AccessPointId: r.PunktDostepuKod,
		UpdatedAt:     r.Zaktualizowano,
	}
	if r.Konfiguracja != "" {
		rozszerzenie.Config = json.RawMessage(r.Konfiguracja)
	}
	return rozszerzenie
}

// zrodlaRozszerzenia wylicza dwa jedynie dopuszczalne źródła pochodzenia
// pozycji katalogu tego kontraktu.
var zrodlaRozszerzenia = map[shared.ExtensionOrigin]struct{}{
	shared.ExtensionOriginDanaco:   {},
	shared.ExtensionOriginPersonal: {},
}

// zrodloPochodzenia czyta źródło z żądania. Pominięte znaczy `personal`, bo
// komendę woła Operator ze swojego urządzenia — pakiet serwera nie instaluje się
// przez WebSocket.
func zrodloPochodzenia(zrodlo *shared.ExtensionOrigin) (string, error) {
	if zrodlo == nil {
		return shared.ExtensionOriginPersonal, nil
	}
	if _, jest := zrodlaRozszerzenia[*zrodlo]; !jest {
		return "", bladWskazaniaRozszerzenia("źródło " + string(*zrodlo) +
			" nie należy do kontraktu")
	}
	return string(*zrodlo), nil
}

// pochodzenieKontraktu przekłada źródło z bazy. Wartość nierozpoznana czyta się jako
// `personal` — ogłoszenie „to część pakietu serwera" ma padać wyłącznie wtedy,
// gdy ktoś to zapisał.
func pochodzenieKontraktu(zrodlo string) shared.ExtensionOrigin {
	wartosc := shared.ExtensionOrigin(strings.TrimSpace(zrodlo))
	if _, jest := zrodlaRozszerzenia[wartosc]; !jest {
		return shared.ExtensionOriginPersonal
	}
	return wartosc
}

// rodzajRozszerzenia przepuszcza wyłącznie cztery wartości `ExtensionKind`.
// Kolumna `rozszerzenie.rodzaj` niesie wartość kontraktu wprost, więc przekładu
// tu nie ma — jest sprawdzenie.
func rodzajRozszerzenia(rodzaj string) (string, error) {
	switch rodzaj {
	case shared.ExtensionKindMcp, shared.ExtensionKindPlugin,
		shared.ExtensionKindApi, shared.ExtensionKindSkill:
		return rodzaj, nil
	}
	return "", bladWskazaniaRozszerzenia("rodzaj rozszerzenia " + rodzaj + " nie należy do kontraktu")
}

// znacznikRozszerzenia oddaje wskaźnik na wartość logiczną. Pola `dane.ZmianaRozszerzenia`
// są wskaźnikami, bo nil znaczy tam „bez zmiany", a nie „fałsz".
func znacznikRozszerzenia(wartosc bool) *bool {
	kopia := wartosc
	return &kopia
}

// sprawdzKatalog odmawia czynności, gdy katalog rozszerzeń nie został wpięty.
// Pięć komend ma wtedy odmówić głośno, a nie odpowiedzieć ciszą pustym
// wynikiem wyglądającym jak stan platformy.
func (a *adapterRozszerzen) sprawdzKatalog() error {
	if a == nil || a.rejestr == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"katalog rozszerzeń: rejestr rozszerzeń nie jest wpięty"))
	}
	return nil
}

// bladRozszerzenia przekłada niepowodzenie warstwy danych na odmowę kontraktu.
// Brak wiersza jest `not_found`, reszta — `internal_error`. Kodu spoza
// zamkniętego katalogu kontraktu tu nie ma i być nie może.
func bladRozszerzenia(err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"katalog rozszerzeń: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"katalog rozszerzeń: "+err.Error()))
}

// bladWskazaniaRozszerzenia odmawia żądaniu niezgodnemu z kontraktem tego
// katalogu rozszerzeń platformy.
func bladWskazaniaRozszerzenia(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"katalog rozszerzeń: "+powod))
}

// bladStanuRozszerzenia odmawia czynności, której stan tego katalogu
// rozszerzeń obecnie nie dopuszcza.
func bladStanuRozszerzenia(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"katalog rozszerzeń: "+powod))
}
