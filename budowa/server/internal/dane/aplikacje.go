// Plik definiuje architekturę produktu i jej komponenty: kontrakt
// RepozytoriumAplikacji obszaru Apps, wspólny z aplikacje_warsztat.go i
// aplikacje_wdrozenia.go.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// ArchitekturaApp to wiersz tabeli `architektura_apps`. Kod jest
// identyfikatorem, którym architektura wychodzi kontraktem (`AppArchitecture.id`).
type ArchitekturaApp struct {
	ID                    int64
	Kod                   string
	Okno                  string
	Nazwa                 *string
	Szablon               string
	Wersja                int
	ZastrzezeniaWalidacji []string
	Utworzono             string
	Zaktualizowano        string
	// RoznicaWersji jest polem wyłącznie zapisu: różnica wobec poprzedniej
	// wersji idzie do historii.
	RoznicaWersji *string
}

// KomponentArchitektury to wiersz tabeli `komponent_architektury_apps`.
// Rodzaj jest wartością kontraktu (AppComponentKind), bez tłumaczenia.
type KomponentArchitektury struct {
	ID            int64
	KodZewnetrzny string
	Nazwa         string
	Rodzaj        string
	Stos          *string
	Opis          *string
	KontraktAPI   *string
}

// ZaleznoscKomponentu to jeden łuk grafu zależności między komponentami
// jednej architektury — pole `AppComponent.DependsOn` kontraktu powstaje
// z tych wierszy.
type ZaleznoscKomponentu struct {
	KomponentZ  string
	KomponentDo string
}

