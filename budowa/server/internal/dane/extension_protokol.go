// Rodzina extension.*: warstwa protokołu. Narzędzia odkryte u integracji,
// dziennik ramek JSON-RPC, wywołania wraz z ich czasem oraz wyniki sprawdzeń
// kondycji; wszystkie cztery byty powstają z pracy, nie z żądania odczytu.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// NarzedzieRozszerzenia to wiersz tabeli `narzedzie_rozszerzenia` — jeden wpis
// jednego z trzech wykazów protokołu MCP.
type NarzedzieRozszerzenia struct {
	ID              int64
	RozszerzenieKod string
	Rodzaj          string
	Nazwa           string
	Opis            *string
	SchematWejscia  *string
	Adres           *string
	Odkryto         int64
	WersjaProtokolu *string
}

// RamkaProtokolu to wiersz tabeli `ramka_protokolu_rozszerzenia`: jedna
// ramka JSON-RPC odnotowana w dzienniku wywołania.
type RamkaProtokolu struct {
	ID              int64
	Kod             string
	RozszerzenieKod string
	Kierunek        string
	Metoda          *string
	Korelacja       *string
	Tresc           *string
	KodBledu        *string
	Zaszlo          int64
}

// WywolanieRozszerzenia to wiersz tabeli `wywolanie_rozszerzenia`: jedno
// wywołanie narzędzia integracji wraz z jego czasem i wynikiem.
type WywolanieRozszerzenia struct {
	ID              int64
	RozszerzenieKod string
	AgentKod        *string
	Narzedzie       *string
	Udane           bool
	CzasMs          int64
	Szczegol        *string
	Zaszlo          int64
}

// KondycjaRozszerzenia to wiersz tabeli `kondycja_rozszerzenia`: wynik
// jednego sprawdzenia kondycji integracji.
type KondycjaRozszerzenia struct {
	ID              int64
	RozszerzenieKod string
	Stan            string
	CzasMs          *int64
	LiczbaNarzedzi  *int64
	BladPowitania   *string
	Sprawdzono      int64
}

const (
	usunNarzedziaRozszerzenia = `DELETE FROM narzedzie_rozszerzenia WHERE rozszerzenie_kod = ?`

	wstawNarzedzieRozszerzenia = `INSERT INTO narzedzie_rozszerzenia
	                              (rozszerzenie_kod, rodzaj, nazwa, opis, schemat_wejscia,
	                               adres, odkryto, wersja_protokolu)
	                              VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	                              ON CONFLICT(rozszerzenie_kod, rodzaj, nazwa) DO UPDATE SET
	                                  opis = excluded.opis,
	                                  schemat_wejscia = excluded.schemat_wejscia,
	                                  adres = excluded.adres,
	                                  odkryto = excluded.odkryto,
	                                  wersja_protokolu = excluded.wersja_protokolu`

	listaNarzedziRozszerzenia = `SELECT id, rozszerzenie_kod, rodzaj, nazwa, opis,
	                                    schemat_wejscia, adres, odkryto, wersja_protokolu
	                             FROM narzedzie_rozszerzenia
	                             WHERE rozszerzenie_kod = ? AND (? = '' OR rodzaj = ?)
	                             ORDER BY rodzaj, nazwa`

	wstawRamkeProtokolu = `INSERT INTO ramka_protokolu_rozszerzenia
	                       (identyfikator_zewnetrzny, rozszerzenie_kod, kierunek, metoda,
	                        korelacja, tresc, kod_bledu, zaszlo)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	listaRamekProtokolu = `SELECT id, identyfikator_zewnetrzny, rozszerzenie_kod, kierunek,
	                              metoda, korelacja, tresc, kod_bledu, zaszlo
	                       FROM ramka_protokolu_rozszerzenia
	                       WHERE rozszerzenie_kod = ? AND zaszlo >= ?
	                       ORDER BY zaszlo DESC, id DESC LIMIT ?`

	policzRamkiProtokolu = `SELECT COUNT(*) FROM ramka_protokolu_rozszerzenia
	                        WHERE rozszerzenie_kod = ? AND zaszlo >= ?`

	wstawWywolanieRozszerzenia = `INSERT INTO wywolanie_rozszerzenia
	                              (rozszerzenie_kod, agent_kod, narzedzie, udane, czas_ms,
	                               szczegol, zaszlo)
	                              VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Metryka użycia liczy się z wierszy, nie z licznika: liczba wywołań,
	// liczba niepowodzeń i średni czas w oknie czasu wymagają trzech różnych
	// agregatów nad tym samym zbiorem.
	metrykaWywolanRozszerzenia = `SELECT rozszerzenie_kod, COUNT(*),
	                                     SUM(CASE WHEN udane = 0 THEN 1 ELSE 0 END),
	                                     AVG(czas_ms)
	                              FROM wywolanie_rozszerzenia
	                              WHERE (? = '' OR rozszerzenie_kod = ?)
	                                AND zaszlo >= ? AND zaszlo <= ?
	                              GROUP BY rozszerzenie_kod ORDER BY rozszerzenie_kod`

	listaAudytuRozszerzenia = `SELECT rozszerzenie_kod, agent_kod, narzedzie, zaszlo
	                           FROM wywolanie_rozszerzenia
	                           WHERE (? = '' OR rozszerzenie_kod = ?)
	                             AND (? = '' OR agent_kod = ?) AND zaszlo >= ?
	                           ORDER BY zaszlo DESC, id DESC LIMIT ?`

	policzAudytRozszerzenia = `SELECT COUNT(*) FROM wywolanie_rozszerzenia
	                           WHERE (? = '' OR rozszerzenie_kod = ?)
	                             AND (? = '' OR agent_kod = ?) AND zaszlo >= ?`

	wstawKondycjeRozszerzenia = `INSERT INTO kondycja_rozszerzenia
	                             (rozszerzenie_kod, stan, czas_ms, liczba_narzedzi,
	                              blad_powitania, sprawdzono)
	                             VALUES (?, ?, ?, ?, ?, ?)`

	ostatniaKondycjaRozszerzenia = `SELECT id, rozszerzenie_kod, stan, czas_ms, liczba_narzedzi,
	                                       blad_powitania, sprawdzono
	                                FROM kondycja_rozszerzenia WHERE rozszerzenie_kod = ?
	                                ORDER BY sprawdzono DESC, id DESC LIMIT 1`
)

