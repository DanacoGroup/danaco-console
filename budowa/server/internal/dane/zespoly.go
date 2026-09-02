// Trwałość zespołów ekspertów modułu Agents (tabele `zespol` i `zespol_sklad`) wraz z zapisem pełnego składu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Zespol to wiersz `zespol` wraz ze składem; `Kod` jest identyfikatorem trwałym, polem `Team.id` kontraktu.
type Zespol struct {
	ID             int64
	Kod            string
	Nazwa          string
	Opis           string
	Sklad          []string
	Pominieci      []string
	Utworzono      int64
	Zaktualizowano *int64
}

// FiltrZespolow zawęża wykaz zespołów zwracany zapytaniem `Lista` modułu Agents; pole puste znaczy „bez zawężenia”.
type FiltrZespolow struct {
	Fraza        string
	Granica      int
	Przesuniecie int
}

// RepozytoriumZespolow jest kontraktem trwałości zespołów ekspertów; `team.duplicate` składa warstwa wyższa.
type RepozytoriumZespolow interface {
	Lista(ctx context.Context, filtr FiltrZespolow) ([]Zespol, int, error)
	PoKodzie(ctx context.Context, kod string) (Zespol, error)
	Dodaj(ctx context.Context, zespol Zespol) (Zespol, error)
	Zapisz(ctx context.Context, zespol Zespol) (Zespol, error)
}

const (
	kolumnyZespolu = `id, kod, nazwa, opis, utworzono, zaktualizowano`

	// Fraza wchodzi do LIKE jako treść, nie wzorzec — stąd osłona i klauzula ESCAPE.
	listaZespolow = `SELECT ` + kolumnyZespolu + ` FROM zespol
	                 WHERE (? = '' OR lower(nazwa) LIKE ? ESCAPE '\' OR lower(opis) LIKE ? ESCAPE '\')
	                       AND ` + WarunekKonta + `
	                 ORDER BY nazwa, kod`

	zespolPoKodzie = `SELECT ` + kolumnyZespolu + ` FROM zespol WHERE kod = ? AND ` + WarunekKonta

	// Lewe złączenie z biblioteką: kod bez czynnego eksperta wraca jako pominięty, a nie znika.
	skladZespolow = `SELECT s.zespol_id, s.agent_kod,
	                        CASE WHEN a.id IS NULL THEN 0 ELSE 1 END
	                 FROM zespol_sklad s
	                 LEFT JOIN agent a
	                        ON a.kod = s.agent_kod AND a.zarchiwizowano_o IS NULL
	                 ORDER BY s.zespol_id, s.kolejnosc`

	wstawZespol = `INSERT INTO zespol (kod, nazwa, opis, utworzono, konto_id) VALUES (?, ?, ?, ?, ` + WskazanieKonta + `)`

	aktualizujZespol = `UPDATE zespol SET nazwa = ?, opis = ?, zaktualizowano = ? WHERE kod = ? AND ` + WarunekKonta

	usunSkladZespolu = `DELETE FROM zespol_sklad WHERE zespol_id = ?`

	wstawSkladZespolu = `INSERT INTO zespol_sklad (zespol_id, agent_kod, kolejnosc) VALUES (?, ?, ?)`

	numerZespolu = `SELECT id FROM zespol WHERE kod = ? AND ` + WarunekKonta
)

type repozytoriumZespolow struct {
	zapytania *zapytania
	db        *sql.DB
}

var _ RepozytoriumZespolow = (*repozytoriumZespolow)(nil)

func noweRepozytoriumZespolow(z *zapytania, db *sql.DB) *repozytoriumZespolow {
	return &repozytoriumZespolow{zapytania: z, db: db}
}