// RepozytoriumAplikacji jest kontraktem obszaru Apps: architektura,
// warsztat, wdrożenia, produkt, środowiska i pozostałe zasoby modułu.
type RepozytoriumAplikacji interface {
	// --- architektura ---
	ZapiszArchitekture(ctx context.Context, architektura ArchitekturaApp,
		komponenty []KomponentArchitektury, zaleznosci []ZaleznoscKomponentu) (ArchitekturaApp, error)
	Architektura(ctx context.Context, kod string) (ArchitekturaApp, error)
	// ArchitekturaOkna zwraca architekturę okna dla apps.architecture.get,
	// po oknie, nie po kodzie.
	ArchitekturaOkna(ctx context.Context, okno string) (ArchitekturaApp, error)
	Komponenty(ctx context.Context, architekturaID int64) ([]KomponentArchitektury, error)
	ZaleznosciKomponentow(ctx context.Context, architekturaID int64) ([]ZaleznoscKomponentu, error)

	// --- warsztat ---
	ZapiszPlikWarsztatu(ctx context.Context, plik PlikWarsztatu) (PlikWarsztatu, error)
	PlikiWarsztatu(ctx context.Context, okno string) ([]PlikWarsztatu, error)
	// PlikWarsztatu zwraca jeden plik po kluczu naturalnym, do rozróżnienia
	// powstania od zmiany pliku.
	PlikWarsztatu(ctx context.Context, okno string, warstwa shared.AppWorkspaceLayer,
		sciezka string) (PlikWarsztatu, error)

	// --- wdrożenia ---
	ZapiszWdrozenie(ctx context.Context, wdrozenie WdrozenieApp) (WdrozenieApp, error)
	Wdrozenie(ctx context.Context, kod string) (WdrozenieApp, error)
	// Wdrozenia zwraca stronę wdrożeń okna oraz liczbę wszystkich
	// spełniających te same warunki.
	Wdrozenia(ctx context.Context, okno string, srodowisko *shared.AppDeployEnvironment,
		limit int) ([]WdrozenieApp, int, error)

	// --- produkt, etapy, kamienie milowe (aplikacje_produkt.go) ---
	ZapiszProduktApp(ctx context.Context, produkt ProduktApp) (ProduktApp, error)
	ProduktApp(ctx context.Context, okno string) (ProduktApp, error)
	ZapiszEtapApp(ctx context.Context, etap EtapApp) (EtapApp, error)
	EtapApp(ctx context.Context, kod string) (EtapApp, error)
	EtapyApp(ctx context.Context, okno string) ([]EtapApp, error)
	ZapiszKamienMilowyApp(ctx context.Context, kamien KamienMilowyApp) (KamienMilowyApp, error)
	KamienMilowyApp(ctx context.Context, kod string) (KamienMilowyApp, error)
	KamienieMiloweApp(ctx context.Context, okno string) ([]KamienMilowyApp, error)
	UsunKamienMilowyApp(ctx context.Context, kod string) (bool, error)

	// --- historia wersji i adnotacje (aplikacje_architektura_wersje.go) ---
	WersjeArchitekturyApp(ctx context.Context, architekturaID int64) ([]WersjaArchitekturyApp, error)
	ZapiszAdnotacjeApp(ctx context.Context, adnotacja AdnotacjaArchitekturyApp) (AdnotacjaArchitekturyApp, error)
	AdnotacjaApp(ctx context.Context, kod string) (AdnotacjaArchitekturyApp, error)
	AdnotacjeApp(ctx context.Context, architekturaID int64) ([]AdnotacjaArchitekturyApp, error)

	// --- środowiska, zmienne, skalowanie, kondycja (aplikacje_srodowiska.go) ---
	ZapiszSrodowiskoApp(ctx context.Context, srodowisko SrodowiskoApp) (SrodowiskoApp, error)
	SrodowiskoApp(ctx context.Context, okno, kod string) (SrodowiskoApp, error)
	SrodowiskaApp(ctx context.Context, okno string) ([]SrodowiskoApp, error)
	ZapiszZmiennaSrodowiskaApp(ctx context.Context, zmienna ZmiennaSrodowiskaApp) (ZmiennaSrodowiskaApp, error)
	ZmienneSrodowiskaApp(ctx context.Context, okno, srodowisko string) ([]ZmiennaSrodowiskaApp, error)
	ZapiszSkalowanieApp(ctx context.Context, skalowanie SkalowanieApp) (SkalowanieApp, error)
	SkalowanieApp(ctx context.Context, okno, srodowisko string) (SkalowanieApp, error)
	ZapiszKondycjeApp(ctx context.Context, kondycja KondycjaWdrozeniaApp) error
	KondycjeApp(ctx context.Context, okno, srodowisko string, granica int) ([]KondycjaWdrozeniaApp, error)

	// --- podgląd, motyw, dzienniki, artefakty, pakiety (aplikacje_wytwory.go) ---
	ZapiszPodgladApp(ctx context.Context, podglad PodgladApp) error
	PodgladApp(ctx context.Context, okno string) (PodgladApp, error)
	ZapiszMotywApp(ctx context.Context, motyw MotywApp) (MotywApp, error)
	MotywApp(ctx context.Context, okno string) (MotywApp, error)
	DopiszWierszDziennikaApp(ctx context.Context, wiersz WierszDziennikaApp) error
	DziennikUslugiApp(ctx context.Context, okno, komponent string, od int64, granica int) ([]WierszDziennikaApp, int, error)
	DziennikWdrozeniaApp(ctx context.Context, wdrozenie string, od int64, granica int) ([]WierszDziennikaApp, int, error)
	ZalozArtefaktApp(ctx context.Context, artefakt ArtefaktApp) (ArtefaktApp, error)
	ArtefaktApp(ctx context.Context, kod string) (ArtefaktApp, error)
	ArtefaktyApp(ctx context.Context, okno, wdrozenie string) ([]ArtefaktApp, error)
	OstatniArtefaktUdanegoWdrozeniaApp(ctx context.Context, okno string) (ArtefaktApp, error)
	ZapiszPakietApp(ctx context.Context, pakiet PakietApp) (PakietApp, error)
	PakietApp(ctx context.Context, kod string) (PakietApp, error)
	PakietyApp(ctx context.Context, okno string) ([]PakietApp, error)
}

