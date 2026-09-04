// Plik prowadzi reguły alarmowania, skarbiec referencji poświadczeń i dziennik audytu automatyki; skarbiec nie zna
// wartości poświadczeń i nigdy jej nie pozna, bo baza trzyma wyłącznie nazwę, zasięg i odwołanie do sejfu plikowego.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// RegulaAlarmowania to wiersz tabeli `regula_alarmowania_automatyki` niosący warunek oraz próg alarmu.
type RegulaAlarmowania struct {
	Kod            string
	AutomatykaID   int64
	AutomatykaKod  string
	Wyzwalacz      string
	Warunek        *string
	Kanaly         string
	Czynna         bool
	Zaktualizowano string
}

// PoswiadczenieAutomatyki to wiersz tabeli `poswiadczenie_automatyki` — sama
// referencja, bez wartości.
type PoswiadczenieAutomatyki struct {
	Odwolanie      string
	Nazwa          string
	Zasieg         *string
	ZasiegID       *string
	Zaktualizowano string
}

// WpisAudytuAutomatyki to wiersz tabeli `wpis_audytu_automatyki` niosący jeden zapis dziennika audytu.
type WpisAudytuAutomatyki struct {
	Kod           string
	AutomatykaID  *int64
	AutomatykaKod *string
	Wykonawca     string
	Czynnosc      string
	Szczegoly     *string
	Chwila        string
}

