// Plik utrzymuje trwałość katalogu rozszerzeń w tabeli rozszerzenie wraz
// z rodziną extension.*, oddzielony od poziomu eksperta, mostu MCP i konektora
// eksperta, bez pobierania ani uruchamiania zadeklarowanego źródła.
package dane

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// Rozszerzenie to wiersz tabeli rozszerzenie z dwoma identyfikatorami:
// Identyfikator odpowiada polu Extension.id, a Kod polu Extension.code,
// kodowi trwałemu między wydaniami.
type Rozszerzenie struct {
	ID            int64
	Identyfikator string
	Kod           string
	Rodzaj        string
	Nazwa         string
	Opis          *string
	Wersja        *string
	Zainstalowane bool
	Wlaczone      bool
	// PunktDostepuID i PunktDostepuKod wskazują punkt dostępu: klucz wewnętrzny oraz jego kod trwały.
	PunktDostepuID  *int64
	PunktDostepuKod *string
	// ZrodloDeklarowane to napis podany w install.source. Nikt go nie pobiera.
	ZrodloDeklarowane string
	// ZrodloPochodzenia niesie shared.ExtensionOrigin: danaco albo personal, nie ZrodloDeklarowane.
	ZrodloPochodzenia string
	Konfiguracja      string
	Zaktualizowano    int64
}

// FiltrRozszerzen zawęża wykaz rozszerzeń i obsługuje oba pola zawężające
// polecenia extension.list: rodzaj pozycji oraz ograniczenie do zainstalowanych.
type FiltrRozszerzen struct {
	// Rodzaj pusty znaczy „wszystkie rodzaje".
	Rodzaj string
	// TylkoZainstalowane odpowiada polu installedOnly; domyślnie pokazuje też pozycje odinstalowane.
	TylkoZainstalowane bool
}

// ZmianaRozszerzenia niesie pola zapisu polecenia extension.configure; wskaźnik
// pusty w polu znaczy brak zmiany tego pola, a nie wyczyszczenie wartości.
type ZmianaRozszerzenia struct {
	// Nazwa, Opis i Wersja to metryki zmieniane publikacją pakietu oraz przypięciem wersji w managerze.
	Nazwa             *string
	Opis              *string
	Wersja            *string
	Zainstalowane     *bool
	Wlaczone          *bool
	PunktDostepuID    *int64
	ZrodloDeklarowane *string
	ZrodloPochodzenia *string
	Konfiguracja      *string
}

