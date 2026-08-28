// Plik prowadzi to, co budzik harmonogramu musi zastać po restarcie rdzenia: okna wykonania, nadzór obecności
// uruchomień, klucz podpisu webhooka oraz historię rzeczywistych wyzwoleń; pamięć procesu nie przeżywa restartu, więc wszystko idzie do bazy.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// OknoWykonania to wiersz tabeli `okno_wykonania_harmonogramu` — przedział
// czasu, w którym uruchomienie w ogóle następuje.
type OknoWykonania struct {
	DniTygodnia   string
	MinutaOd      int
	MinutaDo      int
	StrefaCzasowa *string
	Kolejnosc     int
}

// WyzwolenieAutomatyki to wiersz tabeli `wyzwolenie_automatyki` — moment
// RZECZYWISTY wraz z przyczyną.
type WyzwolenieAutomatyki struct {
	AutomatykaID   int64
	AutomatykaKod  string
	HarmonogramID  *int64
	HarmonogramKod *string
	Przyczyna      string
	PrzebiegID     *int64
	PrzebiegKod    *string
	Chwila         string
}

const (
	usunOknaWykonania = `DELETE FROM okno_wykonania_harmonogramu WHERE harmonogram_id = ?`

	wstawOknoWykonania = `INSERT INTO okno_wykonania_harmonogramu
	                      (harmonogram_id, dni_tygodnia, minuta_od, minuta_do,
	                       strefa_czasowa, kolejnosc)
	                      VALUES (?, ?, ?, ?, ?, ?)`

	listaOkienWykonania = `SELECT dni_tygodnia, minuta_od, minuta_do, strefa_czasowa, kolejnosc
	                       FROM okno_wykonania_harmonogramu WHERE harmonogram_id = ?
	                       ORDER BY kolejnosc, id`

	ustawNadzorHarmonogramu = `UPDATE harmonogram_automatyki
	                           SET tolerancja_sekundy = ?, regula_nadzoru = ?,
	                               zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                           WHERE id = ?`

	ustawPodpisHarmonogramu = `UPDATE harmonogram_automatyki
	                           SET odwolanie_podpisu = ?,
	                               zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                           WHERE id = ?`

	dopiszWyzwolenie = `INSERT INTO wyzwolenie_automatyki
	                    (automatyka_id, harmonogram_id, przyczyna, przebieg_id)
	                    VALUES (?, ?, ?, ?)`

	// Zero w miejscu automatyki albo harmonogramu znaczy „bez zawężenia” — tak
	// mówi kontrakt o żądaniu bez pola `workflowId` czy `scheduleId`.
	listaWyzwolen = `SELECT w.automatyka_id, a.identyfikator_zewnetrzny, w.harmonogram_id,
	                        h.identyfikator_zewnetrzny, w.przyczyna, w.przebieg_id,
	                        p.identyfikator_zewnetrzny, w.chwila
	                 FROM wyzwolenie_automatyki w
	                 JOIN automatyka a ON a.id = w.automatyka_id
	                 LEFT JOIN harmonogram_automatyki h ON h.id = w.harmonogram_id
	                 LEFT JOIN przebieg_automatyki p ON p.id = w.przebieg_id
	                 WHERE (? = 0 OR w.automatyka_id = ?)
	                   AND (? = 0 OR w.harmonogram_id = ?)
	                 ORDER BY w.chwila DESC, w.id DESC LIMIT ?`
)

