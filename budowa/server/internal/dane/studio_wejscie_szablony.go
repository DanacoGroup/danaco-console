// Warsztat szablonów Studia: wiersz tabeli `szablon_studio` w pełnym kształcie
// z migracji 367, z granicą konta z migracji 484.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type SzablonWarsztatuStudia struct {
	ID                  int64
	Kod                 string
	Nazwa               string
	Opis                *string
	Format              string
	Tresc               string
	PolaJSON            *string
	PostacJSON          *string
	Kategoria           *string
	MiniaturaZasobKod   *string
	ZrodloPliku         *string
	DokumentZrodlowyKod *string
	Fabryczny           bool
	Utworzono           string
	Zaktualizowano      *string
}

type WarsztatSzablonowStudia interface {
	ZapiszSzablonWarsztatu(ctx context.Context,
		szablon SzablonWarsztatuStudia) (SzablonWarsztatuStudia, error)
	SzablonWarsztatu(ctx context.Context, kod string) (SzablonWarsztatuStudia, error)
	SzablonyWarsztatu(ctx context.Context, kategoria string) ([]SzablonWarsztatuStudia, error)
	UsunSzablonWlasny(ctx context.Context, kod string) (bool, error)
}

const (
	wejscieKolumnySzablonu = `id, identyfikator_zewnetrzny, nazwa, opis, format, tresc,
	                          pola_json, postac_json, kategoria, miniatura_zasob_kod,
	                          zrodlo_pliku, dokument_zrodlowy_kod, fabryczny,
	                          utworzono, zaktualizowano`

	// Klucz `identyfikator_zewnetrzny` jest jeden na całą tabelę, więc kod szablonu
	// cudzego konta trafia w konflikt; warunek przy DO UPDATE zostawia wiersz nietknięty.
	wejscieZapiszSzablon = `INSERT INTO szablon_studio
	                        (identyfikator_zewnetrzny, nazwa, opis, format, tresc, pola_json,
	                         postac_json, kategoria, miniatura_zasob_kod, zrodlo_pliku,
	                         dokument_zrodlowy_kod, fabryczny, zaktualizowano, konto_id)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	                                strftime('%Y-%m-%dT%H:%M:%fZ','now'), ` + WskazanieKonta + `)
	                        ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                            nazwa = excluded.nazwa,
	                            opis = excluded.opis,
	                            format = excluded.format,
	                            tresc = excluded.tresc,
	                            pola_json = excluded.pola_json,
	                            postac_json = excluded.postac_json,
	                            kategoria = excluded.kategoria,
	                            miniatura_zasob_kod = excluded.miniatura_zasob_kod,
	                            zrodlo_pliku = excluded.zrodlo_pliku,
	                            dokument_zrodlowy_kod = excluded.dokument_zrodlowy_kod,
	                            zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                        WHERE ` + WarunekKonta

	// Szablon fabryczny zakłada migracja 131 bez wskazania konta i jest wyposażeniem
	// instalacji widocznym dla każdego konta; zawężenie obejmuje szablony własne.
	wejsciePobierzSzablon = `SELECT ` + wejscieKolumnySzablonu + ` FROM szablon_studio
	                         WHERE identyfikator_zewnetrzny = ?
	                           AND (fabryczny = 1 OR ` + WarunekKonta + `)`

	wejscieListaSzablonow = `SELECT ` + wejscieKolumnySzablonu + ` FROM szablon_studio
	                         WHERE (? = '' OR kategoria = ?)
	                           AND (fabryczny = 1 OR ` + WarunekKonta + `)
	                         ORDER BY fabryczny DESC, nazwa`

	// Warunek `fabryczny = 0` stoi w zapytaniu: szablonu fabrycznego nie da się
	// odtworzyć bez migracji.
	wejscieUsunSzablon = `DELETE FROM szablon_studio
	                      WHERE identyfikator_zewnetrzny = ? AND fabryczny = 0
	                        AND ` + WarunekKonta
)