// Lista zwraca zespoły wraz z ich liczbą przed ucięciem granicą i przesunięciem; wykaz pusty nie jest błędem.
func (r *repozytoriumZespolow) Lista(ctx context.Context, filtr FiltrZespolow) ([]Zespol, int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaZespolow)
	if err != nil {
		return nil, 0, err
	}
	fraza := strings.ToLower(strings.TrimSpace(filtr.Fraza))
	wzorzec := "%" + oslonWieloznaczniki(fraza) + "%"
	wiersze, err := polecenie.QueryContext(ctx, fraza, wzorzec, wzorzec, KontoOperatora(ctx))
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać wykazu zespołów: %w", err)
	}
	defer wiersze.Close()

	wszystkie := []Zespol{}
	for wiersze.Next() {
		zespol, err := odczytajZespol(wiersze)
		if err != nil {
			return nil, 0, err
		}
		wszystkie = append(wszystkie, zespol)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt wykazu zespołów: %w", err)
	}
	if err := r.dolaczSklad(ctx, wszystkie); err != nil {
		return nil, 0, err
	}
	razem := len(wszystkie)
	return utnijWykaz(wszystkie, filtr), razem, nil
}

// oslonWieloznaczniki osłania `%`, `_` i sam znak `\`; replacer idzie napisem jednym przebiegiem.
func oslonWieloznaczniki(fraza string) string {
	zastepnik := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return zastepnik.Replace(fraza)
}

// utnijWykaz nakłada przesunięcie i granicę; przesunięcie za końcem wykazu oddaje pustkę, nie błąd.
func utnijWykaz(wszystkie []Zespol, filtr FiltrZespolow) []Zespol {
	if filtr.Przesuniecie > 0 {
		if filtr.Przesuniecie >= len(wszystkie) {
			return []Zespol{}
		}
		wszystkie = wszystkie[filtr.Przesuniecie:]
	}
	if filtr.Granica > 0 && filtr.Granica < len(wszystkie) {
		wszystkie = wszystkie[:filtr.Granica]
	}
	return wszystkie
}

// PoKodzie zwraca jeden zespół wraz ze składem; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumZespolow) PoKodzie(ctx context.Context, kod string) (Zespol, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zespolPoKodzie)
	if err != nil {
		return Zespol{}, err
	}
	zespol, err := odczytajZespol(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return Zespol{}, fmt.Errorf("dane: zespół %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	if err != nil {
		return Zespol{}, fmt.Errorf("dane: nie można odczytać zespołu %q: %w", kod, err)
	}
	jeden := []Zespol{zespol}
	if err := r.dolaczSklad(ctx, jeden); err != nil {
		return Zespol{}, err
	}
	return jeden[0], nil
}

// Dodaj zakłada nowy zespół wraz ze składem w jednej transakcji i oddaje go odczytanego z bazy po zapisie.
func (r *repozytoriumZespolow) Dodaj(ctx context.Context, zespol Zespol) (Zespol, error) {
	if strings.TrimSpace(zespol.Nazwa) == "" {
		return Zespol{}, fmt.Errorf("dane: zespół %q wymaga nazwy", zespol.Kod)
	}
	teraz := time.Now().UTC().UnixMilli()
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawZespol)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, zespol.Kod, zespol.Nazwa, zespol.Opis, teraz, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można założyć zespołu %q: %w", zespol.Kod, err)
		}
		numer, err := wynik.LastInsertId()
		if err != nil {
			return fmt.Errorf("dane: nieznany numer założonego zespołu %q: %w", zespol.Kod, err)
		}
		return r.zapiszSklad(ctx, transakcja, numer, zespol.Sklad)
	})
	if err != nil {
		return Zespol{}, err
	}
	return r.PoKodzie(ctx, zespol.Kod)
}

