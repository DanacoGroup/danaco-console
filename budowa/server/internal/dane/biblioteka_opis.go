// Plik prowadzi opis zasobu w schemacie Dublin Core oraz definicje pól niestandardowych schematu metadanych; dwa
// byty w jednym pliku, bo definicja mówi, jakie pole wolno wypełnić, a opis niesie wartość — dwie strony jednego pytania.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// OpisZasobuBiblioteki to wiersz tabeli `opis_zasobu_biblioteki` — piętnaście
// pól Dublin Core i mapa pól niestandardowych w zapisie JSON.
type OpisZasobuBiblioteki struct {
	Tytul              *string
	Tworca             *string
	Temat              *string
	Opis               *string
	Wydawca            *string
	Wspoltworca        *string
	DataZasobu         *string
	Rodzaj             *string
	Format             *string
	Identyfikator      *string
	Zrodlo             *string
	Jezyk              *string
	Powiazanie         *string
	Zakres             *string
	Prawa              *string
	PolaNiestandardowe *string
	Zaktualizowano     string
}

// PoleSchematuBiblioteki to wiersz tabeli `pole_schematu_biblioteki` niosący definicję jednego pola metadanych.
type PoleSchematuBiblioteki struct {
	Kod         string
	Etykieta    string
	Rodzaj      string
	Wymagane    bool
	MimeType    *string
	KolekcjaKod *string
	// Opcje niesie słownik dopuszczalnych wartości w zapisie JSON, wypełniony wyłącznie dla rodzaju lista.
	Opcje     *string
	Utworzono string
}

