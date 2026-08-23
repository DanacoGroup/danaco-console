// Odpowiedzialność pliku: rodzina `extension.*` — warstwa zaufania.
// Uprawnienia deklarowane i nadane (`uprawnienie_rozszerzenia`), podpis pozycji
// (`podpis_rozszerzenia`) oraz rejestr referencji sekretów wraz z zakresem
// współdzielenia (`sekret_rozszerzenia`,
// `udostepnienie_sekretu_rozszerzenia`) —
// `store/migracja_210_rozszerzenia_zaufanie.sql`.
//
// Żaden wiersz tej warstwy niczego nie blokuje. Uprawnienie jest zapisem tego,
// co manifest deklaruje i co Operator nadał; podpis — zapisem wyniku
// weryfikacji; referencja sekretu — kluczem jawnym, nigdy treścią poświadczenia.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// UprawnienieRozszerzenia to wiersz tabeli `uprawnienie_rozszerzenia`.
// Kolumna `Nadane` rozstrzyga, czy wiersz opisuje deklarację manifestu, czy
// nadanie Operatora — patrz czoło migracji 210.
type UprawnienieRozszerzenia struct {
	ID              int64
	RozszerzenieKod string
	Zakres          string
	Byt             *string
	Tryb            *string
	Objasnienie     *string
	Nadane          bool
	AgentKod        *string
	Nadano          *int64
}

// PodpisRozszerzenia to wiersz tabeli `podpis_rozszerzenia`.
type PodpisRozszerzenia struct {
	RozszerzenieKod string
	Algorytm        *string
	SumaKontrolna   *string
	Wydawca         *string
	PodpisBase64    *string
	KluczBase64     *string
	Zaktualizowano  int64
}

// SekretRozszerzenia to wiersz tabeli `sekret_rozszerzenia` wraz z zakresem
// współdzielenia — kontrakt oddaje referencję zawsze razem z nim.
type SekretRozszerzenia struct {
	ID              int64
	Odwolanie       string
	Etykieta        *string
	SposobLogowania *string
	Wygasa          *int64
	KodyRozszerzen  []string
	KodyRol         []string
	Zaktualizowano  int64
}

const (
	usunUprawnieniaRozszerzenia = `DELETE FROM uprawnienie_rozszerzenia
	                               WHERE rozszerzenie_kod = ? AND nadane = ?`

	wstawUprawnienieRozszerzenia = `INSERT INTO uprawnienie_rozszerzenia
	                                (rozszerzenie_kod, zakres, byt, tryb, objasnienie,
	                                 nadane, agent_kod, nadano)
	                                VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	                                ON CONFLICT(rozszerzenie_kod, zakres, byt, nadane, agent_kod)
	                                DO UPDATE SET
	                                    tryb = excluded.tryb,
	                                    objasnienie = excluded.objasnienie,
	                                    nadano = excluded.nadano`

	listaUprawnienRozszerzenia = `SELECT id, rozszerzenie_kod, zakres, byt, tryb, objasnienie,
	                                     nadane, agent_kod, nadano
	                              FROM uprawnienie_rozszerzenia
	                              WHERE rozszerzenie_kod = ? ORDER BY nadane, zakres, id`

	zapiszPodpisRozszerzenia = `INSERT INTO podpis_rozszerzenia
	                            (rozszerzenie_kod, algorytm, suma_kontrolna, wydawca,
	                             podpis_base64, klucz_base64, zaktualizowano)
	                            VALUES (?, ?, ?, ?, ?, ?, ?)
	                            ON CONFLICT(rozszerzenie_kod) DO UPDATE SET
	                                algorytm = excluded.algorytm,
	                                suma_kontrolna = excluded.suma_kontrolna,
	                                wydawca = excluded.wydawca,
	                                podpis_base64 = excluded.podpis_base64,
	                                klucz_base64 = excluded.klucz_base64,
	                                zaktualizowano = excluded.zaktualizowano`

	pobierzPodpisRozszerzenia = `SELECT rozszerzenie_kod, algorytm, suma_kontrolna, wydawca,
	                                    podpis_base64, klucz_base64, zaktualizowano
	                             FROM podpis_rozszerzenia WHERE rozszerzenie_kod = ?`

	zapiszSekretRozszerzenia = `INSERT INTO sekret_rozszerzenia
	                            (odwolanie, etykieta, sposob_logowania, wygasa, zaktualizowano)
	                            VALUES (?, ?, ?, ?, ?)
	                            ON CONFLICT(odwolanie) DO UPDATE SET
	                                etykieta = IFNULL(excluded.etykieta, sekret_rozszerzenia.etykieta),
	                                sposob_logowania = IFNULL(excluded.sposob_logowania,
	                                                          sekret_rozszerzenia.sposob_logowania),
	                                wygasa = IFNULL(excluded.wygasa, sekret_rozszerzenia.wygasa),
	                                zaktualizowano = excluded.zaktualizowano`

	pobierzSekretRozszerzenia = `SELECT id, odwolanie, etykieta, sposob_logowania, wygasa,
	                                    zaktualizowano
	                             FROM sekret_rozszerzenia WHERE odwolanie = ?`

	// Zawężenie po terminie ważności: `expiringWithinDays` kontraktu przekłada
	// się na górną granicę czasu, a wartość zerowa granicy wyłącza warunek.
	listaSekretowRozszerzenia = `SELECT id, odwolanie, etykieta, sposob_logowania, wygasa,
	                                    zaktualizowano
	                             FROM sekret_rozszerzenia
	                             WHERE (? = 0 OR (wygasa IS NOT NULL AND wygasa <= ?))
	                             ORDER BY odwolanie`

	usunUdostepnieniaSekretu = `DELETE FROM udostepnienie_sekretu_rozszerzenia WHERE sekret_id = ?`

	wstawUdostepnienieSekretu = `INSERT INTO udostepnienie_sekretu_rozszerzenia
	                             (sekret_id, rodzaj, byt_kod) VALUES (?, ?, ?)
	                             ON CONFLICT(sekret_id, rodzaj, byt_kod) DO NOTHING`

	listaUdostepnienSekretu = `SELECT s.odwolanie, u.rodzaj, u.byt_kod
	                           FROM udostepnienie_sekretu_rozszerzenia u
	                           JOIN sekret_rozszerzenia s ON s.id = u.sekret_id
	                           ORDER BY u.byt_kod`
)

