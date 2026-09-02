// Odpowiedzialność pliku: źródła zebrane w toku przeglądania (tabela zrodlo_przegladania), szuflada Sources okna.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ZrodloPrzegladania struct {
	ID                  int64
	Kod                 string
	Okno                string
	Url                 string
	Tytul               *string
	MigawkaZewnetrznaID *string
	Kluczowe            bool
	// Grupa niesie zestaw tematyczny źródła; wskaźnik pusty znaczy poza zestawami, nie zestaw bez nazwy.
	Grupa     *string
	Utworzono string
}

const (
	kolumnyZrodlaPrzegladania = `id, identyfikator_zewnetrzny, okno, url, tytul,
	                             migawka_zewnetrzna_id, kluczowe, grupa, utworzono`

	// UNIQUE na identyfikatorze obejmuje całą tabelę: człon DO UPDATE bez zawężenia nadpisałby źródło konta cudzego.
	zapiszZrodloPrzegladania = `INSERT INTO zrodlo_przegladania
	                            (identyfikator_zewnetrzny, okno, url, tytul,
	                             migawka_zewnetrzna_id, kluczowe, grupa, konto_id)
	                            VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                            ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                okno = excluded.okno,
	                                url = excluded.url,
	                                tytul = excluded.tytul,
	                                migawka_zewnetrzna_id = excluded.migawka_zewnetrzna_id,
	                                kluczowe = excluded.kluczowe,
	                                grupa = excluded.grupa
	                            WHERE ` + WarunekKonta

	pobierzZrodloPrzegladania = `SELECT ` + kolumnyZrodlaPrzegladania + `
	                             FROM zrodlo_przegladania
	                             WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaZrodelPrzegladania = `SELECT ` + kolumnyZrodlaPrzegladania + `
	                           FROM zrodlo_przegladania
	                           WHERE okno = ? AND (? = 0 OR kluczowe = 1)
	                             AND (? = '' OR grupa = ?)
	                             AND (? = '' OR url LIKE ? OR IFNULL(tytul,'') LIKE ?)
	                             AND ` + WarunekKonta + `
	                           ORDER BY utworzono DESC, id DESC LIMIT ?`
)

// FiltrZrodelPrzegladania niesie pola żądania `browser.source.list`.
type FiltrZrodelPrzegladania struct {
	Okno          string
	TylkoKluczowe bool
	Zestaw        string
	Szukaj        string
	// Limit 0 lub ujemny znaczy wykaz pełny, nie wykaz pusty.
	Limit int
}

func (r *repozytoriumPrzegladania) ZapiszZrodlo(ctx context.Context, zrodlo ZrodloPrzegladania) (ZrodloPrzegladania, error) {
	if zrodlo.Kod == "" || zrodlo.Okno == "" || zrodlo.Url == "" {
		return ZrodloPrzegladania{}, fmt.Errorf("dane: źródło przeglądania bez identyfikatora, okna albo adresu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZrodloPrzegladania)
	if err != nil {
		return ZrodloPrzegladania{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, zrodlo.Kod, zrodlo.Okno, zrodlo.Url,
		tekstDoKolumny(zrodlo.Tytul), tekstDoKolumny(zrodlo.MigawkaZewnetrznaID),
		liczbaLogiczna(zrodlo.Kluczowe), tekstDoKolumny(zrodlo.Grupa),
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return ZrodloPrzegladania{}, fmt.Errorf("dane: nie można zapisać źródła przeglądania %q: %w", zrodlo.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "źródło przeglądania", zrodlo.Kod); err != nil {
		return ZrodloPrzegladania{}, err
	}
	return r.jednoZrodlo(ctx, zrodlo.Kod)
}

func (r *repozytoriumPrzegladania) jednoZrodlo(ctx context.Context, kod string) (ZrodloPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZrodloPrzegladania)
	if err != nil {
		return ZrodloPrzegladania{}, err
	}
	zrodlo, err := odczytajZrodloPrzegladania(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return ZrodloPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return ZrodloPrzegladania{}, fmt.Errorf("dane: nieczytelne źródło przeglądania %q: %w", kod, err)
	}
	return zrodlo, nil
}

func (r *repozytoriumPrzegladania) Zrodla(ctx context.Context,
	filtr FiltrZrodelPrzegladania) ([]ZrodloPrzegladania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZrodelPrzegladania)
	if err != nil {
		return nil, err
	}
	okno := filtr.Okno
	wzorzec := ""
	if filtr.Szukaj != "" {
		wzorzec = "%" + filtr.Szukaj + "%"
	}
	wiersze, err := polecenie.QueryContext(ctx, okno,
		liczbaLogiczna(filtr.TylkoKluczowe), filtr.Zestaw, filtr.Zestaw,
		filtr.Szukaj, wzorzec, wzorzec, KontoOperatora(ctx), granicaWykazu(filtr.Limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać źródeł przeglądania okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []ZrodloPrzegladania{}
	for wiersze.Next() {
		zrodlo, err := odczytajZrodloPrzegladania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz źródła przeglądania: %w", err)
		}
		lista = append(lista, zrodlo)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt źródeł przeglądania: %w", err)
	}
	return lista, nil
}

func odczytajZrodloPrzegladania(wiersz skaner) (ZrodloPrzegladania, error) {
	var zrodlo ZrodloPrzegladania
	var tytul, migawkaZewnetrznaID, grupa sql.NullString
	var kluczowe int
	err := wiersz.Scan(&zrodlo.ID, &zrodlo.Kod, &zrodlo.Okno, &zrodlo.Url, &tytul,
		&migawkaZewnetrznaID, &kluczowe, &grupa, &zrodlo.Utworzono)
	if err != nil {
		return ZrodloPrzegladania{}, err
	}
	zrodlo.Tytul = tekstZKolumny(tytul)
	zrodlo.MigawkaZewnetrznaID = tekstZKolumny(migawkaZewnetrznaID)
	zrodlo.Kluczowe = kluczowe != 0
	zrodlo.Grupa = tekstZKolumny(grupa)
	return zrodlo, nil
}
