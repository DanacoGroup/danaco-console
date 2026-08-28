// Plik prowadzi obszar Research: deklarację całego interfejsu RepozytoriumBadan oraz obsługę źródeł badania;
// ustalenia leżą w badania_ustalenia.go, raport w badania_raport.go, a interfejs deklaruje wyłącznie ten plik, w całości.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// ZrodloBadania to wiersz tabeli `zrodlo_badania`; kod jest identyfikatorem, którym źródło wychodzi kontraktem, a plik biblioteki nie niesie więzu obcego.
type ZrodloBadania struct {
	ID               int64
	Kod              string
	Okno             string
	Tytul            string
	Rodzaj           shared.ResearchSourceKind
	Adres            *string
	Pochodzenie      *string
	Wiarygodnosc     shared.ResearchCredibility
	PlikBibliotekiID *string
	PozyskanoO       string
}

// RepozytoriumBadan jest kontraktem obszaru Research: źródła, ustalenia, raport, eksport i przestrzeń badania.
type RepozytoriumBadan interface {
	// źródła
	ZapiszZrodlo(ctx context.Context, zrodlo ZrodloBadania) (ZrodloBadania, error)
	Zrodlo(ctx context.Context, kod string) (ZrodloBadania, error)
	Zrodla(ctx context.Context, okno string) ([]ZrodloBadania, error)

	// ustalenia
	ZapiszUstalenie(ctx context.Context, ustalenie UstalenieBadania,
		kodyZrodel []string) (UstalenieBadania, error)
	Ustalenie(ctx context.Context, kod string) (UstalenieBadania, error)
	Ustalenia(ctx context.Context, okno string) ([]UstalenieBadania, error)
	ZrodlaUstalenia(ctx context.Context, ustalenieID int64) ([]ZrodloBadania, error)

	// raport, eksport, przestrzeń
	ZapiszRaport(ctx context.Context, raport RaportBadania,
		sekcje []SekcjaRaportu) (RaportBadania, error)
	Raport(ctx context.Context, kod string) (RaportBadania, error)
	Sekcje(ctx context.Context, raportID int64) ([]SekcjaRaportu, error)
	ZapiszEksport(ctx context.Context, eksport EksportRaportu) (EksportRaportu, error)
	UstawPrzestrzen(ctx context.Context, zakres string, etapy []string) (string, []string, error)
	Przestrzen(ctx context.Context) (string, []string, error)

	// Dobudowa modułu: katalogowanie źródeł, lektura, adnotacje, kodowanie, sprzeczności i odkrywanie.
	RepozytoriumBadanDobudowa
}

const (
	kolumnyZrodlaBadania = `id, identyfikator_zewnetrzny, okno, tytul, rodzaj, adres,
	                        pochodzenie, wiarygodnosc, plik_biblioteki_id, pozyskano_o`

	zapiszZrodloBadania = `INSERT INTO zrodlo_badania
	                       (identyfikator_zewnetrzny, okno, tytul, rodzaj, adres,
	                        pochodzenie, wiarygodnosc, plik_biblioteki_id)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                           tytul = excluded.tytul,
	                           rodzaj = excluded.rodzaj,
	                           adres = excluded.adres,
	                           pochodzenie = excluded.pochodzenie,
	                           wiarygodnosc = excluded.wiarygodnosc,
	                           plik_biblioteki_id = excluded.plik_biblioteki_id`

	pobierzZrodloBadania = `SELECT ` + kolumnyZrodlaBadania + ` FROM zrodlo_badania
	                        WHERE identyfikator_zewnetrzny = ?`

	pobierzZrodlaBadaniaOkna = `SELECT ` + kolumnyZrodlaBadania + ` FROM zrodlo_badania
	                            WHERE okno = ?
	                            ORDER BY pozyskano_o DESC, id DESC`
)

type repozytoriumBadan struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumBadan(z *zapytania, db *sql.DB) *repozytoriumBadan {
	return &repozytoriumBadan{zapytania: z, db: db}
}