func (r *repozytoriumStudia) ZapiszSzablonWarsztatu(ctx context.Context,
	szablon SzablonWarsztatuStudia) (SzablonWarsztatuStudia, error) {

	if szablon.Kod == "" {
		return SzablonWarsztatuStudia{}, fmt.Errorf("dane: szablon studio bez identyfikatora")
	}
	if szablon.Nazwa == "" {
		return SzablonWarsztatuStudia{}, fmt.Errorf("dane: szablon studio %q bez nazwy", szablon.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wejscieZapiszSzablon)
	if err != nil {
		return SzablonWarsztatuStudia{}, err
	}
	konto := KontoOperatora(ctx)
	wynik, err := polecenie.ExecContext(ctx, szablon.Kod, szablon.Nazwa, tekstDoKolumny(szablon.Opis),
		szablon.Format, szablon.Tresc, tekstDoKolumny(szablon.PolaJSON),
		tekstDoKolumny(szablon.PostacJSON), tekstDoKolumny(szablon.Kategoria),
		tekstDoKolumny(szablon.MiniaturaZasobKod), tekstDoKolumny(szablon.ZrodloPliku),
		tekstDoKolumny(szablon.DokumentZrodlowyKod), liczbaLogiczna(szablon.Fabryczny),
		konto, konto)
	if err != nil {
		return SzablonWarsztatuStudia{}, fmt.Errorf("dane: nie można zapisać szablonu %q: %w",
			szablon.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "szablon studio", szablon.Kod); err != nil {
		return SzablonWarsztatuStudia{}, err
	}
	return r.SzablonWarsztatu(ctx, szablon.Kod)
}

func (r *repozytoriumStudia) SzablonWarsztatu(ctx context.Context,
	kod string) (SzablonWarsztatuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wejsciePobierzSzablon)
	if err != nil {
		return SzablonWarsztatuStudia{}, err
	}
	szablon, err := wejscieOdczytajSzablon(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return SzablonWarsztatuStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return SzablonWarsztatuStudia{}, fmt.Errorf("dane: nieczytelny szablon %q: %w", kod, err)
	}
	return szablon, nil
}

func (r *repozytoriumStudia) SzablonyWarsztatu(ctx context.Context,
	kategoria string) ([]SzablonWarsztatuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wejscieListaSzablonow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kategoria, kategoria, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać szablonów studio: %w", err)
	}
	defer wiersze.Close()

	lista := []SzablonWarsztatuStudia{}
	for wiersze.Next() {
		szablon, err := wejscieOdczytajSzablon(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz szablonu studio: %w", err)
		}
		lista = append(lista, szablon)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt szablonów studio: %w", err)
	}
	return lista, nil
}

func (r *repozytoriumStudia) UsunSzablonWlasny(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wejscieUsunSzablon)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć szablonu %q: %w", kod, err)
	}
	ile, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie wiadomo, czy szablon %q zszedł: %w", kod, err)
	}
	return ile > 0, nil
}

func wejscieOdczytajSzablon(wiersz skaner) (SzablonWarsztatuStudia, error) {
	var szablon SzablonWarsztatuStudia
	var opis, pola, postac, kategoria, miniatura, zrodlo, dokument, zmieniono sql.NullString
	var fabryczny int
	err := wiersz.Scan(&szablon.ID, &szablon.Kod, &szablon.Nazwa, &opis, &szablon.Format,
		&szablon.Tresc, &pola, &postac, &kategoria, &miniatura, &zrodlo, &dokument,
		&fabryczny, &szablon.Utworzono, &zmieniono)
	if err != nil {
		return SzablonWarsztatuStudia{}, err
	}
	szablon.Opis = tekstZKolumny(opis)
	szablon.PolaJSON = tekstZKolumny(pola)
	szablon.PostacJSON = tekstZKolumny(postac)
	szablon.Kategoria = tekstZKolumny(kategoria)
	szablon.MiniaturaZasobKod = tekstZKolumny(miniatura)
	szablon.ZrodloPliku = tekstZKolumny(zrodlo)
	szablon.DokumentZrodlowyKod = tekstZKolumny(dokument)
	szablon.Zaktualizowano = tekstZKolumny(zmieniono)
	szablon.Fabryczny = fabryczny == 1
	return szablon, nil
}
