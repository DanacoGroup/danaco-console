// Odpowiedzialność pliku: głosowania debaty, czyli otwarcie, warianty i oddane głosy; wyniku agregacji tu nie ma, liczy go rdzeń przy odczycie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// GlosowanieDebaty to wiersz tabeli debata_glosowanie, niosący stan i metodę agregacji tego głosowania.
type GlosowanieDebaty struct {
	Kod        string
	Okno       string
	Tura       string
	Metoda     string
	Stan       string
	Prog       float64
	Uprawnieni []string
	Rozpoczeto string
	Zamknieto  *string
}

// WariantDebaty to jeden wariant poddany pod głosowanie, wraz z jego treścią i kolejnością jego podania.
type WariantDebaty struct {
	Kod        string
	Glosowanie string
	Etykieta   string
	Wypowiedz  string
	Kolejnosc  int
}

// GlosDebaty to jeden oddany głos; kształt zależy od metody agregacji, więc wypełnione bywa jedno z trzech pól.
type GlosDebaty struct {
	Kod        string
	Glosowanie string
	Wyborca    string
	Aprobaty   []string
	Ranking    []string
	PunktyJson string
	Oddano     string
}

// RepozytoriumDebatyGlosowan jest częścią kontraktu całego obszaru Roundtable odpowiadającą za głosowania.
type RepozytoriumDebatyGlosowan interface {
	ZalozGlosowanieDebaty(ctx context.Context, glosowanie GlosowanieDebaty,
		warianty []WariantDebaty) (GlosowanieDebaty, error)
	GlosowanieDebatyPoKodzie(ctx context.Context, kod string) (GlosowanieDebaty, error)
	OstatnieGlosowanieDebaty(ctx context.Context, okno string) (GlosowanieDebaty, error)
	WariantyDebaty(ctx context.Context, glosowanie string) ([]WariantDebaty, error)
	OddajGlosDebaty(ctx context.Context, glos GlosDebaty) (GlosDebaty, error)
	GlosyDebaty(ctx context.Context, glosowanie string) ([]GlosDebaty, error)
	UstawStanGlosowaniaDebaty(ctx context.Context, kod, stan string, zamknieto *string) error
}