const (
	kolumnyOpisuZasobuBiblioteki = `tytul, tworca, temat, opis, wydawca, wspoltworca, data_zasobu,
	                      rodzaj, format, identyfikator, zrodlo, jezyk, powiazanie, zakres,
	                      prawa, pola_niestandardowe, zaktualizowano`

	pobierzOpisZasobuBiblioteki = `SELECT ` + kolumnyOpisuZasobuBiblioteki + ` FROM opis_zasobu_biblioteki
	                     WHERE plik_id = ?`

	zapiszOpisZasobuBiblioteki = `INSERT INTO opis_zasobu_biblioteki
	                    (plik_id, tytul, tworca, temat, opis, wydawca, wspoltworca, data_zasobu,
	                     rodzaj, format, identyfikator, zrodlo, jezyk, powiazanie, zakres, prawa,
	                     pola_niestandardowe)
	                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                    ON CONFLICT(plik_id) DO UPDATE SET
	                        tytul = excluded.tytul, tworca = excluded.tworca,
	                        temat = excluded.temat, opis = excluded.opis,
	                        wydawca = excluded.wydawca, wspoltworca = excluded.wspoltworca,
	                        data_zasobu = excluded.data_zasobu, rodzaj = excluded.rodzaj,
	                        format = excluded.format, identyfikator = excluded.identyfikator,
	                        zrodlo = excluded.zrodlo, jezyk = excluded.jezyk,
	                        powiazanie = excluded.powiazanie, zakres = excluded.zakres,
	                        prawa = excluded.prawa,
	                        pola_niestandardowe = excluded.pola_niestandardowe,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	kolumnyPolaSchematuBiblioteki = `kod, etykieta, rodzaj, wymagane, mime_type, kolekcja_kod, opcje, utworzono`

	zapiszPoleSchematuBiblioteki = `INSERT INTO pole_schematu_biblioteki
	                                (kod, etykieta, rodzaj, wymagane, mime_type, kolekcja_kod, opcje)
	                                VALUES (?, ?, ?, ?, ?, ?, ?)
	                                ON CONFLICT(kod) DO UPDATE SET
	                                    etykieta = excluded.etykieta,
	                                    rodzaj = excluded.rodzaj,
	                                    wymagane = excluded.wymagane,
	                                    mime_type = excluded.mime_type,
	                                    kolekcja_kod = excluded.kolekcja_kod,
	                                    opcje = excluded.opcje`

	pobierzPoleSchematuBiblioteki = `SELECT ` + kolumnyPolaSchematuBiblioteki + ` FROM pole_schematu_biblioteki WHERE kod = ?`

	usunPoleSchematuBiblioteki = `DELETE FROM pole_schematu_biblioteki WHERE kod = ?`

	// Wartość pola niestandardowego leży w zapisie JSON kolumny
	// `pola_niestandardowe`, więc zliczenie idzie funkcją `json_extract`
	// SQLite — nie po tekście, żeby kod pola będący fragmentem innego kodu nie
	// dawał trafienia.
	policzZasobyZPolemBiblioteki = `SELECT COUNT(*) FROM opis_zasobu_biblioteki
	                      WHERE pola_niestandardowe IS NOT NULL
	                        AND json_valid(pola_niestandardowe)
	                        AND json_extract(pola_niestandardowe, '$.' || ?) IS NOT NULL`
)

// Opis zwraca opis zasobu. Zasób bez opisu oddaje strukturę pustą — brak opisu
// jest stanem, nie usterką.
func (r *repozytoriumBiblioteki) Opis(ctx context.Context, plikID int64) (OpisZasobuBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOpisZasobuBiblioteki)
	if err != nil {
		return OpisZasobuBiblioteki{}, err
	}
	opis, err := odczytajOpisZasobuBiblioteki(polecenie.QueryRowContext(ctx, plikID))
	if errors.Is(err, sql.ErrNoRows) {
		return OpisZasobuBiblioteki{}, nil
	}
	if err != nil {
		return OpisZasobuBiblioteki{}, fmt.Errorf("dane: nieczytelny opis zasobu %d: %w", plikID, err)
	}
	return opis, nil
}

// ZapiszOpis utrwala opis zasobu w całości — wołający przysyła stan po zmianie,
// a nie różnicę. Scalanie z opisem zastanym należy do rdzenia, bo to on zna
// znaczenie żądania (`replace` kontraktu).
func (r *repozytoriumBiblioteki) ZapiszOpis(ctx context.Context, plikID int64,
	opis OpisZasobuBiblioteki) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszOpisZasobuBiblioteki)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, plikID,
		tekstDoKolumny(opis.Tytul), tekstDoKolumny(opis.Tworca), tekstDoKolumny(opis.Temat),
		tekstDoKolumny(opis.Opis), tekstDoKolumny(opis.Wydawca), tekstDoKolumny(opis.Wspoltworca),
		tekstDoKolumny(opis.DataZasobu), tekstDoKolumny(opis.Rodzaj), tekstDoKolumny(opis.Format),
		tekstDoKolumny(opis.Identyfikator), tekstDoKolumny(opis.Zrodlo), tekstDoKolumny(opis.Jezyk),
		tekstDoKolumny(opis.Powiazanie), tekstDoKolumny(opis.Zakres), tekstDoKolumny(opis.Prawa),
		tekstDoKolumny(opis.PolaNiestandardowe))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać opisu zasobu %d: %w", plikID, err)
	}
	return nil
}

// PolaSchematu zwraca definicje pól, zawężone do rodzaju treści albo kolekcji; zawężenie jest miękkie.
func (r *repozytoriumBiblioteki) PolaSchematu(ctx context.Context,
	mimeType, kolekcjaKod *string) ([]PoleSchematuBiblioteki, error) {

	warunki := []string{"1 = 1"}
	argumenty := []any{}
	if mimeType != nil && *mimeType != "" {
		warunki = append(warunki, "(mime_type IS NULL OR mime_type = '' OR mime_type = ?)")
		argumenty = append(argumenty, *mimeType)
	}
	if kolekcjaKod != nil && *kolekcjaKod != "" {
		warunki = append(warunki, "(kolekcja_kod IS NULL OR kolekcja_kod = '' OR kolekcja_kod = ?)")
		argumenty = append(argumenty, *kolekcjaKod)
	}
	zapytanie := `SELECT ` + kolumnyPolaSchematuBiblioteki + ` FROM pole_schematu_biblioteki
	              WHERE ` + strings.Join(warunki, " AND ") + ` ORDER BY etykieta, kod`

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pól schematu biblioteki: %w", err)
	}
	defer wiersze.Close()

	lista := []PoleSchematuBiblioteki{}
	for wiersze.Next() {
		pole, err := odczytajPoleSchematuBiblioteki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pola schematu biblioteki: %w", err)
		}
		lista = append(lista, pole)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pól schematu biblioteki: %w", err)
	}
	return lista, nil
}

// ZapiszPoleSchematu zakłada definicję pola albo zmienia zastaną — pole o kodzie
// już istniejącym jest zmieniane, nie dublowane.
func (r *repozytoriumBiblioteki) ZapiszPoleSchematu(ctx context.Context,
	pole PoleSchematuBiblioteki) (PoleSchematuBiblioteki, error) {

	if strings.TrimSpace(pole.Kod) == "" {
		return PoleSchematuBiblioteki{}, fmt.Errorf("dane: pole schematu biblioteki bez kodu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPoleSchematuBiblioteki)
	if err != nil {
		return PoleSchematuBiblioteki{}, err
	}
	_, err = polecenie.ExecContext(ctx, pole.Kod, pole.Etykieta, pole.Rodzaj,
		liczbaLogiczna(pole.Wymagane), tekstDoKolumny(pole.MimeType),
		tekstDoKolumny(pole.KolekcjaKod), tekstDoKolumny(pole.Opcje))
	if err != nil {
		return PoleSchematuBiblioteki{}, fmt.Errorf("dane: nie można zapisać pola schematu %q: %w",
			pole.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzPoleSchematuBiblioteki)
	if err != nil {
		return PoleSchematuBiblioteki{}, err
	}
	zapisane, err := odczytajPoleSchematuBiblioteki(odczyt.QueryRowContext(ctx, pole.Kod))
	if err != nil {
		return PoleSchematuBiblioteki{}, fmt.Errorf("dane: nieczytelne pole schematu %q: %w", pole.Kod, err)
	}
	return zapisane, nil
}

// UsunPoleSchematu zdejmuje definicję pola. Wartości zapisane przy zasobach
// zostają nietknięte — kontrakt mówi wprost, że wracają, gdy pole zostanie
// założone ponownie.
func (r *repozytoriumBiblioteki) UsunPoleSchematu(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunPoleSchematuBiblioteki)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć pola schematu %q: %w", kod, err)
	}
	zdjete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć zdjętych pól schematu: %w", err)
	}
	return zdjete > 0, nil
}

// ZasobyZPolem liczy zasoby, przy których pole niestandardowe ma już wypełnioną wartość w bazie danych.
func (r *repozytoriumBiblioteki) ZasobyZPolem(ctx context.Context, kodPola string) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, policzZasobyZPolemBiblioteki)
	if err != nil {
		return 0, err
	}
	var liczba int
	if err := polecenie.QueryRowContext(ctx, kodPola).Scan(&liczba); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć zasobów z polem %q: %w", kodPola, err)
	}
	return liczba, nil
}

// odczytajOpisZasobuBiblioteki składa opis wprost z jednego wiersza wyniku zapytania SQL do bazy danych.
func odczytajOpisZasobuBiblioteki(wiersz skaner) (OpisZasobuBiblioteki, error) {
	var opis OpisZasobuBiblioteki
	var tytul, tworca, temat, tresc, wydawca, wspoltworca, data sql.NullString
	var rodzaj, format, identyfikator, zrodlo, jezyk, powiazanie, zakres, prawa, pola sql.NullString
	err := wiersz.Scan(&tytul, &tworca, &temat, &tresc, &wydawca, &wspoltworca, &data,
		&rodzaj, &format, &identyfikator, &zrodlo, &jezyk, &powiazanie, &zakres, &prawa, &pola,
		&opis.Zaktualizowano)
	if err != nil {
		return OpisZasobuBiblioteki{}, err
	}
	opis.Tytul, opis.Tworca, opis.Temat = tekstZKolumny(tytul), tekstZKolumny(tworca), tekstZKolumny(temat)
	opis.Opis, opis.Wydawca = tekstZKolumny(tresc), tekstZKolumny(wydawca)
	opis.Wspoltworca, opis.DataZasobu = tekstZKolumny(wspoltworca), tekstZKolumny(data)
	opis.Rodzaj, opis.Format = tekstZKolumny(rodzaj), tekstZKolumny(format)
	opis.Identyfikator, opis.Zrodlo = tekstZKolumny(identyfikator), tekstZKolumny(zrodlo)
	opis.Jezyk, opis.Powiazanie = tekstZKolumny(jezyk), tekstZKolumny(powiazanie)
	opis.Zakres, opis.Prawa = tekstZKolumny(zakres), tekstZKolumny(prawa)
	opis.PolaNiestandardowe = tekstZKolumny(pola)
	return opis, nil
}

// odczytajPoleSchematuBiblioteki składa definicję pola wprost z jednego wiersza wyniku zapytania SQL do bazy.
func odczytajPoleSchematuBiblioteki(wiersz skaner) (PoleSchematuBiblioteki, error) {
	var pole PoleSchematuBiblioteki
	var wymagane int
	var mimeType, kolekcja, opcje sql.NullString
	err := wiersz.Scan(&pole.Kod, &pole.Etykieta, &pole.Rodzaj, &wymagane,
		&mimeType, &kolekcja, &opcje, &pole.Utworzono)
	if err != nil {
		return PoleSchematuBiblioteki{}, err
	}
	pole.Wymagane = wymagane == 1
	pole.MimeType, pole.KolekcjaKod, pole.Opcje = tekstZKolumny(mimeType), tekstZKolumny(kolekcja),
		tekstZKolumny(opcje)
	return pole, nil
}