// RepozytoriumRozszerzen jest kontraktem katalogu rozszerzeń obejmującym
// cykl życia, warstwę protokołu, integracje zewnętrzne oraz warstwę zaufania pozycji.
type RepozytoriumRozszerzen interface {
	Rozszerzenia(ctx context.Context, filtr FiltrRozszerzen) ([]Rozszerzenie, error)
	// Rozszerzenie odczytuje pozycję po Extension.id.
	Rozszerzenie(ctx context.Context, identyfikator string) (Rozszerzenie, error)
	// RozszerzeniePoKodzie odczytuje pozycję po Extension.code; brak wraca jako ErrBrakWiersza.
	RozszerzeniePoKodzie(ctx context.Context, kod string) (Rozszerzenie, error)
	// ZalozRozszerzenie wstawia pozycję katalogu; czas zmiany podaje warstwa wyższa w milisekundach epoki.
	ZalozRozszerzenie(ctx context.Context, rozszerzenie Rozszerzenie) (Rozszerzenie, error)
	ZmienRozszerzenie(ctx context.Context, identyfikator string,
		zmiana ZmianaRozszerzenia, teraz int64) (Rozszerzenie, error)

	// --- cykl życia pozycji (extension_cykl.go) ---
	ZapiszKolekcjeRozszerzen(ctx context.Context, kolekcja KolekcjaRozszerzen) (KolekcjaRozszerzen, error)
	KolekcjaRozszerzen(ctx context.Context, kod string) (KolekcjaRozszerzen, error)
	KolekcjeRozszerzen(ctx context.Context) ([]KolekcjaRozszerzen, error)
	DopiszHistorieRozszerzenia(ctx context.Context, wpis WpisHistoriiRozszerzenia) error
	HistoriaRozszerzenia(ctx context.Context, rozszerzenie string, od int64, granica int) ([]WpisHistoriiRozszerzenia, int, error)
	ZapiszWersjeRozszerzenia(ctx context.Context, wersja WersjaRozszerzenia) error
	WersjeRozszerzenia(ctx context.Context, rozszerzenie string) ([]WersjaRozszerzenia, error)
	PrzypnijWersjeRozszerzenia(ctx context.Context, rozszerzenie, wersja string, teraz int64) error
	PrzypiecieWersjiRozszerzenia(ctx context.Context, rozszerzenie string) (string, error)
	ZalozPaczkeRozszerzenia(ctx context.Context, paczka PaczkaRozszerzenia) error
	PaczkaRozszerzenia(ctx context.Context, kod string) (PaczkaRozszerzenia, error)

	// --- warstwa protokołu (extension_protokol.go) ---
	ZapiszNarzedziaRozszerzenia(ctx context.Context, rozszerzenie string, wpisy []NarzedzieRozszerzenia) error
	NarzedziaRozszerzenia(ctx context.Context, rozszerzenie, rodzaj string) ([]NarzedzieRozszerzenia, error)
	DopiszRamkeProtokolu(ctx context.Context, ramka RamkaProtokolu) error
	RamkiProtokolu(ctx context.Context, rozszerzenie string, od int64, granica int) ([]RamkaProtokolu, int, error)
	DopiszWywolanieRozszerzenia(ctx context.Context, wywolanie WywolanieRozszerzenia) error
	MetrykiUzyciaRozszerzen(ctx context.Context, rozszerzenie string, od, do int64) ([]MetrykaUzyciaRozszerzenia, error)
	AudytRozszerzen(ctx context.Context, rozszerzenie, agent string, od int64, granica int) ([]WywolanieRozszerzenia, int, error)
	DopiszKondycjeRozszerzenia(ctx context.Context, kondycja KondycjaRozszerzenia) error
	OstatniaKondycjaRozszerzenia(ctx context.Context, rozszerzenie string) (KondycjaRozszerzenia, error)

	// --- integracje zewnętrzne (extension_integracje.go) ---
	ZapiszIntegracjeRozszerzenia(ctx context.Context, integracja IntegracjaRozszerzenia) (IntegracjaRozszerzenia, error)
	IntegracjaRozszerzenia(ctx context.Context, rozszerzenie string) (IntegracjaRozszerzenia, error)
	ZapiszWebhookRozszerzenia(ctx context.Context, webhook WebhookRozszerzenia) (WebhookRozszerzenia, error)
	WebhookRozszerzenia(ctx context.Context, kod string) (WebhookRozszerzenia, error)
	WebhookiRozszerzen(ctx context.Context, rozszerzenie, kierunek string) ([]WebhookRozszerzenia, error)
	ZapiszMapowanieRozszerzenia(ctx context.Context, mapowanie MapowanieRozszerzenia) (MapowanieRozszerzenia, error)
	MapowanieRozszerzenia(ctx context.Context, kod string) (MapowanieRozszerzenia, error)

	// --- warstwa zaufania (extension_zaufanie.go) ---
	ZapiszUprawnieniaRozszerzenia(ctx context.Context, rozszerzenie string, nadane bool, uprawnienia []UprawnienieRozszerzenia) error
	UprawnieniaRozszerzenia(ctx context.Context, rozszerzenie string) ([]UprawnienieRozszerzenia, error)
	ZapiszPodpisRozszerzenia(ctx context.Context, podpis PodpisRozszerzenia) error
	PodpisRozszerzenia(ctx context.Context, rozszerzenie string) (PodpisRozszerzenia, error)
	ZapiszSekretRozszerzenia(ctx context.Context, sekret SekretRozszerzenia, wymienZakres bool) (SekretRozszerzenia, error)
	SekretRozszerzenia(ctx context.Context, odwolanie string) (SekretRozszerzenia, error)
	SekretyRozszerzen(ctx context.Context, doCzasu int64) ([]SekretRozszerzenia, error)
}

