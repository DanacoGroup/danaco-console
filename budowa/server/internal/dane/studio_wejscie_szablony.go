// Repozytorium ogłasza SzablonWarsztatuStudia jako wiersz tabeli `szablon_studio`
// wzbogacony migracją 367 o kategorię, postać wzorcową, miniaturę, źródło pliku,
// dokument źródłowy, czas zmiany oraz usunięcie szablonu własnego.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SzablonWarsztatuStudia to wiersz tabeli `szablon_studio` widziany w pełni —
// wraz z kolumnami warsztatu szablonów dobudowanymi migracją 367.
type SzablonWarsztatuStudia struct {
	ID     int64
	Kod    string
	Nazwa  string
	Opis   *string
	Format string
	Tresc  string
	// PolaJSON niesie wykaz pól do wypełnienia w kształcie
	// `shared.StudioTemplateFieldSpec[]`.
	PolaJSON *string
	// PostacJSON niesie postać wzorcową w tym samym kształcie, co postać
	// dokumentu, bez przekładu.
	PostacJSON          *string
	Kategoria           *string
	MiniaturaZasobKod   *string
	ZrodloPliku         *string
	DokumentZrodlowyKod *string
	Fabryczny           bool
	Utworzono           string
	Zaktualizowano      *string
}

// WarsztatSzablonowStudia jest kontraktem tej warstwy. Rdzeń bierze go
// rzutowaniem, więc repozytorium bez tych kolumn nazywa brak wprost, zamiast
// wywracać montaż.
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

	// Zapis jest upsertem po kluczu zewnętrznym, tym samym, którym jedzie zapis
	// z `studio_katalogi.go`, do tabeli `szablon_studio`.
	wejscieZapiszSzablon = `INSERT INTO szablon_studio
	                        (identyfikator_zewnetrzny, nazwa, opis, format, tresc, pola_json,
	                         postac_json, kategoria, miniatura_zasob_kod, zrodlo_pliku,
	                         dokument_zrodlowy_kod, fabryczny, zaktualizowano)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	                                strftime('%Y-%m-%dT%H:%M:%fZ','now'))
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
	                            zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	wejsciePobierzSzablon = `SELECT ` + wejscieKolumnySzablonu + ` FROM szablon_studio
	                         WHERE identyfikator_zewnetrzny = ?`

	wejscieListaSzablonow = `SELECT ` + wejscieKolumnySzablonu + ` FROM szablon_studio
	                         WHERE (? = '' OR kategoria = ?)
	                         ORDER BY fabryczny DESC, nazwa`

	// Usunięcie obejmuje WYŁĄCZNIE szablon własny. Warunek stoi w zapytaniu,
	// nie tylko w rdzeniu: szablon fabryczny usunięty inną drogą zabrałby
	// Operatorowi wykaz, którego nie da się odtworzyć bez migracji.
	wejscieUsunSzablon = `DELETE FROM szablon_studio
	                      WHERE identyfikator_zewnetrzny = ? AND fabryczny = 0`
)

// ZapiszSzablonWarsztatu zakłada szablon pisma albo nadpisuje zastany wraz
// z postacią wzorcową i polami do wypełnienia.
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
	_, err = polecenie.ExecContext(ctx, szablon.Kod, szablon.Nazwa, tekstDoKolumny(szablon.Opis),
		szablon.Format, szablon.Tresc, tekstDoKolumny(szablon.PolaJSON),
		tekstDoKolumny(szablon.PostacJSON), tekstDoKolumny(szablon.Kategoria),
		tekstDoKolumny(szablon.MiniaturaZasobKod), tekstDoKolumny(szablon.ZrodloPliku),
		tekstDoKolumny(szablon.DokumentZrodlowyKod), liczbaLogiczna(szablon.Fabryczny))
	if err != nil {
		return SzablonWarsztatuStudia{}, fmt.Errorf("dane: nie można zapisać szablonu %q: %w",
			szablon.Kod, err)
	}
	return r.SzablonWarsztatu(ctx, szablon.Kod)
}

// SzablonWarsztatu oddaje szablon warsztatu wraz z postacią wzorcową i polami
// do wypełnienia z tabeli `szablon_studio`.
func (r *repozytoriumStudia) SzablonWarsztatu(ctx context.Context,
	kod string) (SzablonWarsztatuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wejsciePobierzSzablon)
	if err != nil {
		return SzablonWarsztatuStudia{}, err
	}
	szablon, err := wejscieOdczytajSzablon(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SzablonWarsztatuStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return SzablonWarsztatuStudia{}, fmt.Errorf("dane: nieczytelny szablon %q: %w", kod, err)
	}
	return szablon, nil
}

// SzablonyWarsztatu oddaje wykaz szablonów, fabryczne na początku. Kategoria
// pusta znaczy „wszystkie" — zawężenie jest zawężeniem Operatora, nie warunkiem.
func (r *repozytoriumStudia) SzablonyWarsztatu(ctx context.Context,
	kategoria string) ([]SzablonWarsztatuStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wejscieListaSzablonow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kategoria, kategoria)
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

// UsunSzablonWlasny usuwa szablon własny i mówi, czy wiersz zszedł. Fałsz przy
// istniejącym szablonie znaczy szablon fabryczny — i tak to nazywa rdzeń.
func (r *repozytoriumStudia) UsunSzablonWlasny(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wejscieUsunSzablon)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
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