// ZapiszZrodlo zakłada wiersz źródła albo nadpisuje zastane i zwraca stan po
// zapisie. Wiarygodność bierze wartość kontraktu (`ResearchCredibility`)
// wprost, bez tłumaczenia.
func (r *repozytoriumBadan) ZapiszZrodlo(ctx context.Context, zrodlo ZrodloBadania) (ZrodloBadania, error) {
	if zrodlo.Kod == "" {
		return ZrodloBadania{}, fmt.Errorf("dane: źródło badania bez identyfikatora")
	}
	if zrodlo.Okno == "" {
		return ZrodloBadania{}, fmt.Errorf("dane: źródło badania %q bez okna", zrodlo.Kod)
	}
	if zrodlo.Tytul == "" {
		return ZrodloBadania{}, fmt.Errorf("dane: źródło badania %q bez tytułu", zrodlo.Kod)
	}
	rodzaj := string(zrodlo.Rodzaj)
	if rodzaj == "" {
		rodzaj = string(shared.ResearchSourceKindWeb)
	}
	wiarygodnosc := string(zrodlo.Wiarygodnosc)
	if wiarygodnosc == "" {
		wiarygodnosc = "unverified"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZrodloBadania)
	if err != nil {
		return ZrodloBadania{}, err
	}
	_, err = polecenie.ExecContext(ctx, zrodlo.Kod, zrodlo.Okno, zrodlo.Tytul, rodzaj,
		tekstDoKolumny(zrodlo.Adres), tekstDoKolumny(zrodlo.Pochodzenie), wiarygodnosc,
		tekstDoKolumny(zrodlo.PlikBibliotekiID))
	if err != nil {
		return ZrodloBadania{}, fmt.Errorf("dane: nie można zapisać źródła badania %q: %w", zrodlo.Kod, err)
	}
	return r.Zrodlo(ctx, zrodlo.Kod)
}

// Zrodlo zwraca źródło o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się nie powiódł”.
func (r *repozytoriumBadan) Zrodlo(ctx context.Context, kod string) (ZrodloBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZrodloBadania)
	if err != nil {
		return ZrodloBadania{}, err
	}
	zrodlo, err := odczytajZrodloBadania(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ZrodloBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return ZrodloBadania{}, fmt.Errorf("dane: nieczytelny wiersz źródła badania %q: %w", kod, err)
	}
	return zrodlo, nil
}

// Zrodla zwraca źródła okna badania, posortowane od najświeżej pozyskanych, wprost z bazy danych repozytorium.
func (r *repozytoriumBadan) Zrodla(ctx context.Context, okno string) ([]ZrodloBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZrodlaBadaniaOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać źródeł badania okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []ZrodloBadania{}
	for wiersze.Next() {
		zrodlo, err := odczytajZrodloBadania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz źródła badania: %w", err)
		}
		lista = append(lista, zrodlo)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt źródeł badania okna %q: %w", okno, err)
	}
	return lista, nil
}

// odczytajZrodloBadania składa strukturę źródła badania wprost z jednego wiersza wyniku zapytania do bazy.
func odczytajZrodloBadania(wiersz skaner) (ZrodloBadania, error) {
	var zrodlo ZrodloBadania
	var rodzaj, wiarygodnosc string
	var adres, pochodzenie, plikBibliotekiID sql.NullString
	err := wiersz.Scan(&zrodlo.ID, &zrodlo.Kod, &zrodlo.Okno, &zrodlo.Tytul, &rodzaj,
		&adres, &pochodzenie, &wiarygodnosc, &plikBibliotekiID, &zrodlo.PozyskanoO)
	if err != nil {
		return ZrodloBadania{}, err
	}
	zrodlo.Rodzaj = shared.ResearchSourceKind(rodzaj)
	zrodlo.Wiarygodnosc = shared.ResearchCredibility(wiarygodnosc)
	zrodlo.Adres = tekstZKolumny(adres)
	zrodlo.Pochodzenie = tekstZKolumny(pochodzenie)
	zrodlo.PlikBibliotekiID = tekstZKolumny(plikBibliotekiID)
	return zrodlo, nil
}
