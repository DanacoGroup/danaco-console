// Pamięć tłumaczeń modułu Translate (`pamiec_tlumaczen`, migracja 160) i polityka
// pamięci okna (`polityka_pamieci_okna`, dziecko `okno_tlumaczenia`).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type WpisPamieciTlumaczenPelny struct {
	ID                int64
	Kod               string
	PanelID           *int64
	PanelKod          *string
	Jezyk             string
	SegmentZrodlowy   string
	SegmentDocelowy   string
	Projekt           *string
	Klient            *string
	Autor             *string
	KontekstPoprzedni *string
	KontekstNastepny  *string
	Zasieg            string
	Utworzono         int64
	Zaktualizowano    int64
}

// Pole puste filtru nie zawęża wykazu.
type FiltrPamieciTlumaczen struct {
	Jezyk   string
	Fraza   string
	Projekt string
	Zasieg  string
	Limit   int
	Offset  int
}

type PolitykaPamieciOkna struct {
	OknoID               int64
	Zasieg               string
	Prog                 int64
	DopasowanieKontekstu bool
	WstepneTlumaczenie   bool
	Zaktualizowano       int64
}

const kolumnyWpisuPamieciTlumaczen = `w.id, w.identyfikator_zewnetrzny, w.panel_id, p.identyfikator_zewnetrzny,
	w.jezyk, w.segment_zrodlowy, w.segment_docelowy, w.projekt, w.klient, w.autor,
	w.kontekst_poprzedni, w.kontekst_nastepny, w.zasieg, w.utworzono, w.zaktualizowano`

const zrodloWpisuPamieciTlumaczen = ` FROM pamiec_tlumaczen w
	LEFT JOIN panel_tlumaczenia p ON p.id = w.panel_id`

// Tabela `polityka_pamieci_okna` własnej kolumny konta nie ma; granica idzie
// drogą po `okno_id` do `okno_tlumaczenia` (konto_id z migracji 484).
const warunekOknaPolitykiPamieci = `EXISTS (SELECT 1 FROM okno_tlumaczenia
	WHERE okno_tlumaczenia.id = polityka_pamieci_okna.okno_id AND ` + WarunekKonta + `)`

func warunkiPamieci(filtr FiltrPamieciTlumaczen) (string, []any) {
	warunki := []string{}
	argumenty := []any{}
	if strings.TrimSpace(filtr.Jezyk) != "" {
		warunki = append(warunki, "w.jezyk = ?")
		argumenty = append(argumenty, filtr.Jezyk)
	}
	if strings.TrimSpace(filtr.Projekt) != "" {
		warunki = append(warunki, "w.projekt = ?")
		argumenty = append(argumenty, filtr.Projekt)
	}
	if strings.TrimSpace(filtr.Zasieg) != "" {
		warunki = append(warunki, "w.zasieg = ?")
		argumenty = append(argumenty, filtr.Zasieg)
	}
	if fraza := strings.TrimSpace(filtr.Fraza); fraza != "" {
		// Kontrakt pola `query`: dopasowanie w dowolnym miejscu obu segmentów.
		warunki = append(warunki, "(w.segment_zrodlowy LIKE ? OR w.segment_docelowy LIKE ?)")
		wzorzec := "%" + fraza + "%"
		argumenty = append(argumenty, wzorzec, wzorzec)
	}
	if len(warunki) == 0 {
		return "", argumenty
	}
	return " WHERE " + strings.Join(warunki, " AND "), argumenty
}

