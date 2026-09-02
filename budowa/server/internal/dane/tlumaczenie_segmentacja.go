// Plik obsługuje segmentację modułu Translate: zestawy reguł podziału (zestaw_regul_segmentacji,
// regula_segmentacji) oraz trwałe segmenty okna (segment_okna_tlumaczenia), wprowadzone migracją 161.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// `Srx` niesie treść pliku SRX w całości — standard branżowy, którego rdzeń nie rozkłada.
type ZestawRegulSegmentacji struct {
	ID             int64
	Kod            string
	Nazwa          string
	Jezyk          *string
	Srx            *string
	Reguly         []RegulaSegmentacji
	Zaktualizowano int64
}

type RegulaSegmentacji struct {
	Kolejnosc int64
	Przed     string
	Po        string
	Lamie     bool
}

// warunekOknaSegmentu prowadzi segment do konta przez okno_tlumaczenia (konto_id od migracji 484).
const warunekOknaSegmentu = `EXISTS (SELECT 1 FROM okno_tlumaczenia
	WHERE okno_tlumaczenia.id = segment_okna_tlumaczenia.okno_id AND ` + WarunekKonta + `)`

func (r *repozytoriumTlumaczen) ZestawyRegulSegmentacji(ctx context.Context,
	jezyk string) ([]ZestawRegulSegmentacji, error) {

	zapytanie := `SELECT id, identyfikator_zewnetrzny, nazwa, jezyk, srx, zaktualizowano
	                FROM zestaw_regul_segmentacji WHERE ` + WarunekKonta
	argumenty := []any{KontoOperatora(ctx)}
	if strings.TrimSpace(jezyk) != "" {
		zapytanie += " AND jezyk = ?"
		argumenty = append(argumenty, jezyk)
	}
	zapytanie += " ORDER BY nazwa"

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zestawów reguł segmentacji: %w", err)
	}
	defer wiersze.Close()

	zestawy := []ZestawRegulSegmentacji{}
	for wiersze.Next() {
		var zestaw ZestawRegulSegmentacji
		var jezykKolumna, srx sql.NullString
		if err := wiersze.Scan(&zestaw.ID, &zestaw.Kod, &zestaw.Nazwa, &jezykKolumna,
			&srx, &zestaw.Zaktualizowano); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny zestaw reguł segmentacji: %w", err)
		}
		zestaw.Jezyk = tekstZKolumny(jezykKolumna)
		zestaw.Srx = tekstZKolumny(srx)
		zestawy = append(zestawy, zestaw)
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	for i := range zestawy {
		reguly, err := r.regulySegmentacji(ctx, zestawy[i].ID)
		if err != nil {
			return nil, err
		}
		zestawy[i].Reguly = reguly
	}
	return zestawy, nil
}

func (r *repozytoriumTlumaczen) regulySegmentacji(ctx context.Context,
	zestawID int64) ([]RegulaSegmentacji, error) {

	wiersze, err := r.db.QueryContext(ctx,
		`SELECT kolejnosc, przed, po, lamie FROM regula_segmentacji
		  WHERE zestaw_id = ? ORDER BY kolejnosc`, zestawID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać reguł segmentacji: %w", err)
	}
	defer wiersze.Close()

	reguly := []RegulaSegmentacji{}
	for wiersze.Next() {
		var regula RegulaSegmentacji
		var lamie int64
		if err := wiersze.Scan(&regula.Kolejnosc, &regula.Przed, &regula.Po, &lamie); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna reguła segmentacji: %w", err)
		}
		regula.Lamie = lamie == 1
		reguly = append(reguly, regula)
	}
	return reguly, wiersze.Err()
}