var (
	kolumnyReguly = `r.identyfikator_zewnetrzny, r.automatyka_id, a.identyfikator_zewnetrzny,
	                 r.wyzwalacz, r.warunek, r.kanaly, r.czynna, r.zaktualizowano`

	zrodloReguly = ` FROM regula_alarmowania_automatyki r
	                 JOIN automatyka a ON a.id = r.automatyka_id`

	zapiszRegule = `INSERT INTO regula_alarmowania_automatyki
	                (identyfikator_zewnetrzny, automatyka_id, wyzwalacz, warunek, kanaly, czynna)
	                VALUES (?, ?, ?, ?, ?, ?)
	                ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                    automatyka_id = excluded.automatyka_id,
	                    wyzwalacz = excluded.wyzwalacz,
	                    warunek = excluded.warunek,
	                    kanaly = excluded.kanaly,
	                    czynna = excluded.czynna,
	                    zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzRegule = `SELECT ` + kolumnyReguly + zrodloReguly +
		` WHERE r.identyfikator_zewnetrzny = ?`

	// Zero w miejscu automatyki znaczy „reguły wszystkich automatyk” — tak
	// mówi kontrakt o żądaniu bez pola `workflowId`.
	listaRegul = `SELECT ` + kolumnyReguly + zrodloReguly +
		` WHERE (? = 0 OR r.automatyka_id = ?) ORDER BY a.nazwa, r.id`

	// Odwołanie jest UNIQUE w całej tabeli, więc warunek przy DO UPDATE
	// zatrzymuje nadpisanie wiersza należącego do innego konta.
	zapiszPoswiadczenieAutomatyki = `INSERT INTO poswiadczenie_automatyki
	                                 (odwolanie, nazwa, zasieg, zasieg_id, konto_id)
	                                 VALUES (?, ?, ?, ?, ` + WskazanieKonta + `)
	                                 ON CONFLICT(odwolanie) DO UPDATE SET
	                                     nazwa = excluded.nazwa,
	                                     zasieg = excluded.zasieg,
	                                     zasieg_id = excluded.zasieg_id,
	                                     zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                 WHERE ` + WarunekKonta

	listaPoswiadczenAutomatyki = `SELECT odwolanie, nazwa, zasieg, zasieg_id, zaktualizowano
	                              FROM poswiadczenie_automatyki
	                              WHERE (? = '' OR zasieg = ?) AND (? = '' OR zasieg_id = ?)
	                                AND ` + WarunekKonta + `
	                              ORDER BY nazwa, odwolanie`

	pobierzPoswiadczenieAutomatyki = `SELECT odwolanie, nazwa, zasieg, zasieg_id, zaktualizowano
	                                  FROM poswiadczenie_automatyki
	                                  WHERE odwolanie = ? AND ` + WarunekKonta

	usunPoswiadczenieAutomatyki = `DELETE FROM poswiadczenie_automatyki
	                               WHERE odwolanie = ? AND ` + WarunekKonta

	dopiszWpisAudytu = `INSERT INTO wpis_audytu_automatyki
	                    (identyfikator_zewnetrzny, automatyka_id, wykonawca, czynnosc, szczegoly)
	                    VALUES (?, ?, ?, ?, ?)`

	// Zakres dat pusty znaczy „bez zawężenia”: znacznik pusty jest
	// leksykograficznie mniejszy od każdego znacznika ISO, a górna granica
	// pusta zdejmuje warunek jawnym porównaniem.
	listaAudytuAutomatyki = `SELECT w.identyfikator_zewnetrzny, w.automatyka_id,
	                                a.identyfikator_zewnetrzny, w.wykonawca, w.czynnosc,
	                                w.szczegoly, w.chwila
	                         FROM wpis_audytu_automatyki w
	                         LEFT JOIN automatyka a ON a.id = w.automatyka_id
	                         WHERE (? = 0 OR w.automatyka_id = ?)
	                           AND (? = '' OR w.chwila >= ?)
	                           AND (? = '' OR w.chwila <= ?)
	                           AND ` + warunekKontaAutomatyki + `
	                         ORDER BY w.chwila DESC, w.id DESC LIMIT ?`
)

// ZapiszRegulealarmowania zakłada regułę alarmowania albo nadpisuje zastaną regułę tego samego zasięgu.
func (r *repozytoriumAutomatyk) ZapiszRegulealarmowania(ctx context.Context,
	regula RegulaAlarmowania) (RegulaAlarmowania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszRegule)
	if err != nil {
		return RegulaAlarmowania{}, err
	}
	_, err = polecenie.ExecContext(ctx, regula.Kod, regula.AutomatykaID, regula.Wyzwalacz,
		tekstDoKolumny(regula.Warunek), regula.Kanaly, liczbaLogiczna(regula.Czynna))
	if err != nil {
		return RegulaAlarmowania{}, fmt.Errorf("dane: nie można zapisać reguły alarmowania %q: %w",
			regula.Kod, err)
	}
	pobranie, err := r.zapytania.przygotuj(ctx, pobierzRegule)
	if err != nil {
		return RegulaAlarmowania{}, err
	}
	zapisana, err := odczytajRegule(pobranie.QueryRowContext(ctx, regula.Kod))
	if err != nil {
		return RegulaAlarmowania{}, fmt.Errorf("dane: nieczytelna reguła alarmowania %q: %w",
			regula.Kod, err)
	}
	return zapisana, nil
}

// RegulyAlarmowania zwraca reguły alarmowania jednej automatyki, albo pełny komplet reguł Operatora z bazy.
func (r *repozytoriumAutomatyk) RegulyAlarmowania(ctx context.Context,
	automatykaID int64) ([]RegulaAlarmowania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaRegul)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID, automatykaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać reguł alarmowania: %w", err)
	}
	defer wiersze.Close()

	lista := []RegulaAlarmowania{}
	for wiersze.Next() {
		regula, err := odczytajRegule(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz reguły alarmowania: %w", err)
		}
		lista = append(lista, regula)
	}
	return lista, wiersze.Err()
}

// ZapiszPoswiadczenieAutomatyki zapisuje REFERENCJĘ poświadczenia. Wartość idzie
// osobno, do sejfu — ta metoda jej nie widzi i widzieć nie może.
func (r *repozytoriumAutomatyk) ZapiszPoswiadczenieAutomatyki(ctx context.Context,
	poswiadczenie PoswiadczenieAutomatyki) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPoswiadczenieAutomatyki)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, poswiadczenie.Odwolanie, poswiadczenie.Nazwa,
		tekstDoKolumny(poswiadczenie.Zasieg), tekstDoKolumny(poswiadczenie.ZasiegID),
		KontoOperatora(ctx), KontoOperatora(ctx))
	// Trójka nazwa-zasięg-zasięg_id jest UNIQUE poza celem ON CONFLICT, więc
	// nazwa zajęta przez inne konto rozbija wstawienie o wiąz, a nie o warunek.
	if czyKolizja(err) {
		return fmt.Errorf("dane: nazwa poświadczenia %q jest zajęta w tym zasięgu: %w",
			poswiadczenie.Nazwa, ErrKolizjaWiersza)
	}
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać referencji poświadczenia %q: %w",
			poswiadczenie.Nazwa, err)
	}
	return sprawdzTrafienieZapisu(wynik, "odwołanie poświadczenia", poswiadczenie.Odwolanie)
}

// PoswiadczeniaAutomatyki zwraca referencje poświadczeń w kolejności nazw wprost z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) PoswiadczeniaAutomatyki(ctx context.Context,
	zasieg, zasiegID string) ([]PoswiadczenieAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPoswiadczenAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, zasieg, zasieg, zasiegID, zasiegID,
		KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać referencji poświadczeń: %w", err)
	}
	defer wiersze.Close()

	lista := []PoswiadczenieAutomatyki{}
	for wiersze.Next() {
		poswiadczenie, err := odczytajPoswiadczenieAutomatyki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz referencji poświadczenia: %w", err)
		}
		lista = append(lista, poswiadczenie)
	}
	return lista, wiersze.Err()
}

// PoswiadczenieAutomatykiPoOdwolaniu zwraca jedną referencję poświadczenia po jej odwołaniu w skarbcu.
func (r *repozytoriumAutomatyk) PoswiadczenieAutomatykiPoOdwolaniu(ctx context.Context,
	odwolanie string) (PoswiadczenieAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPoswiadczenieAutomatyki)
	if err != nil {
		return PoswiadczenieAutomatyki{}, err
	}
	poswiadczenie, err := odczytajPoswiadczenieAutomatyki(
		polecenie.QueryRowContext(ctx, odwolanie, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return PoswiadczenieAutomatyki{}, ErrBrakWiersza
	}
	if err != nil {
		return PoswiadczenieAutomatyki{}, fmt.Errorf("dane: nieczytelna referencja %q: %w", odwolanie, err)
	}
	return poswiadczenie, nil
}

// UsunPoswiadczenieAutomatyki kasuje referencję i mówi, czy była. Wartość
// z sejfu zdejmuje wołający — sejf leży poza bazą i poza transakcją.
func (r *repozytoriumAutomatyk) UsunPoswiadczenieAutomatyki(ctx context.Context,
	odwolanie string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunPoswiadczenieAutomatyki)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, odwolanie, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć referencji %q: %w", odwolanie, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia referencji %q: %w", odwolanie, err)
	}
	return zmienione > 0, nil
}

// DopiszAudytAutomatyki nanosi wiersz dziennika audytu. Dziennik jest zapisem
// niezmiennym — żadna komenda modułu go nie zmienia ani nie kasuje.
func (r *repozytoriumAutomatyk) DopiszAudytAutomatyki(ctx context.Context,
	wpis WpisAudytuAutomatyki) error {

	polecenie, err := r.zapytania.przygotuj(ctx, dopiszWpisAudytu)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wpis.Kod, liczbaDoKolumny(wpis.AutomatykaID),
		wpis.Wykonawca, wpis.Czynnosc, tekstDoKolumny(wpis.Szczegoly))
	if err != nil {
		return fmt.Errorf("dane: nie można dopisać wpisu audytu %q: %w", wpis.Czynnosc, err)
	}
	return nil
}

// AudytAutomatyki zwraca dziennik audytu automatyki od najnowszego wpisu wprost z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) AudytAutomatyki(ctx context.Context, automatykaID int64,
	od, do string, limit int) ([]WpisAudytuAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaAudytuAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID, automatykaID,
		od, od, do, do, KontoOperatora(ctx), granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dziennika audytu: %w", err)
	}
	defer wiersze.Close()

	lista := []WpisAudytuAutomatyki{}
	for wiersze.Next() {
		var wpis WpisAudytuAutomatyki
		var automatyka sql.NullInt64
		var kodAutomatyki, szczegoly sql.NullString
		err := wiersze.Scan(&wpis.Kod, &automatyka, &kodAutomatyki, &wpis.Wykonawca,
			&wpis.Czynnosc, &szczegoly, &wpis.Chwila)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz audytu: %w", err)
		}
		wpis.AutomatykaID = liczbaZKolumny(automatyka)
		wpis.AutomatykaKod = tekstZKolumny(kodAutomatyki)
		wpis.Szczegoly = tekstZKolumny(szczegoly)
		lista = append(lista, wpis)
	}
	return lista, wiersze.Err()
}

// odczytajRegule składa regułę alarmowania automatyki wprost z jednego wiersza wyniku zapytania do bazy.
func odczytajRegule(wiersz skaner) (RegulaAlarmowania, error) {
	var regula RegulaAlarmowania
	var warunek sql.NullString
	var czynna int
	err := wiersz.Scan(&regula.Kod, &regula.AutomatykaID, &regula.AutomatykaKod,
		&regula.Wyzwalacz, &warunek, &regula.Kanaly, &czynna, &regula.Zaktualizowano)
	if err != nil {
		return RegulaAlarmowania{}, err
	}
	regula.Warunek = tekstZKolumny(warunek)
	regula.Czynna = czynna == 1
	return regula, nil
}

// odczytajPoswiadczenieAutomatyki składa referencję wprost z jednego wiersza wyniku zapytania do bazy.
func odczytajPoswiadczenieAutomatyki(wiersz skaner) (PoswiadczenieAutomatyki, error) {
	var poswiadczenie PoswiadczenieAutomatyki
	var zasieg, zasiegID sql.NullString
	err := wiersz.Scan(&poswiadczenie.Odwolanie, &poswiadczenie.Nazwa, &zasieg,
		&zasiegID, &poswiadczenie.Zaktualizowano)
	if err != nil {
		return PoswiadczenieAutomatyki{}, err
	}
	poswiadczenie.Zasieg = tekstZKolumny(zasieg)
	poswiadczenie.ZasiegID = tekstZKolumny(zasiegID)
	return poswiadczenie, nil
}
