// Zasoby wizualne i ich etykiety (tabele `zasob_design`,
// `etykieta_zasobu_design`) — obszar Assets Panel modułu Design. Prompt
// strukturalny leży w `design.go` (tam też interfejs całego obszaru),
// kompozycje w `design_kompozycje.go`.
//
// Zasoby oddają `Total` oddzielnie od strony: `design.asset.list` niesie
// `Total` obok `Assets` przyciętych limitem (`DesignAssetListResponse`), żeby
// panel pokazał „X z Y” bez drugiego zapytania po stronie klienta. `Zasoby`
// stosuje więc te same warunki dwa razy — raz do stronicowanej listy
// (z LIMIT), raz do liczby całkowitej (bez LIMIT).
//
// Filtr po etykietach jest koniunkcją. `DesignAssetListRequest.Tags` nie
// rozstrzyga, czy zasób ma nieść wszystkie wskazane etykiety, czy choć jedną,
// a panel filtrujący po wielu etykietach naraz zawęża wynik, więc `Zasoby`
// wymaga zestawu pełnego: `HAVING COUNT(DISTINCT etykieta) = len(etykiety)`.
//
// Etykiety są wymianą, nie dokładaniem — ten sam wzorzec co `UstawEtykiety`
// w `library_kolekcje.go`: `UstawEtykietyZasobu` usuwa zastane etykiety
// zasobu i wstawia nadesłane w jednej transakcji (`transakcja.go`).
//
// Zasób odczytany po kodzie jest osobnym wejściem (`Zasob`). Komenda
// `design.asset.tag.set` dostaje z zewnątrz identyfikator kontraktu, a
// `etykieta_zasobu_design.zasob_id` wskazuje klucz wiersza; bez tego przekładu
// uchwyt nie odróżni zasobu nieznanego (odmowa) od zasobu bez etykiet
// (droga udana).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ZasobDesignu to wiersz tabeli `zasob_design`. PromptID niesie wskaźnik, bo
// `prompt_id` dopuszcza NULL (ON DELETE SET NULL) — zasób przeżywa usunięcie
// promptu, z którego powstał.
type ZasobDesignu struct {
	ID       int64
	Kod      string
	Okno     string
	Nazwa    *string
	Rodzaj   string
	Format   *string
	URI      *string
	PromptID *int64
	// PromptKod jest identyfikatorem ZEWNĘTRZNYM promptu, z którego zasób
	// powstał — tym, którego chce kontrakt (`DesignAsset.PromptId`). PromptID
	// jest kluczem wiersza i na zewnątrz nie wychodzi. Odczyt bierze go
	// podzapytaniem obok wiersza zasobu, więc każdy czytelnik zasobu dostaje
	// prowenancję bez drugiego wywołania i bez zgadywania.
	PromptKod       *string
	WariantZasobuID *string
	Ulubiony        bool
	Szerokosc       *int
	Wysokosc        *int
	Utworzono       string
}

// FiltrZasobow niesie dokładnie to, czego wymaga `design.asset.list`
// (`DesignAssetListRequest`): okno, rodzaj, etykiety (koniunkcja), wyłącznie
// ulubione i granicę strony. Brak okna znaczy „wszystkie okna” — pole jest
// opcjonalne w żądaniu.
type FiltrZasobow struct {
	Okno          *string
	Rodzaj        *string
	Etykiety      []string
	TylkoUlubione bool
	Limit         int
}