// ZapiszNarzedziaRozszerzenia wymienia komplet wpisów odkrytych u integracji.
// Wymiana, nie dokładanie: serwer, który przestał udostępniać narzędzie, ma
// przestać je pokazywać, a wpis pozostawiony byłby obietnicą bez pokrycia.
func (r *repozytoriumRozszerzen) ZapiszNarzedziaRozszerzenia(ctx context.Context,
	rozszerzenie string, wpisy []NarzedzieRozszerzenia) error {

	if rozszerzenie == "" {
		return fmt.Errorf("dane: wykaz narzędzi bez pozycji katalogu")
	}
	return wTransakcji(ctx, r.baza, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunNarzedziaRozszerzenia)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, rozszerzenie); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić narzędzi pozycji %q: %w", rozszerzenie, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawNarzedzieRozszerzenia)
		if err != nil {
			return err
		}
		for _, wpis := range wpisy {
			_, err := wstawienie.ExecContext(ctx, rozszerzenie, wpis.Rodzaj, wpis.Nazwa,
				tekstDoKolumny(wpis.Opis), tekstDoKolumny(wpis.SchematWejscia),
				tekstDoKolumny(wpis.Adres), wpis.Odkryto, tekstDoKolumny(wpis.WersjaProtokolu))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać narzędzia %q pozycji %q: %w",
					wpis.Nazwa, rozszerzenie, err)
			}
		}
		return nil
	})
}

// NarzedziaRozszerzenia zwraca wpisy pozycji katalogu, opcjonalnie zawężone
// do jednego rodzaju narzędzia.
func (r *repozytoriumRozszerzen) NarzedziaRozszerzenia(ctx context.Context,
	rozszerzenie, rodzaj string) ([]NarzedzieRozszerzenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaNarzedziRozszerzenia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rozszerzenie, rodzaj, rodzaj)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać narzędzi pozycji %q: %w", rozszerzenie, err)
	}
	defer wiersze.Close()

	lista := []NarzedzieRozszerzenia{}
	for wiersze.Next() {
		var wpis NarzedzieRozszerzenia
		var opis, schemat, adres, protokol sql.NullString
		err := wiersze.Scan(&wpis.ID, &wpis.RozszerzenieKod, &wpis.Rodzaj, &wpis.Nazwa,
			&opis, &schemat, &adres, &wpis.Odkryto, &protokol)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz narzędzia rozszerzenia: %w", err)
		}
		wpis.Opis = tekstZKolumny(opis)
		wpis.SchematWejscia = tekstZKolumny(schemat)
		wpis.Adres = tekstZKolumny(adres)
		wpis.WersjaProtokolu = tekstZKolumny(protokol)
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt narzędzi pozycji %q: %w", rozszerzenie, err)
	}
	return lista, nil
}