// ZapiszUprawnieniaRozszerzenia wymienia komplet uprawnień jednej strony —
// deklarowanych albo nadanych. Wymiana, nie dokładanie: manifest, który
// przestał deklarować uprawnienie, ma przestać je pokazywać, a nadanie zdjęte
// przez Operatora ma zniknąć, nie zostać.
func (r *repozytoriumRozszerzen) ZapiszUprawnieniaRozszerzenia(ctx context.Context,
	rozszerzenie string, nadane bool, uprawnienia []UprawnienieRozszerzenia) error {

	if rozszerzenie == "" {
		return fmt.Errorf("dane: uprawnienia bez pozycji katalogu")
	}
	return wTransakcji(ctx, r.baza, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunUprawnieniaRozszerzenia)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, rozszerzenie, liczbaLogiczna(nadane)); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić uprawnień pozycji %q: %w", rozszerzenie, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawUprawnienieRozszerzenia)
		if err != nil {
			return err
		}
		for _, uprawnienie := range uprawnienia {
			_, err := wstawienie.ExecContext(ctx, rozszerzenie, uprawnienie.Zakres,
				tekstDoKolumny(uprawnienie.Byt), tekstDoKolumny(uprawnienie.Tryb),
				tekstDoKolumny(uprawnienie.Objasnienie), liczbaLogiczna(nadane),
				tekstDoKolumny(uprawnienie.AgentKod), liczbaDoKolumny(uprawnienie.Nadano))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać uprawnienia %q pozycji %q: %w",
					uprawnienie.Zakres, rozszerzenie, err)
			}
		}
		return nil
	})
}

// UprawnieniaRozszerzenia zwraca uprawnienia pozycji — obie strony naraz;
// warstwa wyższa rozdziela je po kolumnie `Nadane`.
func (r *repozytoriumRozszerzen) UprawnieniaRozszerzenia(ctx context.Context,
	rozszerzenie string) ([]UprawnienieRozszerzenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaUprawnienRozszerzenia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rozszerzenie)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać uprawnień pozycji %q: %w", rozszerzenie, err)
	}
	defer wiersze.Close()

	lista := []UprawnienieRozszerzenia{}
	for wiersze.Next() {
		var uprawnienie UprawnienieRozszerzenia
		var byt, tryb, objasnienie, agent sql.NullString
		var nadane int
		var nadano sql.NullInt64
		err := wiersze.Scan(&uprawnienie.ID, &uprawnienie.RozszerzenieKod, &uprawnienie.Zakres,
			&byt, &tryb, &objasnienie, &nadane, &agent, &nadano)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz uprawnienia rozszerzenia: %w", err)
		}
		uprawnienie.Byt = tekstZKolumny(byt)
		uprawnienie.Tryb = tekstZKolumny(tryb)
		uprawnienie.Objasnienie = tekstZKolumny(objasnienie)
		uprawnienie.Nadane = nadane == 1
		uprawnienie.AgentKod = tekstZKolumny(agent)
		uprawnienie.Nadano = liczbaZKolumny(nadano)
		lista = append(lista, uprawnienie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt uprawnień pozycji %q: %w", rozszerzenie, err)
	}
	return lista, nil
}