// Zapisz zmienia nazwę, opis i cały skład; skład wymienia się w całości, bo `team.save` niesie komplet `agentIds`.
func (r *repozytoriumZespolow) Zapisz(ctx context.Context, zespol Zespol) (Zespol, error) {
	if strings.TrimSpace(zespol.Nazwa) == "" {
		return Zespol{}, fmt.Errorf("dane: zespół %q wymaga nazwy", zespol.Kod)
	}
	numer, err := r.numer(ctx, zespol.Kod)
	if err != nil {
		return Zespol{}, err
	}
	teraz := time.Now().UTC().UnixMilli()
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, aktualizujZespol)
		if err != nil {
			return err
		}
		if _, err := polecenie.ExecContext(ctx, zespol.Nazwa, zespol.Opis, teraz, zespol.Kod, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można zapisać zespołu %q: %w", zespol.Kod, err)
		}
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunSkladZespolu)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, numer); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić składu zespołu %q: %w", zespol.Kod, err)
		}
		return r.zapiszSklad(ctx, transakcja, numer, zespol.Sklad)
	})
	if err != nil {
		return Zespol{}, err
	}
	return r.PoKodzie(ctx, zespol.Kod)
}

func (r *repozytoriumZespolow) numer(ctx context.Context, kod string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, numerZespolu)
	if err != nil {
		return 0, err
	}
	var numer int64
	err = polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&numer)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("dane: zespół %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odczytać numeru zespołu %q: %w", kod, err)
	}
	return numer, nil
}

// zapiszSklad wpisuje skład w kolejności żądania; kod powtórzony wchodzi tylko raz.
func (r *repozytoriumZespolow) zapiszSklad(ctx context.Context, transakcja *sql.Tx,
	numer int64, sklad []string) error {

	if len(sklad) == 0 {
		return nil
	}
	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawSkladZespolu)
	if err != nil {
		return err
	}
	widziane := make(map[string]struct{}, len(sklad))
	kolejnosc := 0
	for _, kod := range sklad {
		kod = strings.TrimSpace(kod)
		if kod == "" {
			continue
		}
		if _, jest := widziane[kod]; jest {
			continue
		}
		widziane[kod] = struct{}{}
		if _, err := polecenie.ExecContext(ctx, numer, kod, kolejnosc); err != nil {
			return fmt.Errorf("dane: nie można zapisać składu zespołu (ekspert %q): %w", kod, err)
		}
		kolejnosc++
	}
	return nil
}

func (r *repozytoriumZespolow) dolaczSklad(ctx context.Context, zespoly []Zespol) error {
	if len(zespoly) == 0 {
		return nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, skladZespolow)
	if err != nil {
		return err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return fmt.Errorf("dane: nie można odczytać składu zespołów: %w", err)
	}
	defer wiersze.Close()

	obecni := map[int64][]string{}
	pominieci := map[int64][]string{}
	for wiersze.Next() {
		var numer int64
		var kod string
		var czynny int
		if err := wiersze.Scan(&numer, &kod, &czynny); err != nil {
			return fmt.Errorf("dane: nie można odczytać wiersza składu zespołu: %w", err)
		}
		if czynny == 1 {
			obecni[numer] = append(obecni[numer], kod)
			continue
		}
		pominieci[numer] = append(pominieci[numer], kod)
	}
	if err := wiersze.Err(); err != nil {
		return fmt.Errorf("dane: przerwany odczyt składu zespołów: %w", err)
	}
	for i := range zespoly {
		zespoly[i].Sklad = obecni[zespoly[i].ID]
		zespoly[i].Pominieci = pominieci[zespoly[i].ID]
	}
	return nil
}

func odczytajZespol(wiersz skaner) (Zespol, error) {
	var zespol Zespol
	var zaktualizowano sql.NullInt64
	if err := wiersz.Scan(&zespol.ID, &zespol.Kod, &zespol.Nazwa, &zespol.Opis,
		&zespol.Utworzono, &zaktualizowano); err != nil {
		return Zespol{}, err
	}
	if zaktualizowano.Valid {
		chwila := zaktualizowano.Int64
		zespol.Zaktualizowano = &chwila
	}
	return zespol, nil
}