// DopiszRamkeProtokolu odnotowuje jedną ramkę JSON-RPC w dzienniku
// protokołu wskazanej pozycji katalogu.
func (r *repozytoriumRozszerzen) DopiszRamkeProtokolu(ctx context.Context, ramka RamkaProtokolu) error {
	if ramka.Kod == "" || ramka.RozszerzenieKod == "" {
		return fmt.Errorf("dane: ramka protokołu bez identyfikatora albo pozycji")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawRamkeProtokolu)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, ramka.Kod, ramka.RozszerzenieKod, ramka.Kierunek,
		tekstDoKolumny(ramka.Metoda), tekstDoKolumny(ramka.Korelacja),
		tekstDoKolumny(ramka.Tresc), tekstDoKolumny(ramka.KodBledu), ramka.Zaszlo)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać ramki protokołu %q: %w", ramka.Kod, err)
	}
	return nil
}

// RamkiProtokolu zwraca stronę dziennika ramek protokołu wraz z liczbą
// wszystkich pasujących wierszy dziennika.
func (r *repozytoriumRozszerzen) RamkiProtokolu(ctx context.Context, rozszerzenie string,
	od int64, granica int) ([]RamkaProtokolu, int, error) {

	if granica <= 0 {
		granica = 200
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaRamekProtokolu)
	if err != nil {
		return nil, 0, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rozszerzenie, od, granica)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać ramek pozycji %q: %w", rozszerzenie, err)
	}
	defer wiersze.Close()

	lista := []RamkaProtokolu{}
	for wiersze.Next() {
		var ramka RamkaProtokolu
		var metoda, korelacja, tresc, blad sql.NullString
		err := wiersze.Scan(&ramka.ID, &ramka.Kod, &ramka.RozszerzenieKod, &ramka.Kierunek,
			&metoda, &korelacja, &tresc, &blad, &ramka.Zaszlo)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz ramki protokołu: %w", err)
		}
		ramka.Metoda = tekstZKolumny(metoda)
		ramka.Korelacja = tekstZKolumny(korelacja)
		ramka.Tresc = tekstZKolumny(tresc)
		ramka.KodBledu = tekstZKolumny(blad)
		lista = append(lista, ramka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt ramek pozycji %q: %w", rozszerzenie, err)
	}

	liczenie, err := r.zapytania.przygotuj(ctx, policzRamkiProtokolu)
	if err != nil {
		return nil, 0, err
	}
	var razem int
	if err := liczenie.QueryRowContext(ctx, rozszerzenie, od).Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć ramek pozycji %q: %w", rozszerzenie, err)
	}
	return lista, razem, nil
}

// DopiszWywolanieRozszerzenia odnotowuje jedno wywołanie narzędzia
// integracji wraz z jego wynikiem i czasem.
func (r *repozytoriumRozszerzen) DopiszWywolanieRozszerzenia(ctx context.Context,
	wywolanie WywolanieRozszerzenia) error {

	if wywolanie.RozszerzenieKod == "" {
		return fmt.Errorf("dane: wywołanie rozszerzenia bez pozycji katalogu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawWywolanieRozszerzenia)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wywolanie.RozszerzenieKod,
		tekstDoKolumny(wywolanie.AgentKod), tekstDoKolumny(wywolanie.Narzedzie),
		liczbaLogiczna(wywolanie.Udane), wywolanie.CzasMs,
		tekstDoKolumny(wywolanie.Szczegol), wywolanie.Zaszlo)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać wywołania pozycji %q: %w",
			wywolanie.RozszerzenieKod, err)
	}
	return nil
}

// MetrykaUzyciaRozszerzenia to policzony wynik jednego okna czasu: liczba
// wywołań, niepowodzeń i średni czas.
type MetrykaUzyciaRozszerzenia struct {
	RozszerzenieKod string
	Wywolan         int
	Niepowodzen     int
	SredniCzasMs    *int64
}