// ZapiszPodpisRozszerzenia utrwala podpis pozycji wraz z materiałem do jego
// ponownej weryfikacji.
func (r *repozytoriumRozszerzen) ZapiszPodpisRozszerzenia(ctx context.Context,
	podpis PodpisRozszerzenia) error {

	if podpis.RozszerzenieKod == "" {
		return fmt.Errorf("dane: podpis bez pozycji katalogu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPodpisRozszerzenia)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, podpis.RozszerzenieKod, tekstDoKolumny(podpis.Algorytm),
		tekstDoKolumny(podpis.SumaKontrolna), tekstDoKolumny(podpis.Wydawca),
		tekstDoKolumny(podpis.PodpisBase64), tekstDoKolumny(podpis.KluczBase64),
		podpis.Zaktualizowano)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać podpisu pozycji %q: %w",
			podpis.RozszerzenieKod, err)
	}
	return nil
}

// PodpisRozszerzenia zwraca podpis pozycji; brak wraca jako ErrBrakWiersza.
func (r *repozytoriumRozszerzen) PodpisRozszerzenia(ctx context.Context,
	rozszerzenie string) (PodpisRozszerzenia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPodpisRozszerzenia)
	if err != nil {
		return PodpisRozszerzenia{}, err
	}
	var podpis PodpisRozszerzenia
	var algorytm, suma, wydawca, bajty, klucz sql.NullString
	err = polecenie.QueryRowContext(ctx, rozszerzenie).Scan(&podpis.RozszerzenieKod,
		&algorytm, &suma, &wydawca, &bajty, &klucz, &podpis.Zaktualizowano)
	if err == sql.ErrNoRows {
		return PodpisRozszerzenia{}, ErrBrakWiersza
	}
	if err != nil {
		return PodpisRozszerzenia{}, fmt.Errorf("dane: nieczytelny podpis pozycji %q: %w",
			rozszerzenie, err)
	}
	podpis.Algorytm = tekstZKolumny(algorytm)
	podpis.SumaKontrolna = tekstZKolumny(suma)
	podpis.Wydawca = tekstZKolumny(wydawca)
	podpis.PodpisBase64 = tekstZKolumny(bajty)
	podpis.KluczBase64 = tekstZKolumny(klucz)
	return podpis, nil
}

// ZapiszSekretRozszerzenia zapisuje referencję sekretu wraz z zakresem
// współdzielenia; zakres wymieniany jest w całości, bo kontrakt nadsyła oba
// wykazy w komplecie.
func (r *repozytoriumRozszerzen) ZapiszSekretRozszerzenia(ctx context.Context,
	sekret SekretRozszerzenia, wymienZakres bool) (SekretRozszerzenia, error) {

	if sekret.Odwolanie == "" {
		return SekretRozszerzenia{}, fmt.Errorf("dane: referencja sekretu bez klucza jawnego")
	}
	err := wTransakcji(ctx, r.baza, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszSekretRozszerzenia)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, sekret.Odwolanie, tekstDoKolumny(sekret.Etykieta),
			tekstDoKolumny(sekret.SposobLogowania), liczbaDoKolumny(sekret.Wygasa),
			sekret.Zaktualizowano); err != nil {
			return fmt.Errorf("dane: nie można zapisać referencji %q: %w", sekret.Odwolanie, err)
		}
		if !wymienZakres {
			return nil
		}

		var sekretID int64
		wiersz := transakcja.QueryRowContext(ctx,
			`SELECT id FROM sekret_rozszerzenia WHERE odwolanie = ?`, sekret.Odwolanie)
		if err := wiersz.Scan(&sekretID); err != nil {
			return fmt.Errorf("dane: nie można odczytać id referencji %q: %w", sekret.Odwolanie, err)
		}
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunUdostepnieniaSekretu)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, sekretID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić zakresu referencji %q: %w",
				sekret.Odwolanie, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawUdostepnienieSekretu)
		if err != nil {
			return err
		}
		for _, kod := range sekret.KodyRozszerzen {
			if kod == "" {
				continue
			}
			if _, err := wstawienie.ExecContext(ctx, sekretID, "rozszerzenie", kod); err != nil {
				return fmt.Errorf("dane: nie można zapisać zakresu referencji %q: %w",
					sekret.Odwolanie, err)
			}
		}
		for _, kod := range sekret.KodyRol {
			if kod == "" {
				continue
			}
			if _, err := wstawienie.ExecContext(ctx, sekretID, "rola", kod); err != nil {
				return fmt.Errorf("dane: nie można zapisać zakresu referencji %q: %w",
					sekret.Odwolanie, err)
			}
		}
		return nil
	})
	if err != nil {
		return SekretRozszerzenia{}, err
	}
	return r.SekretRozszerzenia(ctx, sekret.Odwolanie)
}