const (
	kolumnyZasobuDesign = `z.id, z.identyfikator_zewnetrzny, z.okno, z.nazwa, z.rodzaj,
	                       z.format, z.uri, z.prompt_id,
	                       (SELECT p.identyfikator_zewnetrzny FROM prompt_design p
	                         WHERE p.id = z.prompt_id),
	                       z.wariant_zasobu_id, z.ulubiony,
	                       z.szerokosc, z.wysokosc, z.utworzono`

	// Zapis zakłada zasób albo nadpisuje zastany po identyfikatorze
	// zewnętrznym — regeneracja tego samego wariantu (np. ponowny odczyt po
	// zakończeniu procesu generowania) jest normalną ścieżką, nie usterką.
	zapiszZasobDesign = `INSERT INTO zasob_design
	                     (identyfikator_zewnetrzny, okno, nazwa, rodzaj, format, uri,
	                      prompt_id, wariant_zasobu_id, ulubiony, szerokosc, wysokosc)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                     ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                         okno = excluded.okno,
	                         nazwa = excluded.nazwa,
	                         rodzaj = excluded.rodzaj,
	                         format = excluded.format,
	                         uri = excluded.uri,
	                         prompt_id = excluded.prompt_id,
	                         wariant_zasobu_id = excluded.wariant_zasobu_id,
	                         ulubiony = excluded.ulubiony,
	                         szerokosc = excluded.szerokosc,
	                         wysokosc = excluded.wysokosc`

	pobierzZasobDesign = `SELECT ` + kolumnyZasobuDesign + ` FROM zasob_design z
	                      WHERE z.identyfikator_zewnetrzny = ?`

	usunEtykietyZasobuDesign = `DELETE FROM etykieta_zasobu_design WHERE zasob_id = ?`

	wstawEtykieteZasobuDesign = `INSERT INTO etykieta_zasobu_design (zasob_id, etykieta)
	                             VALUES (?, ?)
	                             ON CONFLICT(zasob_id, etykieta) DO NOTHING`

	listaEtykietZasobuDesign = `SELECT etykieta FROM etykieta_zasobu_design
	                            WHERE zasob_id = ? ORDER BY etykieta`

	// Oznaczenie ulubionego jest zapisem jednej kolumny — reszta wiersza
	// zostaje nietknięta (patrz `UstawUlubionyZasobu` w interfejsie,
	// `design.go`).
	ustawUlubionyZasobuDesign = `UPDATE zasob_design SET ulubiony = ? WHERE id = ?`

	// Usunięcie idzie po identyfikatorze zewnętrznym, bo tym wskazuje komenda.
	// Etykiety zasobu znikają same — `etykieta_zasobu_design.zasob_id` niesie
	// ON DELETE CASCADE, więc drugiego polecenia tu nie ma. Warstwy kompozycji
	// wskazujące ten zasób zostają: `warstwa_kompozycji_design.zasob_id` jest
	// kolumną TEXT bez klucza obcego, a kompozycja przeżywa usunięcie zasobu,
	// który się w niej znalazł.
	usunZasobDesign = `DELETE FROM zasob_design WHERE identyfikator_zewnetrzny = ?`
)

