// Adapter rodziny `extension.*` — katalogu rozszerzeń platformy. Port i wpięcie
// stoją w `handlers_extension.go`, schemat katalogu —
// w `migracja_070_katalog_rozszerzen.sql`.
//
// ── Czym jest rozszerzenie w tej platformie ──────────────────────────────────
//
// Rozszerzenie jest pozycją katalogu platformy, którą ekspert dopiero bierze —
// ani mostem MCP, ani konektorem eksperta:
//
//   - Mostem MCP nie jest. Most mieszka w `punkt_dostepu`
//     (`migracja_013_punkty_dostepu.sql`) i z tej
//     tabeli `most_okna.go` składa `mcpServers` procesu modelu. `ExtensionKind`
//     niesie cztery wartości — `mcp`, `plugin`, `api`, `skill` — a w punkcie
//     dostępu da się zapisać wyłącznie maszynę mostu albo katalog lokalny.
//     Pozycja rodzaju `mcp` wskazuje most polem `accessPointId` i tyle; drugiego
//     rejestru serwerów MCP platforma nie ma.
//   - Konektorem eksperta nie jest. `agent_konektor`
//     (`migracja_038_agenci_konektory.sql`) ma
//     `agent_id NOT NULL`, więc konektor należy do jednego eksperta; katalog
//     stoi poziom wyżej i do eksperta nie należy. Pole `extension.list.agentId`
//     pyta o przypisywalność pozycji, a takie pytanie ma sens tylko wtedy, gdy
//     pozycja istnieje niezależnie od eksperta. Do tego `AgentConnectorKind` ma
//     trzy wartości, a `ExtensionKind` cztery — skilla konektorem zapisać się
//     nie da (umiejętności trzyma `agent_umiejetnosc`).
//
// Katalog jest więc warstwą nad tymi dwiema tabelami, a nie ich kopią: pozycję
// rodzaju `mcp`/`plugin`/`api` ekspert bierze zakładając sobie `agent_konektor`
// (`agent.connector.add`), pozycję rodzaju `skill` — zakładając
// `agent_umiejetnosc` (`agent.skill.add`). Ten adapter żadnej z tych dwóch tabel
// nie tyka: prawdą o zasobach eksperta pozostają jego własne komendy.
//
// ── Instalacja niczego nie pobiera i niczego nie uruchamia ───────────────────
//
// Instalacja zakłada pozycję w katalogu platformy i znaczy ją jako
// zainstalowaną. Nic więcej: nie sięga do sieci, nie pobiera paczki, nie
// rozpakowuje archiwum, nie uruchamia procesu i nie zakłada punktu dostępu.
// Szersze znaczenie zmieniłoby klasę bezpieczeństwa całego produktu —
// z programu, który wykonuje wyłącznie własny kod, na program, który wykonuje
// kod przyniesiony z zewnątrz.
//
// Napis podany w polu `source` nie jest wyrzucany — ląduje w kolumnie
// `zrodlo_deklarowane` jako zapis faktu, że Operator taki adres podał. Kształt
// `Extension` pola `source` nie ma, więc nie wychodzi w odpowiedzi.
//
// Odinstalowanie zdejmuje znaczniki, nie wiersz. Gdyby kasowało pozycję,
// `extension.list` nie miałby czego zawężać polem `installedOnly` — katalog
// zawierałby wtedy wyłącznie pozycje zainstalowane i przełącznik byłby zawsze
// prawdziwy. Odinstalowana pozycja zostaje w katalogu ze swoją konfiguracją,
// a `extension.install` tym samym kodem ją przywraca.
//
// Katalog jest rejestrem i nikt w rdzeniu go jeszcze nie czyta: włączenie
// pozycji rodzaju `mcp` nie dokłada mostu do procesu modelu, bo `mcpServers`
// składa się z `punkt_dostepu` i drugi pisarz tej listy byłby drugą prawdą.
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