const (
	kolumnyRozszerzenia = `r.id, r.identyfikator_zewnetrzny, r.kod, r.rodzaj, r.nazwa, r.opis,
	                       r.wersja, r.zainstalowane, r.wlaczone, r.punkt_dostepu_id, p.kod,
	                       r.zrodlo_deklarowane, r.zrodlo_pochodzenia, r.konfiguracja,
	                       r.zaktualizowano`

	zrodloRozszerzenia = ` FROM rozszerzenie r LEFT JOIN punkt_dostepu p ON p.id = r.punkt_dostepu_id`

	pobierzRozszerzenie = `SELECT ` + kolumnyRozszerzenia + zrodloRozszerzenia +
		` WHERE r.identyfikator_zewnetrzny = ?`

	rozszerzeniePoKodzie = `SELECT ` + kolumnyRozszerzenia + zrodloRozszerzenia + ` WHERE r.kod = ?`

	// Jedno zapytanie obsługuje cztery warianty żądania: puste zawężenie rodzaju
	// wyłącza pierwszy warunek, a zamknięty parametr installedOnly wyłącza drugi.
	// Porządek wynika z indeksu wykazu.
	listaRozszerzen = `SELECT ` + kolumnyRozszerzenia + zrodloRozszerzenia +
		` WHERE (? = '' OR r.rodzaj = ?) AND (? = 0 OR r.zainstalowane = 1)
		  ORDER BY r.rodzaj, r.nazwa, r.id`

	wstawRozszerzenie = `INSERT INTO rozszerzenie
	                     (identyfikator_zewnetrzny, kod, rodzaj, nazwa, opis, wersja,
	                      zainstalowane, wlaczone, punkt_dostepu_id, zrodlo_deklarowane,
	                      zrodlo_pochodzenia, konfiguracja, zaktualizowano)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Zmiana idzie jednym poleceniem: NULL w argumencie zostawia kolumnę bez
	// zmiany. Dzięki temu „pola pominięte zostają bez zmian" jest własnością
	// zapytania, a nie kolejnością gałęzi w Go.
	zmienRozszerzenie = `UPDATE rozszerzenie SET
	                        nazwa              = COALESCE(?, nazwa),
	                        opis               = COALESCE(?, opis),
	                        wersja             = COALESCE(?, wersja),
	                        zainstalowane      = COALESCE(?, zainstalowane),
	                        wlaczone           = COALESCE(?, wlaczone),
	                        punkt_dostepu_id   = COALESCE(?, punkt_dostepu_id),
	                        zrodlo_deklarowane = COALESCE(?, zrodlo_deklarowane),
	                        zrodlo_pochodzenia = COALESCE(?, zrodlo_pochodzenia),
	                        konfiguracja       = COALESCE(?, konfiguracja),
	                        zaktualizowano     = ?
	                     WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumRozszerzen struct {
	zapytania *zapytania
	// baza jest połączeniem potrzebnym transakcjom wielotabelowym repozytorium.
	baza *sql.DB
}

// Zgodność implementacji z kontraktem sprawdzana jest przy kompilacji, a nie
// dopiero przy złożeniu zestawu repozytoriów.
var _ RepozytoriumRozszerzen = (*repozytoriumRozszerzen)(nil)

// Rozszerzenia oddaje repozytorium katalogu rozszerzeń złożone nad
// współdzieloną pamięcią zapytań zestawu, bez własnego stanu poza wskaźnikiem na tę pamięć.
func (z *Zestaw) Rozszerzenia() RepozytoriumRozszerzen {
	if z == nil || z.zapytania == nil {
		return nil
	}
	return &repozytoriumRozszerzen{zapytania: z.zapytania, baza: z.baza}
}

// Rozszerzenia zwraca wykaz pozycji katalogu w stałej kolejności wyświetlania,
// zawężony filtrem rodzaju i stanu instalacji.
func (r *repozytoriumRozszerzen) Rozszerzenia(ctx context.Context,
	filtr FiltrRozszerzen) ([]Rozszerzenie, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaRozszerzen)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, filtr.Rodzaj, filtr.Rodzaj,
		liczbaLogiczna(filtr.TylkoZainstalowane))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać katalogu rozszerzeń: %w", err)
	}
	defer wiersze.Close()

	katalog := make([]Rozszerzenie, 0, 16)
	for wiersze.Next() {
		pozycja, err := odczytajRozszerzenie(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: uszkodzony wiersz rozszerzenia: %w", err)
		}
		katalog = append(katalog, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać katalogu rozszerzeń: %w", err)
	}
	return katalog, nil
}