// ZapiszZasob zakłada zasób albo nadpisuje zastany po identyfikatorze
// zewnętrznym i zwraca stan po zapisie — obsługuje część zasobową
// `design.asset.generate`.
func (r *repozytoriumDesignu) ZapiszZasob(ctx context.Context, zasob ZasobDesignu) (ZasobDesignu, error) {
	if zasob.Kod == "" {
		return ZasobDesignu{}, fmt.Errorf("dane: zasob design bez identyfikatora")
	}
	if zasob.Okno == "" {
		return ZasobDesignu{}, fmt.Errorf("dane: zasob design %q bez okna", zasob.Kod)
	}
	if zasob.Rodzaj == "" {
		return ZasobDesignu{}, fmt.Errorf("dane: zasob design %q bez rodzaju", zasob.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZasobDesign)
	if err != nil {
		return ZasobDesignu{}, err
	}
	// Szerokosc i Wysokosc niesie kontrakt jako *int, a liczbaDoKolumny
	// przyjmuje *int64, stąd przekład wprost na wartość kolumny.
	var szerokosc, wysokosc any
	if zasob.Szerokosc != nil {
		szerokosc = int64(*zasob.Szerokosc)
	}
	if zasob.Wysokosc != nil {
		wysokosc = int64(*zasob.Wysokosc)
	}
	ulubiony := 0
	if zasob.Ulubiony {
		ulubiony = 1
	}
	_, err = polecenie.ExecContext(ctx, zasob.Kod, zasob.Okno, tekstDoKolumny(zasob.Nazwa),
		zasob.Rodzaj, tekstDoKolumny(zasob.Format), tekstDoKolumny(zasob.URI),
		liczbaDoKolumny(zasob.PromptID), tekstDoKolumny(zasob.WariantZasobuID), ulubiony,
		szerokosc, wysokosc)
	if err != nil {
		return ZasobDesignu{}, fmt.Errorf("dane: nie można zapisać zasobu design %q: %w", zasob.Kod, err)
	}
	return r.Zasob(ctx, zasob.Kod)
}

// Zasoby zwraca stronę zasobów spełniających filtr (od najnowszych) oraz
// liczbę wszystkich zasobów spełniających ten sam filtr, bez przycięcia
// limitem — druga wartość zasila `DesignAssetListResponse.Total`.
func (r *repozytoriumDesignu) Zasoby(ctx context.Context, filtr FiltrZasobow) ([]ZasobDesignu, int, error) {
	warunki, argumenty := warunkiFiltruZasobow(filtr)

	zapytanieStrony := `SELECT ` + kolumnyZasobuDesign + ` FROM zasob_design z` + warunki +
		` ORDER BY z.utworzono DESC, z.id DESC LIMIT ?`
	polecenieStrony, err := r.zapytania.przygotuj(ctx, zapytanieStrony)
	if err != nil {
		return nil, 0, err
	}
	wiersze, err := polecenieStrony.QueryContext(ctx, append(argumenty, granicaWykazu(filtr.Limit))...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać zasobów design: %w", err)
	}
	defer wiersze.Close()

	lista := []ZasobDesignu{}
	for wiersze.Next() {
		zasob, err := odczytajZasobDesign(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz zasobu design: %w", err)
		}
		lista = append(lista, zasob)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt zasobów design: %w", err)
	}

	zapytanieLiczby := `SELECT COUNT(*) FROM zasob_design z` + warunki
	polecenieLiczby, err := r.zapytania.przygotuj(ctx, zapytanieLiczby)
	if err != nil {
		return nil, 0, err
	}
	var razem int
	if err := polecenieLiczby.QueryRowContext(ctx, argumenty...).Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć zasobów design: %w", err)
	}
	return lista, razem, nil
}

// warunkiFiltruZasobow składa klauzulę WHERE i listę argumentów wspólną dla
// stronicowanego odczytu i liczenia całości — dwa zapytania muszą widzieć
// dokładnie te same warunki, inaczej Total i długość strony rozjadą się
// pozornie losowo.
func warunkiFiltruZasobow(filtr FiltrZasobow) (string, []any) {
	warunki := []string{}
	argumenty := []any{}

	if filtr.Okno != nil {
		warunki = append(warunki, "z.okno = ?")
		argumenty = append(argumenty, *filtr.Okno)
	}
	if filtr.Rodzaj != nil {
		warunki = append(warunki, "z.rodzaj = ?")
		argumenty = append(argumenty, *filtr.Rodzaj)
	}
	if filtr.TylkoUlubione {
		warunki = append(warunki, "z.ulubiony = 1")
	}
	if len(filtr.Etykiety) > 0 {
		// Zasób musi nieść wszystkie wskazane etykiety (koniunkcja) — patrz
		// komentarz nagłówkowy pliku.
		zaslepki := strings.TrimSuffix(strings.Repeat("?,", len(filtr.Etykiety)), ",")
		warunki = append(warunki, fmt.Sprintf(
			`z.id IN (SELECT zasob_id FROM etykieta_zasobu_design
			          WHERE etykieta IN (%s)
			          GROUP BY zasob_id
			          HAVING COUNT(DISTINCT etykieta) = ?)`, zaslepki))
		for _, etykieta := range filtr.Etykiety {
			argumenty = append(argumenty, etykieta)
		}
		argumenty = append(argumenty, len(filtr.Etykiety))
	}

	if len(warunki) == 0 {
		return "", argumenty
	}
	return " WHERE " + strings.Join(warunki, " AND "), argumenty
}