// MetrykiUzyciaRozszerzen liczy metryki użycia integracji w zadanym oknie
// czasu, zgrupowane po pozycji katalogu.
func (r *repozytoriumRozszerzen) MetrykiUzyciaRozszerzen(ctx context.Context, rozszerzenie string,
	od, do int64) ([]MetrykaUzyciaRozszerzenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, metrykaWywolanRozszerzenia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rozszerzenie, rozszerzenie, od, do)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można policzyć metryk użycia rozszerzeń: %w", err)
	}
	defer wiersze.Close()

	lista := []MetrykaUzyciaRozszerzenia{}
	for wiersze.Next() {
		var metryka MetrykaUzyciaRozszerzenia
		var niepowodzen sql.NullInt64
		var sredni sql.NullFloat64
		if err := wiersze.Scan(&metryka.RozszerzenieKod, &metryka.Wywolan,
			&niepowodzen, &sredni); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz metryki użycia: %w", err)
		}
		metryka.Niepowodzen = int(niepowodzen.Int64)
		if sredni.Valid {
			zaokraglony := int64(sredni.Float64 + 0.5)
			metryka.SredniCzasMs = &zaokraglony
		}
		lista = append(lista, metryka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt metryk użycia: %w", err)
	}
	return lista, nil
}

// AudytRozszerzen zwraca stronę wywołań widzianą od strony eksperta wraz
// z liczbą wszystkich spełniających te same warunki.
func (r *repozytoriumRozszerzen) AudytRozszerzen(ctx context.Context, rozszerzenie, agent string,
	od int64, granica int) ([]WywolanieRozszerzenia, int, error) {

	if granica <= 0 {
		granica = 200
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaAudytuRozszerzenia)
	if err != nil {
		return nil, 0, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rozszerzenie, rozszerzenie, agent, agent, od, granica)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać audytu rozszerzeń: %w", err)
	}
	defer wiersze.Close()

	lista := []WywolanieRozszerzenia{}
	for wiersze.Next() {
		var wpis WywolanieRozszerzenia
		var agentKod, narzedzie sql.NullString
		if err := wiersze.Scan(&wpis.RozszerzenieKod, &agentKod, &narzedzie, &wpis.Zaszlo); err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz audytu rozszerzeń: %w", err)
		}
		wpis.AgentKod = tekstZKolumny(agentKod)
		wpis.Narzedzie = tekstZKolumny(narzedzie)
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt audytu rozszerzeń: %w", err)
	}

	liczenie, err := r.zapytania.przygotuj(ctx, policzAudytRozszerzenia)
	if err != nil {
		return nil, 0, err
	}
	var razem int
	if err := liczenie.QueryRowContext(ctx, rozszerzenie, rozszerzenie, agent, agent, od).
		Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć wpisów audytu: %w", err)
	}
	return lista, razem, nil
}

// DopiszKondycjeRozszerzenia odnotowuje wynik jednego sprawdzenia kondycji
// integracji w dzienniku kondycji.
func (r *repozytoriumRozszerzen) DopiszKondycjeRozszerzenia(ctx context.Context,
	kondycja KondycjaRozszerzenia) error {

	if kondycja.RozszerzenieKod == "" || kondycja.Stan == "" {
		return fmt.Errorf("dane: wynik kondycji rozszerzenia bez pozycji albo stanu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKondycjeRozszerzenia)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, kondycja.RozszerzenieKod, kondycja.Stan,
		liczbaDoKolumny(kondycja.CzasMs), liczbaDoKolumny(kondycja.LiczbaNarzedzi),
		tekstDoKolumny(kondycja.BladPowitania), kondycja.Sprawdzono)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać kondycji pozycji %q: %w",
			kondycja.RozszerzenieKod, err)
	}
	return nil
}

// OstatniaKondycjaRozszerzenia zwraca najnowszy wynik sprawdzenia kondycji
// wskazanej pozycji katalogu.
func (r *repozytoriumRozszerzen) OstatniaKondycjaRozszerzenia(ctx context.Context,
	rozszerzenie string) (KondycjaRozszerzenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, ostatniaKondycjaRozszerzenia)
	if err != nil {
		return KondycjaRozszerzenia{}, err
	}
	var kondycja KondycjaRozszerzenia
	var czas, narzedzi sql.NullInt64
	var blad sql.NullString
	err = polecenie.QueryRowContext(ctx, rozszerzenie).Scan(&kondycja.ID, &kondycja.RozszerzenieKod,
		&kondycja.Stan, &czas, &narzedzi, &blad, &kondycja.Sprawdzono)
	if err == sql.ErrNoRows {
		return KondycjaRozszerzenia{}, ErrBrakWiersza
	}
	if err != nil {
		return KondycjaRozszerzenia{}, fmt.Errorf("dane: nieczytelna kondycja pozycji %q: %w",
			rozszerzenie, err)
	}
	kondycja.CzasMs = liczbaZKolumny(czas)
	kondycja.LiczbaNarzedzi = liczbaZKolumny(narzedzi)
	kondycja.BladPowitania = tekstZKolumny(blad)
	return kondycja, nil
}