const (
	kolumnyArchitekturyApp = `id, identyfikator_zewnetrzny, okno, nazwa, szablon, wersja,
	                          zastrzezenia_walidacji, utworzono, zaktualizowano`

	pobierzArchitektureApp = `SELECT ` + kolumnyArchitekturyApp + ` FROM architektura_apps
	                          WHERE identyfikator_zewnetrzny = ?`

	// Okno bierze architekturę najświeższą: w oknie może leżeć więcej niż
	// jedna definicja, a ostatnio zapisana jest tą, którą pokazuje panel.
	pobierzArchitektureOknaApp = `SELECT ` + kolumnyArchitekturyApp + ` FROM architektura_apps
	                          WHERE okno = ?
	                          ORDER BY zaktualizowano DESC, id DESC
	                          LIMIT 1`

	wstawArchitektureApp = `INSERT INTO architektura_apps
	                        (identyfikator_zewnetrzny, okno, nazwa, szablon, wersja, zastrzezenia_walidacji)
	                        VALUES (?, ?, ?, ?, ?, ?)
	                        ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                            nazwa = excluded.nazwa,
	                            szablon = excluded.szablon,
	                            wersja = architektura_apps.wersja + 1,
	                            zastrzezenia_walidacji = excluded.zastrzezenia_walidacji,
	                            zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	usunKomponentyArchitektury = `DELETE FROM komponent_architektury_apps WHERE architektura_id = ?`

	wstawKomponentArchitektury = `INSERT INTO komponent_architektury_apps
	                              (architektura_id, identyfikator_zewnetrzny, nazwa, rodzaj,
	                               stos, opis, kontrakt_api)
	                              VALUES (?, ?, ?, ?, ?, ?, ?)`

	listaKomponentowArchitektury = `SELECT id, identyfikator_zewnetrzny, nazwa, rodzaj,
	                                        stos, opis, kontrakt_api
	                                FROM komponent_architektury_apps WHERE architektura_id = ?
	                                ORDER BY id`

	usunZaleznosciKomponentow = `DELETE FROM zaleznosc_komponentu_apps WHERE architektura_id = ?`

	wstawZaleznoscKomponentu = `INSERT INTO zaleznosc_komponentu_apps
	                            (architektura_id, komponent_z, komponent_do)
	                            VALUES (?, ?, ?)
	                            ON CONFLICT(architektura_id, komponent_z, komponent_do) DO NOTHING`

	listaZaleznosciKomponentow = `SELECT komponent_z, komponent_do FROM zaleznosc_komponentu_apps
	                              WHERE architektura_id = ? ORDER BY komponent_do, komponent_z`
)

type repozytoriumAplikacji struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumAplikacji(z *zapytania, db *sql.DB) *repozytoriumAplikacji {
	return &repozytoriumAplikacji{zapytania: z, db: db}
}

// ZapiszArchitekture zapisuje definicję architektury oraz wymienia komplet
// jej komponentów i zależności w jednej transakcji, usuwając zastane i
// wstawiając od nowa.
func (r *repozytoriumAplikacji) ZapiszArchitekture(ctx context.Context, architektura ArchitekturaApp,
	komponenty []KomponentArchitektury, zaleznosci []ZaleznoscKomponentu) (ArchitekturaApp, error) {

	if architektura.Kod == "" {
		return ArchitekturaApp{}, fmt.Errorf("dane: architektura aplikacji bez identyfikatora")
	}
	if architektura.Okno == "" {
		return ArchitekturaApp{}, fmt.Errorf("dane: architektura aplikacji %q bez okna", architektura.Kod)
	}
	szablon := architektura.Szablon
	if szablon == "" {
		szablon = "monolith"
	}
	wersja := architektura.Wersja
	if wersja == 0 {
		wersja = 1
	}

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, wstawArchitektureApp)
		if err != nil {
			return err
		}
		_, err = zapis.ExecContext(ctx, architektura.Kod, architektura.Okno,
			tekstDoKolumny(architektura.Nazwa), szablon, wersja,
			listaDoKolumny(architektura.ZastrzezeniaWalidacji))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać architektury %q: %w", architektura.Kod, err)
		}

		var architekturaID int64
		wiersz := transakcja.QueryRowContext(ctx,
			`SELECT id FROM architektura_apps WHERE identyfikator_zewnetrzny = ?`, architektura.Kod)
		if err := wiersz.Scan(&architekturaID); err != nil {
			return fmt.Errorf("dane: nie można odczytać id architektury %q: %w", architektura.Kod, err)
		}

		czyszczenieKomponentow, err := r.zapytania.wTransakcji(ctx, transakcja, usunKomponentyArchitektury)
		if err != nil {
			return err
		}
		if _, err := czyszczenieKomponentow.ExecContext(ctx, architekturaID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić komponentów architektury %d: %w", architekturaID, err)
		}
		wstawienieKomponentow, err := r.zapytania.wTransakcji(ctx, transakcja, wstawKomponentArchitektury)
		if err != nil {
			return err
		}
		for _, komponent := range komponenty {
			_, err := wstawienieKomponentow.ExecContext(ctx, architekturaID, komponent.KodZewnetrzny,
				komponent.Nazwa, komponent.Rodzaj, tekstDoKolumny(komponent.Stos),
				tekstDoKolumny(komponent.Opis), tekstDoKolumny(komponent.KontraktAPI))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać komponentu %q architektury %d: %w",
					komponent.KodZewnetrzny, architekturaID, err)
			}
		}

		czyszczenieZaleznosci, err := r.zapytania.wTransakcji(ctx, transakcja, usunZaleznosciKomponentow)
		if err != nil {
			return err
		}
		if _, err := czyszczenieZaleznosci.ExecContext(ctx, architekturaID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić zależności architektury %d: %w", architekturaID, err)
		}
		wstawienieZaleznosci, err := r.zapytania.wTransakcji(ctx, transakcja, wstawZaleznoscKomponentu)
		if err != nil {
			return err
		}
		for _, zaleznosc := range zaleznosci {
			_, err := wstawienieZaleznosci.ExecContext(ctx, architekturaID,
				zaleznosc.KomponentZ, zaleznosc.KomponentDo)
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać zależności %q→%q architektury %d: %w",
					zaleznosc.KomponentZ, zaleznosc.KomponentDo, architekturaID, err)
			}
		}

		// Wiersz historii idzie tą samą transakcją co wymiana komponentów i
		// niesie ich liczbę po zapisie.
		var wersjaPoZapisie int
		wiersz = transakcja.QueryRowContext(ctx,
			`SELECT wersja FROM architektura_apps WHERE id = ?`, architekturaID)
		if err := wiersz.Scan(&wersjaPoZapisie); err != nil {
			return fmt.Errorf("dane: nie można odczytać wersji architektury %d: %w", architekturaID, err)
		}
		historia, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWersjeArchitekturyApp)
		if err != nil {
			return err
		}
		if _, err := historia.ExecContext(ctx, architekturaID, wersjaPoZapisie,
			len(komponenty), tekstDoKolumny(architektura.RoznicaWersji)); err != nil {
			return fmt.Errorf("dane: nie można zapisać wersji %d architektury %d: %w",
				wersjaPoZapisie, architekturaID, err)
		}
		return nil
	})
	if err != nil {
		return ArchitekturaApp{}, err
	}
	return r.Architektura(ctx, architektura.Kod)
}

// Architektura zwraca architekturę o wskazanym kodzie. Brak wiersza wraca
// jako ErrBrakWiersza — warstwa wyższa odróżnia „nie ma" od „odczyt się nie powiódł".
func (r *repozytoriumAplikacji) Architektura(ctx context.Context, kod string) (ArchitekturaApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzArchitektureApp)
	if err != nil {
		return ArchitekturaApp{}, err
	}
	architektura, err := odczytajArchitektureApp(polecenie.QueryRowContext(ctx, kod))
	if err == sql.ErrNoRows {
		return ArchitekturaApp{}, ErrBrakWiersza
	}
	if err != nil {
		return ArchitekturaApp{}, fmt.Errorf("dane: nieczytelny wiersz architektury %q: %w", kod, err)
	}
	return architektura, nil
}

// ArchitekturaOkna zwraca najświeższą architekturę okna. Okno bez ani
// jednej architektury wraca jako ErrBrakWiersza, co jest stanem normalnym
// świeżego okna.
func (r *repozytoriumAplikacji) ArchitekturaOkna(ctx context.Context, okno string) (ArchitekturaApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzArchitektureOknaApp)
	if err != nil {
		return ArchitekturaApp{}, err
	}
	architektura, err := odczytajArchitektureApp(polecenie.QueryRowContext(ctx, okno))
	if err == sql.ErrNoRows {
		return ArchitekturaApp{}, ErrBrakWiersza
	}
	if err != nil {
		return ArchitekturaApp{}, fmt.Errorf("dane: nieczytelna architektura okna %q: %w", okno, err)
	}
	return architektura, nil
}

// Komponenty zwraca komponenty architektury w kolejności zapisu, od
// pierwszego wstawionego do ostatniego.
func (r *repozytoriumAplikacji) Komponenty(ctx context.Context, architekturaID int64) ([]KomponentArchitektury, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKomponentowArchitektury)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, architekturaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać komponentów architektury %d: %w", architekturaID, err)
	}
	defer wiersze.Close()

	lista := []KomponentArchitektury{}
	for wiersze.Next() {
		var komponent KomponentArchitektury
		var stos, opis, kontraktAPI sql.NullString
		err := wiersze.Scan(&komponent.ID, &komponent.KodZewnetrzny, &komponent.Nazwa,
			&komponent.Rodzaj, &stos, &opis, &kontraktAPI)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz komponentu architektury: %w", err)
		}
		komponent.Stos = tekstZKolumny(stos)
		komponent.Opis = tekstZKolumny(opis)
		komponent.KontraktAPI = tekstZKolumny(kontraktAPI)
		lista = append(lista, komponent)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt komponentów architektury %d: %w", architekturaID, err)
	}
	return lista, nil
}

// ZaleznosciKomponentow zwraca graf zależności między komponentami
// architektury jako listę łuków źródło-cel.
func (r *repozytoriumAplikacji) ZaleznosciKomponentow(ctx context.Context, architekturaID int64) ([]ZaleznoscKomponentu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaZaleznosciKomponentow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, architekturaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zależności architektury %d: %w", architekturaID, err)
	}
	defer wiersze.Close()

	lista := []ZaleznoscKomponentu{}
	for wiersze.Next() {
		var zaleznosc ZaleznoscKomponentu
		if err := wiersze.Scan(&zaleznosc.KomponentZ, &zaleznosc.KomponentDo); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zależności architektury: %w", err)
		}
		lista = append(lista, zaleznosc)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zależności architektury %d: %w", architekturaID, err)
	}
	return lista, nil
}

// odczytajArchitektureApp składa strukturę ArchitekturaApp z jednego wiersza
// wyniku zapytania, zamieniając kolumny nullowalne na wskaźniki.
func odczytajArchitektureApp(wiersz skaner) (ArchitekturaApp, error) {
	var architektura ArchitekturaApp
	var nazwa, zastrzezenia sql.NullString
	err := wiersz.Scan(&architektura.ID, &architektura.Kod, &architektura.Okno, &nazwa,
		&architektura.Szablon, &architektura.Wersja, &zastrzezenia,
		&architektura.Utworzono, &architektura.Zaktualizowano)
	if err != nil {
		return ArchitekturaApp{}, err
	}
	architektura.Nazwa = tekstZKolumny(nazwa)
	architektura.ZastrzezeniaWalidacji = listaZKolumny(zastrzezenia)
	return architektura, nil
}

// listaDoKolumny scala zastrzeżenia walidacji w jeden tekst rozdzielony
// znakiem nowej linii — ten sam wzorzec co `argumenty` w `developer_budowanie`
// (migracja 042); nikt nie filtruje ani nie sortuje po pojedynczym zastrzeżeniu.
func listaDoKolumny(wartosci []string) any {
	if len(wartosci) == 0 {
		return nil
	}
	return strings.Join(wartosci, "\n")
}

// listaZKolumny odwraca listaDoKolumny, rozdzielając zapisany tekst z
// powrotem na listę zastrzeżeń walidacji.
func listaZKolumny(kolumna sql.NullString) []string {
	if !kolumna.Valid || kolumna.String == "" {
		return nil
	}
	return strings.Split(kolumna.String, "\n")
}