// Liczba wszystkich pasujących par liczona jest bez limitu, wymóg kontraktu wykazu.
func (r *repozytoriumTlumaczen) WpisyPamieci(ctx context.Context,
	filtr FiltrPamieciTlumaczen) ([]WpisPamieciTlumaczenPelny, int, error) {

	warunek, argumenty := warunkiPamieci(filtr)

	var razem int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pamiec_tlumaczen w`+warunek, argumenty...).Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć pamięci tłumaczeń: %w", err)
	}

	zapytanie := `SELECT ` + kolumnyWpisuPamieciTlumaczen + zrodloWpisuPamieciTlumaczen + warunek +
		` ORDER BY w.zaktualizowano DESC, w.id DESC`
	if filtr.Limit > 0 {
		zapytanie += " LIMIT ?"
		argumenty = append(argumenty, filtr.Limit)
		if filtr.Offset > 0 {
			zapytanie += " OFFSET ?"
			argumenty = append(argumenty, filtr.Offset)
		}
	}

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać pamięci tłumaczeń: %w", err)
	}
	defer wiersze.Close()

	lista := []WpisPamieciTlumaczenPelny{}
	for wiersze.Next() {
		wpis, err := odczytajWpisPamieciPelny(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz pamięci tłumaczeń: %w", err)
		}
		lista = append(lista, wpis)
	}
	return lista, razem, wiersze.Err()
}

func odczytajWpisPamieciPelny(wiersz skaner) (WpisPamieciTlumaczenPelny, error) {
	var wpis WpisPamieciTlumaczenPelny
	var panelID sql.NullInt64
	var panelKod, projekt, klient, autor, poprzedni, nastepny sql.NullString
	err := wiersz.Scan(&wpis.ID, &wpis.Kod, &panelID, &panelKod, &wpis.Jezyk,
		&wpis.SegmentZrodlowy, &wpis.SegmentDocelowy, &projekt, &klient, &autor,
		&poprzedni, &nastepny, &wpis.Zasieg, &wpis.Utworzono, &wpis.Zaktualizowano)
	if err != nil {
		return WpisPamieciTlumaczenPelny{}, err
	}
	wpis.PanelID = liczbaZKolumny(panelID)
	wpis.PanelKod = tekstZKolumny(panelKod)
	wpis.Projekt = tekstZKolumny(projekt)
	wpis.Klient = tekstZKolumny(klient)
	wpis.Autor = tekstZKolumny(autor)
	wpis.KontekstPoprzedni = tekstZKolumny(poprzedni)
	wpis.KontekstNastepny = tekstZKolumny(nastepny)
	return wpis, nil
}

func (r *repozytoriumTlumaczen) WpisPamieci(ctx context.Context,
	kod string) (WpisPamieciTlumaczenPelny, error) {

	wiersz := r.db.QueryRowContext(ctx,
		`SELECT `+kolumnyWpisuPamieciTlumaczen+zrodloWpisuPamieciTlumaczen+` WHERE w.identyfikator_zewnetrzny = ?`, kod)
	wpis, err := odczytajWpisPamieciPelny(wiersz)
	if errors.Is(err, sql.ErrNoRows) {
		return WpisPamieciTlumaczenPelny{}, ErrBrakWiersza
	}
	if err != nil {
		return WpisPamieciTlumaczenPelny{}, fmt.Errorf("dane: nieczytelny wpis pamięci %q: %w", kod, err)
	}
	return wpis, nil
}

func (r *repozytoriumTlumaczen) ZapiszWpisPamieci(ctx context.Context,
	wpis WpisPamieciTlumaczenPelny) (WpisPamieciTlumaczenPelny, error) {

	if strings.TrimSpace(wpis.Kod) == "" {
		return WpisPamieciTlumaczenPelny{}, fmt.Errorf("dane: wpis pamięci bez identyfikatora")
	}
	teraz := time.Now().UnixMilli()
	utworzono := wpis.Utworzono
	if utworzono == 0 {
		utworzono = teraz
	}
	zasieg := wpis.Zasieg
	if zasieg == "" {
		zasieg = "card"
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO pamiec_tlumaczen
		(identyfikator_zewnetrzny, panel_id, jezyk, segment_zrodlowy, segment_docelowy,
		 projekt, klient, autor, kontekst_poprzedni, kontekst_nastepny, zasieg,
		 utworzono, zaktualizowano)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
			panel_id = excluded.panel_id,
			jezyk = excluded.jezyk,
			segment_zrodlowy = excluded.segment_zrodlowy,
			segment_docelowy = excluded.segment_docelowy,
			projekt = excluded.projekt,
			klient = excluded.klient,
			autor = excluded.autor,
			kontekst_poprzedni = excluded.kontekst_poprzedni,
			kontekst_nastepny = excluded.kontekst_nastepny,
			zasieg = excluded.zasieg,
			zaktualizowano = excluded.zaktualizowano`,
		wpis.Kod, liczbaDoKolumny(wpis.PanelID), wpis.Jezyk, wpis.SegmentZrodlowy,
		wpis.SegmentDocelowy, tekstDoKolumny(wpis.Projekt), tekstDoKolumny(wpis.Klient),
		tekstDoKolumny(wpis.Autor), tekstDoKolumny(wpis.KontekstPoprzedni),
		tekstDoKolumny(wpis.KontekstNastepny), zasieg, utworzono, teraz)
	if err != nil {
		return WpisPamieciTlumaczenPelny{}, fmt.Errorf("dane: nie można zapisać wpisu pamięci %q: %w", wpis.Kod, err)
	}
	return r.WpisPamieci(ctx, wpis.Kod)
}