// adapterRozszerzen wypełnia port Rozszerzenia.
//
// Rejestr jest zależnością obowiązkową. Katalog punktów dostępu i biblioteka
// ekspertów obsługują pojedyncze pola żądań (`accessPointId`, `agentId`) i mogą
// być puste — wtedy żądanie niosące to pole odmawia z powodem, zamiast przyjąć
// wskazanie, którego nikt nie sprawdził.
type adapterRozszerzen struct {
	rejestr dane.RepozytoriumRozszerzen
	punkty  dane.RepozytoriumPunktowDostepu
	agenci  dane.RepozytoriumAgentow
	// rozgloszenie rozgłasza `extension.changed`. Nil znaczy port złożony bez
	// nadajnika — komendy pracują wtedy bez zmian. Reguła rodzajów
	// zmiany i dokładka `ZRozgloszeniem` stoją w
	// `adapter_modul_extension_rozgloszenie.go`.
	rozgloszenie *emiter
	// teraz oddaje czas w milisekundach epoki. Wydzielone w pole, bo kolumna
	// czasu niesie wartość kontraktu wprost i musi pochodzić z jednego zegara.
	teraz func() int64
	// magazyn trzyma bajty paczek przesłanych instalacją Personal oraz logi
	// piaskownicy. Ten sam magazyn, którym jadą wytwory modułu Apps — pozycja
	// katalogu i pakiet, z którego powstała, leżą obok siebie.
	magazyn *magazynTresciBiblioteki
	// katalogDanych jest korzeniem, względem którego liczone są odwołania
	// magazynu.
	katalogDanych string
	// sejf wydaje poświadczenie po jego kluczu jawnym. Rodzina `extension.*`
	// nigdy nie ogląda treści poświadczenia — sprawdza wyłącznie, czy odwołanie
	// w sejfie w ogóle stoi.
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

// nowyAdapterRozszerzen wiąże adapter z katalogiem rozszerzeń.
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
//
// Pole `agentId` nie zawęża wykazu. Opisuje eksperta, dla którego liczona jest
// przypisywalność pozycji, ale kształt `Extension` nie ma ani jednego pola,
// w którym wynik takiego liczenia mógłby wyjechać. Obie oczywiste drogi są złe:
//
//   - użycie pola jako sita odebrałoby Connectors Managerowi widok pozycji już
//     wziętych przez eksperta, czyli dokładnie tych, którymi zarządza;
//   - zignorowanie pola w ciszy dałoby odpowiedź udającą, że policzono coś,
//     czego nie policzono.
//
// Adapter sprawdza więc wskazanego eksperta i odmawia `not_found`, gdy takiego
// nie ma — żądanie o eksperta nieistniejącego nie dostaje wtedy odpowiedzi
// wyglądającej na poprawną.
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

// Zainstaluj obsługuje `extension.install`. Znaczenie instalacji — najwęższe
// z trzech wymienionych w kontrakcie, bez sięgania do sieci — opisuje nagłówek
// pliku.
//
// Nazwą pozycji jest jej kod, bo żądanie nazwy nie niesie. Kształt `Extension`
// wymaga pola `name`, a `ExtensionInstallRequest` ma `code`, `kind`, `source`,
// `accessPointId` i `config` — i ani jednego pola nazwy, opisu czy wersji.
// Wymyślanie ładniejszej nazwy z kodu byłoby zmyślaniem treści, której nikt nie
// podał, więc nazwa równa się kodowi, a opis i wersja zostają puste.
//
// Powtórna instalacja pozycji zainstalowanej jest odmową `conflict`, nie cichym
// nadpisaniem: instalujący drugi raz ma się dowiedzieć, że to już stoi.
// Instalacja pozycji odinstalowanej jest jej przywróceniem — to jedyna droga
// powrotu, bo odinstalowanie wiersza nie kasuje.
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
		// Stan wyjściowy zależy od źródła i jest to jedyne miejsce, w którym
		// rdzeń rozróżnia `danaco` od `personal`. Zestaw wbudowany jest częścią
		// funkcjonalności bazowej, więc staje włączony; rozszerzenie własne staje
		// wyłączone, bo włącza je świadoma decyzja po przejrzeniu konfiguracji
		// i zakresu. Poza tym jednym rozgałęzieniem ładowanie, kontrakt
		// integracji i użycie operacyjne są dla obu źródeł te same.
		//
		// Nie jest to bramka kontrolna: `extension.toggle` włącza pozycję bez
		// pytania kogokolwiek o zgodę.
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

// przywroc obsługuje instalację kodu, który już stoi w katalogu.
//
// Rodzaj jest częścią tożsamości pozycji. Kod jest „stały między wydaniami",
// więc ten sam kod zgłoszony raz jako `mcp`, a raz jako `plugin` opisuje dwie
// różne rzeczy pod jedną nazwą. Odmowa `conflict` mówi to wprost, zamiast
// przepisywać rodzaj pozycji zza pleców katalogu.
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
	// Przywrócenie staje w tym samym stanie co instalacja pierwsza. Pozycja
	// `personal` wraca wyłączona, bo powrót do katalogu jest tą samą rejestracją
	// co za pierwszym razem — inaczej odinstalowanie i zainstalowanie z powrotem
	// byłoby drogą obejścia stanu wyjściowego.
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
	// Przywrócenie pozycji jest dla Operatora jej powrotem do katalogu
	// zainstalowanych, więc `created` — uzasadnienie symetrii z `deleted`
	// w `adapter_modul_extension_rozgloszenie.go`.
	a.rozglosRozszerzenie(shared.ChangeKindCreated, pozycja)
	return shared.ExtensionInstallResponse{Extension: pozycja}, nil
}