// UstawEtykietyZasobu podmienia komplet etykiet zasobu, bo Assets Panel
// nadsyła zawsze pełny zestaw, nie różnicę (kontrakt: `DesignAsset.Tags`
// niesie stan docelowy). Usunięcie i wstawienie zachodzi w jednej
// transakcji — zasób nie zostaje przejściowo bez etykiet przy błędzie
// w trakcie.
func (r *repozytoriumDesignu) UstawEtykietyZasobu(ctx context.Context,
	zasobID int64, etykiety []string) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunEtykietyZasobuDesign)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, zasobID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić etykiet zasobu %d: %w", zasobID, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawEtykieteZasobuDesign)
		if err != nil {
			return err
		}
		for _, etykieta := range etykiety {
			if etykieta == "" {
				continue
			}
			if _, err := wstawienie.ExecContext(ctx, zasobID, etykieta); err != nil {
				return fmt.Errorf("dane: nie można zapisać etykiety %q zasobu %d: %w",
					etykieta, zasobID, err)
			}
		}
		return nil
	})
}

// EtykietyZasobu zwraca etykiety zasobu w porządku alfabetycznym.
func (r *repozytoriumDesignu) EtykietyZasobu(ctx context.Context, zasobID int64) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaEtykietZasobuDesign)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, zasobID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać etykiet zasobu %d: %w", zasobID, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var etykieta string
		if err := wiersze.Scan(&etykieta); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz etykiety zasobu %d: %w", zasobID, err)
		}
		lista = append(lista, etykieta)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt etykiet zasobu %d: %w", zasobID, err)
	}
	return lista, nil
}

// Zasob zwraca zasób o wskazanym identyfikatorze zewnętrznym (tym z
// kontraktu, nie kluczu wiersza). Brak wiersza wraca jako ErrBrakWiersza —
// warstwa wyższa odróżnia „nie ma” od „odczyt się nie powiódł”, bo tylko
// pierwsze z tego jest odmową kontraktu.
func (r *repozytoriumDesignu) Zasob(ctx context.Context, kod string) (ZasobDesignu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZasobDesign)
	if err != nil {
		return ZasobDesignu{}, err
	}
	zasob, err := odczytajZasobDesign(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ZasobDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return ZasobDesignu{}, fmt.Errorf("dane: nieczytelny wiersz zasobu design %q: %w", kod, err)
	}
	return zasob, nil
}

// odczytajZasobDesign składa strukturę z jednego wiersza wyniku.
func odczytajZasobDesign(wiersz skaner) (ZasobDesignu, error) {
	var zasob ZasobDesignu
	var nazwa, format, uri, promptKod, wariantZasobuID sql.NullString
	var promptID sql.NullInt64
	var ulubiony int
	var szerokosc, wysokosc sql.NullInt64
	err := wiersz.Scan(&zasob.ID, &zasob.Kod, &zasob.Okno, &nazwa, &zasob.Rodzaj, &format,
		&uri, &promptID, &promptKod, &wariantZasobuID, &ulubiony, &szerokosc, &wysokosc,
		&zasob.Utworzono)
	if err != nil {
		return ZasobDesignu{}, err
	}
	zasob.Nazwa = tekstZKolumny(nazwa)
	zasob.Format = tekstZKolumny(format)
	zasob.URI = tekstZKolumny(uri)
	zasob.WariantZasobuID = tekstZKolumny(wariantZasobuID)
	zasob.PromptID = liczbaZKolumny(promptID)
	zasob.PromptKod = tekstZKolumny(promptKod)
	zasob.Ulubiony = ulubiony != 0
	if szerokosc.Valid {
		wartosc := int(szerokosc.Int64)
		zasob.Szerokosc = &wartosc
	}
	if wysokosc.Valid {
		wartosc := int(wysokosc.Int64)
		zasob.Wysokosc = &wartosc
	}
	return zasob, nil
}
