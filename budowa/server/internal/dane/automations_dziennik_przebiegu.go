// Plik prowadzi trwały zapis tego, co działo się w przebiegu: log przebiegu, stan i ładunki kroków oraz punkty
// wznowienia; wszystkie trzy istnieją z jednego powodu — mają przeżyć restart rdzenia.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WpisLoguPrzebiegu to wiersz tabeli `log_przebiegu` niosący jeden zapis dziennika przebiegu automatyki.
type WpisLoguPrzebiegu struct {
	PrzebiegID int64
	KrokKod    *string
	Poziom     string
	Tresc      string
	Chwila     string
}

// KrokPrzebiegu to wiersz tabeli `krok_przebiegu_automatyki` — stan jednego
// kroku jednego uruchomienia wraz z jego ładunkami.
type KrokPrzebiegu struct {
	PrzebiegID     int64
	KrokKod        string
	Stan           string
	Proba          int
	KomunikatBledu *string
	SladStosu      *string
	LadunekWejscia *string
	LadunekWyjscia *string
	Kolejnosc      int
	Rozpoczeto     string
	Zakonczono     *string
}

// PunktWznowienia to wiersz tabeli `punkt_wznowienia_przebiegu` niosący stan pozwalający wznowić przebieg.
type PunktWznowienia struct {
	Kod            string
	PrzebiegID     int64
	KrokKod        string
	KrokiUkonczone string
	Utworzono      string
}

const (
	dopiszLogPrzebiegu = `INSERT INTO log_przebiegu (przebieg_id, krok_kod, poziom, tresc)
	                      VALUES (?, ?, ?, ?)`

	// Zawężenia są opcjonalne i idą jednym zapytaniem: pusty poziom i pusty krok
	// znaczą „bez zawężenia”. Drugie zapytanie na każdą kombinację byłoby
	// czterema wariantami tego samego zdania.
	listaLoguPrzebiegu = `SELECT przebieg_id, krok_kod, poziom, tresc, chwila
	                      FROM log_przebiegu
	                      WHERE przebieg_id = ?
	                        AND (? = '' OR poziom = ?)
	                        AND (? = '' OR krok_kod = ?)
	                      ORDER BY id LIMIT ?`

	zapiszKrokPrzebiegu = `INSERT INTO krok_przebiegu_automatyki
	                       (przebieg_id, krok_kod, stan, proba, komunikat_bledu, slad_stosu,
	                        ladunek_wejscia, ladunek_wyjscia, kolejnosc, zakonczono)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(przebieg_id, krok_kod) DO UPDATE SET
	                           stan = excluded.stan,
	                           proba = excluded.proba,
	                           komunikat_bledu = excluded.komunikat_bledu,
	                           slad_stosu = excluded.slad_stosu,
	                           ladunek_wejscia = COALESCE(excluded.ladunek_wejscia,
	                                                      krok_przebiegu_automatyki.ladunek_wejscia),
	                           ladunek_wyjscia = COALESCE(excluded.ladunek_wyjscia,
	                                                      krok_przebiegu_automatyki.ladunek_wyjscia),
	                           kolejnosc = excluded.kolejnosc,
	                           zakonczono = excluded.zakonczono`

	kolumnyKrokuPrzebiegu = `przebieg_id, krok_kod, stan, proba, komunikat_bledu, slad_stosu,
	                         ladunek_wejscia, ladunek_wyjscia, kolejnosc, rozpoczeto, zakonczono`

	listaKrokowPrzebiegu = `SELECT ` + kolumnyKrokuPrzebiegu + ` FROM krok_przebiegu_automatyki
	                        WHERE przebieg_id = ? ORDER BY kolejnosc, id`

	pobierzKrokPrzebiegu = `SELECT ` + kolumnyKrokuPrzebiegu + ` FROM krok_przebiegu_automatyki
	                        WHERE przebieg_id = ? AND krok_kod = ?`

	zapiszPunktWznowienia = `INSERT INTO punkt_wznowienia_przebiegu
	                         (identyfikator_zewnetrzny, przebieg_id, krok_kod, kroki_ukonczone)
	                         VALUES (?, ?, ?, ?)
	                         ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                             krok_kod = excluded.krok_kod,
	                             kroki_ukonczone = excluded.kroki_ukonczone`

	listaPunktowWznowienia = `SELECT identyfikator_zewnetrzny, przebieg_id, krok_kod,
	                                 kroki_ukonczone, utworzono
	                          FROM punkt_wznowienia_przebiegu
	                          WHERE przebieg_id = ? ORDER BY id`

	pobierzPunktWznowienia = `SELECT identyfikator_zewnetrzny, przebieg_id, krok_kod,
	                                 kroki_ukonczone, utworzono
	                          FROM punkt_wznowienia_przebiegu
	                          WHERE identyfikator_zewnetrzny = ?`
)