const (
	kolumnyGlosowaniaDebaty = `identyfikator_zewnetrzny, okno, tura, metoda, stan, prog,
	                           uprawnieni, rozpoczeto, zamknieto`

	zalozGlosowanieDebaty = `INSERT INTO debata_glosowanie
	                         (identyfikator_zewnetrzny, okno, tura, metoda, stan, prog, uprawnieni, konto_id)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	zalozWariantDebaty = `INSERT INTO debata_wariant
	                      (identyfikator_zewnetrzny, glosowanie, etykieta, wypowiedz, kolejnosc)
	                      VALUES (?, ?, ?, ?, ?)`

	pobierzGlosowanieDebaty = `SELECT ` + kolumnyGlosowaniaDebaty + `
	                           FROM debata_glosowanie WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	pobierzOstatnieGlosowanieDebaty = `SELECT ` + kolumnyGlosowaniaDebaty + `
	                                   FROM debata_glosowanie WHERE okno = ? AND ` + WarunekKonta + `
	                                   ORDER BY id DESC LIMIT 1`

	pobierzWariantyDebaty = `SELECT identyfikator_zewnetrzny, glosowanie, etykieta, wypowiedz,
	                                kolejnosc
	                         FROM debata_wariant WHERE glosowanie = ?
	                         ORDER BY kolejnosc ASC, id ASC`

	// Powtórne oddanie głosu zastępuje poprzedni: zmiana zdania jest dozwolona, dwa głosy tej samej osoby nie.
	oddajGlosDebaty = `INSERT INTO debata_glos
	                   (identyfikator_zewnetrzny, glosowanie, wyborca, aprobaty, ranking, punkty_json)
	                   VALUES (?, ?, ?, ?, ?, ?)
	                   ON CONFLICT(glosowanie, wyborca) DO UPDATE SET
	                       aprobaty = excluded.aprobaty,
	                       ranking = excluded.ranking,
	                       punkty_json = excluded.punkty_json,
	                       oddano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzGlosDebaty = `SELECT identyfikator_zewnetrzny, glosowanie, wyborca, aprobaty, ranking,
	                            punkty_json, oddano
	                     FROM debata_glos WHERE glosowanie = ? AND wyborca = ?`

	pobierzGlosyDebaty = `SELECT identyfikator_zewnetrzny, glosowanie, wyborca, aprobaty, ranking,
	                             punkty_json, oddano
	                      FROM debata_glos WHERE glosowanie = ? ORDER BY id ASC`

	ustawStanGlosowaniaDebaty = `UPDATE debata_glosowanie SET stan = ?, zamknieto = ?
	                             WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

// ZalozGlosowanieDebaty otwiera głosowanie wraz z wariantami w jednej transakcji.
func (r *repozytoriumRoundtable) ZalozGlosowanieDebaty(ctx context.Context,
	glosowanie GlosowanieDebaty, warianty []WariantDebaty) (GlosowanieDebaty, error) {

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		naglowek, err := r.zapytania.wTransakcji(ctx, transakcja, zalozGlosowanieDebaty)
		if err != nil {
			return err
		}
		if _, err := naglowek.ExecContext(ctx, glosowanie.Kod, glosowanie.Okno, glosowanie.Tura,
			glosowanie.Metoda, glosowanie.Stan, glosowanie.Prog,
			strings.Join(glosowanie.Uprawnieni, "\n"), KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można otworzyć głosowania debaty %q: %w", glosowanie.Kod, err)
		}
		wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, zalozWariantDebaty)
		if err != nil {
			return err
		}
		for pozycja, wariant := range warianty {
			if _, err := wstaw.ExecContext(ctx, wariant.Kod, glosowanie.Kod, wariant.Etykieta,
				wariant.Wypowiedz, pozycja+1); err != nil {
				return fmt.Errorf("dane: nie można zapisać wariantu głosowania %q: %w",
					wariant.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return GlosowanieDebaty{}, err
	}
	return r.GlosowanieDebatyPoKodzie(ctx, glosowanie.Kod)
}

// GlosowanieDebatyPoKodzie zwraca głosowanie po jego identyfikatorze zewnętrznym, wraz z jego wariantami.
func (r *repozytoriumRoundtable) GlosowanieDebatyPoKodzie(ctx context.Context,
	kod string) (GlosowanieDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzGlosowanieDebaty)
	if err != nil {
		return GlosowanieDebaty{}, err
	}
	return odczytajGlosowanieDebaty(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
}

// OstatnieGlosowanieDebaty zwraca ostatnio otwarte głosowanie danego okna operacyjnego tej samej debaty.
func (r *repozytoriumRoundtable) OstatnieGlosowanieDebaty(ctx context.Context,
	okno string) (GlosowanieDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOstatnieGlosowanieDebaty)
	if err != nil {
		return GlosowanieDebaty{}, err
	}
	return odczytajGlosowanieDebaty(polecenie.QueryRowContext(ctx, okno, KontoOperatora(ctx)))
}

func (r *repozytoriumRoundtable) WariantyDebaty(ctx context.Context,
	glosowanie string) ([]WariantDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWariantyDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, glosowanie)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wariantów głosowania %q: %w", glosowanie, err)
	}
	defer wiersze.Close()

	warianty := make([]WariantDebaty, 0, 8)
	for wiersze.Next() {
		var wariant WariantDebaty
		if err := wiersze.Scan(&wariant.Kod, &wariant.Glosowanie, &wariant.Etykieta,
			&wariant.Wypowiedz, &wariant.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wariantu głosowania: %w", err)
		}
		warianty = append(warianty, wariant)
	}
	return warianty, wiersze.Err()
}