// SekretRozszerzenia zwraca jedną referencję wraz z zakresem współdzielenia.
func (r *repozytoriumRozszerzen) SekretRozszerzenia(ctx context.Context, odwolanie string) (SekretRozszerzenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSekretRozszerzenia)
	if err != nil {
		return SekretRozszerzenia{}, err
	}
	sekret, err := odczytajSekretRozszerzenia(polecenie.QueryRowContext(ctx, odwolanie))
	if err == sql.ErrNoRows {
		return SekretRozszerzenia{}, ErrBrakWiersza
	}
	if err != nil {
		return SekretRozszerzenia{}, fmt.Errorf("dane: nieczytelna referencja %q: %w", odwolanie, err)
	}
	rozszerzenia, role, err := r.zakresySekretowRozszerzen(ctx)
	if err != nil {
		return SekretRozszerzenia{}, err
	}
	sekret.KodyRozszerzen = rozszerzenia[sekret.Odwolanie]
	sekret.KodyRol = role[sekret.Odwolanie]
	return sekret, nil
}

// SekretyRozszerzen zwraca referencje sekretów; granica czasu zeruje warunek
// wygaśnięcia, gdy wynosi zero.
func (r *repozytoriumRozszerzen) SekretyRozszerzen(ctx context.Context, doCzasu int64) ([]SekretRozszerzenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSekretowRozszerzenia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, doCzasu, doCzasu)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać referencji sekretów: %w", err)
	}
	defer wiersze.Close()

	lista := []SekretRozszerzenia{}
	for wiersze.Next() {
		sekret, err := odczytajSekretRozszerzenia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz referencji sekretu: %w", err)
		}
		lista = append(lista, sekret)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt referencji sekretów: %w", err)
	}
	rozszerzenia, role, err := r.zakresySekretowRozszerzen(ctx)
	if err != nil {
		return nil, err
	}
	for indeks := range lista {
		lista[indeks].KodyRozszerzen = rozszerzenia[lista[indeks].Odwolanie]
		lista[indeks].KodyRol = role[lista[indeks].Odwolanie]
	}
	return lista, nil
}

// zakresySekretowRozszerzen zwraca dwie mapy: odwołanie → pozycje i odwołanie →
// role. Jedno zapytanie na cały rejestr — patrz wzorzec kolekcji.
func (r *repozytoriumRozszerzen) zakresySekretowRozszerzen(ctx context.Context) (
	map[string][]string, map[string][]string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaUdostepnienSekretu)
	if err != nil {
		return nil, nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("dane: nie można odczytać zakresu referencji: %w", err)
	}
	defer wiersze.Close()

	rozszerzenia := map[string][]string{}
	role := map[string][]string{}
	for wiersze.Next() {
		var odwolanie, rodzaj, kod string
		if err := wiersze.Scan(&odwolanie, &rodzaj, &kod); err != nil {
			return nil, nil, fmt.Errorf("dane: nieczytelny wiersz zakresu referencji: %w", err)
		}
		if rodzaj == "rola" {
			role[odwolanie] = append(role[odwolanie], kod)
			continue
		}
		rozszerzenia[odwolanie] = append(rozszerzenia[odwolanie], kod)
	}
	if err := wiersze.Err(); err != nil {
		return nil, nil, fmt.Errorf("dane: przerwany odczyt zakresu referencji: %w", err)
	}
	return rozszerzenia, role, nil
}

// odczytajSekretRozszerzenia składa referencję z jednego wiersza wyniku; zakres
// współdzielenia dokłada wołający z osobnego zapytania.
func odczytajSekretRozszerzenia(wiersz skaner) (SekretRozszerzenia, error) {
	var sekret SekretRozszerzenia
	var etykieta, sposob sql.NullString
	var wygasa sql.NullInt64
	err := wiersz.Scan(&sekret.ID, &sekret.Odwolanie, &etykieta, &sposob, &wygasa,
		&sekret.Zaktualizowano)
	if err != nil {
		return SekretRozszerzenia{}, err
	}
	sekret.Etykieta = tekstZKolumny(etykieta)
	sekret.SposobLogowania = tekstZKolumny(sposob)
	sekret.Wygasa = liczbaZKolumny(wygasa)
	return sekret, nil
}