// ── extension.configure ──────────────────────────────────────────────────────

// Skonfiguruj obsługuje `extension.configure`: zapisuje konfigurację pozycji
// i — gdy żądanie je niesie — wskazanie punktu dostępu.
//
// Konfiguracja jest wymagana i musi być JSON-em. Kontrakt oznacza pole `config`
// jako wymagane typu `json`; żądanie bez niego albo z treścią, która JSON-em nie
// jest, dostaje `validation_failed`, a nie zapis pustego obiektu udający zapis
// konfiguracji.
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

// Przestaw obsługuje `extension.toggle`.
//
// Pozycji niezainstalowanej nie da się włączyć. Kontrakt opisuje tę komendę jako
// włączenie albo wyłączenie rozszerzenia „bez odinstalowania go" — czyli mówi
// o pozycji, która stoi zainstalowana. Włączenie pozycji odinstalowanej byłoby
// znacznikiem bez pokrycia; odmowa idzie kodem `conflict`: żądanie jest
// poprawne, lecz stan katalogu wyklucza czynność.
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
// i włączenia, zostawiając pozycję w katalogu — uzasadnienie w nagłówku pliku.
//
// Pozycja, której nie ma, jest odmową `not_found` z jej wskazaniem, a nie
// odpowiedzią `uninstalled: false` udającą wykonaną czynność. `false` ma tu
// jeden uczciwy przypadek — pozycję, która już była odinstalowana; wtedy nic nie
// doszło do skutku i wynik mówi to wprost.
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
	// Ładunek niesie wiersz po zdjęciu znaczników, nie ten sprzed. Kasowania
	// wiersza tu nie ma, więc odczyt uprzedni (wzorzec `memory.changed`) dałby
	// pozycję nadal zainstalowaną przy `change: deleted`.
	a.rozglosRozszerzenie(shared.ChangeKindDeleted, rozszerzenieKontraktu(odinstalowane))
	return shared.ExtensionUninstallResponse{Uninstalled: true}, nil
}

// ── wspólne ustalenia żądań ──────────────────────────────────────────────────

// punktZadania przekłada wskazany punkt dostępu na jego klucz w katalogu.
// Wynik nil bez błędu znaczy „żądanie punktu nie wskazało" — wtedy kolumna
// zostaje bez zmiany (patrz `dane.ZmianaRozszerzenia`).
//
// Punkt sprawdzany jest w katalogu, nie przyjmowany na słowo. Zapis wskazania,
// za którym nie stoi żaden most, dałby pozycję katalogu wskazującą na nic —
// i Operator dowiedziałby się o tym dopiero wtedy, gdy rozszerzenie miałoby
// zadziałać.
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
// puste jest zgodne z kontraktem i nie robi nic — patrz komentarz przy Wykaz.
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

// rozszerzenieKontraktu przekłada wiersz katalogu na kształt kontraktu.
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

// zrodlaRozszerzenia wylicza dwa źródła pochodzenia z kontraktu.
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
//
// Pięć komend ma wtedy odmówić głośno, a nie odpowiedzieć ciszą. Bez rejestru
// `extension.list` oddałby pusty katalog, a `extension.uninstall` —
// `uninstalled: false`; jedno i drugie wyglądałoby dla Operatora jak stan
// platformy, a nie jak nieskładany port. Kod `internal_error`:
// żądanie jest poprawne, wina leży po stronie montażu.
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

// bladWskazaniaRozszerzenia odmawia żądaniu niezgodnemu z kontraktem.
func bladWskazaniaRozszerzenia(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"katalog rozszerzeń: "+powod))
}

// bladStanuRozszerzenia odmawia czynności, której stan katalogu nie dopuszcza.
func bladStanuRozszerzenia(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"katalog rozszerzeń: "+powod))
}