// ZapiszOknaWykonania podmienia komplet okien wykonania harmonogramu. Wykaz
// pusty znaczy „bez ograniczenia” — tak mówi kontrakt.
func (r *repozytoriumAutomatyk) ZapiszOknaWykonania(ctx context.Context, harmonogramID int64,
	okna []OknoWykonania) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja,
			usunOknaWykonania, harmonogramID); err != nil {
			return err
		}
		for numer, okno := range okna {
			err := wykonajWTransakcjiAutomatyzacji(ctx, r.zapytania, transakcja, wstawOknoWykonania,
				harmonogramID, okno.DniTygodnia, okno.MinutaOd, okno.MinutaDo,
				tekstDoKolumny(okno.StrefaCzasowa), numer+1)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// OknaWykonania zwraca okna wykonania harmonogramu w kolejności ich zapisu prosto z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) OknaWykonania(ctx context.Context,
	harmonogramID int64) ([]OknoWykonania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaOkienWykonania)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, harmonogramID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać okien wykonania harmonogramu %d: %w",
			harmonogramID, err)
	}
	defer wiersze.Close()

	lista := []OknoWykonania{}
	for wiersze.Next() {
		var okno OknoWykonania
		var strefa sql.NullString
		if err := wiersze.Scan(&okno.DniTygodnia, &okno.MinutaOd, &okno.MinutaDo,
			&strefa, &okno.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz okna wykonania: %w", err)
		}
		okno.StrefaCzasowa = tekstZKolumny(strefa)
		lista = append(lista, okno)
	}
	return lista, wiersze.Err()
}

// UstawNadzorHarmonogramu zapisuje okno tolerancji i regułę alarmowania
// nadzoru. Zero wyłącza nadzór — tak mówi kontrakt.
func (r *repozytoriumAutomatyk) UstawNadzorHarmonogramu(ctx context.Context, harmonogramID int64,
	tolerancjaSekundy int, regula *string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawNadzorHarmonogramu)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, tolerancjaSekundy, tekstDoKolumny(regula), harmonogramID)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać nadzoru harmonogramu %d: %w", harmonogramID, err)
	}
	return nil
}

// UstawPodpisHarmonogramu zapisuje ODWOŁANIE klucza podpisu webhooka. Wartość
// klucza leży w sejfie poza bazą i tu nie dociera.
func (r *repozytoriumAutomatyk) UstawPodpisHarmonogramu(ctx context.Context, harmonogramID int64,
	odwolanie string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawPodpisHarmonogramu)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, odwolanie, harmonogramID); err != nil {
		return fmt.Errorf("dane: nie można zapisać odwołania podpisu harmonogramu %d: %w",
			harmonogramID, err)
	}
	return nil
}

// DopiszWyzwolenie nanosi rzeczywisty moment wyzwolenia harmonogramu wraz z jego przyczyną w bazie danych.
func (r *repozytoriumAutomatyk) DopiszWyzwolenie(ctx context.Context, automatykaID int64,
	harmonogramID *int64, przyczyna string, przebiegID *int64) error {

	polecenie, err := r.zapytania.przygotuj(ctx, dopiszWyzwolenie)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, automatykaID, liczbaDoKolumny(harmonogramID),
		przyczyna, liczbaDoKolumny(przebiegID))
	if err != nil {
		return fmt.Errorf("dane: nie można dopisać wyzwolenia automatyki %d: %w", automatykaID, err)
	}
	return nil
}

// Wyzwolenia zwraca całą historię wyzwoleń harmonogramu od najnowszego wprost z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) Wyzwolenia(ctx context.Context, automatykaID, harmonogramID int64,
	limit int) ([]WyzwolenieAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWyzwolen)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID, automatykaID,
		harmonogramID, harmonogramID, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać historii wyzwoleń: %w", err)
	}
	defer wiersze.Close()

	lista := []WyzwolenieAutomatyki{}
	for wiersze.Next() {
		var wyzwolenie WyzwolenieAutomatyki
		var harmonogram, przebieg sql.NullInt64
		var kodHarmonogramu, kodPrzebiegu sql.NullString
		err := wiersze.Scan(&wyzwolenie.AutomatykaID, &wyzwolenie.AutomatykaKod, &harmonogram,
			&kodHarmonogramu, &wyzwolenie.Przyczyna, &przebieg, &kodPrzebiegu, &wyzwolenie.Chwila)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wyzwolenia: %w", err)
		}
		wyzwolenie.HarmonogramID = liczbaZKolumny(harmonogram)
		wyzwolenie.HarmonogramKod = tekstZKolumny(kodHarmonogramu)
		wyzwolenie.PrzebiegID = liczbaZKolumny(przebieg)
		wyzwolenie.PrzebiegKod = tekstZKolumny(kodPrzebiegu)
		lista = append(lista, wyzwolenie)
	}
	return lista, wiersze.Err()
}