// Rozszerzenie zwraca pozycję o wskazanym `Extension.id`. Brak wiersza wraca
// jako ErrBrakWiersza — warstwa wyższa odróżnia „nie ma" od „odczyt padł".
func (r *repozytoriumRozszerzen) Rozszerzenie(ctx context.Context,
	identyfikator string) (Rozszerzenie, error) {

	return r.jednaPozycja(ctx, pobierzRozszerzenie, identyfikator, "rozszerzenie")
}

// RozszerzeniePoKodzie zwraca pozycję katalogu o wskazanym Extension.code,
// kodzie trwałym między wydaniami pakietu.
func (r *repozytoriumRozszerzen) RozszerzeniePoKodzie(ctx context.Context,
	kod string) (Rozszerzenie, error) {

	return r.jednaPozycja(ctx, rozszerzeniePoKodzie, kod, "rozszerzenie o kodzie")
}

// jednaPozycja odczytuje jeden wiersz katalogu wskazanym zapytaniem i ujednolica
// obsługę braku wiersza dla obu punktów odczytu.
func (r *repozytoriumRozszerzen) jednaPozycja(ctx context.Context,
	zapytanie, wskazanie, nazwaBytu string) (Rozszerzenie, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return Rozszerzenie{}, err
	}
	pozycja, err := odczytajRozszerzenie(polecenie.QueryRowContext(ctx, wskazanie))
	if errors.Is(err, sql.ErrNoRows) {
		return Rozszerzenie{}, fmt.Errorf("dane: %s %q nie istnieje: %w",
			nazwaBytu, wskazanie, ErrBrakWiersza)
	}
	if err != nil {
		return Rozszerzenie{}, fmt.Errorf("dane: nie można odczytać %s %q: %w",
			nazwaBytu, wskazanie, err)
	}
	return pozycja, nil
}