// OddajGlosDebaty zapisuje głos danego uczestnika tej debaty i oddaje ten sam głos odczytany po zapisie.
func (r *repozytoriumRoundtable) OddajGlosDebaty(ctx context.Context,
	glos GlosDebaty) (GlosDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, oddajGlosDebaty)
	if err != nil {
		return GlosDebaty{}, err
	}
	if _, err := polecenie.ExecContext(ctx, glos.Kod, glos.Glosowanie, glos.Wyborca,
		strings.Join(glos.Aprobaty, "\n"), strings.Join(glos.Ranking, "\n"),
		glos.PunktyJson); err != nil {
		return GlosDebaty{}, fmt.Errorf("dane: nie można zapisać głosu %q: %w", glos.Kod, err)
	}

	odczyt, err := r.zapytania.przygotuj(ctx, pobierzGlosDebaty)
	if err != nil {
		return GlosDebaty{}, err
	}
	zapisany, err := odczytajGlosDebaty(odczyt.QueryRowContext(ctx, glos.Glosowanie, glos.Wyborca))
	if err != nil {
		return GlosDebaty{}, err
	}
	return zapisany, nil
}

func (r *repozytoriumRoundtable) GlosyDebaty(ctx context.Context,
	glosowanie string) ([]GlosDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzGlosyDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, glosowanie)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać głosów głosowania %q: %w", glosowanie, err)
	}
	defer wiersze.Close()

	glosy := make([]GlosDebaty, 0, 8)
	for wiersze.Next() {
		glos, err := odczytajGlosDebaty(wiersze)
		if err != nil {
			return nil, err
		}
		glosy = append(glosy, glos)
	}
	return glosy, wiersze.Err()
}

// UstawStanGlosowaniaDebaty zamyka wskazane głosowanie danego okna albo znakuje je jako zakończone remisem.
func (r *repozytoriumRoundtable) UstawStanGlosowaniaDebaty(ctx context.Context,
	kod, stan string, zamknieto *string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawStanGlosowaniaDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, stan, zamknieto, kod, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić stanu głosowania %q: %w", kod, err)
	}
	return trafienieDebaty(wynik)
}

func odczytajGlosowanieDebaty(wiersz interface{ Scan(...any) error }) (GlosowanieDebaty, error) {
	var glosowanie GlosowanieDebaty
	var uprawnieni string
	err := wiersz.Scan(&glosowanie.Kod, &glosowanie.Okno, &glosowanie.Tura, &glosowanie.Metoda,
		&glosowanie.Stan, &glosowanie.Prog, &uprawnieni, &glosowanie.Rozpoczeto,
		&glosowanie.Zamknieto)
	if errors.Is(err, sql.ErrNoRows) {
		return GlosowanieDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return GlosowanieDebaty{}, fmt.Errorf("dane: nieczytelny wiersz głosowania debaty: %w", err)
	}
	glosowanie.Uprawnieni = rozdzielWierszeDebaty(uprawnieni)
	return glosowanie, nil
}

func odczytajGlosDebaty(wiersz interface{ Scan(...any) error }) (GlosDebaty, error) {
	var glos GlosDebaty
	var aprobaty, ranking string
	err := wiersz.Scan(&glos.Kod, &glos.Glosowanie, &glos.Wyborca, &aprobaty, &ranking,
		&glos.PunktyJson, &glos.Oddano)
	if errors.Is(err, sql.ErrNoRows) {
		return GlosDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return GlosDebaty{}, fmt.Errorf("dane: nieczytelny wiersz głosu debaty: %w", err)
	}
	glos.Aprobaty = rozdzielWierszeDebaty(aprobaty)
	glos.Ranking = rozdzielWierszeDebaty(ranking)
	return glos, nil
}

// rozdzielWierszeDebaty rozbija wykaz z jednej kolumny; tekst pusty oddaje wykaz pusty, nie jedną pustą pozycję.
func rozdzielWierszeDebaty(zapis string) []string {
	if strings.TrimSpace(zapis) == "" {
		return nil
	}
	czesci := strings.Split(zapis, "\n")
	wykaz := make([]string, 0, len(czesci))
	for _, czesc := range czesci {
		if przyciete := strings.TrimSpace(czesc); przyciete != "" {
			wykaz = append(wykaz, przyciete)
		}
	}
	return wykaz
}