// DopiszLogPrzebiegu nanosi jeden nowy wiersz dziennika przebiegu automatyki wraz z jego pełną treścią.
func (r *repozytoriumAutomatyk) DopiszLogPrzebiegu(ctx context.Context, wpis WpisLoguPrzebiegu) error {
	polecenie, err := r.zapytania.przygotuj(ctx, dopiszLogPrzebiegu)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wpis.PrzebiegID, tekstDoKolumny(wpis.KrokKod),
		wpis.Poziom, wpis.Tresc)
	if err != nil {
		return fmt.Errorf("dane: nie można dopisać wiersza logu przebiegu %d: %w", wpis.PrzebiegID, err)
	}
	return nil
}

// LogPrzebiegu zwraca wiersze dziennika w kolejności czasu, zawężone poziomem
// i krokiem. Puste zawężenie znaczy „bez zawężenia”.
func (r *repozytoriumAutomatyk) LogPrzebiegu(ctx context.Context, przebiegID int64,
	poziom, krokKod string, limit int) ([]WpisLoguPrzebiegu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaLoguPrzebiegu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, przebiegID, poziom, poziom,
		krokKod, krokKod, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać logu przebiegu %d: %w", przebiegID, err)
	}
	defer wiersze.Close()

	lista := []WpisLoguPrzebiegu{}
	for wiersze.Next() {
		var wpis WpisLoguPrzebiegu
		var krok sql.NullString
		if err := wiersze.Scan(&wpis.PrzebiegID, &krok, &wpis.Poziom, &wpis.Tresc,
			&wpis.Chwila); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz logu przebiegu: %w", err)
		}
		wpis.KrokKod = tekstZKolumny(krok)
		lista = append(lista, wpis)
	}
	return lista, wiersze.Err()
}

// ZapiszKrokPrzebiegu zapisuje stan kroku przebiegu. Ładunek pusty zostawia
// zapisany: ponowna zmiana samego stanu nie ma kasować danych wejściowych.
func (r *repozytoriumAutomatyk) ZapiszKrokPrzebiegu(ctx context.Context, krok KrokPrzebiegu) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKrokPrzebiegu)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, krok.PrzebiegID, krok.KrokKod, krok.Stan, krok.Proba,
		tekstDoKolumny(krok.KomunikatBledu), tekstDoKolumny(krok.SladStosu),
		tekstDoKolumny(krok.LadunekWejscia), tekstDoKolumny(krok.LadunekWyjscia),
		krok.Kolejnosc, tekstDoKolumny(krok.Zakonczono))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać kroku %q przebiegu %d: %w",
			krok.KrokKod, krok.PrzebiegID, err)
	}
	return nil
}

// KrokiPrzebiegu zwraca bieżący stan każdego kroku danego przebiegu automatyki, w kolejności wykonania.
func (r *repozytoriumAutomatyk) KrokiPrzebiegu(ctx context.Context,
	przebiegID int64) ([]KrokPrzebiegu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaKrokowPrzebiegu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, przebiegID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kroków przebiegu %d: %w", przebiegID, err)
	}
	defer wiersze.Close()

	lista := []KrokPrzebiegu{}
	for wiersze.Next() {
		krok, err := odczytajKrokPrzebiegu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kroku przebiegu: %w", err)
		}
		lista = append(lista, krok)
	}
	return lista, wiersze.Err()
}