// ZalozRozszerzenie wstawia nową pozycję katalogu, sprawdza jej wymagane pola
// i oddaje zapisany wiersz po odczycie.
func (r *repozytoriumRozszerzen) ZalozRozszerzenie(ctx context.Context,
	rozszerzenie Rozszerzenie) (Rozszerzenie, error) {

	if rozszerzenie.Identyfikator == "" {
		return Rozszerzenie{}, fmt.Errorf("dane: rozszerzenie bez identyfikatora")
	}
	if rozszerzenie.Kod == "" {
		return Rozszerzenie{}, fmt.Errorf("dane: rozszerzenie %q bez kodu pozycji",
			rozszerzenie.Identyfikator)
	}
	if rozszerzenie.Nazwa == "" {
		return Rozszerzenie{}, fmt.Errorf("dane: rozszerzenie %q bez nazwy", rozszerzenie.Kod)
	}
	konfiguracja, err := konfiguracjaRozszerzenia(rozszerzenie.Konfiguracja, rozszerzenie.Kod)
	if err != nil {
		return Rozszerzenie{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawRozszerzenie)
	if err != nil {
		return Rozszerzenie{}, err
	}
	if _, err := polecenie.ExecContext(ctx, rozszerzenie.Identyfikator, rozszerzenie.Kod,
		rozszerzenie.Rodzaj, rozszerzenie.Nazwa, tekstDoKolumny(rozszerzenie.Opis),
		tekstDoKolumny(rozszerzenie.Wersja), liczbaLogiczna(rozszerzenie.Zainstalowane),
		liczbaLogiczna(rozszerzenie.Wlaczone), liczbaDoKolumny(rozszerzenie.PunktDostepuID),
		rozszerzenie.ZrodloDeklarowane, zrodloPochodzeniaKolumny(rozszerzenie.ZrodloPochodzenia),
		konfiguracja, rozszerzenie.Zaktualizowano); err != nil {

		return Rozszerzenie{}, fmt.Errorf("dane: nie można założyć rozszerzenia %q: %w",
			rozszerzenie.Kod, err)
	}
	return r.Rozszerzenie(ctx, rozszerzenie.Identyfikator)
}

// ZmienRozszerzenie zmienia wyłącznie pola wskazane w żądaniu, pozostałe
// zostawiając bez zmiany, i oddaje wiersz po zapisie.
func (r *repozytoriumRozszerzen) ZmienRozszerzenie(ctx context.Context, identyfikator string,
	zmiana ZmianaRozszerzenia, teraz int64) (Rozszerzenie, error) {

	// Odczyt przed zapisem: zmiana pozycji nieistniejącej wraca jako ErrBrakWiersza, nie sukces.
	if _, err := r.Rozszerzenie(ctx, identyfikator); err != nil {
		return Rozszerzenie{}, err
	}
	konfiguracja, err := zmienionaKonfiguracjaRozszerzenia(zmiana.Konfiguracja, identyfikator)
	if err != nil {
		return Rozszerzenie{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienRozszerzenie)
	if err != nil {
		return Rozszerzenie{}, err
	}
	var zainstalowane, wlaczone any
	if zmiana.Zainstalowane != nil {
		zainstalowane = liczbaLogiczna(*zmiana.Zainstalowane)
	}
	if zmiana.Wlaczone != nil {
		wlaczone = liczbaLogiczna(*zmiana.Wlaczone)
	}
	if _, err := polecenie.ExecContext(ctx,
		tekstDoKolumny(zmiana.Nazwa), tekstDoKolumny(zmiana.Opis), tekstDoKolumny(zmiana.Wersja),
		zainstalowane, wlaczone,
		liczbaDoKolumny(zmiana.PunktDostepuID), tekstDoKolumny(zmiana.ZrodloDeklarowane),
		tekstDoKolumny(zmiana.ZrodloPochodzenia), konfiguracja, teraz, identyfikator); err != nil {

		return Rozszerzenie{}, fmt.Errorf("dane: nie można zmienić rozszerzenia %q: %w",
			identyfikator, err)
	}
	return r.Rozszerzenie(ctx, identyfikator)
}

// konfiguracjaRozszerzenia sprawdza, że konfiguracja jest poprawnym JSON-em,
// i zamienia brak na pusty obiekt — kolumna jest NOT NULL DEFAULT '{}'.
func konfiguracjaRozszerzenia(tresc, wskazanie string) (string, error) {
	if tresc == "" {
		return "{}", nil
	}
	if !json.Valid([]byte(tresc)) {
		return "", fmt.Errorf("dane: konfiguracja rozszerzenia %q nie jest poprawnym JSON-em", wskazanie)
	}
	return tresc, nil
}

// zmienionaKonfiguracjaRozszerzenia przekłada wskaźnik zmiany na argument
// zapytania: nil znaczy „bez zmiany" i zostawia kolumnę nietkniętą przez COALESCE.
func zmienionaKonfiguracjaRozszerzenia(tresc *string, wskazanie string) (any, error) {
	if tresc == nil {
		return nil, nil
	}
	sprawdzona, err := konfiguracjaRozszerzenia(*tresc, wskazanie)
	if err != nil {
		return nil, err
	}
	return sprawdzona, nil
}

// zrodloPochodzeniaKolumny pilnuje, żeby kolumna zrodlo_pochodzenia nigdy
// nie dostała pustki: nieustawiona wartość czyta się jako personal.
func zrodloPochodzeniaKolumny(zrodlo string) string {
	if strings.TrimSpace(zrodlo) == "" {
		return shared.ExtensionOriginPersonal
	}
	return zrodlo
}

func odczytajRozszerzenie(wiersz skaner) (Rozszerzenie, error) {
	var pozycja Rozszerzenie
	var opis, wersja, punktKod sql.NullString
	var punktID sql.NullInt64
	var zainstalowane, wlaczone int
	if err := wiersz.Scan(&pozycja.ID, &pozycja.Identyfikator, &pozycja.Kod, &pozycja.Rodzaj,
		&pozycja.Nazwa, &opis, &wersja, &zainstalowane, &wlaczone, &punktID, &punktKod,
		&pozycja.ZrodloDeklarowane, &pozycja.ZrodloPochodzenia, &pozycja.Konfiguracja,
		&pozycja.Zaktualizowano); err != nil {

		return Rozszerzenie{}, err
	}
	pozycja.Opis = tekstZKolumny(opis)
	pozycja.Wersja = tekstZKolumny(wersja)
	pozycja.PunktDostepuID = liczbaZKolumny(punktID)
	pozycja.PunktDostepuKod = tekstZKolumny(punktKod)
	pozycja.Zainstalowane = zainstalowane == 1
	pozycja.Wlaczone = wlaczone == 1
	return pozycja, nil
}