// Kontrakt nadsyła wykaz reguł kompletem, więc reguły idą na wymianę w całości.
// Warunek przy DO UPDATE zostawia zestaw cudzego konta nietknięty; zapis kończy się ErrKolizjaWiersza.
func (r *repozytoriumTlumaczen) ZapiszZestawRegulSegmentacji(ctx context.Context,
	zestaw ZestawRegulSegmentacji) (ZestawRegulSegmentacji, error) {

	if strings.TrimSpace(zestaw.Kod) == "" {
		return ZestawRegulSegmentacji{}, fmt.Errorf("dane: zestaw reguł segmentacji bez identyfikatora")
	}
	teraz := time.Now().UnixMilli()
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wynik, err := transakcja.ExecContext(ctx, `INSERT INTO zestaw_regul_segmentacji
			(identyfikator_zewnetrzny, nazwa, jezyk, srx, zaktualizowano, konto_id)
			VALUES (?, ?, ?, ?, ?, `+WskazanieKonta+`)
			ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
				nazwa = excluded.nazwa, jezyk = excluded.jezyk, srx = excluded.srx,
				zaktualizowano = excluded.zaktualizowano
			WHERE `+WarunekKonta,
			zestaw.Kod, zestaw.Nazwa, tekstDoKolumny(zestaw.Jezyk),
			tekstDoKolumny(zestaw.Srx), teraz, KontoOperatora(ctx), KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać zestawu reguł segmentacji %q: %w", zestaw.Kod, err)
		}
		if err := sprawdzTrafienieZapisu(wynik, "zestaw reguł segmentacji", zestaw.Kod); err != nil {
			return err
		}
		var zestawID int64
		if err := transakcja.QueryRowContext(ctx,
			`SELECT id FROM zestaw_regul_segmentacji WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta,
			zestaw.Kod, KontoOperatora(ctx)).Scan(&zestawID); err != nil {
			return fmt.Errorf("dane: nie można odczytać zestawu reguł segmentacji %q: %w", zestaw.Kod, err)
		}
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM regula_segmentacji WHERE zestaw_id = ?`, zestawID); err != nil {
			return fmt.Errorf("dane: nie można zdjąć reguł zestawu %q: %w", zestaw.Kod, err)
		}
		for _, regula := range zestaw.Reguly {
			if _, err := transakcja.ExecContext(ctx, `INSERT INTO regula_segmentacji
				(zestaw_id, kolejnosc, przed, po, lamie) VALUES (?, ?, ?, ?, ?)`,
				zestawID, regula.Kolejnosc, regula.Przed, regula.Po,
				wartoscLogicznaDoKolumny(regula.Lamie)); err != nil {
				return fmt.Errorf("dane: nie można zapisać reguły zestawu %q: %w", zestaw.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return ZestawRegulSegmentacji{}, err
	}
	return r.ZestawRegulSegmentacji(ctx, zestaw.Kod)
}

func (r *repozytoriumTlumaczen) ZestawRegulSegmentacji(ctx context.Context,
	kod string) (ZestawRegulSegmentacji, error) {

	var zestaw ZestawRegulSegmentacji
	var jezyk, srx sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, identyfikator_zewnetrzny, nazwa, jezyk, srx, zaktualizowano
		   FROM zestaw_regul_segmentacji WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta,
		kod, KontoOperatora(ctx)).
		Scan(&zestaw.ID, &zestaw.Kod, &zestaw.Nazwa, &jezyk, &srx, &zestaw.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return ZestawRegulSegmentacji{}, ErrBrakWiersza
	}
	if err != nil {
		return ZestawRegulSegmentacji{}, fmt.Errorf("dane: nieczytelny zestaw reguł %q: %w", kod, err)
	}
	zestaw.Jezyk = tekstZKolumny(jezyk)
	zestaw.Srx = tekstZKolumny(srx)
	reguly, err := r.regulySegmentacji(ctx, zestaw.ID)
	if err != nil {
		return ZestawRegulSegmentacji{}, err
	}
	zestaw.Reguly = reguly
	return zestaw, nil
}

func (r *repozytoriumTlumaczen) SegmentyOkna(ctx context.Context, oknoID int64) ([]string, error) {
	wiersze, err := r.db.QueryContext(ctx,
		`SELECT tresc FROM segment_okna_tlumaczenia
		  WHERE okno_id = ? AND `+warunekOknaSegmentu+` ORDER BY kolejnosc`,
		oknoID, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać segmentów okna %d: %w", oknoID, err)
	}
	defer wiersze.Close()

	segmenty := []string{}
	for wiersze.Next() {
		var tresc string
		if err := wiersze.Scan(&tresc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny segment okna %d: %w", oknoID, err)
		}
		segmenty = append(segmenty, tresc)
	}
	return segmenty, wiersze.Err()
}

// Jedna transakcja: scalenie segmentów przesuwa numerację kolejnych wpisów.
func (r *repozytoriumTlumaczen) UstawSegmentyOkna(ctx context.Context,
	oknoID int64, segmenty []string) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM segment_okna_tlumaczenia WHERE okno_id = ? AND `+warunekOknaSegmentu,
			oknoID, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można zdjąć segmentów okna %d: %w", oknoID, err)
		}
		for numer, tresc := range segmenty {
			wynik, err := transakcja.ExecContext(ctx,
				`INSERT INTO segment_okna_tlumaczenia (okno_id, kolejnosc, tresc)
				 SELECT ?, ?, ? WHERE EXISTS (SELECT 1 FROM okno_tlumaczenia
				                              WHERE okno_tlumaczenia.id = ? AND `+WarunekKonta+`)`,
				oknoID, numer, tresc, oknoID, KontoOperatora(ctx))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać segmentu %d okna %d: %w", numer, oknoID, err)
			}
			if err := sprawdzTrafienieZapisu(wynik, "okno tłumaczenia", fmt.Sprintf("%d", oknoID)); err != nil {
				return err
			}
		}
		return nil
	})
}
