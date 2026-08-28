// Plik definiuje obiekty osadzone w dokumencie Studia, aparat dokumentu i
// pola. Kontrakt obszaru deklaruje studio_postac_dokumentu.go; ten plik
// dokłada metody.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ObiektDokumentuStudia to wiersz tabeli obiekt_dokumentu_studio. PostacJSON
// niesie cechy swoiste dla rodzaju obiektu; kolumny osobne mają tylko cechy,
// po których się pyta.
type ObiektDokumentuStudia struct {
	ID                   int64
	Kod                  string
	DokumentID           int64
	Rodzaj               string
	Zrodlo               *string
	ZasobKod             *string
	DesignWezelKod       *string
	BibliotekaPlikKod    *string
	AdresZrodla          *string
	Zakotwiczenie        string
	ZakotwiczeniePozycja int64
	Warstwa              int64
	TekstZastepczy       *string
	TekstWewnetrzny      *string
	Podpis               *string
	PostacJSON           *string
	Utworzono            string
	Zaktualizowano       string
}

// ElementAparatuStudia to wiersz tabeli element_aparatu_studio: przypis,
// odnośnik albo inny element aparatu przypięty do miejsca w treści dokumentu.
type ElementAparatuStudia struct {
	ID             int64
	Kod            string
	DokumentID     int64
	Rodzaj         string
	KotwicaOd      int64
	KotwicaDo      int64
	Numer          *string
	Etykieta       *string
	Tresc          *string
	CelKod         *string
	CelAdres       *string
	Kolejnosc      int64
	Nieswiezy      bool
	DaneJSON       *string
	Utworzono      string
	Zaktualizowano string
}

// PoleDokumentuStudia to wiersz tabeli pole_dokumentu_studio: pole
// obliczane, przypięte do miejsca w treści dokumentu i odświeżane wykazem.
type PoleDokumentuStudia struct {
	ID               int64
	Kod              string
	DokumentID       int64
	Rodzaj           string
	Kotwica          int64
	Format           *string
	Wyrazenie        *string
	NazwaWlasciwosci *string
	Wartosc          *string
	Nieswieze        bool
	Utworzono        string
	Zaktualizowano   string
}