// KrokPrzebieguPoKodzie zwraca jeden krok przebiegu — źródło podglądu jego ładunku wejścia i wyjścia z bazy.
func (r *repozytoriumAutomatyk) KrokPrzebieguPoKodzie(ctx context.Context, przebiegID int64,
	krokKod string) (KrokPrzebiegu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKrokPrzebiegu)
	if err != nil {
		return KrokPrzebiegu{}, err
	}
	krok, err := odczytajKrokPrzebiegu(polecenie.QueryRowContext(ctx, przebiegID, krokKod))
	if errors.Is(err, sql.ErrNoRows) {
		return KrokPrzebiegu{}, ErrBrakWiersza
	}
	if err != nil {
		return KrokPrzebiegu{}, fmt.Errorf("dane: nieczytelny krok %q przebiegu %d: %w",
			krokKod, przebiegID, err)
	}
	return krok, nil
}

// ZapiszPunktWznowienia zapisuje punkt wznowienia przebiegu automatyki wraz z jego pełnym stanem wykonania kroku.
func (r *repozytoriumAutomatyk) ZapiszPunktWznowienia(ctx context.Context,
	punkt PunktWznowienia) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPunktWznowienia)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, punkt.Kod, punkt.PrzebiegID, punkt.KrokKod,
		punkt.KrokiUkonczone)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać punktu wznowienia przebiegu %d: %w",
			punkt.PrzebiegID, err)
	}
	return nil
}

// PunktyWznowienia zwraca wszystkie punkty wznowienia przebiegu automatyki, w kolejności czasu powstania.
func (r *repozytoriumAutomatyk) PunktyWznowienia(ctx context.Context,
	przebiegID int64) ([]PunktWznowienia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPunktowWznowienia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, przebiegID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać punktów wznowienia przebiegu %d: %w",
			przebiegID, err)
	}
	defer wiersze.Close()

	lista := []PunktWznowienia{}
	for wiersze.Next() {
		punkt, err := odczytajPunktWznowienia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz punktu wznowienia: %w", err)
		}
		lista = append(lista, punkt)
	}
	return lista, wiersze.Err()
}

// PunktWznowieniaPoKodzie zwraca jeden punkt wznowienia przebiegu po jego kodzie zewnętrznym z bazy danych.
func (r *repozytoriumAutomatyk) PunktWznowieniaPoKodzie(ctx context.Context,
	kod string) (PunktWznowienia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPunktWznowienia)
	if err != nil {
		return PunktWznowienia{}, err
	}
	punkt, err := odczytajPunktWznowienia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return PunktWznowienia{}, ErrBrakWiersza
	}
	if err != nil {
		return PunktWznowienia{}, fmt.Errorf("dane: nieczytelny punkt wznowienia %q: %w", kod, err)
	}
	return punkt, nil
}

// odczytajKrokPrzebiegu składa krok przebiegu automatyki wprost z jednego wiersza wyniku zapytania SQL.
func odczytajKrokPrzebiegu(wiersz skaner) (KrokPrzebiegu, error) {
	var krok KrokPrzebiegu
	var blad, slad, wejscie, wyjscie, zakonczono sql.NullString
	err := wiersz.Scan(&krok.PrzebiegID, &krok.KrokKod, &krok.Stan, &krok.Proba,
		&blad, &slad, &wejscie, &wyjscie, &krok.Kolejnosc, &krok.Rozpoczeto, &zakonczono)
	if err != nil {
		return KrokPrzebiegu{}, err
	}
	krok.KomunikatBledu = tekstZKolumny(blad)
	krok.SladStosu = tekstZKolumny(slad)
	krok.LadunekWejscia = tekstZKolumny(wejscie)
	krok.LadunekWyjscia = tekstZKolumny(wyjscie)
	krok.Zakonczono = tekstZKolumny(zakonczono)
	return krok, nil
}

// odczytajPunktWznowienia składa punkt wznowienia przebiegu wprost z jednego wiersza wyniku zapytania.
func odczytajPunktWznowienia(wiersz skaner) (PunktWznowienia, error) {
	var punkt PunktWznowienia
	err := wiersz.Scan(&punkt.Kod, &punkt.PrzebiegID, &punkt.KrokKod,
		&punkt.KrokiUkonczone, &punkt.Utworzono)
	return punkt, err
}