func (r *repozytoriumTlumaczen) UsunWpisPamieci(ctx context.Context, kod string) (bool, error) {
	wynik, err := r.db.ExecContext(ctx,
		`DELETE FROM pamiec_tlumaczen WHERE identyfikator_zewnetrzny = ?`, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć wpisu pamięci %q: %w", kod, err)
	}
	zeszlo, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku usunięcia wpisu pamięci: %w", err)
	}
	return zeszlo > 0, nil
}

// `memory.maintain` kasuje dziesiątki par naraz, stąd jedna transakcja.
func (r *repozytoriumTlumaczen) UsunWpisyPamieci(ctx context.Context, kody []string) (int, error) {
	if len(kody) == 0 {
		return 0, nil
	}
	zeszlo := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := transakcja.PrepareContext(ctx,
			`DELETE FROM pamiec_tlumaczen WHERE identyfikator_zewnetrzny = ?`)
		if err != nil {
			return fmt.Errorf("dane: nie można przygotować usunięcia par pamięci: %w", err)
		}
		defer polecenie.Close()
		for _, kod := range kody {
			wynik, err := polecenie.ExecContext(ctx, kod)
			if err != nil {
				return fmt.Errorf("dane: nie można usunąć pary pamięci %q: %w", kod, err)
			}
			if liczba, err := wynik.RowsAffected(); err == nil {
				zeszlo += int(liczba)
			}
		}
		return nil
	})
	return zeszlo, err
}

func (r *repozytoriumTlumaczen) PolitykaPamieci(ctx context.Context,
	oknoID int64) (PolitykaPamieciOkna, error) {

	var polityka PolitykaPamieciOkna
	var kontekst, wstepne int64
	err := r.db.QueryRowContext(ctx,
		`SELECT okno_id, zasieg, prog, dopasowanie_kontekstu, wstepne_tlumaczenie, zaktualizowano
		   FROM polityka_pamieci_okna WHERE okno_id = ? AND `+warunekOknaPolitykiPamieci, oknoID, KontoOperatora(ctx)).
		Scan(&polityka.OknoID, &polityka.Zasieg, &polityka.Prog, &kontekst, &wstepne,
			&polityka.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return PolitykaPamieciOkna{}, ErrBrakWiersza
	}
	if err != nil {
		return PolitykaPamieciOkna{}, fmt.Errorf("dane: nieczytelna polityka pamięci okna %d: %w", oknoID, err)
	}
	polityka.DopasowanieKontekstu = kontekst == 1
	polityka.WstepneTlumaczenie = wstepne == 1
	return polityka, nil
}

func (r *repozytoriumTlumaczen) ZapiszPolitykePamieci(ctx context.Context,
	polityka PolitykaPamieciOkna) (PolitykaPamieciOkna, error) {

	teraz := time.Now().UnixMilli()
	wynik, err := r.db.ExecContext(ctx, `INSERT INTO polityka_pamieci_okna
		(okno_id, zasieg, prog, dopasowanie_kontekstu, wstepne_tlumaczenie, zaktualizowano)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(okno_id) DO UPDATE SET
			zasieg = excluded.zasieg,
			prog = excluded.prog,
			dopasowanie_kontekstu = excluded.dopasowanie_kontekstu,
			wstepne_tlumaczenie = excluded.wstepne_tlumaczenie,
			zaktualizowano = excluded.zaktualizowano
		WHERE `+warunekOknaPolitykiPamieci,
		polityka.OknoID, polityka.Zasieg, polityka.Prog,
		wartoscLogicznaDoKolumny(polityka.DopasowanieKontekstu),
		wartoscLogicznaDoKolumny(polityka.WstepneTlumaczenie), teraz, KontoOperatora(ctx))
	if err != nil {
		return PolitykaPamieciOkna{}, fmt.Errorf("dane: nie można zapisać polityki pamięci okna %d: %w",
			polityka.OknoID, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "polityka pamięci okna",
		fmt.Sprintf("%d", polityka.OknoID)); err != nil {
		return PolitykaPamieciOkna{}, err
	}
	return r.PolitykaPamieci(ctx, polityka.OknoID)
}

// Kolumny logiczne modułu Translate są INTEGER z CHECK na 0 albo 1 (migracja 160).
func wartoscLogicznaDoKolumny(wartosc bool) int64 {
	if wartosc {
		return 1
	}
	return 0
}