const (
	postacKolumnyObiektu = `id, identyfikator_zewnetrzny, dokument_id, rodzaj, zrodlo,
	                        zasob_kod, design_wezel_kod, biblioteka_plik_kod, adres_zrodla,
	                        zakotwiczenie, zakotwiczenie_pozycja, warstwa, tekst_zastepczy,
	                        tekst_wewnetrzny, podpis, postac_json, utworzono, zaktualizowano`

	postacZapiszObiekt = `INSERT INTO obiekt_dokumentu_studio
	                      (identyfikator_zewnetrzny, dokument_id, rodzaj, zrodlo, zasob_kod,
	                       design_wezel_kod, biblioteka_plik_kod, adres_zrodla, zakotwiczenie,
	                       zakotwiczenie_pozycja, warstwa, tekst_zastepczy, tekst_wewnetrzny,
	                       podpis, postac_json)
	                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          rodzaj = excluded.rodzaj,
	                          zrodlo = excluded.zrodlo,
	                          zasob_kod = excluded.zasob_kod,
	                          design_wezel_kod = excluded.design_wezel_kod,
	                          biblioteka_plik_kod = excluded.biblioteka_plik_kod,
	                          adres_zrodla = excluded.adres_zrodla,
	                          zakotwiczenie = excluded.zakotwiczenie,
	                          zakotwiczenie_pozycja = excluded.zakotwiczenie_pozycja,
	                          warstwa = excluded.warstwa,
	                          tekst_zastepczy = excluded.tekst_zastepczy,
	                          tekst_wewnetrzny = excluded.tekst_wewnetrzny,
	                          podpis = excluded.podpis,
	                          postac_json = excluded.postac_json,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	postacPobierzObiekt = `SELECT ` + postacKolumnyObiektu + ` FROM obiekt_dokumentu_studio
	                       WHERE identyfikator_zewnetrzny = ?`

	// Wykaz idzie w kolejności ZAKOTWICZENIA, nie zapisu: obiekty czyta się
	// razem z treścią, od początku dokumentu, a nie od najnowszego wstawienia.
	postacListaObiektow = `SELECT ` + postacKolumnyObiektu + ` FROM obiekt_dokumentu_studio
	                       WHERE dokument_id = ?
	                       ORDER BY zakotwiczenie_pozycja, warstwa, id`

	postacListaObiektowRodzaju = `SELECT ` + postacKolumnyObiektu + ` FROM obiekt_dokumentu_studio
	                              WHERE dokument_id = ? AND rodzaj = ?
	                              ORDER BY zakotwiczenie_pozycja, warstwa, id`

	postacUsunObiekt = `DELETE FROM obiekt_dokumentu_studio WHERE identyfikator_zewnetrzny = ?`

	postacKolumnyAparatu = `id, identyfikator_zewnetrzny, dokument_id, rodzaj, kotwica_od,
	                        kotwica_do, numer, etykieta, tresc, cel_kod, cel_adres,
	                        kolejnosc, nieswiezy, dane_json, utworzono, zaktualizowano`

	postacZapiszAparat = `INSERT INTO element_aparatu_studio
	                      (identyfikator_zewnetrzny, dokument_id, rodzaj, kotwica_od, kotwica_do,
	                       numer, etykieta, tresc, cel_kod, cel_adres, kolejnosc, nieswiezy, dane_json)
	                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          rodzaj = excluded.rodzaj,
	                          kotwica_od = excluded.kotwica_od,
	                          kotwica_do = excluded.kotwica_do,
	                          numer = excluded.numer,
	                          etykieta = excluded.etykieta,
	                          tresc = excluded.tresc,
	                          cel_kod = excluded.cel_kod,
	                          cel_adres = excluded.cel_adres,
	                          kolejnosc = excluded.kolejnosc,
	                          nieswiezy = excluded.nieswiezy,
	                          dane_json = excluded.dane_json,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	postacPobierzAparat = `SELECT ` + postacKolumnyAparatu + ` FROM element_aparatu_studio
	                       WHERE identyfikator_zewnetrzny = ?`

	// Kolejność wykazu jest kolejnością NUMERACJI: przypis wstawiony przed innym
	// ma przenumerować oba, a numer bierze się z miejsca w treści.
	postacListaAparatu = `SELECT ` + postacKolumnyAparatu + ` FROM element_aparatu_studio
	                      WHERE dokument_id = ?
	                      ORDER BY rodzaj, kolejnosc, kotwica_od, id`

	postacListaAparatuRodzaju = `SELECT ` + postacKolumnyAparatu + ` FROM element_aparatu_studio
	                             WHERE dokument_id = ? AND rodzaj = ?
	                             ORDER BY kolejnosc, kotwica_od, id`

	postacUsunAparat = `DELETE FROM element_aparatu_studio WHERE identyfikator_zewnetrzny = ?`

	postacKolumnyPola = `id, identyfikator_zewnetrzny, dokument_id, rodzaj, kotwica, format,
	                     wyrazenie, nazwa_wlasciwosci, wartosc, nieswieze, utworzono, zaktualizowano`

	postacZapiszPole = `INSERT INTO pole_dokumentu_studio
	                    (identyfikator_zewnetrzny, dokument_id, rodzaj, kotwica, format,
	                     wyrazenie, nazwa_wlasciwosci, wartosc, nieswieze)
	                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                        rodzaj = excluded.rodzaj,
	                        kotwica = excluded.kotwica,
	                        format = excluded.format,
	                        wyrazenie = excluded.wyrazenie,
	                        nazwa_wlasciwosci = excluded.nazwa_wlasciwosci,
	                        wartosc = excluded.wartosc,
	                        nieswieze = excluded.nieswieze,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	postacListaPol = `SELECT ` + postacKolumnyPola + ` FROM pole_dokumentu_studio
	                  WHERE dokument_id = ? ORDER BY kotwica, id`

	postacListaPolRodzaju = `SELECT ` + postacKolumnyPola + ` FROM pole_dokumentu_studio
	                         WHERE dokument_id = ? AND rodzaj = ? ORDER BY kotwica, id`

	postacUsunPole = `DELETE FROM pole_dokumentu_studio WHERE identyfikator_zewnetrzny = ?`
)

// ZapiszObiektDokumentu zakłada obiekt osadzony albo nadpisuje zastany po
// identyfikatorze zewnętrznym, zwracając stan po zapisie.
func (r *repozytoriumStudia) ZapiszObiektDokumentu(ctx context.Context,
	obiekt ObiektDokumentuStudia) (ObiektDokumentuStudia, error) {

	if obiekt.Kod == "" || obiekt.DokumentID == 0 {
		return ObiektDokumentuStudia{}, fmt.Errorf("dane: obiekt dokumentu studio bez identyfikatora albo bez dokumentu")
	}
	if obiekt.Rodzaj == "" {
		return ObiektDokumentuStudia{}, fmt.Errorf("dane: obiekt dokumentu studio %q bez rodzaju", obiekt.Kod)
	}
	if obiekt.Zakotwiczenie == "" {
		obiekt.Zakotwiczenie = "paragraph"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, postacZapiszObiekt)
	if err != nil {
		return ObiektDokumentuStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, obiekt.Kod, obiekt.DokumentID, obiekt.Rodzaj,
		tekstDoKolumny(obiekt.Zrodlo), tekstDoKolumny(obiekt.ZasobKod),
		tekstDoKolumny(obiekt.DesignWezelKod), tekstDoKolumny(obiekt.BibliotekaPlikKod),
		tekstDoKolumny(obiekt.AdresZrodla), obiekt.Zakotwiczenie, obiekt.ZakotwiczeniePozycja,
		obiekt.Warstwa, tekstDoKolumny(obiekt.TekstZastepczy),
		tekstDoKolumny(obiekt.TekstWewnetrzny), tekstDoKolumny(obiekt.Podpis),
		tekstDoKolumny(obiekt.PostacJSON))
	if err != nil {
		return ObiektDokumentuStudia{}, fmt.Errorf("dane: nie można zapisać obiektu %q: %w", obiekt.Kod, err)
	}
	return r.ObiektDokumentu(ctx, obiekt.Kod)
}

// ObiektDokumentu zwraca obiekt osadzony o wskazanym kodzie, zwracając błąd
// ErrBrakWiersza, gdy obiekt nie istnieje.
func (r *repozytoriumStudia) ObiektDokumentu(ctx context.Context, kod string) (ObiektDokumentuStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, postacPobierzObiekt)
	if err != nil {
		return ObiektDokumentuStudia{}, err
	}
	obiekt, err := postacOdczytajObiekt(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ObiektDokumentuStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return ObiektDokumentuStudia{}, fmt.Errorf("dane: nieczytelny wiersz obiektu %q: %w", kod, err)
	}
	return obiekt, nil
}

// ObiektyDokumentu zwraca obiekty dokumentu w kolejności zakotwiczenia w
// treści; puste rodzaj znaczy wszystkie rodzaje.
func (r *repozytoriumStudia) ObiektyDokumentu(ctx context.Context, dokumentID int64,
	rodzaj string) ([]ObiektDokumentuStudia, error) {

	zapytanie, argumenty := postacListaObiektow, []any{dokumentID}
	if rodzaj != "" {
		zapytanie, argumenty = postacListaObiektowRodzaju, []any{dokumentID, rodzaj}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać obiektów dokumentu %d: %w", dokumentID, err)
	}
	defer wiersze.Close()

	lista := []ObiektDokumentuStudia{}
	for wiersze.Next() {
		obiekt, err := postacOdczytajObiekt(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz obiektu: %w", err)
		}
		lista = append(lista, obiekt)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt obiektów dokumentu %d: %w", dokumentID, err)
	}
	return lista, nil
}

// UsunObiektDokumentu usuwa obiekt osadzony o wskazanym kodzie, oddając
// informację, czy wiersz istniał.
func (r *repozytoriumStudia) UsunObiektDokumentu(ctx context.Context, kod string) (bool, error) {
	return r.postacUsunWiersz(ctx, postacUsunObiekt, kod, "obiektu")
}

// ZapiszElementAparatu zakłada element aparatu albo nadpisuje zastany po
// identyfikatorze zewnętrznym, zwracając stan po zapisie.
func (r *repozytoriumStudia) ZapiszElementAparatu(ctx context.Context,
	element ElementAparatuStudia) (ElementAparatuStudia, error) {

	if element.Kod == "" || element.DokumentID == 0 {
		return ElementAparatuStudia{}, fmt.Errorf("dane: element aparatu studio bez identyfikatora albo bez dokumentu")
	}
	if element.Rodzaj == "" {
		return ElementAparatuStudia{}, fmt.Errorf("dane: element aparatu studio %q bez rodzaju", element.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, postacZapiszAparat)
	if err != nil {
		return ElementAparatuStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, element.Kod, element.DokumentID, element.Rodzaj,
		element.KotwicaOd, element.KotwicaDo, tekstDoKolumny(element.Numer),
		tekstDoKolumny(element.Etykieta), tekstDoKolumny(element.Tresc),
		tekstDoKolumny(element.CelKod), tekstDoKolumny(element.CelAdres),
		element.Kolejnosc, liczbaLogiczna(element.Nieswiezy), tekstDoKolumny(element.DaneJSON))
	if err != nil {
		return ElementAparatuStudia{}, fmt.Errorf("dane: nie można zapisać elementu aparatu %q: %w", element.Kod, err)
	}
	return r.ElementAparatu(ctx, element.Kod)
}

// ElementAparatu zwraca element aparatu o wskazanym kodzie, zwracając błąd
// ErrBrakWiersza, gdy element nie istnieje.
func (r *repozytoriumStudia) ElementAparatu(ctx context.Context, kod string) (ElementAparatuStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, postacPobierzAparat)
	if err != nil {
		return ElementAparatuStudia{}, err
	}
	element, err := postacOdczytajAparat(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ElementAparatuStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return ElementAparatuStudia{}, fmt.Errorf("dane: nieczytelny wiersz elementu aparatu %q: %w", kod, err)
	}
	return element, nil
}

// ElementyAparatu zwraca aparat dokumentu w kolejności numeracji; puste
// rodzaj znaczy cały aparat dokumentu.
func (r *repozytoriumStudia) ElementyAparatu(ctx context.Context, dokumentID int64,
	rodzaj string) ([]ElementAparatuStudia, error) {

	zapytanie, argumenty := postacListaAparatu, []any{dokumentID}
	if rodzaj != "" {
		zapytanie, argumenty = postacListaAparatuRodzaju, []any{dokumentID, rodzaj}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać aparatu dokumentu %d: %w", dokumentID, err)
	}
	defer wiersze.Close()

	lista := []ElementAparatuStudia{}
	for wiersze.Next() {
		element, err := postacOdczytajAparat(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz elementu aparatu: %w", err)
		}
		lista = append(lista, element)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt aparatu dokumentu %d: %w", dokumentID, err)
	}
	return lista, nil
}

// UsunElementAparatu usuwa element aparatu. Przenumerowanie pozostałych należy
// do rdzenia — warstwa danych nie wie, który numer jest który.
func (r *repozytoriumStudia) UsunElementAparatu(ctx context.Context, kod string) (bool, error) {
	return r.postacUsunWiersz(ctx, postacUsunAparat, kod, "elementu aparatu")
}

// ZapiszPoleDokumentu zakłada pole obliczane albo nadpisuje zastane po
// identyfikatorze zewnętrznym, zwracając stan po zapisie.
func (r *repozytoriumStudia) ZapiszPoleDokumentu(ctx context.Context,
	pole PoleDokumentuStudia) (PoleDokumentuStudia, error) {

	if pole.Kod == "" || pole.DokumentID == 0 {
		return PoleDokumentuStudia{}, fmt.Errorf("dane: pole dokumentu studio bez identyfikatora albo bez dokumentu")
	}
	if pole.Rodzaj == "" {
		return PoleDokumentuStudia{}, fmt.Errorf("dane: pole dokumentu studio %q bez rodzaju", pole.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, postacZapiszPole)
	if err != nil {
		return PoleDokumentuStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, pole.Kod, pole.DokumentID, pole.Rodzaj, pole.Kotwica,
		tekstDoKolumny(pole.Format), tekstDoKolumny(pole.Wyrazenie),
		tekstDoKolumny(pole.NazwaWlasciwosci), tekstDoKolumny(pole.Wartosc),
		liczbaLogiczna(pole.Nieswieze))
	if err != nil {
		return PoleDokumentuStudia{}, fmt.Errorf("dane: nie można zapisać pola %q: %w", pole.Kod, err)
	}
	pola, err := r.PolaDokumentu(ctx, pole.DokumentID, "")
	if err != nil {
		return PoleDokumentuStudia{}, err
	}
	for _, zapisane := range pola {
		if zapisane.Kod == pole.Kod {
			return zapisane, nil
		}
	}
	return PoleDokumentuStudia{}, ErrBrakWiersza
}

// PolaDokumentu zwraca pola dokumentu w kolejności zakotwiczenia w treści;
// puste rodzaj znaczy wszystkie rodzaje.
func (r *repozytoriumStudia) PolaDokumentu(ctx context.Context, dokumentID int64,
	rodzaj string) ([]PoleDokumentuStudia, error) {

	zapytanie, argumenty := postacListaPol, []any{dokumentID}
	if rodzaj != "" {
		zapytanie, argumenty = postacListaPolRodzaju, []any{dokumentID, rodzaj}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać pól dokumentu %d: %w", dokumentID, err)
	}
	defer wiersze.Close()

	lista := []PoleDokumentuStudia{}
	for wiersze.Next() {
		pole, err := postacOdczytajPole(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz pola dokumentu: %w", err)
		}
		lista = append(lista, pole)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt pól dokumentu %d: %w", dokumentID, err)
	}
	return lista, nil
}

// UsunPoleDokumentu usuwa pole dokumentu o wskazanym kodzie zewnętrznym,
// oddając informację, czy wiersz istniał.
func (r *repozytoriumStudia) UsunPoleDokumentu(ctx context.Context, kod string) (bool, error) {
	return r.postacUsunWiersz(ctx, postacUsunPole, kod, "pola dokumentu")
}

// postacOdczytajObiekt składa obiekt osadzony z jednego wiersza wyniku
// zapytania, zamieniając kolumny nullowalne na wskaźniki.
func postacOdczytajObiekt(wiersz skaner) (ObiektDokumentuStudia, error) {
	var obiekt ObiektDokumentuStudia
	var zrodlo, zasob, wezel, plik, adres, zastepczy, wewnetrzny, podpis, postac sql.NullString
	err := wiersz.Scan(&obiekt.ID, &obiekt.Kod, &obiekt.DokumentID, &obiekt.Rodzaj,
		&zrodlo, &zasob, &wezel, &plik, &adres, &obiekt.Zakotwiczenie,
		&obiekt.ZakotwiczeniePozycja, &obiekt.Warstwa, &zastepczy, &wewnetrzny, &podpis,
		&postac, &obiekt.Utworzono, &obiekt.Zaktualizowano)
	if err != nil {
		return ObiektDokumentuStudia{}, err
	}
	obiekt.Zrodlo = tekstZKolumny(zrodlo)
	obiekt.ZasobKod = tekstZKolumny(zasob)
	obiekt.DesignWezelKod = tekstZKolumny(wezel)
	obiekt.BibliotekaPlikKod = tekstZKolumny(plik)
	obiekt.AdresZrodla = tekstZKolumny(adres)
	obiekt.TekstZastepczy = tekstZKolumny(zastepczy)
	obiekt.TekstWewnetrzny = tekstZKolumny(wewnetrzny)
	obiekt.Podpis = tekstZKolumny(podpis)
	obiekt.PostacJSON = tekstZKolumny(postac)
	return obiekt, nil
}

// postacOdczytajAparat składa element aparatu z jednego wiersza wyniku
// zapytania, zamieniając kolumny nullowalne na wskaźniki.
func postacOdczytajAparat(wiersz skaner) (ElementAparatuStudia, error) {
	var element ElementAparatuStudia
	var numer, etykieta, tresc, cel, adres, dane sql.NullString
	var nieswiezy int64
	err := wiersz.Scan(&element.ID, &element.Kod, &element.DokumentID, &element.Rodzaj,
		&element.KotwicaOd, &element.KotwicaDo, &numer, &etykieta, &tresc, &cel, &adres,
		&element.Kolejnosc, &nieswiezy, &dane, &element.Utworzono, &element.Zaktualizowano)
	if err != nil {
		return ElementAparatuStudia{}, err
	}
	element.Numer = tekstZKolumny(numer)
	element.Etykieta = tekstZKolumny(etykieta)
	element.Tresc = tekstZKolumny(tresc)
	element.CelKod = tekstZKolumny(cel)
	element.CelAdres = tekstZKolumny(adres)
	element.Nieswiezy = nieswiezy != 0
	element.DaneJSON = tekstZKolumny(dane)
	return element, nil
}

// postacOdczytajPole składa pole dokumentu z jednego wiersza wyniku
// zapytania, zamieniając kolumny nullowalne na wskaźniki.
func postacOdczytajPole(wiersz skaner) (PoleDokumentuStudia, error) {
	var pole PoleDokumentuStudia
	var format, wyrazenie, wlasciwosc, wartosc sql.NullString
	var nieswieze int64
	err := wiersz.Scan(&pole.ID, &pole.Kod, &pole.DokumentID, &pole.Rodzaj, &pole.Kotwica,
		&format, &wyrazenie, &wlasciwosc, &wartosc, &nieswieze,
		&pole.Utworzono, &pole.Zaktualizowano)
	if err != nil {
		return PoleDokumentuStudia{}, err
	}
	pole.Format = tekstZKolumny(format)
	pole.Wyrazenie = tekstZKolumny(wyrazenie)
	pole.NazwaWlasciwosci = tekstZKolumny(wlasciwosc)
	pole.Wartosc = tekstZKolumny(wartosc)
	pole.Nieswieze = nieswieze != 0
	return pole, nil
}
